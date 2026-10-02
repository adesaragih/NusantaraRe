# Inventaris Rule — Claim Prop

STEP D1, batch 3 (domain claim non-facultative). Dibangkitkan 2026-09-12 dari korpus READ-ONLY
`D:\XML\RNM_BRD\Claim Prop\`.

**Batas step ini:** hanya struktur + metadata. Tidak ada kesimpulan logika bisnis di sini.
Penelusuran perilaku adalah STEP D2 (`../flows/`).

**Identitas rule** diambil dari `<pxInsName>` (lihat `../README.md` §6.4).  
Kolom **Asal salinan** berisi `<pzOriginalInstanceKey>` dan hanya diisi bila **berbeda**
dari identitas — artinya rule ini hasil *Save As* dari rule lain. Tanda `—` berarti sama.
Class `—` berarti `<pxInsName>` tidak berisi bagian class (berlaku pada `DataPage`, yang
beridentitas global, bukan per class).

Perintah audit (dijalankan dari `D:\XML\RNM_BRD\`):

```
find "Claim Prop" -type f -name "*.xml" | wc -l
find "Claim Prop" -type f -name "*.xml" | awk -F/ '{print $2}' | sort | uniq -c
grep -o "<pxInsName>[^<]*" "Claim Prop/<tipe>/<file>.xml" | head -1
grep -o "<pzOriginalInstanceKey>[^<]*" "Claim Prop/<tipe>/<file>.xml" | head -1
```

## 1. Ringkasan terukur

- Total file XML: **270** — cocok dengan `../README.md` §5.1
- **Identitas rule unik: 270** dari 270 file → tidak ada identitas ganda di dalam modul ini
- Class Pega unik: **38**
- RDBList menurut jenis SQL: PLSQL=13, QUERY=41 (`PLSQL` = blok anonim `BEGIN…END;`, umumnya memanggil stored procedure)
- Prefix pada kunci RDBList: `ASM`=11, `RNM`=21, `GCNM`=22 (arti prefix **belum terverifikasi** — OQ-008)

| Tipe rule | Jumlah |
| --- | ---: |
| Activity | 112 |
| ConnectREST | 6 |
| DataTransform | 11 |
| DecisionTable | 1 |
| Flow | 1 |
| FlowAction | 12 |
| Harness | 11 |
| RDBList | 54 |
| ReportDefinition | 19 |
| Section | 36 |
| SystemSettings | 1 |
| When | 6 |
| **TOTAL** | **270** |

## 2. Sebaran class Pega

Class yang **tidak** berawalan `ASM-` adalah class bawaan Pega atau class portal — ruang
lingkup migrasinya belum ditentukan (**OQ-009**).

| Class | Jumlah rule | Aplikasi? |
| --- | ---: | --- |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | 103 | aplikasi |
| `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | 35 | aplikasi |
| `@BASECLASS` | 19 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GCNMFW-WORK` | 19 | aplikasi |
| `ASM-FW-GCNMFW-INT-V_POLIS` | 12 | aplikasi |
| `ASM-FW-GISFW-INT-POLICYJSON` | 10 | aplikasi |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | 10 | aplikasi |
| `DATA-PORTAL` | 7 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-INT-TREATYGROUP` | 5 | aplikasi |
| `ASM-FW-GISFW-DATA-TREATYINTOTAL` | 4 | aplikasi |
| `ASM-FW-GISFW-INT-BANKACCOUNT` | 4 | aplikasi |
| `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS` | 3 | aplikasi |
| `ASM-FW-GCNMFW-INT-V_M_CAUSE_OF_LOSS` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-ADJUSTERCONSULTANT` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-OFFERJSON` | 3 | aplikasi |
| `ASM-FW-GCNMFW-DATA-CURRENCY` | 2 | aplikasi |
| `ASM-FW-GCNMFW-INT-CATASTROPHE` | 2 | aplikasi |
| `ASM-FW-GCNMFW-INT-EMAILKOMITE` | 2 | aplikasi |
| `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS_BUSINESS` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-CLIENT` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-CURRENCY` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-RW` | 2 | aplikasi |
| `ASM-FW-GCNMFW-DATA-ESTIMASI` | 1 | aplikasi |
| `ASM-FW-GCNMFW-INT-CLAIMREJECTED` | 1 | aplikasi |
| `ASM-FW-GCNMFW-INT-CLAIM_MASTER_TREATY` | 1 | aplikasi |
| `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | 1 | aplikasi |
| `ASM-FW-GCNMFW-INT-OS_AKSEPTASI_KLAIM` | 1 | aplikasi |
| `ASM-FW-GCNMFW-INT-V_MST_USER_TEKNIS` | 1 | aplikasi |
| `ASM-FW-GCNMFW-WORK-KOMITETREATY` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-AGENT` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-BUSINESS` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-CLAUSE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-CURRENCYSTANDARD` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-MARKETINGOFFICER` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_LINK_SERVICE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-REINSURANCETYPE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYINDETAIL` | 1 | aplikasi |
| `LINKSERVICE` | 1 | **bukan aplikasi** → OQ-009 |

Rule di class bukan-aplikasi: **27** dari 270.

## 3. Daftar rule per tipe

### Activity — 112 rule (`RULE-OBJ-ACTIVITY`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AddAdjustment_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `ADDADJUSTMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / ADDESTIMATION_ACT` |
| `AddEstimation_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `ADDESTIMATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | — |
| `AddInterest_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `ADDINTEREST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `AddKomiteTreatyChild_ACT.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `ADDKOMITETREATYCHILD_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=38 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / ADDCHILDANDSAVE_ACT` |
| `AddListClaimAmount.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `ADDLISTCLAIMAMOUNT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | — |
| `AddLossAllocation_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `ADDLOSSALLOCATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | — |
| `AttachmentProtect_ACT.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `ATTACHMENTPROTECT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=30 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / ATTACHMENTPROTECT_ACT` |
| `CNMInsertCauseOfLoss_act.xml` | `@BASECLASS` | `CNMINSERTCAUSEOFLOSS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / CNMINSERTSURVEYORS_ACT` |
| `CNMInsertDetailCauseOfLoss_act.xml` | `@BASECLASS` | `CNMINSERTDETAILCAUSEOFLOSS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / CNMINSERTSTSCLAIM_ACT` |
| `CNMSetDetailCauseOfLoss_act.xml` | `@BASECLASS` | `CNMSETDETAILCAUSEOFLOSS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `@BASECLASS / SETDETAILCAUSEOFLOSS_ACT` |
| `CancelTreatyGroup.xml` | `@BASECLASS` | `CANCELTREATYGROUP` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CekPolisAvailable_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `CEKPOLISAVAILABLE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CekPremiLunas_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `CEKPREMILUNAS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `ChangeData_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `CHANGEDATA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `Check.xml` | `ASM-FW-GCNMFW-DATA-ESTIMASI` | `CHECK` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CheckAnyAcceptationProp.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `CHECKANYACCEPTATIONPROP` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `CheckDateDOL_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `CHECKDATEDOL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=26 | `ASM-FW-GCNMFW-WORK-PNC / CHECKDATE_ACT` |
| `CheckDateReceived_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `CHECKDATERECEIVED_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | `ASM-FW-GCNMFW-WORK-PNC / CHECKDATERECEIVED_ACT` |
| `CheckEstimateDate_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `CHECKESTIMATEDATE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `CheckNoPolicy.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `CHECKNOPOLICY` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `CheckNopolicy_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `CHECKNOPOLICY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CheckPeriodPolicy_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `CHECKPERIODPOLICY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CheckReportDate_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `CHECKREPORTDATE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `ASM-FW-GCNMFW-WORK-PNC / CHECKDATEREPORT_ACT` |
| `CheeckNoRNM_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `CHEECKNORNM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | — |
| `CloseClaimProp.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `CLOSECLAIMPROP` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP / CLOSECLAIMTNONPROP` |
| `CopyOldataCurr_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `COPYOLDATACURR_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / COUNTTOTALINSTEREST_ACT` |
| `CountDeductible_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `COUNTDEDUCTIBLE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CountEstimation_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `COUNTESTIMATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=39 | — |
| `CountGrossAdjTreaty_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `COUNTGROSSADJTREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `CountListClaimAmountIDR.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `COUNTLISTCLAIMAMOUNTIDR` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `CountPersen_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `COUNTPERSEN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `CountSpreadingADJ_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `COUNTSPREADINGADJ_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / SETSPREADINGAJSUTEMENT_ACT` |
| `CountSpreading_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `COUNTSPREADING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=25 | — |
| `CountTotalInsterest_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `COUNTTOTALINSTEREST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=36 | — |
| `CountValueADJTreaty_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `COUNTVALUEADJTREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=23 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / SETNILAIRESIKOSENDIRI` |
| `CountValue_Act.xml` | `ASM-FW-GISFW-DATA-TREATYINTOTAL` | `COUNTVALUE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GCNMFW-DATA-OBJECTITEM / COUNTVALUE_ACT` |
| `CreatClaimAnalysis_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `CREATCLAIMANALYSIS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | — |
| `CurencyEstimation_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `CURENCYESTIMATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `DeleteAjsutment_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `DELETEAJSUTMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `DeleteDataTreatyGroup.xml` | `@BASECLASS` | `DELETEDATATREATYGROUP` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `DeleteEstimation_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `DELETEESTIMATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=18 | — |
| `DeleteInterest_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `DELETEINTEREST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | — |
| `DeleteListClaim_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `DELETELISTCLAIM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `EditMstConsultant_Act.xml` | `ASM-FW-GISFW-INT-ADJUSTERCONSULTANT` | `EDITMSTCONSULTANT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `FilterLossAllocation_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `FILTERLOSSALLOCATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `GetAdders_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `GETADDERS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `ASM-FW-GCNMFW-WORK-PNC / INPUTKODEPOS1_ACT` |
| `GetBase64Attachment.xml` | `ASM-FW-GCNMFW-WORK` | `GETBASE64ATTACHMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | `WORK- / LOADATTACHMENTDATA` |
| `GetDataOustanding.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `GETDATAOUSTANDING` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `GetDataOutsClaim_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `GETDATAOUTSCLAIM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `GetDetailPolis_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `GETDETAILPOLIS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `GetDtlPaymentClaimTreatyin_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `GETDTLPAYMENTCLAIMTREATYIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `ASM-FW-GCNMFW-DATA-OBJECTITEM / GETPAYMENTCLAIM_ACT` |
| `GetDtlPaymentPremi_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `GETDTLPAYMENTPREMI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `ASM-FW-GCNMFW-WORK-PNC / GETPAYMENTLIST_ACT` |
| `GetInvoiceAttachments.xml` | `ASM-FW-GCNMFW-WORK` | `GETINVOICEATTACHMENTS` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `GetLinkService.xml` | `ASM-FW-GISFW-INT-M_LINK_SERVICE` | `GETLINKSERVICE` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetMasterTreaty_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `GETMASTERTREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `GetNameCauseofLoss_Act.xml` | `ASM-FW-GISFW-INT-CLAUSE` | `GETNAMECAUSEOFLOSS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `GetPICAdjutment_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `GETPICADJUTMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `GetPayAttachmentAdj_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GETPAYATTACHMENTADJ_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `GetPayAttachment_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `GETPAYATTACHMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `GetRNMShareTreaty.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `GETRNMSHARETREATY` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetReportStatus_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `GETREPORTSTATUS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `ASM-FW-GCNMFW-WORK-PNC / GETREPORTSTATUS_ACT` |
| `GetStatusKasir_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GETSTATUSKASIR_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetUrlGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETURLGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / INSERTGOOGLESTORAGE_ACT` |
| `HitServiceToKasir_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `HITSERVICETOKASIR_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=37 | — |
| `InputCatastrope.xml` | `ASM-FW-GCNMFW-WORK` | `INPUTCATASTROPE` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `InsertDocument_Act.xml` | `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | `INSERTDOCUMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `InsertGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERTGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `InsertJsonClaimTreaty_act.xml` | `ASM-FW-GCNMFW-WORK-KOMITETREATY` | `INSERTJSONCLAIMTREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GCNMFW-WORK / INSERTJSONCLAIMTREATY_ACT` |
| `InsertLogServiceClaim.xml` | `ASM-FW-GCNMFW-WORK` | `INSERTLOGSERVICECLAIM` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `KonversiKlaim_Act.xml` | `ASM-FW-GCNMFW-WORK` | `KONVERSIKLAIM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `MakeLowercase_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `MAKELOWERCASE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `NewAdjustConsult_Act.xml` | `@BASECLASS` | `NEWADJUSTCONSULT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `PrintDLATreatyIn.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `PRINTDLATREATYIN` | [terverifikasi] pyActivityType=ACTIVITY; steps=38 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / GENERATEDLATREATYIN` |
| `PrintFileAcceptance.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `PRINTFILEACCEPTANCE` | [terverifikasi] pyActivityType=ACTIVITY; steps=16 | — |
| `ProteksiData_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `PROTEKSIDATA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `ProteksiInitialandDate_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `PROTEKSIINITIALANDDATE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `RemoveLossAlloction_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `REMOVELOSSALLOCTION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SaveAcceptationTreaty_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SAVEACCEPTATIONTREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=21 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / SETSAVEACCEPTATION` |
| `SaveAdjusterConsultant_Act.xml` | `@BASECLASS` | `SAVEADJUSTERCONSULTANT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `SaveCatasrtope_Act.xml` | `ASM-FW-GCNMFW-WORK` | `SAVECATASRTOPE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SaveOutstanding_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SAVEOUTSTANDING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=49 | — |
| `SendCloseClaimToKomite.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SENDCLOSECLAIMTOKOMITE` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | `ASM-FW-GCNMFW-WORK-PNC / SENDCLOSECLAIMTOKOMITE` |
| `SendEmailKlaim.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SENDEMAILKLAIM` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `SendEmailKlaimRejectClose.xml` | `ASM-FW-GCNMFW-WORK` | `SENDEMAILKLAIMREJECTCLOSE` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | — |
| `SetAdjsuter_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SETADJSUTER_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GCNMFW-WORK-PNC / SETADJSUTER_ACT` |
| `SetCatastrope_act.xml` | `ASM-FW-GCNMFW-INT-CATASTROPHE` | `SETCATASTROPE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetCauseOfLossValue_act.xml` | `@BASECLASS` | `SETCAUSEOFLOSSVALUE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `@BASECLASS / SETSURVERYORSVALUE_ACT` |
| `SetConsultant_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SETCONSULTANT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GCNMFW-WORK-PNC / SETCONSULTANT_ACT` |
| `SetCurencyInterest_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SETCURENCYINTEREST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / CURENCYESTIMATION_ACT` |
| `SetCurencyList_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SETCURENCYLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `SetCurrency_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SETCURRENCY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SetDLACedingSOB.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETDLACEDINGSOB` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetDefNonCatastrope_Act.xml` | `ASM-FW-GCNMFW-WORK` | `SETDEFNONCATASTROPE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetEditCatastrope.xml` | `ASM-FW-GCNMFW-WORK` | `SETEDITCATASTROPE` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetEndDate_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SETENDDATE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `SetFormat_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SETFORMAT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetKomiteNo_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETKOMITENO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetKomiteTreaty_ACT.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETKOMITETREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / SETKOMITELIST_ACT` |
| `SetMOClaimTreaty_Act.xml` | `ASM-FW-GCNMFW-WORK` | `SETMOCLAIMTREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GCNMFW-WORK-PNC / SETMOCLAIM_ACT` |
| `SetMasterID.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SETMASTERID` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetNameCurrency_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETNAMECURRENCY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=20 | — |
| `SetNameTreaty_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SETNAMETREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `SetOutstanding_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SETOUTSTANDING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SetPayableTo_act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETPAYABLETO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SetPayableTreaty_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETPAYABLETREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / SETPAYABLE_ACT` |
| `SetTreatyNameSpreading_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SETTREATYNAMESPREADING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=18 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / SETNAMETREATY_ACT` |
| `SetValueToClaim_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SETVALUETOCLAIM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=18 | `ASM-FW-GISFW-INT-TREATYINDETAIL / SETVALUETOCLAIM_ACT` |
| `SetViewDataMasterTreaty_Act.xml` | `DATA-PORTAL` | `SETVIEWDATAMASTERTREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SethistoryKlaimTreaty.xml` | `ASM-FW-GCNMFW-WORK` | `SETHISTORYKLAIMTREATY` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | — |
| `TryMakePLA_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `TRYMAKEPLA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=36 | — |
| `UpdateTableOS.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `UPDATETABLEOS` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `getStatusKonversi_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GETSTATUSKONVERSI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |

### ConnectREST — 6 rule (`RULE-CONNECT-REST`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `GetDtlPaymentClaim.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `GETDTLPAYMENTCLAIM` | [terverifikasi] pyServiceName=GetDtlPaymentClaim; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | `ASM-FW-GCNMFW-WORK-PNC / GETPAYMENTCLAIM` |
| `KonversiKlaimNonLife.xml` | `ASM-FW-GCNMFW-WORK` | `KONVERSIKLAIMNONLIFE` | [terverifikasi] pyServiceName=KonversiKlaimNonLife; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | `ASM-FW-GCNMFW-WORK-PNC / INSERTREJECTCLAIM` |
| `SendAcceptationToKasir.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SENDACCEPTATIONTOKASIR` | [terverifikasi] pyServiceName=SendAcceptationToKasir; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | — |
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `SERVICEGOOGLE` | [terverifikasi] pyServiceName=ServiceGoogle; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / GOOGLESTORAGE_UPLOAD` |
| `getPayAttachment.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GETPAYATTACHMENT` | [terverifikasi] pyServiceName=getPayAttachment; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | — |
| `getPremiumPaidOnTreatyIn.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `GETPREMIUMPAIDONTREATYIN` | [terverifikasi] pyServiceName=getPremiumPaidOnTreatyIn; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / GETPREMIUMPAIDON` |

### DataTransform — 11 rule (`RULE-OBJ-MODEL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `CNMRefreshCauseOfLoss_dt.xml` | `@BASECLASS` | `CNMREFRESHCAUSEOFLOSS_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / CNMREFRESHSURVEYORS_DT` |
| `CNMRefreshListDetailCauseOfLoss_dt.xml` | `@BASECLASS` | `CNMREFRESHLISTDETAILCAUSEOFLOSS_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CNMShowInsertCauseOfLoss_dt.xml` | `@BASECLASS` | `CNMSHOWINSERTCAUSEOFLOSS_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / CNMSHOWINSERTSUVERYORS_DT` |
| `CNMShowInsertDetailCauseOfLoss_dt.xml` | `@BASECLASS` | `CNMSHOWINSERTDETAILCAUSEOFLOSS_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / CNMINSERTDETAILCAUSEOFLOSS_DT` |
| `CekInterestListDtl_DT.xml` | `ASM-FW-GISFW-DATA-TREATYINTOTAL` | `CEKINTERESTLISTDTL_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CloseStsSaveTreatyGroup.xml` | `ASM-FW-GISFW-INT-TREATYGROUP` | `CLOSESTSSAVETREATYGROUP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DisableEditRNMShare.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `DISABLEEDITRNMSHARE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InsertChronology_DT.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `INSERTCHRONOLOGY_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-PNC / INSERTCHRONOLOGY_DT` |
| `SetDateAcceptation.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SETDATEACCEPTATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetDateOutstanding.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SETDATEOUTSTANDING` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetLabel_Dt.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SETLABEL_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-PNC / SETLABEL_DT` |

### DecisionTable — 1 rule (`RULE-DECLARE-DECISIONTABLE`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `GetMimeType.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETMIMETYPE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Flow — 1 rule (`RULE-OBJ-FLOW`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `Flow_TreatyIn.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `FLOW_TREATYIN` | [terverifikasi] pyStartActivity=Start1 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / REGISTER_FLOW` |

### FlowAction — 12 rule (`RULE-OBJ-FLOWACTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AdjustmentDetail.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `ADJUSTMENTDETAIL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CatastrofeList.xml` | `ASM-FW-GCNMFW-WORK` | `CATASTROFELIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GenerateDLATreaty.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GENERATEDLATREATY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GeneratePLA.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `GENERATEPLA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputAcceptation.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `INPUTACCEPTATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / INPUTESTIMASI` |
| `MessageBeforeDeleteTreatyGroup.xml` | `ASM-FW-GISFW-INT-TREATYGROUP` | `MESSAGEBEFOREDELETETREATYGROUP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `OutstandingClaim.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `OUTSTANDINGCLAIM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / VIEWPOLIS` |
| `PreventRejectClaimProp.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `PREVENTREJECTCLAIMPROP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PrintFile.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `PRINTFILE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PrintFileDLA.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `PRINTFILEDLA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewDetailInterest.xml` | `ASM-FW-GISFW-DATA-TREATYINTOTAL` | `VIEWDETAILINTEREST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewDetailPayment_FA.xml` | `ASM-FW-GCNMFW-DATA-CURRENCY` | `VIEWDETAILPAYMENT_FA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Harness — 11 rule (`RULE-HTML-HARNESS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `CauseofLoss_Harness.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `CAUSEOFLOSS_HARNESS` | [terverifikasi] pyInclude=1 | `ASM-FW-GCNMFW-WORK-PNC / CAUSEOFLOSS_HARNESS` |
| `CommitteeTreaty.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `COMMITTEETREATY` | [terverifikasi] pyInclude=1 | — |
| `ListPaymentClaim_Harness.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `LISTPAYMENTCLAIM_HARNESS` | [terverifikasi] pyInclude=1 | — |
| `ListPaymentPremi_Harness.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `LISTPAYMENTPREMI_HARNESS` | [terverifikasi] pyInclude=1 | — |
| `ListPolicyNoTreaty_Harness.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `LISTPOLICYNOTREATY_HARNESS` | [terverifikasi] pyInclude=1 | — |
| `MasterTreatyIn.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `MASTERTREATYIN` | [terverifikasi] pyInclude=1 | — |
| `MstAdjusterConsultant.xml` | `DATA-PORTAL` | `MSTADJUSTERCONSULTANT` | [terverifikasi] pyInclude=1 | — |
| `SummaryOutSClaim.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SUMMARYOUTSCLAIM` | [terverifikasi] pyInclude=1 | — |
| `TambahCauseofLoss.xml` | `DATA-PORTAL` | `TAMBAHCAUSEOFLOSS` | [terverifikasi] pyInclude=3 | — |
| `TambahMasterCauseOfLoss.xml` | `DATA-PORTAL` | `TAMBAHMASTERCAUSEOFLOSS` | [terverifikasi] pyInclude=3 | — |
| `ViewAttachment.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `VIEWATTACHMENT` | [terverifikasi] pyInclude=1 | — |

### RDBList — 54 rule (`RULE-CONNECT-SQL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseRW_SQL.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `BROWSERW_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CITY,RW | — |
| `CariHistoryClaim_SQL.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `CARIHISTORYCLAIM_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_KLAIM | — |
| `CekLunasPremi_Sql.xml` | `ASM-FW-GCNMFW-WORK` | `CEKLUNASPREMI_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=ARASAPAS.INVOICE | — |
| `CekPolicyNumber_SQL.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `CEKPOLICYNUMBER_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TREATYINPRODUCTION | — |
| `CekProteksiKlaim.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `CEKPROTEKSIKLAIM` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.OPENPROTEKSI_EDM | — |
| `CurrencyStandard.xml` | `ASM-FW-GISFW-INT-CURRENCYSTANDARD` | `CURRENCYSTANDARD` | [terverifikasi] sqlKind=QUERY; procs=POOLDATA.GETCURRENCYSTANDARD; sqlOps=SELECT | `ASM-FW-GISFW-INT-CURRENCYSTANDARD / ASM!UPDATEMASTERCURRENCYSTANDARD` |
| `DataOutstandingTreatyin.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `DATAOUTSTANDINGTREATYIN` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=OS_AKSEPTASI_KLAIM | — |
| `DeleteDataTreatyGroup_SQL.xml` | `ASM-FW-GISFW-INT-TREATYGROUP` | `DELETEDATATREATYGROUP_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=DELETE; tables=M_TREATYGROUP | — |
| `GETTanggalClosing_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETTANGGALCLOSING_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TANGGAL_CLOSING | — |
| `GenerateImageID_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GENERATEIMAGEID_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GenerateNoCLMTreatyIn.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `GENERATENOCLMTREATYIN` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.GENERATE_NOCLMTREATYIN | `ASM-FW-GCNMFW-INT-V_POLIS / GCNM!GENERATENOPLATREATYIN` |
| `GenerateNoTRTInTemp.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `GENERATENOTRTINTEMP` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.GENERATE_NOCLMTRTYINTEMP | — |
| `GetAddressCeding.xml` | `ASM-FW-GISFW-INT-CLIENT` | `GETADDRESSCEDING` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_CLIENT,AGENT | — |
| `GetAddressTreatyIn.xml` | `ASM-FW-GISFW-INT-CLIENT` | `GETADDRESSTREATYIN` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_CLIENT,AGENT | `ASM-FW-GISFW-INT-CLIENT / ASM!GETADDRESSCEDING` |
| `GetAppName_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETAPPNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.T_FOLDER_IMAGE | — |
| `GetCaseIDNBTretyIn.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `GETCASEIDNBTRETYIN` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | `ASM-FW-GCNMFW-INT-V_POLIS / ASM!GETCASEIDNB` |
| `GetCurrency.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `GETCURRENCY` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CURRENCY | `ASM-FW-GISFW-INT-CURRENCY / ASM!UPDATEMASTERCURRENCY` |
| `GetDataBankAccount2_sql.xml` | `ASM-FW-GISFW-INT-BANKACCOUNT` | `GETDATABANKACCOUNT2_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=BANKACCOUNT | `ASM-FW-GISFW-INT-BANKACCOUNT / GCNM!GETDATABANKACCOUNT_SQL` |
| `GetDataBankAccount_sql.xml` | `ASM-FW-GISFW-INT-BANKACCOUNT` | `GETDATABANKACCOUNT_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=BANKACCOUNT | — |
| `GetDataKlaimTreatyin.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETDATAKLAIMTREATYIN` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=OS_AKSEPTASI_KLAIM | — |
| `GetDatabyClientName.xml` | `ASM-FW-GISFW-INT-BANKACCOUNT` | `GETDATABYCLIENTNAME` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=BANKACCOUNT | — |
| `GetEmailCeding_SQL.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GETEMAILCEDING_SQL` | [terverifikasi] sqlKind=QUERY; procs=GL.F_GET_EMAIL; sqlOps=SELECT | — |
| `GetIDConsultanAdj_SQL.xml` | `ASM-FW-GISFW-INT-ADJUSTERCONSULTANT` | `GETIDCONSULTANADJ_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.M_SITE_DATABASE | `ASM-FW-GISFW-INT-ADJUSTERCONSULTANT / RNM!SAVEMASTERADJUSTERCONSULTANT_SQL` |
| `GetKodeProdNonLife_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETKODEPRODNONLIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.KODE_PRODUKSI | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETSEQUENCENUMBER_SQL` |
| `GetLBUID_SQL.xml` | `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS_BUSINESS` | `GETLBUID_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=V_D_CAUSE_OF_LOSS_BUSINESS | `ASM-FW-GCNMFW-INT-V_M_CAUSE_OF_LOSS_BUSINESS / GCNM!GETLBUID_SQL` |
| `GetLeaderReport.xml` | `ASM-FW-GISFW-INT-AGENT` | `GETLEADERREPORT` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=AGENT | `ASM-FW-GISFW-INT-CLIENT / ASM!GETADDRESSCEDING` |
| `GetLimitDirekturUtama_SQL.xml` | `ASM-FW-GCNMFW-INT-EMAILKOMITE` | `GETLIMITDIREKTURUTAMA_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.EMAILKOMITE | — |
| `GetLimitPLATreatyin.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLIMITPLATREATYIN` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=PROPORTIONALARRG | — |
| `GetLimitsTreatyIn_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `GETLIMITSTREATYIN_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_TREATY_IN,M_TREATY_IN_EDM | `ASM-FW-GISFW-INT-OFFERJSON / ASM!GETCOPYNBFORCLAIM` |
| `GetLinkStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETLINKSTORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GETTOKENSTORAGE_SQL` |
| `GetListRetro_Sql.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETLISTRETRO_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYREINSURER | — |
| `GetMObyNopol_SQL.xml` | `ASM-FW-GISFW-INT-MARKETINGOFFICER` | `GETMOBYNOPOL_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=MARKETINGOFFICER,TREATYINPRODUCTION | — |
| `GetNopolis_SQL.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `GETNOPOLIS_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYINPRODUCTION | — |
| `GetOldIDBusiness.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `GETOLDIDBUSINESS` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=BUSINESS | — |
| `GetPolicyData.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `GETPOLICYDATA` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | — |
| `GetSequenceNumber_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETSEQUENCENUMBER_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETTANGGALCLOSING_SQL` |
| `GetStatusKasir_SQL.xml` | `ASM-FW-GCNMFW-WORK` | `GETSTATUSKASIR_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.DIRECTTOKASIR_LOG | — |
| `GetTokenStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETTOKENSTORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.GET_TOKEN_STORAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |
| `GetTreatyGroupID.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETTREATYGROUPID` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYBUSINESS | — |
| `GetTreatyInMasterProp_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `GETTREATYINMASTERPROP_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TREATYINDETAIL,POOLDATA.TREATYBUSINESS,POOLDATA.TREATYINDETAILEDM | `ASM-FW-GISFW-INT-OFFERJSON / ASM!GETTREATYINMASTER_SQL` |
| `GetYearofQuartal.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETYEAROFQUARTAL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | — |
| `InsertClaimPNC.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `INSERTCLAIMPNC` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_KLAIM_PNC | — |
| `InsertLOGDirectKasir_SQL.xml` | `ASM-FW-GCNMFW-WORK` | `INSERTLOGDIRECTKASIR_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.DIRECTTOKASIR_LOG | `ASM-FW-GISFW-WORK / RNM!INSERTLOGMOP_SQL` |
| `InsertLogServiceClaim.xml` | `ASM-FW-GCNMFW-WORK` | `INSERTLOGSERVICECLAIM` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.MONITORING_KLAIM_LOG | — |
| `Insert_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERT_T_STORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GENERATEIMAGEID_SQL` |
| `SaveDataToOsAkseptasiNP.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `SAVEDATATOOSAKSEPTASINP` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP | `ASM-FW-GCNMFW-INT-V_POLIS / GCNM!SAVEDATATOOSAKSEPTASI` |
| `SaveOSClaim_SQL.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `SAVEOSCLAIM_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM | — |
| `SaveOSKlaimTreaty_SQL.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `SAVEOSKLAIMTREATY_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTRT | `ASM-FW-GCNMFW-INT-V_POLIS / ASM!SAVEOSCLAIM_SQL` |
| `SetPolicyTreatyProp.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `SETPOLICYTREATYPROP` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYINPRODUCTION | — |
| `TreatyYearTreatyin_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `TREATYYEARTREATYIN_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYYEAR | — |
| `UpdateDCauseOfLoss.xml` | `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS` | `UPDATEDCAUSEOFLOSS` | [terverifikasi] sqlKind=PLSQL; procs=DBMS_LOB.CREATETEMPORARY,POOLDATA.PEGA_D_CAUSE_OF_LOSS | `ASM-FW-CNMFW-INT-V_D_CAUSE_OF_LOSS / CNM!UPDATEDCAUSEOFLOSS` |
| `UpdateMCauseOfLoss.xml` | `ASM-FW-GCNMFW-INT-V_M_CAUSE_OF_LOSS` | `UPDATEMCAUSEOFLOSS` | [terverifikasi] sqlKind=PLSQL; procs=DBMS_LOB.CREATETEMPORARY,POOLDATA.PEGA_M_CAUSE_OF_LOSS | — |
| `Update_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `UPDATE_T_STORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |
| `getStatusKonversi_SQL.xml` | `ASM-FW-GCNMFW-INT-OS_AKSEPTASI_KLAIM` | `GETSTATUSKONVERSI_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=REINSURANCE.TRLOSS_DETAIL_T | — |

### ReportDefinition — 19 rule (`RULE-OBJ-REPORT-DEFINITION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseAdjusterConsultant.xml` | `ASM-FW-GISFW-INT-ADJUSTERCONSULTANT` | `BROWSEADJUSTERCONSULTANT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseBankAccount.xml` | `ASM-FW-GISFW-INT-BANKACCOUNT` | `BROWSEBANKACCOUNT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseBusiness_RD.xml` | `ASM-FW-GISFW-INT-BUSINESS` | `BROWSEBUSINESS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCLAIM_MASTER_TREATY.xml` | `ASM-FW-GCNMFW-INT-CLAIM_MASTER_TREATY` | `BROWSECLAIM_MASTER_TREATY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCouseOfLoss_Business.xml` | `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS_BUSINESS` | `BROWSECOUSEOFLOSS_BUSINESS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCurrency_RD.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `BROWSECURRENCY_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseProvince_RD.xml` | `ASM-FW-GISFW-INT-RW` | `BROWSEPROVINCE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseRW_RD.xml` | `ASM-FW-GISFW-INT-RW` | `BROWSERW_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseReinsuranceType_RD.xml` | `ASM-FW-GISFW-INT-REINSURANCETYPE` | `BROWSEREINSURANCETYPE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyGroup_RD.xml` | `ASM-FW-GISFW-INT-TREATYGROUP` | `BROWSETREATYGROUP_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyInDetail.xml` | `ASM-FW-GISFW-INT-TREATYINDETAIL` | `BROWSETREATYINDETAIL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseVDCauseOfLoss_RD.xml` | `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS` | `BROWSEVDCAUSEOFLOSS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseVMCauseOfLoss_RD.xml` | `ASM-FW-GCNMFW-INT-V_M_CAUSE_OF_LOSS` | `BROWSEVMCAUSEOFLOSS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseVMstUserTeknis_RD.xml` | `ASM-FW-GCNMFW-INT-V_MST_USER_TEKNIS` | `BROWSEVMSTUSERTEKNIS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `FilterEmailKomiteWithLimit.xml` | `ASM-FW-GCNMFW-INT-EMAILKOMITE` | `FILTEREMAILKOMITEWITHLIMIT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GetCatastrope_RD.xml` | `ASM-FW-GCNMFW-INT-CATASTROPHE` | `GETCATASTROPE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `RejectedClaim_RD.xml` | `ASM-FW-GCNMFW-INT-CLAIMREJECTED` | `REJECTEDCLAIM_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SelectVDCauseOfLoss_RD.xml` | `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS` | `SELECTVDCAUSEOFLOSS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-CNMFW-INT-V_D_CAUSE_OF_LOSS / SELECTVDCAUSEOFLOSS_RD` |
| `SelectVMCauseOfLoss_RD.xml` | `ASM-FW-GCNMFW-INT-V_M_CAUSE_OF_LOSS` | `SELECTVMCAUSEOFLOSS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Section — 36 rule (`RULE-HTML-SECTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AdjustmentDetail.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `ADJUSTMENTDETAIL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `AdjustmentDetail_Section.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `ADJUSTMENTDETAIL_SECTION` | [terverifikasi] pyInclude=6 | — |
| `AdjustmentType.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `ADJUSTMENTTYPE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCauseOfLoss.xml` | `@BASECLASS` | `BROWSECAUSEOFLOSS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseDetailCauseOfLoss.xml` | `@BASECLASS` | `BROWSEDETAILCAUSEOFLOSS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CatastrofeList_Sec.xml` | `ASM-FW-GCNMFW-WORK` | `CATASTROFELIST_SEC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Catastrope_Sec.xml` | `ASM-FW-GCNMFW-WORK` | `CATASTROPE_SEC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CauseofLoss_Section.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `CAUSEOFLOSS_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-PNC / CAUSEOFLOSS_SECTION` |
| `ComiteeClaimTreaty.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `COMITEECLAIMTREATY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DtlDataCommitte.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `DTLDATACOMMITTE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / COMITEECLAIMTREATY` |
| `GeneratePLA.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `GENERATEPLA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridCauseOfLoss.xml` | `DATA-PORTAL` | `GRIDCAUSEOFLOSS` | [terverifikasi] pyInclude=2 | — |
| `InputAcceptation.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `INPUTACCEPTATION` | [terverifikasi] pyInclude=10 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / INPUTESTIMASIADMIN` |
| `InputAcceptation_Adjs.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `INPUTACCEPTATION_ADJS` | [terverifikasi] pyInclude=1 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / INPUTACCEPTATION_EST` |
| `InputAcceptation_Est.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `INPUTACCEPTATION_EST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / INPUTACCEPTATION_INTRS` |
| `InputAcceptation_Intrs.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `INPUTACCEPTATION_INTRS` | [terverifikasi] pyInclude=2 | — |
| `IsDeleteTreatyGroup.xml` | `ASM-FW-GISFW-INT-TREATYGROUP` | `ISDELETETREATYGROUP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ListDetailCauseOfLoss.xml` | `DATA-PORTAL` | `LISTDETAILCAUSEOFLOSS` | [terverifikasi] pyInclude=2 | — |
| `ListPaymentClaim_SC.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `LISTPAYMENTCLAIM_SC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-DATA-OBJECTITEM / LISTPAYMENTCLAIM_SC` |
| `ListPolicyNoTreaty.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `LISTPOLICYNOTREATY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / MASTERTREATYINLIST` |
| `MasterTreatyInList.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `MASTERTREATYINLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `MstAdjusterConsultant.xml` | `DATA-PORTAL` | `MSTADJUSTERCONSULTANT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `OutstandingClaim.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `OUTSTANDINGCLAIM` | [terverifikasi] pyInclude=8 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / VIEWPOLIS` |
| `OutstandingClaim_Est.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `OUTSTANDINGCLAIM_EST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / OUTSTANDINGCLAIM_INTRS` |
| `OutstandingClaim_Intrs.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `OUTSTANDINGCLAIM_INTRS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `OutstandingClaim_Sprd.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `OUTSTANDINGCLAIM_SPRD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / OUTSTANDINGCLAIM_EST` |
| `PaymentPremiList_Section.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `PAYMENTPREMILIST_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-PNC / PAYMENTPREMILIST_SC` |
| `PreventRejectClaimProp.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `PREVENTREJECTCLAIMPROP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-PNC / PREVENTREJECTCLAIM` |
| `PrintFile.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `PRINTFILE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PrintFileDLA.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `PRINTFILEDLA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / PRINTFILEDLA` |
| `Subjectivity.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SUBJECTIVITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / DTLDATACOMMITTE` |
| `SummaryOutsClaim.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SUMMARYOUTSCLAIM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewAttachment.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `VIEWATTACHMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-DATA-OBJECTITEM / VIEWATTACHMENT` |
| `ViewDetailInterest.xml` | `ASM-FW-GISFW-DATA-TREATYINTOTAL` | `VIEWDETAILINTEREST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewDetailPayment_SC.xml` | `ASM-FW-GCNMFW-DATA-CURRENCY` | `VIEWDETAILPAYMENT_SC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `generateDLATreaty.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GENERATEDLATREATY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / PRINTDLA` |

### SystemSettings — 1 rule (`RULE-ADMIN-SYSTEM-SETTINGS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `LinkService.xml` | `LINKSERVICE` | `LINKSERVICE` | [terverifikasi] pyPurpose=LinkService | — |

### When — 6 rule (`RULE-OBJ-WHEN`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `IsBackStage.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `ISBACKSTAGE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-PNC / ISBACKSTAGE` |
| `IsCLM.xml` | `@BASECLASS` | `ISCLM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCLMNP.xml` | `@BASECLASS` | `ISCLMNP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / ISCLMP` |
| `IsCLMP.xml` | `@BASECLASS` | `ISCLMP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / ISCLM` |
| `IsPEGAPROD.xml` | `@BASECLASS` | `ISPEGAPROD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `WORK- / ISPEGAPROD` |
| `IsPEGASyariah.xml` | `@BASECLASS` | `ISPEGASYARIAH` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

## 4. Tabrakan nama di dalam modul ini

### 4.1 Nama file dipakai di lebih dari satu tipe rule

| Nama file | Tipe | Class |
| --- | --- | --- |
| `ViewDetailInterest.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINTOTAL, ASM-FW-GISFW-DATA-TREATYINTOTAL` |
| `OutstandingClaim.xml` | FlowAction, Section | `ASM-FW-GCNMFW-WORK-CLAIMTREATY, ASM-FW-GCNMFW-WORK-CLAIMTREATY` |
| `InsertLogServiceClaim.xml` | Activity, RDBList | `ASM-FW-GCNMFW-WORK, ASM-FW-GCNMFW-WORK` |
| `GeneratePLA.xml` | FlowAction, Section | `ASM-FW-GCNMFW-WORK-CLAIMTREATY, ASM-FW-GCNMFW-WORK-CLAIMTREATY` |
| `PreventRejectClaimProp.xml` | FlowAction, Section | `ASM-FW-GCNMFW-WORK-CLAIMTREATY, ASM-FW-GCNMFW-WORK-CLAIMTREATY` |
| `ViewAttachment.xml` | Harness, Section | `ASM-FW-GCNMFW-WORK-CLAIMTREATY, ASM-FW-GCNMFW-WORK-CLAIMTREATY` |
| `MstAdjusterConsultant.xml` | Harness, Section | `DATA-PORTAL, DATA-PORTAL` |
| `AdjustmentDetail.xml` | FlowAction, Section | `ASM-FW-GCNMFW-DATA-ADJUSTMENT, ASM-FW-GCNMFW-DATA-ADJUSTMENT` |
| `InputAcceptation.xml` | FlowAction, Section | `ASM-FW-GCNMFW-WORK-CLAIMTREATY, ASM-FW-GCNMFW-WORK-CLAIMTREATY` |
| `PrintFileDLA.xml` | FlowAction, Section | `ASM-FW-GCNMFW-WORK-CLAIMTREATY, ASM-FW-GCNMFW-WORK-CLAIMTREATY` |
| `PrintFile.xml` | FlowAction, Section | `ASM-FW-GCNMFW-WORK-CLAIMTREATY, ASM-FW-GCNMFW-WORK-CLAIMTREATY` |

Total: **11** nama file.

### 4.2 Rule `When` bernama sama di lebih dari satu class

_Tidak ada._

## 5. Keluarga klon (satu asal salinan dipakai >1 rule)

Dibaca dari `<pzOriginalInstanceKey>` yang sama pada beberapa file. Ini **sinyal kandidat
konsolidasi** saat migrasi, bukan bukti bahwa perilakunya identik — kesamaan perilaku baru
dapat dinyatakan setelah isinya dibandingkan di STEP D2.

| Asal salinan (class / nama) | Jumlah | File |
| --- | ---: | --- |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATY / SETNAMETREATY_ACT` | 2 | `Activity/SetNameTreaty_Act.xml`, `Activity/SetTreatyNameSpreading_Act.xml` |
| `ASM-FW-GISFW-INT-CLIENT / ASM!GETADDRESSCEDING` | 3 | `RDBList/GetAddressCeding.xml`, `RDBList/GetAddressTreatyIn.xml`, `RDBList/GetLeaderReport.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATY / CURENCYESTIMATION_ACT` | 2 | `Activity/CurencyEstimation_Act.xml`, `Activity/SetCurencyInterest_act.xml` |
| `ASM-FW-GCNMFW-DATA-ADJUSTMENT / COMITEECLAIMTREATY` | 2 | `Section/ComiteeClaimTreaty.xml`, `Section/DtlDataCommitte.xml` |
| `DATA-PORTAL / MSTADJUSTERCONSULTANT` | 2 | `Harness/MstAdjusterConsultant.xml`, `Section/MstAdjusterConsultant.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATY / VIEWPOLIS` | 2 | `FlowAction/OutstandingClaim.xml`, `Section/OutstandingClaim.xml` |
| `ASM-FW-GCNMFW-INT-V_POLIS / ASM!SAVEOSCLAIM_SQL` | 2 | `RDBList/SaveOSClaim_SQL.xml`, `RDBList/SaveOSKlaimTreaty_SQL.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATY / ADDESTIMATION_ACT` | 2 | `Activity/AddAdjustment_Act.xml`, `Activity/AddEstimation_Act.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATY / GENERATEPLA` | 2 | `FlowAction/GeneratePLA.xml`, `Section/GeneratePLA.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATY / COUNTTOTALINSTEREST_ACT` | 2 | `Activity/CopyOldataCurr_act.xml`, `Activity/CountTotalInsterest_Act.xml` |
| `@BASECLASS / ISCLM` | 2 | `When/IsCLM.xml`, `When/IsCLMP.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / INSERTGOOGLESTORAGE_ACT` | 2 | `Activity/GetUrlGoogleStorage_Act.xml`, `Activity/InsertGoogleStorage_Act.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATY / INPUTACCEPTATION_INTRS` | 2 | `Section/InputAcceptation_Est.xml`, `Section/InputAcceptation_Intrs.xml` |
| `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETTANGGALCLOSING_SQL` | 2 | `RDBList/GETTanggalClosing_SQL.xml`, `RDBList/GetSequenceNumber_SQL.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATY / SUMMARYOUTSCLAIM` | 2 | `Harness/SummaryOutSClaim.xml`, `Section/SummaryOutsClaim.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GENERATEIMAGEID_SQL` | 2 | `RDBList/GenerateImageID_SQL.xml`, `RDBList/Insert_T_Storage_SQL.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATY / PRINTFILE` | 2 | `FlowAction/PrintFile.xml`, `Section/PrintFile.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` | 2 | `RDBList/GetTokenStorage_SQL.xml`, `RDBList/Update_T_Storage_SQL.xml` |
| `ASM-FW-GCNMFW-DATA-ADJUSTMENT / ADJUSTMENTDETAIL` | 2 | `FlowAction/AdjustmentDetail.xml`, `Section/AdjustmentDetail.xml` |
| `ASM-FW-GISFW-INT-BANKACCOUNT / GCNM!GETDATABANKACCOUNT_SQL` | 2 | `RDBList/GetDataBankAccount2_sql.xml`, `RDBList/GetDataBankAccount_sql.xml` |
| `ASM-FW-GISFW-DATA-TREATYINTOTAL / VIEWDETAILINTEREST` | 2 | `FlowAction/ViewDetailInterest.xml`, `Section/ViewDetailInterest.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATY / MASTERTREATYINLIST` | 2 | `Section/ListPolicyNoTreaty.xml`, `Section/MasterTreatyInList.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATY / OUTSTANDINGCLAIM_INTRS` | 2 | `Section/OutstandingClaim_Est.xml`, `Section/OutstandingClaim_Intrs.xml` |

## 6. Objek database yang disentuh

Diambil dari `<pyBrowseSQL>` pada rule `RDBList`. **Nama objek apa adanya dari SQL** —
kualifikasi skema tidak konsisten di korpus. **Tipe kolom, constraint, dan relasi tidak
diketahui** (OQ-001).

| Tabel/view (apa adanya di SQL) | Dirujuk oleh |
| --- | ---: |
| `AGENT` | 3 rule |
| `BANKACCOUNT` | 3 rule |
| `JSON_POLIS` | 3 rule |
| `TREATYINPRODUCTION` | 3 rule |
| `T_STORAGE_IMAGE` | 3 rule |
| `M_CLIENT` | 2 rule |
| `OS_AKSEPTASI_KLAIM` | 2 rule |
| `POOLDATA.DIRECTTOKASIR_LOG` | 2 rule |
| `ARASAPAS.INVOICE` | 1 rule |
| `BUSINESS` | 1 rule |
| `CITY` | 1 rule |
| `CURRENCY` | 1 rule |
| `JSON_KLAIM` | 1 rule |
| `MARKETINGOFFICER` | 1 rule |
| `M_TREATYGROUP` | 1 rule |
| `M_TREATY_IN` | 1 rule |
| `M_TREATY_IN_EDM` | 1 rule |
| `POOLDATA.EMAILKOMITE` | 1 rule |
| `POOLDATA.KODE_PRODUKSI` | 1 rule |
| `POOLDATA.MONITORING_KLAIM_LOG` | 1 rule |
| `POOLDATA.M_SITE_DATABASE` | 1 rule |
| `POOLDATA.OPENPROTEKSI_EDM` | 1 rule |
| `POOLDATA.TANGGAL_CLOSING` | 1 rule |
| `POOLDATA.TREATYBUSINESS` | 1 rule |
| `POOLDATA.TREATYINDETAIL` | 1 rule |
| `POOLDATA.TREATYINDETAILEDM` | 1 rule |
| `POOLDATA.TREATYINPRODUCTION` | 1 rule |
| `POOLDATA.T_FOLDER_IMAGE` | 1 rule |
| `PROPORTIONALARRG` | 1 rule |
| `REINSURANCE.TRLOSS_DETAIL_T` | 1 rule |
| `RW` | 1 rule |
| `TREATYBUSINESS` | 1 rule |
| `TREATYREINSURER` | 1 rule |
| `TREATYYEAR` | 1 rule |
| `V_D_CAUSE_OF_LOSS_BUSINESS` | 1 rule |

### Stored procedure / function yang dipanggil

**Body prosedur ini TIDAK ADA di korpus** — logika di dalamnya tidak dapat direkonstruksi
dari XML. Lihat OQ-002.

| Procedure | Dipanggil oleh |
| --- | --- |
| `DBMS_LOB.CREATETEMPORARY` | `UpdateDCauseOfLoss.xml`, `UpdateMCauseOfLoss.xml` |
| `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTRT` | `SaveOSKlaimTreaty_SQL.xml` |
| `POOLDATA.GENERATE_NOCLMTREATYIN` | `GenerateNoCLMTreatyIn.xml` |
| `POOLDATA.GETCURRENCYSTANDARD` | `CurrencyStandard.xml` |
| `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM` | `SaveOSClaim_SQL.xml` |
| `POOLDATA.PEGA_JSON_KLAIM_PNC` | `InsertClaimPNC.xml` |
| `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` | `GetSequenceNumber_SQL.xml` |
| `POOLDATA.PEGA_M_CAUSE_OF_LOSS` | `UpdateMCauseOfLoss.xml` |
| `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP` | `SaveDataToOsAkseptasiNP.xml` |
| `POOLDATA.GENERATE_NOCLMTRTYINTEMP` | `GenerateNoTRTInTemp.xml` |
| `POOLDATA.PEGA_D_CAUSE_OF_LOSS` | `UpdateDCauseOfLoss.xml` |
| `POOLDATA.GET_TOKEN_STORAGE` | `GetTokenStorage_SQL.xml` |
| `GL.F_GET_EMAIL` | `GetEmailCeding_SQL.xml` |

## 7. Integrasi eksternal (`ConnectREST`)

Alamat tujuan dicatat sebagai **konfigurasi**, bukan URL literal. `baseURL=SETTING` berarti
alamat diambil dari rule `SystemSettings` yang disebut di kolom Setting; nilai sebenarnya
**tidak ada di korpus**. Kolom terakhir menandai rule yang justru memuat URL literal —
itu temuan, bukan fakta arsitektur (lihat `_summary.md`).

| File | Class | pyServiceName | Sumber base URL | Setting | URL literal di dalam rule? |
| --- | --- | --- | --- | --- | --- |
| `GetDtlPaymentClaim.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `GetDtlPaymentClaim` | `SETTING` | `LinkService!LinkService` | tidak |
| `KonversiKlaimNonLife.xml` | `ASM-FW-GCNMFW-WORK` | `KonversiKlaimNonLife` | `SETTING` | `LinkService!LinkService` | tidak |
| `SendAcceptationToKasir.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SendAcceptationToKasir` | `SETTING` | `LinkService!LinkService` | tidak |
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `ServiceGoogle` | `SETTING` | `LinkService!LinkService` | tidak |
| `getPayAttachment.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `getPayAttachment` | `SETTING` | `LinkService!LinkService` | tidak |
| `getPremiumPaidOnTreatyIn.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `getPremiumPaidOnTreatyIn` | `SETTING` | `LinkService!LinkService` | tidak |

