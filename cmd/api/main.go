// Command api adalah entry point backend Nusantara Re.
//
// Fase 0 - scaffold. Berkas ini hanya: baca config -> buka koneksi ->
// daftarkan handler -> dengarkan. ⛔ Nol aturan dagang di sini.
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
		log.Printf("oracle: skema %s", db.Skema())
	} else {
		log.Print("oracle: ORACLE_DSN kosong - berjalan tanpa database")
	}

	if *migrasi {
		jalankanMigrasi(db)
		return
	}

	// Interface dibiarkan nil bila database tidak ada, supaya /healthz
	// melaporkan "tidak dikonfigurasi" dan bukan "tidak terjangkau".
	var kesehatan handlers.Kesehatan
	if db != nil {
		kesehatan = db
	}

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handlers.Router(services.New(db), kesehatan),
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

// jalankanMigrasi adalah titik masuk `make migrate`.
//
// Fase 0 tidak punya migrasi: nol DDL, nol aturan dagang. Berkas migrasi lahir
// bersama tiket yang memilikinya, dan ⛔ tidak pernah dijalankan terhadap
// instance produksi.
func jalankanMigrasi(db *repository.DB) {
	if db == nil {
		log.Fatal("migrasi: ORACLE_DSN wajib terisi")
	}
	ctx, batal := context.WithTimeout(context.Background(), 30*time.Second)
	defer batal()
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("migrasi: tidak dapat menjangkau oracle: %v", err)
	}
	log.Printf("migrasi: belum ada migrasi terdaftar untuk skema %s", db.Skema())
}
