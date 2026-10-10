# Struktur Tabel — NB Treaty In

Acuan bentuk tabel modul `nbtreatyin`. **Dibangkitkan** `docs/alat/skema.py` dari
`backend/models/katalog.go` — jangan disunting tangan; sunting katalognya.

Tipe ditulis sebagai kategori logis: teks · angka desimal · bilangan bulat · DATE.
Uang dan persen **angka desimal** `NUMBER(38,10)` — skala minimal 9 (diagram sheet NB Treaty In Prop F20),
10 supaya bagi rata spreading NB presisi 10 tersimpan utuh (J69); tidak pernah float (ADR-0003).
Golongan (uang / persen / kode / penanda / tanggal) ada di kolom *Golongan*.

Tabel yang **dibaca, tidak dibuat** modul ini dideklarasikan di `MODUL.md`.

## T_GENERAL_POLIS_TREATY

Satu baris per generasi polis; kunci utama bersama `T_WORK_POLIS` (ID-7). Generasi tertutup = ada penerus yang `OLD_POLIS_ID`-nya menunjuk baris ini (ID-10).

| Kolom | Tipe | Null | Kunci | Golongan | Properti Pega |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK, FK T_WORK_POLIS | kode | kasus |
| `NOPOLIS` | teks | ya | UQ (NOPOLIS, PRODKE) | kode | `PolicyTreatyIn.PolicyNo` |
| `PRODKE` | bilangan bulat | tidak | UQ (NOPOLIS, PRODKE) | cacah | generasi; NB = 0 |
| `NOENDORS` | teks | ya |  | kode | json_polis |
| `OLD_POLIS_ID` | teks | ya | UQ, FK T_WORK_POLIS | kode | generasi sebelumnya |
| `TGL_INPUT` | DATE | ya |  | tanggal-waktu | json_polis |
| `USERNAME` | teks | ya |  | kode | identitas akses login (P4) |
| `POSITION_NOTE` | teks | ya |  | teks | `PositionNote` |
| `NB_STATUS` | teks | ya |  | teks | `NBStatus` |
| `TREATY_IN_ID` | teks | ya |  | kode | `TreatyIn.ID` |
| `NO_OFFER` | teks | ya |  | kode | `NoOffer` |
| `MASTER_ID` | teks | ya |  | kode | `MasterID` |
| `IS_APPROVED` | teks | ya |  | penanda | `IsApproved` |
| `SUGGEST` | teks | ya |  | teks | `Suggest` |
| `SUGGEST_DATE` | DATE | ya |  | tanggal-waktu | `SuggestDate` |
| `OPERATOR_NAME` | teks | ya |  | teks | `OperatorName` |
| `IS_NEW_POLICY_NON_PROP` | teks | ya |  | penanda | `IsNewPolicyNonProp` |
| `EDM_TYPE` | teks | ya |  | kode | `EDMType` |
| `IS_EDM_INPUT_ON_NB` | teks | ya |  | penanda | `IsEDMInputOnNB` |
| `HAS_FAC_OUT` | teks | ya |  | penanda | `HasFacOut` |
| `FLAG_PPH` | teks | ya |  | penanda | `FlagPPH` |
| `FLAG_RETRO_TREATY` | teks | ya |  | penanda | `FlagRetroTreaty` |
| `DUE_TO` | teks | ya |  | penanda | `DueTo` |
| `TYPE_TAX` | teks | ya |  | kode | `TypeTax` |
| `STATEMENT_TYPE` | teks | ya |  | kode | `StatementType` |
| `TREATY_GROUP_ID` | teks | ya |  | kode | `TreatyGroupID` |
| `TREATY_GROUP_NAME` | teks | ya |  | teks | `TreatyGroupName` |
| `TREATY_GROUP_OLD_ID` | teks | ya |  | kode | `TreatyGroupOldID` |
| `OJK_BUSINESS_ID` | teks | ya |  | kode | `OJKBusinessID` |
| `ID_NEW_BISNIS` | teks | ya |  | kode | `IDNewBisnis` |
| `BIZ_CODE` | teks | ya |  | kode | `BizCode` |
| `BIZ_NAME` | teks | ya |  | teks | `BizName` |
| `SOB` | teks | ya |  | kode | `SOB` |
| `SOB_NAME` | teks | ya |  | teks | `SOBName` |
| `CEDING_CO` | teks | ya |  | kode | `CedingCo` |
| `CEDING_CO_NAME` | teks | ya |  | teks | `CedingCoName` |
| `INSURED_ID` | teks | ya |  | kode | `InsuredID` |
| `INSURED_NAME` | teks | ya |  | teks | `InsuredName` |
| `MARKETING_OFFICER` | teks | ya |  | teks | `MarketingOfficer` |
| `TREATY_TYPE` | teks | ya |  | kode | `TreatyType` |
| `TREATY_YEAR` | teks | ya |  | kode | `TreatyYear` |
| `CURRENCY` | teks | ya |  | kode | `Currency` |
| `ID_CURRENCY` | teks | ya |  | kode | `IDCurrency` |
| `SHARE_CURRENCY` | teks | ya |  | kode | `ShareCurrency` |
| `QUARTAL` | teks | ya |  | kode | `Quartal` |
| `YEAR_OF_QUARTAL` | teks | ya |  | kode | `YearOfQuartal` |
| `CLAIM_TYPE` | teks | ya |  | kode | `ClaimType` |
| `CLAIM_PAYMENT_TYPE` | teks | ya |  | kode | `ClaimPaymentType` |
| `INSTALLMENT` | teks | ya |  | kode | `Installment` |
| `REMARK` | teks | ya |  | teks | `Remark` |
| `START_DATE` | DATE | ya |  | tanggal | `StartDate` |
| `END_DATE` | DATE | ya |  | tanggal | `EndDate` |
| `STATEMENT_DATE` | DATE | ya |  | tanggal-waktu | `StatementDate` |
| `TGL_PROD` | DATE | ya |  | tanggal-waktu | `ProductionDate` |
| `GROSS_PREMIUM` | angka desimal | ya |  | uang | `GrossPremium` |
| `GROSS_CLAIM` | angka desimal | ya |  | uang | `GrossClaim` |
| `PREMI_OGP` | angka desimal | ya |  | uang | `PremiOgp` |
| `RESULT_OGP1` | angka desimal | ya |  | uang | `ResultOgp1` |
| `RESULT_OGP2` | angka desimal | ya |  | uang | `ResultOgp2` |
| `PREMI_ONP` | angka desimal | ya |  | uang | `PremiOnp` |
| `RESULT_ONP1` | angka desimal | ya |  | uang | `ResultOnp1` |
| `RESULT_ONP2` | angka desimal | ya |  | uang | `ResultOnp2` |
| `CLAIM` | angka desimal | ya |  | uang | `Claim` |
| `OUTSTANDING_CLAIM` | angka desimal | ya |  | uang | `OutstandingClaim` |
| `SALVAGE_VALUE` | angka desimal | ya |  | uang | `SalvageValue` |
| `EXCESS_LOSS` | angka desimal | ya |  | uang | `ExcessLoss` |
| `NET_PREMIUM` | angka desimal | ya |  | uang | `NetPremium` |
| `BALANCE_DUE_TO` | angka desimal | ya |  | uang | `BalanceDueTo` |
| `BALANCE_BEFORE_TAX` | angka desimal | ya |  | uang | `BalanceBeforeTax` |
| `BALANCE_BEFORE_PPH` | angka desimal | ya |  | uang | `BalanceBeforePPH` |
| `DEDUCTION1` | angka desimal | ya |  | uang | `Deduction1` |
| `DEDUCTION2` | angka desimal | ya |  | uang | `Deduction2` |
| `BROKERAGE_FEE_SEBENARNYA` | angka desimal | ya |  | uang | `BrokerageFeeSebenarnya` |
| `PPH_VALUE` | angka desimal | ya |  | uang | `PPHValue` |
| `PPN_VALUE` | angka desimal | ya |  | uang | `PPNValue` |
| `SHARE_VALUE` | angka desimal | ya |  | uang | `ShareValue` |
| `RI_COMM_OGP` | angka desimal | ya |  | persen | `RiCommOgp` |
| `OVERIDDING_COMM_OGP` | angka desimal | ya |  | persen | `OveriddingCommOgp` |
| `RI_COMM_ONP` | angka desimal | ya |  | persen | `RiCommOnp` |
| `OVERIDDING_COMM_ONP` | angka desimal | ya |  | persen | `OveriddingCommOnp` |
| `IS_SOA_UPLOAD` | teks | ya |  | kode | `IsSOAUpload` |

## T_POLIS_QUOTATION

Halaman `Quotation` / `PolicyTreatyIn.QuotationData`, 1:1 (ID-23).

| Kolom | Tipe | Null | Kunci | Golongan | Properti Pega |
| --- | --- | --- | --- | --- | --- |
| `POLIS_ID` | teks | tidak | PK, FK T_GENERAL_POLIS_TREATY | kode | induk |
| `PROPORTIONAL_TYPE` | teks | ya |  | kode | `ProportionalType` |
| `MO_ID` | teks | ya |  | kode | `MOID` |
| `BUSINESS_CODE` | teks | ya |  | kode | `BusinessCode` |
| `BUSINESS_OLD_ID` | teks | ya |  | kode | `BusinessOldId` |
| `GROUP_PANEL` | teks | ya |  | kode | `GroupPanel` |
| `SOURCE_OF_BUSINESS` | teks | ya |  | kode | `SourceOfBusiness` |
| `TYPE` | teks | ya |  | kode | `Type` |
| `EDM_TYPE` | teks | ya |  | kode | `EdmType` |
| `OLD_POLICY_NO` | teks | ya |  | kode | `OldPolicyNo` |
| `MARKETING_NAME` | teks | ya |  | teks | `MarketingName` |
| `BUSINESS_NAME` | teks | ya |  | teks | `BusinessName` |
| `BUSINESS_FAC` | teks | ya |  | kode | `BusinessFac` |
| `INSURED_ID` | teks | ya |  | kode | `InsuredID` |
| `INSURED_NAME` | teks | ya |  | teks | `InsuredName` |
| `NO_OFFER_SLIP` | teks | ya |  | teks | `NoOfferSlip` |
| `IS_SURVEY_REPORT` | teks | ya |  | penanda | `IsSurveyReport` |
| `MARKETING_CODE` | teks | ya |  | kode | `MarketingCode` |
| `TEAM_GROUP` | teks | ya |  | kode | `TeamGroup` |
| `BRANCH_CODE` | teks | ya |  | kode | `BranchCode` |
| `BRANCH_NAME` | teks | ya |  | teks | `BranchName` |

## T_POLIS_CEDING

← `QuotationData.CedingCoList` (ID-24), di bawah `T_POLIS_QUOTATION` (diagram O39).

| Kolom | Tipe | Null | Kunci | Golongan | Properti Pega |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | kode | baris |
| `QUOTATION_ID` | teks | tidak | FK T_POLIS_QUOTATION, UQ (QUOTATION_ID, NOURUT) | kode | induk |
| `NOURUT` | bilangan bulat | tidak | UQ (QUOTATION_ID, NOURUT) | cacah | urutan baris (ID-11) |
| `CEDING_CO_ID` | teks | ya |  | kode | `CedingCo` |
| `CEDING_CO_NAME` | teks | ya |  | teks | `CedingCoName` |

## T_POLIS_INSTALMENT

← `PolicyTreatyIn.ListInstallment` (ID-26).

| Kolom | Tipe | Null | Kunci | Golongan | Properti Pega |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | kode | baris |
| `POLIS_ID` | teks | tidak | FK T_GENERAL_POLIS_TREATY, UQ (POLIS_ID, NOURUT) | kode | induk |
| `NOURUT` | bilangan bulat | tidak | UQ (POLIS_ID, NOURUT) | cacah | urutan baris (ID-11) |
| `INSTALLMENT_NO` | bilangan bulat | ya |  | cacah | `InstallmentNo` |
| `DUE_DATE` | DATE | ya |  | tanggal | `DueDate` |
| `INSTALLMENT_PERCENTAGE` | angka desimal | ya |  | persen | `InstallmentPercentage` |
| `PREMIUM` | angka desimal | ya |  | uang | `Premium` |
| `PAYMENT_TOTAL` | angka desimal | ya |  | uang | `PaymentTotal` |
| `PREMIUM_AFTER_PPH` | angka desimal | ya |  | uang | `PremiumAfterPPH` |
| `PREMIUM_AFTER_PPN` | angka desimal | ya |  | uang | `PremiumAfterPPN` |
| `PREMIUM_AFTER_TAX` | angka desimal | ya |  | uang | `PremiumAfterTax` |
| `CURRENCY` | teks | ya |  | kode | `Currency` |
| `ID_CURRENCY` | teks | ya |  | kode | `IDCurrency` |
| `PPN` | angka desimal | ya |  | uang | `PPN` |
| `PPH` | angka desimal | ya |  | uang | `PPh` |
| `PAYMENT_TOTAL_AFTER_PPN` | angka desimal | ya |  | uang | `PaymentTotalAfterPPN` |
| `PAYMENT_TOTAL_AFTER_TAX` | angka desimal | ya |  | uang | `PaymentTotalAfterTax` |

## T_POLIS_INSTALMENT_DETAIL

← `ListInstallment().InstallmentList`, non-proporsional.

| Kolom | Tipe | Null | Kunci | Golongan | Properti Pega |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | kode | baris |
| `INSTALMENT_ID` | teks | tidak | FK T_POLIS_INSTALMENT, UQ (INSTALMENT_ID, NOURUT) | kode | induk |
| `NOURUT` | bilangan bulat | tidak | UQ (INSTALMENT_ID, NOURUT) | cacah | urutan baris (ID-11) |
| `INSTALLMENT_NO` | bilangan bulat | ya |  | cacah | `InstallmentNo` |
| `DUE_DATE` | DATE | ya |  | tanggal | `DueDate` |
| `PAYMENT_DATE` | DATE | ya |  | tanggal | `PaymentDate` |
| `INSTALLMENT_PERCENTAGE` | angka desimal | ya |  | persen | `InstallmentPercentage` |
| `PREMIUM` | angka desimal | ya |  | uang | `Premium` |
| `PAYMENT_TOTAL` | angka desimal | ya |  | uang | `PaymentTotal` |
| `PREMIUM_AFTER_PPH` | angka desimal | ya |  | uang | `PremiumAfterPPH` |
| `PREMIUM_AFTER_PPN` | angka desimal | ya |  | uang | `PremiumAfterPPN` |
| `PREMIUM_AFTER_TAX` | angka desimal | ya |  | uang | `PremiumAfterTax` |
| `CURRENCY` | teks | ya |  | kode | `Currency` |
| `ID_CURRENCY` | teks | ya |  | kode | `IDCurrency` |

## T_POLIS_SPREADING

← `PolicyTreatyIn.SpreadingRiskList` (ID-28).

| Kolom | Tipe | Null | Kunci | Golongan | Properti Pega |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | kode | baris |
| `POLIS_ID` | teks | tidak | FK T_GENERAL_POLIS_TREATY, UQ (POLIS_ID, NOURUT) | kode | induk |
| `NOURUT` | bilangan bulat | tidak | UQ (POLIS_ID, NOURUT) | cacah | urutan baris (ID-11) |
| `TREATY_TYPE` | teks | ya |  | kode | `TreatyType` |
| `TREATY_NAME` | teks | ya |  | teks | `TreatyName` |
| `CURRENCY` | teks | ya |  | kode | `Currency` |
| `CURRENCY_ID` | teks | ya |  | kode | `CurrencyID` |
| `SHARE_PERCENTAGE` | angka desimal | ya |  | persen | `SharePercentage` |
| `SPLIT_RNM_SHARE_PCT` | angka desimal | ya |  | persen | `SplitRNMSharePct` |
| `CLAIM_PERCENTAGE` | angka desimal | ya |  | persen | `ClaimPercentage` |
| `PREMIUM_SPREADED` | angka desimal | ya |  | uang | `PremiumSpreaded` |
| `CLAIM_SPREADED` | angka desimal | ya |  | uang | `ClaimSpreaded` |

## T_POLIS_XOL

← `PolicyTreatyIn.TreatyXOLList` (ID-29).

| Kolom | Tipe | Null | Kunci | Golongan | Properti Pega |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | kode | baris |
| `POLIS_ID` | teks | tidak | FK T_GENERAL_POLIS_TREATY, UQ (POLIS_ID, NOURUT) | kode | induk |
| `NOURUT` | bilangan bulat | tidak | UQ (POLIS_ID, NOURUT) | cacah | urutan baris (ID-11) |
| `CURRENCY` | teks | ya |  | kode | `Currency` |
| `ID_CURRENCY` | teks | ya |  | kode | `IDCurrency` |
| `GROSS_PREMI` | angka desimal | ya |  | uang | `GrossPremi` |
| `NET_PREMI` | angka desimal | ya |  | uang | `NetPremi` |
| `DEDUCTION` | angka desimal | ya |  | uang | `Deduction` |
| `DUE_TO` | teks | ya |  | penanda | `DueTo` |
| `DUE_TO_VALUE` | angka desimal | ya |  | uang | `DueToValue` |
| `BROKERAGE_FEE_SEBENARNYA` | angka desimal | ya |  | uang | `BrokerageFeeSebenarnya` |
| `PPH_VALUE` | angka desimal | ya |  | uang | `PPHValue` |
| `PPN_VALUE` | angka desimal | ya |  | uang | `PPNValue` |
| `NET_PREMI_AFTER_PPH` | angka desimal | ya |  | uang | `NetPremiAfterPPH` |
| `NET_PREMI_AFTER_PPN` | angka desimal | ya |  | uang | `NetPremiAfterPPN` |
| `NET_PREMI_AFTER_TAX` | angka desimal | ya |  | uang | `NetPremiAfterTax` |

## T_POLIS_XOL_LAYER

← `TreatyXOLList().ValueList` (ID-29).

| Kolom | Tipe | Null | Kunci | Golongan | Properti Pega |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | kode | baris |
| `XOL_ID` | teks | tidak | FK T_POLIS_XOL, UQ (XOL_ID, NOURUT) | kode | induk |
| `NOURUT` | bilangan bulat | tidak | UQ (XOL_ID, NOURUT) | cacah | urutan baris (ID-11) |
| `LAYER` | teks | ya |  | kode | `Layer` |
| `LAYER_TYPE` | teks | ya |  | kode | `LayerType` |
| `LAYER_PART` | teks | ya |  | kode | `LayerPart` |
| `LAYER_PART_TYPE` | teks | ya |  | kode | `LayerPartType` |
| `CURRENCY` | teks | ya |  | kode | `Currency` |
| `ID_CURRENCY` | teks | ya |  | kode | `IDCurrency` |
| `GROSS_PREMI` | angka desimal | ya |  | uang | `GrossPremi` |
| `NET_PREMI` | angka desimal | ya |  | uang | `NetPremi` |
| `DEDUCTION` | angka desimal | ya |  | uang | `Deduction` |
| `DUE_TO` | teks | ya |  | penanda | `DueTo` |
| `DUE_TO_VALUE` | angka desimal | ya |  | uang | `DueToValue` |
| `BROKERAGE_FEE_SEBENARNYA` | angka desimal | ya |  | uang | `BrokerageFeeSebenarnya` |
| `PPH_VALUE` | angka desimal | ya |  | uang | `PPHValue` |
| `PPN_VALUE` | angka desimal | ya |  | uang | `PPNValue` |
| `NET_PREMI_AFTER_PPH` | angka desimal | ya |  | uang | `NetPremiAfterPPH` |
| `NET_PREMI_AFTER_PPN` | angka desimal | ya |  | uang | `NetPremiAfterPPN` |
| `NET_PREMI_AFTER_TAX` | angka desimal | ya |  | uang | `NetPremiAfterTax` |

## T_POLIS_SURVEY

← `PolicyTreatyIn.QuotationData.SurveyReportList` (popup Historical Survey Report; keputusan work owner 06-10-2026).

| Kolom | Tipe | Null | Kunci | Golongan | Properti Pega |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | kode | baris |
| `POLIS_ID` | teks | tidak | FK T_GENERAL_POLIS_TREATY, UQ (POLIS_ID, NOURUT) | kode | induk |
| `NOURUT` | bilangan bulat | tidak | UQ (POLIS_ID, NOURUT) | cacah | urutan baris (ID-11) |
| `DATE_OF_SURVEY` | DATE | ya |  | tanggal | `DateofSurvey` |
| `SURVEYED_BY` | teks | ya |  | teks | `SurveyedBy` |
| `LOSS_PREVENTION` | angka desimal | ya |  | persen | `LossPrevention` |
| `REMARKS` | teks | ya |  | kode | `Remarks` |

## HISTORYAKSEPTASIPRODUCTION

⛔ **Tabel WARISAN `POOLDATA` — tidak dibuat, tidak diubah strukturnya** (MODUL.md *Tabel warisan*;
keputusan work owner K4 03-10-2026; spec-penyimpanan ID-31). Catatan `PolicyTreatyIn.SuggestList` ditulis
ke sini (`repository/usulan.go`, pengganti `RDBList/InsertViewSuggest_SQL`) dan dibaca balik untuk layar.
Tipe fisik milik tabel lama (belum dicek katalog Oracle — butir terbuka).

| Kolom | Tipe | Null | Kunci | Golongan | Properti Pega |
| --- | --- | --- | --- | --- | --- |
| `IDPEGA` | warisan | ya |  | — | `pyWorkPage.pzInsKey` |
| `TYPE_POLIS` | warisan | ya |  | — | `@replaceAll(pyWorkIDPrefix,"-","")` = `NB` |
| `NOURUT` | warisan | ya |  | — | berikutnya per IDPEGA (XML `.pxListSubscript`) |
| `POSISI` | warisan | ya |  | — | `"Policy"` |
| `PIC` | warisan | ya |  | — | `.OperatorName` — nama tampilan (P33) |
| `TGL_INP` | warisan | ya |  | — | `.Date` |
| `DIV` | warisan | ya |  | — | `OperatorID.pyOrgDivision` — NULL, tanpa sumber (butir terbuka) |
| `TYPE` | warisan | ya |  | — | `Quotation.BusinessFac` = `T` |
| `PUTARAN` | warisan | ya |  | — | `"2"` |
| `APPROVAL` | warisan | ya |  | — | `.IsApproved` 1 = Accept, 0 = Reject |
| `KETERANGAN` | warisan | ya |  | — | `substr(.Suggest, 0, 3990)` |
| `AKSES_LOGIN` | warisan | ya |  | — | `OperatorID.pyUserIdentifier` — identitas login (P4) |
| `B2B` | warisan | ya |  | — | `OfferFacIn.IsB2B` — NULL di NB |
| `BUSINESS_CODE` | warisan | ya |  | — | `Quotation.BusinessCode` |
| `PERCENT_RNM` | warisan | ya |  | — | `OfferFacIn.PercentShare` — NULL di NB |
