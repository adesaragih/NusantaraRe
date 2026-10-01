// Command api adalah entry point backend Nusantara Re.
//
// Fase 0 - scaffold. Berkas ini hanya: baca config -> buka koneksi ->
// daftarkan handler -> dengarkan. ⛔ Nol aturan dagang di sini.
//
// main adalah composition root. Refactor bentuk B (30-09-2026): ia membangun
// akar bersama (`inti.Dasar`), meminta daftar modul dari `daftar.Rakit`,
// menyaringnya menurut MODUL_AKTIF (`rakit.go`), lalu memasang rute dan
// pekerja latar modul yang aktif saja. Migrasi selalu dari SEMUA modul
// terdaftar (`daftar.SumberMigrasi`), tidak ikut MODUL_AKTIF.
//
// Struktur tim satu folder per modul (30-09-2026): daftar modulnya
// `inti/backend/daftar`, dibangkitkan dari folder `modul/` (go generate), dan
// kontrak lintas modul disambung perakit menurut pernyataan tiap modul.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/daftar"
	intidb "nusantarare/inti/backend/db"
	"nusantarare/inti/backend/login"
	"nusantarare/inti/backend/migrasi"
)

func main() {
	migrasi := flag.Bool("migrate", false, "jalankan migrasi lalu keluar")
	bongkar := flag.Bool("migrate-down", false,
		"BONGKAR skema uji lalu keluar - MENGHAPUS tabel; perlu ORACLE_SKEMA_UJI=true")
	var baru login.AkunBaru
	var wb string
	flag.StringVar(&baru.ID, "buat-pengguna", "",
		"buat akun M_LOGIN_GO lalu keluar; sandi sementara DICETAK SEKALI dan wajib diganti saat login pertama")
	flag.StringVar(&baru.Nama, "nama", "", "nama tampilan akun (-buat-pengguna)")
	flag.StringVar(&baru.Organisasi, "organisasi", "", "M_ORGANIZATION.CODE (-buat-pengguna)")
	flag.StringVar(&baru.Divisi, "divisi", "", "M_DIVISION.CODE (-buat-pengguna)")
	flag.StringVar(&baru.Unit, "unit", "", "M_UNIT.CODE (-buat-pengguna)")
	flag.StringVar(&wb, "workbasket", "", "WORKBASKET_ID dipisah koma (-buat-pengguna)")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("konfigurasi: %v", err)
	}

	var db *intidb.DB
	if cfg.PunyaOracle() {
		db, err = intidb.Open(cfg)
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
	//
	// Refactor bentuk B (30-09-2026): SATU akar untuk semua modul; setiap
	// modul membangun `Service`-nya sendiri di atasnya (`daftar.Rakit`).
	dasar := inti.NewDasar(db).
		DenganLingkungan(inti.LingkunganDariFlag(cfg.IsPegaProd)).
		DenganUnggahanDir(cfg.UnggahanDir)
	if dasar.PunyaDatabase() {
		log.Printf("oracle: skema %s", dasar.SkemaAktif())
	} else {
		log.Print("oracle: ORACLE_DSN kosong - berjalan tanpa database")
	}

	// Treaty Contract Out tco4 (29-09-2026): `-migrate-data-treaty-contract-out`
	// DIBUANG - modul menulis dan membaca tabel warisan langsung, tidak ada
	// tabel tujuan untuk dipindahi.
	if *migrasi && *bongkar {
		log.Fatal("pilih salah satu: -migrate atau -migrate-down")
	}
	if *bongkar {
		bongkarMigrasi(dasar, cfg)
		return
	}
	if *migrasi {
		jalankanMigrasi(dasar)
		return
	}
	if baru.ID != "" {
		for _, w := range strings.Split(wb, ",") {
			if w = strings.TrimSpace(w); w != "" {
				baru.Workbasket = append(baru.Workbasket, w)
			}
		}
		buatPengguna(dasar, baru)
		return
	}

	catat := func(s string) { log.Print(s) }
	// Refactor bentuk B: modul yang dipasang dipilih MODUL_AKTIF (kosong =
	// semua). Nama yang tidak dikenal menolak menyala. ⛔ SESUDAH cabang
	// -migrate / -migrate-down: migrasi tidak ikut MODUL_AKTIF, jadi salah
	// ketik di sana tidak boleh menghalanginya (temuan /code-review).
	// ⛔ Daftar yang bentuknya salah - mis. kontrak yang dibutuhkan tanpa
	// penyedia - MENOLAK menyala dengan kalimat yang menyebut kontrak dan
	// modulnya, bukan nil diam-diam saat permintaan pertama.
	rakitan, err := daftar.Rakit(dasar, cfg, catat)
	if err != nil {
		log.Fatalf("modul: %v", err)
	}
	terdaftar := rakitan.Modul
	aktif, err := pilihModulAktif(terdaftar, daftar.NamaLama(), cfg.ModulAktif)
	if err != nil {
		log.Fatalf("konfigurasi: %v", err)
	}

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           rakitMux(dasar, terdaftar, aktif, cfg.AuthStub, rakitLogin(dasar, cfg)),
		ReadHeaderTimeout: 10 * time.Second,
	}

	ctx, berhenti := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer berhenti()

	pekerja := jalankanPekerja(ctx, aktif)

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
	tungguPekerja(tutup, pekerja, catat)
}

// buatPengguna adalah titik masuk `-buat-pengguna` (login, keputusan work
// owner 01-10-2026): akun pertama dan akun berikutnya sampai layar kelola
// pengguna ada.
//
// ⛔ Sandi sementara dicetak ke stdout SEKALI dan tidak disimpan di mana pun
// selain sebagai hash; akunnya wajib ganti sandi saat login pertama.
func buatPengguna(svc *inti.Dasar, a login.AkunBaru) {
	if !svc.PunyaDatabase() {
		log.Fatal("buat-pengguna: ORACLE_DSN wajib terisi")
	}
	sandi, err := login.NewLayanan(login.NewGudangOracle(svc.DB()), nil).BuatPengguna(context.Background(), a)
	if err != nil {
		log.Fatalf("buat-pengguna: %v", err)
	}
	fmt.Printf("akun %s dibuat. Sandi sementara (tampil SEKALI, wajib diganti saat login pertama): %s\n", a.ID, sandi)
}

// bongkarMigrasi adalah titik masuk `-migrate-down`.
//
// ⛔ Ia MENGHAPUS tabel, dan karena itu dipagari sama persis dengan test bertag
// db: menolak IS_PEGA_PROD=true, menolak tanpa ORACLE_SKEMA_UJI=true, dan
// menolak skema yang memuat POOLDATA. Pagarnya satu-satunya, tinggal di
// inti/backend/config, supaya jalur ini dan jalur test tidak mungkin berselisih.
//
// Kenapa flag ini ada: sampai 26-09-2026 jalur mundur hanya punya pemanggil
// test. Orang yang ingin membongkar skema uji terpaksa menyalin isi berkas
// *_down.sql ke sqlplus - dan itu MELEWATI pengaman T_MIGRASI, yang membongkar
// hanya langkah yang benar-benar tercatat selesai.
func bongkarMigrasi(svc *inti.Dasar, cfg config.Config) {
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
	// Refactor bentuk B: SEMUA modul terdaftar, tidak bergantung modul aktif.
	lap, err := migrasi.Bongkar(ctx, svc.DB(), daftar.SumberMigrasi()...)
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
func jalankanMigrasi(svc *inti.Dasar) {
	if !svc.PunyaDatabase() {
		log.Fatal("migrasi: ORACLE_DSN wajib terisi")
	}
	ctx, batal := context.WithTimeout(context.Background(), 30*time.Second)
	defer batal()
	if err := svc.CekKesehatan(ctx); err != nil {
		log.Fatalf("migrasi: tidak dapat menjangkau oracle: %v", err)
	}
	// Refactor bentuk B: SEMUA modul terdaftar, tidak bergantung modul aktif.
	lap, err := migrasi.Jalankan(ctx, svc.DB(), daftar.SumberMigrasi()...)
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
