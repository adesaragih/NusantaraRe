# Discovery Endorsement — E-2: verifikasi ulang alur endorsement ke korpus EDM

> **Sumber:** `D:\migrasi\RNM\Endorsment Fac In\` (banding ke `NB FacIn\`). READ-ONLY.
> Korpus Treaty (`RNM_BRD\`) **tidak dibaca** (K-005). Label mengikuti `CLAUDE.md` §3.
>
> **Ini bukan pemetaan dari nol.** Alur endorsement sudah terpetakan di
> `_ARSIP-lintas-siklus\01-flow\03-alur-endorsement.md` (§0–§10). K-042 mewajibkan tiap temuannya
> dikonfirmasi ulang ke korpus EDM. Dokumen ini = **hasil konfirmasi ulang**, bukan penemuan ulang.
>
> Lingkup kerja EDM sudah ditetapkan **K-043**: 866 berkas (512 beda + 354 EDM-only); 1.195 identik
> dipakai ulang dari NB **dengan cek perilaku saat integrasi** (pilihan B), bukan diterima buta.

---

## 1. Hasil verifikasi enam klaim inti

| # | Klaim arsip | Status | Bukti ringkas |
| ---: | --- | --- | --- |
| **a** | `IsEDM` / `IsNotEDM` menandai siklus lewat `StatusBusiness = 3` | ✅ **dikonfirmasi ulang di EDM** — **+1 temuan baru** | §1.1 |
| **b** | Konektor `Start2 → Assignment7` menyetel 6 properti, EDM lahir di fase Policy | ✅ **dikonfirmasi ulang di EDM** | §1.2 |
| **c** | Tiga lapis before-image (A `OfferFacIn.OldData` · B properti `*Old` · C penanda `.IsOldData`) | ✅ **dikonfirmasi ulang di EDM** — 1 koreksi metode baca | §1.3 |
| **d** | Prorata langkah 15 — `@Math.divide(…, 20)` | ✅ **dikonfirmasi ulang di EDM**, persis | §1.4 |
| **e** | Delta kanonik + pasangan kolom `*_MENJADI`/`*_SELISIH` di `facinproduction` | ✅ **dikonfirmasi ulang di EDM** — **1 koreksi tipe rule** | §1.5 |
| **f** | `SetOldData` identik NB↔EDM | ✅ **dikonfirmasi ulang dengan metode 23-tag** — identik **fungsional** | §1.6 |

**Enam dari enam bertahan.** Tidak ada klaim arsip yang gugur. Yang berubah: satu klasifikasi tipe
rule, satu cara membaca nilai, dan satu temuan baru yang arsip belum catat.

---

### 1.1 Klaim (a) — penanda siklus `StatusBusiness = 3` ✅

`[terverifikasi]` `Endorsment Fac In\When\IsEDM.xml`:

```
<pxObjClass>        Rule-Obj-When
<pyClassName>       ASM-FW-GISFW-Work
<pyDesignatedClass> ASM-FW-GISFW-Work-Endorsement
<pxRuleClassName>   ASM-FW-GISFW-Data-OfferFacIn
<pyLogic>           A
pyConditionValue1String[1] = pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3
```

`[terverifikasi]` `Endorsment Fac In\When\IsNotEDM.xml`:

```
<pxObjClass>        Rule-Obj-When
<pyClassName>       ASM-FW-GISFW-Work
<pyDesignatedClass> ASM-FW-GISFW-Work-Endorsement
<pxRuleClassName>   ASM-FW-GISFW-Data-Quotation
<pyLogic>           A
pyConditionValue1String[1] = pyWorkPage.Quotation.StatusBusiness != 3
```

#### ⛔ Temuan baru — `IsNotEDM` **bukan** negasi `IsEDM`

`[terverifikasi]` Keduanya membaca **jalur properti yang berbeda**:

| Rule | Jalur yang dibaca |
| --- | --- |
| `IsEDM` | `pyWorkPage.**OfferFacIn.QuotationData**.StatusBusiness` — agregat tersimpan |
| `IsNotEDM` | `pyWorkPage.**Quotation**.StatusBusiness` — halaman aktif |

⚠️ Karena sumbernya berbeda, **keduanya dapat bernilai benar bersamaan, atau salah bersamaan**, bila
kedua halaman tidak sinkron. Di sistem baru, `IsNotEDM` **tidak boleh** diimplementasikan sebagai
`!IsEDM` — itu akan mengubah perilaku. Diajukan sebagai **E-Q13**.

📌 Konsisten dengan temuan E-1 §5: jalur "halaman aktif" dan "agregat tersimpan" hidup berdampingan,
bahkan di antara dua rule yang namanya saling meniadakan.

#### Catatan banding NB

`[terverifikasi]` `IsEDM` di NB **identik seluruh tag kelas dan kondisinya** dengan EDM. `IsNotEDM`
kondisinya identik, tetapi `pyClassName` berbeda: NB `ASM-FW-GISFW-Data-Vehicle`, EDM
`ASM-FW-GISFW-Work`. `[pertanyaan terbuka]` dampaknya belum ditelusuri.

```powershell
foreach($d in @('NB FacIn','Endorsment Fac In')){ foreach($r in @('IsEDM','IsNotEDM')){
  $x=[xml][IO.File]::ReadAllText("D:\migrasi\RNM\$d\When\$r.xml")
  "$d\$r : " + ($x.SelectSingleNode('//pyConditionValue1String')).InnerText.Trim() } }
```

---

### 1.2 Klaim (b) — gerbang masuk flow ✅

`[terverifikasi]` `Endorsment Fac In\Flow\InputAddendumFacultativeIn.xml` (`pxObjClass` =
`Rule-Obj-Flow`, **69 shape / 138 konektor**). Konektor **`pyFrom=Start2 → pyTo=Assignment7`**:

```
SET .FlagOnGoingPolicy = 1
SET .IsCedingConfirm   = Policy
SET .Position          = 1
SET .PositionNote      = "ReasFacInMarketing"
SET .NBStatus          = "NEW EDM"
SET .NBStatusNew       = "NEW EDM"
```

Keenamnya cocok dengan arsip §1.5, **termasuk urutannya**.

`[terverifikasi]` Shape `Assignment7`:

| Tag | Nilai |
| --- | --- |
| `pxObjClass` | `Data-MO-Activity-Assignment` |
| `pyMOName` | **`MARKETING`** |
| `pyImplementation` | `WorkBasket` |
| `pyRouteTo` | `Custom` |
| `pyIsCustomRouter` | `false` |
| `pyParametersParamName` / `pyCallParams` | `Workbasket` = **`ReasFacInMarketing`** |
| `pyTicketShapes` | `…ExceptionAdminPolicyTicket1…` → tiket **`AdminPolicy`** |

`[terverifikasi]` Token `ToWorkbasket` muncul **4×** di flow, dan `pyIsCustomRouter=false` —
router bawaan Pega, bukan router kustom.

⚠️ `[dugaan]` Arsip menyebut router `ToWorkbasket` untuk `Assignment7` secara spesifik; korpus
membuktikan router bawaan dipakai di flow ini, tetapi **pengikatan ke `Assignment7` khusus belum
saya kutip dari tag**. Diperlakukan `[dugaan]`, bukan `[terverifikasi]`.

```powershell
$x=[xml](Get-Content "D:\migrasi\RNM\Endorsment Fac In\Flow\InputAddendumFacultativeIn.xml" -Raw)
$x.SelectNodes('/pagedata/pyModelProcess/pyConnectors/rowdata') | ? {
  $_.SelectSingleNode('pyFrom').InnerText -eq 'Start2' } | % {
  $_.SelectNodes('pyPropertyAssigns/rowdata') } | % {
  $_.SelectSingleNode('pyPropertiesName').InnerText + ' = ' + $_.SelectSingleNode('pyPropertiesValue').InnerText }
```

---

### 1.3 Klaim (c) — tiga lapis before-image ✅

⚠️ **Istilah.** `[terverifikasi]` Token literal `BeforeImage` / `AfterImage` **tidak ada di korpus
EDM (0 berkas)**. "Before-image" adalah **istilah konsep** di dokumen ini, bukan nama tag maupun
properti. Jangan menulis seolah ada properti bernama itu.

#### Lapis A — `OfferFacIn.OldData` ✅

`[terverifikasi]` `Endorsment Fac In\Activity\SetValueToEDMWork.xml`, langkah **14**
(`pyStepsDescription` = *"Copy policy data to old data"*, `pyStepsPreCondition` = **`false`**),
**13 sub-langkah**:

| Sub | Bukti dari tag |
| --- | --- |
| 14.1 | `RDB-List` · `<RequestType>` = **`GetEDMOldData_SQL`** · `<ClassName>` = `ASM-FW-GISFW-Int-OFFERJSON` · `<BrowsePage>` = **`OldData`** · `pyStepsPreCondition` = `false` |
| 14.2 | `Java` — *"copy NB and remove the comments, and fac retro list"* |
| 14.3 | `Property-Set` — **tepat 54 `<PropertiesName>`**, semuanya bertarget `newWorkPage.OfferFacIn.*` |
| 14.4–14.5 | `RDB-List` + `Property-Set` — *"-- Set team Group"* |
| 14.6 | **`Obj-Save`** |
| 14.7–14.13 | **7× `Call SetOLDValueToEDMWork_<LOB>`** — `_FIRE` · `_MC` · `_Aneka` · `_MBU` · `_LIFE` · `_PA` · `_GOLF`, masing-masing `pyStepsPreCondition` = `true` |

✅ Angka **54** cocok persis dengan arsip §2.1.1. Ketujuh varian LOB ada sebagai berkas di
`Endorsment Fac In\Activity\`.

#### Lapis B — properti `*Old` per baris ✅

`[terverifikasi]` `Endorsment Fac In\Activity\InputAddendumFacIn_PreAct.xml` — **30 langkah**:

```
[5]  Call GetdataOldEDMError        preCond=true      ← jalur perbaikan
[27] Call GetHistoryAkseptasiPega_Act
[28] Call SetOldData                « SetOldData »    ← lapis B
[29] Call MappingKlaimToLossRecord_Act
```

✅ Langkah **28 = `Call SetOldData`**, persis seperti arsip §2.2.

#### Lapis C — penanda `.IsOldData = "old"` ✅ (dengan koreksi cara baca)

`[terverifikasi]` `Endorsment Fac In\Activity\SetOLDValueToEDMWork_FIRE.xml`:

```
L416   <PropertiesName>newWorkPage.OfferFacIn.IsProRate</PropertiesName>
       <PropertiesValue>"Prorate"</PropertiesValue>
L526   <pyStepsDescription>Set IsOldData FIRE</pyStepsDescription>
L585   <PropertiesName>.IsOldData</PropertiesName>   <PropertiesValue>"old"</PropertiesValue>
L763   …idem…      L935 …idem…      L1253 …idem…      L2776 …idem…
L1768  <PropertiesName>TempSpread.pxResults(&lt;APPEND&gt;).IsOldData</PropertiesName>
       <PropertiesValue>.IsOldData</PropertiesValue>      ← diteruskan, bukan diset
```

⚠️ **Koreksi metode baca, bukan koreksi klaim.** Pencarian `>old<` memberi **0 hasil** di ketujuh
varian dan sempat menyesatkan. Nilainya tersimpan **berkutip**: `"old"`. Pola pencarian harus
menyertakan kutip.

`[terverifikasi]` Sebaran `IsOldData` pada ketujuh varian:

| Varian | `IsOldData` | `IsProRate` |
| --- | ---: | ---: |
| `_Aneka` | 14 | 0 |
| `_FIRE` | 12 | 2 |
| `_GOLF` | 12 | 0 |
| `_MBU` | 9 | 0 |
| `_MC` | 8 | 0 |
| `_PA` | 8 | 0 |
| `_LIFE` | **4** | 0 |

📌 `_LIFE` paling sedikit (4) dan **hanya `_FIRE`** yang menyetel `IsProRate = "Prorate"`.
`[pertanyaan terbuka]` apakah itu disengaja. `[terverifikasi]` `IsOldData` disebut di **50 berkas**
EDM seluruhnya.

---

### 1.4 Klaim (d) — prorata ✅ persis

`[terverifikasi]` `Endorsment Fac In\Activity\SetValueToEDMWork.xml` langkah **15**
(*"Count StartProRate and EndProRate"*, `pyStepsPreCondition` = `true`):

```
L5319  <PropertiesName>Local.TotalPeriod</PropertiesName>
L5320  <PropertiesValue>@if(Local.TotalPeriod = 0,1,Local.TotalPeriod)</PropertiesValue>
L5340  <PropertiesName>newWorkPage.OfferFacIn.ProrateStartEDM</PropertiesName>
L5341  <PropertiesValue>@Math.divide(Local.EdmToStart, Local.TotalPeriod, 20)</PropertiesValue>
L5367  <PropertiesName>newWorkPage.OfferFacIn.ProrateEDMEnd</PropertiesName>
L5368  <PropertiesValue>@Math.divide(Local.EdmToEnd, Local.TotalPeriod, 20)</PropertiesValue>
```

✅ Cocok persis dengan arsip §2.4, termasuk **guard pembagian-nol** (`TotalPeriod=0 → 1`) dan
**skala 20 desimal**.

⛔ **Mengikat implementasi:** 20 desimal berarti tipe **desimal presisi tinggi**, bukan `float`
(`CLAUDE.md` §4.1). `ProrateStartEDM`/`ProrateEDMEnd` adalah **Ratio**, bukan Money — ia hanya boleh
mengalikan Money, tidak pernah dijumlahkan dengannya.

```powershell
Select-String -Path "D:\migrasi\RNM\Endorsment Fac In\Activity\SetValueToEDMWork.xml" `
  -Pattern 'ProrateStartEDM|ProrateEDMEnd|Math\.divide|Local\.TotalPeriod'
```

---

### 1.5 Klaim (e) — delta dan pasangan kolom produksi ✅ (+ koreksi tipe rule)

#### ⛔ Koreksi — `InsertTreatyProduction_Sql` **bukan** RDB-List

`[terverifikasi]` Berkas `Endorsment Fac In\**RDBList**\InsertTreatyProduction_Sql.xml`:

```
<pxObjClass>     Rule-Connect-SQL        ← BUKAN Rule-Obj-RDBList
<pyClassName>    ASM-FW-GISFW-Int-policyjson
<pyRequestType>  InsertTreatyProduction_Sql
```

📌 **Jebakan ketujuh keluarga "nama bukan bukti": nama FOLDER bukan tipe rule.** Berkas duduk di
`RDBList\` tetapi `pxObjClass`-nya `Rule-Connect-SQL`. Klasifikasi wajib dari `pxObjClass`.
Ini melanjutkan daftar: beda tipe rule (`isApproved`) · nama sama tipe beda (`GetInsuredID`) ·
posisi kata (`TSIOld`) · sufiks (`InputDtlObject`) · kapitalisasi (K-042) · kelas sumber data
(E-1 §6) · **folder**.

#### Pasangan kolom ✅ — arsip benar, termasuk salah ejaannya

`[terverifikasi]` `INSERT INTO facinproduction` memuat **82 kolom**. Berakhiran suffix:

| | Jumlah | Kolom |
| --- | ---: | --- |
| `*_MENJADI` | **5** | `tsi_menjadi` · `premi_menjadi` · `LOL_MENJADI` · `BROKERAGE_FEE_MENJADI` · `TSI100_MENJADI` |
| `*_SELISIH` | **8** | `tsi_selisih` · `premi_selisih` · `pct_ri_comm_selisih` · `ricomm_selisih` · `LOL_SELISIH` · `BROKERAGE_FEE_SELISIH` · `PCT_BROKERAGE_FEE_SELISIH` · `TSI100_SELISIH` |

⚠️ **Tidak setangkup — dan itu memang begitu.** Tiga `*_SELISIH` tidak punya pasangan `*_MENJADI`
karena nilai "sesudah"-nya disimpan di kolom **tanpa sufiks**. Arsip §3.3.1 **sudah mencatatnya
dengan benar**; verifikasi ini menegaskan, bukan mengoreksi:

| `*_SELISIH` | Kolom "sesudah" pasangannya |
| --- | --- |
| `ricomm_selisih` | `ricomm` (tanpa sufiks) |
| `pct_ri_comm_selisih` | `percent_ri_comm` (**ejaan berbeda**: `percent_` vs `pct_`) |
| `PCT_BROKERAGE_FEE_SELISIH` | **`PCT_BROKERGARE_FEE`** — ⚠️ **salah ejaan di skema** (`BROKERGARE`, bukan `BROKERAGE`) |

⛔ **`PCT_BROKERGARE_FEE` diport apa adanya.** Skema Oracle tidak berubah (`CLAUDE.md` §4.3), dan
memperbaiki ejaan di sisi aplikasi akan memutus `INSERT`. Dicatat sebagai **kandidat perbaikan milik
bisnis/DBA**, bukan perbaikan diam-diam (§1).

⛔ **Mengikat implementasi:** seluruh kolom nilai dibungkus
`To_number(Replace({…},',','.'))` — **koma desimal**, sesuai kontrak K-027. Parser wajib mengikuti.

```powershell
$t=Get-Content "D:\migrasi\RNM\Endorsment Fac In\RDBList\InsertTreatyProduction_Sql.xml" -Raw
([regex]::Match($t,'INSERT INTO\s+facinproduction\s*\(([^)]*)\)','Singleline').Groups[1].Value -split ',').Trim() |
  ? { $_ -match '(?i)_MENJADI$|_SELISIH$|BROKER|ricomm' }
```

---

### 1.6 Klaim (f) — `SetOldData` identik NB↔EDM ✅ **identik fungsional**

`[terverifikasi]` Diuji dengan metode **23-tag reproducible** (kontrak di
`07-edm\02-e1-pola-perbedaan-when.md` §1), **bukan** hash mentah:

| Ukuran | NB | EDM |
| --- | ---: | ---: |
| Ukuran berkas | 450.829 B | **450.829 B** |
| Baris mentah | 10.584 | **10.584** |
| Baris setelah normalisasi | 9.183 | **9.183** |
| SHA-256 **mentah** | `2CEE2807…` | `ED84460D…` — **berbeda** |
| **Verdikt 23-tag** | \<— **IDENTIK FUNGSIONAL** —\> | |
| `pxObjClass` / `pyClassName` / `pxRuleClassName` | `Rule-Obj-Activity` / `ASM-FW-GISFW-Work` / `ASM-FW-GISFW-Data-OfferFacIn` | **sama persis** |

#### ⛔ Temuan metodologis penting — **urutan tag ekspor tidak deterministik**

`[terverifikasi]` Kedua berkas berbeda di **1.458 baris mentah dari 10.584**, tetapi identik
fungsional. Sebabnya bukan isi, melainkan **urutan**:

```
L5  NB : <pyRuleName>SetOldData</pyRuleName>
    EDM: <pxUpdateSystemID>pega</pxUpdateSystemID>
L6  NB : <pxUpdateSystemID>pega</pxUpdateSystemID>
    EDM: <pyRuleName>SetOldData</pyRuleName>          ← tag yang sama, urutan tertukar
```

📌 Ini **membuktikan langkah sort Ordinal** dalam metode 23-tag bukan kosmetik melainkan **wajib**.
Tanpa sort, 1.458 baris tampak "berbeda" padahal tidak ada perbedaan isi sama sekali. Arsip sudah
memperingatkan hal ini (§4.2 baris 641: *"bandingkan pohon langkah — bukan hash mentah, urutan tag
ekspor tidak deterministik"*); verifikasi ini mengkuantifikasinya.

⚠️ Sekaligus penjelasan tambahan untuk **"0 identik byte-per-byte"** pada K-042: penyebabnya bukan
hanya stempel waktu, tetapi juga **pengacakan urutan tag**.

---

## 2. Pertanyaan terbuka §10 arsip — yang ditutup

### ✅ #2 — arti `EdmType` 1/2/3 → **DITUTUP oleh K-029**

Tidak ditebak ulang. K-029 sudah menetapkan (dari `<pyPromptTableList>`):
**1 = Batal Sejak Semula · 2 = Batal Prorata · 4 = Penambahan/Pengurangan/Perubahan**; kode **3
usang**, tetapi **cabang `EdmType==3` diport apa adanya**.

📌 Konflik deskripsi langkah yang mendasari #2 (satu langkah melabeli `EdmType==1` "Batal Sejak
Semula", langkah lain melabeli `EdmType==2` "batal sejak semula") **tetap ada di korpus** — tetapi
kini berstatus **label menyesatkan**, bukan pertanyaan terbuka. Yang mengikat adalah `pyPromptTableList`,
bukan `pyStepsDescription`. Penegasan `CLAUDE.md` §3.3.

### ✅ #3 — blok `EdmType==1` tidak pernah dihitung ulang → **terkonfirmasi, diport apa adanya**

`[terverifikasi]` Pola berulang di **kedua** aktivitas:

| Berkas | `EdmType==1\|\|EdmType==2` | `EdmType==2` saja |
| --- | ---: | ---: |
| `Endorsment Fac In\Activity\CountEndorsementData.xml` | 3× (L1177, L2941, L4623) | 6× (L1547, L1911, L3311, L3675, L4993, L5357) |
| `Endorsment Fac In\Activity\CountDataEDMElse.xml` | 4× (L963, L2432, L4143, L4513) | 4× (L1333, L1697, L2802, L4877, L6449) |

Setiap blok penihilan `==1||==2` diikuti blok perhitungan ulang berkondisi **`==2` saja**. Untuk
`EdmType==1` (Batal Sejak Semula), TSI/premi dinihilkan tetapi `CurrencyList` **tidak pernah**
dihitung ulang.

⛔ **Diport apa adanya** (`CLAUDE.md` §1) — kandidat perbaikan milik bisnis, bukan diperbaiki diam-diam.

### ✅ #4 — langkah berlabel MBU bergerbang `IsMarineCargo` → **terkonfirmasi, diport apa adanya**

`[terverifikasi]` `Endorsment Fac In\Activity\CountDataEDMElse.xml`:

| Langkah | `pyStepsDescription` | Gerbang (`pyStepsPreCondParamsWhen`) |
| ---: | --- | --- |
| 1 | *"Calculate payment if MARINE CARGO business"* | `IsMarineCargo` |
| **2** | *"Calculate payment if **MBU** business"* | **`IsMarineCargo`** ← sama dengan langkah 1 |
| 3 | *"Calculate payment if Life business"* | `IsLife` |

Akibat yang terbaca dari struktur: endorsement **MBU tidak mendapat perhitungan pembayaran**, dan
Marine Cargo **dihitung dua kali**. ⛔ **Diport apa adanya**, kandidat perbaikan.

📌 Ini contoh murni `CLAUDE.md` §3.3: **nama/label langkah bukan bukti perilaku**; yang mengikat
adalah `pyStepsPreCondParamsWhen`.

### ✅ #5 — `EDMPremiMenjadi` tidak konsisten → **terkonfirmasi dan DIPERTAJAM**

> 🔁 **Dituntaskan di E-3** — lihat `04-e3-before-image-selisih.md` §5. Angka yang sah:
> **22 penugasan · 12 rumus literal · 11 rumus semantik**. Kesimpulan "FIRE + LIFE memakai
> `EDMOldPayment`" **bertahan**, tetapi berlaku **hanya pada cabang `x.3` (*"batal"*)**, bukan
> sebagai sifat lini. Tabel di bawah disusun dari **rentang nomor baris** — metode yang tidak
> memadai; E-3 menggantinya dengan penelusuran **pohon langkah + gerbang**.

`[terverifikasi]` Seluruh penugasan `.Policy.Payment.EDMPremiMenjadi` di korpus EDM. Arsip menyebut
"FIRE vs ANEKA/GOLF"; korpus menunjukkan pembagiannya **FIRE + LIFE vs empat lini lain**:

| Berkas · baris | Lini (blok) | Rumus |
| --- | --- | --- |
| `CountEndorsementData` L1435 | FIRE (blok `==1\|\|==2`) | `EDMOldPremi + EDMNewPremi` |
| **`CountEndorsementData` L1799** | **FIRE** (blok `==2`) | **`EDMOldPayment + EDMNewPremi`** |
| `CountEndorsementData` L3199 · L3563 | ANEKA | `EDMOldPremi + EDMNewPremi` |
| `CountEndorsementData` L4881 · L5245 | GOLF | `EDMOldPremi + EDMNewPremi` |
| `CountDataEDMElse` L1221 · L1585 | MARINE CARGO | `EDMOldPremi + EDMNewPremi` |
| `CountDataEDMElse` L2690 · L3054 | MBU | `EDMOldPremi + EDMNewPremi` |
| `CountDataEDMElse` L4401 | LIFE | `EDMOldPremi + EDMNewPremi` |
| **`CountDataEDMElse` L4765** | **LIFE** (blok `==2`) | **`EDMOldPayment + EDMNewPremi`** |

Ditambah **rumus yang sama sekali lain** di berkas lain `[terverifikasi]`:

| Berkas | Rumus `EDMPremiMenjadi` |
| --- | --- |
| `CountDataEDMElse` L6233 | `Local.TotalNewPremiNusantaraRe + Local.TotalOldPremiNusantaraRe` |
| `CountPaymentEdm_Act` L2274 | `OldData.CurrencyList(Index).SumTotalPayment` |
| `CountPaymentEdm_Act` L3331 · L5666 · L5927 · L8527 | ekspresi prorata `(premi − komisi − brokerage − deduction + PPh + PPN) × @Math.divide(…, 20)`, empat varian berbeda |
| `CountPaymentEdmTSIObj_Act` L1795 · L3571 | idem · `Local.oldpayment` |
| `ReCountPremiLifeEDM` L3228 | `.Policy.Payment.Premium` |
| `CopyAllObj_ACT` L5899 | `EDMPremiMenjadi + Local.Ujrah` (akumulasi, bukan penetapan) |

⛔ **Diport apa adanya.** Satu properti, **sebelas bentuk perhitungan berbeda**. Ini lebih besar dari
catatan arsip dan **harus dirancang sebagai tabel rumus per (lini × EdmType × jalur)**, bukan satu
fungsi. Diajukan sebagai **E-Q14**.

### ✅ #9 — `GenerateEndorsementNo` / `GenerateEDMNoLife` → **GAGAL KERAS (`panic`)**

`[terverifikasi]`

| Nama | Dirujuk | Berkas rule-nya |
| --- | --- | --- |
| `GenerateEndorsementNo` | `Activity\SaveEDMToJsonPolicy_Act` (1×) **dan** `Activity\GenerateNopolis_Act` | **NIHIL di seluruh `D:\migrasi\RNM\`** |
| `GenerateEDMNoLife` | `Activity\SaveEDMToJsonPolicy_Act` (1×) | **NIHIL di seluruh `D:\migrasi\RNM\`** |

Per `CLAUDE.md` §4.5: rule yang **dirujuk tetapi berkasnya tidak ada di folder mana pun** →
**`panic`**, bukan ditebak, bukan di-stub, bukan dihapus cabangnya.

```go
// Asal: Endorsment Fac In/Activity/SaveEDMToJsonPolicy_Act.xml
panic("GenerateEndorsementNo: rule tidak ada di korpus mana pun")
```

⚠️ `[terverifikasi]` **`GenerateEndorsementNo` dirujuk di dua berkas, bukan satu** — arsip hanya
menyebut satu. Tambahan rujukan di `GenerateNopolis_Act` memperluas dampaknya ke penomoran polis,
bukan hanya penomoran endorsement.

📌 Catatan lingkup: pencarian **tidak** menyentuh `RNM_BRD\` (K-005). Bila rule itu ada di korpus
Treaty, temuan ini tetap berlaku untuk Fac In.

### ⏸ #24 — `<pyStepsPreCondition>false</pyStepsPreCondition>` → **TETAP TERBUKA** (milik IT)

`[terverifikasi]` Hitungan arsip **cocok persis**:

| Berkas | Jumlah |
| --- | ---: |
| `Endorsment Fac In\Activity\SetValueToEDMWork.xml` | **4** |
| `Endorsment Fac In\Activity\SaveFacinProdEDMFire_Act.xml` | **17** |

⛔ **Tetapi skalanya jauh lebih besar dari catatan arsip.** `[terverifikasi]` di seluruh
`Endorsment Fac In\Activity\`: **617 kemunculan di 173 berkas**.

Ini **perilaku engine Pega**, bukan data korpus — korpus tidak dapat menjawabnya. Tetap
`[pertanyaan terbuka]` untuk **IT**, tetapi dengan bobot yang berubah: jawabannya menentukan apakah
**617 langkah** berjalan atau dilewati. Bila "false" berarti *langkah dinonaktifkan*, sebagian besar
alur yang dipetakan di dokumen ini tidak berjalan seperti yang terbaca.

```powershell
$t=0;$b=0
Get-ChildItem "D:\migrasi\RNM\Endorsment Fac In\Activity" -Filter *.xml -File | % {
  $n=([regex]::Matches([IO.File]::ReadAllText($_.FullName),'<pyStepsPreCondition>false</pyStepsPreCondition>')).Count
  if($n){$t+=$n;$b++} }; "$t kemunculan di $b berkas"   # => 617 kemunculan di 173 berkas
```

---

## 3. Pertanyaan §10 yang **TETAP MEMBLOKIR** — tidak ditutup

Sesuai arahan work owner, keenam butir berikut **tidak disentuh** dan tetap memblokir:

| # | Pertanyaan | Mengapa korpus tidak bisa menjawab |
| ---: | --- | --- |
| **1** | Versi rule mana yang berlaku di produksi (`SetToInbox_ACT` / `SetBanding_ACT`) | Butuh **satu ekspor tunggal dari sistem produksi** — bukan soal membaca lebih teliti |
| **6** | Endorsement Life melewati seluruh tangga persetujuan | Keputusan bisnis, bukan fakta korpus |
| **7** | `Decision39` memutus `SaveEDMToJsonPolicy_Act` untuk case ber-fac-retro | Disengaja atau tidak = keputusan bisnis |
| **8** | Isi 4 stored procedure (`INSERTJSONPOLIS`, `PEGA_DELETE_ERROR_KONVERSI`, `PROC_GENERATE_SEQUENCE_NUMBER`, `FacinForBackup`) | Tidak ada di korpus · **DBA** |
| **10** | Isi tabel `OPENPROTEKSI_EDM` + otorisasi pengisinya | Tidak ada di korpus · **DBA + bisnis** |
| **11** | Kerapatan `PRODKE` (jaminan tanpa lubang) | Constraint basis data · **DBA** |

---

## 4. Verifikasi ulang P2 dan tiga perbedaan perilaku

### 4.1 ⛔ Koreksi yang saya terapkan atas laporan E-1 saya sendiri

E-1 §6 menulis *"beda KELAS DEKLARASI rule"*. **Itu salah**, dan koreksi work owner benar.

`[terverifikasi]` `pxObjClass` keenam predikat COB = **`Rule-Obj-When` di NB maupun EDM** — tipe dan
deklarasinya **sama**. Yang berbeda adalah **kelas halaman yang DIBACA kondisinya**:

| Rule | `pxObjClass` NB → EDM | `pyDesignatedClass` EDM | Sumber data kondisi EDM |
| --- | --- | --- | ---: |
| `IsAneka` | `Rule-Obj-When` → `Rule-Obj-When` | `Code-Pega-List` | 25/25 `OutData.pxResults` |
| `IsFire` | `Rule-Obj-When` → `Rule-Obj-When` | `Code-Pega-List` | 5/5 |
| `IsPA` | `Rule-Obj-When` → `Rule-Obj-When` | `ASM-FW-GISFW-Data-Search` | 1/1 |
| `IsGolfInsurance` | `Rule-Obj-When` → `Rule-Obj-When` | `Code-Pega-List` | 1/1 |
| `IsMarineCargo` | `Rule-Obj-When` → `Rule-Obj-When` | `Code-Pega-List` | 1/1 |
| `IsMBU` | `Rule-Obj-When` → `Rule-Obj-When` | `Code-Pega-List` | 4/4 |

`Code-Pega-List` adalah kelas **halaman hasil query** — konsisten dengan `OutData.pxResults(n)`.
Jadi P2 = **perubahan sumber data yang dibaca**, bukan perubahan tempat rule dideklarasikan.

### 4.2 ⛔ Koreksi kedua — delegasi-vs-literal **TIDAK seragam**

E-1 §6.1 menyimpulkan *"NB mendelegasikan ke sub-rule, EDM membandingkan literal"* dari dua contoh.
**Generalisasi itu salah.** `[terverifikasi]` hitungan penuh keenamnya:

| Rule | NB delegasi | NB literal/properti | EDM delegasi | EDM literal |
| --- | ---: | ---: | ---: | ---: |
| `IsAneka` | **25** | 1 | 0 | 25 |
| `IsFire` | **4** | 1 | 0 | 5 |
| `IsPA` | **0** | 1 | 0 | 1 |
| `IsGolfInsurance` | **0** | 1 | 0 | 1 |
| `IsMarineCargo` | **0** | 2 | 0 | 1 |
| `IsMBU` | **0** | 2 | 0 | 4 |

📌 Pernyataan yang benar: **NB mendelegasikan HANYA di `IsFire` dan `IsAneka`**; empat lainnya
membaca properti langsung. EDM memakai literal pada keenamnya. Dokumen E-1 §6.1 sudah dikoreksi.

### 4.3 Tiga perbedaan perilaku — dikonfirmasi ulang di korpus EDM

`[terverifikasi]` Ketiganya `pxObjClass` = `Rule-Obj-When` di kedua sisi.

**`IsUW`** — `pyDesignatedClass` = `Code-Pega-Requestor` **di kedua sisi** (sumber data sama), jadi
murni perbedaan himpunan kondisi:

| | `pyLogic` | Kondisi |
| --- | --- | --- |
| `NB FacIn\When\IsUW.xml` | `A OR B OR C` | `ReasFacInDirector` · `ReasFacInGroupLeader` · `ReasFacInUnderwriting` |
| `Endorsment Fac In\When\IsUW.xml` | **`A OR B OR C OR D`** | `ReasFacInDirector` · **`ReasFacInMarketing`** · `ReasFacInGroupLeader` · `ReasFacInUnderwriting` |

⛔ Di endorsement, operator ber-workbasket **marketing dihitung sebagai underwriter**. Menyatu dengan
§1.2: EDM **lahir** di workbasket `ReasFacInMarketing`. `[pertanyaan terbuka]` apakah itu sebab
akibat yang disengaja — **E-Q6**, milik Underwriting.

**`IsClaim`** — di sini `pyDesignatedClass` memang berbeda, sejalan dengan kondisinya:

| | `pyDesignatedClass` | Kondisi |
| --- | --- | --- |
| `NB FacIn\When\IsClaim.xml` | `ASM-FW-GCNMFW-Work-PNC` | `pyWorkPage.pyWorkIDPrefix = "CLM-"` |
| `Endorsment Fac In\When\IsClaim.xml` | **`ASM-FW-GISFW-Work-Claim`** | `pyWorkPage.ClaimData has a value` |

**`IsTravel`**:

| | `pyLogic` | Kondisi |
| --- | --- | --- |
| `NB FacIn\When\IsTravel.xml` | `A` | `pyWorkPage.Quotation.BusinessType = "Travel"` |
| `Endorsment Fac In\When\IsTravel.xml` | **`A OR B`** | `.Quotation.BusinessType = "Travel"` **OR** `pyWorkPage.OfferFacIn.QuotationData.BusinessType = "Travel"` |

📌 `IsTravel` **tidak** memakai P2 — ia membaca properti, bukan hasil query. P2 terbatas pada 6 rule COB.

---

## 5. ⚠️ Temuan privasi — tag yang memuat nama orang dan alamat email

`[terverifikasi]` Saat membedah pohon langkah, dua tag terbukti memuat **nama orang dan alamat
email** di dalam nilainya:

| Tag | Isi yang terbawa |
| --- | --- |
| `pyStepsRepeatDef` | nama operator + alamat email, tertanam di string `Embed-ActivityRepeat…` |
| `pyStepsPreCondParams` | idem, di string `Embed-ActivityPreConditions…` |

Nilainya **tidak disalin** ke dokumen ini. ⛔ Konsekuensi mengikat, memperluas **K-025**: kedua tag
ini **wajib dibuang** dari fixture, tiket, spec, test, maupun lampiran apa pun — sama seperti kolom
`NAMA`/`LOGIN`. Tag `pxCreateOperator` / `pxCreateOpName` sudah tertangani (ada di daftar 23-tag
volatile), **kedua tag ini belum**.

Diajukan sebagai **E-Q15**.

```powershell
# memastikan tidak ada nama/email lolos ke OUTPUT
Select-String -Path "D:\migrasi\RNM\OUTPUT\07-edm\*.md" -Pattern 'Embed-ActivityRepeat|Embed-ActivityPreConditions|@\w+\.co\.id'
```

---

## 6. Pertanyaan terbuka baru dari E-2

| # | Pertanyaan | Pemilik | Dampak |
| ---: | --- | --- | --- |
| **E-Q13** | **`IsNotEDM` bukan negasi `IsEDM`** — `IsEDM` membaca `OfferFacIn.QuotationData.StatusBusiness`, `IsNotEDM` membaca `Quotation.StatusBusiness`. Keduanya bisa benar/salah bersamaan. Disengaja? | work owner + IT | Bila diimplementasikan sebagai `!IsEDM`, perilaku berubah |
| **E-Q14** | **`EDMPremiMenjadi` punya 11 bentuk perhitungan berbeda** di 7 berkas — jauh melampaui catatan arsip (#5). Perlu tabel rumus per (lini × `EdmType` × jalur), bukan satu fungsi | Aktuaria + work owner | Menentukan bentuk modul premi endorsement |
| **E-Q15** | **`pyStepsRepeatDef` dan `pyStepsPreCondParams` memuat nama orang + email** — perluasan K-025 ke dua tag ini | work owner | Kepatuhan privasi pada fixture/tiket/spec |
| **E-Q16** | **`pyStepsPreCondition=false` berskala 617 kemunculan di 173 berkas** (arsip #24 hanya mencatat 21). Bila artinya "langkah dinonaktifkan", sebagian besar alur terpetakan tidak berjalan | **IT** | Menentukan validitas seluruh pemetaan alur EDM |
| **E-Q17** | Hanya `SetOLDValueToEDMWork_FIRE` yang menyetel `IsProRate="Prorate"`; `_LIFE` hanya punya 4 penanda `IsOldData` vs 8–14 di varian lain. Disengaja? | work owner | Kelengkapan lapis C per lini bisnis |

---

## 7. Belum dikerjakan (bukan untuk prompt ini)

- **E-3** — before/after image dan `SetOldData` mendalam (K-039). **Tidak dimulai tanpa perintah.**
- Arsip §4 (gerbang akseptasi EDM, 7 perbedaan), §5 (nilai literal tipe endorsement), §6 (jalur
  retro), §7 (konversi produksi), §8 (diagram), §9 (rule eksklusif) — **belum dikonfirmasi ulang**
- Butir §10 **#12–#23** — belum ditinjau di putaran ini
- Isi 16 berkas tipe yang 100 % berbeda (E-1 §2), termasuk `DecisionTable\IsUWAccepted` dan
  `DecisionTree\Tree_ShortPeriod`

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
