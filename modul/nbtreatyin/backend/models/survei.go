package models

// Untuk apa berkas ini: popup Historical Survey Report (`Harness/HistoricalSurveyReport` <- tombol
// Survey Report `Section/DetailPolicyTreatyIn` sel 22). Keputusan work owner 06-10-2026 membatalkan K7:
// daftar `PolicyTreatyIn.QuotationData.SurveyReportList` disimpan di T_POLIS_SURVEY (`TabelSurvei`).
//
//	tombol   pyVisible  `pyWorkPage.Quotation.ProportionalType != 'NonProportional'`
//	         pyDisabledWhen `.QuotationData.IsSurveyReport=='No' || .QuotationData.IsSurveyReport==''`
//	admin    `Section/HistoricalSurveyReportDtl`: Add (addRow), Delete (deleteRow), Submit
//	         (`SetSurveyReport_Act`: hanya Property-Set Quotation = PolicyTreatyIn.QuotationData, nol Obj-Save)
//	atasan   `Section/HistoricalSurveyReportDtlUW`: hanya-baca

// kolomSurvei - medan sel grid `Section/HistoricalSurveyReportDtl` (kolom TabelSurvei).
var kolomSurvei = []string{"DateofSurvey", "SurveyedBy", "LossPrevention", "Remarks"}

// SurveiDapatDisunting - tombol Survey Report tampil DAN aktif: popup dapat dibuka.
func SurveiDapatDisunting(h *Halaman) bool {
	v := h.Ambil(HalamanPolis + ".QuotationData.IsSurveyReport")
	return bukanNonProp(h) && v != "No" && v != ""
}

// barisSurvei - kiriman grid survei: hanya medan sel, baris yang seluruhnya kosong dibuang.
func barisSurvei(masuk []Baris) []Baris {
	out := []Baris{}
	for _, b := range masuk {
		r := Baris{}
		for _, k := range kolomSurvei {
			if v := b[k]; v != "" {
				r[k] = v
			}
		}
		if len(r) > 0 {
			out = append(out, r)
		}
	}
	return out
}
