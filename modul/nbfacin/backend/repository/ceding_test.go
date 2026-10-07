package repository

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// TestSQLCeding - tiket 34: daftar dibaca urut SEQ_NO lewat T_QUOTATIONDATA case; ganti utuh
// (hapus lalu sisip); gabungan ditulis ke baris quotation case. Bind bernomor.
func TestSQLCeding(t *testing.T) {
	for got, mau := range map[string]string{
		sqlBacaCeding("UJI.C", "UJI.Q"): "SELECT c.CEDING_CO, c.CEDING_CO_NAME FROM UJI.C c JOIN UJI.Q q ON q.ID = c.PARENT_ID WHERE q.PARENT_ID = :1 ORDER BY c.SEQ_NO",
		sqlIDQuotation("UJI.Q"):         "SELECT ID FROM UJI.Q WHERE PARENT_ID = :1",
		sqlHapusCeding("UJI.C"):         "DELETE FROM UJI.C WHERE PARENT_ID = :1",
		sqlSisipCeding("UJI.C"):         "INSERT INTO UJI.C (ID, PARENT_ID, SEQ_NO, ROW_UID, CEDING_CO, CEDING_CO_NAME) VALUES (:1, :2, :3, :4, :5, :6)",
		sqlGabunganCeding("UJI.Q"):      "UPDATE UJI.Q SET CEDING_CO = :1, CEDING_CO_NAME = :2 WHERE ID = :3",
	} {
		if got != mau {
			t.Errorf("SQL\n%q\nmau\n%q", got, mau)
		}
	}
	sql := bacaMigrasi(t, "185_t_cedingcolist.sql")
	for _, harus := range []string{"CREATE TABLE {skema}.T_CEDINGCOLIST (", "CREATE SEQUENCE {skema}." + SequenceCedingCoList,
		"ADD (\n  CEDING_CO VARCHAR2(1000)\n)", "MODIFY (CEDING_CO_NAME VARCHAR2(4000))"} {
		if !strings.Contains(sql, harus) {
			t.Errorf("185 tanpa %q", harus)
		}
	}
}

// TestGabungCeding - uji instrumen dengan jawaban yang diketahui (E-8, SetCedingCo_Act):
// urut pilih, `;` tanpa spasi, ganda dipertahankan (E-7), kosong = "".
func TestGabungCeding(t *testing.T) {
	for _, u := range []struct {
		daftar     []models.Ceding
		kode, nama string
	}{
		{nil, "", ""},
		{[]models.Ceding{{ID: "UJI-B", Name: "UJI BETA"}}, "UJI-B", "UJI BETA"},
		{[]models.Ceding{{ID: "UJI-B", Name: "UJI BETA"}, {ID: "UJI-A", Name: "UJI ALFA"}, {ID: "UJI-B", Name: "UJI BETA"}},
			"UJI-B;UJI-A;UJI-B", "UJI BETA;UJI ALFA;UJI BETA"},
		// Semantik @If Pega: nama kosong di depan tanpa ";" awal, di tengah tetap ";;".
		{[]models.Ceding{{ID: "UJI-A", Name: ""}, {ID: "UJI-B", Name: "UJI BETA"}}, "UJI-A;UJI-B", "UJI BETA"},
		{[]models.Ceding{{ID: "UJI-A", Name: "UJI ALFA"}, {ID: "UJI-B", Name: ""}, {ID: "UJI-C", Name: "UJI CE"}},
			"UJI-A;UJI-B;UJI-C", "UJI ALFA;;UJI CE"},
	} {
		if k, n := gabungCeding(u.daftar); k != u.kode || n != u.nama {
			t.Errorf("%v -> %q %q, mau %q %q", u.daftar, k, n, u.kode, u.nama)
		}
	}
}

// TestPeriksaLebarCeding - nama per baris <= 500, gabungan kode <= 1000, nama <= 4000 (butir 80).
func TestPeriksaLebarCeding(t *testing.T) {
	if err := periksaLebarCeding([]models.Ceding{{ID: "UJI-1", Name: strings.Repeat("U", 500)}}); err != nil {
		t.Errorf("tepat 500: %v", err)
	}
	if err := periksaLebarCeding([]models.Ceding{{ID: "UJI-1", Name: strings.Repeat("U", 501)}}); !errors.Is(err, ErrCedingTerlaluPanjang) {
		t.Errorf("nama 501: %v", err)
	}
	banyak := make([]models.Ceding, 0, 120)
	for i := 0; i < 120; i++ { // 120 x 8 + 119 pemisah = 1079 > 1000
		banyak = append(banyak, models.Ceding{ID: "UJI-0000", Name: "N"})
	}
	if err := periksaLebarCeding(banyak); !errors.Is(err, ErrCedingTerlaluPanjang) {
		t.Errorf("gabungan kode 1079: %v", err)
	}
	if err := periksaLebarCeding(banyak[:111]); err != nil { // 111 x 8 + 110 = 998
		t.Errorf("gabungan kode 998: %v", err)
	}
}

// TestUIDAcak - ROW_UID baris aplikasi: UUID v4 36 karakter, berbeda tiap panggilan (A106).
func TestUIDAcak(t *testing.T) {
	a, err1 := uidAcak()
	b, err2 := uidAcak()
	pola := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	if err1 != nil || err2 != nil || !pola.MatchString(a) || a == b {
		t.Errorf("%q %q (%v %v)", a, b, err1, err2)
	}
}
