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
//	       radio Approval dan tombol Submit (SetDueTo_act). Grid S45
//	       `.ListInstallment` ber-pyEditingMode/pyRowEditing `readOnly` (kedua layar):
//	       sel ber-aksi SetValidateInstallment_Act / CountPctInstallment_Act di
//	       dalamnya tidak pernah terpicu.
//	Atasan `Section/DetailDeptHeadTreatyIn_UW`: SELURUH sel ber-refresh
//	       `pyReadOnly` (CountNetPremi_act, CountRiCommOgp_act, CountRiCommOnp_act,
//	       CountOverridingCommOgp_Act, CountOverridingCommOnp_Act, ...) - tidak
//	       pernah terpicu. Yang terbuka: radio `ListSuggest .IsApproved`
//	       (runActivity SetDueTo_act; CekLimitTreatyAcc_Act K2; Protection_Act
//	       tidak ada di korpus). Sel %Share subsection NonProp (wadah S88) TIDAK
//	       terbuka bagi atasan - `[keputusan work owner 07-10-2026]` "HANYA ADMIN
//	       YANG BISA EDIT".
//
// `CountSpreading` hanya bila grid spreadingnya terbuka (`models.SpreadingDariLayar`,
// syarat yang sama dengan penerimaan daftarnya) - kini hanya layar admin.
//
// ⛔ Aksi di luar layar posisi berkas ditolak `ErrTindakanTakAdaDiPosisi`:
// `CheckDataMkt` menyimpan halaman (Obj-Save langkah 5) dan tidak boleh
// menjadi jalan simpan draf bagi atasan (layar atasan tanpa tombol Save).

import "nusantarare/modul/nbtreatyin/backend/models"

var aksiAdmin = map[string]bool{
	"SystemSetOneYear": true, "RemoveTypeTax": true, "ProtectDate": true, "CheckDataMkt": true,
	"SetCurrency": true, "TreatyEnableDisableInput": true, "CalculatePremi": true,
	"CountOGPONP": true, "CountResult1": true, "CountResult2Ogp": true, "CountResult1Onp": true,
	"CountResult2Onp": true, "FillPaymentInstallment": true, "SetDueTo": true, "HitungPajak": true,
}

var aksiAtasan = map[string]bool{"SetDueTo": true}

// aksiTerbuka - aksi refresh `aksi` dapat terpicu di layar posisi `posisi`
// atas halaman `h` (sesudah kiriman layar digabung).
func aksiTerbuka(posisi, aksi string, h *models.Halaman) bool {
	if aksi == "CountSpreading" {
		return models.SpreadingDariLayar(h, posisi)
	}
	if posisi == models.PosisiAdmin {
		return aksiAdmin[aksi]
	}
	return aksiAtasan[aksi]
}
