// Package services memuat aturan dagang, batas transaksi, dan penegakan
// wewenang.
//
// Arah ketergantungan: handlers -> services -> repository. Paket ini tidak
// pernah mengimpor handlers.
//
// Fase 0 - scaffold: nol aturan dagang. Yang ada di sini adalah batas
// transaksi dan bentuk penegakan wewenang yang harus diikuti setiap tiket.
package services

import (
	"context"
	"errors"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/kontrak"
)

var (
	// ErrWajibIsi dikembalikan bila medan wajib kosong. Seluruh kolom
	// nullable di database; kewajiban isi ditegakkan DI SINI (ADR-U-0027).
	ErrWajibIsi = errors.New("services: medan wajib belum terisi")
)

// Service adalah akar seluruh layanan.
//
// Refactor bentuk B (30-09-2026): yang bersama - koneksi, lingkungan efek
// keluar, folder unggahan, transaksi - tinggal di `inti.Dasar` yang
// disematkan. Yang tersisa di sini hanya milik modul.
type Service struct {
	*inti.Dasar
	// pembacaPolis - pembaca polis ringkas PremiumList Life (butir pl4/av),
	// disambung `cmd/api` lewat `inti/backend/kontrak`. nil = belum disambung.
	pembacaPolis kontrak.PembacaPolis
}

// DenganUnggahanDir menyetel folder berkas unggahan - lihat
// `Dasar.DenganUnggahanDir`.
func (s *Service) DenganUnggahanDir(dir string) *Service {
	salin := *s
	salin.Dasar = s.Dasar.DenganUnggahanDir(dir)
	return &salin
}

// New membuat Service. db boleh nil bila proses berjalan tanpa Oracle;
// layanan yang memerlukannya akan menolak saat dipanggil.
func New(db *db.DB) *Service {
	return &Service{Dasar: inti.NewDasar(db)}
}

// DariDasar membuat Service di atas akar yang sudah disetel `cmd/api`
// (lingkungan efek keluar, folder unggahan) - satu akar untuk semua modul.
func DariDasar(d *inti.Dasar) *Service {
	return &Service{Dasar: d}
}

// DenganPembacaPolis menyambung pembaca polis PremiumList Life.
//
// Refactor bentuk B (30-09-2026): Claim Life dulu membaca `T_PREMIUM_LIST`
// lewat repository PremiumList secara langsung. Kini ia hanya mengenal
// antarmuka `kontrak.PembacaPolis`; `cmd/api` yang menyambungnya.
func (s *Service) DenganPembacaPolis(p kontrak.PembacaPolis) *Service {
	salin := *s
	salin.pembacaPolis = p
	return &salin
}

// ErrPembacaPolisBelumDisambung - proses menyentuh Oracle tetapi pembaca
// polis PremiumList tidak pernah disambung (salah rakit, bukan data).
var ErrPembacaPolisBelumDisambung = errors.New(
	"services: pembaca polis PremiumList belum disambung (inti/backend/kontrak.PembacaPolis)")

// PembacaPolis mengembalikan pembaca polis yang disambung, atau pembaca yang
// menolak TERANG bila belum - bukan antarmuka nil yang membuat proses panik.
func (s *Service) PembacaPolis() kontrak.PembacaPolis {
	if s.pembacaPolis == nil {
		return pembacaPolisBelumDisambung{}
	}
	return s.pembacaPolis
}

type pembacaPolisBelumDisambung struct{}

func (pembacaPolisBelumDisambung) Ringkas(context.Context, string) (kontrak.PolisRingkas, error) {
	return kontrak.PolisRingkas{}, ErrPembacaPolisBelumDisambung
}

// DenganLingkungan menyetel lingkungan efek keluarnya - lihat
// `Dasar.DenganLingkungan`.
func (s *Service) DenganLingkungan(l inti.Lingkungan) *Service {
	salin := *s
	salin.Dasar = s.Dasar.DenganLingkungan(l)
	return &salin
}

// PunyaDatabase menyatakan apakah proses dikonfigurasi menyentuh Oracle.
//
// ⚠️ Ditulis ulang di sini, bukan hanya diwarisi dari `Dasar`: metode yang
// diwarisi lewat penyematan menyentuh `s.Dasar` lebih dahulu, sehingga
// `(*Service)(nil).PunyaDatabase()` akan panik - padahal sejak awal ia aman
// dipanggil pada nil. Begitu pula CekKesehatan dan SkemaAktif.
func (s *Service) PunyaDatabase() bool { return s != nil && s.Dasar.PunyaDatabase() }

// CekKesehatan memeriksa apakah Oracle terjangkau - lihat `Dasar.CekKesehatan`.
func (s *Service) CekKesehatan(ctx context.Context) error {
	if !s.PunyaDatabase() {
		return db.ErrTanpaOracle
	}
	return s.Dasar.CekKesehatan(ctx)
}

// SkemaAktif mengembalikan nama skema yang dipakai, atau teks kosong bila
// tidak ada database.
func (s *Service) SkemaAktif() string {
	if !s.PunyaDatabase() {
		return ""
	}
	return s.Dasar.SkemaAktif()
}
