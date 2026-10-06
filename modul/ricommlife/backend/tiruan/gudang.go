// Package tiruan - Gudang R/I Comm Life di memori untuk uji layanan dan handler. Data UJI buatan.
package tiruan

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/ricommlife/backend/models"
	"nusantarare/modul/ricommlife/backend/repository"
)

// Transaksi tiruan - fn dijalankan langsung; galat = keadaan dikembalikan (rollback).
type Transaksi struct{ G *Gudang }

// Jalankan memenuhi services.Transaksi: salinan keadaan dipulihkan bila fn gagal.
func (t Transaksi) Jalankan(_ context.Context, fn func(tx *db.Tx) error) error {
	t.G.mu.Lock()
	r, k := salin(t.G.Ringkasan), salin(t.G.Komisi)
	t.G.mu.Unlock()
	if err := fn(nil); err != nil {
		t.G.mu.Lock()
		t.G.Ringkasan, t.G.Komisi = r, k
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
	Komisi    map[string]models.Komisi
	// Situs - baris `M_SITE_DATABASE` CURRENT_SITE = '1' (kosong = tidak ada; lebih dari satu = galat).
	SitusAktif []string
	// SeqRingkasan / SeqKomisi - nomor sequence berikutnya.
	SeqRingkasan, SeqKomisi int
}

// Contoh - situs 1; ringkasan 1000003 (seperti satu baris DEV, nama UJI) dan 1000004; rincian 1000040-1000042.
// Sequence ringkasan berikutnya 4 (ID 1000004 terpakai: diuji lompat lewat AdaID); sequence rincian berikutnya 44.
func Contoh() *Gudang {
	g := &Gudang{Ringkasan: map[string]models.Ringkasan{}, Komisi: map[string]models.Komisi{}, SitusAktif: []string{"1"},
		SeqRingkasan: 4, SeqKomisi: 44}
	for _, r := range []models.Ringkasan{
		{ID: "1000003", UsedBy: "UJI COMM RETRO", OperatorID: "UJI-LAMA", ModifiedDate: "20181205T073755.559 GMT"},
		{ID: "1000004", UsedBy: "UJI COMM B", OperatorID: "UJI-LAMA", ModifiedDate: "20240103"},
	} {
		g.Ringkasan[r.ID] = r
	}
	for _, k := range []models.Komisi{
		{ID: "1000040", IDUsedBy: "1000003", UsedBy: "UJI COMM RETRO", Contract: "1", Year: "1", Comm: "12.5"},
		{ID: "1000041", IDUsedBy: "1000003", UsedBy: "UJI COMM RETRO", Contract: "1", Year: "2", Comm: "10"},
		{ID: "1000042", IDUsedBy: "1000004", UsedBy: "UJI COMM B", Contract: "2", Year: "1", Comm: "0.5"},
	} {
		g.Komisi[k.ID] = k
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
	awal := min((s.Halaman-1)*models.UkuranHalaman, total)
	return out[awal:min(awal+models.UkuranHalaman, total)], total, nil
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

// Situs memenuhi services.Gudang - seperti repository: tepat satu baris.
func (g *Gudang) Situs(_ context.Context, _ *db.Tx) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if len(g.SitusAktif) != 1 {
		return "", fmt.Errorf("%w (%d baris)", repository.ErrSitus, len(g.SitusAktif))
	}
	return g.SitusAktif[0], nil
}

// NomorBaru memenuhi services.Gudang.
func (g *Gudang) NomorBaru(_ context.Context, _ *db.Tx, ringkasan bool) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if ringkasan {
		g.SeqRingkasan++
		return fmt.Sprint(g.SeqRingkasan - 1), nil
	}
	g.SeqKomisi++
	return fmt.Sprint(g.SeqKomisi - 1), nil
}

// AdaID memenuhi services.Gudang.
func (g *Gudang) AdaID(_ context.Context, _ *db.Tx, ringkasan bool, id string) (bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if ringkasan {
		_, ok := g.Ringkasan[id]
		return ok, nil
	}
	_, ok := g.Komisi[id]
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

// UbahNamaKomisi memenuhi services.Gudang.
func (g *Gudang) UbahNamaKomisi(_ context.Context, _ *db.Tx, idUsedBy, nama string) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	n := 0
	for id, k := range g.Komisi {
		if k.IDUsedBy == idUsedBy {
			k.UsedBy = nama
			g.Komisi[id] = k
			n++
		}
	}
	return n, nil
}

// HapusKomisi memenuhi services.Gudang.
func (g *Gudang) HapusKomisi(_ context.Context, _ *db.Tx, idUsedBy string) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	n := 0
	for id, k := range g.Komisi {
		if k.IDUsedBy == idUsedBy {
			delete(g.Komisi, id)
			n++
		}
	}
	return n, nil
}

func (g *Gudang) milik(idUsedBy string) []models.Komisi {
	var out []models.Komisi
	for _, k := range g.Komisi {
		if k.IDUsedBy == idUsedBy {
			out = append(out, k)
		}
	}
	// Seperti SqlDaftarKomisi: ID angka menaik.
	sort.Slice(out, func(a, b int) bool { return angka(out[a].ID) < angka(out[b].ID) })
	return out
}

// JumlahKomisi memenuhi services.Gudang.
func (g *Gudang) JumlahKomisi(_ context.Context, _ *db.Tx, idUsedBy string) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return len(g.milik(idUsedBy)), nil
}

// DaftarKomisi memenuhi services.Gudang.
func (g *Gudang) DaftarKomisi(_ context.Context, idUsedBy string, halaman int) ([]models.Komisi, int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	semua := g.milik(idUsedBy)
	awal := min((halaman-1)*models.UkuranHalaman, len(semua))
	return semua[awal:min(awal+models.UkuranHalaman, len(semua))], len(semua), nil
}

// KomisiDari memenuhi services.Gudang.
func (g *Gudang) KomisiDari(_ context.Context, _ *db.Tx, ids []string) ([]models.Komisi, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var out []models.Komisi
	for _, id := range ids {
		out = append(out, g.milik(id)...)
	}
	return out, nil
}

// SisipKomisi memenuhi services.Gudang - PK ID seperti tabel flat.
func (g *Gudang) SisipKomisi(_ context.Context, _ *db.Tx, k models.Komisi) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ada := g.Komisi[k.ID]; ada {
		return repository.ErrKembar
	}
	g.Komisi[k.ID] = k
	return nil
}

// UbahKomisi memenuhi services.Gudang - CONTRACT, YEAR, COMM baris milik k.IDUsedBy.
func (g *Gudang) UbahKomisi(_ context.Context, _ *db.Tx, k models.Komisi) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	lama, ada := g.Komisi[k.ID]
	if !ada || lama.IDUsedBy != k.IDUsedBy {
		return repository.ErrTidakAda
	}
	lama.Contract, lama.Year, lama.Comm = k.Contract, k.Year, k.Comm
	g.Komisi[k.ID] = lama
	return nil
}
