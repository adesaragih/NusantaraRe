// Package tiruan - Gudang R/I Rate Life di memori untuk uji layanan dan handler. Data UJI buatan.
package tiruan

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/riratelife/backend/models"
	"nusantarare/modul/riratelife/backend/repository"
)

// Transaksi tiruan - fn dijalankan langsung; galat = keadaan dikembalikan (rollback).
type Transaksi struct{ G *Gudang }

// Jalankan memenuhi services.Transaksi: salinan keadaan dipulihkan bila fn gagal.
func (t Transaksi) Jalankan(_ context.Context, fn func(tx *db.Tx) error) error {
	t.G.mu.Lock()
	r, b := salin(t.G.Ringkasan), salin(t.G.Rate)
	t.G.mu.Unlock()
	if err := fn(nil); err != nil {
		t.G.mu.Lock()
		t.G.Ringkasan, t.G.Rate = r, b
		t.G.mu.Unlock()
		return err
	}
	return nil
}

func salin[T any](m map[string]T) map[string]T {
	out := make(map[string]T, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// Gudang tiruan.
type Gudang struct {
	mu        sync.Mutex
	Ringkasan map[string]models.Ringkasan
	Rate      map[string]models.Rate
	// SeqRingkasan / SeqRate - nomor sequence berikutnya.
	SeqRingkasan, SeqRate int
}

// Contoh - tiga ringkasan (ID 101, 102, 105) dan empat rate. Sequence ringkasan berikutnya 105 (terpakai:
// diuji lompat lewat MaksID); sequence rate berikutnya 9001 (ID 9002 terpakai: diuji lompat lewat AdaID - ID 9002 di
// bawah ID angka tertinggi 9100).
func Contoh() *Gudang {
	g := &Gudang{Ringkasan: map[string]models.Ringkasan{}, Rate: map[string]models.Rate{}, SeqRingkasan: 105, SeqRate: 9001}
	for _, r := range []models.Ringkasan{
		{ID: "101", UsedBy: "UJI RATE A", OperatorID: "UJI-LAMA", ModifiedDate: "20240102T030405.000 GMT"},
		{ID: "102", UsedBy: "UJI RATE B", OperatorID: "UJI-LAMA", ModifiedDate: "20240103"},
		{ID: "105", UsedBy: "UJI RATE C"},
	} {
		g.Ringkasan[r.ID] = r
	}
	for _, r := range []models.Rate{
		{ID: "9000", IDUsedBy: "101", UsedBy: "UJI RATE A", Gender: "U", Contract: "1", Age: "30", Rate: "0,5"},
		{ID: "9002", IDUsedBy: "101", UsedBy: "UJI RATE A", Gender: "U", Contract: "", Age: "05", Rate: "1,25"},
		{ID: "9003", IDUsedBy: "101", UsedBy: "UJI RATE A", Gender: "M", Contract: "1", Age: "30", Rate: "0.75"},
		{ID: "9100", IDUsedBy: "102", UsedBy: "UJI RATE B", Gender: "F", Contract: "2", Age: "40", Rate: "2"},
	} {
		g.Rate[r.ID] = r
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

// Daftar memenuhi services.Gudang - saring memuat, urut sesuai Saringan, satu halaman.
func (g *Gudang) Daftar(_ context.Context, s models.Saringan) ([]models.Ringkasan, int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var out []models.Ringkasan
	for _, r := range g.Ringkasan {
		if s.ID != "" && !strings.Contains(strings.ToUpper(r.ID), strings.ToUpper(s.ID)) {
			continue
		}
		if s.UsedBy != "" && !strings.Contains(strings.ToUpper(r.UsedBy), strings.ToUpper(s.UsedBy)) {
			continue
		}
		out = append(out, r)
	}
	kurang := func(a, b models.Ringkasan) bool {
		switch s.Urut {
		case models.UrutUsedBy:
			return strings.ToUpper(a.UsedBy) < strings.ToUpper(b.UsedBy)
		case models.UrutOperator:
			return strings.ToUpper(a.OperatorID) < strings.ToUpper(b.OperatorID)
		case models.UrutTanggal:
			return a.ModifiedDate < b.ModifiedDate
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
	awal := (s.Halaman - 1) * models.UkuranHalaman
	if awal > total {
		awal = total
	}
	akhir := min(awal+models.UkuranHalaman, total)
	return out[awal:akhir], total, nil
}

// Ambil memenuhi services.Gudang.
func (g *Gudang) Ambil(_ context.Context, _ *db.Tx, id string) (models.Ringkasan, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	r, ok := g.Ringkasan[id]
	if !ok {
		return models.Ringkasan{}, repository.ErrTidakAda
	}
	return r, nil
}

// PemakaiNama memenuhi services.Gudang - urut ID.
func (g *Gudang) PemakaiNama(_ context.Context, _ *db.Tx, nama, kecualiID string) ([]models.Ringkasan, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var out []models.Ringkasan
	for _, r := range g.Ringkasan {
		if r.ID != kecualiID && strings.EqualFold(strings.TrimSpace(r.UsedBy), strings.TrimSpace(nama)) {
			out = append(out, r)
		}
	}
	sort.Slice(out, func(a, b int) bool { return out[a].ID < out[b].ID })
	return out, nil
}

// IDBaru memenuhi services.Gudang.
func (g *Gudang) IDBaru(_ context.Context, _ *db.Tx, ringkasan bool) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if ringkasan {
		g.SeqRingkasan++
		return fmt.Sprint(g.SeqRingkasan - 1), nil
	}
	g.SeqRate++
	return fmt.Sprint(g.SeqRate - 1), nil
}

// MaksID memenuhi services.Gudang.
func (g *Gudang) MaksID(_ context.Context, _ *db.Tx, ringkasan bool) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	maks := 0
	if ringkasan {
		for id := range g.Ringkasan {
			maks = max(maks, angka(id))
		}
		return maks, nil
	}
	for id := range g.Rate {
		maks = max(maks, angka(id))
	}
	return maks, nil
}

// AdaID memenuhi services.Gudang.
func (g *Gudang) AdaID(_ context.Context, _ *db.Tx, ringkasan bool, id string) (bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if ringkasan {
		_, ok := g.Ringkasan[id]
		return ok, nil
	}
	_, ok := g.Rate[id]
	return ok, nil
}

// SisipRingkasan memenuhi services.Gudang.
func (g *Gudang) SisipRingkasan(_ context.Context, _ *db.Tx, r models.Ringkasan) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ada := g.Ringkasan[r.ID]; ada {
		return repository.ErrKembar
	}
	g.Ringkasan[r.ID] = r
	return nil
}

// UbahRingkasan memenuhi services.Gudang.
func (g *Gudang) UbahRingkasan(_ context.Context, _ *db.Tx, r models.Ringkasan) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.Ringkasan[r.ID]; !ok {
		return repository.ErrTidakAda
	}
	g.Ringkasan[r.ID] = r
	return nil
}

// UbahNamaRate memenuhi services.Gudang.
func (g *Gudang) UbahNamaRate(_ context.Context, _ *db.Tx, idUsedBy, nama string) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	n := 0
	for id, r := range g.Rate {
		if r.IDUsedBy == idUsedBy {
			r.UsedBy = nama
			g.Rate[id] = r
			n++
		}
	}
	return n, nil
}

// HapusRingkasan memenuhi services.Gudang.
func (g *Gudang) HapusRingkasan(_ context.Context, _ *db.Tx, id string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.Ringkasan[id]; !ok {
		return repository.ErrTidakAda
	}
	delete(g.Ringkasan, id)
	return nil
}

// HapusRate memenuhi services.Gudang.
func (g *Gudang) HapusRate(_ context.Context, _ *db.Tx, idUsedBy string) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	n := 0
	for id, r := range g.Rate {
		if r.IDUsedBy == idUsedBy {
			delete(g.Rate, id)
			n++
		}
	}
	return n, nil
}

func (g *Gudang) rateMilik(idUsedBy string) []models.Rate {
	var out []models.Rate
	for _, r := range g.Rate {
		if r.IDUsedBy == idUsedBy {
			out = append(out, r)
		}
	}
	sort.Slice(out, func(a, b int) bool {
		x, y := out[a], out[b]
		if x.Gender != y.Gender {
			return x.Gender < y.Gender
		}
		if angka(x.Contract) != angka(y.Contract) {
			return angka(x.Contract) < angka(y.Contract)
		}
		if angka(x.Age) != angka(y.Age) {
			return angka(x.Age) < angka(y.Age)
		}
		return x.ID < y.ID
	})
	return out
}

// JumlahRate memenuhi services.Gudang.
func (g *Gudang) JumlahRate(_ context.Context, _ *db.Tx, idUsedBy string) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return len(g.rateMilik(idUsedBy)), nil
}

// DaftarRate memenuhi services.Gudang.
func (g *Gudang) DaftarRate(_ context.Context, idUsedBy string, halaman int) ([]models.Rate, int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	semua := g.rateMilik(idUsedBy)
	awal := min((halaman-1)*models.UkuranHalaman, len(semua))
	return semua[awal:min(awal+models.UkuranHalaman, len(semua))], len(semua), nil
}

// RateDari memenuhi services.Gudang.
func (g *Gudang) RateDari(_ context.Context, _ *db.Tx, ids []string) ([]models.Rate, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var out []models.Rate
	for _, id := range ids {
		out = append(out, g.rateMilik(id)...)
	}
	return out, nil
}

// SisipRate memenuhi services.Gudang.
func (g *Gudang) SisipRate(_ context.Context, _ *db.Tx, r models.Rate) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ada := g.Rate[r.ID]; ada {
		return repository.ErrKembar
	}
	g.Rate[r.ID] = r
	return nil
}
