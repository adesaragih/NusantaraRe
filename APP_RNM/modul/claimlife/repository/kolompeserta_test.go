package repository

// CLAIM_AMOUNT peserta - GILIRAN-14 butir bp. TANPA Oracle.

import (
	"database/sql"
	"strings"
	"testing"

	"nusantarare/inti/backend/uang"
	"nusantarare/modul/claimlife/models"
)

// TestPesertaMenyimpanJumlahKlaim - 7.7 b3280 menyalin `.CLAIM_AMOUNT` ke
// peserta; kolom 003-nya sudah ada, tetapi tidak pernah ditulis.
func TestPesertaMenyimpanJumlahKlaim(t *testing.T) {
	jumlah, err := uang.NewMoney("25000.1234", "IDR")
	if err != nil {
		t.Fatal(err)
	}
	q, nilai := insertPeserta("SKEMAUJI.T_CLAIMLF_PREMIUMLIST_DETAIL", "P-1", "CLM-1",
		models.Peserta{MataUang: "IDR", JumlahKlaim: jumlah})
	kolom := strings.Split(strings.Split(strings.SplitN(q, "(", 2)[1], ")")[0], ", ")
	for i, k := range kolom {
		if strings.TrimSpace(k) == "CLAIM_AMOUNT" {
			if nilai[i] != "25000.1234" {
				t.Errorf("CLAIM_AMOUNT ditulis %v, mau 25000.1234", nilai[i])
			}
			return
		}
	}
	t.Fatalf("insertPeserta tidak menulis CLAIM_AMOUNT:\n%s", q)
}

// TestPesertaMembacaKembaliJumlahKlaim - tulis dan baca bertemu di satu daftar.
func TestPesertaMembacaKembaliJumlahKlaim(t *testing.T) {
	sel := make([]sql.NullString, len(kolomPeserta))
	for i, k := range kolomPeserta {
		switch k.Nama {
		case "CURRENCY":
			sel[i] = sql.NullString{String: "IDR", Valid: true}
		case "CLAIM_AMOUNT":
			sel[i] = sql.NullString{String: "25000.1234", Valid: true}
		}
	}
	p, err := rakitPeserta("P-1", sel)
	if err != nil {
		t.Fatal(err)
	}
	if p.JumlahKlaim.Kosong() || p.JumlahKlaim.Amount.Text('f') != "25000.1234" ||
		p.JumlahKlaim.Currency != "IDR" {
		t.Errorf("JumlahKlaim terbaca %+v, mau 25000.1234 IDR", p.JumlahKlaim)
	}
}
