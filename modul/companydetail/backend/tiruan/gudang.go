// Package tiruan adalah Gudang Company Detail di memori - untuk uji services dan handlers TANPA Oracle.
//
// Perilakunya meniru SQL repository: daftar urut nama tanpa beda huruf lalu ID, cari di nama/ORG ID/NPWP, nomor ORG
// dari sequence (SeqOrg), ganti PIC dan alamat = hapus lalu sisip. Seluruh data uji berawalan `UJI`.
package tiruan

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/companydetail/backend/models"
	"nusantarare/modul/companydetail/backend/repository"
)

// Gudang menyimpan organisasi, PIC, baris alamat, dan pilihan di memori.
type Gudang struct {
	Org    map[string]models.Organisasi
	PIC    map[string][]models.PIC
	Alamat map[string][]models.BarisAlamat
	// Pilihan - isi M_ENUMERASI; Negara - isi NATION.
	Pilihan []models.Pilihan
	Negara  []models.Negara
	// SeqOrg - nilai SEQ_CLIENT_ORG.NEXTVAL berikutnya (migrasi 810 memulainya dari nomor ORG tertinggi + 1).
	SeqOrg int64
	// Disisip dan Diperbarui - jejak tulisan organisasi untuk asersi uji.
	Disisip, Diperbarui []models.Organisasi
	// Lama - isi popup Copy Old; SalinOrgLama membuang ID yang disalin dari sini. GagalSalin - ID yang menjawab galat.
	Lama       []models.OrgLama
	GagalSalin map[string]bool
	// waktu - pencacah jejak tiruan.
	waktu int
}

// Baru menyusun gudang kosong.
func Baru() *Gudang {
	return &Gudang{Org: map[string]models.Organisasi{}, PIC: map[string][]models.PIC{},
		Alamat: map[string][]models.BarisAlamat{}}
}

// Transaksi tiruan - fn(nil).
func Transaksi(_ context.Context, fn func(tx *db.Tx) error) error { return fn(nil) }

func cocok(o models.Organisasi, kueri string) bool {
	if kueri == "" {
		return true
	}
	k := strings.ToUpper(kueri)
	return strings.Contains(strings.ToUpper(o.Nama), k) || strings.Contains(strings.ToUpper(o.IDView), k) ||
		strings.Contains(strings.ToUpper(o.NPWP), k)
}

func (g *Gudang) urut() []models.Organisasi {
	out := make([]models.Organisasi, 0, len(g.Org))
	for _, o := range g.Org {
		out = append(out, o)
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := strings.ToUpper(out[i].Nama), strings.ToUpper(out[j].Nama)
		if a != b {
			return a < b
		}
		return out[i].ID < out[j].ID
	})
	return out
}

// CariOrg - satu halaman organisasi yang cocok dan jumlah seluruhnya.
func (g *Gudang) CariOrg(_ context.Context, kueri string, offset, ukuran int) ([]models.BarisDaftar, int, error) {
	var semua []models.BarisDaftar
	for _, o := range g.urut() {
		if cocok(o, kueri) {
			semua = append(semua, models.BarisDaftar{ID: o.ID, IDView: o.IDView, Nama: o.Nama, Title: o.Title,
				NPWP: o.NPWP, CountryName: o.CountryName, BusinessField: o.BusinessField, ParentName: o.ParentName})
		}
	}
	out := []models.BarisDaftar{}
	for i := offset; i < len(semua) && i < offset+ukuran; i++ {
		out = append(out, semua[i])
	}
	return out, len(semua), nil
}

// AmbilOrg - satu organisasi; tidak ada = repository.ErrTidakAda.
func (g *Gudang) AmbilOrg(_ context.Context, _ *db.Tx, id string) (models.Organisasi, error) {
	o, ada := g.Org[id]
	if !ada {
		return o, repository.ErrTidakAda
	}
	return o, nil
}

// NomorOrgBerikut - SeqOrg lalu dinaikkan, seperti NEXTVAL.
func (g *Gudang) NomorOrgBerikut(context.Context, *db.Tx) (int64, error) {
	n := g.SeqOrg
	g.SeqOrg++
	return n, nil
}

func (g *Gudang) jejak() string {
	g.waktu++
	return "2026-10-04 09:" + strconv.Itoa(10+g.waktu)
}

// SisipOrg - CREATED_AT dan UPDATED_AT = jam tiruan.
func (g *Gudang) SisipOrg(_ context.Context, _ *db.Tx, o models.Organisasi) error {
	o.CreatedAt = g.jejak()
	o.UpdatedAt = o.CreatedAt
	g.Org[o.ID] = o
	g.Disisip = append(g.Disisip, o)
	return nil
}

// PerbaruiOrg - kolom layar dan jejak perubahan; ID, IDVIEW, dan jejak pembuatan tetap.
func (g *Gudang) PerbaruiOrg(_ context.Context, _ *db.Tx, o models.Organisasi) error {
	lama, ada := g.Org[o.ID]
	if !ada {
		return repository.ErrTidakAda
	}
	o.IDView, o.CreatedBy, o.CreatedAt = lama.IDView, lama.CreatedBy, lama.CreatedAt
	o.UpdatedAt = g.jejak()
	g.Org[o.ID] = o
	g.Diperbarui = append(g.Diperbarui, o)
	return nil
}

// DaftarPIC - PIC satu organisasi apa adanya.
func (g *Gudang) DaftarPIC(_ context.Context, _ *db.Tx, clientID string) ([]models.PIC, error) {
	return append([]models.PIC{}, g.PIC[clientID]...), nil
}

// GantiPIC - hapus lalu sisip.
func (g *Gudang) GantiPIC(_ context.Context, _ *db.Tx, clientID string, pic []models.PIC) error {
	g.PIC[clientID] = append([]models.PIC{}, pic...)
	return nil
}

// DaftarBarisAlamat - baris alamat satu organisasi apa adanya.
func (g *Gudang) DaftarBarisAlamat(_ context.Context, _ *db.Tx, clientID string) ([]models.BarisAlamat, error) {
	return append([]models.BarisAlamat{}, g.Alamat[clientID]...), nil
}

// GantiAlamat - hapus lalu sisip; CLIENTID diisi seperti repository.
func (g *Gudang) GantiAlamat(_ context.Context, _ *db.Tx, clientID string, baris []models.BarisAlamat) error {
	out := make([]models.BarisAlamat, 0, len(baris))
	for _, b := range baris {
		b.ClientID = clientID
		out = append(out, b)
	}
	g.Alamat[clientID] = out
	return nil
}

// DaftarPilihan - isi M_ENUMERASI.
func (g *Gudang) DaftarPilihan(context.Context) ([]models.Pilihan, error) {
	return append([]models.Pilihan{}, g.Pilihan...), nil
}

// DaftarNegara - isi NATION.
func (g *Gudang) DaftarNegara(context.Context) ([]models.Negara, error) {
	return append([]models.Negara{}, g.Negara...), nil
}

// DaftarNamaOrg - seluruh organisasi bernama.
func (g *Gudang) DaftarNamaOrg(context.Context) ([]models.BarisDaftar, error) {
	out := []models.BarisDaftar{}
	for _, o := range g.urut() {
		if o.Nama != "" {
			out = append(out, models.BarisDaftar{ID: o.ID, IDView: o.IDView, Nama: o.Nama, Title: o.Title, CountryName: o.CountryName})
		}
	}
	return out, nil
}

// CariInduk - organisasi bernama mirip kueri selain `kecuali`, paling banyak repository.BatasCariInduk.
func (g *Gudang) CariInduk(_ context.Context, kueri, kecuali string) ([]models.BarisDaftar, error) {
	out := []models.BarisDaftar{}
	for _, o := range g.urut() {
		if o.ID == kecuali || !strings.Contains(strings.ToUpper(o.Nama), strings.ToUpper(kueri)) {
			continue
		}
		if len(out) == repository.BatasCariInduk {
			break
		}
		out = append(out, models.BarisDaftar{ID: o.ID, IDView: o.IDView, Nama: o.Nama})
	}
	return out, nil
}

// pil - satu baris M_ENUMERASI tiruan.
func pil(jenis, kode, label string, aktif bool) models.Pilihan {
	return models.Pilihan{Jenis: jenis, Kode: kode, Label: label, Aktif: aktif}
}

// Contoh - gudang berisi data uji: dua organisasi Pega lama (satu dengan PIC dan alamat ber-nomor), pilihan
// M_ENUMERASI ringkas, dan dua negara.
func Contoh() *Gudang {
	g := Baru()
	// Seperti migrasi 810: nomor ORG tertinggi = 120 (dokumen M_CLIENT di atas CLIENT ORG-115) + 1.
	g.SeqOrg = 121
	g.Pilihan = []models.Pilihan{
		pil(models.JenisTitle, "01", "TN.", false),
		pil(models.JenisTitle, "04", "PT.", true),
		pil(models.JenisTitle, "05", "CV.", true),
		pil(models.JenisBidangUsaha, "01", "LIFE INSURANCE", true),
		pil(models.JenisBidangUsaha, "23", "BUILDING / CONSTRUCTION", true),
		pil(models.JenisPosisi, "002", "DIREKTUR", true),
		pil(models.JenisAlamat, "1", "ALAMAT RUMAH", true),
		pil(models.JenisAlamat, "2", "ALAMAT KANTOR", true),
		pil(models.JenisAlamat, "3", "ALAMAT POLIS", false),
		pil(models.JenisTelfax, "2", "FAX", false),
		pil(models.JenisTelfax, "3", "MOBILE PHONE", true),
		pil(models.JenisTelfax, "5", "OFFICE PHONE", true),
		pil(models.JenisTelfax, "6", "EMAIL", false),
		pil(models.JenisKodeHP, "021", "UJI AREA", true),
		pil(models.JenisGender, "1", "Male", true),
		pil(models.JenisGender, "2", "Female", true),
	}
	g.Negara = []models.Negara{
		{ID: "100901", OldID: "001", Nama: "UJI NEGARA SATU", NationInitial: "UJS"},
		{ID: "100902", OldID: "", Nama: "UJI NEGARA BARU", NationInitial: "UJB"},
	}
	induk := models.Organisasi{ID: "ASM-SFAGIS-WORK-ORG ORG-100", IDView: "ORG-100", Nama: "UJI Induk Grup",
		Title: "PT.", Country: "001", CountryName: "UJI NEGARA SATU", BusinessField: "01"}
	anak := models.Organisasi{ID: "ASM-SFAGIS-WORK-ORG ORG-115", IDView: "ORG-115", Nama: "UJI Anak Usaha",
		Title: "PT.", NPWP: "UJI-NPWP-115", Country: "001", CountryName: "UJI NEGARA SATU", BusinessField: "23",
		ParentID: induk.ID, ParentName: induk.Nama}
	g.Org[induk.ID], g.Org[anak.ID] = induk, anak
	g.PIC[anak.ID] = []models.PIC{
		{UserIdentifier: "UJIPEGA", Nama: "UJI Kontak Lama", Position: "Direktur", Gender: "1"},
		{UserIdentifier: "PIC-2", Nama: "UJI Kontak Dua", Position: "Staff", Gender: "2", DateOfBirth: "19900115"},
	}
	asal := models.BarisAlamat{ClientID: anak.ID, Type: "2", Address: "UJI Jalan Satu", City: "UJIKOTA",
		CityName: "UJI KOTA", ZipCode: "10000", PxCreateOperator: "UJI-PEGA", PxCreateDateTime: "20240101T010101.000 GMT"}
	satu, dua := asal, asal
	satu.TelfaxType, satu.TelfaxCode, satu.TelfaxNo = "2", "021", "UJI-FAX-1"
	dua.TelfaxType, dua.TelfaxNo = "5", "UJI-TELP-2"
	g.Alamat[anak.ID] = []models.BarisAlamat{satu, dua,
		{ClientID: anak.ID, Type: "1", Address: "UJI Jalan Dua", PxCreateOperator: "UJI-PEGA"}}
	return g
}

// SiapkanOrgLama - isi popup Copy Old.
func (g *Gudang) SiapkanOrgLama(context.Context) ([]models.OrgLama, error) {
	return append([]models.OrgLama{}, g.Lama...), nil
}

// SalinOrgLama - ID di Lama: disalin (dibuang dari Lama); ID organisasi lain: tidak ada yang ditulis; ID asing:
// repository.ErrTidakAda; GagalSalin: galat basis data tiruan.
func (g *Gudang) SalinOrgLama(_ context.Context, _ *db.Tx, id string, _ repository.ReferensiPindah) (bool, error) {
	if g.GagalSalin[id] {
		return false, errors.New("tiruan: galat basis data")
	}
	for i, o := range g.Lama {
		if o.ID == id {
			g.Lama = append(g.Lama[:i], g.Lama[i+1:]...)
			return true, nil
		}
	}
	if _, ada := g.Org[id]; ada {
		return false, nil
	}
	return false, repository.ErrTidakAda
}

// ReferensiPindah - pemetaan kosong (tiruan Copy Old tidak memetakan title atau negara).
func (g *Gudang) ReferensiPindah(context.Context) (repository.ReferensiPindah, error) {
	return repository.ReferensiPindah{}, nil
}
