// Package tiruan - Gudang Adjuster Consultant di memori untuk uji layanan dan handler. Data UJI-* buatan.
package tiruan

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/adjusterconsultant/backend/models"
	"nusantarare/modul/adjusterconsultant/backend/repository"
)

// Transaksi tiruan - fn dijalankan langsung.
func Transaksi(_ context.Context, fn func(tx *db.Tx) error) error { return fn(nil) }

// Gudang tiruan.
type Gudang struct {
	mu    sync.Mutex
	Baris map[string]models.Adjuster
	// Seq - nomor sequence berikutnya; ID baru = "1" + 4 digit seperti situs aktif DEV.
	Seq int
	// Jam - EDITDATE yang ditulis.
	Jam string
}

// Contoh - dua baris aktif dan satu nonaktif; sequence berikutnya 2 (ID 10002 sudah terpakai: uji lompat).
func Contoh() *Gudang {
	g := &Gudang{Baris: map[string]models.Adjuster{}, Seq: 2, Jam: "05-10-2026 09:30"}
	for _, a := range []models.Adjuster{
		{ID: "10001", Name: "UJI ADJUSTER SATU", Address: "UJI ALAMAT", TelpNo: "+62 21 000", Active: true},
		{ID: "10002", Name: "UJI ADJUSTER DUA", Active: true},
		{ID: "10003", Name: "UJI KONSULTAN LAMA", Active: false},
	} {
		g.Baris[a.ID] = a
	}
	return g
}

// Daftar memenuhi services.Gudang.
func (g *Gudang) Daftar(_ context.Context, kata, bendera string) ([]models.Adjuster, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	k := strings.ToUpper(strings.TrimSpace(kata))
	out := []models.Adjuster{}
	for _, a := range g.Baris {
		if bendera != "" && (bendera == models.BenderaAktif) != a.Active {
			continue
		}
		teks := strings.ToUpper(a.ID + " " + a.Name + " " + a.Address + " " + a.TelpNo)
		if k != "" && !strings.Contains(teks, k) {
			continue
		}
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// Ambil memenuhi services.Gudang.
func (g *Gudang) Ambil(_ context.Context, _ *db.Tx, id string) (models.Adjuster, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	a, ok := g.Baris[id]
	if !ok {
		return models.Adjuster{}, repository.ErrTidakAda
	}
	return a, nil
}

// PemakaiNama memenuhi services.Gudang.
func (g *Gudang) PemakaiNama(_ context.Context, _ *db.Tx, nama, kecualiID string) ([]models.Adjuster, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var out []models.Adjuster
	for _, a := range g.Baris {
		if a.ID != kecualiID && strings.EqualFold(strings.TrimSpace(a.Name), strings.TrimSpace(nama)) {
			out = append(out, a)
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
func (g *Gudang) Sisip(_ context.Context, _ *db.Tx, a models.Adjuster, username string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ada := g.Baris[a.ID]; ada {
		return repository.ErrIDTerpakai
	}
	a.Username, a.EditDate, a.Active = username, g.Jam, true
	g.Baris[a.ID] = a
	return nil
}

// Ubah memenuhi services.Gudang.
func (g *Gudang) Ubah(_ context.Context, _ *db.Tx, a models.Adjuster, username string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	lama, ok := g.Baris[a.ID]
	if !ok {
		return repository.ErrTidakAda
	}
	lama.Name, lama.Address, lama.TelpNo, lama.Username, lama.EditDate = a.Name, a.Address, a.TelpNo, username, g.Jam
	g.Baris[a.ID] = lama
	return nil
}

// SetelAktif memenuhi services.Gudang.
func (g *Gudang) SetelAktif(_ context.Context, _ *db.Tx, id string, aktif bool, username string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	a, ok := g.Baris[id]
	if !ok {
		return repository.ErrTidakAda
	}
	a.Active, a.Username, a.EditDate = aktif, username, g.Jam
	g.Baris[id] = a
	return nil
}
