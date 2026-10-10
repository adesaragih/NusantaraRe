package repository

// Migrasi inti ricommlife dibaca pengurai PRODUKSI pelari (`migrasi.KolomCreateTable`, `KolomAlterTambah`,
// `BacaPerintahKatalog`, `KolomAlterBuang`) - pengurai yang sama dengan pra-terbang `-migrate` - dan berpasangan dengan
// jalur mundurnya: 924 (tabel flat, riwayat - dibuang 934) dan 931-934 (satu tabel per jenis data, keputusan work owner
// 08-10-2026, RALAT R1).

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

func TestMigrasi924TeruraiDanBerpasangan(t *testing.T) {
	maju := langkahInti(t, "924_ricomm_life", false)
	nama, kolom := migrasi.KolomCreateTable(maju[0])
	if nama != "RICOMM_LIFE" || !slices.Equal(kolom, KolomKomisi) {
		t.Fatalf("KolomCreateTable = %s %v, mau RICOMM_LIFE %v", nama, kolom, KolomKomisi)
	}
	ddl := satuBaris(maju[0])
	for _, mau := range []string{"ID VARCHAR2(10) NOT NULL", "IDUSEDBY VARCHAR2(10),", "USEDBY VARCHAR2(200),",
		"CONTRACT NUMBER(5),", "YEAR NUMBER(5),", "COMM NUMBER(38,8),", "CONSTRAINT PK_RICOMM_LIFE PRIMARY KEY (ID)"} {
		if !strings.Contains(ddl, mau) {
			t.Errorf("DDL tanpa %q:\n%s", mau, ddl)
		}
	}
	// Selain PK, semua kolom NULLABLE - satu-satunya CHECK Oracle (SYS_C...) = NOT NULL ID, tidak disalin 934.
	if strings.Count(ddl, "NOT NULL") != 1 {
		t.Errorf("NOT NULL selain PK:\n%s", ddl)
	}
	mundur := langkahInti(t, "924_ricomm_life", true)
	if len(mundur) != 2 || satuBaris(mundur[0]) != "DROP TABLE {skema}.RICOMM_LIFE CASCADE CONSTRAINTS" {
		t.Fatalf("mundur %q", mundur)
	}
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

// Pelindung gagal-keras jalur mundur 931/933 (pola 929_down): pernyataan PERTAMA merujuk JSONDATA (ORA-00904 bila
// kolom itu sudah dibuang langkah maju yang tidak tercatat), nol baris diubah, BUKAN blok (blok melewati diri).
func periksaPelindung(t *testing.T, kunci, tabel, drop string) {
	t.Helper()
	mundur := langkahInti(t, kunci, true)
	mau := "UPDATE {skema}." + tabel + " SET JSONDATA = JSONDATA WHERE 1 = 0"
	if len(mundur) != 2 || satuBaris(mundur[0]) != mau || satuBaris(mundur[1]) != drop {
		t.Fatalf("%s mundur %q", kunci, mundur)
	}
	if _, blok := migrasi.BacaPerintahKatalog(mundur[0]); blok {
		t.Errorf("%s: pelindung tidak boleh blok berpelindung", kunci)
	}
}

// 931: SATU ALTER … ADD ( (KolomAlterTambah) - tiga kolom view ringkasan; USEDBY 200 = salinan di rincian.
func TestMigrasi931KolomRingkasan(t *testing.T) {
	maju := langkahInti(t, "931_m_ricomm_life_summary_kolom", false)
	if len(maju) != 1 {
		t.Fatalf("931 %d pernyataan, mau 1", len(maju))
	}
	nama, kolom := migrasi.KolomAlterTambah(maju[0])
	if nama != TabelRingkasan || !slices.Equal(kolom, []string{"USEDBY", "MODIFIEDDATE", "OPERATORID"}) {
		t.Fatalf("KolomAlterTambah = %s %v", nama, kolom)
	}
	lebar(t, "931", satuBaris(maju[0]), map[string]string{"USEDBY": "VARCHAR2(200)", "MODIFIEDDATE": "VARCHAR2(50)", "OPERATORID": "VARCHAR2(200)"})
	periksaPelindung(t, "931_m_ricomm_life_summary_kolom", TabelRingkasan,
		"ALTER TABLE {skema}.M_RICOMM_LIFE_SUMMARY DROP (USEDBY, MODIFIEDDATE, OPERATORID)")
}

// 932: isi dari JSONDATA (blok) -> buang JSONDATA (blok, KolomAlterBuang) -> indeks nama -> DROP VIEW TERAKHIR;
// mundurnya aman diulang dan memulihkan view PERSIS.
func TestMigrasi932SatuTabel(t *testing.T) {
	maju := langkahInti(t, "932_m_ricomm_life_summary_satu_tabel", false)
	if len(maju) != 4 {
		t.Fatalf("932 %d pernyataan, mau 4", len(maju))
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
	buang, ok := migrasi.BacaPerintahKatalog(maju[1])
	if !ok || buang.Objek != "JSONDATA" {
		t.Fatalf("blok buang %+v %v", buang, ok)
	}
	if nama, kolom := migrasi.KolomAlterBuang(buang.Perintah); nama != TabelRingkasan || !slices.Equal(kolom, []string{"JSONDATA"}) {
		t.Errorf("KolomAlterBuang = %s %v", nama, kolom)
	}
	if satuBaris(maju[2]) != "CREATE INDEX {skema}.IX_M_RICOMM_LIFE_SUMMARY_NAMA ON {skema}.M_RICOMM_LIFE_SUMMARY (UPPER(TRIM(USEDBY)))" ||
		satuBaris(maju[3]) != "DROP VIEW {skema}.RICOMM_LIFE_SUMMARY" {
		t.Errorf("932 indeks / DROP VIEW %q %q", maju[2], maju[3])
	}
	turun := langkahInti(t, "932_m_ricomm_life_summary_satu_tabel", true)
	if len(turun) != 5 {
		t.Fatalf("932 mundur %d pernyataan, mau 5", len(turun))
	}
	for i, mau := range map[int]struct {
		objek string
		ada   bool
	}{0: {"IX_M_RICOMM_LIFE_SUMMARY_NAMA", true}, 1: {"JSONDATA", false}, 3: {"ENSURE_M_RICOMM_LIFE_SUMMARY_JSON", false}} {
		pk, ok := migrasi.BacaPerintahKatalog(turun[i])
		if !ok || pk.Tabel != TabelRingkasan || pk.Objek != mau.objek || pk.BilaAda != mau.ada {
			t.Errorf("932 mundur pernyataan %d: %+v %v", i, pk, ok)
		}
	}
	memuat(t, "932 mundur", strings.Join(turun, "\n"),
		"JSON_OBJECT('USEDBY' VALUE m.USEDBY, 'MODIFIEDDATE' VALUE m.MODIFIEDDATE, 'OPERATORID' VALUE m.OPERATORID ABSENT ON NULL RETURNING CLOB)",
		"CHECK (JSONDATA IS JSON)")
	if satuBaris(turun[4]) != "CREATE VIEW {skema}.RICOMM_LIFE_SUMMARY AS SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID FROM {skema}.M_RICOMM_LIFE_SUMMARY a" {
		t.Errorf("932 mundur view %q", turun[4])
	}
}

// 933: SATU ALTER … ADD ( - kelima kolom 924, tipe dan lebar SAMA.
func TestMigrasi933KolomKomisi(t *testing.T) {
	maju := langkahInti(t, "933_m_ricomm_life_kolom", false)
	if len(maju) != 1 {
		t.Fatalf("933 %d pernyataan, mau 1", len(maju))
	}
	nama, kolom := migrasi.KolomAlterTambah(maju[0])
	if nama != TabelKomisi || !slices.Equal(kolom, KolomKomisi[1:]) {
		t.Fatalf("KolomAlterTambah = %s %v", nama, kolom)
	}
	lebar(t, "933", satuBaris(maju[0]), map[string]string{"IDUSEDBY": "VARCHAR2(10)", "USEDBY": "VARCHAR2(200)",
		"CONTRACT": "NUMBER(5)", "YEAR": "NUMBER(5)", "COMM": "NUMBER(38,8)"})
	periksaPelindung(t, "933_m_ricomm_life_kolom", TabelKomisi,
		"ALTER TABLE {skema}.M_RICOMM_LIFE DROP (IDUSEDBY, USEDBY, CONTRACT, YEAR, COMM)")
}

// 934: isi dari JSON (blok) -> buang JSONDATA (blok) -> timpa dari flat -> sisip NOT EXISTS (SESUDAH buang JSONDATA) ->
// indeks -> DROP TABLE RICOMM_LIFE TERAKHIR; mundurnya membangun ulang bentuk 924 berisi data, aman diulang.
func TestMigrasi934SatuTabel(t *testing.T) {
	maju := langkahInti(t, "934_m_ricomm_life_satu_tabel", false)
	if len(maju) != 6 {
		t.Fatalf("934 %d pernyataan, mau 6", len(maju))
	}
	isi, ok := migrasi.BacaPerintahKatalog(maju[0])
	if !ok || isi.Tabel != TabelKomisi || isi.Objek != "JSONDATA" || !isi.BilaAda {
		t.Fatalf("blok isi %+v %v", isi, ok)
	}
	buang, ok := migrasi.BacaPerintahKatalog(maju[1])
	if !ok {
		t.Fatal("blok buang")
	}
	if nama, kolom := migrasi.KolomAlterBuang(buang.Perintah); nama != TabelKomisi || !slices.Equal(kolom, []string{"JSONDATA"}) {
		t.Errorf("KolomAlterBuang = %s %v", nama, kolom)
	}
	memuat(t, "934 timpa", maju[2], "UPDATE {skema}.M_RICOMM_LIFE m SET (IDUSEDBY, USEDBY, CONTRACT, YEAR, COMM) =",
		"FROM {skema}.RICOMM_LIFE f WHERE f.ID = m.ID", "WHERE EXISTS")
	memuat(t, "934 sisip", maju[3], "INSERT INTO {skema}.M_RICOMM_LIFE (ID, IDUSEDBY, USEDBY, CONTRACT, YEAR, COMM)",
		"WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_RICOMM_LIFE m WHERE m.ID = f.ID)")
	if satuBaris(maju[4]) != "CREATE INDEX {skema}.IX_M_RICOMM_LIFE_IDUSEDBY ON {skema}.M_RICOMM_LIFE (IDUSEDBY)" ||
		satuBaris(maju[5]) != "DROP TABLE {skema}.RICOMM_LIFE CASCADE CONSTRAINTS" {
		t.Errorf("934 indeks / DROP TABLE %q %q", maju[4], maju[5])
	}
	for _, p := range maju {
		if strings.Contains(p, "SYS_C") || strings.Contains(p, "JSON_OBJECT") {
			t.Errorf("934 maju menyalin CHECK sistem / membangun JSON: %s", p)
		}
	}
	turun := langkahInti(t, "934_m_ricomm_life_satu_tabel", true)
	if len(turun) != 7 {
		t.Fatalf("934 mundur %d pernyataan, mau 7", len(turun))
	}
	for i, mau := range map[int]struct {
		tabel, objek string
		ada          bool
	}{0: {"RICOMM_LIFE", "ID", false}, 1: {"RICOMM_LIFE", "IX_RICOMM_LIFE_IDUSEDBY", false},
		3: {TabelKomisi, "IX_M_RICOMM_LIFE_IDUSEDBY", true}, 4: {TabelKomisi, "JSONDATA", false}, 6: {TabelKomisi, "ENSURE_M_RICOMM_LIFE_JSON", false}} {
		pk, ok := migrasi.BacaPerintahKatalog(turun[i])
		if !ok || pk.Tabel != mau.tabel || pk.Objek != mau.objek || pk.BilaAda != mau.ada {
			t.Errorf("934 mundur pernyataan %d: %+v %v", i, pk, ok)
		}
	}
	// Bentuk 924 PERSIS: kolom, tipe, NOT NULL, dan PK sama dengan CREATE TABLE 924.
	buat, _ := migrasi.BacaPerintahKatalog(turun[0])
	ddl924 := satuBaris(langkahInti(t, "924_ricomm_life", false)[0])
	for _, bagian := range []string{"ID VARCHAR2(10) NOT NULL", "IDUSEDBY VARCHAR2(10)", "USEDBY VARCHAR2(200)", "CONTRACT NUMBER(5)",
		"YEAR NUMBER(5)", "COMM NUMBER(38,8)", "CONSTRAINT PK_RICOMM_LIFE PRIMARY KEY (ID)"} {
		if !strings.Contains(ddl924, bagian) || !strings.Contains(buat.Perintah, bagian) {
			t.Errorf("934 mundur CREATE tidak sama dengan 924 pada %q: %s", bagian, buat.Perintah)
		}
	}
	memuat(t, "934 mundur", strings.Join(turun, "\n"), "ID VARCHAR2(10) NOT NULL, IDUSEDBY VARCHAR2(10), USEDBY VARCHAR2(200), CONTRACT NUMBER(5), YEAR NUMBER(5), COMM NUMBER(38,8), CONSTRAINT PK_RICOMM_LIFE PRIMARY KEY (ID)",
		"WHERE NOT EXISTS (SELECT 1 FROM {skema}.RICOMM_LIFE f WHERE f.ID = m.ID)",
		"JSON_OBJECT('IDUSEDBY' VALUE m.IDUSEDBY, 'USEDBY' VALUE m.USEDBY, 'CONTRACT' VALUE m.CONTRACT, 'YEAR' VALUE m.YEAR, 'COMM' VALUE m.COMM ABSENT ON NULL RETURNING CLOB)")
}
