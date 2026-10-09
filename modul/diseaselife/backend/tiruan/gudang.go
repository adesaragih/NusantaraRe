// Package tiruan - Gudang Disease Life di memori untuk uji layanan dan handler. Data UJI buatan (bentuk DEV: ID angka
// teks, ICD Code huruf besar, Disease huruf besar).
package tiruan

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"sync"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/diseaselife/backend/models"
	"nusantarare/modul/diseaselife/backend/repository"
)

// Transaksi tiruan - fn dijalankan langsung; galat = keadaan dikembalikan (rollback).
type Transaksi struct{ G *Gudang }

// Jalankan memenuhi services.Transaksi.
func (t Transaksi) Jalankan(_ context.Context, fn func(tx *db.Tx) error) error {
	t.G.mu.Lock()
	salinan := make(map[string]models.Penyakit, len(t.G.Baris))
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
	Baris map[string]models.Penyakit
	// Seq - nomor SEQ_DISEASE_LIFE berikutnya.
	Seq int
	// Diminta - saringan terakhir yang sampai ke Daftar (uji: saring dan halaman sampai ke "server").
	Diminta models.Saringan
}

// Contoh - bentuk DEV dengan nilai UJI: lima baris 100001-100005; sequence berikutnya 100006 (ID tertinggi + 1, D1.1).
func Contoh() *Gudang {
	g := &Gudang{Baris: map[string]models.Penyakit{}, Seq: 100006}
	for _, p := range []models.Penyakit{
		{ID: "100001", ICDCode: "UJI01", Disease: "UJI KOLERA"},
		{ID: "100002", ICDCode: "UJI02", Disease: "UJI DEMAM TIFOID"},
		{ID: "100003", ICDCode: "UJI03", Disease: "UJI DEMAM PARATIFOID"},
		{ID: "100004", ICDCode: "UJI04", Disease: "UJI INFEKSI SALMONELLA"},
		{ID: "100005", ICDCode: "UJI05", Disease: "UJI SHIGELOSIS"},
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

func memuat(nilai, kata string) bool {
	return kata == "" || strings.Contains(strings.ToUpper(nilai), strings.ToUpper(kata))
}

// Daftar memenuhi services.Gudang - saring "memuat" tanpa beda huruf, urut ID angka / ICD Code, 10 per halaman.
func (g *Gudang) Daftar(_ context.Context, s models.Saringan) ([]models.Penyakit, int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Diminta = s
	out := make([]models.Penyakit, 0, len(g.Baris))
	for _, p := range g.Baris {
		if memuat(p.ICDCode, s.ICDCode) && memuat(p.Disease, s.Disease) {
			out = append(out, p)
		}
	}
	sort.SliceStable(out, func(a, b int) bool {
		x, y := out[a], out[b]
		if !s.Naik {
			x, y = y, x
		}
		if s.Urut == models.UrutICD && x.ICDCode != y.ICDCode {
			return x.ICDCode < y.ICDCode
		}
		return angka(x.ID) < angka(y.ID)
	})
	total := len(out)
	awal := min((s.Halaman-1)*models.UkuranHalaman, total)
	return out[awal:min(awal+models.UkuranHalaman, total)], total, nil
}

// Ambil memenuhi services.Gudang.
func (g *Gudang) Ambil(_ context.Context, _ *db.Tx, id string) (models.Penyakit, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	p, ok := g.Baris[id]
	if !ok {
		return models.Penyakit{}, repository.ErrTidakAda
	}
	return p, nil
}

// PemakaiICD memenuhi services.Gudang.
func (g *Gudang) PemakaiICD(_ context.Context, _ *db.Tx, icd, kecualiID string) ([]models.Penyakit, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var out []models.Penyakit
	for _, p := range g.Baris {
		if p.ID != kecualiID && p.ICDCode != "" && models.KunciTeks(p.ICDCode) == models.KunciTeks(icd) {
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

// Sisip memenuhi services.Gudang - seperti PK_DISEASE_LIFE.
func (g *Gudang) Sisip(_ context.Context, _ *db.Tx, p models.Penyakit) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ada := g.Baris[p.ID]; ada {
		return repository.ErrKembar
	}
	g.Baris[p.ID] = p
	return nil
}

// Ubah memenuhi services.Gudang.
func (g *Gudang) Ubah(_ context.Context, _ *db.Tx, p models.Penyakit) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ada := g.Baris[p.ID]; !ada {
		return repository.ErrTidakAda
	}
	g.Baris[p.ID] = p
	return nil
}
