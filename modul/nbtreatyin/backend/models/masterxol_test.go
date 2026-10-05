package models

// F1 (keputusan WO putaran 3, 04-10-2026): daftar medan master yang boleh
// dibaca dari dokumen JSON adalah daftar TERTUTUP. Menambah atau membuang satu
// medan berarti mengubah daftar harfiah di bawah DENGAN SADAR - bersama kutipan
// langkah XML-nya di masterxol.go. Penyaring pembaca diuji terpisah
// (`repository.TestUraiMasterXOLHanyaMedanDaftarTertutup`).

import (
	"fmt"
	"testing"
)

func TestDaftarMedanMasterXOLTertutup(t *testing.T) {
	harap := []string{
		"EDMState",
		"FacultativeShare",
		"FacultativeShareList.DeductionList.Currency", "FacultativeShareList.DeductionList.Deduction", "FacultativeShareList.DeductionTotalList.Currency", "FacultativeShareList.DeductionTotalList.Value",
		"FacultativeShareList.GrossPremiumList.Currency", "FacultativeShareList.GrossPremiumList.Value",
		"Installment.AmountTotal", "Installment.Currency", "Installment.InstallmentList.Amount", "Installment.InstallmentList.Currency",
		"Installment.InstallmentList.Installment", "Installment.InstallmentList.InstallmentPct", "Installment.InstallmentList.PaymentDate", "Installment.PctTotal",
		"LimitFacShareSummaryList.Deductible", "LimitFacShareSummaryList.Deductible2", "LimitFacShareSummaryList.Limit", "LimitFacShareSummaryList.Limit2",
		"LimitFacShareSummaryList.MDP", "LimitFacShareSummaryList.MDP2", "LimitFacShareSummaryList.NetPremi", "LimitFacShareSummaryList.NetPremi2",
		"LimitFacShareSummaryList.Note",
		"LimitShareSummaryList.Deductible", "LimitShareSummaryList.Deductible2", "LimitShareSummaryList.Limit", "LimitShareSummaryList.Limit2",
		"LimitShareSummaryList.MDP", "LimitShareSummaryList.MDP2", "LimitShareSummaryList.NetPremi", "LimitShareSummaryList.NetPremi2",
		"LimitShareSummaryList.NetPremiAfterPPH", "LimitShareSummaryList.NetPremiAfterPPH2", "LimitShareSummaryList.NetPremiAfterPPN", "LimitShareSummaryList.NetPremiAfterPPN2",
		"LimitShareSummaryList.Note",
		"LimitSummaryList.Deductible", "LimitSummaryList.Deductible2", "LimitSummaryList.Limit", "LimitSummaryList.Limit2",
		"LimitSummaryList.MDP", "LimitSummaryList.MDP2", "LimitSummaryList.Note",
		"Limits.Limit", "Limits.Limit2", "Limits.MDPList.Currency", "Limits.MDPList.Value",
		"Limits.ReinstatementNote", "Limits.ReinstatementValue", "Limits.Reinstatement_List.ReinstatementValue",
		"ProportionType",
		"RNMShare",
		"RnmShareDeducted",
		"Share.DeductionList.Currency", "Share.DeductionList.Deduction", "Share.DeductionTotalList.Currency", "Share.DeductionTotalList.Value",
		"Share.GrossPremiumList.Currency", "Share.GrossPremiumList.Value", "Share.Layer", "Share.LayerPart",
		"Share.LayerPartType", "Share.LayerType", "Share.NetPremiumList.Currency", "Share.NetPremiumList.Value",
		"Share.RnmLimitList.Currency", "Share.RnmLimitList.Value", "Share.SpreadingListXOL.Pct", "Share.SpreadingListXOL.ReinsTypeID",
		"Share.SpreadingListXOL.ReinsTypeName", "Share.SpreadingTypeIDXOL", "Share.SpreadingTypeXOL",
		"TotalFacShareDeductionNP.Currency", "TotalFacShareDeductionNP.Value",
		"TotalFacShareGrossNP.Currency", "TotalFacShareGrossNP.Value",
		"TotalFacShareNetNP.Currency", "TotalFacShareNetNP.Value",
		"TotalFacShareRnmNP.Currency", "TotalFacShareRnmNP.Value",
		"TotalLimitDeductblNP.Currency", "TotalLimitDeductblNP.Value",
		"TotalLimitIOONP.Currency", "TotalLimitIOONP.Value",
		"TotalLimitMDPNP.Currency", "TotalLimitMDPNP.Value",
		"TotalShareDeductionNP.Currency", "TotalShareDeductionNP.TotalPPHValue", "TotalShareDeductionNP.TotalPPNValue", "TotalShareDeductionNP.Value",
		"TotalShareGrossNP.Currency", "TotalShareGrossNP.Value",
		"TotalShareNetNP.Currency", "TotalShareNetNP.TotalNetPremiAfterPPN", "TotalShareNetNP.TotalNetPremiAfterTax", "TotalShareNetNP.Value",
		"TotalShareRnmNP.Currency", "TotalShareRnmNP.Value",
		"TotalSpreadedNetPremi.Currency", "TotalSpreadedNetPremi.Value",
		"TotalSpreadedNetPremiRI.Currency", "TotalSpreadedNetPremiRI.Value",
	}
	if dapat := MedanMasterXOL(); fmt.Sprint(dapat) != fmt.Sprint(harap) {
		ada := map[string]bool{}
		for _, j := range dapat {
			ada[j] = true
		}
		for _, j := range harap {
			if !ada[j] {
				t.Errorf("medan %s hilang dari daftar tertutup", j)
			}
			delete(ada, j)
		}
		for j := range ada {
			t.Errorf("medan %s di luar daftar tertutup - tambahkan dengan kutipan XML dan ubah uji ini dengan sadar", j)
		}
		t.Fatalf("%d medan, harap %d", len(dapat), len(harap))
	}
}

// Setiap medan dibaca dari halaman `TreatyIn` lewat jalur yang sama dengan
// kunci daftar: tanpa kembar, dan tanpa medan keluaran `TreatyXOLList` (F1).
func TestDaftarMedanMasterXOLTanpaKeluaranDanKembar(t *testing.T) {
	lihat := map[string]bool{}
	for _, j := range MedanMasterXOL() {
		if lihat[j] {
			t.Errorf("medan %s kembar", j)
		}
		lihat[j] = true
		if len(j) >= len("TreatyXOLList") && j[:len("TreatyXOLList")] == "TreatyXOLList" {
			t.Errorf("%s adalah keluaran, bukan medan master (F1)", j)
		}
	}
}
