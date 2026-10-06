package repository

import (
	"regexp"
	"slices"
	"strings"
	"testing"
)

// kolomView - nama kolom view (`a.ID` dan `a.JSONDATA.<KUNCI>`) menurut teks definisinya.
func kolomView(definisi string) []string {
	var k []string
	for _, m := range regexp.MustCompile(`a\.(?:JSONDATA\.)?(\w+)`).FindAllStringSubmatch(definisi, -1) {
		k = append(k, m[1])
	}
	return k
}

// RALAT R1: definisi view RATE_LIFE_SUMMARY (ALL_VIEWS, dibaca WO 06-10-2026) = enam kolom ini.
func TestDefinisiViewRingkasanEnamKolom(t *testing.T) {
	mau := []string{"ID", "USEDBY", "TYPE", "MODIFIEDDATE", "OPERATORID", "FLAG"}
	if k := kolomView(DefinisiViewRingkasan); !slices.Equal(k, mau) {
		t.Fatalf("kolom view %v, mau %v", k, mau)
	}
	if !strings.Contains(DefinisiViewRingkasan, "FROM "+TabelRingkasan+" a") {
		t.Errorf("view bukan atas %s", TabelRingkasan)
	}
	if !strings.Contains(DefinisiViewRate, "FROM "+TabelRate+" a") {
		t.Errorf("view rate bukan atas %s", TabelRate)
	}
}

// Setiap kunci JSON yang ditulis dan setiap kolom yang dibaca ADA di definisi view - nol nama karangan.
// TYPE dan FLAG tidak dirujuk XML InboxSummaryRIRate: tidak ditulis, tidak dibaca.
func TestKunciDanKolomCocokDenganView(t *testing.T) {
	ring, rate := kolomView(DefinisiViewRingkasan), kolomView(DefinisiViewRate)
	for _, k := range []string{JSONUsedBy, JSONOperatorID, JSONModified} {
		if !slices.Contains(ring, k) {
			t.Errorf("kunci ringkasan %s tidak ada di view", k)
		}
	}
	for _, k := range strings.Split(kolomRingkasan, ", ") {
		if !slices.Contains(ring, k) {
			t.Errorf("kolom ringkasan %s tidak ada di view", k)
		}
	}
	for _, k := range []string{JSONIDUsedBy, JSONUsedBy, JSONGender, JSONContract, JSONAge, JSONRate} {
		if !slices.Contains(rate, k) {
			t.Errorf("kunci rate %s tidak ada di view", k)
		}
	}
	for _, k := range strings.Split(kolomRate, ", ") {
		if !slices.Contains(rate, k) {
			t.Errorf("kolom rate %s tidak ada di view", k)
		}
	}
	for _, q := range []string{SqlSisipRingkasan("T"), SqlUbahRingkasan("T"), SqlDaftar("V", "", false), kolomRingkasan} {
		for _, k := range []string{"'TYPE'", "'FLAG'", "TYPE,", "FLAG"} {
			if strings.Contains(q, k) {
				t.Errorf("%q memuat %s:\n%s", k, k, q)
			}
		}
	}
}
