# Inventaris Rule — NB Treaty In

STEP D1, batch 1 (domain treaty inward). Dibangkitkan 2026-09-12 dari korpus READ-ONLY
`D:\XML\RNM_BRD\NB Treaty In\`.

**Batas step ini:** hanya struktur + metadata. Tidak ada kesimpulan logika bisnis di sini.
Penelusuran perilaku adalah STEP D2 (`../flows/`).

**Identitas rule** diambil dari `<pxInsName>` (lihat `../README.md` §6.4).  
Kolom **Asal salinan** berisi `<pzOriginalInstanceKey>` dan hanya diisi bila **berbeda**
dari identitas — artinya rule ini hasil *Save As* dari rule lain. Tanda `—` berarti sama.
Class `—` berarti `<pxInsName>` tidak berisi bagian class (berlaku pada `DataPage`, yang
beridentitas global, bukan per class).

Perintah audit (dijalankan dari `D:\XML\RNM_BRD\`):

```
find "NB Treaty In" -type f -name "*.xml" | wc -l
find "NB Treaty In" -type f -name "*.xml" | awk -F/ '{print $2}' | sort | uniq -c
grep -o "<pxInsName>[^<]*" "NB Treaty In/<tipe>/<file>.xml" | head -1
grep -o "<pzOriginalInstanceKey>[^<]*" "NB Treaty In/<tipe>/<file>.xml" | head -1
```

## 1. Ringkasan terukur

- Total file XML: **278** — cocok dengan `../README.md` §5.1
- **Identitas rule unik: 278** dari 278 file → tidak ada identitas ganda di dalam modul ini
- Class Pega unik: **38**
- RDBList menurut jenis SQL: PLSQL=5, QUERY=36 (`PLSQL` = blok anonim `BEGIN…END;`, umumnya memanggil stored procedure)
- Prefix pada kunci RDBList: `ASM`=25, `RNM`=16 (arti prefix **belum terverifikasi** — OQ-008)

| Tipe rule | Jumlah |
| --- | ---: |
| Activity | 92 |
| DataTransform | 12 |
| DecisionTable | 2 |
| Flow | 1 |
| FlowAction | 9 |
| Harness | 6 |
| RDBList | 41 |
| ReportDefinition | 15 |
| Section | 25 |
| When | 75 |
| **TOTAL** | **278** |

## 2. Sebaran class Pega

Class yang **tidak** berawalan `ASM-` adalah class bawaan Pega atau class portal — ruang
lingkup migrasinya belum ditentukan (**OQ-009**).

| Class | Jumlah rule | Aplikasi? |
| --- | ---: | --- |
| `ASM-FW-GISFW-WORK` | 90 | aplikasi |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN` | 51 | aplikasi |
| `ASM-FW-GISFW-DATA` | 30 | aplikasi |
| `DATA-PORTAL` | 10 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-DATA-COVERAGE` | 9 | aplikasi |
| `@BASECLASS` | 8 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-DATA-QUOTATION` | 8 | aplikasi |
| `ASM-FW-GISFW-INT-POLICYJSON` | 7 | aplikasi |
| `ASM-FW-GISFW-INT-CURRENCY` | 6 | aplikasi |
| `ASM-FW-GISFW-INT-AGENT` | 5 | aplikasi |
| `ASM-FW-GISFW-INT-TREATY_IN` | 5 | aplikasi |
| `ASM-FW-GISFW-DATA-AGENT` | 4 | aplikasi |
| `ASM-FW-GISFW-INT-OFFERJSON` | 4 | aplikasi |
| `ASM-FW-GISFW-INT-PROPORTIONALARRG` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-REINSURANCETYPE` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYINDETAILJOINEDM` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYOUTDETAIL` | 3 | aplikasi |
| `ASM-FW-GISFW-DATA-INSTALLMENT` | 2 | aplikasi |
| `ASM-FW-GISFW-DATA-OFFERTREATYIN` | 2 | aplikasi |
| `ASM-FW-GISFW-DATA-TREATYININSTALLMENT` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-CLIENT` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-POLISTREATYIN` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYGROUP` | 2 | aplikasi |
| `ASM-FW-SFAGISFW-WORK-OPPORTUNITY` | 2 | aplikasi |
| `PEGACRM-PORTAL` | 2 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-DATA-OFFERFACIN` | 1 | aplikasi |
| `ASM-FW-GISFW-DATA-PROPERTYITEM` | 1 | aplikasi |
| `ASM-FW-GISFW-DATA-TREATYINLIMITS` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-MARKETINGOFFICER` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-OCCUPATION` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-RW` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYINDETAIL` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-V_JN_OBJ_ITEM` | 1 | aplikasi |
| `ASM-FW-GISFW-WORK-LIFE` | 1 | aplikasi |
| `ASM-FW-GISFW-WORK-NB` | 1 | aplikasi |
| `ASSIGN-WORKLIST` | 1 | **bukan aplikasi** → OQ-009 |
| `DATA-PARTY-PERSON` | 1 | **bukan aplikasi** → OQ-009 |
| `PEGACRM-WORK-` | 1 | **bukan aplikasi** → OQ-009 |

Rule di class bukan-aplikasi: **23** dari 278.

## 3. Daftar rule per tipe

### Activity — 92 rule (`RULE-OBJ-ACTIVITY`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AgentSourceBizTreatyIn_Act.xml` | `ASM-FW-GISFW-DATA-AGENT` | `AGENTSOURCEBIZTREATYIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `ASM-FW-GISFW-DATA-AGENT / AGENTSOURCEBIZ_ACT` |
| `BreakDownSpreading_Act.xml` | `ASM-FW-GISFW-WORK` | `BREAKDOWNSPREADING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CalculatePremi_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `CALCULATEPREMI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `CekLimitTreatyAcc_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `CEKLIMITTREATYACC_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CheckDataMkt.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `CHECKDATAMKT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `ASM-FW-GISFW-DATA-POLICYTREATYIN / CHECKDATAMO` |
| `CheckDuplicateOffer.xml` | `DATA-PORTAL` | `CHECKDUPLICATEOFFER` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `CheckRISLIP.xml` | `ASM-FW-GISFW-WORK` | `CHECKRISLIP` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CheckRISLIP_EDM.xml` | `ASM-FW-GISFW-WORK` | `CHECKRISLIP_EDM` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `ASM-FW-GISFW-WORK / CHECKRISLIP` |
| `CheckSpreadingProtectAnekaGolf_ACT.xml` | `ASM-FW-GISFW-WORK` | `CHECKSPREADINGPROTECTANEKAGOLF_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=43 | `ASM-FW-GISFW-WORK / CHECKSPREADINGPROTECT_ACT` |
| `CheckSpreadingProtectFire_ACT.xml` | `ASM-FW-GISFW-WORK` | `CHECKSPREADINGPROTECTFIRE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=39 | `ASM-FW-GISFW-WORK / CHECKSPREADINGPROTECT_ACT` |
| `CheckSpreadingProtectMCargoMBU_ACT.xml` | `ASM-FW-GISFW-WORK` | `CHECKSPREADINGPROTECTMCARGOMBU_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=31 | `ASM-FW-GISFW-WORK / CHECKSPREADINGPROTECT_ACT` |
| `CheckSpreadingProtectPATravel_ACT.xml` | `ASM-FW-GISFW-WORK` | `CHECKSPREADINGPROTECTPATRAVEL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=31 | `ASM-FW-GISFW-WORK / CHECKSPREADINGPROTECT_ACT` |
| `CheckSpreadingProtect_ACT.xml` | `ASM-FW-GISFW-WORK` | `CHECKSPREADINGPROTECT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=31 | — |
| `ConcatSlipOfferNo_Act.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `CONCATSLIPOFFERNO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `ConvertHistoryDate.xml` | `DATA-PORTAL` | `CONVERTHISTORYDATE` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CountNetPremi_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTNETPREMI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `CountOGPONP_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTOGPONP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `CountOverridingCommOgp_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTOVERRIDINGCOMMOGP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CountOverridingCommOnp_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTOVERRIDINGCOMMONP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CountPctInstallment_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTPCTINSTALLMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CountResult1Onp_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRESULT1ONP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CountResult1_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRESULT1_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `CountResult2Ogp_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRESULT2OGP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CountResult2Onp_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRESULT2ONP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CountRiCommOgp_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRICOMMOGP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CountRiCommOnp_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRICOMMONP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `CountSpreading_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTSPREADING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `FetchMasterTreatyIn.xml` | `ASM-FW-GISFW-WORK` | `FETCHMASTERTREATYIN` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `FetchTreatyGroupOJK.xml` | `ASM-FW-GISFW-WORK-NB` | `FETCHTREATYGROUPOJK` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `FetchTreatyGroupOldID.xml` | `ASM-FW-GISFW-WORK` | `FETCHTREATYGROUPOLDID` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `FillPaymentInstallment.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `FILLPAYMENTINSTALLMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `GeneratePolicyNoTreaty_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `GENERATEPOLICYNOTREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=32 | — |
| `GetCurrencyMaster.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `GETCURRENCYMASTER` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `GetDateValidity_ACT.xml` | `ASM-FW-GISFW-WORK` | `GETDATEVALIDITY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetTreatyName.xml` | `ASM-FW-GISFW-WORK` | `GETTREATYNAME` | [terverifikasi] pyActivityType=ACTIVITY; steps=27 | — |
| `InputParamUploadReas_act.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `INPUTPARAMUPLOADREAS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `InputPolicyTreatyEDMDetail_NP.xml` | `ASM-FW-GISFW-WORK` | `INPUTPOLICYTREATYEDMDETAIL_NP` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | `ASM-FW-GISFW-WORK / INPUTPOLICYTREATYINDETAIL_NONPROP` |
| `InputPolicyTreatyInDetail_NonProp.xml` | `ASM-FW-GISFW-WORK` | `INPUTPOLICYTREATYINDETAIL_NONPROP` | [terverifikasi] pyActivityType=ACTIVITY; steps=31 | — |
| `InputPolicyTreatyInDetail_preACT.xml` | `ASM-FW-GISFW-WORK` | `INPUTPOLICYTREATYINDETAIL_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=42 | `ASM-FW-GISFW-WORK / INPUTPOLICYTREATYIN_PREACT` |
| `InputPolicyTreatyInPost_Act.xml` | `ASM-FW-GISFW-WORK` | `INPUTPOLICYTREATYINPOST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `InputPolicyTreatyInPre_Act.xml` | `ASM-FW-GISFW-WORK` | `INPUTPOLICYTREATYINPRE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `InputPolicyTreatyOutDetail_NonProp.xml` | `ASM-FW-GISFW-WORK` | `INPUTPOLICYTREATYOUTDETAIL_NONPROP` | [terverifikasi] pyActivityType=ACTIVITY; steps=49 | `ASM-FW-GISFW-WORK / INPUTPOLICYTREATYINDETAIL_NONPROP` |
| `InputPolicyTreatyOutDetail_preACT.xml` | `ASM-FW-GISFW-WORK` | `INPUTPOLICYTREATYOUTDETAIL_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=32 | `ASM-FW-GISFW-WORK / INPUTPOLICYTREATYINDETAIL_PREACT` |
| `InputQuotation_PreAct.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `INPUTQUOTATION_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `InsertHistoryAkseptasiPega.xml` | `ASM-FW-GISFW-WORK` | `INSERTHISTORYAKSEPTASIPEGA` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `InsertToTreatyOutXOLList.xml` | `ASM-FW-GISFW-WORK` | `INSERTTOTREATYOUTXOLLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=37 | `ASM-FW-GISFW-WORK / INSERTTOTREATYXOLLIST` |
| `InsertToTreatyXOLList.xml` | `ASM-FW-GISFW-WORK` | `INSERTTOTREATYXOLLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=20 | — |
| `InsertToTreatyXOLListRetroShare.xml` | `ASM-FW-GISFW-WORK` | `INSERTTOTREATYXOLLISTRETROSHARE` | [terverifikasi] pyActivityType=ACTIVITY; steps=27 | `ASM-FW-GISFW-WORK / INSERTTOTREATYXOLLIST` |
| `ProtectCedingCo.xml` | `ASM-FW-GISFW-WORK` | `PROTECTCEDINGCO` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `ProtectCoverage_Act.xml` | `ASM-FW-GISFW-WORK` | `PROTECTCOVERAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=38 | — |
| `ProtectCurrencyTSI_Act.xml` | `ASM-FW-GISFW-WORK` | `PROTECTCURRENCYTSI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `ProtectDate.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `PROTECTDATE` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `ProtectFIREMBUPA_Act.xml` | `ASM-FW-GISFW-WORK` | `PROTECTFIREMBUPA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=89 | `ASM-FW-GISFW-WORK / PROTECTCOVERAGE_ACT` |
| `ProtectPremiPolicy_Act.xml` | `ASM-FW-GISFW-WORK` | `PROTECTPREMIPOLICY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `ProtectRenewal_Act.xml` | `ASM-FW-GISFW-WORK` | `PROTECTRENEWAL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `ProtectShareCedant_Act.xml` | `ASM-FW-GISFW-WORK` | `PROTECTSHARECEDANT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | — |
| `ProtectShipData_Act.xml` | `ASM-FW-GISFW-WORK` | `PROTECTSHIPDATA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `Protection_Act.xml` | `ASM-FW-GISFW-WORK` | `PROTECTION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=37 | — |
| `RemoveTypeTax_ACT.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `REMOVETYPETAX_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SaveJsonPolisTreatyIn_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SAVEJSONPOLISTREATYIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `SaveViewSuggest.xml` | `ASM-FW-GISFW-WORK` | `SAVEVIEWSUGGEST` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `SearchHierarkiSourceBizAgentTreatyIn_Act.xml` | `ASM-FW-GISFW-DATA-AGENT` | `SEARCHHIERARKISOURCEBIZAGENTTREATYIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GISFW-DATA-AGENT / SEARCHHIERARKISOURCEBIZAGENT_ACT` |
| `SetCategoryAttach.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SETCATEGORYATTACH` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SetCurrency_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SETCURRENCY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetDueTo_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SETDUETO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GISFW-DATA-POLICYTREATYIN / COUNTNETPREMI_ACT` |
| `SetFlagOccupation_ACT.xml` | `ASM-FW-GISFW-WORK` | `SETFLAGOCCUPATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | `ASM-FW-GISFW-WORK / PROTECTFIREMBUPA_ACT` |
| `SetPPNPPH.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SETPPNPPH` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SetReinstatementPct.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `SETREINSTATEMENTPCT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SetSurveyReport_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SETSURVEYREPORT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetTreatyCurrencyID.xml` | `ASM-FW-GISFW-WORK` | `SETTREATYCURRENCYID` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetTreatyIn_Act.xml` | `DATA-PORTAL` | `SETTREATYIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | — |
| `SetValidateInstallment_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SETVALIDATEINSTALLMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SetValueRetro_Act.xml` | `ASM-FW-GISFW-INT-TREATYOUTDETAIL` | `SETVALUERETRO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GISFW-INT-TREATYINDETAILJOINEDM / SETVALUE_ACT` |
| `SetValue_Act.xml` | `ASM-FW-GISFW-INT-TREATYINDETAILJOINEDM` | `SETVALUE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GISFW-INT-TREATYINDETAIL / SETVALUE_ACT` |
| `SumTSIPremiSpreadedRNM_ANEKA_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SUMTSIPREMISPREADEDRNM_ANEKA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=65 | `ASM-FW-GISFW-DATA-COVERAGE / SUMTSIPREMISPREADEDRNM_ACT` |
| `SumTSIPremiSpreadedRNM_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SUMTSIPREMISPREADEDRNM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=52 | — |
| `SumTSIPremiSpreadedRNM_FIRE_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SUMTSIPREMISPREADEDRNM_FIRE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=103 | `ASM-FW-GISFW-DATA-COVERAGE / SUMTSIPREMISPREADEDRNM_ACT` |
| `SumTSIPremiSpreadedRNM_GOLF_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SUMTSIPREMISPREADEDRNM_GOLF_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=24 | `ASM-FW-GISFW-DATA-COVERAGE / SUMTSIPREMISPREADEDRNM_ACT` |
| `SumTSIPremiSpreadedRNM_MARINECARGO_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SUMTSIPREMISPREADEDRNM_MARINECARGO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `SumTSIPremiSpreadedRNM_MBU_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SUMTSIPREMISPREADEDRNM_MBU_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=24 | `ASM-FW-GISFW-DATA-COVERAGE / SUMTSIPREMISPREADEDRNM_ACT` |
| `SumTSIPremiSpreadedRNM_PA_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SUMTSIPREMISPREADEDRNM_PA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `ASM-FW-GISFW-DATA-COVERAGE / SUMTSIPREMISPREADEDRNM_ACT` |
| `SumTSIPremiSpreadedRNM_TRAVEL_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SUMTSIPREMISPREADEDRNM_TRAVEL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `ASM-FW-GISFW-DATA-COVERAGE / SUMTSIPREMISPREADEDRNM_ACT` |
| `TreatyInInputVis.xml` | `DATA-PORTAL` | `TREATYININPUTVIS` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `TreatyInNonSetTotal.xml` | `DATA-PORTAL` | `TREATYINNONSETTOTAL` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `TreatyInputPctCommSpreading.xml` | `ASM-FW-GISFW-WORK` | `TREATYINPUTPCTCOMMSPREADING` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `TreatyNonPropOutSetSpreading.xml` | `ASM-FW-GISFW-WORK` | `TREATYNONPROPOUTSETSPREADING` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | `ASM-FW-GISFW-WORK / TREATYNONPROPSETSPREADING` |
| `TreatyNonPropSetSpreading.xml` | `ASM-FW-GISFW-WORK` | `TREATYNONPROPSETSPREADING` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | — |
| `TreatyRealizationCheckDuplicate.xml` | `ASM-FW-GISFW-WORK` | `TREATYREALIZATIONCHECKDUPLICATE` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `TreatyRealizationCheckXOLList.xml` | `ASM-FW-GISFW-WORK` | `TREATYREALIZATIONCHECKXOLLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `ASM-SFAGIS-WORK / TREATYREALIZATIONCHECKXOLLIST` |
| `TreatySetReinstatement.xml` | `DATA-PORTAL` | `TREATYSETREINSTATEMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `cekSpreadingFactIn.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `CEKSPREADINGFACTIN` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `serviceInsertArasapas_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SERVICEINSERTARASAPAS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |

### DataTransform — 12 rule (`RULE-OBJ-MODEL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AddToListCommentsPolicyTreatyIn_DT.xml` | `ASM-FW-GISFW-WORK` | `ADDTOLISTCOMMENTSPOLICYTREATYIN_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DeptHeadTreatyInUW_preDT.xml` | `ASM-FW-GISFW-WORK` | `DEPTHEADTREATYINUW_PREDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / INPUTPOLICYTREATYIN_PREDT` |
| `DeptHeadTreatyIn_UW_postDT.xml` | `ASM-FW-GISFW-WORK` | `DEPTHEADTREATYIN_UW_POSTDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ADDTOLISTCOMMENTSPOLICYTREATYIN_DT` |
| `InboxPolicyTreatyIn_postDT.xml` | `ASM-FW-GISFW-WORK` | `INBOXPOLICYTREATYIN_POSTDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputPolicyTreatyIn_preDT.xml` | `ASM-FW-GISFW-WORK` | `INPUTPOLICYTREATYIN_PREDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SearchHierarkiSourceBizAgent_PostDT.xml` | `ASM-FW-GISFW-DATA-AGENT` | `SEARCHHIERARKISOURCEBIZAGENT_POSTDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SystemSetOneYear_DT.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SYSTEMSETONEYEAR_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TestTreatyToFacStatus.xml` | `ASM-FW-GISFW-WORK` | `TESTTREATYTOFACSTATUS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyEnableDisableInput.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `TREATYENABLEDISABLEINPUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `btnCedingCO_DT.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `BTNCEDINGCO_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `btnSOB_DT.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `BTNSOB_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `setCategoryAttachment_DT.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SETCATEGORYATTACHMENT_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### DecisionTable — 2 rule (`RULE-DECLARE-DECISIONTABLE`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BusinessType_DeT.xml` | `ASM-FW-GISFW-WORK` | `BUSINESSTYPE_DET` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isApproved.xml` | `ASM-FW-GISFW-WORK` | `ISAPPROVED` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Flow — 1 rule (`RULE-OBJ-FLOW`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `InputRealizationTreatyIn.xml` | `ASM-FW-GISFW-WORK` | `INPUTREALIZATIONTREATYIN` | [terverifikasi] pyStartActivity=Start1 | — |

### FlowAction — 9 rule (`RULE-OBJ-FLOWACTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AgentSourceBizDetails.xml` | `ASM-FW-GISFW-DATA-AGENT` | `AGENTSOURCEBIZDETAILS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DeptHeadTreatyIn_UW.xml` | `ASM-FW-GISFW-WORK` | `DEPTHEADTREATYIN_UW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InboxPolicyTreatyIn.xml` | `ASM-FW-GISFW-WORK` | `INBOXPOLICYTREATYIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputHistoricalSurveyReport.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `INPUTHISTORICALSURVEYREPORT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputHistoricalSurveyReportUW.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `INPUTHISTORICALSURVEYREPORTUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-QUOTATION / INPUTHISTORICALSURVEYREPORT` |
| `InstallmentList.xml` | `ASM-FW-GISFW-DATA-INSTALLMENT` | `INSTALLMENTLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Installments_ReadOnly.xml` | `ASM-FW-GISFW-DATA-TREATYININSTALLMENT` | `INSTALLMENTS_READONLY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYININSTALLMENT / INSTALLMENTS_REALISASI` |
| `PolicyTreatyInDeclineConfirm.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `POLICYTREATYINDECLINECONFIRM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowPolicyNoTreaty.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SHOWPOLICYNOTREATY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Harness — 6 rule (`RULE-HTML-HARNESS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BusinessAndSOBList.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `BUSINESSANDSOBLIST` | [terverifikasi] pyInclude=1 | — |
| `BusinessAndSOBListRetro.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `BUSINESSANDSOBLISTRETRO` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-DATA-POLICYTREATYIN / BUSINESSANDSOBLIST` |
| `HistoricalSurveyReport.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `HISTORICALSURVEYREPORT` | [terverifikasi] pyInclude=1 | — |
| `HistoricalSurveyReportUW.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `HISTORICALSURVEYREPORTUW` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-DATA-POLICYTREATYIN / HISTORICALSURVEYREPORT` |
| `SFAPortalOpportunities.xml` | `PEGACRM-PORTAL` | `SFAPORTALOPPORTUNITIES` | [terverifikasi] pyInclude=6 | — |
| `SOB.xml` | `ASM-FW-GISFW-DATA-OFFERTREATYIN` | `SOB` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-DATA-QUOTATION / SOB` |

### RDBList — 41 rule (`RULE-CONNECT-SQL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AttachmentLife.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `ATTACHMENTLIFE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CATEGORY_ATTACH_REAS | — |
| `BrowseClientEmail_SQL.xml` | `ASM-FW-GISFW-INT-CLIENT` | `BROWSECLIENTEMAIL_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CLIENTEMAIL | — |
| `BrowseTreatyIn.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `BROWSETREATYIN` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_TREATY_IN,POOLDATA.M_TREATY_IN_EDM | `ASM-FW-GISFW-INT-PRODUCT_LIFE / ASM!BROWSEUNDERWRITINGLIST` |
| `BrowseTreatyInDetailJoinEDM.xml` | `ASM-FW-GISFW-INT-TREATYINDETAILJOINEDM` | `BROWSETREATYINDETAILJOINEDM` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_TREATY_IN_DETAIL_EDM | `ASM-FW-GISFW-INT-TREATYINDETAIL / ASM!BROWSETREATYINDETAIL` |
| `BrowseTreatyInJoinEDM.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `BROWSETREATYINJOINEDM` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_TREATY_IN,POOLDATA.M_TREATY_IN_EDM | `ASM-FW-GISFW-DATA-POLICYTREATYIN / RNM!BROWSETREATYIN` |
| `BrowseTreatyOut.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `BROWSETREATYOUT` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_TREATY_OUT | `ASM-FW-GISFW-DATA-POLICYTREATYIN / RNM!BROWSETREATYINJOINEDM` |
| `BrowseTreatyOutDetail.xml` | `ASM-FW-GISFW-INT-TREATYOUTDETAIL` | `BROWSETREATYOUTDETAIL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_TREATY_OUT | `ASM-FW-GISFW-INT-TREATYINDETAILJOINEDM / RNM!BROWSETREATYINDETAILJOINEDM` |
| `CategoryAttach_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `CATEGORYATTACH_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CATEGORY_ATTACH_REAS | — |
| `CheckZipCode_SQL.xml` | `ASM-FW-GISFW-INT-RW` | `CHECKZIPCODE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=RW | — |
| `FetchTreatyGroupOLDID.xml` | `ASM-FW-GISFW-INT-TREATYGROUP` | `FETCHTREATYGROUPOLDID` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TREATYGROUP | — |
| `GETTanggalClosing_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETTANGGALCLOSING_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TANGGAL_CLOSING | — |
| `GenerateNoPolicy.xml` | `ASM-FW-GISFW-INT-POLISTREATYIN` | `GENERATENOPOLICY` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GetAllCurrency.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `GETALLCURRENCY` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CURRENCY | — |
| `GetBreakDownSpread_SQL.xml` | `ASM-FW-GISFW-WORK` | `GETBREAKDOWNSPREAD_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.PROPORTIONALARRG | — |
| `GetCategoryOccupation.xml` | `ASM-FW-GISFW-INT-OCCUPATION` | `GETCATEGORYOCCUPATION` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=OCCUPATION | `ASM-FW-GISFW-INT-OCCUPATION / ASM!GETOCCUPATIONNAME` |
| `GetClientID_SQL.xml` | `ASM-FW-GISFW-INT-CLIENT` | `GETCLIENTID_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CLIENT | — |
| `GetCountClaim.xml` | `ASSIGN-WORKLIST` | `GETCOUNTCLAIM` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=DATAPEGA.PC_ASM_FW_GCNMFW_WORK | — |
| `GetCurrency.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `GETCURRENCY` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CURRENCY | `ASM-FW-GISFW-INT-CURRENCY / ASM!UPDATEMASTERCURRENCY` |
| `GetCurrencyIDByName.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `GETCURRENCYIDBYNAME` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.CURRENCY | `ASM-FW-GISFW-INT-CURRENCY / ASM!GETDATACURRENCYBYNAME_SQL` |
| `GetCurrentDate.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `GETCURRENTDATE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GetDataAgentByNameNonLife_SQL.xml` | `ASM-FW-GISFW-INT-AGENT` | `GETDATAAGENTBYNAMENONLIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=AGENT | `ASM-FW-GISFW-INT-AGENT / ASM!GETDATAAGENTBYNAME_SQL` |
| `GetDataByID_SQL.xml` | `ASM-FW-GISFW-INT-REINSURANCETYPE` | `GETDATABYID_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_REINSURANCETYPE | — |
| `GetDataCurrencyByName_SQL.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `GETDATACURRENCYBYNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CURRENCY | — |
| `GetDataDoubleCase.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETDATADOUBLECASE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=FACINPRODUCTION | — |
| `GetFlagReject_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETFLAGREJECT_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=HISTORYAKSEPTASIPEGA | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETVIEWBANDING_SQL` |
| `GetKodeProdNonLife_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETKODEPRODNONLIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.KODE_PRODUKSI | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETSEQUENCENUMBER_SQL` |
| `GetKursLimitSpreading_SQL.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `GETKURSLIMITSPREADING_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYEXCHANGEYEARLY | — |
| `GetObjectItembyName_SQL.xml` | `ASM-FW-GISFW-INT-V_JN_OBJ_ITEM` | `GETOBJECTITEMBYNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=V_JN_OBJ_ITEM | `ASM-FW-GISFW-INT-V_JN_OBJ_ITEM / ASM!GETOBJECTITEMBYID_SQL` |
| `GetOldIDBusiness_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `GETOLDIDBUSINESS_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=BUSINESS | — |
| `GetSQLDate.xml` | `ASM-FW-GISFW-WORK` | `GETSQLDATE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GetSequenceNumber_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETSEQUENCENUMBER_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETTANGGALCLOSING_SQL` |
| `GetTgl_InputJsonPolis_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETTGL_INPUTJSONPOLIS_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | — |
| `GetTreatyName.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `GETTREATYNAME` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=REINSURANCETYPE | `ASM-FW-GISFW-INT-PROPORTIONALARRG / ASM!GETTREATYNAME_SQL` |
| `GetTreatyName_SQL.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `GETTREATYNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.REINSURANCETYPE,PROPORTIONALARRG,TREATYBUSINESS,TREATYCONTRACT | — |
| `GetTreatyName_SQL2.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `GETTREATYNAME_SQL2` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.REINSURANCETYPE,PROPORTIONALARRG,TREATYBUSINESS | `ASM-FW-GISFW-INT-PROPORTIONALARRG / ASM!GETTREATYNAME_SQL` |
| `InsertHistoryAkseptasiPega_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTHISTORYAKSEPTASIPEGA_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=HISTORYAKSEPTASIPEGA | — |
| `InsertViewSuggest_SQL.xml` | `ASM-FW-GISFW-WORK` | `INSERTVIEWSUGGEST_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.HISTORYAKSEPTASIPRODUCTION | — |
| `SavePolisTreatyIn_SQL.xml` | `ASM-FW-GISFW-INT-POLISTREATYIN` | `SAVEPOLISTREATYIN_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_POLIS_TREATYIN | `ASM-FW-GISFW-INT-JSON_POLIS_TREATYIN / ASM!SAVEPOLISTREATYIN_SQL` |
| `SaveTreatyIn.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `SAVETREATYIN` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_TREATY_IN | — |
| `SelectSpreadingTreatyInProduction.xml` | `ASM-FW-GISFW-INT-REINSURANCETYPE` | `SELECTSPREADINGTREATYINPRODUCTION` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=REINSURANCETYPE | — |
| `TreatyRealizationCheckDuplicate.xml` | `ASM-FW-GISFW-WORK` | `TREATYREALIZATIONCHECKDUPLICATE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TREATYINPRODUCTION | — |

### ReportDefinition — 15 rule (`RULE-OBJ-REPORT-DEFINITION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseAgentHierarkiList_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSEAGENTHIERARKILIST_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseAgentNusaRe_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSEAGENTNUSARE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCedingCo_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSECEDINGCO_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-AGENT / BROWSEAGENTHIERARKILIST_RD` |
| `BrowseClientName_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSECLIENTNAME_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-AGENT / BROWSESEARCHSOBCEDING_RD` |
| `BrowseCurrencyTreatyIn_RD.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `BROWSECURRENCYTREATYIN_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCurrency_RD.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `BROWSECURRENCY_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseMarketingOfficer_RD.xml` | `ASM-FW-GISFW-INT-MARKETINGOFFICER` | `BROWSEMARKETINGOFFICER_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseReinsuranceType_RD.xml` | `ASM-FW-GISFW-INT-REINSURANCETYPE` | `BROWSEREINSURANCETYPE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTREATY_IN.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `BROWSETREATY_IN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyGroup_RD.xml` | `ASM-FW-GISFW-INT-TREATYGROUP` | `BROWSETREATYGROUP_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyInDetail.xml` | `ASM-FW-GISFW-INT-TREATYINDETAIL` | `BROWSETREATYINDETAIL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyJoinEDM.xml` | `ASM-FW-GISFW-INT-TREATYINDETAILJOINEDM` | `BROWSETREATYJOINEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-TREATYINDETAIL / BROWSETREATYINDETAIL` |
| `BrowseTreatyOutDetail.xml` | `ASM-FW-GISFW-INT-TREATYOUTDETAIL` | `BROWSETREATYOUTDETAIL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-TREATYINDETAIL / BROWSETREATYINDETAIL` |
| `GetListOpportunity.xml` | `ASM-FW-SFAGISFW-WORK-OPPORTUNITY` | `GETLISTOPPORTUNITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `crmOpportunitiesList.xml` | `ASM-FW-SFAGISFW-WORK-OPPORTUNITY` | `CRMOPPORTUNITIESLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Section — 25 rule (`RULE-HTML-SECTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BusinessAndSOBList.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `BUSINESSANDSOBLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BusinessAndSOBListRetro.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `BUSINESSANDSOBLISTRETRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-POLICYTREATYIN / BUSINESSANDSOBLIST` |
| `DetailDeptHeadTreatyIn_UW.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILDEPTHEADTREATYIN_UW` | [terverifikasi] pyInclude=5 | — |
| `DetailPoliciesNonProportional.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILPOLICIESNONPROPORTIONAL` | [terverifikasi] pyInclude=4 | — |
| `DetailPolicyTreatyIn.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILPOLICYTREATYIN` | [terverifikasi] pyInclude=6 | — |
| `DetailPolicyTreatyInNonProportional.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILPOLICYTREATYINNONPROPORTIONAL` | [terverifikasi] pyInclude=2 | — |
| `DetailPolicyTreatyInNonProportionalEDM.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILPOLICYTREATYINNONPROPORTIONALEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-POLICYTREATYIN / DETAILPOLICYTREATYINNONPROPORTIONAL` |
| `DetailPolicyTreatyOutNonProportional.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILPOLICYTREATYOUTNONPROPORTIONAL` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-POLICYTREATYIN / DETAILPOLICYTREATYINNONPROPORTIONAL` |
| `GeneralDeptHeadTreatyIn_UW.xml` | `ASM-FW-GISFW-WORK` | `GENERALDEPTHEADTREATYIN_UW` | [terverifikasi] pyInclude=6 | — |
| `GeneralPolicyTreatyIn.xml` | `ASM-FW-GISFW-WORK` | `GENERALPOLICYTREATYIN` | [terverifikasi] pyInclude=7 | — |
| `HistoricalSurveyReportDtl.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `HISTORICALSURVEYREPORTDTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN / HISTORICALSURVEYREPORTDTL` |
| `HistoricalSurveyReportDtlUW.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `HISTORICALSURVEYREPORTDTLUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-POLICYTREATYIN / HISTORICALSURVEYREPORTDTL` |
| `InputHistoricalSurveyReportDtl.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `INPUTHISTORICALSURVEYREPORTDTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputHistoricalSurveyReportDtlUW.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `INPUTHISTORICALSURVEYREPORTDTLUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-QUOTATION / INPUTHISTORICALSURVEYREPORTDTL` |
| `InstallmentList.xml` | `ASM-FW-GISFW-DATA-INSTALLMENT` | `INSTALLMENTLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Installments_ReadOnly.xml` | `ASM-FW-GISFW-DATA-TREATYININSTALLMENT` | `INSTALLMENTS_READONLY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYININSTALLMENT / INSTALLMENTS_REALISASI` |
| `ListSuggest.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `LISTSUGGEST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PolicyTreatyInDeclineConfirm.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `POLICYTREATYINDECLINECONFIRM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SFAPortalOpportunitiesHeader.xml` | `PEGACRM-PORTAL` | `SFAPORTALOPPORTUNITIESHEADER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SFAPortal_Opportunities.xml` | `DATA-PORTAL` | `SFAPORTAL_OPPORTUNITIES` | [terverifikasi] pyInclude=4 | — |
| `SFAPortal_OpportunitiesList.xml` | `DATA-PORTAL` | `SFAPORTAL_OPPORTUNITIESLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SFAPortal_OpportunitiesList_Header.xml` | `DATA-PORTAL` | `SFAPORTAL_OPPORTUNITIESLIST_HEADER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowPolicyNoTreaty_SC.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SHOWPOLICYNOTREATY_SC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SourceHierarki.xml` | `ASM-FW-GISFW-DATA-OFFERTREATYIN` | `SOURCEHIERARKI` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-QUOTATION / SOURCEHIERARKI` |
| `SpreadingRiskList.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SPREADINGRISKLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### When — 75 rule (`RULE-OBJ-WHEN`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `IsAneka.xml` | `ASM-FW-GISFW-DATA` | `ISANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsBillboardNeon.xml` | `ASM-FW-GISFW-WORK` | `ISBILLBOARDNEON` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISBILLBOARDNEON` |
| `IsBoiler.xml` | `ASM-FW-GISFW-DATA` | `ISBOILER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISCIT` |
| `IsBonding.xml` | `ASM-FW-GISFW-WORK` | `ISBONDING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsBondingAndCustomBonds.xml` | `ASM-FW-GISFW-DATA` | `ISBONDINGANDCUSTOMBONDS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISBONDING` |
| `IsBondingKBG.xml` | `ASM-FW-GISFW-WORK` | `ISBONDINGKBG` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsBuilderRisk.xml` | `ASM-FW-GISFW-DATA` | `ISBUILDERRISK` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCIS.xml` | `ASM-FW-GISFW-DATA` | `ISCIS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISCIT` |
| `IsCIT.xml` | `ASM-FW-GISFW-DATA` | `ISCIT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISBOLIER` |
| `IsCMI.xml` | `ASM-FW-GISFW-WORK` | `ISCMI` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCar.xml` | `ASM-FW-GISFW-WORK` | `ISCAR` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsClaim.xml` | `ASM-FW-GISFW-DATA` | `ISCLAIM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsContractorsPlantMachinery.xml` | `ASM-FW-GISFW-WORK` | `ISCONTRACTORSPLANTMACHINERY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCrime.xml` | `ASM-FW-GISFW-DATA` | `ISCRIME` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCustomBonds.xml` | `ASM-FW-GISFW-DATA` | `ISCUSTOMBONDS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsEDM.xml` | `ASM-FW-GISFW-WORK` | `ISEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsEDMRiSlip.xml` | `ASM-FW-GISFW-WORK` | `ISEDMRISLIP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsEar.xml` | `@BASECLASS` | `ISEAR` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsEdmAdjShareCedant.xml` | `ASM-FW-GISFW-WORK` | `ISEDMADJSHARECEDANT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISEDMADJREFNO` |
| `IsEdmAdjSpreading.xml` | `ASM-FW-GISFW-WORK` | `ISEDMADJSPREADING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISEDMADJSPREADING` |
| `IsEngineering.xml` | `ASM-FW-GISFW-WORK` | `ISENGINEERING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISENGINEERING` |
| `IsEnvironmental.xml` | `ASM-FW-GISFW-DATA` | `ISENVIRONMENTAL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISCRIME` |
| `IsErrorSpreading.xml` | `@BASECLASS` | `ISERRORSPREADING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsExclusion.xml` | `ASM-FW-GISFW-DATA` | `ISEXCLUSION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsFire.xml` | `ASM-FW-GISFW-WORK` | `ISFIRE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsFireStyle1.xml` | `ASM-FW-GISFW-WORK` | `ISFIRESTYLE1` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsFireStyle2.xml` | `ASM-FW-GISFW-WORK` | `ISFIRESTYLE2` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsGlass.xml` | `ASM-FW-GISFW-DATA` | `ISGLASS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsGolfInsurance.xml` | `ASM-FW-GISFW-DATA` | `ISGOLFINSURANCE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsGrowingTrees.xml` | `@BASECLASS` | `ISGROWINGTREES` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsHE.xml` | `ASM-FW-GISFW-DATA` | `ISHE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsKPR.xml` | `ASM-FW-GISFW-WORK` | `ISKPR` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsLandRig.xml` | `ASM-FW-GISFW-DATA` | `ISLANDRIG` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsLiability.xml` | `ASM-FW-GISFW-WORK` | `ISLIABILITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsLife.xml` | `ASM-FW-GISFW-DATA` | `ISLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / ISLIFE` |
| `IsMBD.xml` | `ASM-FW-GISFW-DATA` | `ISMBD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsMBU.xml` | `ASM-FW-GISFW-DATA` | `ISMBU` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsMBUCar.xml` | `ASM-FW-GISFW-DATA` | `ISMBUCAR` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsMaintenance.xml` | `ASM-FW-GISFW-DATA` | `ISMAINTENANCE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsMarineCargo.xml` | `ASM-FW-GISFW-WORK` | `ISMARINECARGO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsMarineHull.xml` | `ASM-FW-GISFW-DATA` | `ISMARINEHULL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsMarineHullOffshore.xml` | `ASM-FW-GISFW-DATA` | `ISMARINEHULLOFFSHORE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISWHETHERINDEX` |
| `IsNotAdmin.xml` | `@BASECLASS` | `ISNOTADMIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsObjectWithQuantityYear.xml` | `ASM-FW-GISFW-DATA` | `ISOBJECTWITHQUANTITYYEAR` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsOilGas.xml` | `ASM-FW-GISFW-WORK` | `ISOILGAS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsOperatorLife.xml` | `DATA-PORTAL` | `ISOPERATORLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / ISAUCTIONUSER` |
| `IsPA.xml` | `DATA-PARTY-PERSON` | `ISPA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / ISTRAVEL` |
| `IsPortRisk.xml` | `ASM-FW-GISFW-DATA` | `ISPORTRISK` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISBUILDERRISK` |
| `IsProductsLiability.xml` | `ASM-FW-GISFW-WORK` | `ISPRODUCTSLIABILITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISPRODUCTSLIABILITY` |
| `IsProfessionalLiability.xml` | `ASM-FW-GISFW-WORK` | `ISPROFESSIONALLIABILITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISPROFESSIONALLIABILITY` |
| `IsRenewal.xml` | `ASM-FW-GISFW-WORK` | `ISRENEWAL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISEDM` |
| `IsSPVCreate.xml` | `ASM-FW-GISFW-WORK` | `ISSPVCREATE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsSPVTreaty1.xml` | `ASM-FW-GISFW-WORK` | `ISSPVTREATY1` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISSPVCREATE` |
| `IsTBonding.xml` | `ASM-FW-GISFW-WORK` | `ISTBONDING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / IST1T3` |
| `IsTravel.xml` | `ASM-FW-GISFW-WORK` | `ISTRAVEL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsTreaty1.xml` | `ASM-FW-GISFW-WORK` | `ISTREATY1` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISSPVCREATE` |
| `IsUW.xml` | `ASM-FW-GISFW-WORK` | `ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISUW` |
| `IsWorkmenCompensation.xml` | `ASM-FW-GISFW-WORK` | `ISWORKMENCOMPENSATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISCIT` |
| `IsYieldShortfall.xml` | `ASM-FW-GISFW-DATA` | `ISYIELDSHORTFALL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISGROWINGTREES` |
| `NopolisEmpty.xml` | `ASM-FW-GISFW-WORK` | `NOPOLISEMPTY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ToTREATYDEPTHEAD.xml` | `ASM-FW-GISFW-WORK` | `TOTREATYDEPTHEAD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / TOMANAGERTEKNIK` |
| `TreatyMasterInEDM.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `TREATYMASTERINEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / TREATYMASTERINEDM` |
| `crmCreateOpportunity.xml` | `PEGACRM-WORK-` | `CRMCREATEOPPORTUNITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isAllRisk.xml` | `ASM-FW-GISFW-DATA` | `ISALLRISK` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isApproved.xml` | `ASM-FW-GISFW-WORK` | `ISAPPROVED` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isAviationHull.xml` | `ASM-FW-GISFW-DATA` | `ISAVIATIONHULL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isBillboardNeonSyariah.xml` | `ASM-FW-GISFW-DATA` | `ISBILLBOARDNEONSYARIAH` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isBurglary.xml` | `ASM-FW-GISFW-DATA` | `ISBURGLARY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isClaimTreaty.xml` | `ASM-FW-GISFW-WORK` | `ISCLAIMTREATY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isElectronicEquipment.xml` | `ASM-FW-GISFW-WORK` | `ISELECTRONICEQUIPMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isFidelity.xml` | `ASM-FW-GISFW-DATA` | `ISFIDELITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isSellingModeB2B.xml` | `@BASECLASS` | `ISSELLINGMODEB2B` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isSellingModeB2BB2C.xml` | `@BASECLASS` | `ISSELLINGMODEB2BB2C` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isSellingModeB2C.xml` | `@BASECLASS` | `ISSELLINGMODEB2C` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `pyIsIpadOrDesktop.xml` | `@BASECLASS` | `PYISIPADORDESKTOP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

## 4. Tabrakan nama di dalam modul ini

### 4.1 Nama file dipakai di lebih dari satu tipe rule

| Nama file | Tipe | Class |
| --- | --- | --- |
| `GetTreatyName.xml` | Activity, RDBList | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-INT-PROPORTIONALARRG` |
| `PolicyTreatyInDeclineConfirm.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-POLICYTREATYIN, ASM-FW-GISFW-DATA-POLICYTREATYIN` |
| `BusinessAndSOBListRetro.xml` | Harness, Section | `ASM-FW-GISFW-DATA-POLICYTREATYIN, ASM-FW-GISFW-DATA-POLICYTREATYIN` |
| `Installments_ReadOnly.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYININSTALLMENT, ASM-FW-GISFW-DATA-TREATYININSTALLMENT` |
| `TreatyRealizationCheckDuplicate.xml` | Activity, RDBList | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-WORK` |
| `InstallmentList.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-INSTALLMENT, ASM-FW-GISFW-DATA-INSTALLMENT` |
| `isApproved.xml` | DecisionTable, When | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-WORK` |
| `BusinessAndSOBList.xml` | Harness, Section | `ASM-FW-GISFW-DATA-POLICYTREATYIN, ASM-FW-GISFW-DATA-POLICYTREATYIN` |
| `BrowseTreatyOutDetail.xml` | RDBList, ReportDefinition | `ASM-FW-GISFW-INT-TREATYOUTDETAIL, ASM-FW-GISFW-INT-TREATYOUTDETAIL` |

Total: **9** nama file.

### 4.2 Rule `When` bernama sama di lebih dari satu class

_Tidak ada._

## 5. Keluarga klon (satu asal salinan dipakai >1 rule)

Dibaca dari `<pzOriginalInstanceKey>` yang sama pada beberapa file. Ini **sinyal kandidat
konsolidasi** saat migrasi, bukan bukti bahwa perilakunya identik — kesamaan perilaku baru
dapat dinyatakan setelah isinya dibandingkan di STEP D2.

| Asal salinan (class / nama) | Jumlah | File |
| --- | ---: | --- |
| `ASM-FW-GISFW-INT-CURRENCY / ASM!GETDATACURRENCYBYNAME_SQL` | 2 | `RDBList/GetCurrencyIDByName.xml`, `RDBList/GetDataCurrencyByName_SQL.xml` |
| `ASM-FW-GISFW-WORK / PROTECTCOVERAGE_ACT` | 2 | `Activity/ProtectCoverage_Act.xml`, `Activity/ProtectFIREMBUPA_Act.xml` |
| `ASM-FW-GISFW-WORK / INSERTTOTREATYXOLLIST` | 3 | `Activity/InsertToTreatyOutXOLList.xml`, `Activity/InsertToTreatyXOLList.xml`, `Activity/InsertToTreatyXOLListRetroShare.xml` |
| `ASM-FW-GISFW-WORK / INPUTPOLICYTREATYIN_PREDT` | 2 | `DataTransform/DeptHeadTreatyInUW_preDT.xml`, `DataTransform/InputPolicyTreatyIn_preDT.xml` |
| `ASM-FW-GISFW-WORK / ADDTOLISTCOMMENTSPOLICYTREATYIN_DT` | 2 | `DataTransform/AddToListCommentsPolicyTreatyIn_DT.xml`, `DataTransform/DeptHeadTreatyIn_UW_postDT.xml` |
| `ASM-FW-GISFW-DATA-INSTALLMENT / INSTALLMENTLIST` | 2 | `FlowAction/InstallmentList.xml`, `Section/InstallmentList.xml` |
| `ASM-FW-GISFW-WORK / CHECKRISLIP` | 2 | `Activity/CheckRISLIP.xml`, `Activity/CheckRISLIP_EDM.xml` |
| `ASM-FW-GISFW-DATA / ISCIT` | 2 | `When/IsBoiler.xml`, `When/IsCIS.xml` |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN / DETAILPOLICYTREATYINNONPROPORTIONAL` | 3 | `Section/DetailPolicyTreatyInNonProportional.xml`, `Section/DetailPolicyTreatyInNonProportionalEDM.xml`, `Section/DetailPolicyTreatyOutNonProportional.xml` |
| `ASM-FW-GISFW-DATA / ISBUILDERRISK` | 2 | `When/IsBuilderRisk.xml`, `When/IsPortRisk.xml` |
| `ASM-FW-GISFW-WORK / ISEDM` | 2 | `When/IsEDM.xml`, `When/IsRenewal.xml` |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN / COUNTNETPREMI_ACT` | 2 | `Activity/CountNetPremi_act.xml`, `Activity/SetDueTo_act.xml` |
| `ASM-FW-GISFW-DATA-QUOTATION / INPUTHISTORICALSURVEYREPORT` | 2 | `FlowAction/InputHistoricalSurveyReport.xml`, `FlowAction/InputHistoricalSurveyReportUW.xml` |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN / HISTORICALSURVEYREPORT` | 2 | `Harness/HistoricalSurveyReport.xml`, `Harness/HistoricalSurveyReportUW.xml` |
| `ASM-FW-GISFW-INT-TREATYINDETAIL / BROWSETREATYINDETAIL` | 3 | `ReportDefinition/BrowseTreatyInDetail.xml`, `ReportDefinition/BrowseTreatyJoinEDM.xml`, `ReportDefinition/BrowseTreatyOutDetail.xml` |
| `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETTANGGALCLOSING_SQL` | 2 | `RDBList/GETTanggalClosing_SQL.xml`, `RDBList/GetSequenceNumber_SQL.xml` |
| `ASM-FW-GISFW-WORK / ISSPVCREATE` | 3 | `When/IsSPVCreate.xml`, `When/IsSPVTreaty1.xml`, `When/IsTreaty1.xml` |
| `ASM-FW-GISFW-WORK / ISAPPROVED` | 2 | `DecisionTable/isApproved.xml`, `When/isApproved.xml` |
| `ASM-FW-GISFW-WORK / INPUTPOLICYTREATYINDETAIL_NONPROP` | 3 | `Activity/InputPolicyTreatyEDMDetail_NP.xml`, `Activity/InputPolicyTreatyInDetail_NonProp.xml`, `Activity/InputPolicyTreatyOutDetail_NonProp.xml` |
| `ASM-FW-GISFW-DATA-TREATYININSTALLMENT / INSTALLMENTS_REALISASI` | 2 | `FlowAction/Installments_ReadOnly.xml`, `Section/Installments_ReadOnly.xml` |
| `ASM-FW-GISFW-INT-PROPORTIONALARRG / ASM!GETTREATYNAME_SQL` | 3 | `RDBList/GetTreatyName.xml`, `RDBList/GetTreatyName_SQL.xml`, `RDBList/GetTreatyName_SQL2.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / SUMTSIPREMISPREADEDRNM_ACT` | 7 | `Activity/SumTSIPremiSpreadedRNM_ANEKA_Act.xml`, `Activity/SumTSIPremiSpreadedRNM_Act.xml`, `Activity/SumTSIPremiSpreadedRNM_FIRE_Act.xml`, `Activity/SumTSIPremiSpreadedRNM_GOLF_Act.xml`, `Activity/SumTSIPremiSpreadedRNM_MBU_Act.xml`, `Activity/SumTSIPremiSpreadedRNM_PA_Act.xml`, `Activity/SumTSIPremiSpreadedRNM_TRAVEL_Act.xml` |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN / POLICYTREATYINDECLINECONFIRM` | 2 | `FlowAction/PolicyTreatyInDeclineConfirm.xml`, `Section/PolicyTreatyInDeclineConfirm.xml` |
| `ASM-FW-GISFW-DATA / ISCRIME` | 2 | `When/IsCrime.xml`, `When/IsEnvironmental.xml` |
| `ASM-FW-GISFW-INT-AGENT / BROWSEAGENTHIERARKILIST_RD` | 2 | `ReportDefinition/BrowseAgentHierarkiList_RD.xml`, `ReportDefinition/BrowseCedingCo_RD.xml` |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN / BUSINESSANDSOBLIST` | 4 | `Harness/BusinessAndSOBList.xml`, `Harness/BusinessAndSOBListRetro.xml`, `Section/BusinessAndSOBList.xml`, `Section/BusinessAndSOBListRetro.xml` |
| `ASM-FW-GISFW-WORK / CHECKSPREADINGPROTECT_ACT` | 5 | `Activity/CheckSpreadingProtectAnekaGolf_ACT.xml`, `Activity/CheckSpreadingProtectFire_ACT.xml`, `Activity/CheckSpreadingProtectMCargoMBU_ACT.xml`, `Activity/CheckSpreadingProtectPATravel_ACT.xml`, `Activity/CheckSpreadingProtect_ACT.xml` |
| `ASM-FW-GISFW-WORK / TREATYNONPROPSETSPREADING` | 2 | `Activity/TreatyNonPropOutSetSpreading.xml`, `Activity/TreatyNonPropSetSpreading.xml` |
| `ASM-FW-GISFW-DATA-QUOTATION / INPUTHISTORICALSURVEYREPORTDTL` | 2 | `Section/InputHistoricalSurveyReportDtl.xml`, `Section/InputHistoricalSurveyReportDtlUW.xml` |

## 6. Objek database yang disentuh

Diambil dari `<pyBrowseSQL>` pada rule `RDBList`. **Nama objek apa adanya dari SQL** —
kualifikasi skema tidak konsisten di korpus. **Tipe kolom, constraint, dan relasi tidak
diketahui** (OQ-001).

| Tabel/view (apa adanya di SQL) | Dirujuk oleh |
| --- | ---: |
| `CURRENCY` | 3 rule |
| `CATEGORY_ATTACH_REAS` | 2 rule |
| `HISTORYAKSEPTASIPEGA` | 2 rule |
| `POOLDATA.M_TREATY_IN` | 2 rule |
| `POOLDATA.M_TREATY_IN_EDM` | 2 rule |
| `POOLDATA.M_TREATY_OUT` | 2 rule |
| `POOLDATA.REINSURANCETYPE` | 2 rule |
| `PROPORTIONALARRG` | 2 rule |
| `REINSURANCETYPE` | 2 rule |
| `TREATYBUSINESS` | 2 rule |
| `AGENT` | 1 rule |
| `BUSINESS` | 1 rule |
| `CLIENT` | 1 rule |
| `CLIENTEMAIL` | 1 rule |
| `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` | 1 rule |
| `FACINPRODUCTION` | 1 rule |
| `JSON_POLIS` | 1 rule |
| `M_REINSURANCETYPE` | 1 rule |
| `OCCUPATION` | 1 rule |
| `POOLDATA.CURRENCY` | 1 rule |
| `POOLDATA.HISTORYAKSEPTASIPRODUCTION` | 1 rule |
| `POOLDATA.KODE_PRODUKSI` | 1 rule |
| `POOLDATA.M_TREATY_IN_DETAIL_EDM` | 1 rule |
| `POOLDATA.PROPORTIONALARRG` | 1 rule |
| `POOLDATA.TANGGAL_CLOSING` | 1 rule |
| `POOLDATA.TREATYGROUP` | 1 rule |
| `POOLDATA.TREATYINPRODUCTION` | 1 rule |
| `RW` | 1 rule |
| `TREATYCONTRACT` | 1 rule |
| `TREATYEXCHANGEYEARLY` | 1 rule |
| `V_JN_OBJ_ITEM` | 1 rule |

### Stored procedure / function yang dipanggil

**Body prosedur ini TIDAK ADA di korpus** — logika di dalamnya tidak dapat direkonstruksi
dari XML. Lihat OQ-002.

| Procedure | Dipanggil oleh |
| --- | --- |
| `POOLDATA.PEGA_JSON_POLIS_TREATYIN` | `SavePolisTreatyIn_SQL.xml` |
| `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` | `GetSequenceNumber_SQL.xml` |
| `POOLDATA.PEGA_TREATY_IN` | `SaveTreatyIn.xml` |

## 7. Integrasi eksternal (`ConnectREST`)

_Tidak ada rule ConnectREST di modul ini._

