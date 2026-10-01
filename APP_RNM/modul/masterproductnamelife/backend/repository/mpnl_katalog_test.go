package repository

// Uji katalog (lanjutan 1, L1): setiap kolom yang DITULIS repository ada di katalog
// DEV (`testdata/katalog-dev.json`, `ALL_TAB_COLUMNS` 01-10-2026). Kolom yang tidak
// ada di DEV = `ORA-00904` pada setiap simpan - penjaga ini merah lebih dulu.

import (
	"encoding/json"
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// katalogDEV - kolom per tabel menurut katalog DEV.
func katalogDEV(t *testing.T) map[string][]string {
	t.Helper()
	isi, err := os.ReadFile("testdata/katalog-dev.json")
	if err != nil {
		t.Fatal(err)
	}
	var k struct {
		Tabel map[string][]string `json:"tabel"`
	}
	if err := json.Unmarshal(isi, &k); err != nil {
		t.Fatal(err)
	}
	return k.Tabel
}

var (
	polaSisip  = regexp.MustCompile(`(?s)INSERT\s+INTO\s+\S+\s*\(([^)]*)\)`)
	polaUbah   = regexp.MustCompile(`(?s)UPDATE\s+\S+\s+SET\s+(.*?)\s+WHERE\b`)
	polaSetKol = regexp.MustCompile(`(?:^|,)\s*([A-Z_][A-Z0-9_]*)\s*=`)
)

// kolomDitulis - kolom yang disebut INSERT (daftar kolom) atau UPDATE (klausa SET).
func kolomDitulis(q string) []string {
	var hasil []string
	if m := polaSisip.FindStringSubmatch(q); m != nil {
		for _, k := range strings.Split(m[1], ",") {
			hasil = append(hasil, strings.TrimSpace(k))
		}
	}
	if m := polaUbah.FindStringSubmatch(q); m != nil {
		for _, s := range polaSetKol.FindAllStringSubmatch(m[1], -1) {
			hasil = append(hasil, s[1])
		}
	}
	return hasil
}

// kolomAsing - kolom yang ditulis q tetapi tidak ada di katalog tabelnya.
func kolomAsing(q string, katalog []string) []string {
	ada := map[string]bool{}
	for _, k := range katalog {
		ada[k] = true
	}
	var asing []string
	for _, k := range kolomDitulis(q) {
		if !ada[k] {
			asing = append(asing, k)
		}
	}
	return asing
}

func TestKolomDitulisAdaDiKatalogDEV(t *testing.T) {
	kat := katalogDEV(t)
	kasus := []struct {
		tabel string
		sql   string
	}{
		{TabelProduk, sqlSisipUmum("S.T")},
		{TabelProduk, sqlPerbaruiUmum("S.T")},
		{TabelInward, sqlSisipInward("S.T")},
		{TabelInward, sqlPerbaruiInward("S.T")},
	}
	for _, k := range kasus {
		if len(kolomDitulis(k.sql)) == 0 {
			t.Errorf("%s: pembaca kolom tidak menemukan kolom apa pun - instrumennya yang rusak:\n%s", k.tabel, rata(k.sql))
		}
		if asing := kolomAsing(k.sql, kat[k.tabel]); len(asing) > 0 {
			t.Errorf("%s: kolom %v tidak ada di katalog DEV %v (ORA-00904):\n%s", k.tabel, asing, kat[k.tabel], rata(k.sql))
		}
	}
	// Daftar kolom penjaga kata cadangan = katalog persis.
	for tabel, kolom := range map[string][]string{TabelProduk: KolomProduk, TabelInward: KolomInward} {
		a, b := append([]string{}, kolom...), append([]string{}, kat[tabel]...)
		sort.Strings(a)
		sort.Strings(b)
		if strings.Join(a, ",") != strings.Join(b, ",") {
			t.Errorf("%s: daftar kolom %v ≠ katalog DEV %v", tabel, kolom, kat[tabel])
		}
	}
}

// Uji gigit: menambah PRODUCTNAME ke INSERT / UPDATE = merah.
func TestAturanKatalogMenggigit(t *testing.T) {
	kat := katalogDEV(t)[TabelProduk]
	sisip := strings.Replace(sqlSisipUmum("S.T"), "RIRISK)", "RIRISK, PRODUCTNAME)", 1)
	if asing := kolomAsing(sisip, kat); len(asing) != 1 || asing[0] != "PRODUCTNAME" {
		t.Errorf("INSERT ber-PRODUCTNAME harus tertangkap: %v\n%s", asing, rata(sisip))
	}
	ubah := strings.Replace(sqlPerbaruiUmum("S.T"), " WHERE", ", BEGIN_DATE = TO_DATE(:9, 'DD/MM/YYYY') WHERE", 1)
	if asing := kolomAsing(ubah, kat); len(asing) != 1 || asing[0] != "BEGIN_DATE" {
		t.Errorf("UPDATE ber-BEGIN_DATE harus tertangkap: %v\n%s", asing, rata(ubah))
	}
	// Butir yang jawabannya diketahui - instrumen diuji lebih dulu.
	if got := kolomDitulis(`UPDATE S.T SET A = :1, B = TO_DATE(:2, 'DD/MM/YYYY') WHERE ID = :3`); strings.Join(got, ",") != "A,B" {
		t.Errorf("pembaca SET: %v", got)
	}
	if got := kolomDitulis("INSERT INTO S.T (ID, X)\n VALUES (:1, :2)"); strings.Join(got, ",") != "ID,X" {
		t.Errorf("pembaca INSERT: %v", got)
	}
}
