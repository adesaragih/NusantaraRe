// Package tiruan - Gudang Treaty Group OJK di memori untuk uji layanan dan handler. Data UJI buatan.
package tiruan

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"sync"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatygroupojk/backend/models"
	"nusantarare/modul/treatygroupojk/backend/repository"
)

// Transaksi tiruan - fn dijalankan langsung.
func Transaksi(_ context.Context, fn func(tx *db.Tx) error) error { return fn(nil) }

// Gudang tiruan.
type Gudang struct {
	mu    sync.Mutex
	Baris map[string]models.Ojk
	// Dikunci - berapa kali KunciNomorTertinggi dipanggil (uji: Add mengunci, Edit tidak).
	Dikunci int
}

// Contoh - tiga OJK; ID `09` tertinggi, Order No `1`, `2`, `11`.
func Contoh() *Gudang {
	g := &Gudang{Baris: map[string]models.Ojk{}}
	for _, o := range []models.Ojk{
		{ID: "01", Name: "UJI PROPERTY", NameIDN: "UJI HARTA BENDA", OrderNo: "1"},
		{ID: "09", Name: "UJI ENGINEERING", NameIDN: "UJI REKAYASA", OrderNo: "2"},
		{ID: "05", Name: "UJI AVIATION", NameIDN: "UJI PESAWAT", OrderNo: "11"},
	} {
		g.Baris[o.ID] = o
	}
	return g
}

func angka(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 1 << 30
	}
	return n
}

// Daftar memenuhi services.Gudang - urut Order No sebagai angka, lalu ID.
func (g *Gudang) Daftar(_ context.Context, kata string) ([]models.Ojk, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	k := strings.ToUpper(strings.TrimSpace(kata))
	out := []models.Ojk{}
	for _, o := range g.Baris {
		if k != "" && !strings.Contains(strings.ToUpper(o.ID+" "+o.Name+" "+o.NameIDN), k) {
			continue
		}
		out = append(out, o)
	}
	sort.Slice(out, func(i, j int) bool {
		if a, b := angka(out[i].OrderNo), angka(out[j].OrderNo); a != b {
			return a < b
		}
		return out[i].ID < out[j].ID
	})
	return out, nil
}

// Ambil memenuhi services.Gudang.
func (g *Gudang) Ambil(_ context.Context, _ *db.Tx, id string) (models.Ojk, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	o, ok := g.Baris[id]
	if !ok {
		return models.Ojk{}, repository.ErrTidakAda
	}
	return o, nil
}

func (g *Gudang) pemakai(nilai, kecualiID string, kolom func(models.Ojk) string) []models.Ojk {
	g.mu.Lock()
	defer g.mu.Unlock()
	var out []models.Ojk
	for _, o := range g.Baris {
		if o.ID != kecualiID && strings.EqualFold(strings.TrimSpace(kolom(o)), strings.TrimSpace(nilai)) {
			out = append(out, o)
		}
	}
	return out
}

// PemakaiNama memenuhi services.Gudang.
func (g *Gudang) PemakaiNama(_ context.Context, _ *db.Tx, nama, kecualiID string) ([]models.Ojk, error) {
	return g.pemakai(nama, kecualiID, func(o models.Ojk) string { return o.Name }), nil
}

// KunciNomorTertinggi memenuhi services.Gudang.
func (g *Gudang) KunciNomorTertinggi(context.Context, *db.Tx) (id, order int, err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Dikunci++
	for k, o := range g.Baris {
		if v, err := strconv.Atoi(k); err == nil && v > id {
			id = v
		}
		if v, err := strconv.Atoi(o.OrderNo); err == nil && v > order {
			order = v
		}
	}
	return id, order, nil
}

// Sisip memenuhi services.Gudang.
func (g *Gudang) Sisip(_ context.Context, _ *db.Tx, o models.Ojk) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Baris[o.ID] = o
	return nil
}

// Ubah memenuhi services.Gudang.
func (g *Gudang) Ubah(_ context.Context, _ *db.Tx, o models.Ojk) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	lama, ok := g.Baris[o.ID]
	if !ok {
		return repository.ErrTidakAda
	}
	lama.Name, lama.NameIDN = o.Name, o.NameIDN
	g.Baris[o.ID] = lama
	return nil
}
