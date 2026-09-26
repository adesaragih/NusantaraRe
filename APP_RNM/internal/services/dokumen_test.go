package services_test

// Gerbang dokumen - TANPA Oracle.
//
// Pemilik: tiket 03. Dibaca sesudah: dokumen.go.

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/internal/models"
	"nusantarare/internal/services"
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

// TestKelengkapanMembandingkanCACAHKategoriBerbeda mengunci langkah 12.2.2
// sampai 12.2.4: kategori di-dedup, lalu CACAHnya dibandingkan.
func TestKelengkapanMembandingkanCacahKategoriBerbeda(t *testing.T) {
	wajib := []string{"UJI-A", "UJI-B"}

	// Dua dokumen berkategori SAMA: sesudah dedup cacahnya 1, bukan 2.
	kembar := []models.Peserta{pesertaDipilih(
		models.Dokumen{Kategori2: "UJI-A"},
		models.Dokumen{Kategori2: "UJI-A"},
	)}
	if err := services.PeriksaDokumenLengkap(kembar, wajib); !errors.Is(
		err, services.ErrDokumenTidakLengkap) {
		t.Fatalf("galat = %v, mau ErrDokumenTidakLengkap; dua dokumen "+
			"berkategori sama tidak menutup dua kategori wajib", err)
	}

	lengkap := []models.Peserta{pesertaDipilih(
		models.Dokumen{Kategori2: "UJI-A"},
		models.Dokumen{Kategori2: "UJI-B"},
	)}
	if err := services.PeriksaDokumenLengkap(lengkap, wajib); err != nil {
		t.Errorf("dua kategori berbeda ditolak: %v", err)
	}
}

// TestKategoriWajibKosongMenggagalkan - nol kategori wajib berarti aturannya
// belum ada, bukan berarti semua orang lengkap.
func TestKategoriWajibKosongMenggagalkan(t *testing.T) {
	peserta := []models.Peserta{pesertaDipilih(models.Dokumen{Kategori2: "UJI-A"})}
	err := services.PeriksaDokumenLengkap(peserta, nil)
	if !errors.Is(err, services.ErrKategoriWajibBelumDiketahui) {
		t.Fatalf("galat = %v, mau ErrKategoriWajibBelumDiketahui", err)
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

// TestKategoriKosongTidakDihitung - dokumen tanpa kategori bukan kategori.
func TestKategoriKosongTidakDihitung(t *testing.T) {
	p := pesertaDipilih(
		models.Dokumen{Kategori2: "UJI-A"},
		models.Dokumen{Kategori2: "   "},
		models.Dokumen{},
	)
	if got := services.KategoriBerbeda(p); len(got) != 1 || got[0] != "UJI-A" {
		t.Errorf("KategoriBerbeda = %v, mau [UJI-A]", got)
	}
}

// TestPesanGerbangKeduaPersisSepertiXML - pesan galat adalah logika bisnis,
// dan XML hanya menuliskan kalimat ini. Nomor peserta TIDAK ada di dalamnya;
// ia disediakan terpisah supaya layar tetap dapat menunjukkannya.
func TestPesanGerbangKeduaPersisSepertiXML(t *testing.T) {
	wajib := []string{"UJI-A", "UJI-B"}
	peserta := []models.Peserta{
		pesertaDipilih(models.Dokumen{Kategori2: "UJI-A"}, models.Dokumen{Kategori2: "UJI-B"}),
		pesertaDipilih(models.Dokumen{Kategori2: "UJI-A"}),
	}
	err := services.PeriksaDokumenLengkap(peserta, wajib)
	if err == nil {
		t.Fatal("peserta kedua lolos padahal kategorinya kurang")
	}
	const mau = "Documents are incomplete, please complete the documents"
	if err.Error() != mau {
		t.Errorf("pesan = %q, mau persis %q", err.Error(), mau)
	}
	nomor := services.PesertaDokumenTidakLengkap(peserta, wajib)
	if len(nomor) != 1 || nomor[0] != 2 {
		t.Errorf("nomor peserta gagal = %v, mau [2]", nomor)
	}
}
