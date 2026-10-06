package repository

// Migrasi inti 926 (tabel flat RATE_LIFE_SUMMARY, K-F1/K-F2) dibaca pengurai PRODUKSI pelari (`migrasi.KolomCreateTable`)
// dan berpasangan dengan jalur mundurnya; lebar kolom yang dipakai alat pindah = DDL.

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

const (
	berkas926     = "926_rate_life_summary_flat.sql"
	berkas926Down = "926_rate_life_summary_flat_down.sql"
)

func langkah926(t *testing.T, mundur bool) []string {
	t.Helper()
	sumber := fstest.MapFS{}
	for _, n := range []string{berkas926, berkas926Down} {
		isi, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "inti", "backend", "migrations", n))
		if err != nil {
			t.Fatal(err)
		}
		sumber["migrations/"+n] = &fstest.MapFile{Data: isi}
	}
	l, err := migrasi.Daftar(mundur, sumber)
	if err != nil || len(l) != 1 || migrasi.KunciLangkah(l[0].Nama) != "926_rate_life_summary_flat" {
		t.Fatalf("membaca 926 (mundur=%v): %v %v", mundur, l, err)
	}
	return l[0].Pernyataan
}

func TestMigrasi926TeruraiDanBerpasangan(t *testing.T) {
	maju := langkah926(t, false)
	if len(maju) != 2 {
		t.Fatalf("926 maju %d pernyataan, mau 2 (CREATE TABLE + indeks pengaman)", len(maju))
	}
	nama, kolom := migrasi.KolomCreateTable(maju[0])
	if nama != TabelRingkasan || !slices.Equal(kolom, KolomRingkasanFlat) {
		t.Fatalf("KolomCreateTable = %s %v, mau %s %v", nama, kolom, TabelRingkasan, KolomRingkasanFlat)
	}
	// RALAT R5: FLAG tidak digunakan - tidak ada di DDL tabel flat (jalur mundur tetap memulihkan view lengkap).
	if slices.Contains(kolom, "FLAG") || strings.Contains(maju[0], "FLAG") {
		t.Errorf("DDL 926 memuat FLAG:\n%s", maju[0])
	}
	if len(LebarKolomFlat) != len(KolomRingkasanFlat) {
		t.Errorf("LebarKolomFlat %v, kolom flat %v", LebarKolomFlat, KolomRingkasanFlat)
	}
	ddl := satuBaris(maju[0])
	if strings.Count(ddl, "NOT NULL") != 1 || !strings.Contains(ddl, "CONSTRAINT PK_RATE_LIFE_SUMMARY PRIMARY KEY (ID)") {
		t.Errorf("PK / NULLABLE:\n%s", ddl)
	}
	// Lebar alat pindah = lebar DDL, dan semua kolom teks (view: semua VARCHAR2).
	for _, k := range KolomRingkasanFlat {
		m := regexp.MustCompile(k + ` VARCHAR2\((\d+)\)`).FindStringSubmatch(ddl)
		if m == nil || m[1] != fmt.Sprint(LebarKolomFlat[k]) {
			t.Errorf("%s: DDL %v, LebarKolomFlat %d", k, m, LebarKolomFlat[k])
		}
	}
	// Pengaman view: pernyataan kedua = indeks fungsi atas nama (dipakai SqlPemakaiNama) - gagal ORA-01702 atas view.
	if satuBaris(maju[1]) != "CREATE INDEX {skema}.IX_RATE_LIFE_SUMMARY_NAMA ON {skema}.RATE_LIFE_SUMMARY (UPPER(TRIM(USEDBY)))" {
		t.Errorf("indeks pengaman %q", maju[1])
	}
	if !strings.Contains(SqlPemakaiNama("V"), "UPPER(TRIM(USEDBY))") {
		t.Error("indeks nama tidak lagi dipakai SqlPemakaiNama")
	}
	for _, p := range maju {
		if strings.Contains(strings.ToUpper(p), "VIEW") || strings.Contains(p, "M_RATE_LIFE_SUMMARY") {
			t.Errorf("jalur maju menyentuh view atau M_RATE_LIFE_SUMMARY (DROP VIEW = berkas DBA): %s", p)
		}
	}

	mundur := langkah926(t, true)
	if len(mundur) != 2 || satuBaris(mundur[0]) != "DROP TABLE {skema}.RATE_LIFE_SUMMARY CASCADE CONSTRAINTS" {
		t.Fatalf("mundur %q", mundur)
	}
	mauView := "CREATE VIEW {skema}.RATE_LIFE_SUMMARY AS SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, " +
		"a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID, a.JSONDATA.FLAG FROM {skema}.M_RATE_LIFE_SUMMARY a"
	if satuBaris(mundur[1]) != mauView {
		t.Errorf("view dipulihkan\n%s\nmau\n%s", satuBaris(mundur[1]), mauView)
	}
}

// Berkas DBA (bukan migrasi): ALL_DEPENDENCIES mendahului DROP VIEW, dan hanya view RATE_LIFE_SUMMARY yang di-DROP.
func TestBerkasDBALepasView(t *testing.T) {
	isi, err := os.ReadFile(filepath.Join("..", "..", "docs", "DBA-LEPAS-VIEW-RATE_LIFE_SUMMARY.sql"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(isi)
	if strings.Contains(s, "\r") {
		t.Error("berkas DBA memuat CR")
	}
	i, j := strings.Index(s, "ALL_DEPENDENCIES"), strings.Index(s, "DROP VIEW POOLDATA.RATE_LIFE_SUMMARY;")
	if i < 0 || j < 0 || i > j {
		t.Errorf("ALL_DEPENDENCIES (%d) harus mendahului DROP VIEW (%d)", i, j)
	}
	for _, b := range strings.Split(s, "\n") {
		u := strings.ToUpper(strings.TrimSpace(b))
		if !strings.HasPrefix(u, "--") && strings.HasPrefix(u, "DROP ") && u != "DROP VIEW POOLDATA.RATE_LIFE_SUMMARY;" {
			t.Errorf("DROP lain: %s", b)
		}
	}
}
