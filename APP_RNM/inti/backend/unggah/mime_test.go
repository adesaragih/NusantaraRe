package unggah

import (
	"testing"
)

// Uji aturan lahirnya satu baris dokumen — kelompok Dokumen.

func TestMimeDariNamaFile(t *testing.T) {
	t.Run("tabel keputusan disalin utuh", func(t *testing.T) {
		// ⛔ Cacahnya dikunci. Baris yang HILANG tidak berbunyi sendiri:
		// berkasnya diam-diam menjadi application/octet-stream, dan tidak
		// ada yang tahu ia pernah punya jenis sendiri.
		if len(petaMime) != JumlahBarisMime {
			t.Fatalf("petaMime = %d baris, mau %d", len(petaMime), JumlahBarisMime)
		}
		// Beberapa baris diperiksa harfiah, termasuk yang tampak ganjil -
		// justru yang ganjil itu yang paling mudah "dirapikan" orang.
		for ext, mau := range map[string]string{
			"pdf":  "application/pdf",
			"jfif": "image/jpeg",
			"xlsb": "application/vnd.ms-excel.sheet.binary.macroEnabled.12",
			"et":   "application/et",
			"lnk":  "application/x-ms-shortcut",
			"mht":  "message/rfc822",
			"avi":  "video/x-msvideo",
		} {
			if got := MimeDariNamaFile("UJI-berkas." + ext); got != mau {
				t.Errorf("MimeDariNamaFile(.%s) = %q, mau %q", ext, got, mau)
			}
		}
	})

	t.Run("ekstensi diambil dari titik TERAKHIR", func(t *testing.T) {
		// Nama berkas bertitik banyak lumrah pada lampiran klaim.
		if got := MimeDariNamaFile("UJI-laporan.2026.pdf"); got != "application/pdf" {
			t.Errorf("= %q, mau application/pdf", got)
		}
	})

	t.Run("huruf besar tetap dikenali", func(t *testing.T) {
		if got := MimeDariNamaFile("UJI-LAPORAN.PDF"); got != "application/pdf" {
			t.Errorf("= %q, mau application/pdf", got)
		}
	})

	t.Run("tak dikenal mendapat bawaan, bukan galat", func(t *testing.T) {
		// ⛔ Tabelnya punya `otherwise` (b89). Menolak berkas yang sistem
		// lama terima berarti menutup pintu yang terbuka.
		for _, nama := range []string{"UJI-berkas.xyz", "UJI-tanpa-titik", "UJI-berkas."} {
			if got := MimeDariNamaFile(nama); got != MimeBawaan {
				t.Errorf("MimeDariNamaFile(%q) = %q, mau %q", nama, got, MimeBawaan)
			}
		}
	})
}
