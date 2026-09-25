# Inventaris Rule — PremiumList Life

STEP D1, batch 4 (domain life, master & outward). Dibangkitkan 2026-09-12 dari korpus READ-ONLY
`D:\XML\RNM_BRD\PremiumList Life\`.

**Batas step ini:** hanya struktur + metadata. Tidak ada kesimpulan logika bisnis di sini.
Penelusuran perilaku adalah STEP D2 (`../flows/`).

**Identitas rule** diambil dari `<pxInsName>` (lihat `../README.md` §6.4).  
Kolom **Asal salinan** berisi `<pzOriginalInstanceKey>` dan hanya diisi bila **berbeda**
dari identitas — artinya rule ini hasil *Save As* dari rule lain. Tanda `—` berarti sama.
Class `—` berarti `<pxInsName>` tidak berisi bagian class (berlaku pada `DataPage`, yang
beridentitas global, bukan per class).

Perintah audit (dijalankan dari `D:\XML\RNM_BRD\`):

```
find "PremiumList Life" -type f -name "*.xml" | wc -l
find "PremiumList Life" -type f -name "*.xml" | awk -F/ '{print $2}' | sort | uniq -c
grep -o "<pxInsName>[^<]*" "PremiumList Life/<tipe>/<file>.xml" | head -1
grep -o "<pzOriginalInstanceKey>[^<]*" "PremiumList Life/<tipe>/<file>.xml" | head -1
```

## 1. Ringkasan terukur

- Total file XML: **124** — **TIDAK COCOK** dengan `../README.md` §5.1 yang mencatat 123
  **Selisih ini dijelaskan, bukan kesalahan:** matriks §5.1 menghitung file **per folder tipe**,
  sedangkan `PremiumList Life/InputPolicyHolder.xml` berada langsung di root modul tanpa folder
  tipe. 123 + 1 = 124. Lihat §3 baris "(di luar folder tipe)" dan §8 di bawah.
- **Identitas rule unik: 124** dari 124 file → tidak ada identitas ganda di dalam modul ini
- Class Pega unik: **22**
- RDBList menurut jenis SQL: PLSQL=8, QUERY=17 (`PLSQL` = blok anonim `BEGIN…END;`, umumnya memanggil stored procedure)
- Prefix pada kunci RDBList: `ASM`=19, `RNM`=6 (arti prefix **belum terverifikasi** — OQ-008)

| Tipe rule | Jumlah |
| --- | ---: |
| (di luar folder tipe) | 1 |
| Activity | 37 |
| ConnectREST | 1 |
| DataTransform | 7 |
| DecisionTable | 2 |
| FlowAction | 10 |
| HTMLRule | 1 |
| Harness | 9 |
| RDBList | 25 |
| ReportDefinition | 9 |
| Section | 19 |
| SystemSettings | 1 |
| When | 2 |
| **TOTAL** | **124** |

## 2. Sebaran class Pega

Class yang **tidak** berawalan `ASM-` adalah class bawaan Pega atau class portal — ruang
lingkup migrasinya belum ditentukan (**OQ-009**).

| Class | Jumlah rule | Aplikasi? |
| --- | ---: | --- |
| `ASM-FW-GISFW-WORK-LIFE` | 64 | aplikasi |
| `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | 11 | aplikasi |
| `@BASECLASS` | 6 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-INT-PRODUCT_LIFE` | 6 | aplikasi |
| `ASM-FW-GISFW-INT-OFFERJSON` | 5 | aplikasi |
| `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` | 4 | aplikasi |
| `ASM-FW-GISFW-INT-POLICYJSON` | 4 | aplikasi |
| `ASM-FW-GISFW-WORK` | 4 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYCONTRACT_LIFE` | 3 | aplikasi |
| `DATA-PORTAL` | 3 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-INT-AGENT` | 2 | aplikasi |
| `ASSIGN-WORKLIST` | 2 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-DATA-OFFERFACIN` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-CLIENT` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-MARKETINGOFFICER` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_LINK_SERVICE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_RATE_LIFE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-RI_COMM_LIFE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` | 1 | aplikasi |
| `CODE-PEGA-PDF` | 1 | **bukan aplikasi** → OQ-009 |
| `LINKSERVICE` | 1 | **bukan aplikasi** → OQ-009 |
| `WORK-` | 1 | **bukan aplikasi** → OQ-009 |

Rule di class bukan-aplikasi: **14** dari 124.

## 3. Daftar rule per tipe

### (di luar folder tipe) — 1 rule (`RULE-OBJ-FLOW`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `InputPolicyHolder.xml` | `ASM-FW-GISFW-WORK-LIFE` | `INPUTPOLICYHOLDER` | [terverifikasi] pyStartActivity=Start1 | — |

### Activity — 37 rule (`RULE-OBJ-ACTIVITY`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AddHistorySuggest.xml` | `ASM-FW-GISFW-WORK-LIFE` | `ADDHISTORYSUGGEST` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `AttachRISlipToWork.xml` | `CODE-PEGA-PDF` | `ATTACHRISLIPTOWORK` | [terverifikasi] pyActivityType=ACTIVITY; steps=26 | `CODE-PEGA-PDF / ATTACHTOWORK` |
| `Calculate1_Act.xml` | `ASM-FW-GISFW-WORK-LIFE` | `CALCULATE1_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=24 | — |
| `CreateInputLife.xml` | `DATA-PORTAL` | `CREATEINPUTLIFE` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `ASM-FW-GISFW-INT-POLICYHOLDER / CREATEINPUTLIFE` |
| `GenerateDataDtlLife_act.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `GENERATEDATADTLLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `GenerateDetailPeserta.xml` | `ASM-FW-GISFW-WORK-LIFE` | `GENERATEDETAILPESERTA` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `GeneratePDFOffer.xml` | `ASM-FW-GISFW-WORK-LIFE` | `GENERATEPDFOFFER` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `GetLinkService.xml` | `ASM-FW-GISFW-INT-M_LINK_SERVICE` | `GETLINKSERVICE` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetOfferLife_Act.xml` | `ASM-FW-GISFW-WORK-LIFE` | `GETOFFERLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `GetPLNumber_Act.xml` | `ASM-FW-GISFW-WORK-LIFE` | `GETPLNUMBER_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `HTMLRISlipToPDF.xml` | `@BASECLASS` | `HTMLRISLIPTOPDF` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `@BASECLASS / HTMLTOPDF` |
| `InputOfferLife_ACT.xml` | `ASM-FW-GISFW-WORK-LIFE` | `INPUTOFFERLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `InputOfferLife_preAct.xml` | `ASM-FW-GISFW-WORK-LIFE` | `INPUTOFFERLIFE_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `InputParamUploadReas_act.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `INPUTPARAMUPLOADREAS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `InsertJsonPolisLife_Act.xml` | `ASM-FW-GISFW-WORK-LIFE` | `INSERTJSONPOLISLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=21 | `ASM-FW-GISFW-WORK-LIFE / INSERTJSONPOLIS_ACT` |
| `InsertLogServiceProd.xml` | `ASM-FW-GISFW-WORK` | `INSERTLOGSERVICEPROD` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GCNMFW-WORK / INSERTLOGSERVICECLAIM` |
| `LoadAttachmentData.xml` | `WORK-` | `LOADATTACHMENTDATA` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `ProtectAccept.xml` | `ASM-FW-GISFW-WORK-LIFE` | `PROTECTACCEPT` | [terverifikasi] pyActivityType=ACTIVITY; steps=24 | — |
| `ProtectProductName_Act.xml` | `ASM-FW-GISFW-WORK-LIFE` | `PROTECTPRODUCTNAME_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SavePremiumList_Act.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SAVEPREMIUMLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=23 | `ASM-FW-GISFW-WORK-LIFE / SUBMITPREMIUMLIST_ACT` |
| `SearchPolicyHolder_act.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SEARCHPOLICYHOLDER_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GISFW-WORK-LIFE / INPUTPOLICYHOLDER_PREACT` |
| `SendEmailNotification.xml` | `@BASECLASS` | `SENDEMAILNOTIFICATION` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetCategoryAttach.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SETCATEGORYATTACH` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SetMaxTBCLife_Act.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SETMAXTBCLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GISFW-WORK-LIFE / INPUTOFFERLIFE_ACT` |
| `SetProdNametoPolis.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `SETPRODNAMETOPOLIS` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetProductInwardLife_Act.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SETPRODUCTINWARDLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SubmitPremiumList_Act.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SUBMITPREMIUMLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `UploadCSVLifePremium_Act.xml` | `@BASECLASS` | `UPLOADCSVLIFEPREMIUM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | — |
| `ValidasiUploadPL_act.xml` | `ASM-FW-GISFW-WORK-LIFE` | `VALIDASIUPLOADPL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=102 | — |
| `WPCLife_Act.xml` | `ASM-FW-GISFW-WORK-LIFE` | `WPCLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `ASM-FW-GISFW-WORK-LIFE / INSERTJSONPOLISLIFE_ACT` |
| `countCategoryAttachment_act.xml` | `ASM-FW-GISFW-WORK-LIFE` | `COUNTCATEGORYATTACHMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `serviceInsertArasapasLife_act.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SERVICEINSERTARASAPASLIFE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `ASM-FW-GISFW-WORK-LIFE / SERVICEINSERTARASAPAS_ACT` |
| `setCeding_act.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SETCEDING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `setNoOffer_Act.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SETNOOFFER_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `setPolicyHolder_act.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SETPOLICYHOLDER_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `setSOB_act.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SETSOB_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GISFW-WORK-LIFE / SETCEDING_ACT` |
| `setSecurityReinsurer_act.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SETSECURITYREINSURER_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |

### ConnectREST — 1 rule (`RULE-CONNECT-REST`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `ConvertJsonNusareToProduction.xml` | `ASM-FW-GISFW-WORK` | `CONVERTJSONNUSARETOPRODUCTION` | [terverifikasi] pyServiceName=convertJsonNusareToProduction; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | — |

### DataTransform — 7 rule (`RULE-OBJ-MODEL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AppendCurrencySummary_DT.xml` | `ASM-FW-GISFW-WORK-LIFE` | `APPENDCURRENCYSUMMARY_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetCoBName_Act.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SETCOBNAME_ACT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetReinsuranceType.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SETREINSURANCETYPE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetRetro_DT.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SETRETRO_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / SETRETRO_DT` |
| `SetSecurityReinsurer_DT.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SETSECURITYREINSURER_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-LIFE / SETRETRO_DT` |
| `SetStatusAkseptasi.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SETSTATUSAKSEPTASI` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `setCategoryAttachment_DT.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SETCATEGORYATTACHMENT_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### DecisionTable — 2 rule (`RULE-DECLARE-DECISIONTABLE`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `IsFlagOnGoingPolicy.xml` | `ASM-FW-GISFW-WORK-LIFE` | `ISFLAGONGOINGPOLICY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-LIFE / ISLIFEACCEPTED` |
| `IsLifeAccepted.xml` | `ASM-FW-GISFW-WORK-LIFE` | `ISLIFEACCEPTED` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### FlowAction — 10 rule (`RULE-OBJ-FLOWACTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `ChooseProdName.xml` | `ASM-FW-GISFW-WORK-LIFE` | `CHOOSEPRODNAME` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ChooseRetroName.xml` | `ASM-FW-GISFW-WORK-LIFE` | `CHOOSERETRONAME` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PRODUCT_LIFE / CHOOSECEDING` |
| `ChooseSecurityReinsurer.xml` | `ASM-FW-GISFW-WORK-LIFE` | `CHOOSESECURITYREINSURER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-LIFE / CHOOSERETRONAME` |
| `InputDataOfferLife.xml` | `ASM-FW-GISFW-WORK-LIFE` | `INPUTDATAOFFERLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PL_DetailAction.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `PL_DETAILACTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `RetroLife.xml` | `ASM-FW-GISFW-INT-TREATYCONTRACT_LIFE` | `RETROLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowLifePremiumDetail.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SHOWLIFEPREMIUMDETAIL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowLifePremiumSummary.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SHOWLIFEPREMIUMSUMMARY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowPLNumber.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SHOWPLNUMBER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `UploadCSV_LifePremium.xml` | `@BASECLASS` | `UPLOADCSV_LIFEPREMIUM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / UPLOADCSV_LIFEPREMIUM` |

### HTMLRule — 1 rule (`RULE-OBJ-HTML`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `GeneratePdfOfferLife.xml` | `ASM-FW-GISFW-WORK-LIFE` | `GENERATEPDFOFFERLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Harness — 9 rule (`RULE-HTML-HARNESS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `Ceding_Harness.xml` | `ASM-FW-GISFW-WORK-LIFE` | `CEDING_HARNESS` | [terverifikasi] pyInclude=1 | — |
| `FollowingOffer_Harness.xml` | `ASM-FW-GISFW-WORK-LIFE` | `FOLLOWINGOFFER_HARNESS` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-WORK-LIFE / POLICYHOLDER_HARNESS` |
| `PolicyHolder_Harness.xml` | `ASM-FW-GISFW-WORK-LIFE` | `POLICYHOLDER_HARNESS` | [terverifikasi] pyInclude=1 | — |
| `PremiumLife_harness.xml` | `DATA-PORTAL` | `PREMIUMLIFE_HARNESS` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-WORK-LIFE / PREMIUMLIFE_HARNESS` |
| `SOB_Harness.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SOB_HARNESS` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-WORK-LIFE / CEDING_HARNESS` |
| `ViewCSVResult_LifePremium.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `VIEWCSVRESULT_LIFEPREMIUM` | [terverifikasi] pyInclude=1 | `RNM-FW-LIFEFW-INT-LIFE_PREMIUM_DETAIL / VIEWCSVRESULT_LIFEPREMIUM` |
| `ViewCSVResult_LifePremium_QP.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `VIEWCSVRESULT_LIFEPREMIUM_QP` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / VIEWCSVRESULT_LIFEPREMIUM` |
| `ViewCSVResult_LifePremium_TP.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `VIEWCSVRESULT_LIFEPREMIUM_TP` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / VIEWCSVRESULT_LIFEPREMIUM` |
| `ViewCSVResult_LifePremium_TR.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `VIEWCSVRESULT_LIFEPREMIUM_TR` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / VIEWCSVRESULT_LIFEPREMIUM_TP` |

### RDBList — 25 rule (`RULE-CONNECT-SQL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AttachmentLife.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `ATTACHMENTLIFE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CATEGORY_ATTACH_REAS | — |
| `CategoryAttach_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `CATEGORYATTACH_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CATEGORY_ATTACH_REAS | — |
| `CekDoubleInsured.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `CEKDOUBLEINSURED` | [terverifikasi] sqlKind=PLSQL; sqlOps=SELECT; tables=V_COUNT,POOLDATA.M_TEMPUPLOADLIFE,V_MIN,V_TEMPRETENSICEDING | — |
| `DeleteTempUploadDataLife.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `DELETETEMPUPLOADDATALIFE` | [terverifikasi] sqlKind=QUERY; sqlOps=DELETE; tables=M_TEMPUPLOADLIFE | — |
| `GETTanggalClosing_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETTANGGALCLOSING_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TANGGAL_CLOSING | — |
| `GetIdOffer_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `GETIDOFFER_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.JSON_OFFER_LIFE | `ASM-FW-GISFW-INT-OFFERJSON / ASM!SAVEOFFERJSONLIFE_SQL` |
| `GetJsonProductLife.xml` | `ASM-FW-GISFW-INT-TREATYCONTRACT_LIFE` | `GETJSONPRODUCTLIFE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TREATYCONTRACT_LIFE | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE / ASM!GETJSONPRODUCTLIFE` |
| `GetKodeProdLife_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETKODEPRODLIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.KODE_PRODUKSI | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETKODEPRODNONLIFE_SQL` |
| `GetNopolisByIDPega.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` | `GETNOPOLISBYIDPEGA` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | — |
| `GetOfferLife_sql.xml` | `ASM-FW-GISFW-WORK-LIFE` | `GETOFFERLIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT,UPDATE; tables=POOLDATA.JSON_OFFER_LIFE | `ASM-FW-GISFW-WORK-LIFE / RNM!GETPLANDNOPOLIS_SQL` |
| `GetPLandNopolis_sql.xml` | `ASM-FW-GISFW-WORK-LIFE` | `GETPLANDNOPOLIS_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.JSON_POLIS | — |
| `GetPolicyNoByCaseId.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETPOLICYNOBYCASEID` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | `ASM-FW-GISFW-INT-POLICYJSON / ASM!BROWSECLIENTIDBYNOPOLICY` |
| `GetProductDtlPL.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `GETPRODUCTDTLPL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=PRODUCT_LIFE | — |
| `GetRateLifePM.xml` | `ASM-FW-GISFW-INT-M_RATE_LIFE` | `GETRATELIFEPM` | [terverifikasi] sqlKind=PLSQL; sqlOps=SELECT; tables=VCOUNT,RATE_LIFE | — |
| `GetRateProductLife.xml` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` | `GETRATEPRODUCTLIFE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.PRODUCTINWARD_LIFE | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE / ASM!GETJSONPRODUCTLIFE` |
| `GetSequenceNumber_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETSEQUENCENUMBER_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETTANGGALCLOSING_SQL` |
| `InsertDataUploadLife.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `INSERTDATAUPLOADLIFE` | [terverifikasi] sqlKind=QUERY; sqlOps=INSERT; tables=POOLDATA.M_TEMPUPLOADLIFE | — |
| `InsertJsonPolis.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` | `INSERTJSONPOLIS` | [terverifikasi] sqlKind=PLSQL; procs=DBMS_LOB.CREATETEMPORARY,POOLDATA.INSERTJSONPOLISLIFE | — |
| `InsertLogServiceProd.xml` | `ASM-FW-GISFW-WORK` | `INSERTLOGSERVICEPROD` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.MONITORING_PROD_LOG | `ASM-FW-GCNMFW-WORK / RNM!INSERTLOGSERVICECLAIM` |
| `InsertPLSummary.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` | `INSERTPLSUMMARY` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY | — |
| `SaveLifeinProduction_SQL.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SAVELIFEINPRODUCTION_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.LIFEINPRODUCTION | `ASM-FW-GISFW-INT-OFFERJSON / ASM!SAVEOFFERJSONLIFE_SQL` |
| `SaveOfferJsonLife_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `SAVEOFFERJSONLIFE_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.INSERTJSONOFFERLIFE; sqlOps=UPDATE | `ASM-FW-GISFW-INT-OFFERJSON / ASM!SAVEOFFERJSON_SQL` |
| `SelectComm_SQL.xml` | `ASM-FW-GISFW-INT-RI_COMM_LIFE` | `SELECTCOMM_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=RICOMM_LIFE | `ASM-FW-GISFW-INT-RI_COMM_LIFE / ASM!CHECKDOUBLERATELIFE_SQL` |
| `SetWPCFunction_Life.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SETWPCFUNCTION_LIFE` | [terverifikasi] sqlKind=QUERY; procs=POOLDATA.GETQUARTER; sqlOps=SELECT | `ASM-FW-GISFW-WORK-LIFE / ASM!BROWSEPREMIUMUPLOAD` |
| `SetWPCFunction_LifeRetro.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SETWPCFUNCTION_LIFERETRO` | [terverifikasi] sqlKind=QUERY; procs=POOLDATA.GETQUARTERRETRO; sqlOps=SELECT | `ASM-FW-GISFW-WORK-LIFE / ASM!SETWPCFUNCTION_LIFE` |

### ReportDefinition — 9 rule (`RULE-OBJ-REPORT-DEFINITION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseAgent_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSEAGENT_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-AGENT / BROWSEAGENTNUSARE_RD` |
| `BrowseCedingCoLife_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSECEDINGCOLIFE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-AGENT / BROWSECEDINGCO_RD` |
| `BrowseClientNusaRe_RD.xml` | `ASM-FW-GISFW-INT-CLIENT` | `BROWSECLIENTNUSARE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseMarketingOfficer_RD.xml` | `ASM-FW-GISFW-INT-MARKETINGOFFICER` | `BROWSEMARKETINGOFFICER_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowsePremiumList_RD.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` | `BROWSEPREMIUMLIST_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseProductForNB_Life.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` | `BROWSEPRODUCTFORNB_LIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PRODUCT_LIFE / BROWSEPRODUCT_LIFE` |
| `GetPLNumber_rd.xml` | `ASM-FW-GISFW-WORK-LIFE` | `GETPLNUMBER_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-LIFE / INBOXPREMIUMLIST` |
| `InboxPremiumList.xml` | `ASSIGN-WORKLIST` | `INBOXPREMIUMLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InboxPremiumListOLD.xml` | `ASSIGN-WORKLIST` | `INBOXPREMIUMLISTOLD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASSIGN-WORKLIST / INBOXPREMIUMLIST` |

### Section — 19 rule (`RULE-HTML-SECTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `Ceding_Section.xml` | `ASM-FW-GISFW-WORK-LIFE` | `CEDING_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-LIFE / POLICYHOLDER_SECTION` |
| `ChooseProdName.xml` | `ASM-FW-GISFW-WORK-LIFE` | `CHOOSEPRODNAME` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ConfirmSection.xml` | `ASM-FW-GISFW-WORK-LIFE` | `CONFIRMSECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `FollowingOffer_Section.xml` | `ASM-FW-GISFW-WORK-LIFE` | `FOLLOWINGOFFER_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-LIFE / POLICYHOLDER_SECTION` |
| `InputOfferLife.xml` | `ASM-FW-GISFW-WORK-LIFE` | `INPUTOFFERLIFE` | [terverifikasi] pyInclude=1 | — |
| `PL_Detail_Sec.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `PL_DETAIL_SEC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PolicyHolder_Section.xml` | `ASM-FW-GISFW-WORK-LIFE` | `POLICYHOLDER_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-QUOTATION / CEDINGCOHIERARKI` |
| `PremiumList.xml` | `DATA-PORTAL` | `PREMIUMLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-LIFE / PREMIUMLIST` |
| `RetroDetailLife.xml` | `ASM-FW-GISFW-INT-TREATYCONTRACT_LIFE` | `RETRODETAILLIFE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Retro_Section.xml` | `ASM-FW-GISFW-WORK-LIFE` | `RETRO_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-LIFE / CEDING_SECTION` |
| `SOB_Section.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SOB_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-LIFE / CEDING_SECTION` |
| `SecurityReinsurer_Section.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SECURITYREINSURER_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-LIFE / RETRO_SECTION` |
| `ShowLifePremiumDetail.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SHOWLIFEPREMIUMDETAIL` | [terverifikasi] pyInclude=2 | — |
| `ShowLifePremiumSummary.xml` | `ASM-FW-GISFW-WORK-LIFE` | `SHOWLIFEPREMIUMSUMMARY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowPolis_sc.xml` | `ASM-FW-GISFW-WORK` | `SHOWPOLIS_SC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewCSVResult_LifePremium.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `VIEWCSVRESULT_LIFEPREMIUM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewCSVResult_LifePremium_QP.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `VIEWCSVRESULT_LIFEPREMIUM_QP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewCSVResult_LifePremium_TP.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `VIEWCSVRESULT_LIFEPREMIUM_TP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewCSVResult_LifePremium_TR.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` | `VIEWCSVRESULT_LIFEPREMIUM_TR` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

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
| `ShowLifePremiumSummary.xml` | FlowAction, Section | `ASM-FW-GISFW-WORK-LIFE, ASM-FW-GISFW-WORK-LIFE` |
| `ViewCSVResult_LifePremium_TP.xml` | Harness, Section | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL, ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` |
| `ShowLifePremiumDetail.xml` | FlowAction, Section | `ASM-FW-GISFW-WORK-LIFE, ASM-FW-GISFW-WORK-LIFE` |
| `InsertLogServiceProd.xml` | Activity, RDBList | `ASM-FW-GISFW-WORK, ASM-FW-GISFW-WORK` |
| `ViewCSVResult_LifePremium.xml` | Harness, Section | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL, ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` |
| `ViewCSVResult_LifePremium_TR.xml` | Harness, Section | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL, ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` |
| `ViewCSVResult_LifePremium_QP.xml` | Harness, Section | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL, ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` |
| `ChooseProdName.xml` | FlowAction, Section | `ASM-FW-GISFW-WORK-LIFE, ASM-FW-GISFW-WORK-LIFE` |

Total: **8** nama file.

### 4.2 Rule `When` bernama sama di lebih dari satu class

_Tidak ada._

## 5. Keluarga klon (satu asal salinan dipakai >1 rule)

Dibaca dari `<pzOriginalInstanceKey>` yang sama pada beberapa file. Ini **sinyal kandidat
konsolidasi** saat migrasi, bukan bukti bahwa perilakunya identik — kesamaan perilaku baru
dapat dinyatakan setelah isinya dibandingkan di STEP D2.

| Asal salinan (class / nama) | Jumlah | File |
| --- | ---: | --- |
| `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETTANGGALCLOSING_SQL` | 2 | `RDBList/GETTanggalClosing_SQL.xml`, `RDBList/GetSequenceNumber_SQL.xml` |
| `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / VIEWCSVRESULT_LIFEPREMIUM_TP` | 2 | `Harness/ViewCSVResult_LifePremium_TR.xml`, `Section/ViewCSVResult_LifePremium_TP.xml` |
| `ASM-FW-GISFW-WORK-LIFE / CHOOSEPRODNAME` | 2 | `FlowAction/ChooseProdName.xml`, `Section/ChooseProdName.xml` |
| `ASM-FW-GISFW-WORK-LIFE / POLICYHOLDER_HARNESS` | 2 | `Harness/FollowingOffer_Harness.xml`, `Harness/PolicyHolder_Harness.xml` |
| `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL / VIEWCSVRESULT_LIFEPREMIUM` | 3 | `Harness/ViewCSVResult_LifePremium_QP.xml`, `Harness/ViewCSVResult_LifePremium_TP.xml`, `Section/ViewCSVResult_LifePremium.xml` |
| `ASM-FW-GISFW-WORK-LIFE / CEDING_HARNESS` | 2 | `Harness/Ceding_Harness.xml`, `Harness/SOB_Harness.xml` |
| `ASM-FW-GISFW-WORK-LIFE / ISLIFEACCEPTED` | 2 | `DecisionTable/IsFlagOnGoingPolicy.xml`, `DecisionTable/IsLifeAccepted.xml` |
| `ASM-FW-GISFW-WORK-LIFE / SETCEDING_ACT` | 2 | `Activity/setCeding_act.xml`, `Activity/setSOB_act.xml` |
| `ASM-FW-GISFW-WORK-LIFE / SUBMITPREMIUMLIST_ACT` | 2 | `Activity/SavePremiumList_Act.xml`, `Activity/SubmitPremiumList_Act.xml` |
| `ASM-FW-GISFW-WORK-LIFE / INPUTOFFERLIFE_ACT` | 2 | `Activity/InputOfferLife_ACT.xml`, `Activity/SetMaxTBCLife_Act.xml` |
| `ASM-FW-GISFW-INT-TREATYYEAR_LIFE / ASM!GETJSONPRODUCTLIFE` | 2 | `RDBList/GetJsonProductLife.xml`, `RDBList/GetRateProductLife.xml` |
| `ASM-FW-GISFW-INT-OFFERJSON / ASM!SAVEOFFERJSONLIFE_SQL` | 2 | `RDBList/GetIdOffer_SQL.xml`, `RDBList/SaveLifeinProduction_SQL.xml` |
| `ASM-FW-GISFW-WORK-LIFE / SHOWLIFEPREMIUMDETAIL` | 2 | `FlowAction/ShowLifePremiumDetail.xml`, `Section/ShowLifePremiumDetail.xml` |
| `ASM-FW-GISFW-WORK-LIFE / POLICYHOLDER_SECTION` | 2 | `Section/Ceding_Section.xml`, `Section/FollowingOffer_Section.xml` |
| `ASM-FW-GISFW-WORK-LIFE / RNM!GETPLANDNOPOLIS_SQL` | 2 | `RDBList/GetOfferLife_sql.xml`, `RDBList/GetPLandNopolis_sql.xml` |
| `ASM-FW-GISFW-WORK-LIFE / SHOWLIFEPREMIUMSUMMARY` | 2 | `FlowAction/ShowLifePremiumSummary.xml`, `Section/ShowLifePremiumSummary.xml` |
| `ASM-FW-GISFW-WORK-LIFE / CEDING_SECTION` | 2 | `Section/Retro_Section.xml`, `Section/SOB_Section.xml` |
| `ASSIGN-WORKLIST / INBOXPREMIUMLIST` | 2 | `ReportDefinition/InboxPremiumList.xml`, `ReportDefinition/InboxPremiumListOLD.xml` |

## 6. Objek database yang disentuh

Diambil dari `<pyBrowseSQL>` pada rule `RDBList`. **Nama objek apa adanya dari SQL** —
kualifikasi skema tidak konsisten di korpus. **Tipe kolom, constraint, dan relasi tidak
diketahui** (OQ-001).

| Tabel/view (apa adanya di SQL) | Dirujuk oleh |
| --- | ---: |
| `CATEGORY_ATTACH_REAS` | 2 rule |
| `JSON_POLIS` | 2 rule |
| `POOLDATA.JSON_OFFER_LIFE` | 2 rule |
| `POOLDATA.M_TEMPUPLOADLIFE` | 2 rule |
| `M_TEMPUPLOADLIFE` | 1 rule |
| `POOLDATA.JSON_POLIS` | 1 rule |
| `POOLDATA.KODE_PRODUKSI` | 1 rule |
| `POOLDATA.LIFEINPRODUCTION` | 1 rule |
| `POOLDATA.MONITORING_PROD_LOG` | 1 rule |
| `POOLDATA.PRODUCTINWARD_LIFE` | 1 rule |
| `POOLDATA.TANGGAL_CLOSING` | 1 rule |
| `POOLDATA.TREATYCONTRACT_LIFE` | 1 rule |
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
| `POOLDATA.GETQUARTERRETRO` | `SetWPCFunction_LifeRetro.xml` |
| `POOLDATA.GETQUARTER` | `SetWPCFunction_Life.xml` |
| `DBMS_LOB.CREATETEMPORARY` | `InsertJsonPolis.xml` |
| `POOLDATA.INSERTJSONOFFERLIFE` | `SaveOfferJsonLife_SQL.xml` |
| `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` | `GetSequenceNumber_SQL.xml` |
| `POOLDATA.INSERTJSONPOLISLIFE` | `InsertJsonPolis.xml` |
| `POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY` | `InsertPLSummary.xml` |

## 7. Integrasi eksternal (`ConnectREST`)

Alamat tujuan dicatat sebagai **konfigurasi**, bukan URL literal. `baseURL=SETTING` berarti
alamat diambil dari rule `SystemSettings` yang disebut di kolom Setting; nilai sebenarnya
**tidak ada di korpus**. Kolom terakhir menandai rule yang justru memuat URL literal —
itu temuan, bukan fakta arsitektur (lihat `_summary.md`).

| File | Class | pyServiceName | Sumber base URL | Setting | URL literal di dalam rule? |
| --- | --- | --- | --- | --- | --- |
| `ConvertJsonNusareToProduction.xml` | `ASM-FW-GISFW-WORK` | `convertJsonNusareToProduction` | `SETTING` | `LinkService!LinkService` | tidak |


## 8. Catatan khusus modul ini — OQ-004 terjawab

`[terverifikasi]` `PremiumList Life/InputPolicyHolder.xml` adalah satu-satunya file XML di seluruh
korpus 9.369 file yang berada **di luar folder tipe rule**. Identitasnya kini diketahui:

| Hal | Nilai |
| --- | --- |
| `<pzOriginalInstanceKey>` | `RULE-OBJ-FLOW ASM-FW-GISFW-WORK-LIFE INPUTPOLICYHOLDER #20180807T032319.287 GMT` |
| `<pxInsName>` | `ASM-FW-GISFW-WORK-LIFE!INPUTPOLICYHOLDER` |
| `<pxObjClass>` | `Rule-Obj-Flow` |
| `<pyStartActivity>` | `Start1` |
| Ukuran | 106.880 byte |

Kedua sumber tipe sepakat: ini rule **`Flow`**.

**Konsekuensi untuk OQ-005.** Daftar "modul tanpa rule `Flow`" yang disusun di STEP D0 memasukkan
`PremiumList Life`, karena folder `Flow/` di modul ini memang tidak ada. Itu **keliru** — modul ini
punya rule `Flow`, hanya filenya tidak diletakkan di folder `Flow/`. Jumlah modul tanpa `Flow`
turun dari 6 menjadi **5**, dan jumlah rule `Flow` di korpus naik dari 23 menjadi **24**.

Untuk STEP D2, berarti `PremiumList Life` **dapat** ditelusur dengan metode standar (mulai dari
`<pyStartActivity>`), tidak perlu titik masuk alternatif.

Perintah audit:

```
grep -o "<pzOriginalInstanceKey>[^<]*" "PremiumList Life/InputPolicyHolder.xml" | head -1
grep -o "<pxObjClass>[^<]*" "PremiumList Life/InputPolicyHolder.xml" | head -1
find . -maxdepth 2 -name "*.xml" -not -path "./OUTPUT_HASIL_RNM/*"
```

Yang **masih belum terverifikasi**: mengapa file ini ditempatkan di luar struktur folder. Itu
pertanyaan proses ekspor, bukan pertanyaan isi rule — tidak memblokir apa pun.
