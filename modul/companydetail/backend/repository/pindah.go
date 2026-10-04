package repository

// Pindah SEKALI dari dokumen organisasi Pega (`M_CLIENT`) ke tabel datar Company Detail - disetujui work owner
// 03-10-2026 ("SETUJU": JSON hanya dibaca sekali oleh skrip pindah, sesudah itu aplikasi tidak pernah membacanya).
// Alatnya `backend/alat/pindahclient`.
//
// ⛔ BUKAN berkas migrasi SQL: teks jalur dokumen tidak lolos penjaga inti, dan di skema uji tabel warisan hanya
// tiruan. Yang diisi (tidak pernah menimpa nilai yang sudah ada):
//   - organisasi ber-ID `ASM-SFAGIS-WORK-ORG ORG-n` di M_CLIENT yang belum punya baris CLIENT -> baris CLIENT baru;
//   - CLIENT yang sudah ada: PARENT_ID, NOTE, TITLE yang masih kosong;
//   - CLIENT_PICLIST: PIC dokumen yang belum ada (trigger lama menghapus PIC tanpa ASMUserIdentifier);
//   - CLIENT_ADDRESS: nomor Phone and Fax dokumen - nomor pertama mengisi baris alamatnya, nomor berikutnya baris
//     baru beralamat sama (satu baris per nomor).
//
//	-uji      (bawaan) hanya SELECT; nol tulisan; laporan.
//	-jalankan SATU transaksi yang lebih dulu mengunci ketiga tabel; aman diulang (yang sudah terisi dilewati).
//
// ⛔ Laporan AGREGAT: jumlah per jenis, tidak pernah ID kasus, nama, atau nilai.

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/companydetail/backend/models"
)

// ErrPindahDiProduksi - `-jalankan` di lingkungan IS_PEGA_PROD=true.
var ErrPindahDiProduksi = errors.New("repository: moving client documents to the flat tables is refused when IS_PEGA_PROD=true")

// Batas lebar kolom tujuan (BYTE).
const (
	lebarNote     = 4000
	lebarPICNama  = 500
	lebarPICEmail = 500
	lebarPICPos   = 100
	lebarPICTelp  = 100
	lebarPICUID   = 500
	lebarTelfaxNo = 50
	lebarJenis    = 10
	lebarGroup    = 1000
	lebarCountryN = 50
)

// teks menerima nilai dokumen apa pun (teks, angka, boolean, null) sebagai teks.
type teks string

func (t *teks) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case nil:
		*t = ""
	case string:
		*t = teks(x)
	case float64:
		*t = teks(strconv.FormatFloat(x, 'f', -1, 64))
	case bool:
		*t = teks(strconv.FormatBool(x))
	default:
		*t = "" // objek atau larik: bukan nilai medan
	}
	return nil
}

func (t teks) s() string { return strings.TrimSpace(string(t)) }

// dokumenMentah - bagian dokumen organisasi Pega yang dibaca (jalur dari katalog DEV 03-10-2026).
type dokumenMentah struct {
	Name          teks `json:"Name"`
	ParentID      teks `json:"ParentID"`
	ParentName    teks `json:"ParentName"`
	PyDescription teks `json:"pyDescription"`
	DataNB        struct {
		CustomerC struct {
			ASMNPWP          teks `json:"ASMNPWP"`
			ASMBusinessField teks `json:"ASMBusinessField"`
			ASMNationality   teks `json:"ASMNationality"`
			PyTitle          teks `json:"pyTitle"`
		} `json:"Customer_C"`
	} `json:"DataNB"`
	PICList []struct {
		ASMUserIdentifier teks `json:"ASMUserIdentifier"`
		ASMNickName       teks `json:"ASMNickName"`
		ASMPosition       teks `json:"ASMPosition"`
		ASMGender         teks `json:"ASMGender"`
		ASMNoTemplate     teks `json:"ASMNoTemplate"`
		DateOfBirth       teks `json:"DateOfBirth"`
		PyPhoneNumber     teks `json:"pyPhoneNumber"`
	} `json:"PICList"`
	AddressList []struct {
		ASMAddress string `json:"ASMAddress"`
		ASMTelfax  []struct {
			TelfaxType   teks `json:"TelfaxType"`
			TelFaxCode   teks `json:"TelFaxCode"`
			TelfaxNumber teks `json:"TelfaxNumber"`
		} `json:"ASMTelfax"`
	} `json:"AddressList"`
}

// DokumenOrg - isi satu dokumen organisasi yang dipakai pindah.
type DokumenOrg struct {
	ID, Nama, NPWP, BusinessField, Negara, Title, ParentID, ParentName, Note string
	PIC                                                                      []models.PIC
	// Telfax - per ASMADDRESS (teks apa adanya, kunci baris CLIENT_ADDRESS), nomor berurutan.
	Telfax     map[string][]models.Telfax
	urutAlamat []string
}

// UraiDokumen mengurai satu dokumen organisasi.
func UraiDokumen(id string, isi []byte) (DokumenOrg, error) {
	var m dokumenMentah
	if err := json.Unmarshal(isi, &m); err != nil {
		return DokumenOrg{}, err
	}
	c := m.DataNB.CustomerC
	d := DokumenOrg{ID: id, Nama: m.Name.s(), NPWP: c.ASMNPWP.s(), BusinessField: c.ASMBusinessField.s(),
		Negara: c.ASMNationality.s(), Title: c.PyTitle.s(), ParentID: m.ParentID.s(), ParentName: m.ParentName.s(),
		Note: m.PyDescription.s(), Telfax: map[string][]models.Telfax{}}
	for _, p := range m.PICList {
		d.PIC = append(d.PIC, models.PIC{UserIdentifier: p.ASMUserIdentifier.s(), Nama: p.ASMNickName.s(),
			Position: p.ASMPosition.s(), Gender: p.ASMGender.s(), Email: p.ASMNoTemplate.s(),
			DateOfBirth: p.DateOfBirth.s(), Phone: p.PyPhoneNumber.s()})
	}
	for _, a := range m.AddressList {
		var nomor []models.Telfax
		for _, t := range a.ASMTelfax {
			nomor = append(nomor, models.Telfax{Jenis: t.TelfaxType.s(), Code: t.TelFaxCode.s(), No: t.TelfaxNumber.s()})
		}
		if len(nomor) > 0 {
			if _, ada := d.Telfax[a.ASMAddress]; !ada {
				d.urutAlamat = append(d.urutAlamat, a.ASMAddress)
			}
			d.Telfax[a.ASMAddress] = append(d.Telfax[a.ASMAddress], nomor...)
		}
	}
	return d, nil
}

// KeadaanFlat - isi tabel datar sebelum pindah.
type KeadaanFlat struct {
	// SemuaID - seluruh CLIENT.ID (Org, Contact, dan lainnya).
	SemuaID map[string]bool
	// Org - baris CLIENT Org: kolom yang boleh diisi pindah.
	Org    map[string]models.Organisasi
	PIC    map[string][]models.PIC
	Alamat map[string][]models.BarisAlamat
}

// ReferensiPindah - pemetaan title dan negara.
type ReferensiPindah struct {
	// TitleKode - kode title -> LABEL; TitleLabel - LABEL yang sah.
	TitleKode  map[string]string
	TitleLabel map[string]bool
	// Negara - kode CLIENT.COUNTRY (OLDID, atau ID) -> negara.
	Negara map[string]models.Negara
}

// IsiOrg - kolom kosong CLIENT yang diisi (kosong = tidak diisi).
type IsiOrg struct{ ID, ParentID, Note, Title string }

// UbahTelfax - baris alamat yang mendapat nomor pertamanya.
type UbahTelfax struct {
	ClientID, Address string
	Telfax            models.Telfax
}

// RencanaPindah - seluruh tulisan satu putaran.
type RencanaPindah struct {
	OrgBaru     []models.Organisasi
	OrgIsi      []IsiOrg
	PICBaru     map[string][]models.PIC
	TelfaxUbah  []UbahTelfax
	TelfaxSisip []models.BarisAlamat
}

// LaporanPindah - laporan agregat satu putaran.
type LaporanPindah struct {
	Mode    string
	Dokumen int
	Jumlah  map[string]int
}

func (l *LaporanPindah) tambah(jenis string, n int) {
	if l.Jumlah == nil {
		l.Jumlah = map[string]int{}
	}
	l.Jumlah[jenis] += n
}

// Teks - laporan untuk konsol, urut jenis.
func (l LaporanPindah) Teks() string {
	var b strings.Builder
	fmt.Fprintf(&b, "pindahclient %s: %d dokumen organisasi dibaca\n", l.Mode, l.Dokumen)
	jenis := make([]string, 0, len(l.Jumlah))
	for k := range l.Jumlah {
		jenis = append(jenis, k)
	}
	sort.Strings(jenis)
	for _, k := range jenis {
		fmt.Fprintf(&b, "  %-44s %d\n", k, l.Jumlah[k])
	}
	return b.String()
}

var (
	polaIDOrg  = regexp.MustCompile(polaID)
	polaPICNo  = regexp.MustCompile(`^` + models.AwalanPIC + `([0-9]+)$`)
	polaTgl8   = regexp.MustCompile(`^[0-9]{8}$`)
	awalanIDOk = models.AwalanID
)

// cariInduk - ParentID dokumen sebagai CLIENT.ID: apa adanya, atau `ORG-n` diberi awalan ID.
func cariInduk(parent string, ada func(string) bool) string {
	if parent == "" {
		return ""
	}
	if ada(parent) {
		return parent
	}
	if ada(awalanIDOk + parent) {
		return awalanIDOk + parent
	}
	return ""
}

func (r ReferensiPindah) title(nilai string) (string, bool) {
	if nilai == "" {
		return "", true
	}
	if l, ok := r.TitleKode[nilai]; ok {
		return l, true
	}
	if r.TitleLabel[nilai] {
		return nilai, true
	}
	return "", false
}

// SusunRencana menyusun tulisan pindah dari dokumen dan keadaan tabel datar - murni, diuji tanpa Oracle.
func SusunRencana(dok []DokumenOrg, k KeadaanFlat, r ReferensiPindah, lap *LaporanPindah) RencanaPindah {
	rc := RencanaPindah{PICBaru: map[string][]models.PIC{}}
	baru := map[string]bool{}
	for _, d := range dok {
		if polaIDOrg.MatchString(d.ID) && !k.SemuaID[d.ID] {
			baru[d.ID] = true
		}
	}
	ada := func(id string) bool { return k.Org[id].ID != "" || baru[id] }
	for _, d := range dok {
		if !polaIDOrg.MatchString(d.ID) {
			lap.tambah("dokumen ID bukan ORG-n (dilewati)", 1)
			continue
		}
		title, sah := r.title(d.Title)
		if !sah {
			lap.tambah("title dokumen tidak dikenal (tidak diisi)", 1)
		}
		induk := cariInduk(d.ParentID, ada)
		if d.ParentID != "" && induk == "" {
			lap.tambah("ParentID tanpa baris CLIENT (tidak diisi)", 1)
		}
		note := d.Note
		if len(note) > lebarNote {
			note = ""
			lap.tambah("Note lebih dari 4000 byte (tidak diisi)", 1)
		}
		switch {
		case baru[d.ID]:
			o := models.Organisasi{ID: d.ID, IDView: strings.TrimPrefix(d.ID, awalanIDOk), Nama: d.Nama, NPWP: d.NPWP,
				BusinessField: d.BusinessField, Title: title, ParentID: induk, Note: note}
			if len(d.ParentName) <= lebarGroup {
				o.ParentName = d.ParentName
			}
			if n, ok := r.Negara[d.Negara]; ok && len(n.Nama) <= lebarCountryN {
				o.Country, o.CountryName = n.Kode(), n.Nama
			} else if d.Negara != "" {
				o.Country = d.Negara
				lap.tambah("organisasi baru: negara tidak ada di NATION", 1)
			}
			rc.OrgBaru = append(rc.OrgBaru, o)
			lap.tambah("CLIENT baru (dokumen tanpa baris CLIENT)", 1)
		case k.Org[d.ID].ID != "":
			lama := k.Org[d.ID]
			isi := IsiOrg{ID: d.ID}
			if lama.ParentID == "" && induk != "" {
				isi.ParentID = induk
				lap.tambah("CLIENT.PARENT_ID diisi", 1)
			}
			if lama.Note == "" && note != "" {
				isi.Note = note
				lap.tambah("CLIENT.NOTE diisi", 1)
			}
			if lama.Title == "" && title != "" {
				isi.Title = title
				lap.tambah("CLIENT.TITLE diisi", 1)
			}
			if isi.ParentID != "" || isi.Note != "" || isi.Title != "" {
				rc.OrgIsi = append(rc.OrgIsi, isi)
			}
		default:
			lap.tambah("dokumen ber-ID CLIENT bukan Org (dilewati)", 1)
			continue
		}
		susunPICPindah(d, k.PIC[d.ID], &rc, lap)
		susunTelfaxPindah(d, k.Alamat[d.ID], &rc, lap)
	}
	return rc
}

func susunPICPindah(d DokumenOrg, lama []models.PIC, rc *RencanaPindah, lap *LaporanPindah) {
	if len(d.PIC) == 0 {
		return
	}
	uid, nama := map[string]bool{}, map[string]bool{}
	maks := 0
	for _, p := range lama {
		uid[p.UserIdentifier] = true
		nama[strings.ToUpper(p.Nama)] = true
		if m := polaPICNo.FindStringSubmatch(p.UserIdentifier); m != nil {
			if n, _ := strconv.Atoi(m[1]); n > maks {
				maks = n
			}
		}
	}
	for _, p := range d.PIC {
		switch {
		case p.Nama == "":
			lap.tambah("PIC dokumen tanpa nama (dilewati)", 1)
			continue
		case (p.UserIdentifier != "" && uid[p.UserIdentifier]) || (p.UserIdentifier == "" && nama[strings.ToUpper(p.Nama)]):
			lap.tambah("PIC sudah ada (dilewati)", 1)
			continue
		case len(p.Nama) > lebarPICNama || len(p.UserIdentifier) > lebarPICUID:
			lap.tambah("PIC nama/identifier terlalu panjang (dilewati)", 1)
			continue
		}
		if p.UserIdentifier == "" {
			maks++
			p.UserIdentifier = models.AwalanPIC + strconv.Itoa(maks)
		}
		uid[p.UserIdentifier], nama[strings.ToUpper(p.Nama)] = true, true
		if p.Gender != "1" && p.Gender != "2" {
			p.Gender = ""
		}
		if len(p.Position) > lebarPICPos {
			p.Position = ""
			lap.tambah("PIC Position terlalu panjang (dikosongkan)", 1)
		}
		if len(p.Email) > lebarPICEmail {
			p.Email = ""
			lap.tambah("PIC Email terlalu panjang (dikosongkan)", 1)
		}
		if len(p.Phone) > lebarPICTelp {
			p.Phone = ""
			lap.tambah("PIC Phone terlalu panjang (dikosongkan)", 1)
		}
		if p.DateOfBirth != "" && !polaTgl8.MatchString(p.DateOfBirth) {
			if t, err := time.Parse("2006-01-02", p.DateOfBirth); err == nil {
				p.DateOfBirth = t.Format("20060102")
			}
		}
		rc.PICBaru[d.ID] = append(rc.PICBaru[d.ID], p)
		lap.tambah("CLIENT_PICLIST baris baru", 1)
	}
}

func susunTelfaxPindah(d DokumenOrg, baris []models.BarisAlamat, rc *RencanaPindah, lap *LaporanPindah) {
	for _, alamat := range d.urutAlamat {
		var milik []models.BarisAlamat
		for _, b := range baris {
			if b.Address == alamat {
				milik = append(milik, b)
			}
		}
		if len(milik) == 0 {
			lap.tambah("alamat bernomor tanpa baris CLIENT_ADDRESS (dilewati)", 1)
			continue
		}
		sudah := false
		for _, b := range milik {
			sudah = sudah || b.TelfaxType != "" || b.TelfaxNo != ""
		}
		if sudah {
			lap.tambah("alamat yang nomornya sudah dipindah (dilewati)", 1)
			continue
		}
		pertama := true
		for _, t := range d.Telfax[alamat] {
			if t.No == "" || len(t.No) > lebarTelfaxNo || len(t.Jenis) > lebarJenis || len(t.Code) > lebarJenis {
				lap.tambah("nomor kosong atau terlalu panjang (dilewati)", 1)
				continue
			}
			if pertama {
				rc.TelfaxUbah = append(rc.TelfaxUbah, UbahTelfax{ClientID: d.ID, Address: alamat, Telfax: t})
				lap.tambah("CLIENT_ADDRESS nomor pertama diisi", 1)
				pertama = false
				continue
			}
			b := milik[0]
			b.TelfaxType, b.TelfaxCode, b.TelfaxNo = t.Jenis, t.Code, t.No
			rc.TelfaxSisip = append(rc.TelfaxSisip, b)
			lap.tambah("CLIENT_ADDRESS baris nomor berikutnya", 1)
		}
	}
}

// --- Oracle ---

func sqlDokumenOrg(t string) string {
	return fmt.Sprintf(`SELECT ID, JSONDATA FROM %s WHERE REGEXP_LIKE(ID, :1) ORDER BY ID`, t)
}

func sqlSemuaClient(t string) string {
	return fmt.Sprintf(`SELECT ID, FLAG, PARENT_ID, NOTE, TITLE FROM %s`, t)
}

func sqlSemuaPIC(t string) string {
	return fmt.Sprintf(`SELECT CLIENTID, USERIDENTIFIER, NICKNAME FROM %s`, t)
}

func sqlSemuaAlamat(t string) string { return fmt.Sprintf(`SELECT %s FROM %s`, kolomAlamat, t) }

// sqlIsiOrg - NVL: nilai yang sudah ada TIDAK pernah ditimpa.
func sqlIsiOrg(t string) string {
	return fmt.Sprintf(`UPDATE %s SET PARENT_ID = NVL(PARENT_ID, :1), NOTE = NVL(NOTE, :2), TITLE = NVL(TITLE, :3)
	  WHERE ID = :4 AND FLAG = :5`, t)
}

// sqlUbahTelfax - nomor pertama mengisi SATU baris alamat yang belum bernomor.
func sqlUbahTelfax(t string) string {
	return fmt.Sprintf(`UPDATE %s SET TELFAX_TYPE = :1, TELFAX_CODE = :2, TELFAX_NO = :3
	  WHERE CLIENTID = :4 AND ASMADDRESS = :5 AND TELFAX_TYPE IS NULL AND TELFAX_NO IS NULL AND ROWNUM = 1`, t)
}

func sqlKunciTabel(t string) string { return fmt.Sprintf(`LOCK TABLE %s IN EXCLUSIVE MODE`, t) }

// bacaKeadaan membaca tabel datar dan referensi pindah.
func (g *Gudang) bacaKeadaan(ctx context.Context, j penjalan) (KeadaanFlat, ReferensiPindah, error) {
	k := KeadaanFlat{SemuaID: map[string]bool{}, Org: map[string]models.Organisasi{}, PIC: map[string][]models.PIC{},
		Alamat: map[string][]models.BarisAlamat{}}
	r := ReferensiPindah{TitleKode: map[string]string{}, TitleLabel: map[string]bool{}, Negara: map[string]models.Negara{}}
	tc, err := g.nama(TabelClient)
	if err != nil {
		return k, r, err
	}
	client, err := daftar(ctx, j, sqlSemuaClient(tc), func(p pemindai) ([]string, error) { return pindaiTeks(p, 5) })
	if err != nil {
		return k, r, err
	}
	for _, v := range client {
		k.SemuaID[v[0]] = true
		if v[1] == models.FlagOrg {
			k.Org[v[0]] = models.Organisasi{ID: v[0], ParentID: v[2], Note: v[3], Title: v[4]}
		}
	}
	tp, err := g.nama(TabelPIC)
	if err != nil {
		return k, r, err
	}
	pic, err := daftar(ctx, j, sqlSemuaPIC(tp), func(p pemindai) ([]string, error) { return pindaiTeks(p, 3) })
	if err != nil {
		return k, r, err
	}
	for _, v := range pic {
		k.PIC[v[0]] = append(k.PIC[v[0]], models.PIC{UserIdentifier: v[1], Nama: v[2]})
	}
	ta, err := g.nama(TabelAlamat)
	if err != nil {
		return k, r, err
	}
	alamat, err := daftar(ctx, j, sqlSemuaAlamat(ta), pindaiAlamat)
	if err != nil {
		return k, r, err
	}
	for _, b := range alamat {
		k.Alamat[b.ClientID] = append(k.Alamat[b.ClientID], b)
	}
	pil, err := g.DaftarPilihan(ctx)
	if err != nil {
		return k, r, err
	}
	for _, p := range pil {
		if p.Jenis == models.JenisTitle {
			r.TitleKode[p.Kode] = p.Label
			r.TitleLabel[p.Label] = true
		}
	}
	neg, err := g.DaftarNegara(ctx)
	if err != nil {
		return k, r, err
	}
	for _, n := range neg {
		r.Negara[n.Kode()] = n
	}
	return k, r, nil
}

// bacaDokumen mengurai setiap dokumen organisasi; dokumen yang tidak terurai dihitung, tidak menghentikan putaran.
func (g *Gudang) bacaDokumen(ctx context.Context, j penjalan, lap *LaporanPindah) ([]DokumenOrg, error) {
	td, err := g.nama(TabelDokumen)
	if err != nil {
		return nil, err
	}
	q := sqlDokumenOrg(td)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := j.QueryContext(ctx, q, polaID)
	if err != nil {
		return nil, bungkus(err, "membaca dokumen organisasi")
	}
	defer func() { _ = rows.Close() }()
	var out []DokumenOrg
	for rows.Next() {
		var id, isi sql.NullString
		if err := rows.Scan(&id, &isi); err != nil {
			return nil, bungkus(err, "memindai dokumen organisasi")
		}
		lap.Dokumen++
		d, err := UraiDokumen(id.String, []byte(isi.String))
		if err != nil {
			lap.tambah("dokumen tidak terurai (dilewati)", 1)
			continue
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// PindahClient menjalankan satu putaran pindah. `tulis` false = uji kering (SELECT saja).
func (g *Gudang) PindahClient(ctx context.Context, tulis bool) (LaporanPindah, error) {
	lap := LaporanPindah{Mode: "-uji"}
	if !tulis {
		dok, err := g.bacaDokumen(ctx, g.db, &lap)
		if err != nil {
			return lap, err
		}
		k, r, err := g.bacaKeadaan(ctx, g.db)
		if err != nil {
			return lap, err
		}
		SusunRencana(dok, k, r, &lap)
		return lap, nil
	}
	lap.Mode = "-jalankan"
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return lap, err
	}
	selesai := false
	defer func() {
		if !selesai {
			_ = tx.Rollback()
		}
	}()
	for _, t := range []string{TabelClient, TabelPIC, TabelAlamat} {
		nama, err := g.nama(t)
		if err != nil {
			return lap, err
		}
		if _, err := jalankan(ctx, tx, sqlKunciTabel(nama), "mengunci "+t); err != nil {
			return lap, err
		}
	}
	dok, err := g.bacaDokumen(ctx, tx, &lap)
	if err != nil {
		return lap, err
	}
	k, r, err := g.bacaKeadaan(ctx, tx)
	if err != nil {
		return lap, err
	}
	rc := SusunRencana(dok, k, r, &lap)
	if err := g.terapkan(ctx, tx, rc); err != nil {
		return lap, err
	}
	if err := tx.Commit(); err != nil {
		return lap, err
	}
	selesai = true
	return lap, nil
}

func (g *Gudang) terapkan(ctx context.Context, tx *db.Tx, rc RencanaPindah) error {
	tc, err := g.nama(TabelClient)
	if err != nil {
		return err
	}
	for _, o := range rc.OrgBaru {
		if _, err := jalankan(ctx, tx, sqlSisipOrg(tc), "pindah: menyisipkan organisasi", NilaiSisipOrg(o)...); err != nil {
			return err
		}
	}
	n := db.KosongJadiNil
	for _, i := range rc.OrgIsi {
		if _, err := jalankan(ctx, tx, sqlIsiOrg(tc), "pindah: mengisi organisasi",
			n(i.ParentID), n(i.Note), n(i.Title), i.ID, models.FlagOrg); err != nil {
			return err
		}
	}
	tp, err := g.nama(TabelPIC)
	if err != nil {
		return err
	}
	klien := make([]string, 0, len(rc.PICBaru))
	for c := range rc.PICBaru {
		klien = append(klien, c)
	}
	sort.Strings(klien)
	for _, c := range klien {
		for _, p := range rc.PICBaru[c] {
			if _, err := jalankan(ctx, tx, sqlSisipPIC(tp), "pindah: menyisipkan PIC", NilaiSisipPIC(c, p)...); err != nil {
				return err
			}
		}
	}
	ta, err := g.nama(TabelAlamat)
	if err != nil {
		return err
	}
	for _, u := range rc.TelfaxUbah {
		if _, err := jalankan(ctx, tx, sqlUbahTelfax(ta), "pindah: mengisi nomor alamat", n(u.Telfax.Jenis),
			n(u.Telfax.Code), u.Telfax.No, u.ClientID, u.Address); err != nil {
			return err
		}
	}
	for _, b := range rc.TelfaxSisip {
		if _, err := jalankan(ctx, tx, sqlSisipAlamat(ta), "pindah: menyisipkan nomor alamat", NilaiSisipAlamat(b)...); err != nil {
			return err
		}
	}
	return nil
}
