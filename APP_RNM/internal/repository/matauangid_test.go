package repository

// Uji bentuk query pengenal mata uang - celah sensus 28-09-2026.

import (
	"strings"
	"testing"
)

func TestQueryPengenalMataUangMeniruRuleNya(t *testing.T) {
	q := sqlPengenalMataUang("SKEMAUJI.CURRENCY")
	// VERBATIM `GetCurrencyID.xml` b85:
	//   SELECT ID AS CARI1 FROM POOLDATA.CURRENCY WHERE CURRENCY = {...}
	if !strings.Contains(q, "SELECT ID FROM") {
		t.Errorf("kolom yang dibaca bukan ID:\n%s", q)
	}
	if !strings.Contains(q, "WHERE CURRENCY = :1") {
		t.Errorf("syaratnya bukan kode mata uang, atau tidak di-bind:\n%s", q)
	}
	// ⛔ Nilainya lewat BIND, tidak pernah ditempel.
	if strings.Contains(q, "'") {
		t.Errorf("query memuat literal berkutip:\n%s", q)
	}
	if err := PeriksaSQL(q); err != nil {
		t.Errorf("%v\n%s", err, q)
	}
	// ⚠️ Batas satu baris DITAMBAHKAN; rule aslinya membaca `pxResults(1)`
	// saja. Menuliskannya membuat niat itu terbaca dan menahan tabel yang
	// ternyata berisi kode kembar.
	if !strings.Contains(q, "FETCH FIRST 1 ROWS ONLY") {
		t.Errorf("query tidak berbatas satu baris:\n%s", q)
	}
}

func TestPengenalMataUangKosongBukanGalat(t *testing.T) {
	// ⛔ Kode KOSONG bukan kode yang SALAH (ADR-U-0027). Ia dijawab kosong
	// tanpa menyentuh basis data sama sekali - dan uji ini membuktikannya
	// dengan pembaca yang db-nya nil: bila ia menyentuh db, ia akan panik.
	id, err := (&MataUang{}).Pengenal(nil, "   ")
	if err != nil || id != "" {
		t.Errorf("kode kosong -> (%q, %v), mau (\"\", nil)", id, err)
	}
}
