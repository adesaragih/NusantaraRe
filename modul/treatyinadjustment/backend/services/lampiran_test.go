package services_test

// Uji panel Attachment modul Adjustment - NOL koneksi Oracle.
//
// ⛔ Aturannya SAMA PERSIS dengan modul Treaty In, dan uji ini ada supaya
// kesamaan itu TERJAGA di kedua tempat. Menebak pasangan kode↔nama di salah
// satu modul saja sudah cukup untuk menaruh berkas di kategori yang salah.

import (
	"context"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyinadjustment/backend/models"
	"nusantarare/modul/treatyinadjustment/backend/services"
)

func katalogUji() map[string]string {
	return map[string]string{
		"00000": "Others",
		"00001": "Analysed Email",
		"00002": "Approval Email",
		"00005": "Summary Treaty Leader",
		"00006": "Assessment Inward Treaty Form / Format Analisa Treaty",
		"00007": "Pega Proportional Calculation /Perhitungan Pega Proportional",
		"00010": "Offer Email",
	}
}

// ⛔ UJI PALING PENTING: keempat kode tanpa nama TIDAK ditebak namanya.
func TestKeempatKodeTanpaNamaTIDAKDitebakNamanya(t *testing.T) {
	baris := services.SusunKategoriLampiran(katalogUji(), nil)

	belum := 0
	for _, b := range baris {
		if b.Dipastikan {
			continue
		}
		belum++
		if b.Nama != "" {
			t.Errorf("kode %s diberi nama %q - pasangannya TIDAK diketahui, dan "+
				"menebaknya menaruh berkas di kategori yang salah", b.Kode, b.Nama)
		}
	}
	if belum != 4 {
		t.Errorf("%d kategori belum dipastikan, mau 4", belum)
	}
	// Keempat nama yang belum berumah tidak boleh bocor ke baris mana pun.
	for _, b := range baris {
		for _, n := range models.NamaKategoriBelumDipastikan {
			if b.Nama == n {
				t.Errorf("kode %s diberi nama %q yang kodenya belum diketahui", b.Kode, n)
			}
		}
	}
}

// Kategori BERNOL berkas tetap tampil - kolom `Count` tidak akan pernah
// berbunyi `0` kalau barisnya disembunyikan saat kosong.
func TestKategoriNolBerkasTetapTampil(t *testing.T) {
	lampiran := []models.BarisLampiranWarisan{{KodeKategori: "00001"}, {KodeKategori: "00001"}}
	baris := services.SusunKategoriLampiran(katalogUji(), lampiran)
	if len(baris) != 11 {
		t.Fatalf("%d baris, mau 11 (7 terbukti + 4 belum)", len(baris))
	}
	cacah := map[string]int{}
	for _, b := range baris {
		cacah[b.Kode] = b.Cacah
	}
	if cacah["00001"] != 2 || cacah["00000"] != 0 || cacah["00003"] != 0 {
		t.Errorf("cacah salah: %v", cacah)
	}
}

// ⭐ Spanduk biru itu ATURAN. Diuji DUA ARAH.
func TestNamaBerkasAmanDuaArah(t *testing.T) {
	for _, n := range []string{
		"laporan treaty.pdf", "laporan(1).pdf", "a/b.pdf", "a\b.pdf",
		"../rahasia.pdf", "surat&lampiran.pdf", "", "   ", ".", "..",
	} {
		if services.NamaBerkasAman(n) {
			t.Errorf("%q DITERIMA; spanduk menyebut titik atau garis bawah sebagai pengganti aman", n)
		}
	}
	// Pemeriksa yang menolak SEGALANYA lulus uji di atas dengan sempurna,
	// dan menutup pintu yang sistem lama buka.
	for _, n := range []string{
		"laporan.pdf", "laporan_treaty.pdf", "laporan-treaty.pdf",
		"laporan.2026.pdf", "LAPORAN.PDF", "Treaty_In_2026-01.xlsx",
	} {
		if !services.NamaBerkasAman(n) {
			t.Errorf("%q DITOLAK; ia hanya huruf, angka, titik, garis bawah, tanda hubung", n)
		}
	}
}

// Tanpa identitas ditolak SEBELUM gudang disentuh.
func TestLampiranMenolakTanpaIdentitasSebelumGudang(t *testing.T) {
	g := &gudangTiruan{katalogKategori: katalogUji()}
	l := services.LayananDengan(g)

	if _, err := l.LampiranKontrakWarisan(context.Background(), inti.Pelaku{}, "1001851"); err == nil {
		t.Fatal("permintaan tanpa identitas DITERIMA")
	}
}

// Pengenal kosong ditolak - "tidak ada" berbeda dari "tidak ditanyakan".
func TestLampiranMenolakPengenalKosong(t *testing.T) {
	g := &gudangTiruan{katalogKategori: katalogUji()}
	l := services.LayananDengan(g)

	for _, id := range []string{"", "   "} {
		if _, err := l.LampiranKontrakWarisan(context.Background(), pelakuAda, id); err == nil {
			t.Errorf("pengenal %q DITERIMA", id)
		}
	}
}

// Uji POSITIF: kategori tersusun, cacahnya benar, berkasnya diteruskan.
func TestLampiranPositif(t *testing.T) {
	g := &gudangTiruan{
		katalogKategori: katalogUji(),
		lampiran: []models.BarisLampiranWarisan{
			{ID: "1", KodeKategori: "00001", NamaBerkas: "a.pdf"},
			{ID: "2", KodeKategori: "00001", NamaBerkas: "b.pdf"},
		},
	}
	l := services.LayananDengan(g)

	hasil, err := l.LampiranKontrakWarisan(context.Background(), pelakuAda, "1001851")
	if err != nil {
		t.Fatalf("mau diterima, dapat %v", err)
	}
	if len(hasil.Berkas) != 2 {
		t.Errorf("%d berkas, mau 2", len(hasil.Berkas))
	}
	if len(hasil.Kategori) != 11 {
		t.Errorf("%d kategori, mau 11", len(hasil.Kategori))
	}
	for _, k := range hasil.Kategori {
		if k.Kode == "00001" && k.Cacah != 2 {
			t.Errorf("kategori 00001 cacahnya %d, mau 2", k.Cacah)
		}
	}
}
