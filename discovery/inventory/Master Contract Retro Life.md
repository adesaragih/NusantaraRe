# Inventaris Rule — Master Contract Retro Life

STEP D1, batch 4 (domain life, master & outward). Dibangkitkan 2026-09-12 dari korpus READ-ONLY
`D:\XML\RNM_BRD\Master Contract Retro Life\`.

**Batas step ini:** hanya struktur + metadata. Tidak ada kesimpulan logika bisnis di sini.
Penelusuran perilaku adalah STEP D2 (`../flows/`).

**Identitas rule** diambil dari `<pxInsName>` (lihat `../README.md` §6.4).  
Kolom **Asal salinan** berisi `<pzOriginalInstanceKey>` dan hanya diisi bila **berbeda**
dari identitas — artinya rule ini hasil *Save As* dari rule lain. Tanda `—` berarti sama.
Class `—` berarti `<pxInsName>` tidak berisi bagian class (berlaku pada `DataPage`, yang
beridentitas global, bukan per class).

Perintah audit (dijalankan dari `D:\XML\RNM_BRD\`):

```
find "Master Contract Retro Life" -type f -name "*.xml" | wc -l
find "Master Contract Retro Life" -type f -name "*.xml" | awk -F/ '{print $2}' | sort | uniq -c
grep -o "<pxInsName>[^<]*" "Master Contract Retro Life/<tipe>/<file>.xml" | head -1
grep -o "<pzOriginalInstanceKey>[^<]*" "Master Contract Retro Life/<tipe>/<file>.xml" | head -1
```

## 1. Ringkasan terukur

- Total file XML: **66** — cocok dengan `../README.md` §5.1
- **Identitas rule unik: 66** dari 66 file → tidak ada identitas ganda di dalam modul ini
- Class Pega unik: **14**
- RDBList menurut jenis SQL: PLSQL=6, QUERY=6 (`PLSQL` = blok anonim `BEGIN…END;`, umumnya memanggil stored procedure)
- Prefix pada kunci RDBList: `ASM`=12 (arti prefix **belum terverifikasi** — OQ-008)

| Tipe rule | Jumlah |
| --- | ---: |
| Activity | 27 |
| DataTransform | 1 |
| FlowAction | 2 |
| Harness | 4 |
| RDBList | 12 |
| ReportDefinition | 12 |
| Section | 8 |
| **TOTAL** | **66** |

## 2. Sebaran class Pega

Class yang **tidak** berawalan `ASM-` adalah class bawaan Pega atau class portal — ruang
lingkup migrasinya belum ditentukan (**OQ-009**).

| Class | Jumlah rule | Aplikasi? |
| --- | ---: | --- |
| `@BASECLASS` | 38 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` | 6 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE` | 5 | aplikasi |
| `DATA-PORTAL` | 4 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-INT-RETROCESSIONLIFE` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYCONTRACT_LIFE` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-AGENT` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-BUSINESS` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_RATE_LIFE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-RATE_LIFE_SUMMARY` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-REINSURANCETYPE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYSECURITYREINSURER_LIFE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYYEAR` | 1 | aplikasi |

Rule di class bukan-aplikasi: **42** dari 66.

## 3. Daftar rule per tipe

### Activity — 27 rule (`RULE-OBJ-ACTIVITY`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `CancelActivityTreatyContract.xml` | `@BASECLASS` | `CANCELACTIVITYTREATYCONTRACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CountingPercentShare_Act.xml` | `@BASECLASS` | `COUNTINGPERCENTSHARE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `DeleteRowBusiness.xml` | `@BASECLASS` | `DELETEROWBUSINESS` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `DeleteSecurityLife_Act.xml` | `@BASECLASS` | `DELETESECURITYLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `@BASECLASS / DELETERETROCESSIONLIFE_ACT` |
| `DeleteSecurityReinsurerLife_Act.xml` | `@BASECLASS` | `DELETESECURITYREINSURERLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `@BASECLASS / DELETESECURITYLIFE_ACT` |
| `DeleteTreatyLimit_Act.xml` | `@BASECLASS` | `DELETETREATYLIMIT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / DELETERETROCESSIONLIFE_ACT` |
| `NewInputBusinessLife_Act.xml` | `@BASECLASS` | `NEWINPUTBUSINESSLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `NewInputSecurityLife_Act.xml` | `@BASECLASS` | `NEWINPUTSECURITYLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `@BASECLASS / NEWINPUTBUSINESSLIFE_ACT` |
| `NewInputTreatyLimit_Life.xml` | `@BASECLASS` | `NEWINPUTTREATYLIMIT_LIFE` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `NewInputTreatyYear_Life_Act.xml` | `@BASECLASS` | `NEWINPUTTREATYYEAR_LIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `NewTreatyReinsurerDetail_Act.xml` | `@BASECLASS` | `NEWTREATYREINSURERDETAIL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SaveBusinessLife_Act.xml` | `@BASECLASS` | `SAVEBUSINESSLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `SaveBusinessToAllLife_Act.xml` | `@BASECLASS` | `SAVEBUSINESSTOALLLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `@BASECLASS / SAVEBUSINESSLIFE_ACT` |
| `SaveSecurityLife_Act.xml` | `@BASECLASS` | `SAVESECURITYLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `@BASECLASS / SAVEBUSINESSLIFE_ACT` |
| `SaveSecurityReinsurerLife_Act.xml` | `@BASECLASS` | `SAVESECURITYREINSURERLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `@BASECLASS / SAVESECURITYLIFE_ACT` |
| `SaveTreatyLimit_Act.xml` | `@BASECLASS` | `SAVETREATYLIMIT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `SaveTreatyYearLife_Act.xml` | `@BASECLASS` | `SAVETREATYYEARLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SetBusinessListLife_Act.xml` | `@BASECLASS` | `SETBUSINESSLISTLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `@BASECLASS / SETRETROLISTLIFE_ACT` |
| `SetErrorMessageReinsurer.xml` | `@BASECLASS` | `SETERRORMESSAGEREINSURER` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetParamRate.xml` | `@BASECLASS` | `SETPARAMRATE` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetParamRateTable.xml` | `@BASECLASS` | `SETPARAMRATETABLE` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / SETPARAMRATE` |
| `SetRetroListLife_Act.xml` | `@BASECLASS` | `SETRETROLISTLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetSecurityLife_Act.xml` | `@BASECLASS` | `SETSECURITYLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `@BASECLASS / SETRETROLISTLIFE_ACT` |
| `SetSecurityReinsurerLife_Act.xml` | `@BASECLASS` | `SETSECURITYREINSURERLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `@BASECLASS / SETSECURITYLIFE_ACT` |
| `SetTreatyYearLife_Act.xml` | `@BASECLASS` | `SETTREATYYEARLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetValueRetroLimit_TreatyYearLife.xml` | `@BASECLASS` | `SETVALUERETROLIMIT_TREATYYEARLIFE` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `TreatyLimit_TypeProtect.xml` | `@BASECLASS` | `TREATYLIMIT_TYPEPROTECT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |

### DataTransform — 1 rule (`RULE-OBJ-MODEL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `SetOutputParam_DT.xml` | `@BASECLASS` | `SETOUTPUTPARAM_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### FlowAction — 2 rule (`RULE-OBJ-FLOWACTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `ViewRate.xml` | `@BASECLASS` | `VIEWRATE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewRateTable.xml` | `@BASECLASS` | `VIEWRATETABLE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / VIEWRATE` |

### Harness — 4 rule (`RULE-HTML-HARNESS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `InboxBusinessLifeReinsurers.xml` | `DATA-PORTAL` | `INBOXBUSINESSLIFEREINSURERS` | [terverifikasi] pyInclude=1 | `DATA-PORTAL / INBOXRETROLIFEREINSURERSBUSINESS` |
| `InboxRetroLifeReinsurersList.xml` | `DATA-PORTAL` | `INBOXRETROLIFEREINSURERSLIST` | [terverifikasi] pyInclude=1 | `DATA-PORTAL / INBOXRETROLIFEREINSURERS` |
| `InboxRetroLimitReinsurers.xml` | `DATA-PORTAL` | `INBOXRETROLIMITREINSURERS` | [terverifikasi] pyInclude=1 | `DATA-PORTAL / INBOXRETROLIFEREINSURERS` |
| `InboxSecurityReinsurerLife.xml` | `DATA-PORTAL` | `INBOXSECURITYREINSURERLIFE` | [terverifikasi] pyInclude=1 | `DATA-PORTAL / INBOXRETROLIFEREINSURERSLIST` |

### RDBList — 12 rule (`RULE-CONNECT-SQL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `DeleteRowBusinessList.xml` | `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE` | `DELETEROWBUSINESSLIST` | [terverifikasi] sqlKind=QUERY; sqlOps=DELETE; tables=POOLDATA.TREATYBUSINESS_LIFE | `ASM-FW-GISFW-INT-TREATYBUSINESS / ASM!DELETEROWBUSINESSLIST` |
| `DeleteSecurityReinsurerLife_SQL.xml` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` | `DELETESECURITYREINSURERLIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=DELETE; tables=POOLDATA.TREATYSECURITYREINSURER_LIFE | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE / ASM!DELETESECURITYREINSURER_SQL` |
| `DeleteSecurityReinsurer_SQL.xml` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` | `DELETESECURITYREINSURER_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=DELETE; tables=POOLDATA.TREATYREINSURER_LIFE | `ASM-FW-GISFW-INT-TREATYYEAR / ASM!DELETE2MASTERCOPYDATA_SQL` |
| `DeleteTreatyLimit_SQL.xml` | `ASM-FW-GISFW-INT-RETROCESSIONLIFE` | `DELETETREATYLIMIT_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=DELETE; tables=POOLDATA.TREATYCONTRACT_LIFE | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE / ASM!DELETETREATYLIMIT_SQL` |
| `GetMasterReinsurerLifeList_SQl.xml` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` | `GETMASTERREINSURERLIFELIST_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TREATYREINSURER_LIFE | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE / LIFE!GETMASTERREINSURERLIFELIST_SQL` |
| `GetTreatyContract_life.xml` | `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE` | `GETTREATYCONTRACT_LIFE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TREATYCONTRACT_LIFE | `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE / ASM!SAVEMASTERTREATYBUSINESS_LIFE_SQL` |
| `SaveMasterTreatyBusiness_Life_SQL.xml` | `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE` | `SAVEMASTERTREATYBUSINESS_LIFE_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.INSERTBUSINESS_LIFE; sqlOps=UPDATE | — |
| `SaveMasterTreatyContract_Life_SQL.xml` | `ASM-FW-GISFW-INT-TREATYCONTRACT_LIFE` | `SAVEMASTERTREATYCONTRACT_LIFE_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.INSERTTREATYCONTRACT_LIFE; sqlOps=UPDATE | — |
| `SaveMasterTreatyReinsurer_Life_SQL.xml` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` | `SAVEMASTERTREATYREINSURER_LIFE_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.INSERTREINSURER_LIFE; sqlOps=UPDATE | `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE / ASM!SAVEMASTERTREATYREINSURER_LIFE_SQL` |
| `SaveMasterTreatySecurityReinsurer_Life_SQL.xml` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` | `SAVEMASTERTREATYSECURITYREINSURER_LIFE_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.INSERTSECURITYREINSURER_LIFE; sqlOps=UPDATE | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE / ASM!SAVEMASTERTREATYREINSURER_LIFE_SQL` |
| `SaveMasterTreatyYear_Life_SQL.xml` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` | `SAVEMASTERTREATYYEAR_LIFE_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.INSERTTREATYYEAR_LIFE; sqlOps=UPDATE | — |
| `SaveTreatyBusinessAll_Life_SQL.xml` | `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE` | `SAVETREATYBUSINESSALL_LIFE_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT,UPDATE; tables=POOLDATA.TREATYBUSINESS_LIFE | `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE / ASM!SAVEMASTERTREATYBUSINESS_LIFE_SQL` |

### ReportDefinition — 12 rule (`RULE-OBJ-REPORT-DEFINITION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseBusinessLife_RD.xml` | `ASM-FW-GISFW-INT-BUSINESS` | `BROWSEBUSINESSLIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-BUSINESS / BROWSEBUSINESS_RD` |
| `BrowseCedingCoLife_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSECEDINGCOLIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-AGENT / BROWSECEDINGCO_RD` |
| `BrowseDetailTreatyReisurerLife_RD.xml` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` | `BROWSEDETAILTREATYREISURERLIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseRateLifeSummary.xml` | `ASM-FW-GISFW-INT-RATE_LIFE_SUMMARY` | `BROWSERATELIFESUMMARY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseRateLife_RD.xml` | `ASM-FW-GISFW-INT-M_RATE_LIFE` | `BROWSERATELIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseReinsuranceTypeLimit_RD.xml` | `ASM-FW-GISFW-INT-REINSURANCETYPE` | `BROWSEREINSURANCETYPELIMIT_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseRetrocessionLife_RD.xml` | `ASM-FW-GISFW-INT-RETROCESSIONLIFE` | `BROWSERETROCESSIONLIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseSecurityReinsurer_Life_RD.xml` | `ASM-FW-GISFW-INT-TREATYSECURITYREINSURER_LIFE` | `BROWSESECURITYREINSURER_LIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyBusiness_Life_RD.xml` | `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE` | `BROWSETREATYBUSINESS_LIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyContract_Life_RD.xml` | `ASM-FW-GISFW-INT-TREATYCONTRACT_LIFE` | `BROWSETREATYCONTRACT_LIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyYear_Life_RD.xml` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` | `BROWSETREATYYEAR_LIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyYear_RD.xml` | `ASM-FW-GISFW-INT-TREATYYEAR` | `BROWSETREATYYEAR_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Section — 8 rule (`RULE-HTML-SECTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `GridRetrocessionLife.xml` | `@BASECLASS` | `GRIDRETROCESSIONLIFE` | [terverifikasi] pyInclude=2 | `@BASECLASS / GRIDTREATYCONTRACT` |
| `InputBusinessLifeReinsurers.xml` | `@BASECLASS` | `INPUTBUSINESSLIFEREINSURERS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlRetrocessionLife.xml` | `@BASECLASS` | `INPUTDTLRETROCESSIONLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputRetroLimitReinsurers.xml` | `@BASECLASS` | `INPUTRETROLIMITREINSURERS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputRetrocessionLife.xml` | `@BASECLASS` | `INPUTRETROCESSIONLIFE` | [terverifikasi] pyInclude=1 | — |
| `InputSecurityLifeReinsurers.xml` | `@BASECLASS` | `INPUTSECURITYLIFEREINSURERS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputSecurityReinsurerLife.xml` | `@BASECLASS` | `INPUTSECURITYREINSURERLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / INPUTSECURITYLIFEREINSURERS` |
| `ViewRate.xml` | `@BASECLASS` | `VIEWRATE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PLAN / VIEWRATE` |

## 4. Tabrakan nama di dalam modul ini

### 4.1 Nama file dipakai di lebih dari satu tipe rule

| Nama file | Tipe | Class |
| --- | --- | --- |
| `ViewRate.xml` | FlowAction, Section | `@BASECLASS, @BASECLASS` |

Total: **1** nama file.

### 4.2 Rule `When` bernama sama di lebih dari satu class

_Tidak ada._

## 5. Keluarga klon (satu asal salinan dipakai >1 rule)

Dibaca dari `<pzOriginalInstanceKey>` yang sama pada beberapa file. Ini **sinyal kandidat
konsolidasi** saat migrasi, bukan bukti bahwa perilakunya identik — kesamaan perilaku baru
dapat dinyatakan setelah isinya dibandingkan di STEP D2.

| Asal salinan (class / nama) | Jumlah | File |
| --- | ---: | --- |
| `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE / ASM!SAVEMASTERTREATYBUSINESS_LIFE_SQL` | 3 | `RDBList/GetTreatyContract_life.xml`, `RDBList/SaveMasterTreatyBusiness_Life_SQL.xml`, `RDBList/SaveTreatyBusinessAll_Life_SQL.xml` |
| `@BASECLASS / SETRETROLISTLIFE_ACT` | 3 | `Activity/SetBusinessListLife_Act.xml`, `Activity/SetRetroListLife_Act.xml`, `Activity/SetSecurityLife_Act.xml` |
| `@BASECLASS / VIEWRATE` | 2 | `FlowAction/ViewRate.xml`, `FlowAction/ViewRateTable.xml` |
| `@BASECLASS / DELETERETROCESSIONLIFE_ACT` | 2 | `Activity/DeleteSecurityLife_Act.xml`, `Activity/DeleteTreatyLimit_Act.xml` |
| `@BASECLASS / SAVEBUSINESSLIFE_ACT` | 3 | `Activity/SaveBusinessLife_Act.xml`, `Activity/SaveBusinessToAllLife_Act.xml`, `Activity/SaveSecurityLife_Act.xml` |
| `@BASECLASS / NEWINPUTBUSINESSLIFE_ACT` | 2 | `Activity/NewInputBusinessLife_Act.xml`, `Activity/NewInputSecurityLife_Act.xml` |
| `@BASECLASS / INPUTSECURITYLIFEREINSURERS` | 2 | `Section/InputSecurityLifeReinsurers.xml`, `Section/InputSecurityReinsurerLife.xml` |
| `@BASECLASS / SETPARAMRATE` | 2 | `Activity/SetParamRate.xml`, `Activity/SetParamRateTable.xml` |
| `DATA-PORTAL / INBOXRETROLIFEREINSURERS` | 2 | `Harness/InboxRetroLifeReinsurersList.xml`, `Harness/InboxRetroLimitReinsurers.xml` |

## 6. Objek database yang disentuh

Diambil dari `<pyBrowseSQL>` pada rule `RDBList`. **Nama objek apa adanya dari SQL** —
kualifikasi skema tidak konsisten di korpus. **Tipe kolom, constraint, dan relasi tidak
diketahui** (OQ-001).

| Tabel/view (apa adanya di SQL) | Dirujuk oleh |
| --- | ---: |
| `POOLDATA.TREATYBUSINESS_LIFE` | 2 rule |
| `POOLDATA.TREATYCONTRACT_LIFE` | 2 rule |
| `POOLDATA.TREATYREINSURER_LIFE` | 2 rule |
| `POOLDATA.TREATYSECURITYREINSURER_LIFE` | 1 rule |

### Stored procedure / function yang dipanggil

**Body prosedur ini TIDAK ADA di korpus** — logika di dalamnya tidak dapat direkonstruksi
dari XML. Lihat OQ-002.

| Procedure | Dipanggil oleh |
| --- | --- |
| `POOLDATA.INSERTTREATYYEAR_LIFE` | `SaveMasterTreatyYear_Life_SQL.xml` |
| `POOLDATA.INSERTTREATYCONTRACT_LIFE` | `SaveMasterTreatyContract_Life_SQL.xml` |
| `POOLDATA.INSERTREINSURER_LIFE` | `SaveMasterTreatyReinsurer_Life_SQL.xml` |
| `POOLDATA.INSERTSECURITYREINSURER_LIFE` | `SaveMasterTreatySecurityReinsurer_Life_SQL.xml` |
| `POOLDATA.INSERTBUSINESS_LIFE` | `SaveMasterTreatyBusiness_Life_SQL.xml` |

## 7. Integrasi eksternal (`ConnectREST`)

_Tidak ada rule ConnectREST di modul ini._

