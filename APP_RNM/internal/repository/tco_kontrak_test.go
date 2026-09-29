package repository

import (
	"strings"
	"testing"
)

// Kolom VERBATIM `PEGA_TREATYCONTRACT` (8 parameter), urut ID DESC, dibatasi tahun.
func TestSQLKontrakTCO(t *testing.T) {
	d := sqlDaftarKontrakTCO("S.T")
	for _, mau := range []string{"WHERE IDTREATYYEAR = :1", "ORDER BY ID DESC", "REINSTYPEID", "REINSTYPENAME",
		"TO_CHAR(TREATYSTARTDATE", "TO_CHAR(TREATYENDDATE", "USERID", "TO_CHAR(TGLUPDATE"} {
		if !strings.Contains(d, mau) {
			t.Errorf("daftar tanpa %q:\n%s", mau, d)
		}
	}
	if s := sqlSisipKontrakTCO("S.T"); !strings.Contains(s,
		"(ID, IDTREATYYEAR, REINSTYPEID, REINSTYPENAME, TREATYSTARTDATE, TREATYENDDATE, USERID, TGLUPDATE)") ||
		!strings.Contains(s, ":8)") {
		t.Errorf("sisip: %s", s)
	}
	// AC 8: pembaruan menimpa SELURUH medan dan dibatasi tahunnya.
	p := sqlPerbaruiKontrakTCO("S.T")
	for _, mau := range []string{"REINSTYPEID = :1", "REINSTYPENAME = :2", "TREATYSTARTDATE = :3",
		"TREATYENDDATE = :4", "USERID = :5", "TGLUPDATE = :6", "WHERE ID = :7 AND IDTREATYYEAR = :8"} {
		if !strings.Contains(p, mau) {
			t.Errorf("perbarui tanpa %q:\n%s", mau, p)
		}
	}
	c := sqlCariDobelKontrakTCO("S.T")
	for _, mau := range []string{"IDTREATYYEAR = :1", "REINSTYPEID = :2", "(:3 IS NULL OR ID <> :4)",
		"FETCH FIRST 1 ROWS ONLY"} {
		if !strings.Contains(c, mau) {
			t.Errorf("dobel tanpa %q:\n%s", mau, c)
		}
	}
	for _, q := range []string{d, sqlAmbilKontrakTCO("S.T"), p, c} {
		if err := PeriksaSQL(q); err != nil {
			t.Errorf("PeriksaSQL: %v", err)
		}
	}
}

// Setiap placeholder muncul SEKALI di kueri modul ini: pengikatan per nama dan
// per kemunculan memberi hasil yang sama.
func TestTCOPlaceholderTidakBerulang(t *testing.T) {
	for nama, q := range map[string]string{
		"dobel tahun":    sqlCariDobelTahunTreaty("S.T"),
		"dobel kontrak":  sqlCariDobelKontrakTCO("S.T"),
		"lampiran":       sqlDaftarLampiranTCO("S.T", "S.O", true),
		"share lain":     sqlShareLainTCO("S.T"),
		"perbarui reas":  sqlPerbaruiReinsurerTCO("S.T"),
		"dobel bisnis":   sqlCariDobelBusinessTCO("S.T"),
		"perbarui biz":   sqlPerbaruiBusinessTCO("S.T"),
		"daftar klausul": sqlDaftarKlausulTCO("S.T"),
		"pct anak":       sqlPctAnakLainTCO("S.T"),
		"perbarui klaus": sqlPerbaruiKlausulTCO("S.T"),
		"sisip klausul":  sqlSisipKlausulTCO("S.T"),
		"sisip security": sqlSisipSecurityTCO("S.T"),
		"perbarui sec":   sqlPerbaruiSecurityTCO("S.T"),
		"dobel security": sqlCariDobelSecurityTCO("S.T"),
		"daftar sec":     sqlDaftarSecurityTCO("S.T", "S.A"),
	} {
		lihat := map[string]bool{}
		for _, b := range strings.FieldsFunc(q, func(r rune) bool {
			return !(r == ':' || (r >= '0' && r <= '9'))
		}) {
			if !strings.HasPrefix(b, ":") || len(b) < 2 {
				continue
			}
			if lihat[b] {
				t.Errorf("%s: placeholder %s berulang", nama, b)
			}
			lihat[b] = true
		}
	}
}

func TestPindaiKontrakTCO(t *testing.T) {
	baris := barisPalsu{nilai: []any{"1000003", "1000001", "10003", "UJI QS",
		"2026-01-01 00:00:00", "2027-01-01 00:00:00", "UJI-ADMIN", "2026-09-29 10:00:00"}}
	k, err := pindaiKontrakTCO(baris)
	if err != nil {
		t.Fatal(err)
	}
	if k.ID != "1000003" || k.IDTreatyYear != "1000001" || k.ReinsTypeID != "10003" ||
		k.TreatyStartDate.Format("2006-01-02") != "2026-01-01" || k.TreatyEndDate.Year() != 2027 ||
		k.TglUpdate.Hour() != 10 {
		t.Errorf("pindai: %+v", k)
	}
}
