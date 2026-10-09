package repository

// Migrasi inti benefitlife 942-945 (keputusan work owner 08-10-2026 K1-K4) dibaca pengurai PRODUKSI pelari
// (`migrasi.Daftar`, `BacaPerintahKatalog`, `KolomAlterTambah`, `KolomAlterBuang`, `KolomCreateTable`) dan berpasangan
// dengan jalur mundurnya. Nol CREATE TABLE BENEFIT_LIFE (K1: RENAME, bukan tabel baru).

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

// langkahInti - pernyataan migrasi inti `kunci` (maju atau mundur), dibaca dari berkasnya bersama pasangan `_down`.
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

var semuaLangkah = []string{"942_benefit_life_ganti_nama", "943_benefit_life_kolom", "944_benefit_life_satu_tabel",
	"945_m_nav_menu_benefitlife"}

// viewLama - teks view DEV (fakta WO 08-10-2026, bab 1 prompt), ber-{skema}: dibuat ulang PERSIS oleh 942_down.
const viewLama = "CREATE VIEW {skema}.BENEFIT_LIFE AS SELECT a.ID, a.JSONDATA.Benefit FROM {skema}.M_BENEFIT_LIFE a"

// Objek milik aplikasi lain yang TIDAK boleh disentuh (prompt bab 1).
var polaJanganDisentuh = regexp.MustCompile(`\b(M_BENEFIT|MBENEFIT\w*|VJ_\w*BENEFIT\w*|VH_\w*BENEFIT\w*|HCD_\w*BENEFIT\w*|` +
	`T_BENEFIT|DET_GROUP_BENEFIT|M_PLAN_BENEFIT|M_PROPERTY_BENEFIT)\b`)

// K1: nol CREATE TABLE, nol tabel lain (staging, _OLD, _BARU), nol objek "jangan disentuh".
func TestMigrasiBenefitTanpaTabelLain(t *testing.T) {
	for _, k := range semuaLangkah {
		for _, mundur := range []bool{false, true} {
			for _, p := range langkahInti(t, k, mundur) {
				if nama, _ := migrasi.KolomCreateTable(p); nama != "" {
					t.Errorf("%s (mundur=%v) membuat tabel %s", k, mundur, nama)
				}
				if regexp.MustCompile(`BENEFIT_LIFE\w*_(OLD|BARU|TMP|TEMP)`).MatchString(p) || polaJanganDisentuh.MatchString(p) {
					t.Errorf("%s menyentuh objek lain: %s", k, p)
				}
			}
		}
	}
}

// 942: blok ALL_VIEWS (DROP VIEW hanya bila VIEW) lalu blok RENAME berpelindung tabel sumber; mundurnya: pelindung
// ORA-00904, RENAME balik berpelindung JSONDATA, CREATE VIEW PERSIS teks DEV.
func TestMigrasi942GantiNama(t *testing.T) {
	maju := langkahInti(t, "942_benefit_life_ganti_nama", false)
	if len(maju) != 2 {
		t.Fatalf("942 %d pernyataan, mau 2", len(maju))
	}
	v, ok := migrasi.BacaPerintahKatalog(maju[0])
	if !ok || v != (migrasi.PerintahKatalog{Katalog: "ALL_VIEWS", Tabel: Tabel, Objek: Tabel, BilaAda: true, Perintah: "DROP VIEW {skema}.BENEFIT_LIFE"}) {
		t.Errorf("blok view %+v %v", v, ok)
	}
	r, ok := migrasi.BacaPerintahKatalog(maju[1])
	if !ok || r != (migrasi.PerintahKatalog{Katalog: "ALL_TAB_COLUMNS", Tabel: "M_BENEFIT_LIFE", Objek: "ID", BilaAda: true,
		Perintah: "ALTER TABLE {skema}.M_BENEFIT_LIFE RENAME TO BENEFIT_LIFE"}) {
		t.Errorf("blok rename %+v %v", r, ok)
	}
	turun := langkahInti(t, "942_benefit_life_ganti_nama", true)
	if len(turun) != 3 || satuBaris(turun[0]) != "UPDATE {skema}.BENEFIT_LIFE SET JSONDATA = JSONDATA WHERE 1 = 0" {
		t.Fatalf("942 mundur %q", turun)
	}
	if _, blok := migrasi.BacaPerintahKatalog(turun[0]); blok {
		t.Error("pelindung ORA-00904 tidak boleh blok (blok melewati diri)")
	}
	b, ok := migrasi.BacaPerintahKatalog(turun[1])
	if !ok || b != (migrasi.PerintahKatalog{Katalog: "ALL_TAB_COLUMNS", Tabel: Tabel, Objek: "JSONDATA", BilaAda: true,
		Perintah: "ALTER TABLE {skema}.BENEFIT_LIFE RENAME TO M_BENEFIT_LIFE"}) {
		t.Errorf("mundur rename %+v %v", b, ok)
	}
	if satuBaris(turun[2]) != viewLama {
		t.Errorf("view dipulihkan\n%s\nmau\n%s", satuBaris(turun[2]), viewLama)
	}
}

// 943: SATU ALTER … ADD ( (KolomAlterTambah) - BENEFIT VARCHAR2(200) (K1: pola ririsklife USEDBY); mundurnya pelindung
// ORA-00904 lalu DROP (BENEFIT).
func TestMigrasi943Kolom(t *testing.T) {
	maju := langkahInti(t, "943_benefit_life_kolom", false)
	nama, kolom := migrasi.KolomAlterTambah(maju[0])
	if len(maju) != 1 || nama != Tabel || !slices.Equal(kolom, KolomTabel[1:]) {
		t.Fatalf("943 = %s %v", nama, kolom)
	}
	if !regexp.MustCompile(`\bBENEFIT\s+VARCHAR2\(200\)\s*\)`).MatchString(satuBaris(maju[0])) {
		t.Errorf("943 lebar: %s", satuBaris(maju[0]))
	}
	turun := langkahInti(t, "943_benefit_life_kolom", true)
	if len(turun) != 2 || satuBaris(turun[0]) != "UPDATE {skema}.BENEFIT_LIFE SET JSONDATA = JSONDATA WHERE 1 = 0" ||
		satuBaris(turun[1]) != "ALTER TABLE {skema}.BENEFIT_LIFE DROP (BENEFIT)" {
		t.Fatalf("943 mundur %q", turun)
	}
}

// Ekspresi pengisian 944 (K2) dan pemeriksaannya - persis seperti yang dikirim pelari.
const (
	perintahIsi      = "UPDATE {skema}.BENEFIT_LIFE m SET m.BENEFIT = m.JSONDATA.Benefit"
	perintahPeriksa  = "UPDATE {skema}.BENEFIT_LIFE m SET m.ID = NULL WHERE m.JSONDATA.Benefit IS NOT NULL AND (m.BENEFIT IS NULL OR m.BENEFIT <> m.JSONDATA.Benefit)"
	perintahBuang    = "ALTER TABLE {skema}.BENEFIT_LIFE DROP COLUMN JSONDATA CASCADE CONSTRAINTS"
	ekspresiViewLama = "a.JSONDATA.Benefit"
)

// 944 (K2): isi -> pemeriksaan -> buang, ketiganya blok ALL_TAB_COLUMNS BENEFIT_LIFE.JSONDATA n > 0 (aman diulang).
// Ekspresi isi = ekspresi view lama dengan alias berbeda SAJA (peka huruf: `Benefit`, bukan `BENEFIT` / `"benefit"`),
// jadi nilainya identik dengan yang view keluarkan.
func TestMigrasi944SatuTabel(t *testing.T) {
	maju := langkahInti(t, "944_benefit_life_satu_tabel", false)
	if len(maju) != 3 {
		t.Fatalf("944 %d pernyataan, mau 3", len(maju))
	}
	for i, mau := range []string{perintahIsi, perintahPeriksa, perintahBuang} {
		pk, ok := migrasi.BacaPerintahKatalog(maju[i])
		if !ok || pk != (migrasi.PerintahKatalog{Katalog: "ALL_TAB_COLUMNS", Tabel: Tabel, Objek: "JSONDATA", BilaAda: true, Perintah: mau}) {
			t.Errorf("944 blok %d: %+v %v", i, pk, ok)
		}
	}
	// Ekspresi isi = teks view lama (942_down) dengan alias `m` menggantikan `a` - huruf kunci persis.
	if !strings.Contains(viewLama, ekspresiViewLama) || !strings.Contains(perintahIsi, "m."+strings.TrimPrefix(ekspresiViewLama, "a.")) {
		t.Errorf("ekspresi isi %q tidak sama dengan view %q", perintahIsi, ekspresiViewLama)
	}
	for _, salah := range []string{"JSONDATA.BENEFIT", "JSONDATA.benefit", `JSONDATA."`} {
		if strings.Contains(perintahIsi+perintahPeriksa, salah) {
			t.Errorf("kunci JSON ditulis %q - notasi titik peka huruf, mau persis seperti view", salah)
		}
	}
	if nama, kolom := migrasi.KolomAlterBuang(perintahBuang); nama != Tabel || !slices.Equal(kolom, []string{"JSONDATA"}) {
		t.Errorf("KolomAlterBuang = %s %v", nama, kolom)
	}
	// Urutan: isi SEBELUM pemeriksaan SEBELUM buang.
	if !(strings.Index(strings.Join(maju, "\n"), perintahIsi) < strings.Index(strings.Join(maju, "\n"), "SET m.ID = NULL") &&
		strings.Index(strings.Join(maju, "\n"), "SET m.ID = NULL") < strings.Index(strings.Join(maju, "\n"), perintahBuang)) {
		t.Error("urutan 944 harus isi -> pemeriksaan -> buang")
	}
	turun := langkahInti(t, "944_benefit_life_satu_tabel", true)
	if len(turun) != 3 {
		t.Fatalf("944 mundur %d pernyataan, mau 3", len(turun))
	}
	for i, mau := range map[int]struct {
		objek string
		ada   bool
	}{0: {"JSONDATA", false}, 2: {"ENSURE_M_BENEFIT_LIFE_JSON", false}} {
		pk, ok := migrasi.BacaPerintahKatalog(turun[i])
		if !ok || pk.Tabel != Tabel || pk.Objek != mau.objek || pk.BilaAda != mau.ada {
			t.Errorf("944 mundur %d: %+v %v", i, pk, ok)
		}
	}
	// Kunci "Benefit" huruf persis - view yang dibuat ulang 942_down membacanya.
	memuat(t, "944 mundur", turun[1], "SET JSONDATA = JSON_OBJECT('Benefit' VALUE m.BENEFIT ABSENT ON NULL RETURNING CLOB)")
}

// K2 (berhenti sebelum DROP JSONDATA): bentuk yang diterima pengurai produksi TIDAK boleh bertanda kutip di teks
// EXECUTE IMMEDIATE - versi JSON_EXISTS(…, '$.Benefit') ditolak pengurai (jadi tidak dapat dipakai), versi notasi titik
// diterima. Pemeriksaan gagal keras lewat ORA-01407 (ID NOT NULL) - bukan galat yang ditoleransi pelari (hanya
// ORA-00955 pada CREATE maju, ORA-00942 / ORA-02289 mundur).
func TestBentukPemeriksaanK2(t *testing.T) {
	blok := func(perintah string) string {
		return "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS\n" +
			"   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'BENEFIT_LIFE' AND COLUMN_NAME = 'JSONDATA';\n" +
			"  IF n > 0 THEN\n    EXECUTE IMMEDIATE '" + perintah + "';\n  END IF;\nEND;"
	}
	if _, ok := migrasi.BacaPerintahKatalog(blok("UPDATE {skema}.BENEFIT_LIFE m SET m.ID = NULL WHERE JSON_EXISTS(m.JSONDATA, '$.Benefit') AND m.BENEFIT IS NULL")); ok {
		t.Error("blok bertanda kutip diterima pengurai - aturan polaPerintahKatalog berubah?")
	}
	if pk, ok := migrasi.BacaPerintahKatalog(blok(perintahPeriksa)); !ok || pk.Perintah != perintahPeriksa {
		t.Errorf("bentuk pemeriksaan ditolak pengurai: %+v %v", pk, ok)
	}
	// Pemeriksaan hanya MENGHENTIKAN: ia tidak menyebut kolom BENEFIT di SET (nilai tidak pernah diubah) dan targetnya
	// kolom NOT NULL (ID, PK SYS_C009031).
	if !strings.Contains(perintahPeriksa, "SET m.ID = NULL WHERE") || strings.Contains(perintahPeriksa, "SET m.BENEFIT") {
		t.Errorf("bentuk pemeriksaan %q", perintahPeriksa)
	}
}

// 945 (K4): KODE / MODUL benefitlife, LABEL 'Benefit', MASTER TREATY URUTAN 12, DIMIGRASI '1'; mundurnya membuang hak
// lalu barisnya.
func TestMigrasi945Menu(t *testing.T) {
	maju := langkahInti(t, "945_m_nav_menu_benefitlife", false)
	if len(maju) != 1 || !strings.Contains(satuBaris(maju[0]),
		"'benefitlife', 'Benefit', 'MASTER TREATY', 'benefitlife', 12, '1' FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'benefitlife')") {
		t.Errorf("945 %q", maju)
	}
	if turun := langkahInti(t, "945_m_nav_menu_benefitlife", true); len(turun) != 2 ||
		satuBaris(turun[0]) != "DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'benefitlife'" ||
		satuBaris(turun[1]) != "DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'benefitlife'" {
		t.Errorf("945 mundur %q", turun)
	}
}
