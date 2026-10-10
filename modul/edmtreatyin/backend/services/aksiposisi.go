package services

// Untuk apa berkas ini: AKSI REFRESH YANG DAPAT TERPICU PER LAYAR - hanya sel / tombol yang TERBUKA (tidak
// `pyReadOnly`, tidak `pyDisabled`, wadahnya tampil) yang menjalankan action set-nya. Pola: salinan
// `modul/nbtreatyin/backend/services/aksiposisi.go`, isinya menurut section EDM.
//
//	Semua posisi  `ListSuggestEDM` radio Approval (SetDueTo_act, Protection_Act; CekLimitTreatyAcc_Act tidak
//	              dibangun - tangga tiga jenjang)
//	Admin         + header `DetailPolicyTreatyInAddendum`: With Tax `.FlagPPH` (click -> RemoveTypeTax_ACT), Type Tax
//	              (`HitungPajak`, ketetapan NB), Marketing Officer `.QuotationData.MOID` (change -> CheckDataMkt).
//	              ⛔ XML tidak mengunci sel ini per posisi (pyEditOptions Auto); keputusan work owner 07-10-2026:
//	              atasan TIDAK boleh mengubahnya - penyimpangan sadar.
//	              + tab New Data `PropNewData2` (wadah `.IsNewPolicyNonProp != 1`): sel uang (CountOGPONP_Act,
//	              CountResult1_Act, CountResult2Ogp_act, CountResult1Onp_Act, CountResult2Onp_act + refresh
//	              PropValueDifference = EDMTCalculateTreatyDifference), grid spreading %Share (CountSpreading_Act)
//	              dan `.Installment` (FillPaymentInstallment) - keduanya + EDMTCalculateTreatyDifference karena tombol
//	              Calculate Value Difference (S24) DIBUANG (WO 08-10-2026); S12 `.Installment`
//	              NonProp baru (FillPaymentInstallmentEDMT)
//
// Grid `.ListInstallment` ber-edit mode readOnly (semua tab): CountPctInstallment_Act / SetValidateInstallment_Act di
// selnya tidak pernah terpicu. Sel tab Old Data / Value Difference / New Data atasan (`PropNewData`) seluruhnya RO -
// aksi Count* warisan salinannya tidak pernah terpicu.

import "nusantarare/modul/edmtreatyin/backend/models"

var aksiSemua = map[string]bool{"SetDueTo": true, "Protection": true}

// aksiAdminHeader - aksi sel header yang hanya terbuka bagi Admin (keputusan work owner 07-10-2026).
var aksiAdminHeader = map[string]bool{"RemoveTypeTax": true, "HitungPajak": true, "CheckDataMkt": true}

var aksiAdminTabBaru = map[string]bool{
	"CountOGPONP": true, "CountResult1": true, "CountResult2Ogp": true, "CountResult1Onp": true,
	"CountResult2Onp": true, "EDMTCalculateTreatyDifference": true, "CountSpreading": true,
	"FillPaymentInstallment": true,
}

// aksiTerbuka - aksi refresh `aksi` dapat terpicu di layar posisi `posisi` atas halaman `h`.
func aksiTerbuka(posisi, aksi string, h *models.Halaman) bool {
	if aksiSemua[aksi] {
		return true
	}
	if posisi != models.PosisiAdmin {
		return false
	}
	if aksiAdminHeader[aksi] {
		return true
	}
	if models.PolisNonPropBaru(h) {
		return aksi == "FillPaymentInstallmentEDMT"
	}
	return aksiAdminTabBaru[aksi]
}
