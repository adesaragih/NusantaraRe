// Package tiruan memuat gudang DI MEMORI modul Master Contract Retro Life -
// HANYA untuk uji services dan handlers (tidak diimpor kode produksi; penjaga
// modul `TestMCRLTiruanHanyaDiUji`).
//
// Ia meniru perilaku yang diuji terhadap Oracle oleh uji bertag `db`:
// saringan kunci induk, urutan RD, dan (paket 2+) penulisan. Fixture selalu
// berawalan `UJI-`.
package tiruan

import (
	"context"
	"sort"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/mastercontractretrolife/backend/models"
	"nusantarare/modul/mastercontractretrolife/backend/repository"
)

// Gudang adalah pengganti `repository.Gudang` di memori.
type Gudang struct {
	Tahun     map[string]models.TahunTreaty
	Kontrak   map[string]models.Kontrak
	Reinsurer map[string]models.Reinsurer
	Security  map[string]models.SecurityReinsurer
	Business  map[string]models.Business

	Jenis     []models.JenisReasuransi
	MasterRe  []models.MasterReinsurer
	MasterBiz []models.MasterBusiness

	// GalatMaster - bila terisi, setiap pembacaan master gagal dengan galat ini.
	GalatMaster error
	// Panggilan mencatat argumen kunci induk yang diterima pembaca grid.
	Panggilan []string
	// Komit mencacah transaksi yang ditutup sukses.
	Komit int
}

// Baru menyusun gudang kosong.
func Baru() *Gudang {
	return &Gudang{
		Tahun: map[string]models.TahunTreaty{}, Kontrak: map[string]models.Kontrak{},
		Reinsurer: map[string]models.Reinsurer{}, Security: map[string]models.SecurityReinsurer{},
		Business: map[string]models.Business{},
	}
}

func (g *Gudang) catat(s ...string) { g.Panggilan = append(g.Panggilan, strings.Join(s, "|")) }

func urutID[T any](m map[string]T, saring func(T) bool, id func(T) string, menurun bool) []T {
	var hasil []T
	for _, v := range m {
		if saring(v) {
			hasil = append(hasil, v)
		}
	}
	sort.Slice(hasil, func(i, j int) bool {
		if menurun {
			return id(hasil[i]) > id(hasil[j])
		}
		return id(hasil[i]) < id(hasil[j])
	})
	return hasil
}

func ambil[T any](m map[string]T, id string) (T, error) {
	v, ada := m[id]
	if !ada {
		var nol T
		return nol, repository.ErrTidakAda
	}
	return v, nil
}

// DaftarTahun - `ID ASC`.
func (g *Gudang) DaftarTahun(context.Context) ([]models.TahunTreaty, error) {
	return urutID(g.Tahun, func(models.TahunTreaty) bool { return true },
		func(t models.TahunTreaty) string { return t.ID }, false), nil
}

// AmbilTahun - satu tahun.
func (g *Gudang) AmbilTahun(_ context.Context, _ *db.Tx, id string) (models.TahunTreaty, error) {
	return ambil(g.Tahun, id)
}

// DaftarKontrak - `IDTREATYYEAR =`, `ID ASC`.
func (g *Gudang) DaftarKontrak(_ context.Context, _ *db.Tx, tahunID string) ([]models.Kontrak, error) {
	g.catat("kontrak", tahunID)
	return urutID(g.Kontrak, func(k models.Kontrak) bool { return k.IDTreatyYear == tahunID },
		func(k models.Kontrak) string { return k.ID }, false), nil
}

// AmbilKontrak - satu kontrak.
func (g *Gudang) AmbilKontrak(_ context.Context, _ *db.Tx, id string) (models.Kontrak, error) {
	return ambil(g.Kontrak, id)
}

// DaftarReinsurer - dua kunci induk, `ID DESC`.
func (g *Gudang) DaftarReinsurer(_ context.Context, _ *db.Tx, tahunID, kontrakID string) ([]models.Reinsurer, error) {
	g.catat("reinsurer", tahunID, kontrakID)
	return urutID(g.Reinsurer, func(r models.Reinsurer) bool {
		return r.TreatyYearID == tahunID && r.TreatyContractID == kontrakID
	}, func(r models.Reinsurer) string { return r.ID }, true), nil
}

// AmbilReinsurer - satu reinsurer.
func (g *Gudang) AmbilReinsurer(_ context.Context, _ *db.Tx, id string) (models.Reinsurer, error) {
	return ambil(g.Reinsurer, id)
}

// DaftarSecurity - tiga kunci induk, `ID ASC`.
func (g *Gudang) DaftarSecurity(_ context.Context, _ *db.Tx, tahunID, kontrakID, reinsurerID string) (
	[]models.SecurityReinsurer, error) {

	g.catat("security", tahunID, kontrakID, reinsurerID)
	return urutID(g.Security, func(s models.SecurityReinsurer) bool {
		return s.TreatyYearID == tahunID && s.TreatyContractID == kontrakID && s.TreatyReinsurerID == reinsurerID
	}, func(s models.SecurityReinsurer) string { return s.ID }, false), nil
}

// AmbilSecurity - satu security.
func (g *Gudang) AmbilSecurity(_ context.Context, _ *db.Tx, id string) (models.SecurityReinsurer, error) {
	return ambil(g.Security, id)
}

// DaftarBusiness - dua kunci induk, `TGLUPDATE ASC, ID ASC`.
func (g *Gudang) DaftarBusiness(_ context.Context, _ *db.Tx, tahunID, kontrakID string) ([]models.Business, error) {
	g.catat("business", tahunID, kontrakID)
	d := urutID(g.Business, func(b models.Business) bool {
		return b.TreatyYearID == tahunID && b.TreatyContractID == kontrakID
	}, func(b models.Business) string { return b.ID }, false)
	sort.SliceStable(d, func(i, j int) bool { return d[i].TglUpdate.Before(d[j].TglUpdate) })
	return d, nil
}

// AmbilBusiness - satu business.
func (g *Gudang) AmbilBusiness(_ context.Context, _ *db.Tx, id string) (models.Business, error) {
	return ambil(g.Business, id)
}

// JenisReasuransiLife - daftar tetap uji.
func (g *Gudang) JenisReasuransiLife(context.Context) ([]models.JenisReasuransi, error) {
	return g.Jenis, g.GalatMaster
}

// CariMasterReinsurer - `Contains`, tidak peka huruf besar-kecil.
func (g *Gudang) CariMasterReinsurer(_ context.Context, kata string) ([]models.MasterReinsurer, error) {
	var hasil []models.MasterReinsurer
	for _, m := range g.MasterRe {
		if strings.Contains(strings.ToUpper(m.ClientName), strings.ToUpper(strings.TrimSpace(kata))) {
			hasil = append(hasil, m)
		}
	}
	return hasil, g.GalatMaster
}

// AmbilMasterReinsurer - menurut ID.
func (g *Gudang) AmbilMasterReinsurer(_ context.Context, id string) (models.MasterReinsurer, bool, error) {
	for _, m := range g.MasterRe {
		if m.ID == id {
			return m, true, g.GalatMaster
		}
	}
	return models.MasterReinsurer{}, false, g.GalatMaster
}

// CariMasterBusiness - `Contains` atas Note.
func (g *Gudang) CariMasterBusiness(_ context.Context, kata string) ([]models.MasterBusiness, error) {
	var hasil []models.MasterBusiness
	for _, m := range g.MasterBiz {
		if strings.Contains(strings.ToUpper(m.Note), strings.ToUpper(strings.TrimSpace(kata))) {
			hasil = append(hasil, m)
		}
	}
	return hasil, g.GalatMaster
}

// AmbilMasterBusiness - menurut ID.
func (g *Gudang) AmbilMasterBusiness(_ context.Context, id string) (models.MasterBusiness, bool, error) {
	for _, m := range g.MasterBiz {
		if m.ID == id {
			return m, true, g.GalatMaster
		}
	}
	return models.MasterBusiness{}, false, g.GalatMaster
}

// Transaksi meniru `inti.Dasar.DalamTransaksi` untuk `services.BaruLayanan`:
// fn berjalan dengan tx nil; bila ia gagal, isi kelima tabel DIPULIHKAN ke
// keadaan sebelum fn - sehingga uji atomisitas (kaskade, salin-semua) jujur.
// `Komit` mencacah transaksi yang berhasil ditutup.
func (g *Gudang) Transaksi(_ context.Context, fn func(tx *db.Tx) error) error {
	potret := g.potret()
	if err := fn(nil); err != nil {
		g.pulihkan(potret)
		return err
	}
	g.Komit++
	return nil
}

type potretTabel struct {
	tahun     map[string]models.TahunTreaty
	kontrak   map[string]models.Kontrak
	reinsurer map[string]models.Reinsurer
	security  map[string]models.SecurityReinsurer
	business  map[string]models.Business
}

func salin[T any](m map[string]T) map[string]T {
	s := make(map[string]T, len(m))
	for k, v := range m {
		s[k] = v
	}
	return s
}

func (g *Gudang) potret() potretTabel {
	return potretTabel{salin(g.Tahun), salin(g.Kontrak), salin(g.Reinsurer), salin(g.Security), salin(g.Business)}
}

func (g *Gudang) pulihkan(p potretTabel) {
	g.Tahun, g.Kontrak, g.Reinsurer, g.Security, g.Business = p.tahun, p.kontrak, p.reinsurer, p.security, p.business
}
