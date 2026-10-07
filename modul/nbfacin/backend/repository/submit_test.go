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
		"konteks":     {sqlKonteksSubmit(tbSubmit), 2},
		"akun":        {sqlAkunSubmit("L"), 1},
		"riwayat":     {sqlRiwayatAkseptasi("P"), 1},
		"banding":     {sqlOpsiBanding("LIM", "HP"), 3},
		"rw":          {sqlCekRW("RW"), 4},
		"okupasi":     {sqlKategoriOkupasi("OC", 3), 3},
		"produksi":    {sqlSisipRiwayatProduksi("P"), 13},
		"pega":        {sqlSisipRiwayatPega("HP"), 4},
		"performance": {sqlSisipPerformance("F"), 8},
		"posisi":      {sqlPosisiSubmit("W"), 5},
		"nbstatus":    {sqlNBStatus("G"), 2},
		"sisipNB":     {sqlSisipNBStatus("G"), 2},
		"kunciKor":    {sqlKunciKorespondensi("K"), 1},
		"kunciPen":    {sqlKunciPenerima("PN", "K"), 1},
		"sebab":       {sqlSebabDecline("C"), 0},
	} {
		if n := bindBerurutan(t, nama, c.q); n != c.mau {
			t.Errorf("%s: %d bind, mau %d", nama, n, c.mau)
		}
	}
	for _, harus := range []string{"SELECT JABATAN FROM LIM WHERE WORKBASKET IN (SELECT WORKBASKET FROM (",
		"SELECT WORKBASKET FROM HP WHERE ID_PEGA = :1 AND WORKBASKET != :2 ORDER BY TGL_TRANSFER DESC", ") WHERE ROWNUM = 1) AND TEAM_GROUP = :3"} {
		if !strings.Contains(sqlOpsiBanding("LIM", "HP"), harus) {
			t.Errorf("GetAksepBanding_SQL tanpa %q", harus)
		}
	}
	if sqlSebabDecline("C") != "SELECT ID, TERMCONDITION, DESCRIPTION, TERMCONDITION2 FROM C" {
		t.Error("CAUSEOFDECLINE")
	}
	if !strings.Contains(sqlSisipRiwayatProduksi("P"), "TO_DATE(:5, 'DD/MM/YYYY HH24:MI:SS')") ||
		!strings.Contains(sqlKonteksSubmit(tbSubmit), "m.MOSTATUS = '1' AND m.BRANCHSTATUS IS NULL AND m.CLIENTID2 <> 'LEADER'") ||
		!strings.Contains(sqlPosisiSubmit("W"), "STATUS_WORK = NVL(:2, STATUS_WORK)") || !strings.Contains(sqlPosisiSubmit("W"), "AND POSITION = :5") {
		t.Error("TGL_INP jam 24 / saringan MOLEADER / penjaga posisi")
	}
	if !regexp.MustCompile(`(?m)^\s+NB_STATUS VARCHAR2\(1000\)$`).MatchString(bacaMigrasi(t, "208_nb_status.sql")) ||
		sqlNBStatus("G") != "UPDATE G SET NB_STATUS = :1 WHERE ID = :2" {
		t.Error("NB_STATUS 208")
	}
}
