package services

// Untuk apa berkas ini: AKSI REFRESH YANG DAPAT TERPICU PER LAYAR - hanya sel
// yang TERBUKA (tidak `pyReadOnly`, tidak `pyDisabled`) yang menjalankan action
// set-nya. Diperiksa ulang ke korpus XML 04-10-2026 (audit silang putaran 3).
//
//	Admin  `Section/DetailPolicyTreatyIn` (+ `ListSuggest`, `SpreadingRiskList`
//	       NonProp): StartDate (pra-DT SystemSetOneYear_DT), FlagPPH
//	       (RemoveTypeTax_ACT), EndDate (ProtectDate), QuotationData.MOID
//	       (CheckDataMkt), IDCurrency (SetCurrency_act), tombol Enable / Disable
//	       Input Type (TreatyEnableDisableInput), GrossPremium/GrossClaim
//	       (CalculatePremi_Act), sel uang (CountOGPONP_Act, CountResult1_Act,
//	       CountResult2Ogp_act, CountResult1Onp_Act, CountResult2Onp_act), grid
//	       spreading (CountSpreading_Act), Installment (FillPaymentInstallment),
//	       grid angsuran (SetValidateInstallment_Act, CountPctInstallment_Act),
//	       radio Approval dan tombol Submit (SetDueTo_act).
//	Atasan `Section/DetailDeptHeadTreatyIn_UW`: SELURUH sel ber-refresh
//	       `pyReadOnly` (CountNetPremi_act, CountRiCommOgp_act, CountRiCommOnp_act,
//	       CountOverridingCommOgp_Act, CountOverridingCommOnp_Act, ...) - tidak
//	       pernah terpicu. Yang terbuka hanya radio `ListSuggest .IsApproved`
//	       (runActivity SetDueTo_act; CekLimitTreatyAcc_Act K2; Protection_Act
//	       tidak ada di korpus).
//
// ⛔ Aksi di luar layar posisi berkas ditolak `ErrTindakanTakAdaDiPosisi`:
// `CheckDataMkt` menyimpan halaman (Obj-Save langkah 5) dan tidak boleh
// menjadi jalan simpan draf bagi atasan (layar atasan tanpa tombol Save).

import "nusantarare/modul/nbtreatyin/backend/models"

var aksiAdmin = map[string]bool{
	"SystemSetOneYear": true, "RemoveTypeTax": true, "ProtectDate": true, "CheckDataMkt": true,
	"SetCurrency": true, "TreatyEnableDisableInput": true, "CalculatePremi": true,
	"CountOGPONP": true, "CountResult1": true, "CountResult2Ogp": true, "CountResult1Onp": true,
	"CountResult2Onp": true, "CountSpreading": true, "FillPaymentInstallment": true,
	"SetValidateInstallment": true, "CountPctInstallment": true, "SetDueTo": true,
}

var aksiAtasan = map[string]bool{"SetDueTo": true}

// aksiTerbuka - aksi refresh `aksi` dapat terpicu di layar posisi `posisi`.
func aksiTerbuka(posisi, aksi string) bool {
	if posisi == models.PosisiAdmin {
		return aksiAdmin[aksi]
	}
	return aksiAtasan[aksi]
}
