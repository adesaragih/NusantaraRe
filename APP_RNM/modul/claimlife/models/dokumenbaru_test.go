package models

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"nusantarare/inti/unggah"
)

func TestMimeDokumen(t *testing.T) {
	// ⛔ Padanan prasyarat b586 `Param.MIME==""` WhenTrue=2 LANJUT: tabel
	// hanya mengisi yang KOSONG. Yang disebut pemanggil menang.
	t.Run("yang disebut pemanggil MENANG atas tabel", func(t *testing.T) {
		got := unggah.MimeDokumen("image/png", "UJI-berkas.pdf")
		if got != "image/png" {
			t.Errorf("= %q, mau image/png - tabel tidak boleh menimpa", got)
		}
	})
	t.Run("kosong diisi dari nama berkas", func(t *testing.T) {
		if got := unggah.MimeDokumen("", "UJI-berkas.pdf"); got != "application/pdf" {
			t.Errorf("= %q, mau application/pdf", got)
		}
		if got := unggah.MimeDokumen("   ", "UJI-berkas.pdf"); got != "application/pdf" {
			t.Errorf("spasi = %q, mau application/pdf", got)
		}
	})
	t.Run("hasilnya selalu huruf kecil", func(t *testing.T) {
		// Padanan `@toLowerCase(Param.MIME)` b782. Tanpa itu
		// "APPLICATION/PDF" dan "application/pdf" menjadi dua jenis berbeda
		// di kolom yang sama.
		if got := unggah.MimeDokumen("APPLICATION/PDF", ""); got != "application/pdf" {
			t.Errorf("= %q, mau application/pdf", got)
		}
	})
}

func TestKunciKelompokDokumen(t *testing.T) {
	saat := time.Date(2026, 3, 17, 14, 5, 9, 123000000, time.UTC)

	t.Run("lahir berawalan DL- dan hanya angka sesudahnya", func(t *testing.T) {
		got := KunciKelompokDokumen("", saat)
		if !strings.HasPrefix(got, AwalanKunciDokumen) {
			t.Fatalf("%q tidak berawalan %q", got, AwalanKunciDokumen)
		}
		sisa := strings.TrimPrefix(got, AwalanKunciDokumen)
		if regexp.MustCompile(`[^0-9]`).MatchString(sisa) {
			t.Errorf("%q memuat bukan angka; b596 membuang [^0-9]", sisa)
		}
		if sisa == "" {
			t.Error("bagian angkanya kosong")
		}
	})

	t.Run("yang SUDAH ADA tidak pernah ditimpa", func(t *testing.T) {
		// ⛔ Inti berkas ini. Menimpanya pada lampiran kedua memutus lampiran
		// pertama dari pesertanya: baris lama tetap menyimpan kunci lama di
		// KATEGORI_1, dan saringan layar hanya mencari kunci baru. Dokumen
		// itu tidak hilang dari basis data - ia hilang dari LAYAR.
		lama := "DL-20260101000000000"
		if got := KunciKelompokDokumen(lama, saat); got != lama {
			t.Errorf("= %q, mau %q apa adanya", got, lama)
		}
	})
}

func TestIDDokumenBaru(t *testing.T) {
	jakarta := time.FixedZone("WIB", 7*60*60)
	saat := time.Date(2026, 3, 17, 7, 5, 9, 123000000, time.UTC) // 14:05:09 WIB

	got := IDDokumenBaru(saat, jakarta)
	if regexp.MustCompile(`[^0-9]`).MatchString(got) {
		t.Errorf("%q memuat bukan angka", got)
	}
	if len(got) != 17 {
		t.Errorf("panjang %q = %d, mau 17 (yyyyMMdd + hhmmss + SSS)", got, len(got))
	}
	if !strings.HasPrefix(got, "20260317") {
		t.Errorf("%q tidak berawalan tanggal Jakarta 20260317", got)
	}

	// ⚠️ CACAT WARISAN, DITIRU DAN DINYATAKAN: `hh` adalah jam 12-jam, jadi
	// 14:05 dan 02:05 menghasilkan pengenal yang SAMA. Uji ini mengunci
	// cacat itu supaya ia tidak diam-diam "diperbaiki" - memperbaikinya
	// membuat pengenal baris baru berbeda bentuk dari pengenal baris lama,
	// dan keduanya hidup di kolom yang sama.
	pagi := time.Date(2026, 3, 16, 19, 5, 9, 123000000, time.UTC) // 02:05:09 WIB 17-03
	if IDDokumenBaru(pagi, jakarta) != got {
		t.Errorf("jam 12-jam warisan tidak lagi ditiru: %q vs %q",
			IDDokumenBaru(pagi, jakarta), got)
	}

	// Zona waktunya Jakarta, bukan zona server.
	utc := IDDokumenBaru(saat, time.UTC)
	if utc == got {
		t.Error("zona waktu tidak berpengaruh; b648 menyebut Asia/Jakarta")
	}
}

func TestUrutanUnggahDanSimpan(t *testing.T) {
	// ⛔ Unggah DULU, baris basis data KEMUDIAN (prasyarat b1283,
	// WhenTrue=3 LEWATI). Membaliknya menghasilkan baris yang menunjuk
	// berkas yang tidak pernah naik - layar menampilkannya sebagai dokumen
	// yang ada, tautannya gagal, dan tidak ada yang dapat membedakannya
	// dari gangguan jaringan sesaat.
	if BolehSimpanBarisDokumen("") {
		t.Error("baris disimpan tanpa T_STORAGE_ID")
	}
	if BolehSimpanBarisDokumen("   ") {
		t.Error("spasi diterima sebagai penunjuk penyimpanan")
	}
	if !BolehSimpanBarisDokumen("UJI-STORAGE-1") {
		t.Error("penunjuk yang sah ditolak")
	}
}

func TestUrutanHapus(t *testing.T) {
	// ⛔ Penghapusan di penyimpanan DILEWATI bila tidak ada berkasnya
	// (prasyarat b472), tetapi `Obj-Delete` b513 berjalan TANPA prasyarat:
	// barisnya tetap dihapus. Arah yang benar - baris yatim tanpa berkas
	// tidak berguna bagi siapa pun.
	if PerluHapusDiPenyimpanan("") {
		t.Error("penghapusan penyimpanan dipanggil untuk baris tanpa berkas")
	}
	if !PerluHapusDiPenyimpanan("UJI-STORAGE-1") {
		t.Error("berkas yang ada tidak ikut dihapus - ia akan menjadi yatim di penyimpanan")
	}
}
