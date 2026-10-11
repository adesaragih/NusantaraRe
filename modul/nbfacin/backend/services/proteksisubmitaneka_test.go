package services

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
)

// lokasiSubmitAneka - satu lokasi IDR: okupasi 1 -> objek A (coverage sah ber-deductible, spreading OR 100 % premi = premi RNM);
// objek Golf G (USD) satu coverage.
func lokasiSubmitAneka() []models.LokasiAneka {
	return []models.LokasiAneka{{ObjectNo: "1", Okupasi: []models.OkupasiAneka{{OccupationName: "OKU", Objek: []models.ObjekAneka{{
		ObjectName: "A", Currency: "IDR", Coverages: []models.CoverageAneka{{OldID: "1", Coverage: "C1", CoverageNote: "UJI COV",
			PremiNusantaraRe: d("4"), Deductibles: []models.Deductible{{TypeDeductible: "1"}},
			Spreading: []models.BarisSpreading{{TreatyType: "OR", SharePercentage: d("100"), PremiumSpreaded: d("4")}}}}}}}},
		ObjekGolf: []models.ObjekAneka{{ObjectName: "G", Currency: "USD", Coverages: []models.CoverageAneka{{OldID: "9", Coverage: "G1",
			CoverageNote: "UJI GOLF", PremiNusantaraRe: d("2"),
			Spreading: []models.BarisSpreading{{TreatyType: "OR", SharePercentage: d("100"), PremiumSpreaded: d("2")}}}}}}}}
}

// layananSubmitAneka - layananSubmit (proteksi tiket 57) + case ber-BUSINESS (003, old) dan objek Aneka lokasi; objek FIRE =
// satu lokasi tanpa okupasi / item (baris T_LOCATIONLIST yang sama dengan lokasi Aneka).
func layananSubmitAneka(st *submitTiruan, old string, lokasi []models.LokasiAneka, ubah func(*models.Kasus)) *Service {
	svc := layananSubmit(st, models.ObjekFire{ObjectNo: "1"})
	k := &models.Kasus{CaseID: "UJI-NB-1", Opportunity: models.Opportunity{InsuredID: "UJI-ORG", ClassOfBusiness: "UJI KELAS",
		GroupBusinessID: "3"}, General: models.General{MarketingID: "UJI-MO", ReffNumber: "UJI/SLIP/1",
		StartDateTime: "20261001T170000.000 GMT", EndDateTime: "20271001T170000.000 GMT"}}
	if ubah != nil {
		ubah(k)
	}
	return svc.DenganKasus(kasusTiruan{ada: map[string]*models.Kasus{"UJI-NB-1": k}}).
		DenganTableOfLimit(tolTiruan{bisnis: satuBisnis("003", old, "10048")}).
		DenganObjekAneka(anekaTiruan{ada: map[string][]models.LokasiAneka{"UJI-NB-1": lokasi}})
}

// TestProteksiAcceptAnekaTanpaFire - K73-28: case Aneka sah tanpa Scoring Risk lolos; tidak ada pesan FIRE (okupasi / Province /
// Loss record / Scoring) walau baris lokasi FIRE kosong. Case FIRE (lini tak terbaca) tetap mendapat pesan FIRE.
func TestProteksiAcceptAnekaTanpaFire(t *testing.T) {
	ctx := context.Background()
	pesan, err := layananSubmitAneka(submitUji(), "30", lokasiSubmitAneka(), nil).ProteksiAccept(ctx, "UJI-NB-1", nil)
	if err != nil || len(pesan) != 0 {
		t.Fatalf("Aneka sah: %v %q", err, pesan)
	}
	fire, err := layananSubmit(submitUji(), models.ObjekFire{ObjectNo: "1"}).ProteksiAccept(ctx, "UJI-NB-1", nil)
	if err != nil || !strings.Contains(strings.Join(fire, "|"), "Occupation on location 1 can't be null!") {
		t.Errorf("FIRE tanpa okupasi: %v %q", err, fire)
	}
}

// TestProteksiCoverageAneka - ProtectCoverage_Act langkah 5 / 6, pesan verbatim + flag: coverage kosong, mata uang coverage
// (kosong -> objek), deductible hanya bila Binding flag 2 dan bukan IsTBonding, mata uang objek; Golf coverage / mata uang.
func TestProteksiCoverageAneka(t *testing.T) {
	ctx := context.Background()
	lokasi := lokasiSubmitAneka()
	o := &lokasi[0].Okupasi[0].Objek[0]
	o.Currency = ""
	o.Coverages = append(o.Coverages, models.CoverageAneka{OldID: "2", Coverage: "C2"},
		models.CoverageAneka{OldID: "3", Coverage: "C3", CoverageNote: "N", Currency: "USD", Deductibles: []models.Deductible{{}}})
	lokasi[0].ObjekGolf[0].Currency = ""
	lokasi[0].ObjekGolf[0].Coverages[0].CoverageNote = ""
	isi := isianProteksiAneka{lini: liniSpreadingAneka{aneka: true, golf: true}, deductible: true}
	var p pengumpulPesan
	if err := Baru(nil).proteksiCoverageAneka(ctx, lokasi, isi, &p); err != nil {
		t.Fatal(err)
	}
	mau := []string{
		" Currency on location 1 object item 1 Coverage 1 can't be empty!",
		" Coverage on location 1 can't be empty!",
		" Currency on location 1 object item 1 Coverage 2 can't be empty!",
		" Deductible on location 1 object item 1 Coverage 2 should be added!",
		" Deductible on location 1 Coverage 3  DeductibleList 1 can't be empty!",
		" Currency on location 1 can't be empty!",
	}
	if !reflect.DeepEqual(p.pesan, mau) {
		t.Fatalf("pesan:\n%q\nmau\n%q", p.pesan, mau)
	}
	if !reflect.DeepEqual(p.flag[3], []flagProteksi{flagDeductible4}) || !reflect.DeepEqual(p.flag[4], []flagProteksi{flagDeductible2}) ||
		!reflect.DeepEqual(p.flag[0], []flagProteksi{flagMarine2}) {
		t.Errorf("flag %v", p.flag)
	}
	isi.deductible = false
	var q pengumpulPesan
	_ = Baru(nil).proteksiCoverageAneka(ctx, lokasi, isi, &q)
	if strings.Contains(strings.Join(q.pesan, "|"), "Deductible") {
		t.Errorf("deductible di luar Binding: %q", q.pesan)
	}
	isi.lini = liniSpreadingAneka{golf: true}
	var g pengumpulPesan
	_ = Baru(nil).proteksiCoverageAneka(ctx, lokasi, isi, &g)
	if !reflect.DeepEqual(g.pesan, []string{" Coverage on location 1 can't be empty!", " Currency on location 1 can't be empty!"}) {
		t.Errorf("Golf: %q", g.pesan)
	}
}

// TestKasusGandaAneka - GetDataDoubleCase per lokasi: isian ParamInProdCase (InsuredName tanpa "." / spasi, BusinessCode,
// NoOfferSlip, periode dd/MM/yyyy WIB, case sendiri dikecualikan K73-31); ditemukan -> pesan lalu lompat (lokasi berikutnya
// dan Golf tidak diperiksa).
func TestKasusGandaAneka(t *testing.T) {
	ctx := context.Background()
	st := submitUji()
	k := st.konteks["UJI-NB-1"]
	k.InsuredName = "P.T. UJI  JAYA"
	st.konteks["UJI-NB-1"] = k
	st.ganda = []string{"ASM-FW-GISFW-WORK NB-777"}
	lokasi := append(lokasiSubmitAneka(), lokasiSubmitAneka()...)
	lokasi[1].Okupasi[0].Objek[0].Coverages[0].Coverage = ""
	svc := layananSubmitAneka(st, "30", lokasi, nil)
	pesan, err := svc.ProteksiAccept(ctx, "UJI-NB-1", nil)
	if err != nil || !reflect.DeepEqual(pesan, []string{"Double Input Case ASM-FW-GISFW-WORK NB-777"}) {
		t.Fatalf("pesan %v %q", err, pesan)
	}
	mau := models.CariKasusGanda{InsuredName: "PTUJIJAYA", SobName: "UJI SOB", BusinessCode: "10048", NoOfferSlip: "UJI/SLIP/1",
		Mulai: "02/10/2026", Akhir: "02/10/2027", IDPegaSendiri: idPega("UJI-NB-1")}
	if len(st.cariGanda) != 1 || st.cariGanda[0] != mau {
		t.Errorf("isian: %+v", st.cariGanda)
	}
}

// TestDeductibleAnekaTahapBinding - deductible ProtectCoverage Aneka hanya tahap Binding (flag 2) dan bukan IsTBonding (team 5).
func TestDeductibleAnekaTahapBinding(t *testing.T) {
	ctx := context.Background()
	lokasi := lokasiSubmitAneka()
	lokasi[0].Okupasi[0].Objek[0].Coverages[0].Deductibles = nil
	for team, ada := range map[string]bool{"1": true, "5": false} {
		st := submitUji()
		k := st.konteks["UJI-NB-1"]
		k.FlagOnGoingPolicy, k.TeamGroup = "2", team
		st.konteks["UJI-NB-1"] = k
		svc := layananSubmitAneka(st, "30", lokasi, nil)
		p, err := svc.proteksiBerflag(ctx, "UJI-NB-1", nil, true)
		if err != nil {
			t.Fatal(err)
		}
		if got := strings.Contains(strings.Join(p.pesan, "|"), " Deductible on location 1 object item 1 Coverage 1 should be added!"); got != ada {
			t.Errorf("team %s: %q", team, p.pesan)
		}
	}
}

// TestDeductibleBedaMataUangAneka - K73-30 `[dugaan]`: deductible ber-Currency beda dari Currency coverage (kosong -> objek).
func TestDeductibleBedaMataUangAneka(t *testing.T) {
	lokasi := lokasiSubmitAneka()
	if deductibleBedaMataUangAneka(lokasi) {
		t.Fatal("tanpa mata uang deductible")
	}
	lokasi[0].Okupasi[0].Objek[0].Coverages[0].Deductibles[0].Currency = "IDR"
	if deductibleBedaMataUangAneka(lokasi) {
		t.Error("sama dengan objek")
	}
	lokasi[0].Okupasi[0].Objek[0].Coverages[0].Deductibles[0].Currency = "USD"
	if !deductibleBedaMataUangAneka(lokasi) {
		t.Error("beda")
	}
	lokasi[0].Okupasi[0].Objek[0].Coverages[0].Deductibles[0].Currency = "UNKNOWN"
	if deductibleBedaMataUangAneka(lokasi) {
		t.Error("UNKNOWN")
	}
}

// TestPesanSpreadingAneka - CheckSpreadingProtectAnekaGolf: Aneka (object item = nomor okupasi, skala 4) dan Golf (nomor objek
// diisi K73-32, skala 2).
func TestPesanSpreadingAneka(t *testing.T) {
	lokasi := lokasiSubmitAneka()
	if p := pesanSpreadingAneka(lokasi, liniSpreadingAneka{aneka: true, golf: true}); len(p) != 0 {
		t.Fatalf("sah: %q", p)
	}
	c := &lokasi[0].Okupasi[0].Objek[0].Coverages[0]
	c.Spreading = []models.BarisSpreading{{TreatyType: "OR", SharePercentage: d("60"), PremiumSpreaded: d("2.40006")},
		{SharePercentage: d("30"), PremiumSpreaded: d("1.6")}}
	g := &lokasi[0].ObjekGolf[0].Coverages[0]
	g.Spreading = []models.BarisSpreading{{TreatyType: "OR", SharePercentage: d("100"), PremiumSpreaded: d("2.004")}}
	mau := []string{
		"% spreading on location 1 object item 1 coverage 1 can't be empty!",
		"% Total spreading on location 1 occupation 1 object 1 coverage 1 must be 100%!",
		"Total premium spreaded on location 1 occupation 1 object 1 coverage 1 is not equals premium RNM!",
	}
	if p := pesanSpreadingAneka(lokasi, liniSpreadingAneka{aneka: true, golf: true}); !reflect.DeepEqual(p, mau) {
		t.Errorf("Aneka:\n%q", p)
	}
	g.Spreading[0].PremiumSpreaded = d("2.006")
	if p := pesanSpreadingAneka(lokasi, liniSpreadingAneka{golf: true}); !reflect.DeepEqual(p,
		[]string{"Total premium spreaded on location 1 object 1 coverage 1 is not equals premium RNM!"}) {
		t.Errorf("Golf skala 2: %q", p)
	}
	g.Spreading = nil
	if p := pesanSpreadingAneka(lokasi, liniSpreadingAneka{golf: true}); !reflect.DeepEqual(p, []string{
		"% Total spreading on location 1 object 1 coverage 1 must be 100%!",
		"Total premium spreaded on location 1 object 1 coverage 1 is not equals premium RNM!"}) {
		t.Errorf("Golf kosong: %q", p)
	}
	g.Spreading = []models.BarisSpreading{{SharePercentage: d("100"), PremiumSpreaded: d("2")}}
	if p := pesanSpreadingAneka(lokasi, liniSpreadingAneka{golf: true}); !reflect.DeepEqual(p,
		[]string{"% spreading on location 1 object item 1 coverage 1 can't be empty!"}) {
		t.Errorf("Golf treaty kosong: %q", p)
	}
}

// TestSpreadingUWAneka - tahap uw: CheckSpreadingProtectAnekaGolf ikut (case Aneka); case FIRE tidak membaca objek Aneka.
func TestSpreadingUWAneka(t *testing.T) {
	ctx := context.Background()
	lokasi := lokasiSubmitAneka()
	lokasi[0].Okupasi[0].Objek[0].Coverages[0].Spreading = nil
	svc := layananSubmitAneka(submitUji(), "30", lokasi, nil).DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{
		{PercentShare: d("50"), StartDateTime: "20261001T170000.000 GMT", Lokasi: []models.LokasiSpreading{{ObjectNo: "1"}}}}})
	pesan, err := svc.proteksiSpreadingUW(ctx, "UJI-NB-1")
	if err != nil || !reflect.DeepEqual(pesan, []string{"% Total spreading on location 1 occupation 1 object 1 coverage 1 must be 100%!",
		"Total premium spreaded on location 1 occupation 1 object 1 coverage 1 is not equals premium RNM!"}) {
		t.Errorf("Aneka: %v %q", err, pesan)
	}
}

// TestSpreadingBindingFire - K73-29 (regresi FIRE, keputusan work owner 09-10-2026): di tahap Binding CheckSpreadingProtect
// berjalan untuk FIRE juga - coverage tanpa baris spreading menahan Accept, tombol Binding, dan Banding.
func TestSpreadingBindingFire(t *testing.T) {
	ctx := context.Background()
	svc := layananBinding(submitUji(), "1", nil).DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{kasusDuaItem("50", "USD")}})
	k, err := svc.konteksSubmit(ctx, "UJI-NB-1")
	if err != nil {
		t.Fatal(err)
	}
	m := "% Total spreading on location 1 object item 1 coverage 1 must be 100%!"
	pesan, pesanBinding, err := svc.periksaBinding(ctx, "UJI-NB-1", k, models.EmailTypeAccept)
	if err != nil || !strings.Contains(strings.Join(pesan, "|"), m) || !strings.Contains(strings.Join(pesanBinding, "|"), m) {
		t.Errorf("Accept: %v %q / %q", err, pesan, pesanBinding)
	}
	banding, _, err := svc.periksaBinding(ctx, "UJI-NB-1", k, models.EmailTypeBanding)
	if err != nil || !strings.Contains(strings.Join(banding, "|"), m) {
		t.Errorf("Banding: %v %q", err, banding)
	}
	// Spreading sah (fixture layananBinding) -> tanpa pesan spreading.
	pesan, _, _ = layananBinding(submitUji(), "1", nil).periksaBinding(ctx, "UJI-NB-1", k, models.EmailTypeAccept)
	if strings.Contains(strings.Join(pesan, "|"), "spreading") {
		t.Errorf("sah: %q", pesan)
	}
}

// layananJenjangAneka - Submit tahap uw `posisi` untuk case Aneka (BUSINESS 003 / old): spreading FIRE tanpa item (PS 50), objek
// Aneka `lokasi`, tabel limit PRY (limitUji) + NONPROPANDENGG / ENGINEERINGG team 1 yang jenjangnya berbeda dari PRY.
func layananJenjangAneka(st *submitTiruan, posisi, old string, lokasi []models.LokasiAneka, kurs map[string]*apd.Decimal) *Service {
	sp := models.KasusSpreading{PercentShare: d("50"), StartDateTime: "20261001T170000.000 GMT", Lokasi: []models.LokasiSpreading{{ObjectNo: "1"}}}
	svc := layananSubmitUW(st, posisi, sp).DenganTableOfLimit(tolTiruan{oldID: []string{oldIDUji}, bisnis: satuBisnis("003", old)}).
		DenganObjekAneka(anekaTiruan{ada: map[string][]models.LokasiAneka{"UJI-NB-1": lokasi}}).
		DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{sp}, idUang: map[string]string{"IDR": idMataUangIDR, "USD": "10001"}, kurs: kurs})
	svc.limit = limitTiruan{a: append(append([]models.BarisLimitA{}, limitUji...),
		barisLimit("M_LIMIT_NONPROPANDENGG", "1", "JUW_B", "0"), barisLimit("M_LIMIT_NONPROPANDENGG", "1", "JUW_A", "1"),
		barisLimit("M_LIMIT_ENGINEERINGG", "1", "JUW_B", "0"), barisLimit("M_LIMIT_ENGINEERINGG", "1", "SENIORUW", "1"))}
	return svc
}

// lokasiJenjangAneka - lokasiSubmitAneka dengan TSI 1000 IDR (coverage pertama) dan spreading ORS (lolos jenis treaty).
func lokasiJenjangAneka() []models.LokasiAneka {
	l := lokasiSubmitAneka()
	setiapCoverageAneka(l, true, true, func(_ int, _ *models.ObjekAneka, _ int, c *models.CoverageAneka) {
		c.TSI, c.TSILiability = d("1000"), d("1000")
		for i := range c.Spreading {
			c.Spreading[i].TreatyType = treatyORS
		}
	})
	return l
}

// TestJenjangAneka - K73 P4 butir 4: JUW B Accept case Aneka memakai M_LIMIT_NONPROPANDENGG (langkah 15) -> JUW_A; MBD
// (IsEngineering) memakai M_LIMIT_ENGINEERINGG (langkah 14) -> SENIORUW; PRY (FIRE) akan memberi UNDERWRITER.
func TestJenjangAneka(t *testing.T) {
	ctx := context.Background()
	for old, mau := range map[string]string{"30": posisiJuniorUWA, "10": posisiSeniorUW} {
		st := submitUji()
		svc := layananJenjangAneka(st, posisiJuniorUWB, old, lokasiJenjangAneka(), nil)
		h, err := svc.AjukanSubmit(ctx, pelakuPosisi(posisiJuniorUWB), "UJI-NB-1", models.IsiSubmit{EmailType: "1", Comment: "x"})
		if err != nil || h.Position != mau {
			t.Errorf("BUSINESS %s: %v %+v", old, err, h)
		}
	}
}

// TestCarid2Aneka - CountTotalTSIPremiNusaRe_Act langkah 6 / 7: Aneka = CoverageList(1).TSI × kurs objek, top risk =
// CoverageList(1).TSILiability × kurs (K73-33); Golf = SETIAP coverage .TSI × kurs, tanpa top risk.
func TestCarid2Aneka(t *testing.T) {
	ctx := context.Background()
	kurs := map[string]*apd.Decimal{"10001": d("2")}
	hitung := func(old string, ubah func([]models.LokasiAneka)) string {
		l := lokasiJenjangAneka()
		ubah(l)
		svc := layananJenjangAneka(submitUji(), posisiJuniorUWB, old, l, kurs)
		d, err := svc.bacaDasarJenjang(ctx, "UJI-NB-1")
		if err != nil {
			t.Fatal(err)
		}
		if d.pesanKurs != "" {
			return d.pesanKurs
		}
		return utils.FormatDecimal(d.carid2())
	}
	if c := hitung("30", func([]models.LokasiAneka) {}); c != "500" {
		t.Errorf("Aneka IDR: %s", c)
	}
	usd := func(l []models.LokasiAneka) {
		o := &l[0].Okupasi[0].Objek[0]
		o.Currency = "USD"
		o.Coverages = append(o.Coverages, models.CoverageAneka{OldID: "2", TSI: d("999999")}) // bukan coverage pertama: tidak dihitung
	}
	if c := hitung("30", usd); c != "1000" {
		t.Errorf("Aneka USD: %s", c)
	}
	if c := hitung("30", func(l []models.LokasiAneka) {
		usd(l)
		l[0].IsTopRisk = true
		l[0].Okupasi[0].Objek[0].Coverages[0].TSILiability = d("300")
	}); c != "300" {
		t.Errorf("Aneka top risk berkurs: %s", c)
	}
	golf := func(l []models.LokasiAneka) {
		g := &l[0].ObjekGolf[0]
		g.Coverages = append(g.Coverages, models.CoverageAneka{OldID: "8", TSI: d("500")})
	}
	if c := hitung("27", golf); c != "1500" { // USD 2 × (1000 + 500) × 50 %
		t.Errorf("Golf: %s", c)
	}
	if c := hitung("27", func(l []models.LokasiAneka) { golf(l); l[0].ObjekGolf[0].Coverages[1].Currency = "IDR" }); c != "1250" {
		t.Errorf("Golf coverage IDR kurs 1: %s", c)
	}
}

// TestFacRetroAneka - K73-36: FACOUT (10015) di spreading objek mana pun - termasuk Golf di lokasi BUKAN terakhir - = IsFacRetro.
func TestFacRetroAneka(t *testing.T) {
	ctx := context.Background()
	l := append(lokasiJenjangAneka(), lokasiJenjangAneka()...)
	l[0].ObjekGolf[0].Coverages[0].Spreading[0].TreatyType = treatyFACOUT
	d, err := layananJenjangAneka(submitUji(), posisiJuniorUWB, "27", l, nil).bacaDasarJenjang(ctx, "UJI-NB-1")
	if err != nil || !d.adaFacRetro() {
		t.Errorf("Golf lokasi 1: %v %v", err, d.adaFacRetro())
	}
	d, _ = layananJenjangAneka(submitUji(), posisiJuniorUWB, "30", l, nil).bacaDasarJenjang(ctx, "UJI-NB-1")
	if d.adaFacRetro() {
		t.Error("Aneka tanpa FACOUT di objek okupasi")
	}
}

// TestDeductibleAcceptanceTLAneka - Protection_Act 9-10 untuk case Aneka: IsBedaCurr dari deductible Aneka (K73-30), bukan
// objek FIRE.
func TestDeductibleAcceptanceTLAneka(t *testing.T) {
	ctx := context.Background()
	for cur, ada := range map[string]bool{"USD": true, "IDR": false} {
		l := lokasiSubmitAneka()
		l[0].Okupasi[0].Objek[0].Coverages[0].Deductibles[0].Currency = cur
		p, err := layananSubmitAneka(submitUji(), "30", l, nil).proteksiBerflag(ctx, "UJI-NB-1", &isianDeductibleTL{}, false)
		if err != nil {
			t.Fatal(err)
		}
		if got := strings.Contains(strings.Join(p.pesan, "|"), "Deductible Acceptance cannot be empty"); got != ada {
			t.Errorf("%s: %q", cur, p.pesan)
		}
	}
}

// TestCarid2AnekaTSIBukanLiability - langkah 6: CARI22 = CoverageList(1).TSI (bukan TSILiability) di luar top risk.
func TestCarid2AnekaTSIBukanLiability(t *testing.T) {
	l := lokasiJenjangAneka()
	l[0].Okupasi[0].Objek[0].Coverages[0].TSILiability = d("777")
	d, err := layananJenjangAneka(submitUji(), posisiJuniorUWB, "30", l, nil).bacaDasarJenjang(context.Background(), "UJI-NB-1")
	if err != nil || utils.FormatDecimal(d.carid2()) != "500" {
		t.Errorf("%v %v", err, d.carid2())
	}
}

// TestKasusGandaMelompatiGolf - TRANS L6372 T=1 -> AA: sesudah "Double Input Case", blok GOLF (langkah 6) tidak dijalankan.
func TestKasusGandaMelompatiGolf(t *testing.T) {
	st := submitUji()
	st.ganda = []string{"NB-1"}
	l := lokasiSubmitAneka()
	l[0].ObjekGolf[0].Currency = ""
	var p pengumpulPesan
	err := Baru(nil).DenganSubmit(st).proteksiCoverageAneka(context.Background(), l, isianProteksiAneka{
		lini: liniSpreadingAneka{aneka: true, golf: true}, cariGanda: &models.CariKasusGanda{}}, &p)
	if err != nil || !reflect.DeepEqual(p.pesan, []string{"Double Input Case NB-1"}) {
		t.Errorf("%v %q", err, p.pesan)
	}
}

// TestJenjangLanjutAneka - ActFlow (Senior) case Aneka: When COB tangga membaca BusinessType case (bukan "Fire") ->
// M_LIMIT_NONPROPANDENGG: SENIORUW 1 -> DEPHEADUNDERWRITER 2 (< TSI 500); tabel PRY (SENIORUW 300, DEPHEAD 1000) -> selesai.
func TestJenjangLanjutAneka(t *testing.T) {
	st := submitUji()
	svc := layananJenjangAneka(st, posisiSeniorUW, "30", lokasiJenjangAneka(), nil)
	svc.limit = limitTiruan{a: append(append([]models.BarisLimitA{}, limitUji...),
		barisLimit("M_LIMIT_NONPROPANDENGG", "1", "SENIORUW", "1"), barisLimit("M_LIMIT_NONPROPANDENGG", "1", "DEPHEADUNDERWRITER", "2"))}
	h, err := svc.AjukanSubmit(context.Background(), pelakuPosisi(posisiSeniorUW), "UJI-NB-1", models.IsiSubmit{EmailType: "1", Comment: "x"})
	if err != nil || h.Position != posisiDepHeadUW {
		t.Errorf("%v %+v", err, h)
	}
}
