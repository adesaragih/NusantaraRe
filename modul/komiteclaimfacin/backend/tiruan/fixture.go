package tiruan

// Untuk apa berkas ini: FIXTURE `UJI-` - roster FACIN ke workbasket (migrasi 640; JABATAN dan batas DEV 10-10-2026,
// akun orang tidak disalin) dan kasus klaim induk Fac In (satu objek, dua item, adjustment TT2 di item kedua).

import (
	"time"

	"nusantarare/modul/komiteclaimfacin/backend/models"
)

// RosterUji - roster EMAILKOMITE STS_KLAIM FACIN sesudah migrasi 640 (DEGREE 1-5, batas DEV).
func RosterUji() []models.AnggotaRoster {
	return []models.AnggotaRoster{
		{ID: "1", OperatorID: models.WorkbasketSPVA, Jabatan: "Claim Supervisor", Degree: "1", LimitBottom: "-9999999999999"},
		{ID: "2", OperatorID: models.WorkbasketDeptHead, Jabatan: "Claim Dept. Head", Degree: "2", LimitBottom: "57750001"},
		{ID: "10", OperatorID: "ReasClaimTechDivHead", Jabatan: "Technic Div. Head", Degree: "3", LimitBottom: "189750001"},
		{ID: "3", OperatorID: "ReasClaimOpsDir", Jabatan: "Operational Director", Degree: "4", LimitBottom: "495000001"},
		{ID: "4", OperatorID: "ReasClaimTechDir", Jabatan: "Technical Director", Degree: "5", LimitBottom: "660000001"},
	}
}

// Nilai tetap fixture.
const (
	KlaimUji   = "CLM-900001"
	AdjUji     = "UJI-ADJ-21"
	KomiteUji  = "KMT-900001"
	PembuatUji = "UJI-ADMIN"
)

// KlaimFacInUji - kasus klaim induk Fire: objek 1 (dua item), adjustment PT 1 bernilai IDR `nilaiIDR` di item 2 (sudah
// diserahkan ke `KomiteUji`), estimasi item 1 ber-PrintFaceClaim 1 dan item 2 ber-PrintFaceClaim 0.
func KlaimFacInUji(nilaiIDR string) (map[string]string, map[string][]map[string]string) {
	nilai := map[string]string{
		"ClaimData.NoClaim": "UJI-K-0001", "ClaimData.CauseOfLoss": "UJI KEBAKARAN", "ClaimData.CauseOfLossID": "UJI-COL",
		"ClaimData.DateOfLoss": "2026-09-30", "ClaimData.Location": "UJI LOKASI", "ClaimData.Remark": "UJI ALASAN",
		models.OQ + "BusinessType": "Fire", models.OQ + "BusinessOldId": "12", models.OQ + "InsuredName": "UJI TERTANGGUNG",
		models.OQ + "BusinessName": "UJI FIRE", models.OQ + "CedingCo": "UJI-CEDING-01",
		"OfferFacIn.PolicyData.PolicyNo": "UJI-RNM-F.001", "OfferFacIn.PercentShare": "25", "IsTreatyIn": "0",
		"pxCreateOperator": PembuatUji, "pxCreateOpName": "UJI Admin",
	}
	daftar := map[string][]map[string]string{
		models.DaftarObjek: {{"ID": "UJI-OBJ-1", "ObjectID": "1", "ObjectName": "UJI GEDUNG", "PlaStatus": "1"}},
		models.DaftarItem(1): {
			{"ID": "UJI-IT-1", "ObjectItemName": "UJI ITEM 1", "CoverageID": "UJI-COV", "CoverageNote": "UJI CATATAN COV",
				"Currency": "IDR", "CurrencyID": "10026", "TSIPerObject": "1000000000", "TSINusare": "250000000"},
			{"ID": "UJI-IT-2", "ObjectItemName": "UJI ITEM 2", "CoverageID": "UJI-COV2", "CoverageNote": "UJI CATATAN COV2",
				"Currency": "IDR", "CurrencyID": "10026", "OccupationName": "UJI OKUPASI"},
		},
		models.DaftarDiItem(1, 1, models.AnakEstimasi): {
			{"PrintFaceClaim": "1", "EstimationValue": "1000000", "GrossEstimationPct": "4000000", "Currency": "IDR",
				"CurrencyID": "10026"},
		},
		models.DaftarDiItem(1, 2, models.AnakEstimasi): {
			{"PrintFaceClaim": "0", "EstimationValue": "500000", "Currency": "IDR", "CurrencyID": "10026"},
		},
		models.DaftarDiItem(1, 2, models.AnakSpreadPL): {
			{"TreatyType": "10003", "TreatyName": "UJI QS", "SharePercentage": "100", "TSISpreaded": "250000000"},
		},
		models.DaftarAdj(1, 1): {
			{"ID": "UJI-ADJ-11", "PaymentType": "2", "Currency": "IDR", "CurrencyID": "10026", "GrossAdjustment": "8000000",
				"AdjustmentValue": "2000000", "ValueAdjustment": "2000000", "AcceptanceStatus": "1", "AcceptedNo": "UJI-LAMA"},
		},
		models.DaftarAdj(1, 2): {
			{"ID": AdjUji, "PaymentType": "1", "Currency": "IDR", "CurrencyID": "10026", "CurrencyDol": "1",
				"GrossAdjustment": nilaiIDR, "AdjustmentValue": nilaiIDR, "ValueAdjustment": nilaiIDR, "PersenRNM": "25",
				"IndividualRiskValue": "0", "IsKomite": "1", "KomiteID": KomiteUji, "pxCreateOperator": PembuatUji,
				"pxCreateOpName": "UJI Admin", "DirectToKasir": "true", "NameOfBank": "UJI BANK", "BranchOfBank": "UJI CABANG",
				"NoAccount": "12-34", "PayableTo": "UJI PENERIMA", "DataCommitteFacin.Remarks": "UJI CATATAN KOMITE"},
		},
		models.DaftarDiAdj(1, 2, 1, models.AnakSpread): {
			{"TreatyType": "10003", "TreatyName": "UJI QS", "SharePercentage": "100", "ClaimSpreaded": nilaiIDR,
				"Currency": "IDR"},
		},
	}
	return nilai, daftar
}

// SiapkanTT2 - kasus klaim induk + kasus komite TT2 lahir dengan tangga calon tingkat 1 (Claim Supervisor = SPVA).
func (g *Gudang) SiapkanTT2(nilaiIDR string, saat time.Time) {
	n, d := KlaimFacInUji(nilaiIDR)
	g.Klaim.Setel(KlaimUji, n, d)
	g.Lahirkan(KomiteUji, KlaimUji, AdjUji, models.TransferAdjustment, PembuatUji, "UJI Admin",
		[]models.Anggota{{OperatorID: models.WorkbasketSPVA, Jabatan: "Claim Supervisor"}}, saat)
}

// SiapkanTutup - kasus klaim induk + kasus komite TT3 / TT4 satu tingkat (ReasClaimDeptHead, KCF-03) tanpa adjustment.
func (g *Gudang) SiapkanTutup(transfer string, saat time.Time) {
	n, d := KlaimFacInUji("10000000")
	g.Klaim.Setel(KlaimUji, n, d)
	g.Lahirkan(KomiteUji, KlaimUji, "", transfer, PembuatUji, "UJI Admin",
		[]models.Anggota{{OperatorID: models.WorkbasketDeptHead, Jabatan: "Claim Dept. Head"}}, saat)
}
