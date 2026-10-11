// Command pemuatlama memuat kasus Claim Prop warisan Pega (`OS_AKSEPTASI_KLAIM` baris berlaku per AC 123, plus
// halaman `JSON_KLAIM` bila ada) ke tabel Claim Prop lewat penyimpanan yang sama dengan aplikasi - prompt
// implementasi §6 butir 11, MODUL.md bab "Pemuat data lama".
//
//	go run ./modul/claimprop/backend/alat/pemuatlama -keluaran D:\laporan-pemuat            # uji-kering (bawaan)
//	go run ./modul/claimprop/backend/alat/pemuatlama -keluaran D:\laporan-pemuat -jalankan  # tulis
//
// Uji-kering hanya MEMBACA kedua tabel sumber (nol pernyataan tulis) dan menulis kedua berkas laporan. `-jalankan`
// hanya oleh work owner: satu transaksi per kasus, kasus yang ID-nya sudah ada dilewati (aman diulang), dan ditolak
// bila IS_PEGA_PROD=true (pola `pindahflat`, sama dengan `-migrate`). Konfigurasi dari lingkungan yang sama dengan
// `cmd/api` (`ORACLE_DSN`, `ORACLE_SCHEMA`, `IS_PEGA_PROD`).
//
// Berkas keluaran per jalankan di folder `-keluaran`:
//
//	claimprop-arsip-medan-<stempel>.csv  KASUS,JALUR,NILAI,SEBAB - medan dokumen lama yang tidak masuk kolom
//	claimprop-galat-<stempel>.csv        KASUS,JALUR,NILAI,SEBAB - kasus yang tidak dimuat dan sebabnya
//
// Keluar 0 hanya bila nol kasus gagal.
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
	"nusantarare/modul/claimprop/backend/services"
)

func main() {
	keluaran := flag.String("keluaran", "", "folder berkas laporan CSV (wajib; dibuat bila belum ada)")
	jalankan := flag.Bool("jalankan", false, "tulis ke tabel Claim Prop, satu transaksi per kasus (bawaan: uji-kering, baca saja)")
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

// errBelumSelesai - jalankan tuntas, tetapi ada kasus gagal.
var errBelumSelesai = errors.New("belum selesai: ada kasus gagal (lihat berkas galat)")

func jalan(folder string, tulis bool) (hasil error) {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("konfigurasi: %w", err)
	}
	if !cfg.PunyaOracle() {
		return errors.New("ORACLE_DSN belum dikonfigurasi")
	}
	if tulis && cfg.IsPegaProd {
		return errors.New("-jalankan ditolak: IS_PEGA_PROD=true (pemuatan produksi oleh work owner / DBA, MODUL.md bab Pemuat data lama)")
	}
	if err := os.MkdirAll(folder, 0o750); err != nil {
		return err
	}
	stempel := time.Now().Format("20060102-150405")
	jalurArsip := filepath.Join(folder, "claimprop-arsip-medan-"+stempel+".csv")
	jalurGalat := filepath.Join(folder, "claimprop-galat-"+stempel+".csv")
	fArsip, err := os.Create(jalurArsip)
	if err != nil {
		return err
	}
	defer func() { hasil = errors.Join(hasil, fArsip.Close()) }()
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
	lap := services.LaporanPemuatBaru(fArsip, fGalat, tulis)
	ctx, henti := signal.NotifyContext(context.Background(), os.Interrupt)
	defer henti()
	errJalan := p.Jalankan(ctx, tulis, lap)
	errTutup := lap.Tutup()
	r := lap.Ringkasan()
	fmt.Print(r.Teks())
	fmt.Printf("Berkas arsip medan : %s\nBerkas galat       : %s\n", jalurArsip, jalurGalat)
	if err := errors.Join(errJalan, errTutup); err != nil {
		return err
	}
	if !r.Selesai() {
		return errBelumSelesai
	}
	return nil
}
