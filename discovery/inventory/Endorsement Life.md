# Inventaris Rule — Endorsement Life

STEP D1, batch 4 (domain life, master & outward). Dibangkitkan 2026-09-12 dari korpus READ-ONLY
`D:\XML\RNM_BRD\Endorsement Life\`.

**Batas step ini:** hanya struktur + metadata. Tidak ada kesimpulan logika bisnis di sini.
Penelusuran perilaku adalah STEP D2 (`../flows/`).

**Identitas rule** diambil dari `<pxInsName>` (lihat `../README.md` §6.4).  
Kolom **Asal salinan** berisi `<pzOriginalInstanceKey>` dan hanya diisi bila **berbeda**
dari identitas — artinya rule ini hasil *Save As* dari rule lain. Tanda `—` berarti sama.
Class `—` berarti `<pxInsName>` tidak berisi bagian class (berlaku pada `DataPage`, yang
beridentitas global, bukan per class).

Perintah audit (dijalankan dari `D:\XML\RNM_BRD\`):

```
find "Endorsement Life" -type f -name "*.xml" | wc -l
find "Endorsement Life" -type f -name "*.xml" | awk -F/ '{print $2}' | sort | uniq -c
grep -o "<pxInsName>[^<]*" "Endorsement Life/<tipe>/<file>.xml" | head -1
grep -o "<pzOriginalInstanceKey>[^<]*" "Endorsement Life/<tipe>/<file>.xml" | head -1
```

## 1. Ringkasan terukur

- Total file XML: **75** — cocok dengan `../README.md` §5.1
- **Identitas rule unik: 75** dari 75 file → tidak ada identitas ganda di dalam modul ini
- Class Pega unik: **19**
- RDBList menurut jenis SQL: PLSQL=6, QUERY=11 (`PLSQL` = blok anonim `BEGIN…END;`, umumnya memanggil stored procedure)
- Prefix pada kunci RDBList: `ASM`=14, `RNM`=3 (arti prefix **belum terverifikasi** — OQ-008)

| Tipe rule | Jumlah |
| --- | ---: |
| Activity | 16 |
| ConnectREST | 1 |
| DataTransform | 2 |
| DecisionTable | 1 |
| Flow | 1 |
| FlowAction | 6 |
| Harness | 8 |
| RDBList | 17 |
| ReportDefinition | 6 |
| Section | 14 |
| SystemSettings | 1 |
| When | 2 |
| **TOTAL** | **75** |

## 2. Sebaran class Pega

Class yang **tidak** berawalan `ASM-` adalah class bawaan Pega atau class portal — ruang
lingkup migrasinya belum ditentukan (**OQ-009**).

| Class | Jumlah rule | Aplikasi? |
| --- | ---: | --- |
| `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | 38 | aplikasi |
| `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | 7 | aplikasi |
| `@BASECLASS` | 4 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE` | 4 | aplikasi |
| `DATA-PORTAL` | 4 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-POLICYJSON` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYCONTRACT_LIFE` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-AGENT` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-CLIENT` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-MARKETINGOFFICER` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_LINK_SERVICE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_RATE_LIFE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-OFFERJSON` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-RI_COMM_LIFE` | 1 | aplikasi |
| `ASM-FW-GISFW-WORK` | 1 | aplikasi |
| `ASM-FW-GISFW-WORK-LIFE` | 1 | aplikasi |
| `ASSIGN-WORKLIST` | 1 | **bukan aplikasi** → OQ-009 |
| `LINKSERVICE` | 1 | **bukan aplikasi** → OQ-009 |

Rule di class bukan-aplikasi: **10** dari 75.

## 3. Daftar rule per tipe

### Activity — 16 rule (`RULE-OBJ-ACTIVITY`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AddHistorySuggest.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `ADDHISTORYSUGGEST` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GISFW-WORK-LIFE / ADDHISTORYSUGGEST` |
| `Calculate1_Act.xml` | `ASM-FW-GISFW-WORK-LIFE` | `CALCULATE1_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=24 | — |
| `CancelCreateCaseEDML.xml` | `DATA-PORTAL` | `CANCELCREATECASEEDML` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CreateCaseEMDL.xml` | `DATA-PORTAL` | `CREATECASEEMDL` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | — |
| `GenerateDataDtlLife_act.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `GENERATEDATADTLLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `GenerateNoEDM_Life.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `GENERATENOEDM_LIFE` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / INSERTJSONPOLISLIFE_ACT` |
| `GetLinkService.xml` | `ASM-FW-GISFW-INT-M_LINK_SERVICE` | `GETLINKSERVICE` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetOldDetail_EDM.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `GETOLDDETAIL_EDM` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / SETPOLICYNO` |
| `InsertJsonPolisLife_Act.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `INSERTJSONPOLISLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=22 | `ASM-FW-GISFW-WORK-LIFE / INSERTJSONPOLISLIFE_ACT` |
| `MappingEDMLife.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `MAPPINGEDMLIFE` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / SETPOLICYNO` |
| `SaveCSVEDMLife.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `SAVECSVEDMLIFE` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / SETPREMI_EDM` |
| `SelectAllEdmLife_act.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `SELECTALLEDMLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / SELLECTALLEDMLIFE_ACT` |
| `SetErrorBatalEndorsement_Act.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `SETERRORBATALENDORSEMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | — |
| `SetPremi_EDM.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `SETPREMI_EDM` | [terverifikasi] pyActivityType=ACTIVITY; steps=16 | — |
| `UploadCSVEDMLifePremium_Act.xml` | `@BASECLASS` | `UPLOADCSVEDMLIFEPREMIUM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | `@BASECLASS / UPLOADCSVLIFEPREMIUM_ACT` |
| `serviceInsertArasapasLife_act.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `SERVICEINSERTARASAPASLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `ASM-FW-GISFW-WORK-LIFE / SERVICEINSERTARASAPASLIFE_ACT` |

### ConnectREST — 1 rule (`RULE-CONNECT-REST`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `ConvertJsonNusareToProduction.xml` | `ASM-FW-GISFW-WORK` | `CONVERTJSONNUSARETOPRODUCTION` | [terverifikasi] pyServiceName=convertJsonNusareToProduction; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | — |

### DataTransform — 2 rule (`RULE-OBJ-MODEL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AppendCurrencySummary_DT.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `APPENDCURRENCYSUMMARY_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-LIFE / APPENDCURRENCYSUMMARY_DT` |
| `SetCoBName_Act.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `SETCOBNAME_ACT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-LIFE / SETCOBNAME_ACT` |

### DecisionTable — 1 rule (`RULE-DECLARE-DECISIONTABLE`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `IsLifeAccepted.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `ISLIFEACCEPTED` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Flow — 1 rule (`RULE-OBJ-FLOW`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `InputEDMLife.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `INPUTEDMLIFE` | [terverifikasi] pyStartActivity=Start1 | — |

### FlowAction — 6 rule (`RULE-OBJ-FLOWACTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `InputEDMLife.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `INPUTEDMLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputEDMLife_Summary.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `INPUTEDMLIFE_SUMMARY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PL_DetailAction.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `PL_DETAILACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `RetroLife.xml` | `ASM-FW-GISFW-INT-TREATYCONTRACT_LIFE` | `RETROLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetJsonPolisEDMLife_Confirm.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `SETJSONPOLISEDMLIFE_CONFIRM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / INPUTEDMLIFE_SUMMARY` |
| `UploadCSV_LifeEndorsement.xml` | `@BASECLASS` | `UPLOADCSV_LIFEENDORSEMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / UPLOADCSV_LIFEPREMIUM` |

### Harness — 8 rule (`RULE-HTML-HARNESS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `EditDetail_Harness.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `EDITDETAIL_HARNESS` | [terverifikasi] pyInclude=1 | — |
| `EndorsmentLife_harnes.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `ENDORSMENTLIFE_HARNES` | [terverifikasi] pyInclude=1 | — |
| `InboxEndorsementLife.xml` | `DATA-PORTAL` | `INBOXENDORSEMENTLIFE` | [terverifikasi] pyInclude=1 | — |
| `ViewCSVResult_LifeEDM.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `VIEWCSVRESULT_LIFEEDM` | [terverifikasi] pyInclude=1 | — |
| `ViewOldPolicy_EDM.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `VIEWOLDPOLICY_EDM` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / EDITDETAIL_HARNESS` |
| `ViewOldPolicy_EDM_QP.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `VIEWOLDPOLICY_EDM_QP` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / VIEWOLDPOLICY_EDM` |
| `ViewOldPolicy_EDM_TP.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `VIEWOLDPOLICY_EDM_TP` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / VIEWOLDPOLICY_EDM_QP` |
| `ViewOldPolicy_EDM_TR.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `VIEWOLDPOLICY_EDM_TR` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / VIEWOLDPOLICY_EDM_TP` |

### RDBList — 17 rule (`RULE-CONNECT-SQL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `CekDoubleInsured.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `CEKDOUBLEINSURED` | [terverifikasi] sqlKind=PLSQL; sqlOps=SELECT; tables=V_COUNT,POOLDATA.M_TEMPUPLOADLIFE,V_MIN,V_TEMPRETENSICEDING | — |
| `DeleteTempUploadDataLife.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `DELETETEMPUPLOADDATALIFE` | [terverifikasi] sqlKind=QUERY; sqlOps=DELETE; tables=M_TEMPUPLOADLIFE | — |
| `Generate_NoEndorsmentLife.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `GENERATE_NOENDORSMENTLIFE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.JSON_POLIS | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / ASM!GENERATE_NOPREMIUMLIST` |
| `GetEdmTypeLife.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `GETEDMTYPELIFE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.JSON_POLIS | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / RNM!GETPL_NUMBERLIFE` |
| `GetPL_NumberLife.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `GETPL_NUMBERLIFE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.JSON_POLIS | `ASM-FW-GISFW-WORK-LIFE / RNM!GETPLANDNOPOLIS_SQL` |
| `GetPolicyNoByCaseId.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETPOLICYNOBYCASEID` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | `ASM-FW-GISFW-INT-POLICYJSON / ASM!BROWSECLIENTIDBYNOPOLICY` |
| `GetProdKeOldData_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `GETPRODKEOLDDATA_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.JSON_POLIS | `ASM-FW-GISFW-INT-OFFERJSON / ASM!GETEDMOLDDATA_SQL` |
| `GetProdkeNopolis.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `GETPRODKENOPOLIS` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.JSON_POLIS | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / RNM!GETPL_NUMBERLIFE` |
| `GetProductDtlPL.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `GETPRODUCTDTLPL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=PRODUCT_LIFE | — |
| `GetRateLifePM.xml` | `ASM-FW-GISFW-INT-M_RATE_LIFE` | `GETRATELIFEPM` | [terverifikasi] sqlKind=PLSQL; sqlOps=SELECT; tables=VCOUNT,RATE_LIFE | — |
| `InsertDataUploadLife.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `INSERTDATAUPLOADLIFE` | [terverifikasi] sqlKind=QUERY; sqlOps=INSERT; tables=POOLDATA.M_TEMPUPLOADLIFE | — |
| `InsertJsonPolisEDM.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` | `INSERTJSONPOLISEDM` | [terverifikasi] sqlKind=PLSQL; procs=DBMS_LOB.CREATETEMPORARY; sqlOps=INSERT; tables=POOLDATA.JSON_POLIS | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY / ASM!INSERTJSONPOLIS` |
| `InsertPLSummary.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` | `INSERTPLSUMMARY` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY | — |
| `SaveLifeinProduction_SQL.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `SAVELIFEINPRODUCTION_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.LIFEINPRODUCTION | `ASM-FW-GISFW-WORK-LIFE / ASM!SAVELIFEINPRODUCTION_SQL` |
| `SaveMasterLPDet.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `SAVEMASTERLPDET` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.M_LIFE_PREMIUM_DETAIL | — |
| `SearcStatusBayarArasaps_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `SEARCSTATUSBAYARARASAPS_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=ARASAPAS.DETAIL_INVOICE | — |
| `SelectComm_SQL.xml` | `ASM-FW-GISFW-INT-RI_COMM_LIFE` | `SELECTCOMM_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=RICOMM_LIFE | `ASM-FW-GISFW-INT-RI_COMM_LIFE / ASM!CHECKDOUBLERATELIFE_SQL` |

### ReportDefinition — 6 rule (`RULE-OBJ-REPORT-DEFINITION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseCedingCoLife_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSECEDINGCOLIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-AGENT / BROWSECEDINGCO_RD` |
| `BrowseClientNusaRe_RD.xml` | `ASM-FW-GISFW-INT-CLIENT` | `BROWSECLIENTNUSARE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseMarketingOfficer_RD.xml` | `ASM-FW-GISFW-INT-MARKETINGOFFICER` | `BROWSEMARKETINGOFFICER_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowsePremiumList_RD.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` | `BROWSEPREMIUMLIST_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `FilterProteksiEDMLife.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `FILTERPROTEKSIEDMLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / INBOXEDMLIFE` |
| `InboxEDMLife.xml` | `ASSIGN-WORKLIST` | `INBOXEDMLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Section — 14 rule (`RULE-HTML-SECTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `ConfirmSection.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `CONFIRMSECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ConfirmSubmitEDM.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `CONFIRMSUBMITEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / CONFIRMCLAIMAMOUNT` |
| `EditDetail_Section.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `EDITDETAIL_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `EndorsmentLife_Section.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `ENDORSMENTLIFE_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InboxEndorsementLife.xml` | `DATA-PORTAL` | `INBOXENDORSEMENTLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputEDMLife.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `INPUTEDMLIFE` | [terverifikasi] pyInclude=1 | — |
| `PL_Detail_Sec.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `PL_DETAIL_SEC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `RetroDetailLife.xml` | `ASM-FW-GISFW-INT-TREATYCONTRACT_LIFE` | `RETRODETAILLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowLifePremiumSummary_EDM.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `SHOWLIFEPREMIUMSUMMARY_EDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-LIFE / SHOWLIFEPREMIUMSUMMARY` |
| `ViewCSVResult_LifeEDM.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `VIEWCSVRESULT_LIFEEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / VIEWCSVRESULT_LIFEPREMIUM` |
| `ViewOldPolicy_EDM.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `VIEWOLDPOLICY_EDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewOldPolicy_EDM_QP.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `VIEWOLDPOLICY_EDM_QP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / VIEWOLDPOLICY_EDM` |
| `ViewOldPolicy_EDM_TP.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `VIEWOLDPOLICY_EDM_TP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / VIEWOLDPOLICY_EDM_TR` |
| `ViewOldPolicy_EDM_TR.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` | `VIEWOLDPOLICY_EDM_TR` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / VIEWOLDPOLICY_EDM_QP` |

### SystemSettings — 1 rule (`RULE-ADMIN-SYSTEM-SETTINGS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `LinkService.xml` | `LINKSERVICE` | `LINKSERVICE` | [terverifikasi] pyPurpose=LinkService | — |

### When — 2 rule (`RULE-OBJ-WHEN`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `IsPEGAPROD.xml` | `@BASECLASS` | `ISPEGAPROD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `WORK- / ISPEGAPROD` |
| `recordEvent.xml` | `@BASECLASS` | `RECORDEVENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

## 4. Tabrakan nama di dalam modul ini

### 4.1 Nama file dipakai di lebih dari satu tipe rule

| Nama file | Tipe | Class |
| --- | --- | --- |
| `ViewOldPolicy_EDM_TR.xml` | Harness, Section | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE, ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` |
| `ViewOldPolicy_EDM_TP.xml` | Harness, Section | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE, ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` |
| `InputEDMLife.xml` | Flow, FlowAction, Section | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE, ASM-FW-GISFW-WORK-ENDORSEMENTLIFE, ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` |
| `InboxEndorsementLife.xml` | Harness, Section | `DATA-PORTAL, DATA-PORTAL` |
| `ViewOldPolicy_EDM_QP.xml` | Harness, Section | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE, ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` |
| `ViewCSVResult_LifeEDM.xml` | Harness, Section | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL, ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` |
| `ViewOldPolicy_EDM.xml` | Harness, Section | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE, ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` |

Total: **7** nama file.

### 4.2 Rule `When` bernama sama di lebih dari satu class

_Tidak ada._

## 5. Keluarga klon (satu asal salinan dipakai >1 rule)

Dibaca dari `<pzOriginalInstanceKey>` yang sama pada beberapa file. Ini **sinyal kandidat
konsolidasi** saat migrasi, bukan bukti bahwa perilakunya identik — kesamaan perilaku baru
dapat dinyatakan setelah isinya dibandingkan di STEP D2.

| Asal salinan (class / nama) | Jumlah | File |
| --- | ---: | --- |
| `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / EDITDETAIL_HARNESS` | 2 | `Harness/EditDetail_Harness.xml`, `Harness/ViewOldPolicy_EDM.xml` |
| `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / SETPREMI_EDM` | 2 | `Activity/SaveCSVEDMLife.xml`, `Activity/SetPremi_EDM.xml` |
| `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / INPUTEDMLIFE` | 3 | `Flow/InputEDMLife.xml`, `FlowAction/InputEDMLife.xml`, `Section/InputEDMLife.xml` |
| `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / INPUTEDMLIFE_SUMMARY` | 2 | `FlowAction/InputEDMLife_Summary.xml`, `FlowAction/SetJsonPolisEDMLife_Confirm.xml` |
| `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / RNM!GETPL_NUMBERLIFE` | 2 | `RDBList/GetEdmTypeLife.xml`, `RDBList/GetProdkeNopolis.xml` |
| `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / VIEWOLDPOLICY_EDM_QP` | 2 | `Harness/ViewOldPolicy_EDM_TP.xml`, `Section/ViewOldPolicy_EDM_TR.xml` |
| `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / SETPOLICYNO` | 2 | `Activity/GetOldDetail_EDM.xml`, `Activity/MappingEDMLife.xml` |
| `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE / VIEWOLDPOLICY_EDM` | 3 | `Harness/ViewOldPolicy_EDM_QP.xml`, `Section/ViewOldPolicy_EDM.xml`, `Section/ViewOldPolicy_EDM_QP.xml` |
| `DATA-PORTAL / INBOXENDORSEMENTLIFE` | 2 | `Harness/InboxEndorsementLife.xml`, `Section/InboxEndorsementLife.xml` |

## 6. Objek database yang disentuh

Diambil dari `<pyBrowseSQL>` pada rule `RDBList`. **Nama objek apa adanya dari SQL** —
kualifikasi skema tidak konsisten di korpus. **Tipe kolom, constraint, dan relasi tidak
diketahui** (OQ-001).

| Tabel/view (apa adanya di SQL) | Dirujuk oleh |
| --- | ---: |
| `POOLDATA.JSON_POLIS` | 6 rule |
| `POOLDATA.M_TEMPUPLOADLIFE` | 2 rule |
| `ARASAPAS.DETAIL_INVOICE` | 1 rule |
| `JSON_POLIS` | 1 rule |
| `M_TEMPUPLOADLIFE` | 1 rule |
| `POOLDATA.LIFEINPRODUCTION` | 1 rule |
| `POOLDATA.M_LIFE_PREMIUM_DETAIL` | 1 rule |
| `PRODUCT_LIFE` | 1 rule |
| `RATE_LIFE` | 1 rule |
| `RICOMM_LIFE` | 1 rule |
| `VCOUNT` | 1 rule |
| `V_COUNT` | 1 rule |
| `V_MIN` | 1 rule |
| `V_TEMPRETENSICEDING` | 1 rule |

### Stored procedure / function yang dipanggil

**Body prosedur ini TIDAK ADA di korpus** — logika di dalamnya tidak dapat direkonstruksi
dari XML. Lihat OQ-002.

| Procedure | Dipanggil oleh |
| --- | --- |
| `DBMS_LOB.CREATETEMPORARY` | `InsertJsonPolisEDM.xml` |
| `POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY` | `InsertPLSummary.xml` |

## 7. Integrasi eksternal (`ConnectREST`)

Alamat tujuan dicatat sebagai **konfigurasi**, bukan URL literal. `baseURL=SETTING` berarti
alamat diambil dari rule `SystemSettings` yang disebut di kolom Setting; nilai sebenarnya
**tidak ada di korpus**. Kolom terakhir menandai rule yang justru memuat URL literal —
itu temuan, bukan fakta arsitektur (lihat `_summary.md`).

| File | Class | pyServiceName | Sumber base URL | Setting | URL literal di dalam rule? |
| --- | --- | --- | --- | --- | --- |
| `ConvertJsonNusareToProduction.xml` | `ASM-FW-GISFW-WORK` | `convertJsonNusareToProduction` | `SETTING` | `LinkService!LinkService` | tidak |

