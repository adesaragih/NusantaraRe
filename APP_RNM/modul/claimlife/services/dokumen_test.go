package services_test

// Gerbang dokumen - TANPA Oracle.
//
// Pemilik: tiket 03. Dibaca sesudah: dokumen.go.

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/claimlife/models"
	"nusantarare/modul/claimlife/services"
)

func pesertaDipilih(dok ...models.Dokumen) models.Peserta {
	return models.Peserta{IsCheck: "true", Dokumen: dok}
}

// TestPesertaTanpaDokumenDitolakDenganNomorUrut mengunci langkah 3.1: pesannya
// wajib menyebut peserta MANA, dan nomornya nomor urut 1-based.
func TestPesertaTanpaDokumenDitolakDenganNomorUrut(t *testing.T) {
	peserta := []models.Peserta{
		pesertaDipilih(models.Dokumen{Kategori2: "UJI-KTP"}),
		pesertaDipilih(), // peserta kedua, nol dokumen
	}
	err := services.PeriksaDokumenAda("QR", peserta)
	if !errors.Is(err, services.ErrDokumenBelumDiunggah) {
		t.Fatalf("galat = %v, mau ErrDokumenBelumDiunggah", err)
	}
	if !strings.Contains(err.Error(), "person number 2") {
		t.Errorf("pesan tidak menyebut peserta kedua: %v", err)
	}
	if strings.Contains(err.Error(), "person number 1") {
		t.Errorf("pesan menuduh peserta pertama yang dokumennya ada: %v", err)
	}
}

// TestTipeTreatyMelewatiGerbangDokumen mengunci precondition
// `Type=="TP"||Type=="TR"` dengan WhenTrue=3 - benar berarti LEWATI.
func TestTipeTreatyMelewatiGerbangDokumen(t *testing.T) {
	peserta := []models.Peserta{pesertaDipilih()}
	for _, tipe := range []string{"TP", "TR", "tp"} {
		if err := services.PeriksaDokumenAda(tipe, peserta); err != nil {
			t.Errorf("Type %q: %v; treaty tidak dituntut berdokumen", tipe, err)
		}
	}
	if err := services.PeriksaDokumenAda("QP", peserta); err == nil {
		t.Error("Type QP lolos gerbang; hanya TP dan TR yang dilewati")
	}
}

// TestPesertaTidakDipilihDilewati - gerbangnya hanya mengenai peserta yang
// benar-benar diklaim.
func TestPesertaTidakDipilihDilewati(t *testing.T) {
	peserta := []models.Peserta{{IsCheck: "", Dokumen: nil}}
	if err := services.PeriksaDokumenAda("QR", peserta); err != nil {
		t.Errorf("peserta yang tidak dipilih ikut dituntut: %v", err)
	}
}

// TestSumberKategoriBawaanGagalTerang - yang belum diketahui harus terlihat
// sebagai satu galat yang menyebut apa yang ditunggu, bukan sebagai daftar
// kosong yang diam-diam meloloskan semua.
func TestSumberKategoriBawaanGagalTerang(t *testing.T) {
	_, err := services.KategoriWajibBelumDiketahui{}.KategoriWajib(context.Background())
	if !errors.Is(err, services.ErrKategoriWajibBelumDiketahui) {
		t.Fatalf("galat = %v, mau ErrKategoriWajibBelumDiketahui", err)
	}
	if !strings.Contains(err.Error(), "GetCategoryLife_SQL") {
		t.Errorf("pesan tidak menyebut rule yang ditunggu: %v", err)
	}
}
