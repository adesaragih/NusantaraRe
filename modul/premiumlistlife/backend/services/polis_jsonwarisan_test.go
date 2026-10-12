package services

// JSON_OFFER_LIFE (Save Offer) dan JSON_POLIS (Confirm) - keputusan work owner
// 06-10-2026. TANPA Oracle: penjaga statik urutan dan transaksi (alasan sama
// dengan polis_summary_test.go). ⚠️ Bukti tingkat Oracle belum ada - kedua tabel
// warisan tidak dipasang skema uji.

import (
	"os"
	"strings"
	"testing"
)

// urutDalam - setiap jejak ada di badan SESUDAH jejak sebelumnya (dicari berurutan).
func urutDalam(t *testing.T, badan string, urut []string) {
	t.Helper()
	badan = strings.ReplaceAll(badan, "\r\n", "\n")
	dari := 0
	for _, jejak := range urut {
		i := strings.Index(badan[dari:], jejak)
		if i < 0 {
			t.Fatalf("%s tidak ada sesudah langkah sebelumnya; urutan terkunci: %v", jejak, urut)
		}
		dari += i + len(jejak)
	}
}

// TestSaveOfferMenulisJSONOfferSatuTransaksi - T_PREMIUM_LIST lalu JSON_OFFER_LIFE di
// DalamTransaksi yang sama; galatnya dikembalikan (transaksi batal), nol galat ditelan.
func TestSaveOfferMenulisJSONOfferSatuTransaksi(t *testing.T) {
	isi, err := os.ReadFile("polis_isianpenawaran.go")
	if err != nil {
		t.Fatal(err)
	}
	badan := badanFungsi(t, string(isi), "func (f *FormPenawaran) Simpan(")
	urutDalam(t, badan, []string{
		"f.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {",
		"repo.PindahkanBarisStatus(ctx, tx,",
		"repo.Simpan(ctx, tx, polisID, siap)",
		"repository.NewJSONOfferLife(f.svc.DB()).Ganti(ctx, tx, polisID)",
		"ErrPenawaranMelebihiKolom",
		"return err",
		"repo.SisipSuggest(ctx, tx,",
	})
	if strings.Count(badan, "DalamTransaksi(") != 1 {
		t.Error("Save Offer harus SATU transaksi")
	}
}

// TestConfirmMenulisJSONPolisSesudahProduksi - JSON_POLIS sesudah nomor terbit dan
// LIFEINPRODUCTION, di transaksi Confirm yang sama; galat membatalkan Confirm.
func TestConfirmMenulisJSONPolisSesudahProduksi(t *testing.T) {
	badan := badanFungsi(t, sumberSummary(t), "func (s *SummaryPremiumList) simpanDalam(")
	urutDalam(t, badan, []string{
		"s.nomor.terbitkanDalam(",
		"NewProduksiLife(s.svc.DB()).Ganti(ctx, tx, polisID, nomor.Nomor)",
		"tglProd, err := s.tanggalProduksi(ctx, tx, polisID, nomor.Periode, saat)",
		"return HasilSubmitSummary{}, err",
		"repository.NewJSONPolis(s.svc.DB()).Ganti(ctx, tx, polisID, nomor.Nomor, tglProd); err != nil {\n\t\treturn HasilSubmitSummary{}, err",
	})
	if strings.Contains(badan, "DalamTransaksi(") {
		t.Error("simpanDalam memakai transaksi pemanggil, bukan transaksi baru")
	}
	// TGL_PROD: nomor baru → periodenya; nomor lama → TGL_PROD lama, lalu aturan PL Number.
	tp := badanFungsi(t, sumberSummary(t), "func (s *SummaryPremiumList) tanggalProduksi(")
	urutDalam(t, tp, []string{"TglProdLama(", "HariClosing(", "penomor.HitungPeriodeNomor(saat, hari)",
		"models.TanggalProduksi(saat, periode)"})
}
