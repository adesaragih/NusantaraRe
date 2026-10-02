# Inventaris Rule — Master Product Name Life

STEP D1, batch 4 (domain life, master & outward). Dibangkitkan 2026-09-12 dari korpus READ-ONLY
`D:\XML\RNM_BRD\Master Product Name Life\`.

**Batas step ini:** hanya struktur + metadata. Tidak ada kesimpulan logika bisnis di sini.
Penelusuran perilaku adalah STEP D2 (`../flows/`).

**Identitas rule** diambil dari `<pxInsName>` (lihat `../README.md` §6.4).  
Kolom **Asal salinan** berisi `<pzOriginalInstanceKey>` dan hanya diisi bila **berbeda**
dari identitas — artinya rule ini hasil *Save As* dari rule lain. Tanda `—` berarti sama.
Class `—` berarti `<pxInsName>` tidak berisi bagian class (berlaku pada `DataPage`, yang
beridentitas global, bukan per class).

Perintah audit (dijalankan dari `D:\XML\RNM_BRD\`):

```
find "Master Product Name Life" -type f -name "*.xml" | wc -l
find "Master Product Name Life" -type f -name "*.xml" | awk -F/ '{print $2}' | sort | uniq -c
grep -o "<pxInsName>[^<]*" "Master Product Name Life/<tipe>/<file>.xml" | head -1
grep -o "<pzOriginalInstanceKey>[^<]*" "Master Product Name Life/<tipe>/<file>.xml" | head -1
```

## 1. Ringkasan terukur

- Total file XML: **114** — cocok dengan `../README.md` §5.1
- **Identitas rule unik: 114** dari 114 file → tidak ada identitas ganda di dalam modul ini
- Class Pega unik: **27**
- RDBList menurut jenis SQL: PLSQL=6, QUERY=16 (`PLSQL` = blok anonim `BEGIN…END;`, umumnya memanggil stored procedure)
- Prefix pada kunci RDBList: `ASM`=13, `RNM`=9 (arti prefix **belum terverifikasi** — OQ-008)

| Tipe rule | Jumlah |
| --- | ---: |
| Activity | 37 |
| ConnectREST | 1 |
| DataTransform | 10 |
| DecisionTable | 1 |
| FlowAction | 11 |
| Harness | 1 |
| RDBList | 22 |
| ReportDefinition | 17 |
| Section | 12 |
| SystemSettings | 1 |
| When | 1 |
| **TOTAL** | **114** |

## 2. Sebaran class Pega

Class yang **tidak** berawalan `ASM-` adalah class bawaan Pega atau class portal — ruang
lingkup migrasinya belum ditentukan (**OQ-009**).

| Class | Jumlah rule | Aplikasi? |
| --- | ---: | --- |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE` | 55 | aplikasi |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | 12 | aplikasi |
| `DATA-PORTAL` | 8 | **bukan aplikasi** → OQ-009 |
| `@BASECLASS` | 7 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-INT-TREATY_IN` | 6 | aplikasi |
| `ASM-FW-GISFW-DATA-PLAN` | 3 | aplikasi |
| `ASM-FW-GISFW-DATA-UNDERWRITINGLIMIT` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-CURRENCY` | 2 | aplikasi |
| `ASM-FW-GISFW-DATA-OFFERFACIN` | 1 | aplikasi |
| `ASM-FW-GISFW-DATA-TREATYINLIMITS` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-AGENT` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-BENEFIT_LIFE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-BUSINESS` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-CAUSEOFLOSS_LIFE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-CLIENT` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_LINK_SERVICE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_RATE_LIFE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-OFFERJSON` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-PRODUCTINWARD_LIFE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-PRODUCT_TYPE_LIFE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-RATE_LIFE_SUMMARY` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-RICOMM_LIFE_SUMMARY` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-RIRISK_LIFE_SUMMARY` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYYEAR` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` | 1 | aplikasi |
| `ASSIGN-WORKLIST` | 1 | **bukan aplikasi** → OQ-009 |
| `LINKSERVICE` | 1 | **bukan aplikasi** → OQ-009 |

Rule di class bukan-aplikasi: **17** dari 114.

## 3. Daftar rule per tipe

### Activity — 37 rule (`RULE-OBJ-ACTIVITY`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AddCommentList_Act.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `ADDCOMMENTLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CheckDuplicateOffer.xml` | `DATA-PORTAL` | `CHECKDUPLICATEOFFER` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `ConvertHistoryDate.xml` | `DATA-PORTAL` | `CONVERTHISTORYDATE` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CopyFinancialWriting.xml` | `ASM-FW-GISFW-DATA-UNDERWRITINGLIMIT` | `COPYFINANCIALWRITING` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CopyUnderWritingLimit.xml` | `ASM-FW-GISFW-DATA-UNDERWRITINGLIMIT` | `COPYUNDERWRITINGLIMIT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CountMaxReasured_Act.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `COUNTMAXREASURED_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `CountMaxSumReasured_Act.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `COUNTMAXSUMREASURED_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GISFW-INT-PRODUCT_LIFE / COUNTMAXREASURED_ACT` |
| `DeleteAttacProdName_act.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `DELETEATTACPRODNAME_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `DATA-PORTAL / DELETE_ACT` |
| `DeleteGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `DELETEGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / GETURLGOOGLESTORAGE_ACT` |
| `DownloadAll_Act.xml` | `DATA-PORTAL` | `DOWNLOADALL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | — |
| `DownloadAttProdName_Act.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `DOWNLOADATTPRODNAME_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `ASM-FW-GISFW-INT-TREATY_IN / DOWNLOADATTACHMENT2` |
| `GenerateUpload_Act.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `GENERATEUPLOAD_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `GetLinkService.xml` | `ASM-FW-GISFW-INT-M_LINK_SERVICE` | `GETLINKSERVICE` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetReinsTypeOR_Life.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `GETREINSTYPEOR_LIFE` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GISFW-INT-PRODUCT_LIFE / GETREINSTYPELIST_LIFE` |
| `GetUrlGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETURLGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / INSERTGOOGLESTORAGE_ACT` |
| `InsertGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERTGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `LoadAttachment.xml` | `DATA-PORTAL` | `LOADATTACHMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `WORK- / LOADATTACHMENT` |
| `LoadAttachmentProdName.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `LOADATTACHMENTPRODNAME` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `DATA-PORTAL / LOADATTACHMENT` |
| `NewProductLife.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `NEWPRODUCTLIFE` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `@BASECLASS / NEWINPUTTREATYYEAR_LIFE_ACT` |
| `ProductNameSaveAttachment.xml` | `@BASECLASS` | `PRODUCTNAMESAVEATTACHMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | `@BASECLASS / TREATYSAVEATTACHMENT` |
| `ProteksiPlanListLife.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `PROTEKSIPLANLISTLIFE` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `ASM-FW-GISFW-INT-PRODUCT_LIFE / VALIDASIBENEFITLIFE` |
| `SaveInwardProductName_Act.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `SAVEINWARDPRODUCTNAME_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `ASM-FW-GISFW-INT-PRODUCT_LIFE / SAVEPRODUCTNAME_ACT` |
| `SaveProductName_Act.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `SAVEPRODUCTNAME_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=16 | — |
| `SearchPolicyHolder_act.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `SEARCHPOLICYHOLDER_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GISFW-WORK-LIFE / SEARCHPOLICYHOLDER_ACT` |
| `SetCategoryAttachTreatyin.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SETCATEGORYATTACHTREATYIN` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `ASM-FW-GISFW-DATA-OFFERFACIN / SETCATEGORYATTACH` |
| `SetCategory_act.xml` | `@BASECLASS` | `SETCATEGORY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `DATA-PORTAL / SETCATEGORY_ACT` |
| `SetParamRate.xml` | `@BASECLASS` | `SETPARAMRATE` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetProductName.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `SETPRODUCTNAME` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | `@BASECLASS / SETTREATYYEARLIFE_ACT` |
| `SetProductNameInward.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `SETPRODUCTNAMEINWARD` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `ASM-FW-GISFW-INT-PRODUCT_LIFE / SETPRODUCTNAME` |
| `SetRIRate.xml` | `ASM-FW-GISFW-DATA-PLAN` | `SETRIRATE` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetReinstatementPct.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `SETREINSTATEMENTPCT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SetTreatyIn_Act.xml` | `DATA-PORTAL` | `SETTREATYIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | — |
| `SetTreatyName_Act.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `SETTREATYNAME_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GISFW-INT-PRODUCT_LIFE / COUNTMAXSUMREASURED_ACT` |
| `TreatyInDownloadAll.xml` | `DATA-PORTAL` | `TREATYINDOWNLOADALL` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `TreatyInInputVis.xml` | `DATA-PORTAL` | `TREATYININPUTVIS` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `TreatyInitAttach.xml` | `@BASECLASS` | `TREATYINITATTACH` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / SFAGISINITATTACH` |
| `TreatySetReinstatement.xml` | `DATA-PORTAL` | `TREATYSETREINSTATEMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |

### ConnectREST — 1 rule (`RULE-CONNECT-REST`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `SERVICEGOOGLE` | [terverifikasi] pyServiceName=ServiceGoogle; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / GOOGLESTORAGE_UPLOAD` |

### DataTransform — 10 rule (`RULE-OBJ-MODEL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `CopyProduct.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `COPYPRODUCT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `HideCreateLife.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `HIDECREATELIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetViewEdit.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `SETVIEWEDIT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PRODUCT_LIFE / COPYPRODUCT` |
| `TreatyInIDSetPyPortal.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `TREATYINIDSETPYPORTAL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / TREATYINIDSETPYPORTAL` |
| `setCauseOfLoss_DT.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `SETCAUSEOFLOSS_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PRODUCT_LIFE / SETRIRISK_DT` |
| `setCeding_DT.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `SETCEDING_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `setCurrency_DT.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `SETCURRENCY_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PRODUCT_LIFE / SETPOLICYHOLDER_DT` |
| `setPolicyHolder_DT.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `SETPOLICYHOLDER_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PRODUCT_LIFE / SETRIRISK_DT` |
| `setRIRISK_DT.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `SETRIRISK_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `setSOB_DT.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `SETSOB_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### DecisionTable — 1 rule (`RULE-DECLARE-DECISIONTABLE`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `GetMimeType.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETMIMETYPE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### FlowAction — 11 rule (`RULE-OBJ-FLOWACTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `ChooseCauseOfLoss.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `CHOOSECAUSEOFLOSS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PRODUCT_LIFE / CHOOSERIRISK` |
| `ChooseCeding.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `CHOOSECEDING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChooseCurrency.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `CHOOSECURRENCY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PRODUCT_LIFE / CHOOSEPOLICYHOLDER` |
| `ChoosePolicyHolder.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `CHOOSEPOLICYHOLDER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PRODUCT_LIFE / CHOOSERIRISK` |
| `ChooseRIRate.xml` | `ASM-FW-GISFW-DATA-PLAN` | `CHOOSERIRATE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChooseRIRisk.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `CHOOSERIRISK` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChooseSOB.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `CHOOSESOB` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `EditProductName_Confirm.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `EDITPRODUCTNAME_CONFIRM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ProductNameAttachContent.xml` | `@BASECLASS` | `PRODUCTNAMEATTACHCONTENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / TREATYATTACHCONTENT` |
| `SaveProductName_Confirm.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `SAVEPRODUCTNAME_CONFIRM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewRate.xml` | `@BASECLASS` | `VIEWRATE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Harness — 1 rule (`RULE-HTML-HARNESS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `InwardProductName.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `INWARDPRODUCTNAME` | [terverifikasi] pyInclude=1 | — |

### RDBList — 22 rule (`RULE-CONNECT-SQL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseReinstypeOR_SQL.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `BROWSEREINSTYPEOR_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TREATYCONTRACT_LIFE,POOLDATA.TREATYYEAR_LIFE | `ASM-FW-GISFW-INT-PRODUCT_LIFE / ASM!BROWSEREINSTYPELIST_SQL` |
| `BrowseTreatyIn.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `BROWSETREATYIN` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_TREATY_IN,POOLDATA.M_TREATY_IN_EDM | `ASM-FW-GISFW-INT-PRODUCT_LIFE / ASM!BROWSEUNDERWRITINGLIST` |
| `BrowseUnderwritingList.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `BROWSEUNDERWRITINGLIST` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_PRODUCT_LIFE | — |
| `CategoryAttach_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `CATEGORYATTACH_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CATEGORY_ATTACH_REAS | — |
| `DeleteAttachProdName_Sql.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `DELETEATTACHPRODNAME_SQL` | [terverifikasi] sqlKind=QUERY | `ASM-FW-GISFW-INT-PRODUCT_LIFE / ASM!GETATTACHMENTPRODNAME_SQL` |
| `DeleteStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `DELETESTORAGE_SQL` | [terverifikasi] sqlKind=QUERY | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GETLINKSTORAGE_SQL` |
| `GenerateImageID_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GENERATEIMAGEID_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GetAllAttachment2_Sql.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `GETALLATTACHMENT2_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_ATTACHMENTTREATY_2 | `ASM-FW-GISFW-INT-TREATY_IN / ASM!GETATTACHMENT2_SQL` |
| `GetAppName_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETAPPNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.T_FOLDER_IMAGE | — |
| `GetAttachment2_Sql.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `GETATTACHMENT2_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_ATTACHMENTTREATY_2 | `ASM-FW-GISFW-INT-TREATY_IN / ASM!GETATTACHMENT_SQL` |
| `GetAttachmentProdName_Sql.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `GETATTACHMENTPRODNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_ATTACHMENTPRODUCTNAME | `ASM-FW-GISFW-INT-TREATY_IN / ASM!GETATTACHMENT_SQL` |
| `GetCountClaim.xml` | `ASSIGN-WORKLIST` | `GETCOUNTCLAIM` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=DATAPEGA.PC_ASM_FW_GCNMFW_WORK | — |
| `GetCurrentDate.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `GETCURRENTDATE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GetLinkStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETLINKSTORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GETTOKENSTORAGE_SQL` |
| `GetTokenStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETTOKENSTORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.GET_TOKEN_STORAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |
| `InsertAttachProdName_Sql.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `INSERTATTACHPRODNAME_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=M_ATTACHMENTPRODUCTNAME | `ASM-FW-GISFW-INT-TREATY_IN / ASM!INSERTATATCHMENT_SQL` |
| `Insert_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERT_T_STORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GENERATEIMAGEID_SQL` |
| `SaveProductNameInwardLIfe.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `SAVEPRODUCTNAMEINWARDLIFE` | [terverifikasi] sqlKind=PLSQL; procs=DBMS_LOB.CREATETEMPORARY,POOLDATA.PEGA_M_PRODUCT_INWARD_LIFE | `ASM-FW-GISFW-INT-PRODUCT_LIFE / ASM!SAVEPRODUCTNAMELIFE` |
| `SaveProductNameLIfe.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `SAVEPRODUCTNAMELIFE` | [terverifikasi] sqlKind=PLSQL; procs=DBMS_LOB.CREATETEMPORARY,POOLDATA.PEGA_M_PRODUCT_LIFE | — |
| `SaveProductNameLIfeFlat.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `SAVEPRODUCTNAMELIFEFLAT` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=POOLDATA.M_PRODUCT_LIFE | `ASM-FW-GISFW-INT-PRODUCT_LIFE / ASM!SAVEPRODUCTNAMELIFE` |
| `SaveTreatyIn.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `SAVETREATYIN` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_TREATY_IN | — |
| `Update_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `UPDATE_T_STORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |

### ReportDefinition — 17 rule (`RULE-OBJ-REPORT-DEFINITION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseBenefitLife_RD.xml` | `ASM-FW-GISFW-INT-BENEFIT_LIFE` | `BROWSEBENEFITLIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseBusinessLife_RD.xml` | `ASM-FW-GISFW-INT-BUSINESS` | `BROWSEBUSINESSLIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-BUSINESS / BROWSEBUSINESS_RD` |
| `BrowseCauseofLossLife_RD.xml` | `ASM-FW-GISFW-INT-CAUSEOFLOSS_LIFE` | `BROWSECAUSEOFLOSSLIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCedingCoLife_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSECEDINGCOLIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-AGENT / BROWSECEDINGCO_RD` |
| `BrowseClientNusaRe_RD.xml` | `ASM-FW-GISFW-INT-CLIENT` | `BROWSECLIENTNUSARE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCurrencyFacIn_RD.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `BROWSECURRENCYFACIN_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCurrencyLIFE_RD.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `BROWSECURRENCYLIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-CURRENCY / BROWSECURRENCYFACIN_RD` |
| `BrowseProductInward.xml` | `ASM-FW-GISFW-INT-PRODUCTINWARD_LIFE` | `BROWSEPRODUCTINWARD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseProductTypeLife_RD.xml` | `ASM-FW-GISFW-INT-PRODUCT_TYPE_LIFE` | `BROWSEPRODUCTTYPELIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseProduct_Life.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `BROWSEPRODUCT_LIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseRICommSummary.xml` | `ASM-FW-GISFW-INT-RICOMM_LIFE_SUMMARY` | `BROWSERICOMMSUMMARY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseRIRiskSummary.xml` | `ASM-FW-GISFW-INT-RIRISK_LIFE_SUMMARY` | `BROWSERIRISKSUMMARY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseRateLifeSummary.xml` | `ASM-FW-GISFW-INT-RATE_LIFE_SUMMARY` | `BROWSERATELIFESUMMARY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseRateLife_RD.xml` | `ASM-FW-GISFW-INT-M_RATE_LIFE` | `BROWSERATELIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTREATY_IN.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `BROWSETREATY_IN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyYear_Life_RD.xml` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` | `BROWSETREATYYEAR_LIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyYear_RD.xml` | `ASM-FW-GISFW-INT-TREATYYEAR` | `BROWSETREATYYEAR_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Section — 12 rule (`RULE-HTML-SECTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `CauseOfLoss_Section.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `CAUSEOFLOSS_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PRODUCT_LIFE / RIRISK_SECTION` |
| `Ceding_Section.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `CEDING_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-LIFE / CEDING_SECTION` |
| `Currency_Section.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `CURRENCY_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PRODUCT_LIFE / POLICYHOLDER_SECTION` |
| `EditProductName_Confirm.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `EDITPRODUCTNAME_CONFIRM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InboxProductName.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `INBOXPRODUCTNAME` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InwardProductName.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `INWARDPRODUCTNAME` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PolicyHolder_Section.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `POLICYHOLDER_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PRODUCT_LIFE / RIRISK_SECTION` |
| `RIRISK_Section.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `RIRISK_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PRODUCT_LIFE / CEDING_SECTION` |
| `RIRate_Section.xml` | `ASM-FW-GISFW-DATA-PLAN` | `RIRATE_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SOB_Section.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `SOB_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PRODUCT_LIFE / CEDING_SECTION` |
| `SaveProductName_Confirm.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `SAVEPRODUCTNAME_CONFIRM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewRate.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `VIEWRATE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PRODUCT_LIFE / VIEWOUTWARD` |

### SystemSettings — 1 rule (`RULE-ADMIN-SYSTEM-SETTINGS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `LinkService.xml` | `LINKSERVICE` | `LINKSERVICE` | [terverifikasi] pyPurpose=LinkService | — |

### When — 1 rule (`RULE-OBJ-WHEN`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `recordEvent.xml` | `@BASECLASS` | `RECORDEVENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

## 4. Tabrakan nama di dalam modul ini

### 4.1 Nama file dipakai di lebih dari satu tipe rule

| Nama file | Tipe | Class |
| --- | --- | --- |
| `InwardProductName.xml` | Harness, Section | `ASM-FW-GISFW-INT-PRODUCT_LIFE, ASM-FW-GISFW-INT-PRODUCT_LIFE` |
| `EditProductName_Confirm.xml` | FlowAction, Section | `ASM-FW-GISFW-INT-PRODUCT_LIFE, ASM-FW-GISFW-INT-PRODUCT_LIFE` |
| `ViewRate.xml` | FlowAction, Section | `@BASECLASS, ASM-FW-GISFW-INT-PRODUCT_LIFE` |
| `SaveProductName_Confirm.xml` | FlowAction, Section | `ASM-FW-GISFW-INT-PRODUCT_LIFE, ASM-FW-GISFW-INT-PRODUCT_LIFE` |

Total: **4** nama file.

### 4.2 Rule `When` bernama sama di lebih dari satu class

_Tidak ada._

## 5. Keluarga klon (satu asal salinan dipakai >1 rule)

Dibaca dari `<pzOriginalInstanceKey>` yang sama pada beberapa file. Ini **sinyal kandidat
konsolidasi** saat migrasi, bukan bukti bahwa perilakunya identik — kesamaan perilaku baru
dapat dinyatakan setelah isinya dibandingkan di STEP D2.

| Asal salinan (class / nama) | Jumlah | File |
| --- | ---: | --- |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE / CHOOSERIRISK` | 3 | `FlowAction/ChooseCauseOfLoss.xml`, `FlowAction/ChoosePolicyHolder.xml`, `FlowAction/ChooseRIRisk.xml` |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE / SETRIRISK_DT` | 3 | `DataTransform/setCauseOfLoss_DT.xml`, `DataTransform/setPolicyHolder_DT.xml`, `DataTransform/setRIRISK_DT.xml` |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE / COPYPRODUCT` | 2 | `DataTransform/CopyProduct.xml`, `DataTransform/SetViewEdit.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` | 2 | `RDBList/GetTokenStorage_SQL.xml`, `RDBList/Update_T_Storage_SQL.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / INSERTGOOGLESTORAGE_ACT` | 2 | `Activity/GetUrlGoogleStorage_Act.xml`, `Activity/InsertGoogleStorage_Act.xml` |
| `ASM-FW-GISFW-INT-TREATY_IN / ASM!GETATTACHMENT_SQL` | 2 | `RDBList/GetAttachment2_Sql.xml`, `RDBList/GetAttachmentProdName_Sql.xml` |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE / COUNTMAXREASURED_ACT` | 2 | `Activity/CountMaxReasured_Act.xml`, `Activity/CountMaxSumReasured_Act.xml` |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE / ASM!BROWSEUNDERWRITINGLIST` | 2 | `RDBList/BrowseTreatyIn.xml`, `RDBList/BrowseUnderwritingList.xml` |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE / SAVEPRODUCTNAME_ACT` | 2 | `Activity/SaveInwardProductName_Act.xml`, `Activity/SaveProductName_Act.xml` |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE / RIRISK_SECTION` | 2 | `Section/CauseOfLoss_Section.xml`, `Section/PolicyHolder_Section.xml` |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE / ASM!SAVEPRODUCTNAMELIFE` | 3 | `RDBList/SaveProductNameInwardLIfe.xml`, `RDBList/SaveProductNameLIfe.xml`, `RDBList/SaveProductNameLIfeFlat.xml` |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE / INWARDPRODUCTNAME` | 2 | `Harness/InwardProductName.xml`, `Section/InwardProductName.xml` |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE / EDITPRODUCTNAME_CONFIRM` | 2 | `FlowAction/EditProductName_Confirm.xml`, `Section/EditProductName_Confirm.xml` |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE / CEDING_SECTION` | 2 | `Section/RIRISK_Section.xml`, `Section/SOB_Section.xml` |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE / SAVEPRODUCTNAME_CONFIRM` | 2 | `FlowAction/SaveProductName_Confirm.xml`, `Section/SaveProductName_Confirm.xml` |
| `ASM-FW-GISFW-INT-CURRENCY / BROWSECURRENCYFACIN_RD` | 2 | `ReportDefinition/BrowseCurrencyFacIn_RD.xml`, `ReportDefinition/BrowseCurrencyLIFE_RD.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GENERATEIMAGEID_SQL` | 2 | `RDBList/GenerateImageID_SQL.xml`, `RDBList/Insert_T_Storage_SQL.xml` |

## 6. Objek database yang disentuh

Diambil dari `<pyBrowseSQL>` pada rule `RDBList`. **Nama objek apa adanya dari SQL** —
kualifikasi skema tidak konsisten di korpus. **Tipe kolom, constraint, dan relasi tidak
diketahui** (OQ-001).

| Tabel/view (apa adanya di SQL) | Dirujuk oleh |
| --- | ---: |
| `T_STORAGE_IMAGE` | 3 rule |
| `M_ATTACHMENTPRODUCTNAME` | 2 rule |
| `M_ATTACHMENTTREATY_2` | 2 rule |
| `CATEGORY_ATTACH_REAS` | 1 rule |
| `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` | 1 rule |
| `M_PRODUCT_LIFE` | 1 rule |
| `POOLDATA.M_PRODUCT_LIFE` | 1 rule |
| `POOLDATA.M_TREATY_IN` | 1 rule |
| `POOLDATA.M_TREATY_IN_EDM` | 1 rule |
| `POOLDATA.TREATYCONTRACT_LIFE` | 1 rule |
| `POOLDATA.TREATYYEAR_LIFE` | 1 rule |
| `POOLDATA.T_FOLDER_IMAGE` | 1 rule |

### Stored procedure / function yang dipanggil

**Body prosedur ini TIDAK ADA di korpus** — logika di dalamnya tidak dapat direkonstruksi
dari XML. Lihat OQ-002.

| Procedure | Dipanggil oleh |
| --- | --- |
| `DBMS_LOB.CREATETEMPORARY` | `SaveProductNameInwardLIfe.xml`, `SaveProductNameLIfe.xml` |
| `POOLDATA.PEGA_M_PRODUCT_LIFE` | `SaveProductNameLIfe.xml` |
| `POOLDATA.PEGA_TREATY_IN` | `SaveTreatyIn.xml` |
| `POOLDATA.PEGA_M_PRODUCT_INWARD_LIFE` | `SaveProductNameInwardLIfe.xml` |
| `POOLDATA.GET_TOKEN_STORAGE` | `GetTokenStorage_SQL.xml` |

## 7. Integrasi eksternal (`ConnectREST`)

Alamat tujuan dicatat sebagai **konfigurasi**, bukan URL literal. `baseURL=SETTING` berarti
alamat diambil dari rule `SystemSettings` yang disebut di kolom Setting; nilai sebenarnya
**tidak ada di korpus**. Kolom terakhir menandai rule yang justru memuat URL literal —
itu temuan, bukan fakta arsitektur (lihat `_summary.md`).

| File | Class | pyServiceName | Sumber base URL | Setting | URL literal di dalam rule? |
| --- | --- | --- | --- | --- | --- |
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `ServiceGoogle` | `SETTING` | `LinkService!LinkService` | tidak |

