package repository

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// TestSQLCoverage - tiket 43: coverage berinduk item (PARENT_TABLE literal), dibaca lewat case urut item lalu SEQ_NO;
// dihapus SEBELUM item; PARENT_TABLE / SRC_PATH = isi loader.
func TestSQLCoverage(t *testing.T) {
	baca := sqlBacaCoverage(tabelObjekUji())
	for _, harus := range []string{"FROM UJI.V v", "JOIN UJI.I i ON i.ID = v.PARENT_ID", "WHERE v.PARENT_TABLE = 'T_PROPERTYITEMLIST' AND l.PARENT_ID = :1",
		"ORDER BY v.PARENT_ID, v.SEQ_NO"} {
		if !strings.Contains(baca, harus) {
			t.Errorf("baca coverage tanpa %q", harus)
		}
	}
	posisi := map[string]int{}
	for i, q := range sqlHapusObjek(tabelObjekUji()) {
		posisi[strings.Fields(q)[2]] = i
	}
	if posisi["UJI.V"] >= posisi["UJI.I"] {
		t.Errorf("coverage harus dihapus sebelum item: %v", posisi)
	}
	if indukCoverage != "T_PROPERTYITEMLIST" || jalurCoverage != "LocationList/Property/PropertyItemList/CoverageList" {
		t.Error("PARENT_TABLE / SRC_PATH harus sama dengan loader")
	}
}

// TestSisipCoverageMenurutKolom - SQL, bind, dan kolom baca dari SATU tabel kolom: tiap kolom NUMBER lewat TO_NUMBER,
// kolom VARCHAR2 (termasuk FIRST_LOSS / FIRST_SCALE / LOST_LIMIT / EML_PML) tanpa; bind dipasangkan lewat KUNCI.
func TestSisipCoverageMenurutKolom(t *testing.T) {
	m := regexp.MustCompile(`\(([^)]*)\) VALUES \((.*)\)$`).FindStringSubmatch(sqlSisipCoverage("UJI.V"))
	kolom := strings.Split(m[1], ", ")
	nilai := regexp.MustCompile(`TO_NUMBER\(:\d+[^)]*\)|:\d+`).FindAllString(m[2], -1)
	if len(kolom) != 6+len(kolomCoverageData)+1 || len(nilai) != len(kolom) || len(kolomCoverageData) != 28 {
		t.Fatalf("%d kolom / %d nilai / %d data", len(kolom), len(nilai), len(kolomCoverageData))
	}
	angka := map[string]bool{"TSI": true, "RATE": true, "RATE_OJK": true, "DISCOUNT_PERCENTAGE": true, "TSI_LIABILITY": true, "NET_RATE": true,
		"LIMITOF_LIABILITY": true, "PCT_LO_L": true, "PRO_RATE_PERCENT": true, "INDEMNITY_PERCENTAGE": true, "SUBLIMIT": true, "DISCOUNT": true,
		"PREMIUM": true, "PCT_ADJUSTMENT": true}
	for i, k := range kolom {
		if !strings.Contains(nilai[i], fmt.Sprintf(":%d", i+1)) || strings.HasPrefix(nilai[i], "TO_NUMBER(") != angka[k] {
			t.Errorf("kolom %s nilai %q", k, nilai[i])
		}
	}
	// Nilai uji = nomor urut kolom data: teks "K<n>", desimal n - dipasangkan lewat kunci tabel kolom.
	var c models.CoverageObjek
	mau := map[string]any{}
	for n, k := range kolomCoverageData {
		if k.jenis == jenisTeks {
			*k.teks(&c) = fmt.Sprintf("K%d", n)
			mau[k.nama] = fmt.Sprintf("K%d", n)
		} else {
			*k.des(&c) = desimalUji(fmt.Sprint(n + 1))
			mau[k.nama] = fmt.Sprint(n + 1)
		}
	}
	arg := argCoverage(c, "IDR")
	for i, k := range kolom[6:] {
		if k == "CURRENCY_CODE" {
			if arg[i] != "IDR" {
				t.Errorf("CURRENCY_CODE = %v", arg[i])
			}
			continue
		}
		if arg[i] != mau[k] {
			t.Errorf("bind kolom %s = %v, mau %v", k, arg[i], mau[k])
		}
	}
	// Kolom baca: NUMBER lewat angkaKeluar, VARCHAR2 apa adanya; kunci induk lalu kunci coverage (tiket 45) di depan.
	baca := kolomBacaCoverage()
	if baca[0] != "TO_CHAR(v.PARENT_ID)" || baca[1] != "TO_CHAR(v.ID)" || len(baca) != 30 {
		t.Fatalf("kolom baca %v", baca)
	}
	for i, k := range kolomCoverageData {
		if angka[k.nama] != strings.HasPrefix(baca[i+2], "TO_CHAR(v.") {
			t.Errorf("kolom baca %s = %s", k.nama, baca[i+2])
		}
	}
}
