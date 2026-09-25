# Inventaris Rule — Komite Claim Prop

STEP D1, batch 3 (domain claim non-facultative). Dibangkitkan 2026-09-12 dari korpus READ-ONLY
`D:\XML\RNM_BRD\Komite Claim Prop\`.

**Batas step ini:** hanya struktur + metadata. Tidak ada kesimpulan logika bisnis di sini.
Penelusuran perilaku adalah STEP D2 (`../flows/`).

**Identitas rule** diambil dari `<pxInsName>` (lihat `../README.md` §6.4).  
Kolom **Asal salinan** berisi `<pzOriginalInstanceKey>` dan hanya diisi bila **berbeda**
dari identitas — artinya rule ini hasil *Save As* dari rule lain. Tanda `—` berarti sama.
Class `—` berarti `<pxInsName>` tidak berisi bagian class (berlaku pada `DataPage`, yang
beridentitas global, bukan per class).

Perintah audit (dijalankan dari `D:\XML\RNM_BRD\`):

```
find "Komite Claim Prop" -type f -name "*.xml" | wc -l
find "Komite Claim Prop" -type f -name "*.xml" | awk -F/ '{print $2}' | sort | uniq -c
grep -o "<pxInsName>[^<]*" "Komite Claim Prop/<tipe>/<file>.xml" | head -1
grep -o "<pzOriginalInstanceKey>[^<]*" "Komite Claim Prop/<tipe>/<file>.xml" | head -1
```

## 1. Ringkasan terukur

- Total file XML: **80** — cocok dengan `../README.md` §5.1
- **Identitas rule unik: 80** dari 80 file → tidak ada identitas ganda di dalam modul ini
- Class Pega unik: **19**
- RDBList menurut jenis SQL: PLSQL=10, QUERY=13 (`PLSQL` = blok anonim `BEGIN…END;`, umumnya memanggil stored procedure)
- Prefix pada kunci RDBList: `ASM`=3, `RNM`=12, `GCNM`=8 (arti prefix **belum terverifikasi** — OQ-008)

| Tipe rule | Jumlah |
| --- | ---: |
| Activity | 35 |
| ConnectREST | 3 |
| DataTransform | 2 |
| DecisionTable | 1 |
| Flow | 1 |
| FlowAction | 2 |
| RDBList | 23 |
| ReportDefinition | 4 |
| Section | 2 |
| SystemSettings | 1 |
| When | 6 |
| **TOTAL** | **80** |

## 2. Sebaran class Pega

Class yang **tidak** berawalan `ASM-` adalah class bawaan Pega atau class portal — ruang
lingkup migrasinya belum ditentukan (**OQ-009**).

| Class | Jumlah rule | Aplikasi? |
| --- | ---: | --- |
| `ASM-FW-GCNMFW-WORK-KOMITETREATY` | 14 | aplikasi |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | 12 | aplikasi |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | 10 | aplikasi |
| `@BASECLASS` | 7 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | 7 | aplikasi |
| `ASM-FW-GCNMFW-WORK` | 7 | aplikasi |
| `ASM-FW-GISFW-INT-POLICYJSON` | 7 | aplikasi |
| `ASM-FW-GCNMFW-INT-V_POLIS` | 4 | aplikasi |
| `ASM-FW-GISFW-DATA-TREATYINTOTAL` | 2 | aplikasi |
| `ASM-FW-GCNMFW-INT-CLAIMREJECTED` | 1 | aplikasi |
| `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | 1 | aplikasi |
| `ASM-FW-GCNMFW-INT-OS_AKSEPTASI_KLAIM` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-CURRENCY` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-CURRENCYSTANDARD` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-LST_BANK_GROUP` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_LINK_SERVICE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-REINSURANCETYPE` | 1 | aplikasi |
| `DATA-ADMIN-OPERATOR-ID` | 1 | **bukan aplikasi** → OQ-009 |
| `LINKSERVICE` | 1 | **bukan aplikasi** → OQ-009 |

Rule di class bukan-aplikasi: **9** dari 80.

## 3. Daftar rule per tipe

### Activity — 35 rule (`RULE-OBJ-ACTIVITY`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AddEstimation_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `ADDESTIMATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | — |
| `AddLossAllocation_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `ADDLOSSALLOCATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | — |
| `CheckEstimateDate_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `CHECKESTIMATEDATE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `CountEstimation_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `COUNTESTIMATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=39 | — |
| `CountListClaimAmountIDR.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `COUNTLISTCLAIMAMOUNTIDR` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `CountPersen_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `COUNTPERSEN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `CountSpreading_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `COUNTSPREADING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=25 | — |
| `CurencyEstimation_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `CURENCYESTIMATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `GetBase64Attachment.xml` | `ASM-FW-GCNMFW-WORK` | `GETBASE64ATTACHMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | `WORK- / LOADATTACHMENTDATA` |
| `GetLinkService.xml` | `ASM-FW-GISFW-INT-M_LINK_SERVICE` | `GETLINKSERVICE` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetUrlGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETURLGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / INSERTGOOGLESTORAGE_ACT` |
| `HTMLToPDF.xml` | `@BASECLASS` | `HTMLTOPDF` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `HitServiceToKasirKMT_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `HITSERVICETOKASIRKMT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=42 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / HITSERVICETOKASIR_ACT` |
| `InsertDocument_Act.xml` | `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | `INSERTDOCUMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `InsertGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERTGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `InsertJsonClaimTreaty_act.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATY` | `INSERTJSONCLAIMTREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GCNMFW-WORK / INSERTJSONCLAIMTREATY_ACT` |
| `InsertLogServiceClaim.xml` | `ASM-FW-GCNMFW-WORK` | `INSERTLOGSERVICECLAIM` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `KomitePost.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATY` | `KOMITEPOST` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `KomitePostAdjustment.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATY` | `KOMITEPOSTADJUSTMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=49 | — |
| `KomitePost_Close.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATY` | `KOMITEPOST_CLOSE` | [terverifikasi] pyActivityType=ACTIVITY; steps=34 | — |
| `KomitePost_Reject.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATY` | `KOMITEPOST_REJECT` | [terverifikasi] pyActivityType=ACTIVITY; steps=32 | `ASM-FW-GCNMFW-WORK-KOMITE / KOMITEPOST_REJECT` |
| `KomiteRouter.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATY` | `KOMITEROUTER` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `KonversiKlaim_Act.xml` | `ASM-FW-GCNMFW-WORK` | `KONVERSIKLAIM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `PrintFileAcceptance_TKMT.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `PRINTFILEACCEPTANCE_TKMT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / PRINTFILEACCEPTANCE` |
| `SaveAcceptationTreaty_TKMT.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SAVEACCEPTATIONTREATY_TKMT` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / SAVEACCEPTATIONTREATY_ACT` |
| `SaveAcceptation_Act.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATY` | `SAVEACCEPTATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `SaveRejectTreatyIn_Act_KMT.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SAVEREJECTTREATYIN_ACT_KMT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SendEmailKlaim_KMT.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SENDEMAILKLAIM_KMT` | [terverifikasi] pyActivityType=ACTIVITY; steps=21 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / SENDEMAILKLAIM` |
| `SendEmailWithAttachments.xml` | `@BASECLASS` | `SENDEMAILWITHATTACHMENTS` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SendErrorDirectKasir.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SENDERRORDIRECTKASIR` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SetCurencyList_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SETCURENCYLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `SetDataAcceptationTreaty_Act.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATY` | `SETDATAACCEPTATIONTREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `SetDefNonCatastrope_Act.xml` | `ASM-FW-GCNMFW-WORK` | `SETDEFNONCATASTROPE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetKomiteList_Act.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATY` | `SETKOMITELIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `getStatusKonversi_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GETSTATUSKONVERSI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |

### ConnectREST — 3 rule (`RULE-CONNECT-REST`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `KonversiKlaimNonLife.xml` | `ASM-FW-GCNMFW-WORK` | `KONVERSIKLAIMNONLIFE` | [terverifikasi] pyServiceName=KonversiKlaimNonLife; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | `ASM-FW-GCNMFW-WORK-PNC / INSERTREJECTCLAIM` |
| `SendAcceptationToKasir.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SENDACCEPTATIONTOKASIR` | [terverifikasi] pyServiceName=SendAcceptationToKasir; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | — |
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `SERVICEGOOGLE` | [terverifikasi] pyServiceName=ServiceGoogle; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / GOOGLESTORAGE_UPLOAD` |

### DataTransform — 2 rule (`RULE-OBJ-MODEL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `CekInterestListDtl_DT.xml` | `ASM-FW-GISFW-DATA-TREATYINTOTAL` | `CEKINTERESTLISTDTL_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InsertChronology_DT.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `INSERTCHRONOLOGY_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-PNC / INSERTCHRONOLOGY_DT` |

### DecisionTable — 1 rule (`RULE-DECLARE-DECISIONTABLE`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `GetMimeType.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETMIMETYPE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Flow — 1 rule (`RULE-OBJ-FLOW`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `KomiteTreaty_Flow.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATY` | `KOMITETREATY_FLOW` | [terverifikasi] pyStartActivity=Start1 | — |

### FlowAction — 2 rule (`RULE-OBJ-FLOWACTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `ViewDetailInterest.xml` | `ASM-FW-GISFW-DATA-TREATYINTOTAL` | `VIEWDETAILINTEREST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewTransferDtl.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATY` | `VIEWTRANSFERDTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### RDBList — 23 rule (`RULE-CONNECT-SQL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `CurrencyStandard.xml` | `ASM-FW-GISFW-INT-CURRENCYSTANDARD` | `CURRENCYSTANDARD` | [terverifikasi] sqlKind=QUERY; procs=POOLDATA.GETCURRENCYSTANDARD; sqlOps=SELECT | `ASM-FW-GISFW-INT-CURRENCYSTANDARD / ASM!UPDATEMASTERCURRENCYSTANDARD` |
| `GenerateImageID_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GENERATEIMAGEID_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GetAppName_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETAPPNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.T_FOLDER_IMAGE | — |
| `GetEmailCeding_SQL.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GETEMAILCEDING_SQL` | [terverifikasi] sqlKind=QUERY; procs=GL.F_GET_EMAIL; sqlOps=SELECT | — |
| `GetKodeProdNonLife_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETKODEPRODNONLIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.KODE_PRODUKSI | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETSEQUENCENUMBER_SQL` |
| `GetLimitPLATreatyin.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLIMITPLATREATYIN` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=PROPORTIONALARRG | — |
| `GetLinkStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETLINKSTORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GETTOKENSTORAGE_SQL` |
| `GetListRetro_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLISTRETRO_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYREINSURER | — |
| `GetOldIDBusiness.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `GETOLDIDBUSINESS` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=BUSINESS | — |
| `GetSequenceNumber_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETSEQUENCENUMBER_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETTANGGALCLOSING_SQL` |
| `GetTokenStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETTOKENSTORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.GET_TOKEN_STORAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |
| `GetTreatyGroupID.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETTREATYGROUPID` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYBUSINESS | — |
| `InsertClaimPNC.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `INSERTCLAIMPNC` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_KLAIM_PNC | — |
| `InsertClaimRejected_Sql.xml` | `ASM-FW-GCNMFW-INT-CLAIMREJECTED` | `INSERTCLAIMREJECTED_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.CLAIMREJECTED | — |
| `InsertHistoryAkseptasiPega_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTHISTORYAKSEPTASIPEGA_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=HISTORYAKSEPTASIPEGA | — |
| `InsertLOGDirectKasir_SQL.xml` | `ASM-FW-GCNMFW-WORK` | `INSERTLOGDIRECTKASIR_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.DIRECTTOKASIR_LOG | `ASM-FW-GISFW-WORK / RNM!INSERTLOGMOP_SQL` |
| `InsertLogServiceClaim.xml` | `ASM-FW-GCNMFW-WORK` | `INSERTLOGSERVICECLAIM` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.MONITORING_KLAIM_LOG | — |
| `Insert_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERT_T_STORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GENERATEIMAGEID_SQL` |
| `SaveDataToOsAkseptasiNP.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `SAVEDATATOOSAKSEPTASINP` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP | `ASM-FW-GCNMFW-INT-V_POLIS / GCNM!SAVEDATATOOSAKSEPTASI` |
| `SaveOSClaim_SQL.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `SAVEOSCLAIM_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM | — |
| `TreatyYearTreatyin_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `TREATYYEARTREATYIN_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYYEAR | — |
| `Update_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `UPDATE_T_STORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |
| `getStatusKonversi_SQL.xml` | `ASM-FW-GCNMFW-INT-OS_AKSEPTASI_KLAIM` | `GETSTATUSKONVERSI_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=REINSURANCE.TRLOSS_DETAIL_T | — |

### ReportDefinition — 4 rule (`RULE-OBJ-REPORT-DEFINITION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseBankGroup.xml` | `ASM-FW-GISFW-INT-LST_BANK_GROUP` | `BROWSEBANKGROUP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCurrency_RD.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `BROWSECURRENCY_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseReinsuranceType_RD.xml` | `ASM-FW-GISFW-INT-REINSURANCETYPE` | `BROWSEREINSURANCETYPE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GetEmailUser_RD.xml` | `DATA-ADMIN-OPERATOR-ID` | `GETEMAILUSER_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-ADMIN-OPERATOR-ID / OPERATORSBYSKILL` |

### Section — 2 rule (`RULE-HTML-SECTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `ShowTransfer.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATY` | `SHOWTRANSFER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewDetailInterest.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATY` | `VIEWDETAILINTEREST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYINTOTAL / VIEWDETAILINTEREST` |

### SystemSettings — 1 rule (`RULE-ADMIN-SYSTEM-SETTINGS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `LinkService.xml` | `LINKSERVICE` | `LINKSERVICE` | [terverifikasi] pyPurpose=LinkService | — |

### When — 6 rule (`RULE-OBJ-WHEN`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `IsCLM.xml` | `@BASECLASS` | `ISCLM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCLMNP.xml` | `@BASECLASS` | `ISCLMNP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / ISCLMP` |
| `IsCLMP.xml` | `@BASECLASS` | `ISCLMP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / ISCLM` |
| `IsKomiteLoop.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATY` | `ISKOMITELOOP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsPEGAPROD.xml` | `@BASECLASS` | `ISPEGAPROD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `WORK- / ISPEGAPROD` |
| `IsPEGASyariah.xml` | `@BASECLASS` | `ISPEGASYARIAH` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

## 4. Tabrakan nama di dalam modul ini

### 4.1 Nama file dipakai di lebih dari satu tipe rule

| Nama file | Tipe | Class |
| --- | --- | --- |
| `InsertLogServiceClaim.xml` | Activity, RDBList | `ASM-FW-GCNMFW-WORK, ASM-FW-GCNMFW-WORK` |
| `ViewDetailInterest.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINTOTAL, ASM-FW-GCNMFW-WORK-KOMITETREATY` |

Total: **2** nama file.

### 4.2 Rule `When` bernama sama di lebih dari satu class

_Tidak ada._

## 5. Keluarga klon (satu asal salinan dipakai >1 rule)

Dibaca dari `<pzOriginalInstanceKey>` yang sama pada beberapa file. Ini **sinyal kandidat
konsolidasi** saat migrasi, bukan bukti bahwa perilakunya identik — kesamaan perilaku baru
dapat dinyatakan setelah isinya dibandingkan di STEP D2.

| Asal salinan (class / nama) | Jumlah | File |
| --- | ---: | --- |
| `@BASECLASS / ISCLM` | 2 | `When/IsCLM.xml`, `When/IsCLMP.xml` |
| `ASM-FW-GISFW-DATA-TREATYINTOTAL / VIEWDETAILINTEREST` | 2 | `FlowAction/ViewDetailInterest.xml`, `Section/ViewDetailInterest.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` | 2 | `RDBList/GetTokenStorage_SQL.xml`, `RDBList/Update_T_Storage_SQL.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / INSERTGOOGLESTORAGE_ACT` | 2 | `Activity/GetUrlGoogleStorage_Act.xml`, `Activity/InsertGoogleStorage_Act.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GENERATEIMAGEID_SQL` | 2 | `RDBList/GenerateImageID_SQL.xml`, `RDBList/Insert_T_Storage_SQL.xml` |

## 6. Objek database yang disentuh

Diambil dari `<pyBrowseSQL>` pada rule `RDBList`. **Nama objek apa adanya dari SQL** —
kualifikasi skema tidak konsisten di korpus. **Tipe kolom, constraint, dan relasi tidak
diketahui** (OQ-001).

| Tabel/view (apa adanya di SQL) | Dirujuk oleh |
| --- | ---: |
| `T_STORAGE_IMAGE` | 3 rule |
| `BUSINESS` | 1 rule |
| `HISTORYAKSEPTASIPEGA` | 1 rule |
| `POOLDATA.CLAIMREJECTED` | 1 rule |
| `POOLDATA.DIRECTTOKASIR_LOG` | 1 rule |
| `POOLDATA.KODE_PRODUKSI` | 1 rule |
| `POOLDATA.MONITORING_KLAIM_LOG` | 1 rule |
| `POOLDATA.T_FOLDER_IMAGE` | 1 rule |
| `PROPORTIONALARRG` | 1 rule |
| `REINSURANCE.TRLOSS_DETAIL_T` | 1 rule |
| `TREATYBUSINESS` | 1 rule |
| `TREATYREINSURER` | 1 rule |
| `TREATYYEAR` | 1 rule |

### Stored procedure / function yang dipanggil

**Body prosedur ini TIDAK ADA di korpus** — logika di dalamnya tidak dapat direkonstruksi
dari XML. Lihat OQ-002.

| Procedure | Dipanggil oleh |
| --- | --- |
| `POOLDATA.GETCURRENCYSTANDARD` | `CurrencyStandard.xml` |
| `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM` | `SaveOSClaim_SQL.xml` |
| `POOLDATA.PEGA_JSON_KLAIM_PNC` | `InsertClaimPNC.xml` |
| `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` | `GetSequenceNumber_SQL.xml` |
| `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP` | `SaveDataToOsAkseptasiNP.xml` |
| `POOLDATA.GET_TOKEN_STORAGE` | `GetTokenStorage_SQL.xml` |
| `GL.F_GET_EMAIL` | `GetEmailCeding_SQL.xml` |

## 7. Integrasi eksternal (`ConnectREST`)

Alamat tujuan dicatat sebagai **konfigurasi**, bukan URL literal. `baseURL=SETTING` berarti
alamat diambil dari rule `SystemSettings` yang disebut di kolom Setting; nilai sebenarnya
**tidak ada di korpus**. Kolom terakhir menandai rule yang justru memuat URL literal —
itu temuan, bukan fakta arsitektur (lihat `_summary.md`).

| File | Class | pyServiceName | Sumber base URL | Setting | URL literal di dalam rule? |
| --- | --- | --- | --- | --- | --- |
| `KonversiKlaimNonLife.xml` | `ASM-FW-GCNMFW-WORK` | `KonversiKlaimNonLife` | `SETTING` | `LinkService!LinkService` | tidak |
| `SendAcceptationToKasir.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SendAcceptationToKasir` | `SETTING` | `LinkService!LinkService` | tidak |
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `ServiceGoogle` | `SETTING` | `LinkService!LinkService` | tidak |

