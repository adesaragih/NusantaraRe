// Package services memuat aturan dagang modul PremiumList Life.
//
// Arah ketergantungan: handlers -> services -> repository. Paket ini tidak
// pernah mengimpor handlers, dan tidak pernah mengimpor modul lain - yang
// bersama datang dari `inti/`, dan yang diberikan ke modul lain lewat
// `inti/kontrak`.
package services

// Akar layanan modul PremiumList Life - refactor bentuk B.
//
// Untuk apa berkas ini: sampai 30-09-2026 metode-metode PremiumList menumpang
// di `Service` bersama milik Claim Life. Kini modul ini punya `Service`
// sendiri; yang bersama - koneksi, lingkungan efek keluar, folder unggahan,
// transaksi - datang dari `inti.Dasar` yang disematkan. Nol perilaku baru.

import (
	"context"

	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/inti/kontrak"
	"nusantarare/modul/premiumlist/repository"
)

// Service adalah akar layanan PremiumList Life.
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

// PembacaPolis menyediakan pembaca polis ringkas modul ini untuk Claim Life
// (butir pl4/av) - `inti/kontrak.PembacaPolis`.
//
// Refactor bentuk B (30-09-2026): Claim Life dulu memanggil
// `repository.NewRingkasPolisLife` langsung. Kini `cmd/api` menyerahkan hasil
// fungsi ini kepada Claim Life; pembacaan dan kalimat galatnya sama.
func PembacaPolis(svc *Service) kontrak.PembacaPolis {
	return repository.NewRingkasPolisLife(svc.DB())
}
