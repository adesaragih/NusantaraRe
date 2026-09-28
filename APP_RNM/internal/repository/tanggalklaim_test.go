package repository

import (
	"regexp"
	"strings"
	"testing"
	"time"
)

// TestSQLTanggalKlaimBerurutPosisi - driver Oracle mengikat penampung
// BERPOSISI menurut urutan kemunculannya. Penampung yang tertukar menulis
// tanggal konfirmasi ke kolom tanggal terima, tanpa satu pun galat.
func TestSQLTanggalKlaimBerurutPosisi(t *testing.T) {
	q := sqlTanggalKlaim("T")
	if err := PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	urut := regexp.MustCompile(`:(\d)`).FindAllStringSubmatch(q, -1)
	var got []string
	for _, m := range urut {
		got = append(got, m[1])
	}
	if strings.Join(got, "") != "12345" {
		t.Fatalf("urutan penampung %v, mau 1..5", got)
	}
	for i, kolom := range []string{"CLAIM_RECEIVED_DATE", "COMPLETE_DATE", "CONFIRMATION_DATE"} {
		pola := regexp.MustCompile(kolom + `\s*=\s*TO_DATE\(:` + string(rune('1'+i)))
		if !pola.MatchString(q) {
			t.Errorf("%s tidak terikat ke :%d", kolom, i+1)
		}
	}
	if !strings.Contains(q, "WHERE ID = :4 AND CLAIM_ID = :5") {
		t.Error("WHERE harus mengikat peserta DAN klaimnya")
	}
}

// TestKolomTanggalKlaimAdaDiMigrasi003 - kolom yang ditulis harus ada.
func TestKolomTanggalKlaimAdaDiMigrasi003(t *testing.T) {
	isi, err := berkasMigrasi.ReadFile("migrations/003_t_claimlf_premiumlist_detail.sql")
	if err != nil {
		t.Fatal(err)
	}
	for _, kolom := range []string{"CLAIM_RECEIVED_DATE", "COMPLETE_DATE", "CONFIRMATION_DATE"} {
		if !regexp.MustCompile(`(?m)^\s*` + kolom + `\s+DATE`).Match(isi) {
			t.Errorf("%s bukan kolom DATE di migrasi 003", kolom)
		}
	}
}

// TestArgTanggalKosongMenjadiNULL - kosong bukan tanggal nol (ADR-U-0027).
func TestArgTanggalKosongMenjadiNULL(t *testing.T) {
	if argTanggal(nil) != nil {
		t.Error("nil harus menjadi NULL")
	}
	w := time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC)
	if got := argTanggal(&w); got != "2026-03-02 00:00:00" {
		t.Errorf("argTanggal = %v", got)
	}
}
