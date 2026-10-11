package repository

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// TestPetaBindSisip - instrumen diuji atas jawaban yang diketahui: sqlSisipSpreading (14 bind; :9-:13 angkaMasuk; ID /
// PARENT_ID / SEQ_NO nomor), T_LOCATIONLIST.PARENT_ID (VARCHAR2 ID case) tidak dipetakan, ekspresi NVL / literal dilewati.
func TestPetaBindSisip(t *testing.T) {
	peta, n, err := petaBindSisip(sqlSisipSpreading("S.SPR"), TabelSpreadingList)
	if err != nil || n != 14 {
		t.Fatalf("%v %d", err, n)
	}
	mau := map[int]kolomBind{0: {nama: "ID"}, 1: {nama: "PARENT_ID"}, 4: {nama: "SEQ_NO"}, 8: {"SHARE_PERCENTAGE", true},
		9: {"TSI_SPREADED", true}, 10: {"TSI_GROSS_SPREADED", true}, 11: {"PREMIUM_SPREADED", true}, 12: {"CLAIM_ESTIMATION", true}}
	if !reflect.DeepEqual(peta, mau) {
		t.Errorf("peta: %v", peta)
	}
	q := "INSERT INTO S.LOC (ID, PARENT_ID, SEQ_NO, X, Y) VALUES (:1, :2, :3, NVL(:4, 'A,B'), 'NULL')"
	if peta, n, err := petaBindSisip(q, TabelLocationList); err != nil || n != 4 || !reflect.DeepEqual(peta, map[int]kolomBind{0: {nama: "ID"}, 2: {nama: "SEQ_NO"}}) {
		t.Errorf("lokasi: %v %d %v", peta, n, err)
	}
	if _, _, err := petaBindSisip("UPDATE S.X SET A = :1", "T"); !errors.Is(err, ErrBindTidakSah) {
		t.Errorf("bukan INSERT: %v", err)
	}
}

// TestPeriksaBindBaris - bentuk nilai yang ditolak (tabel, kolom, baris, BENTUK; nilai mentah tidak dicetak) dan yang lolos.
func TestPeriksaBindBaris(t *testing.T) {
	q := sqlSisipSpreading("S.SPR")
	baris := func(ubah func([]any)) [][]any {
		b := []any{"11", "12", "T_COVERAGELIST", "J", 1, "uid", "QS", nil, "100", "137500000000000", "500000000000000",
			"1375000000000", nil, "IDR"}
		ubah(b)
		return [][]any{{"21", "12", "T_COVERAGELIST", "J", 1, "uid", "OR", nil, "0", "0", nil, "0", nil, "IDR"}, b}
	}
	if err := periksaBindBaris(q, TabelSpreadingList, baris(func([]any) {}), 0); err != nil {
		t.Fatalf("bentuk sah: %v", err)
	}
	for _, k := range []struct {
		i      int
		v      any
		bagian string
	}{
		{9, "137500000000000.123456789", "kolom TSI_SPREADED baris 2: panjang bulat 15, desimal 9"},
		{11, "1.375E+12", "kolom PREMIUM_SPREADED baris 2: memuat E, 2 karakter bukan angka, panjang bulat 1, desimal 7"},
		{8, "", "kolom SHARE_PERCENTAGE baris 2: teks kosong"},
		{8, "100,5", "memuat koma"},
		{8, " 100", "memuat spasi"},
		{9, strings.Repeat("9", 32) + "." + strings.Repeat("1", 20), "panjang bulat 32, desimal 20"},
		{9, 1.5, "kolom TSI_SPREADED baris 2: tipe Go float64"},
		{0, "K1", "kolom ID baris 2: 1 karakter bukan angka, panjang bulat 2"},
		{1, nil, "kolom PARENT_ID baris 2: NULL (kolom NOT NULL)"},
		{4, "1.5", "kolom SEQ_NO baris 2: panjang bulat 1, desimal 1"},
	} {
		err := periksaBindBaris(q, TabelSpreadingList, baris(func(b []any) { b[k.i] = k.v }), 0)
		if !errors.Is(err, ErrBindTidakSah) || !strings.Contains(err.Error(), k.bagian) || !strings.Contains(err.Error(), TabelSpreadingList) {
			t.Errorf("%v: %v", k.v, err)
		}
		if s, ok := k.v.(string); ok && len(s) > 3 && strings.Contains(err.Error(), s) {
			t.Errorf("nilai mentah tercetak: %v", err)
		}
	}
	if err := periksaBindBaris(q, TabelSpreadingList, [][]any{{"1"}}, 4); !errors.Is(err, ErrBindTidakSah) || !strings.Contains(err.Error(), "baris 5: 1 nilai, teks SQL 14 bind") {
		t.Errorf("jumlah nilai: %v", err)
	}
}

// TestSisipAnekaPerBaris - bug DEV 10-10-2026: coverage / deductible / spreading Aneka satu eksekusi per baris (bind skalar
// seperti spreading FIRE), setelah induknya; tabel lain tetap array-bind; galat Oracle dibungkus tabel + baris + jumlah bind.
func TestSisipAnekaPerBaris(t *testing.T) {
	ta := tabelAnekaUji()
	share := desimalUji("100")
	cov := func() []models.CoverageAneka {
		return []models.CoverageAneka{{Coverage: "PL", Deductibles: []models.Deductible{{TypeDeductible: "1"}},
			Spreading: []models.BarisSpreading{{TreatyType: "QS", TreatyName: "2024 QS 155M FAC", SharePercentage: share}}}}
	}
	lokasi := []models.LokasiAneka{
		{ObjectNo: "1", Okupasi: []models.OkupasiAneka{{Objek: []models.ObjekAneka{{ObjectName: "A", Currency: "IDR", Coverages: cov()}}}}},
		{ObjectNo: "2", Okupasi: []models.OkupasiAneka{{Objek: []models.ObjekAneka{{ObjectName: "B", Currency: "IDR", Coverages: cov()}}}}},
	}
	e := &eksekutorTiruan{u: &urutTiruan{n: map[string]int{}}, seragam: true}
	if _, _, err := sisipAnekaBerkelompok(context.Background(), e, ta, "K1", lokasi, uidUrut()); err != nil {
		t.Fatal(err)
	}
	eks := map[string]int{}
	pertama := map[string]int{}
	for i, q := range e.urutan {
		tabel := strings.Fields(q)[2]
		if _, ada := pertama[tabel]; !ada {
			pertama[tabel] = i
		}
		eks[tabel]++
	}
	if eks["S.SPR"] != 2 || eks["S.COV"] != 2 || eks["S.DED"] != 2 || eks["S.LOC"] != 1 || eks["S.ANEKA"] != 1 {
		t.Errorf("eksekusi per tabel: %v", eks)
	}
	if !(pertama["S.ANEKA"] < pertama["S.COV"] && pertama["S.COV"] < pertama["S.DED"] && pertama["S.COV"] < pertama["S.SPR"]) {
		t.Errorf("urutan induk-anak: %v", pertama)
	}
	gagal := &eksekutorTiruan{u: &urutTiruan{n: map[string]int{}}, gagal: TabelSpreadingList}
	ta.spread = TabelSpreadingList
	_, _, err := sisipAnekaBerkelompok(context.Background(), gagal, ta, "K1", lokasi, uidUrut())
	if err == nil || !strings.Contains(err.Error(), "ORA-01722") || !strings.Contains(err.Error(), "(T_SPREADINGLIST baris 1-1 dari 2, 14 bind, per baris)") {
		t.Errorf("pembungkus: %v", err)
	}
}
