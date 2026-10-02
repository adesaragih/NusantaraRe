# Inventaris Rule — Komite Claim FacIn

STEP D1, batch 2 (domain facultative inward). Dibangkitkan 2026-09-12 dari korpus READ-ONLY
`D:\XML\RNM_BRD\Komite Claim FacIn\`.

**Batas step ini:** hanya struktur + metadata. Tidak ada kesimpulan logika bisnis di sini.
Penelusuran perilaku adalah STEP D2 (`../flows/`).

**Identitas rule** diambil dari `<pxInsName>` (lihat `../README.md` §6.4).  
Kolom **Asal salinan** berisi `<pzOriginalInstanceKey>` dan hanya diisi bila **berbeda**
dari identitas — artinya rule ini hasil *Save As* dari rule lain. Tanda `—` berarti sama.
Class `—` berarti `<pxInsName>` tidak berisi bagian class (berlaku pada `DataPage`, yang
beridentitas global, bukan per class).

Perintah audit (dijalankan dari `D:\XML\RNM_BRD\`):

```
find "Komite Claim FacIn" -type f -name "*.xml" | wc -l
find "Komite Claim FacIn" -type f -name "*.xml" | awk -F/ '{print $2}' | sort | uniq -c
grep -o "<pxInsName>[^<]*" "Komite Claim FacIn/<tipe>/<file>.xml" | head -1
grep -o "<pzOriginalInstanceKey>[^<]*" "Komite Claim FacIn/<tipe>/<file>.xml" | head -1
```

## 1. Ringkasan terukur

- Total file XML: **114** — cocok dengan `../README.md` §5.1
- **Identitas rule unik: 114** dari 114 file → tidak ada identitas ganda di dalam modul ini
- Class Pega unik: **22**
- RDBList menurut jenis SQL: PLSQL=9, QUERY=8 (`PLSQL` = blok anonim `BEGIN…END;`, umumnya memanggil stored procedure)
- Prefix pada kunci RDBList: `ASM`=2, `RNM`=13, `GCNM`=2 (arti prefix **belum terverifikasi** — OQ-008)

| Tipe rule | Jumlah |
| --- | ---: |
| Activity | 28 |
| ConnectREST | 3 |
| DataTransform | 1 |
| DecisionTable | 1 |
| Flow | 1 |
| FlowAction | 5 |
| RDBList | 17 |
| ReportDefinition | 3 |
| Section | 5 |
| SystemSettings | 1 |
| When | 49 |
| **TOTAL** | **114** |

## 2. Sebaran class Pega

Class yang **tidak** berawalan `ASM-` adalah class bawaan Pega atau class portal — ruang
lingkup migrasinya belum ditentukan (**OQ-009**).

| Class | Jumlah rule | Aplikasi? |
| --- | ---: | --- |
| `@BASECLASS` | 26 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GCNMFW-WORK-KOMITE` | 15 | aplikasi |
| `ASM-FW-GISFW-DATA` | 13 | aplikasi |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | 10 | aplikasi |
| `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | 9 | aplikasi |
| `ASM-FW-GCNMFW-WORK` | 8 | aplikasi |
| `ASM-FW-GISFW-WORK` | 6 | aplikasi |
| `ASM-FW-GCNMFW-DATA-OBJECTITEM` | 4 | aplikasi |
| `ASM-FW-GISFW-DATA-SPREADINGRISK` | 4 | aplikasi |
| `ASM-FW-GISFW-DATA-FACOFFER` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-POLICYJSON` | 3 | aplikasi |
| `ASM-FW-GCNMFW-INT-V_POLIS` | 2 | aplikasi |
| `ASM-FW-GCNMFW-WORK-PNC` | 2 | aplikasi |
| `ASM-FW-GCNMFW-INT-CLAIMREJECTED` | 1 | aplikasi |
| `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | 1 | aplikasi |
| `ASM-FW-GCNMFW-INT-OS_AKSEPTASI_KLAIM` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-LST_BANK_GROUP` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_LINK_SERVICE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-REINSURANCETYPE` | 1 | aplikasi |
| `DATA-ADMIN-OPERATOR-ID` | 1 | **bukan aplikasi** → OQ-009 |
| `DATA-PARTY-PERSON` | 1 | **bukan aplikasi** → OQ-009 |
| `LINKSERVICE` | 1 | **bukan aplikasi** → OQ-009 |

Rule di class bukan-aplikasi: **29** dari 114.

## 3. Daftar rule per tipe

### Activity — 28 rule (`RULE-OBJ-ACTIVITY`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `ApprovalKomite_Act.xml` | `ASM-FW-GCNMFW-WORK-KOMITE` | `APPROVALKOMITE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | — |
| `CallSpreadingView.xml` | `ASM-FW-GISFW-DATA-SPREADINGRISK` | `CALLSPREADINGVIEW` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `GetBase64Attachment.xml` | `ASM-FW-GCNMFW-WORK` | `GETBASE64ATTACHMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | `WORK- / LOADATTACHMENTDATA` |
| `GetLinkService.xml` | `ASM-FW-GISFW-INT-M_LINK_SERVICE` | `GETLINKSERVICE` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetUrlGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETURLGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / INSERTGOOGLESTORAGE_ACT` |
| `HitServiceToKasirKMT_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `HITSERVICETOKASIRKMT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=42 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / HITSERVICETOKASIR_ACT` |
| `InsertDocument_Act.xml` | `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | `INSERTDOCUMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `InsertGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERTGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `InsertJsonClaimNonMBU_act.xml` | `ASM-FW-GCNMFW-WORK-KOMITE` | `INSERTJSONCLAIMNONMBU_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `ASM-FW-GCNMFW-WORK / INSERTJSONCLAIMNONMBU_ACT` |
| `InsertLogServiceClaim.xml` | `ASM-FW-GCNMFW-WORK` | `INSERTLOGSERVICECLAIM` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `KomitePostAct.xml` | `ASM-FW-GCNMFW-WORK-KOMITE` | `KOMITEPOSTACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `KomitePost_Adjustment.xml` | `ASM-FW-GCNMFW-WORK-KOMITE` | `KOMITEPOST_ADJUSTMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=63 | — |
| `KomitePost_CloseClaim.xml` | `ASM-FW-GCNMFW-WORK-KOMITE` | `KOMITEPOST_CLOSECLAIM` | [terverifikasi] pyActivityType=ACTIVITY; steps=33 | — |
| `KomitePost_Reject.xml` | `ASM-FW-GCNMFW-WORK-KOMITE` | `KOMITEPOST_REJECT` | [terverifikasi] pyActivityType=ACTIVITY; steps=39 | — |
| `KomiteRouter.xml` | `ASM-FW-GCNMFW-WORK-KOMITE` | `KOMITEROUTER` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `KonversiKlaim_Act.xml` | `ASM-FW-GCNMFW-WORK` | `KONVERSIKLAIM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `PreSecurityReas_Act.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `PRESECURITYREAS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `PreShowRetro_Act.xml` | `ASM-FW-GISFW-DATA-SPREADINGRISK` | `PRESHOWRETRO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `PrintPDFAccep_MultiAksep_KMT.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `PRINTPDFACCEP_MULTIAKSEP_KMT` | [terverifikasi] pyActivityType=ACTIVITY; steps=47 | `ASM-FW-GCNMFW-WORK-PNC / PRINTPDFACCEP_MULTIAKSEP` |
| `SaveAccept_ACT.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `SAVEACCEPT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=51 | — |
| `SaveAcceptation_KMT.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SAVEACCEPTATION_KMT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / SAVEACCEPTATION` |
| `SaveReject_ACT_KMT.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `SAVEREJECT_ACT_KMT` | [terverifikasi] pyActivityType=ACTIVITY; steps=20 | `ASM-FW-GCNMFW-DATA-OBJECTITEM / SAVEREJECT_ACT` |
| `SendEmailKlaim_KMT.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SENDEMAILKLAIM_KMT` | [terverifikasi] pyActivityType=ACTIVITY; steps=21 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / SENDEMAILKLAIM` |
| `SendErrorDirectKasir.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SENDERRORDIRECTKASIR` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SetDataForInformation_Act.xml` | `ASM-FW-GCNMFW-WORK-KOMITE` | `SETDATAFORINFORMATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `SetProteksiSubmiteKomite.xml` | `ASM-FW-GCNMFW-WORK-KOMITE` | `SETPROTEKSISUBMITEKOMITE` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GCNMFW-WORK-KOMITE / SETVALUEKOMITE` |
| `SetValueKomite.xml` | `ASM-FW-GCNMFW-WORK-KOMITE` | `SETVALUEKOMITE` | [terverifikasi] pyActivityType=ACTIVITY; steps=24 | — |
| `getStatusKonversi_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GETSTATUSKONVERSI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |

### ConnectREST — 3 rule (`RULE-CONNECT-REST`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `KonversiKlaimNonLife.xml` | `ASM-FW-GCNMFW-WORK` | `KONVERSIKLAIMNONLIFE` | [terverifikasi] pyServiceName=KonversiKlaimNonLife; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | `ASM-FW-GCNMFW-WORK-PNC / INSERTREJECTCLAIM` |
| `SendAcceptationToKasir.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SENDACCEPTATIONTOKASIR` | [terverifikasi] pyServiceName=SendAcceptationToKasir; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | — |
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `SERVICEGOOGLE` | [terverifikasi] pyServiceName=ServiceGoogle; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / GOOGLESTORAGE_UPLOAD` |

### DataTransform — 1 rule (`RULE-OBJ-MODEL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `ChronologyInsertion_DT.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `CHRONOLOGYINSERTION_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### DecisionTable — 1 rule (`RULE-DECLARE-DECISIONTABLE`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `GetMimeType.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETMIMETYPE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Flow — 1 rule (`RULE-OBJ-FLOW`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `Komite_Flow.xml` | `ASM-FW-GCNMFW-WORK-KOMITE` | `KOMITE_FLOW` | [terverifikasi] pyStartActivity=Start1 | — |

### FlowAction — 5 rule (`RULE-OBJ-FLOWACTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `DetailAdjustmentFac.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `DETAILADJUSTMENTFAC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowRetro.xml` | `ASM-FW-GISFW-DATA-SPREADINGRISK` | `SHOWRETRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowSecurityReinsurer.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `SHOWSECURITYREINSURER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SpreadingDetail.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `SPREADINGDETAIL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewTransferDtl.xml` | `ASM-FW-GCNMFW-WORK-KOMITE` | `VIEWTRANSFERDTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### RDBList — 17 rule (`RULE-CONNECT-SQL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `GenerateImageID_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GENERATEIMAGEID_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GetAppName_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETAPPNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.T_FOLDER_IMAGE | — |
| `GetEmailCeding_SQL.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GETEMAILCEDING_SQL` | [terverifikasi] sqlKind=QUERY; procs=GL.F_GET_EMAIL; sqlOps=SELECT | — |
| `GetKodeProdNonLife_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETKODEPRODNONLIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.KODE_PRODUKSI | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETSEQUENCENUMBER_SQL` |
| `GetLinkStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETLINKSTORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GETTOKENSTORAGE_SQL` |
| `GetSequenceNumber_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETSEQUENCENUMBER_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETTANGGALCLOSING_SQL` |
| `GetTokenStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETTOKENSTORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.GET_TOKEN_STORAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |
| `InsertClaimPNC.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `INSERTCLAIMPNC` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_KLAIM_PNC | — |
| `InsertClaimRejected_Sql.xml` | `ASM-FW-GCNMFW-INT-CLAIMREJECTED` | `INSERTCLAIMREJECTED_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.CLAIMREJECTED | — |
| `InsertHistoryAkseptasiPega_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTHISTORYAKSEPTASIPEGA_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=HISTORYAKSEPTASIPEGA | — |
| `InsertLOGDirectKasir_SQL.xml` | `ASM-FW-GCNMFW-WORK` | `INSERTLOGDIRECTKASIR_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.DIRECTTOKASIR_LOG | `ASM-FW-GISFW-WORK / RNM!INSERTLOGMOP_SQL` |
| `InsertLogServiceClaim.xml` | `ASM-FW-GCNMFW-WORK` | `INSERTLOGSERVICECLAIM` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.MONITORING_KLAIM_LOG | — |
| `Insert_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERT_T_STORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GENERATEIMAGEID_SQL` |
| `SaveOSClaim_SQL.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `SAVEOSCLAIM_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM | — |
| `UpdateSubProgresKlaim.xml` | `ASM-FW-GCNMFW-WORK-KOMITE` | `UPDATESUBPROGRESKLAIM` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=POOLDATA.SUBPROGRESSCLAIM | — |
| `Update_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `UPDATE_T_STORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |
| `getStatusKonversi_SQL.xml` | `ASM-FW-GCNMFW-INT-OS_AKSEPTASI_KLAIM` | `GETSTATUSKONVERSI_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=REINSURANCE.TRLOSS_DETAIL_T | — |

### ReportDefinition — 3 rule (`RULE-OBJ-REPORT-DEFINITION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseBankGroup.xml` | `ASM-FW-GISFW-INT-LST_BANK_GROUP` | `BROWSEBANKGROUP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseReinsuranceType_RD.xml` | `ASM-FW-GISFW-INT-REINSURANCETYPE` | `BROWSEREINSURANCETYPE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GetEmailUser_RD.xml` | `DATA-ADMIN-OPERATOR-ID` | `GETEMAILUSER_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-ADMIN-OPERATOR-ID / OPERATORSBYSKILL` |

### Section — 5 rule (`RULE-HTML-SECTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `DetailAdjustmentFac.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `DETAILADJUSTMENTFAC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowRetro_Sec.xml` | `ASM-FW-GISFW-DATA-SPREADINGRISK` | `SHOWRETRO_SEC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowSecurityReinsurer.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `SHOWSECURITYREINSURER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowTransfer.xml` | `ASM-FW-GCNMFW-WORK-KOMITE` | `SHOWTRANSFER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SpreadingDetail.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `SPREADINGDETAIL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### SystemSettings — 1 rule (`RULE-ADMIN-SYSTEM-SETTINGS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `LinkService.xml` | `LINKSERVICE` | `LINKSERVICE` | [terverifikasi] pyPurpose=LinkService | — |

### When — 49 rule (`RULE-OBJ-WHEN`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `IsAllRisk.xml` | `ASM-FW-GISFW-DATA` | `ISALLRISK` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsAneka.xml` | `ASM-FW-GISFW-DATA` | `ISANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsAviationHull.xml` | `ASM-FW-GISFW-DATA` | `ISAVIATIONHULL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsBillboardNeon.xml` | `ASM-FW-GISFW-WORK` | `ISBILLBOARDNEON` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISBILLBOARDNEON` |
| `IsBoiler.xml` | `@BASECLASS` | `ISBOILER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsBondingAndCustomBonds.xml` | `ASM-FW-GISFW-DATA` | `ISBONDINGANDCUSTOMBONDS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISBONDING` |
| `IsBondingKBG.xml` | `@BASECLASS` | `ISBONDINGKBG` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsBurglary.xml` | `ASM-FW-GISFW-DATA` | `ISBURGLARY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCAR.xml` | `@BASECLASS` | `ISCAR` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCIS.xml` | `@BASECLASS` | `ISCIS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCIT.xml` | `@BASECLASS` | `ISCIT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCLM.xml` | `@BASECLASS` | `ISCLM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCLMNP.xml` | `@BASECLASS` | `ISCLMNP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / ISCLMP` |
| `IsCLMP.xml` | `@BASECLASS` | `ISCLMP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / ISCLM` |
| `IsContractorsPlantMachinery.xml` | `@BASECLASS` | `ISCONTRACTORSPLANTMACHINERY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCrime.xml` | `ASM-FW-GISFW-DATA` | `ISCRIME` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCustomBond.xml` | `ASM-FW-GCNMFW-WORK` | `ISCUSTOMBOND` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISCUSTOMBOND` |
| `IsCustomBonds.xml` | `ASM-FW-GCNMFW-WORK` | `ISCUSTOMBONDS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsEar.xml` | `@BASECLASS` | `ISEAR` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsEnvironmental.xml` | `ASM-FW-GISFW-DATA` | `ISENVIRONMENTAL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISCRIME` |
| `IsExclusion.xml` | `@BASECLASS` | `ISEXCLUSION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsFidelity.xml` | `@BASECLASS` | `ISFIDELITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK / ISFIDELITY` |
| `IsFire.xml` | `ASM-FW-GISFW-WORK` | `ISFIRE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsFireStyle1.xml` | `@BASECLASS` | `ISFIRESTYLE1` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISFIRESTYLE1` |
| `IsFireStyle2.xml` | `@BASECLASS` | `ISFIRESTYLE2` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISFIRESTYLE2` |
| `IsGlass.xml` | `@BASECLASS` | `ISGLASS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsGrowingTrees.xml` | `@BASECLASS` | `ISGROWINGTREES` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsHE.xml` | `ASM-FW-GISFW-DATA` | `ISHE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsKPR.xml` | `@BASECLASS` | `ISKPR` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISKPR` |
| `IsKomiteLoop.xml` | `ASM-FW-GCNMFW-WORK-KOMITE` | `ISKOMITELOOP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsLandRig.xml` | `@BASECLASS` | `ISLANDRIG` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsLiability.xml` | `@BASECLASS` | `ISLIABILITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsMBD.xml` | `ASM-FW-GISFW-DATA` | `ISMBD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsMBU.xml` | `ASM-FW-GISFW-DATA` | `ISMBU` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsMarineCargo.xml` | `ASM-FW-GISFW-WORK` | `ISMARINECARGO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsMarineHull.xml` | `ASM-FW-GISFW-DATA` | `ISMARINEHULL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsOilGas.xml` | `@BASECLASS` | `ISOILGAS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISOILGAS` |
| `IsPA.xml` | `DATA-PARTY-PERSON` | `ISPA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / ISTRAVEL` |
| `IsPEGAPROD.xml` | `@BASECLASS` | `ISPEGAPROD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `WORK- / ISPEGAPROD` |
| `IsPEGASyariah.xml` | `@BASECLASS` | `ISPEGASYARIAH` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsProductsLiability.xml` | `ASM-FW-GISFW-WORK` | `ISPRODUCTSLIABILITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISPRODUCTSLIABILITY` |
| `IsProfessionalLiability.xml` | `ASM-FW-GISFW-WORK` | `ISPROFESSIONALLIABILITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISPROFESSIONALLIABILITY` |
| `IsTravel.xml` | `@BASECLASS` | `ISTRAVEL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISTRAVEL` |
| `IsWorkmenCompensation.xml` | `ASM-FW-GISFW-WORK` | `ISWORKMENCOMPENSATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISCIT` |
| `IsYieldShortfall.xml` | `ASM-FW-GISFW-DATA` | `ISYIELDSHORTFALL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISGROWINGTREES` |
| `isBillboardNeonSyariah.xml` | `@BASECLASS` | `ISBILLBOARDNEONSYARIAH` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isElectronicEquipment.xml` | `@BASECLASS` | `ISELECTRONICEQUIPMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isGolfInsurance.xml` | `ASM-FW-GISFW-DATA` | `ISGOLFINSURANCE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isMaintenance.xml` | `@BASECLASS` | `ISMAINTENANCE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

## 4. Tabrakan nama di dalam modul ini

### 4.1 Nama file dipakai di lebih dari satu tipe rule

| Nama file | Tipe | Class |
| --- | --- | --- |
| `SpreadingDetail.xml` | FlowAction, Section | `ASM-FW-GCNMFW-DATA-OBJECTITEM, ASM-FW-GCNMFW-DATA-OBJECTITEM` |
| `InsertLogServiceClaim.xml` | Activity, RDBList | `ASM-FW-GCNMFW-WORK, ASM-FW-GCNMFW-WORK` |
| `DetailAdjustmentFac.xml` | FlowAction, Section | `ASM-FW-GCNMFW-DATA-ADJUSTMENT, ASM-FW-GCNMFW-DATA-ADJUSTMENT` |
| `ShowSecurityReinsurer.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-FACOFFER, ASM-FW-GISFW-DATA-FACOFFER` |

Total: **4** nama file.

### 4.2 Rule `When` bernama sama di lebih dari satu class

_Tidak ada._

## 5. Keluarga klon (satu asal salinan dipakai >1 rule)

Dibaca dari `<pzOriginalInstanceKey>` yang sama pada beberapa file. Ini **sinyal kandidat
konsolidasi** saat migrasi, bukan bukti bahwa perilakunya identik — kesamaan perilaku baru
dapat dinyatakan setelah isinya dibandingkan di STEP D2.

| Asal salinan (class / nama) | Jumlah | File |
| --- | ---: | --- |
| `@BASECLASS / ISCLM` | 2 | `When/IsCLM.xml`, `When/IsCLMP.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` | 2 | `RDBList/GetTokenStorage_SQL.xml`, `RDBList/Update_T_Storage_SQL.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / INSERTGOOGLESTORAGE_ACT` | 2 | `Activity/GetUrlGoogleStorage_Act.xml`, `Activity/InsertGoogleStorage_Act.xml` |
| `ASM-FW-GCNMFW-DATA-OBJECTITEM / SPREADINGDETAIL` | 2 | `FlowAction/SpreadingDetail.xml`, `Section/SpreadingDetail.xml` |
| `ASM-FW-GISFW-DATA-FACOFFER / SHOWSECURITYREINSURER` | 2 | `FlowAction/ShowSecurityReinsurer.xml`, `Section/ShowSecurityReinsurer.xml` |
| `ASM-FW-GCNMFW-DATA-ADJUSTMENT / DETAILADJUSTMENTFAC` | 2 | `FlowAction/DetailAdjustmentFac.xml`, `Section/DetailAdjustmentFac.xml` |
| `ASM-FW-GCNMFW-WORK-KOMITE / SETVALUEKOMITE` | 2 | `Activity/SetProteksiSubmiteKomite.xml`, `Activity/SetValueKomite.xml` |
| `ASM-FW-GISFW-DATA / ISCRIME` | 2 | `When/IsCrime.xml`, `When/IsEnvironmental.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GENERATEIMAGEID_SQL` | 2 | `RDBList/GenerateImageID_SQL.xml`, `RDBList/Insert_T_Storage_SQL.xml` |

## 6. Objek database yang disentuh

Diambil dari `<pyBrowseSQL>` pada rule `RDBList`. **Nama objek apa adanya dari SQL** —
kualifikasi skema tidak konsisten di korpus. **Tipe kolom, constraint, dan relasi tidak
diketahui** (OQ-001).

| Tabel/view (apa adanya di SQL) | Dirujuk oleh |
| --- | ---: |
| `T_STORAGE_IMAGE` | 3 rule |
| `HISTORYAKSEPTASIPEGA` | 1 rule |
| `POOLDATA.CLAIMREJECTED` | 1 rule |
| `POOLDATA.DIRECTTOKASIR_LOG` | 1 rule |
| `POOLDATA.KODE_PRODUKSI` | 1 rule |
| `POOLDATA.MONITORING_KLAIM_LOG` | 1 rule |
| `POOLDATA.SUBPROGRESSCLAIM` | 1 rule |
| `POOLDATA.T_FOLDER_IMAGE` | 1 rule |
| `REINSURANCE.TRLOSS_DETAIL_T` | 1 rule |

### Stored procedure / function yang dipanggil

**Body prosedur ini TIDAK ADA di korpus** — logika di dalamnya tidak dapat direkonstruksi
dari XML. Lihat OQ-002.

| Procedure | Dipanggil oleh |
| --- | --- |
| `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM` | `SaveOSClaim_SQL.xml` |
| `POOLDATA.PEGA_JSON_KLAIM_PNC` | `InsertClaimPNC.xml` |
| `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` | `GetSequenceNumber_SQL.xml` |
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

