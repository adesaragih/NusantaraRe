package repository

// Daftar kolom penyalinan versi - dibandingkan dengan DDL PremiumList 052-054
// oleh `TestKolomSalinanSamaDenganDDL`, supaya kolom yang ditambahkan
// PremiumList tidak diam-diam tertinggal dari salinan endorsement.

// kolomIdentitasPeserta - kolom `T_PREMIUM_LIST_DETAIL` yang DIRAKIT, bukan
// disalin: identitas baru, induk versi, pengenal kasus, nomor, dan status EDM.
var kolomIdentitasPeserta = []string{
	"ID", "PREMIUM_LIST_ID", "PARENT_ID", "ID_PEGA", "PL_NUMBER",
	"PL_NUMBER_EDM", "EDM_STATUS", "STATUS_OLD", "STATUS",
}

// kolomNilaiPeserta - kolom nilai `T_PREMIUM_LIST_DETAIL` (052), urut DDL,
// disalin apa adanya dari versi lama (`MappingEDMLife` 10 b2339).
var kolomNilaiPeserta = []string{
	"POLICY_NO", "POLICY_HOLDER", "CERTIFICATE_NO", "NAME_OF_INSURED", "DESCRIPTION", "SEX", "DOB",
	"AGE", "ENTRY_AGE", "CURRENT_AGE", "PLAN", "RISK", "MEDICAL_STATUS", "STNC", "WPC", "CURRENCY",
	"PERIOD_YY", "PERIOD_MM", "PASSED_PERIOD", "BEGIN_DATE", "EFFECTIVE_DATE", "EXPIRED_DATE", "LAPSE_DATE",
	"GROSS_VALUATION_BEGIN_DATE", "GROSS_VALUATION_EXPIRED_DATE", "RETRO_VALUATION_BEGIN_DATE",
	"RETRO_VALUATION_EXPIRED_DATE", "SUM_INSURED", "SUM_REASURED", "SUM_AT_RISK_GROSS", "SUM_AT_RISK_RETRO",
	"CEDING_RETENTION", "CEDING_CO", "SHARE_NUSANTARA_RE", "SHARE_NUSANTARA_RE_GROSS", "SHARE_RETRO",
	"RETROCEDED_SHARE", "RATE", "FACTOR", "EM_PERCENT", "PRO_RATE_TYPE", "GROSS_PREMIUM", "NET_PREMIUM",
	"COMM", "PROF_COMM", "OVR_COMM", "BROKERAGE_FEE", "TAX", "FLEET_DISCOUNT", "DEDUCTION", "RI_ADMIN_FEE",
	"CLAIM", "CLAIM_AMOUNT", "GROSS_PREMIUM_REFUND", "NET_PREMIUM_REFUND", "COMM_REFUND", "OVR_COMM_REFUND",
	"BROKERAGE_FEE_REFUND", "TAX_REFUND", "DEDUCTION_REFUND", "RI_ADMIN_FEE_REFUND", "GROSS_PREMIUM_RETRO",
	"NET_PREMIUM_RETRO", "DISCOUNT_PREMIUM_RETRO", "OVR_COMM_RETRO", "BROKERAGE_FEE_RETRO",
	"RI_ADMIN_FEE_RETRO", "GROSS_PREMIUM_REFUND_RETRO", "NET_PREMIUM_REFUND_RETRO",
	"DISCOUNT_PREMIUM_REFUND_RETRO", "OVR_COMM_REFUND_RETRO", "BROKERAGE_FEE_REFUND_RETRO",
	"RI_ADMIN_FEE_REFUND_RETRO",
}

// kolomNilaiSpreading - `T_PREMIUM_LIST_SPREADING` (053) selain `ID`, `DETAIL_ID`.
var kolomNilaiSpreading = []string{
	"TREATY_TYPE_ID", "TREATY_TYPE_NAME", "TREATY_YEAR_LIFE", "RETROCADED_SHARE", "IDR", "USD",
	"IDR_SELISIH", "USD_SELISIH", "B_IDR", "B_USD", "TGL_UPDATE", "USER_ID",
}

// kolomNilaiSpreadingRetro - `T_PREMIUM_LIST_SPREADING_RETRO` (054) selain
// `ID`, `SPREADING_ID`.
var kolomNilaiSpreadingRetro = []string{
	"REINSURER_NAME", "PERCENT_SHARE", "AMOUNT", "COMMISION", "OVR_COMM", "RATE",
	"PREMIUM_SPREADED_GROSS", "PREMIUM_SPREADED_NET", "TREATY_TYPE_ID", "TREATY_TYPE_NAME",
	"TREATY_START_DATE", "TREATY_END_DATE", "TGL_UPDATE", "USER_ID",
}

// sumberWarisanPeserta - ekspresi kolom tabel peserta warisan (alias `m`)
// untuk kolom nilai yang namanya ATAU bentuknya berbeda di tabel kita.
//
// `[terverifikasi]` pemetaan kebalikan `SaveMasterLPDet` (`RDBList/SaveMasterLPDet.xml`
// b87, daftar `INSERT`): `PRORATETYPE` ← `CARI34 = ProRateType`; `STNC`/`WPC`
// bertipe DATE di tabel warisan (`[data DBA]`, katalog Claim Life) dan TEKS
// `dd/mm/yyyy` di tabel kita (052); `RISK` dikirim `@toDecimal(.RISK)` ke kolom
// warisan, teks di tabel kita. Kolom lain bernama dan berbentuk sama.
var sumberWarisanPeserta = map[string]string{
	"PRO_RATE_TYPE": "m.PRORATETYPE",
	"STNC":          "TO_CHAR(m.STNC, 'DD/MM/YYYY')",
	"WPC":           "TO_CHAR(m.WPC, 'DD/MM/YYYY')",
	"RISK":          "TO_CHAR(m.RISK)",
}

// ekspresiWarisan - ekspresi sumber warisan untuk satu kolom nilai kita.
func ekspresiWarisan(kolom string) string {
	if e, ada := sumberWarisanPeserta[kolom]; ada {
		return e
	}
	return "m." + kolom
}
