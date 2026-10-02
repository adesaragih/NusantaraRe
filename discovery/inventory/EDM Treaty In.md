# Inventaris Rule — EDM Treaty In

STEP D1, batch 1 (domain treaty inward). Dibangkitkan 2026-09-12 dari korpus READ-ONLY
`D:\XML\RNM_BRD\EDM Treaty In\`.

**Batas step ini:** hanya struktur + metadata. Tidak ada kesimpulan logika bisnis di sini.
Penelusuran perilaku adalah STEP D2 (`../flows/`).

**Identitas rule** diambil dari `<pxInsName>` (lihat `../README.md` §6.4).  
Kolom **Asal salinan** berisi `<pzOriginalInstanceKey>` dan hanya diisi bila **berbeda**
dari identitas — artinya rule ini hasil *Save As* dari rule lain. Tanda `—` berarti sama.
Class `—` berarti `<pxInsName>` tidak berisi bagian class (berlaku pada `DataPage`, yang
beridentitas global, bukan per class).

Perintah audit (dijalankan dari `D:\XML\RNM_BRD\`):

```
find "EDM Treaty In" -type f -name "*.xml" | wc -l
find "EDM Treaty In" -type f -name "*.xml" | awk -F/ '{print $2}' | sort | uniq -c
grep -o "<pxInsName>[^<]*" "EDM Treaty In/<tipe>/<file>.xml" | head -1
grep -o "<pzOriginalInstanceKey>[^<]*" "EDM Treaty In/<tipe>/<file>.xml" | head -1
```

## 1. Ringkasan terukur

- Total file XML: **163** — cocok dengan `../README.md` §5.1
- **Identitas rule unik: 163** dari 163 file → tidak ada identitas ganda di dalam modul ini
- Class Pega unik: **25**
- RDBList menurut jenis SQL: PLSQL=8, QUERY=28 (`PLSQL` = blok anonim `BEGIN…END;`, umumnya memanggil stored procedure)
- Prefix pada kunci RDBList: `ASM`=14, `RNM`=22 (arti prefix **belum terverifikasi** — OQ-008)

| Tipe rule | Jumlah |
| --- | ---: |
| Activity | 66 |
| ConnectREST | 1 |
| DataTransform | 10 |
| DecisionTable | 1 |
| Flow | 1 |
| FlowAction | 6 |
| Harness | 3 |
| RDBList | 36 |
| ReportDefinition | 9 |
| Section | 18 |
| SystemSettings | 1 |
| When | 11 |
| **TOTAL** | **163** |

## 2. Sebaran class Pega

Class yang **tidak** berawalan `ASM-` adalah class bawaan Pega atau class portal — ruang
lingkup migrasinya belum ditentukan (**OQ-009**).

| Class | Jumlah rule | Aplikasi? |
| --- | ---: | --- |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN` | 46 | aplikasi |
| `ASM-FW-GISFW-WORK` | 33 | aplikasi |
| `ASM-FW-GISFW-INT-TREATY_IN_EDM` | 14 | aplikasi |
| `DATA-PORTAL` | 14 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-INT-POLICYJSON` | 13 | aplikasi |
| `ASM-FW-GISFW-WORK-ENDORSEMENTTREATY` | 10 | aplikasi |
| `ASM-FW-GISFW-INT-TREATY_IN` | 5 | aplikasi |
| `@BASECLASS` | 3 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-INT-CURRENCY` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-REINSURANCETYPE` | 3 | aplikasi |
| `ASM-FW-GISFW-DATA-INSTALLMENT` | 2 | aplikasi |
| `ASM-FW-GISFW-DATA-XOLREALISASIDATA` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-POLISTREATYIN` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYGROUP` | 2 | aplikasi |
| `ASM-FW-GISFW-DATA` | 1 | aplikasi |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | 1 | aplikasi |
| `ASM-FW-GISFW-DATA-TREATYINLIMITS` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-AGENT` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-MARKETINGOFFICER` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_LINK_SERVICE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-OFFERJSON` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYOUTDETAIL` | 1 | aplikasi |
| `ASM-FW-GISFW-WORK-NB` | 1 | aplikasi |
| `ASSIGN-WORKLIST` | 1 | **bukan aplikasi** → OQ-009 |
| `LINKSERVICE` | 1 | **bukan aplikasi** → OQ-009 |

Rule di class bukan-aplikasi: **19** dari 163.

## 3. Daftar rule per tipe

### Activity — 66 rule (`RULE-OBJ-ACTIVITY`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `CalculateDifferenceEDM_act.xml` | `ASM-FW-GISFW-INT-TREATY_IN_EDM` | `CALCULATEDIFFERENCEEDM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `CekLimitTreatyAcc_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `CEKLIMITTREATYACC_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CheckDataMkt.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `CHECKDATAMKT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `ASM-FW-GISFW-DATA-POLICYTREATYIN / CHECKDATAMO` |
| `CheckDuplicateOffer.xml` | `DATA-PORTAL` | `CHECKDUPLICATEOFFER` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `CheckNopolisAvailability.xml` | `DATA-PORTAL` | `CHECKNOPOLISAVAILABILITY` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `ConvertHistoryDate.xml` | `DATA-PORTAL` | `CONVERTHISTORYDATE` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CopyGeneralDataEDM_act.xml` | `ASM-FW-GISFW-INT-TREATY_IN_EDM` | `COPYGENERALDATAEDM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `CountNetPremi_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTNETPREMI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `CountOGPONP_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTOGPONP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `CountOverridingCommOgp_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTOVERRIDINGCOMMOGP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CountOverridingCommOnp_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTOVERRIDINGCOMMONP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CountPctInstallment_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTPCTINSTALLMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CountResult1Onp_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRESULT1ONP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CountResult1_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRESULT1_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `CountResult2Ogp_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRESULT2OGP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CountResult2Onp_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRESULT2ONP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CountRiCommOgp_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRICOMMOGP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CountRiCommOnp_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRICOMMONP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `CountSpreading_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTSPREADING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `CreateEDMT.xml` | `DATA-PORTAL` | `CREATEEDMT` | [terverifikasi] pyActivityType=ACTIVITY; steps=20 | — |
| `EDMChooseBusiness_Act.xml` | `ASM-FW-GISFW-INT-TREATY_IN_EDM` | `EDMCHOOSEBUSINESS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | — |
| `EDMTCalculateTreatyDifference.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `EDMTCALCULATETREATYDIFFERENCE` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `FetchTreatyGroupOJK.xml` | `ASM-FW-GISFW-WORK-NB` | `FETCHTREATYGROUPOJK` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `FetchTreatyGroupOldID.xml` | `ASM-FW-GISFW-WORK` | `FETCHTREATYGROUPOLDID` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `FillMasterInstallment.xml` | `ASM-FW-GISFW-INT-TREATY_IN_EDM` | `FILLMASTERINSTALLMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `FillPaymentInstallment.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `FILLPAYMENTINSTALLMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `FillPaymentInstallmentEDMT.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `FILLPAYMENTINSTALLMENTEDMT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `ASM-FW-GISFW-DATA-POLICYTREATYIN / FILLPAYMENTINSTALLMENT` |
| `FillSpreading.xml` | `ASM-FW-GISFW-INT-TREATY_IN_EDM` | `FILLSPREADING` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `GeneratePolicyNoTreatyAddendum_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `GENERATEPOLICYNOTREATYADDENDUM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | `ASM-FW-GISFW-DATA-POLICYTREATYIN / GENERATEPOLICYNOTREATY_ACT` |
| `GeneratePolicyNoTreaty_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `GENERATEPOLICYNOTREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=32 | — |
| `GetLinkService.xml` | `ASM-FW-GISFW-INT-M_LINK_SERVICE` | `GETLINKSERVICE` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `InputPolicyTreatyEDMDetail_NP.xml` | `ASM-FW-GISFW-WORK` | `INPUTPOLICYTREATYEDMDETAIL_NP` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | `ASM-FW-GISFW-WORK / INPUTPOLICYTREATYINDETAIL_NONPROP` |
| `InputPolicyTreatyEDMDetail_NP_AdjPremi.xml` | `ASM-FW-GISFW-WORK` | `INPUTPOLICYTREATYEDMDETAIL_NP_ADJPREMI` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `ASM-FW-GISFW-WORK / INPUTPOLICYTREATYEDMDETAIL_NP` |
| `InputPolicyTreatyInPre_Act.xml` | `ASM-FW-GISFW-WORK` | `INPUTPOLICYTREATYINPRE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `InsertHistoryAkseptasiPega.xml` | `ASM-FW-GISFW-WORK` | `INSERTHISTORYAKSEPTASIPEGA` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `InsertToTreatyOutXOLList.xml` | `ASM-FW-GISFW-WORK` | `INSERTTOTREATYOUTXOLLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=37 | `ASM-FW-GISFW-WORK / INSERTTOTREATYXOLLIST` |
| `InsertToTreatyOutXOLListEDMOldData.xml` | `ASM-FW-GISFW-WORK` | `INSERTTOTREATYOUTXOLLISTEDMOLDDATA` | [terverifikasi] pyActivityType=ACTIVITY; steps=30 | `ASM-FW-GISFW-WORK / INSERTTOTREATYOUTXOLLIST` |
| `InsertToTreatyXOLList.xml` | `ASM-FW-GISFW-WORK` | `INSERTTOTREATYXOLLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=20 | — |
| `InsertToTreatyXOLListEDM.xml` | `ASM-FW-GISFW-WORK` | `INSERTTOTREATYXOLLISTEDM` | [terverifikasi] pyActivityType=ACTIVITY; steps=24 | `ASM-FW-GISFW-WORK / INSERTTOTREATYXOLLIST` |
| `InsertToTreatyXOLListEDMOldData.xml` | `ASM-FW-GISFW-WORK` | `INSERTTOTREATYXOLLISTEDMOLDDATA` | [terverifikasi] pyActivityType=ACTIVITY; steps=23 | `ASM-FW-GISFW-WORK / INSERTTOTREATYXOLLISTEDM` |
| `InsetTreatyInProdAddendum_Act.xml` | `ASM-FW-GISFW-WORK` | `INSETTREATYINPRODADDENDUM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=22 | `ASM-FW-GISFW-WORK / INSETTREATYINPROD_ACT` |
| `ProtectionNonProp_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `PROTECTIONNONPROP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `Protection_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `PROTECTION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=16 | — |
| `RemoveTypeTax_ACT.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `REMOVETYPETAX_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SaveJsonPolisTreatyInEDM_Act.xml` | `ASM-FW-GISFW-WORK` | `SAVEJSONPOLISTREATYINEDM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=22 | `ASM-FW-GISFW-WORK / SAVEJSONPOLISTREATYIN_ACT` |
| `SendEmailWithAttachments.xml` | `@BASECLASS` | `SENDEMAILWITHATTACHMENTS` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetCategoryAttach.xml` | `ASM-FW-GISFW-WORK` | `SETCATEGORYATTACH` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GISFW-WORK / INPUTPOLICYTREATYINPRE_ACT` |
| `SetCurrency_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SETCURRENCY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetDueTo_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SETDUETO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GISFW-DATA-POLICYTREATYIN / COUNTNETPREMI_ACT` |
| `SetEDMAchivementValue.xml` | `ASM-FW-GISFW-WORK` | `SETEDMACHIVEMENTVALUE` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `ASM-FW-GISFW-WORK / SETACHIVEMENTVALUE` |
| `SetEDMTCancel.xml` | `ASM-FW-GISFW-WORK` | `SETEDMTCANCEL` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `DATA-PORTAL / SETEDMTCANCEL` |
| `SetEDMTNoPolis.xml` | `DATA-PORTAL` | `SETEDMTNOPOLIS` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `ASM-FW-GISFW-WORK-ENDORSEMENTTREATY / SETEDMNOPOLIS` |
| `SetPPNPPH.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SETPPNPPH` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SetReinstatementPct.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `SETREINSTATEMENTPCT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SetTreatyInEDM_Act.xml` | `DATA-PORTAL` | `SETTREATYINEDM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `DATA-PORTAL / SETTREATYIN_ACT` |
| `SetTreatyIn_Act.xml` | `DATA-PORTAL` | `SETTREATYIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | — |
| `SetValidateInstallment_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SETVALIDATEINSTALLMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SetValueEDM_Act.xml` | `ASM-FW-GISFW-INT-TREATY_IN_EDM` | `SETVALUEEDM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | — |
| `SetValueOldTax.xml` | `ASM-FW-GISFW-INT-TREATY_IN_EDM` | `SETVALUEOLDTAX` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `TreatyInInputVis.xml` | `DATA-PORTAL` | `TREATYININPUTVIS` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `TreatyLoadMasterJoinEdmChooseBusiness.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `TREATYLOADMASTERJOINEDMCHOOSEBUSINESS` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `TreatyRealizationCheckXOLList.xml` | `ASM-FW-GISFW-WORK` | `TREATYREALIZATIONCHECKXOLLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `ASM-SFAGIS-WORK / TREATYREALIZATIONCHECKXOLLIST` |
| `TreatyRealizationCheckXOLListEDM.xml` | `ASM-FW-GISFW-WORK` | `TREATYREALIZATIONCHECKXOLLISTEDM` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `ASM-FW-GISFW-WORK / TREATYREALIZATIONCHECKXOLLIST` |
| `TreatySetReinstatement.xml` | `DATA-PORTAL` | `TREATYSETREINSTATEMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `TrtEdmCheckPolicyError.xml` | `DATA-PORTAL` | `TRTEDMCHECKPOLICYERROR` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `serviceInsertArasapas_act.xml` | `ASM-FW-GISFW-WORK` | `SERVICEINSERTARASAPAS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=20 | — |

### ConnectREST — 1 rule (`RULE-CONNECT-REST`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `convertJsonNusareToProduction.xml` | `ASM-FW-GISFW-WORK` | `CONVERTJSONNUSARETOPRODUCTION` | [terverifikasi] pyServiceName=convertJsonNusareToProduction; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | — |

### DataTransform — 10 rule (`RULE-OBJ-MODEL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AddToListCommentsPolicyTreatyIn_DT.xml` | `ASM-FW-GISFW-WORK` | `ADDTOLISTCOMMENTSPOLICYTREATYIN_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DeptHeadTreatyInAddendum_PreDT.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTTREATY` | `DEPTHEADTREATYINADDENDUM_PREDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-ENDORSEMENT / DEPTHEADTREATYINADDENDUM_PREDT` |
| `ExpandAllExpandables.xml` | `ASM-FW-GISFW-INT-TREATY_IN_EDM` | `EXPANDALLEXPANDABLES` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InboxPolicyTreatyInAddendum_postDT.xml` | `ASM-FW-GISFW-WORK` | `INBOXPOLICYTREATYINADDENDUM_POSTDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InboxPolicyTreatyIn_UW_postDT.xml` | `ASM-FW-GISFW-WORK` | `INBOXPOLICYTREATYIN_UW_POSTDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputPolicyTreatyInAddendum_preAddDT.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTTREATY` | `INPUTPOLICYTREATYINADDENDUM_PREADDDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-ENDORSEMENT / INPUTPOLICYTREATYINADDENDUM_PREADDDT` |
| `InputPolicyTreatyIn_preAddDT.xml` | `ASM-FW-GISFW-WORK` | `INPUTPOLICYTREATYIN_PREADDDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetInstallmentValue.xml` | `ASM-FW-GISFW-INT-TREATY_IN_EDM` | `SETINSTALLMENTVALUE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetOldMasterNo.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SETOLDMASTERNO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SystemSetOneYear_DT.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SYSTEMSETONEYEAR_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### DecisionTable — 1 rule (`RULE-DECLARE-DECISIONTABLE`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `isApproved.xml` | `ASM-FW-GISFW-WORK` | `ISAPPROVED` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Flow — 1 rule (`RULE-OBJ-FLOW`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `InputAddendumTreatyIn.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTTREATY` | `INPUTADDENDUMTREATYIN` | [terverifikasi] pyStartActivity=Start1 | — |

### FlowAction — 6 rule (`RULE-OBJ-FLOWACTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `DeptHeadTreatyIn_UWAddendum.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTTREATY` | `DEPTHEADTREATYIN_UWADDENDUM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-ENDORSEMENT / DEPTHEADTREATYIN_UWADDENDUM` |
| `DetailPolicyAddPremiDetail.xml` | `ASM-FW-GISFW-DATA-XOLREALISASIDATA` | `DETAILPOLICYADDPREMIDETAIL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InboxPolicyTreatyInAddendum.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTTREATY` | `INBOXPOLICYTREATYINADDENDUM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InstallmentList.xml` | `ASM-FW-GISFW-DATA-INSTALLMENT` | `INSTALLMENTLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PolicyTreatyInDeclineConfirm.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `POLICYTREATYINDECLINECONFIRM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowPolicyNoTreaty.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SHOWPOLICYNOTREATY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Harness — 3 rule (`RULE-HTML-HARNESS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BusinessAndSOBListEDM.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `BUSINESSANDSOBLISTEDM` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-DATA-POLICYTREATYIN / BUSINESSANDSOBLIST` |
| `SFAPortal_Endorsement_Treaty.xml` | `DATA-PORTAL` | `SFAPORTAL_ENDORSEMENT_TREATY` | [terverifikasi] pyInclude=1 | — |
| `TreatyCreateEdm.xml` | `DATA-PORTAL` | `TREATYCREATEEDM` | [terverifikasi] pyInclude=1 | — |

### RDBList — 36 rule (`RULE-CONNECT-SQL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseTreatyIn.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `BROWSETREATYIN` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_TREATY_IN,POOLDATA.M_TREATY_IN_EDM | `ASM-FW-GISFW-INT-PRODUCT_LIFE / ASM!BROWSEUNDERWRITINGLIST` |
| `BrowseTreatyInEDM.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `BROWSETREATYINEDM` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_TREATY_IN_EDM | `ASM-FW-GISFW-DATA-POLICYTREATYIN / RNM!BROWSETREATYIN` |
| `BrowseTreatyInEDM_Int_treaty_in_edm.xml` | `ASM-FW-GISFW-INT-TREATY_IN_EDM` | `BROWSETREATYINEDM` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_TREATY_IN_EDM | `ASM-FW-GISFW-INT-TREATY_IN / ASM!BROWSETREATYINEDM` |
| `BrowseTreatyOutDetailEDM.xml` | `ASM-FW-GISFW-INT-TREATYOUTDETAIL` | `BROWSETREATYOUTDETAILEDM` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_TREATY_OUT | `ASM-FW-GISFW-INT-TREATYOUTDETAIL / RNM!BROWSETREATYOUTDETAIL` |
| `CekSTSKonversiJson.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `CEKSTSKONVERSIJSON` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.JSON_POLIS | — |
| `DeleteDataProduction.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `DELETEDATAPRODUCTION` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_DELETE_ERROR_KONVERSI | `ASM-FW-GISFW-INT-POLICYJSON / ASM!DELETEJSONPOLIS` |
| `FetchNoOfferFromNoPolis.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `FETCHNOOFFERFROMNOPOLIS` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TREATYINPRODUCTION | — |
| `FetchNopolisCount.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTTREATY` | `FETCHNOPOLISCOUNT` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.JSON_POLIS | `ASM-FW-GISFW-WORK-ENDORSEMENTTREATY / RNM!FETCHDUPEADJPREMI` |
| `FetchPolisJsonPolis.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTTREATY` | `FETCHPOLISJSONPOLIS` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.JSON_POLIS | — |
| `FetchTreatyGroupOLDID.xml` | `ASM-FW-GISFW-INT-TREATYGROUP` | `FETCHTREATYGROUPOLDID` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TREATYGROUP | — |
| `GETTanggalClosing_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETTANGGALCLOSING_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TANGGAL_CLOSING | — |
| `GenerateNoEDMTreaty.xml` | `ASM-FW-GISFW-INT-POLISTREATYIN` | `GENERATENOEDMTREATY` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | `ASM-FW-GISFW-INT-POLISTREATYIN / ASM!GENERATENOPOLICYTREATY` |
| `GetCountClaim.xml` | `ASSIGN-WORKLIST` | `GETCOUNTCLAIM` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=DATAPEGA.PC_ASM_FW_GCNMFW_WORK | — |
| `GetCurrency.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `GETCURRENCY` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CURRENCY | `ASM-FW-GISFW-INT-CURRENCY / ASM!UPDATEMASTERCURRENCY` |
| `GetCurrentDate.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `GETCURRENTDATE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GetDataCurrencyByName_SQL.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `GETDATACURRENCYBYNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CURRENCY | — |
| `GetDataTreatyInProd_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETDATATREATYINPROD_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYINPRODUCTION | — |
| `GetKodeProdNonLife_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETKODEPRODNONLIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.KODE_PRODUKSI | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETSEQUENCENUMBER_SQL` |
| `GetOldIDBusiness_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `GETOLDIDBUSINESS_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=BUSINESS | — |
| `GetPolicyNoByCaseId.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETPOLICYNOBYCASEID` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | `ASM-FW-GISFW-INT-POLICYJSON / ASM!BROWSECLIENTIDBYNOPOLICY` |
| `GetReinstypeIDbyName_SQL.xml` | `ASM-FW-GISFW-INT-REINSURANCETYPE` | `GETREINSTYPEIDBYNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=REINSURANCETYPE | — |
| `GetSQLDate.xml` | `ASM-FW-GISFW-WORK` | `GETSQLDATE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GetSequenceNumber_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETSEQUENCENUMBER_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETTANGGALCLOSING_SQL` |
| `GetTglInputCaseId.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETTGLINPUTCASEID` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | `ASM-FW-GISFW-INT-POLICYJSON / ASM!GETPOLICYNOBYCASEID` |
| `INSERTJSON_JSONPOLISMONITORING_FACIN.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTJSON_JSONPOLISMONITORING_FACIN` | [terverifikasi] sqlKind=PLSQL; procs=DBMS_LOB.CREATETEMPORARY,POOLDATA.INSERTJSONPOLISMONITORING | `ASM-FW-GISFW-INT-POLICYJSON / ASM!INSERTJSON_JSONPOLIS_FACIN` |
| `InsertHistoryAkseptasiPega_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTHISTORYAKSEPTASIPEGA_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=HISTORYAKSEPTASIPEGA | — |
| `InsertTreatyInProdEDMT_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `INSERTTREATYINPRODEDMT_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=TREATYINPRODUCTION | `ASM-FW-GISFW-INT-POLICYJSON / ASM!INSERTTREATYINPROD_SQL` |
| `SaveAchievementSQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `SAVEACHIEVEMENTSQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.INSERTUPDATEACHIEVMENT | `ASM-FW-GISFW-INT-POLICYJSON / ASM!INSERTTREATYINPROD_SQL` |
| `SavePolisTreatyInEDM_SQL.xml` | `ASM-FW-GISFW-INT-POLISTREATYIN` | `SAVEPOLISTREATYINEDM_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_POLIS_TREATYIN | `ASM-FW-GISFW-INT-POLISTREATYIN / ASM!SAVEPOLISTREATYIN_SQL` |
| `SaveTreatyIn.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `SAVETREATYIN` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_TREATY_IN | — |
| `SelectProdKe.xml` | `ASM-FW-GISFW-INT-TREATY_IN_EDM` | `SELECTPRODKE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | — |
| `SelectSpreadingTreatyInProduction.xml` | `ASM-FW-GISFW-INT-REINSURANCETYPE` | `SELECTSPREADINGTREATYINPRODUCTION` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=REINSURANCETYPE | — |
| `TreatyInSearchProdKe.xml` | `ASM-FW-GISFW-WORK` | `TREATYINSEARCHPRODKE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.JSON_POLIS | — |
| `TreatyLoadMasterJoinEdmChooseBusiness.xml` | `ASM-FW-GISFW-INT-TREATY_IN_EDM` | `TREATYLOADMASTERJOINEDMCHOOSEBUSINESS` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TREATY_IN,POOLDATA.TREATY_IN_EDM | — |
| `TreatyLoadMasterJoinEdmChooseBusinessRetro.xml` | `ASM-FW-GISFW-INT-TREATY_IN_EDM` | `TREATYLOADMASTERJOINEDMCHOOSEBUSINESSRETRO` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TREATY_OUT2 | `ASM-FW-GISFW-INT-TREATY_IN_EDM / RNM!TREATYLOADMASTERJOINEDMCHOOSEBUSINESS` |
| `UpdateErrorNoteJsonPolisMonitoring.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `UPDATEERRORNOTEJSONPOLISMONITORING` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=POOLDATA.JSON_POLIS_MONITORING | `ASM-FW-GISFW-INT-POLICYJSON / ASM!UPDATEERRORNOTEJSONPOLIS` |

### ReportDefinition — 9 rule (`RULE-OBJ-REPORT-DEFINITION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseClientName_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSECLIENTNAME_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-AGENT / BROWSESEARCHSOBCEDING_RD` |
| `BrowseCurrencyTreatyIn_RD.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `BROWSECURRENCYTREATYIN_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseMarketingOfficer_RD.xml` | `ASM-FW-GISFW-INT-MARKETINGOFFICER` | `BROWSEMARKETINGOFFICER_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseReinsuranceType_RD.xml` | `ASM-FW-GISFW-INT-REINSURANCETYPE` | `BROWSEREINSURANCETYPE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTREATY_IN.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `BROWSETREATY_IN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTREATY_IN_EDM.xml` | `ASM-FW-GISFW-INT-TREATY_IN_EDM` | `BROWSETREATY_IN_EDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyGroup_RD.xml` | `ASM-FW-GISFW-INT-TREATYGROUP` | `BROWSETREATYGROUP_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GetListEdmTreaty.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTTREATY` | `GETLISTEDMTREATY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-ENDORSEMENT / GETLISTEDM` |
| `InboxEDM_RD2.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTTREATY` | `INBOXEDM_RD2` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-ENDORSEMENT / INBOXEDM_RD2` |

### Section — 18 rule (`RULE-HTML-SECTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BusinessAndSOBListEDM.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `BUSINESSANDSOBLISTEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DetailPolicyAddPremiDetail.xml` | `ASM-FW-GISFW-DATA-XOLREALISASIDATA` | `DETAILPOLICYADDPREMIDETAIL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DetailPolicyTreatyInAddGeneral.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILPOLICYTREATYINADDGENERAL` | [terverifikasi] pyInclude=2 | — |
| `DetailPolicyTreatyInAddGeneralEditable.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILPOLICYTREATYINADDGENERALEDITABLE` | [terverifikasi] pyInclude=10 | `ASM-FW-GISFW-DATA-POLICYTREATYIN / DETAILPOLICYTREATYINADDGENERAL` |
| `DetailPolicyTreatyInAddPremi.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILPOLICYTREATYINADDPREMI` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DetailPolicyTreatyInAddendum.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILPOLICYTREATYINADDENDUM` | [terverifikasi] pyInclude=6 | `ASM-FW-GISFW-DATA-POLICYTREATYIN / DETAILPOLICYTREATYINADDENDUM1` |
| `DetailPolicyTreatyInPropNewData.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILPOLICYTREATYINPROPNEWDATA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-POLICYTREATYIN / DETAILPOLICYTREATYINPROPOLDDATA2` |
| `DetailPolicyTreatyInPropNewData2.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILPOLICYTREATYINPROPNEWDATA2` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-POLICYTREATYIN / DETAILPOLICYTREATYINPROPNEWDATA` |
| `DetailPolicyTreatyInPropOldData.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILPOLICYTREATYINPROPOLDDATA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-POLICYTREATYIN / DETAILPOLICYTREATYINADDGENERALEDITABLE` |
| `DetailPolicyTreatyInPropOldData2.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILPOLICYTREATYINPROPOLDDATA2` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-POLICYTREATYIN / DETAILPOLICYTREATYINPROPOLDDATA` |
| `DetailPolicyTreatyInPropValueDifference.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILPOLICYTREATYINPROPVALUEDIFFERENCE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-POLICYTREATYIN / DETAILPOLICYTREATYINPROPNEWDATA2` |
| `GeneralPolicyTreatyInAddendum.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENTTREATY` | `GENERALPOLICYTREATYINADDENDUM` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-WORK-ENDORSEMENT / GENERALPOLICYTREATYINADDENDUM` |
| `InstallmentList.xml` | `ASM-FW-GISFW-DATA-INSTALLMENT` | `INSTALLMENTLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ListSuggestEDM.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `LISTSUGGESTEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PolicyTreatyInDeclineConfirm.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `POLICYTREATYINDECLINECONFIRM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SFAPortal_Endorsement_Treaty.xml` | `DATA-PORTAL` | `SFAPORTAL_ENDORSEMENT_TREATY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowPolicyNoTreaty_SC.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SHOWPOLICYNOTREATY_SC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyCreateEdm.xml` | `DATA-PORTAL` | `TREATYCREATEEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### SystemSettings — 1 rule (`RULE-ADMIN-SYSTEM-SETTINGS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `LinkService.xml` | `LINKSERVICE` | `LINKSERVICE` | [terverifikasi] pyPurpose=LinkService | — |

### When — 11 rule (`RULE-OBJ-WHEN`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `IsClaim.xml` | `ASM-FW-GISFW-DATA` | `ISCLAIM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsFacRetro.xml` | `ASM-FW-GISFW-WORK` | `ISFACRETRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsPEGAPROD.xml` | `@BASECLASS` | `ISPEGAPROD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `WORK- / ISPEGAPROD` |
| `IsSPVCreate.xml` | `ASM-FW-GISFW-WORK` | `ISSPVCREATE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsSPVTreaty1.xml` | `ASM-FW-GISFW-WORK` | `ISSPVTREATY1` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISSPVCREATE` |
| `IsSuccessHitService.xml` | `ASM-FW-GISFW-WORK` | `ISSUCCESSHITSERVICE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsTreaty1.xml` | `ASM-FW-GISFW-WORK` | `ISTREATY1` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISSPVCREATE` |
| `IsTreatyIn.xml` | `ASM-FW-GISFW-WORK` | `ISTREATYIN` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ToTREATYDEPTHEAD.xml` | `ASM-FW-GISFW-WORK` | `TOTREATYDEPTHEAD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / TOMANAGERTEKNIK` |
| `recordEvent.xml` | `@BASECLASS` | `RECORDEVENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

## 4. Tabrakan nama di dalam modul ini

### 4.1 Nama file dipakai di lebih dari satu tipe rule

| Nama file | Tipe | Class |
| --- | --- | --- |
| `PolicyTreatyInDeclineConfirm.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-POLICYTREATYIN, ASM-FW-GISFW-DATA-POLICYTREATYIN` |
| `SFAPortal_Endorsement_Treaty.xml` | Harness, Section | `DATA-PORTAL, DATA-PORTAL` |
| `InstallmentList.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-INSTALLMENT, ASM-FW-GISFW-DATA-INSTALLMENT` |
| `DetailPolicyAddPremiDetail.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-XOLREALISASIDATA, ASM-FW-GISFW-DATA-XOLREALISASIDATA` |
| `TreatyCreateEdm.xml` | Harness, Section | `DATA-PORTAL, DATA-PORTAL` |
| `BusinessAndSOBListEDM.xml` | Harness, Section | `ASM-FW-GISFW-DATA-POLICYTREATYIN, ASM-FW-GISFW-DATA-POLICYTREATYIN` |
| `TreatyLoadMasterJoinEdmChooseBusiness.xml` | Activity, RDBList | `ASM-FW-GISFW-DATA-POLICYTREATYIN, ASM-FW-GISFW-INT-TREATY_IN_EDM` |

Total: **7** nama file.

### 4.2 Rule `When` bernama sama di lebih dari satu class

_Tidak ada._

## 5. Keluarga klon (satu asal salinan dipakai >1 rule)

Dibaca dari `<pzOriginalInstanceKey>` yang sama pada beberapa file. Ini **sinyal kandidat
konsolidasi** saat migrasi, bukan bukti bahwa perilakunya identik — kesamaan perilaku baru
dapat dinyatakan setelah isinya dibandingkan di STEP D2.

| Asal salinan (class / nama) | Jumlah | File |
| --- | ---: | --- |
| `DATA-PORTAL / SETTREATYIN_ACT` | 2 | `Activity/SetTreatyInEDM_Act.xml`, `Activity/SetTreatyIn_Act.xml` |
| `ASM-FW-GISFW-WORK / INPUTPOLICYTREATYINPRE_ACT` | 2 | `Activity/InputPolicyTreatyInPre_Act.xml`, `Activity/SetCategoryAttach.xml` |
| `ASM-FW-GISFW-WORK / INSERTTOTREATYXOLLIST` | 3 | `Activity/InsertToTreatyOutXOLList.xml`, `Activity/InsertToTreatyXOLList.xml`, `Activity/InsertToTreatyXOLListEDM.xml` |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN / GENERATEPOLICYNOTREATY_ACT` | 2 | `Activity/GeneratePolicyNoTreatyAddendum_Act.xml`, `Activity/GeneratePolicyNoTreaty_Act.xml` |
| `ASM-FW-GISFW-DATA-INSTALLMENT / INSTALLMENTLIST` | 2 | `FlowAction/InstallmentList.xml`, `Section/InstallmentList.xml` |
| `ASM-FW-GISFW-INT-POLICYJSON / ASM!INSERTTREATYINPROD_SQL` | 2 | `RDBList/InsertTreatyInProdEDMT_SQL.xml`, `RDBList/SaveAchievementSQL.xml` |
| `DATA-PORTAL / TREATYCREATEEDM` | 2 | `Harness/TreatyCreateEdm.xml`, `Section/TreatyCreateEdm.xml` |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN / COUNTNETPREMI_ACT` | 2 | `Activity/CountNetPremi_act.xml`, `Activity/SetDueTo_act.xml` |
| `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETTANGGALCLOSING_SQL` | 2 | `RDBList/GETTanggalClosing_SQL.xml`, `RDBList/GetSequenceNumber_SQL.xml` |
| `ASM-FW-GISFW-WORK / ISSPVCREATE` | 3 | `When/IsSPVCreate.xml`, `When/IsSPVTreaty1.xml`, `When/IsTreaty1.xml` |
| `ASM-FW-GISFW-INT-TREATY_IN_EDM / RNM!TREATYLOADMASTERJOINEDMCHOOSEBUSINESS` | 2 | `RDBList/TreatyLoadMasterJoinEdmChooseBusiness.xml`, `RDBList/TreatyLoadMasterJoinEdmChooseBusinessRetro.xml` |
| `DATA-PORTAL / SFAPORTAL_ENDORSEMENT_TREATY` | 2 | `Harness/SFAPortal_Endorsement_Treaty.xml`, `Section/SFAPortal_Endorsement_Treaty.xml` |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN / DETAILPOLICYTREATYINADDGENERAL` | 2 | `Section/DetailPolicyTreatyInAddGeneral.xml`, `Section/DetailPolicyTreatyInAddGeneralEditable.xml` |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN / POLICYTREATYINDECLINECONFIRM` | 2 | `FlowAction/PolicyTreatyInDeclineConfirm.xml`, `Section/PolicyTreatyInDeclineConfirm.xml` |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN / FILLPAYMENTINSTALLMENT` | 2 | `Activity/FillPaymentInstallment.xml`, `Activity/FillPaymentInstallmentEDMT.xml` |
| `ASM-FW-GISFW-DATA-XOLREALISASIDATA / DETAILPOLICYADDPREMIDETAIL` | 2 | `FlowAction/DetailPolicyAddPremiDetail.xml`, `Section/DetailPolicyAddPremiDetail.xml` |

## 6. Objek database yang disentuh

Diambil dari `<pyBrowseSQL>` pada rule `RDBList`. **Nama objek apa adanya dari SQL** —
kualifikasi skema tidak konsisten di korpus. **Tipe kolom, constraint, dan relasi tidak
diketahui** (OQ-001).

| Tabel/view (apa adanya di SQL) | Dirujuk oleh |
| --- | ---: |
| `POOLDATA.JSON_POLIS` | 4 rule |
| `JSON_POLIS` | 3 rule |
| `POOLDATA.M_TREATY_IN_EDM` | 3 rule |
| `CURRENCY` | 2 rule |
| `REINSURANCETYPE` | 2 rule |
| `TREATYINPRODUCTION` | 2 rule |
| `BUSINESS` | 1 rule |
| `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` | 1 rule |
| `HISTORYAKSEPTASIPEGA` | 1 rule |
| `POOLDATA.JSON_POLIS_MONITORING` | 1 rule |
| `POOLDATA.KODE_PRODUKSI` | 1 rule |
| `POOLDATA.M_TREATY_IN` | 1 rule |
| `POOLDATA.M_TREATY_OUT` | 1 rule |
| `POOLDATA.TANGGAL_CLOSING` | 1 rule |
| `POOLDATA.TREATYGROUP` | 1 rule |
| `POOLDATA.TREATYINPRODUCTION` | 1 rule |
| `POOLDATA.TREATY_IN` | 1 rule |
| `POOLDATA.TREATY_IN_EDM` | 1 rule |
| `POOLDATA.TREATY_OUT2` | 1 rule |

### Stored procedure / function yang dipanggil

**Body prosedur ini TIDAK ADA di korpus** — logika di dalamnya tidak dapat direkonstruksi
dari XML. Lihat OQ-002.

| Procedure | Dipanggil oleh |
| --- | --- |
| `POOLDATA.INSERTUPDATEACHIEVMENT` | `SaveAchievementSQL.xml` |
| `DBMS_LOB.CREATETEMPORARY` | `INSERTJSON_JSONPOLISMONITORING_FACIN.xml` |
| `POOLDATA.INSERTJSONPOLISMONITORING` | `INSERTJSON_JSONPOLISMONITORING_FACIN.xml` |
| `POOLDATA.PEGA_DELETE_ERROR_KONVERSI` | `DeleteDataProduction.xml` |
| `POOLDATA.PEGA_JSON_POLIS_TREATYIN` | `SavePolisTreatyInEDM_SQL.xml` |
| `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` | `GetSequenceNumber_SQL.xml` |
| `POOLDATA.PEGA_TREATY_IN` | `SaveTreatyIn.xml` |

## 7. Integrasi eksternal (`ConnectREST`)

Alamat tujuan dicatat sebagai **konfigurasi**, bukan URL literal. `baseURL=SETTING` berarti
alamat diambil dari rule `SystemSettings` yang disebut di kolom Setting; nilai sebenarnya
**tidak ada di korpus**. Kolom terakhir menandai rule yang justru memuat URL literal —
itu temuan, bukan fakta arsitektur (lihat `_summary.md`).

| File | Class | pyServiceName | Sumber base URL | Setting | URL literal di dalam rule? |
| --- | --- | --- | --- | --- | --- |
| `convertJsonNusareToProduction.xml` | `ASM-FW-GISFW-WORK` | `convertJsonNusareToProduction` | `SETTING` | `LinkService!LinkService` | tidak |

