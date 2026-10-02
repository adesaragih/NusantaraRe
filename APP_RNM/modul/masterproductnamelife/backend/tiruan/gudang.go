// Package tiruan memuat gudang DI MEMORI modul Master Product Name Life -
// HANYA untuk uji services dan handlers (tidak diimpor kode produksi; penjaga
// modul `TestMPNLTiruanHanyaDiUji`).
//
// ⭐ Sejak 02-10-2026 (tabel flat, tiket 01 bab bertanggal) tiruan menyimpan produk BERBENTUK FLAT lewat pemetaan
// repository yang sungguhan (`repository.NormalkanFlat`): angka kanonik, nilai tak muat = `ErrNilaiTidakMuat`, medan
// tanpa kolom kosong - persis yang dibaca kembali dari Oracle (`mpnl_flat_db_test.go` membuktikan kesamaannya).
// Fixture berbentuk JSON lama diisi lewat `IsiJSON` (kodek `repository.UraiProduk`, seperti alat pindah).
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
	// Produk - tabel flat `M_PRODUCTNAME_LIFE` + tujuh anak: ID → produk berbentuk flat (NormalkanFlat).
	Produk map[string]models.Produk

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
	// GagalTulis - bila terisi, setiap penulis menulis LALU gagal dengan galat
	// ini (uji: transaksi gagal = nol tulisan).
	GagalTulis error
	// GagalTulisAnak - penulis tabel anak gagal SESUDAH baris induk tertulis (uji P4: satu transaksi).
	GagalTulisAnak error

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
	return &Gudang{Produk: map[string]models.Produk{},
		Master: map[models.JenisMaster][]models.NilaiMaster{}, Seq: 44, Rate: map[string][]models.BarisRate{}}
}

// Transaksi - tiruan `DalamTransaksi`: fn(nil); sukses = Komit++, gagal =
// isi tabel dipulihkan (rollback). `Seq` TIDAK dipulihkan: `NEXTVAL` Oracle tidak ikut
// rollback, jadi simpan yang gagal meninggalkan celah nomor (audit 02-10-2026).
func (g *Gudang) Transaksi(_ context.Context, fn func(tx *db.Tx) error) error {
	produk := make(map[string]models.Produk, len(g.Produk))
	for k, v := range g.Produk {
		produk[k] = v
	}
	if err := fn(nil); err != nil {
		g.Produk = produk
		return err
	}
	g.Komit++
	return nil
}

// IsiJSON - fixture produk dari JSON lama kedua tabel (`jsonInward` kosong = tanpa baris inward), dipindah seperti alat
// pindah: kodek `repository.UraiProduk` lalu bentuk flat. Medan tanpa kolom yang terisi atau nilai tak muat = panik
// (fixture yang keliru, bukan keadaan yang diuji).
func (g *Gudang) IsiJSON(id, jsonUmum, jsonInward string) models.Produk {
	idInward := ""
	if jsonInward != "" {
		idInward = id
	}
	p, err := repository.UraiProduk(id, jsonUmum, idInward, jsonInward)
	if err != nil {
		panic(fmt.Sprintf("tiruan: fixture %s: %v", id, err))
	}
	if m := repository.MedanTanpaKolom(p); len(m) > 0 {
		panic(fmt.Sprintf("tiruan: fixture %s mengisi medan tanpa kolom flat %v", id, m))
	}
	n, masalah := repository.NormalkanFlat(p)
	if len(masalah) > 0 {
		panic(fmt.Sprintf("tiruan: fixture %s: %+v", id, masalah))
	}
	g.Produk[id] = n
	return salinProduk(n)
}

// salinProduk - salinan dalam (daftar baru), supaya pemanggil tidak mengubah isi tabel tiruan.
func salinProduk(p models.Produk) models.Produk {
	p.LienClause = append([]models.BarisLien{}, p.LienClause...)
	p.DocumentClaim = append([]models.BarisDokumen{}, p.DocumentClaim...)
	p.PlanList = append([]models.BarisPlan{}, p.PlanList...)
	p.FinancialUnderwriting = append([]models.BarisFinUW{}, p.FinancialUnderwriting...)
	p.UnderwritingLimit = append([]models.BarisUWLimit{}, p.UnderwritingLimit...)
	p.OutwardList = append([]models.BarisOutward{}, p.OutwardList...)
	p.CommentList = append([]models.BarisKomentar{}, p.CommentList...)
	return p
}

// tulis - satu produk ke tabel flat tiruan (induk + anak), seperti `tulisFlat` repository.
func (g *Gudang) tulis(p models.Produk) error {
	n, masalah := repository.NormalkanFlat(p)
	if len(masalah) > 0 {
		return fmt.Errorf("%w: %s.%s (%s)", repository.ErrNilaiTidakMuat, masalah[0].Tabel, masalah[0].Kolom, masalah[0].Jenis)
	}
	induk := n
	induk.LienClause, induk.DocumentClaim, induk.PlanList = nil, nil, nil
	induk.FinancialUnderwriting, induk.UnderwritingLimit, induk.OutwardList, induk.CommentList = nil, nil, nil, nil
	g.Produk[n.ID] = induk
	if g.GagalTulis != nil {
		return g.GagalTulis
	}
	if g.GagalTulisAnak != nil {
		return g.GagalTulisAnak
	}
	g.Produk[n.ID] = n
	return nil
}

// KunciProduk - tiruan `FOR UPDATE`: sama dengan AmbilProduk.
func (g *Gudang) KunciProduk(ctx context.Context, tx *db.Tx, id string) (models.Produk, error) {
	return g.AmbilProduk(ctx, tx, id)
}

// SisipProduk - ID dari Seq lewat `repository.PilihIdentitasBebas` (aturan yang sama dengan gudang Oracle): ID yang
// dipakai induk flat ATAU tabel JSON warisan dilewati, batasnya gagal terang.
func (g *Gudang) SisipProduk(_ context.Context, _ *db.Tx, p models.Produk) (string, error) {
	if p.SalinanDari != "" {
		if _, ada := g.Produk[p.SalinanDari]; !ada {
			return "", fmt.Errorf("%w: %s", repository.ErrTidakAda, p.SalinanDari)
		}
	}
	id, _, err := repository.PilihIdentitasBebas(
		func() (int64, error) {
			n := g.Seq
			g.Seq++
			return n, nil
		},
		func(id string) (bool, error) {
			_, flat := g.Produk[id]
			return flat, nil
		},
	)
	if err != nil {
		return "", err
	}
	p.ID = id
	p.Inward.ID, p.Inward.ProductID = id, id
	if err := g.tulis(p); err != nil {
		return "", err
	}
	return id, nil
}

// PerbaruiProduk - baris induk diperbarui (harus ada), anak ditulis ulang.
func (g *Gudang) PerbaruiProduk(_ context.Context, _ *db.Tx, p models.Produk) error {
	if _, ada := g.Produk[p.ID]; !ada {
		return fmt.Errorf("%w: %s", repository.ErrTidakAda, p.ID)
	}
	p.Inward.ID, p.Inward.ProductID = p.ID, p.ID
	return g.tulis(p)
}

// DaftarProduk - urut ID, kelima kolom grid dari induk.
func (g *Gudang) DaftarProduk(context.Context) ([]models.RingkasanProduk, error) {
	if g.GagalBaca != nil {
		return nil, g.GagalBaca
	}
	id := make([]string, 0, len(g.Produk))
	for k := range g.Produk {
		id = append(id, k)
	}
	sort.Strings(id)
	hasil := []models.RingkasanProduk{}
	for _, k := range id {
		u := g.Produk[k].Umum
		hasil = append(hasil, models.RingkasanProduk{ID: k, Ceding: u.Ceding, TreatyNumber: u.TreatyNumber,
			InwardName: u.InwardName, CreateOp: u.CreateOp, UpdateOp: u.UpdateOp})
	}
	return hasil, nil
}

// AmbilProduk - satu produk utuh (salinan).
func (g *Gudang) AmbilProduk(_ context.Context, _ *db.Tx, id string) (models.Produk, error) {
	if g.GagalBaca != nil {
		return models.Produk{}, g.GagalBaca
	}
	p, ada := g.Produk[id]
	if !ada {
		return models.Produk{}, fmt.Errorf("%w: %s", repository.ErrTidakAda, id)
	}
	return salinProduk(p), nil
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
