package services

import (
	"context"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// Ralat butir 78.2 (work owner 09-10-2026): Policy Type = standard value 0 / 1 / 2 (DDL\PolicyType.xml).

// TestKodePolicyType - label -> kode; kode / nilai lain apa adanya.
func TestKodePolicyType(t *testing.T) {
	for masuk, mau := range map[string]string{"Individual Policy": "0", "Master Policy": "1", "Declaration Policy": "2", "0": "0",
		"1": "1", "2": "2", "": "", "master policy": "master policy", "X": "X"} {
		if got := kodePolicyType(masuk); got != mau {
			t.Errorf("%q: %q mau %q", masuk, got, mau)
		}
	}
	for k, sah := range map[string]bool{"": true, "0": true, "1": true, "2": true, "3": false, "Master Policy": false} {
		if kodePolicyTypeSah[k] != sah {
			t.Errorf("sah %q", k)
		}
	}
}

// TestPolicyTypeLabelLama - data lama berlabel: GET kasus memberi kode, selesai Policy menulis FACINPRODUCTION.POLICYTYPE "1".
func TestPolicyTypeLabelLama(t *testing.T) {
	ctx := context.Background()
	st := submitUji()
	svc := layananPolicy(st, posisiTeamLeader, spreadingPolicy("1000", "10001", "QS UJI"), 0)
	svc.kasus.(kasusTiruan).ada["UJI-NB-1"].General.PolicyType = "Master Policy"
	if k, err := svc.BacaKasus(ctx, "UJI-NB-1"); err != nil || k.General.PolicyType != "1" {
		t.Errorf("GET kasus: %v %q", err, k.General.PolicyType)
	}
	if _, err := svc.AjukanSubmit(ctx, pelakuTL, "UJI-NB-1", models.IsiSubmit{EmailType: "1", Comment: "UJI"}); err != nil ||
		len(st.tulis) != 1 || st.tulis[0].Policy == nil || len(st.tulis[0].Policy.Produksi) == 0 {
		t.Fatalf("%v", err)
	}
	for _, b := range st.tulis[0].Policy.Produksi {
		if b.PolicyType != "1" {
			t.Errorf("POLICYTYPE %q", b.PolicyType)
		}
	}
}

// TestProteksiPremiNolMasterPolicy - ProtectPremiPolicy_Act langkah 7 / 11 dilewati bila PolicyType == 1 (Master Policy, label lama
// maupun kode); instrumen: Individual Policy tetap ditahan.
func TestProteksiPremiNolMasterPolicy(t *testing.T) {
	ctx := context.Background()
	for pt, ditahan := range map[string]bool{"0": true, "Individual Policy": true, "": true, "1": false, "Master Policy": false} {
		svc := layananSubmit(submitUji())
		svc.spreading.(*spreadingTiruan).baca = []models.KasusSpreading{kasusDuaItem("0", "USD")}
		svc.kasus.(kasusTiruan).ada["UJI-NB-1"].General.PolicyType = pt
		a, err := svc.ProteksiAccept(ctx, "UJI-NB-1", nil)
		if err != nil || dalam(a, "Premi on paymnet can't be 0") != ditahan {
			t.Errorf("%q: %v %q", pt, err, a)
		}
	}
}
