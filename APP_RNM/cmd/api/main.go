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

	svc := services.New(db)
	if svc.PunyaDatabase() {
		log.Printf("oracle: skema %s", svc.SkemaAktif())
	} else {
		log.Print("oracle: ORACLE_DSN kosong - berjalan tanpa database")
	}

	if *migrasi && *bongkar {
		log.Fatal("pilih salah satu: -migrate atau -migrate-down, tidak keduanya")
	}
	if *bongkar {
		bongkarMigrasi(svc, cfg)
		return
	}
	if *migrasi {
		jalankanMigrasi(svc)
		return
	}

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handlers.Router(svc),
		ReadHeaderTimeout: 10 * time.Second,
	}

	ctx, berhenti := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer berhenti()

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
