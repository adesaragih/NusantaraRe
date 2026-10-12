package repository

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/nbfacin/backend/models"
)

// TestSQLSelesaiPolicy - tiket 69a-2: FACINPRODUCTION 79 kolom urut InsertTreatyProduction_Sql (sensus dua cara sesi c3: urai
// daftar kolom korpus = 79, token = 79, VALUES = 79; A\B = B\A = ∅, urutan sama), bind berurutan = jumlah argumen, teks tetap
// (K69a2-4 "NULL", PRORATE 1, NULL halaman tak dikenal), GetMOID_Sql; JSON_POLIS / CEDING / M_RISK_LOSS_PROFILE.
func TestSQLSelesaiPolicy(t *testing.T) {
	unik := map[string]bool{}
	for _, k := range kolomFacInProduction {
		unik[k] = true
	}
	if len(kolomFacInProduction) != 79 || len(unik) != 79 || kolomFacInProduction[0] != "IDPEGA" ||
		kolomFacInProduction[78] != "DEDUCTION2_SELISIH" {
		t.Fatalf("kolom %d / unik %d", len(kolomFacInProduction), len(unik))
	}
	q := sqlSisipFacInProduction("UJI.P", "UJI.M")
	n := bindBerurutan(t, "FACINPRODUCTION", q)
	arg := argFacInProduction(models.BarisFacInProduction{IDPega: "ID", LayerNo: "NULL", Rate: apd.New(15, -2)}, "NP", "01/11/2026 12:00:00")
	// 79 kolom - 10 NULL - 2 'NULL' - PRORATE 1 = 66 kolom ber-bind; MARKETINGOFFICERCODE dua bind (NVL MOID case, tiket 73 P5) = 67.
	// Tiket 74 P5: enam kolom Marine (TRADING_NAME / CONVEYANCE_NAME / IDSHIP / FROMRUTE / TORUTE / SAILDATE) ber-bind.
	if n != len(arg) || n != 67 {
		t.Errorf("bind %d argumen %d", n, len(arg))
	}
	for _, x := range []string{"INSERT INTO UJI.P (" + strings.Join(kolomFacInProduction, ", ") + ") VALUES (:1, :2, NULL, :3,",
		"TO_DATE(:6, 'DD/MM/YYYY HH24:MI:SS'), TO_DATE(:7, 'DD/MM/YYYY HH24:MI:SS')", "'NULL', NULL, NULL, :", ", :41, 'NULL', :42,",
		"NVL(:33, (SELECT MIN(m.ID) FROM UJI.M m WHERE m.CLIENTID = :34", " AND m.MOSTATUS = '1'))", ", 1, TO_NUMBER(:"} {
		if !strings.Contains(q, x) {
			t.Errorf("tanpa %q:\n%s", x, q)
		}
	}
	if strings.Count(q, "'NULL'") != 2 || strings.Count(q, "TO_NUMBER(") != 19 || strings.Count(q, "TO_DATE(") != 4 {
		t.Errorf("teks tetap 'NULL' %d / TO_NUMBER %d / TO_DATE %d", strings.Count(q, "'NULL'"), strings.Count(q, "TO_NUMBER("),
			strings.Count(q, "TO_DATE("))
	}
	if arg[0] != "ID" || arg[1] != "NP" {
		t.Errorf("argumen awal %v", arg[:2])
	}
	// FIRE: MOID case NULL (NVL -> GetMOID), NOSERIAL teks "NULL" (K69a2-4) - hasil sama dengan sebelum tiket 73 P5.
	if arg[32] != nil || arg[40] != "NULL" || arg[35] != nil {
		t.Errorf("MOID / NOSERIAL / TRADING_NAME FIRE: %v / %v / %v", arg[32], arg[40], arg[35])
	}
	arg = argFacInProduction(models.BarisFacInProduction{MOIDKasus: "MO-UJI", NoSerial: "SN-UJI", TradingName: "TR", IDShip: "S1",
		SailDate: "15/03/2026 12:00:00"}, "", "")
	if arg[32] != "MO-UJI" || arg[40] != "SN-UJI" || arg[35] != "TR" {
		t.Errorf("MOID / NOSERIAL / TRADING_NAME: %v / %v / %v", arg[32], arg[40], arg[35])
	}
	if !strings.Contains(q, "TO_DATE(:") || !strings.Contains(sqlSisipFacInProduction("U.P", "U.M"), "SAILDATE, POLICYTYPE") {
		t.Error("SAILDATE")
	}
	for nama, c := range map[string]struct {
		q    string
		mau  string
		bind int
	}{
		"cari polis":   {sqlCariJSONPolis("UJI.J"), "SELECT NOPOLIS, TO_CHAR(TGL_PROD, 'DD/MM/YYYY HH24:MI:SS') FROM UJI.J WHERE IDPEGA IN (:1, :2) ORDER BY TGL_INPUT DESC FETCH FIRST 1 ROWS ONLY", 2},
		"sisip polis":  {sqlSisipJSONPolis("UJI.J"), "INSERT INTO UJI.J (IDPEGA, DATA_JSON, TGL_INPUT, NOPOLIS, PRODKE, OLDNOPOLIS, TGL_PROD, USERNAME, NOENDORS) VALUES (:1, :2, SYSDATE, :3, '0', NULL, TO_DATE(:4, 'DD/MM/YYYY HH24:MI:SS'), :5, NULL)", 5},
		"hapus ceding": {sqlHapusCedingProduksi("UJI.C"), "DELETE FROM UJI.C WHERE IDPEGA IN (:1, :2)", 2},
		"ada risiko":   {sqlAdaProfilRisiko("UJI.R"), "SELECT COUNT(*) FROM UJI.R WHERE NOPOLIS = :1", 1},
		"cek facin":    {sqlAdaFacInProduction("UJI.P"), "SELECT COUNT(*) FROM UJI.P WHERE IDPEGA IN (:1, :2)", 2},
	} {
		if c.q != c.mau || bindBerurutan(t, nama, c.q) != c.bind {
			t.Errorf("%s: %q", nama, c.q)
		}
	}
	ced := sqlSisipCedingProduksi("UJI.C", "UJI.J")
	// NOENDORS_ARAS tidak ada di skema DEV (ORA-00904, 09-10-2026) - tidak ditulis.
	if bindBerurutan(t, "ceding", ced) != 8 || !strings.HasPrefix(ced, "INSERT INTO UJI.C (IDPEGA, NOPOLIS, NOENDORS, CEDINGID, CEDINGNAME, PCT_SHARE, PRODKE, NOPOLIS_ARAS) VALUES (:1, :2, NULL, :3, :4, TO_NUMBER(:5,") ||
		!strings.HasSuffix(ced, "(SELECT MAX(j.PRODKE) FROM UJI.J j WHERE j.IDPEGA IN (:6, :7)), REPLACE(:8, '.', ''))") || strings.Contains(ced, "UPDATE") ||
		strings.Contains(ced, "NOENDORS_ARAS") {
		t.Errorf("ceding %q", ced)
	}
	ris := sqlSisipProfilRisiko("UJI.R")
	if bindBerurutan(t, "risiko", ris) != 7 || !strings.HasPrefix(ris, "INSERT INTO UJI.R (NOPOLIS, TANGGALSPREADING, TSINUSANTARARE, PREMINUSANTARARE, PREMIUMOR, PREMIUMQS, PREMIUMFACOUT, COB) VALUES (:1, SYSDATE, TO_NUMBER(:2,") {
		t.Errorf("risiko %q", ris)
	}
	if !strings.Contains(sqlPosisiSubmit("W"), "STATUS_WORK NOT LIKE 'Resolved-%'") {
		t.Error("penjaga Resolved")
	}
}

// TestIdentitasPolisTahap1 - keputusan work owner 12-10-2026 (migrasi 252): UPDATE T_GENERAL_POLIS POLICY_NO / PROD_DATE_TIME /
// PROD_KE (dari JSON_POLIS IDPEGA dua bentuk) per ID case; TGL_PROD WIB -> teks Pega UTC; ORA-00904 dikenali; tulisPolicy
// memanggilnya SESUDAH JSON_POLIS (PRODKE terbaca) dan SEBELUM CEDING dengan argumen berurutan.
func TestIdentitasPolisTahap1(t *testing.T) {
	q := sqlIdentitasPolis("S.G", "S.J")
	if q != "UPDATE S.G SET POLICY_NO = :1, PROD_DATE_TIME = :2, PROD_KE = (SELECT CASE WHEN REGEXP_LIKE(MAX(j.PRODKE), '^[0-9]{1,5}$')"+
		" THEN TO_NUMBER(MAX(j.PRODKE)) END FROM S.J j WHERE j.IDPEGA IN (:3, :4)), GENERATE_DATE_POLICY_NO = NVL(:5,"+
		" GENERATE_DATE_POLICY_NO) WHERE ID = :6" || bindBerurutan(t, "identitas", q) != 6 {
		t.Errorf("%q", q)
	}
	for masuk, mau := range map[string]any{"12/10/2026 14:30:05": "20261012T073005.000 GMT", "01/11/2026 05:00:00": "20261031T220000.000 GMT",
		" 01/12/2026 12:00:00 ": "20261201T050000.000 GMT", "": nil, "2026-10-12": nil} {
		if got := prodDateTimePega(masuk); got != mau {
			t.Errorf("%q: %v, mau %v", masuk, got, mau)
		}
	}
	if !kolomBelumAda(errors.New("repository: x: ORA-00904: \"POLICY_NO\": invalid identifier")) || kolomBelumAda(nil) ||
		kolomBelumAda(errors.New("ORA-00942: table or view does not exist")) {
		t.Error("kolomBelumAda")
	}
	src, err := os.ReadFile("selesaipolicy.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(src)
	i, j, k := strings.Index(s, "sqlSisipJSONPolis(t.jp), TabelJSONPolis"), strings.Index(s, "sqlIdentitasPolis(t.general, t.jp)"),
		strings.Index(s, "sqlHapusCedingProduksi(t.ceding), \"hapus \"")
	if i < 0 || j < i || k < j || !strings.Contains(s, "jalankanKolomMigrasi(ctx, tx, \"252\", sqlIdentitasPolis(t.general, t.jp), \"identitas polis \"+TabelGeneralPolis, caseID,\n\t\tnoPolis, prodDateTimePega(prod), idPega, awalan, terbit, caseID); err != nil {") ||
		!strings.Contains(s, "if baru {\n\t\tterbit = p.Sekarang.UTC().Format(bentukWaktuPegaPolis)") {
		t.Errorf("urutan / argumen tulisPolicy: JSON_POLIS %d, identitas %d, ceding %d", i, j, k)
	}
}

// TestBackfillIdentitasPolis252 - docs/alat/isi_identitas_polis_252.sql (keputusan work owner 12-10-2026): satu UPDATE, nol COMMIT
// di luar komentar, cakupan LINI FAC + NB-% + POLICY_NO NULL + ada JSON_POLIS, IDPEGA dua bentuk, tiga kolom tahap 1 saja, dan
// rumus sama dengan penulisan Selesai Policy (teks Pega UTC dari TGL_PROD WIB, PRODKE bukan angka -> NULL).
func TestBackfillIdentitasPolis252(t *testing.T) {
	b, err := os.ReadFile("../../docs/alat/isi_identitas_polis_252.sql")
	if err != nil {
		t.Fatal(err)
	}
	var kode []string
	for _, baris := range strings.Split(strings.ReplaceAll(string(b), "\r", ""), "\n") {
		if s := strings.TrimSpace(baris); s != "" && !strings.HasPrefix(s, "--") {
			kode = append(kode, s)
		}
	}
	s := strings.Join(kode, "\n")
	if strings.Contains(strings.ToUpper(s), "COMMIT") || strings.Count(s, "UPDATE ") != 1 || strings.Contains(s, "DELETE") ||
		strings.Contains(s, "INSERT") {
		t.Fatalf("pernyataan tulis / COMMIT:\n%s", s)
	}
	upd := s[strings.Index(s, "UPDATE "):]
	upd = upd[:strings.Index(upd, ";")]
	for _, x := range []string{"SET g.POLICY_NO = (SELECT MAX(j.NOPOLIS) KEEP (DENSE_RANK LAST ORDER BY j.TGL_INPUT)",
		"g.PROD_DATE_TIME = (SELECT TO_CHAR(SYS_EXTRACT_UTC(FROM_TZ(CAST(MAX(j.TGL_PROD) KEEP (DENSE_RANK LAST ORDER BY j.TGL_INPUT)",
		"'+07:00')), 'YYYYMMDD\"T\"HH24MISS\".000 GMT\"')", "REGEXP_LIKE(MAX(j.PRODKE) KEEP (DENSE_RANK LAST ORDER BY j.TGL_INPUT), '^[0-9]{1,5}$')",
		"WHERE g.POLICY_NO IS NULL", "w.LINI = 'FAC' AND w.ID LIKE 'NB-%'"} {
		if !strings.Contains(upd, x) {
			t.Errorf("UPDATE tanpa %q", x)
		}
	}
	if strings.Count(upd, "j.IDPEGA IN (g.ID, 'ASM-FW-GISFW-WORK ' || g.ID)") != 4 || strings.Count(upd, " = (SELECT") != 3 {
		t.Errorf("IDPEGA dua bentuk / tiga kolom:\n%s", upd)
	}
	for _, k := range []string{"GENERATE_DATE_POLICY_NO", "POLICY_MASTER", "BINDER_RNM", "DATE_OFFER"} {
		if strings.Contains(s, k) {
			t.Errorf("kolom di luar sumber terverifikasi ikut diisi: %s", k)
		}
	}
	// Rumus SQL = prodDateTimePega: '+07:00' = models.ZonaWIB, format = teks Pega UTC.
	if prodDateTimePega("12/10/2026 14:30:05") != "20261012T073005.000 GMT" {
		t.Error("prodDateTimePega")
	}
}

// TestTotalDanDateOffer252 - tiket 76 (keputusan work owner 12-10-2026): UPDATE total (empat kolom NUMBER(38,8) lewat mask angkaMasuk)
// dan DATE_OFFER hanya bila NULL; TulisSubmit memanggil keduanya lewat jalankanKolom252 (toleran ORA-00904) dengan argumen berurutan.
func TestTotalDanDateOffer252(t *testing.T) {
	q := sqlRingkasanTotalPolis("S.G")
	if bindBerurutan(t, "total", q) != 5 || !strings.HasPrefix(q, "UPDATE S.G SET TOTAL_TSI_NUSA_RE = "+angkaMasuk(":1")+", TOTAL_TSI_TOP_RISK = "+
		angkaMasuk(":2")+", TOTAL_TSI_NUSA_RE_SPREADING = "+angkaMasuk(":3")+", TOTAL_PREMI_NUSA_RE = "+angkaMasuk(":4")) ||
		!strings.HasSuffix(q, " WHERE ID = :5") {
		t.Errorf("%q", q)
	}
	if q := sqlDateOffer("S.G"); q != "UPDATE S.G SET DATE_OFFER = NVL(DATE_OFFER, :1) WHERE ID = :2" {
		t.Errorf("%q", q)
	}
	src, err := os.ReadFile("submit.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(src)
	for _, x := range []string{"jalankanKolomMigrasi(ctx, tx, \"252\", sqlRingkasanTotalPolis(t.general), \"total \"+TabelGeneralPolis, w.CaseID,\n" +
		"\t\t\tikatDesimal(g.TotalTSINusaRe), ikatDesimal(g.TotalTSITopRisk), ikatDesimal(g.TotalTSINusaReSpreading),\n" +
		"\t\t\tikatDesimal(g.TotalPremiNusaRe), w.CaseID)",
		"jalankanKolomMigrasi(ctx, tx, \"252\", sqlDateOffer(t.general), \"DATE_OFFER\", w.CaseID, w.DateOffer,\n\t\t\tw.CaseID)",
		"\tif kolomBelumAda(err) {\n\t\tlog.Printf("} {
		if !strings.Contains(s, x) {
			t.Errorf("submit.go tanpa %q", x)
		}
	}
}

// TestSalinanQuotation253 - tiket 76b (keputusan work owner 12-10-2026): satu UPDATE T_QUOTATIONDATA per PARENT_ID di Save General:
// marketing / cabang dari MARKETINGOFFICER per MOID (CheckDataQuotation langkah 6), tertanggung dari rujukan opportunity / akun,
// IS_MOP SetMOP_Act, BUSINESS_FAC "F", STATUS_BUSINESS "1", bisnis :1-:6 NVL (tetap bila tak diketahui), penjaga lebar (bukan
// ORA-12899); dipanggil SimpanGeneral SESUDAH upsert T_QUOTATIONDATA dan SEBELUM ceding, toleran ORA-00904 migrasi 253.
func TestSalinanQuotation253(t *testing.T) {
	q := sqlSalinanQuotation("S.Q", "S.M", "S.O", "S.A")
	if bindBerurutan(t, "salinan", q) != 7 || !strings.HasPrefix(q, "UPDATE S.Q q SET ") || !strings.HasSuffix(q, " WHERE q.PARENT_ID = :7") {
		t.Fatalf("%q", q)
	}
	for _, x := range []string{
		"MARKETING_CODE = (SELECT MAX(CASE WHEN LENGTHB(m.CLIENTID) <= 50 THEN m.CLIENTID END) FROM S.M m WHERE m.ID = q.MOID)",
		"MARKETING_NAME = (SELECT MAX(CASE WHEN LENGTHB(m.CLIENTNAME) <= 500 THEN m.CLIENTNAME END) FROM S.M m WHERE m.ID = q.MOID)",
		"TEAM_GROUP = (SELECT MAX(CASE WHEN LENGTHB(m.TEAMGROUP) <= 50 THEN m.TEAMGROUP END) FROM S.M m WHERE m.ID = q.MOID)",
		"BRANCH_CODE = (SELECT MAX(CASE WHEN LENGTHB(m.BRANCHDETAILID) <= 50 THEN m.BRANCHDETAILID END)",
		"BRANCH_NAME = (SELECT MAX(CASE WHEN LENGTHB(m.BRANCHDETAILNAME) <= 500 THEN m.BRANCHDETAILNAME END)",
		"INSURED_ID = (SELECT MAX(CASE WHEN LENGTHB(o.INSURED_ID) <= 50 THEN o.INSURED_ID END) FROM S.O o WHERE o.ID = q.PARENT_ID)",
		"INSURED_NAME = (SELECT MAX(CASE WHEN LENGTHB(a.INSUREDNAME) <= 500 THEN a.INSUREDNAME END) FROM S.O o JOIN S.A a ON a.ID = o.ACCOUNT_ID WHERE o.ID = q.PARENT_ID)",
		"IS_MOP = CASE WHEN q.POLICY_TYPE = '0' THEN 'NonMOP' WHEN q.POLICY_TYPE IN ('1', '2') THEN 'MOP' ELSE q.IS_MOP END",
		"BUSINESS_FAC = 'F', STATUS_BUSINESS = '1', BUSINESS_CODE = NVL(:1, BUSINESS_CODE), BUSINESS_NAME = NVL(:2, BUSINESS_NAME)",
		"BUSINESS_OLD_ID = NVL(:3, BUSINESS_OLD_ID), GROUP_PANEL = NVL(:4, GROUP_PANEL), BUSINESS_TYPE = NVL(:5, BUSINESS_TYPE)",
		"BUSINESS_TYPE2 = NVL(:6, BUSINESS_TYPE2)"} {
		if !strings.Contains(q, x) {
			t.Errorf("tanpa %q", x)
		}
	}
	src, err := os.ReadFile("kasus.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(src)
	i, j, k := strings.Index(s, "sqlSisipQuotation(quo), TabelQuotationData"), strings.Index(s, "r.salinanQuotation(ctx, tx, quo, id, g.Bisnis)"),
		strings.Index(s, "return r.simpanCeding(ctx, tx, quo, id, g.CedingIDs)")
	if i < 0 || j < i || k < j || !strings.Contains(s, "jalankanKolomMigrasi(ctx, tx, \"253\", sqlSalinanQuotation(quo, t[0], t[1], t[2]), "+
		"\"salinan \"+TabelQuotationData, id,\n\t\tk(b.Code), k(b.Name), k(b.OldID), k(b.GroupPanel), k(b.Type), k(b.Type2), id)") {
		t.Errorf("urutan / argumen SimpanGeneral: upsert %d, salinan %d, ceding %d", i, j, k)
	}
}

// TestBackfillSalinanQuotation253 - docs/alat/isi_salinan_quotation_253.sql (tiket 76b): satu UPDATE, nol COMMIT / DELETE / INSERT di
// luar komentar, peringatan "nilai sekarang" di kepala, cakupan LINI FAC + NB-% + STATUS_BUSINESS NULL (hasil biner tidak ditimpa),
// rumus marketing / tertanggung / IS_MOP sama dengan sqlSalinanQuotation, BUSINESS tepat satu baris ber-GROUPPANEL, BUSINESS_TYPE /
// TYPE2 tidak diisi (turunan DT di Go).
func TestBackfillSalinanQuotation253(t *testing.T) {
	b, err := os.ReadFile("../../docs/alat/isi_salinan_quotation_253.sql")
	if err != nil {
		t.Fatal(err)
	}
	teks := strings.ReplaceAll(string(b), "\r", "")
	if !strings.Contains(teks, "NILAI SEKARANG, BUKAN NILAI SAAT CASE DIBUAT") {
		t.Error("peringatan nilai master sekarang")
	}
	var kode []string
	for _, baris := range strings.Split(teks, "\n") {
		if s := strings.TrimSpace(baris); s != "" && !strings.HasPrefix(s, "--") {
			kode = append(kode, s)
		}
	}
	s := strings.Join(kode, " ")
	if strings.Contains(strings.ToUpper(s), "COMMIT") || strings.Count(s, "UPDATE ") != 1 || strings.Contains(s, "DELETE") ||
		strings.Contains(s, "INSERT") {
		t.Fatalf("pernyataan tulis / COMMIT:\n%s", s)
	}
	upd := s[strings.Index(s, "UPDATE "):]
	upd = upd[:strings.Index(upd, ";")]
	for _, x := range []string{"q.MARKETING_CODE = (SELECT MAX(CASE WHEN LENGTHB(m.CLIENTID) <= 50 THEN m.CLIENTID END)",
		"q.TEAM_GROUP = (SELECT MAX(CASE WHEN LENGTHB(m.TEAMGROUP) <= 50 THEN m.TEAMGROUP END)",
		"q.IS_MOP = CASE WHEN q.POLICY_TYPE = '0' THEN 'NonMOP' WHEN q.POLICY_TYPE IN ('1', '2') THEN 'MOP' ELSE q.IS_MOP END",
		"q.BUSINESS_FAC = 'F'", "q.STATUS_BUSINESS = '1'", "WHERE q.STATUS_BUSINESS IS NULL",
		"w.LINI = 'FAC' AND w.ID LIKE 'NB-%'", "JOIN POOLDATA.T_M_ACCOUNT a ON a.ID = o.ACCOUNT_ID"} {
		if !strings.Contains(upd, x) {
			t.Errorf("UPDATE tanpa %q", x)
		}
	}
	if strings.Count(upd, "COUNT(*) = 1") != 3 || strings.Count(upd, "b.GROUPPANEL IS NOT NULL") != 3 || strings.Count(upd, "WHERE m.ID = q.MOID") != 5 {
		t.Errorf("BUSINESS tepat satu / MOID:\n%s", upd)
	}
	if strings.Contains(upd, "BUSINESS_TYPE") {
		t.Error("BUSINESS_TYPE / TYPE2 ikut diisi SQL")
	}
}
