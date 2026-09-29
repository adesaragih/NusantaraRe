package services_test

// Baris adjustment yang lahir saat Submit Register - GILIRAN-14 butir bp.
// TANPA Oracle.
//
// Dibaca sesudah: adjustment.go (BarisPendaftaran, LahirkanBarisPendaftaran).

import (
	"testing"

	"nusantarare/internal/models"
	"nusantarare/internal/services"
	"nusantarare/pkg/utils"
)

// pesertaTerpilih adalah peserta contoh bernilai LITERAL, sebagian lebih
// dari empat angka desimal supaya pembulatan `@divide(…,1,4)` terlihat.
func pesertaTerpilih(t *testing.T) models.Peserta {
	t.Helper()
	return models.Peserta{
		ID: "P-1", IsCheck: models.PenandaDipilih, MataUang: "IDR",
		CedingRetention:  uang(t, "100000.123456", "IDR"),
		ShareNusantaraRe: uang(t, "800000", "IDR"),
		SumInsured:       uang(t, "1000000.12345", "IDR"),
		SumReasured:      uang(t, "900000", "IDR"),
		ShareRetro:       uang(t, "200000.00005", "IDR"),
		RetrocededShare:  uang(t, "150000.99999", "IDR"),
		JumlahKlaim:      uang(t, "25000.123449", "IDR"),
		// Bukan bagian 7.8 - tidak boleh ikut ke baris.
		GrossPremium: uang(t, "50000", "IDR"),
	}
}

// TestBarisPendaftaranMenyalinDelapanMedan7_8 - `SavePesertaClaim` 7.8.
//
// `[terverifikasi]` b3696-b3869: tujuh medan uang
// `@divide(@toDecimal(@replaceAll(.X,",",".")),1,4)` ditambah `.CURRENCY`.
// `SHARE_NUSANTARA_RE` b3743 memakai `@if` yang SAMA dengan peserta b2845 -
// nilai peserta sudah memilihnya (`repository.ShareNusantaraReTeks`).
func TestBarisPendaftaranMenyalinDelapanMedan7_8(t *testing.T) {
	baru, lahir, err := services.BarisPendaftaran(pesertaTerpilih(t))
	if err != nil {
		t.Fatal(err)
	}
	if !lahir {
		t.Fatal("peserta ber-IsCheck \"true\" tidak melahirkan baris")
	}
	for _, k := range []struct {
		medan string
		got   models.Money
		mau   string
	}{
		{"CEDING_RETENTION", baru.CedingRetention, "100000.1235"},
		{"SHARE_NUSANTARA_RE", baru.ShareNusantaraRe, "800000.0000"},
		{"SUM_INSURED", baru.SumInsured, "1000000.1235"},
		{"SUM_REASURED", baru.SumReasured, "900000.0000"},
		{"SHARE_RETRO", baru.ShareRetro, "200000.0001"},
		{"RETROCEDED_SHARE", baru.RetrocededShare, "150001.0000"},
		{"CLAIM_AMOUNT", baru.JumlahKlaim, "25000.1234"},
	} {
		// Dibandingkan sebagai BILANGAN: `bulat` merapatkan nol di ekor.
		if k.got.Amount == nil || k.got.Amount.Cmp(uang(t, k.mau, "IDR").Amount) != 0 {
			t.Errorf("%s = %s, mau %s (empat angka, setengah ke atas)",
				k.medan, utils.FormatDecimal(k.got.Amount), k.mau)
		}
	}
	if baru.JumlahKlaim.Currency != "IDR" {
		t.Errorf("CURRENCY = %q, mau IDR (b3868 `.CURRENCY`)", baru.JumlahKlaim.Currency)
	}
	// Status TIDAK ditulis 7.8 - `Save to RNM` 22.1.3.2 yang menulis "0".
	if baru.KodeStatus != "" || baru.ID != "" || baru.NomorAkseptasi != "" ||
		baru.KomiteID != "" || baru.CurrencyID != "" || baru.NamaBank != "" {
		t.Errorf("baris membawa medan di luar 7.8: %+v", baru)
	}
}

// TestBarisPendaftaranHanyaBagiPesertaTerpilih - WHEN b3919.
//
// ⛔ `.IsCheck=="true"` - TEKS, persis. "false" dan kosong tidak melahirkan
// baris (WhenFalse=3, langkahnya dilewati).
func TestBarisPendaftaranHanyaBagiPesertaTerpilih(t *testing.T) {
	for _, isCheck := range []string{"false", ""} {
		p := pesertaTerpilih(t)
		p.IsCheck = isCheck
		_, lahir, err := services.BarisPendaftaran(p)
		if err != nil {
			t.Fatal(err)
		}
		if lahir {
			t.Errorf("IsCheck %q melahirkan baris", isCheck)
		}
	}
}

// TestBarisPendaftaranKosongDibacaNol - keputusan work owner 29-09-2026.
//
// ⛔ PENYIMPANGAN BERTANGGAL terhadap ADR-U-0027, HANYA di langkah 7.8: ketujuh
// medan `@divide(@toDecimal(…),1,4)` membaca sumber KOSONG sebagai NOL - seperti
// `@toDecimal("")` Pega. `CURRENCY` b3868 bukan `@toDecimal` dan tetap teks
// apa adanya.
func TestBarisPendaftaranKosongDibacaNol(t *testing.T) {
	p := models.Peserta{ID: "P-1", IsCheck: models.PenandaDipilih, MataUang: ""}
	baru, lahir, err := services.BarisPendaftaran(p)
	if err != nil || !lahir {
		t.Fatalf("lahir %v, galat %v", lahir, err)
	}
	for _, k := range []struct {
		medan string
		got   models.Money
	}{
		{"CEDING_RETENTION", baru.CedingRetention}, {"SHARE_NUSANTARA_RE", baru.ShareNusantaraRe},
		{"SUM_INSURED", baru.SumInsured}, {"SUM_REASURED", baru.SumReasured},
		{"SHARE_RETRO", baru.ShareRetro}, {"CLAIM_AMOUNT", baru.JumlahKlaim},
		{"RETROCEDED_SHARE", baru.RetrocededShare},
	} {
		if k.got.Kosong() || k.got.Amount.Sign() != 0 {
			t.Errorf("%s dari sumber kosong = %s, mau 0 (@toDecimal(\"\"))",
				k.medan, utils.FormatDecimal(k.got.Amount))
		}
	}
	if baru.JumlahKlaim.Currency != "" {
		t.Errorf("CURRENCY kosong menjadi %q; ia bukan @toDecimal", baru.JumlahKlaim.Currency)
	}
}

// TestPesertaDibulatkanSeperti7_7 - OQ-N10 ditutup (keputusan 29-09-2026).
//
// `[terverifikasi]` `SavePesertaClaim.xml` 7.7 b2744: sepuluh medan peserta
// `@divide(@toDecimal(@replaceAll(.X,",",".")),1,4)` - GROSS_PREMIUM b2770,
// NET_PREMIUM b2824, SHARE_NUSANTARA_RE b2845, SUM_INSURED b2872,
// CEDING_RETENTION b2893, SUM_REASURED b3107, EM_PERCENT b3134, CLAIM_AMOUNT
// b3281, SHARE_RETRO b3401, RETROCEDED_SHARE b3561.
//
// ⛔ Kosong TETAP kosong di sini: penyimpangan "kosong = nol" berlaku HANYA di
// 7.8 (ADR-U-0027 tetap di tempat lain).
func TestPesertaDibulatkanSeperti7_7(t *testing.T) {
	p := pesertaTerpilih(t)
	p.NetPremium = uang(t, "45000.00005", "IDR")
	em, err := models.NewRatio("0.123456", 6)
	if err != nil {
		t.Fatal(err)
	}
	p.EMPercent = em
	p.SumReasured = models.Money{Currency: "IDR"} // kosong - tetap kosong
	daftar := []models.Peserta{p}
	if err := services.BulatkanPesertaPendaftaran(daftar); err != nil {
		t.Fatal(err)
	}
	q := daftar[0]
	for _, k := range []struct {
		medan string
		got   models.Money
		mau   string
	}{
		{"CEDING_RETENTION", q.CedingRetention, "100000.1235"},
		{"SUM_INSURED", q.SumInsured, "1000000.1235"},
		{"SHARE_RETRO", q.ShareRetro, "200000.0001"},
		{"RETROCEDED_SHARE", q.RetrocededShare, "150001"},
		{"CLAIM_AMOUNT", q.JumlahKlaim, "25000.1234"},
		{"NET_PREMIUM", q.NetPremium, "45000.0001"},
		{"GROSS_PREMIUM", q.GrossPremium, "50000"},
		{"SHARE_NUSANTARA_RE", q.ShareNusantaraRe, "800000"},
	} {
		if k.got.Amount == nil || k.got.Amount.Cmp(uang(t, k.mau, "IDR").Amount) != 0 {
			t.Errorf("%s = %s, mau %s", k.medan, utils.FormatDecimal(k.got.Amount), k.mau)
		}
	}
	if q.EMPercent.Value == nil || q.EMPercent.Value.Cmp(uang(t, "0.1235", "").Amount) != 0 {
		t.Errorf("EM_PERCENT = %s, mau 0.1235", utils.FormatDecimal(q.EMPercent.Value))
	}
	if !q.SumReasured.Kosong() {
		t.Errorf("SUM_REASURED kosong menjadi %s; di 7.7 kosong tetap kosong",
			utils.FormatDecimal(q.SumReasured.Amount))
	}
}

// TestLahirkanBarisPendaftaranSatuPerPeserta - satu baris per peserta
// terpilih, dan baris itu `.AdjustmentList(1)` yang diwarisi putaran.
func TestLahirkanBarisPendaftaranSatuPerPeserta(t *testing.T) {
	terpilih, tidak := pesertaTerpilih(t), pesertaTerpilih(t)
	tidak.ID, tidak.IsCheck = "P-2", "false"
	daftar := []models.Peserta{terpilih, tidak}
	if err := services.LahirkanBarisPendaftaran(daftar); err != nil {
		t.Fatal(err)
	}
	if len(daftar[0].Baris) != 1 {
		t.Errorf("peserta terpilih berbaris %d, mau 1", len(daftar[0].Baris))
	}
	if len(daftar[1].Baris) != 0 {
		t.Errorf("peserta tidak terpilih berbaris %d, mau 0", len(daftar[1].Baris))
	}
	// Putaran berikutnya mewarisi dari baris ini (SetIndexAdjustmentList
	// langkah 3) - bukti bahwa ia benar-benar baris PERTAMA.
	lanjut := daftar[0]
	lanjut.Baris[0].KodeStatus = models.KodeDitolak
	baru, err := services.BarisLanjutan(lanjut)
	if err != nil {
		t.Fatal(err)
	}
	if baru.SumInsured.Amount.Cmp(uang(t, "1000000.1235", "IDR").Amount) != 0 {
		t.Errorf("putaran mewarisi SUM_INSURED %s, mau 1000000.1235",
			utils.FormatDecimal(baru.SumInsured.Amount))
	}
}
