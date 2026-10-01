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
