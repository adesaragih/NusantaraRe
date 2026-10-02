# Inventaris Rule — Komite Claim Life

STEP D1, batch 3 (domain claim non-facultative). Dibangkitkan 2026-09-12 dari korpus READ-ONLY
`D:\XML\RNM_BRD\Komite Claim Life\`.

**Batas step ini:** hanya struktur + metadata. Tidak ada kesimpulan logika bisnis di sini.
Penelusuran perilaku adalah STEP D2 (`../flows/`).

**Identitas rule** diambil dari `<pxInsName>` (lihat `../README.md` §6.4).  
Kolom **Asal salinan** berisi `<pzOriginalInstanceKey>` dan hanya diisi bila **berbeda**
dari identitas — artinya rule ini hasil *Save As* dari rule lain. Tanda `—` berarti sama.
Class `—` berarti `<pxInsName>` tidak berisi bagian class (berlaku pada `DataPage`, yang
beridentitas global, bukan per class).

Perintah audit (dijalankan dari `D:\XML\RNM_BRD\`):

```
find "Komite Claim Life" -type f -name "*.xml" | wc -l
find "Komite Claim Life" -type f -name "*.xml" | awk -F/ '{print $2}' | sort | uniq -c
grep -o "<pxInsName>[^<]*" "Komite Claim Life/<tipe>/<file>.xml" | head -1
grep -o "<pzOriginalInstanceKey>[^<]*" "Komite Claim Life/<tipe>/<file>.xml" | head -1
```

## 1. Ringkasan terukur

- Total file XML: **47** — cocok dengan `../README.md` §5.1
- **Identitas rule unik: 47** dari 47 file → tidak ada identitas ganda di dalam modul ini
- Class Pega unik: **18**
- RDBList menurut jenis SQL: PLSQL=6, QUERY=8 (`PLSQL` = blok anonim `BEGIN…END;`, umumnya memanggil stored procedure)
- Prefix pada kunci RDBList: `ASM`=1, `RNM`=12, `GCNM`=1 (arti prefix **belum terverifikasi** — OQ-008)

| Tipe rule | Jumlah |
| --- | ---: |
| Activity | 18 |
| ConnectREST | 1 |
| DecisionTable | 1 |
| Flow | 1 |
| FlowAction | 3 |
| RDBList | 14 |
| ReportDefinition | 3 |
| Section | 3 |
| SystemSettings | 1 |
| When | 2 |
| **TOTAL** | **47** |

## 2. Sebaran class Pega

Class yang **tidak** berawalan `ASM-` adalah class bawaan Pega atau class portal — ruang
lingkup migrasinya belum ditentukan (**OQ-009**).

| Class | Jumlah rule | Aplikasi? |
| --- | ---: | --- |
| `ASM-FW-GCNMFW-WORK-KOMITELIFE` | 12 | aplikasi |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | 10 | aplikasi |
| `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | 5 | aplikasi |
| `ASM-FW-GISFW-INT-POLICYJSON` | 3 | aplikasi |
| `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | 2 | aplikasi |
| `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` | 2 | aplikasi |
| `@BASECLASS` | 1 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | 1 | aplikasi |
| `ASM-FW-GCNMFW-INT-V_POLIS` | 1 | aplikasi |
| `ASM-FW-GCNMFW-WORK` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-DISEASE_LIFE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_LINK_SERVICE` | 1 | aplikasi |
| `DATA-ADMIN-OPERATOR-ID` | 1 | **bukan aplikasi** → OQ-009 |
| `DATA-WORKATTACH-FILE` | 1 | **bukan aplikasi** → OQ-009 |
| `LINK-ATTACHMENT` | 1 | **bukan aplikasi** → OQ-009 |
| `LINKSERVICE` | 1 | **bukan aplikasi** → OQ-009 |

Rule di class bukan-aplikasi: **5** dari 47.

## 3. Daftar rule per tipe

### Activity — 18 rule (`RULE-OBJ-ACTIVITY`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `DownloadDocumentClaim.xml` | `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | `DOWNLOADDOCUMENTCLAIM` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `GetLinkService.xml` | `ASM-FW-GISFW-INT-M_LINK_SERVICE` | `GETLINKSERVICE` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetUrlGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETURLGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / INSERTGOOGLESTORAGE_ACT` |
| `HitServiceToKasirKMTLife_Act.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` | `HITSERVICETOKASIRKMTLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | — |
| `InsertDocument_Act.xml` | `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | `INSERTDOCUMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `InsertGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERTGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `InsertJsonClaimLife_Act.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` | `INSERTJSONCLAIMLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `KomitePostAdjustment.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` | `KOMITEPOSTADJUSTMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=36 | — |
| `KomiteRouter.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` | `KOMITEROUTER` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `LoadDocumentKomiteLife_ACT.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` | `LOADDOCUMENTKOMITELIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `LoadDocumentLife_ACT.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `LOADDOCUMENTLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `NewAttachLife.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `NEWATTACHLIFE` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `PreKomiteLife.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` | `PREKOMITELIFE` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `PrintAkseptasiPDF.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` | `PRINTAKSEPTASIPDF` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `SaveAttachLife.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `SAVEATTACHLIFE` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `SendEmailKlaimLife.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` | `SENDEMAILKLAIMLIFE` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | — |
| `SetInformationData.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` | `SETINFORMATIONDATA` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SetRemarkKomiteLife.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` | `SETREMARKKOMITELIFE` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |

### ConnectREST — 1 rule (`RULE-CONNECT-REST`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `SERVICEGOOGLE` | [terverifikasi] pyServiceName=ServiceGoogle; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / GOOGLESTORAGE_UPLOAD` |

### DecisionTable — 1 rule (`RULE-DECLARE-DECISIONTABLE`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `GetMimeType.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETMIMETYPE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Flow — 1 rule (`RULE-OBJ-FLOW (pxObjClass)`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `KomiteLife_Flow.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` | `KOMITELIFE_FLOW` | [terverifikasi] pyStartActivity=Start2 | — |

### FlowAction — 3 rule (`RULE-OBJ-FLOWACTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AttachDocumentLife.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `ATTACHDOCUMENTLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `RetroLife.xml` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` | `RETROLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewTransferDtl.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` | `VIEWTRANSFERDTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### RDBList — 14 rule (`RULE-CONNECT-SQL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseCurrency.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `BROWSECURRENCY` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CURRENCY | — |
| `GETTanggalClosing_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETTANGGALCLOSING_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TANGGAL_CLOSING | — |
| `GenerateImageID_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GENERATEIMAGEID_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GetAppName_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETAPPNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.T_FOLDER_IMAGE | — |
| `GetEmailCeding_SQL.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GETEMAILCEDING_SQL` | [terverifikasi] sqlKind=QUERY; procs=GL.F_GET_EMAIL; sqlOps=SELECT | — |
| `GetKodeProdLife_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETKODEPRODLIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.KODE_PRODUKSI | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETKODEPRODNONLIFE_SQL` |
| `GetLinkStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETLINKSTORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GETTOKENSTORAGE_SQL` |
| `GetSequenceNumber_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETSEQUENCENUMBER_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETTANGGALCLOSING_SQL` |
| `GetTokenStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETTOKENSTORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.GET_TOKEN_STORAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |
| `InsertJsonClaimLifeGCNM.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` | `INSERTJSONCLAIMLIFEGCNM` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_KLAIM_PNC | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY / ASM!INSERTJSONCLAIMLIFE` |
| `InsertLOGDirectKasir_SQL.xml` | `ASM-FW-GCNMFW-WORK` | `INSERTLOGDIRECTKASIR_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.DIRECTTOKASIR_LOG | `ASM-FW-GISFW-WORK / RNM!INSERTLOGMOP_SQL` |
| `Insert_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERT_T_STORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GENERATEIMAGEID_SQL` |
| `UpdateOsAkseptasiClaimLife_sql.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `UPDATEOSAKSEPTASICLAIMLIFE_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.OS_AKSEPTASI_KLAIM_LIFE | — |
| `Update_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `UPDATE_T_STORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |

### ReportDefinition — 3 rule (`RULE-OBJ-REPORT-DEFINITION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseDiseaseLife_RD.xml` | `ASM-FW-GISFW-INT-DISEASE_LIFE` | `BROWSEDISEASELIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GetAllAttachmentsLife.xml` | `LINK-ATTACHMENT` | `GETALLATTACHMENTSLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `LINK-ATTACHMENT / GCNMGETALLATTACHMENTS` |
| `GetEmailUser_RD.xml` | `DATA-ADMIN-OPERATOR-ID` | `GETEMAILUSER_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-ADMIN-OPERATOR-ID / OPERATORSBYSKILL` |

### Section — 3 rule (`RULE-HTML-SECTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AttachDocScreenLife.xml` | `DATA-WORKATTACH-FILE` | `ATTACHDOCSCREENLIFE` | [terverifikasi] pyInclude=1 | `DATA-WORKATTACH-FILE / PYATTACHMENTSCREEN` |
| `RetroDetailLife.xml` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` | `RETRODETAILLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowTransfer.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` | `SHOWTRANSFER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### SystemSettings — 1 rule (`RULE-ADMIN-SYSTEM-SETTINGS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `LinkService.xml` | `LINKSERVICE` | `LINKSERVICE` | [terverifikasi] pyPurpose=LinkService | — |

### When — 2 rule (`RULE-OBJ-WHEN`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `IsKomiteLoop.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` | `ISKOMITELOOP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsPEGAPROD.xml` | `@BASECLASS` | `ISPEGAPROD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `WORK- / ISPEGAPROD` |

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
| `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETTANGGALCLOSING_SQL` | 2 | `RDBList/GETTanggalClosing_SQL.xml`, `RDBList/GetSequenceNumber_SQL.xml` |
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
| `CURRENCY` | 1 rule |
| `POOLDATA.DIRECTTOKASIR_LOG` | 1 rule |
| `POOLDATA.KODE_PRODUKSI` | 1 rule |
| `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` | 1 rule |
| `POOLDATA.TANGGAL_CLOSING` | 1 rule |
| `POOLDATA.T_FOLDER_IMAGE` | 1 rule |

### Stored procedure / function yang dipanggil

**Body prosedur ini TIDAK ADA di korpus** — logika di dalamnya tidak dapat direkonstruksi
dari XML. Lihat OQ-002.

| Procedure | Dipanggil oleh |
| --- | --- |
| `POOLDATA.PEGA_JSON_KLAIM_PNC` | `InsertJsonClaimLifeGCNM.xml` |
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
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `ServiceGoogle` | `SETTING` | `LinkService!LinkService` | tidak |

