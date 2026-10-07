// Command pindahflat memindah rincian R/I Comm Life dari tabel JSON warisan `M_RICOMM_LIFE` ke tabel flat
// `RICOMM_LIFE` (migrasi inti 924) - keputusan work owner 06-10-2026 butir 3, pola masterproductnamelife.
//
//	go run ./modul/ricommlife/backend/alat/pindahflat                        # -uji (bawaan): SELECT saja, laporan
//	go run ./modul/ricommlife/backend/alat/pindahflat -jalankan              # satu transaksi; aman diulang
//	  … -jalankan -terima-normalisasi   # bila laporan menyebut normalisasi angka dan work owner menerimanya
//
// Konfigurasi dari lingkungan yang sama dengan `cmd/api` (`ORACLE_DSN`, `ORACLE_SCHEMA`, `IS_PEGA_PROD`).
// `-jalankan` ditolak bila IS_PEGA_PROD=true (repository.PeriksaMode), bila ada kegagalan, dan bila tabel flat memuat
// ID yang sama dengan isi berbeda. Laporan AGREGAT: ID baris dan nama kolom, tidak pernah nilainya. M_RICOMM_LIFE
// tidak pernah ditulis.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/ricommlife/backend/repository"
)

func main() {
	jalankan := flag.Bool("jalankan", false, "tulis tabel flat dalam satu transaksi (bawaan: -uji, SELECT saja)")
	_ = flag.Bool("uji", true, "uji kering: SELECT saja, nol tulisan (bawaan)")
	terima := flag.Bool("terima-normalisasi", false, "terima normalisasi teks angka (mis. 05 -> 5, 0,5 -> 0.5) yang SUDAH diputuskan work owner")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("konfigurasi: %v", err)
	}
	if !cfg.PunyaOracle() {
		log.Fatal("pindahflat: ORACLE_DSN belum dikonfigurasi")
	}
	if err := repository.PeriksaMode(*jalankan, cfg.IsPegaProd); err != nil {
		log.Fatal(err)
	}
	d, err := db.Open(cfg)
	if err != nil {
		log.Fatalf("oracle: %v", err)
	}
	defer func() { _ = d.Close() }()

	lap, err := repository.Baru(d).PindahFlat(context.Background(), *jalankan, *terima)
	fmt.Print(lap.Teks())
	if err != nil {
		fmt.Fprintf(os.Stderr, "pindahflat: %v\n", err)
		os.Exit(1)
	}
	if !lap.BolehDitulis() {
		os.Exit(1)
	}
}
