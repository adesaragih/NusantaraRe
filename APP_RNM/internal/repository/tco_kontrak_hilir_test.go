package repository

// Uji kontrak baca hilir Treaty Contract Out - TANPA Oracle (tiket 01, tco3).
//
// Tiga sisi dikunci:
//  1. kolom yang pembaca sebut ADA di DDL tabel T_* (300-305);
//  2. kolom yang pembaca sebut ADA di kueri hilir korpus - VERBATIM;
//  3. nol kata kerja tulis di pembaca.
//
// Sisi 2 DILEWATI bila korpus tidak terjangkau (mesin lain) - bukan gagal.

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// kolomDDLTCO membaca kolom tiap tabel T_* dari berkas migrasi 300-306.
func kolomDDLTCO(t *testing.T) map[string]map[string]bool {
	t.Helper()
	hasil := map[string]map[string]bool{}
	for nama, teks := range seluruhSQL(t, false) {
		if !strings.HasPrefix(nama, "30") {
			continue
		}
		for _, p := range strings.Split(teks, "\n/") {
			tabel, kolom := KolomCreateTable(p)
			if tabel == "" {
				continue
			}
			hasil[tabel] = map[string]bool{}
			for _, k := range kolom {
				hasil[tabel][k] = true
			}
		}
	}
	if len(hasil) == 0 {
		t.Fatal("nol tabel T_* terbaca dari migrasi; pembacanya yang rusak")
	}
	return hasil
}

var pembacaHilir = []struct {
	nama  string
	sql   string
	tabel string
	kolom []string
	// where adalah kolom yang dipakai penyaring, ikut wajib ada di DDL.
	where []string
}{
	{"KlausulUntukHilir", sqlKlausulHilir("S.T_PROPORTIONALARRG"), TabelKlausulTCO, KolomKlausulHilir,
		[]string{"TREATYDESCID", "TREATYYEAR", "TREATYGROUPID", "REINSTYPEID"}},
	{"KlausulIndukUntukHilir", sqlKlausulIndukHilir("S.T_PROPORTIONALARRG"), TabelKlausulTCO, KolomKlausulHilir,
		[]string{"PARENTREINSTYPEID"}},
	{"ReinsurerUntukHilir", sqlReinsurerHilir("S.T_TREATYREINSURER"), TabelReinsurerTCO, KolomReinsurerHilir,
		[]string{"REINSTYPEID", "TREATYYEAR", "TREATYGROUPID"}},
	{"BusinessUntukHilir", sqlBusinessHilir("S.T_TREATYBUSINESS"), TabelBusinessTCO, KolomBusinessHilir,
		[]string{"BIZCODE"}},
	{"GrupTreatyAktifUntukHilir", sqlGrupTreatyAktifHilir("S.T_TREATYBUSINESS"), TabelBusinessTCO,
		[]string{"TREATYGROUPID"}, []string{"BIZCODE", "TREATYYEAR", "ISACTIVE"}},
}

func TestKontrakHilirTCOKolomAdaDiDDL(t *testing.T) {
	ddl := kolomDDLTCO(t)
	for _, p := range pembacaHilir {
		kolomTabel := ddl[p.tabel]
		if len(kolomTabel) == 0 {
			t.Fatalf("%s: tabel %s tidak ada di DDL", p.nama, p.tabel)
		}
		for _, k := range append(append([]string{}, p.kolom...), p.where...) {
			if !kolomTabel[k] {
				t.Errorf("%s menyebut kolom %s yang tidak ada di %s", p.nama, k, p.tabel)
			}
			if !strings.Contains(p.sql, k) {
				t.Errorf("%s: SQL tidak menyebut kolom %s", p.nama, k)
			}
		}
	}
	// Pembaca gabungan menyentuh tiga tabel.
	gabung := sqlLimitTreatyHilir("S.T_TREATYBUSINESS", "S.T_PROPORTIONALARRG", "S.T_TREATYCONTRACT")
	for tabel, kolom := range map[string][]string{
		TabelBusinessTCO: {"TREATYYEARID", "REINSTYPEID", "BIZCODE"},
		TabelKlausulTCO:  {"TREATYYEARID", "REINSTYPEID", "TREATYDESCID", "PCT", "RP", "USD"},
		TabelKontrakTCO:  {"IDTREATYYEAR", "REINSTYPEID", "TREATYSTARTDATE", "TREATYENDDATE"},
	} {
		for _, k := range kolom {
			if !ddl[tabel][k] {
				t.Errorf("LimitTreatyUntukHilir: kolom %s tidak ada di %s", k, tabel)
			}
			if !strings.Contains(gabung, k) {
				t.Errorf("LimitTreatyUntukHilir: SQL tidak menyebut %s", k)
			}
		}
	}
}

func TestKontrakHilirTCOHanyaMembaca(t *testing.T) {
	sumber, ada := cariBerkas(berkasSumberProduksi(t), "repository/tco_kontrak_hilir.go")
	if !ada {
		t.Fatal("tco_kontrak_hilir.go tidak terbaca")
	}
	kode := strings.ToUpper(buangKomentarSumber("x.go", sumber))
	for _, tulis := range []string{"INSERT ", "UPDATE ", "DELETE ", "MERGE "} {
		if strings.Contains(kode, tulis) {
			t.Errorf("pembaca hilir memuat %q - hilir READ-ONLY (AC 1, 2)", tulis)
		}
	}
	// Literal jenis klausul hilir TIDAK ditanam (parameter pemanggil).
	if strings.Contains(kode, "'10001'") {
		t.Error("'10001' ditanam sebagai literal; ia parameter milik hilir")
	}
}

// Sisi korpus: kolom pembaca kita ADA di kueri hilir yang sebenarnya.
func TestKontrakHilirTCOVerbatimTerhadapKorpus(t *testing.T) {
	const korpus = `D:\XML\RNM_BRD`
	if _, err := os.Stat(korpus); err != nil {
		t.Skipf("korpus tidak terjangkau: %v", err)
	}
	polaSQL := regexp.MustCompile(`(?s)<pyBrowseSQL>(.*?)</pyBrowseSQL>`)
	bacaSQL := func(relatif string) string {
		isi, err := os.ReadFile(filepath.Join(korpus, filepath.FromSlash(relatif)))
		if err != nil {
			t.Fatalf("%s: %v", relatif, err)
		}
		m := polaSQL.FindStringSubmatch(string(isi))
		if m == nil {
			t.Fatalf("%s: pyBrowseSQL tidak ditemukan", relatif)
		}
		return strings.ToUpper(m[1])
	}
	kasus := []struct {
		berkas string
		kolom  []string
		tabel  string
	}{
		{"Claim Prop/RDBList/GetLimitPLATreatyin.xml", KolomKlausulHilir, "PROPORTIONALARRG"},
		{"Claim Fac In/RDBList/GetLimitPLADLA_Sql.xml", KolomKlausulHilir, "PROPORTIONALARRG"},
		{"Claim Fac In/RDBList/GetQuotaShare.xml", append(append([]string{}, KolomKlausulHilir...), "PARENTREINSTYPEID"), "PROPORTIONALARRG"},
		{"Claim Prop/RDBList/GetListRetro_Sql.xml", KolomReinsurerHilir, "TREATYREINSURER"},
		{"Komite Claim Prop/RDBList/GetListRetro_Sql.xml", KolomReinsurerHilir, "TREATYREINSURER"},
		{"Claim Fac In/RDBList/GetListRetro_Sql.xml", KolomReinsurerHilir, "TREATYREINSURER"},
		{"Claim Fac In/RDBList/GetTreatyGroup_Sql.xml", KolomBusinessHilir, "TREATYBUSINESS"},
		{"Claim Prop/RDBList/GetTreatyGroupID.xml", []string{"TREATYGROUPID", "BIZCODE", "TREATYYEAR", "ISACTIVE"}, "TREATYBUSINESS"},
		{"Claim Fac In/RDBList/GetDataTreatyLimit_Sql.xml",
			[]string{"TREATYYEARID", "REINSTYPEID", "BIZCODE", "TREATYDESCID", "IDTREATYYEAR", "TREATYSTARTDATE", "TREATYENDDATE"},
			"PROPORTIONALARRG"},
	}
	for _, k := range kasus {
		teks := bacaSQL(k.berkas)
		if !strings.Contains(teks, k.tabel) {
			t.Errorf("%s tidak membaca %s", k.berkas, k.tabel)
		}
		for _, kolom := range k.kolom {
			if !strings.Contains(teks, kolom) {
				t.Errorf("%s: kolom %s tidak ada di kueri hilir - kontrak kita mengarang kolom", k.berkas, kolom)
			}
		}
	}
}
