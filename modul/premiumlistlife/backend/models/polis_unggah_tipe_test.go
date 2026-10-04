package models

// Kolom wajib unggahan per Type - keputusan work owner 03-10-2026.

import (
	"errors"
	"slices"
	"sort"
	"testing"
)

func TestKolomWajibPerTipePersisDaftarWorkOwner(t *testing.T) {
	umum := []string{"POLICY_NO", "CERTIFICATE_NO", "NAME_OF_INSURED", "DOB", "ENTRY_AGE",
		"BEGIN_DATE", "EXPIRED_DATE", "CURRENCY", "SUM_INSURED"}
	for tipe, valuasi := range map[string][]string{
		"QR": {"GROSS_VALUATION_BEGIN_DATE", "GROSS_VALUATION_EXPIRED_DATE"},
		"QP": {"GROSS_VALUATION_BEGIN_DATE", "GROSS_VALUATION_EXPIRED_DATE"},
		"TR": {"RETROCESSION_VALUATION_BEGIN_DATE", "RETROCESSION_VALUATION_EXPIRED_DATE"},
		"TP": {"RETROCESSION_VALUATION_BEGIN_DATE", "RETROCESSION_VALUATION_EXPIRED_DATE"},
	} {
		mau := append(slices.Clone(umum), valuasi...)
		sort.Strings(mau)
		got, err := KolomWajibUnggah(tipe)
		if err != nil || !slices.Equal(got, mau) {
			t.Errorf("%s: %v, %v - mau %v", tipe, got, err, mau)
		}
	}
}

func TestTipeKosongAtauAsingDitolak(t *testing.T) {
	for _, tipe := range []string{"", "  ", "XX"} {
		if _, err := ValidasiUnggah(tipe, []BarisUnggah{barisSah(1)}); !errors.Is(err, ErrTipeUnggahTakDikenal) {
			t.Errorf("Type %q: galat %v, mau ErrTipeUnggahTakDikenal", tipe, err)
		}
	}
}

func TestValuasiMengikutiTipe(t *testing.T) {
	b := barisSah(1)
	delete(b.Nilai, "RETROCESSION_VALUATION_BEGIN_DATE")
	delete(b.Nilai, "RETROCESSION_VALUATION_EXPIRED_DATE")
	if h, _ := ValidasiUnggah("QR", []BarisUnggah{b}); !h.Lolos() {
		t.Errorf("QR tanpa Retrocession Valuation ditolak: %+v", h.Ditolak)
	}
	if h, _ := ValidasiUnggah("TP", []BarisUnggah{b}); h.Lolos() || len(h.Ditolak) != 2 {
		t.Errorf("TP tanpa Retrocession Valuation mau 2 penolakan: %+v", h.Ditolak)
	}
	c := barisSah(2)
	delete(c.Nilai, "GROSS_VALUATION_BEGIN_DATE")
	delete(c.Nilai, "GROSS_VALUATION_EXPIRED_DATE")
	if h, _ := ValidasiUnggah("TR", []BarisUnggah{c}); !h.Lolos() {
		t.Errorf("TR tanpa Gross Valuation ditolak: %+v", h.Ditolak)
	}
}

func TestKolomWajibLamaTidakLagiWajib(t *testing.T) {
	b := barisSah(1)
	for _, k := range []string{"PLAN", "WPC", "CEDING_RETENTION", "SUM_REASURED",
		"SHARE_NUSANTARA_RE", "GROSS_PREMIUM", "NET_PREMIUM"} {
		delete(b.Nilai, k)
	}
	baris := []BarisUnggah{b}
	IsiNolUangKosong(baris)
	if h := validasiQR(baris); !h.Lolos() {
		t.Errorf("kolom wajib lama yang kosong ditolak: %+v", h.Ditolak)
	}
}

// SUM_INSURED wajib: sel kosongnya DITOLAK, bukan diisi 0 lebih dahulu.
func TestSumInsuredKosongDitolakWalauAdaIsiNol(t *testing.T) {
	b := barisSah(1)
	b.Nilai["SUM_INSURED"] = " "
	baris := []BarisUnggah{b}
	IsiNolUangKosong(baris)
	h := validasiQR(baris)
	if len(h.Ditolak) != 1 || h.Ditolak[0].Kolom != "SUM_INSURED" || h.Ditolak[0].Sebab != "kolom kosong" {
		t.Errorf("SUM_INSURED kosong: %+v", h.Ditolak)
	}
}

func TestEntryAgeHarusBilanganBulat(t *testing.T) {
	for v, sah := range map[string]bool{"30": true, "0": true, "3a": false, "30.5": false, "-1": false} {
		b := barisSah(1)
		b.Nilai["ENTRY_AGE"] = v
		if h := validasiQR([]BarisUnggah{b}); h.Lolos() != sah {
			t.Errorf("ENTRY_AGE %q: lolos=%v, mau %v (%+v)", v, h.Lolos(), sah, h.Ditolak)
		}
	}
}
