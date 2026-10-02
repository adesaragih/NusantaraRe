# 03-ui / 01 — Inventaris layar & komponen Facultative Inward (target React)

**Sumber tunggal:** `D:\migrasi\RNM\NB FacIn\`, `D:\migrasi\RNM\RNW Fac In\`,
`D:\migrasi\RNM\Endorsment Fac In\` — ekspor rule Pega, dibaca read-only.
Tidak ada bahan dari luar ketiga folder itu. Discovery lama (`RNM_BRD/`) **tidak dipakai**.

**Tanggal audit:** 2026-09-15 · setiap angka di bawah disertai perintah PowerShell yang
menghasilkannya. Tidak ada angka yang disalin dari dokumen lain.

**Singkatan folder** yang dipakai sepanjang dokumen: `NB` = `NB FacIn`, `RNW` = `RNW Fac In`,
`EDM` = `Endorsment Fac In`.

---

## 0. Cara dokumen ini membaca korpus UI

Tiga tipe rule membentuk lapisan UI Pega. Hubungan antar-ketiganya **tidak** terbaca dari nama
berkas; harus dibaca dari tag di dalam XML. Tag yang dipakai dokumen ini:

| Tag | Ada di | Arti | Label |
| --- | --- | --- | --- |
| `<pySection>` | Section, Harness | section yang disisipkan (embedded section) | [terverifikasi] |
| `<pyStreamName>` | Harness | section yang dirender di region harness | [terverifikasi] |
| `<pySectionReference>` | FlowAction | section yang dirender oleh flow action itu | [terverifikasi] |
| `<pyPageListProperty>` | Section | properti PageList yang menjadi sumber baris grid | [terverifikasi] |
| `<pyRepeatDirection>RepeatGrid</…>` | Section | layout berulang bertipe grid | [terverifikasi] |
| `<pyShowModalDialog>true</…>` | Section | aksi yang membuka jendela modal | [terverifikasi] |
| `<pyCondition>` | Section | ekspresi visibilitas/read-only kustom | [terverifikasi] |

### 0.1 Jebakan yang sudah ditemukan dan dihindari — [terverifikasi]

**(a) `pyStreamName` pada FlowAction BUKAN nama section.** Pada seluruh 754 berkas FlowAction,
`<pyStreamName>` berisi nama flow action itu sendiri (karena `pyRuleFormType` = `Harness`), bukan
section yang dirender. Section sesungguhnya ada di `<pySectionReference>`.

Bukti — `D:\migrasi\RNM\NB FacIn\FlowAction\CauseOfLossGrid.xml`:

```
<pyStreamName>CauseOfLossGrid</pyStreamName>        <- nama flow action
<pySectionReference>InputLossExperiance</pySectionReference>   <- section sesungguhnya
<pyLabel>Input Pengalaman Kerugian</pyLabel>
```

```powershell
Select-String -Path "D:\migrasi\RNM\NB FacIn\FlowAction\CauseOfLossGrid.xml" `
  -Pattern '<pyStreamName>|<pySectionReference>|<pyLabel>|<pyRuleFormType>' |
  ForEach-Object { "L{0}: {1}" -f $_.LineNumber, $_.Line.Trim() }
```

Jumlah FlowAction yang `pyStreamName` = nama sendiri, dari 754 berkas: **754**.

```powershell
$n=0; foreach($f in @("NB FacIn","RNW Fac In","Endorsment Fac In")){
  Get-ChildItem "D:\migrasi\RNM\$f\FlowAction" -File -Filter *.xml | ForEach-Object {
    $c=[IO.File]::ReadAllText($_.FullName)
    if([regex]::Match($c,'<pyStreamName>([^<]*)</pyStreamName>').Groups[1].Value -eq $_.BaseName){$n++} } }
"pyStreamName = nama sendiri : $n"
```

**(b) Perbandingan `Get-FileHash` mentah menghasilkan nol layar identik — dan itu artefak, bukan
temuan.** Setiap berkas memuat metadata ekspor yang berbeda per instance Pega
(`pxHostId`, `pxCommitDateTime`, `pyShowJavaWindowName`, `pzIndexCount`, urutan serialisasi
elemen). Dokumen ini karena itu memakai **dua** ukuran yang dilaporkan terpisah (§5).

**(c) Korpus dirakit dari lebih dari satu instance Pega** — sebaran `pxHostId` pada 2.136 berkas UI:

```powershell
foreach($f in @("NB FacIn","RNW Fac In","Endorsment Fac In")){ "=== $f ==="; $h=@{}
 foreach($t in @("Harness","Section","FlowAction")){
  Get-ChildItem "D:\migrasi\RNM\$f\$t" -File -Filter *.xml | ForEach-Object {
   $c=[IO.File]::ReadAllText($_.FullName)
   $v=[regex]::Match($c,'<pxHostId>([^<]*)</pxHostId>').Groups[1].Value
   if(-not $v){$v='(kosong)'}; $h[$v]=1+$h[$v] } }
 $h.GetEnumerator()|Sort-Object Value -Descending|ForEach-Object{"  {0,-6} {1}" -f $_.Value,$_.Key} }
```

| Host | NB | RNW | EDM |
| --- | --- | --- | --- |
| `pega-nusre` | 229 | 210 | 208 |
| `819804e741a72afc8f52d7483eb2af0d` | 206 | 190 | 241 |
| `a11207acce0fa39e49c3b08f8b22d60f` | 205 | 196 | 217 |
| `jboss1073` | 76 | 75 | 0 |
| `jboss117` | 1 | 0 | 44 |
| `jboss122117` | 2 | 2 | 22 |
| `e19f736f394dfb2b34c997a08a129e91` | 3 | 3 | 0 |
| `c489159d74914edfec1e1fff24865df7` | 0 | 0 | 3 |
| `10.225.68.27_envhyd83-web-2` | 1 | 1 | 1 |

`jboss1073` hanya menyentuh NB+RNW; `jboss117` praktis hanya EDM. **[dugaan]** ekspor NB dan RNW
diambil dari lingkungan yang sama, EDM dari lingkungan lain — konsisten dengan temuan §5.

---

## 1. Jumlah berkas dan identitas

### 1.1 Per folder

```powershell
foreach ($f in @("NB FacIn","RNW Fac In","Endorsment Fac In")) {
  Write-Output "=== $f ==="
  foreach ($t in @("Harness","Section","FlowAction")) {
    "{0,-12} {1}" -f $t, (Get-ChildItem "D:\migrasi\RNM\$f\$t" -File -Filter *.xml).Count } }
```

| Tipe rule | NB | RNW | EDM | **Total berkas** |
| --- | ---: | ---: | ---: | ---: |
| `Harness` | 41 | 41 | 37 | **119** |
| `Section` | 432 | 397 | 434 | **1.263** |
| `FlowAction` | 250 | 239 | 265 | **754** |
| **Jumlah** | **723** | **677** | **736** | **2.136** |

### 1.2 Identitas unik gabungan (nama berkas di-dedup lintas tiga folder)

```powershell
$folders=@("NB FacIn","RNW Fac In","Endorsment Fac In"); $g=0
foreach($t in @("Harness","Section","FlowAction")){
  $all=@(); foreach($f in $folders){ $all += (Get-ChildItem "D:\migrasi\RNM\$f\$t" -File -Filter *.xml).BaseName }
  $u=($all|Sort-Object -Unique).Count; $g+=$u
  "{0,-12} berkas={1,-5} identitas-unik={2}" -f $t,$all.Count,$u }
"TOTAL identitas unik = $g"
```

| Tipe rule | Berkas | **Identitas unik** |
| --- | ---: | ---: |
| `Harness` | 119 | **45** |
| `Section` | 1.263 | **520** |
| `FlowAction` | 754 | **307** |
| **Jumlah** | **2.136** | **872** |

Konteks: 2.136 berkas UI dari **6.071** berkas `.xml` di ketiga folder = **35,2 %**.

```powershell
(Get-ChildItem "D:\migrasi\RNM\NB FacIn","D:\migrasi\RNM\RNW Fac In","D:\migrasi\RNM\Endorsment Fac In" `
  -Recurse -File -Filter *.xml | Measure-Object).Count
```

### 1.3 Nama dipakai ulang lintas tipe rule — [terverifikasi]

**139** nama muncul pada lebih dari satu tipe rule (umumnya FlowAction + Section dengan nama sama;
sebagian Harness + Section). Konsekuensi: **nama saja tidak mengidentifikasi artefak**; di React,
kunci identitas harus `tipe + nama`.

```powershell
$hn=@{}; foreach($t in @("Harness","Section","FlowAction")){
 foreach($f in @("NB FacIn","RNW Fac In","Endorsment Fac In")){
  foreach($n in (Get-ChildItem "D:\migrasi\RNM\$f\$t" -File -Filter *.xml).BaseName){
   if(-not $hn.ContainsKey($n)){$hn[$n]=@{}}; $hn[$n][$t]=$true } } }
($hn.GetEnumerator()|Where-Object{$_.Value.Count -gt 1}).Count
```

---

## 2. Daftar HARNESS — kandidat halaman utama React

45 identitas harness. Kolom **Section** diisi dari `<pySection>` + `<pyStreamName>` di berkas
harness, disaring ke nama yang benar-benar ada sebagai berkas `Section\*.xml` di korpus
**[terverifikasi]**. Kolom **Kelas** dari `<pxInsName>` (bagian sebelum `!`) **[terverifikasi]**.
Kolom **Cabang** dari §5.

```powershell
$folders=@("NB FacIn","RNW Fac In","Endorsment Fac In")
$secSet=@{}; foreach($f in $folders){ Get-ChildItem "D:\migrasi\RNM\$f\Section" -File -Filter *.xml |
  ForEach-Object { $secSet[$_.BaseName.ToUpper()]=$true } }
foreach($f in $folders){ Get-ChildItem "D:\migrasi\RNM\$f\Harness" -File -Filter *.xml | ForEach-Object {
  $c=[IO.File]::ReadAllText($_.FullName)
  $refs = @()
  $refs += [regex]::Matches($c,'<pySection>([^<]+)</pySection>')|ForEach-Object{$_.Groups[1].Value}
  $refs += [regex]::Matches($c,'<pyStreamName>([^<]+)</pyStreamName>')|ForEach-Object{$_.Groups[1].Value}
  $hit = $refs | Sort-Object -Unique | Where-Object { $secSet.ContainsKey($_.ToUpper()) }
  "{0,-18} {1,-28} {2}" -f $f,$_.BaseName,($hit -join ', ') } }
```

| # | Harness | Folder | Cabang | Kelas Pega | Section yang dirujuk |
| ---: | --- | --- | --- | --- | --- |
| 1 | `AccumulationRisk` | NB/RNW/EDM | sama | `DATA-PORTAL` | `AccumulationRisk` |
| 2 | `AdjustmentRiskAccumulation` | NB/RNW/EDM | sama | `DATA-PORTAL` | `AdjustmentRiskAccumulation` |
| 3 | `CedingCedant` | NB/RNW/EDM | sama | `ASM-FW-GISFW-DATA-QUOTATION` | `CedingCedant`, `CedingCedantHierarki` |
| 4 | `CedingCompany` | NB/RNW/EDM | **EDM beda** | NB/RNW `…DATA-OFFERTREATYIN` · EDM `…DATA-QUOTATION` | `CedingCoHierarki`, `CedingCompany`, `InputQuotation` |
| 5 | `ChooseAccumulation_FacIn` | NB/RNW/EDM | sama | `ASM-FW-GISFW-DATA-COVERAGE` | `ChooseAccumulation_FacIn` |
| 6 | `ChooseClauseFire` | NB/RNW/EDM | **EDM beda** | NB/RNW `ASM-FW-GISFW-WORK` · EDM `…DATA-CLAUSE` | `ChooseClauseFire`, `InputDtlClause` |
| 7 | `ChooseDeductible` | NB/RNW/EDM | sama | `ASM-FW-GISFW-DATA-DEDUCTIBLE` | `ChooseDeductible` |
| 8 | `ChooseInsured` | RNW/EDM | bercabang | `ASM-FW-GISFW-DATA-QUOTATION` | `ChooseInsured`, `ChooseInsuredDtl` |
| 9 | `ChooseOccupation` | NB/RNW/EDM | sama | `ASM-FW-GISFW-DATA-OCCUPATION` | `ChooseOccupation`, `ChooseOccupation_Dtl` |
| 10 | `ChooseRiskAddress` | NB/RNW/EDM | sama | `…DATA-OFFERFACIN-LOCATIONREINSURANCE` | `ChooseRiskAddress` |
| 11 | `ChooseRiskLocation` | NB/RNW/EDM | sama | `…DATA-OFFERFACIN-LOCATIONREINSURANCE` | `ChooseRiskLocation`, `GridRiskAddress`, `InputRiskAddress` |
| 12 | `HistoricalSurveyReport` | NB/RNW/EDM | sama | `ASM-FW-GISFW-DATA-OFFERFACIN` | `HistoricalSurveyReport`, `HistoricalSurveyReportDtl` |
| 13 | `HistoricalSurveyReportUW` | NB/RNW/EDM | **EDM beda** | NB/RNW `…DATA-POLICYTREATYIN` · EDM `…DATA-OFFERFACIN` | `HistoricalSurveyReportDtlUW`, `HistoricalSurveyReportUW` |
| 14 | `InputCauseOfDecline` | NB/RNW | – | `ASM-FW-GISFW-WORK` | `InputCauseOfDecline` |
| 15 | `Medical_Harnes` | **EDM saja** | – | `DATA-PARTY-PERSON` | `InputMedical1_Sec`, `InputMedical2_Sec`, `Medical_Harnes`, `Medical_Sec` |
| 16 | `PrintRISlip` | NB/RNW/EDM | sama | `ASM-FW-GISFW-DATA-FACOFFER` | `PrintRISlip`, `PrintRISlip_Endorsement` |
| 17 | `PrintRISlips` | NB/RNW/EDM | sama | `ASM-FW-GISFW-WORK` | `PrintRISlip`, `PrintRISlips`, `PrintRISlips_Endorsement` |
| 18 | `RiskAccumulationReport` | NB/RNW/EDM | sama | `DATA-PORTAL` | `RiskAccumulationReport` |
| 19 | `SelectAgent` | NB/RNW/EDM | sama | `ASM-FW-GISFW-INT-LLOYDAGENT` | `SelectAgent` |
| 20 | `SelectCoverage` | NB/RNW | – | `ASM-FW-GISFW-INT-CONDITION` | `SelectCoverage` |
| 21 | `SelectShip` | NB/RNW/EDM | sama | `ASM-FW-GISFW-INT-SHIP` | `SelectShip` |
| 22 | `SFAPortalEndorsement` | **EDM saja** | – | `DATA-PORTAL` | `SFAPortal_EndorsementHeader`, `SFAPortalEndorsement` |
| 23 | `SFAPortalOpportunities` | **NB saja** | – | `PEGACRM-PORTAL` | `SFAPortal_Opportunities`, `SFAPortalOpportunities`, `SFAPortalOpportunitiesHeader` |
| 24 | `ShowCedingCoList` | NB/RNW/EDM | sama | `ASM-FW-GISFW-DATA-OFFERFACIN` | `ShowCedingCoList` |
| 25 | `ShowPolis` | NB/RNW | – | `ASM-FW-GISFW-WORK` | `ShowPolis`, `ShowPolis_sc` |
| 26 | `SOB` | NB/RNW | – | `ASM-FW-GISFW-DATA-OFFERTREATYIN` | `SOB`, `SourceHierarki` |
| 27 | `SummaryRiskAccumulation` | NB/RNW/EDM | sama | `DATA-PORTAL` | `SummaryRiskAccumulation` |
| 28 | `TabbedScreenFlow7` | NB/RNW | – | `ASM-FW-GISFW-WORK` | hanya section bawaan Pega (`pyCaseHeaderOuter`, `pyTabbedScreenFlow7Main`, `pyTabbedScreenFlow7Nav`, `pyScreenFlow7Footer`, …) |
| 29 | `TotalAccumulation_FacIn` | NB/RNW/EDM | sama | `ASM-FW-GISFW-WORK` | `TotalAccumulation_FacIn` |
| 30 | `TotalAccumulationDtl` | NB/RNW/EDM | **EDM beda** | `ASM-FW-GISFW-DATA-COVERAGE` (sama) | `TotalAccumulationDtl` |
| 31 | `ViewClaimList` | NB/RNW/EDM | sama | `ASM-FW-GISFW-WORK` | `ViewClaimList` |
| 32 | `ViewCoverageFacOut` | NB/RNW/EDM | **EDM beda** | NB/RNW `…DATA-PROPERTYITEM` · EDM `…DATA-ANEKA` | `ViewCoverageFacOut`, `ViewCoverageFacOutShow` |
| 33 | `ViewCoverageFacOut_isUW` | NB/RNW/EDM | **EDM beda** | NB/RNW `…DATA-PROPERTYITEM` · EDM `…DATA-ANEKA` | `ViewCoverageFacOut_isUW`, `ViewCoverageFacOutShow_isUW` (EDM: `…_Is_UW`) |
| 34 | `ViewCSVResult_B2BHost` | NB/RNW | – | `ASM-FW-GISFW-DATA-OFFERFACIN` | `ViewCSVAneka`, `ViewCSVCredit`, `ViewCSVMarine`, `ViewCSVResult_B2BHost`, `ViewCSVResult_TableB2B` |
| 35 | `ViewDataOfferFacIn` | NB/RNW/EDM | sama | `ASM-FW-GISFW-WORK` | `CoverageCommisionList`, `CoverageSpreadingList`, `InputDtlObject`, `InputDtlObject_FacIn`, `InputOtherObjectAneka_UW`, `PaymentCurrencyList`, `ViewInwardFacultativeDtl_IsUW` |
| 36 | `ViewFacretro` | NB/RNW/EDM | sama | `ASM-FW-GISFW-WORK` | `InputFacOffer_IsUW`, `ViewFacretro` |
| 37 | `ViewFollowingNB` | NB/RNW | – | `ASM-FW-GISFW-WORK` | `GridViewFollowingNB`, `ViewFollowingNB`, `ViewOpenFollowingPolicy` |
| 38 | `ViewLetter` | NB/RNW/EDM | sama | `ASM-FW-GISFW-WORK` | `FacultativeLetter`, `ViewLetter` |
| 39 | `ViewOfferStatus` | NB/RNW/EDM | sama | `ASM-FW-GISFW-WORK` | `OfferStatusFacOut`, `ViewOfferStatus` |
| 40 | `ViewOldData` | **EDM saja** | – | `ASM-FW-GISFW-WORK` | `CoverageCommisionList`, `InputDtlObject(_FacIn)`, `InputInwardFacultative`, `InputOtherObjectAneka_FacIn`, `OldDataEndorsementDtl`, `PaymentCurrencyList` |
| 41 | `ViewOldDeductible` | NB/RNW/EDM | sama | `ASM-FW-GISFW-DATA-COVERAGE` | `ViewOldDeductible` |
| 42 | `ViewOldEndorsement` | NB/RNW/EDM | sama | `ASM-FW-GISFW-WORK` | sama isi dengan `ViewOldData` |
| 43 | `ViewOldOffer` | NB/RNW/EDM | sama | `ASM-FW-GISFW-WORK` | `OldOfferStatusFacOut`, `ViewOldOffer` |
| 44 | `ViewPictureList` | NB/RNW/EDM | sama | `ASM-FW-GISFW-INT` | `ReasViewAttachment`, `ViewPictureList` |
| 45 | `ViewPolis` | NB/RNW/EDM | sama | `ASM-FW-GISFW-WORK` | `CoverageCommisionList`, `CoverageSpreadingList`, `InputDtlObject(_FacIn)`, `InputInwardFacultativeDtl_IsUW`, `InputOtherObjectAneka_UW`, `PaymentCurrencyList`, `ViewDtlCoverageFacOut`, `ViewObjectSavior`, `ViewTabGroupPolisFacOutSavior` |

### 2.1 Harness BUKAN peta halaman yang lengkap — [terverifikasi]

Hanya **4 dari 45** harness merujuk lebih dari 5 section (`ViewPolis` 9, `ViewDataOfferFacIn` 6,
`ViewOldData` 6, `ViewOldEndorsement` 6); **39 harness merujuk 1–3 section**. Penutupan
transitif (harness → section → sub-section) terbesar adalah **44 section**
(`ViewFollowingNB` di NB, `ViewPolis` di EDM):

```powershell
$e = @()   # bangun graf section dulu (lihat §3.1), lalu:
# lihat skrip lengkap penutupan di §3.1; keluarannya:
#   ViewFollowingNB (NB)  langsung=2  closure=44
#   ViewPolis (EDM)       langsung=9  closure=44
```

Artinya **layar kerja sesungguhnya tidak digantung di Harness**, melainkan di FlowAction
(§4): 98 FlowAction menjangkau ≥ 20 section secara transitif. Harness di korpus ini berperan
sebagai **jendela pendukung** (pemilih/modal, cetak, tampilan riwayat) dan **portal**
(`DATA-PORTAL`, `PEGACRM-PORTAL`).

**[pertanyaan terbuka]** Harness mana yang menjadi titik masuk portal underwriter sehari-hari tidak
dapat ditentukan dari `Harness\`, `Section\`, `FlowAction\` saja — informasi itu ada di rule
`Rule-Portal` / access group, yang **tidak ada** di korpus.

---

## 3. Daftar SECTION per tema

520 identitas section. Pengelompokan berikut dibuat **dari pola penamaan** dan karena itu
berlabel **[dugaan]** secara keseluruhan; baris yang isinya sudah dibuka diberi label
[terverifikasi] tersendiri di §6.

```powershell
$o = (Get-ChildItem "D:\migrasi\RNM\NB FacIn\Section","D:\migrasi\RNM\RNW Fac In\Section",`
  "D:\migrasi\RNM\Endorsment Fac In\Section" -File -Filter *.xml).BaseName | Sort-Object -Unique
$pat=[ordered]@{
 'Coverage/Peril'='Coverage|Peril|Clause|Warranty|Deductible'
 'Obyek/Risiko'='Object|Obyek|Okupasi|Occupation|Property|Vehicle|Ship|Cargo|Goods|Conveyance|Location|Address'
 'Akseptasi/UW'='Scoring|Akseptasi|Approval|Decline|Accept|Limit|Treaty|UW'
 'Spreading/Share'='Spreading|Share|Layer|FacOffer|FacOut|Facretro'
 'Peserta/Person'='Person|Participant|Insured|Member|Medical|Benefit|Plan|Life|PA_|Travel'
 'Premi/Pembayaran'='Payment|Premi|Commision|Commission|Installment|Currency'
 'Korespondensi/Dokumen'='Email|Letter|Correspondence|Attach|Print|RISlip|Report|Upload|CSV'
 'Endorsement'='Endorsement|Endorsment|Edm|OldData'
 'Akumulasi'='Accumulation|Akumulasi'
 'Mitra/Cedant'='Ceding|Cedant|Agent|Source|Marketing|Client|Broker'
 'Wilayah/Lookup'='City|District|Province|ZipCode|Country|Periode|Hierarki'
 'Renewal'='Renewal|Renew' }
$as=@{}; foreach($k in $pat.Keys){ $m=$o|Where-Object{$_ -match $pat[$k] -and -not $as.ContainsKey($_)}
  foreach($x in $m){$as[$x]=$k}; "{0,-24} {1}" -f $k,$m.Count }
"{0,-24} {1}" -f 'Lain-lain',(($o|Where-Object{-not $as.ContainsKey($_)}).Count)
```

> Catatan metode: pola diterapkan **berurutan** dan setiap nama hanya masuk satu kelompok
> (kelompok pertama yang cocok menang). Jadi `InputCoverageLife_FacIn` masuk *Coverage/Peril*,
> bukan *Peserta/Person*. Angka di bawah hanya sah dengan urutan pola persis seperti skrip di atas.

| Tema (urutan pola) | Jumlah section unik | Contoh |
| --- | ---: | --- |
| Coverage / Peril / Klausula / Deductible | **127** | `InputCoverageFire`, `CoverageItem`, `CoverageSpreadingList`, `InputDtlClause_FacIn`, `InputDeductible_FacIn` |
| Obyek / Risiko / Okupasi / Lokasi | **106** | `InputDtlObject_FacIn`, `InputRiskAddress`, `ObjectList`, `InputOkupasiAneka_FacIn`, `VehicleGrid` |
| Akseptasi / UW / Limit / Treaty | **49** | `ScoringRisk`, `InputDtlMemorandumAccept`, `FormulaTreatyCapacityDesc`, `TableOfLimit` |
| Spreading / Share / Layer / FacOut | **41** | `SpreadingItem`, `InputDtlSpreadingCoverage_FacIn`, `OfferStatusFacOut`, `AddFacOfferList` |
| Peserta / Person / Benefit / Plan | **38** | `InputPerson`, `InputDtlParticipantPA`, `SelectBenefitList`, `Medical_Sec` |
| Premi / Pembayaran / Komisi / Mata uang | **22** | `PaymentCurrencyList`, `CoverageCommisionList`, `InputDtlPayment_FacIn`, `InstallmentList` |
| Korespondensi / Dokumen / Unggah | **17** | `EmailSection`, `FacultativeLetter`, `PrintRISlip`, `ReasViewAttachment`, `UploadMember` |
| Endorsement | **12** | `InputEndorsementDtl`, `OldDataEndorsementDtl`, `EdmType1/2/3`, `PeriodeEndorsement` |
| Akumulasi | **10** | `AccumulationRisk`, `TotalAccumulation_FacIn`, `SummaryRiskAccumulation` |
| Mitra / Cedant / Sumber bisnis | **8** | `CedingCedant`, `CedingCoHierarki`, `SourceHierarki`, `EditMarketing` |
| Wilayah / Lookup | **6** | `InputCity`, `InputDistrict`, `InputProvince`, `Periode` |
| Renewal | **4** | `InputRenewal`, `InputRenewalDtl`, `PeriodeRenewal` |
| **Lain-lain** | **80** | `InwardFacIn`, `OfferFacIn_NusaRe`, `InputFEA`, `InputLossExperiance`, `crm*`, `*SummarySection` |

### 3.1 Section paling banyak dirujuk = komponen bersama — [terverifikasi]

Graf rujukan dibangun dari `<pySection>` (Section + Harness), `<pyStreamName>` (Harness) dan
`<pySectionReference>` (FlowAction), lalu disaring ke target yang ada sebagai berkas Section.
Hasil: **2.040 sisi**, **1.798** di antaranya menunjuk section yang ada di korpus.

```powershell
$folders=@("NB FacIn","RNW Fac In","Endorsment Fac In")
$secSet=@{}; foreach($f in $folders){ Get-ChildItem "D:\migrasi\RNM\$f\Section" -File -Filter *.xml |
  ForEach-Object { $secSet[$_.BaseName.ToUpper()]=$true } }
$edges=New-Object System.Collections.ArrayList
foreach($f in $folders){ foreach($t in @("Harness","Section","FlowAction")){
  Get-ChildItem "D:\migrasi\RNM\$f\$t" -File -Filter *.xml | ForEach-Object {
   $c=[IO.File]::ReadAllText($_.FullName); $src=$_.BaseName; $refs=@()
   $refs += [regex]::Matches($c,'<pySection>([^<]+)</pySection>')|ForEach-Object{$_.Groups[1].Value}
   $refs += [regex]::Matches($c,'<pySectionReference>([^<]+)</pySectionReference>')|ForEach-Object{$_.Groups[1].Value}
   if($t -eq 'Harness'){ $refs += [regex]::Matches($c,'<pyStreamName>([^<]+)</pyStreamName>')|ForEach-Object{$_.Groups[1].Value} }
   foreach($r in ($refs|Sort-Object -Unique)){
     if($t -eq 'Section' -and $r -ceq $src){continue}
     [void]$edges.Add([PSCustomObject]@{Folder=$f;SrcType=$t;Src=$src;Ref=$r;Known=$secSet.ContainsKey($r.ToUpper())}) } } } }
"sisi=$($edges.Count) ; ke section dikenal=$(($edges|Where-Object{$_.Known}).Count)"
$edges | Where-Object{$_.Known} | Group-Object Ref | Sort-Object Count -Descending |
  Select-Object -First 30 Count,Name
```

> **Peringatan metode:** rujukan harness→section dengan **nama yang sama** (mis. harness
> `AccumulationRisk` → section `AccumulationRisk`) mudah hilang bila penyaring "buang self-reference"
> diterapkan tanpa melihat tipe. Skrip di atas hanya membuang self-reference **Section→Section**.
> Tanpa koreksi itu, jumlah section yatim salah naik dari 141 menjadi 157.

| Section | Rujukan (pasangan folder+sumber) | Rujukan (sumber unik, lintas-folder digabung) |
| --- | ---: | ---: |
| `PaymentCurrencyList` | **83** | 38 |
| `InputInwardFacultative` | **62** | 28 |
| `CoverageSpreadingList` | **60** | 27 |
| `CoverageCommisionList` | **53** | 22 |
| `InputDtlObject_FacIn` | **41** | 19 |
| `InputEndorsementDtl` | **36** | 16 |
| `InputInwardFacultativeDtl` | **34** | 14 |
| `InwardFacIn` | **31** | 14 |
| `OfferFacIn_NusaRe` | **31** | 13 |
| `InputDtlCoverage_FacIn` | 23 | 8 |
| `CoverageItem` | 23 | 10 |
| `EmailSection` | 22 | 8 |
| `PropertyItemListCoverage` | 22 | 9 |
| `InputCoverageFire` | 19 | 9 |
| `AddFacOfferList` | 19 | 7 |
| `SpreadingItem` | 19 | 8 |
| `InputRiskAddress` | 18 | 6 |
| `InputDtlSpreadingCoverage_FacIn` | 13 | 5 |
| `InputFacOffer` / `Property` / `PropertyItemList` / `InputDtlTrading` / `PropertyItemListCoverageSpreading` / `ShowCoverageFacOut` / `ScoringRisk` | 12 | 4 |

### 3.2 Section yatim — [terverifikasi]

**141 dari 520** identitas section tidak dirujuk oleh harness, section, maupun flow action mana pun
di dalam korpus.

```powershell
# lanjutan skrip §3.1
$ref=@{}; foreach($r in $edges){ $ref[$r.Ref.ToUpper()]=$true }
$all=(Get-ChildItem "D:\migrasi\RNM\*\Section" -File -Filter *.xml).BaseName | Sort-Object -Unique
"unik=$($all.Count) dirujuk=$(($all|Where-Object{$ref.ContainsKey($_.ToUpper())}).Count) yatim=$(($all|Where-Object{-not $ref.ContainsKey($_.ToUpper())}).Count)"
```

Ini **[dugaan]**, bukan bukti layar mati: rujukan bisa datang dari `Flow`, `Activity`,
`DataPage`, atau `Harness` di luar tiga folder ini. Contoh yatim yang jelas masih terpakai:
`ScoringRiskForm1..4`, `ChooseRiskAddress_ResultList`, `crmDisplayStages`.
Namun kelompok berikut tetap layak ditanyakan ke bisnis sebelum ikut dimigrasikan:

| Kelompok yatim | Jumlah | Catatan |
| --- | ---: | --- |
| Varian `*_IsUW` yang pasangan non-UW-nya **juga** ada | banyak (mis. `CoverageList_IsUW`, `EmailSection_IsUW`, `PaymentCurrencyList_IsUW`, `OfferFacIn_NusaRe_IsUW`, `Periode_IsUW`) | rujukan mungkin lewat `Flow`/`Activity` |
| `*SummarySection` per lini (`FireSummarySection`, `GolfSummarySection`, `GrowingTreesSummarySection`, `PASummarySection`, `AllSummarySection`) | 5 | pola "ringkasan per COB" |
| `View*FacOut*` (`ViewObjectFireFacOut`, `ViewObjectPAFacOut`, `ViewPerilsFacOut`, …) | ±15 | tampilan Fac **Out**ward di dalam korpus Fac **In** |
| `DetailPolicyTreatyIn*NonProportional*`, `FormulaTreatyCapacityDesc*` | 5 | ranah Treaty, bukan Facultative |

**[pertanyaan terbuka]** Apakah layar `*FacOut*` dan `*TreatyIn*` termasuk lingkup migrasi
Facultative Inward, atau ikut terbawa karena satu ruleset?

---

## 4. Daftar FLOW ACTION — kandidat aksi/form UI

307 identitas flow action. **Seluruh 754 berkas** ber-`pyRuleFormType` = `Harness`
(tak satu pun `pyHTMLStream` kustom) — artinya semuanya dirender lewat section, sehingga
**pemetaan 1 aksi → 1 form React bersifat langsung**.

```powershell
$folders=@("NB FacIn","RNW Fac In","Endorsment Fac In")
foreach($f in $folders){ Get-ChildItem "D:\migrasi\RNM\$f\FlowAction" -File -Filter *.xml | ForEach-Object {
 $c=[IO.File]::ReadAllText($_.FullName)
 "{0,-18} {1,-45} {2,-42} {3}" -f $f,$_.BaseName,
   [regex]::Match($c,'<pySectionReference>([^<]*)</pySectionReference>').Groups[1].Value,
   [regex]::Match($c,'<pyLabel>([^<]*)</pyLabel>').Groups[1].Value } }
```

| Metrik | Nilai |
| --- | ---: |
| Berkas FlowAction | 754 |
| Identitas unik | **307** |
| Punya `pySectionReference` terisi | **749** (5 kosong) |
| Section berbeda yang dirujuk | **280** |
| Section referensi ada sebagai berkas `Section\*.xml` di folder yang sama | **728 / 749** |
| Identitas yang `pySectionReference` ≠ nama sendiri | **185 / 307** |
| Varian `*_IsUW` | **71** (68 punya pasangan non-UW) |

5 flow action tanpa `pySectionReference`: `AgentSourceBizDetails` (NB, RNW, EDM) dan
`StartScreenFlowAuto` (NB, RNW). **[pertanyaan terbuka]** apakah keduanya merender apa pun.

21 nama section yang dirujuk tetapi berkasnya tidak ada di folder yang sama:
`InputClauseFire_ViewDtl`, `InputDeductibleDtlAneka_GCNM`, `InputPackageDM`,
`pxUploadCSVResults` (bawaan Pega), `ViewAdditionalCoverage`.

### 4.1 Pengelompokan flow action — [dugaan] (dari nama)

```powershell
# $o = daftar 307 nama unik FlowAction; pola diterapkan berurutan, satu nama satu kelompok
```

| Tema | Jumlah |
| --- | ---: |
| Coverage / Klausula / Deductible | **81** |
| Obyek / Risiko / Okupasi / Lokasi | **64** |
| Spreading / Share / FacOut / Retro | **32** |
| Peserta / Person / Benefit / Plan | **21** |
| Premi / Komisi / Pembayaran / Formula | **20** |
| Pemilih (`Choose*` / `Select*`) | **13** |
| Dokumen / Email / Cetak / Notifikasi | **10** |
| Endorsement / Renewal | **7** |
| Akseptasi / UW / Decline / Scoring | **6** |
| Unggah CSV | **6** |
| Akumulasi | **2** |
| Lain-lain | **45** |

### 4.2 Aksi dengan label manusia — [terverifikasi] (`<pyLabel>`)

Label diambil apa adanya dari `<pyLabel>` berkas FlowAction. Hanya baris yang labelnya **berbeda**
dari nama rule yang ditampilkan (sisanya label = nama rule).

| FlowAction | Label di UI | Section yang dirender |
| --- | --- | --- |
| `InwardFacultative` | Inward Facultative | `InputInwardFacultative` |
| `InwardFacultative_IsUW` | Inward Facultative | `InputInwardFacultative_2` |
| `Endorsement_FlowAct` | Endorsement | `InputEndorsement` |
| `Endorsement_FlowAct_IsUW` | Endorsement untuk UW | `InputEndorsement_IsUW` |
| `Renewal_FlowAct` | Renewal | `InputRenewal` |
| `Renewal_FlowAct_IsUW` | Renewal unutk UW *(ejaan asli)* | `InputRenewal_IsUW` |
| `CoverageListFire_FlowAction` | Input Coverage | `PropertyItemListCoverage` |
| `CoverageListFire_FlowAction_IsUW` | Input Coverage_IsUW | `PropertyItemListCoverage_IsUW` |
| `FlowActCoverageListFire` | Input Coverage | `InputCoverageFire` |
| `CoverageFireGrid_IsUW` | Input Detail Jaminan_IsUW | `InputPerCoverageFire_IsUW` |
| `InputDtlDeductible_FacIn` | Input Own Risk | `InputDeductible_FacIn` |
| `InputDtlDeductibleFire` | Isi Resiko Sendiri | `InputDtlDeductibleFire` |
| `InputDtlSpreadingCoverage` | Isi Spreading | `InputDtlSpreadingCoverage` |
| `InputCoverageSpreading` | Input Coverage Spreading | `PropertyItemListCoverageSpreading` |
| `InputCoverageCommision` | Input Coverage Commision | `PropertyItemListCoverageCommision` |
| `InputDtlObjectAneka_FacIn` | Input Obyek Fac In | `InputDtlObjectAneka_FacIn` |
| `ObjItemGrid` | Input Obyek Item | `InputObjItem` |
| `OccupationGrid` | Input Okupasi | `InputOkupasi` |
| `OccupationAnekaGrid_FacIn` | Input Okupasi Aneka | `InputOkupasiAneka_FacIn` |
| `OccupationAnekaGridGCNM` | Input Okupasi Aneka GCNM | `InputOkupasiAneka_GCNM` |
| `InputAnekaPolicySchedule` | Schedule Polis | `InputDtlAnekaPolicySchedule` |
| `ObjectOtherSchedule` | Schedule Lain | `InputOtherSchedule` |
| `InputDtlOtherMaintenanceAneka` | Aneka Maintenance | `InputDtlMaintenance` |
| `InputDtlOtherObjectAneka*` (3 varian) | General Information | `InputDtlOtherObjectAneka*` |
| `CauseOfLossGrid` | Input Pengalaman Kerugian | `InputLossExperiance` |
| `InputCauseOfLoss_FacIn` | Input Cause Of Loss | `InputCauseOfLoss_FacIn` |
| `InputLossRecord` | Input Loss Record | `InputLossRecord_Sec` |
| `InputFEA` | Input Fire Extinguisher Availability | `InputFEA` |
| `FEAGrid` | Input Data FEA | `InputFEA` |
| `InputClause_ViewDtl` | Lihat Detail Klausula | `InputClause_ViewDtl` |
| `InputClauseFire_ViewDtl` | Isi Klasula *(ejaan asli)* | `InputClauseFire_ViewDtl` |
| `InpBenefit` | Pilih Benefit | `CreateListBenefit` |
| `ChooseAccumulation` | Pilih Akumulasi | `ChooseAccumulation` |
| `ModalAccumulation_FacIn` | Accumulation Type | `InputAccumulatedType_FacIn` |
| `ChooseRiskAddress` | Search Risk Address | `ChooseRiskAddress` |
| `CopyCoverageFrom` | Copy Coverage From | `CopyCoverageFrom` |
| `TemplateCoverage_FacIn` | Copy Premi from PKS | `TemplateCoverage_FacIn` |
| `PolicyTreatyInDeclineConfirm` | Confirm Decline NB | `PolicyTreatyInDeclineConfirm` |
| `RejectNotificationFlow` | Reject Notification | `RejectNotificationSection` |
| `PrintRISlip_FlowAction` | Print RI Slip | `PrintRISlip` |
| `ReasViewAttachment` | Reas View Attachment | `ReasViewAttachment` |
| `TableInwardScale` / `TableOfLimit` | Table Inward Scale / Table Of Limit | idem |
| `ViewIndemnity_LA` | View Indemnity Table | `ViewIndemnity_Section` |
| `UploadCSV_DataPolis` | Data Polis - Upload CSV | `pxUploadCSVResults` |
| `UploadCSV_Participant` | Participant - Upload CSV Records | `pxUploadCSVResults` |
| `UploadCSV_PolicyMember` | Member - Upload CSV Records | `pxUploadCSVResults` |
| `UploadCSV_Vehicle` | Vehicle - Upload CSV Records | `pxUploadCSVResults` |

### 4.3 Satu section dipakai oleh beberapa aksi — [terverifikasi]

21 section dirender oleh lebih dari satu flow action. Yang paling ekstrem:

| Section | Dipakai oleh |
| --- | --- |
| `pxUploadCSVResults` (bawaan Pega) | **6** aksi: `UploadCSV_Aneka`, `UploadCSV_DataPolis`, `UploadCSV_Participant`, `UploadCSV_PolicyMember`, `UploadCSV_Vehicle`, `UploadCSVEDM_BenefitLimit` |
| `ShowCoverageFacOut` | 3 aksi |
| `InputCoverageFire`, `InputCoverageFire_IsUW`, `InputProvince`, `InputCity`, `InputDistrict`, `ListPayment_SC`, `InputFEA`, `pyAttachmentScreen`, `InputCoverageLife_FacIn(+_IsUW)`, `InputCoverageAneka_FacIn(+_IsUW)`, `InputSpreadingLife_FacIn(+_IsUW)`, `InputCommissionLife_FacIn`, `InputCommissionPA_FacIn`, `InputClause_ViewDtl`, `ViewCoverageFire`, `ShowCoverageFacOut_IsUW` | 2 aksi masing-masing |

Pasangan yang menarik untuk React: `InputCity` ← `InputCity` + `InputKabupaten`,
`InputDistrict` ← `InputDistrict` + `InputDistrik`, `InputProvince` ← `InputProvince` +
`InputKotaProvince`, `ViewPaymentPremi` + `ViewPaymentRetro` → `ListPayment_SC`.
Satu komponen, beberapa pintu masuk.

---

## 5. Tumpang tindih & percabangan antar-siklus

### 5.1 Dua ukuran, karena `Get-FileHash` mentah menyesatkan

| Ukuran | Definisi | Tujuan |
| --- | --- | --- |
| **A — hash mentah** | `Get-FileHash -Algorithm SHA256` apa adanya | mendeteksi ekspor yang benar-benar byte-identik |
| **B — hash ternormalisasi** | buang baris metadata ekspor + indeks, samarkan timestamp, urutkan baris secara **ordinal**, lalu hash | membandingkan **isi desain**, mengabaikan urutan serialisasi XML dan jejak instance |

Normalisasi B membuang tag: `pxCommitDateTime`, `pxUpdateDateTime`, `pxSaveDateTime`,
`pxCreateDateTime`, `pxOriginalCreateDateTime`, `pxMoveImportDateTime`, `pyRuleFormStatusTime`,
`pyShowJavaWindowName`, `pxHostId`, `pzIndexCount`, `pxUpdateOpName`, `pxCreateOpName`,
`pxMoveImportOperName`, `pxOriginalCreateOpName`, `pySPRuleSetName`, `pzInsKey`,
`pzIndexOwnerKey`, `pxInstanceLockedBy`, `pxInstanceLockedCreateDateTime`, `pxCommitSystemID`,
`pxUpdateSystemID`, `pxOriginalCreateSystemID`, `pyJavaStream`, `pyHTMLStream`; membuang baris
`<rowdata REPEATINGINDEX="…Reference">n</rowdata>`; dan mengganti `YYYYMMDDTHHMMSS.mmm GMT` → `TS`.

> Empat tag terakhir memuat **nama orang** (`pxCreateOpName`, `pxUpdateOpName`,
> `pxMoveImportOperName`, `pxOriginalCreateOpName`). Nilainya tidak disalin ke mana pun di
> dokumen ini; tag hanya dibuang dari perbandingan.

**Ukuran A:**

```powershell
$all=@{}
foreach($f in @("NB FacIn","RNW Fac In","Endorsment Fac In")){ foreach($t in @("Harness","Section","FlowAction")){
  Get-ChildItem "D:\migrasi\RNM\$f\$t" -File -Filter *.xml | ForEach-Object {
    $k="$t|"+$_.BaseName; if(-not $all.ContainsKey($k)){$all[$k]=@{}}
    $all[$k][$f]=(Get-FileHash $_.FullName -Algorithm SHA256).Hash } } }
foreach($t in @("Harness","Section","FlowAction")){
  $s=$all.Keys|Where-Object{$_ -like "$t|*"}
  $tiga=$s|Where-Object{$all[$_].Count -eq 3}
  "{0,-12} di-3={1,-4} identik-di-3={2}" -f $t,$tiga.Count,($tiga|Where-Object{(@($all[$_].Values|Sort-Object -Unique)).Count -eq 1}).Count }
```

Hasil ukuran A: **0** layar byte-identik di ketiga folder (Harness 0/33, Section 0/337,
FlowAction 0/205). **Ini artefak ekspor, bukan temuan desain** — lihat §0.1(b).

**Ukuran B:**

```powershell
$volatile='pxCommitDateTime|pxUpdateDateTime|pxSaveDateTime|pxCreateDateTime|pxOriginalCreateDateTime|pxMoveImportDateTime|pyRuleFormStatusTime|pyShowJavaWindowName|pxHostId|pzIndexCount|pxUpdateOpName|pxCreateOpName|pxMoveImportOperName|pxOriginalCreateOpName|pySPRuleSetName|pzInsKey|pzIndexOwnerKey|pxInstanceLockedBy|pxInstanceLockedCreateDateTime|pxCommitSystemID|pxUpdateSystemID|pxOriginalCreateSystemID|pyJavaStream|pyHTMLStream'
function OrdHash([string]$path){
  $md5=[System.Security.Cryptography.MD5]::Create()
  $keep=New-Object System.Collections.Generic.List[string]
  foreach($l in [IO.File]::ReadAllLines($path)){
    if($l -match "<($volatile)[>/]"){continue}
    if($l -match '^<rowdata REPEATINGINDEX="[A-Za-z]+">\d+</rowdata>$'){continue}
    [void]$keep.Add(($l -replace '\d{8}T\d{6}\.\d{3} GMT','TS')) }
  $arr=$keep.ToArray(); [Array]::Sort($arr,[StringComparer]::Ordinal)
  ([BitConverter]::ToString($md5.ComputeHash([Text.Encoding]::UTF8.GetBytes(($arr -join "`n"))))) -replace '-','' }
$all=@{}
foreach($f in @("NB FacIn","RNW Fac In","Endorsment Fac In")){ foreach($t in @("Harness","Section","FlowAction")){
  Get-ChildItem "D:\migrasi\RNM\$f\$t" -File -Filter *.xml | ForEach-Object {
    $k="$t|"+$_.BaseName; if(-not $all.ContainsKey($k)){$all[$k]=@{}}; $all[$k][$f]=(OrdHash $_.FullName) } } }
foreach($t in @("Harness","Section","FlowAction")){
  $tiga=$all.Keys|Where-Object{$_ -like "$t|*" -and $all[$_].Count -eq 3}
  $nb='NB FacIn';$rnw='RNW Fac In';$edm='Endorsment Fac In'
  "{0,-12} di-3={1,-4} NB=RNW:{2,-4} NB=EDM:{3,-4} identik-3:{4,-4} bercabang:{5}" -f $t,$tiga.Count,
    ($tiga|Where-Object{$all[$_][$nb] -ceq $all[$_][$rnw]}).Count,
    ($tiga|Where-Object{$all[$_][$nb] -ceq $all[$_][$edm]}).Count,
    ($tiga|Where-Object{$all[$_][$nb] -ceq $all[$_][$rnw] -and $all[$_][$nb] -ceq $all[$_][$edm]}).Count,
    ($tiga|Where-Object{-not($all[$_][$nb] -ceq $all[$_][$rnw] -and $all[$_][$nb] -ceq $all[$_][$edm])}).Count }
```

> Catatan implementasi: pengurutan **wajib** `[StringComparer]::Ordinal`. `Sort-Object` bawaan
> PowerShell membandingkan tanpa peduli huruf besar/kecil dan tidak stabil untuk baris yang hanya
> berbeda kapitalisasi; dengan pengurutan itu, Section identik pun terbaca bercabang (0 identik).
> Demikian pula `Compare-Object` **tidak** dapat dipakai sebagai verifikasi: ia mengabaikan
> jumlah duplikat dan case, sehingga melaporkan `diff=0` untuk berkas yang isinya memang berbeda.

### 5.2 Hasil — temuan utama dokumen ini

| Tipe | Ada di 3 folder | NB = RNW | NB = EDM | **Identik di 3** | **Bercabang** |
| --- | ---: | ---: | ---: | ---: | ---: |
| `Harness` | 33 | **33 (100 %)** | 27 | **27** | **6** |
| `Section` | 337 | **337 (100 %)** | 261 | **261** | **76** |
| `FlowAction` | 205 | **205 (100 %)** | 161 | **161** | **44** |
| **Jumlah** | **575** | **575 (100 %)** | 449 | **449 (78,1 %)** | **126 (21,9 %)** |

| Tipe | Ada di 2 folder | identik | bercabang | Hanya 1 folder |
| --- | ---: | ---: | ---: | ---: |
| `Harness` | 8 | 7 | 1 | 4 |
| `Section` | 69 | 64 | 5 | 114 |
| `FlowAction` | 37 | 35 | 2 | 65 |

**Temuan #1 — [terverifikasi]: NB dan RNW tidak berbeda sama sekali di lapisan UI.**
Untuk **575 dari 575** identitas yang ada di ketiga folder, salinan NB dan salinan RNW identik
sebagai multiset baris ternormalisasi. Tidak ada satu pun pengecualian.

Konsekuensi untuk target React: **New Business dan Renewal berbagi satu himpunan layar yang sama
persis.** Perbedaan siklus NB vs RNW **tidak** hidup di lapisan UI; ia harus dicari di `Flow`,
`When`, `Activity` — di luar cakupan dokumen ini. Membangun dua pohon komponen terpisah untuk
NB dan RNW tidak didukung korpus.

**Batas metode — [dugaan]:** ukuran B mengabaikan **urutan** elemen XML. Bila dua salinan memuat
himpunan baris yang sama tetapi susunan kolom/urutan field berbeda, ukuran B akan menyebutnya
identik. Verifikasi urutan memerlukan pembandingan pohon XML per-node, yang belum dikerjakan.

**Temuan #2 — [terverifikasi]: seluruh percabangan ada di EDM.**
126 layar bercabang; dalam **setiap** kasus, polanya NB = RNW ≠ EDM.

### 5.3 Daftar HARNESS bercabang (ada di 3 folder, EDM berbeda) — 6

| Harness | Sifat perbedaan yang terbaca |
| --- | --- |
| `CedingCompany` | kelas berpindah: NB/RNW `ASM-FW-GISFW-Data-OfferTreatyIn` → EDM `ASM-FW-GISFW-Data-Quotation` **[terverifikasi]** |
| `ChooseClauseFire` | kelas NB/RNW `ASM-FW-GISFW-Work` → EDM `ASM-FW-GISFW-Data-Clause`; label NB/RNW "Pilih Klausula Fire" → EDM "ChooseClauseFire" **[terverifikasi]** |
| `HistoricalSurveyReportUW` | kelas NB/RNW `…Data-PolicyTreatyIn` → EDM `…Data-OfferFacIn` **[terverifikasi]** |
| `ViewCoverageFacOut` | kelas NB/RNW `…Data-PropertyItem` → EDM `…Data-Aneka` **[terverifikasi]** |
| `ViewCoverageFacOut_isUW` | kelas NB/RNW `…Data-PropertyItem` → EDM `…Data-Aneka`; **section berubah nama**: `ViewCoverageFacOutShow_isUW` → `ViewCoverageFacOutShow_Is_UW` **[terverifikasi]** |
| `TotalAccumulationDtl` | kelas sama (`…Data-Coverage`), perbedaan ada di isi layout — **belum terverifikasi** jenisnya |

```powershell
foreach($f in @("NB FacIn","RNW Fac In","Endorsment Fac In")){
 foreach($n in @('CedingCompany','ChooseClauseFire','HistoricalSurveyReportUW','TotalAccumulationDtl','ViewCoverageFacOut','ViewCoverageFacOut_isUW')){
  $c=[IO.File]::ReadAllText("D:\migrasi\RNM\$f\Harness\$n.xml")
  "{0,-18} {1,-26} {2,-34} {3}" -f $f,$n,
    [regex]::Match($c,'<pxInsName>([^!<]*)!').Groups[1].Value,
    [regex]::Match($c,'<pyLabelOld>([^<]*)</pyLabelOld>').Groups[1].Value } }
```

Tiga harness hanya di satu folder: `Medical_Harnes`, `SFAPortalEndorsement`, `ViewOldData`
(semua EDM); `SFAPortalOpportunities` (NB). `ChooseInsured` hanya RNW+EDM dan **bercabang**;
`InputCauseOfDecline`, `SelectCoverage`, `ShowPolis`, `SOB`, `TabbedScreenFlow7`,
`ViewCSVResult_B2BHost`, `ViewFollowingNB` hanya NB+RNW dan identik.

### 5.4 Daftar SECTION bercabang (ada di 3 folder, EDM berbeda) — 76

```
additemButton                     InputDtlObjectLocation_FacIn_IsUW   PaymentCurrencyList
AttachmentGridReas                InputDtlObjFire_IsUW                PaymentCurrencyListLife
CallDualScoringRisk               InputDtlOtherObjectAneka_FacIn      Property_IsUW
CauseOfLoss_FacIn                 InputDtlOtherSchedule               PropertyItemListCoverage_IsUW
CedingCoHierarki                  InputDtlParticipantLife_FacIn_IsUW  ScoringRisk
ChooseClassofContraction          InputDtlParticipantPA_FacIn_IsUW    SelClauseList
ChooseClauseFire                  InputDtlParticipantTravel_IsUW      ShowCoverageFacOut
ChooseCoverage                    InputDtlPayment_FacIn               ShowCoverageFacOut_IsUW
ChooseObject_Ship                 InputDtlPayment_FacInLife           ShowPolis_sc
ChooseOccupation                  InputFacOffer                       TotalAccumulationDtl
Correspondence                    InputHistoricalSurveyReportDtl      UploadMember
CoverageList                      InputInwardFacultativeDtl_IsUW      ViewCoverage
CoverageSpreadingList             InputMarketingOfficer               ViewCoverageFacOutShow
CreateListClause                  InputObjectDtlFire_IsUW             ViewCoverageFire
EmailSection                      InputObjectSubContract_FacIn_IsUW   ViewCoveragePA
HistoricalSurveyReportDtlUW       InputOtherObjectAneka_FacIn         ViewCoverageTravel
InputAccumulationCov              InputOtherObjectAneka_FacIn_ISUW    ViewDeliveryAddressFacOut
InputClause                       InputOtherSchedule                  ViewDetailPayment_SC
InputCommissionLife_FacIn         InputPerson                         ViewGeneralPolisFacOut
InputCoverageFire                 InputSpreadingLife_FacIn_IsUW       ViewInwardFacultativeDtl_IsUW
InputCoverageFire_IsUW            ObjectList                          ViewObjectMarineCargoFacOut
InputDeductible_FacIn_IsUW        OfferStatusFacOut                   ViewObjectPAFacOut
InputDeductible_GCNM              InputDtlAnekaPolicySchedule         ViewPerilsFacOut
InputDeductibleHull_GCNM          InputDtlCoverage_FacIn              ViewVehicleGridFacOut
InputDtlCoverage_FacIn_IsUW       InputDtlObject_FacIn                InputDtlObject_FacIn_IsUW
InputDtlObjectLocation_FacIn
```

```powershell
# $all dari skrip ukuran B di §5.1
$nb='NB FacIn';$edm='Endorsment Fac In'
$all.Keys | Where-Object{ $_ -like 'Section|*' -and $all[$_].Count -eq 3 -and -not ($all[$_][$nb] -ceq $all[$_][$edm]) } |
  ForEach-Object{ $_.Split('|')[1] } | Sort-Object
```

**Yang paling mahal untuk target React** — section bercabang yang sekaligus komponen bersama
peringkat atas (§3.1): `PaymentCurrencyList` (83 rujukan), `CoverageSpreadingList` (60),
`InputDtlObject_FacIn` (41), `InputCoverageFire` (19), `EmailSection` (22),
`InputDtlCoverage_FacIn` (23), `ScoringRisk` (12). Komponen bersama yang **berbeda isi** antar
siklus adalah kombinasi terburuk: satu komponen React tidak bisa melayani ketiganya tanpa cabang
internal, dan cabangnya belum dipetakan.

Section bercabang yang hanya ada di 2 folder — 5: `FormulaTreatyCapacityDesc_IsUW`,
`OfferFacIn_NusaRe`, `PeriodeEndorsement`, `InputCoverageAneka_FacIn` (NB+EDM);
`ChooseInsuredDtl` (RNW+EDM).

### 5.5 Daftar FLOW ACTION bercabang (3 folder, EDM berbeda) — 44

```
CedingCedant                      InputCoverageCommision            ObjectGridLocation_FacIn_IsUW
ChooseClassofContraction          InputDistrict                     ObjectOtherSchedule
ChooseCoverage                    InputDtlCargo_FacIn_IsUW          PersonCoverage
ChooseObject_Ship                 InputDtlCoverage_FacIn            PersonGridDM
ChooseOccupation                  InputDtlGoods_IsUW                PersonGridLife_FacIn_IsUW
ChooseSubContract_FacIn           InputDtlOtherObjectAneka_FacIn_ISUW  PersonGridPA_FacIn_IsUW
ChooseZipCode                     InputDtlPaymentFlow_FacIn         PersonGridTravel_IsUW
CoverageListFire_FlowAction_IsUW  InputDtlPaymentFlow_FacIn_IsUW    Property_FlowAction_IsUW
InpClause                         InputDtlTrading_IsUW              SelectionPlanRO
InputAnekaPolicySchedule          InputObjFireGrid_IsUW             ShowCoveragePAFacOutIsUW_FlowAction
InputClause_ViewArg               InputCommissionLife_FacIn         UploadCSV_Aneka
InputClause_ViewDtl               InputCommissionMBU_FacIn          UploadCSV_Participant
InputCommissionMC_FacIn           InputCommissionPA_FacIn           UploadCSV_PolicyMember
InputCommissionTravel_FacIn       ViewObjectDetailOccupation        UploadCSV_Vehicle
ViewPaymentPremi                  VehicleGrid_FacIn_IsUW
```

```powershell
# $all dari skrip ukuran B di §5.1
$nb='NB FacIn';$edm='Endorsment Fac In'
$all.Keys | Where-Object{ $_ -like 'FlowAction|*' -and $all[$_].Count -eq 3 -and -not ($all[$_][$nb] -ceq $all[$_][$edm]) } |
  ForEach-Object{ $_.Split('|')[1] } | Sort-Object
```

Flow action bercabang di 2 folder — 2: `LimitTreaty`, `InputDtlFormulaFlow_FacIn` (NB+EDM).

### 5.6 Layar khas per siklus (hanya 1 folder)

| Folder | Harness | Section | FlowAction |
| --- | ---: | ---: | ---: |
| NB | 1 | 27 | 8 |
| RNW | 0 | 8 | 3 |
| EDM | 3 | 79 | 54 |
| **Jumlah** | **4** | **114** | **65** |

```powershell
# $all dari skrip ukuran B di §5.1
foreach($t in @('Harness','Section','FlowAction')){
 $s=@($all.Keys|Where-Object{$_ -like "$t|*" -and $all[$_].Count -eq 1})
 "{0,-12} total1={1,-4} NB={2,-4} RNW={3,-4} EDM={4}" -f $t,$s.Count,
   @($s|Where-Object{$all[$_].ContainsKey('NB FacIn')}).Count,
   @($s|Where-Object{$all[$_].ContainsKey('RNW Fac In')}).Count,
   @($s|Where-Object{$all[$_].ContainsKey('Endorsment Fac In')}).Count }
```

**RNW menyumbang paling sedikit layar khas** (8 section, 3 flow action): `InputRenewal`,
`InputRenewal_IsUW`, `InputRenewalDtl`, `InputRenewalDtl_IsUW`, `PeriodeRenewal`,
`PeriodeRenewal_IsUW`, `EditMarketing`, `SFAPortal_Renewal`; aksi `Renewal_FlowAct`,
`Renewal_FlowAct_IsUW`, `EditMarketing`. Menggabungkan dengan temuan #1: **Renewal = NB + satu
lapis tipis form periode/renewal.**

**EDM menyumbang paling banyak** (79 section, 54 flow action). Kelompok yang terbaca dari nama
**[dugaan]**: rumpun *medical / benefit / plan* (`Medical_Sec`, `InputMedical1_Sec`,
`InputMedical2_Sec`, `HistoryofDisease_Sec`, `BenefitClause*`, `SelectPlanList*`,
`CreateListPlan`, `PlanContainsAll`, `InpMedPackage`, `Pkg_Detail`, `ShowDetPckge`), rumpun
*endorsement* (`EdmType1/2/3`, `InputEndorsement_IsUW`, `PeriodeEndorsement_IsUW`,
`AcceptNotificationEdm`, `EditMarketingEndorsment`, `extensiondtl(_IsUW)`, `FinalPolicy`,
`ClientApproval`), dan rumpun *coverage per COB* (`InputCoverageLife_FacIn`,
`InputCoverageMBU_FacIn`, `InputCoveragePA_FacIn`, `InputCoverageTravel`,
`InputCoverageGridMBU_FacIn`, masing-masing dengan varian `_IsUW`).

**NB** menyumbang 27 section khas, didominasi rumpun **Treaty** (`DetailPolicyTreatyIn*`,
`GeneralPolicyTreatyIn`, `FormulaTreatyCapacityDesc`, `PolicyTreatyInDeclineConfirm`,
`DeptHeadTreatyIn_UW`), `SFAPortal_Opportunities*`, `Installment*`, `Warranty*`, `*Suggest`.

---

## 6. Komponen berulang yang layak jadi komponen React bersama

Bukti di bawah berasal dari **jumlah rujukan** (§3.1) dan **pengikatan data grid**
(`<pyPageListProperty>`), bukan dari kesan atas nama.

```powershell
$rows=New-Object System.Collections.ArrayList
foreach($f in @("NB FacIn","RNW Fac In","Endorsment Fac In")){
 Get-ChildItem "D:\migrasi\RNM\$f\Section" -File -Filter *.xml | ForEach-Object {
  $c=[IO.File]::ReadAllText($_.FullName)
  foreach($p in ([regex]::Matches($c,'<pyPageListProperty>([^<]+)</pyPageListProperty>')|ForEach-Object{$_.Groups[1].Value}|Sort-Object -Unique)){
    [void]$rows.Add([PSCustomObject]@{Folder=$f;Section=$_.BaseName;Prop=$p}) } } }
"pasangan section-properti=$($rows.Count) ; properti unik=$((($rows.Prop)|Sort-Object -Unique).Count)"
$rows | Group-Object Prop | Sort-Object Count -Descending | Select-Object -First 30 Count,Name
```

**1.381** pasangan (folder × section × properti-list), **219** properti list berbeda.

Kolom **pasangan** menghitung folder × section (satu section yang ada di tiga folder terhitung
tiga kali); kolom **section unik** menggabungkan lintas folder. Keduanya dilaporkan karena
keduanya menjawab pertanyaan berbeda: beban migrasi vs jumlah komponen React.

```powershell
$rows | Select-Object Section,Prop -Unique | Group-Object Prop |
  Sort-Object Count -Descending | Select-Object -First 25 Count,Name
```

| Properti PageList (sumber baris grid) | Pasangan | Section unik |
| --- | ---: | ---: |
| `.CoverageList` | **101** | **43** |
| `.LocationList` | 45 | 15 |
| `.PersonList` | 42 | 14 |
| `.VehicleList` | 41 | 14 |
| `.AnekaList` | 35 | 13 |
| `.DeductibleList` | 35 | **17** |
| `.ASMCoverage` | 33 | 15 |
| `.Property.PropertyItemList` | 33 | 11 |
| `.ASMHeir` | 30 | 11 |
| `.Property.RiskLocation.OccupationList` | 30 | 11 |
| `.CargoList` | 27 | 9 |
| `.PropertyList` | 24 | 8 |
| `.Property.RiskLocation.AnekaList` | 24 | 9 |
| `.OccupationList` | 23 | 8 |
| `TotalSpreadAll.pxResults` | 21 | 9 |
| `.OfferFacIn.CurrencyList` + `pyWorkPage.OfferFacIn.CurrencyList` | 18 + 17 | 8 + 7 (**15** section unik) |
| `.ClauseList` | 18 | 8 |
| `SpreadingList.pxResults` | 18 | 8 |
| `TempViewSuggest.pxResults` | 17 | 8 |
| `.OfferFacIn.CedingCedantList` | 15 | 7 |
| `.AdditionalCoverage` | – | 6 |
| `.OfferFacIn.LocationList` / `TotalSpreadingCurrency.pxResults` | 14 | 6 |
| `.LayerList` | 14 | 5 |

### 6.1 Kandidat komponen bersama, dengan bukti isi

Enam berkas berikut **dibaca isinya** (bukan hanya namanya) untuk dokumen ini, lewat ekstraksi
seluruh `<pyPageListProperty>` dan seluruh FieldValue `pyCaption …` yang dirujuknya:

- `D:\migrasi\RNM\NB FacIn\Section\PaymentCurrencyList.xml` (303.728 byte)
- `D:\migrasi\RNM\NB FacIn\Section\CoverageSpreadingList.xml` (743.451 byte)
- `D:\migrasi\RNM\NB FacIn\Section\CoverageCommisionList.xml` (737.820 byte)
- `D:\migrasi\RNM\NB FacIn\Section\InputRiskAddress.xml` (305.967 byte)
- `D:\migrasi\RNM\NB FacIn\Section\OfferFacIn_NusaRe.xml` (407.982 byte)
- `D:\migrasi\RNM\NB FacIn\Section\EmailSection.xml` (650.060 byte)

Ditambah dua berkas yang dibaca untuk membuktikan struktur tag (§0.1):
`D:\migrasi\RNM\NB FacIn\FlowAction\CauseOfLossGrid.xml`,
`D:\migrasi\RNM\NB FacIn\Harness\CedingCedant.xml`.

```powershell
foreach($n in @('PaymentCurrencyList','CoverageSpreadingList','CoverageCommisionList',
                'InputRiskAddress','OfferFacIn_NusaRe','EmailSection')){
 $c=[IO.File]::ReadAllText("D:\migrasi\RNM\NB FacIn\Section\$n.xml")
 "### $n"
 "  binding : " + (([regex]::Matches($c,'<pyPageListProperty>([^<]+)</pyPageListProperty>')|
                    ForEach-Object{$_.Groups[1].Value}|Sort-Object -Unique) -join ', ')
 "  kolom   : " + (([regex]::Matches($c,'<pyRuleName>pyCaption ([^<]*)</pyRuleName>')|
                    ForEach-Object{$_.Groups[1].Value}|Sort-Object -Unique) -join ' | ') }
```

| Komponen | Pengikatan data | Kolom/field yang terbaca | Rujukan |
| --- | --- | --- | ---: |
| **`PaymentCurrencyList`** — daftar mata uang & premi | `pyWorkPage.OfferFacIn.CurrencyList`, kelas baris `ASM-FW-GISFW-Data-OfferFacIn-Currency` | Currency · Installment · Premium RNM · Actual Net Premium · Old Net Premium · New Net Premium · Total Payment | **83** |
| **`CoverageSpreadingList`** — daftar obyek per jenis, untuk spreading | `.PersonList`, `.LocationList`, `.VehicleList`, `.CargoList` (4 grid dalam 1 section) | No · Object Name · Location · Name · Age · Date of Birth · Gender · ID Card · Height (Cm) · Weight (Kg) · Left-handed · Partisipant Status *(ejaan asli)* · Brand · Model · Type · Chassis Number · Engine Number · License Plate · Goods · Packing · Conveyance · Trading | **60** |
| **`CoverageCommisionList`** — kembaran untuk komisi | `.PersonList`, `.LocationList`, `.VehicleList`, `.CargoList` — **identik dengan di atas** | sama, kecuali "Participant Status" (ejaan benar) | **53** |
| **`OfferFacIn_NusaRe`** — panel share/limit Nusa Re | `pyWorkPage.OfferFacIn.CurrencyList` | Currency · Nusa Re · %ASM Share · %Share in TSI · ASM TSI Share · ASM Gross Premium · ASM Net Premium · %R/I Comm · %Rate · Treaty Capacity · %Max Treaty Capacity · Max Treaty Capacity · %Additional Capital · TSI Additional Capital · TSI Liability · TSI Top Risk · Status · Remarks | **31** |
| **`InputRiskAddress`** — alamat risiko | (tanpa grid) | Title · Type · Address · City · District · Province · Country · Zip Code · Code · Territory · Status Trans Pusat · Button | **18** |
| **`EmailSection`** — panel korespondensi/approval | `TempViewSuggest.pxResults` | Email Type · PIC · Date · Status · Approval · Banding to · Comment · Choose · Button | **22** |

**`CoverageSpreadingList` dan `CoverageCommisionList` adalah duplikasi yang nyaris sempurna:**
pengikatan grid identik (`.PersonList`, `.LocationList`, `.VehicleList`, `.CargoList`), 23 caption
yang sama, dan satu-satunya beda tekstual yang terbaca adalah ejaan "Partisipant"/"Participant".
**Usulan:** satu komponen React `<ObjectListGrid variant="spreading" | "commission">`.

**`OfferFacIn_NusaRe` menegaskan aturan uang di CLAUDE.md** — kolom `%ASM Share`, `ASM TSI Share`,
`TSI Liability`, `TSI Top Risk` dan `Max Treaty Capacity` semuanya hidup di satu baris yang
dikunci `Currency`. Komponen React-nya **wajib** membawa mata uang bersama nilainya.

### 6.2 Duplikasi UW / non-UW — [terverifikasi]

| Tipe | Identitas unik | Varian `*_IsUW` | Punya pasangan non-UW |
| --- | ---: | ---: | ---: |
| `Section` | 520 | **101** | **97** |
| `FlowAction` | 307 | **71** | **68** |
| `Harness` | 45 | 1 | 1 |

```powershell
foreach($t in @('Section','FlowAction','Harness')){
 $n=(Get-ChildItem "D:\migrasi\RNM\*\$t" -File -Filter *.xml).BaseName | Sort-Object -Unique
 $uw=$n|Where-Object{$_ -match '(?i)_is_?uw$'}
 $p=0; foreach($x in $uw){ $b=$x -replace '(?i)_is_?uw$',''; if($n|Where-Object{$_ -ieq $b}){$p++} }
 "{0,-12} total={1,-4} varian_IsUW={2,-4} punya-pasangan={3}" -f $t,$n.Count,$uw.Count,$p }
```

**97 dari 520 section (18,7 %) adalah salinan UW dari section non-UW.** Kemiripan diukur dengan
Jaccard atas himpunan baris ternormalisasi:

```powershell
$volatile='pxCommitDateTime|pxUpdateDateTime|pxSaveDateTime|pxCreateDateTime|pxOriginalCreateDateTime|pxMoveImportDateTime|pyRuleFormStatusTime|pyShowJavaWindowName|pxHostId|pzIndexCount|pxUpdateOpName|pxCreateOpName|pxMoveImportOperName|pxOriginalCreateOpName|pySPRuleSetName|pzInsKey|pzIndexOwnerKey|pxInstanceLockedBy|pxInstanceLockedCreateDateTime|pxCommitSystemID|pxUpdateSystemID|pxOriginalCreateSystemID|pyJavaStream|pyHTMLStream|pyAutomationID|pxInsName'
function NL([string]$p){ foreach($x in [IO.File]::ReadAllLines($p)){
  if($x -match "<($volatile)[>/]"){continue}; $x -replace '\d{8}T\d{6}\.\d{3} GMT','TS' } }
foreach($n in @('CoverageItem','EmailSection','InputCoverageFire','InputDtlObject_FacIn',
                'InputInwardFacultativeDtl','PeriodeRenewal','SpreadingItem')){
 foreach($f in @('NB FacIn','RNW Fac In','Endorsment Fac In')){
  $a="D:\migrasi\RNM\$f\Section\$n.xml"; $b="D:\migrasi\RNM\$f\Section\${n}_IsUW.xml"
  if((Test-Path $a) -and (Test-Path $b)){
   $sa=[System.Collections.Generic.HashSet[string]](NL $a); $sb=[System.Collections.Generic.HashSet[string]](NL $b)
   $i=[System.Collections.Generic.HashSet[string]]::new($sa); $i.IntersectWith($sb)
   $u=[System.Collections.Generic.HashSet[string]]::new($sa); $u.UnionWith($sb)
   "{0,-28} {1,-18} Jaccard={2:P1}" -f $n,$f,($i.Count/$u.Count); break } } }
```

| Pasangan | Folder sampel | Jaccard |
| --- | --- | ---: |
| `InputDtlObject_FacIn` vs `…_IsUW` | NB | **79,9 %** |
| `CoverageItem` vs `…_IsUW` | NB | **77,7 %** |
| `SpreadingItem` vs `…_IsUW` | NB | **75,1 %** |
| `PeriodeRenewal` vs `…_IsUW` | RNW | **73,4 %** |
| `InputCoverageFire` vs `…_IsUW` | NB | **72,2 %** |
| `InputInwardFacultativeDtl` vs `…_IsUW` | NB | **69,2 %** |
| `EmailSection` vs `…_IsUW` | NB | **64,6 %** |

**Usulan:** satu komponen React dengan prop `mode: 'input' | 'uw'`, bukan dua komponen. Sisa 20–35 %
perbedaan **belum terverifikasi** isinya dan harus dipetakan field-per-field sebelum digabung.

---

## 7. Pola UI yang sulit dipindahkan ke React

### 7.1 Grid berulang — [terverifikasi]

```powershell
$n=0;$f2=0;$tot=0
foreach($f in @("NB FacIn","RNW Fac In","Endorsment Fac In")){
 Get-ChildItem "D:\migrasi\RNM\$f\Section" -File -Filter *.xml | ForEach-Object {
  $c=[IO.File]::ReadAllText($_.FullName)
  $m=[regex]::Matches($c,'<pyRepeatDirection>RepeatGrid</pyRepeatDirection>').Count
  $tot+=$m; if($m -gt 0){$n++} } }
"berkas Section ber-RepeatGrid = $n ; total grid = $tot"
```

| Metrik | Nilai |
| --- | ---: |
| Berkas Section yang memuat ≥1 `RepeatGrid` | **822 / 1.263 (65,1 %)** |
| Total layout `RepeatGrid` | **1.810** |
| Berkas dengan baris dapat di-expand (`pyExpandable=true`) | total 1.280 kemunculan |
| Berkas dengan deferred load (`pyLoadDeferred=true`) | hanya 50 kemunculan |

Section dengan grid terbanyak:

| Section | RepeatGrid | Sub-section disisipkan | Ukuran |
| --- | ---: | ---: | ---: |
| `InputDtlObject_FacIn` | **22** | 52 | 3,59 MB |
| `InputDtlObject_FacIn_IsUW` | 22 | 21 | 3,15 MB |
| `DetailPolicyTreatyInNonProportional` (NB saja) | 17 | 0 | 1,66 MB |
| `InputDtlCoverage_FacIn_IsUW` | 13 | 0 | 1,69 MB |
| `ViewInwardFacultativeDtl_IsUW` | 12 | 22 | 2,32 MB |
| `PropertyItemListCoverageSpreading` | 12 | 0 | 1,16 MB |

### 7.2 Grid bersarang di dalam grid — [terverifikasi]

Didefinisikan operasional: section **A** memuat ≥1 `RepeatGrid` **dan** menyisipkan section **B**
yang juga memuat ≥1 `RepeatGrid`.

```powershell
# $edges dari §3.1, $grid = himpunan "Folder|Section" yang RepeatGrid>0
$nest = $edges | Where-Object { $_.SrcType -eq 'Section' -and
  $grid.ContainsKey($_.Folder+'|'+$_.Src) -and $grid.ContainsKey($_.Folder+'|'+$_.Ref) }
"pasangan grid-dalam-grid = $($nest.Count) ; section induk unik = $((($nest.Src)|Sort-Object -Unique).Count)"
```

**414 pasangan**, **68 section induk unik**. Terparah:

| Section induk | Section-grid yang disisipkan |
| --- | ---: |
| `InputInwardFacultativeDtl` | **31** |
| `PropertyItemListCoverage` | 22 |
| `CoverageItem` | 22 |
| `InputDtlObject_FacIn` | 19 |
| `InputCoverageFire` | 18 |
| `InputDtlSpreadingCoverage_FacIn` | 18 |

Contoh nyata (`NB FacIn`):
`InputInwardFacultativeDtl` → `CoverageCommisionList`, `CoverageSpreadingList`,
`InputCoverageAneka_FacIn`, `InputDtlObject_FacIn`, `InputDtlPayment_FacInLife`,
`InputEndorsementDtl`, `InputSpreadingLife_FacIn`, `ObjectDtlAneka_FacIn`,
`ObjectOccupation_FacIn`, `PaymentCurrencyList`, `SummaryLossRecord_Section`,
`SummarySpreading_Section` — masing-masing membawa gridnya sendiri.

**Risiko React:** grid bersarang 3–4 lapis dengan ribuan baris. Di Pega, setiap grid melakukan
partial refresh server-side (`pyPartialRefresh=true`, `pyGridPreActivity=pyPreGridUpdate`).
Di React, padanannya adalah virtualisasi + state normalisasi; **tidak ada padanan langsung** untuk
"refresh satu grid dari server" tanpa merancang ulang kontrak API.

### 7.3 Section yang saling merujuk (siklus) — [terverifikasi]

```powershell
# $edgesSection = sisi Section→Section dari §3.1
$pairs=@{}; foreach($r in $edgesSection){ $pairs[$r.Folder+'|'+$r.Src+'>'+$r.Ref]=$true }
$cyc=@(); foreach($r in $edgesSection){ if($pairs.ContainsKey($r.Folder+'|'+$r.Ref+'>'+$r.Src)){
  $a,$b=($r.Src,$r.Ref)|Sort-Object; $cyc += "$a <-> $b" } }
($cyc|Sort-Object -Unique)
```

**5 pasangan saling merujuk:**

| Pasangan |
| --- |
| `CoverageItem` ↔ `InputCoverageFire` |
| `CoverageItem` ↔ `PropertyItemListCoverage` |
| `InputCoverageFire` ↔ `PropertyItemListCoverage` |
| `GridViewFollowingNB` ↔ `ViewOpenFollowingPolicy` |
| `InputDtlObject_FacIn` ↔ `InputEndorsementDtl` |

Tiga yang pertama membentuk satu **komponen terhubung kuat beranggota 3**
(`CoverageItem`, `InputCoverageFire`, `PropertyItemListCoverage`) di graf section `NB FacIn`
(116 node bersumber, 302 sisi).

```powershell
# reach maju & mundur dari CoverageItem pada graf Section NB FacIn
# hasil: maju=40, mundur=20, irisan (SCC) = 3
```

**Risiko React:** komponen yang me-render dirinya sendiri secara tidak langsung. Di Pega ini aman
karena rendering dipotong oleh kondisi visibilitas server-side; di React ini rekursi tanpa batas
kecuali kedalaman dibatasi secara eksplisit. **Ini harus dirancang, bukan diterjemahkan.**

### 7.4 Modal bertingkat — [terverifikasi]

```powershell
$tot=0;$files=0
foreach($f in @("NB FacIn","RNW Fac In","Endorsment Fac In")){
 Get-ChildItem "D:\migrasi\RNM\$f\Section" -File -Filter *.xml | ForEach-Object {
  $m=[regex]::Matches([IO.File]::ReadAllText($_.FullName),'<pyShowModalDialog>true</pyShowModalDialog>').Count
  $tot+=$m; if($m -gt 0){$files++} } }
"aksi buka-modal = $tot ; berkas Section terdampak = $files"
```

**1.045** aksi membuka jendela modal, tersebar di **174** berkas Section. Nilai
`<pyShowModalDialog>` hanya pernah `true` (tidak ada `false` eksplisit).

| Section | Aksi buka-modal |
| --- | ---: |
| `InputDtlObject_FacIn` / `InputDtlObject_FacIn_IsUW` | **30** masing-masing (× 3 folder) |
| `InputOtherObjectAneka_FacIn`, `InputEndorsementDtl`, `InputInwardFacultativeDtl_IsUW`, `InputRenewalDtl(_IsUW)`, `OldDataEndorsementDtl`, `InputOtherObjectAneka_UW` | 12 masing-masing |

Bukti bahwa ini modal sungguhan (bukan panel inline) — kelas animasi modal ada di berkas:
`<pxObjClass>Embed-DesktopAPI-OpenModalWindow-ModalAnimationSettings</pxObjClass>` dengan
`<pyLaunch>anim-bottom</pyLaunch>` di `D:\migrasi\RNM\NB FacIn\Section\OfferFacIn_NusaRe.xml`.

**Bertingkat — [terverifikasi]:** `InputDtlObject_FacIn` membuka 30 modal, **dan** section itu
sendiri dirender oleh flow action yang dibuka dari modal induk (`ViewPolis`, `ViewDataOfferFacIn`
merujuknya). Kedalaman pasti **belum terverifikasi**; yang terbukti adalah ada ≥2 tingkat.

### 7.5 Visibilitas kondisional — [terverifikasi]

```powershell
$v=@{}; $cond=@{}
foreach($f in @("NB FacIn","RNW Fac In","Endorsment Fac In")){
 Get-ChildItem "D:\migrasi\RNM\$f\Section" -File -Filter *.xml | ForEach-Object {
  $c=[IO.File]::ReadAllText($_.FullName)
  foreach($m in [regex]::Matches($c,'<pyVisible>([^<]+)</pyVisible>')){ $v[$m.Groups[1].Value]=1+$v[$m.Groups[1].Value] }
  foreach($m in [regex]::Matches($c,'<pyCondition>([^<]+)</pyCondition>')){ $cond[$m.Groups[1].Value]=1+$cond[$m.Groups[1].Value] } } }
$v.GetEnumerator()|Sort-Object Value -Descending
"ekspresi pyCondition unik = $($cond.Count) ; total kemunculan = $((($cond.Values)|Measure-Object -Sum).Sum)"
$cond.GetEnumerator()|Sort-Object Value -Descending|Select-Object -First 20
```

| `pyVisible` | Kemunculan |
| --- | ---: |
| `ALWAYS` | 73.165 |
| `OTHER` (ekspresi kustom di `<pyCondition>`) | **4.976** |
| `NOTBLANK` | 719 |
| `NOTZERO` | 9 |

**632** ekspresi `<pyCondition>` berbeda, **6.971** kemunculan. 20 teratas:

| Kemunculan | Ekspresi |
| ---: | --- |
| 1.083 | `Other Property` *(placeholder editor, bukan ekspresi)* |
| **857** | **`1=2`** |
| 276 | `When Rule` *(placeholder)* |
| 128 | `!IsSpreadingUW` |
| **127** | **`Never`** |
| 87 | `IsObjectSectionAneka` |
| 86 | `IsFire` |
| 74 | `IsHE` |
| 73 | `IsEDM` |
| 73 | `!IsUW` |
| 72 | `IsObjectWithQuantityYear \|\| IsHE` |
| **62** | **`1 = 2`** |
| 62 | `IsCreditBriguna` |
| 62 | `!IsFire` |
| 60 | `!IsClaim && .FlagDelete != '1'` |
| 50 | `.CoverageBasis==2` |
| 47 | `.FlagDelete!=1` |
| 44 | `IsFac` |
| 42 | `pyWorkPage.ProposalPosition = 'H1' \|\| … = '1' \|\| … = '1A'` |
| 42 | `pyWorkPage.ProposalPosition = '2' \|\| … = '6'` |

**Temuan #3 — [terverifikasi]: UI mati yang dikodekan keras.** `1=2` + `1 = 2` + `Never` =
**936** kemunculan, tersebar di **244** berkas Section.

```powershell
$n=0;$files=0
foreach($f in @("NB FacIn","RNW Fac In","Endorsment Fac In")){
 Get-ChildItem "D:\migrasi\RNM\$f\Section" -File -Filter *.xml | ForEach-Object {
  $m=[regex]::Matches([IO.File]::ReadAllText($_.FullName),'<pyCondition>\s*(1\s*=\s*2|Never)\s*</pyCondition>').Count
  $n+=$m; if($m -gt 0){$files++} } }
"kemunculan '1=2'/'Never' = $n ; berkas terdampak = $files"
```

Terbanyak: `InputInwardFacultativeDtl` (32 × 3 folder), `InputDtlObject_FacIn_IsUW` (25 × 3),
`InputEndorsementDtl` (23), `InputRenewalDtl(_IsUW)` (21 masing-masing).

Ini **kandidat perbaikan, bukan keputusan migrasi** (lihat CLAUDE.md §1). Elemen ber-`1=2`
tidak pernah tampil di Pega; membawanya ke React berarti membawa mati-suri. Membuangnya berarti
memutuskan sendiri bahwa bisnis tidak akan pernah menyalakannya lagi. **Harus ditanyakan.**

### 7.6 Read-only kondisional — [terverifikasi]

100 ekspresi `<pyReadOnlyCondition>` berbeda; 15 teratas:

| Kemunculan | Ekspresi |
| ---: | --- |
| 1.087 | `IsSpreadingUW \|\| IsGroupFac \|\| IsGroupUWFac` |
| **344** | **`1!=2`** *(selalu benar → selalu read-only)* |
| 296 | `IsSpreadingUW` |
| 222 | `.FlagDelete==1` |
| 202 | `pyWorkPage.FlagViewPolicy=1` |
| 198 | `IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW` |
| 144 | `IsUW` |
| 93 | `.FlagDelete = 1` |
| 90 | `.FlagDelete = '1'` |
| 63 | `IsGroupFac \|\| IsGroupUWFac` |
| **47** | **`1=1`** *(selalu benar)* |
| 45 | `pyWorkPage.IsFacRetroOffer=='1'` |
| 42 | `ALWAYS` |
| 38 | `IsUW \|\| IsGroupFac \|\| IsGroupUWFac \|\| IsSpreadingUW` |
| 36 | `IsUW \|\| IsSpreadingUW` |

Perhatikan **tiga ejaan berbeda untuk satu maksud**: `.FlagDelete==1`, `.FlagDelete = 1`,
`.FlagDelete = '1'` — totalnya 405 kemunculan. Yang terakhir membandingkan **sebagai string**.
Ini sejalan dengan peringatan di CLAUDE.md §4.1 dan **tidak boleh diseragamkan diam-diam**.

Kondisi read-only bergantung pada rule `When`: `IsSpreadingUW`, `IsGroupFac`, `IsGroupUWFac`,
`IsUW`. Perilaku aktualnya ada di `OUTPUT/04-aturan/01-katalog-when.md`, bukan di sini.

### 7.7 Rendering sisi-klien vs sisi-server — [terverifikasi]

| Penanda | Kemunculan (semua Section) |
| --- | ---: |
| `<pyIsClientWhen>true</…>` — kondisi dievaluasi di browser | **2.586** |
| `<pyPreDataTransform>` dengan `<pyName>` terisi — transform server sebelum render | **2.306** |
| `<pyLoadDeferred>true</…>` | 50 |
| `<pyRefreshWhen>` terisi | 6 |

2.306 pre-data-transform berarti **isi layar dihitung di server sebelum dirender**. Di React,
setiap satu dari itu menjadi keputusan: panggil endpoint, atau pindahkan logika ke klien.
Memindahkannya ke klien mengubah perilaku dan melanggar prinsip migrasi. **[pertanyaan terbuka]**
apakah backend Go akan mengekspos satu endpoint "render state" per layar, atau logikanya
dipindah ke frontend.

---

## 8. Usulan pemetaan ke struktur React target

> **Seluruh isi §8 adalah USULAN**, bukan temuan korpus. Korpus tidak memuat struktur folder React
> mana pun. Dasar usulan adalah angka §3.1, §6, dan §7; setiap baris menyebut angka yang
> mendasarinya agar dapat dibantah.

```
frontend/src/
├── pages/            <- dari Harness (§2) + FlowAction berclosure besar (§4)
├── components/
│   ├── grids/        <- dari pyPageListProperty teratas (§6)
│   ├── panels/       <- dari Section komponen bersama (§3.1)
│   ├── pickers/      <- dari Harness Choose*/Select* + FlowAction Choose*/Select*
│   └── layout/
├── hooks/
├── services/
└── store/
```

### 8.1 `pages/` — usulan

| Halaman React usulan | Asal korpus | Dasar angka |
| --- | --- | --- |
| `OfferFacInPage` | FlowAction `InwardFacultative` / `InwardFacultative_IsUW` → `InputInwardFacultative(_2)` | `InputInwardFacultative` dirujuk 62× — tertinggi kedua |
| `OfferDetailPage` | `InputInwardFacultativeDtl(_IsUW)` | 34× rujukan; 31 grid bersarang |
| `EndorsementPage` | FlowAction `Endorsement_FlowAct(_IsUW)` → `InputEndorsement(_IsUW)` | EDM saja |
| `RenewalPage` | FlowAction `Renewal_FlowAct(_IsUW)` → `InputRenewal(_IsUW)` | RNW saja, 8 section khas |
| `PolicyViewPage` | Harness `ViewPolis` | closure 33–44 section |
| `OldDataViewPage` | Harness `ViewOldEndorsement` / `ViewOldData` | closure 36–41 |
| `AccumulationPage` | Harness `AccumulationRisk`, `TotalAccumulation_FacIn`, `SummaryRiskAccumulation`, `RiskAccumulationReport`, `AdjustmentRiskAccumulation` | 5 harness `DATA-PORTAL` |
| `PrintRISlipPage` | Harness `PrintRISlip` / `PrintRISlips` | 2 harness |
| `LetterPage` | Harness `ViewLetter` → `FacultativeLetter` | 1 harness |
| `AttachmentPage` | Harness `ViewPictureList` → `ReasViewAttachment` | 1 harness |
| `ClaimListPage` | Harness `ViewClaimList` | 1 harness |
| `CsvImportResultPage` | Harness `ViewCSVResult_B2BHost` | NB/RNW saja |

**Tidak diusulkan sebagai `pages/`:** `TabbedScreenFlow7` (kerangka navigasi bawaan Pega — di React
digantikan router, bukan diterjemahkan), `SFAPortal*` (portal CRM Pega).

### 8.2 `components/grids/` — usulan, diurutkan menurut bukti rujukan

| Komponen usulan | Dari Section | Binding korpus | Rujukan |
| --- | --- | --- | ---: |
| `<CurrencyPaymentGrid>` | `PaymentCurrencyList` | `.OfferFacIn.CurrencyList` (15 section unik) | **83** |
| `<ObjectListGrid variant>` | `CoverageSpreadingList` + `CoverageCommisionList` | `.PersonList` `.LocationList` `.VehicleList` `.CargoList` | **60 + 53** |
| `<CoverageGrid>` | `PropertyItemListCoverage`, `CoverageItem`, `InputCoverageFire` | `.CoverageList` (43 section unik), `.Property.PropertyItemList` (11) | 22 + 23 + 19 |
| `<RiskLocationGrid>` | `InputRiskAddress`, `ObjectList` | `.LocationList` (15), `.Property.RiskLocation.OccupationList` (11) | 18 + 12 |
| `<SpreadingGrid>` | `SpreadingItem`, `InputDtlSpreadingCoverage_FacIn` | `SpreadingList.pxResults` (8), `TotalSpreadAll.pxResults` (9) | 19 + 13 |
| `<DeductibleGrid>` | `InputDeductible_FacIn` | `.DeductibleList` (17) | 3 |
| `<ParticipantGrid>` | `InputPerson`, `InputDtlParticipant*` | `.PersonList` (14), `.ASMHeir` (11) | 9 |
| `<VehicleGrid>` | `VehicleGrid` | `.VehicleList` (14) | 1 |
| `<ClauseGrid>` | `CreateListClause`, `SelClauseList` | `.ClauseList` (8) | – |

### 8.3 `components/pickers/` — usulan

Semua harness/flow action `Choose*` dan `Select*` berbagi satu bentuk: cari → daftar hasil →
pilih → tutup. 13 flow action + 10 harness.

**Usulan:** satu `<EntityPicker>` generik, di-parameterkan sumber data dan kolom, menggantikan:
`ChooseAccumulation_FacIn`, `ChooseClauseFire`, `ChooseDeductible`, `ChooseInsured`,
`ChooseOccupation`, `ChooseRiskAddress`, `ChooseRiskLocation`, `ChooseZipCode`,
`ChooseClassofContraction`, `ChooseCoverage`, `ChooseObject_Ship`, `ChooseSubContract_FacIn`,
`SelectAgent`, `SelectCoverage`, `SelectShip`, `SelectBenefitList`, `SelectClauseList`,
`SelectPlanList`, `ShowCedingCoList`, `CedingCedant`, `CedingCompany`.

**Catatan jujur:** 8 dari daftar itu **bercabang** antara NB/RNW dan EDM (§5.3, §5.5). Generalisasi
ini hanya sah bila perbedaan itu sudah dipetakan lebih dulu.

### 8.4 `components/panels/` — usulan

| Komponen usulan | Dari Section | Rujukan |
| --- | --- | ---: |
| `<NusaReSharePanel>` | `OfferFacIn_NusaRe` | 31 |
| `<CorrespondencePanel>` | `EmailSection`, `Correspondence` | 22 |
| `<ScoringRiskPanel>` | `ScoringRisk` (+ `ScoringRiskForm1..4`) | 12 |
| `<FacOfferPanel>` | `AddFacOfferList`, `InputFacOffer` | 19 + 12 |
| `<OfferStatusPanel>` | `OfferStatusFacOut`, `OldOfferStatusFacOut` | – |
| `<PeriodPanel>` | `Periode`, `PeriodeEndorsement`, `PeriodeRenewal` (+ varian `_IsUW`) | 9 |
| `<LossRecordPanel>` | `InputLossExperiance`, `InputLossRecord_Sec`, `SummaryLossRecord_Section` | – |

### 8.5 `hooks/` dan `store/` — usulan

| Usulan | Dasar |
| --- | --- |
| `useUwMode()` | 97 pasangan section UW/non-UW (§6.2); mode menjadi prop, bukan komponen kedua |
| `useCycle()` — `'NB' \| 'RNW' \| 'EDM'` | NB = RNW pada 575/575 layar (§5.2); percabangan hanya EDM. Satu flag, bukan tiga pohon komponen |
| `useVisibility(expr)` | 632 ekspresi `<pyCondition>` (§7.5) — evaluator, bukan 632 `if` tersebar |
| `useReadOnly(expr)` | 100 ekspresi `<pyReadOnlyCondition>` (§7.6) |
| store: normalisasi `CoverageList` / `LocationList` / `PersonList` / `VehicleList` / `CargoList` / `DeductibleList` | 6 properti list teratas mengikat 43 + 15 + 14 + 14 + 9 + 17 section unik (§6) |
| store: `currencyList` sebagai agregat mandiri | `.OfferFacIn.CurrencyList` (dua ejaan) dipakai **15 section unik**; setiap nilai uang terikat mata uang (CLAUDE.md §4.1) |

### 8.6 `services/` — usulan

2.306 `pyPreDataTransform` dan 1.810 grid dengan `pyGridPreActivity=pyPreGridUpdate` /
`pyGridPostActivity=pyPostGridUpdate` (§7.1, §7.7) berarti UI lama **mengandalkan server untuk
menghitung isi layar sebelum dan sesudah setiap perubahan grid**. Usulan: `services/` memaparkan
satu kontrak per layar (`loadX`, `recalcX`) alih-alih CRUD per-entitas, agar urutan perhitungan
lama dapat direproduksi saat paralel run.

**Ini usulan arsitektur, bukan temuan.** Korpus tidak menyebut satu pun endpoint; sesuai
CLAUDE.md §4.4, daftar endpoint sesungguhnya ada di `M_LINK_SERVICE` yang tidak ada di korpus.

---

## 9. Pertanyaan terbuka

1. **Portal & titik masuk.** Harness mana yang menjadi layar pertama underwriter? Tidak terjawab
   dari `Harness\`, `Section\`, `FlowAction\`; informasi itu ada di `Rule-Portal` / access group
   yang tidak ada di korpus. — **[pertanyaan terbuka]**

2. **Urutan langkah layar.** 307 flow action tidak membawa informasi urutan; screen-flow ada di
   `Flow\` (6 + 4 + 2 berkas). Dokumen ini sengaja tidak menyentuhnya. Tanpa itu, `pages/` di §8.1
   masih berupa kumpulan, bukan alur. — **[pertanyaan terbuka]**

3. **Mengapa NB dan RNW identik 575/575 di lapisan UI?** Apakah Renewal memang memakai layar NB
   apa adanya, ataukah ekspor RNW sebenarnya salinan ruleset NB? Jawabannya menentukan apakah
   `useCycle()` (§8.5) diperlukan sama sekali. — **[pertanyaan terbuka]**

4. **Isi 126 percabangan EDM.** Untuk 6 harness, 5 di antaranya berbeda **kelas Pega**
   [terverifikasi]; sisa perbedaan (76 section, 44 flow action, `TotalAccumulationDtl`) **belum
   dipetakan field-per-field**. Sebelum dipetakan, tidak boleh diklaim "hanya beda kosmetik".
   — **[pertanyaan terbuka]**

5. **936 elemen `1=2` / `Never`.** Sengaja dimatikan untuk sementara, atau sisa fitur yang
   dibatalkan? Migrasi tidak boleh memutuskannya sendiri (CLAUDE.md §1). — **[pertanyaan terbuka]**

6. **344 `1!=2` + 47 `1=1` pada `pyReadOnlyCondition`.** Keduanya selalu benar → field selalu
   read-only. Apakah itu maksudnya, atau sisa debugging? — **[pertanyaan terbuka]**

7. **Tiga ejaan `.FlagDelete`** (`==1`, `= 1`, `= '1'`, total 405 kemunculan). Yang terakhir
   perbandingan string. Apakah ketiganya berperilaku sama di Pega? — **[pertanyaan terbuka]**

8. **141 section yatim.** Benar-benar mati, atau dirujuk dari `Flow`/`Activity`/`DataPage` di luar
   cakupan dokumen ini? Khususnya rumpun `*FacOut*` (±15) dan `*TreatyIn*` (5): apakah termasuk
   lingkup migrasi Facultative **Inward**? — **[pertanyaan terbuka]**

9. **97 pasangan UW/non-UW berselisih 20–35 %** (§6.2). Perbedaannya read-only saja, atau ada field
   yang hanya ada di satu sisi? Menentukan apakah `mode` cukup sebagai prop.
   — **[pertanyaan terbuka]**

10. **Kedalaman modal sesungguhnya.** Terbukti ada ≥2 tingkat; jumlah tingkat maksimum belum
    diukur. Menentukan apakah React memerlukan manajer modal bertumpuk. — **[pertanyaan terbuka]**

11. **Batas ukuran B (§5.1).** Perbandingan mengabaikan **urutan** elemen XML. Dua salinan dengan
    urutan kolom berbeda akan terbaca identik. Apakah perlu pembandingan pohon XML per-node
    sebelum angka 575/575 dijadikan dasar keputusan? — **[pertanyaan terbuka]**

12. **5 flow action tanpa `pySectionReference`** (`AgentSourceBizDetails` ×3,
    `StartScreenFlowAuto` ×2) dan **21 rujukan section yang berkasnya tidak ada**
    (`InputClauseFire_ViewDtl`, `InputDeductibleDtlAneka_GCNM`, `InputPackageDM`,
    `ViewAdditionalCoverage`). Ekspor tidak lengkap, atau rule memang tidak ada?
    — **[pertanyaan terbuka]**

13. **Kontrak render server-side.** 2.306 `pyPreDataTransform` + 1.810 `pyGridPreActivity`:
    apakah backend Go akan mereproduksinya sebagai endpoint per-layar, atau logika dipindah ke
    frontend? Yang kedua mengubah perilaku. — **[pertanyaan terbuka]**

14. **Flow action bernuansa AI di jalur UI.** `AnalysLocationbyAI` dan `AttachDoc_AI` ada di
    `NB FacIn\FlowAction\` dan `RNW Fac In\FlowAction\`, **tidak ada** di `Endorsment Fac In\`
    **[terverifikasi]**. Isinya tidak dibaca dokumen ini — nama bukan bukti perilaku. Bila
    keduanya benar-benar memanggil model AI pihak ketiga, CLAUDE.md §6 mengharuskan tinjauan
    keamanan sebelum dipindahkan. — **[pertanyaan terbuka]**

    ```powershell
    foreach($f in @("NB FacIn","RNW Fac In","Endorsment Fac In")){
     foreach($n in @('AnalysLocationbyAI','AttachDoc_AI')){
      "{0,-18} {1,-22} {2}" -f $f,$n,(Test-Path "D:\migrasi\RNM\$f\FlowAction\$n.xml") } }
    ```
