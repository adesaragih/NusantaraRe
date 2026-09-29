package repository

import (
	"regexp"
	"strings"
	"testing"
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

// AC 23: pembaruan menulis SELURUH medan non-kunci - bukan hanya lima kolom
// prosedur warisan. Setiap kolom sisip (selain kunci) harus ada di SET.
func TestSQLBusinessTCOPerbaruiSeluruhMedan(t *testing.T) {
	kunci := map[string]bool{"ID": true, "TREATYYEAR": true, "TREATYGROUPID": true, "REINSTYPEID": true}
	sisip := sqlSisipBusinessTCO("S.T")
	daftarKolom := sisip[strings.Index(sisip, "(")+1 : strings.Index(sisip, ")")]
	set := sqlPerbaruiBusinessTCO("S.T")
	set = set[strings.Index(set, "SET"):strings.Index(set, "WHERE")]
	cacah := 0
	for _, k := range strings.Split(daftarKolom, ",") {
		k = strings.TrimSpace(k)
		if kunci[k] {
			continue
		}
		cacah++
		if !regexp.MustCompile(`\b` + k + ` = :`).MatchString(set) {
			t.Errorf("kolom %s tidak diperbarui - field terkirim yang tidak tersimpan (AC 23)", k)
		}
	}
	if cacah != 8 {
		t.Errorf("kolom non-kunci %d, mau 8 (12 parameter prosedur dikurangi 4 kunci)", cacah)
	}
	if !strings.Contains(sqlPerbaruiBusinessTCO("S.T"), "WHERE ID = :9 AND TREATYYEAR = :10 AND TREATYGROUPID = :11 AND REINSTYPEID = :12") {
		t.Error("pembaruan tidak dibatasi kombinasi")
	}
}

// AC 63/64: hapus menyentuh SATU tabel; nol tabel dokumen kembar.
func TestSQLBusinessTCOHapusSatuTabel(t *testing.T) {
	h := sqlHapusBusinessTCO("S.T_TREATYBUSINESS")
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
		if err := PeriksaSQL(s); err != nil {
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
