// Package tiruan adalah Gudang Accounts di memori - untuk uji services dan handlers TANPA Oracle.
//
// Perilakunya meniru SQL repository: daftar terbaru (nomor ACC terbesar) dulu, cari tanpa beda huruf di Insured
// Name, Group Business, ID, dan Insured ID, organisasi hanya ber-FLAG Org. Seluruh data uji berawalan `UJI`.
package tiruan

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/accounts/backend/models"
	"nusantarare/modul/accounts/backend/repository"
)

// Gudang menyimpan akun, organisasi, dan Group Business di memori.
type Gudang struct {
	Akun map[string]models.Account
	// Org - organisasi CLIENT ber-FLAG Org; NonOrg - CLIENT lain (orang), tidak boleh terpilih.
	Org, NonOrg map[string]models.Organisasi
	GB          map[string]models.GroupBusiness
	// Nomor - nilai SEQ_T_M_ACCOUNT berikutnya.
	Nomor int64
	// Disisip - jejak tulisan untuk asersi uji.
	Disisip []models.Account
	// TanggalBaru - CREATEDATE yang "SYSDATE" tiruan berikan.
	TanggalBaru string
}

// Baru menyusun gudang kosong.
func Baru() *Gudang {
	return &Gudang{Akun: map[string]models.Account{}, Org: map[string]models.Organisasi{},
		NonOrg: map[string]models.Organisasi{}, GB: map[string]models.GroupBusiness{}, Nomor: 4916562,
		TanggalBaru: "2026-10-04 09:00"}
}

// Contoh - dua organisasi, satu orang, dua Group Business, dua akun lama.
func Contoh() *Gudang {
	g := Baru()
	for _, o := range []models.Organisasi{
		{ID: models.AwalanOrg + "ORG-101", IDView: "ORG-101", Nama: "UJI PT SATU"},
		{ID: models.AwalanOrg + "ORG-102", IDView: "ORG-102", Nama: "UJI PT DUA"},
	} {
		g.Org[o.ID] = o
	}
	g.NonOrg["UJI-ORANG-1"] = models.Organisasi{ID: "UJI-ORANG-1", IDView: "PER-1", Nama: "UJI ORANG"}
	g.GB["10001"] = models.GroupBusiness{ID: "10001", Note: "UJI GROUP A"}
	g.GB["10002"] = models.GroupBusiness{ID: "10002", Note: "UJI GROUP B"}
	for _, a := range []models.Account{
		{ID: models.AwalanID + "ACC-4916560", GroupBusinessID: "10001", GroupBusiness: "UJI GROUP A",
			InsuredID: models.AwalanOrg + "ORG-101", InsuredName: "UJI PT SATU"},
		{ID: models.AwalanID + "ACC-4916561", GroupBusinessID: "10002", GroupBusiness: "UJI GROUP B LAMA",
			InsuredID: models.AwalanOrg + "ORG-102", InsuredName: "UJI PT DUA"},
	} {
		g.simpan(a)
	}
	return g
}

// Transaksi tiruan - fn(nil).
func Transaksi(_ context.Context, fn func(tx *db.Tx) error) error { return fn(nil) }

func (g *Gudang) simpan(a models.Account) {
	a.IDView = models.BagianAkhir(a.ID, models.AwalanID)
	a.OrgID = models.BagianAkhir(a.InsuredID, models.AwalanOrg)
	g.Akun[a.ID] = a
}

func nomor(id string) int64 {
	i := strings.LastIndexAny(id, "-")
	n, err := strconv.ParseInt(id[i+1:], 10, 64)
	if err != nil {
		return -1
	}
	return n
}

// Daftar - nomor ACC terbesar dulu; cari tanpa beda huruf.
func (g *Gudang) Daftar(_ context.Context, s models.Saringan) ([]models.Account, int, error) {
	kata := strings.ToUpper(strings.TrimSpace(s.Cari))
	cocok := []models.Account{}
	for _, a := range g.Akun {
		if kata == "" || strings.Contains(strings.ToUpper(a.InsuredName), kata) ||
			strings.Contains(strings.ToUpper(a.GroupBusiness), kata) || strings.Contains(strings.ToUpper(a.ID), kata) ||
			strings.Contains(strings.ToUpper(a.InsuredID), kata) {
			cocok = append(cocok, a)
		}
	}
	sort.Slice(cocok, func(i, j int) bool { return nomor(cocok[i].ID) > nomor(cocok[j].ID) })
	awal := (s.Halaman - 1) * s.Ukuran
	if awal > len(cocok) {
		awal = len(cocok)
	}
	akhir := awal + s.Ukuran
	if akhir > len(cocok) {
		akhir = len(cocok)
	}
	return cocok[awal:akhir], len(cocok), nil
}

// Ambil - satu akun.
func (g *Gudang) Ambil(_ context.Context, _ *db.Tx, id string) (models.Account, error) {
	a, ada := g.Akun[id]
	if !ada {
		return models.Account{}, repository.ErrTidakAda
	}
	return a, nil
}

// AdaID - ID sudah terpakai.
func (g *Gudang) AdaID(_ context.Context, _ *db.Tx, id string) (bool, error) {
	_, ada := g.Akun[id]
	return ada, nil
}

// NomorBerikut - sequence tiruan.
func (g *Gudang) NomorBerikut(context.Context, *db.Tx) (int64, error) {
	n := g.Nomor
	g.Nomor++
	return n, nil
}

// Sisip - baris baru; CREATEDATE = TanggalBaru.
func (g *Gudang) Sisip(_ context.Context, _ *db.Tx, a models.Account) error {
	a.CreateDate = g.TanggalBaru
	g.Disisip = append(g.Disisip, a)
	g.simpan(a)
	return nil
}

// PasanganLain - akun yang sudah ada ber-Insured + Group Business sama.
func (g *Gudang) PasanganLain(_ context.Context, _ *db.Tx, insuredID, groupBusinessID string) (string, error) {
	var ids []string
	for id, a := range g.Akun {
		if a.InsuredID == insuredID && a.GroupBusinessID == groupBusinessID {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	if len(ids) == 0 {
		return "", nil
	}
	return ids[0], nil
}

// CariOrganisasi - nama atau ORG-n mengandung kata cari, urut nama.
func (g *Gudang) CariOrganisasi(_ context.Context, kata string, batas int) ([]models.Organisasi, error) {
	k := strings.ToUpper(strings.TrimSpace(kata))
	out := []models.Organisasi{}
	for _, o := range g.Org {
		if k == "" || strings.Contains(strings.ToUpper(o.Nama), k) || strings.Contains(strings.ToUpper(o.IDView), k) {
			out = append(out, o)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Nama < out[j].Nama })
	if len(out) > batas {
		out = out[:batas]
	}
	return out, nil
}

// AmbilOrganisasi - hanya FLAG Org.
func (g *Gudang) AmbilOrganisasi(_ context.Context, _ *db.Tx, id string) (models.Organisasi, error) {
	o, ada := g.Org[id]
	if !ada {
		return models.Organisasi{}, repository.ErrTidakAda
	}
	return o, nil
}

// DaftarGroupBusiness - urut NOTE.
func (g *Gudang) DaftarGroupBusiness(context.Context) ([]models.GroupBusiness, error) {
	out := []models.GroupBusiness{}
	for _, b := range g.GB {
		out = append(out, b)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Note < out[j].Note })
	return out, nil
}

// AmbilGroupBusiness - satu Group Business.
func (g *Gudang) AmbilGroupBusiness(_ context.Context, _ *db.Tx, id string) (models.GroupBusiness, error) {
	b, ada := g.GB[id]
	if !ada {
		return models.GroupBusiness{}, repository.ErrTidakAda
	}
	return b, nil
}
