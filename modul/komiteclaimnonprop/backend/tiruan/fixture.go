package tiruan

// Untuk apa berkas ini: FIXTURE UJI- - satu kasus klaim treaty non proporsional induk (CLMNP-) dengan satu baris
// akseptasi yang sudah diserahkan (layer XoL non-UR + UR, Spreading In), dan satu kasus komite KMTNP- bertangga
// UJI-K1..UJI-Kn (seperti `CreateChildKomiteCNP_Act` Claim Non Prop). Dipakai uji seam HTTP dan server tiruan uji
// manual. Nol data orang / klaim asli.

import (
	"fmt"
	"time"

	"nusantarare/modul/komiteclaimnonprop/backend/models"
)

// ID fixture.
const (
	KlaimUji  = "CLMNP-UJI001"
	AdjUji    = "UJI-ADJ-1"
	KomiteUji = "KMTNP-UJI001"
	// PembuatUji - admin klaim yang menyerahkan ke komite.
	PembuatUji = "UJI-ADMIN"
)

// SaatUji - jam tetap uji (Jakarta).
var SaatUji = time.Date(2026, 10, 8, 10, 0, 0, 0, models.Jakarta)

// PenyetujuUji - operator tangga ke-n (1..).
func PenyetujuUji(n int) string { return fmt.Sprintf("UJI-K%d", n) }

// SiapkanUji melahirkan klaim induk + kasus komite `id` bertangga `n` penyetuju (akseptasi `adjID`).
func SiapkanUji(g *Gudang, a *Acuan, id, adjID string, n int) {
	nilai := map[string]string{
		"pyID":                                   KlaimUji,
		"ClaimData.NoClaim":                      "UJI-K-0001",
		"ClaimData.PolicyData.PolicyNo":          "UJI-POL-0001",
		"ClaimData.PolicyNo":                     "UJI-POLCED-0001",
		"ClaimData.PolicyData.StartDateTime":     "2026-01-01",
		"ClaimData.PolicyData.EndDateTime":       "2026-12-31",
		"ClaimData.InsuredName":                  "UJI TERTANGGUNG",
		"ClaimData.CauseOfLoss":                  "UJI SEBAB",
		"ClaimData.CauseOfLossID":                "UJI-COL",
		"ClaimData.ReportType":                   "1",
		"ClaimData.Location":                     "UJI LOKASI",
		"ClaimData.DateOfLoss":                   "2026-09-15",
		"ClaimData.ReportDate":                   "2026-09-20",
		"ClaimData.DateReceived":                 "2026-09-21",
		"ClaimData.ReporterName":                 "UJI PELAPOR",
		"ClaimData.ReporterStatus":               "1",
		"ClaimData.NoPla":                        "UJI-PLA-0001",
		"ClaimData.IDMaster":                     "UJI-MASTER",
		"ClaimData.Occupation":                   "UJI OKUPASI",
		"ClaimData.Payable":                      "1",
		"ClaimData.PayableTo":                    "UJI CEDING",
		"OfferFacIn.QuotationData.BusinessName":  "UJI COB",
		"OfferFacIn.QuotationData.BusinessOldId": "123",
		"TreatyInMaster.Ceding":                  "UJI CEDING",
		"TreatyInMaster.CedingID":                "UJI-CED",
		"TreatyInMaster.LeadingReinsSource":      "UJI SOB",
		"TreatyInMaster.RNMShare":                "25",
	}
	adj := func(anak string) string { return "ClaimData.AdjustmentList(2)." + anak }
	daftar := map[string][]map[string]string{
		"ClaimData.AdjustmentList": {
			{"ID": "UJI-ADJ-0", "Type": "1", "PaymentType": "1", "Currency": "IDR", "AcceptanceStatus": "1",
				"AcceptedNo": "UJIA123.09.2026.TX00009"},
			{"ID": adjID, "Type": "1", "PaymentType": "1", "Currency": "IDR", "PersenRNM": "25", "IsKomite": "1",
				"DirectToKasir": "true", "NameOfBank": "UJI BANK", "BranchOfBank": "UJI CABANG", "NoAccount": "UJI-123-456",
				"pxCreateOperator": PembuatUji, "DataCommitteeTreaty.CircumCauseOfLoss": "UJI KRONOLOGI",
				"DataCommitteeTreaty.Remarks": "UJI CATATAN KOMITE"},
		},
		adj(models.AnakClaimAcc): {
			{"CurrencyID": "IDR", "Currency": "IDR", "AltValue": "1", "Value": "400", "TPL": "0", "AdjusterFee": "10",
				"Salvage": "0", "CNPOthersFee": "0", "PctProrateClaim": "100", "USD": "400", "ClaimAmountCedant": "400"},
		},
		adj(models.AnakLossAlloc): {
			{"CurrencyID": "IDR", "Currency": "IDR", "TreatyName": "UJI XOL 1", "ClaimPercentage": "100",
				"ClaimAmountAdjust": "400", "AdjusterFee": "10", "Salvage": "0", "CNPOthersFee": "0", "CNPFlagXOL": "true"},
		},
		adj(models.AnakXOL): {
			{"CurrencyID": "IDR", "Currency": "IDR", "TreatyName": models.TreatyUR, "TotalClaim": "100",
				"ClaimPercentage": "0", "ClaimSpreaded": "0"},
			{"CurrencyID": "IDR", "Currency": "IDR", "TreatyName": "UJI XOL 1", "TreatyType": "UJI-XOL1",
				"TotalClaim": "300", "ClaimEstimation": "300", "ClaimPercentage": "25", "ClaimSpreaded": "75",
				"AdjusterFee": "2.5", "Salvage": "0", "CNPOthersFee": "0", "CNPReinstatement": "10",
				"CNPReinstatementRNM": "2.5", "Kurs": "1", "KursIDR": "1"},
		},
		adj(models.AnakSpreadIn): {
			{"TreatyType": "UJI-QS", "TreatyName": "UJI Quota Share", "Currency": "IDR", "CurrencyID": "IDR",
				"SharePercentage": "100", "ClaimSpreaded": "75", "TotalClaim": "77.5", "PremiumSpreaded": "2.5",
				"NoAccount": "UJI-123-456"},
		},
	}
	g.Klaim.Setel(KlaimUji, nilai, daftar)
	var tangga []models.Anggota
	for i := 1; i <= n; i++ {
		tangga = append(tangga, models.Anggota{OperatorID: PenyetujuUji(i), Jabatan: fmt.Sprintf("UJI-JABATAN-%d", i)})
	}
	g.Lahirkan(id, KlaimUji, adjID, PembuatUji, "UJI Admin", tangga, SaatUji)
	if a != nil {
		a.Nama[PenyetujuUji(1)] = "UJI Penyetuju Satu"
	}
}
