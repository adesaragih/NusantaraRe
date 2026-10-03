package models

// Uji seam 3 - hari tutup buku (tiket 02). Hari batas DIBERIKAN pemanggil, yang
// membacanya dari `POOLDATA.TANGGAL_CLOSING` (`inti/backend/penomor`); di sini
// batas tiruan 25. Yang diuji: hari = batas TIDAK digeser (pembanding `>`).
//
//	InputPolicyTreatyInPre_Act langkah 9 ("Set Production Date kalau diatas tanggal 25"):
//	  syarat `@substring(pyWorkPage.PolicyTreatyIn.StatementDate,6,2)>25`
//	  ProductionDate = @addCalendar(StatementDate,'0','1','0','0','0','0','0')
//	  ProductionDate = @substring(ProductionDate,0,6) + "01" + @substring(ProductionDate,8)
//	  -> tanggal 1 bulan berikut; bagian waktu (`@substring(..,8)`) dipertahankan
//
//	GeneratePolicyNoTreaty_Act langkah 5.3 (`Local.TglProd` = GETTanggalClosing_SQL):
//	  syarat `@toInt(@FormatDateTime(ProductionDate,"dd","Asia/Jakarta","in_ID"))>Local.TglProd`
//	  ProductionDate = @FormatDateTime(ProductionDate,"yyyyMM",..)+"01T050000.000 GMT", lalu +1 bulan

import (
	"testing"
	"time"

	"nusantarare/inti/backend/utils"
)

const batasTiruan = 25

func TestPraprosesTanggalBatasTutupBuku(t *testing.T) {
	for _, tt := range []struct {
		saat  time.Time
		harap string
	}{
		{time.Date(2026, 10, 24, 9, 30, 0, 0, time.UTC), "2026-10-24 09:30:00"},
		{time.Date(2026, 10, 25, 9, 30, 0, 0, time.UTC), "2026-10-25 09:30:00"}, // hari = batas: tetap
		{time.Date(2026, 10, 26, 9, 30, 0, 0, time.UTC), "2026-11-01 09:30:00"},
		{time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC), "2027-01-01 23:59:59"}, // pergantian tahun
		{time.Date(2027, 1, 31, 8, 0, 0, 0, time.UTC), "2027-02-01 08:00:00"},     // bulan pendek
	} {
		h := HalamanBaru()
		PraprosesTanggal(h, tt.saat, batasTiruan)
		if got := h.Ambil("PolicyTreatyIn.StatementDate"); got != utils.FormatTanggalWaktu(tt.saat) {
			t.Errorf("%s: StatementDate = %q", tt.saat, got)
		}
		if got := h.Ambil("PolicyTreatyIn.ProductionDate"); got != tt.harap {
			t.Errorf("%s: ProductionDate = %q, harap %q", tt.saat, got, tt.harap)
		}
	}
}

// Batas bukan 25 tertanam: dengan batas 26, hari 26 tidak digeser.
func TestPraprosesTanggalMengikutiBatasTabel(t *testing.T) {
	h := HalamanBaru()
	saat := time.Date(2026, 10, 26, 9, 30, 0, 0, time.UTC)
	PraprosesTanggal(h, saat, 26)
	if got := h.Ambil("PolicyTreatyIn.ProductionDate"); got != "2026-10-26 09:30:00" {
		t.Fatalf("batas 26, hari 26: ProductionDate = %q, harap tetap", got)
	}
}

func TestTanggalProduksiNomorBatasTutupBuku(t *testing.T) {
	jkt := zonaJakarta()
	for _, tt := range []struct {
		hari  int
		harap time.Time
	}{
		{24, time.Date(2026, 10, 24, 10, 0, 0, 0, jkt)},
		{25, time.Date(2026, 10, 25, 10, 0, 0, 0, jkt)}, // hari = batas: tetap
		{26, time.Date(2026, 11, 1, 5, 0, 0, 0, time.UTC)},
	} {
		sekarang := time.Date(2026, 10, tt.hari, 10, 0, 0, 0, jkt)
		if got := TanggalProduksiNomor(sekarang, sekarang, batasTiruan); !got.Equal(tt.harap) {
			t.Errorf("hari %d: %v, harap %v", tt.hari, got, tt.harap)
		}
	}
	// 5.3 membaca hari di Asia/Jakarta: 25 Oktober 18:00 GMT = 26 Oktober 01:00 WIB -> digeser.
	malam := time.Date(2026, 10, 25, 18, 0, 0, 0, time.UTC)
	if got := TanggalProduksiNomor(malam, malam, batasTiruan); !got.Equal(time.Date(2026, 11, 1, 5, 0, 0, 0, time.UTC)) {
		t.Errorf("25 Okt 18:00 GMT (26 Okt WIB): %v", got)
	}
}
