// Command pemuatlama - data lama Komite Claim Prop (tiket 13): sensus UJI-KERING kasus komite warisan Pega untuk klaim
// yang ditunda pemuat Claim Prop (baris berlaku OS_AKSEPTASI_KLAIM ber-STS_REJECT 1).
//
//	go run ./modul/komiteclaimprop/backend/alat/pemuatlama -keluaran D:\laporan-pemuat            # uji-kering (bawaan)
//	go run ./modul/komiteclaimprop/backend/alat/pemuatlama -keluaran D:\laporan-pemuat -jalankan  # DITOLAK (OQ)
//
// Uji-kering hanya MEMBACA OS_AKSEPTASI_KLAIM, HISTORYAKSEPTASIPEGA (tanpa USERNAME), dan work object KomiteTreaty di
// DATAPEGA (tanpa BLOB), lalu menulis satu berkas laporan. `-jalankan` hanya oleh work owner, ditolak bila
// IS_PEGA_PROD=true, dan saat ini ditolak di mana pun: tangga lama belum punya sumber yang dapat dimuat tanpa mengarang
// (MODUL.md, docs/PARITAS.md). Konfigurasi dari lingkungan yang sama dengan `cmd/api` (`ORACLE_DSN`, `ORACLE_SCHEMA`,
// `IS_PEGA_PROD`).
//
// Berkas keluaran di folder `-keluaran`:
//
//	komiteclaimprop-sensus-<stempel>.csv  KLAIM,BARIS_OS,RIWAYAT_AKSEPTASI,KOMITE_DI_RIWAYAT,WORK_OBJECT_KOMITE,SEBAB
package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"time"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/komiteclaimprop/backend/repository"
	"nusantarare/modul/komiteclaimprop/backend/services"
)

func main() {
	keluaran := flag.String("keluaran", "", "folder berkas laporan CSV (wajib; dibuat bila belum ada)")
	jalankan := flag.Bool("jalankan", false, "tulis ke tabel komite (DITOLAK sampai sumber tangga lama diputuskan)")
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

func jalan(folder string, tulis bool) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("konfigurasi: %w", err)
	}
	if tulis {
		return services.JalankanLama(cfg.IsPegaProd)
	}
	if !cfg.PunyaOracle() {
		return fmt.Errorf("ORACLE_DSN belum dikonfigurasi")
	}
	d, err := db.Open(cfg)
	if err != nil {
		return fmt.Errorf("membuka basis data: %w", err)
	}
	defer func() { _ = d.Close() }()
	ctx, batal := signal.NotifyContext(context.Background(), os.Interrupt)
	defer batal()
	s, pegaGagal, err := services.SensusLama(ctx, repository.Baru(d))
	if err != nil {
		return err
	}
	if err := os.MkdirAll(folder, 0o755); err != nil {
		return err
	}
	nama := filepath.Join(folder, "komiteclaimprop-sensus-"+time.Now().Format("20060102-150405")+".csv")
	f, err := os.Create(nama)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	_ = w.Write([]string{"KLAIM", "BARIS_OS", "RIWAYAT_AKSEPTASI", "KOMITE_DI_RIWAYAT", "WORK_OBJECT_KOMITE", "SEBAB"})
	for _, b := range s.Baris {
		_ = w.Write([]string{b.Klaim, strconv.Itoa(b.BarisOS), strconv.Itoa(b.Riwayat), strconv.Itoa(b.KomiteRiwayat),
			strconv.Itoa(b.KomitePega), b.Sebab})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	fmt.Printf("klaim CLMP lama: %d; ditunda (STS_REJECT 1): %d; ber-BLOB KomiteTreaty: %d; riwayat tanpa tangga: %d; "+
		"tanpa jejak komite: %d; dimuat: 0 (uji-kering)\nlaporan: %s\n", s.KlaimLama, s.Ditunda, s.DenganBlob,
		s.TanpaTangga, s.TanpaKomite, nama)
	if pegaGagal != nil {
		fmt.Printf("peringatan: DATAPEGA tidak terbaca - work object KomiteTreaty dihitung 0 (%v)\n", pegaGagal)
	}
	return nil
}
