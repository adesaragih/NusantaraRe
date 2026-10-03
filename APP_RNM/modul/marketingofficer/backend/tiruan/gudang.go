// Package tiruan adalah Gudang Marketing Officer di memori - untuk uji services dan handlers TANPA Oracle.
//
// Perilakunya meniru SQL repository: ID baru `1` + 7 digit dari nomor urut, AKSES_LOGIN dibandingkan tanpa beda
// huruf, cabang dan akun dicari persis. Seluruh data uji berawalan `UJI-`.
package tiruan

import (
	"context"
	"sort"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/marketingofficer/backend/models"
	"nusantarare/modul/marketingofficer/backend/repository"
)

// Gudang menyimpan baris MO, akun, dan cabang di memori.
type Gudang struct {
	MO     map[string]models.MarketingOfficer
	Akun   map[string]models.Akun
	Cabang map[string]models.Cabang
	// Nomor - nilai CURRENCY_SEQ berikutnya.
	Nomor int64
	// Disisip dan Diperbarui - jejak tulisan untuk asersi uji.
	Disisip, Diperbarui []models.MarketingOfficer
}

// Baru menyusun gudang kosong.
func Baru() *Gudang {
	return &Gudang{MO: map[string]models.MarketingOfficer{}, Akun: map[string]models.Akun{},
		Cabang: map[string]models.Cabang{}, Nomor: 201}
}

// Transaksi tiruan - fn(nil).
func Transaksi(_ context.Context, fn func(tx *db.Tx) error) error { return fn(nil) }

// DaftarMO - urut MOSTATUS, nama, ID seperti SQL.
func (g *Gudang) DaftarMO(context.Context) ([]models.MarketingOfficer, error) {
	out := []models.MarketingOfficer{}
	for _, m := range g.MO {
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if a.MOStatus != b.MOStatus {
			return a.MOStatus < b.MOStatus
		}
		if strings.ToUpper(a.ClientName) != strings.ToUpper(b.ClientName) {
			return strings.ToUpper(a.ClientName) < strings.ToUpper(b.ClientName)
		}
		return a.ID < b.ID
	})
	return out, nil
}

// AmbilMO membaca satu baris.
func (g *Gudang) AmbilMO(_ context.Context, _ *db.Tx, id string) (models.MarketingOfficer, error) {
	m, ada := g.MO[id]
	if !ada {
		return m, repository.ErrTidakAda
	}
	return m, nil
}

// SisipMO menyisipkan baris dengan ID dari nomor urut; TANGGAL tiruan.
func (g *Gudang) SisipMO(_ context.Context, _ *db.Tx, m models.MarketingOfficer) (string, error) {
	id, err := repository.FormatIDMO(g.Nomor)
	if err != nil {
		return "", err
	}
	g.Nomor++
	m.ID, m.Tanggal = id, "2026-10-03 09:00"
	g.MO[id] = m
	g.Disisip = append(g.Disisip, m)
	return id, nil
}

// PerbaruiMO menulis ulang kolom yang boleh berubah - CLIENTID dan CLIENTNAME TIDAK, seperti SQL.
func (g *Gudang) PerbaruiMO(_ context.Context, _ *db.Tx, m models.MarketingOfficer) error {
	lama, ada := g.MO[m.ID]
	if !ada {
		return repository.ErrTidakAda
	}
	m.ClientID, m.ClientName, m.BranchStatus, m.Tanggal = lama.ClientID, lama.ClientName, lama.BranchStatus, "2026-10-03 10:00"
	g.MO[m.ID] = m
	g.Diperbarui = append(g.Diperbarui, m)
	return nil
}

func (g *Gudang) aktifBila(cocok func(models.MarketingOfficer) bool, kecuali string) []string {
	out := []string{}
	for id, m := range g.MO {
		if m.MOStatus == models.StatusAktif && id != kecuali && cocok(m) {
			out = append(out, id)
		}
	}
	sort.Strings(out)
	return out
}

// AktifDenganClientID - baris aktif lain ber-Marketing Code itu.
func (g *Gudang) AktifDenganClientID(_ context.Context, _ *db.Tx, clientID, kecuali string) ([]string, error) {
	return g.aktifBila(func(m models.MarketingOfficer) bool { return m.ClientID == clientID }, kecuali), nil
}

// AktifDenganAkses - baris aktif lain ber-akun itu, tanpa beda huruf.
func (g *Gudang) AktifDenganAkses(_ context.Context, _ *db.Tx, login, kecuali string) ([]string, error) {
	return g.aktifBila(func(m models.MarketingOfficer) bool { return strings.EqualFold(m.AksesLogin, login) }, kecuali), nil
}

// ClientIDDariAkses - baris aktif lebih dulu, lalu ID terbesar.
func (g *Gudang) ClientIDDariAkses(_ context.Context, _ *db.Tx, login string) (string, bool, error) {
	var calon []models.MarketingOfficer
	for _, m := range g.MO {
		if strings.EqualFold(m.AksesLogin, login) && m.ClientID != "" {
			calon = append(calon, m)
		}
	}
	if len(calon) == 0 {
		return "", false, nil
	}
	sort.Slice(calon, func(i, j int) bool {
		if calon[i].MOStatus != calon[j].MOStatus {
			return calon[i].MOStatus < calon[j].MOStatus
		}
		return calon[i].ID > calon[j].ID
	})
	return calon[0].ClientID, true, nil
}

// DaftarAkun - seluruh akun, urut nama.
func (g *Gudang) DaftarAkun(context.Context) ([]models.Akun, error) {
	out := []models.Akun{}
	for _, a := range g.Akun {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Nama < out[j].Nama })
	return out, nil
}

// AmbilAkun - persis LOGIN_ID.
func (g *Gudang) AmbilAkun(_ context.Context, _ *db.Tx, login string) (models.Akun, error) {
	a, ada := g.Akun[login]
	if !ada {
		return a, repository.ErrTidakAda
	}
	return a, nil
}

// DaftarCabang - cabang aktif, urut nama.
func (g *Gudang) DaftarCabang(context.Context) ([]models.Cabang, error) {
	out := []models.Cabang{}
	for _, c := range g.Cabang {
		if c.Aktif {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Nama < out[j].Nama })
	return out, nil
}

// AmbilCabang - persis ID, aktif atau tidak.
func (g *Gudang) AmbilCabang(_ context.Context, _ *db.Tx, id string) (models.Cabang, error) {
	c, ada := g.Cabang[id]
	if !ada {
		return c, repository.ErrTidakAda
	}
	return c, nil
}

// Contoh mengisi gudang dengan data UJI: leader aktif 10000101, satu leader nonaktif, satu baris lama Pega,
// tiga akun aktif, satu akun nonaktif, dan cabang induk UJI-P00 dengan dua anak (satu nonaktif).
func Contoh() *Gudang {
	g := Baru()
	for _, a := range []models.Akun{
		{LoginID: "UJI-LEAD01", Nama: "UJI Leader Satu", ContactID: "CON-1001", Email: "uji.lead01@nusantara.example", Aktif: true},
		{LoginID: "UJI-MKT01", Nama: "UJI Marketing Satu", ContactID: "CON-1002", Email: "uji.mkt01@nusantara.example", Aktif: true},
		{LoginID: "UJI-MKT02", Nama: "UJI Marketing Dua", ContactID: "CON-1003", Email: "uji.mkt02@nusantara.example", Aktif: true},
		{LoginID: "UJI-OFF", Nama: "UJI Nonaktif", ContactID: "CON-1004", Aktif: false},
	} {
		g.Akun[a.LoginID] = a
	}
	for _, c := range []models.Cabang{
		{ID: "UJI-P00", Nama: "UJI Kantor Pusat", Induk: "UJI-P00", Aktif: true},
		{ID: "UJI-B01", Nama: "UJI Cabang Satu", Induk: "UJI-P00", KanwilGroup: "1", Aktif: true},
		{ID: "UJI-B02", Nama: "UJI Cabang Dua", Induk: "UJI-P00", KanwilGroup: "2", Aktif: true},
		{ID: "UJI-B09", Nama: "UJI Cabang Tutup", Induk: "UJI-P00", KanwilGroup: "9", Aktif: false},
	} {
		g.Cabang[c.ID] = c
	}
	for _, m := range []models.MarketingOfficer{
		{ID: "10000101", ClientID: "CON-1001", ClientName: "UJI Leader Satu", ClientID2: models.NilaiLeader,
			MOLeader: "UJI Leader Satu", MOStatus: models.StatusAktif, BranchParent: "UJI-P00", BranchDetailID: "UJI-B01",
			BranchDetailName: "UJI Cabang Satu", TeamGroup: "1", AksesLogin: "UJI-LEAD01"},
		{ID: "10000102", ClientID: "UJI-KONTAK-9", ClientName: "UJI Leader Lama", ClientID2: models.NilaiLeader,
			MOLeader: "UJI Leader Lama", MOStatus: models.StatusNonaktif},
		// Baris lama Pega: kode kontak, AKSES_LOGIN Operator ID yang tidak ada di M_LOGIN_GO.
		{ID: "10000103", ClientID: "UJI-KONTAK-3", ClientName: "UJI Pega Lama", ClientID2: "10000101",
			MOLeader: "UJI Leader Satu", MOStatus: models.StatusAktif, BranchParent: "UJI-P00", BranchDetailID: "UJI-B02",
			BranchDetailName: "UJI Cabang Dua Lama", TeamGroup: "2", AksesLogin: "UJIOPERATORLAMA"},
	} {
		g.MO[m.ID] = m
	}
	return g
}
