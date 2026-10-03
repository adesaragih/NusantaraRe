// Package services memuat aturan dagang modul Komite Claim Life.
//
// Arah ketergantungan: handlers -> services -> repository. Paket ini tidak
// pernah mengimpor handlers, dan tidak pernah mengimpor modul lain - yang
// bersama datang dari `inti/`, dan yang dibutuhkan dari Claim Life datang
// lewat `inti/backend/kontrak.KlaimKomite`.
package services

// Akar layanan modul Komite Claim Life - refactor bentuk B.
//
// Untuk apa berkas ini: sampai 30-09-2026 metode-metode Komite menumpang di
// `Service` bersama milik Claim Life, dan memanggil repository Claim Life
// langsung. Kini modul ini punya `Service` sendiri; yang bersama datang dari
// `inti.Dasar`, dan baris klaim Claim Life dibaca serta ditulis lewat kontrak
// yang disambung `cmd/api`. Nol perilaku baru.

import (
	"context"
	"errors"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/kontrak"
)

// Service adalah akar layanan Komite Claim Life.
type Service struct {
	*inti.Dasar
	// klaim - kontrak Claim Life (km3). nil = belum disambung.
	klaim kontrak.KlaimKomite
}

// New membuat Service. db boleh nil bila proses berjalan tanpa Oracle;
// layanan yang memerlukannya akan menolak saat dipanggil.
func New(d *db.DB) *Service {
	return &Service{Dasar: inti.NewDasar(d)}
}

// DariDasar membuat Service di atas akar yang sudah disetel `cmd/api`
// (lingkungan efek keluar, folder unggahan) - satu akar untuk semua modul.
func DariDasar(d *inti.Dasar) *Service {
	return &Service{Dasar: d}
}

// DenganKlaim menyambung kontrak Claim Life yang dipakai Komite.
func (s *Service) DenganKlaim(k kontrak.KlaimKomite) *Service {
	salin := *s
	salin.klaim = k
	return &salin
}

// ErrKlaimBelumDisambung - proses menyentuh Oracle tetapi kontrak Claim Life
// tidak pernah disambung (salah rakit, bukan data).
var ErrKlaimBelumDisambung = errors.New(
	"services: kontrak Claim Life untuk Komite belum disambung (inti/backend/kontrak.KlaimKomite)")

// Klaim mengembalikan kontrak Claim Life yang disambung, atau penolak TERANG
// bila belum - bukan antarmuka nil yang membuat proses panik.
func (s *Service) Klaim() kontrak.KlaimKomite {
	if s.klaim == nil {
		return klaimBelumDisambung{}
	}
	return s.klaim
}

// DenganUnggahanDir menyetel folder berkas unggahan - lihat
// `inti.Dasar.DenganUnggahanDir`.
func (s *Service) DenganUnggahanDir(dir string) *Service {
	salin := *s
	salin.Dasar = s.Dasar.DenganUnggahanDir(dir)
	return &salin
}

// DenganLingkungan menyetel lingkungan efek keluarnya - lihat
// `inti.Dasar.DenganLingkungan`.
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

// CekKesehatan memeriksa apakah Oracle terjangkau - lihat `inti.Dasar.CekKesehatan`.
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
