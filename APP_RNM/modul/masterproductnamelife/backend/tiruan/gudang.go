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
	// CariTerakhir - kata cari terakhir yang diterima CariMaster.
	CariTerakhir string
	// GagalMaster - bila terisi, pembacaan master gagal dengan galat ini.
	GagalMaster error

	// Seq - nomor urut berikut `M_PRODUCT_LIFE_SEQ` (bawaan 44, seperti `START WITH 44`).
	Seq int64
	// Datar - kolom datar `M_PRODUCT_LIFE`: RIRISKID, RIRISK, PRODUCTNAME, BEGIN_DATE (`dd/MM/yyyy`).
	Datar map[string][4]string
	// GagalTulis - bila terisi, setiap penulis menulis LALU gagal dengan galat
	// ini (uji: transaksi gagal = nol tulisan).
	GagalTulis error

	// Komit mencacah transaksi yang ditutup sukses.
	Komit int
	// GagalBaca - bila terisi, setiap pembacaan gagal dengan galat ini.
	GagalBaca error
}

// Baru menyusun gudang kosong.
func Baru() *Gudang {
	return &Gudang{Umum: map[string]string{}, Inward: map[string]string{},
		Master: map[models.JenisMaster][]models.NilaiMaster{}, Seq: 44, Datar: map[string][4]string{}}
}

// Transaksi - tiruan `DalamTransaksi`: fn(nil); sukses = Komit++, gagal =
// seluruh isi dipulihkan (rollback).
func (g *Gudang) Transaksi(_ context.Context, fn func(tx *db.Tx) error) error {
	umum, inward, datar, seq := salin(g.Umum), salin(g.Inward), map[string][4]string{}, g.Seq
	for k, v := range g.Datar {
		datar[k] = v
	}
	if err := fn(nil); err != nil {
		g.Umum, g.Inward, g.Datar, g.Seq = umum, inward, datar, seq
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
	g.Datar[p.ID] = [4]string{p.Umum.RIRiskID, p.Umum.RIRisk, p.Umum.ProductName, repository.TanggalKePega(p.Inward.Begin)}
}

// SisipProduk - ID dari Seq lewat `repository.FormatIdentitas`; ID terpakai = bentrok.
func (g *Gudang) SisipProduk(_ context.Context, _ *db.Tx, p models.Produk) (string, error) {
	id, err := repository.FormatIdentitas(g.Seq)
	if err != nil {
		return "", err
	}
	g.Seq++
	if _, ada := g.Umum[id]; ada {
		return "", fmt.Errorf("%w: %s", repository.ErrIdentitasBentrok, id)
	}
	if _, ada := g.Inward[id]; ada {
		return "", fmt.Errorf("%w: %s", repository.ErrIdentitasBentrok, id)
	}
	p.ID = id
	umum, err := repository.RakitUmum(p, "", true)
	if err != nil {
		return "", err
	}
	g.Umum[id] = umum
	g.tulisDatar(p)
	if g.GagalTulis != nil {
		return "", g.GagalTulis
	}
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
	return g.GagalTulis
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
	if isi, ada := g.Inward[id]; ada {
		return id, isi
	}
	return "", ""
}

// CariMaster - "Contains" tanpa membedakan huruf, urutan isian.
func (g *Gudang) CariMaster(_ context.Context, jenis models.JenisMaster, kata string) ([]models.NilaiMaster, error) {
	g.CariTerakhir = kata
	if g.GagalMaster != nil {
		return nil, g.GagalMaster
	}
	if _, dikenal := map[models.JenisMaster]bool{models.MasterCeding: true, models.MasterSOB: true,
		models.MasterPemegangPolis: true, models.MasterMataUang: true, models.MasterRIRisk: true,
		models.MasterPenyebab: true}[jenis]; !dikenal {
		return nil, fmt.Errorf("%w: %q", repository.ErrJenisMasterTidakDikenal, jenis)
	}
	hasil := []models.NilaiMaster{}
	for _, m := range g.Master[jenis] {
		if strings.Contains(strings.ToUpper(m.Nama), strings.ToUpper(kata)) {
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

// GalatMasterUji - galat master tak terbaca berbentuk repository (sebab "ORA-" hanya di log).
func GalatMasterUji(objek string) error {
	return repository.GalatMaster{Objek: objek, Sebab: errors.New("ORA-00942: table or view does not exist")}
}

// ErrTiruan - galat buatan uji.
var ErrTiruan = errors.New("tiruan: injected failure")
