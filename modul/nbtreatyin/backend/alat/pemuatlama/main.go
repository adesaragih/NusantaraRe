// Command pemuatlama memuat dokumen polis lama NB Treaty In (`POOLDATA.JSON_POLIS`,
// generasi NB = PRODKE 0) ke 8 tabel diagram grilling lewat antarmuka penyimpanan
// yang sama dengan aplikasi - tiket 22, MODUL.md "Pemuat dokumen lama".
//
//	go run ./modul/nbtreatyin/backend/alat/pemuatlama -keluaran D:\laporan-pemuat            # uji-kering (bawaan)
//	go run ./modul/nbtreatyin/backend/alat/pemuatlama -keluaran D:\laporan-pemuat -jalankan  # tulis
//
// Uji-kering hanya MEMBACA JSON_POLIS (nol pernyataan ke tabel mana pun) dan menulis
// kedua berkas laporan. `-jalankan` menulis satu transaksi per dokumen; aman
// diulang (kasus yang sudah dimuat dilewati). Konfigurasi dari lingkungan yang
// sama dengan `cmd/api` (`ORACLE_DSN`, `ORACLE_SCHEMA`, `IS_PEGA_PROD`);
// `-jalankan` ditolak bila IS_PEGA_PROD=true (pola `pindahflat`, ADR-U-0005, sama
// dengan `-migrate`). Urutan resmi pemuatan (F7, WO 04-10-2026): migrasi 320-327
// oleh WO -> uji-kering di skema uji (K11) -> F3 tuntas (nol medan BELUM
// DIPUTUSKAN) -> pemuatan produksi oleh WO/DBA; tidak pernah dijalankan agen.
//
// Berkas keluaran per jalankan di folder `-keluaran`:
//
//	nbtreatyin-arsip-medan-<stempel>.csv  POLIS_ID,JALUR,NILAI,KEPUTUSAN - arsip audit pemuatan (F3):
//	                                      setiap medan yang tidak masuk kolom, "dibuang: <alasan>"
//	                                      atau "BELUM DIPUTUSKAN"
//	nbtreatyin-galat-<stempel>.csv        IDPEGA,NOPOLIS,JALUR,NILAI,SEBAB (AC 58, K15)
//
// Keluar 0 hanya bila nol dokumen gagal DAN nol medan BELUM DIPUTUSKAN (AC 59
// RALAT F3); medan yang dibuang menurut keputusan tertulis tidak menahan.
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

// errBelumSelesai - jalankan tuntas, tetapi masih ada galat atau medan belum diputuskan.
var errBelumSelesai = errors.New("belum selesai: dokumen gagal atau medan BELUM DIPUTUSKAN > 0 (lihat berkas laporan)")

func jalan(folder string, tulis bool) (hasil error) {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("konfigurasi: %w", err)
	}
	if !cfg.PunyaOracle() {
		return errors.New("ORACLE_DSN belum dikonfigurasi")
	}
	if tulis && cfg.IsPegaProd {
		return errors.New("-jalankan ditolak: IS_PEGA_PROD=true (ADR-U-0005; pemuatan produksi oleh WO/DBA menurut urutan F7, MODUL.md bab Migrasi)")
	}
	if err := os.MkdirAll(folder, 0o750); err != nil {
		return err
	}
	stempel := time.Now().Format("20060102-150405")
	jalurArsip := filepath.Join(folder, "nbtreatyin-arsip-medan-"+stempel+".csv")
	jalurGalat := filepath.Join(folder, "nbtreatyin-galat-"+stempel+".csv")
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
	lap, err := models.LaporanPemuatBaru(fArsip, fGalat, tulis)
	if err != nil {
		return err
	}
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
