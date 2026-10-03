package repository

// Teks SQL pembaca tabel flat mengikuti RD Pega (PARITAS §2, §7) - tanpa Oracle.

import (
	"fmt"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

func rata(s string) string { return strings.Join(strings.Fields(s), " ") }

func TestSQLBacaMengikutiRD(t *testing.T) {
	kasus := []struct {
		nama  string
		sql   string
		wajib []string
	}{
		{"grid BrowseProduct_Life b1096", sqlDaftarFlat("S.M_PRODUCTNAME_LIFE"),
			[]string{"SELECT ID, CEDING, TREATYNUMBER, INWARDNAME, CREATEOP, UPDATEOP FROM S.M_PRODUCTNAME_LIFE",
				"ORDER BY ID ASC", "FETCH FIRST 500 ROWS ONLY"}},
		{"produk satu baris induk", sqlBacaInduk("S.M_PRODUCTNAME_LIFE", false),
			[]string{"FROM S.M_PRODUCTNAME_LIFE WHERE ID = :1", "TO_CHAR(BEGIN_DATE, 'YYYY-MM-DD')",
				fmt.Sprintf(db.FmtDesimal, "CEDINGLIMIT"), "TO_CHAR(IS_ORS)"}},
		{"anak urut URUT", sqlBacaAnak(AnakUWLimit, "S.M_PRODUCTNAME_LIFE_UWLIMIT"),
			[]string{"FROM S.M_PRODUCTNAME_LIFE_UWLIMIT WHERE PRODUCTID = :1 ORDER BY URUT"}},
		{"komentar stempel GMT", sqlBacaAnak(AnakKomentar, "S.M_PRODUCTNAME_LIFE_COMMENT"),
			[]string{`TO_CHAR(TANGGAL, 'YYYYMMDD"T"HH24MISS.FF3')`}},
		{"sumber alat pindah", sqlSemuaJSON("S.M_PRODUCTINWARD_LIFE"), []string{"SELECT ID, JSONDATA FROM S.M_PRODUCTINWARD_LIFE ORDER BY ID ASC"}},
		{"sumber alat pindah + kolom datar", sqlSemuaJSONUmum("S.M_PRODUCT_LIFE"),
			[]string{"SELECT ID, JSONDATA, RIRISKID, RIRISK FROM S.M_PRODUCT_LIFE ORDER BY ID ASC"}},
	}
	for _, k := range kasus {
		for _, w := range k.wajib {
			if !strings.Contains(rata(k.sql), w) {
				t.Errorf("%s: SQL tanpa %q:\n%s", k.nama, w, rata(k.sql))
			}
		}
		if err := db.PeriksaSQL(k.sql); err != nil {
			t.Errorf("%s: %v", k.nama, err)
		}
		if strings.Contains(k.sql, "FOR UPDATE") {
			t.Errorf("%s: pembaca biasa tidak mengunci", k.nama)
		}
	}
	if !strings.HasSuffix(rata(sqlBacaInduk("S.T", true)), "FOR UPDATE") {
		t.Error("pembaca di dalam simpan harus mengunci baris induknya")
	}
	for _, q := range []string{sqlSisipInduk("S.T"), sqlPerbaruiInduk("S.T"), sqlSisipAnak(AnakPlan, "S.T"),
		sqlHapusAnak("S.T"), sqlHapusInduk("S.T"), sqlKunciInduk("S.T"), sqlTransaksiBacaSaja} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Errorf("%v: %s", err, rata(q))
		}
	}
	// Angka ditulis kebal NLS (koefisien bulat + skala), tanggal dan stempel berbentuk tetap.
	if q := rata(sqlSisipInduk("S.T")); !strings.Contains(q, "(TO_NUMBER(:7) / POWER(10, :8))") ||
		!strings.Contains(q, "TO_DATE(:") {
		t.Errorf("penulis induk: %s", q)
	}
	if q := sqlKunciInduk("S.T"); q != "LOCK TABLE S.T IN EXCLUSIVE MODE" || sqlHapusInduk("S.T") != "DELETE FROM S.T WHERE ID = :1" {
		t.Errorf("pindah: kunci induk dan hapus per ID: %q %q", q, sqlHapusInduk("S.T"))
	}
	if q := rata(sqlSisipAnak(AnakKomentar, "S.T")); !strings.Contains(q, "TO_TIMESTAMP(:3, 'YYYYMMDDHH24MISS.FF3')") {
		t.Errorf("penulis komentar: %s", q)
	}
}
