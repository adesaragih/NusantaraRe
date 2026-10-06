// Command pindahflat memindah ringkasan R/I Rate Life dari tabel JSON warisan `M_RATE_LIFE_SUMMARY` (cacah dibaca saat
// berjalan; data DEV masih berubah) ke tabel flat `RATE_LIFE_SUMMARY` (migrasi inti 926) - keputusan work owner 06-10-2026 K-F1/K-F2, pola ricommlife.
//
//	go run ./modul/riratelife/backend/alat/pindahflat              # -uji (bawaan): SELECT saja, laporan
//	go run ./modul/riratelife/backend/alat/pindahflat -jalankan    # pemindahan penuh: satu koneksi, satu transaksi
//	go run ./modul/riratelife/backend/alat/pindahflat -jalankan -sejak="<batas delta berikutnya dari laporan sebelumnya>"  # delta
//
// Konfigurasi dari lingkungan yang sama dengan `cmd/api` (`ORACLE_DSN`, `ORACLE_SCHEMA`, `IS_PEGA_PROD`).
// `-jalankan` ditolak bila IS_PEGA_PROD=true (repository.PeriksaMode, SEBELUM koneksi dibuka), bila ada nilai yang tidak
// muat kolom flat (laporan menyebut panjang maksimum tiap kolom; nol pemotongan), dan bila ada kegagalan. Aturan delta
// (baru / berubah / sama / konflik / dilewati): kepala `repository/pindah.go` - tulisan aplikasi yang lebih baru tidak
// pernah ditimpa. Laporan AGREGAT: ID baris dan nama kolom, tidak pernah nilainya.
// M_RATE_LIFE_SUMMARY tidak pernah ditulis.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/riratelife/backend/repository"
)

func main() {
	jalankan := flag.Bool("jalankan", false, "tulis tabel flat dalam satu transaksi (bawaan: -uji, SELECT saja)")
	_ = flag.Bool("uji", true, "uji kering: SELECT saja, nol tulisan (bawaan)")
	sejak := flag.String("sejak", "", "batas delta: MODIFIEDDATE Pega (mis. 20261006T040628.169 GMT) dari laporan putaran sebelumnya; ringkasan lama yang tidak ada di flat tidak dihidupkan lagi")
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

	lap, err := repository.Baru(d).PindahFlat(context.Background(), *jalankan, *sejak)
	fmt.Print(lap.Teks())
	if err != nil {
		fmt.Fprintf(os.Stderr, "pindahflat: %v\n", err)
		os.Exit(1)
	}
	if !lap.BolehDitulis() {
		os.Exit(1)
	}
}
