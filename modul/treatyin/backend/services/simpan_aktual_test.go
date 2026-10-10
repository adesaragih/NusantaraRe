package services_test

// Halaman `ActualValue` Save layar Adjustment — MASTERID sendiri
// (`repository.AkhiranSisiAktual`), keputusan pemakai 8 Oktober 2026.

import (
	"context"
	"testing"

	"nusantarare/modul/treatyin/backend/services"
)

// EDMState 3 (Adjust Premium): `ActualValue` ISIAN pemakai — dikirim apa
// adanya ke sisi Actual, dan TIDAK ikut dokumen sisi New.
func TestSaveEDM3MendaratkanActualValueSendiri(t *testing.T) {
	g := gudangEDM()
	_, err := services.LayananDengan(g).SimpanPenyesuaian(context.Background(), admin, services.MasukanPenyesuaian{
		ID: "1001001/R02", Draf: true,
		Baru: services.SisiKiriman{
			Medan: map[string]any{"OLDID": "1001001/R01", "EDMState": "3", "RNMShare": "10", "ActualValue.TotalEgnpiAmount": "500"},
			Larik: map[string]any{"ActualValue.EGNPI": []any{map[string]any{"Amount": "500", "Currency": "IDR"}}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.disimpanEDM) != 1 {
		t.Fatalf("rencana = %d, ingin 1", len(g.disimpanEDM))
	}
	r := g.disimpanEDM[0]
	if _, ada := r.Baru["ActualValue"]; ada {
		t.Error("ActualValue masih ikut dokumen sisi New")
	}
	if r.Aktual["TotalEgnpiAmount"] != "500" {
		t.Errorf("Aktual.TotalEgnpiAmount = %v", r.Aktual["TotalEgnpiAmount"])
	}
	egnpi, _ := r.Aktual["EGNPI"].([]any)
	if len(egnpi) != 1 {
		t.Errorf("Aktual.EGNPI = %v", r.Aktual["EGNPI"])
	}
	if r.Aktual["RNMShare"] != nil {
		t.Error("EDMState 3 TIDAK menyalin TreatyIn ke ActualValue (SaveTreatyIn_EDM_Act [2] lompat ke jmp)")
	}
}

// EDMState 1/2: `SaveTreatyIn_EDM_Act` [2]–[4] — ActualValue = salinan
// TreatyIn tanpa OLDDATA, ValueDifference, dan ActualValue lamanya.
func TestSaveEDM12MenyalinTreatyInKeActualValue(t *testing.T) {
	g := gudangEDM()
	_, err := services.LayananDengan(g).SimpanPenyesuaian(context.Background(), admin, services.MasukanPenyesuaian{
		ID: "1001001/R02", Draf: true,
		Baru: services.SisiKiriman{
			Medan: map[string]any{"OLDID": "1001001/R01", "EDMState": "2", "RNMShare": "12", "ValueDifference.RNMShare": "2", "ActualValue.Lama": "x"},
			Larik: map[string]any{"Portfolio": []any{map[string]any{"Description": "p"}}},
		},
		Lama: &services.SisiKiriman{Medan: map[string]any{"RNMShare": "10"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	r := g.disimpanEDM[0]
	if r.Aktual["RNMShare"] != "12" || r.Aktual["EDMState"] != "2" {
		t.Errorf("salinan skalar salah: %v", r.Aktual)
	}
	if p, _ := r.Aktual["Portfolio"].([]any); len(p) != 1 {
		t.Errorf("salinan larik salah: %v", r.Aktual["Portfolio"])
	}
	for _, k := range []string{"OLDDATA", "ValueDifference", "ActualValue"} {
		if _, ada := r.Aktual[k]; ada {
			t.Errorf("%s ikut tersalin ke ActualValue", k)
		}
	}
	// Sisi New tetap membawa ValueDifference-nya.
	if _, ada := r.Baru["ValueDifference"]; !ada {
		t.Error("ValueDifference hilang dari sisi New")
	}
}
