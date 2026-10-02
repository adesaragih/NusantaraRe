# Inventaris Rule — Claim Life

STEP D1, batch 3 (domain claim non-facultative). Dibangkitkan 2026-09-12 dari korpus READ-ONLY
`D:\XML\RNM_BRD\Claim Life\`.

**Batas step ini:** hanya struktur + metadata. Tidak ada kesimpulan logika bisnis di sini.
Penelusuran perilaku adalah STEP D2 (`../flows/`).

**Identitas rule** diambil dari `<pxInsName>` (lihat `../README.md` §6.4).  
Kolom **Asal salinan** berisi `<pzOriginalInstanceKey>` dan hanya diisi bila **berbeda**
dari identitas — artinya rule ini hasil *Save As* dari rule lain. Tanda `—` berarti sama.
Class `—` berarti `<pxInsName>` tidak berisi bagian class (berlaku pada `DataPage`, yang
beridentitas global, bukan per class).

Perintah audit (dijalankan dari `D:\XML\RNM_BRD\`):

```
find "Claim Life" -type f -name "*.xml" | wc -l
find "Claim Life" -type f -name "*.xml" | awk -F/ '{print $2}' | sort | uniq -c
grep -o "<pxInsName>[^<]*" "Claim Life/<tipe>/<file>.xml" | head -1
grep -o "<pzOriginalInstanceKey>[^<]*" "Claim Life/<tipe>/<file>.xml" | head -1
```

## 1. Ringkasan terukur

- Total file XML: **136** — cocok dengan `../README.md` §5.1
- **Identitas rule unik: 136** dari 136 file → tidak ada identitas ganda di dalam modul ini
- Class Pega unik: **23**
- RDBList menurut jenis SQL: PLSQL=7, QUERY=22 (`PLSQL` = blok anonim `BEGIN…END;`, umumnya memanggil stored procedure)
- Prefix pada kunci RDBList: `ASM`=7, `RNM`=22 (arti prefix **belum terverifikasi** — OQ-008)

| Tipe rule | Jumlah |
| --- | ---: |
| Activity | 51 |
| ConnectREST | 2 |
| DataTransform | 1 |
| DecisionTable | 1 |
| Flow | 1 |
| FlowAction | 16 |
| Harness | 3 |
| RDBList | 29 |
| ReportDefinition | 8 |
| Section | 20 |
| SystemSettings | 1 |
| When | 3 |
| **TOTAL** | **136** |

## 2. Sebaran class Pega

Class yang **tidak** berawalan `ASM-` adalah class bawaan Pega atau class portal — ruang
lingkup migrasinya belum ditentukan (**OQ-009**).

| Class | Jumlah rule | Aplikasi? |
| --- | ---: | --- |
| `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | 44 | aplikasi |
| `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | 27 | aplikasi |
| `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` | 13 | aplikasi |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | 12 | aplikasi |
| `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | 5 | aplikasi |
| `ASM-FW-GCNMFW-WORK` | 5 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` | 5 | aplikasi |
| `@BASECLASS` | 4 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-DATA-DIAGNOSELIFE` | 4 | aplikasi |
| `ASM-FW-GISFW-INT-POLICYJSON` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-BUSINESS` | 2 | aplikasi |
| `ASM-FW-GCNMFW-INT-EMAILKOMITE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-AGENT` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-DISEASE_LIFE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-MARKETINGOFFICER` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_LINK_SERVICE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-PRODUCTINWARD_LIFE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-RETROCESSIONLIFE` | 1 | aplikasi |
| `ASM-FW-GISFW-WORK-LIFE` | 1 | aplikasi |
| `ASSIGN-WORKLIST` | 1 | **bukan aplikasi** → OQ-009 |
| `DATA-WORKATTACH-FILE` | 1 | **bukan aplikasi** → OQ-009 |
| `LINKSERVICE` | 1 | **bukan aplikasi** → OQ-009 |

Rule di class bukan-aplikasi: **7** dari 136.

## 3. Daftar rule per tipe

### Activity — 51 rule (`RULE-OBJ-ACTIVITY`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `CountClaimAmountLife_Act.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` | `COUNTCLAIMAMOUNTLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE / SAVEOUTSTANDINGLIFE_ACT` |
| `CreateKMTLife_Act.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` | `CREATEKMTLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `DeleteDocument_Act.xml` | `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | `DELETEDOCUMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `DeleteGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `DELETEGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / GETURLGOOGLESTORAGE_ACT` |
| `DeletePesertaClaimLife.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `DELETEPESERTACLAIMLIFE` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `DownloadDocumentClaim.xml` | `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | `DOWNLOADDOCUMENTCLAIM` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `GetLinkService.xml` | `ASM-FW-GISFW-INT-M_LINK_SERVICE` | `GETLINKSERVICE` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetListKomiteLife.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` | `GETLISTKOMITELIFE` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE / CREATEKMTLIFE_ACT` |
| `GetUrlGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETURLGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / INSERTGOOGLESTORAGE_ACT` |
| `InsertDocument_Act.xml` | `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | `INSERTDOCUMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `InsertGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERTGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `InsertJsonClaimLife_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `INSERTJSONCLAIMLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GISFW-WORK-CLAIMLIFE / INSERTJSONPOLISLIFE_ACT` |
| `InsertLogServiceClaim.xml` | `ASM-FW-GCNMFW-WORK` | `INSERTLOGSERVICECLAIM` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `LoadDataPesertaSpesifik_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `LOADDATAPESERTASPESIFIK_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `LoadDataPeserta_Act.xml` | `ASM-FW-GCNMFW-WORK` | `LOADDATAPESERTA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `ASM-FW-GCNMFW-WORK-CLAIMLIFE / LOADDATAPESERTA_ACT` |
| `LoadDocumentLife_ACT.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `LOADDOCUMENTLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `NewAttachLife.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `NEWATTACHLIFE` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `NextPrev.xml` | `ASM-FW-GCNMFW-WORK` | `NEXTPREV` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `ObjSave_Act.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `OBJSAVE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `PreCaimLife_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `PRECAIMLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `ProtectCloseClaim_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `PROTECTCLOSECLAIM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `RejectOSClaimLife_Act.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` | `REJECTOSCLAIMLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SaveAdjustment_Act.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `SAVEADJUSTMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / ADDADJUSTMENTLIFE` |
| `SaveAttachLife.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `SAVEATTACHLIFE` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `SaveInsuredClaim_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `SAVEINSUREDCLAIM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SaveOutStandingLife_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `SAVEOUTSTANDINGLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=58 | — |
| `SavePesertaClaim.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `SAVEPESERTACLAIM` | [terverifikasi] pyActivityType=ACTIVITY; steps=18 | `ASM-FW-GISFW-WORK-CLAIMLIFE / SAVEPESERTACLAIM` |
| `SearchDiagnose_act.xml` | `ASM-FW-GISFW-DATA-DIAGNOSELIFE` | `SEARCHDIAGNOSE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / SEARCHDIAGNOSE_ACT` |
| `SearchPolicyHolder_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `SEARCHPOLICYHOLDER_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SelectAllClaimLife_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `SELECTALLCLAIMLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SendEmailKlaimLF.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` | `SENDEMAILKLAIMLF` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE / SENDEMAILKLAIMLIFE` |
| `SendEmailWithAttachments.xml` | `@BASECLASS` | `SENDEMAILWITHATTACHMENTS` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SendtoAdmin_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `SENDTOADMIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GCNMFW-WORK-CLAIMLIFE / DELETEDOCUMENTLIST_ACT` |
| `SendtoAdmin_Act1.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `SENDTOADMIN_ACT1` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GCNMFW-WORK-CLAIMLIFE / SENDTOADMIN_ACT` |
| `SendtoMedical_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `SENDTOMEDICAL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GCNMFW-WORK-CLAIMLIFE / SENDTOADMIN_ACT` |
| `SetClaimXOL_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `SETCLAIMXOL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GCNMFW-WORK-CLAIMLIFE / STNCLIFE_ACT` |
| `SetCurrencyID_Act.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` | `SETCURRENCYID_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetDisease.xml` | `ASM-FW-GISFW-DATA-DIAGNOSELIFE` | `SETDISEASE` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetIndexAdjustmentList.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `SETINDEXADJUSTMENTLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetMOClaim_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `SETMOCLAIM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `ASM-FW-GCNMFW-WORK-PNC / SETMOCLAIM_ACT` |
| `SetSTS_Reject.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `SETSTS_REJECT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / ADDADJUSTMENTLIFE` |
| `SpreadingClaimLife_Act.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` | `SPREADINGCLAIMLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=24 | — |
| `UpdateDateClaimLife_Act.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `UPDATEDATECLAIMLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `UploadCSVClaimLife_Act.xml` | `@BASECLASS` | `UPLOADCSVCLAIMLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `@BASECLASS / UPLOADCSVLIFEPREMIUM_ACT` |
| `ValidasiClaimReceived_Act.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `VALIDASICLAIMRECEIVED_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / VALIDASISTNC_ACT` |
| `ValidasiDOL_Act.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `VALIDASIDOL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / ADDADJUSTMENTLIFE` |
| `ValidasiSTNC_Act.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `VALIDASISTNC_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / VALIDASIDOL_ACT` |
| `getMaxPagination_Act.xml` | `ASM-FW-GCNMFW-WORK` | `GETMAXPAGINATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GCNMFW-WORK-CLAIMLIFE / GETMAXPAGINATION_ACT` |
| `serviceInsertArasapasClaimLife_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `SERVICEINSERTARASAPASCLAIMLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `setDetailClaim_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `SETDETAILCLAIM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `setVisibility_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `SETVISIBILITY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |

### ConnectREST — 2 rule (`RULE-CONNECT-REST`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `SERVICEGOOGLE` | [terverifikasi] pyServiceName=ServiceGoogle; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / GOOGLESTORAGE_UPLOAD` |
| `convertJsonNusareToProductionClaimLife.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `CONVERTJSONNUSARETOPRODUCTIONCLAIMLIFE` | [terverifikasi] pyServiceName=convertJsonNusareToProductionClaimLife; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | — |

### DataTransform — 1 rule (`RULE-OBJ-MODEL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `SetDisableAddButton.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `SETDISABLEADDBUTTON` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### DecisionTable — 1 rule (`RULE-DECLARE-DECISIONTABLE`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `GetMimeType.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETMIMETYPE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Flow — 1 rule (`RULE-OBJ-FLOW (pxObjClass)`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `Register_Flow.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `REGISTER_FLOW` | [terverifikasi] pyStartActivity=Start2 | — |

### FlowAction — 16 rule (`RULE-OBJ-FLOWACTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `Adjustment_Detail.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` | `ADJUSTMENT_DETAIL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `AkseptasiClaimLife.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `AKSEPTASICLAIMLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `AttachDocumentLife.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `ATTACHDOCUMENTLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CloseClaim.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `CLOSECLAIM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ConfirmDeleteAttachment.xml` | `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | `CONFIRMDELETEATTACHMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputRegisterClaimLife.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `INPUTREGISTERCLAIMLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `MedicalCheck.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `MEDICALCHECK` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-CLAIMLIFE / AKSEPTASICLAIMLIFE` |
| `OSClaimLife.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `OSCLAIMLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PL_DetailAction_ViewPolis.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `PL_DETAILACTION_VIEWPOLIS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / PL_DETAILACTION` |
| `RejectOSClaimLife.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` | `REJECTOSCLAIMLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `RetroClaimLife.xml` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` | `RETROCLAIMLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE / RETROLIFE` |
| `SendtoAdmin.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `SENDTOADMIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-CLAIMLIFE / AKSEPTASICLAIMLIFE` |
| `SendtoMedical.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `SENDTOMEDICAL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-CLAIMLIFE / SENDTOADMIN` |
| `ShowEditClaimLife.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `SHOWEDITCLAIMLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / PL_DETAILACTION` |
| `UploadCSV_ClaimLife.xml` | `@BASECLASS` | `UPLOADCSV_CLAIMLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / UPLOADCSV_LIFEPREMIUM` |
| `ViewClaimDetailLifeGCNM.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `VIEWCLAIMDETAILLIFEGCNM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / VIEWCLAIMDETAILLIFE` |

### Harness — 3 rule (`RULE-HTML-HARNESS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `Committe_Life.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` | `COMMITTE_LIFE` | [terverifikasi] pyInclude=1 | — |
| `Diagnose_Harness.xml` | `ASM-FW-GISFW-DATA-DIAGNOSELIFE` | `DIAGNOSE_HARNESS` | [terverifikasi] pyInclude=1 | — |
| `SearchPolicy_Harness.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `SEARCHPOLICY_HARNESS` | [terverifikasi] pyInclude=1 | — |

### RDBList — 29 rule (`RULE-CONNECT-SQL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `CountPesertaAkseptasiLifeHealth_SQL.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `COUNTPESERTAAKSEPTASILIFEHEALTH_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.OS_AKSEPTASI_KLAIM_LIFE | `ASM-FW-GCNMFW-WORK-CLAIMLIFE / RNM!COUNTPESERTAAKSEPTASILIFE_SQL` |
| `CountPesertaAkseptasiLife_SQL.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `COUNTPESERTAAKSEPTASILIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.OS_AKSEPTASI_KLAIM_LIFE | — |
| `DeleteStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `DELETESTORAGE_SQL` | [terverifikasi] sqlKind=QUERY | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GETLINKSTORAGE_SQL` |
| `GETTanggalClosing_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETTANGGALCLOSING_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TANGGAL_CLOSING | — |
| `GenerateImageID_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GENERATEIMAGEID_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `Generate_NoAccept_Life.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `GENERATE_NOACCEPT_LIFE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / ASM!GENERATE_NOKLAIM_LIFE` |
| `Generate_NoAccept_LifeRetro.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `GENERATE_NOACCEPT_LIFERETRO` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / ASM!GENERATE_NOACCEPT_LIFE` |
| `GetAcceptedNoCL.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `GETACCEPTEDNOCL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.OS_AKSEPTASI_KLAIM_LIFE | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / ASM!GENERATE_NOACCEPT_LIFE` |
| `GetAppName_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETAPPNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.T_FOLDER_IMAGE | — |
| `GetCurrencyID.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `GETCURRENCYID` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.CURRENCY | `ASM-FW-GCNMFW-WORK-CLAIMLIFE / RNM!GETPESERTACLAIM_SQL1` |
| `GetJsonProductLife.xml` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` | `GETJSONPRODUCTLIFE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYYEAR_LIFE,M_PRODUCT_LIFE | `ASM-FW-GISFW-INT-PRODUCT_LIFE / ASM!GETJSONPRODUCTLIFE` |
| `GetKodeProdLife_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETKODEPRODLIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.KODE_PRODUKSI | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETKODEPRODNONLIFE_SQL` |
| `GetLinkStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETLINKSTORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GETTOKENSTORAGE_SQL` |
| `GetPesertaClaim_sql.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `GETPESERTACLAIM_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_LIFE_PREMIUM_DETAIL | `ASM-FW-GCNMFW-WORK-CLAIMLIFE / RNM!GETPESERTACLAIM_SQL` |
| `GetPesertaClaim_sql1.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `GETPESERTACLAIM_SQL1` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_LIFE_PREMIUM_DETAIL | `ASM-FW-GCNMFW-WORK-CLAIMLIFE / RNM!GETPESERTACLAIM_SQL` |
| `GetProductLife.xml` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` | `GETPRODUCTLIFE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=PRODUCT_LIFE,M_PRODUCT_LIFE | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE / ASM!GETRATEPRODUCTLIFE` |
| `GetProductName.xml` | `ASM-FW-GISFW-INT-PRODUCTINWARD_LIFE` | `GETPRODUCTNAME` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.PRODUCTINWARD_LIFE | `ASM-FW-GISFW-INT-PRODUCTINWARD_LIFE / ASM!GETDATASUMMARYTREATY` |
| `GetRateRetro.xml` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` | `GETRATERETRO` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.RATE_LIFE | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE / ASM!GETRATEGENDER` |
| `GetRetroLife_SQL.xml` | `ASM-FW-GISFW-INT-RETROCESSIONLIFE` | `GETRETROLIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=RETROCESSIONLIFE | — |
| `GetSequenceNumber_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETSEQUENCENUMBER_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETTANGGALCLOSING_SQL` |
| `GetTokenStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETTOKENSTORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.GET_TOKEN_STORAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |
| `InsertJsonClaimLifeGCNM.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` | `INSERTJSONCLAIMLIFEGCNM` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_KLAIM_PNC | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY / ASM!INSERTJSONCLAIMLIFE` |
| `InsertJsonKlaimLife_sql.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `INSERTJSONKLAIMLIFE_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.OS_AKSEPTASI_KLAIM_LIFE | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / ASM!SAVEMASTERLPDET` |
| `InsertLogServiceClaim.xml` | `ASM-FW-GCNMFW-WORK` | `INSERTLOGSERVICECLAIM` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.MONITORING_KLAIM_LOG | — |
| `Insert_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERT_T_STORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GENERATEIMAGEID_SQL` |
| `UpdateDateClaimLife_SQL.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `UPDATEDATECLAIMLIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=POOLDATA.OS_AKSEPTASI_KLAIM_LIFE | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / RNM!INSERTJSONKLAIMLIFE_SQL` |
| `UpdateOsAkseptasiClaimLife_sql.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `UPDATEOSAKSEPTASICLAIMLIFE_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.OS_AKSEPTASI_KLAIM_LIFE | — |
| `Update_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `UPDATE_T_STORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |
| `getMaxPagination_sql.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `GETMAXPAGINATION_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_LIFE_PREMIUM_DETAIL | — |

### ReportDefinition — 8 rule (`RULE-OBJ-REPORT-DEFINITION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseBusinessLife_RD.xml` | `ASM-FW-GISFW-INT-BUSINESS` | `BROWSEBUSINESSLIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-BUSINESS / BROWSEBUSINESS_RD` |
| `BrowseCedingCoLife_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSECEDINGCOLIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-AGENT / BROWSECEDINGCO_RD` |
| `BrowseDiseaseLife_RD.xml` | `ASM-FW-GISFW-INT-DISEASE_LIFE` | `BROWSEDISEASELIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseFilterBusiness_RD.xml` | `ASM-FW-GISFW-INT-BUSINESS` | `BROWSEFILTERBUSINESS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseMarketingOfficer_RD.xml` | `ASM-FW-GISFW-INT-MARKETINGOFFICER` | `BROWSEMARKETINGOFFICER_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `FilterEmailKomiteWithLimit.xml` | `ASM-FW-GCNMFW-INT-EMAILKOMITE` | `FILTEREMAILKOMITEWITHLIMIT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InboxPremiumList.xml` | `ASSIGN-WORKLIST` | `INBOXPREMIUMLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InboxPremiumList_Claim.xml` | `ASM-FW-GISFW-WORK-LIFE` | `INBOXPREMIUMLIST_CLAIM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-LIFE / INBOXPREMIUMLIST` |

### Section — 20 rule (`RULE-HTML-SECTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AdjustmentDetail_Section.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` | `ADJUSTMENTDETAIL_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `AttachDocScreenLife.xml` | `DATA-WORKATTACH-FILE` | `ATTACHDOCSCREENLIFE` | [terverifikasi] pyInclude=1 | `DATA-WORKATTACH-FILE / PYATTACHMENTSCREEN` |
| `ClaimComite.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` | `CLAIMCOMITE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ClaimLifeDetailGCNM.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `CLAIMLIFEDETAILGCNM` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / CLAIMLIFEDETAIL` |
| `CloseClaim_Section.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `CLOSECLAIM_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ConfirmDeleteAttachment.xml` | `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | `CONFIRMDELETEATTACHMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DetailPolisLife.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `DETAILPOLISLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Diagnose_Section.xml` | `ASM-FW-GISFW-DATA-DIAGNOSELIFE` | `DIAGNOSE_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / DIAGNOSE_SECTION` |
| `DocumentLife.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `DOCUMENTLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `EditDateClaimLife_Section.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `EDITDATECLAIMLIFE_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputAkseptasiClaimLife.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `INPUTAKSEPTASICLAIMLIFE` | [terverifikasi] pyInclude=2 | `ASM-FW-GCNMFW-WORK-CLAIMLIFE / MEDICALCHECKCLAIMLIFE` |
| `InputOSClaimLife.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `INPUTOSCLAIMLIFE` | [terverifikasi] pyInclude=2 | `ASM-FW-GCNMFW-WORK-CLAIMLIFE / INPUTAKSEPTASICLAIMLIFE` |
| `InputRegisterClaimLife.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `INPUTREGISTERCLAIMLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `MedicalCheckClaimLife.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `MEDICALCHECKCLAIMLIFE` | [terverifikasi] pyInclude=2 | `ASM-FW-GCNMFW-WORK-CLAIMLIFE / INPUTOSCLAIMLIFE` |
| `PL_DetailViewPolis_Sec.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `PL_DETAILVIEWPOLIS_SEC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / PL_DETAIL_SEC` |
| `RejectOSClaimLife_Sec.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` | `REJECTOSCLAIMLIFE_SEC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `RetroDetailClaimLife.xml` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` | `RETRODETAILCLAIMLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE / RETRODETAILLIFE` |
| `SearchPolicy_Section.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `SEARCHPOLICY_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SendtoAdmin_Section.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `SENDTOADMIN_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-CLAIMLIFE / SEARCHPOLICY_SECTION` |
| `SendtoMedical_Section.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `SENDTOMEDICAL_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-CLAIMLIFE / SENDTOADMIN_SECTION` |

### SystemSettings — 1 rule (`RULE-ADMIN-SYSTEM-SETTINGS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `LinkService.xml` | `LINKSERVICE` | `LINKSERVICE` | [terverifikasi] pyPurpose=LinkService | — |

### When — 3 rule (`RULE-OBJ-WHEN`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `IsPEGAPROD.xml` | `@BASECLASS` | `ISPEGAPROD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `WORK- / ISPEGAPROD` |
| `IsSendtoAdmin.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `ISSENDTOADMIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsSendtoMedical.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `ISSENDTOMEDICAL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-CLAIMLIFE / ISSENDTOADMIN` |

## 4. Tabrakan nama di dalam modul ini

### 4.1 Nama file dipakai di lebih dari satu tipe rule

| Nama file | Tipe | Class |
| --- | --- | --- |
| `InsertLogServiceClaim.xml` | Activity, RDBList | `ASM-FW-GCNMFW-WORK, ASM-FW-GCNMFW-WORK` |
| `InputRegisterClaimLife.xml` | FlowAction, Section | `ASM-FW-GCNMFW-WORK-CLAIMLIFE, ASM-FW-GCNMFW-WORK-CLAIMLIFE` |
| `ConfirmDeleteAttachment.xml` | FlowAction, Section | `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM, ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` |

Total: **3** nama file.

### 4.2 Rule `When` bernama sama di lebih dari satu class

_Tidak ada._

## 5. Keluarga klon (satu asal salinan dipakai >1 rule)

Dibaca dari `<pzOriginalInstanceKey>` yang sama pada beberapa file. Ini **sinyal kandidat
konsolidasi** saat migrasi, bukan bukti bahwa perilakunya identik — kesamaan perilaku baru
dapat dinyatakan setelah isinya dibandingkan di STEP D2.

| Asal salinan (class / nama) | Jumlah | File |
| --- | ---: | --- |
| `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETTANGGALCLOSING_SQL` | 2 | `RDBList/GETTanggalClosing_SQL.xml`, `RDBList/GetSequenceNumber_SQL.xml` |
| `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM / CONFIRMDELETEATTACHMENT` | 2 | `FlowAction/ConfirmDeleteAttachment.xml`, `Section/ConfirmDeleteAttachment.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` | 2 | `RDBList/GetTokenStorage_SQL.xml`, `RDBList/Update_T_Storage_SQL.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / INSERTGOOGLESTORAGE_ACT` | 2 | `Activity/GetUrlGoogleStorage_Act.xml`, `Activity/InsertGoogleStorage_Act.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMLIFE / RNM!COUNTPESERTAAKSEPTASILIFE_SQL` | 2 | `RDBList/CountPesertaAkseptasiLifeHealth_SQL.xml`, `RDBList/CountPesertaAkseptasiLife_SQL.xml` |
| `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / ADDADJUSTMENTLIFE` | 3 | `Activity/SaveAdjustment_Act.xml`, `Activity/SetSTS_Reject.xml`, `Activity/ValidasiDOL_Act.xml` |
| `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE / CREATEKMTLIFE_ACT` | 2 | `Activity/CreateKMTLife_Act.xml`, `Activity/GetListKomiteLife.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMLIFE / RNM!GETPESERTACLAIM_SQL` | 2 | `RDBList/GetPesertaClaim_sql.xml`, `RDBList/GetPesertaClaim_sql1.xml` |
| `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / ASM!GENERATE_NOACCEPT_LIFE` | 2 | `RDBList/Generate_NoAccept_LifeRetro.xml`, `RDBList/GetAcceptedNoCL.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMLIFE / INPUTREGISTERCLAIMLIFE` | 2 | `FlowAction/InputRegisterClaimLife.xml`, `Section/InputRegisterClaimLife.xml` |
| `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / PL_DETAILACTION` | 2 | `FlowAction/PL_DetailAction_ViewPolis.xml`, `FlowAction/ShowEditClaimLife.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMLIFE / AKSEPTASICLAIMLIFE` | 3 | `FlowAction/AkseptasiClaimLife.xml`, `FlowAction/MedicalCheck.xml`, `FlowAction/SendtoAdmin.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMLIFE / SENDTOADMIN_ACT` | 2 | `Activity/SendtoAdmin_Act1.xml`, `Activity/SendtoMedical_Act.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMLIFE / SEARCHPOLICY_SECTION` | 2 | `Section/SearchPolicy_Section.xml`, `Section/SendtoAdmin_Section.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMLIFE / ISSENDTOADMIN` | 2 | `When/IsSendtoAdmin.xml`, `When/IsSendtoMedical.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GENERATEIMAGEID_SQL` | 2 | `RDBList/GenerateImageID_SQL.xml`, `RDBList/Insert_T_Storage_SQL.xml` |

## 6. Objek database yang disentuh

Diambil dari `<pyBrowseSQL>` pada rule `RDBList`. **Nama objek apa adanya dari SQL** —
kualifikasi skema tidak konsisten di korpus. **Tipe kolom, constraint, dan relasi tidak
diketahui** (OQ-001).

| Tabel/view (apa adanya di SQL) | Dirujuk oleh |
| --- | ---: |
| `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` | 6 rule |
| `POOLDATA.M_LIFE_PREMIUM_DETAIL` | 3 rule |
| `T_STORAGE_IMAGE` | 3 rule |
| `M_PRODUCT_LIFE` | 2 rule |
| `POOLDATA.CURRENCY` | 1 rule |
| `POOLDATA.KODE_PRODUKSI` | 1 rule |
| `POOLDATA.MONITORING_KLAIM_LOG` | 1 rule |
| `POOLDATA.PRODUCTINWARD_LIFE` | 1 rule |
| `POOLDATA.RATE_LIFE` | 1 rule |
| `POOLDATA.TANGGAL_CLOSING` | 1 rule |
| `POOLDATA.T_FOLDER_IMAGE` | 1 rule |
| `PRODUCT_LIFE` | 1 rule |
| `RETROCESSIONLIFE` | 1 rule |
| `TREATYYEAR_LIFE` | 1 rule |

### Stored procedure / function yang dipanggil

**Body prosedur ini TIDAK ADA di korpus** — logika di dalamnya tidak dapat direkonstruksi
dari XML. Lihat OQ-002.

| Procedure | Dipanggil oleh |
| --- | --- |
| `POOLDATA.PEGA_JSON_KLAIM_PNC` | `InsertJsonClaimLifeGCNM.xml` |
| `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` | `GetSequenceNumber_SQL.xml` |
| `POOLDATA.GET_TOKEN_STORAGE` | `GetTokenStorage_SQL.xml` |

## 7. Integrasi eksternal (`ConnectREST`)

Alamat tujuan dicatat sebagai **konfigurasi**, bukan URL literal. `baseURL=SETTING` berarti
alamat diambil dari rule `SystemSettings` yang disebut di kolom Setting; nilai sebenarnya
**tidak ada di korpus**. Kolom terakhir menandai rule yang justru memuat URL literal —
itu temuan, bukan fakta arsitektur (lihat `_summary.md`).

| File | Class | pyServiceName | Sumber base URL | Setting | URL literal di dalam rule? |
| --- | --- | --- | --- | --- | --- |
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `ServiceGoogle` | `SETTING` | `LinkService!LinkService` | tidak |
| `convertJsonNusareToProductionClaimLife.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` | `convertJsonNusareToProductionClaimLife` | `SETTING` | `LinkService!LinkService` | tidak |

