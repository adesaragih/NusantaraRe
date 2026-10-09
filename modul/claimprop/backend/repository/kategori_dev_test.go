//go:build ujidev

package repository

// Uji DEV migrasi 536 (isi T_KATEGORI_DOC_KLAIM): pernyataan dipecah PEMECAH RUNNER MIGRASI yang sama
// (`migrasi.PernyataanLangkah`), diperiksa `db.PeriksaSQL`, lalu DIURAI (DBMS_SQL.PARSE) - tidak dieksekusi. 08-10-2026
// sebuah baris komentar terpotong membuat 536 gagal di `-migrate`; uji ini menangkap bentuk itu sebelum dijalankan.

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/migrasi"
)

func TestMigrasiKategoriDokumenDiDEV(t *testing.T) {
	cfg, err := config.Load()
	if err != nil || !cfg.PunyaOracle() {
		t.Skip("ORACLE_DSN belum dikonfigurasi")
	}
	d, err := db.Open(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = d.Close() }()
	ctx, batal := context.WithTimeout(context.Background(), time.Minute)
	defer batal()
	p, err := migrasi.PernyataanLangkah(os.DirFS(".."), "536_t_kategori_doc_klaim_isi.sql")
	if err != nil {
		t.Fatal(err)
	}
	if len(p) != 1 {
		t.Fatalf("536: %d pernyataan, mau 1 (INSERT ALL)", len(p))
	}
	q := strings.ReplaceAll(p[0], "{skema}", d.Skema())
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	if !polaDML.MatchString(q) || strings.Count(q, " INTO ") != 49 {
		t.Fatalf("536 bukan INSERT ALL 49 baris: %.120s", q)
	}
	var h hasilDEV
	_, err = d.ExecContext(ctx, sqlUrai, q)
	h.catat(t, "urai 536 INSERT ALL T_KATEGORI_DOC_KLAIM", err)
	if h.gagal > 0 {
		t.Fatal("536 tidak lolos urai Oracle")
	}
}
