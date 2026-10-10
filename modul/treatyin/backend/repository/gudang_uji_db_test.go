//go:build db

package repository_test

// Pembantu bersama seluruh uji Oracle modul ini.
//
// ⛔ Satu transaksi per test, diakhiri `Rollback`. Nol `Commit`, nol
// teardown, nol `DELETE` di luar transaksi.
//
// ⚠️ Berkas ini lahir 8 Oktober 2026 ketika berkas yang dulu memuat kedua
// pembantu ini DIBUANG bersama alat pemuat korpus. Pembantunya sendiri tidak
// punya urusan dengan korpus itu — ia hanya membuka Gudang dan transaksi —
// jadi ia dipindahkan ke sini alih-alih ikut terbuang.

import (
	"context"
	"testing"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/repository"
)

// bukaDB membuka Oracle, atau melewati test bila ia tidak terjangkau.
func bukaDB(t *testing.T) (*db.DB, context.Context) {
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
	d, err := db.Open(cfg)
	if err != nil {
		t.Fatalf("membuka oracle: %v", err)
	}
	t.Cleanup(func() { _ = d.Close() })
	ctx := context.Background()
	if err := d.Ping(ctx); err != nil {
		t.Skipf("lewati: oracle tidak terjangkau: %v", err)
	}
	return d, ctx
}

// gudangBaca - Gudang untuk uji yang HANYA membaca; nol transaksi.
func gudangBaca(t *testing.T) (*repository.Gudang, context.Context) {
	t.Helper()
	d, ctx := bukaDB(t)
	return repository.Baru(d), ctx
}

// gudangUji - Gudang beserta satu transaksi yang SELALU dibatalkan.
func gudangUji(t *testing.T) (*repository.Gudang, *db.Tx, context.Context) {
	t.Helper()
	d, ctx := bukaDB(t)
	tx, err := d.Mulai(ctx)
	if err != nil {
		t.Fatalf("memulai transaksi: %v", err)
	}
	// ⛔ Rollback, SELALU.
	t.Cleanup(func() { _ = tx.Rollback() })
	return repository.Baru(d), tx, ctx
}
