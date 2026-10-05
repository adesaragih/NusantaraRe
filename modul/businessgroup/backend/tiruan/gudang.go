// Package tiruan - Gudang Business Group di memori untuk uji layanan dan handler. Data UJI buatan.
package tiruan

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/businessgroup/backend/models"
	"nusantarare/modul/businessgroup/backend/repository"
)

// Transaksi tiruan - fn dijalankan langsung.
func Transaksi(_ context.Context, fn func(tx *db.Tx) error) error { return fn(nil) }

// Gudang tiruan.
type Gudang struct {
	mu     sync.Mutex
	Baris  map[string]models.BisnisGrup
	Treaty []models.TreatyGroup
	// Seq - nomor sequence berikutnya; ID baru = "1" + 4 digit seperti situs aktif DEV.
	Seq int
}

// Contoh - empat grup bisnis (satu berakhiran SYARIAH, satu ber-TOPID yatim `10019`), dua Treaty Group. Sequence
// berikutnya 13 - ID 10013 sudah terpakai (uji lompat).
func Contoh() *Gudang {
	g := &Gudang{Baris: map[string]models.BisnisGrup{}, Seq: 13}
	g.Treaty = []models.TreatyGroup{{ID: "10007", Name: "UJI PROPERTY"}, {ID: "10009", Name: "UJI MARINE CARGO"}}
	for _, b := range []models.BisnisGrup{
		{ID: "10013", Name: "UJI FIRE", Alias: "UJI FIRE", TopID: "10007", TreatyName: "UJI PROPERTY"},
		{ID: "10014", Name: "UJI EARTHQUAKE", Alias: "UJI EQ", TopID: "10007", TreatyName: "UJI PROPERTY LAMA"},
		{ID: "10020", Name: "UJI FIRE SYARIAH", Alias: "UJI FIRE SYARIAH", TopID: "10007", TreatyName: "UJI PROPERTY"},
		{ID: "10028", Name: "UJI YATIM", Alias: "UJI YATIM", TopID: "10019", TreatyName: "UJI HILANG"},
	} {
		g.Baris[b.ID] = b
	}
	return g
}

func syariah(nama string) bool {
	return strings.HasSuffix(strings.ToUpper(strings.TrimSpace(nama)), models.AkhiranSyariah)
}

// Daftar memenuhi services.Gudang - tanpa SYARIAH, urut nama lalu ID.
func (g *Gudang) Daftar(_ context.Context, kata, topID string) ([]models.BisnisGrup, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	k := strings.ToUpper(strings.TrimSpace(kata))
	out := []models.BisnisGrup{}
	for _, b := range g.Baris {
		if syariah(b.Name) || (topID != "" && b.TopID != topID) {
			continue
		}
		if k != "" && !strings.Contains(strings.ToUpper(b.ID+" "+b.Name+" "+b.Alias+" "+b.TreatyName), k) {
			continue
		}
		out = append(out, b)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Name != out[j].Name {
			return out[i].Name < out[j].Name
		}
		return out[i].ID < out[j].ID
	})
	return out, nil
}

// Ambil memenuhi services.Gudang (juga yang SYARIAH, seperti SQL-nya).
func (g *Gudang) Ambil(_ context.Context, _ *db.Tx, id string) (models.BisnisGrup, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	b, ok := g.Baris[id]
	if !ok {
		return models.BisnisGrup{}, repository.ErrTidakAda
	}
	return b, nil
}

// PemakaiNama memenuhi services.Gudang.
func (g *Gudang) PemakaiNama(_ context.Context, _ *db.Tx, nama, kecualiID string) ([]string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var out []string
	for _, b := range g.Baris {
		if b.ID != kecualiID && strings.EqualFold(strings.TrimSpace(b.Name), strings.TrimSpace(nama)) {
			out = append(out, b.ID)
		}
	}
	return out, nil
}

// DaftarTreaty memenuhi services.Gudang.
func (g *Gudang) DaftarTreaty(context.Context) ([]models.TreatyGroup, error) {
	return append([]models.TreatyGroup{}, g.Treaty...), nil
}

// AmbilTreaty memenuhi services.Gudang.
func (g *Gudang) AmbilTreaty(_ context.Context, _ *db.Tx, id string) (models.TreatyGroup, error) {
	for _, t := range g.Treaty {
		if t.ID == id {
			return t, nil
		}
	}
	return models.TreatyGroup{}, repository.ErrTidakAda
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
func (g *Gudang) Sisip(_ context.Context, _ *db.Tx, b models.BisnisGrup) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Baris[b.ID] = b
	return nil
}

// Ubah memenuhi services.Gudang.
func (g *Gudang) Ubah(_ context.Context, _ *db.Tx, b models.BisnisGrup) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.Baris[b.ID]; !ok {
		return repository.ErrTidakAda
	}
	g.Baris[b.ID] = b
	return nil
}
