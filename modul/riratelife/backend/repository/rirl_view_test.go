package repository

import (
	"slices"
	"strconv"
	"strings"
	"testing"

	"nusantarare/modul/riratelife/backend/models"
)

// RALAT R1: definisi view RATE_LIFE_SUMMARY (ALL_VIEWS, dibaca WO 06-10-2026) = enam kolom ini (sejarah; dipulihkan
// jalur mundur 926). Kolom ringkasan M_RATE_LIFE_SUMMARY = kolom ini TANPA FLAG (RALAT R5, R6).
func TestKolomViewRingkasanEnamKolom(t *testing.T) {
	mau := []string{"ID", "USEDBY", "TYPE", "MODIFIEDDATE", "OPERATORID", "FLAG"}
	if !slices.Equal(KolomViewRingkasan, mau) {
		t.Fatalf("kolom view %v, mau %v", KolomViewRingkasan, mau)
	}
	if !slices.Equal(KolomRingkasan, mau[:5]) {
		t.Errorf("kolom ringkasan %v, mau %v", KolomRingkasan, mau[:5])
	}
}

// Setiap kolom ringkasan dan rincian yang ditulis / dibaca ADA di definisi kolomnya - nol nama karangan.
// TYPE dan FLAG tidak dirujuk XML InboxSummaryRIRate: tidak ditulis, tidak dibaca.
func TestKunciDanKolomCocokDenganView(t *testing.T) {
	ring, rate := KolomViewRingkasan, KolomRate
	for _, k := range []string{KolomUsedBy, KolomOperatorID, KolomModified} {
		if !slices.Contains(ring, k) {
			t.Errorf("kolom tulis ringkasan %s tidak ada di view", k)
		}
	}
	for _, k := range strings.Split(kolomRingkasan, ", ") {
		if !slices.Contains(ring, k) {
			t.Errorf("kolom ringkasan %s tidak ada di view", k)
		}
	}
	for _, k := range strings.Split(kolomRate, ", ") {
		if !slices.Contains(rate, k) {
			t.Errorf("kolom rate %s tidak ada di M_RATE_LIFE", k)
		}
	}
	for _, q := range []string{SqlSisipRingkasan("T"), SqlUbahRate("T"), SqlSisipRate("T"), SqlDaftar("V", "", false), kolomRingkasan, kolomRate} {
		for _, k := range []string{"'TYPE'", "'FLAG'", "TYPE,", "FLAG"} {
			if strings.Contains(q, k) {
				t.Errorf("%q memuat %s:\n%s", k, k, q)
			}
		}
	}
}

// RALAT R7: M_RATE_LIFE flat 8 kolom dengan lebar ALL_TAB_COLUMNS DEV 10-10-2026; batas validasi Go muat di kolomnya.
func TestLebarKolomRateFlat(t *testing.T) {
	if !slices.Equal(KolomRate, []string{"ID", "IDUSEDBY", "USEDBY", "TYPE", "GENDER", "CONTRACT", "AGE", "RATE"}) {
		t.Fatalf("kolom M_RATE_LIFE %v", KolomRate)
	}
	for _, k := range KolomRate {
		if LebarKolomRate[k] == 0 {
			t.Errorf("lebar %s tidak tercatat", k)
		}
	}
	if models.LebarRate != LebarKolomRate["RATE"] || models.BatasNama > LebarKolomRate["USEDBY"] ||
		len(strconv.Itoa(models.UmurMaks)) > LebarKolomRate["AGE"] || len(strconv.Itoa(models.UmurMaks)) > LebarKolomRate["CONTRACT"] {
		t.Errorf("batas Go melebihi kolom: RATE %d, USEDBY %d, umur %d", models.LebarRate, models.BatasNama, models.UmurMaks)
	}
	for _, gd := range models.Gender {
		if len(gd) > LebarKolomRate["GENDER"] {
			t.Errorf("gender %q", gd)
		}
	}
}
