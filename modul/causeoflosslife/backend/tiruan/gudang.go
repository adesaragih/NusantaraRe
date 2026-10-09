// Package tiruan - Gudang Cause Of Loss Life di memori untuk uji layanan dan handler. Data UJI buatan.
package tiruan

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"sync"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/causeoflosslife/backend/models"
	"nusantarare/modul/causeoflosslife/backend/repository"
)

// Transaksi tiruan - fn dijalankan langsung; galat = keadaan dikembalikan (rollback).
type Transaksi struct{ G *Gudang }

// Jalankan memenuhi services.Transaksi.
func (t Transaksi) Jalankan(_ context.Context, fn func(tx *db.Tx) error) error {
	t.G.mu.Lock()
	salinan := make(map[string]models.CauseOfLoss, len(t.G.Baris))
	for k, v := range t.G.Baris {
		salinan[k] = v
	}
	t.G.mu.Unlock()
	if err := fn(nil); err != nil {
		t.G.mu.Lock()
		t.G.Baris = salinan
		t.G.mu.Unlock()
		return err
	}
	return nil
}

// Gudang tiruan.
type Gudang struct {
	mu    sync.Mutex
	Baris map[string]models.CauseOfLoss
	// Seq - nomor M_CAUSEOFLOSS_LIFE_SEQ berikutnya.
	Seq int
}

// Contoh - bentuk DEV dengan nilai UJI: 100001 kosong (NULL, seperti DEV), 100002-100004 terisi; sequence berikutnya 5.
func Contoh() *Gudang {
	g := &Gudang{Baris: map[string]models.CauseOfLoss{}, Seq: 5}
	for _, c := range []models.CauseOfLoss{
		{ID: "100001", CauseOfLoss: ""},
		{ID: "100002", CauseOfLoss: "UJI SAKIT"},
		{ID: "100003", CauseOfLoss: "UJI KECELAKAAN"},
		{ID: "100004", CauseOfLoss: "UJI SEMUA SEBAB"},
	} {
		g.Baris[c.ID] = c
	}
	return g
}

func angka(s string) int {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return -1
	}
	return n
}

// Daftar memenuhi services.Gudang - ID angka menaik, 10 per halaman.
func (g *Gudang) Daftar(_ context.Context, s models.Saringan) ([]models.CauseOfLoss, int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := make([]models.CauseOfLoss, 0, len(g.Baris))
	for _, c := range g.Baris {
		out = append(out, c)
	}
	sort.SliceStable(out, func(a, b int) bool { return angka(out[a].ID) < angka(out[b].ID) })
	total := len(out)
	awal := min((s.Halaman-1)*models.UkuranHalaman, total)
	return out[awal:min(awal+models.UkuranHalaman, total)], total, nil
}

// Ambil memenuhi services.Gudang.
func (g *Gudang) Ambil(_ context.Context, _ *db.Tx, id string) (models.CauseOfLoss, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	c, ok := g.Baris[id]
	if !ok {
		return models.CauseOfLoss{}, repository.ErrTidakAda
	}
	return c, nil
}

// PemakaiNama memenuhi services.Gudang - seperti SQL: nama NULL / kosong tidak pernah cocok.
func (g *Gudang) PemakaiNama(_ context.Context, _ *db.Tx, nama, kecualiID string) ([]models.CauseOfLoss, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var out []models.CauseOfLoss
	for _, c := range g.Baris {
		if c.ID != kecualiID && c.CauseOfLoss != "" && models.KunciTeks(c.CauseOfLoss) == models.KunciTeks(nama) {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(a, b int) bool { return out[a].ID < out[b].ID })
	return out, nil
}

// NomorBaru memenuhi services.Gudang - NEXTVAL.
func (g *Gudang) NomorBaru(_ context.Context, _ *db.Tx) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	n := g.Seq
	g.Seq++
	return strconv.Itoa(n), nil
}

// AdaID memenuhi services.Gudang.
func (g *Gudang) AdaID(_ context.Context, _ *db.Tx, id string) (bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	_, ada := g.Baris[id]
	return ada, nil
}

// Sisip memenuhi services.Gudang - seperti PK.
func (g *Gudang) Sisip(_ context.Context, _ *db.Tx, c models.CauseOfLoss) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ada := g.Baris[c.ID]; ada {
		return repository.ErrKembar
	}
	g.Baris[c.ID] = c
	return nil
}

// Ubah memenuhi services.Gudang.
func (g *Gudang) Ubah(_ context.Context, _ *db.Tx, c models.CauseOfLoss) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ada := g.Baris[c.ID]; !ada {
		return repository.ErrTidakAda
	}
	g.Baris[c.ID] = c
	return nil
}
