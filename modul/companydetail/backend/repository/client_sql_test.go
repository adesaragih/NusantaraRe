package repository

// SQL Company Detail TANPA Oracle: kolom, urutan bind, kunci baris, dan nol COMMIT.

import (
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/companydetail/backend/models"
)

func satuBaris(q string) string { return strings.Join(strings.Fields(q), " ") }

// Baris organisasi baru menulis kolom yang ditulis RDBINSERTCLIENT (tanpa IDNUMBER dan BU_NOTE) + kolom 805.
func TestSisipOrgKolomDanBind(t *testing.T) {
	q := satuBaris(sqlSisipOrg("S.CLIENT"))
	mau := "INSERT INTO S.CLIENT (ID, NAME, FLAG, BU_ID, IDVIEW, NPWP, GROUPNAME, TITLE, COUNTRY, COUNTRYNAME, " +
		"PARENT_ID, NOTE, CREATED_BY, CREATED_AT, UPDATED_BY, UPDATED_AT) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, " +
		":10, :11, :12, :13, SYSTIMESTAMP, :14, SYSTIMESTAMP)"
	if q != mau {
		t.Errorf("sisip:\n dapat %s\n mau   %s", q, mau)
	}
	o := models.Organisasi{ID: "i", Nama: "n", BusinessField: "bu", IDView: "iv", NPWP: "np", ParentName: "pn",
		Title: "t", Country: "c", CountryName: "cn", ParentID: "pi", Note: "no", CreatedBy: "cb", UpdatedBy: "ub"}
	n := NilaiSisipOrg(o)
	mauNilai := []any{"i", "n", "Org", "bu", "iv", "np", "pn", "t", "c", "cn", "pi", "no", "cb", "ub"}
	if len(n) != len(mauNilai) {
		t.Fatalf("nilai %d, mau %d", len(n), len(mauNilai))
	}
	for i := range n {
		if n[i] != mauNilai[i] {
			t.Errorf("bind :%d = %v, mau %v", i+1, n[i], mauNilai[i])
		}
	}
	if kosong := NilaiSisipOrg(models.Organisasi{ID: "i"}); kosong[1] != nil || kosong[13] != nil {
		t.Errorf("kosong = NULL: %v", kosong)
	}
}

// UPDATE tidak pernah menyentuh ID, IDVIEW, FLAG, IDNUMBER, BU_NOTE, dan jejak pembuatan.
func TestPerbaruiOrgTanpaKolomTetap(t *testing.T) {
	q := satuBaris(sqlPerbaruiOrg("S.CLIENT"))
	for _, k := range []string{"IDVIEW", "FLAG =", "IDNUMBER", "BU_NOTE", "CREATED_"} {
		if strings.Contains(strings.SplitN(q, "WHERE", 2)[0], k) {
			t.Errorf("perbarui menyentuh %s: %s", k, q)
		}
	}
	if !strings.HasSuffix(q, "WHERE ID = :11 AND FLAG = :12") || !strings.Contains(q, "UPDATED_AT = SYSTIMESTAMP") {
		t.Errorf("perbarui: %s", q)
	}
	n := NilaiPerbaruiOrg(models.Organisasi{ID: "ASM-SFAGIS-WORK-ORG ORG-1", Nama: "x"})
	if len(n) != 12 || n[10] != "ASM-SFAGIS-WORK-ORG ORG-1" || n[11] != "Org" {
		t.Errorf("bind perbarui: %v", n)
	}
}

// Satu baris alamat: empat belas kolom, CLIENTID di bind pertama, kolom telfax di tiga terakhir.
func TestSisipAlamatKolomDanBind(t *testing.T) {
	q := satuBaris(sqlSisipAlamat("S.CLIENT_ADDRESS"))
	if !strings.Contains(q, "(CLIENTID, ASMADDRESSTYPE, ASMADDRESS, ASMCITY, CITYNAME, ASMZIPCODE, DISTRICTNAME, "+
		"PROVINCENAME, RWNAME, PXCREATEOPERATOR, PXCREATEDATETIME, TELFAX_TYPE, TELFAX_CODE, TELFAX_NO)") ||
		!strings.HasSuffix(q, ":13, :14)") {
		t.Errorf("sisip alamat: %s", q)
	}
	n := NilaiSisipAlamat(models.BarisAlamat{ClientID: "c", Address: "a", TelfaxType: "5", TelfaxNo: "1"})
	if len(n) != 14 || n[0] != "c" || n[2] != "a" || n[11] != "5" || n[12] != nil || n[13] != "1" {
		t.Errorf("bind alamat: %v", n)
	}
}

// Pola cari: huruf besar, `%` `_` `\` ketikan pemakai di-escape; kueri kosong = :2 NULL (tanpa saringan).
func TestPolaCariDanSaringan(t *testing.T) {
	if p := PolaCari(` a%b_c\ `); p != `%A\%B\_C\\%` {
		t.Errorf("pola %q", p)
	}
	if n := NilaiCari(" "); n[0] != "Org" || n[1] != nil {
		t.Errorf("tanpa saringan: %v", n)
	}
	if n := NilaiCari("uji"); n[1] != "%UJI%" || n[4] != "%UJI%" {
		t.Errorf("saringan: %v", n)
	}
	q := satuBaris(sqlCariOrg("S.CLIENT"))
	if !strings.Contains(q, "OFFSET :6 ROWS FETCH NEXT :7 ROWS ONLY") || !strings.Contains(q, "ORDER BY UPPER(NAME), ID") {
		t.Errorf("cari: %s", q)
	}
}

// Setiap SQL modul ini lolos PeriksaSQL (nol COMMIT/ROLLBACK).
func TestSeluruhSQLLolosPeriksa(t *testing.T) {
	for _, q := range []string{sqlCariOrg("S.T"), sqlHitungOrg("S.T"), sqlAmbilOrg("S.T"), sqlNomorOrgBerikut("S.SEQ"),
		sqlSisipOrg("S.T"), sqlPerbaruiOrg("S.T"), sqlDaftarPIC("S.T"), sqlHapusPIC("S.T"),
		sqlSisipPIC("S.T"), sqlDaftarAlamat("S.T"), sqlHapusAlamat("S.T"), sqlSisipAlamat("S.T"),
		sqlDaftarPilihan("S.T"), sqlDaftarNegara("S.T"), sqlCariInduk("S.T"), sqlNamaOrg("S.T")} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Error(err)
		}
	}
}

// Pilihan PIC Name: akun login AKTIF saja, urut nama; M_LOGIN_GO hanya dibaca.
func TestDaftarAkunAktif(t *testing.T) {
	q := satuBaris(sqlDaftarAkunAktif("S.M_LOGIN_GO"))
	mau := "SELECT LOGIN_ID, NAME, JOB_POSITION FROM S.M_LOGIN_GO WHERE IS_ACTIVE = :1 ORDER BY UPPER(NAME), LOGIN_ID"
	if q != mau {
		t.Errorf("akun:\n dapat %s\n mau   %s", q, mau)
	}
	for _, tb := range DaftarTabelDitulis {
		if tb == TabelLogin {
			t.Errorf("%s tidak boleh ditulis", TabelLogin)
		}
	}
}
