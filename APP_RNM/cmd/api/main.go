// Command api adalah entry point backend Nusantara Re.
//
// Fase 0 - scaffold. Berkas ini hanya: baca config -> buka koneksi ->
// daftarkan handler -> dengarkan. ⛔ Nol aturan dagang di sini.
//
// main adalah composition root: ia satu-satunya tempat yang merakit
// repository dan services. Sesudah dirakit, seluruh pertanyaan lewat services.
package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"nusantarare/internal/config"
	"nusantarare/internal/handlers"
	"nusantarare/internal/repository"
	"nusantarare/internal/services"
)

func main() {
	migrasi := flag.Bool("migrate", false, "jalankan migrasi lalu keluar")
	bongkar := flag.Bool("migrate-down", false,
		"BONGKAR skema uji lalu keluar - MENGHAPUS tabel; perlu ORACLE_SKEMA_UJI=true")
	// Tiket 01 Treaty Contract Out (tco2). Ditambahkan ADITIF 28-09-2026.
	pindahTCO := flag.Bool("migrate-data-treaty-contract-out", false,
		"pindahkan enam tabel warisan master arrangement Treaty Contract Out ke tabel T_* "+
			"(satu transaksi, rekonsiliasi tepat, sequence diselaraskan) lalu keluar")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("konfigurasi: %v", err)
	}

	var db *repository.DB
	if cfg.PunyaOracle() {
		db, err = repository.Open(cfg)
		if err != nil {
			log.Fatalf("oracle: %v", err)
		}
		defer func() { _ = db.Close() }()
	}

	// ⛔ Lingkungan efek keluar dipasang SEKALI di sini, dari `config.IsPegaProd`
	// (ADR-U-0005). Tanpa baris ini tidak ada efek keluar yang pernah berjalan
	// di mana pun - dan AC lingkungan akan tercentang secara hampa.
	//
	// ⚠️ Ia menggerbangi EFEK KELUAR saja, tidak pernah penyimpanan: klaim
	// tetap tersimpan di lingkungan non-produksi.
	svc := services.New(db).
		DenganLingkungan(services.LingkunganDariFlag(cfg.IsPegaProd)).
		DenganUnggahanDir(cfg.UnggahanDir).
		// OQ-TCO-08: bawaan stub; ⛔ garam tidak pernah dicetak.
		DenganPenyimpananLampiranTCO(cfg.PelaksanaStorage == config.PelaksanaStorageNyata, cfg.StorageTokenSalt)
	if svc.PunyaDatabase() {
		log.Printf("oracle: skema %s", svc.SkemaAktif())
	} else {
		log.Print("oracle: ORACLE_DSN kosong - berjalan tanpa database")
	}

	if (*migrasi && *bongkar) || (*pindahTCO && (*migrasi || *bongkar)) {
		log.Fatal("pilih salah satu: -migrate, -migrate-down, atau -migrate-data-treaty-contract-out")
	}
	if *bongkar {
		bongkarMigrasi(svc, cfg)
		return
	}
	if *migrasi {
		jalankanMigrasi(svc)
		return
	}
	if *pindahTCO {
		pindahkanDataTreatyContractOut(svc, cfg)
		return
	}

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handlers.Router(svc, cfg.AuthStub),
		ReadHeaderTimeout: 10 * time.Second,
	}

	ctx, berhenti := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer berhenti()

	log.Printf("lampiran treaty contract out: pelaksana penyimpanan %s", cfg.PelaksanaStorage)
	pekerjaSelesai := jalankanPekerjaLampiranTCO(ctx, svc, cfg)

	go func() {
		log.Printf("http: mendengarkan di %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http: %v", err)
		}
	}()

	<-ctx.Done()
	log.Print("http: menutup")

	tutup, batal := context.WithTimeout(context.Background(), 10*time.Second)
	defer batal()
	if err := srv.Shutdown(tutup); err != nil {
		log.Printf("http: penutupan tidak bersih: %v", err)
	}
	// ⛔ Ditunggu SEBELUM db ditutup (defer di atas): putaran yang sedang
	// berjalan menuntaskan atau membatalkan transaksinya sendiri.
	select {
	case <-pekerjaSelesai:
	case <-tutup.Done():
		log.Print("lampiran treaty contract out: pekerja latar belum berhenti saat batas penutupan")
	}
}

// jalankanPekerjaLampiranTCO menyalakan pekerja latar antrean lampiran Treaty
// Contract Out (OQ-TCO-09, keputusan work owner 29-09-2026).
//
// Mati bila TCO_PEKERJA_LAMPIRAN_INTERVAL kosong/0 atau tanpa Oracle. Ia
// berhenti bersama ctx proses; kanal yang dikembalikan tertutup saat ia
// benar-benar berhenti (langsung tertutup bila tidak dinyalakan).
func jalankanPekerjaLampiranTCO(ctx context.Context, svc *services.Service, cfg config.Config) <-chan struct{} {
	selesai := make(chan struct{})
	if cfg.IntervalPekerjaLampiranTCO <= 0 {
		log.Print("lampiran treaty contract out: pekerja latar mati (interval kosong)")
		close(selesai)
		return selesai
	}
	if !svc.PunyaDatabase() {
		log.Print("lampiran treaty contract out: pekerja latar mati (tanpa oracle)")
		close(selesai)
		return selesai
	}
	log.Printf("lampiran treaty contract out: pekerja latar tiap %s", cfg.IntervalPekerjaLampiranTCO)
	go func() {
		defer close(selesai)
		handlers.LayananLampiranTCO(svc).JalankanPekerja(ctx, cfg.IntervalPekerjaLampiranTCO,
			func(s string) { log.Print(s) })
	}()
	return selesai
}

// bongkarMigrasi adalah titik masuk `-migrate-down`.
//
// ⛔ Ia MENGHAPUS tabel, dan karena itu dipagari sama persis dengan test bertag
// db: menolak IS_PEGA_PROD=true, menolak tanpa ORACLE_SKEMA_UJI=true, dan
// menolak skema yang memuat POOLDATA. Pagarnya satu-satunya, tinggal di
// internal/config, supaya jalur ini dan jalur test tidak mungkin berselisih.
//
// Kenapa flag ini ada: sampai 26-09-2026 jalur mundur hanya punya pemanggil
// test. Orang yang ingin membongkar skema uji terpaksa menyalin isi berkas
// *_down.sql ke sqlplus - dan itu MELEWATI pengaman T_MIGRASI, yang membongkar
// hanya langkah yang benar-benar tercatat selesai.
func bongkarMigrasi(svc *services.Service, cfg config.Config) {
	if !svc.PunyaDatabase() {
		log.Fatal("bongkar: ORACLE_DSN wajib terisi")
	}
	if cfg.IsPegaProd {
		log.Fatal("bongkar: menolak berjalan saat IS_PEGA_PROD=true (ADR-U-0005)")
	}
	if err := cfg.PastikanSkemaUji(); err != nil {
		log.Fatalf("bongkar: %v", err)
	}
	ctx, batal := context.WithTimeout(context.Background(), 30*time.Second)
	defer batal()
	if err := svc.CekKesehatan(ctx); err != nil {
		log.Fatalf("bongkar: tidak dapat menjangkau oracle: %v", err)
	}
	lap, err := svc.BongkarMigrasi(ctx)
	if err != nil {
		log.Fatalf("bongkar: %v", err)
	}
	log.Printf("bongkar skema %s: %d langkah dibongkar, %d dilewati",
		svc.SkemaAktif(), len(lap.Dijalankan), len(lap.Dilewati))
	for _, n := range lap.Dijalankan {
		log.Printf("  dibongkar: %s", n)
	}
}

// jalankanMigrasi adalah titik masuk `make migrate`.
//
// Fase 0 tidak punya migrasi: nol DDL, nol aturan dagang. Berkas migrasi lahir
// bersama tiket yang memilikinya, dan ⛔ tidak pernah dijalankan terhadap
// instance produksi.
func jalankanMigrasi(svc *services.Service) {
	if !svc.PunyaDatabase() {
		log.Fatal("migrasi: ORACLE_DSN wajib terisi")
	}
	ctx, batal := context.WithTimeout(context.Background(), 30*time.Second)
	defer batal()
	if err := svc.CekKesehatan(ctx); err != nil {
		log.Fatalf("migrasi: tidak dapat menjangkau oracle: %v", err)
	}
	lap, err := svc.JalankanMigrasi(ctx)
	if err != nil {
		log.Fatalf("migrasi: %v", err)
	}
	log.Printf("migrasi skema %s: %d langkah dijalankan, %d dilewati, %d pernyataan",
		svc.SkemaAktif(), len(lap.Dijalankan), len(lap.Dilewati), lap.Pernyataan)
	for _, n := range lap.Dijalankan {
		log.Printf("  dijalankan: %s", n)
	}
	// Daftar ini biasanya kosong. Kalau berisi, ada langkah yang dulu gagal di
	// tengah jalan: sebagian objeknya sudah berdiri sebelum percobaan ini.
	for _, n := range lap.ObjekSudahAda {
		log.Printf("  dilewati, objeknya sudah ada: %s", n)
	}
}

// pindahkanDataTreatyContractOut adalah titik masuk
// `-migrate-data-treaty-contract-out` (tiket 01 Treaty Contract Out, tco2).
//
// Ia MENULIS ke tabel T_TREATY* dari enam tabel warisan yang hanya dibaca.
// Menolak IS_PEGA_PROD=true (ADR-U-0005). Laporannya dicetak SELALU - juga
// saat dibatalkan - sebab temuan dan selisihnya adalah alasan pembatalan.
func pindahkanDataTreatyContractOut(svc *services.Service, cfg config.Config) {
	if !svc.PunyaDatabase() {
		log.Fatal("migrasi data treaty contract out: ORACLE_DSN wajib terisi")
	}
	if cfg.IsPegaProd {
		log.Fatal("migrasi data treaty contract out: menolak berjalan saat IS_PEGA_PROD=true (ADR-U-0005)")
	}
	ctx, batal := context.WithTimeout(context.Background(), 10*time.Minute)
	defer batal()
	if err := svc.CekKesehatan(ctx); err != nil {
		log.Fatalf("migrasi data treaty contract out: tidak dapat menjangkau oracle: %v", err)
	}
	lap, err := svc.PindahkanDataTreatyContractOut(ctx)
	log.Printf("laporan migrasi data Treaty Contract Out (skema %s):", svc.SkemaAktif())
	log.Print(lap.String())
	if err != nil {
		log.Fatalf("migrasi data treaty contract out: %v", err)
	}
	log.Print("migrasi data treaty contract out: selesai; sequence diselaraskan")
}
