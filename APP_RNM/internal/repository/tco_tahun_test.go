package repository

// Uji bentuk SQL tahun treaty, jejak, dan grup treaty - TANPA Oracle (tiket 03).

import (
	"strings"
	"testing"
)

func TestSQLTahunTreatyUrutIDDesc(t *testing.T) {
	q := sqlDaftarTahunTreaty("S.T_TREATYYEAR")
	for _, mau := range []string{"ORDER BY ID DESC", "OFFSET :1 ROWS FETCH NEXT :2 ROWS ONLY",
		"TO_CHAR(STARTDATE, 'YYYY-MM-DD HH24:MI:SS')", "TO_CHAR(TGLUPDATE, 'YYYY-MM-DD HH24:MI:SS')"} {
		if !strings.Contains(q, mau) {
			t.Errorf("SQL daftar tanpa %q:\n%s", mau, q)
		}
	}
	if strings.Contains(q, "SELECT *") {
		t.Error("SELECT * dilarang")
	}
	if !strings.Contains(sqlCacahTahunTreaty("S.T"), "COUNT(*)") {
		t.Error("pencacah tersendiri hilang")
	}
	for _, q := range []string{sqlDaftarTahunTreaty("S.T"), sqlAmbilTahunTreaty("S.T"),
		sqlSisipTahunTreaty("S.T"), sqlPerbaruiTahunTreaty("S.T"), sqlCariDobelTahunTreaty("S.T"),
		sqlSisipJejakTCO("S.J"), sqlBacaJejakTCO("S.J"), sqlGrupTreatyTCO("S.G")} {
		if err := PeriksaSQL(q); err != nil {
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
	for _, mau := range []string{"TREATYGROUPID = :1", "TRUNC(STARTDATE) = TRUNC(:2)",
		"TRUNC(ENDDATE) = TRUNC(:3)", "(:4 IS NULL OR ID <> :4)", "FETCH FIRST 1 ROWS ONLY"} {
		if !strings.Contains(q, mau) {
			t.Errorf("SQL dobel tanpa %q:\n%s", mau, q)
		}
	}
}

func TestSQLJejakTCOBernamaLengkap(t *testing.T) {
	q := sqlSisipJejakTCO("S.J")
	for _, k := range []string{"ID", "WAKTU", "AKUN_ID", "TABEL", "BARIS_ID", "AKSI", "KETERANGAN"} {
		if !strings.Contains(q, k) {
			t.Errorf("jejak tanpa kolom %s", k)
		}
	}
	// Kolom jejak ada di DDL 306.
	ddl := kolomDDLTCO(t)[TabelJejakTCO]
	for _, k := range []string{"ID", "WAKTU", "AKUN_ID", "TABEL", "BARIS_ID", "AKSI", "KETERANGAN"} {
		if !ddl[k] {
			t.Errorf("kolom %s tidak ada di DDL T_TREATYCO_JEJAK", k)
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
	baris := barisPalsu{nilai: []any{"1000001", "2026", "2026", "10001", "UJI GRUP", "P",
		"2026-01-01 00:00:00", "", "UJI-OP", "2026-09-15 10:00:00"}}
	th, err := pindaiTahunTreaty(baris)
	if err != nil {
		t.Fatal(err)
	}
	if th.ID != "1000001" || th.TreatyGroupName != "UJI GRUP" || th.Proportion != "P" {
		t.Errorf("teks: %+v", th)
	}
	if th.StartDate.Year() != 2026 || !th.EndDate.IsZero() || th.TglUpdate.Hour() != 10 {
		t.Errorf("tanggal: %v %v %v", th.StartDate, th.EndDate, th.TglUpdate)
	}
	buruk := barisPalsu{nilai: []any{"1", "", "", "", "", "", "kapan-kapan", "", "", ""}}
	if _, err := pindaiTahunTreaty(buruk); err == nil {
		t.Error("tanggal yang tidak terurai harus galat, bukan kosong diam-diam")
	}
}
