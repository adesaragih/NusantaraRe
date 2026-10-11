package repository

import (
	"testing"

	"nusantarare/modul/treatyinadjustment/backend/models"
)

// Halaman Actual (`#AKTUAL`) kembali ke sisi New sebagai kunci BERTITIK —
// medan, larik datar, dan pohon; yang sudah ada di sisi New tidak ditimpa.
func TestGabungAktualBertitik(t *testing.T) {
	baru := sisiKosong()
	baru.Medan["RNMShare"] = "10"
	baru.Medan["ActualValue.TotalEgnpiAmount"] = "sudah"
	aktual := sisiKosong()
	aktual.Medan["TotalEgnpiAmount"] = "500"
	aktual.Medan["RNMShare"] = "12"
	aktual.Larik["EGNPI"] = []map[string]string{{"Amount": "500"}}
	aktual.Pohon["Share"] = []map[string]any{{"RNMShare": "12", "GrossPremiumList": []map[string]any{{"Value": "5"}}}}

	gabungAktual(baru, aktual)

	if baru.Medan["RNMShare"] != "10" {
		t.Errorf("medan sisi New tertimpa: %q", baru.Medan["RNMShare"])
	}
	if baru.Medan["ActualValue.RNMShare"] != "12" {
		t.Errorf("ActualValue.RNMShare = %q", baru.Medan["ActualValue.RNMShare"])
	}
	if baru.Medan["ActualValue.TotalEgnpiAmount"] != "sudah" {
		t.Errorf("kunci bertitik yang sudah ada tertimpa: %q", baru.Medan["ActualValue.TotalEgnpiAmount"])
	}
	if got := baru.Larik["ActualValue.EGNPI"]; len(got) != 1 || got[0]["Amount"] != "500" {
		t.Errorf("ActualValue.EGNPI = %v", got)
	}
	if got := baru.Pohon["ActualValue.Share"]; len(got) != 1 {
		t.Errorf("ActualValue.Share (pohon) = %v", got)
	}
	var _ models.SisiPenyesuaian = baru
}
