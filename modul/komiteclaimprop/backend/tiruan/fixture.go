package tiruan

// Untuk apa berkas ini: FIXTURE UJI- - satu kasus klaim treaty induk (CLMP-) dengan satu baris adjustment yang sudah
// diserahkan, dan satu kasus komite TKMT- bertangga UJI-K1..UJI-Kn (seperti `AddKomiteTreatyChild_ACT` Claim Prop).
// Dipakai uji seam HTTP dan server tiruan `ujimanual`. Nol data orang / klaim asli.

import (
	"fmt"
	"time"

	"nusantarare/modul/komiteclaimprop/backend/models"
)

// ID fixture.
const (
	KlaimUji  = "CLMP-UJI001"
	AdjUji    = "UJI-ADJ-1"
	KomiteUji = "TKMT-UJI001"
	// PembuatUji - admin klaim yang menyerahkan ke komite.
	PembuatUji = "UJI-ADMIN"
)

// SaatUji - jam tetap uji (Jakarta).
var SaatUji = time.Date(2026, 10, 8, 10, 0, 0, 0, models.Jakarta)

// PenyetujuUji - operator tangga ke-n (1..).
func PenyetujuUji(n int) string { return fmt.Sprintf("UJI-K%d", n) }

// SiapkanUji melahirkan klaim induk + kasus komite `id` bertangga `n` penyetuju (adjustment `adjID`).
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
		"ClaimData.StsKatastrofe":                "Non-Catastrophe",
		"ClaimData.NonKatastrofeType":            "Claim",
		"ClaimData.DateOfLoss":                   "2026-09-15",
		"ClaimData.ReportDate":                   "2026-09-20",
		"ClaimData.DateReceived":                 "2026-09-21",
		"ClaimData.ReporterName":                 "UJI PELAPOR",
		"ClaimData.ReporterStatus":               "1",
		"ClaimData.TreatyGroupID":                "UJI-TG",
		"ClaimData.IDMaster":                     "UJI-MASTER",
		"ClaimData.Occupation":                   "UJI OKUPASI",
		"ClaimData.TotalGrossEstimateIDR":        "1000000",
		"ClaimData.TotalEstimasiIDR":             "250000",
		"OfferFacIn.QuotationData.BusinessName":  "UJI COB",
		"OfferFacIn.QuotationData.BusinessOldId": "12",
		"TreatyInMaster.Ceding":                  "UJI CEDING",
		"TreatyInMaster.CedingID":                "UJI-CED",
		"TreatyInMaster.LeadingReinsSource":      "UJI SOB",
		"TreatyInMaster.RNMShareP":               "25",
		"AktifButton":                            "1",
	}
	daftar := map[string][]map[string]string{
		"ClaimData.AdjustmentList": {
			{"ID": "UJI-ADJ-0", "Type": "1", "Currency": "IDR", "CurrencyID": "IDR", "GrossAdjustment": "100",
				"AdjustmentValue": "25", "KursIDR": "1", "AcceptanceStatus": "1", "AcceptedNo": "UJIA12.09.2026.TP00009"},
			{"ID": adjID, "Type": "1", "Currency": "IDR", "CurrencyID": "IDR", "GrossAdjustment": "400",
				"AdjustmentValue": "100", "KursIDR": "1", "PersenRNM": "25", "IsKomite": "1", "KomiteID": id,
				"DataCommitteeTreaty.Remarks": "UJI CATATAN KOMITE"},
		},
		"ClaimData.AdjustmentList(2).SpreadingAdjustment": {
			{"TreatyType": "UJI-QS", "TreatyName": "UJI Quota Share", "Currency": "IDR", "SharePercentage": "100",
				"ClaimSpreaded": "100"},
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
		a.Tahun["UJI-TG|20260101"] = "2026"
		a.Batas["2026|UJI-TG|UJI-QS"] = "50"
		a.Retro["UJI-QS|2026|UJI-TG"] = []models.BarisRetro{{ReinsurerID: "UJI-R1", Name: "UJI RETRO", RiComm: "10",
			PctShare: "20"}}
	}
}
