package repository

// Migrasi inti planlife 946-949 (keputusan work owner 08-10-2026 K1-K6) dibaca pengurai PRODUKSI pelari
// (`migrasi.Daftar`, `BacaPerintahKatalog`, `KolomAlterTambah`, `KolomAlterBuang`, `KolomCreateTable`) dan berpasangan
// dengan jalur mundurnya. Nol CREATE TABLE PRODUCT_TYPE_LIFE (K1: RENAME, bukan tabel baru).

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	"nusantarare/inti/backend/migrasi"
)

func langkahInti(t *testing.T, kunci string, mundur bool) []string {
	t.Helper()
	sumber := fstest.MapFS{}
	for _, n := range []string{kunci + ".sql", kunci + "_down.sql"} {
		isi, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "inti", "backend", "migrations", n))
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(isi), "\r") {
			t.Errorf("%s memuat CR", n)
		}
		sumber["migrations/"+n] = &fstest.MapFile{Data: isi}
	}
	l, err := migrasi.Daftar(mundur, sumber)
	if err != nil || len(l) != 1 || migrasi.KunciLangkah(l[0].Nama) != kunci {
		t.Fatalf("membaca %s (mundur=%v): %v %v", kunci, mundur, l, err)
	}
	return l[0].Pernyataan
}

var semuaLangkah = []string{"946_product_type_life_ganti_nama", "947_product_type_life_kolom",
	"948_product_type_life_satu_tabel", "949_m_nav_menu_planlife"}

// kunciView - kunci JSON view lama, huruf PERSIS (bab 1 prompt); kolom tujuan = huruf besar.
var kunciView = []string{"ID", "CoverName", "Business", "BusinessID", "Benefit", "BenefitID"}

// viewLama - teks view DEV ber-{skema}: dibuat ulang PERSIS oleh 946_down.
const viewLama = "CREATE VIEW {skema}.PRODUCT_TYPE_LIFE AS SELECT a.JSONDATA.ID, a.JSONDATA.CoverName, a.JSONDATA.Business, " +
	"a.JSONDATA.BusinessID, a.JSONDATA.Benefit, a.JSONDATA.BenefitID FROM {skema}.M_PRODUCT_TYPE_LIFE a"

const (
	paksaBerhenti = "SET m.COVERNAME = RPAD(CHR(88), 201, CHR(88)) WHERE "
	perintahIsi   = "UPDATE {skema}.PRODUCT_TYPE_LIFE m SET m.COVERNAME = m.JSONDATA.CoverName, m.BUSINESS = m.JSONDATA.Business, " +
		"m.BUSINESSID = m.JSONDATA.BusinessID, m.BENEFIT = m.JSONDATA.Benefit, m.BENEFITID = m.JSONDATA.BenefitID"
	perintahPeriksaID = "UPDATE {skema}.PRODUCT_TYPE_LIFE m " + paksaBerhenti + "m.ID IS NULL OR m.JSONDATA.ID IS NULL OR " +
		"m.ID <> m.JSONDATA.ID OR m.ID IN (SELECT d.ID FROM {skema}.PRODUCT_TYPE_LIFE d GROUP BY d.ID HAVING COUNT(*) > 1)"
	perintahPK    = "ALTER TABLE {skema}.PRODUCT_TYPE_LIFE ADD CONSTRAINT PK_PRODUCT_TYPE_LIFE PRIMARY KEY (ID)"
	perintahBuang = "ALTER TABLE {skema}.PRODUCT_TYPE_LIFE DROP COLUMN JSONDATA CASCADE CONSTRAINTS"
)

// perintahPeriksaIsi - pemeriksaan K2 kelima kolom (dibangun dari kunciView supaya hurufnya sama dengan view).
func perintahPeriksaIsi() string {
	var syarat []string
	for _, k := range kunciView[1:] {
		kol := strings.ToUpper(k)
		syarat = append(syarat, "(m.JSONDATA."+k+" IS NOT NULL AND (m."+kol+" IS NULL OR m."+kol+" <> m.JSONDATA."+k+"))")
	}
	return "UPDATE {skema}.PRODUCT_TYPE_LIFE m " + paksaBerhenti + strings.Join(syarat, " OR ")
}

// K1: nol CREATE TABLE, nol tabel lain, nol tulisan ke objek "jangan disentuh" (BUSINESS / M_BUSINESS / BENEFIT_LIFE /
// PLAN_LIFE_SUMMARY - kata BUSINESS / BENEFIT sebagai NAMA KOLOM PRODUCT_TYPE_LIFE dikecualikan lewat awalan `{skema}.`).
func TestMigrasiPlanTanpaTabelLain(t *testing.T) {
	for _, k := range semuaLangkah {
		for _, mundur := range []bool{false, true} {
			for _, p := range langkahInti(t, k, mundur) {
				if nama, _ := migrasi.KolomCreateTable(p); nama != "" {
					t.Errorf("%s membuat tabel %s", k, nama)
				}
				for _, o := range []string{"{skema}.BUSINESS", "{skema}.M_BUSINESS", "{skema}.BENEFIT_LIFE", "PLAN_LIFE_SUMMARY", "_OLD", "_BARU"} {
					if strings.Contains(p, o) {
						t.Errorf("%s menyentuh %s: %s", k, o, p)
					}
				}
			}
		}
	}
}

// 946: DROP VIEW berpelindung ALL_VIEWS, RENAME berpelindung tabel SUMBER (JSONDATA = tabel Pega asli); target TABLE
// yang sudah ada = ORA-00955 di dalam blok (tidak ditoleransi pelari). Mundur: pelindung ORA-00904, RENAME balik,
// CREATE VIEW persis teks DEV.
func TestMigrasi946GantiNama(t *testing.T) {
	maju := langkahInti(t, "946_product_type_life_ganti_nama", false)
	if len(maju) != 2 {
		t.Fatalf("946 %d pernyataan", len(maju))
	}
	if v, ok := migrasi.BacaPerintahKatalog(maju[0]); !ok || v != (migrasi.PerintahKatalog{Katalog: "ALL_VIEWS", Tabel: Tabel, Objek: Tabel,
		BilaAda: true, Perintah: "DROP VIEW {skema}.PRODUCT_TYPE_LIFE"}) {
		t.Errorf("blok view %+v %v", v, ok)
	}
	if r, ok := migrasi.BacaPerintahKatalog(maju[1]); !ok || r != (migrasi.PerintahKatalog{Katalog: "ALL_TAB_COLUMNS", Tabel: "M_PRODUCT_TYPE_LIFE",
		Objek: "JSONDATA", BilaAda: true, Perintah: "ALTER TABLE {skema}.M_PRODUCT_TYPE_LIFE RENAME TO PRODUCT_TYPE_LIFE"}) {
		t.Errorf("blok rename %+v %v", r, ok)
	}
	turun := langkahInti(t, "946_product_type_life_ganti_nama", true)
	if len(turun) != 3 || satuBaris(turun[0]) != "UPDATE {skema}.PRODUCT_TYPE_LIFE SET JSONDATA = JSONDATA WHERE 1 = 0" {
		t.Fatalf("946 mundur %q", turun)
	}
	if b, ok := migrasi.BacaPerintahKatalog(turun[1]); !ok || b.Perintah != "ALTER TABLE {skema}.PRODUCT_TYPE_LIFE RENAME TO M_PRODUCT_TYPE_LIFE" ||
		b.Objek != "JSONDATA" || !b.BilaAda {
		t.Errorf("mundur rename %+v %v", b, ok)
	}
	if satuBaris(turun[2]) != viewLama {
		t.Errorf("view\n%s\nmau\n%s", satuBaris(turun[2]), viewLama)
	}
}

// 947: SATU ALTER … ADD ( lima kolom, urutan view lama; lebar nama 200, ID master 10.
func TestMigrasi947Kolom(t *testing.T) {
	maju := langkahInti(t, "947_product_type_life_kolom", false)
	nama, kolom := migrasi.KolomAlterTambah(maju[0])
	if len(maju) != 1 || nama != Tabel || !slices.Equal(kolom, KolomTabel[1:]) {
		t.Fatalf("947 = %s %v", nama, kolom)
	}
	for k, tipe := range map[string]string{"COVERNAME": "VARCHAR2(200)", "BUSINESS": "VARCHAR2(200)", "BUSINESSID": "VARCHAR2(10)",
		"BENEFIT": "VARCHAR2(200)", "BENEFITID": "VARCHAR2(10)"} {
		if !regexp.MustCompile(`\b` + k + `\s+` + regexp.QuoteMeta(tipe) + `\s*[,)]`).MatchString(satuBaris(maju[0])) {
			t.Errorf("947 %s bukan %s", k, tipe)
		}
	}
	turun := langkahInti(t, "947_product_type_life_kolom", true)
	if len(turun) != 2 || satuBaris(turun[1]) != "ALTER TABLE {skema}.PRODUCT_TYPE_LIFE DROP (COVERNAME, BUSINESS, BUSINESSID, BENEFIT, BENEFITID)" {
		t.Fatalf("947 mundur %q", turun)
	}
}

// 948: isi -> periksa ID (K1.4) -> periksa isi (K2) -> PK -> buang; setiap blok terbaca pengurai produksi, persis.
// Ekspresi = teks view lama (alias berbeda saja, huruf kunci persis).
func TestMigrasi948SatuTabel(t *testing.T) {
	maju := langkahInti(t, "948_product_type_life_satu_tabel", false)
	if len(maju) != 5 {
		t.Fatalf("948 %d pernyataan, mau 5", len(maju))
	}
	mau := []migrasi.PerintahKatalog{
		{Katalog: "ALL_TAB_COLUMNS", Tabel: Tabel, Objek: "JSONDATA", BilaAda: true, Perintah: perintahIsi},
		{Katalog: "ALL_TAB_COLUMNS", Tabel: Tabel, Objek: "JSONDATA", BilaAda: true, Perintah: perintahPeriksaID},
		{Katalog: "ALL_TAB_COLUMNS", Tabel: Tabel, Objek: "JSONDATA", BilaAda: true, Perintah: perintahPeriksaIsi()},
		{Katalog: "ALL_CONSTRAINTS", Tabel: Tabel, Objek: "PK_PRODUCT_TYPE_LIFE", BilaAda: false, Perintah: perintahPK},
		{Katalog: "ALL_TAB_COLUMNS", Tabel: Tabel, Objek: "JSONDATA", BilaAda: true, Perintah: perintahBuang},
	}
	for i, m := range mau {
		if pk, ok := migrasi.BacaPerintahKatalog(maju[i]); !ok || pk != m {
			t.Errorf("948 blok %d:\n%+v %v\nmau\n%+v", i, pk, ok, m)
		}
	}
	// Setiap ekspresi isi = ekspresi view lama (a.JSONDATA.<k> -> m.JSONDATA.<k>), huruf persis.
	for _, k := range kunciView[1:] {
		if !strings.Contains(viewLama, "a.JSONDATA."+k+",") && !strings.Contains(viewLama, "a.JSONDATA."+k+" ") {
			t.Fatalf("kunci %s tidak di view", k)
		}
		if !strings.Contains(perintahIsi, "m."+strings.ToUpper(k)+" = m.JSONDATA."+k) {
			t.Errorf("isi %s bukan notasi titik view", k)
		}
	}
	if !strings.Contains(perintahPeriksaID, "m.JSONDATA.ID") || !strings.Contains(viewLama, "a.JSONDATA.ID,") {
		t.Error("ID dibandingkan dengan JSONDATA.ID (kolom ID view)")
	}
	for _, salah := range []string{"JSONDATA.COVERNAME", "JSONDATA.BENEFIT ", "JSONDATA.BUSINESSID", `JSONDATA."`} {
		if strings.Contains(strings.Join(maju, "\n"), salah) {
			t.Errorf("kunci JSON %q bukan huruf view", salah)
		}
	}
	if nama, kolom := migrasi.KolomAlterBuang(perintahBuang); nama != Tabel || !slices.Equal(kolom, []string{"JSONDATA"}) {
		t.Errorf("KolomAlterBuang %s %v", nama, kolom)
	}
	turun := langkahInti(t, "948_product_type_life_satu_tabel", true)
	if len(turun) != 4 {
		t.Fatalf("948 mundur %d", len(turun))
	}
	for i, m := range map[int]struct {
		objek string
		ada   bool
	}{0: {"PK_PRODUCT_TYPE_LIFE", true}, 1: {"JSONDATA", false}, 3: {"ENSURE_M_PRODUCT_TYPE_LIFE", false}} {
		if pk, ok := migrasi.BacaPerintahKatalog(turun[i]); !ok || pk.Tabel != Tabel || pk.Objek != m.objek || pk.BilaAda != m.ada {
			t.Errorf("948 mundur %d %+v %v", i, pk, ok)
		}
	}
	memuat(t, "948 mundur", turun[2], "JSON_OBJECT('ID' VALUE m.ID, 'CoverName' VALUE m.COVERNAME, 'Business' VALUE m.BUSINESS, "+
		"'BusinessID' VALUE m.BUSINESSID, 'Benefit' VALUE m.BENEFIT, 'BenefitID' VALUE m.BENEFITID ABSENT ON NULL RETURNING CLOB)")
}

// Bentuk pemaksa berhenti: tanpa tanda kutip (versi JSON_EXISTS bertanda kutip DITOLAK pengurai), tidak mengubah ID,
// target kolom yang lebarnya 200 diberi 201 huruf -> ORA-12899 (bukan galat yang ditoleransi pelari).
func TestBentukPemaksaBerhenti(t *testing.T) {
	blok := func(perintah string) string {
		return "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS\n" +
			"   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'PRODUCT_TYPE_LIFE' AND COLUMN_NAME = 'JSONDATA';\n" +
			"  IF n > 0 THEN\n    EXECUTE IMMEDIATE '" + perintah + "';\n  END IF;\nEND;"
	}
	if _, ok := migrasi.BacaPerintahKatalog(blok("UPDATE {skema}.PRODUCT_TYPE_LIFE m " + paksaBerhenti + "JSON_EXISTS(m.JSONDATA, '$.ID')")); ok {
		t.Error("blok bertanda kutip diterima pengurai")
	}
	for _, p := range []string{perintahPeriksaID, perintahPeriksaIsi()} {
		if pk, ok := migrasi.BacaPerintahKatalog(blok(p)); !ok || pk.Perintah != p {
			t.Errorf("pemeriksaan ditolak pengurai: %s", p)
		}
		if strings.Contains(p, "SET m.ID") || !strings.Contains(p, "RPAD(CHR(88), 201, CHR(88))") {
			t.Errorf("bentuk %s", p)
		}
	}
	if !regexp.MustCompile(`COVERNAME\s+VARCHAR2\(200\)`).MatchString(strings.Join(langkahInti(t, "947_product_type_life_kolom", false), " ")) {
		t.Error("201 huruf hanya memaksa ORA-12899 bila COVERNAME VARCHAR2(200)")
	}
}

// 949 (K6): planlife 'Plan' MASTER TREATY URUTAN 13 '1'; mundur hak lalu baris.
func TestMigrasi949Menu(t *testing.T) {
	maju := langkahInti(t, "949_m_nav_menu_planlife", false)
	if len(maju) != 1 || !strings.Contains(satuBaris(maju[0]),
		"'planlife', 'Plan', 'MASTER TREATY', 'planlife', 13, '1' FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'planlife')") {
		t.Errorf("949 %q", maju)
	}
	if turun := langkahInti(t, "949_m_nav_menu_planlife", true); len(turun) != 2 ||
		satuBaris(turun[0]) != "DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'planlife'" {
		t.Errorf("949 mundur %q", turun)
	}
}
