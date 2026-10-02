# Inventaris Rule — Claim Non Prop

STEP D1, batch 3 (domain claim non-facultative). Dibangkitkan 2026-09-12 dari korpus READ-ONLY
`D:\XML\RNM_BRD\Claim Non Prop\`.

**Batas step ini:** hanya struktur + metadata. Tidak ada kesimpulan logika bisnis di sini.
Penelusuran perilaku adalah STEP D2 (`../flows/`).

**Identitas rule** diambil dari `<pxInsName>` (lihat `../README.md` §6.4).  
Kolom **Asal salinan** berisi `<pzOriginalInstanceKey>` dan hanya diisi bila **berbeda**
dari identitas — artinya rule ini hasil *Save As* dari rule lain. Tanda `—` berarti sama.
Class `—` berarti `<pxInsName>` tidak berisi bagian class (berlaku pada `DataPage`, yang
beridentitas global, bukan per class).

Perintah audit (dijalankan dari `D:\XML\RNM_BRD\`):

```
find "Claim Non Prop" -type f -name "*.xml" | wc -l
find "Claim Non Prop" -type f -name "*.xml" | awk -F/ '{print $2}' | sort | uniq -c
grep -o "<pxInsName>[^<]*" "Claim Non Prop/<tipe>/<file>.xml" | head -1
grep -o "<pzOriginalInstanceKey>[^<]*" "Claim Non Prop/<tipe>/<file>.xml" | head -1
```

## 1. Ringkasan terukur

- Total file XML: **279** — cocok dengan `../README.md` §5.1
- **Identitas rule unik: 279** dari 279 file → tidak ada identitas ganda di dalam modul ini
- Class Pega unik: **48**
- RDBList menurut jenis SQL: PLSQL=8, QUERY=42 (`PLSQL` = blok anonim `BEGIN…END;`, umumnya memanggil stored procedure)
- Prefix pada kunci RDBList: `ASM`=20, `RNM`=12, `GCNM`=18 (arti prefix **belum terverifikasi** — OQ-008)

| Tipe rule | Jumlah |
| --- | ---: |
| Activity | 114 |
| ConnectREST | 7 |
| DataTransform | 10 |
| DecisionTable | 1 |
| Flow | 1 |
| FlowAction | 16 |
| Harness | 14 |
| RDBList | 50 |
| ReportDefinition | 22 |
| Section | 36 |
| SystemSettings | 1 |
| When | 7 |
| **TOTAL** | **279** |

## 2. Sebaran class Pega

Class yang **tidak** berawalan `ASM-` adalah class bawaan Pega atau class portal — ruang
lingkup migrasinya belum ditentukan (**OQ-009**).

| Class | Jumlah rule | Aplikasi? |
| --- | ---: | --- |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | 84 | aplikasi |
| `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | 27 | aplikasi |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN` | 19 | aplikasi |
| `@BASECLASS` | 18 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GCNMFW-WORK` | 16 | aplikasi |
| `ASM-FW-GCNMFW-INT-V_POLIS` | 14 | aplikasi |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | 10 | aplikasi |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | 7 | aplikasi |
| `ASM-FW-GISFW-DATA-TREATYINLIMITLAYER` | 6 | aplikasi |
| `DATA-PORTAL` | 6 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-INT-BANKACCOUNT` | 5 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYINDETAIL` | 5 | aplikasi |
| `ASM-FW-GCNMFW-DATA-OBJECTITEM` | 3 | aplikasi |
| `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS` | 3 | aplikasi |
| `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS_BUSINESS` | 3 | aplikasi |
| `ASM-FW-GCNMFW-INT-V_M_CAUSE_OF_LOSS` | 3 | aplikasi |
| `ASM-FW-GISFW-DATA-QUOTATION` | 3 | aplikasi |
| `ASM-FW-GISFW-DATA-TREATYINLIMITS` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-CURRENCY` | 3 | aplikasi |
| `ASM-FW-GISFW-INT-OFFERJSON` | 3 | aplikasi |
| `ASM-FW-GCNMFW-DATA-CURRENCY` | 2 | aplikasi |
| `ASM-FW-GCNMFW-INT-CATASTROPHE` | 2 | aplikasi |
| `ASM-FW-GISFW-DATA-SPREADINGRISK` | 2 | aplikasi |
| `ASM-FW-GISFW-DATA-TREATYININSTALLMENT` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-AGENT` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-BUSINESS` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-CLIENT` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-POLICYJSON` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-REINSURANCETYPE` | 2 | aplikasi |
| `LINK-ATTACHMENT` | 2 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GCNMFW-INT-CLAIMREJECTED` | 1 | aplikasi |
| `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | 1 | aplikasi |
| `ASM-FW-GCNMFW-INT-EMAILKOMITE` | 1 | aplikasi |
| `ASM-FW-GCNMFW-INT-OS_AKSEPTASI_KLAIM` | 1 | aplikasi |
| `ASM-FW-GCNMFW-INT-V_MST_USER_TEKNIS` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-ADJUSTERCONSULTANT` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-CURRENCYSTANDARD` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-LST_BANK_GROUP` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-MARKETINGOFFICER` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_LINK_SERVICE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-PROPORTIONALARRG` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-PROVINCE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-RW` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYBUSINESS` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYGROUP` | 1 | aplikasi |
| `CODE-PEGA-PDF` | 1 | **bukan aplikasi** → OQ-009 |
| `LINKSERVICE` | 1 | **bukan aplikasi** → OQ-009 |
| `WORK-` | 1 | **bukan aplikasi** → OQ-009 |

Rule di class bukan-aplikasi: **29** dari 279.

## 3. Daftar rule per tipe

### Activity — 114 rule (`RULE-OBJ-ACTIVITY`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `ASMForceCaseClose.xml` | `WORK-` | `ASMFORCECASECLOSE` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | — |
| `AddAkseptasiCNP_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `ADDAKSEPTASICNP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=27 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / ADDADJUSTMENT_ACT` |
| `AddInterestListCNP_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `ADDINTERESTLISTCNP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `AddListClaimNP_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `ADDLISTCLAIMNP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | — |
| `AddLossAlocation_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `ADDLOSSALOCATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `AdjClaimAmount_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `ADJCLAIMAMOUNT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=24 | — |
| `AdjClaimCNP_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `ADJCLAIMCNP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=35 | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP / ADJCLAIMAMOUNT_ACT` |
| `AttachCAPDF_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `ATTACHCAPDF_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=18 | — |
| `AttachRISlipToWork.xml` | `CODE-PEGA-PDF` | `ATTACHRISLIPTOWORK` | [terverifikasi] pyActivityType=ACTIVITY; steps=26 | `CODE-PEGA-PDF / ATTACHTOWORK` |
| `BrowseDtlMasterTNP_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `BROWSEDTLMASTERTNP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `BrowseMasterTNP_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `BROWSEMASTERTNP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `CNMInsertCauseOfLoss_act.xml` | `@BASECLASS` | `CNMINSERTCAUSEOFLOSS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / CNMINSERTSURVEYORS_ACT` |
| `CNMInsertDetailCauseOfLoss_act.xml` | `@BASECLASS` | `CNMINSERTDETAILCAUSEOFLOSS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / CNMINSERTSTSCLAIM_ACT` |
| `CNMSetDetailCauseOfLoss_act.xml` | `@BASECLASS` | `CNMSETDETAILCAUSEOFLOSS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `@BASECLASS / SETDETAILCAUSEOFLOSS_ACT` |
| `CheckDateDOL_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `CHECKDATEDOL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=41 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / CHECKDATEDOL_ACT` |
| `CheckDateReceived_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `CHECKDATERECEIVED_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / CHECKDATERECEIVED_ACT` |
| `CheckNoPolicy.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `CHECKNOPOLICY` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / CHECKNOPOLICY` |
| `CheckPeriodPolicy_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `CHECKPERIODPOLICY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / CHECKPERIODPOLICY_ACT` |
| `CheckReportDate_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `CHECKREPORTDATE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / CHECKREPORTDATE_ACT` |
| `CloseClaimNP_preAct.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `CLOSECLAIMNP_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP / REJECTCLAIMNP_PREACT` |
| `CloseClaimTNonProp.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `CLOSECLAIMTNONPROP` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `ConcatSlipOfferNo_Act.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `CONCATSLIPOFFERNO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `CopyOldataCurr_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `COPYOLDATACURR_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / COUNTTOTALINSTEREST_ACT` |
| `CountClaimTNP_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `COUNTCLAIMTNP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=27 | — |
| `CountLossAllocation_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `COUNTLOSSALLOCATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=46 | — |
| `CountNetPremi_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTNETPREMI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `CountOverridingCommOgp_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTOVERRIDINGCOMMOGP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CountOverridingCommOnp_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTOVERRIDINGCOMMONP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CountPctInstallment_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTPCTINSTALLMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CountReinstatement_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `COUNTREINSTATEMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `CountResult1Onp_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRESULT1ONP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CountResult1_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRESULT1_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `CountResult2Ogp_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRESULT2OGP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CountResult2Onp_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRESULT2ONP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CountRiCommOgp_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRICOMMOGP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CountRiCommOnp_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `COUNTRICOMMONP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `CountSpreadingXOL.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `COUNTSPREADINGXOL` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `CountSpreading_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `COUNTSPREADING_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=22 | — |
| `CountTotalInsterest_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `COUNTTOTALINSTEREST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=36 | — |
| `CountTotalInterest_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `COUNTTOTALINTEREST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `CreateChildKomiteCNP_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `CREATECHILDKOMITECNP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=71 | — |
| `CreateChildKomiteCloseNP_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `CREATECHILDKOMITECLOSENP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | — |
| `DeleteAkseptasi_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `DELETEAKSEPTASI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / DELETEAJSUTMENT_ACT` |
| `DetailCalculation.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITLAYER` | `DETAILCALCULATION` | [terverifikasi] pyActivityType=ACTIVITY; steps=34 | — |
| `EditXOLAlokasi.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `EDITXOLALOKASI` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `FillPaymentInstallment.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `FILLPAYMENTINSTALLMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `GenerateCACNP_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GENERATECACNP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=44 | — |
| `GenerateCFS_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `GENERATECFS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=33 | — |
| `GeneratePlaCNP2_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `GENERATEPLACNP2_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP / GENERATEPLACNP_ACT` |
| `GeneratePlaCNP_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `GENERATEPLACNP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=34 | — |
| `GetAdders_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `GETADDERS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / GETADDERS_ACT` |
| `GetBase64Attachment.xml` | `ASM-FW-GCNMFW-WORK` | `GETBASE64ATTACHMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | `WORK- / LOADATTACHMENTDATA` |
| `GetDataOldAllocation.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `GETDATAOLDALLOCATION` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `GetDetailPaymentStatus_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `GETDETAILPAYMENTSTATUS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `GetDetailPolis_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `GETDETAILPOLIS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / GETDETAILPOLIS_ACT` |
| `GetHistoryMasterID_NP.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `GETHISTORYMASTERID_NP` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP / INPUTOUTSTANDINGCLMTNP_PREACT` |
| `GetInvoiceAttachments.xml` | `ASM-FW-GCNMFW-WORK` | `GETINVOICEATTACHMENTS` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `GetLinkService.xml` | `ASM-FW-GISFW-INT-M_LINK_SERVICE` | `GETLINKSERVICE` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetNameCauseofLoss_Act.xml` | `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS_BUSINESS` | `GETNAMECAUSEOFLOSS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `GetPayAttachmentAdj_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GETPAYATTACHMENTADJ_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `GetPayAttachmentNP_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `GETPAYATTACHMENTNP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / GETPAYATTACHMENT_ACT` |
| `GetReportStatus_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `GETREPORTSTATUS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / GETREPORTSTATUS_ACT` |
| `GetSelisihActual_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `GETSELISIHACTUAL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=16 | — |
| `GetUrlGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETURLGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / INSERTGOOGLESTORAGE_ACT` |
| `HTMLRISlipToPDF.xml` | `@BASECLASS` | `HTMLRISLIPTOPDF` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `@BASECLASS / HTMLTOPDF` |
| `HTMLToPDF.xml` | `@BASECLASS` | `HTMLTOPDF` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `HitServiceToKasir_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `HITSERVICETOKASIR_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=37 | — |
| `InputAkseptasi_PreAct.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `INPUTAKSEPTASI_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `InputCatastrope.xml` | `ASM-FW-GCNMFW-WORK` | `INPUTCATASTROPE` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `InputOutStandingCTNP_PostAct.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `INPUTOUTSTANDINGCTNP_POSTACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `InputOutStandingClmTNP_PreAct.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `INPUTOUTSTANDINGCLMTNP_PREACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | — |
| `InsertDocument_Act.xml` | `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM` | `INSERTDOCUMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `InsertGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERTGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `InsertJsonClaimTreaty_act.xml` | `ASM-FW-GCNMFW-WORK` | `INSERTJSONCLAIMTREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GCNMFW-WORK / INSERTJSONCLAIMNONMBU_ACT` |
| `IntIsSaveToOs.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `INTISSAVETOOS` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `KonversiKlaim_Act.xml` | `ASM-FW-GCNMFW-WORK` | `KONVERSIKLAIM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `MakeLowercase_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `MAKELOWERCASE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / MAKELOWERCASE_ACT` |
| `ProtectNilaiClaim.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `PROTECTNILAICLAIM` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / SAVECNPLAYERLIST_ACT` |
| `ProteksiNilaiClaim.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `PROTEKSINILAICLAIM` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / PROTEKSINILAICLAIM` |
| `ProteksiSendKomiteCNP_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `PROTEKSISENDKOMITECNP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=31 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / PROTEKSIINITIALANDDATE_ACT` |
| `SaveAdjustmentToOSAksep_Act_Tes.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SAVEADJUSTMENTTOOSAKSEP_ACT_TES` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / SAVECNPLAYERLIST_ACT` |
| `SaveCNPLayerList_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SAVECNPLAYERLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / SAVETOISSUERNM_ACT` |
| `SaveCatasrtope_Act.xml` | `ASM-FW-GCNMFW-WORK` | `SAVECATASRTOPE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SaveDataToJClaim_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `SAVEDATATOJCLAIM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SaveDataToOSAksep_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `SAVEDATATOOSAKSEP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=59 | — |
| `SaveToOS.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `SAVETOOS` | [terverifikasi] pyActivityType=ACTIVITY; steps=25 | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP / HITUNG_SELISIH_TEST` |
| `SendEmailKlaim.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SENDEMAILKLAIM` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | — |
| `SendEmailKlaimRejectClose.xml` | `ASM-FW-GCNMFW-WORK` | `SENDEMAILKLAIMREJECTCLOSE` | [terverifikasi] pyActivityType=ACTIVITY; steps=13 | — |
| `SetAccoutNo_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETACCOUTNO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=17 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / SETPAYABLETREATYNP_ACT` |
| `SetActualPremium_ACT.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `SETACTUALPREMIUM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `SetAdjsuter_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `SETADJSUTER_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / SETADJSUTER_ACT` |
| `SetCatastrope_act.xml` | `ASM-FW-GCNMFW-INT-CATASTROPHE` | `SETCATASTROPE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetCauseOfLossValue_act.xml` | `@BASECLASS` | `SETCAUSEOFLOSSVALUE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `@BASECLASS / SETSURVERYORSVALUE_ACT` |
| `SetConsultant_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `SETCONSULTANT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / SETCONSULTANT_ACT` |
| `SetCurencyInterest_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SETCURENCYINTEREST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / CURENCYESTIMATION_ACT` |
| `SetCurrency_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `SETCURRENCY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | — |
| `SetDataDeductible_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `SETDATADEDUCTIBLE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SetDefNonCatastrope_Act.xml` | `ASM-FW-GCNMFW-WORK` | `SETDEFNONCATASTROPE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetDueTo_act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SETDUETO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GISFW-DATA-POLICYTREATYIN / COUNTNETPREMI_ACT` |
| `SetEditCatastrope.xml` | `ASM-FW-GCNMFW-WORK` | `SETEDITCATASTROPE` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetEndDate_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `SETENDDATE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / SETENDDATE_ACT` |
| `SetFormat_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `SETFORMAT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / SETFORMAT_ACT` |
| `SetInterimXOL_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETINTERIMXOL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `SetMOClaimTreaty_Act.xml` | `ASM-FW-GCNMFW-WORK` | `SETMOCLAIMTREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `ASM-FW-GCNMFW-WORK-PNC / SETMOCLAIM_ACT` |
| `SetPPNPPH.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SETPPNPPH` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SetPayableTreatyNP_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SETPAYABLETREATYNP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=35 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / SETPAYABLETREATY_ACT` |
| `SetTPLNote_Act.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `SETTPLNOTE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=12 | — |
| `SetValidateInstallment_Act.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SETVALIDATEINSTALLMENT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SetValueClaimTNP_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `SETVALUECLAIMTNP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `SetViewDataMasterTreaty_Act.xml` | `DATA-PORTAL` | `SETVIEWDATAMASTERTREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SethistoryKlaimTreaty.xml` | `ASM-FW-GCNMFW-WORK` | `SETHISTORYKLAIMTREATY` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | — |
| `TotalEgnpi.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `TOTALEGNPI` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | `ASM-FW-GISFW-DATA-TREATYINLIMITLAYER / TOTALEGNPI` |
| `TreatyInNonSetTotal.xml` | `DATA-PORTAL` | `TREATYINNONSETTOTAL` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | — |
| `getStatusKonversi_Act.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GETSTATUSKONVERSI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |

### ConnectREST — 7 rule (`RULE-CONNECT-REST`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `InsertClaimOutstanding_NP.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `INSERTCLAIMOUTSTANDING_NP` | [terverifikasi] pyServiceName=InsertClaimOutstanding_NP; baseURL=URL; resolusi=JNDIName; URL_LITERAL_DITEMUKAN=2 | — |
| `KonversiKlaimNonLife.xml` | `ASM-FW-GCNMFW-WORK` | `KONVERSIKLAIMNONLIFE` | [terverifikasi] pyServiceName=KonversiKlaimNonLife; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | `ASM-FW-GCNMFW-WORK-PNC / INSERTREJECTCLAIM` |
| `SendAcceptationToKasir.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SENDACCEPTATIONTOKASIR` | [terverifikasi] pyServiceName=SendAcceptationToKasir; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | — |
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `SERVICEGOOGLE` | [terverifikasi] pyServiceName=ServiceGoogle; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / GOOGLESTORAGE_UPLOAD` |
| `getPayAttachment.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GETPAYATTACHMENT` | [terverifikasi] pyServiceName=getPayAttachment; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | — |
| `getPremiumPaidOnTreatyIn.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `GETPREMIUMPAIDONTREATYIN` | [terverifikasi] pyServiceName=getPremiumPaidOnTreatyIn; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / GETPREMIUMPAIDONTREATYIN` |
| `insertClaimFinalOrClosed_NP.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `INSERTCLAIMFINALORCLOSED_NP` | [terverifikasi] pyServiceName=insertClaimFinalOrClosed_NP; baseURL=URL; resolusi=JNDIName; URL_LITERAL_DITEMUKAN=2 | — |

### DataTransform — 10 rule (`RULE-OBJ-MODEL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AddCoB.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITLAYER` | `ADDCOB` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CNMRefreshCauseOfLoss_dt.xml` | `@BASECLASS` | `CNMREFRESHCAUSEOFLOSS_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / CNMREFRESHSURVEYORS_DT` |
| `CNMRefreshListDetailCauseOfLoss_dt.xml` | `@BASECLASS` | `CNMREFRESHLISTDETAILCAUSEOFLOSS_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CNMShowInsertCauseOfLoss_dt.xml` | `@BASECLASS` | `CNMSHOWINSERTCAUSEOFLOSS_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / CNMSHOWINSERTSUVERYORS_DT` |
| `CNMShowInsertDetailCauseOfLoss_dt.xml` | `@BASECLASS` | `CNMSHOWINSERTDETAILCAUSEOFLOSS_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / CNMINSERTDETAILCAUSEOFLOSS_DT` |
| `InsertChronology_DT.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `INSERTCHRONOLOGY_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetCoB.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITLAYER` | `SETCOB` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetIndexLayer_DT.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITLAYER` | `SETINDEXLAYER_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYINLIMITS / SETINDEXLAYER_DT` |
| `SetLabel_Dt.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `SETLABEL_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-PNC / SETLABEL_DT` |
| `SystemSetOneYear_DT.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `SYSTEMSETONEYEAR_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### DecisionTable — 1 rule (`RULE-DECLARE-DECISIONTABLE`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `GetMimeType.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETMIMETYPE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Flow — 1 rule (`RULE-OBJ-FLOW`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `Flow_TreatyIn.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `FLOW_TREATYIN` | [terverifikasi] pyStartActivity=Start1 | — |

### FlowAction — 16 rule (`RULE-OBJ-FLOWACTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AdjustmentDetailNP.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `ADJUSTMENTDETAILNP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / ADJUSTMENTDETAIL` |
| `CatastrofeList.xml` | `ASM-FW-GCNMFW-WORK` | `CATASTROFELIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CloseClaimMD.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `CLOSECLAIMMD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CloseClaimNP.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `CLOSECLAIMNP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP / CLOSECLAIMMD` |
| `CoBList.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITLAYER` | `COBLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GenerateDLACNP.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GENERATEDLACNP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GeneratePLACNP.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `GENERATEPLACNP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputAcceptation.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `INPUTACCEPTATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputDtlInterest.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `INPUTDTLINTEREST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputHistoricalSurveyReportUW.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `INPUTHISTORICALSURVEYREPORTUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-QUOTATION / INPUTHISTORICALSURVEYREPORT` |
| `Installments_ReadOnly.xml` | `ASM-FW-GISFW-DATA-TREATYININSTALLMENT` | `INSTALLMENTS_READONLY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYININSTALLMENT / INSTALLMENTS_REALISASI` |
| `OutstandingClaim.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `OUTSTANDINGCLAIM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ShowDetailXOL.xml` | `ASM-FW-GISFW-DATA-SPREADINGRISK` | `SHOWDETAILXOL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewClaimLayerDetail.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `VIEWCLAIMLAYERDETAIL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewDetailPayment_FA.xml` | `ASM-FW-GCNMFW-DATA-CURRENCY` | `VIEWDETAILPAYMENT_FA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewListPolicyCNP.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `VIEWLISTPOLICYCNP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Harness — 14 rule (`RULE-HTML-HARNESS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `CauseofLoss_Harness.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `CAUSEOFLOSS_HARNESS` | [terverifikasi] pyInclude=1 | `ASM-FW-GCNMFW-WORK-PNC / CAUSEOFLOSS_HARNESS` |
| `ChooseMasterTNonProp.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `CHOOSEMASTERTNONPROP` | [terverifikasi] pyInclude=1 | — |
| `DetailPaymentStsCNP.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `DETAILPAYMENTSTSCNP` | [terverifikasi] pyInclude=1 | — |
| `EditXOLAlokasi.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `EDITXOLALOKASI` | [terverifikasi] pyInclude=1 | — |
| `HistoricalSurveyReportUW.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `HISTORICALSURVEYREPORTUW` | [terverifikasi] pyInclude=1 | `ASM-FW-GISFW-DATA-POLICYTREATYIN / HISTORICALSURVEYREPORT` |
| `Hitung_Test.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `HITUNG_TEST` | [terverifikasi] pyInclude=1 | — |
| `KomiteCNP.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `KOMITECNP` | [terverifikasi] pyInclude=1 | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / COMMITTEETREATY` |
| `OutstandingClaim.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `OUTSTANDINGCLAIM` | [terverifikasi] pyInclude=9 | — |
| `TambahCauseofLoss.xml` | `DATA-PORTAL` | `TAMBAHCAUSEOFLOSS` | [terverifikasi] pyInclude=3 | — |
| `TambahMasterCauseOfLoss.xml` | `DATA-PORTAL` | `TAMBAHMASTERCAUSEOFLOSS` | [terverifikasi] pyInclude=3 | — |
| `ViewAttachmentNP.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `VIEWATTACHMENTNP` | [terverifikasi] pyInclude=1 | — |
| `ViewHistoryMasterID_NP.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `VIEWHISTORYMASTERID_NP` | [terverifikasi] pyInclude=1 | — |
| `ViewOldAllocation.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `VIEWOLDALLOCATION` | [terverifikasi] pyInclude=1 | — |
| `ViewPolisNonProp.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `VIEWPOLISNONPROP` | [terverifikasi] pyInclude=3 | — |

### RDBList — 50 rule (`RULE-CONNECT-SQL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseDtlTreatyNP.xml` | `ASM-FW-GISFW-INT-TREATYINDETAIL` | `BROWSEDTLTREATYNP` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TREATYINDETAIL,POOLDATA.TREATYINDETAILEDM | — |
| `BrowseDtlTreatyNPEDM.xml` | `ASM-FW-GISFW-INT-TREATYINDETAIL` | `BROWSEDTLTREATYNPEDM` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TREATYINDETAILEDM | — |
| `BrowseDtlTreatyOutNP.xml` | `ASM-FW-GISFW-INT-TREATYINDETAIL` | `BROWSEDTLTREATYOUTNP` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_TREATY_OUT_DETAIL | — |
| `BrowseRW_SQL.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `BROWSERW_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CITY,RW | — |
| `BrowseTreatyNP.xml` | `ASM-FW-GISFW-INT-TREATYINDETAIL` | `BROWSETREATYNP` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TREATYINDETAIL | `ASM-FW-GISFW-INT-OFFERJSON / ASM!GETHISTORYMASTERID` |
| `BrowseTreatyNP_EDM.xml` | `ASM-FW-GISFW-INT-TREATYINDETAIL` | `BROWSETREATYNP_EDM` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.TREATYINDETAILEDM | — |
| `CariHistoryClaim_SQL.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `CARIHISTORYCLAIM_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_KLAIM | — |
| `CekHistoryClaimNonProp_SQL.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `CEKHISTORYCLAIMNONPROP_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_KLAIM,CLAIMREJECTED | `ASM-FW-GCNMFW-INT-V_POLIS / GCNM!CARIHISTORYCLAIM_SQL` |
| `CekOSClaimNonProp_SQL.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `CEKOSCLAIMNONPROP_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=OS_AKSEPTASI_KLAIM | — |
| `CurrencyStandard.xml` | `ASM-FW-GISFW-INT-CURRENCYSTANDARD` | `CURRENCYSTANDARD` | [terverifikasi] sqlKind=QUERY; procs=POOLDATA.GETCURRENCYSTANDARD; sqlOps=SELECT | `ASM-FW-GISFW-INT-CURRENCYSTANDARD / ASM!UPDATEMASTERCURRENCYSTANDARD` |
| `GenerateImageID_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GENERATEIMAGEID_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GenerateNoPLATNP.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `GENERATENOPLATNP` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT | — |
| `GetAddressCeding.xml` | `ASM-FW-GISFW-INT-CLIENT` | `GETADDRESSCEDING` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_CLIENT,AGENT | — |
| `GetAddressTreatyIn.xml` | `ASM-FW-GISFW-INT-CLIENT` | `GETADDRESSTREATYIN` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_CLIENT,AGENT | `ASM-FW-GISFW-INT-CLIENT / ASM!GETADDRESSCEDING` |
| `GetAppName_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETAPPNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.T_FOLDER_IMAGE | — |
| `GetCaseIDNBTretyIn.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `GETCASEIDNBTRETYIN` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | `ASM-FW-GCNMFW-INT-V_POLIS / ASM!GETCASEIDNB` |
| `GetClientName.xml` | `ASM-FW-GISFW-INT-BANKACCOUNT` | `GETCLIENTNAME` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=BANKACCOUNT | `ASM-FW-GISFW-INT-BANKACCOUNT / GCNM!GETDATABYCLIENTNAME2` |
| `GetCurrency.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `GETCURRENCY` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CURRENCY | `ASM-FW-GISFW-INT-CURRENCY / ASM!UPDATEMASTERCURRENCY` |
| `GetDataBankAccount_sql.xml` | `ASM-FW-GISFW-INT-BANKACCOUNT` | `GETDATABANKACCOUNT_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=BANKACCOUNT | — |
| `GetDataBusiness_SQL.xml` | `ASM-FW-GISFW-INT-BUSINESS` | `GETDATABUSINESS_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=BUSINESS | — |
| `GetDataCNPOS.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `GETDATACNPOS` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=OS_AKSEPTASI_KLAIM | — |
| `GetDataMasterTOutNP.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `GETDATAMASTERTOUTNP` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_TREATY_OUT | — |
| `GetDataOS.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `GETDATAOS` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=OS_AKSEPTASI_KLAIM | `ASM-FW-GCNMFW-INT-V_POLIS / ASM!GETDATACNPOS` |
| `GetDataPolisNonProp_SQL.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `GETDATAPOLISNONPROP_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYINPRODUCTION | — |
| `GetDatabyClientName.xml` | `ASM-FW-GISFW-INT-BANKACCOUNT` | `GETDATABYCLIENTNAME` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=BANKACCOUNT | — |
| `GetDatabyClientName2.xml` | `ASM-FW-GISFW-INT-BANKACCOUNT` | `GETDATABYCLIENTNAME2` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=BANKACCOUNT | `ASM-FW-GISFW-INT-BANKACCOUNT / GCNM!GETDATABYCLIENTNAME` |
| `GetDatabyClientName3.xml` | `ASM-FW-GISFW-INT-BANKACCOUNT` | `GETDATABYCLIENTNAME3` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=BANKACCOUNT | `ASM-FW-GISFW-INT-BANKACCOUNT / GCNM!GETDATABYCLIENTNAME2` |
| `GetEmailCeding_SQL.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GETEMAILCEDING_SQL` | [terverifikasi] sqlKind=QUERY; procs=GL.F_GET_EMAIL; sqlOps=SELECT | — |
| `GetHistoryMasterID.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `GETHISTORYMASTERID` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.CLAIMXOL2,POOLDATA.OS_AKSEPTASI_KLAIM,POOLDATA.CLAIMXOL | `ASM-FW-GISFW-INT-OFFERJSON / ASM!GETTREATYINMASTER_SQL` |
| `GetKodeProdNonLife_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETKODEPRODNONLIFE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.KODE_PRODUKSI | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETSEQUENCENUMBER_SQL` |
| `GetLBUID_SQL.xml` | `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS_BUSINESS` | `GETLBUID_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=V_D_CAUSE_OF_LOSS_BUSINESS | `ASM-FW-GCNMFW-INT-V_M_CAUSE_OF_LOSS_BUSINESS / GCNM!GETLBUID_SQL` |
| `GetLeaderReport.xml` | `ASM-FW-GISFW-INT-AGENT` | `GETLEADERREPORT` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=AGENT | `ASM-FW-GISFW-INT-CLIENT / ASM!GETADDRESSCEDING` |
| `GetLimitTONPPLA.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `GETLIMITTONPPLA` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATY_OUT | `ASM-FW-GCNMFW-INT-V_POLIS / ASM!GETLIMITTREATYOUT` |
| `GetLimitsTreatyIn_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `GETLIMITSTREATYIN_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_TREATY_IN,M_TREATY_IN_EDM | `ASM-FW-GISFW-INT-OFFERJSON / ASM!GETCOPYNBFORCLAIM` |
| `GetLinkStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETLINKSTORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GETTOKENSTORAGE_SQL` |
| `GetMObyNopol_SQL.xml` | `ASM-FW-GISFW-INT-MARKETINGOFFICER` | `GETMOBYNOPOL_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=MARKETINGOFFICER,TREATYINPRODUCTION | — |
| `GetNopolis_SQL.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `GETNOPOLIS_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYINPRODUCTION | — |
| `GetPolicyData.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `GETPOLICYDATA` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=JSON_POLIS | — |
| `GetReinsuranceTypeBYName_SQL.xml` | `ASM-FW-GISFW-INT-REINSURANCETYPE` | `GETREINSURANCETYPEBYNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=REINSURANCETYPE | — |
| `GetSequenceNumber_SQL.xml` | `ASM-FW-GISFW-INT-POLICYJSON` | `GETSEQUENCENUMBER_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER | `ASM-FW-GISFW-INT-POLICYJSON / RNM!GETTANGGALCLOSING_SQL` |
| `GetTokenStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETTOKENSTORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.GET_TOKEN_STORAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |
| `GetTreatyName_SQL.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `GETTREATYNAME_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=POOLDATA.REINSURANCETYPE,PROPORTIONALARRG,TREATYBUSINESS,TREATYCONTRACT | — |
| `InsertClaimPNC.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `INSERTCLAIMPNC` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_KLAIM_PNC | — |
| `InsertLOGDirectKasir_SQL.xml` | `ASM-FW-GCNMFW-WORK` | `INSERTLOGDIRECTKASIR_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=POOLDATA.DIRECTTOKASIR_LOG | `ASM-FW-GISFW-WORK / RNM!INSERTLOGMOP_SQL` |
| `Insert_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `INSERT_T_STORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=INSERT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GENERATEIMAGEID_SQL` |
| `SaveDataToOsAkseptasiNP.xml` | `ASM-FW-GCNMFW-INT-V_POLIS` | `SAVEDATATOOSAKSEPTASINP` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP | `ASM-FW-GCNMFW-INT-V_POLIS / GCNM!SAVEDATATOOSAKSEPTASI` |
| `UpdateDCauseOfLoss.xml` | `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS` | `UPDATEDCAUSEOFLOSS` | [terverifikasi] sqlKind=PLSQL; procs=DBMS_LOB.CREATETEMPORARY,POOLDATA.PEGA_D_CAUSE_OF_LOSS | `ASM-FW-CNMFW-INT-V_D_CAUSE_OF_LOSS / CNM!UPDATEDCAUSEOFLOSS` |
| `UpdateMCauseOfLoss.xml` | `ASM-FW-GCNMFW-INT-V_M_CAUSE_OF_LOSS` | `UPDATEMCAUSEOFLOSS` | [terverifikasi] sqlKind=PLSQL; procs=DBMS_LOB.CREATETEMPORARY,POOLDATA.PEGA_M_CAUSE_OF_LOSS | — |
| `Update_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `UPDATE_T_STORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |
| `getStatusKonversi_SQL.xml` | `ASM-FW-GCNMFW-INT-OS_AKSEPTASI_KLAIM` | `GETSTATUSKONVERSI_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=REINSURANCE.TRLOSS_DETAIL_T | — |

### ReportDefinition — 22 rule (`RULE-OBJ-REPORT-DEFINITION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseAdjusterConsultant.xml` | `ASM-FW-GISFW-INT-ADJUSTERCONSULTANT` | `BROWSEADJUSTERCONSULTANT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseBankGroup.xml` | `ASM-FW-GISFW-INT-LST_BANK_GROUP` | `BROWSEBANKGROUP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseBusiness_RD.xml` | `ASM-FW-GISFW-INT-BUSINESS` | `BROWSEBUSINESS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseClientName_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSECLIENTNAME_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-AGENT / BROWSESEARCHSOBCEDING_RD` |
| `BrowseCouseOfLoss_Business.xml` | `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS_BUSINESS` | `BROWSECOUSEOFLOSS_BUSINESS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCurrencyTreatyIn_RD.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `BROWSECURRENCYTREATYIN_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseCurrency_RD.xml` | `ASM-FW-GISFW-INT-CURRENCY` | `BROWSECURRENCY_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseProvince_RD.xml` | `ASM-FW-GISFW-INT-PROVINCE` | `BROWSEPROVINCE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseRW_RD.xml` | `ASM-FW-GISFW-INT-RW` | `BROWSERW_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseReinsuranceType_RD.xml` | `ASM-FW-GISFW-INT-REINSURANCETYPE` | `BROWSEREINSURANCETYPE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyBusinessWOType_RD.xml` | `ASM-FW-GISFW-INT-TREATYBUSINESS` | `BROWSETREATYBUSINESSWOTYPE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-TREATYBUSINESS / BROWSETREATYBUSINESS_RD` |
| `BrowseTreatyGroup_RD.xml` | `ASM-FW-GISFW-INT-TREATYGROUP` | `BROWSETREATYGROUP_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseVDCauseOfLoss_RD.xml` | `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS` | `BROWSEVDCAUSEOFLOSS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseVMCauseOfLoss_RD.xml` | `ASM-FW-GCNMFW-INT-V_M_CAUSE_OF_LOSS` | `BROWSEVMCAUSEOFLOSS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseVMstUserTeknis_RD.xml` | `ASM-FW-GCNMFW-INT-V_MST_USER_TEKNIS` | `BROWSEVMSTUSERTEKNIS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `FilterEmailKomiteWithLimit.xml` | `ASM-FW-GCNMFW-INT-EMAILKOMITE` | `FILTEREMAILKOMITEWITHLIMIT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GCNMGetAllAttachments.xml` | `LINK-ATTACHMENT` | `GCNMGETALLATTACHMENTS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GCNMGetInvoiceAttachments.xml` | `LINK-ATTACHMENT` | `GCNMGETINVOICEATTACHMENTS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `LINK-ATTACHMENT / GCNMGETALLATTACHMENTS` |
| `GetCatastrope_RD.xml` | `ASM-FW-GCNMFW-INT-CATASTROPHE` | `GETCATASTROPE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `RejectedClaim_RD.xml` | `ASM-FW-GCNMFW-INT-CLAIMREJECTED` | `REJECTEDCLAIM_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SelectVDCauseOfLoss_RD.xml` | `ASM-FW-GCNMFW-INT-V_D_CAUSE_OF_LOSS` | `SELECTVDCAUSEOFLOSS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-CNMFW-INT-V_D_CAUSE_OF_LOSS / SELECTVDCAUSEOFLOSS_RD` |
| `SelectVMCauseOfLoss_RD.xml` | `ASM-FW-GCNMFW-INT-V_M_CAUSE_OF_LOSS` | `SELECTVMCAUSEOFLOSS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Section — 36 rule (`RULE-HTML-SECTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `AdjustmentDetailNP.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `ADJUSTMENTDETAILNP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / ADJUSTMENTDETAIL` |
| `AdjustmentDetailNP_Section.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `ADJUSTMENTDETAILNP_SECTION` | [terverifikasi] pyInclude=4 | — |
| `BrowseCauseOfLoss.xml` | `@BASECLASS` | `BROWSECAUSEOFLOSS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseDetailCauseOfLoss.xml` | `@BASECLASS` | `BROWSEDETAILCAUSEOFLOSS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CatastrofeList_Sec.xml` | `ASM-FW-GCNMFW-WORK` | `CATASTROFELIST_SEC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Catastrope_Sec.xml` | `ASM-FW-GCNMFW-WORK` | `CATASTROPE_SEC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CauseofLoss_Section.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATY` | `CAUSEOFLOSS_SECTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-PNC / CAUSEOFLOSS_SECTION` |
| `ChooseMasterTNonProp.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `CHOOSEMASTERTNONPROP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CloseClaimMD.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `CLOSECLAIMMD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `CloseClaimNP.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `CLOSECLAIMNP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP / CLOSECLAIMMD` |
| `CoBList.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITLAYER` | `COBLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DetailPaymentCNP.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `DETAILPAYMENTCNP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / PAYMENTPREMILIST_SECTION` |
| `DetailPolisCNP.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `DETAILPOLISCNP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-POLICYTREATYIN / DETAILPOLICYTREATYINNONPROPORTIONAL` |
| `EditXOLAlokasi.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `EDITXOLALOKASI` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GenerateDLACNP.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `GENERATEDLACNP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridCauseOfLoss.xml` | `DATA-PORTAL` | `GRIDCAUSEOFLOSS` | [terverifikasi] pyInclude=2 | — |
| `HistoricalSurveyReportDtlUW.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `HISTORICALSURVEYREPORTDTLUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-POLICYTREATYIN / HISTORICALSURVEYREPORTDTL` |
| `Hitung_Test.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `HITUNG_TEST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputAcceptation.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `INPUTACCEPTATION` | [terverifikasi] pyInclude=1 | — |
| `InputDtlInterest.xml` | `ASM-FW-GCNMFW-DATA-OBJECTITEM` | `INPUTDTLINTEREST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputHistoricalSurveyReportDtlUW.xml` | `ASM-FW-GISFW-DATA-QUOTATION` | `INPUTHISTORICALSURVEYREPORTDTLUW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-QUOTATION / INPUTHISTORICALSURVEYREPORTDTL` |
| `Installments_ReadOnly.xml` | `ASM-FW-GISFW-DATA-TREATYININSTALLMENT` | `INSTALLMENTS_READONLY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-DATA-TREATYININSTALLMENT / INSTALLMENTS_REALISASI` |
| `KomiteCLMNP.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `KOMITECLMNP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / COMITEECLAIMTREATY` |
| `ListDetailCauseOfLoss.xml` | `DATA-PORTAL` | `LISTDETAILCAUSEOFLOSS` | [terverifikasi] pyInclude=2 | — |
| `OutstandingClaim(1).xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `OUTSTANDINGCLAIM` | [terverifikasi] pyInclude=1 | — |
| `PreviewPLA.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `PREVIEWPLA` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ReinstatementPremiumDetails.xml` | `ASM-FW-GISFW-DATA-SPREADINGRISK` | `REINSTATEMENTPREMIUMDETAILS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `Subjectivity.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SUBJECTIVITY` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-DATA-ADJUSTMENT / DTLDATACOMMITTE` |
| `ViewAttachmentNP.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `VIEWATTACHMENTNP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / VIEWATTACHMENT` |
| `ViewClaimLayerDetail.xml` | `ASM-FW-GISFW-DATA-TREATYINLIMITS` | `VIEWCLAIMLAYERDETAIL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewDetailDeptHeadTreatyIn_UW.xml` | `ASM-FW-GISFW-DATA-POLICYTREATYIN` | `VIEWDETAILDEPTHEADTREATYIN_UW` | [terverifikasi] pyInclude=2 | `ASM-FW-GISFW-DATA-POLICYTREATYIN / DETAILDEPTHEADTREATYIN_UW` |
| `ViewDetailPayment_SC.xml` | `ASM-FW-GCNMFW-DATA-CURRENCY` | `VIEWDETAILPAYMENT_SC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewDetailPolisTNonProp.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `VIEWDETAILPOLISTNONPROP` | [terverifikasi] pyInclude=2 | — |
| `ViewHistoryMasterID_NP.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `VIEWHISTORYMASTERID_NP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `DATA-PORTAL / INBOXCLAIMNONPROP_HARNESS` |
| `ViewListPolicyCNP.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `VIEWLISTPOLICYCNP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewOldAllocation.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `VIEWOLDALLOCATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### SystemSettings — 1 rule (`RULE-ADMIN-SYSTEM-SETTINGS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `LinkService.xml` | `LINKSERVICE` | `LINKSERVICE` | [terverifikasi] pyPurpose=LinkService | — |

### When — 7 rule (`RULE-OBJ-WHEN`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `IsBackStage.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `ISBACKSTAGE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GCNMFW-WORK-CLAIMTREATY / ISBACKSTAGE` |
| `IsCLM.xml` | `@BASECLASS` | `ISCLM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsCLMNP.xml` | `@BASECLASS` | `ISCLMNP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / ISCLMP` |
| `IsCLMP.xml` | `@BASECLASS` | `ISCLMP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / ISCLM` |
| `IsClaim.xml` | `@BASECLASS` | `ISCLAIM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `IsPEGAPROD.xml` | `@BASECLASS` | `ISPEGAPROD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `WORK- / ISPEGAPROD` |
| `IsPEGASyariah.xml` | `@BASECLASS` | `ISPEGASYARIAH` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

## 4. Tabrakan nama di dalam modul ini

### 4.1 Nama file dipakai di lebih dari satu tipe rule

| Nama file | Tipe | Class |
| --- | --- | --- |
| `ViewListPolicyCNP.xml` | FlowAction, Section | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP, ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` |
| `ViewClaimLayerDetail.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINLIMITS, ASM-FW-GISFW-DATA-TREATYINLIMITS` |
| `AdjustmentDetailNP.xml` | FlowAction, Section | `ASM-FW-GCNMFW-DATA-ADJUSTMENT, ASM-FW-GCNMFW-DATA-ADJUSTMENT` |
| `CloseClaimNP.xml` | FlowAction, Section | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP, ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` |
| `Installments_ReadOnly.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYININSTALLMENT, ASM-FW-GISFW-DATA-TREATYININSTALLMENT` |
| `ChooseMasterTNonProp.xml` | Harness, Section | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP, ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` |
| `OutstandingClaim.xml` | FlowAction, Harness | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP, ASM-FW-GCNMFW-WORK-CLAIMTREATY` |
| `EditXOLAlokasi.xml` | Activity, Harness, Section | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP, ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP, ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` |
| `GenerateDLACNP.xml` | FlowAction, Section | `ASM-FW-GCNMFW-DATA-ADJUSTMENT, ASM-FW-GCNMFW-DATA-ADJUSTMENT` |
| `ViewOldAllocation.xml` | Harness, Section | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP, ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` |
| `CoBList.xml` | FlowAction, Section | `ASM-FW-GISFW-DATA-TREATYINLIMITLAYER, ASM-FW-GISFW-DATA-TREATYINLIMITLAYER` |
| `InputDtlInterest.xml` | FlowAction, Section | `ASM-FW-GCNMFW-DATA-OBJECTITEM, ASM-FW-GCNMFW-DATA-OBJECTITEM` |
| `InputAcceptation.xml` | FlowAction, Section | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP, ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` |
| `ViewHistoryMasterID_NP.xml` | Harness, Section | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP, ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` |
| `CloseClaimMD.xml` | FlowAction, Section | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP, ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` |
| `ViewAttachmentNP.xml` | Harness, Section | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP, ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` |
| `Hitung_Test.xml` | Harness, Section | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP, ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` |

Total: **17** nama file.

### 4.2 Rule `When` bernama sama di lebih dari satu class

_Tidak ada._

## 5. Keluarga klon (satu asal salinan dipakai >1 rule)

Dibaca dari `<pzOriginalInstanceKey>` yang sama pada beberapa file. Ini **sinyal kandidat
konsolidasi** saat migrasi, bukan bukti bahwa perilakunya identik — kesamaan perilaku baru
dapat dinyatakan setelah isinya dibandingkan di STEP D2.

| Asal salinan (class / nama) | Jumlah | File |
| --- | ---: | --- |
| `ASM-FW-GISFW-INT-CLIENT / ASM!GETADDRESSCEDING` | 3 | `RDBList/GetAddressCeding.xml`, `RDBList/GetAddressTreatyIn.xml`, `RDBList/GetLeaderReport.xml` |
| `ASM-FW-GCNMFW-INT-V_POLIS / GCNM!CARIHISTORYCLAIM_SQL` | 2 | `RDBList/CariHistoryClaim_SQL.xml`, `RDBList/CekHistoryClaimNonProp_SQL.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP / ADJCLAIMAMOUNT_ACT` | 2 | `Activity/AdjClaimAmount_Act.xml`, `Activity/AdjClaimCNP_Act.xml` |
| `LINK-ATTACHMENT / GCNMGETALLATTACHMENTS` | 2 | `ReportDefinition/GCNMGetAllAttachments.xml`, `ReportDefinition/GCNMGetInvoiceAttachments.xml` |
| `ASM-FW-GISFW-INT-BANKACCOUNT / GCNM!GETDATABYCLIENTNAME2` | 2 | `RDBList/GetClientName.xml`, `RDBList/GetDatabyClientName3.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP / CLOSECLAIMMD` | 4 | `FlowAction/CloseClaimMD.xml`, `FlowAction/CloseClaimNP.xml`, `Section/CloseClaimMD.xml`, `Section/CloseClaimNP.xml` |
| `ASM-FW-GCNMFW-INT-V_POLIS / ASM!GETDATACNPOS` | 2 | `RDBList/GetDataCNPOS.xml`, `RDBList/GetDataOS.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATY / COUNTTOTALINSTEREST_ACT` | 2 | `Activity/CopyOldataCurr_act.xml`, `Activity/CountTotalInsterest_Act.xml` |
| `@BASECLASS / ISCLM` | 2 | `When/IsCLM.xml`, `When/IsCLMP.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / INSERTGOOGLESTORAGE_ACT` | 2 | `Activity/GetUrlGoogleStorage_Act.xml`, `Activity/InsertGoogleStorage_Act.xml` |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN / COUNTNETPREMI_ACT` | 2 | `Activity/CountNetPremi_act.xml`, `Activity/SetDueTo_act.xml` |
| `ASM-FW-GCNMFW-DATA-ADJUSTMENT / GENERATEDLACNP` | 2 | `FlowAction/GenerateDLACNP.xml`, `Section/GenerateDLACNP.xml` |
| `ASM-FW-GISFW-DATA-TREATYININSTALLMENT / INSTALLMENTS_REALISASI` | 2 | `FlowAction/Installments_ReadOnly.xml`, `Section/Installments_ReadOnly.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GENERATEIMAGEID_SQL` | 2 | `RDBList/GenerateImageID_SQL.xml`, `RDBList/Insert_T_Storage_SQL.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP / VIEWOLDALLOCATION` | 2 | `Harness/ViewOldAllocation.xml`, `Section/ViewOldAllocation.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP / INPUTACCEPTATION` | 2 | `FlowAction/InputAcceptation.xml`, `Section/InputAcceptation.xml` |
| `@BASECLASS / HTMLTOPDF` | 2 | `Activity/HTMLRISlipToPDF.xml`, `Activity/HTMLToPDF.xml` |
| `ASM-FW-GISFW-DATA-TREATYINLIMITLAYER / COBLIST` | 2 | `FlowAction/CoBList.xml`, `Section/CoBList.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP / INPUTOUTSTANDINGCLMTNP_PREACT` | 2 | `Activity/GetHistoryMasterID_NP.xml`, `Activity/InputOutStandingClmTNP_PreAct.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP / GENERATEPLACNP_ACT` | 2 | `Activity/GeneratePlaCNP2_Act.xml`, `Activity/GeneratePlaCNP_Act.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` | 2 | `RDBList/GetTokenStorage_SQL.xml`, `RDBList/Update_T_Storage_SQL.xml` |
| `ASM-FW-GCNMFW-DATA-OBJECTITEM / INPUTDTLINTEREST` | 2 | `FlowAction/InputDtlInterest.xml`, `Section/InputDtlInterest.xml` |
| `ASM-FW-GCNMFW-DATA-ADJUSTMENT / SAVECNPLAYERLIST_ACT` | 2 | `Activity/ProtectNilaiClaim.xml`, `Activity/SaveAdjustmentToOSAksep_Act_Tes.xml` |
| `ASM-FW-GCNMFW-DATA-ADJUSTMENT / ADJUSTMENTDETAIL` | 2 | `FlowAction/AdjustmentDetailNP.xml`, `Section/AdjustmentDetailNP.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP / EDITXOLALOKASI` | 3 | `Activity/EditXOLAlokasi.xml`, `Harness/EditXOLAlokasi.xml`, `Section/EditXOLAlokasi.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP / CHOOSEMASTERTNONPROP` | 2 | `Harness/ChooseMasterTNonProp.xml`, `Section/ChooseMasterTNonProp.xml` |
| `ASM-FW-GISFW-INT-BANKACCOUNT / GCNM!GETDATABYCLIENTNAME` | 2 | `RDBList/GetDatabyClientName.xml`, `RDBList/GetDatabyClientName2.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP / VIEWLISTPOLICYCNP` | 2 | `FlowAction/ViewListPolicyCNP.xml`, `Section/ViewListPolicyCNP.xml` |
| `ASM-FW-GISFW-DATA-TREATYINLIMITS / VIEWCLAIMLAYERDETAIL` | 2 | `FlowAction/ViewClaimLayerDetail.xml`, `Section/ViewClaimLayerDetail.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP / HITUNG_TEST` | 2 | `Harness/Hitung_Test.xml`, `Section/Hitung_Test.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP / OUTSTANDINGCLAIM` | 2 | `FlowAction/OutstandingClaim.xml`, `Section/OutstandingClaim(1).xml` |

## 6. Objek database yang disentuh

Diambil dari `<pyBrowseSQL>` pada rule `RDBList`. **Nama objek apa adanya dari SQL** —
kualifikasi skema tidak konsisten di korpus. **Tipe kolom, constraint, dan relasi tidak
diketahui** (OQ-001).

| Tabel/view (apa adanya di SQL) | Dirujuk oleh |
| --- | ---: |
| `BANKACCOUNT` | 5 rule |
| `AGENT` | 3 rule |
| `OS_AKSEPTASI_KLAIM` | 3 rule |
| `POOLDATA.TREATYINDETAILEDM` | 3 rule |
| `TREATYINPRODUCTION` | 3 rule |
| `T_STORAGE_IMAGE` | 3 rule |
| `JSON_KLAIM` | 2 rule |
| `JSON_POLIS` | 2 rule |
| `M_CLIENT` | 2 rule |
| `POOLDATA.TREATYINDETAIL` | 2 rule |
| `BUSINESS` | 1 rule |
| `CITY` | 1 rule |
| `CLAIMREJECTED` | 1 rule |
| `CURRENCY` | 1 rule |
| `MARKETINGOFFICER` | 1 rule |
| `M_TREATY_IN` | 1 rule |
| `M_TREATY_IN_EDM` | 1 rule |
| `M_TREATY_OUT` | 1 rule |
| `M_TREATY_OUT_DETAIL` | 1 rule |
| `POOLDATA.CLAIMXOL` | 1 rule |
| `POOLDATA.CLAIMXOL2` | 1 rule |
| `POOLDATA.DIRECTTOKASIR_LOG` | 1 rule |
| `POOLDATA.KODE_PRODUKSI` | 1 rule |
| `POOLDATA.OS_AKSEPTASI_KLAIM` | 1 rule |
| `POOLDATA.REINSURANCETYPE` | 1 rule |
| `POOLDATA.T_FOLDER_IMAGE` | 1 rule |
| `PROPORTIONALARRG` | 1 rule |
| `REINSURANCE.TRLOSS_DETAIL_T` | 1 rule |
| `REINSURANCETYPE` | 1 rule |
| `RW` | 1 rule |
| `TREATYBUSINESS` | 1 rule |
| `TREATYCONTRACT` | 1 rule |
| `TREATY_OUT` | 1 rule |
| `V_D_CAUSE_OF_LOSS_BUSINESS` | 1 rule |

### Stored procedure / function yang dipanggil

**Body prosedur ini TIDAK ADA di korpus** — logika di dalamnya tidak dapat direkonstruksi
dari XML. Lihat OQ-002.

| Procedure | Dipanggil oleh |
| --- | --- |
| `DBMS_LOB.CREATETEMPORARY` | `UpdateDCauseOfLoss.xml`, `UpdateMCauseOfLoss.xml` |
| `POOLDATA.GETCURRENCYSTANDARD` | `CurrencyStandard.xml` |
| `POOLDATA.PEGA_JSON_KLAIM_PNC` | `InsertClaimPNC.xml` |
| `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` | `GetSequenceNumber_SQL.xml` |
| `POOLDATA.PEGA_M_CAUSE_OF_LOSS` | `UpdateMCauseOfLoss.xml` |
| `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP` | `SaveDataToOsAkseptasiNP.xml` |
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
| `InsertClaimOutstanding_NP.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `InsertClaimOutstanding_NP` | `URL` | — | **YA** |
| `KonversiKlaimNonLife.xml` | `ASM-FW-GCNMFW-WORK` | `KonversiKlaimNonLife` | `SETTING` | `LinkService!LinkService` | tidak |
| `SendAcceptationToKasir.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `SendAcceptationToKasir` | `SETTING` | `LinkService!LinkService` | tidak |
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `ServiceGoogle` | `SETTING` | `LinkService!LinkService` | tidak |
| `getPayAttachment.xml` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT` | `getPayAttachment` | `SETTING` | `LinkService!LinkService` | tidak |
| `getPremiumPaidOnTreatyIn.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `getPremiumPaidOnTreatyIn` | `SETTING` | `LinkService!LinkService` | tidak |
| `insertClaimFinalOrClosed_NP.xml` | `ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP` | `insertClaimFinalOrClosed_NP` | `URL` | — | **YA** |

