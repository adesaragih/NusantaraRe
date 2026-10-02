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
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
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
	// RingkasanRate - view `RATE_LIFE_SUMMARY` tiruan; Rate - view `RATE_LIFE` per `IDUSEDBY`.
	RingkasanRate []models.RingkasanRate
	Rate          map[string][]models.BarisRate

	// GalatMaster - bila terisi, setiap pembacaan master gagal dengan galat ini.
	GalatMaster error
	// Panggilan mencatat argumen kunci induk yang diterima pembaca grid.
	Panggilan []string
	// Komit mencacah transaksi yang ditutup sukses.
	Komit int
	// Seq - nomor urut berikut per tabel (bawaan 44, seperti `START WITH 44`).
	Seq map[string]int64
	// Jam - cap `TGLUPDATE` (SYSDATE tiruan).
	Jam time.Time
	// GagalTulis - bila terisi, penulis yang namanya disebut gagal dengan galat
	// itu (uji atomisitas dan K8).
	GagalTulis map[string]error
	hitung     map[string]int
	// SelaHapus - dijalankan HapusKontrak sebelum menghapus: meniru baris yang
	// lahir di antara pencacahan dan penghapusan (lapis kedua kaskade).
	SelaHapus func()
	// Kunci mencatat baris yang dikunci `KunciBaris` ("jenis|id"), berurutan.
	Kunci []string
}

// KunciBaris - tiruan `SELECT … FOR UPDATE`: mencatat kuncinya; baris tidak ada = ErrTidakAda.
func (g *Gudang) KunciBaris(_ context.Context, _ *db.Tx, jenis, id string) error {
	g.Kunci = append(g.Kunci, jenis+"|"+id)
	var ada bool
	switch jenis {
	case "kontrak":
		_, ada = g.Kontrak[id]
	case "reinsurer":
		_, ada = g.Reinsurer[id]
	case "security":
		_, ada = g.Security[id]
	case "business":
		_, ada = g.Business[id]
	default:
		return fmt.Errorf("tiruan: unknown row kind to lock %q", jenis)
	}
	if !ada {
		return fmt.Errorf("%w: %s %s", repository.ErrTidakAda, jenis, id)
	}
	return nil
}

func (g *Gudang) nomorBaru(tabel string) (string, error) {
	n, ada := g.Seq[tabel]
	if !ada {
		n = 44
	}
	g.Seq[tabel] = n + 1
	return repository.FormatIdentitas(n)
}

// gagal - galat tiruan untuk penulis `nama`: `GagalTulis[nama]` selalu, atau
// `GagalTulis[nama#n]` pada panggilan ke-n (uji gagal di tengah).
func (g *Gudang) gagal(nama string) error {
	g.hitung[nama]++
	if err := g.GagalTulis[fmt.Sprintf("%s#%d", nama, g.hitung[nama])]; err != nil {
		return err
	}
	return g.GagalTulis[nama]
}

// SisipBusiness - ID dari sequence tiruan.
func (g *Gudang) SisipBusiness(_ context.Context, _ *db.Tx, b models.Business) (string, error) {
	if err := g.gagal("SisipBusiness"); err != nil {
		return "", err
	}
	id, err := g.nomorBaru(repository.TabelBusiness)
	if err != nil {
		return "", err
	}
	b.ID, b.TglUpdate = id, g.Jam
	g.Business[id] = b
	return id, nil
}

// PerbaruiBusiness - kunci induk tidak berpindah.
func (g *Gudang) PerbaruiBusiness(_ context.Context, _ *db.Tx, b models.Business) error {
	if err := g.gagal("PerbaruiBusiness"); err != nil {
		return err
	}
	lama, ada := g.Business[b.ID]
	if !ada {
		return repository.ErrTidakAda
	}
	b.TreatyYearID, b.TreatyContractID, b.TglUpdate = lama.TreatyYearID, lama.TreatyContractID, g.Jam
	g.Business[b.ID] = b
	return nil
}

// SisipTahun - ID dari sequence tiruan, TGLUPDATE = Jam.
func (g *Gudang) SisipTahun(_ context.Context, _ *db.Tx, t models.TahunTreaty) (string, error) {
	if err := g.gagal("SisipTahun"); err != nil {
		return "", err
	}
	id, err := g.nomorBaru(repository.TabelTahun)
	if err != nil {
		return "", err
	}
	t.ID, t.TglUpdate = id, g.Jam
	g.Tahun[id] = t
	return id, nil
}

// PerbaruiTahun - seluruh medan; ID tak ada = ErrTidakAda.
func (g *Gudang) PerbaruiTahun(_ context.Context, _ *db.Tx, t models.TahunTreaty) error {
	if err := g.gagal("PerbaruiTahun"); err != nil {
		return err
	}
	if _, ada := g.Tahun[t.ID]; !ada {
		return repository.ErrTidakAda
	}
	t.TglUpdate = g.Jam
	g.Tahun[t.ID] = t
	return nil
}

// SisipKontrak - ID dari sequence tiruan.
func (g *Gudang) SisipKontrak(_ context.Context, _ *db.Tx, k models.Kontrak) (string, error) {
	if err := g.gagal("SisipKontrak"); err != nil {
		return "", err
	}
	id, err := g.nomorBaru(repository.TabelKontrak)
	if err != nil {
		return "", err
	}
	k.ID, k.TglUpdate = id, g.Jam
	g.Kontrak[id] = k
	return id, nil
}

// PerbaruiKontrak - IDTREATYYEAR tidak berpindah.
func (g *Gudang) PerbaruiKontrak(_ context.Context, _ *db.Tx, k models.Kontrak) error {
	if err := g.gagal("PerbaruiKontrak"); err != nil {
		return err
	}
	lama, ada := g.Kontrak[k.ID]
	if !ada {
		return repository.ErrTidakAda
	}
	k.IDTreatyYear, k.TglUpdate = lama.IDTreatyYear, g.Jam
	g.Kontrak[k.ID] = k
	return nil
}

// SalinKontrakKeAnak - K4 atas reinsurer dan business kontrak itu.
func (g *Gudang) SalinKontrakKeAnak(_ context.Context, _ *db.Tx, k models.Kontrak) (int64, error) {
	var n int64
	for id, r := range g.Reinsurer {
		if r.TreatyContractID == k.ID && (r.ReinsTypeID != k.ReinsTypeID || r.ReinsTypeName != k.ReinsTypeName) {
			r.ReinsTypeID, r.ReinsTypeName, r.UserID, r.TglUpdate = k.ReinsTypeID, k.ReinsTypeName, k.UserID, g.Jam
			g.Reinsurer[id] = r
			n++
		}
	}
	for id, b := range g.Business {
		if b.TreatyContractID == k.ID && (b.ReinsTypeID != k.ReinsTypeID || b.ReinsTypeName != k.ReinsTypeName) {
			b.ReinsTypeID, b.ReinsTypeName, b.UserID, b.TglUpdate = k.ReinsTypeID, k.ReinsTypeName, k.UserID, g.Jam
			g.Business[id] = b
			n++
		}
	}
	return n, nil
}

// SisipReinsurer - ID dari sequence tiruan.
func (g *Gudang) SisipReinsurer(_ context.Context, _ *db.Tx, r models.Reinsurer) (string, error) {
	if err := g.gagal("SisipReinsurer"); err != nil {
		return "", err
	}
	id, err := g.nomorBaru(repository.TabelReinsurer)
	if err != nil {
		return "", err
	}
	r.ID, r.TglUpdate = id, g.Jam
	g.Reinsurer[id] = r
	return id, nil
}

// PerbaruiReinsurer - kunci induk tidak berpindah.
func (g *Gudang) PerbaruiReinsurer(_ context.Context, _ *db.Tx, r models.Reinsurer) error {
	if err := g.gagal("PerbaruiReinsurer"); err != nil {
		return err
	}
	lama, ada := g.Reinsurer[r.ID]
	if !ada {
		return repository.ErrTidakAda
	}
	r.TreatyYearID, r.TreatyContractID, r.TglUpdate = lama.TreatyYearID, lama.TreatyContractID, g.Jam
	g.Reinsurer[r.ID] = r
	return nil
}

// SisipSecurity - ID dari sequence tiruan.
func (g *Gudang) SisipSecurity(_ context.Context, _ *db.Tx, s models.SecurityReinsurer) (string, error) {
	if err := g.gagal("SisipSecurity"); err != nil {
		return "", err
	}
	id, err := g.nomorBaru(repository.TabelSecurity)
	if err != nil {
		return "", err
	}
	s.ID, s.TglUpdate = id, g.Jam
	g.Security[id] = s
	return id, nil
}

// PerbaruiSecurity - kunci induk tidak berpindah.
func (g *Gudang) PerbaruiSecurity(_ context.Context, _ *db.Tx, s models.SecurityReinsurer) error {
	if err := g.gagal("PerbaruiSecurity"); err != nil {
		return err
	}
	lama, ada := g.Security[s.ID]
	if !ada {
		return repository.ErrTidakAda
	}
	s.TreatyYearID, s.TreatyContractID, s.TreatyReinsurerID, s.TglUpdate =
		lama.TreatyYearID, lama.TreatyContractID, lama.TreatyReinsurerID, g.Jam
	g.Security[s.ID] = s
	return nil
}

// dampakKontrak - predikat sama dengan repository: security lewat reinsurer
// kontrak, reinsurer/business lewat TREATYCONTRACTID.
func (g *Gudang) dampakKontrak(id string, hapus bool) models.Dampak {
	var d models.Dampak
	reinsurer := map[string]bool{}
	for rid, r := range g.Reinsurer {
		if r.TreatyContractID == id {
			reinsurer[rid] = true
		}
	}
	for sid, s := range g.Security {
		if reinsurer[s.TreatyReinsurerID] {
			d.Security++
			if hapus {
				delete(g.Security, sid)
			}
		}
	}
	for rid := range reinsurer {
		d.Reinsurer++
		if hapus {
			delete(g.Reinsurer, rid)
		}
	}
	for bid, b := range g.Business {
		if b.TreatyContractID == id {
			d.Business++
			if hapus {
				delete(g.Business, bid)
			}
		}
	}
	return d
}

func (g *Gudang) dampakReinsurer(id string, hapus bool) models.Dampak {
	var d models.Dampak
	for sid, s := range g.Security {
		if s.TreatyReinsurerID == id {
			d.Security++
			if hapus {
				delete(g.Security, sid)
			}
		}
	}
	return d
}

// DampakHapusKontrak - pencacah tiruan.
func (g *Gudang) DampakHapusKontrak(_ context.Context, _ *db.Tx, id string) (models.Dampak, error) {
	return g.dampakKontrak(id, false), nil
}

// HapusKontrak - anak lebih dulu, lalu kontrak.
func (g *Gudang) HapusKontrak(_ context.Context, _ *db.Tx, id string) (models.Dampak, error) {
	if g.SelaHapus != nil {
		g.SelaHapus()
	}
	d := g.dampakKontrak(id, true)
	if err := g.gagal("HapusKontrak"); err != nil {
		return d, err
	}
	if _, ada := g.Kontrak[id]; !ada {
		return d, repository.ErrTidakAda
	}
	delete(g.Kontrak, id)
	return d, nil
}

// DampakHapusReinsurer - pencacah tiruan.
func (g *Gudang) DampakHapusReinsurer(_ context.Context, _ *db.Tx, id string) (models.Dampak, error) {
	return g.dampakReinsurer(id, false), nil
}

// HapusReinsurer - security lebih dulu, lalu reinsurer.
func (g *Gudang) HapusReinsurer(_ context.Context, _ *db.Tx, id string) (models.Dampak, error) {
	d := g.dampakReinsurer(id, true)
	if err := g.gagal("HapusReinsurer"); err != nil {
		return d, err
	}
	if _, ada := g.Reinsurer[id]; !ada {
		return d, repository.ErrTidakAda
	}
	delete(g.Reinsurer, id)
	return d, nil
}

// HapusSecurity - daun.
func (g *Gudang) HapusSecurity(_ context.Context, _ *db.Tx, id string) error {
	if err := g.gagal("HapusSecurity"); err != nil {
		return err
	}
	if _, ada := g.Security[id]; !ada {
		return repository.ErrTidakAda
	}
	delete(g.Security, id)
	return nil
}

// HapusBusiness - daun.
func (g *Gudang) HapusBusiness(_ context.Context, _ *db.Tx, id string) error {
	if err := g.gagal("HapusBusiness"); err != nil {
		return err
	}
	if _, ada := g.Business[id]; !ada {
		return repository.ErrTidakAda
	}
	delete(g.Business, id)
	return nil
}

// TotalSharePerKontrak - LEFT JOIN tiruan atas dua kunci induk.
func (g *Gudang) TotalSharePerKontrak(_ context.Context, tahunID string) ([]models.TotalShareKontrak, error) {
	kontrak := urutID(g.Kontrak, func(k models.Kontrak) bool { return tahunID == "" || k.IDTreatyYear == tahunID },
		func(k models.Kontrak) string { return k.IDTreatyYear + "|" + k.ID }, false)
	var hasil []models.TotalShareKontrak
	for _, k := range kontrak {
		t := models.TotalShareKontrak{KontrakID: k.ID, TahunID: k.IDTreatyYear, TreatyYear: g.Tahun[k.IDTreatyYear].TreatyYear,
			ReinsTypeName: k.ReinsTypeName}
		for _, r := range g.Reinsurer {
			if r.TreatyContractID != k.ID || r.TreatyYearID != k.IDTreatyYear || r.PctShare == nil {
				continue
			}
			if t.Total == nil {
				t.Total = new(apd.Decimal)
			}
			if _, err := utils.DecimalContext().Add(t.Total, t.Total, r.PctShare); err != nil {
				return nil, err
			}
		}
		hasil = append(hasil, t)
	}
	return hasil, nil
}

// SalinTahunKeAnak - K4/R5 atas business dan kontrak tahun itu yang berbeda.
func (g *Gudang) SalinTahunKeAnak(_ context.Context, _ *db.Tx, t models.TahunTreaty) (int64, error) {
	var n int64
	for id, b := range g.Business {
		// Hanya business yang kontrak induknya ada di tahun itu (perbaikan 01-10-2026, sama dengan SQL-nya).
		induk, ada := g.Kontrak[b.TreatyContractID]
		if !ada || induk.IDTreatyYear != t.ID {
			continue
		}
		if b.TreatyYearID == t.ID && b.TreatyYear != t.TreatyYear {
			b.TreatyYear, b.UserID, b.TglUpdate = t.TreatyYear, t.UserID, g.Jam
			g.Business[id] = b
			n++
		}
	}
	for id, k := range g.Kontrak {
		if k.IDTreatyYear == t.ID && (!k.TreatyStartDate.Equal(t.StartDate) || !k.TreatyEndDate.Equal(t.EndDate)) {
			k.TreatyStartDate, k.TreatyEndDate, k.UserID, k.TglUpdate = t.StartDate, t.EndDate, t.UserID, g.Jam
			g.Kontrak[id] = k
			n++
		}
	}
	return n, nil
}

// Baru menyusun gudang kosong.
func Baru() *Gudang {
	return &Gudang{
		Tahun: map[string]models.TahunTreaty{}, Kontrak: map[string]models.Kontrak{},
		Reinsurer: map[string]models.Reinsurer{}, Security: map[string]models.SecurityReinsurer{},
		Business: map[string]models.Business{}, Seq: map[string]int64{},
		Jam: time.Date(2026, 9, 30, 10, 0, 0, 0, time.UTC), GagalTulis: map[string]error{}, hitung: map[string]int{},
		// Satu tabel rate bawaan: business baru di uji memilih `UJI-RATE` (pilihan baru wajib ada di view).
		RingkasanRate: []models.RingkasanRate{{ID: "UJI-RATE", UsedBy: "UJI R"}},
		Rate:          map[string][]models.BarisRate{},
	}
}

// galatMaster - pembungkus yang sama dengan repository.
func (g *Gudang) galatMaster(objek string) error {
	if g.GalatMaster == nil {
		return nil
	}
	return repository.MasterTidakTerbaca(objek, g.GalatMaster)
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
	return g.Jenis, g.galatMaster(repository.MasterJenisReasuransi)
}

// CariMasterReinsurer - `Contains`, tidak peka huruf besar-kecil.
func (g *Gudang) CariMasterReinsurer(_ context.Context, kata string) ([]models.MasterReinsurer, error) {
	var hasil []models.MasterReinsurer
	for _, m := range g.MasterRe {
		if strings.Contains(strings.ToUpper(m.ClientName), strings.ToUpper(strings.TrimSpace(kata))) {
			hasil = append(hasil, m)
		}
	}
	return hasil, g.galatMaster(repository.MasterReinsurer)
}

// AmbilMasterReinsurer - menurut ID.
func (g *Gudang) AmbilMasterReinsurer(_ context.Context, id string) (models.MasterReinsurer, bool, error) {
	for _, m := range g.MasterRe {
		if m.ID == id {
			return m, true, g.galatMaster(repository.MasterReinsurer)
		}
	}
	return models.MasterReinsurer{}, false, g.galatMaster(repository.MasterReinsurer)
}

// CariMasterBusiness - `Contains` atas Note.
func (g *Gudang) CariMasterBusiness(_ context.Context, kata string) ([]models.MasterBusiness, error) {
	var hasil []models.MasterBusiness
	for _, m := range g.MasterBiz {
		if strings.Contains(strings.ToUpper(m.Note), strings.ToUpper(strings.TrimSpace(kata))) {
			hasil = append(hasil, m)
		}
	}
	return hasil, g.galatMaster(repository.MasterBusiness)
}

// AmbilMasterBusiness - menurut ID.
func (g *Gudang) AmbilMasterBusiness(_ context.Context, id string) (models.MasterBusiness, bool, error) {
	for _, m := range g.MasterBiz {
		if m.ID == id {
			return m, true, g.galatMaster(repository.MasterBusiness)
		}
	}
	return models.MasterBusiness{}, false, g.galatMaster(repository.MasterBusiness)
}

// CariRingkasanRate - `Contains` atas USEDBY, urut `ID ASC` (`BrowseRateLifeSummary` b692).
func (g *Gudang) CariRingkasanRate(_ context.Context, kata string) ([]models.RingkasanRate, error) {
	var hasil []models.RingkasanRate
	for _, m := range g.RingkasanRate {
		if strings.Contains(strings.ToUpper(m.UsedBy), strings.ToUpper(strings.TrimSpace(kata))) {
			hasil = append(hasil, m)
		}
	}
	sort.SliceStable(hasil, func(i, j int) bool { return hasil[i].ID < hasil[j].ID })
	return hasil, g.galatMaster(repository.MasterRingkasanRate)
}

// AmbilRingkasanRate - menurut ID.
func (g *Gudang) AmbilRingkasanRate(_ context.Context, id string) (models.RingkasanRate, bool, error) {
	for _, m := range g.RingkasanRate {
		if m.ID == id {
			return m, true, g.galatMaster(repository.MasterRingkasanRate)
		}
	}
	return models.RingkasanRate{}, false, g.galatMaster(repository.MasterRingkasanRate)
}

// DaftarRate - satu `IDUSEDBY`, urut `ID DESC, RATE ASC` (`BrowseRateLife_RD` b747, b784), dipotong
// `repository.BatasRate`.
func (g *Gudang) DaftarRate(_ context.Context, idUsedBy string) ([]models.BarisRate, bool, error) {
	g.catat("rate", idUsedBy)
	d := append([]models.BarisRate{}, g.Rate[idUsedBy]...)
	sort.SliceStable(d, func(i, j int) bool {
		if d[i].ID != d[j].ID {
			return d[i].ID > d[j].ID
		}
		return d[i].Rate < d[j].Rate
	})
	terpotong := len(d) > repository.BatasRate
	if terpotong {
		d = d[:repository.BatasRate]
	}
	return d, terpotong, g.galatMaster(repository.MasterRate)
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
