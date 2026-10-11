package services

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
)

// posisiUWUji - PositionNote When IsUWSpread (A..F; D = ReasFacInJuniorUnderwritingB, ADR-0032).
var posisiUWUji = []string{"ReasFacInUnderwriting", "ReasFacInSeniorUnderwriting", "ReasFacInUnderwritingFinancial",
	"ReasFacInJuniorUnderwritingB", "ReasFacInDepHeadUnderwriting", "ReasFacInJuniorUnderwritingA"}

// TestTahapUWDariIsUWSpread - tahap uw = registry ISUWSPREAD: enam posisi, FlagOnGoingPolicy apa pun; nama lama
// ReasFacInJuniorUnderwriting, posisi marketing / TL / Decline, beda huruf / spasi -> bukan uw.
func TestTahapUWDariIsUWSpread(t *testing.T) {
	for _, p := range posisiUWUji {
		for _, f := range []string{"0", "1", "2", ""} {
			if !tahapUW(models.TahapKasus{Position: p, FlagOnGoingPolicy: f}) {
				t.Errorf("%s flag %q bukan uw", p, f)
			}
		}
	}
	for _, p := range []string{"ReasFacInJuniorUnderwriting", "Offer", "ReasFacInMarketing", "ReasFacInTeamLeader", "Decline", "",
		"reasfacinunderwriting", "ReasFacInUnderwriting "} {
		if tahapUW(models.TahapKasus{Position: p, FlagOnGoingPolicy: "0"}) {
			t.Errorf("%q terbaca uw", p)
		}
	}
}

// kasusMatriks - satu baris matriks tahap: id kasus, posisi, workbasket pemegang, tahap yang diharapkan.
type kasusMatriks struct{ id, posisi, pemegang, tahap string }

func matriksTahap() ([]kasusMatriks, tahapTiruan) {
	m := []kasusMatriks{{"M-MKT", "Offer", "ReasFacInMarketing", TahapMarketing},
		{"M-TL", "ReasFacInTeamLeader", "ReasFacInTeamLeader", TahapTeamLeader},
		{"M-DEC", "Decline", "ReasFacInMarketing", ""},
		{"M-JUW-LAMA", "ReasFacInJuniorUnderwriting", "ReasFacInJuniorUnderwriting", ""}}
	for _, p := range posisiUWUji {
		m = append(m, kasusMatriks{"M-" + p, p, p, TahapUW})
	}
	tt := tahapTiruan{}
	for _, k := range m {
		tt[k.id] = models.TahapKasus{Position: k.posisi, FlagOnGoingPolicy: "0"}
	}
	return m, tt
}

func badanSatuCoverage(baris ...models.BarisSpreading) []models.LokasiSpreading {
	return []models.LokasiSpreading{{Items: []models.ItemSpreading{{Coverages: []models.CoverageSpreading{{OldID: "C1", Spreading: baris}}}}}}
}

// galatSama - nil diharapkan nil; selain itu errors.Is.
func galatSama(err, mau error) bool {
	if mau == nil {
		return err == nil
	}
	return errors.Is(err, mau)
}

// TestMatriksPenjagaSpreading - tiket 59: tahap × workbasket untuk PUT spreading / salin / hitung-baris / hitung-share dan
// akses GET kasus. Pemegang = workbasket posisi; "lain" = akun tanpa workbasket itu.
func TestMatriksPenjagaSpreading(t *testing.T) {
	ctx := context.Background()
	matriks, tahap := matriksTahap()
	for _, k := range matriks {
		for _, pegang := range []bool{true, false} {
			peran := "ReasFacInLainUji"
			if pegang {
				peran = k.pemegang
			}
			p := inti.Pelaku{AkunID: "UJI-M", Peran: []string{peran}}
			src := &spreadingTiruan{baca: []models.KasusSpreading{kasusSatuCoverage("50")}, kap: &kapasitasUji,
				idUang: map[string]string{"IDR": idMataUangIDR}}
			svc := layananSpreading(src).DenganTahap(tahap)
			var mauPut, mauSalin, mauBaris error
			mauSalin, mauBaris = ErrSpreadingHanyaUW, ErrSpreadingHanyaUW
			switch {
			case k.tahap == "":
				mauPut = ErrKasusBacaSaja
			case !pegang:
				mauPut = inti.ErrTanpaWewenang
			}
			if k.tahap == TahapUW {
				mauSalin, mauBaris = mauPut, mauPut
			}
			fac := models.BarisSpreading{TreatyType: "10015", TreatyName: "FACOUT", SharePercentage: d("40")}
			_, err := svc.SimpanSpreading(ctx, p, k.id, d("50"), badanSatuCoverage(fac))
			if !galatSama(err, mauPut) {
				t.Errorf("PUT %s pegang=%v: %v, mau %v", k.posisi, pegang, err, mauPut)
			}
			_, err = svc.SalinSpreading(ctx, p, k.id, d("50"), []models.TemplateSpreading{{TreatyType: "10015", TreatyName: "FACOUT",
				SharePercentage: d("100")}})
			if !galatSama(err, mauSalin) {
				t.Errorf("salin %s pegang=%v: %v, mau %v", k.posisi, pegang, err, mauSalin)
			}
			_, err = svc.HitungBarisSpreading(ctx, p, k.id, d("50"), badanSatuCoverage(fac))
			if !galatSama(err, mauBaris) {
				t.Errorf("hitung-baris %s pegang=%v: %v, mau %v", k.posisi, pegang, err, mauBaris)
			}
			if _, err := svc.HitungShareSpreading(ctx, k.id, d("50")); err != nil {
				t.Errorf("hitung-share %s: %v", k.posisi, err)
			}
			boleh, baris := svc.AksesSpreading(ctx, p, k.id)
			if boleh != (mauPut == nil) || baris != (mauPut == nil && k.tahap == TahapUW) {
				t.Errorf("akses %s pegang=%v: boleh %v baris %v", k.posisi, pegang, boleh, baris)
			}
			if ubah, _ := svc.AksesUbahKasus(ctx, p, k.id); k.tahap == TahapUW && ubah {
				t.Errorf("bolehUbah %s di tahap uw", k.posisi)
			}
		}
	}
	// 401 lebih dulu; penjaga tidak aktif (tanpa pembaca tahap) = perilaku tiket 48.
	src := &spreadingTiruan{baca: []models.KasusSpreading{kasusSatuCoverage("50")}, idUang: map[string]string{"IDR": idMataUangIDR}}
	svc := layananSpreading(src).DenganTahap(tahap)
	if _, err := svc.HitungBarisSpreading(ctx, inti.Pelaku{}, "M-MKT", d("50"), nil); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("hitung-baris 401: %v", err)
	}
	if _, err := svc.SalinSpreading(ctx, inti.Pelaku{}, "M-MKT", d("50"), nil); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("salin 401: %v", err)
	}
	if boleh, baris := layananSpreading(src).AksesSpreading(ctx, pelakuUji, "NB-1"); !boleh || !baris {
		t.Errorf("tanpa pembaca tahap: %v %v", boleh, baris)
	}
	if boleh, baris := svc.AksesSpreading(ctx, inti.Pelaku{AkunID: "x", Peran: []string{"ReasFacInMarketing"}}, "UJI-GALAT"); boleh || baris {
		t.Errorf("galat baca tahap harus baca-saja: %v %v", boleh, baris)
	}
}

// TestTahapUWHanyaSpreading - di tahap uw setiap endpoint tulis SELAIN spreading tetap 409 (pemegang workbasket posisi,
// berjabatan IT untuk tab IT-only).
func TestTahapUWHanyaSpreading(t *testing.T) {
	_, tahap := matriksTahap()
	svc := Baru(nil).DenganTahap(tahap).DenganJabatan(jabatanTiruan{"UJI-IT": "IT"})
	for _, pos := range posisiUWUji {
		p := inti.Pelaku{AkunID: "UJI-IT", Peran: []string{pos}}
		for nama, m := range metodeTulis(svc) {
			if nama == "spreading" {
				continue
			}
			if err := m(p, "M-"+pos); !errors.Is(err, ErrKasusBacaSaja) {
				t.Errorf("%s di %s: %v, mau 409", nama, pos, err)
			}
		}
	}
}

// barisTersimpan - kasus % Share 50 dengan satu baris spreading tersimpan Q1 60 % bernilai basi.
func barisTersimpan() models.KasusSpreading {
	k := kasusSatuCoverage("50")
	k.Lokasi[0].Items[0].Coverages[0].Spreading = []models.BarisSpreading{{TreatyType: "Q1", TreatyName: "QS UJI",
		SharePercentage: d("60"), TSISpreaded: d("1"), PremiumSpreaded: d("1"), TSIGrossSpreaded: d("1")}}
	return k
}

// TestSimpanSpreadingMarketingAbaikanBarisBadan - keputusan work owner 08-10-2026: di marketing / Team Leader baris badan
// yang dimanipulasi tidak mengubah spreading tersimpan (jumlah, jenis, share); nilainya dihitung ulang dari % Share badan;
// tanpa spreading otomatis (meski kapasitas treaty ada). Di uw baris badan dipakai.
func TestSimpanSpreadingMarketingAbaikanBarisBadan(t *testing.T) {
	ctx := context.Background()
	_, tahap := matriksTahap()
	manipulasi := badanSatuCoverage(models.BarisSpreading{TreatyType: "X9", TreatyName: "SISIPAN", SharePercentage: d("30")},
		models.BarisSpreading{TreatyType: "X8", TreatyName: "SISIPAN", SharePercentage: d("30")})
	for _, c := range []struct{ id, peran string }{{"M-MKT", "ReasFacInMarketing"}, {"M-TL", "ReasFacInTeamLeader"}} {
		src := &spreadingTiruan{baca: []models.KasusSpreading{barisTersimpan()}, kap: &kapasitasUji, idUang: map[string]string{"IDR": idMataUangIDR}}
		p := inti.Pelaku{AkunID: "UJI", Peran: []string{c.peran}}
		if _, err := layananSpreading(src).DenganTahap(tahap).SimpanSpreading(ctx, p, c.id, d("40"), manipulasi); err != nil {
			t.Fatal(err)
		}
		b := src.tulis[0].Lokasi[0].Items[0].Coverages[0].Spreading
		// PS 40: TSINusantaraRe 400, PremiNR 3.2 -> TSISpreaded 60 % = 240, Premium 1.92, Gross 60 × 1200 × 40 / 10000 = 288.
		if len(b) != 1 || b[0].TreatyType != "Q1" || t8(b[0].SharePercentage) != "60" || t8(b[0].TSISpreaded) != "240" ||
			t8(b[0].PremiumSpreaded) != "1.92" || t8(b[0].TSIGrossSpreaded) != "288" || t8(src.tulis[0].PercentShare) != "40" {
			t.Errorf("%s ditulis %+v", c.id, b)
		}
	}
	// Case baru di marketing (tanpa baris): tetap tanpa baris.
	src := &spreadingTiruan{baca: []models.KasusSpreading{kasusSatuCoverage("")}, kap: &kapasitasUji, idUang: map[string]string{"IDR": idMataUangIDR}}
	mkt := inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInMarketing"}}
	if _, err := layananSpreading(src).DenganTahap(tahap).SimpanSpreading(ctx, mkt, "M-MKT", d("50"), manipulasi); err != nil ||
		len(src.tulis[0].Lokasi[0].Items[0].Coverages[0].Spreading) != 0 {
		t.Errorf("case baru marketing: %v %+v", err, src.tulis)
	}
	// uw: baris badan dipakai (jenis treaty sah - K59-8).
	src = &spreadingTiruan{baca: []models.KasusSpreading{barisTersimpan()}, idUang: map[string]string{"IDR": idMataUangIDR}}
	uw := inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInUnderwriting"}}
	sah := badanSatuCoverage(models.BarisSpreading{TreatyType: "10015", TreatyName: "FACOUT", SharePercentage: d("30")},
		models.BarisSpreading{TreatyType: "10007", TreatyName: "ORS", SharePercentage: d("30")})
	if _, err := layananSpreading(src).DenganTahap(tahap).SimpanSpreading(ctx, uw, "M-ReasFacInUnderwriting", d("40"), sah); err != nil {
		t.Fatal(err)
	}
	if b := src.tulis[0].Lokasi[0].Items[0].Coverages[0].Spreading; len(b) != 2 || b[0].TreatyType != "10015" || t8(b[1].TSISpreaded) != "120" {
		t.Errorf("uw ditulis %+v", b)
	}
}

// TestHitungShareOtomatisHanyaUW - keputusan work owner 08-10-2026: spreading otomatis (GetKapasitasTreaty) hanya di uw;
// marketing / TL tanpa baris tetap tanpa baris, baris tersimpan dipertahankan (jumlah tetap) dan dihitung ulang.
func TestHitungShareOtomatisHanyaUW(t *testing.T) {
	ctx := context.Background()
	_, tahap := matriksTahap()
	baru := func(k models.KasusSpreading) *spreadingTiruan {
		return &spreadingTiruan{baca: []models.KasusSpreading{k}, kap: &kapasitasUji, idUang: map[string]string{"IDR": idMataUangIDR}}
	}
	for _, id := range []string{"M-MKT", "M-TL", "M-DEC"} {
		src := baru(kasusSatuCoverage(""))
		got, err := layananSpreading(src).DenganTahap(tahap).HitungShareSpreading(ctx, id, d("50"))
		if c := got.Lokasi[0].Items[0].Coverages[0]; err != nil || len(c.Spreading) != 0 || t8(c.TSINusantaraRe) != "500" || len(src.tulis) != 0 {
			t.Errorf("%s tanpa baris: %v %+v", id, err, c)
		}
		src = baru(barisTersimpan())
		got, err = layananSpreading(src).DenganTahap(tahap).HitungShareSpreading(ctx, id, d("40"))
		if b := got.Lokasi[0].Items[0].Coverages[0].Spreading; err != nil || len(b) != 1 || b[0].TreatyType != "Q1" || t8(b[0].TSISpreaded) != "240" {
			t.Errorf("%s baris tersimpan: %v %+v", id, err, b)
		}
	}
	for _, pos := range posisiUWUji {
		got, err := layananSpreading(baru(barisTersimpan())).DenganTahap(tahap).HitungShareSpreading(ctx, "M-"+pos, d("50"))
		if b := got.Lokasi[0].Items[0].Coverages[0].Spreading; err != nil || len(b) != 2 || b[0].TreatyType != "Q1" || b[1].TreatyType != "S1" {
			t.Errorf("%s otomatis: %v %+v", pos, err, b)
		}
	}
}

// TestHitungBarisSpreading - tiket 59: hitungan Save tahap uw TANPA menulis; "Invalid Treaty Type Spreading!" (cekSpreadingFactIn
// 8-9) untuk jenis di luar GetTreatyName; Σ > 100 -> 400 seperti PUT.
func TestHitungBarisSpreading(t *testing.T) {
	ctx := context.Background()
	_, tahap := matriksTahap()
	uw := inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInSeniorUnderwriting"}}
	src := &spreadingTiruan{baca: []models.KasusSpreading{barisTersimpan()}, idUang: map[string]string{"IDR": idMataUangIDR}}
	svc := layananSpreading(src).DenganTahap(tahap)
	got, err := svc.HitungBarisSpreading(ctx, uw, "M-ReasFacInSeniorUnderwriting", d("50"),
		badanSatuCoverage(models.BarisSpreading{TreatyType: "10015", TreatyName: "FACOUT", SharePercentage: d("40")}))
	b := got.Lokasi[0].Items[0].Coverages[0].Spreading
	if err != nil || len(src.tulis) != 0 || len(b) != 1 || t8(b[0].TSISpreaded) != "200" || t8(b[0].TSIGrossSpreaded) != "240" ||
		t8(b[0].ClaimEstimation) != "160" || len(got.Pesan) != 0 {
		t.Fatalf("%v tulis %d %+v pesan %v", err, len(src.tulis), b, got.Pesan)
	}
	got, err = svc.HitungBarisSpreading(ctx, uw, "M-ReasFacInSeniorUnderwriting", d("50"),
		badanSatuCoverage(models.BarisSpreading{TreatyType: "X9", SharePercentage: d("40")}))
	if err != nil || len(got.Pesan) != 1 || got.Pesan[0] != pesanJenisTreaty || len(src.tulis) != 0 {
		t.Errorf("jenis treaty: %v %v", err, got.Pesan)
	}
	_, err = svc.HitungBarisSpreading(ctx, uw, "M-ReasFacInSeniorUnderwriting", d("50"), badanSatuCoverage(
		models.BarisSpreading{TreatyType: "A", SharePercentage: d("60")}, models.BarisSpreading{TreatyType: "B", SharePercentage: d("50")}))
	if !errors.Is(err, ErrMasukanSpreading) || !strings.Contains(err.Error(), pesanShareLebih100) {
		t.Errorf("Σ > 100: %v", err)
	}
	if _, err := svc.HitungBarisSpreading(ctx, uw, "M-ReasFacInSeniorUnderwriting", d("50"), nil); !errors.Is(err, ErrSpreadingBerubah) {
		t.Errorf("bentuk: %v", err)
	}
}

// urutanTx - transaksi tiruan yang mencatat mulai / commit / rollback, dan penyimpan yang mencatat tulisan ke urutan yang sama.
type urutanTx struct{ urutan []string }

func (u *urutanTx) transaksi(_ context.Context, fn func(*db.Tx) error) error {
	u.urutan = append(u.urutan, "mulai")
	err := fn(nil)
	u.urutan = append(u.urutan, map[bool]string{true: "rollback", false: "commit"}[err != nil])
	return err
}

type spreadingTercatat struct {
	*spreadingTiruan
	u *urutanTx
}

func (s spreadingTercatat) TulisSpreading(ctx context.Context, tx *db.Tx, id string, k models.KasusSpreading, ubah bool) error {
	s.u.urutan = append(s.u.urutan, "spreading")
	return s.spreadingTiruan.TulisSpreading(ctx, tx, id, k, ubah)
}

type deductionTercatat struct {
	*deductionTiruan
	u *urutanTx
}

func (d deductionTercatat) TulisDeduction(ctx context.Context, tx *db.Tx, id string, k models.KasusDeduction) error {
	d.u.urutan = append(d.u.urutan, "deduction")
	return d.deductionTiruan.TulisDeduction(ctx, tx, id, k)
}

type paymentTercatat struct {
	*paymentTiruan
	u     *urutanTx
	gagal error
}

func (p paymentTercatat) TulisPayment(ctx context.Context, tx *db.Tx, id string, k models.SimpanPayment) error {
	p.u.urutan = append(p.u.urutan, "payment")
	if p.gagal != nil {
		return p.gagal
	}
	return p.paymentTiruan.TulisPayment(ctx, tx, id, k)
}

// badanDuaItem - badan PUT bentuk kasusDuaItem (C1, C2 | C3), satu baris FACOUT 100 % per coverage.
func badanDuaItem() []models.LokasiSpreading {
	b := func(id string) models.CoverageSpreading {
		return models.CoverageSpreading{OldID: id, Spreading: []models.BarisSpreading{{TreatyType: "10015", TreatyName: "FACOUT", SharePercentage: d("100")}}}
	}
	return []models.LokasiSpreading{{Items: []models.ItemSpreading{{Coverages: []models.CoverageSpreading{b("C1"), b("C2")}},
		{Coverages: []models.CoverageSpreading{b("C3")}}}}}
}

// layananLangkah11 - kasusDuaItem PS 50 (premi NR IDR 6, USD 1); Premium Deduction tersimpan 10 % / 5 %, PPnCheck "false",
// Deduction2 IDR 0.6 (nilai lama atas premi 6); Payment 2 cicilan 50 / 50.
func layananLangkah11(tahap tahapTiruan, gagal error) (*Service, *urutanTx, *deductionTiruan, *paymentTiruan) {
	u := &urutanTx{}
	dt := &deductionTiruan{tarif: &tarifUji, ada: models.KasusDeduction{RICommision: d("10"), PctBrokerageFee: d("5"), PPnCheck: "false",
		MataUang: []models.BarisDeduction{{Currency: "IDR", Premium: d("6"), Commision: d("0.6"), Deduction2: d("0.6"), BrokerageFee: d("0.3")},
			{Currency: "USD", Premium: d("1"), Commision: d("0.1"), BrokerageFee: d("0.05")}}}}
	cicilan := []models.CicilanPayment{{InstallmentNo: 1, Persen: d("50")}, {InstallmentNo: 2, Persen: d("50")}}
	pt := &paymentTiruan{ada: models.KasusPayment{Installment: "2", Cicilan: map[string][]models.CicilanPayment{"IDR": cicilan, "USD": cicilan}}}
	src := &spreadingTiruan{baca: []models.KasusSpreading{kasusDuaItem("50", "USD")}, idUang: map[string]string{"IDR": idMataUangIDR}}
	svc := Baru(nil).DenganSpreading(spreadingTercatat{src, u}).DenganDeduction(deductionTercatat{dt, u}).
		DenganPayment(paymentTercatat{pt, u, gagal}).DenganTransaksi(u.transaksi).DenganTahap(tahap)
	return svc, u, dt, pt
}

// TestSimpanSpreadingUWLangkah11 - PUT tahap uw: share 50 -> 25 -> Premium Deduction (SetRIComm_Act dari persen tersimpan,
// PPnCheck "true", Deduction2 tetap) dan cicilan Payment tersimpan (CountPaymentInstallment_Act, Net mengurangkan Deduction2
// K51-1) ditulis ulang di SATU transaksi bersama spreading; gagal tulis Payment -> satu rollback, galat diteruskan.
func TestSimpanSpreadingUWLangkah11(t *testing.T) {
	ctx := context.Background()
	_, tahap := matriksTahap()
	uw := inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInUnderwriting"}}
	svc, u, dt, pt := layananLangkah11(tahap, nil)
	if _, err := svc.SimpanSpreading(ctx, uw, "M-ReasFacInUnderwriting", d("25"), badanDuaItem()); err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(u.urutan, ","); got != "mulai,spreading,deduction,payment,commit" {
		t.Fatalf("urutan %s", got)
	}
	// Premi NR PS 25: IDR 3, USD 0.5. Commision 10 %, Brokerage 5 %, PPh 2 % / PPN 11 % × Brokerage.
	w := dt.tulis[0]
	if got := barisUji(t, w.MataUang); got != "IDR/3/0.3/0.6/0.15/0.003/0.0165 USD/0.5/0.05//0.025/0.0005/0.00275" ||
		w.PPnCheck != "true" || t8(w.RICommision) != "10" || t8(w.PctBrokerageFee) != "5" {
		t.Errorf("deduction %s %+v", got, w)
	}
	// Net IDR = 3 − 0.3 − 0.6 − 0.15 + 0.003 + 0.0165 = 1.9695; USD = 0.5 − 0.05 − 0.025 + 0.0005 + 0.00275 = 0.42825.
	p := pt.tulis[0]
	if p.Installment != 2 || len(p.MataUang) != 2 || t8(p.MataUang[0].NetPremium) != "1.9695" || t8(p.MataUang[1].NetPremium) != "0.42825" ||
		t8(p.MataUang[0].Cicilan[0].Premium) != "1.5" || t8(p.MataUang[0].Cicilan[1].PaymentTotal) != "0.98475" {
		t.Errorf("payment %+v", p)
	}
	// Gagal di tengah: satu transaksi, rollback, galat sampai pemanggil.
	gagal := errors.New("UJI payment gagal")
	svc, u, _, _ = layananLangkah11(tahap, gagal)
	if _, err := svc.SimpanSpreading(ctx, uw, "M-ReasFacInUnderwriting", d("25"), badanDuaItem()); !errors.Is(err, gagal) {
		t.Fatalf("galat %v", err)
	}
	if got := strings.Join(u.urutan, ","); got != "mulai,spreading,deduction,payment,rollback" {
		t.Errorf("urutan gagal %s", got)
	}
}

// TestSimpanSpreadingLangkah11HanyaUW - marketing / Team Leader / penjaga tidak aktif: perilaku tidak berubah (hanya
// spreading ditulis); uw tanpa Premium Deduction tersimpan: tidak ada yang ditulis selain spreading.
func TestSimpanSpreadingLangkah11HanyaUW(t *testing.T) {
	ctx := context.Background()
	_, tahap := matriksTahap()
	for _, c := range []struct {
		tahap tahapTiruan
		id    string
		peran string
	}{{tahap, "M-MKT", "ReasFacInMarketing"}, {tahap, "M-TL", "ReasFacInTeamLeader"}, {nil, "NB-1", "x"}} {
		svc, u, dt, pt := layananLangkah11(c.tahap, nil)
		if c.tahap == nil {
			svc.DenganTahap(nil)
		}
		if _, err := svc.SimpanSpreading(ctx, inti.Pelaku{AkunID: "UJI", Peran: []string{c.peran}}, c.id, d("25"), badanDuaItem()); err != nil {
			t.Fatal(err)
		}
		if got := strings.Join(u.urutan, ","); got != "mulai,spreading,commit" || len(dt.tulis) != 0 || len(pt.tulis) != 0 {
			t.Errorf("%s: %s", c.id, got)
		}
	}
	svc, u, dt, _ := layananLangkah11(tahap, nil)
	dt.ada = models.KasusDeduction{}
	uw := inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInUnderwriting"}}
	if _, err := svc.SimpanSpreading(ctx, uw, "M-ReasFacInUnderwriting", d("25"), badanDuaItem()); err != nil ||
		strings.Join(u.urutan, ",") != "mulai,spreading,commit" {
		t.Errorf("uw tanpa deduction tersimpan: %v %v", err, u.urutan)
	}
}

// TestSimpanSpreadingUWJenisTreatyTidakSah - K59-8: PUT uw dengan jenis treaty di luar GetTreatyName / ORS / FACOUT -> 400
// "Invalid Treaty Type Spreading!" dan TIDAK ada tulisan (spreading maupun langkah 11); jenis sah -> lolos. Marketing tidak
// diperiksa (baris badan diabaikan).
func TestSimpanSpreadingUWJenisTreatyTidakSah(t *testing.T) {
	ctx := context.Background()
	_, tahap := matriksTahap()
	tidakSah := badanDuaItem()
	tidakSah[0].Items[1].Coverages[0].Spreading[0].TreatyType = "X9"
	uw := inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInUnderwriting"}}
	svc, u, dt, pt := layananLangkah11(tahap, nil)
	_, err := svc.SimpanSpreading(ctx, uw, "M-ReasFacInUnderwriting", d("25"), tidakSah)
	if !errors.Is(err, ErrMasukanSpreading) || !strings.HasSuffix(err.Error(), ": "+pesanJenisTreaty) || len(u.urutan) != 0 ||
		len(dt.tulis) != 0 || len(pt.tulis) != 0 {
		t.Fatalf("%v urutan %v", err, u.urutan)
	}
	svc, u, _, _ = layananLangkah11(tahap, nil)
	if _, err := svc.SimpanSpreading(ctx, uw, "M-ReasFacInUnderwriting", d("25"), badanDuaItem()); err != nil || len(u.urutan) == 0 {
		t.Errorf("sah: %v %v", err, u.urutan)
	}
	svc, u, _, _ = layananLangkah11(tahap, nil)
	mkt := inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInMarketing"}}
	if _, err := svc.SimpanSpreading(ctx, mkt, "M-MKT", d("25"), tidakSah); err != nil || strings.Join(u.urutan, ",") != "mulai,spreading,commit" {
		t.Errorf("marketing: %v %v", err, u.urutan)
	}
}

// TestSimpanSpreadingUWDeduction2LebihPremi - K59-9: Deduction 2 tersimpan IDR 0.6. Share 50 -> 5 (premi IDR 0.6, sama) lolos;
// 50 -> 4 (premi IDR 0.48 < 0.6) -> 400 verbatim dan nol tulisan; tanpa Premium Deduction tersimpan -> lolos; marketing tidak
// diperiksa.
func TestSimpanSpreadingUWDeduction2LebihPremi(t *testing.T) {
	ctx := context.Background()
	_, tahap := matriksTahap()
	uw := inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInUnderwriting"}}
	svc, u, dt, pt := layananLangkah11(tahap, nil)
	_, err := svc.SimpanSpreading(ctx, uw, "M-ReasFacInUnderwriting", d("4"), badanDuaItem())
	if !errors.Is(err, ErrMasukanSpreading) || !strings.HasSuffix(err.Error(), ": "+pesanDeduction2Lebih) || len(u.urutan) != 0 ||
		len(dt.tulis) != 0 || len(pt.tulis) != 0 {
		t.Fatalf("%v urutan %v", err, u.urutan)
	}
	// Batas Deduction 2 = premi (0.6) lolos K59-9, tetapi Net IDR = 0.6 − 0.06 − 0.6 − 0.03 + 0.0006 + 0.0033 < 0 -> K59-10.
	svc, u, _, _ = layananLangkah11(tahap, nil)
	if _, err := svc.SimpanSpreading(ctx, uw, "M-ReasFacInUnderwriting", d("5"), badanDuaItem()); err == nil ||
		strings.Contains(err.Error(), pesanDeduction2Lebih) || len(u.urutan) != 0 {
		t.Errorf("sama dengan premi: %v %v", err, u.urutan)
	}
	// Share turun, Deduction 2 < premi (IDR 1.2): lolos, Payment ikut ditulis.
	svc, u, _, _ = layananLangkah11(tahap, nil)
	if _, err := svc.SimpanSpreading(ctx, uw, "M-ReasFacInUnderwriting", d("10"), badanDuaItem()); err != nil ||
		strings.Join(u.urutan, ",") != "mulai,spreading,deduction,payment,commit" {
		t.Errorf("share 10: %v %v", err, u.urutan)
	}
	svc, u, dt, _ = layananLangkah11(tahap, nil)
	dt.ada = models.KasusDeduction{}
	if _, err := svc.SimpanSpreading(ctx, uw, "M-ReasFacInUnderwriting", d("4"), badanDuaItem()); err != nil {
		t.Errorf("tanpa deduction tersimpan: %v", err)
	}
	svc, _, _, _ = layananLangkah11(tahap, nil)
	if _, err := svc.SimpanSpreading(ctx, inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInMarketing"}}, "M-MKT", d("4"), badanDuaItem()); err != nil {
		t.Errorf("marketing: %v", err)
	}
}

// TestSimpanSpreadingUWNetNegatif - K59-10. Deduction tersimpan: 10 % / 5 %, PPH 2 / PPN 11, Deduction2 IDR 0.6. Net IDR =
// premi × (1 − 0.10 − 0.05 + 0.05 × 0.13) − 0.6 = premi × 0.8565 − 0.6: premi 0.6 (share 5) -> < 0 -> 400 pesan K51-7 dan nol
// tulisan. Batas Net = 0: Deduction2 diset = premi × 0.8565 pada share 10 (premi IDR 1.2 -> 1.0278) -> lolos. Tanpa
// deduction tersimpan dan marketing: tidak diperiksa.
func TestSimpanSpreadingUWNetNegatif(t *testing.T) {
	ctx := context.Background()
	_, tahap := matriksTahap()
	uw := inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInUnderwriting"}}
	// Teks K51-7 dipaku (bersama Save Payment), bukan dibangun dari fungsinya sendiri.
	if got := pesanNetNegatif("IDR"); got != "Total Payment IDR negatif - potongan melebihi premi" {
		t.Fatalf("pesan K51-7 berubah: %q", got)
	}
	svc, u, dt, pt := layananLangkah11(tahap, nil)
	_, err := svc.SimpanSpreading(ctx, uw, "M-ReasFacInUnderwriting", d("5"), badanDuaItem())
	if !errors.Is(err, ErrMasukanSpreading) || !strings.HasSuffix(err.Error(), ": "+pesanNetNegatif("IDR")) || len(u.urutan) != 0 ||
		len(dt.tulis) != 0 || len(pt.tulis) != 0 {
		t.Fatalf("%v urutan %v", err, u.urutan)
	}
	svc, u, dt, pt = layananLangkah11(tahap, nil)
	dt.ada.MataUang[0].Deduction2 = d("1.0278")
	if _, err := svc.SimpanSpreading(ctx, uw, "M-ReasFacInUnderwriting", d("10"), badanDuaItem()); err != nil ||
		strings.Join(u.urutan, ",") != "mulai,spreading,deduction,payment,commit" || t8(pt.tulis[0].MataUang[0].NetPremium) != "0" {
		t.Errorf("Net = 0: %v %v %+v", err, u.urutan, pt.tulis)
	}
	svc, _, dt, _ = layananLangkah11(tahap, nil)
	dt.ada = models.KasusDeduction{}
	if _, err := svc.SimpanSpreading(ctx, uw, "M-ReasFacInUnderwriting", d("5"), badanDuaItem()); err != nil {
		t.Errorf("tanpa deduction tersimpan: %v", err)
	}
	svc, _, _, _ = layananLangkah11(tahap, nil)
	if _, err := svc.SimpanSpreading(ctx, inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInMarketing"}}, "M-MKT", d("5"), badanDuaItem()); err != nil {
		t.Errorf("marketing: %v", err)
	}
}

// TestRingkasanTidakMengubahMasukan - bug DEV 08-10-2026 (ORA-01722): sesudah hitungSaveSpreading, premiPerMataUang /
// ringkasanRNM / barisLangkah11 / totalMataUang tidak boleh mengubah NR atau baris spreading kasus pemanggil. Angka DEV:
// TSILiability 1e11, Premium 15924826.9040554, share 50.
func TestRingkasanTidakMengubahMasukan(t *testing.T) {
	k := kasusSatuCoverage("")
	c := &k.Lokasi[0].Items[0].Coverages[0]
	c.TSILiability, c.Premium, c.Discount, c.TSI = d("100000000000"), d("15924826.9040554"), d("0"), d("100000000000")
	c.Spreading = []models.BarisSpreading{{TreatyType: "10015", SharePercentage: d("100")}}
	adaTop := hitungSaveSpreading(&k, d("50"))
	potret := func() string {
		c := k.Lokasi[0].Items[0].Coverages[0]
		b := c.Spreading[0]
		return strings.Join([]string{utils.FormatDecimal(c.TSINusantaraRe), utils.FormatDecimal(c.PremiNusantaraRe),
			utils.FormatDecimal(b.TSISpreaded), utils.FormatDecimal(b.PremiumSpreaded), utils.FormatDecimal(b.TSIGrossSpreaded)}, "|")
	}
	awal := potret()
	if awal != "50000000000|7962413.4520277|50000000000|7962413.4520277|50000000000" || maskSpreading(k) != nil {
		t.Fatalf("hitungSaveSpreading: %s %v", awal, maskSpreading(k))
	}
	premiPerMataUang(k)
	ringkasanRNM(k)
	barisLangkah11(k, &dasarLangkah11{deduction: models.KasusDeduction{RICommision: d("10")}})
	totalMataUang(k, adaTop)
	if got := potret(); got != awal || maskSpreading(k) != nil {
		t.Errorf("masukan berubah: %s, mau %s (%v)", got, awal, maskSpreading(k))
	}
	// salinDalamSpreading: mengubah salinan (baris, coverage, total lokasi) tidak mengubah asal.
	k.Lokasi[0].Total = []models.TotalSpreading{{TreatyType: "10015"}}
	s := salinDalamSpreading(k)
	s.Lokasi[0].Items[0].Coverages[0].Spreading[0].TreatyType = "UBAH"
	s.Lokasi[0].Items[0].Coverages[0].TSINusantaraRe = d("1")
	s.Lokasi[0].Total[0].TreatyType = "UBAH"
	if potret() != awal || k.Lokasi[0].Items[0].Coverages[0].Spreading[0].TreatyType != "10015" || k.Lokasi[0].Total[0].TreatyType != "10015" {
		t.Errorf("salinan dalam berbagi slice")
	}
}

// TestTampilanTidakMengubahMasukan - tampilan (by value) menghitung total / membulatkan atas salinan dalam.
func TestTampilanTidakMengubahMasukan(t *testing.T) {
	src := &spreadingTiruan{baca: []models.KasusSpreading{kasusSatuCoverage("50")}, idUang: map[string]string{"IDR": idMataUangIDR}}
	k := kasusSatuCoverage("50")
	k.Lokasi[0].Items[0].Coverages[0].Spreading = []models.BarisSpreading{{TreatyType: "10015", SharePercentage: d("40")}}
	hitungNR(&k) // skala 20, sengaja belum dibulatkan
	sebelum := utils.FormatDecimal(k.Lokasi[0].Items[0].Coverages[0].TSINusantaraRe)
	if _, err := layananSpreading(src).tampilan(context.Background(), "NB-1", k, nil, false); err != nil {
		t.Fatal(err)
	}
	if got := utils.FormatDecimal(k.Lokasi[0].Items[0].Coverages[0].TSINusantaraRe); got != sebelum || k.Lokasi[0].Total != nil {
		t.Errorf("tampilan mengubah masukan: %s -> %s, total %v", sebelum, got, k.Lokasi[0].Total)
	}
}

// TestPeriksaTemplateSpreading - POST …/spreading/periksa-template: proteksiTemplate yang sama dengan salin (QS baris pertama
// kecuali ORS / FACOUT, kembar), baris tanpa treatyType diabaikan, tidak menulis; penjaga = salin.
func TestPeriksaTemplateSpreading(t *testing.T) {
	ctx := context.Background()
	_, tahap := matriksTahap()
	src := &spreadingTiruan{baca: []models.KasusSpreading{kasusSatuCoverage("50")}, idUang: map[string]string{"IDR": idMataUangIDR}}
	svc := layananSpreading(src).DenganTahap(tahap)
	uw := inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInUnderwriting"}}
	tpl := func(jenis ...string) []models.TemplateSpreading {
		out := []models.TemplateSpreading{}
		for _, j := range jenis {
			out = append(out, models.TemplateSpreading{TreatyType: j})
		}
		return out
	}
	for nama, c := range map[string]struct {
		jenis []string
		mau   []string
	}{
		"QS dulu":           {[]string{"Q1", "S1"}, []string{}},
		"SPL saja":          {[]string{"S1"}, []string{pesanTanpaQS}},
		"ORS pertama":       {[]string{"10007", "S1"}, []string{}},
		"FACOUT pertama":    {[]string{" 10015 "}, []string{}},
		"kembar":            {[]string{"Q1", "Q1"}, []string{pesanJenisKembar}},
		"kembar & tanpa QS": {[]string{"S1", "S1"}, []string{pesanJenisKembar, pesanTanpaQS}},
		"kosong diabaikan":  {[]string{"", "Q1", " "}, []string{}},
		"semua kosong":      {[]string{"", ""}, []string{}},
		"tanpa baris":       {nil, []string{}},
	} {
		got, err := svc.PeriksaTemplateSpreading(ctx, uw, "M-ReasFacInUnderwriting", tpl(c.jenis...))
		if err != nil || !reflect.DeepEqual(got, c.mau) {
			t.Errorf("%s: %v %q", nama, err, got)
		}
	}
	if len(src.tulis) != 0 {
		t.Errorf("menulis %d", len(src.tulis))
	}
	for nama, c := range map[string]struct {
		pl  inti.Pelaku
		id  string
		tpl []models.TemplateSpreading
		mau error
	}{
		"401":       {inti.Pelaku{}, "M-ReasFacInUnderwriting", tpl("Q1"), inti.ErrTanpaIdentitas},
		"409 mkt":   {inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInMarketing"}}, "M-MKT", tpl("Q1"), ErrSpreadingHanyaUW},
		"409 TL":    {inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInTeamLeader"}}, "M-TL", tpl("Q1"), ErrSpreadingHanyaUW},
		"403":       {inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInMarketing"}}, "M-ReasFacInUnderwriting", tpl("Q1"), inti.ErrTanpaWewenang},
		"400 baris": {uw, "M-ReasFacInUnderwriting", make([]models.TemplateSpreading, banyakBarisSpreading+1), ErrMasukanSpreading},
	} {
		if _, err := svc.PeriksaTemplateSpreading(ctx, c.pl, c.id, c.tpl); !errors.Is(err, c.mau) {
			t.Errorf("%s: %v", nama, err)
		}
	}
}
