package repository

import (
	"html"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"nusantarare/modul/endorsementlife/backend/models"
)

const uRiwayat = "UJISKEMA.T_VIEW_SUGGEST"

func TestSQLPutusan(t *testing.T) {
	if q := sqlRiwayat(uRiwayat); penampungUnik(t, "sqlRiwayat", q) != 1 || !strings.Contains(q, "ORDER BY s.NO DESC FETCH FIRST 200 ROWS ONLY") {
		t.Errorf("riwayat: %s", q)
	}
	if q := sqlSisipRiwayat(uRiwayat); penampungUnik(t, "sqlSisipRiwayat", q) != 7 || !strings.Contains(q, "NVL(MAX(s.NO), 0) + 1") {
		t.Errorf("sisip riwayat: %s", q)
	}
	if q := sqlAdaVersiResmi(uPolis); penampungUnik(t, "sqlAdaVersiResmi", q) != 2 || !strings.Contains(q, "p.NO_POLIS = :1 AND p.PROD_KE = :2") {
		t.Errorf("anti-dobel: %s", q)
	}
	k := sqlResmikanKepala(uPolis)
	if penampungUnik(t, "sqlResmikanKepala", k) != 6 || !strings.HasSuffix(k, "WHERE p.ID = :6 AND p.STATUSS IS NULL") {
		t.Errorf("resmikan kepala: %s", k)
	}
	for _, kol := range []string{"p.NO_POLIS", "p.PROD_KE", "p.NO_ENDORS", "p.PL_NUMBER_EDM", "p.STATUSS"} {
		if !strings.Contains(k, kol+" = :") {
			t.Errorf("resmikan kepala tanpa %s", kol)
		}
	}
	if q := sqlResmikanPeserta(uPeserta); penampungUnik(t, "sqlResmikanPeserta", q) != 6 || !strings.Contains(q, "CASE WHEN d.EDM_STATUS = :2 THEN :3 ELSE :4 END") {
		t.Errorf("resmikan peserta: %s", q)
	}
	if q := sqlTolak(uPolis); penampungUnik(t, "sqlTolak", q) != 2 || !strings.HasSuffix(q, "WHERE p.ID = :2 AND p.STATUSS IS NULL") {
		t.Errorf("tolak: %s", q)
	}
	if q := sqlHapusRekapWarisan(uRekapW); penampungUnik(t, "sqlHapusRekapWarisan", q) != 2 || !strings.Contains(q, "w.PL_NUMBER = :1 AND w.IDPEGA = :2") {
		t.Errorf("hapus rekap warisan: %s", q)
	}
	s := sqlSisipRekapWarisan(uRekapW, "UJISKEMA.M_LIFE_PREMIUM_SUMMARY_SEQ", uRekap)
	if penampungUnik(t, "sqlSisipRekapWarisan", s) != 5 || !strings.Contains(s, "TO_CHAR(UJISKEMA.M_LIFE_PREMIUM_SUMMARY_SEQ.NEXTVAL)") ||
		!strings.HasSuffix(s, "WHERE r.PREMIUM_LIST_ID = :5") || strings.Contains(s, "COMMIT") {
		t.Errorf("sisip rekap warisan: %s", s)
	}
	if strings.Count(s, "NVL(r.") != len(KolomRekapWarisan)-5 || strings.Contains(s, "TO_CHAR(r.") {
		t.Errorf("uang rekap warisan harus NUMBER langsung ber-NVL 0 (kebal NLS)")
	}
}

// TestRekapWarisanDariProsedurDanKorpus - (1) daftar `INSERT` VERBATIM badan
// prosedur (dokumen DBA PremiumList); (2) parameter `P_X` mengisi kolom `X`;
// (3) `InsertPLSummary` b85 + `InsertJsonPolisLife_Act` 12.1: `CARIn` ←
// properti `CurrencyList` bernama SAMA dengan parameternya - sehingga rekap
// kasus (kolom bernama sama) dapat disalin langsung.
func TestRekapWarisanDariProsedurDanKorpus(t *testing.T) {
	b, err := os.ReadFile(filepath.Join("..", "..", "..", "premiumlistlife", "docs", "dba-procedure-PEGA_M_LIFE_PREMIUM_SUMMARY.md"))
	if err != nil {
		t.Fatalf("dokumen prosedur: %v", err)
	}
	doc := strings.ReplaceAll(string(b), "\r", "") // salinan kerja dapat ber-CRLF, repo LF
	param := regexp.MustCompile(`(?m)^\s*\d+\s+P_(\w+) IN VARCHAR2,?$`).FindAllStringSubmatch(doc, -1)
	nilai := regexp.MustCompile(`(?m)^\s*\d+\s+P_(\w+)[,)]`).FindAllStringSubmatch(doc[strings.Index(doc, "VALUES"):], -1)
	if len(param) != 37 || len(nilai) != 37 {
		t.Fatalf("prosedur: %d parameter, %d nilai; mau 37", len(param), len(nilai))
	}
	var kolom []string
	for _, n := range nilai {
		kolom = append(kolom, n[1])
	}
	if strings.Join(kolom, ",") != strings.Join(KolomRekapWarisan, ",") {
		t.Errorf("KolomRekapWarisan bukan urutan VALUES prosedur:\n kode %v\n doc  %v", KolomRekapWarisan, kolom)
	}
	ipls, err := os.ReadFile(akarKorpus + `\RDBList\InsertPLSummary.xml`)
	if err != nil {
		t.Skipf("korpus tidak terjangkau (%v)", err)
	}
	argumen := regexp.MustCompile(`\{([\w.]+)(?: out )?\}`).FindAllStringSubmatch(html.UnescapeString(string(ipls)), -1)
	ijp, err := os.ReadFile(akarKorpus + `\Activity\InsertJsonPolisLife_Act.xml`)
	if err != nil {
		t.Skipf("korpus tidak terjangkau (%v)", err)
	}
	isi := string(ijp)
	awal := strings.Index(isi, "<pyStepsDescription>Set PremiumListSummary to Pega (Temporary)</pyStepsDescription>")
	if awal < 0 {
		t.Fatal("langkah 12.1 tidak ditemukan")
	}
	potong := isi[awal:]
	if j := strings.Index(potong, "<pyStepsMethod>RDB-List</pyStepsMethod>"); j > 0 {
		potong = potong[:j]
	}
	cari := map[string]string{}
	for _, m := range regexp.MustCompile(`<PropertiesName>TempInputData\.(CARI\d+)</PropertiesName>\s*<PropertiesValue>\.(\w+)</PropertiesValue>`).FindAllStringSubmatch(potong, -1) {
		cari[m[1]] = m[2]
	}
	if len(cari) != 33 {
		t.Fatalf("12.1 memberi %d CARI, mau 33", len(cari))
	}
	for i, a := range argumen[:37] {
		p := param[i][1]
		if c, ok := strings.CutPrefix(a[1], "TempInputData."); ok {
			if cari[c] != p {
				t.Errorf("P_%s ← %s ← .%s, mau .%s", p, c, cari[c], p)
			}
		}
	}
	if argumen[0][1] != "pyWorkPage.BusinessName" || argumen[36][1] != "pyWorkPage.pzInsKey" {
		t.Errorf("COB/IDPEGA: %s / %s", argumen[0][1], argumen[36][1])
	}
	for _, k := range KolomRekapWarisan {
		switch k {
		case "COB", "PL_NUMBER", "PL_NUMBER_EDM", "IDPEGA":
			continue
		}
		found := false
		for _, r := range models.KolomRekapKasus {
			found = found || r == k
		}
		if !found {
			t.Errorf("kolom warisan %s tidak ada di rekap kasus", k)
		}
	}
}

const uProduksiW = "UJISKEMA.LIFEINPRODUCTION"

// K4 keputusan work owner 01-10-2026 (OQ-EDM-010): `LIFEINPRODUCTION` ditulis seperti
// `RDBList/SaveLifeinProduction_SQL.xml` b86/b87 - 27 kolom VERBATIM urutan korpus, satu baris per
// kasus yang diresmikan, nilai dari kepala kasus `T_PREMIUM_LIST` (sumber yang sama dengan properti
// `pyWorkPage.*` langkah 8 b1805), nol `COMMIT`.
func TestSQLProduksiWarisan(t *testing.T) {
	q := sqlSisipProduksiWarisan(uProduksiW, uPolis)
	if penampungUnik(t, "sqlSisipProduksiWarisan", q) != 4 || !strings.HasSuffix(q, "FROM "+uPolis+" p WHERE p.ID = :4") ||
		strings.Contains(q, "COMMIT") {
		t.Errorf("sisip produksi warisan: %s", q)
	}
	if len(KolomProduksiWarisan) != 27 || KolomProduksiWarisan[0] != "IDPEGA" || KolomProduksiWarisan[26] != "TGL_INPUT" {
		t.Fatalf("27 kolom VERBATIM: %v", KolomProduksiWarisan)
	}
	if !strings.Contains(q, "INSERT INTO "+uProduksiW+" ("+strings.Join(KolomProduksiWarisan, ", ")+")") {
		t.Errorf("daftar kolom bukan urutan korpus: %s", q)
	}
	for _, w := range []string{"SELECT :1, :2, :3, p.BUSINESS_CODE,", "TRUNC(p.DATE_RECEIVED)", "p.CREATE_OP_NAME",
		"CASE p.TYPE_CEDING WHEN '1' THEN 'QS' WHEN '2' THEN 'SURPLUS' WHEN '3' THEN 'QS + SURPLUS' WHEN '4' THEN 'XOL' END",
		"p.SECURITY_REINSURER_ID, p.SECURITY_REINSURER, SYSDATE FROM"} {
		if !strings.Contains(q, w) {
			t.Errorf("sisip produksi warisan tanpa %q:\n%s", w, q)
		}
	}
}

// TestProduksiWarisanDariKorpus - (1) daftar kolom `INSERT` = `SaveLifeinProduction_SQL.xml` VERBATIM; (2) VALUES ke-n
// = `{InputDataLife.CARI1}` / `{TempInputDataLife.CARIn}`; (3) `InsertJsonPolisLife_Act` 7–8: `CARIn` ← `pyWorkPage.X`,
// dan kolom kepala sumber kode ini = kolom `KolomKepalaSalin` berproperti `X` yang sama.
func TestProduksiWarisanDariKorpus(t *testing.T) {
	b, err := os.ReadFile(akarKorpus + `\RDBList\SaveLifeinProduction_SQL.xml`)
	if err != nil {
		t.Skipf("korpus tidak terjangkau (%v)", err)
	}
	sql := html.UnescapeString(string(b))
	sql = sql[strings.Index(sql, "INSERT INTO POOLDATA.LIFEINPRODUCTION"):strings.Index(sql, "</pyBrowseSQL>")]
	kolom := strings.FieldsFunc(sql[strings.Index(sql, "(")+1:strings.Index(sql, ")")], func(r rune) bool {
		return r == ',' || r == ' ' || r == '\n' || r == '\t'
	})
	if strings.Join(kolom, ",") != strings.Join(KolomProduksiWarisan, ",") {
		t.Fatalf("KolomProduksiWarisan bukan urutan korpus:\n kode %v\n XML  %v", KolomProduksiWarisan, kolom)
	}
	nilai := regexp.MustCompile(`\{(InputDataLife|TempInputDataLife)\.(CARI\d+)\}|sysdate`).FindAllStringSubmatch(sql[strings.Index(sql, "VALUES"):], -1)
	if len(nilai) != 27 {
		t.Fatalf("VALUES: %d butir, mau 27", len(nilai))
	}
	ijp, err := os.ReadFile(akarKorpus + `\Activity\InsertJsonPolisLife_Act.xml`)
	if err != nil {
		t.Skipf("korpus tidak terjangkau (%v)", err)
	}
	properti := map[string]string{}
	for _, m := range regexp.MustCompile(`<PropertiesName>((?:Temp)?InputDataLife)\.(CARI\d+)</PropertiesName>\s*<PropertiesValue>([^<]*)</PropertiesValue>`).FindAllStringSubmatch(html.UnescapeString(string(ijp)), -1) {
		properti[m[1]+"."+m[2]] = m[3]
	}
	kepala := map[string]string{}
	for _, k := range models.KolomKepalaSalin {
		kepala["pyWorkPage."+k.Properti] = k.Kolom
	}
	for i, n := range nilai {
		k := KolomProduksiWarisan[i]
		mau := nilaiProduksiWarisan[k]
		if n[0] == "sysdate" {
			if mau != "SYSDATE" {
				t.Errorf("%s: korpus sysdate, kode %q", k, mau)
			}
			continue
		}
		asal := properti[n[1]+"."+n[2]]
		if kol, ada := kepala[asal]; ada && !strings.Contains(mau, "p."+kol) {
			t.Errorf("%s: korpus %s ← %s (kolom kepala %s), kode %q", k, n[2], asal, kol, mau)
		}
		if _, ada := kepala[asal]; !ada {
			// Bukan properti kepala: pengenal kasus, nomor polis, nomor endorsement, pembuat, nama jenis ceding.
			t.Logf("%s ← %s %s: %s", k, n[2], asal, mau)
		}
	}
}

// K5 keputusan work owner 01-10-2026 (OQ-EDM-016): peserta kasus yang diresmikan ditulis ke
// `M_LIFE_PREMIUM_DETAIL` seperti `RDBList/SaveMasterLPDet.xml` b86/b87 - 80 kolom VERBATIM, `ID` dari
// sequence, sumber berkunci `IDX_PLD_PL` (`PREMIUM_LIST_ID`), tabel warisan hanya SASARAN sisip (nol baca,
// nol pemindaian), nol `COMMIT`.
func TestSQLPesertaWarisanEDM(t *testing.T) {
	q := sqlSisipPesertaWarisanEDM("UJISKEMA.M_LIFE_PREMIUM_DETAIL", "UJISKEMA.M_LIFE_PREMIUM_DETAIL_SEQ", uPeserta, uPolis)
	if penampungUnik(t, "sqlSisipPesertaWarisanEDM", q) != 4 || strings.Contains(q, "COMMIT") ||
		!strings.HasSuffix(q, "FROM "+uPeserta+" d JOIN "+uPolis+" p ON p.ID = d.PREMIUM_LIST_ID WHERE d.PREMIUM_LIST_ID = :4") {
		t.Errorf("sisip peserta warisan: %s", q)
	}
	if len(KolomPesertaWarisanEDM) != 80 || KolomPesertaWarisanEDM[0] != "ID" || KolomPesertaWarisanEDM[79] != "RISK" {
		t.Fatalf("80 kolom VERBATIM: %d", len(KolomPesertaWarisanEDM))
	}
	for _, w := range []string{"SELECT TO_CHAR(UJISKEMA.M_LIFE_PREMIUM_DETAIL_SEQ.NEXTVAL), NVL(d.SHARE_NUSANTARA_RE, 0), d.SEX,",
		"TRUNC(d.EXPIRED_DATE)", "TO_DATE(d.STNC, 'DD/MM/YYYY')", "TO_DATE(d.WPC, 'DD/MM/YYYY')", ":1, :2, p.CEDING_CO,",
		"p.PRO_RATE_TYPE", "TRUNC(d.RETRO_VALUATION_BEGIN_DATE)", ":3, d.EDM_STATUS, d.STATUS_OLD, d.STATUS, NULL, NULL FROM"} {
		if !strings.Contains(q, w) {
			t.Errorf("sisip peserta warisan tanpa %q:\n%s", w, q)
		}
	}
	// Tabel 66,8 juta baris tidak pernah dibaca: nol `m.`, nol `WHERE` atasnya.
	if strings.Contains(q, "m.") || strings.Count(q, "WHERE") != 1 {
		t.Errorf("tabel warisan dibaca/dipindai: %s", q)
	}
}

// TestPesertaWarisanEDMDariKorpus - (1) daftar kolom `INSERT` = `SaveMasterLPDet.xml` VERBATIM; (2) VALUES ke-n
// diterjemahkan dari `InsertJsonPolisLife_Act` 11.1 (b2889): `TempInputDetail.CARIn ← @toDecimal(.X)` (penetapan
// TERAKHIR menang - CARI12) = `NVL(d.X, 0)` (`@toDecimal("")` = 0, OQ-PL-10); `TempValue.X` = `d.X`;
// `To_date(TempValue.X, …)` = tanggal; CARIn yang tidak pernah ditetapkan = NULL.
func TestPesertaWarisanEDMDariKorpus(t *testing.T) {
	b, err := os.ReadFile(akarKorpus + `\RDBList\SaveMasterLPDet.xml`)
	if err != nil {
		t.Skipf("korpus tidak terjangkau (%v)", err)
	}
	sql := html.UnescapeString(string(b))
	sql = sql[strings.Index(sql, "INSERT INTO POOLDATA.M_LIFE_PREMIUM_DETAIL"):strings.Index(sql, "</pyBrowseSQL>")]
	kolom := strings.FieldsFunc(sql[strings.Index(sql, "(")+1:strings.Index(sql, ")")], func(r rune) bool {
		return r == ',' || r == ' ' || r == '\n' || r == '\t' || r == '\r'
	})
	if strings.Join(kolom, ",") != strings.Join(KolomPesertaWarisanEDM, ",") {
		t.Fatalf("KolomPesertaWarisanEDM bukan urutan korpus:\n kode %v\n XML  %v", KolomPesertaWarisanEDM, kolom)
	}
	vals := sql[strings.Index(sql, "VALUES"):]
	vals = vals[strings.Index(vals, "(")+1 : strings.LastIndex(vals, ")")]
	butir := regexp.MustCompile(`(?:To_date\()?\{([\w.]+)\}|to_Char\(M_LIFE_PREMIUM_DETAIL_SEQ\.nextval\)`).FindAllStringSubmatch(vals, -1)
	if len(butir) != 80 {
		t.Fatalf("VALUES: %d butir, mau 80", len(butir))
	}
	ijp, err := os.ReadFile(akarKorpus + `\Activity\InsertJsonPolisLife_Act.xml`)
	if err != nil {
		t.Skipf("korpus tidak terjangkau (%v)", err)
	}
	isi := html.UnescapeString(string(ijp))
	awal, akhir := strings.Index(isi, "insert nilai dari data-batch"), strings.Index(isi, "<pyStepsDescription>Status QR/QP</pyStepsDescription>")
	if awal < 0 || akhir < awal {
		t.Fatal("langkah 11.1 / 11.2 tidak ditemukan")
	}
	isi = isi[awal:akhir]
	cari := map[string]string{}
	for _, m := range regexp.MustCompile(`<PropertiesName>TempInputDetail\.(CARI\d+)</PropertiesName>\s*<PropertiesValue>([^<]*)</PropertiesValue>`).FindAllStringSubmatch(isi, -1) {
		cari[m[1]] = m[2] // terakhir menang, seperti Pega
	}
	desimal := regexp.MustCompile(`^@toDecimal\(\.(\w+)\)$`)
	for i, m := range butir {
		k := KolomPesertaWarisanEDM[i]
		got := nilaiPesertaWarisanEDM[k]
		var mau string
		switch {
		case m[1] == "":
			mau = "SEQ"
		case strings.HasPrefix(m[0], "To_date("):
			mau = "TANGGAL"
		case m[1] == "TempInputDetail.CARI47":
			mau = "d.STATUS_OLD" // 11.4/11.5: Old → 1, selain itu 0 - ditetapkan Resmikan lebih dulu
		case m[1] == "TempInputDetail.CARI48":
			mau = "d.STATUS" // 11.2/11.3: QR/QP → 0, TR/TP → 1 - ditetapkan Resmikan lebih dulu
		case strings.HasPrefix(m[1], "TempInputDetail."):
			v, ada := cari[strings.TrimPrefix(m[1], "TempInputDetail.")]
			switch {
			case !ada:
				mau = "NULL"
			case desimal.MatchString(v):
				mau = "NVL(d." + desimal.FindStringSubmatch(v)[1] + ", 0)"
			default:
				mau = v
			}
		case strings.HasPrefix(m[1], "TempValue."):
			mau = "d." + strings.TrimPrefix(m[1], "TempValue.")
		default:
			mau = m[1]
		}
		ok := false
		switch mau {
		case "SEQ":
			ok = k == "ID" && got == ""
		case "TANGGAL":
			ok = strings.HasPrefix(got, "TRUNC(d.") || strings.HasPrefix(got, "TO_DATE(d.")
		case "NULL":
			ok = got == "NULL"
		case "d.EDMStatus":
			ok = got == "d.EDM_STATUS"
		case "pyWorkPage.CedingCo":
			ok = got == "p.CEDING_CO"
		case "pyWorkPage.ProRateType":
			ok = got == "p.PRO_RATE_TYPE"
		case "pyWorkPage.PremiumListSummary.PL_NUMBER":
			ok = got == ":1"
		case "pyWorkPage.PremiumListSummary.PL_NUMBER_EDM":
			ok = got == ":2"
		case "pyWorkPage.pzInsKey":
			ok = got == ":3"
		default:
			ok = got == mau
		}
		if !ok {
			t.Errorf("%s: korpus %q (%s), kode %q", k, m[0], mau, got)
		}
	}
}
