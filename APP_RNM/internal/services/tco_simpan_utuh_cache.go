package services

// Cache master per permintaan simpan utuh - temuan /code-review (tiket 09).
//
// Untuk apa berkas ini: penulis klausul membaca master utuh per baris
// (`TREATYDESC`, daftar jenis reasuransi, dan seluruh riwayat kurs). Dalam
// simpan utuh itu berulang untuk tiap baris SELAGI kontrak/tahun terkunci.
// Pembungkus ini membaca tiap master SEKALI per permintaan. Hasil galat tidak
// disimpan - galat tetap terang pada baris yang memicunya.

import (
	"context"
	"sync"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
)

type kursSekaliTCO struct {
	dasar PembacaKursTCO
	mu    sync.Mutex
	isi   map[string]models.KursTCO
}

func (k *kursSekaliTCO) Berlaku(ctx context.Context, tahun models.TahunTreaty) (models.KursTCO, error) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if v, ada := k.isi[tahun.ID]; ada {
		return v, nil
	}
	v, err := k.dasar.Berlaku(ctx, tahun)
	if err != nil {
		return models.KursTCO{}, err
	}
	k.isi[tahun.ID] = v
	return v, nil
}

type jenisSekaliTCO struct {
	dasar PembacaJenisReasuransiTCO
	mu    sync.Mutex
	isi   []repository.JenisReasuransiTCO
	ada   bool
}

func (j *jenisSekaliTCO) DaftarNonLife(ctx context.Context) ([]repository.JenisReasuransiTCO, error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.ada {
		return j.isi, nil
	}
	v, err := j.dasar.DaftarNonLife(ctx)
	if err != nil {
		return nil, err
	}
	j.isi, j.ada = v, true
	return v, nil
}

type masterKlausulSekaliTCO struct {
	PembacaMasterKlausulTCO
	mu  sync.Mutex
	isi map[string][]repository.JenisKlausulMasterTCO
}

func (m *masterKlausulSekaliTCO) JenisKlausul(ctx context.Context, isXOL string) ([]repository.JenisKlausulMasterTCO, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if v, ada := m.isi[isXOL]; ada {
		return v, nil
	}
	v, err := m.PembacaMasterKlausulTCO.JenisKlausul(ctx, isXOL)
	if err != nil {
		return nil, err
	}
	m.isi[isXOL] = v
	return v, nil
}

// denganMasterSekali mengembalikan salinan penulis klausul yang membaca master
// sekali per permintaan.
func (l *KlausulTCO) denganMasterSekali() *KlausulTCO {
	s := l.salin()
	s.kurs = &kursSekaliTCO{dasar: l.kurs, isi: map[string]models.KursTCO{}}
	s.jenis = &jenisSekaliTCO{dasar: l.jenis}
	s.master = &masterKlausulSekaliTCO{PembacaMasterKlausulTCO: l.master, isi: map[string][]repository.JenisKlausulMasterTCO{}}
	return s
}
