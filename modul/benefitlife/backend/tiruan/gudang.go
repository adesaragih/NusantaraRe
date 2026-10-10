// Package tiruan - Gudang Benefit di memori untuk uji layanan dan handler. Data UJI buatan.
package tiruan

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"sync"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/benefitlife/backend/models"
	"nusantarare/modul/benefitlife/backend/repository"
)

// Transaksi tiruan - fn dijalankan langsung; galat = keadaan dikembalikan (rollback).
type Transaksi struct{ G *Gudang }

// Jalankan memenuhi services.Transaksi: salinan keadaan dipulihkan bila fn gagal.
func (t Transaksi) Jalankan(_ context.Context, fn func(tx *db.Tx) error) error {
	t.G.mu.Lock()
	salinan := make(map[string]models.Benefit, len(t.G.Baris))
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
	Baris map[string]models.Benefit
	// Seq - nomor M_BENEFIT_LIFE_SEQ berikutnya.
	Seq int
}

// Contoh - ID UJI 100001, 100002, 100004 (rumus '1' || LPAD(seq, 5)); sequence berikutnya 12 (ID 100012 kosong).
func Contoh() *Gudang {
	g := &Gudang{Baris: map[string]models.Benefit{}, Seq: 12}
	for _, b := range []models.Benefit{
		{ID: "100001", Benefit: "UJI RAWAT INAP"},
		{ID: "100002", Benefit: "UJI MENINGGAL DUNIA"},
		{ID: "100004", Benefit: "UJI CACAT TETAP"},
	} {
		g.Baris[b.ID] = b
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

// Daftar memenuhi services.Gudang - saring memuat, ID angka (bawaan menurun), satu halaman.
func (g *Gudang) Daftar(_ context.Context, s models.Saringan) ([]models.Benefit, int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var out []models.Benefit
	for _, b := range g.Baris {
		if s.ID != "" && !strings.Contains(strings.ToUpper(b.ID), strings.ToUpper(s.ID)) {
			continue
		}
		if s.Benefit != "" && !strings.Contains(strings.ToUpper(b.Benefit), strings.ToUpper(s.Benefit)) {
			continue
		}
		out = append(out, b)
	}
	sort.SliceStable(out, func(a, b int) bool {
		if s.Naik {
			return angka(out[a].ID) < angka(out[b].ID)
		}
		return angka(out[a].ID) > angka(out[b].ID)
	})
	total := len(out)
	awal := min((s.Halaman-1)*models.UkuranHalaman, total)
	return out[awal:min(awal+models.UkuranHalaman, total)], total, nil
}

// Ambil memenuhi services.Gudang.
func (g *Gudang) Ambil(_ context.Context, _ *db.Tx, id string) (models.Benefit, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	b, ok := g.Baris[id]
	if !ok {
		return models.Benefit{}, repository.ErrTidakAda
	}
	return b, nil
}

// NomorBaru memenuhi services.Gudang - NEXTVAL (naik satu).
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

// Sisip memenuhi services.Gudang - seperti PK: ID kembar = ErrKembar.
func (g *Gudang) Sisip(_ context.Context, _ *db.Tx, b models.Benefit) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ada := g.Baris[b.ID]; ada {
		return repository.ErrKembar
	}
	g.Baris[b.ID] = b
	return nil
}

// Ubah memenuhi services.Gudang.
func (g *Gudang) Ubah(_ context.Context, _ *db.Tx, b models.Benefit) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ada := g.Baris[b.ID]; !ada {
		return repository.ErrTidakAda
	}
	g.Baris[b.ID] = b
	return nil
}
