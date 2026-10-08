package repository

// Peta larik JSON -> tabel pendaratan. SATU sumber kebenaran.
//
// ⛔ Peta ini yang membangkitkan SQL sisip, SQL hapus, dan rekonsiliasi.
// Ketiganya pernah ditulis terpisah di tempat lain, dan yang ketiga selalu
// yang paling akhir diperbarui - sehingga rekonsiliasi membuktikan bentuk
// yang sudah tidak dipakai. Di sini ketiganya tidak dapat berselisih.
//
// Urutan kunci di bawah = urutan kolom di DDL migrasi 430, dan
// `TestPetaPendaratanCocokDenganDDL` mengadu keduanya.
//
// Ukuran yang menentukan isi peta ini: sapuan 3 Oktober 2026 atas SELURUH
// 1.854 dokumen `POOLDATA.M_TREATY_IN.JSONDATA`, diurai utuh sebagai JSON.

// Pendaratan adalah satu larik JSON yang mendarat ke satu tabel.
type Pendaratan struct {
	// Kunci larik di puncak dokumen - kosong untuk tabel anak.
	Larik string
	// Nama tabel INDUK - kosong untuk tabel tingkat pertama.
	//
	// ⭐ Sejak migrasi 437 pohonnya TIGA tingkat, dan penanganan anak tidak
	// lagi ditanam untuk satu kasus: `Induk` + `KunciAnak` memerikan
	// hubungannya sebagai DATA, dan pemuatnya menapakinya secara umum.
	Induk string
	// Kunci larik DI DALAM elemen induk. Hanya berarti bila `Induk` terisi.
	KunciAnak string
	// ⭐ BEBERAPA larik ke SATU tabel — migrasi 438 dan 439.
	//
	// xlsx menggabungkan larik yang BENTUKNYA sama ke satu kotak:
	// `IOOLimitList | RetentionList | CessionList | EPIList` semuanya
	// `{Currency, Value, ...}` dan semuanya mendarat di
	// `T_TREATY_LIMIT_AMOUNT`. Yang membedakan barisnya sesudah mendarat
	// adalah kolom `JENIS`, yang diisi NAMA LARIK ASALNYA.
	//
	// ⛔ `JENIS` TIDAK ADA di dokumen. Ia dipaksa oleh bentuk xlsx, dan
	// tanpanya baris retensi tidak dapat dibedakan dari baris cession.
	//
	// Berlaku di kedua tingkat: dengan `Induk` kosong ia membaca larik AKAR
	// dokumen, dengan `Induk` terisi ia membaca larik di dalam elemen induk.
	// Tidak boleh dipakai bersama `Larik` maupun `KunciAnak`.
	LarikGabung []string
	// ⭐ SKALAR AKAR dokumen sebagai SATU baris — `T_TREATY_REVISION`.
	//
	// Dokumennya sendiri yang menjadi "elemen", dan `URUTAN` selalu 0.
	// Tidak boleh dipakai bersama medan larik mana pun.
	Akar  bool
	Tabel string
	Seq   string
	// Kunci JSON elemen, berpasangan urut dengan Kolom.
	Kunci []string
	Kolom []string
	// Cacah baris yang sapuan 3 Oktober 2026 temukan di seluruh korpus.
	// Dipakai rekonsiliasi penuh, bukan rekonsiliasi per kontrak.
	CacahTerukur int
}

// PetaPendaratan adalah kedelapan tabel, dalam urutan MUAT: induk sebelum
// anaknya. Urutan HAPUS adalah kebalikannya.
//
// ⭐ Tabel ANAK tidak punya `Larik` di puncak dokumen; ia memerikan
// induknya lewat `Induk` dan kunci lariknya lewat `KunciAnak`. Pemuatnya
// menapaki hubungan itu secara umum — nol cabang kode per pasangan.
var PetaPendaratan = []Pendaratan{
	{
		Larik: "ReportingPeriodList", Tabel: "T_TREATY_REPORTING_PERIOD", Seq: "SEQ_TT_REPORTING_PERIOD",
		Kunci:        []string{"AutoCalculate", "ConfirmationDue", "InitialDate", "Period", "SettlementDue", "SubmissionDue"},
		Kolom:        []string{"AUTOCALCULATE", "CONFIRMATIONDUE", "INITIALDATE", "PERIOD", "SETTLEMENTDUE", "SUBMISSIONDUE"},
		CacahTerukur: 4548,
	},
	{
		Larik: "Portfolio", Tabel: "T_TREATY_PORTFOLIO", Seq: "SEQ_TT_PORTFOLIO",
		Kunci:        []string{"Description", "Type", "TypePortfolio"},
		Kolom:        []string{"DESCRIPTION", "TYPE", "TYPEPORTFOLIO"},
		CacahTerukur: 1925,
	},
	{
		Larik: "AccumulationList", Tabel: "T_TREATY_ACCUMULATION", Seq: "SEQ_TT_ACCUMULATION",
		Kunci:        []string{"Period", "ReportDate", "SubDays", "SubDueDate"},
		Kolom:        []string{"PERIOD", "REPORTDATE", "SUBDAYS", "SUBDUEDATE"},
		CacahTerukur: 60,
	},
	{
		Larik: "EGNPI", Tabel: "T_TREATY_EGNPI", Seq: "SEQ_TT_EGNPI",
		Kunci:        []string{"Amount", "AmountIDR", "AsDate", "ClassOfBusiness", "Currency", "CurrencyID", "Note", "Proportion", "TreatyGroup", "TreatyGroupID", "pyTemplateRichTextEditor"},
		Kolom:        []string{"AMOUNT", "AMOUNTIDR", "ASDATE", "CLASSOFBUSINESS", "CURRENCY", "CURRENCYID", "NOTE", "PROPORTION", "TREATYGROUP", "TREATYGROUPID", "PYTEMPLATERICHTEXTEDITOR"},
		CacahTerukur: 2298,
	},
	{
		Larik: "Retention", Tabel: "T_TREATY_RETENTION", Seq: "SEQ_TT_RETENTION",
		Kunci:        []string{"Amount", "ClassOfBusiness", "Currency", "CurrencyID", "Note", "TreatyGroup", "TreatyGroupID"},
		Kolom:        []string{"AMOUNT", "CLASSOFBUSINESS", "CURRENCY", "CURRENCYID", "NOTE", "TREATYGROUP", "TREATYGROUPID"},
		CacahTerukur: 2511,
	},
	{
		Larik: "Installment", Tabel: "T_TREATY_INSTALLMENT", Seq: "SEQ_TT_INSTALLMENT",
		Kunci:        []string{"AmountTotal", "Currency", "PctTotal", "pxListSubscript"},
		Kolom:        []string{"AMOUNTTOTAL", "CURRENCY", "PCTTOTAL", "PXLISTSUBSCRIPT"},
		CacahTerukur: 796,
	},
	{
		// Anak `Installment`; lihat catatan di atas.
		Induk: "T_TREATY_INSTALLMENT", KunciAnak: LarikAnakAngsuran,
		Tabel: "T_TREATY_INSTALLMENT_ITEM", Seq: "SEQ_TT_INSTALLMENT_ITEM",
		Kunci:        []string{"Amount", "Currency", "DueDate", "Installment", "InstallmentPct", "PaymentDate", "WPC"},
		Kolom:        []string{"AMOUNT", "CURRENCY", "DUEDATE", "INSTALLMENT", "INSTALLMENTPCT", "PAYMENTDATE", "WPC"},
		CacahTerukur: 3033,
	},
	{
		// Tab Co-Ins Scale - migrasi 432, 3 Oktober 2026.
		//
		// ⚠️ ENAM medan skalar, bukan dua. Rancangannya menyebut `CoInShare`
		// dan `PctLimit` saja; sapuan menemukan keempat medan jejak Pega ada
		// pada 701 dari 702 elemen, dan ia satu-satunya catatan siapa yang
		// menyusun skala itu dan kapan.
		Larik: "CoInScale", Tabel: "M_TREATYIN_COINSCALE", Seq: "SEQ_MTI_COINSCALE",
		Kunci:        []string{"CoInShare", "PctLimit", "pxCreateDateTime", "pxCreateOpName", "pxCreateOperator", "pxCreateSystemID"},
		Kolom:        []string{"COINSHARE", "PCTLIMIT", "PXCREATEDATETIME", "PXCREATEOPNAME", "PXCREATEOPERATOR", "PXCREATESYSTEMID"},
		CacahTerukur: 702,
	},
	{
		// ⚠️ `Date` -> `TANGGAL`: `DATE` kata cadangan Oracle (ORA-00923).
		Larik: "CommentList", Tabel: "T_VIEW_COMMENT", Seq: "SEQ_TV_COMMENT",
		Kunci:        []string{"Date", "ConvertDate", "HasHistory", "IsApproved", "OperatorName", "Suggest"},
		Kolom:        []string{"TANGGAL", "CONVERTDATE", "HASHISTORY", "ISAPPROVED", "OPERATORNAME", "SUGGEST"},
		CacahTerukur: 11365,
	},

	// ⭐ TIGA BELAS anak `TREATY_IN` dari migrasi 437 — pohon menurut
	// `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx`, keputusan pemilik proses
	// 6 Oktober 2026. Urutannya mengikat: induk mendahului anaknya.
	//
	// Kunci dan kolomnya DITURUNKAN dari sapuan 1.854 dokumen, bukan daftar
	// tangan — xlsx menyatakan "KOLOM SENGAJA TIDAK DIMUAT".
	//
	// ⚠ Kunci dokumen `ID` menjadi kolom `IDX`: `ID` sudah dipakai kunci
	// utama tabelnya sendiri. Begitu pula `Comment` → `COMMENT_`, kata
	// cadangan Oracle.
	{
		Larik: "Limits",
		Tabel: "T_TREATY_LIMITS", Seq: "SEQ_TT_LIMITS",
		Kunci:        []string{"AdjRate", "AgregateLimit", "AgregateLimit2", "Cover", "Currency", "Currency2", "CurrencyID", "CurrencyRelation", "Deductible", "Deductible2", "ID", "IsCombineMDP", "Layer", "LayerPart", "LayerPartType", "LayerType", "Limit", "Limit2", "MDPMinPct", "MDPPct", "NoRIPCalculation", "ROLPct", "ReinstatementNote", "ReinstatementPct", "ReinstatementValue", "SumLossRatio", "SumTotalAchievIncured", "SumTotalAchievNetPremium", "TreatyType", "TreatyTypeID"},
		Kolom:        []string{"ADJRATE", "AGREGATELIMIT", "AGREGATELIMIT2", "COVER", "CURRENCY", "CURRENCY2", "CURRENCYID", "CURRENCYRELATION", "DEDUCTIBLE", "DEDUCTIBLE2", "IDX", "ISCOMBINEMDP", "LAYER", "LAYERPART", "LAYERPARTTYPE", "LAYERTYPE", "LIMIT", "LIMIT2", "MDPMINPCT", "MDPPCT", "NORIPCALCULATION", "ROLPCT", "REINSTATEMENTNOTE", "REINSTATEMENTPCT", "REINSTATEMENTVALUE", "SUMLOSSRATIO", "SUMTOTALACHIEVINCURED", "SUMTOTALACHIEVNETPREMIUM", "TREATYTYPE", "TREATYTYPEID"},
		CacahTerukur: 4210,
	},
	{
		Larik: "Share",
		Tabel: "T_TREATY_SHARE", Seq: "SEQ_TT_SHARE",
		Kunci:        []string{"Cover", "Layer", "LayerPart", "LayerPartType", "LayerType", "RNMShare", "SpreadingTotalPctXOL", "SpreadingTypeIDXOL", "SpreadingTypeXOL"},
		Kolom:        []string{"COVER", "LAYER", "LAYERPART", "LAYERPARTTYPE", "LAYERTYPE", "RNMSHARE", "SPREADINGTOTALPCTXOL", "SPREADINGTYPEIDXOL", "SPREADINGTYPEXOL"},
		CacahTerukur: 2923,
	},
	{
		Larik: "ShareReins",
		Tabel: "T_TREATY_RETRO_SHARE", Seq: "SEQ_TT_RETRO_SHARE",
		Kunci:        []string{"Layer", "ReinsID", "ReinsName", "SharePct"},
		Kolom:        []string{"LAYER", "REINSID", "REINSNAME", "SHAREPCT"},
		CacahTerukur: 9,
	},
	{
		Larik: "FacultativeShareList",
		Tabel: "T_TREATY_FAC_SHARE", Seq: "SEQ_TT_FAC_SHARE",
		Kunci:        []string{"Cover", "Layer", "LayerPart", "LayerPartType", "LayerType", "Limit", "Limit2", "SpreadingTotalPctXOL", "SpreadingTypeIDXOL", "SpreadingTypeXOL", "SpreadingTypeXOLRetro", "SpreadingTypeXOLRetroID"},
		Kolom:        []string{"COVER", "LAYER", "LAYERPART", "LAYERPARTTYPE", "LAYERTYPE", "LIMIT", "LIMIT2", "SPREADINGTOTALPCTXOL", "SPREADINGTYPEIDXOL", "SPREADINGTYPEXOL", "SPREADINGTYPEXOLRETRO", "SPREADINGTYPEXOLRETROID"},
		CacahTerukur: 12,
	},
	{
		Larik: "ShareFacultativeReinsurers",
		Tabel: "T_TREATY_FAC_REINSURER", Seq: "SEQ_TT_FAC_REINSURER",
		Kunci:        []string{"Layer", "ReinsID", "ReinsName", "SharePct"},
		Kolom:        []string{"LAYER", "REINSID", "REINSNAME", "SHAREPCT"},
		CacahTerukur: 15,
	},
	{
		Induk: "T_TREATY_LIMITS", KunciAnak: "Detail",
		Tabel: "T_TREATY_LIMIT_DETAIL", Seq: "SEQ_TT_LIMIT_DETAIL",
		Kunci:        []string{"AchievementPct", "Brokerage", "CashLoss", "CessionPct", "ClaimCooperation", "ClassOfBusiness", "CurrencyCashLoss", "CurrencyClaimCooperation", "CurrencyEPI", "CurrencyEarthquake", "CurrencyFloodJab", "CurrencyFloodNat", "CurrencyID", "CurrencyPLA", "CurrencyRSMD", "EPI", "Earthquake", "FloodJab", "FloodNation", "ID", "IOOPct", "LossRatio", "LowerBand", "PLA", "ParentID", "Periode", "PremiumReservePct", "ProfitCommision", "ProfitME", "ProfitYDCF", "QSPct", "RIOGR", "RIONR", "RNMShare", "RSMDLimit", "ReisuredParticipant", "RetentionPct", "ShareNote", "SpreadingTotalPct", "SpreadingType", "SpreadingTypeID", "SumLossRatio", "SumTotalAchievIncured", "SumTotalAchievNetPremium", "Surplus", "TotalAchAfterClaim", "TotalAchCASHCALL", "TotalAchEstCASHCALL", "TotalAchIncured", "TotalAchNetPremium", "TotalAchOsCASHCALL", "TotalAchOsClaim", "TotalAchOuts", "TotalAchPaid", "TotalAchPremium", "TreatyGroup", "TreatyGroupID", "TreatyType", "UpperBand", "Currency", "NusaReLimit"},
		Kolom:        []string{"ACHIEVEMENTPCT", "BROKERAGE", "CASHLOSS", "CESSIONPCT", "CLAIMCOOPERATION", "CLASSOFBUSINESS", "CURRENCYCASHLOSS", "CURRENCYCLAIMCOOPERATION", "CURRENCYEPI", "CURRENCYEARTHQUAKE", "CURRENCYFLOODJAB", "CURRENCYFLOODNAT", "CURRENCYID", "CURRENCYPLA", "CURRENCYRSMD", "EPI", "EARTHQUAKE", "FLOODJAB", "FLOODNATION", "IDX", "IOOPCT", "LOSSRATIO", "LOWERBAND", "PLA", "PARENTID", "PERIODE", "PREMIUMRESERVEPCT", "PROFITCOMMISION", "PROFITME", "PROFITYDCF", "QSPCT", "RIOGR", "RIONR", "RNMSHARE", "RSMDLIMIT", "REISUREDPARTICIPANT", "RETENTIONPCT", "SHARENOTE", "SPREADINGTOTALPCT", "SPREADINGTYPE", "SPREADINGTYPEID", "SUMLOSSRATIO", "SUMTOTALACHIEVINCURED", "SUMTOTALACHIEVNETPREMIUM", "SURPLUS", "TOTALACHAFTERCLAIM", "TOTALACHCASHCALL", "TOTALACHESTCASHCALL", "TOTALACHINCURED", "TOTALACHNETPREMIUM", "TOTALACHOSCASHCALL", "TOTALACHOSCLAIM", "TOTALACHOUTS", "TOTALACHPAID", "TOTALACHPREMIUM", "TREATYGROUP", "TREATYGROUPID", "TREATYTYPE", "UPPERBAND", "CURRENCY", "NUSARELIMIT"},
		CacahTerukur: 2868,
	},
	{
		Induk: "T_TREATY_LIMITS", KunciAnak: "TreatyGroupList",
		Tabel: "T_TREATY_LIMIT_GROUP", Seq: "SEQ_TT_LIMIT_GROUP",
		Kunci:        []string{"ClassOfBusiness", "ClassOfBusinessID", "IsROLProfile", "TreatyGroup", "TreatyGroupID"},
		Kolom:        []string{"CLASSOFBUSINESS", "CLASSOFBUSINESSID", "ISROLPROFILE", "TREATYGROUP", "TREATYGROUPID"},
		CacahTerukur: 8433,
	},
	{
		Induk: "T_TREATY_SHARE", KunciAnak: "SpreadingListXOL",
		Tabel: "T_TREATY_SHARE_SPREADING", Seq: "SEQ_TT_SHARE_SPREADING",
		Kunci:        []string{"ParentReinsTypeID", "Pct", "ReinsTypeID", "ReinsTypeName", "Rp", "Usd"},
		Kolom:        []string{"PARENTREINSTYPEID", "PCT", "REINSTYPEID", "REINSTYPENAME", "RP", "USD"},
		CacahTerukur: 5550,
	},
	{
		Induk: "T_TREATY_SHARE", KunciAnak: "DeductionList",
		Tabel: "T_TREATY_SHARE_DEDUCTION", Seq: "SEQ_TT_SHARE_DEDUCTION",
		Kunci:        []string{"Comment", "Currency", "Deduction", "DeductionPct", "DeductionPctCalculate"},
		Kolom:        []string{"COMMENT_", "CURRENCY", "DEDUCTION", "DEDUCTIONPCT", "DEDUCTIONPCTCALCULATE"},
		CacahTerukur: 2419,
	},
	{
		Induk: "T_TREATY_FAC_SHARE", KunciAnak: "DeductionList",
		Tabel: "T_TREATY_FAC_SHARE_DEDUCTION", Seq: "SEQ_TT_FAC_SHARE_DED",
		Kunci:        []string{"Comment", "Currency", "Deduction", "DeductionPct", "DeductionPctCalculate"},
		Kolom:        []string{"COMMENT_", "CURRENCY", "DEDUCTION", "DEDUCTIONPCT", "DEDUCTIONPCTCALCULATE"},
		CacahTerukur: 12,
	},
	{
		Induk: "T_TREATY_LIMIT_DETAIL", KunciAnak: "COBList",
		Tabel: "T_TREATY_LIMIT_COB", Seq: "SEQ_TT_LIMIT_COB",
		Kunci:        []string{"ClassOfBusiness", "ClassOfBusinessID", "TreatyGroup", "TreatyGroupID"},
		Kolom:        []string{"CLASSOFBUSINESS", "CLASSOFBUSINESSID", "TREATYGROUP", "TREATYGROUPID"},
		CacahTerukur: 6547,
	},
	{
		Induk: "T_TREATY_LIMIT_DETAIL", KunciAnak: "AchievementLists",
		Tabel: "T_TREATY_LIMIT_ACHIEVEMENT", Seq: "SEQ_TT_LIMIT_ACHIEVE",
		Kunci:        []string{"BROKERAGE", "CASHCALL", "CASHCALLtoIDR", "Conversion", "Currency", "CurrencyID", "EstimationCASHCALL", "EstimationCASHCALLtoIDR", "IDPEGA", "IncuredClaim", "IncuredClaimtoIDR", "LossRatio", "NETPREMIUM", "NETPREMIUMtoIDR", "NOOFFER", "NOPOLIS", "OutstandingCASHCALL", "OutstandingCASHCALLtoIDR", "OutstandingClaim", "OutstandingClaimtoIDR", "PREMIUM", "PREMIUMtoIDR", "PaidClaim", "PaidClaimtoIDR", "Period", "QUARTERYEAR", "Quarter", "RICOMM", "SOBNAME", "TREATYGROUPNAME", "TREATYTYPE", "Total", "TotalAchievCASHCALL", "TotalAchievEstCASHCALL", "TotalAchievIncured", "TotalAchievNetPremium", "TotalAchievOsCASHCALL", "TotalAchievOsClaim", "TotalAchievOuts", "TotalAchievPaid", "TotalAchievPremium", "TotalAfterClaim", "TotaltoIDR", "Value"},
		Kolom:        []string{"BROKERAGE", "CASHCALL", "CASHCALLTOIDR", "CONVERSION", "CURRENCY", "CURRENCYID", "ESTIMATIONCASHCALL", "ESTIMATIONCASHCALLTOIDR", "IDPEGA", "INCUREDCLAIM", "INCUREDCLAIMTOIDR", "LOSSRATIO", "NETPREMIUM", "NETPREMIUMTOIDR", "NOOFFER", "NOPOLIS", "OUTSTANDINGCASHCALL", "OUTSTANDINGCASHCALLTOIDR", "OUTSTANDINGCLAIM", "OUTSTANDINGCLAIMTOIDR", "PREMIUM", "PREMIUMTOIDR", "PAIDCLAIM", "PAIDCLAIMTOIDR", "PERIOD", "QUARTERYEAR", "QUARTER", "RICOMM", "SOBNAME", "TREATYGROUPNAME", "TREATYTYPE", "TOTAL", "TOTALACHIEVCASHCALL", "TOTALACHIEVESTCASHCALL", "TOTALACHIEVINCURED", "TOTALACHIEVNETPREMIUM", "TOTALACHIEVOSCASHCALL", "TOTALACHIEVOSCLAIM", "TOTALACHIEVOUTS", "TOTALACHIEVPAID", "TOTALACHIEVPREMIUM", "TOTALAFTERCLAIM", "TOTALTOIDR", "VALUE"},
		CacahTerukur: 2031,
	},
	{
		// ⭐ 449 (8 Oktober 2026) — `Detail.SpreadingList` tab Share Prop:
		// anak susunan treaty `PROPORTIONALARRG` (`FetchQSfromMaster` [9]).
		Induk: "T_TREATY_LIMIT_DETAIL", KunciAnak: "SpreadingList",
		Tabel: "T_TREATY_LIMIT_SPREADING", Seq: "SEQ_TT_LIMIT_SPREADING",
		Kunci: []string{"ParentReinsTypeID", "Pct", "ReinsTypeID", "ReinsTypeName", "Rp", "Usd", "Value"},
		Kolom: []string{"PARENTREINSTYPEID", "PCT", "REINSTYPEID", "REINSTYPENAME", "RP", "USD", "VALUE"},
	},
	{
		// ⭐ 450 (8 Oktober 2026) — grid Deduction rincian Limits Prop.
		Induk: "T_TREATY_LIMIT_DETAIL", KunciAnak: "DeductionList",
		Tabel: "T_TREATY_LIMIT_DEDUCTION", Seq: "SEQ_TT_LIMIT_DEDUCTION",
		Kunci: []string{"Comment", "Currency", "CurrencyID", "Deduction", "DeductionPct", "DeductionPctCalculate"},
		Kolom: []string{"COMMENT_", "CURRENCY", "CURRENCYID", "DEDUCTION", "DEDUCTIONPCT", "DEDUCTIONPCTCALCULATE"},
	},
	{
		// ⭐ 450 — grid Parameter tab Achievement (`GetAchievement`).
		Induk: "T_TREATY_LIMIT_DETAIL", KunciAnak: "CurrencyList",
		Tabel: "T_TREATY_LIMIT_ACH_PARAM", Seq: "SEQ_TT_LIMIT_ACH_PARAM",
		Kunci: []string{"Parameter", "AchievementPctGross", "LossRatioGross"},
		Kolom: []string{"PARAMETER", "ACHIEVEMENTPCTGROSS", "LOSSRATIOGROSS"},
	},
	{
		// ⭐ 450 — grid Reinstatement layer Non-Prop.
		Induk: "T_TREATY_LIMITS", KunciAnak: "Reinstatement_List",
		Tabel: "T_TREATY_LIMIT_REINSTATEMENT", Seq: "SEQ_TT_LIMIT_REINST",
		Kunci: []string{"AdditionalAmount1", "AdditionalAmount2", "AdditionalPct", "ID", "ReinstatementAmount1", "ReinstatementAmount2", "ReinstatementNote", "ReinstatementPct", "ReinstatementValue"},
		Kolom: []string{"ADDITIONALAMOUNT1", "ADDITIONALAMOUNT2", "ADDITIONALPCT", "IDX", "REINSTATEMENTAMOUNT1", "REINSTATEMENTAMOUNT2", "REINSTATEMENTNOTE", "REINSTATEMENTPCT", "REINSTATEMENTVALUE"},
	},
	{
		Induk: "T_TREATY_LIMIT_GROUP", KunciAnak: "ClassOfBusinessList",
		Tabel: "T_TREATY_LIMIT_GROUP_COB", Seq: "SEQ_TT_LIMIT_GRP_COB",
		Kunci:        []string{"ClassOfBusiness", "ClassOfBusinessID", "TreatyGroup", "TreatyGroupID"},
		Kolom:        []string{"CLASSOFBUSINESS", "CLASSOFBUSINESSID", "TREATYGROUP", "TREATYGROUPID"},
		CacahTerukur: 15746,
	},
	// =====================================================================
	// MIGRASI 438 — tabel NILAI, beberapa larik ke satu tabel lewat `JENIS`
	// =====================================================================
	{
		Induk: "T_TREATY_LIMIT_DETAIL",
		// ⭐ 449 (8 Oktober 2026) — tiga larik `{Currency, Value}` tab Share
		// Prop (`CalculateShareList`, `FetchQSfromMaster`) — TANPA DDL.
		LarikGabung: []string{"IOOLimitList", "RetentionList", "CessionList", "EPIList",
			"RNMShareList", "RNMSpreadedList", "RNMSpreadedListRI",
			// ⭐ 450 — grid nilai rincian Limits Prop, TANPA DDL.
			"ReserveList", "PLAList", "CashLossList", "ClaimCoopList", "DeductionTotalList"},
		Tabel: "T_TREATY_LIMIT_AMOUNT", Seq: "SEQ_TT_LIMIT_AMOUNT",
		Kunci:        []string{"Currency", "CurrencyID", "Layer", "Note", "Value"},
		Kolom:        []string{"CURRENCY", "CURRENCYID", "LAYER", "NOTE", "VALUE"},
		CacahTerukur: 11475,
	},
	{
		Induk:       "T_TREATY_SHARE",
		LarikGabung: []string{"GrossPremiumList", "NetPremiumList"},
		Tabel:       "T_TREATY_SHARE_AMOUNT", Seq: "SEQ_TT_SHARE_AMOUNT",
		Kunci:        []string{"Currency", "Value"},
		Kolom:        []string{"CURRENCY", "VALUE"},
		CacahTerukur: 6097,
	},
	{
		Induk:       "T_TREATY_FAC_SHARE",
		LarikGabung: []string{"GrossPremiumList", "NetPremiumList"},
		Tabel:       "T_TREATY_FAC_SHARE_AMOUNT", Seq: "SEQ_TT_FAC_SHARE_AMOUNT",
		Kunci:        []string{"Currency", "Value"},
		Kolom:        []string{"CURRENCY", "VALUE"},
		CacahTerukur: 24,
	},
	// =====================================================================
	// MIGRASI 439 — skalar akar, kurs, besaran layer, ringkasan, dan total
	// =====================================================================
	//
	// ⭐ KETIGA `CacahTerukur: 0` DIISI 6 Oktober 2026 — sapuan atas
	// SELURUH kedua korpus, tanpa `ROWNUM` dan tanpa sampel:
	// 1.855 dokumen `M_TREATY_IN` + 280 dokumen `M_TREATY_IN_EDM`
	// (kedua sisi `New`/`OLDDATA`), 2.135 dokumen terurai, NOL gagal urai.
	//
	//                              M_TREATY_IN   EDM(2 sisi)    jumlah
	//   T_TREATY_REVISION                1.855           560     2.415
	//   T_TREATY_LIMIT_MEASURE           8.849         6.383    15.232
	//   T_TREATY_TOTAL                   6.395         5.098    11.493
	//
	// ⛔ YANG DIISI DI BAWAH ADALAH KOLOM `M_TREATY_IN`, bukan jumlahnya —
	// dan itu BUKAN pilihan selera. Sapuan yang sama mengukur ulang ke-27
	// entri lain di peta ini, dan ke-27-nya cocok PERSIS dengan kolom
	// `M_TREATY_IN` (702, 2.298, 4.210, 11.475, 15.746, …). Jadi arti
	// `CacahTerukur` di peta ini sudah tertetapkan oleh 27 saudaranya:
	// korpus `M_TREATY_IN` saja. Mengisi tiga entri dengan jumlah kedua
	// korpus akan membuat baris `TOTAL` di `pemuat/jalankan.go`
	// membandingkan dua hal yang berbeda.
	//
	// ⚠️ AKIBATNYA, DAN IA DINYATAKAN BUKAN DISEMBUNYIKAN: begitu
	// `M_TREATY_IN_EDM` ikut dimuat (ia mendarat ke tabel yang SAMA,
	// dibedakan `MASTERID` berakhiran `#LAMA`), cacah nyata tiap tabel
	// melampaui angka di bawah dan `cetakCacah` akan menandai ⚠️ pada
	// SETIAP baris. Tanda itu benar secara mekanis dan menyesatkan secara
	// arti. Memperbaikinya menuntut penyebut KEDUA di struktur ini, dan itu
	// keputusan yang belum diambil — lihat
	// `docs/PERTANYAAN-TERBUKA-PENYEBUT-REKONSILIASI.md`.
	{
		// ⭐ ENAM BELAS kunci TERAKHIR ditambahkan migrasi `444`, 6 Oktober
		// 2026 — medan yang layar tampilkan tetapi tabelnya belum punya
		// tempatnya. Yang terbesar akibatnya ketujuh `Reporting*`: terisi
		// di 1.219 dari 1.855 dokumen, dan selama ini tab Reporting Period
		// memperlihatkan kepala KOSONG pada seribu dua ratus kontrak yang
		// sungguh punya nilainya (`PERTANYAAN-TERBUKA-LAYAR-PEGA.md` §17).
		//
		// ⭐ EMPAT kunci TERAKHIR ditambahkan migrasi `445` (parkir di folder
		// `treatyinadjustment`), 7 Oktober 2026 — skalar akar tab Share
		// Non-Prop, tempat Save menyimpannya (keputusan pemakai: skema v2,
		// akar `TreatyIn` → `T_TREATY_REVISION`). Layar MEMBACANYA lewat
		// `BacaShareAkarRevisi` yang toleran terhadap kolom yang belum ada.
		//
		// ⚠️ `Comment` -> `COMMENTTEKS`: `COMMENT` kata tercadang Oracle.
		// Ini SATU-SATUNYA pasangan di peta ini yang nama kolomnya bukan
		// nama kuncinya dalam huruf besar, dan `TestPetaPendaratanCocokDenganDDL`
		// mengadu keduanya terhadap DDL — bukan terhadap aturan penamaan.
		Akar:  true,
		Tabel: "T_TREATY_REVISION", Seq: "SEQ_TT_REVISION",
		Kunci: []string{"ID", "OLDID", "ProportionType", "TreatyContractName", "TeritorialScope", "TreatyYear", "Ceding", "LeadingReinsSource", "Commencement", "Termination", "ContractRefNo", "Bordeaux", "BordereauxNote", "AccountingMode", "AccountingModeNonProp", "TreatyLeader", "IsMultipleRetro", "EDMState", "EDMMaterialType", "EDMEffective", "StatusAkseptasi", "IsProRate", "ProRateDays", "ProRateTotalDays", "ProRatePercent", "IsEditData", "TotalEgnpiAmount", "TotalEgnpiProportion", "TotalLimitsROL", "FacultativeShare", "FacultativeShareBrokerage", "ValueDifference.RNMShare", "ValueDifference.BrokeragePercent", "ValueDifference.TotalLimitsROL", "Exclusions", "ExclusionsP", "SpecialConditions", "SpecialConditionsP", "SpecialConditionsp", "CedingID", "LeadingReinsSourceID", "LeadingReinsID", "Information", "Position", "PositionUsername", "ChooseStatusAkseptasi", "AccumulationPeriod", "Comment", "ReportingStart", "ReportingEnd", "ReportingPeriod", "ReportingInterval", "ReportingSubmission", "ReportingConfirmation", "ReportingSettlement", "RNMShare", "BrokeragePercent", "RNMShareAcrossTheBoard", "RnmShareDeducted", "RNMShareP", "BrokeragePercentP", "OptionLimit", "InstallmentNo", "RevisionState", "ViewState", "RevisionDate"},
		Kolom: []string{"IDX", "OLDID", "PROPORTIONTYPE", "TREATYCONTRACTNAME", "TERITORIALSCOPE", "TREATYYEAR", "CEDING", "LEADINGREINSSOURCE", "COMMENCEMENT", "TERMINATION", "CONTRACTREFNO", "BORDEAUX", "BORDEREAUXNOTE", "ACCOUNTINGMODE", "ACCOUNTINGMODENONPROP", "TREATYLEADER", "ISMULTIPLERETRO", "EDMSTATE", "EDMMATERIALTYPE", "EDMEFFECTIVE", "STATUSAKSEPTASI", "ISPRORATE", "PRORATEDAYS", "PRORATETOTALDAYS", "PRORATEPERCENT", "ISEDITDATA", "TOTALEGNPIAMOUNT", "TOTALEGNPIPROPORTION", "TOTALLIMITSROL", "FACULTATIVESHARE", "FACULTATIVESHAREBROKERAGE", "VALUEDIFF_RNMSHARE", "VALUEDIFF_BROKERAGEPCT", "VALUEDIFF_TOTALLIMITSROL", "EXCLUSIONS", "EXCLUSIONSP", "SPECIALCONDITIONS", "SPECIALCONDITIONSP", "SPECIALCONDITIONSLC", "CEDINGID", "LEADINGREINSSOURCEID", "LEADINGREINSID", "INFORMATION", "POSITION", "POSITIONUSERNAME", "CHOOSESTATUSAKSEPTASI", "ACCUMULATIONPERIOD", "COMMENTTEKS", "REPORTINGSTART", "REPORTINGEND", "REPORTINGPERIOD", "REPORTINGINTERVAL", "REPORTINGSUBMISSION", "REPORTINGCONFIRMATION", "REPORTINGSETTLEMENT", "RNMSHARE", "BROKERAGEPERCENT", "RNMSHAREACROSSTHEBOARD", "RNMSHAREDEDUCTED", "RNMSHAREP", "BROKERAGEPERCENTP", "OPTIONLIMIT", "INSTALLMENTNO", "REVISIONSTATE", "VIEWSTATE", "REVISIONDATE"},
		// Satu baris per dokumen. 1.855 (EDM: +560).
		CacahTerukur: 1855,
	},
	{
		// ⭐ `T_TREATY_HAZARD_LIMIT` — migrasi `446` (parkir di folder
		// `treatyinadjustment`), 7 Oktober 2026. Diagram v2
		// (`Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx`): `TreatyIn
		// [BATAS_BAHAYA]`, 1:1. Sepuluh skalar AKAR yang sebelumnya nol
		// kolom pendaratan: Event Limits Non-Prop (akar berisi di 49 dari 772
		// kontrak Non-Prop) dan Max Co-Insurance Panel tab Co-Ins Scale (173 /
		// 87 dari 1.855 dokumen). Tabel akar KEDUA sesudah `T_TREATY_REVISION`
		// — `KunciTakTerpetakan` menggabungkan kunci keduanya.
		Akar:  true,
		Tabel: "T_TREATY_HAZARD_LIMIT", Seq: "SEQ_TT_HAZARD_LIMIT",
		Kunci: []string{"RSMDLimit", "CurrencyRSMD", "Earthquake", "CurrencyEarthquake", "FloodJab", "CurrencyFloodJab", "FloodNation", "CurrencyFloodNat", "MaxCoGroup", "MaxCoNonGroup"},
		Kolom: []string{"RSMDLIMIT", "CURRENCYRSMD", "EARTHQUAKE", "CURRENCYEARTHQUAKE", "FLOODJAB", "CURRENCYFLOODJAB", "FLOODNATION", "CURRENCYFLOODNAT", "MAXCOGROUP", "MAXCONONGROUP"},
		// Satu baris per dokumen, seperti `T_TREATY_REVISION`.
		CacahTerukur: 1855,
	},
	{
		Induk: "T_TREATY_LIMITS",
		// ⭐ 450 — `MDPMinList` layer Non-Prop (`LimitCalculation`), TANPA DDL.
		LarikGabung: []string{"MDPList", "PremiumEarnedList", "EgnpiTotalList",
			"MDPMinList"},
		Tabel: "T_TREATY_LIMIT_MEASURE", Seq: "SEQ_TT_LIMIT_MEASURE",
		Kunci: []string{"Currency", "CurrencyID", "Value"},
		Kolom: []string{"CURRENCY", "CURRENCYID", "VALUE"},
		// Tiga larik di dalam tiap `Limits[]`. 8.849 (EDM: +6.383).
		CacahTerukur: 8849,
	},
	{
		Larik: "LimitSummaryList",
		Tabel: "T_TREATY_LIMIT_SUMMARY", Seq: "SEQ_TT_LIMIT_SUMMARY",
		Kunci:        []string{"Layer", "LayerPart", "LayerPartType", "LayerType", "Currency", "Currency2", "Limit", "Limit2", "AggregateLimit", "AggregateLimit2", "Deductible", "Deductible2", "MDP", "MDP2", "Note"},
		Kolom:        []string{"LAYER", "LAYERPART", "LAYERPARTTYPE", "LAYERTYPE", "CURRENCY", "CURRENCY2", "LIMITVAL", "LIMITVAL2", "AGGREGATELIMIT", "AGGREGATELIMIT2", "DEDUCTIBLE", "DEDUCTIBLE2", "MDP", "MDP2", "NOTE"},
		CacahTerukur: 2855,
	},
	{
		LarikGabung: []string{"TotalRetentionAmountNP", "TotalEgnpiAmountNP", "TotalLimitIOONP",
			"TotalLimitDeductblNP", "TotalLimitPremiEarnNP", "TotalLimitMDPNP",
			// ⭐ 448 (7 Oktober 2026) — total tab Share Prop/Non-Prop dan
			// Installment yang tombol Save kirim; bentuknya sama, `JENIS`
			// yang membedakan, jadi nol DDL.
			"TotalShareRnmProp", "TotalSpreadedRnmProp", "TotalSpreadedRnmRIProp",
			"TotalShareRnmNP", "TotalShareGrossNP", "TotalShareGrossMinNP", "TotalShareDeductionNP",
			"TotalShareNetNP", "TotalSpreadedNetPremi", "TotalSpreadedNetPremiRI", "TotalInstallmentNP"},
		Tabel: "T_TREATY_TOTAL", Seq: "SEQ_TT_TOTAL",
		Kunci: []string{"Currency", "CurrencyID", "Value"},
		Kolom: []string{"CURRENCY", "CURRENCYID", "VALUE"},
		// Enam larik akar. 6.395 (EDM: +5.098).
		CacahTerukur: 6395,
	},
	{
		// ⭐ 448 (7 Oktober 2026) — dua larik ringkasan tab Share Non-Prop
		// (`SummaryLimitShare`), dibedakan `JENIS`.
		LarikGabung: []string{
			"LimitShareSummaryList", "LimitFacShareSummaryList"},
		Tabel: "T_TREATY_SHARE_SUMMARY", Seq: "SEQ_TT_SHARE_SUMMARY",
		Kunci: []string{"LayerType", "Layer", "LayerPartType", "LayerPart", "Note", "Limit", "Limit2", "MDP", "MDP2", "Deductible", "Deductible2", "NetPremi", "NetPremi2"},
		Kolom: []string{"LAYERTYPE", "LAYER", "LAYERPARTTYPE", "LAYERPART", "NOTE", "LIMITVAL", "LIMITVAL2", "MDP", "MDP2", "DEDUCTIBLE", "DEDUCTIBLE2", "NETPREMI", "NETPREMI2"},
	},
}

// Indeks tetap ke dalam `PetaPendaratan` — dipakai rekonsiliasi dan uji.
//
// ⛔ Pasangan Installment adalah SATU-SATUNYA yang indeksnya disebut, sebab
// rekonsiliasi butir angsuran harus menapaki induknya lebih dulu. Pasangan
// lain ditapaki secara umum lewat `Induk`/`KunciAnak`.
const (
	// IndeksAngsuran menunjuk entri `T_TREATY_INSTALLMENT`.
	IndeksAngsuran = 5
	// IndeksButirAngsuran menunjuk entri `T_TREATY_INSTALLMENT_ITEM`,
	// dan ia WAJIB sesudah induknya.
	IndeksButirAngsuran = 6
)

// LarikAnakAngsuran - kunci larik butir di dalam elemen `Installment`.
const LarikAnakAngsuran = "InstallmentList"

// AkhiranSisiLama membedakan sisi `Old` sebuah dokumen Adjustment dari sisi
// `New`-nya di dalam `MASTERID`.
//
// ⛔ SATU tetapan, dipakai pemuat DAN pembaca. Dua tempat yang menulis
// akhiran yang sama akan berselisih suatu hari, dan selisihnya terbaca
// sebagai "sisi Old kosong" — bukan sebagai galat.
//
// ⚠️ `#` dipilih sebab nol pengenal di korpus memuatnya: pengenal Treaty In
// berupa tujuh angka, pengenal Adjustment berbentuk `<asal>/R<nn>`.
const AkhiranSisiLama = "#LAMA"
