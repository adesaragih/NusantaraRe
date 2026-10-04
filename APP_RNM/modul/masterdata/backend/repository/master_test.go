package repository

import (
	"strings"
	"testing"

	"nusantarare/modul/masterdata/backend/models"
)

func master(t *testing.T, kunci string) models.TabelMaster {
	t.Helper()
	m, ada := models.CariMaster(kunci)
	if !ada {
		t.Fatalf("master %s tidak ada", kunci)
	}
	return m
}

// TestSQLDaftarMaster - kolom status tabel sendiri vs status terpisah (NATION, MD-2); kata Contains atas KolomCari
// dengan bind berurutan; halaman = dua bind terakhir.
func TestSQLDaftarMaster(t *testing.T) {
	p := master(t, "province")
	eks := ekspresiStatus(p, "UJI.S")
	where, arg := saring(p, eks, " jaw_a ", "aktif")
	if where != ` WHERE (UPPER(m.ID) LIKE :1 ESCAPE '\' OR UPPER(m.NOTE) LIKE :2 ESCAPE '\') AND m.STS_AKTIF = '1'` || len(arg) != 2 || arg[0] != `%JAW\_A%` {
		t.Fatalf("saring: %q %v", where, arg)
	}
	if q := sqlDaftar("UJI.P", p, eks, where, len(arg)); q != "SELECT m.ID, m.NATIONID, m.NOTE, m.NATIONNAME, m.STS_AKTIF FROM UJI.P m"+where+
		" ORDER BY m.ID OFFSET :3 ROWS FETCH NEXT :4 ROWS ONLY" {
		t.Errorf("daftar: %q", q)
	}
	n := master(t, "nation")
	eks = ekspresiStatus(n, "UJI.S")
	if eks != "NVL((SELECT s.STS_AKTIF FROM UJI.S s WHERE s.NAMA_TABEL = 'NATION' AND s.ID_BARIS = m.ID), '1')" {
		t.Errorf("status terpisah: %q", eks)
	}
	if where, arg := saring(n, eks, "", "nonaktif"); where != " WHERE ("+eks+" IS NULL OR "+eks+" <> '1')" || len(arg) != 0 {
		t.Errorf("nonaktif: %q", where)
	}
	if where, _ := saring(n, eks, "", ""); where != "" {
		t.Errorf("tanpa saringan: %q", where)
	}
}

// TestSQLTulisMaster - INSERT seluruh kolom (+ status '1' bila tabel berkolom status), UPDATE kolom selain ID dengan
// ID bind terakhir, kolom kata cadangan dikutip, ID akumulasi persis prosedur.
func TestSQLTulisMaster(t *testing.T) {
	p := master(t, "province")
	if q := sqlSisip("UJI.P", p); q != "INSERT INTO UJI.P (ID, NATIONID, NOTE, NATIONNAME, STS_AKTIF) VALUES (:1, :2, :3, :4, '1')" {
		t.Errorf("sisip: %q", q)
	}
	if q := sqlSisip("UJI.N", master(t, "nation")); q != "INSERT INTO UJI.N (ID, OLDID, NOTE, NATIONINITIAL) VALUES (:1, :2, :3, :4)" {
		t.Errorf("sisip nation: %q", q)
	}
	if q := sqlUbah("UJI.P", p); q != "UPDATE UJI.P SET NATIONID = :1, NOTE = :2, NATIONNAME = :3 WHERE ID = :4" {
		t.Errorf("ubah: %q", q)
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
	if q := sqlUbahStatusTerpisah("UJI.S"); q != "UPDATE UJI.S SET STS_AKTIF = :1 WHERE NAMA_TABEL = :2 AND ID_BARIS = :3" {
		t.Errorf("status terpisah: %q", q)
	}
	if PolaCari("  ") != "" || PolaCari(`a%b\`) != `%A\%B\\%` {
		t.Error("pola cari")
	}
}
