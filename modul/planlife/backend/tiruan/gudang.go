// Package tiruan - Gudang Plan di memori untuk uji layanan dan handler. Data UJI buatan.
package tiruan

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"sync"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/planlife/backend/models"
	"nusantarare/modul/planlife/backend/repository"
)

// Transaksi tiruan - fn dijalankan langsung; galat = keadaan dikembalikan (rollback).
type Transaksi struct{ G *Gudang }

// Jalankan memenuhi services.Transaksi.
func (t Transaksi) Jalankan(_ context.Context, fn func(tx *db.Tx) error) error {
	t.G.mu.Lock()
	salinan := make(map[string]models.Plan, len(t.G.Baris))
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

// Business - satu baris master BUSINESS tiruan (dengan grupnya).
type Business struct {
	models.PilihanBusiness
	Grup string
}

// Gudang tiruan.
type Gudang struct {
	mu       sync.Mutex
	Baris    map[string]models.Plan
	Business []Business
	Benefit  []models.PilihanBenefit
	// Seq - nomor M_PRODUCT_TYPE_LIFE_SEQ berikutnya.
	Seq int
}

// Contoh - plan UJI 100001, 100002; business UJI grup 009 (dan satu grup lain); benefit UJI; sequence berikutnya 44.
func Contoh() *Gudang {
	g := &Gudang{Baris: map[string]models.Plan{}, Seq: 44,
		Business: []Business{
			{models.PilihanBusiness{ID: "9001", OldID: "L1", Note: "UJI KREDIT"}, "009"},
			{models.PilihanBusiness{ID: "9002", OldID: "L2", Note: "UJI Jiwa Kumpulan"}, "009"},
			{models.PilihanBusiness{ID: "9101", OldID: "F1", Note: "UJI KEBAKARAN"}, "001"},
		},
		Benefit: []models.PilihanBenefit{{ID: "100001", Benefit: "UJI RAWAT INAP"}, {ID: "100002", Benefit: "UJI MENINGGAL"}},
	}
	for _, p := range []models.Plan{
		{ID: "100001", CoverName: "UJI PLAN A", Business: "UJI KREDIT", BusinessID: "9001", Benefit: "UJI MENINGGAL", BenefitID: "100002"},
		{ID: "100002", CoverName: "UJI PLAN B", Business: "UJI Jiwa Kumpulan", BusinessID: "9002", Benefit: "UJI RAWAT INAP", BenefitID: "100001"},
	} {
		g.Baris[p.ID] = p
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

// Daftar memenuhi services.Gudang.
func (g *Gudang) Daftar(_ context.Context, s models.Saringan) ([]models.Plan, int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := make([]models.Plan, 0, len(g.Baris))
	for _, p := range g.Baris {
		out = append(out, p)
	}
	kurang := func(a, b models.Plan) bool {
		switch s.Urut {
		case models.UrutCoverName:
			return strings.ToUpper(a.CoverName) < strings.ToUpper(b.CoverName)
		case models.UrutBenefit:
			return strings.ToUpper(a.Benefit) < strings.ToUpper(b.Benefit)
		}
		return angka(a.ID) < angka(b.ID)
	}
	sort.SliceStable(out, func(a, b int) bool {
		if s.Turun {
			return kurang(out[b], out[a])
		}
		return kurang(out[a], out[b])
	})
	total := len(out)
	awal := min((s.Halaman-1)*models.UkuranHalaman, total)
	return out[awal:min(awal+models.UkuranHalaman, total)], total, nil
}

// Ambil memenuhi services.Gudang.
func (g *Gudang) Ambil(_ context.Context, _ *db.Tx, id string) (models.Plan, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	p, ok := g.Baris[id]
	if !ok {
		return models.Plan{}, repository.ErrTidakAda
	}
	return p, nil
}

// PemakaiNama memenuhi services.Gudang.
func (g *Gudang) PemakaiNama(_ context.Context, _ *db.Tx, nama, kecualiID string) ([]models.Plan, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var out []models.Plan
	for _, p := range g.Baris {
		if p.ID != kecualiID && models.KunciTeks(p.CoverName) == models.KunciTeks(nama) {
			out = append(out, p)
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
func (g *Gudang) Sisip(_ context.Context, _ *db.Tx, p models.Plan) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ada := g.Baris[p.ID]; ada {
		return repository.ErrKembar
	}
	g.Baris[p.ID] = p
	return nil
}

// Ubah memenuhi services.Gudang.
func (g *Gudang) Ubah(_ context.Context, _ *db.Tx, p models.Plan) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ada := g.Baris[p.ID]; !ada {
		return repository.ErrTidakAda
	}
	g.Baris[p.ID] = p
	return nil
}

// PilihanBusiness memenuhi services.Gudang - grup `grup`, urut ID.
func (g *Gudang) PilihanBusiness(_ context.Context, grup string) ([]models.PilihanBusiness, error) {
	out := []models.PilihanBusiness{}
	for _, b := range g.Business {
		if b.Grup == grup {
			out = append(out, b.PilihanBusiness)
		}
	}
	return out, nil
}

// PilihanBenefit memenuhi services.Gudang.
func (g *Gudang) PilihanBenefit(_ context.Context) ([]models.PilihanBenefit, error) {
	return append([]models.PilihanBenefit{}, g.Benefit...), nil
}
