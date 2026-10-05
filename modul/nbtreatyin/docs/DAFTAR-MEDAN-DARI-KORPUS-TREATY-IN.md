# Daftar medan Treaty In — disusun dari Activity dan Section

> `[keputusan work owner]` **23 September 2026** — *"Abaikan `JSON_DATAGUIDE`, ikuti dari data
> yang digunakan di Activity dan Section."*
>
> Berkas ini **menggantikan** data guide sebagai dasar penetapan daftar kolom.

---

## Cara menyusunnya

Disapu dari **302 berkas aturan** di dua modul:

| Modul | Activity | Section | DataTransform | Flow | RDBList |
| --- | ---: | ---: | ---: | ---: | ---: |
| NB Treaty In | 92 | 25 | 12 | 1 | 41 |
| EDM Treaty In | 66 | 18 | 10 | 1 | 36 |

⭐ **Dua titik buta yang sebelumnya meleset empat kali, kali ini ditutup:**

1. **Dua konvensi tag yang berlawanan.** `Activity` memakai `PropertiesName` **tanpa** awalan
   `py`; `DataTransform` dan `Flow` memakai `pyPropertiesName` **dengan** awalan. Sapuan yang
   hanya mengenal satu konvensi kehilangan separuh isi.
2. **Rujukan relatif di dalam loop.** Sapuan berjangkar `PolicyTreatyIn.…` tidak melihat
   `.PremiumSpreaded` atau `.CedingCoName` yang ditulis relatif. Sapuan ini menangkap keduanya.

⛔ `pyExpressionGadget` dibuang sebelum pencocokan. ⛔ Entitas HTML **tidak** di-unescape sebelum
pencocokan struktur — `&lt;` memecahkan pola tag.

---

## Hasil

| | Jumlah |
| --- | ---: |
| Medan unik NB Treaty In | **376** |
| Medan unik EDM Treaty In | **217** |
| ⭐ **Gabungan unik** | **394** |
| Di antaranya **tampil di layar** *(muncul di Section)* | **196** |
| Nama daun unik pada data guide | 126 |

⭐ **Korpus memuat 281 nama yang tidak ada di data guide.** Itu bukan selisih kecil — itu
sebagian besar daftarnya.

⚠️ **Tetapi data guide juga memuat 13 nama yang tidak tersapu korpus**, termasuk `OldData` ·
`ProdKe` · `EDMNo` · `OldPolicyNo` · `BranchCode`. Keduanya tidak lengkap sendirian.
⇒ **Yang dipakai adalah gabungannya**, dengan korpus sebagai dasar dan data guide sebagai
penambal.

---

## Daftar medan

Kolom **layar** bertanda ⭐ bila medan itu muncul di `Section` — artinya benar-benar ditampilkan
atau diketik pengguna, bukan sekadar variabel antara.

Kolom **DG** menunjukkan apakah nama itu juga ada di data guide.

| Medan | NB | EDM | layar | DG |
| --- | ---: | ---: | :---: | :---: |
| `AcceptStatus` | 2 | — |  | ⛔ |
| `AccumulationCode` | 1 | — |  | ⛔ |
| `AchievementLists` | 1 | 1 |  | ⛔ |
| `AdditionalAmount1` | 1 | 1 |  | ⛔ |
| `AdditionalAmount2` | 1 | 1 |  | ⛔ |
| `AdditionalPct` | 1 | 1 |  | ⛔ |
| `AdjRate` | 1 | — |  | ⛔ |
| `Amount` | 6 | 3 | ⭐ | ⛔ |
| `AmountIDR` | 1 | — |  | ⛔ |
| `AmountTotal` | 3 | 3 | ⭐ | ⛔ |
| `AnekaList` | 3 | — |  | ⛔ |
| `Approval` | 1 | — |  | ⛔ |
| `ASMCoverage` | 4 | — |  | ⛔ |
| `BalanceBeforePPH` | 6 | 5 | ⭐ | ✅ |
| `BalanceBeforeTax` | 6 | 5 | ⭐ | ✅ |
| `BalanceDueTo` | 15 | 15 | ⭐ | ✅ |
| `BizCode` | 3 | 4 |  | ⛔ |
| `BizName` | 6 | 2 | ⭐ | ⛔ |
| `BranchDetailID` | 1 | 1 |  | ⛔ |
| `BranchDetailName` | 1 | 1 |  | ⛔ |
| `BreakDownSpreadList` | 5 | — | ⭐ | ✅ |
| `BrokerageFee` | 1 | 1 |  | ⛔ |
| `BrokerageFeeSebenarnya` | 3 | 5 |  | ✅ |
| `BusinessCode` | 1 | 2 |  | ⛔ |
| `BusinessOldId` | 2 | 2 |  | ⛔ |
| `CARI1` | 18 | 14 | ⭐ | ⛔ |
| `CARI10` | 1 | — |  | ⛔ |
| `CARI12` | 1 | — |  | ⛔ |
| `CARI2` | 8 | 5 | ⭐ | ⛔ |
| `CARI3` | 7 | 3 | ⭐ | ⛔ |
| `CARIDATETIME` | 1 | — |  | ⛔ |
| `Category` | 2 | — |  | ⛔ |
| `Ceding` | — | 1 | ⭐ | ⛔ |
| `CEDING` | 4 | — | ⭐ | ⛔ |
| `CedingCo` | 2 | 2 |  | ✅ |
| `CedingCoName` | 8 | 4 | ⭐ | ✅ |
| `CEDINGID` | 2 | — |  | ⛔ |
| `CedingID` | 1 | 1 |  | ⛔ |
| `ChildCount` | 3 | — | ⭐ | ⛔ |
| `CITYNAME` | 1 | — |  | ⛔ |
| `Claim` | 8 | 10 | ⭐ | ✅ |
| `ClaimAmountIDR` | 2 | — |  | ⛔ |
| `ClaimEstimation` | 5 | — |  | ⛔ |
| `ClaimPaymentType` | 5 | 5 | ⭐ | ⛔ |
| `ClaimPercentage` | 9 | 11 | ⭐ | ✅ |
| `ClaimSpreaded` | 14 | 10 | ⭐ | ✅ |
| `ClaimType` | 9 | 11 | ⭐ | ⛔ |
| `CLASSOFBUSINESS` | 4 | — | ⭐ | ⛔ |
| `CLASSOFBUSINESSID` | 2 | — |  | ⛔ |
| `ClientID` | 2 | 1 |  | ⛔ |
| `ClientName` | 6 | 2 | ⭐ | ⛔ |
| `Commencement` | 5 | 2 | ⭐ | ⛔ |
| `CommentSuggest` | 1 | — |  | ⛔ |
| `Condition` | 2 | — |  | ⛔ |
| `Conversion` | 1 | 1 |  | ⛔ |
| `ConvertDate` | 1 | 1 |  | ⛔ |
| `CountAttach` | 4 | — |  | ⛔ |
| `Coverage` | 3 | — |  | ⛔ |
| `CoverageBasis` | 2 | — |  | ⛔ |
| `CoverageList` | 11 | — |  | ⛔ |
| `CoverageNote` | 3 | — |  | ⛔ |
| `CURR` | 2 | — |  | ⛔ |
| `Currency` | 35 | 22 | ⭐ | ✅ |
| `CURRENCYID` | 2 | — |  | ⛔ |
| `CurrencyID` | 9 | 1 | ⭐ | ✅ |
| `CurrencyList` | 2 | 1 |  | ⛔ |
| `CurrencyOldID` | 1 | — |  | ⛔ |
| `CurrentDateTime` | 1 | — |  | ⛔ |
| `CurrentMode` | 1 | — | ⭐ | ⛔ |
| `Date` | 5 | 4 | ⭐ | ✅ |
| `DateofSurvey` | 4 | — | ⭐ | ⛔ |
| `DateSuggest` | 2 | — |  | ⛔ |
| `Deductible` | 7 | — | ⭐ | ⛔ |
| `Deductible2` | 6 | — | ⭐ | ⛔ |
| `DeductibleList` | 2 | — |  | ⛔ |
| `Deduction` | 3 | 11 | ⭐ | ✅ |
| `DEDUCTION1` | 2 | — |  | ⛔ |
| `Deduction1` | 11 | 11 | ⭐ | ✅ |
| `DEDUCTION2` | 2 | — |  | ⛔ |
| `Deduction2` | 10 | 9 | ⭐ | ✅ |
| `DeductionList` | 3 | 5 |  | ⛔ |
| `DeductionTotalList` | 3 | 2 |  | ⛔ |
| `Detail` | 2 | 1 |  | ⛔ |
| `DistrictName` | 1 | — |  | ⛔ |
| `DueDate` | 12 | 12 | ⭐ | ✅ |
| `DueTo` | 12 | 13 | ⭐ | ✅ |
| `DueToValue` | 3 | 8 |  | ✅ |
| `EDMNo` | — | 8 | ⭐ | ✅ |
| `EdmType` | — | 1 | ⭐ | ⛔ |
| `EDMType` | 1 | 8 | ⭐ | ✅ |
| `EMAILADDRESS` | 1 | — |  | ⛔ |
| `EmailTypeUW` | 1 | — |  | ⛔ |
| `EmailTypeUWEDM` | 1 | — |  | ⛔ |
| `EmailTypeUWPolicy` | 1 | — |  | ⛔ |
| `EndDate` | 11 | 7 | ⭐ | ✅ |
| `EngineNumber` | 1 | — |  | ⛔ |
| `EPICURRENCY` | 2 | — | ⭐ | ⛔ |
| `EPIVALUE` | 2 | — | ⭐ | ⛔ |
| `ExcessLoss` | 6 | 8 | ⭐ | ✅ |
| `FillPaymentInstallmentEDMT` | — | 1 |  | ⛔ |
| `FilterApplied` | 1 | — | ⭐ | ⛔ |
| `FilterForOpportunityInd` | 1 | — | ⭐ | ⛔ |
| `FilterTermForEndorsement` | — | 1 | ⭐ | ⛔ |
| `FilterTermForOpportunity` | 1 | — | ⭐ | ⛔ |
| `FirstLoss` | 1 | — |  | ⛔ |
| `FlagOnGoingPolicy` | 1 | 1 |  | ⛔ |
| `FlagPPH` | 11 | 8 | ⭐ | ✅ |
| `FlagRetroTreaty` | 7 | 4 | ⭐ | ✅ |
| `getPageValue` | — | 1 |  | ⛔ |
| `getStringValue` | 1 | — |  | ⛔ |
| `GrossClaim` | 4 | 1 | ⭐ | ⛔ |
| `GrossPremi` | 3 | 11 | ⭐ | ✅ |
| `GrossPremium` | 7 | 5 | ⭐ | ✅ |
| `GrossPremiumList` | 6 | 7 |  | ⛔ |
| `GroupPanel` | 1 | 1 |  | ⛔ |
| `GuaranteeFund` | — | 1 |  | ⛔ |
| `HasFacOut` | 2 | — |  | ✅ |
| `HASIL` | 1 | 3 |  | ⛔ |
| `HASIL1` | 2 | 1 |  | ⛔ |
| `IDCurrency` | 16 | 15 | ⭐ | ✅ |
| `IDNewBisnis` | — | 1 |  | ✅ |
| `IDPEGA` | 3 | — |  | ⛔ |
| `IdxCargo` | 1 | — |  | ⛔ |
| `IdxLocation` | 1 | — |  | ⛔ |
| `IndexCargo` | 1 | — |  | ⛔ |
| `IndexLocation` | 1 | — |  | ⛔ |
| `IndexProperty` | 1 | — |  | ⛔ |
| `IndexPropertyItem` | 1 | — |  | ⛔ |
| `Installment` | 13 | 14 | ⭐ | ✅ |
| `InstallmentList` | 7 | 6 | ⭐ | ✅ |
| `InstallmentNo` | 11 | 11 | ⭐ | ✅ |
| `INSTALLMENTNO` | 2 | — |  | ⛔ |
| `InstallmentPct` | 7 | 3 | ⭐ | ⛔ |
| `InstallmentPercentage` | 13 | 14 | ⭐ | ✅ |
| `InsuredID` | 2 | 2 |  | ✅ |
| `InsuredName` | 8 | 3 | ⭐ | ✅ |
| `IsApproved` | 13 | 9 | ⭐ | ✅ |
| `isApprovedtoDeptHead` | 3 | — |  | ⛔ |
| `IsCedingConfirm` | 1 | — |  | ⛔ |
| `IsEDMInputOnNB` | 3 | — | ⭐ | ⛔ |
| `IsInputPct` | 1 | — |  | ⛔ |
| `IsNewPolicyListFormat` | 2 | — | ⭐ | ⛔ |
| `IsNewPolicyNonProp` | 8 | 11 | ⭐ | ✅ |
| `IsOJKNopolis` | 1 | 1 |  | ✅ |
| `IsOldData` | 1 | — |  | ⛔ |
| `IsSave` | 1 | — |  | ⛔ |
| `IsTopRisk` | 1 | — |  | ⛔ |
| `ItemType` | 2 | — |  | ⛔ |
| `ItemTypeID` | 1 | — |  | ⛔ |
| `KATEGORI_1` | — | 1 |  | ⛔ |
| `KATEGORI_2` | 1 | 1 |  | ⛔ |
| `LAYER` | 4 | — | ⭐ | ⛔ |
| `Layer` | 9 | 11 | ⭐ | ✅ |
| `LayerList` | 4 | 2 |  | ⛔ |
| `LayerPart` | 9 | 11 | ⭐ | ✅ |
| `LAYERPART` | 4 | — | ⭐ | ⛔ |
| `LAYERPARTTYPE` | 4 | — | ⭐ | ⛔ |
| `LayerPartType` | 9 | 11 | ⭐ | ✅ |
| `LAYERTYPE` | 4 | — | ⭐ | ⛔ |
| `LayerType` | 9 | 11 | ⭐ | ✅ |
| `Leader0` | 2 | — |  | ⛔ |
| `Leader1` | 2 | — |  | ⛔ |
| `LeadingReinsSource` | — | 1 | ⭐ | ⛔ |
| `LeadingReinsSourceID` | 1 | 1 |  | ⛔ |
| `LicensePlate` | 1 | — |  | ⛔ |
| `Limit` | 6 | 1 | ⭐ | ⛔ |
| `Limit2` | 5 | 1 | ⭐ | ⛔ |
| `LIMITCURRENCY` | 4 | — | ⭐ | ⛔ |
| `LimitofLiability` | 3 | — |  | ⛔ |
| `LIMITVALUE` | 2 | — | ⭐ | ⛔ |
| `ListInstallment` | 14 | 13 | ⭐ | ✅ |
| `LossPrevention` | 1 | — | ⭐ | ⛔ |
| `MarketingCode` | 1 | 1 |  | ✅ |
| `MarketingName` | 1 | 2 |  | ✅ |
| `MarketingOfficer` | 4 | 5 | ⭐ | ✅ |
| `MasterID` | 2 | 1 |  | ⛔ |
| `MDP` | 5 | — | ⭐ | ⛔ |
| `MDP2` | 4 | — | ⭐ | ⛔ |
| `MDPCURRENCY` | 2 | — | ⭐ | ⛔ |
| `MDPList` | 1 | 1 |  | ⛔ |
| `MDPVALUE` | 4 | — | ⭐ | ⛔ |
| `MinMax` | 2 | — |  | ⛔ |
| `MOID` | 1 | 1 |  | ✅ |
| `NAME` | 1 | — | ⭐ | ⛔ |
| `Name` | 8 | — | ⭐ | ⛔ |
| `NBStatus` | 2 | 2 |  | ⛔ |
| `NetPremi` | 9 | 14 | ⭐ | ✅ |
| `NetPremi2` | 6 | — | ⭐ | ⛔ |
| `NetPremiAfterPPH` | 3 | 9 | ⭐ | ✅ |
| `NetPremiAfterPPH2` | 2 | — | ⭐ | ⛔ |
| `NetPremiAfterPPN` | 3 | 9 | ⭐ | ✅ |
| `NetPremiAfterPPN2` | 2 | — | ⭐ | ⛔ |
| `NetPremiAfterTax` | 1 | 9 | ⭐ | ✅ |
| `NETPREMICURRENCY` | 2 | — | ⭐ | ⛔ |
| `NetPremium` | 10 | 9 | ⭐ | ✅ |
| `NetPremiumList` | 6 | 7 |  | ⛔ |
| `NETPREMIVALUE` | 2 | — | ⭐ | ⛔ |
| `NoOffer` | 12 | 15 | ⭐ | ✅ |
| `NoOfferSlip` | 1 | — |  | ⛔ |
| `NOTE` | 5 | — |  | ⛔ |
| `Note` | 8 | 6 | ⭐ | ⛔ |
| `NusareLimit` | 1 | — |  | ⛔ |
| `ObjectIndex` | 1 | — |  | ⛔ |
| `OccupationId` | 2 | — |  | ⛔ |
| `OfferFacIn` | 1 | — |  | ⛔ |
| `OJKBusinessID` | 2 | 2 |  | ✅ |
| `OldData` | 1 | 29 | ⭐ | ✅ |
| `OldID` | 7 | — |  | ⛔ |
| `OLDID` | 3 | 1 | ⭐ | ⛔ |
| `OldPolicyNo` | — | 1 |  | ✅ |
| `OldTSI` | 1 | — |  | ⛔ |
| `OperatorName` | 8 | 8 | ⭐ | ✅ |
| `OutstandingClaim` | 4 | 4 | ⭐ | ⛔ |
| `OveriddingCommOgp` | 8 | 10 | ⭐ | ⛔ |
| `OveriddingCommOnp` | 8 | 10 | ⭐ | ⛔ |
| `PaymentDate` | 6 | 3 | ⭐ | ⛔ |
| `PaymentTotal` | 13 | 13 | ⭐ | ✅ |
| `PaymentTotalAfterPPN` | 1 | — |  | ✅ |
| `PaymentTotalAfterTax` | 1 | — |  | ✅ |
| `Pct` | 3 | — |  | ⛔ |
| `PctTotal` | 4 | 2 | ⭐ | ⛔ |
| `PICSuggest` | 1 | — |  | ⛔ |
| `PNOTE` | 1 | — |  | ⛔ |
| `Policy` | 3 | 2 |  | ⛔ |
| `PolicyData` | 1 | — |  | ⛔ |
| `PolicyNo` | 6 | 14 | ⭐ | ✅ |
| `PolicyTreatyIn` | 16 | 11 | ⭐ | ⛔ |
| `PolicyTreatyInDetail` | 1 | 1 |  | ⛔ |
| `Position` | 1 | 1 |  | ⛔ |
| `PositionNote` | 1 | 2 | ⭐ | ⛔ |
| `PPh` | 1 | — |  | ✅ |
| `PPHValue` | 10 | 13 | ⭐ | ✅ |
| `PPHValue2` | 1 | — |  | ⛔ |
| `PPN` | 1 | — |  | ✅ |
| `PPNValue` | 10 | 13 | ⭐ | ✅ |
| `PPNValue2` | 1 | — |  | ⛔ |
| `PremiNusantaraRe` | 5 | — |  | ⛔ |
| `PremiOgp` | 17 | 16 | ⭐ | ✅ |
| `PremiOnp` | 11 | 13 | ⭐ | ✅ |
| `Premium` | 23 | 16 | ⭐ | ✅ |
| `PremiumAfterPPH` | — | 1 |  | ✅ |
| `PremiumAfterPPN` | 2 | 2 | ⭐ | ✅ |
| `PremiumAfterTax` | 2 | 2 | ⭐ | ✅ |
| `PremiumEarned` | 1 | — |  | ⛔ |
| `PremiumSpreaded` | 24 | 11 | ⭐ | ✅ |
| `ProdKe` | — | 1 |  | ✅ |
| `ProductionDate` | 6 | 6 | ⭐ | ✅ |
| `Property` | 11 | — |  | ⛔ |
| `property` | 85 | 64 |  | ⛔ |
| `Proportion` | 1 | — |  | ⛔ |
| `ProportionalType` | 2 | 5 | ⭐ | ✅ |
| `ProportionType` | 1 | 2 | ⭐ | ⛔ |
| `PROPORTIONTYPE` | 4 | — | ⭐ | ⛔ |
| `PROVINCENAME` | 1 | — |  | ⛔ |
| `pyBestBet` | 3 | — | ⭐ | ⛔ |
| `pyCells` | 5 | 8 | ⭐ | ⛔ |
| `pyCustomerID` | 4 | 5 | ⭐ | ⛔ |
| `pyExpanded` | 3 | 4 |  | ✅ |
| `pyID` | 1 | 1 | ⭐ | ⛔ |
| `pyMessageLabel` | 1 | — |  | ⛔ |
| `pySectionBody` | 5 | 8 | ⭐ | ⛔ |
| `pySections` | 5 | 8 | ⭐ | ⛔ |
| `pySteps` | 56 | 31 |  | ⛔ |
| `pySummaryCount` | 1 | — | ⭐ | ⛔ |
| `pySummaryDateTime` | 1 | — | ⭐ | ⛔ |
| `pySummaryValue` | 1 | — | ⭐ | ⛔ |
| `pyTable` | 5 | 8 | ⭐ | ⛔ |
| `pyTemplateAutoComplete` | 2 | — | ⭐ | ⛔ |
| `pyTemplateButton` | 11 | 9 | ⭐ | ⛔ |
| `pyTemplateCalendar` | 5 | 6 | ⭐ | ⛔ |
| `pyTemplateCheckbox` | 4 | 1 | ⭐ | ⛔ |
| `pyTemplateDisplayText` | 4 | — | ⭐ | ⛔ |
| `pyTemplateInputBox` | 24 | 14 | ⭐ | ⛔ |
| `pyTemplateRadioButton` | 5 | 2 | ⭐ | ⛔ |
| `pyTemplateSelect` | 5 | 9 | ⭐ | ⛔ |
| `pyWorkBasketName` | 2 | 1 | ⭐ | ⛔ |
| `pzInsKey` | 2 | 2 | ⭐ | ⛔ |
| `Quartal` | 4 | 4 | ⭐ | ⛔ |
| `Quotation` | 2 | 2 |  | ⛔ |
| `QuotationData` | 23 | 17 | ⭐ | ✅ |
| `QuotationList` | 2 | — | ⭐ | ⛔ |
| `ReinsName` | 2 | 2 |  | ⛔ |
| `Reinstatement_List` | 2 | 2 |  | ⛔ |
| `ReinstatementAmount1` | 1 | 1 |  | ⛔ |
| `ReinstatementAmount2` | 1 | 1 |  | ⛔ |
| `ReinstatementNote` | 1 | 1 |  | ⛔ |
| `ReinstatementPct` | 1 | 1 |  | ⛔ |
| `ReinstatementValue` | 2 | 2 |  | ⛔ |
| `ReinsTypeID` | 1 | — |  | ⛔ |
| `ReinsTypeName` | 1 | — |  | ⛔ |
| `ReinsuranceListTONP` | 3 | 2 |  | ⛔ |
| `Remark` | 4 | 1 | ⭐ | ✅ |
| `Remarks` | 4 | — | ⭐ | ⛔ |
| `ResultOgp1` | 9 | 11 | ⭐ | ✅ |
| `ResultOgp2` | 7 | 9 | ⭐ | ✅ |
| `ResultOnp1` | 8 | 10 | ⭐ | ✅ |
| `ResultOnp2` | 7 | 9 | ⭐ | ✅ |
| `RETENTIONCURRENCY` | 2 | — | ⭐ | ⛔ |
| `RETENTIONVALUE` | 2 | — | ⭐ | ⛔ |
| `RetroPct` | 2 | — |  | ⛔ |
| `RetroTypeID` | 2 | — |  | ⛔ |
| `RICommision` | 1 | — |  | ⛔ |
| `RiCommOgp` | 9 | 9 | ⭐ | ⛔ |
| `RiCommOnp` | 7 | 9 | ⭐ | ⛔ |
| `RIOGR` | 1 | — |  | ⛔ |
| `RIONR` | 1 | — |  | ⛔ |
| `RnmLimitList` | 3 | 2 |  | ⛔ |
| `SalvageValue` | 7 | 9 | ⭐ | ✅ |
| `ShareCeding` | 1 | — |  | ⛔ |
| `ShareCurrency` | 6 | — | ⭐ | ✅ |
| `SHARECURRENCY` | 4 | — | ⭐ | ⛔ |
| `SharePercentage` | 25 | 12 | ⭐ | ✅ |
| `ShareValue` | 9 | 3 | ⭐ | ✅ |
| `SHAREVALUE` | 4 | — | ⭐ | ⛔ |
| `Show` | 1 | 3 |  | ✅ |
| `SOB` | 4 | 2 | ⭐ | ✅ |
| `SOBID` | 1 | — |  | ⛔ |
| `SOBName` | 8 | 7 | ⭐ | ✅ |
| `SobName` | 1 | 1 |  | ✅ |
| `SourceOfBusiness` | 2 | 2 |  | ✅ |
| `SplitRNMSharePct` | 1 | 1 |  | ⛔ |
| `SpreadingList` | 14 | — |  | ⛔ |
| `SpreadingListXOL` | 3 | — |  | ⛔ |
| `SpreadingRiskList` | 15 | 10 | ⭐ | ✅ |
| `SpreadingTotalPct` | 1 | — |  | ⛔ |
| `SpreadingType` | 1 | — |  | ⛔ |
| `SpreadingTypeID` | 1 | — |  | ⛔ |
| `SpreadingTypeIDXOL` | 4 | 3 |  | ⛔ |
| `SpreadingTypeXOL` | 1 | 1 |  | ⛔ |
| `StartDate` | 11 | 7 | ⭐ | ✅ |
| `StatementDate` | 7 | 8 | ⭐ | ✅ |
| `StatementType` | 4 | 1 | ⭐ | ⛔ |
| `STS_PKP` | 1 | 1 |  | ⛔ |
| `Suggest` | 8 | 8 | ⭐ | ✅ |
| `SuggestDate` | 4 | 5 |  | ⛔ |
| `SuggestList` | 4 | 3 | ⭐ | ✅ |
| `SumTotalPayment` | 2 | 1 |  | ⛔ |
| `SurveyedBy` | 4 | — | ⭐ | ⛔ |
| `TANGGAL` | 1 | 3 |  | ⛔ |
| `TeamGroup` | 1 | 1 |  | ✅ |
| `TERITORYNAME` | 1 | — |  | ⛔ |
| `Termination` | 5 | 2 | ⭐ | ⛔ |
| `TextNoQuotation` | 1 | — | ⭐ | ⛔ |
| `TimeExcess` | 2 | — |  | ⛔ |
| `TotalClaim` | 7 | 9 | ⭐ | ✅ |
| `TotalGrossPremi` | 1 | — |  | ⛔ |
| `TotalNet` | 2 | — |  | ⛔ |
| `TotalNetPremiAfterPPH` | 1 | — |  | ⛔ |
| `TotalNetPremiAfterPPN` | 2 | — | ⭐ | ⛔ |
| `TotalNetPremiAfterTax` | 2 | — | ⭐ | ⛔ |
| `TotalNetRate` | 5 | — |  | ⛔ |
| `TotalPPHValue` | 2 | — | ⭐ | ⛔ |
| `TotalPPNValue` | 2 | — | ⭐ | ⛔ |
| `TotalPremium` | 8 | 11 | ⭐ | ✅ |
| `TotalSharePercentageClaim` | 6 | 8 | ⭐ | ✅ |
| `TotalSharePercentagePremium` | 6 | 8 | ⭐ | ✅ |
| `TotalSpread` | 1 | — |  | ⛔ |
| `TotalTSIPremiSpreadRNM` | 1 | — |  | ⛔ |
| `TreatyContractName` | — | 1 | ⭐ | ⛔ |
| `TREATYCONTRACTNAME` | 2 | — | ⭐ | ⛔ |
| `TreatyDifference` | 1 | 12 | ⭐ | ✅ |
| `TreatyGroup` | 1 | 1 |  | ⛔ |
| `TREATYGROUP` | 4 | — | ⭐ | ⛔ |
| `TreatyGroupID` | 5 | 3 |  | ✅ |
| `TREATYGROUPID` | 2 | — |  | ⛔ |
| `TreatyGroupList` | — | 1 |  | ⛔ |
| `TreatyGroupName` | 7 | 4 | ⭐ | ✅ |
| `TreatyGroupOldID` | 4 | 2 |  | ✅ |
| `TREATYID` | 4 | — | ⭐ | ⛔ |
| `TreatyIn` | 3 | — |  | ⛔ |
| `TreatyName` | 13 | 5 | ⭐ | ⛔ |
| `TreatyType` | 28 | 15 | ⭐ | ✅ |
| `TREATYTYPE` | 4 | — | ⭐ | ⛔ |
| `TreatyTypeID` | 1 | — |  | ⛔ |
| `TreatyXOLDifferenceList` | 2 | 7 | ⭐ | ✅ |
| `TreatyXOLList` | 4 | 12 | ⭐ | ✅ |
| `TreatyYear` | 9 | 5 | ⭐ | ✅ |
| `TREATYYEAR` | 4 | — | ⭐ | ⛔ |
| `TSI` | 11 | — |  | ⛔ |
| `TSIAddCap` | 1 | — |  | ⛔ |
| `TSIGrossSpreaded` | 1 | — |  | ⛔ |
| `TSINusantaraRe` | 5 | — |  | ⛔ |
| `TSIObjectItem` | 1 | — |  | ⛔ |
| `TSISpreaded` | 13 | — |  | ⛔ |
| `Type` | — | 1 | ⭐ | ⛔ |
| `TypeDeductible` | 2 | — |  | ⛔ |
| `TypeTax` | 11 | 9 | ⭐ | ✅ |
| `URL` | — | 1 |  | ⛔ |
| `Value` | 12 | 8 | ⭐ | ⛔ |
| `ValueList` | 3 | 10 | ⭐ | ✅ |
| `ViewState` | — | 2 |  | ✅ |
| `WPC` | 1 | — | ⭐ | ⛔ |
| `YearOfQuartal` | 4 | 3 | ⭐ | ⛔ |
| `ZipCode` | 1 | — |  | ⛔ |

---

## Medan yang hanya ada di data guide

Tidak tersapu dari aturan, tetapi ada pada dokumen tersimpan. Sebagian besar adalah medan yang
**ditulis sistem**, bukan diketik atau dibaca aturan — karena itu tidak muncul di Activity
maupun Section.

`BranchCode` · `BranchName` · `BusinessFac` · `BusinessType2` · `IsSurveyReport` · `OperatorID` · `pxListSubscript` · `pxObjClass` · `SobLsg` · `StatusBusiness` · `StatusSyariah` · `SurveyReportList` · `TypeFacultative`

⚠️ Medan di daftar ini **tetap dimigrasi**. Ketiadaannya di aturan berarti tidak ada layar yang
memakainya, bukan berarti isinya tidak bernilai.

---

## Batas berkas ini

⛔ Ia mendaftar **nama medan**, bukan tipe dan bukan panjang. Panjang maksimum per medan masih
diambil dari data guide — satu-satunya kegunaannya yang tetap sah.

⛔ Cacah di sini adalah **batas bawah**, bukan total. Pemecah dokumen tetap wajib menyediakan
penampung medan tak dikenal, dan penampung itu wajib **kosong** sebelum pekerjaan dinyatakan
selesai.

⛔ Nol nilai data produksi disalin ke berkas ini. Hanya nama medan dan cacah kemunculan.

---

*Disusun 23 September 2026 sesudah work owner memutuskan data guide tidak dipakai sebagai dasar.*
