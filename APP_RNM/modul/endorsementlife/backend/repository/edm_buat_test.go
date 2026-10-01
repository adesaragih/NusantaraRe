package repository

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"nusantarare/modul/endorsementlife/backend/models"
)

const (
	uSpreading = "UJISKEMA.T_PREMIUM_LIST_SPREADING"
	uRetro     = "UJISKEMA.T_PREMIUM_LIST_SPREADING_RETRO"
	uWarisan   = "UJISKEMA.M_LIFE_PREMIUM_DETAIL"
	uRekapW    = "UJISKEMA.M_LIFE_PREMIUM_SUMMARY"
)

func TestSQLKasusTerbukaDanSisipKasus(t *testing.T) {
	q := sqlKasusTerbuka(uPolis)
	penampungUnik(t, "sqlKasusTerbuka", q)
	if !strings.Contains(q, "p.OLD_POLICY_NO = :1 AND p.ID LIKE :2 AND p.STATUSS IS NULL") {
		t.Errorf("gerbang 3 bukan kasus EDM terbuka atas polis: %s", q)
	}
	s := sqlSisipKasus(uPolis)
	if n := penampungUnik(t, "sqlSisipKasus", s); n != 8+len(models.KolomKepalaSalin) {
		t.Fatalf("%d penampung", n)
	}
	for _, wajib := range []string{"OLD_POLICY_NO", "EDM_TYPE", "TO_DATE(:6, 'YYYY-MM-DD')", "PROD_KE", "CREATE_OP_NAME", "SYSDATE"} {
		if !strings.Contains(s, wajib) {
			t.Errorf("sqlSisipKasus tanpa %q", wajib)
		}
	}
	if strings.Contains(s, "NO_POLIS") || strings.Contains(s, "STATUSS") {
		t.Error("kasus baru tidak boleh mengisi NO_POLIS/STATUSS: ia belum versi resmi dan masih terbuka")
	}
}

func TestSQLSalinVersiAplikasi(t *testing.T) {
	q := sqlSalinPesertaAplikasi(uPeserta)
	if n := penampungUnik(t, "sqlSalinPesertaAplikasi", q); n != 6 {
		t.Fatalf("%d penampung", n)
	}
	for _, wajib := range []string{
		"RAWTOHEX(STANDARD_HASH(:1 || '/D/' || d.ID, 'MD5'))", ":2, d.ID, :3, d.PL_NUMBER, :4",
		"d.PREMIUM_LIST_ID = :5", "(d.EDM_STATUS IS NULL OR d.EDM_STATUS <> :6)",
	} {
		if !strings.Contains(q, wajib) {
			t.Errorf("salin peserta tanpa %q", wajib)
		}
	}
	for _, k := range kolomNilaiPeserta {
		if !strings.Contains(q, "d."+k) {
			t.Errorf("salin peserta melewatkan %s", k)
		}
	}
	sp := sqlSalinSpreading(uSpreading, uPeserta)
	penampungUnik(t, "sqlSalinSpreading", sp)
	if !strings.Contains(sp, "STANDARD_HASH(:2 || '/D/' || s.DETAIL_ID, 'MD5')") || !strings.Contains(sp, "d.EDM_STATUS <> :4") {
		t.Errorf("spreading tidak menunjuk peserta versi baru / tidak menyaring Delete: %s", sp)
	}
	rt := sqlSalinSpreadingRetro(uRetro, uSpreading, uPeserta)
	penampungUnik(t, "sqlSalinSpreadingRetro", rt)
	if !strings.Contains(rt, "STANDARD_HASH(:2 || '/S/' || r.SPREADING_ID, 'MD5')") {
		t.Errorf("retro tidak menunjuk spreading versi baru: %s", rt)
	}
}

func TestSQLSalinVersiWarisanBerindex(t *testing.T) {
	q := sqlSalinPesertaWarisan(uPeserta, uWarisan)
	if n := penampungUnik(t, "sqlSalinPesertaWarisan", q); n != 6 {
		t.Fatalf("%d penampung", n)
	}
	for _, wajib := range []string{"m.PL_NUMBER = :5 AND m.IDPEGA = :6", ":2, NULL, :3", "m.PRORATETYPE",
		"TO_CHAR(m.STNC, 'DD/MM/YYYY')", "TO_CHAR(m.WPC, 'DD/MM/YYYY')"} {
		if !strings.Contains(q, wajib) {
			t.Errorf("salin warisan tanpa %q", wajib)
		}
	}
	if strings.Contains(q, "m.PRO_RATE_TYPE") {
		t.Error("PRO_RATE_TYPE bernama PRORATETYPE di tabel warisan")
	}
}

func TestSQLKepalaWarisanMembacaDataJSON(t *testing.T) {
	q, err := sqlKepalaWarisan(uJSON)
	if err != nil {
		t.Fatal(err)
	}
	penampungUnik(t, "sqlKepalaWarisan", q)
	for _, k := range models.KolomKepalaSalin {
		if !strings.Contains(q, "JSON_VALUE(j.DATA_JSON, '$."+k.Properti+"' NULL ON ERROR)") {
			t.Errorf("kepala warisan tidak membaca %s", k.Properti)
		}
	}
	if !strings.Contains(q, "j.IDPEGA = :1 AND j.NOPOLIS = :2") {
		t.Error("kepala warisan tidak berkunci IDPEGA + NOPOLIS")
	}
}

func TestTanggalPega(t *testing.T) {
	for masuk, mau := range map[string]string{
		"20240115": "2024-01-15", "20240115T170000.000 GMT": "2024-01-15", "2024-01-15": "2024-01-15", "": "",
	} {
		if g, ok := TanggalPega(masuk); !ok || g != mau {
			t.Errorf("TanggalPega(%q) = %q, %v", masuk, g, ok)
		}
	}
	for _, rusak := range []string{"15/01/2024", "kemarin", "2024011"} {
		if _, ok := TanggalPega(rusak); ok {
			t.Errorf("%q diterima", rusak)
		}
	}
}

func TestSQLRincianDanPolisLama(t *testing.T) {
	penampungUnik(t, "sqlRincianPeserta", sqlRincianPeserta(uPeserta))
	if !strings.Contains(sqlRincianPeserta(uPeserta), "d.ID = :1 AND d.PREMIUM_LIST_ID = :2") {
		t.Error("rincian peserta tidak dikurung kasusnya")
	}
	sp := sqlSpreadingPeserta(uSpreading, uRetro)
	penampungUnik(t, "sqlSpreadingPeserta", sp)
	if !strings.Contains(sp, "LEFT JOIN") || !strings.Contains(sp, "s.DETAIL_ID = :1") {
		t.Error("spreading peserta")
	}
	h := sqlPesertaWarisanHalaman(uWarisan)
	if n := penampungUnik(t, "sqlPesertaWarisanHalaman", h); n != 4 {
		t.Fatalf("%d penampung", n)
	}
	r := sqlRekapPolis(uRekapW)
	penampungUnik(t, "sqlRekapPolis", r)
	if !strings.Contains(r, "r.PL_NUMBER = :1") || !strings.Contains(r, "FETCH FIRST 500 ROWS ONLY") || strings.Contains(r, "UPPER(") {
		t.Errorf("rekap polis: %s", r)
	}
}

func TestSQLArasapasHanyaAdaTidak(t *testing.T) {
	q := sqlSudahDibayar()
	penampungUnik(t, "sqlSudahDibayar", q)
	if !strings.Contains(q, "ARASAPAS.DETAIL_INVOICE WHERE INV_INV_NO = :1 AND IVD_JR_ID = :2 FETCH FIRST 1 ROWS ONLY") ||
		strings.Contains(q, "*") || models.KodeJurnalPelunasan != "5" {
		t.Errorf("gerbang 5: %s", q)
	}
}

// berkasRepository - teks setiap berkas Go NON-uji paket ini.
func berkasRepository(t *testing.T) map[string]string {
	t.Helper()
	cocok, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	hasil := map[string]string{}
	for _, f := range cocok {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		hasil[f] = string(b)
	}
	if len(hasil) < 5 {
		t.Fatalf("hanya %d berkas terbaca", len(hasil))
	}
	return hasil
}

// TestArasapasHanyaDiSatuBerkas - spec §13 / AC 47: kueri skema Arasapas
// dikurung satu repository bertanda batas lintas sistem.
func TestArasapasHanyaDiSatuBerkas(t *testing.T) {
	var ada []string
	for nama, isi := range berkasRepository(t) {
		if strings.Contains(strings.ToUpper(isi), "ARASAPAS.") {
			ada = append(ada, nama)
		}
	}
	sort.Strings(ada)
	if strings.Join(ada, ",") != "edm_arasapas.go,edm_tabel.go" {
		t.Errorf("skema Arasapas disebut di %v; hanya edm_tabel.go (konstanta) dan edm_arasapas.go (kueri)", ada)
	}
}

// TestKueriWarisanBerindex - E4 brief gelombang 2: setiap fungsi yang
// menyentuh tabel peserta warisan (±66,8 juta baris) memakai kunci ber-index
// (`PL_NUMBER = :` → `_INDEX4`, atau `PL_NUMBER_EDM = :` → `_INDEX21`) dan
// MENYEBUT indeksnya di komentarnya. Fungsi baru yang lupa = merah.
func TestKueriWarisanBerindex(t *testing.T) {
	fset := token.NewFileSet()
	diperiksa := 0
	for nama, isi := range berkasRepository(t) {
		f, err := parser.ParseFile(fset, nama, isi, parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		for _, d := range f.Decls {
			fn, ok := d.(*ast.FuncDecl)
			if !ok || !strings.HasPrefix(fn.Name.Name, "sql") {
				continue
			}
			badan := isi[fset.Position(fn.Pos()).Offset:fset.Position(fn.End()).Offset]
			komentar := ""
			if fn.Doc != nil {
				komentar = fn.Doc.Text()
			}
			// Hanya pembangun SQL yang alias tabel warisannya `m` (konvensi berkas ini).
			if !regexp.MustCompile(`\bm\.(PL_NUMBER|IDPEGA|ID)\b`).MatchString(badan) {
				continue
			}
			diperiksa++
			if !regexp.MustCompile(`m\.PL_NUMBER(_EDM)? = :\d`).MatchString(badan) {
				t.Errorf("%s: kueri tabel peserta warisan tanpa kunci PL_NUMBER/PL_NUMBER_EDM", fn.Name.Name)
			}
			if !strings.Contains(komentar, "M_LIFE_PREMIUM_DETAIL_INDEX4") && !strings.Contains(komentar, "M_LIFE_PREMIUM_DETAIL_INDEX21") {
				t.Errorf("%s: komentar tidak menyebut index yang dipakai (E4)", fn.Name.Name)
			}
		}
	}
	if diperiksa < 3 {
		t.Fatalf("hanya %d pembangun SQL warisan diperiksa; pembacanya yang rusak", diperiksa)
	}
}

// kolomDDL membaca kolom satu CREATE TABLE migrasi PremiumList (urut DDL).
func kolomDDL(t *testing.T, berkas string) []string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "..", "premiumlistlife", "backend", "migrations", berkas))
	if err != nil {
		t.Fatalf("DDL PremiumList %s: %v", berkas, err)
	}
	var kolom []string
	for _, baris := range strings.Split(string(b), "\n") {
		baris = strings.TrimSpace(baris)
		if m := regexp.MustCompile(`^([A-Z][A-Z0-9_]*)\s+(VARCHAR2|NUMBER|DATE|CLOB|TIMESTAMP)`).FindStringSubmatch(baris); m != nil {
			kolom = append(kolom, m[1])
		}
	}
	return kolom
}

// TestKolomSalinanSamaDenganDDL - kolom yang PremiumList tambahkan kelak tidak
// boleh diam-diam tertinggal dari salinan versi endorsement.
func TestKolomSalinanSamaDenganDDL(t *testing.T) {
	cek := func(nama string, ddl, identitas, nilai []string) {
		t.Helper()
		gabung := append(append([]string{}, identitas...), nilai...)
		sort.Strings(gabung)
		d := append([]string{}, ddl...)
		sort.Strings(d)
		if strings.Join(gabung, ",") != strings.Join(d, ",") {
			t.Errorf("%s: kolom salinan ≠ DDL\n salinan %v\n DDL     %v", nama, gabung, d)
		}
	}
	cek("T_PREMIUM_LIST_DETAIL", kolomDDL(t, "052_t_premium_list_detail.sql"), kolomIdentitasPeserta, kolomNilaiPeserta)
	cek("T_PREMIUM_LIST_SPREADING", kolomDDL(t, "053_t_premium_list_spreading.sql"), []string{"ID", "DETAIL_ID"}, kolomNilaiSpreading)
	cek("T_PREMIUM_LIST_SPREADING_RETRO", kolomDDL(t, "054_t_premium_list_spreading_retro.sql"), []string{"ID", "SPREADING_ID"}, kolomNilaiSpreadingRetro)
	for _, k := range models.KolomPesertaRinci {
		found := false
		for _, n := range kolomNilaiPeserta {
			found = found || n == k.Nama
		}
		if !found {
			t.Errorf("kolom rincian %s bukan kolom nilai peserta", k.Nama)
		}
	}
}
