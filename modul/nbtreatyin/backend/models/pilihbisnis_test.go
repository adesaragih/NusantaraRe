package models

// Uji audit silang putaran 3 (bab 3.2 butir 1): aktivitas refresh sel layar
// admin dan When yang dibangun tanpa uji berharapan XML. Harapan dari langkah
// XML yang dikutip per uji.

import "testing"

// Activity/ProtectDate:
//
//	1  local.Errmsg = "End Date Cannot be less than Start Date"
//	2  [@CompareDates(.StartDate,.EndDate)] -> Property-Set-Messages .EndDate = local.Errmsg
//
// @CompareDates benar bila tanggal pertama SESUDAH tanggal kedua (bunyi
// pesannya: tanggal akhir tidak boleh sebelum tanggal mulai). Sama hari = tanpa pesan.
func TestProtectDatePesanBilaAkhirSebelumMulai(t *testing.T) {
	for _, tt := range []struct {
		mulai, akhir string
		pesan        bool
	}{
		{"2026-12-31", "2026-10-01", true},
		{"2026-10-01", "2026-10-01", false},
		{"2026-10-01", "2026-12-31", false},
	} {
		h := HalamanBaru()
		h.Setel("PolicyTreatyIn.StartDate", tt.mulai)
		h.Setel("PolicyTreatyIn.EndDate", tt.akhir)
		ProtectDate(h)
		got := h.Pesan["PolicyTreatyIn.EndDate"]
		if tt.pesan && (len(got) != 1 || got[0] != "End Date Cannot be less than Start Date") {
			t.Errorf("%s..%s: harap pesan verbatim pada EndDate, dapat %v", tt.mulai, tt.akhir, h.SemuaPesan())
		}
		if !tt.pesan && h.AdaPesan() {
			t.Errorf("%s..%s: harap tanpa pesan, dapat %v", tt.mulai, tt.akhir, h.SemuaPesan())
		}
	}
}

// Activity/RemoveTypeTax_ACT langkah 2:
// [pyWorkPage.PolicyTreatyIn.FlagPPH==false] -> Property-Remove .TypeTax.
// Properti TrueFalse kosong = false (sel `.FlagPPH` pxCheckbox).
func TestRemoveTypeTaxHanyaBilaFlagPPHTidakBenar(t *testing.T) {
	for flag, tetap := range map[string]bool{"true": true, "false": false, "": false} {
		h := HalamanBaru()
		h.Setel("PolicyTreatyIn.FlagPPH", flag)
		h.Setel("PolicyTreatyIn.TypeTax", TypeTaxInclusive)
		RemoveTypeTax(h)
		if got := h.Ambil("PolicyTreatyIn.TypeTax"); (got == TypeTaxInclusive) != tetap {
			t.Errorf("FlagPPH %q: TypeTax %q, harap tetap=%v", flag, got, tetap)
		}
	}
}

// Activity/CheckDataMkt:
//
//	2  .MarketingOfficer = ""
//	3  [MOID==""] -> lewati Obj-Browse (DataMarketing tanpa pxResults)
//	4  pyWorkPage.Quotation.{MOID,MarketingCode,MarketingName,TeamGroup,BranchCode,BranchName}
//	   = pxResults(1).{ID,ClientID,ClientName,TeamGroup,BranchDetailID,BranchDetailName};
//	   .QuotationData.{MOID,MarketingCode,MarketingName,TeamGroup} sama; .MarketingOfficer = ClientName
func TestTerapkanMOMenurutCheckDataMkt(t *testing.T) {
	h := HalamanBaru()
	h.Setel("PolicyTreatyIn.MarketingOfficer", "UJI-LAMA")
	TerapkanMO(h, BarisMO{ID: "UJI-MO-1", ClientID: "UJI-K1", ClientName: "UJI-NAMA", TeamGroup: "UJI-T",
		BranchDetailID: "UJI-C1", BranchDetailName: "UJI-CABANG"})
	harap := map[string]string{
		"Quotation.MOID": "UJI-MO-1", "Quotation.MarketingCode": "UJI-K1", "Quotation.MarketingName": "UJI-NAMA",
		"Quotation.TeamGroup": "UJI-T", "Quotation.BranchCode": "UJI-C1", "Quotation.BranchName": "UJI-CABANG",
		"PolicyTreatyIn.QuotationData.MOID": "UJI-MO-1", "PolicyTreatyIn.QuotationData.MarketingCode": "UJI-K1",
		"PolicyTreatyIn.QuotationData.MarketingName": "UJI-NAMA", "PolicyTreatyIn.QuotationData.TeamGroup": "UJI-T",
		"PolicyTreatyIn.MarketingOfficer": "UJI-NAMA",
	}
	for j, w := range harap {
		if got := h.Ambil(j); got != w {
			t.Errorf("%s = %q, harap %q", j, got, w)
		}
	}
	// MOID kosong: langkah 3 dilewati, langkah 4 menulis nilai kosong
	TerapkanMO(h, BarisMO{})
	for j := range harap {
		if got := h.Ambil(j); got != "" {
			t.Errorf("MOID kosong: %s = %q, harap kosong", j, got)
		}
	}
}

// When/TreatyMasterInEDM: `pyWorkPage.TreatyIn.EDMState` = "1" OR "2" OR "3".
func TestTreatyMasterInEDMTigaNilai(t *testing.T) {
	for v, harap := range map[string]bool{"1": true, "2": true, "3": true, "": false, "0": false, "4": false} {
		h := HalamanBaru()
		h.Setel("TreatyIn.EDMState", v)
		if got := TreatyMasterInEDM(h); got != harap {
			t.Errorf("EDMState %q: %v, harap %v", v, got, harap)
		}
	}
}
