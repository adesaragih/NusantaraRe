package repository

import (
	"strings"
	"testing"
)

// Daftar per KOMBINASI (GetMasterReinsurerList), urut ID ASC (RD b730); uang
// dibaca sebagai teks TM9 - tidak pernah float.
func TestSQLReinsurerTCO(t *testing.T) {
	d := sqlDaftarReinsurerTCO("S.T")
	for _, mau := range []string{"TREATYYEAR = :1 AND TREATYGROUPID = :2 AND REINSTYPEID = :3", "ORDER BY ID ASC",
		"TO_CHAR(RICOMM, 'TM9'", "TO_CHAR(PCTSHARE, 'TM9'", "IUDATE", "STDRATING", "OPERATORNAME"} {
		if !strings.Contains(d, mau) {
			t.Errorf("daftar tanpa %q:\n%s", mau, d)
		}
	}
	if strings.Contains(d, "IDTREATYCONTRACT") || strings.Contains(d, "KONTRAKID") {
		t.Error("reinsurer disaring per ID kontrak - harus per kombinasi")
	}
	s := sqlSisipReinsurerTCO("S.T")
	for _, k := range []string{"REINSURERID", "CLIENTID", "NAME", "RICOMM", "PCTSHARE", "IUDATE", "USERID",
		"STARTDATE", "ENDDATE", "STATUSON", "STDRATING", "OPERATORNAME", "TGLUPDATE", ":19)"} {
		if !strings.Contains(s, k) {
			t.Errorf("sisip tanpa %s", k)
		}
	}
	p := sqlPerbaruiReinsurerTCO("S.T")
	if !strings.Contains(p, "WHERE ID = :16 AND TREATYYEAR = :17 AND TREATYGROUPID = :18 AND REINSTYPEID = :19") {
		t.Errorf("perbarui tidak dibatasi kombinasi:\n%s", p)
	}
	if !strings.Contains(sqlShareLainTCO("S.T"), "(:4 IS NULL OR ID <> :5) FOR UPDATE") {
		t.Error("share lain tanpa pengecualian baris sendiri atau tanpa kunci")
	}
	for _, q := range []string{d, s, p, sqlShareLainTCO("S.T"), sqlAmbilReinsurerTCO("S.T"),
		sqlCariReinsurerMasterTCO("S.A"), sqlAmbilReinsurerMasterTCO("S.A"), sqlKunciKontrakTCO("S.K")} {
		if err := PeriksaSQL(q); err != nil {
			t.Errorf("PeriksaSQL: %v", err)
		}
	}
}

// Pemilih: ClientName Contains tanpa beda huruf, StatusActive = 1, ketikan
// pemakai tidak menjadi wildcard.
func TestCariReinsurerMasterTCO(t *testing.T) {
	q := sqlCariReinsurerMasterTCO("S.A")
	for _, mau := range []string{"STATUSACTIVE = :1", "UPPER(CLIENTNAME) LIKE :2 ESCAPE", "ORDER BY CLIENTNAME, ID"} {
		if !strings.Contains(q, mau) {
			t.Errorf("cari tanpa %q:\n%s", mau, q)
		}
	}
	kasus := map[string]string{"re": "%RE%", " Uji Re ": "%UJI RE%", "50%_x": `%50\%\_X%`, `a\b`: `%A\\B%`}
	for masuk, mau := range kasus {
		if dapat := polaLikeTCO(masuk); dapat != mau {
			t.Errorf("polaLikeTCO(%q) = %q, mau %q", masuk, dapat, mau)
		}
	}
}

func TestPindaiReinsurerTCO(t *testing.T) {
	baris := barisPalsu{nilai: []any{"1000005", "2026", "10001", "UJI GRUP", "10003", "UJI QS", "UJI-R1",
		"UJI-C1", "UJI REAS", "7.5", "33.33333333", "", "UJI-ADMIN", "", "", "", "A", "UJI-ADMIN",
		"2026-09-29 10:00:00"}}
	r, err := pindaiReinsurerTCO(baris)
	if err != nil {
		t.Fatal(err)
	}
	if r.PctShare.Text('f') != "33.33333333" || r.Ricomm.Text('f') != "7.5" || r.Name != "UJI REAS" ||
		!r.StartDate.IsZero() || r.TglUpdate.Hour() != 10 {
		t.Errorf("pindai: %+v", r)
	}
	buruk := barisPalsu{nilai: []any{"1", "", "", "", "", "", "", "", "", "tujuh", "", "", "", "", "", "", "", "", ""}}
	if _, err := pindaiReinsurerTCO(buruk); err == nil {
		t.Error("RICOMM bukan angka harus galat, bukan nol diam-diam")
	}
}
