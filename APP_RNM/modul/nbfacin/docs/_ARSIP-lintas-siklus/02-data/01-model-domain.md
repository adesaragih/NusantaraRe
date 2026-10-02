# 01 — Model Domain Facultative Inward (dari korpus Pega)

**Sumber tunggal:** `D:\migrasi\RNM\NB FacIn\`, `D:\migrasi\RNM\RNW Fac In\`, `D:\migrasi\RNM\Endorsment Fac In\` — 6.071 berkas `.xml`.
**Status korpus:** READ-ONLY. Dokumen ini tidak memuat nilai nama orang, hostname, maupun secret.

> **Peringatan metodologis yang mengikat seluruh dokumen ini.**
> Korpus **tidak memuat satu pun rule `Rule-Obj-Property`** (lihat §1.2). Karena itu:
> tipe data properti, panjang, nilai wajib/opsional, dan daftar nilai yang sah **tidak dapat dibaca**.
> Struktur objek di bawah ini disusun dari **jalur properti yang tertulis harfiah** di dalam
> Activity / DataTransform / Section / When. Jalur yang tertulis harfiah = `[terverifikasi]`
> (objeknya memang punya anak itu). Arti, tipe, dan kardinalitasnya = `[dugaan]` atau
> `[pertanyaan terbuka]`.

---

## 1. Kelas Pega yang dipakai

### 1.1 Sebaran rule menurut tipe

| `pxObjClass` | Jumlah berkas |
| --- | --- |
| Rule-Obj-Activity | 1.747 |
| Rule-HTML-Section | 1.263 |
| Rule-Obj-FlowAction | 754 |
| Rule-Connect-SQL | 605 |
| Rule-Obj-When | 601 |
| Rule-Obj-Model (Data Transform) | 443 |
| Rule-Obj-Report-Definition | 379 |
| Rule-HTML-Harness | 119 |
| Rule-Declare-Pages (Data Page) | 101 |
| Rule-Declare-DecisionTable | 32 |
| Rule-Obj-Flow | 12 |
| Rule-Connect-REST | 9 |
| Rule-Declare-DecisionTree | 3 |
| Rule-Admin-System-Settings | 3 |
| **Total** | **6.071** |

`[terverifikasi]` — dibaca dari tag `<pxObjClass>` pertama tiap berkas.

**Perintah audit**

```powershell
$files = Get-ChildItem -Path "D:\migrasi\RNM\NB FacIn","D:\migrasi\RNM\RNW Fac In","D:\migrasi\RNM\Endorsment Fac In" -Recurse -File -Filter *.xml
$tal=@{}
foreach($f in $files){ $t=[System.IO.File]::ReadAllText($f.FullName)
  $m=[regex]::Match($t,'(?m)^<pxObjClass>([^<]*)</pxObjClass>')
  if($m.Success){ $v=$m.Groups[1].Value; if($tal.ContainsKey($v)){$tal[$v]++}else{$tal[$v]=1} } }
$tal.GetEnumerator()|Sort-Object Value -Descending|ForEach-Object{"{0,6}  {1}" -f $_.Value,$_.Key}
```

### 1.2 Konsekuensi: tidak ada rule Property di korpus

`Rule-Obj-Property` **tidak muncul sama sekali** pada daftar di atas (0 dari 6.071).
Ini `[terverifikasi]` dan merupakan batas pengetahuan terbesar di area data — lihat
`03-batas-pengetahuan.md` §3.

### 1.3 Rule per kelas — 40 teratas

Jumlah berkas menurut `<pyClassName>` (kelas tempat rule di-*resolve*).
Total **211 kelas unik**, 6.071 berkas.

| # | Kelas | Rule |
| --- | --- | --- |
| 1 | `ASM-FW-GISFW-Work` | 1.592 |
| 2 | `ASM-FW-GISFW-Data-Coverage` | 554 |
| 3 | `@baseclass` | 283 |
| 4 | `ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance` | 259 |
| 5 | `Data-Party-Person` | 255 |
| 6 | `ASM-FW-GISFW-Int-policyjson` | 241 |
| 7 | `ASM-FW-GISFW-Data-OfferFacIn` | 206 |
| 8 | `ASM-FW-GISFW-Data-PropertyItem` | 157 |
| 9 | `ASM-FW-GISFW-Data-Aneka` | 157 |
| 10 | `ASM-FW-GISFW-Data-Occupation` | 143 |
| 11 | `ASM-FW-GISFW-Data-FacOffer` | 130 |
| 12 | `ASM-FW-GISFW-Data-Vehicle` | 97 |
| 13 | `ASM-FW-GISFW-Data-Deductible` | 87 |
| 14 | `ASM-FW-GISFW-Data-Cargo` | 87 |
| 15 | `ASM-FW-GISFW-Data` | 86 |
| 16 | `Data-Portal` | 78 |
| 17 | `ASM-FW-GISFW-Int-OFFERJSON` | 75 |
| 18 | `ASM-FW-GISFW-Data-Clause` | 67 |
| 19 | `ASM-FW-GISFW-Data-Quotation` | 67 |
| 20 | `ASM-FW-GISFW-Data-Property` | 60 |
| 21 | `ASM-FW-GISFW-Data-OfferFacIn-Currency` | 52 |
| 22 | `ASM-FW-GISFW-Data-PolicyTreatyIn` | 50 |
| 23 | `ASM-FW-GISFW-Int-T_STORAGE_IMAGE` | 40 |
| 24 | `Code-Pega-List` | 38 |
| 25 | `ASM-FW-GISFW-Int-RW` | 35 |
| 26 | `ASM-FW-GISFW-Data-Policy` | 32 |
| 27 | `Data-Address` | 31 |
| 28 | `ASM-SFAGIS-Work-Endorsement` | 31 |
| 29 | `ASM-FW-GISFW-Int-OCCUPATION` | 30 |
| 30 | `ASM-FW-GISFW-Int-CURRENCY` | 28 |
| 31 | `ASM-FW-GISFW-Int-ACCUMULATION` | 27 |
| 32 | `ASM-FW-GISFW-Data-CauseOfLoss` | 26 |
| 33 | `ASM-FW-GISFW-Int-AGENT` | 24 |
| 34 | `ASM-FW-GISFW-Data-ScoringRisk` | 21 |
| 35 | `ASM-FW-GISFW-Data-Outgo` | 21 |
| 36 | `ASM-FW-GISFW-Data-SubContract` | 21 |
| 37 | `ASM-FW-GISFW-Int-marketingofficer` | 21 |
| 38 | `Work-` | 20 |
| 39 | `ASM-FW-GISFW-Data-Good` | 20 |
| 40 | `ASM-FW-GISFW-Int-BRANDDETAIL` | 18 |

**Perintah audit**

```powershell
$files = Get-ChildItem -Path "D:\migrasi\RNM\NB FacIn","D:\migrasi\RNM\RNW Fac In","D:\migrasi\RNM\Endorsment Fac In" -Recurse -File -Filter *.xml
$res = foreach ($f in $files) {
  $m = Select-String -Path $f.FullName -Pattern '^<pyClassName>([^<]*)</pyClassName>' | Select-Object -First 1
  if ($m) { $m.Matches[0].Groups[1].Value }
}
Write-Output ("class unik: " + ($res | Group-Object).Count + " ; total: " + $res.Count)
$res | Group-Object | Sort-Object Count -Descending | Select-Object -First 40 Count,Name
```

### 1.4 Kelas Work — identifikasi class Work utama

| Kelas Work | Rule | Catatan |
| --- | --- | --- |
| **`ASM-FW-GISFW-Work`** | **1.592** | **kelas Work utama** — 26,2 % seluruh korpus |
| `ASM-SFAGIS-Work-Endorsement` | 31 | |
| `ASM-FW-GISFW-Work-NB` | 14 | |
| `ASM-FW-SFAGISFW-Work-Account` | 5 | |
| `ASM-FW-GISFW-Work-TKUSales` | 5 | |
| `ASM-FW-GISFW-Work-Endorsement` | 4 | |
| `ASM-FW-GISFW-WORK-INT-V_ZIPCODE` | 3 | |
| `ASM-FW-SFAGISFW-Work-Opportunity` | 3 | |
| `ASM-FW-GISFW-Work-LIFE` | 2 | |
| `ASM-FW-GCNMFW-Work` | 1 | |
| `ASM-FW-GISFW-Work-NewBusinessGuarantee` | 1 | |
| `ASM-FW-GISFW-Work-Renewal` | 1 | |

`[terverifikasi]` `ASM-FW-GISFW-Work` adalah kelas Work utama: seluruh flow inti, seluruh
rule `When` tangga akseptasi, dan seluruh Activity produksi di-*resolve* di sana.

`[terverifikasi]` Pembuatan case endorsement memakai kelas **`ASM-FW-GISFW-Work-Endorsement`**,
bukan `ASM-SFAGIS-Work-Endorsement`. Bukti — `Endorsment Fac In/Activity/GetEDMData.xml`:

```xml
<PropertiesName>param.classname</PropertiesName>
<PropertiesValue>"ASM-FW-GISFW-Work-Endorsement"</PropertiesValue>
<PropertiesName>param.IDPrefix</PropertiesName>
<PropertiesValue>"EDM-"</PropertiesValue>
```

`[pertanyaan terbuka]` Peran `ASM-SFAGIS-Work-Endorsement` (31 rule) tidak terbaca dari korpus.

**Perintah audit**

```powershell
$files = Get-ChildItem -Path "D:\migrasi\RNM\NB FacIn","D:\migrasi\RNM\RNW Fac In","D:\migrasi\RNM\Endorsment Fac In" -Recurse -File -Filter *.xml
$files | ForEach-Object {
  $m = Select-String -Path $_.FullName -Pattern '^<pyClassName>([^<]*)</pyClassName>' | Select-Object -First 1
  if ($m) { $m.Matches[0].Groups[1].Value } } |
  Where-Object { $_ -like 'ASM-*Work*' } | Group-Object | Sort-Object Count -Descending
```

### 1.5 Konvensi penamaan kelas (pola, bukan bukti)

`[dugaan]` — dari pola nama, tidak dikonfirmasi oleh rule Property mana pun:

| Segmen | Dugaan peran | Contoh |
| --- | --- | --- |
| `…-Work` | kelas case | `ASM-FW-GISFW-Work` |
| `…-Data-*` | kelas embedded/page di dalam agregat case | `ASM-FW-GISFW-Data-Coverage` |
| `…-Int-*` | kelas terpeta ke **tabel/view Oracle** eksternal | `ASM-FW-GISFW-Int-M_LINK_SERVICE` |

Dugaan `…-Int-*` = tabel Oracle **didukung satu bukti kuat**: `NB FacIn/Activity/GetLinkService.xml`
melakukan `Obj-Browse` atas `<ObjClass>ASM-FW-GISFW-Int-M_LINK_SERVICE</ObjClass>` dengan filter
`<Field>.KATEGORI_1</Field>` dan `<Field>.KATEGORI_2</Field>` — kolom bergaya Oracle, bukan properti
bergaya Pega. Namun **pemetaan kelas→tabel tidak ada di korpus**, jadi tetap `[dugaan]`.

---

## 2. Agregat utama case dan sub-agregatnya

### 2.1 Halaman puncak pada work page

`[terverifikasi]` Ada **tiga** halaman pembawa data bisnis yang hidup berdampingan di `pyWorkPage`:

| Halaman | Rujukan | Berkas | Peran terbaca |
| --- | --- | --- | --- |
| `pyWorkPage.OfferFacIn` | 20.570 | 1.391 | **agregat case utama** |
| `pyWorkPage.Quotation` | 5.180 | 569 | halaman quotation/entry paralel |
| `pyWorkPage.Policy` | 823 | 186 | halaman polis paralel |
| `newWorkPage.OfferFacIn` | 297 | 12 | agregat pada case baru (spin-off EDM/RNW) |

Rujukan token `OfferFacIn` di seluruh korpus (semua konteks): **71.014** dalam **2.087** berkas.

**Perintah audit**

```powershell
$files = Get-ChildItem -Path "D:\migrasi\RNM\NB FacIn","D:\migrasi\RNM\RNW Fac In","D:\migrasi\RNM\Endorsment Fac In" -Recurse -File -Filter *.xml
$pats = @{
  'pyWorkPage.OfferFacIn' = 'pyWorkPage\.OfferFacIn(?![A-Za-z0-9_])'
  'pyWorkPage.Quotation'  = 'pyWorkPage\.Quotation(?![A-Za-z0-9_])'
  'pyWorkPage.Policy'     = 'pyWorkPage\.Policy(?![A-Za-z0-9_])'
  'newWorkPage.OfferFacIn'= 'newWorkPage\.OfferFacIn(?![A-Za-z0-9_])'
  'OfferFacIn (token)'    = '(?<![A-Za-z0-9_])OfferFacIn(?![A-Za-z0-9_])' }
foreach($k in $pats.Keys){ $rx=[regex]$pats[$k]; $c=0; $fc=0
  foreach($f in $files){ $t=[System.IO.File]::ReadAllText($f.FullName)
    $n=$rx.Matches($t).Count; if($n -gt 0){$fc++;$c+=$n} }
  Write-Output ("{0,-24} {1,7} rujukan / {2,5} berkas" -f $k,$c,$fc) }
```

`[pertanyaan terbuka]` **Duplikasi `Quotation` vs `OfferFacIn.QuotationData`.**
Keduanya membawa field bernama sama (`BusinessCode`, `BusinessType`, `StatusBusiness`,
`BusinessOldId`, `OldPolicyNo`, `TeamGroup`, `EdmType`, …) dan **keduanya sama-sama dibaca**.
Contoh `[terverifikasi]` keduanya ditulis ke kolom yang sama di SQL yang sama:
`NB FacIn/RDBList/InsertTreatyProduction_Sql.xml` memakai
`{pyWorkPage.Quotation.StatusBusiness}` untuk kolom `STATUS_BUSINESS`, sementara
`{pyWorkPage.OfferFacIn.QuotationData.GroupName}` untuk `GROUPNAME`.
Mana yang otoritatif, dan apakah keduanya wajib sinkron, **tidak terjawab dari korpus**.

### 2.2 Sub-agregat langsung di bawah `OfferFacIn` — urut jumlah rujukan

Dihitung dari pola harfiah `OfferFacIn.<Nama>` di seluruh korpus.
**Urutan ini menunjukkan tulang punggung data.**

| # | Sub-agregat | Rujukan | Bentuk terbaca |
| --- | --- | --- | --- |
| 1 | `QuotationData` | 7.197 | page |
| 2 | `LocationList` | 3.544 | **page-list** |
| 3 | `PolicyData` | 2.501 | page |
| 4 | `OldData` | 1.565 | page (cermin agregat — lihat §5) |
| 5 | `FacRetroList` | 1.050 | **page-list** |
| 6 | `CurrencyList` | 855 | **page-list** |
| 7 | `PercentShare` | 802 | skalar |
| 8 | `PersonList` | 783 | **page-list** |
| 9 | `Obligee` | 588 | page/skalar (belum terverifikasi) |
| 10 | `CargoList` | 493 | **page-list** |
| 11 | `VehicleList` | 426 | **page-list** |
| 12 | `ProRatePercent` | 353 | skalar |
| 13 | `ViewSuggest` | 258 | belum terverifikasi |
| 14 | `ProrateEDMEnd` | 249 | skalar |
| 15 | `ShareRNML` | 200 | skalar |
| 16 | `ScoringRisk` | 192 | page |
| 17 | `CedingCedantList` | 172 | **page-list** |
| 18 | `TotalTSINusaReSpreading` | 152 | skalar (uang) |
| 19 | `IsInputFacRetro` | 144 | bendera |
| 20 | `IsOccupException` | 143 | bendera |
| 21 | `FacRetro` | 141 | page |
| 22 | `TotalPremiNusaRe` | 140 | skalar (uang) |
| 23 | `IsB2B` | 131 | bendera bernilai `""` / `"ASM"` |
| 24 | `ProRateType` | 126 | enumerasi |
| 25 | `IsRISlip` | 124 | bendera |
| 26 | `IsPreferredRisk` | 123 | bendera |
| 27 | `IsSpecialAcceptance` | 106 | bendera |
| 28 | `FacRetroDetails` | 102 | belum terverifikasi |
| 29 | `IsFacRetro` | 93 | bendera |
| 30 | `ShareCedantType` | 76 | enumerasi |
| 31 | `TotalTSINusaRe` | 69 | skalar (uang) |
| 32 | `TotalTSITopRisk` | 61 | skalar (uang) |
| 33 | `ConfirmBinding` | 54 | belum terverifikasi |
| 34 | `CedingRetention` | 51 | skalar |
| 35 | `ProrateStartEDM` | 49 | skalar |
| 36 | `GoodsList` | 20 | **page-list** |
| 37 | `CorrespondenceList` | 18 | **page-list** |
| 38 | `ListCauseOfLoss` | 34 | **page-list** |
| 39 | `Parameters` | 36 | page |
| 40 | `PolicyMasterNumber` | 38 | skalar |

Label "page-list" `[terverifikasi]` bila korpus memuat pemakaian berindeks `Nama(<indeks>)`;
"page" bila hanya `Nama.<anak>`; "skalar" bila tidak pernah diikuti titik maupun kurung.
Kolom yang tertulis *belum terverifikasi* berarti korpus menyebut namanya tanpa konteks
yang membuktikan bentuknya.

**Perintah audit**

```powershell
$files = Get-ChildItem -Path "D:\migrasi\RNM\NB FacIn","D:\migrasi\RNM\RNW Fac In","D:\migrasi\RNM\Endorsment Fac In" -Recurse -File -Filter *.xml
$rx=[regex]'OfferFacIn\.([A-Za-z][A-Za-z0-9_]*)'
$tal=@{}
foreach ($f in $files) { $t=[System.IO.File]::ReadAllText($f.FullName)
  foreach ($m in $rx.Matches($t)) { $k=$m.Groups[1].Value
    if ($tal.ContainsKey($k)) { $tal[$k]++ } else { $tal[$k]=1 } } }
$tal.GetEnumerator() | Sort-Object Value -Descending | Select-Object -First 60 |
  ForEach-Object { "{0,7}  {1}" -f $_.Value, $_.Key }
```

### 2.3 Isi `QuotationData` — 25 anak teratas

| Anak | Rujukan | | Anak | Rujukan |
| --- | --- | --- | --- | --- |
| `BusinessOldId` | 1.190 | | `NoOfferSlip` | 82 |
| `StatusBusiness` | 1.096 | | `EdmDate` | 78 |
| `Type` | 956 | | `InsuredID` | 75 |
| `EdmType` | 655 | | `IsMOP` | 66 |
| `IsGroup` | 404 | | `IsSurveyReport` | 65 |
| `BusinessType` | 397 | | `QQName` | 61 |
| `MarketingName` | 350 | | `EdmTypeNew` | 58 |
| `TeamGroup` | 291 | | `MOID` | 58 |
| `BusinessCode` | 216 | | `Period` | 51 |
| `MarketingCode` | 206 | | `PeriodMM` | 50 |
| `PolicyType` | 146 | | `SourceOfBusiness` | 49 |
| `OldPolicyNo` | 117 | | `EDMDay` | 46 |
| `InsuredName` | 115 | | `CedingCo` | 44 |
| `CedingCoName` | 111 | | `IsMultiCoB` | 36 |
| `BusinessName` | 111 | | `GroupName` | 36 |
| `SobName` | 110 | | `SurveyReportList` | 32 |

Ekor yang penting untuk EDM: `EndorsementInternalRetro` (9), `NPLStatus` (11), `EdmNote` (25),
`FacOutStatus` (24), `IsTBA` (24), `ProportionalType` (21), `ShareOfCeding` (21),
`StatusSyariah` (19), `IsTemplate` (15), `TypeFacultative` (12), `CommentOldPolicy` (11).

### 2.4 Isi `PolicyData`

| Anak | Rujukan | Catatan |
| --- | --- | --- |
| `Payment` | 1.136 | **sub-agregat sendiri** — lihat §2.7 |
| `StartDateTime` | 782 | |
| `EndDateTime` | 593 | |
| `Ship` | 381 | page (marine) |
| `ProdDateTime` | 203 | tanggal produksi |
| `PolicyNo` | 192 | |
| `SurveyAgent` | 135 | |
| `LC` | 74 | |
| `EndorsementNo` | 61 | |
| `SailDate` | 55 | |
| `MaxAgeDetail` / `MaxAge` | 24 / 12 | |
| `OfferingDate` | 24 | |
| `InvoiceNumber` / `InvoiceDate` | 18 / 12 | |
| `BLNumber` | 14 | |

Terbaca juga salah-ketik yang bertahan di korpus: `EndtDateTime` (3), `SttDateTim` (3),
`EndtDateTim` (3), `E` (1). `[terverifikasi]` bahwa string itu ada; artinya `[pertanyaan terbuka]`
(cacat ketik yang tak pernah terpakai, atau properti nyata yang dipakai satu tempat).

### 2.5 Isi sub-agregat berbentuk daftar

Dihitung dari pola berindeks `Nama(<apa saja>).<anak>` — ini sekaligus **bukti bahwa
sub-agregat tersebut adalah page-list**, bukan page.

**`LocationList(i).`**

| Anak | Rujukan |
| --- | --- |
| `Property` | 3.361 |
| `OccupationList` | 423 |
| `IsOldData` | 36 |
| `IndexProperty` | 21 |
| `FlagDelete` | 12 |
| `ASMAddress`, `FEAList`, `ASMStatus`, `ASMDistrict`, `ASMRW`, `ASMZipCode`, `ASMCity` | 3–7 |

**`CurrencyList(i).`** — `Policy` 213 · `Name` 104 · `TSI` 88 · `OldID` 59 ·
`SumTotalPayment` 38 · `Premium` 31 · `TSIOld` 21 · `PremiumOld` 21 · `ID` 12 ·
`TSINusantaraRe` 3 · `PremiNusantaraRe` 3 · `IndexCurrency` 1

**`FacRetroList(i).`** — `LocationList` 324 · `CargoList` 120 · `CurrencyList` 115 ·
`PersonList` 114 · `VehicleList` 84 · `PrintRISlip` 54 · `PctPremiAllObjUSD` 23 ·
`StartPeriod` 21 · `EndPeriod` 21 · `PctPremiAllObj` 19 · `ReinsurerName` 17 ·
`OurRef` 12 · `ReinsurerID` 12 · `PCTPremiIDR` 8 · `TFAllObj` 3 · `PctShareAllObj` 3 ·
`UjrahAllObj` 3 · `RiCommAllObj` 3

> `[terverifikasi]` **`FacRetroList` mereplikasi seluruh pohon objek risiko.**
> Setiap baris retrosesi memuat `LocationList`, `CargoList`, `PersonList`, `VehicleList`,
> dan `CurrencyList` sendiri — bukan sekadar persentase share. Konsekuensi migrasi: volume
> data per case berlipat sebanyak jumlah retrosesionaris.

**`PersonList(i).`** — `ASMCoverage` 1.150 · `CoverageList` 63 · `pyFullName` 51 ·
`ASMStatusCoverage` 39 · `ASMBankData` 36 · `ASMRegNo` 36 · `FlagOldData` 30 ·
`TotalTSIPremiGrossList` 24 · `ASMGender` 18 · `FlagDelete` 18 · `ASMDateOfBirth` 17 ·
`AccumulationName` 16 · `ASMIDCard` 15 · `ASMJoinDate` 12 · (sisanya ≤ 9)

**`CargoList(i).`** — `CoverageList` 296 · `PolicyData` 97 · `FlagOldData` 33 ·
`TotalTSIPremiGrossList` 33 · `ConveyanceList`, `TradingList`, `GoodList`,
`FromRute`, `ToRute`, `PackingID`, … (2–4)

**`VehicleList(i).`** — `CoverageList` 878 · `TotalTSIPremiGrossList` 33 ·
`AccessoryList` 30 · `TotalTSIPremiSpreadRNM` 15 · `Occupation` 11 ·
`LicensePlate`, `ManufactureYear`, `VehicleCodeId`, `ChassisNumber`, `EngineNumber`,
`BrandName`, `Model`, `Type`, `ColorName`, `TransmissionType` (3–7)

**Perintah audit** (ganti `$keys` sesuai daftar yang mau dihitung)

```powershell
$files = Get-ChildItem -Path "D:\migrasi\RNM\NB FacIn","D:\migrasi\RNM\RNW Fac In","D:\migrasi\RNM\Endorsment Fac In" -Recurse -File -Filter *.xml
$keys = @('LocationList','CurrencyList','FacRetroList','PersonList','CargoList','VehicleList')
$tal=@{}; $rxs=@{}
foreach ($k in $keys) { $tal[$k]=@{}; $rxs[$k]=[regex]($k+'\(([^)]*)\)\.([A-Za-z][A-Za-z0-9_]*)') }
foreach ($f in $files) { $t=[System.IO.File]::ReadAllText($f.FullName)
  foreach ($k in $keys) { foreach ($m in $rxs[$k].Matches($t)) { $v=$m.Groups[2].Value
      if ($tal[$k].ContainsKey($v)) { $tal[$k][$v]++ } else { $tal[$k][$v]=1 } } } }
foreach ($k in $keys) { Write-Output "===== $k ====="
  $tal[$k].GetEnumerator()|Sort-Object Value -Descending|Select-Object -First 30 |
    ForEach-Object { "{0,7}  {1}" -f $_.Value, $_.Key } }
```

### 2.6 Lapis dalam — `Property`, `CoverageList`, `SpreadingList`, `FacOutObjectList`

**`Property.`** (page tunggal, **bukan** list — pola `Property(i).` tidak pernah muncul,
0 kecocokan)

| Anak | Rujukan |
| --- | --- |
| `RiskLocation` | 2.867 |
| `PropertyItemList` | 1.768 |
| `ObjectNo` | 382 |
| `RoadName` | 253 |
| `TotalTSIPremiSpreadRNM` | 231 |
| `SurroundingRisk` | 214 |
| `Province` | 176 |
| `TotalTSIPremiGrossList` | 147 |
| `OccupationList` | 141 |
| `ListCauseOfLossClaim` | 132 |
| `IsTopRisk` | 131 |
| `ObjectName` | 128 |
| `RoadType` | 117 |
| `AlmRiskID` | 113 |
| `BuildingNo` | 111 |
| `Country` | 71 |
| `BuildingConstruction` | 56 |
| `TotalTSIList` | 54 |
| `ObjectNoFacIn` | 42 |
| `ObjectType` | 34 |
| `MultiCOBOldID` | 30 |

**`RiskLocation.`** — `OccupationList` 1.216 · `Earthquake` 948 · `AnekaList` 509 ·
`Riot` 393 · `Tsunami` 393 · `Flood` 393 · `ASMZipCode` 337 · `ASMAddress` 298 ·
`ASMCity` 177 · `ASMDistrict` 159 · `ASMRW` 153 · `CityName` 48 · `DistrictName` 48 ·
`RWName` 21 · `Territory` 21

**`PropertyItemList(i).`** — `CoverageList` 947 · `Currency` 76 · `IsOldData` 63 ·
`ItemType` 42 · `TSIObjectItem` 42 · `PctAdjustOther` 41 · `FlagDelete` 39 ·
`IsAdjustableFlag` 23 · `PctAdjust1` 18 · `CurrencyID` 14 · `ItemTypeID` 14 ·
`FlagNetRate` 13 · `TotalNetRate` 13 · `TotalGrossPremi` 9 ·
`TSIObjectItemOld` 3 · `TotalPremiumNusantaraReOld` 3 · `TotalGrossPremiOld` 3

**`OccupationList(i).`** — `AnekaList` 1.281 · `IsOldData` 42 · `OccupationId` 33 ·
`OccupationName` 28 · `TableOfLimit` 21 · `FlagDelete` 18 · `BondList` 18 ·
`CashList` 12 · `CARList` 6 · `PlanList` 2

**`AnekaList(i).`** — `CoverageList` 1.224 · `Currency` 87 · `TSI` 63 · `IsOldData` 54 ·
`MarineHull` 42 · `VehicleHE` 19 · `AviationHull` 18 · `FlagDelete` 15 · `Section` 9 ·
`ObjectName` 8 · `TSIOld` 6 · `LimitOfLiability` 3 · `GeographicalLimits` 3 · `NoOfTree` 3 ·
`AreaHectar` 3

**`CoverageList(i).`** — daun perhitungan premi

| Anak | Rujukan | | Anak | Rujukan |
| --- | --- | --- | --- | --- |
| `FacOutObjectList` | 1.486 | | `NominalShareOffered` | 51 |
| `SpreadingList` | 1.133 | | `PremiumRetro` | 48 |
| `Coverage` | 323 | | `PlanList` | 48 |
| `AdditionalCoverage` | 298 | | `CoverageNote` | 42 |
| `DeductibleList` | 182 | | `TSINusantaraRe` | 40 |
| `TSI` | 147 | | `TypeOfDiscount` | 33 |
| `Currency` | 144 | | `PercentFacOut` | 33 |
| `OutGoList` | 128 | | `RateLifeAverage` | 30 |
| `Premium` | 121 | | `FlagDelete` | 24 |
| `IsOldData` | 117 | | `NetRate` | 23 |
| `FacOutTSI` | 117 | | `OldCoverage` | 63 |
| `Discount` | 110 | | `IndexCoverage` | 60 |
| `Rate` | 73 | | `TSILiability` | 57 |
| `LayerList` | 69 | | `DiscountPercentage` | 56 |

**`ASMCoverage(i).`** (padanan `CoverageList` pada cabang `PersonList`) —
`FacOutObjectList` 336 · `PlanList` 234 · `SpreadingList` 203 · `Premium` 138 ·
`Coverage` 120 · `Currency` 60 · `TSI` 57 · `IsOldData` 54 · `Loading` 36 ·
`Discount` 31 · `FacOutTSI` 24 · `TempPremium` 24 · `NominalShareOffered` 24 ·
`PremiRp` 24 · `OldPlanPremium` 12 · `TSIOld` 6 · `PremiumOld` 6

**`SpreadingList(i).`** (daun sebaran treaty) — `SharePercentage` 426 · `TSISpreaded` 356 ·
`TreatyType` 274 · `PremiumSpreaded` 194 · `IsOldData` 192 · `TreatyName` 69 ·
`FlagDelete` 63 · `TSIGrossSpreaded` 36 · `ClaimEstimation` 18 · `ClaimSpreaded` 18 ·
`ClaimAmountIDR` 6

**`FacOutObjectList(i).`** (daun retrosesi per objek) — `ShareOffered` 422 ·
`PercentOffered` 366 · `ObjectPremi` 299 · `Rate` 291 · `RiComm` 249 · `commision` 210 ·
`TF` 3 · `Ujrah` 3

**`DeductibleList(i).`** — `Descriptions` 29 · `Coverage` 28 · `ClaimCategory` 28 ·
`Currency` 23 · `DeductibleType` 21 · `BasisType` 21 · `Type` 21 · `Value` 21

**`PlanList(i).`** — `Name` 98 · `PremiumList` 81 · `BenefitList` 42 · `Limit` 25 ·
`CoverageType` 16 · `ClauseList` 9

### 2.7 `PolicyData.Payment` — sub-agregat uang

| Anak | Rujukan | | Anak | Rujukan |
| --- | --- | --- | --- | --- |
| `Premium` | 682 | | `Installment` | 145 |
| `ListInstallment` | 668 | | `NetPremium` | 132 |
| `RICommision` | 591 | | `Deduction2` | 113 |
| `PctBrokerageFee` | 361 | | `Diskon` | 108 |
| `Commision` | 359 | | `EDMOldPremi` | 92 |
| `PctDeduction2` | 273 | | `EDMNewPremi` | 85 |
| `BrokerageFee` | 246 | | `EDMOldCommision` | 72 |
| `PPh` | 195 | | `AdminFee` | 48 |
| `EDMPremiMenjadi` | 186 | | `StampPolicy` | 45 |
| `EDMOldPayment` | 184 | | `StampReceipts` | 39 |
| `PPN` | 183 | | `OriginalPremium` | 36 |

`[terverifikasi]` `Payment` membawa **pasangan sebelum/sesudah khusus endorsement**:
`EDMOldPremi` / `EDMNewPremi` / `EDMPremiMenjadi` / `EDMOldPayment` / `EDMOldCommision`.
Ini padanan di sisi aplikasi untuk pasangan kolom `_MENJADI` / `_SELISIH` di Oracle
(lihat `02-skema-oracle.md` §4.1).

**Perintah audit**

```powershell
$files = Get-ChildItem -Path "D:\migrasi\RNM\NB FacIn","D:\migrasi\RNM\RNW Fac In","D:\migrasi\RNM\Endorsment Fac In" -Recurse -File -Filter *.xml
$rx=[regex]'(?<![A-Za-z0-9_])Payment\.([A-Za-z][A-Za-z0-9_]*)'
$tal=@{}
foreach($f in $files){ $t=[System.IO.File]::ReadAllText($f.FullName)
  foreach($m in $rx.Matches($t)){ $v=$m.Groups[1].Value
    if($tal.ContainsKey($v)){$tal[$v]++}else{$tal[$v]=1} } }
$tal.GetEnumerator()|Sort-Object Value -Descending|Select-Object -First 25 |
  ForEach-Object { "{0,6}  {1}" -f $_.Value,$_.Key }
```

---

## 3. Properti paling sering dirujuk dan nilai literal yang diuji terhadapnya

### 3.1 Properti daun terpopuler di bawah `OfferFacIn`

**513 nama daun unik**. 30 teratas:

| # | Daun | Rujukan | | # | Daun | Rujukan |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | `BusinessOldId` | 1.188 | | 16 | `TSISpreaded` | 336 |
| 2 | `StatusBusiness` | 1.076 | | 17 | `PctBrokerageFee` | 309 |
| 3 | `Type` | 954 | | 18 | `FacRetroList` | 298 |
| 4 | `PercentShare` | 807 | | 19 | `TSI` | 296 |
| 5 | `LocationList` | 709 | | 20 | `TeamGroup` | 290 |
| 6 | `StartDateTime` | 698 | | 21 | `PctDeduction2` | 261 |
| 7 | `EdmType` | 646 | | 22 | `ProrateEDMEnd` | 249 |
| 8 | `IsOldData` | 576 | | 23 | `Coverage` | 233 |
| 9 | `CurrencyList` | 563 | | 24 | `FlagDelete` | 222 |
| 10 | `EndDateTime` | 498 | | 25 | `BusinessCode` | 210 |
| 11 | `RICommision` | 485 | | 26 | `MarketingCode` | 205 |
| 12 | `IsGroup` | 399 | | 27 | `ShareRNML` | 200 |
| 13 | `BusinessType` | 396 | | 28 | `PolicyNo` | 189 |
| 14 | `ProRatePercent` | 361 | | 29 | `TreatyType` | 177 |
| 15 | `PersonList` | 356 | | 30 | `PremiumSpreaded` | 177 |

**Perintah audit**

```powershell
$files = Get-ChildItem -Path "D:\migrasi\RNM\NB FacIn","D:\migrasi\RNM\RNW Fac In","D:\migrasi\RNM\Endorsment Fac In" -Recurse -File -Filter *.xml
$rx=[regex]'OfferFacIn((?:\.[A-Za-z][A-Za-z0-9_]*(?:\([^)]*\))?)+)'
$tal=@{}
foreach($f in $files){ $t=[System.IO.File]::ReadAllText($f.FullName)
  foreach($m in $rx.Matches($t)){ $leaf=(($m.Groups[1].Value) -split '\.')[-1]
    if($leaf -match '\('){ continue }
    if($tal.ContainsKey($leaf)){$tal[$leaf]++}else{$tal[$leaf]=1} } }
Write-Output ("daun unik: " + $tal.Count)
$tal.GetEnumerator()|Sort-Object Value -Descending|Select-Object -First 45 |
  ForEach-Object { "{0,6}  {1}" -f $_.Value,$_.Key }
```

### 3.2 Nilai literal yang diuji — dari rule `When`

Korpus memuat **1.863** blok `<pyConditionValue1String>` non-kosong; **1.590** di antaranya
berbentuk `<properti> <operator> <nilai>` dan dapat diurai otomatis.

**Perintah audit**

```powershell
$files = Get-ChildItem -Path "D:\migrasi\RNM\NB FacIn","D:\migrasi\RNM\RNW Fac In","D:\migrasi\RNM\Endorsment Fac In" -Recurse -File -Filter *.xml
$rx=[regex]'<pyConditionValue1String>([^<]*)</pyConditionValue1String>'
$rows=New-Object System.Collections.ArrayList
foreach ($f in $files) { $t=[System.IO.File]::ReadAllText($f.FullName)
  foreach ($m in $rx.Matches($t)) { [void]$rows.Add([pscustomobject]@{File=$f.FullName;Cond=$m.Groups[1].Value}) } }
Write-Output ("total kondisi: " + $rows.Count)
$p=[regex]'^\s*(?<lhs>[A-Za-z0-9_.()\[\]]+)\s*(?<op>=|==|!=|>=|<=|>|<)\s*(?<rhs>.+?)\s*$'
$parsed = foreach ($r in $rows) { $m=$p.Match($r.Cond)
  if ($m.Success) { [pscustomobject]@{ Prop=($m.Groups['lhs'].Value -replace '^pyWorkPage\.','')
                                       Op=$m.Groups['op'].Value; Val=$m.Groups['rhs'].Value.Trim(); File=$r.File } } }
Write-Output ("kondisi terparse: " + ($parsed|Measure-Object).Count)
$parsed | Group-Object Prop | Sort-Object Count -Descending | Select-Object -First 40 |
  ForEach-Object { "{0,5}  {1}  ->  {2}" -f $_.Count, $_.Name,
    (($_.Group | Select-Object -ExpandProperty Val -Unique | Sort-Object) -join ' | ') }
```

**Properti dengan literal terbanyak** (nilai literal dikutip apa adanya; **artinya belum
terverifikasi** — lihat `03-batas-pengetahuan.md`):

| Properti | Uji | Nilai literal yang terbaca |
| --- | --- | --- |
| `Quotation.BusinessCode` | 236 | 64 nilai berbeda: `"05" "06" "07" "15" "17" "18" "19" "21" "24" "27" "32" "90" "91" "B7" "SK" "SL" "SQ" "SV" "SX" "T9"` + blok lima-digit `"100xx"–"102xx"` |
| `OfferFacIn.QuotationData.BusinessOldId` | 210 | 70 nilai: `30 36–52 79–88 95–99` dan kode huruf `A1–A9 B4–B9 C1–C9 D1–D9 F1–F4 F6` |
| `.Quotation.BusinessCode` (relatif) | 181 | 34 nilai lain lagi: `"41" "61" "63" "66" "68" "72" "75" "78"–"89" "99" "S7" "S8" "T1" "TA"` + `"100xx"` |
| `Quotation.BusinessType` | 132 | 35 nilai nama COB: `"AllRisk" "Aneka" "AviationHull" "BillboardNeonSyariah" "Boiler" "Bonding" "BondingKBG" "Burglary" "Car" "CIS" "CIT" "ContractorsPM" "CustomBond" "Ear" " ElectronicEquipment " "Exclusion" "Fidelity" "FireStyle1" "FireStyle2" "Glass" "GolfInsurance" "GrowingTrees" "HE" "KPR" "LandRig" "Liability" "Maintenance" "MarineCargo" "MarineHull" "MBD" "MBUCar" "OilGas" "PA" "Travel" "Workmen"` |
| `.CoverageID` | 48 | `"100846" "100848"–"100850" "100855" "100856" "100874"–"100879"` |
| `Quotation.BusinessOldId` | 48 | `"L1"–"L16"` (blok huruf L) |
| `.Coverage` | 48 | `"100143" "100847"–"100850" "100855"–"100858" "100861" "100867" "100869" "100872" "100875" "100878" "100879"` |
| `OfferFacIn.QuotationData.Type` | 43 | `0 1 2 3 4 5 6 7 8 9 11 12` (**10 tidak pernah muncul**) |
| `OfferFacIn.QuotationData.StatusBusiness` | 40 | `3` |
| `OfferFacIn.QuotationData.EdmType` | 37 | `4` |
| `LetterNo` | 34 | `"" "DEPHEADUNDERWRITER" "DEPTHEADUWLIFE" "DIREKTURMARKETING" "DIREKTURTEKNIK" "JUW_A" "KADIVFACULTATIVE" "KADIVFINANCIAL" "KADIVTEKNIK" "MANAGERTEKNIK" "SENIORUW" "TREATYINDEPTHEAD" "UNDERWRITER"` |
| `.TabPosition` | 29 | `"Cedant" "Coverage" "Deduction" "LossRecord" "Object" "Payment" "PaymentLife" "ScoringRisk" "Spreading" "SpreadingLife"` |
| `.DeductibleType` | 27 | `"100413" "100414" "100416"–"100420"` (**100415 tidak pernah muncul**) |
| `OfferFacIn.QuotationData.TeamGroup` | 24 | `1 2 3 4` (nilai `5` muncul di luar rule When — §3.4) |
| `Quotation.SobLeader1` | 21 | 7 ID numerik 8-digit |
| `PositionNote` | 18 | `"ReasFacInDepHeadUnderwriting" "ReasFacInJuniorUnderwriting" "ReasFacInJuniorUnderwritingA" "ReasFacInSeniorUnderwriting" "ReasFacInUnderwriting" "ReasFacInUnderwritingFinancial"` |
| `OfferFacIn.IsB2B` | 15 | `""` dan `"ASM"` |
| `.IsCedingConfirm` | 12 | `"accepted"`, `Offer`, `Policy` |
| `pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName` | 10 | `ReasFacInDirector ReasFacInGroupLeader ReasFacInMarketing ReasFacInUnderwriting` |
| `Quotation.StatusBusiness` | 9 | `1 2 3` |
| `.PropertyItemGroup` | 9 | `"STOCK" "STOCK(S)" "STOCKS"` |
| `.Quotation.BusinessFac` | 7 | `"F"` `"T"` |
| `.Policy.TypeOfCoins` | 6 | `"F"` `0` |
| `.Policy.TypeOfPolicy` | 6 | `"02"` |
| `.CustomerType` | 5 | `"1"` `"2"` |

### 3.3 Guard berbasis identitas — dicatat sebagai mekanisme, bukan nilai

`[terverifikasi]` Korpus memuat rule `When` yang bercabang atas **identitas orang**, bukan peran:

| Mekanisme | Jumlah identitas | Lokasi |
| --- | --- | --- |
| `OperatorID.pyUserIdentifier = <identitas>` | **3 identitas operator** | rule `When` (9 pemakaian kondisi) |
| `pxCreateOperator = <identitas>` | **4 identitas operator** | rule `When` (8 pemakaian kondisi) |
| `.OfferFacIn.QuotationData.MarketingName = <nama>` | **3 nama marketing** | rule `When` (9 pemakaian kondisi) |
| `Quotation.SobLeader1 = <ID klien>` | **7 ID klien 8-digit** | rule `When` (21 pemakaian) |
| `Quotation.OldPolicyNo = <nomor polis>` | **2 nomor polis produksi** | rule `When` (6 pemakaian) |
| `PXADDEDBYID not in (…)` di SQL | **1 identitas operator + 1 literal `'System'`** | `NB FacIn/RDBList/GetHistoryAccPega_SQL.xml` |

**Nilai tidak disalin.** Semuanya **kandidat perbaikan**, keputusannya milik bisnis.

### 3.4 Enumerasi yang terbaca di luar rule `When`

Pemindaian seluruh korpus (bukan hanya `When`) atas pola `<Prop> = <nilai>`:

| Properti | Nilai yang pernah muncul | Catatan |
| --- | --- | --- |
| `QuotationData.Type` | 0(6) 1(3) 2(3) 3(108) 4(234) 5(9) 6(24) 7(147) 8(3) 9(3) 11(12) 12(22) | **10 kosong** |
| `EdmType` | 1 2 3 4 | dominan 4 |
| `ProRateType` | 3 4 (dan di-*set* ke `1`) | 1 dan 2 tidak pernah **diuji** |
| `StatusBusiness` | 1(24) 2(15) 3(447) | |
| `TeamGroup` | 1 2 3 4 **5** | nilai 5 hanya di luar `When` |
| `DeductibleType` | 100413 100414 100416 100417 100418 100419 100420 **dan** 4 7 9 | dua skema penomoran berbeda |
| `TreatyType` | 10001 10003 10004 10007 10011 10012 10014 10015 10019 10021 10022 10024–10026 10028 10184–10186 10218 1000036 1001414 **dan** `"FAC"` `"Local"` | campur ID numerik dan teks |
| `PolicyType` | 1 2 | |
| `CustomerType` | 1 2 | |
| `IsB2B` | `""` `"ASM"` | |
| `ProportionalType` | `"Proportional"` | hanya satu nilai terbaca |

**Perintah audit**

```powershell
$files = Get-ChildItem -Path "D:\migrasi\RNM\NB FacIn","D:\migrasi\RNM\RNW Fac In","D:\migrasi\RNM\Endorsment Fac In" -Recurse -File -Filter *.xml
$rx=[regex]'QuotationData\.Type\s*(?:=|==|!=)\s*"?(\d+)"?'
$tal=@{}
foreach($f in $files){ $t=[System.IO.File]::ReadAllText($f.FullName)
  foreach($m in $rx.Matches($t)){ $v=$m.Groups[1].Value
    if($tal.ContainsKey($v)){$tal[$v]++}else{$tal[$v]=1} } }
$tal.GetEnumerator()|Sort-Object {[int]$_.Key}|ForEach-Object{"Type={0}  n={1}" -f $_.Key,$_.Value}
```

---

## 4. Kedalaman bersarang

### 4.1 Jalur terdalam yang tertulis **harfiah** di korpus — `[terverifikasi]`

Ini bukan simpulan dari penamaan Section. Ini string yang ada apa adanya di dalam XML.

**Cabang PROPERTY / FIRE (7 lapis di bawah work page):**

```
pyWorkPage.OfferFacIn
  .LocationList(i)
  .Property
  .PropertyItemList(j)
  .CoverageList(k)
  .SpreadingList(l)
```

> `NB FacIn/Activity/CopyToAllLocSpreading_ACT.xml` baris 1996
> `<Page>pyWorkPage.OfferFacIn.LocationList(Param.IdxLoc).Property.PropertyItemList(Local.IdxObjItem).CoverageList(Local.IdxCov).SpreadingList(Local.IdxSpreading)</Page>`

Ada satu lapis tambahan opsional `LayerList` di antara `CoverageList` dan `SpreadingList`:

> `Endorsment Fac In/Activity/CopyToAllSpreadingFire_ACT.xml`
> `pyWorkPage.OfferFacIn.OldData.LocationList(…).Property.PropertyItemList(…).CoverageList(…).LayerList(Local.IdxLayer).SpreadingList`

**Cabang ANEKA / GOLF (8 lapis):**

```
pyWorkPage.OfferFacIn
  .LocationList(i)
  .Property
  .RiskLocation
  .OccupationList(j)
  .AnekaList(k)
  .CoverageList(l)
  .SpreadingList / .FacOutObjectList
```

> `NB FacIn/Activity/addDeductibleCargo_ACT.xml` baris 240
> `pyWorkPage.OfferFacIn.LocationList(1).Property.RiskLocation.OccupationList(1).AnekaList(1).CoverageList(1).Currency`

**Cabang RETROSESI (10 lapis) — kedalaman maksimum yang terbaca:**

```
pyWorkPage.OfferFacIn
  .FacRetroList(a)
  .LocationList(b)
  .Property
  .RiskLocation
  .OccupationList(c)
  .AnekaList(d)
  .CoverageList(e)
  .FacOutTSI
```

> `NB FacIn/Activity/CopyFacRetroAnekaGolf_ACT.xml` baris 2212
> `<PropertiesName>pyWorkPage.OfferFacIn.FacRetroList(local.counterFacRetro).LocationList(local.counterLocation).Property.RiskLocation.OccupationList(local.counterOccupation).AnekaList(local.counterAneka).CoverageList(1).FacOutTSI</PropertiesName>`

**Cabang PA / LIFE:** `pyWorkPage.OfferFacIn.PersonList(i).ASMCoverage(j).SpreadingList(k)`
(`Endorsment Fac In/Activity/CopyToAllSpreading_ACT.xml`).

**Cabang CARGO / MBU:** `pyWorkPage.OfferFacIn.CargoList(i).CoverageList(j).SpreadingList(k)` dan
`…VehicleList(i).CoverageList(j).SpreadingList(k)`
(`Endorsment Fac In/Activity/CopyToAllSpreadingCargoMBU_ACT.xml`).

**Perintah audit**

```powershell
Select-String -Path "D:\migrasi\RNM\NB FacIn\*\*.xml","D:\migrasi\RNM\RNW Fac In\*\*.xml","D:\migrasi\RNM\Endorsment Fac In\*\*.xml" `
  -Pattern 'LocationList\([^)]*\)\.Property\.PropertyItemList\([^)]*\)\.CoverageList\([^)]*\)\.SpreadingList' |
  Select-Object -First 3 Path,LineNumber
Select-String -Path "D:\migrasi\RNM\NB FacIn\*\*.xml" `
  -Pattern 'FacRetroList\([^)]*\)\.LocationList\([^)]*\)\.Property\.RiskLocation\.OccupationList' |
  Select-Object -First 3 Path,LineNumber
```

### 4.2 Bentuk yang sama dikonfirmasi ulang dari sisi Oracle — `[terverifikasi]`

Ini bukti independen: query Oracle membongkar JSON hasil serialisasi agregat dan
**menyebut jalur yang persis sama**.

> `NB FacIn/RDBList/GetSummaryRiskAccumPolis_Sql.xml`
> ```sql
> from pooldata.json_polis a,
>      json_table(data_json, '$.LocationList[*]' columns(
>        Address varchar2(10000) path '$.Property.RiskLocation.ASMAddress',
>        nested path '$.Property.PropertyItemList[*]' columns(
>          Currency varchar2(100) path '$.Currency',
>          nested path '$.CoverageList[*]' columns(
>            AccumulationCode varchar2(100) path '$.AccumulationCode',
>            Coverage         varchar2(100) path '$.CoverageNote',
>            nested path '$.SpreadingList[*]' columns(
>              TreatyType      varchar2(100) path '$.TreatyType',
>              TSISpreaded     varchar2(100) path '$.TSISpreaded',
>              PremiumSpreaded varchar2(100) path '$.PremiumSpreaded' )))))
> ```

Yang ini membuktikan sekaligus **tiga** hal:

1. Agregat `OfferFacIn` **diserialkan utuh** ke `JSON_POLIS.DATA_JSON` dengan nama properti
   yang identik (`LocationList`, `Property`, `RiskLocation`, `PropertyItemList`,
   `CoverageList`, `SpreadingList`).
2. `LocationList`, `PropertyItemList`, `CoverageList`, `SpreadingList` adalah **array**
   (`[*]`), `Property` dan `RiskLocation` adalah **objek tunggal**.
3. `TSISpreaded` dan `PremiumSpreaded` dibaca sebagai **`varchar2(100)`**, bukan `NUMBER`.
   Lihat §6.

### 4.3 Yang **hanya** dari penamaan Section — `[dugaan]`

Nilai `.TabPosition` (`"Object"`, `"Coverage"`, `"Spreading"`, `"Payment"`, `"Cedant"`,
`"Deduction"`, `"LossRecord"`, `"ScoringRisk"`, `"PaymentLife"`, `"SpreadingLife"`)
**menyiratkan** urutan tab layar lokasi → objek → coverage → spreading → payment.
Itu hanya nama tab; korpus tidak memuat rule yang membuktikan urutan itu mengikat data.
**Ini `[dugaan]`, bukan bukti struktur.**

---

## 5. Mekanisme before-image untuk endorsement

`[terverifikasi]` Korpus memuat **tiga mekanisme berbeda** yang semuanya menyimpan "nilai
sebelumnya". Ketiganya **tidak boleh disamakan**.

### 5.1 Mekanisme A — riwayat versi di database (baris bertambah)

**Bukan** salinan kerja. Ini versioning polis di Oracle.

`JSON_POLIS` menyimpan **satu baris per versi polis**, dibedakan oleh `PRODKE`.
Nomor versi dihitung sebagai **cacah baris dikurangi satu**:

> `NB FacIn/RDBList/GetProdKeEDM_SQL.xml`
> ```sql
> select PRODKE as HASIL2 from json_polis
> where NOPOLIS = {InputData.CARI17}
>   and PRODKE = (SELECT COUNT(NOPOLIS)-1 FROM JSON_POLIS b WHERE NOPOLIS={InputData.CARI17})
> ```

Versi sebelumnya diambil utuh sebagai JSON:

> `NB FacIn/RDBList/GetOldDataEDM_SQL.xml`
> `SELECT a.DATA_JSON AS HASIL1, a.IDPEGA AS HASIL2 FROM JSON_POLIS a WHERE NOPOLIS={TempPolis.CARI4} and PRODKE={TempPolis.CARI22}`

**Sifat:** append-only, lintas-case, permanen. Detail di `02-skema-oracle.md` §4.2.

### 5.2 Mekanisme B — salinan kerja agregat: `OfferFacIn.OldData`

`[terverifikasi]` `OldData` adalah **halaman di dalam agregat yang mencerminkan agregat itu
sendiri**. Anak-anaknya adalah sub-agregat yang sama:

| Anak `OldData` | Rujukan |
| --- | --- |
| `PolicyData` | 508 |
| `LocationList` | 258 |
| `QuotationData` | 239 |
| `FacRetroList` | 207 |
| `CurrencyList` | 110 |
| `PersonList` | 52 |
| `VehicleList` | 50 |
| `CargoList` | 35 |
| `IDNewBisnis` | 35 |
| `OfferFacIn` | 56 |
| `OldData` | 22 |
| `ProRatePercent`, `TotalTSINusaRe`, `TotalTSITopRisk`, `PercentShare`, `IsProRate`, `ProRateType`, `MaxTreatyCapacity`, `MaxPctTreatyCapacity`, `InwardScale`, `AdditionalCapital`, `BinderRNM`, `TreatyXOLList`, `PolicyMasterNumber`, `IsSpecialAcceptance`, `Parameters`, … | 1–8 |

Total rujukan `OfferFacIn.OldData` = **1.565** dalam **175 berkas**.
Kemunculan `OldData.OldData` (22) dan `OldData.OfferFacIn` (56) `[pertanyaan terbuka]` —
bisa jadi penyalinan berlapis, bisa jadi cacat penulisan ekspresi.

Jalur bersarang penuh ada di dalam `OldData` juga — `[terverifikasi]`:

> `NB FacIn/Activity/CopyToAllSpreadingFire_ACT.xml`
> `pyWorkPage.OfferFacIn.OldData.LocationList(Local.IdxLocation).Property.PropertyItemList(Local.IdxObjItem).CoverageList(Local.IdxCov).SpreadingList`

`OldData` dipakai langsung sebagai pembanding periode saat hitung pro-rata endorsement:

> `Endorsment Fac In/DataTransform/CountPremiEDM_DT.xml` baris 1138 / 1287 / 1317 / 1345
> `@DateTimeDifference(pyWorkPage.OfferFacIn.PolicyData.StartDateTime, pyWorkPage.OfferFacIn.OldData.PolicyData.StartDateTime,"D") == 0 …`
> `@DateTimeDifference(pyWorkPage.OfferFacIn.OldData.PolicyData.StartDateTime, pyWorkPage.OfferFacIn.OldData.QuotationData.EdmDate,"D")`

**Sifat:** hidup di dalam satu case, dipakai untuk menghitung selisih.

### 5.3 Mekanisme C — before-image **di dalam baris** dengan akhiran `Old`

`[terverifikasi]` Berbeda dari A dan B: di sini nilai lama disimpan **di baris yang sama**
sebagai field bersebelahan.

> `Endorsment Fac In/Activity/SetOldData.xml` (rule `SetOldData`)
> ```
> pyWorkPage.OfferFacIn.LocationList(idxlocation).Property.PropertyItemList(idx).TSIObjectItemOld
>     = @If(.TSIObjectItem!="", .TSIObjectItem, 0)
> …TotalGrossPremiOld            = @If(.TotalGrossPremi!="", .TotalGrossPremi, 0)
> …TotalPremiumNusantaraReOld    = @If(.TotalPremiumNusantaraRe!="", .TotalPremiumNusantaraRe, 0)
> …TotalTSIPremiGrossList(idx3).TSIOld     = @If(.TSI!="", .TSI, 0)
> …TotalTSIPremiGrossList(idx3).PremiumOld = @If(.PremiumOld!="", .Premium, 0)
> …TotalTSIPremiGrossList(idx3).RateOld    = @If(.RateOld!="", .Rate, 0)
> …TotalTSIList(idx).TSIOld                = @If(.TSI!="", .TSI, 0)
> …CedingCedantList(idxcedant).CurrencyList(idx).TSIOld     = @If(.TSI!="", .TSI, 0)
> …CedingCedantList(idxcedant).CurrencyList(idx).PremiumOld = @If(.Premium!="", .Premium, 0)
> …LocationList(idxlocation).Property.RiskLocation.AnekaList(idx).TSIOld = @If(.TSI!="", .TSI, 0)
> ```

> **Temuan yang harus dibawa ke bisnis, bukan diperbaiki diam-diam.**
> Dua baris di atas menguji **field tujuan**, bukan field sumber:
> `PremiumOld = @If(.PremiumOld!="", .Premium, 0)` dan `RateOld = @If(.RateOld!="", .Rate, 0)`.
> Tujuh baris lainnya menguji field sumber (`@If(.TSI!="", .TSI, 0)`).
> Akibat yang **[dugaan]**: pada endorsement pertama (`PremiumOld` masih kosong)
> `PremiumOld` dan `RateOld` akan di-set `0`, bukan nilai lama.
> Ini **kandidat perbaikan** — migrasi harus mereplikasi perilaku apa adanya sampai bisnis
> memutuskan, agar rekonsiliasi paralel run tidak menghasilkan selisih tak terjelaskan.

Guard masuk `SetOldData` `[terverifikasi]`:
`IsLife exit` · `!IsEDM exit` · `@Utilities.SizeOfPropertyList(.OfferFacIn.LocationList) > 100 exit`.
Batas **100 lokasi** itu literal di rule.

**Bendera pendamping:** `IsOldData` (1.767 rujukan / 133 berkas) dan
`FlagOldData` (514 rujukan / 56 berkas) menandai **baris** mana yang berasal dari versi lama.
`[pertanyaan terbuka]` beda tegas `IsOldData` vs `FlagOldData` tidak terbaca; keduanya muncul
di kelas yang berbeda (`IsOldData` di cabang Property/Coverage, `FlagOldData` di cabang
Person/Cargo/Vehicle) tetapi tidak ada rule yang menyatakan semantiknya.

**Perintah audit**

```powershell
$files = Get-ChildItem -Path "D:\migrasi\RNM\NB FacIn","D:\migrasi\RNM\RNW Fac In","D:\migrasi\RNM\Endorsment Fac In" -Recurse -File -Filter *.xml
foreach($pat in @('OfferFacIn\.OldData','(?<![A-Za-z0-9_])IsOldData(?![A-Za-z0-9_])','(?<![A-Za-z0-9_])FlagOldData(?![A-Za-z0-9_])')){
  $rx=[regex]$pat; $c=0; $fc=0
  foreach($f in $files){ $t=[System.IO.File]::ReadAllText($f.FullName)
    $n=$rx.Matches($t).Count; if($n -gt 0){$fc++; $c+=$n} }
  Write-Output ("{0,-48} {1,6} rujukan / {2,4} berkas" -f $pat,$c,$fc) }
```

### 5.4 Mekanisme D — pasangan kolom di Oracle

Bukan mekanisme aplikasi, tetapi **hilir** dari ketiga mekanisme di atas:
kolom `*_MENJADI` (nilai sesudah) dan `*_SELISIH` (delta) berdampingan di tabel produksi.
22 kolom, rinciannya di `02-skema-oracle.md` §4.1.

### 5.5 Ringkasan perbedaan — jangan dicampur saat migrasi

| | A — riwayat versi | B — salinan kerja agregat | C — before-image dalam baris |
| --- | --- | --- | --- |
| **Wadah** | baris `JSON_POLIS` per `PRODKE` | halaman `OfferFacIn.OldData` | field berakhiran `…Old` di baris yang sama |
| **Umur** | permanen, lintas case | selama case EDM hidup | selama case EDM hidup |
| **Granularitas** | seluruh polis | seluruh agregat | per baris objek/coverage/currency |
| **Diisi oleh** | prosedur `POOLDATA.INSERTJSONPOLIS` | belum terverifikasi (lihat catatan) | rule `SetOldData` |
| **Dipakai untuk** | mengambil polis versi n−1 | menghitung selisih periode & premi | menampilkan lama vs baru per baris |

`[pertanyaan terbuka]` **Titik pengisian `OfferFacIn.OldData` tidak terbaca.**
Korpus memperlihatkan `OldData` **dibaca** di 175 berkas dan **ditulis per-cabang** di
beberapa Activity (mis. `Endorsment Fac In/Activity/GetOldDataRetro_ACT.xml` menulis
`pyWorkPage.OfferFacIn.OldData.FacRetroList`), tetapi **tidak ada satu pun `Page-Copy`
dengan `<CopyInto>` menuju `OfferFacIn.OldData`** di seluruh korpus (0 kecocokan).
Bagaimana `OldData` terisi utuh dari `DATA_JSON` versi sebelumnya **belum terverifikasi**.

```powershell
Select-String -Path "D:\migrasi\RNM\*\Activity\*.xml" -Pattern '<CopyInto>[^<]*OfferFacIn\.OldData[^<]*</CopyInto>' | Measure-Object
```

---

## 6. Uang di dalam model domain

`[terverifikasi]` Nilai uang **tidak pernah boleh menjadi `float`** dalam implementasi Go.
Bukti langsung dari korpus:

1. **Uang dibaca dari JSON sebagai teks.**
   `NB FacIn/RDBList/GetSummaryRiskAccumPolis_Sql.xml` mendeklarasikan
   `TSISpreaded varchar2(100) path '$.TSISpreaded'` dan
   `PremiumSpreaded varchar2(100) path '$.PremiumSpreaded'`.

2. **Uang masuk ke Oracle sebagai string berkoma desimal, lalu dikonversi.**
   `NB FacIn/RDBList/InsertTreatyProduction_Sql.xml` membungkus hampir setiap nilai uang dengan
   `To_number(Replace({Datain.CARIxx}, ',', '.'))` — artinya nilai yang datang dari aplikasi
   memakai **koma** sebagai pemisah desimal.

3. **Arah konversi tidak konsisten antar tabel.**
   `NB FacIn/RDBList/InsertOfferProduction_Sql.xml` justru memakai
   `Replace({Datain.CARI34}, '.', ',')` — mengubah titik **menjadi** koma, **tanpa** `To_number`.
   Begitu pula `InsertIntoFacinSPreadLife_Sql.xml` (`Replace(…,'.',',')`) sementara versi
   `…Monthly_Sql.xml` memakai `To_number(Replace(…,',','.'))`.
   `[pertanyaan terbuka]` apakah kolom di `FACINOFFER` / `FACINSPREADLIFE` bertipe teks.

4. **Pembagian premi memakai pembagi `1e11` dengan 20 digit presisi.**
   `NB FacIn/Activity/CountPremi_ACT.xml` baris 4224:
   `@Math.divide((.TSI*.Rate*Local.Prorate*.IndemnityPercentage*.PctAdjustment*Local.LossLimit), @toDecimal("100000000000"), 20)`

5. **Ambang eksposur dinyatakan sebagai literal string lalu dikonversi.**
   `NB FacIn/Activity/CheckSpreadingProtect_ACT.xml`:
   `pyWorkPage.OfferFacIn.RNMShareTopLoc > @toDecimal("50000000000")` dan
   `pyWorkPage.OfferFacIn.TotalTSINusaRe > toDecimal("30000000000")`.

**Perintah audit**

```powershell
$files = Get-ChildItem -Path "D:\migrasi\RNM\NB FacIn","D:\migrasi\RNM\RNW Fac In","D:\migrasi\RNM\Endorsment Fac In" -Recurse -File -Filter *.xml
$rx=[regex]'"[0-9]{9,}"'
$tal=@{}
foreach($f in $files){ $t=[System.IO.File]::ReadAllText($f.FullName)
  foreach($m in $rx.Matches($t)){ $k = $m.Value + '|' + (Split-Path $f.FullName -Leaf)
    if($tal.ContainsKey($k)){$tal[$k]++}else{$tal[$k]=1} } }
$tal.GetEnumerator()|Sort-Object Value -Descending|Select-Object -First 25 |
  ForEach-Object { "{0,5}  {1}" -f $_.Value,$_.Key }
```

`[pertanyaan terbuka]` **Mata uang tidak melekat pada nilai.**
`CoverageList(i).Currency`, `PropertyItemList(i).Currency`, dan `CurrencyList(i).Name`
ada sebagai properti terpisah, sedangkan `TSI`, `Premium`, `TSISpreaded`,
`PremiumSpreaded` adalah skalar telanjang. Model target **wajib** memasangkan keduanya
(`Money{Amount decimal, Currency string}`) karena korpus tidak menjamin pemasangan itu.

---

## 7. Ringkasan untuk perancangan model Go

| Temuan | Label | Konsekuensi |
| --- | --- | --- |
| Agregat tunggal `OfferFacIn` menampung seluruh case | terverifikasi | satu *aggregate root*, bukan banyak tabel normal |
| Kedalaman 10 lapis pada cabang retrosesi | terverifikasi | struktur bersarang harus dipertahankan, bukan diratakan |
| `FacRetroList` mereplikasi seluruh pohon objek | terverifikasi | volume data per case berlipat; perlu strategi indeks |
| Agregat diserialkan ke `JSON_POLIS.DATA_JSON` dengan nama properti identik | terverifikasi | nama field Go/JSON **tidak boleh diubah** |
| `PRODKE = COUNT(NOPOLIS) − 1` | terverifikasi | versioning append-only, bukan update in-place |
| Tiga mekanisme "nilai lama" berbeda (A/B/C) | terverifikasi | jangan disatukan jadi satu konsep |
| `PremiumOld`/`RateOld` menguji field tujuan | terverifikasi | replikasi apa adanya; catat sebagai kandidat perbaikan |
| Tidak ada rule Property di korpus | terverifikasi | seluruh tipe data = `belum terverifikasi` |
| `Quotation` vs `QuotationData` ganda | terverifikasi (keberadaan), terbuka (otoritas) | **memblokir** desain model kanonik |
| Uang sebagai string berkoma desimal, konversi tidak konsisten | terverifikasi | `decimal` + mata uang eksplisit; JSON API pakai string desimal |
| Enumerasi `Type` bolong di 10, `DeductibleType` bolong di 100415 | terverifikasi | celah = temuan, harus dikonfirmasi bisnis |
| Guard berbasis identitas orang di 6 tempat | terverifikasi | kandidat perbaikan; nilai tidak disalin |

Lanjutkan ke `02-skema-oracle.md` dan `03-batas-pengetahuan.md`.
