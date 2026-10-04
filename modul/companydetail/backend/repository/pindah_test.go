package repository

// Pindah dokumen organisasi ke tabel datar TANPA Oracle: pengurai dokumen dan penyusun rencana.

import (
	"reflect"
	"strings"
	"testing"

	"nusantarare/modul/companydetail/backend/models"
)

const dokUji = `{"Name":"UJI Organisasi","ParentID":"ORG-100","ParentName":"UJI Induk","pyDescription":"UJI catatan",
 "DataNB":{"Customer_C":{"ASMNPWP":"UJI-NPWP","ASMBusinessField":"23","ASMNationality":"001","pyTitle":"04"}},
 "PICList":[{"ASMNickName":"UJI PIC Satu","ASMPosition":"Direktur","ASMGender":"1","ASMNoTemplate":"uji@contoh.id","pyPhoneNumber":812},
            {"ASMUserIdentifier":"UJIPEGA","ASMNickName":"UJI PIC Lama"},
            {"ASMNickName":""}],
 "AddressList":[{"ASMAddress":"UJI Jalan","ASMTelfax":[{"TelfaxType":"2","TelFaxCode":"021","TelfaxNumber":"UJI-FAX"},
                                                     {"TelfaxType":"5","TelfaxNumber":"UJI-TELP"}]},
                {"ASMAddress":"UJI Jalan Kosong","ASMTelfax":[]}],
 "pxFlow":{"x":{"pyStatus":"abaikan"}}}`

func TestUraiDokumenMembacaJalurYangDipakai(t *testing.T) {
	d, err := UraiDokumen("ASM-SFAGIS-WORK-ORG ORG-201", []byte(dokUji))
	if err != nil {
		t.Fatal(err)
	}
	if d.Nama != "UJI Organisasi" || d.ParentID != "ORG-100" || d.Note != "UJI catatan" || d.NPWP != "UJI-NPWP" ||
		d.BusinessField != "23" || d.Negara != "001" || d.Title != "04" {
		t.Errorf("organisasi %+v", d)
	}
	if len(d.PIC) != 3 || d.PIC[0].Phone != "812" || d.PIC[0].Email != "uji@contoh.id" || d.PIC[1].UserIdentifier != "UJIPEGA" {
		t.Errorf("PIC %+v", d.PIC)
	}
	if len(d.Telfax["UJI Jalan"]) != 2 || d.Telfax["UJI Jalan"][0] != (models.Telfax{Jenis: "2", Code: "021", No: "UJI-FAX"}) ||
		len(d.Telfax) != 1 {
		t.Errorf("telfax %+v", d.Telfax)
	}
	if _, err := UraiDokumen("x", []byte("{bukan")); err == nil {
		t.Error("dokumen rusak diterima")
	}
}

func referensiUji() ReferensiPindah {
	return ReferensiPindah{TitleKode: map[string]string{"04": "PT."}, TitleLabel: map[string]bool{"PT.": true},
		Negara: map[string]models.Negara{"001": {ID: "100901", OldID: "001", Nama: "UJI NEGARA"}}}
}

// Dokumen tanpa baris CLIENT menjadi CLIENT baru; title kode -> LABEL, negara -> NATION, ParentID ORG-n -> CLIENT.ID.
func TestRencanaOrganisasiBaru(t *testing.T) {
	d, _ := UraiDokumen("ASM-SFAGIS-WORK-ORG ORG-201", []byte(dokUji))
	k := KeadaanFlat{SemuaID: map[string]bool{"ASM-SFAGIS-WORK-ORG ORG-100": true},
		Org: map[string]models.Organisasi{"ASM-SFAGIS-WORK-ORG ORG-100": {ID: "ASM-SFAGIS-WORK-ORG ORG-100"}}}
	var lap LaporanPindah
	rc := SusunRencana([]DokumenOrg{d}, k, referensiUji(), &lap)
	if len(rc.OrgBaru) != 1 {
		t.Fatalf("org baru %d", len(rc.OrgBaru))
	}
	o := rc.OrgBaru[0]
	if o.IDView != "ORG-201" || o.Title != "PT." || o.Country != "001" || o.CountryName != "UJI NEGARA" ||
		o.ParentID != "ASM-SFAGIS-WORK-ORG ORG-100" || o.ParentName != "UJI Induk" || o.Note != "UJI catatan" {
		t.Errorf("org baru %+v", o)
	}
	pic := rc.PICBaru["ASM-SFAGIS-WORK-ORG ORG-201"]
	if len(pic) != 2 || pic[0].UserIdentifier != "PIC-1" || pic[1].UserIdentifier != "UJIPEGA" || pic[0].Gender != "1" {
		t.Errorf("PIC baru %+v", pic)
	}
	if lap.Jumlah["PIC dokumen tanpa nama (dilewati)"] != 1 || lap.Jumlah["CLIENT baru (dokumen tanpa baris CLIENT)"] != 1 {
		t.Errorf("laporan %v", lap.Jumlah)
	}
	// Alamat bernomor tanpa baris CLIENT_ADDRESS: dilewati (trigger lama selalu membuat barisnya).
	if len(rc.TelfaxUbah) != 0 || lap.Jumlah["alamat bernomor tanpa baris CLIENT_ADDRESS (dilewati)"] != 1 {
		t.Errorf("telfax tanpa baris %+v %v", rc.TelfaxUbah, lap.Jumlah)
	}
}

// CLIENT yang sudah ada: hanya kolom KOSONG yang diisi; PIC yang sudah ada dilewati; nomor pertama mengisi baris
// alamatnya, nomor berikutnya baris baru dengan kolom tersembunyi yang sama; putaran kedua tidak menulis apa pun.
func TestRencanaClientLamaDanUlang(t *testing.T) {
	id := "ASM-SFAGIS-WORK-ORG ORG-201"
	d, _ := UraiDokumen(id, []byte(dokUji))
	alamat := models.BarisAlamat{ClientID: id, Type: "2", Address: "UJI Jalan", City: "UJIKOTA", PxCreateOperator: "UJI-PEGA"}
	k := KeadaanFlat{SemuaID: map[string]bool{id: true},
		Org:    map[string]models.Organisasi{id: {ID: id, Title: "CV."}},
		PIC:    map[string][]models.PIC{id: {{UserIdentifier: "UJIPEGA", Nama: "UJI PIC Lama"}}},
		Alamat: map[string][]models.BarisAlamat{id: {alamat}}}
	var lap LaporanPindah
	rc := SusunRencana([]DokumenOrg{d}, k, referensiUji(), &lap)
	if len(rc.OrgBaru) != 0 || !reflect.DeepEqual(rc.OrgIsi, []IsiOrg{{ID: id, Note: "UJI catatan"}}) {
		t.Errorf("isi org %+v (ParentID ORG-100 tanpa baris CLIENT; TITLE sudah ada)", rc.OrgIsi)
	}
	if pic := rc.PICBaru[id]; len(pic) != 1 || pic[0].Nama != "UJI PIC Satu" || pic[0].UserIdentifier != "PIC-1" {
		t.Errorf("PIC %+v", pic)
	}
	if len(rc.TelfaxUbah) != 1 || rc.TelfaxUbah[0].Telfax.No != "UJI-FAX" || len(rc.TelfaxSisip) != 1 ||
		rc.TelfaxSisip[0].TelfaxNo != "UJI-TELP" || rc.TelfaxSisip[0].City != "UJIKOTA" || rc.TelfaxSisip[0].PxCreateOperator != "UJI-PEGA" {
		t.Errorf("telfax ubah %+v sisip %+v", rc.TelfaxUbah, rc.TelfaxSisip)
	}
	// Putaran kedua atas keadaan sesudah putaran pertama.
	k.Org[id] = models.Organisasi{ID: id, Title: "CV.", Note: "UJI catatan"}
	k.PIC[id] = append(k.PIC[id], rc.PICBaru[id]...)
	isi := alamat
	isi.TelfaxType, isi.TelfaxNo = "2", "UJI-FAX"
	k.Alamat[id] = []models.BarisAlamat{isi, rc.TelfaxSisip[0]}
	var lap2 LaporanPindah
	rc2 := SusunRencana([]DokumenOrg{d}, k, referensiUji(), &lap2)
	if len(rc2.OrgIsi) != 0 || len(rc2.PICBaru) != 0 || len(rc2.TelfaxUbah) != 0 || len(rc2.TelfaxSisip) != 0 {
		t.Errorf("putaran kedua menulis: %+v", rc2)
	}
}

// Laporan agregat: nol ID kasus dan nol nama.
func TestLaporanTanpaIDDanNama(t *testing.T) {
	d, _ := UraiDokumen("ASM-SFAGIS-WORK-ORG ORG-201", []byte(dokUji))
	lap := LaporanPindah{Mode: "-uji"}
	SusunRencana([]DokumenOrg{d}, KeadaanFlat{SemuaID: map[string]bool{}, Org: map[string]models.Organisasi{}}, referensiUji(), &lap)
	teks := lap.Teks()
	for _, terlarang := range []string{"ORG-201", "UJI Organisasi", "UJI PIC", "UJI-FAX"} {
		if strings.Contains(teks, terlarang) {
			t.Errorf("laporan memuat %q:\n%s", terlarang, teks)
		}
	}
}

// SQL pindah: NVL tidak menimpa nilai yang sudah ada; nomor pertama hanya mengisi baris yang belum bernomor.
func TestSQLPindahTidakMenimpa(t *testing.T) {
	q := satuBaris(sqlIsiOrg("S.CLIENT"))
	if q != "UPDATE S.CLIENT SET PARENT_ID = NVL(PARENT_ID, :1), NOTE = NVL(NOTE, :2), TITLE = NVL(TITLE, :3) WHERE ID = :4 AND FLAG = :5" {
		t.Errorf("isi org: %s", q)
	}
	if q := satuBaris(sqlUbahTelfax("S.A")); !strings.HasSuffix(q, "TELFAX_TYPE IS NULL AND TELFAX_NO IS NULL AND ROWNUM = 1") {
		t.Errorf("ubah telfax: %s", q)
	}
}

// Baris popup Copy Old: organisasi baru, PIC, nomor, dan kolom yang diisi; TITLE saja bukan baris popup.
func TestRingkasRencanaCopyOld(t *testing.T) {
	d, _ := UraiDokumen("ASM-SFAGIS-WORK-ORG ORG-201", []byte(dokUji))
	k := KeadaanFlat{SemuaID: map[string]bool{}, Org: map[string]models.Organisasi{}}
	var lap LaporanPindah
	o, ada := RingkasRencana(d, SusunRencana([]DokumenOrg{d}, k, referensiUji(), &lap))
	if !ada || !o.Baru || o.IDView != "ORG-201" || o.PIC != 2 || o.Nama != "UJI Organisasi" {
		t.Errorf("ringkas org baru %+v", o)
	}
	hanyaTitle := DokumenOrg{ID: "ASM-SFAGIS-WORK-ORG ORG-202", Title: "04", Telfax: map[string][]models.Telfax{}}
	k2 := KeadaanFlat{SemuaID: map[string]bool{hanyaTitle.ID: true},
		Org: map[string]models.Organisasi{hanyaTitle.ID: {ID: hanyaTitle.ID}}}
	if o, ada := RingkasRencana(hanyaTitle, SusunRencana([]DokumenOrg{hanyaTitle}, k2, referensiUji(), &lap)); ada {
		t.Errorf("TITLE saja menjadi baris popup: %+v", o)
	}
}
