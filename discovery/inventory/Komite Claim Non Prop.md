# Inventaris Rule — Komite Claim Non Prop

STEP D1, batch 3 (domain claim non-facultative). Dibangkitkan 2026-09-12 dari korpus READ-ONLY
`D:\XML\RNM_BRD\Komite Claim Non Prop\`.

**Batas step ini:** hanya struktur + metadata. Tidak ada kesimpulan logika bisnis di sini.
Penelusuran perilaku adalah STEP D2 (`../flows/`).

**Identitas rule** diambil dari `<pxInsName>` (lihat `../README.md` §6.4).  
Kolom **Asal salinan** berisi `<pzOriginalInstanceKey>` dan hanya diisi bila **berbeda**
dari identitas — artinya rule ini hasil *Save As* dari rule lain. Tanda `—` berarti sama.
Class `—` berarti `<pxInsName>` tidak berisi bagian class (berlaku pada `DataPage`, yang
beridentitas global, bukan per class).

Perintah audit (dijalankan dari `D:\XML\RNM_BRD\`):

```
find "Komite Claim Non Prop" -type f -name "*.xml" | wc -l
find "Komite Claim Non Prop" -type f -name "*.xml" | awk -F/ '{print $2}' | sort | uniq -c
grep -o "<pxInsName>[^<]*" "Komite Claim Non Prop/<tipe>/<file>.xml" | head -1
grep -o "<pzOriginalInstanceKey>[^<]*" "Komite Claim Non Prop/<tipe>/<file>.xml" | head -1
```

## 1. Ringkasan terukur

- Total file XML: **59** — cocok dengan `../README.md` §5.1
- **Identitas rule unik: 59** dari 59 file → tidak ada identitas ganda di dalam modul ini
- Class Pega unik: **16**
- RDBList menurut jenis SQL: PLSQL=10, QUERY=6 (`PLSQL` = blok anonim `BEGIN…END;`, umumnya memanggil stored procedure)
- Prefix pada kunci RDBList: `ASM`=5, `RNM`=9, `GCNM`=2 (arti prefix **belum terverifikasi** — OQ-008)

| Tipe rule | Jumlah |
| --- | ---: |
| Activity | 20 |
| ConnectREST | 5 |
| DataTransform | 1 |
| DecisionTable | 1 |
| Flow | 1 |
| FlowAction | 2 |
| RDBList | 16 |
| ReportDefinition | 3 |
| Section | 2 |
| SystemSettings | 1 |
| When | 6 |
| excludeXML | 1 |
| **TOTAL** | **59** |

## 2. Sebaran class Pega

Class yang **tidak** berawalan `ASM-` adalah class bawaan Pega atau class portal — ruang
lingkup migrasinya belum ditentukan (**OQ-009**).

| Class | Jumlah rule | Aplikasi? |
| --- | ---: | --- |
| `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` | 12 | aplikasi |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | 7 | aplikasi |
| `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | 6 | aplikasi |
| `ASM-FW-GCNMFW-INT-V_POLIS` | 6 | aplikasi |
| `ASM-FW-GCNMFW-WORK` | 6 | aplikasi |
| `@BASECLASS` | 5 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | 5 | aplikasi |
| `ASM-FW-GISFW-INT-POLICYJSON` | 3 | aplikasi |
| `ASM-FW-GISFW-DATA-SPREADINGRISK` | 2 | aplikasi |
| `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | 1 | aplikasi |
| `ASM-FW-GCNMFW-INT-OS_AKSEPTASI_KLAIM` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-CURRENCY` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-LST_BANK_GROUP` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_LINK_SERVICE` | 1 | aplikasi |
| `DATA-ADMIN-OPERATOR-ID` | 1 | **bukan aplikasi** → OQ-009 |
| `LINKSERVICE` | 1 | **bukan aplikasi** → OQ-009 |

Rule di class bukan-aplikasi: **7** dari 59.

## 3. Daftar rule per tipe

### Activity — 20 rule (`RULE-OBJ-ACTIVITY`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `GenerateAccCNP_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `GENERATEACCCNP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=25 | — |
| `GetLinkService.xml` | `ASM-FW-GISFW-INT-M_LINK_SERVICE` | `GETLINKSERVICE` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `HitServiceToKasirKMT_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `HITSERVICETOKASIRKMT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=42 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / HITSERVICETOKASIR_ACT` |
| `InsertDocument_Act.xml` | `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | `INSERTDOCUMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `InsertGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERTGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `InsertJsonClaimTreatyNonProp_act.xml` | `ASM-FW-GCNMFW-WORK` | `INSERTJSONCLAIMTREATYNONPROP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `ASM-FW-GCNMFW-WORK / INSERTJSONCLAIMTREATY_ACT` |
| `InsertOSKlaimCNP.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` | `INSERTOSKLAIMCNP` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `InsertOSSubjectivityCNP.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` | `INSERTOSSUBJECTIVITYCNP` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP / INSERTOSKLAIMCNP` |
| `InsertXOLKlaimCNP.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` | `INSERTXOLKLAIMCNP` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP / INSERTOSKLAIMCNP` |
| `KomitePostAdjustment.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` | `KOMITEPOSTADJUSTMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=53 | — |
| `KomitePostAdjustmentCWP.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` | `KOMITEPOSTADJUSTMENTCWP` | [terverifikasi] pyActivityType=ACTIVITY; steps=35 | `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP / KOMITEPOSTADJUSTMENT` |
| `KomiteRouter.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` | `KOMITEROUTER` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `KonversiKlaim_Act.xml` | `ASM-FW-GCNMFW-WORK` | `KONVERSIKLAIM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `ResetSubjectivityNote.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` | `RESETSUBJECTIVITYNOTE` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SaveRejectOSKomiteCNP.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `SAVEREJECTOSKOMITECNP` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP / SAVEREJECTOSCNP` |
| `SendEmailKlaimRejectClose_KMT.xml` | `ASM-FW-GCNMFW-WORK` | `SENDEMAILKLAIMREJECTCLOSE_KMT` | [terverifikasi] pyActivityType=ACTIVITY; steps=16 | `ASM-FW-GCNMFW-WORK / SENDEMAILKLAIMREJECTCLOSE` |
| `SendEmailKlaim_KMT.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SENDEMAILKLAIM_KMT` | [terverifikasi] pyActivityType=ACTIVITY; steps=21 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / SENDEMAILKLAIM` |
| `SendErrorDirectKasir.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SENDERRORDIRECTKASIR` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SetKomiteList_Act.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` | `SETKOMITELIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `ASM-FW-GCNMFW-WORK-KOMITETREATY / SETKOMITELIST_ACT` |
| `getStatusKonversi_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GETSTATUSKONVERSI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |

### ConnectREST — 5 rule (`RULE-CONNECT-REST`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `KonversiKlaimNonLife.xml` | `ASM-FW-GCNMFW-WORK` | `KONVERSIKLAIMNONLIFE` | [terverifikasi] pyServiceName=KonversiKlaimNonLife; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | `ASM-FW-GCNMFW-WORK-PNC / INSERTREJECTCLAIM` |
| `SendAcceptationToKasir.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SENDACCEPTATIONTOKASIR` | [terverifikasi] pyServiceName=SendAcceptationToKasir; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | — |
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `SERVICEGOOGLE` | [terverifikasi] pyServiceName=ServiceGoogle; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / GOOGLESTORAGE_UPLOAD` |
| `insertClaimFinalOrClosed_NP.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `INSERTCLAIMFINALORCLOSED_NP` | [terverifikasi] pyServiceName=insertClaimFinalOrClosed_NP; baseURL=URL; resolusi=JNDIName; URL_LITERAL_DITEMUKAN=2 | — |
| `insertClaimReject_NP.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `INSERTCLAIMREJECT_NP` | [terverifikasi] pyServiceName=insertClaimReject_NP; baseURL=URL; resolusi=JNDIName; URL_LITERAL_DITEMUKAN=2 | — |

### DataTransform — 1 rule (`RULE-OBJ-MODEL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `InsertChronology_DT.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `INSERTCHRONOLOGY_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### DecisionTable — 1 rule (`RULE-DECLARE-DECISIONTABLE`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `GetMimeType.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETMIMETYPE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Flow — 1 rule (`RULE-OBJ-FLOW`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `KomiteTreaty_Flow.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` | `KOMITETREATY_FLOW` | [terverifikasi] pyStartActivity=Start1 | `ASM-FW-GCNMFW-WORK-KOMITETREATY / KOMITETREATY_FLOW` |

### FlowAction — 2 rule (`RULE-OBJ-FLOWACTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `ShowDetailXOL.xml` | `ASM-FW-GISFW-DATA-SPREADINGRISK` | `SHOWDETAILXOL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewTransferDtl.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` | `VIEWTRANSFERDTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-KOMITETREATY / VIEWTRANSFERDTL` |

### RDBList — 16 rule (`RULE-CONNECT-SQL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `GenerateImageID_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GENERATEIMAGEID_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GetAppName_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETAPPNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.T_FOLDER_IMAGE | — |
| `GetDataCNPOS.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `GETDATACNPOS` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=OS_AKSEPTASI_KLAIM | — |
| `GetEmailCeding_SQL.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GETEMAILCEDING_SQL` | [terverifikasi] sqlKind=QUERY; procs=GL.F_GET_EMAIL; sqlOps=SELECT | — |
| `GetKodeProdNonLife_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETKODEPRODNONLIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.KODE_PRODUKSI | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETSEQUENCENUMBER_SQL` |
| `GetSequenceNumber_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETSEQUENCENUMBER_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETTANGGALCLOSING_SQL` |
| `GetTokenStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETTOKENSTORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.GET_TOKEN_STORAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |
| `InsertClaimPNC.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `INSERTCLAIMPNC` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_KLAIM_PNC | — |
| `InsertHistoryAkseptasiPega_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTHISTORYAKSEPTASIPEGA_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=HISTORYAKSEPTASIPEGA | — |
| `InsertLOGDirectKasir_SQL.xml` | `ASM-FW-GCNMFW-WORK` | `INSERTLOGDIRECTKASIR_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.DIRECTTOKASIR_LOG | `ASM-FW-GISFW-WORK / RNM!INSERTLOGMOP_SQL` |
| `Insert_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERT_T_STORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GENERATEIMAGEID_SQL` |
| `SaveDataToOsAkseptasiNP.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `SAVEDATATOOSAKSEPTASINP` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP | `ASM-FW-GCNMFW-INT-V_POLIS / GCNM!SAVEDATATOOSAKSEPTASI` |
| `SaveOSClaim_SQL.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `SAVEOSCLAIM_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM | — |
| `SaveOSSubjectivity_SQL.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `SAVEOSSUBJECTIVITY_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_OS_AKSEP_SUBJECTIVITY | `ASM-FW-GCNMFW-INT-V_POLIS / ASM!SAVEOSCLAIM_SQL` |
| `SaveXOLClaim_SQL.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `SAVEXOLCLAIM_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.XOL2_AKSEP_KLAIM | `ASM-FW-GCNMFW-INT-V_POLIS / ASM!SAVEOSCLAIM_SQL` |
| `getStatusKonversi_SQL.xml` | `ASM-FW-GCNMFW-INT-OS_AKSEPTASI_KLAIM` | `GETSTATUSKONVERSI_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=REINSURANCE.TRLOSS_DETAIL_T | — |

### ReportDefinition — 3 rule (`RULE-OBJ-REPORT-DEFINITION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseBankGroup.xml` | `ASM-FW-GISFW-INT-LST_BANK_GROUP` | `BROWSEBANKGROUP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCurrency_RD.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `BROWSECURRENCY_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GetEmailUser_RD.xml` | `DATA-ADMIN-OPERATOR-ID` | `GETEMAILUSER_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-ADMIN-OPERATOR-ID / OPERATORSBYSKILL` |

### Section — 2 rule (`RULE-HTML-SECTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `ReinstatementPremiumDetails.xml` | `ASM-FW-GISFW-DATA-SPREADINGRISK` | `REINSTATEMENTPREMIUMDETAILS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowTransfer.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` | `SHOWTRANSFER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-KOMITETREATY / SHOWTRANSFER` |

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
| `IsKomiteLoop.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP` | `ISKOMITELOOP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-KOMITETREATY / ISKOMITELOOP` |
| `IsPEGAPROD.xml` | `@BASECLASS` | `ISPEGAPROD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `WORK- / ISPEGAPROD` |
| `IsPEGASyariah.xml` | `@BASECLASS` | `ISPEGASYARIAH` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### excludeXML — 1 rule (`RULE-OBJ-ACTIVITY`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `GetBase64Attachment.xml` | `ASM-FW-GCNMFW-WORK` | `GETBASE64ATTACHMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `WORK- / LOADATTACHMENTDATA` |

## 4. Tabrakan nama di dalam modul ini

### 4.1 Nama file dipakai di lebih dari satu tipe rule

_Tidak ada._

### 4.2 Rule `When` bernama sama di lebih dari satu class

_Tidak ada._

## 5. Keluarga klon (satu asal salinan dipakai >1 rule)

Dibaca dari `<pzOriginalInstanceKey>` yang sama pada beberapa file. Ini **sinyal kandidat
konsolidasi** saat migrasi, bukan bukti bahwa perilakunya identik — kesamaan perilaku baru
dapat dinyatakan setelah isinya dibandingkan di STEP D2.

| Asal salinan (class / nama) | Jumlah | File |
| --- | ---: | --- |
| `@BASECLASS / ISCLM` | 2 | `When/IsCLM.xml`, `When/IsCLMP.xml` |
| `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP / INSERTOSKLAIMCNP` | 3 | `Activity/InsertOSKlaimCNP.xml`, `Activity/InsertOSSubjectivityCNP.xml`, `Activity/InsertXOLKlaimCNP.xml` |
| `ASM-FW-GCNMFW-INT-V_POLIS / ASM!SAVEOSCLAIM_SQL` | 3 | `RDBList/SaveOSClaim_SQL.xml`, `RDBList/SaveOSSubjectivity_SQL.xml`, `RDBList/SaveXOLClaim_SQL.xml` |
| `ASM-FW-GCNMFW-WORK-KOMITETREATYNONPROP / KOMITEPOSTADJUSTMENT` | 2 | `Activity/KomitePostAdjustment.xml`, `Activity/KomitePostAdjustmentCWP.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GENERATEIMAGEID_SQL` | 2 | `RDBList/GenerateImageID_SQL.xml`, `RDBList/Insert_T_Storage_SQL.xml` |

## 6. Objek database yang disentuh

Diambil dari `<pyBrowseSQL>` pada rule `RDBList`. **Nama objek apa adanya dari SQL** —
kualifikasi skema tidak konsisten di korpus. **Tipe kolom, constraint, dan relasi tidak
diketahui** (OQ-001).

| Tabel/view (apa adanya di SQL) | Dirujuk oleh |
| --- | ---: |
| `HISTORYAKSEPTASIPEGA` | 1 rule |
| `OS_AKSEPTASI_KLAIM` | 1 rule |
| `POOLDATA.DIRECTTOKASIR_LOG` | 1 rule |
| `POOLDATA.KODE_PRODUKSI` | 1 rule |
| `POOLDATA.T_FOLDER_IMAGE` | 1 rule |
| `REINSURANCE.TRLOSS_DETAIL_T` | 1 rule |
| `T_STORAGE_IMAGE` | 1 rule |

### Stored procedure / function yang dipanggil

**Body prosedur ini TIDAK ADA di korpus** — logika di dalamnya tidak dapat direkonstruksi
dari XML. Lihat OQ-002.

| Procedure | Dipanggil oleh |
| --- | --- |
| `POOLDATA.XOL2_AKSEP_KLAIM` | `SaveXOLClaim_SQL.xml` |
| `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM` | `SaveOSClaim_SQL.xml` |
| `POOLDATA.PEGA_JSON_KLAIM_PNC` | `InsertClaimPNC.xml` |
| `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` | `GetSequenceNumber_SQL.xml` |
| `POOLDATA.PEGA_JSON_OS_AKSEP_SUBJECTIVITY` | `SaveOSSubjectivity_SQL.xml` |
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
| `insertClaimFinalOrClosed_NP.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `insertClaimFinalOrClosed_NP` | `URL` | — | **YA** |
| `insertClaimReject_NP.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `insertClaimReject_NP` | `URL` | — | **YA** |

