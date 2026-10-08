//go:build db

package repository_test

// Harness uji Oracle PERTAMA modul `treatyinadjustment`.
//
// ⛔ POLANYA DISALIN dari `modul/treatyin/backend/repository`, bukan dikarang
// ulang — dan tiga aturannya dibawa apa adanya:
//
//	1. SATU TRANSAKSI PER TEST, diakhiri `Rollback` lewat `t.Cleanup`.
//	   Nol `Commit`, di cabang mana pun.
//	2. NOL `uji/skemauji`. `Bongkar`-nya menjalankan
//	   `DROP TABLE … CASCADE CONSTRAINTS` atas tabel WARISAN sungguhan, dan
//	   `ORACLE_SCHEMA` di sini POOLDATA — skema aplikasi, bukan skema buangan.
//	3. NOL pemasangan dan NOL pembongkaran. Yang dibaca tabel yang sudah
//	   hidup; yang ditulis (bila kelak ada) dibatalkan.
//
// ⚠️ Produksi ditolak keras. Itu satu-satunya pagar yang tidak pernah
// dicabut dari jalur mana pun di repositori ini.

import (
	"context"
	"database/sql"
	"testing"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyinadjustment/backend/repository"
)

// siapkan membuka Gudang sungguhan beserta satu transaksi yang SELALU
// dibatalkan.
//
// Mengembalikan `*sql.Tx` telanjang di samping `*db.Tx`: yang pertama untuk
// pernyataan mentah di dalam uji, yang kedua untuk memanggil repository.
// Keduanya transaksi YANG SAMA, jadi satu `Rollback` menutup keduanya.
func siapkan(t *testing.T) (*repository.Gudang, *db.Tx, context.Context) {
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
	tx, err := d.Mulai(ctx)
	if err != nil {
		t.Fatalf("memulai transaksi: %v", err)
	}
	// ⛔ Rollback, SELALU — tidak ada cabang yang melakukan Commit.
	t.Cleanup(func() { _ = tx.Rollback() })
	return repository.Baru(d), tx, ctx
}

// bacaSaja membuka Gudang TANPA transaksi — untuk uji yang hanya membaca.
//
// Dipisah dengan sengaja: transaksi yang dibuka lalu tidak dipakai menahan
// sumber daya tanpa membeli apa pun, dan transaksi yang dibuka membuat uji
// baca terlihat seolah ia menulis.
func bacaSaja(t *testing.T) (*repository.Gudang, context.Context) {
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
	return repository.Baru(d), ctx
}

// sqlMentah membuka koneksi kedua untuk pernyataan pembanding.
//
// ⚠️ Dipakai HANYA untuk membaca angka pembanding dari Oracle — cacah baris
// yang diadu dengan hasil repository. Ia tidak pernah menulis.
func sqlMentah(t *testing.T) (*sql.DB, string) {
	t.Helper()
	cfg, _ := config.Load()
	h, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = h.Close() })
	return h, cfg.OracleSchema
}

// skemaUji mengembalikan nama skema dari konfigurasi — dipakai pernyataan
// mentah yang memuat penanda `{skema}`.
func skemaUji(t *testing.T) string {
	t.Helper()
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("konfigurasi: %v", err)
	}
	return cfg.OracleSchema
}
