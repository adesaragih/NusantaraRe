# Audit 03 — Struktur `FacRetro` dan populasi Fac Out

> **Diukur ulang 21 September 2026.** Sumber angka bagi `05-tickets\facout\F05-empat-struktur-staging.md`
> dan `F11-jalur-facout-endorsement.md`.

---

## 0. Uji silang wajib — dijalankan LEBIH DULU

⛔ Hasil banding 23 tag di §3 **tidak boleh dipakai** sebelum uji ini lulus.

| Berkas | Harus | Hasil |
| --- | --- | --- |
| `When\FlagOldData` | **BEDA** | BEDA ✅ |
| `When\IsAdmin` | identik | identik ✅ |
| `When\IsBonding` | identik | identik ✅ |
| `When\IsAddButton` | identik | identik ✅ |
| `DataPage\D_AnekaList` | identik | identik ✅ |

**LULUS 5/5.** Validasi kedua: metode ini mereproduksi angka resmi K-043 — NB↔EDM **1.353 identik /
354 berbeda** (lihat `07-hitung-berkas-dan-drift.md`).

---

## 1. `FacRetroDetails` dan keempat propertinya

`[terverifikasi]` **117 kemunculan** string `FacRetroDetails`, **39 per folder, identik di ketiganya**
(22 Activity + 6 Harness + 11 Section). **Nol kemunculan di `DDL\`** → struktur clipboard/JSON, **bukan
tabel Oracle**.

`[terverifikasi]` **Tepat empat properti** di bawah `FacRetroDetails`:

| Properti | Ditulis oleh (`<PropertiesName>`) | Dibaca oleh |
| --- | --- | --- |
| `.BackUpStatus` | `OfferFacOut_PreAct` · `SetPropertyToOfferFacIn_Act` · `GetRISlipDataFromDB_Act` | `CheckProtectFacout_Act` `<pyStepsPreCondParamsWhen>` `==3` **telanjang** · `OfferFacOut_PostAct` `=="3"` **berkutip** · `Harness\ViewLetter` · `ViewOfferStatus` · `ViewOldOffer` · `Section\FacOutPrintRISlipSectionInside` · `FacultativeLetter` · `OfferStatusFacOut` · `OldOfferStatusFacOut` |
| `.FacNo` | `OfferFacOut_PreAct` · `SetPropertyToOfferFacIn_Act` | — |
| `.AdditionalInfo` | `OfferFacOut_PreAct` · `SetPropertyToOfferFacIn_Act` | `Section\FacOutPrintRISlipSectionInside` |
| `.DocumentPosition` | `OfferFacOut_PreAct` · `SetPropertyToOfferFacIn_Act` | `Harness\ViewLetter` · `Section\FacOutPrintRISlipSectionInside` · `FacultativeLetter` |

⚠️ **Nilai yang sama dibandingkan dua gaya** — `==3` telanjang dan `=="3"` berkutip. Kandidat K-046,
diport apa adanya.

### 1.1 ⛔ Nama properti itu TIDAK eksklusif milik `FacRetroDetails`

`[terverifikasi]` Keempat nama juga hidup di struktur lain, dan **tidak boleh disamakan**:

| Nama | Pemilik lain | Bukti |
| --- | --- | --- |
| `BackUpStatus` | `pyWorkPage.Policy.FacOfferList(n)` · `FacOfferList` · halaman `ViewOfferStatus` | `OfferFacOut_PreAct` `<PropertiesName>` `pyWorkPage.Policy.FacOfferList(local.index1).BackUpStatus` ; `Activity\SetIDBackUpStatus` `<PropertiesName>` `.BackUpStatus` |
| `FacNo` | `FacOfferList` | `OfferFacOut_PreAct` `<PropertiesValue>` `FacOfferList.FacNo` |
| `AdditionalInfo` | `.Policy.FacOfferList(1)` · `ViewOfferStatus` | `ShowInformationOfferFacOut` `<pyStepsPreCondParamsWhen>` `@String.equals(ViewOfferStatus.AdditionalInfo,"true")` |
| `DocumentPosition` | `pyWorkPage.Policy.FacOfferList(n)` | `OfferFacOut_PreAct` `<PropertiesName>` `pyWorkPage.Policy.FacOfferList(local.index1).DocumentPosition` |

📌 Pola penulisannya: `OfferFacOut_PreAct` dan `SetPropertyToOfferFacIn_Act` **menyalin dari
`FacOfferList` ke `FacRetroDetails`**. Total kemunculan per nama: `BackUpStatus` **183** ·
`AdditionalInfo` **75** · `DocumentPosition` **42** · `FacNo` **36**; seluruhnya **nol di `DDL\`**.

---

## 2. Anak langsung `FacRetroList` dan letak `FacOutTSI`

### 2.1 ⛔ KOREKSI — anak langsungnya **19**, bukan enam

`[terverifikasi]` Dipindai dari pola `FacRetroList(<indeks>).<X>` di ketiga folder, **19 ruas pertama
unik**:

| Berbentuk koleksi (`*List`) — **5** | Skalar / lainnya — **14** |
| --- | --- |
| `CargoList` · `CurrencyList` · `LocationList` · `PersonList` · `VehicleList` | `EndPeriod` · `OurRef` · `PctPremiAllObj` · `PctPremiAllObjUSD` · `PCTPremiIDR` · `PctShareAllObj` · `PrintRISlip` · `ReinsurerID` · `ReinsurerName` · `RiCommAllObj` · `StartPeriod` · `TFAllObj` · `UjrahAllObj` · `pyExpanded` |

⛔ **Angka "enam koleksi anak" yang beredar sebelumnya SALAH pada dua hal sekaligus:**

1. Ia memasukkan `Property.RiskLocation.AnekaList` dan
   `Property.RiskLocation.OccupationList(n).AnekaList` sebagai "anak `FacRetroList`". Keduanya
   **bukan anak langsung** — letaknya bersarang di bawah `LocationList(n).Property.RiskLocation`.
2. Ia **melewatkan `CurrencyList`**, yang justru **memang** anak langsung.

**Yang benar:** `FacRetroList(n)` punya **5 koleksi anak langsung** dan **14 properti non-koleksi**.

⚠️ `pyExpanded` adalah properti tampilan Pega, bukan data bisnis — dicatat apa adanya agar sensusnya
dapat direproduksi.

### 2.2 `FacOutTSI` melekat di coverage — dikonfirmasi

`[terverifikasi]` **18 jalur properti unik** memuat `FacOutTSI`. Seluruhnya bersandar pada
`.CoverageList(1).FacOutTSI` atau saudara PA-nya `.ASMCoverage(1).FacOutTSI`.

> ⛔ **Jalur berbentuk `FacOutObjectList(..).FacOutTSI` = 0.**

Enam jalur panjang memperlihatkan induknya per lini bisnis:

```
FacRetroList(i).LocationList(j).Property.PropertyItemList(k).CoverageList(1).FacOutTSI
FacRetroList(i).LocationList(j).Property.RiskLocation.AnekaList(m).CoverageList(1).FacOutTSI
FacRetroList(i).LocationList(j).Property.RiskLocation.OccupationList(o).AnekaList(m).CoverageList(1).FacOutTSI
FacRetroList(i).CargoList(c).CoverageList(1).FacOutTSI
FacRetroList(i).VehicleList(v).CoverageList(1).FacOutTSI
FacRetroList(i).PersonList(p).CoverageList(1).FacOutTSI   dan   .PersonList(p).ASMCoverage(1).FacOutTSI
```

📌 `PersonList` memakai **dua** wadah coverage — `CoverageList(1)` **dan** `ASMCoverage(1)`.

---

## 3. Populasi kerja Fac Out

### 3.1 Berbasis nama

`[terverifikasi]` Berkas `Activity\` yang namanya memuat `FacOut`/`FacRetro`/`Facout`/`Facretro`
(Ordinal, **peka huruf**):

| Folder | Bernama Fac Out | Pembawa `10015` **tanpa** nama Fac Out | Permukaan fungsional |
| --- | ---: | ---: | ---: |
| NB FacIn | **41** | **16** | ≈57 |
| RNW Fac In | **40** | **15** | ≈55 |
| Endorsment Fac In | **39** | **16** | ≈55 |
| **Hadir di ketiga folder** | **39** | — | — |

⛔ Populasi berbasis nama adalah **batas kemudahan, bukan batas fungsi** (`PANDUAN-KERJA` §3).

### 3.2 Banding 23 tag atas ke-39 yang hadir di ketiganya

`[terverifikasi]`

| Pasangan | Identik | Berbeda |
| --- | ---: | ---: |
| **NB vs RNW** | **39 / 39** | **0** |
| **NB vs EDM** | **30 / 39** | **9** |

Kesembilan yang berbeda:

`CopyAllObjFacOutFireAneka_ACT` · `CopyFacRetroAnekaGolf_ACT` · `InsertFacoutProd` ·
`InsertFacoutProduction` · `InsertFacoutProductionEDM` · `InsertFacoutProductionEDMCurr` ·
`OfferFacOut_PreAct` · `SetDataFacOutAnekaGolf_Act` · `SetDataFacOutFire_Act`

📌 **NB dan RNW berbagi satu implementasi** — nol perbedaan pada seluruh 39.

⛔ Angka lama **18/18 dan 10/18 tetap DIBATALKAN** — populasinya tidak pernah terdefinisi.

---

## 4. Perintah audit

```powershell
# populasi + banding 23 tag (fungsi Get-Sig sama dengan 07-hitung-berkas-dan-drift.md §4)
$set=@{}
foreach($f in @('NB FacIn','RNW Fac In','Endorsment Fac In')){
  $h=New-Object 'System.Collections.Generic.HashSet[string]' ([StringComparer]::Ordinal)
  foreach($x in [IO.Directory]::EnumerateFiles("D:\migrasi\RNM\$f\Activity",'*.xml')){
    $bn=[IO.Path]::GetFileNameWithoutExtension($x)
    if($bn.IndexOf('FacOut',[StringComparison]::Ordinal) -ge 0 -or
       $bn.IndexOf('FacRetro',[StringComparison]::Ordinal) -ge 0 -or
       $bn.IndexOf('Facout',[StringComparison]::Ordinal) -ge 0 -or
       $bn.IndexOf('Facretro',[StringComparison]::Ordinal) -ge 0){ [void]$h.Add($bn) } }
  $set[$f]=$h; "$f = $($h.Count)" }
# -> 41 / 40 / 39
```

```powershell
# anak langsung FacRetroList -> 19
$anak=New-Object 'System.Collections.Generic.HashSet[string]' ([StringComparer]::Ordinal)
foreach($f in @('NB FacIn','RNW Fac In','Endorsment Fac In')){
  foreach($s in @('Activity','DataTransform','Section','Harness','FlowAction','Flow','RDBList',
                  'ReportDefinition','DataPage','When','DecisionTable','DecisionTree','ConnectREST','SystemSettings')){
    $p="D:\migrasi\RNM\$f\$s"; if(-not [IO.Directory]::Exists($p)){continue}
    foreach($x in [IO.Directory]::EnumerateFiles($p,'*.xml')){
      $t=[IO.File]::ReadAllText($x)
      if($t.IndexOf('FacRetroList',[StringComparison]::Ordinal) -lt 0){ continue }
      foreach($m in [regex]::Matches($t,'FacRetroList\([^)]*\)\.([A-Za-z0-9_.]+)')){
        [void]$anak.Add((($m.Groups[1].Value) -split '\.')[0]) } } } }
$anak.Count      # -> 19
```

```powershell
# jalur FacOutObjectList(..).FacOutTSI -> 0
$n=0
foreach($f in @('NB FacIn','RNW Fac In','Endorsment Fac In')){
  foreach($x in [IO.Directory]::EnumerateFiles("D:\migrasi\RNM\$f\Activity",'*.xml')){
    $t=[IO.File]::ReadAllText($x)
    $n += ([regex]::Matches($t,'FacOutObjectList\([^)]*\)\.FacOutTSI')).Count } }
$n               # -> 0
```

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
