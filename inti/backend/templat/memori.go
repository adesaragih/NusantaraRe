package templat

import (
	"context"
	"sort"
	"sync"
)

// GudangMemori adalah Gudang di memori - untuk uji dan untuk proses tanpa Oracle yang tetap ingin mencoba layar.
type GudangMemori struct {
	mu    sync.Mutex
	versi map[string][]Versi
	isi   map[string]map[int][]byte
	// Waktu - nilai TglUnggah versi baru.
	Waktu string
}

// NewGudangMemori membuat gudang kosong.
func NewGudangMemori() *GudangMemori {
	return &GudangMemori{versi: map[string][]Versi{}, isi: map[string]map[int][]byte{}, Waktu: "04-10-2026 09:30"}
}

// Aktif memenuhi Gudang.
func (g *GudangMemori) Aktif(context.Context) (map[string]Versi, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	hasil := map[string]Versi{}
	for k, vs := range g.versi {
		for _, v := range vs {
			if v.Aktif {
				hasil[k] = v
			}
		}
	}
	return hasil, nil
}

// Riwayat memenuhi Gudang.
func (g *GudangMemori) Riwayat(_ context.Context, kode string) ([]Versi, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	vs := append([]Versi(nil), g.versi[kode]...)
	sort.Slice(vs, func(i, j int) bool { return vs[i].Versi > vs[j].Versi })
	return vs, nil
}

// Isi memenuhi Gudang.
func (g *GudangMemori) Isi(_ context.Context, kode string, versi int) ([]byte, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	b, ok := g.isi[kode][versi]
	if !ok {
		return nil, ErrVersiTidakAda
	}
	return append([]byte(nil), b...), nil
}

// Sisip memenuhi Gudang.
func (g *GudangMemori) Sisip(_ context.Context, kode string, v Versi, isi []byte) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for i := range g.versi[kode] {
		g.versi[kode][i].Aktif = false
	}
	v.Versi = len(g.versi[kode]) + 1
	v.Ukuran = len(isi)
	v.Aktif = true
	v.TglUnggah = g.Waktu
	g.versi[kode] = append(g.versi[kode], v)
	if g.isi[kode] == nil {
		g.isi[kode] = map[int][]byte{}
	}
	g.isi[kode][v.Versi] = append([]byte(nil), isi...)
	return v.Versi, nil
}

// Aktifkan memenuhi Gudang.
func (g *GudangMemori) Aktifkan(_ context.Context, kode string, versi int) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	ada := versi == 0
	for i := range g.versi[kode] {
		g.versi[kode][i].Aktif = g.versi[kode][i].Versi == versi
		ada = ada || g.versi[kode][i].Aktif
	}
	if !ada {
		return ErrVersiTidakAda
	}
	return nil
}
