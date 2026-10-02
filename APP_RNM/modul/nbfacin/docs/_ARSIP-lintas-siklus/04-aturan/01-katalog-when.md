# 04-aturan / 01 — Katalog rule `When` Facultative Inward

**Sumber tunggal:** `D:\migrasi\RNM\NB FacIn\When\`, `D:\migrasi\RNM\RNW Fac In\When\`,
`D:\migrasi\RNM\Endorsment Fac In\When\` — ekspor rule Pega, dibaca read-only.
Tidak ada bahan dari luar ketiga folder itu. Discovery lama (`RNM_BRD/`) **tidak dipakai**.

**Tanggal audit:** 2026-09-15 · seluruh angka di bawah dihasilkan ulang oleh perintah PowerShell
yang tercantum di sebelahnya; tidak ada angka yang disalin dari dokumen lain.

---

## 0. Cara dokumen ini membaca korpus

Satu berkas `When/*.xml` memuat kondisi di **dua tempat yang tidak selalu sinkron**:

| Tag | Letak | Sifat |
| --- | --- | --- |
| `<pyConditionString>` | di dalam `<pyConditionViewer>/<pyNestedConditions>/<rowdata>` | representasi editor; **sering kosong atau berisi placeholder** |
| `<pyConditionValue1String>` | di dalam `<pyCondition>/<rowdata>` | label kondisi per baris; **selalu terisi di korpus ini** |
| `<pyConditionValue1>` | di dalam `<pyCondition>/<rowdata>` | ekspresi ter-generate yang benar-benar dieksekusi |
| `<pyLogic>` | akar `<pagedata>` | perangkaian antar-baris (`A OR B`, `A AND (B OR C)`, …) |

Bukti bahwa kedua tempat itu dapat berbeda — `IsPKSASM` **[terverifikasi]**:

```
D:\migrasi\RNM\NB FacIn\When\IsPKSASM.xml
  L161: <pyConditionString>[Double click to add condition]</pyConditionString>
  L336: <pyConditionValue1String>pyWorkPage.OfferFacIn.IsB2B = "ASM" </pyConditionValue1String>
```

```powershell
Select-String -Path "D:\migrasi\RNM\NB FacIn\When\IsPKSASM.xml" `
  -Pattern '<pyConditionString>|<pyConditionValue1String>' |
  ForEach-Object { "L{0}: {1}" -f $_.LineNumber, $_.Line.Trim() }
```

**Konsekuensi metodologis:** dokumen ini membaca `pyConditionValue1String` + `pyConditionValue1`
sebagai sumber utama, dan `pyConditionString` hanya sebagai konfirmasi. Discovery yang hanya
membaca `pyConditionString` akan kehilangan **170 dari 601 berkas** (lihat §6).

### Korpus dirakit dari lebih dari satu server Pega — [terverifikasi]

```powershell
Select-String -Path "D:\migrasi\RNM\*\When\*.xml" -Pattern '<pxCreateSystemID>([^<]*)</pxCreateSystemID>' |
  ForEach-Object { $_.Matches[0].Groups[1].Value } | Group-Object | Sort-Object Count -Descending |
  ForEach-Object { "{0,5}  {1}" -f $_.Count, $_.Name }
```

| Jumlah baris | `pxCreateSystemID` |
| --- | --- |
| 1550 | `pega` |
| 354 | `pegaprdnusare` |
| 275 | `pegadevnusare2` |
| 3 | `sde` |
| 1 | `sls-envcam82` |

Sekurang-kurangnya satu host **DEV** (`pegadevnusare2`) menyumbang isi korpus. Karena itu setiap
klaim di dokumen ini menyebut **path berkas**, bukan sekadar nama rule.

---

## 1. Jumlah rule `When`

```powershell
foreach ($f in @("NB FacIn","RNW Fac In","Endorsment Fac In")) {
  $c = (Get-ChildItem -Path "D:\migrasi\RNM\$f\When" -Filter *.xml -File -Recurse | Measure-Object).Count
  "$f : $c"
}
```

| Folder siklus | Berkas `When\*.xml` |
| --- | --- |
| `NB FacIn` | **210** |
| `RNW Fac In` | **189** |
| `Endorsment Fac In` | **202** |
| **Total berkas** | **601** |

Seluruh 601 berkas berekstensi `.xml`, tidak ada subfolder di bawah `When\` **[terverifikasi]**:

```powershell
Get-ChildItem -Path "D:\migrasi\RNM\NB FacIn\When" -File -Recurse | Group-Object Extension |
  Select-Object Name,Count
```

### Identitas unik (nama rule)

```powershell
$all = foreach ($f in @("NB FacIn","RNW Fac In","Endorsment Fac In")) {
  Get-ChildItem -Path "D:\migrasi\RNM\$f\When" -Filter *.xml -File |
    ForEach-Object { [PSCustomObject]@{Folder=$f; Name=$_.BaseName} } }
"Unik case-insensitive: " + ($all.Name | Sort-Object -Unique).Count
"Unik case-sensitive  : " + ($all.Name | Sort-Object -Unique -CaseSensitive).Count
```

| Ukuran | Nilai |
| --- | --- |
| Nama unik, **case-insensitive** | **226** |
| Nama unik, **case-sensitive** | **237** |

Selisih 11 berasal dari **11 nama rule yang dieja berbeda kapitalisasinya antar siklus**
**[terverifikasi]**:

```powershell
$d = Import-Csv <hasil-ekstraksi>   # lihat §A-lampiran perintah
$d | Group-Object FileName |
  Where-Object { (@($_.Group.FileName | Sort-Object -Unique -CaseSensitive)).Count -gt 1 } |
  Sort-Object Name | ForEach-Object {
    $pairs = $_.Group | Sort-Object Folder | ForEach-Object { "{0}={1}" -f $_.Folder.Substring(0,3), $_.FileName }
    "  {0,-30} {1}" -f $_.Name, ($pairs -join '  ') }
```

| Rule | EDM | NB | RNW |
| --- | --- | --- | --- |
| isAllRisk | `IsAllRisk` | `isAllRisk` | `isAllRisk` |
| isAviationHull | `IsAviationHull` | `isAviationHull` | `isAviationHull` |
| isBillboardNeonSyariah | `IsBillboardNeonSyariah` | `isBillboardNeonSyariah` | `isBillboardNeonSyariah` |
| isBurglary | `IsBurglary` | `isBurglary` | `isBurglary` |
| IsCAR | `IsCar` | `IsCAR` | `IsCAR` |
| isElectronicEquipment | `IsElectronicEquipment` | `isElectronicEquipment` | `isElectronicEquipment` |
| IsExclusion | `isExclusion` | `IsExclusion` | `IsExclusion` |
| isFidelity | `IsFidelity` | `isFidelity` | `isFidelity` |
| IsGroupUwFac | `IsGroupUWFac` | `IsGroupUwFac` | `IsGroupUwFac` |
| IsMaintenance | `isMaintenance` | `IsMaintenance` | `IsMaintenance` |
| isTravelTime | `IsTravelTime` | `isTravelTime` | `isTravelTime` |

Pola konsisten: **selalu EDM yang berbeda dari NB+RNW**, 11 dari 11. Nilai `<pyRuleName>` di dalam
berkas cocok dengan nama berkasnya di semua kasus **[terverifikasi]**.

> **[pertanyaan terbuka]** Apakah resolusi rule Pega di instalasi ini case-sensitive? Bila ya,
> `IsCar` (EDM) dan `IsCAR` (NB/RNW) adalah **dua rule berbeda**, bukan satu rule dengan dua ejaan,
> dan kesebelasnya harus diimplementasikan sebagai predikat terpisah. Korpus tidak menjawab ini.

---

## 2. Tumpang tindih antar siklus — klaim "tiga siklus di atas satu basis rule"

```powershell
$g = $all | Group-Object Name
$g | Group-Object { ($_.Group.Folder | Sort-Object -Unique).Count } |
  Select-Object @{n='JumlahFolder';e={$_.Name}}, Count | Sort-Object JumlahFolder
```

| Muncul di | Jumlah nama rule | % dari 226 |
| --- | --- | --- |
| **3 folder** (NB + RNW + EDM) | **183** | 81,0 % |
| **2 folder** | **9** | 4,0 % |
| **1 folder saja** | **34** | 15,0 % |

**Verdict: klaim "tiga siklus di atas satu basis rule" TERKONFIRMASI pada tataran identitas rule
(81 % nama dipakai ketiga siklus), tetapi TIDAK pada tataran isi** — 39 dari 192 nama yang dipakai
di ≥2 folder ternyata isinya berbeda antar siklus (§3). Basisnya bersama; sebagian perilakunya
tidak.

### Rule eksklusif satu siklus (34)

```powershell
$g | Where-Object { ($_.Group.Folder | Sort-Object -Unique).Count -eq 1 } |
  Group-Object { $_.Group[0].Folder } | ForEach-Object { "{0,3}  {1}" -f $_.Count, $_.Name }
```

| Folder | Jumlah rule eksklusif |
| --- | --- |
| `NB FacIn` | **18** |
| `Endorsment Fac In` | **16** |
| `RNW Fac In` | **0** |

**RNW tidak memiliki satu pun rule `When` yang hanya miliknya** — himpunan rule RNW adalah
himpunan bagian murni dari gabungan NB+EDM **[terverifikasi]**. Ini temuan struktural: siklus
renewal tidak menambah predikat baru apa pun, ia memakai ulang predikat NB.

**Eksklusif `NB FacIn` (18):** `crmCreateOpportunity`, `isApproved`, `isClaimTreaty`, `IsNotAdmin`,
`IsOfferFacIn`, `IsOperatorLife`, `isSellingModeB2B`, `isSellingModeB2BB2C`, `isSellingModeB2C`,
`IsSPVCreate`, `IsSPVTreaty1`, `IsTreaty1`, `NopolisEmpty`, `pyIsIpadOrDesktop`, `pyIsMobile`,
`ToDeptHeadUWLife`, `ToTREATYDEPTHEAD`, `TreatyMasterInEDM`.

**Eksklusif `Endorsment Fac In` (16):** `IsCustomBond`, `IsEdmAdjCeding`, `IsEdmInternalRetro`,
`IsHealthSSC`, `IsMaterialDamage`, `IsMultipleObyekMBU`, `isNew`, `IsNotActive`, `IsNotFire`,
`IsRequired`, `IsShortPeriod`, `IsSpreadingDepan`, `IsTahun1`, `IsTJHInRange`, `IsTJHNotInRange`,
`recordEvent`.

> Catatan: `TreatyMasterInEDM` ada **hanya di folder NB**, meski namanya menyebut EDM. Nama rule
> bukan bukti tempat pakainya **[terverifikasi: berkas hanya ada di `NB FacIn\When\`]**.

### Rule di tepat 2 folder (9)

| Rule | Folder |
| --- | --- |
| `hasPrimaryPage` | NB, RNW |
| `IsDeclarationPolicy` | NB, RNW |
| `IsFlagUW` | NB, RNW |
| `IsIT` | NB, RNW |
| `IsTABActive_LossRecord` | NB, RNW |
| `LetterNoNull` | NB, RNW |
| `IsProposalTransfer` | NB, EDM |
| `IsUWAccepted` | NB, EDM |
| `StepStatusFail` | NB, EDM |

---

## 3. Perbandingan hash: identik vs bercabang

### 3.1 Hash berkas mentah (`Get-FileHash`) hampir tidak berguna

```powershell
# Hash dihitung di skrip ekstraksi: (Get-FileHash -LiteralPath $file.FullName -Algorithm SHA256).Hash
$multi = $d | Group-Object FileName | Where-Object { ($_.Group.Folder | Sort-Object -Unique).Count -ge 2 }
"Nama rule di >=2 folder : " + $multi.Count
"  identik byte-for-byte : " + ($multi | Where-Object { ($_.Group.Hash | Sort-Object -Unique).Count -eq 1 }).Count
"  hash BERBEDA          : " + ($multi | Where-Object { ($_.Group.Hash | Sort-Object -Unique).Count -gt 1 }).Count
```

| Ukuran | Nilai |
| --- | --- |
| Nama rule hadir di ≥2 folder | **192** |
| SHA-256 berkas **identik** di semua folder | **6** |
| SHA-256 berkas **berbeda** | **186** |

Enam rule yang identik byte-for-byte **[terverifikasi]**: `hasPrimaryPage`, `IsDeclarationPolicy`,
`IsFlagUW`, `IsIT`, `IsTABActive_LossRecord`, `LetterNoNull` — semuanya pasangan **NB ↔ RNW**.

186 sisanya berbeda hash karena ekspor Pega membawa metadata per-instance (`pxUpdateDateTime`,
`pzInsKey`, `pxCommitDateTime`, `pxHostId`, `pzChecksum`, `pyJavaClassName`, metadata operator).
Perbedaan hash **bukan** bukti perbedaan perilaku.

### 3.2 Hash logika ternormalisasi — inilah yang menjawab pertanyaan

Dua normalisasi dihitung terpisah:

```powershell
# H1 = ClassName + pyLogic + seluruh ekspresi pyConditionValue1 (spasi diciutkan, lower-case)
# H2 =             pyLogic + seluruh ekspresi pyConditionValue1
$sha = [System.Security.Cryptography.SHA256]::Create()
foreach ($r in $d) {
  $h1 = (($r.ClassName + '#' + $r.Logic + '#' + $r.V1Exprs) -replace '\s+',' ').Trim().ToLowerInvariant()
  $h2 = ((                 $r.Logic + '#' + $r.V1Exprs) -replace '\s+',' ').Trim().ToLowerInvariant()
  # ... ComputeHash(UTF8.GetBytes($hN))
}
```

| Ukuran | Nilai |
| --- | --- |
| Nama rule hadir di ≥2 folder | **192** |
| H2 identik (pyLogic + ekspresi sama) | **153** |
| H2 **bercabang** | **39** |
| dari 153 yang identik: `pyClassName` tetap berbeda | **10** |
| H1 identik (ikut memperhitungkan `pyClassName`) | **143** |
| H1 bercabang | **49** |

**39 rule benar-benar bercabang antar siklus.** Seluruhnya melibatkan `Endorsment Fac In` sebagai
pihak yang berbeda; NB dan RNW selalu sepakat satu sama lain **[terverifikasi]**.

### 3.3 Tiga pola percabangan — temuan bernilai tinggi

**Pola P1 — referensi properti berbeda (30 rule).** Bentuk paling lazim: EDM membaca agregat
tersimpan `.OfferFacIn.QuotationData.*` sedangkan NB/RNW membaca halaman quotation aktif
`pyWorkPage.Quotation.*`. Contoh `isBurglary` **[terverifikasi]**:

| Siklus | Berkas | `pyConditionValue1String` |
| --- | --- | --- |
| EDM | `Endorsment Fac In\When\IsBurglary.xml` | `.OfferFacIn.QuotationData.BusinessType = "Burglary"` |
| NB | `NB FacIn\When\isBurglary.xml` | `pyWorkPage.Quotation.BusinessType = "Burglary"` |
| RNW | `RNW Fac In\When\isBurglary.xml` | `pyWorkPage.Quotation.BusinessType = "Burglary"` |

Pola yang sama (29 rule lain): `isAllRisk`, `isAviationHull`, `isBillboardNeonSyariah`, `IsBoiler`,
`IsBondingKBG`, `IsCIS`, `IsCIT`, `IsContractorsPlantMachinery`, `IsCustomBonds`, `IsEar`,
`isElectronicEquipment`, `IsExclusion`, `isFidelity`, `IsFireStyle1`, `IsFireStyle2`, `IsGlass`,
`IsGrowingTrees`, `IsHE`, `IsKPR`, `IsLandRig`, `IsLiability`, `IsMaintenance`, `IsMBD`,
`IsMBUCar`, `IsOilGas`, `IsFacRetro`, `IsCorporate`, `IsPASSG`, `IsProposalTransfer`.

Empat di antaranya menyimpang **ke arah berlawanan** — EDM memakai path absolut sementara NB/RNW
memakai path relatif **[terverifikasi]**, mis. `IsCorporate`: EDM `pyWorkPage.Policy.CustomerType = "2"`
(kelas `ASM-FW-GCNMFW-Work`) vs NB/RNW `.CustomerType = "2"` (kelas `ASM-FW-GISFW-Data-Policy`);
demikian pula `IsFacRetro`, `IsPASSG`, `IsProposalTransfer`. Arahnya berbeda, kelas kerentanannya
sama: **dua ekspresi yang tidak dijamin membaca nilai yang sama**.

> **[pertanyaan terbuka]** `.Quotation.BusinessType` dan `.OfferFacIn.QuotationData.BusinessType`
> adalah dua properti berbeda. Apakah keduanya selalu bernilai sama saat rule dievaluasi? Bila
> tidak, klasifikasi lini bisnis di endorsement **dapat berbeda** dari klasifikasi di NB untuk polis
> yang sama — dan itu mengubah jalur akseptasi.

**Pola P2 — EDM membaca hasil query eksternal `OutData.pxResults(1).CARI2`, bukan properti kasus
(6 rule) [terverifikasi]:**

| Rule | EDM (`ASM-SFAGIS-Work-Endorsement`) | NB/RNW |
| --- | --- | --- |
| `IsAneka` | 25 baris `OutData.pxResults(1).CARI2 = …` | 26 baris, mayoritas `Rule <X> evaluates to true` |
| `IsFire` | 5 baris `OutData.pxResults(1).CARI2 = …` | 4 baris `Rule …` + 1 baris `.OfferFacIn.QuotationData.BusinessType = "Fire"` |
| `IsGolfInsurance` | `OutData.pxResults(1).CARI2 = "GolfInsurance"` | `pyWorkPage.Quotation.BusinessType = "GolfInsurance"` |
| `IsMarineCargo` | `OutData.pxResults(1).CARI2 = "MarineCargo"` | `.OfferFacIn.QuotationData.BusinessType` **OR** `pyWorkPage.Quotation.BusinessType` |
| `IsMBU` | 4 baris CARI2 (2 di antaranya **duplikat**) | 2 baris `pyWorkPage.OfferFacIn.QuotationData.BusinessType` |
| `IsPA` | `OutData.pxResults(1).CARI2 = "PA"` | `pyWorkPage.Quotation.BusinessType = "PA"` |

> **[pertanyaan terbuka]** Apa itu `OutData.pxResults(1).CARI2`? Nama halaman `OutData` dan kolom
> `CARI2` menunjukkan hasil sebuah RDB/Report, tetapi **isi query-nya tidak ada di folder `When\`**.
> Selama isinya belum diketahui, kelima rule ini tidak dapat diimplementasikan untuk siklus EDM.
> `pxResults(1)` juga berarti **hanya baris pertama** yang dibaca — bila query mengembalikan lebih
> dari satu baris, perilakunya bergantung pada urutan hasil yang tidak ditentukan di rule ini.

**Pola P3 — himpunan nilai yang diuji berbeda (percabangan semantik sesungguhnya; 3 rule).**
`IsUW` **[terverifikasi]**:

| Siklus | `pyLogic` | Workbasket yang diuji |
| --- | --- | --- |
| EDM | `A OR B OR C OR D` | `ReasFacInDirector`, **`ReasFacInMarketing`**, `ReasFacInGroupLeader`, `ReasFacInUnderwriting` |
| NB | `A OR B OR C` | `ReasFacInDirector`, `ReasFacInGroupLeader`, `ReasFacInUnderwriting` |
| RNW | `A OR B OR C` | `ReasFacInDirector`, `ReasFacInGroupLeader`, `ReasFacInUnderwriting` |

Berkas: `Endorsment Fac In\When\IsUW.xml` vs `NB FacIn\When\IsUW.xml` vs `RNW Fac In\When\IsUW.xml`.
Properti yang diuji: `pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName`.

> **Ini bukan perbedaan kosmetik.** Operator di workbasket `ReasFacInMarketing` dianggap "UW" saat
> endorsement tetapi **tidak** saat new business / renewal. Migrasi yang menyeragamkan keduanya akan
> mengubah siapa yang boleh bertindak. **[pertanyaan terbuka]** — keputusan bisnis, bukan keputusan
> migrasi.

Dua percabangan semantik lain **[terverifikasi]**:

| Rule | EDM | NB / RNW |
| --- | --- | --- |
| `IsClaim` | `pyWorkPage.ClaimData has a value` (`pyConditionValue1Purpose` = `PropertyHasValue`, kelas `Work-`) | `pyWorkPage.pyWorkIDPrefix = "CLM-"` (kelas `ASM-FW-GISFW-Data`) |
| `IsTravel` | `A OR B`: `.Quotation.BusinessType = "Travel"` **OR** `pyWorkPage.OfferFacIn.QuotationData.BusinessType = "Travel"` | `A`: `pyWorkPage.Quotation.BusinessType = "Travel"` saja |

**Daftar lengkap 39 rule bercabang** (kolom "Kondisi terbaca" di §4 menandainya dengan
`**BERCABANG**` dan menampilkan tiap varian):

`isAllRisk`, `IsAneka`, `isAviationHull`, `isBillboardNeonSyariah`, `IsBoiler`, `IsBondingKBG`,
`isBurglary`, `IsCIS`, `IsCIT`, `IsClaim`, `IsContractorsPlantMachinery`, `IsCorporate`,
`IsCustomBonds`, `IsEar`, `isElectronicEquipment`, `IsExclusion`, `IsFacRetro`, `isFidelity`,
`IsFire`, `IsFireStyle1`, `IsFireStyle2`, `IsGlass`, `IsGolfInsurance`, `IsGrowingTrees`, `IsHE`,
`IsKPR`, `IsLandRig`, `IsLiability`, `IsMaintenance`, `IsMarineCargo`, `IsMBD`, `IsMBU`, `IsMBUCar`,
`IsOilGas`, `IsPA`, `IsPASSG`, `IsProposalTransfer`, `IsTravel`, `IsUW`.

### 3.4 Sepuluh rule yang ekspresinya sama tetapi `pyClassName`-nya berbeda

`FlagOldData`, `IsBondingAndCustomBonds`, `IsCrime`, `IsEngineering`, `IsEnvironmental`,
`IsFlagDelete`, `IsFlagOldData`, `IsGroupFac`, `IsLife`, `IsNotEDM` **[terverifikasi]**.

Contoh `IsNotEDM`: `pyClassName` = `ASM-FW-GISFW-Data-Vehicle` di NB, sementara varian lain
memakai kelas berbeda, dengan ekspresi identik `pyWorkPage.Quotation.StatusBusiness != 3`.

> **[dugaan]** `pyClassName` pada rule `When` Pega menentukan konteks resolusi rule (kelas mana yang
> "memiliki" rule itu). Bila benar, dua rule dengan nama sama tetapi kelas berbeda **tidak saling
> menimpa** dan bisa aktif bersamaan pada kasus yang berbeda. Korpus `When\` sendiri tidak
> membuktikan mekanisme resolusinya → **[pertanyaan terbuka]**.

---

## 4. Katalog lengkap — 226 rule unik

Satu baris per nama rule unik (case-insensitive). Kolom **Siklus** memakai `NB` / `RNW` / `EDM`.
Kolom kondisi mengutip isi `<pyConditionValue1String>` **apa adanya**, didahului `pyLogic` dalam
kurung siku. Bila sebuah rule punya >2 baris kondisi yang menguji **properti dan operator yang
sama**, baris-baris itu diciutkan menjadi `properti op` + himpunan nilai — himpunan nilainya lengkap,
tidak dipotong.

Kolom **Sumber kondisi**:

- `keduanya` — `<pyConditionString>` terisi teks nyata **dan** `<pyConditionValue1String>` terisi
- `V1-saja (CS placeholder)` — `<pyConditionString>` berisi `[Double click to add condition]`
- `V1-saja (CS kosong)` — tag `<pyConditionString/>` self-closing

Nilai literal pada rule bergantung identitas orang **disensor** sesuai CLAUDE.md §3.5; mekanismenya
dicatat di §7.

Seluruh **226** baris hadir; tema tercantum di §4.1 sebagai indeks silang, bukan sebagai pengganti
tabel.


| Rule | Siklus | pyLogic + kondisi terbaca (isi tag apa adanya) | Sumber kondisi |
| --- | --- | --- | --- |
| `crmCreateOpportunity` | NB | **[(A0 Or A1) And A2]** `true = @(Pega-RULES:ExpressionEvaluators).evaluateWhen("crmBypassOperatorAccessChecks")` ; `"true" equals Declare_crmOperatorAccess.canCreateOpportunity` ; `true = @(Pega-RULES:ExpressionEvaluators).evaluateWhen("crmIsOpen")` | V1-saja (CS placeholder) |
| `FlagOldData` | NB RNW EDM | **[A]** `.FlagOldData = 1` | keduanya |
| `hasPrimaryPage` | NB RNW | **[A0]** `@hasPrimaryPage(tools)` | keduanya |
| `IsAddButton` | NB RNW EDM | **[A]** `pyPortal.IsAddButton = "True"` | keduanya |
| `IsAddCoverage` | NB RNW EDM | **[A OR B OR C OR D OR E OR F OR G OR H]** `.Coverage =` salah satu dari { "100143", "100847", "100857", "100858", "100861", "100867", "100869", "100872" } | keduanya |
| `IsAddMedEx` | NB RNW EDM | **[A OR B]** `.Coverage = "100855"` ; `.Coverage = "100856"` | keduanya |
| `IsAddMedExPA` | NB RNW EDM | **[A OR B]** `Rule IsAddMedEx evaluates to true` ; `Rule IsAddPA evaluates to true` | keduanya |
| `IsAddPA` | NB RNW EDM | **[A OR B]** `.Coverage = "100848"` ; `.Coverage = "100849"` | keduanya |
| `IsAddTJHPLL` | NB RNW EDM | **[A OR B]** `Rule IsTJH evaluates to true` ; `Rule IsPLL evaluates to true` | keduanya |
| `IsAdjustableFlag` | NB RNW EDM | **[A]** `.IsAdjustableFlag = true` | keduanya |
| `IsAdmin` | NB RNW EDM | **[A OR B]** `OperatorID.pyWorkBasketList(1).pyWorkBasketName = "ReasFacInAdmin"` ; `OperatorID.pyWorkBasketList(1).pyWorkBasketName = ""` | keduanya |
| `isAllRisk / IsAllRisk` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "AllRisk"`<br>**EDM** **[A OR B]** `.Quotation.BusinessType = "AllRisk"` ; `pyWorkPage.Quotation.BusinessType = "AllRisk"` | keduanya |
| `IsAneka` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A OR B OR C OR E OR F OR G OR H OR I OR J OR K OR L OR M OR N OR O OR P OR Q OR R OR S OR T OR U OR V OR W OR X OR Y OR Z OR AA]** `Rule IsLiability evaluates to true` ; `Rule IsMarineHull evaluates to true` ; `Rule isGrowingTrees evaluates to true` ; `Rule IsExclusion evaluates to true` ; `Rule isCAR evaluates to true` ; `Rule IsMaintenance evaluates to true` ; `Rule isElectronicEquipment evaluates to true` ; `Rule isAviationHull evaluates to true` ; `Rule isEAR evaluates to true` ; `Rule isAllRisk evaluates to true` ; `Rule IsGlass evaluates to true` ; `Rule isFidelity evaluates to true` ; `Rule isBillboardNeonSyariah evaluates to true` ; `Rule isBurglary evaluates to true` ; `Rule IsCIT evaluates to true` ; `Rule IsCIS evaluates to true` ; `Rule IsMBD evaluates to true` ; `Rule IsBoiler evaluates to true` ; `Rule IsHE evaluates to true` ; `Rule IsLandRig evaluates to true` ; `Rule IsContractorsPlantMachinery evaluates to true` ; `Rule IsYieldShortfall evaluates to true` ; `Rule IsCrime evaluates to true` ; `Rule IsEnvironmental evaluates to true` ; `pyWorkPage.Quotation.BusinessType = "Aneka"` ; `Rule IsBondingAndCustomBonds evaluates to true`<br>**EDM** **[A OR B OR C  OR E OR F OR G OR H OR I OR J OR K OR L OR M OR N OR O OR P OR Q OR R OR S OR T OR U OR V OR W OR X OR Y OR Z]** `OutData.pxResults(1).CARI2 =` salah satu dari { "AllRisk", "Aneka", "AviationHull", "BillboardNeonSyariah", "Boiler", "Bonding", "BondingKBG", "Burglary", "Car", "CIS", "CIT", "ContractorsPM", "CustomBond", "Ear", "ElectronicEquipment ", "Exclusion", "Fidelity", "Glass", "GrowingTrees", "HE", "LandRig", "Liability", "Maintenance", "MarineHull", "MBD" } | keduanya / V1-saja (CS placeholder) |
| `IsAnekaNotBonding` | NB RNW EDM | **[A OR B OR C OR E OR F OR G OR H OR I OR J OR K OR L OR M OR N OR O OR P OR Q OR R OR S OR T OR U OR V OR W]** `Rule IsLiability evaluates to true` ; `Rule IsMarineHull evaluates to true` ; `Rule IsGrowingTrees evaluates to true` ; `Rule isExclusion evaluates to true` ; `Rule isMaintenance evaluates to true` ; `Rule IsCar evaluates to true` ; `Rule IsElectronicEquipment evaluates to true` ; `Rule IsAviationHull evaluates to true` ; `Rule IsEar evaluates to true` ; `Rule IsAllRisk evaluates to true` ; `Rule IsGlass evaluates to true` ; `Rule isFidelity evaluates to true` ; `Rule IsBillboardNeonSyariah evaluates to true` ; `Rule IsBurglary evaluates to true` ; `Rule IsCIT evaluates to true` ; `Rule IsCIS evaluates to true` ; `Rule IsBoiler evaluates to true` ; `Rule IsMBD evaluates to true` ; `Rule IsHE evaluates to true` ; `Rule IsLandRig evaluates to true` ; `Rule IsContractorsPlantMachinery evaluates to true` ; `.Quotation.BusinessType = "Aneka"` | keduanya |
| `isApproved` | NB | **[A]** `pyWorkPage.PolicyTreatyIn.IsApproved = 1` | keduanya |
| `IsAssociationPolicy` | NB RNW EDM | **[A]** `.Policy.TypeOfPolicy = "02"` | keduanya |
| `IsAsuransiKredit` | NB RNW EDM | **[A OR B OR C OR D]** `pyWorkPage.Quotation.BusinessCode =` salah satu dari { "10053", "10243", "10244", "10245" } | keduanya |
| `isAviationHull / IsAviationHull` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "AviationHull"`<br>**EDM** **[A]** `.OfferFacIn.QuotationData.BusinessType = "AviationHull"` | keduanya |
| `IsBillboardNeon` | NB RNW EDM | **[A]** `pyWorkPage.Quotation.BusinessName = "BILLBOARD/NEON SIGN"` | keduanya |
| `isBillboardNeonSyariah` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A OR B]** `pyWorkPage.Quotation.BusinessType = "BillboardNeonSyariah"` ; `pyWorkPage.Quotation.BusinessCode = "10106"`<br>**EDM** **[A OR B]** `.OfferFacIn.QuotationData.BusinessType = "BillboardNeonSyariah"` ; `.OfferFacIn.QuotationData.BusinessCode = "10106"` | keduanya |
| `isBisnisOtherObjAneka` | NB RNW EDM | **[A OR B OR D OR E]** `pyWorkPage.Quotation.BusinessCode =` salah satu dari { "15", "21", "32", "SQ" } | V1-saja (CS placeholder) |
| `IsBoiler` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "Boiler"`<br>**EDM** **[A]** `.OfferFacIn.QuotationData.BusinessType = "Boiler"` | keduanya |
| `IsBonding` | NB RNW EDM | **[A OR B OR C]** `.Quotation.BusinessType = "Bonding"` ; `Rule IsBondingKBG evaluates to true` ; `pyWorkPage.Quotation.BusinessType = "Bonding"` | V1-saja (CS placeholder) |
| `IsBondingAndCustomBonds` | NB RNW EDM | **[A OR B OR C]** `pyWorkPage.Quotation.BusinessType = "Bonding"` ; `Rule IsBondingKBG evaluates to true` ; `Rule IsCustomBonds evaluates to true` | keduanya |
| `IsBondingKBG` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A OR B]** `pyWorkPage.Quotation.BusinessType = "BondingKBG"` ; `pyWorkPage.Quotation.BusinessType = "BondingKBG"`<br>**EDM** **[A OR B]** `.Quotation.BusinessType = "BondingKBG"` ; `pyWorkPage.Quotation.BusinessType = "BondingKBG"` | keduanya |
| `IsBuilderRisk` | NB RNW EDM | **[A OR B]** `pyWorkPage.OfferFacIn.QuotationData.BusinessCode = "10085"` ; `pyWorkPage.OfferFacIn.QuotationData.BusinessName = "BUILDER RISK"` | V1-saja (CS placeholder) |
| `isBurglary / IsBurglary` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "Burglary"`<br>**EDM** **[A]** `.OfferFacIn.QuotationData.BusinessType = "Burglary"` | keduanya |
| `IsButtonShow` | NB RNW EDM | **[A]** `pyPortal.IsButtonShow = "True"` | keduanya |
| `IsCAR` | NB RNW EDM | **[A]** `pyWorkPage.Quotation.BusinessType = "Car"` | keduanya |
| `IsCedingConfirmOffer` | NB RNW EDM | **[A OR B]** `.IsCedingConfirm = Offer` ; `.IsCedingConfirm != "accepted"` | keduanya |
| `IsCedingConfirmPolicy` | NB RNW EDM | **[A OR B]** `.IsCedingConfirm = Policy` ; `.IsCedingConfirm = "accepted"` | keduanya |
| `IsChekCity` | NB RNW EDM | **[A AND B]** `InputCity.ID != ""` ; `InputCity.Note != ""` | V1-saja (CS placeholder) |
| `IsChekNation` | NB RNW EDM | **[A AND B]** `InputNation.ID != ""` ; `InputNation.Note != ""` | V1-saja (CS placeholder) |
| `IsChekProvince` | NB RNW EDM | **[A AND B]** `InputProvince.ID != ""` ; `InputProvince.Note != ""` | V1-saja (CS placeholder) |
| `IsChekRW` | NB RNW EDM | **[A AND B]** `InputRW.ID != ""` ; `InputRW.Note != ""` | V1-saja (CS placeholder) |
| `IsCIS` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "CIS"`<br>**EDM** **[A]** `.OfferFacIn.QuotationData.BusinessType = "CIS"` | keduanya |
| `IsCIT` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "CIT"`<br>**EDM** **[A]** `.OfferFacIn.QuotationData.BusinessType = "CIT"` | keduanya |
| `isCitCis` | NB RNW EDM | **[A OR B OR C OR D]** `pyWorkPage.Quotation.BusinessCode =` salah satu dari { "10040", "10041", "10103", "10104" } | keduanya |
| `IsCivilEngineeringCompletedRisks` | NB RNW EDM | **[A AND (B OR C)]** `pyWorkPage.Quotation.BusinessType = "Aneka"` ; `pyWorkPage.Quotation.BusinessCode = "10166"` ; `pyWorkPage.Quotation.BusinessName = "CIVIL ENGINEERING COMPLETED RISKS"` | keduanya |
| `IsClaim` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.pyWorkIDPrefix = "CLM-"`<br>**EDM** **[A]** `pyWorkPage.ClaimData has a value` | keduanya |
| `isClaimTreaty` | NB | **[G AND (A OR B OR C OR D OR E OR F)]** `.PolicyTreatyIn.ClaimType = "XOL"` ; `.PolicyTreatyIn.Claim != 0` ; `.PolicyTreatyIn.SalvageValue != 0` ; `.PolicyTreatyIn.ExcessLoss != 0` ; `.PolicyTreatyIn.ClaimPaymentType = "Salvage"` ; `.PolicyTreatyIn.ClaimPaymentType = "Claim"` ; `.Quotation.ProportionalType != "Proportional"` | V1-saja (CS placeholder) |
| `IsCMI` | NB RNW EDM | **[A]** `.OfferFacIn.QuotationData.BusinessName = "COMPREHENSIVE MACHINERIES INSURANCE"` | keduanya |
| `IsCoas` | NB RNW EDM | **[A AND B]** `.Policy.TypeOfCoins != 0` ; `.Policy.TypeOfCoins != "F"` | keduanya |
| `IsComprehensiveMachineriesInsurance` | NB RNW EDM | **[A AND (B OR C)]** `pyWorkPage.Quotation.BusinessType = "Aneka"` ; `pyWorkPage.Quotation.BusinessCode = "10187"` ; `pyWorkPage.Quotation.BusinessName = "COMPREHENSIVE MACHINERIES INSURANCE"` | keduanya |
| `IsContractorsPlantMachinery` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "ContractorsPM"`<br>**EDM** **[A]** `.OfferFacIn.QuotationData.BusinessType = "ContractorsPM"` | keduanya |
| `IsCorporate` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `.CustomerType = "2"`<br>**EDM** **[A]** `pyWorkPage.Policy.CustomerType = "2"` | keduanya |
| `isCoverageOtherObjAneka` | NB RNW EDM | **[A OR B OR C]** `Rule IsCIS evaluates to true` ; `Rule isGolfInsurance evaluates to true` ; `Rule IsAviationHull evaluates to true` | keduanya |
| `IsCreditBriguna` | NB RNW EDM | **[A]** `pyWorkPage.Quotation.BusinessCode = "10235"` | keduanya |
| `IsCrime` | NB RNW EDM | **[A AND B]** `pyWorkPage.Quotation.BusinessCode = 10168` ; `pyWorkPage.Quotation.BusinessType = "Aneka"` | V1-saja (CS placeholder) |
| `IsCustomBond` | EDM | **[A]** `.OfferFacIn.QuotationData.BusinessType = "CustomBond"` | keduanya |
| `IsCustomBonds` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A OR B OR C]** `pyWorkPage.Quotation.BusinessType = "CustomBond"` ; `pyWorkPage.Policy.Quotation.BusinessType = "CustomBond"` ; `pyWorkPage.OfferFacIn.QuotationData.BusinessType = "CustomBond"`<br>**EDM** **[A OR B OR C]** `.Quotation.BusinessType = "CustomBond"` ; `pyWorkPage.Policy.Quotation.BusinessType = "CustomBond"` ; `.OfferFacIn.QuotationData.BusinessType = "CustomBond"` | keduanya |
| `IsDeclarationPolicy` | NB RNW | **[A]** `pyWorkPage.Quotation.PolicyType = 2` | keduanya |
| `IsDeducType` | NB RNW EDM | **[A]** `@String.notEquals(.DeductibleType,"3")&&@String.notEquals(.DeductibleType,"1")&&@String.notEquals(.DeductibleType,"2")&&@String.notEquals(.DeductibleType,"5")` | keduanya |
| `IsDescDeductType` | NB RNW EDM | **[A OR B OR C OR D]** `.DeductibleType =` salah satu dari { "100413", "100416", "100419", "100420" } | V1-saja (CS placeholder) |
| `IsDM` | NB RNW EDM | **[A OR B]** `.Quotation.BusinessCode = 86` ; `.Quotation.BusinessCode = "10135"` | keduanya |
| `IsEar` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "Ear"`<br>**EDM** **[A]** `.OfferFacIn.QuotationData.BusinessType = "Ear"` | keduanya |
| `IsEDM` | NB RNW EDM | **[A]** `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3` | keduanya |
| `IsEdmAddObject` | NB RNW EDM | **[(A AND B AND C) OR D]** `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3` ; `pyWorkPage.OfferFacIn.QuotationData.EdmType = 4` ; `pyWorkPage.OfferFacIn.QuotationData.Type = 5` ; `Rule IsEdmAdjTSI evaluates to true` | V1-saja (CS placeholder) |
| `IsEdmAdjCeding` | EDM | **[A AND B AND C]** `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3` ; `pyWorkPage.OfferFacIn.QuotationData.EdmType = 4` ; `pyWorkPage.OfferFacIn.QuotationData.Type = 12` | V1-saja (CS placeholder) |
| `IsEdmAdjCurrency` | NB RNW EDM | **[A AND B AND C]** `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3` ; `pyWorkPage.OfferFacIn.QuotationData.EdmType = 4` ; `pyWorkPage.OfferFacIn.QuotationData.Type = 7` | V1-saja (CS placeholder) |
| `IsEdmAdjInsured` | NB RNW EDM | **[A AND B AND C]** `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3` ; `pyWorkPage.OfferFacIn.QuotationData.EdmType = 4` ; `pyWorkPage.OfferFacIn.QuotationData.Type = 9` | V1-saja (CS placeholder) |
| `IsEdmAdjPeriod` | NB RNW EDM | **[A AND B AND C]** `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3` ; `pyWorkPage.OfferFacIn.QuotationData.EdmType = 4` ; `pyWorkPage.OfferFacIn.QuotationData.Type = 6` | V1-saja (CS placeholder) |
| `IsEdmAdjRate` | NB RNW EDM | **[(A AND B AND C) OR D]** `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3` ; `pyWorkPage.OfferFacIn.QuotationData.EdmType = 4` ; `pyWorkPage.OfferFacIn.QuotationData.Type = 3` ; `Rule IsEdmAdjTSI evaluates to true` | V1-saja (CS placeholder) |
| `IsEdmAdjRefNo` | NB RNW EDM | **[A AND B AND C]** `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3` ; `pyWorkPage.OfferFacIn.QuotationData.EdmType = 4` ; `pyWorkPage.OfferFacIn.QuotationData.Type = 0` | V1-saja (CS placeholder) |
| `IsEdmAdjRIC` | NB RNW EDM | **[A AND B AND (C OR D)]** `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3` ; `pyWorkPage.OfferFacIn.QuotationData.EdmType = 4` ; `pyWorkPage.OfferFacIn.QuotationData.Type = 8` ; `pyWorkPage.OfferFacIn.QuotationData.Type = 12` | V1-saja (CS placeholder) |
| `IsEdmAdjShareCedant` | NB RNW EDM | **[A AND B AND C]** `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3` ; `pyWorkPage.OfferFacIn.QuotationData.EdmType = 4` ; `pyWorkPage.OfferFacIn.QuotationData.Type = 11` | V1-saja (CS placeholder) |
| `IsEdmAdjSpreading` | NB RNW EDM | **[A AND B AND C]** `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3` ; `pyWorkPage.OfferFacIn.QuotationData.EdmType = 4` ; `pyWorkPage.OfferFacIn.QuotationData.Type = 4` | V1-saja (CS placeholder) |
| `IsEdmAdjTSI` | NB RNW EDM | **[A AND B AND C]** `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3` ; `pyWorkPage.OfferFacIn.QuotationData.EdmType = 4` ; `pyWorkPage.OfferFacIn.QuotationData.Type = 2` | V1-saja (CS placeholder) |
| `IsEdmExtendPeriod` | NB RNW EDM | **[A AND B AND C]** `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3` ; `pyWorkPage.OfferFacIn.QuotationData.EdmType = 4` ; `pyWorkPage.OfferFacIn.QuotationData.Type = 1` | V1-saja (CS placeholder) |
| `IsEdmInternalRetro` | EDM | **[A]** `pyWorkPage.OfferFacIn.QuotationData.EndorsementInternalRetro = 1` | V1-saja (CS kosong) |
| `IsEdmPerubahan` | NB RNW EDM | **[A OR B OR C OR D OR E OR F OR G OR H OR I OR J OR K OR L]** `Rule IsEdmAddObject evaluates to true` ; `Rule IsEdmAdjCurrency evaluates to true` ; `Rule IsEdmAdjPeriod evaluates to true` ; `Rule IsEdmAdjRate evaluates to true` ; `Rule IsEdmAdjRIC evaluates to true` ; `Rule IsEdmAdjSpreading evaluates to true` ; `Rule IsEdmAdjTSI evaluates to true` ; `Rule IsEdmExtendPeriod evaluates to true` ; `Rule IsEdmAdjInsured evaluates to true` ; `Rule IsEdmAdjRefNo evaluates to true` ; `Rule IsEdmAdjShareCedant evaluates to true` ; `Rule IsEDMRiSlip evaluates to true` | V1-saja (CS placeholder) |
| `IsEdmPPNPPH` | NB RNW EDM | **[A AND B AND C]** `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3` ; `pyWorkPage.OfferFacIn.QuotationData.EdmType = 4` ; `pyWorkPage.OfferFacIn.QuotationData.Type = 12` | V1-saja (CS placeholder) |
| `IsEDMRiSlip` | NB RNW EDM | **[A OR B]** `pyWorkPage.OfferFacIn.QuotationData.EdmTypeNew = 4` ; `pyWorkPage.OfferFacIn.QuotationData.Type = 0` | keduanya |
| `isElectronicEquipment` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = " ElectronicEquipment "`<br>**EDM** **[A]** `.OfferFacIn.QuotationData.BusinessType = " ElectronicEquipment "` | keduanya |
| `IsEngineering` | NB RNW EDM | **[A OR B OR C OR D OR E]** `pyWorkPage.Quotation.BusinessType = "Car"` ; `pyWorkPage.Quotation.BusinessType = "Ear"` ; `pyWorkPage.Quotation.BusinessType = "MBD"` ; `pyWorkPage.Quotation.BusinessCode = "10166"` ; `pyWorkPage.Quotation.BusinessCode = "10032"` | V1-saja (CS placeholder) |
| `IsEnvironmental` | NB RNW EDM | **[A AND B]** `pyWorkPage.Quotation.BusinessCode = 10169` ; `pyWorkPage.Quotation.BusinessType = "Aneka"` | V1-saja (CS placeholder) |
| `IsErrorSpreading` | NB RNW EDM | **[A OR B]** `pyWorkPage.Quotation.OldPolicyNo = "<nomor-polis-1 disamarkan 01-10-2026>"` ; `pyWorkPage.Quotation.OldPolicyNo = "<nomor-polis-2 disamarkan 01-10-2026>"` | keduanya |
| `IsExclusion` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "Exclusion"`<br>**EDM** **[A]** `.OfferFacIn.QuotationData.BusinessType = "Exclusion"` | keduanya |
| `IsFac` | NB RNW EDM | **[A OR B]** `.Quotation.BusinessFac = "F"` ; `.OfferFacIn.QuotationData.BusinessFac = "F"` | keduanya |
| `IsFacout` | NB RNW EDM | **[A]** `pyWorkPage.ProposalAcceptStatus = 4` | keduanya |
| `IsFacRetro` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `.IsFacRetro = 1`<br>**EDM** **[A]** `pyWorkPage.OfferFacIn.IsFacRetro = 1` | keduanya |
| `isFidelity` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "Fidelity"`<br>**EDM** **[A]** `.OfferFacIn.QuotationData.BusinessType = "Fidelity"` | keduanya |
| `IsFinanceInsurance` | NB RNW EDM | **[A OR B OR C OR D OR E OR F]** `Rule IsLimitSBondKBG evaluates to true` ; `Rule IsLimitCreditCL evaluates to true` ; `Rule IsLimitCreditNCL evaluates to true` ; `Rule IsLimitCustomBond evaluates to true` ; `Rule IsLimitTradeCredit evaluates to true` ; `pyWorkPage.OfferFacIn.QuotationData.BusinessOldId = 30` | V1-saja (CS placeholder) |
| `IsFire` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A OR B OR C OR D OR E]** `Rule IsKPR evaluates to true` ; `Rule IsOilGas evaluates to true` ; `Rule IsFireStyle1 evaluates to true` ; `Rule IsFireStyle2 evaluates to true` ; `pyWorkPage.OfferFacIn.QuotationData.BusinessType = "Fire"`<br>**EDM** **[A OR B OR C OR D OR E]** `OutData.pxResults(1).CARI2 =` salah satu dari { "Fire", "FireStyle1", "FireStyle2", "KPR", "OilGas" } | V1-saja (CS placeholder) |
| `IsFireStyle1` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "FireStyle1"`<br>**EDM** **[A]** `pyWorkPage.OfferFacIn.QuotationData.BusinessType = "FireStyle1"` | keduanya / V1-saja (CS placeholder) |
| `IsFireStyle2` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "FireStyle2"`<br>**EDM** **[A]** `pyWorkPage.OfferFacIn.QuotationData.BusinessType = "FireStyle2"` | keduanya / V1-saja (CS placeholder) |
| `IsFlagDelete` | NB RNW EDM | **[A]** `.FlagDelete = 1` | keduanya |
| `IsFlagOldData` | NB RNW EDM | **[A]** `.FlagOldData = 1` | keduanya |
| `IsFlagOnGoingPolicy` | NB RNW EDM | **[A]** `pyWorkPage.FlagOnGoingPolicy = 1` | keduanya |
| `IsFlagUW` | NB RNW | **[A]** `pyWorkPage.IsFlagUW = 1` | keduanya |
| `IsGabungan` | NB RNW EDM | **[A]** `.Policy.TypeOfPolicy = "02"` | keduanya |
| `IsGlass` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "Glass"`<br>**EDM** **[A]** `.OfferFacIn.QuotationData.BusinessType = "Glass"` | keduanya |
| `IsGolfInsurance` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "GolfInsurance"`<br>**EDM** **[A]** `OutData.pxResults(1).CARI2 = "GolfInsurance"` | keduanya |
| `IsGroup` | NB RNW EDM | **[A OR B OR C OR D OR E]** `.OfferFacIn.QuotationData.IsGroup = "Group"` ; `OperatorID.pyUserIdentifier = <nilai-identitas-disensor>` ; `OperatorID.pyUserIdentifier = <nilai-identitas-disensor>` ; `OperatorID.pyUserIdentifier = <nilai-identitas-disensor>` ; `pyWorkPage.OfferFacIn.QuotationData.MarketingCode = <nilai-identitas-disensor>` | keduanya |
| `IsGroupCreate` | NB RNW EDM | **[A OR B]** `pyWorkPage.pxCreateOperator = <nilai-identitas-disensor>` ; `pyWorkPage.pxCreateOperator = <nilai-identitas-disensor>` | keduanya |
| `IsGroupFac` | NB RNW EDM | **[A]** `pyWorkPage.FacultativePosition = "FAC"` | keduanya |
| `IsGroupUwFac` | NB RNW EDM | **[A]** `.FacultativePosition = "UWFAC"` | keduanya |
| `IsGrowingTrees` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "GrowingTrees"`<br>**EDM** **[A]** `.OfferFacIn.QuotationData.BusinessType = "GrowingTrees"` | keduanya |
| `IsHE` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "HE"`<br>**EDM** **[A]** `.OfferFacIn.QuotationData.BusinessType = "HE"` | keduanya |
| `IsHealth` | NB RNW EDM | **[A OR B OR C OR D OR E OR F OR G OR H OR I OR J OR K OR L OR M OR N OR O OR P OR Q OR R OR S OR T OR U OR V OR W OR X OR SSI OR SSIS OR SSGS OR SSES]** `.Quotation.BusinessCode =` salah satu dari { "10016", "10024", "10057", "10078", "10101", "10135", "10154", "41", "61", "63", "66", "68", "72", "75", "78", "79", "81", "82", "84", "85", "86", "87", "88", "89", "99", "S8", "T1", "TA" } | keduanya |
| `IsHealthCorporate` | NB RNW EDM | **[A OR B]** `.Quotation.BusinessCode = "S7"` ; `.Quotation.BusinessCode = "10057"` | keduanya |
| `IsHealthIndividu` | NB RNW EDM | **[A OR B OR C OR D OR E OR F OR G OR H OR I OR J OR K OR L OR M OR N OR O OR P OR Q OR R OR S OR T OR U OR V OR W OR X OR Y]** `.Quotation.BusinessCode =` salah satu dari { "10078", "10079", "10101", "10135", "41", "61", "63", "66", "68", "72", "75", "78", "79", "81", "82", "84", "85", "86", "87", "88", "89", "99", "S8", "T1", "TA" } | keduanya |
| `IsHealthSSC` | EDM | **[A]** `.Quotation.BusinessCode = "10057"` | V1-saja (CS placeholder) |
| `IsIndividual` | NB RNW EDM | **[A]** `.CustomerType = "1"` | keduanya |
| `IsInputFacRetro` | NB RNW EDM | **[A]** `pyWorkPage.OfferFacIn.IsInputFacRetro = 1` | keduanya |
| `isInterestInsuredOtherObjAneka` | NB RNW EDM | **[A OR B OR C OR D OR E OR F OR G OR H OR I OR J OR K]** `pyWorkPage.Quotation.BusinessCode =` salah satu dari { "05", "06", "15", "17", "18", "19", "32", "B7", "SK", "SL", "SQ" } | keduanya |
| `IsIT` | NB RNW | **[A OR B]** `OperatorID.pyPosition = "IT Developer"` ; `OperatorID.pyOrgDivision = "IT"` | keduanya |
| `IsKPR` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "KPR"`<br>**EDM** **[A]** `pyWorkPage.OfferFacIn.QuotationData.BusinessType = "KPR"` | keduanya / V1-saja (CS placeholder) |
| `IsLandRig` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "LandRig"`<br>**EDM** **[A]** `.OfferFacIn.QuotationData.BusinessType = "LandRig"` | keduanya |
| `IsLiability` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A OR B OR C OR D OR E OR F OR G]** `pyWorkPage.Quotation.BusinessType = "Liability"` ; `Rule IsProductsLiability evaluates to true` ; `Rule IsProfessionalLiability evaluates to true` ; `Rule IsWorkmenCompensation evaluates to true` ; `pyWorkPage.Quotation.BusinessCode = "10048"` ; `pyWorkPage.Quotation.BusinessCode = "10184"` ; `Rule IsBillboardNeon evaluates to true`<br>**EDM** **[A OR B OR C OR D OR E OR F OR G]** `.OfferFacIn.QuotationData.BusinessType = "Liability"` ; `Rule IsProductsLiability evaluates to true` ; `Rule IsProfessionalLiability evaluates to true` ; `Rule IsWorkmenCompensation evaluates to true` ; `pyWorkPage.Quotation.BusinessCode = "10048"` ; `pyWorkPage.Quotation.BusinessCode = "10184"` ; `Rule IsBillboardNeon evaluates to true` | keduanya |
| `IsLife` | NB RNW EDM | **[A OR B OR C OR D OR E OR F OR G OR H OR I OR J OR K OR L OR M OR N OR O OR P]** `pyWorkPage.Quotation.BusinessOldId =` salah satu dari { "L1", "L10", "L11", "L12", "L13", "L14", "L15", "L16", "L2", "L3", "L4", "L5", "L6", "L7", "L8", "L9" } | keduanya |
| `IsLimitCreditCL` | NB RNW EDM | **[A OR B OR C OR D OR E OR F OR G OR H OR I OR J OR K]** `pyWorkPage.OfferFacIn.QuotationData.BusinessOldId =` salah satu dari { C2, C3, C4, C5, C6, C7, C8, C9, D1, D2, D3 } | V1-saja (CS placeholder) |
| `IsLimitCreditNCL` | NB RNW EDM | **[A OR B OR C]** `pyWorkPage.OfferFacIn.QuotationData.BusinessOldId =` salah satu dari { 46, 47, 48 } | V1-saja (CS placeholder) |
| `IsLimitCustomBond` | NB RNW EDM | **[A OR B OR C OR D OR E OR F OR G OR H OR I OR J]** `pyWorkPage.OfferFacIn.QuotationData.BusinessOldId =` salah satu dari { 79, 80, 81, 82, 83, 84, 85, 86, 87, 88 } | V1-saja (CS placeholder) |
| `IsLimitSBondKBG` | NB RNW EDM | **[A OR B OR C OR D OR E OR F OR G OR H OR I OR J OR K OR L OR M OR N OR O OR P OR Q OR R OR S OR T OR U OR V OR W OR X OR Y OR Z OR  AA OR AB OR AC OR AD OR AE OR AF OR AG OR AH OR AI OR AJ OR AK OR AL OR AM OR AN OR AO]** `pyWorkPage.OfferFacIn.QuotationData.BusinessOldId =` salah satu dari { 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 95, 96, 97, 98, 99, A1, A2, A3, A4, A5, A6, A7, A8, A9, B4, B5, B6, B8, B9, C1, D4, D5, D6, D7, D8, D9, F1, F2, F3, F4, F6 } | V1-saja (CS placeholder) |
| `IsLimitTradeCredit` | NB RNW EDM | **[A OR B OR C OR D]** `pyWorkPage.OfferFacIn.QuotationData.BusinessOldId =` salah satu dari { 49, 50, 51, 52 } | V1-saja (CS placeholder) |
| `IsMaintenance` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "Maintenance"`<br>**EDM** **[A]** `.OfferFacIn.QuotationData.BusinessType = "Maintenance"` | keduanya |
| `IsMarineCargo` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[G OR H]** `.OfferFacIn.QuotationData.BusinessType = "MarineCargo"` ; `pyWorkPage.Quotation.BusinessType = "MarineCargo"`<br>**EDM** **[G]** `OutData.pxResults(1).CARI2 = "MarineCargo"` | keduanya |
| `IsMarineHull` | NB RNW EDM | **[A]** `pyWorkPage.Quotation.BusinessType = "MarineHull"` | keduanya |
| `IsMarineHullOffshore` | NB RNW EDM | **[A AND (B OR C)]** `pyWorkPage.Quotation.BusinessType = "Aneka"` ; `pyWorkPage.Quotation.BusinessCode = "10157"` ; `pyWorkPage.Quotation.BusinessName = "MARINE HULL OFFSHORE"` | keduanya |
| `IsMaterialDamage` | EDM | **[A]** `pyWorkPage.PropertyList(1).IsMaterialDamage = false` | keduanya |
| `IsMBD` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "MBD"`<br>**EDM** **[A]** `.OfferFacIn.QuotationData.BusinessType = "MBD"` | keduanya |
| `IsMBU` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A OR B]** `pyWorkPage.OfferFacIn.QuotationData.BusinessType = "MBUCar"` ; `pyWorkPage.OfferFacIn.QuotationData.BusinessType = "MBUMotorCycle"`<br>**EDM** **[A OR B OR C OR D]** `OutData.pxResults(1).CARI2 =` salah satu dari { "MBUCar", "MBUMotorCycle" } | keduanya / V1-saja (CS placeholder) |
| `IsMBUCar` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "MBUCar"`<br>**EDM** **[A]** `.Quotation.BusinessType = "MBUCar"` | keduanya |
| `IsMultiCOB` | NB RNW EDM | **[A AND B]** `.OfferFacIn.QuotationData.IsMultiCoB = 1` ; `.OfferFacIn.QuotationData.BusinessOldId = 22` | keduanya |
| `IsMultipleObyekMBU` | EDM | **[A]** `@Utilities.SizeOfPropertyList(.VehicleList) > 1` | keduanya |
| `IsNB` | NB RNW EDM | **[A]** `pyWorkPage.Quotation.StatusBusiness = 1` | keduanya |
| `isNew` | EDM | **[A]** `[first value] [relation] [second value]` | V1-saja (CS kosong) |
| `IsNoDeduc` | NB RNW EDM | **[A OR B OR C OR D OR E OR F OR G OR H OR I OR J OR K OR L]** `.CoverageID =` salah satu dari { "100846", "100848", "100849", "100850", "100855", "100856", "100874", "100875", "100876", "100877", "100878", "100879" } | keduanya |
| `IsNonPropertyandNonEngineering` | NB RNW EDM | **[!A AND !B]** `Rule IsFire evaluates to true` ; `Rule IsEngineering evaluates to true` | V1-saja (CS placeholder) |
| `IsNotActive` | EDM | **[A]** `@DateTime.CompareDates(pyWorkPage.Quotation.EdmDate,.EndDate)` | keduanya |
| `IsNotAdmin` | NB | **[A]** `OperatorID.pyPosition != "Admin"` | keduanya |
| `IsNotDirect` | NB RNW EDM | **[A AND B AND C AND D AND E AND F AND G AND H]** `pyWorkPage.Quotation.SobLeader1 != "10000137"` ; `pyWorkPage.Quotation.SobLeader1 != "10000711"` ; `pyWorkPage.Quotation.SobLeader1 != "10001756"` ; `pyWorkPage.Quotation.SobLeader1 != "10002313"` ; `pyWorkPage.Quotation.SobLeader1 != "10004812"` ; `pyWorkPage.Quotation.SobLeader1 != "10027341"` ; `pyWorkPage.Quotation.SobLeader1 != "10027720"` ; `pyWorkPage.Quotation.SobLeader0 != "10001692"` | keduanya |
| `IsNotEDM` | NB RNW EDM | **[A]** `pyWorkPage.Quotation.StatusBusiness != 3` | keduanya |
| `IsNotFire` | EDM | **[A AND B]** `pyWorkPage.Policy.Quotation.BusinessType != "FireGeneral"` ; `pyWorkPage.Policy.Quotation.BusinessType != "Fire"` | keduanya |
| `IsNotPAandNotMBU` | NB RNW EDM | **[!A AND !B]** `Rule IsPA evaluates to true` ; `Rule IsMBU evaluates to true` | V1-saja (CS placeholder) |
| `IsNotPrintRISlip` | NB RNW EDM | **[A]** `pyWorkPage.OfferFacIn.IsRISlip != 1` | keduanya |
| `IsObjectDeleteVisible` | NB RNW EDM | **[A OR B OR C]** `.Quotation.BusinessCode =` salah satu dari { "10039", "10040", "10041" } | keduanya |
| `IsObjectOtherScheduleVisible` | NB RNW EDM | **[A OR B OR C OR D OR E OR F OR G OR V OR H OR I OR J OR K OR L OR M OR N OR O OR P OR Q OR R OR S OR T OR U OR W OR X]** `pyWorkPage.Quotation.BusinessCode =` salah satu dari { "10020", "10031", "10032", "10033", "10034", "10035", "10037", "10038", "10039", "10042", "10043", "10051", "10105", "10107", "10109", "10115", "10118", "10120", "10121", "10124", "10140", "10148", "10151" } | V1-saja (CS placeholder) |
| `IsObjectParticipantVisible` | NB RNW EDM | **[A]** `pyWorkPage.Quotation.BusinessCode = "10042"` | V1-saja (CS placeholder) |
| `IsObjectPolicyScheduleVisible` | NB RNW EDM | **[A OR B OR C OR D]** `pyWorkPage.Quotation.BusinessCode =` salah satu dari { "10031", "10115", "10138", "10140" } | V1-saja (CS placeholder) |
| `IsObjectSectionAneka` | NB RNW EDM | **[A OR B OR C OR D OR E OR F OR G OR H OR I OR J OR K OR L]** `pyWorkPage.Quotation.BusinessCode =` salah satu dari { "07", "15", "17", "18", "24", "27", "32", "90", "91", "SV", "SX", "T9" } | keduanya |
| `IsObjectWithQuantityYear` | NB RNW EDM | **[A OR B OR C OR D]** `pyWorkPage.Quotation.BusinessType =` salah satu dari { "Boiler", "ContractorsPM", "LandRig", "MBD" } | keduanya |
| `IsOfferFacIn` | NB | **[A]** `.Quotation.BusinessFac = "F"` | keduanya |
| `IsOilGas` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "OilGas"`<br>**EDM** **[A]** `pyWorkPage.OfferFacIn.QuotationData.BusinessType = "OilGas"` | keduanya |
| `IsOperatorLife` | NB | **[A]** `OperatorID.pyWorkGroup = "ReasLife"` | keduanya |
| `IsPA` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "PA"`<br>**EDM** **[A]** `OutData.pxResults(1).CARI2 = "PA"` | keduanya |
| `IsPASSG` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `@Utilities.SizeOfPropertyList(pyWorkPage.PersonListPA) > 0`<br>**EDM** **[A]** `@Utilities.SizeOfPropertyList(.PersonListPA) > 0` | V1-saja (CS placeholder) |
| `IsPEGAPROD` | NB RNW EDM | **[A]** `pxProcess.pzProductionLevel = "5"` | keduanya |
| `IsPKSASM` | NB RNW EDM | **[A]** `pyWorkPage.OfferFacIn.IsB2B = "ASM"` | V1-saja (CS placeholder) |
| `IsPKSEmpty` | NB RNW EDM | **[A]** `.Quotation.TemplateNo != ""` | keduanya |
| `IsPLL` | NB RNW EDM | **[A OR B OR C OR D]** `.Coverage =` salah satu dari { "100850", "100875", "100878", "100879" } | keduanya |
| `IsPolicyPeriodeNotEmpty` | NB RNW EDM | **[A]** `pyWorkPage.Policy.StartDateTime != ""` | keduanya |
| `IsPortRisk` | NB RNW EDM | **[A]** `pyWorkPage.OfferFacIn.QuotationData.BusinessCode = "10188"` | V1-saja (CS placeholder) |
| `IsProductsLiability` | NB RNW EDM | **[A AND B]** `pyWorkPage.Quotation.BusinessType = "Aneka"` ; `pyWorkPage.Quotation.BusinessCode = "10021"` | keduanya |
| `IsProfessionalLiability` | NB RNW EDM | **[A AND B]** `pyWorkPage.Quotation.BusinessType = "Aneka"` ; `pyWorkPage.Quotation.BusinessCode = "10023"` | keduanya |
| `IsPropertyandEngineering` | NB RNW EDM | **[A OR B]** `Rule IsFire evaluates to true` ; `Rule IsEngineering evaluates to true` | V1-saja (CS placeholder) |
| `IsProposalTransfer` | NB EDM | **BERCABANG** - **NB** **[(A OR B) AND C]** `pyWorkPage.ProposalTransfer = "1"` ; `pyWorkPage.ProposalTransfer = "2"` ; `pyWorkPage.ProposalPosition != "1"`<br>**EDM** **[(A OR B) AND C]** `.ProposalTransfer = "1"` ; `.ProposalTransfer = "2"` ; `.ProposalPosition != "1"` | keduanya |
| `IsRenewal` | NB RNW EDM | **[A]** `pyWorkPage.Quotation.StatusBusiness = 2` | keduanya |
| `IsRequired` | EDM | **[A]** `.ClauseType = "TC Policy"` | keduanya |
| `IsRiskIDNull` | NB RNW EDM | **[A AND B AND C]** `.Property.AlmRiskID = ""` ; `.Property.RiskLocation.ASMZipCode =` ; `.Property.RoadName =` | keduanya |
| `IsRISlip` | NB RNW EDM | **[A]** `pyWorkPage.OfferFacIn.IsRISlip = 1` | keduanya |
| `isSellingModeB2B` | NB | **[A]** `@(Pega-RULES:Utilities).getDataSystemSetting("PegaCRM-", "SellingMode") EQUALS "B2B"` | keduanya |
| `isSellingModeB2BB2C` | NB | **[A]** `@(Pega-RULES:Utilities).getDataSystemSetting("PegaCRM-", "SellingMode") equals "B2B_B2C"` | keduanya |
| `isSellingModeB2C` | NB | **[A]** `@(Pega-RULES:Utilities).getDataSystemSetting("PegaCRM-", "SellingMode") equals "B2C"` | keduanya |
| `IsShortPeriod` | EDM | **[A]** `.CalculateMethod = 2` | keduanya |
| `IsShowInput` | NB RNW EDM | **[A OR B]** `pyPortal.IsEdit` ; `OutputData.HASIL1 != ""` | V1-saja (CS placeholder) |
| `IsSpecialCase` | NB RNW EDM | **[A]** `pyWorkPage.IsSpecialCase = true` | keduanya |
| `IsSpreadingDepan` | EDM | **[A OR B OR C OR D]** `Rule IsPA evaluates to true` ; `Rule IsTravel evaluates to true` ; `Rule IsMBU evaluates to true` ; `Rule IsFire evaluates to true` | V1-saja (CS placeholder) |
| `IsSpreadingUW` | NB RNW EDM | **[A]** `pyWorkPage.IsSpreadingUW = "1"` | keduanya |
| `IsSPVCreate` | NB | **[A OR B]** `pyWorkPage.pxCreateOperator = <nilai-identitas-disensor>` ; `pyWorkPage.pxCreateOperator = <nilai-identitas-disensor>` | V1-saja (CS placeholder) |
| `IsSPVTreaty1` | NB | **[A]** `OperatorID.pyTelephone = <nilai-identitas-disensor>` | V1-saja (CS placeholder) |
| `IsSuccessHitService` | NB RNW EDM | **[(A AND B) OR (A AND C)]** `pyWorkPage.StatusService.StsKonversiFacIn = 1` ; `pyWorkPage.StatusService.StsKonversiFacOut = ""` ; `pyWorkPage.StatusService.StsKonversiFacOut = 1` | keduanya |
| `IsT1T3` | NB RNW EDM | **[A OR (B AND C)]** `pyWorkPage.OfferFacIn.QuotationData.TeamGroup = 1` ; `pyWorkPage.OfferFacIn.QuotationData.TeamGroup = 2` ; `pyWorkPage.OfferFacIn.IsB2B = "ASM"` | V1-saja (CS placeholder) |
| `IsT1T4` | NB RNW EDM | **[A OR B]** `pyWorkPage.OfferFacIn.QuotationData.TeamGroup = 1` ; `pyWorkPage.OfferFacIn.QuotationData.TeamGroup = 4` | V1-saja (CS kosong) |
| `IsT2T3` | NB RNW EDM | **[A OR B]** `pyWorkPage.OfferFacIn.QuotationData.TeamGroup = 2` ; `pyWorkPage.OfferFacIn.QuotationData.TeamGroup = 3` | V1-saja (CS kosong) |
| `IsT2T4` | NB RNW EDM | **[(A AND B ) OR (C AND D )]** `pyWorkPage.OfferFacIn.QuotationData.TeamGroup = 2` ; `pyWorkPage.OfferFacIn.IsSpecialAcceptance = "false"` ; `pyWorkPage.OfferFacIn.QuotationData.TeamGroup = 2` ; `pyWorkPage.OfferFacIn.IsB2B = "ASM"` | V1-saja (CS placeholder) |
| `IsTABActive_Cedant` | NB RNW EDM | **[A]** `.TabPosition = "Cedant"` | keduanya |
| `IsTABActive_Coverage` | NB RNW EDM | **[A]** `.TabPosition = "Coverage"` | keduanya |
| `IsTABActive_Deduction` | NB RNW EDM | **[A]** `.TabPosition = "Deduction"` | keduanya |
| `IsTABActive_LossRecord` | NB RNW | **[A]** `.TabPosition = "LossRecord"` | keduanya |
| `IsTABActive_Object` | NB RNW EDM | **[A]** `.TabPosition = "Object"` | keduanya |
| `IsTABActive_Payment` | NB RNW EDM | **[A]** `.TabPosition = "Payment"` | keduanya |
| `IsTABActive_PaymentLife` | NB RNW EDM | **[A]** `.TabPosition = "PaymentLife"` | keduanya |
| `IsTABActive_ScoringRisk` | NB RNW EDM | **[A]** `.TabPosition = "ScoringRisk"` | keduanya |
| `IsTABActive_Spreading` | NB RNW EDM | **[A]** `.TabPosition = "Spreading"` | keduanya |
| `IsTABActive_SpreadingLife` | NB RNW EDM | **[A]** `.TabPosition = "SpreadingLife"` | keduanya |
| `IsTahun1` | EDM | **[A]** `.Year = 1` | keduanya |
| `IsTBonding` | NB RNW EDM | **[A OR B OR C OR D]** `.OfferFacIn.QuotationData.TeamGroup = 5` ; `.OfferFacIn.QuotationData.MarketingName = <nilai-identitas-disensor>` ; `.OfferFacIn.QuotationData.MarketingName = <nilai-identitas-disensor>` ; `.OfferFacIn.QuotationData.MarketingName = <nilai-identitas-disensor>` | V1-saja (CS placeholder) |
| `IsTJH` | NB RNW EDM | **[A OR B OR C OR D]** `.CoverageID =` salah satu dari { "100846", "100874", "100876", "100877" } | keduanya |
| `IsTJHInRange` | EDM | **[A AND B]** `Rule IsTJH evaluates to true` ; `.MinTSI <= Param.TSITJH` | V1-saja (CS placeholder) |
| `IsTJHNotInRange` | EDM | **[A AND B]** `Rule IsTJH evaluates to true` ; `.MinTSI > Param.TSITJH` | V1-saja (CS placeholder) |
| `IsTravel` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A]** `pyWorkPage.Quotation.BusinessType = "Travel"`<br>**EDM** **[A OR B]** `.Quotation.BusinessType = "Travel"` ; `pyWorkPage.OfferFacIn.QuotationData.BusinessType = "Travel"` | keduanya |
| `isTravelTime` | NB RNW EDM | **[A && B]** `pyWorkPage.Quotation.BusinessType = "Travel"` ; `pyWorkPage.FlagOnGoingPolicy = 1` | V1-saja (CS placeholder) |
| `IsTreaty1` | NB | **[A]** `OperatorID.pyTelephone = <nilai-identitas-disensor>` | V1-saja (CS placeholder) |
| `IsTreatyIn` | NB RNW EDM | **[A]** `.Quotation.BusinessFac = "T"` | keduanya |
| `IsTypeDeductType` | NB RNW EDM | **[A]** `@String.notEquals(.DeductibleType,"4")&&@String.notEquals(.DeductibleType,"7")&&@String.notEquals(.DeductibleType,"8")&&@String.notEquals(.DeductibleType,"1")&&@String.notEquals(.DeductibleType,"2")&&@String.notEquals(.DeductibleType,"5")` | keduanya |
| `IsTypeStock` | NB RNW EDM | **[A OR B OR C]** `.PropertyItemGroup =` salah satu dari { "STOCK", "STOCK(S)", "STOCKS" } | keduanya |
| `IsUW` | NB RNW EDM | **BERCABANG** - **NB/RNW** **[A OR B OR C]** `pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName =` salah satu dari { ReasFacInDirector, ReasFacInGroupLeader, ReasFacInUnderwriting }<br>**EDM** **[A OR B OR C OR D]** `pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName =` salah satu dari { ReasFacInDirector, ReasFacInGroupLeader, ReasFacInMarketing, ReasFacInUnderwriting } | keduanya |
| `IsUWAcceptance` | NB RNW EDM | **[A]** `pyWorkPage.ProposalPosition = "2"` | keduanya |
| `IsUWAccepted` | NB EDM | **[A]** `.ProposalAcceptStatus = "1"` | keduanya |
| `IsUWSpread` | NB RNW EDM | **[A OR B OR C OR D OR E OR F]** `pyWorkPage.PositionNote =` salah satu dari { "ReasFacInDepHeadUnderwriting", "ReasFacInJuniorUnderwriting", "ReasFacInJuniorUnderwritingA", "ReasFacInSeniorUnderwriting", "ReasFacInUnderwriting", "ReasFacInUnderwritingFinancial" } | V1-saja (CS placeholder) |
| `IsValueDeductType` | NB RNW EDM | **[A OR B OR C OR D OR E]** `.DeductibleType =` salah satu dari { "100414", "100416", "100417", "100418", "100419" } | keduanya |
| `IsVisible` | NB RNW EDM | **[A]** `IsVisible = True` | keduanya |
| `IsWhetherIndex` | NB RNW EDM | **[A AND (B OR C)]** `pyWorkPage.Quotation.BusinessType = "Aneka"` ; `pyWorkPage.Quotation.BusinessCode = "10165"` ; `pyWorkPage.Quotation.BusinessName = "WHETHER INDEX"` | keduanya |
| `IsWorkmenCompensation` | NB RNW EDM | **[A]** `pyWorkPage.Quotation.BusinessType = "Workmen"` | keduanya |
| `IsYieldShortfall` | NB RNW EDM | **[A AND B]** `pyWorkPage.Quotation.BusinessType = "Aneka"` ; `pyWorkPage.Quotation.BusinessCode = 10167` | keduanya |
| `LetterNoNull` | NB RNW | **[A]** `pyWorkPage.LetterNo = ""` | V1-saja (CS kosong) |
| `NopolisEmpty` | NB | **[A]** `pyWorkPage.PolicyTreatyIn.PolicyNo = ""` | keduanya |
| `pyIsIpadOrDesktop` | NB | **[A OR B]** `Rule pyIsIPad evaluates to true` ; `pxRequestor.pxDeviceType = desktop` | V1-saja (CS placeholder) |
| `pyIsMobile` | NB | **[I1]** `@(Pega-RulesEngine:Utilities).pzIsMobile(tools)` | V1-saja (CS placeholder) |
| `recordEvent` | EDM | **[A0]** `[first value] [relation] [second value]` | V1-saja (CS kosong) |
| `StepStatusFail` | NB EDM | **[a]** `[first value] [relation] [second value]` | V1-saja (CS kosong) |
| `ToDepHeadUW` | NB RNW EDM | **[A AND ( B OR C )]** `pyWorkPage.LetterNo = "DEPHEADUNDERWRITER"` ; `pyWorkPage.OfferFacIn.IsB2B != "ASM"` ; `pyWorkPage.OfferFacIn.IsB2B = ""` | keduanya |
| `ToDeptHeadUWLife` | NB | **[A]** `pyWorkPage.LetterNo = "DEPTHEADUWLIFE"` | keduanya |
| `ToDirMarketing` | NB RNW EDM | **[A]** `pyWorkPage.LetterNo = "DIREKTURMARKETING"` | keduanya |
| `ToDirTeknik` | NB RNW EDM | **[A]** `pyWorkPage.LetterNo = "DIREKTURTEKNIK"` | keduanya |
| `ToJUW_A` | NB RNW EDM | **[A]** `pyWorkPage.LetterNo = "JUW_A"` | keduanya |
| `ToKadivFacultative` | NB RNW EDM | **[A]** `pyWorkPage.LetterNo = "KADIVFACULTATIVE"` | keduanya |
| `ToKadivFin` | NB RNW EDM | **[A]** `pyWorkPage.LetterNo = "KADIVFINANCIAL"` | keduanya |
| `ToKadivTeknik` | NB RNW EDM | **[A]** `pyWorkPage.LetterNo = "KADIVTEKNIK"` | keduanya |
| `ToManagerTeknik` | NB RNW EDM | **[A]** `pyWorkPage.LetterNo = "MANAGERTEKNIK"` | keduanya |
| `ToSeniorUW` | NB RNW EDM | **[A]** `pyWorkPage.LetterNo = "SENIORUW"` | keduanya |
| `ToTREATYDEPTHEAD` | NB | **[A]** `pyWorkPage.LetterNo = "TREATYINDEPTHEAD"` | keduanya |
| `ToUW` | NB RNW EDM | **[A]** `pyWorkPage.LetterNo = "UNDERWRITER"` | keduanya |
| `TreatyMasterInEDM` | NB | **[A OR B OR C]** `pyWorkPage.TreatyIn.EDMState =` salah satu dari { "1", "2", "3" } | V1-saja (CS placeholder) |


### 4.1 Indeks tema — [dugaan]

Pengelompokan di bawah dihitung otomatis dari **isi kondisi** (bukan dari nama rule), dengan urutan
prioritas kata kunci yang dicantumkan penuh di bawah agar dapat diaudit ulang. Klasifikasinya
berlabel **[dugaan]**: ia mencerminkan properti apa yang diuji, bukan pernyataan resmi tentang
maksud bisnis rule.

```powershell
# prioritas dievaluasi dari atas ke bawah; $txt = seluruh V1Strings + V1Exprs untuk nama rule tsb
if     ($txt -match 'pyUserIdentifier|pxCreateOperator|pyTelephone|pyPosition|pyOrgDivision|pyWorkGroup') { '4-Guard identitas' }
elseif ($txt -match 'pyWorkBasketName|\.LetterNo|\.PositionNote|pyWorkBasket')                            { '2-Routing approver' }
elseif ($txt -match 'EdmType|StatusBusiness|QuotationData\.Type|TeamGroup|FlagOldData|ProdKe|IsCedingConfirm|pyWorkIDPrefix|StepStatus|pyID') { '3-Diskriminator siklus' }
elseif ($txt -match 'BusinessType|BusinessCode|BusinessOldId|CARI2|BusinessName')                         { '1-Lini bisnis' }
elseif ($txt -match '\.Coverage|CoverageID|DeductibleType|Deduc|Premium|Rate|TSI|Spread|Discount|Komisi|Brokerage') { '5-Cakupan/deduksi/premi' }
elseif ($txt -match 'TabPosition|pyPortal|Visible|IsEdit|pxDeviceType|pzIsMobile|hasPrimaryPage|IsAddButton|IsRequired') { '6-UI' }
else { '7-Lain-lain' }
```

| Tema | Jumlah |
| --- | --- |
| 1 — Lini bisnis (COB) | 80 |
| 2 — Routing approver / jabatan | 16 |
| 3 — Diskriminator siklus / status | 29 |
| 4 — Guard identitas / atribut operator | 8 |
| 5 — Cakupan / deduksi / premi / spreading | 17 |
| 6 — UI / visibilitas layar | 14 |
| 7 — Lain-lain | 62 |
| **Total** | **226** |

**Tema 2 — Routing approver (16):** `IsAdmin`, `IsUW`, `IsUWSpread`, `LetterNoNull`, `ToDepHeadUW`,
`ToDeptHeadUWLife`, `ToDirMarketing`, `ToDirTeknik`, `ToJUW_A`, `ToKadivFacultative`, `ToKadivFin`,
`ToKadivTeknik`, `ToManagerTeknik`, `ToSeniorUW`, `ToTREATYDEPTHEAD`, `ToUW`.

**Tema 3 — Diskriminator siklus/status (29):** `FlagOldData`, `IsCedingConfirmOffer`,
`IsCedingConfirmPolicy`, `IsClaim`, `IsEDM`, `IsEdmAddObject`, `IsEdmAdjCeding`, `IsEdmAdjCurrency`,
`IsEdmAdjInsured`, `IsEdmAdjPeriod`, `IsEdmAdjRate`, `IsEdmAdjRefNo`, `IsEdmAdjRIC`,
`IsEdmAdjShareCedant`, `IsEdmAdjSpreading`, `IsEdmAdjTSI`, `IsEdmExtendPeriod`, `IsEdmPPNPPH`,
`IsEDMRiSlip`, `IsFlagOldData`, `IsNB`, `isNew`, `IsNotEDM`, `IsRenewal`, `IsT1T3`, `IsT1T4`,
`IsT2T3`, `IsT2T4`, `IsTBonding`.

**Tema 4 — Guard identitas (8):** `IsGroup`, `IsGroupCreate`, `IsIT`, `IsNotAdmin`,
`IsOperatorLife`, `IsSPVCreate`, `IsSPVTreaty1`, `IsTreaty1` — rincian di §7.

**Tema 6 — UI (14):** `hasPrimaryPage`, `IsAddButton`, `IsButtonShow`, `IsShowInput`,
`IsTABActive_Cedant`, `IsTABActive_Coverage`, `IsTABActive_LossRecord`, `IsTABActive_Object`,
`IsTABActive_Payment`, `IsTABActive_PaymentLife`, `IsTABActive_ScoringRisk`, `IsVisible`,
`pyIsIpadOrDesktop`, `pyIsMobile`.

Tema 1, 5, dan 7 selengkapnya dapat dibaca langsung dari tabel §4 (kolom kondisi menampilkan
properti yang menentukan klasifikasinya).

### 4.2 Rantai rule → rule

245 dari 1.863 baris kondisi berbentuk `Rule <X> evaluates to true`
(`pyConditionValue1Purpose` = `Rule-Obj-When`, ekspresi
`@(Pega-RULES:ExpressionEvaluators).evaluateWhen("<X>")`) **[terverifikasi]**:

```powershell
$refs = $c | Where-Object { $_.Prop -match '^Rule\s+(.+?)\s+evaluates to true$' }
"Baris : " + $refs.Count
$names = $refs | ForEach-Object { if ($_.Prop -match '^Rule\s+(.+?)\s+evaluates to true$') { $Matches[1] } }
"Rule dirujuk (unik) : " + ($names | Sort-Object -Unique).Count
($names | Sort-Object -Unique) | Where-Object { $have -notcontains $_ }
```

| Ukuran | Nilai |
| --- | --- |
| Total baris kondisi (semua berkas) | **1.863** |
| Baris berbentuk `Rule X evaluates to true` | **245** |
| Rule yang dirujuk (unik) | **63** |
| Baris berbentuk `Rule X evaluates to **false**` | **0** |
| Rule dirujuk yang **tidak ada** di korpus `When\` | **1** → `pyIsIPad` (dirujuk oleh `pyIsIpadOrDesktop`) |

Grafnya praktis tertutup: hanya satu dependensi hilang, dan itu rule platform Pega
(`pyIsIPad`), bukan rule bisnis. **[dugaan]** rule platform disediakan produk Pega dan tidak
diekspor bersama ruleset aplikasi.

### 4.3 Label kondisi tidak kontinu — sisa penghapusan baris

Beberapa rule memakai `pyLogic` yang melompati huruf, mis. `IsAneka` (NB):
`A OR B OR C OR E OR F …` (tanpa `D`), atau `IsMarineCargo` (NB): `G OR H` (mulai dari `G`).

Audit menunjukkan **tidak ada baris kondisi yatim** — setiap `<pyConditionLabel>` yang ada selalu
dirujuk `pyLogic`, dan sebaliknya **[terverifikasi]**:

```powershell
$dead = foreach ($r in $d) {
  $labs = @($r.CondLabels -split ',' | Where-Object { $_ -ne '' })
  $lg = ' ' + ($r.Logic -replace '[()!]',' ') + ' '
  $unused = @(); foreach ($l in $labs) { if ($lg -notmatch ('(?i)\s' + [regex]::Escape($l) + '\s')) { $unused += $l } }
  if ($unused.Count -gt 0) { $r.FileName }
}
"Berkas dengan baris kondisi mati: " + (@($dead)).Count      # -> 0
```

Huruf yang hilang adalah bekas baris yang **dihapus** dari rule, bukan baris yang masih ada tetapi
tidak dievaluasi. Implementasi Go boleh mengabaikan penomorannya dan mengurutkan ulang.

---

## 5. Rule yang kondisinya TIDAK terbaca sama sekali → wajib `panic()`

**Jumlahnya NOL.** Seluruh 601 berkas memiliki `<pyConditionValue1String>` **dan**
`<pyConditionValue1>` yang terisi **[terverifikasi]**:

```powershell
"V1StrCount = 0  : " + ($d | Where-Object { [int]$_.V1StrCount -eq 0 }).Count   # -> 0
"V1ExprCount = 0 : " + ($d | Where-Object { [int]$_.V1ExprCount -eq 0 }).Count  # -> 0
"Kedua sumber tak terbaca : " + ($d | Where-Object {
  ($_.CondStrings -replace '\|\|\|','').Trim() -eq '' -and
  ($_.V1Strings   -replace '\|\|\|','').Trim() -eq '' }).Count                  # -> 0
```

### 5.1 Koreksi terhadap catatan awal work owner

Lima rule pernah dinyatakan "kondisi kosong di ekspor" dan diarahkan menjadi `panic()`. **Keempat
hingga kelimanya terbaca penuh** **[terverifikasi]**:

| Rule | Kondisi sesungguhnya | Berkas + baris |
| --- | --- | --- |
| `IsPKSASM` | `pyWorkPage.OfferFacIn.IsB2B = "ASM"` | `NB FacIn\When\IsPKSASM.xml` L336 (juga di RNW & EDM, ekspresi identik) |
| `ToUW` | `pyWorkPage.LetterNo = "UNDERWRITER"` | `NB FacIn\When\ToUW.xml` L158 **dan** L313 — kedua tag terisi |
| `ToJUW_A` | `pyWorkPage.LetterNo = "JUW_A"` | `NB FacIn\When\ToJUW_A.xml` L159 **dan** L303 — kedua tag terisi |
| `LetterNoNull` | `pyWorkPage.LetterNo = ""` | `NB FacIn\When\LetterNoNull.xml` L301 (`<pyConditionString/>` self-closing di L163) |
| `IsEdmInternalRetro` | `pyWorkPage.OfferFacIn.QuotationData.EndorsementInternalRetro = 1` | `Endorsment Fac In\When\IsEdmInternalRetro.xml` L363 (`<pyConditionString/>` di L150) |

`ToUW` dan `ToJUW_A` bahkan **tidak pernah kosong** — `<pyConditionString>` keduanya terisi.
Kesimpulan "kondisi belum diketahui" untuk kelima rule ini **dibantah oleh korpus**.

**Konsekuensi untuk implementasi:** aturan CLAUDE.md §4.5 yang memerintahkan
`panic("IsPKSASM: kondisi belum diketahui")` **tidak lagi berdasar** untuk kelima rule tersebut.
Aturan §4.5 sendiri tetap berlaku sebagai prinsip — hanya daftar rule-nya yang kosong sekarang.

### 5.2 Empat berkas yang label kondisinya belum ter-resolve (bukan kondisi kosong)

```powershell
Select-String -Path "D:\migrasi\RNM\*\When\*.xml" -Pattern '<pyConditionValue1String>\[first value\]' -List |
  ForEach-Object { $_.Path }
```

| Berkas | `pyConditionValue1String` | `pyConditionValue1` (yang dieksekusi) |
| --- | --- | --- |
| `Endorsment Fac In\When\isNew.xml` | `[first value] [relation] [second value]` | `compareTwoValues(.pyID, "=", "")` |
| `Endorsment Fac In\When\recordEvent.xml` | idem | `compareTwoValues(true, "=", @HavePrivilege(tools, "pyShowClientTimer", "@baseclass", null))` |
| `Endorsment Fac In\When\StepStatusFail.xml` | idem | `compareTwoValues(Lib(Pega-RULES:String).inString(pxThread.pxMethodStatus, "Fail"), "=", 0)` |
| `NB FacIn\When\StepStatusFail.xml` | idem | identik dengan baris di atas |

Kondisinya **terbaca** lewat `pyConditionValue1`; yang belum ter-resolve hanyalah label tampilannya.

> **[pertanyaan terbuka]** `StepStatusFail` membandingkan `inString(pxThread.pxMethodStatus,"Fail")`
> dengan `0`. `inString` Pega mengembalikan **posisi** substring. Nilai kembali `0` berarti
> "ditemukan pada indeks 0" atau "tidak ditemukan"? Arah predikatnya terbalik bergantung jawabannya,
> dan nama rule (`…Fail`) bukan bukti. Korpus `When\` tidak memuat definisi `inString`.

> **[pertanyaan terbuka]** `recordEvent` menggerbangi atas privilege `pyShowClientTimer`
> (`@HavePrivilege`). Privilege ini tidak didefinisikan di korpus `When\`; siapa yang memilikinya
> tidak dapat ditentukan.

### 5.3 `pyTempText` BUKAN penanda kondisi setengah tersunting

```powershell
"Berkas dengan <pyTempText>true</pyTempText> : " + ((Select-String -Path "D:\migrasi\RNM\*\When\*.xml" `
  -Pattern '<pyTempText>true</pyTempText>' -List) | Measure-Object).Count
```

| Ukuran | Nilai |
| --- | --- |
| Berkas dengan `<pyTempText>true</pyTempText>` | **294** dari 601 (48,9 %) |
| Berkas yang `pyConditionValue1String`-nya benar-benar belum ter-resolve | **4** |

Di `IsPKSASM`, `<pyTempText>true</pyTempText>` berdampingan dengan
`<pyConditionValue1StringLabel>[first value][relation][second value]</pyConditionValue1StringLabel>`
sementara `<pyConditionValue1String>` di baris lain **terisi penuh** **[terverifikasi]**. Tag itu
menandai *label* yang belum di-render, bukan kondisi yang belum ditulis. Memakai `pyTempText`
sebagai tanda bahaya akan menandai separuh korpus secara keliru.

---

## 6. Rule yang kondisinya TERSEMBUNYI — dilewatkan discovery yang hanya membaca `pyConditionString`

```powershell
$hid = $d | Where-Object { $cs = ($_.CondStrings -replace '\|\|\|','').Trim()
                           $cs -eq '' -or $cs -eq '[Double click to add condition]' }
"Berkas : " + $hid.Count
"Nama rule terdampak : " + ($hid.FileName | Sort-Object -Unique).Count
```

| Kategori | Berkas | Nama rule unik |
| --- | --- | --- |
| `<pyConditionString/>` **self-closing / kosong** | **13** | **7** |
| `<pyConditionString>[Double click to add condition]</pyConditionString>` | **157** | **64** |
| **Total tersembunyi** | **170** (28,3 % dari 601) | **71** (31,4 % dari 226) |

### 6.1 Tujuh rule dengan `<pyConditionString/>` kosong

| Rule | Folder | `pyLogic` | Kondisi dari `pyConditionValue1String` |
| --- | --- | --- | --- |
| `IsEdmInternalRetro` | EDM | `A` | `pyWorkPage.OfferFacIn.QuotationData.EndorsementInternalRetro = 1` |
| `isNew` | EDM | `A` | label belum resolve → ekspresi `compareTwoValues(.pyID, "=", "")` |
| `IsT1T4` | NB, RNW, EDM | `A OR B` | `pyWorkPage.OfferFacIn.QuotationData.TeamGroup = 1` ‖ `… = 4` |
| `IsT2T3` | NB, RNW, EDM | `A OR B` | `pyWorkPage.OfferFacIn.QuotationData.TeamGroup = 2` ‖ `… = 3` |
| `LetterNoNull` | NB, RNW | `A` | `pyWorkPage.LetterNo = ""` |
| `recordEvent` | EDM | `A0` | label belum resolve → `compareTwoValues(true, "=", @HavePrivilege(tools,"pyShowClientTimer","@baseclass",null))` |
| `StepStatusFail` | NB, EDM | `a` | label belum resolve → `compareTwoValues(Lib(Pega-RULES:String).inString(pxThread.pxMethodStatus,"Fail"), "=", 0)` |

Catatan `IsT1T4` / `IsT2T3`: `<pyConditionViewer>` hanya memuat **satu** `rowdata` (kosong)
sementara `<pyCondition>` memuat **dua** baris terisi — viewer-nya basi terhadap kondisi
sesungguhnya **[terverifikasi: `NB FacIn\When\IsT1T4.xml` L163 `<pyConditionString/>`, sementara
dua `<pyConditionValue1String>` terisi]**.

### 6.2 Enam puluh empat rule dengan `pyConditionString` = placeholder

`crmCreateOpportunity` (NB), `IsAneka` (NB/RNW), `isBisnisOtherObjAneka`, `IsBonding`,
`IsBuilderRisk`, `IsChekCity`, `IsChekNation`, `IsChekProvince`, `IsChekRW`, `isClaimTreaty` (NB),
`IsCrime`, `IsDescDeductType`, `IsEdmAddObject`, `IsEdmAdjCeding` (EDM), `IsEdmAdjCurrency`,
`IsEdmAdjInsured`, `IsEdmAdjPeriod`, `IsEdmAdjRate`, `IsEdmAdjRefNo`, `IsEdmAdjRIC`,
`IsEdmAdjShareCedant`, `IsEdmAdjSpreading`, `IsEdmAdjTSI`, `IsEdmExtendPeriod`, `IsEdmPerubahan`,
`IsEdmPPNPPH`, `IsEngineering`, `IsEnvironmental`, `IsFinanceInsurance`, `IsFire`,
`IsFireStyle1` (EDM), `IsFireStyle2` (EDM), `IsHealthSSC` (EDM), `IsKPR` (EDM), `IsLimitCreditCL`,
`IsLimitCreditNCL`, `IsLimitCustomBond`, `IsLimitSBondKBG`, `IsLimitTradeCredit`, `IsMBU` (EDM),
`IsNonPropertyandNonEngineering`, `IsNotPAandNotMBU`, `IsObjectOtherScheduleVisible`,
`IsObjectParticipantVisible`, `IsObjectPolicyScheduleVisible`, `IsPASSG`, `IsPKSASM`, `IsPortRisk`,
`IsPropertyandEngineering`, `IsShowInput`, `IsSpreadingDepan` (EDM), `IsSPVCreate` (NB),
`IsSPVTreaty1` (NB), `IsT1T3`, `IsT2T4`, `IsTBonding`, `IsTJHInRange` (EDM),
`IsTJHNotInRange` (EDM), `isTravelTime`, `IsTreaty1` (NB), `IsUWSpread`, `pyIsIpadOrDesktop` (NB),
`pyIsMobile` (NB), `TreatyMasterInEDM` (NB).

Tanpa keterangan folder = hadir di ketiga siklus. Isi kondisi masing-masing ada di tabel §4.

Yang tercakup di sini mencakup **seluruh rangkaian `IsEdmAdj*`** — sepuluh predikat yang
membedakan jenis endorsement. Discovery yang hanya membaca `pyConditionString` akan melaporkan
mesin endorsement sebagai "tidak diketahui" padahal kondisinya terbaca seluruhnya.

---

## 7. Guard berbasis identitas orang — jumlah dan mekanisme saja

Sesuai CLAUDE.md §3.5, **tidak ada nilai identitas yang disalin**. Yang dicatat: properti yang
diuji, jumlah cabang, dan bentuk nilainya.

```powershell
$c | Where-Object { $_.PropNorm -match 'MarketingName|MarketingCode|pyUserIdentifier|pxCreateOperator|pyTelephone' } |
  Group-Object Rule,PropNorm | Sort-Object Name | ForEach-Object { "{0,3}  {1}" -f $_.Count, $_.Name }
```

### 7.1 Guard yang memuat identitas orang (6 rule)

| Rule | Siklus | Mekanisme | Jumlah cabang identitas |
| --- | --- | --- | --- |
| `IsGroup` | NB, RNW, EDM | `pyLogic = A OR B OR C OR D OR E`; A menguji `.OfferFacIn.QuotationData.IsGroup = "Group"`, B–D menguji `pxRequestor.OperatorID.pyUserIdentifier`, E menguji `pyWorkPage.OfferFacIn.QuotationData.MarketingCode` | **3 identitas operator + 1 kode kontak marketing** (berpola `xxx-xxxxxx-xxxx-xxxxxxx xxx-99`) |
| `IsGroupCreate` | NB, RNW, EDM | `A OR B` atas `pyWorkPage.pxCreateOperator` | **2 identitas operator** |
| `IsSPVCreate` | NB | `A OR B` atas `pyWorkPage.pxCreateOperator` | **2 identitas operator** |
| `IsSPVTreaty1` | NB | `A` atas `pxRequestor.OperatorID.pyTelephone` | **1 nomor telepon operator** |
| `IsTreaty1` | NB | `A` atas `pxRequestor.OperatorID.pyTelephone` | **1 nomor telepon operator** |
| `IsTBonding` | NB, RNW, EDM | `A OR B OR C OR D`; A menguji `.OfferFacIn.QuotationData.TeamGroup = 5`, B–D menguji `.OfferFacIn.QuotationData.MarketingName` | **3 nama marketing** |

**Semua enam identik di setiap siklus tempat mereka hadir** (tidak termasuk 39 rule bercabang)
**[terverifikasi]**.

**Risiko migrasi:** keenamnya mengeraskan identitas orang ke dalam logika akseptasi. `IsTBonding`
adalah yang paling tajam — ia menyamakan "tim bonding" dengan `TeamGroup = 5` **atau** tiga nama
orang tertentu; bila salah satu dari mereka berganti peran, perilakunya berubah diam-diam.
`IsSPVTreaty1` dan `IsTreaty1` memakai **nomor telepon** sebagai kunci identitas — nilai yang
dapat berubah tanpa kaitan dengan otorisasi.

> **[pertanyaan terbuka]** Keenam guard ini **kandidat perbaikan**, bukan kandidat "diperbaiki
> diam-diam saat migrasi". Bisnis harus memutuskan: dipertahankan apa adanya (paralel run cocok
> angka), atau dipindah ke atribut peran/workgroup (paralel run akan berbeda).

### 7.2 Guard berbasis konstruk organisasi (bukan identitas orang — nilainya boleh dikutip)

| Rule | Siklus | Kondisi |
| --- | --- | --- |
| `IsAdmin` | NB, RNW, EDM | `A OR B` atas `pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName` = `"ReasFacInAdmin"` ‖ `""` |
| `IsIT` | NB, RNW | `A OR B`: `OperatorID.pyPosition = "IT Developer"` ‖ `OperatorID.pyOrgDivision = "IT"` |
| `IsNotAdmin` | NB | `A`: `OperatorID.pyPosition != "Admin"` |
| `IsOperatorLife` | NB | `A`: `OperatorID.pyWorkGroup = "ReasLife"` |
| `IsUW` | NB, RNW, EDM | lihat §3.3 pola P3 — **bercabang** |
| `IsUWSpread` | NB, RNW, EDM | `A OR … OR F` atas `pyWorkPage.PositionNote` |

Perhatikan cabang `pyWorkBasketName = ""` pada `IsAdmin`: **operator tanpa workbasket dihitung
sebagai admin** **[terverifikasi]**. Itu perilaku fail-open. **[pertanyaan terbuka]** apakah
disengaja.

Perhatikan juga `pyWorkBasketList(1)` — **hanya elemen pertama** daftar workbasket yang diuji, pada
`IsAdmin` maupun `IsUW`. Operator yang memiliki beberapa workbasket akan dinilai hanya dari yang
pertama, dan urutannya tidak ditentukan di rule ini **[terverifikasi]**.

---

## 8. Properti yang paling sering diuji + enumerasi nilainya

Ekstraksi memecah setiap `pyConditionValue1String` menjadi `properti / operator / nilai`, lalu
menormalkan awalan halaman (`pyWorkPage.` → `.`) agar `pyWorkPage.Quotation.BusinessType` dan
`.Quotation.BusinessType` terhitung sebagai properti yang sama.

```powershell
$c | Group-Object PropNorm | Sort-Object Count -Descending | Select-Object -First 30 |
  ForEach-Object { "{0,5}  {1}" -f $_.Count, $_.Name }
```

Total baris kondisi terurai: **1.863**.

| Baris | Properti (ternormalisasi) |
| --- | --- |
| 417 | `.Quotation.BusinessCode` |
| 213 | `.OfferFacIn.QuotationData.BusinessOldId` |
| 143 | `.Quotation.BusinessType` |
| 48 | `.CoverageID` |
| 48 | `.Quotation.BusinessOldId` |
| 48 | `.Coverage` |
| 43 | `.OfferFacIn.QuotationData.Type` |
| 40 | `.OfferFacIn.QuotationData.StatusBusiness` |
| 37 | `.OfferFacIn.QuotationData.EdmType` |
| 37 | `OutData.pxResults(1).CARI2` |
| 35 | `.OfferFacIn.QuotationData.BusinessType` |
| 34 | `.LetterNo` |
| 29 | `.TabPosition` |
| 27 | `.DeductibleType` |
| 27 | `.OfferFacIn.QuotationData.TeamGroup` |
| 21 | `.Quotation.SobLeader1` |
| 18 | `.PositionNote` |
| 15 | `.Quotation.BusinessName` |
| 15 | `.OfferFacIn.IsB2B` |
| 12 | `.IsCedingConfirm` |
| 10 | `pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName` |
| 9 | `.Quotation.StatusBusiness` |
| 9 | `OperatorID.pyUserIdentifier` |
| 9 | `.PropertyItemGroup` |
| 9 | `.OfferFacIn.QuotationData.MarketingName` |
| 8 | `.pxCreateOperator` |

### 8.1 Enumerasi lengkap per properti kunci — [terverifikasi]

```powershell
foreach ($p in @('.OfferFacIn.QuotationData.Type', '.OfferFacIn.QuotationData.StatusBusiness', …)) {
  $v = $c | Where-Object { $_.PropNorm -eq $p }
  $v | Group-Object Op,Val | Sort-Object Count -Descending | ForEach-Object { "{0,4}  {1}" -f $_.Count, $_.Name }
}
```

**`QuotationData.StatusBusiness`** — hanya satu nilai yang pernah diuji lewat properti ini: `3`
(40 baris). Lewat `.Quotation.StatusBusiness` (9 baris) diuji `= 1`, `= 2`, `!= 3`.

| Rule | Kondisi |
| --- | --- |
| `IsNB` | `pyWorkPage.Quotation.StatusBusiness = 1` |
| `IsRenewal` | `pyWorkPage.Quotation.StatusBusiness = 2` |
| `IsEDM` | `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3` |
| `IsNotEDM` | `pyWorkPage.Quotation.StatusBusiness != 3` |

> **[dugaan]** `StatusBusiness` ∈ {1,2,3} adalah diskriminator siklus NB / RNW / EDM. Bukti
> langsungnya hanya nama rule, dan nama rule bukan bukti. Yang **[terverifikasi]** adalah nilai yang
> diuji dan properti yang dipakai.
> **[pertanyaan terbuka]** `IsEDM` membaca `.OfferFacIn.QuotationData.StatusBusiness` sementara
> `IsNB`/`IsRenewal`/`IsNotEDM` membaca `.Quotation.StatusBusiness` — **properti berbeda**. Keempat
> predikat itu tidak saling eksklusif secara struktural.

**`QuotationData.Type`** — 12 nilai, `0`–`12` kecuali `10`, dengan `12` diuji 7 kali dan `0` 6 kali:
`0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12`. Arti tiap kode: **belum terverifikasi**.

**`QuotationData.EdmType`** — hanya nilai `4` yang pernah diuji (37 baris). Nilai lain tidak muncul
di korpus `When\`. Arti `4`: **belum terverifikasi**.

**`QuotationData.TeamGroup`** — nilai yang diuji: `1, 2, 3, 4, 5` (`2` paling sering, 12 baris).

| Rule | Kondisi |
| --- | --- |
| `IsT1T4` | `TeamGroup = 1` OR `TeamGroup = 4` |
| `IsT2T3` | `TeamGroup = 2` OR `TeamGroup = 3` |
| `IsT1T3` | `A OR (B AND C)` = `TeamGroup = 1` OR (`TeamGroup = 2` AND `IsB2B = "ASM"`) |
| `IsT2T4` | `(A AND B) OR (C AND D)` = (`TeamGroup = 2` AND `IsSpecialAcceptance = "false"`) OR (`TeamGroup = 2` AND `IsB2B = "ASM"`) |
| `IsTBonding` | `TeamGroup = 5` OR 3 cabang identitas (§7.1) |

> **[pertanyaan terbuka]** Nama `IsT1T3` menyiratkan "TeamGroup 1 atau 3", tetapi isinya menguji
> `1` dan `2`. Nama `IsT2T4` menyiratkan "2 atau 4", isinya menguji `2` dua kali. Nama rule
> **bertentangan dengan isinya** pada kedua kasus. Mana yang benar — nama atau isi — hanya bisnis
> yang tahu. Isi adalah yang dieksekusi Pega.

**`.LetterNo`** — 13 nilai, sekaligus **enumerasi jabatan approver** [terverifikasi]:

| Nilai `LetterNo` | Rule penguji | Siklus |
| --- | --- | --- |
| `"UNDERWRITER"` | `ToUW` | NB, RNW, EDM |
| `"SENIORUW"` | `ToSeniorUW` | NB, RNW, EDM |
| `"JUW_A"` | `ToJUW_A` | NB, RNW, EDM |
| `"DEPHEADUNDERWRITER"` | `ToDepHeadUW` | NB, RNW, EDM |
| `"MANAGERTEKNIK"` | `ToManagerTeknik` | NB, RNW, EDM |
| `"KADIVTEKNIK"` | `ToKadivTeknik` | NB, RNW, EDM |
| `"KADIVFACULTATIVE"` | `ToKadivFacultative` | NB, RNW, EDM |
| `"KADIVFINANCIAL"` | `ToKadivFin` | NB, RNW, EDM |
| `"DIREKTURTEKNIK"` | `ToDirTeknik` | NB, RNW, EDM |
| `"DIREKTURMARKETING"` | `ToDirMarketing` | NB, RNW, EDM |
| `"DEPTHEADUWLIFE"` | `ToDeptHeadUWLife` | **NB saja** |
| `"TREATYINDEPTHEAD"` | `ToTREATYDEPTHEAD` | **NB saja** |
| `""` | `LetterNoNull` | NB, RNW |

`ToDepHeadUW` adalah satu-satunya rule `To*` dengan kondisi majemuk:
`A AND (B OR C)` = `LetterNo = "DEPHEADUNDERWRITER"` AND (`OfferFacIn.IsB2B != "ASM"` OR
`OfferFacIn.IsB2B = ""`) **[terverifikasi]**. Artinya jalur Dep-Head UW **dilewati untuk kasus
ber-`IsB2B = "ASM"`**, sementara sebelas jabatan lain tidak punya syarat tambahan.

> **[pertanyaan terbuka]** Tidak ada rule `To*` untuk `"DEPTHEADUWLIFE"` dan `"TREATYINDEPTHEAD"` di
> siklus RNW maupun EDM. Apakah kedua jabatan itu memang tidak pernah muncul di renewal/endorsement,
> atau routing-nya ditangani di luar rule `When`?

**`.PositionNote`** — 6 nilai, dipakai hanya oleh `IsUWSpread` [terverifikasi]:
`"ReasFacInUnderwriting"`, `"ReasFacInSeniorUnderwriting"`, `"ReasFacInUnderwritingFinancial"`,
`"ReasFacInJuniorUnderwriting"`, `"ReasFacInDepHeadUnderwriting"`, `"ReasFacInJuniorUnderwritingA"`.

**`pyWorkBasketName`** — 5 nilai: `"ReasFacInAdmin"`, `"ReasFacInDirector"`,
`"ReasFacInGroupLeader"`, `"ReasFacInUnderwriting"`, `"ReasFacInMarketing"` (yang terakhir **hanya
di varian EDM `IsUW`**, §3.3), plus `""` pada `IsAdmin`.

**`.OfferFacIn.IsB2B`** — 3 bentuk uji: `= "ASM"` (9 baris), `!= "ASM"` (3), `= ""` (3).
Dipakai oleh `IsPKSASM`, `IsT1T3`, `IsT2T4`, `ToDepHeadUW`.

**`.TabPosition`** — 10 nilai, satu per rule `IsTABActive_*` [terverifikasi]:
`"Cedant"`, `"Object"`, `"Coverage"`, `"Deduction"`, `"Spreading"`, `"SpreadingLife"`, `"Payment"`,
`"PaymentLife"`, `"ScoringRisk"`, `"LossRecord"`.

**`.DeductibleType`** — 7 kode numerik-string: `"100413"`, `"100414"`, `"100416"`, `"100417"`,
`"100418"`, `"100419"`, `"100420"`. Selain itu dua rule menguji kode **satu digit** lewat ekspresi
bebas:

```
IsDeducType      : @String.notEquals(.DeductibleType,"3")&&@String.notEquals(.DeductibleType,"1")
                   &&@String.notEquals(.DeductibleType,"2")&&@String.notEquals(.DeductibleType,"5")
IsTypeDeductType : @String.notEquals(.DeductibleType,"4")&&@String.notEquals(.DeductibleType,"7")
                   &&@String.notEquals(.DeductibleType,"8")&&@String.notEquals(.DeductibleType,"1")
                   &&@String.notEquals(.DeductibleType,"2")&&@String.notEquals(.DeductibleType,"5")
```

> **[pertanyaan terbuka]** Satu properti yang sama diuji terhadap dua ruang nilai yang tidak
> beririsan: kode 6-digit (`100413`…) dan kode 1-digit (`1`…`8`). Mana yang benar-benar tersimpan di
> `DeductibleType`? Salah satu dari dua kelompok rule ini **selalu** mengembalikan hasil yang sama —
> tetapi korpus tidak mengatakan yang mana.

**`.IsCedingConfirm`** — diuji terhadap `Offer`, `Policy` (tanpa tanda kutip di ekspor) dan
`"accepted"` (dengan tanda kutip) [terverifikasi]:

```
IsCedingConfirmOffer  [A OR B] : .IsCedingConfirm = Offer   ; .IsCedingConfirm != "accepted"
IsCedingConfirmPolicy [A OR B] : .IsCedingConfirm = Policy  ; .IsCedingConfirm = "accepted"
```

> **[pertanyaan terbuka]** `Offer` dan `Policy` tanpa kutip di ekspresi Pega adalah **referensi
> properti**, bukan literal string. Bila benar, kedua rule membandingkan `IsCedingConfirm` dengan
> isi properti `Offer`/`Policy`, bukan dengan teks `"Offer"`/`"Policy"`. Ini mengubah artinya
> sepenuhnya. Korpus `When\` tidak memuat definisi properti `Offer`/`Policy`.

**`BusinessType` / `CARI2`** — 39 nilai unik lintas ketiga siklus [terverifikasi]:

```
" ElectronicEquipment " · "AllRisk" · "Aneka" · "AviationHull" · "BillboardNeonSyariah" · "Boiler" ·
"Bonding" · "BondingKBG" · "Burglary" · "Car" · "CIS" · "CIT" · "ContractorsPM" · "CustomBond" ·
"Ear" · "ElectronicEquipment " · "Exclusion" · "Fidelity" · "Fire" · "FireGeneral" · "FireStyle1" ·
"FireStyle2" · "Glass" · "GolfInsurance" · "GrowingTrees" · "HE" · "KPR" · "LandRig" · "Liability" ·
"Maintenance" · "MarineCargo" · "MarineHull" · "MBD" · "MBUCar" · "MBUMotorCycle" · "OilGas" · "PA" ·
"Travel" · "Workmen"
```

> **Dua nilai bermasalah [terverifikasi]:** `" ElectronicEquipment "` (spasi **di depan dan
> belakang**) dipakai `isElectronicEquipment` di ketiga siklus; `"ElectronicEquipment "` (spasi
> belakang saja) dipakai varian EDM `IsAneka` lewat `CARI2`. Keduanya tidak akan pernah cocok dengan
> nilai yang sama. Ini bug tipografis di sistem lama — **kandidat perbaikan, bukan sesuatu yang
> boleh dibetulkan diam-diam saat migrasi** (CLAUDE.md §1). Di Go, literalnya harus disalin
> **berikut spasinya**, dengan komentar yang menunjuk berkas asalnya.

**`.Quotation.BusinessCode`** — 424 baris kondisi, **98 nilai kode unik**. Pemakai terbesar:
`IsHealth` (84 baris), `IsHealthIndividu` (75), `IsObjectOtherScheduleVisible` (72),
`IsObjectSectionAneka` (36), `isInterestInsuredOtherObjAneka` (33). Daftar nilai lengkapnya ada di
kolom kondisi tabel §4 untuk masing-masing rule.

**`.OfferFacIn.QuotationData.BusinessOldId`** — 261 baris, **87 nilai unik**. Pemakai terbesar:
`IsLimitSBondKBG` (123 baris), `IsLife` (48), `IsLimitCreditCL` (33), `IsLimitCustomBond` (30).

> **[pertanyaan terbuka]** `BusinessCode` (98 nilai) dan `BusinessOldId` (87 nilai) adalah dua
> kamus kode produk yang berdampingan, keduanya tanpa tabel referensi di korpus `When\`. Tanpa
> kamusnya, 80 rule bertema lini bisnis tidak dapat diverifikasi — hanya disalin.

### 8.2 Bentuk ekspresi selain `compareTwoValues`

```powershell
$d | Group-Object Purposes | Sort-Object Count -Descending | ForEach-Object { "{0,4}  {1}" -f $_.Count, $_.Name }
```

| Berkas | `pyConditionValue1Purpose` |
| --- | --- |
| 533 | `CompareTwoValues` |
| 28 | `CompareTwoValues` + `Rule-Obj-When` |
| 22 | `Rule-Obj-When` |
| 10 | `1FreeFormExpressionBoolean` |
| 3 | `1FreeFormExpressionBoolean` + `CompareTwoValues` |
| 2 | `StringEqualsIgnoreCase` |
| 1 | `PropertyHasValue` |
| 1 | `CompareTwoValues` + `StringEqualsIgnoreCase` |
| 1 | `compareTwoStrings` |

Ekspresi bebas (`1FreeFormExpressionBoolean`) yang perlu penanganan khusus di Go
**[terverifikasi]**:

| Rule | Ekspresi |
| --- | --- |
| `hasPrimaryPage` | `@hasPrimaryPage(tools)` |
| `pyIsMobile` | `@(Pega-RulesEngine:Utilities).pzIsMobile(tools)` |
| `IsNotActive` | `@DateTime.CompareDates(pyWorkPage.Quotation.EdmDate, .EndDate)` |
| `IsDeducType`, `IsTypeDeductType` | rantai `@String.notEquals(...)&&...` (lihat di atas) |
| `IsShowInput` | `A OR B`: `pyPortal.IsEdit` ‖ `compareTwoValues(OutputData.HASIL1, "!=", "")` |
| `isSellingModeB2B` | `compareTwoStrings(@(Pega-RULES:Utilities).getDataSystemSetting("PegaCRM-","SellingMode"), "EQUALS", "B2B")` |
| `isSellingModeB2BB2C` | `equalsIgnoreCase(getDataSystemSetting("PegaCRM-","SellingMode"), "B2B_B2C")` |
| `isSellingModeB2C` | `equalsIgnoreCase(getDataSystemSetting("PegaCRM-","SellingMode"), "B2C")` |
| `crmCreateOpportunity` | `(A0 Or A1) And A2` atas `evaluateWhen("crmBypassOperatorAccessChecks")`, `Declare_crmOperatorAccess.canCreateOpportunity`, `evaluateWhen("crmIsOpen")` |

> **[pertanyaan terbuka]** `IsNotActive` mengembalikan hasil `@DateTime.CompareDates(a,b)`
> **langsung sebagai boolean**. `CompareDates` pada umumnya mengembalikan −1/0/1. Bila begitu,
> predikatnya bernilai true untuk nilai apa pun yang bukan nol/false — perilakunya bergantung pada
> aturan koersi Pega yang tidak ada di korpus.

> **[pertanyaan terbuka]** `isSellingModeB2B` memakai `compareTwoStrings(..., "EQUALS", ...)`
> sementara `isSellingModeB2C` memakai `equalsIgnoreCase(...)`. Satu peka huruf besar-kecil, satu
> tidak, untuk keputusan yang setara. Korpus tidak mengatakan apakah itu disengaja.

> **[pertanyaan terbuka]** `getDataSystemSetting("PegaCRM-", "SellingMode")` membaca konfigurasi
> sistem yang **isinya tidak ada di korpus**. Ketiga rule `isSellingMode*` tidak dapat dievaluasi
> tanpa nilai itu. Sesuai CLAUDE.md §4.4, di target ia menjadi env var / konfigurasi.

---

## 9. Ringkasan angka — semua dengan perintah auditnya

| # | Ukuran | Nilai | Perintah |
| --- | --- | --- | --- |
| 1 | Berkas `When` — NB / RNW / EDM | 210 / 189 / 202 | §1 |
| 2 | Total berkas `When` | 601 | §1 |
| 3 | Nama rule unik (case-insensitive) | 226 | §1 |
| 4 | Nama rule unik (case-sensitive) | 237 | §1 |
| 5 | Nama dieja beda kapitalisasi antar siklus | 11 (semuanya EDM vs NB+RNW) | §1 |
| 6 | Nama hadir di 3 folder | 183 | §2 |
| 7 | Nama hadir di 2 folder | 9 | §2 |
| 8 | Nama hadir di 1 folder | 34 (NB 18, EDM 16, **RNW 0**) | §2 |
| 9 | Nama hadir di ≥2 folder | 192 | §3.1 |
| 10 | SHA-256 berkas identik antar folder | 6 | §3.1 |
| 11 | `pyLogic` + ekspresi identik antar folder | 153 | §3.2 |
| 12 | **Rule bercabang antar siklus** | **39** | §3.2 |
| 13 | Ekspresi sama, `pyClassName` beda | 10 | §3.4 |
| 14 | Berkas dengan kondisi **tidak terbaca sama sekali** | **0** | §5 |
| 15 | Berkas dengan `pyConditionString` tersembunyi | 170 (28,3 %) → 71 nama rule | §6 |
| 16 | — di antaranya `<pyConditionString/>` kosong | 13 berkas / 7 nama | §6.1 |
| 17 | — di antaranya placeholder | 157 berkas / 64 nama | §6.2 |
| 18 | Berkas dengan `pyTempText = true` | 294 (48,9 %) — **bukan penanda** | §5.3 |
| 19 | Berkas dengan label `pyConditionValue1String` belum resolve | 4 | §5.2 |
| 20 | Total baris kondisi terurai | 1.863 | §8 |
| 21 | Baris `Rule X evaluates to true` | 245, merujuk 63 rule | §4.2 |
| 22 | Dependensi rule yang hilang dari korpus | 1 (`pyIsIPad`) | §4.2 |
| 23 | Baris kondisi yatim / mati | 0 | §4.3 |
| 24 | `pyRuleAvailable` = `Yes` | 601 / 601 | perintah di bawah |
| 25 | Rule dengan guard identitas orang | 6 | §7.1 |
| 26 | Nilai `BusinessType`/`CARI2` unik | 39 | §8.1 |
| 27 | Nilai `BusinessCode` unik | 98 | §8.1 |
| 28 | Nilai `BusinessOldId` unik | 87 | §8.1 |

```powershell
# #24 — tidak ada rule yang Blocked / Withdrawn
Select-String -Path "D:\migrasi\RNM\*\When\*.xml" -Pattern '<pyRuleAvailable>([^<]*)</pyRuleAvailable>' |
  ForEach-Object { $_.Matches[0].Groups[1].Value } | Group-Object | Select-Object Name,Count
```

### 9.1 Skrip ekstraksi yang memasok tabel `$d` dan `$c`

Seluruh angka bertanda `$d` / `$c` berasal dari dua CSV yang dibangun sekali di awal:

```powershell
# $d : satu baris per BERKAS When (601)
#   kolom: Folder, FileName, Path, Hash(SHA256), RuleName(<pyRuleName>), ClassName(<pyClassName>),
#          Logic(<pyLogic>), Available(<pyRuleAvailable>),
#          CondStrings  = seluruh //pyConditionViewer//pyConditionString    digabung " ||| "
#          V1Strings    = seluruh /pagedata/pyCondition/rowdata/pyConditionValue1String  digabung " ||| "
#          V1Exprs      = seluruh /pagedata/pyCondition/rowdata/pyConditionValue1        digabung " ||| "
#          CondLabels   = seluruh pyConditionLabel, Purposes = pyConditionValue1Purpose unik
# dibangun dengan: $x = [xml](Get-Content -LiteralPath $f -Raw -Encoding UTF8)
#                  $x.SelectNodes("//pyConditionViewer//pyConditionString")
#                  $x.SelectNodes("/pagedata/pyCondition/rowdata")

# $c : satu baris per BARIS KONDISI (1.863) — hasil pemecahan V1Strings
#   Prop / Op / Val, plus PropNorm = Prop dengan awalan 'pyWorkPage.' diganti '.'
```

Nol galat parse pada 601 berkas (`PARSE ERRORS: 0`).

---

## 10. Pertanyaan terbuka yang MEMBLOKIR implementasi

Urut berdasarkan besar dampaknya terhadap angka rekonsiliasi paralel run.

### B1 — Isi query di balik `OutData.pxResults(1).CARI2` (EDM)
Enam rule klasifikasi lini bisnis di siklus endorsement (`IsAneka`, `IsFire`, `IsGolfInsurance`,
`IsMarineCargo`, `IsMBU`, `IsPA`) membaca hasil query yang **tidak ada di korpus `When\`**.
Tanpa definisi query-nya, klasifikasi lini bisnis untuk endorsement tidak dapat diimplementasikan.
`pxResults(1)` juga hanya membaca baris pertama tanpa jaminan urutan. **(§3.3 pola P2)**

### B2 — Kamus `BusinessCode` (98 kode) dan `BusinessOldId` (87 kode)
Tidak ada tabel referensi di korpus. 80 rule bertema lini bisnis — sepertiga dari seluruh katalog —
hanya bisa disalin, tidak bisa diverifikasi. Bila satu kode saja salah salin, seluruh jalur
akseptasi untuk produk itu meleset. **(§8.1)**

### B3 — `IsUW` bercabang: workbasket `ReasFacInMarketing` diakui di EDM, tidak di NB/RNW
Bukan bug tampilan — ini menentukan siapa yang dianggap underwriter. Bisnis harus memutuskan
apakah percabangan ini disengaja sebelum satu baris Go ditulis. **(§3.3 pola P3)**

### B4 — `.Quotation.BusinessType` vs `.OfferFacIn.QuotationData.BusinessType`
30 rule membaca properti yang berbeda tergantung siklus. Apakah keduanya selalu sinkron? Bila
tidak, klasifikasi produk berbeda antara NB dan EDM untuk kasus yang sama. **(§3.3 pola P1)**

### B5 — `DeductibleType`: dua ruang nilai yang tidak beririsan
7 kode 6-digit (`100413`…`100420`) vs 6 kode 1-digit (`1`…`8`) diuji terhadap properti yang sama.
Satu kelompok rule pasti selalu mengembalikan hasil konstan. **(§8.1)**

### B6 — `IsCedingConfirm` dibandingkan dengan `Offer` / `Policy` **tanpa tanda kutip**
Bila itu referensi properti dan bukan literal, arti kedua rule berubah total. **(§8.1)**

### B7 — `IsT1T3` dan `IsT2T4`: nama rule bertentangan dengan isi
`IsT1T3` menguji TeamGroup 1 dan 2. `IsT2T4` menguji TeamGroup 2 dua kali (tidak pernah 4).
Nama atau isi — mana yang mencerminkan maksud bisnis? **(§8.1)**

### B8 — `" ElectronicEquipment "` dengan spasi di depan dan belakang
Dua rule membandingkan properti yang sama dengan dua literal yang berbeda hanya pada spasi.
Kandidat perbaikan yang **tidak boleh** dibetulkan diam-diam. **(§8.1)**

### B9 — `StepStatusFail`: arah predikat `inString(...) = 0`
Apakah `0` berarti "ditemukan di indeks 0" atau "tidak ditemukan"? Predikatnya terbalik bergantung
jawabannya, dan rule ini menggerbangi penanganan kegagalan langkah. **(§5.2)**

### B10 — Resolusi rule Pega: case-sensitive atau tidak?
11 nama rule punya dua ejaan. Bila resolusi case-sensitive, itu 22 predikat berbeda, bukan 11.
**(§1)**

### B11 — `pyWorkBasketList(1)` — hanya workbasket pertama yang diuji
`IsAdmin` dan `IsUW` mengabaikan workbasket kedua dan seterusnya, dengan urutan yang tidak
ditentukan. Operator multi-workbasket berperilaku tidak deterministik. **(§7.2)**

### B12 — `IsAdmin` fail-open: `pyWorkBasketName = ""` dihitung sebagai admin
Operator tanpa workbasket lolos sebagai admin. Disengaja atau tidak — ini keputusan authz, dan
CLAUDE.md §6 mensyaratkan persetujuan manusia sebelum perubahan authn/authz. **(§7.2)**

### B13 — Enam guard identitas orang yang dikeraskan ke dalam logika
`IsGroup`, `IsGroupCreate`, `IsSPVCreate`, `IsSPVTreaty1`, `IsTreaty1`, `IsTBonding`. Dipertahankan
apa adanya (paralel run cocok) atau dipindah ke atribut peran (paralel run berbeda)? **(§7.1)**

### B14 — `getDataSystemSetting("PegaCRM-", "SellingMode")`
Nilai konfigurasi tidak ada di korpus. Tiga rule `isSellingMode*` tidak dapat dievaluasi. **(§8.2)**

### B15 — Arti kode `QuotationData.Type` (0–9, 11, 12) dan `EdmType` (hanya `4` yang diuji)
Tidak ada satu pun rule yang menjelaskan artinya. Bahwa `EdmType` hanya pernah diuji `= 4` pada 37
baris menunjukkan kemungkinan ada nilai lain yang tidak pernah diperiksa — atau bahwa `4` adalah
satu-satunya nilai yang pernah ada. Korpus tidak memisahkan keduanya. **(§8.1)**

### B16 — `IsNotActive` memakai hasil `CompareDates` langsung sebagai boolean
Perilakunya bergantung aturan koersi Pega yang tidak ada di korpus. **(§8.2)**

---

*Dokumen ini hanya mencakup subfolder `When\`. Rule `Flow`, `Activity`, `DataTransform`,
`DecisionTable`, `DecisionTree`, `RDBList`, dan `ReportDefinition` di ketiga folder korpus belum
dibaca untuk dokumen ini; klaim mana pun tentang **kapan** sebuah predikat dipanggil harus berasal
dari dokumen yang membaca subfolder-subfolder itu.*
