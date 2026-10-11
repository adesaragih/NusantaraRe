// Package tiruan - Gudang Endorsment Fac In di memori untuk uji layanan dan handler. Data UJI buatan: nomor polis
// sintetis, tanpa nilai produksi.
package tiruan

import (
	"context"
	"strconv"
	"sync"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/endorsmentfacin/backend/models"
	"nusantarare/modul/endorsmentfacin/backend/repository"
)

// Gudang tiruan: setiap medan = jawaban query; Galat* = galat yang dikembalikan query itu. Panggilan dicatat.
type Gudang struct {
	mu sync.Mutex

	StatusEDMNilai   string
	JenisBisnisNilai string
	KlaimNilai       []string
	Pembayaran       bool
	Renewal          bool
	Facout           bool
	Periode          *repository.PeriodeTeks
	PengecualianPer  map[string][]string
	// Proteksi - klep OPENPROTEKSI_EDM per TYPE (1-4).
	Proteksi map[int]bool
	// GalatProteksi - galat AdaProteksi.
	GalatProteksi error

	GalatStatusEDM, GalatJenisBisnis, GalatKlaim, GalatPembayaran error
	GalatRenewal, GalatFacout, GalatPeriode, GalatPengecualian    error

	// EDMTerbuka - jawaban gerbang 2 (ID kasus EDM terbuka; "" = tidak ada); GalatEDMBelumSelesai - galatnya.
	EDMTerbuka           string
	GalatEDMBelumSelesai error

	// Tulisan I3. Seq - nomor sequence terakhir (NEXTVAL = Seq+1; seperti Oracle, TIDAK ikut rollback).
	Seq       int64
	WorkPolis map[string]string // ID → akun pembuat
	Kasus     map[string]models.KepalaKasus
	// Galat tulisan.
	GalatNomor, GalatSisipWorkPolis, GalatSisipKasus error

	// Panggilan - nama metode menurut urutan dipanggil.
	Panggilan []string
}

// Bersih - polis FIRE tanpa alasan tolak, periode 01-01-2025 s/d 01-01-2026 pukul 12:00 WIB (bentuk tulisan NB).
func Bersih() *Gudang {
	return &Gudang{JenisBisnisNilai: "Fire",
		Periode: &repository.PeriodeTeks{Mulai: "20250101120000", Akhir: "20260101120000"}}
}

// Transaksi - fn dijalankan langsung; galat → WorkPolis dan Kasus kembali ke keadaan sebelum fn (rollback).
func (g *Gudang) Transaksi(_ context.Context, fn func(tx *db.Tx) error) error {
	g.mu.Lock()
	wp, ks := map[string]string{}, map[string]models.KepalaKasus{}
	for k, v := range g.WorkPolis {
		wp[k] = v
	}
	for k, v := range g.Kasus {
		ks[k] = v
	}
	g.mu.Unlock()
	if err := fn(nil); err != nil {
		g.mu.Lock()
		g.WorkPolis, g.Kasus = wp, ks
		g.mu.Unlock()
		return err
	}
	return nil
}

func (g *Gudang) EDMBelumSelesai(context.Context, string) (string, bool, error) {
	g.catat("EDMBelumSelesai")
	return g.EDMTerbuka, g.EDMTerbuka != "", g.GalatEDMBelumSelesai
}

func (g *Gudang) NomorKasusBerikut(context.Context, *db.Tx) (string, error) {
	g.catat("NomorKasusBerikut")
	if g.GalatNomor != nil {
		return "", g.GalatNomor
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Seq++
	return models.RakitIDKasusEDM(g.Seq), nil
}

func (g *Gudang) SisipWorkPolis(_ context.Context, _ *db.Tx, id, akun string) error {
	g.catat("SisipWorkPolis")
	if g.GalatSisipWorkPolis != nil {
		return g.GalatSisipWorkPolis
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.WorkPolis == nil {
		g.WorkPolis = map[string]string{}
	}
	g.WorkPolis[id] = akun
	return nil
}

func (g *Gudang) SisipKasus(_ context.Context, _ *db.Tx, k models.KepalaKasus) error {
	g.catat("SisipKasus")
	if g.GalatSisipKasus != nil {
		return g.GalatSisipKasus
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Kasus == nil {
		g.Kasus = map[string]models.KepalaKasus{}
	}
	g.Kasus[k.ID] = k
	return nil
}

func (g *Gudang) catat(nama string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Panggilan = append(g.Panggilan, nama)
}

// Dipanggil - apakah metode itu pernah dipanggil.
func (g *Gudang) Dipanggil(nama string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, p := range g.Panggilan {
		if p == nama {
			return true
		}
	}
	return false
}

func (g *Gudang) StatusEDM(context.Context, string) (string, error) {
	g.catat("StatusEDM")
	return g.StatusEDMNilai, g.GalatStatusEDM
}

func (g *Gudang) JenisBisnis(context.Context, string) (string, error) {
	g.catat("JenisBisnis")
	return g.JenisBisnisNilai, g.GalatJenisBisnis
}

func (g *Gudang) Klaim(context.Context, string) ([]string, error) {
	g.catat("Klaim")
	return g.KlaimNilai, g.GalatKlaim
}

func (g *Gudang) AdaPembayaran(context.Context, string) (bool, error) {
	g.catat("AdaPembayaran")
	return g.Pembayaran, g.GalatPembayaran
}

func (g *Gudang) AdaRenewal(context.Context, string) (bool, error) {
	g.catat("AdaRenewal")
	return g.Renewal, g.GalatRenewal
}

func (g *Gudang) AdaFacout(context.Context, string) (bool, error) {
	g.catat("AdaFacout")
	return g.Facout, g.GalatFacout
}

// AdaProteksi dicatat sebagai "AdaProteksi:<tipe>".
func (g *Gudang) AdaProteksi(_ context.Context, _ string, tipe int) (bool, error) {
	g.catat("AdaProteksi:" + strconv.Itoa(tipe))
	return g.Proteksi[tipe], g.GalatProteksi
}

func (g *Gudang) PeriodeProduksi(context.Context, string) (*repository.PeriodeTeks, error) {
	g.catat("PeriodeProduksi")
	return g.Periode, g.GalatPeriode
}

func (g *Gudang) Pengecualian(_ context.Context, jenis string) ([]string, error) {
	g.catat("Pengecualian")
	return g.PengecualianPer[jenis], g.GalatPengecualian
}
