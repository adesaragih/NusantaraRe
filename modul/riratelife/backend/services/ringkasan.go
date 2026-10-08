package services

import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"nusantarare/modul/riratelife/backend/models"
	"nusantarare/modul/riratelife/backend/repository"
)

var wib = time.FixedZone("WIB", 7*60*60)

// FormatWaktuPega - bentuk MODIFIEDDATE yang ditulis (`20261005T102343.866 GMT`, ASUMSI A3 - seperti TGLUPDATE
// Pega di tabel warisan lain).
func FormatWaktuPega(t time.Time) string { return t.UTC().Format("20060102T150405.000") + " GMT" }

// TampilTanggal - `Date-Short-Custom-YYYY` b10929: teks Pega (`YYYYMMDDTHHMMSS.mmm GMT` atau `YYYYMMDD`) ke
// `DD-MM-YYYY` WIB; bentuk lain apa adanya.
func TampilTanggal(s string) string {
	s = strings.TrimSpace(s)
	if t, err := time.Parse("20060102T150405.000 MST", s); err == nil {
		return t.In(wib).Format("02-01-2006")
	}
	if t, err := time.Parse("20060102", s); err == nil {
		return t.Format("02-01-2006")
	}
	return s
}

func lengkapi(r models.Ringkasan) models.Ringkasan {
	r.Diubah = TampilTanggal(r.ModifiedDate)
	return r
}

// potongByte memotong teks ke n byte tanpa memotong satu karakter.
func potongByte(s string, n int) string {
	if len(s) <= n {
		return s
	}
	s = s[:n]
	for !utf8.ValidString(s) {
		s = s[:len(s)-1]
	}
	return s
}

// RapikanSaringan - halaman >= 1; kolom urut di luar daftar = UrutBawaan (ID menurun).
func RapikanSaringan(s models.Saringan) models.Saringan {
	s.ID, s.UsedBy, s.Urut = strings.TrimSpace(s.ID), strings.TrimSpace(s.UsedBy), strings.ToLower(strings.TrimSpace(s.Urut))
	if !slices.Contains(models.KolomUrut, s.Urut) {
		s.Urut, s.Turun = models.UrutBawaan, true
	}
	if s.Halaman < 1 {
		s.Halaman = 1
	}
	return s
}

// Daftar - satu halaman grid ringkasan.
func (l *Layanan) Daftar(ctx context.Context, s models.Saringan) (models.Halaman[models.Ringkasan], error) {
	s = RapikanSaringan(s)
	d, total, err := l.gudang.Daftar(ctx, s)
	if err != nil {
		return models.Halaman[models.Ringkasan]{}, err
	}
	for i := range d {
		d[i] = lengkapi(d[i])
	}
	if d == nil {
		d = []models.Ringkasan{}
	}
	return models.Halaman[models.Ringkasan]{Daftar: d, Total: total, Halaman: s.Halaman, Ukuran: models.UkuranHalaman}, nil
}

func petaTidakAda(err error) error {
	if errors.Is(err, repository.ErrTidakAda) {
		return ErrTidakAda
	}
	return err
}

// Buka - satu ringkasan beserta jumlah rate-nya (dialog Delete).
func (l *Layanan) Buka(ctx context.Context, id string) (models.Ringkasan, error) {
	id = strings.TrimSpace(id)
	r, err := l.gudang.Ambil(ctx, nil, id)
	if err != nil {
		return models.Ringkasan{}, petaTidakAda(err)
	}
	n, err := l.gudang.JumlahRate(ctx, nil, id)
	if err != nil {
		return models.Ringkasan{}, err
	}
	r = lengkapi(r)
	r.JumlahRate = &n
	return r, nil
}

// periksaNama - R/I RATE NAME wajib (b1285), dipangkas, paling panjang BatasNama byte.
func periksaNama(nama string) error {
	switch {
	case nama == "":
		return tolak("R/I RATE NAME is required")
	case len(nama) > models.BatasNama:
		return tolak("R/I RATE NAME is longer than %d characters", models.BatasNama)
	}
	return nil
}

// Simpan - Save (`AddToListSummary_Act`): Add bila id kosong, selain itu Edit (`EditListSummary_DT` b11136: ID,
// UsedBy). Edit nama ikut mengganti salinan nama di baris rate ringkasan itu (ASUMSI A6).
func (l *Layanan) Simpan(ctx context.Context, a Aktor, id string, isi models.Isian) (models.Ringkasan, error) {
	if err := wajibPenuh(a); err != nil {
		return models.Ringkasan{}, err
	}
	id, nama := strings.TrimSpace(id), strings.TrimSpace(isi.UsedBy)
	if err := periksaNama(nama); err != nil {
		return models.Ringkasan{}, err
	}
	err := l.tx(ctx, func(tx *dbTx) error {
		var lama *models.Ringkasan
		if id != "" {
			r, err := l.gudang.Ambil(ctx, tx, id)
			if err != nil {
				return err
			}
			lama = &r
		}
		ganda, err := l.gudang.PemakaiNama(ctx, tx, nama, id)
		if err != nil {
			return err
		}
		if len(ganda) > 0 {
			return tolak("R/I RATE NAME %s is already used by ID %s", nama, ganda[0].ID)
		}
		r := models.Ringkasan{ID: id, UsedBy: nama, OperatorID: potongByte(a.AkunID, models.BatasNama),
			ModifiedDate: FormatWaktuPega(l.jam())}
		if lama != nil {
			if err := l.gudang.UbahRingkasan(ctx, tx, r); err != nil {
				return err
			}
			if lama.UsedBy != nama {
				_, err := l.gudang.UbahNamaRate(ctx, tx, id, nama)
				return err
			}
			return nil
		}
		baru, err := (&pemberiID{l: l, ringkasan: true}).baru(ctx, tx)
		if err != nil {
			return err
		}
		id, r.ID = baru, baru
		return l.gudang.SisipRingkasan(ctx, tx, r)
	})
	if err != nil {
		return models.Ringkasan{}, petaTidakAda(err)
	}
	r, err := l.gudang.Ambil(ctx, nil, id)
	return lengkapi(r), petaTidakAda(err)
}

// Hapus - Delete (`DeleteSummaryDetail`, DeleteID=.ID): ringkasan BESERTA baris rate ber-IDUSEDBY sama, satu transaksi.
func (l *Layanan) Hapus(ctx context.Context, a Aktor, id string) (models.HasilHapus, error) {
	if err := wajibPenuh(a); err != nil {
		return models.HasilHapus{}, err
	}
	id = strings.TrimSpace(id)
	var n int
	err := l.tx(ctx, func(tx *dbTx) error {
		if _, err := l.gudang.Ambil(ctx, tx, id); err != nil {
			return err
		}
		var err error
		if n, err = l.gudang.HapusRate(ctx, tx, id); err != nil {
			return err
		}
		return l.gudang.HapusRingkasan(ctx, tx, id)
	})
	if err != nil {
		return models.HasilHapus{}, petaTidakAda(err)
	}
	return models.HasilHapus{ID: id, RateTerhapus: n}, nil
}

// DaftarRate - satu halaman Rate Detail ringkasan id (20 per halaman, `pyPageSizeOther` b10206).
func (l *Layanan) DaftarRate(ctx context.Context, id string, halaman int) (models.Halaman[models.Rate], error) {
	id = strings.TrimSpace(id)
	if halaman < 1 {
		halaman = 1
	}
	if _, err := l.gudang.Ambil(ctx, nil, id); err != nil {
		return models.Halaman[models.Rate]{}, petaTidakAda(err)
	}
	d, total, err := l.gudang.DaftarRate(ctx, id, halaman)
	if err != nil {
		return models.Halaman[models.Rate]{}, err
	}
	if d == nil {
		d = []models.Rate{}
	}
	return models.Halaman[models.Rate]{Daftar: d, Total: total, Halaman: halaman, Ukuran: models.UkuranHalamanRate}, nil
}
