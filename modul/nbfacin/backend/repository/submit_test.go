package repository

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// kolomDDLLama - kolom tabel lama POOLDATA yang ditulis / dibaca Submit, disalin dari `D:\migrasi\RNM\DDL\<TABEL>.txt`
// (HISTORYAKSEPTASIPRODUCTION, HISTORYAKSEPTASIPEGA, FACINPERFORMANCE, M_LIMIT_PROPERTYY, MARKETINGOFFICER, RW); bernilai
// true = NOT NULL.
var kolomDDLLama = map[string]map[string]bool{
	TabelRiwayatProduksi: {"IDPEGA": false, "TYPE_POLIS": false, "NOURUT": false, "POSISI": false, "PIC": false, "TGL_INP": false,
		"DIV": false, "TYPE": false, "PUTARAN": false, "APPROVAL": false, "KETERANGAN": false, "AKSES_LOGIN": false, "B2B": false,
		"BUSINESS_CODE": false, "PERCENT_RNM": false},
	TabelRiwayatPega: {"ID_PEGA": false, "TGL_TRANSFER": false, "STATUS": false, "USERNAME": false, "WORKBASKET": false,
		"ID_KOMITE": false, "OPERATORID": false},
	TabelPerformanceFacIn: {"IDPEGA": true, "POSITIONDOC": true, "STATUSDOC": true, "DATETIME": true, "USERNAME": true,
		"INSUREDNAME": true, "SOB": true, "MOLEADER": true, "COB": true, "RISLIPRECEIVED": false},
}

var tbSubmit = tabelSubmit{work: "W", opp: "O", general: "G", quo: "Q", akun: "A", mo: "MO", login: "L", produksi: "P",
	pega: "HP", perf: "F", limit: "LIM", sebab: "C", rw: "RW", okupasi: "OC", kor: "K", pen: "PN"}

// bindBerurutan - go-ora mengikat menurut urutan KEMUNCULAN: :1..:n wajib muncul berurutan, sekali masing-masing.
func bindBerurutan(t *testing.T, nama, q string) int {
	t.Helper()
	m := regexp.MustCompile(`:(\d+)`).FindAllStringSubmatch(q, -1)
	for i, x := range m {
		if n, _ := strconv.Atoi(x[1]); n != i+1 {
			t.Errorf("%s: bind ke-%d adalah :%s", nama, i+1, x[1])
		}
	}
	return len(m)
}

// TestSQLSubmit - kolom tulis = DDL lama (NOT NULL FACINPERFORMANCE diisi semua), bind berurutan dan jumlahnya = argumen
// TulisSubmit / baca, GetAksepBanding_SQL / GetCauseOfDecline verbatim, NB_STATUS ada di migrasi 208.
func TestSQLSubmit(t *testing.T) {
	sisip := regexp.MustCompile(`INSERT INTO (\w+) \(([^)]*)\)`)
	for tabel, q := range map[string]string{TabelRiwayatProduksi: sqlSisipRiwayatProduksi("P"), TabelRiwayatPega: sqlSisipRiwayatPega("HP"),
		TabelPerformanceFacIn: sqlSisipPerformance("F")} {
		m := sisip.FindStringSubmatch(strings.ReplaceAll(q, "\n", ""))
		ditulis := map[string]bool{}
		for _, c := range strings.Split(m[2], ",") {
			c = strings.TrimSpace(c)
			ditulis[c] = true
			if _, ada := kolomDDLLama[tabel][c]; !ada {
				t.Errorf("%s.%s tidak ada di DDL", tabel, c)
			}
		}
		for c, wajib := range kolomDDLLama[tabel] {
			if wajib && !ditulis[c] {
				t.Errorf("%s.%s NOT NULL tidak ditulis", tabel, c)
			}
		}
	}
	for nama, c := range map[string]struct {
		q   string
		mau int
	}{
		"konteks":      {sqlKonteksSubmit(tbSubmit), 2},
		"akun":         {sqlAkunSubmit("L"), 1},
		"riwayat":      {sqlRiwayatAkseptasi("P"), 2},
		"banding":      {sqlOpsiBanding("LIM", "HP"), 4},
		"rw":           {sqlCekRW("RW"), 4},
		"okupasi":      {sqlKategoriOkupasi("OC", 3), 3},
		"produksi":     {sqlSisipRiwayatProduksi("P"), 14},
		"pega":         {sqlSisipRiwayatPega("HP"), 4},
		"performance":  {sqlSisipPerformance("F"), 9},
		"posisi":       {sqlPosisiSubmit("W"), 5},
		"nbstatus":     {sqlNBStatus("G"), 2},
		"sisipNB":      {sqlSisipNBStatus("G"), 2},
		"deductibleTL": {sqlDeductibleTL("G"), 3},
		"keUW":         {sqlKeUW("G"), 3},
		"validitasUW":  {sqlValiditasUW("G"), 5},
		"kunciKor":     {sqlKunciKorespondensi("K"), 1},
		"kunciPen":     {sqlKunciPenerima("PN", "K"), 1},
		"sebab":        {sqlSebabDecline("C"), 0},
	} {
		if n := bindBerurutan(t, nama, c.q); n != c.mau {
			t.Errorf("%s: %d bind, mau %d", nama, n, c.mau)
		}
	}
	for _, harus := range []string{"SELECT JABATAN, JABATAN_ATASAN FROM LIM WHERE WORKBASKET IN (SELECT WORKBASKET FROM (",
		"SELECT WORKBASKET FROM HP WHERE ID_PEGA IN (:1, :2) AND WORKBASKET != :3 ORDER BY TGL_TRANSFER DESC", ") WHERE ROWNUM = 1) AND TEAM_GROUP = :4"} {
		if !strings.Contains(sqlOpsiBanding("LIM", "HP"), harus) {
			t.Errorf("GetAksepBanding_SQL tanpa %q", harus)
		}
	}
	if sqlSebabDecline("C") != "SELECT ID, TERMCONDITION, DESCRIPTION, TERMCONDITION2 FROM C" {
		t.Error("CAUSEOFDECLINE")
	}
	if !strings.Contains(sqlSisipRiwayatProduksi("P"), "TO_DATE(:5, 'DD/MM/YYYY HH24:MI:SS')") ||
		!strings.Contains(sqlKonteksSubmit(tbSubmit), "m.MOSTATUS = '1' AND m.BRANCHSTATUS IS NULL AND m.CLIENTID2 <> 'LEADER'") ||
		!strings.Contains(sqlPosisiSubmit("W"), "STATUS_WORK = NVL(:2, STATUS_WORK)") || !strings.Contains(sqlPosisiSubmit("W"), "AND POSITION = :5 AND (STATUS_WORK IS NULL OR STATUS_WORK NOT LIKE 'Resolved-%')") {
		t.Error("TGL_INP jam 24 / saringan MOLEADER / penjaga posisi")
	}
	if !regexp.MustCompile(`(?m)^\s+NB_STATUS VARCHAR2\(1000\)$`).MatchString(bacaMigrasi(t, "208_nb_status.sql")) ||
		sqlNBStatus("G") != "UPDATE G SET NB_STATUS = :1 WHERE ID = :2" {
		t.Error("NB_STATUS 208")
	}
	// Tiket 58 (migrasi 210): isian Deductible TL dan IsJuniorUW_B (SubmitToUW hanya bila kosong).
	m210 := bacaMigrasi(t, "210_submit_team_leader.sql")
	for _, k := range []string{"IS_DEDUCTIBLE_ACCEPTANCE VARCHAR2(10)", "COMMENT_DEDUCTIBLE       VARCHAR2(500)", "DATE_TO_UW               VARCHAR2(30)",
		"SUBMIT_TO_UW             VARCHAR2(30)"} {
		if !strings.Contains(m210, k) {
			t.Errorf("210 tanpa %s", k)
		}
	}
	if sqlKeUW("G") != "UPDATE G SET DATE_TO_UW = :1, SUBMIT_TO_UW = NVL(SUBMIT_TO_UW, :2) WHERE ID = :3" ||
		sqlDeductibleTL("G") != "UPDATE G SET IS_DEDUCTIBLE_ACCEPTANCE = :1, COMMENT_DEDUCTIBLE = :2 WHERE ID = :3" ||
		!strings.Contains(sqlKonteksSubmit(tbSubmit), "g.IS_DEDUCTIBLE_ACCEPTANCE, g.COMMENT_DEDUCTIBLE,\n w.CREATE_OP, g.START_DATE_TIME, "+
			"g.DATE_TO_UW, g.DAYS_VALIDITY, g.DATE_VALIDITY, g.WAITING_BIND_DATE,\n") {
		t.Error("SQL tiket 58")
	}
}

// TestMigrasiPosisiJUWB - 211 (ADR-0032): hanya baris LINI Fac In, nama lama -> ReasFacInJuniorUnderwritingB; jalur mundur
// kebalikannya.
func TestMigrasiPosisiJUWB(t *testing.T) {
	naik, turun := bacaMigrasi(t, "211_posisi_juw_b.sql"), bacaMigrasi(t, "211_posisi_juw_b_down.sql")
	if !strings.Contains(naik, "UPDATE {skema}.T_WORK_POLIS SET POSITION = 'ReasFacInJuniorUnderwritingB'\n WHERE LINI = '"+LiniFacIn+"' AND POSITION = 'ReasFacInJuniorUnderwriting'\n") ||
		!strings.Contains(turun, "UPDATE {skema}.T_WORK_POLIS SET POSITION = 'ReasFacInJuniorUnderwriting'\n WHERE LINI = '"+LiniFacIn+"' AND POSITION = 'ReasFacInJuniorUnderwritingB'\n") {
		t.Error("211")
	}
}

// TestSQLValiditasUW - tiket 61 (migrasi 212): empat kolom T_GENERAL_POLIS, bind berurutan.
func TestSQLValiditasUW(t *testing.T) {
	if q := sqlValiditasUW("UJI.G"); q != "UPDATE UJI.G SET DAYS_VALIDITY = :1, DATE_VALIDITY = :2, WAITING_BIND_DATE = :3, COMMENT_TEXT = :4 WHERE ID = :5" {
		t.Errorf("%q", q)
	}
}

// TestSQLCatatanPegaNonMarketing - tiket 57 "IS ada Ask ?": catatan terbaru (TGL_TRANSFER) di luar workbasket marketing, pola
// GetAksepBanding_SQL; dua bind.
func TestSQLCatatanPegaNonMarketing(t *testing.T) {
	q := sqlCatatanPegaNonMarketing("S.P")
	if bindTerbesar(t, q) != 3 || !strings.Contains(q, "WHERE ID_PEGA IN (:1, :2) AND WORKBASKET != :3 ORDER BY TGL_TRANSFER DESC) WHERE ROWNUM = 1") ||
		!strings.HasPrefix(q, "SELECT STATUS, WORKBASKET FROM (SELECT STATUS, WORKBASKET FROM S.P") {
		t.Error(q)
	}
}
