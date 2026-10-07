// Command pemuatlama memuat dokumen polis lama GENERASI ENDORSEMEN EDM Treaty In (`POOLDATA.JSON_POLIS` ber-PRODKE
// > 0) ke tabel generasi + proyeksi selisih SUMBER 'PEGA' lewat antarmuka penyimpanan yang sama dengan aplikasi -
// tiket EDM 10 (dan penanda migrasi tiket 09). Asal: pola `modul/nbtreatyin/backend/alat/pemuatlama` (tiket NB 22).
//
//	go run ./modul/edmtreatyin/backend/alat/pemuatlama -keluaran D:\laporan-pemuat-edm            # uji-kering (bawaan)
//	go run ./modul/edmtreatyin/backend/alat/pemuatlama -keluaran D:\laporan-pemuat-edm -jalankan  # tulis
//
// Urutan: generasi NB (PRODKE 0) WAJIB sudah dimuat pemuat NB (`modul/nbtreatyin/backend/alat/pemuatlama
// -jalankan`) - generasi endorsemen pertama menunjuknya lewat OLD_POLIS_ID. Uji-kering hanya MEMBACA JSON_POLIS
// (nol pernyataan ke tabel mana pun) dan menulis kedua berkas laporan. `-jalankan` menulis satu transaksi per
// dokumen, generasi menurut nomornya (NOPOLIS, PRODKE naik); aman diulang (generasi yang sudah dimuat dilewati).
// Konfigurasi dari lingkungan yang sama dengan `cmd/api` (`ORACLE_DSN`, `ORACLE_SCHEMA`, `IS_PEGA_PROD`);
// `-jalankan` ditolak bila IS_PEGA_PROD=true (pola pemuat NB, ADR-U-0005). Dijalankan MANUSIA, tidak pernah agen.
//
// Berkas keluaran per jalankan di folder `-keluaran`:
//
//	edmtreatyin-arsip-medan-<stempel>.csv  POLIS_ID,JALUR,NILAI,KEPUTUSAN - arsip audit pemuatan: setiap medan
//	                                       yang tidak masuk kolom, "dibuang: <alasan>" atau "BELUM DIPUTUSKAN"
//	edmtreatyin-galat-<stempel>.csv        IDPEGA,NOPOLIS,PRODKE,JALUR,NILAI,SEBAB
//
// Keluar 0 hanya bila nol dokumen gagal DAN nol medan BELUM DIPUTUSKAN; medan yang dibuang menurut keputusan
// tertulis tidak menahan.
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
	"nusantarare/modul/edmtreatyin/backend/models"
	"nusantarare/modul/edmtreatyin/backend/services"
)

func main() {
	keluaran := flag.String("keluaran", "", "folder berkas laporan CSV (wajib; dibuat bila belum ada)")
	jalankan := flag.Bool("jalankan", false, "tulis ke tabel, satu transaksi per dokumen (bawaan: uji-kering, baca saja)")
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
		return errors.New("-jalankan ditolak: IS_PEGA_PROD=true (ADR-U-0005; pemuatan produksi oleh WO/DBA)")
	}
	if err := os.MkdirAll(folder, 0o750); err != nil {
		return err
	}
	stempel := time.Now().Format("20060102-150405")
	jalurArsip := filepath.Join(folder, "edmtreatyin-arsip-medan-"+stempel+".csv")
	jalurGalat := filepath.Join(folder, "edmtreatyin-galat-"+stempel+".csv")
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
