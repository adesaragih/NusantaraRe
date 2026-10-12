package models

// POLICY_HOLDER / PLAN peserta dari data polis (keputusan work owner 05-10-2026).

import (
	"errors"
	"strings"
	"testing"
)

func TestTerapkanDataPolisKePeserta(t *testing.T) {
	baris := []BarisUnggah{
		{Nomor: 1, Nilai: map[string]string{"POLICY_HOLDER": "UJI-CSV-LAIN", "PLAN": "UJI-PLAN-CSV"}},
		{Nomor: 2, Nilai: map[string]string{"NAME_OF_INSURED": "UJI-B"}}, // kolom tidak ada di CSV
		{Nomor: 3},
	}
	if err := TerapkanDataPolisKePeserta("  UJI-PEMEGANG  ", " UJI-GTL ", baris); err != nil {
		t.Fatal(err)
	}
	for _, b := range baris {
		if b.Nilai["POLICY_HOLDER"] != "UJI-PEMEGANG" || b.Nilai["PLAN"] != "UJI-GTL" {
			t.Errorf("baris %d: %+v", b.Nomor, b.Nilai)
		}
	}
	if baris[1].Nilai["NAME_OF_INSURED"] != "UJI-B" {
		t.Error("kolom lain ikut berubah")
	}
}

func TestDataPolisPesertaKosongDitolak(t *testing.T) {
	for _, k := range []struct{ ph, plan, pesan string }{
		{"", "UJI-GTL", "Policy Holder is empty on the policy"},
		{"UJI-PEMEGANG", "  ", "Class of Business is empty on the policy"},
		{"", "", "Policy Holder and Class of Business are empty on the policy"},
	} {
		baris := []BarisUnggah{{Nomor: 1, Nilai: map[string]string{"PLAN": "UJI-CSV"}}}
		err := TerapkanDataPolisKePeserta(k.ph, k.plan, baris)
		if !errors.Is(err, ErrDataPolisPesertaKosong) || !strings.Contains(err.Error(), k.pesan) {
			t.Errorf("%q/%q: %v, mau memuat %q", k.ph, k.plan, err, k.pesan)
		}
		if baris[0].Nilai["PLAN"] != "UJI-CSV" {
			t.Error("baris disentuh walau ditolak")
		}
	}
}

// POLICY_HOLDER dan PLAN bukan kolom wajib judul / baris di Type mana pun.
func TestPolicyHolderDanPlanTidakWajib(t *testing.T) {
	for _, tipe := range []string{"QR", "QP", "TR", "TP"} {
		k, err := KolomWajibUnggah(tipe)
		if err != nil {
			t.Fatal(err)
		}
		for _, x := range k {
			if x == "POLICY_HOLDER" || x == "PLAN" {
				t.Errorf("%s: %s wajib", tipe, x)
			}
		}
	}
}
