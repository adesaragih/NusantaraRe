package services

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
)

// hariUji - tanggal acuan Loss Ratio (WIB).
var hariUji = time.Date(2026, 10, 3, 9, 0, 0, 0, WIB)

func catatan(tgl, amount, claim string) models.CatatanKerugian {
	return models.CatatanKerugian{DateOfLoss: tgl, Currency: "IDR", Amount: ds(amount), Claim: ds(claim)}
}

// TestHitungLossRatio - tiket 42, SetLossRatio_Act + W-4: <= 365 hari masuk LR1 DAN LR35 (kumulatif), <= 1825 hanya
// LR35; ΣClaim 0 -> 0; ΣAmount 0 -> 0 (W-4); tanpa tanggal tidak dijumlah (A147); 8 desimal setengah-ke-atas (A148).
func TestHitungLossRatio(t *testing.T) {
	nol := lrUji("0", "0", "0", "0")
	for nama, u := range map[string]struct {
		rugi []models.CatatanKerugian
		mau  models.LossRatio
	}{
		"kosong": {nil, nol},
		"satu tahun dan lima tahun": {[]models.CatatanKerugian{catatan("03-10-2025", "1000", "250"), catatan("04-10-2022", "1000", "750")},
			lrUji("0.25", "25", "0.5", "50")},
		"tepat 365 hari masuk": {[]models.CatatanKerugian{catatan("03-10-2025", "400", "100")},
			lrUji("0.25", "25", "0.25", "25")},
		"366 hari hanya LR35": {[]models.CatatanKerugian{catatan("02-10-2025", "400", "100")},
			lrUji("0", "0", "0.25", "25")},
		"lebih dari 1825 hari": {[]models.CatatanKerugian{catatan("01-01-2020", "400", "100")}, nol},
		"claim nol":            {[]models.CatatanKerugian{catatan("01-10-2026", "400", "0")}, nol},
		"amount nol W-4":       {[]models.CatatanKerugian{catatan("01-10-2026", "", "100")}, nol},
		"tanpa tanggal":        {[]models.CatatanKerugian{catatan("", "400", "100")}, nol},
		"tanggal sesudah hari ini ikut": {[]models.CatatanKerugian{catatan("31-12-2026", "400", "100")},
			lrUji("0.25", "25", "0.25", "25")},
		"pembulatan 8 desimal": {[]models.CatatanKerugian{catatan("01-10-2026", "3", "2")},
			lrUji("0.66666667", "66.66666667", "0.66666667", "66.66666667")},
		"uang besar eksak": {[]models.CatatanKerugian{catatan("01-10-2026", "123456789012345678901234567890.12345678", "123456789012345678901234567890.12345678")},
			lrUji("1", "100", "1", "100")},
	} {
		if got, err := hitungLossRatio(u.rugi, hariUji); err != nil || lrTeks(got) != lrTeks(u.mau) {
			t.Errorf("%s: %v (%v), mau %v", nama, lrTeks(got), err, lrTeks(u.mau))
		}
	}
}

// TestLossRatioTerlalu - ΣClaim / ΣAmount melebihi NUMBER(38,8) -> 400 ber-indeks, bukan galat Oracle (500).
func TestLossRatioTerlalu(t *testing.T) {
	besar := "123456789012345678901234567890"
	if _, err := hitungLossRatio([]models.CatatanKerugian{catatan("01-10-2026", "0.00000001", besar)}, hariUji); !errors.Is(err, errLRTerlalu) {
		t.Errorf("mau errLRTerlalu: %v", err)
	}
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-USER"}
	var panggil int
	svc := Baru(nil).DenganObjek(objekTiruan{ada: map[string][]models.ObjekFire{"UJI-NB-1": {}}}).DenganTransaksi(tanpaTx).
		DenganPilihanItem(pilihanItemTiruan{&panggil}).DenganJam(func() time.Time { return hariUji }).
		DenganKasus(kasusTiruan{ada: map[string]*models.Kasus{"UJI-NB-1": {CaseID: "UJI-NB-1"}}})
	obj := []models.ObjekFire{{ObjectType: "UJI"}, {ObjectType: "UJI", LossRecords: []models.CatatanKerugian{catatan("01-10-2026", "0.00000001", besar)}}}
	if _, err := svc.GantiObjek(ctx, akun, "UJI-NB-1", obj); !errors.Is(err, ErrMasukanObjek) || !strings.Contains(err.Error(), "baris[1].lossRatio") {
		t.Errorf("400: %v", err)
	}
}

// TestPeriksaKerugian - currency wajib, uang desimal <= 8, tanggal DD-MM-YYYY sah, lebar kolom; detail tidak wajib.
func TestPeriksaKerugian(t *testing.T) {
	if m := periksaKerugian(0, []models.CatatanKerugian{catatan("29-02-2024", "1.5", "0"), {Currency: "USD"}}); len(m) != 0 {
		t.Fatalf("sah ditolak: %v", m)
	}
	for nama, u := range map[string]struct {
		c     models.CatatanKerugian
		pesan string
	}{
		"tanpa mata uang": {models.CatatanKerugian{Currency: " "}, "baris[1].lossRecords[0].currency wajib"},
		"claim 31 digit":  {models.CatatanKerugian{Currency: "IDR", Claim: ds(strings.Repeat("9", 31))}, "baris[1].lossRecords[0].claim harus"},
		"prevention minus": {models.CatatanKerugian{Currency: "IDR", PreventionOfLoss: ds("-1")},
			"baris[1].lossRecords[0].preventionOfLoss harus"},
		"amount 9 desimal": {models.CatatanKerugian{Currency: "IDR", Amount: ds("1.123456789")}, "baris[1].lossRecords[0].amount harus"},
		"tanggal 31-02":    {models.CatatanKerugian{Currency: "IDR", DateOfLoss: "31-02-2026"}, "baris[1].lossRecords[0].dateOfLoss bukan tanggal"},
		"tanggal ISO":      {models.CatatanKerugian{Currency: "IDR", DateOfLoss: "2026-01-01"}, "dateOfLoss bukan tanggal"},
		"detail 501":       {models.CatatanKerugian{Currency: "IDR", Detail: strings.Repeat("U", 501)}, "baris[1].lossRecords[0].detail paling banyak 500"},
	} {
		if m := strings.Join(periksaKerugian(1, []models.CatatanKerugian{u.c}), "; "); !strings.Contains(m, u.pesan) {
			t.Errorf("%s: %q, mau %q", nama, m, u.pesan)
		}
	}
}

// TestGantiObjekKerugian - tiket 42: simpan menghitung LR (lossRatio badan diabaikan), CoinsName = tertanggung case,
// dateOfLoss tersimpan teks Pega 12:00 WIB dan dibaca ulang DD-MM-YYYY; mata uang diperiksa ke CURRENCY; 503 tanpa
// pembaca case bila ada catatan.
func TestGantiObjekKerugian(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-USER"}
	var panggil int
	tersimpan := objekTiruan{ada: map[string][]models.ObjekFire{"UJI-NB-1": {}}}
	kasus := kasusTiruan{ada: map[string]*models.Kasus{"UJI-NB-1": {CaseID: "UJI-NB-1", InsuredName: "UJI TERTANGGUNG"}}}
	svc := Baru(nil).DenganObjek(tersimpan).DenganTransaksi(tanpaTx).DenganPilihanItem(pilihanItemTiruan{&panggil}).
		DenganKasus(kasus).DenganJam(func() time.Time { return hariUji })
	c := catatan("01-10-2026", "400", "100")
	c.CoinsName = "UJI DIABAIKAN"
	obj := []models.ObjekFire{{ObjectType: "UJI", LossRecords: []models.CatatanKerugian{c},
		LossRatio: models.LossRatio{OneYearAmount: ds("9"), OneYearPercent: ds("9")}}}
	d, err := svc.GantiObjek(ctx, akun, "UJI-NB-1", obj)
	if err != nil || len(d) != 1 || len(d[0].LossRecords) != 1 {
		t.Fatalf("%+v (%v)", d, err)
	}
	if s := tersimpan.ada["UJI-NB-1"][0]; s.LossRecords[0].DateOfLoss != "20261001T050000.000 GMT" ||
		s.LossRecords[0].CoinsName != "UJI TERTANGGUNG" || utils.FormatDecimal(s.LossRatio.OneYearAmount) != "0.25" || utils.FormatDecimal(s.LossRatio.OneYearPercent) != "25" {
		t.Errorf("tersimpan: %+v", s)
	}
	if got := d[0].LossRecords[0]; got.DateOfLoss != "01-10-2026" || got.CoinsName != "UJI TERTANGGUNG" {
		t.Errorf("dibaca ulang: %+v", got)
	}
	if obj[0].LossRecords[0].CoinsName != "UJI DIABAIKAN" {
		t.Error("masukan pemanggil tidak boleh diubah")
	}
	obj[0].LossRecords[0].Currency = "ITL"
	if _, err := svc.GantiObjek(ctx, akun, "UJI-NB-1", obj); !errors.Is(err, ErrMasukanObjek) ||
		!strings.Contains(err.Error(), `baris[0].lossRecords[0].currency "ITL" tidak ada`) {
		t.Errorf("ITL: %v", err)
	}
	obj[0].LossRecords[0].Currency = "IDR"
	tanpaKasus := Baru(nil).DenganObjek(tersimpan).DenganTransaksi(tanpaTx).DenganPilihanItem(pilihanItemTiruan{&panggil})
	if _, err := tanpaKasus.GantiObjek(ctx, akun, "UJI-NB-1", obj); !errors.Is(err, ErrObjekTanpaDatabase) {
		t.Errorf("tanpa pembaca case: %v", err)
	}
}

// TestTanggalKerugian - teks Pega DateTime / Date -> DD-MM-YYYY (WIB); teks lain apa adanya; pulang-pergi.
func TestTanggalKerugian(t *testing.T) {
	for masuk, mau := range map[string]string{"20261001T050000.000 GMT": "01-10-2026", "20250131T170000.000 GMT": "01-02-2025",
		"20240229": "29-02-2024", "": "", "UJI": "UJI"} {
		if got := tanggalKerugianKeKabel(masuk); got != mau {
			t.Errorf("%q -> %q, mau %q", masuk, got, mau)
		}
	}
	if got := tanggalKerugianKeKabel(tanggalKerugianKePega("15-08-2026")); got != "15-08-2026" {
		t.Errorf("pulang-pergi %q", got)
	}
}

// lrUji / lrTeks - LossRatio uji dari / ke teks (perbandingan nilai lewat teks kanonik apd).
func lrUji(a1, p1, a35, p35 string) models.LossRatio {
	return models.LossRatio{OneYearAmount: ds(a1), OneYearPercent: ds(p1), ThreeFiveYearAmount: ds(a35), ThreeFiveYearPercent: ds(p35)}
}

func lrTeks(lr models.LossRatio) [4]string {
	return [4]string{utils.FormatDecimal(lr.OneYearAmount), utils.FormatDecimal(lr.OneYearPercent),
		utils.FormatDecimal(lr.ThreeFiveYearAmount), utils.FormatDecimal(lr.ThreeFiveYearPercent)}
}
