// Package tiruan - Gudang Treaty Description di memori untuk uji layanan dan handler. Data UJI buatan.
package tiruan

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatydescription/backend/models"
	"nusantarare/modul/treatydescription/backend/repository"
)

// Transaksi tiruan - fn dijalankan langsung.
func Transaksi(_ context.Context, fn func(tx *db.Tx) error) error { return fn(nil) }

// Gudang tiruan.
type Gudang struct {
	mu    sync.Mutex
	Baris map[string]models.Desc
	// Seq - nomor TREATY_DESCRIPTION_SEQ berikutnya; ID baru = "1" + 4 digit seperti PEGA_TREATYDESC.
	Seq int
}

// Contoh - lima baris: satu berstatus NULL (baris lama Pega), satu XOL nonaktif; sequence berikutnya 17 dan ID 10017
// sudah terpakai (uji lompat), jadi Add pertama menjadi 10018.
func Contoh() *Gudang {
	g := &Gudang{Baris: map[string]models.Desc{}, Seq: 17}
	for _, r := range []models.Desc{
		{ID: "10001", DescName: "UJI TREATY LIMIT", IsXOL: models.NonXOL},
		{ID: "10004", DescName: "UJI CASH LOSS LIMIT", IsXOL: models.NonXOL, StatusAktif: models.Aktif},
		{ID: "10013", DescName: "UJI EXCLUSION TREATY", IsXOL: models.NonXOL, StatusAktif: models.Aktif},
		{ID: "10015", DescName: "UJI XOL LAYER", IsXOL: models.XOL, StatusAktif: models.Nonaktif},
		{ID: "10017", DescName: "UJI MB CAPACITY", IsXOL: models.NonXOL, StatusAktif: models.Aktif},
	} {
		g.Baris[r.ID] = r
	}
	return g
}

func (g *Gudang) lengkap(r models.Desc) models.Desc {
	if r.IsXOL == "" {
		r.IsXOL = models.NonXOL
	}
	r.Aktif = r.StatusAktif != models.Nonaktif
	return r
}

// Daftar memenuhi services.Gudang - sama dengan SqlDaftar, urut ID.
func (g *Gudang) Daftar(_ context.Context, kata, xol, status string) ([]models.Desc, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	k := strings.ToUpper(strings.TrimSpace(kata))
	out := []models.Desc{}
	for _, r := range g.Baris {
		r = g.lengkap(r)
		if k != "" && !strings.Contains(strings.ToUpper(r.ID), k) && !strings.Contains(strings.ToUpper(r.DescName), k) {
			continue
		}
		if xol != "" && r.IsXOL != xol {
			continue
		}
		if status != "" && (status == models.Nonaktif) == r.Aktif {
			continue
		}
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

// Ambil memenuhi services.Gudang.
func (g *Gudang) Ambil(_ context.Context, _ *db.Tx, id string) (models.Desc, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	r, ok := g.Baris[id]
	if !ok {
		return models.Desc{}, repository.ErrTidakAda
	}
	return g.lengkap(r), nil
}

// PemakaiNama memenuhi services.Gudang.
func (g *Gudang) PemakaiNama(_ context.Context, _ *db.Tx, nama, kecualiID string) ([]string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var out []string
	for _, r := range g.Baris {
		if r.ID != kecualiID && strings.EqualFold(strings.TrimSpace(r.DescName), strings.TrimSpace(nama)) {
			out = append(out, r.ID)
		}
	}
	sort.Strings(out)
	return out, nil
}

// IDBaru memenuhi services.Gudang.
func (g *Gudang) IDBaru(context.Context, *db.Tx) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	id := fmt.Sprintf("%s%04d", repository.AwalanID, g.Seq)
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

// Sisip memenuhi services.Gudang - hanya keempat kolom tabel.
func (g *Gudang) Sisip(_ context.Context, _ *db.Tx, r models.Desc) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Baris[r.ID] = models.Desc{ID: r.ID, DescName: r.DescName, IsXOL: r.IsXOL, StatusAktif: r.StatusAktif}
	return nil
}

// Ubah memenuhi services.Gudang - ID tetap seperti SQL-nya.
func (g *Gudang) Ubah(_ context.Context, _ *db.Tx, r models.Desc) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.Baris[r.ID]; !ok {
		return repository.ErrTidakAda
	}
	g.Baris[r.ID] = models.Desc{ID: r.ID, DescName: r.DescName, IsXOL: r.IsXOL, StatusAktif: r.StatusAktif}
	return nil
}
