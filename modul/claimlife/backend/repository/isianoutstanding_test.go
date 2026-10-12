package repository

import (
	"context"
	"os"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

// 10-10-2026 - SQL layar Outstanding. TANPA Oracle.
func TestSQLIsianOutstanding(t *testing.T) {
	for nama, q := range map[string]string{
		"tertanggung": sqlTertanggungPeserta("S.P", "S.M"),
		"accept":      sqlSetelAcceptPeserta("S.P"),
		"katastrofe":  sqlSetelKatastrofe("S.H"),
		"kpr":         sqlSetelKPR("S.H"),
	} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Errorf("%s: %v", nama, err)
		}
		if strings.Contains(strings.ToUpper(q), "COMMIT") {
			t.Errorf("%s memuat COMMIT", nama)
		}
	}
	// Catastrophe / Note terkunci sesudah bernomor (b3515, b5059) - syaratnya di DALAM pernyataan.
	if !strings.Contains(sqlSetelKatastrofe("S.H"), "WHERE ID = :4 AND CLAIM_NO IS NULL") {
		t.Error("penulisan Catastrophe tidak dikunci CLAIM_NO IS NULL")
	}
	if strings.Contains(sqlSetelKPR("S.H"), "STS_KATASTROFE") || strings.Contains(sqlSetelKPR("S.H"), "KATASTROFE_NOTE") {
		t.Error("varian bernomor menulis Catastrophe")
	}
	// Nama dan DOB hanya DIBACA dari baris sumber; tidak ada penulisan.
	q := strings.ToUpper(sqlTertanggungPeserta("S.P", "S.M"))
	if strings.Contains(q, "UPDATE") || strings.Contains(q, "INSERT") {
		t.Error("pembaca tertanggung menulis")
	}
}

// View Detail peserta polis (11-10-2026): berpaginasi server, disaring PL_NUMBER (berindeks) + IDPEGA, urutan Pega.
func TestSQLPesertaPolisDetail(t *testing.T) {
	q := sqlPesertaPolisDetail("S.M")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	satu := strings.Join(strings.Fields(q), " ")
	for _, mau := range []string{
		"WHERE PL_NUMBER = :1 AND IDPEGA = :2",
		"ORDER BY BEGIN_DATE, CERTIFICATE_NO, NAME_OF_INSURED DESC, ID",
		"OFFSET :3 ROWS FETCH FIRST :4 ROWS ONLY",
		"COUNT(*) OVER ()",
	} {
		if !strings.Contains(satu, mau) {
			t.Errorf("SQL tanpa %q:\n%s", mau, satu)
		}
	}
	for _, k := range []string{"POLICY_NO", "CERTIFICATE_NO", "NAME_OF_INSURED", "SEX", "DOB", "AGE", "PLAN",
		"BEGIN_DATE", "EFFECTIVE_DATE", "EXPIRED_DATE"} {
		if !strings.Contains(satu, k) {
			t.Errorf("kolom grid %s tidak dibaca", k)
		}
	}
	if strings.Contains(strings.ToUpper(satu), "UPDATE") || strings.Contains(strings.ToUpper(satu), "INSERT") {
		t.Error("pembaca peserta polis menulis")
	}
	if UkuranHalamanPesertaPolis != 500 {
		t.Errorf("ukuran halaman %d, mau 500 (LoadDataPeserta_Act b572)", UkuranHalamanPesertaPolis)
	}
}

// Tanpa PL atau IDPEGA: halaman kosong TANPA menyentuh basis data (tidak memindai tabel tanpa penyaring).
func TestPesertaPolisDetailTanpaKunciTidakMembaca(t *testing.T) {
	h, err := (&KlaimLife{}).PesertaPolisDetail(context.Background(), "", "NBLF-UJI-1", 0)
	if err != nil || len(h.Baris) != 0 || h.Halaman != 1 || h.Ukuran != 500 {
		t.Errorf("hasil %+v, galat %v", h, err)
	}
}

// 11-10-2026 - Admin Notes dan Claim RNM: syarat di DALAM pernyataan.
func TestSQLCatatanAdminDanClaimRNM(t *testing.T) {
	ca := strings.Join(strings.Fields(sqlSetelCatatanAdmin("S.P")), " ")
	if !strings.Contains(ca, "(STS_REJECT IS NULL OR STS_REJECT NOT IN ('1', '2'))") || !strings.Contains(ca, "STS_HAPUS IS NULL") {
		t.Errorf("Admin Notes tidak dibekukan untuk peserta diputus (b2628):\n%s", ca)
	}
	cr := strings.Join(strings.Fields(sqlSetelClaimRNM("S.A", "S.P", "S.H")), " ")
	// 11-10-2026 "semua read only kecuali name of bank": hanya Claim RNM + Name of Bank yang ditulis.
	for _, mau := range []string{"SET a.CLAIM_AMOUNT = :1, a.NAME_OF_BANK = :2 WHERE",
		"a.STS_REJECT IS NULL", "p.CLAIM_ID = :4 AND p.STS_HAPUS IS NULL", "h.ID = :5 AND h.CLAIM_NO IS NULL"} {
		if !strings.Contains(cr, mau) {
			t.Errorf("Claim RNM tanpa %q:\n%s", mau, cr)
		}
	}
	if strings.Contains(cr, "BRANCH_OF_BANK") || strings.Contains(cr, "ACCOUNT_NO") {
		t.Errorf("Branch of Bank / Account No baca saja, tidak boleh ditulis:\n%s", cr)
	}
	for _, q := range []string{ca, cr} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Error(err)
		}
	}
}

// Lima tanggal peserta (11-10-2026): satu pernyataan, hanya selama klaim belum bernomor.
func TestSQLSetelTanggalPeserta(t *testing.T) {
	q := strings.Join(strings.Fields(sqlSetelTanggalPeserta("S.P", "S.H")), " ")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	for _, mau := range []string{"p.DATE_OF_LOSS", "p.CLAIM_RECEIVED_DATE", "p.COMPLETE_DATE", "p.MAX_DATA_RECEIVED",
		"p.CONFIRMATION_DATE", "p.STS_HAPUS IS NULL", "h.ID = :8 AND h.CLAIM_NO IS NULL"} {
		if !strings.Contains(q, mau) {
			t.Errorf("SQL tanggal tanpa %q:\n%s", mau, q)
		}
	}
}

// OutStanding Date (11-10-2026) diturunkan dari jejak transisi ” -> '0'; AmbilBaris membaca data bank.
func TestAmbilBarisMembacaJejakDanBank(t *testing.T) {
	isi, err := os.ReadFile("klaimlife.go")
	if err != nil {
		t.Fatal(err)
	}
	s := strings.Join(strings.Fields(string(isi)), " ")
	for _, mau := range []string{"SELECT MIN(j.WAKTU) FROM %s j WHERE j.ADJUSTMENT_ID = a.ID AND j.KE = '0' AND j.DARI IS NULL",
		"a.PAYABLE_TO, a.NAME_OF_BANK, a.ID_BANK, a.BRANCH_OF_BANK, a.ACCOUNT_NO"} {
		if !strings.Contains(s, mau) {
			t.Errorf("AmbilBaris tanpa %q", mau)
		}
	}
}

// Bawaan Register 11-10-2026: header ditulis bersama STS_KATASTROFE dan IS_KPR.
func TestSimpanHeaderMenulisKatastrofeDanKPR(t *testing.T) {
	isi, err := os.ReadFile("pohonklaim.go")
	if err != nil {
		t.Fatal(err)
	}
	s := strings.Join(strings.Fields(string(isi)), " ")
	if !strings.Contains(s, "(ID, CLAIM_NO, POLICY_NO, BUSINESS_NAME, STS_REJECT, CURRENCY, TYPE, STS_KATASTROFE, IS_KPR) VALUES (:1,:2,:3,:4,:5,:6,:7,:8,:9)") {
		t.Error("INSERT header tidak menulis STS_KATASTROFE / IS_KPR")
	}
	if !strings.Contains(s, "db.KosongJadiNil(p.Klaim.StsKatastrofe), db.KosongJadiNil(p.Klaim.IsKPR))") {
		t.Error("nilai Catastrophe / Is KPR tidak diikat ke INSERT header")
	}
}
