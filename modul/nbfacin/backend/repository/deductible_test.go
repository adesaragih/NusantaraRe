package repository

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// TestSQLDeductible - tiket 45: deductible berinduk coverage dibaca lewat case urut coverage lalu SEQ_NO; dihapus SEBELUM
// coverage; PARENT_TABLE / SRC_PATH = isi loader.
func TestSQLDeductible(t *testing.T) {
	baca := sqlBacaDeductible(tabelObjekUji())
	for _, harus := range []string{"FROM UJI.E e", "JOIN UJI.V v ON v.ID = e.PARENT_ID", "WHERE v.PARENT_TABLE = 'T_PROPERTYITEMLIST' AND l.PARENT_ID = :1",
		"ORDER BY e.PARENT_ID, e.SEQ_NO"} {
		if !strings.Contains(baca, harus) {
			t.Errorf("baca deductible tanpa %q", harus)
		}
	}
	posisi := map[string]int{}
	for i, q := range sqlHapusObjek(tabelObjekUji()) {
		posisi[strings.Fields(q)[2]] = i
	}
	if posisi["UJI.E"] >= posisi["UJI.V"] {
		t.Errorf("deductible harus dihapus sebelum coverage: %v", posisi)
	}
	if indukDeductible != "T_COVERAGELIST" || jalurDeductible != "LocationList/Property/PropertyItemList/CoverageList/DeductibleList" {
		t.Error("PARENT_TABLE / SRC_PATH harus sama dengan loader")
	}
}

// TestSisipDeductibleMenurutKolom - kolom NUMBER(38,8) lewat angkaMasuk, kode lewat TO_NUMBER, VARCHAR2 (termasuk
// TIME_EXCESS) apa adanya; bind dipasangkan lewat KUNCI; CURRENCY hanya bila terisi (A169).
func TestSisipDeductibleMenurutKolom(t *testing.T) {
	for _, denganMataUang := range []bool{true, false} {
		m := regexp.MustCompile(`\(([^)]*)\) VALUES \((.*)\)$`).FindStringSubmatch(sqlSisipDeductible("UJI.E", denganMataUang))
		kolom := strings.Split(m[1], ", ")
		nilai := regexp.MustCompile(`TO_NUMBER\(:\d+[^)]*\)|:\d+`).FindAllString(m[2], -1)
		mau := 6 + len(kolomDeductibleData)
		if denganMataUang {
			mau++
		}
		if len(kolom) != mau || len(nilai) != mau || len(kolomDeductibleData) != 9 {
			t.Fatalf("mata uang %v: %d kolom / %d nilai", denganMataUang, len(kolom), len(nilai))
		}
		if (kolom[len(kolom)-1] == "CURRENCY") != denganMataUang {
			t.Errorf("mata uang %v: kolom terakhir %s", denganMataUang, kolom[len(kolom)-1])
		}
		angka := map[string]bool{"AMOUNT": true, "PCT_DEDUCTIBLE": true, "PCT_DEDUCTIBLE2": true}
		kode := map[string]bool{"TYPE_DEDUCTIBLE": true, "TYPE_DEDUCTIBLE2": true}
		for i, k := range kolom {
			if !strings.Contains(nilai[i], fmt.Sprintf(":%d", i+1)) {
				t.Errorf("kolom %s nilai %q", k, nilai[i])
			}
			nls := strings.Contains(nilai[i], "NLS_NUMERIC_CHARACTERS")
			if nls != angka[k] || (strings.HasPrefix(nilai[i], "TO_NUMBER(") && !nls) != kode[k] {
				t.Errorf("kolom %s nilai %q", k, nilai[i])
			}
		}
	}
	d := models.Deductible{TypeDeductible: "7", MinMax: "3", TypeDeductible2: "1", Condition: "5", InputCondition: "UJI",
		PctDeductible: desimalUji("10"), PctDeductible2: desimalUji("2.5"), Amount: desimalUji("1000000"), TimeExcess: desimalUji("14")}
	mau := map[string]any{"AMOUNT": "1000000", "CONDITION": "5", "INPUT_CONDITION": "UJI", "MIN_MAX": "3", "PCT_DEDUCTIBLE": "10",
		"PCT_DEDUCTIBLE2": "2.5", "TIME_EXCESS": "14", "TYPE_DEDUCTIBLE": "7", "TYPE_DEDUCTIBLE2": "1"}
	arg := argDeductible(d)
	if len(arg) != len(kolomDeductibleData) {
		t.Fatalf("tanpa mata uang: %d bind", len(arg))
	}
	for i, k := range kolomDeductibleData {
		if arg[i] != mau[k.nama] {
			t.Errorf("bind %s = %v, mau %v", k.nama, arg[i], mau[k.nama])
		}
	}
	d.Currency = "USD"
	if arg = argDeductible(d); len(arg) != len(kolomDeductibleData)+1 || arg[len(arg)-1] != "USD" {
		t.Errorf("dengan mata uang: %v", arg)
	}
	if arg = argDeductible(models.Deductible{}); arg[1] != nil || arg[0] != nil {
		t.Errorf("kosong mestinya NULL: %v", arg)
	}
	baca := kolomBacaDeductible()
	if baca[0] != "TO_CHAR(e.PARENT_ID)" || baca[1] != "e.CURRENCY" || len(baca) != 11 || baca[9] != "TO_CHAR(e.TYPE_DEDUCTIBLE)" || baca[8] != "e.TIME_EXCESS" {
		t.Errorf("kolom baca %v", baca)
	}
	if mataUangBaca("UNKNOWN") != "" || mataUangBaca("USD") != "USD" {
		t.Error("UNKNOWN mestinya dibaca kosong")
	}
}
