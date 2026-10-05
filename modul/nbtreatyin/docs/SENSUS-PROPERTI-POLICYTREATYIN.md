# Sensus properti `PolicyTreatyIn` dan `Quotation` — dari XML yang terjangkau

> Bangkitan `docs/alat/sensus.py` atas korpus.json (rule TERJANGKAU saja, `graf.py`). Tiket 00: *"Cocokkan ulang ke XML (termasuk DataTransform, When, dan DecisionTable)"*.

- Skalar `PolicyTreatyIn`: **153**
- Daftar (PageList) `PolicyTreatyIn`: **7**
- Properti `Quotation` / `QuotationData`: **23**
- Tampil di layar (sel berhalaman PolicyTreatyIn): **120**

## Skalar `PolicyTreatyIn`

| Properti | Layar | Dirujuk (rule) |
| --- | :---: | --- |
| `BalanceBeforePPH` | ⭐ | 4: CountNetPremi_act, InputPolicyTreatyInDetail_NonProp, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `BalanceBeforeTax` | ⭐ | 4: CountNetPremi_act, InputPolicyTreatyInDetail_NonProp, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `BalanceDueTo` | ⭐ | 13: CountNetPremi_act, CountPctInstallment_Act, FillPaymentInstallment, InputPolicyTreatyEDMDetail_NP, InputPolicyTreatyInDetail_NonProp, InputPolicyTreatyInDetail_preACT … |
| `BizCode` |  | 3: InputPolicyTreatyInDetail_preACT, InputPolicyTreatyInPre_Act, InputPolicyTreatyOutDetail_preACT |
| `BizName` | ⭐ | 4: InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `BranchDetailID` |  | 1: CheckDataMkt |
| `BranchDetailName` |  | 1: CheckDataMkt |
| `BrokerageFee` |  | 1: SetPPNPPH |
| `BrokerageFeeSebenarnya` |  | 1: SetPPNPPH |
| `CEDING` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `CedingCo` |  | 2: InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT |
| `CedingCoName` | ⭐ | 5: InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT, TreatyRealizationCheckDuplicate, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `Claim` | ⭐ | 7: CalculatePremi_Act, CountNetPremi_act, CountOGPONP_Act, CountSpreading_Act, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn … |
| `ClaimPaymentType` | ⭐ | 4: CountOGPONP_Act, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn, isClaimTreaty |
| `ClaimPercentage` | ⭐ | 3: DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn, SpreadingRiskList |
| `ClaimSpreaded` | ⭐ | 3: DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn, SpreadingRiskList |
| `ClaimType` | ⭐ | 8: CountOGPONP_Act, GeneratePolicyNoTreaty_Act, InputPolicyTreatyEDMDetail_NP, TreatyNonPropSetSpreading, TreatyRealizationCheckDuplicate, DetailDeptHeadTreatyIn_UW … |
| `CLASSOFBUSINESS` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `ClientID` |  | 1: CheckDataMkt |
| `ClientName` |  | 1: CheckDataMkt |
| `Currency` | ⭐ | 15: CekLimitTreatyAcc_Act, InputPolicyTreatyEDMDetail_NP, InputPolicyTreatyInDetail_NonProp, InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_NonProp, InputPolicyTreatyOutDetail_preACT … |
| `Date` | ⭐ | 1: ListSuggest |
| `DateofSurvey` | ⭐ | 2: HistoricalSurveyReportDtl, HistoricalSurveyReportDtlUW |
| `Deductible` | ⭐ | 3: DetailPolicyTreatyInNonProportional, DetailPolicyTreatyInNonProportionalEDM, DetailPolicyTreatyOutNonProportional |
| `Deductible2` | ⭐ | 3: DetailPolicyTreatyInNonProportional, DetailPolicyTreatyInNonProportionalEDM, DetailPolicyTreatyOutNonProportional |
| `Deduction1` | ⭐ | 9: CountNetPremi_act, InputPolicyTreatyEDMDetail_NP, InputPolicyTreatyInDetail_NonProp, InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_NonProp, InputPolicyTreatyOutDetail_preACT … |
| `Deduction2` | ⭐ | 8: CountNetPremi_act, InputPolicyTreatyEDMDetail_NP, InputPolicyTreatyInDetail_NonProp, InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_NonProp, InputPolicyTreatyOutDetail_preACT … |
| `DueDate` | ⭐ | 2: DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `DueTo` | ⭐ | 7: GeneratePolicyNoTreaty_Act, InputPolicyTreatyEDMDetail_NP, InputPolicyTreatyInDetail_NonProp, InputPolicyTreatyOutDetail_NonProp, SetDueTo_act, DetailDeptHeadTreatyIn_UW … |
| `EDMType` |  | 1: InputPolicyTreatyInPre_Act |
| `EndDate` | ⭐ | 9: InputPolicyTreatyEDMDetail_NP, InputPolicyTreatyInDetail_NonProp, InputPolicyTreatyOutDetail_NonProp, ProtectDate, TreatyRealizationCheckDuplicate, InputPolicyTreatyIn_preDT … |
| `EPICURRENCY` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `EPIVALUE` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `ExcessLoss` | ⭐ | 5: CountNetPremi_act, CountSpreading_Act, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn, isClaimTreaty |
| `FlagPPH` | ⭐ | 9: InputPolicyTreatyInDetail_NonProp, InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_NonProp, InsertToTreatyXOLList, InsertToTreatyXOLListRetroShare, RemoveTypeTax_ACT … |
| `FlagRetroTreaty` | ⭐ | 6: InputPolicyTreatyInDetail_NonProp, InputPolicyTreatyOutDetail_NonProp, InsertToTreatyXOLListRetroShare, TreatyNonPropOutSetSpreading, TreatyNonPropSetSpreading, DetailPolicyTreatyIn |
| `GrossClaim` | ⭐ | 3: CalculatePremi_Act, CountNetPremi_act, DetailPolicyTreatyIn |
| `GrossPremium` | ⭐ | 4: CalculatePremi_Act, CountResult1_Act, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `HasFacOut` |  | 2: InboxPolicyTreatyIn_postDT, TestTreatyToFacStatus |
| `ID` |  | 1: CheckDataMkt |
| `IDCurrency` | ⭐ | 7: CekLimitTreatyAcc_Act, InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT, SetTreatyCurrencyID, TreatyNonPropOutSetSpreading, TreatyNonPropSetSpreading … |
| `Installment` | ⭐ | 7: FillPaymentInstallment, InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn, DetailPolicyTreatyInNonProportional … |
| `InstallmentNo` | ⭐ | 2: DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `InstallmentPercentage` | ⭐ | 2: DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `InsuredID` |  | 2: InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT |
| `InsuredName` | ⭐ | 6: InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn, HistoricalSurveyReportDtl, HistoricalSurveyReportDtlUW |
| `IsApproved` | ⭐ | 9: InputPolicyTreatyInPost_Act, InsertHistoryAkseptasiPega, DeptHeadTreatyInUW_preDT, DeptHeadTreatyIn_UW_postDT, InboxPolicyTreatyIn_postDT, isApproved … |
| `isApprovedtoDeptHead` |  | 1: DeptHeadTreatyInUW_preDT |
| `IsEDMInputOnNB` |  | 2: InputPolicyTreatyInDetail_NonProp, InputPolicyTreatyOutDetail_NonProp |
| `IsNewPolicyNonProp` |  | 6: InputPolicyTreatyInPre_Act, InsertToTreatyXOLListRetroShare, InputPolicyTreatyIn_preDT, TreatyEnableDisableInput, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `IsOJKNopolis` |  | 1: GeneratePolicyNoTreaty_Act |
| `Layer` | ⭐ | 4: InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `LAYER` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `LayerPart` | ⭐ | 4: InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `LAYERPART` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `LayerPartType` | ⭐ | 4: InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `LAYERPARTTYPE` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `LayerType` | ⭐ | 4: InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `LAYERTYPE` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `Limit` | ⭐ | 2: DetailPolicyTreatyInNonProportional, DetailPolicyTreatyOutNonProportional |
| `Limit2` | ⭐ | 2: DetailPolicyTreatyInNonProportional, DetailPolicyTreatyOutNonProportional |
| `LIMITCURRENCY` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `LIMITVALUE` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `LossPrevention` | ⭐ | 1: HistoricalSurveyReportDtl |
| `MarketingCode` |  | 1: CheckDataMkt |
| `MarketingName` |  | 1: CheckDataMkt |
| `MarketingOfficer` | ⭐ | 3: CheckDataMkt, InputPolicyTreatyIn_preDT, DetailDeptHeadTreatyIn_UW |
| `MasterID` |  | 2: DeptHeadTreatyInUW_preDT, InputPolicyTreatyIn_preDT |
| `MDP` | ⭐ | 3: DetailPolicyTreatyInNonProportional, DetailPolicyTreatyInNonProportionalEDM, DetailPolicyTreatyOutNonProportional |
| `MDP2` | ⭐ | 3: DetailPolicyTreatyInNonProportional, DetailPolicyTreatyInNonProportionalEDM, DetailPolicyTreatyOutNonProportional |
| `MDPCURRENCY` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `MDPVALUE` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `NetPremi` | ⭐ | 3: DetailPolicyTreatyInNonProportional, DetailPolicyTreatyInNonProportionalEDM, DetailPolicyTreatyOutNonProportional |
| `NetPremi2` | ⭐ | 3: DetailPolicyTreatyInNonProportional, DetailPolicyTreatyInNonProportionalEDM, DetailPolicyTreatyOutNonProportional |
| `NetPremiAfterPPH` | ⭐ | 1: DetailPolicyTreatyInNonProportional |
| `NetPremiAfterPPH2` | ⭐ | 1: DetailPolicyTreatyInNonProportional |
| `NetPremiAfterPPN` | ⭐ | 1: DetailPolicyTreatyInNonProportional |
| `NetPremiAfterPPN2` | ⭐ | 1: DetailPolicyTreatyInNonProportional |
| `NETPREMICURRENCY` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `NetPremium` | ⭐ | 7: CountNetPremi_act, CountSpreading_Act, InputPolicyTreatyEDMDetail_NP, InputPolicyTreatyInDetail_NonProp, InputPolicyTreatyOutDetail_NonProp, DetailDeptHeadTreatyIn_UW … |
| `NETPREMIVALUE` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `NoOffer` | ⭐ | 9: InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT, TreatyRealizationCheckDuplicate, TreatyRealizationCheckXOLList, BrowseTreatyInJoinEDM, BrowseTreatyOut … |
| `Note` |  | 3: DetailPolicyTreatyInNonProportional, DetailPolicyTreatyInNonProportionalEDM, DetailPolicyTreatyOutNonProportional |
| `OJKBusinessID` |  | 2: FetchTreatyGroupOJK, GeneratePolicyNoTreaty_Act |
| `OldData` |  | 1: TreatyRealizationCheckXOLList |
| `OperatorName` | ⭐ | 3: DeptHeadTreatyInUW_preDT, InputPolicyTreatyIn_preDT, ListSuggest |
| `OutstandingClaim` | ⭐ | 2: DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `OveriddingCommOgp` | ⭐ | 6: CountNetPremi_act, CountOGPONP_Act, CountOverridingCommOgp_Act, CountResult2Ogp_act, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `OveriddingCommOnp` | ⭐ | 6: CountNetPremi_act, CountOGPONP_Act, CountOverridingCommOnp_Act, CountResult2Onp_act, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `PaymentTotal` | ⭐ | 2: DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `Policy` |  | 1: SaveJsonPolisTreatyIn_Act |
| `PolicyNo` |  | 6: GeneratePolicyNoTreaty_Act, DeptHeadTreatyIn_UW_postDT, SavePolisTreatyIn_SQL, DetailDeptHeadTreatyIn_UW, ShowPolicyNoTreaty_SC, NopolisEmpty |
| `PPHValue` | ⭐ | 5: CountNetPremi_act, InputPolicyTreatyInDetail_NonProp, SetPPNPPH, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `PPNValue` | ⭐ | 5: CountNetPremi_act, InputPolicyTreatyInDetail_NonProp, SetPPNPPH, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `PremiOgp` | ⭐ | 15: CalculatePremi_Act, CountNetPremi_act, CountOGPONP_Act, CountOverridingCommOgp_Act, CountResult1_Act, CountResult2Ogp_act … |
| `PremiOnp` | ⭐ | 9: CountNetPremi_act, CountOGPONP_Act, CountOverridingCommOnp_Act, CountResult1Onp_Act, CountResult2Onp_act, CountRiCommOnp_act … |
| `Premium` | ⭐ | 2: DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `PremiumSpreaded` | ⭐ | 3: DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn, SpreadingRiskList |
| `ProductionDate` | ⭐ | 4: GeneratePolicyNoTreaty_Act, InputPolicyTreatyInPre_Act, DetailDeptHeadTreatyIn_UW, ListSuggest |
| `PROPORTIONTYPE` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `Quartal` | ⭐ | 2: DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `Remark` | ⭐ | 2: DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `Remarks` | ⭐ | 2: HistoricalSurveyReportDtl, HistoricalSurveyReportDtlUW |
| `ResultOgp1` | ⭐ | 7: CalculatePremi_Act, CountNetPremi_act, CountOGPONP_Act, CountResult1_Act, CountRiCommOgp_act, DetailDeptHeadTreatyIn_UW … |
| `ResultOgp2` | ⭐ | 5: CountOGPONP_Act, CountOverridingCommOgp_Act, CountResult2Ogp_act, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `ResultOnp1` | ⭐ | 6: CountNetPremi_act, CountOGPONP_Act, CountResult1Onp_Act, CountRiCommOnp_act, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `ResultOnp2` | ⭐ | 5: CountOGPONP_Act, CountOverridingCommOnp_Act, CountResult2Onp_act, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `RETENTIONCURRENCY` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `RETENTIONVALUE` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `RiCommOgp` | ⭐ | 7: CalculatePremi_Act, CountOGPONP_Act, CountResult1_Act, CountRiCommOgp_act, TreatyInputPctCommSpreading, DetailDeptHeadTreatyIn_UW … |
| `RiCommOnp` | ⭐ | 5: CountOGPONP_Act, CountResult1Onp_Act, CountRiCommOnp_act, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `SalvageValue` | ⭐ | 6: CountNetPremi_act, CountOGPONP_Act, CountSpreading_Act, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn, isClaimTreaty |
| `ShareCurrency` | ⭐ | 4: InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `SHARECURRENCY` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `SharePercentage` | ⭐ | 3: DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn, SpreadingRiskList |
| `ShareValue` | ⭐ | 7: InputPolicyTreatyEDMDetail_NP, InputPolicyTreatyInDetail_NonProp, InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_NonProp, InputPolicyTreatyOutDetail_preACT, DetailDeptHeadTreatyIn_UW … |
| `SHAREVALUE` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `Show` |  | 1: InputRealizationTreatyIn |
| `SOB` | ⭐ | 4: InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT, BusinessAndSOBList, BusinessAndSOBListRetro |
| `SobName` |  | 1: CheckDataMkt |
| `SOBName` | ⭐ | 6: InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_NonProp, InputPolicyTreatyOutDetail_preACT, InsertToTreatyOutXOLList, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `SourceOfBusiness` |  | 1: CheckDataMkt |
| `StartDate` | ⭐ | 9: InputPolicyTreatyEDMDetail_NP, InputPolicyTreatyInDetail_NonProp, InputPolicyTreatyOutDetail_NonProp, ProtectDate, TreatyRealizationCheckDuplicate, InputPolicyTreatyIn_preDT … |
| `StatementDate` | ⭐ | 5: GeneratePolicyNoTreaty_Act, InputPolicyTreatyInPre_Act, InputPolicyTreatyIn_preDT, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `StatementType` | ⭐ | 2: DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `Suggest` | ⭐ | 3: DeptHeadTreatyInUW_preDT, InputPolicyTreatyIn_preDT, ListSuggest |
| `SuggestDate` |  | 2: DeptHeadTreatyInUW_preDT, InputPolicyTreatyIn_preDT |
| `SurveyedBy` |  | 2: HistoricalSurveyReportDtl, HistoricalSurveyReportDtlUW |
| `T` |  | 1: GeneratePolicyNoTreaty_Act |
| `TeamGroup` |  | 1: CheckDataMkt |
| `TotalClaim` | ⭐ | 5: BreakDownSpreading_Act, CountSpreading_Act, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn, SpreadingRiskList |
| `TotalNetPremiAfterPPN` | ⭐ | 1: DetailPolicyTreatyInNonProportional |
| `TotalNetPremiAfterTax` | ⭐ | 1: DetailPolicyTreatyInNonProportional |
| `TotalPPHValue` | ⭐ | 1: DetailPolicyTreatyInNonProportional |
| `TotalPPNValue` | ⭐ | 1: DetailPolicyTreatyInNonProportional |
| `TotalPremium` | ⭐ | 6: BreakDownSpreading_Act, CekLimitTreatyAcc_Act, CountSpreading_Act, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn, SpreadingRiskList |
| `TotalSharePercentageClaim` | ⭐ | 4: CountSpreading_Act, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn, SpreadingRiskList |
| `TotalSharePercentagePremium` | ⭐ | 4: CountSpreading_Act, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn, SpreadingRiskList |
| `TREATYCONTRACTNAME` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `TreatyDifference` |  | 1: CekLimitTreatyAcc_Act |
| `TREATYGROUP` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `TreatyGroupID` |  | 5: BreakDownSpreading_Act, FetchTreatyGroupOJK, FetchTreatyGroupOldID, InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT |
| `TreatyGroupName` | ⭐ | 5: InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT, TreatyInputPctCommSpreading, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `TreatyGroupOldID` |  | 4: FetchTreatyGroupOldID, GeneratePolicyNoTreaty_Act, InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT |
| `TREATYID` |  | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `TreatyIn` |  | 1: TreatyInputPctCommSpreading |
| `TreatyType` | ⭐ | 6: InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT, TreatyInputPctCommSpreading, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn, SpreadingRiskList |
| `TREATYTYPE` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `TreatyYear` | ⭐ | 6: BreakDownSpreading_Act, InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_preACT, TreatyRealizationCheckDuplicate, DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |
| `TREATYYEAR` | ⭐ | 2: BusinessAndSOBList, BusinessAndSOBListRetro |
| `TypeTax` | ⭐ | 9: InputPolicyTreatyInDetail_NonProp, InputPolicyTreatyInDetail_preACT, InputPolicyTreatyOutDetail_NonProp, InsertToTreatyOutXOLList, InsertToTreatyXOLList, RemoveTypeTax_ACT … |
| `Value` | ⭐ | 3: DetailPolicyTreatyInNonProportional, DetailPolicyTreatyInNonProportionalEDM, DetailPolicyTreatyOutNonProportional |
| `YearOfQuartal` | ⭐ | 2: DetailDeptHeadTreatyIn_UW, DetailPolicyTreatyIn |

## Daftar `PolicyTreatyIn` dan anggotanya

- `BreakDownSpreadList` — 1 rule; anggota: —
- `ListInstallment` — 8 rule; anggota: `Currency`, `DueDate`, `IDCurrency`, `InstallmentList`, `InstallmentNo`, `InstallmentPercentage`, `PaymentTotal`, `Premium`
- `pxResults` — 1 rule; anggota: —
- `SpreadingRiskList` — 7 rule; anggota: `ClaimPercentage`, `ClaimSpreaded`, `Currency`, `CurrencyID`, `PremiumSpreaded`, `SharePercentage`, `SplitRNMSharePct`, `TreatyName`, `TreatyType`
- `SuggestList` — 1 rule; anggota: `IsSave`
- `TreatyXOLDifferenceList` — 2 rule; anggota: —
- `TreatyXOLList` — 3 rule; anggota: `BrokerageFeeSebenarnya`, `Currency`, `Deduction`, `DueTo`, `DueToValue`, `GrossPremi`, `IDCurrency`, `NetPremi`, `NetPremiAfterPPH`, `NetPremiAfterPPN`, `NetPremiAfterTax`, `PPHValue`, `PPNValue`, `ValueList`

## `Quotation` / `QuotationData`

`BranchCode`, `BranchName`, `btnQuotation`, `BusinessCode`, `BusinessFac`, `BusinessName`, `BusinessOldId`, `BusinessType`, `CedingCo`, `CedingCoName`, `GroupPanel`, `InsuredID`, `InsuredName`, `MarketingCode`, `MarketingName`, `MOID`, `NoOfferSlip`, `ProportionalType`, `SobLeader0`, `SobLeader1`, `SobName`, `SourceOfBusiness`, `TeamGroup`

