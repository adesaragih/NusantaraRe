package services

import (
	"context"
	"errors"
	"reflect"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
)

// Tiket 69a-1 - jalur Policy: TL + jenjang UW tahap Policy (FlagOnGoingPolicy 1).

// spreadingPolicy - satu coverage PS 50, TSILiability `tsi` (TotalTSINusaRe = tsi / 2), satu baris spreading 100 % treaty itu.
func spreadingPolicy(tsi, tipe, nama string) models.KasusSpreading {
	k := spreadingTanpaFacRetro(tsi)
	k.Lokasi[0].Items[0].Coverages[0].Spreading = []models.BarisSpreading{{TreatyType: tipe, TreatyName: nama, SharePercentage: d("100")}}
	hitungSaveSpreading(&k, d("50"))
	return k
}

// limitPolicy - tabel PRY team 1 dengan LEADER (GetLimitAkseptasi_JUW_UW langkah 7).
var limitPolicy = []models.BarisLimitA{barisLimit(tabelPRY, "1", "LEADER", "0"), barisLimit(tabelPRY, "1", "JUW_A", "0"),
	barisLimit(tabelPRY, "1", "UNDERWRITER", "1"), barisLimit(tabelPRY, "1", "SENIORUW", "300"),
	barisLimit(tabelPRY, "1", "DEPHEADUNDERWRITER", "1000")}

// layananPolicy - kasus UJI-NB-1 di posisi `posisi` FLAG 1, spreading sp, tabel limitPolicy, ACCUMULATIONCODE "UJI-AKU"
// (lokasiSah) berjalan sebanyak `akum`; General / cedant terisi untuk selesai Policy (tiket 69a-2).
func layananPolicy(st *submitTiruan, posisi string, sp models.KasusSpreading, akum int) *Service {
	svc := layananSubmitUW(st, posisi, sp)
	k := st.konteks["UJI-NB-1"]
	k.FlagOnGoingPolicy, k.ConfirmBinding, k.MarketingCode = "1", "1", "UJI-MC"
	st.konteks["UJI-NB-1"] = k
	st.akumulasi = map[string]int{"UJI-AKU": akum}
	svc.limit = limitTiruan{a: limitPolicy}
	g := &svc.kasus.(kasusTiruan).ada["UJI-NB-1"].General
	g.StartDateTime, g.EndDateTime, g.ReffNumber, g.QQName = "20261019T170000.000 GMT", "20271019T170000.000 GMT", "UJI SLIP", "UJI QQ"
	g.PolicyType, g.SourceOfBusinessID, g.CedingCoName, g.GroupName = "2", "UJI-SOB", "UJI CEDING A;UJI CEDING B", "UJI GRUP"
	g.CedingList = []models.Ceding{{ID: "C1", Name: "UJI CEDING A"}, {ID: "C2", Name: "UJI CEDING B"}}
	return svc.DenganCedant(&cedantTiruan{ada: models.KasusCedant{ShareCedantType: "1", Cedant: []models.BarisCedant{
		{CedingCo: "C1;", CedingCoName: "UJI CEDING A;", ShareCeding: d("60")}, {CedingCo: "C2", CedingCoName: "UJI CEDING B", ShareCeding: d("40")}}}})
}

// TestTahapPolicy - TL flag 1 = teamLeaderPolicy (tahapDari: tab dapat diubah TL); UW / Senior / DepHead / UW Financial / JUW A flag 1
// = uwPolicy; JUW B flag 1 / flag 2 bukan; kunciPeriode keduanya.
func TestTahapPolicy(t *testing.T) {
	ctx := context.Background()
	for posisi, mau := range map[string]string{posisiTeamLeader: TahapTeamLeaderPolicy, posisiUWUtama: TahapUWPolicy, posisiSeniorUW: TahapUWPolicy,
		posisiDepHeadUW: TahapUWPolicy, posisiUWFinancial: TahapUWPolicy, posisiJuniorUWA: TahapUWPolicy, posisiJuniorUWB: "",
		"ReasFacInMarketing": ""} {
		if got := tahapKonteks(models.KonteksSubmit{Position: posisi, FlagOnGoingPolicy: "1"}); got != mau {
			t.Errorf("%s: %q mau %q", posisi, got, mau)
		}
		svc := Baru(nil).DenganTahap(tahapTiruan{"K": {Position: posisi, FlagOnGoingPolicy: "1"}})
		if got := svc.KunciPeriode(ctx, "K"); got != (mau != "") {
			t.Errorf("%s kunciPeriode %v", posisi, got)
		}
	}
	if tahapKonteks(models.KonteksSubmit{Position: posisiSeniorUW, FlagOnGoingPolicy: "2"}) != "" ||
		tahapDari(models.TahapKasus{Position: posisiTeamLeader, FlagOnGoingPolicy: "1"}) != "" || // K69-8: tab baca-saja
		tahapDari(models.TahapKasus{Position: posisiSeniorUW, FlagOnGoingPolicy: "1"}) != "" {
		t.Error("tahap lain")
	}
	if !reflect.DeepEqual(emailTypeTahap[TahapUWPolicy], []string{"1", "7"}) ||
		!reflect.DeepEqual(emailTypeTahap[TahapTeamLeaderPolicy], []string{"1", "2", "3", "7"}) {
		t.Errorf("email type %v / %v", emailTypeTahap[TahapUWPolicy], emailTypeTahap[TahapTeamLeaderPolicy])
	}
}

// TestIsFlagUW - K69-5: GetCountAccumulation > 0 (argumen: kode, StartDateTime dd/MM/yyyy WIB, pzInsKey, pyID) atau spreading
// 10007 / 10015 / TreatyName "SPL"; selain itu 0; StartDateTime kosong -> tanpa kueri.
func TestIsFlagUW(t *testing.T) {
	ctx := context.Background()
	for nama, c := range map[string]struct {
		sp   models.KasusSpreading
		akum int
		mau  bool
	}{
		"QS tanpa akumulasi": {spreadingPolicy("1000", "10001", "QS UJI"), 0, false},
		"akumulasi":          {spreadingPolicy("1000", "10001", "QS UJI"), 2, true},
		"ORS":                {spreadingPolicy("1000", "10007", "ORS"), 0, true},
		"FACOUT":             {spreadingPolicy("1000", "10015", "FACOUT"), 0, true},
		"SPL":                {spreadingPolicy("1000", "10002", "SPL UJI"), 0, true},
		"spl huruf kecil":    {spreadingPolicy("1000", "10002", "spl uji"), 0, false},
	} {
		st := submitUji()
		svc := layananPolicy(st, posisiTeamLeader, c.sp, c.akum)
		got, err := svc.isFlagUW(ctx, "UJI-NB-1")
		if err != nil || got != c.mau {
			t.Errorf("%s: %v %v", nama, got, err)
		}
		if nama == "akumulasi" && !reflect.DeepEqual(st.argAkumulasi, []string{"UJI-AKU", "02/01/2025", "ASM-FW-GISFW-WORK UJI-NB-1", "UJI-NB-1"}) {
			t.Errorf("argumen GetCountAccumulation %q", st.argAkumulasi)
		}
	}
	st := submitUji()
	sp := spreadingPolicy("1000", "10001", "QS UJI")
	sp.StartDateTime = ""
	if got, err := layananPolicy(st, posisiTeamLeader, sp, 5).isFlagUW(ctx, "UJI-NB-1"); err != nil || got || st.argAkumulasi != nil {
		t.Errorf("tanpa StartDateTime: %v %v %q", got, err, st.argAkumulasi)
	}
}

// TestAjukanSubmitTLPolicy - TL tahap Policy: Accept IsFlagUW 1 -> jenjang LEADER (posisi, flag tetap, riwayat Policy / 2, tanpa
// DATE_TO_UW / performance); IsFlagUW 0 / Else -> selesai Policy (tiket 69a-2); Else + FacRetro -> 409 FacRetro; group -> Else;
// Reject / Ask -> Marketing Binding flag 2 (ConfirmBinding tidak direset); Decline -> Resolved-Decline; proteksi 422; marketing 403.
func TestAjukanSubmitTLPolicy(t *testing.T) {
	ctx := context.Background()
	qs := spreadingPolicy("1000", "10001", "QS UJI")
	accept := models.IsiSubmit{EmailType: "1", Comment: "UJI setuju"}
	st := submitUji()
	h, err := layananPolicy(st, posisiTeamLeader, qs, 1).AjukanSubmit(ctx, pelakuTL, "UJI-NB-1", accept)
	if err != nil || h.Position != posisiUWUtama || h.NBStatus != "NB IS IN UNDERWRITER 'S INBOX" || len(st.tulis) != 1 || h.NomorPolis != "" {
		t.Fatalf("accept jenjang: %v %+v", err, h)
	}
	w := st.tulis[0]
	if w.FlagBaru != "" || w.Riwayat.Posisi != "Policy" || w.Riwayat.Putaran != "2" || w.Riwayat.Approval != "Accept" ||
		w.Pega.Workbasket != posisiTeamLeader || w.KeUW != "" || w.Deductible == nil || w.Performance != nil || w.Validitas != nil ||
		w.StatusBinding != nil || w.Binding != nil {
		t.Errorf("tulisan accept: %+v", w)
	}
	for nama, c := range map[string]struct {
		sp   models.KasusSpreading
		akum int
		ubah func(*Service)
		mau  error // nil = selesai Policy
	}{
		"IsFlagUW 0":           {qs, 0, nil, nil},
		"tanpa LEADER":         {qs, 1, func(s *Service) { s.limit = limitTiruan{a: limitUji} }, nil},
		"Else FacRetro":        {spreadingPolicy("1000", "10015", "FACOUT"), 0, func(s *Service) { s.limit = limitTiruan{a: limitUji} }, ErrFacRetroBelumTersedia},
		"TSI di bawah jenjang": {spreadingPolicy("1", "10007", "ORS"), 0, nil, nil},
		"group":                {qs, 1, func(s *Service) { s.DenganGrup(KonfigGrup{Login: []string{"UJI-TL"}}, false) }, nil},
	} {
		st := submitUji()
		svc := layananPolicy(st, posisiTeamLeader, c.sp, c.akum)
		if c.ubah != nil {
			c.ubah(svc)
		}
		h, err := svc.AjukanSubmit(ctx, pelakuTL, "UJI-NB-1", accept)
		switch {
		case c.mau != nil:
			if !errors.Is(err, c.mau) || len(st.tulis) != 0 {
				t.Errorf("%s: %v (%d tulis)", nama, err, len(st.tulis))
			}
		case err != nil || h.Position != posisiTeamLeader || len(st.tulis) != 1 || st.tulis[0].Policy == nil ||
			st.tulis[0].StatusWorkBaru != statusSelesaiPolicy || !st.tulis[0].NBStatusTetap || st.tulis[0].PosisiBaru != posisiTeamLeader:
			t.Errorf("%s selesai: %v %+v %d", nama, err, h, len(st.tulis))
		}
		if _, err := svc.PeriksaSubmit(ctx, pelakuTL, "UJI-NB-1", IsiPeriksa{EmailType: "1"}); !errors.Is(err, c.mau) {
			t.Errorf("%s periksa: %v", nama, err)
		}
	}
	for _, e := range []string{"2", "3"} {
		st := submitUji()
		h, err := layananPolicy(st, posisiTeamLeader, qs, 0).AjukanSubmit(ctx, pelakuTL, "UJI-NB-1", models.IsiSubmit{EmailType: e, Comment: "UJI"})
		if err != nil || h.Position != "ReasFacInMarketing" || h.NBStatus != "NB IS IN Uji - Pemasar'S INBOX" || st.tulis[0].FlagBaru != "2" ||
			st.tulis[0].StatusBinding != nil || st.tulis[0].Riwayat.Posisi != "Policy" {
			t.Errorf("%s: %v %+v %+v", e, err, h, st.tulis)
		}
	}
	st = submitUji()
	decline := models.IsiSubmit{EmailType: "7", SebabNonTc: "Accumulation / Capacity"}
	h, err = layananPolicy(st, posisiTeamLeader, qs, 0).AjukanSubmit(ctx, pelakuTL, "UJI-NB-1", decline)
	if err != nil || h.Position != posisiDecline || h.NBStatus != "NB WAS DECLINED BY  UJI PENGGUNA" || st.tulis[0].StatusWorkBaru != statusDecline ||
		st.tulis[0].FlagBaru != "" {
		t.Errorf("decline: %v %+v", err, h)
	}
	if _, err := layananPolicy(submitUji(), posisiTeamLeader, qs, 0).AjukanSubmit(ctx, pelakuTL, "UJI-NB-1",
		models.IsiSubmit{EmailType: "7"}); !errors.Is(err, ErrMasukanSubmit) {
		t.Errorf("decline tanpa sebab: %v", err)
	}
	// Proteksi Accept TL (MOID kosong) -> 422 sebelum jenjang.
	st = submitUji()
	svc := layananPolicy(st, posisiTeamLeader, qs, 1)
	svc.kasus.(kasusTiruan).ada["UJI-NB-1"].General.MarketingID = ""
	if _, err := svc.AjukanSubmit(ctx, pelakuTL, "UJI-NB-1", accept); !errors.Is(err, ErrProteksiSubmit) || len(st.tulis) != 0 {
		t.Errorf("proteksi: %v", err)
	}
	if _, err := layananPolicy(submitUji(), posisiTeamLeader, qs, 1).AjukanSubmit(ctx, pelakuMkt, "UJI-NB-1", accept); !errors.Is(err, inti.ErrTanpaWewenang) {
		t.Errorf("marketing: %v", err)
	}
	if _, err := layananPolicy(submitUji(), posisiTeamLeader, qs, 1).AjukanSubmit(ctx, pelakuTL, "UJI-NB-1",
		models.IsiSubmit{EmailType: "4", Comment: "UJI"}); !errors.Is(err, ErrMasukanSubmit) {
		t.Errorf("banding: %v", err)
	}
}

// TestAjukanSubmitUWPolicy - UW tahap Policy: Underwriting / JUW A Accept -> jenjang (POLICY/ACCEPT tanpa OUTSTANDING, Validity,
// riwayat Policy / 2); Senior / UW Financial Accept -> selesai Policy (tiket 69a-2), DepHead FacRetro -> 409; Decline tanpa Cause of decline
// (POLICY/DECLINE); Reject / Ask 400; Comment wajib.
func TestAjukanSubmitUWPolicy(t *testing.T) {
	ctx := context.Background()
	qs := spreadingTanpaFacRetro("1000") // ORS: lolos proteksi spreading UW, bukan FacRetro
	accept := models.IsiSubmit{EmailType: "1", Comment: "UJI setuju"}
	for _, c := range []struct {
		posisi, mauPosisi, mauNB string
	}{
		{posisiUWUtama, posisiSeniorUW, "NB IS IN SENIOR UNDERWRITER'S INBOX"},
		{posisiJuniorUWA, posisiUWUtama, "NB IS IN UNDERWRITER 'S INBOX"},
	} {
		st := submitUji()
		h, err := layananPolicy(st, c.posisi, qs, 0).AjukanSubmit(ctx, pelakuPosisi(c.posisi), "UJI-NB-1", accept)
		if err != nil || h.Position != c.mauPosisi || h.NBStatus != c.mauNB || len(st.tulis) != 1 {
			t.Fatalf("%s: %v %+v", c.posisi, err, h)
		}
		w := st.tulis[0]
		if w.Performance == nil || w.Performance.PositionDoc != "POLICY" || w.Performance.StatusDoc != "ACCEPT" || w.PerformanceLanjut != nil ||
			w.Validitas == nil || w.Riwayat.Posisi != "Policy" || w.Riwayat.Putaran != "2" || w.Pega.Workbasket != c.posisi || w.FlagBaru != "" ||
			w.Deductible != nil {
			t.Errorf("%s tulisan: %+v / %+v", c.posisi, w, w.Performance)
		}
	}
	for _, posisi := range []string{posisiSeniorUW, posisiUWFinancial} {
		st := submitUji()
		h, err := layananPolicy(st, posisi, qs, 0).AjukanSubmit(ctx, pelakuPosisi(posisi), "UJI-NB-1", accept)
		if err != nil || h.Position != posisi || len(st.tulis) != 1 || st.tulis[0].Policy == nil || st.tulis[0].StatusWorkBaru != statusSelesaiPolicy ||
			st.tulis[0].Performance == nil || st.tulis[0].Performance.StatusDoc != "ACCEPT" || st.tulis[0].Validitas == nil {
			t.Errorf("%s selesai: %v %+v", posisi, err, h)
		}
	}
	st := submitUji()
	if _, err := layananPolicy(st, posisiDepHeadUW, spreadingSah(), 0).AjukanSubmit(ctx, pelakuPosisi(posisiDepHeadUW), "UJI-NB-1",
		accept); !errors.Is(err, ErrFacRetroBelumTersedia) || len(st.tulis) != 0 {
		t.Errorf("DepHead FacRetro: %v", err)
	}
	st = submitUji()
	h, err := layananPolicy(st, posisiSeniorUW, qs, 0).AjukanSubmit(ctx, pelakuPosisi(posisiSeniorUW), "UJI-NB-1",
		models.IsiSubmit{EmailType: "7", Comment: "UJI tolak"})
	if err != nil || h.Position != posisiDecline || st.tulis[0].Performance == nil || st.tulis[0].Performance.StatusDoc != "DECLINE" ||
		st.tulis[0].Performance.PositionDoc != "POLICY" || st.tulis[0].StatusWorkBaru != statusDecline {
		t.Errorf("decline: %v %+v %+v", err, h, st.tulis)
	}
	for nama, isi := range map[string]models.IsiSubmit{"reject": {EmailType: "2", Comment: "UJI"}, "ask": {EmailType: "3", Comment: "UJI"},
		"banding": {EmailType: "4", Comment: "UJI"}, "tanpa comment": {EmailType: "7"}} {
		if _, err := layananPolicy(submitUji(), posisiSeniorUW, qs, 0).AjukanSubmit(ctx, pelakuPosisi(posisiSeniorUW), "UJI-NB-1", isi); !errors.Is(err, ErrMasukanSubmit) {
			t.Errorf("%s: %v", nama, err)
		}
	}
	// Validity Of Acceptance tetap terhitung di tahap Policy.
	if v, err := layananPolicy(submitUji(), posisiSeniorUW, qs, 0).HitungValiditasSubmit(ctx, pelakuPosisi(posisiSeniorUW), "UJI-NB-1", "5"); err != nil ||
		v.DaysValidity != "5" {
		t.Errorf("validitas: %v %+v", err, v)
	}
}

// TestDeductibleTLPolicy - Protection_Act 9-12 (PositionNote ReasFacInTeamLeader) berlaku juga di tahap Policy: IsBedaCurr tanpa
// Deductible Acceptance -> 422; TampilanSubmit menghitung DeductibleBedaMataUang.
func TestDeductibleTLPolicy(t *testing.T) {
	ctx := context.Background()
	st := submitUji()
	svc := layananSubmitTL(st, true)
	k := st.konteks["UJI-NB-1"]
	k.FlagOnGoingPolicy = "1"
	st.konteks["UJI-NB-1"] = k
	if v, err := svc.TampilanSubmit(ctx, pelakuTL, "UJI-NB-1"); err != nil || v.Tahap != TahapTeamLeaderPolicy || !v.DeductibleBedaMataUang {
		t.Errorf("tampilan: %v %+v", err, v)
	}
	h, err := svc.PeriksaSubmit(ctx, pelakuTL, "UJI-NB-1", IsiPeriksa{EmailType: "1"})
	if err != nil || !reflect.DeepEqual(h.Pesan, []string{"Deductible Acceptance cannot be empty"}) {
		t.Errorf("periksa: %v %+v", err, h)
	}
}

// TestK698TahapPolicyBacaSaja - K69-8 + RALAT (work owner 09-10-2026): flag 1 seluruh tab baca-saja ("Case baca-saja: tahap Policy
// (sudah di-bind)."), kecuali form Submit (TL / UW pemegang workbasket) dan tab Spreading tahap uw (IsUWSpread, seperti UW
// Offer); Accumulation 409; Following false; kunciPeriode tetap true. Instrumen: posisi sama di flag 0 terbuka.
func TestK698TahapPolicyBacaSaja(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct {
		posisi                 string
		pelaku                 inti.Pelaku
		spreading, barisSpread bool
	}{
		{posisiTeamLeader, pelakuTL, false, false},
		{posisiUWUtama, pelakuPosisi(posisiUWUtama), true, true},
	} {
		tk := models.TahapKasus{Position: c.posisi, FlagOnGoingPolicy: "1", ConfirmBinding: "1"}
		svc := Baru(nil).DenganTahap(tahapTiruan{"K": tk})
		boleh, alasan := svc.AksesUbahKasus(ctx, c.pelaku, "K")
		sp, baris := svc.AksesSpreading(ctx, c.pelaku, "K")
		if boleh || alasan != "Case baca-saja: tahap Policy (sudah di-bind)." || sp != c.spreading || baris != c.barisSpread ||
			svc.AksesAkumulasi(ctx, c.pelaku, "K") || svc.AksesFollowing(ctx, c.pelaku, "K") || !svc.KunciPeriode(ctx, "K") {
			t.Errorf("%s: boleh %v %q spreading %v/%v", c.posisi, boleh, alasan, sp, baris)
		}
		if err := svc.wajibBolehUbah(ctx, c.pelaku, "K"); !errors.Is(err, ErrTahapPolicyBacaSaja) || !errors.Is(err, ErrKasusBacaSaja) {
			t.Errorf("%s tulis tab: %v", c.posisi, err)
		}
		if _, err := svc.SimpanAkumulasi(ctx, c.pelaku, "K", nil); !errors.Is(err, ErrTahapPolicyBacaSaja) {
			t.Errorf("%s akumulasi: %v", c.posisi, err)
		}
		if _, err := alasanTidakBoleh(models.KonteksSubmit{Position: c.posisi, FlagOnGoingPolicy: "1"}, c.pelaku); err != nil {
			t.Errorf("%s Submit harus boleh: %v", c.posisi, err)
		}
		// Instrumen: posisi yang sama di flag 0 (tahap Offer) dapat diubah / Submit.
		nol := Baru(nil).DenganTahap(tahapTiruan{"K": {Position: c.posisi, FlagOnGoingPolicy: "0"}})
		if b, _ := nol.AksesUbahKasus(ctx, c.pelaku, "K"); b != (c.posisi == posisiTeamLeader) {
			t.Errorf("%s instrumen flag 0: %v", c.posisi, b)
		}
		if c.spreading && !nol.AksesAkumulasi(ctx, c.pelaku, "K") {
			t.Errorf("%s instrumen akumulasi flag 0", c.posisi)
		}
	}
	// Form Submit TL Policy tetap tampil dan dapat diisi.
	st := submitUji()
	v, err := layananPolicy(st, posisiTeamLeader, spreadingPolicy("1000", "10001", "QS UJI"), 0).TampilanSubmit(ctx, pelakuTL, "UJI-NB-1")
	if err != nil || v.Tahap != TahapTeamLeaderPolicy || !v.BolehSubmit || v.Alasan != "" {
		t.Errorf("TampilanSubmit TL Policy: %v %+v", err, v)
	}
}
