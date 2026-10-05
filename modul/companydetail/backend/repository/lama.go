package repository

// Copy Old (perintah work owner 04-10-2026, `models/lama.go`): popup berisi organisasi dokumen `M_CLIENT` yang datanya
// belum pindah, dan salinan per organisasi. Aturan salin SAMA dengan alat pindah (`SusunRencana`, `terapkan`).

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/companydetail/backend/models"
)

// sqlDokumenPerlu - dokumen organisasi yang MUNGKIN masih punya data untuk dipindah: tanpa baris CLIENT, ber-PIC,
// ber-nomor, atau ber-ParentID / pyDescription yang kolom CLIENT-nya masih kosong. Penyaring kasar di Oracle (isi
// dokumen yang dibaca Go adalah biaya terbesar popup - DEV 04-10-2026: 1.009 baris, 16,6 detik); penentunya
// `SusunRencana`.
func sqlDokumenPerlu(dokumen, client string) string {
	return fmt.Sprintf(`SELECT m.ID, m.JSONDATA FROM %s m
	  LEFT JOIN %s c ON c.ID = m.ID
	  WHERE REGEXP_LIKE(m.ID, :1)
	    AND (c.ID IS NULL
	         OR JSON_EXISTS(m.JSONDATA, '$.PICList[0]')
	         OR JSON_EXISTS(m.JSONDATA, '$.AddressList[*].ASMTelfax[0]')
	         OR (c.PARENT_ID IS NULL AND JSON_VALUE(m.JSONDATA, '$.ParentID') IS NOT NULL)
	         OR (c.NOTE IS NULL AND JSON_VALUE(m.JSONDATA, '$.pyDescription' RETURNING VARCHAR2(4000) NULL ON ERROR) IS NOT NULL))
	  ORDER BY m.ID`, dokumen, client)
}

func sqlDokumenSatu(t string) string {
	return fmt.Sprintf(`SELECT ID, JSONDATA FROM %s WHERE ID = :1`, t)
}

// sqlClientSatu - baris CLIENT organisasi itu dan calon induknya (ParentID apa adanya / berawalan ID).
func sqlClientSatu(t string) string {
	return fmt.Sprintf(`SELECT ID, FLAG, PARENT_ID, NOTE, TITLE FROM %s WHERE ID IN (:1, :2, :3)`, t)
}

func sqlPICSatu(t string) string {
	return fmt.Sprintf(`SELECT CLIENTID, USERIDENTIFIER, NICKNAME FROM %s WHERE CLIENTID = :1`, t)
}

// RingkasRencana - baris popup Copy Old dari rencana satu dokumen; `ada` false = tidak ada yang ditulis.
func RingkasRencana(d DokumenOrg, rc RencanaPindah) (models.OrgLama, bool) {
	o := models.OrgLama{ID: d.ID, IDView: strings.TrimPrefix(d.ID, awalanIDOk), Nama: d.Nama, Isi: []string{}}
	o.Baru = len(rc.OrgBaru) > 0
	for _, i := range rc.OrgIsi {
		if i.ParentID != "" {
			o.Isi = append(o.Isi, "PARENT_ID")
		}
		if i.Note != "" {
			o.Isi = append(o.Isi, "NOTE")
		}
		if i.Title != "" {
			o.Isi = append(o.Isi, "TITLE")
		}
	}
	for _, p := range rc.PICBaru {
		o.PIC += len(p)
	}
	o.Nomor = len(rc.TelfaxUbah) + len(rc.TelfaxSisip)
	// TITLE saja tidak menjadikan baris popup: 12 ribuan dokumen ber-title, alat pindah yang mengisinya. Bila baris
	// itu disalin karena alasan lain, title ikut diisi.
	adaLain := o.Baru || o.PIC > 0 || o.Nomor > 0
	for _, k := range o.Isi {
		adaLain = adaLain || k != "TITLE"
	}
	return o, adaLain
}

// ReferensiPindah - title dan negara pemetaan salinan; dibaca sekali per Process Copy.
func (g *Gudang) ReferensiPindah(ctx context.Context) (ReferensiPindah, error) {
	r := ReferensiPindah{TitleKode: map[string]string{}, TitleLabel: map[string]bool{}, Negara: map[string]models.Negara{}}
	pil, err := g.DaftarPilihan(ctx)
	if err != nil {
		return r, err
	}
	for _, p := range pil {
		if p.Jenis == models.JenisTitle {
			r.TitleKode[p.Kode] = p.Label
			r.TitleLabel[p.Label] = true
		}
	}
	neg, err := g.DaftarNegara(ctx)
	if err != nil {
		return r, err
	}
	for _, n := range neg {
		r.Negara[n.Kode()] = n
	}
	return r, nil
}

// SiapkanOrgLama - isi popup Copy Old. SELECT saja.
func (g *Gudang) SiapkanOrgLama(ctx context.Context) ([]models.OrgLama, error) {
	td, err := g.nama(TabelDokumen)
	if err != nil {
		return nil, err
	}
	tc, err := g.nama(TabelClient)
	if err != nil {
		return nil, err
	}
	q := sqlDokumenPerlu(td, tc)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, polaID)
	if err != nil {
		return nil, bungkus(err, "membaca dokumen organisasi")
	}
	defer func() { _ = rows.Close() }()
	var dok []DokumenOrg
	for rows.Next() {
		var id, isi sql.NullString
		if err := rows.Scan(&id, &isi); err != nil {
			return nil, bungkus(err, "memindai dokumen organisasi")
		}
		if d, err := UraiDokumen(id.String, []byte(isi.String)); err == nil {
			dok = append(dok, d)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, bungkus(err, "membaca dokumen organisasi")
	}
	k, r, err := g.bacaKeadaan(ctx, g.db)
	if err != nil {
		return nil, err
	}
	out := []models.OrgLama{}
	for _, d := range dok {
		var lap LaporanPindah
		if o, ada := RingkasRencana(d, SusunRencana([]DokumenOrg{d}, k, r, &lap)); ada {
			out = append(out, o)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

// SalinOrgLama - salin satu organisasi dokumen ke tabel datar di transaksi pemanggil (tabel dikunci lebih dulu).
// `ditulis` false = tidak ada lagi yang perlu ditulis; ErrTidakAda = ID bukan organisasi dokumen Pega.
func (g *Gudang) SalinOrgLama(ctx context.Context, tx *db.Tx, id string, r ReferensiPindah) (bool, error) {
	if !polaIDOrg.MatchString(id) {
		return false, ErrTidakAda
	}
	nama := map[string]string{}
	for _, t := range []string{TabelClient, TabelPIC, TabelAlamat, TabelDokumen} {
		n, err := g.nama(t)
		if err != nil {
			return false, err
		}
		nama[t] = n
	}
	d, err := satu(ctx, tx, sqlDokumenSatu(nama[TabelDokumen]), func(p pemindai) (DokumenOrg, error) {
		var i, isi sql.NullString
		if err := p.Scan(&i, &isi); err != nil {
			return DokumenOrg{}, err
		}
		return UraiDokumen(i.String, []byte(isi.String))
	}, id)
	if err != nil {
		return false, err
	}
	k := KeadaanFlat{SemuaID: map[string]bool{}, Org: map[string]models.Organisasi{}, PIC: map[string][]models.PIC{},
		Alamat: map[string][]models.BarisAlamat{}}
	baris, err := daftar(ctx, tx, sqlClientSatu(nama[TabelClient]), func(p pemindai) ([]string, error) { return pindaiTeks(p, 5) },
		id, db.KosongJadiNil(d.ParentID), db.KosongJadiNil(awalanIDOk+d.ParentID))
	if err != nil {
		return false, err
	}
	for _, v := range baris {
		k.SemuaID[v[0]] = true
		if v[1] == models.FlagOrg {
			k.Org[v[0]] = models.Organisasi{ID: v[0], ParentID: v[2], Note: v[3], Title: v[4]}
		}
	}
	pic, err := daftar(ctx, tx, sqlPICSatu(nama[TabelPIC]), func(p pemindai) ([]string, error) { return pindaiTeks(p, 3) }, id)
	if err != nil {
		return false, err
	}
	for _, v := range pic {
		k.PIC[id] = append(k.PIC[id], models.PIC{UserIdentifier: v[1], Nama: v[2]})
	}
	if k.Alamat[id], err = daftar(ctx, tx, sqlDaftarAlamat(nama[TabelAlamat]), pindaiAlamat, id); err != nil {
		return false, err
	}
	var lap LaporanPindah
	rc := SusunRencana([]DokumenOrg{d}, k, r, &lap)
	if len(rc.OrgBaru) == 0 && len(rc.OrgIsi) == 0 && len(rc.PICBaru) == 0 && len(rc.TelfaxUbah) == 0 && len(rc.TelfaxSisip) == 0 {
		return false, nil
	}
	if err := g.terapkan(ctx, tx, rc); err != nil {
		return false, err
	}
	return true, nil
}
