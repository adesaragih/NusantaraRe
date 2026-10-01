package repository

// Kode Claim Life sesudah migrasi 023 (keputusan work owner 01-10-2026).
// TANPA Oracle.
//
// Untuk apa berkas ini: kolom yang dibuang 023 tidak boleh tersisa di teks SQL
// mana pun - kompilator tidak membaca isi literal SQL, dan pernyataan yang
// masih menyebut kolom itu baru mati di ORA-00904 saat jalan. Dan `TypeKlaim`
// membaca rumah baru `TYPE`: header klaim.

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// polaLiteralSQL - literal backtick Go: tempat seluruh SQL modul ini ditulis.
var polaLiteralSQL = regexp.MustCompile("(?s)`[^`]*`")

func TestNolSQLMenyebutKolomYangDibuang023(t *testing.T) {
	dibuang := regexp.MustCompile(`\b(PY_POSITION|CASE_ID)\b`)
	diperiksa := 0
	for _, akar := range []string{"..", filepath.FromSlash("../../../komiteclaimlife/backend")} {
		err := filepath.Walk(akar, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && info.Name() == "migrations" {
				return filepath.SkipDir
			}
			if info.IsDir() || !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
				return nil
			}
			isi, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			diperiksa++
			// Baris komentar dibuang dulu: komentar memakai backtick untuk
			// menyebut nama kolom, dan itu bukan SQL.
			var kode []string
			for _, b := range strings.Split(string(isi), "\n") {
				if !strings.HasPrefix(strings.TrimSpace(b), "//") {
					kode = append(kode, b)
				}
			}
			for _, lit := range polaLiteralSQL.FindAllString(strings.Join(kode, "\n"), -1) {
				if m := dibuang.FindString(lit); m != "" {
					t.Errorf("%s: literal SQL masih menyebut %s (dibuang migrasi 023):\n%s",
						filepath.ToSlash(p), m, lit)
				}
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if diperiksa < 20 {
		t.Fatalf("hanya %d berkas Go terbaca; pemindainya yang rusak", diperiksa)
	}
}

func TestTypeKlaimMembacaHeaderKlaim(t *testing.T) {
	isi, err := os.ReadFile("klaimlife.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(isi)
	awal := strings.Index(s, "func (r *KlaimLife) TypeKlaim(")
	if awal < 0 {
		t.Fatal("TypeKlaim tidak ditemukan")
	}
	akhir := strings.Index(s[awal:], "\n}\n")
	tubuh := s[awal : awal+akhir]
	if !strings.Contains(tubuh, `Qualify("T_GENERAL_CLAIM")`) || !strings.Contains(tubuh, "SELECT TYPE FROM %s WHERE ID = :1") {
		t.Errorf("TypeKlaim tidak membaca T_GENERAL_CLAIM.TYPE (migrasi 023):\n%s", tubuh)
	}
}
