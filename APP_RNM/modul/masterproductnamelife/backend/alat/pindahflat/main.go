// Command pindahflat memindah produk Master Product Name Life dari kedua tabel JSON warisan (`M_PRODUCT_LIFE`,
// `M_PRODUCTINWARD_LIFE`) ke tabel flat `M_PRODUCTNAME_LIFE` + tujuh anak - tiket 01 bab 02-10-2026 (T2/T3),
// panduan `docs/PANDUAN-PINDAH-FLAT.md`.
//
//	go run ./modul/masterproductnamelife/backend/alat/pindahflat            # -uji (bawaan): SELECT saja, laporan
//	go run ./modul/masterproductnamelife/backend/alat/pindahflat -jalankan  # satu transaksi; aman diulang
//	  … -jalankan -terima-normalisasi="nol depan,koma desimal"  # hanya jenis yang diputuskan work owner (OQ-FLAT-07)
//
// Konfigurasi dari lingkungan yang sama dengan `cmd/api` (`ORACLE_DSN`, `ORACLE_SCHEMA`, `IS_PEGA_PROD`).
// `-jalankan` ditolak bila IS_PEGA_PROD=true, bila rekonsiliasi tidak lolos, dan bila tabel flat sudah memuat produk
// yang berbeda dari sumber JSON. Laporan AGREGAT: ID produk dan nama kolom, tidak pernah nilainya. Kedua tabel JSON
// tidak pernah ditulis.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterproductnamelife/backend/repository"
)

func main() {
	jalankan := flag.Bool("jalankan", false, "tulis tabel flat dalam satu transaksi (bawaan: -uji, SELECT saja)")
	_ = flag.Bool("uji", true, "uji kering: SELECT saja, nol tulisan (bawaan)")
	terima := flag.String("terima-normalisasi", "",
		"jenis normalisasi (dipisah koma, mis. \"nol depan,koma desimal\") yang SUDAH diputuskan work owner (OQ-FLAT-07)")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("konfigurasi: %v", err)
	}
	if !cfg.PunyaOracle() {
		log.Fatal("pindahflat: ORACLE_DSN belum dikonfigurasi")
	}
	if *jalankan && cfg.IsPegaProd {
		log.Fatal(repository.ErrPindahDiProduksi)
	}
	d, err := db.Open(cfg)
	if err != nil {
		log.Fatalf("oracle: %v", err)
	}
	defer func() { _ = d.Close() }()

	var jenis []string
	for _, j := range strings.Split(*terima, ",") {
		if j = strings.TrimSpace(j); j != "" {
			jenis = append(jenis, j)
		}
	}
	lap, err := repository.Baru(d).PindahFlat(context.Background(), *jalankan, jenis)
	fmt.Print(lap.Teks())
	if err != nil {
		fmt.Fprintf(os.Stderr, "pindahflat: %v\n", err)
		os.Exit(1)
	}
	if !lap.BolehDitulis() {
		os.Exit(1)
	}
}
