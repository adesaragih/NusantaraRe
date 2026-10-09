//go:build ignore

// Mengosongkan SELURUH tabel pendaratan Treaty In - tanpa saringan kontrak.
//
//	go run modul/treatyin/backend/pemuat/kosongkan_semua.go          (kering)
//	go run modul/treatyin/backend/pemuat/kosongkan_semua.go -ikat    (mengikat)
//
// ---------------------------------------------------------------------
// ⛔ MENGAPA ADA DI SAMPING `-kosongkan` MILIK PEMUAT
// ---------------------------------------------------------------------
// `jalankan.go -kosongkan` menghapus PER KONTRAK: ia mendaftar `MASTERID`
// dari korpus dokumen lalu menjalankan 30 `DELETE … WHERE MASTERID = :1`
// untuk masing-masing. Pada 2.415 kontrak itu 72.450 pernyataan, dan pada
// 7 Oktober 2026 ia belum selesai sesudah lima menit.
//
// ⚠️ Dan ia TIDAK LENGKAP: baris yang `MASTERID`-nya sudah tidak ada di
// korpus dokumen - kontrak yang dokumennya dibuang, atau baris yang
// ditulis aplikasi sendiri - nol tersentuh. Untuk "kosongkan semuanya",
// saringan itu justru lubangnya.
//
// Berkas ini menjalankan 30 `DELETE` TANPA `WHERE`. Tiga puluh pernyataan,
// bukan 72.450.
//
// ---------------------------------------------------------------------
// ⛔ PAGAR
// ---------------------------------------------------------------------
//  1. Daftar tabelnya DARI `repository.PetaPendaratan`, urutan terbalik
//     (anak sebelum induk). Nol nama tabel diketik di berkas ini - jadi nol
//     peluang salah ketik mengenai tabel warisan.
//  2. `IS_PEGA_PROD` menolak jalan, sama seperti pemuat.
//  3. Bawaannya KERING: tanpa `-ikat` seluruhnya dibatalkan dan yang
//     tercetak adalah apa yang AKAN terjadi.
//  4. `DELETE`, bukan `TRUNCATE` - agar `ROLLBACK` masih berarti.
//
// ⛔ `TREATY_IN`, `M_TREATY_IN`, `TREATYEXCHANGEYEARLY`, `REINSURANCETYPE`,
// `ACHIEVEMENT`, `PROPORTIONALARRG` NOL disentuh - nol di antaranya ada di
// `PetaPendaratan`, dan `TestKosongkanNolMenyentuhTabelLuar` menjaganya.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/repository"
)

func main() {
	ikat := flag.Bool("ikat", false, "COMMIT; tanpa ini seluruhnya dibatalkan")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("konfigurasi: %v", err)
	}
	if cfg.IsPegaProd {
		log.Fatal("menolak berjalan saat IS_PEGA_PROD=true (ADR-U-0005)")
	}
	if !cfg.PunyaOracle() {
		log.Fatal("ORACLE_DSN kosong")
	}
	d, err := db.Open(cfg)
	if err != nil {
		log.Fatalf("membuka oracle: %v", err)
	}
	defer func() { _ = d.Close() }()

	ctx := context.Background()
	if err := d.Ping(ctx); err != nil {
		log.Fatalf("oracle tidak terjangkau: %v", err)
	}
	fmt.Printf("skema %s · ikat=%v · %d tabel\n", cfg.OracleSchema, *ikat, len(repository.PetaPendaratan))

	tx, err := d.Mulai(ctx)
	if err != nil {
		log.Fatalf("membuka transaksi: %v", err)
	}

	var total int64
	// Anak sebelum induk — kebalikan urutan muat.
	for i := len(repository.PetaPendaratan) - 1; i >= 0; i-- {
		p := repository.PetaPendaratan[i]
		nama, err := d.Qualify(p.Tabel)
		if err != nil {
			_ = tx.Rollback()
			log.Fatalf("qualify %s: %v", p.Tabel, err)
		}
		q := fmt.Sprintf("DELETE FROM %s", nama)
		if err := db.PeriksaSQL(q); err != nil {
			_ = tx.Rollback()
			log.Fatalf("periksa SQL %s: %v", p.Tabel, err)
		}
		hasil, err := tx.ExecContext(ctx, q)
		if err != nil {
			_ = tx.Rollback()
			log.Fatalf("mengosongkan %s: %v", p.Tabel, err)
		}
		n, _ := hasil.RowsAffected()
		total += n
		fmt.Printf("  %-32s %8d\n", p.Tabel, n)
	}

	if *ikat {
		if err := tx.Commit(); err != nil {
			log.Fatalf("mengikat: %v", err)
		}
		fmt.Printf("TOTAL %d baris DIBUANG dan DIIKAT\n", total)
		return
	}
	if err := tx.Rollback(); err != nil {
		log.Fatalf("membatalkan: %v", err)
	}
	fmt.Printf("TOTAL %d baris AKAN dibuang — DIBATALKAN (pakai -ikat)\n", total)
}
