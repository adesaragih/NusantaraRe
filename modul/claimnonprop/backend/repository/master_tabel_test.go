package repository

import (
	"testing"

	"nusantarare/modul/claimnonprop/backend/models"
)

// TestAngkaMasterMeluruskanFormatTampilan - angka tabel `T_TREATY_*` berformat tampilan Indonesia diluruskan, angka
// polos tidak disentuh (DEV 11-10-2026: Limit "7.100.000.000,00", RNMSHARE "2,5").
func TestAngkaMasterMeluruskanFormatTampilan(t *testing.T) {
	for masuk, harap := range map[string]string{
		"7.100.000.000,00":   "7100000000.00",
		"150.000.000,00":     "150000000.00",
		"2,5":                "2.5",
		"-1.250.000,75":      "-1250000.75",
		"7.100.000":          "7100000",
		"4675000000.00":      "4675000000.00",
		"2893967952.0000075": "2893967952.0000075",
		"1.500":              "1.500",
		"22.5":               "22.5",
		"40":                 "40",
		"":                   "",
	} {
		if got := angkaMaster(masuk); got != harap {
			t.Errorf("angkaMaster(%q) = %q, harap %q", masuk, got, harap)
		}
	}
}

// TestModePembukuanMasterNonProp - Accounting Mode master Non Prop dari tabel = AccountingModeNonProp bila
// AccountingMode kosong (CLMNP-000017 / master 1002320); nilai yang ada dipertahankan.
func TestModePembukuanMasterNonProp(t *testing.T) {
	for _, k := range []struct{ mode, nonProp, harap string }{
		{"", "loss", "loss"}, {" ", "risk", "risk"}, {"accounting", "loss", "accounting"}, {"underwriting", "", "underwriting"},
		{"", "", ""},
	} {
		if got := modePembukuanMaster(k.mode, k.nonProp); got != k.harap {
			t.Errorf("modePembukuanMaster(%q, %q) = %q, harap %q", k.mode, k.nonProp, got, k.harap)
		}
	}
}

// TestLengkapiJenisSpreadingDariRincianTunggal - Share tanpa jenis + rincian ORS tunggal = jenis ORS / 10007 (bentuk
// master Pega ber-ORS); jenis yang ada dan rincian majemuk tidak disentuh.
func TestLengkapiJenisSpreadingDariRincianTunggal(t *testing.T) {
	ors := models.SpreadingMaster{ReinsTypeID: "10007", ReinsTypeName: "ORS", Pct: "40"}
	sh := models.ShareXOL{SpreadingTotalPctXOL: "40", SpreadingListXOL: []models.SpreadingMaster{ors}}
	lengkapiJenisSpreading(&sh)
	if sh.SpreadingTypeIDXOL != "10007" || sh.SpreadingTypeXOL != "ORS" {
		t.Errorf("jenis = %q / %q, harap 10007 / ORS", sh.SpreadingTypeIDXOL, sh.SpreadingTypeXOL)
	}
	bernama := models.ShareXOL{SpreadingTypeIDXOL: "10260", SpreadingTypeXOL: "2026 QS 212M TRT HR",
		SpreadingListXOL: []models.SpreadingMaster{ors}}
	lengkapiJenisSpreading(&bernama)
	if bernama.SpreadingTypeXOL != "2026 QS 212M TRT HR" {
		t.Errorf("jenis bernama tertimpa: %q", bernama.SpreadingTypeXOL)
	}
	majemuk := models.ShareXOL{SpreadingListXOL: []models.SpreadingMaster{
		{ReinsTypeID: "10028", ReinsTypeName: "QS (OR)", Pct: "40"}, {ReinsTypeID: "10004", ReinsTypeName: "QS (R/I)", Pct: "60"}}}
	lengkapiJenisSpreading(&majemuk)
	if majemuk.SpreadingTypeXOL != "" || majemuk.SpreadingTypeIDXOL != "" {
		t.Errorf("rincian majemuk diberi jenis %q", majemuk.SpreadingTypeXOL)
	}
}
