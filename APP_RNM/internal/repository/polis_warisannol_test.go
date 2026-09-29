package repository

// OQ-PL-10 DITUTUP (GILIRAN-17): kolom uang kosong di M_LIFE_PREMIUM_DETAIL
// warisan = 0, seperti Pega - SATU fungsi di tepi repository warisan.
// ADR-U-0027 (kosong bukan nol) tetap berlaku untuk tabel `T_*`.

import (
	"database/sql"
	"html"
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
)

func TestKosongWarisanJadiNolSepertiToDecimal(t *testing.T) {
	b := BarisWarisan{Tipe: "QP", Nilai: map[string]sql.NullString{
		"GROSS_PREMIUM": {String: "  ", Valid: true},
		"SEX":           {},
	}}
	arg := nilaiSalinWarisan("UJI-PL-1", "UJI-POLIS-1", b)
	for i, k := range kolomPesertaWarisan {
		switch {
		case kolomNolBilaKosongWarisan[k.Kolom]:
			if arg[i] != "0" {
				t.Errorf("%s kosong = %v, mau \"0\" (@toDecimal)", k.Kolom, arg[i])
			}
		case k.Kolom == "SEX" || k.Kolom == "DOB" || k.Kolom == "PERIOD_YY" || k.Kolom == "AGE":
			if arg[i] != nil {
				t.Errorf("%s kosong = %v, mau NULL (bukan @toDecimal)", k.Kolom, arg[i])
			}
		}
	}
	b.Nilai["GROSS_PREMIUM"] = sql.NullString{String: "12.5", Valid: true}
	for i, k := range kolomPesertaWarisan {
		if k.Kolom == "GROSS_PREMIUM" && nilaiSalinWarisan("P", "I", b)[i] != "12.5" {
			t.Error("nilai terisi diubah")
		}
	}
}

// Himpunannya DITURUNKAN dari korpus, dua arah: kolom VALUES `SaveMasterLPDet`
// yang memakai `TempInputDetail.CARIn`, dengan `CARIn` ditetapkan
// `@toDecimal(...)` di `InsertLifePremiumDetail_act` langkah 3.3.3 (b1843,
// hidup). `@toDecimal("")` = 0.
func TestKolomNolBilaKosongWarisanDariKorpus(t *testing.T) {
	const akar = `D:\XML\RNM_BRD\PremiumList Life`
	rdb, err := os.ReadFile(akar + `\RDBList\SaveMasterLPDet.xml`)
	if err != nil {
		t.Skipf("korpus tidak terjangkau (%v)", err)
	}
	act, err := os.ReadFile(akar + `\Activity\InsertLifePremiumDetail_act.xml`)
	if err != nil {
		t.Skipf("korpus tidak terjangkau (%v)", err)
	}
	toDec := map[string]bool{}
	for _, m := range regexp.MustCompile(`<PropertiesName>TempInputDetail\.(CARI\d+)</PropertiesName>\s*<PropertiesValue>@toDecimal\(`).
		FindAllStringSubmatch(html.UnescapeString(string(act)), -1) {
		toDec[m[1]] = true
	}
	m := regexp.MustCompile(`(?s)INSERT INTO POOLDATA\.M_LIFE_PREMIUM_DETAIL\s*\((.*?)\)\s*VALUES\s*\((.*)\)`).
		FindStringSubmatch(html.UnescapeString(string(rdb)))
	if m == nil || len(toDec) == 0 {
		t.Fatal("korpus tidak terbaca; pembacanya yang rusak")
	}
	kolom := strings.Split(m[1], ",")
	var nilai []string
	tingkat, cur := 0, ""
	for _, c := range m[2] {
		if c == ')' && tingkat == 0 {
			break
		}
		switch c {
		case '(':
			tingkat++
		case ')':
			tingkat--
		}
		if c == ',' && tingkat == 0 {
			nilai, cur = append(nilai, cur), ""
			continue
		}
		cur += string(c)
	}
	nilai = append(nilai, cur)
	if len(kolom) != len(nilai) {
		t.Fatalf("kolom %d lawan nilai %d", len(kolom), len(nilai))
	}
	var dariKorpus []string
	cari := regexp.MustCompile(`TempInputDetail\.(CARI\d+)`)
	for i, k := range kolom {
		if c := cari.FindStringSubmatch(nilai[i]); c != nil && toDec[c[1]] {
			dariKorpus = append(dariKorpus, strings.TrimSpace(k))
		}
	}
	var diKode []string
	for k := range kolomNolBilaKosongWarisan {
		diKode = append(diKode, k)
	}
	sort.Strings(dariKorpus)
	sort.Strings(diKode)
	if strings.Join(dariKorpus, ",") != strings.Join(diKode, ",") {
		t.Errorf("himpunan nol-bila-kosong berbeda dari korpus:\nkorpus %v\nkode   %v", dariKorpus, diKode)
	}
}
