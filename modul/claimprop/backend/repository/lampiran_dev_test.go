//go:build ujidev

package repository

// Uji DEV lampiran klaim (pola sqldev_test.go): INSERT dokumen klaim hanya DIURAI (DBMS_SQL.PARSE), SELECT dijalankan
// dengan masukan UJI-*. T_KATEGORI_DOC_KLAIM baru ada sesudah migrasi 535/536 (`-migrate` work owner).

import (
	"context"
	"testing"
	"time"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
)

func TestSQLLampiranDiDEV(t *testing.T) {
	cfg, err := config.Load()
	if err != nil || !cfg.PunyaOracle() {
		t.Skip("ORACLE_DSN belum dikonfigurasi")
	}
	d, err := db.Open(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = d.Close() }()
	ctx, batal := context.WithTimeout(context.Background(), 2*time.Minute)
	defer batal()
	g := Baru(d)
	var h hasilDEV
	tabel, err := d.Qualify(TabelDokumenKlaim)
	if err != nil {
		t.Fatal(err)
	}
	q := sqlSisipDokumenKlaim(tabel)
	if !polaDML.MatchString(q) {
		t.Fatal("bukan DML - DBMS_SQL.PARSE mengeksekusi DDL, ditolak")
	}
	_, err = d.ExecContext(ctx, sqlUrai, q)
	h.catat(t, "urai SisipDokumenKlaim (tabel warisan)", err)
	_, err = d.ExecContext(ctx, sqlUrai, sqlHapusDokumenKlaim(tabel))
	h.catat(t, "urai HapusDokumenKlaim (tabel warisan)", err)
	_, err = d.ExecContext(ctx, sqlUrai, sqlPindahKategoriDokumen(tabel))
	h.catat(t, "urai PindahKategoriDokumen (tabel warisan)", err)
	_, err = g.DaftarLampiran(ctx, "UJI-CLMP-0")
	h.catat(t, "baca DaftarLampiran", err)
	kat, err := g.KategoriLampiran(ctx, "UJI-CLMP-0")
	h.catat(t, "baca KategoriLampiran (T_KATEGORI_DOC_KLAIM, migrasi 535)", err)
	// Master PROP terisi (536, 15 baris): kategori klaim mana pun = seluruh master, cacah 0. Nol baris = bind salah urut
	// (Oracle mengikat bind SQL menurut urutan kemunculan, bukan nomor `:n`).
	if err == nil && len(kat) == 0 {
		t.Errorf("GAGAL  KategoriLampiran: nol kategori PROP padahal master terisi")
	}
	t.Logf("kategori PROP terbaca: %d", len(kat))
	t.Logf("ringkasan: ok %d, belum-migrasi %d, gagal %d", h.ok, h.belum, h.gagal)
}
