package repository

import (
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/db"
)

// AC 22: daftar memuat baris NONAKTIF (tanpa saringan ISACTIVE); per kombinasi.
func TestSQLBusinessTCODaftarTanpaSaringanAktif(t *testing.T) {
	d := sqlDaftarBusinessTCO("S.T")
	if strings.Contains(strings.ToUpper(d), "ISACTIVE =") {
		t.Errorf("daftar menyaring ISACTIVE - baris nonaktif harus tetap terbaca (AC 22):\n%s", d)
	}
	for _, mau := range []string{"TREATYYEAR = :1 AND TREATYGROUPID = :2 AND REINSTYPEID = :3", "ORDER BY ID ASC", "TREATYYEARID"} {
		if !strings.Contains(d, mau) {
			t.Errorf("daftar tanpa %q", mau)
		}
	}
}

// tco4 (RALAT AC 23): UPDATE PERSIS `PEGA_TREATYBUSINESS` [data DBA] - hanya
// ISACTIVE, BIZCODE, BIZNAME, USERID, TGLUPDATE; dibatasi kombinasi.
func TestSQLBusinessTCOPerbaruiSepertiProcedure(t *testing.T) {
	set := sqlPerbaruiBusinessTCO("S.T")
	set = set[strings.Index(set, "SET"):strings.Index(set, "WHERE")]
	var kolom []string
	for _, m := range regexp.MustCompile(`(\w+) = :\d+`).FindAllStringSubmatch(set, -1) {
		kolom = append(kolom, m[1])
	}
	if strings.Join(kolom, ",") != "ISACTIVE,BIZCODE,BIZNAME,USERID,TGLUPDATE" {
		t.Errorf("SET %v, mau lima kolom UPDATE procedure", kolom)
	}
	if !strings.Contains(sqlPerbaruiBusinessTCO("S.T"), "WHERE ID = :6 AND TREATYYEAR = :7 AND TREATYGROUPID = :8 AND REINSTYPEID = :9") {
		t.Error("pembaruan tidak dibatasi kombinasi")
	}
}

// AC 63/64: hapus menyentuh SATU tabel; nol tabel dokumen kembar.
func TestSQLBusinessTCOHapusSatuTabel(t *testing.T) {
	h := sqlHapusBusinessTCO("S.TREATYBUSINESS")
	if strings.Count(strings.ToUpper(h), "DELETE") != 1 || strings.Contains(strings.ToUpper(h), "M_TREATYBUSINESS") {
		t.Errorf("hapus: %s", h)
	}
	if !strings.Contains(h, "WHERE ID = :1 AND TREATYYEAR = :2 AND TREATYGROUPID = :3 AND REINSTYPEID = :4") {
		t.Errorf("hapus tidak dibatasi kombinasi: %s", h)
	}
}

func TestSQLBusinessMasterTCO(t *testing.T) {
	q := sqlDaftarBusinessMasterTCO("S.B")
	for _, mau := range []string{"SELECT ID, NOTE", "BUSINESSGROUPID IS NOT NULL", "ORDER BY NOTE ASC"} {
		if !strings.Contains(q, mau) {
			t.Errorf("master tanpa %q:\n%s", mau, q)
		}
	}
	for _, s := range []string{q, sqlAmbilBusinessMasterTCO("S.B"), sqlDaftarBusinessTCO("S.T"), sqlAmbilBusinessTCO("S.T"),
		sqlSisipBusinessTCO("S.T"), sqlPerbaruiBusinessTCO("S.T"), sqlHapusBusinessTCO("S.T"), sqlCariDobelBusinessTCO("S.T")} {
		if err := db.PeriksaSQL(s); err != nil {
			t.Errorf("PeriksaSQL: %v", err)
		}
	}
	if !strings.Contains(sqlCariDobelBusinessTCO("S.T"), "BIZCODE = :4 AND (:5 IS NULL OR ID <> :6)") {
		t.Error("dobel bisnis tanpa pengecualian baris sendiri")
	}
}

func TestPindaiBusinessTCO(t *testing.T) {
	b, err := pindaiBusinessTCO(barisPalsu{nilai: []any{"1000002", "0", "2026", "", "10001", "UJI GRUP",
		"10003", "UJI QS", "UJI-B1", "UJI BISNIS", "UJI-ADMIN", "2026-09-29 10:00:00"}})
	if err != nil {
		t.Fatal(err)
	}
	// TREATYYEARID NULL warisan terbaca kosong, bukan galat.
	if b.IsActive != "0" || b.TreatyYearID != "" || b.BizName != "UJI BISNIS" || b.TglUpdate.Hour() != 10 {
		t.Errorf("pindai: %+v", b)
	}
}
