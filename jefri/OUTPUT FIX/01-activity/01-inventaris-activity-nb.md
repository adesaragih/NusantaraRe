# Inventaris Activity — siklus New Business (`NB FacIn/Activity/`)

**Lingkup pass ini:** HANYA `D:\migrasi\RNM\NB FacIn\`. Folder `RNW Fac In\` dan
`Endorsment Fac In\` **tidak dibuka** — tidak ada satu pun angka di dokumen ini yang berasal
dari sana.

**Metode.** Seluruh 609 berkas diurai massal dengan parser `[System.Xml.XmlDocument]`
(bukan regex teks) melalui skrip di scratchpad; hasilnya diekspor ke CSV lalu diagregasi.
Perintah audit di bawah setiap angka dapat dijalankan ulang secara mandiri.

**Berkas yang dibaca PENUH (isi mentah, bukan agregat):**

| Berkas | Alasan |
| --- | --- |
| `NB FacIn/Activity/CallAturNilaiIndex.xml` | berkas Activity terkecil — dipakai memetakan struktur tag `pySteps`, `pxRuleReferences`, `pyParameters` |
| `NB FacIn/Activity/AddCoverageAutoFire.xml` (node langkah 1 saja) | membuktikan adanya **langkah bersarang** (`pySteps` di dalam `pySteps`) |
| `NB FacIn/Activity/InsertCedingProduction.xml` (node langkah 1 saja) | membuktikan letak parameter metode RDB (`pyStepsCallParams/RequestType`, `/Access`) |
| `NB FacIn/RDBList/AttachmentLife.xml` | membuktikan struktur rule SQL (`pyBrowseSQL`, `pyRWAccess`, `pyRequestType`) |

Sisanya dibaca lewat ekstraksi terstruktur atas seluruh 609 (dan 2.083) berkas.

---

## 0. Ringkasan temuan yang mengubah rencana migrasi

| # | Temuan | Label |
| --- | --- | --- |
| T-1 | **Tidak ada satu pun langkah `RDB-Save` di 609 Activity NB.** Seluruh tulisan ke Oracle dijalankan lewat `RDB-List` — metode yang namanya berarti *baca*. | [terverifikasi] |
| T-2 | **55 dari 57 pernyataan SQL yang menulis disimpan di tag `pyBrowseSQL`** (tag "browse"/baca), bukan `pySaveSQL`/`pyDeleteSQL`. Nama tag bukan bukti arah operasi. | [terverifikasi] |
| T-3 | **`COMMIT` berada di dalam SQL, bukan di Pega.** 36 rule SQL memuat `COMMIT;` di badan blok PL/SQL. Hanya 5 langkah `Commit` dan 4 langkah `RollBack` di seluruh 609 Activity. Batas transaksi milik database, bukan aplikasi. | [terverifikasi] |
| T-4 | Struktur langkah **bersarang sampai kedalaman 9**. Menghitung hanya `pySteps` tingkat atas menghasilkan 4.069 langkah; angka sebenarnya **11.090**. Pass yang hanya membaca tingkat atas kehilangan 63 % logika. | [terverifikasi] |
| T-5 | **577 langkah memiliki `pyStepsPreCondition` bernilai `false`**, 404 di antaranya tetap menyimpan ekspresi kondisi. Arti flag ini tidak dijelaskan korpus — dua tafsir yang mungkin menghasilkan perilaku berlawanan. **MEMBLOKIR.** | [pertanyaan terbuka] |
| T-6 | **37 nama activity dirujuk tetapi berkasnya tidak ada** di ekspor NB (di luar 15 activity platform ber-prefiks `px`/`pz`/`py`). Ekspor NB tidak menutup sendiri. | [terverifikasi] |
| T-7 | Hanya **5 activity yang tidak terjangkau** dari titik masuk UI/Flow; **1 kelompok klon** (4 activity byte-identik pada `pySteps`). Kode mati praktis nihil — 609 activity memang dipakai. | [terverifikasi] |
| T-8 | `SearchJobID.xml` memiliki `pyRequestType` = `SearchJobIDSQL`, sedangkan activity memanggil `SearchJobID`. **Nama berkas ≠ kunci rule.** | [terverifikasi] |

---

## 1. Jumlah dan sebaran

### 1.1 Jumlah berkas

609 berkas `.xml` di `NB FacIn/Activity/`; 2.083 berkas di seluruh `NB FacIn/`.

```powershell
# jumlah Activity NB
(Get-ChildItem 'D:\migrasi\RNM\NB FacIn\Activity' -File -Filter *.xml).Count          # 609
# sebaran seluruh folder NB
Get-ChildItem 'D:\migrasi\RNM\NB FacIn' -Directory | ForEach-Object {
  [PSCustomObject]@{ Folder=$_.Name; Files=(Get-ChildItem $_.FullName -Recurse -File -Filter *.xml).Count } }
```

| Folder NB | Berkas |
| --- | --- |
| Activity | **609** |
| Section | 432 |
| FlowAction | 250 |
| RDBList | 217 |
| When | 210 |
| DataTransform | 148 |
| ReportDefinition | 123 |
| Harness | 41 |
| DataPage | 30 |
| DecisionTable | 12 |
| Flow | 6 |
| ConnectREST | 3 |
| SystemSettings | 1 |
| DecisionTree | 1 |
| **Total** | **2.083** |

[terverifikasi]

### 1.2 Sebaran per kelas Pega (`pyClassName`)

56 nilai `pyClassName` berbeda. 15 teratas:

| Jumlah | `pyClassName` |
| --- | --- |
| 246 | `ASM-FW-GISFW-Work` |
| 51 | `ASM-FW-GISFW-Data-Coverage` |
| 32 | `ASM-FW-GISFW-Data-OfferFacIn` |
| 26 | `ASM-FW-GISFW-Data-PolicyTreatyIn` |
| 26 | `@baseclass` |
| 25 | `ASM-FW-GISFW-Data-FacOffer` |
| 22 | `ASM-FW-GISFW-Data-PropertyItem` |
| 19 | `ASM-FW-GISFW-Data-Aneka` |
| 17 | `ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance` |
| 17 | `Data-Portal` |
| 11 | `Code-Pega-List` |
| 10 | `Data-Party-Person` |
| 9 | `ASM-FW-GISFW-Data-Vehicle` |
| 8 | `ASM-FW-GISFW-Data-Clause` |
| 6 | `ASM-FW-GISFW-Data-Occupation` / `Work-` / `ASM-FW-GISFW-Data-Cargo` / `ASM-FW-GISFW-Data-Quotation` |

Sisanya 1–4 berkas per kelas, termasuk 20 kelas integrasi `ASM-FW-GISFW-Int-*`
(`OFFERJSON`, `policyjson`, `T_STORAGE_IMAGE`, `TABLE_B2B`, `SHIP`, `M_LINK_SERVICE`, dll.).

```powershell
$src='D:\migrasi\RNM\NB FacIn\Activity'
Get-ChildItem $src -File -Filter *.xml | ForEach-Object {
  $x=New-Object System.Xml.XmlDocument; $x.Load($_.FullName)
  $x.DocumentElement.SelectSingleNode('pyClassName').InnerText
} | Group-Object | Sort-Object Count -Descending | Select-Object Count,Name
```

[terverifikasi]

### 1.3 RuleSet dan versi

| `pyRuleSet` | Jumlah |
| --- | --- |
| `GISFW` | 595 |
| `Pega-ProcessEngine` | 6 |
| `GISFWInt` | 3 |
| `Pega-API`, `Pega-IntegrationEngine`, `Pega-IntegrationArchitect`, `Pega-ProCom`, `Pega-RulesEngine` | 1 masing-masing |

→ **14 activity berasal dari ruleset platform Pega**, bukan aplikasi. Kandidat *tidak*
dimigrasikan sebagai kode bisnis.

`pyRuleSetVersion`: 60 nilai berbeda. Terbanyak `01-01-52` (105), `01-01-53` (75),
`01-01-91` (43), `01-01-81` (35), `01-01-55` (29), `01-01-95` (27). Rentang aplikatif
`01-01-01` … `01-01-96`; 11 berkas memakai versi platform seri `08-*`.

> Catatan pelajaran pass sebelumnya: 86 rule berbeda versi antar folder siklus. Pass ini
> **tidak membandingkan antar folder**, jadi versi di atas berlaku khusus untuk ekspor NB.

```powershell
Get-ChildItem 'D:\migrasi\RNM\NB FacIn\Activity' -File -Filter *.xml | ForEach-Object {
  $x=New-Object System.Xml.XmlDocument; $x.Load($_.FullName)
  [PSCustomObject]@{ RS=$x.DocumentElement.SelectSingleNode('pyRuleSet').InnerText
                     Ver=$x.DocumentElement.SelectSingleNode('pyRuleSetVersion').InnerText }
} | Group-Object RS | Select-Object Count,Name
```

[terverifikasi]

### 1.4 Jenis dan ketersediaan

| Tag | Nilai | Jumlah |
| --- | --- | --- |
| `pyActivityType` | `ACTIVITY` | 601 |
| | `RULECONNECT` | 7 |
| | `UTILITY` | 1 |
| `pyRuleAvailable` | `Yes` | 598 |
| | `Final` | 11 |
| `pzStatus` | `valid` | 609 |

Tidak ada satu pun activity ber-status `Withdrawn`/`Blocked` — artinya **tidak ada kode mati
yang ditandai secara eksplisit oleh Pega**. [terverifikasi]

### 1.5 Pola penamaan

Sufiks (dihitung **case-sensitive**, `Group-Object -CaseSensitive`):

| Jumlah | Sufiks |
| --- | --- |
| 249 | `_Act` |
| 194 | (tanpa sufiks) |
| 75 | `_ACT` |
| 43 | `_act` |
| 16 | `_PreAct` |
| 13 | `_PostAct` |
| 5 | `_FacIn` |
| 2 | `_EDM`, `_preACT` |
| 1 masing-masing | `_RD` `_Sql` `_DT` `_Reas` `_UW` `_ZIPCODE` `_ANEKA` `_MODEL` `_ActFlow` `_COINS` |

**367 activity memakai sufiks "act" dengan 5 ejaan huruf besar-kecil yang berbeda**
(`_Act`/`_ACT`/`_act`/`_PreAct`/`_PostAct`/`_preACT`). Konvensi tidak pernah ditegakkan.
Untuk Go ini berarti pemetaan nama **tidak boleh** bergantung pada kapitalisasi sufiks.

Kata kerja awal (top 10): `Set` 94 · `Get` 47 · `Count` 46 · `Save` 45 · `Copy` 31 ·
`Check` 21 · `Input` 16 · `Insert` 16 · `Protection` 15 · `Protect` 14.

```powershell
$o='<scratchpad>'   # acts.csv hasil ekstraksi
Import-Csv "$o\acts.csv" | ForEach-Object {
  if($_.Rule -cmatch '(_[A-Za-z]{2,8})$'){$Matches[1]}else{'(tanpa)'}
} | Group-Object -CaseSensitive | Sort-Object Count -Descending
```

**Ketidaksesuaian nama berkas vs `pyRuleName`: 1 kasus** —
`pzChangeStageWrapper.xml` memuat rule `pzChangeStageWrapperV2`. [terverifikasi]

### 1.6 Struktur langkah

| Ukuran | Nilai |
| --- | --- |
| Total langkah (termasuk bersarang) | **11.090** |
| Langkah tingkat atas (`depth 0`) | 4.069 |
| Langkah bersarang | 7.021 (63,3 %) |
| Kedalaman maksimum | 9 |

Sebaran kedalaman: d0 4.069 · d1 2.494 · d2 1.586 · d3 1.085 · d4 791 · d5 546 · d6 323 ·
d7 140 · d8 50 · d9 6.

**Ini jebakan utama.** Langkah bersarang muncul sebagai anak `pySteps` di dalam sebuah
`rowdata` langkah induk yang `pyStepsActivityName`-nya **kosong**; induk tersebut adalah
kontainer perulangan (`pyStepsRepeatDef/pyStepsRepeatDefHasRepeat` = `REPEAT` 484 kali,
`EMBEDDED` 1.814 kali, `PROPERTYLIST` 1, `PAGE` 1). [terverifikasi — dibuktikan pada
`AddCoverageAutoFire.xml` langkah 1: `pyStepsActivityName` kosong,
`pyStepsRepeatDefHasRepeat`=`REPEAT`, berisi 3 anak `pySteps`]

### 1.7 Sebaran metode langkah (11.090 langkah)

| Jumlah | Metode | Jumlah | Metode |
| --- | --- | --- | --- |
| 5.771 | `Property-Set` | 26 | `Obj-Refresh-And-Lock` |
| 1.814 | *(blok `EMBEDDED`)* | 19 | `Property-Map-DecisionTable` |
| 677 | `Call` | 11 | `Exit-Activity` |
| 614 | `RDB-List` | 11 | `Obj-Open-By-Handle` |
| 484 | *(blok `REPEAT`)* | 8 | `Link-Objects` |
| 222 | `Page-Remove` | 6 | `Connect-REST` |
| 187 | `Page-Copy` | 6 | `Branch` |
| 179 | `Property-Remove` | 5 | `History-Add` |
| 167 | `Page-New` | 5 | `Commit` |
| 161 | `Property-Set-Messages` | 4 | `Property-Validate` |
| 150 | `Page-Set-Messages` | 4 | `RollBack` |
| 146 | *(kosong, tanpa blok)* | 3 | `Property-Set-Special` / `Obj-Sort` / `Obj-Delete` / `Show-Page` |
| 97 | `Property-Set-HTML` | 2 | `RDB-Delete` / `Obj-Set-Tickets` / `Call-Automation` / `Activity-Clear-Status` / `Property-Set-Stream` |
| 82 | `Apply-DataTransform` | 1 | `Property-Map-DecisionTree` / `Page-Validate` / `Page-Change-Class` / `Activity-End` / `Obj-Open` / `Log-Message` |
| 49 | `java` | | |
| 48 | `Obj-Browse` | | |
| 45 | `Obj-Save` | | |
| 35 | `Page-Clear-Messages` | | |
| 27 | `Page-Rename` | | |

**`RDB-Save` = 0.** Diverifikasi dua kali, dengan parser XML dan dengan pencarian teks:

```powershell
$src='D:\migrasi\RNM\NB FacIn\Activity'
foreach($m in 'RDB-Save','RDB-List','RDB-Delete','Obj-Save','Commit','RollBack'){
  $n=(Select-String -Path "$src\*.xml" -SimpleMatch "<pyStepsActivityName>$m<" -AllMatches | Measure-Object).Count
  "{0,-12} {1}" -f $m,$n }
# RDB-Save 0 | RDB-List 614 | RDB-Delete 2 | Obj-Save 45 | Commit 5 | RollBack 4
```

Catatan ejaan: korpus memakai **dua ejaan** untuk metode rollback — `RollBack`
(`ASMForceCaseClose.xml`) dan `Rollback` (`Save.xml`, `DeleteAttachment.xml`,
`svcAddWorkObject.xml`). Registry metode di Go harus *case-insensitive*. [terverifikasi]

---

## 2. Graf panggilan

### 2.1 Dua sumber sisi, disatukan

| Sumber | Sisi | Keandalan |
| --- | --- | --- |
| Langkah `Call …` / `Branch …` (`pyStepsActivityName`) | **683** | eksplisit di badan activity |
| Indeks `pxRuleReferences` (`pxRuleObjClass` = `Rule-Obj-Activity`, tidak self-ref) | **563** | dihasilkan mesin Pega, sudah ter-dedup per berkas |

Format target pada langkah: `Call [<Kelas>.]<NamaActivity>` — kelas boleh dihilangkan
(`Call CountPremi_ACT`) atau ditulis penuh (`Call  ASM-FW-GISFW-Work.SetRIComm_Act`;
perhatikan **spasi ganda**). Satu kasus memakai huruf kecil (`call  ASM-FW-…`).

Referensi activity di seluruh 2.083 berkas NB (`pxRuleReferences` → `Rule-Obj-Activity`,
total 690):

| Folder pemanggil | Ref |
| --- | --- |
| Activity | 567 |
| FlowAction | 63 |
| Flow | 35 |
| DataPage | 18 |
| Harness | 4 |
| Section | 3 |

Karena indeks `pxRuleReferences` pada Section jelas kurang lengkap (hanya 3, padahal
Section adalah folder terbesar), graf luar dilengkapi **pemindaian node teks** atas 1.474
berkas non-Activity, dibatasi pada tag pemanggil sah:
`pyActivity`, `pyPreProcessingActivity`, `pyPostProcessingActivity`, `pyLocalActionActivity`,
`pyLoadActivity`, `pyDeferLoadActivity`, `pyDeferLoadRetrievalActivity`, `pyActivityName`,
`pyValidateActivity`, `pyImplementation`, `pyDataSourceHolder`.

Hasil tag kuat:

| Jumlah | Folder, tag |
| --- | --- |
| 648 | Section, `pyActivity` |
| 99 | Harness, `pyActivity` |
| 36 | FlowAction, `pyPreProcessingActivity` |
| 25 | FlowAction, `pyLocalActionActivity` |
| 25 | Flow, `pyImplementation` (shape *Utility*) |
| 14 | DataPage, `pyLoadActivity` |
| 14 | DataPage, `pyDataSourceHolder` |
| 15 | Section/Harness, `pyDeferLoad*` |

> **Jebakan yang dihindari:** pemindaian teks mentah memberi 425 hit palsu karena ada
> activity bernama `Save` yang cocok dengan nilai label UI (`pyToolbarSaveLabel`,
> `pyButtonCaption`, `pyMemo`). Tag-tag label itu dibuang.

```powershell
# rekonstruksi sisi dari langkah Call/Branch
Import-Csv "$o\steps2.csv" | Where-Object { $_.Method -match '^(?i)(call|branch)\s+\S' } | Measure-Object
# 683
```

[terverifikasi]

### 2.2 Klasifikasi 609 activity menurut siapa yang merujuknya

| Kelas | Jumlah | Arti |
| --- | --- | --- |
| **A** — hanya dipanggil activity lain | 212 | helper internal murni → fungsi tak-diekspor di Go |
| **B** — hanya dirujuk UI/Flow/DataPage | 330 | **titik masuk** → handler / service method publik |
| **C** — keduanya | 64 | layanan bersama yang juga dipanggil UI |
| **D** — tidak dirujuk sama sekali | **3** | lihat §4 |

Titik masuk (B + C = 394) memang dirujuk dari `Flow\`/`FlowAction\`/`Harness\`/`Section\`/
`DataPage\`, sesuai dugaan awal tugas. **Pembedaan "titik masuk vs kode mati" karena itu
bisa dijawab tegas: kode mati hampir tidak ada.**

```powershell
$sum=Import-Csv "$o\summary.csv"
'A={0} B={1} C={2} D={3}' -f `
  ($sum|?{[int]$_.FanInAct -gt 0 -and [int]$_.FanInExt -eq 0}).Count,
  ($sum|?{[int]$_.FanInAct -eq 0 -and [int]$_.FanInExt -gt 0}).Count,
  ($sum|?{[int]$_.FanInAct -gt 0 -and [int]$_.FanInExt -gt 0}).Count,
  ($sum|?{[int]$_.FanInAct -eq 0 -and [int]$_.FanInExt -eq 0}).Count
```

[terverifikasi]

### 2.3 Titik masuk proses (shape *Utility* di 6 berkas `Flow\`)

Ini adalah daftar terpenting untuk lapisan `handlers` Go: activity yang dijalankan
**oleh mesin alur**, bukan oleh layar.

| Flow | Activity |
| --- | --- |
| `InputInwardFacultativeOffer.xml` | `GetLimitAkseptasi_ActFlow`, `GetLimitAkseptasi_JUW_UW`, `GetLimitAkseptasiLife_Act`, `SaveJsonOfferFacIn_Act`, `SaveToProduction_ACT`, `SendEmailPolicy`, `SetBanding_ACT`, `SetNBStatus_Act`, `SetTicket`, `SetToJsonOffer_ACT` |
| `InputInwardFacultativeRISlip.xml` | `CekLimitSpreading_Act`, `GetLimitAkseptasi_Act`, `GetLimitAkseptasi_JUW_UW`, `SaveJsonPolicyFacIn_Act`, `SendEmailPolicy`, `serviceInsertArasapas_act`, `SetNBStatus_Act`, `SetToInbox_ACT`, `UpdateStsKonversiFacOut_Act` |
| `InputQuotation.xml` | `SetBusinessType_Act` |
| `InputRealizationTreatyIn.xml` | `SaveJsonPolisTreatyIn_Act`, `serviceInsertArasapas_act` |
| `OfferFacRetro.xml` | `InsertFacoutProd`, `SendEmailPolicy`, `SetNBStatus_Act` |
| `OfferFacOut.xml` | *(tidak ada shape Utility yang merujuk activity)* |

25 rujukan, 19 activity berbeda. [terverifikasi]

### 2.4 Fan-in tertinggi — kandidat layanan bersama di Go

Fan-in internal = jumlah **activity berbeda** yang memanggilnya.

| Activity | Kelas | fan-in activity | fan-in UI | Langkah | Kandidat paket |
| --- | --- | --- | --- | --- | --- |
| `GetCurrencyMaster` | `Data-PropertyItem` | 19 | 4 | 3 | `repository/lookup` |
| `cekSpreadingFactIn` | `Data-Coverage` | 16 | 6 | 10 | `spreading` |
| `CountPremi_ACT` | `Data-Coverage` | 14 | 8 | 57 | `premium` |
| `SetRIComm_Act` | `Work` | 11 | 10 | 10 | `premium` |
| `CountNetPremi_act` | `Data-PolicyTreatyIn` | 10 | 2 | 7 | `premium` |
| `GetTreatyName` | `Work` | 9 | 0 | 20 | `spreading` |
| `CountPremiAndTSINusantaraReEDM_ACT` | `Work` | 8 | 1 | 14 | `premium` |
| `SetCategoryAttach` | `Data-OfferFacIn` | 8 | 5 | 4 | `integration` |
| `CountTotalTSIPremiNusaRe_Act` | `Work` | 7 | 0 | 15 | `premium` |
| `AddCurrencyList__ACT` | `Data-Party-Person` | 7 | 1 | 5 | `faccase` |
| `SumTotalTSIPremiGross_Act` | `Data-Coverage` | 6 | 1 | 6 | `premium` |
| `SendEmailWithAttachments` | `@baseclass` | 6 | 0 | 1 | `integration` |
| `GetLinkService` | `Int-M_LINK_SERVICE` | 6 | 0 | 4 | `config` / `integration` |
| `CountPremiAndTSINusantaraRe_ACT` | `Work` | 6 | 10 | 12 | `premium` |
| `GetHistoryAkseptasiPega_Act` | `Work` | 5 | 0 | 1 | `acceptance` |
| `CountPremiCoverageAneka` | `Data-Coverage` | 5 | 1 | 18 | `premium` |
| `DeleteAttachment` | `Work-` | 5 | 0 | 30 | `integration` |
| `InsertHistoryAkseptasiPega` | `Work` | 5 | 1 | 9 | `acceptance` |
| `SumTreatyCapacity_Act` | `Data-OfferFacIn` | 4 | 0 | 11 | `spreading` |
| `InsertDocument_Act` | `Int-DOCUMENT_POLIS` | 4 | 0 | 6 | `integration` |

Fan-in **eksternal** tertinggi (paling sering ditanam di layar) — memasok prioritas
endpoint API: `downloadFile` 17 · `IsThereAnyObjectLocation_Act` 15 ·
`CountPremiAndTSINusantaraRe_ACT` 10 · `ViewOldDataEDM_Act` 10 · `SetRIComm_Act` 10 ·
`CalcultePersentageSpeading_Act` 9 · `ChangeSpreadingPercentage_Act` 9 ·
`CopyToAllSpreading_ACT` 9 · `cekSpreadingFactIn_Act` 9 · `CountPremi_ACT` 8.

> **Perhatikan:** `cekSpreadingFactIn` dan `cekSpreadingFactIn_Act` adalah **dua rule
> berbeda**. Jangan digabung.

[terverifikasi]

### 2.5 Activity dirujuk tetapi berkasnya TIDAK ADA (ekspor tidak lengkap)

52 nama target berbeda tidak punya berkas di `NB FacIn/Activity/`.

**15 berprefiks `px`/`pz`/`py`** — activity platform Pega, wajar tidak diekspor:
`pxShowReport` (48 pemanggil-langkah), `pxRetrieveReportData` (23), `pxUploadCSVResults` (7),
`pzCheckFieldSecurity`, `pzRemoveAttachments`, `pxConvertResultsToCSV`,
`pzMapAutomationErrorsToResponse`, `pyDeleteAttachmentContent`,
`pxValidateAndDeleteThumbnail`, `pzApplyPageInstructions`, `pzStoreAttachmentsSentOnCorr`,
`pzUpdateAndDeleteAssignments`, `pzGetDescendants`, `pzCheckWorkObjectID`,
`pzAttachUploadedFileWrapper`.

**37 nama lain** — sebagian tampak activity standar Pega tanpa prefiks
(`addWork`, `createWorkPage`, `WorkCommit`, `WorkUnlock`, `DisplayHarness`, `Show-Harness`,
`commitWithErrorHandling`, `HandleCommitFailure`, `UpdateWorkObject`, `SetPageErrors`,
`SetOutput`, `setProcessState`, `performAssignmentCheck`, `CorrSend`, `CorrAttach`,
`GetEmailSenderInfo`, `SendEmailNotification`, `MSOGenerateExcelFile`) [dugaan — disimpulkan
dari pola nama Pega OOTB, **belum terverifikasi** karena berkasnya tidak ada],
sebagian jelas milik aplikasi:

| Target hilang | Dipanggil oleh |
| --- | --- |
| `AttachAsPDFCFO` | `GeneratePDFFacOutMemoPlacing`, `GeneratePDFSlipKomisi` |
| `AveragePremium_Act` | `FillActInstallment` |
| `CheckCeding_Act` | `DeleteCeding_Act` |
| `CheckLimitAdditionalTreatyType_Act` | `SpreadingAdditionalProtection` |
| `CountASMGrossPremi_ACT` | `SetErrorMessage_Act` |
| `CountASMNetPremi_ACT` | `SetErrorMessage_Act` |
| `CountTreatyCapacity_Act` | `SumTreatyCapacity_Act` |
| `EDMPayment` | `InputDtlPayment_PreAct` |
| `fillActAdminFee` | `fillActBussiness` |
| `GetLimitAkseptasi_Act2` | `InputOfferFacInEngineerUW_preACT`, `SetValidateDate_PostAct` |
| `InsertJsonH2HNusare_act` | `ChooseActivityMappingCSV_Act` |
| `InserttoRouteNotif` | `SetNBStatus_Act` |
| `InsertUploadCredit_Act` | `ChooseActivityMappingCSV_Act` |
| `MoreThan100_Act_FacIn` | `CountPremiFacIn_act` |
| `SaveFacinProdEDMBonding_Act` | `SaveFacinProdAllEDM_Act` |
| `SearchTemplateAdditionalCoverageSQL_PreAct` | `UploadCSVVehicle_PostAct` |
| `SetCurrencyServis_PreAct` | `InputOfferFacInEngineer_preACT` |
| `SetTsiCovAct` | `UploadCSVVehicle_PostAct` |
| `SumTSIPremiSpreadRNMMultiCob_Act` | `SumTSIPremiSpreadedRNM_Act` (kelas `ASM-FW-GISFW-Data-Coverage`) |

**Konsekuensi:** `GetLimitAkseptasi_Act2` adalah lubang di **mesin tangga persetujuan** —
dua activity NB memanggilnya dan badannya tidak ada di ekspor ini. Demikian pula
`CountASMGrossPremi_ACT` / `CountASMNetPremi_ACT` adalah lubang di **rumus premi**.
Lihat §9 (memblokir).

```powershell
Import-Csv "$o\missing.csv" | Measure-Object                                   # 52
(Import-Csv "$o\missing.csv" | ?{ $_.Name -cnotmatch '^(px|pz|py)' }).Count     # 37
```

[terverifikasi]

---

## 3. Siapa menulis ke Oracle

### 3.1 Bentuk akses: semuanya lewat `RDB-List`

| Ukuran | Nilai |
| --- | --- |
| Langkah `RDB-List` + `RDB-Delete` | 616 |
| Activity yang memuat langkah RDB | **192** dari 609 |
| `pyRequestType` berbeda yang dirujuk | 248 |
| … yang rule-nya ada di `NB FacIn/RDBList/` | 216 |
| … yang rule-nya **tidak ada** | **32** |
| Berkas rule di `NB FacIn/RDBList/` | 217 (216 nama `pyRequestType` unik) |
| Pernyataan SQL (satu berkas bisa punya >1 tag) | 221 |

Parameter langkah RDB ada di **anak langsung `pyStepsCallParams`**, bukan di `pyParamArray`:

```xml
<pyStepsCallParams>
  <Access>RNM</Access>
  <ClassName>ASM-FW-GISFW-Work</ClassName>
  <RequestType>InsertCedingProduction_SQL</RequestType>
  <RunInParallel>false</RunInParallel>
  …
</pyStepsCallParams>
```
*(Asal: `NB FacIn/Activity/InsertCedingProduction.xml`, langkah 1, metode `RDB-Delete`.)*
[terverifikasi]

**Dua alias koneksi database** dipakai: `ASM` (529 langkah) dan `RNM` (87 langkah).
Kunci rule SQL = `pyClassName` + `pyRWAccess` + `pyRequestType`. Lapisan repository Go
karena itu butuh **dua pool koneksi terpisah**, atau minimal pemetaan alias→DSN yang eksplisit.
Nilai DSN tidak ada di korpus. [terverifikasi untuk alias; **belum terverifikasi** untuk DSN]

### 3.2 Klasifikasi 221 pernyataan SQL

| Verb awal | Jumlah | Tag penyimpan |
| --- | --- | --- |
| `SELECT` | 163 | `pyBrowseSQL` |
| `BEGIN …` (blok PL/SQL) | 31 | `pyBrowseSQL` |
| `DECLARE …` | 14 | 12 `pyBrowseSQL`, 2 `pySaveSQL` |
| `UPDATE` | 7 | `pyBrowseSQL` |
| `DELETE` | 3 | 1 `pyBrowseSQL`, 2 `pyDeleteSQL` |
| `INSERT` | 2 | `pyBrowseSQL` |
| `WITH` | 1 | `pyBrowseSQL` |

**57 pernyataan menulis. 55 di antaranya berada di `pyBrowseSQL`.** Hanya 2 rule memakai
`pySaveSQL` (`ConvertBusinessField`, `ConvertNationality`) dan 2 memakai `pyDeleteSQL`
(`InsertCedingProduction_SQL`, `InsertOfferProduction_Sql`).

Dua rule bernama `Insert…` justru menyimpan **DELETE** di `pyDeleteSQL` dan blok PL/SQL
di `pyBrowseSQL` — sesuai pelajaran "nama bukan bukti" yang sudah tercatat di pass
sebelumnya, dan **terkonfirmasi ulang** di sini. [terverifikasi]

```powershell
foreach($f in (Get-ChildItem 'D:\migrasi\RNM\NB FacIn\RDBList' -File -Filter *.xml)){
  $x=New-Object System.Xml.XmlDocument; $x.Load($f.FullName)
  foreach($t in 'pyBrowseSQL','pySaveSQL','pyDeleteSQL'){
    $n=$x.DocumentElement.SelectSingleNode($t); if($n -and $n.InnerText.Trim()){
      $s=($n.InnerText -replace '\s+',' ').Trim()
      [PSCustomObject]@{Rule=$x.DocumentElement.SelectSingleNode('pyRequestType').InnerText
                        Tag=$t; Verb=($s -split '\s+')[0].ToUpper()} } }
} | Group-Object Verb | Sort-Object Count -Descending
```

### 3.3 `COMMIT` ada di dalam SQL

36 dari 221 pernyataan memuat `COMMIT;`. Contoh (dikutip apa adanya):

```sql
-- Asal: NB FacIn/RDBList/InsertHistoryAkseptasiPega_Sql.xml  (pyBrowseSQL, Access=ASM)
BEGIN INSERT INTO HISTORYAKSEPTASIPEGA
  (ID_PEGA, Tgl_Transfer, Status, Username, Workbasket, ID_KOMITE)
  VALUES ({InsertHistory.CARI1}, sysdate, {InsertHistory.CARI5},
          {InsertHistory.CARI4}, {InsertHistory.CARI2}, {InsertHistory.CARI6});
  COMMIT; END;
```

```sql
-- Asal: NB FacIn/RDBList/GetSequenceNumber_SQL.xml  (pyBrowseSQL, Access=RNM)
BEGIN POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER({ParamSeq.CARI1},{ParamSeq.CARI2},
  TO_DATE({ParamSeq.CARI3},'DD/MM/YYYY'), {ParamSeq.HASIL1 out},{ParamSeq.HASIL2 out});
  COMMIT; END;
```

**21 pernyataan penulis TIDAK memuat `COMMIT`** — di antaranya
`InsertTreatyProduction_Sql`, `InsertFacinProductionBackup_Sql`, `InsertOfferProduction_Sql`,
`ForInputCurrencyAdj_Sql`, `UpdateJsonOffer_SQL`, `UpdatePolisEndorsement_SQL`,
`UpdateSTSKonversi_SQL`, `UpdateFacProdObjItemID_SQL`, `DeleteDataProduction`,
`DeleteStorage_SQL`, `Update_T_Storage_SQL`. Bagaimana pernyataan-pernyataan ini
di-commit **tidak dapat dijawab dari korpus** — hanya 5 langkah `Commit` yang ada di
seluruh 609 Activity, dan tidak satu pun berada di activity penulis besar seperti
`SaveTreatyProduction_Act` atau `SaveOfferProduction_Act`. Lihat §9 (memblokir).
[terverifikasi untuk jumlah · [pertanyaan terbuka] untuk mekanismenya]

### 3.4 62 activity yang menulis ke Oracle

Pemetaan **activity → rule SQL penulis**. Ini memasok batas transaksi untuk
`internal/repository/oracle/`.

| Activity (berkas) | Rule SQL penulis |
| --- | --- |
| `ConvertPropViewPolisFacOut_Act.xml` | `ConvertBusinessField`, `ConvertNationality` |
| `createNewAccumulation_act.xml` | `SaveNewAccumulation_SQL` |
| `DeleteGoogleStorage_Act.xml` | `DeleteStorage_SQL`, `GetTokenStorage_SQL` |
| `GeminiAIGoogle_Act.xml` | `GetTokenStorage_SQL` |
| `GenerateNopolis_Act.xml` | `GetSequenceNumber_SQL` |
| `GeneratePolicyNoTreaty_Act.xml` | `GetSequenceNumber_SQL` |
| `GetInsuredID.xml` | `INSERTCONORGJSON_MCLIENT` |
| `GetUrlGoogleStorage_Act.xml` | `GetTokenStorage_SQL`, `Update_T_Storage_SQL` |
| `InsertCedingProduction.xml` | `InsertCedingProduction_SQL` |
| `InsertFacInPerformance_Act.xml` | `InsertDataFacinPerformance_SQL` |
| `InsertFacoutProduction.xml` / `…EDM.xml` / `…EDMCurr.xml` | `InsertTreatyProd_Sql` |
| `InsertGoogleStorage_Act.xml` | `GetTokenStorage_SQL`, `Insert_T_Storage_SQL` |
| `InsertHistoryAkseptasiPega.xml` | `InsertHistoryAkseptasiPega_Sql` |
| `InsertLogServiceProd.xml` | `InsertLogServiceProd` |
| `InsertUploadFire_act.xml` | `UpdateMasterAccumulation`, `UpdateMasterRiskAddress_SQL` |
| `ProtectUploadMarineCargo_act.xml` | `UpdateMasterShip` |
| `SaveAccumulatedType_Act.xml` | `UpdateMasterAccumulatedType` |
| `SaveAccumulation_Act.xml` / `SaveAccumulationByRiskAddress_Act.xml` / `SearchAccumulationData_Act.xml` | `UpdateMasterAccumulation` |
| `SaveBranch_Act.xml` | `UpdateMasterBranch` |
| `SaveCity_Act.xml` | `UpdateMasterCity` |
| `SaveDataToJsonFollowing_Act.xml` | `InsertJsonFollowing_SQL` |
| `SaveDistrict_Act.xml` | `UpdateMasterDistrict` |
| `SaveEDMToJsonPolicy_Act.xml` | `GetSequenceNumber_SQL`, `INSERTJSON_JSONPOLISEDM_FACIN` |
| `SaveFacinLive_Act.xml` | `InsertFacinLife_Sql`, `InsertFacinLifeMonthly_Sql` |
| `SaveFacinOfferLife_Act.xml` | `InsertFacinOfferLife_Sql` |
| `SaveFacinProdCurr{Aneka,Fire,Golf,MarineCargo,MBU,PA}_Act.xml` (6) | `ForInputCurrencyAdj_Sql`, `ForInputCurrencyAdjBackup_Sql`, `INSERTERRORFACIN_Sql`, `MachingDataFacin_Sql` |
| `SaveFacinProdEDM{Fire,Golf,MarineCargo,MBU,PA}_Act.xml` (5) | `INSERTERRORFACIN_Sql`, `InsertFacinProductionBackup_Sql`, `InsertTreatyProduction_Sql`, `MachingDataFacin_Sql` |
| `SaveFacinProdFireNB_Act.xml` | `InsertTreatyProduction_Sql` |
| `SaveFacinRNWProd_Act.xml` | idem 4 rule di atas |
| `SaveFacinSpreadLife_Sql.xml` | `InsertIntoFacinSPreadLife_Sql`, `InsertIntoFacinSPreadLifeMonthly_Sql` |
| `SaveJsonOfferFacIn_Act.xml` | `SaveOfferJson_SQL` |
| `SaveJsonPolicyFacIn_Act.xml` | `GetSequenceNumber_SQL`, `InsertIntoJsonError_SQL`, `INSERTJSON_JSONPOLIS_FACIN`, `InsertRiskAndLossProfile_SQL`, `UpdateFacProdObjItemID_SQL` |
| `SaveJsonPolisTreatyIn_Act.xml` | `SavePolisTreatyIn_SQL` |
| `SaveMarketingOfficer_Act.xml` | `UpdateMasterMarketingOfficer` |
| `SaveNation_Act.xml` | `UpdateMasterNation` |
| `SaveOfferProduction_Act.xml` | `InsertOfferProduction_Sql` |
| `SaveProvince_Act.xml` | `UpdateMasterProvince` |
| `SaveRiskAddress_Act.xml` | `UpdateMasterRiskAddress_SQL` |
| `SaveRW_Act.xml` | `UpdateMasterRW` |
| `SaveTreatyProduction_Act.xml` / `SaveTreatyProductionEDMAneka_Act.xml` | `INSERTERRORFACIN_Sql`, `InsertFacinProductionBackup_Sql`, `InsertTreatyProduction_Sql`, `MachingDataFacin_Sql` |
| `SaveViewSuggest.xml` | `InsertViewSuggest_SQL` |
| `serviceInsertArasapasEDM_act.xml` | `DeleteDataProduction`, `INSERTJSON_JSONPOLISMONITORING_FACIN`, `UpdateErrorNoteJsonPolisMonitoring` |
| `SetTreatyIn_Act.xml` | `SaveTreatyIn` |
| `UpdateAccumulationLife_Act.xml` | `Update_sql` |
| `UpdateJSONOffer_ACT.xml` | `UpdateJsonOffer_SQL` |
| `UpdatePolisAddendum_Act.xml` | `UpdatePolisEndorsement_SQL` |
| `UpdateStsKonversiFacOut_Act.xml` | `UpdateSTSKonversi_SQL` |
| `UpdateTglBindJsonOffer_Act.xml` | `UpdateTgl_BindJsonOffer_SQL` |

Total **62 activity**, **225 langkah RDB** yang menunjuk rule SQL penulis.

> Catatan kehati-hatian: `GetSequenceNumber_SQL` dan `GetTokenStorage_SQL` bernama `Get…`
> tetapi memanggil prosedur dengan parameter `OUT` **dan** `COMMIT`. Keduanya dihitung
> sebagai penulis. Apakah prosedurnya benar-benar mengubah baris **belum terverifikasi** —
> isi stored procedure tidak ada di korpus.

```powershell
$q=Import-Csv "$o\sql.csv"; $c=Import-Csv "$o\stepcp.csv"
$w=[Collections.Generic.HashSet[string]]::new([StringComparer]::OrdinalIgnoreCase)
$q|?{$_.Verb -notin 'SELECT','WITH'}|%{[void]$w.Add($_.RuleName)}
$wr=$c|?{$_.Method -match '^(RDB-List|RDB-Delete)' -and $w.Contains($_.RequestType)}
$wr.Count; ($wr|Select-Object -Expand File -Unique).Count     # 225 ; 62
```

[terverifikasi]

### 3.5 Objek Oracle yang disentuh

**Prosedur/paket yang dipanggil** (nama diikuti `(` di dalam blok PL/SQL) — 35 objek
`POOLDATA.*` plus `GENERAL.F_GET_NM_ASURADUR`, `MBU.F_CEK_HURUF`, `DBMS_LOB.CREATETEMPORARY`
(10 rule), `UTL_MATCH.EDIT_DISTANCE_SIMILARITY`, `STANDARD_HASH`, `SYS_GUID`, `JSON_TABLE`.

Prosedur tulis yang penting:

| Prosedur | Rule SQL pemanggil |
| --- | --- |
| `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` | `GetSequenceNumber_SQL` |
| `POOLDATA.INSERTJSONPOLIS` | `INSERTJSON_JSONPOLIS_FACIN`, `INSERTJSON_JSONPOLISEDM_FACIN` |
| `POOLDATA.INSERTJSONPOLISMONITORING` | `INSERTJSON_JSONPOLISMONITORING_FACIN` |
| `POOLDATA.INSERTUPDATECEDINGPRODUCTION` | `InsertCedingProduction_SQL` |
| `POOLDATA.INSERTUPDATERISKADDRESS` | `UpdateMasterRiskAddress_SQL` |
| `POOLDATA.FACINFORBACKUP` | `MachingDataFacin_Sql` |
| `POOLDATA.PEGA_DELETE_ERROR_KONVERSI` | `DeleteDataProduction` |
| `POOLDATA.PEGA_JSON_POLIS_TREATYIN` | `SavePolisTreatyIn_SQL` |
| `POOLDATA.PEGA_M_ACCUMULATION_LIFE` | `SaveNewAccumulation_SQL`, `Update_sql` |
| `POOLDATA.PEGA_M_JSON_OFFER` | `SaveOfferJson_SQL` |
| `POOLDATA.PEGA_MARKETINGOFFICER` | `UpdateMasterMarketingOfficer` |
| `POOLDATA.PEGA_TREATY_IN` | `SaveTreatyIn` |
| `POOLDATA.PROSESCOPY` | `ConvertBusinessField`, `ConvertNationality` |
| `POOLDATA.RDBINSERTCLIENT` | `INSERTCONORGJSON_MCLIENT` |
| `POOLDATA.RDBMASTER{ACCUMULATEDTYPE,ACCUMULATION,BRANCH,CITY,DISTRICT,NATION,PROVINCE,RW,SHIP}` | 9 rule `UpdateMaster*` |
| `POOLDATA.GET_TOKEN_STORAGE` | `GetTokenStorage_SQL` |
| `POOLDATA.GENERATE_FACRETRO_NO` | `GenerateRISlipNumber` |
| `POOLDATA.GETCURRENCYSTANDARD` | `CurrencyStandard` |

**Isi seluruh stored procedure di atas belum terverifikasi** — tidak ada satu pun badan
prosedur di korpus. Ini area terbesar yang tidak dapat direkam dari ekspor Pega.

**Tabel/view** pada klausa `FROM/JOIN/INTO/UPDATE`: **117 objek berbeda**. Terbanyak
`JSON_POLIS` (18 rule), `DUAL` (12), `FACINPRODUCTION` (7), `ACCUMULATION` (6), `RW` (6),
`CURRENCY` (5), `OCCUPATION` (5), `TREATYBUSINESS` (5), `POOLDATA.JSON_POLIS` (5),
`POOLDATA.REINSURANCETYPE` (5), `POOLDATA.M_LIMIT_PROPERTYY` (4), `MARKETINGOFFICER` (4),
`PROPORTIONALARRG` (4), `FACOUTPRODUCTION` (3), `HISTORYAKSEPTASIPEGA` (3),
`JSON_FOLLOWING` (3), `M_CLAUSE` (3), `T_STORAGE_IMAGE` (3).

Dua tabel internal Pega tersentuh langsung dari SQL aplikasi:
`DATAPEGA.PC_ASM_FW_GCNMFW_WORK` (1 rule) dan `DATAPEGA.PC_HISTORY_ASM_FW_GISFW_WORK`
(2 rule). Sesuai CLAUDE.md §4.3, tabel `DATAPEGA.PC_*` **hilang bersama Pega** — dua/tiga
query ini harus dirancang ulang, bukan diporting. [terverifikasi]

```powershell
$q=Import-Csv "$o\sql.csv"; $t=@{}
foreach($r in $q){ foreach($m in [regex]::Matches($r.Sql,'(?i)\b(from|join|into|update)\s+([A-Za-z_][\w$#]*(?:\.[\w$#]+)?)')){
  $k=$m.Groups[2].Value.ToUpper(); if(-not $t.ContainsKey($k)){$t[$k]=@()}; $t[$k]+=$r.RuleName } }
$t.Count            # 117
```

### 3.6 Penulis jalur objek Pega (`Obj-Save` / `Commit` / `Obj-Delete`)

Ini jalur kedua, terpisah dari RDB: menyimpan *work object* Pega, bukan tabel bisnis.

| Kategori | Jumlah langkah | Activity |
| --- | --- | --- |
| `Obj-Save` | 45 | 41 |
| `Obj-Refresh-And-Lock` | 26 | 25 |
| `Commit` | 5 | `ASMForceCaseClose`, `AttachFileGIS`, `AttachRISlipToWork`, `AttachToWork`, `DeleteAttachment` |
| `RollBack`/`Rollback` | 4 | `ASMForceCaseClose`, `DeleteAttachment`, `Save`, `svcAddWorkObject` |
| `Obj-Delete` | 3 | `DeleteAttachment` (2), `DeleteDocumentPolis_Act` (1) |
| `Link-Objects` | 8 | `AttachFileGIS`, `AttachRISlipToWork`, `AttachToWork` |
| `History-Add` | 5 | `AttachRISlipToWork`, `AttachToWork`, `DeleteAttachment` |

**Seluruh 5 `Commit` berada di activity lampiran/penutupan kasus.** Tidak ada satu pun
`Commit` di jalur produksi/premi. [terverifikasi]

---

## 4. Kode mati

Dibedakan tegas antara **terbukti mati** dan **tidak ditemukan pemanggilnya**.

### 4.1 Tidak dirujuk sama sekali — 3 activity

| Activity | Kelas | Byte | Langkah (bersarang) |
| --- | --- | --- | --- |
| `InputDtlAdditional_PreAct` | `ASM-FW-GISFW-Data-Coverage` | 67.935 | 8 |
| `ValidateCoverage` | `ASM-FW-GISFW-Data-Coverage` | 119.863 | 18 |
| `pzChangeStageWrapperV2` | `Pega-API-CaseManagement-Case` | 92.559 | 5 (ruleset platform) |

Status: **bukan "terbukti mati"** — hanya "tidak ada pemanggil di dalam korpus NB".
Ketiganya bisa dipanggil dari ekspor siklus lain, dari ruleset yang tidak diekspor,
atau dari pemanggilan dinamis (nama activity dirakit saat runtime).
`pzChangeStageWrapperV2` adalah rule platform Pega — memang tidak diharapkan dirujuk
aplikasi. [terverifikasi untuk fakta "tidak ada rujukan" · [dugaan] untuk penyebabnya]

### 4.2 Tidak terjangkau secara transitif dari titik masuk — 5 activity

Menelusuri maju dari 394 titik masuk (activity yang dirujuk `Flow`/`FlowAction`/`Harness`/
`Section`/`DataPage`) melalui semua sisi `Call`/`Branch`:

| Activity | Kelas | Byte | Langkah (bersarang) | Pemanggil internal |
| --- | --- | --- | --- | --- |
| `ValidateCoverage` | `Data-Coverage` | 119.863 | 18 | — |
| `CopyTemplateCov_PostAct` | `Data-Vehicle` | 300.938 | 36 | `ValidateCoverage` |
| `SpreadingAdditionalProtection` | `Data-Coverage` | 99.754 | 13 | `ValidateCoverage` |
| `InputDtlAdditional_PreAct` | `Data-Coverage` | 67.935 | 8 | — |
| `pzChangeStageWrapperV2` | platform Pega | 92.559 | 5 | — |

Inilah kasus "**activity yang hanya dipanggil dari activity mati**" yang diminta:
`CopyTemplateCov_PostAct` dan `SpreadingAdditionalProtection` hidup semata-mata karena
`ValidateCoverage`, yang sendirinya tak berpemanggil.

Isi kelimanya konsisten satu tema: **coverage tambahan pada objek kendaraan**
(`pyWorkPage.VehicleList(…).CoverageList`, `.AdditionalCoverage`, `TempOldCovAddId`).
`SpreadingAdditionalProtection` bahkan memanggil `CheckLimitAdditionalTreatyType_Act`
yang **berkasnya juga tidak ada** (§2.5). Pola ini konsisten dengan satu fitur yang
kabelnya ke UI tidak ikut terekspor, atau fitur yang belum selesai.

**Bukan "terbukti mati".** Untuk membuktikan mati diperlukan pencarian nama activity di
ekspor siklus lain dan di log runtime — keduanya di luar batas pass ini. [terverifikasi
untuk ketakterjangkauan · [pertanyaan terbuka] untuk status hidup/mati]

```powershell
# skrip reach.ps1 di scratchpad: BFS dari 394 root
# roots(entry)=394  activity TAK TERJANGKAU=5
```

### 4.3 Langkah dengan pra-kondisi dimatikan — 577 langkah

| `pyStepsPreCondition` | `pyStepsPreCondParamsWhen` | Jumlah langkah |
| --- | --- | --- |
| *(tidak ada tag isi)* | kosong | 5.763 |
| `true` | ada | 4.744 |
| `false` | **ada** | **404** |
| `false` | kosong | 173 |
| `true` | kosong | 2 |
| `0` | kosong | 4 |

404 langkah menyimpan ekspresi kondisi yang lengkap sambil flag-nya `false`. Contoh
(dikutip apa adanya):

| Berkas | Langkah | Metode | Kondisi tersimpan |
| --- | --- | --- | --- |
| `AddCoverageAutoFire.xml` | 3 | `Property-Set` | `@LengthOfPageList(.CoverageList)>0` |
| `AddCurencyList_ACT.xml` | 9.2 | `Call …SetRIComm_Act` | `IsEdmAdjTSI` |
| `AttachAsPDFCFOLife.xml` | 3 | `Call AttachRISlipToWork` | `Param.IsViewPolis !="1"` |
| `CalculatePremiPA_FacIn.xml` | 1, 7 | `Property-Set` | `.CalculateMethod_FacIn==1` |
| `ASMForceCaseClose.xml` | 6 | `Obj-Open-By-Handle` | `local.isPrimaryPageUsed` |

Berkas dengan pra-kondisi dimatikan terbanyak:
`UploadCSVPolicyMemberEDM_PostAct` 37 · `SaveOfferProduction_Act` 19 ·
`SaveFacinRNWProd_Act` 16 · `SaveFacinProdEDMFire_Act` 15 ·
`ChangeEmailTextCeding_ACT` 14 · `SaveTreatyProduction_Act` 13.

**Arti flag `pyStepsPreCondition` tidak dijelaskan korpus.** Dua tafsir:
(a) *pra-kondisi dinonaktifkan* → langkah **selalu** jalan;
(b) *langkah dinonaktifkan* → langkah **tidak pernah** jalan.
Bukti tak-langsung condong ke (a): flag `true` hampir selalu berpasangan dengan kondisi
terisi (4.744 vs 2), yang berarti flag melacak keberadaan kondisi, bukan keberadaan langkah.
Tetapi ini **[dugaan]**, dan salah memilih akan mengubah hasil 577 langkah — 19 di antaranya
di `SaveOfferProduction_Act`, activity penulis produksi. **MEMBLOKIR** (§9).

### 4.4 Pra-kondisi konstan

11 langkah memakai pra-kondisi literal `1==1`:

| Berkas | Langkah | Flag |
| --- | --- | --- |
| `AttachFileGIS.xml` | 1.2, 1.3, 1.5 | `true` |
| `GetInsuredID.xml` | 1, 3, 4.3, 4.6, 4.8, 4.11 | `false` |
| `SetBusinessType_Act.xml` | 2 | `true` |
| `SetCauseOfDecline_ACT.xml` | 1 | `true` |

`1==1` selalu benar. Dengan tafsir (a) maupun (b) di §4.3, langkah-langkah ber-flag `true`
selalu jalan; jadi kondisinya **tidak punya efek** dan aman dihapus saat porting.
Tidak ditemukan satu pun pra-kondisi yang konstan **salah** (`1==2`, `"A"=="B"`),
jadi **tidak ada langkah yang terbukti mati lewat jalur pra-kondisi**. [terverifikasi]

### 4.5 Duplikasi

Normalisasi: ambil subtree `pySteps`, buang semua elemen ber-prefiks `px`/`pz`
(operator, timestamp, host), ratakan menjadi `jalur=nilai`, urutkan
`[StringComparer]::Ordinal`, SHA-256.

Hasil: **1 kelompok klon, 4 activity byte-identik**:

`ProtectionObjectBoilerPresure_Act` = `ProtectionObjectContractorPlant_Act` =
`ProtectionObjectLandRig_Act` = `ProtectionObjectMB_Act`

→ di Go cukup **satu** fungsi dengan parameter jenis objek. Tidak ada activity dengan
`pySteps` kosong. [terverifikasi]

---

## 5. Activity terbesar

### 5.1 Ukuran berkas berbanding lurus dengan jumlah langkah

Median byte per langkah di seluruh 609 = **9.981**. Sepuluh terbesar justru **di bawah**
median (5.681–10.185 B/langkah), artinya ukuran memang digerakkan oleh **banyaknya langkah**,
bukan oleh blob tertanam. Ukuran karena itu memang proksi yang sah untuk beban logika.
[terverifikasi]

### 5.2 Sepuluh terbesar

| # | Activity | Byte | Langkah (bersarang) | Kelas | Isi garis besar |
| --- | --- | --- | --- | --- | --- |
| 1 | `SetFlagDeleteEDM_Act` | 1.885.967 | 332 | `Data` | `Property-Set`×129, blok×123, `Page-Remove`×30, `Call`×25, `Apply-DataTransform`×20. Menandai penghapusan per lini bisnis (deskripsi langkah: `FIRE`, `location`, `objitem`, `coverage`). Dirujuk **8 Section**, tanpa pemanggil activity. |
| 2 | `UploadCSVPolicyMemberEDM_PostAct` | 1.775.722 | 271 | `Data-Policy` | `Property-Set`×178, blok×58, `Page-Copy`×12. Pasca-proses unggah CSV anggota polis; deskripsi menyebut pengambilan "regno terakhir" dan pemetaan Plan (`PlanDE`, `PlanGS`, `PlanIP`). Dipanggil `CallUploadCSVPolicyMember_PostAct`. |
| 3 | `GetKapasitasTreaty` | 1.323.132 | 167 | `Work` | `Property-Set`×94, blok×39, `RDB-List`×5, `Obj-Browse`×3. Kapasitas treaty; deskripsi menyebut cek "toprisk" dan satu ambang numerik hard-code. Dipanggil `CountPremiAndTSINusantaraRe_ACT`. |
| 4 | `SaveFacinProdEDMFire_Act` | 1.246.177 | 132 | `Work` | `Property-Set`×86, blok×25, **`RDB-List`×19**. Penulis produksi FIRE untuk EDM. Dipanggil `SaveFacinProdAllEDM_Act`. |
| 5 | `UploadCSVVehicle_PostAct` | 1.213.503 | 158 | `Work` | `Property-Set`×99, blok×21, `Call`×13, `Property-Map-DecisionTable`×9, `RDB-List`×4. Impor CSV kendaraan; memanggil 2 activity yang **hilang** (`SearchTemplateAdditionalCoverageSQL_PreAct`, `SetTsiCovAct`). Titik masuk dari FlowAction. |
| 6 | `SaveFacinRNWProd_Act` | 1.120.324 | 110 | `Work` | `Property-Set`×50, **`RDB-List`×31**, blok×29. Penulis produksi. Dipanggil `SaveTreatyProduction_Act`. |
| 7 | `SaveTreatyProduction_Act` | 1.048.214 | 103 | `Work` | `Property-Set`×44, **`RDB-List`×28**, blok×25. Penulis produksi treaty. Dipanggil `SaveEDMToJsonPolicy_Act`, `SaveJsonPolicyFacIn_Act`. |
| 8 | `SaveOfferProduction_Act` | 1.031.929 | 110 | `Work` | `Property-Set`×48, **`RDB-List`×31**, `RDB-Delete`×1. Penulis produksi offer. Langkah 11 `RDB-Delete` berdeskripsi **"insert to db"** — pengulangan persis jebakan nama yang sudah tercatat. Dipanggil `SaveJsonOfferFacIn_Act`, `UpdateJSONOffer_ACT`. |
| 9 | `SumTSIPremiSpreadedRNM_FIRE_Act` | 996.052 | 139 | `Data-Coverage` | `Property-Set`×100, blok×36. Penjumlahan TSI & premi ter-spread per mata uang. Dipanggil `SumTSIPremiSpreadedRNM_Act`. |
| 10 | `ProtectFIREMBUPA_Act` | 875.119 | 104 | `Work` | `Property-Set`×49, **`Page-Set-Messages`×25**, blok×15, `RDB-List`×5, `Obj-Browse`×4. Validasi/proteksi input lintas lini (MBU: plat & mesin; coverage & currency). Dipanggil `ProtectCoverage_Act`. |

**Pola:** enam dari sepuluh terbesar adalah **penulis produksi** (`Save*Prod*`,
`Save*Production*`), dan lima di antaranya mengandung 19–31 langkah `RDB-List`.
Beban bisnis terberat memang menumpuk di jalur tulis produksi.

```powershell
Get-ChildItem 'D:\migrasi\RNM\NB FacIn\Activity' -File -Filter *.xml |
  Sort-Object Length -Descending | Select-Object -First 10 Name,Length
```

[terverifikasi]

---

## 6. Java inline

**49 langkah Java di 32 activity** (5,3 % dari 609). Panjang total ±77 KB.

```powershell
(Import-Csv "$o\java.csv").Count                                   # 49
(Import-Csv "$o\java.csv" | Select-Object -Expand File -Unique).Count  # 32
```

### 6.1 Sebaran pola

| Pola | Langkah | Apa yang dilakukan |
| --- | --- | --- |
| **Dedup pagelist** | 8 | Iterasi mundur `ClipboardProperty.pxResults`, kunci dedup dirakit dari 1–2 properti dan digabung `#`, disimpan di `HashSet<String>`, baris duplikat dihapus `Estimasi.remove(i)`. **Padanan Go langsung: `map[string]struct{}` + filter slice.** |
| **PDF / lampiran** | 13 | Rangkaian `HTMLToPDF`, `HTMLRISlipToPDF`, `GeneratePDFFacOutMemoPlacing`, `GeneratePDFSlipKomisi`, `GenerateEndorsementPDF`, `AttachToWork`, `AttachRISlipToWork`. |
| **JSON** | 7 | `adoptJSONObject` — menempelkan string JSON dari hasil SQL ke halaman clipboard (`SetTreatyIn_Act` langkah 5, dibungkus `try/catch` dengan `oLog.error`). **Tidak ada padanan Go; harus ditulis ulang sebagai unmarshal bertipe.** |
| **Manipulasi string** | 2 | |
| **Ticket** | 1 | |
| **File I/O** | 1 | |
| **Lain** | 17 | `downloadFile` (2.578 char), `SendSimpleEmail`, `svcAddWorkObject`, `CallVirusCheck`, `ASMForceCaseClose`, `DeleteAttachment`, `UploadDoc_AI_V2`, dst. |

### 6.2 Dua blok Java terbesar bukan logika

`InputClause_ViewDtl_PreAct.xml` langkah 2 dan `InputClause_ViewArg_PreAct.xml` langkah 2:
**25.694 karakter, byte-identik satu sama lain**, dan isinya adalah **teks klausul polis
berformat HTML** yang di-hard-code sebagai literal string Java (baris pertama:
`curText = "<P …><B>KLAUSUL HURU-HARA, TERORISME DAN SABOTASE</B>…"`).

→ Ini **konten, bukan kode**. Di target harus pindah ke tabel/berkas template, bukan
diporting sebagai string Go. [terverifikasi]

### 6.3 Risiko porting

Contoh yang dikutip apa adanya (dedup, pola paling sering):

```java
// Asal: NB FacIn/Activity/CekTotalShareOffered_Act.xml, langkah 5
ClipboardPage Estimation = tools.findPage("testsi");
ClipboardProperty Estimasi = Estimation.getProperty(".pxResults");
HashSet<String> values = new HashSet<String>();
for(int i = Estimasi.size() ; i > 0 ; i--){
  ClipboardPage aPage = Estimasi.getPageValue(i);
  String prop2 = aPage.getProperty(".CARI2").getStringValue();
  String aValue = prop2;
  boolean isAdded = values.add(aValue);
  if(!isAdded){ Estimasi.remove(i); }
}
```

Seluruh 49 blok bergantung pada API clipboard Pega (`tools.findPage`,
`ClipboardPage`, `ClipboardProperty`, `oLog`). **Tidak ada padanan 1:1 di Go.** Delapan blok
dedup dan tujuh blok JSON dapat ditulis ulang secara mekanis; 13 blok PDF/lampiran dan
`downloadFile`/`SendSimpleEmail`/`CallVirusCheck` menyentuh subsistem Pega
(Link-Attachment, korespondensi, pemindai virus) yang **tidak ada padanannya** dan perlu
keputusan arsitektur tersendiri. [terverifikasi untuk isi · [dugaan] untuk tingkat kesulitan]

---

## 7. Parameter dan halaman kerja

### 7.1 `pyParameters` — memasok signature fungsi Go

| Ukuran | Nilai |
| --- | --- |
| Activity yang punya parameter | **269** dari 609 (44 %) |
| Total baris parameter | 672 |
| Rata-rata (yang punya) | 2,5 |

Sebaran jumlah parameter: 1 param → 121 activity · 2 → 70 · 3 → 37 · 4 → 12 · 5 → 9 ·
6 → 7 · 7 → 3 · 8 → 3 · 9 → 1 · 13 → 1 · 14 → 1 · 18 → 1 · 21 → 1 · 22 → 1 · 23 → 1.

Tipe: `STRING` 570 (85 %) · `INTEGER` 45 · `BOOLEAN` 21 · `Decimal` 14 · `DateTime` 5 ·
`PAGE` 5 · `JAVAOBJECT` 4 · `Double` 4 · `TrueFalse` 4.

Arah: `IN` 640 · `OUT` 28 · kosong 4.

Parameter tersering: `Nopolis` 22 · `Status` 18 · `ID` 18 · `NoEndors` 17 · `FromWhere` 13 ·
`ERRMSGOBJECT` 13 · `note` 9 · `idx` 9 · `coverage` 8 · `InSpreading` 7 · `IdxCoverage` 7 ·
`IdxFac` 7 · `enddate` 6 · `targetPage` 6 · `DiscountStatus` 6 · `IdxPropertyItem` 6 ·
`idxlocation` 6.

**Implikasi uang.** Hanya 14 parameter bertipe `Decimal` dan 4 bertipe `Double`; nilai uang
lewat sebagai `STRING`. Ini **menguatkan** aturan CLAUDE.md §4.1: di Go seluruh nilai uang
masuk sebagai string desimal lalu di-parse ke `decimal.Decimal`, **tidak pernah** `float64`.
Empat parameter `Double` yang ada adalah kandidat perbaikan — **keputusan bisnis**, bukan
keputusan migrasi. [terverifikasi]

**Implikasi penamaan.** Penomoran indeks memakai empat ejaan berbeda untuk konsep yang sama
(`idx`, `IdxCoverage`, `IdxFac`, `IdxLoc`, `idxlocation`, `IdxProperty`,
`IdxPropertyItem`). Normalisasi di Go harus dipetakan eksplisit, tidak boleh otomatis.

```powershell
$pr=Import-Csv "$o\params.csv"
($pr|Select-Object -Expand File -Unique).Count; $pr.Count      # 269 ; 672
$pr|Group-Object Type|Sort-Object Count -Descending
```

### 7.2 `pyPagesAndClasses` — halaman kerja

3.164 baris di **551 activity** (90 %). Halaman tersering:

| Jumlah | Halaman | Jumlah | Halaman |
| --- | --- | --- | --- |
| 433 | `pyWorkPage` | 31 | `Report`, `DataSearch` |
| 76 | `InputData` | 29 | `Track` |
| 65 | `pyReportContentPage` | 27 | `OldID`, `InputParam` |
| 50 | `OutputData` | 26 | `PolicyFacIn`(+`.pxResults`) |
| 48 | `pyReportContentPage.pxResults` | 25 | `OutputParam` |
| 40 | `CurrencyMaster`(+`.pxResults`) | 21 | `DataIN`, `OldIDBusiness`(+`.pxResults`) |
| 37 | `ProtectSpreading` | 20 | `DataIN1`, `DataOut`(+`.pxResults`) |

Kelas halaman tersering: `Code-Pega-List` 874 · `ASM-FW-GISFW-Data-Search` 861 ·
`ASM-FW-GISFW-Work` 293 · `ASM-FW-GISFW-Work-NB` 181 ·
`ASM-FW-GISFW-Data-SpreadingRisk` 59 · `ASM-FW-GISFW-Int-CURRENCY` 53 ·
`ASM-FW-GISFW-Int-policyjson` 46 · `Rule-Obj-Report-Definition` 42.

**Pola yang mengikat desain Go:** `ASM-FW-GISFW-Data-Search` muncul 861 kali. Ini kelas
halaman parameter SQL generik dengan properti bernomor (`CARI1`…`CARI23`, `HASIL1`…`HASIL10`
— terlihat langsung di badan SQL `MachingDataFacin_Sql`, `InsertViewSuggest_SQL`,
`GetSequenceNumber_SQL`). Repository Go **tidak boleh** meniru pola `CARIn/HASILn`; setiap
query harus diberi struct parameter bernama, dengan komentar asal rule seperti diwajibkan
CLAUDE.md §4.6. Pemetaan `CARIn` → makna kolom harus dibaca dari SQL masing-masing rule —
**belum dilakukan di pass ini**. [terverifikasi untuk angka · [pertanyaan terbuka] untuk
pemetaan `CARIn`]

```powershell
$pg=Import-Csv "$o\pages.csv"
$pg.Count; ($pg|Select-Object -Expand File -Unique).Count       # 3164 ; 551
```

### 7.3 Ketergantungan rule lain

Dari 20.703 `pxRuleReferences` yang dipancarkan 609 Activity:

| Tipe rule dirujuk | Jumlah ref |
| --- | --- |
| `Rule-Obj-Property` | 11.416 |
| `Rule-Obj-Class` | 4.503 |
| `Rule-Method` | 1.454 |
| `Rule-Obj-When` | 684 |
| `Rule-Obj-Activity` | 567 |
| `Rule-Utility-Library` | 557 |
| `Rule-Utility-Function` | 377 |
| `Rule-Connect-SQL` | 372 |
| `Rule-RDB-SQL` | 372 |
| `Rule-Obj-FieldValue` | 224 |
| `Rule-Obj-HTML` | 80 |
| `Rule-Obj-Model` | 53 |
| `Rule-Declare-Pages` | 14 |
| `Rule-Declare-DecisionTable` | 14 |
| `Rule-Message` | 8 |
| `Rule-Connect-REST` | 6 |
| `Rule-HTML-Fragment` / `Rule-Declare-DecisionTree` | 1 masing-masing |

**`When` yang dipakai Activity: 89 rule berbeda.** Dari 89 itu, **84 punya berkas** di
`NB FacIn/When/` dan **5 tidak**: `IsSpreadingDepan` (rule aplikasi — hilang dari ekspor),
serta 4 rule platform `pxWebStorageEnabled`, `pxCMISEnabled`, `pxRepositoryEnabled`,
`pzIsRequest_V2API`. Sebaliknya, 126 dari 210 berkas `When` di ekspor NB tidak dirujuk
activity mana pun — dirujuk Section/Flow/FlowAction.

Sepuluh teratas: `IsFire` 70 · `IsEDM` 64 · `IsMarineCargo` 55 · `IsAneka` 53 · `IsMBU` 51 ·
`IsPA` 41 · `IsGolfInsurance` 36 · `IsTravel` 28 · `IsLife` 18 · `IsMBD` 16.

→ Registry predikat `internal/rules/` harus menyediakan minimal 89 predikat agar
Activity NB dapat dijalankan; sisanya diperlukan lapisan UI. Satu di antaranya
(`IsSpreadingDepan`) **kondisinya tidak ada di ekspor** — perlakukan seperti lima rule
berkondisi kosong di CLAUDE.md §4.5: `panic`, bukan `return false`.

**`Apply-DataTransform`: 82 langkah, 36 DataTransform berbeda** dari 148 berkas
DataTransform. Terbanyak `SetTSIPremiSpreaded_FacIn` (31×) — kandidat kuat untuk
`internal/services/spreading/`.

**`Property-Map-DecisionTable`: 19 langkah** menunjuk antara lain `MappingCoverage` (6×),
`SetUploadHubAW1/2/3` (2× masing-masing), `MappingOutgoIndex`, `MappingOutgoIndex2`,
`MappingAdditionalCoverageIndex`, `LicensePlatRegion_DeT`, `GetMimeType`, `BusinessType_DeT`.
Satu `Property-Map-DecisionTree` menunjuk `Tree_ShortPeriod`. [terverifikasi]

### 7.4 Integrasi keluar

| Ukuran | Nilai |
| --- | --- |
| Langkah `Connect-REST` | 6 |
| Rule di `NB FacIn/ConnectREST/` | 3 |
| Activity yang menyebut `M_LINK_SERVICE` | 7 |

| Activity | `ServiceName` | `EndPointURL` di langkah |
| --- | --- | --- |
| `DeleteGoogleStorage_Act` (lang. 9) | `ServiceGoogle` | **kosong** |
| `GeminiAIGoogle_Act` (lang. 11) | `ServiceGoogle` | **kosong** |
| `GetUrlGoogleStorage_Act` (lang. 6.5) | `ServiceGoogle` | **kosong** |
| `InsertGoogleStorage_Act` (lang. 11) | `ServiceGoogle` | **kosong** |
| `GetPaymentList_Act` (lang. 6) | `getPremiumPaidOn` | **kosong** |
| `serviceInsertArasapasEDM_act` (lang. 6) | `convertJsonNusareToProduction` | **kosong** |

Seluruh enam langkah memakai `MethodName=POST`, `ExecutionMode=Run`, dan
**`EndPointURL` kosong** — URL diselesaikan saat runtime. Tujuh activity menyentuh
`M_LINK_SERVICE`, dan `GetLinkService` (fan-in 6, 4 langkah:
`Page-New` → `Obj-Browse` → `Property-Set` → `Page-Remove`) adalah pencari alamatnya.

Ini **menguatkan** CLAUDE.md §4.4: tidak ada satu pun endpoint literal di korpus Activity NB.
Isi `M_LINK_SERVICE` tetap **belum terverifikasi**. Tidak ada nilai endpoint, token, atau
alamat yang disalin ke dokumen ini. [terverifikasi]

`GeminiAIGoogle_Act` — model AI pihak ketiga di alur underwriting — memanggil
`GetTokenStorage_SQL` (prosedur `POOLDATA.GET_TOKEN_STORAGE`, `COMMIT`) lalu `Connect-REST`
ke `ServiceGoogle`. **Isinya tidak dibaca lebih jauh**; sesuai CLAUDE.md §6 pemindahannya
butuh persetujuan manusia dan tinjauan keamanan.

### 7.5 Guard berbasis identitas operator

Sesuai CLAUDE.md §3.5 dicatat **jumlah dan mekanismenya saja; tidak ada nama yang disalin**.

| Ukuran | Nilai |
| --- | --- |
| Langkah dengan guard identitas di pra-kondisi | 37 |
| Berkas terlibat | 11 |
| Bentuk ekspresi berbeda | 4 |
| Identitas literal berbeda | **2** |
| Guard identitas di transisi (`pyStepsTransParamsWhen`) | 0 |

Mekanisme: perbandingan `OperatorID.pyUserIdentifier` terhadap literal, di-OR dengan
sebuah kondisi bisnis. Berkas: `CountGrossPremiEDM_Act` (8), `CountGrossPremi_Act` (6),
`CountGPWMarinePAMbu_Act` (6), `InputOfferFacInEngineerUW_preACT` (3),
`CheckSpreadingProtect_ACT` (3), `InputOfferFacInEngineer_preACT` (2),
`CheckSpreadingProtectMCargoMBU_ACT` (2), `CheckSpreadingProtectFire_ACT` (2),
`CheckSpreadingProtectAnekaGolf_ACT` (2), `CheckSpreadingProtectPATravel_ACT` (2),
`SetValidateDateUW_PostAct` (1).

**20 dari 37 berada di activity perhitungan premi bruto.** Artinya identitas operator
ikut menentukan angka premi. Kandidat perbaikan — **keputusan bisnis**, bukan keputusan
migrasi (CLAUDE.md §1).

```powershell
$p=Import-Csv "$o\precond.csv"
$g=$p|?{$_.PCwhen -match '(?i)pyUserIdentifier|pxCreateOperator|OperatorID\.'}
$g.Count; ($g|Select-Object -Expand File -Unique).Count       # 37 ; 11
```

[terverifikasi]

---

## 8. Usulan pemetaan ke `internal/services/` — **[usulan]**

> **Seluruh seksi ini adalah [usulan], bukan temuan.** Pemetaan diturunkan dari
> (a) fakta terverifikasi "menulis ke Oracle" / "memanggil REST", dan (b) pola nama.
> Pola nama **bukan bukti perilaku**; setiap penempatan wajib diverifikasi ulang saat
> activity yang bersangkutan benar-benar diporting.

Aturan penempatan yang dipakai, berurutan (yang pertama cocok menang), lalu dua aturan
penimpa berbasis fakta:

1. nama memuat `Spreading|Speading|Protect|Protek|KapasitasTreaty|BreakDownSpread|CekLimit|LayerList|Retro` → `spreading`
2. nama diawali `Count|Calculate|Hitung|Sum|Average|ReCount|Total`, **atau** memuat `Premi|Rate|RIComm|Brokerage|Discount|AdminFee|Installment|Prorate|Loading` → `premium`
3. nama diawali `GetLimitAkseptasi|GetAksep|SetBanding|SetNBStatus|SetToInbox|…HistoryAkseptasi|…ViewSuggest|ASMForceCaseClose|SetTicket|SetCauseOfDecline|SetScore|ScoringResult` → `acceptance`
4. nama diawali `Save|Insert|Update|Delete|serviceInsert`, atau memuat `Production|Nopolis|GeneratePolicyNo` → `production`
5. nama memuat `Upload|Attach|Email|Correspondence|PDF|GoogleStorage|Gemini|LinkService|Document|VirusCheck|Export|Download|CSV|JSON|Arasapas|Storage|Print|svc|Service` → `integration`
6. sisanya → `faccase`
7. **penimpa fakta A:** activity yang terbukti menulis Oracle (§3.4) → `production`
8. **penimpa fakta B:** activity yang memuat langkah `Connect-REST` → `integration`

Hasil:

| Paket usulan | Activity | Catatan |
| --- | --- | --- |
| `faccase` | **275** | dipecah: `faccase/case` 173 · `faccase/lookup` 76 · `faccase/validation` 26 |
| `premium` | 97 | inti rumus premi; termasuk seluruh `Count*`/`Sum*` |
| `production` | 81 | **57 di antaranya terverifikasi menulis Oracle**; 5 penulis sisanya jatuh ke `integration` (Google Storage, Gemini, Arasapas) |
| `spreading` | 79 | termasuk `Protect*`/`Protection*` (15+14 berkas) |
| `integration` | 61 | unggah CSV, PDF, e-mail, Google Storage, REST |
| `acceptance` | 16 | mesin tangga persetujuan |
| **Total** | **609** | |

```powershell
& "$o\bucket2.ps1"      # mencetak tabel di atas; hasil di bucket2.csv
```

### 8.1 Catatan penting atas usulan

- **`acceptance` hanya 16 activity** padahal ini inti aplikasi. Penyebabnya: sebagian
  besar logika akseptasi hidup di rule `When` (89 dipakai) + `RDBList`
  (`GetLimitAkseptasi_SQL`, `GetAksepBanding_SQL`, `GetLimitAkseptasiJUWA_SQL`,
  `GetLimitAkseptasiNonFire*_SQL`, `GetLimitAkseptasiPreferedCommJUWA_SQL`,
  `GetLimitAkseptasiNonPreferJUWA_SQL`, `GetLimitAccEngineeringJUWA_SQL`) dan di Flow —
  **bukan** di badan Activity. Jangan simpulkan "akseptasi itu ringan".
- **`GetLimitAkseptasi_Act2` hilang dari ekspor** (§2.5) — lubang di paket `acceptance`.
- **`faccase` 275 masih terlalu gemuk** untuk satu paket Go. Pemecahan tiga arah
  (`case` / `lookup` / `validation`) di atas adalah titik awal; pemecahan final
  sebaiknya mengikuti agregat `OfferFacIn` dan bukan kata kerja nama.
- **`repository/lookup`**: `GetCurrencyMaster` (fan-in 19, tertinggi di seluruh korpus)
  dan `GetLinkService` (fan-in 6) bukan logika bisnis — keduanya pencari referensi murni
  dan lebih tepat di `internal/repository/lookup/` daripada di `services/`.
- **`money`**: seluruh activity `premium` (97) dan `spreading` (79) menyentuh nilai uang.
  Tidak satu pun boleh memakai `float64` (CLAUDE.md §4.1).

---

## 9. Pertanyaan terbuka

### MEMBLOKIR

| # | Pertanyaan | Mengapa memblokir | Bukti |
| --- | --- | --- | --- |
| **B-1** | Apa arti `pyStepsPreCondition` = `false`? Pra-kondisi dinonaktifkan (langkah selalu jalan) **atau** langkah dinonaktifkan (tidak pernah jalan)? | Mengubah perilaku **577 langkah**, termasuk 19 di `SaveOfferProduction_Act`, 16 di `SaveFacinRNWProd_Act`, 13 di `SaveTreatyProduction_Act` — semuanya penulis produksi. Salah pilih = selisih data produksi yang tidak dapat dijelaskan saat paralel run. | §4.3 |
| **B-2** | Bagaimana 21 pernyataan SQL penulis **tanpa `COMMIT`** di-commit? | Termasuk `InsertTreatyProduction_Sql`, `InsertOfferProduction_Sql`, `UpdateJsonOffer_SQL`, `UpdatePolisEndorsement_SQL`, `DeleteDataProduction`. Hanya ada 5 langkah `Commit` di seluruh 609 Activity dan tak satu pun di jalur produksi. Batas transaksi `internal/repository/oracle/` tidak dapat dirancang tanpa jawaban ini. | §3.3 |
| **B-3** | Di mana badan activity aplikasi yang hilang? **Sebagian ditutup K-006 (15 Sep 2026):** kesebelasnya tidak diminta diekspor ulang secara khusus. **Tetapi enam cabang pemanggilnya berstatus ⏸ DITANGGUHKAN, bukan kode mati** — Special Acceptance, tangga akseptasi putaran kedua, limit tambahan treaty type, kapasitas treaty, simpan produksi endorsement bonding, penanganan galat konversi produksi. | **Masih memblokir sebagian.** "Hilang dari ekspor" ≠ "usang": ekspor diambil dari titik waktu berbeda antar folder (R0), jadi ketiadaan bukan bukti tidak-dipakai. Diverifikasi terhadap **ekspor produksi tunggal**; bila activity-nya ada di sana → **wajib diport**. | `00-KEPUTUSAN-WORK-OWNER.md` K-006 |
| **B-4** | Apa isi 35 stored procedure `POOLDATA.*` yang dipanggil? | 62 activity menulis ke Oracle **melalui** prosedur ini. Tidak ada satu pun badan prosedur di korpus. Ini blok terbesar perilaku yang tidak dapat direkam dari ekspor Pega. | §3.5 |
| **B-5** | Alias koneksi `ASM` (529 langkah) vs `RNM` (87 langkah) menunjuk instance/skema Oracle apa? | Menentukan berapa pool koneksi yang dibutuhkan dan apakah transaksi lintas-alias mungkin. Nilai DSN tidak ada di korpus. | §3.1 |

### TIDAK MEMBLOKIR

| # | Pertanyaan | Catatan |
| --- | --- | --- |
| T-1 | Apa arti kode 1–6 pada `pyStepsPreCondParamsWhenTrue`/`WhenFalse` dan `pyStepsTransParams…`? | Pola terbanyak `true→2 false→3` (2.528×), lalu `true→(kosong) false→3` (1.480×), `true→3 false→2` (304×), `true→5 false→2` (100×). Enumerasi ini **belum terverifikasi** — tidak dijelaskan korpus. Tidak memblokir inventaris, tetapi **wajib dijawab sebelum menerjemahkan alur kendali** langkah mana pun. |
| T-2 | Apakah 5 activity tak terjangkau (§4.2) benar-benar mati? | Semuanya bertema coverage tambahan kendaraan. Perlu pencarian nama di ekspor siklus lain + log runtime. Di luar batas pass ini. |
| T-3 | Pemetaan `CARI1`…`CARI23` / `HASIL1`…`HASIL10` → makna kolom | Harus dibaca per rule SQL. Belum dilakukan; diperlukan sebelum struct parameter repository Go ditulis. |
| T-4 | 32 `pyRequestType` dirujuk tetapi rule SQL-nya tidak ada di ekspor NB | Termasuk `GetLimitAkseptasi1SA_Act`, `GetLimitAkseptasi2SA_Act`, `GenerateEndorsementNo`, `GenerateNoPolicyOJKID`, `SearchPctTreatyLimit`, `MachingDataOffer_Sql`, `INSERTERROROFFER_Sql`, `InsertOfferProductionBackup_Sql`. Dua yang pertama menyentuh akseptasi. |
| T-5 | `SearchJobID.xml` memuat `pyRequestType` = `SearchJobIDSQL`, sedangkan activity memanggil `SearchJobID` | Satu-satunya rule RDBList di ekspor NB yang tidak dipakai activity NB. Kemungkinan salah ketik yang sudah lama hidup, atau pemanggilnya ada di siklus lain. **Nama berkas ≠ kunci rule** — konfirmasi ulang sebelum memetakan query. |
| T-6 | Apakah `GetSequenceNumber_SQL` / `GetTokenStorage_SQL` benar-benar mengubah baris? | Keduanya memanggil prosedur dengan parameter `OUT` **dan** `COMMIT`. Dihitung sebagai penulis secara konservatif. Jawabannya bergantung pada B-4. |
| T-7 | `GeminiAIGoogle_Act` — model AI pihak ketiga di alur underwriting | Isi tidak dibaca (CLAUDE.md §6). Perlu tinjauan keamanan tersendiri sebelum dipindahkan. |
| T-8 | Empat parameter bertipe `Double` dan 14 `Decimal` | Kandidat perbaikan tipe uang. Keputusan bisnis, bukan keputusan migrasi. |
| T-9 | 37 guard identitas operator, 20 di antaranya di activity premi bruto | Kandidat perbaikan. Keputusan bisnis. |
| T-10 | Rule `When` **`IsSpreadingDepan`** dirujuk Activity NB tetapi berkasnya tidak ada di `NB FacIn/When/` | Kondisinya tidak dapat dibaca. Implementasikan sebagai `panic("IsSpreadingDepan: kondisi belum diketahui")` sesuai CLAUDE.md §4.5, **bukan** `return false`. Empat rule `When` lain yang hilang (`pxWebStorageEnabled`, `pxCMISEnabled`, `pxRepositoryEnabled`, `pzIsRequest_V2API`) adalah rule platform Pega. |

---

## Lampiran — artefak analisis

Skrip dan CSV perantara ada di scratchpad sesi
(`…\scratchpad\`), **tidak** di dalam repo maupun korpus:

| Berkas | Isi |
| --- | --- |
| `extract.ps1` → `acts.csv` | 609 baris metadata activity |
| `extract2.ps1` → `steps2.csv` | 11.090 langkah, rekursif, dengan jalur bersarang |
| `extract3.ps1` → `allrefs.csv` | 53.097 `pxRuleReferences` dari 2.083 berkas NB |
| `extract4.ps1` → `textrefs.csv` | 1.817 kecocokan nama activity di 1.474 berkas non-Activity |
| `extract5.ps1` → `stepcp.csv`, `java.csv` | parameter `pyStepsCallParams` per langkah; 49 blok Java |
| `precond.ps1` → `precond.csv` | pra-kondisi & transisi per langkah |
| `sqlx.ps1` → `sql.csv` | 221 pernyataan SQL dari 217 rule RDBList |
| `graph.ps1` → `edges.csv`, `edges_idx.csv`, `summary.csv`, `missing.csv` | graf panggilan |
| `reach.ps1` → `unreach.csv` | BFS dari 394 titik masuk |
| `clone.ps1` → `clone.csv` | hash `pySteps` ternormalisasi |
| `bucket2.ps1` → `bucket2.csv` | usulan pemetaan paket |

Tidak ada berkas di dalam `D:\migrasi\RNM\NB FacIn\`, `RNW Fac In\`,
`Endorsment Fac In\`, atau `OUTPUT\_ARSIP-lintas-siklus\` yang dibuat, diubah, atau dihapus.
