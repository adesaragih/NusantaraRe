//go:build db

package repository_test

// Bukti Oracle untuk migrasi `451` — kolom `REVISIONDATE`.
//
// ---------------------------------------------------------------------
// ⛔ YANG DIBUKTIKAN
// ---------------------------------------------------------------------
// Laporan pemakai 8 Oktober 2026 sesudah Save di layar Treaty In Adjustment:
//
//	Data Sudah Disimpan Dengan ID : 1000001/R01
//	TIDAK tersimpan (belum punya kolom/tabel): RevisionDate
//
// `TreatyInEDMNew` langkah [3] menyetel `TreatyIn.RevisionDate =
// @CurrentDateTime()` pada SETIAP draf penyesuaian, jadi setiap Save
// melaporkannya hilang.
//
// ⛔ Satu transaksi, diakhiri `Rollback`. Nol `Commit`, nol teardown.
//
// ⛔ NOL SENTUHAN `M_TREATY_IN` DAN NOL BACA JSON — dokumennya disusun di
// sini sebagai data Go, pengenal kontraknya karangan.

import (
	"database/sql"
	"testing"

	"nusantarare/inti/backend/config"
	"nusantarare/modul/treatyin/backend/repository"
)

const idUjiRevisionDate = "UJI-REVISIONDATE/R01"

// skemaRev - skema Oracle yang sedang dipakai.
func skemaRev(t *testing.T) string {
	t.Helper()
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("konfigurasi: %v", err)
	}
	return cfg.OracleSchema
}

// ⭐ Nilainya MENDARAT, dan dibaca kembali apa adanya.
func TestRevisionDateMendaratDiTabelRevisi(t *testing.T) {
	g, tx, ctx := gudangUji(t)

	const stempel = "20261008T083015.000 GMT"
	doc := map[string]any{
		"ID":                 idUjiRevisionDate,
		"TreatyContractName": "KONTRAK UJI",
		// ⭐ Persis yang `SusunDraf` tulis.
		"RevisionDate": stempel,
	}
	if _, err := g.MuatKontrak(ctx, tx, idUjiRevisionDate, doc); err != nil {
		t.Fatalf("memuat: %v", err)
	}

	var terbaca sql.NullString
	q := `SELECT REVISIONDATE FROM ` + skemaRev(t) + `.T_TREATY_REVISION WHERE MASTERID = :1`
	if err := tx.QueryRowContext(ctx, q, idUjiRevisionDate).Scan(&terbaca); err != nil {
		t.Fatalf("membaca REVISIONDATE: %v", err)
	}
	if !terbaca.Valid || terbaca.String != stempel {
		t.Fatalf("REVISIONDATE = %v, mau %q", terbaca, stempel)
	}
}

// ⛔ DAN ia tidak lagi dilaporkan hilang — inilah kalimat yang pemakai baca.
//
// ⚠️ DUA jalur memutuskan laporan itu, dan keduanya harus setuju:
//
//	KunciTakTerpetakan   peta tidak mengenal kuncinya sama sekali
//	KunciBelumTerpasang  peta kenal, tetapi kolomnya belum ada di basis data
//
// Sebelum `451` yang pertama menjeratnya (nol di `Kunci` akar); sesudah peta
// diperbaiki tanpa migrasi, yang kedua akan menjeratnya. Uji ini menutup
// keduanya sekaligus.
func TestRevisionDateTidakLagiDilaporkanHilang(t *testing.T) {
	g, _, ctx := gudangUji(t)

	doc := map[string]any{"ID": idUjiRevisionDate, "RevisionDate": "20261008T083015.000 GMT"}

	for tabel, kunci := range repository.KunciTakTerpetakan(doc) {
		for _, k := range kunci {
			if k == "RevisionDate" {
				t.Fatalf("%s masih menganggap RevisionDate tak terpetakan", tabel)
			}
		}
	}

	belum, err := g.KunciBelumTerpasang(ctx, doc)
	if err != nil {
		t.Fatal(err)
	}
	for _, k := range belum {
		if k == "RevisionDate" {
			t.Fatal("kolom REVISIONDATE belum terpasang — migrasi 451 belum dijalankan")
		}
	}
}
