package repository

import (
	"slices"
	"strings"
	"testing"
)

// RALAT R1: definisi view RATE_LIFE_SUMMARY (ALL_VIEWS, dibaca WO 06-10-2026) = enam kolom ini.
func TestKolomViewRingkasanEnamKolom(t *testing.T) {
	mau := []string{"ID", "USEDBY", "TYPE", "MODIFIEDDATE", "OPERATORID", "FLAG"}
	if !slices.Equal(KolomViewRingkasan, mau) {
		t.Fatalf("kolom view %v, mau %v", KolomViewRingkasan, mau)
	}
}

// Setiap kunci JSON yang ditulis dan setiap kolom yang dibaca ADA di definisi view - nol nama karangan.
// TYPE dan FLAG tidak dirujuk XML InboxSummaryRIRate: tidak ditulis, tidak dibaca.
func TestKunciDanKolomCocokDenganView(t *testing.T) {
	ring, rate := KolomViewRingkasan, KolomViewRate
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
