package services

// Tiket 03 - popup `View Old Policy` (`ViewOldPolicy_EDM`, `_QP`, `_TP`, `_TR`)
// dan rincian peserta (`PL_DetailAction` → `PL_Detail_Sec`, `RetroLife`).
//
// ⚠️ `GetOldDetail_EDM` (4 dari 5 langkah `//`) tidak ditiru apa adanya; yang
// hidup hanya langkah 2 b384 `Obj-Open-By-Handle` `OldCaseID` - membuka versi
// LAMA. Di sini: versi sebelumnya dibaca lagi dari sumbernya (spec §4, §16:
// `OldData` tidak disimpan). Kasus terbuka → versi berjalan polis (belum ada
// versi lain di antaranya, gerbang 3); kasus resmi → versi di bawah `PROD_KE`-nya.

import (
	"context"
	"errors"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/repository"
)

// ErrPesertaTidakAda - peserta bukan milik kasus itu, atau tidak ada.
var ErrPesertaTidakAda = errors.New("services: participant not found in this endorsement case")

// PolisLama membaca isi popup polis lama satu kasus.
func (l *Layanan) PolisLama(ctx context.Context, p inti.Pelaku, id string, halaman, ukuran int) (models.PolisLama, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.PolisLama{}, err
	}
	k, err := l.muatKasus(ctx, id)
	if err != nil {
		return models.PolisLama{}, err
	}
	sebelum := 0
	if k.Status == models.StatusKasusSelesai {
		sebelum = k.ProdKe
	}
	v, ada, err := l.gudang.VersiBerjalan(ctx, nil, k.NomorPolis, sebelum)
	if err != nil {
		return models.PolisLama{}, err
	}
	hasil := models.PolisLama{Tipe: k.Kepala["TYPE"], Peserta: []models.Peserta{}, Rekap: []map[string]string{}}
	if !ada {
		return hasil, nil
	}
	halaman, ukuran = jepit(halaman, ukuran)
	peserta, total, err := l.gudang.PesertaVersi(ctx, v, k.NomorPolis, halaman, ukuran)
	if err != nil {
		return models.PolisLama{}, err
	}
	rekap, err := l.gudang.RekapPolis(ctx, k.NomorPolis)
	if err != nil {
		return models.PolisLama{}, err
	}
	hasil.Sumber, hasil.ProdKe, hasil.Total = v.Jenis, v.ProdKe, total
	if peserta != nil {
		hasil.Peserta = peserta
	}
	if rekap != nil {
		hasil.Rekap = rekap
	}
	return hasil, nil
}

// RincianPeserta membaca `PL_Detail_Sec` satu peserta kasus.
func (l *Layanan) RincianPeserta(ctx context.Context, p inti.Pelaku, id, pesertaID string) (models.RincianPeserta, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.RincianPeserta{}, err
	}
	if _, err := l.muatKasus(ctx, id); err != nil {
		return models.RincianPeserta{}, err
	}
	r, err := l.gudang.RincianPeserta(ctx, id, pesertaID)
	if errors.Is(err, repository.ErrTidakAda) {
		return models.RincianPeserta{}, fmt.Errorf("%w: %q", ErrPesertaTidakAda, pesertaID)
	}
	return r, err
}
