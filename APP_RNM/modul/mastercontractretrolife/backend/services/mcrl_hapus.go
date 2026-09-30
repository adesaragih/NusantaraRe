package services

// Hapus berjenjang (paket 7, tiket 09) - padanan `DeleteTreatyLimit_Act`,
// `DeleteSecurityLife_Act`, `DeleteSecurityReinsurerLife_Act`,
// `DeleteRowBusiness` (keempatnya `DELETE … WHERE ID =` datar di Pega) dengan
// penyimpangan 3: kaskade + popup konfirmasi SEBELUM apa pun terhapus.
//
//	kontrak    (1) security reinsurernya, (2) reinsurer dan business, (3) kontrak
//	reinsurer  (1) securitynya, (2) reinsurer
//	security, business - daun
//	tahun      TIDAK ADA - tahun treaty abadi
//
// ⛔ K2: DEV nol FK - kaskade DI GO, satu transaksi, anak lebih dulu.
// ⛔ Popup membawa cacahan yang dilihat pengguna; cacahan dihitung ulang DI
// DALAM transaksi, dan yang berbeda = 409 tanpa satu baris pun terhapus.
// ⛔ Pesan sukses VERBATIM Pega; jejak = satu baris log berisi cacah (K6).

import (
	"context"
	"errors"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/mastercontractretrolife/backend/models"
)

// JenisHapus - entitas yang punya jalur hapus.
type JenisHapus string

// Empat jenis yang dapat dihapus - tahun treaty tidak termasuk.
const (
	HapusKontrak   JenisHapus = "kontrak"
	HapusReinsurer JenisHapus = "reinsurer"
	HapusSecurity  JenisHapus = "security"
	HapusBusiness  JenisHapus = "business"
)

// ErrJenisHapusTidakAda - jenis tanpa jalur hapus (tahun treaty abadi).
var ErrJenisHapusTidakAda = errors.New("services: this kind of row has no delete path (treaty years are permanent)")

// JawabanDampak - isi popup konfirmasi hapus.
type JawabanDampak struct {
	Jenis  JenisHapus    `json:"jenis"`
	ID     string        `json:"id"`
	Dampak models.Dampak `json:"dampak"`
}

// HasilHapus - jawaban hapus yang sudah dijalankan.
type HasilHapus struct {
	Pesan    string        `json:"pesan"`
	Terhapus models.Dampak `json:"terhapus"`
}

// dampak menghitung anak yang ikut terhapus; induk yang tidak ada = 404.
func (l *Layanan) dampak(ctx context.Context, tx *db.Tx, jenis JenisHapus, id string) (models.Dampak, error) {
	switch jenis {
	case HapusKontrak:
		if _, err := l.ambilKontrak(ctx, tx, id); err != nil {
			return models.Dampak{}, err
		}
		return l.gudang.DampakHapusKontrak(ctx, tx, id)
	case HapusReinsurer:
		if _, err := l.ambilReinsurer(ctx, tx, id); err != nil {
			return models.Dampak{}, err
		}
		return l.gudang.DampakHapusReinsurer(ctx, tx, id)
	case HapusSecurity:
		_, err := l.gudang.AmbilSecurity(ctx, tx, id)
		return models.Dampak{}, tidakAda(err, ErrSecurityTidakAda, id)
	case HapusBusiness:
		_, err := l.gudang.AmbilBusiness(ctx, tx, id)
		return models.Dampak{}, tidakAda(err, ErrBusinessTidakAda, id)
	}
	return models.Dampak{}, fmt.Errorf("%w: %q", ErrJenisHapusTidakAda, jenis)
}

// DampakHapus - isi popup: apa dan berapa yang ikut terhapus (nol tulisan).
func (l *Layanan) DampakHapus(ctx context.Context, p inti.Pelaku, jenis JenisHapus, id string) (JawabanDampak, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return JawabanDampak{}, err
	}
	d, err := l.dampak(ctx, nil, jenis, id)
	return JawabanDampak{Jenis: jenis, ID: id, Dampak: d}, err
}

// Hapus menjalankan hapus yang SUDAH dikonfirmasi dengan cacahan `dikonfirmasi`.
func (l *Layanan) Hapus(ctx context.Context, p inti.Pelaku, jenis JenisHapus, id string, dikonfirmasi models.Dampak) (HasilHapus, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilHapus{}, err
	}
	var hasil HasilHapus
	err := l.tx(ctx, func(tx *db.Tx) error {
		kini, err := l.dampak(ctx, tx, jenis, id)
		if err != nil {
			return err
		}
		if kini != dikonfirmasi {
			return fmt.Errorf("%w (confirmed %+v, now %+v)", ErrDampakBerubah, dikonfirmasi, kini)
		}
		var terhapus models.Dampak
		switch jenis {
		case HapusKontrak:
			terhapus, err = l.gudang.HapusKontrak(ctx, tx, id)
			hasil.Pesan = PesanHapusBerhasil
		case HapusReinsurer:
			terhapus, err = l.gudang.HapusReinsurer(ctx, tx, id)
			hasil.Pesan = PesanHapusBerhasil
		case HapusSecurity:
			err = l.gudang.HapusSecurity(ctx, tx, id)
			hasil.Pesan = PesanHapusBerhasil
		case HapusBusiness:
			err = l.gudang.HapusBusiness(ctx, tx, id)
			hasil.Pesan = PesanHapusBusiness(id)
		}
		if err != nil {
			return err
		}
		if terhapus != kini {
			// Lapis kedua: penghapus menyentuh jumlah lain dari yang dihitung -
			// batalkan seluruhnya, jangan biarkan yatim atau kelebihan.
			return fmt.Errorf("services: delete touched %+v but %+v was counted: %w", terhapus, kini, ErrDampakBerubah)
		}
		hasil.Terhapus = terhapus
		return nil
	})
	if err != nil {
		return HasilHapus{}, err
	}
	l.catat(fmt.Sprintf("master contract retro life: deleted %s %s with %d security, %d reinsurer, %d business rows",
		jenis, id, hasil.Terhapus.Security, hasil.Terhapus.Reinsurer, hasil.Terhapus.Business))
	return hasil, nil
}
