package repository

// Peta tabel pendaratan yang layar Penyesuaian baca.
//
// ---------------------------------------------------------------------
// ⛔ MENGAPA DAFTAR KEDUA, padahal modul Treaty In sudah punya petanya
// ---------------------------------------------------------------------
// Penjaga arsitektur `TestModulTidakMengimporModulLain` melarang modul
// mengimpor modul lain: ketergantungan lintas modul lewat
// `inti/backend/kontrak`, atau tidak sama sekali. Mengimpor
// `treatyin/backend/repository` dari sini karena itu DITOLAK — dan menolaknya
// benar: modul yang saling impor tidak dapat dilepas sendiri-sendiri.
//
// ⭐ Daftar kedua yang menyalin daftar pertama SELALU menjadi basi. Karena
// itu kesamaannya DIBUKTIKAN, bukan dipercaya:
// `uji/lintasmodul/peta_pendaratan_test.go` mengadu daftar ini dengan
// `treatyin.PetaPendaratan` nama demi nama, kunci demi kunci — dan lapisan
// `uji/` memang boleh mengenal kedua modul.
//
// ⚠️ Isinya DIBANGKITKAN dari peta Treaty In, bukan diketik ulang. Yang
// mengetik ulang puluhan nama kolom akan salah satu, dan yang salah satu
// tidak terbaca sampai kolomnya kosong di layar.

// larikPendaratan - satu tabel pendaratan yang layar ini baca.
type larikPendaratan struct {
	// Larik - nama larik Pega di akar dokumen. Kosong bila `Akar` atau
	// `Gabung` terisi.
	Larik string
	// Gabung - beberapa larik DI SATU tabel, dibedakan kolom `JENIS`: larik
	// AKAR bila `Induk` kosong, larik di dalam simpul induk bila terisi.
	Gabung []string
	// Induk - tabel INDUK bagi tabel anak (`IDINDUK` → `ID` induk). Tabel anak
	// TIDAK dibaca ke larik datar; ia dirangkai ke `Pohon` (lihat
	// `pendaratan_pohon.go`).
	Induk string
	// KunciAnak - nama larik di dalam simpul induk (`Detail`, `DeductionList`, …).
	KunciAnak string
	// Akar - skalar akar dokumen; satu baris per sisi.
	Akar  bool
	Tabel string
	Kunci []string
	Kolom []string
}

// petaPendaratanPenyesuaian - tabel yang layar Penyesuaian baca, kedua sisi.
var petaPendaratanPenyesuaian = []larikPendaratan{
	{Larik: "ReportingPeriodList", Tabel: "T_TREATY_REPORTING_PERIOD",
		Kunci: []string{"AutoCalculate", "ConfirmationDue", "InitialDate", "Period", "SettlementDue", "SubmissionDue"},
		Kolom: []string{"AUTOCALCULATE", "CONFIRMATIONDUE", "INITIALDATE", "PERIOD", "SETTLEMENTDUE", "SUBMISSIONDUE"}},
	{Larik: "Portfolio", Tabel: "T_TREATY_PORTFOLIO",
		Kunci: []string{"Description", "Type", "TypePortfolio"},
		Kolom: []string{"DESCRIPTION", "TYPE", "TYPEPORTFOLIO"}},
	{Larik: "AccumulationList", Tabel: "T_TREATY_ACCUMULATION",
		Kunci: []string{"Period", "ReportDate", "SubDays", "SubDueDate"},
		Kolom: []string{"PERIOD", "REPORTDATE", "SUBDAYS", "SUBDUEDATE"}},
	{Larik: "EGNPI", Tabel: "T_TREATY_EGNPI",
		Kunci: []string{"Amount", "AmountIDR", "AsDate", "ClassOfBusiness", "Currency", "CurrencyID", "Note", "Proportion", "TreatyGroup", "TreatyGroupID", "pyTemplateRichTextEditor"},
		Kolom: []string{"AMOUNT", "AMOUNTIDR", "ASDATE", "CLASSOFBUSINESS", "CURRENCY", "CURRENCYID", "NOTE", "PROPORTION", "TREATYGROUP", "TREATYGROUPID", "PYTEMPLATERICHTEXTEDITOR"}},
	{Larik: "Retention", Tabel: "T_TREATY_RETENTION",
		Kunci: []string{"Amount", "ClassOfBusiness", "Currency", "CurrencyID", "Note", "TreatyGroup", "TreatyGroupID"},
		Kolom: []string{"AMOUNT", "CLASSOFBUSINESS", "CURRENCY", "CURRENCYID", "NOTE", "TREATYGROUP", "TREATYGROUPID"}},
	{Larik: "CommentList", Tabel: "T_VIEW_COMMENT",
		Kunci: []string{"Date", "ConvertDate", "HasHistory", "IsApproved", "OperatorName", "Suggest"},
		Kolom: []string{"TANGGAL", "CONVERTDATE", "HASHISTORY", "ISAPPROVED", "OPERATORNAME", "SUGGEST"}},
	{Larik: "Limits", Tabel: "T_TREATY_LIMITS",
		Kunci: []string{"AdjRate", "AgregateLimit", "AgregateLimit2", "Cover", "Currency", "Currency2", "CurrencyID", "CurrencyRelation", "Deductible", "Deductible2", "ID", "IsCombineMDP", "Layer", "LayerPart", "LayerPartType", "LayerType", "Limit", "Limit2", "MDPMinPct", "MDPPct", "NoRIPCalculation", "ROLPct", "ReinstatementNote", "ReinstatementPct", "ReinstatementValue", "SumLossRatio", "SumTotalAchievIncured", "SumTotalAchievNetPremium", "TreatyType", "TreatyTypeID"},
		Kolom: []string{"ADJRATE", "AGREGATELIMIT", "AGREGATELIMIT2", "COVER", "CURRENCY", "CURRENCY2", "CURRENCYID", "CURRENCYRELATION", "DEDUCTIBLE", "DEDUCTIBLE2", "IDX", "ISCOMBINEMDP", "LAYER", "LAYERPART", "LAYERPARTTYPE", "LAYERTYPE", "LIMIT", "LIMIT2", "MDPMINPCT", "MDPPCT", "NORIPCALCULATION", "ROLPCT", "REINSTATEMENTNOTE", "REINSTATEMENTPCT", "REINSTATEMENTVALUE", "SUMLOSSRATIO", "SUMTOTALACHIEVINCURED", "SUMTOTALACHIEVNETPREMIUM", "TREATYTYPE", "TREATYTYPEID"}},
	{Akar: true, Tabel: "T_TREATY_REVISION",
		Kunci: []string{"ID", "OLDID", "ProportionType", "TreatyContractName", "TeritorialScope", "TreatyYear", "Ceding", "LeadingReinsSource", "Commencement", "Termination", "ContractRefNo", "Bordeaux", "BordereauxNote", "AccountingMode", "AccountingModeNonProp", "TreatyLeader", "IsMultipleRetro", "EDMState", "EDMMaterialType", "EDMEffective", "StatusAkseptasi", "IsProRate", "ProRateDays", "ProRateTotalDays", "ProRatePercent", "IsEditData", "TotalEgnpiAmount", "TotalEgnpiProportion", "TotalLimitsROL", "FacultativeShare", "FacultativeShareBrokerage", "ValueDifference.RNMShare", "ValueDifference.BrokeragePercent", "Exclusions", "ExclusionsP", "SpecialConditions", "SpecialConditionsP", "SpecialConditionsp", "CedingID", "LeadingReinsSourceID", "LeadingReinsID", "Information", "Position", "PositionUsername", "ChooseStatusAkseptasi", "AccumulationPeriod", "Comment", "ReportingStart", "ReportingEnd", "ReportingPeriod", "ReportingInterval", "ReportingSubmission", "ReportingConfirmation", "ReportingSettlement", "RNMShare", "BrokeragePercent", "RNMShareAcrossTheBoard", "RnmShareDeducted", "RNMShareP", "BrokeragePercentP", "OptionLimit", "InstallmentNo", "RevisionState", "ViewState", "RevisionDate"},
		Kolom: []string{"IDX", "OLDID", "PROPORTIONTYPE", "TREATYCONTRACTNAME", "TERITORIALSCOPE", "TREATYYEAR", "CEDING", "LEADINGREINSSOURCE", "COMMENCEMENT", "TERMINATION", "CONTRACTREFNO", "BORDEAUX", "BORDEREAUXNOTE", "ACCOUNTINGMODE", "ACCOUNTINGMODENONPROP", "TREATYLEADER", "ISMULTIPLERETRO", "EDMSTATE", "EDMMATERIALTYPE", "EDMEFFECTIVE", "STATUSAKSEPTASI", "ISPRORATE", "PRORATEDAYS", "PRORATETOTALDAYS", "PRORATEPERCENT", "ISEDITDATA", "TOTALEGNPIAMOUNT", "TOTALEGNPIPROPORTION", "TOTALLIMITSROL", "FACULTATIVESHARE", "FACULTATIVESHAREBROKERAGE", "VALUEDIFF_RNMSHARE", "VALUEDIFF_BROKERAGEPCT", "EXCLUSIONS", "EXCLUSIONSP", "SPECIALCONDITIONS", "SPECIALCONDITIONSP", "SPECIALCONDITIONSLC", "CEDINGID", "LEADINGREINSSOURCEID", "LEADINGREINSID", "INFORMATION", "POSITION", "POSITIONUSERNAME", "CHOOSESTATUSAKSEPTASI", "ACCUMULATIONPERIOD", "COMMENTTEKS", "REPORTINGSTART", "REPORTINGEND", "REPORTINGPERIOD", "REPORTINGINTERVAL", "REPORTINGSUBMISSION", "REPORTINGCONFIRMATION", "REPORTINGSETTLEMENT", "RNMSHARE", "BROKERAGEPERCENT", "RNMSHAREACROSSTHEBOARD", "RNMSHAREDEDUCTED", "RNMSHAREP", "BROKERAGEPERCENTP", "OPTIONLIMIT", "INSTALLMENTNO", "REVISIONSTATE", "VIEWSTATE", "REVISIONDATE"}},
	{Larik: "LimitSummaryList", Tabel: "T_TREATY_LIMIT_SUMMARY",
		Kunci: []string{"Layer", "LayerPart", "LayerPartType", "LayerType", "Currency", "Currency2", "Limit", "Limit2", "AggregateLimit", "AggregateLimit2", "Deductible", "Deductible2", "MDP", "MDP2", "Note"},
		Kolom: []string{"LAYER", "LAYERPART", "LAYERPARTTYPE", "LAYERTYPE", "CURRENCY", "CURRENCY2", "LIMITVAL", "LIMITVAL2", "AGGREGATELIMIT", "AGGREGATELIMIT2", "DEDUCTIBLE", "DEDUCTIBLE2", "MDP", "MDP2", "NOTE"}},
	// ⭐ 7 Oktober 2026 — enam larik tingkat atas yang kerangka layar ini
	// ikat dan pemuat SUDAH daratkan untuk kedua korpus, tetapi tidak
	// pernah dibaca di sini: tab Installment, Co-Ins Scale, dan Share
	// tampil "No items" padahal barisnya ada. Disalin dari peta Treaty In
	// apa adanya; `uji/lintasmodul` mengadu kesamaannya.
	{Larik: "Installment", Tabel: "T_TREATY_INSTALLMENT",
		Kunci: []string{"AmountTotal", "Currency", "PctTotal", "pxListSubscript"},
		Kolom: []string{"AMOUNTTOTAL", "CURRENCY", "PCTTOTAL", "PXLISTSUBSCRIPT"}},
	{Larik: "CoInScale", Tabel: "M_TREATYIN_COINSCALE",
		Kunci: []string{"CoInShare", "PctLimit", "pxCreateDateTime", "pxCreateOpName", "pxCreateOperator", "pxCreateSystemID"},
		Kolom: []string{"COINSHARE", "PCTLIMIT", "PXCREATEDATETIME", "PXCREATEOPNAME", "PXCREATEOPERATOR", "PXCREATESYSTEMID"}},
	{Larik: "Share", Tabel: "T_TREATY_SHARE",
		Kunci: []string{"Cover", "Layer", "LayerPart", "LayerPartType", "LayerType", "RNMShare", "SpreadingTotalPctXOL", "SpreadingTypeIDXOL", "SpreadingTypeXOL"},
		Kolom: []string{"COVER", "LAYER", "LAYERPART", "LAYERPARTTYPE", "LAYERTYPE", "RNMSHARE", "SPREADINGTOTALPCTXOL", "SPREADINGTYPEIDXOL", "SPREADINGTYPEXOL"}},
	{Larik: "ShareReins", Tabel: "T_TREATY_RETRO_SHARE",
		Kunci: []string{"Layer", "ReinsID", "ReinsName", "SharePct"},
		Kolom: []string{"LAYER", "REINSID", "REINSNAME", "SHAREPCT"}},
	{Larik: "FacultativeShareList", Tabel: "T_TREATY_FAC_SHARE",
		Kunci: []string{"Cover", "Layer", "LayerPart", "LayerPartType", "LayerType", "Limit", "Limit2", "SpreadingTotalPctXOL", "SpreadingTypeIDXOL", "SpreadingTypeXOL", "SpreadingTypeXOLRetro", "SpreadingTypeXOLRetroID"},
		Kolom: []string{"COVER", "LAYER", "LAYERPART", "LAYERPARTTYPE", "LAYERTYPE", "LIMIT", "LIMIT2", "SPREADINGTOTALPCTXOL", "SPREADINGTYPEIDXOL", "SPREADINGTYPEXOL", "SPREADINGTYPEXOLRETRO", "SPREADINGTYPEXOLRETROID"}},
	{Larik: "ShareFacultativeReinsurers", Tabel: "T_TREATY_FAC_REINSURER",
		Kunci: []string{"Layer", "ReinsID", "ReinsName", "SharePct"},
		Kolom: []string{"LAYER", "REINSID", "REINSNAME", "SHAREPCT"}},
	// ⭐ 7 Oktober 2026 — TIGA BELAS tabel ANAK, disalin dari peta Treaty In
	// apa adanya (Induk · KunciAnak/Gabung · Kunci · Kolom). Dibaca ke
	// `Pohon`, bukan ke larik datar: rumus Limits/Share/Installment
	// (rute `/hitung/*` Treaty In) menuntut simpul bersarang.
	// `T_TREATY_INSTALLMENT_ITEM` tidak punya `KunciAnak` di peta Treaty In —
	// larik anaknya tetapan `LarikAnakAngsuran` = `InstallmentList`.
	{Induk: "T_TREATY_INSTALLMENT", KunciAnak: "InstallmentList", Tabel: "T_TREATY_INSTALLMENT_ITEM",
		Kunci: []string{"Amount", "Currency", "DueDate", "Installment", "InstallmentPct", "PaymentDate", "WPC"},
		Kolom: []string{"AMOUNT", "CURRENCY", "DUEDATE", "INSTALLMENT", "INSTALLMENTPCT", "PAYMENTDATE", "WPC"}},
	{Induk: "T_TREATY_LIMITS", KunciAnak: "Detail", Tabel: "T_TREATY_LIMIT_DETAIL",
		Kunci: []string{"AchievementPct", "Brokerage", "CashLoss", "CessionPct", "ClaimCooperation", "ClassOfBusiness", "CurrencyCashLoss", "CurrencyClaimCooperation", "CurrencyEPI", "CurrencyEarthquake", "CurrencyFloodJab", "CurrencyFloodNat", "CurrencyID", "CurrencyPLA", "CurrencyRSMD", "EPI", "Earthquake", "FloodJab", "FloodNation", "ID", "IOOPct", "LossRatio", "LowerBand", "PLA", "ParentID", "Periode", "PremiumReservePct", "ProfitCommision", "ProfitME", "ProfitYDCF", "QSPct", "RIOGR", "RIONR", "RNMShare", "RSMDLimit", "ReisuredParticipant", "RetentionPct", "ShareNote", "SpreadingTotalPct", "SpreadingType", "SpreadingTypeID", "SumLossRatio", "SumTotalAchievIncured", "SumTotalAchievNetPremium", "Surplus", "TotalAchAfterClaim", "TotalAchCASHCALL", "TotalAchEstCASHCALL", "TotalAchIncured", "TotalAchNetPremium", "TotalAchOsCASHCALL", "TotalAchOsClaim", "TotalAchOuts", "TotalAchPaid", "TotalAchPremium", "TreatyGroup", "TreatyGroupID", "TreatyType", "UpperBand", "Currency", "NusaReLimit"},
		Kolom: []string{"ACHIEVEMENTPCT", "BROKERAGE", "CASHLOSS", "CESSIONPCT", "CLAIMCOOPERATION", "CLASSOFBUSINESS", "CURRENCYCASHLOSS", "CURRENCYCLAIMCOOPERATION", "CURRENCYEPI", "CURRENCYEARTHQUAKE", "CURRENCYFLOODJAB", "CURRENCYFLOODNAT", "CURRENCYID", "CURRENCYPLA", "CURRENCYRSMD", "EPI", "EARTHQUAKE", "FLOODJAB", "FLOODNATION", "IDX", "IOOPCT", "LOSSRATIO", "LOWERBAND", "PLA", "PARENTID", "PERIODE", "PREMIUMRESERVEPCT", "PROFITCOMMISION", "PROFITME", "PROFITYDCF", "QSPCT", "RIOGR", "RIONR", "RNMSHARE", "RSMDLIMIT", "REISUREDPARTICIPANT", "RETENTIONPCT", "SHARENOTE", "SPREADINGTOTALPCT", "SPREADINGTYPE", "SPREADINGTYPEID", "SUMLOSSRATIO", "SUMTOTALACHIEVINCURED", "SUMTOTALACHIEVNETPREMIUM", "SURPLUS", "TOTALACHAFTERCLAIM", "TOTALACHCASHCALL", "TOTALACHESTCASHCALL", "TOTALACHINCURED", "TOTALACHNETPREMIUM", "TOTALACHOSCASHCALL", "TOTALACHOSCLAIM", "TOTALACHOUTS", "TOTALACHPAID", "TOTALACHPREMIUM", "TREATYGROUP", "TREATYGROUPID", "TREATYTYPE", "UPPERBAND", "CURRENCY", "NUSARELIMIT"}},
	{Induk: "T_TREATY_LIMITS", KunciAnak: "TreatyGroupList", Tabel: "T_TREATY_LIMIT_GROUP",
		Kunci: []string{"ClassOfBusiness", "ClassOfBusinessID", "IsROLProfile", "TreatyGroup", "TreatyGroupID"},
		Kolom: []string{"CLASSOFBUSINESS", "CLASSOFBUSINESSID", "ISROLPROFILE", "TREATYGROUP", "TREATYGROUPID"}},
	{Induk: "T_TREATY_SHARE", KunciAnak: "SpreadingListXOL", Tabel: "T_TREATY_SHARE_SPREADING",
		Kunci: []string{"ParentReinsTypeID", "Pct", "ReinsTypeID", "ReinsTypeName", "Rp", "Usd"},
		Kolom: []string{"PARENTREINSTYPEID", "PCT", "REINSTYPEID", "REINSTYPENAME", "RP", "USD"}},
	{Induk: "T_TREATY_SHARE", KunciAnak: "DeductionList", Tabel: "T_TREATY_SHARE_DEDUCTION",
		Kunci: []string{"Comment", "Currency", "Deduction", "DeductionPct", "DeductionPctCalculate"},
		Kolom: []string{"COMMENT_", "CURRENCY", "DEDUCTION", "DEDUCTIONPCT", "DEDUCTIONPCTCALCULATE"}},
	{Induk: "T_TREATY_FAC_SHARE", KunciAnak: "DeductionList", Tabel: "T_TREATY_FAC_SHARE_DEDUCTION",
		Kunci: []string{"Comment", "Currency", "Deduction", "DeductionPct", "DeductionPctCalculate"},
		Kolom: []string{"COMMENT_", "CURRENCY", "DEDUCTION", "DEDUCTIONPCT", "DEDUCTIONPCTCALCULATE"}},
	{Induk: "T_TREATY_LIMIT_DETAIL", KunciAnak: "COBList", Tabel: "T_TREATY_LIMIT_COB",
		Kunci: []string{"ClassOfBusiness", "ClassOfBusinessID", "TreatyGroup", "TreatyGroupID"},
		Kolom: []string{"CLASSOFBUSINESS", "CLASSOFBUSINESSID", "TREATYGROUP", "TREATYGROUPID"}},
	{Induk: "T_TREATY_LIMIT_DETAIL", KunciAnak: "AchievementLists", Tabel: "T_TREATY_LIMIT_ACHIEVEMENT",
		Kunci: []string{"BROKERAGE", "CASHCALL", "CASHCALLtoIDR", "Conversion", "Currency", "CurrencyID", "EstimationCASHCALL", "EstimationCASHCALLtoIDR", "IDPEGA", "IncuredClaim", "IncuredClaimtoIDR", "LossRatio", "NETPREMIUM", "NETPREMIUMtoIDR", "NOOFFER", "NOPOLIS", "OutstandingCASHCALL", "OutstandingCASHCALLtoIDR", "OutstandingClaim", "OutstandingClaimtoIDR", "PREMIUM", "PREMIUMtoIDR", "PaidClaim", "PaidClaimtoIDR", "Period", "QUARTERYEAR", "Quarter", "RICOMM", "SOBNAME", "TREATYGROUPNAME", "TREATYTYPE", "Total", "TotalAchievCASHCALL", "TotalAchievEstCASHCALL", "TotalAchievIncured", "TotalAchievNetPremium", "TotalAchievOsCASHCALL", "TotalAchievOsClaim", "TotalAchievOuts", "TotalAchievPaid", "TotalAchievPremium", "TotalAfterClaim", "TotaltoIDR", "Value"},
		Kolom: []string{"BROKERAGE", "CASHCALL", "CASHCALLTOIDR", "CONVERSION", "CURRENCY", "CURRENCYID", "ESTIMATIONCASHCALL", "ESTIMATIONCASHCALLTOIDR", "IDPEGA", "INCUREDCLAIM", "INCUREDCLAIMTOIDR", "LOSSRATIO", "NETPREMIUM", "NETPREMIUMTOIDR", "NOOFFER", "NOPOLIS", "OUTSTANDINGCASHCALL", "OUTSTANDINGCASHCALLTOIDR", "OUTSTANDINGCLAIM", "OUTSTANDINGCLAIMTOIDR", "PREMIUM", "PREMIUMTOIDR", "PAIDCLAIM", "PAIDCLAIMTOIDR", "PERIOD", "QUARTERYEAR", "QUARTER", "RICOMM", "SOBNAME", "TREATYGROUPNAME", "TREATYTYPE", "TOTAL", "TOTALACHIEVCASHCALL", "TOTALACHIEVESTCASHCALL", "TOTALACHIEVINCURED", "TOTALACHIEVNETPREMIUM", "TOTALACHIEVOSCASHCALL", "TOTALACHIEVOSCLAIM", "TOTALACHIEVOUTS", "TOTALACHIEVPAID", "TOTALACHIEVPREMIUM", "TOTALAFTERCLAIM", "TOTALTOIDR", "VALUE"}},
	// ⭐ 449 — `Detail.SpreadingList` tab Share Prop (salinan peta Treaty In).
	{Induk: "T_TREATY_LIMIT_DETAIL", KunciAnak: "SpreadingList", Tabel: "T_TREATY_LIMIT_SPREADING",
		Kunci: []string{"ParentReinsTypeID", "Pct", "ReinsTypeID", "ReinsTypeName", "Rp", "Usd", "Value"},
		Kolom: []string{"PARENTREINSTYPEID", "PCT", "REINSTYPEID", "REINSTYPENAME", "RP", "USD", "VALUE"}},
	// ⭐ 450 — Deduction, Parameter Achievement, Reinstatement (salinan peta Treaty In).
	{Induk: "T_TREATY_LIMIT_DETAIL", KunciAnak: "DeductionList", Tabel: "T_TREATY_LIMIT_DEDUCTION",
		Kunci: []string{"Comment", "Currency", "CurrencyID", "Deduction", "DeductionPct", "DeductionPctCalculate"},
		Kolom: []string{"COMMENT_", "CURRENCY", "CURRENCYID", "DEDUCTION", "DEDUCTIONPCT", "DEDUCTIONPCTCALCULATE"}},
	{Induk: "T_TREATY_LIMIT_DETAIL", KunciAnak: "CurrencyList", Tabel: "T_TREATY_LIMIT_ACH_PARAM",
		Kunci: []string{"Parameter", "AchievementPctGross", "LossRatioGross"},
		Kolom: []string{"PARAMETER", "ACHIEVEMENTPCTGROSS", "LOSSRATIOGROSS"}},
	{Induk: "T_TREATY_LIMITS", KunciAnak: "Reinstatement_List", Tabel: "T_TREATY_LIMIT_REINSTATEMENT",
		Kunci: []string{"AdditionalAmount1", "AdditionalAmount2", "AdditionalPct", "ID", "ReinstatementAmount1", "ReinstatementAmount2", "ReinstatementNote", "ReinstatementPct", "ReinstatementValue"},
		Kolom: []string{"ADDITIONALAMOUNT1", "ADDITIONALAMOUNT2", "ADDITIONALPCT", "IDX", "REINSTATEMENTAMOUNT1", "REINSTATEMENTAMOUNT2", "REINSTATEMENTNOTE", "REINSTATEMENTPCT", "REINSTATEMENTVALUE"}},
	{Induk: "T_TREATY_LIMIT_GROUP", KunciAnak: "ClassOfBusinessList", Tabel: "T_TREATY_LIMIT_GROUP_COB",
		Kunci: []string{"ClassOfBusiness", "ClassOfBusinessID", "TreatyGroup", "TreatyGroupID"},
		Kolom: []string{"CLASSOFBUSINESS", "CLASSOFBUSINESSID", "TREATYGROUP", "TREATYGROUPID"}},
	{Induk: "T_TREATY_LIMIT_DETAIL", Gabung: []string{"IOOLimitList", "RetentionList", "CessionList", "EPIList", "RNMShareList", "RNMSpreadedList", "RNMSpreadedListRI", "ReserveList", "PLAList", "CashLossList", "ClaimCoopList", "DeductionTotalList"}, Tabel: "T_TREATY_LIMIT_AMOUNT",
		Kunci: []string{"Currency", "CurrencyID", "Layer", "Note", "Value"},
		Kolom: []string{"CURRENCY", "CURRENCYID", "LAYER", "NOTE", "VALUE"}},
	{Induk: "T_TREATY_SHARE", Gabung: []string{"GrossPremiumList", "NetPremiumList"}, Tabel: "T_TREATY_SHARE_AMOUNT",
		Kunci: []string{"Currency", "Value"},
		Kolom: []string{"CURRENCY", "VALUE"}},
	{Induk: "T_TREATY_FAC_SHARE", Gabung: []string{"GrossPremiumList", "NetPremiumList"}, Tabel: "T_TREATY_FAC_SHARE_AMOUNT",
		Kunci: []string{"Currency", "Value"},
		Kolom: []string{"CURRENCY", "VALUE"}},
	{Induk: "T_TREATY_LIMITS", Gabung: []string{"MDPList", "PremiumEarnedList", "EgnpiTotalList", "MDPMinList"}, Tabel: "T_TREATY_LIMIT_MEASURE",
		Kunci: []string{"Currency", "CurrencyID", "Value"},
		Kolom: []string{"CURRENCY", "CURRENCYID", "VALUE"}},
	{Gabung: []string{"TotalRetentionAmountNP", "TotalEgnpiAmountNP", "TotalLimitIOONP", "TotalLimitDeductblNP", "TotalLimitPremiEarnNP", "TotalLimitMDPNP", "TotalShareRnmProp", "TotalSpreadedRnmProp", "TotalSpreadedRnmRIProp", "TotalShareRnmNP", "TotalShareGrossNP", "TotalShareGrossMinNP", "TotalShareDeductionNP", "TotalShareNetNP", "TotalSpreadedNetPremi", "TotalSpreadedNetPremiRI", "TotalInstallmentNP"}, Tabel: "T_TREATY_TOTAL",
		Kunci: []string{"Currency", "CurrencyID", "Value"},
		Kolom: []string{"CURRENCY", "CURRENCYID", "VALUE"}},
}

// LarikPendaratanUji - bentuk terbuka `larikPendaratan`, untuk uji lintas
// modul saja.
//
// ⛔ Dibuka HANYA untuk diadu. Jalur baca di modul ini memakai bentuk yang
// tertutup; yang terbuka ada supaya `uji/lintasmodul` dapat membuktikan
// salinan peta ini tidak basi terhadap milik Treaty In.
type LarikPendaratanUji struct {
	Larik     string
	Gabung    []string
	Akar      bool
	Induk     string
	KunciAnak string
	Tabel     string
	Kunci     []string
	Kolom     []string
}

// PetaPendaratanPenyesuaianUntukUji membuka salinan peta untuk diadu.
func PetaPendaratanPenyesuaianUntukUji() []LarikPendaratanUji {
	out := make([]LarikPendaratanUji, 0, len(petaPendaratanPenyesuaian))
	for _, p := range petaPendaratanPenyesuaian {
		out = append(out, LarikPendaratanUji{
			Larik: p.Larik, Gabung: p.Gabung, Akar: p.Akar,
			Induk: p.Induk, KunciAnak: p.KunciAnak,
			Tabel: p.Tabel, Kunci: p.Kunci, Kolom: p.Kolom,
		})
	}
	return out
}
