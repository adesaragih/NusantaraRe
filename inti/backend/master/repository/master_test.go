package repository

import (
	"strings"
	"testing"

	"nusantarare/inti/backend/master/models"
)

func master(t *testing.T, kunci string) models.TabelMaster {
	t.Helper()
	m, ada := models.CariMaster(kunci)
	if !ada {
		t.Fatalf("master %s tidak ada", kunci)
	}
	return m
}

// TestSQLDaftarMaster - tabel flat: kolom data + jejak (m.) + status sendiri; tabel warisan (MD-2): LEFT JOIN
// T_MASTER_STATUS untuk jejak (s.) dan - NATION - status; kata Contains atas KolomCari; halaman = dua bind terakhir.
func TestSQLDaftarMaster(t *testing.T) {
	p := master(t, "province")
	where, arg := saring(p, " jaw_a ", "aktif")
	if where != ` WHERE (UPPER(m.ID) LIKE :1 ESCAPE '\' OR UPPER(m.NOTE) LIKE :2 ESCAPE '\') AND m.STS_AKTIF = '1'` || len(arg) != 2 || arg[0] != `%JAW\_A%` {
		t.Fatalf("saring: %q %v", where, arg)
	}
	dari := sumber(p, "UJI.P", "UJI.S")
	if q := sqlDaftar(p, dari, where, len(arg)); q != "SELECT m.ID, m.NATIONID, m.NOTE, m.NATIONNAME, m.CREATE_OP, TO_CHAR(m.TGL_CREATE, 'YYYY-MM-DD HH24:MI:SS'), "+
		"m.UPDATE_OP, TO_CHAR(m.TGL_UPDATE, 'YYYY-MM-DD HH24:MI:SS'), m.STS_AKTIF FROM UJI.P m"+where+" ORDER BY m.ID OFFSET :3 ROWS FETCH NEXT :4 ROWS ONLY" {
		t.Errorf("daftar: %q", q)
	}
	n := master(t, "nation")
	if d := sumber(n, "UJI.N", "UJI.S"); d != "UJI.N m LEFT JOIN UJI.S s ON s.NAMA_TABEL = 'NATION' AND s.ID_BARIS = m.ID" {
		t.Errorf("sumber nation: %q", d)
	}
	if where, arg := saring(n, "", "nonaktif"); where != " WHERE (NVL(s.STS_AKTIF, '1') IS NULL OR NVL(s.STS_AKTIF, '1') <> '1')" || len(arg) != 0 {
		t.Errorf("nonaktif: %q", where)
	}
	o := master(t, "objectitemtype")
	if q := sqlDaftar(o, sumber(o, "UJI.O", "UJI.S"), "", 0); !strings.Contains(q, "s.CREATE_OP, TO_CHAR(s.TGL_CREATE") || !strings.Contains(q, ", m.ISACTIVE FROM UJI.O m LEFT JOIN UJI.S s") {
		t.Errorf("objectitemtype: %q", q)
	}
	if where, _ := saring(n, "", ""); where != "" {
		t.Errorf("tanpa saringan: %q", where)
	}
}

// TestSQLTulisMaster - INSERT / UPDATE / status membawa jejak (pelaku = bind, tanggal = SYSDATE) di tabel flat;
// tabel warisan tanpa kolom jejak, jejaknya di T_MASTER_STATUS (ubah lalu sisip bila belum ada).
func TestSQLTulisMaster(t *testing.T) {
	p := master(t, "province")
	if q := sqlSisip("UJI.P", p); q != "INSERT INTO UJI.P (ID, NATIONID, NOTE, NATIONNAME, STS_AKTIF, CREATE_OP, TGL_CREATE) VALUES (:1, :2, :3, :4, '1', :5, SYSDATE)" {
		t.Errorf("sisip: %q", q)
	}
	if q := sqlSisip("UJI.N", master(t, "nation")); q != "INSERT INTO UJI.N (ID, OLDID, NOTE, NATIONINITIAL) VALUES (:1, :2, :3, :4)" {
		t.Errorf("sisip nation: %q", q)
	}
	if q := sqlUbah("UJI.P", p); q != "UPDATE UJI.P SET NATIONID = :1, NOTE = :2, NATIONNAME = :3, UPDATE_OP = :4, TGL_UPDATE = SYSDATE WHERE ID = :5" {
		t.Errorf("ubah: %q", q)
	}
	if q := sqlUbah("UJI.N", master(t, "nation")); q != "UPDATE UJI.N SET OLDID = :1, NOTE = :2, NATIONINITIAL = :3 WHERE ID = :4" {
		t.Errorf("ubah nation: %q", q)
	}
	if q := sqlStatus("UJI.P", p); q != "UPDATE UJI.P SET STS_AKTIF = :1, UPDATE_OP = :2, TGL_UPDATE = SYSDATE WHERE ID = :3" {
		t.Errorf("status: %q", q)
	}
	if q := sqlStatus("UJI.O", master(t, "objectitemtype")); q != "UPDATE UJI.O SET ISACTIVE = :1 WHERE ID = :2" {
		t.Errorf("status objectitemtype: %q", q)
	}
	if q := sqlUbahJejakTerpisah("UJI.S", true); q != "UPDATE UJI.S SET STS_AKTIF = :1, UPDATE_OP = :2, TGL_UPDATE = SYSDATE WHERE NAMA_TABEL = :3 AND ID_BARIS = :4" {
		t.Errorf("jejak terpisah: %q", q)
	}
	if q := sqlUbahJejakTerpisah("UJI.S", false); q != "UPDATE UJI.S SET UPDATE_OP = :1, TGL_UPDATE = SYSDATE WHERE NAMA_TABEL = :2 AND ID_BARIS = :3" {
		t.Errorf("jejak terpisah tanpa status: %q", q)
	}
	if q := sqlSisipJejakTerpisah("UJI.S", false); q != "INSERT INTO UJI.S (UPDATE_OP, TGL_UPDATE, NAMA_TABEL, ID_BARIS) VALUES (:1, SYSDATE, :2, :3)" {
		t.Errorf("sisip jejak: %q", q)
	}
	if q := sqlBuatJejakTerpisah("UJI.S"); q != "INSERT INTO UJI.S (CREATE_OP, TGL_CREATE, NAMA_TABEL, ID_BARIS) VALUES (:1, SYSDATE, :2, :3)" {
		t.Errorf("buat jejak: %q", q)
	}
	// Baris sisa (master dihapus di luar menu, ID ditambah lagi): diatur ulang, bukan tabrakan PK.
	if q := sqlUbahBuatJejakTerpisah("UJI.S"); q != "UPDATE UJI.S SET STS_AKTIF = '1', CREATE_OP = :1, TGL_CREATE = SYSDATE, "+
		"UPDATE_OP = NULL, TGL_UPDATE = NULL WHERE NAMA_TABEL = :2 AND ID_BARIS = :3" {
		t.Errorf("buat jejak (baris sisa): %q", q)
	}
	b := models.Baris{"id": "P1", "nationId": "", "note": "N", "nationName": "X"}
	if arg := append(argKolom(b, p.Kolom[1:]), b["id"]); len(arg) != 4 || arg[0] != nil || arg[1] != "N" || arg[3] != "P1" {
		t.Errorf("bind ubah %v", arg)
	}
	if q := sqlSisip("UJI.O", master(t, "objectitemtype")); !strings.Contains(q, `"GROUP"`) || !strings.HasSuffix(q, ", '1')") || !strings.Contains(q, ", ISACTIVE)") {
		t.Errorf("objectitemtype: %q", q)
	}
	if q := sqlIDAkumulasi("UJI.SEQ"); q != "SELECT :1 || '-' || :2 || '-' || LPAD(TO_CHAR(UJI.SEQ.NEXTVAL), 6, '0') FROM DUAL" {
		t.Errorf("ID akumulasi: %q", q)
	}
	if q := sqlCatatanAkumulasi("UJI.A"); q != "SELECT COUNT(*) FROM UJI.A WHERE UPPER(NOTE) = UPPER(:1) AND ZIPCODE = :2" {
		t.Errorf("catatan akumulasi: %q", q)
	}
	if q := sqlRujukan("UJI.N", p.Rujukan[0]); q != "SELECT COUNT(*), MAX(NOTE) FROM UJI.N WHERE ID = :1" {
		t.Errorf("rujukan: %q", q)
	}
	if PolaCari("  ") != "" || PolaCari(`a%b\`) != `%A\%B\\%` {
		t.Error("pola cari")
	}
}

// TestJejakDiSetiapMaster - MD-7: tabel flat (880) membawa jejak sendiri; tepat NATION dan OBJECTITEMTYPE (warisan,
// MD-2) memakai T_MASTER_STATUS; keempat kunci jejak turunan.
func TestJejakDiSetiapMaster(t *testing.T) {
	var terpisah []string
	for _, m := range models.DaftarMaster {
		if m.JejakTerpisah {
			terpisah = append(terpisah, m.Nama)
		}
		if m.KolomStatus == "" && !m.JejakTerpisah {
			t.Errorf("%s tanpa kolom status harus JejakTerpisah", m.Nama)
		}
	}
	if strings.Join(terpisah, ",") != "NATION,OBJECTITEMTYPE" {
		t.Errorf("jejak terpisah %v", terpisah)
	}
	var kunci []string
	for _, k := range models.KolomAudit {
		if !k.Turunan {
			t.Errorf("%s harus turunan", k.Nama)
		}
		kunci = append(kunci, k.JSON)
	}
	if strings.Join(kunci, ",") != "createOp,tglCreate,updateOp,tglUpdate" {
		t.Errorf("kunci jejak %v", kunci)
	}
}

// TestCityTanpaEmailMoID - "di menu city kolom email dan mo id di hapus saja" (04-10-2026): SQL City tidak menyebut
// EMAIL / MOID, jadi tambah membiarkannya NULL dan ubah tidak mengosongkan isi lamanya; PROVINCENAME (883) ikut ditulis.
func TestCityTanpaEmailMoID(t *testing.T) {
	c := master(t, "city")
	for _, q := range []string{sqlSisip("UJI.C", c), sqlUbah("UJI.C", c), sqlDaftar(c, sumber(c, "UJI.C", "UJI.S"), "", 0)} {
		if strings.Contains(q, "EMAIL") || strings.Contains(q, "MOID") || !strings.Contains(q, "PROVINCENAME") {
			t.Errorf("SQL City: %q", q)
		}
	}
}
