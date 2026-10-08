package repository

// Migrasi inti ringkasan R/I Rate Life dibaca pengurai PRODUKSI pelari (`migrasi.KolomCreateTable`,
// `KolomAlterTambah`, `KolomAlterBuang`, `BacaPerintahKatalog`) - pengurai yang sama dengan pra-terbang `-migrate`:
// 926 (tabel flat, sudah jalan di DEV - riwayat), 927 (kolom M_RATE_LIFE_SUMMARY), 928 (satu tabel, RALAT R6), dan
// pasangan `_down` masing-masing.

import (
	"fmt"
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

// lebarDDL - lebar VARCHAR2 setiap kolom ringkasan di teks DDL.
func periksaLebar(t *testing.T, nama, ddl string, kolom []string) {
	t.Helper()
	for _, k := range kolom {
		m := regexp.MustCompile(`\b` + k + ` VARCHAR2\((\d+)\)`).FindStringSubmatch(ddl)
		if m == nil || m[1] != fmt.Sprint(LebarKolomRingkasan[k]) {
			t.Errorf("%s %s: DDL %v, LebarKolomRingkasan %d", nama, k, m, LebarKolomRingkasan[k])
		}
	}
}

// 926 (riwayat, sudah jalan di DEV): tabel flat lima kolom, jalur mundur memulihkan view lengkap.
func TestMigrasi926Riwayat(t *testing.T) {
	maju := langkahInti(t, "926_rate_life_summary_flat", false)
	nama, kolom := migrasi.KolomCreateTable(maju[0])
	if nama != "RATE_LIFE_SUMMARY" || !slices.Equal(kolom, KolomRingkasan) {
		t.Fatalf("926 = %s %v", nama, kolom)
	}
	periksaLebar(t, "926", satuBaris(maju[0]), KolomRingkasan)
	mundur := langkahInti(t, "926_rate_life_summary_flat", true)
	mauView := "CREATE VIEW {skema}.RATE_LIFE_SUMMARY AS SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, " +
		"a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID, a.JSONDATA.FLAG FROM {skema}.M_RATE_LIFE_SUMMARY a"
	if len(mundur) != 2 || satuBaris(mundur[1]) != mauView {
		t.Errorf("926 mundur %q", mundur)
	}
}

// 927: SATU pernyataan ALTER … ADD ( biasa (terbaca KolomAlterTambah), keempat kolom ringkasan selebar 926; mundurnya
// membuang keempatnya.
func TestMigrasi927KolomRingkasan(t *testing.T) {
	maju := langkahInti(t, "927_m_rate_life_summary_kolom", false)
	if len(maju) != 1 {
		t.Fatalf("927 %d pernyataan, mau 1 (berdiri sendiri - ADD tidak aman diulang)", len(maju))
	}
	nama, kolom := migrasi.KolomAlterTambah(maju[0])
	if nama != TabelRingkasan || !slices.Equal(kolom, KolomRingkasan[1:]) {
		t.Fatalf("KolomAlterTambah = %s %v", nama, kolom)
	}
	periksaLebar(t, "927", satuBaris(maju[0]), KolomRingkasan[1:])
	mundur := langkahInti(t, "927_m_rate_life_summary_kolom", true)
	if len(mundur) != 1 || satuBaris(mundur[0]) != "ALTER TABLE {skema}.M_RATE_LIFE_SUMMARY DROP (USEDBY, TYPE, MODIFIEDDATE, OPERATORID)" {
		t.Errorf("927 mundur %q", mundur)
	}
}

// 928: urutan isi -> buang JSONDATA (blok berpelindung katalog, terbaca KolomAlterBuang) -> sisip baris flat-saja ->
// indeks -> DROP TABLE flat TERAKHIR; setiap pernyataan aman diulang. Mundurnya membangun ulang tabel flat dan JSONDATA.
func TestMigrasi928SatuTabel(t *testing.T) {
	maju := langkahInti(t, "928_m_rate_life_summary_satu_tabel", false)
	if len(maju) != 6 {
		t.Fatalf("928 %d pernyataan, mau 6", len(maju))
	}
	awal := []string{"UPDATE {skema}.M_RATE_LIFE_SUMMARY m SET (USEDBY, TYPE, MODIFIEDDATE, OPERATORID) =",
		"DELETE FROM {skema}.M_RATE_LIFE_SUMMARY m WHERE NOT EXISTS (SELECT 1 FROM {skema}.RATE_LIFE_SUMMARY f",
		"DECLARE", "INSERT INTO {skema}.M_RATE_LIFE_SUMMARY (ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID)",
		"CREATE INDEX {skema}.IX_M_RATE_LIFE_SUMMARY_NAMA ON {skema}.M_RATE_LIFE_SUMMARY (UPPER(TRIM(USEDBY)))",
		"DROP TABLE {skema}.RATE_LIFE_SUMMARY CASCADE CONSTRAINTS"}
	for i, a := range awal {
		if !strings.HasPrefix(satuBaris(maju[i]), a) {
			t.Errorf("928 pernyataan %d:\n%s\nmau berawal\n%s", i, satuBaris(maju[i]), a)
		}
	}
	pk, ok := migrasi.BacaPerintahKatalog(maju[2])
	if !ok || pk.Katalog != "ALL_TAB_COLUMNS" || pk.Tabel != TabelRingkasan || pk.Objek != "JSONDATA" || !pk.BilaAda {
		t.Errorf("blok buang JSONDATA %+v %v", pk, ok)
	}
	if nama, kolom := migrasi.KolomAlterBuang(pk.Perintah); nama != TabelRingkasan || !slices.Equal(kolom, []string{"JSONDATA"}) {
		t.Errorf("KolomAlterBuang = %s %v", nama, kolom)
	}
	if !strings.Contains(maju[3], "WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_RATE_LIFE_SUMMARY m WHERE m.ID = f.ID)") {
		t.Error("sisipan baris flat-saja harus NOT EXISTS (aman diulang)")
	}
	for _, p := range maju {
		if strings.Contains(p, "JSON_OBJECT") || strings.Contains(p, "FLAG") {
			t.Errorf("jalur maju membangun JSON / menyebut FLAG: %s", p)
		}
	}

	mundur := langkahInti(t, "928_m_rate_life_summary_satu_tabel", true)
	nama, kolom := migrasi.KolomCreateTable(mundur[0])
	if nama != "RATE_LIFE_SUMMARY" || !slices.Equal(kolom, KolomRingkasan) {
		t.Fatalf("928 mundur CREATE = %s %v", nama, kolom)
	}
	periksaLebar(t, "928_down", satuBaris(mundur[0]), KolomRingkasan)
	semua := satuBaris(strings.Join(mundur, "\n"))
	for _, mau := range []string{"CREATE INDEX {skema}.IX_RATE_LIFE_SUMMARY_NAMA ON {skema}.RATE_LIFE_SUMMARY",
		"INSERT INTO {skema}.RATE_LIFE_SUMMARY (ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID)",
		"DROP INDEX {skema}.IX_M_RATE_LIFE_SUMMARY_NAMA",
		"JSON_OBJECT('USEDBY' VALUE m.USEDBY, 'TYPE' VALUE m.TYPE, 'MODIFIEDDATE' VALUE m.MODIFIEDDATE, 'OPERATORID' VALUE m.OPERATORID ABSENT ON NULL RETURNING CLOB)",
		"ADD CONSTRAINT ENSURE_M_RATE_LIFE_SUMMARY_JSON CHECK (JSONDATA IS JSON)"} {
		if !strings.Contains(semua, mau) {
			t.Errorf("928 mundur tanpa %q", mau)
		}
	}
	if nama, kolom := migrasi.KolomAlterTambah(mundur[4]); nama != TabelRingkasan || !slices.Equal(kolom, []string{"JSONDATA"}) {
		t.Errorf("928 mundur ADD JSONDATA = %s %v", nama, kolom)
	}
}
