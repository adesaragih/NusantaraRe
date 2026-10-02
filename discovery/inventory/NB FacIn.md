# Inventaris Rule — NB FacIn

STEP D1, batch 2 (domain facultative inward). Dibangkitkan 2026-09-12 dari korpus READ-ONLY
`D:\XML\RNM_BRD\NB FacIn\`.

**Batas step ini:** hanya struktur + metadata. Tidak ada kesimpulan logika bisnis di sini.
Penelusuran perilaku adalah STEP D2 (`../flows/`).

**Identitas rule** diambil dari `<pxInsName>` (lihat `../README.md` §6.4).  
Kolom **Asal salinan** berisi `<pzOriginalInstanceKey>` dan hanya diisi bila **berbeda**
dari identitas — artinya rule ini hasil *Save As* dari rule lain. Tanda `—` berarti sama.
Class `—` berarti `<pxInsName>` tidak berisi bagian class (berlaku pada `DataPage`, yang
beridentitas global, bukan per class).

Perintah audit (dijalankan dari `D:\XML\RNM_BRD\`):

```
find "NB FacIn" -type f -name "*.xml" | wc -l
find "NB FacIn" -type f -name "*.xml" | awk -F/ '{print $2}' | sort | uniq -c
grep -o "<pxInsName>[^<]*" "NB FacIn/<tipe>/<file>.xml" | head -1
grep -o "<pzOriginalInstanceKey>[^<]*" "NB FacIn/<tipe>/<file>.xml" | head -1
```

## 1. Ringkasan terukur

- Total file XML: **2083** — cocok dengan `../README.md` §5.1
- **Identitas rule unik: 2081** dari 2083 file → 2 file adalah identitas ganda **di dalam modul ini**
- Class Pega unik: **173**
- RDBList menurut jenis SQL: PLSQL=42, QUERY=175 (`PLSQL` = blok anonim `BEGIN…END;`, umumnya memanggil stored procedure)
- Prefix pada kunci RDBList: `ASM`=178, `RNM`=39 (arti prefix **belum terverifikasi** — OQ-008)

| Tipe rule | Jumlah |
| --- | ---: |
| Activity | 609 |
| ConnectREST | 3 |
| DataPage | 30 |
| DataTransform | 148 |
| DecisionTable | 12 |
| DecisionTree | 1 |
| Flow | 6 |
| FlowAction | 250 |
| Harness | 41 |
| RDBList | 217 |
| ReportDefinition | 123 |
| Section | 432 |
| SystemSettings | 1 |
| When | 210 |
| **TOTAL** | **2083** |

## 2. Sebaran class Pega

Class yang **tidak** berawalan `ASM-` adalah class bawaan Pega atau class portal — ruang
lingkup migrasinya belum ditentukan (**OQ-009**).

| Class | Jumlah rule | Aplikasi? |
| --- | ---: | --- |
| `ASM-FW-GISFW-WORK` | 551 | aplikasi |
| `ASM-FW-GISFW-DATA-COVERAGE` | 166 | aplikasi |
| `@BASECLASS` | 107 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-INT-POLICYJSON` | 85 | aplikasi |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | 83 | aplikasi |
| `ASM-FW-GISFW-DATA-OFFERFACIN` | 81 | aplikasi |
| `DATA-PARTY-PERSON` | 74 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-DATA-PROPERTYITEM` | 54 | aplikasi |
| `ASM-FW-GISFW-DATA-ANEKA` | 46 | aplikasi |
| `ASM-FW-GISFW-DATA-OCCUPATION` | 45 | aplikasi |
| `ASM-FW-GISFW-DATA-FACOFFER` | 43 | aplikasi |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN` | 43 | aplikasi |
| `ASM-FW-GISFW-DATA` | 34 | aplikasi |
| `ASM-FW-GISFW-DATA-DEDUCTIBLE` | 34 | aplikasi |
| `DATA-PORTAL` | 33 | **bukan aplikasi** → OQ-009 |
| _(tanpa class)_ | 30 | n/a (DataPage) |
| `ASM-FW-GISFW-DATA-CARGO` | 30 | aplikasi |
| `ASM-FW-GISFW-DATA-VEHICLE` | 28 | aplikasi |
| `ASM-FW-GISFW-INT-OFFERJSON` | 26 | aplikasi |
| `ASM-FW-GISFW-DATA-CLAUSE` | 23 | aplikasi |
| `ASM-FW-GISFW-DATA-QUOTATION` | 21 | aplikasi |
| `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | 19 | aplikasi |
| `ASM-FW-GISFW-DATA-PROPERTY` | 18 | aplikasi |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | 14 | aplikasi |
| `ASM-FW-GISFW-INT-RW` | 12 | aplikasi |
| `CODE-PEGA-LIST` | 12 | **bukan aplikasi** → OQ-009 |
| `DATA-ADDRESS` | 12 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-DATA-CAUSEOFLOSS` | 10 | aplikasi |
| `ASM-FW-GISFW-DATA-POLICY` | 10 | aplikasi |
| `ASM-FW-GISFW-INT-CURRENCY` | 10 | aplikasi |
| `ASM-FW-GISFW-INT-ACCUMULATION` | 9 | aplikasi |
| `ASM-FW-GISFW-INT-AGENT` | 8 | aplikasi |
| `ASM-FW-GISFW-INT-OCCUPATION` | 8 | aplikasi |
| `ASM-FW-GISFW-DATA-OUTGO` | 7 | aplikasi |
| `ASM-FW-GISFW-DATA-SCORINGRISK` | 7 | aplikasi |
| `ASM-FW-GISFW-DATA-SUBCONTRACT` | 7 | aplikasi |
| `ASM-FW-GISFW-INT-MARKETINGOFFICER` | 7 | aplikasi |
| `ASM-FW-GISFW-INT-SHIP` | 7 | aplikasi |
| `ASM-FW-GISFW-WORK-NB` | 7 | aplikasi |
| `ASM-FW-GISFW-DATA-GOOD` | 6 | aplikasi |
| `ASM-FW-GISFW-INT-ACCUMULATION_LIFE` | 6 | aplikasi |
| `ASM-FW-GISFW-INT-RISKADDRESS` | 6 | aplikasi |
| `WORK-` | 6 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-DATA-CONVEYANCE` | 5 | aplikasi |
| `ASM-FW-GISFW-DATA-CORRESPONDENCE` | 5 | aplikasi |
| `ASM-FW-GISFW-DATA-TRADING` | 5 | aplikasi |
| `ASM-FW-GISFW-INT-BRANDDETAIL` | 5 | aplikasi |
| `ASM-FW-GISFW-INT-CLAUSE` | 5 | aplikasi |
| `ASM-FW-GISFW-INT-PROPORTIONALARRG` | 5 | aplikasi |
| `ASM-FW-GISFW-DATA-AGENT` | 4 | aplikasi |
| `ASM-FW-GISFW-DATA-PLAN` | 4 | aplikasi |
| `ASM-FW-GISFW-DATA-SPREADINGRISK` | 4 | aplikasi |
| `ASM-FW-GISFW-INT` | 4 | aplikasi |
| `ASM-FW-GISFW-INT-TREATY_IN` | 4 | aplikasi |
| `ASM-FW-GISFW-DATA-BATCHPALIST` | 3 | aplikasi |
| `ASM-FW-GISFW-DATA-INSTALLMENT` | 3 | aplikasi |
| `ASM-FW-GISFW-DATA-MEMORANDUMACCEPT` | 3 | aplikasi |
| `ASM-FW-GISFW-DATA-OFFERFACIN-OFFERFEALIST` | 3 | aplikasi |
| `ASM-FW-GISFW-DATA-PREMIUMTEMPLATE` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-BRANCH` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-BUSINESS` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-CLIENT` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-CONDITION` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-COVERAGE` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-DEDUCTIBLE` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-DISTRICT` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-DOCUMENT_POLIS` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-LLOYDAGENT` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-NATION` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-REINSURANCETYPE` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-TABLE_B2B` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-V_DUPLICATE_VEHICLE` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-V_JN_OBJ_ITEM` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-V_JOB_PA` | 3 | aplikasi |
| `ASM-FW-SFAGISFW-WORK-OPPORTUNITY` | 3 | aplikasi |
| `DATA-WORKATTACH-FILE` | 3 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-DATA-ACCESSORY` | 2 | aplikasi |
| `ASM-FW-GISFW-DATA-CAUSEOFDECLINE` | 2 | aplikasi |
| `ASM-FW-GISFW-DATA-COINS` | 2 | aplikasi |
| `ASM-FW-GISFW-DATA-ENUMERATION` | 2 | aplikasi |
| `ASM-FW-GISFW-DATA-FEA` | 2 | aplikasi |
| `ASM-FW-GISFW-DATA-OFFERTREATYIN` | 2 | aplikasi |
| `ASM-FW-GISFW-DATA-SPLIT` | 2 | aplikasi |
| `ASM-FW-GISFW-DATA-TREATYININSTALLMENT` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-ACCUMULATEDTYPE` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-BRAND` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-CITY` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-COVERAGE_FACIN` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-M_KONSTRUKSI` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-M_RATE_LIFE` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-OBJECTITEMTYPE` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-POLISTREATYIN` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-PREMITRAVEL` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-PROVINCE` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-V_COINS` | 2 | aplikasi |
| `ASM-FW-GISFW-WORK-TKUSALES` | 2 | aplikasi |
| `CODE-PEGA-PDF` | 2 | **bukan aplikasi** → OQ-009 |
| `DATA-ADMIN-OPERATOR-ID` | 2 | **bukan aplikasi** → OQ-009 |
| `LINK-ATTACHMENT` | 2 | **bukan aplikasi** → OQ-009 |
| `PEGACRM-PORTAL` | 2 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-DATA-ANEKAOTHERSCHEDULE` | 1 | aplikasi |
| `ASM-FW-GISFW-DATA-ANEKAPARTICIPANT` | 1 | aplikasi |
| `ASM-FW-GISFW-DATA-BENEFIT` | 1 | aplikasi |
| `ASM-FW-GISFW-DATA-FACOFFEROBJECTLIST` | 1 | aplikasi |
| `ASM-FW-GISFW-DATA-PRINTRISLIP` | 1 | aplikasi |
| `ASM-FW-GISFW-DATA-SEARCH` | 1 | aplikasi |
| `ASM-FW-GISFW-DATA-TREATYINLIMITS` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-ACCESORY` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-ACCESORYBRAND` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-ACCESORYTYPE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-ARGCLAUSEFIRE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-BRANCHDETAIL` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-BUSINESSFIELD` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-COLOR` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-CONSTRUCTION` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-CONVEYANCE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-CURRENCYSTANDARD` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-CZONE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-GOODSTYPE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-INWARDSCALE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-JOB` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-JSONOFFER` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-LST_BANK_GROUP` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-LST_BUSINESS_DUE_DATE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-MARINECONDITION` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-MAXAGE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-MCLIENT` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-MCONVEYANCE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-MV_AGEN` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_KLAUSUL_PLAN` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_LINK_SERVICE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_OKUPASI_ANEKA` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_TYPE_PROPERTY_PLAN` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-OBJECT` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-OC` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-OCCUPATIONS` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-PACKING` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-PLANTRAVEL` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-RATE_LIFE_SUMMARY` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-RETROCESSIONLIFE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-RISK` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-SPREADSYARIAH` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-SUBCONTRACT` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-SUBJECTTO` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-TABLEOFLIMIT` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-TRADING` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYBUSINESS` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYGROUP` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-TYPESHIP` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-VIEW_BENEFIT_PROPERTY` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-VIEW_PREMI_TRAVEL` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-VJ_BENEFIT_PROPERTY` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-VJ_M_TYPE_PROPERTY_PLAN` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-VJ_PACKAGE_AGE_KLAUSUL_DM` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-VJ_PKG_BENEFIT` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-V_ALM_RISIKO` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-V_CLAUSES_ARG` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-V_COVERAGE_ANEKA` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-V_DRAFTWORDING` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-V_MARKETING_LEADER` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-V_PACKAGE_AGE_KLAUSUL` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-V_PKG_BENEFIT` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-V_USER_ASSIGNMENT` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-WARRANTY` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-ZONES` | 1 | aplikasi |
| `ASM-FW-GISFW-WORK-ENDORSEMENT` | 1 | aplikasi |
| `ASM-FW-GISFW-WORK-LIFE` | 1 | aplikasi |
| `ASSIGN-WORKLIST` | 1 | **bukan aplikasi** → OQ-009 |
| `LINKSERVICE` | 1 | **bukan aplikasi** → OQ-009 |
| `PEGA-API-CASEMANAGEMENT-CASE` | 1 | **bukan aplikasi** → OQ-009 |
| `PEGACRM-WORK-` | 1 | **bukan aplikasi** → OQ-009 |
| `RULE-FILE-BINARY` | 1 | **bukan aplikasi** → OQ-009 |

Rule di class bukan-aplikasi: **260** dari 2083.

## 3. Daftar rule per tipe

### Activity — 609 rule (`RULE-OBJ-ACTIVITY`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `ASMForceCaseClose.xml` | `WORK-` | `ASMFORCECASECLOSE` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `ASMLogActivity_Act.xml` | `@BASECLASS` | `ASMLOGACTIVITY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `AddCedantList_act.xml` | `ASM-FW-GISFW-WORK` | `ADDCEDANTLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `AddCedingList_act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `ADDCEDINGLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `AddCoverageAutoFire.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `ADDCOVERAGEAUTOFIRE` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | — |
| `AddCurencyList_ACT.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `ADDCURENCYLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=37 | `ASM-FW-GISFW-DATA-PROPERTYITEM / ADDPAYMENTCURENCY_ACT` |
| `AddCurrencyListPA_ACT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `ADDCURRENCYLISTPA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=20 | — |
| `AddCurrencyList__ACT.xml` | `DATA-PARTY-PERSON` | `ADDCURRENCYLIST__ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `ASM-FW-GISFW-DATA-CARGO / ADDCURRENCYLIST__ACT` |
| `AddLocation_act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `ADDLOCATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `AddPaymentCurency_ACT.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `ADDPAYMENTCURENCY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=24 | `ASM-FW-GISFW-WORK / ADDPAYMENTCURENCY_ACT` |
| `AddStatusCoverageNew_Act.xml` | `ASM-FW-GISFW-WORK` | `ADDSTATUSCOVERAGENEW_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `AdjustmentRiskAccumulation_Act.xml` | `DATA-PORTAL` | `ADJUSTMENTRISKACCUMULATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `AgentSourceBiz_Act.xml` | `ASM-FW-GISFW-DATA-AGENT` | `AGENTSOURCEBIZ_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `AppendCoverageLife_ACT.xml` | `DATA-PARTY-PERSON` | `APPENDCOVERAGELIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `AppendCoverageMarineCargo_ACT.xml` | `ASM-FW-GISFW-DATA-CARGO` | `APPENDCOVERAGEMARINECARGO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `AttachAsPDFC.xml` | `@BASECLASS` | `ATTACHASPDFC` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `AttachAsPDFCFOLife.xml` | `ASM-FW-GISFW-WORK` | `ATTACHASPDFCFOLIFE` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `ASM-FW-GISFW-DATA-FACOFFER / ATTACHASPDFCFOLIFE` |
| `AttachFileGIS.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `ATTACHFILEGIS` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `ASM-FW-GISFW-WORK-NB / ATTACHSERVISH2H_ACT` |
| `AttachRISlipToWork.xml` | `CODE-PEGA-PDF` | `ATTACHRISLIPTOWORK` | [terverifikasi] pyActivityType=ACTIVITY; steps=26 | `CODE-PEGA-PDF / ATTACHTOWORK` |
| `AttachToWork.xml` | `CODE-PEGA-PDF` | `ATTACHTOWORK` | [terverifikasi] pyActivityType=ACTIVITY; steps=26 | — |
| `AturNilaiIndex.xml` | `ASM-FW-GISFW-WORK` | `ATURNILAIINDEX` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `BreakDownSpreading_Act.xml` | `ASM-FW-GISFW-WORK` | `BREAKDOWNSPREADING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `BrowseActV_COINS.xml` | `ASM-FW-GISFW-INT-V_COINS` | `BROWSEACTV_COINS` | [terverifikasi] pyActivityType=RULECONNECT; steps=9 | `ASM-FW-GISFW-INT-V_COINS / BROWSEACTV_OKUPASI` |
| `BrowseActV_VEHICLE_MODEL.xml` | `ASM-FW-GISFW-INT-BRANDDETAIL` | `BROWSEACTV_VEHICLE_MODEL` | [terverifikasi] pyActivityType=RULECONNECT; steps=5 | — |
| `BrowseCoverageSupport.xml` | `ASM-FW-GISFW-WORK` | `BROWSECOVERAGESUPPORT` | [terverifikasi] pyActivityType=RULECONNECT; steps=3 | `ASM-FW-GISFW-WORK / BROWSEM_OKUPASI_ANEKA` |
| `BrowseLocationByZipCodeGolf_ACT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `BROWSELOCATIONBYZIPCODEGOLF_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `BrowseM_OKUPASI_ANEKA.xml` | `ASM-FW-GISFW-WORK` | `BROWSEM_OKUPASI_ANEKA` | [terverifikasi] pyActivityType=RULECONNECT; steps=3 | — |
| `BrowsePolicyJSONMarine_Act.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `BROWSEPOLICYJSONMARINE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `ASM-FW-GISFW-INT-POLICYJSON / BROWSEPOLICYJSON_ACT` |
| `BrowseV_ZIPCODE.xml` | `ASM-FW-GISFW-WORK` | `BROWSEV_ZIPCODE` | [terverifikasi] pyActivityType=RULECONNECT; steps=3 | — |
| `BrowseZipCodeFireFacIn_Act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `BROWSEZIPCODEFIREFACIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `BrowseZipCodeFire_Act.xml` | `ASM-FW-GISFW-DATA-PROPERTY` | `BROWSEZIPCODEFIRE_ACT` | [terverifikasi] pyActivityType=RULECONNECT; steps=4 | — |
| `BrowseZipCode_Act.xml` | `ASM-FW-GISFW-WORK-TKUSALES` | `BROWSEZIPCODE_ACT` | [terverifikasi] pyActivityType=RULECONNECT; steps=4 | — |
| `BuildingYear_Act.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `BUILDINGYEAR_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `CalculateNetRate_ACT.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `CALCULATENETRATE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GISFW-DATA-COVERAGE / CALCULATENETRATE_ACT` |
| `CalculatePremiPA_FacIn.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `CALCULATEPREMIPA_FACIN` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `CalculatePremi_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `CALCULATEPREMI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `CalcultePersentageSpeadingLife_Act.xml` | `ASM-FW-GISFW-WORK` | `CALCULTEPERSENTAGESPEADINGLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `ASM-FW-GISFW-WORK / CALCULTEPERSENTAGESPEADING_ACT` |
| `CalcultePersentageSpeading_Act.xml` | `ASM-FW-GISFW-WORK` | `CALCULTEPERSENTAGESPEADING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=22 | — |
| `CalcultePersentageSpeading_Act_FacIn.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `CALCULTEPERSENTAGESPEADING_ACT_FACIN` | [terverifikasi] pyActivityType=ACTIVITY; steps=35 | `ASM-FW-GISFW-DATA-COVERAGE / CALCULTEPERSENTAGESPEADING_ACT` |
| `CallAturNilaiIndex.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `CALLATURNILAIINDEX` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CallCheckDtlObject.xml` | `DATA-PARTY-PERSON` | `CALLCHECKDTLOBJECT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CallCountPremiFacIn_act.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `CALLCOUNTPREMIFACIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `CallUploadCSVPolicyMember_PostAct.xml` | `ASM-FW-GISFW-DATA-POLICY` | `CALLUPLOADCSVPOLICYMEMBER_POSTACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `CallVirusCheck.xml` | `DATA-WORKATTACH-FILE` | `CALLVIRUSCHECK` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `CancelInputAccumulation_Act.xml` | `@BASECLASS` | `CANCELINPUTACCUMULATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CekCauseOfDecline_ACT.xml` | `ASM-FW-GISFW-WORK` | `CEKCAUSEOFDECLINE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `ASM-FW-GISFW-WORK / SETCAUSEOFDECLINE_ACT` |
| `CekDuplicateLocationOffer.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `CEKDUPLICATELOCATIONOFFER` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `CekLimitSpreading_Act.xml` | `ASM-FW-GISFW-WORK` | `CEKLIMITSPREADING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=29 | — |
| `CekLimitTreatyAcc_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `CEKLIMITTREATYACC_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CekNetRate_ACT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `CEKNETRATE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `CekTotalShareOffered_Act.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `CEKTOTALSHAREOFFERED_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | — |
| `ChangeCoveragePremium.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `CHANGECOVERAGEPREMIUM` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `ChangeEmailText.xml` | `ASM-FW-GISFW-WORK` | `CHANGEEMAILTEXT` | [terverifikasi] pyActivityType=ACTIVITY; steps=55 | — |
| `ChangeEmailTextCeding_ACT.xml` | `ASM-FW-GISFW-WORK` | `CHANGEEMAILTEXTCEDING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=55 | `ASM-FW-GISFW-WORK / CHANGEEMAILTEXT` |
| `ChangeListViewOfferFacOut.xml` | `ASM-FW-GISFW-WORK` | `CHANGELISTVIEWOFFERFACOUT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `ChangeProRatePreActVehicle_FacIn.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `CHANGEPRORATEPREACTVEHICLE_FACIN` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GISFW-DATA-ANEKA / CHANGEPRORATEPREACT_FACIN` |
| `ChangeSpreadingPercentageLife_Act.xml` | `ASM-FW-GISFW-WORK` | `CHANGESPREADINGPERCENTAGELIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `ASM-FW-GISFW-WORK / CHANGESPREADINGPERCENTAGE_ACT` |
| `ChangeSpreadingPercentage_Act.xml` | `ASM-FW-GISFW-WORK` | `CHANGESPREADINGPERCENTAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `ChangeUpKendaraan.xml` | `ASM-FW-GISFW-WORK` | `CHANGEUPKENDARAAN` | [terverifikasi] pyActivityType=ACTIVITY; steps=21 | — |
| `CheckDataFacOut_Act.xml` | `ASM-FW-GISFW-WORK` | `CHECKDATAFACOUT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CheckDataMarketing.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `CHECKDATAMARKETING` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `CheckDataMkt.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `CHECKDATAMKT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `ASM-FW-GISFW-DATA-POLICYTREATYIN / CHECKDATAMO` |
| `CheckDataQuotation.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `CHECKDATAQUOTATION` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | `ASM-FW-GISFW-DATA-QUOTATION / CHECKQUOTATION` |
| `CheckDeductible_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `CHECKDEDUCTIBLE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=23 | — |
| `CheckDtlObject.xml` | `ASM-FW-GISFW-WORK` | `CHECKDTLOBJECT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `CheckDtlObjectTravel.xml` | `ASM-FW-GISFW-WORK` | `CHECKDTLOBJECTTRAVEL` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `ASM-FW-GISFW-WORK / CHECKDTLOBJECT` |
| `CheckDuplicateOffer.xml` | `DATA-PORTAL` | `CHECKDUPLICATEOFFER` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `CheckLocation_FacIn.xml` | `ASM-FW-GISFW-WORK` | `CHECKLOCATION_FACIN` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `CheckObject.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `CHECKOBJECT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CheckOkupasi.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `CHECKOKUPASI` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `CheckPeriodLife_ACT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `CHECKPERIODLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `CheckPremiSpreading_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `CHECKPREMISPREADING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `CheckProtectFacout_Act.xml` | `ASM-FW-GISFW-WORK` | `CHECKPROTECTFACOUT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `CheckRISLIP.xml` | `ASM-FW-GISFW-WORK` | `CHECKRISLIP` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CheckRISLIP_EDM.xml` | `ASM-FW-GISFW-WORK` | `CHECKRISLIP_EDM` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `ASM-FW-GISFW-WORK / CHECKRISLIP` |
| `CheckSpreadingProtectAnekaGolf_ACT.xml` | `ASM-FW-GISFW-WORK` | `CHECKSPREADINGPROTECTANEKAGOLF_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=43 | `ASM-FW-GISFW-WORK / CHECKSPREADINGPROTECT_ACT` |
| `CheckSpreadingProtectFire_ACT.xml` | `ASM-FW-GISFW-WORK` | `CHECKSPREADINGPROTECTFIRE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=39 | `ASM-FW-GISFW-WORK / CHECKSPREADINGPROTECT_ACT` |
| `CheckSpreadingProtectMCargoMBU_ACT.xml` | `ASM-FW-GISFW-WORK` | `CHECKSPREADINGPROTECTMCARGOMBU_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=31 | `ASM-FW-GISFW-WORK / CHECKSPREADINGPROTECT_ACT` |
| `CheckSpreadingProtectPATravel_ACT.xml` | `ASM-FW-GISFW-WORK` | `CHECKSPREADINGPROTECTPATRAVEL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=31 | `ASM-FW-GISFW-WORK / CHECKSPREADINGPROTECT_ACT` |
| `CheckSpreadingProtect_ACT.xml` | `ASM-FW-GISFW-WORK` | `CHECKSPREADINGPROTECT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=31 | — |
| `ChooseActivityMappingCSV_Act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `CHOOSEACTIVITYMAPPINGCSV_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CleanActAneka.xml` | `ASM-FW-GISFW-WORK` | `CLEANACTANEKA` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `ASM-FW-GISFW-WORK / CLEANACTMBU` |
| `CleanActMBU.xml` | `ASM-FW-GISFW-WORK` | `CLEANACTMBU` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CleanActPA.xml` | `ASM-FW-GISFW-WORK` | `CLEANACTPA` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GISFW-WORK / CLEANACTMBU` |
| `ClearPageIdxCurr_ACT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | `CLEARPAGEIDXCURR_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `ClearPageRiskAddress.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `CLEARPAGERISKADDRESS` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `ClearRiskLocation_act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `CLEARRISKLOCATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CloseCheckViewFacOut.xml` | `ASM-FW-GISFW-WORK` | `CLOSECHECKVIEWFACOUT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CloseViewOfferFacOut.xml` | `ASM-FW-GISFW-WORK` | `CLOSEVIEWOFFERFACOUT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `ConcatSlipOfferNo_Act.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `CONCATSLIPOFFERNO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `ConvertHistoryDate.xml` | `DATA-PORTAL` | `CONVERTHISTORYDATE` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `ConvertPropViewPolisFacOut_Act.xml` | `ASM-FW-GISFW-WORK` | `CONVERTPROPVIEWPOLISFACOUT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | — |
| `ConvertResultsToTXT.xml` | `@BASECLASS` | `CONVERTRESULTSTOTXT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `@BASECLASS / PXCONVERTRESULTSTOCSV` |
| `CopyAccumulationCode_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `COPYACCUMULATIONCODE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `CopyAllObjFacOutFireAneka_ACT.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `COPYALLOBJFACOUTFIREANEKA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=23 | `ASM-FW-GISFW-DATA-FACOFFER / COPYALLOBJ_ACT` |
| `CopyAllObjFacOutGolfCargo_ACT.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `COPYALLOBJFACOUTGOLFCARGO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=18 | `ASM-FW-GISFW-DATA-FACOFFER / COPYALLOBJ_ACT` |
| `CopyAllObjFacOutPAMBU_ACT.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `COPYALLOBJFACOUTPAMBU_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | `ASM-FW-GISFW-DATA-FACOFFER / COPYALLOBJ_ACT` |
| `CopyAllObj_ACT.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `COPYALLOBJ_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=27 | — |
| `CopyClauseFromTemplate_Act.xml` | `ASM-FW-GISFW-WORK` | `COPYCLAUSEFROMTEMPLATE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `CopyCoverageObjItem_Act.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `COPYCOVERAGEOBJITEM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `ASM-FW-GISFW-DATA-PROPERTYITEM / COPYCOVERAGE_ACT` |
| `CopyCoverageToAllLocation_Act.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `COPYCOVERAGETOALLLOCATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=21 | — |
| `CopyCoverage_Act.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `COPYCOVERAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `CopyDeductibleToAllLocation_Act.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `COPYDEDUCTIBLETOALLLOCATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | `ASM-FW-GISFW-DATA-PROPERTYITEM / COPYCOVERAGETOALLLOCATION_ACT` |
| `CopyFacRetroAnekaGolf_ACT.xml` | `ASM-FW-GISFW-WORK` | `COPYFACRETROANEKAGOLF_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=24 | `ASM-FW-GISFW-WORK / COPYFACRETRO_ACT` |
| `CopyFacRetroCargoMBU_ACT.xml` | `ASM-FW-GISFW-WORK` | `COPYFACRETROCARGOMBU_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | `ASM-FW-GISFW-WORK / COPYFACRETRO_ACT` |
| `CopyFacRetroEDMLoc_ACT.xml` | `ASM-FW-GISFW-WORK` | `COPYFACRETROEDMLOC_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CopyFacRetroFire_ACT.xml` | `ASM-FW-GISFW-WORK` | `COPYFACRETROFIRE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | `ASM-FW-GISFW-WORK / COPYFACRETRO_ACT` |
| `CopyFacRetroLife_Act.xml` | `ASM-FW-GISFW-WORK` | `COPYFACRETROLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `ASM-FW-GISFW-WORK / COPYFACRETROPATRAVEL_ACT` |
| `CopyFacRetroPATravel_ACT.xml` | `ASM-FW-GISFW-WORK` | `COPYFACRETROPATRAVEL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `ASM-FW-GISFW-WORK / COPYFACRETRO_ACT` |
| `CopyFacRetro_ACT.xml` | `ASM-FW-GISFW-WORK` | `COPYFACRETRO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `CopyInstallmentRetro.xml` | `ASM-FW-GISFW-WORK` | `COPYINSTALLMENTRETRO` | [terverifikasi] pyActivityType=ACTIVITY; steps=18 | — |
| `CopyLocationFollowing_Act.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `COPYLOCATIONFOLLOWING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=25 | `ASM-FW-GISFW-DATA-OFFERFACIN / COPYLOCATIONFOLLOWING_ACT` |
| `CopyNopolisFollowing_Act.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `COPYNOPOLISFOLLOWING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `CopyObjectRiSlip.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `COPYOBJECTRISLIP` | [terverifikasi] pyActivityType=ACTIVITY; steps=22 | `ASM-FW-GISFW-DATA-FACOFFER / GETRISLIPDATA_ACT` |
| `CopyPaymentTemplate_Act.xml` | `ASM-FW-GISFW-WORK` | `COPYPAYMENTTEMPLATE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CopyTemplateCov_PostAct.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `COPYTEMPLATECOV_POSTACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=32 | — |
| `CopyToAllLocSpreading_ACT.xml` | `ASM-FW-GISFW-DATA` | `COPYTOALLLOCSPREADING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=42 | — |
| `CopyToAllSpreadingAnekaGolf_ACT.xml` | `ASM-FW-GISFW-WORK` | `COPYTOALLSPREADINGANEKAGOLF_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=34 | `ASM-FW-GISFW-WORK / COPYTOALLSPREADING_ACT` |
| `CopyToAllSpreadingCargoMBU_ACT.xml` | `ASM-FW-GISFW-WORK` | `COPYTOALLSPREADINGCARGOMBU_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=31 | `ASM-FW-GISFW-WORK / COPYTOALLSPREADING_ACT` |
| `CopyToAllSpreadingFire_ACT.xml` | `ASM-FW-GISFW-WORK` | `COPYTOALLSPREADINGFIRE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=27 | `ASM-FW-GISFW-WORK / COPYTOALLSPREADING_ACT` |
| `CopyToAllSpreading_ACT.xml` | `ASM-FW-GISFW-WORK` | `COPYTOALLSPREADING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=32 | — |
| `CopyToPolicy.xml` | `ASM-FW-GISFW-WORK` | `COPYTOPOLICY` | [terverifikasi] pyActivityType=ACTIVITY; steps=65 | — |
| `CopyToPolicyList.xml` | `ASM-FW-GISFW-WORK` | `COPYTOPOLICYLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=64 | `ASM-FW-GISFW-WORK / COPYTOPOLICY` |
| `CopyToPolicyListPASSG.xml` | `ASM-FW-GISFW-WORK` | `COPYTOPOLICYLISTPASSG` | [terverifikasi] pyActivityType=ACTIVITY; steps=70 | `ASM-FW-GISFW-WORK / COPYTOPOLICYLIST` |
| `CopypyWorkPage.xml` | `ASM-FW-GISFW-WORK` | `COPYPYWORKPAGE` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CountASMShareTotal_ACT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `COUNTASMSHARETOTAL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GISFW-DATA-OFFERFACIN / COUNTASMSHARETOTAL` |
| `CountAgeShip.xml` | `ASM-FW-GISFW-DATA-CARGO` | `COUNTAGESHIP` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `CountCoverageMarine.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `COUNTCOVERAGEMARINE` | [terverifikasi] pyActivityType=ACTIVITY; steps=21 | — |
| `CountEdmAdjTSI_Act.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `COUNTEDMADJTSI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `CountFormulaRNM_Act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `COUNTFORMULARNM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `CountGPWMarinePAMbu_Act.xml` | `ASM-FW-GISFW-WORK` | `COUNTGPWMARINEPAMBU_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=33 | `ASM-FW-GISFW-WORK / COUNTGROSSPREMI_ACT` |
| `CountGrossPremiEDM_Act.xml` | `ASM-FW-GISFW-WORK` | `COUNTGROSSPREMIEDM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=63 | `ASM-FW-GISFW-WORK / COUNTGROSSPREMI_ACT` |
| `CountGrossPremi_Act.xml` | `ASM-FW-GISFW-WORK` | `COUNTGROSSPREMI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=51 | — |
| `CountNetPremiFOFire.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `COUNTNETPREMIFOFIRE` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CountNetPremi_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTNETPREMI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `CountOGPONP_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTOGPONP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `CountOverridingCommOgp_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTOVERRIDINGCOMMOGP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CountOverridingCommOnp_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTOVERRIDINGCOMMONP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CountPCT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `COUNTPCT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CountPaymentEdmTSIObj_Act.xml` | `ASM-FW-GISFW-WORK` | `COUNTPAYMENTEDMTSIOBJ_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `CountPaymentEdm_Act.xml` | `ASM-FW-GISFW-WORK` | `COUNTPAYMENTEDM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=39 | — |
| `CountPaymentInstallment_Act.xml` | `ASM-FW-GISFW-WORK` | `COUNTPAYMENTINSTALLMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `CountPctInstallment_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTPCTINSTALLMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CountPremiAndTSINusantaraReEDM_ACT.xml` | `ASM-FW-GISFW-WORK` | `COUNTPREMIANDTSINUSANTARAREEDM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=70 | `ASM-FW-GISFW-WORK / COUNTPREMIANDTSINUSANTARARE_ACT` |
| `CountPremiAndTSINusantaraRe_ACT.xml` | `ASM-FW-GISFW-WORK` | `COUNTPREMIANDTSINUSANTARARE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=16 | `ASM-FW-GISFW-DATA-COVERAGE / COUNTPREMIANDTSINUSANTARARE_ACT` |
| `CountPremiAndTSIRNMCargo_ACT.xml` | `ASM-FW-GISFW-WORK` | `COUNTPREMIANDTSIRNMCARGO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=25 | `ASM-FW-GISFW-WORK / COUNTPREMIANDTSIRNMPALIFECARGO_ACT` |
| `CountPremiAndTSIRNMFireAnekaGolf_ACT.xml` | `ASM-FW-GISFW-WORK` | `COUNTPREMIANDTSIRNMFIREANEKAGOLF_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=33 | `ASM-FW-GISFW-WORK / COUNTPREMIANDTSINUSANTARARE_ACT` |
| `CountPremiAndTSIRNMFireMBU_ACT.xml` | `ASM-FW-GISFW-WORK` | `COUNTPREMIANDTSIRNMFIREMBU_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=40 | `ASM-FW-GISFW-WORK / COUNTPREMIANDTSIRNMFIREANEKAGOLF_ACT` |
| `CountPremiAndTSIRNMPALifeCargo_ACT.xml` | `ASM-FW-GISFW-WORK` | `COUNTPREMIANDTSIRNMPALIFECARGO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=20 | `ASM-FW-GISFW-WORK / COUNTPREMIANDTSINUSANTARARE_ACT` |
| `CountPremiCedant_Act.xml` | `ASM-FW-GISFW-WORK` | `COUNTPREMICEDANT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CountPremiCoverageAneka.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `COUNTPREMICOVERAGEANEKA` | [terverifikasi] pyActivityType=ACTIVITY; steps=28 | `ASM-FW-GISFW-DATA-COVERAGE / FILLPREMIGOLF` |
| `CountPremiFacIn_act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `COUNTPREMIFACIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | — |
| `CountPremiNusareRetro_Act.xml` | `ASM-FW-GISFW-WORK` | `COUNTPREMINUSARERETRO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CountPremiRate_ACT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `COUNTPREMIRATE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `ASM-FW-GISFW-DATA-COVERAGE / COUNTPREMI_ACT` |
| `CountPremi_ACT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `COUNTPREMI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=60 | — |
| `CountPremiumNet.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `COUNTPREMIUMNET` | [terverifikasi] pyActivityType=ACTIVITY; steps=30 | — |
| `CountPremiumNetElse.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `COUNTPREMIUMNETELSE` | [terverifikasi] pyActivityType=ACTIVITY; steps=22 | `ASM-FW-GISFW-DATA-FACOFFER / COUNTPREMIUMNET` |
| `CountProrateExtension_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `COUNTPRORATEEXTENSION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `CountRIComm_Act.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `COUNTRICOMM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GISFW-DATA-ANEKA / COUNTRICOMM_ACT` |
| `CountRateRetroCov.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `COUNTRATERETROCOV` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | — |
| `CountResult1Onp_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRESULT1ONP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CountResult1_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRESULT1_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `CountResult2Ogp_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRESULT2OGP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CountResult2Onp_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRESULT2ONP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CountRiCommOgp_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRICOMMOGP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CountRiCommOnp_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRICOMMONP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `CountShareOffered_ACT.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `COUNTSHAREOFFERED_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | — |
| `CountSpreading_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTSPREADING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `CountTSIRNMMultiCob_Act.xml` | `ASM-FW-GISFW-WORK` | `COUNTTSIRNMMULTICOB_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | `ASM-FW-GISFW-WORK / COUNTTOTALTSIPREMINUSARE_ACT` |
| `CountTSIandRateMarineCargo_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `COUNTTSIANDRATEMARINECARGO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CountTotalTSIPremiNusaRe_Act.xml` | `ASM-FW-GISFW-WORK` | `COUNTTOTALTSIPREMINUSARE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=48 | — |
| `DeleteAttachment.xml` | `WORK-` | `DELETEATTACHMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=30 | — |
| `DeleteCeding_Act.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `DELETECEDING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `DeleteDocumentPolis_Act.xml` | `ASM-FW-GISFW-INT-DOCUMENT_POLIS` | `DELETEDOCUMENTPOLIS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `DeleteGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `DELETEGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / GETURLGOOGLESTORAGE_ACT` |
| `DeletePDFFacOutMemoPlacing.xml` | `ASM-FW-GISFW-WORK` | `DELETEPDFFACOUTMEMOPLACING` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GISFW-WORK / GENERATEPDFFACOUTMEMOPLACING` |
| `DeletePDFPenawaranFac.xml` | `ASM-FW-GISFW-WORK` | `DELETEPDFPENAWARANFAC` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `ASM-FW-GISFW-WORK / DELETEPDFRISLIP` |
| `DeletePDFRiSlip.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `DELETEPDFRISLIP` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `DeletePDFSlipKomisi.xml` | `ASM-FW-GISFW-WORK` | `DELETEPDFSLIPKOMISI` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GISFW-WORK / DELETEPDFFACOUTMEMOPLACING` |
| `DetailEDMPaymentHealth.xml` | `ASM-FW-GISFW-WORK` | `DETAILEDMPAYMENTHEALTH` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | — |
| `DetailEDMPaymentHealthCorporate.xml` | `ASM-FW-GISFW-WORK` | `DETAILEDMPAYMENTHEALTHCORPORATE` | [terverifikasi] pyActivityType=ACTIVITY; steps=21 | `ASM-FW-GISFW-WORK / DETAILEDMPAYMENTHEALTH` |
| `DisableCauseOfDecline.xml` | `ASM-FW-GISFW-WORK` | `DISABLECAUSEOFDECLINE` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | — |
| `DoneEditRISlip_Act.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `DONEEDITRISLIP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `DownloadDocumentPolis.xml` | `ASM-FW-GISFW-INT` | `DOWNLOADDOCUMENTPOLIS` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `DownloadMasterLossRecord.xml` | `ASM-FW-GISFW-WORK` | `DOWNLOADMASTERLOSSRECORD` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `ASM-FW-GISFW-DATA-OFFERFACIN / DOWNLOADMASTERLOSSRECORD` |
| `EDM_CekNoMesin.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `EDM_CEKNOMESIN` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `ASM-FW-GISFW-DATA-VEHICLE / EDM_CEKNORANGKA` |
| `EDM_CekNoRangka.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `EDM_CEKNORANGKA` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `ASM-FW-GISFW-DATA-VEHICLE / EDM_CEKNORANGKAMESIN` |
| `EditRISlip_Act.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `EDITRISLIP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `ExportDataMember.xml` | `@BASECLASS` | `EXPORTDATAMEMBER` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GISFW-DATA-POLICY / EXPORTDATAMEMBER` |
| `ExportDataMemberTEST.xml` | `@BASECLASS` | `EXPORTDATAMEMBERTEST` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `@BASECLASS / EXPORTDATAMEMBER` |
| `FetchTreatyGroupOJK.xml` | `ASM-FW-GISFW-WORK-NB` | `FETCHTREATYGROUPOJK` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `FetchTreatyGroupOldID.xml` | `ASM-FW-GISFW-WORK` | `FETCHTREATYGROUPOLDID` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `FillActInstallment.xml` | `ASM-FW-GISFW-WORK` | `FILLACTINSTALLMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=28 | — |
| `FillPaymentInstallment.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `FILLPAYMENTINSTALLMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `FillPremiGolf.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `FILLPREMIGOLF` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `FillPremiMBU_FacIn.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `FILLPREMIMBU_FACIN` | [terverifikasi] pyActivityType=ACTIVITY; steps=29 | — |
| `FilterObjectItem.xml` | `CODE-PEGA-LIST` | `FILTEROBJECTITEM` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GISFW-DATA-PROPERTYITEM / FILTEROBJECTITEM` |
| `FindTemplateAdditionalCoverageSQL_PreAct.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `FINDTEMPLATEADDITIONALCOVERAGESQL_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | `ASM-FW-GISFW-DATA-VEHICLE / SEARCHTEMPLATEADDITIONALCOVERAGESQL_PREACT` |
| `FirstRowSinarmas_Act.xml` | `ASM-FW-GISFW-WORK` | `FIRSTROWSINARMAS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `GISInitAttach.xml` | `@BASECLASS` | `GISINITATTACH` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GISFW-INT-OFFERJSON / GISINITATTACH` |
| `GeminiAIGoogle_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GEMINIAIGOOGLE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | — |
| `GenerateEndorsementPDF.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `GENERATEENDORSEMENTPDF` | [terverifikasi] pyActivityType=ACTIVITY; steps=18 | — |
| `GenerateLayerList_ACT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `GENERATELAYERLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `GenerateNopolis_Act.xml` | `ASM-FW-GISFW-WORK` | `GENERATENOPOLIS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=24 | — |
| `GeneratePDFFacOutMemoPlacing.xml` | `ASM-FW-GISFW-WORK` | `GENERATEPDFFACOUTMEMOPLACING` | [terverifikasi] pyActivityType=ACTIVITY; steps=47 | — |
| `GeneratePDFFacOutOfferStatus.xml` | `ASM-FW-GISFW-WORK` | `GENERATEPDFFACOUTOFFERSTATUS` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `GeneratePDFFacOutViewOffer.xml` | `ASM-FW-GISFW-WORK` | `GENERATEPDFFACOUTVIEWOFFER` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `GeneratePDFSlipKomisi.xml` | `ASM-FW-GISFW-WORK` | `GENERATEPDFSLIPKOMISI` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `ASM-FW-GISFW-WORK / GENERATEPDFFACOUTMEMOPLACING` |
| `GeneratePolicyNoTreaty_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `GENERATEPOLICYNOTREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=32 | — |
| `GenerateRiskAccumReport_Act.xml` | `DATA-PORTAL` | `GENERATERISKACCUMREPORT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `GenerateSpreadingLifeOR_act.xml` | `ASM-FW-GISFW-WORK` | `GENERATESPREADINGLIFEOR_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=23 | `ASM-FW-GISFW-WORK / GENERATESPREADINGLIFE_ACT` |
| `GenerateSpreadingLife_act.xml` | `ASM-FW-GISFW-WORK` | `GENERATESPREADINGLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=25 | — |
| `GetAccumulationRiskDtl_Act.xml` | `DATA-PORTAL` | `GETACCUMULATIONRISKDTL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | — |
| `GetAksepBanding.xml` | `ASM-FW-GISFW-WORK` | `GETAKSEPBANDING` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetAnekaList.xml` | `CODE-PEGA-LIST` | `GETANEKALIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetCLaimList_Act.xml` | `ASM-FW-GISFW-WORK` | `GETCLAIMLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetCauseOfDecline.xml` | `ASM-FW-GISFW-WORK` | `GETCAUSEOFDECLINE` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `GetCoverage_Act.xml` | `CODE-PEGA-LIST` | `GETCOVERAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `GetCurrMasterCargo.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `GETCURRMASTERCARGO` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `GetCurrencyMaster.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `GETCURRENCYMASTER` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `GetCzone_Act.xml` | `DATA-PORTAL` | `GETCZONE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetDataAccumulation_act.xml` | `DATA-PORTAL` | `GETDATAACCUMULATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | — |
| `GetData_ACT.xml` | `DATA-PORTAL` | `GETDATA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `ASSIGN-WORKBASKET / GETDATA_ACT` |
| `GetDateValidity_ACT.xml` | `ASM-FW-GISFW-WORK` | `GETDATEVALIDITY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetDtlTotalAccumulation_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `GETDTLTOTALACCUMULATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | — |
| `GetEDMData.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `GETEDMDATA` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | — |
| `GetEdmProdKe_Act.xml` | `ASM-FW-GISFW-WORK` | `GETEDMPRODKE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `GetGolfCoverage_Act.xml` | `CODE-PEGA-LIST` | `GETGOLFCOVERAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `CODE-PEGA-LIST / GETCOVERAGE_ACT` |
| `GetGolfLocationTopRisk_Act.xml` | `CODE-PEGA-LIST` | `GETGOLFLOCATIONTOPRISK_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `CODE-PEGA-LIST / GETLOCATIONTOPRISK_ACT` |
| `GetGrowingTreeLocationTopRisk_Act.xml` | `CODE-PEGA-LIST` | `GETGROWINGTREELOCATIONTOPRISK_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `CODE-PEGA-LIST / GETGOLFLOCATIONTOPRISK_ACT` |
| `GetGrowingTreesCoverage_Act.xml` | `CODE-PEGA-LIST` | `GETGROWINGTREESCOVERAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `CODE-PEGA-LIST / GETGOLFCOVERAGE_ACT` |
| `GetHistoryAkseptasiPega_Act.xml` | `ASM-FW-GISFW-WORK` | `GETHISTORYAKSEPTASIPEGA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | — |
| `GetInsuredID.xml` | `ASM-FW-GISFW-WORK` | `GETINSUREDID` | [terverifikasi] pyActivityType=ACTIVITY; steps=18 | — |
| `GetKapasitasTreaty.xml` | `ASM-FW-GISFW-WORK` | `GETKAPASITASTREATY` | [terverifikasi] pyActivityType=ACTIVITY; steps=128 | — |
| `GetLicensePlateRegion.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `GETLICENSEPLATEREGION` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `GetLimitAkseptasiLife_Act.xml` | `ASM-FW-GISFW-WORK` | `GETLIMITAKSEPTASILIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `GetLimitAkseptasi_Act.xml` | `ASM-FW-GISFW-WORK` | `GETLIMITAKSEPTASI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=49 | — |
| `GetLimitAkseptasi_ActFlow.xml` | `ASM-FW-GISFW-WORK` | `GETLIMITAKSEPTASI_ACTFLOW` | [terverifikasi] pyActivityType=ACTIVITY; steps=48 | `ASM-FW-GISFW-WORK / GETLIMITAKSEPTASI_ACT` |
| `GetLimitAkseptasi_JUW_UW.xml` | `ASM-FW-GISFW-WORK` | `GETLIMITAKSEPTASI_JUW_UW` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | `ASM-FW-GISFW-WORK / GETLIMITAKSEPTASI_ACTFLOW` |
| `GetLinkService.xml` | `ASM-FW-GISFW-INT-M_LINK_SERVICE` | `GETLINKSERVICE` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetListViewFollowing_act.xml` | `ASM-FW-GISFW-WORK` | `GETLISTVIEWFOLLOWING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=38 | — |
| `GetLocationTopRisk_Act.xml` | `CODE-PEGA-LIST` | `GETLOCATIONTOPRISK_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetLossAmount.xml` | `ASM-FW-GISFW-DATA-CAUSEOFLOSS` | `GETLOSSAMOUNT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `GetLowestPctLimit_ACT.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `GETLOWESTPCTLIMIT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GISFW-DATA-OFFERFACIN / GETLOWESTPCTLIMIT_ACT` |
| `GetOccupation_act.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `GETOCCUPATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `GetOldDataRetro_ACT.xml` | `ASM-FW-GISFW-WORK` | `GETOLDDATARETRO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetPaymentList_Act.xml` | `ASM-FW-GISFW-WORK` | `GETPAYMENTLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `ASM-FW-GISFW-DATA-PROPERTYITEM / GETPAYMENTLIST_ACT` |
| `GetPersonListCoverage_Act.xml` | `CODE-PEGA-LIST` | `GETPERSONLISTCOVERAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `GetPersonListList_Act.xml` | `CODE-PEGA-LIST` | `GETPERSONLISTLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `GetPolicyData_ACT.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETPOLICYDATA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=23 | — |
| `GetPremiPAR_ACT.xml` | `ASM-FW-GISFW-WORK` | `GETPREMIPAR_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `GetRISlipDataFromDB_Act.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `GETRISLIPDATAFROMDB_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetRISlipData_ACT.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `GETRISLIPDATA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=23 | — |
| `GetTeamGroup_Act.xml` | `ASM-FW-GISFW-WORK` | `GETTEAMGROUP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `GetTotalAccumulationLoader_Act.xml` | `ASM-FW-GISFW-DATA-SEARCH` | `GETTOTALACCUMULATIONLOADER_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `GetTreatyName.xml` | `ASM-FW-GISFW-WORK` | `GETTREATYNAME` | [terverifikasi] pyActivityType=ACTIVITY; steps=27 | — |
| `GetUrlGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETURLGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / INSERTGOOGLESTORAGE_ACT` |
| `GetUwYear.xml` | `ASM-FW-GISFW-DATA-CAUSEOFLOSS` | `GETUWYEAR` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `GetZipCodeFromAddress.xml` | `ASM-FW-GISFW-WORK-TKUSALES` | `GETZIPCODEFROMADDRESS` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `HTMLRISlipToPDF.xml` | `@BASECLASS` | `HTMLRISLIPTOPDF` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `@BASECLASS / HTMLTOPDF` |
| `HTMLToPDF.xml` | `@BASECLASS` | `HTMLTOPDF` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `Height_Message_Act.xml` | `DATA-PARTY-PERSON` | `HEIGHT_MESSAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `HitungPremiAndTSINusantaraRe_ACT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `HITUNGPREMIANDTSINUSANTARARE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `IndMedPrintLoading_PreAct.xml` | `ASM-FW-GISFW-WORK` | `INDMEDPRINTLOADING_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=18 | — |
| `IndMedPrintProposal_PreAct.xml` | `ASM-FW-GISFW-WORK` | `INDMEDPRINTPROPOSAL_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=38 | — |
| `InputClause_ViewArg_PreAct.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `INPUTCLAUSE_VIEWARG_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `ASM-FW-GISFW-DATA-CLAUSE / INPUTCLAUSE_VIEWDTL_PREACT` |
| `InputClause_ViewDtl_PreAct.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `INPUTCLAUSE_VIEWDTL_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `InputCurrencyValueAct.xml` | `ASM-FW-GISFW-DATA-POLICY` | `INPUTCURRENCYVALUEACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `InputDtlAdditional_PreAct.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDTLADDITIONAL_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `InputDtlCargo_PreAct.xml` | `ASM-FW-GISFW-DATA-CARGO` | `INPUTDTLCARGO_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `InputDtlClause_PreActMBU.xml` | `ASM-FW-GISFW-WORK` | `INPUTDTLCLAUSE_PREACTMBU` | [terverifikasi] pyActivityType=ACTIVITY; steps=18 | — |
| `InputDtlObjectAneka_PreAct.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `INPUTDTLOBJECTANEKA_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `InputDtlPayment_PreAct.xml` | `ASM-FW-GISFW-WORK` | `INPUTDTLPAYMENT_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=59 | — |
| `InputObjFireGrid_PreAct.xml` | `ASM-FW-GISFW-DATA-PROPERTY` | `INPUTOBJFIREGRID_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `InputOfferFacInEngineerUW_preACT.xml` | `ASM-FW-GISFW-WORK` | `INPUTOFFERFACINENGINEERUW_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=35 | `ASM-FW-GISFW-WORK / INPUTOFFERFACINENGINEER_PREACT` |
| `InputOfferFacInEngineer_preACT.xml` | `ASM-FW-GISFW-WORK` | `INPUTOFFERFACINENGINEER_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=41 | — |
| `InputParamUploadReas_act.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `INPUTPARAMUPLOADREAS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `InputPolicyTreatyInPost_Act.xml` | `ASM-FW-GISFW-WORK` | `INPUTPOLICYTREATYINPOST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `InputPolicyTreatyInPre_Act.xml` | `ASM-FW-GISFW-WORK` | `INPUTPOLICYTREATYINPRE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `InputQuotation_PreAct.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `INPUTQUOTATION_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `InputRateFO.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `INPUTRATEFO` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `ASM-FW-GISFW-DATA-COVERAGE / INPUTRATEFO` |
| `InsertCedingProduction.xml` | `ASM-FW-GISFW-WORK` | `INSERTCEDINGPRODUCTION` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `InsertDocument_Act.xml` | `ASM-FW-GISFW-INT-DOCUMENT_POLIS` | `INSERTDOCUMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM / INSERTDOCUMENT_ACT` |
| `InsertFacInPerformance_Act.xml` | `ASM-FW-GISFW-WORK` | `INSERTFACINPERFORMANCE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=26 | — |
| `InsertFacoutProd.xml` | `ASM-FW-GISFW-WORK` | `INSERTFACOUTPROD` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `InsertFacoutProduction.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `INSERTFACOUTPRODUCTION` | [terverifikasi] pyActivityType=ACTIVITY; steps=36 | — |
| `InsertFacoutProductionEDM.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `INSERTFACOUTPRODUCTIONEDM` | [terverifikasi] pyActivityType=ACTIVITY; steps=59 | `ASM-FW-GISFW-DATA-FACOFFER / INSERTFACOUTPRODUCTION` |
| `InsertFacoutProductionEDMCurr.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `INSERTFACOUTPRODUCTIONEDMCURR` | [terverifikasi] pyActivityType=ACTIVITY; steps=30 | `ASM-FW-GISFW-DATA-FACOFFER / INSERTFACOUTPRODUCTION` |
| `InsertGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERTGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `InsertHistoryAkseptasiPega.xml` | `ASM-FW-GISFW-WORK` | `INSERTHISTORYAKSEPTASIPEGA` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `InsertLogServiceProd.xml` | `ASM-FW-GISFW-WORK` | `INSERTLOGSERVICEPROD` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GCNMFW-WORK / INSERTLOGSERVICECLAIM` |
| `InsertToTreatyXOLList.xml` | `ASM-FW-GISFW-WORK` | `INSERTTOTREATYXOLLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=20 | — |
| `InsertUploadAneka_act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `INSERTUPLOADANEKA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=22 | — |
| `InsertUploadFire_act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `INSERTUPLOADFIRE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=35 | — |
| `InsertUploadMarineCargo_act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `INSERTUPLOADMARINECARGO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=29 | `ASM-FW-GISFW-DATA-OFFERFACIN / INSERTJSONH2HNUSAREMARINECARGO_ACT` |
| `IsThereAnyObjectLocation_Act.xml` | `ASM-FW-GISFW-WORK` | `ISTHEREANYOBJECTLOCATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | — |
| `LoadAttachmentData.xml` | `WORK-` | `LOADATTACHMENTDATA` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `LoadPackingItem.xml` | `ASM-FW-GISFW-DATA-CARGO` | `LOADPACKINGITEM` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `LogActivity_Act.xml` | `ASM-FW-GISFW-WORK` | `LOGACTIVITY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `MappingCoverage.xml` | `ASM-FW-GISFW-INT-TABLE_B2B` | `MAPPINGCOVERAGE` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `MappingJsonH2HNusare2_act.xml` | `ASM-FW-GISFW-INT-TABLE_B2B` | `MAPPINGJSONH2HNUSARE2_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=53 | `ASM-FW-GISFW-INT-TABLE_B2B / MAPPINGJSONH2HNUSARE_ACT` |
| `MappingKlaimToLossRecord_Act.xml` | `ASM-FW-GISFW-WORK` | `MAPPINGKLAIMTOLOSSRECORD_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `MoreThan100_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `MORETHAN100_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `NegativeIsNotAllowed.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `NEGATIVEISNOTALLOWED` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `OfferFacOut_PostAct.xml` | `ASM-FW-GISFW-WORK` | `OFFERFACOUT_POSTACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `OfferFacOut_PreAct.xml` | `ASM-FW-GISFW-WORK` | `OFFERFACOUT_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=31 | — |
| `PercentFireOutgoProtect_Act.xml` | `ASM-FW-GISFW-DATA-OUTGO` | `PERCENTFIREOUTGOPROTECT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `PremiPaymentMarine.xml` | `ASM-FW-GISFW-WORK` | `PREMIPAYMENTMARINE` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `ASM-FW-GISFW-WORK / PREMIPAYMENT` |
| `PrintRISlilpPost_Act.xml` | `ASM-FW-GISFW-WORK` | `PRINTRISLILPPOST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `PrintRISlipPre_act.xml` | `ASM-FW-GISFW-WORK` | `PRINTRISLIPPRE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `PropertyPlanList_AddPropertyAct.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `PROPERTYPLANLIST_ADDPROPERTYACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `ASM-FW-GISFW-DATA-PLAN / PROPERTYPLANLIST_ADDPROPERTYACT` |
| `ProtectAttachRISlip_Act.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `PROTECTATTACHRISLIP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `ProtectCedingCo.xml` | `ASM-FW-GISFW-WORK` | `PROTECTCEDINGCO` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `ProtectCoverage_Act.xml` | `ASM-FW-GISFW-WORK` | `PROTECTCOVERAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=38 | — |
| `ProtectCurrencyTSI_Act.xml` | `ASM-FW-GISFW-WORK` | `PROTECTCURRENCYTSI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `ProtectDate.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `PROTECTDATE` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `ProtectFIREMBUPA_Act.xml` | `ASM-FW-GISFW-WORK` | `PROTECTFIREMBUPA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=89 | `ASM-FW-GISFW-WORK / PROTECTCOVERAGE_ACT` |
| `ProtectObjectMarine.xml` | `ASM-FW-GISFW-DATA-CARGO` | `PROTECTOBJECTMARINE` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | — |
| `ProtectPremiPolicy_Act.xml` | `ASM-FW-GISFW-WORK` | `PROTECTPREMIPOLICY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `ProtectRenewal_Act.xml` | `ASM-FW-GISFW-WORK` | `PROTECTRENEWAL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `ProtectShareCedant_Act.xml` | `ASM-FW-GISFW-WORK` | `PROTECTSHARECEDANT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | — |
| `ProtectShareRNM.xml` | `ASM-FW-GISFW-WORK` | `PROTECTSHARERNM` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `ProtectShipData_Act.xml` | `ASM-FW-GISFW-WORK` | `PROTECTSHIPDATA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `ProtectUploadAneka_Act.xml` | `ASM-FW-GISFW-INT-TABLE_B2B` | `PROTECTUPLOADANEKA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | — |
| `ProtectUploadMarineCargo_act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `PROTECTUPLOADMARINECARGO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=33 | `ASM-FW-GISFW-DATA-OFFERFACIN / INSERTUPLOADMARINECARGO_ACT` |
| `ProtectionEDMFacout_Act.xml` | `ASM-FW-GISFW-WORK` | `PROTECTIONEDMFACOUT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | — |
| `ProtectionObjectAllRisk_Act.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `PROTECTIONOBJECTALLRISK_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `ASM-FW-GISFW-DATA-ANEKA / PROTECTIONOBJECTAVIATION_ACT` |
| `ProtectionObjectAneka_Act.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `PROTECTIONOBJECTANEKA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | — |
| `ProtectionObjectAviation_Act.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `PROTECTIONOBJECTAVIATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | `ASM-FW-GISFW-DATA-ANEKA / PROTECTIONOBJECTGOLFINSURANCE_ACT` |
| `ProtectionObjectBoilerPresure_Act.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `PROTECTIONOBJECTBOILERPRESURE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `ASM-FW-GISFW-DATA-ANEKA / PROTECTIONOBJECTCONTRACTORPLANT_ACT` |
| `ProtectionObjectBurglary_Act.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `PROTECTIONOBJECTBURGLARY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `ASM-FW-GISFW-DATA-ANEKA / PROTECTIONOBJECTALLRISK_ACT` |
| `ProtectionObjectContractorPlant_Act.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `PROTECTIONOBJECTCONTRACTORPLANT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `ASM-FW-GISFW-DATA-ANEKA / PROTECTIONOBJECTLANDRIG_ACT` |
| `ProtectionObjectElectronicEq_Act.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `PROTECTIONOBJECTELECTRONICEQ_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | — |
| `ProtectionObjectFidelity_Act.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `PROTECTIONOBJECTFIDELITY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `ASM-FW-GISFW-DATA-ANEKA / PROTECTIONOBJECTGOLFINSURANCE_ACT` |
| `ProtectionObjectGolfInsurance_Act.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `PROTECTIONOBJECTGOLFINSURANCE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `ASM-FW-GISFW-DATA-ANEKA / PROTECTIONOBJECTHE_ACT` |
| `ProtectionObjectGrowingTree_Act.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `PROTECTIONOBJECTGROWINGTREE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | — |
| `ProtectionObjectHE_Act.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `PROTECTIONOBJECTHE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=21 | — |
| `ProtectionObjectLandRig_Act.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `PROTECTIONOBJECTLANDRIG_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `ASM-FW-GISFW-DATA-ANEKA / PROTECTIONOBJECTHE_ACT` |
| `ProtectionObjectMB_Act.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `PROTECTIONOBJECTMB_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `ASM-FW-GISFW-DATA-ANEKA / PROTECTIONOBJECTCONTRACTORPLANT_ACT` |
| `Protection_Act.xml` | `ASM-FW-GISFW-WORK` | `PROTECTION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=37 | — |
| `ProtekDeduction2_Act.xml` | `ASM-FW-GISFW-WORK` | `PROTEKDEDUCTION2_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `ProtekReinsurerList_Act.xml` | `ASM-FW-GISFW-WORK` | `PROTEKREINSURERLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `ReCountPremiLifeAge_EDM.xml` | `DATA-PARTY-PERSON` | `RECOUNTPREMILIFEAGE_EDM` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | — |
| `ReceiverSumbisProtection_Act.xml` | `ASM-FW-GISFW-DATA-OUTGO` | `RECEIVERSUMBISPROTECTION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `RemoveIFCurrencyIDNull.xml` | `ASM-FW-GISFW-WORK` | `REMOVEIFCURRENCYIDNULL` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `RemoveTypeTax_ACT.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `REMOVETYPETAX_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `ReplaceClauseArgument.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `REPLACECLAUSEARGUMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `ReplaceClauseArgumentFireAct.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `REPLACECLAUSEARGUMENTFIREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `ResetFlagSpreading_ACT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `RESETFLAGSPREADING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `ResetPct_Adjustment.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `RESETPCT_ADJUSTMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `Save.xml` | `WORK-` | `SAVE` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | — |
| `SaveAccumulatedType_Act.xml` | `@BASECLASS` | `SAVEACCUMULATEDTYPE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `@BASECLASS / SAVETYPESHIP_ACT` |
| `SaveAccumulationByRiskAddress_Act.xml` | `@BASECLASS` | `SAVEACCUMULATIONBYRISKADDRESS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `SaveAccumulation_Act.xml` | `@BASECLASS` | `SAVEACCUMULATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `SaveBranch_Act.xml` | `@BASECLASS` | `SAVEBRANCH_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SaveCity_Act.xml` | `@BASECLASS` | `SAVECITY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / SAVECOLOR_ACT` |
| `SaveDataToJsonFollowing_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEDATATOJSONFOLLOWING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SaveDeductible.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `SAVEDEDUCTIBLE` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SaveDeductibleCargo.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `SAVEDEDUCTIBLECARGO` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SaveDistrict_Act.xml` | `@BASECLASS` | `SAVEDISTRICT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / SAVECOLOR_ACT` |
| `SaveEDMToJsonPolicy_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEEDMTOJSONPOLICY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=32 | — |
| `SaveFacIn_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEFACIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GISFW-DATA-OFFERFACIN / SAVEFACIN_ACT` |
| `SaveFacinLive_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEFACINLIVE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `SaveFacinOfferLife_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEFACINOFFERLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SaveFacinProdAllEDM_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEFACINPRODALLEDM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | `ASM-FW-GISFW-WORK / SAVETREATYPRODUCTIONEDM_ACT` |
| `SaveFacinProdCurrAneka_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEFACINPRODCURRANEKA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=30 | — |
| `SaveFacinProdCurrFire_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEFACINPRODCURRFIRE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=37 | — |
| `SaveFacinProdCurrGolf_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEFACINPRODCURRGOLF_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=25 | — |
| `SaveFacinProdCurrMBU_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEFACINPRODCURRMBU_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=25 | — |
| `SaveFacinProdCurrMarineCargo_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEFACINPRODCURRMARINECARGO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=23 | — |
| `SaveFacinProdCurrPA_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEFACINPRODCURRPA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=24 | — |
| `SaveFacinProdEDMFire_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEFACINPRODEDMFIRE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=107 | `ASM-FW-GISFW-WORK / SAVETREATYPRODUCTIONEDMFIRE_ACT` |
| `SaveFacinProdEDMGolf_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEFACINPRODEDMGOLF_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=60 | `ASM-FW-GISFW-WORK / SAVETREATYPRODUCTIONEDMGOLF_ACT` |
| `SaveFacinProdEDMMBU_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEFACINPRODEDMMBU_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=59 | `ASM-FW-GISFW-WORK / SAVETREATYPRODUCTIONEDMMBU_ACT` |
| `SaveFacinProdEDMMarineCargo_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEFACINPRODEDMMARINECARGO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=58 | `ASM-FW-GISFW-WORK / SAVETREATYPRODUCTIONEDMMARINECARGO_ACT` |
| `SaveFacinProdEDMPA_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEFACINPRODEDMPA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=58 | `ASM-FW-GISFW-WORK / SAVETREATYPRODUCTIONEDMPA_ACT` |
| `SaveFacinProdFireNB_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEFACINPRODFIRENB_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=25 | `ASM-FW-GISFW-WORK / SAVETREATYPRODUCTION_ACT` |
| `SaveFacinRNWProd_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEFACINRNWPROD_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=81 | `ASM-FW-GISFW-WORK / SAVETREATYPRODUCTION_ACT` |
| `SaveFacinSpreadLife_Sql.xml` | `ASM-FW-GISFW-WORK` | `SAVEFACINSPREADLIFE_SQL` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SaveFillPaymentInstallment.xml` | `ASM-FW-GISFW-WORK` | `SAVEFILLPAYMENTINSTALLMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | `ASM-FW-GISFW-WORK / FILLPAYMENTINSTALLMENT` |
| `SaveJsonOfferFacIn_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEJSONOFFERFACIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=18 | — |
| `SaveJsonPolicyFacIn_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEJSONPOLICYFACIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=37 | — |
| `SaveJsonPolisTreatyIn_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SAVEJSONPOLISTREATYIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `SaveLossRecord.xml` | `ASM-FW-GISFW-WORK` | `SAVELOSSRECORD` | [terverifikasi] pyActivityType=ACTIVITY; steps=21 | — |
| `SaveMarketingOfficer_Act.xml` | `@BASECLASS` | `SAVEMARKETINGOFFICER_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `@BASECLASS / SAVEASURADUR_ACT` |
| `SaveNation_Act.xml` | `@BASECLASS` | `SAVENATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SaveOfferProduction_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEOFFERPRODUCTION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=80 | `ASM-FW-GISFW-WORK / SAVETREATYPRODUCTION_ACT` |
| `SaveProvince_Act.xml` | `@BASECLASS` | `SAVEPROVINCE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SaveRW_Act.xml` | `@BASECLASS` | `SAVERW_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `@BASECLASS / SAVECOLOR_ACT` |
| `SaveRiskAccumReportToExcel_Act.xml` | `DATA-PORTAL` | `SAVERISKACCUMREPORTTOEXCEL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SaveRiskAddress_Act.xml` | `@BASECLASS` | `SAVERISKADDRESS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `SaveToProduction_ACT.xml` | `ASM-FW-GISFW-WORK` | `SAVETOPRODUCTION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | `ASM-FW-GISFW-WORK / SETTOJSONOFFER_ACT` |
| `SaveTreatyProductionEDMAneka_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVETREATYPRODUCTIONEDMANEKA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=68 | — |
| `SaveTreatyProduction_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVETREATYPRODUCTION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=78 | — |
| `SaveViewSuggest.xml` | `ASM-FW-GISFW-WORK` | `SAVEVIEWSUGGEST` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `ScoringResult.xml` | `ASM-FW-GISFW-DATA-SCORINGRISK` | `SCORINGRESULT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `SearchAccumAct.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SEARCHACCUMACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SearchAccumulationData_Act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SEARCHACCUMULATIONDATA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=16 | — |
| `SearchAccumulationPre_Act.xml` | `DATA-PORTAL` | `SEARCHACCUMULATIONPRE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SearchClauseFireSQL_PostAct.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `SEARCHCLAUSEFIRESQL_POSTACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | — |
| `SearchClauseFireSQL_PreAct.xml` | `ASM-FW-GISFW-WORK` | `SEARCHCLAUSEFIRESQL_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SearchClauseSQL_PostAct.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `SEARCHCLAUSESQL_POSTACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `SearchClauseSQL_PreAct.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `SEARCHCLAUSESQL_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `SearchGroupName_ACT.xml` | `ASM-FW-GISFW-WORK` | `SEARCHGROUPNAME_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SearchInwardScale_Act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SEARCHINWARDSCALE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SearchRiskAddressAct.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `SEARCHRISKADDRESSACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SearchSubContractSQL_PostAct.xml` | `ASM-FW-GISFW-DATA-SUBCONTRACT` | `SEARCHSUBCONTRACTSQL_POSTACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `SearchSubContractSQL_PreAct.xml` | `ASM-FW-GISFW-DATA-SUBCONTRACT` | `SEARCHSUBCONTRACTSQL_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SearchTreatyGroup_Act.xml` | `ASM-FW-GISFW-WORK` | `SEARCHTREATYGROUP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SelectCoverageMarine_Act.xml` | `ASM-FW-GISFW-INT-CONDITION` | `SELECTCOVERAGEMARINE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SelectShip_act.xml` | `ASM-FW-GISFW-INT-SHIP` | `SELECTSHIP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SendCorrespondence.xml` | `ASM-FW-GISFW-DATA-CORRESPONDENCE` | `SENDCORRESPONDENCE` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SendEmailPolicy.xml` | `ASM-FW-GISFW-WORK` | `SENDEMAILPOLICY` | [terverifikasi] pyActivityType=ACTIVITY; steps=16 | — |
| `SendEmailToMO.xml` | `ASM-FW-GISFW-WORK` | `SENDEMAILTOMO` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `SendEmailToUw.xml` | `ASM-FW-GISFW-WORK` | `SENDEMAILTOUW` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `SendEmailWithAttachments.xml` | `@BASECLASS` | `SENDEMAILWITHATTACHMENTS` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SendEmailtoCeding_ACT.xml` | `ASM-FW-GISFW-WORK` | `SENDEMAILTOCEDING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `SendOfferEmail.xml` | `ASM-FW-GISFW-WORK` | `SENDOFFEREMAIL` | [terverifikasi] pyActivityType=ACTIVITY; steps=25 | — |
| `SendSimpleEmail.xml` | `WORK-` | `SENDSIMPLEEMAIL` | [terverifikasi] pyActivityType=ACTIVITY; steps=30 | — |
| `SetAccumulationCode_Act.xml` | `DATA-PORTAL` | `SETACCUMULATIONCODE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetAccumulationData_Act.xml` | `@BASECLASS` | `SETACCUMULATIONDATA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetAutoAccept_Act.xml` | `ASM-FW-GISFW-WORK` | `SETAUTOACCEPT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `SetBanding_ACT.xml` | `ASM-FW-GISFW-WORK` | `SETBANDING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=24 | — |
| `SetBusinessType_Act.xml` | `ASM-FW-GISFW-WORK` | `SETBUSINESSTYPE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetCODRate.xml` | `ASM-FW-GISFW-WORK` | `SETCODRATE` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetCategoryAttach.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SETCATEGORYATTACH` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SetCategoryAttachment_Reas.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `SETCATEGORYATTACHMENT_REAS` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GISFW-INT-OFFERJSON / SETCATEGORYATTACHMENT_REASBU` |
| `SetCauseOfDecline_ACT.xml` | `ASM-FW-GISFW-WORK` | `SETCAUSEOFDECLINE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=20 | `ASM-FW-GISFW-WORK / DISABLECAUSEOFDECLINE` |
| `SetCedingCo_Act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SETCEDINGCO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetClassConstructionEarthquake_ACT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `SETCLASSCONSTRUCTIONEARTHQUAKE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetConveyence_Act.xml` | `ASM-FW-GISFW-DATA-CONVEYANCE` | `SETCONVEYENCE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetCopyRemark_Act.xml` | `ASM-FW-GISFW-WORK` | `SETCOPYREMARK_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `DATA-PARTY-PERSON / SETCOPYREMARK_ACT` |
| `SetCountryID_Act.xml` | `@BASECLASS` | `SETCOUNTRYID_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetCoverageFromMasterPolis_Act.xml` | `ASM-FW-GISFW-DATA-CARGO` | `SETCOVERAGEFROMMASTERPOLIS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetCoverageFrom_PreAct.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SETCOVERAGEFROM_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SetCurrencyCoverage_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SETCURRENCYCOVERAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SetCurrency_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SETCURRENCY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetDataAccum_Act.xml` | `DATA-PORTAL` | `SETDATAACCUM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `ASM-FW-GISFW-DATA-SEARCH / SETDATAACCUM_ACT` |
| `SetDataClassofConstraction.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `SETDATACLASSOFCONSTRACTION` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetDataCoverage.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SETDATACOVERAGE` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetDataFacOutAnekaGolf_Act.xml` | `ASM-FW-GISFW-WORK` | `SETDATAFACOUTANEKAGOLF_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `ASM-FW-GISFW-WORK / SETDATAFACOUT_ACT` |
| `SetDataFacOutCargoMBU_Act.xml` | `ASM-FW-GISFW-WORK` | `SETDATAFACOUTCARGOMBU_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `ASM-FW-GISFW-WORK / SETDATAFACOUT_ACT` |
| `SetDataFacOutFire_Act.xml` | `ASM-FW-GISFW-WORK` | `SETDATAFACOUTFIRE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | `ASM-FW-GISFW-WORK / SETDATAFACOUT_ACT` |
| `SetDataFacOutPATravel_Act.xml` | `ASM-FW-GISFW-WORK` | `SETDATAFACOUTPATRAVEL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | `ASM-FW-GISFW-WORK / SETDATAFACOUT_ACT` |
| `SetDataFacOut_Act.xml` | `ASM-FW-GISFW-WORK` | `SETDATAFACOUT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `ASM-FW-GISFW-WORK / SETVALIDATEDATE_POSTACT` |
| `SetDataKapal_Act.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `SETDATAKAPAL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetDataOccupation.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `SETDATAOCCUPATION` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetDataOccupationScoringRisk_Act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SETDATAOCCUPATIONSCORINGRISK_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `SetDataScoringRisk_act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SETDATASCORINGRISK_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=25 | `ASM-FW-GISFW-WORK / SETDATASCORINGRISK_ACT` |
| `SetDataSobCeding_Act.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `SETDATASOBCEDING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `SetDatePolicy.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `SETDATEPOLICY` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `SetDay.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SETDAY` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetDefaultEmailTo.xml` | `ASM-FW-GISFW-WORK` | `SETDEFAULTEMAILTO` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | — |
| `SetDefaultSearchRiskLocation_Act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `SETDEFAULTSEARCHRISKLOCATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetDownloadMasterCSV.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SETDOWNLOADMASTERCSV` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `SetDueTo_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SETDUETO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GISFW-DATA-POLICYTREATYIN / COUNTNETPREMI_ACT` |
| `SetEndorsementRISlipData.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `SETENDORSEMENTRISLIPDATA` | [terverifikasi] pyActivityType=ACTIVITY; steps=27 | — |
| `SetErrorMessageFacOut.xml` | `ASM-FW-GISFW-DATA-FACOFFEROBJECTLIST` | `SETERRORMESSAGEFACOUT` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | — |
| `SetErrorMessageFloorNumber_Act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `SETERRORMESSAGEFLOORNUMBER_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetErrorMessageTSIObjectItem_Act.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `SETERRORMESSAGETSIOBJECTITEM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `SetErrorMessageUnit_Act.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `SETERRORMESSAGEUNIT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetErrorMessage_Act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SETERRORMESSAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SetFirstLossScale_ACT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SETFIRSTLOSSSCALE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `SetFlagDeleteEDM_Act.xml` | `ASM-FW-GISFW-DATA` | `SETFLAGDELETEEDM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=209 | `ASM-FW-GISFW-DATA / SETFLAGDELETE_ACT` |
| `SetFlagOccupation_ACT.xml` | `ASM-FW-GISFW-WORK` | `SETFLAGOCCUPATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | `ASM-FW-GISFW-WORK / PROTECTFIREMBUPA_ACT` |
| `SetFlagSaveFO_act.xml` | `ASM-FW-GISFW-WORK` | `SETFLAGSAVEFO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SetFloodParamFacIn_ACT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `SETFLOODPARAMFACIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetGoods_Act.xml` | `ASM-FW-GISFW-DATA-GOOD` | `SETGOODS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetIDBackUpStatus.xml` | `ASM-FW-GISFW-WORK` | `SETIDBACKUPSTATUS` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SetIdxCargo.xml` | `ASM-FW-GISFW-INT-SHIP` | `SETIDXCARGO` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GISFW-INT-LLOYDAGENT / SETIDXCARGO` |
| `SetIndemnityRate_ACT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SETINDEMNITYRATE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SetIndexCoverage.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `SETINDEXCOVERAGE` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetIndexMarineCargo_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SETINDEXMARINECARGO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetInsuredName_Act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SETINSUREDNAME_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GISFW-INT-CLIENT / SETINSUREDNAME_ACT` |
| `SetLayerSpreading_act.xml` | `ASM-FW-GISFW-WORK` | `SETLAYERSPREADING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `SetLossRatio_Act.xml` | `ASM-FW-GISFW-DATA-SCORINGRISK` | `SETLOSSRATIO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetMOP_Act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SETMOP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetNBStatus_Act.xml` | `ASM-FW-GISFW-WORK` | `SETNBSTATUS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | `ASM-FW-GISFW-WORK / SETTOINBOX_ACT` |
| `SetNewRate_DT.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `SETNEWRATE_DT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GISFW-DATA-PROPERTYITEM / SETNEWRATE_DT` |
| `SetObjItemType_Act.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `SETOBJITEMTYPE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetOldData.xml` | `ASM-FW-GISFW-WORK` | `SETOLDDATA` | [terverifikasi] pyActivityType=ACTIVITY; steps=59 | — |
| `SetPPNPPH.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SETPPNPPH` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SetPackingNote_Act.xml` | `ASM-FW-GISFW-DATA-GOOD` | `SETPACKINGNOTE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetPolicyType_Act.xml` | `ASM-FW-GISFW-WORK` | `SETPOLICYTYPE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetPostalCode_Act.xml` | `ASM-FW-GISFW-INT-RISKADDRESS` | `SETPOSTALCODE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetPropertyItemNo.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `SETPROPERTYITEMNO` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GISFW-DATA-PROPERTYITEM / SETPROPERTYITEMNO` |
| `SetPropertyToOfferFacIn_Act.xml` | `ASM-FW-GISFW-WORK` | `SETPROPERTYTOOFFERFACIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetProvinceID_Act.xml` | `@BASECLASS` | `SETPROVINCEID_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetRIComm_Act.xml` | `ASM-FW-GISFW-WORK` | `SETRICOMM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=16 | — |
| `SetRatePolis_ACT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SETRATEPOLIS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `SetReinstatementPct.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `SETREINSTATEMENTPCT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SetReinsurerEndorsement_Act.xml` | `ASM-FW-GISFW-WORK` | `SETREINSURERENDORSEMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=33 | — |
| `SetScore_Act.xml` | `ASM-FW-GISFW-DATA-SCORINGRISK` | `SETSCORE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=31 | — |
| `SetShareOfCeding.xml` | `ASM-FW-GISFW-WORK` | `SETSHAREOFCEDING` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetStockAdjustment_ACT.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `SETSTOCKADJUSTMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `SetSubjectToAct.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `SETSUBJECTTOACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SetSurveyAgent_Act.xml` | `ASM-FW-GISFW-WORK` | `SETSURVEYAGENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetSurveyReport_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SETSURVEYREPORT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetTAXNote.xml` | `ASM-FW-GISFW-WORK` | `SETTAXNOTE` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetTSIPremiCedant_Act.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `SETTSIPREMICEDANT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SetTabActive.xml` | `ASM-FW-GISFW-WORK` | `SETTABACTIVE` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetTicket.xml` | `@BASECLASS` | `SETTICKET` | [terverifikasi] pyActivityType=UTILITY; steps=4 | — |
| `SetToInbox_ACT.xml` | `ASM-FW-GISFW-WORK` | `SETTOINBOX_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=30 | `ASM-FW-GISFW-WORK / SETINBOX_ACT` |
| `SetToJsonOffer_ACT.xml` | `ASM-FW-GISFW-WORK` | `SETTOJSONOFFER_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SetTotalFormula_ACT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SETTOTALFORMULA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SetTotalPremium_Act.xml` | `ASM-FW-GISFW-WORK` | `SETTOTALPREMIUM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetTrading_Act.xml` | `ASM-FW-GISFW-DATA-TRADING` | `SETTRADING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetTreatyIn_Act.xml` | `DATA-PORTAL` | `SETTREATYIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | — |
| `SetValidateDateUW_PostAct.xml` | `ASM-FW-GISFW-WORK` | `SETVALIDATEDATEUW_POSTACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=33 | — |
| `SetValidateDate_Act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SETVALIDATEDATE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `SetValidateDate_PostAct.xml` | `ASM-FW-GISFW-WORK` | `SETVALIDATEDATE_POSTACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=30 | — |
| `SetValidateInstallment_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SETVALIDATEINSTALLMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SetValueOnObjectName_Act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `SETVALUEONOBJECTNAME_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `ShowCoverageFacOut.xml` | `ASM-FW-GISFW-WORK` | `SHOWCOVERAGEFACOUT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GISFW-WORK / SHOWSUMINSUREDTSIFACOUT` |
| `ShowGeneratedRISlip.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `SHOWGENERATEDRISLIP` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `ShowInformationOfferFacOut.xml` | `ASM-FW-GISFW-WORK` | `SHOWINFORMATIONOFFERFACOUT` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | — |
| `ShowLocationFacOut.xml` | `ASM-FW-GISFW-WORK` | `SHOWLOCATIONFACOUT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `ShowOccupationFacOut.xml` | `ASM-FW-GISFW-WORK` | `SHOWOCCUPATIONFACOUT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GISFW-WORK / SHOWCOVERAGEFACOUT` |
| `ShowSumInsuredTSIFacOut.xml` | `ASM-FW-GISFW-WORK` | `SHOWSUMINSUREDTSIFACOUT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GISFW-WORK / SHOWLOCATIONFACOUT` |
| `ShowViewCheckListFacOut.xml` | `ASM-FW-GISFW-WORK` | `SHOWVIEWCHECKLISTFACOUT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SortCommentCeding_Act.xml` | `CODE-PEGA-LIST` | `SORTCOMMENTCEDING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `CODE-PEGA-LIST / SORTCOMMETN_ACT` |
| `SpreadingAdditionalProtection.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SPREADINGADDITIONALPROTECTION` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | — |
| `SumCurrencyListAllRetro_Act.xml` | `ASM-FW-GISFW-DATA` | `SUMCURRENCYLISTALLRETRO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `SumFacOutPA_Act.xml` | `ASM-FW-GISFW-WORK` | `SUMFACOUTPA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=21 | — |
| `SumTSIObjItem_act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `SUMTSIOBJITEM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `ASM-FW-GISFW-DATA-OFFERFACIN / SUMTSIOBJITEM_ACT` |
| `SumTSIPremiSpreadedRNM_ANEKA_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SUMTSIPREMISPREADEDRNM_ANEKA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=65 | `ASM-FW-GISFW-DATA-COVERAGE / SUMTSIPREMISPREADEDRNM_ACT` |
| `SumTSIPremiSpreadedRNM_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SUMTSIPREMISPREADEDRNM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=52 | — |
| `SumTSIPremiSpreadedRNM_FIRE_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SUMTSIPREMISPREADEDRNM_FIRE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=103 | `ASM-FW-GISFW-DATA-COVERAGE / SUMTSIPREMISPREADEDRNM_ACT` |
| `SumTSIPremiSpreadedRNM_GOLF_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SUMTSIPREMISPREADEDRNM_GOLF_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=24 | `ASM-FW-GISFW-DATA-COVERAGE / SUMTSIPREMISPREADEDRNM_ACT` |
| `SumTSIPremiSpreadedRNM_MARINECARGO_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SUMTSIPREMISPREADEDRNM_MARINECARGO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `SumTSIPremiSpreadedRNM_MBU_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SUMTSIPREMISPREADEDRNM_MBU_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=24 | `ASM-FW-GISFW-DATA-COVERAGE / SUMTSIPREMISPREADEDRNM_ACT` |
| `SumTSIPremiSpreadedRNM_PA_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SUMTSIPREMISPREADEDRNM_PA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `ASM-FW-GISFW-DATA-COVERAGE / SUMTSIPREMISPREADEDRNM_ACT` |
| `SumTSIPremiSpreadedRNM_TRAVEL_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SUMTSIPREMISPREADEDRNM_TRAVEL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `ASM-FW-GISFW-DATA-COVERAGE / SUMTSIPREMISPREADEDRNM_ACT` |
| `SumTotalTSIPremiGross_Act.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SUMTOTALTSIPREMIGROSS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=39 | — |
| `SumTreatyCapacity_Act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SUMTREATYCAPACITY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | — |
| `SummaryRiskAccumulation_Act.xml` | `DATA-PORTAL` | `SUMMARYRISKACCUMULATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `Test_Act.xml` | `DATA-PARTY-PERSON` | `TEST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `TotalAccumulation_Act.xml` | `ASM-FW-GISFW-WORK` | `TOTALACCUMULATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `TotalBrokerageFee_Act.xml` | `ASM-FW-GISFW-WORK` | `TOTALBROKERAGEFEE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | `ASM-FW-GISFW-WORK / TESTHITUNG_ACT` |
| `TotalEDMPaymentHealth_act.xml` | `ASM-FW-GISFW-WORK` | `TOTALEDMPAYMENTHEALTH_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `TotalKlaimToLossRecord_Act.xml` | `ASM-FW-GISFW-WORK` | `TOTALKLAIMTOLOSSRECORD_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | — |
| `TreatyInInputVis.xml` | `DATA-PORTAL` | `TREATYININPUTVIS` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `TreatyInNonSetTotal.xml` | `DATA-PORTAL` | `TREATYINNONSETTOTAL` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `TreatyRealizationCheckDuplicate.xml` | `ASM-FW-GISFW-WORK` | `TREATYREALIZATIONCHECKDUPLICATE` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `TreatyRealizationCheckXOLList.xml` | `ASM-FW-GISFW-WORK` | `TREATYREALIZATIONCHECKXOLLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `ASM-SFAGIS-WORK / TREATYREALIZATIONCHECKXOLLIST` |
| `TreatySetReinstatement.xml` | `DATA-PORTAL` | `TREATYSETREINSTATEMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `UpdateAccumulationLife_Act.xml` | `ASM-FW-GISFW-WORK` | `UPDATEACCUMULATIONLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | — |
| `UpdateJSONOffer_ACT.xml` | `ASM-FW-GISFW-WORK` | `UPDATEJSONOFFER_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `ASM-FW-GISFW-DATA-FACOFFER / UPDATEJSONPOLICY_ACT` |
| `UpdatePolisAddendum_Act.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `UPDATEPOLISADDENDUM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `UpdatePropertyItemNo.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `UPDATEPROPERTYITEMNO` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `UpdateSectionLocation_ACT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `UPDATESECTIONLOCATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `ASM-FW-GISFW-DATA-OFFERFACIN / UPDATESECTIONBREAKDOWN_ACT` |
| `UpdateStsKonversiFacOut_Act.xml` | `ASM-FW-GISFW-WORK` | `UPDATESTSKONVERSIFACOUT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `UpdateTglBindJsonOffer_Act.xml` | `ASM-FW-GISFW-WORK` | `UPDATETGLBINDJSONOFFER_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `UploadCSVAneka_PostAct.xml` | `ASM-FW-GISFW-WORK` | `UPLOADCSVANEKA_POSTACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=45 | — |
| `UploadCSVCopyTemplateCoverage.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `UPLOADCSVCOPYTEMPLATECOVERAGE` | [terverifikasi] pyActivityType=ACTIVITY; steps=21 | `ASM-FW-GISFW-DATA-VEHICLE / COPYTEMPLATECOVERAGE` |
| `UploadCSVDataPolis_Act.xml` | `@BASECLASS` | `UPLOADCSVDATAPOLIS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `@BASECLASS / UPLOADCSVLIFEPREMIUM_ACT` |
| `UploadCSVPerson_PostAct.xml` | `ASM-FW-GISFW-WORK` | `UPLOADCSVPERSON_POSTACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=91 | — |
| `UploadCSVPolicyMemberEDM_PostAct.xml` | `ASM-FW-GISFW-DATA-POLICY` | `UPLOADCSVPOLICYMEMBEREDM_POSTACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=213 | `ASM-FW-GISFW-DATA-POLICY / UPLOADCSVPOLICYMEMBER_POSTACT` |
| `UploadCSVPolicyMember_PostAct.xml` | `ASM-FW-GISFW-DATA-POLICY` | `UPLOADCSVPOLICYMEMBER_POSTACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=96 | — |
| `UploadCSVSupportCoverageT.xml` | `ASM-FW-GISFW-WORK` | `UPLOADCSVSUPPORTCOVERAGET` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `UploadCSVSupportJob.xml` | `ASM-FW-GISFW-WORK` | `UPLOADCSVSUPPORTJOB` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `UploadCSVSupportVehicle.xml` | `ASM-FW-GISFW-WORK` | `UPLOADCSVSUPPORTVEHICLE` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `UploadCSVSupportVehicleColor.xml` | `ASM-FW-GISFW-WORK` | `UPLOADCSVSUPPORTVEHICLECOLOR` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `UploadCSVSupportVehicleOccu.xml` | `ASM-FW-GISFW-WORK` | `UPLOADCSVSUPPORTVEHICLEOCCU` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `UploadCSVVehicle_PostAct.xml` | `ASM-FW-GISFW-WORK` | `UPLOADCSVVEHICLE_POSTACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=137 | — |
| `UploadDoc_AI_V2.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `UPLOADDOC_AI_V2` | [terverifikasi] pyActivityType=ACTIVITY; steps=44 | `ASM-FW-GISFW-DATA-OFFERFACIN / UPLOADDOC_AI` |
| `ValidateAdjustPct.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `VALIDATEADJUSTPCT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `ValidateCoverage.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `VALIDATECOVERAGE` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | — |
| `ValidateDuplicateVehicle.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `VALIDATEDUPLICATEVEHICLE` | [terverifikasi] pyActivityType=ACTIVITY; steps=33 | — |
| `ValidatePolicyFeeValue.xml` | `ASM-FW-GISFW-WORK` | `VALIDATEPOLICYFEEVALUE` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `ValueScoringRisk_act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `VALUESCORINGRISK_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=27 | — |
| `ViewAttach_RD.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `VIEWATTACH_RD` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `ViewClauseArgFireSQL.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `VIEWCLAUSEARGFIRESQL` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `ViewEndorsementData_Act.xml` | `ASM-FW-GISFW-WORK` | `VIEWENDORSEMENTDATA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `ViewEndorsementRISlipData.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `VIEWENDORSEMENTRISLIPDATA` | [terverifikasi] pyActivityType=ACTIVITY; steps=21 | `ASM-FW-GISFW-DATA-FACOFFER / SETENDORSEMENTRISLIPDATA` |
| `ViewIndemnity2.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `VIEWINDEMNITY2` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `ViewIndemnity_ACT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `VIEWINDEMNITY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `ViewMasterCargo.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `VIEWMASTERCARGO` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `ViewOldDataEDM_Act.xml` | `ASM-FW-GISFW-WORK` | `VIEWOLDDATAEDM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=26 | — |
| `ViewpyWorkPage.xml` | `ASM-FW-GISFW-WORK` | `VIEWPYWORKPAGE` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GISFW-WORK / COPYPYWORKPAGE` |
| `Weight_Message_Act.xml` | `DATA-PARTY-PERSON` | `WEIGHT_MESSAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `addDeductibleCargo_ACT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `ADDDEDUCTIBLECARGO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `addDeductible_ACT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `ADDDEDUCTIBLE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `cekSpreadingFactIn.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `CEKSPREADINGFACTIN` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `cekSpreadingFactIn_Act.xml` | `ASM-FW-GISFW-WORK` | `CEKSPREADINGFACTIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `ASM-FW-GISFW-DATA-COVERAGE / CEKSPREADINGFACTIN` |
| `changepercentrnmLife_act.xml` | `ASM-FW-GISFW-WORK` | `CHANGEPERCENTRNMLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=16 | `ASM-FW-GISFW-WORK / CHANGEPERCENTRNM_ACT` |
| `changepercentrnm_act.xml` | `ASM-FW-GISFW-WORK` | `CHANGEPERCENTRNM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=16 | — |
| `createNewAccumulation_act.xml` | `DATA-PARTY-PERSON` | `CREATENEWACCUMULATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `downloadFile.xml` | `RULE-FILE-BINARY` | `DOWNLOADFILE` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `downloadPenawaranFac_act.xml` | `ASM-FW-GISFW-WORK` | `DOWNLOADPENAWARANFAC_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `fillActAksesoriStandard.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `FILLACTAKSESORISTANDARD` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `fillActBussiness.xml` | `ASM-FW-GISFW-INT-BUSINESS` | `FILLACTBUSSINESS` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `ASM-FW-GISFW-INT-V_TEMPLATE_INSTALLMENT / FILLACTBUSSINESS` |
| `fillActPremiTravel.xml` | `ASM-FW-GISFW-INT-PREMITRAVEL` | `FILLACTPREMITRAVEL` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | `ASM-FW-GISFW-INT-VIEW_PREMI_TRAVEL / FILLACTPREMITRAVEL` |
| `fillPremiAneka.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `FILLPREMIANEKA` | [terverifikasi] pyActivityType=ACTIVITY; steps=24 | `ASM-FW-GISFW-DATA-COVERAGE / FILLPREMIBOND` |
| `getOldDeduct.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `GETOLDDEDUCT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `hitungspreadingotomatis_act.xml` | `ASM-FW-GISFW-WORK` | `HITUNGSPREADINGOTOMATIS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `pzChangeStageWrapper.xml` | `PEGA-API-CASEMANAGEMENT-CASE` | `PZCHANGESTAGEWRAPPERV2` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `removeErrMsg_act.xml` | `DATA-PARTY-PERSON` | `REMOVEERRMSG_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `serviceInsertArasapasEDM_act.xml` | `ASM-FW-GISFW-WORK` | `SERVICEINSERTARASAPASEDM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=20 | `ASM-FW-GISFW-WORK / SERVICEINSERTARASAPAS_ACT` |
| `serviceInsertArasapas_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SERVICEINSERTARASAPAS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `setAdjustmentPercentProtection.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `SETADJUSTMENTPERCENTPROTECTION` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `setDataAkumulasi_act.xml` | `DATA-PARTY-PERSON` | `SETDATAAKUMULASI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `setDiscountRetro_act.xml` | `ASM-FW-GISFW-DATA-SPREADINGRISK` | `SETDISCOUNTRETRO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `setProRatePercent_Act.xml` | `ASM-FW-GISFW-WORK` | `SETPRORATEPERCENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `setRateLifeRetro_act.xml` | `ASM-FW-GISFW-WORK` | `SETRATELIFERETRO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | — |
| `setTSIAdjustment_Act.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `SETTSIADJUSTMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `setTSITopRisk_Act.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SETTSITOPRISK_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `seterrormessagesublimit_ACT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SETERRORMESSAGESUBLIMIT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `svcAddWorkObject.xml` | `WORK-` | `SVCADDWORKOBJECT` | [terverifikasi] pyActivityType=ACTIVITY; steps=23 | — |

### ConnectREST — 3 rule (`RULE-CONNECT-REST`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `SERVICEGOOGLE` | [terverifikasi] pyServiceName=ServiceGoogle; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / GOOGLESTORAGE_UPLOAD` |
| `convertJsonNusareToProduction.xml` | `ASM-FW-GISFW-WORK` | `CONVERTJSONNUSARETOPRODUCTION` | [terverifikasi] pyServiceName=convertJsonNusareToProduction; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | — |
| `getPremiumPaidOn.xml` | `ASM-FW-GISFW-WORK-NB` | `GETPREMIUMPAIDON` | [terverifikasi] pyServiceName=getPremiumPaidOn; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | — |

### DataPage — 30 rule (`RULE-DECLARE-PAGES`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `D_AnekaList.xml` | — | `D_ANEKALIST` | [terverifikasi] pyPageName=D_AnekaList; scope=thread; class=ASM-FW-GISFW-Data-Aneka; struktur=list | `D_ANEKALIST / #20171207T065339.976` |
| `D_BENEFIT_CLAUSE.xml` | — | `D_BENEFIT_CLAUSE` | [terverifikasi] pyPageName=D_BENEFIT_CLAUSE; scope=thread; class=ASM-FW-GISFW-Int-VIEW_BENEFIT_PROPERTY; struktur=list | `D_BENEFIT_CLAUSE / #20151214T040907.404` |
| `D_BankGroup.xml` | — | `D_BANKGROUP` | [terverifikasi] pyPageName=D_BankGroup; scope=node; class=ASM-FW-GISFW-Int-LST_BANK_GROUP; struktur=list | `D_BANKGROUP / #20160427T113654.318` |
| `D_BrowseOccupationFacInFIRE.xml` | — | `D_BROWSEOCCUPATIONFACINFIRE` | [terverifikasi] pyPageName=D_BrowseOccupationFacInFIRE; scope=thread; class=ASM-FW-GISFW-Int-OCCUPATION; struktur=list | `D_BROWSEOCCUPATIONFACINFIRE / #20180305T073030.836` |
| `D_BrowseTableOfLimit.xml` | — | `D_BROWSETABLEOFLIMIT` | [terverifikasi] pyPageName=D_BrowseTableOfLimit; scope=thread; class=ASM-FW-GISFW-Int-TABLEOFLIMIT; struktur=list | `D_BROWSETABLEOFLIMIT / #20180305T081138.390` |
| `D_CITY.xml` | — | `D_CITY` | [terverifikasi] pyPageName=D_CITY; scope=thread; class=ASM-FW-GISFW-Int-CITY; struktur=list | `D_CITY / #20180226T035421.300` |
| `D_CoinsList.xml` | — | `D_COINSLIST` | [terverifikasi] pyPageName=D_CoinsList; scope=requestor; class=ASM-FW-GISFW-Int-V_COINS; struktur=page | `D_COINSLIST / #20151211T024925.977` |
| `D_CommentCeding.xml` | — | `D_COMMENTCEDING` | [terverifikasi] pyPageName=D_CommentCeding; scope=thread; class=ASM-FW-GISFW-Data-OfferFacIn-SuggestList; struktur=list | `D_COMMENT / #20171124T091515.900` |
| `D_Coverage.xml` | — | `D_COVERAGE` | [terverifikasi] pyPageName=D_Coverage; scope=thread; class=ASM-FW-GISFW-Data-Coverage; struktur=page | `D_COVERAGE / #20170731T101130.513` |
| `D_CoverageFacOut.xml` | — | `D_COVERAGEFACOUT` | [terverifikasi] pyPageName=D_CoverageFacOut; scope=thread; class=ASM-FW-GISFW-Data-Coverage; struktur=list | `D_COVERAGEFACOUT / #20170829T041222.879` |
| `D_CoverageSummary.xml` | — | `D_COVERAGESUMMARY` | [terverifikasi] pyPageName=D_CoverageSummary; scope=thread; class=ASM-FW-GISFW-Data-Coverage; struktur=list | `D_COVERAGESUMMARY / #20171120T073200.324` |
| `D_DISTRICT.xml` | — | `D_DISTRICT` | [terverifikasi] pyPageName=D_DISTRICT; scope=thread; class=ASM-FW-GISFW-Int-DISTRICT; struktur=list | `D_DISTRICT / #20180226T035949.969` |
| `D_EnumerationList.xml` | — | `D_ENUMERATIONLIST` | [terverifikasi] pyPageName=D_EnumerationList; scope=thread; class=ASM-FW-GISFW-Data-Enumeration; struktur=list | `D_ENUMERATIONLIST / #20151201T084130.039` |
| `D_FacOfferObjectSection.xml` | — | `D_FACOFFEROBJECTSECTION` | [terverifikasi] pyPageName=D_FacOfferObjectSection; scope=thread; class=ASM-FW-GISFW-Data-FacOfferObjectList; struktur=list | `D_FACOFFEROBJECTSECTION / #20170809T023930.183` |
| `D_FacOutFromVehicle.xml` | — | `D_FACOUTFROMVEHICLE` | [terverifikasi] pyPageName=D_FacOutFromVehicle; scope=thread; class=ASM-FW-GISFW-Data-FacOffer; struktur=list | `D_FACOUTFROMVEHICLE / #20170816T043929.768` |
| `D_FilterObjectItem.xml` | — | `D_FILTEROBJECTITEM` | [terverifikasi] pyPageName=D_FilterObjectItem; scope=thread; class=ASM-FW-GISFW-Data-PropertyItem; struktur=list | `D_FILTEROBJECTITEM / #20171103T042249.055` |
| `D_GolfCoverageSummary.xml` | — | `D_GOLFCOVERAGESUMMARY` | [terverifikasi] pyPageName=D_GolfCoverageSummary; scope=thread; class=ASM-FW-GISFW-Data-Coverage; struktur=list | `D_COVERAGESUMMARY / #20171120T073200.324` |
| `D_GolfLocationSummary.xml` | — | `D_GOLFLOCATIONSUMMARY` | [terverifikasi] pyPageName=D_GolfLocationSummary; scope=thread; class=ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance; struktur=list | `D_LOCATIONSUMMARY / #20171120T044142.723` |
| `D_GrowingTreesCoverageSummary.xml` | — | `D_GROWINGTREESCOVERAGESUMMARY` | [terverifikasi] pyPageName=D_GrowingTreesCoverageSummary; scope=thread; class=ASM-FW-GISFW-Data-Coverage; struktur=list | `D_GOLFCOVERAGESUMMARY / #20171204T030802.539` |
| `D_GrowingTreesLocationSummary.xml` | — | `D_GROWINGTREESLOCATIONSUMMARY` | [terverifikasi] pyPageName=D_GrowingTreesLocationSummary; scope=thread; class=ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance; struktur=list | `D_GOLFLOCATIONSUMMARY / #20171204T072543.165` |
| `D_J_PlanProperty.xml` | — | `D_J_PLANPROPERTY` | [terverifikasi] pyPageName=D_J_PlanProperty; scope=thread; class=ASM-FW-GISFW-Int-VJ_M_TYPE_PROPERTY_PLAN; struktur=list | `D_J_PLANPROPERTY / #20170407T082256.391` |
| `D_LocationSummary.xml` | — | `D_LOCATIONSUMMARY` | [terverifikasi] pyPageName=D_LocationSummary; scope=thread; class=ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance; struktur=list | `D_LOCATIONSUMMARY / #20171120T044142.723` |
| `D_ModelList.xml` | — | `D_MODELLIST` | [terverifikasi] pyPageName=D_ModelList; scope=requestor; class=ASM-FW-GISFW-Int-BRANDDETAIL; struktur=page | `D_MODELLIST / #20170418T071931.623` |
| `D_Occupation.xml` | — | `D_OCCUPATION` | [terverifikasi] pyPageName=D_Occupation; scope=thread; class=ASM-FW-GISFW-Int-OCCUPATION; struktur=list | `D_OCCUPATION / #20180301T041641.366` |
| `D_PersonListCoverageSummary.xml` | — | `D_PERSONLISTCOVERAGESUMMARY` | [terverifikasi] pyPageName=D_PersonListCoverageSummary; scope=thread; class=ASM-FW-GISFW-Data-Coverage; struktur=list | `D_GROWINGTREESCOVERAGESUMMARY / #20171205T100245.349` |
| `D_PersonListSummary.xml` | — | `D_PERSONLISTSUMMARY` | [terverifikasi] pyPageName=D_PersonListSummary; scope=thread; class=Data-Party-Person; struktur=list | `D_GROWINGTREESLOCATIONSUMMARY / #20171205T092615.000` |
| `D_PlanClause.xml` | — | `D_PLANCLAUSE` | [terverifikasi] pyPageName=D_PlanClause; scope=thread; class=ASM-FW-GISFW-Int-M_KLAUSUL_PLAN; struktur=list | `D_PLANCLAUSE / #20160420T142440.434` |
| `D_PlanProperty.xml` | — | `D_PLANPROPERTY` | [terverifikasi] pyPageName=D_PlanProperty; scope=thread; class=ASM-FW-GISFW-Int-M_TYPE_PROPERTY_PLAN; struktur=list | `D_PLANPROPERTY / #20151223T072706.631` |
| `D_RWNAME.xml` | — | `D_RWNAME` | [terverifikasi] pyPageName=D_RWNAME; scope=thread; class=ASM-FW-GISFW-Int-RW; struktur=list | `D_RWNAME / #20180226T040430.984` |
| `D_ZipCodeList.xml` | — | `D_ZIPCODELIST` | [terverifikasi] pyPageName=D_ZipCodeList; scope=node; class=ASM-FW-GISFW-WORK-INT-V_ZIPCODE; struktur=page | `D_ZIPCODELIST / #20151201T102457.729` |

### DataTransform — 148 rule (`RULE-OBJ-MODEL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AddCoverageProRate.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `ADDCOVERAGEPRORATE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `AddPolicySchedule.xml` | `ASM-FW-GISFW-WORK` | `ADDPOLICYSCHEDULE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `AddToListCommentsPolicyTreatyIn_DT.xml` | `ASM-FW-GISFW-WORK` | `ADDTOLISTCOMMENTSPOLICYTREATYIN_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `AddToListSuggestOfferFacIn_DT.xml` | `ASM-FW-GISFW-WORK` | `ADDTOLISTSUGGESTOFFERFACIN_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BackViewPolicy_Dt.xml` | `ASM-FW-GISFW-WORK` | `BACKVIEWPOLICY_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BlankChangeSublimitFacIn_DT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `BLANKCHANGESUBLIMITFACIN_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / BLANKCHANGESUBLIMIT_DT` |
| `BlankValuePropertyItem_Act.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `BLANKVALUEPROPERTYITEM_ACT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CalcLoading_DT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `CALCLOADING_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CalculateAge_DT.xml` | `DATA-PARTY-PERSON` | `CALCULATEAGE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CargoPosition_DT.xml` | `ASM-FW-GISFW-DATA-POLICY` | `CARGOPOSITION_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / CARGOPOSITION_DT` |
| `ChangeCoinsShare.xml` | `ASM-FW-GISFW-DATA-COINS` | `CHANGECOINSSHARE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChangeDiskon.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `CHANGEDISKON` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChangeDiskon_FacIn.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `CHANGEDISKON_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / CHANGEDISKON` |
| `ChangeKomisi.xml` | `ASM-FW-GISFW-DATA-OUTGO` | `CHANGEKOMISI` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CloseInputRW_PreDT.xml` | `@BASECLASS` | `CLOSEINPUTRW_PREDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / CLOSEINPUTRW_PREDT` |
| `CloseMarketingOfficerDT.xml` | `@BASECLASS` | `CLOSEMARKETINGOFFICERDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ConcatAllEmail_DT.xml` | `ASM-FW-GISFW-WORK` | `CONCATALLEMAIL_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ConcatAllPayment_DT.xml` | `ASM-FW-GISFW-WORK` | `CONCATALLPAYMENT_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CorrespondenceAdd.xml` | `ASM-FW-GISFW-WORK` | `CORRESPONDENCEADD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CountAge_DT.xml` | `DATA-PARTY-PERSON` | `COUNTAGE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CountPremiEDMFacOut_DT.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `COUNTPREMIEDMFACOUT_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTYITEM / COUNTPREMIEDMFACOUT_DT` |
| `CountPremiEDM_DT.xml` | `ASM-FW-GISFW-WORK` | `COUNTPREMIEDM_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `D_BenefitCLause_DT.xml` | `CODE-PEGA-LIST` | `D_BENEFITCLAUSE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DeptHeadTreatyInUW_preDT.xml` | `ASM-FW-GISFW-WORK` | `DEPTHEADTREATYINUW_PREDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / INPUTPOLICYTREATYIN_PREDT` |
| `DeptHeadTreatyIn_UW_postDT.xml` | `ASM-FW-GISFW-WORK` | `DEPTHEADTREATYIN_UW_POSTDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ADDTOLISTCOMMENTSPOLICYTREATYIN_DT` |
| `DisplayGridOpenPolicy_DT.xml` | `DATA-PORTAL` | `DISPLAYGRIDOPENPOLICY_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GetLowestPctLimit_DT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `GETLOWESTPCTLIMIT_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `HilangNilaiOutGo_Dt.xml` | `ASM-FW-GISFW-DATA-OUTGO` | `HILANGNILAIOUTGO_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `HitungPremiDT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `HITUNGPREMIDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `HitungPremi_FacInDT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `HITUNGPREMI_FACINDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / HITUNGPREMIDT` |
| `InboxPolicyTreatyIn_postDT.xml` | `ASM-FW-GISFW-WORK` | `INBOXPOLICYTREATYIN_POSTDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IndMedPremium_PreDT.xml` | `ASM-FW-GISFW-WORK` | `INDMEDPREMIUM_PREDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputAccumulatedTypeAdd_DT.xml` | `@BASECLASS` | `INPUTACCUMULATEDTYPEADD_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / INPUTACCUMULATEDTYPEADD_PREDT` |
| `InputAccumulatedType_PreDT.xml` | `@BASECLASS` | `INPUTACCUMULATEDTYPE_PREDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-ACCUMULATEDTYPE / INPUTACCUMULATEDTYPE_PREDT` |
| `InputAccumulationAdd_PreDT.xml` | `DATA-PORTAL` | `INPUTACCUMULATIONADD_PREDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / INPUTTYPESHIPADD_PREDT` |
| `InputAgent_DT.xml` | `@BASECLASS` | `INPUTAGENT_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / INPUTRW_DT` |
| `InputAnekaParticipant_PreDT.xml` | `ASM-FW-GISFW-WORK` | `INPUTANEKAPARTICIPANT_PREDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputBank_DT.xml` | `@BASECLASS` | `INPUTBANK_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCityAdd_DT.xml` | `@BASECLASS` | `INPUTCITYADD_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / INPUTCITYADD_PREDT` |
| `InputCity_PreDT.xml` | `@BASECLASS` | `INPUTCITY_PREDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlObjectAneka_PreDT.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `INPUTDTLOBJECTANEKA_PREDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlOtherMaintenanceAneka_preDT.xml` | `DATA-ADDRESS` | `INPUTDTLOTHERMAINTENANCEANEKA_PREDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlPaymentFire_PreDT.xml` | `ASM-FW-GISFW-WORK` | `INPUTDTLPAYMENTFIRE_PREDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / INPUTDTLPAYMENT_PREDT` |
| `InputDtlPayment_PreDT.xml` | `ASM-FW-GISFW-WORK` | `INPUTDTLPAYMENT_PREDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlWarrantyList_DT.xml` | `ASM-FW-GISFW-WORK-NB` | `INPUTDTLWARRANTYLIST_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-NB / INPCLAUSELIST_DT` |
| `InputPolicyTreatyIn_preDT.xml` | `ASM-FW-GISFW-WORK` | `INPUTPOLICYTREATYIN_PREDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputProvinceAdd_DT.xml` | `@BASECLASS` | `INPUTPROVINCEADD_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / INPUTPROVINCEADD_PREDT` |
| `InputQuotation_PreDT.xml` | `ASM-FW-GISFW-WORK` | `INPUTQUOTATION_PREDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputRWAdd_DT.xml` | `@BASECLASS` | `INPUTRWADD_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / INPUTCARCLASSADD_DT` |
| `InputRiskAddress_DT.xml` | `@BASECLASS` | `INPUTRISKADDRESS_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputScoringRisk_PreDT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `INPUTSCORINGRISK_PREDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InwardFacultative_PreDT.xml` | `ASM-FW-GISFW-WORK` | `INWARDFACULTATIVE_PREDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsJuniorUW_B.xml` | `ASM-FW-GISFW-WORK` | `ISJUNIORUW_B` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `LoadPackingItem.xml` | `ASM-FW-GISFW-DATA-CARGO` | `LOADPACKINGITEM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PersonPlanAndDeduct_PostDT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `PERSONPLANANDDEDUCT_POSTDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `RecipientListAdd.xml` | `ASM-FW-GISFW-DATA-CORRESPONDENCE` | `RECIPIENTLISTADD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `RemoveClauseList.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `REMOVECLAUSELIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `RemoveDiscount.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `REMOVEDISCOUNT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `RemoveSubContarctList.xml` | `ASM-FW-GISFW-DATA-SUBCONTRACT` | `REMOVESUBCONTARCTLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ResetFilterValue_DT.xml` | `ASM-FW-GISFW-WORK` | `RESETFILTERVALUE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ResetNumberofLayer_DT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `RESETNUMBEROFLAYER_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SearchHierarkiSourceBizAgent_PostDT.xml` | `ASM-FW-GISFW-DATA-AGENT` | `SEARCHHIERARKISOURCEBIZAGENT_POSTDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetAccumSearchPreDT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `SETACCUMSEARCHPREDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetAdditionalPremiumDeductible_DT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SETADDITIONALPREMIUMDEDUCTIBLE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / SETPREMIUMDEDUCTIBLE_DT` |
| `SetAkseptasiCeding_DT.xml` | `ASM-FW-GISFW-WORK` | `SETAKSEPTASICEDING_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / SETAKSEPTASIPROPOSAL` |
| `SetAkseptasiProposal.xml` | `ASM-FW-GISFW-WORK` | `SETAKSEPTASIPROPOSAL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetAkseptasiProposal_DT.xml` | `ASM-FW-GISFW-WORK` | `SETAKSEPTASIPROPOSAL_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetAskProposal_DT.xml` | `ASM-FW-GISFW-WORK` | `SETASKPROPOSAL_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetBandingProposal_DT.xml` | `ASM-FW-GISFW-WORK` | `SETBANDINGPROPOSAL_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / SETASKPROPOSAL_DT` |
| `SetCalculateMethodNDay_DT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SETCALCULATEMETHODNDAY_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetCedingCedant_DT.xml` | `ASM-FW-GISFW-DATA-AGENT` | `SETCEDINGCEDANT_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-AGENT / SETCEDINGCO_DT` |
| `SetCodeRiskExposure_DT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `SETCODERISKEXPOSURE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetCorrSubscript.xml` | `ASM-FW-GISFW-DATA-CORRESPONDENCE` | `SETCORRSUBSCRIPT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetCorrespondenceReadOnly.xml` | `ASM-FW-GISFW-WORK` | `SETCORRESPONDENCEREADONLY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetCoverageDay_DT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SETCOVERAGEDAY_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetCoverageFromTemplate_DT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SETCOVERAGEFROMTEMPLATE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetDataQuotation_DT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SETDATAQUOTATION_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN / SETNOOFFER_DT` |
| `SetDeclineProposal_DT.xml` | `ASM-FW-GISFW-WORK` | `SETDECLINEPROPOSAL_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetDeductible_DT.xml` | `ASM-FW-GISFW-INT-DEDUCTIBLE` | `SETDEDUCTIBLE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetDetailViewOpenFollowingPolicy.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `SETDETAILVIEWOPENFOLLOWINGPOLICY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-POLICYJSON / SETDETAILVIEWOPENPOLICYMARINE_PREDT` |
| `SetDistirctFacIn_DT.xml` | `@BASECLASS` | `SETDISTIRCTFACIN_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / SETDISTIRCT_DT` |
| `SetFilterAccum_DT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SETFILTERACCUM_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetFlagDelete_DT.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `SETFLAGDELETE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-VEHICLE / SETFLAGDELETE_DT` |
| `SetFlagWpcExtend_DT.xml` | `ASM-FW-GISFW-DATA-INSTALLMENT` | `SETFLAGWPCEXTEND_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetIdxCoverage_DT.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `SETIDXCOVERAGE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTY / SETCOVERAGENODT` |
| `SetInFacRetroStatus_PreDT.xml` | `ASM-FW-GISFW-WORK` | `SETINFACRETROSTATUS_PREDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetIndexCoverage_DT.xml` | `ASM-FW-GISFW-WORK` | `SETINDEXCOVERAGE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / SETINDEXCOVERAGE_DT` |
| `SetIndexCurrency_DT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | `SETINDEXCURRENCY_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN / SETINDEXCURRENCY_DT` |
| `SetIndexDeductible_DT.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `SETINDEXDEDUCTIBLE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetIndexDt.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `SETINDEXDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetIndexFacRetroList_DT.xml` | `ASM-FW-GISFW-WORK` | `SETINDEXFACRETROLIST_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetIndexHE_DT.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `SETINDEXHE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetIndexObject_DT.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `SETINDEXOBJECT_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetIndexSpreading_Dt.xml` | `ASM-FW-GISFW-DATA-SPREADINGRISK` | `SETINDEXSPREADING_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetInputNationAdd_DT.xml` | `@BASECLASS` | `SETINPUTNATIONADD_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / INPUTNATIONADD_PREDT` |
| `SetLetterNo.xml` | `ASM-FW-GISFW-WORK` | `SETLETTERNO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISJUNIORUW_B` |
| `SetNation_DT.xml` | `@BASECLASS` | `SETNATION_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / SETSHIP_DT` |
| `SetObjNameDT.xml` | `ASM-FW-GISFW-DATA-PROPERTY` | `SETOBJNAMEDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetObjNoDT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SETOBJNODT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / SETOBJNODT` |
| `SetObjectDescHE_DT.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `SETOBJECTDESCHE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetOcId_DT.xml` | `ASM-FW-GISFW-DATA-OUTGO` | `SETOCID_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / SETOCID_DT` |
| `SetOccupation_DT.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `SETOCCUPATION_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetOutgoByTemplate_DT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SETOUTGOBYTEMPLATE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / SETOUTGOTEMPLATE_DT` |
| `SetParamIndex_DT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `SETPARAMINDEX_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-ANEKA / SETPARAMINDEX_DT` |
| `SetPassword_DT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SETPASSWORD_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetPercentPremiSpreaded.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SETPERCENTPREMISPREADED` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetPremiumCoverage.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `SETPREMIUMCOVERAGE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetPremiumDeductible_DT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SETPREMIUMDEDUCTIBLE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetPropertyCoverageDT_FacIn.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SETPROPERTYCOVERAGEDT_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetPropertyItemIndex_DT.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `SETPROPERTYITEMINDEX_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetProvince_DT.xml` | `@BASECLASS` | `SETPROVINCE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / SETCITY_DT` |
| `SetRateFacOffer_DT.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `SETRATEFACOFFER_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-ANEKA / SETRATEFACOFFER_DT` |
| `SetReinsurer_DT.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `SETREINSURER_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetRejectProposal.xml` | `ASM-FW-GISFW-WORK` | `SETREJECTPROPOSAL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetReviseProposal.xml` | `ASM-FW-GISFW-WORK` | `SETREVISEPROPOSAL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / SETREJECTPROPOSAL` |
| `SetRiskIdDT_FacIn.xml` | `ASM-FW-GISFW-INT-RISKADDRESS` | `SETRISKIDDT_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetRoadNametoUpper_DT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `SETROADNAMETOUPPER_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetShowZone_DT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SETSHOWZONE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetSpreadingCoverage_DT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SETSPREADINGCOVERAGE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetStatusCoverage_DT.xml` | `ASM-FW-GISFW-WORK-NB` | `SETSTATUSCOVERAGE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetTSICoverage_DT.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `SETTSICOVERAGE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetTSIPremiSpreaded_FacIn.xml` | `ASM-FW-GISFW-WORK` | `SETTSIPREMISPREADED_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetTSISublimitDT_FacIn.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SETTSISUBLIMITDT_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / SETTSISUBLIMITDT` |
| `SetTSIValue_DT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `SETTSIVALUE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetTemplateAdditionalCoverageIndex_DT.xml` | `ASM-FW-GISFW-DATA-PREMIUMTEMPLATE` | `SETTEMPLATEADDITIONALCOVERAGEINDEX_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PREMIUMTEMPLATE / SETTEMPLATEADDITIONALCOVERAGE_DT` |
| `SetTemplateAdditionalCoverage_DT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SETTEMPLATEADDITIONALCOVERAGE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetTemplateIndex_DT.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SETTEMPLATEINDEX_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetToUW_DT.xml` | `ASM-FW-GISFW-WORK` | `SETTOUW_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Set_BMI_DT.xml` | `DATA-PARTY-PERSON` | `SET_BMI_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Set_PersonNoteDokter_DT.xml` | `DATA-PARTY-PERSON` | `SET_PERSONNOTEDOKTER_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / SET_PERSONNOTE_DT` |
| `Set_PersonNote_DT.xml` | `DATA-PARTY-PERSON` | `SET_PERSONNOTE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowRadioButtonAccept_DT.xml` | `ASM-FW-GISFW-WORK` | `SHOWRADIOBUTTONACCEPT_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SystemSetOneYear_DT.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SYSTEMSETONEYEAR_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TestTreatyToFacStatus.xml` | `ASM-FW-GISFW-WORK` | `TESTTREATYTOFACSTATUS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyEnableDisableInput.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `TREATYENABLEDISABLEINPUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `VehicleAge_DT.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `VEHICLEAGE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewOfferFacInUW_PreDT.xml` | `ASM-FW-GISFW-WORK` | `VIEWOFFERFACINUW_PREDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `btnCedingCO_DT.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `BTNCEDINGCO_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `btnCedingCedant_DT.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `BTNCEDINGCEDANT_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-QUOTATION / BTNCEDINGCO_DT` |
| `btnSOB_DT.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `BTNSOB_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `fillPremiBond.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `FILLPREMIBOND` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `insertCalculatedMethod.xml` | `DATA-PARTY-PERSON` | `INSERTCALCULATEDMETHOD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `pyDefault.xml` | `DATA-WORKATTACH-FILE` | `PYDEFAULT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `setCategoryAttachment_DT.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SETCATEGORYATTACHMENT_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `setIndexProperty.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `SETINDEXPROPERTY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `setIndexPropertyFacOut_DT.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `SETINDEXPROPERTYFACOUT_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTYITEM / SETINDEXPROPERTY` |
| `setOutputValue_DT.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `SETOUTPUTVALUE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `setParamIdxVehicle_DT.xml` | `ASM-FW-GISFW-WORK` | `SETPARAMIDXVEHICLE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### DecisionTable — 12 rule (`RULE-DECLARE-DECISIONTABLE`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BusinessType_DeT.xml` | `ASM-FW-GISFW-WORK` | `BUSINESSTYPE_DET` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GetMimeType.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETMIMETYPE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsUWAccepted.xml` | `ASM-FW-GISFW-WORK` | `ISUWACCEPTED` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `LicensePlatRegion_DeT.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `LICENSEPLATREGION_DET` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `MappingAdditionalCoverageIndex.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `MAPPINGADDITIONALCOVERAGEINDEX` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `MappingCoverage.xml` | `ASM-FW-GISFW-WORK` | `MAPPINGCOVERAGE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `MappingOutgoIndex.xml` | `ASM-FW-GISFW-DATA-SPLIT` | `MAPPINGOUTGOINDEX` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `MappingOutgoIndex2.xml` | `ASM-FW-GISFW-DATA-SPLIT` | `MAPPINGOUTGOINDEX2` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-SPLIT / MAPPINGOUTGOINDEX` |
| `SetUploadHubAW1.xml` | `ASM-FW-GISFW-DATA-BATCHPALIST` | `SETUPLOADHUBAW1` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-BATCHTRAVELLIST / SETUPLOADHUBAW1` |
| `SetUploadHubAW2.xml` | `ASM-FW-GISFW-DATA-BATCHPALIST` | `SETUPLOADHUBAW2` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-BATCHTRAVELLIST / SETUPLOADHUBAW2` |
| `SetUploadHubAW3.xml` | `ASM-FW-GISFW-DATA-BATCHPALIST` | `SETUPLOADHUBAW3` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-BATCHTRAVELLIST / SETUPLOADHUBAW3` |
| `isApproved.xml` | `ASM-FW-GISFW-WORK` | `ISAPPROVED` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### DecisionTree — 1 rule (`RULE-DECLARE-DECISIONTREE`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `Tree_ShortPeriod.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `TREE_SHORTPERIOD` | [terverifikasi] class=ASM-FW-GISFW-Data-Coverage | — |

### Flow — 6 rule (`RULE-OBJ-FLOW`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `InputInwardFacultativeOffer.xml` | `ASM-FW-GISFW-WORK` | `INPUTINWARDFACULTATIVEOFFER` | [terverifikasi] pyStartActivity=Start1 | — |
| `InputInwardFacultativeRISlip.xml` | `ASM-FW-GISFW-WORK` | `INPUTINWARDFACULTATIVERISLIP` | [terverifikasi] pyStartActivity=Start1 | — |
| `InputQuotation.xml` | `ASM-FW-GISFW-WORK-NB` | `INPUTQUOTATION` | [terverifikasi] pyStartActivity=Start1 | — |
| `InputRealizationTreatyIn.xml` | `ASM-FW-GISFW-WORK` | `INPUTREALIZATIONTREATYIN` | [terverifikasi] pyStartActivity=Start1 | — |
| `OfferFacOut.xml` | `ASM-FW-GISFW-WORK` | `OFFERFACOUT` | [terverifikasi] pyStartActivity=Start62 | — |
| `OfferFacRetro.xml` | `ASM-FW-GISFW-WORK` | `OFFERFACRETRO` | [terverifikasi] pyStartActivity=Start1 | — |

### FlowAction — 250 rule (`RULE-OBJ-FLOWACTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AddFacOfferList.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `ADDFACOFFERLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-FACOFFER / TESFACOUT` |
| `AddFacOfferList_IsUW.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `ADDFACOFFERLIST_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-FACOFFER / ADDFACOFFERLIST` |
| `AgentSourceBizDetails.xml` | `ASM-FW-GISFW-DATA-AGENT` | `AGENTSOURCEBIZDETAILS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `AnalysLocationbyAI.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `ANALYSLOCATIONBYAI` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `AnekaListFacOut_FlowAction.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `ANEKALISTFACOUT_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `AnekaListFacOut_FlowAction_IsUW.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `ANEKALISTFACOUT_FLOWACTION_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OCCUPATION / ANEKALISTFACOUT_FLOWACTION` |
| `AttachContentGIS.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `ATTACHCONTENTGIS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `AttachDoc_AI.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `ATTACHDOC_AI` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BenefitListRO.xml` | `ASM-FW-GISFW-DATA-PLAN` | `BENEFITLISTRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CauseOfLossGrid.xml` | `ASM-FW-GISFW-DATA-CAUSEOFLOSS` | `CAUSEOFLOSSGRID` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CedingCedant.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `CEDINGCEDANT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CedingCedant_IsUW.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `CEDINGCEDANT_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-QUOTATION / CEDINGCEDANT` |
| `ChangeInsuredName.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `CHANGEINSUREDNAME` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChooseClassofContraction.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `CHOOSECLASSOFCONTRACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChooseClause.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `CHOOSECLAUSE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChooseCoverage.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `CHOOSECOVERAGE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChooseObject_Ship.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `CHOOSEOBJECT_SHIP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChooseOccupation.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `CHOOSEOCCUPATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChooseRiskAddress.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `CHOOSERISKADDRESS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChooseSubContract.xml` | `ASM-FW-GISFW-DATA-SUBCONTRACT` | `CHOOSESUBCONTRACT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChooseSubContract_FacIn.xml` | `ASM-FW-GISFW-DATA-SUBCONTRACT` | `CHOOSESUBCONTRACT_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-SUBCONTRACT / CHOOSESUBCONTRACT` |
| `ChooseZipCode.xml` | `@BASECLASS` | `CHOOSEZIPCODE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CommisionCoverageAneka_FlowAction.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `COMMISIONCOVERAGEANEKA_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-ANEKA / SPREADINGCOVERAGEANEKA_FLOWACTION` |
| `CommisionItem.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `COMMISIONITEM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / SPREADINGITEM` |
| `CommisionItemAneka.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `COMMISIONITEMANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CommisionObjectAneka_FlowAction.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `COMMISIONOBJECTANEKA_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OCCUPATION / SPREADINGOBJECTANEKA_FLOWACTION` |
| `ConfirmOutgo.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `CONFIRMOUTGO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CopyCoverageFrom.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `COPYCOVERAGEFROM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CorrespondenceContent.xml` | `ASM-FW-GISFW-DATA-CORRESPONDENCE` | `CORRESPONDENCECONTENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CoverageFireGrid_GCNM.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `COVERAGEFIREGRID_GCNM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CoverageItem.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `COVERAGEITEM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CoverageItem_IsUW.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `COVERAGEITEM_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / COVERAGEITEM` |
| `CoverageListFire_FlowAction.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `COVERAGELISTFIRE_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CoverageListFire_FlowAction_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `COVERAGELISTFIRE_FLOWACTION_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / COVERAGELISTFIRE_FLOWACTION` |
| `DeductibleAnekaGridGCNM.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `DEDUCTIBLEANEKAGRIDGCNM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DeductibleDtl.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `DEDUCTIBLEDTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DeductibleDtlAnekaGridGCNM.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `DEDUCTIBLEDTLANEKAGRIDGCNM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DeductibleHullGridGCNM.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `DEDUCTIBLEHULLGRIDGCNM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DeptHeadTreatyIn_UW.xml` | `ASM-FW-GISFW-WORK` | `DEPTHEADTREATYIN_UW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `FEAGrid.xml` | `ASM-FW-GISFW-DATA-FEA` | `FEAGRID` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `FlowActCoverageListFacOut.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `FLOWACTCOVERAGELISTFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-VEHICLE / FLOWACTCOVERAGELIST` |
| `FlowActCoverageListFacOut_IsUW.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `FLOWACTCOVERAGELISTFACOUT_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-VEHICLE / FLOWACTCOVERAGELISTFACOUT` |
| `FlowActViewCoverageList.xml` | `ASM-FW-GISFW-DATA-CARGO` | `FLOWACTVIEWCOVERAGELIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-VEHICLE / FLOWACTVIEWCOVERAGELIST` |
| `FlowActViewDtlOutgo.xml` | `ASM-FW-GISFW-WORK` | `FLOWACTVIEWDTLOUTGO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / FLOWACTINPUTDTLOUTGO` |
| `GolfLocationDetail_FlowAction.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `GOLFLOCATIONDETAIL_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / LOCATIONDETAIL_FLOWACTION` |
| `GrowingTreesLocationDetail_FlowAction.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `GROWINGTREESLOCATIONDETAIL_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / GOLFLOCATIONDETAIL_FLOWACTION` |
| `InboxPolicyTreatyIn.xml` | `ASM-FW-GISFW-WORK` | `INBOXPOLICYTREATYIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InpClause.xml` | `ASM-FW-GISFW-WORK` | `INPCLAUSE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-NB / INPCLAUSE` |
| `InputAnekaParticipant.xml` | `ASM-FW-GISFW-WORK` | `INPUTANEKAPARTICIPANT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputAnekaPolicySchedule.xml` | `ASM-FW-GISFW-WORK` | `INPUTANEKAPOLICYSCHEDULE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCauseOfLoss_FacIn.xml` | `ASM-FW-GISFW-DATA-CAUSEOFLOSS` | `INPUTCAUSEOFLOSS_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCauseOfLoss_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-CAUSEOFLOSS` | `INPUTCAUSEOFLOSS_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-CAUSEOFLOSS / INPUTCAUSEOFLOSS_FACIN` |
| `InputCity.xml` | `@BASECLASS` | `INPUTCITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / MODALMAP` |
| `InputClauseFire_ViewDtl.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `INPUTCLAUSEFIRE_VIEWDTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputClause_ViewArg.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `INPUTCLAUSE_VIEWARG` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputClause_ViewDtl.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `INPUTCLAUSE_VIEWDTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCommentDeductible.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `INPUTCOMMENTDEDUCTIBLE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN / EDITMARKETING` |
| `InputCommissionLife_FacIn.xml` | `DATA-PARTY-PERSON` | `INPUTCOMMISSIONLIFE_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCommissionMBU_FacIn.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `INPUTCOMMISSIONMBU_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCommissionMC_FacIn.xml` | `ASM-FW-GISFW-DATA-CARGO` | `INPUTCOMMISSIONMC_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCommissionPA_FacIn.xml` | `DATA-PARTY-PERSON` | `INPUTCOMMISSIONPA_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCommissionTravel_FacIn.xml` | `DATA-PARTY-PERSON` | `INPUTCOMMISSIONTRAVEL_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / INPUTCOMMISSIONPA_FACIN` |
| `InputCoverageCommision.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `INPUTCOVERAGECOMMISION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / INPUTCOVERAGESPREADING` |
| `InputCoverageSpreading.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `INPUTCOVERAGESPREADING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCoverageSpreading_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `INPUTCOVERAGESPREADING_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / INPUTCOVERAGESPREADING` |
| `InputCoverageTravelGrid_GCNM.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTCOVERAGETRAVELGRID_GCNM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / DEDUCTIBLETRAVELGRID_GCNM` |
| `InputDistrict.xml` | `@BASECLASS` | `INPUTDISTRICT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / INPUTCITY` |
| `InputDistrik.xml` | `@BASECLASS` | `INPUTDISTRIK` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / INPUTKABUPATEN` |
| `InputDtlCargo_FacIn.xml` | `ASM-FW-GISFW-DATA-CARGO` | `INPUTDTLCARGO_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-CARGO / INPUTDTLCARGO` |
| `InputDtlCargo_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-CARGO` | `INPUTDTLCARGO_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-CARGO / INPUTDTLCARGO_FACIN` |
| `InputDtlCommissionLife_FacIn.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDTLCOMMISSIONLIFE_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlConveyance.xml` | `ASM-FW-GISFW-DATA-CONVEYANCE` | `INPUTDTLCONVEYANCE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlConveyance_IsUW.xml` | `ASM-FW-GISFW-DATA-CONVEYANCE` | `INPUTDTLCONVEYANCE_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-CONVEYANCE / INPUTDTLCONVEYANCE` |
| `InputDtlCoverage_FacIn.xml` | `ASM-FW-GISFW-DATA-CARGO` | `INPUTDTLCOVERAGE_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-CARGO / INPUTDTLCOVERAGE` |
| `InputDtlCoverage_Marine.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDTLCOVERAGE_MARINE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlCoverage_Marine_IsUW.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDTLCOVERAGE_MARINE_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlDeductibleCargo_FacIn.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `INPUTDTLDEDUCTIBLECARGO_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlDeductibleCargo_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `INPUTDTLDEDUCTIBLECARGO_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlDeductibleFire_FacIn.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `INPUTDTLDEDUCTIBLEFIRE_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlDeductibleFire_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `INPUTDTLDEDUCTIBLEFIRE_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-DEDUCTIBLE / INPUTDTLDEDUCTIBLEFIRE_FACIN` |
| `InputDtlDeductible_FacIn.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `INPUTDTLDEDUCTIBLE_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlDeductible_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `INPUTDTLDEDUCTIBLE_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-DEDUCTIBLE / INPUTDTLDEDUCTIBLE_FACIN` |
| `InputDtlFormulaFlow_FacIn.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | `INPUTDTLFORMULAFLOW_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN / INPUTDTLFORMULAFLOW_FACIN` |
| `InputDtlGoods.xml` | `ASM-FW-GISFW-DATA-GOOD` | `INPUTDTLGOODS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlGoods_IsUW.xml` | `ASM-FW-GISFW-DATA-GOOD` | `INPUTDTLGOODS_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-GOOD / INPUTDTLGOODS` |
| `InputDtlLayerListCommision_FacIn.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDTLLAYERLISTCOMMISION_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLLAYERLISTSPREADING_FACIN` |
| `InputDtlLayerListFire_FacIn.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDTLLAYERLISTFIRE_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlLayerListFire_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDTLLAYERLISTFIRE_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLLAYERLISTFIRE_FACIN` |
| `InputDtlLayerListSpreading_FacIn.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDTLLAYERLISTSPREADING_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlLayerListSpreading_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDTLLAYERLISTSPREADING_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLLAYERLISTSPREADING_FACIN` |
| `InputDtlLayerSpreading_FacIn.xml` | `ASM-FW-GISFW-DATA-SPREADINGRISK` | `INPUTDTLLAYERSPREADING_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlObjectAnekaGCNM.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `INPUTDTLOBJECTANEKAGCNM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlObjectAneka_FacIn.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `INPUTDTLOBJECTANEKA_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlObjectAneka_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `INPUTDTLOBJECTANEKA_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-ANEKA / INPUTDTLOBJECTANEKA_FACIN` |
| `InputDtlObjectHull_GCNM.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `INPUTDTLOBJECTHULL_GCNM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlOtherMaintenanceAneka.xml` | `DATA-ADDRESS` | `INPUTDTLOTHERMAINTENANCEANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlOtherObjectAneka.xml` | `ASM-FW-GISFW-WORK` | `INPUTDTLOTHEROBJECTANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlOtherObjectAneka_FacIn.xml` | `ASM-FW-GISFW-WORK` | `INPUTDTLOTHEROBJECTANEKA_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / INPUTDTLOTHEROBJECTANEKA` |
| `InputDtlOtherObjectAneka_FacIn_ISUW.xml` | `ASM-FW-GISFW-WORK` | `INPUTDTLOTHEROBJECTANEKA_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / INPUTDTLOTHEROBJECTANEKA_FACIN` |
| `InputDtlPaymentFacRetro_FlowAct.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | `INPUTDTLPAYMENTFACRETRO_FLOWACT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlPaymentFacRetro_FlowAct_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | `INPUTDTLPAYMENTFACRETRO_FLOWACT_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY / INPUTDTLPAYMENTFACRETRO_FLOWACT` |
| `InputDtlPaymentFlow_FacIn.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | `INPUTDTLPAYMENTFLOW_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlPaymentFlow_FacInLife.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | `INPUTDTLPAYMENTFLOW_FACINLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY / INPUTDTLPAYMENTFLOW_FACIN` |
| `InputDtlPaymentFlow_FacInLife_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | `INPUTDTLPAYMENTFLOW_FACINLIFE_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY / INPUTDTLPAYMENTFLOW_FACINLIFE` |
| `InputDtlPaymentFlow_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | `INPUTDTLPAYMENTFLOW_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY / INPUTDTLPAYMENTFLOW_FACIN` |
| `InputDtlSpreadingLife_FacIn.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDTLSPREADINGLIFE_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLCOVERAGELIFE_FACIN` |
| `InputDtlSpreadingLife_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDTLSPREADINGLIFE_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLSPREADINGLIFE_FACIN` |
| `InputDtlSpreadingPA_FacIn.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDTLSPREADINGPA_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlSpreadingPA_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDTLSPREADINGPA_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLSPREADINGPA_FACIN` |
| `InputDtlSpreadingTravel_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDTLSPREADINGTRAVEL_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLSPREADINGTRAVEL_FACIN` |
| `InputDtlTrading.xml` | `ASM-FW-GISFW-DATA-TRADING` | `INPUTDTLTRADING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlTrading_IsUW.xml` | `ASM-FW-GISFW-DATA-TRADING` | `INPUTDTLTRADING_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TRADING / INPUTDTLTRADING` |
| `InputDtlWarranty.xml` | `ASM-FW-GISFW-WORK` | `INPUTDTLWARRANTY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputFEA.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-OFFERFEALIST` | `INPUTFEA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputFEA_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-OFFERFEALIST` | `INPUTFEA_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-OFFERFEALIST / INPUTFEA` |
| `InputHistoricalSurveyReport.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `INPUTHISTORICALSURVEYREPORT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputHistoricalSurveyReportUW.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `INPUTHISTORICALSURVEYREPORTUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-QUOTATION / INPUTHISTORICALSURVEYREPORT` |
| `InputKabupaten.xml` | `@BASECLASS` | `INPUTKABUPATEN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputKotaBranchName.xml` | `@BASECLASS` | `INPUTKOTABRANCHNAME` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / INPUTKOTA` |
| `InputKotaMOName.xml` | `@BASECLASS` | `INPUTKOTAMONAME` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / INPUTKOTABRANCHNAME` |
| `InputKotaProvince.xml` | `@BASECLASS` | `INPUTKOTAPROVINCE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / INPUTKOTA` |
| `InputLossRecord.xml` | `ASM-FW-GISFW-DATA-CAUSEOFLOSS` | `INPUTLOSSRECORD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-CAUSEOFLOSS / INPUTCAUSEOFLOSS_FACIN` |
| `InputNation.xml` | `@BASECLASS` | `INPUTNATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputObjFireGrid.xml` | `ASM-FW-GISFW-DATA-PROPERTY` | `INPUTOBJFIREGRID` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputObjFireGrid_GCNM.xml` | `ASM-FW-GISFW-DATA-PROPERTY` | `INPUTOBJFIREGRID_GCNM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTY / INPUTOBJFIREGRID` |
| `InputObjFireGrid_IsUW.xml` | `ASM-FW-GISFW-DATA-PROPERTY` | `INPUTOBJFIREGRID_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTY / INPUTOBJFIREGRID` |
| `InputProvince.xml` | `@BASECLASS` | `INPUTPROVINCE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputRW.xml` | `@BASECLASS` | `INPUTRW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputSpreadingLife_FacIn.xml` | `DATA-PARTY-PERSON` | `INPUTSPREADINGLIFE_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / INPUTCOVERAGELIFE_FACIN` |
| `InputSpreadingLife_FacIn_IsUW.xml` | `DATA-PARTY-PERSON` | `INPUTSPREADINGLIFE_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / INPUTSPREADINGLIFE_FACIN` |
| `InputSpreadingMBU_FacISUW.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `INPUTSPREADINGMBU_FACISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-VEHICLE / INPUTSPREADINGMBU_FACIN` |
| `InputSpreadingMBU_FacIn.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `INPUTSPREADINGMBU_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputSpreadingMarineCargo_FacIn.xml` | `ASM-FW-GISFW-DATA-CARGO` | `INPUTSPREADINGMARINECARGO_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputSpreadingMarineCargo_FacIn_ISUW.xml` | `ASM-FW-GISFW-DATA-CARGO` | `INPUTSPREADINGMARINECARGO_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-CARGO / INPUTSPREADINGMARINECARGO_FACIN` |
| `InputSpreadingPA_FacIn.xml` | `DATA-PARTY-PERSON` | `INPUTSPREADINGPA_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputSpreadingPA_FacIn_IsUW.xml` | `DATA-PARTY-PERSON` | `INPUTSPREADINGPA_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / INPUTSPREADINGPA_FACIN` |
| `InputSpreadingTravel_FacIn.xml` | `DATA-PARTY-PERSON` | `INPUTSPREADINGTRAVEL_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / INPUTSPREADINGPA_FACIN` |
| `InputSpreadingTravel_FacIn_IsUW.xml` | `DATA-PARTY-PERSON` | `INPUTSPREADINGTRAVEL_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / INPUTSPREADINGTRAVEL_FACIN` |
| `InstallmentList.xml` | `ASM-FW-GISFW-DATA-INSTALLMENT` | `INSTALLMENTLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Installments_ReadOnly.xml` | `ASM-FW-GISFW-DATA-TREATYININSTALLMENT` | `INSTALLMENTS_READONLY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYININSTALLMENT / INSTALLMENTS_REALISASI` |
| `InwardFacultative.xml` | `ASM-FW-GISFW-WORK` | `INWARDFACULTATIVE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InwardFacultative_IsUW.xml` | `ASM-FW-GISFW-WORK` | `INWARDFACULTATIVE_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / INWARDFACULTATIVE` |
| `LimitTreaty.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `LIMITTREATY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `LocationDetail_FlowAction.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `LOCATIONDETAIL_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ModalAccumulation_FacIn.xml` | `@BASECLASS` | `MODALACCUMULATION_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / MODALACCUMULATION` |
| `ObjItemGrid.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `OBJITEMGRID` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ObjectGridLocationHull_GISFW.xml` | `DATA-ADDRESS` | `OBJECTGRIDLOCATIONHULL_GISFW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-ADDRESS / OBJECTGRIDLOCATIONHULL_GCNM` |
| `ObjectGridLocation_FacIn.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `OBJECTGRIDLOCATION_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ObjectGridLocation_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `OBJECTGRIDLOCATION_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / OBJECTGRIDLOCATION_FACIN` |
| `ObjectGridLocation_GISFW.xml` | `DATA-ADDRESS` | `OBJECTGRIDLOCATION_GISFW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-ADDRESS / OBJECTGRIDLOCATION_GCNM` |
| `ObjectOtherSchedule.xml` | `ASM-FW-GISFW-WORK` | `OBJECTOTHERSCHEDULE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `OccupationAnekaGridGCNM.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `OCCUPATIONANEKAGRIDGCNM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OCCUPATION / OCCUPATIONANEKAGRID` |
| `OccupationAnekaGridView.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `OCCUPATIONANEKAGRIDVIEW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OCCUPATION / OCCUPATIONANEKAGRID` |
| `OccupationAnekaGrid_FacIn.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `OCCUPATIONANEKAGRID_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OCCUPATION / OCCUPATIONANEKAGRID` |
| `OccupationAnekaGrid_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `OCCUPATIONANEKAGRID_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OCCUPATION / OCCUPATIONANEKAGRID_FACIN` |
| `OccupationFacOut_FlowAction.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `OCCUPATIONFACOUT_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `OccupationFacOut_FlowAction_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `OCCUPATIONFACOUT_FLOWACTION_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / OCCUPATIONFACOUT_FLOWACTION` |
| `OccupationGrid.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `OCCUPATIONGRID` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `OccupationHullGrid_GCNM.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `OCCUPATIONHULLGRID_GCNM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `OccupationItemFacIn_FlowAction.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `OCCUPATIONITEMFACIN_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `OccupationItemFacIn_FlowAction_IsUW.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `OCCUPATIONITEMFACIN_FLOWACTION_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OCCUPATION / OCCUPATIONITEMFACIN_FLOWACTION` |
| `OfferFacOut.xml` | `ASM-FW-GISFW-WORK` | `OFFERFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `OfferFacinLifeConfirm.xml` | `ASM-FW-GISFW-WORK` | `OFFERFACINLIFECONFIRM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ACCEPTNOTIFICATIONEDM` |
| `PersonCoverage.xml` | `DATA-PARTY-PERSON` | `PERSONCOVERAGE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PersonGridDM.xml` | `DATA-PARTY-PERSON` | `PERSONGRIDDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / PERSONGRIDPA` |
| `PersonGridLife_FacIn.xml` | `DATA-PARTY-PERSON` | `PERSONGRIDLIFE_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PersonGridLife_FacIn_IsUW.xml` | `DATA-PARTY-PERSON` | `PERSONGRIDLIFE_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / PERSONGRIDLIFE_FACIN` |
| `PersonGridPA_FacIn.xml` | `DATA-PARTY-PERSON` | `PERSONGRIDPA_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / PERSONGRIDPA` |
| `PersonGridPA_FacIn_IsUW.xml` | `DATA-PARTY-PERSON` | `PERSONGRIDPA_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / PERSONGRIDPA_FACIN` |
| `PersonGridTravel.xml` | `DATA-PARTY-PERSON` | `PERSONGRIDTRAVEL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PersonGridTravel_GISFW.xml` | `DATA-PARTY-PERSON` | `PERSONGRIDTRAVEL_GISFW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / PERSONGRIDTRAVEL_GCNM` |
| `PersonGridTravel_IsUW.xml` | `DATA-PARTY-PERSON` | `PERSONGRIDTRAVEL_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / PERSONGRIDTRAVEL` |
| `PersonListDetail_FlowAct.xml` | `DATA-PARTY-PERSON` | `PERSONLISTDETAIL_FLOWACT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PersonPlanAndDeduct.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `PERSONPLANANDDEDUCT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PlanContains_UWRO.xml` | `ASM-FW-GISFW-DATA-PLAN` | `PLANCONTAINS_UWRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PlanListRO.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `PLANLISTRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PolicyTreatyInDeclineConfirm.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `POLICYTREATYINDECLINECONFIRM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PrintRISlip_FlowAction.xml` | `ASM-FW-GISFW-WORK` | `PRINTRISLIP_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / MEMOPLACINGFACOUT` |
| `PropertyItemCoverageCommisionFacIn_FlowAction.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `PROPERTYITEMCOVERAGECOMMISIONFACIN_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTYITEM / PROPERTYITEMCOVERAGESPREADINGFACIN_FLOWACTION` |
| `PropertyItemCoverageFacIn_FlowAction.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `PROPERTYITEMCOVERAGEFACIN_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTYITEM / PROPERTYITEMFACIN_FLOWACTION` |
| `PropertyItemCoverageFacIn_FlowAction_IsUW.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `PROPERTYITEMCOVERAGEFACIN_FLOWACTION_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTYITEM / PROPERTYITEMCOVERAGEFACIN_FLOWACTION` |
| `PropertyItemCoverageSpreadingFacIn_FlowAction.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `PROPERTYITEMCOVERAGESPREADINGFACIN_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTYITEM / PROPERTYITEMCOVERAGEFACIN_FLOWACTION` |
| `PropertyItemCoverageSpreadingFacIn_FlowAction_IsUW.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `PROPERTYITEMCOVERAGESPREADINGFACIN_FLOWACTION_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTYITEM / PROPERTYITEMCOVERAGESPREADINGFACIN_FLOWACTION` |
| `PropertyItemFacIn_FlowAction.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `PROPERTYITEMFACIN_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PropertyItemFacIn_FlowAction_IsUW.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `PROPERTYITEMFACIN_FLOWACTION_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTYITEM / PROPERTYITEMFACIN_FLOWACTION` |
| `PropertyItemFacOut_FlowAction.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `PROPERTYITEMFACOUT_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PropertyItemFacOut_FlowAction_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `PROPERTYITEMFACOUT_FLOWACTION_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / PROPERTYITEMFACOUT_FLOWACTION` |
| `PropertyItemSpreadingAneka_Flow.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `PROPERTYITEMSPREADINGANEKA_FLOW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PropertyItemSpreadingAneka_Flow_IsUW.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `PROPERTYITEMSPREADINGANEKA_FLOW_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-ANEKA / PROPERTYITEMSPREADINGANEKA_FLOW` |
| `Property_FlowAction.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `PROPERTY_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Property_FlowAction_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `PROPERTY_FLOWACTION_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / PROPERTY_FLOWACTION` |
| `ProtectCurrency.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `PROTECTCURRENCY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ProtectCurrencyCargo.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `PROTECTCURRENCYCARGO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ProtectCurrencyCargo_IsUW.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `PROTECTCURRENCYCARGO_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ProtectCurrency_IsUW.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `PROTECTCURRENCY_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-DEDUCTIBLE / PROTECTCURRENCY` |
| `ReasViewAttachment.xml` | `ASM-FW-GISFW-INT` | `REASVIEWATTACHMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `RejectNotificationFlow.xml` | `ASM-FW-GISFW-WORK` | `REJECTNOTIFICATIONFLOW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ReviseNotificationFlow.xml` | `ASM-FW-GISFW-WORK` | `REVISENOTIFICATIONFLOW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SelectionPlanRO.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `SELECTIONPLANRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OCCUPATION / SELECTIONPLAN` |
| `ShowCoverageFacOut.xml` | `ASM-FW-GISFW-DATA-PROPERTY` | `SHOWCOVERAGEFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowCoverageFacOut_FlowAction.xml` | `ASM-FW-GISFW-DATA-CARGO` | `SHOWCOVERAGEFACOUT_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowCoverageFacOut_FlowAction_IsUW.xml` | `ASM-FW-GISFW-DATA-CARGO` | `SHOWCOVERAGEFACOUT_FLOWACTION_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-CARGO / SHOWCOVERAGEFACOUT_FLOWACTION` |
| `ShowCoveragePAFacOutIsUW_FlowAction.xml` | `DATA-PARTY-PERSON` | `SHOWCOVERAGEPAFACOUTISUW_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / SHOWCOVERAGEPAFACOUT_FLOWACTION` |
| `ShowCoveragePAFacOut_FlowAction.xml` | `DATA-PARTY-PERSON` | `SHOWCOVERAGEPAFACOUT_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowCoverageTravelFacOut_FlowAction.xml` | `DATA-PARTY-PERSON` | `SHOWCOVERAGETRAVELFACOUT_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / SHOWCOVERAGEPAFACOUT_FLOWACTION` |
| `ShowCoverageTravelFacOut_FlowAction_IsUW.xml` | `DATA-PARTY-PERSON` | `SHOWCOVERAGETRAVELFACOUT_FLOWACTION_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / SHOWCOVERAGETRAVELFACOUT_FLOWACTION` |
| `ShowLetter.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `SHOWLETTER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowPolicyNoTreaty.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SHOWPOLICYNOTREATY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowPolis.xml` | `ASM-FW-GISFW-WORK` | `SHOWPOLIS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SpreadingCoverageAneka_FlowAction.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `SPREADINGCOVERAGEANEKA_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SpreadingCoverageAneka_FlowAction_IsUW.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `SPREADINGCOVERAGEANEKA_FLOWACTION_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-ANEKA / SPREADINGCOVERAGEANEKA_FLOWACTION` |
| `SpreadingItem.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SPREADINGITEM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / COVERAGEITEM` |
| `SpreadingItemAneka.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SPREADINGITEMANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SpreadingItemAneka_IsUW.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SPREADINGITEMANEKA_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / SPREADINGITEMANEKA` |
| `SpreadingItem_IsUW.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SPREADINGITEM_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / SPREADINGITEM` |
| `SpreadingObjectAneka_FlowAction.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `SPREADINGOBJECTANEKA_FLOWACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SpreadingObjectAneka_FlowAction_IsUW.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `SPREADINGOBJECTANEKA_FLOWACTION_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OCCUPATION / SPREADINGOBJECTANEKA_FLOWACTION` |
| `StartScreenFlowAuto.xml` | `ASM-FW-GISFW-WORK` | `STARTSCREENFLOWAUTO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TableInwardScale.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `TABLEINWARDSCALE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TableOfLimit.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `TABLEOFLIMIT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `UploadCSV_Aneka.xml` | `ASM-FW-GISFW-WORK` | `UPLOADCSV_ANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `UploadCSV_DataPolis.xml` | `@BASECLASS` | `UPLOADCSV_DATAPOLIS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / UPLOADCSV_LIFEPREMIUM` |
| `UploadCSV_Participant.xml` | `ASM-FW-GISFW-WORK` | `UPLOADCSV_PARTICIPANT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / UPLOADCSV_VEHICLE` |
| `UploadCSV_PolicyMember.xml` | `ASM-FW-GISFW-DATA-POLICY` | `UPLOADCSV_POLICYMEMBER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `UploadCSV_Vehicle.xml` | `ASM-FW-GISFW-WORK` | `UPLOADCSV_VEHICLE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `VehicleGridFacOut.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `VEHICLEGRIDFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-VEHICLE / VEHICLEGRID` |
| `VehicleGrid_FacIn.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `VEHICLEGRID_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-VEHICLE / VEHICLEGRID` |
| `VehicleGrid_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `VEHICLEGRID_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-VEHICLE / VEHICLEGRID_FACIN` |
| `ViewCoverageAneka.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `VIEWCOVERAGEANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-ANEKA / INPUTCOVERAGEANEKA` |
| `ViewCoverageFireGrid.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `VIEWCOVERAGEFIREGRID` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / COVERAGEFIREGRID` |
| `ViewCoveragePA.xml` | `DATA-PARTY-PERSON` | `VIEWCOVERAGEPA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / INPUTCOVERAGEPA` |
| `ViewCoverageTravel.xml` | `DATA-PARTY-PERSON` | `VIEWCOVERAGETRAVEL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / INPUTCOVERAGETRAVEL` |
| `ViewDetailPayment.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | `VIEWDETAILPAYMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewDtlAddCargo.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `VIEWDTLADDCARGO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLADDCARGO` |
| `ViewDtlAdditionalDetail.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `VIEWDTLADDITIONALDETAIL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLADDITIONALDETAIL` |
| `ViewDtlCargo.xml` | `ASM-FW-GISFW-DATA-CARGO` | `VIEWDTLCARGO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-CARGO / INPUTDTLCARGO` |
| `ViewDtlDeductible.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `VIEWDTLDEDUCTIBLE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewDtlObjectAneka.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `VIEWDTLOBJECTANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-ANEKA / INPUTDTLOBJECTANEKA` |
| `ViewDtlOtherObjectAneka.xml` | `ASM-FW-GISFW-WORK` | `VIEWDTLOTHEROBJECTANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / INPUTDTLOTHEROBJECTANEKA` |
| `ViewFireCoverage.xml` | `ASM-FW-GISFW-DATA-PROPERTY` | `VIEWFIRECOVERAGE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTY / FLOWACTCOVERAGELISTFIRE` |
| `ViewIndemnity_LA.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `VIEWINDEMNITY_LA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewObjectDetailGridLocation.xml` | `DATA-ADDRESS` | `VIEWOBJECTDETAILGRIDLOCATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-ADDRESS / OBJECTGRIDLOCATION` |
| `ViewObjectDetailOccupation.xml` | `DATA-ADDRESS` | `VIEWOBJECTDETAILOCCUPATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-ADDRESS / OBJECTDTLOCCUPATION` |
| `ViewObjectDtlAneka.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `VIEWOBJECTDTLANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OCCUPATION / OBJECTDTLANEKA` |
| `ViewOldFacOfferList.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `VIEWOLDFACOFFERLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-FACOFFER / ADDFACOFFERLIST_ISUW` |
| `ViewPaymentPremi.xml` | `ASM-FW-GISFW-WORK` | `VIEWPAYMENTPREMI` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewPaymentRetro.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `VIEWPAYMENTRETRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewPersonGridPA.xml` | `DATA-PARTY-PERSON` | `VIEWPERSONGRIDPA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / PERSONGRIDPA` |
| `ViewPolicyListRO.xml` | `ASM-FW-GISFW-DATA-POLICY` | `VIEWPOLICYLISTRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `confirmCreateAccumulation.xml` | `DATA-PARTY-PERSON` | `CONFIRMCREATEACCUMULATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Harness — 41 rule (`RULE-HTML-HARNESS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AccumulationRisk.xml` | `DATA-PORTAL` | `ACCUMULATIONRISK` | [terverifikasi] pyInclude=1 | — |
| `AdjustmentRiskAccumulation.xml` | `DATA-PORTAL` | `ADJUSTMENTRISKACCUMULATION` | [terverifikasi] pyInclude=1 | — |
| `CedingCedant.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `CEDINGCEDANT` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-DATA-QUOTATION / CEDINGCOMPANY` |
| `CedingCompany.xml` | `ASM-FW-GISFW-DATA-OFFERTREATYIN` | `CEDINGCOMPANY` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-DATA-QUOTATION / CEDINGCOMPANY` |
| `ChooseAccumulation_FacIn.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `CHOOSEACCUMULATION_FACIN` | [terverifikasi] pyInclude=3 | — |
| `ChooseClauseFire.xml` | `ASM-FW-GISFW-WORK` | `CHOOSECLAUSEFIRE` | [terverifikasi] pyInclude=1 | — |
| `ChooseDeductible.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `CHOOSEDEDUCTIBLE` | [terverifikasi] pyInclude=1 | — |
| `ChooseOccupation.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `CHOOSEOCCUPATION` | [terverifikasi] pyInclude=1 | — |
| `ChooseRiskAddress.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `CHOOSERISKADDRESS` | [terverifikasi] pyInclude=3 | — |
| `ChooseRiskLocation.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `CHOOSERISKLOCATION` | [terverifikasi] pyInclude=1 | — |
| `HistoricalSurveyReport.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `HISTORICALSURVEYREPORT` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-DATA-QUOTATION / HISTORICALSURVEYREPORT` |
| `HistoricalSurveyReportUW.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `HISTORICALSURVEYREPORTUW` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-DATA-POLICYTREATYIN / HISTORICALSURVEYREPORT` |
| `InputCauseOfDecline.xml` | `ASM-FW-GISFW-WORK` | `INPUTCAUSEOFDECLINE` | [terverifikasi] pyInclude=1 | — |
| `PrintRISlip.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `PRINTRISLIP` | [terverifikasi] pyInclude=6 | — |
| `PrintRISlips.xml` | `ASM-FW-GISFW-WORK` | `PRINTRISLIPS` | [terverifikasi] pyInclude=10 | — |
| `RiskAccumulationReport.xml` | `DATA-PORTAL` | `RISKACCUMULATIONREPORT` | [terverifikasi] pyInclude=1 | — |
| `SFAPortalOpportunities.xml` | `PEGACRM-PORTAL` | `SFAPORTALOPPORTUNITIES` | [terverifikasi] pyInclude=6 | — |
| `SOB.xml` | `ASM-FW-GISFW-DATA-OFFERTREATYIN` | `SOB` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-DATA-QUOTATION / SOB` |
| `SelectAgent.xml` | `ASM-FW-GISFW-INT-LLOYDAGENT` | `SELECTAGENT` | [terverifikasi] pyInclude=1 | — |
| `SelectCoverage.xml` | `ASM-FW-GISFW-INT-CONDITION` | `SELECTCOVERAGE` | [terverifikasi] pyInclude=1 | — |
| `SelectShip.xml` | `ASM-FW-GISFW-INT-SHIP` | `SELECTSHIP` | [terverifikasi] pyInclude=1 | — |
| `ShowCedingCoList.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SHOWCEDINGCOLIST` | [terverifikasi] pyInclude=1 | — |
| `ShowPolis.xml` | `ASM-FW-GISFW-WORK` | `SHOWPOLIS` | [terverifikasi] pyInclude=1 | — |
| `SummaryRiskAccumulation.xml` | `DATA-PORTAL` | `SUMMARYRISKACCUMULATION` | [terverifikasi] pyInclude=1 | — |
| `TabbedScreenFlow7.xml` | `ASM-FW-GISFW-WORK` | `TABBEDSCREENFLOW7` | [terverifikasi] pyInclude=34 | `WORK- / TABBEDSCREENFLOW7` |
| `TotalAccumulationDtl.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `TOTALACCUMULATIONDTL` | [terverifikasi] pyInclude=1 | — |
| `TotalAccumulation_FacIn.xml` | `ASM-FW-GISFW-WORK` | `TOTALACCUMULATION_FACIN` | [terverifikasi] pyInclude=1 | — |
| `ViewCSVResult_B2BHost.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `VIEWCSVRESULT_B2BHOST` | [terverifikasi] pyInclude=4 | `ASM-FW-GISFW-INT-TABLE_B2B / VIEWCSVRESULT_B2BHOST` |
| `ViewClaimList.xml` | `ASM-FW-GISFW-WORK` | `VIEWCLAIMLIST` | [terverifikasi] pyInclude=2 | — |
| `ViewCoverageFacOut.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `VIEWCOVERAGEFACOUT` | [terverifikasi] pyInclude=1 | — |
| `ViewCoverageFacOut_isUW.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `VIEWCOVERAGEFACOUT_ISUW` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-DATA-PROPERTYITEM / VIEWCOVERAGEFACOUT` |
| `ViewDataOfferFacIn.xml` | `ASM-FW-GISFW-WORK` | `VIEWDATAOFFERFACIN` | [terverifikasi] pyInclude=22 | — |
| `ViewFacretro.xml` | `ASM-FW-GISFW-WORK` | `VIEWFACRETRO` | [terverifikasi] pyInclude=7 | — |
| `ViewFollowingNB.xml` | `ASM-FW-GISFW-WORK` | `VIEWFOLLOWINGNB` | [terverifikasi] pyInclude=8 | — |
| `ViewLetter.xml` | `ASM-FW-GISFW-WORK` | `VIEWLETTER` | [terverifikasi] pyInclude=1 | — |
| `ViewOfferStatus.xml` | `ASM-FW-GISFW-WORK` | `VIEWOFFERSTATUS` | [terverifikasi] pyInclude=1 | — |
| `ViewOldDeductible.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `VIEWOLDDEDUCTIBLE` | [terverifikasi] pyInclude=1 | — |
| `ViewOldEndorsement.xml` | `ASM-FW-GISFW-WORK` | `VIEWOLDENDORSEMENT` | [terverifikasi] pyInclude=22 | — |
| `ViewOldOffer.xml` | `ASM-FW-GISFW-WORK` | `VIEWOLDOFFER` | [terverifikasi] pyInclude=1 | — |
| `ViewPictureList.xml` | `ASM-FW-GISFW-INT` | `VIEWPICTURELIST` | [terverifikasi] pyInclude=1 | — |
| `ViewPolis.xml` | `ASM-FW-GISFW-WORK` | `VIEWPOLIS` | [terverifikasi] pyInclude=54 | — |

### RDBList — 217 rule (`RULE-CONNECT-SQL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AttachmentLife.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `ATTACHMENTLIFE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CATEGORY_ATTACH_REAS | — |
| `BrowseDataViewPolis1.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `BROWSEDATAVIEWPOLIS1` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=FACINPRODUCTION,POOLDATA.FACINPRODUCTION | `ASM-FW-GISFW-INT-POLICYJSON / ASM!BROWSEDATAVIEWPOLIS` |
| `BrowseLifeRateRetro_SQL.xml` | `ASM-FW-GISFW-INT-M_RATE_LIFE` | `BROWSELIFERATERETRO_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=RATE_LIFE | `ASM-FW-GISFW-INT-M_RATE_LIFE / ASM!BROWSELIFERATE_SQL` |
| `BrowseLifeRate_SQL.xml` | `ASM-FW-GISFW-INT-M_RATE_LIFE` | `BROWSELIFERATE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=RATE_LIFE | `ASM-FW-GISFW-INT-RATE_LIFE_SUMMARY / ASM!BROWSELIFERATE_SQL` |
| `BrowseMaxValueTreatyType.xml` | `ASM-FW-GISFW-INT-TREATYBUSINESS` | `BROWSEMAXVALUETREATYTYPE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_TREATYBUSINESS | `ASM-FW-GISFW-DATA-SPREADINGRISK / ASM!BROWSEMAXVALUETREATYTYPE` |
| `BrowseMaxValueTreatyType2.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSEMAXVALUETREATYTYPE2` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=PROPORTIONALARRG | `ASM-FW-GISFW-DATA-SPREADINGRISK / ASM!BROWSEMAXVALUETREATYTYPE2` |
| `BrowseRW2_SQL.xml` | `ASM-FW-GISFW-INT-RW` | `BROWSERW2_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.RW | `ASM-FW-GISFW-INT-RW / ASM!BROWSERW_SQL` |
| `BrowseRW_SQL.xml` | `ASM-FW-GISFW-INT-RW` | `BROWSERW_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.RW | — |
| `BrowseTreatyIn.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `BROWSETREATYIN` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_TREATY_IN,POOLDATA.M_TREATY_IN_EDM | `ASM-FW-GISFW-INT-PRODUCT_LIFE / ASM!BROWSEUNDERWRITINGLIST` |
| `BrowseZipCodeAndCzone.xml` | `ASM-FW-GISFW-INT-RW` | `BROWSEZIPCODEANDCZONE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=RW | — |
| `CategoryAttach_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `CATEGORYATTACH_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CATEGORY_ATTACH_REAS | — |
| `CekDuplicateLocationOffer_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `CEKDUPLICATELOCATIONOFFER_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=FACINOFFER,WITH,WORDLIST | — |
| `CekFacin_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `CEKFACIN_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.FACINPRODUCTION | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETNOPOLIS_SQL` |
| `CekFacoutProd_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `CEKFACOUTPROD_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=FACOUTPRODUCTION | `ASM-FW-GISFW-INT-POLICYJSON / RNM!INSERTTREATYPROD_SQL` |
| `CekSTSKonversiJson.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `CEKSTSKONVERSIJSON` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.JSON_POLIS | — |
| `CheckZipCode_SQL.xml` | `ASM-FW-GISFW-INT-RW` | `CHECKZIPCODE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=RW | — |
| `ConvertBusinessField.xml` | `ASM-FW-GISFW-INT-BUSINESSFIELD` | `CONVERTBUSINESSFIELD` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_BUSINESSFIELD | `ASM-FW-GISFW-INT-SERVICESFA-CUSTOMERCOMPANY / ASM!CONVERTBUSINESSFIELD` |
| `ConvertCurrency.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `CONVERTCURRENCY` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_CURRENCY | — |
| `ConvertNationality.xml` | `ASM-FW-GISFW-INT-NATION` | `CONVERTNATIONALITY` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_COUNTRY | `ASM-FW-GISFW-INT-REINSURANCETYPE / ASM!CARINOTESPREADINGRISK` |
| `CurrencyStandard.xml` | `ASM-FW-GISFW-INT-CURRENCYSTANDARD` | `CURRENCYSTANDARD` | [terverifikasi] sqlKind=QUERY; procs=POOLDATA.GETCURRENCYSTANDARD; sqlOps=SELECT | — |
| `DeleteDataProduction.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `DELETEDATAPRODUCTION` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_DELETE_ERROR_KONVERSI | `ASM-FW-GISFW-INT-POLICYJSON / ASM!DELETEJSONPOLIS` |
| `DeleteStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `DELETESTORAGE_SQL` | [terverifikasi] sqlKind=QUERY | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GETLINKSTORAGE_SQL` |
| `FetchTreatyGroupOLDID.xml` | `ASM-FW-GISFW-INT-TREATYGROUP` | `FETCHTREATYGROUPOLDID` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TREATYGROUP | — |
| `FilterAttachmentByPositionEDM_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `FILTERATTACHMENTBYPOSITIONEDM_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CATEGORY_ATTACH_REAS | `ASM-FW-GISFW-INT-OFFERJSON / ASM!FILTERATTACHMENTBYPOSITION_SQL` |
| `ForInputCurrencyAdjBackup_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `FORINPUTCURRENCYADJBACKUP_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=TREATYPRODUCTION_BACKUP | — |
| `ForInputCurrencyAdj_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `FORINPUTCURRENCYADJ_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=FACINPRODUCTION | — |
| `GETTanggalClosing_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETTANGGALCLOSING_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TANGGAL_CLOSING | — |
| `GenerateImageID_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GENERATEIMAGEID_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GenerateNoPolicy.xml` | `ASM-FW-GISFW-INT-POLISTREATYIN` | `GENERATENOPOLICY` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GenerateOurRefFacOut_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `GENERATEOURREFFACOUT_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GenerateRISlipNumber.xml` | `ASM-FW-GISFW-DATA-PRINTRISLIP` | `GENERATERISLIPNUMBER` | [terverifikasi] sqlKind=QUERY; procs=POOLDATA.GENERATE_FACRETRO_NO; sqlOps=SELECT | — |
| `Get5YearAgo.xml` | `ASM-FW-GISFW-WORK` | `GET5YEARAGO` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | `ASM-FW-GISFW-WORK / ASM!GETDETAILSERVISASM` |
| `GetAccumulationByNote_SQL.xml` | `ASM-FW-GISFW-INT-ACCUMULATION` | `GETACCUMULATIONBYNOTE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=ACCUMULATION | `ASM-FW-GISFW-INT-ACCUMULATION / ASM!GETACCUMULATION_SQL` |
| `GetAccumulationDistrict_SQL.xml` | `ASM-FW-GISFW-INT-ACCUMULATION` | `GETACCUMULATIONDISTRICT_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=ACCUMULATION,RW | `ASM-FW-GISFW-INT-ACCUMULATION / ASM!GETACCUMULATIONPROVINCE_SQL` |
| `GetAccumulationProvince_SQL.xml` | `ASM-FW-GISFW-INT-ACCUMULATION` | `GETACCUMULATIONPROVINCE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=ACCUMULATION,RW | `ASM-FW-GISFW-INT-ACCUMULATION / ASM!GETACCUMULATION_SQL` |
| `GetAccumulationZipcode_SQL.xml` | `ASM-FW-GISFW-INT-ACCUMULATION` | `GETACCUMULATIONZIPCODE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CZONE,RW | `ASM-FW-GISFW-INT-ACCUMULATION / ASM!GETACCUMULATIONDISTRICT_SQL` |
| `GetAccumulation_Sql.xml` | `ASM-FW-GISFW-INT-ACCUMULATION_LIFE` | `GETACCUMULATION_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_ACCUMULATION_LIFE | — |
| `GetAksepBanding_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETAKSEPBANDING_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_LIMIT_PROPERTYY,HISTORYAKSEPTASIPEGA | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETVIEWBANDING_SQL` |
| `GetAllCurrency.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `GETALLCURRENCY` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CURRENCY | — |
| `GetAppName_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETAPPNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.T_FOLDER_IMAGE | — |
| `GetBreakDownSpread_SQL.xml` | `ASM-FW-GISFW-WORK` | `GETBREAKDOWNSPREAD_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.PROPORTIONALARRG | — |
| `GetCLaimList_SQL.xml` | `ASM-FW-GISFW-WORK` | `GETCLAIMLIST_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=DETAIL_INVOICE | — |
| `GetCategoryOccupation.xml` | `ASM-FW-GISFW-INT-OCCUPATION` | `GETCATEGORYOCCUPATION` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=OCCUPATION | `ASM-FW-GISFW-INT-OCCUPATION / ASM!GETOCCUPATIONNAME` |
| `GetCauseOfDecline.xml` | `ASM-FW-GISFW-DATA-CAUSEOFDECLINE` | `GETCAUSEOFDECLINE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.CAUSEOFDECLINE | `ASM-FW-GISFW-INT-OFFERJSON / ASM!GETCOPYNB` |
| `GetCauseOfDecline2.xml` | `ASM-FW-GISFW-DATA-CAUSEOFDECLINE` | `GETCAUSEOFDECLINE2` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.CAUSEOFDECLINE | `ASM-FW-GISFW-DATA-CAUSEOFDECLINE / RNM!GETCAUSEOFDECLINE` |
| `GetCountAccumulation.xml` | `ASM-FW-GISFW-WORK` | `GETCOUNTACCUMULATION` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.FACINPRODUCTION,POOLDATA.FACINOFFER | `ASM-FW-GISFW-WORK / ASM!GETBRANCHCODE` |
| `GetCountClaim.xml` | `ASSIGN-WORKLIST` | `GETCOUNTCLAIM` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=DATAPEGA.PC_ASM_FW_GCNMFW_WORK | — |
| `GetCoverageFacin.xml` | `ASM-FW-GISFW-INT-COVERAGE_FACIN` | `GETCOVERAGEFACIN` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=COVERAGE_FACIN | — |
| `GetCurrency.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `GETCURRENCY` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CURRENCY | `ASM-FW-GISFW-INT-CURRENCY / ASM!UPDATEMASTERCURRENCY` |
| `GetCurrencyIDByName.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `GETCURRENCYIDBYNAME` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.CURRENCY | `ASM-FW-GISFW-INT-CURRENCY / ASM!GETDATACURRENCYBYNAME_SQL` |
| `GetCurrencySymbol.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `GETCURRENCYSYMBOL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CURRENCY | `ASM-FW-GISFW-INT-CURRENCY / ASM!GETCURRENCY` |
| `GetCurrencyToIDR_SQL.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `GETCURRENCYTOIDR_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYEXCHANGEYEARLY | — |
| `GetCurrentDate.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `GETCURRENTDATE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GetDataAgentByNameNonLife_SQL.xml` | `ASM-FW-GISFW-INT-AGENT` | `GETDATAAGENTBYNAMENONLIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=AGENT | `ASM-FW-GISFW-INT-AGENT / ASM!GETDATAAGENTBYNAME_SQL` |
| `GetDataByID_SQL.xml` | `ASM-FW-GISFW-INT-REINSURANCETYPE` | `GETDATABYID_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_REINSURANCETYPE | — |
| `GetDataClientByName1_SQL.xml` | `ASM-FW-GISFW-INT-CLIENT` | `GETDATACLIENTBYNAME1_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CLIENT | `ASM-FW-GISFW-INT-CLIENT / ASM!GETDATACLIENTBYNAME_SQL` |
| `GetDataCurrencyByName_SQL.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `GETDATACURRENCYBYNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CURRENCY | — |
| `GetDataDoubleCase.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETDATADOUBLECASE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=FACINPRODUCTION | — |
| `GetDataKlaim_SQL.xml` | `ASM-FW-GISFW-WORK` | `GETDATAKLAIM_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.DATAKLAIM | — |
| `GetDataMarketingByName_SQL.xml` | `ASM-FW-GISFW-INT-MARKETINGOFFICER` | `GETDATAMARKETINGBYNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=MARKETINGOFFICER | — |
| `GetDataMarketing_SQL.xml` | `ASM-FW-GISFW-INT-MARKETINGOFFICER` | `GETDATAMARKETING_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=MARKETINGOFFICER | — |
| `GetDataObjectItem.xml` | `ASM-FW-GISFW-INT-OBJECTITEMTYPE` | `GETDATAOBJECTITEM` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=OBJECTITEMTYPE | — |
| `GetDataSyariah_SQL.xml` | `ASM-FW-GISFW-INT-SPREADSYARIAH` | `GETDATASYARIAH_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=SPREADSYARIAH,TREATYBUSINESS,TREATYYEAR | — |
| `GetEDMOldIDPEGA.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `GETEDMOLDIDPEGA` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS,FACOUTPRODUCTION | — |
| `GetFlagReject_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETFLAGREJECT_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=HISTORYAKSEPTASIPEGA | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETVIEWBANDING_SQL` |
| `GetGroupName_SQL.xml` | `ASM-FW-GISFW-WORK` | `GETGROUPNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.CLIENT | `ASM-FW-GISFW-WORK / ASM!GETNPWPBYSOB_SQL` |
| `GetHistoryAccPega_SQL.xml` | `ASM-FW-GISFW-WORK` | `GETHISTORYACCPEGA_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=DATAPEGA.PC_HISTORY_ASM_FW_GISFW_WORK | — |
| `GetIDFollowing_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETIDFOLLOWING_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_FOLLOWING | — |
| `GetIDPegaJsonFollowing_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETIDPEGAJSONFOLLOWING_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.JSON_FOLLOWING | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETJSONDATAFOLLOWING_SQL` |
| `GetIDPega_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETIDPEGA_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | — |
| `GetInsuredID.xml` | `ASM-FW-GISFW-INT-CLIENT` | `GETINSUREDID` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.CLIENT | — |
| `GetJsonDataFollowingView_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETJSONDATAFOLLOWINGVIEW_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETJSONDATAFOLLOWING_SQL` |
| `GetJsonDataFollowing_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETJSONDATAFOLLOWING_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_FOLLOWING | — |
| `GetKodeProdLife_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETKODEPRODLIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.KODE_PRODUKSI | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETKODEPRODNONLIFE_SQL` |
| `GetKodeProdNonLife_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETKODEPRODNONLIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.KODE_PRODUKSI | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETSEQUENCENUMBER_SQL` |
| `GetKurs.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | `GETKURS` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=LST_KURS_STANDARD@ASMD.SINARMAS.CO.ID | — |
| `GetKursLimitByName_SQL.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `GETKURSLIMITBYNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TREATYEXCHANGEYEARLY | `ASM-FW-GISFW-INT-CURRENCY / ASM!GETKURSLIMITSPREADING_SQL` |
| `GetKursLimitSpreading_SQL.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `GETKURSLIMITSPREADING_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYEXCHANGEYEARLY | — |
| `GetLastPPNCheckEDM.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENT` | `GETLASTPPNCHECKEDM` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.JSON_POLIS | — |
| `GetLimitAccEngineeringBanding_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLIMITACCENGINEERINGBANDING_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_LIMIT_ENGINEERINGG | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETLIMITACCENGINEERINGUW_SQL` |
| `GetLimitAccEngineeringJUWA_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLIMITACCENGINEERINGJUWA_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_LIMIT_ENGINEERINGG | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETLIMITACCENGINEERINGUW_SQL` |
| `GetLimitAccEngineeringUW_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLIMITACCENGINEERINGUW_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_LIMIT_ENGINEERINGG | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETLIMITAKSEPTASI_SQL` |
| `GetLimitAkseptasiBanding_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLIMITAKSEPTASIBANDING_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_LIMIT_PROPERTYY | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETLIMITAKSEPTASI_SQL` |
| `GetLimitAkseptasiBond_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLIMITAKSEPTASIBOND_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_LIMIT_FINANCIALINS | — |
| `GetLimitAkseptasiJUWA_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLIMITAKSEPTASIJUWA_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_LIMIT_PROPERTYY | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETLIMITAKSEPTASI_SQL` |
| `GetLimitAkseptasiKreditCL_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLIMITAKSEPTASIKREDITCL_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_LIMIT_FINANCIALINS | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETLIMITAKSEPTASIBOND_SQL` |
| `GetLimitAkseptasiKreditNCL_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLIMITAKSEPTASIKREDITNCL_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_LIMIT_FINANCIALINS | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETLIMITAKSEPTASIBOND_SQL` |
| `GetLimitAkseptasiLife_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLIMITAKSEPTASILIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_LIMIT_LIFE | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETLIMITAKSEPTASI_SQL` |
| `GetLimitAkseptasiNonFireBanding_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLIMITAKSEPTASINONFIREBANDING_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_LIMIT_NONPROPANDENGG | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETLIMITAKSEPTASINONFIRE_SQL` |
| `GetLimitAkseptasiNonFireJUWA_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLIMITAKSEPTASINONFIREJUWA_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_LIMIT_NONPROPANDENGG | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETLIMITAKSEPTASINONFIRE_SQL` |
| `GetLimitAkseptasiNonFire_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLIMITAKSEPTASINONFIRE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_LIMIT_NONPROPANDENGG | — |
| `GetLimitAkseptasiNonPreferBanding_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLIMITAKSEPTASINONPREFERBANDING_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_LIMIT_PROPERTY_NON_PREFERREDD | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETLIMITAKSEPTASINONPREFER_SQL` |
| `GetLimitAkseptasiNonPreferJUWA_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLIMITAKSEPTASINONPREFERJUWA_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_LIMIT_PROPERTY_NON_PREFERREDD | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETLIMITAKSEPTASINONPREFER_SQL` |
| `GetLimitAkseptasiNonPrefer_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLIMITAKSEPTASINONPREFER_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_LIMIT_PROPERTY_NON_PREFERREDD | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETLIMITAKSEPTASI_SQL` |
| `GetLimitAkseptasiPreferedCommJUWA_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLIMITAKSEPTASIPREFEREDCOMMJUWA_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETLIMITAKSEPTASIPREFEREDCOMM_SQL` |
| `GetLimitAkseptasiPreferedComm_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLIMITAKSEPTASIPREFEREDCOMM_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETLIMITAKSEPTASINONPREFER_SQL` |
| `GetLimitAkseptasi_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLIMITAKSEPTASI_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_LIMIT_PROPERTYY | — |
| `GetLinkStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETLINKSTORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GETTOKENSTORAGE_SQL` |
| `GetMKTandLeader_SQL.xml` | `ASM-FW-GISFW-INT-MARKETINGOFFICER` | `GETMKTANDLEADER_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=MARKETINGOFFICER | — |
| `GetMOID_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETMOID_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=MARKETINGOFFICER | — |
| `GetMasterBenefit.xml` | `ASM-FW-GISFW-INT-VJ_PKG_BENEFIT` | `GETMASTERBENEFIT` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.VJ_PKG_BENEFIT | `ASM-FW-GISFW-INT-VJ_PACKAGE_AGE_KLAUSUL / ASM!GETMASTERKLAUSULAGE` |
| `GetMasterJPlanBenefit.xml` | `ASM-FW-GISFW-INT-VJ_BENEFIT_PROPERTY` | `GETMASTERJPLANBENEFIT` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.VJ_BENEFIT_PROPERTY | `ASM-FW-GISFW-INT-VIEW_BENEFIT_PROPERTY / ASM!GETMASTERJPLANBENEFIT` |
| `GetMasterKlausulAge.xml` | `ASM-FW-GISFW-INT-VJ_PACKAGE_AGE_KLAUSUL_DM` | `GETMASTERKLAUSULAGE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.VJ_PACKAGE_AGE_KLAUSUL_DM | `ASM-FW-GISFW-INT-VJ_PACKAGE_AGE_KLAUSUL / ASM!GETMASTERKLAUSULAGE` |
| `GetMasterPolicyJSON.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETMASTERPOLICYJSON` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | `ASM-FW-GISFW-INT-PLANTRAVEL / ASM!GETMASTERPOLICYJSON` |
| `GetNPWPbyCedingCo_SQL.xml` | `ASM-FW-GISFW-WORK` | `GETNPWPBYCEDINGCO_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.CLIENT,POOLDATA.AGENT | `ASM-FW-GISFW-WORK / ASM!GETNPWPBYSOB_SQL` |
| `GetNPWPbySOB_SQL.xml` | `ASM-FW-GISFW-WORK` | `GETNPWPBYSOB_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.CLIENT,POOLDATA.AGENT | `ASM-FW-GISFW-WORK / ASM!GETBRANCHCODE` |
| `GetNoEndors_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETNOENDORS_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | — |
| `GetNoEndorsement_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETNOENDORSEMENT_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | — |
| `GetNoOffer_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETNOOFFER_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.JSON_OFFER | — |
| `GetNoRangkaMesinDiff.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETNORANGKAMESINDIFF` | [terverifikasi] sqlKind=QUERY; procs=MBU.F_CEK_HURUF; sqlOps=SELECT | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GENERATECOUNTER` |
| `GetNopolis1_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETNOPOLIS1_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | — |
| `GetNopolis_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETNOPOLIS_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | — |
| `GetOPFacOut_Sql.xml` | `ASM-FW-GISFW-WORK` | `GETOPFACOUT_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=DATAPEGA.PC_HISTORY_ASM_FW_GISFW_WORK | — |
| `GetObjectItembyName_SQL.xml` | `ASM-FW-GISFW-INT-V_JN_OBJ_ITEM` | `GETOBJECTITEMBYNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=V_JN_OBJ_ITEM | `ASM-FW-GISFW-INT-V_JN_OBJ_ITEM / ASM!GETOBJECTITEMBYID_SQL` |
| `GetOldDataEDM_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETOLDDATAEDM_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | — |
| `GetOldIDBusiness_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `GETOLDIDBUSINESS_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=BUSINESS | — |
| `GetPolicyNoByCaseId.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETPOLICYNOBYCASEID` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | `ASM-FW-GISFW-INT-POLICYJSON / ASM!BROWSECLIENTIDBYNOPOLICY` |
| `GetProdKeEDM_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `GETPRODKEEDM_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | — |
| `GetProdKeOldData_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `GETPRODKEOLDDATA_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.JSON_POLIS | `ASM-FW-GISFW-INT-OFFERJSON / ASM!GETEDMOLDDATA_SQL` |
| `GetPromptAI_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETPROMPTAI_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_PROMPT_AI | — |
| `GetProrateEdm_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETPRORATEEDM_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | — |
| `GetRiskExposureOccupationFire.xml` | `ASM-FW-GISFW-INT-OCCUPATION` | `GETRISKEXPOSUREOCCUPATIONFIRE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.OCCUPATION | — |
| `GetSQLDate.xml` | `ASM-FW-GISFW-WORK` | `GETSQLDATE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GetSequenceNumber_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETSEQUENCENUMBER_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETTANGGALCLOSING_SQL` |
| `GetStsKonversiPolicy_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETSTSKONVERSIPOLICY_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | — |
| `GetSummaryRiskAccumPolis_Sql.xml` | `ASM-FW-GISFW-INT-ACCUMULATION` | `GETSUMMARYRISKACCUMPOLIS_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=ACCUMULATION,POOLDATA.REINSURANCETYPE,POOLDATA.JSON_POLIS | `ASM-FW-GISFW-INT-ACCUMULATION / ASM!GETSUMMARYRISKACCUMULATION_SQL` |
| `GetSummaryRiskAccumulation_Sql.xml` | `ASM-FW-GISFW-INT-ACCUMULATION` | `GETSUMMARYRISKACCUMULATION_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=ACCUMULATION,POOLDATA.REINSURANCETYPE,CURRENCY,FACINPRODUCTION | — |
| `GetTglInputOffer_sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETTGLINPUTOFFER_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.JSON_OFFER | — |
| `GetTgl_InputJsonPolis_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETTGL_INPUTJSONPOLIS_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | — |
| `GetTokenStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETTOKENSTORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.GET_TOKEN_STORAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |
| `GetTotalAccumulation_sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETTOTALACCUMULATION_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.REINSURANCETYPE,FACINPRODUCTION | — |
| `GetTreatyName.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `GETTREATYNAME` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=REINSURANCETYPE | `ASM-FW-GISFW-INT-PROPORTIONALARRG / ASM!GETTREATYNAME_SQL` |
| `GetTreatyName_SQL.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `GETTREATYNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.REINSURANCETYPE,PROPORTIONALARRG,TREATYBUSINESS,TREATYCONTRACT | — |
| `GetTreatyName_SQL2.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `GETTREATYNAME_SQL2` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.REINSURANCETYPE,PROPORTIONALARRG,TREATYBUSINESS | `ASM-FW-GISFW-INT-PROPORTIONALARRG / ASM!GETTREATYNAME_SQL` |
| `INSERTCONORGJSON_MCLIENT.xml` | `ASM-FW-GISFW-INT-MCLIENT` | `INSERTCONORGJSON_MCLIENT` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.RDBINSERTCLIENT | `ASM-FW-GISFW-INT-MCLIENT / ASM!UPDATECONORGJSON_MCLIENT` |
| `INSERTERRORFACIN_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTERRORFACIN_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.ERRORFACINPROD | — |
| `INSERTJSON_JSONPOLISEDM_FACIN.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTJSON_JSONPOLISEDM_FACIN` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.INSERTJSONPOLIS | `ASM-FW-GISFW-INT-POLICYJSON / ASM!INSERTJSON_JSONPOLIS_FACIN` |
| `INSERTJSON_JSONPOLISMONITORING_FACIN.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTJSON_JSONPOLISMONITORING_FACIN` | [terverifikasi] sqlKind=PLSQL; procs=DBMS_LOB.CREATETEMPORARY,POOLDATA.INSERTJSONPOLISMONITORING | `ASM-FW-GISFW-INT-POLICYJSON / ASM!INSERTJSON_JSONPOLIS_FACIN` |
| `INSERTJSON_JSONPOLIS_FACIN.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTJSON_JSONPOLIS_FACIN` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.INSERTJSONPOLIS | `ASM-FW-GISFW-INT-POLICYJSON / ASM!INSERTJSON_JSONPOLIS` |
| `InsertCedingProduction_SQL.xml` | `ASM-FW-GISFW-WORK` | `INSERTCEDINGPRODUCTION_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.INSERTUPDATECEDINGPRODUCTION | — |
| `InsertDataFacinPerformance_SQL.xml` | `ASM-FW-GISFW-INT-JSONOFFER` | `INSERTDATAFACINPERFORMANCE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=INSERT; tables=POOLDATA.FACINPERFORMANCE | — |
| `InsertFacinLifeMonthly_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTFACINLIFEMONTHLY_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.FACINLIFE | `ASM-FW-GISFW-INT-POLICYJSON / ASM!INSERTFACINLIFE_SQL` |
| `InsertFacinLife_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTFACINLIFE_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.FACINLIFE | — |
| `InsertFacinOfferLife_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTFACINOFFERLIFE_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.FACINOFFERLIFE | — |
| `InsertFacinProductionBackup_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTFACINPRODUCTIONBACKUP_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.TREATYPRODUCTION_BACKUP | — |
| `InsertHistoryAkseptasiPega_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTHISTORYAKSEPTASIPEGA_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=HISTORYAKSEPTASIPEGA | — |
| `InsertIntoFacinSPreadLifeMonthly_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTINTOFACINSPREADLIFEMONTHLY_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.FACINSPREADLIFE | `ASM-FW-GISFW-INT-POLICYJSON / ASM!INSERTINTOFACINSPREADLIFE_SQL` |
| `InsertIntoFacinSPreadLife_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTINTOFACINSPREADLIFE_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.FACINSPREADLIFE | — |
| `InsertIntoJsonError_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTINTOJSONERROR_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=INSERT; tables=JSON_POLIS_ERROR | — |
| `InsertJsonFollowing_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTJSONFOLLOWING_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=JSON_FOLLOWING | — |
| `InsertLogServiceProd.xml` | `ASM-FW-GISFW-WORK` | `INSERTLOGSERVICEPROD` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.MONITORING_PROD_LOG | `ASM-FW-GCNMFW-WORK / RNM!INSERTLOGSERVICECLAIM` |
| `InsertOfferProduction_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTOFFERPRODUCTION_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.FACINOFFER | — |
| `InsertRiskAndLossProfile_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTRISKANDLOSSPROFILE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=INSERT; tables=M_RISK_LOSS_PROFILE | `ASM-FW-GISFW-INT-POLICYJSON / ASM!INSERTRISK&amp;LOSSPROFILE_SQL` |
| `InsertTreatyProd_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTTREATYPROD_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=FACOUTPRODUCTION | `ASM-FW-GISFW-INT-POLICYJSON / ASM!INSERTTREATYPRODUCTION_SQL` |
| `InsertTreatyProduction_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTTREATYPRODUCTION_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=FACINPRODUCTION | — |
| `InsertViewSuggest_SQL.xml` | `ASM-FW-GISFW-WORK` | `INSERTVIEWSUGGEST_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.HISTORYAKSEPTASIPRODUCTION | — |
| `Insert_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERT_T_STORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GENERATEIMAGEID_SQL` |
| `MachingDataFacin_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `MACHINGDATAFACIN_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.FACINFORBACKUP | — |
| `RateEarthquake_Sql.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `RATEEARTHQUAKE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.ZONEOJK | `ASM-FW-GISFW-INT-OFFERJSON / ASM!RATEFLEXAS_SQL` |
| `RateFlexas_Sql.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `RATEFLEXAS_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=OCCUPATION | `ASM-FW-GISFW-INT-OFFERJSON / ASM!SEARCHRATEPOLIS_SQL` |
| `RetrieveArgumentNumberClauseSQL.xml` | `ASM-FW-GISFW-INT-CLAUSE` | `RETRIEVEARGUMENTNUMBERCLAUSESQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_ARGCLAUSEFIRE,M_CLAUSE | — |
| `RetrieveDataAccumulationLife_SQL.xml` | `ASM-FW-GISFW-INT-ACCUMULATION_LIFE` | `RETRIEVEDATAACCUMULATIONLIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_ACCUMULATION_LIFE | `ASM-FW-GISFW-INT-ACCUMULATION_LIFE / ASM!GETACCUMULATION_SQL` |
| `SaveNewAccumulation_SQL.xml` | `ASM-FW-GISFW-INT-ACCUMULATION_LIFE` | `SAVENEWACCUMULATION_SQL` | [terverifikasi] sqlKind=PLSQL; procs=DBMS_LOB.CREATETEMPORARY,POOLDATA.PEGA_M_ACCUMULATION_LIFE | `ASM-FW-GISFW-INT-ACCUMULATION_LIFE / ASM!RETRIEVEDATAACCUMULATIONLIFE_SQL` |
| `SaveOfferJson_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `SAVEOFFERJSON_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_M_JSON_OFFER | — |
| `SavePolisTreatyIn_SQL.xml` | `ASM-FW-GISFW-INT-POLISTREATYIN` | `SAVEPOLISTREATYIN_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_POLIS_TREATYIN | `ASM-FW-GISFW-INT-JSON_POLIS_TREATYIN / ASM!SAVEPOLISTREATYIN_SQL` |
| `SaveTreatyIn.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `SAVETREATYIN` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_TREATY_IN | — |
| `SearchAccumulationbypersetase_SQL.xml` | `ASM-FW-GISFW-INT-ACCUMULATION` | `SEARCHACCUMULATIONBYPERSETASE_SQL` | [terverifikasi] sqlKind=QUERY; procs=UTL_MATCH.EDIT_DISTANCE_SIMILARITY; sqlOps=SELECT; tables=ACCUMULATION,PREP_DATA,PREP_PARAM | — |
| `SearchClassConstructionReas_SQL.xml` | `ASM-FW-GISFW-INT-CONSTRUCTION` | `SEARCHCLASSCONSTRUCTIONREAS_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CONSTRUCTION | — |
| `SearchClauseArgFireSQL.xml` | `ASM-FW-GISFW-INT-ARGCLAUSEFIRE` | `SEARCHCLAUSEARGFIRESQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_ARGCLAUSEFIRE,M_CLAUSE | `ASM-FW-GISFW-INT-V_CLAUSES_ARG / ASM!SEARCHCLAUSEARGFIRESQL` |
| `SearchClobClauseSQL.xml` | `ASM-FW-GISFW-INT-CLAUSE` | `SEARCHCLOBCLAUSESQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_CLAUSE | — |
| `SearchCoinsSQL.xml` | `ASM-FW-GISFW-INT-V_COINS` | `SEARCHCOINSSQL` | [terverifikasi] sqlKind=QUERY; procs=GENERAL.F_GET_NM_ASURADUR; sqlOps=SELECT; tables=GENERAL.LST_ASURADUR,GENERAL.LST_DET_CABANG | — |
| `SearchCoverageIDSQL.xml` | `ASM-FW-GISFW-INT-COVERAGE` | `SEARCHCOVERAGEIDSQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=COVERAGE | `ASM-FW-GISFW-WORK-INT-V_ZIPCODE / ASM!SEARCHKOTAIDSQL` |
| `SearchFirstLossScale1_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `SEARCHFIRSTLOSSSCALE1_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_EQS_MULTIPLIER | — |
| `SearchIndemnityRate_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `SEARCHINDEMNITYRATE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_BI_INDEMNITY | — |
| `SearchJobID.xml` | `ASM-FW-GISFW-INT-V_JOB_PA` | `SEARCHJOBIDSQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=GENERAL.V_JOB_PA | `ASM-FW-GISFW-INT-CLAUSE_TEMPLATE / ASM!SEARCHTEMPLATEMAINCLAUSESQL` |
| `SearchJobIDSQL.xml` | `ASM-FW-GISFW-INT-V_JOB_PA` | `SEARCHJOBIDSQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=GENERAL.V_JOB_PA | `ASM-FW-GISFW-INT-CLAUSE_TEMPLATE / ASM!SEARCHTEMPLATEMAINCLAUSESQL` |
| `SearchMKonstruksiSQL.xml` | `ASM-FW-GISFW-INT-M_KONSTRUKSI` | `SEARCHMKONSTRUKSISQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_KONSTRUKSI | `ASM-FW-GISFW-INT-PREMIUM_PERL_TEMPLATE / ASM!SEARCHTEMPLATEADDITIONALCOVERAGESQL` |
| `SearchOccupationIDSQL.xml` | `ASM-FW-GISFW-INT-OCCUPATION` | `SEARCHOCCUPATIONIDSQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=OCCUPATION | `ASM-FW-GISFW-INT-COVERAGE / ASM!SEARCHCOVERAGEIDSQL` |
| `SearchPaymentSQL.xml` | `ASM-FW-GISFW-INT-BUSINESS` | `SEARCHPAYMENTSQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.BUSINESS | — |
| `SearchPlanTravel.xml` | `ASM-FW-GISFW-INT-PLANTRAVEL` | `SEARCHPLANTRAVEL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_PLANTRAVEL | — |
| `SearchRWIDSQL.xml` | `ASM-FW-GISFW-INT-RW` | `SEARCHRWIDSQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=RW | `ASM-FW-GISFW-WORK-INT-V_ZIPCODE / ASM!SEARCHKOTAIDSQL` |
| `SearchRatePolisEQS_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `SEARCHRATEPOLISEQS_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_EQS_RATE@ASMD.SINARMAS.CO.ID,ZONES | — |
| `SearchRatePolisFlood_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `SEARCHRATEPOLISFLOOD_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=FIRE.M_FLOOD_RATE@ASMD.SINARMAS.CO.ID,M_FLOOD_AREA@ASMD.SINARMAS.CO.ID | — |
| `SearchRatePolisRsmd_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `SEARCHRATEPOLISRSMD_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_RSMD_RATE@ASMD.SINARMAS.CO.ID,OCCUPATION,ZONES | — |
| `SearchRatePolisTerrorism_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `SEARCHRATEPOLISTERRORISM_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_TERORISME_RATE@ASMD.SINARMAS.CO.ID,OCCUPATION | — |
| `SearchRatePolis_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `SEARCHRATEPOLIS_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_FLEXAS_RATE@ASMD.SINARMAS.CO.ID | — |
| `SearchRiskID.xml` | `ASM-FW-GISFW-INT-RISK` | `SEARCHRISKID` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_RISK | — |
| `SearchStatusCabSQL.xml` | `ASM-FW-GISFW-INT-BRANCH` | `SEARCHSTATUSCABSQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_BRANCH | — |
| `SearchTreatyGroupName_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `SEARCHTREATYGROUPNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYBUSINESS | — |
| `SelectSpreadingTreatyInProduction.xml` | `ASM-FW-GISFW-INT-REINSURANCETYPE` | `SELECTSPREADINGTREATYINPRODUCTION` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=REINSURANCETYPE | — |
| `SumAmountMaxTreatyCapacity.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `SUMAMOUNTMAXTREATYCAPACITY` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=PROPORTIONALARRG,TREATYBUSINESS | — |
| `TreatyRealizationCheckDuplicate.xml` | `ASM-FW-GISFW-WORK` | `TREATYREALIZATIONCHECKDUPLICATE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TREATYINPRODUCTION | — |
| `UpdateErrorNoteJsonPolisMonitoring.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `UPDATEERRORNOTEJSONPOLISMONITORING` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=POOLDATA.JSON_POLIS_MONITORING | `ASM-FW-GISFW-INT-POLICYJSON / ASM!UPDATEERRORNOTEJSONPOLIS` |
| `UpdateFacProdObjItemID_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `UPDATEFACPRODOBJITEMID_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT,UPDATE; tables=FACINPRODUCTION,V_JN_OBJ_ITEM | — |
| `UpdateJsonOffer_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `UPDATEJSONOFFER_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=JSON_OFFER | `ASM-FW-GISFW-INT-POLICYJSON / ASM!UPDATEPOLISENDORSEMENT_SQL` |
| `UpdateMasterAccumulatedType.xml` | `ASM-FW-GISFW-INT-ACCUMULATEDTYPE` | `UPDATEMASTERACCUMULATEDTYPE` | [terverifikasi] sqlKind=PLSQL; procs=DBMS_LOB.CREATETEMPORARY,POOLDATA.RDBMASTERACCUMULATEDTYPE | `ASM-FW-GISFW-INT-TYPESHIP / ASM!UPDATEMASTERTYPESHIP` |
| `UpdateMasterAccumulation.xml` | `ASM-FW-GISFW-INT-ACCUMULATION` | `UPDATEMASTERACCUMULATION` | [terverifikasi] sqlKind=PLSQL; procs=DBMS_LOB.CREATETEMPORARY,POOLDATA.RDBMASTERACCUMULATION | — |
| `UpdateMasterBranch.xml` | `ASM-FW-GISFW-INT-BRANCH` | `UPDATEMASTERBRANCH` | [terverifikasi] sqlKind=PLSQL; procs=DBMS_LOB.CREATETEMPORARY,POOLDATA.RDBMASTERBRANCH | — |
| `UpdateMasterCity.xml` | `ASM-FW-GISFW-INT-CITY` | `UPDATEMASTERCITY` | [terverifikasi] sqlKind=PLSQL; procs=DBMS_LOB.CREATETEMPORARY,POOLDATA.RDBMASTERCITY | — |
| `UpdateMasterDistrict.xml` | `ASM-FW-GISFW-INT-DISTRICT` | `UPDATEMASTERDISTRICT` | [terverifikasi] sqlKind=PLSQL; procs=DBMS_LOB.CREATETEMPORARY,POOLDATA.RDBMASTERDISTRICT | — |
| `UpdateMasterMarketingOfficer.xml` | `ASM-FW-GISFW-INT-MARKETINGOFFICER` | `UPDATEMASTERMARKETINGOFFICER` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_MARKETINGOFFICER | — |
| `UpdateMasterNation.xml` | `ASM-FW-GISFW-INT-NATION` | `UPDATEMASTERNATION` | [terverifikasi] sqlKind=PLSQL; procs=DBMS_LOB.CREATETEMPORARY,POOLDATA.RDBMASTERNATION | — |
| `UpdateMasterProvince.xml` | `ASM-FW-GISFW-INT-PROVINCE` | `UPDATEMASTERPROVINCE` | [terverifikasi] sqlKind=PLSQL; procs=DBMS_LOB.CREATETEMPORARY,POOLDATA.RDBMASTERPROVINCE | — |
| `UpdateMasterRW.xml` | `ASM-FW-GISFW-INT-RW` | `UPDATEMASTERRW` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.RDBMASTERRW | — |
| `UpdateMasterRiskAddress_SQL.xml` | `ASM-FW-GISFW-INT-RISKADDRESS` | `UPDATEMASTERRISKADDRESS_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.INSERTUPDATERISKADDRESS | `ASM-FW-GISFW-INT-RISKADDRESS / ASM!UPDATEMASTERRISKADDRESS` |
| `UpdateMasterShip.xml` | `ASM-FW-GISFW-INT-SHIP` | `UPDATEMASTERSHIP` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.RDBMASTERSHIP | `ASM-FW-GISFW-INT-CARGROUP / ASM!UPDATEMASTERCARGROUP` |
| `UpdatePolisEndorsement_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `UPDATEPOLISENDORSEMENT_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=JSON_POLIS | — |
| `UpdateSTSKonversi_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `UPDATESTSKONVERSI_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=JSON_POLIS | — |
| `UpdateTgl_BindJsonOffer_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `UPDATETGL_BINDJSONOFFER_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=JSON_OFFER | — |
| `Update_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `UPDATE_T_STORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |
| `Update_sql.xml` | `ASM-FW-GISFW-INT-ACCUMULATION_LIFE` | `UPDATE_SQL` | [terverifikasi] sqlKind=PLSQL; procs=DBMS_LOB.CREATETEMPORARY,POOLDATA.PEGA_M_ACCUMULATION_LIFE | — |
| `ViewIndemnity2_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `VIEWINDEMNITY2_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=FIRE.M_BI_INDEMNITY@ASMD.SINARMAS.CO.ID | — |
| `ViewIndemnity_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `VIEWINDEMNITY_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=FIRE.M_BI_INDEMNITY@ASMD.SINARMAS.CO.ID | — |
| `ViewMasterCargo.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `VIEWMASTERCARGO` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.MASTERCARGO | — |
| `ViewMasterCargo2.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `VIEWMASTERCARGO2` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.JSON_POLIS | `ASM-FW-GISFW-DATA-OFFERFACIN / ASM!VIEWMASTERCARGO` |
| `getTglInputJsonPolis1_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETTGLINPUTJSONPOLIS1_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | — |
| `getTglInputJsonPolis_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETTGLINPUTJSONPOLIS_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | — |

### ReportDefinition — 123 rule (`RULE-OBJ-REPORT-DEFINITION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseAccesorryBrand_RD.xml` | `ASM-FW-GISFW-INT-ACCESORYBRAND` | `BROWSEACCESORRYBRAND_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseAccesory_RD.xml` | `ASM-FW-GISFW-INT-ACCESORY` | `BROWSEACCESORY_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseAccumulatedType_RD.xml` | `ASM-FW-GISFW-INT-ACCUMULATEDTYPE` | `BROWSEACCUMULATEDTYPE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseAccumulationLife_RD.xml` | `ASM-FW-GISFW-INT-ACCUMULATION_LIFE` | `BROWSEACCUMULATIONLIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseAgentHierarkiList_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSEAGENTHIERARKILIST_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseAgentNonLife_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSEAGENTNONLIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-AGENT / BROWSEAGENTNUSARE_RD` |
| `BrowseAgentNusaReLife_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSEAGENTNUSARELIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-AGENT / BROWSEAGENTNUSARE_RD` |
| `BrowseAgentNusaRe_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSEAGENTNUSARE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseAgent_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSEAGENT_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-AGENT / BROWSEAGENTNUSARE_RD` |
| `BrowseBankGroup.xml` | `ASM-FW-GISFW-INT-LST_BANK_GROUP` | `BROWSEBANKGROUP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseBenefitClause.xml` | `ASM-FW-GISFW-INT-VIEW_BENEFIT_PROPERTY` | `BROWSEBENEFITCLAUSE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseBenefitPkg.xml` | `ASM-FW-GISFW-INT-V_PKG_BENEFIT` | `BROWSEBENEFITPKG` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseBranchDetail_RD.xml` | `ASM-FW-GISFW-INT-BRANCHDETAIL` | `BROWSEBRANCHDETAIL_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseBranch_RD.xml` | `ASM-FW-GISFW-INT-BRANCH` | `BROWSEBRANCH_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseBrandDetail_RD.xml` | `ASM-FW-GISFW-INT-BRANDDETAIL` | `BROWSEBRANDDETAIL_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseBrand_RD.xml` | `ASM-FW-GISFW-INT-BRAND` | `BROWSEBRAND_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-BRAND / BROWSEBRAND` |
| `BrowseBusiness_RD.xml` | `ASM-FW-GISFW-INT-BUSINESS` | `BROWSEBUSINESS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCZoneIsNotNull_RD.xml` | `ASM-FW-GISFW-INT-CZONE` | `BROWSECZONEISNOTNULL_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-CZONE / BROWSECZONEISNOTENULL_RD` |
| `BrowseCedingCo_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSECEDINGCO_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-AGENT / BROWSEAGENTHIERARKILIST_RD` |
| `BrowseCityInput_RD.xml` | `ASM-FW-GISFW-INT-CITY` | `BROWSECITYINPUT_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-CITY / BROWSECITY_RD` |
| `BrowseCity_RD.xml` | `ASM-FW-GISFW-INT-RW` | `BROWSECITY_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-RW / BROWSEPROVINCE_RD` |
| `BrowseClauseAneka_RD.xml` | `ASM-FW-GISFW-INT-CLAUSE` | `BROWSECLAUSEANEKA_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseClauseMBU_RD.xml` | `ASM-FW-GISFW-INT-CLAUSE` | `BROWSECLAUSEMBU_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-CLAUSE / BROWSECLAUSEPA_RD` |
| `BrowseClausePA_RD.xml` | `ASM-FW-GISFW-INT-CLAUSE` | `BROWSECLAUSEPA_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-CLAUSE / BROWSECLAUSE_RD` |
| `BrowseClausePlan.xml` | `ASM-FW-GISFW-INT-M_KLAUSUL_PLAN` | `BROWSECLAUSEPLAN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseClieent_RD.xml` | `ASM-FW-GISFW-INT-CLIENT` | `BROWSECLIEENT_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseColor_RD.xml` | `ASM-FW-GISFW-INT-COLOR` | `BROWSECOLOR_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseConveyanceType.xml` | `ASM-FW-GISFW-INT-MCONVEYANCE` | `BROWSECONVEYANCETYPE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseConveyance_RD.xml` | `ASM-FW-GISFW-INT-CONVEYANCE` | `BROWSECONVEYANCE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCoverageAneka.xml` | `ASM-FW-GISFW-INT-V_COVERAGE_ANEKA` | `BROWSECOVERAGEANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCoverageFacIn_RD.xml` | `ASM-FW-GISFW-INT-COVERAGE_FACIN` | `BROWSECOVERAGEFACIN_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCoverageMBU.xml` | `ASM-FW-GISFW-INT-COVERAGE` | `BROWSECOVERAGEMBU` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCoverage_RD.xml` | `ASM-FW-GISFW-INT-COVERAGE` | `BROWSECOVERAGE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCurrencyTreatyIn_RD.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `BROWSECURRENCYTREATYIN_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCurrency_RD.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `BROWSECURRENCY_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseDeductMBU_RD.xml` | `ASM-FW-GISFW-INT-DEDUCTIBLE` | `BROWSEDEDUCTMBU_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseDeductibleFacIn_RD.xml` | `ASM-FW-GISFW-INT-DEDUCTIBLE` | `BROWSEDEDUCTIBLEFACIN_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseDistrictInputC_RD.xml` | `ASM-FW-GISFW-INT-DISTRICT` | `BROWSEDISTRICTINPUTC_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseDistrictInput_RD.xml` | `ASM-FW-GISFW-INT-DISTRICT` | `BROWSEDISTRICTINPUT_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-DISTRICT / BROWSEDISTRICT_RD` |
| `BrowseDistrict_RD.xml` | `ASM-FW-GISFW-INT-RW` | `BROWSEDISTRICT_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-RW / BROWSECITY_RD` |
| `BrowseEnumData.xml` | `ASM-FW-GISFW-DATA-ENUMERATION` | `BROWSEENUMDATA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseGoodsType.xml` | `ASM-FW-GISFW-INT-GOODSTYPE` | `BROWSEGOODSTYPE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseJTypePropertyPlan.xml` | `ASM-FW-GISFW-INT-VJ_M_TYPE_PROPERTY_PLAN` | `BROWSEJTYPEPROPERTYPLAN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseJobPA.xml` | `ASM-FW-GISFW-INT-V_JOB_PA` | `BROWSEJOBPA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseJob_RD.xml` | `ASM-FW-GISFW-INT-JOB` | `BROWSEJOB_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseLLOYDAGENT_RD.xml` | `ASM-FW-GISFW-INT-LLOYDAGENT` | `BROWSELLOYDAGENT_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseLimit.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSELIMIT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PROPORTIONALARRG / BROWSETREATYARRANGEMENT_PARENTREINS` |
| `BrowseM_KONSTRUKSIList.xml` | `ASM-FW-GISFW-INT-M_KONSTRUKSI` | `BROWSEM_KONSTRUKSILIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseM_OKUPASI_ANEKA.xml` | `ASM-FW-GISFW-INT-M_OKUPASI_ANEKA` | `BROWSEM_OKUPASI_ANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseMarineCondition_RD.xml` | `ASM-FW-GISFW-INT-MARINECONDITION` | `BROWSEMARINECONDITION_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseMarketingOfficer_RD.xml` | `ASM-FW-GISFW-INT-MARKETINGOFFICER` | `BROWSEMARKETINGOFFICER_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseMaxAge_RD.xml` | `ASM-FW-GISFW-INT-MAXAGE` | `BROWSEMAXAGE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseMaxHariTravel_RD.xml` | `ASM-FW-GISFW-INT-PREMITRAVEL` | `BROWSEMAXHARITRAVEL_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseMerk_RD.xml` | `ASM-FW-GISFW-INT-BRAND` | `BROWSEMERK_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseModelList.xml` | `ASM-FW-GISFW-INT-BRANDDETAIL` | `BROWSEMODELLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseNation_RD.xml` | `ASM-FW-GISFW-INT-NATION` | `BROWSENATION_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseOC_RD.xml` | `ASM-FW-GISFW-INT-OC` | `BROWSEOC_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseObjectItemType_RD.xml` | `ASM-FW-GISFW-INT-OBJECTITEMTYPE` | `BROWSEOBJECTITEMTYPE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseObject_RD.xml` | `ASM-FW-GISFW-INT-OBJECT` | `BROWSEOBJECT_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseOccupationANEKA_RD.xml` | `ASM-FW-GISFW-INT-OCCUPATION` | `BROWSEOCCUPATIONANEKA_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseOccupationFIRE_RD.xml` | `ASM-FW-GISFW-INT-OCCUPATION` | `BROWSEOCCUPATIONFIRE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseOccupationFacInFIRE_RD.xml` | `ASM-FW-GISFW-INT-OCCUPATION` | `BROWSEOCCUPATIONFACINFIRE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseOccupationList_RD.xml` | `ASM-FW-GISFW-INT-OCCUPATION` | `BROWSEOCCUPATIONLIST_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseOccupationMBU_RD.xml` | `ASM-FW-GISFW-INT-OCCUPATION` | `BROWSEOCCUPATIONMBU_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseOccupations_RD.xml` | `ASM-FW-GISFW-INT-OCCUPATIONS` | `BROWSEOCCUPATIONS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowsePacking_RD.xml` | `ASM-FW-GISFW-INT-PACKING` | `BROWSEPACKING_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowsePlanByAge.xml` | `ASM-FW-GISFW-INT-V_PACKAGE_AGE_KLAUSUL` | `BROWSEPLANBYAGE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowsePlanProperty.xml` | `ASM-FW-GISFW-INT-M_TYPE_PROPERTY_PLAN` | `BROWSEPLANPROPERTY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowsePremiTravel.xml` | `ASM-FW-GISFW-INT-VIEW_PREMI_TRAVEL` | `BROWSEPREMITRAVEL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseProvince2_RD.xml` | `ASM-FW-GISFW-INT-PROVINCE` | `BROWSEPROVINCE2_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseProvince_RD.xml` | `ASM-FW-GISFW-INT-RW` | `BROWSEPROVINCE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseRWInput_RD.xml` | `ASM-FW-GISFW-INT-RW` | `BROWSERWINPUT_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseRW_RD.xml` | `ASM-FW-GISFW-INT-RW` | `BROWSERW_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseRateLifeSummary.xml` | `ASM-FW-GISFW-INT-RATE_LIFE_SUMMARY` | `BROWSERATELIFESUMMARY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseReinsuranceType_RD.xml` | `ASM-FW-GISFW-INT-REINSURANCETYPE` | `BROWSEREINSURANCETYPE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseRetrocessionLife_RD.xml` | `ASM-FW-GISFW-INT-RETROCESSIONLIFE` | `BROWSERETROCESSIONLIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseRiskAddressZipCode_RD.xml` | `ASM-FW-GISFW-INT-RISKADDRESS` | `BROWSERISKADDRESSZIPCODE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-RISKADDRESS / BROWSERISKADDRESS_RD` |
| `BrowseRiskAddress_RD.xml` | `ASM-FW-GISFW-INT-RISKADDRESS` | `BROWSERISKADDRESS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseRisksAddress_RD.xml` | `ASM-FW-GISFW-INT-RISKADDRESS` | `BROWSERISKSADDRESS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-RISKADDRESS / BROWSERISKADDRESS_RD` |
| `BrowseSameAccumulationLife_RD.xml` | `ASM-FW-GISFW-INT-ACCUMULATION_LIFE` | `BROWSESAMEACCUMULATIONLIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-ACCUMULATION_LIFE / BROWSEACCUMULATIONLIFE_RD` |
| `BrowseSearchSobCeding_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSESEARCHSOBCEDING_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-AGENT / BROWSECEDINGCO_RD` |
| `BrowseShip2_RD.xml` | `ASM-FW-GISFW-INT-SHIP` | `BROWSESHIP2_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-SHIP / BROWSESHIP_RD` |
| `BrowseShip_RD.xml` | `ASM-FW-GISFW-INT-SHIP` | `BROWSESHIP_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseSourceBizList.xml` | `ASM-FW-GISFW-INT-MV_AGEN` | `BROWSESOURCEBIZLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-MV_AGEN / BROWSESUMBISLIST` |
| `BrowseSubContract_RD.xml` | `ASM-FW-GISFW-INT-SUBCONTRACT` | `BROWSESUBCONTRACT_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseSubjectTo_RD.xml` | `ASM-FW-GISFW-INT-SUBJECTTO` | `BROWSESUBJECTTO_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTableOfLimit_RD.xml` | `ASM-FW-GISFW-INT-TABLEOFLIMIT` | `BROWSETABLEOFLIMIT_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTeritory_RD.xml` | `ASM-FW-GISFW-INT-RW` | `BROWSETERITORY_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-RW / BROWSEDISTRICT_RD` |
| `BrowseTrading_RD.xml` | `ASM-FW-GISFW-INT-TRADING` | `BROWSETRADING_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyYearOR_Life_RD.xml` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` | `BROWSETREATYYEAROR_LIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE / BROWSETREATYYEAR_LIFE_RD` |
| `BrowseTreatyYear_Life_RD.xml` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` | `BROWSETREATYYEAR_LIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTypeList_RD.xml` | `ASM-FW-GISFW-INT-BRANDDETAIL` | `BROWSETYPELIST_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTypeShip_RD.xml` | `ASM-FW-GISFW-INT-TYPESHIP` | `BROWSETYPESHIP_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseType_RD.xml` | `ASM-FW-GISFW-INT-BRANDDETAIL` | `BROWSETYPE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseUserAssignment.xml` | `ASM-FW-GISFW-INT-V_USER_ASSIGNMENT` | `BROWSEUSERASSIGNMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseVClauseArg.xml` | `ASM-FW-GISFW-INT-V_CLAUSES_ARG` | `BROWSEVCLAUSEARG` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseVDraftWordingList.xml` | `ASM-FW-GISFW-INT-V_DRAFTWORDING` | `BROWSEVDRAFTWORDINGLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseVMarketingLeaderList.xml` | `ASM-FW-GISFW-INT-V_MARKETING_LEADER` | `BROWSEVMARKETINGLEADERLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseV_ALM_RISIKOList.xml` | `ASM-FW-GISFW-INT-V_ALM_RISIKO` | `BROWSEV_ALM_RISIKOLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseV_JN_OBJ_ITEM.xml` | `ASM-FW-GISFW-INT-V_JN_OBJ_ITEM` | `BROWSEV_JN_OBJ_ITEM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseWarranty_RD.xml` | `ASM-FW-GISFW-INT-WARRANTY` | `BROWSEWARRANTY_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseZones_RD.xml` | `ASM-FW-GISFW-INT-ZONES` | `BROWSEZONES_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowserAccessoryType_RD.xml` | `ASM-FW-GISFW-INT-ACCESORYTYPE` | `BROWSERACCESSORYTYPE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CheckMST_KENDARAANByChas.xml` | `ASM-FW-GISFW-INT-V_DUPLICATE_VEHICLE` | `CHECKMST_KENDARAANBYCHAS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-V_DUPLICATE_VEHICLE / CHECKMST_KENDARAANBYPLATCHASENGINE` |
| `CheckMST_KENDARAANByEngine.xml` | `ASM-FW-GISFW-INT-V_DUPLICATE_VEHICLE` | `CHECKMST_KENDARAANBYENGINE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-V_DUPLICATE_VEHICLE / CHECKMST_KENDARAANBYPLATCHASENGINE` |
| `CheckMST_KENDARAANByPlat.xml` | `ASM-FW-GISFW-INT-V_DUPLICATE_VEHICLE` | `CHECKMST_KENDARAANBYPLAT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-V_DUPLICATE_VEHICLE / CHECKMST_KENDARAANBYPLATCHASENGINE` |
| `DataTableEditorReport.xml` | `ASM-FW-GISFW-DATA-ENUMERATION` | `DATATABLEEDITORREPORT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GetAllDocument.xml` | `ASM-FW-GISFW-INT-DOCUMENT_POLIS` | `GETALLDOCUMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GetListOpportunity.xml` | `ASM-FW-SFAGISFW-WORK-OPPORTUNITY` | `GETLISTOPPORTUNITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GetListOpportunityF.xml` | `ASM-FW-SFAGISFW-WORK-OPPORTUNITY` | `GETLISTOPPORTUNITYF` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GetListOpportunityLife.xml` | `ASM-FW-SFAGISFW-WORK-OPPORTUNITY` | `GETLISTOPPORTUNITYLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GetMasterPolicyLife.xml` | `ASM-FW-GISFW-WORK-NB` | `GETMASTERPOLICYLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-NB / GETINBOXSFA` |
| `GetObjectItem.xml` | `ASM-FW-GISFW-INT-V_JN_OBJ_ITEM` | `GETOBJECTITEM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-V_JN_OBJ_ITEM / BROWSEV_JN_OBJ_ITEM` |
| `GetOperatorIDbyWorkBasket_RD.xml` | `DATA-ADMIN-OPERATOR-ID` | `GETOPERATORIDBYWORKBASKET_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-ADMIN-OPERATOR-ID / GETEMAILUSER_RD` |
| `GetTeamGroup_RD.xml` | `ASM-FW-GISFW-INT-MARKETINGOFFICER` | `GETTEAMGROUP_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InwardScale_RD.xml` | `ASM-FW-GISFW-INT-INWARDSCALE` | `INWARDSCALE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ReasGetAllAttachments.xml` | `LINK-ATTACHMENT` | `REASGETALLATTACHMENTS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SearchDueDate_RD.xml` | `ASM-FW-GISFW-INT-LST_BUSINESS_DUE_DATE` | `SEARCHDUEDATE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SearchRiskAccumulation_RD.xml` | `ASM-FW-GISFW-INT-ACCUMULATION` | `SEARCHRISKACCUMULATION_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-ACCUMULATION / BROWSEACCUMULATION_RD` |
| `SelectLeader_RD.xml` | `ASM-FW-GISFW-INT-MARKETINGOFFICER` | `SELECTLEADER_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-MARKETINGOFFICER / BROWSEMARKETINGOFFICERALL_RD` |
| `pyDefaultReport.xml` | `@BASECLASS` | `PYDEFAULTREPORT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `pyGetAllAttachments.xml` | `LINK-ATTACHMENT` | `PYGETALLATTACHMENTS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `pyGetListOfOperators.xml` | `DATA-ADMIN-OPERATOR-ID` | `PYGETLISTOFOPERATORS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Section — 432 rule (`RULE-HTML-SECTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AccumulationRisk.xml` | `DATA-PORTAL` | `ACCUMULATIONRISK` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `AddFacOfferList.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `ADDFACOFFERLIST` | [terverifikasi] pyInclude=2 | — |
| `AddFacOfferList2.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `ADDFACOFFERLIST2` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-FACOFFER / ADDFACOFFERLIST` |
| `AddFacOfferList_IsUW.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `ADDFACOFFERLIST_ISUW` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-FACOFFER / ADDFACOFFERLIST` |
| `AdjustmentRiskAccumulation.xml` | `DATA-PORTAL` | `ADJUSTMENTRISKACCUMULATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `AllSummarySection.xml` | `ASM-FW-GISFW-WORK` | `ALLSUMMARYSECTION` | [terverifikasi] pyInclude=8 | `ASM-FW-GISFW-WORK / FIRESUMMARYSECTION` |
| `AnalysLocationbyAI_Sec.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `ANALYSLOCATIONBYAI_SEC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `AnekaListFacOut.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `ANEKALISTFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `AnekaListFacOut_IsUW.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `ANEKALISTFACOUT_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OCCUPATION / ANEKALISTFACOUT` |
| `AttachmentGridReas.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `ATTACHMENTGRIDREAS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ATTACHMENTGRIDREAS` |
| `BenefitClauseRO.xml` | `ASM-FW-GISFW-DATA-BENEFIT` | `BENEFITCLAUSERO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-BENEFIT / BENEFITCLAUSE` |
| `BenefitListRO.xml` | `ASM-FW-GISFW-DATA-PLAN` | `BENEFITLISTRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CallDualScoringRisk.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `CALLDUALSCORINGRISK` | [terverifikasi] pyInclude=4 | — |
| `CauseOfLossClaim_FacIn.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `CAUSEOFLOSSCLAIM_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CauseOfLoss_FacIn.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `CAUSEOFLOSS_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CauseOfLoss_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `CAUSEOFLOSS_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / CAUSEOFLOSS_FACIN` |
| `CedingCedant.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `CEDINGCEDANT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CedingCedantHierarki.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `CEDINGCEDANTHIERARKI` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-QUOTATION / CEDINGCOHIERARKI` |
| `CedingCedant_IsUW.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `CEDINGCEDANT_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-QUOTATION / CEDINGCEDANT` |
| `CedingCoHierarki.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `CEDINGCOHIERARKI` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-QUOTATION / SOURCEHIERARKI` |
| `ChooseAccumulation_FacIn.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `CHOOSEACCUMULATION_FACIN` | [terverifikasi] pyInclude=2 | — |
| `ChooseClassofContraction.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `CHOOSECLASSOFCONTRACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChooseClause.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `CHOOSECLAUSE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChooseClauseFire.xml` | `ASM-FW-GISFW-WORK` | `CHOOSECLAUSEFIRE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChooseCoverage.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `CHOOSECOVERAGE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChooseDeductible.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `CHOOSEDEDUCTIBLE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChooseObject_Ship.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `CHOOSEOBJECT_SHIP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChooseOccupation.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `CHOOSEOCCUPATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChooseOccupation_Dtl.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `CHOOSEOCCUPATION_DTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChooseRiskAddress.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `CHOOSERISKADDRESS` | [terverifikasi] pyInclude=2 | — |
| `ChooseRiskAddress_ResultList.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `CHOOSERISKADDRESS_RESULTLIST` | [terverifikasi] pyInclude=1 | — |
| `ChooseSubContract.xml` | `ASM-FW-GISFW-DATA-SUBCONTRACT` | `CHOOSESUBCONTRACT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChooseSubContract_FacIn.xml` | `ASM-FW-GISFW-DATA-SUBCONTRACT` | `CHOOSESUBCONTRACT_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-SUBCONTRACT / CHOOSESUBCONTRACT` |
| `ChooseZipCodeDtl.xml` | `@BASECLASS` | `CHOOSEZIPCODEDTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ClauseList.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `CLAUSELIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ClauseValue.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `CLAUSEVALUE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ComfirmPolis_Life.xml` | `ASM-FW-GISFW-WORK` | `COMFIRMPOLIS_LIFE` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-WORK / COMFIRMPOLIS_NONLIFE` |
| `ComfirmPolis_NonLife.xml` | `ASM-FW-GISFW-WORK` | `COMFIRMPOLIS_NONLIFE` | [terverifikasi] pyInclude=4 | — |
| `CommisionCoverageAneka.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `COMMISIONCOVERAGEANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-ANEKA / SPREADINGCOVERAGEANEKA` |
| `CommisionItem.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `COMMISIONITEM` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-COVERAGE / SPREADINGITEM` |
| `CommisionItemAneka.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `COMMISIONITEMANEKA` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-COVERAGE / SPREADINGITEMANEKA` |
| `CommisionObjectAneka.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `COMMISIONOBJECTANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OCCUPATION / SPREADINGOBJECTANEKA` |
| `ConfirmChangeOutgo.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `CONFIRMCHANGEOUTGO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CopyCoverageFrom.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `COPYCOVERAGEFROM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Correspondence.xml` | `ASM-FW-GISFW-WORK` | `CORRESPONDENCE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CorrespondenceContent.xml` | `ASM-FW-GISFW-DATA-CORRESPONDENCE` | `CORRESPONDENCECONTENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CoverageCommisionList.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `COVERAGECOMMISIONLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN / COVERAGESPREADINGLIST` |
| `CoverageCommisionList_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `COVERAGECOMMISIONLIST_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN / COVERAGECOMMISIONLIST` |
| `CoverageFormula.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `COVERAGEFORMULA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CoverageItem.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `COVERAGEITEM` | [terverifikasi] pyInclude=2 | — |
| `CoverageItem_IsUW.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `COVERAGEITEM_ISUW` | [terverifikasi] pyInclude=2 | — |
| `CoverageList.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `COVERAGELIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CoverageList_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `COVERAGELIST_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN / COVERAGELIST` |
| `CoveragePropertyFacOut.xml` | `ASM-FW-GISFW-DATA-PROPERTY` | `COVERAGEPROPERTYFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CoverageSpreadingList.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `COVERAGESPREADINGLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CoverageSpreadingList_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `COVERAGESPREADINGLIST_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN / COVERAGESPREADINGLIST` |
| `CoverageSummary.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `COVERAGESUMMARY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CreateListClause.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `CREATELISTCLAUSE` | [terverifikasi] pyInclude=4 | `ASM-FW-GISFW-DATA-PLAN / CREATELISTCLAUSE` |
| `CreateListWarranty.xml` | `ASM-FW-GISFW-WORK` | `CREATELISTWARRANTY` | [terverifikasi] pyInclude=4 | — |
| `DeductibleDetail.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `DEDUCTIBLEDETAIL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DetailCoverageListFire_GCNM.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `DETAILCOVERAGELISTFIRE_GCNM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DetailCoverageTravelGrid_GCNM.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `DETAILCOVERAGETRAVELGRID_GCNM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / DETAILDEDUCTIBLETRAVELGRID_GCNM` |
| `DetailDeptHeadTreatyIn_UW.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILDEPTHEADTREATYIN_UW` | [terverifikasi] pyInclude=5 | — |
| `DetailLocation.xml` | `ASM-FW-GISFW-WORK` | `DETAILLOCATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DetailPoliciesNonProportional.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILPOLICIESNONPROPORTIONAL` | [terverifikasi] pyInclude=4 | — |
| `DetailPolicyTreatyIn.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILPOLICYTREATYIN` | [terverifikasi] pyInclude=6 | — |
| `DetailPolicyTreatyInNonProportional.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILPOLICYTREATYINNONPROPORTIONAL` | [terverifikasi] pyInclude=2 | — |
| `DetailPolicyTreatyInNonProportionalEDM.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILPOLICYTREATYINNONPROPORTIONALEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-POLICYTREATYIN / DETAILPOLICYTREATYINNONPROPORTIONAL` |
| `DetailPolicyTreatyOutNonProportional.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILPOLICYTREATYOUTNONPROPORTIONAL` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-POLICYTREATYIN / DETAILPOLICYTREATYINNONPROPORTIONAL` |
| `EmailSection.xml` | `ASM-FW-GISFW-WORK` | `EMAILSECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `EmailSectionCeding.xml` | `ASM-FW-GISFW-WORK` | `EMAILSECTIONCEDING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `EmailSection_IsUW.xml` | `ASM-FW-GISFW-WORK` | `EMAILSECTION_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `EmailSection_Life.xml` | `ASM-FW-GISFW-WORK` | `EMAILSECTION_LIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / EMAILSECTION` |
| `FEAList.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `FEALIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `FEAList_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `FEALIST_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / FEALIST` |
| `FacOutOffer.xml` | `ASM-FW-GISFW-WORK` | `FACOUTOFFER` | [terverifikasi] pyInclude=2 | — |
| `FacOutPrintRISlipSectionInside.xml` | `ASM-FW-GISFW-WORK` | `FACOUTPRINTRISLIPSECTIONINSIDE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / FACOUTMEMOPLACINGSECTIONINSIDE` |
| `FacultativeLetter.xml` | `ASM-FW-GISFW-WORK` | `FACULTATIVELETTER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `FinalPolicyRO.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `FINALPOLICYRO` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-DATA-OCCUPATION / FINALPOLICY` |
| `FireSummarySection.xml` | `ASM-FW-GISFW-WORK` | `FIRESUMMARYSECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `FormulaDtl.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | `FORMULADTL` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-OFFERFACIN / FORMULADTL` |
| `FormulaTreatyCapacityDesc.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `FORMULATREATYCAPACITYDESC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `FormulaTreatyCapacityDesc_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `FORMULATREATYCAPACITYDESC_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN / FORMULATREATYCAPACITYDESC` |
| `GeneralDeptHeadTreatyIn_UW.xml` | `ASM-FW-GISFW-WORK` | `GENERALDEPTHEADTREATYIN_UW` | [terverifikasi] pyInclude=6 | — |
| `GeneralPolicyTreatyIn.xml` | `ASM-FW-GISFW-WORK` | `GENERALPOLICYTREATYIN` | [terverifikasi] pyInclude=7 | — |
| `GolfCoverageSummary.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `GOLFCOVERAGESUMMARY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / COVERAGESUMMARY` |
| `GolfLocationDetail.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `GOLFLOCATIONDETAIL` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / LOCATIONDETAIL` |
| `GolfObjectItemSummary.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `GOLFOBJECTITEMSUMMARY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / OBJECTITEMSUMMARY` |
| `GolfSummarySection.xml` | `ASM-FW-GISFW-WORK` | `GOLFSUMMARYSECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / FIRESUMMARYSECTION` |
| `GridViewFollowingNB.xml` | `ASM-FW-GISFW-WORK` | `GRIDVIEWFOLLOWINGNB` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / VIEWFOLLOWINGNB` |
| `GrowingTreesLocationDetail.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `GROWINGTREESLOCATIONDETAIL` | [terverifikasi] pyInclude=3 | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / GOLFLOCATIONDETAIL` |
| `GrowingTreesOccupationSummary.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `GROWINGTREESOCCUPATIONSUMMARY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / GOLFOBJECTITEMSUMMARY` |
| `GrowingTreesSummarySection.xml` | `ASM-FW-GISFW-WORK` | `GROWINGTREESSUMMARYSECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / GOLFSUMMARYSECTION` |
| `HistoricalSurveyReportDtl.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `HISTORICALSURVEYREPORTDTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `HistoricalSurveyReportDtlUW.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `HISTORICALSURVEYREPORTDTLUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-POLICYTREATYIN / HISTORICALSURVEYREPORTDTL` |
| `InputAccumulatedType_FacIn.xml` | `@BASECLASS` | `INPUTACCUMULATEDTYPE_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / INPUTACCUMULATEDTYPE` |
| `InputAccumulationCov.xml` | `@BASECLASS` | `INPUTACCUMULATIONCOV` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / INPUTACCUMULATION` |
| `InputAnekaParticipant.xml` | `ASM-FW-GISFW-WORK` | `INPUTANEKAPARTICIPANT` | [terverifikasi] pyInclude=1 | — |
| `InputBranch.xml` | `@BASECLASS` | `INPUTBRANCH` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCauseOfDecline.xml` | `ASM-FW-GISFW-WORK` | `INPUTCAUSEOFDECLINE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCauseOfLoss_FacIn.xml` | `ASM-FW-GISFW-DATA-CAUSEOFLOSS` | `INPUTCAUSEOFLOSS_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCauseOfLoss_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-CAUSEOFLOSS` | `INPUTCAUSEOFLOSS_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-CAUSEOFLOSS / INPUTCAUSEOFLOSS_FACIN` |
| `InputCity.xml` | `@BASECLASS` | `INPUTCITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputClause.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `INPUTCLAUSE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputClauseFire_FacIn.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `INPUTCLAUSEFIRE_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputClauseFire_FacInUW.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `INPUTCLAUSEFIRE_FACINUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-CLAUSE / INPUTCLAUSEFIRE_FACIN` |
| `InputClauseFire_ViewDtl.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `INPUTCLAUSEFIRE_VIEWDTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputClause_ViewDtl.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `INPUTCLAUSE_VIEWDTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCommentDeductible.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `INPUTCOMMENTDEDUCTIBLE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN / EDITMARKETING` |
| `InputCommissionLife_FacIn.xml` | `DATA-PARTY-PERSON` | `INPUTCOMMISSIONLIFE_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCommissionMBU_FacIn.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `INPUTCOMMISSIONMBU_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCommissionMC_FacIn.xml` | `ASM-FW-GISFW-DATA-CARGO` | `INPUTCOMMISSIONMC_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCommissionPA_FacIn.xml` | `DATA-PARTY-PERSON` | `INPUTCOMMISSIONPA_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCoverageAneka_FacIn.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTCOVERAGEANEKA_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCoverageCargo_FacIn.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTCOVERAGECARGO_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCoverageCargo_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTCOVERAGECARGO_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / INPUTCOVERAGECARGO_FACIN` |
| `InputCoverageCommisionFire.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `INPUTCOVERAGECOMMISIONFIRE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTYITEM / INPUTCOVERAGESPREADINGFIRE` |
| `InputCoverageDeductibleAneka_GCNM.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `INPUTCOVERAGEDEDUCTIBLEANEKA_GCNM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCoverageDeductibleHull_GCNM.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `INPUTCOVERAGEDEDUCTIBLEHULL_GCNM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-ANEKA / INPUTCOVERAGEDEDUCTIBLEANEKA_GCNM` |
| `InputCoverageFire.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `INPUTCOVERAGEFIRE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCoverageFire_IsUW.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `INPUTCOVERAGEFIRE_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTYITEM / INPUTCOVERAGEFIRE` |
| `InputCoverageSpreadingAneka.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `INPUTCOVERAGESPREADINGANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTYITEM / INPUTCOVERAGESPREADINGFIRE` |
| `InputCoverageSpreadingAneka_IsUW.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `INPUTCOVERAGESPREADINGANEKA_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-ANEKA / INPUTCOVERAGESPREADINGANEKA` |
| `InputCoverageSpreadingFire.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `INPUTCOVERAGESPREADINGFIRE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputCoverageSpreadingFire_IsUW.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `INPUTCOVERAGESPREADINGFIRE_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTYITEM / INPUTCOVERAGESPREADINGFIRE` |
| `InputCoverageSpreadingMBU.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `INPUTCOVERAGESPREADINGMBU` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-VEHICLE / INPUTCOVERAGESPREADINGMARINECARGO` |
| `InputCoverageSpreadingMBUISUW.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `INPUTCOVERAGESPREADINGMBUISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-VEHICLE / INPUTCOVERAGESPREADINGMBU` |
| `InputCoverageSpreadingMarineCargo.xml` | `ASM-FW-GISFW-DATA-CARGO` | `INPUTCOVERAGESPREADINGMARINECARGO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTYITEM / INPUTCOVERAGESPREADINGMARINECARGO` |
| `InputCoverageSpreadingMarineCargoISUW.xml` | `ASM-FW-GISFW-DATA-CARGO` | `INPUTCOVERAGESPREADINGMARINECARGOISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-CARGO / INPUTCOVERAGESPREADINGMARINECARGO` |
| `InputDeductibleHull_GCNM.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDEDUCTIBLEHULL_GCNM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / INPUTDEDUCTIBLE_GCNM` |
| `InputDeductible_FacIn.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `INPUTDEDUCTIBLE_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDeductible_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `INPUTDEDUCTIBLE_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-DEDUCTIBLE / INPUTDEDUCTIBLE_FACIN` |
| `InputDeductible_GCNM.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDEDUCTIBLE_GCNM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDistrict.xml` | `@BASECLASS` | `INPUTDISTRICT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlAnekaPolicySchedule.xml` | `ASM-FW-GISFW-WORK` | `INPUTDTLANEKAPOLICYSCHEDULE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlCargo.xml` | `ASM-FW-GISFW-DATA-CARGO` | `INPUTDTLCARGO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlCargo_FacIn.xml` | `ASM-FW-GISFW-DATA-CARGO` | `INPUTDTLCARGO_FACIN` | [terverifikasi] pyInclude=1 | — |
| `InputDtlCargo_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-CARGO` | `INPUTDTLCARGO_FACIN_ISUW` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-DATA-CARGO / INPUTDTLCARGO_FACIN` |
| `InputDtlClause_FacIn.xml` | `ASM-FW-GISFW-WORK` | `INPUTDTLCLAUSE_FACIN` | [terverifikasi] pyInclude=8 | — |
| `InputDtlClause_FacIn_IsUW.xml` | `ASM-FW-GISFW-WORK` | `INPUTDTLCLAUSE_FACIN_ISUW` | [terverifikasi] pyInclude=8 | — |
| `InputDtlConveyance.xml` | `ASM-FW-GISFW-DATA-CONVEYANCE` | `INPUTDTLCONVEYANCE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlConveyance_IsUW.xml` | `ASM-FW-GISFW-DATA-CONVEYANCE` | `INPUTDTLCONVEYANCE_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-CONVEYANCE / INPUTDTLCONVEYANCE` |
| `InputDtlCoverage_FacIn.xml` | `ASM-FW-GISFW-DATA-CARGO` | `INPUTDTLCOVERAGE_FACIN` | [terverifikasi] pyInclude=1 | — |
| `InputDtlCoverage_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-CARGO` | `INPUTDTLCOVERAGE_FACIN_ISUW` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-DATA-CARGO / INPUTDTLCOVERAGE_FACIN` |
| `InputDtlGoods.xml` | `ASM-FW-GISFW-DATA-GOOD` | `INPUTDTLGOODS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlGoods_IsUW.xml` | `ASM-FW-GISFW-DATA-GOOD` | `INPUTDTLGOODS_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-GOOD / INPUTDTLGOODS` |
| `InputDtlLayerSpreading_FacIn.xml` | `ASM-FW-GISFW-DATA-SPREADINGRISK` | `INPUTDTLLAYERSPREADING_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlMaintenance.xml` | `DATA-ADDRESS` | `INPUTDTLMAINTENANCE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlMemorandumAccept.xml` | `ASM-FW-GISFW-DATA-MEMORANDUMACCEPT` | `INPUTDTLMEMORANDUMACCEPT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlMemorandumAccept_IsUW.xml` | `ASM-FW-GISFW-DATA-MEMORANDUMACCEPT` | `INPUTDTLMEMORANDUMACCEPT_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-MEMORANDUMACCEPT / INPUTDTLMEMORANDUMACCEPT` |
| `InputDtlObjFire.xml` | `ASM-FW-GISFW-DATA-PROPERTY` | `INPUTDTLOBJFIRE` | [terverifikasi] pyInclude=2 | — |
| `InputDtlObjFire_IsUW.xml` | `ASM-FW-GISFW-DATA-PROPERTY` | `INPUTDTLOBJFIRE_ISUW` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-PROPERTY / INPUTDTLOBJFIRE` |
| `InputDtlObjectAneka_FacIn.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `INPUTDTLOBJECTANEKA_FACIN` | [terverifikasi] pyInclude=2 | — |
| `InputDtlObjectAneka_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `INPUTDTLOBJECTANEKA_FACIN_ISUW` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-ANEKA / INPUTDTLOBJECTANEKA_FACIN` |
| `InputDtlObjectLocationHull_GISFW.xml` | `DATA-ADDRESS` | `INPUTDTLOBJECTLOCATIONHULL_GISFW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-ADDRESS / INPUTDTLOBJECTLOCATIONHULL_GCNM` |
| `InputDtlObjectLocation_FacIn.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `INPUTDTLOBJECTLOCATION_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlObjectLocation_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `INPUTDTLOBJECTLOCATION_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / INPUTDTLOBJECTLOCATION_FACIN` |
| `InputDtlObjectLocation_GISFW.xml` | `DATA-ADDRESS` | `INPUTDTLOBJECTLOCATION_GISFW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-ADDRESS / INPUTDTLOBJECTLOCATION_GCNM` |
| `InputDtlObject_FacIn.xml` | `ASM-FW-GISFW-WORK` | `INPUTDTLOBJECT_FACIN` | [terverifikasi] pyInclude=11 | — |
| `InputDtlObject_FacIn_IsUW.xml` | `ASM-FW-GISFW-WORK` | `INPUTDTLOBJECT_FACIN_ISUW` | [terverifikasi] pyInclude=8 | `ASM-FW-GISFW-WORK / INPUTDTLOBJECT_FACIN` |
| `InputDtlOtherObjectAneka.xml` | `ASM-FW-GISFW-WORK` | `INPUTDTLOTHEROBJECTANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlOtherObjectAneka_FacIn.xml` | `ASM-FW-GISFW-WORK` | `INPUTDTLOTHEROBJECTANEKA_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlOtherObjectAneka_FacIn_ISUW.xml` | `ASM-FW-GISFW-WORK` | `INPUTDTLOTHEROBJECTANEKA_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / INPUTDTLOTHEROBJECTANEKA_FACIN` |
| `InputDtlOtherSchedule.xml` | `ASM-FW-GISFW-DATA-ANEKAOTHERSCHEDULE` | `INPUTDTLOTHERSCHEDULE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlParticipant.xml` | `ASM-FW-GISFW-DATA-ANEKAPARTICIPANT` | `INPUTDTLPARTICIPANT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlParticipantDM.xml` | `DATA-PARTY-PERSON` | `INPUTDTLPARTICIPANTDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlParticipantLife_FacIn.xml` | `DATA-PARTY-PERSON` | `INPUTDTLPARTICIPANTLIFE_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlParticipantLife_FacIn_IsUW.xml` | `DATA-PARTY-PERSON` | `INPUTDTLPARTICIPANTLIFE_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / INPUTDTLPARTICIPANTLIFE_FACIN` |
| `InputDtlParticipantPA_FacIn.xml` | `DATA-PARTY-PERSON` | `INPUTDTLPARTICIPANTPA_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlParticipantPA_FacIn_IsUW.xml` | `DATA-PARTY-PERSON` | `INPUTDTLPARTICIPANTPA_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / INPUTDTLPARTICIPANTPA_FACIN` |
| `InputDtlParticipantTravel.xml` | `DATA-PARTY-PERSON` | `INPUTDTLPARTICIPANTTRAVEL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlParticipantTravel_GISFW.xml` | `DATA-PARTY-PERSON` | `INPUTDTLPARTICIPANTTRAVEL_GISFW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / INPUTDTLPARTICIPANTTRAVEL_GCNM` |
| `InputDtlParticipantTravel_IsUW.xml` | `DATA-PARTY-PERSON` | `INPUTDTLPARTICIPANTTRAVEL_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / INPUTDTLPARTICIPANTTRAVEL` |
| `InputDtlPaymentFacRetro.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | `INPUTDTLPAYMENTFACRETRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlPaymentFacRetro_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | `INPUTDTLPAYMENTFACRETRO_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlPayment_FacIn.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | `INPUTDTLPAYMENT_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlPayment_FacInLife.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | `INPUTDTLPAYMENT_FACINLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlPayment_FacInLife_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | `INPUTDTLPAYMENT_FACINLIFE_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY / INPUTDTLPAYMENT_FACINLIFE` |
| `InputDtlPayment_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | `INPUTDTLPAYMENT_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY / INPUTDTLPAYMENT_FACIN` |
| `InputDtlSpreadingCommisionCoverage_FacIn.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDTLSPREADINGCOMMISIONCOVERAGE_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLSPREADINGCOVERAGE_FACIN` |
| `InputDtlSpreadingCoverage_FacIn.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDTLSPREADINGCOVERAGE_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlSpreadingCoverage_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDTLSPREADINGCOVERAGE_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLSPREADINGCOVERAGE_FACIN` |
| `InputDtlSpreadingPA_FacIn.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDTLSPREADINGPA_FACIN` | [terverifikasi] pyInclude=2 | — |
| `InputDtlSpreadingPA_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDTLSPREADINGPA_FACIN_ISUW` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLSPREADINGPA_FACIN` |
| `InputDtlSpreadingTravel_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `INPUTDTLSPREADINGTRAVEL_FACIN_ISUW` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLSPREADINGTRAVEL_FACIN` |
| `InputDtlTrading.xml` | `ASM-FW-GISFW-DATA-TRADING` | `INPUTDTLTRADING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlTrading_IsUW.xml` | `ASM-FW-GISFW-DATA-TRADING` | `INPUTDTLTRADING_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TRADING / INPUTDTLTRADING` |
| `InputEndorsement.xml` | `ASM-FW-GISFW-WORK` | `INPUTENDORSEMENT` | [terverifikasi] pyInclude=10 | `ASM-FW-GISFW-WORK / INPUTINWARDFACULTATIVE` |
| `InputEndorsementDtl.xml` | `ASM-FW-GISFW-WORK` | `INPUTENDORSEMENTDTL` | [terverifikasi] pyInclude=31 | — |
| `InputEndorsementDtl_IsUW.xml` | `ASM-FW-GISFW-WORK` | `INPUTENDORSEMENTDTL_ISUW` | [terverifikasi] pyInclude=27 | — |
| `InputFEA.xml` | `ASM-FW-GISFW-DATA-FEA` | `INPUTFEA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputFEA_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-OFFERFEALIST` | `INPUTFEA_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-OFFERFEALIST / INPUTFEA` |
| `InputFacOffer.xml` | `ASM-FW-GISFW-WORK` | `INPUTFACOFFER` | [terverifikasi] pyInclude=4 | — |
| `InputFacOffer_IsUW.xml` | `ASM-FW-GISFW-WORK` | `INPUTFACOFFER_ISUW` | [terverifikasi] pyInclude=6 | — |
| `InputHistoricalSurveyReportDtl.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `INPUTHISTORICALSURVEYREPORTDTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputHistoricalSurveyReportDtlUW.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `INPUTHISTORICALSURVEYREPORTDTLUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-QUOTATION / INPUTHISTORICALSURVEYREPORTDTL` |
| `InputInwardFacultative.xml` | `ASM-FW-GISFW-WORK` | `INPUTINWARDFACULTATIVE` | [terverifikasi] pyInclude=12 | — |
| `InputInwardFacultativeDtl.xml` | `ASM-FW-GISFW-WORK` | `INPUTINWARDFACULTATIVEDTL` | [terverifikasi] pyInclude=38 | — |
| `InputInwardFacultativeDtl_IsUW.xml` | `ASM-FW-GISFW-WORK` | `INPUTINWARDFACULTATIVEDTL_ISUW` | [terverifikasi] pyInclude=30 | — |
| `InputInwardFacultativeSuggest.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `INPUTINWARDFACULTATIVESUGGEST` | [terverifikasi] pyInclude=1 | — |
| `InputInwardFacultative_2.xml` | `ASM-FW-GISFW-WORK` | `INPUTINWARDFACULTATIVE_2` | [terverifikasi] pyInclude=12 | `ASM-FW-GISFW-WORK / INPUTINWARDFACULTATIVE` |
| `InputLossExperiance.xml` | `ASM-FW-GISFW-DATA-CAUSEOFLOSS` | `INPUTLOSSEXPERIANCE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputLossRecord_Sec.xml` | `ASM-FW-GISFW-DATA-CAUSEOFLOSS` | `INPUTLOSSRECORD_SEC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-CAUSEOFLOSS / INPUTCAUSEOFLOSS_FACIN` |
| `InputMarketingOfficer.xml` | `@BASECLASS` | `INPUTMARKETINGOFFICER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputNation.xml` | `@BASECLASS` | `INPUTNATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputObjItem.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `INPUTOBJITEM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputObjectDtlFire.xml` | `ASM-FW-GISFW-DATA-PROPERTY` | `INPUTOBJECTDTLFIRE` | [terverifikasi] pyInclude=2 | — |
| `InputObjectDtlFire_GCNM.xml` | `ASM-FW-GISFW-DATA-PROPERTY` | `INPUTOBJECTDTLFIRE_GCNM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTY / INPUTOBJECTDTLFIRE` |
| `InputObjectDtlFire_IsUW.xml` | `ASM-FW-GISFW-DATA-PROPERTY` | `INPUTOBJECTDTLFIRE_ISUW` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-PROPERTY / INPUTOBJECTDTLFIRE` |
| `InputObjectSubContract_FacIn.xml` | `ASM-FW-GISFW-WORK` | `INPUTOBJECTSUBCONTRACT_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputObjectSubContract_FacIn_IsUW.xml` | `ASM-FW-GISFW-WORK` | `INPUTOBJECTSUBCONTRACT_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / INPUTOBJECTSUBCONTRACT_FACIN` |
| `InputOfferFacInLossRatio.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `INPUTOFFERFACINLOSSRATIO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputOkupasi.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `INPUTOKUPASI` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputOkupasiAneka_FacIn.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `INPUTOKUPASIANEKA_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputOkupasiAneka_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `INPUTOKUPASIANEKA_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OCCUPATION / INPUTOKUPASIANEKA_FACIN` |
| `InputOkupasiAneka_GCNM.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `INPUTOKUPASIANEKA_GCNM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OCCUPATION / INPUTOKUPASIANEKA` |
| `InputOkupasiHull_GCNM.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `INPUTOKUPASIHULL_GCNM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OCCUPATION / INPUTOKUPASIANEKA_GCNM` |
| `InputOtherObjectAneka_FacIn.xml` | `ASM-FW-GISFW-WORK` | `INPUTOTHEROBJECTANEKA_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputOtherObjectAneka_FacIn_ISUW.xml` | `ASM-FW-GISFW-WORK` | `INPUTOTHEROBJECTANEKA_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / INPUTOTHEROBJECTANEKA_FACIN` |
| `InputOtherObjectAneka_UW.xml` | `ASM-FW-GISFW-WORK` | `INPUTOTHEROBJECTANEKA_UW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / INPUTOTHEROBJECTANEKA` |
| `InputOtherSchedule.xml` | `ASM-FW-GISFW-WORK` | `INPUTOTHERSCHEDULE` | [terverifikasi] pyInclude=1 | — |
| `InputOutgo.xml` | `ASM-FW-GISFW-DATA-OUTGO` | `INPUTOUTGO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputPerson.xml` | `ASM-FW-GISFW-WORK` | `INPUTPERSON` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-NB / INPUTPERSON` |
| `InputProvince.xml` | `@BASECLASS` | `INPUTPROVINCE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputRW.xml` | `@BASECLASS` | `INPUTRW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputRiskAddress.xml` | `@BASECLASS` | `INPUTRISKADDRESS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputSpreadingLife_FacIn.xml` | `DATA-PARTY-PERSON` | `INPUTSPREADINGLIFE_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputSpreadingLife_FacIn_IsUW.xml` | `DATA-PARTY-PERSON` | `INPUTSPREADINGLIFE_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / INPUTSPREADINGLIFE_FACIN` |
| `InputSpreadingPA_FacIn.xml` | `DATA-PARTY-PERSON` | `INPUTSPREADINGPA_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputSpreadingPA_FacIn_IsUW.xml` | `DATA-PARTY-PERSON` | `INPUTSPREADINGPA_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / INPUTSPREADINGPA_FACIN` |
| `InputSpreadingTravel_FacIn.xml` | `DATA-PARTY-PERSON` | `INPUTSPREADINGTRAVEL_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / INPUTSPREADINGPA_FACIN` |
| `InputSpreadingTravel_FacIn_IsUW.xml` | `DATA-PARTY-PERSON` | `INPUTSPREADINGTRAVEL_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / INPUTSPREADINGTRAVEL_FACIN` |
| `InstallmentList.xml` | `ASM-FW-GISFW-DATA-INSTALLMENT` | `INSTALLMENTLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Installments_ReadOnly.xml` | `ASM-FW-GISFW-DATA-TREATYININSTALLMENT` | `INSTALLMENTS_READONLY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYININSTALLMENT / INSTALLMENTS_REALISASI` |
| `InsuredNameTable.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `INSUREDNAMETABLE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InwardFacIn.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `INWARDFACIN` | [terverifikasi] pyInclude=3 | — |
| `LayerListDtl.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `LAYERLISTDTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `LayerListDtl_IsUW.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `LAYERLISTDTL_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / LAYERLISTDTL` |
| `LimitTreaty.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `LIMITTREATY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ListPayment_SC.xml` | `ASM-FW-GISFW-WORK` | `LISTPAYMENT_SC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-NB / LISTPAYMENT_SC` |
| `ListSuggest.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `LISTSUGGEST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `LocationDetail.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `LOCATIONDETAIL` | [terverifikasi] pyInclude=3 | — |
| `LocationPropertyFacOut.xml` | `ASM-FW-GISFW-DATA-PROPERTY` | `LOCATIONPROPERTYFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `MarineCargoDtl.xml` | `ASM-FW-GISFW-DATA-CARGO` | `MARINECARGODTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `MarineCargoDtl_IsUW.xml` | `ASM-FW-GISFW-DATA-CARGO` | `MARINECARGODTL_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ObjectDetails.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `OBJECTDETAILS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ObjectDetails_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `OBJECTDETAILS_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / OBJECTDETAILS` |
| `ObjectDtlAneka_FacIn.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `OBJECTDTLANEKA_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ObjectItemSummary.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `OBJECTITEMSUMMARY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ObjectList.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `OBJECTLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ObjectList_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `OBJECTLIST_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN / OBJECTLIST` |
| `ObjectList_LossRecord.xml` | `ASM-FW-GISFW-WORK` | `OBJECTLIST_LOSSRECORD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ObjectOccupation_FacIn.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `OBJECTOCCUPATION_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `OccupationDetail.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `OCCUPATIONDETAIL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `OccupationFacOut.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `OCCUPATIONFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `OccupationFacOut_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `OCCUPATIONFACOUT_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / OCCUPATIONFACOUT` |
| `OccupationItemFacIn_Section.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `OCCUPATIONITEMFACIN_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `OccupationItemFacIn_Section_IsUW.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `OCCUPATIONITEMFACIN_SECTION_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OCCUPATION / OCCUPATIONITEMFACIN_SECTION` |
| `OccupationList.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `OCCUPATIONLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `OccupationList_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `OCCUPATIONLIST_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / OCCUPATIONLIST` |
| `OccupationPropertyFacOut.xml` | `ASM-FW-GISFW-DATA-PROPERTY` | `OCCUPATIONPROPERTYFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `OfferFacIn_NusaRe.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `OFFERFACIN_NUSARE` | [terverifikasi] pyInclude=2 | — |
| `OfferFacIn_NusaRe_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `OFFERFACIN_NUSARE_ISUW` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-OFFERFACIN / OFFERFACIN_NUSARE` |
| `OfferFacinLifeNotification.xml` | `ASM-FW-GISFW-WORK` | `OFFERFACINLIFENOTIFICATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ACCEPTNOTIFICATIONEDM` |
| `OfferStatusFacOut.xml` | `ASM-FW-GISFW-WORK` | `OFFERSTATUSFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `OldDataEndorsementDtl.xml` | `ASM-FW-GISFW-WORK` | `OLDDATAENDORSEMENTDTL` | [terverifikasi] pyInclude=21 | `ASM-FW-GISFW-WORK / INPUTENDORSEMENTDTL` |
| `OldOfferStatusFacOut.xml` | `ASM-FW-GISFW-WORK` | `OLDOFFERSTATUSFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / OFFERSTATUSFACOUT` |
| `PASummarySection.xml` | `ASM-FW-GISFW-WORK` | `PASUMMARYSECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / GROWINGTREESSUMMARYSECTION` |
| `PaymentCurrencyList.xml` | `ASM-FW-GISFW-WORK` | `PAYMENTCURRENCYLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PaymentCurrencyListLife.xml` | `ASM-FW-GISFW-WORK` | `PAYMENTCURRENCYLISTLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PaymentCurrencyListLife_IsUW.xml` | `ASM-FW-GISFW-WORK` | `PAYMENTCURRENCYLISTLIFE_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / PAYMENTCURRENCYLISTLIFE` |
| `PaymentCurrencyList_IsUW.xml` | `ASM-FW-GISFW-WORK` | `PAYMENTCURRENCYLIST_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / PAYMENTCURRENCYLIST` |
| `Periode.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `PERIODE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PeriodeEndorsement.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `PERIODEENDORSEMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PeriodePolicy.xml` | `ASM-FW-GISFW-DATA-POLICY` | `PERIODEPOLICY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Periode_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `PERIODE_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN / PERIODE` |
| `PersonCoverage.xml` | `DATA-PARTY-PERSON` | `PERSONCOVERAGE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PersonCoverageRO.xml` | `DATA-PARTY-PERSON` | `PERSONCOVERAGERO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PersonListDetail.xml` | `DATA-PARTY-PERSON` | `PERSONLISTDETAIL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PersonPlanAndDeduct.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `PERSONPLANANDDEDUCT` | [terverifikasi] pyInclude=2 | — |
| `PlanContains_UWRO.xml` | `ASM-FW-GISFW-DATA-PLAN` | `PLANCONTAINS_UWRO` | [terverifikasi] pyInclude=1 | — |
| `PlanListRO.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `PLANLISTRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PolicyTreatyInDeclineConfirm.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `POLICYTREATYINDECLINECONFIRM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PrintRISlip.xml` | `ASM-FW-GISFW-WORK` | `PRINTRISLIP` | [terverifikasi] pyInclude=6 | `ASM-FW-GISFW-WORK / FACOUTMEMOPLACING` |
| `PrintRISlip_Endorsement.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `PRINTRISLIP_ENDORSEMENT` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-FACOFFER / PRINTRISLIP` |
| `PrintRISlips_Endorsement.xml` | `ASM-FW-GISFW-WORK` | `PRINTRISLIPS_ENDORSEMENT` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-FACOFFER / PRINTRISLIP_ENDORSEMENT` |
| `Property.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `PROPERTY` | [terverifikasi] pyInclude=16 | — |
| `PropertyItemFacIn_Section.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `PROPERTYITEMFACIN_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PropertyItemFacIn_Section_IsUW.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `PROPERTYITEMFACIN_SECTION_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTYITEM / PROPERTYITEMFACIN_SECTION` |
| `PropertyItemList.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `PROPERTYITEMLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PropertyItemListCoverage.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `PROPERTYITEMLISTCOVERAGE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / PROPERTYITEMLIST` |
| `PropertyItemListCoverageCommision.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `PROPERTYITEMLISTCOVERAGECOMMISION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / PROPERTYITEMLISTCOVERAGESPREADING` |
| `PropertyItemListCoverageSpreading.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `PROPERTYITEMLISTCOVERAGESPREADING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / PROPERTYITEMLISTCOVERAGE` |
| `PropertyItemListCoverageSpreading_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `PROPERTYITEMLISTCOVERAGESPREADING_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / PROPERTYITEMLISTCOVERAGESPREADING` |
| `PropertyItemListCoverage_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `PROPERTYITEMLISTCOVERAGE_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PropertyItemListFacOut.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `PROPERTYITEMLISTFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / PROPERTYITEMLISTCOVERAGE` |
| `PropertyItemListFacOut_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `PROPERTYITEMLISTFACOUT_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / PROPERTYITEMLISTFACOUT` |
| `PropertyItemList_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `PROPERTYITEMLIST_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / PROPERTYITEMLIST` |
| `Property_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `PROPERTY_ISUW` | [terverifikasi] pyInclude=16 | — |
| `ProtectCurrency.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `PROTECTCURRENCY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ProtectCurrencyCargo.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `PROTECTCURRENCYCARGO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ProtectCurrencyCargo_IsUW.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `PROTECTCURRENCYCARGO_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ProtectCurrency_IsUW.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `PROTECTCURRENCY_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-DEDUCTIBLE / PROTECTCURRENCY` |
| `ReasViewAttachment.xml` | `ASM-FW-GISFW-INT` | `REASVIEWATTACHMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `RejectNotificationSection.xml` | `ASM-FW-GISFW-WORK` | `REJECTNOTIFICATIONSECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `RemarksFacLetter.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `REMARKSFACLETTER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ReviseNotification.xml` | `ASM-FW-GISFW-WORK` | `REVISENOTIFICATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `RiskAccumulationReport.xml` | `DATA-PORTAL` | `RISKACCUMULATIONREPORT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `RiskAround.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `RISKAROUND` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `RiskAround_IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `RISKAROUND_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / RISKAROUND` |
| `SFAPortalOpportunitiesHeader.xml` | `PEGACRM-PORTAL` | `SFAPORTALOPPORTUNITIESHEADER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SFAPortal_Opportunities.xml` | `DATA-PORTAL` | `SFAPORTAL_OPPORTUNITIES` | [terverifikasi] pyInclude=4 | — |
| `SFAPortal_OpportunitiesList.xml` | `DATA-PORTAL` | `SFAPORTAL_OPPORTUNITIESLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SFAPortal_OpportunitiesList_Header.xml` | `DATA-PORTAL` | `SFAPORTAL_OPPORTUNITIESLIST_HEADER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ScoringRisk.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SCORINGRISK` | [terverifikasi] pyInclude=4 | — |
| `ScoringRiskForm1.xml` | `ASM-FW-GISFW-DATA-SCORINGRISK` | `SCORINGRISKFORM1` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ScoringRiskForm2.xml` | `ASM-FW-GISFW-DATA-SCORINGRISK` | `SCORINGRISKFORM2` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-SCORINGRISK / SCORINGRISKFORM1` |
| `ScoringRiskForm3.xml` | `ASM-FW-GISFW-DATA-SCORINGRISK` | `SCORINGRISKFORM3` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-SCORINGRISK / SCORINGRISKFORM2` |
| `ScoringRiskForm4.xml` | `ASM-FW-GISFW-DATA-SCORINGRISK` | `SCORINGRISKFORM4` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-SCORINGRISK / SCORINGRISKFORM1` |
| `SearchRiskAccumCov.xml` | `DATA-PORTAL` | `SEARCHRISKACCUMCOV` | [terverifikasi] pyInclude=2 | `DATA-PORTAL / SEARCHRISKACCUMULATION` |
| `SelClauseList.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SELCLAUSELIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PLAN / SELCLAUSELIST` |
| `SelectAgent.xml` | `ASM-FW-GISFW-INT-LLOYDAGENT` | `SELECTAGENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SelectCoverage.xml` | `ASM-FW-GISFW-INT-CONDITION` | `SELECTCOVERAGE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SelectPlanListRO.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `SELECTPLANLISTRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SelectShip.xml` | `ASM-FW-GISFW-INT-SHIP` | `SELECTSHIP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SelectWarrantyList.xml` | `ASM-FW-GISFW-WORK` | `SELECTWARRANTYLIST` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-WORK-NB / SELECTWARRANTYLIST` |
| `ShowCedingCoList.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SHOWCEDINGCOLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowCoverageFacOut.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `SHOWCOVERAGEFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowCoverageFacOut_IsUW.xml` | `ASM-FW-GISFW-DATA-CARGO` | `SHOWCOVERAGEFACOUT_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-CARGO / SHOWCOVERAGEFACOUT` |
| `ShowCoveragePAFacOutIsUW_Dtl.xml` | `DATA-PARTY-PERSON` | `SHOWCOVERAGEPAFACOUTISUW_DTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / SHOWCOVERAGEPAFACOUT_DTL` |
| `ShowCoveragePAFacOut_Dtl.xml` | `DATA-PARTY-PERSON` | `SHOWCOVERAGEPAFACOUT_DTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowCoverageTravelFacOut_Dtl.xml` | `DATA-PARTY-PERSON` | `SHOWCOVERAGETRAVELFACOUT_DTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / SHOWCOVERAGEPAFACOUT_DTL` |
| `ShowCoverageTravelFacOut_Dtl_IsUW.xml` | `DATA-PARTY-PERSON` | `SHOWCOVERAGETRAVELFACOUT_DTL_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / SHOWCOVERAGETRAVELFACOUT_DTL` |
| `ShowObjectFacOut.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `SHOWOBJECTFACOUT` | [terverifikasi] pyInclude=2 | — |
| `ShowObjectFacOut_IsUW.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `SHOWOBJECTFACOUT_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowPolicyNoTreaty_SC.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SHOWPOLICYNOTREATY_SC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowPolis_sc.xml` | `ASM-FW-GISFW-WORK` | `SHOWPOLIS_SC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SourceHierarki.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `SOURCEHIERARKI` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SpreadingCoverageAneka.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `SPREADINGCOVERAGEANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SpreadingCoverageAneka_IsUW.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `SPREADINGCOVERAGEANEKA_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-ANEKA / SPREADINGCOVERAGEANEKA` |
| `SpreadingItem.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SPREADINGITEM` | [terverifikasi] pyInclude=2 | — |
| `SpreadingItemAneka.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SPREADINGITEMANEKA` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-COVERAGE / SPREADINGITEM` |
| `SpreadingItemAneka_IsUW.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SPREADINGITEMANEKA_ISUW` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-COVERAGE / SPREADINGITEMANEKA` |
| `SpreadingItemLayer.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SPREADINGITEMLAYER` | [terverifikasi] pyInclude=2 | — |
| `SpreadingItemLayerCommision.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SPREADINGITEMLAYERCOMMISION` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-COVERAGE / SPREADINGITEMLAYER` |
| `SpreadingItemLayer_IsUW.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SPREADINGITEMLAYER_ISUW` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-COVERAGE / SPREADINGITEMLAYER` |
| `SpreadingItem_IsUW.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `SPREADINGITEM_ISUW` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-COVERAGE / SPREADINGITEM` |
| `SpreadingObjectAneka.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `SPREADINGOBJECTANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SpreadingObjectAneka_IsUW.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `SPREADINGOBJECTANEKA_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OCCUPATION / SPREADINGOBJECTANEKA` |
| `SpreadingRiskList.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SPREADINGRISKLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SummaryCoverage_Section.xml` | `ASM-FW-GISFW-WORK` | `SUMMARYCOVERAGE_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / COVERAGESPREADING_SECTION` |
| `SummaryLossRecord_Section.xml` | `ASM-FW-GISFW-WORK` | `SUMMARYLOSSRECORD_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / SUMMARYCOVERAGE_SECTION` |
| `SummaryRiskAccumulation.xml` | `DATA-PORTAL` | `SUMMARYRISKACCUMULATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SummarySpreading_Section.xml` | `ASM-FW-GISFW-WORK` | `SUMMARYSPREADING_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TSIForFacOut.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `TSIFORFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TSIPropertyFacOut.xml` | `ASM-FW-GISFW-DATA-PROPERTY` | `TSIPROPERTYFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TableInwardScale.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `TABLEINWARDSCALE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TableOfLimit.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `TABLEOFLIMIT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TotalAccumulationDtl.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `TOTALACCUMULATIONDTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TotalAccumulation_FacIn.xml` | `ASM-FW-GISFW-WORK` | `TOTALACCUMULATION_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `UploadMember.xml` | `ASM-FW-GISFW-WORK-NB` | `UPLOADMEMBER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `VehicleGrid.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `VEHICLEGRID` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `VehicleGrid_FacIn.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `VEHICLEGRID_FACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `VehicleGrid_FacIn_IsUW.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `VEHICLEGRID_FACIN_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-VEHICLE / VEHICLEGRID_FACIN` |
| `ViewAlamatKirim.xml` | `DATA-ADDRESS` | `VIEWALAMATKIRIM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewCSVAneka.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `VIEWCSVANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewCSVCredit.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `VIEWCSVCREDIT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewCSVMarine.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `VIEWCSVMARINE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OFFERFACIN / VIEWCSVCREDIT` |
| `ViewCSVResult_TableB2B.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `VIEWCSVRESULT_TABLEB2B` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewCheckListOffer.xml` | `ASM-FW-GISFW-WORK` | `VIEWCHECKLISTOFFER` | [terverifikasi] pyInclude=6 | `ASM-FW-GISFW-WORK / VIEWCHECKLISTOFFERMEMOPLACING` |
| `ViewCheckListOfferFacOut.xml` | `ASM-FW-GISFW-WORK` | `VIEWCHECKLISTOFFERFACOUT` | [terverifikasi] pyInclude=6 | — |
| `ViewClaimList.xml` | `ASM-FW-GISFW-WORK` | `VIEWCLAIMLIST` | [terverifikasi] pyInclude=1 | — |
| `ViewClauseFire.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `VIEWCLAUSEFIRE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-CLAUSE / INPUTCLAUSEFIRE` |
| `ViewCoins.xml` | `ASM-FW-GISFW-DATA-COINS` | `VIEWCOINS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewCompanyInformationFacOut.xml` | `ASM-FW-GISFW-WORK` | `VIEWCOMPANYINFORMATIONFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewCoverage.xml` | `ASM-FW-GISFW-DATA-CARGO` | `VIEWCOVERAGE` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-DATA-VEHICLE / VIEWCOVERAGE` |
| `ViewCoverageAneka.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `VIEWCOVERAGEANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / INPUTCOVERAGEANEKA` |
| `ViewCoverageCargoFacOut.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `VIEWCOVERAGECARGOFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / INPUTCOVERAGECARGO` |
| `ViewCoverageFacOutShow.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `VIEWCOVERAGEFACOUTSHOW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewCoverageFacOutShow_isUW.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `VIEWCOVERAGEFACOUTSHOW_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTYITEM / VIEWCOVERAGEFACOUTSHOW` |
| `ViewCoverageFire.xml` | `ASM-FW-GISFW-DATA-PROPERTY` | `VIEWCOVERAGEFIRE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTY / INPUTCOVERAGEFIRE` |
| `ViewCoveragePA.xml` | `DATA-PARTY-PERSON` | `VIEWCOVERAGEPA` | [terverifikasi] pyInclude=1 | `DATA-PARTY-PERSON / INPUTCOVERAGEPA` |
| `ViewCoverageTravel.xml` | `DATA-PARTY-PERSON` | `VIEWCOVERAGETRAVEL` | [terverifikasi] pyInclude=1 | `DATA-PARTY-PERSON / INPUTCOVERAGETRAVEL` |
| `ViewDataPolicy.xml` | `ASM-FW-GISFW-WORK` | `VIEWDATAPOLICY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewDataPolicy_IsUW.xml` | `ASM-FW-GISFW-WORK` | `VIEWDATAPOLICY_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / VIEWDATAPOLICY` |
| `ViewDeductible.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `VIEWDEDUCTIBLE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-DEDUCTIBLE / INPUTDEDUCTIBLE` |
| `ViewDeliveryAddressFacOut.xml` | `ASM-FW-GISFW-WORK` | `VIEWDELIVERYADDRESSFACOUT` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-DATA-POLICY / VIEWDELIVERYADDRESSFACOUT` |
| `ViewDetailCargo.xml` | `ASM-FW-GISFW-DATA-CARGO` | `VIEWDETAILCARGO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-CARGO / INPUTDTLCARGO` |
| `ViewDetailPayment_SC.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | `VIEWDETAILPAYMENT_SC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewDtlAddCargo.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `VIEWDTLADDCARGO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLADDCARGO` |
| `ViewDtlAdditional.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `VIEWDTLADDITIONAL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLADDITIONAL` |
| `ViewDtlCoverageFacOut.xml` | `ASM-FW-GISFW-WORK` | `VIEWDTLCOVERAGEFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / INPUTDTLCOVERAGE` |
| `ViewDtlDeductible.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `VIEWDTLDEDUCTIBLE` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLDEDUCTIBLE` |
| `ViewDtlMemorandumAccept.xml` | `ASM-FW-GISFW-DATA-MEMORANDUMACCEPT` | `VIEWDTLMEMORANDUMACCEPT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-MEMORANDUMACCEPT / INPUTDTLMEMORANDUMACCEPT` |
| `ViewDtlObjectAneka.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `VIEWDTLOBJECTANEKA` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-ANEKA / INPUTDTLOBJECTANEKA` |
| `ViewDtlObjectLocation.xml` | `DATA-ADDRESS` | `VIEWDTLOBJECTLOCATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-ADDRESS / INPUTDTLOBJECTLOCATION` |
| `ViewDtlOtherObjectAneka.xml` | `ASM-FW-GISFW-WORK` | `VIEWDTLOTHEROBJECTANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / INPUTDTLOTHEROBJECTANEKA` |
| `ViewDtlOutgo.xml` | `ASM-FW-GISFW-WORK` | `VIEWDTLOUTGO` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-WORK / INPUTDTLOUTGO` |
| `ViewDtlParticipantPA.xml` | `DATA-PARTY-PERSON` | `VIEWDTLPARTICIPANTPA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / INPUTDTLPARTICIPANTPA` |
| `ViewGeneralPolisFacOut.xml` | `ASM-FW-GISFW-WORK` | `VIEWGENERALPOLISFACOUT` | [terverifikasi] pyInclude=6 | — |
| `ViewIndemnity_Section.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `VIEWINDEMNITY_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewIndividualInformationFacOut.xml` | `ASM-FW-GISFW-WORK` | `VIEWINDIVIDUALINFORMATIONFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewInwardFacultativeDtl_IsUW.xml` | `ASM-FW-GISFW-WORK` | `VIEWINWARDFACULTATIVEDTL_ISUW` | [terverifikasi] pyInclude=21 | — |
| `ViewObjectAnekaFacOut.xml` | `ASM-FW-GISFW-WORK` | `VIEWOBJECTANEKAFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-POLICY / VIEWOBJECTANEKAFACOUT` |
| `ViewObjectDtlAneka.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `VIEWOBJECTDTLANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OCCUPATION / OBJECTDTLANEKA` |
| `ViewObjectFireFacOut.xml` | `ASM-FW-GISFW-WORK` | `VIEWOBJECTFIREFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-POLICY / VIEWOBJECTFIREFACOUT` |
| `ViewObjectMarineCargoFacOut.xml` | `ASM-FW-GISFW-WORK` | `VIEWOBJECTMARINECARGOFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-POLICY / VIEWOBJECTMARINECARGOFACOUT` |
| `ViewObjectMarineHullFacOut.xml` | `ASM-FW-GISFW-WORK` | `VIEWOBJECTMARINEHULLFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-POLICY / VIEWOBJECTMARINEHULLFACOUT` |
| `ViewObjectOccupation.xml` | `DATA-ADDRESS` | `VIEWOBJECTOCCUPATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-ADDRESS / OBJECTOCCUPATION` |
| `ViewObjectPAFacOut.xml` | `ASM-FW-GISFW-WORK` | `VIEWOBJECTPAFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-POLICY / VIEWOBJECTPAFACOUT` |
| `ViewObjectSavior.xml` | `ASM-FW-GISFW-WORK` | `VIEWOBJECTSAVIOR` | [terverifikasi] pyInclude=16 | — |
| `ViewObjectSubContract.xml` | `ASM-FW-GISFW-WORK` | `VIEWOBJECTSUBCONTRACT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / INPUTOBJECTSUBCONTRACT` |
| `ViewObjectTravelFacOut.xml` | `ASM-FW-GISFW-WORK` | `VIEWOBJECTTRAVELFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-POLICY / VIEWOBJECTTRAVELFACOUT` |
| `ViewObjectVehicleFacOut.xml` | `ASM-FW-GISFW-WORK` | `VIEWOBJECTVEHICLEFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-POLICY / VIEWOBJECTVEHICLEFACOUT` |
| `ViewOkupasiAneka.xml` | `ASM-FW-GISFW-DATA-OCCUPATION` | `VIEWOKUPASIANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OCCUPATION / INPUTOKUPASIANEKA` |
| `ViewOldDataPayment_IsUW.xml` | `ASM-FW-GISFW-WORK` | `VIEWOLDDATAPAYMENT_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / OLDDATAPAYMENT_ISUW` |
| `ViewOldDeductible.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `VIEWOLDDEDUCTIBLE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewOpenFollowingPolicy.xml` | `ASM-FW-GISFW-WORK` | `VIEWOPENFOLLOWINGPOLICY` | [terverifikasi] pyInclude=7 | — |
| `ViewOutgo.xml` | `ASM-FW-GISFW-DATA-OUTGO` | `VIEWOUTGO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-OUTGO / INPUTOUTGO` |
| `ViewPerilsFacOut.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `VIEWPERILSFACOUT` | [terverifikasi] pyInclude=1 | — |
| `ViewPolicyListRO.xml` | `ASM-FW-GISFW-DATA-POLICY` | `VIEWPOLICYLISTRO` | [terverifikasi] pyInclude=1 | — |
| `ViewTabGroupPolisFacOutSavior.xml` | `ASM-FW-GISFW-WORK` | `VIEWTABGROUPPOLISFACOUTSAVIOR` | [terverifikasi] pyInclude=22 | `ASM-FW-GISFW-WORK / VIEWTABGROUPPOLISFACOUT` |
| `ViewVehicleGridFacOut.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `VIEWVEHICLEGRIDFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-VEHICLE / VEHICLEGRID` |
| `WarrantyList.xml` | `ASM-FW-GISFW-WORK` | `WARRANTYLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-NB / WARRANTYLIST` |
| `addDeductible.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `ADDDEDUCTIBLE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `addDeductibleCargo.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `ADDDEDUCTIBLECARGO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `addDeductibleCargo_IsUW.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `ADDDEDUCTIBLECARGO_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `addDeductible_IsUW.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `ADDDEDUCTIBLE_ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-DEDUCTIBLE / ADDDEDUCTIBLE` |
| `additemButton.xml` | `@BASECLASS` | `ADDITEMBUTTON` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `confirmCreateAccumulation.xml` | `DATA-PARTY-PERSON` | `CONFIRMCREATEACCUMULATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `deleteItemButton.xml` | `@BASECLASS` | `DELETEITEMBUTTON` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `pyAttachmentScreen.xml` | `DATA-WORKATTACH-FILE` | `PYATTACHMENTSCREEN` | [terverifikasi] pyInclude=10 | — |

### SystemSettings — 1 rule (`RULE-ADMIN-SYSTEM-SETTINGS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `LinkService.xml` | `LINKSERVICE` | `LINKSERVICE` | [terverifikasi] pyPurpose=LinkService | — |

### When — 210 rule (`RULE-OBJ-WHEN`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `FlagOldData.xml` | `ASM-FW-GISFW-DATA-ACCESSORY` | `ISFLAGOLDDATA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-VEHICLE / ISFLAGOLDDATA` |
| `IsAddButton.xml` | `@BASECLASS` | `ISADDBUTTON` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / ISADDCITY` |
| `IsAddCoverage.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `ISADDCOVERAGE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / ISADDMEDEXPA` |
| `IsAddMedEx.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `ISADDMEDEX` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / ISADDPA` |
| `IsAddMedExPA.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `ISADDMEDEXPA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / ISADDMEDEX` |
| `IsAddPA.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `ISADDPA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / ISTJH` |
| `IsAddTJHPLL.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `ISADDTJHPLL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / ISPLL` |
| `IsAdjustableFlag.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `ISADJUSTABLEFLAG` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsAdmin.xml` | `ASM-FW-GISFW-WORK` | `ISADMIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsAneka.xml` | `ASM-FW-GISFW-DATA` | `ISANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsAnekaNotBonding.xml` | `ASM-FW-GISFW-WORK` | `ISANEKANOTBONDING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsAssociationPolicy.xml` | `ASM-FW-GISFW-WORK` | `ISASSOCIATIONPOLICY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISNOTINDIVIDUALPOLICY` |
| `IsAsuransiKredit.xml` | `ASM-FW-GISFW-DATA` | `ISASURANSIKREDIT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsBillboardNeon.xml` | `ASM-FW-GISFW-WORK` | `ISBILLBOARDNEON` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISBILLBOARDNEON` |
| `IsBoiler.xml` | `@BASECLASS` | `ISBOILER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsBonding.xml` | `ASM-FW-GISFW-WORK` | `ISBONDING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsBondingAndCustomBonds.xml` | `ASM-FW-GISFW-DATA` | `ISBONDINGANDCUSTOMBONDS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISBONDING` |
| `IsBondingKBG.xml` | `@BASECLASS` | `ISBONDINGKBG` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsBuilderRisk.xml` | `ASM-FW-GISFW-DATA` | `ISBUILDERRISK` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsButtonShow.xml` | `@BASECLASS` | `ISBUTTONSHOW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / ISCITYSHOW` |
| `IsCAR.xml` | `@BASECLASS` | `ISCAR` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCIS.xml` | `@BASECLASS` | `ISCIS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCIT.xml` | `@BASECLASS` | `ISCIT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCMI.xml` | `ASM-FW-GISFW-WORK` | `ISCMI` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCedingConfirmOffer.xml` | `ASM-FW-GISFW-WORK` | `ISCEDINGCONFIRMOFFER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCedingConfirmPolicy.xml` | `ASM-FW-GISFW-WORK` | `ISCEDINGCONFIRMPOLICY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsChekCity.xml` | `@BASECLASS` | `ISCHEKCITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsChekNation.xml` | `@BASECLASS` | `ISCHEKNATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsChekProvince.xml` | `@BASECLASS` | `ISCHEKPROVINCE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsChekRW.xml` | `@BASECLASS` | `ISCHEKRW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCivilEngineeringCompletedRisks.xml` | `ASM-FW-GISFW-DATA` | `ISCIVILENGINEERINGCOMPLETEDRISKS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISMARINEHULLOFFSHORE` |
| `IsClaim.xml` | `ASM-FW-GISFW-DATA` | `ISCLAIM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCoas.xml` | `ASM-FW-GISFW-WORK` | `ISCOAS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsComprehensiveMachineriesInsurance.xml` | `ASM-FW-GISFW-DATA` | `ISCOMPREHENSIVEMACHINERIESINSURANCE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISCIVILENGINEERINGCOMPLETEDRISKS` |
| `IsContractorsPlantMachinery.xml` | `@BASECLASS` | `ISCONTRACTORSPLANTMACHINERY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCorporate.xml` | `ASM-FW-GISFW-DATA-POLICY` | `ISCORPORATE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCreditBriguna.xml` | `ASM-FW-GISFW-DATA` | `ISCREDITBRIGUNA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISASURANSIKREDIT` |
| `IsCrime.xml` | `ASM-FW-GISFW-DATA` | `ISCRIME` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCustomBonds.xml` | `@BASECLASS` | `ISCUSTOMBONDS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsDM.xml` | `ASM-FW-GISFW-WORK` | `ISDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsDeclarationPolicy.xml` | `ASM-FW-GISFW-WORK` | `ISDECLARATIONPOLICY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsDeducType.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `ISDEDUCTYPE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsDescDeductType.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `ISDESCDEDUCTTYPE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsEDM.xml` | `ASM-FW-GISFW-WORK` | `ISEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsEDMRiSlip.xml` | `ASM-FW-GISFW-WORK` | `ISEDMRISLIP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsEar.xml` | `@BASECLASS` | `ISEAR` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsEdmAddObject.xml` | `ASM-FW-GISFW-WORK` | `ISEDMADDOBJECT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISEDMADDOBJECT` |
| `IsEdmAdjCurrency.xml` | `ASM-FW-GISFW-WORK` | `ISEDMADJCURRENCY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISEDMADJCURRENCY` |
| `IsEdmAdjInsured.xml` | `ASM-FW-GISFW-WORK` | `ISEDMADJINSURED` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISEDMADJINSURED` |
| `IsEdmAdjPeriod.xml` | `ASM-FW-GISFW-WORK` | `ISEDMADJPERIOD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISEDMADJPERIOD` |
| `IsEdmAdjRIC.xml` | `ASM-FW-GISFW-WORK` | `ISEDMADJRIC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-MSHIP / ISEDMADJRIC` |
| `IsEdmAdjRate.xml` | `ASM-FW-GISFW-WORK` | `ISEDMADJRATE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISEDMADJRATE` |
| `IsEdmAdjRefNo.xml` | `ASM-FW-GISFW-WORK` | `ISEDMADJREFNO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISEDMADJREFNO` |
| `IsEdmAdjShareCedant.xml` | `ASM-FW-GISFW-WORK` | `ISEDMADJSHARECEDANT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISEDMADJREFNO` |
| `IsEdmAdjSpreading.xml` | `ASM-FW-GISFW-WORK` | `ISEDMADJSPREADING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISEDMADJSPREADING` |
| `IsEdmAdjTSI.xml` | `ASM-FW-GISFW-WORK` | `ISEDMADJTSI` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISEDMADJTSI` |
| `IsEdmExtendPeriod.xml` | `ASM-FW-GISFW-WORK` | `ISEDMEXTENDPERIOD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISEDMEXTENDPERIOD` |
| `IsEdmPPNPPH.xml` | `ASM-FW-GISFW-WORK` | `ISEDMPPNPPH` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-MSHIP / ISEDMPPNPPH` |
| `IsEdmPerubahan.xml` | `ASM-FW-GISFW-WORK` | `ISEDMPERUBAHAN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsEngineering.xml` | `ASM-FW-GISFW-DATA` | `ISENGINEERING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsEnvironmental.xml` | `ASM-FW-GISFW-DATA` | `ISENVIRONMENTAL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISCRIME` |
| `IsErrorSpreading.xml` | `@BASECLASS` | `ISERRORSPREADING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsExclusion.xml` | `@BASECLASS` | `ISEXCLUSION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsFac.xml` | `ASM-FW-GISFW-WORK` | `ISFAC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsFacRetro.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `ISFACRETRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISUW` |
| `IsFacout.xml` | `ASM-FW-GISFW-DATA` | `ISFACOUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsFinanceInsurance.xml` | `ASM-FW-GISFW-DATA` | `ISFINANCEINSURANCE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsFire.xml` | `ASM-FW-GISFW-WORK` | `ISFIRE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsFireStyle1.xml` | `@BASECLASS` | `ISFIRESTYLE1` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISFIRESTYLE1` |
| `IsFireStyle2.xml` | `@BASECLASS` | `ISFIRESTYLE2` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISFIRESTYLE2` |
| `IsFlagDelete.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `ISFLAGDELETE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-SUBJECTTO / ISFLAGDELETE` |
| `IsFlagOldData.xml` | `ASM-FW-GISFW-DATA-ACCESSORY` | `ISFLAGOLDDATA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-VEHICLE / ISFLAGOLDDATA` |
| `IsFlagOnGoingPolicy.xml` | `ASM-FW-GISFW-WORK` | `ISFLAGONGOINGPOLICY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsFlagUW.xml` | `ASM-FW-GISFW-WORK` | `ISFLAGUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsGabungan.xml` | `ASM-FW-GISFW-WORK` | `ISGABUNGAN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsGlass.xml` | `@BASECLASS` | `ISGLASS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsGolfInsurance.xml` | `ASM-FW-GISFW-DATA` | `ISGOLFINSURANCE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsGroup.xml` | `ASM-FW-GISFW-WORK` | `ISGROUP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsGroupCreate.xml` | `ASM-FW-GISFW-WORK` | `ISGROUPCREATE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsGroupFac.xml` | `ASM-FW-GISFW-DATA` | `ISGROUPFAC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISGROUPMKTFAC` |
| `IsGroupUwFac.xml` | `ASM-FW-GISFW-WORK` | `ISGROUPUWFAC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISUWFAC` |
| `IsGrowingTrees.xml` | `@BASECLASS` | `ISGROWINGTREES` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsHE.xml` | `ASM-FW-GISFW-DATA` | `ISHE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsHealth.xml` | `ASM-FW-GISFW-WORK` | `ISHEALTH` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsHealthCorporate.xml` | `ASM-FW-GISFW-WORK` | `ISHEALTHCORPORATE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsHealthIndividu.xml` | `ASM-FW-GISFW-WORK` | `ISHEALTHINDIVIDU` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsIT.xml` | `DATA-PORTAL` | `ISIT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsIndividual.xml` | `ASM-FW-GISFW-WORK` | `ISINDIVIDUAL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsInputFacRetro.xml` | `ASM-FW-GISFW-WORK` | `ISINPUTFACRETRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISNOTPRINTRISLIP` |
| `IsKPR.xml` | `@BASECLASS` | `ISKPR` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISKPR` |
| `IsLandRig.xml` | `@BASECLASS` | `ISLANDRIG` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsLiability.xml` | `@BASECLASS` | `ISLIABILITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsLife.xml` | `ASM-FW-GISFW-DATA` | `ISLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / ISLIFE` |
| `IsLimitCreditCL.xml` | `ASM-FW-GISFW-WORK` | `ISLIMITCREDITCL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsLimitCreditNCL.xml` | `ASM-FW-GISFW-WORK` | `ISLIMITCREDITNCL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsLimitCustomBond.xml` | `ASM-FW-GISFW-WORK` | `ISLIMITCUSTOMBOND` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsLimitSBondKBG.xml` | `ASM-FW-GISFW-WORK` | `ISLIMITSBONDKBG` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsLimitTradeCredit.xml` | `ASM-FW-GISFW-WORK` | `ISLIMITTRADECREDIT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsMBD.xml` | `ASM-FW-GISFW-DATA` | `ISMBD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsMBU.xml` | `ASM-FW-GISFW-DATA` | `ISMBU` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsMBUCar.xml` | `ASM-FW-GISFW-DATA` | `ISMBUCAR` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsMaintenance.xml` | `@BASECLASS` | `ISMAINTENANCE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsMarineCargo.xml` | `ASM-FW-GISFW-WORK` | `ISMARINECARGO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsMarineHull.xml` | `ASM-FW-GISFW-DATA` | `ISMARINEHULL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsMarineHullOffshore.xml` | `ASM-FW-GISFW-DATA` | `ISMARINEHULLOFFSHORE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISWHETHERINDEX` |
| `IsMultiCOB.xml` | `ASM-FW-GISFW-WORK` | `ISMULTICOB` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsNB.xml` | `ASM-FW-GISFW-WORK` | `ISNB` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISNOTEDM` |
| `IsNoDeduc.xml` | `ASM-FW-GISFW-DATA-PREMIUMTEMPLATE` | `ISNODEDUC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsNonPropertyandNonEngineering.xml` | `ASM-FW-GISFW-WORK` | `ISNONPROPERTYANDNONENGINEERING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsNotAdmin.xml` | `@BASECLASS` | `ISNOTADMIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsNotDirect.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `ISNOTDIRECT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsNotEDM.xml` | `ASM-FW-GISFW-DATA-VEHICLE` | `ISNOTEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-VEHICLE / ISEDM` |
| `IsNotPAandNotMBU.xml` | `ASM-FW-GISFW-DATA-CLAUSE` | `ISNOTPAANDNOTMBU` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsNotPrintRISlip.xml` | `ASM-FW-GISFW-WORK` | `ISNOTPRINTRISLIP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISPRINTRISLIP` |
| `IsObjectDeleteVisible.xml` | `ASM-FW-GISFW-WORK` | `ISOBJECTDELETEVISIBLE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsObjectOtherScheduleVisible.xml` | `ASM-FW-GISFW-WORK` | `ISOBJECTOTHERSCHEDULEVISIBLE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsObjectParticipantVisible.xml` | `ASM-FW-GISFW-WORK` | `ISOBJECTPARTICIPANTVISIBLE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsObjectPolicyScheduleVisible.xml` | `ASM-FW-GISFW-WORK` | `ISOBJECTPOLICYSCHEDULEVISIBLE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsObjectSectionAneka.xml` | `ASM-FW-GISFW-DATA` | `ISOBJECTSECTIONANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-ANEKA / ISOBJECTSECTIONANEKA` |
| `IsObjectWithQuantityYear.xml` | `ASM-FW-GISFW-DATA` | `ISOBJECTWITHQUANTITYYEAR` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsOfferFacIn.xml` | `ASM-FW-GISFW-WORK` | `ISOFFERFACIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISFAC` |
| `IsOilGas.xml` | `@BASECLASS` | `ISOILGAS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISOILGAS` |
| `IsOperatorLife.xml` | `DATA-PORTAL` | `ISOPERATORLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / ISAUCTIONUSER` |
| `IsPA.xml` | `DATA-PARTY-PERSON` | `ISPA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / ISTRAVEL` |
| `IsPASSG.xml` | `DATA-PARTY-PERSON` | `ISPASSG` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsPEGAPROD.xml` | `@BASECLASS` | `ISPEGAPROD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `WORK- / ISPEGAPROD` |
| `IsPKSASM.xml` | `ASM-FW-GISFW-WORK` | `ISPKSASM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / IST2T4` |
| `IsPKSEmpty.xml` | `ASM-FW-GISFW-WORK` | `ISPKSEMPTY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-NB / ISPKSEMPTY` |
| `IsPLL.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `ISPLL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-COVERAGE / ISTJH` |
| `IsPolicyPeriodeNotEmpty.xml` | `DATA-PARTY-PERSON` | `ISPOLICYPERIODENOTEMPTY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsPortRisk.xml` | `ASM-FW-GISFW-DATA` | `ISPORTRISK` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISBUILDERRISK` |
| `IsProductsLiability.xml` | `ASM-FW-GISFW-WORK` | `ISPRODUCTSLIABILITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISPRODUCTSLIABILITY` |
| `IsProfessionalLiability.xml` | `ASM-FW-GISFW-WORK` | `ISPROFESSIONALLIABILITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISPROFESSIONALLIABILITY` |
| `IsPropertyandEngineering.xml` | `ASM-FW-GISFW-WORK` | `ISPROPERTYANDENGINEERING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsProposalTransfer.xml` | `ASM-FW-GISFW-DATA-COVERAGE` | `ISPROPOSALTRANSFER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsRISlip.xml` | `ASM-FW-GISFW-WORK` | `ISRISLIP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsRenewal.xml` | `ASM-FW-GISFW-WORK` | `ISRENEWAL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISEDM` |
| `IsRiskIDNull.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `ISRISKIDNULL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsSPVCreate.xml` | `ASM-FW-GISFW-WORK` | `ISSPVCREATE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsSPVTreaty1.xml` | `ASM-FW-GISFW-WORK` | `ISSPVTREATY1` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISSPVCREATE` |
| `IsShowInput.xml` | `@BASECLASS` | `ISSHOWINPUT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsSpecialCase.xml` | `ASM-FW-GISFW-WORK` | `ISSPECIALCASE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsSpreadingUW.xml` | `ASM-FW-GISFW-WORK` | `ISSPREADINGUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsSuccessHitService.xml` | `ASM-FW-GISFW-WORK` | `ISSUCCESSHITSERVICE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsT1T3.xml` | `ASM-FW-GISFW-WORK` | `IST1T3` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsT1T4.xml` | `ASM-FW-GISFW-WORK` | `IST1T4` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsT2T3.xml` | `ASM-FW-GISFW-WORK` | `IST2T3` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / IST1T4` |
| `IsT2T4.xml` | `ASM-FW-GISFW-WORK` | `IST2T4` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / IST1T3` |
| `IsTABActive_Cedant.xml` | `ASM-FW-GISFW-WORK` | `ISTABACTIVE_CEDANT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISTABACTIVE_SPREADING` |
| `IsTABActive_Coverage.xml` | `ASM-FW-GISFW-WORK` | `ISTABACTIVE_COVERAGE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISTABACTIVE_OBJECT` |
| `IsTABActive_Deduction.xml` | `ASM-FW-GISFW-WORK` | `ISTABACTIVE_DEDUCTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISTABACTIVE_CEDANT` |
| `IsTABActive_LossRecord.xml` | `ASM-FW-GISFW-WORK` | `ISTABACTIVE_LOSSRECORD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISTABACTIVE_OBJECT` |
| `IsTABActive_Object.xml` | `ASM-FW-GISFW-WORK` | `ISTABACTIVE_OBJECT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISTABACTIVE` |
| `IsTABActive_Payment.xml` | `ASM-FW-GISFW-WORK` | `ISTABACTIVE_PAYMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISTABACTIVE_DEDUCTION` |
| `IsTABActive_PaymentLife.xml` | `ASM-FW-GISFW-WORK` | `ISTABACTIVE_PAYMENTLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISTABACTIVE_PAYMENT` |
| `IsTABActive_ScoringRisk.xml` | `ASM-FW-GISFW-WORK` | `ISTABACTIVE_SCORINGRISK` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISTABACTIVE_LOSSRECORD` |
| `IsTABActive_Spreading.xml` | `ASM-FW-GISFW-WORK` | `ISTABACTIVE_SPREADING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISTABACTIVE_COVERAGE` |
| `IsTABActive_SpreadingLife.xml` | `ASM-FW-GISFW-WORK` | `ISTABACTIVE_SPREADINGLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISTABACTIVE_SPREADING` |
| `IsTBonding.xml` | `ASM-FW-GISFW-WORK` | `ISTBONDING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / IST1T3` |
| `IsTJH.xml` | `ASM-FW-GISFW-DATA-PREMIUMTEMPLATE` | `ISTJH` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsTravel.xml` | `@BASECLASS` | `ISTRAVEL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISTRAVEL` |
| `IsTreaty1.xml` | `ASM-FW-GISFW-WORK` | `ISTREATY1` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISSPVCREATE` |
| `IsTreatyIn.xml` | `ASM-FW-GISFW-WORK` | `ISTREATYIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsTypeDeductType.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `ISTYPEDEDUCTTYPE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsTypeStock.xml` | `ASM-FW-GISFW-DATA-PROPERTYITEM` | `ISTYPESTOCK` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTYITEM / ISFLAGDELETE` |
| `IsUW.xml` | `ASM-FW-GISFW-WORK` | `ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISUW` |
| `IsUWAcceptance.xml` | `DATA-PARTY-PERSON` | `ISUWACCEPTANCE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsUWAccepted.xml` | `ASM-FW-GISFW-WORK` | `ISUWACCEPTED` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-NB / ISUWACCEPTED` |
| `IsUWSpread.xml` | `ASM-FW-GISFW-WORK` | `ISUWSPREAD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISUWSPREAD` |
| `IsValueDeductType.xml` | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `ISVALUEDEDUCTTYPE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsVisible.xml` | `ASM-FW-GISFW-WORK` | `ISVISIBLE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-NB / ISVISIBLE` |
| `IsWhetherIndex.xml` | `ASM-FW-GISFW-DATA` | `ISWHETHERINDEX` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISPRODUCTSLIABILITY` |
| `IsWorkmenCompensation.xml` | `ASM-FW-GISFW-WORK` | `ISWORKMENCOMPENSATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISCIT` |
| `IsYieldShortfall.xml` | `ASM-FW-GISFW-DATA` | `ISYIELDSHORTFALL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISGROWINGTREES` |
| `LetterNoNull.xml` | `ASM-FW-GISFW-WORK` | `LETTERNONULL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `NopolisEmpty.xml` | `ASM-FW-GISFW-WORK` | `NOPOLISEMPTY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `StepStatusFail.xml` | `@BASECLASS` | `STEPSTATUSFAIL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ToDepHeadUW.xml` | `ASM-FW-GISFW-WORK` | `TODEPHEADUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / TOMANAGERTEKNIK` |
| `ToDeptHeadUWLife.xml` | `ASM-FW-GISFW-WORK` | `TODEPTHEADUWLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / TOMANAGERTEKNIK` |
| `ToDirMarketing.xml` | `ASM-FW-GISFW-WORK` | `TODIRMARKETING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / TODIVHEAD` |
| `ToDirTeknik.xml` | `ASM-FW-GISFW-WORK` | `TODIRTEKNIK` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / TODIRMARKETING` |
| `ToJUW_A.xml` | `ASM-FW-GISFW-WORK` | `TOJUW_A` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / TOSENIORUW` |
| `ToKadivFacultative.xml` | `ASM-FW-GISFW-WORK` | `TOKADIVFACULTATIVE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / TODIVHEAD` |
| `ToKadivFin.xml` | `ASM-FW-GISFW-WORK` | `TOKADIVFIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ToKadivTeknik.xml` | `ASM-FW-GISFW-WORK` | `TOKADIVTEKNIK` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / TODIVHEAD` |
| `ToManagerTeknik.xml` | `ASM-FW-GISFW-WORK` | `TOMANAGERTEKNIK` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ToSeniorUW.xml` | `ASM-FW-GISFW-WORK` | `TOSENIORUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ToTREATYDEPTHEAD.xml` | `ASM-FW-GISFW-WORK` | `TOTREATYDEPTHEAD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / TOMANAGERTEKNIK` |
| `ToUW.xml` | `ASM-FW-GISFW-WORK` | `TOUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / TOJUW_A` |
| `TreatyMasterInEDM.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `TREATYMASTERINEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / TREATYMASTERINEDM` |
| `crmCreateOpportunity.xml` | `PEGACRM-WORK-` | `CRMCREATEOPPORTUNITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `hasPrimaryPage.xml` | `@BASECLASS` | `HASPRIMARYPAGE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isAllRisk.xml` | `ASM-FW-GISFW-DATA` | `ISALLRISK` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isApproved.xml` | `ASM-FW-GISFW-WORK` | `ISAPPROVED` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isAviationHull.xml` | `ASM-FW-GISFW-DATA` | `ISAVIATIONHULL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isBillboardNeonSyariah.xml` | `@BASECLASS` | `ISBILLBOARDNEONSYARIAH` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isBisnisOtherObjAneka.xml` | `ASM-FW-GISFW-WORK` | `ISBISNISOTHEROBJANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isBurglary.xml` | `ASM-FW-GISFW-DATA` | `ISBURGLARY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isCitCis.xml` | `ASM-FW-GISFW-DATA` | `ISCITCIS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isClaimTreaty.xml` | `ASM-FW-GISFW-WORK` | `ISCLAIMTREATY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isCoverageOtherObjAneka.xml` | `ASM-FW-GISFW-WORK` | `ISCOVERAGEOTHEROBJANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isElectronicEquipment.xml` | `@BASECLASS` | `ISELECTRONICEQUIPMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isFidelity.xml` | `@BASECLASS` | `ISFIDELITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK / ISFIDELITY` |
| `isInterestInsuredOtherObjAneka.xml` | `ASM-FW-GISFW-WORK` | `ISINTERESTINSUREDOTHEROBJANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isSellingModeB2B.xml` | `@BASECLASS` | `ISSELLINGMODEB2B` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isSellingModeB2BB2C.xml` | `@BASECLASS` | `ISSELLINGMODEB2BB2C` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isSellingModeB2C.xml` | `@BASECLASS` | `ISSELLINGMODEB2C` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isTravelTime.xml` | `ASM-FW-GISFW-WORK` | `ISTRAVELTIME` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `pyIsIpadOrDesktop.xml` | `@BASECLASS` | `PYISIPADORDESKTOP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `pyIsMobile.xml` | `@BASECLASS` | `PYISMOBILE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

## 4. Tabrakan nama di dalam modul ini

### 4.1 Nama file dipakai di lebih dari satu tipe rule

| Nama file | Tipe | Class |
| --- | --- | --- |
| `InputCommissionMC_FacIn.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-CARGO, ASM-FW-GISFW-DATA-CARGO` |
| `InputDtlTrading_IsUW.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TRADING, ASM-FW-GISFW-DATA-TRADING` |
| `InputDtlCargo_FacIn_IsUW.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-CARGO, ASM-FW-GISFW-DATA-CARGO` |
| `TableInwardScale.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-OFFERFACIN, ASM-FW-GISFW-DATA-OFFERFACIN` |
| `InputDtlObjectAneka_FacIn_IsUW.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-ANEKA, ASM-FW-GISFW-DATA-ANEKA` |
| `ViewCoveragePA.xml` | FlowAction, Section | `DATA-PARTY-PERSON, DATA-PARTY-PERSON` |
| `VehicleGrid_FacIn_IsUW.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-VEHICLE, ASM-FW-GISFW-DATA-VEHICLE` |
| `MappingCoverage.xml` | Activity, DecisionTable | `ASM-FW-GISFW-INT-TABLE_B2B, ASM-FW-GISFW-WORK` |
| `ChooseAccumulation_FacIn.xml` | Harness, Section | `ASM-FW-GISFW-DATA-COVERAGE, ASM-FW-GISFW-DATA-COVERAGE` |
| `InputClauseFire_ViewDtl.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-CLAUSE, ASM-FW-GISFW-DATA-CLAUSE` |
| `RiskAccumulationReport.xml` | Harness, Section | `DATA-PORTAL, DATA-PORTAL` |
| `SelectCoverage.xml` | Harness, Section | `ASM-FW-GISFW-INT-CONDITION, ASM-FW-GISFW-INT-CONDITION` |
| `AddFacOfferList.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-FACOFFER, ASM-FW-GISFW-DATA-FACOFFER` |
| `InputDtlCargo_FacIn.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-CARGO, ASM-FW-GISFW-DATA-CARGO` |
| `PolicyTreatyInDeclineConfirm.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-POLICYTREATYIN, ASM-FW-GISFW-DATA-POLICYTREATYIN` |
| `InputDtlCoverage_FacIn.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-CARGO, ASM-FW-GISFW-DATA-CARGO` |
| `CoverageItem.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-COVERAGE, ASM-FW-GISFW-DATA-COVERAGE` |
| `InputDtlLayerSpreading_FacIn.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-SPREADINGRISK, ASM-FW-GISFW-DATA-SPREADINGRISK` |
| `InputDtlObjectAneka_FacIn.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-ANEKA, ASM-FW-GISFW-DATA-ANEKA` |
| `ViewDtlDeductible.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-COVERAGE, ASM-FW-GISFW-DATA-COVERAGE` |
| `TotalAccumulation_FacIn.xml` | Harness, Section | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-WORK` |
| `ChooseClause.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-CLAUSE, ASM-FW-GISFW-DATA-CLAUSE` |
| `InputClause_ViewDtl.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-CLAUSE, ASM-FW-GISFW-DATA-CLAUSE` |
| `ViewOldDeductible.xml` | Harness, Section | `ASM-FW-GISFW-DATA-COVERAGE, ASM-FW-GISFW-DATA-COVERAGE` |
| `IsUWAccepted.xml` | DecisionTable, When | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-WORK` |
| `InputCommissionPA_FacIn.xml` | FlowAction, Section | `DATA-PARTY-PERSON, DATA-PARTY-PERSON` |
| `SummaryRiskAccumulation.xml` | Harness, Section | `DATA-PORTAL, DATA-PORTAL` |
| `ChooseOccupation.xml` | FlowAction, Harness, Section | `ASM-FW-GISFW-DATA-OCCUPATION, ASM-FW-GISFW-DATA-OCCUPATION, ASM-FW-GISFW-DATA-OCCUPATION` |
| `ViewDtlOtherObjectAneka.xml` | FlowAction, Section | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-WORK` |
| `ViewCoverageAneka.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-ANEKA, ASM-FW-GISFW-DATA-COVERAGE` |
| `SpreadingItem.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-COVERAGE, ASM-FW-GISFW-DATA-COVERAGE` |
| `InputFEA_IsUW.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-OFFERFACIN-OFFERFEALIST, ASM-FW-GISFW-DATA-OFFERFACIN-OFFERFEALIST` |
| `VehicleGrid_FacIn.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-VEHICLE, ASM-FW-GISFW-DATA-VEHICLE` |
| `ProtectCurrencyCargo_IsUW.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-DEDUCTIBLE, ASM-FW-GISFW-DATA-DEDUCTIBLE` |
| `isApproved.xml` | DecisionTable, When | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-WORK` |
| `ViewObjectDtlAneka.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-OCCUPATION, ASM-FW-GISFW-DATA-OCCUPATION` |
| `ChooseRiskAddress.xml` | FlowAction, Harness, Section | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE, ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE, ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` |
| `InputDtlTrading.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TRADING, ASM-FW-GISFW-DATA-TRADING` |
| `TableOfLimit.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-OFFERFACIN, ASM-FW-GISFW-DATA-OFFERFACIN` |
| `SpreadingItemAneka.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-COVERAGE, ASM-FW-GISFW-DATA-COVERAGE` |
| `InputSpreadingLife_FacIn.xml` | FlowAction, Section | `DATA-PARTY-PERSON, DATA-PARTY-PERSON` |
| `ChooseObject_Ship.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-ANEKA, ASM-FW-GISFW-DATA-ANEKA` |
| `InputCauseOfLoss_FacIn_IsUW.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-CAUSEOFLOSS, ASM-FW-GISFW-DATA-CAUSEOFLOSS` |
| `ViewClaimList.xml` | Harness, Section | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-WORK` |
| `ViewDtlObjectAneka.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-ANEKA, ASM-FW-GISFW-DATA-ANEKA` |
| `PersonPlanAndDeduct.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-COVERAGE, ASM-FW-GISFW-DATA-COVERAGE` |
| `InputDtlOtherObjectAneka_FacIn_ISUW.xml` | FlowAction, Section | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-WORK` |
| `ShowCedingCoList.xml` | Harness, Section | `ASM-FW-GISFW-DATA-OFFERFACIN, ASM-FW-GISFW-DATA-OFFERFACIN` |
| `CommisionItem.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-COVERAGE, ASM-FW-GISFW-DATA-COVERAGE` |
| `SelectShip.xml` | Harness, Section | `ASM-FW-GISFW-INT-SHIP, ASM-FW-GISFW-INT-SHIP` |
| `PrintRISlip.xml` | Harness, Section | `ASM-FW-GISFW-DATA-FACOFFER, ASM-FW-GISFW-WORK` |
| `LimitTreaty.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-OFFERFACIN, ASM-FW-GISFW-DATA-OFFERFACIN` |
| `InputCity.xml` | FlowAction, Section | `@BASECLASS, @BASECLASS` |
| `InputFEA.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-OFFERFACIN-OFFERFEALIST, ASM-FW-GISFW-DATA-FEA` |
| `InputDtlConveyance_IsUW.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-CONVEYANCE, ASM-FW-GISFW-DATA-CONVEYANCE` |
| `ChooseClassofContraction.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-OCCUPATION, ASM-FW-GISFW-DATA-OCCUPATION` |
| `GetTreatyName.xml` | Activity, RDBList | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-INT-PROPORTIONALARRG` |
| `InputDtlGoods.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-GOOD, ASM-FW-GISFW-DATA-GOOD` |
| `ReasViewAttachment.xml` | FlowAction, Section | `ASM-FW-GISFW-INT, ASM-FW-GISFW-INT` |
| `ProtectCurrencyCargo.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-DEDUCTIBLE, ASM-FW-GISFW-DATA-DEDUCTIBLE` |
| `InputCommissionLife_FacIn.xml` | FlowAction, Section | `DATA-PARTY-PERSON, DATA-PARTY-PERSON` |
| `CedingCedant_IsUW.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-QUOTATION, ASM-FW-GISFW-DATA-QUOTATION` |
| `ViewPolicyListRO.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-POLICY, ASM-FW-GISFW-DATA-POLICY` |
| `InputSpreadingTravel_FacIn_IsUW.xml` | FlowAction, Section | `DATA-PARTY-PERSON, DATA-PARTY-PERSON` |
| `InputNation.xml` | FlowAction, Section | `@BASECLASS, @BASECLASS` |
| `InputProvince.xml` | FlowAction, Section | `@BASECLASS, @BASECLASS` |
| `InputRW.xml` | FlowAction, Section | `@BASECLASS, @BASECLASS` |
| `ViewCoverageTravel.xml` | FlowAction, Section | `DATA-PARTY-PERSON, DATA-PARTY-PERSON` |
| `PlanListRO.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-COVERAGE, ASM-FW-GISFW-DATA-COVERAGE` |
| `PersonCoverage.xml` | FlowAction, Section | `DATA-PARTY-PERSON, DATA-PARTY-PERSON` |
| `CorrespondenceContent.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-CORRESPONDENCE, ASM-FW-GISFW-DATA-CORRESPONDENCE` |
| `ChooseSubContract_FacIn.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-SUBCONTRACT, ASM-FW-GISFW-DATA-SUBCONTRACT` |
| `ChooseCoverage.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-COVERAGE, ASM-FW-GISFW-DATA-COVERAGE` |
| `InputDistrict.xml` | FlowAction, Section | `@BASECLASS, @BASECLASS` |
| `PlanContains_UWRO.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-PLAN, ASM-FW-GISFW-DATA-PLAN` |
| `InsertLogServiceProd.xml` | Activity, RDBList | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-WORK` |
| `InputDtlGoods_IsUW.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-GOOD, ASM-FW-GISFW-DATA-GOOD` |
| `SpreadingItem_IsUW.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-COVERAGE, ASM-FW-GISFW-DATA-COVERAGE` |
| `ChooseSubContract.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-SUBCONTRACT, ASM-FW-GISFW-DATA-SUBCONTRACT` |
| `ShowCoverageFacOut.xml` | Activity, FlowAction, Section | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-DATA-PROPERTY, ASM-FW-GISFW-DATA-PROPERTYITEM` |
| `confirmCreateAccumulation.xml` | FlowAction, Section | `DATA-PARTY-PERSON, DATA-PARTY-PERSON` |
| `AccumulationRisk.xml` | Harness, Section | `DATA-PORTAL, DATA-PORTAL` |
| `InputCommissionMBU_FacIn.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-VEHICLE, ASM-FW-GISFW-DATA-VEHICLE` |
| `AddFacOfferList_IsUW.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-FACOFFER, ASM-FW-GISFW-DATA-FACOFFER` |
| `InputDtlConveyance.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-CONVEYANCE, ASM-FW-GISFW-DATA-CONVEYANCE` |
| `ProtectCurrency_IsUW.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-DEDUCTIBLE, ASM-FW-GISFW-DATA-DEDUCTIBLE` |
| `ShowPolis.xml` | FlowAction, Harness | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-WORK` |
| `BenefitListRO.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-PLAN, ASM-FW-GISFW-DATA-PLAN` |
| `AdjustmentRiskAccumulation.xml` | Harness, Section | `DATA-PORTAL, DATA-PORTAL` |
| `BrowseM_OKUPASI_ANEKA.xml` | Activity, ReportDefinition | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-INT-M_OKUPASI_ANEKA` |
| `ViewDtlAddCargo.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-COVERAGE, ASM-FW-GISFW-DATA-COVERAGE` |
| `CoverageItem_IsUW.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-COVERAGE, ASM-FW-GISFW-DATA-COVERAGE` |
| `ChooseClauseFire.xml` | Harness, Section | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-WORK` |
| `InputCauseOfLoss_FacIn.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-CAUSEOFLOSS, ASM-FW-GISFW-DATA-CAUSEOFLOSS` |
| `ViewMasterCargo.xml` | Activity, RDBList | `ASM-FW-GISFW-DATA-OFFERFACIN, ASM-FW-GISFW-DATA-OFFERFACIN` |
| `TreatyRealizationCheckDuplicate.xml` | Activity, RDBList | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-WORK` |
| `CommisionItemAneka.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-COVERAGE, ASM-FW-GISFW-DATA-COVERAGE` |
| `InputCauseOfDecline.xml` | Harness, Section | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-WORK` |
| `CopyCoverageFrom.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-PROPERTYITEM, ASM-FW-GISFW-DATA-PROPERTYITEM` |
| `InputDtlSpreadingPA_FacIn_IsUW.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-COVERAGE, ASM-FW-GISFW-DATA-COVERAGE` |
| `LoadPackingItem.xml` | Activity, DataTransform | `ASM-FW-GISFW-DATA-CARGO, ASM-FW-GISFW-DATA-CARGO` |
| `InputSpreadingPA_FacIn.xml` | FlowAction, Section | `DATA-PARTY-PERSON, DATA-PARTY-PERSON` |
| `SelectAgent.xml` | Harness, Section | `ASM-FW-GISFW-INT-LLOYDAGENT, ASM-FW-GISFW-INT-LLOYDAGENT` |
| `InputSpreadingPA_FacIn_IsUW.xml` | FlowAction, Section | `DATA-PARTY-PERSON, DATA-PARTY-PERSON` |
| `GetCauseOfDecline.xml` | Activity, RDBList | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-DATA-CAUSEOFDECLINE` |
| `InputDtlOtherObjectAneka.xml` | FlowAction, Section | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-WORK` |
| `InputDtlSpreadingPA_FacIn.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-COVERAGE, ASM-FW-GISFW-DATA-COVERAGE` |
| `OfferFacOut.xml` | Flow, FlowAction | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-WORK` |
| `SpreadingItemAneka_IsUW.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-COVERAGE, ASM-FW-GISFW-DATA-COVERAGE` |
| `ProtectCurrency.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-DEDUCTIBLE, ASM-FW-GISFW-DATA-DEDUCTIBLE` |
| `InputDtlSpreadingTravel_FacIn_IsUW.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-COVERAGE, ASM-FW-GISFW-DATA-COVERAGE` |
| `InputDtlOtherObjectAneka_FacIn.xml` | FlowAction, Section | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-WORK` |
| `Installments_ReadOnly.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYININSTALLMENT, ASM-FW-GISFW-DATA-TREATYININSTALLMENT` |
| `TotalAccumulationDtl.xml` | Harness, Section | `ASM-FW-GISFW-DATA-COVERAGE, ASM-FW-GISFW-DATA-COVERAGE` |
| `InstallmentList.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-INSTALLMENT, ASM-FW-GISFW-DATA-INSTALLMENT` |
| `ChooseDeductible.xml` | Harness, Section | `ASM-FW-GISFW-DATA-DEDUCTIBLE, ASM-FW-GISFW-DATA-DEDUCTIBLE` |
| `InputSpreadingLife_FacIn_IsUW.xml` | FlowAction, Section | `DATA-PARTY-PERSON, DATA-PARTY-PERSON` |
| `InputSpreadingTravel_FacIn.xml` | FlowAction, Section | `DATA-PARTY-PERSON, DATA-PARTY-PERSON` |
| `GetInsuredID.xml` | Activity, RDBList | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-INT-CLIENT` |
| `InputAnekaParticipant.xml` | FlowAction, Section | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-WORK` |
| `CedingCedant.xml` | FlowAction, Harness, Section | `ASM-FW-GISFW-DATA-QUOTATION, ASM-FW-GISFW-DATA-QUOTATION, ASM-FW-GISFW-DATA-QUOTATION` |
| `InputCommentDeductible.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-OFFERFACIN, ASM-FW-GISFW-DATA-OFFERFACIN` |

Total: **122** nama file.

### 4.2 Rule `When` bernama sama di lebih dari satu class

_Tidak ada._

## 5. Keluarga klon (satu asal salinan dipakai >1 rule)

Dibaca dari `<pzOriginalInstanceKey>` yang sama pada beberapa file. Ini **sinyal kandidat
konsolidasi** saat migrasi, bukan bukti bahwa perilakunya identik — kesamaan perilaku baru
dapat dinyatakan setelah isinya dibandingkan di STEP D2.

| Asal salinan (class / nama) | Jumlah | File |
| --- | ---: | --- |
| `ASM-FW-GISFW-DATA-VEHICLE / INPUTSPREADINGMBU_FACIN` | 2 | `FlowAction/InputSpreadingMBU_FacISUW.xml`, `FlowAction/InputSpreadingMBU_FacIn.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / PERSONPLANANDDEDUCT` | 2 | `FlowAction/PersonPlanAndDeduct.xml`, `Section/PersonPlanAndDeduct.xml` |
| `ASM-FW-GISFW-DATA-CARGO / INPUTDTLCARGO_FACIN` | 3 | `FlowAction/InputDtlCargo_FacIn_IsUW.xml`, `Section/InputDtlCargo_FacIn.xml`, `Section/InputDtlCargo_FacIn_IsUW.xml` |
| `DATA-PARTY-PERSON / INPUTDTLPARTICIPANTLIFE_FACIN` | 2 | `Section/InputDtlParticipantLife_FacIn.xml`, `Section/InputDtlParticipantLife_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETJSONDATAFOLLOWING_SQL` | 3 | `RDBList/GetIDPegaJsonFollowing_SQL.xml`, `RDBList/GetJsonDataFollowingView_SQL.xml`, `RDBList/GetJsonDataFollowing_SQL.xml` |
| `ASM-FW-GISFW-WORK / ADDTOLISTCOMMENTSPOLICYTREATYIN_DT` | 2 | `DataTransform/AddToListCommentsPolicyTreatyIn_DT.xml`, `DataTransform/DeptHeadTreatyIn_UW_postDT.xml` |
| `ASM-FW-GISFW-DATA-CLAUSE / INPUTCLAUSEFIRE_FACIN` | 2 | `Section/InputClauseFire_FacIn.xml`, `Section/InputClauseFire_FacInUW.xml` |
| `ASM-FW-GISFW-WORK / ISTABACTIVE_OBJECT` | 2 | `When/IsTABActive_Coverage.xml`, `When/IsTABActive_LossRecord.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / SPREADINGITEM` | 6 | `FlowAction/CommisionItem.xml`, `FlowAction/SpreadingItem_IsUW.xml`, `Section/CommisionItem.xml`, `Section/SpreadingItem.xml`, `Section/SpreadingItemAneka.xml`, `Section/SpreadingItem_IsUW.xml` |
| `ASM-FW-GISFW-DATA-QUOTATION / BTNCEDINGCO_DT` | 2 | `DataTransform/btnCedingCO_DT.xml`, `DataTransform/btnCedingCedant_DT.xml` |
| `ASM-FW-GISFW-WORK / INPUTDTLOTHEROBJECTANEKA` | 5 | `FlowAction/InputDtlOtherObjectAneka.xml`, `FlowAction/InputDtlOtherObjectAneka_FacIn.xml`, `FlowAction/ViewDtlOtherObjectAneka.xml`, `Section/InputDtlOtherObjectAneka.xml`, `Section/ViewDtlOtherObjectAneka.xml` |
| `ASM-FW-GISFW-WORK / ISFAC` | 2 | `When/IsFac.xml`, `When/IsOfferFacIn.xml` |
| `ASM-FW-GISFW-DATA-OCCUPATION / ANEKALISTFACOUT` | 2 | `Section/AnekaListFacOut.xml`, `Section/AnekaListFacOut_IsUW.xml` |
| `DATA-PARTY-PERSON / PERSONCOVERAGE` | 2 | `FlowAction/PersonCoverage.xml`, `Section/PersonCoverage.xml` |
| `ASM-FW-GISFW-WORK / ASM!GETBRANCHCODE` | 2 | `RDBList/GetCountAccumulation.xml`, `RDBList/GetNPWPbySOB_SQL.xml` |
| `ASM-FW-GISFW-DATA-QUOTATION / INPUTHISTORICALSURVEYREPORT` | 2 | `FlowAction/InputHistoricalSurveyReport.xml`, `FlowAction/InputHistoricalSurveyReportUW.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / LAYERLISTDTL` | 2 | `Section/LayerListDtl.xml`, `Section/LayerListDtl_IsUW.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLLAYERLISTFIRE_FACIN` | 2 | `FlowAction/InputDtlLayerListFire_FacIn.xml`, `FlowAction/InputDtlLayerListFire_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-DATA-DEDUCTIBLE / ADDDEDUCTIBLE` | 2 | `Section/addDeductible.xml`, `Section/addDeductible_IsUW.xml` |
| `ASM-FW-GISFW-WORK / INPUTANEKAPARTICIPANT` | 2 | `FlowAction/InputAnekaParticipant.xml`, `Section/InputAnekaParticipant.xml` |
| `ASM-FW-GISFW-WORK / INPUTOTHEROBJECTANEKA_FACIN` | 2 | `Section/InputOtherObjectAneka_FacIn.xml`, `Section/InputOtherObjectAneka_FacIn_ISUW.xml` |
| `DATA-PARTY-PERSON / SHOWCOVERAGEPAFACOUT_FLOWACTION` | 3 | `FlowAction/ShowCoveragePAFacOutIsUW_FlowAction.xml`, `FlowAction/ShowCoveragePAFacOut_FlowAction.xml`, `FlowAction/ShowCoverageTravelFacOut_FlowAction.xml` |
| `ASM-FW-GISFW-DATA-PLAN / PLANCONTAINS_UWRO` | 2 | `FlowAction/PlanContains_UWRO.xml`, `Section/PlanContains_UWRO.xml` |
| `ASM-FW-GISFW-INT-AGENT / BROWSEAGENTNUSARE_RD` | 4 | `ReportDefinition/BrowseAgentNonLife_RD.xml`, `ReportDefinition/BrowseAgentNusaReLife_RD.xml`, `ReportDefinition/BrowseAgentNusaRe_RD.xml`, `ReportDefinition/BrowseAgent_RD.xml` |
| `ASM-FW-GISFW-WORK / ASM!GETNPWPBYSOB_SQL` | 2 | `RDBList/GetGroupName_SQL.xml`, `RDBList/GetNPWPbyCedingCo_SQL.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / INPUTCOVERAGESPREADING` | 3 | `FlowAction/InputCoverageCommision.xml`, `FlowAction/InputCoverageSpreading.xml`, `FlowAction/InputCoverageSpreading_IsUW.xml` |
| `ASM-FW-GISFW-INT-RW / BROWSEPROVINCE_RD` | 2 | `ReportDefinition/BrowseCity_RD.xml`, `ReportDefinition/BrowseProvince_RD.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLADDCARGO` | 2 | `FlowAction/ViewDtlAddCargo.xml`, `Section/ViewDtlAddCargo.xml` |
| `ASM-FW-GISFW-WORK / CHECKRISLIP` | 2 | `Activity/CheckRISLIP.xml`, `Activity/CheckRISLIP_EDM.xml` |
| `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETTANGGALCLOSING_SQL` | 2 | `RDBList/GETTanggalClosing_SQL.xml`, `RDBList/GetSequenceNumber_SQL.xml` |
| `ASM-FW-GISFW-DATA-PROPERTYITEM / INPUTCOVERAGESPREADINGFIRE` | 4 | `Section/InputCoverageCommisionFire.xml`, `Section/InputCoverageSpreadingAneka.xml`, `Section/InputCoverageSpreadingFire.xml`, `Section/InputCoverageSpreadingFire_IsUW.xml` |
| `ASM-FW-GISFW-DATA-PROPERTY / INPUTOBJFIREGRID` | 3 | `FlowAction/InputObjFireGrid.xml`, `FlowAction/InputObjFireGrid_GCNM.xml`, `FlowAction/InputObjFireGrid_IsUW.xml` |
| `@BASECLASS / INPUTCITY` | 2 | `FlowAction/InputDistrict.xml`, `Section/InputCity.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY / INPUTDTLPAYMENTFACRETRO_FLOWACT` | 2 | `FlowAction/InputDtlPaymentFacRetro_FlowAct.xml`, `FlowAction/InputDtlPaymentFacRetro_FlowAct_IsUW.xml` |
| `ASM-FW-GISFW-WORK / UPLOADCSV_VEHICLE` | 2 | `FlowAction/UploadCSV_Participant.xml`, `FlowAction/UploadCSV_Vehicle.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / FEALIST` | 2 | `Section/FEAList.xml`, `Section/FEAList_IsUW.xml` |
| `DATA-PARTY-PERSON / SHOWCOVERAGEPAFACOUT_DTL` | 3 | `Section/ShowCoveragePAFacOutIsUW_Dtl.xml`, `Section/ShowCoveragePAFacOut_Dtl.xml`, `Section/ShowCoverageTravelFacOut_Dtl.xml` |
| `DATA-PARTY-PERSON / INPUTSPREADINGLIFE_FACIN` | 3 | `FlowAction/InputSpreadingLife_FacIn_IsUW.xml`, `Section/InputSpreadingLife_FacIn.xml`, `Section/InputSpreadingLife_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / INPUTDTLOBJECTLOCATION_FACIN` | 2 | `Section/InputDtlObjectLocation_FacIn.xml`, `Section/InputDtlObjectLocation_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN / PERIODE` | 2 | `Section/Periode.xml`, `Section/Periode_IsUW.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / PROPERTYITEMLISTCOVERAGESPREADING` | 2 | `Section/PropertyItemListCoverageCommision.xml`, `Section/PropertyItemListCoverageSpreading_IsUW.xml` |
| `ASM-FW-GISFW-DATA-PROPERTYITEM / PROPERTYITEMFACIN_SECTION` | 2 | `Section/PropertyItemFacIn_Section.xml`, `Section/PropertyItemFacIn_Section_IsUW.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY / INPUTDTLPAYMENT_FACINLIFE` | 2 | `Section/InputDtlPayment_FacInLife.xml`, `Section/InputDtlPayment_FacInLife_IsUW.xml` |
| `ASM-FW-GISFW-WORK / INPUTPOLICYTREATYIN_PREDT` | 2 | `DataTransform/DeptHeadTreatyInUW_preDT.xml`, `DataTransform/InputPolicyTreatyIn_preDT.xml` |
| `ASM-FW-GISFW-DATA-ANEKA / INPUTCOVERAGEDEDUCTIBLEANEKA_GCNM` | 2 | `Section/InputCoverageDeductibleAneka_GCNM.xml`, `Section/InputCoverageDeductibleHull_GCNM.xml` |
| `ASM-FW-GISFW-WORK / PAYMENTCURRENCYLIST` | 2 | `Section/PaymentCurrencyList.xml`, `Section/PaymentCurrencyList_IsUW.xml` |
| `ASM-FW-GISFW-INT-ACCUMULATION / ASM!GETSUMMARYRISKACCUMULATION_SQL` | 2 | `RDBList/GetSummaryRiskAccumPolis_Sql.xml`, `RDBList/GetSummaryRiskAccumulation_Sql.xml` |
| `ASM-FW-GISFW-DATA-DEDUCTIBLE / CHOOSEDEDUCTIBLE` | 2 | `Harness/ChooseDeductible.xml`, `Section/ChooseDeductible.xml` |
| `ASM-FW-GISFW-DATA-SCORINGRISK / SCORINGRISKFORM1` | 3 | `Section/ScoringRiskForm1.xml`, `Section/ScoringRiskForm2.xml`, `Section/ScoringRiskForm4.xml` |
| `ASM-FW-GISFW-INT-LLOYDAGENT / SELECTAGENT` | 2 | `Harness/SelectAgent.xml`, `Section/SelectAgent.xml` |
| `ASM-FW-GISFW-INT-V_DUPLICATE_VEHICLE / CHECKMST_KENDARAANBYPLATCHASENGINE` | 3 | `ReportDefinition/CheckMST_KENDARAANByChas.xml`, `ReportDefinition/CheckMST_KENDARAANByEngine.xml`, `ReportDefinition/CheckMST_KENDARAANByPlat.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLSPREADINGTRAVEL_FACIN` | 2 | `FlowAction/InputDtlSpreadingTravel_FacIn_IsUW.xml`, `Section/InputDtlSpreadingTravel_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / PROPERTYITEMLISTCOVERAGE` | 2 | `Section/PropertyItemListCoverageSpreading.xml`, `Section/PropertyItemListFacOut.xml` |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN / POLICYTREATYINDECLINECONFIRM` | 2 | `FlowAction/PolicyTreatyInDeclineConfirm.xml`, `Section/PolicyTreatyInDeclineConfirm.xml` |
| `ASM-FW-GISFW-INT-AGENT / BROWSEAGENTHIERARKILIST_RD` | 2 | `ReportDefinition/BrowseAgentHierarkiList_RD.xml`, `ReportDefinition/BrowseCedingCo_RD.xml` |
| `ASM-FW-GISFW-WORK / GENERATEPDFFACOUTMEMOPLACING` | 3 | `Activity/DeletePDFFacOutMemoPlacing.xml`, `Activity/GeneratePDFFacOutMemoPlacing.xml`, `Activity/GeneratePDFSlipKomisi.xml` |
| `@BASECLASS / HTMLTOPDF` | 2 | `Activity/HTMLRISlipToPDF.xml`, `Activity/HTMLToPDF.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / OBJECTGRIDLOCATION_FACIN` | 2 | `FlowAction/ObjectGridLocation_FacIn.xml`, `FlowAction/ObjectGridLocation_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-WORK / ISEDM` | 2 | `When/IsEDM.xml`, `When/IsRenewal.xml` |
| `ASM-FW-GISFW-WORK / SETREJECTPROPOSAL` | 2 | `DataTransform/SetRejectProposal.xml`, `DataTransform/SetReviseProposal.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / PROPERTYITEMLIST` | 3 | `Section/PropertyItemList.xml`, `Section/PropertyItemListCoverage.xml`, `Section/PropertyItemList_IsUW.xml` |
| `ASM-FW-GISFW-WORK / FIRESUMMARYSECTION` | 3 | `Section/AllSummarySection.xml`, `Section/FireSummarySection.xml`, `Section/GolfSummarySection.xml` |
| `ASM-FW-GISFW-INT-SHIP / BROWSESHIP_RD` | 2 | `ReportDefinition/BrowseShip2_RD.xml`, `ReportDefinition/BrowseShip_RD.xml` |
| `ASM-FW-GISFW-DATA-OCCUPATION / SPREADINGOBJECTANEKA` | 3 | `Section/CommisionObjectAneka.xml`, `Section/SpreadingObjectAneka.xml`, `Section/SpreadingObjectAneka_IsUW.xml` |
| `ASM-FW-GISFW-DATA-ANEKA / INPUTDTLOBJECTANEKA_FACIN` | 4 | `FlowAction/InputDtlObjectAneka_FacIn.xml`, `FlowAction/InputDtlObjectAneka_FacIn_IsUW.xml`, `Section/InputDtlObjectAneka_FacIn.xml`, `Section/InputDtlObjectAneka_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-DATA-CARGO / LOADPACKINGITEM` | 2 | `Activity/LoadPackingItem.xml`, `DataTransform/LoadPackingItem.xml` |
| `ASM-FW-GISFW-DATA-VEHICLE / INPUTCOMMISSIONMBU_FACIN` | 2 | `FlowAction/InputCommissionMBU_FacIn.xml`, `Section/InputCommissionMBU_FacIn.xml` |
| `ASM-FW-GISFW-DATA-OCCUPATION / CHOOSEOCCUPATION` | 3 | `FlowAction/ChooseOccupation.xml`, `Harness/ChooseOccupation.xml`, `Section/ChooseOccupation.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY / INPUTDTLPAYMENTFLOW_FACIN` | 3 | `FlowAction/InputDtlPaymentFlow_FacIn.xml`, `FlowAction/InputDtlPaymentFlow_FacInLife.xml`, `FlowAction/InputDtlPaymentFlow_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-DATA-CLAUSE / INPUTCLAUSE_VIEWDTL` | 2 | `FlowAction/InputClause_ViewDtl.xml`, `Section/InputClause_ViewDtl.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / SUMTSIPREMISPREADEDRNM_ACT` | 7 | `Activity/SumTSIPremiSpreadedRNM_ANEKA_Act.xml`, `Activity/SumTSIPremiSpreadedRNM_Act.xml`, `Activity/SumTSIPremiSpreadedRNM_FIRE_Act.xml`, `Activity/SumTSIPremiSpreadedRNM_GOLF_Act.xml`, `Activity/SumTSIPremiSpreadedRNM_MBU_Act.xml`, `Activity/SumTSIPremiSpreadedRNM_PA_Act.xml`, `Activity/SumTSIPremiSpreadedRNM_TRAVEL_Act.xml` |
| `ASM-FW-GISFW-WORK-INT-V_ZIPCODE / ASM!SEARCHKOTAIDSQL` | 2 | `RDBList/SearchCoverageIDSQL.xml`, `RDBList/SearchRWIDSQL.xml` |
| `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETLIMITAKSEPTASINONFIRE_SQL` | 3 | `RDBList/GetLimitAkseptasiNonFireBanding_SQL.xml`, `RDBList/GetLimitAkseptasiNonFireJUWA_SQL.xml`, `RDBList/GetLimitAkseptasiNonFire_SQL.xml` |
| `ASM-FW-GISFW-WORK / INPUTDTLOTHEROBJECTANEKA_FACIN` | 3 | `FlowAction/InputDtlOtherObjectAneka_FacIn_ISUW.xml`, `Section/InputDtlOtherObjectAneka_FacIn.xml`, `Section/InputDtlOtherObjectAneka_FacIn_ISUW.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / INSERTGOOGLESTORAGE_ACT` | 2 | `Activity/GetUrlGoogleStorage_Act.xml`, `Activity/InsertGoogleStorage_Act.xml` |
| `DATA-PORTAL / SUMMARYRISKACCUMULATION` | 2 | `Harness/SummaryRiskAccumulation.xml`, `Section/SummaryRiskAccumulation.xml` |
| `ASM-FW-GISFW-WORK / COPYTOALLSPREADING_ACT` | 4 | `Activity/CopyToAllSpreadingAnekaGolf_ACT.xml`, `Activity/CopyToAllSpreadingCargoMBU_ACT.xml`, `Activity/CopyToAllSpreadingFire_ACT.xml`, `Activity/CopyToAllSpreading_ACT.xml` |
| `ASM-FW-GISFW-DATA-POLICY / VIEWPOLICYLISTRO` | 2 | `FlowAction/ViewPolicyListRO.xml`, `Section/ViewPolicyListRO.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / CEKSPREADINGFACTIN` | 2 | `Activity/cekSpreadingFactIn.xml`, `Activity/cekSpreadingFactIn_Act.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLSPREADINGCOVERAGE_FACIN` | 3 | `Section/InputDtlSpreadingCommisionCoverage_FacIn.xml`, `Section/InputDtlSpreadingCoverage_FacIn.xml`, `Section/InputDtlSpreadingCoverage_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-DATA-PROPERTYITEM / PROPERTYITEMCOVERAGESPREADINGFACIN_FLOWACTION` | 2 | `FlowAction/PropertyItemCoverageCommisionFacIn_FlowAction.xml`, `FlowAction/PropertyItemCoverageSpreadingFacIn_FlowAction_IsUW.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / FILLPREMIBOND` | 2 | `Activity/fillPremiAneka.xml`, `DataTransform/fillPremiBond.xml` |
| `ASM-FW-GISFW-WORK / COUNTGROSSPREMI_ACT` | 3 | `Activity/CountGPWMarinePAMbu_Act.xml`, `Activity/CountGrossPremiEDM_Act.xml`, `Activity/CountGrossPremi_Act.xml` |
| `ASM-FW-GISFW-DATA-FACOFFER / SETENDORSEMENTRISLIPDATA` | 2 | `Activity/SetEndorsementRISlipData.xml`, `Activity/ViewEndorsementRISlipData.xml` |
| `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETVIEWBANDING_SQL` | 2 | `RDBList/GetAksepBanding_SQL.xml`, `RDBList/GetFlagReject_SQL.xml` |
| `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETLIMITAKSEPTASIBOND_SQL` | 3 | `RDBList/GetLimitAkseptasiBond_SQL.xml`, `RDBList/GetLimitAkseptasiKreditCL_SQL.xml`, `RDBList/GetLimitAkseptasiKreditNCL_SQL.xml` |
| `ASM-FW-GISFW-WORK / VIEWDATAPOLICY` | 2 | `Section/ViewDataPolicy.xml`, `Section/ViewDataPolicy_IsUW.xml` |
| `ASM-FW-GISFW-DATA-MEMORANDUMACCEPT / INPUTDTLMEMORANDUMACCEPT` | 3 | `Section/InputDtlMemorandumAccept.xml`, `Section/InputDtlMemorandumAccept_IsUW.xml`, `Section/ViewDtlMemorandumAccept.xml` |
| `ASM-FW-GISFW-INT-RISKADDRESS / BROWSERISKADDRESS_RD` | 3 | `ReportDefinition/BrowseRiskAddressZipCode_RD.xml`, `ReportDefinition/BrowseRiskAddress_RD.xml`, `ReportDefinition/BrowseRisksAddress_RD.xml` |
| `ASM-FW-GISFW-DATA-CLAUSE / CHOOSECLAUSE` | 2 | `FlowAction/ChooseClause.xml`, `Section/ChooseClause.xml` |
| `DATA-PARTY-PERSON / INPUTSPREADINGPA_FACIN` | 6 | `FlowAction/InputSpreadingPA_FacIn.xml`, `FlowAction/InputSpreadingPA_FacIn_IsUW.xml`, `FlowAction/InputSpreadingTravel_FacIn.xml`, `Section/InputSpreadingPA_FacIn.xml`, `Section/InputSpreadingPA_FacIn_IsUW.xml`, `Section/InputSpreadingTravel_FacIn.xml` |
| `ASM-FW-GISFW-DATA-QUOTATION / INPUTHISTORICALSURVEYREPORTDTL` | 2 | `Section/InputHistoricalSurveyReportDtl.xml`, `Section/InputHistoricalSurveyReportDtlUW.xml` |
| `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETLIMITACCENGINEERINGUW_SQL` | 2 | `RDBList/GetLimitAccEngineeringBanding_SQL.xml`, `RDBList/GetLimitAccEngineeringJUWA_SQL.xml` |
| `ASM-FW-GISFW-INT-RW / ASM!BROWSERW_SQL` | 2 | `RDBList/BrowseRW2_SQL.xml`, `RDBList/BrowseRW_SQL.xml` |
| `ASM-FW-GISFW-DATA-OUTGO / INPUTOUTGO` | 2 | `Section/InputOutgo.xml`, `Section/ViewOutgo.xml` |
| `DATA-PARTY-PERSON / CONFIRMCREATEACCUMULATION` | 2 | `FlowAction/confirmCreateAccumulation.xml`, `Section/confirmCreateAccumulation.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / COVERAGEITEM` | 4 | `FlowAction/CoverageItem.xml`, `FlowAction/CoverageItem_IsUW.xml`, `FlowAction/SpreadingItem.xml`, `Section/CoverageItem.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / PROPERTYITEMFACOUT_FLOWACTION` | 2 | `FlowAction/PropertyItemFacOut_FlowAction.xml`, `FlowAction/PropertyItemFacOut_FlowAction_IsUW.xml` |
| `ASM-FW-GISFW-INT-CONDITION / SELECTCOVERAGE` | 2 | `Harness/SelectCoverage.xml`, `Section/SelectCoverage.xml` |
| `ASM-FW-GISFW-INT-OFFERJSON / ASM!SEARCHRATEPOLIS_SQL` | 2 | `RDBList/RateFlexas_Sql.xml`, `RDBList/SearchRatePolis_SQL.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / COVERAGESUMMARY` | 2 | `Section/CoverageSummary.xml`, `Section/GolfCoverageSummary.xml` |
| `@BASECLASS / INPUTPROVINCE` | 2 | `FlowAction/InputProvince.xml`, `Section/InputProvince.xml` |
| `ASM-FW-GISFW-DATA-OCCUPATION / ANEKALISTFACOUT_FLOWACTION` | 2 | `FlowAction/AnekaListFacOut_FlowAction.xml`, `FlowAction/AnekaListFacOut_FlowAction_IsUW.xml` |
| `ASM-FW-GISFW-DATA-TREATYININSTALLMENT / INSTALLMENTS_REALISASI` | 2 | `FlowAction/Installments_ReadOnly.xml`, `Section/Installments_ReadOnly.xml` |
| `ASM-FW-GISFW-DATA-OCCUPATION / OCCUPATIONITEMFACIN_SECTION` | 2 | `Section/OccupationItemFacIn_Section.xml`, `Section/OccupationItemFacIn_Section_IsUW.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / CHOOSECOVERAGE` | 2 | `FlowAction/ChooseCoverage.xml`, `Section/ChooseCoverage.xml` |
| `ASM-FW-GISFW-WORK / GENERATESPREADINGLIFE_ACT` | 2 | `Activity/GenerateSpreadingLifeOR_act.xml`, `Activity/GenerateSpreadingLife_act.xml` |
| `ASM-FW-GISFW-WORK / DETAILEDMPAYMENTHEALTH` | 2 | `Activity/DetailEDMPaymentHealth.xml`, `Activity/DetailEDMPaymentHealthCorporate.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / OCCUPATIONFACOUT` | 2 | `Section/OccupationFacOut.xml`, `Section/OccupationFacOut_IsUW.xml` |
| `ASM-FW-GISFW-INT-POLICYJSON / ASM!INSERTINTOFACINSPREADLIFE_SQL` | 2 | `RDBList/InsertIntoFacinSPreadLifeMonthly_Sql.xml`, `RDBList/InsertIntoFacinSPreadLife_Sql.xml` |
| `ASM-FW-GISFW-DATA-OCCUPATION / OBJECTDTLANEKA` | 2 | `FlowAction/ViewObjectDtlAneka.xml`, `Section/ViewObjectDtlAneka.xml` |
| `ASM-FW-GISFW-DATA-PROPERTYITEM / VIEWCOVERAGEFACOUTSHOW` | 2 | `Section/ViewCoverageFacOutShow.xml`, `Section/ViewCoverageFacOutShow_isUW.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / CHANGEDISKON` | 2 | `DataTransform/ChangeDiskon.xml`, `DataTransform/ChangeDiskon_FacIn.xml` |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN / COUNTNETPREMI_ACT` | 2 | `Activity/CountNetPremi_act.xml`, `Activity/SetDueTo_act.xml` |
| `ASM-FW-GISFW-WORK / COPYPYWORKPAGE` | 2 | `Activity/CopypyWorkPage.xml`, `Activity/ViewpyWorkPage.xml` |
| `ASM-FW-GISFW-WORK / COUNTTOTALTSIPREMINUSARE_ACT` | 2 | `Activity/CountTSIRNMMultiCob_Act.xml`, `Activity/CountTotalTSIPremiNusaRe_Act.xml` |
| `ASM-FW-GISFW-WORK / CHOOSECLAUSEFIRE` | 2 | `Harness/ChooseClauseFire.xml`, `Section/ChooseClauseFire.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / COVERAGELISTFIRE_FLOWACTION` | 2 | `FlowAction/CoverageListFire_FlowAction.xml`, `FlowAction/CoverageListFire_FlowAction_IsUW.xml` |
| `ASM-FW-GISFW-WORK / CHANGEPERCENTRNM_ACT` | 2 | `Activity/changepercentrnmLife_act.xml`, `Activity/changepercentrnm_act.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / ISTJH` | 2 | `When/IsAddPA.xml`, `When/IsPLL.xml` |
| `ASM-FW-GISFW-INT-VJ_PACKAGE_AGE_KLAUSUL / ASM!GETMASTERKLAUSULAGE` | 2 | `RDBList/GetMasterBenefit.xml`, `RDBList/GetMasterKlausulAge.xml` |
| `CODE-PEGA-PDF / ATTACHTOWORK` | 2 | `Activity/AttachRISlipToWork.xml`, `Activity/AttachToWork.xml` |
| `DATA-PARTY-PERSON / INPUTDTLPARTICIPANTPA_FACIN` | 2 | `Section/InputDtlParticipantPA_FacIn.xml`, `Section/InputDtlParticipantPA_FacIn_IsUW.xml` |
| `D_COVERAGESUMMARY / #20171120T073200.324` | 2 | `DataPage/D_CoverageSummary.xml`, `DataPage/D_GolfCoverageSummary.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / HITUNGPREMIDT` | 2 | `DataTransform/HitungPremiDT.xml`, `DataTransform/HitungPremi_FacInDT.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLLAYERLISTSPREADING_FACIN` | 3 | `FlowAction/InputDtlLayerListCommision_FacIn.xml`, `FlowAction/InputDtlLayerListSpreading_FacIn.xml`, `FlowAction/InputDtlLayerListSpreading_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN / FORMULATREATYCAPACITYDESC` | 2 | `Section/FormulaTreatyCapacityDesc.xml`, `Section/FormulaTreatyCapacityDesc_IsUW.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN / ASM!VIEWMASTERCARGO` | 2 | `RDBList/ViewMasterCargo.xml`, `RDBList/ViewMasterCargo2.xml` |
| `ASM-FW-GISFW-DATA-CLAUSE / INPUTCLAUSEFIRE_VIEWDTL` | 2 | `FlowAction/InputClauseFire_ViewDtl.xml`, `Section/InputClauseFire_ViewDtl.xml` |
| `ASM-FW-GISFW-WORK / PROTECTCOVERAGE_ACT` | 2 | `Activity/ProtectCoverage_Act.xml`, `Activity/ProtectFIREMBUPA_Act.xml` |
| `@BASECLASS / INPUTNATION` | 2 | `FlowAction/InputNation.xml`, `Section/InputNation.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN / OBJECTLIST` | 2 | `Section/ObjectList.xml`, `Section/ObjectList_IsUW.xml` |
| `ASM-FW-GISFW-WORK / COPYFACRETRO_ACT` | 5 | `Activity/CopyFacRetroAnekaGolf_ACT.xml`, `Activity/CopyFacRetroCargoMBU_ACT.xml`, `Activity/CopyFacRetroFire_ACT.xml`, `Activity/CopyFacRetroPATravel_ACT.xml`, `Activity/CopyFacRetro_ACT.xml` |
| `ASM-FW-GISFW-INT-TREATYYEAR_LIFE / BROWSETREATYYEAR_LIFE_RD` | 2 | `ReportDefinition/BrowseTreatyYearOR_Life_RD.xml`, `ReportDefinition/BrowseTreatyYear_Life_RD.xml` |
| `DATA-PARTY-PERSON / INPUTCOVERAGETRAVEL` | 2 | `FlowAction/ViewCoverageTravel.xml`, `Section/ViewCoverageTravel.xml` |
| `DATA-PARTY-PERSON / INPUTCOMMISSIONPA_FACIN` | 3 | `FlowAction/InputCommissionPA_FacIn.xml`, `FlowAction/InputCommissionTravel_FacIn.xml`, `Section/InputCommissionPA_FacIn.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN / TABLEINWARDSCALE` | 2 | `FlowAction/TableInwardScale.xml`, `Section/TableInwardScale.xml` |
| `ASM-FW-GISFW-DATA-OCCUPATION / CHOOSECLASSOFCONTRACTION` | 2 | `FlowAction/ChooseClassofContraction.xml`, `Section/ChooseClassofContraction.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / PLANLISTRO` | 2 | `FlowAction/PlanListRO.xml`, `Section/PlanListRO.xml` |
| `ASM-FW-GISFW-DATA-ANEKA / PROPERTYITEMSPREADINGANEKA_FLOW` | 2 | `FlowAction/PropertyItemSpreadingAneka_Flow.xml`, `FlowAction/PropertyItemSpreadingAneka_Flow_IsUW.xml` |
| `ASM-FW-GISFW-DATA / ISCRIME` | 2 | `When/IsCrime.xml`, `When/IsEnvironmental.xml` |
| `ASM-FW-GISFW-DATA-PLAN / BENEFITLISTRO` | 2 | `FlowAction/BenefitListRO.xml`, `Section/BenefitListRO.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN / EDITMARKETING` | 2 | `FlowAction/InputCommentDeductible.xml`, `Section/InputCommentDeductible.xml` |
| `ASM-FW-GISFW-WORK / CHECKSPREADINGPROTECT_ACT` | 5 | `Activity/CheckSpreadingProtectAnekaGolf_ACT.xml`, `Activity/CheckSpreadingProtectFire_ACT.xml`, `Activity/CheckSpreadingProtectMCargoMBU_ACT.xml`, `Activity/CheckSpreadingProtectPATravel_ACT.xml`, `Activity/CheckSpreadingProtect_ACT.xml` |
| `ASM-FW-GISFW-WORK / VIEWCLAIMLIST` | 2 | `Harness/ViewClaimList.xml`, `Section/ViewClaimList.xml` |
| `DATA-PARTY-PERSON / INPUTCOMMISSIONLIFE_FACIN` | 2 | `FlowAction/InputCommissionLife_FacIn.xml`, `Section/InputCommissionLife_FacIn.xml` |
| `ASM-FW-GISFW-INT-CLAUSE_TEMPLATE / ASM!SEARCHTEMPLATEMAINCLAUSESQL` | 2 | `RDBList/SearchJobID.xml`, `RDBList/SearchJobIDSQL.xml` |
| `ASM-FW-GISFW-DATA-TRADING / INPUTDTLTRADING` | 4 | `FlowAction/InputDtlTrading.xml`, `FlowAction/InputDtlTrading_IsUW.xml`, `Section/InputDtlTrading.xml`, `Section/InputDtlTrading_IsUW.xml` |
| `ASM-FW-GISFW-WORK / CALCULTEPERSENTAGESPEADING_ACT` | 2 | `Activity/CalcultePersentageSpeadingLife_Act.xml`, `Activity/CalcultePersentageSpeading_Act.xml` |
| `ASM-FW-GISFW-WORK / COMFIRMPOLIS_NONLIFE` | 2 | `Section/ComfirmPolis_Life.xml`, `Section/ComfirmPolis_NonLife.xml` |
| `@BASECLASS / SAVECOLOR_ACT` | 3 | `Activity/SaveCity_Act.xml`, `Activity/SaveDistrict_Act.xml`, `Activity/SaveRW_Act.xml` |
| `ASM-FW-GISFW-DATA-ANEKA / SPREADINGCOVERAGEANEKA_FLOWACTION` | 3 | `FlowAction/CommisionCoverageAneka_FlowAction.xml`, `FlowAction/SpreadingCoverageAneka_FlowAction.xml`, `FlowAction/SpreadingCoverageAneka_FlowAction_IsUW.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / OBJECTITEMSUMMARY` | 2 | `Section/GolfObjectItemSummary.xml`, `Section/ObjectItemSummary.xml` |
| `ASM-FW-GISFW-DATA-ANEKA / PROTECTIONOBJECTGOLFINSURANCE_ACT` | 2 | `Activity/ProtectionObjectAviation_Act.xml`, `Activity/ProtectionObjectFidelity_Act.xml` |
| `DATA-PARTY-PERSON / SET_PERSONNOTE_DT` | 2 | `DataTransform/Set_PersonNoteDokter_DT.xml`, `DataTransform/Set_PersonNote_DT.xml` |
| `ASM-FW-GISFW-WORK / INPUTDTLOBJECT_FACIN` | 2 | `Section/InputDtlObject_FacIn.xml`, `Section/InputDtlObject_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-DATA-PROPERTYITEM / COPYCOVERAGEFROM` | 2 | `FlowAction/CopyCoverageFrom.xml`, `Section/CopyCoverageFrom.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN / TABLEOFLIMIT` | 2 | `FlowAction/TableOfLimit.xml`, `Section/TableOfLimit.xml` |
| `ASM-FW-GISFW-DATA-CAUSEOFLOSS / INPUTCAUSEOFLOSS_FACIN` | 6 | `FlowAction/InputCauseOfLoss_FacIn.xml`, `FlowAction/InputCauseOfLoss_FacIn_IsUW.xml`, `FlowAction/InputLossRecord.xml`, `Section/InputCauseOfLoss_FacIn.xml`, `Section/InputCauseOfLoss_FacIn_IsUW.xml`, `Section/InputLossRecord_Sec.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / INPUTDTLSPREADINGPA_FACIN` | 4 | `FlowAction/InputDtlSpreadingPA_FacIn.xml`, `FlowAction/InputDtlSpreadingPA_FacIn_IsUW.xml`, `Section/InputDtlSpreadingPA_FacIn.xml`, `Section/InputDtlSpreadingPA_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / INPUTCOVERAGECARGO_FACIN` | 2 | `Section/InputCoverageCargo_FacIn.xml`, `Section/InputCoverageCargo_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-WORK / CLEANACTMBU` | 3 | `Activity/CleanActAneka.xml`, `Activity/CleanActMBU.xml`, `Activity/CleanActPA.xml` |
| `ASM-FW-GISFW-DATA / ISBUILDERRISK` | 2 | `When/IsBuilderRisk.xml`, `When/IsPortRisk.xml` |
| `DATA-PARTY-PERSON / PERSONGRIDLIFE_FACIN` | 2 | `FlowAction/PersonGridLife_FacIn.xml`, `FlowAction/PersonGridLife_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-WORK / GETLIMITAKSEPTASI_ACT` | 2 | `Activity/GetLimitAkseptasi_Act.xml`, `Activity/GetLimitAkseptasi_ActFlow.xml` |
| `ASM-FW-GISFW-DATA / ISASURANSIKREDIT` | 2 | `When/IsAsuransiKredit.xml`, `When/IsCreditBriguna.xml` |
| `ASM-FW-GISFW-DATA-FACOFFER / COUNTPREMIUMNET` | 2 | `Activity/CountPremiumNet.xml`, `Activity/CountPremiumNetElse.xml` |
| `ASM-FW-GISFW-INT-ACCUMULATION_LIFE / BROWSEACCUMULATIONLIFE_RD` | 2 | `ReportDefinition/BrowseAccumulationLife_RD.xml`, `ReportDefinition/BrowseSameAccumulationLife_RD.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN / OFFERFACIN_NUSARE` | 2 | `Section/OfferFacIn_NusaRe.xml`, `Section/OfferFacIn_NusaRe_IsUW.xml` |
| `ASM-FW-GISFW-INT-CURRENCY / ASM!GETKURSLIMITSPREADING_SQL` | 2 | `RDBList/GetKursLimitByName_SQL.xml`, `RDBList/GetKursLimitSpreading_SQL.xml` |
| `ASM-FW-GISFW-DATA-OCCUPATION / SPREADINGOBJECTANEKA_FLOWACTION` | 3 | `FlowAction/CommisionObjectAneka_FlowAction.xml`, `FlowAction/SpreadingObjectAneka_FlowAction.xml`, `FlowAction/SpreadingObjectAneka_FlowAction_IsUW.xml` |
| `DATA-PORTAL / ACCUMULATIONRISK` | 2 | `Harness/AccumulationRisk.xml`, `Section/AccumulationRisk.xml` |
| `@BASECLASS / INPUTRW` | 2 | `FlowAction/InputRW.xml`, `Section/InputRW.xml` |
| `ASM-FW-GISFW-WORK / TOSENIORUW` | 2 | `When/ToJUW_A.xml`, `When/ToSeniorUW.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / CHOOSEACCUMULATION_FACIN` | 2 | `Harness/ChooseAccumulation_FacIn.xml`, `Section/ChooseAccumulation_FacIn.xml` |
| `ASM-FW-GISFW-WORK / DISABLECAUSEOFDECLINE` | 2 | `Activity/DisableCauseOfDecline.xml`, `Activity/SetCauseOfDecline_ACT.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / LOCATIONDETAIL_FLOWACTION` | 2 | `FlowAction/GolfLocationDetail_FlowAction.xml`, `FlowAction/LocationDetail_FlowAction.xml` |
| `ASM-FW-GISFW-DATA-CONVEYANCE / INPUTDTLCONVEYANCE` | 4 | `FlowAction/InputDtlConveyance.xml`, `FlowAction/InputDtlConveyance_IsUW.xml`, `Section/InputDtlConveyance.xml`, `Section/InputDtlConveyance_IsUW.xml` |
| `ASM-FW-GISFW-INT-ACCUMULATION / ASM!GETACCUMULATION_SQL` | 2 | `RDBList/GetAccumulationByNote_SQL.xml`, `RDBList/GetAccumulationProvince_SQL.xml` |
| `ASM-FW-GISFW-DATA-ANEKA / PROTECTIONOBJECTHE_ACT` | 3 | `Activity/ProtectionObjectGolfInsurance_Act.xml`, `Activity/ProtectionObjectHE_Act.xml`, `Activity/ProtectionObjectLandRig_Act.xml` |
| `ASM-FW-GISFW-WORK / ISSPVCREATE` | 3 | `When/IsSPVCreate.xml`, `When/IsSPVTreaty1.xml`, `When/IsTreaty1.xml` |
| `ASM-FW-GISFW-DATA-ANEKA / CHOOSEOBJECT_SHIP` | 2 | `FlowAction/ChooseObject_Ship.xml`, `Section/ChooseObject_Ship.xml` |
| `ASM-FW-GISFW-WORK / TOTALACCUMULATION_FACIN` | 2 | `Harness/TotalAccumulation_FacIn.xml`, `Section/TotalAccumulation_FacIn.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / OCCUPATIONLIST` | 2 | `Section/OccupationList.xml`, `Section/OccupationList_IsUW.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / CHOOSERISKADDRESS` | 3 | `FlowAction/ChooseRiskAddress.xml`, `Harness/ChooseRiskAddress.xml`, `Section/ChooseRiskAddress.xml` |
| `ASM-FW-GISFW-WORK / OFFERFACOUT` | 2 | `Flow/OfferFacOut.xml`, `FlowAction/OfferFacOut.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GENERATEIMAGEID_SQL` | 2 | `RDBList/GenerateImageID_SQL.xml`, `RDBList/Insert_T_Storage_SQL.xml` |
| `ASM-FW-GISFW-DATA-QUOTATION / CEDINGCEDANT` | 4 | `FlowAction/CedingCedant.xml`, `FlowAction/CedingCedant_IsUW.xml`, `Section/CedingCedant.xml`, `Section/CedingCedant_IsUW.xml` |
| `ASM-FW-GISFW-DATA-DEDUCTIBLE / PROTECTCURRENCYCARGO_ISUW` | 2 | `FlowAction/ProtectCurrencyCargo_IsUW.xml`, `Section/ProtectCurrencyCargo_IsUW.xml` |
| `ASM-FW-GISFW-WORK / ACCEPTNOTIFICATIONEDM` | 2 | `FlowAction/OfferFacinLifeConfirm.xml`, `Section/OfferFacinLifeNotification.xml` |
| `ASM-FW-GISFW-DATA-CARGO / INPUTCOMMISSIONMC_FACIN` | 2 | `FlowAction/InputCommissionMC_FacIn.xml`, `Section/InputCommissionMC_FacIn.xml` |
| `ASM-FW-GISFW-DATA-SPREADINGRISK / INPUTDTLLAYERSPREADING_FACIN` | 2 | `FlowAction/InputDtlLayerSpreading_FacIn.xml`, `Section/InputDtlLayerSpreading_FacIn.xml` |
| `ASM-FW-GISFW-WORK / IST1T3` | 3 | `When/IsT1T3.xml`, `When/IsT2T4.xml`, `When/IsTBonding.xml` |
| `ASM-FW-GISFW-WORK / IST1T4` | 2 | `When/IsT1T4.xml`, `When/IsT2T3.xml` |
| `ASM-FW-GISFW-DATA-GOOD / INPUTDTLGOODS` | 4 | `FlowAction/InputDtlGoods.xml`, `FlowAction/InputDtlGoods_IsUW.xml`, `Section/InputDtlGoods.xml`, `Section/InputDtlGoods_IsUW.xml` |
| `DATA-PARTY-PERSON / INPUTCOVERAGEPA` | 2 | `FlowAction/ViewCoveragePA.xml`, `Section/ViewCoveragePA.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / RISKAROUND` | 2 | `Section/RiskAround.xml`, `Section/RiskAround_IsUW.xml` |
| `ASM-FW-GISFW-DATA-VEHICLE / VEHICLEGRID_FACIN` | 3 | `FlowAction/VehicleGrid_FacIn_IsUW.xml`, `Section/VehicleGrid_FacIn.xml`, `Section/VehicleGrid_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-DATA-VEHICLE / ISFLAGOLDDATA` | 2 | `When/FlagOldData.xml`, `When/IsFlagOldData.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN / COVERAGELIST` | 2 | `Section/CoverageList.xml`, `Section/CoverageList_IsUW.xml` |
| `DATA-PORTAL / ADJUSTMENTRISKACCUMULATION` | 2 | `Harness/AdjustmentRiskAccumulation.xml`, `Section/AdjustmentRiskAccumulation.xml` |
| `@BASECLASS / INPUTKABUPATEN` | 2 | `FlowAction/InputDistrik.xml`, `FlowAction/InputKabupaten.xml` |
| `ASM-FW-GISFW-WORK / ISJUNIORUW_B` | 2 | `DataTransform/IsJuniorUW_B.xml`, `DataTransform/SetLetterNo.xml` |
| `ASM-FW-GISFW-DATA-CARGO / SHOWCOVERAGEFACOUT_FLOWACTION` | 2 | `FlowAction/ShowCoverageFacOut_FlowAction.xml`, `FlowAction/ShowCoverageFacOut_FlowAction_IsUW.xml` |
| `ASM-FW-GISFW-DATA-PROPERTYITEM / SETINDEXPROPERTY` | 2 | `DataTransform/setIndexProperty.xml`, `DataTransform/setIndexPropertyFacOut_DT.xml` |
| `ASM-FW-GISFW-WORK / SETVALIDATEDATE_POSTACT` | 2 | `Activity/SetDataFacOut_Act.xml`, `Activity/SetValidateDate_PostAct.xml` |
| `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETLIMITAKSEPTASINONPREFER_SQL` | 3 | `RDBList/GetLimitAkseptasiNonPreferBanding_SQL.xml`, `RDBList/GetLimitAkseptasiNonPreferJUWA_SQL.xml`, `RDBList/GetLimitAkseptasiPreferedComm_SQL.xml` |
| `@BASECLASS / INPUTKOTA` | 2 | `FlowAction/InputKotaBranchName.xml`, `FlowAction/InputKotaProvince.xml` |
| `ASM-FW-GISFW-WORK / EMAILSECTION` | 2 | `Section/EmailSection.xml`, `Section/EmailSection_Life.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / SETPREMIUMDEDUCTIBLE_DT` | 2 | `DataTransform/SetAdditionalPremiumDeductible_DT.xml`, `DataTransform/SetPremiumDeductible_DT.xml` |
| `ASM-FW-GISFW-INT-PROPORTIONALARRG / ASM!GETTREATYNAME_SQL` | 3 | `RDBList/GetTreatyName.xml`, `RDBList/GetTreatyName_SQL.xml`, `RDBList/GetTreatyName_SQL2.xml` |
| `ASM-FW-GISFW-DATA-CORRESPONDENCE / CORRESPONDENCECONTENT` | 2 | `FlowAction/CorrespondenceContent.xml`, `Section/CorrespondenceContent.xml` |
| `DATA-PARTY-PERSON / INPUTDTLPARTICIPANTTRAVEL` | 2 | `Section/InputDtlParticipantTravel.xml`, `Section/InputDtlParticipantTravel_IsUW.xml` |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN / DETAILPOLICYTREATYINNONPROPORTIONAL` | 3 | `Section/DetailPolicyTreatyInNonProportional.xml`, `Section/DetailPolicyTreatyInNonProportionalEDM.xml`, `Section/DetailPolicyTreatyOutNonProportional.xml` |
| `ASM-FW-GISFW-DATA-POLICY / UPLOADCSVPOLICYMEMBER_POSTACT` | 2 | `Activity/UploadCSVPolicyMemberEDM_PostAct.xml`, `Activity/UploadCSVPolicyMember_PostAct.xml` |
| `ASM-FW-GISFW-DATA-PROPERTYITEM / PROPERTYITEMCOVERAGEFACIN_FLOWACTION` | 2 | `FlowAction/PropertyItemCoverageFacIn_FlowAction_IsUW.xml`, `FlowAction/PropertyItemCoverageSpreadingFacIn_FlowAction.xml` |
| `ASM-FW-GISFW-DATA-OCCUPATION / INPUTOKUPASIANEKA_FACIN` | 2 | `Section/InputOkupasiAneka_FacIn.xml`, `Section/InputOkupasiAneka_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / TOTALACCUMULATIONDTL` | 2 | `Harness/TotalAccumulationDtl.xml`, `Section/TotalAccumulationDtl.xml` |
| `ASM-FW-GISFW-DATA-OCCUPATION / INPUTOKUPASIANEKA` | 2 | `Section/InputOkupasiAneka_GCNM.xml`, `Section/ViewOkupasiAneka.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / SPREADINGITEMLAYER` | 3 | `Section/SpreadingItemLayer.xml`, `Section/SpreadingItemLayerCommision.xml`, `Section/SpreadingItemLayer_IsUW.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / INPUTDEDUCTIBLE_GCNM` | 2 | `Section/InputDeductibleHull_GCNM.xml`, `Section/InputDeductible_GCNM.xml` |
| `ASM-FW-GISFW-DATA-CARGO / INPUTDTLCOVERAGE_FACIN` | 2 | `Section/InputDtlCoverage_FacIn.xml`, `Section/InputDtlCoverage_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-INT-POLICYJSON / ASM!INSERTJSON_JSONPOLIS_FACIN` | 2 | `RDBList/INSERTJSON_JSONPOLISEDM_FACIN.xml`, `RDBList/INSERTJSON_JSONPOLISMONITORING_FACIN.xml` |
| `ASM-FW-GISFW-DATA-CLAUSE / INPUTCLAUSE_VIEWDTL_PREACT` | 2 | `Activity/InputClause_ViewArg_PreAct.xml`, `Activity/InputClause_ViewDtl_PreAct.xml` |
| `DATA-PARTY-PERSON / PERSONGRIDTRAVEL` | 2 | `FlowAction/PersonGridTravel.xml`, `FlowAction/PersonGridTravel_IsUW.xml` |
| `DATA-PORTAL / RISKACCUMULATIONREPORT` | 2 | `Harness/RiskAccumulationReport.xml`, `Section/RiskAccumulationReport.xml` |
| `ASM-FW-GISFW-WORK / CHANGESPREADINGPERCENTAGE_ACT` | 2 | `Activity/ChangeSpreadingPercentageLife_Act.xml`, `Activity/ChangeSpreadingPercentage_Act.xml` |
| `ASM-FW-GISFW-WORK / INPUTINWARDFACULTATIVE` | 3 | `Section/InputEndorsement.xml`, `Section/InputInwardFacultative.xml`, `Section/InputInwardFacultative_2.xml` |
| `DATA-PARTY-PERSON / INPUTSPREADINGTRAVEL_FACIN` | 2 | `FlowAction/InputSpreadingTravel_FacIn_IsUW.xml`, `Section/InputSpreadingTravel_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-INT-POLICYJSON / ASM!UPDATEPOLISENDORSEMENT_SQL` | 2 | `RDBList/UpdateJsonOffer_SQL.xml`, `RDBList/UpdatePolisEndorsement_SQL.xml` |
| `ASM-FW-GISFW-DATA-QUOTATION / SOURCEHIERARKI` | 2 | `Section/CedingCoHierarki.xml`, `Section/SourceHierarki.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN / VIEWCSVCREDIT` | 2 | `Section/ViewCSVCredit.xml`, `Section/ViewCSVMarine.xml` |
| `ASM-FW-GISFW-WORK / SHOWPOLIS` | 2 | `FlowAction/ShowPolis.xml`, `Harness/ShowPolis.xml` |
| `ASM-FW-GISFW-DATA-PROPERTYITEM / VIEWCOVERAGEFACOUT` | 2 | `Harness/ViewCoverageFacOut.xml`, `Harness/ViewCoverageFacOut_isUW.xml` |
| `ASM-FW-GISFW-INT-SHIP / SELECTSHIP` | 2 | `Harness/SelectShip.xml`, `Section/SelectShip.xml` |
| `ASM-FW-GISFW-WORK / CHECKDTLOBJECT` | 2 | `Activity/CheckDtlObject.xml`, `Activity/CheckDtlObjectTravel.xml` |
| `ASM-FW-GISFW-INT-CURRENCY / ASM!GETDATACURRENCYBYNAME_SQL` | 2 | `RDBList/GetCurrencyIDByName.xml`, `RDBList/GetDataCurrencyByName_SQL.xml` |
| `ASM-FW-GISFW-WORK / INWARDFACULTATIVE` | 2 | `FlowAction/InwardFacultative.xml`, `FlowAction/InwardFacultative_IsUW.xml` |
| `CODE-PEGA-LIST / GETCOVERAGE_ACT` | 2 | `Activity/GetCoverage_Act.xml`, `Activity/GetGolfCoverage_Act.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN / SHOWCEDINGCOLIST` | 2 | `Harness/ShowCedingCoList.xml`, `Section/ShowCedingCoList.xml` |
| `ASM-FW-GISFW-DATA-PROPERTYITEM / COPYCOVERAGE_ACT` | 2 | `Activity/CopyCoverageObjItem_Act.xml`, `Activity/CopyCoverage_Act.xml` |
| `ASM-FW-GISFW-DATA-FACOFFER / COPYALLOBJ_ACT` | 4 | `Activity/CopyAllObjFacOutFireAneka_ACT.xml`, `Activity/CopyAllObjFacOutGolfCargo_ACT.xml`, `Activity/CopyAllObjFacOutPAMBU_ACT.xml`, `Activity/CopyAllObj_ACT.xml` |
| `ASM-FW-GISFW-WORK / INPUTOBJECTSUBCONTRACT_FACIN` | 2 | `Section/InputObjectSubContract_FacIn.xml`, `Section/InputObjectSubContract_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-WORK / OFFERSTATUSFACOUT` | 2 | `Section/OfferStatusFacOut.xml`, `Section/OldOfferStatusFacOut.xml` |
| `ASM-FW-GISFW-WORK / SETAKSEPTASIPROPOSAL` | 2 | `DataTransform/SetAkseptasiCeding_DT.xml`, `DataTransform/SetAkseptasiProposal.xml` |
| `CODE-PEGA-LIST / GETLOCATIONTOPRISK_ACT` | 2 | `Activity/GetGolfLocationTopRisk_Act.xml`, `Activity/GetLocationTopRisk_Act.xml` |
| `ASM-FW-GISFW-DATA-PROPERTYITEM / INPUTCOVERAGEFIRE` | 2 | `Section/InputCoverageFire.xml`, `Section/InputCoverageFire_IsUW.xml` |
| `ASM-FW-GISFW-WORK / PAYMENTCURRENCYLISTLIFE` | 2 | `Section/PaymentCurrencyListLife.xml`, `Section/PaymentCurrencyListLife_IsUW.xml` |
| `ASM-FW-GISFW-DATA-PROPERTY / INPUTDTLOBJFIRE` | 2 | `Section/InputDtlObjFire.xml`, `Section/InputDtlObjFire_IsUW.xml` |
| `ASM-FW-GISFW-DATA-PROPERTY / INPUTOBJECTDTLFIRE` | 3 | `Section/InputObjectDtlFire.xml`, `Section/InputObjectDtlFire_GCNM.xml`, `Section/InputObjectDtlFire_IsUW.xml` |
| `ASM-FW-GISFW-WORK / SHOWLOCATIONFACOUT` | 2 | `Activity/ShowLocationFacOut.xml`, `Activity/ShowSumInsuredTSIFacOut.xml` |
| `ASM-FW-GISFW-INT-POLICYJSON / ASM!INSERTFACINLIFE_SQL` | 2 | `RDBList/InsertFacinLifeMonthly_Sql.xml`, `RDBList/InsertFacinLife_Sql.xml` |
| `ASM-FW-GISFW-DATA-ANEKA / PROTECTIONOBJECTCONTRACTORPLANT_ACT` | 2 | `Activity/ProtectionObjectBoilerPresure_Act.xml`, `Activity/ProtectionObjectMB_Act.xml` |
| `ASM-FW-GISFW-DATA-CARGO / INPUTSPREADINGMARINECARGO_FACIN` | 2 | `FlowAction/InputSpreadingMarineCargo_FacIn.xml`, `FlowAction/InputSpreadingMarineCargo_FacIn_ISUW.xml` |
| `ASM-FW-GISFW-WORK / INPUTOFFERFACINENGINEER_PREACT` | 2 | `Activity/InputOfferFacInEngineerUW_preACT.xml`, `Activity/InputOfferFacInEngineer_preACT.xml` |
| `ASM-FW-GISFW-WORK / COUNTPREMIANDTSINUSANTARARE_ACT` | 3 | `Activity/CountPremiAndTSINusantaraReEDM_ACT.xml`, `Activity/CountPremiAndTSIRNMFireAnekaGolf_ACT.xml`, `Activity/CountPremiAndTSIRNMPALifeCargo_ACT.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / OBJECTDETAILS` | 2 | `Section/ObjectDetails.xml`, `Section/ObjectDetails_IsUW.xml` |
| `ASM-FW-GISFW-WORK / VIEWFOLLOWINGNB` | 2 | `Harness/ViewFollowingNB.xml`, `Section/GridViewFollowingNB.xml` |
| `ASM-FW-GISFW-DATA-PROPERTYITEM / COPYCOVERAGETOALLLOCATION_ACT` | 2 | `Activity/CopyCoverageToAllLocation_Act.xml`, `Activity/CopyDeductibleToAllLocation_Act.xml` |
| `ASM-FW-GISFW-DATA-OCCUPATION / OCCUPATIONITEMFACIN_FLOWACTION` | 2 | `FlowAction/OccupationItemFacIn_FlowAction.xml`, `FlowAction/OccupationItemFacIn_FlowAction_IsUW.xml` |
| `ASM-FW-GISFW-WORK / SETASKPROPOSAL_DT` | 2 | `DataTransform/SetAskProposal_DT.xml`, `DataTransform/SetBandingProposal_DT.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / VIEWOLDDEDUCTIBLE` | 2 | `Harness/ViewOldDeductible.xml`, `Section/ViewOldDeductible.xml` |
| `ASM-FW-GISFW-INT-ACCUMULATION_LIFE / ASM!GETACCUMULATION_SQL` | 2 | `RDBList/GetAccumulation_Sql.xml`, `RDBList/RetrieveDataAccumulationLife_SQL.xml` |
| `ASM-FW-GISFW-DATA-QUOTATION / CEDINGCOMPANY` | 2 | `Harness/CedingCedant.xml`, `Harness/CedingCompany.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / FILLPREMIGOLF` | 2 | `Activity/CountPremiCoverageAneka.xml`, `Activity/FillPremiGolf.xml` |
| `DATA-PARTY-PERSON / PERSONGRIDPA` | 3 | `FlowAction/PersonGridDM.xml`, `FlowAction/PersonGridPA_FacIn.xml`, `FlowAction/ViewPersonGridPA.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / LOCATIONDETAIL` | 2 | `Section/GolfLocationDetail.xml`, `Section/LocationDetail.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / SPREADINGITEMANEKA` | 4 | `FlowAction/SpreadingItemAneka.xml`, `FlowAction/SpreadingItemAneka_IsUW.xml`, `Section/CommisionItemAneka.xml`, `Section/SpreadingItemAneka_IsUW.xml` |
| `ASM-FW-GISFW-WORK / SETTOJSONOFFER_ACT` | 2 | `Activity/SaveToProduction_ACT.xml`, `Activity/SetToJsonOffer_ACT.xml` |
| `ASM-FW-GISFW-DATA-OCCUPATION / OCCUPATIONANEKAGRID` | 3 | `FlowAction/OccupationAnekaGridGCNM.xml`, `FlowAction/OccupationAnekaGridView.xml`, `FlowAction/OccupationAnekaGrid_FacIn.xml` |
| `ASM-FW-GISFW-WORK / INPUTDTLPAYMENT_PREDT` | 2 | `DataTransform/InputDtlPaymentFire_PreDT.xml`, `DataTransform/InputDtlPayment_PreDT.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN / COVERAGESPREADINGLIST` | 3 | `Section/CoverageCommisionList.xml`, `Section/CoverageSpreadingList.xml`, `Section/CoverageSpreadingList_IsUW.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY / INPUTDTLPAYMENT_FACIN` | 2 | `Section/InputDtlPayment_FacIn.xml`, `Section/InputDtlPayment_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-DATA-SPLIT / MAPPINGOUTGOINDEX` | 2 | `DecisionTable/MappingOutgoIndex.xml`, `DecisionTable/MappingOutgoIndex2.xml` |
| `ASM-FW-GISFW-DATA-VEHICLE / VEHICLEGRID` | 4 | `FlowAction/VehicleGridFacOut.xml`, `FlowAction/VehicleGrid_FacIn.xml`, `Section/VehicleGrid.xml`, `Section/ViewVehicleGridFacOut.xml` |
| `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETNOPOLIS_SQL` | 2 | `RDBList/CekFacin_Sql.xml`, `RDBList/GetNopolis_Sql.xml` |
| `ASM-FW-GISFW-DATA-DEDUCTIBLE / PROTECTCURRENCY` | 4 | `FlowAction/ProtectCurrency.xml`, `FlowAction/ProtectCurrency_IsUW.xml`, `Section/ProtectCurrency.xml`, `Section/ProtectCurrency_IsUW.xml` |
| `ASM-FW-GISFW-DATA-FACOFFER / ADDFACOFFERLIST` | 4 | `FlowAction/AddFacOfferList_IsUW.xml`, `Section/AddFacOfferList.xml`, `Section/AddFacOfferList2.xml`, `Section/AddFacOfferList_IsUW.xml` |
| `ASM-FW-GISFW-WORK / ISTABACTIVE_SPREADING` | 2 | `When/IsTABActive_Cedant.xml`, `When/IsTABActive_SpreadingLife.xml` |
| `ASM-FW-GISFW-WORK / COPYTOPOLICY` | 2 | `Activity/CopyToPolicy.xml`, `Activity/CopyToPolicyList.xml` |
| `ASM-FW-GISFW-DATA-COVERAGE / COUNTPREMI_ACT` | 2 | `Activity/CountPremiRate_ACT.xml`, `Activity/CountPremi_ACT.xml` |
| `ASM-FW-GISFW-DATA-DEDUCTIBLE / INPUTDTLDEDUCTIBLE_FACIN` | 2 | `FlowAction/InputDtlDeductible_FacIn.xml`, `FlowAction/InputDtlDeductible_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-INT-POLICYJSON / ASM!INSERTTREATYPRODUCTION_SQL` | 2 | `RDBList/InsertTreatyProd_Sql.xml`, `RDBList/InsertTreatyProduction_Sql.xml` |
| `ASM-FW-GISFW-DATA-DEDUCTIBLE / INPUTDEDUCTIBLE_FACIN` | 2 | `Section/InputDeductible_FacIn.xml`, `Section/InputDeductible_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / PROPERTY_FLOWACTION` | 2 | `FlowAction/Property_FlowAction.xml`, `FlowAction/Property_FlowAction_IsUW.xml` |
| `ASM-FW-GISFW-DATA / ISUW` | 2 | `When/IsFacRetro.xml`, `When/IsUW.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / CAUSEOFLOSS_FACIN` | 2 | `Section/CauseOfLoss_FacIn.xml`, `Section/CauseOfLoss_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-DATA-FACOFFER / INSERTFACOUTPRODUCTION` | 3 | `Activity/InsertFacoutProduction.xml`, `Activity/InsertFacoutProductionEDM.xml`, `Activity/InsertFacoutProductionEDMCurr.xml` |
| `ASM-FW-GISFW-DATA-ANEKA / SPREADINGCOVERAGEANEKA` | 3 | `Section/CommisionCoverageAneka.xml`, `Section/SpreadingCoverageAneka.xml`, `Section/SpreadingCoverageAneka_IsUW.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN / LIMITTREATY` | 2 | `FlowAction/LimitTreaty.xml`, `Section/LimitTreaty.xml` |
| `ASM-FW-GISFW-DATA-ANEKA / INPUTDTLOBJECTANEKA` | 2 | `FlowAction/ViewDtlObjectAneka.xml`, `Section/ViewDtlObjectAneka.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE / OCCUPATIONFACOUT_FLOWACTION` | 2 | `FlowAction/OccupationFacOut_FlowAction.xml`, `FlowAction/OccupationFacOut_FlowAction_IsUW.xml` |
| `ASM-FW-GISFW-DATA-INSTALLMENT / INSTALLMENTLIST` | 2 | `FlowAction/InstallmentList.xml`, `Section/InstallmentList.xml` |
| `ASM-FW-GISFW-WORK / CHANGEEMAILTEXT` | 2 | `Activity/ChangeEmailText.xml`, `Activity/ChangeEmailTextCeding_ACT.xml` |
| `ASM-FW-GISFW-DATA-DEDUCTIBLE / INPUTDTLDEDUCTIBLEFIRE_FACIN` | 2 | `FlowAction/InputDtlDeductibleFire_FacIn.xml`, `FlowAction/InputDtlDeductibleFire_FacIn_IsUW.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` | 2 | `RDBList/GetTokenStorage_SQL.xml`, `RDBList/Update_T_Storage_SQL.xml` |
| `ASM-FW-GISFW-DATA-OFFERFACIN-OFFERFEALIST / INPUTFEA` | 3 | `FlowAction/InputFEA.xml`, `FlowAction/InputFEA_IsUW.xml`, `Section/InputFEA_IsUW.xml` |
| `ASM-FW-GISFW-DATA-SUBCONTRACT / CHOOSESUBCONTRACT` | 4 | `FlowAction/ChooseSubContract.xml`, `FlowAction/ChooseSubContract_FacIn.xml`, `Section/ChooseSubContract.xml`, `Section/ChooseSubContract_FacIn.xml` |
| `ASM-FW-GISFW-DATA-FACOFFER / GETRISLIPDATA_ACT` | 2 | `Activity/CopyObjectRiSlip.xml`, `Activity/GetRISlipData_ACT.xml` |
| `ASM-FW-GISFW-WORK / TODIVHEAD` | 3 | `When/ToDirMarketing.xml`, `When/ToKadivFacultative.xml`, `When/ToKadivTeknik.xml` |
| `ASM-FW-GISFW-DATA-PROPERTYITEM / PROPERTYITEMFACIN_FLOWACTION` | 3 | `FlowAction/PropertyItemCoverageFacIn_FlowAction.xml`, `FlowAction/PropertyItemFacIn_FlowAction.xml`, `FlowAction/PropertyItemFacIn_FlowAction_IsUW.xml` |
| `ASM-FW-GISFW-WORK / TOMANAGERTEKNIK` | 4 | `When/ToDepHeadUW.xml`, `When/ToDeptHeadUWLife.xml`, `When/ToManagerTeknik.xml`, `When/ToTREATYDEPTHEAD.xml` |
| `ASM-FW-GISFW-INT / REASVIEWATTACHMENT` | 2 | `FlowAction/ReasViewAttachment.xml`, `Section/ReasViewAttachment.xml` |
| `ASM-FW-GISFW-DATA-CARGO / INPUTDTLCARGO` | 4 | `FlowAction/InputDtlCargo_FacIn.xml`, `FlowAction/ViewDtlCargo.xml`, `Section/InputDtlCargo.xml`, `Section/ViewDetailCargo.xml` |
| `ASM-FW-GISFW-INT-V_JN_OBJ_ITEM / BROWSEV_JN_OBJ_ITEM` | 2 | `ReportDefinition/BrowseV_JN_OBJ_ITEM.xml`, `ReportDefinition/GetObjectItem.xml` |
| `ASM-FW-GISFW-WORK / SETDATAFACOUT_ACT` | 4 | `Activity/SetDataFacOutAnekaGolf_Act.xml`, `Activity/SetDataFacOutCargoMBU_Act.xml`, `Activity/SetDataFacOutFire_Act.xml`, `Activity/SetDataFacOutPATravel_Act.xml` |
| `D_LOCATIONSUMMARY / #20171120T044142.723` | 2 | `DataPage/D_GolfLocationSummary.xml`, `DataPage/D_LocationSummary.xml` |
| `ASM-FW-GISFW-WORK / INPUTCAUSEOFDECLINE` | 2 | `Harness/InputCauseOfDecline.xml`, `Section/InputCauseOfDecline.xml` |
| `ASM-FW-GISFW-DATA / ISPRODUCTSLIABILITY` | 2 | `When/IsProductsLiability.xml`, `When/IsWhetherIndex.xml` |
| `ASM-FW-GISFW-DATA-DEDUCTIBLE / PROTECTCURRENCYCARGO` | 2 | `FlowAction/ProtectCurrencyCargo.xml`, `Section/ProtectCurrencyCargo.xml` |
| `ASM-FW-GISFW-WORK / BROWSEM_OKUPASI_ANEKA` | 2 | `Activity/BrowseCoverageSupport.xml`, `Activity/BrowseM_OKUPASI_ANEKA.xml` |
| `ASM-FW-GISFW-WORK / SAVETREATYPRODUCTION_ACT` | 4 | `Activity/SaveFacinProdFireNB_Act.xml`, `Activity/SaveFacinRNWProd_Act.xml`, `Activity/SaveOfferProduction_Act.xml`, `Activity/SaveTreatyProduction_Act.xml` |
| `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETLIMITAKSEPTASI_SQL` | 6 | `RDBList/GetLimitAccEngineeringUW_SQL.xml`, `RDBList/GetLimitAkseptasiBanding_SQL.xml`, `RDBList/GetLimitAkseptasiJUWA_SQL.xml`, `RDBList/GetLimitAkseptasiLife_SQL.xml`, `RDBList/GetLimitAkseptasiNonPrefer_SQL.xml`, `RDBList/GetLimitAkseptasi_SQL.xml` |
| `ASM-FW-GISFW-DATA-FACOFFER / PRINTRISLIP` | 2 | `Harness/PrintRISlip.xml`, `Section/PrintRISlip_Endorsement.xml` |
| `ASM-FW-GISFW-WORK / ISAPPROVED` | 2 | `DecisionTable/isApproved.xml`, `When/isApproved.xml` |
| `ASM-FW-GISFW-WORK / INPUTENDORSEMENTDTL` | 2 | `Section/InputEndorsementDtl.xml`, `Section/OldDataEndorsementDtl.xml` |

## 6. Objek database yang disentuh

Diambil dari `<pyBrowseSQL>` pada rule `RDBList`. **Nama objek apa adanya dari SQL** —
kualifikasi skema tidak konsisten di korpus. **Tipe kolom, constraint, dan relasi tidak
diketahui** (OQ-001).

| Tabel/view (apa adanya di SQL) | Dirujuk oleh |
| --- | ---: |
| `JSON_POLIS` | 18 rule |
| `FACINPRODUCTION` | 7 rule |
| `ACCUMULATION` | 6 rule |
| `RW` | 6 rule |
| `CURRENCY` | 5 rule |
| `OCCUPATION` | 5 rule |
| `POOLDATA.JSON_POLIS` | 5 rule |
| `POOLDATA.REINSURANCETYPE` | 5 rule |
| `TREATYBUSINESS` | 5 rule |
| `MARKETINGOFFICER` | 4 rule |
| `POOLDATA.CLIENT` | 4 rule |
| `POOLDATA.M_LIMIT_PROPERTYY` | 4 rule |
| `PROPORTIONALARRG` | 4 rule |
| `CATEGORY_ATTACH_REAS` | 3 rule |
| `FACOUTPRODUCTION` | 3 rule |
| `HISTORYAKSEPTASIPEGA` | 3 rule |
| `JSON_FOLLOWING` | 3 rule |
| `M_CLAUSE` | 3 rule |
| `M_LIMIT_FINANCIALINS` | 3 rule |
| `POOLDATA.FACINPRODUCTION` | 3 rule |
| `POOLDATA.M_LIMIT_NONPROPANDENGG` | 3 rule |
| `T_STORAGE_IMAGE` | 3 rule |
| `DATAPEGA.PC_HISTORY_ASM_FW_GISFW_WORK` | 2 rule |
| `FIRE.M_BI_INDEMNITY@ASMD.SINARMAS.CO.ID` | 2 rule |
| `GENERAL.V_JOB_PA` | 2 rule |
| `JSON_OFFER` | 2 rule |
| `M_ARGCLAUSEFIRE` | 2 rule |
| `M_LIMIT_ENGINEERINGG` | 2 rule |
| `M_LIMIT_PROPERTY_NON_PREFERREDD` | 2 rule |
| `POOLDATA.AGENT` | 2 rule |
| `POOLDATA.CAUSEOFDECLINE` | 2 rule |
| `POOLDATA.FACINLIFE` | 2 rule |
| `POOLDATA.FACINOFFER` | 2 rule |
| `POOLDATA.FACINSPREADLIFE` | 2 rule |
| `POOLDATA.JSON_OFFER` | 2 rule |
| `POOLDATA.KODE_PRODUKSI` | 2 rule |
| `POOLDATA.RW` | 2 rule |
| `RATE_LIFE` | 2 rule |
| `REINSURANCETYPE` | 2 rule |
| `TREATYEXCHANGEYEARLY` | 2 rule |
| `V_JN_OBJ_ITEM` | 2 rule |
| `ZONES` | 2 rule |
| `AGENT` | 1 rule |
| `BUSINESS` | 1 rule |
| `CLIENT` | 1 rule |
| `CONSTRUCTION` | 1 rule |
| `COVERAGE` | 1 rule |
| `COVERAGE_FACIN` | 1 rule |
| `CZONE` | 1 rule |
| `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` | 1 rule |
| `DETAIL_INVOICE` | 1 rule |
| `FACINOFFER` | 1 rule |
| `FIRE.M_FLOOD_RATE@ASMD.SINARMAS.CO.ID` | 1 rule |
| `GENERAL.LST_ASURADUR` | 1 rule |
| `GENERAL.LST_DET_CABANG` | 1 rule |
| `JSON_POLIS_ERROR` | 1 rule |
| `LST_KURS_STANDARD@ASMD.SINARMAS.CO.ID` | 1 rule |
| `M_ACCUMULATION_LIFE` | 1 rule |
| `M_BI_INDEMNITY` | 1 rule |
| `M_BRANCH` | 1 rule |
| `M_BUSINESSFIELD` | 1 rule |
| `M_COUNTRY` | 1 rule |
| `M_CURRENCY` | 1 rule |
| `M_EQS_RATE@ASMD.SINARMAS.CO.ID` | 1 rule |
| `M_FLEXAS_RATE@ASMD.SINARMAS.CO.ID` | 1 rule |
| `M_FLOOD_AREA@ASMD.SINARMAS.CO.ID` | 1 rule |
| `M_KONSTRUKSI` | 1 rule |
| `M_LIMIT_LIFE` | 1 rule |
| `M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL` | 1 rule |
| `M_PROMPT_AI` | 1 rule |
| `M_REINSURANCETYPE` | 1 rule |
| `M_RISK_LOSS_PROFILE` | 1 rule |
| `M_RSMD_RATE@ASMD.SINARMAS.CO.ID` | 1 rule |
| `M_TERORISME_RATE@ASMD.SINARMAS.CO.ID` | 1 rule |
| `M_TREATYBUSINESS` | 1 rule |
| `OBJECTITEMTYPE` | 1 rule |
| `POOLDATA.BUSINESS` | 1 rule |
| `POOLDATA.CURRENCY` | 1 rule |
| `POOLDATA.DATAKLAIM` | 1 rule |
| `POOLDATA.ERRORFACINPROD` | 1 rule |
| `POOLDATA.FACINOFFERLIFE` | 1 rule |
| `POOLDATA.FACINPERFORMANCE` | 1 rule |
| `POOLDATA.HISTORYAKSEPTASIPRODUCTION` | 1 rule |
| `POOLDATA.JSON_FOLLOWING` | 1 rule |
| `POOLDATA.JSON_POLIS_MONITORING` | 1 rule |
| `POOLDATA.MASTERCARGO` | 1 rule |
| `POOLDATA.MONITORING_PROD_LOG` | 1 rule |
| `POOLDATA.M_ACCUMULATION_LIFE` | 1 rule |
| `POOLDATA.M_EQS_MULTIPLIER` | 1 rule |
| `POOLDATA.M_LIMIT_ENGINEERINGG` | 1 rule |
| `POOLDATA.M_LIMIT_PROPERTY_NON_PREFERREDD` | 1 rule |
| `POOLDATA.M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL` | 1 rule |
| `POOLDATA.M_PLANTRAVEL` | 1 rule |
| `POOLDATA.M_RISK` | 1 rule |
| `POOLDATA.M_TREATY_IN` | 1 rule |
| `POOLDATA.M_TREATY_IN_EDM` | 1 rule |
| `POOLDATA.OCCUPATION` | 1 rule |
| `POOLDATA.PROPORTIONALARRG` | 1 rule |
| `POOLDATA.TANGGAL_CLOSING` | 1 rule |
| `POOLDATA.TREATYEXCHANGEYEARLY` | 1 rule |
| `POOLDATA.TREATYGROUP` | 1 rule |
| `POOLDATA.TREATYINPRODUCTION` | 1 rule |
| `POOLDATA.TREATYPRODUCTION_BACKUP` | 1 rule |
| `POOLDATA.T_FOLDER_IMAGE` | 1 rule |
| `POOLDATA.VJ_BENEFIT_PROPERTY` | 1 rule |
| `POOLDATA.VJ_PACKAGE_AGE_KLAUSUL_DM` | 1 rule |
| `POOLDATA.VJ_PKG_BENEFIT` | 1 rule |
| `POOLDATA.ZONEOJK` | 1 rule |
| `PREP_DATA` | 1 rule |
| `PREP_PARAM` | 1 rule |
| `SPREADSYARIAH` | 1 rule |
| `TREATYCONTRACT` | 1 rule |
| `TREATYPRODUCTION_BACKUP` | 1 rule |
| `TREATYYEAR` | 1 rule |
| `WITH` | 1 rule |
| `WORDLIST` | 1 rule |

### Stored procedure / function yang dipanggil

**Body prosedur ini TIDAK ADA di korpus** — logika di dalamnya tidak dapat direkonstruksi
dari XML. Lihat OQ-002.

| Procedure | Dipanggil oleh |
| --- | --- |
| `DBMS_LOB.CREATETEMPORARY` | `INSERTJSON_JSONPOLISMONITORING_FACIN.xml`, `SaveNewAccumulation_SQL.xml`, `UpdateMasterAccumulatedType.xml`, `UpdateMasterAccumulation.xml`, `UpdateMasterBranch.xml`, `UpdateMasterCity.xml`, `UpdateMasterDistrict.xml`, `UpdateMasterNation.xml`, `UpdateMasterProvince.xml`, `Update_sql.xml` |
| `POOLDATA.RDBMASTERPROVINCE` | `UpdateMasterProvince.xml` |
| `POOLDATA.PEGA_M_ACCUMULATION_LIFE` | `SaveNewAccumulation_SQL.xml`, `Update_sql.xml` |
| `POOLDATA.INSERTJSONPOLISMONITORING` | `INSERTJSON_JSONPOLISMONITORING_FACIN.xml` |
| `POOLDATA.GETCURRENCYSTANDARD` | `CurrencyStandard.xml` |
| `POOLDATA.FACINFORBACKUP` | `MachingDataFacin_Sql.xml` |
| `POOLDATA.PEGA_DELETE_ERROR_KONVERSI` | `DeleteDataProduction.xml` |
| `POOLDATA.RDBMASTERSHIP` | `UpdateMasterShip.xml` |
| `POOLDATA.RDBMASTERACCUMULATEDTYPE` | `UpdateMasterAccumulatedType.xml` |
| `POOLDATA.PEGA_JSON_POLIS_TREATYIN` | `SavePolisTreatyIn_SQL.xml` |
| `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` | `GetSequenceNumber_SQL.xml` |
| `POOLDATA.INSERTUPDATERISKADDRESS` | `UpdateMasterRiskAddress_SQL.xml` |
| `POOLDATA.RDBMASTERACCUMULATION` | `UpdateMasterAccumulation.xml` |
| `POOLDATA.RDBMASTERNATION` | `UpdateMasterNation.xml` |
| `GENERAL.F_GET_NM_ASURADUR` | `SearchCoinsSQL.xml` |
| `POOLDATA.PEGA_TREATY_IN` | `SaveTreatyIn.xml` |
| `POOLDATA.PEGA_MARKETINGOFFICER` | `UpdateMasterMarketingOfficer.xml` |
| `MBU.F_CEK_HURUF` | `GetNoRangkaMesinDiff.xml` |
| `POOLDATA.PEGA_M_JSON_OFFER` | `SaveOfferJson_SQL.xml` |
| `POOLDATA.RDBMASTERBRANCH` | `UpdateMasterBranch.xml` |
| `UTL_MATCH.EDIT_DISTANCE_SIMILARITY` | `SearchAccumulationbypersetase_SQL.xml` |
| `POOLDATA.RDBINSERTCLIENT` | `INSERTCONORGJSON_MCLIENT.xml` |
| `POOLDATA.GENERATE_FACRETRO_NO` | `GenerateRISlipNumber.xml` |
| `POOLDATA.RDBMASTERRW` | `UpdateMasterRW.xml` |
| `POOLDATA.RDBMASTERDISTRICT` | `UpdateMasterDistrict.xml` |
| `POOLDATA.RDBMASTERCITY` | `UpdateMasterCity.xml` |
| `POOLDATA.INSERTUPDATECEDINGPRODUCTION` | `InsertCedingProduction_SQL.xml` |
| `POOLDATA.INSERTJSONPOLIS` | `INSERTJSON_JSONPOLISEDM_FACIN.xml`, `INSERTJSON_JSONPOLIS_FACIN.xml` |
| `POOLDATA.GET_TOKEN_STORAGE` | `GetTokenStorage_SQL.xml` |

## 7. Integrasi eksternal (`ConnectREST`)

Alamat tujuan dicatat sebagai **konfigurasi**, bukan URL literal. `baseURL=SETTING` berarti
alamat diambil dari rule `SystemSettings` yang disebut di kolom Setting; nilai sebenarnya
**tidak ada di korpus**. Kolom terakhir menandai rule yang justru memuat URL literal —
itu temuan, bukan fakta arsitektur (lihat `_summary.md`).

| File | Class | pyServiceName | Sumber base URL | Setting | URL literal di dalam rule? |
| --- | --- | --- | --- | --- | --- |
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `ServiceGoogle` | `SETTING` | `LinkService!LinkService` | tidak |
| `convertJsonNusareToProduction.xml` | `ASM-FW-GISFW-WORK` | `convertJsonNusareToProduction` | `SETTING` | `LinkService!LinkService` | tidak |
| `getPremiumPaidOn.xml` | `ASM-FW-GISFW-WORK-NB` | `getPremiumPaidOn` | `SETTING` | `LinkService!LinkService` | tidak |

