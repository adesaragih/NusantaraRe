# Turunan kolom 28 anak `TREATY_IN` yang belum dibangun

Disapu dari `POOLDATA.M_TREATY_IN.JSONDATA`, **1.854 dokumen**, 6 Oktober 2026.
Jalurnya diambil dari `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx`.

⛔ `pxObjClass`, `pxListSubscript`, dan `pyTemplate*` adalah PERABOT Pega — tidak
menjadi kolom. Keputusan yang sama sudah berlaku pada migrasi `436`.

## `T_TREATY_CURRENCY`

Jalur `TreatyIn.CurrencyList` · **3465 elemen** · 5 kolom calon

| kolom | ada pada |
| --- | ---: |
| `Conversion` | 3465 |
| `Currency` | 3465 |
| `PeriodStart` | 3461 |
| `PeriodEnd` | 3460 |
| `CurrencyID` | 3357 |

## `T_TREATY_FAC_LIMITS`

Jalur `TreatyIn.ShareFacultativeReinsurers.FacultativeLimits` · **0 elemen** · 0 kolom calon

⛔ **NOL elemen di 1.854 dokumen.** Jalur ini tidak pernah berisi.

## `T_TREATY_FAC_LIMIT_DETAIL`

Jalur `TreatyIn.ShareFacultativeReinsurers.FacultativeLimits.Detail` · **0 elemen** · 0 kolom calon

⛔ **NOL elemen di 1.854 dokumen.** Jalur ini tidak pernah berisi.

## `T_TREATY_FAC_REINSURER`

Jalur `TreatyIn.ShareFacultativeReinsurers` · **15 elemen** · 4 kolom calon

| kolom | ada pada |
| --- | ---: |
| `ReinsID` | 15 |
| `ReinsName` | 15 |
| `SharePct` | 15 |
| `Layer` | 13 |

## `T_TREATY_FAC_SHARE`

Jalur `TreatyIn.FacultativeShareList` · **12 elemen** · 31 kolom calon

| kolom | ada pada |
| --- | ---: |
| `ClassofBusinessList` | 12 |
| `Cover` | 12 |
| `DeductionList` | 12 |
| `DeductionTotalList` | 12 |
| `GrossPremiumList` | 12 |
| `Layer` | 12 |
| `LayerPart` | 12 |
| `LayerPartType` | 12 |
| `LayerType` | 12 |
| `NetPremiumList` | 12 |
| `RNMSpreadedListDeductRIXOL` | 12 |
| `RNMSpreadedListDeductXOL` | 12 |
| `RNMSpreadedListGrossRIXOL` | 12 |
| `RNMSpreadedListGrossXOL` | 12 |
| `RNMSpreadedListNetRIXOL` | 12 |
| `RNMSpreadedListNetXOL` | 12 |
| `RNMSpreadedListRIXOL` | 12 |
| `RNMSpreadedListXOL` | 12 |
| `RnmGrossPremiDisplay` | 12 |
| `RnmLimitList` | 12 |
| `RnmLimitListDisplay` | 12 |
| `SpreadingListXOL` | 12 |
| `SpreadingTotalPctXOL` | 12 |
| `SpreadingTypeIDXOL` | 12 |
| `SpreadingTypeXOL` | 12 |
| `TreatyGroupList` | 12 |
| `Limit` | 4 |
| `Limit2` | 4 |
| `ShareFacultativeReinsurers` | 4 |
| `SpreadingTypeXOLRetro` | 2 |
| `SpreadingTypeXOLRetroID` | 2 |

## `T_TREATY_FAC_SHARE_AMOUNT`

Jalur `TreatyIn.FacultativeShareList.GrossPremiumList` · **12 elemen** · 2 kolom calon

| kolom | ada pada |
| --- | ---: |
| `Currency` | 12 |
| `Value` | 12 |

## `T_TREATY_FAC_SHARE_DEDUCTION`

Jalur `TreatyIn.FacultativeShareList.DeductionList` · **12 elemen** · 5 kolom calon

| kolom | ada pada |
| --- | ---: |
| `Comment` | 12 |
| `Currency` | 12 |
| `Deduction` | 12 |
| `DeductionPct` | 12 |
| `DeductionPctCalculate` | 12 |

## `T_TREATY_HAZARD_LIMIT`

Jalur `TreatyIn` · **1854 elemen** · 134 kolom calon

| kolom | ada pada |
| --- | ---: |
| `AccountingMode` | 1851 |
| `AccountingModeNonProp` | 1851 |
| `AccumulationList` | 1851 |
| `AccumulationPeriod` | 1851 |
| `Bordeaux` | 1851 |
| `Ceding` | 1851 |
| `CedingID` | 1851 |
| `CommentList` | 1851 |
| `CurrencyList` | 1851 |
| `EGNPI` | 1851 |
| `ID` | 1851 |
| `Installment` | 1851 |
| `LeadingReinsSource` | 1851 |
| `LeadingReinsSourceID` | 1851 |
| `LimitShareSummaryList` | 1851 |
| `LimitSummaryList` | 1851 |
| `Limits` | 1851 |
| `Portfolio` | 1851 |
| `ProportionType` | 1851 |
| `ReportingPeriod` | 1851 |
| `ReportingPeriodList` | 1851 |
| `Retention` | 1851 |
| `Share` | 1851 |
| `ShareReins` | 1851 |
| `TotalEgnpiAmountNP` | 1851 |
| `TotalInstallmentNP` | 1851 |
| `TotalLimitDeductblNP` | 1851 |
| `TotalLimitIOONP` | 1851 |
| `TotalLimitMDPNP` | 1851 |
| `TotalLimitPremiEarnNP` | 1851 |
| `TotalRetentionAmountNP` | 1851 |
| `TotalShareDeductionNP` | 1851 |
| `TotalShareGrossNP` | 1851 |
| `TotalShareNetNP` | 1851 |
| `TotalShareRnmNP` | 1851 |
| `TotalShareRnmProp` | 1851 |
| `Commencement` | 1850 |
| `Termination` | 1850 |
| `TreatyContractName` | 1850 |
| `TreatyYear` | 1849 |
| `ChooseStatusAkseptasi` | 1847 |
| `TotalSpreadedNetPremi` | 1842 |
| `TotalSpreadedNetPremiRI` | 1842 |
| `TotalSpreadedRnmProp` | 1842 |
| `TotalSpreadedRnmRIProp` | 1842 |
| `StatusAkseptasi` | 1838 |
| `TeritorialScope` | 1809 |
| `ViewState` | 1677 |
| `ActualValue` | 1614 |
| `ValueDifference` | 1612 |
| `IsMultipleRetro` | 1536 |
| `RNMShareAcrossTheBoard` | 1390 |
| `OLDID` | 1221 |
| `RNMShareP` | 1221 |
| `ReportingConfirmation` | 1219 |
| `ReportingEnd` | 1219 |
| `ReportingSettlement` | 1219 |
| `ReportingStart` | 1219 |
| `ReportingSubmission` | 1219 |
| `ExclusionsP` | 1155 |
| `FacultativeShareList` | 1150 |
| `LimitFacShareSummaryList` | 1149 |
| `TotalFacShareDeductionNP` | 1149 |
| `TotalFacShareGrossNP` | 1149 |
| `TotalFacShareNetNP` | 1149 |
| `TotalFacShareRnmNP` | 1149 |
| `SpecialConditionsP` | 1022 |
| `BordereauxNote` | 1018 |
| `InstallmentNo` | 952 |
| `RNMShare` | 952 |
| `TotalEgnpiAmount` | 952 |
| `TotalEgnpiProportion` | 952 |
| `TotalLimitsROL` | 952 |
| `MDPSummaryList` | 944 |
| `BrokeragePercentP` | 940 |
| `Exclusions` | 896 |
| `OptionLimit` | 828 |
| `BrokeragePercent` | 821 |
| `ContractRefNo` | 742 |
| `SpecialConditions` | 686 |
| `CoInScale` | 659 |
| `TreatyLeader` | 659 |
| `FacultativeShare` | 611 |
| `FacShare` | 367 |
| `Comment` | 314 |
| `ReportingInterval` | 300 |
| `LeadingReinsID` | 297 |
| `MDPSumarry` | 297 |
| `SpecialConditionsp` | 292 |
| `Information` | 196 |
| `MaxCoNonGroup` | 173 |
| `TotalShareGrossMinNP` | 147 |
| `MaxCoGroup` | 87 |
| `TotalLimitMDPMinNP` | 74 |
| `CurrencyEarthquake` | 48 |
| `CurrencyFloodNat` | 48 |
| `Earthquake` | 48 |
| `FloodNation` | 48 |
| `PositionUsername` | 34 |
| `Position` | 33 |
| `RSMDLimit` | 31 |
| `CurrencyRSMD` | 30 |
| `RnmShareDeducted` | 14 |
| `ShareFacultativeReinsurers` | 14 |
| `CurrencyFloodJab` | 13 |
| `FloodJab` | 13 |
| `FacultativeShareBrokerage` | 10 |
| `IsEditData` | 10 |
| `FacShareBrokerage` | 4 |
| `RevisionHistory` | 4 |
| `RevisionState` | 4 |
| `pyActiveWorkGroup` | 3 |
| `pyActiveWorkGroupLabel` | 3 |
| `pyCommonParams` | 3 |
| `pyCurrentPage` | 3 |
| `pyCustomPortalParams` | 3 |
| `pyGadget` | 3 |
| `pyLabel` | 3 |
| `pyOverridePreferences` | 3 |
| `pyPortalPages` | 3 |
| `pyPortalSummary` | 3 |
| `pyRefreshInterval` | 3 |
| `pyReportingDateTimeInterval` | 3 |
| `pyReportingTimeInterval` | 3 |
| `pyRuleHarness` | 3 |
| `pyStartPage` | 3 |
| `pyThemeName` | 3 |
| `pyTopPerformers` | 3 |
| `pyUrgentWorkHeader` | 3 |
| `pyWelcomeHTML` | 3 |
| `EDMState` | 2 |
| `RetroList` | 2 |
| `RevisionDate` | 2 |
| `ValueBeforeProrate` | 2 |

## `T_TREATY_LIMITS`

Jalur `TreatyIn.Limits` · **4210 elemen** · 38 kolom calon

| kolom | ada pada |
| --- | ---: |
| `ID` | 4100 |
| `Currency` | 2850 |
| `EgnpiTotalList` | 2850 |
| `LayerPartType` | 2850 |
| `LayerType` | 2850 |
| `MDPList` | 2850 |
| `PremiumEarnedList` | 2850 |
| `ReinstatementPct` | 2850 |
| `TreatyGroupList` | 2850 |
| `LayerList` | 2849 |
| `Limit` | 2849 |
| `AdjRate` | 2848 |
| `Cover` | 2848 |
| `Currency2` | 2848 |
| `Deductible` | 2847 |
| `ROLPct` | 2847 |
| `MDPPct` | 2846 |
| `ReinstatementValue` | 2845 |
| `Layer` | 2840 |
| `LayerPart` | 2840 |
| `Reinstatement_List` | 2838 |
| `Limit2` | 2441 |
| `Deductible2` | 2436 |
| `IsCombineMDP` | 2155 |
| `CurrencyRelation` | 2116 |
| `NoRIPCalculation` | 2096 |
| `Detail` | 1365 |
| `TreatyType` | 1362 |
| `ReinstatementNote` | 798 |
| `AgregateLimit` | 664 |
| `AgregateLimit2` | 597 |
| `MDPMinList` | 310 |
| `MDPMinPct` | 310 |
| `TreatyTypeID` | 19 |
| `CurrencyID` | 10 |
| `SumLossRatio` | 4 |
| `SumTotalAchievIncured` | 4 |
| `SumTotalAchievNetPremium` | 4 |

## `T_TREATY_LIMIT_ACHIEVEMENT`

Jalur `TreatyIn.Limits.Detail.AchievementLists` · **2031 elemen** · 44 kolom calon

| kolom | ada pada |
| --- | ---: |
| `Currency` | 98 |
| `IncuredClaim` | 98 |
| `Period` | 88 |
| `Value` | 88 |
| `CASHCALL` | 10 |
| `EstimationCASHCALL` | 10 |
| `NETPREMIUM` | 10 |
| `OutstandingCASHCALL` | 10 |
| `OutstandingClaim` | 10 |
| `PREMIUM` | 10 |
| `PaidClaim` | 10 |
| `Total` | 10 |
| `BROKERAGE` | 4 |
| `CASHCALLtoIDR` | 4 |
| `Conversion` | 4 |
| `CurrencyID` | 4 |
| `EstimationCASHCALLtoIDR` | 4 |
| `IDPEGA` | 4 |
| `IncuredClaimtoIDR` | 4 |
| `LossRatio` | 4 |
| `NETPREMIUMtoIDR` | 4 |
| `NOOFFER` | 4 |
| `NOPOLIS` | 4 |
| `OutstandingCASHCALLtoIDR` | 4 |
| `OutstandingClaimtoIDR` | 4 |
| `PREMIUMtoIDR` | 4 |
| `PaidClaimtoIDR` | 4 |
| `QUARTERYEAR` | 4 |
| `Quarter` | 4 |
| `RICOMM` | 4 |
| `SOBNAME` | 4 |
| `TREATYGROUPNAME` | 4 |
| `TREATYTYPE` | 4 |
| `TotalAchievCASHCALL` | 4 |
| `TotalAchievEstCASHCALL` | 4 |
| `TotalAchievIncured` | 4 |
| `TotalAchievNetPremium` | 4 |
| `TotalAchievOsCASHCALL` | 4 |
| `TotalAchievOsClaim` | 4 |
| `TotalAchievOuts` | 4 |
| `TotalAchievPaid` | 4 |
| `TotalAchievPremium` | 4 |
| `TotalAfterClaim` | 4 |
| `TotaltoIDR` | 4 |

## `T_TREATY_LIMIT_AMOUNT`

Jalur `TreatyIn.Limits.Detail.EPIList` · **2867 elemen** · 3 kolom calon

| kolom | ada pada |
| --- | ---: |
| `Value` | 2866 |
| `Currency` | 2812 |
| `CurrencyID` | 17 |

## `T_TREATY_LIMIT_COB`

Jalur `TreatyIn.Limits.Detail.COBList` · **6547 elemen** · 4 kolom calon

| kolom | ada pada |
| --- | ---: |
| `ClassOfBusiness` | 6544 |
| `ClassOfBusinessID` | 6484 |
| `TreatyGroup` | 5624 |
| `TreatyGroupID` | 5624 |

## `T_TREATY_LIMIT_DETAIL`

Jalur `TreatyIn.Limits.Detail` · **2868 elemen** · 79 kolom calon

| kolom | ada pada |
| --- | ---: |
| `COBList` | 2868 |
| `CashLossList` | 2868 |
| `CessionList` | 2868 |
| `EPIList` | 2868 |
| `IOOLimitList` | 2868 |
| `PLAList` | 2868 |
| `RetentionList` | 2868 |
| `TreatyGroup` | 2868 |
| `TreatyType` | 2868 |
| `TreatyGroupID` | 2867 |
| `RNMShareList` | 2865 |
| `RNMSpreadedList` | 2865 |
| `RNMSpreadedListRI` | 2865 |
| `SpreadingList` | 2865 |
| `SpreadingTotalPct` | 2865 |
| `RIOGR` | 2861 |
| `CessionPct` | 2860 |
| `SpreadingTypeID` | 2856 |
| `ClaimCoopList` | 2851 |
| `AchievementLists` | 2820 |
| `DeductionList` | 2820 |
| `DeductionTotalList` | 2820 |
| `SpreadingType` | 2641 |
| `ParentID` | 2629 |
| `ReserveList` | 2305 |
| `IOOPct` | 2236 |
| `ProfitCommision` | 2143 |
| `ProfitME` | 2125 |
| `ProfitYDCF` | 2103 |
| `RNMShare` | 1924 |
| `QSPct` | 1469 |
| `RetentionPct` | 1469 |
| `Brokerage` | 1463 |
| `CurrencyList` | 1450 |
| `Surplus` | 1396 |
| `ID` | 1328 |
| `ShareNote` | 1120 |
| `CurrencyEarthquake` | 572 |
| `Earthquake` | 569 |
| `CurrencyFloodNat` | 550 |
| `FloodNation` | 550 |
| `RSMDLimit` | 493 |
| `CurrencyRSMD` | 492 |
| `RIONR` | 356 |
| `CurrencyFloodJab` | 224 |
| `FloodJab` | 218 |
| `LowerBand` | 179 |
| `ReisuredParticipant` | 179 |
| `Periode` | 161 |
| `CurrencyID` | 151 |
| `UpperBand` | 149 |
| `IOOLimitListSurplus` | 66 |
| `CessionListSurplus` | 62 |
| `RetentionListSurplus` | 62 |
| `PremiumReservePct` | 16 |
| `AchievementPct` | 6 |
| `LossRatio` | 6 |
| `SumLossRatio` | 6 |
| `SumTotalAchievIncured` | 6 |
| `SumTotalAchievNetPremium` | 6 |
| `TotalAchAfterClaim` | 6 |
| `TotalAchCASHCALL` | 6 |
| `TotalAchEstCASHCALL` | 6 |
| `TotalAchIncured` | 6 |
| `TotalAchNetPremium` | 6 |
| `TotalAchOsCASHCALL` | 6 |
| `TotalAchOsClaim` | 6 |
| `TotalAchOuts` | 6 |
| `TotalAchPaid` | 6 |
| `TotalAchPremium` | 6 |
| `CashLoss` | 2 |
| `ClaimCooperation` | 2 |
| `ClassOfBusiness` | 2 |
| `CurrencyCashLoss` | 2 |
| `CurrencyClaimCooperation` | 2 |
| `CurrencyEPI` | 2 |
| `CurrencyPLA` | 2 |
| `EPI` | 2 |
| `PLA` | 2 |

## `T_TREATY_LIMIT_GROUP`

Jalur `TreatyIn.Limits.TreatyGroupList` · **8433 elemen** · 6 kolom calon

| kolom | ada pada |
| --- | ---: |
| `TreatyGroup` | 8433 |
| `ClassOfBusinessList` | 8432 |
| `TreatyGroupID` | 8413 |
| `IsROLProfile` | 7535 |
| `ClassOfBusiness` | 9 |
| `ClassOfBusinessID` | 5 |

## `T_TREATY_LIMIT_GROUP_COB`

Jalur `TreatyIn.Limits.TreatyGroupList.ClassOfBusinessList` · **15746 elemen** · 4 kolom calon

| kolom | ada pada |
| --- | ---: |
| `ClassOfBusiness` | 15742 |
| `ClassOfBusinessID` | 15708 |
| `TreatyGroup` | 15488 |
| `TreatyGroupID` | 15488 |

## `T_TREATY_LIMIT_MEASURE`

Jalur `TreatyIn.Limits.MDPList` · **2953 elemen** · 3 kolom calon

| kolom | ada pada |
| --- | ---: |
| `Currency` | 2953 |
| `Value` | 2951 |
| `CurrencyID` | 16 |

## `T_TREATY_LIMIT_SUMMARY`

Jalur `TreatyIn.LimitSummaryList` · **2855 elemen** · 15 kolom calon

| kolom | ada pada |
| --- | ---: |
| `AggregateLimit` | 2855 |
| `AggregateLimit2` | 2855 |
| `Currency` | 2855 |
| `Deductible` | 2855 |
| `Deductible2` | 2855 |
| `LayerPartType` | 2855 |
| `LayerType` | 2855 |
| `Limit` | 2855 |
| `Limit2` | 2855 |
| `Note` | 2855 |
| `MDP` | 2845 |
| `MDP2` | 2845 |
| `Layer` | 2844 |
| `LayerPart` | 2844 |
| `Currency2` | 2779 |

## `T_TREATY_REINSTATEMENT`

Jalur `TreatyIn.Limits.Reinstatement_List` · **6348 elemen** · 9 kolom calon

| kolom | ada pada |
| --- | ---: |
| `ReinstatementNote` | 6348 |
| `ReinstatementPct` | 6348 |
| `ReinstatementValue` | 6348 |
| `ReinstatementAmount1` | 4788 |
| `ReinstatementAmount2` | 4788 |
| `ID` | 4500 |
| `AdditionalAmount1` | 3358 |
| `AdditionalAmount2` | 3358 |
| `AdditionalPct` | 3358 |

## `T_TREATY_RETRO_SHARE`

Jalur `TreatyIn.ShareReins` · **9 elemen** · 4 kolom calon

| kolom | ada pada |
| --- | ---: |
| `ReinsName` | 4 |
| `ReinsID` | 3 |
| `SharePct` | 2 |
| `Layer` | 1 |

## `T_TREATY_REVISION`

Jalur `TreatyIn` · **1854 elemen** · 134 kolom calon

| kolom | ada pada |
| --- | ---: |
| `AccountingMode` | 1851 |
| `AccountingModeNonProp` | 1851 |
| `AccumulationList` | 1851 |
| `AccumulationPeriod` | 1851 |
| `Bordeaux` | 1851 |
| `Ceding` | 1851 |
| `CedingID` | 1851 |
| `CommentList` | 1851 |
| `CurrencyList` | 1851 |
| `EGNPI` | 1851 |
| `ID` | 1851 |
| `Installment` | 1851 |
| `LeadingReinsSource` | 1851 |
| `LeadingReinsSourceID` | 1851 |
| `LimitShareSummaryList` | 1851 |
| `LimitSummaryList` | 1851 |
| `Limits` | 1851 |
| `Portfolio` | 1851 |
| `ProportionType` | 1851 |
| `ReportingPeriod` | 1851 |
| `ReportingPeriodList` | 1851 |
| `Retention` | 1851 |
| `Share` | 1851 |
| `ShareReins` | 1851 |
| `TotalEgnpiAmountNP` | 1851 |
| `TotalInstallmentNP` | 1851 |
| `TotalLimitDeductblNP` | 1851 |
| `TotalLimitIOONP` | 1851 |
| `TotalLimitMDPNP` | 1851 |
| `TotalLimitPremiEarnNP` | 1851 |
| `TotalRetentionAmountNP` | 1851 |
| `TotalShareDeductionNP` | 1851 |
| `TotalShareGrossNP` | 1851 |
| `TotalShareNetNP` | 1851 |
| `TotalShareRnmNP` | 1851 |
| `TotalShareRnmProp` | 1851 |
| `Commencement` | 1850 |
| `Termination` | 1850 |
| `TreatyContractName` | 1850 |
| `TreatyYear` | 1849 |
| `ChooseStatusAkseptasi` | 1847 |
| `TotalSpreadedNetPremi` | 1842 |
| `TotalSpreadedNetPremiRI` | 1842 |
| `TotalSpreadedRnmProp` | 1842 |
| `TotalSpreadedRnmRIProp` | 1842 |
| `StatusAkseptasi` | 1838 |
| `TeritorialScope` | 1809 |
| `ViewState` | 1677 |
| `ActualValue` | 1614 |
| `ValueDifference` | 1612 |
| `IsMultipleRetro` | 1536 |
| `RNMShareAcrossTheBoard` | 1390 |
| `OLDID` | 1221 |
| `RNMShareP` | 1221 |
| `ReportingConfirmation` | 1219 |
| `ReportingEnd` | 1219 |
| `ReportingSettlement` | 1219 |
| `ReportingStart` | 1219 |
| `ReportingSubmission` | 1219 |
| `ExclusionsP` | 1155 |
| `FacultativeShareList` | 1150 |
| `LimitFacShareSummaryList` | 1149 |
| `TotalFacShareDeductionNP` | 1149 |
| `TotalFacShareGrossNP` | 1149 |
| `TotalFacShareNetNP` | 1149 |
| `TotalFacShareRnmNP` | 1149 |
| `SpecialConditionsP` | 1022 |
| `BordereauxNote` | 1018 |
| `InstallmentNo` | 952 |
| `RNMShare` | 952 |
| `TotalEgnpiAmount` | 952 |
| `TotalEgnpiProportion` | 952 |
| `TotalLimitsROL` | 952 |
| `MDPSummaryList` | 944 |
| `BrokeragePercentP` | 940 |
| `Exclusions` | 896 |
| `OptionLimit` | 828 |
| `BrokeragePercent` | 821 |
| `ContractRefNo` | 742 |
| `SpecialConditions` | 686 |
| `CoInScale` | 659 |
| `TreatyLeader` | 659 |
| `FacultativeShare` | 611 |
| `FacShare` | 367 |
| `Comment` | 314 |
| `ReportingInterval` | 300 |
| `LeadingReinsID` | 297 |
| `MDPSumarry` | 297 |
| `SpecialConditionsp` | 292 |
| `Information` | 196 |
| `MaxCoNonGroup` | 173 |
| `TotalShareGrossMinNP` | 147 |
| `MaxCoGroup` | 87 |
| `TotalLimitMDPMinNP` | 74 |
| `CurrencyEarthquake` | 48 |
| `CurrencyFloodNat` | 48 |
| `Earthquake` | 48 |
| `FloodNation` | 48 |
| `PositionUsername` | 34 |
| `Position` | 33 |
| `RSMDLimit` | 31 |
| `CurrencyRSMD` | 30 |
| `RnmShareDeducted` | 14 |
| `ShareFacultativeReinsurers` | 14 |
| `CurrencyFloodJab` | 13 |
| `FloodJab` | 13 |
| `FacultativeShareBrokerage` | 10 |
| `IsEditData` | 10 |
| `FacShareBrokerage` | 4 |
| `RevisionHistory` | 4 |
| `RevisionState` | 4 |
| `pyActiveWorkGroup` | 3 |
| `pyActiveWorkGroupLabel` | 3 |
| `pyCommonParams` | 3 |
| `pyCurrentPage` | 3 |
| `pyCustomPortalParams` | 3 |
| `pyGadget` | 3 |
| `pyLabel` | 3 |
| `pyOverridePreferences` | 3 |
| `pyPortalPages` | 3 |
| `pyPortalSummary` | 3 |
| `pyRefreshInterval` | 3 |
| `pyReportingDateTimeInterval` | 3 |
| `pyReportingTimeInterval` | 3 |
| `pyRuleHarness` | 3 |
| `pyStartPage` | 3 |
| `pyThemeName` | 3 |
| `pyTopPerformers` | 3 |
| `pyUrgentWorkHeader` | 3 |
| `pyWelcomeHTML` | 3 |
| `EDMState` | 2 |
| `RetroList` | 2 |
| `RevisionDate` | 2 |
| `ValueBeforeProrate` | 2 |

## `T_TREATY_SHARE`

Jalur `TreatyIn.Share` · **2923 elemen** · 30 kolom calon

| kolom | ada pada |
| --- | ---: |
| `ClassofBusinessList` | 2923 |
| `GrossPremiumList` | 2923 |
| `LayerPartType` | 2923 |
| `LayerType` | 2923 |
| `NetPremiumList` | 2923 |
| `RnmLimitList` | 2923 |
| `TreatyGroupList` | 2922 |
| `Layer` | 2913 |
| `LayerPart` | 2913 |
| `RnmGrossPremiDisplay` | 2908 |
| `RnmLimitListDisplay` | 2908 |
| `Cover` | 2847 |
| `RNMSpreadedListGrossRIXOL` | 2823 |
| `RNMSpreadedListGrossXOL` | 2823 |
| `RNMSpreadedListNetRIXOL` | 2823 |
| `RNMSpreadedListNetXOL` | 2823 |
| `RNMSpreadedListRIXOL` | 2823 |
| `RNMSpreadedListXOL` | 2823 |
| `SpreadingListXOL` | 2823 |
| `SpreadingTotalPctXOL` | 2823 |
| `SpreadingTypeXOL` | 2822 |
| `DeductionList` | 2813 |
| `DeductionTotalList` | 2813 |
| `RNMSpreadedListDeductRIXOL` | 2805 |
| `RNMSpreadedListDeductXOL` | 2805 |
| `SpreadingTypeIDXOL` | 2586 |
| `RNMShare` | 1003 |
| `GrossPremiumMinList` | 305 |
| `RNMSpreadedListGrossMinXOL` | 305 |
| `RNMSpreadedListGrossRIMinXOL` | 305 |

## `T_TREATY_SHARE_AMOUNT`

Jalur `TreatyIn.Share.GrossPremiumList` · **3049 elemen** · 2 kolom calon

| kolom | ada pada |
| --- | ---: |
| `Currency` | 3025 |
| `Value` | 3025 |

## `T_TREATY_SHARE_DEDUCTION`

Jalur `TreatyIn.Share.DeductionList` · **2419 elemen** · 5 kolom calon

| kolom | ada pada |
| --- | ---: |
| `Currency` | 2419 |
| `Deduction` | 2419 |
| `Comment` | 2414 |
| `DeductionPct` | 2414 |
| `DeductionPctCalculate` | 2352 |

## `T_TREATY_SHARE_SPREADING`

Jalur `TreatyIn.Share.SpreadingListXOL` · **5550 elemen** · 6 kolom calon

| kolom | ada pada |
| --- | ---: |
| `ReinsTypeName` | 5547 |
| `Pct` | 5529 |
| `Usd` | 5479 |
| `ParentReinsTypeID` | 5078 |
| `ReinsTypeID` | 5014 |
| `Rp` | 5011 |

## `T_TREATY_SHARE_SPREAD_AMOUNT`

Jalur `TreatyIn.Share.SpreadingListXOL.RNMSpreadedListXOL` · **0 elemen** · 0 kolom calon

⛔ **NOL elemen di 1.854 dokumen.** Jalur ini tidak pernah berisi.

## `T_TREATY_TOTAL`

Jalur `TreatyIn.TotalShareNetNP` · **872 elemen** · 2 kolom calon

| kolom | ada pada |
| --- | ---: |
| `Currency` | 872 |
| `Value` | 872 |

## `T_TREATY_VALUE_BEFORE_PRORATE`

Jalur `TreatyIn.ValueBeforeProrate` · **2 elemen** · 12 kolom calon

| kolom | ada pada |
| --- | ---: |
| `Installment` | 2 |
| `LimitShareSummaryList` | 2 |
| `Share` | 2 |
| `TotalInstallmentNP` | 2 |
| `TotalShareDeductionNP` | 2 |
| `TotalShareGrossNP` | 2 |
| `TotalShareNetNP` | 2 |
| `TotalShareRnmNP` | 2 |
| `TotalSpreadedNetPremi` | 2 |
| `TotalSpreadedNetPremiRI` | 2 |
| `TotalSpreadedRnmProp` | 2 |
| `TotalSpreadedRnmRIProp` | 2 |

## `T_TREATY_VALUE_DIFFERENCE`

Jalur `TreatyIn.ValueDifference` · **1612 elemen** · 23 kolom calon

| kolom | ada pada |
| --- | ---: |
| `EGNPI` | 1612 |
| `LimitShareSummaryList` | 1612 |
| `LimitSummaryList` | 1612 |
| `Limits` | 1612 |
| `Share` | 1612 |
| `TotalEgnpiAmountNP` | 1612 |
| `TotalLimitMDPNP` | 1612 |
| `TotalLimitPremiEarnNP` | 1612 |
| `TotalShareDeductionNP` | 1612 |
| `TotalShareGrossNP` | 1612 |
| `TotalShareNetNP` | 1612 |
| `TotalSpreadedNetPremi` | 1612 |
| `TotalSpreadedNetPremiRI` | 1612 |
| `FacultativeShareList` | 500 |
| `LimitFacShareSummaryList` | 500 |
| `TotalFacShareDeductionNP` | 500 |
| `TotalFacShareGrossNP` | 500 |
| `TotalFacShareNetNP` | 500 |
| `Installment` | 2 |
| `TotalInstallmentNP` | 2 |
| `TotalShareRnmNP` | 2 |
| `TotalSpreadedRnmProp` | 2 |
| `TotalSpreadedRnmRIProp` | 2 |

