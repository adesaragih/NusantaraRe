package repository

// Migrasi inti ririsklife 935-941 (keputusan work owner 08-10-2026 K1-K4) dibaca pengurai PRODUKSI pelari
// (`migrasi.KolomAlterTambah`, `BacaPerintahKatalog`, `KolomAlterBuang`, `KolomCreateTable`) dan berpasangan dengan
// jalur mundurnya. Nol CREATE TABLE RIRISK_LIFE* (K1: RENAME, bukan tabel baru).

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

var semuaLangkah = []string{"935_ririsk_life_summary_ganti_nama", "936_ririsk_life_summary_kolom",
	"937_ririsk_life_summary_satu_tabel", "938_ririsk_life_ganti_nama", "939_ririsk_life_kolom",
	"940_ririsk_life_satu_tabel", "941_m_nav_menu_ririsklife"}

// K1: nol CREATE TABLE (pra-terbang pelari menganggap VIEW = objek ada; RENAME tidak membuat tabel baru) dan nol tabel
// lain (staging, _OLD, _BARU); M_RIRISK_LIFE_TEMP tidak disentuh.
func TestMigrasiRiRiskTanpaCreateTable(t *testing.T) {
	for _, k := range semuaLangkah {
		for _, mundur := range []bool{false, true} {
			for _, p := range langkahInti(t, k, mundur) {
				if nama, _ := migrasi.KolomCreateTable(p); nama != "" {
					t.Errorf("%s (mundur=%v) membuat tabel %s", k, mundur, nama)
				}
				if strings.Contains(p, "M_RIRISK_LIFE_TEMP") || regexp.MustCompile(`RIRISK_LIFE\w*_(OLD|BARU|TMP)`).MatchString(p) {
					t.Errorf("%s menyentuh tabel lain: %s", k, p)
				}
			}
		}
	}
}

// periksaGantiNama - 935 / 938: blok ALL_VIEWS (DROP VIEW hanya bila VIEW) lalu blok RENAME berpelindung tabel sumber;
// mundurnya: pelindung ORA-00904, RENAME balik berpelindung JSONDATA, CREATE VIEW persis teks DEV terakhir.
func periksaGantiNama(t *testing.T, kunci, lama, baru, view string) {
	t.Helper()
	maju := langkahInti(t, kunci, false)
	if len(maju) != 2 {
		t.Fatalf("%s %d pernyataan, mau 2", kunci, len(maju))
	}
	v, ok := migrasi.BacaPerintahKatalog(maju[0])
	if !ok || v != (migrasi.PerintahKatalog{Katalog: "ALL_VIEWS", Tabel: baru, Objek: baru, BilaAda: true, Perintah: "DROP VIEW {skema}." + baru}) {
		t.Errorf("%s blok view %+v %v", kunci, v, ok)
	}
	r, ok := migrasi.BacaPerintahKatalog(maju[1])
	if !ok || r != (migrasi.PerintahKatalog{Katalog: "ALL_TAB_COLUMNS", Tabel: lama, Objek: "ID", BilaAda: true,
		Perintah: "ALTER TABLE {skema}." + lama + " RENAME TO " + baru}) {
		t.Errorf("%s blok rename %+v %v", kunci, r, ok)
	}
	turun := langkahInti(t, kunci, true)
	if len(turun) != 3 || satuBaris(turun[0]) != "UPDATE {skema}."+baru+" SET JSONDATA = JSONDATA WHERE 1 = 0" {
		t.Fatalf("%s mundur %q", kunci, turun)
	}
	b, ok := migrasi.BacaPerintahKatalog(turun[1])
	if !ok || b != (migrasi.PerintahKatalog{Katalog: "ALL_TAB_COLUMNS", Tabel: baru, Objek: "JSONDATA", BilaAda: true,
		Perintah: "ALTER TABLE {skema}." + baru + " RENAME TO " + lama}) {
		t.Errorf("%s mundur rename %+v %v", kunci, b, ok)
	}
	if satuBaris(turun[2]) != view {
		t.Errorf("%s view dipulihkan\n%s\nmau\n%s", kunci, satuBaris(turun[2]), view)
	}
}

func TestMigrasi935938GantiNama(t *testing.T) {
	periksaGantiNama(t, "935_ririsk_life_summary_ganti_nama", "M_RIRISK_LIFE_SUMMARY", "RIRISK_LIFE_SUMMARY",
		"CREATE VIEW {skema}.RIRISK_LIFE_SUMMARY AS SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID FROM {skema}.M_RIRISK_LIFE_SUMMARY a")
	periksaGantiNama(t, "938_ririsk_life_ganti_nama", "M_RIRISK_LIFE", "RIRISK_LIFE",
		"CREATE VIEW {skema}.RIRISK_LIFE AS SELECT a.ID, a.JSONDATA.IDUSEDBY, a.JSONDATA.USEDBY, a.JSONDATA.AGE, a.JSONDATA.YEAR, a.JSONDATA.MONTH, a.JSONDATA.RISK, a.JSONDATA.CONTRACT FROM {skema}.M_RIRISK_LIFE a")
}

// lebar - `KOLOM TIPE` di DDL (satu baris).
func lebar(t *testing.T, nama, ddl string, mau map[string]string) {
	t.Helper()
	for k, tipe := range mau {
		if !regexp.MustCompile(`\b` + k + `\s+` + regexp.QuoteMeta(tipe) + `\s*[,)]`).MatchString(ddl) {
			t.Errorf("%s: %s bukan %s:\n%s", nama, k, tipe, ddl)
		}
	}
}

// periksaPelindung - mundur 936 / 939: pernyataan PERTAMA merujuk JSONDATA (ORA-00904 bila langkah maju yang tidak
// tercatat sudah membuangnya), bukan blok (blok melewati diri).
func periksaPelindung(t *testing.T, kunci, tabel, drop string) {
	t.Helper()
	turun := langkahInti(t, kunci, true)
	if len(turun) != 2 || satuBaris(turun[0]) != "UPDATE {skema}."+tabel+" SET JSONDATA = JSONDATA WHERE 1 = 0" || satuBaris(turun[1]) != drop {
		t.Fatalf("%s mundur %q", kunci, turun)
	}
	if _, blok := migrasi.BacaPerintahKatalog(turun[0]); blok {
		t.Errorf("%s: pelindung tidak boleh blok", kunci)
	}
}

// 936 / 939: SATU ALTER … ADD ( (KolomAlterTambah). Kolom baru ringkasan = pola ricommlife 931; AGE = teks seperti
// saudaranya. Kolom warisan (IDUSEDBY, USEDBY, YEAR, MONTH, RISK NUMBER, CONTRACT) TIDAK ditambah ulang.
func TestMigrasi936939Kolom(t *testing.T) {
	maju := langkahInti(t, "936_ririsk_life_summary_kolom", false)
	nama, kolom := migrasi.KolomAlterTambah(maju[0])
	if len(maju) != 1 || nama != TabelRingkasan || !slices.Equal(kolom, KolomRingkasan[1:]) {
		t.Fatalf("936 = %s %v", nama, kolom)
	}
	lebar(t, "936", satuBaris(maju[0]), map[string]string{"USEDBY": "VARCHAR2(200)", "MODIFIEDDATE": "VARCHAR2(50)", "OPERATORID": "VARCHAR2(200)"})
	periksaPelindung(t, "936_ririsk_life_summary_kolom", TabelRingkasan, "ALTER TABLE {skema}.RIRISK_LIFE_SUMMARY DROP (USEDBY, MODIFIEDDATE, OPERATORID)")

	maju = langkahInti(t, "939_ririsk_life_kolom", false)
	nama, kolom = migrasi.KolomAlterTambah(maju[0])
	if len(maju) != 1 || nama != TabelRincian || !slices.Equal(kolom, []string{"AGE"}) {
		t.Fatalf("939 = %s %v", nama, kolom)
	}
	lebar(t, "939", satuBaris(maju[0]), map[string]string{"AGE": "VARCHAR2(10)"})
	periksaPelindung(t, "939_ririsk_life_kolom", TabelRincian, "ALTER TABLE {skema}.RIRISK_LIFE DROP (AGE)")
}

// 937: isi dari JSONDATA (blok) -> buang JSONDATA (blok, KolomAlterBuang) -> indeks nama; mundurnya aman diulang.
func TestMigrasi937SatuTabel(t *testing.T) {
	maju := langkahInti(t, "937_ririsk_life_summary_satu_tabel", false)
	if len(maju) != 3 {
		t.Fatalf("937 %d pernyataan, mau 3", len(maju))
	}
	isi, ok := migrasi.BacaPerintahKatalog(maju[0])
	if !ok || isi.Tabel != TabelRingkasan || isi.Objek != "JSONDATA" || !isi.BilaAda {
		t.Fatalf("blok isi %+v %v", isi, ok)
	}
	for _, k := range KolomRingkasan[1:] {
		if !strings.Contains(isi.Perintah, "m."+k+" = m.JSONDATA."+k) {
			t.Errorf("blok isi tanpa %s apa adanya: %s", k, isi.Perintah)
		}
	}
	buang, _ := migrasi.BacaPerintahKatalog(maju[1])
	if nama, kolom := migrasi.KolomAlterBuang(buang.Perintah); nama != TabelRingkasan || !slices.Equal(kolom, []string{"JSONDATA"}) {
		t.Errorf("KolomAlterBuang = %s %v", nama, kolom)
	}
	if satuBaris(maju[2]) != "CREATE INDEX {skema}.IX_RIRISK_LIFE_SUMMARY_NAMA ON {skema}.RIRISK_LIFE_SUMMARY (UPPER(TRIM(USEDBY)))" {
		t.Errorf("937 indeks %q", maju[2])
	}
	turun := langkahInti(t, "937_ririsk_life_summary_satu_tabel", true)
	for i, mau := range map[int]struct {
		objek string
		ada   bool
	}{0: {"IX_RIRISK_LIFE_SUMMARY_NAMA", true}, 1: {"JSONDATA", false}, 3: {"ENSURE_M_RIRISK_LIFE_SUMMARY_JSON", false}} {
		pk, ok := migrasi.BacaPerintahKatalog(turun[i])
		if !ok || pk.Tabel != TabelRingkasan || pk.Objek != mau.objek || pk.BilaAda != mau.ada {
			t.Errorf("937 mundur %d: %+v %v", i, pk, ok)
		}
	}
	memuat(t, "937 mundur", strings.Join(turun, "\n"),
		"JSON_OBJECT('USEDBY' VALUE m.USEDBY, 'MODIFIEDDATE' VALUE m.MODIFIEDDATE, 'OPERATORID' VALUE m.OPERATORID ABSENT ON NULL RETURNING CLOB)")
}

// 940: isi ULANG semua kolom dari JSONDATA (timpa kolom datar basi; RISK tanpa NLS) -> PK (sebelum JSONDATA dibuang) ->
// buang JSONDATA -> indeks IDUSEDBY hanya bila belum ada indeks berkolom pertama IDUSEDBY (INDEX4). Mundurnya aman
// diulang.
func TestMigrasi940SatuTabel(t *testing.T) {
	maju := langkahInti(t, "940_ririsk_life_satu_tabel", false)
	if len(maju) != 4 {
		t.Fatalf("940 %d pernyataan, mau 4", len(maju))
	}
	isi, ok := migrasi.BacaPerintahKatalog(maju[0])
	if !ok || isi.Tabel != TabelRincian || isi.Objek != "JSONDATA" || !isi.BilaAda {
		t.Fatalf("blok isi %+v %v", isi, ok)
	}
	for _, k := range []string{"IDUSEDBY", "USEDBY", "AGE", "YEAR", "MONTH", "CONTRACT"} {
		if !strings.Contains(isi.Perintah, "m."+k+" = m.JSONDATA."+k) {
			t.Errorf("blok isi tanpa %s dari JSON: %s", k, isi.Perintah)
		}
	}
	// RISK: koma DAN titik -> pemisah desimal sesi sendiri, TO_NUMBER tanpa format (tanpa bergantung NLS).
	memuat(t, "940 RISK", isi.Perintah, "m.RISK = TO_NUMBER(TRANSLATE(TRIM(m.JSONDATA.RISK), CHR(44) || CHR(46), "+
		"SUBSTR(TO_CHAR(0.5), 1, 1) || SUBSTR(TO_CHAR(0.5), 1, 1)))")
	pk, ok := migrasi.BacaPerintahKatalog(maju[1])
	if !ok || pk.Katalog != "ALL_CONSTRAINTS" || pk.Objek != "PK_RIRISK_LIFE" || pk.BilaAda ||
		pk.Perintah != "ALTER TABLE {skema}.RIRISK_LIFE ADD CONSTRAINT PK_RIRISK_LIFE PRIMARY KEY (ID)" {
		t.Errorf("940 PK %+v %v", pk, ok)
	}
	buang, _ := migrasi.BacaPerintahKatalog(maju[2])
	if nama, kolom := migrasi.KolomAlterBuang(buang.Perintah); nama != TabelRincian || !slices.Equal(kolom, []string{"JSONDATA"}) {
		t.Errorf("KolomAlterBuang = %s %v", nama, kolom)
	}
	ix, ok := migrasi.BacaPerintahKatalog(maju[3])
	if !ok || ix != (migrasi.PerintahKatalog{Katalog: "ALL_IND_COLUMNS", Tabel: TabelRincian, Objek: "IDUSEDBY", BilaAda: false,
		Perintah: "CREATE INDEX {skema}.IX_RIRISK_LIFE_IDUSEDBY ON {skema}.RIRISK_LIFE (IDUSEDBY)"}) {
		t.Errorf("940 indeks %+v %v", ix, ok)
	}
	turun := langkahInti(t, "940_ririsk_life_satu_tabel", true)
	if len(turun) != 5 {
		t.Fatalf("940 mundur %d pernyataan, mau 5", len(turun))
	}
	for i, mau := range map[int]struct {
		objek string
		ada   bool
	}{0: {"IX_RIRISK_LIFE_IDUSEDBY", true}, 1: {"PK_RIRISK_LIFE", true}, 2: {"JSONDATA", false}, 4: {"ENSURE_M_RIRISK_LIFE_JSON", false}} {
		pk, ok := migrasi.BacaPerintahKatalog(turun[i])
		if !ok || pk.Tabel != TabelRincian || pk.Objek != mau.objek || pk.BilaAda != mau.ada {
			t.Errorf("940 mundur %d: %+v %v", i, pk, ok)
		}
	}
	memuat(t, "940 mundur", turun[3], "'RISK' VALUE TO_CHAR(m.RISK, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')",
		"'IDUSEDBY' VALUE m.IDUSEDBY", "'MONTH' VALUE m.MONTH", "ABSENT ON NULL RETURNING CLOB")
}

// 941 (K4): KODE / MODUL ririsklife, LABEL 'R/I Risk', MASTER TREATY URUTAN 11, DIMIGRASI '1'; mundurnya membuang hak lalu
// barisnya.
func TestMigrasi941Menu(t *testing.T) {
	maju := langkahInti(t, "941_m_nav_menu_ririsklife", false)
	if len(maju) != 1 || !strings.Contains(satuBaris(maju[0]),
		"'ririsklife', 'R/I Risk', 'MASTER TREATY', 'ririsklife', 11, '1' FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'ririsklife')") {
		t.Errorf("941 %q", maju)
	}
	if turun := langkahInti(t, "941_m_nav_menu_ririsklife", true); len(turun) != 2 ||
		satuBaris(turun[0]) != "DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'ririsklife'" {
		t.Errorf("941 mundur %q", turun)
	}
}
