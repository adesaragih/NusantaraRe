//go:build db

package repository_test

// Bukti Oracle untuk isi kedua pemilih "Choose …".
//
// ⛔ BACA SAJA. Nol transaksi, nol tulisan; tidak ada yang perlu dibatalkan.
//
// ⛔ Yang dibuktikan di sini bukan "ada datanya" melainkan DUA hal yang
// menjadi dasar keputusan 4 Oktober 2026 menyalakan tombolnya:
//
//  1. nilainya memang ada — alasan lama "tidak ada yang dapat dipilih" keliru;
//  2. satu nama memang dapat ber-pengenal lebih dari satu — alasan mengapa
//     layar WAJIB menandainya alih-alih menyatukannya.
//
// Bila suatu hari butir 2 menjadi nol, penandaan di layar boleh dicabut —
// dan uji ini yang akan memberitahu.

import (
	"context"
	"database/sql"
	"testing"

	"nusantarare/inti/backend/config"
)

func bukaOracle(t *testing.T) (*sql.DB, string) {
	t.Helper()
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("konfigurasi: %v", err)
	}
	if cfg.IsPegaProd {
		t.Fatal("menolak berjalan saat IS_PEGA_PROD=true (ADR-U-0005)")
	}
	if !cfg.PunyaOracle() {
		t.Skip("lewati: ORACLE_DSN kosong")
	}
	db, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db, cfg.OracleSchema
}

// Kuerinya SAMA PERSIS dengan `bacaPilihan` di `warisan_pemilih.go`. Salinan
// yang menyimpang diam-diam membuktikan kueri yang tidak dijalankan siapa pun.
func cacahPilihan(t *testing.T, db *sql.DB, skema, kolomID, kolomNama string) int {
	t.Helper()
	q := `SELECT COUNT(*) FROM (SELECT DISTINCT ` + kolomID + `, ` + kolomNama +
		` FROM ` + skema + `.TREATY_IN WHERE ` + kolomNama + ` IS NOT NULL)`
	var n int
	if err := db.QueryRowContext(context.Background(), q).Scan(&n); err != nil {
		t.Fatalf("menghitung pilihan %s: %v", kolomNama, err)
	}
	return n
}

// Kolomnya ADA dan terisi. Nama kolomnya bukan `CEDANT`/`ASAL_BISNIS`
// seperti yang ERD sebut, melainkan `CEDING`/`LEADINGREINSSOURCE` —
// perbedaan itu yang dulu membuat orang menyimpulkan datanya tidak ada.
func TestPilihanWarisanAdaIsinya(t *testing.T) {
	db, skema := bukaOracle(t)

	ced := cacahPilihan(t, db, skema, "CEDINGID", "CEDING")
	asal := cacahPilihan(t, db, skema, "LEADINGREINSSOURCEID", "LEADINGREINSSOURCE")

	if ced == 0 || asal == 0 {
		t.Fatalf("pemilih kosong — tombolnya tidak layak hidup: cedant=%d asal=%d", ced, asal)
	}
	// Terukur 4 Oktober 2026: 94 dan 91. Dinyatakan sebagai BATAS BAWAH,
	// bukan angka tepat: cedant baru boleh bertambah, dan uji yang pecah
	// karena data bertambah hanya mengajari orang mengabaikannya.
	if ced < 90 || asal < 85 {
		t.Errorf("pilihan menyusut jauh di bawah ukuran 4 Okt 2026 (94/91): cedant=%d asal=%d", ced, asal)
	}
	t.Logf("cedant=%d  asal bisnis=%d", ced, asal)
}

// ⚠️ Inilah sebab layar menandai, bukan menyatukan.
func TestNamaCedantMasihBerpengenalGanda(t *testing.T) {
	db, skema := bukaOracle(t)

	q := `SELECT COUNT(*) FROM (
	        SELECT CEDING FROM ` + skema + `.TREATY_IN
	         WHERE CEDING IS NOT NULL
	         GROUP BY CEDING HAVING COUNT(DISTINCT CEDINGID) > 1)`
	var n int
	if err := db.QueryRowContext(context.Background(), q).Scan(&n); err != nil {
		t.Fatalf("menghitung nama kembar: %v", err)
	}
	if n == 0 {
		t.Log("nol nama kembar — penandaan `kembar` di layar boleh ditinjau ulang")
		return
	}
	t.Logf("%d nama cedant ber-pengenal lebih dari satu — penandaan WAJIB tetap ada", n)
}
