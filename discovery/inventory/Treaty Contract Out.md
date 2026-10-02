# Inventaris Rule — Treaty Contract Out

STEP D1, batch 4 (domain life, master & outward). Dibangkitkan 2026-09-12 dari korpus READ-ONLY
`D:\XML\RNM_BRD\Treaty Contract Out\`.

**Batas step ini:** hanya struktur + metadata. Tidak ada kesimpulan logika bisnis di sini.
Penelusuran perilaku adalah STEP D2 (`../flows/`).

**Identitas rule** diambil dari `<pxInsName>` (lihat `../README.md` §6.4).  
Kolom **Asal salinan** berisi `<pzOriginalInstanceKey>` dan hanya diisi bila **berbeda**
dari identitas — artinya rule ini hasil *Save As* dari rule lain. Tanda `—` berarti sama.
Class `—` berarti `<pxInsName>` tidak berisi bagian class (berlaku pada `DataPage`, yang
beridentitas global, bukan per class).

Perintah audit (dijalankan dari `D:\XML\RNM_BRD\`):

```
find "Treaty Contract Out" -type f -name "*.xml" | wc -l
find "Treaty Contract Out" -type f -name "*.xml" | awk -F/ '{print $2}' | sort | uniq -c
grep -o "<pxInsName>[^<]*" "Treaty Contract Out/<tipe>/<file>.xml" | head -1
grep -o "<pzOriginalInstanceKey>[^<]*" "Treaty Contract Out/<tipe>/<file>.xml" | head -1
```

## 1. Ringkasan terukur

- Total file XML: **303** — cocok dengan `../README.md` §5.1
- **Identitas rule unik: 303** dari 303 file → tidak ada identitas ganda di dalam modul ini
- Class Pega unik: **25**
- RDBList menurut jenis SQL: PLSQL=11, QUERY=26 (`PLSQL` = blok anonim `BEGIN…END;`, umumnya memanggil stored procedure)
- Prefix pada kunci RDBList: `ASM`=31, `RNM`=6 (arti prefix **belum terverifikasi** — OQ-008)

| Tipe rule | Jumlah |
| --- | ---: |
| Activity | 168 |
| ConnectREST | 1 |
| DataPage | 2 |
| DataTransform | 10 |
| FlowAction | 2 |
| Harness | 3 |
| RDBList | 37 |
| ReportDefinition | 30 |
| Section | 49 |
| SystemSettings | 1 |
| **TOTAL** | **303** |

## 2. Sebaran class Pega

Class yang **tidak** berawalan `ASM-` adalah class bawaan Pega atau class portal — ruang
lingkup migrasinya belum ditentukan (**OQ-009**).

| Class | Jumlah rule | Aplikasi? |
| --- | ---: | --- |
| `@BASECLASS` | 202 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-INT-PROPORTIONALARRG` | 40 | aplikasi |
| `DATA-PORTAL` | 8 | **bukan aplikasi** → OQ-009 |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | 7 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYREINSURER` | 5 | aplikasi |
| `ASM-FW-GISFW-INT-TREATY_IN` | 5 | aplikasi |
| `ASM-FW-GISFW-INT` | 4 | aplikasi |
| `ASM-FW-GISFW-INT-MTREATYSECURITY` | 4 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYBUSINESS` | 4 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYYEAR` | 4 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYCONTRACT` | 3 | aplikasi |
| _(tanpa class)_ | 2 | n/a (DataPage) |
| `ASM-FW-GISFW-INT-REINSURANCETYPE` | 2 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYDESC` | 2 | aplikasi |
| `ASM-FW-GISFW-DATA-ENUMERATION` | 1 | aplikasi |
| `ASM-FW-GISFW-DATA-OFFERFACIN` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-AGENT` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-BUSINESS` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-CLAUSE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-M_LINK_SERVICE` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-OCCUPATION` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-OFFERJSON` | 1 | aplikasi |
| `ASM-FW-GISFW-INT-TREATYGROUP` | 1 | aplikasi |
| `CODE-PEGA-LIST` | 1 | **bukan aplikasi** → OQ-009 |
| `LINKSERVICE` | 1 | **bukan aplikasi** → OQ-009 |

Rule di class bukan-aplikasi: **212** dari 303.

## 3. Daftar rule per tipe

### Activity — 168 rule (`RULE-OBJ-ACTIVITY`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseCopyData.xml` | `@BASECLASS` | `BROWSECOPYDATA` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `BrowseDeleteRowTreatyInContract.xml` | `@BASECLASS` | `BROWSEDELETEROWTREATYINCONTRACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `BrowseDescriptionLimit.xml` | `ASM-FW-GISFW-INT-TREATYDESC` | `BROWSEDESCRIPTIONLIMIT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `BrowsePortfolioList_Act.xml` | `@BASECLASS` | `BROWSEPORTFOLIOLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `BrowseReinsTypeYear.xml` | `ASM-FW-GISFW-INT` | `BROWSEREINSTYPEYEAR` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `BrowseTreatyArrCLaimCoorpParentList.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRCLAIMCOORPPARENTLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `ASM-FW-GISFW-INT-PROPORTIONALARRG / BROWSETREATYARRLIMITPARENTLIST` |
| `BrowseTreatyArrCashLossParentList.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRCASHLOSSPARENTLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `ASM-FW-GISFW-INT-PROPORTIONALARRG / BROWSETREATYARRLIMITPARENTLIST` |
| `BrowseTreatyArrEpiParentList.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARREPIPARENTLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `ASM-FW-GISFW-INT-PROPORTIONALARRG / BROWSETREATYARRLIMITPARENTLIST` |
| `BrowseTreatyArrFacInParentList.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRFACINPARENTLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `ASM-FW-GISFW-INT-PROPORTIONALARRG / BROWSETREATYARRLIMITPARENTLIST` |
| `BrowseTreatyArrLimitParentList.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRLIMITPARENTLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `BrowseTreatyArrPLAParentList.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRPLAPARENTLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `ASM-FW-GISFW-INT-PROPORTIONALARRG / BROWSETREATYARRLIMITPARENTLIST` |
| `BrowseTreatyArrangementExGratiaChild_Act.xml` | `@BASECLASS` | `BROWSETREATYARRANGEMENTEXGRATIACHILD_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / BROWSETREATYARRANGEMENTEXGRATIA_ACT` |
| `BrowseTreatyArrangementExGratia_Act.xml` | `ASM-FW-GISFW-INT` | `BROWSETREATYARRANGEMENTEXGRATIA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `ASM-FW-GISFW-INT / BROWSETREATYARRANGEMENTFACIN_ACT` |
| `BrowseTreatyBusinessList_Act.xml` | `ASM-FW-GISFW-INT` | `BROWSETREATYBUSINESSLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `BrowseTreatyReinsurerList_Act.xml` | `ASM-FW-GISFW-INT` | `BROWSETREATYREINSURERLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `CalculateTSIExcludeTreaty.xml` | `@BASECLASS` | `CALCULATETSIEXCLUDETREATY` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `CancelActivityBordereAux.xml` | `@BASECLASS` | `CANCELACTIVITYBORDEREAUX` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / CANCELACTIVITYCLAIMCOORPERATION` |
| `CancelActivityCashLossLimit.xml` | `@BASECLASS` | `CANCELACTIVITYCASHLOSSLIMIT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / CANCELACTIVITYEXGRATIA` |
| `CancelActivityClaimCoorperation.xml` | `@BASECLASS` | `CANCELACTIVITYCLAIMCOORPERATION` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / CANCELACTIVITYPROFITCOMMISION` |
| `CancelActivityEpi.xml` | `@BASECLASS` | `CANCELACTIVITYEPI` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / CANCELACTIVITYBORDEREAUX` |
| `CancelActivityExGratia.xml` | `@BASECLASS` | `CANCELACTIVITYEXGRATIA` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / CANCELACTIVITYEPI` |
| `CancelActivityExGratiaLimitChild.xml` | `@BASECLASS` | `CANCELACTIVITYEXGRATIALIMITCHILD` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / CANCELACTIVITYCASHLOSSLIMITCHILD` |
| `CancelActivityFacIn.xml` | `@BASECLASS` | `CANCELACTIVITYFACIN` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / CANCELACTIVITYPLA` |
| `CancelActivityFacInList.xml` | `@BASECLASS` | `CANCELACTIVITYFACINLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / CANCELACTIVITYPLA` |
| `CancelActivityPLA.xml` | `@BASECLASS` | `CANCELACTIVITYPLA` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / CANCELACTIVITYCASHLOSSLIMIT` |
| `CancelActivityProfitCommision.xml` | `@BASECLASS` | `CANCELACTIVITYPROFITCOMMISION` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / CANCELACTIVITYRICOMM` |
| `CancelActivityRicomm.xml` | `@BASECLASS` | `CANCELACTIVITYRICOMM` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / CANCELACTIVITYTREATYLIMIT` |
| `CancelActivityTreatyContract.xml` | `@BASECLASS` | `CANCELACTIVITYTREATYCONTRACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `CancelActivityTreatyLimit.xml` | `@BASECLASS` | `CANCELACTIVITYTREATYLIMIT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / CANCELACTIVITYTREATYLIMITCHILD` |
| `CancelActivityTreatyLimitChild.xml` | `@BASECLASS` | `CANCELACTIVITYTREATYLIMITCHILD` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / CANCELACTIVITY` |
| `CancelActivitypPortfolio.xml` | `@BASECLASS` | `CANCELACTIVITYPPORTFOLIO` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / CANCELACTIVITYPLA` |
| `CancelActivitypTerrLimit.xml` | `@BASECLASS` | `CANCELACTIVITYPTERRLIMIT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / CANCELACTIVITYPPORTFOLIO` |
| `CheckYear.xml` | `@BASECLASS` | `CHECKYEAR` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `D_TreatyContract_Act.xml` | `CODE-PEGA-LIST` | `D_TREATYCONTRACT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `DeleteAttachmentTreaty.xml` | `@BASECLASS` | `DELETEATTACHMENTTREATY` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `DeleteGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `DELETEGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / GETURLGOOGLESTORAGE_ACT` |
| `DeleteRowBusiness.xml` | `@BASECLASS` | `DELETEROWBUSINESS` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `DeleteSecurityReinsurer.xml` | `@BASECLASS` | `DELETESECURITYREINSURER` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `DeleteTreatyReins_Act.xml` | `ASM-FW-GISFW-INT-TREATYREINSURER` | `DELETETREATYREINS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `Delete_act.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `DELETE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `DownloadAll_Act.xml` | `DATA-PORTAL` | `DOWNLOADALL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=11 | — |
| `Download_Act.xml` | `DATA-PORTAL` | `DOWNLOAD_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetLinkService.xml` | `ASM-FW-GISFW-INT-M_LINK_SERVICE` | `GETLINKSERVICE` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetMaxCoinsPanel.xml` | `@BASECLASS` | `GETMAXCOINSPANEL` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / GETMINIMUMLOL` |
| `GetMinimumLOL.xml` | `@BASECLASS` | `GETMINIMUMLOL` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / GETPERIODE` |
| `GetMinimumLOLMB.xml` | `@BASECLASS` | `GETMINIMUMLOLMB` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / GETMINIMUMLOL` |
| `GetObject.xml` | `@BASECLASS` | `GETOBJECT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / GETPERIODE` |
| `GetPeriode.xml` | `@BASECLASS` | `GETPERIODE` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `GetUrlGoogleStorage_Act.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETURLGOOGLESTORAGE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / INSERTGOOGLESTORAGE_ACT` |
| `HitungRpUsd.xml` | `@BASECLASS` | `HITUNGRPUSD` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `HitungRpUsd_depan.xml` | `@BASECLASS` | `HITUNGRPUSD_DEPAN` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `KirimDataRPUsd.xml` | `@BASECLASS` | `KIRIMDATARPUSD` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `KirimIdMethod.xml` | `@BASECLASS` | `KIRIMIDMETHOD` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `KirimTahunGroupID.xml` | `@BASECLASS` | `KIRIMTAHUNGROUPID` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `LoadAttachment.xml` | `DATA-PORTAL` | `LOADATTACHMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `WORK- / LOADATTACHMENT` |
| `LoadAttachmentTreatyOut.xml` | `@BASECLASS` | `LOADATTACHMENTTREATYOUT` | [terverifikasi] pyActivityType=ACTIVITY; steps=0 | `DATA-PORTAL / LOADATTACHMENTTREATYOUT` |
| `NewInputTreatyContract_Act.xml` | `@BASECLASS` | `NEWINPUTTREATYCONTRACT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `NewInputTreatyYearMultiple_Act.xml` | `@BASECLASS` | `NEWINPUTTREATYYEARMULTIPLE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / NEWINPUTTREATYYEAR_ACT` |
| `NewInputTreatyYear_Act.xml` | `@BASECLASS` | `NEWINPUTTREATYYEAR_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `@BASECLASS / NEWINPUTTREATYCONTRACT_ACT` |
| `NewTreatyArrBordereAux.xml` | `@BASECLASS` | `NEWTREATYARRBORDEREAUX` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `NewTreatyArrCashLossLimit.xml` | `@BASECLASS` | `NEWTREATYARRCASHLOSSLIMIT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / NEWTREATYARRPLA` |
| `NewTreatyArrCashLossLimitList.xml` | `@BASECLASS` | `NEWTREATYARRCASHLOSSLIMITLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / NEWTREATYARRCASHLOSSLIMIT` |
| `NewTreatyArrClaimCoorp.xml` | `@BASECLASS` | `NEWTREATYARRCLAIMCOORP` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `NewTreatyArrClaimCoorpList.xml` | `@BASECLASS` | `NEWTREATYARRCLAIMCOORPLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / NEWTREATYARRTREATYLIMITLIST` |
| `NewTreatyArrCoins.xml` | `@BASECLASS` | `NEWTREATYARRCOINS` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / NEWTREATYARREXCLUTIONTREATY` |
| `NewTreatyArrEpi.xml` | `@BASECLASS` | `NEWTREATYARREPI` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / NEWTREATYARRCASHLOSSLIMIT` |
| `NewTreatyArrEpiList.xml` | `@BASECLASS` | `NEWTREATYARREPILIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / NEWTREATYARRTREATYLIMITLIST` |
| `NewTreatyArrExGratia.xml` | `@BASECLASS` | `NEWTREATYARREXGRATIA` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / NEWTREATYARRRICOMM` |
| `NewTreatyArrExGratiaList.xml` | `@BASECLASS` | `NEWTREATYARREXGRATIALIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / NEWTREATYARRTREATYLIMITLIST` |
| `NewTreatyArrExclutionTreaty.xml` | `@BASECLASS` | `NEWTREATYARREXCLUTIONTREATY` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / NEWTREATYARRTERRLIMIT` |
| `NewTreatyArrFacIn.xml` | `@BASECLASS` | `NEWTREATYARRFACIN` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / NEWTREATYARREPI` |
| `NewTreatyArrFacInList.xml` | `@BASECLASS` | `NEWTREATYARRFACINLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / NEWTREATYARREPILIST` |
| `NewTreatyArrLimitMB.xml` | `@BASECLASS` | `NEWTREATYARRLIMITMB` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `NewTreatyArrMaxCoinsPanel.xml` | `@BASECLASS` | `NEWTREATYARRMAXCOINSPANEL` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / NEWTREATYARRMINLOL` |
| `NewTreatyArrMinLOL.xml` | `@BASECLASS` | `NEWTREATYARRMINLOL` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / NEWTREATYARREXCLUTIONTREATY` |
| `NewTreatyArrMinLOLMB.xml` | `@BASECLASS` | `NEWTREATYARRMINLOLMB` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / NEWTREATYARRMINLOL` |
| `NewTreatyArrPLA.xml` | `@BASECLASS` | `NEWTREATYARRPLA` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / NEWTREATYARRTREATYLIMIT` |
| `NewTreatyArrPLAList.xml` | `@BASECLASS` | `NEWTREATYARRPLALIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `@BASECLASS / NEWTREATYARRTREATYLIMITLIST` |
| `NewTreatyArrProfitComm.xml` | `@BASECLASS` | `NEWTREATYARRPROFITCOMM` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `@BASECLASS / NEWTREATYARRBORDEREAUX` |
| `NewTreatyArrRicomm.xml` | `@BASECLASS` | `NEWTREATYARRRICOMM` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `@BASECLASS / NEWTREATYARRTERRLIMIT` |
| `NewTreatyArrTerrLimit.xml` | `@BASECLASS` | `NEWTREATYARRTERRLIMIT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `NewTreatyArrTreatyLimit.xml` | `@BASECLASS` | `NEWTREATYARRTREATYLIMIT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / NEWTREATYARRRICOMM` |
| `NewTreatyArrTreatyLimitList.xml` | `@BASECLASS` | `NEWTREATYARRTREATYLIMITLIST` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `NewTreatyBusinessDetail_Act.xml` | `@BASECLASS` | `NEWTREATYBUSINESSDETAIL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | — |
| `NewTreatyReinsurerDetail_Act.xml` | `@BASECLASS` | `NEWTREATYREINSURERDETAIL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `PanggilID.xml` | `@BASECLASS` | `PANGGILID` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | `@BASECLASS / BROWSEPORTFOLIOLIST_ACT` |
| `RefreshErrorProportionalarrg.xml` | `@BASECLASS` | `REFRESHERRORPROPORTIONALARRG` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `RefreshKurs.xml` | `@BASECLASS` | `REFRESHKURS` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `RefreshYear.xml` | `@BASECLASS` | `REFRESHYEAR` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SaveSecurityReinsurer_Act.xml` | `@BASECLASS` | `SAVESECURITYREINSURER_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SaveTreatyArrBordereAux_Act.xml` | `@BASECLASS` | `SAVETREATYARRBORDEREAUX_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | `@BASECLASS / SAVETREATYARRRICOMM_ACT` |
| `SaveTreatyArrCashLossLimitList_Act.xml` | `@BASECLASS` | `SAVETREATYARRCASHLOSSLIMITLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | `@BASECLASS / SAVETREATYARRPLALIST_ACT` |
| `SaveTreatyArrCashLossLimit_Act.xml` | `@BASECLASS` | `SAVETREATYARRCASHLOSSLIMIT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | `@BASECLASS / SAVETREATYARRPLA_ACT` |
| `SaveTreatyArrClaimCoorpChild_Act.xml` | `@BASECLASS` | `SAVETREATYARRCLAIMCOORPCHILD_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `@BASECLASS / SAVETREATYARRTREATYLIMITCHILD_ACT` |
| `SaveTreatyArrClaimCoorp_Act.xml` | `@BASECLASS` | `SAVETREATYARRCLAIMCOORP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | `@BASECLASS / SAVETREATYARRBORDEREAUX_ACT` |
| `SaveTreatyArrCoinsPanel_Act.xml` | `@BASECLASS` | `SAVETREATYARRCOINSPANEL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | `@BASECLASS / SAVETREATYARREXCLUTIONTREATY_ACT` |
| `SaveTreatyArrEPI_Act.xml` | `@BASECLASS` | `SAVETREATYARREPI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | `@BASECLASS / SAVETREATYARRCASHLOSSLIMIT_ACT` |
| `SaveTreatyArrEpiList_Act.xml` | `@BASECLASS` | `SAVETREATYARREPILIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `@BASECLASS / SAVETREATYARRCASHLOSSLIMITLIST_ACT` |
| `SaveTreatyArrExGratiaChildList_Act.xml` | `@BASECLASS` | `SAVETREATYARREXGRATIACHILDLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `@BASECLASS / SAVETREATYARRPLALIST_ACT` |
| `SaveTreatyArrExGratia_Act.xml` | `@BASECLASS` | `SAVETREATYARREXGRATIA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `@BASECLASS / SAVETREATYARRRICOMM_ACT` |
| `SaveTreatyArrExclutionTreaty_Act.xml` | `@BASECLASS` | `SAVETREATYARREXCLUTIONTREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | `@BASECLASS / SAVETREATYARRTERRLIMIT_ACT` |
| `SaveTreatyArrFacInList_Act.xml` | `@BASECLASS` | `SAVETREATYARRFACINLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=14 | `@BASECLASS / SAVETREATYARREPILIST_ACT` |
| `SaveTreatyArrFacIn_Act.xml` | `@BASECLASS` | `SAVETREATYARRFACIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `@BASECLASS / SAVETREATYARREPI_ACT` |
| `SaveTreatyArrLimitMB_Act.xml` | `@BASECLASS` | `SAVETREATYARRLIMITMB_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SaveTreatyArrMaxCoinsPanel.xml` | `@BASECLASS` | `SAVETREATYARRMAXCOINSPANEL` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | `@BASECLASS / SAVETREATYARRMINLOL` |
| `SaveTreatyArrMinLOL.xml` | `@BASECLASS` | `SAVETREATYARRMINLOL` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | `@BASECLASS / SAVETREATYARREXCLUTIONTREATY_ACT` |
| `SaveTreatyArrMinLOLMB.xml` | `@BASECLASS` | `SAVETREATYARRMINLOLMB` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | `@BASECLASS / SAVETREATYARRMINLOL` |
| `SaveTreatyArrPLAList_Act.xml` | `@BASECLASS` | `SAVETREATYARRPLALIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=15 | `@BASECLASS / SAVETREATYARRTREATYLIMITCHILD_ACT` |
| `SaveTreatyArrPLA_Act.xml` | `@BASECLASS` | `SAVETREATYARRPLA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | `@BASECLASS / SAVETREATYARRTREATYLIMIT_ACT` |
| `SaveTreatyArrPortfolio_Act.xml` | `@BASECLASS` | `SAVETREATYARRPORTFOLIO_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `@BASECLASS / SAVETREATYARREXGRATIA_ACT` |
| `SaveTreatyArrProfitComm_Act.xml` | `@BASECLASS` | `SAVETREATYARRPROFITCOMM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | `@BASECLASS / SAVETREATYARRBORDEREAUX_ACT` |
| `SaveTreatyArrRicomm_Act.xml` | `@BASECLASS` | `SAVETREATYARRRICOMM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `@BASECLASS / SAVETREATYARRTERRLIMIT_ACT` |
| `SaveTreatyArrTerrLimit_Act.xml` | `@BASECLASS` | `SAVETREATYARRTERRLIMIT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=8 | `@BASECLASS / SAVETREATYGROUP_ACT` |
| `SaveTreatyArrTreatyLimitChild_Act.xml` | `@BASECLASS` | `SAVETREATYARRTREATYLIMITCHILD_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=18 | — |
| `SaveTreatyArrTreatyLimit_Act.xml` | `@BASECLASS` | `SAVETREATYARRTREATYLIMIT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `@BASECLASS / SAVETREATYARRRICOMM_ACT` |
| `SaveTreatyBusinessDetail_Act.xml` | `@BASECLASS` | `SAVETREATYBUSINESSDETAIL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=7 | — |
| `SaveTreatyContract_Act.xml` | `@BASECLASS` | `SAVETREATYCONTRACT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `SaveTreatyReinsurerDetail1_Act.xml` | `@BASECLASS` | `SAVETREATYREINSURERDETAIL1_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=19 | `@BASECLASS / SAVETREATYREINSURERDETAIL_ACT` |
| `SaveTreatyYearMultiple_Act.xml` | `@BASECLASS` | `SAVETREATYYEARMULTIPLE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `@BASECLASS / SAVETREATYYEAR_ACT` |
| `SaveTreatyYear_Act.xml` | `@BASECLASS` | `SAVETREATYYEAR_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `@BASECLASS / SAVETREATYCONTRACT_ACT` |
| `SetCategoryAttachTreatyin.xml` | `ASM-FW-GISFW-DATA-OFFERFACIN` | `SETCATEGORYATTACHTREATYIN` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `ASM-FW-GISFW-DATA-OFFERFACIN / SETCATEGORYATTACH` |
| `SetCategory_act.xml` | `DATA-PORTAL` | `SETCATEGORY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetErrorMessage.xml` | `@BASECLASS` | `SETERRORMESSAGE` | [terverifikasi] pyActivityType=ACTIVITY; steps=10 | — |
| `SetErrorMessageReinsurer.xml` | `@BASECLASS` | `SETERRORMESSAGEREINSURER` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |
| `SetKirimIDDesc.xml` | `@BASECLASS` | `SETKIRIMIDDESC` | [terverifikasi] pyActivityType=ACTIVITY; steps=28 | — |
| `SetOccupationID.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `SETOCCUPATIONID` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetOccupationLimitMB.xml` | `@BASECLASS` | `SETOCCUPATIONLIMITMB` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SetTanggalTreatyContract.xml` | `@BASECLASS` | `SETTANGGALTREATYCONTRACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | — |
| `SetTreatyArrCashLossLimitList_Act.xml` | `@BASECLASS` | `SETTREATYARRCASHLOSSLIMITLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `SetTreatyArrCashLossLimit_Act.xml` | `@BASECLASS` | `SETTREATYARRCASHLOSSLIMIT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `@BASECLASS / SETTREATYARRPLA_ACT` |
| `SetTreatyArrClaimCoorpList_Act.xml` | `@BASECLASS` | `SETTREATYARRCLAIMCOORPLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `SetTreatyArrClaimCoorp_Act.xml` | `@BASECLASS` | `SETTREATYARRCLAIMCOORP_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `@BASECLASS / SETTREATYARRTERRBORD_ACT` |
| `SetTreatyArrEpiList_Act.xml` | `@BASECLASS` | `SETTREATYARREPILIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `SetTreatyArrEpi_Act.xml` | `@BASECLASS` | `SETTREATYARREPI_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `@BASECLASS / SETTREATYARRCLAIMCOORP_ACT` |
| `SetTreatyArrExGratiaList_Act.xml` | `@BASECLASS` | `SETTREATYARREXGRATIALIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SetTreatyArrExGratia_Act.xml` | `@BASECLASS` | `SETTREATYARREXGRATIA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `@BASECLASS / SETTREATYARRRICOMM_ACT` |
| `SetTreatyArrExclustionCoins_Act.xml` | `@BASECLASS` | `SETTREATYARREXCLUSTIONCOINS_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `@BASECLASS / SETTREATYARREXCLUSTIONTREATYOBJECT_ACT` |
| `SetTreatyArrExclustionLimitLimitMB_Act.xml` | `@BASECLASS` | `SETTREATYARREXCLUSTIONLIMITLIMITMB_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SetTreatyArrExclustionTreatyClause_Act.xml` | `@BASECLASS` | `SETTREATYARREXCLUSTIONTREATYCLAUSE_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `@BASECLASS / SETTREATYARREXCLUSTIONCOINS_ACT` |
| `SetTreatyArrExclustionTreatyOccupation_Act.xml` | `@BASECLASS` | `SETTREATYARREXCLUSTIONTREATYOCCUPATION_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `@BASECLASS / SETTREATYARREXCLUSTIONTREATYCLAUSE_ACT` |
| `SetTreatyArrExclutionTreaty_Act.xml` | `@BASECLASS` | `SETTREATYARREXCLUTIONTREATY_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / SETTREATYARREXCLUSTIONCOINS_ACT` |
| `SetTreatyArrFacInList_Act.xml` | `@BASECLASS` | `SETTREATYARRFACINLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SetTreatyArrFacIn_Act.xml` | `@BASECLASS` | `SETTREATYARRFACIN_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `@BASECLASS / SETTREATYARREPI_ACT` |
| `SetTreatyArrPLAList_Act.xml` | `@BASECLASS` | `SETTREATYARRPLALIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `SetTreatyArrPLA_Act.xml` | `@BASECLASS` | `SETTREATYARRPLA_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `@BASECLASS / SETTREATYARRTREATYLIMIT_ACT` |
| `SetTreatyArrProfitComm_Act.xml` | `@BASECLASS` | `SETTREATYARRPROFITCOMM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SetTreatyArrRicomm_Act.xml` | `@BASECLASS` | `SETTREATYARRRICOMM_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetTreatyArrTerrBord_Act.xml` | `@BASECLASS` | `SETTREATYARRTERRBORD_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | `@BASECLASS / SETTREATYARRTERRLIMIT_ACT` |
| `SetTreatyArrTerrLimit_Act.xml` | `@BASECLASS` | `SETTREATYARRTERRLIMIT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | `@BASECLASS / SETTREATYGROUPDATA_ACT` |
| `SetTreatyArrTreatyLimitList_Act.xml` | `@BASECLASS` | `SETTREATYARRTREATYLIMITLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `SetTreatyArrTreatyLimit_Act.xml` | `@BASECLASS` | `SETTREATYARRTREATYLIMIT_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / SETTREATYARRRICOMM_ACT` |
| `SetTreatyArrangementDesc_Act.xml` | `@BASECLASS` | `SETTREATYARRANGEMENTDESC_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `SetTreatyBusinessList_Act.xml` | `@BASECLASS` | `SETTREATYBUSINESSLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | `@BASECLASS / SETTREATYGROUPDATA_ACT` |
| `SetTreatyReinsurerList_Act.xml` | `@BASECLASS` | `SETTREATYREINSURERLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=5 | — |
| `SetTreatyYear.xml` | `@BASECLASS` | `SETTREATYYEAR` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetTreatyYear_Act.xml` | `@BASECLASS` | `SETTREATYYEAR_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SetUbahTreatyBusinessList_Act.xml` | `@BASECLASS` | `SETUBAHTREATYBUSINESSLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `SetUbahTreatyContract.xml` | `@BASECLASS` | `SETUBAHTREATYCONTRACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `SetUbahTreatyReinsurerList_Act.xml` | `@BASECLASS` | `SETUBAHTREATYREINSURERLIST_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | — |
| `ShowEditSecurityReinsurer.xml` | `@BASECLASS` | `SHOWEDITSECURITYREINSURER` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `TreatyInitAttach.xml` | `@BASECLASS` | `TREATYINITATTACH` | [terverifikasi] pyActivityType=ACTIVITY; steps=4 | `@BASECLASS / SFAGISINITATTACH` |
| `TreatyOutDownloadAll_Act.xml` | `@BASECLASS` | `TREATYOUTDOWNLOADALL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `TreatyOutDownloadOne.xml` | `DATA-PORTAL` | `TREATYOUTDOWNLOADONE` | [terverifikasi] pyActivityType=ACTIVITY; steps=2 | `DATA-PORTAL / TREATYINDOWNLOADONE` |
| `TreatyOutSaveAttachment.xml` | `@BASECLASS` | `TREATYOUTSAVEATTACHMENT` | [terverifikasi] pyActivityType=ACTIVITY; steps=9 | `@BASECLASS / TREATYSAVEATTACHMENT` |
| `TreatyTestChildTotal_Act.xml` | `@BASECLASS` | `TREATYTESTCHILDTOTAL_ACT` | [terverifikasi] pyActivityType=ACTIVITY; steps=6 | — |
| `kosongkanListProportionalArrg.xml` | `@BASECLASS` | `KOSONGKANLISTPROPORTIONALARRG` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `setValueTreatyDesc.xml` | `@BASECLASS` | `SETVALUETREATYDESC` | [terverifikasi] pyActivityType=ACTIVITY; steps=1 | — |
| `testingKurs.xml` | `@BASECLASS` | `TESTINGKURS` | [terverifikasi] pyActivityType=ACTIVITY; steps=3 | — |

### ConnectREST — 1 rule (`RULE-CONNECT-REST`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `SERVICEGOOGLE` | [terverifikasi] pyServiceName=ServiceGoogle; baseURL=SETTING; setting=LinkService!LinkService; resolusi=JNDIName | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / GOOGLESTORAGE_UPLOAD` |

### DataPage — 2 rule (`RULE-DECLARE-PAGES`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `D_EnumerationList.xml` | — | `D_ENUMERATIONLIST` | [terverifikasi] pyPageName=D_EnumerationList; scope=thread; class=ASM-FW-GISFW-Data-Enumeration; struktur=list | `D_ENUMERATIONLIST / #20151201T084130.039` |
| `D_TreatyContract.xml` | — | `D_TREATYCONTRACT` | [terverifikasi] pyPageName=D_TreatyContract; scope=thread; class=ASM-FW-GISFW-Int-TREATYCONTRACT; struktur=list | `D_TREATYCONTRACT / #20250115T041117.614` |

### DataTransform — 10 rule (`RULE-OBJ-MODEL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `CancelDatShow.xml` | `@BASECLASS` | `CANCELDATSHOW` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridTreatyArrangementRetentionEdit_DT.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTRETENTIONEDIT_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridTreatyArrangementRetentionUSD_DT.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTRETENTIONUSD_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputNewSecurityReinsurer.xml` | `@BASECLASS` | `INPUTNEWSECURITYREINSURER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetOperatorName_DT.xml` | `@BASECLASS` | `SETOPERATORNAME_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetOutputParam1_DT.xml` | `@BASECLASS` | `SETOUTPUTPARAM1_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetOutputParam_DT.xml` | `@BASECLASS` | `SETOUTPUTPARAM_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetSecurityReinsurer.xml` | `@BASECLASS` | `SETSECURITYREINSURER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / SETOUTPUTPARAM1_DT` |
| `SetTreatyContractReinsTypeEnddate.xml` | `@BASECLASS` | `SETTREATYCONTRACTREINSTYPEENDDATE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SetTreatyYearTreatyContract_DT.xml` | `@BASECLASS` | `SETTREATYYEARTREATYCONTRACT_DT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### FlowAction — 2 rule (`RULE-OBJ-FLOWACTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `DetailTreatyExclustion.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `DETAILTREATYEXCLUSTION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `TreatyOutAttachContent.xml` | `@BASECLASS` | `TREATYOUTATTACHCONTENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / TREATYATTACHCONTENT` |

### Harness — 3 rule (`RULE-HTML-HARNESS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `InboxTreatyContract.xml` | `DATA-PORTAL` | `INBOXTREATYCONTRACT` | [terverifikasi] pyInclude=3 | — |
| `InboxTreatyContractDescription.xml` | `DATA-PORTAL` | `INBOXTREATYCONTRACTDESCRIPTION` | [terverifikasi] pyInclude=9 | — |
| `InboxTreatyContractReinsType.xml` | `DATA-PORTAL` | `INBOXTREATYCONTRACTREINSTYPE` | [terverifikasi] pyInclude=3 | `DATA-PORTAL / INBOXTREATYCONTRACT` |

### RDBList — 37 rule (`RULE-CONNECT-SQL`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `CategoryAttach_SQL.xml` | `ASM-FW-GISFW-INT-OFFERJSON` | `CATEGORYATTACH_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=CATEGORY_ATTACH_REAS | — |
| `DeleteAttachment2_Sql.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `DELETEATTACHMENT2_SQL` | [terverifikasi] sqlKind=QUERY | `ASM-FW-GISFW-INT-TREATY_IN / ASM!GETATTACHMENT2_SQL` |
| `DeleteFromTREATYCONTRACT_SQL.xml` | `ASM-FW-GISFW-INT-TREATYCONTRACT` | `DELETEFROMTREATYCONTRACT_SQL` | [terverifikasi] sqlKind=PLSQL; sqlOps=SELECT,DELETE; tables=TREATYCONTRACT,TREATYBUSINESS,MTREATYSECURITY,TREATYREINSURER | — |
| `DeleteFromTreatyReinsurer_Act.xml` | `ASM-FW-GISFW-INT-TREATYREINSURER` | `DELETEFROMTREATYREINSURER_ACT` | [terverifikasi] sqlKind=PLSQL; sqlOps=DELETE; tables=MTREATYSECURITY,TREATYREINSURER | — |
| `DeleteRowBusinessList.xml` | `ASM-FW-GISFW-INT-TREATYBUSINESS` | `DELETEROWBUSINESSLIST` | [terverifikasi] sqlKind=QUERY; sqlOps=DELETE; tables=TREATYBUSINESS | — |
| `DeleteSecurityReinsurer.xml` | `ASM-FW-GISFW-INT-MTREATYSECURITY` | `DELETESECURITYREINSURER` | [terverifikasi] sqlKind=QUERY; sqlOps=DELETE; tables=MTREATYSECURITY | `ASM-FW-GISFW-INT-MTREATYSECURITY / ASM!UPDATEMTREATYSECURITY` |
| `DeleteStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `DELETESTORAGE_SQL` | [terverifikasi] sqlKind=QUERY | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GETLINKSTORAGE_SQL` |
| `GetAllAttachment2_Sql.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `GETALLATTACHMENT2_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_ATTACHMENTTREATY_2 | `ASM-FW-GISFW-INT-TREATY_IN / ASM!GETATTACHMENT2_SQL` |
| `GetAttachment2_Sql.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `GETATTACHMENT2_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_ATTACHMENTTREATY_2 | `ASM-FW-GISFW-INT-TREATY_IN / ASM!GETATTACHMENT_SQL` |
| `GetLinkStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETLINKSTORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!GETTOKENSTORAGE_SQL` |
| `GetMasterBusinessList.xml` | `ASM-FW-GISFW-INT-TREATYBUSINESS` | `GETMASTERBUSINESSLIST` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYBUSINESS | `ASM-FW-GISFW-INT-TREATYCONTRACT / ASM!GETMASTERTREATYCONTRACTBUSINESSLIST` |
| `GetMasterDescriptionCashParentList.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `GETMASTERDESCRIPTIONCASHPARENTLIST` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_PROPORTIONALARRG | `ASM-FW-GISFW-INT-PROPORTIONALARRG / ASM!GETMASTERDESCRIPTIONLIMITPARENTLIST` |
| `GetMasterDescriptionClaimCoorpParentList.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `GETMASTERDESCRIPTIONCLAIMCOORPPARENTLIST` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_PROPORTIONALARRG | `ASM-FW-GISFW-INT-PROPORTIONALARRG / ASM!GETMASTERDESCRIPTIONLIMITPARENTLIST` |
| `GetMasterDescriptionEPIParentList.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `GETMASTERDESCRIPTIONEPIPARENTLIST` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_PROPORTIONALARRG | `ASM-FW-GISFW-INT-PROPORTIONALARRG / ASM!GETMASTERDESCRIPTIONCLAIMCOORPPARENTLIST` |
| `GetMasterDescriptionExGratiaChildList.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `GETMASTERDESCRIPTIONEXGRATIACHILDLIST` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_PROPORTIONALARRG | `ASM-FW-GISFW-INT-PROPORTIONALARRG / ASM!GETMASTERDESCRIPTIONEXGRATIALIST` |
| `GetMasterDescriptionExGratiaList.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `GETMASTERDESCRIPTIONEXGRATIALIST` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_PROPORTIONALARRG | `ASM-FW-GISFW-INT-PROPORTIONALARRG / ASM!GETMASTERDESCRIPTIONLIMITLIST` |
| `GetMasterDescriptionFACINParentList.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `GETMASTERDESCRIPTIONFACINPARENTLIST` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_PROPORTIONALARRG | `ASM-FW-GISFW-INT-PROPORTIONALARRG / ASM!GETMASTERDESCRIPTIONPLAPARENTLIST` |
| `GetMasterDescriptionLimitList.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `GETMASTERDESCRIPTIONLIMITLIST` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_PROPORTIONALARRG | `ASM-FW-GISFW-INT-TREATYBUSINESS / ASM!GETMASTERBUSINESSLIST` |
| `GetMasterDescriptionLimitParentList.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `GETMASTERDESCRIPTIONLIMITPARENTLIST` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=PROPORTIONALARRG | `ASM-FW-GISFW-INT-PROPORTIONALARRG / ASM!GETMASTERDESCRIPTIONLIMITLIST` |
| `GetMasterDescriptionPLAParentList.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `GETMASTERDESCRIPTIONPLAPARENTLIST` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_PROPORTIONALARRG | `ASM-FW-GISFW-INT-PROPORTIONALARRG / ASM!GETMASTERDESCRIPTIONLIMITPARENTLIST` |
| `GetMasterKursList.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `GETMASTERKURSLIST` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYEXCHANGE | `ASM-FW-GISFW-INT-PROPORTIONALARRG / ASM!GETMASTERTREATYYEARLIST` |
| `GetMasterPanggilID.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `GETMASTERPANGGILID` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_PROPORTIONALARRG | `ASM-FW-GISFW-INT-PROPORTIONALARRG / ASM!GETMASTERDESCRIPTIONEXGRATIACHILDLIST` |
| `GetMasterPortfolioListDetail.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `GETMASTERPORTFOLIOLISTDETAIL` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_PROPORTIONALARRG | `ASM-FW-GISFW-INT-PLANTRAVEL / ASM!GETMASTERPORTFOLIOLISTDETAIL` |
| `GetMasterReinsTypeContract.xml` | `ASM-FW-GISFW-INT-TREATYCONTRACT` | `GETMASTERREINSTYPECONTRACT` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=M_TREATYYEAR | `ASM-FW-GISFW-INT-PROPORTIONALARRG / ASM!GETMASTERDESCRIPTIONLIMITLIST` |
| `GetMasterReinsurerList.xml` | `ASM-FW-GISFW-INT-TREATYREINSURER` | `GETMASTERREINSURERLIST` | [terverifikasi] sqlKind=QUERY; sqlOps=SELECT; tables=TREATYREINSURER | — |
| `GetTokenStorage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `GETTOKENSTORAGE_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.GET_TOKEN_STORAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |
| `InsertAtatchment_Sql.xml` | `ASM-FW-GISFW-INT-TREATY_IN` | `INSERTATATCHMENT_SQL` | [terverifikasi] sqlKind=PLSQL; procs=DBMS_LOB.CREATETEMPORARY,POOLDATA.PEGA_M_ATTACHMENT | — |
| `InsertToMTreatySecurity.xml` | `ASM-FW-GISFW-INT-MTREATYSECURITY` | `INSERTTOMTREATYSECURITY` | [terverifikasi] sqlKind=QUERY; sqlOps=INSERT; tables=MTREATYSECURITY | — |
| `SaveMasterCopyData_SQL.xml` | `ASM-FW-GISFW-INT-TREATYYEAR` | `SAVEMASTERCOPYDATA_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PROSESCOPY,DBMS_OUTPUT.PUT_LINE | — |
| `SaveMasterProportionalArrg.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `SAVEMASTERPROPORTIONALARRG` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_PROPORTIONALARRG | `ASM-FW-GISFW-INT-PROPORTIONALARRG / ASM!GETMASTERDESCRIPTIONLIMITLIST` |
| `SaveMasterProportionalArrgChild.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `SAVEMASTERPROPORTIONALARRGCHILD` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_M_PROPORTIONALARRG_CHILD | `ASM-FW-GISFW-INT-PROPORTIONALARRG / ASM!SAVEMASTERPROPORTIONALARRG` |
| `SaveMasterTreatyBusiness_SQL.xml` | `ASM-FW-GISFW-INT-TREATYBUSINESS` | `SAVEMASTERTREATYBUSINESS_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_TREATYBUSINESS | `ASM-FW-GRSFW-INT-TREATYBUSINESS / ASM!SAVEMASTERTREATYBUSINESS_SQL` |
| `SaveMasterTreatyContract_SQL.xml` | `ASM-FW-GISFW-INT-TREATYCONTRACT` | `SAVEMASTERTREATYCONTRACT_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_TREATYCONTRACT | `ASM-FW-GRSFW-INT-TREATYCONTRACT / ASM!SAVEMASTERTREATYCONTRACT_SQL` |
| `SaveMasterTreatyReinsurer_SQL.xml` | `ASM-FW-GISFW-INT-TREATYREINSURER` | `SAVEMASTERTREATYREINSURER_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_TREATYREINSURER | — |
| `SaveMasterTreatyYear_SQL.xml` | `ASM-FW-GISFW-INT-TREATYYEAR` | `SAVEMASTERTREATYYEAR_SQL` | [terverifikasi] sqlKind=PLSQL; procs=POOLDATA.PEGA_TREATYYEAR | `ASM-FW-GISFW-INT-TREATYCONTRACT / ASM!SAVEMASTERTREATYCONTRACT_SQL` |
| `UpdateMTreatySecurity.xml` | `ASM-FW-GISFW-INT-MTREATYSECURITY` | `UPDATEMTREATYSECURITY` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=MTREATYSECURITY | `ASM-FW-GISFW-INT-MTREATYSECURITY / ASM!INSERTTOMTREATYSECURITY` |
| `Update_T_Storage_SQL.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `UPDATE_T_STORAGE_SQL` | [terverifikasi] sqlKind=QUERY; sqlOps=UPDATE; tables=T_STORAGE_IMAGE | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` |

### ReportDefinition — 30 rule (`RULE-OBJ-REPORT-DEFINITION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `BrowseAgentReinsSOA_RD.xml` | `ASM-FW-GISFW-INT-AGENT` | `BROWSEAGENTREINSSOA_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseDetailTreatyReisurer_RD.xml` | `ASM-FW-GISFW-INT-TREATYREINSURER` | `BROWSEDETAILTREATYREISURER_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-TREATYREINSURER / BROWSEDETAILTREAYREISURER_RD` |
| `BrowseFilterBusiness_RD.xml` | `ASM-FW-GISFW-INT-BUSINESS` | `BROWSEFILTERBUSINESS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseFireClauseFacIn_RD.xml` | `ASM-FW-GISFW-INT-CLAUSE` | `BROWSEFIRECLAUSEFACIN_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-CLAUSE / BROWSEFIRECLAUSE_RD` |
| `BrowseOccupationFIRE_RD.xml` | `ASM-FW-GISFW-INT-OCCUPATION` | `BROWSEOCCUPATIONFIRE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseReinsuranceType_RD.xml` | `ASM-FW-GISFW-INT-REINSURANCETYPE` | `BROWSEREINSURANCETYPE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseReinsuranceType_RD_Old_Ljt_id_isnotnull.xml` | `ASM-FW-GISFW-INT-REINSURANCETYPE` | `BROWSEREINSURANCETYPE_RD_OLD_LJT_ID_ISNOTNULL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-REINSURANCETYPE / BROWSEREINSURANCETYPE_RD` |
| `BrowseTreatyArrangement_BordereAux_RD.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRANGEMENT_BORDEREAUX_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyArrangement_ClaimCoorp_RD.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRANGEMENT_CLAIMCOORP_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyArrangement_CoinsPanel_RD.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRANGEMENT_COINSPANEL_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PROPORTIONALARRG / BROWSETREATYARRANGEMENT_EXCLUTIONTREATYOCCUPATION_RD` |
| `BrowseTreatyArrangement_EPI_RD.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRANGEMENT_EPI_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyArrangement_ExGratiaChild_RD.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRANGEMENT_EXGRATIACHILD_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PROPORTIONALARRG / BROWSETREATYARRANGEMENT_LIMIT_RD` |
| `BrowseTreatyArrangement_ExGratiaList_RD.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRANGEMENT_EXGRATIALIST_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PROPORTIONALARRG / BROWSETREATYARRANGEMENT_EXGRATIA_RD` |
| `BrowseTreatyArrangement_ExclutionTreatyClause_RD.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRANGEMENT_EXCLUTIONTREATYCLAUSE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PROPORTIONALARRG / BROWSETREATYARRANGEMENT_EXCLUTIONTREATYOCCUPATION_RD` |
| `BrowseTreatyArrangement_ExclutionTreatyOccupation_ClassCont_RD.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRANGEMENT_EXCLUTIONTREATYOCCUPATION_CLASSCONT_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PROPORTIONALARRG / BROWSETREATYARRANGEMENT_EXCLUTIONTREATYOCCUPATION_RD` |
| `BrowseTreatyArrangement_ExclutionTreatyOccupation_RD.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRANGEMENT_EXCLUTIONTREATYOCCUPATION_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PROPORTIONALARRG / BROWSETREATYARRANGEMENT_EXCLUTIONTREATY_RD` |
| `BrowseTreatyArrangement_FacIn_RD.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRANGEMENT_FACIN_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyArrangement_LimitMB_RD.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRANGEMENT_LIMITMB_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyArrangement_Limit_RD.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRANGEMENT_LIMIT_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyArrangement_ParentReins.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRANGEMENT_PARENTREINS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyArrangement_Profit.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRANGEMENT_PROFIT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-PROPORTIONALARRG / BROWSETREATYARRANGEMENT_PARENTREINS` |
| `BrowseTreatyArrangement_Ricomm_RD.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRANGEMENT_RICOMM_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyArrangement_TerritorialLimit_RD.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `BROWSETREATYARRANGEMENT_TERRITORIALLIMIT_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyBusiness_RD.xml` | `ASM-FW-GISFW-INT-TREATYBUSINESS` | `BROWSETREATYBUSINESS_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyDesc_RD.xml` | `ASM-FW-GISFW-INT-TREATYDESC` | `BROWSETREATYDESC_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyGroup_RD.xml` | `ASM-FW-GISFW-INT-TREATYGROUP` | `BROWSETREATYGROUP_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `BrowseTreatyYearMultiple_RD.xml` | `ASM-FW-GISFW-INT-TREATYYEAR` | `BROWSETREATYYEARMULTIPLE_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `ASM-FW-GISFW-INT-TREATYYEAR / BROWSETREATYYEAR_RD` |
| `BrowseTreatyYear_RD.xml` | `ASM-FW-GISFW-INT-TREATYYEAR` | `BROWSETREATYYEAR_RD` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DataTableEditorReport.xml` | `ASM-FW-GISFW-DATA-ENUMERATION` | `DATATABLEEDITORREPORT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `SelectSecurityReinsurer.xml` | `ASM-FW-GISFW-INT-MTREATYSECURITY` | `SELECTSECURITYREINSURER` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |

### Section — 49 rule (`RULE-HTML-SECTION`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `DetailCoinsShare.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `DETAILCOINSSHARE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `DetailTreatyExclustion_Sec.xml` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` | `DETAILTREATYEXCLUSTION_SEC` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridTreatyArrTreatyCashLossLimitList.xml` | `@BASECLASS` | `GRIDTREATYARRTREATYCASHLOSSLIMITLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridTreatyArrTreatyEpiList.xml` | `@BASECLASS` | `GRIDTREATYARRTREATYEPILIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridTreatyArrTreatyExGratiaChildList.xml` | `@BASECLASS` | `GRIDTREATYARRTREATYEXGRATIACHILDLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridTreatyArrTreatyFacInList.xml` | `@BASECLASS` | `GRIDTREATYARRTREATYFACINLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridTreatyArrTreatyLimitList.xml` | `@BASECLASS` | `GRIDTREATYARRTREATYLIMITLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridTreatyArrTreatyPLAList.xml` | `@BASECLASS` | `GRIDTREATYARRTREATYPLALIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridTreatyArrangementAttachment.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTATTACHMENT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridTreatyArrangementBordereAux.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTBORDEREAUX` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridTreatyArrangementCashLossLimit.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTCASHLOSSLIMIT` | [terverifikasi] pyInclude=2 | — |
| `GridTreatyArrangementClaimCoorp.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTCLAIMCOORP` | [terverifikasi] pyInclude=2 | — |
| `GridTreatyArrangementClaimCoorpList.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTCLAIMCOORPLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridTreatyArrangementCoins.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTCOINS` | [terverifikasi] pyInclude=2 | — |
| `GridTreatyArrangementEpi.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTEPI` | [terverifikasi] pyInclude=2 | — |
| `GridTreatyArrangementExGratia.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTEXGRATIA` | [terverifikasi] pyInclude=2 | — |
| `GridTreatyArrangementExclutionTreaty.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTEXCLUTIONTREATY` | [terverifikasi] pyInclude=8 | `@BASECLASS / GRIDTREATYARRANGEMENTTERRLIMIT` |
| `GridTreatyArrangementExclutionTreatyClausule.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTEXCLUTIONTREATYCLAUSULE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / GRIDTREATYARRANGEMENTEXCLUTIONTREATYOCCUPATION` |
| `GridTreatyArrangementExclutionTreatyObject.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTEXCLUTIONTREATYOBJECT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridTreatyArrangementExclutionTreatyOccupation.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTEXCLUTIONTREATYOCCUPATION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridTreatyArrangementExclutionTreatyPeriode.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTEXCLUTIONTREATYPERIODE` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridTreatyArrangementFacIn.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTFACIN` | [terverifikasi] pyInclude=2 | — |
| `GridTreatyArrangementLIMITMB.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTLIMITMB` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / GRIDTREATYARRANGEMENTMAXCOINSPANEL` |
| `GridTreatyArrangementMInLOLMB.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTMINLOLMB` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / GRIDTREATYARRANGEMENTMINLOL` |
| `GridTreatyArrangementMaxCoinsPanel.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTMAXCOINSPANEL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridTreatyArrangementMinLOL.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTMINLOL` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridTreatyArrangementPLA.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTPLA` | [terverifikasi] pyInclude=2 | — |
| `GridTreatyArrangementPortfolio.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTPORTFOLIO` | [terverifikasi] pyInclude=2 | `@BASECLASS / GRIDTREATYARRANGEMENTEXGRATIA` |
| `GridTreatyArrangementPortfolioList.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTPORTFOLIOLIST` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridTreatyArrangementProfitCommision.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTPROFITCOMMISION` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridTreatyArrangementRetentionNP.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTRETENTIONNP` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | `@BASECLASS / GRIDTREATYARRANGEMENTTREATYLIMITNP` |
| `GridTreatyArrangementRicomm.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTRICOMM` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridTreatyArrangementTerrLimit.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTTERRLIMIT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `GridTreatyArrangementTreatyLimit.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTTREATYLIMIT` | [terverifikasi] pyInclude=2 | — |
| `GridTreatyContract.xml` | `@BASECLASS` | `GRIDTREATYCONTRACT` | [terverifikasi] pyInclude=2 | — |
| `GridTreatyContractReinsType.xml` | `@BASECLASS` | `GRIDTREATYCONTRACTREINSTYPE` | [terverifikasi] pyInclude=2 | `@BASECLASS / GRIDTREATYCONTRACT` |
| `InputDtlTreatyContact.xml` | `@BASECLASS` | `INPUTDTLTREATYCONTACT` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `InputTreatyContract.xml` | `@BASECLASS` | `INPUTTREATYCONTRACT` | [terverifikasi] pyInclude=8 | — |
| `InputTreatyContractReinsType.xml` | `@BASECLASS` | `INPUTTREATYCONTRACTREINSTYPE` | [terverifikasi] pyInclude=4 | — |
| `NitipKurs.xml` | `@BASECLASS` | `NITIPKURS` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `PanggilReinsType.xml` | `@BASECLASS` | `PANGGILREINSTYPE` | [terverifikasi] pyInclude=2 | — |
| `SubViewDetailDescription.xml` | `@BASECLASS` | `SUBVIEWDETAILDESCRIPTION` | [terverifikasi] pyInclude=40 | `@BASECLASS / VIEWDETAILDESCRIPTION` |
| `ViewDetailDescription.xml` | `@BASECLASS` | `VIEWDETAILDESCRIPTION` | [terverifikasi] pyInclude=8 | — |
| `ViewDetailDescription2.xml` | `@BASECLASS` | `VIEWDETAILDESCRIPTION2` | [terverifikasi] pyInclude=24 | `@BASECLASS / VIEWDETAILDESCRIPTION` |
| `ViewDetailDescriptionNonProp.xml` | `@BASECLASS` | `VIEWDETAILDESCRIPTIONNONPROP` | [terverifikasi] pyInclude=4 | `@BASECLASS / VIEWDETAILDESCRIPTIONPROP` |
| `ViewDetailDescriptionProp.xml` | `@BASECLASS` | `VIEWDETAILDESCRIPTIONPROP` | [terverifikasi] pyInclude=24 | — |
| `ViewDetailTreatyBusinessGrid.xml` | `@BASECLASS` | `VIEWDETAILTREATYBUSINESSGRID` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `ViewDetailTreatyReinsurerGrid1.xml` | `@BASECLASS` | `VIEWDETAILTREATYREINSURERGRID1` | [dugaan] peran belum dikonfirmasi tag; hanya dari nama file + tipe rule | — |
| `gridTreatyArrangementExGratiaList.xml` | `@BASECLASS` | `GRIDTREATYARRANGEMENTEXGRATIALIST` | [terverifikasi] pyInclude=2 | — |

### SystemSettings — 1 rule (`RULE-ADMIN-SYSTEM-SETTINGS`)

| File | Class | Nama rule | Fakta dari tag | Asal salinan (bila beda) |
| --- | --- | --- | --- | --- |
| `LinkService.xml` | `LINKSERVICE` | `LINKSERVICE` | [terverifikasi] pyPurpose=LinkService | — |

## 4. Tabrakan nama di dalam modul ini

### 4.1 Nama file dipakai di lebih dari satu tipe rule

| Nama file | Tipe | Class |
| --- | --- | --- |
| `DeleteSecurityReinsurer.xml` | Activity, RDBList | `@BASECLASS, ASM-FW-GISFW-INT-MTREATYSECURITY` |

Total: **1** nama file.

### 4.2 Rule `When` bernama sama di lebih dari satu class

_Tidak ada._

## 5. Keluarga klon (satu asal salinan dipakai >1 rule)

Dibaca dari `<pzOriginalInstanceKey>` yang sama pada beberapa file. Ini **sinyal kandidat
konsolidasi** saat migrasi, bukan bukti bahwa perilakunya identik — kesamaan perilaku baru
dapat dinyatakan setelah isinya dibandingkan di STEP D2.

| Asal salinan (class / nama) | Jumlah | File |
| --- | ---: | --- |
| `ASM-FW-GISFW-INT-MTREATYSECURITY / ASM!INSERTTOMTREATYSECURITY` | 2 | `RDBList/InsertToMTreatySecurity.xml`, `RDBList/UpdateMTreatySecurity.xml` |
| `@BASECLASS / NEWTREATYARRTERRLIMIT` | 3 | `Activity/NewTreatyArrExclutionTreaty.xml`, `Activity/NewTreatyArrRicomm.xml`, `Activity/NewTreatyArrTerrLimit.xml` |
| `ASM-FW-GISFW-INT-TREATY_IN / ASM!GETATTACHMENT2_SQL` | 2 | `RDBList/DeleteAttachment2_Sql.xml`, `RDBList/GetAllAttachment2_Sql.xml` |
| `@BASECLASS / NEWTREATYARRRICOMM` | 2 | `Activity/NewTreatyArrExGratia.xml`, `Activity/NewTreatyArrTreatyLimit.xml` |
| `@BASECLASS / VIEWDETAILDESCRIPTION` | 3 | `Section/SubViewDetailDescription.xml`, `Section/ViewDetailDescription.xml`, `Section/ViewDetailDescription2.xml` |
| `@BASECLASS / GRIDTREATYCONTRACT` | 2 | `Section/GridTreatyContract.xml`, `Section/GridTreatyContractReinsType.xml` |
| `@BASECLASS / SAVETREATYARRBORDEREAUX_ACT` | 2 | `Activity/SaveTreatyArrClaimCoorp_Act.xml`, `Activity/SaveTreatyArrProfitComm_Act.xml` |
| `ASM-FW-GISFW-INT-PROPORTIONALARRG / ASM!GETMASTERDESCRIPTIONLIMITPARENTLIST` | 3 | `RDBList/GetMasterDescriptionCashParentList.xml`, `RDBList/GetMasterDescriptionClaimCoorpParentList.xml`, `RDBList/GetMasterDescriptionPLAParentList.xml` |
| `@BASECLASS / SETTREATYARREXCLUSTIONCOINS_ACT` | 2 | `Activity/SetTreatyArrExclustionTreatyClause_Act.xml`, `Activity/SetTreatyArrExclutionTreaty_Act.xml` |
| `@BASECLASS / GRIDTREATYARRANGEMENTMINLOL` | 2 | `Section/GridTreatyArrangementMInLOLMB.xml`, `Section/GridTreatyArrangementMinLOL.xml` |
| `@BASECLASS / GRIDTREATYARRANGEMENTEXCLUTIONTREATYOCCUPATION` | 2 | `Section/GridTreatyArrangementExclutionTreatyClausule.xml`, `Section/GridTreatyArrangementExclutionTreatyOccupation.xml` |
| `@BASECLASS / NEWTREATYARRMINLOL` | 2 | `Activity/NewTreatyArrMaxCoinsPanel.xml`, `Activity/NewTreatyArrMinLOLMB.xml` |
| `@BASECLASS / SAVETREATYARRPLALIST_ACT` | 2 | `Activity/SaveTreatyArrCashLossLimitList_Act.xml`, `Activity/SaveTreatyArrExGratiaChildList_Act.xml` |
| `@BASECLASS / SETOUTPUTPARAM1_DT` | 2 | `DataTransform/SetOutputParam1_DT.xml`, `DataTransform/SetSecurityReinsurer.xml` |
| `ASM-FW-GISFW-INT-TREATYYEAR / BROWSETREATYYEAR_RD` | 2 | `ReportDefinition/BrowseTreatyYearMultiple_RD.xml`, `ReportDefinition/BrowseTreatyYear_RD.xml` |
| `ASM-FW-GISFW-INT-PROPORTIONALARRG / ASM!GETMASTERDESCRIPTIONLIMITLIST` | 4 | `RDBList/GetMasterDescriptionExGratiaList.xml`, `RDBList/GetMasterDescriptionLimitParentList.xml`, `RDBList/GetMasterReinsTypeContract.xml`, `RDBList/SaveMasterProportionalArrg.xml` |
| `@BASECLASS / SAVETREATYARRMINLOL` | 2 | `Activity/SaveTreatyArrMaxCoinsPanel.xml`, `Activity/SaveTreatyArrMinLOLMB.xml` |
| `@BASECLASS / SAVETREATYARREXCLUTIONTREATY_ACT` | 2 | `Activity/SaveTreatyArrCoinsPanel_Act.xml`, `Activity/SaveTreatyArrMinLOL.xml` |
| `@BASECLASS / SAVETREATYCONTRACT_ACT` | 2 | `Activity/SaveTreatyContract_Act.xml`, `Activity/SaveTreatyYear_Act.xml` |
| `@BASECLASS / GETMINIMUMLOL` | 2 | `Activity/GetMaxCoinsPanel.xml`, `Activity/GetMinimumLOLMB.xml` |
| `@BASECLASS / NEWINPUTTREATYCONTRACT_ACT` | 2 | `Activity/NewInputTreatyContract_Act.xml`, `Activity/NewInputTreatyYear_Act.xml` |
| `@BASECLASS / SETTREATYARRRICOMM_ACT` | 3 | `Activity/SetTreatyArrExGratia_Act.xml`, `Activity/SetTreatyArrRicomm_Act.xml`, `Activity/SetTreatyArrTreatyLimit_Act.xml` |
| `@BASECLASS / BROWSEPORTFOLIOLIST_ACT` | 2 | `Activity/BrowsePortfolioList_Act.xml`, `Activity/PanggilID.xml` |
| `@BASECLASS / NEWTREATYARREXCLUTIONTREATY` | 2 | `Activity/NewTreatyArrCoins.xml`, `Activity/NewTreatyArrMinLOL.xml` |
| `@BASECLASS / NEWTREATYARRCASHLOSSLIMIT` | 2 | `Activity/NewTreatyArrCashLossLimitList.xml`, `Activity/NewTreatyArrEpi.xml` |
| `@BASECLASS / SETTREATYGROUPDATA_ACT` | 2 | `Activity/SetTreatyArrTerrLimit_Act.xml`, `Activity/SetTreatyBusinessList_Act.xml` |
| `DATA-PORTAL / INBOXTREATYCONTRACT` | 2 | `Harness/InboxTreatyContract.xml`, `Harness/InboxTreatyContractReinsType.xml` |
| `ASM-FW-GISFW-INT-PROPORTIONALARRG / BROWSETREATYARRANGEMENT_PARENTREINS` | 2 | `ReportDefinition/BrowseTreatyArrangement_ParentReins.xml`, `ReportDefinition/BrowseTreatyArrangement_Profit.xml` |
| `ASM-FW-GISFW-INT-PROPORTIONALARRG / BROWSETREATYARRANGEMENT_EXCLUTIONTREATYOCCUPATION_RD` | 3 | `ReportDefinition/BrowseTreatyArrangement_CoinsPanel_RD.xml`, `ReportDefinition/BrowseTreatyArrangement_ExclutionTreatyClause_RD.xml`, `ReportDefinition/BrowseTreatyArrangement_ExclutionTreatyOccupation_ClassCont_RD.xml` |
| `@BASECLASS / SAVETREATYARRTREATYLIMITCHILD_ACT` | 3 | `Activity/SaveTreatyArrClaimCoorpChild_Act.xml`, `Activity/SaveTreatyArrPLAList_Act.xml`, `Activity/SaveTreatyArrTreatyLimitChild_Act.xml` |
| `ASM-FW-GISFW-INT-REINSURANCETYPE / BROWSEREINSURANCETYPE_RD` | 2 | `ReportDefinition/BrowseReinsuranceType_RD.xml`, `ReportDefinition/BrowseReinsuranceType_RD_Old_Ljt_id_isnotnull.xml` |
| `@BASECLASS / NEWTREATYARRTREATYLIMITLIST` | 5 | `Activity/NewTreatyArrClaimCoorpList.xml`, `Activity/NewTreatyArrEpiList.xml`, `Activity/NewTreatyArrExGratiaList.xml`, `Activity/NewTreatyArrPLAList.xml`, `Activity/NewTreatyArrTreatyLimitList.xml` |
| `@BASECLASS / GRIDTREATYARRANGEMENTTERRLIMIT` | 2 | `Section/GridTreatyArrangementExclutionTreaty.xml`, `Section/GridTreatyArrangementTerrLimit.xml` |
| `@BASECLASS / SAVETREATYARRRICOMM_ACT` | 3 | `Activity/SaveTreatyArrBordereAux_Act.xml`, `Activity/SaveTreatyArrExGratia_Act.xml`, `Activity/SaveTreatyArrTreatyLimit_Act.xml` |
| `@BASECLASS / SAVETREATYARRTERRLIMIT_ACT` | 2 | `Activity/SaveTreatyArrExclutionTreaty_Act.xml`, `Activity/SaveTreatyArrRicomm_Act.xml` |
| `@BASECLASS / GRIDTREATYARRANGEMENTMAXCOINSPANEL` | 2 | `Section/GridTreatyArrangementLIMITMB.xml`, `Section/GridTreatyArrangementMaxCoinsPanel.xml` |
| `ASM-FW-GISFW-INT-T_STORAGE_IMAGE / RNM!INSERT_T_STORAGE_SQL` | 2 | `RDBList/GetTokenStorage_SQL.xml`, `RDBList/Update_T_Storage_SQL.xml` |
| `ASM-FW-GISFW-INT-PROPORTIONALARRG / BROWSETREATYARRANGEMENT_LIMIT_RD` | 2 | `ReportDefinition/BrowseTreatyArrangement_ExGratiaChild_RD.xml`, `ReportDefinition/BrowseTreatyArrangement_Limit_RD.xml` |
| `ASM-FW-GISFW-INT-PROPORTIONALARRG / BROWSETREATYARRLIMITPARENTLIST` | 6 | `Activity/BrowseTreatyArrCLaimCoorpParentList.xml`, `Activity/BrowseTreatyArrCashLossParentList.xml`, `Activity/BrowseTreatyArrEpiParentList.xml`, `Activity/BrowseTreatyArrFacInParentList.xml`, `Activity/BrowseTreatyArrLimitParentList.xml`, `Activity/BrowseTreatyArrPLAParentList.xml` |
| `@BASECLASS / GETPERIODE` | 3 | `Activity/GetMinimumLOL.xml`, `Activity/GetObject.xml`, `Activity/GetPeriode.xml` |
| `@BASECLASS / VIEWDETAILDESCRIPTIONPROP` | 2 | `Section/ViewDetailDescriptionNonProp.xml`, `Section/ViewDetailDescriptionProp.xml` |
| `@BASECLASS / GRIDTREATYARRANGEMENTEXGRATIA` | 2 | `Section/GridTreatyArrangementExGratia.xml`, `Section/GridTreatyArrangementPortfolio.xml` |
| `@BASECLASS / NEWTREATYARRBORDEREAUX` | 2 | `Activity/NewTreatyArrBordereAux.xml`, `Activity/NewTreatyArrProfitComm.xml` |
| `@BASECLASS / CANCELACTIVITYPLA` | 3 | `Activity/CancelActivityFacIn.xml`, `Activity/CancelActivityFacInList.xml`, `Activity/CancelActivitypPortfolio.xml` |

## 6. Objek database yang disentuh

Diambil dari `<pyBrowseSQL>` pada rule `RDBList`. **Nama objek apa adanya dari SQL** —
kualifikasi skema tidak konsisten di korpus. **Tipe kolom, constraint, dan relasi tidak
diketahui** (OQ-001).

| Tabel/view (apa adanya di SQL) | Dirujuk oleh |
| --- | ---: |
| `M_PROPORTIONALARRG` | 10 rule |
| `MTREATYSECURITY` | 5 rule |
| `TREATYBUSINESS` | 3 rule |
| `TREATYREINSURER` | 3 rule |
| `M_ATTACHMENTTREATY_2` | 2 rule |
| `T_STORAGE_IMAGE` | 2 rule |
| `CATEGORY_ATTACH_REAS` | 1 rule |
| `M_TREATYYEAR` | 1 rule |
| `PROPORTIONALARRG` | 1 rule |
| `TREATYCONTRACT` | 1 rule |
| `TREATYEXCHANGE` | 1 rule |

### Stored procedure / function yang dipanggil

**Body prosedur ini TIDAK ADA di korpus** — logika di dalamnya tidak dapat direkonstruksi
dari XML. Lihat OQ-002.

| Procedure | Dipanggil oleh |
| --- | --- |
| `POOLDATA.PEGA_TREATYBUSINESS` | `SaveMasterTreatyBusiness_SQL.xml` |
| `POOLDATA.PEGA_M_ATTACHMENT` | `InsertAtatchment_Sql.xml` |
| `DBMS_LOB.CREATETEMPORARY` | `InsertAtatchment_Sql.xml` |
| `POOLDATA.PEGA_TREATYYEAR` | `SaveMasterTreatyYear_SQL.xml` |
| `POOLDATA.PEGA_TREATYCONTRACT` | `SaveMasterTreatyContract_SQL.xml` |
| `DBMS_OUTPUT.PUT_LINE` | `SaveMasterCopyData_SQL.xml` |
| `POOLDATA.PEGA_M_PROPORTIONALARRG_CHILD` | `SaveMasterProportionalArrgChild.xml` |
| `POOLDATA.PEGA_PROPORTIONALARRG` | `SaveMasterProportionalArrg.xml` |
| `POOLDATA.PROSESCOPY` | `SaveMasterCopyData_SQL.xml` |
| `POOLDATA.PEGA_TREATYREINSURER` | `SaveMasterTreatyReinsurer_SQL.xml` |
| `POOLDATA.GET_TOKEN_STORAGE` | `GetTokenStorage_SQL.xml` |

## 7. Integrasi eksternal (`ConnectREST`)

Alamat tujuan dicatat sebagai **konfigurasi**, bukan URL literal. `baseURL=SETTING` berarti
alamat diambil dari rule `SystemSettings` yang disebut di kolom Setting; nilai sebenarnya
**tidak ada di korpus**. Kolom terakhir menandai rule yang justru memuat URL literal —
itu temuan, bukan fakta arsitektur (lihat `_summary.md`).

| File | Class | pyServiceName | Sumber base URL | Setting | URL literal di dalam rule? |
| --- | --- | --- | --- | --- | --- |
| `ServiceGoogle.xml` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` | `ServiceGoogle` | `SETTING` | `LinkService!LinkService` | tidak |

