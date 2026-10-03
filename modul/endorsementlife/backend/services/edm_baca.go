package services

// Pembacaan: kotak masuk (`InboxEndorsementLife`), kepala kasus dan peserta
// (`InputEDMLife`).
//
// ⚠️ Nol gerbang PERAN: `InboxEDMLife` (kelas `Assign-Worklist`) tidak
// menyaring per operator atau peran - ia daftar kerja seluruh inputor Life.

import (
	"context"
	"errors"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/repository"
)

// Halaman adalah satu halaman grid beserta cacah seluruh barisnya.
type Halaman[T any] struct {
	Baris   []T `json:"baris"`
	Total   int `json:"total"`
	Halaman int `json:"halaman"`
	Ukuran  int `json:"ukuran"`
}

// ErrKasusTidakAda - pengenal bukan kasus Endorsement Life, atau tidak ada.
var ErrKasusTidakAda = errors.New("services: endorsement case not found")

// Inbox membaca satu halaman kasus terbuka.
func (l *Layanan) Inbox(ctx context.Context, p inti.Pelaku, halaman, ukuran int) (Halaman[models.BarisInbox], error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return Halaman[models.BarisInbox]{}, err
	}
	halaman, ukuran = jepit(halaman, ukuran)
	b, total, err := l.gudang.Inbox(ctx, halaman, ukuran)
	if err != nil {
		return Halaman[models.BarisInbox]{}, err
	}
	if b == nil {
		b = []models.BarisInbox{}
	}
	return Halaman[models.BarisInbox]{Baris: b, Total: total, Halaman: halaman, Ukuran: ukuran}, nil
}

// muatKasus membaca kepala kasus beserta tanda-tanda layarnya.
func (l *Layanan) muatKasus(ctx context.Context, id string) (models.Kasus, error) {
	if !models.KasusEDM(id) {
		return models.Kasus{}, fmt.Errorf("%w: %q", ErrKasusTidakAda, id)
	}
	k, err := l.gudang.AmbilKasus(ctx, nil, id, false)
	if errors.Is(err, repository.ErrTidakAda) {
		return models.Kasus{}, fmt.Errorf("%w: %q", ErrKasusTidakAda, id)
	}
	if err != nil {
		return models.Kasus{}, err
	}
	if k.Cacah, err = l.gudang.CacahPeserta(ctx, nil, id); err != nil {
		return models.Kasus{}, err
	}
	k.CSVTerkunci = k.Cacah[models.StatusNew] > 0
	if k.SudahSimpan, err = l.gudang.AdaRekap(ctx, nil, id); err != nil {
		return models.Kasus{}, err
	}
	if k.Rekap, err = l.gudang.RekapKasus(ctx, nil, id); err != nil {
		return models.Kasus{}, err
	}
	if k.Rekap == nil {
		k.Rekap = []map[string]string{}
	}
	if k.Riwayat, err = l.gudang.Riwayat(ctx, nil, id); err != nil {
		return models.Kasus{}, err
	}
	if k.Riwayat == nil {
		k.Riwayat = []models.BarisRiwayat{}
	}
	return k, nil
}

// BacaKasus membaca kepala kasus `InputEDMLife`.
func (l *Layanan) BacaKasus(ctx context.Context, p inti.Pelaku, id string) (models.Kasus, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.Kasus{}, err
	}
	return l.muatKasus(ctx, id)
}

// DaftarPeserta membaca satu halaman peserta kasus.
func (l *Layanan) DaftarPeserta(ctx context.Context, p inti.Pelaku, id string, halaman, ukuran int) (Halaman[models.Peserta], error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return Halaman[models.Peserta]{}, err
	}
	k, err := l.muatKasus(ctx, id)
	if err != nil {
		return Halaman[models.Peserta]{}, err
	}
	halaman, ukuran = jepit(halaman, ukuran)
	b, err := l.gudang.DaftarPeserta(ctx, id, halaman, ukuran)
	if err != nil {
		return Halaman[models.Peserta]{}, err
	}
	if b == nil {
		b = []models.Peserta{}
	}
	total := 0
	for _, c := range k.Cacah {
		total += c
	}
	return Halaman[models.Peserta]{Baris: b, Total: total, Halaman: halaman, Ukuran: ukuran}, nil
}
