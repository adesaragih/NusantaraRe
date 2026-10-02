# Inventaris Rule — Treaty In

STEP D1, batch 1 (domain treaty inward). Dibangkitkan 2026-09-12 dari korpus READ-ONLY
`D:\XML\RNM_BRD\Treaty In\`.

**Batas step ini:** hanya struktur + metadata. Tidak ada kesimpulan logika bisnis di sini.
Penelusuran perilaku adalah STEP D2 (`../flows/`).

**Identitas rule** diambil dari `<pxInsName>` (lihat `../README.md` §6.4).  
Kolom **Asal salinan** berisi `<pzOriginalInstanceKey>` dan hanya diisi bila **berbeda**
dari identitas — artinya rule ini hasil *Save As* dari rule lain. Tanda `—` berarti sama.
Class `—` berarti `<pxInsName>` tidak berisi bagian class (berlaku pada `DataPage`, yang
beridentitas global, bukan per class).

Perintah audit (dijalankan dari `D:\XML\RNM_BRD\`):

```
find "Treaty In" -type f -name "*.xml" | wc -l
find "Treaty In" -type f -name "*.xml" | awk -F/ '{print $2}' | sort | uniq -c
grep -o "<pxInsName>[^<]*" "Treaty In/<tipe>/<file>.xml" | head -1
grep -o "<pzOriginalInstanceKey>[^<]*" "Treaty In/<tipe>/<file>.xml" | head -1
```

## 1. Ringkasan terukur

- Total file XML: **329** — cocok dengan `../README.md` §5.1
- **Identitas rule unik: 329** dari 329 file → tidak ada identitas ganda di dalam modul ini
- Class Pega unik: **38**
- RDBList menurut jenis SQL: PLSQL=9, QUERY=31 (`PLSQL` = blok anonim `BEGIN…END;`, umumnya memanggil stored procedure)
- Prefix pada kunci RDBList: `ASM`=24, `RNM`=16 (arti prefix **belum terverifikasi** — OQ-008)

| Tipe rule | Jumlah |
| --- | ---: |
| Activity | 149 |
| ConnectREST | 1 |
| DataTransform | 35 |
| DecisionTable | 1 |
| FlowAction | 32 |
| Harness | 3 |
| RDBList | 40 |
| ReportDefinition | 16 |
| Section | 50 |
| SystemSettings | 1 |
| When | 1 |
| **TOTAL** | **329** |

## 2. Sebaran class Pega

Class yang **tidak** berawalan `ASM-` adalah class bawaan Pega atau class portal — ruang
lingkup migrasinya belum ditentukan (**OQ-009**).

| Class | Jumlah rule | Aplikasi? |
| --- | ---: | --- |
| `DATA-PORTAL` | 146 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | 31 | aplikasi |
| `ASM-FW-GISFW-INT-TREATY_IN` | 26 | aplikasi |
| `ASM-FW-GISFW-DATA-TREATYINLIMITS` | 23 | aplikasi |
| `ASM-FW-GISFW-DATA-TREATYINSHARE` | 18 | aplikasi |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | 12 | aplikasi |
| `ASM-FW-GISFW-DATA-TREATYINSHAREREINS` | 7 | aplikasi |
| `ASM-FW-GISFW-DATA-TREATYINRETROSHARE` | 6 | aplikasi |
| `@BASECLASS` | 5 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-DATA-TREATYINLIMITLAYER` | 5 | aplikasi |
| `ASM-FW-GISFW-DATA-TREATYINLIMITSSPREADING` | 4 | aplikasi |
| `ASM-FW-GISFW-INT-PROPORTIONALARRG` | 4 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYINDETAIL` | 4 | aplikasi |
| `ASM-FW-GISFW-INT-TREATY_IN_EDM` | 4 | aplikasi |
| `ASM-FW-GISFW-DATA-TREATYINEGNPI` | 3 | aplikasi |
| `ASM-FW-GISFW-DATA-TREATYININSTALLMENT` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-CURRENCY` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-OFFERJSON` | 3 | aplikasi |
| `ASM-FW-GISFW-DATA-TREATYINRETENTION` | 2 | aplikasi |
| `ASM-FW-GISFW-WORK` | 2 | aplikasi |
| `ASM-FW-GISFW-DATA-OFFERFACIN` | 1 | aplikasi |
| `ASM-FW-GISFW-DATA-TREATYINCURRENCYLIST` | 1 | aplikasi |
| `ASM-FW-GISFW-DATA-TREATYINDEDUCTION` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-AGENT` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-BUSINESS` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-BUSINESSGROUP` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-CLIENT` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_LINK_SERVICE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-REINSURANCETYPE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYBUSINESS` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYEXCHANGEYEARLY` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYGROUP` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYINOFFER` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYREINSURER` | 1 | aplikasi |
| `ASM-FW-GISFW-WORK-LIFE` | 1 | aplikasi |
| `ASSIGN-WORKLIST` | 1 | **bukan aplikasi** → OQ-009 |
| `LINKSERVICE` | 1 | **bukan aplikasi** → OQ-009 |
| `WORK-` | 1 | **bukan aplikasi** → OQ-009 |

Rule di class bukan-aplikasi: **154** dari 329.

## 3. Daftar rule per tipe

### Activity — 149 rule (`RULE-OBJ-ACTIVITY`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AddClassofBusiness.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `ADDCLASSOFBUSINESS` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `AddCommentList_Act.xml` | `DATA-PORTAL` | `ADDCOMMENTLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `AddDeduction.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `ADDDEDUCTION` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GISFW-DATA-TREATYINSHARE / ADDDEDUCTION` |
| `AddDelSpreadingTreatyin.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `ADDDELSPREADINGTREATYIN` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `AddFacRetroProp.xml` | `DATA-PORTAL` | `ADDFACRETROPROP` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `AddSpreadingXOL.xml` | `ASM-FW-GISFW-DATA-TREATYINSHARE` | `ADDSPREADINGXOL` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `AddValue.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `ADDVALUE` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `Akseptasi_Act.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `AKSEPTASI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CalculateDeduction.xml` | `ASM-FW-GISFW-DATA-TREATYINSHARE` | `CALCULATEDEDUCTION` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | `ASM-FW-GISFW-DATA-TREATYINSHARE / DEDUCTIONTOTAL` |
| `CalculateORRItoRetro.xml` | `ASM-FW-GISFW-DATA-TREATYINSHARE` | `CALCULATEORRITORETRO` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `CalculateRetroSumary.xml` | `ASM-FW-GISFW-DATA-TREATYINRETROSHARE` | `CALCULATERETROSUMARY` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | `DATA-PORTAL / TREATYINSUMMARYLIMITSHARE` |
| `CalculateRetroTotal.xml` | `ASM-FW-GISFW-DATA-TREATYINRETROSHARE` | `CALCULATERETROTOTAL` | [terverifikasi] pyActivityType=ACTIVITY; steps=25 | — |
| `CalculateShareList.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `CALCULATESHARELIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `CalculateSharePctToRetro.xml` | `ASM-FW-GISFW-DATA-TREATYINRETROSHARE` | `CALCULATESHAREPCTTORETRO` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `ChangeDokument_Act.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `CHANGEDOKUMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CheckDuplicateOffer.xml` | `DATA-PORTAL` | `CHECKDUPLICATEOFFER` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `ConvertHistoryDate.xml` | `DATA-PORTAL` | `CONVERTHISTORYDATE` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CopyLastLimitNP.xml` | `DATA-PORTAL` | `COPYLASTLIMITNP` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CountRetroShare_Act.xml` | `ASM-FW-GISFW-DATA-TREATYINSHAREREINS` | `COUNTRETROSHARE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `DelSpreadingRetro.xml` | `DATA-PORTAL` | `DELSPREADINGRETRO` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `DeleteGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `DELETEGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / GETURLGOOGLESTORAGE_ACT` |
| `Delete_act.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `DELETE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `DetailCalculation.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `DETAILCALCULATION` | [terverifikasi] pyActivityType=ACTIVITY; steps=24 | — |
| `DetailCalculationROL.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `DETAILCALCULATIONROL` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | — |
| `DownloadAll_Act.xml` | `DATA-PORTAL` | `DOWNLOADALL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | — |
| `DownloadAttachmentTreaty.xml` | `@BASECLASS` | `DOWNLOADATTACHMENTTREATY` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | `@BASECLASS / DOWNLOADATTACHMENT2` |
| `FetchQSfromMaster.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `FETCHQSFROMMASTER` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `FetchQSfromMasterXOL.xml` | `ASM-FW-GISFW-DATA-TREATYINSHARE` | `FETCHQSFROMMASTERXOL` | [terverifikasi] pyActivityType=ACTIVITY; steps=49 | — |
| `FetchTreatyContractDetail.xml` | `ASM-FW-GISFW-DATA-TREATYINRETROSHARE` | `FETCHTREATYCONTRACTDETAIL` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `FetchTreatyExistingProduction.xml` | `DATA-PORTAL` | `FETCHTREATYEXISTINGPRODUCTION` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GenerateCSVTreaty.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `GENERATECSVTREATY` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `GetAchievement.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `GETACHIEVEMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=38 | — |
| `GetEstimasiClaim_act.xml` | `ASM-FW-GISFW-WORK` | `GETESTIMASICLAIM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=16 | — |
| `GetHistoryClaim_act.xml` | `ASM-FW-GISFW-WORK` | `GETHISTORYCLAIM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=16 | — |
| `GetLinkService.xml` | `ASM-FW-GISFW-INT-M_LINK_SERVICE` | `GETLINKSERVICE` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetMasterTreatyCategory_Act.xml` | `DATA-PORTAL` | `GETMASTERTREATYCATEGORY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `GetNilaiTotal.xml` | `ASM-FW-GISFW-DATA-TREATYINSHARE` | `GETNILAITOTAL` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | `ASM-FW-GISFW-DATA-TREATYINSHARE / FETCHQSFROMMASTERXOL` |
| `GetNilaiTotalAdjustment.xml` | `ASM-FW-GISFW-DATA-TREATYINSHARE` | `GETNILAITOTALADJUSTMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=36 | `ASM-FW-GISFW-DATA-TREATYINSHARE / GETNILAITOTAL` |
| `GetSpreadingRetro.xml` | `ASM-FW-GISFW-DATA-TREATYINSHARE` | `GETSPREADINGRETRO` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `GetUrlGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETURLGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / INSERTGOOGLESTORAGE_ACT` |
| `InputParamUploadReas_act.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `INPUTPARAMUPLOADREAS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `InsertGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERTGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `InsertToLogAchievement.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `INSERTTOLOGACHIEVEMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `LimitCalculation.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `LIMITCALCULATION` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | — |
| `LoadAttachment.xml` | `DATA-PORTAL` | `LOADATTACHMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `WORK- / LOADATTACHMENT` |
| `LoadAttachmentData.xml` | `WORK-` | `LOADATTACHMENTDATA` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `PremiumReserveCalculate.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `PREMIUMRESERVECALCULATE` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `SaveTreatyInDetailEdm_Act.xml` | `DATA-PORTAL` | `SAVETREATYINDETAILEDM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=71 | `DATA-PORTAL / SAVETREATYINDETAIL_ACT` |
| `SaveTreatyInDetail_Act.xml` | `DATA-PORTAL` | `SAVETREATYINDETAIL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=57 | — |
| `SaveTreatyInOfferNonProp1_Act.xml` | `DATA-PORTAL` | `SAVETREATYINOFFERNONPROP1_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | `DATA-PORTAL / SAVETREATYINOFFER_ACT` |
| `SaveTreatyInOfferProportional_Act.xml` | `DATA-PORTAL` | `SAVETREATYINOFFERPROPORTIONAL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | `DATA-PORTAL / SAVETREATYINOFFERNONPROP2_ACT` |
| `SaveTreatyInOffer_Act.xml` | `DATA-PORTAL` | `SAVETREATYINOFFER_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SaveTreatyIn_Act.xml` | `DATA-PORTAL` | `SAVETREATYIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | — |
| `SaveTreatyIn_EDM_Act.xml` | `DATA-PORTAL` | `SAVETREATYIN_EDM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | `DATA-PORTAL / SAVETREATYIN_ACT` |
| `SetAmountConversion.xml` | `ASM-FW-GISFW-DATA-TREATYINEGNPI` | `SETAMOUNTCONVERSION` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SetCategoryAttach.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SETCATEGORYATTACH` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SetCurrNameMasterTreaty_Act.xml` | `ASM-FW-GISFW-DATA-TREATYINCURRENCYLIST` | `SETCURRNAMEMASTERTREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetCurrName_Act.xml` | `ASM-FW-GISFW-DATA-TREATYINDEDUCTION` | `SETCURRNAME_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `ASM-FW-GISFW-DATA-TREATYINCURRENCYLIST / SETCURRNAMEMASTERTREATY_ACT` |
| `SetReinstatementPct.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `SETREINSTATEMENTPCT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SetSpreadName.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `SETSPREADNAME` | [terverifikasi] pyActivityType=ACTIVITY; steps=31 | — |
| `SetSpreadingXOL.xml` | `ASM-FW-GISFW-DATA-TREATYINSHARE` | `SETSPREADINGXOL` | [terverifikasi] pyActivityType=ACTIVITY; steps=47 | — |
| `SetTotalInstallment.xml` | `ASM-FW-GISFW-DATA-TREATYININSTALLMENT` | `SETTOTALINSTALLMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `SetTreatyGroupName_Act.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `SETTREATYGROUPNAME_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetTreatyIn_Act.xml` | `DATA-PORTAL` | `SETTREATYIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | — |
| `SetTreatyTypeName_Act.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `SETTREATYTYPENAME_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GISFW-DATA-TREATYINLIMITS / SETTREATYGROUPNAME_ACT` |
| `SetTreatyinRetro_Act.xml` | `ASM-FW-GISFW-DATA-TREATYINSHAREREINS` | `SETTREATYINRETRO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetViewTreatyIn_Act.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `SETVIEWTREATYIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `DATA-PORTAL / SETVIEWTREATYIN_ACT` |
| `SetkategoriDoc.xml` | `@BASECLASS` | `SETKATEGORIDOC` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `DATA-PORTAL / SETKATEGORIDOC` |
| `TotalEgnpi.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `TOTALEGNPI` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | `ASM-FW-GISFW-DATA-TREATYINLIMITLAYER / TOTALEGNPI` |
| `TreatyEDMCalculateDifference.xml` | `DATA-PORTAL` | `TREATYEDMCALCULATEDIFFERENCE` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | — |
| `TreatyEDMDifferenceDeduction.xml` | `DATA-PORTAL` | `TREATYEDMDIFFERENCEDEDUCTION` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | `DATA-PORTAL / TREATYINDIFFERENCEDEDUCTION` |
| `TreatyEDMDifferenceLimits.xml` | `DATA-PORTAL` | `TREATYEDMDIFFERENCELIMITS` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | `DATA-PORTAL / TREATYINDIFFERENCELIMITS` |
| `TreatyEDMDifferencePremium.xml` | `DATA-PORTAL` | `TREATYEDMDIFFERENCEPREMIUM` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `DATA-PORTAL / TREATYINDIFFERENCEPREMIUM` |
| `TreatyEDMDifferenceShare.xml` | `DATA-PORTAL` | `TREATYEDMDIFFERENCESHARE` | [terverifikasi] pyActivityType=ACTIVITY; steps=18 | `DATA-PORTAL / TREATYINDIFFERENCESHARE` |
| `TreatyEDMProRateCalculation.xml` | `DATA-PORTAL` | `TREATYEDMPRORATECALCULATION` | [terverifikasi] pyActivityType=ACTIVITY; steps=27 | — |
| `TreatyInAccumulationSetSubDue.xml` | `DATA-PORTAL` | `TREATYINACCUMULATIONSETSUBDUE` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `TreatyInActualShare.xml` | `DATA-PORTAL` | `TREATYINACTUALSHARE` | [terverifikasi] pyActivityType=ACTIVITY; steps=27 | — |
| `TreatyInActualUpdateValue.xml` | `DATA-PORTAL` | `TREATYINACTUALUPDATEVALUE` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `TreatyInActualUpdateValueLimits.xml` | `DATA-PORTAL` | `TREATYINACTUALUPDATEVALUELIMITS` | [terverifikasi] pyActivityType=ACTIVITY; steps=31 | `DATA-PORTAL / TREATYINACTUALUPDATEVALUE` |
| `TreatyInActualUpdateValuePremiumIncome.xml` | `DATA-PORTAL` | `TREATYINACTUALUPDATEVALUEPREMIUMINCOME` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `TreatyInActualUpdateValueShare.xml` | `DATA-PORTAL` | `TREATYINACTUALUPDATEVALUESHARE` | [terverifikasi] pyActivityType=ACTIVITY; steps=16 | `DATA-PORTAL / TREATYINACTUALUPDATEVALUE` |
| `TreatyInAddCurrency.xml` | `DATA-PORTAL` | `TREATYINADDCURRENCY` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `TreatyInAkseptasiEDM_Act.xml` | `DATA-PORTAL` | `TREATYINAKSEPTASIEDM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `DATA-PORTAL / TREATYINAKSEPTASI_ACT` |
| `TreatyInAkseptasi_Act.xml` | `DATA-PORTAL` | `TREATYINAKSEPTASI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GISFW-INT-TREATY_IN / AKSEPTASI_ACT` |
| `TreatyInCheckCedingBlacklist.xml` | `DATA-PORTAL` | `TREATYINCHECKCEDINGBLACKLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `TreatyInCheckError.xml` | `DATA-PORTAL` | `TREATYINCHECKERROR` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `TreatyInCheckID.xml` | `DATA-PORTAL` | `TREATYINCHECKID` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `TreatyInConvertCallData_act.xml` | `DATA-PORTAL` | `TREATYINCONVERTCALLDATA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | `DATA-PORTAL / SETTREATYIN_ACT` |
| `TreatyInCopy.xml` | `DATA-PORTAL` | `TREATYINCOPY` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `TreatyInDeclineConfirmation_postact.xml` | `DATA-PORTAL` | `TREATYINDECLINECONFIRMATION_POSTACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `TreatyInDeclineConfirmation_postactEDM.xml` | `DATA-PORTAL` | `TREATYINDECLINECONFIRMATION_POSTACTEDM` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | `DATA-PORTAL / TREATYINDECLINECONFIRMATION_POSTACT` |
| `TreatyInDifferenceDeduction.xml` | `DATA-PORTAL` | `TREATYINDIFFERENCEDEDUCTION` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | `DATA-PORTAL / TREATYINSETBROKERAGE` |
| `TreatyInDifferenceFacShare.xml` | `DATA-PORTAL` | `TREATYINDIFFERENCEFACSHARE` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | `DATA-PORTAL / TREATYINDIFFERENCESHARE` |
| `TreatyInDifferenceFacShareTotal.xml` | `DATA-PORTAL` | `TREATYINDIFFERENCEFACSHARETOTAL` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `DATA-PORTAL / TREATYINDIFFERENCESHARETOTAL` |
| `TreatyInDifferenceLimits.xml` | `DATA-PORTAL` | `TREATYINDIFFERENCELIMITS` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `TreatyInDifferenceLimitsSetTotal.xml` | `DATA-PORTAL` | `TREATYINDIFFERENCELIMITSSETTOTAL` | [terverifikasi] pyActivityType=ACTIVITY; steps=18 | `DATA-PORTAL / TREATYINLIMITSACTUALSETTOTAL` |
| `TreatyInDifferenceLimitsSumary.xml` | `DATA-PORTAL` | `TREATYINDIFFERENCELIMITSSUMARY` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `DATA-PORTAL / TREATYINSUMMARYLIMITACTUAL` |
| `TreatyInDifferencePremium.xml` | `DATA-PORTAL` | `TREATYINDIFFERENCEPREMIUM` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `TreatyInDifferenceShare.xml` | `DATA-PORTAL` | `TREATYINDIFFERENCESHARE` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `TreatyInDifferenceShareSumary.xml` | `DATA-PORTAL` | `TREATYINDIFFERENCESHARESUMARY` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | `DATA-PORTAL / TREATYINSUMMARYLIMITSHARE` |
| `TreatyInDifferenceShareTotal.xml` | `DATA-PORTAL` | `TREATYINDIFFERENCESHARETOTAL` | [terverifikasi] pyActivityType=ACTIVITY; steps=25 | — |
| `TreatyInDifferenceSummaryFacShare.xml` | `DATA-PORTAL` | `TREATYINDIFFERENCESUMMARYFACSHARE` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | `DATA-PORTAL / TREATYINSUMMARYLIMITACTUALFACSHARE` |
| `TreatyInDownloadAll.xml` | `DATA-PORTAL` | `TREATYINDOWNLOADALL` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `TreatyInEGNPIListValue.xml` | `DATA-PORTAL` | `TREATYINEGNPILISTVALUE` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `DATA-PORTAL / TREATYINEGNPICONVERTIONLIST` |
| `TreatyInFixSpl.xml` | `DATA-PORTAL` | `TREATYINFIXSPL` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | — |
| `TreatyInInputVis.xml` | `DATA-PORTAL` | `TREATYININPUTVIS` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `TreatyInLimitsActualSetTotal.xml` | `DATA-PORTAL` | `TREATYINLIMITSACTUALSETTOTAL` | [terverifikasi] pyActivityType=ACTIVITY; steps=26 | `DATA-PORTAL / TREATYINACTUALSETTOTAL` |
| `TreatyInLimitsListValue.xml` | `DATA-PORTAL` | `TREATYINLIMITSLISTVALUE` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `TreatyInLimitsListValueProportional.xml` | `DATA-PORTAL` | `TREATYINLIMITSLISTVALUEPROPORTIONAL` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `TreatyInMappingDataconvert.xml` | `DATA-PORTAL` | `TREATYINMAPPINGDATACONVERT` | [terverifikasi] pyActivityType=ACTIVITY; steps=46 | — |
| `TreatyInMappingDataconvertProp.xml` | `DATA-PORTAL` | `TREATYINMAPPINGDATACONVERTPROP` | [terverifikasi] pyActivityType=ACTIVITY; steps=29 | — |
| `TreatyInNPSetTotal.xml` | `DATA-PORTAL` | `TREATYINNPSETTOTAL` | [terverifikasi] pyActivityType=ACTIVITY; steps=81 | — |
| `TreatyInNPSetTotalActual.xml` | `DATA-PORTAL` | `TREATYINNPSETTOTALACTUAL` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `DATA-PORTAL / TREATYINNPSETTOTAL` |
| `TreatyInNPSetTotalActualLimits.xml` | `DATA-PORTAL` | `TREATYINNPSETTOTALACTUALLIMITS` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | `DATA-PORTAL / TREATYINNPSETTOTALACTUAL` |
| `TreatyInNPSetTotalActualShare.xml` | `DATA-PORTAL` | `TREATYINNPSETTOTALACTUALSHARE` | [terverifikasi] pyActivityType=ACTIVITY; steps=26 | `DATA-PORTAL / TREATYINNPSETTOTALACTUAL` |
| `TreatyInNonAddItem.xml` | `DATA-PORTAL` | `TREATYINNONADDITEM` | [terverifikasi] pyActivityType=ACTIVITY; steps=33 | — |
| `TreatyInNonSetTotal.xml` | `DATA-PORTAL` | `TREATYINNONSETTOTAL` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `TreatyInPopulateDetail.xml` | `DATA-PORTAL` | `TREATYINPOPULATEDETAIL` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `TreatyInPropAdd.xml` | `DATA-PORTAL` | `TREATYINPROPADD` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `TreatyInPropshare.xml` | `DATA-PORTAL` | `TREATYINPROPSHARE` | [terverifikasi] pyActivityType=ACTIVITY; steps=25 | — |
| `TreatyInPropshareDetail.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `TREATYINPROPSHAREDETAIL` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `TreatyInSaveROL.xml` | `DATA-PORTAL` | `TREATYINSAVEROL` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | — |
| `TreatyInSetAccountReport.xml` | `DATA-PORTAL` | `TREATYINSETACCOUNTREPORT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `DATA-PORTAL / TREATYINSETREPORT` |
| `TreatyInSetBrokerage.xml` | `DATA-PORTAL` | `TREATYINSETBROKERAGE` | [terverifikasi] pyActivityType=ACTIVITY; steps=30 | — |
| `TreatyInSetBrokerageActual.xml` | `DATA-PORTAL` | `TREATYINSETBROKERAGEACTUAL` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | `DATA-PORTAL / TREATYINSETBROKERAGE` |
| `TreatyInSetReport.xml` | `DATA-PORTAL` | `TREATYINSETREPORT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | — |
| `TreatyInSetToDirector.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `TREATYINSETTODIRECTOR` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `TreatyInSetValue.xml` | `DATA-PORTAL` | `TREATYINSETVALUE` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | `DATA-PORTAL / TREATYINPOPULATEDETAIL` |
| `TreatyInSetValueDifferenceInstallment.xml` | `DATA-PORTAL` | `TREATYINSETVALUEDIFFERENCEINSTALLMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | `DATA-PORTAL / TREATYINSETVALUEINSTALLMENT` |
| `TreatyInSetValueInstallment.xml` | `DATA-PORTAL` | `TREATYINSETVALUEINSTALLMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=20 | — |
| `TreatyInShareActualSetTotal.xml` | `DATA-PORTAL` | `TREATYINSHAREACTUALSETTOTAL` | [terverifikasi] pyActivityType=ACTIVITY; steps=39 | — |
| `TreatyInShareListValue.xml` | `DATA-PORTAL` | `TREATYINSHARELISTVALUE` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `TreatyInSubmit.xml` | `DATA-PORTAL` | `TREATYINSUBMIT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `TreatyInSubmitEDM.xml` | `DATA-PORTAL` | `TREATYINSUBMITEDM` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | `DATA-PORTAL / TREATYINSUBMIT` |
| `TreatyInSummaryLimit.xml` | `DATA-PORTAL` | `TREATYINSUMMARYLIMIT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `DATA-PORTAL / TREATYINSUMMARYMDP` |
| `TreatyInSummaryLimitActual.xml` | `DATA-PORTAL` | `TREATYINSUMMARYLIMITACTUAL` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `DATA-PORTAL / TREATYINSUMMARYLIMIT` |
| `TreatyInSummaryLimitActualFacShare.xml` | `DATA-PORTAL` | `TREATYINSUMMARYLIMITACTUALFACSHARE` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | `DATA-PORTAL / TREATYINSUMMARYLIMITFACSHARE` |
| `TreatyInSummaryLimitFacShare.xml` | `DATA-PORTAL` | `TREATYINSUMMARYLIMITFACSHARE` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | `DATA-PORTAL / TREATYINSUMMARYLIMITSHARE` |
| `TreatyInSummaryLimitShare.xml` | `DATA-PORTAL` | `TREATYINSUMMARYLIMITSHARE` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | `DATA-PORTAL / TREATYINSUMMARYMDPSHARE` |
| `TreatyInSummaryLimitShareActual.xml` | `DATA-PORTAL` | `TREATYINSUMMARYLIMITSHAREACTUAL` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | `DATA-PORTAL / TREATYINSUMMARYLIMITSHARE` |
| `TreatyInSummaryMDP.xml` | `DATA-PORTAL` | `TREATYINSUMMARYMDP` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | `DATA-PORTAL / TREATYINNPSETTOTAL` |
| `TreatyInTestAgent.xml` | `DATA-PORTAL` | `TREATYINTESTAGENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `TreatyInUploadCSV.xml` | `DATA-PORTAL` | `TREATYINUPLOADCSV` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `TreatyInXOLAddSpreading.xml` | `DATA-PORTAL` | `TREATYINXOLADDSPREADING` | [terverifikasi] pyActivityType=ACTIVITY; steps=25 | — |
| `TreatyInXOLAddSpreadingDetail.xml` | `ASM-FW-GISFW-DATA-TREATYINSHARE` | `TREATYINXOLADDSPREADINGDETAIL` | [terverifikasi] pyActivityType=ACTIVITY; steps=32 | `DATA-PORTAL / TREATYINXOLADDSPREADING` |
| `TreatyInXOLAddSpreadingDetailActual.xml` | `ASM-FW-GISFW-DATA-TREATYINSHARE` | `TREATYINXOLADDSPREADINGDETAILACTUAL` | [terverifikasi] pyActivityType=ACTIVITY; steps=32 | `ASM-FW-GISFW-DATA-TREATYINSHARE / TREATYINXOLADDSPREADINGDETAIL` |
| `TreatyInitAttach.xml` | `@BASECLASS` | `TREATYINITATTACH` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / SFAGISINITATTACH` |
| `TreatySaveAttachment.xml` | `@BASECLASS` | `TREATYSAVEATTACHMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | — |
| `TreatySetReinstatement.xml` | `DATA-PORTAL` | `TREATYSETREINSTATEMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |

### ConnectREST — 1 rule (`RULE-CONNECT-REST`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `SERVICEGOOGLE` | [terverifikasi] pyServiceName=ServiceGoogle; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / GOOGLESTORAGE_UPLOAD` |

### DataTransform — 35 rule (`RULE-OBJ-MODEL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AddCob.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `ADDCOB` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `AddLimitRetentionCession.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `ADDLIMITRETENTIONCESSION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Akseptasi_DT.xml` | `DATA-PORTAL` | `AKSEPTASI_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CalculateReinstatement.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `CALCULATEREINSTATEMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CalculateReinstatementPct.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `CALCULATEREINSTATEMENTPCT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYINLIMITS / CALCULATEREINSTATEMENT` |
| `CountTotalPctSpead.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `COUNTTOTALPCTSPEAD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CountTotalPctSpreadXOL.xml` | `ASM-FW-GISFW-DATA-TREATYINSHARE` | `COUNTTOTALPCTSPREADXOL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DeleteCoB.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `DELETECOB` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ReCalculateReinstatement.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `RECALCULATEREINSTATEMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYINLIMITS / CALCULATEREINSTATEMENT` |
| `Reset_DT.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `RESET_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetCoB.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITLAYER` | `SETCOB` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetCoBID.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `SETCOBID` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetDetailsID.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `SETDETAILSID` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetIndexLayer_DT.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITLAYER` | `SETINDEXLAYER_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYINLIMITS / SETINDEXLAYER_DT` |
| `SetParamRetro.xml` | `ASM-FW-GISFW-DATA-TREATYINSHAREREINS` | `SETPARAMRETRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShareRetroFacultative_DT.xml` | `DATA-PORTAL` | `SHARERETROFACULTATIVE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyCalculateProratePct.xml` | `DATA-PORTAL` | `TREATYCALCULATEPRORATEPCT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInAddAccumulation.xml` | `DATA-PORTAL` | `TREATYINADDACCUMULATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / TREATYINADDRETENTIONPERIOD` |
| `TreatyInAddNew.xml` | `DATA-PORTAL` | `TREATYINADDNEW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInCedingSetToUppercase.xml` | `DATA-PORTAL` | `TREATYINCEDINGSETTOUPPERCASE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / SETTOUPPERCASE` |
| `TreatyInCopyConditions.xml` | `DATA-PORTAL` | `TREATYINCOPYCONDITIONS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInDeleteAccumulationLists.xml` | `DATA-PORTAL` | `TREATYINDELETEACCUMULATIONLISTS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInForceEdit.xml` | `DATA-PORTAL` | `TREATYINFORCEEDIT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInForceResolveComplete.xml` | `DATA-PORTAL` | `TREATYINFORCERESOLVECOMPLETE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInIDSetPyPortal.xml` | `DATA-PORTAL` | `TREATYINIDSETPYPORTAL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInRemoveLastComment.xml` | `DATA-PORTAL` | `TREATYINREMOVELASTCOMMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInReturntoInputor.xml` | `DATA-PORTAL` | `TREATYINRETURNTOINPUTOR` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / TREATYINFORCERESOLVECOMPLETE` |
| `TreatyInSetEdit.xml` | `DATA-PORTAL` | `TREATYINSETEDIT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / TREATYINSETADM` |
| `TreatyInSetPeriod.xml` | `DATA-PORTAL` | `TREATYINSETPERIOD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInSetReinsured.xml` | `DATA-PORTAL` | `TREATYINSETREINSURED` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInSetTreatyYear.xml` | `DATA-PORTAL` | `TREATYINSETTREATYYEAR` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInShowHide.xml` | `DATA-PORTAL` | `TREATYINSHOWHIDE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInShowHideFacShare.xml` | `DATA-PORTAL` | `TREATYINSHOWHIDEFACSHARE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyTypeSetIndex.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `TREATYTYPESETINDEX` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `setCategoryAttachment_DT.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SETCATEGORYATTACHMENT_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### DecisionTable — 1 rule (`RULE-DECLARE-DECISIONTABLE`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `GetMimeType.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETMIMETYPE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### FlowAction — 32 rule (`RULE-OBJ-FLOWACTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AchievementCombine.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `ACHIEVEMENTCOMBINE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Action.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `ACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CoBList.xml` | `ASM-FW-GISFW-DATA-TREATYINSHARE` | `COBLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CoBListReadOnly.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITLAYER` | `COBLISTREADONLY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYINLIMITLAYER / COBLIST` |
| `DetailEGNPI.xml` | `ASM-FW-GISFW-DATA-TREATYINEGNPI` | `DETAILEGNPI` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DetailLimits.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `DETAILLIMITS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DetailShare.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `DETAILSHARE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DetailShareRetro.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `DETAILSHARERETRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL / DETAILSHARE` |
| `Installments.xml` | `ASM-FW-GISFW-DATA-TREATYININSTALLMENT` | `INSTALLMENTS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Layers.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `LAYERS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `LayersEDM.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `LAYERSEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `LimitFacRetro.xml` | `ASM-FW-GISFW-DATA-TREATYINSHAREREINS` | `LIMITFACRETRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `LimitProportional.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `LIMITPROPORTIONAL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `MaxRetention.xml` | `ASM-FW-GISFW-DATA-TREATYINRETENTION` | `MAXRETENTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Share.xml` | `ASM-FW-GISFW-DATA-TREATYINSHARE` | `SHARE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShareOldData.xml` | `ASM-FW-GISFW-DATA-TREATYINSHARE` | `SHAREOLDDATA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYINSHARE / SHARE` |
| `ShareRetro.xml` | `ASM-FW-GISFW-DATA-TREATYINSHARE` | `SHARERETRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYINSHARE / SHARE` |
| `ShowSummary.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `SHOWSUMMARY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SpreadingTPDtl.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSSPREADING` | `SPREADINGTPDTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYINLIMITSSPREADING / SPREADINGTONPDTL` |
| `SpreadingTXOLDtl.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSSPREADING` | `SPREADINGTXOLDTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYINLIMITSSPREADING / SPREADINGTPDTL` |
| `TotalLimits.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `TOTALLIMITS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TotalLimitsRetro.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `TOTALLIMITSRETRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYINLIMITS / TOTALLIMITS` |
| `TreatyAttachContent.xml` | `@BASECLASS` | `TREATYATTACHCONTENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / SFAGISATTACHCONTENT` |
| `TreatyInAction.xml` | `DATA-PORTAL` | `TREATYINACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInActionEDM.xml` | `DATA-PORTAL` | `TREATYINACTIONEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / TREATYINACTION` |
| `TreatyInDeclineConfirmation.xml` | `DATA-PORTAL` | `TREATYINDECLINECONFIRMATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInDeclineConfirmationEDM.xml` | `DATA-PORTAL` | `TREATYINDECLINECONFIRMATIONEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / TREATYINDECLINECONFIRMATION` |
| `TreatyInSearchReinsured.xml` | `DATA-PORTAL` | `TREATYINSEARCHREINSURED` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInSearchSoB.xml` | `DATA-PORTAL` | `TREATYINSEARCHSOB` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / TREATYINSEARCHREINSURED` |
| `TreatyInUploadCSV.xml` | `DATA-PORTAL` | `TREATYINUPLOADCSV` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyRetroList.xml` | `ASM-FW-GISFW-DATA-TREATYINRETROSHARE` | `TREATYRETROLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyinChooseRetro.xml` | `ASM-FW-GISFW-DATA-TREATYINSHAREREINS` | `TREATYINCHOOSERETRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Harness — 3 rule (`RULE-HTML-HARNESS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `InputTreatyInOffer.xml` | `DATA-PORTAL` | `INPUTTREATYINOFFER` | [terverifikasi] pyInclude=33 | — |
| `ShowAttachmentTreaty.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `SHOWATTACHMENTTREATY` | [terverifikasi] pyInclude=1 | — |
| `TreatyInFacultativeShareCalculation.xml` | `DATA-PORTAL` | `TREATYINFACULTATIVESHARECALCULATION` | [terverifikasi] pyInclude=1 | — |

### RDBList — 40 rule (`RULE-CONNECT-SQL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AttachmentLife.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `ATTACHMENTLIFE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CATEGORY_ATTACH_REAS | — |
| `BrowseOffer.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `BROWSEOFFER(CONVERT)` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.JSON_OFFER | `ASM-FW-GISFW-INT-TREATY_IN / ASM!BROWSETREATYIN` |
| `BrowseTreatyIn.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `BROWSETREATYIN` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_TREATY_IN,POOLDATA.M_TREATY_IN_EDM | `ASM-FW-GISFW-INT-PRODUCT_LIFE / ASM!BROWSEUNDERWRITINGLIST` |
| `CategoryAttach_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `CATEGORYATTACH_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CATEGORY_ATTACH_REAS | — |
| `ChangeKateAttachment2_Sql.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `CHANGEKATEATTACHMENT2_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT,UPDATE; tables=M_ATTACHMENTTREATY_2,POOLDATA.M_KATEGORIMASTERTREATY | `ASM-FW-GISFW-INT-TREATY_IN / ASM!DELETEATTACHMENT2_SQL` |
| `CountDocTreaty_SQL.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `COUNTDOCTREATY_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_ATTACHMENTTREATY_2 | `ASM-FW-GISFW-INT-TREATY_IN / ASM!GETMASTERTREATYCATEGORY_SQL` |
| `DeleteAttachment2_Sql.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `DELETEATTACHMENT2_SQL` | [terverifikasi] sqlKind=QUERY | `ASM-FW-GISFW-INT-TREATY_IN / ASM!GETATTACHMENT2_SQL` |
| `DeleteStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `DELETESTORAGE_SQL` | [terverifikasi] sqlKind=QUERY | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GETLINKSTORAGE_SQL` |
| `EDMCheckExistingData.xml` | `ASM-FW-GISFW-INT-TREATY_IN_EDM` | `EDMCHECKEXISTINGDATA` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATY_IN_EDM | — |
| `FetchTreatyContractDetail.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `FETCHTREATYCONTRACTDETAIL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=PROPORTIONALARRG | — |
| `FetchTreatyInProductionUsingNooffer.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `FETCHTREATYINPRODUCTIONUSINGNOOFFER` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TREATYINPRODUCTION | — |
| `GenerateImageID_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GENERATEIMAGEID_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GetAchievement.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `GETACHIEVEMENT` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.ACHIEVEMENT | — |
| `GetAllAttachment2_Sql.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `GETALLATTACHMENT2_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_ATTACHMENTTREATY_2 | `ASM-FW-GISFW-INT-TREATY_IN / ASM!GETATTACHMENT2_SQL` |
| `GetAppName_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETAPPNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.T_FOLDER_IMAGE | — |
| `GetAttachment2_Sql.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `GETATTACHMENT2_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_ATTACHMENTTREATY_2 | `ASM-FW-GISFW-INT-TREATY_IN / ASM!GETATTACHMENT_SQL` |
| `GetCountClaim.xml` | `ASSIGN-WORKLIST` | `GETCOUNTCLAIM` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=DATAPEGA.PC_ASM_FW_GCNMFW_WORK | — |
| `GetCurrencyToIDR_SQL.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `GETCURRENCYTOIDR_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYEXCHANGEYEARLY | — |
| `GetCurrentDate.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `GETCURRENTDATE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GetIDClaimAchievement.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `GETIDCLAIMACHIEVEMENT` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.OS_AKSEPTASI_KLAIM,POOLDATA.TREATYINPRODUCTION | — |
| `GetIDClaimEstimationAchievement.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `GETIDCLAIMESTIMATIONACHIEVEMENT` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.OS_AKSEPTASI_KLAIM,POOLDATA.TREATYINPRODUCTION | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL / ASM!GETIDCLAIMACHIEVEMENT` |
| `GetLinkStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETLINKSTORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GETTOKENSTORAGE_SQL` |
| `GetMasterTreatyCategory_SQL.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `GETMASTERTREATYCATEGORY_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_KATEGORIMASTERTREATY | — |
| `GetQuarter.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `GETQUARTER` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.ACHIEVEMENT | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL / ASM!GETACHIEVEMENT` |
| `GetQuarterYear.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `GETQUARTERYEAR` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.ACHIEVEMENT | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL / ASM!GETQUARTER` |
| `GetTokenStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETTOKENSTORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.GET_TOKEN_STORAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |
| `InsertAttachment2_Sql.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `INSERTATTACHMENT2_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=M_ATTACHMENTTREATY_2 | `ASM-FW-GISFW-INT-TREATY_IN / ASM!GETATTACHMENT_SQL` |
| `InsertToLogAchievement_SQL.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `INSERTTOLOGACHIEVEMENT_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=INSERT; tables=POOLDATA.LOG_ACHIEVEMENT | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL / ASM!GETQUARTERYEAR` |
| `InsertToTable.xml` | `ASM-FW-GISFW-INT-TREATYINOFFER` | `INSERTTOTABLE` | [terverifikasi] sqlKind=QUERY; sqlOps=INSERT; tables=POOLDATA.TREATYINOFFER | — |
| `Insert_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERT_T_STORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GENERATEIMAGEID_SQL` |
| `RemoveTreatyInDetail.xml` | `ASM-FW-GISFW-INT-TREATYINDETAIL` | `REMOVETREATYINDETAIL` | [terverifikasi] sqlKind=PLSQL; sqlOps=DELETE; tables=TREATYINDETAIL,M_TREATY_IN_DETAIL | — |
| `RemoveTreatyInDetailEdm.xml` | `ASM-FW-GISFW-INT-TREATYINDETAIL` | `REMOVETREATYINDETAILEDM` | [terverifikasi] sqlKind=PLSQL; sqlOps=DELETE; tables=TREATYINDETAILEDM,M_TREATY_IN_DETAIL_EDM | `ASM-FW-GISFW-INT-TREATYINDETAIL / ASM!REMOVETREATYINDETAIL` |
| `RemoveTreatyInEDM.xml` | `ASM-FW-GISFW-INT-TREATY_IN_EDM` | `REMOVETREATYINEDM` | [terverifikasi] sqlKind=QUERY; sqlOps=DELETE; tables=M_TREATY_IN_EDM | — |
| `RemoveTreatyInEDM2.xml` | `ASM-FW-GISFW-INT-TREATY_IN_EDM` | `REMOVETREATYINEDM2` | [terverifikasi] sqlKind=QUERY; sqlOps=DELETE; tables=TREATY_IN_EDM | `ASM-FW-GISFW-INT-TREATY_IN_EDM / RNM!REMOVETREATYINEDM` |
| `SaveTreatyIn.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `SAVETREATYIN` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_TREATY_IN | — |
| `SaveTreatyInDetail.xml` | `ASM-FW-GISFW-INT-TREATYINDETAIL` | `SAVETREATYINDETAIL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_M_TREATY_IN_DETAIL | `ASM-FW-GISFW-INT-TREATY_IN / ASM!SAVETREATYIN` |
| `SaveTreatyInDetailEdm.xml` | `ASM-FW-GISFW-INT-TREATYINDETAIL` | `SAVETREATYINDETAILEDM` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_M_TREATY_IN_DETAIL_EDM | `ASM-FW-GISFW-INT-TREATYINDETAIL / ASM!SAVETREATYINDETAIL` |
| `SaveTreatyInEDM.xml` | `ASM-FW-GISFW-INT-TREATY_IN_EDM` | `SAVETREATYINEDM` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_M_TREATY_IN_EDM | `ASM-FW-GISFW-INT-TREATY_IN / ASM!SAVETREATYIN` |
| `TreatyInSelectAll.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `TREATYINSELECTALL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_TREATY_IN | — |
| `Update_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `UPDATE_T_STORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |

### ReportDefinition — 16 rule (`RULE-OBJ-REPORT-DEFINITION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseAgentNusaRe_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSEAGENTNUSARE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseBusinessGroup_RD.xml` | `ASM-FW-GISFW-INT-BUSINESSGROUP` | `BROWSEBUSINESSGROUP_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseBusiness_RD.xml` | `ASM-FW-GISFW-INT-BUSINESS` | `BROWSEBUSINESS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseClientNusaRe_RD.xml` | `ASM-FW-GISFW-INT-CLIENT` | `BROWSECLIENTNUSARE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCurrencyTreatyIn_RD.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `BROWSECURRENCYTREATYIN_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCurrency_RD.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `BROWSECURRENCY_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseDetailTreatyReisurer_RD.xml` | `ASM-FW-GISFW-INT-TREATYREINSURER` | `BROWSEDETAILTREATYREISURER_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-TREATYREINSURER / BROWSEDETAILTREAYREISURER_RD` |
| `BrowseReinsuranceType_RD.xml` | `ASM-FW-GISFW-INT-REINSURANCETYPE` | `BROWSEREINSURANCETYPE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTREATY_IN.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `BROWSETREATY_IN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyArrangement_Limit_MstTrt_RD.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRANGEMENT_LIMIT_MSTTRT_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyArrangement_Limit_RD.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRANGEMENT_LIMIT_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyArrangement_ParentReinsMasterTrt.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRANGEMENT_PARENTREINSMASTERTRT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyBusinessWOType_RD.xml` | `ASM-FW-GISFW-INT-TREATYBUSINESS` | `BROWSETREATYBUSINESSWOTYPE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-TREATYBUSINESS / BROWSETREATYBUSINESS_RD` |
| `BrowseTreatyExchangeYearly_RD.xml` | `ASM-FW-GISFW-INT-TREATYEXCHANGEYEARLY` | `BROWSETREATYEXCHANGEYEARLY_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyGroup_RD.xml` | `ASM-FW-GISFW-INT-TREATYGROUP` | `BROWSETREATYGROUP_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyTypeParent.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYTYPEPARENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PROPORTIONALARRG / BROWSETREATYARRANGEMENT_PARENTREINSMASTERTRT` |

### Section — 50 rule (`RULE-HTML-SECTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AchievementCombine.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `ACHIEVEMENTCOMBINE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Action.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `ACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CoBList.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITLAYER` | `COBLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CoBListReadOnly.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITLAYER` | `COBLISTREADONLY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYINLIMITLAYER / COBLIST` |
| `DetailEGNPI.xml` | `ASM-FW-GISFW-DATA-TREATYINEGNPI` | `DETAILEGNPI` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DetailLimits.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `DETAILLIMITS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DetailShare.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `DETAILSHARE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DetailShareRetro.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` | `DETAILSHARERETRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL / DETAILSHARE` |
| `InputTreatyInOffer.xml` | `DATA-PORTAL` | `INPUTTREATYINOFFER` | [terverifikasi] pyInclude=32 | — |
| `Installments.xml` | `ASM-FW-GISFW-DATA-TREATYININSTALLMENT` | `INSTALLMENTS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Layers.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `LAYERS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `LayersEDM.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `LAYERSEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYINLIMITS / LAYERS` |
| `LimitFacRetro_Sec.xml` | `ASM-FW-GISFW-DATA-TREATYINSHAREREINS` | `LIMITFACRETRO_SEC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `LimitProportional.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `LIMITPROPORTIONAL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `MaxRetention.xml` | `ASM-FW-GISFW-DATA-TREATYINRETENTION` | `MAXRETENTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Share.xml` | `ASM-FW-GISFW-DATA-TREATYINSHARE` | `SHARE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShareOldData.xml` | `ASM-FW-GISFW-DATA-TREATYINSHARE` | `SHAREOLDDATA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYINSHARE / SHARE` |
| `ShareRetro.xml` | `ASM-FW-GISFW-DATA-TREATYINSHARE` | `SHARERETRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYINSHARE / SHARE` |
| `ShowAttachmentTreaty.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `SHOWATTACHMENTTREATY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / SHOWATTACHMENTTREATY` |
| `ShowSummary.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `SHOWSUMMARY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SpreadingTPDtl.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSSPREADING` | `SPREADINGTPDTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYINLIMITSSPREADING / SPREADINGTONPDTL` |
| `SpreadingTXOLDtl.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITSSPREADING` | `SPREADINGTXOLDTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYINLIMITSSPREADING / SPREADINGTPDTL` |
| `TotalLimits.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `TOTALLIMITS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TotalLimitsRetro.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `TOTALLIMITSRETRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYINLIMITS / TOTALLIMITS` |
| `TreatyInAction.xml` | `DATA-PORTAL` | `TREATYINACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInActionButtons.xml` | `DATA-PORTAL` | `TREATYINACTIONBUTTONS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInActionEDM.xml` | `DATA-PORTAL` | `TREATYINACTIONEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / TREATYINACTION` |
| `TreatyInActualFacultativeShareCalculation.xml` | `DATA-PORTAL` | `TREATYINACTUALFACULTATIVESHARECALCULATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / TREATYINFACULTATIVESHARECALCULATION` |
| `TreatyInActualLimits.xml` | `DATA-PORTAL` | `TREATYINACTUALLIMITS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / TREATYINACTUALSHARE` |
| `TreatyInActualShare.xml` | `DATA-PORTAL` | `TREATYINACTUALSHARE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInActualSumary.xml` | `DATA-PORTAL` | `TREATYINACTUALSUMARY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInDeclineConfirmation.xml` | `DATA-PORTAL` | `TREATYINDECLINECONFIRMATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInDeclineConfirmationEDM.xml` | `DATA-PORTAL` | `TREATYINDECLINECONFIRMATIONEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / TREATYINDECLINECONFIRMATION` |
| `TreatyInFacultativeRetro.xml` | `DATA-PORTAL` | `TREATYINFACULTATIVERETRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInFacultativeShareCalculation.xml` | `DATA-PORTAL` | `TREATYINFACULTATIVESHARECALCULATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInNONProportional.xml` | `DATA-PORTAL` | `TREATYINNONPROPORTIONAL` | [terverifikasi] pyInclude=27 | — |
| `TreatyInSearchReinsured.xml` | `DATA-PORTAL` | `TREATYINSEARCHREINSURED` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInSearchSoB.xml` | `DATA-PORTAL` | `TREATYINSEARCHSOB` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / TREATYINSEARCHREINSURED` |
| `TreatyInShareProp.xml` | `DATA-PORTAL` | `TREATYINSHAREPROP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInTabsAchievement.xml` | `DATA-PORTAL` | `TREATYINTABSACHIEVEMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyInTabsNPValueDifferenceProRate.xml` | `DATA-PORTAL` | `TREATYINTABSNPVALUEDIFFERENCEPRORATE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / TREATYINTABSNONPROPORTIONALVALUEDIFFERENCE` |
| `TreatyInTabsNPValueDifference_NoProRate.xml` | `DATA-PORTAL` | `TREATYINTABSNPVALUEDIFFERENCE_NOPRORATE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / TREATYINTABSNPVALUEDIFFERENCEPRORATE` |
| `TreatyInTabsNonProportional.xml` | `DATA-PORTAL` | `TREATYINTABSNONPROPORTIONAL` | [terverifikasi] pyInclude=6 | — |
| `TreatyInTabsNonProportionalAdjustPremi.xml` | `DATA-PORTAL` | `TREATYINTABSNONPROPORTIONALADJUSTPREMI` | [terverifikasi] pyInclude=10 | `DATA-PORTAL / TREATYINTABSNONPROPORTIONAL` |
| `TreatyInTabsNonProportionalValueDifference.xml` | `DATA-PORTAL` | `TREATYINTABSNONPROPORTIONALVALUEDIFFERENCE` | [terverifikasi] pyInclude=6 | `DATA-PORTAL / TREATYINTABSNONPROPORTIONAL` |
| `TreatyInTabsProportional.xml` | `DATA-PORTAL` | `TREATYINTABSPROPORTIONAL` | [terverifikasi] pyInclude=8 | — |
| `TreatyInfoSubmit.xml` | `DATA-PORTAL` | `TREATYINFOSUBMIT` | [terverifikasi] pyInclude=2 | — |
| `TreatyRetroList.xml` | `ASM-FW-GISFW-DATA-TREATYINRETROSHARE` | `TREATYRETROLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyinChooseRetro_Sec.xml` | `ASM-FW-GISFW-DATA-TREATYINSHAREREINS` | `TREATYINCHOOSERETRO_SEC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `WorkAttachments.xml` | `DATA-PORTAL` | `WORKATTACHMENTS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### SystemSettings — 1 rule (`RULE-ADMIN-SYSTEM-SETTINGS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `LinkService.xml` | `LINKSERVICE` | `LINKSERVICE` | [terverifikasi] pyPurpose=LinkService | — |

### When — 1 rule (`RULE-OBJ-WHEN`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `TreatyMasterInEDM.xml` | `DATA-PORTAL` | `TREATYMASTERINEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

## 4. Tabrakan nama di dalam modul ini

### 4.1 Nama file dipakai di lebih dari satu tipe rule

| Nama file | Tipe | Class |
| --- | --- | --- |
| `SpreadingTPDtl.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINLIMITSSPREADING, ASM-FW-GISFW-DATA-TREATYINLIMITSSPREADING` |
| `TreatyInSearchReinsured.xml` | FlowAction, Section | `DATA-PORTAL, DATA-PORTAL` |
| `TreatyInActionEDM.xml` | FlowAction, Section | `DATA-PORTAL, DATA-PORTAL` |
| `LimitProportional.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINLIMITS, ASM-FW-GISFW-DATA-TREATYINLIMITS` |
| `Layers.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINLIMITS, ASM-FW-GISFW-DATA-TREATYINLIMITS` |
| `TreatyInDeclineConfirmation.xml` | FlowAction, Section | `DATA-PORTAL, DATA-PORTAL` |
| `ShareOldData.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINSHARE, ASM-FW-GISFW-DATA-TREATYINSHARE` |
| `Share.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINSHARE, ASM-FW-GISFW-DATA-TREATYINSHARE` |
| `GetAchievement.xml` | Activity, RDBList | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL, ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` |
| `LayersEDM.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINLIMITS, ASM-FW-GISFW-DATA-TREATYINLIMITS` |
| `TreatyInSearchSoB.xml` | FlowAction, Section | `DATA-PORTAL, DATA-PORTAL` |
| `TreatyInDeclineConfirmationEDM.xml` | FlowAction, Section | `DATA-PORTAL, DATA-PORTAL` |
| `Installments.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYININSTALLMENT, ASM-FW-GISFW-DATA-TREATYININSTALLMENT` |
| `DetailShareRetro.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL, ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` |
| `SpreadingTXOLDtl.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINLIMITSSPREADING, ASM-FW-GISFW-DATA-TREATYINLIMITSSPREADING` |
| `TreatyInFacultativeShareCalculation.xml` | Harness, Section | `DATA-PORTAL, DATA-PORTAL` |
| `Action.xml` | FlowAction, Section | `ASM-FW-GISFW-INT-TREATY_IN, ASM-FW-GISFW-INT-TREATY_IN` |
| `ShowAttachmentTreaty.xml` | Harness, Section | `ASM-FW-GISFW-INT-TREATY_IN, ASM-FW-GISFW-INT-TREATY_IN` |
| `TreatyInActualShare.xml` | Activity, Section | `DATA-PORTAL, DATA-PORTAL` |
| `DetailEGNPI.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINEGNPI, ASM-FW-GISFW-DATA-TREATYINEGNPI` |
| `DetailShare.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL, ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` |
| `CoBList.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINSHARE, ASM-FW-GISFW-DATA-TREATYINLIMITLAYER` |
| `TreatyInAction.xml` | FlowAction, Section | `DATA-PORTAL, DATA-PORTAL` |
| `CoBListReadOnly.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINLIMITLAYER, ASM-FW-GISFW-DATA-TREATYINLIMITLAYER` |
| `ShowSummary.xml` | FlowAction, Section | `ASM-FW-GISFW-INT-TREATY_IN, ASM-FW-GISFW-INT-TREATY_IN` |
| `FetchTreatyContractDetail.xml` | Activity, RDBList | `ASM-FW-GISFW-DATA-TREATYINRETROSHARE, ASM-FW-GISFW-INT-TREATY_IN` |
| `MaxRetention.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINRETENTION, ASM-FW-GISFW-DATA-TREATYINRETENTION` |
| `TreatyRetroList.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINRETROSHARE, ASM-FW-GISFW-DATA-TREATYINRETROSHARE` |
| `TotalLimitsRetro.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINLIMITS, ASM-FW-GISFW-DATA-TREATYINLIMITS` |
| `DetailLimits.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL, ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL` |
| `AchievementCombine.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINLIMITS, ASM-FW-GISFW-DATA-TREATYINLIMITS` |
| `ShareRetro.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINSHARE, ASM-FW-GISFW-DATA-TREATYINSHARE` |
| `TotalLimits.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINLIMITS, ASM-FW-GISFW-DATA-TREATYINLIMITS` |
| `InputTreatyInOffer.xml` | Harness, Section | `DATA-PORTAL, DATA-PORTAL` |
| `TreatyInUploadCSV.xml` | Activity, FlowAction | `DATA-PORTAL, DATA-PORTAL` |

Total: **35** nama file.

### 4.2 Rule `When` bernama sama di lebih dari satu class

_Tidak ada._

## 5. Keluarga klon (satu asal salinan dipakai >1 rule)

Dibaca dari `<pzOriginalInstanceKey>` yang sama pada beberapa file. Ini **sinyal kandidat
konsolidasi** saat migrasi, bukan bukti bahwa perilakunya identik — kesamaan perilaku baru
dapat dinyatakan setelah isinya dibandingkan di STEP D2.

| Asal salinan (class / nama) | Jumlah | File |
| --- | ---: | --- |
| `DATA-PORTAL / TREATYINUPLOADCSV` | 2 | `Activity/TreatyInUploadCSV.xml`, `FlowAction/TreatyInUploadCSV.xml` |
| `ASM-FW-GISFW-INT-TREATY_IN / ASM!GETMASTERTREATYCATEGORY_SQL` | 2 | `RDBList/CountDocTreaty_SQL.xml`, `RDBList/GetMasterTreatyCategory_SQL.xml` |
| `DATA-PORTAL / SETTREATYIN_ACT` | 2 | `Activity/SetTreatyIn_Act.xml`, `Activity/TreatyInConvertCallData_act.xml` |
| `ASM-FW-GISFW-DATA-TREATYINLIMITS / ACHIEVEMENTCOMBINE` | 2 | `FlowAction/AchievementCombine.xml`, `Section/AchievementCombine.xml` |
| `DATA-PORTAL / TREATYINDIFFERENCELIMITS` | 2 | `Activity/TreatyEDMDifferenceLimits.xml`, `Activity/TreatyInDifferenceLimits.xml` |
| `DATA-PORTAL / TREATYINSETBROKERAGE` | 3 | `Activity/TreatyInDifferenceDeduction.xml`, `Activity/TreatyInSetBrokerage.xml`, `Activity/TreatyInSetBrokerageActual.xml` |
| `DATA-PORTAL / TREATYINXOLADDSPREADING` | 2 | `Activity/TreatyInXOLAddSpreading.xml`, `Activity/TreatyInXOLAddSpreadingDetail.xml` |
| `ASM-FW-GISFW-INT-TREATY_IN / ASM!GETATTACHMENT2_SQL` | 2 | `RDBList/DeleteAttachment2_Sql.xml`, `RDBList/GetAllAttachment2_Sql.xml` |
| `ASM-FW-GISFW-DATA-TREATYINLIMITS / LIMITPROPORTIONAL` | 2 | `FlowAction/LimitProportional.xml`, `Section/LimitProportional.xml` |
| `ASM-FW-GISFW-DATA-TREATYINEGNPI / DETAILEGNPI` | 2 | `FlowAction/DetailEGNPI.xml`, `Section/DetailEGNPI.xml` |
| `DATA-PORTAL / TREATYINSETREPORT` | 2 | `Activity/TreatyInSetAccountReport.xml`, `Activity/TreatyInSetReport.xml` |
| `ASM-FW-GISFW-DATA-TREATYINLIMITS / LAYERS` | 3 | `FlowAction/Layers.xml`, `Section/Layers.xml`, `Section/LayersEDM.xml` |
| `DATA-PORTAL / TREATYINTABSNONPROPORTIONAL` | 3 | `Section/TreatyInTabsNonProportional.xml`, `Section/TreatyInTabsNonProportionalAdjustPremi.xml`, `Section/TreatyInTabsNonProportionalValueDifference.xml` |
| `ASM-FW-GISFW-DATA-TREATYINLIMITSSPREADING / SPREADINGTPDTL` | 2 | `FlowAction/SpreadingTXOLDtl.xml`, `Section/SpreadingTXOLDtl.xml` |
| `ASM-FW-GISFW-INT-TREATYINDETAIL / ASM!REMOVETREATYINDETAIL` | 2 | `RDBList/RemoveTreatyInDetail.xml`, `RDBList/RemoveTreatyInDetailEdm.xml` |
| `DATA-PORTAL / TREATYINFACULTATIVESHARECALCULATION` | 3 | `Harness/TreatyInFacultativeShareCalculation.xml`, `Section/TreatyInActualFacultativeShareCalculation.xml`, `Section/TreatyInFacultativeShareCalculation.xml` |
| `ASM-FW-GISFW-INT-TREATY_IN_EDM / RNM!REMOVETREATYINEDM` | 2 | `RDBList/RemoveTreatyInEDM.xml`, `RDBList/RemoveTreatyInEDM2.xml` |
| `DATA-PORTAL / SAVETREATYINOFFER_ACT` | 2 | `Activity/SaveTreatyInOfferNonProp1_Act.xml`, `Activity/SaveTreatyInOffer_Act.xml` |
| `DATA-PORTAL / TREATYINFORCERESOLVECOMPLETE` | 2 | `DataTransform/TreatyInForceResolveComplete.xml`, `DataTransform/TreatyInReturntoInputor.xml` |
| `DATA-PORTAL / TREATYINNPSETTOTALACTUAL` | 2 | `Activity/TreatyInNPSetTotalActualLimits.xml`, `Activity/TreatyInNPSetTotalActualShare.xml` |
| `DATA-PORTAL / TREATYINDECLINECONFIRMATION` | 4 | `FlowAction/TreatyInDeclineConfirmation.xml`, `FlowAction/TreatyInDeclineConfirmationEDM.xml`, `Section/TreatyInDeclineConfirmation.xml`, `Section/TreatyInDeclineConfirmationEDM.xml` |
| `DATA-PORTAL / TREATYINPOPULATEDETAIL` | 2 | `Activity/TreatyInPopulateDetail.xml`, `Activity/TreatyInSetValue.xml` |
| `ASM-FW-GISFW-DATA-TREATYINSHARE / FETCHQSFROMMASTERXOL` | 2 | `Activity/FetchQSfromMasterXOL.xml`, `Activity/GetNilaiTotal.xml` |
| `ASM-FW-GISFW-INT-PROPORTIONALARRG / BROWSETREATYARRANGEMENT_PARENTREINSMASTERTRT` | 2 | `ReportDefinition/BrowseTreatyArrangement_ParentReinsMasterTrt.xml`, `ReportDefinition/BrowseTreatyTypeParent.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / INSERTGOOGLESTORAGE_ACT` | 2 | `Activity/GetUrlGoogleStorage_Act.xml`, `Activity/InsertGoogleStorage_Act.xml` |
| `DATA-PORTAL / TREATYINDECLINECONFIRMATION_POSTACT` | 2 | `Activity/TreatyInDeclineConfirmation_postact.xml`, `Activity/TreatyInDeclineConfirmation_postactEDM.xml` |
| `ASM-FW-GISFW-DATA-TREATYINRETENTION / MAXRETENTION` | 2 | `FlowAction/MaxRetention.xml`, `Section/MaxRetention.xml` |
| `DATA-PORTAL / TREATYINSETVALUEINSTALLMENT` | 2 | `Activity/TreatyInSetValueDifferenceInstallment.xml`, `Activity/TreatyInSetValueInstallment.xml` |
| `DATA-PORTAL / TREATYINDIFFERENCEPREMIUM` | 2 | `Activity/TreatyEDMDifferencePremium.xml`, `Activity/TreatyInDifferencePremium.xml` |
| `DATA-PORTAL / TREATYINSEARCHREINSURED` | 4 | `FlowAction/TreatyInSearchReinsured.xml`, `FlowAction/TreatyInSearchSoB.xml`, `Section/TreatyInSearchReinsured.xml`, `Section/TreatyInSearchSoB.xml` |
| `DATA-PORTAL / TREATYINACTION` | 4 | `FlowAction/TreatyInAction.xml`, `FlowAction/TreatyInActionEDM.xml`, `Section/TreatyInAction.xml`, `Section/TreatyInActionEDM.xml` |
| `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL / DETAILLIMITS` | 2 | `FlowAction/DetailLimits.xml`, `Section/DetailLimits.xml` |
| `ASM-FW-GISFW-INT-TREATY_IN / ACTION` | 2 | `FlowAction/Action.xml`, `Section/Action.xml` |
| `DATA-PORTAL / TREATYINACTUALSHARE` | 3 | `Activity/TreatyInActualShare.xml`, `Section/TreatyInActualLimits.xml`, `Section/TreatyInActualShare.xml` |
| `ASM-FW-GISFW-DATA-TREATYINCURRENCYLIST / SETCURRNAMEMASTERTREATY_ACT` | 2 | `Activity/SetCurrNameMasterTreaty_Act.xml`, `Activity/SetCurrName_Act.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GENERATEIMAGEID_SQL` | 2 | `RDBList/GenerateImageID_SQL.xml`, `RDBList/Insert_T_Storage_SQL.xml` |
| `ASM-FW-GISFW-DATA-TREATYINLIMITSSPREADING / SPREADINGTONPDTL` | 2 | `FlowAction/SpreadingTPDtl.xml`, `Section/SpreadingTPDtl.xml` |
| `DATA-PORTAL / TREATYINSUMMARYLIMITSHARE` | 4 | `Activity/CalculateRetroSumary.xml`, `Activity/TreatyInDifferenceShareSumary.xml`, `Activity/TreatyInSummaryLimitFacShare.xml`, `Activity/TreatyInSummaryLimitShareActual.xml` |
| `DATA-PORTAL / INPUTTREATYINOFFER` | 2 | `Harness/InputTreatyInOffer.xml`, `Section/InputTreatyInOffer.xml` |
| `ASM-FW-GISFW-DATA-TREATYINRETROSHARE / TREATYRETROLIST` | 2 | `FlowAction/TreatyRetroList.xml`, `Section/TreatyRetroList.xml` |
| `ASM-FW-GISFW-DATA-TREATYINLIMITS / CALCULATEREINSTATEMENT` | 3 | `DataTransform/CalculateReinstatement.xml`, `DataTransform/CalculateReinstatementPct.xml`, `DataTransform/ReCalculateReinstatement.xml` |
| `ASM-FW-GISFW-DATA-TREATYINLIMITLAYER / COBLIST` | 3 | `FlowAction/CoBListReadOnly.xml`, `Section/CoBList.xml`, `Section/CoBListReadOnly.xml` |
| `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL / DETAILSHARE` | 4 | `FlowAction/DetailShare.xml`, `FlowAction/DetailShareRetro.xml`, `Section/DetailShare.xml`, `Section/DetailShareRetro.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` | 2 | `RDBList/GetTokenStorage_SQL.xml`, `RDBList/Update_T_Storage_SQL.xml` |
| `ASM-FW-GISFW-INT-TREATY_IN / AKSEPTASI_ACT` | 2 | `Activity/Akseptasi_Act.xml`, `Activity/TreatyInAkseptasi_Act.xml` |
| `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL / ASM!GETACHIEVEMENT` | 2 | `RDBList/GetAchievement.xml`, `RDBList/GetQuarter.xml` |
| `ASM-FW-GISFW-DATA-TREATYININSTALLMENT / INSTALLMENTS` | 2 | `FlowAction/Installments.xml`, `Section/Installments.xml` |
| `ASM-FW-GISFW-INT-TREATY_IN / ASM!GETATTACHMENT_SQL` | 2 | `RDBList/GetAttachment2_Sql.xml`, `RDBList/InsertAttachment2_Sql.xml` |
| `DATA-PORTAL / TREATYINDIFFERENCESHARETOTAL` | 2 | `Activity/TreatyInDifferenceFacShareTotal.xml`, `Activity/TreatyInDifferenceShareTotal.xml` |
| `DATA-PORTAL / TREATYINDIFFERENCESHARE` | 3 | `Activity/TreatyEDMDifferenceShare.xml`, `Activity/TreatyInDifferenceFacShare.xml`, `Activity/TreatyInDifferenceShare.xml` |
| `DATA-PORTAL / SAVETREATYIN_ACT` | 2 | `Activity/SaveTreatyIn_Act.xml`, `Activity/SaveTreatyIn_EDM_Act.xml` |
| `ASM-FW-GISFW-DATA-TREATYINLIMITSDETAIL / ASM!GETIDCLAIMACHIEVEMENT` | 2 | `RDBList/GetIDClaimAchievement.xml`, `RDBList/GetIDClaimEstimationAchievement.xml` |
| `DATA-PORTAL / TREATYINNPSETTOTAL` | 3 | `Activity/TreatyInNPSetTotal.xml`, `Activity/TreatyInNPSetTotalActual.xml`, `Activity/TreatyInSummaryMDP.xml` |
| `ASM-FW-GISFW-INT-TREATY_IN / ASM!SAVETREATYIN` | 3 | `RDBList/SaveTreatyIn.xml`, `RDBList/SaveTreatyInDetail.xml`, `RDBList/SaveTreatyInEDM.xml` |
| `ASM-FW-GISFW-DATA-TREATYINSHARE / SHARE` | 6 | `FlowAction/Share.xml`, `FlowAction/ShareOldData.xml`, `FlowAction/ShareRetro.xml`, `Section/Share.xml`, `Section/ShareOldData.xml`, `Section/ShareRetro.xml` |
| `ASM-FW-GISFW-DATA-TREATYINLIMITS / TOTALLIMITS` | 4 | `FlowAction/TotalLimits.xml`, `FlowAction/TotalLimitsRetro.xml`, `Section/TotalLimits.xml`, `Section/TotalLimitsRetro.xml` |
| `DATA-PORTAL / TREATYINSUBMIT` | 2 | `Activity/TreatyInSubmit.xml`, `Activity/TreatyInSubmitEDM.xml` |
| `ASM-FW-GISFW-INT-TREATY_IN / SHOWSUMMARY` | 2 | `FlowAction/ShowSummary.xml`, `Section/ShowSummary.xml` |
| `DATA-PORTAL / TREATYINACTUALUPDATEVALUE` | 3 | `Activity/TreatyInActualUpdateValue.xml`, `Activity/TreatyInActualUpdateValueLimits.xml`, `Activity/TreatyInActualUpdateValueShare.xml` |
| `DATA-PORTAL / SAVETREATYINDETAIL_ACT` | 2 | `Activity/SaveTreatyInDetailEdm_Act.xml`, `Activity/SaveTreatyInDetail_Act.xml` |

## 6. Objek database yang disentuh

Diambil dari `<pyBrowseSQL>` pada rule `RDBList`. **Nama objek apa adanya dari SQL** —
kualifikasi skema tidak konsisten di korpus. **Tipe kolom, constraint, dan relasi tidak
diketahui** (OQ-001).

| Tabel/view (apa adanya di SQL) | Dirujuk oleh |
| --- | ---: |
| `M_ATTACHMENTTREATY_2` | 4 rule |
| `POOLDATA.ACHIEVEMENT` | 3 rule |
| `POOLDATA.TREATYINPRODUCTION` | 3 rule |
| `T_STORAGE_IMAGE` | 3 rule |
| `CATEGORY_ATTACH_REAS` | 2 rule |
| `POOLDATA.M_KATEGORIMASTERTREATY` | 2 rule |
| `POOLDATA.OS_AKSEPTASI_KLAIM` | 2 rule |
| `TREATY_IN_EDM` | 2 rule |
| `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` | 1 rule |
| `M_TREATY_IN` | 1 rule |
| `M_TREATY_IN_DETAIL` | 1 rule |
| `M_TREATY_IN_DETAIL_EDM` | 1 rule |
| `M_TREATY_IN_EDM` | 1 rule |
| `POOLDATA.JSON_OFFER` | 1 rule |
| `POOLDATA.LOG_ACHIEVEMENT` | 1 rule |
| `POOLDATA.M_ATTACHMENTTREATY_2` | 1 rule |
| `POOLDATA.M_TREATY_IN` | 1 rule |
| `POOLDATA.M_TREATY_IN_EDM` | 1 rule |
| `POOLDATA.TREATYINOFFER` | 1 rule |
| `POOLDATA.T_FOLDER_IMAGE` | 1 rule |
| `PROPORTIONALARRG` | 1 rule |
| `TREATYEXCHANGEYEARLY` | 1 rule |
| `TREATYINDETAIL` | 1 rule |
| `TREATYINDETAILEDM` | 1 rule |

### Stored procedure / function yang dipanggil

**Body prosedur ini TIDAK ADA di korpus** — logika di dalamnya tidak dapat direkonstruksi
dari XML. Lihat OQ-002.

| Procedure | Dipanggil oleh |
| --- | --- |
| `POOLDATA.PEGA_M_TREATY_IN_DETAIL` | `SaveTreatyInDetail.xml` |
| `POOLDATA.PEGA_M_TREATY_IN_DETAIL_EDM` | `SaveTreatyInDetailEdm.xml` |
| `POOLDATA.PEGA_TREATY_IN` | `SaveTreatyIn.xml` |
| `POOLDATA.PEGA_M_TREATY_IN_EDM` | `SaveTreatyInEDM.xml` |
| `POOLDATA.GET_TOKEN_STORAGE` | `GetTokenStorage_SQL.xml` |

## 7. Integrasi eksternal (`ConnectREST`)

Alamat tujuan dicatat sebagai **konfigurasi**, bukan URL literal. `baseURL=SETTING` berarti
alamat diambil dari rule `SystemSettings` yang disebut di kolom Setting; nilai sebenarnya
**tidak ada di korpus**. Kolom terakhir menandai rule yang justru memuat URL literal —
itu temuan, bukan fakta arsitektur (lihat `_summary.md`).

| File | Class | pyServiceName | Sumber base URL | Setting | URL literal di dalam rule? |
| --- | --- | --- | --- | --- | --- |
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `ServiceGoogle` | `SETTING` | `LinkService!LinkService` | tidak |

