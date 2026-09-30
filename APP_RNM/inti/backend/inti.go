// Package inti memuat akar layanan yang dipakai bersama oleh setiap modul:
// `Dasar` (koneksi, lingkungan efek keluar, folder unggahan, transaksi),
// antarmuka `Akar`, `Lingkungan`, dan `Pelaku`.
//
// Refactor bentuk B (30-09-2026). Paket di bawah `inti/` adalah SATU-SATUNYA
// kode bersama; paket ini tidak pernah mengimpor `modul/` mana pun.
package backend

// Akar layanan yang dipakai bersama oleh setiap modul - refactor bentuk B.
//
// Untuk apa berkas ini: sampai 30-09-2026 SATU `Service` memegang koneksi,
// lingkungan efek keluar, dan folder unggahan untuk keempat modul sekaligus.
// Bentuk B memberi tiap modul `Service`-nya sendiri; yang benar-benar bersama
// dikumpulkan di `Dasar` dan disematkan (embedded) ke `Service` setiap modul.
//
// ⛔ Nol perilaku baru. Isi setiap metode dipindah apa adanya dari
// `services.go`; yang berubah hanya tempat tinggalnya.
//
// Istilah:
//   - disematkan (embedded) : medan tanpa nama di dalam struct Go. Metode
//     tipe yang disematkan ikut terpanggil lewat tipe luarnya, jadi
//     `s.DalamTransaksi(...)` tetap berlaku di setiap modul.

import (
	"context"

	"nusantarare/inti/backend/db"
)

// Akar adalah yang dibawa setiap `Service` modul.
//
// Fungsi bersama (outbox, resolver layanan, jejak) menerima `Akar`, bukan
// `Service` modul tertentu, supaya paket bersama tidak pernah mengimpor modul.
// Setiap `Service` yang menyematkan `*Dasar` memenuhinya dengan sendirinya.
type Akar interface {
	PunyaDatabase() bool
	DB() *db.DB
	DalamTransaksi(ctx context.Context, fn func(tx *db.Tx) error) error
	Lingkungan() Lingkungan
	UnggahanDir() string
}

// Dasar memegang koneksi, lingkungan efek keluar, dan folder unggahan.
type Dasar struct {
	db *db.DB
	// lingkungan menggerbangi EFEK KELUAR saja, tidak pernah penyimpanan.
	//
	// ⛔ Bawaannya BUKAN produksi, dan itu disengaja: gagal tertutup, bukan
	// gagal terbuka. Proses yang lupa menyetelnya tidak akan mengirim email
	// kepada orang sungguhan.
	lingkungan Lingkungan
	// unggahanDir adalah folder lokal berkas unggahan (butir be). Kosong
	// berarti unggahan GAGAL TERANG - bukan bawaan diam-diam.
	unggahanDir string
}

// NewDasar membuat akar layanan. db boleh nil bila proses berjalan tanpa
// Oracle; layanan yang memerlukannya akan menolak saat dipanggil.
func NewDasar(db *db.DB) *Dasar {
	return &Dasar{db: db, lingkungan: BukanProduksi}
}

// DenganUnggahanDir menyetel folder berkas unggahan.
//
// ⚠️ Dipanggil SEKALI saat proses menyala, dari `cmd/api`, dengan nilai
// dari `config.UnggahanDir`. Tanpa pemanggilan itu tidak ada unggahan yang
// pernah berhasil di mana pun - dan tidak adanya pemanggil itulah yang
// membuat AC lingkungan sempat tercentang secara hampa (lihat
// DenganLingkungan).
func (d *Dasar) DenganUnggahanDir(dir string) *Dasar {
	salin := *d
	salin.unggahanDir = dir
	return &salin
}

// DenganLingkungan menyetel lingkungan efek keluarnya.
//
// ⚠️ Dipanggil SEKALI saat proses menyala, dari `cmd/api`, dengan nilai dari
// `config.IsPegaProd` (ADR-U-0005). Tanpa pemanggilan itu tidak ada efek
// keluar yang pernah berjalan di mana pun - dan tidak adanya pemanggil itulah
// yang membuat AC lingkungan sempat tercentang secara hampa.
func (d *Dasar) DenganLingkungan(l Lingkungan) *Dasar {
	salin := *d
	salin.lingkungan = l
	return &salin
}

// Lingkungan menyebut lingkungan efek keluar proses ini.
func (d *Dasar) Lingkungan() Lingkungan { return d.lingkungan }

// UnggahanDir menyebut folder berkas unggahan; kosong = belum disetel.
func (d *Dasar) UnggahanDir() string { return d.unggahanDir }

// DB mengembalikan koneksi; nil bila proses berjalan tanpa Oracle.
func (d *Dasar) DB() *db.DB { return d.db }

// PunyaDatabase menyatakan apakah proses dikonfigurasi menyentuh Oracle.
func (d *Dasar) PunyaDatabase() bool { return d != nil && d.db != nil }

// CekKesehatan memeriksa apakah Oracle terjangkau.
//
// Ia ada DI SINI, bukan di handlers, supaya arah handlers -> services ->
// repository tidak dipotong: handlers tidak pernah memegang koneksi.
// Mengembalikan ErrTanpaOracle bila database memang tidak dikonfigurasi -
// itu keadaan yang sah di Fase 0, bukan kegagalan.
func (d *Dasar) CekKesehatan(ctx context.Context) error {
	if !d.PunyaDatabase() {
		return db.ErrTanpaOracle
	}
	return d.db.Ping(ctx)
}

// SkemaAktif mengembalikan nama skema yang dipakai, atau teks kosong bila
// tidak ada database.
func (d *Dasar) SkemaAktif() string {
	if !d.PunyaDatabase() {
		return ""
	}
	return d.db.Skema()
}

// DalamTransaksi menjalankan fn di dalam satu transaksi.
//
// Transaksi dibuka dan ditutup DI SINI, bukan di dalam teks SQL
// (ADR-U-0029). Galat apa pun dari fn membatalkan seluruhnya.
func (d *Dasar) DalamTransaksi(ctx context.Context, fn func(tx *db.Tx) error) (err error) {
	if d.db == nil {
		return db.ErrTanpaOracle
	}
	tx, err := d.db.Mulai(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	if err = fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}

// Lingkungan membedakan produksi dari selainnya.
type Lingkungan int

const (
	// BukanProduksi - efek keluar TIDAK dijalankan.
	BukanProduksi Lingkungan = iota
	// Produksi - efek keluar dijalankan.
	Produksi
)

// AdalahProduksi menyatakan lingkungan ini produksi.
//
// ⚠️ Namanya BUKAN `Produksi`: metode dan konstanta bernama sama membuat
// `l.Produksi()` dan `Produksi` terbaca seolah hal yang sama.
func (l Lingkungan) AdalahProduksi() bool { return l == Produksi }

// LingkunganDariFlag menerjemahkan flag konfigurasi menjadi lingkungan.
//
// ⚠️ Terpisah dari `DariTingkatProduksi`. Yang satu membaca kolom Pega
// (`pzProductionLevel`), yang lain membaca konfigurasi kita sendiri
// (`IS_PEGA_PROD`, ADR-U-0005). Keduanya sengaja tidak disatukan: menyatukan
// dua sumber kebenaran atas pertanyaan "apakah ini produksi" berarti salah
// satunya diam-diam menang.
func LingkunganDariFlag(produksi bool) Lingkungan {
	if produksi {
		return Produksi
	}
	return BukanProduksi
}
