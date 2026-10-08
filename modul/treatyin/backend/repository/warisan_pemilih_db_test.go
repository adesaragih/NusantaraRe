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
	"nusantarare/modul/treatyin/backend/repository"
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

// ⛔ RALAT 6 Oktober 2026 — sumber pemilih DIGANTI ke daftar agen Pega
// (`BrowseAgentNusaRe_RD`), dan uji di bawah mengikutinya.
//
// Kuerinya memakai `repository.SaringanAgenPega` — konstanta YANG SAMA dengan
// yang `bacaAgenAktif` jalankan, bukan salinannya. Bentuk sebelumnya menyalin
// teks kuerinya, dan salinan itu akan terus lulus sesudah kueri aslinya
// berubah.
func cacahAgenPega(t *testing.T, db *sql.DB, skema, saringanTambahan string) int {
	t.Helper()
	q := `SELECT COUNT(*) FROM ` + skema + `.` + repository.TabelAgen +
		` WHERE ` + repository.SaringanAgenPega + saringanTambahan
	var n int
	if err := db.QueryRowContext(context.Background(), q).Scan(&n); err != nil {
		t.Fatalf("menghitung agen: %v", err)
	}
	return n
}

// ⭐ Pemilih berisi, dan ukurannya ukuran PEGA — bukan ukuran nilai yang
// pernah dipakai.
func TestPilihanWarisanAdaIsinya(t *testing.T) {
	db, skema := bukaOracle(t)
	n := cacahAgenPega(t, db, skema, "")
	if n == 0 {
		t.Fatal("nol agen lolos saringan Pega — pemilih Ceding/SoB kosong")
	}
	// Terukur 6 Oktober 2026: 255 dari 429 baris `AGENT`. BATAS BAWAH, bukan
	// angka tepat: agen baru boleh bertambah, dan yang menyusut jauh berarti
	// saringannya rusak — misalnya ejaan `AGENTTPYE2` yang "dibetulkan".
	if n < 240 {
		t.Errorf("agen lolos saringan menyusut ke %d (terukur 255 pada 6 Okt 2026)", n)
	}
	t.Logf("agen yang Pega tawarkan: %d", n)
}

// ⭐ Ruang pengenalnya `AGENT.ID` — DIUKUR, dan uji ini yang menguncinya.
//
// Bila suatu hari pengenal tersimpan berhenti cocok dengan `AGENT.ID`, maka
// memilih di layar menulis pengenal dari ruang yang SALAH ke kontrak.
func TestPengenalTersimpanAdalahAgentID(t *testing.T) {
	db, skema := bukaOracle(t)
	for _, k := range []string{"CEDINGID", "LEADINGREINSSOURCEID"} {
		q := `SELECT COUNT(*) FROM ` + skema + `.TREATY_IN t WHERE EXISTS (SELECT 1 FROM ` +
			skema + `.` + repository.TabelAgen + ` a WHERE TO_CHAR(a.ID) = TO_CHAR(t.` + k + `))`
		var n int
		if err := db.QueryRowContext(context.Background(), q).Scan(&n); err != nil {
			t.Fatalf("%s: %v", k, err)
		}
		// Terukur 6 Oktober 2026 atas 1.854 baris: CEDINGID 1.744,
		// LEADINGREINSSOURCEID 1.854.
		if n < 1700 {
			t.Errorf("%s cocok AGENT.ID hanya pada %d baris — ruang pengenalnya berubah", k, n)
		}
		t.Logf("%s = AGENT.ID pada %d baris", k, n)
	}
}

// ⚠️ Inilah sebab layar menandai, bukan menyatukan.
//
// Terukur 6 Oktober 2026 di bawah saringan Pega: 4 nama dipakai lebih dari
// satu `AGENT.ID` (8 baris). Satu pun cukup untuk membuat pencarian balik
// nama->ID menunjuk agen yang keliru.
func TestNamaCedantMasihBerpengenalGanda(t *testing.T) {
	db, skema := bukaOracle(t)

	q := `SELECT COUNT(*) FROM (
	        SELECT CLIENTNAME FROM ` + skema + `.` + repository.TabelAgen + `
	         WHERE ` + repository.SaringanAgenPega + `
	         GROUP BY CLIENTNAME HAVING COUNT(DISTINCT ID) > 1)`
	var n int
	if err := db.QueryRowContext(context.Background(), q).Scan(&n); err != nil {
		t.Fatalf("menghitung nama kembar: %v", err)
	}
	if n == 0 {
		t.Log("nol nama kembar — penandaan `kembar` di layar boleh ditinjau ulang")
		return
	}
	t.Logf("%d nama agen ber-pengenal lebih dari satu — penandaan WAJIB tetap ada", n)
}
