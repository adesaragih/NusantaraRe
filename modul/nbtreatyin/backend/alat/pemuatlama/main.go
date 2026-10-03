// Command pemuatlama memuat dokumen polis lama NB Treaty In (`POOLDATA.JSON_POLIS`,
// generasi NB = PRODKE 0) ke 8 tabel diagram grilling lewat antarmuka penyimpanan
// yang sama dengan aplikasi - tiket 22, MODUL.md "Pemuat dokumen lama".
//
//	go run ./modul/nbtreatyin/backend/alat/pemuatlama -keluaran D:\laporan-pemuat            # uji-kering (bawaan)
//	go run ./modul/nbtreatyin/backend/alat/pemuatlama -keluaran D:\laporan-pemuat -jalankan  # tulis
//
// Uji-kering hanya MEMBACA JSON_POLIS (nol pernyataan ke tabel baru) dan menulis
// kedua berkas laporan. `-jalankan` menulis satu transaksi per dokumen; aman
// diulang (kasus yang sudah dimuat dilewati). Konfigurasi dari lingkungan yang
// sama dengan `cmd/api` (`ORACLE_DSN`, `ORACLE_SCHEMA`, `IS_PEGA_PROD`);
// `-jalankan` ditolak bila IS_PEGA_PROD=true (pola `pindahflat`, ADR-U-0005).
//
// Berkas keluaran per jalankan di folder `-keluaran`:
//
//	nbtreatyin-medan-tak-dikenal-<stempel>.csv  POLIS_ID,JALUR,NILAI (K17) - WAJIB 0 baris data
//	nbtreatyin-galat-<stempel>.csv              IDPEGA,NOPOLIS,JALUR,NILAI,SEBAB (AC 58, K15)
//
// Keluar 0 hanya bila nol dokumen gagal DAN nol medan tak dikenal (AC 59).
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/services"
)

func main() {
	keluaran := flag.String("keluaran", "", "folder berkas laporan CSV (wajib; dibuat bila belum ada)")
	jalankan := flag.Bool("jalankan", false, "tulis ke tabel baru, satu transaksi per dokumen (bawaan: uji-kering, baca saja)")
	flag.Parse()
	if *keluaran == "" {
		fmt.Fprintln(os.Stderr, "pemuatlama: -keluaran <folder> wajib diisi")
		flag.Usage()
		os.Exit(2)
	}
	if err := jalan(*keluaran, *jalankan); err != nil {
		fmt.Fprintf(os.Stderr, "pemuatlama: %v\n", err)
		os.Exit(1)
	}
}

// errBelumSelesai - jalankan tuntas, tetapi masih ada galat atau medan tak dikenal.
var errBelumSelesai = errors.New("belum selesai: dokumen gagal atau medan tak dikenal > 0 (lihat berkas laporan)")

func jalan(folder string, tulis bool) (hasil error) {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("konfigurasi: %w", err)
	}
	if !cfg.PunyaOracle() {
		return errors.New("ORACLE_DSN belum dikonfigurasi")
	}
	if tulis && cfg.IsPegaProd {
		return errors.New("-jalankan ditolak: IS_PEGA_PROD=true (pemindahan di produksi butuh keputusan work owner)")
	}
	if err := os.MkdirAll(folder, 0o750); err != nil {
		return err
	}
	stempel := time.Now().Format("20060102-150405")
	jalurTak := filepath.Join(folder, "nbtreatyin-medan-tak-dikenal-"+stempel+".csv")
	jalurGalat := filepath.Join(folder, "nbtreatyin-galat-"+stempel+".csv")
	fTak, err := os.Create(jalurTak)
	if err != nil {
		return err
	}
	defer func() { hasil = errors.Join(hasil, fTak.Close()) }()
	fGalat, err := os.Create(jalurGalat)
	if err != nil {
		return err
	}
	defer func() { hasil = errors.Join(hasil, fGalat.Close()) }()

	d, err := db.Open(cfg)
	if err != nil {
		return fmt.Errorf("oracle: %w", err)
	}
	defer func() { _ = d.Close() }()
	p, err := services.PemuatDariDasar(inti.NewDasar(d))
	if err != nil {
		return err
	}
	lap, err := models.LaporanPemuatBaru(fTak, fGalat, tulis)
	if err != nil {
		return err
	}
	ctx, henti := signal.NotifyContext(context.Background(), os.Interrupt)
	defer henti()
	errJalan := p.Jalankan(ctx, tulis, lap)
	errTutup := lap.Tutup()
	r := lap.Ringkasan()
	fmt.Print(r.Teks())
	fmt.Printf("Berkas medan tak dikenal : %s\nBerkas galat             : %s\n", jalurTak, jalurGalat)
	if err := errors.Join(errJalan, errTutup); err != nil {
		return err
	}
	if !r.Selesai() {
		return errBelumSelesai
	}
	return nil
}
