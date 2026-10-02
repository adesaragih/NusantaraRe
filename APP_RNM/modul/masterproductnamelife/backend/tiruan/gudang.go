// Package tiruan memuat gudang DI MEMORI modul Master Product Name Life -
// HANYA untuk uji services dan handlers (tidak diimpor kode produksi; penjaga
// modul `TestMPNLTiruanHanyaDiUji`).
//
// ⭐ Tiruan menyimpan `JSONDATA` MENTAH dan memakai kodek repository yang
// sungguhan (`repository.UraiProduk`, `RingkasanDari`), sehingga uji services
// dan handlers ikut menguji bentuk JSON Pega - bukan bentuk karangan tiruan.
// Fixture selalu berawalan `UJI-`.
package tiruan

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/repository"
)

// Gudang adalah pengganti `repository.Gudang` di memori.
type Gudang struct {
	// Umum - `M_PRODUCT_LIFE`: ID → JSONDATA.
	Umum map[string]string
	// Inward - `M_PRODUCTINWARD_LIFE`: ID baris → JSONDATA.
	Inward map[string]string

	// Master - isi tiap pemilih, urutan RD.
	Master map[models.JenisMaster][]models.NilaiMaster
	// Plan - isi `PRODUCT_TYPE_LIFE`.
	Plan []models.JenisPlan
	// Rate - view `RATE_LIFE` per IDUSEDBY (dialog `View Rate`, K1 01-10-2026).
	Rate map[string][]models.BarisRate
	// CariTerakhir - kata cari terakhir yang diterima CariMaster.
	CariTerakhir string
	// GagalMaster - bila terisi, pembacaan master gagal dengan galat ini.
	GagalMaster error

	// Seq - nomor urut berikut `M_PRODUCT_LIFE_SEQ` (bawaan 44, seperti `START WITH 44`).
	Seq int64
	// Datar - kolom datar `M_PRODUCT_LIFE`: RIRISKID, RIRISK (katalog DEV, lanjutan 1 L1).
	Datar map[string][2]string
	// GagalTulis - bila terisi, setiap penulis menulis LALU gagal dengan galat
	// ini (uji: transaksi gagal = nol tulisan).
	GagalTulis error
	// GagalTulisInward - penulis sisi inward gagal SESUDAH sisi umum tertulis (uji P4).
	GagalTulisInward error

	// Lampiran, Objek, Outbox, AppName - tabel lampiran (paket 8, `lampiran.go`).
	Lampiran      map[string]models.Lampiran
	Objek         map[string]models.ObjekPenyimpanan
	Outbox        []Efek
	AppName       string
	nomorLampiran int

	// KontrakOR - `TREATYCONTRACT_LIFE` × `TREATYYEAR_LIFE` (paket 9, `outward.go`).
	KontrakOR []KontrakOR
	// MintaOR - tanggal (`dd/MM/yyyy`) yang diterima DaftarReinstypeOR terakhir.
	MintaOR [2]string

	// Komit mencacah transaksi yang ditutup sukses.
	Komit int
	// GagalBaca - bila terisi, setiap pembacaan gagal dengan galat ini.
	GagalBaca error
}

// Baru menyusun gudang kosong.
func Baru() *Gudang {
	return &Gudang{Umum: map[string]string{}, Inward: map[string]string{},
		Master: map[models.JenisMaster][]models.NilaiMaster{}, Seq: 44, Datar: map[string][2]string{},
		Rate: map[string][]models.BarisRate{}}
}

// Transaksi - tiruan `DalamTransaksi`: fn(nil); sukses = Komit++, gagal =
// isi tabel dipulihkan (rollback). `Seq` TIDAK dipulihkan: `NEXTVAL` Oracle tidak ikut
// rollback, jadi simpan yang gagal meninggalkan celah nomor (audit 02-10-2026).
func (g *Gudang) Transaksi(_ context.Context, fn func(tx *db.Tx) error) error {
	umum, inward, datar := salin(g.Umum), salin(g.Inward), map[string][2]string{}
	for k, v := range g.Datar {
		datar[k] = v
	}
	if err := fn(nil); err != nil {
		g.Umum, g.Inward, g.Datar = umum, inward, datar
		return err
	}
	g.Komit++
	return nil
}

func salin(m map[string]string) map[string]string {
	h := make(map[string]string, len(m))
	for k, v := range m {
		h[k] = v
	}
	return h
}

// KunciProduk - tiruan `FOR UPDATE`: sama dengan AmbilProduk.
func (g *Gudang) KunciProduk(ctx context.Context, tx *db.Tx, id string) (models.Produk, error) {
	return g.AmbilProduk(ctx, tx, id)
}

func (g *Gudang) tulisDatar(p models.Produk) {
	g.Datar[p.ID] = [2]string{p.Umum.RIRiskID, p.Umum.RIRisk}
}

// SisipProduk - ID dari Seq lewat `repository.PilihIdentitasBebas` (aturan yang sama dengan
// gudang Oracle): ID terpakai di salah satu tabel dilewati, batasnya gagal terang.
func (g *Gudang) SisipProduk(_ context.Context, _ *db.Tx, p models.Produk) (string, error) {
	id, _, err := repository.PilihIdentitasBebas(
		func() (int64, error) {
			n := g.Seq
			g.Seq++
			return n, nil
		},
		func(id string) (bool, error) {
			_, umum := g.Umum[id]
			_, inward := g.Inward[id]
			return umum || inward, nil
		},
	)
	if err != nil {
		return "", err
	}
	dasarUmum, dasarInward := "", ""
	if p.SalinanDari != "" {
		if _, ada := g.Umum[p.SalinanDari]; !ada {
			return "", fmt.Errorf("%w: %s", repository.ErrTidakAda, p.SalinanDari)
		}
		dasarUmum = g.Umum[p.SalinanDari]
		_, dasarInward = g.cariInward(p.SalinanDari)
	}
	p.ID = id
	p.Inward.ID, p.Inward.ProductID = id, id
	umum, err := repository.RakitUmum(p, dasarUmum, true)
	if err != nil {
		return "", err
	}
	g.Umum[id] = umum
	g.tulisDatar(p)
	if g.GagalTulis != nil {
		return "", g.GagalTulis
	}
	if g.GagalTulisInward != nil {
		return "", g.GagalTulisInward
	}
	inward, err := repository.RakitInward(p, dasarInward)
	if err != nil {
		return "", err
	}
	g.Inward[id] = inward
	return id, nil
}

// PerbaruiProduk - kunci JSON lama dipertahankan lewat `repository.RakitUmum`.
func (g *Gudang) PerbaruiProduk(_ context.Context, _ *db.Tx, p models.Produk) error {
	lama, ada := g.Umum[p.ID]
	if !ada {
		return fmt.Errorf("%w: %s", repository.ErrTidakAda, p.ID)
	}
	umum, err := repository.RakitUmum(p, lama, false)
	if err != nil {
		return err
	}
	g.Umum[p.ID] = umum
	g.tulisDatar(p)
	if g.GagalTulis != nil {
		return g.GagalTulis
	}
	if g.GagalTulisInward != nil {
		return g.GagalTulisInward
	}
	p.Inward.ProductID = p.ID
	idIn, isiIn := g.cariInward(p.ID)
	if idIn == "" {
		if lain := g.milikLain(p.ID); lain != "" {
			return fmt.Errorf("%w: inward row %s belongs to product %s", repository.ErrIdentitasBentrok, p.ID, lain)
		}
		idIn = p.ID
	}
	p.Inward.ID = idIn
	inward, err := repository.RakitInward(p, isiIn)
	if err != nil {
		return err
	}
	g.Inward[idIn] = inward
	return nil
}

// DaftarProduk - urut ID.
func (g *Gudang) DaftarProduk(context.Context) ([]models.RingkasanProduk, error) {
	if g.GagalBaca != nil {
		return nil, g.GagalBaca
	}
	id := make([]string, 0, len(g.Umum))
	for k := range g.Umum {
		id = append(id, k)
	}
	sort.Strings(id)
	hasil := []models.RingkasanProduk{}
	for _, k := range id {
		r, err := repository.RingkasanDari(k, g.Umum[k])
		if err != nil {
			return nil, err
		}
		hasil = append(hasil, r)
	}
	return hasil, nil
}

// AmbilProduk - kedua sisi; inward dicari lewat `PRODUCTID` lalu `ID`, seperti repository.
func (g *Gudang) AmbilProduk(_ context.Context, _ *db.Tx, id string) (models.Produk, error) {
	if g.GagalBaca != nil {
		return models.Produk{}, g.GagalBaca
	}
	umum, ada := g.Umum[id]
	if !ada {
		return models.Produk{}, fmt.Errorf("%w: %s", repository.ErrTidakAda, id)
	}
	idIn, isiIn := g.cariInward(id)
	return repository.UraiProduk(id, umum, idIn, isiIn)
}

func (g *Gudang) cariInward(id string) (string, string) {
	kunci := make([]string, 0, len(g.Inward))
	for k := range g.Inward {
		kunci = append(kunci, k)
	}
	sort.Strings(kunci)
	idPilih, isiPilih := "", ""
	for _, k := range kunci {
		p, err := repository.UraiProduk("", "", k, g.Inward[k])
		if err == nil && p.Inward.ProductID == id {
			idPilih, isiPilih = k, g.Inward[k]
		}
	}
	if idPilih != "" {
		return idPilih, isiPilih
	}
	if isi, ada := g.Inward[id]; ada && g.milikLain(id) == "" {
		return id, isi
	}
	return "", ""
}

// milikLain - PRODUCTID baris inward ber-ID = id yang milik produk lain (seperti repository).
func (g *Gudang) milikLain(id string) string {
	isi, ada := g.Inward[id]
	if !ada {
		return ""
	}
	p, err := repository.UraiProduk("", "", id, isi)
	if err != nil || p.Inward.ProductID == "" || p.Inward.ProductID == id {
		return ""
	}
	return p.Inward.ProductID
}

// CariMaster - "Contains" tanpa membedakan huruf, urutan isian.
func (g *Gudang) CariMaster(_ context.Context, jenis models.JenisMaster, kata string, batas int) ([]models.NilaiMaster, error) {
	g.CariTerakhir = kata
	if g.GagalMaster != nil {
		return nil, g.GagalMaster
	}
	if _, dikenal := map[models.JenisMaster]bool{models.MasterCeding: true, models.MasterSOB: true,
		models.MasterPemegangPolis: true, models.MasterMataUang: true, models.MasterRIRisk: true,
		models.MasterPenyebab: true, models.MasterRIRate: true}[jenis]; !dikenal {
		return nil, fmt.Errorf("%w: %q", repository.ErrJenisMasterTidakDikenal, jenis)
	}
	hasil := []models.NilaiMaster{}
	for _, m := range g.Master[jenis] {
		if (batas <= 0 || len(hasil) < batas) && strings.Contains(strings.ToUpper(m.Nama), strings.ToUpper(kata)) {
			hasil = append(hasil, m)
		}
	}
	return hasil, nil
}

// AmbilMaster - satu nilai menurut ID.
func (g *Gudang) AmbilMaster(_ context.Context, jenis models.JenisMaster, id string) (models.NilaiMaster, bool, error) {
	if g.GagalMaster != nil {
		return models.NilaiMaster{}, false, g.GagalMaster
	}
	for _, m := range g.Master[jenis] {
		if m.ID == id {
			return m, true, nil
		}
	}
	return models.NilaiMaster{}, false, nil
}

// CariPlan - "Contains" pada CoverName atau Business.
func (g *Gudang) CariPlan(_ context.Context, kata string) ([]models.JenisPlan, error) {
	if g.GagalMaster != nil {
		return nil, g.GagalMaster
	}
	hasil := []models.JenisPlan{}
	for _, p := range g.Plan {
		if strings.Contains(strings.ToUpper(p.CoverName), kata) || strings.Contains(strings.ToUpper(p.Business), kata) {
			hasil = append(hasil, p)
		}
	}
	return hasil, nil
}

// AmbilPlan - satu jenis plan menurut ID.
func (g *Gudang) AmbilPlan(_ context.Context, id string) (models.JenisPlan, bool, error) {
	if g.GagalMaster != nil {
		return models.JenisPlan{}, false, g.GagalMaster
	}
	for _, p := range g.Plan {
		if p.ID == id {
			return p, true, nil
		}
	}
	return models.JenisPlan{}, false, nil
}

// DaftarRate - satu IDUSEDBY, urut `ID DESC, RATE ASC` (`BrowseRateLife_RD` b748, b786), dipotong
// `repository.BatasRate`.
func (g *Gudang) DaftarRate(_ context.Context, idUsedBy string) ([]models.BarisRate, bool, error) {
	if g.GagalMaster != nil {
		return nil, false, g.GagalMaster
	}
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
	return d, terpotong, nil
}

// GalatMasterUji - galat master tak terbaca berbentuk repository (sebab "ORA-" hanya di log).
func GalatMasterUji(objek string) error {
	return repository.GalatMaster{Objek: objek, Sebab: errors.New("ORA-00942: table or view does not exist")}
}

// ErrTiruan - galat buatan uji.
var ErrTiruan = errors.New("tiruan: injected failure")
