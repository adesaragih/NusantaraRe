# Discovery Endorsement — E-3: before-image, `SetOldData`, dan perhitungan SELISIH

> **Sumber:** `D:\migrasi\RNM\Endorsment Fac In\` (banding ke `NB FacIn\`). READ-ONLY.
> Korpus Treaty (`RNM_BRD\`) **tidak dibaca** (K-005). Label mengikuti `CLAUDE.md` §3.
>
> **Inti endorsement:** nilai dasar akseptasi endorsement adalah **SELISIH**, bukan nilai penuh.
> Tanpa lapis before-image, selisih tidak dapat dihitung — karena itu ketiga lapisnya dipetakan lebih
> dulu, baru rumusnya.
>
> ⚠️ **Perubahan metode dari E-2.** Seluruh penelusuran cabang di dokumen ini memakai **pohon langkah
> XML** (`pySteps/rowdata` bersarang + `pyStepsPreCondParams/rowdata/pyStepsPreCondParamsWhen`),
> **bukan kedekatan nomor baris**. E-2 memakai rentang nomor baris untuk menebak keanggotaan blok —
> itu tidak memadai, dan diganti. Akibatnya lihat §4.1.
>
> ⚠️ **Privasi (E-Q15).** Tag `pyStepsRepeatDef` dan `pyStepsPreCondParams` memuat nama orang dan
> alamat email. Nilainya **tidak dibaca dan tidak disalin**; hanya sub-tag
> `pyStepsPreCondParamsWhen` (ekspresi kondisi, tanpa PII) yang dikutip.

---

## 1. Peta tiga lapis before-image

⚠️ **Istilah.** `[terverifikasi]` (E-2) token `BeforeImage`/`AfterImage` **tidak ada di korpus EDM**.
"Before-image" adalah istilah konsep, bukan nama properti.

| Lapis | Apa | Diisi oleh | Kapan | Status verifikasi |
| --- | --- | --- | --- | --- |
| **A** | Dokumen polis versi terakhir, utuh → `OfferFacIn.OldData` | `SetValueToEDMWork` langkah 14 | sekali, saat case EDM lahir | ✅ §2 |
| **B** | Angka lama **per baris** → properti bersaudara `*Old` | `SetOldData` | **setiap kali layar dibuka** | ✅ §3 |
| **C** | Penanda baris warisan → `.IsOldData = "old"` | 7× `SetOLDValueToEDMWork_<LOB>` | sekali, saat case lahir | ✅ §4 |

**Ketiganya berbeda siklus hidup.** A dan C sekali seumur case; **B diulang tiap buka layar** — dan
itulah yang membuat temuan §3.4 berdampak.

---

## 2. Lapis A — `OfferFacIn.OldData` ✅

`[terverifikasi]` `Endorsment Fac In\Activity\SetValueToEDMWork.xml` langkah **14**
(`pyStepsDescription` = *"Copy policy data to old data"*), **13 sub-langkah**.

⛔ `pyStepsPreCondition` langkah 14 = **`false`** — lihat §6.

| Sub | Bukti |
| --- | --- |
| **14.1** | `RDB-List` · `<RequestType>` = `GetEDMOldData_SQL` · `<ClassName>` = `ASM-FW-GISFW-Int-OFFERJSON` · `<BrowsePage>` = `OldData` · `pyStepsPreCondition` = **`false`** |
| 14.2 | `Java` — *"copy NB and remove the comments, and fac retro list"* |
| **14.3** | `Property-Set` — **tepat 54** `<PropertiesName>`, seluruhnya bertarget `newWorkPage.OfferFacIn.*` |
| 14.4–14.5 | `RDB-List` + `Property-Set` — *"-- Set team Group"* |
| 14.6 | `Obj-Save` |
| **14.7–14.13** | **7× `Call SetOLDValueToEDMWork_<LOB>`** — `_FIRE` `_MC` `_Aneka` `_MBU` `_LIFE` `_PA` `_GOLF`, semua `pyStepsPreCondition` = `true` |

### Query lapis A — dikutip apa adanya

`[terverifikasi]` `Endorsment Fac In\RDBList\GetEDMOldData_SQL.xml`:

```
<pxObjClass>     Rule-Connect-SQL            ← BUKAN Rule-Obj-RDBList
<pyClassName>    ASM-FW-GISFW-Int-OFFERJSON
<pyRequestType>  GetEDMOldData_SQL
```

```sql
SELECT a.DATA_JSON AS HASIL1 FROM JSON_POLIS a
 WHERE NOPOLIS={newWorkPage.OfferFacIn.QuotationData.OldPolicyNo}
   and PRODKE=(SELECT COUNT(NOPOLIS)-1 FROM JSON_POLIS a
                WHERE NOPOLIS={newWorkPage.OfferFacIn.QuotationData.OldPolicyNo})
```

📌 **Jebakan folder (ketujuh) berulang.** Berkas duduk di `RDBList\` tetapi `pxObjClass`-nya
`Rule-Connect-SQL` — sama seperti `InsertTreatyProduction_Sql` (E-2 §1.5). Ini **bukan kebetulan
satu berkas**; klasifikasi wajib dari `pxObjClass`.

⛔ **`PRODKE = COUNT(NOPOLIS)-1` benar hanya bila `PRODKE` rapat dari 0 tanpa lubang.** Bila satu
baris pernah dihapus, before-image yang diambil **versi yang salah** — dan seluruh SELISIH ikut
salah. Tetap **#11, MEMBLOKIR**, milik DBA.

```powershell
$x=[xml](Get-Content "D:\migrasi\RNM\Endorsment Fac In\RDBList\GetEDMOldData_SQL.xml" -Raw)
$x.SelectSingleNode('//pxObjClass').InnerText; $x.SelectSingleNode('//pyBrowseSQL').InnerText
```

---

## 3. Lapis B — `SetOldData` dibaca tuntas ✅

`[terverifikasi]` `Endorsment Fac In\Activity\SetOldData.xml` · `pxObjClass` = `Rule-Obj-Activity` ·
`pyClassName` = `ASM-FW-GISFW-Work` · **4 langkah tingkat atas**. Dipanggil
`InputAddendumFacIn_PreAct` langkah **28** (E-2 §1.3).

### 3.1 Tiga gerbang keluar di depan — dengan kode transisinya

`[terverifikasi]` (kode transisi Pega: `2` = lanjut langkah berikut, `6` = keluar dari activity)

| Langkah | `pyStepsPreCondParamsWhen` | `WhenTrue` | `WhenFalse` | `pyStepsPreCondition` | Arti |
| ---: | --- | ---: | ---: | --- | --- |
| **1** | `IsLife` | **6** | 2 | `true` | Life → **keluar**; Life punya jalur sendiri |
| **2** | `IsEDM` | 2 | **6** | `true` | bukan EDM → **keluar** |
| **3** | `@Utilities.SizeOfPropertyList(.OfferFacIn.LocationList)>100` | **6** | 2 | `true` | >100 lokasi → **keluar** |
| 4 | `IsFire` | 2 | 3 | ⛔ **`false`** | lihat §3.5 |

⛔ **Batas 100 lokasi mengikat:** polis dengan lebih dari 100 lokasi **tidak mendapat lapis B sama
sekali** — seluruh `*Old` kosong, sehingga SELISIH per baris tidak dapat dihitung. Ambangnya
`[pertanyaan terbuka]` (arsip #18).

📌 Gerbang 2 memakai **`IsEDM`**, dan `IsEDM` membaca `OfferFacIn.QuotationData.StatusBusiness`
(E-Q13). Jadi lapis B bergantung pada **agregat tersimpan**, bukan halaman aktif.

### 3.2 Blok 4 — 14 sub-langkah, tujuh lini bisnis berpasangan

`[terverifikasi]` Seluruh 14 sub-langkah `pyStepsPreCondition` = **`true`** (aktif). Polanya
berpasangan: satu untuk objek, satu untuk cedant.

| Sub | `pyStepsDescription` | `pyStepsPreCondParamsWhen` | Indeks lokal |
| ---: | --- | --- | --- |
| 4.1 / 4.2 | *"Untuk FIRE"* / *"Untuk FIRE Cedant"* | `IsFire` | `idxlocation` / `idxcedant` |
| 4.3 / 4.4 | *"Untuk Golf Insurance"* / *"Untuk Golf Ceding"* | `IsGolfInsurance` | `idxlocation` / `idxcedant` |
| 4.5 / 4.6 | *"Untuk Aneka"* / *"Untuk Aneka Ceding"* | `IsAneka` | `idxlocation` / `idxcedant` |
| 4.7 / 4.8 | *"Untuk PA"* / *"Untuk PA Ceding"* | `IsPA` | `idxperson` / `idxcedant` |
| 4.9 / 4.10 | *"Untuk Marine Cargo"* / *"…Ceding"* | `IsMarineCargo` | `idxcargo` / `idxcedant` |
| 4.11 / 4.12 | *"Untuk MBU"* / *"Untuk MBU Ceding"* | `IsMBU` | `idxvehicle` / `idxcedant` |
| 4.13 / 4.14 | *"Untuk Travel"* / *"Untuk Travel Ceding"* | `IsTravel` | `idxtravel` / `idxcedant` |

⚠️ **Tujuh lini di `SetOldData`, tetapi bukan tujuh yang sama dengan lapis C.** `SetOldData` menangani
**Travel** dan **tidak** menangani Life (keluar di gerbang 1). Lapis C menangani **Life** dan
**tidak** menangani Travel. Himpunannya berbeda — `[pertanyaan terbuka]` **E-Q22**.

### 3.3 Properti `*Old` yang diisi per lini `[terverifikasi]`

| Lini | Properti `*Old` (jalur diringkas dari akar `pyWorkPage.OfferFacIn`) |
| --- | --- |
| **FIRE** (4.1) | `LocationList(i).Property.PropertyItemList(j).{TSIObjectItemOld, TotalGrossPremiOld, TotalPremiumNusantaraReOld}` · `…Property.TotalTSIList(j).TSIOld` · `…Property.TotalTSIPremiGrossList(k).{TSIOld, PremiumOld, RateOld}` — **7** |
| **GOLF** (4.3) | `LocationList(i).Property.RiskLocation.AnekaList(j).TSIOld` · `…AnekaList(j).CoverageList(k).{TSIOld, PremiumOld}` · `…TotalTSIPremiGrossList(j).{TSIOld, PremiumOld, RateOld}` — **6** |
| **ANEKA** (4.5) | `LocationList(i).Property.RiskLocation.OccupationList(j).AnekaList(k).TSIOld` · `…TotalTSIPremiGrossList(j).{TSIOld, PremiumOld, RateOld}` — **4** |
| **PA** (4.7) | `PersonList(i).ASMCoverage(j).{TSIOld, PremiumOld}` · `PersonList(i).TotalTSIPremiGrossList(j).{TSIOld, PremiumOld, RateOld}` — **5** |
| **MARINE CARGO** (4.9) | `CargoList(i).CoverageList(j).{TSIOld, PremiumOld, PremiNusantaraReOld}` · `CargoList(i).TotalTSIPremiGrossList(j).{TSIOld, PremiumOld, RateOld}` — **6** |
| **MBU** (4.11) | `VehicleList(i).CoverageList(j).{TSIOld, PremiumOld, PremiNusantaraReOld, PremiRpOld, PremiumGrossDiscountFleetOld}` · `VehicleList(i).TotalTSIPremiGrossList(j).{TSIOld, PremiumOld, RateOld}` — **8** |
| **TRAVEL** (4.13) | `PersonList(i).ASMCoverage(j).{TSIOld, PremiumOld}` — **2** |
| **Cedant** (4.2/4/6/8/10/12/14) | `CedingCedantList(i).CurrencyList(j).{TSIOld, PremiumOld}` — **2**, identik di ketujuh lini |

`[terverifikasi]` **52 penugasan `*Old`** seluruhnya, dari **89 penugasan** total di `SetOldData`.

### 3.4 Pola nilai — dan ⛔ tujuh guard yang menguji dirinya sendiri

`[terverifikasi]` **52 dari 52** memakai pola `@If(...)`. **Nol pengecualian.** Bentuk bakunya:

```
TSIObjectItemOld           = @If(.TSIObjectItem!="",           .TSIObjectItem,           0)
TSIOld                     = @If(.TSI!="",                     .TSI,                     0)
TotalGrossPremiOld         = @If(.TotalGrossPremi!="",         .TotalGrossPremi,         0)
TotalPremiumNusantaraReOld = @If(.TotalPremiumNusantaraRe!="", .TotalPremiumNusantaraRe, 0)
```

✅ **Kosong dipetakan ke `0`, bukan NULL** — dikonfirmasi ulang, 52/52.

⛔ **Tetapi 7 dari 52 mengujinya terbalik** `[terverifikasi]` — guard membaca **properti tujuan**,
bukan properti sumber:

| Properti daun | Jumlah | Ekspresi |
| --- | ---: | --- |
| **`RateOld`** | **6 dari 6** — *selalu* | `@If(.RateOld!="", .Rate, 0)` ← menguji `.RateOld` |
| **`PremiumOld`** | **1 dari 18** — *tidak konsisten* | `@If(.PremiumOld!="", .Premium, 0)` |
| `TSIOld` | 0 dari 21 | `@If(.TSI!="", .TSI, 0)` — benar |
| sisanya | 0 | benar |

**Akibat yang terbaca dari ekspresi** `[terverifikasi]`: pada eksekusi ketika `.RateOld` masih kosong,
guard bernilai salah → `RateOld` diisi **`0`**, bukan `.Rate`.

`[dugaan]` Karena `SetOldData` dipanggil **setiap kali layar dibuka** (§1), pembukaan berikutnya
menemukan `.RateOld` sudah terisi sehingga menyalin `.Rate` **saat itu** — yaitu rate **baru**, bukan
rate polis lama. Semantik evaluasi persisnya tidak dapat dipastikan dari korpus → **E-Q18**.

⛔ **Diport apa adanya** (`CLAUDE.md` §1). Ini kandidat perbaikan milik bisnis — memperbaikinya
diam-diam akan mengubah angka SELISIH rate dan memutus rekonsiliasi paralel run.

#### ⛔ E-Q19 — fallback `"0"` berkutip · **TERBUKTI, dengan koreksi angka**

> 🔁 **Koreksi.** Versi pertama menulis *"`PremiNusantaraReOld` memakai `"0"`, **51 lainnya** `0`
> numerik"*. Angkanya **salah**: yang berkutip ada **2**, bukan 1, sehingga pembandingnya **50**,
> bukan 51. Substansi klaim (ketidakkonsistenan tipe) **bertahan**; angkanya diperbaiki.

`[terverifikasi]` `Endorsment Fac In\Activity\SetOldData.xml`:

| Bentuk fallback | Jumlah |
| --- | ---: |
| `0` — **numerik** | **50** |
| `"0"` — **berkutip (string)** | **2** |
| tidak terpola | 0 |
| **total penugasan `*Old`** | **52** |

**Kedua yang berkutip — dikutip persis beserta nomor barisnya:**

```
L6293  <PropertiesName>pyWorkPage.OfferFacIn.CargoList(local.idxcargo).CoverageList(local.idx).PremiNusantaraReOld</PropertiesName>
L6294  <PropertiesValue>@If(.PremiNusantaraRe!="",.PremiNusantaraRe,"0")</PropertiesValue>

L7526  <PropertiesName>pyWorkPage.OfferFacIn.VehicleList(Local.idxvehicle).CoverageList(Local.idx).PremiNusantaraReOld</PropertiesName>
L7527  <PropertiesValue>@If(.PremiNusantaraRe!="",.PremiNusantaraRe,"0")</PropertiesValue>
```

⚠️ **Mengapa nilainya tidak terlihat saat memeriksa L6293/L7526 saja:** `<PropertiesValue>` adalah
**tag saudara berikutnya**, bukan bagian dari baris `<PropertiesName>`. Struktur `rowdata`
induknya (L6289 dan L7522) berurutan: `pxObjClass` → `pxCreateOperator` → `pxCreateDateTime` →
**`PropertiesName`** → **`PropertiesValue`** → `pxCreateSystemID` → `pxCreateOpName` →
`pyExpressionGadget`. Pemeriksaan wajib mengambil **pasangan dalam satu `rowdata`**, bukan satu baris.

**Tiga pembanding fallback numerik `[terverifikasi]`:**

```
L895  → L896   TSIObjectItemOld            = @If(.TSIObjectItem!="",.TSIObjectItem,0)
L961  → L962   TotalGrossPremiOld          = @If(.TotalGrossPremi!="",.TotalGrossPremi,0)
L987  → L988   TotalPremiumNusantaraReOld  = @If(.TotalPremiumNusantaraRe!="",.TotalPremiumNusantaraRe,0)
```

📌 **Temuan tambahan — cache editor tidak sepakat dengan ekspresi tersimpan** `[terverifikasi]`.
Pada `rowdata` yang sama, `<pyExpressionGadget>` memuat bentuk **tanpa kutip**:

| Tag | Baris | Isi |
| --- | ---: | --- |
| `<PropertiesValue>` (dieksekusi) | L6294 / L7527 | `@If(.PremiNusantaraRe!="",.PremiNusantaraRe,` **`"0"`** `)` |
| `<pyExpressionGadget>` (cache form) | L6297 / L7530 | `@If(.PremiNusantaraRe!="",.PremiNusantaraRe,` **`0`** `)` |

⛔ Yang mengikat adalah **`<PropertiesValue>`**; `pyExpressionGadget` adalah cache tampilan editor.
Penegasan `CLAUDE.md` §3.3 dalam bentuk baru: **bukan hanya nama yang bukan bukti — cache render pun
bukan bukti.**

⛔ **Diport apa adanya.** Sistem lama sudah terbukti membandingkan angka sebagai string di tempat lain
(`CLAUDE.md` §4.1); menyeragamkan `"0"` → `0` diam-diam adalah perbaikan terlarang (§1).

```powershell
# menghasilkan 50 numerik / 2 berkutip / 52 total
Add-Type -AssemblyName System.Xml.Linq
$XN=[System.Xml.Linq.XName]
$xd=[System.Xml.Linq.XDocument]::Load("D:\migrasi\RNM\Endorsment Fac In\Activity\SetOldData.xml",
      [System.Xml.Linq.LoadOptions]::SetLineInfo)
$xd.Descendants($XN::Get('PropertiesName')) | ? { $_.Value.Trim() -match 'Old$' } | % {
  $pv=$_.Parent.Element($XN::Get('PropertiesValue')); if($pv){
    $m=[regex]::Match($pv.Value.Trim(),',\s*("0"|0)\s*\)\s*$')
    "{0}`t{1}" -f $(if($m.Success){$m.Groups[1].Value}else{'(lain)'}),
                  ([System.Xml.IXmlLineInfo]$pv).LineNumber } } |
  Group-Object { ($_ -split "`t")[0] } | Select-Object Count,Name
```

```powershell
$d=[xml](Get-Content "D:\migrasi\RNM\Endorsment Fac In\Activity\SetOldData.xml" -Raw)
$d.SelectNodes('//pyParamArray/rowdata') | % {
  $n=$_.SelectSingleNode('PropertiesName'); $v=$_.SelectSingleNode('PropertiesValue')
  if($n -and $v -and $n.InnerText -match 'Old$'){ "$($n.InnerText.Trim()) = $($v.InnerText.Trim())" } } |
  ? { $_ -match '@If\(\.(\w+Old)!=' }      # => 7 baris: RateOld ×6, PremiumOld ×1
```

### 3.5 ⛔ Langkah 4 bergerbang `IsFire` — **dengan prakondisi dinonaktifkan**

`[terverifikasi]`

```
langkah 4 : pyStepsPreCondParamsWhen = IsFire
            pyStepsPreCondition      = false      ← DINONAKTIFKAN
```

Ini **titik paling menentukan di seluruh lapis B**, dan jawabannya bukan milik korpus:

| Bila `false` berarti… | Akibat pada `SetOldData` |
| --- | --- |
| **prakondisi tidak dievaluasi, langkah tetap jalan** | Blok 4 jalan untuk semua lini; 4.1–4.14 menyaring sendiri per lini → **perilaku wajar** |
| **langkah dinonaktifkan** | Blok 4 **tidak pernah jalan** → **lapis B tidak pernah terisi untuk lini mana pun** |

⛔ Kedua bacaan berlawanan total. Ini menaikkan **E-Q16** dari catatan kaki menjadi **pemblokir bagi
seluruh perhitungan SELISIH**. Lihat §6.

### 3.6 Banding NB ✅

`[terverifikasi]` (E-2 §1.6, diulang di sini sebagai titik pijak) `NB FacIn\Activity\SetOldData.xml`
**identik fungsional** dengan salinan EDM di bawah metode 23-tag: 450.829 B, 10.584 baris, 9.183
baris ternormalisasi, `pxObjClass`/`pyClassName`/`pxRuleClassName` sama persis. Perbedaan 1.458 baris
mentah **hanya urutan tag ekspor**, bukan isi.

⛔ **Konsekuensi K-043:** `SetOldData` termasuk 1.195 berkas identik yang **dipakai ulang dari NB** —
tetapi tetap **dicek perilakunya saat integrasi** (pilihan B), dan §3.4–§3.5 adalah tepat jenis hal
yang harus dicek.

---

## 4. Lapis C — tujuh varian `SetOLDValueToEDMWork_<LOB>` ✅

`[terverifikasi]` Nilai penanda tersimpan **berkutip**: `<PropertiesValue>"old"</PropertiesValue>`.

| Varian | `IsOldData` | `IsProRate` | Kedalaman list |
| --- | ---: | ---: | --- |
| `_Aneka` | 14 | 0 | LocationList → OccupationList → AnekaList → CoverageList |
| `_FIRE` | 12 | **2** | LocationList → PropertyItemList → CoverageList → LayerList |
| `_GOLF` | 12 | 0 | LocationList → AnekaList → CoverageList |
| `_MBU` | 9 | 0 | VehicleList → CoverageList |
| `_MC` | 8 | 0 | CargoList → CoverageList |
| `_PA` | 8 | 0 | PersonList → ASMCoverage |
| **`_LIFE`** | **4** | 0 | **PersonList saja** |

### 4.1 Mengapa `_LIFE` hanya 4 `[terverifikasi]`

Pohon langkah `Endorsment Fac In\Activity\SetOLDValueToEDMWork_LIFE.xml`:

```
[1]     IF: IsLife  (aktif=true)
  [1.1] Property-Set   « Copy PersonList LIFE »      IF: IsLife  (aktif=false)
        SET newWorkPage.OfferFacIn.PersonList = newWorkPage.OfferFacIn.OldData.PersonList
  [1.2]                « Set flag old data LIFE »    IF: IsLife  (aktif=false)
    [1.2.1] Property-Set   SET .FlagOldData = "old"
    [1.2.2]
      [1.2.2.1] Property-Set   SET .IsOldData = "old"
```

**Sebabnya struktural, bukan kelalaian:** model data Life hanya punya **satu tingkat list**
(`PersonList`), sedangkan FIRE punya empat tingkat bersarang. Jumlah penanda mengikuti kedalaman
list — bukan kelengkapan implementasi.

⛔ **Temuan baru:** `_LIFE` menyetel **dua penanda berbeda** — `.FlagOldData = "old"` **dan**
`.IsOldData = "old"`. Enam varian lain hanya `.IsOldData`. `[pertanyaan terbuka]` apakah `FlagOldData`
dibaca di tempat lain → **E-Q22**.

📌 Konsisten dengan E-1: `FlagOldData` adalah rule `When` yang **berbeda kelas** antara NB
(`…Data-Accessory`) dan EDM (`…Data-Coverage`) — berkas uji silang metode 23-tag. Kini terlihat ia
juga dipakai sebagai **properti** di jalur Life.

### 4.2 Hanya `_FIRE` menyetel `IsProRate` — konsekuensinya

`[terverifikasi]` `SetOLDValueToEDMWork_FIRE.xml` L416:

```
SET newWorkPage.OfferFacIn.IsProRate = "Prorate"
```

Enam varian lain: **0 kemunculan**.

⛔ **Konsekuensi yang terbaca:** `IsProRate` adalah properti **tingkat `OfferFacIn`** (bukan per
baris), jadi ia menandai **seluruh case**. Untuk case non-FIRE, `IsProRate` **tidak pernah disetel
oleh lapis C** — nilainya bergantung pada apa pun yang tersisa dari lapis A (langkah 14.3 menyalin 54
field dari `OldData`, dan `IsProRate` **tidak** termasuk di antaranya).

`[pertanyaan terbuka]` apakah lini non-FIRE memang tidak diprorata, atau `IsProRate` diisi di jalur
lain → **E-Q17** (diangkat di E-2, kini dengan sebab yang lebih jelas).

---

## 5. E-Q14 TUNTAS — tabel rumus `EDMPremiMenjadi` per cabang

### 5.1 Angka yang benar

`[terverifikasi]` Penelusuran **pohon langkah** atas keenam berkas:

| Ukuran | Jumlah |
| --- | ---: |
| Penugasan `EDMPremiMenjadi` | **22** |
| Rumus **literal** berbeda | **12** |
| Rumus berbeda **setelah normalisasi spasi** | **11** |

⚠️ **Koreksi angka E-2.** E-2 menyebut "11 bentuk" — angka itu **kebetulan benar** secara semantik,
tetapi diperoleh dari penghitungan yang tidak memadai. Angka yang dapat dipertanggungjawabkan:
**22 penugasan · 12 literal · 11 semantik**. Dua literal yang menyatu hanya berbeda spasi di
`(… BrokerageFee -.Policy…)` vs `(… BrokerageFee-.Policy…)`.

### 5.2 Tabel rumus per (lini × `EdmType` × jalur)

⚠️ Arti `EdmType` mengikuti **K-029**: `1` = Batal Sejak Semula · `2` = Batal Prorata ·
`4` = Penambahan/Pengurangan/Perubahan · `3` usang (cabang tetap diport).

#### Kelompok I — `CountEndorsementData` / `CountDataEDMElse` (jalur "batal")

| Berkas | Langkah | Lini (gerbang) | Gerbang `EdmType` | Label langkah | Rumus |
| --- | --- | --- | --- | --- | --- |
| `CountEndorsementData` | 1.2.1 | `IsFire` | `==2` | *"batal sejak semula"* | `EDMOldPremi + EDMNewPremi` |
| `CountEndorsementData` | **1.3.1** | `IsFire` | `==2` | *"batal"* | ⚠️ **`EDMOldPayment + EDMNewPremi`** |
| `CountEndorsementData` | 2.2.1 | `IsAneka` | `==2` | *"batal sejak semula"* | `EDMOldPremi + EDMNewPremi` |
| `CountEndorsementData` | 2.3.1 | `IsAneka` | `==2` | *"batal"* | `EDMOldPremi + EDMNewPremi` |
| `CountEndorsementData` | 3.2.1 | `isGolfInsurance` | `==2` | *"batal sejak semula"* | `EDMOldPremi + EDMNewPremi` |
| `CountEndorsementData` | 3.3.1 | `isGolfInsurance` | `==2` | *"batal"* | `EDMOldPremi + EDMNewPremi` |
| `CountDataEDMElse` | 1.2.1 | `IsMarineCargo` | `==2` | *"batal sejak semula"* | `EDMOldPremi + EDMNewPremi` |
| `CountDataEDMElse` | 1.3.1 | `IsMarineCargo` | `==2` | *"batal"* | `EDMOldPremi + EDMNewPremi` |
| `CountDataEDMElse` | 2.2.1 | ⚠️ `IsMarineCargo` *(label MBU)* | `==2` | *"batal sejak semula"* | `EDMOldPremi + EDMNewPremi` |
| `CountDataEDMElse` | 2.3.1 | ⚠️ `IsMarineCargo` *(label MBU)* | `==2` | *"batal"* | `EDMOldPremi + EDMNewPremi` |
| `CountDataEDMElse` | 3.2.1 | `IsLife` | ⚠️ **`==1 \|\| ==2`** | *"batal sejak semula"* | `EDMOldPremi + EDMNewPremi` |
| `CountDataEDMElse` | **3.3.1** | `IsLife` | `==2` | *"batal"* | ⚠️ **`EDMOldPayment + EDMNewPremi`** |
| `CountDataEDMElse` | 3.4.1.3 | `IsLife` | `==2` + `IsSameCurrency==1` | *"batal"* | `Local.TotalNewPremiNusantaraRe + Local.TotalOldPremiNusantaraRe` |

#### Kelompok II — `CountPaymentEdm_Act` (jalur per **jenis endorsement**, bukan `EdmType`)

| Langkah | Gerbang tingkat atas | Gerbang dalam | Rumus |
| --- | --- | --- | --- |
| 12.1 | `EdmTypeNew==4` *("EDM RI SLIP")* | — | `OldData.CurrencyList(Index).SumTotalPayment` |
| 13.2.3 | `IsEdmExtendPeriod` | `IsFire` ⛔ aktif=`false` | `(preminew − Commision − BrokerageFee − Deduction2 + PPh + PPN) × @Math.divide(datedif, day, 20)` |
| 14.4.3 | `IsEdmAdjRate` | `IsFire` ⛔ aktif=`false` | idem **`+ (oldpayment × @Math.divide(datedifbefore, day, 20))`** |
| 14.4.4 | `IsEdmAdjRate` | `IsFire` + `Local.check==0` *("kalau currency baru")* | `(preminew − Commision) × @Math.divide(datedif, day, 20)` |
| 15.5.3 | `IsEdmAdjPeriod` | `IsFire` ⛔ aktif=`false` | idem 13.2.3 tetapi **`@Math.divide(edmdate, startdate, 20)`** |

#### Kelompok III — sisanya

| Berkas | Langkah | Gerbang | Rumus |
| --- | --- | --- | --- |
| `CountPaymentEdmTSIObj_Act` | 1.4.3 | `IsEdmAdjTSI` + `CalcBrokerFee=="true"` | `(preminew − Commision − BrokerageFee − Deduction2 + PPh + PPN) × @Math.divide(datedif, day, 20)` |
| `CountPaymentEdmTSIObj_Act` | 2.2.3 | `IsEdmAdjInsured` | `Local.oldpayment` |
| `ReCountPremiLifeEDM` | 4.1.9.3 | **(tanpa prakondisi)** | `.Policy.Payment.Premium` |
| `CopyAllObj_ACT` | 9.2.2.1 | `Local.idxOcc==.Name` | `EDMPremiMenjadi + Local.Ujrah` ← **akumulator**, bukan penetapan |

### 5.3 Kesimpulan yang bertahan, dan yang diperbaiki

✅ **Bertahan:** pada cabang berlabel *"batal"*, **FIRE dan LIFE memakai `EDMOldPayment`**; ANEKA,
GOLF, MARINE CARGO, dan MBU memakai `EDMOldPremi`. Kesimpulan E-2 **benar** — tetapi diperolehnya
dengan menebak keanggotaan blok dari rentang nomor baris. Kini dibuktikan dari pohon langkah beserta
gerbangnya, dan **itu** yang membuatnya sah.

⛔ **Diperbaiki:** pembagiannya **bukan** "FIRE+LIFE vs empat lini lain" sebagai sifat lini. Ia
berlaku **hanya pada cabang `x.3` (*"batal"*)**. Pada cabang `x.2` (*"batal sejak semula"*), ketujuh
lini memakai rumus yang **sama**: `EDMOldPremi + EDMNewPremi`.

#### ⛔ E-Q21 — cabang yang menerima `EdmType==1` · **TERBUKTI, dengan koreksi lingkup**

> 🔁 **Koreksi.** Versi pertama menulis *"Dari **13 cabang**, hanya `IsLife` 3.2 bergerbang
> `==1||==2`; 12 lainnya `==2` saja."* Lingkupnya **salah**: 13 itu hanya cabang yang **menugaskan
> `EDMPremiMenjadi`** — bukan seluruh cabang bergerbang `EdmType`. Di kedua berkas sebenarnya ada
> **7 gerbang `==1||==2`**, enam di antaranya langkah **penihilan** yang tidak menugaskan
> `EDMPremiMenjadi` sehingga luput dari pencacahan pertama. Kesimpulan intinya **bertahan**, tetapi
> hanya berlaku untuk **cabang perhitungan pembayaran**.

`[terverifikasi]` Cacah gerbang di kedua berkas:

| Berkas | Gerbang memuat `EdmType==1` | Gerbang `EdmType==2` saja |
| --- | ---: | ---: |
| `CountEndorsementData` | **3** | 6 |
| `CountDataEDMElse` | **4** | 6 |
| **total** | **7** | **12** |

##### Struktur sebenarnya: tiap lini punya tiga langkah

`[terverifikasi]` Pola `x.1` / `x.2` / `x.3` berulang di keenam blok lini. Nomor baris merujuk tag
`<pyStepsPreCondParamsWhen>`; **seluruh `<pyStepsPreCondition>` bernilai `'true'` (aktif)**:

| Lini (gerbang `When`) | `x.1` — *"Set TSI and Premi to 0…"* | `x.2` — *"…'batal sejak semula'"* | `x.3` — *"…'batal'"* |
| --- | --- | --- | --- |
| **FIRE** `IsFire` @L1963 | **`==1\|\|==2`** @L1177 | `==2` @L1547 | `==2` @L1911 |
| **ANEKA** `IsAneka` @L3727 | **`==1\|\|==2`** @L2941 | `==2` @L3311 | `==2` @L3675 |
| **GOLF** `isGolfInsurance` @L5409 | **`==1\|\|==2`** @L4623 | `==2` @L4993 | `==2` @L5357 |
| **MARINE CARGO** `IsMarineCargo` @L1749 | **`==1\|\|==2`** @L963 | `==2` @L1333 | `==2` @L1697 |
| **MBU** ⚠️ `IsMarineCargo` @L3218 | **`==1\|\|==2`** @L2432 | `==2` @L2802 | `==2` @L3166 |
| **LIFE** `IsLife` @L6501 | **`==1\|\|==2`** @L4143 | ⚠️ **`==1\|\|==2`** @L4513 | `==2` @L4877 *(+ `x.4` `==2` @L6449)* |

*(Tiga baris pertama di `CountEndorsementData`; tiga terakhir di `CountDataEDMElse`.)*

**Isi tag persis untuk dua cabang penentu:**

```
# LIFE x.2 — SATU-SATUNYA cabang PERHITUNGAN yang menerima EdmType==1
CountDataEDMElse
  L4202  <pyStepsDescription>Calculate payment if edm type is 'batal sejak semula'</pyStepsDescription>
  L4208  <pyStepsPreCondition>true</pyStepsPreCondition>
  L4513  <pyStepsPreCondParamsWhen>newWorkPage.OfferFacIn.QuotationData.EdmType==1||newWorkPage.OfferFacIn.QuotationData.EdmType==2</pyStepsPreCondParamsWhen>

# FIRE x.2 — pembanding: langkah sejenis, tetapi ==2 saja
CountEndorsementData
  L1236  <pyStepsDescription>Calculate payment if edm type is 'batal sejak semula'</pyStepsDescription>
  L1241  <pyStepsPreCondition>true</pyStepsPreCondition>
  L1547  <pyStepsPreCondParamsWhen>newWorkPage.OfferFacIn.QuotationData.EdmType==2</pyStepsPreCondParamsWhen>
```

##### Gerbang lini = rule `When`; gerbang `EdmType` = literal

`[terverifikasi]` Keduanya, bertingkat:

- **Tingkat lini** memakai **rule `When`**. `Endorsment Fac In\When\IsLife.xml`:
  `<pxObjClass>` = **`Rule-Obj-When`** · `pyClassName` = `ASM-FW-GISFW-Work` ·
  `pxRuleClassName` = `ASM-FW-GISFW-Data-Quotation` · `pyLogic` = `A OR B OR … OR P` (16 cabang),
  seluruhnya `pyWorkPage.Quotation.BusinessOldId = "L1"` … `"L16"`.
- **Tingkat `EdmType`** memakai **ekspresi literal** di `<pyStepsPreCondParamsWhen>`, bukan rule.

📌 Jadi cabang Life dijaga **dua lapis**: rule `IsLife` (L6501) lalu literal `EdmType==1||==2` (L4513).

##### Pernyataan yang benar

✅ **Enam langkah penihilan (`x.1`) berjalan untuk `EdmType==1`** — satu per lini.
⛔ **Dari dua belas cabang perhitungan pembayaran (`x.2`/`x.3`, plus `x.4` Life), hanya LIFE `x.2`
yang berjalan untuk `EdmType==1`.** Sebelas lainnya `==2` saja.

Ini **mempertajam arsip #3**, dan arah temuannya tidak berubah: untuk `EdmType==1` di enam lini
selain Life, TSI/premi **dinihilkan** tetapi **tidak pernah dihitung ulang**. Life adalah
pengecualian yang arsip belum catat. → **E-Q21**.

##### ⚠️ Catatan metodologis — mengapa kedekatan baris tidak bisa dipakai

`[terverifikasi]` Tag-tag milik **satu langkah** terpencar ribuan baris. Contoh langkah `[1]`
`CountEndorsementData`: `<pyStepsDescription>` @L336 · `<pyStepsPreCondition>` @L342 ·
`<pyStepsPreCondParamsWhen>` @**L1963** — selisih **1.627 baris**, karena `pyStepsPreCondParams`
diserialisasi **setelah** seluruh `pySteps` bersarang. Inilah sebab konkret metode E-2 (rentang nomor
baris) gagal, dan mengapa seluruh E-3 memakai penelusuran pohon.

```powershell
# menghasilkan 3/6 dan 4/6
Add-Type -AssemblyName System.Xml.Linq
$XN=[System.Xml.Linq.XName]
foreach($nm in @('CountEndorsementData','CountDataEDMElse')){
  $xd=[System.Xml.Linq.XDocument]::Load("D:\migrasi\RNM\Endorsment Fac In\Activity\$nm.xml",
        [System.Xml.Linq.LoadOptions]::SetLineInfo)
  $a=0;$b=0
  foreach($w in $xd.Descendants($XN::Get('pyStepsPreCondParamsWhen'))){
    if($w.Value -match 'EdmType==1'){$a++} elseif($w.Value -match 'EdmType==2'){$b++} }
  "$nm : memuat ==1 -> $a ; hanya ==2 -> $b" }
```

⛔ **Paradoks label yang terekam** `[terverifikasi]`: langkah `x.2` berlabel
*"Calculate payment if edm type is 'batal sejak semula'"* — yang per **K-029** adalah `EdmType==1` —
bergerbang **`EdmType==2`** (Batal Prorata) pada **5 dari 6 lini**: FIRE (@L1547) · ANEKA (@L3311) ·
GOLF (@L4993) · MARINE CARGO (@L1333) · MBU (@L2802). Hanya **LIFE** (@L4513) yang gerbangnya memuat
`==1`. Label dan gerbang **bertentangan** di lima lini. Penegasan `CLAUDE.md` §3.3 — **yang mengikat
gerbangnya, bukan labelnya**.

```powershell
# menghasilkan 22 penugasan / 12 literal / 11 setelah normalisasi spasi
# (penelusur pohon langkah lengkap ada di OUTPUT\07-edm\ — lihat metode §0 dokumen ini)
```

---

## 6. Prorata dan rumus SELISIH ✅

### 6.1 Prorata dipakai sebagai **pengali nilai baru**, bukan penjumlah

`[terverifikasi]` `Endorsment Fac In\Activity\SaveFacinProdEDMFire_Act.xml` — **19 penugasan**
memakai `ProrateStartEDM`/`ProrateEDMEnd`. Titik pemakaian nyata:

```
SET Datain1.CARI8       = pyWorkPage.OfferFacIn.ProrateEDMEnd
SET Local.ProrateStart  = pyWorkPage.OfferFacIn.ProrateStartEDM

SET Local.TsiSpreadEDM   = .TSISpreaded    * pyWorkPage.OfferFacIn.ProrateEDMEnd
SET Local.PremiSpreadEDM = .PremiumSpreaded * pyWorkPage.OfferFacIn.ProrateEDMEnd
SET Local.NewRIComm      = (.PremiumSpreaded * …Payment.RICommision/100) * …ProrateEDMEnd
SET Local.PremiStart     = @if(…QuotationData.Type=="3", (Local.ProrateStart * .PremiumSpreaded), .PremiumSpreaded)
```

📌 `ProrateStartEDM` **hanya** dipakai lewat `Local.ProrateStart`, dan **hanya bila `Type=="3"`**
(Adj Rate). `ProrateEDMEnd` dipakai luas sebagai pengali seluruh nilai baru.

⛔ **Mengikat implementasi:** `ProrateStartEDM`/`ProrateEDMEnd` adalah **Ratio** (skala 20 desimal,
`@Math.divide(…, 20)` — E-2 §1.4). Ia **hanya mengalikan Money**, tidak pernah dijumlahkan dengannya.
Ini persis jembatan `Money × Ratio → Money`.

### 6.2 Rumus kanonik SELISIH = MENJADI − SEBELUM ✅

`[terverifikasi]` **57 penugasan delta berbentuk pengurangan** di `SaveFacinProdEDMFire_Act`. Tiga
titik pemakaian nyata:

```
# (1) kasus baku — nilai baru sudah diprorata lebih dulu
SET Local.TsiSpreadEDM   = Local.TsiSpreadEDM   - Local.TsiSpreadNB
SET Local.PremiSpreadEDM = Local.PremiSpreadEDM - Local.PremiSpreadNB
SET Datain1.CARI11       = Local.NewRIComm       - Local.OldRIComm
SET Datain1.CARI15       = Local.NewBrokerageFree - Local.OldBrokerageFree

# (2) varian Adj Rate — porsi "sebelum tanggal EDM" ikut ditambahkan sebelum dikurangi
SET Local.PremiSpreadEDM = (Local.TotalPremiNew + Local.PremiStart) - Local.PremiSpreadNB
SET Datain1.CARI11       = (Local.NewRIComm + Local.RICommSTart)    - Local.OldRIComm

# (3) pembalikan tanda untuk EdmType 1/2/3
SET Datain1.CARI11 = (@toDecimal(Datain1.CARI9)  * Local.PremiSpreadEDM/100) * -1
SET Datain1.CARI15 = (@toDecimal(Datain.CARI52) * Local.PremiSpreadEDM/100) * -1
```

✅ Rumus kanoniknya dengan demikian bukan `baru − lama` polos, melainkan:

> **`SELISIH = (nilai_baru × ProrateEDMEnd) − nilai_lama_dengan_TreatyType_sama`**

dengan `nilai_lama` diambil dari `OfferFacIn.OldData` (lapis A), bukan dari `*Old` per baris (lapis B).

📌 **Lapis A dan lapis B melayani konsumen berbeda:** delta produksi (`facinproduction`) membaca
**lapis A**; layar dan perhitungan per baris membaca **lapis B**. Keduanya tidak boleh disatukan.

---

## 7. ⛔ E-Q16 naik menjadi pemblokir perhitungan

`[terverifikasi]` Sebaran prakondisi di seluruh `Endorsment Fac In\Activity\`:

| | Jumlah | di antaranya bergerbang **lini bisnis** |
| --- | ---: | ---: |
| `pyStepsPreCondition` = `true` | 4.905 | 433 |
| `pyStepsPreCondition` = **`false`** | **617** | **47** |

⛔ **47 gerbang lini bisnis dinonaktifkan.** Contoh yang menyentuh langsung pekerjaan E-3:

| Berkas | Langkah | Gerbang yang dinonaktifkan |
| --- | --- | --- |
| `SetOldData` | 4 | `IsFire` ← **seluruh lapis B** |
| `CountPaymentEdm_Act` | 13.2 · 14.4 · 15.5 | `IsFire` ← perhitungan Ext Periode / Adj Rate / Adj Periode |
| `GetdataOldEDMError` | 3.4.1 · 3.4.2 · 3.5.1 | `IsFire` · `IsMarineCargo` ← jalur perbaikan lapis A |
| `CalcultePersentageSpeading_Act` | 2 | `IsPA` |
| `CopyFacRetroAnekaGolf_ACT` | 2.1 · 3.1 | `IsAneka` |

**Dua bacaan, akibat berlawanan:**

| Bila `false` = prakondisi tidak dievaluasi | Bila `false` = langkah dinonaktifkan |
| --- | --- |
| 47 gerbang lini **diabaikan** → langkah jalan untuk **semua** lini | 47 langkah **tidak pernah jalan** → lapis B kosong, perhitungan Adj Rate/Periode tak jalan |

⛔ Tidak ada jalan menulis modul premi endorsement sebelum ini dijawab. **Milik IT.**

---

## 8. Kejanggalan yang diport apa adanya — kandidat perbaikan, bukan perbaikan

Sesuai `CLAUDE.md` §1: perbaikan dipisahkan dari migrasi; memutuskannya milik bisnis.

| # | Kejanggalan | Status verifikasi E-3 | Sikap |
| --- | --- | --- | --- |
| arsip **#3** | Blok `EdmType==1` tak pernah dihitung ulang | ✅ diperkuat (§5.3, bukti baris): 6 langkah **penihilan** menerima `==1`; dari 12 cabang **perhitungan**, **hanya LIFE `x.2`** menerima `==1` | **diport apa adanya** |
| arsip **#4** | `CountDataEDMElse` langkah 2 berlabel MBU bergerbang `IsMarineCargo` | ✅ dikonfirmasi di pohon langkah: `[2] IsMarineCargo` dengan deskripsi *"Calculate payment if MBU business"* | **diport apa adanya** |
| arsip **#5** | `EDMPremiMenjadi` tidak konsisten | ✅ dituntaskan: **22 penugasan / 11 rumus semantik** (§5) | **diport apa adanya** |
| **baru** | `RateOld` guard menguji dirinya sendiri (6/6); `PremiumOld` 1/18 | ✅ §3.4 | **diport apa adanya** · E-Q18 |
| **baru** | `PremiNusantaraReOld` memakai `"0"` string (**2** penugasan, L6294/L7527), **50** lainnya `0` numerik | ✅ §3.4, bukti baris | **diport apa adanya** · E-Q19 |
| **baru** | `<pyExpressionGadget>` (cache editor) memuat `0` tanpa kutip, berbeda dari `<PropertiesValue>` yang dieksekusi | ✅ §3.4 | cache **bukan bukti**; yang mengikat `PropertiesValue` |
| **baru** | Label *"batal sejak semula"* bergerbang `EdmType==2` (= Batal Prorata per K-029) | ✅ §5.3 | **diport apa adanya**; label tidak mengikat |
| **baru** | `SetOldData` menangani Travel tetapi bukan Life; lapis C sebaliknya | ✅ §3.2 | **diport apa adanya** · E-Q22 |

⛔ **Tidak satu pun diperbaiki di sini.** Memperbaikinya akan menghasilkan selisih angka yang tidak
dapat dijelaskan saat rekonsiliasi paralel run.

---

## 9. Pertanyaan terbuka baru dari E-3

| # | Pertanyaan | Pemilik | Dampak |
| ---: | --- | --- | --- |
| **E-Q18** | **`RateOld` guard self-referential** — `@If(.RateOld!="", .Rate, 0)` pada 6 dari 6 kemunculan; `PremiumOld` 1 dari 18. Pada pembukaan layar pertama `RateOld` menjadi `0`. Disengaja? | work owner + Aktuaria | SELISIH rate dan seluruh delta berbasis rate |
| **E-Q19** ✅ **terbukti** | **`PremiNusantaraReOld` memakai `"0"` berkutip** — **2** penugasan (L6294, L7527); **50** penugasan `*Old` lain memakai `0` numerik. Cache `pyExpressionGadget` bahkan memuat bentuk **tanpa kutip**, berbeda dari yang dieksekusi. Perbandingan angka-sebagai-string sudah jadi bahaya terverifikasi (`CLAUDE.md` §4.1) | IT + Aktuaria | Tipe kolom dan perbandingan ambang |
| **E-Q20** | **47 gerbang lini bisnis dinonaktifkan** (`pyStepsPreCondition=false`) dari 617 total. Menaikkan **E-Q16** menjadi pemblokir perhitungan | **IT** | Menentukan apakah lapis B terisi sama sekali |
| **E-Q21** ✅ **terbukti (lingkup dikoreksi)** | Dari **12 cabang perhitungan pembayaran**, **hanya LIFE `x.2`** (`CountDataEDMElse`, gerbang @L4513) menerima `EdmType==1`; sebelas lainnya `==2` saja. Enam langkah **penihilan** (`x.1`) menerima `==1` di semua lini. Disengaja, atau Life satu-satunya yang benar? | work owner + Underwriting | Perhitungan premi pembatalan per lini |
| **E-Q22** | **Himpunan lini lapis B ≠ lapis C** — `SetOldData` menangani **Travel** tanpa Life; `SetOLDValueToEDMWork_*` menangani **Life** tanpa Travel. Dan `_LIFE` menyetel **dua** penanda (`.FlagOldData` + `.IsOldData`) | work owner | Kelengkapan before-image per lini |

Yang **tetap** dari putaran sebelumnya dan belum terjawab: **E-Q13** (`IsNotEDM` bukan negasi
`IsEDM`), **E-Q15** (PII pada `pyStepsRepeatDef`/`pyStepsPreCondParams`), **E-Q16**/**E-Q17**.

---

## 10. Belum dikerjakan (bukan untuk prompt ini)

- **E-4** — 354 berkas EDM-only (Activity 94, Section 81). **Tidak dimulai tanpa perintah.**
- Arsip §4 (gerbang akseptasi EDM, 7 perbedaan) — belum dikonfirmasi ulang; **di sinilah klaim
  "nilai dasar akseptasi = SELISIH TSI" (§4.3) harus diuji**
- Arsip §5 (nilai literal tipe endorsement), §6 (jalur retro), §7 (konversi produksi), §8, §9
- Butir §10 arsip **#12–#23**
- Enam varian `SaveFacinProdEDM*_Act` selain `Fire` — apakah rumus deltanya sama
- Isi 16 berkas tipe yang 100 % berbeda (E-1 §2), termasuk `DecisionTable\IsUWAccepted`

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
