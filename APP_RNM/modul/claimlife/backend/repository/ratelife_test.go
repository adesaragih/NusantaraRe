package repository

// Pembaca RATE_LIFE - OQ-M7 (GILIRAN-17), izin sempit seperti butir bh.

import (
	"context"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

func TestRateLifeKolomTetapBerkunciIDUSEDBY(t *testing.T) {
	q := sqlRateLife("S.R")
	mau := "SELECT r.ID, r.AGE, r.CONTRACT, r.GENDER, r.RATE FROM S.R r WHERE r.IDUSEDBY = :1 ORDER BY r.ID"
	if strings.Join(strings.Fields(q), " ") != mau {
		t.Errorf("SQL:\n%s\nmau:\n%s", q, mau)
	}
	if err := db.PeriksaSQL(q); err != nil {
		t.Error(err)
	}
	if len(KolomRateLife) != 5 {
		t.Errorf("kolom %v, mau tepat lima", KolomRateLife)
	}
}

// Pengenal kosong gagal terang sebelum basis data - `OUTWARDRATEID` produk
// yang tidak ada bukan izin membaca seluruh view.
func TestRateLifePengenalKosongGagalTerang(t *testing.T) {
	if _, err := NewRateLife(nil).Baca(context.Background(), "  "); err == nil {
		t.Error("IDUSEDBY kosong diterima")
	}
}
