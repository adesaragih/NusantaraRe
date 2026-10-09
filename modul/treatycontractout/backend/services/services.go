// Package services memuat aturan dagang modul Treaty Contract Out.
//
// Arah ketergantungan: handlers -> services -> repository. Paket ini tidak
// pernah mengimpor handlers, dan tidak pernah mengimpor modul lain - yang
// bersama datang dari `inti/`.
package services

// Akar layanan modul Treaty Contract Out - refactor bentuk B.
//
// Untuk apa berkas ini: sampai 30-09-2026 metode-metode Treaty menumpang di
// `Service` bersama milik Claim Life. Kini modul ini punya `Service` sendiri;
// yang bersama - koneksi, lingkungan efek keluar, folder unggahan, transaksi -
// datang dari `inti.Dasar` yang disematkan. Nol perilaku baru.

import (
	"context"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
)

// Service adalah akar layanan Treaty Contract Out.
type Service struct {
	*inti.Dasar
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
