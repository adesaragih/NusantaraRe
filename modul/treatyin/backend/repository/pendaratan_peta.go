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
// ⚠️ `M_TREATYIN_INSTALLMENTITEM` tidak punya `Larik` di puncak dokumen:
// jalurnya `Installment[].InstallmentList`. Pemuatnya memperlakukannya
// khusus, dan `Larik` yang kosong adalah penandanya.
var PetaPendaratan = []Pendaratan{
	{
		Larik: "ReportingPeriodList", Tabel: "M_TREATYIN_REPORTINGPERIOD", Seq: "SEQ_MTI_REPORTINGPERIOD",
		Kunci:        []string{"AutoCalculate", "ConfirmationDue", "InitialDate", "Period", "SettlementDue", "SubmissionDue", "pxObjClass"},
		Kolom:        []string{"AUTOCALCULATE", "CONFIRMATIONDUE", "INITIALDATE", "PERIOD", "SETTLEMENTDUE", "SUBMISSIONDUE", "PXOBJCLASS"},
		CacahTerukur: 4548,
	},
	{
		Larik: "Portfolio", Tabel: "M_TREATYIN_PORTFOLIO", Seq: "SEQ_MTI_PORTFOLIO",
		Kunci:        []string{"Description", "Type", "TypePortfolio", "pxObjClass"},
		Kolom:        []string{"DESCRIPTION", "TYPE", "TYPEPORTFOLIO", "PXOBJCLASS"},
		CacahTerukur: 1925,
	},
	{
		Larik: "AccumulationList", Tabel: "M_TREATYIN_ACCUMULATION", Seq: "SEQ_MTI_ACCUMULATION",
		Kunci:        []string{"Period", "ReportDate", "SubDays", "SubDueDate", "pxObjClass"},
		Kolom:        []string{"PERIOD", "REPORTDATE", "SUBDAYS", "SUBDUEDATE", "PXOBJCLASS"},
		CacahTerukur: 60,
	},
	{
		Larik: "EGNPI", Tabel: "M_TREATYIN_EGNPI", Seq: "SEQ_MTI_EGNPI",
		Kunci:        []string{"Amount", "AmountIDR", "AsDate", "ClassOfBusiness", "Currency", "CurrencyID", "Note", "Proportion", "TreatyGroup", "TreatyGroupID", "pyTemplateRichTextEditor", "pxObjClass"},
		Kolom:        []string{"AMOUNT", "AMOUNTIDR", "ASDATE", "CLASSOFBUSINESS", "CURRENCY", "CURRENCYID", "NOTE", "PROPORTION", "TREATYGROUP", "TREATYGROUPID", "PYTEMPLATERICHTEXTEDITOR", "PXOBJCLASS"},
		CacahTerukur: 2298,
	},
	{
		Larik: "Retention", Tabel: "M_TREATYIN_RETENTION", Seq: "SEQ_MTI_RETENTION",
		Kunci:        []string{"Amount", "ClassOfBusiness", "Currency", "CurrencyID", "Note", "TreatyGroup", "TreatyGroupID", "pxObjClass"},
		Kolom:        []string{"AMOUNT", "CLASSOFBUSINESS", "CURRENCY", "CURRENCYID", "NOTE", "TREATYGROUP", "TREATYGROUPID", "PXOBJCLASS"},
		CacahTerukur: 2511,
	},
	{
		Larik: "Installment", Tabel: "M_TREATYIN_INSTALLMENT", Seq: "SEQ_MTI_INSTALLMENT",
		Kunci:        []string{"AmountTotal", "Currency", "PctTotal", "pxListSubscript", "pxObjClass"},
		Kolom:        []string{"AMOUNTTOTAL", "CURRENCY", "PCTTOTAL", "PXLISTSUBSCRIPT", "PXOBJCLASS"},
		CacahTerukur: 796,
	},
	{
		// Anak `Installment`; lihat catatan di atas.
		Larik: "", Tabel: "M_TREATYIN_INSTALLMENTITEM", Seq: "SEQ_MTI_INSTALLMENTITEM",
		Kunci:        []string{"Amount", "Currency", "DueDate", "Installment", "InstallmentPct", "PaymentDate", "WPC", "pxObjClass"},
		Kolom:        []string{"AMOUNT", "CURRENCY", "DUEDATE", "INSTALLMENT", "INSTALLMENTPCT", "PAYMENTDATE", "WPC", "PXOBJCLASS"},
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
		Kunci:        []string{"CoInShare", "PctLimit", "pxCreateDateTime", "pxCreateOpName", "pxCreateOperator", "pxCreateSystemID", "pxObjClass"},
		Kolom:        []string{"COINSHARE", "PCTLIMIT", "PXCREATEDATETIME", "PXCREATEOPNAME", "PXCREATEOPERATOR", "PXCREATESYSTEMID", "PXOBJCLASS"},
		CacahTerukur: 702,
	},
	{
		// ⚠️ `Date` -> `TANGGAL`: `DATE` kata cadangan Oracle (ORA-00923).
		Larik: "CommentList", Tabel: "M_TREATYIN_COMMENT", Seq: "SEQ_MTI_COMMENT",
		Kunci:        []string{"Date", "ConvertDate", "HasHistory", "IsApproved", "OperatorName", "Suggest", "pxObjClass"},
		Kolom:        []string{"TANGGAL", "CONVERTDATE", "HASHISTORY", "ISAPPROVED", "OPERATORNAME", "SUGGEST", "PXOBJCLASS"},
		CacahTerukur: 11365,
	},
}

// LarikAnakAngsuran adalah kunci larik bersarang di dalam `Installment`.
const LarikAnakAngsuran = "InstallmentList"

// IndeksAngsuran dan IndeksButirAngsuran menunjuk pasangan induk-anak di
// dalam `PetaPendaratan`. Dinyatakan sebagai tetapan supaya pemuat tidak
// mencarinya dengan nama tabel - pencarian yang gagal diam-diam melewati
// seluruh butir angsuran.
const (
	IndeksAngsuran      = 5
	IndeksButirAngsuran = 6
)
