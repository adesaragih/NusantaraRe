// Command pindahclient memindah SEKALI isi dokumen organisasi Pega (`M_CLIENT`) ke tabel datar Company Detail -
// disetujui work owner 03-10-2026. Rinciannya: `repository/pindah.go`.
//
//	go run ./modul/companydetail/backend/alat/pindahclient             # -uji (bawaan): SELECT saja, laporan
//	go run ./modul/companydetail/backend/alat/pindahclient -jalankan   # satu transaksi; aman diulang
//
// Butuh migrasi modul 800-807 sudah dijalankan (kolom tujuan, M_ENUMERASI, NATION). Konfigurasi dari lingkungan yang
// sama dengan `cmd/api` (`ORACLE_DSN`, `ORACLE_SCHEMA`, `IS_PEGA_PROD`). `-jalankan` ditolak bila IS_PEGA_PROD=true.
// Laporan AGREGAT: jumlah per jenis, tidak pernah ID kasus, nama, atau nilai. Dokumen M_CLIENT tidak pernah ditulis.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/companydetail/backend/repository"
)

func main() {
	jalankan := flag.Bool("jalankan", false, "tulis tabel datar dalam satu transaksi (bawaan: -uji, SELECT saja)")
	_ = flag.Bool("uji", true, "uji kering: SELECT saja, nol tulisan (bawaan)")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("konfigurasi: %v", err)
	}
	if !cfg.PunyaOracle() {
		log.Fatal("pindahclient: ORACLE_DSN belum dikonfigurasi")
	}
	if *jalankan && cfg.IsPegaProd {
		log.Fatal(repository.ErrPindahDiProduksi)
	}
	d, err := db.Open(cfg)
	if err != nil {
		log.Fatalf("oracle: %v", err)
	}
	defer func() { _ = d.Close() }()

	lap, err := repository.Baru(d).PindahClient(context.Background(), *jalankan)
	fmt.Print(lap.Teks())
	if err != nil {
		fmt.Fprintf(os.Stderr, "pindahclient: %v\n", err)
		os.Exit(1)
	}
}
