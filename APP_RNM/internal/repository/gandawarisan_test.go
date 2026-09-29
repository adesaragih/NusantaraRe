package repository

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"
)

func urutanPenampung(q string) string {
	var got []string
	for _, m := range regexp.MustCompile(`:(\d)`).FindAllStringSubmatch(q, -1) {
		got = append(got, m[1])
	}
	return strings.Join(got, "")
}

// TestSQLGandaBerurutPosisi - driver mengikat berposisi menurut kemunculan.
func TestSQLGandaBerurutPosisi(t *testing.T) {
	for nama, u := range map[string]struct{ q, mau string }{
		"dob":    {sqlDOBSumberKosong("S"), "123"},
		"death":  {sqlStatusWarisanTerakhir("L", "S"), "12345"},
		"health": {sqlAdaWarisanSamaDOL("L", "S"), "123456"},
		"cermin": {sqlIsiTertanggungCermin("L", "S"), "1234567"},
		"ceding": {sqlIsiCedingCermin("L", "P"), "12"},
	} {
		if err := PeriksaSQL(u.q); err != nil {
			t.Errorf("%s: %v", nama, err)
		}
		if got := urutanPenampung(u.q); got != u.mau {
			t.Errorf("%s: urutan penampung %q, mau %q\n%s", nama, got, u.mau, u.q)
		}
	}
}

// TestSQLGandaTidakMengeluarkanNama - nama dan DOB hanya di WHERE.
func TestSQLGandaTidakMengeluarkanNama(t *testing.T) {
	for _, q := range []string{sqlStatusWarisanTerakhir("L", "S"), sqlAdaWarisanSamaDOL("L", "S"), sqlDOBSumberKosong("S")} {
		pilih := q[:strings.Index(q, "FROM")]
		if strings.Contains(pilih, "NAME_OF_INSURED") || regexp.MustCompile(`\bDOB\b`).MatchString(strings.ReplaceAll(pilih, "DOB IS NULL", "")) {
			t.Errorf("SELECT mengeluarkan data pribadi:\n%s", q)
		}
	}
}

// TestSQLGandaVerbatimKunci - kelima kunci SQL warisan, dan urutannya.
func TestSQLGandaVerbatimKunci(t *testing.T) {
	q := sqlStatusWarisanTerakhir("L", "S")
	for _, k := range []string{"o.CEDINGCO", "o.NAME_OF_INSURED", "o.DOB", "o.CERTIFICATE_NO", "o.PL_NUMBER",
		"ORDER BY o.ACCEPTATION_DATE DESC", "FETCH FIRST 1 ROWS ONLY"} {
		if !strings.Contains(q, k) {
			t.Errorf("kunci %q hilang:\n%s", k, q)
		}
	}
	if !strings.Contains(sqlAdaWarisanSamaDOL("L", "S"), "o.LAPSE_DATE = TO_DATE(:5") {
		t.Error("health tidak menyaring LAPSE_DATE = DOL")
	}
}

// OQ-N2 DITUTUP (GILIRAN-17): sesudah cermin mengisi nama/DOB/CEDINGCO, baris
// milik klaim SENDIRI dan baris yang belum pernah di-Save to RNM (status NULL
// - di Pega baris cermin baru lahir saat Save Outstanding dengan '0', b176)
// tidak boleh ikut: `ORDER BY ... DESC` Oracle menaruh NULL di depan, dan
// baris seperti itu akan MENUTUPI baris era Pega yang sah.
func TestSQLGandaMengecualikanKlaimSendiriDanBarisTanpaStatus(t *testing.T) {
	death := sqlStatusWarisanTerakhir("L", "S")
	for _, mau := range []string{"AND o.STS_REJECT IS NOT NULL", "AND (o.CASEID IS NULL OR o.CASEID <> :5)"} {
		if !strings.Contains(death, mau) {
			t.Errorf("death tanpa %q:\n%s", mau, death)
		}
	}
	if health := sqlAdaWarisanSamaDOL("L", "S"); !strings.Contains(health, "AND (o.CASEID IS NULL OR o.CASEID <> :6)") {
		t.Errorf("health tidak mengecualikan klaim sendiri:\n%s", health)
	}
}

// OQ-N2 DITUTUP (GILIRAN-17): cermin OS_AKSEPTASI_KLAIM_LIFE mengisi
// NAME_OF_INSURED, DOB, CEDINGCO seperti Pega (SaveOutStandingLife_Act 22.1.1
// b8906/b8946/b9226 -> InsertJsonKlaimLife_sql b93/b95/b114) - DI DALAM SQL:
// nama dan tanggal lahir disalin dari baris sumber M_LIFE_PREMIUM_DETAIL oleh
// Oracle, tidak pernah dibaca atau diikat dari Go. CEDINGCO dari polis
// (`PolicyDataLife.CedingCo` b9226), sumber yang sama dengan pemeriksa ganda.
func TestSQLIsiTertanggungCerminDariSumber(t *testing.T) {
	q := sqlIsiTertanggungCermin("L", "S")
	for _, mau := range []string{"UPDATE L o", "SET (NAME_OF_INSURED, DOB) =",
		"SELECT m.NAME_OF_INSURED, TRUNC(m.DOB)", "FROM S m",
		"m.PL_NUMBER = :1 AND m.CERTIFICATE_NO = :2 AND m.ID = :3",
		"WHERE o.ID = :4", "m2.PL_NUMBER = :5 AND m2.CERTIFICATE_NO = :6 AND m2.ID = :7"} {
		if !strings.Contains(q, mau) {
			t.Errorf("tanpa %q:\n%s", mau, q)
		}
	}
	c := sqlIsiCedingCermin("L", "P")
	for _, mau := range []string{"UPDATE L o", "SET CEDINGCO = (SELECT pl.CEDING_CO FROM P pl WHERE pl.NO_POLIS = :1",
		"ORDER BY NVL(pl.PROD_KE, 0) DESC FETCH FIRST 1 ROWS ONLY", "WHERE o.ID = :2"} {
		if !strings.Contains(c, mau) {
			t.Errorf("ceding tanpa %q:\n%s", mau, c)
		}
	}
	if err := PeriksaSQL(c); err != nil {
		t.Error(err)
	}
	if strings.Contains(strings.ToUpper(q), "RETURNING") {
		t.Error("nama/DOB dikembalikan ke Go")
	}
	if err := PeriksaSQL(q); err != nil {
		t.Error(err)
	}
}

// TestStatusBarisKodeLamaKosongMenjadiISNULL - `= NULL` tidak pernah benar.
func TestStatusBarisKodeLamaKosongMenjadiISNULL(t *testing.T) {
	q := sqlStatusBarisAdjustment("A", true)
	if !strings.Contains(q, "STS_REJECT IS NULL") || strings.Contains(q, ":5") {
		t.Errorf("kode lama kosong:\n%s", q)
	}
	if n := len(argStatusBarisAdjustment("0", "", time.Time{}, "ADJ", "")); n != 4 {
		t.Errorf("argumen kode lama kosong = %d, mau 4 (penampung :1-:4)", n)
	}
	q = sqlStatusBarisAdjustment("A", false)
	if !strings.Contains(q, "STS_REJECT = :5") {
		t.Errorf("kode lama terisi:\n%s", q)
	}
	if n := len(argStatusBarisAdjustment("1", "N", time.Time{}, "ADJ", "0")); n != 5 {
		t.Errorf("argumen kode lama terisi = %d, mau 5", n)
	}
	for _, lamaKosong := range []bool{true, false} {
		if err := PeriksaSQL(sqlStatusBarisAdjustment("A", lamaKosong)); err != nil {
			t.Error(err)
		}
	}
}

// TestTandaiBarisOutstandingHanyaBarisAdjustment - langkah 22.1.3.2 menulis
// `.STS_REJECT = 0` baris adjustment SAJA: bukan peserta, bukan nomor atau
// tanggal akseptasi (sensus di PerbaruiStatusBaris).
func TestTandaiBarisOutstandingHanyaBarisAdjustment(t *testing.T) {
	q := fmt.Sprintf(sqlTandaiBarisOutstanding, "A")
	if err := PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(q, "STS_REJECT IS NULL") {
		t.Errorf("baris berstatus dapat ditimpa:\n%s", q)
	}
	for _, terlarang := range []string{"ACCEPTED_NO", "ACCEPTATION_DATE", "PREMIUMLIST_DETAIL", ":3"} {
		if strings.Contains(q, terlarang) {
			t.Errorf("SQL memuat %q:\n%s", terlarang, q)
		}
	}
}
