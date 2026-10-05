// Package tiruan - Gudang Reinsurance Type di memori untuk uji layanan dan handler. Data UJI buatan.
package tiruan

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/reinsurancetype/backend/models"
	"nusantarare/modul/reinsurancetype/backend/repository"
)

// Transaksi tiruan - fn dijalankan langsung.
func Transaksi(_ context.Context, fn func(tx *db.Tx) error) error { return fn(nil) }

// Gudang tiruan.
type Gudang struct {
	mu    sync.Mutex
	Baris map[string]models.Jenis
	// Seq - nomor sequence berikutnya; ID baru = "1" + 4 digit.
	Seq int
}

// Contoh - lima baris: Flag `active`, `inactive`, `1` (Contract Retro Life), dan kosong; satu Type kosong. Sequence
// berikutnya 262 - ID 10262 sudah terpakai (uji lompat).
func Contoh() *Gudang {
	g := &Gudang{Baris: map[string]models.Jenis{}, Seq: 262}
	for _, j := range []models.Jenis{
		{ID: "10007", Name: "UJI ORS", Type: "1", Code: "00", Flag: "active", GroupType: "OR"},
		{ID: "10196", Name: "UJI QS", Type: "2", Code: "00", Flag: "1"},
		{ID: "10030", Name: "UJI PURE FACULTATIVE", Type: "4", Code: "30", Flag: "inactive", TglUpdate: "20191029T164018.868 GMT"},
		{ID: "10033", Name: "UJI POOL SURPLUS", Code: "00"},
		{ID: "10262", Name: "UJI 2026 QS TRT", Type: "2", SoaName: "UJI QUOTA SHARE", Code: "00", Flag: "active", NoUrut: "1"},
	} {
		g.Baris[j.ID] = j
	}
	return g
}

// Daftar memenuhi services.Gudang - urut Type, nama, ID.
func (g *Gudang) Daftar(_ context.Context, kata, tipe, flag string) ([]models.Jenis, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	k := strings.ToUpper(strings.TrimSpace(kata))
	out := []models.Jenis{}
	for _, j := range g.Baris {
		if (tipe != "" && j.Type != tipe) || (flag != "" && j.Flag != flag) {
			continue
		}
		if k != "" && !strings.Contains(strings.ToUpper(j.ID+" "+j.Name+" "+j.SoaName+" "+j.Code), k) {
			continue
		}
		out = append(out, j)
	}
	sort.Slice(out, func(a, b int) bool {
		if out[a].Type != out[b].Type {
			return out[a].Type < out[b].Type
		}
		if out[a].Name != out[b].Name {
			return out[a].Name < out[b].Name
		}
		return out[a].ID < out[b].ID
	})
	return out, nil
}

// Ambil memenuhi services.Gudang.
func (g *Gudang) Ambil(_ context.Context, _ *db.Tx, id string) (models.Jenis, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	j, ok := g.Baris[id]
	if !ok {
		return models.Jenis{}, repository.ErrTidakAda
	}
	return j, nil
}

// PemakaiNama memenuhi services.Gudang.
func (g *Gudang) PemakaiNama(_ context.Context, _ *db.Tx, nama, kecualiID string) ([]models.Jenis, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var out []models.Jenis
	for _, j := range g.Baris {
		if j.ID != kecualiID && strings.EqualFold(strings.TrimSpace(j.Name), strings.TrimSpace(nama)) {
			out = append(out, j)
		}
	}
	return out, nil
}

// IDBaru memenuhi services.Gudang.
func (g *Gudang) IDBaru(context.Context, *db.Tx) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	id := fmt.Sprintf("1%04d", g.Seq)
	g.Seq++
	return id, nil
}

// AdaID memenuhi services.Gudang.
func (g *Gudang) AdaID(_ context.Context, _ *db.Tx, id string) (bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	_, ok := g.Baris[id]
	return ok, nil
}

// Sisip memenuhi services.Gudang.
func (g *Gudang) Sisip(_ context.Context, _ *db.Tx, j models.Jenis) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Baris[j.ID] = j
	return nil
}

// Ubah memenuhi services.Gudang.
func (g *Gudang) Ubah(_ context.Context, _ *db.Tx, j models.Jenis) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.Baris[j.ID]; !ok {
		return repository.ErrTidakAda
	}
	g.Baris[j.ID] = j
	return nil
}
