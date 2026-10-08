package repository

// Migrasi inti 924 (tabel flat RICOMM_LIFE) dibaca pengurai PRODUKSI pelari (`migrasi.KolomCreateTable`) - pengurai
// yang sama dengan pra-terbang `-migrate` - dan berpasangan dengan jalur mundurnya (keputusan work owner 06-10-2026
// butir 2).

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	"nusantarare/inti/backend/migrasi"
)

const (
	berkas924     = "924_ricomm_life.sql"
	berkas924Down = "924_ricomm_life_down.sql"
)

func langkah924(t *testing.T, mundur bool) []string {
	t.Helper()
	sumber := fstest.MapFS{}
	for _, n := range []string{berkas924, berkas924Down} {
		isi, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "inti", "backend", "migrations", n))
		if err != nil {
			t.Fatal(err)
		}
		sumber["migrations/"+n] = &fstest.MapFile{Data: isi}
	}
	l, err := migrasi.Daftar(mundur, sumber)
	if err != nil || len(l) != 1 {
		t.Fatalf("membaca 924 (mundur=%v): %d langkah, %v", mundur, len(l), err)
	}
	if migrasi.KunciLangkah(l[0].Nama) != "924_ricomm_life" {
		t.Fatalf("kunci langkah %s", l[0].Nama)
	}
	return l[0].Pernyataan
}

func TestMigrasi924TeruraiDanBerpasangan(t *testing.T) {
	maju := langkah924(t, false)
	nama, kolom := migrasi.KolomCreateTable(maju[0])
	if nama != TabelKomisi || !slices.Equal(kolom, KolomViewKomisiLama) {
		t.Fatalf("KolomCreateTable = %s %v, mau %s %v", nama, kolom, TabelKomisi, KolomViewKomisiLama)
	}
	ddl := satuBaris(maju[0])
	for _, mau := range []string{"ID VARCHAR2(10) NOT NULL", "IDUSEDBY VARCHAR2(10),", "USEDBY VARCHAR2(200),",
		"CONTRACT NUMBER(5),", "YEAR NUMBER(5),", "COMM NUMBER(38,8),", "CONSTRAINT PK_RICOMM_LIFE PRIMARY KEY (ID)"} {
		if !strings.Contains(ddl, mau) {
			t.Errorf("DDL tanpa %q:\n%s", mau, ddl)
		}
	}
	// Selain PK, semua kolom NULLABLE (wajib-isi di Go).
	if strings.Count(ddl, "NOT NULL") != 1 {
		t.Errorf("NOT NULL selain PK:\n%s", ddl)
	}
	for _, p := range maju {
		if strings.Contains(strings.ToUpper(p), "VIEW") || strings.Contains(p, "M_RICOMM_LIFE") {
			t.Errorf("jalur maju menyentuh view atau M_RICOMM_LIFE (DROP VIEW = berkas DBA terpisah): %s", p)
		}
	}

	mundur := langkah924(t, true)
	if len(mundur) != 2 || satuBaris(mundur[0]) != "DROP TABLE {skema}.RICOMM_LIFE CASCADE CONSTRAINTS" {
		t.Fatalf("mundur %q", mundur)
	}
	mauView := "CREATE VIEW {skema}.RICOMM_LIFE AS SELECT a.ID, a.JSONDATA.IDUSEDBY, a.JSONDATA.USEDBY, " +
		"a.JSONDATA.CONTRACT, a.JSONDATA.YEAR, a.JSONDATA.COMM FROM {skema}.M_RICOMM_LIFE a"
	if satuBaris(mundur[1]) != mauView {
		t.Errorf("view dipulihkan\n%s\nmau\n%s", satuBaris(mundur[1]), mauView)
	}
}

// Berkas DBA (bukan migrasi) memeriksa ALL_DEPENDENCIES lalu DROP VIEW - dan HANYA view RICOMM_LIFE.
func TestBerkasDBALepasView(t *testing.T) {
	isi, err := os.ReadFile(filepath.Join("..", "..", "docs", "DBA-LEPAS-VIEW-RICOMM_LIFE.sql"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(isi)
	if strings.Contains(s, "\r") {
		t.Error("berkas DBA memuat CR")
	}
	i, j := strings.Index(s, "ALL_DEPENDENCIES"), strings.Index(s, "DROP VIEW POOLDATA.RICOMM_LIFE")
	if i < 0 || j < 0 || i > j {
		t.Errorf("ALL_DEPENDENCIES (%d) harus mendahului DROP VIEW POOLDATA.RICOMM_LIFE (%d)", i, j)
	}
	for _, b := range strings.Split(s, "\n") {
		u := strings.ToUpper(strings.TrimSpace(b))
		if strings.HasPrefix(u, "--") {
			continue
		}
		if strings.HasPrefix(u, "DROP ") && u != "DROP VIEW POOLDATA.RICOMM_LIFE;" {
			t.Errorf("DROP lain di berkas DBA: %s", b)
		}
	}
}
