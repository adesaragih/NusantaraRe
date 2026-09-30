package repository

// Uji bentuk SQL tahun treaty dan grup treaty - TANPA Oracle (tiket 03, tco4).

import (
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

func TestSQLTahunTreatyUrutIDDesc(t *testing.T) {
	q := sqlDaftarTahunTreaty("S.TREATYYEAR")
	for _, mau := range []string{"ORDER BY ID DESC", "OFFSET :1 ROWS FETCH NEXT :2 ROWS ONLY",
		"STARTDATE, ENDDATE, USERID, TGLUPDATE"} {
		if !strings.Contains(q, mau) {
			t.Errorf("SQL daftar tanpa %q:\n%s", mau, q)
		}
	}
	// tco4: kolom tanggal TREATYYEAR VARCHAR2 - TO_CHAR berformat tanggal
	// akan gagal (ORA-01722) atau salah baca.
	if strings.Contains(q, "TO_CHAR(") {
		t.Errorf("TO_CHAR pada kolom VARCHAR2 warisan:\n%s", q)
	}
	if strings.Contains(q, "SELECT *") {
		t.Error("SELECT * dilarang")
	}
	if !strings.Contains(sqlCacahTahunTreaty("S.T"), "COUNT(*)") {
		t.Error("pencacah tersendiri hilang")
	}
	for _, q := range []string{sqlDaftarTahunTreaty("S.T"), sqlAmbilTahunTreaty("S.T"),
		sqlSisipTahunTreaty("S.T"), sqlPerbaruiTahunTreaty("S.T"), sqlCariDobelTahunTreaty("S.T"),
		sqlGrupTreatyTCO("S.G")} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Errorf("%v:\n%s", err, q)
		}
	}
}

// AC 8: pembaruan menimpa SELURUH kolom; AC 6: sisipan bernama sepuluh kolom.
func TestSQLTahunTreatySisipDanPerbaruiSeluruhKolom(t *testing.T) {
	sisip := sqlSisipTahunTreaty("S.T")
	perbarui := sqlPerbaruiTahunTreaty("S.T")
	for _, k := range KolomWarisanTCO(warisanTahunTCO) {
		if !strings.Contains(sisip, k) {
			t.Errorf("INSERT tanpa kolom %s", k)
		}
		if k != "ID" && !strings.Contains(perbarui, k+" = :") {
			t.Errorf("UPDATE tidak menimpa kolom %s", k)
		}
	}
	if !strings.Contains(perbarui, "WHERE ID = :10") {
		t.Errorf("UPDATE tidak berkunci ID:\n%s", perbarui)
	}
	if strings.Count(sisip, ":") != 10 {
		t.Errorf("INSERT penampung %d, mau 10", strings.Count(sisip, ":"))
	}
}

// AC 73: anti-dobel berkunci (grup, mulai, akhir), mengecualikan diri
// sendiri, dan tahan terhadap ID kosong (NULL Oracle).
func TestSQLCariDobelTahunTreaty(t *testing.T) {
	q := sqlCariDobelTahunTreaty("S.T")
	for _, mau := range []string{"SELECT ID, STARTDATE, ENDDATE", "TREATYGROUPID = :1", "(:2 IS NULL OR ID <> :3)"} {
		if !strings.Contains(q, mau) {
			t.Errorf("SQL dobel tanpa %q:\n%s", mau, q)
		}
	}
}

func TestSQLGrupTreatyTCO(t *testing.T) {
	q := sqlGrupTreatyTCO("S.TREATYGROUP")
	if !strings.Contains(q, "SELECT ID, TREATYGROUPNAME FROM S.TREATYGROUP ORDER BY ID DESC") {
		t.Errorf("SQL grup treaty:\n%s", q)
	}
}

// Instrumen pindai: sepuluh kolom, tanggal terurai, kosong tetap kosong.
type barisPalsu struct{ nilai []any }

func (b barisPalsu) Scan(tujuan ...any) error { return isiNullString(tujuan, b.nilai) }

func TestPindaiTahunTreaty(t *testing.T) {
	// Lanjutan 4 (OQ-TCO-01): tanggal YYYYMMDD seperti data warisan; TGLUPDATE stempel Pega.
	baris := barisPalsu{nilai: []any{"1000001", "2026", "2026", "10001", "UJI GRUP", "P",
		"20260101", "", "UJI-OP", "20260915T100000.000 GMT"}}
	th, err := pindaiTahunTreaty(baris)
	if err != nil {
		t.Fatal(err)
	}
	if th.ID != "1000001" || th.TreatyGroupName != "UJI GRUP" || th.Proportion != "P" {
		t.Errorf("teks: %+v", th)
	}
	if th.StartDate.Year() != 2026 || th.StartDate.Month() != 1 || th.StartDate.Day() != 1 ||
		!th.EndDate.IsZero() || th.TglUpdate.Hour() != 17 { // 10:00 GMT = 17:00 WIB
		t.Errorf("tanggal: %v %v %v", th.StartDate, th.EndDate, th.TglUpdate)
	}
	for _, bentuk := range []string{"kapan-kapan", "20251231T170000.000 GMT"} {
		buruk := barisPalsu{nilai: []any{"1", "", "", "", "", "", bentuk, "", "", ""}}
		if _, err := pindaiTahunTreaty(buruk); err == nil {
			t.Errorf("%q harus galat, bukan kosong atau ditebak", bentuk)
		}
	}
}
