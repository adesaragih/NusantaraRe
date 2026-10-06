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
	// Gabung - beberapa larik akar DI SATU tabel, dibedakan kolom `JENIS`.
	Gabung []string
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
		Kunci: []string{"ID", "OLDID", "ProportionType", "TreatyContractName", "TeritorialScope", "TreatyYear", "Ceding", "LeadingReinsSource", "Commencement", "Termination", "ContractRefNo", "Bordeaux", "BordereauxNote", "AccountingMode", "AccountingModeNonProp", "TreatyLeader", "IsMultipleRetro", "EDMState", "EDMMaterialType", "EDMEffective", "StatusAkseptasi", "IsProRate", "ProRateDays", "ProRateTotalDays", "ProRatePercent", "IsEditData", "TotalEgnpiAmount", "TotalEgnpiProportion", "TotalLimitsROL", "FacultativeShare", "FacultativeShareBrokerage", "ValueDifference.RNMShare", "ValueDifference.BrokeragePercent", "Exclusions", "ExclusionsP", "SpecialConditions", "SpecialConditionsP", "SpecialConditionsp", "CedingID", "LeadingReinsSourceID", "LeadingReinsID", "Information", "Position", "PositionUsername", "ChooseStatusAkseptasi", "AccumulationPeriod", "Comment", "ReportingStart", "ReportingEnd", "ReportingPeriod", "ReportingInterval", "ReportingSubmission", "ReportingConfirmation", "ReportingSettlement"},
		Kolom: []string{"IDX", "OLDID", "PROPORTIONTYPE", "TREATYCONTRACTNAME", "TERITORIALSCOPE", "TREATYYEAR", "CEDING", "LEADINGREINSSOURCE", "COMMENCEMENT", "TERMINATION", "CONTRACTREFNO", "BORDEAUX", "BORDEREAUXNOTE", "ACCOUNTINGMODE", "ACCOUNTINGMODENONPROP", "TREATYLEADER", "ISMULTIPLERETRO", "EDMSTATE", "EDMMATERIALTYPE", "EDMEFFECTIVE", "STATUSAKSEPTASI", "ISPRORATE", "PRORATEDAYS", "PRORATETOTALDAYS", "PRORATEPERCENT", "ISEDITDATA", "TOTALEGNPIAMOUNT", "TOTALEGNPIPROPORTION", "TOTALLIMITSROL", "FACULTATIVESHARE", "FACULTATIVESHAREBROKERAGE", "VALUEDIFF_RNMSHARE", "VALUEDIFF_BROKERAGEPCT", "EXCLUSIONS", "EXCLUSIONSP", "SPECIALCONDITIONS", "SPECIALCONDITIONSP", "SPECIALCONDITIONSLC", "CEDINGID", "LEADINGREINSSOURCEID", "LEADINGREINSID", "INFORMATION", "POSITION", "POSITIONUSERNAME", "CHOOSESTATUSAKSEPTASI", "ACCUMULATIONPERIOD", "COMMENTTEKS", "REPORTINGSTART", "REPORTINGEND", "REPORTINGPERIOD", "REPORTINGINTERVAL", "REPORTINGSUBMISSION", "REPORTINGCONFIRMATION", "REPORTINGSETTLEMENT"}},
	{Larik: "LimitSummaryList", Tabel: "T_TREATY_LIMIT_SUMMARY",
		Kunci: []string{"Layer", "LayerPart", "LayerPartType", "LayerType", "Currency", "Currency2", "Limit", "Limit2", "AggregateLimit", "AggregateLimit2", "Deductible", "Deductible2", "MDP", "MDP2", "Note"},
		Kolom: []string{"LAYER", "LAYERPART", "LAYERPARTTYPE", "LAYERTYPE", "CURRENCY", "CURRENCY2", "LIMITVAL", "LIMITVAL2", "AGGREGATELIMIT", "AGGREGATELIMIT2", "DEDUCTIBLE", "DEDUCTIBLE2", "MDP", "MDP2", "NOTE"}},
	{Gabung: []string{"TotalRetentionAmountNP", "TotalEgnpiAmountNP", "TotalLimitIOONP", "TotalLimitDeductblNP", "TotalLimitPremiEarnNP", "TotalLimitMDPNP"}, Tabel: "T_TREATY_TOTAL",
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
	Larik  string
	Gabung []string
	Akar   bool
	Tabel  string
	Kunci  []string
	Kolom  []string
}

// PetaPendaratanPenyesuaianUntukUji membuka salinan peta untuk diadu.
func PetaPendaratanPenyesuaianUntukUji() []LarikPendaratanUji {
	out := make([]LarikPendaratanUji, 0, len(petaPendaratanPenyesuaian))
	for _, p := range petaPendaratanPenyesuaian {
		out = append(out, LarikPendaratanUji{
			Larik: p.Larik, Gabung: p.Gabung, Akar: p.Akar,
			Tabel: p.Tabel, Kunci: p.Kunci, Kolom: p.Kolom,
		})
	}
	return out
}
