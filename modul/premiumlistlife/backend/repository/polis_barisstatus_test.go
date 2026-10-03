package repository

// Uji satu baris T_PREMIUM_LIST per status penawaran (migrasi 062). TANPA Oracle.

import (
	"context"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

func TestPengenalBarisStatusUnikDanBukanNomorKasus(t *testing.T) {
	a, p := PengenalBarisStatus("NBLF-1", "Accept"), PengenalBarisStatus("NBLF-1", "Pending")
	if len(a) != 32 || a == p || a == PengenalBarisStatus("NBLF-2", "Accept") || a != PengenalBarisStatus("NBLF-1", "Accept") {
		t.Errorf("pengenal %q / %q - harus 32 heksa, unik per (kasus, status), deterministik", a, p)
	}
	if a == "NBLF-1" {
		t.Error("baris status menimpa baris utama")
	}
}

func TestSalinanBarisStatusHanyaIsianPenawaran(t *testing.T) {
	q := sqlSalinKeBarisStatus("SKEMAUJI.T_PREMIUM_LIST")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"INSERT INTO SKEMAUJI.T_PREMIUM_LIST (ID, ID_PEGA, TGL_INPUT, CREATE_OP_NAME",
		"STATUS_PENAWARAN", "CEDING_CO_NAME", "SELECT :1, ID_PEGA", "WHERE ID = :2"} {
		if !strings.Contains(q, k) {
			t.Errorf("bentuk %q tidak ada:\n%s", k, q)
		}
	}
	// Kolom milik tahap lain TIDAK disalin: baris status tidak boleh ditemukan
	// pencarian polis (Claim Life mencari lewat NO_POLIS).
	for _, k := range []string{"NO_POLIS", " TYPE,", "PRODUCT_NAME", "PROD_KE", "SOB"} {
		if strings.Contains(q, k) {
			t.Errorf("kolom %s ikut disalin:\n%s", k, q)
		}
	}
	if !strings.Contains(sqlKunciStatusPenawaran("S.T"), "FOR UPDATE") {
		t.Error("status lama dibaca tanpa kunci - dua simpanan serentak dapat menggandakan baris")
	}
	if !strings.Contains(sqlBuangBarisStatus("S.T"), "ID = :1 AND ID_PEGA = :2") {
		t.Error("pembuangan baris status tidak dibatasi kasusnya")
	}
}

// Status sama, atau belum pernah ada: tidak ada yang dipindah (dan tanpa Oracle).
func TestPindahBarisStatusDilewatiBilaStatusSama(t *testing.T) {
	r := NewPenawaran(nil)
	for _, c := range [][2]string{{"", "Accept"}, {"Accept", "Accept"}} {
		pindah, err := r.PindahkanBarisStatus(context.Background(), nil, "NBLF-1", c[0], c[1])
		if err != nil || pindah {
			t.Errorf("lama %q baru %q: pindah=%v err=%v", c[0], c[1], pindah, err)
		}
	}
}

func TestKotakMasukHanyaBarisUtama(t *testing.T) {
	q := sqlInboxPolis("S.W", "S.P", "S.D")
	if !strings.Contains(q, "p.ID_PEGA = w.ID AND p.ID = p.ID_PEGA") {
		t.Errorf("kotak masuk menggabung baris status - kasus tampil berulang:\n%s", q)
	}
}
