package repository

// Uji saringan jenis reasuransi non-life - TANPA Oracle (tiket 02).

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/db"
)

func TestSaringanNonLifeTCOTabelKebenaran(t *testing.T) {
	kasus := []struct {
		id, flag, tipe string
		mau            bool
	}{
		{"10003", "active", "1", true},
		{"10003", "active", "2", true},
		{"10003", "active", "3", true},
		{"10003", "active", "4", false},   // Type di luar 1,2,3
		{"10003", "inactive", "1", false}, // Flag bukan active
		{"10003", "1", "1", false},        // ejaan Life (`1`) BUKAN ejaan sini
		{"10004", "active", "1", false},   // blacklist
		{"10217", "active", "3", false},   // blacklist terakhir
		{"100041", "active", "1", false},  // NotStartsWith: berawalan 10004
		{"20004", "active", "1", true},    // bukan awalan
		{"", "active", "1", true},         // ID kosong tidak berawalan apa pun
	}
	for _, k := range kasus {
		if dapat := LolosSaringanNonLifeTCO(k.id, k.flag, k.tipe); dapat != k.mau {
			t.Errorf("(%q,%q,%q) = %v, mau %v", k.id, k.flag, k.tipe, dapat, k.mau)
		}
	}
}

func TestBlacklistJenisReasuransiVerbatimDuaBelas(t *testing.T) {
	mau := []string{"10004", "10011", "10012", "10021", "10022", "10025",
		"10026", "10028", "10248", "10249", "10018", "10217"}
	if len(BlacklistJenisReasuransiNonLife) != 12 {
		t.Fatalf("blacklist %d, mau 12", len(BlacklistJenisReasuransiNonLife))
	}
	for i, id := range mau {
		if BlacklistJenisReasuransiNonLife[i] != id {
			t.Errorf("urutan ke-%d = %s, mau %s (VERBATIM b573)", i+1, BlacklistJenisReasuransiNonLife[i], id)
		}
	}
}

func TestSQLJenisReasuransiNonLifeTCO(t *testing.T) {
	q := sqlJenisReasuransiNonLifeTCO("S.REINSURANCETYPE")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	for _, mau := range []string{"SELECT ID, NOTE, TYPE FROM S.REINSURANCETYPE",
		"FLAG = :1", "TYPE IN (:2, :3, :4)", "ID NOT LIKE :5", "ID NOT LIKE :16", "ORDER BY NOTE ASC"} {
		if !strings.Contains(q, mau) {
			t.Errorf("SQL tanpa %q:\n%s", mau, q)
		}
	}
	// NOT LIKE, bukan NOT IN: operatornya NotStartsWith (b581).
	if strings.Contains(q, "NOT IN") {
		t.Error("operator harus NOT LIKE (NotStartsWith b581), bukan NOT IN")
	}
	if strings.Count(q, "NOT LIKE") != 12 {
		t.Errorf("NOT LIKE %d kali, mau 12", strings.Count(q, "NOT LIKE"))
	}
	if strings.Contains(q, "SELECT *") {
		t.Error("SELECT * dilarang")
	}
	arg := argJenisReasuransiNonLifeTCO()
	if len(arg) != 16 {
		t.Fatalf("argumen %d, mau 16", len(arg))
	}
	if arg[0] != "active" || arg[1] != "1" || arg[3] != "3" || arg[4] != "10004%" || arg[15] != "10217%" {
		t.Errorf("argumen tidak urut: %v", arg)
	}
	// Nol literal nilai di teks SQL: semuanya bind.
	for _, id := range BlacklistJenisReasuransiNonLife {
		if strings.Contains(q, id) {
			t.Errorf("ID %s ditanam sebagai literal SQL, harus bind", id)
		}
	}
}

// Sisi korpus: dua belas nilai, flag, tipe, dan operatornya ADA di RD.
func TestSaringanNonLifeTCOVerbatimTerhadapKorpus(t *testing.T) {
	berkas := filepath.Join(`D:\XML\RNM_BRD`, "Treaty Contract Out", "ReportDefinition",
		"BrowseReinsuranceType_RD_Old_Ljt_id_isnotnull.xml")
	isi, err := os.ReadFile(berkas)
	if err != nil {
		t.Skipf("korpus tidak terjangkau: %v", err)
	}
	teks := string(isi)
	nilai := regexp.MustCompile(`<pyFilterValue>([^<]*)</pyFilterValue>`).FindAllStringSubmatch(teks, -1)
	gabung := ""
	for _, m := range nilai {
		gabung += m[1] + "|"
	}
	gabung = strings.ReplaceAll(gabung, "&quot;", `"`)
	for _, id := range BlacklistJenisReasuransiNonLife {
		if !strings.Contains(gabung, `"`+id+`"`) {
			t.Errorf("ID %s tidak ada di pyFilterValue RD", id)
		}
	}
	if !strings.Contains(gabung, `"active"`) || !strings.Contains(gabung, `"1","2","3"`) {
		t.Errorf("flag/tipe tidak ada di pyFilterValue RD: %s", gabung)
	}
	if !strings.Contains(teks, "<pyFilterOperation>NotStartsWith</pyFilterOperation>") {
		t.Error("operator NotStartsWith tidak ada di RD - SQL NOT LIKE kehilangan dasarnya")
	}
	if !strings.Contains(teks, "<pyFilterLogic>A AND B AND C</pyFilterLogic>") {
		t.Error("logika A AND B AND C tidak ada di RD")
	}
}
