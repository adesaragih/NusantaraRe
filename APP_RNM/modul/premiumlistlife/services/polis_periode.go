package services

// Periode produksi - tiket 02 PremiumList Life.
//
// Untuk apa berkas ini: menyatukan pembacaan tanggal tutup buku dengan
// aturan periodenya, dan menyediakan jam yang dapat dikendalikan uji.
//
// ⛔ DIBACA SETIAP KALI, nol cache - lihat `repository/polis_periode.go`.
//
// Dibaca sesudah: models/polis_periode.go (aturannya).

import (
	"context"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/penomor"
	"nusantarare/modul/premiumlistlife/repository"
)

// Periode melayani penentuan periode produksi.
type Periode struct {
	svc *Service
	// jam diserahkan supaya aturan periode dapat diuji tanpa menunggu
	// tanggal nyata (AC tiket 02).
	jam func() time.Time
}

// Periode menyusun layanannya.
func (s *Service) Periode() *Periode {
	return &Periode{svc: s, jam: time.Now}
}

// DenganJam mengganti sumber waktunya - dipakai uji.
func (p *Periode) DenganJam(j func() time.Time) *Periode {
	salin := *p
	salin.jam = j
	return &salin
}

// Sekarang menentukan periode produksi untuk saat ini.
//
// Mengembalikan `ErrTanggalTutupBukuKosong` bila tabelnya kosong - dan
// pesannya menyebut tabel sumbernya.
func (p *Periode) Sekarang(ctx context.Context, pelaku inti.Pelaku) (time.Time, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return time.Time{}, err
	}
	if p == nil || p.svc == nil || !p.svc.PunyaDatabase() {
		return time.Time{}, db.ErrTanpaOracle
	}
	tgl, err := repository.NewTutupBuku(p.svc.DB()).Tanggal(ctx)
	if err != nil {
		return time.Time{}, err
	}
	return penomor.PeriodeProduksi(tgl, p.jam())
}
