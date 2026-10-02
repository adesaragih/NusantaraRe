# Inventaris Rule — Claim Fac In

STEP D1, batch 2 (domain facultative inward). Dibangkitkan 2026-09-12 dari korpus READ-ONLY
`D:\XML\RNM_BRD\Claim Fac In\`.

**Batas step ini:** hanya struktur + metadata. Tidak ada kesimpulan logika bisnis di sini.
Penelusuran perilaku adalah STEP D2 (`../flows/`).

**Identitas rule** diambil dari `<pxInsName>` (lihat `../README.md` §6.4).  
Kolom **Asal salinan** berisi `<pzOriginalInstanceKey>` dan hanya diisi bila **berbeda**
dari identitas — artinya rule ini hasil *Save As* dari rule lain. Tanda `—` berarti sama.
Class `—` berarti `<pxInsName>` tidak berisi bagian class (berlaku pada `DataPage`, yang
beridentitas global, bukan per class).

Perintah audit (dijalankan dari `D:\XML\RNM_BRD\`):

```
find "Claim Fac In" -type f -name "*.xml" | wc -l
find "Claim Fac In" -type f -name "*.xml" | awk -F/ '{print $2}' | sort | uniq -c
grep -o "<pxInsName>[^<]*" "Claim Fac In/<tipe>/<file>.xml" | head -1
grep -o "<pzOriginalInstanceKey>[^<]*" "Claim Fac In/<tipe>/<file>.xml" | head -1
```

## 1. Ringkasan terukur

- Total file XML: **482** — cocok dengan `../README.md` §5.1
- **Identitas rule unik: 481** dari 482 file → 1 file adalah identitas ganda **di dalam modul ini**
- Class Pega unik: **55**
- RDBList menurut jenis SQL: PLSQL=10, QUERY=54 (`PLSQL` = blok anonim `BEGIN…END;`, umumnya memanggil stored procedure)
- Prefix pada kunci RDBList: `ASM`=13, `RNM`=24, `GCNM`=27 (arti prefix **belum terverifikasi** — OQ-008)

| Tipe rule | Jumlah |
| --- | ---: |
| Activity | 179 |
| ConnectREST | 8 |
| DataPage | 8 |
| DataTransform | 28 |
| DecisionTable | 1 |
| Flow | 1 |
| FlowAction | 33 |
| Harness | 16 |
| RDBList | 63 |
| ReportDefinition | 27 |
| Section | 57 |
| SystemSettings | 1 |
| When | 60 |
| **TOTAL** | **482** |

## 2. Sebaran class Pega

Class yang **tidak** berawalan `ASM-` adalah class bawaan Pega atau class portal — ruang
lingkup migrasinya belum ditentukan (**OQ-009**).

| Class | Jumlah rule | Aplikasi? |
| --- | ---: | --- |
| `ASM-FW-GCNMFW-WORK-PNC` | 89 | aplikasi |
| `ASM-FW-GCNMFW-DATA-OBJECTITEM` | 56 | aplikasi |
| `ASM-FW-GCNMFW-DATA-OBJECT` | 50 | aplikasi |
| `@BASECLASS` | 44 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | 43 | aplikasi |
| `ASM-FW-GCNMFW-WORK` | 25 | aplikasi |
| `ASM-FW-GISFW-INT-POLICYJSON` | 18 | aplikasi |
| `ASM-FW-GISFW-DATA` | 13 | aplikasi |
| `ASM-FW-GISFW-DATA-QUOTATION` | 10 | aplikasi |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | 10 | aplikasi |
| `ASM-FW-GISFW-WORK` | 9 | aplikasi |
| _(tanpa class)_ | 8 | n/a (DataPage) |
| `CODE-PEGA-LIST` | 8 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GCNMFW-INT-OS_AKSEPTASI_KLAIM` | 7 | aplikasi |
| `ASM-FW-GCNMFW-DATA-CLAIMDATA` | 6 | aplikasi |
| `ASM-FW-GCNMFW-INT-V_POLIS` | 6 | aplikasi |
| `DATA-PORTAL` | 6 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-INT-AGENT` | 5 | aplikasi |
| `ASM-FW-GISFW-INT-RW` | 5 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYGROUP` | 5 | aplikasi |
| `ASM-FW-GISFW-DATA-AGENT` | 4 | aplikasi |
| `ASM-FW-GISFW-DATA-SPREADINGRISK` | 4 | aplikasi |
| `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS` | 3 | aplikasi |
| `ASM-FW-GISFW-DATA-FACOFFER` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-ADJUSTERCONSULTANT` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-CURRENCY` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-OFFERJSON` | 3 | aplikasi |
| `ASM-FW-GCNMFW-DATA-CURRENCY` | 2 | aplikasi |
| `ASM-FW-GCNMFW-INT-CATASTROPHE` | 2 | aplikasi |
| `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | 2 | aplikasi |
| `ASM-FW-GCNMFW-INT-EMAILKOMITE` | 2 | aplikasi |
| `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS_BUSINESS` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-BANKACCOUNT` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-CLIENT` | 2 | aplikasi |
| `ASM-FW-GISFW-WORK-NB` | 2 | aplikasi |
| `ASM-FW-GCNMFW-DATA-ESTIMASI` | 1 | aplikasi |
| `ASM-FW-GCNMFW-INT-CLAIMREJECTED` | 1 | aplikasi |
| `ASM-FW-GCNMFW-INT-V_LST_DOC_TRAVEL_COVERAGE` | 1 | aplikasi |
| `ASM-FW-GCNMFW-INT-V_MST_USER_TEKNIS` | 1 | aplikasi |
| `ASM-FW-GCNMFW-INT-V_M_CAUSE_OF_LOSS` | 1 | aplikasi |
| `ASM-FW-GCNMFW-INT-V_STS_CLAIM` | 1 | aplikasi |
| `ASM-FW-GCNMFW-WORK-OPENPROTECTION` | 1 | aplikasi |
| `ASM-FW-GISFW-DATA-ANEKA` | 1 | aplikasi |
| `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | 1 | aplikasi |
| `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-BUSINESS` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-CLAUSE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-COUNTRY` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-CURRENCYSTANDARD` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-DISTRICT` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_LINK_SERVICE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-REINSURANCETYPE` | 1 | aplikasi |
| `ASM-FW-GISFW-WORK-ENDORSEMENT` | 1 | aplikasi |
| `DATA-PARTY-PERSON` | 1 | **bukan aplikasi** → OQ-009 |
| `LINKSERVICE` | 1 | **bukan aplikasi** → OQ-009 |

Rule di class bukan-aplikasi: **60** dari 482.

## 3. Daftar rule per tipe

### Activity — 179 rule (`RULE-OBJ-ACTIVITY`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AddDelProgress.xml` | `ASM-FW-GCNMFW-DATA-CLAIMDATA` | `ADDDELPROGRESS` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `AgentSourceBiz_Act.xml` | `ASM-FW-GISFW-DATA-AGENT` | `AGENTSOURCEBIZ_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `BackToRegister_act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `BACKTOREGISTER_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `BrowseDocTravel.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `BROWSEDOCTRAVEL` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `CLaimFaceSheet_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `CLAIMFACESHEET_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=82 | — |
| `CNMInsertCauseOfLoss_act.xml` | `@BASECLASS` | `CNMINSERTCAUSEOFLOSS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / CNMINSERTSURVEYORS_ACT` |
| `CNMInsertDetailCauseOfLoss_act.xml` | `@BASECLASS` | `CNMINSERTDETAILCAUSEOFLOSS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / CNMINSERTSTSCLAIM_ACT` |
| `CNMSetDetailCauseOfLoss_act.xml` | `@BASECLASS` | `CNMSETDETAILCAUSEOFLOSS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `@BASECLASS / SETDETAILCAUSEOFLOSS_ACT` |
| `CallActivityInputRegister.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `CALLACTIVITYINPUTREGISTER` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | — |
| `CallSpreadingView.xml` | `ASM-FW-GISFW-DATA-SPREADINGRISK` | `CALLSPREADINGVIEW` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CancelTreatyGroup.xml` | `@BASECLASS` | `CANCELTREATYGROUP` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CekCoverageNote_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `CEKCOVERAGENOTE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `CekExGratia.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `CEKEXGRATIA` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `CekPremiLunas_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `CEKPREMILUNAS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `CheckAnyAcceptation.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `CHECKANYACCEPTATION` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `CheckCurrency_ACT.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `CHECKCURRENCY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | — |
| `CheckDateReceived_Act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `CHECKDATERECEIVED_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `ASM-FW-GCNMFW-WORK-PNC / CHECKDATE_ACT` |
| `CheckDateReport_Act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `CHECKDATEREPORT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `ASM-FW-GCNMFW-WORK-PNC / CHECKDATERECEIVED_ACT` |
| `CheckDate_Act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `CHECKDATE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `CheckDoubleClaim_Act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `CHECKDOUBLECLAIM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `ASM-FW-GCNMFW-WORK-PNC / CHECKDATE_ACT` |
| `CheckEndorsment.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `CHECKENDORSMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CheckEstimateValue.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `CHECKESTIMATEVALUE` | [terverifikasi] pyActivityType=ACTIVITY; steps=48 | — |
| `CheckIDXAneka_act.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `CHECKIDXANEKA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CheckLimitSpreadingTreaty_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `CHECKLIMITSPREADINGTREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=56 | `ASM-FW-GCNMFW-DATA-OBJECTITEM / CHECKLIMIT_ACT` |
| `CheckLimit_Act1.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `CHECKLIMIT_ACT1` | [terverifikasi] pyActivityType=ACTIVITY; steps=21 | `ASM-FW-GCNMFW-DATA-OBJECTITEM / CHECKLIMIT_ACT` |
| `CheckListEstimasi_Act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `CHECKLISTESTIMASI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CheckPeriodePolicy.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `CHECKPERIODEPOLICY` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `ASM-FW-GCNMFW-WORK-PNC / CHECKVIEWPOLISRNM_ACT` |
| `CheckTotalSpreadingPct_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `CHECKTOTALSPREADINGPCT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `ASM-FW-GCNMFW-DATA-OBJECTITEM / SETTOTALSPREADING` |
| `ChooseDla_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `CHOOSEDLA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=16 | — |
| `CloseClaim.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `CLOSECLAIM` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | — |
| `ConvertTSINusare_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `CONVERTTSINUSARE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `CopyNB_Act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `COPYNB_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | — |
| `CopySpreading_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `COPYSPREADING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=18 | — |
| `CountSpreadingAdjustment_ACT.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `COUNTSPREADINGADJUSTMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `CountSpreadingClaim_ACT.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `COUNTSPREADINGCLAIM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=27 | `ASM-FW-GCNMFW-DATA-OBJECTITEM / CHECKESTIMATEVALUE` |
| `CountTSI_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `COUNTTSI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `CountTotalEstimasi_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `COUNTTOTALESTIMASI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=35 | `ASM-FW-GCNMFW-DATA-OBJECTITEM / SETESTIMATIONVALUE_ACT` |
| `CreateKMTNo_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `CREATEKMTNO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=33 | — |
| `CreateRemarksNotePLA_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `CREATEREMARKSNOTEPLA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `CreateRemarksNote_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `CREATEREMARKSNOTE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `DLAFacintoTreaty_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `DLAFACINTOTREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=66 | `ASM-FW-GCNMFW-DATA-OBJECT / GENERATEDLATREATY_ACT1` |
| `DeleteDataTreatyGroup.xml` | `@BASECLASS` | `DELETEDATATREATYGROUP` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `DeleteLocation_Act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `DELETELOCATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `DeleteObjectItemMBU.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `DELETEOBJECTITEMMBU` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `DeleteTest.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `DELETETEST` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `DeleteValueEstimation.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `DELETEVALUEESTIMATION` | [terverifikasi] pyActivityType=ACTIVITY; steps=31 | — |
| `DisableSendComite_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `DISABLESENDCOMITE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `DownloadDocumentClaim.xml` | `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | `DOWNLOADDOCUMENTCLAIM` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `WORK- / DOWNLOADDOCUMENTCLAIM` |
| `DownloadDocumentKasir.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `DOWNLOADDOCUMENTKASIR` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `DraftGenerateDLAFacin_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `DRAFTGENERATEDLAFACIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=62 | `ASM-FW-GCNMFW-DATA-OBJECT / GENERATEDLAFACIN_ACT` |
| `EditMstConsultant_Act.xml` | `ASM-FW-GISFW-INT-ADJUSTERCONSULTANT` | `EDITMSTCONSULTANT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `FilterCoverage_Act.xml` | `CODE-PEGA-LIST` | `FILTERCOVERAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `CODE-PEGA-LIST / FILTERPROPERTYITEM_ACT` |
| `FilterPropertyItem_Act.xml` | `CODE-PEGA-LIST` | `FILTERPROPERTYITEM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `GenerateDLAFacin_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `GENERATEDLAFACIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=70 | `ASM-FW-GCNMFW-DATA-OBJECT / GENERATEDLAFACIN` |
| `GeneratePLA.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `GENERATEPLA` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | — |
| `GeneratePLATreaty_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `GENERATEPLATREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=64 | `ASM-FW-GCNMFW-DATA-OBJECTITEM / GENERATEPLA` |
| `GetAllData_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `GETALLDATA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=39 | — |
| `GetAnekaList.xml` | `CODE-PEGA-LIST` | `GETANEKALIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `GetBase64Attachment.xml` | `ASM-FW-GCNMFW-WORK` | `GETBASE64ATTACHMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | `WORK- / LOADATTACHMENTDATA` |
| `GetCeding_act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `GETCEDING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `GetCoverageAnekaList.xml` | `CODE-PEGA-LIST` | `GETCOVERAGEANEKALIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `CODE-PEGA-LIST / GETANEKALIST` |
| `GetCoverageAneka_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `GETCOVERAGEANEKA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=26 | — |
| `GetCoverageListMBUClaim.xml` | `CODE-PEGA-LIST` | `GETCOVERAGELISTMBUCLAIM` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `GetCoverageListPA.xml` | `CODE-PEGA-LIST` | `GETCOVERAGELISTPA` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `GetCoverageListTravelClaim.xml` | `CODE-PEGA-LIST` | `GETCOVERAGELISTTRAVELCLAIM` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `GetCurencyCoverage_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `GETCURENCYCOVERAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=18 | — |
| `GetCurrencyName.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `GETCURRENCYNAME` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetInvoiceAttachments.xml` | `ASM-FW-GCNMFW-WORK` | `GETINVOICEATTACHMENTS` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `GetKursObjectItemMBU_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `GETKURSOBJECTITEMMBU_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `ASM-FW-GCNMFW-DATA-OBJECTITEM / GETKURSOBJECTITEM_ACT` |
| `GetKursObjectItemMarine_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `GETKURSOBJECTITEMMARINE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `ASM-FW-GCNMFW-DATA-OBJECTITEM / GETKURSOBJECTITEMMBU_ACT` |
| `GetKursObjectItem_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `GETKURSOBJECTITEM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetLinkService.xml` | `ASM-FW-GISFW-INT-M_LINK_SERVICE` | `GETLINKSERVICE` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetNameCauseofLoss_Act.xml` | `ASM-FW-GISFW-INT-CLAUSE` | `GETNAMECAUSEOFLOSS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `GetNameCurrency_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `GETNAMECURRENCY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GCNMFW-DATA-OBJECTITEM / GETNAMECURRENCY_ACTE` |
| `GetPayAttachmentAdj_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GETPAYATTACHMENTADJ_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `GetPayAttachment_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `GETPAYATTACHMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `GetPaymentClaim_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `GETPAYMENTCLAIM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `GetPaymentList_Act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `GETPAYMENTLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `ASM-FW-GCNMFW-DATA-OBJECTITEM / GETPAYMENTLIST_ACT` |
| `GetProgresClaim_ACT.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `GETPROGRESCLAIM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetReportStatus_Act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `GETREPORTSTATUS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | — |
| `GetSpreadingMarine_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `GETSPREADINGMARINE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | — |
| `GetStatusKasir_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GETSTATUSKASIR_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetTSITreatyIn_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `GETTSITREATYIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=21 | — |
| `GetTotalJumlahAdj_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `GETTOTALJUMLAHADJ_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `GetUrlGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETURLGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / INSERTGOOGLESTORAGE_ACT` |
| `GettsiAneka_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `GETTSIANEKA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `HitServiceToKasir_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `HITSERVICETOKASIR_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=37 | — |
| `InputCatastrope.xml` | `ASM-FW-GCNMFW-WORK` | `INPUTCATASTROPE` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `InputCurrencyValueAct_Register.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `INPUTCURRENCYVALUEACT_REGISTER` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `InputEstimationPre.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `INPUTESTIMATIONPRE` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `InputKodePos1_Act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `INPUTKODEPOS1_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `InputQuotation_PreAct.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `INPUTQUOTATION_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `InsertDocument_Act.xml` | `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | `INSERTDOCUMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `InsertGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERTGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `InsertJsonClaimNonMBU_act.xml` | `ASM-FW-GCNMFW-WORK` | `INSERTJSONCLAIMNONMBU_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `InsertLogServiceClaim.xml` | `ASM-FW-GCNMFW-WORK` | `INSERTLOGSERVICECLAIM` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.MONITORING_KLAIM_LOG | — |
| `InsertProgressClaim.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `INSERTPROGRESSCLAIM` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | `ASM-FW-GCNMFW-WORK / INSERTPROGRESSCLAIM` |
| `KonversiKlaim_Act.xml` | `ASM-FW-GCNMFW-WORK` | `KONVERSIKLAIM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `NewAdjustConsult_Act.xml` | `@BASECLASS` | `NEWADJUSTCONSULT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `Occupation_Act.xml` | `CODE-PEGA-LIST` | `OCCUPATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `PreReject.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `PREREJECT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `PreSecurityReas_Act.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `PRESECURITYREAS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `PreShowRetro_Act.xml` | `ASM-FW-GISFW-DATA-SPREADINGRISK` | `PRESHOWRETRO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `PrintPDFAccep_MultiAksep.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `PRINTPDFACCEP_MULTIAKSEP` | [terverifikasi] pyActivityType=ACTIVITY; steps=46 | — |
| `ProtectCoverage_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `PROTECTCOVERAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | `ASM-FW-GCNMFW-DATA-OBJECT / SETSPREADING_ACT` |
| `ProtectCurrAdjustment.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `PROTECTCURRADJUSTMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `ProtectDownloadFaceClaim.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `PROTECTDOWNLOADFACECLAIM` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `ProtectPrint.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `PROTECTPRINT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `ProtectionDate_Act.xml` | `ASM-FW-GCNMFW-DATA-ESTIMASI` | `PROTECTIONDATE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `ProtectionObjectItem_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `PROTECTIONOBJECTITEM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `ASM-FW-GCNMFW-DATA-OBJECTITEM / SETCONVERTVALUEKURS_OBJECTITEM` |
| `ProteksiDataRegister_Act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `PROTEKSIDATAREGISTER_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `ProtetDuplicateSpread.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `PROTETDUPLICATESPREAD` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SaveAcceptation.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SAVEACCEPTATION` | [terverifikasi] pyActivityType=ACTIVITY; steps=16 | — |
| `SaveAdjusterConsultant_Act.xml` | `@BASECLASS` | `SAVEADJUSTERCONSULTANT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `SaveAdjustmenttoDB_ACT.xml` | `ASM-FW-GCNMFW-WORK` | `SAVEADJUSTMENTTODB_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `ASM-FW-GCNMFW-WORK / INSERTJSONCLAIMNONMBU_ACT` |
| `SaveCFS_ACT.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `SAVECFS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `SaveCatasrtope_Act.xml` | `ASM-FW-GCNMFW-WORK` | `SAVECATASRTOPE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SaveSpreadingSP_act.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `SAVESPREADINGSP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SearchPolis_act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SEARCHPOLIS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | — |
| `SendCloseClaimToKomite.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SENDCLOSECLAIMTOKOMITE` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | `ASM-FW-GCNMFW-WORK-PNC / SENDREJECTCLAIMTOKOMITE2` |
| `SendEmailDLA_ACT.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `SENDEMAILDLA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `SendEmailKlaim.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SENDEMAILKLAIM` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `SendEmailKlaimRejectClose.xml` | `ASM-FW-GCNMFW-WORK` | `SENDEMAILKLAIMREJECTCLOSE` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | — |
| `SendEmail_ACT.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `SENDEMAIL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=27 | — |
| `SendPICProtect_Act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SENDPICPROTECT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=27 | — |
| `SendRejectClaimToKomite2.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SENDREJECTCLAIMTOKOMITE2` | [terverifikasi] pyActivityType=ACTIVITY; steps=16 | `ASM-FW-GCNMFW-WORK-PNC / SENDREJECTCLAIMTOKOMITE` |
| `Set7Hours.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SET7HOURS` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetAdjTypePayment_act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETADJTYPEPAYMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=23 | — |
| `SetAdjsuter_act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SETADJSUTER_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / SETADJSUTER_ACT` |
| `SetCatastrope_act.xml` | `ASM-FW-GCNMFW-INT-CATASTROPHE` | `SETCATASTROPE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetCauseOfLossValue_act.xml` | `@BASECLASS` | `SETCAUSEOFLOSSVALUE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `@BASECLASS / SETSURVERYORSVALUE_ACT` |
| `SetCedant_act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETCEDANT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `SetConsultant_Act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SETCONSULTANT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / SETCONSULTANT_ACT` |
| `SetConvertCurrencyValue_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `SETCONVERTCURRENCYVALUE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetConvertValueKurs_Estimation.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `SETCONVERTVALUEKURS_ESTIMATION` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SetConvertValueKurs_ObjectItem.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `SETCONVERTVALUEKURS_OBJECTITEM` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GCNMFW-DATA-OBJECTITEM / SETCONVERTVALUEKURS_ESTIMATION` |
| `SetDLACedingSOB.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETDLACEDINGSOB` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetDataSobCeding_Act.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `SETDATASOBCEDING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `SetDefNonCatastrope_Act.xml` | `ASM-FW-GCNMFW-WORK` | `SETDEFNONCATASTROPE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetDisable_ACT.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SETDISABLE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | — |
| `SetEditCatastrope.xml` | `ASM-FW-GCNMFW-WORK` | `SETEDITCATASTROPE` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetGrossAdjustment_act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETGROSSADJUSTMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / SETNILAIRESIKOSENDIRI` |
| `SetIndexObj_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `SETINDEXOBJ_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetIndexObject_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `SETINDEXOBJECT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `SetIndex_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `SETINDEX_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | — |
| `SetInitial_ACT.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `SETINITIAL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `SetInputParam_Act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SETINPUTPARAM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SetKomiteList_ACT.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETKOMITELIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SetListKomite_act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETLISTKOMITE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / SETKOMITELIST_ACT` |
| `SetMOClaim_Act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SETMOCLAIM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SetNilaiResikoSendiri.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETNILAIRESIKOSENDIRI` | [terverifikasi] pyActivityType=ACTIVITY; steps=29 | — |
| `SetObjectItem_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `SETOBJECTITEM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `CODE-PEGA-LIST / FILTERPROPERTYITEM_ACT` |
| `SetPayableTo_act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETPAYABLETO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SetPayable_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETPAYABLE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `SetProtectionEstimation.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `SETPROTECTIONESTIMATION` | [terverifikasi] pyActivityType=ACTIVITY; steps=38 | — |
| `SetQQName.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SETQQNAME` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetRejectClaim_Cancel.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SETREJECTCLAIM_CANCEL` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GCNMFW-WORK-PNC / SETREJECTCLAIM_PRE` |
| `SetRejectClaim_pre.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SETREJECTCLAIM_PRE` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetRemarksKomite.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETREMARKSKOMITE` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetSalvageValue.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETSALVAGEVALUE` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `SetSpreadingAjsutement_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETSPREADINGAJSUTEMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=29 | — |
| `SetSpreading_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `SETSPREADING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SetTSIPremiCedant_Act.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `SETTSIPREMICEDANT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SetTempLocation_Act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SETTEMPLOCATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetTreatNameAdjustment.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETTREATNAMEADJUSTMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetTreatyname.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `SETTREATYNAME` | [terverifikasi] pyActivityType=ACTIVITY; steps=16 | — |
| `SetTryeatyName_ACT.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `SETTRYEATYNAME_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GCNMFW-DATA-OBJECTITEM / SENDEMAIL_ACT` |
| `SetTypePDFAdjustment.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SETTYPEPDFADJUSTMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `SetValueAdjusterFee.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETVALUEADJUSTERFEE` | [terverifikasi] pyActivityType=ACTIVITY; steps=18 | — |
| `SetValueSaveSpreading.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `SETVALUESAVESPREADING` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GCNMFW-DATA-OBJECT / SETVALUESAVESPREADING` |
| `SetchronologyKlaimFacIn.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SETCHRONOLOGYKLAIMFACIN` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SpreadingCheckSP.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `SPREADINGCHECKSP` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | `ASM-FW-GCNMFW-WORK-PNC / SPREADINGCHECK` |
| `TravelDocument_act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `TRAVELDOCUMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `ASM-FW-GCNMFW-WORK-PNC / REQUIREDDOCUMENT_ACT` |
| `UpdateDataPolisFacin_Act.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `UPDATEDATAPOLISFACIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | — |
| `ValidateInputEstimate_act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `VALIDATEINPUTESTIMATE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=18 | — |
| `ValidationAdjustmentKomite.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `VALIDATIONADJUSTMENTKOMITE` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `ViewKomite_act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `VIEWKOMITE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / SETVALUEADJUSTERFEE_ACT` |
| `getStatusKonversi_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GETSTATUSKONVERSI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `setNullPassword_act.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `SETNULLPASSWORD_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GCNMFW-WORK-PNC / SETNULLPASSWORD_ACT` |

### ConnectREST — 8 rule (`RULE-CONNECT-REST`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `HitDLAClaimFacin.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `HITDLACLAIMFACIN` | [terverifikasi] pyServiceName=HitDLAClaimFacin; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | — |
| `KonversiKlaimNonLife.xml` | `ASM-FW-GCNMFW-WORK` | `KONVERSIKLAIMNONLIFE` | [terverifikasi] pyServiceName=KonversiKlaimNonLife; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | `ASM-FW-GCNMFW-WORK-PNC / INSERTREJECTCLAIM` |
| `SendAcceptationToKasir.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SENDACCEPTATIONTOKASIR` | [terverifikasi] pyServiceName=SendAcceptationToKasir; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | — |
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `SERVICEGOOGLE` | [terverifikasi] pyServiceName=ServiceGoogle; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / GOOGLESTORAGE_UPLOAD` |
| `getPayAttachment.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GETPAYATTACHMENT` | [terverifikasi] pyServiceName=getPayAttachment; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | — |
| `getPaymentClaim.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `GETPAYMENTCLAIM` | [terverifikasi] pyServiceName=getPaymentClaim; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | — |
| `getPremiumPaidOn.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `GETPREMIUMPAIDON` | [terverifikasi] pyServiceName=getPremiumPaidOn; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | — |
| `getPremiumPaidOnMarine.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `GETPREMIUMPAIDONMARINE` | [terverifikasi] pyServiceName=getPremiumPaidOnMarine; baseURL=URL; resolusi=JNDIName; URL_LITERAL_DITEMUKAN=2 | `ASM-FW-GCNMFW-WORK-PNC / GETPREMIUMPAIDON` |

### DataPage — 8 rule (`RULE-DECLARE-PAGES`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `D_AnekaList.xml` | — | `D_ANEKALIST` | [terverifikasi] pyPageName=D_AnekaList; scope=thread; class=ASM-FW-GISFW-Data-Aneka; struktur=list | `D_ANEKALIST / #20171207T065339.976` |
| `D_CoverageMBUClaimList.xml` | — | `D_COVERAGEMBUCLAIMLIST` | [terverifikasi] pyPageName=D_CoverageMBUClaimList; scope=thread; class=ASM-FW-GISFW-Data-Coverage; struktur=list | `D_COVERAGEMBUCLAIM / #20180725T031153.950` |
| `D_CoveragePAClaimList.xml` | — | `D_COVERAGEPACLAIMLIST` | [terverifikasi] pyPageName=D_CoveragePAClaimList; scope=thread; class=ASM-FW-GISFW-Data-Coverage; struktur=list | `D_COVERAGEPACLAIMLIST / #20180815T040322.463` |
| `D_CoverageTravelClaimList.xml` | — | `D_COVERAGETRAVELCLAIMLIST` | [terverifikasi] pyPageName=D_CoverageTravelClaimList; scope=thread; class=ASM-FW-GISFW-Data-Coverage; struktur=list | `D_COVERAGETRAVELCLAIMLIST / #20180814T032128.872` |
| `D_FilteredCoverageAnekaList.xml` | — | `D_FILTEREDCOVERAGEANEKALIST` | [terverifikasi] pyPageName=D_FilteredCoverageAnekaList; scope=thread; class=ASM-FW-GISFW-Data-Coverage; struktur=list | `D_FILTEREDCOVERAGEANEKALIST / #20180409T034518.539` |
| `D_FilteredCoverageList.xml` | — | `D_FILTEREDCOVERAGELIST` | [terverifikasi] pyPageName=D_FilteredCoverageList; scope=thread; class=ASM-FW-GISFW-Data-Coverage; struktur=list | `D_FILTEREDPROPERTYITEMLIST / #20180112T065752.733` |
| `D_FilteredPropertyItemList.xml` | — | `D_FILTEREDPROPERTYITEMLIST` | [terverifikasi] pyPageName=D_FilteredPropertyItemList; scope=thread; class=ASM-FW-GISFW-Data-PropertyItem; struktur=list | `D_FILTEREDPROPERTYITEMLIST / #20180112T065752.733` |
| `D_OccupationList.xml` | — | `D_OCCUPATIONLIST` | [terverifikasi] pyPageName=D_OccupationList; scope=thread; class=ASM-FW-GISFW-Data-Occupation; struktur=list | `D_OCCUPATIONLIST / #20210721T023606.426` |

### DataTransform — 28 rule (`RULE-OBJ-MODEL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BackFromRegister.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `BACKFROMREGISTER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BackToEstimasi.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `BACKTOESTIMASI` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CNMRefreshCauseOfLoss_dt.xml` | `@BASECLASS` | `CNMREFRESHCAUSEOFLOSS_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / CNMREFRESHSURVEYORS_DT` |
| `CNMRefreshListDetailCauseOfLoss_dt.xml` | `@BASECLASS` | `CNMREFRESHLISTDETAILCAUSEOFLOSS_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CNMShowInsertCauseOfLoss_dt.xml` | `@BASECLASS` | `CNMSHOWINSERTCAUSEOFLOSS_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / CNMSHOWINSERTSUVERYORS_DT` |
| `CNMShowInsertDetailCauseOfLoss_dt.xml` | `@BASECLASS` | `CNMSHOWINSERTDETAILCAUSEOFLOSS_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / CNMINSERTDETAILCAUSEOFLOSS_DT` |
| `ChronologyInsertion_DT.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `CHRONOLOGYINSERTION_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CloseStsSaveTreatyGroup.xml` | `ASM-FW-GISFW-INT-TREATYGROUP` | `CLOSESTSSAVETREATYGROUP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CopyCurrency.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `COPYCURRENCY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputRegisterPostDT.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `INPUTREGISTERPOSTDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InsertObjectItemList_DT.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `INSERTOBJECTITEMLIST_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InsertObjects_dt.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `INSERTOBJECTS_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `RemoveDataClaimData_DT.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `REMOVEDATACLAIMDATA_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SearchHierarkiSourceBizAgent_PostDT.xml` | `ASM-FW-GISFW-DATA-AGENT` | `SEARCHHIERARKISOURCEBIZAGENT_POSTDT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetCedingCedant_DT.xml` | `ASM-FW-GISFW-DATA-AGENT` | `SETCEDINGCEDANT_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-AGENT / SETCEDINGCO_DT` |
| `SetCoverageID_DT.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `SETCOVERAGEID_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetCurrencyAdjustment_DT.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETCURRENCYADJUSTMENT_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetDisable_DT.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SETDISABLE_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetEstimation_DT.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SETESTIMATION_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetIndexObject_DT.xml` | `ASM-FW-GISFW-DATA-ANEKA` | `SETINDEXOBJECT_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetLabel_Dt.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SETLABEL_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetNilaiEstimasi_DT.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `SETNILAIESTIMASI_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetObjectID.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `SETOBJECTID` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetStartDateAdjustment.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SETSTARTDATEADJUSTMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetTempAdjustID.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETTEMPADJUSTID` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `btnCedingCO_DT.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `BTNCEDINGCO_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `btnCedingCedant_DT.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `BTNCEDINGCEDANT_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-QUOTATION / BTNCEDINGCO_DT` |
| `btnSOB_DT.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `BTNSOB_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### DecisionTable — 1 rule (`RULE-DECLARE-DECISIONTABLE`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `GetMimeType.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETMIMETYPE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Flow — 1 rule (`RULE-OBJ-FLOW`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `Register_Flow.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `REGISTER_FLOW` | [terverifikasi] pyStartActivity=Start1 | — |

### FlowAction — 33 rule (`RULE-OBJ-FLOWACTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `Adjusment_FA.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `ADJUSMENT_FA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `AgentSourceBizDetails.xml` | `ASM-FW-GISFW-DATA-AGENT` | `AGENTSOURCEBIZDETAILS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CatastrofeList.xml` | `ASM-FW-GCNMFW-WORK` | `CATASTROFELIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CedingCedant.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `CEDINGCEDANT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Estimasi.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `ESTIMASI` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `EstimasiMarine_FA.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `ESTIMASIMARINE_FA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `EstimasiPA.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `ESTIMASIPA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-DATA-OBJECTITEM / ESTIMASI` |
| `InputAdjustment.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `INPUTADJUSTMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputEstimasi.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `INPUTESTIMASI` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputRegister.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `INPUTREGISTER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputSubProgressClaim.xml` | `ASM-FW-GCNMFW-DATA-CLAIMDATA` | `INPUTSUBPROGRESSCLAIM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputSurveyor.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `INPUTSURVEYOR` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `MessageBeforeDeleteTreatyGroup.xml` | `ASM-FW-GISFW-INT-TREATYGROUP` | `MESSAGEBEFOREDELETETREATYGROUP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PreventRejectClaim.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `PREVENTREJECTCLAIM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PropertyItemListAdjustementTravel_FW.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `PROPERTYITEMLISTADJUSTEMENTTRAVEL_FW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-DATA-OBJECT / SHOWPROPERTYITEMLISTADJUSTEMENTMBU_FW` |
| `PropertyItemListGridEstimationMBU_FA.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `PROPERTYITEMLISTGRIDESTIMATIONMBU_FA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PropertyItemListGridEstimationTravel_FA.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `PROPERTYITEMLISTGRIDESTIMATIONTRAVEL_FA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ProtectDOL.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `PROTECTDOL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `RejectSurveyClaim.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `REJECTSURVEYCLAIM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetPassWordSP.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `SETPASSWORDSP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowPropertyItemListAdjustementMBU_FW.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `SHOWPROPERTYITEMLISTADJUSTEMENTMBU_FW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-DATA-OBJECT / SHOWPROPERTYITEMLISTADJUSTEMENT_FW` |
| `ShowPropertyItemListAdjustementPA_FW.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `SHOWPROPERTYITEMLISTADJUSTEMENTPA_FW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowPropertyItemListAdjustement_FW.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `SHOWPROPERTYITEMLISTADJUSTEMENT_FW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-DATA-OBJECT / SHOWPROPERTYITEMLISTESTIMATION_FLOWACT` |
| `ShowPropertyItemListEstimationMarine_FlowAct.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `SHOWPROPERTYITEMLISTESTIMATIONMARINE_FLOWACT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowPropertyItemListEstimationPA_FlowAct.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `SHOWPROPERTYITEMLISTESTIMATIONPA_FLOWACT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowPropertyItemListEstimationTravel_FA.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `SHOWPROPERTYITEMLISTESTIMATIONTRAVEL_FA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-DATA-OBJECT / SHOWPROPERTYITEMLISTESTIMATIONMARINE_FLOWACT` |
| `ShowPropertyItemListEstimation_FlowAct.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `SHOWPROPERTYITEMLISTESTIMATION_FLOWACT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowRetro.xml` | `ASM-FW-GISFW-DATA-SPREADINGRISK` | `SHOWRETRO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowSecurityReinsurer.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `SHOWSECURITYREINSURER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SubProgresClaim.xml` | `ASM-FW-GCNMFW-DATA-CLAIMDATA` | `SUBPROGRESCLAIM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SureRejectClaim.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SUREREJECTCLAIM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewDetailPayment.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY` | `VIEWDETAILPAYMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewDetailPayment_FA.xml` | `ASM-FW-GCNMFW-DATA-CURRENCY` | `VIEWDETAILPAYMENT_FA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Harness — 16 rule (`RULE-HTML-HARNESS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `CauseofLoss_Harness.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `CAUSEOFLOSS_HARNESS` | [terverifikasi] pyInclude=1 | — |
| `CedingCedant.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `CEDINGCEDANT` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-DATA-QUOTATION / CEDINGCOMPANY` |
| `ChoosePolis.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `CHOOSEPOLIS` | [terverifikasi] pyInclude=5 | — |
| `Comittee.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `COMITTEE` | [terverifikasi] pyInclude=1 | — |
| `ListPaymentClaim_Harness.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `LISTPAYMENTCLAIM_HARNESS` | [terverifikasi] pyInclude=1 | — |
| `ListPaymentPremi_Harness.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `LISTPAYMENTPREMI_HARNESS` | [terverifikasi] pyInclude=1 | — |
| `ListPayment_Harness.xml` | `ASM-FW-GISFW-WORK-NB` | `LISTPAYMENT_HARNESS` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-DATA-PROPERTYITEM / LISTPAYMENT_HARNESS` |
| `MstAdjusterConsultant.xml` | `DATA-PORTAL` | `MSTADJUSTERCONSULTANT` | [terverifikasi] pyInclude=1 | — |
| `Outstanding.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `OUTSTANDING` | [terverifikasi] pyInclude=1 | — |
| `Pla_Dtl_Harness.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `PLA_DTL_HARNESS` | [terverifikasi] pyInclude=1 | — |
| `PrintDLA_dtl_Harness.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `PRINTDLA_DTL_HARNESS` | [terverifikasi] pyInclude=1 | — |
| `RetroList_Harnness.xml` | `ASM-FW-GCNMFW-WORK` | `RETROLIST_HARNNESS` | [terverifikasi] pyInclude=1 | — |
| `TambahCauseofLoss.xml` | `DATA-PORTAL` | `TAMBAHCAUSEOFLOSS` | [terverifikasi] pyInclude=3 | — |
| `TambahMasterCauseOfLoss.xml` | `DATA-PORTAL` | `TAMBAHMASTERCAUSEOFLOSS` | [terverifikasi] pyInclude=3 | — |
| `ViewAttachment.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `VIEWATTACHMENT` | [terverifikasi] pyInclude=1 | — |
| `ViewCedantpanels.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `VIEWCEDANTPANELS` | [terverifikasi] pyInclude=1 | — |

### RDBList — 63 rule (`RULE-CONNECT-SQL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseHistoryClaim.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `BROWSEHISTORYCLAIM` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_KLAIM | `ASM-FW-GCNMFW-INT-V_STS_CLAIM / GCNM!BROWSEHISTORYCLAIM` |
| `BrowseRW_SQL.xml` | `ASM-FW-GISFW-INT-RW` | `BROWSERW_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.RW | — |
| `CariHistoryClaim_SQL.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `CARIHISTORYCLAIM_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_KLAIM | — |
| `CekLunasPremi_Sql.xml` | `ASM-FW-GCNMFW-WORK` | `CEKLUNASPREMI_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=ARASAPAS.INVOICE | — |
| `CekProteksiKlaim.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `CEKPROTEKSIKLAIM` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.OPENPROTEKSI_EDM | — |
| `CheckEDMPolicy.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `CHECKEDMPOLICY` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | — |
| `CurrencyStandard.xml` | `ASM-FW-GISFW-INT-CURRENCYSTANDARD` | `CURRENCYSTANDARD` | [terverifikasi] sqlKind=QUERY; procs=POOLDATA.GETCURRENCYSTANDARD; sqlOps=SELECT | — |
| `DeleteDataTreatyGroup_SQL.xml` | `ASM-FW-GISFW-INT-TREATYGROUP` | `DELETEDATATREATYGROUP_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=DELETE; tables=M_TREATYGROUP | — |
| `GenerateImageID_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GENERATEIMAGEID_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GetAcceptation_SQL.xml` | `ASM-FW-GCNMFW-INT-OS_AKSEPTASI_KLAIM` | `GETACCEPTATION_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=OS_AKSEPTASI_KLAIM | — |
| `GetAddressCeding.xml` | `ASM-FW-GISFW-INT-CLIENT` | `GETADDRESSCEDING` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_CLIENT,AGENT | — |
| `GetAddressTreatyIn.xml` | `ASM-FW-GISFW-INT-CLIENT` | `GETADDRESSTREATYIN` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_CLIENT,AGENT | `ASM-FW-GISFW-INT-CLIENT / ASM!GETADDRESSCEDING` |
| `GetAllClaimWithSameNopolis1_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETALLCLAIMWITHSAMENOPOLIS1_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=OS_AKSEPTASI_KLAIM | — |
| `GetAllClaimWithSameNopolis_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETALLCLAIMWITHSAMENOPOLIS_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=REINSURANCE.TRLOSS_DETAIL_T,REINSURANCE.TREATY_LOSS | — |
| `GetAllDataPA_SQL.xml` | `ASM-FW-GCNMFW-INT-OS_AKSEPTASI_KLAIM` | `GETALLDATAPA_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=OS_AKSEPTASI_KLAIM | — |
| `GetAllDataTravel_Act.xml` | `ASM-FW-GCNMFW-INT-OS_AKSEPTASI_KLAIM` | `GETALLDATATRAVEL_ACT` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=OS_AKSEPTASI_KLAIM | — |
| `GetAppName_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETAPPNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.T_FOLDER_IMAGE | — |
| `GetCaseIDNB.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `GETCASEIDNB` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | — |
| `GetCedingName.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETCEDINGNAME` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=FACINPRODUCTION | — |
| `GetCedingNameAneka.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETCEDINGNAMEANEKA` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=FACINPRODUCTION | `ASM-FW-GISFW-INT-POLICYJSON / GCNM!GETCEDINGNAME` |
| `GetCedingNameMBU.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETCEDINGNAMEMBU` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=FACINPRODUCTION | `ASM-FW-GISFW-INT-POLICYJSON / GCNM!GETCEDINGNAMEMARINE` |
| `GetCedingNameMarine.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETCEDINGNAMEMARINE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=FACINPRODUCTION | `ASM-FW-GISFW-INT-POLICYJSON / GCNM!GETCEDINGNAMEANEKA` |
| `GetCopyNBForClaim.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `GETCOPYNBFORCLAIM` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | `ASM-FW-GISFW-INT-OFFERJSON / ASM!GETCOPYNB` |
| `GetCurrency.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `GETCURRENCY` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CURRENCY | `ASM-FW-GISFW-INT-CURRENCY / ASM!UPDATEMASTERCURRENCY` |
| `GetDataBankAccount2_sql.xml` | `ASM-FW-GISFW-INT-BANKACCOUNT` | `GETDATABANKACCOUNT2_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=BANKACCOUNT | `ASM-FW-GISFW-INT-BANKACCOUNT / GCNM!GETDATABANKACCOUNT_SQL` |
| `GetDataBankAccount_sql.xml` | `ASM-FW-GISFW-INT-BANKACCOUNT` | `GETDATABANKACCOUNT_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=BANKACCOUNT | — |
| `GetDataEstimation.xml` | `ASM-FW-GCNMFW-INT-OS_AKSEPTASI_KLAIM` | `GETDATAESTIMATION` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=REINSURANCE.TRLOSS_DETAIL_T | — |
| `GetDataOfferTreatyIn_SQl.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETDATAOFFERTREATYIN_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYINOFFER | — |
| `GetDataOutstandingMBU_Sql.xml` | `ASM-FW-GCNMFW-INT-OS_AKSEPTASI_KLAIM` | `GETDATAOUTSTANDINGMBU_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=OS_AKSEPTASI_KLAIM | — |
| `GetDataOutstanding_SQL.xml` | `ASM-FW-GCNMFW-INT-OS_AKSEPTASI_KLAIM` | `GETDATAOUTSTANDING_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=OS_AKSEPTASI_KLAIM | — |
| `GetDataTreatyLimit_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETDATATREATYLIMIT_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYBUSINESS,PROPORTIONALARRG,TREATYCONTRACT | `ASM-FW-GISFW-INT-POLICYJSON / GCNM!GETYEARFORGROUPIID_SQL` |
| `GetEmailCeding_SQL.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GETEMAILCEDING_SQL` | [terverifikasi] sqlKind=QUERY; procs=GL.F_GET_EMAIL; sqlOps=SELECT | — |
| `GetIDConsultanAdj_SQL.xml` | `ASM-FW-GISFW-INT-ADJUSTERCONSULTANT` | `GETIDCONSULTANADJ_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_SITE_DATABASE | `ASM-FW-GISFW-INT-ADJUSTERCONSULTANT / RNM!SAVEMASTERADJUSTERCONSULTANT_SQL` |
| `GetKodeProdNonLife_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETKODEPRODNONLIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.KODE_PRODUKSI | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETSEQUENCENUMBER_SQL` |
| `GetLBUID_SQL.xml` | `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS_BUSINESS` | `GETLBUID_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=V_D_CAUSE_OF_LOSS_BUSINESS | `ASM-FW-GCNMFW-INT-V_M_CAUSE_OF_LOSS_BUSINESS / GCNM!GETLBUID_SQL` |
| `GetLeaderReport.xml` | `ASM-FW-GISFW-INT-AGENT` | `GETLEADERREPORT` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=AGENT | `ASM-FW-GISFW-INT-CLIENT / ASM!GETADDRESSCEDING` |
| `GetLimitDirekturUtama_SQL.xml` | `ASM-FW-GCNMFW-INT-EMAILKOMITE` | `GETLIMITDIREKTURUTAMA_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.EMAILKOMITE | — |
| `GetLimitPLADLA_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLIMITPLADLA_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=PROPORTIONALARRG | — |
| `GetLinkStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETLINKSTORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GETTOKENSTORAGE_SQL` |
| `GetListRetro_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLISTRETRO_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYREINSURER | — |
| `GetNopolisNoklaim.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETNOPOLISNOKLAIM` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_KLAIM | — |
| `GetPeriodeTreaty_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETPERIODETREATY_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYYEAR | — |
| `GetPolisForClaim_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETPOLISFORCLAIM_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_BUSINESS,JSON_POLIS,POOLDATA.FACINPRODUCTION | — |
| `GetProgressClaim_SQL.xml` | `ASM-FW-GCNMFW-DATA-CLAIMDATA` | `GETPROGRESSCLAIM_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.PROGRESSCLAIM | `ASM-FW-GCNMFW-WORK / RNM!GETPROGRESSCLAIM_SQL` |
| `GetQuotaShare.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETQUOTASHARE` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=PROPORTIONALARRG | — |
| `GetSequenceNumber_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETSEQUENCENUMBER_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETTANGGALCLOSING_SQL` |
| `GetStatusKasir_SQL.xml` | `ASM-FW-GCNMFW-WORK` | `GETSTATUSKASIR_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.DIRECTTOKASIR_LOG | — |
| `GetSubProgressClaim_SQL.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `GETSUBPROGRESSCLAIM_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.SUBPROGRESSCLAIM | `ASM-FW-GCNMFW-DATA-CLAIMDATA / RNM!GETPROGRESSCLAIM_SQL` |
| `GetTokenStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETTOKENSTORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.GET_TOKEN_STORAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |
| `GetTreatyGroup_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETTREATYGROUP_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYBUSINESS | — |
| `GetUpdateNBForClaim.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `GETUPDATENBFORCLAIM` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | — |
| `GetYearforGroupiID_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETYEARFORGROUPIID_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYCONTRACT | — |
| `InsertClaimPNC.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `INSERTCLAIMPNC` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_KLAIM_PNC | — |
| `InsertDLA_OS_SQL.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `INSERTDLA_OS_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=POOLDATA.OS_AKSEPTASI_KLAIM | — |
| `InsertLOGDirectKasir_SQL.xml` | `ASM-FW-GCNMFW-WORK` | `INSERTLOGDIRECTKASIR_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.DIRECTTOKASIR_LOG | `ASM-FW-GISFW-WORK / RNM!INSERTLOGMOP_SQL` |
| `InsertProgressClaim_SQL.xml` | `ASM-FW-GCNMFW-WORK` | `INSERTPROGRESSCLAIM_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_PROGRESSCLAIM | — |
| `InsertSUBProgressClaim_SQL.xml` | `ASM-FW-GCNMFW-WORK` | `INSERTSUBPROGRESSCLAIM_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_SUBPROGRESSCLAIM | `ASM-FW-GCNMFW-WORK / RNM!INSERTPROGRESSCLAIM_SQL` |
| `Insert_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERT_T_STORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GENERATEIMAGEID_SQL` |
| `SaveOSClaim_SQL.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `SAVEOSCLAIM_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM | — |
| `UpdateDCauseOfLoss.xml` | `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS` | `UPDATEDCAUSEOFLOSS` | [terverifikasi] sqlKind=PLSQL; procs=DBMS_LOB.CREATETEMPORARY,POOLDATA.PEGA_D_CAUSE_OF_LOSS | `ASM-FW-CNMFW-INT-V_D_CAUSE_OF_LOSS / CNM!UPDATEDCAUSEOFLOSS` |
| `UpdateTotalJob_sql.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `UPDATETOTALJOB_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=MST_USER_TEKNIS | `ASM-FW-GCNMFW-INT-V_POLIS / GCNM!UPDATELASTTEAMGETSURVEYOR_SQL` |
| `Update_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `UPDATE_T_STORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |
| `getStatusKonversi_SQL.xml` | `ASM-FW-GCNMFW-INT-OS_AKSEPTASI_KLAIM` | `GETSTATUSKONVERSI_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=REINSURANCE.TRLOSS_DETAIL_T | — |

### ReportDefinition — 27 rule (`RULE-OBJ-REPORT-DEFINITION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseAdjusterConsultant.xml` | `ASM-FW-GISFW-INT-ADJUSTERCONSULTANT` | `BROWSEADJUSTERCONSULTANT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseAgentHierarkiList_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSEAGENTHIERARKILIST_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseAgentNusaRe_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSEAGENTNUSARE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseBusiness_RD.xml` | `ASM-FW-GISFW-INT-BUSINESS` | `BROWSEBUSINESS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCaseOPCList.xml` | `ASM-FW-GCNMFW-WORK-OPENPROTECTION` | `BROWSECASEOPCLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCedingCo_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSECEDINGCO_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-AGENT / BROWSEAGENTHIERARKILIST_RD` |
| `BrowseCity_RD.xml` | `ASM-FW-GISFW-INT-RW` | `BROWSECITY_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-RW / BROWSEPROVINCE_RD` |
| `BrowseCountry_RD.xml` | `ASM-FW-GISFW-INT-COUNTRY` | `BROWSECOUNTRY_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCouseOfLoss_Business.xml` | `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS_BUSINESS` | `BROWSECOUSEOFLOSS_BUSINESS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCurrency_RD.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `BROWSECURRENCY_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseDistrictInputC_RD.xml` | `ASM-FW-GISFW-INT-DISTRICT` | `BROWSEDISTRICTINPUTC_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseDocumentTravel_Rd.xml` | `ASM-FW-GCNMFW-INT-V_LST_DOC_TRAVEL_COVERAGE` | `BROWSEDOCUMENTTRAVEL_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseProvince_RD.xml` | `ASM-FW-GISFW-INT-RW` | `BROWSEPROVINCE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseRWInput_RD.xml` | `ASM-FW-GISFW-INT-RW` | `BROWSERWINPUT_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseRW_RD.xml` | `ASM-FW-GISFW-INT-RW` | `BROWSERW_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseReinsuranceType_RD.xml` | `ASM-FW-GISFW-INT-REINSURANCETYPE` | `BROWSEREINSURANCETYPE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseSearchSobCeding_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSESEARCHSOBCEDING_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-AGENT / BROWSECEDINGCO_RD` |
| `BrowseTreatyGroup_RD.xml` | `ASM-FW-GISFW-INT-TREATYGROUP` | `BROWSETREATYGROUP_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseVDCauseOfLoss_RD.xml` | `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS` | `BROWSEVDCAUSEOFLOSS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseVMCauseOfLoss_RD.xml` | `ASM-FW-GCNMFW-INT-V_M_CAUSE_OF_LOSS` | `BROWSEVMCAUSEOFLOSS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseVMstUserTeknis_RD.xml` | `ASM-FW-GCNMFW-INT-V_MST_USER_TEKNIS` | `BROWSEVMSTUSERTEKNIS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseVStsClaim_RD.xml` | `ASM-FW-GCNMFW-INT-V_STS_CLAIM` | `BROWSEVSTSCLAIM_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `FilterEmailKomiteWithLimit.xml` | `ASM-FW-GCNMFW-INT-EMAILKOMITE` | `FILTEREMAILKOMITEWITHLIMIT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GetCatastrope_RD.xml` | `ASM-FW-GCNMFW-INT-CATASTROPHE` | `GETCATASTROPE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GetListEdm.xml` | `ASM-FW-GISFW-WORK-ENDORSEMENT` | `GETLISTEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK-ENDORSEMENT / INBOXEDM_RD2` |
| `RejectedClaim_RD.xml` | `ASM-FW-GCNMFW-INT-CLAIMREJECTED` | `REJECTEDCLAIM_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SelectVDCauseOfLoss_RD.xml` | `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS` | `SELECTVDCAUSEOFLOSS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-CNMFW-INT-V_D_CAUSE_OF_LOSS / SELECTVDCAUSEOFLOSS_RD` |

### Section — 57 rule (`RULE-HTML-SECTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `Adjusment_SC.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `ADJUSMENT_SC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCauseOfLoss.xml` | `@BASECLASS` | `BROWSECAUSEOFLOSS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseDetailCauseOfLoss.xml` | `@BASECLASS` | `BROWSEDETAILCAUSEOFLOSS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CatastrofeList_Sec.xml` | `ASM-FW-GCNMFW-WORK` | `CATASTROFELIST_SEC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Catastrope_Sec.xml` | `ASM-FW-GCNMFW-WORK` | `CATASTROPE_SEC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CauseofLoss_Section.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `CAUSEOFLOSS_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CedingCedant.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `CEDINGCEDANT` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-DATA-QUOTATION / CEDINGCOMPANY` |
| `CedingCedantHierarki.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `CEDINGCEDANTHIERARKI` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-QUOTATION / CEDINGCOHIERARKI` |
| `ClaimComite.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `CLAIMCOMITE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ClaimComiteeReject.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `CLAIMCOMITEEREJECT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ClaimSurvey.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `CLAIMSURVEY` | [terverifikasi] pyInclude=6 | — |
| `Estimasi.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `ESTIMASI` | [terverifikasi] pyInclude=4 | — |
| `EstimasiMarine.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `ESTIMASIMARINE` | [terverifikasi] pyInclude=2 | `ASM-FW-GCNMFW-DATA-OBJECTITEM / ESTIMASI` |
| `EstimasiPA.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `ESTIMASIPA` | [terverifikasi] pyInclude=2 | — |
| `GridCauseOfLoss.xml` | `DATA-PORTAL` | `GRIDCAUSEOFLOSS` | [terverifikasi] pyInclude=2 | — |
| `InputAdjustment.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `INPUTADJUSTMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputEstimasi.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `INPUTESTIMASI` | [terverifikasi] pyInclude=6 | — |
| `InputEstimasiAdmin.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `INPUTESTIMASIADMIN` | [terverifikasi] pyInclude=8 | — |
| `InputEstimasiDetail.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `INPUTESTIMASIDETAIL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputInwardFacultativeDtl.xml` | `ASM-FW-GCNMFW-WORK` | `INPUTINWARDFACULTATIVEDTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputRegister.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `INPUTREGISTER` | [terverifikasi] pyInclude=6 | — |
| `InputRegisterDetail.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `INPUTREGISTERDETAIL` | [terverifikasi] pyInclude=2 | — |
| `InputSubProgressClaim_Sec.xml` | `ASM-FW-GCNMFW-DATA-CLAIMDATA` | `INPUTSUBPROGRESSCLAIM_SEC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsDeleteTreatyGroup.xml` | `ASM-FW-GISFW-INT-TREATYGROUP` | `ISDELETETREATYGROUP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ItemListEstimation.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `ITEMLISTESTIMATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ItemListEstimationMBU.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `ITEMLISTESTIMATIONMBU` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-DATA-OBJECT / ITEMLISTESTIMATION` |
| `ItemListEstimationPA.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `ITEMLISTESTIMATIONPA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ItemListEstimationTravel.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `ITEMLISTESTIMATIONTRAVEL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-DATA-OBJECT / ITEMLISTESTIMATIONMBU` |
| `ListDetailCauseOfLoss.xml` | `DATA-PORTAL` | `LISTDETAILCAUSEOFLOSS` | [terverifikasi] pyInclude=2 | — |
| `ListPaymentClaim_SC.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `LISTPAYMENTCLAIM_SC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ListPayment_SC.xml` | `ASM-FW-GISFW-WORK-NB` | `LISTPAYMENT_SC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-PROPERTYITEM / LISTPAYMENT_SC` |
| `MstAdjusterConsultant.xml` | `DATA-PORTAL` | `MSTADJUSTERCONSULTANT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Outstanding_SC.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `OUTSTANDING_SC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PaymentPremiList_SC.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `PAYMENTPREMILIST_SC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-DATA-OBJECTITEM / PAYMENTPREMILIST_SC` |
| `Pla_Dtl.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `PLA_DTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PreventRejectClaim.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `PREVENTREJECTCLAIM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PrintDLA_dtl.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `PRINTDLA_DTL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ProgresClaim_Sec.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `PROGRESCLAIM_SEC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PropertyItemListGridEstimation.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `PROPERTYITEMLISTGRIDESTIMATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PropertyItemListGridEstimationMBU.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `PROPERTYITEMLISTGRIDESTIMATIONMBU` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PropertyItemListGridEstimationMarine.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `PROPERTYITEMLISTGRIDESTIMATIONMARINE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-DATA-OBJECT / PROPERTYITEMLISTGRIDESTIMATION` |
| `PropertyItemListGridEstimationTravel.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `PROPERTYITEMLISTGRIDESTIMATIONTRAVEL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ProtectDOL.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `PROTECTDOL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `RetroList_SC.xml` | `ASM-FW-GCNMFW-WORK` | `RETROLIST_SC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetDetPassWordSP_sc.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `SETDETPASSWORDSP_SC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowItemPA.xml` | `ASM-FW-GCNMFW-DATA-OBJECT` | `SHOWITEMPA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowObjectAdj.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SHOWOBJECTADJ` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowRetro_Sec.xml` | `ASM-FW-GISFW-DATA-SPREADINGRISK` | `SHOWRETRO_SEC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowSecurityReinsurer.xml` | `ASM-FW-GISFW-DATA-FACOFFER` | `SHOWSECURITYREINSURER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SubProgresClaim_Sec.xml` | `ASM-FW-GCNMFW-DATA-CLAIMDATA` | `SUBPROGRESCLAIM_SEC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SureRejectClaim_section.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `SUREREJECTCLAIM_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewAttachment.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `VIEWATTACHMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / VIEWADJUSTMENTATTACHMENT` |
| `ViewCedantPanel.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `VIEWCEDANTPANEL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewDetailPayment_SC.xml` | `ASM-FW-GCNMFW-DATA-CURRENCY` | `VIEWDETAILPAYMENT_SC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewHistoryClaim.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `VIEWHISTORYCLAIM` | [terverifikasi] pyInclude=2 | — |
| `ViewPolis.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `VIEWPOLIS` | [terverifikasi] pyInclude=2 | — |
| `deleteItemButton.xml` | `@BASECLASS` | `DELETEITEMBUTTON` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### SystemSettings — 1 rule (`RULE-ADMIN-SYSTEM-SETTINGS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `LinkService.xml` | `LINKSERVICE` | `LINKSERVICE` | [terverifikasi] pyPurpose=LinkService | — |

### When — 60 rule (`RULE-OBJ-WHEN`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `IsAllRisk.xml` | `ASM-FW-GISFW-DATA` | `ISALLRISK` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsAneka.xml` | `ASM-FW-GISFW-DATA` | `ISANEKA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsAviationHull.xml` | `ASM-FW-GISFW-DATA` | `ISAVIATIONHULL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsBackStage.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `ISBACKSTAGE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsBillboardNeon.xml` | `ASM-FW-GISFW-WORK` | `ISBILLBOARDNEON` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISBILLBOARDNEON` |
| `IsBoiler.xml` | `@BASECLASS` | `ISBOILER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsBonding.xml` | `ASM-FW-GISFW-WORK` | `ISBONDING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsBondingAndCustomBonds.xml` | `ASM-FW-GISFW-DATA` | `ISBONDINGANDCUSTOMBONDS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISBONDING` |
| `IsBondingKBG.xml` | `@BASECLASS` | `ISBONDINGKBG` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsBurglary.xml` | `ASM-FW-GISFW-DATA` | `ISBURGLARY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCAR.xml` | `@BASECLASS` | `ISCAR` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCIS.xml` | `@BASECLASS` | `ISCIS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCIT.xml` | `@BASECLASS` | `ISCIT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCLM.xml` | `@BASECLASS` | `ISCLM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCLMNP.xml` | `@BASECLASS` | `ISCLMNP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / ISCLMP` |
| `IsCLMP.xml` | `@BASECLASS` | `ISCLMP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / ISCLM` |
| `IsClaim.xml` | `@BASECLASS` | `ISCLAIM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsContractorsPlantMachinery.xml` | `@BASECLASS` | `ISCONTRACTORSPLANTMACHINERY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCrime.xml` | `ASM-FW-GISFW-DATA` | `ISCRIME` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCustomBond.xml` | `ASM-FW-GCNMFW-WORK` | `ISCUSTOMBOND` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISCUSTOMBOND` |
| `IsEDM.xml` | `@BASECLASS` | `ISEDM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISEDM` |
| `IsEar.xml` | `@BASECLASS` | `ISEAR` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsEdmAdjRefNo.xml` | `ASM-FW-GISFW-WORK` | `ISEDMADJREFNO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISEDMADJREFNO` |
| `IsEdmAdjShareCedant.xml` | `ASM-FW-GISFW-WORK` | `ISEDMADJSHARECEDANT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISEDMADJREFNO` |
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
| `IsLandRig.xml` | `@BASECLASS` | `ISLANDRIG` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsLiability.xml` | `@BASECLASS` | `ISLIABILITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsMBD.xml` | `ASM-FW-GISFW-DATA` | `ISMBD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsMBU.xml` | `ASM-FW-GISFW-DATA` | `ISMBU` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsMarineCargo.xml` | `ASM-FW-GISFW-WORK` | `ISMARINECARGO` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsMarineHull.xml` | `ASM-FW-GISFW-DATA` | `ISMARINEHULL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsNotTravelPA.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `ISNOTTRAVELPA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsOilGas.xml` | `@BASECLASS` | `ISOILGAS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISOILGAS` |
| `IsPA.xml` | `DATA-PARTY-PERSON` | `ISPA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PARTY-PERSON / ISTRAVEL` |
| `IsPEGAPROD.xml` | `@BASECLASS` | `ISPEGAPROD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `WORK- / ISPEGAPROD` |
| `IsPEGASyariah.xml` | `@BASECLASS` | `ISPEGASYARIAH` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsProductsLiability.xml` | `ASM-FW-GISFW-WORK` | `ISPRODUCTSLIABILITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISPRODUCTSLIABILITY` |
| `IsProfessionalLiability.xml` | `ASM-FW-GISFW-WORK` | `ISPROFESSIONALLIABILITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISPROFESSIONALLIABILITY` |
| `IsSPK.xml` | `ASM-FW-GCNMFW-WORK` | `ISSPK` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsSpreadingUW.xml` | `@BASECLASS` | `ISSPREADINGUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsTravel.xml` | `@BASECLASS` | `ISTRAVEL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISTRAVEL` |
| `IsUW.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN-LOCATIONREINSURANCE` | `ISUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsWorkmenCompensation.xml` | `ASM-FW-GISFW-WORK` | `ISWORKMENCOMPENSATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-WORK / ISCIT` |
| `IsYieldShortfall.xml` | `ASM-FW-GISFW-DATA` | `ISYIELDSHORTFALL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA / ISGROWINGTREES` |
| `isAnalistorTransfer.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `ISANALISTORTRANSFER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isBillboardNeonSyariah.xml` | `@BASECLASS` | `ISBILLBOARDNEONSYARIAH` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isElectronicEquipment.xml` | `@BASECLASS` | `ISELECTRONICEQUIPMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isGolfInsurance.xml` | `ASM-FW-GISFW-DATA` | `ISGOLFINSURANCE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isMaintenance.xml` | `@BASECLASS` | `ISMAINTENANCE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `isPA_PNC.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `ISPA_PNC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `whenRiskType.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `WHENRISKTYPE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

## 4. Tabrakan nama di dalam modul ini

### 4.1 Nama file dipakai di lebih dari satu tipe rule

| Nama file | Tipe | Class |
| --- | --- | --- |
| `InputEstimasi.xml` | FlowAction, Section | `ASM-FW-GCNMFW-WORK-PNC, ASM-FW-GCNMFW-WORK-PNC` |
| `InputRegister.xml` | FlowAction, Section | `ASM-FW-GCNMFW-WORK-PNC, ASM-FW-GCNMFW-WORK-PNC` |
| `InputAdjustment.xml` | FlowAction, Section | `ASM-FW-GCNMFW-DATA-ADJUSTMENT, ASM-FW-GCNMFW-DATA-ADJUSTMENT` |
| `Estimasi.xml` | FlowAction, Section | `ASM-FW-GCNMFW-DATA-OBJECTITEM, ASM-FW-GCNMFW-DATA-OBJECTITEM` |
| `ProtectDOL.xml` | FlowAction, Section | `ASM-FW-GCNMFW-WORK-PNC, ASM-FW-GCNMFW-WORK-PNC` |
| `PreventRejectClaim.xml` | FlowAction, Section | `ASM-FW-GCNMFW-WORK-PNC, ASM-FW-GCNMFW-WORK-PNC` |
| `CedingCedant.xml` | FlowAction, Harness, Section | `ASM-FW-GISFW-DATA-QUOTATION, ASM-FW-GISFW-DATA-QUOTATION, ASM-FW-GISFW-DATA-QUOTATION` |
| `ShowSecurityReinsurer.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-FACOFFER, ASM-FW-GISFW-DATA-FACOFFER` |
| `EstimasiPA.xml` | FlowAction, Section | `ASM-FW-GCNMFW-DATA-OBJECTITEM, ASM-FW-GCNMFW-DATA-OBJECTITEM` |
| `ViewAttachment.xml` | Harness, Section | `ASM-FW-GCNMFW-DATA-OBJECTITEM, ASM-FW-GCNMFW-DATA-OBJECTITEM` |
| `MstAdjusterConsultant.xml` | Harness, Section | `DATA-PORTAL, DATA-PORTAL` |

Total: **11** nama file.

### 4.2 Rule `When` bernama sama di lebih dari satu class

_Tidak ada._

## 5. Keluarga klon (satu asal salinan dipakai >1 rule)

Dibaca dari `<pzOriginalInstanceKey>` yang sama pada beberapa file. Ini **sinyal kandidat
konsolidasi** saat migrasi, bukan bukti bahwa perilakunya identik — kesamaan perilaku baru
dapat dinyatakan setelah isinya dibandingkan di STEP D2.

| Asal salinan (class / nama) | Jumlah | File |
| --- | ---: | --- |
| `ASM-FW-GISFW-INT-CLIENT / ASM!GETADDRESSCEDING` | 3 | `RDBList/GetAddressCeding.xml`, `RDBList/GetAddressTreatyIn.xml`, `RDBList/GetLeaderReport.xml` |
| `ASM-FW-GCNMFW-DATA-OBJECTITEM / ESTIMASI` | 4 | `FlowAction/Estimasi.xml`, `FlowAction/EstimasiPA.xml`, `Section/Estimasi.xml`, `Section/EstimasiMarine.xml` |
| `ASM-FW-GCNMFW-DATA-ADJUSTMENT / SETNILAIRESIKOSENDIRI` | 2 | `Activity/SetGrossAdjustment_act.xml`, `Activity/SetNilaiResikoSendiri.xml` |
| `ASM-FW-GCNMFW-WORK-PNC / PREVENTREJECTCLAIM` | 2 | `FlowAction/PreventRejectClaim.xml`, `Section/PreventRejectClaim.xml` |
| `ASM-FW-GCNMFW-DATA-OBJECTITEM / CHECKLIMIT_ACT` | 2 | `Activity/CheckLimitSpreadingTreaty_Act.xml`, `Activity/CheckLimit_Act1.xml` |
| `DATA-PORTAL / MSTADJUSTERCONSULTANT` | 2 | `Harness/MstAdjusterConsultant.xml`, `Section/MstAdjusterConsultant.xml` |
| `ASM-FW-GCNMFW-DATA-OBJECT / SHOWPROPERTYITEMLISTESTIMATIONMARINE_FLOWACT` | 2 | `FlowAction/ShowPropertyItemListEstimationMarine_FlowAct.xml`, `FlowAction/ShowPropertyItemListEstimationTravel_FA.xml` |
| `ASM-FW-GCNMFW-WORK-PNC / SETREJECTCLAIM_PRE` | 2 | `Activity/SetRejectClaim_Cancel.xml`, `Activity/SetRejectClaim_pre.xml` |
| `ASM-FW-GISFW-DATA-QUOTATION / BTNCEDINGCO_DT` | 2 | `DataTransform/btnCedingCO_DT.xml`, `DataTransform/btnCedingCedant_DT.xml` |
| `ASM-FW-GCNMFW-DATA-OBJECT / PROPERTYITEMLISTGRIDESTIMATION` | 2 | `Section/PropertyItemListGridEstimation.xml`, `Section/PropertyItemListGridEstimationMarine.xml` |
| `CODE-PEGA-LIST / GETANEKALIST` | 2 | `Activity/GetAnekaList.xml`, `Activity/GetCoverageAnekaList.xml` |
| `ASM-FW-GCNMFW-DATA-OBJECTITEM / GETKURSOBJECTITEM_ACT` | 2 | `Activity/GetKursObjectItemMBU_Act.xml`, `Activity/GetKursObjectItem_Act.xml` |
| `ASM-FW-GISFW-DATA-QUOTATION / CEDINGCOMPANY` | 2 | `Harness/CedingCedant.xml`, `Section/CedingCedant.xml` |
| `ASM-FW-GCNMFW-DATA-OBJECTITEM / SETCONVERTVALUEKURS_ESTIMATION` | 2 | `Activity/SetConvertValueKurs_Estimation.xml`, `Activity/SetConvertValueKurs_ObjectItem.xml` |
| `ASM-FW-GISFW-INT-POLICYJSON / GCNM!GETCEDINGNAME` | 2 | `RDBList/GetCedingName.xml`, `RDBList/GetCedingNameAneka.xml` |
| `ASM-FW-GISFW-INT-POLICYJSON / GCNM!GETYEARFORGROUPIID_SQL` | 2 | `RDBList/GetDataTreatyLimit_Sql.xml`, `RDBList/GetYearforGroupiID_Sql.xml` |
| `@BASECLASS / ISCLM` | 2 | `When/IsCLM.xml`, `When/IsCLMP.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / INSERTGOOGLESTORAGE_ACT` | 2 | `Activity/GetUrlGoogleStorage_Act.xml`, `Activity/InsertGoogleStorage_Act.xml` |
| `ASM-FW-GCNMFW-DATA-OBJECT / ITEMLISTESTIMATION` | 2 | `Section/ItemListEstimation.xml`, `Section/ItemListEstimationMBU.xml` |
| `ASM-FW-GCNMFW-WORK-PNC / INPUTREGISTER` | 2 | `FlowAction/InputRegister.xml`, `Section/InputRegister.xml` |
| `ASM-FW-GISFW-INT-RW / BROWSEPROVINCE_RD` | 2 | `ReportDefinition/BrowseCity_RD.xml`, `ReportDefinition/BrowseProvince_RD.xml` |
| `ASM-FW-GISFW-DATA-FACOFFER / SHOWSECURITYREINSURER` | 2 | `FlowAction/ShowSecurityReinsurer.xml`, `Section/ShowSecurityReinsurer.xml` |
| `ASM-FW-GCNMFW-WORK-PNC / CHECKDATE_ACT` | 3 | `Activity/CheckDateReceived_Act.xml`, `Activity/CheckDate_Act.xml`, `Activity/CheckDoubleClaim_Act.xml` |
| `ASM-FW-GCNMFW-WORK-PNC / GETPREMIUMPAIDON` | 2 | `ConnectREST/getPremiumPaidOn.xml`, `ConnectREST/getPremiumPaidOnMarine.xml` |
| `CODE-PEGA-LIST / FILTERPROPERTYITEM_ACT` | 3 | `Activity/FilterCoverage_Act.xml`, `Activity/FilterPropertyItem_Act.xml`, `Activity/SetObjectItem_Act.xml` |
| `ASM-FW-GCNMFW-DATA-OBJECT / SETSPREADING_ACT` | 2 | `Activity/ProtectCoverage_Act.xml`, `Activity/SetSpreading_Act.xml` |
| `ASM-FW-GCNMFW-DATA-OBJECTITEM / SENDEMAIL_ACT` | 2 | `Activity/SendEmail_ACT.xml`, `Activity/SetTryeatyName_ACT.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GENERATEIMAGEID_SQL` | 2 | `RDBList/GenerateImageID_SQL.xml`, `RDBList/Insert_T_Storage_SQL.xml` |
| `ASM-FW-GCNMFW-DATA-ADJUSTMENT / INPUTADJUSTMENT` | 2 | `FlowAction/InputAdjustment.xml`, `Section/InputAdjustment.xml` |
| `ASM-FW-GCNMFW-WORK / RNM!INSERTPROGRESSCLAIM_SQL` | 2 | `RDBList/InsertProgressClaim_SQL.xml`, `RDBList/InsertSUBProgressClaim_SQL.xml` |
| `ASM-FW-GISFW-DATA / ISCRIME` | 2 | `When/IsCrime.xml`, `When/IsEnvironmental.xml` |
| `ASM-FW-GCNMFW-WORK-PNC / PROTECTDOL` | 2 | `FlowAction/ProtectDOL.xml`, `Section/ProtectDOL.xml` |
| `ASM-FW-GCNMFW-WORK-PNC / INPUTESTIMASI` | 2 | `FlowAction/InputEstimasi.xml`, `Section/InputEstimasi.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` | 2 | `RDBList/GetTokenStorage_SQL.xml`, `RDBList/Update_T_Storage_SQL.xml` |
| `ASM-FW-GISFW-INT-AGENT / BROWSEAGENTHIERARKILIST_RD` | 2 | `ReportDefinition/BrowseAgentHierarkiList_RD.xml`, `ReportDefinition/BrowseCedingCo_RD.xml` |
| `ASM-FW-GISFW-INT-BANKACCOUNT / GCNM!GETDATABANKACCOUNT_SQL` | 2 | `RDBList/GetDataBankAccount2_sql.xml`, `RDBList/GetDataBankAccount_sql.xml` |
| `ASM-FW-GCNMFW-DATA-ADJUSTMENT / SETKOMITELIST_ACT` | 2 | `Activity/SetKomiteList_ACT.xml`, `Activity/SetListKomite_act.xml` |
| `ASM-FW-GCNMFW-DATA-OBJECTITEM / CHECKESTIMATEVALUE` | 2 | `Activity/CheckEstimateValue.xml`, `Activity/CountSpreadingClaim_ACT.xml` |
| `ASM-FW-GCNMFW-WORK / INSERTJSONCLAIMNONMBU_ACT` | 2 | `Activity/InsertJsonClaimNonMBU_act.xml`, `Activity/SaveAdjustmenttoDB_ACT.xml` |
| `D_FILTEREDPROPERTYITEMLIST / #20180112T065752.733` | 2 | `DataPage/D_FilteredCoverageList.xml`, `DataPage/D_FilteredPropertyItemList.xml` |
| `ASM-FW-GCNMFW-DATA-OBJECT / SHOWPROPERTYITEMLISTESTIMATION_FLOWACT` | 2 | `FlowAction/ShowPropertyItemListAdjustement_FW.xml`, `FlowAction/ShowPropertyItemListEstimation_FlowAct.xml` |

## 6. Objek database yang disentuh

Diambil dari `<pyBrowseSQL>` pada rule `RDBList`. **Nama objek apa adanya dari SQL** —
kualifikasi skema tidak konsisten di korpus. **Tipe kolom, constraint, dan relasi tidak
diketahui** (OQ-001).

| Tabel/view (apa adanya di SQL) | Dirujuk oleh |
| --- | ---: |
| `OS_AKSEPTASI_KLAIM` | 6 rule |
| `JSON_POLIS` | 5 rule |
| `FACINPRODUCTION` | 4 rule |
| `AGENT` | 3 rule |
| `JSON_KLAIM` | 3 rule |
| `PROPORTIONALARRG` | 3 rule |
| `REINSURANCE.TRLOSS_DETAIL_T` | 3 rule |
| `T_STORAGE_IMAGE` | 3 rule |
| `BANKACCOUNT` | 2 rule |
| `M_CLIENT` | 2 rule |
| `POOLDATA.DIRECTTOKASIR_LOG` | 2 rule |
| `TREATYBUSINESS` | 2 rule |
| `TREATYCONTRACT` | 2 rule |
| `ARASAPAS.INVOICE` | 1 rule |
| `CURRENCY` | 1 rule |
| `MST_USER_TEKNIS` | 1 rule |
| `M_BUSINESS` | 1 rule |
| `M_TREATYGROUP` | 1 rule |
| `POOLDATA.EMAILKOMITE` | 1 rule |
| `POOLDATA.FACINPRODUCTION` | 1 rule |
| `POOLDATA.KODE_PRODUKSI` | 1 rule |
| `POOLDATA.MONITORING_KLAIM_LOG` | 1 rule |
| `POOLDATA.M_SITE_DATABASE` | 1 rule |
| `POOLDATA.OPENPROTEKSI_EDM` | 1 rule |
| `POOLDATA.OS_AKSEPTASI_KLAIM` | 1 rule |
| `POOLDATA.PROGRESSCLAIM` | 1 rule |
| `POOLDATA.RW` | 1 rule |
| `POOLDATA.SUBPROGRESSCLAIM` | 1 rule |
| `POOLDATA.T_FOLDER_IMAGE` | 1 rule |
| `REINSURANCE.TREATY_LOSS` | 1 rule |
| `TREATYINOFFER` | 1 rule |
| `TREATYREINSURER` | 1 rule |
| `TREATYYEAR` | 1 rule |
| `V_D_CAUSE_OF_LOSS_BUSINESS` | 1 rule |

### Stored procedure / function yang dipanggil

**Body prosedur ini TIDAK ADA di korpus** — logika di dalamnya tidak dapat direkonstruksi
dari XML. Lihat OQ-002.

| Procedure | Dipanggil oleh |
| --- | --- |
| `DBMS_LOB.CREATETEMPORARY` | `UpdateDCauseOfLoss.xml` |
| `POOLDATA.GETCURRENCYSTANDARD` | `CurrencyStandard.xml` |
| `POOLDATA.PEGA_PROGRESSCLAIM` | `InsertProgressClaim_SQL.xml` |
| `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM` | `SaveOSClaim_SQL.xml` |
| `POOLDATA.PEGA_JSON_KLAIM_PNC` | `InsertClaimPNC.xml` |
| `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` | `GetSequenceNumber_SQL.xml` |
| `POOLDATA.PEGA_D_CAUSE_OF_LOSS` | `UpdateDCauseOfLoss.xml` |
| `POOLDATA.PEGA_SUBPROGRESSCLAIM` | `InsertSUBProgressClaim_SQL.xml` |
| `POOLDATA.GET_TOKEN_STORAGE` | `GetTokenStorage_SQL.xml` |
| `GL.F_GET_EMAIL` | `GetEmailCeding_SQL.xml` |

## 7. Integrasi eksternal (`ConnectREST`)

Alamat tujuan dicatat sebagai **konfigurasi**, bukan URL literal. `baseURL=SETTING` berarti
alamat diambil dari rule `SystemSettings` yang disebut di kolom Setting; nilai sebenarnya
**tidak ada di korpus**. Kolom terakhir menandai rule yang justru memuat URL literal —
itu temuan, bukan fakta arsitektur (lihat `_summary.md`).

| File | Class | pyServiceName | Sumber base URL | Setting | URL literal di dalam rule? |
| --- | --- | --- | --- | --- | --- |
| `HitDLAClaimFacin.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `HitDLAClaimFacin` | `SETTING` | `LinkService!LinkService` | tidak |
| `KonversiKlaimNonLife.xml` | `ASM-FW-GCNMFW-WORK` | `KonversiKlaimNonLife` | `SETTING` | `LinkService!LinkService` | tidak |
| `SendAcceptationToKasir.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SendAcceptationToKasir` | `SETTING` | `LinkService!LinkService` | tidak |
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `ServiceGoogle` | `SETTING` | `LinkService!LinkService` | tidak |
| `getPayAttachment.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `getPayAttachment` | `SETTING` | `LinkService!LinkService` | tidak |
| `getPaymentClaim.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `getPaymentClaim` | `SETTING` | `LinkService!LinkService` | tidak |
| `getPremiumPaidOn.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `getPremiumPaidOn` | `SETTING` | `LinkService!LinkService` | tidak |
| `getPremiumPaidOnMarine.xml` | `ASM-FW-GCNMFW-WORK-PNC` | `getPremiumPaidOnMarine` | `URL` | — | **YA** |

