# 04-aturan / 02 — Formula, Enumerasi Status, dan Token Routing

**Sumber tunggal:** korpus ekspor rule Pega di `D:\migrasi\RNM\` — `NB FacIn\`, `RNW Fac In\`,
`Endorsment Fac In\` (READ-ONLY). Tidak ada sumber lain yang dipakai.

**Label bukti yang dipakai di seluruh dokumen:**

| Label | Arti |
| --- | --- |
| `[terverifikasi]` | Ekspresi/tag dikutip langsung dari berkas yang disebut. |
| `[dugaan]` | Disimpulkan dari pola/aritmetika, belum ada pernyataan eksplisit di korpus. |
| `[pertanyaan terbuka]` | Tidak dapat dijawab dari korpus. |

**Aturan yang ditaati:** nilai nama orang tidak pernah disalin; guard berbasis identitas dilaporkan
sebagai jumlah + mekanisme saja. Nama rule tidak pernah dipakai sebagai bukti perilaku — setiap
klaim perilaku menyertakan isi langkah yang terbaca.

---

## 0. Ukuran korpus dan perintah audit dasar

| Angka | Nilai | Perintah audit |
| --- | --- | --- |
| Berkas `.xml` total (3 korpus) | **6.071** | `Get-ChildItem -Path "D:\migrasi\RNM" -Recurse -File -Filter *.xml \| Measure-Object` |
| NB FacIn | **2.083** | `(Get-ChildItem -Path "D:\migrasi\RNM\NB FacIn" -Recurse -File -Filter *.xml \| Measure-Object).Count` |
| RNW Fac In | **1.927** | idem, ganti folder |
| Endorsment Fac In | **2.061** | idem, ganti folder |
| Rule `When` (nama unik lintas 3 korpus) | **226** | `Get-ChildItem "D:\migrasi\RNM\NB FacIn\When","D:\migrasi\RNM\RNW Fac In\When","D:\migrasi\RNM\Endorsment Fac In\When" -File \| Select-Object -ExpandProperty Name \| Sort-Object -Unique \| Measure-Object` |
| Rule `Activity` | **1.747** | `Get-ChildItem "D:\migrasi\RNM" -Recurse -File -Filter *.xml \| Group-Object { (Split-Path $_.DirectoryName -Leaf) }` |
| Rule `DataTransform` | **443** | idem |
| Rule `RDBList` | **605** | idem |
| `DecisionTable` | **32** | idem |
| `DecisionTree` | **3** | idem |

Tiga korpus sebagian besar berisi **salinan rule yang sama**. Contoh terverifikasi:

```powershell
foreach ($c in @("NB FacIn","RNW Fac In","Endorsment Fac In")) {
  (Get-FileHash "D:\migrasi\RNM\$c\DataTransform\HitungPremi_FacInDT.xml" -Algorithm MD5).Hash }
# NB FacIn        431C2599EB1775B28AD79EF54BFCA9AD
# RNW Fac In      431C2599EB1775B28AD79EF54BFCA9AD
# Endorsment      95388978D4F823B088B841866447E12D   <- beda hash, TAPI isi langkah identik (metadata saja yang beda)
```

`[terverifikasi]` Isi langkah `HitungPremi_FacInDT` di korpus Endorsment **identik** dengan NB/RNW
(dibandingkan langkah demi langkah); perbedaan hash berasal dari metadata `pxUpdateDateTime`/
`pzChecksum`, bukan logika.

---

## 1. Rumus premi

### 1.1 `HitungPremi_FacInDT` — DataTransform premi per-coverage (FacIn)

**Berkas:** `NB FacIn\DataTransform\HitungPremi_FacInDT.xml`
(identik di `RNW Fac In\DataTransform\HitungPremi_FacInDT.xml` dan
`Endorsment Fac In\DataTransform\HitungPremi_FacInDT.xml`)
**Kelas:** `ASM-FW-GISFW-Data-Coverage` · **Parameter:** `typediskon` (STRING, IN)

Struktur langkah lengkap `[terverifikasi]`:

| Step | Aksi | Target / kondisi |
| --- | --- | --- |
| 1 | SET | `Param.TotalPremi = 0` |
| 2 | WHEN | `pyWorkPage.Quotation.StatusBusiness=='3'` |
| 2.1 | SET | `.Premium` = (lihat di bawah) |
| 2.2 | SET | `.PremiRp` = (lihat di bawah) |
| 3 | WHEN | `pyWorkPage.Quotation.StatusBusiness !='3'` |
| 3.1 | SET | `.Premium` |
| 3.2 | SET | `.PremiRp` |
| 4 | FOR_EACH_PAGE_IN | `pyWorkPage.OfferFacIn.LocationList` |
| 4.1 | FOR_EACH_PAGE_IN | `.Property.CoverageList` |
| 4.1.1 | SET | `Param.TotalPremi = Param.TotalPremi + .Premium` |
| 5 | SET | `pyWorkPage.Policy.Payment.Premium` (nilai kosong di ekspor) |
| 6 | APPLY_MODEL | `ChangeDiskon_FacIn` |

**Cabang non-endorsement (`StatusBusiness != '3'`)** `[terverifikasi]`:

```
Step 3.1  .Premium =
  @Math.divide((.TSI*.Rate*pyWorkPage.OfferFacIn.ProRatePercent
               *@if(.FirstLossScale=="",100,.FirstLossScale)
               *@if(.IndemnityPercentage=="",100,.IndemnityPercentage)),1000000000,4)

Step 3.2  .PremiRp =
  @Math.divide((.TSI*.Rate*pyWorkPage.OfferFacIn.ProRatePercent
               *.FirstLossScale
               *@if(.IndemnityPercentage=="",100,.IndemnityPercentage)),1000000000,4)
  *@if(pyWorkPage.Policy.CurrencyValue=="",1,pyWorkPage.Policy.CurrencyValue)
```

**Cabang endorsement (`StatusBusiness == '3'`)** `[terverifikasi]`:

```
Step 2.1  .Premium =
  ((@Math.divide((.TSI*.Rate*.ProratePercentEDMEnd
                 *@if(.FirstLossScale=="",100,.FirstLossScale)
                 *@if(.IndemnityPercentage=="",100,.IndemnityPercentage)),1000000000,4))
  +(@Math.divide((.OldCoverage(1).TSI*.OldCoverage(1).Rate
                 *@if(.ProratePercentStartEDM=="",100,.ProratePercentStartEDM)
                 *@if(.OldCoverage(1).FirstLossScale=="",100,.OldCoverage(1).FirstLossScale)
                 *@if(.OldCoverage(1).IndemnityPercentage=="",100,.OldCoverage(1).IndemnityPercentage)),1000000000,4)))

Step 2.2  .PremiRp =
  ((@Math.divide((.TSI*.Rate*.ProratePercentStartEDM
                 *@if(.FirstLossScale=="",100,.FirstLossScale)
                 *@if(.IndemnityPercentage=="",100,.IndemnityPercentage)),1000000000,4))
  +(@Math.divide((.OldCoverage(1).TSI*.OldCoverage(1).Rate
                 *@if(.ProratePercentEDMEnd=="",100,.ProratePercentEDMEnd)
                 *@if(.OldCoverage(1).FirstLossScale=="",100,.OldCoverage(1).FirstLossScale)
                 *@if(.OldCoverage(1).IndemnityPercentage=="",100,.OldCoverage(1).IndemnityPercentage)),1000000000,4)))
  *@if(pyWorkPage.Policy.CurrencyValue=="",1,pyWorkPage.Policy.CurrencyValue)
```

#### Tiga asimetri yang material untuk rekonsiliasi

1. **`FirstLossScale` dijaga di `.Premium` tetapi TIDAK di `.PremiRp`** `[terverifikasi]`.
   Step 3.1 memakai `@if(.FirstLossScale=="",100,.FirstLossScale)`; step 3.2 memakai `.FirstLossScale`
   telanjang. Bila `FirstLossScale` kosong, `.Premium` memakai 100 sedangkan `.PremiRp` memakai
   nilai kosong → hasil dua properti berbeda drastis.
2. **Properti prorate tertukar antara `.Premium` dan `.PremiRp` di cabang endorsement**
   `[terverifikasi]`. Step 2.1 memakai `.ProratePercentEDMEnd` untuk coverage baru dan
   `.ProratePercentStartEDM` untuk `OldCoverage(1)`; step 2.2 memakai **pasangan terbalik**
   (`.ProratePercentStartEDM` untuk coverage baru, `.ProratePercentEDMEnd` untuk `OldCoverage(1)`).
   Apakah ini disengaja: `[pertanyaan terbuka]`.
3. **Guard `@if(...=="",100,...)` tidak konsisten pada cabang 2.1**: `.ProratePercentEDMEnd` (coverage
   baru) tidak dijaga, sedangkan `.ProratePercentStartEDM` (OldCoverage) dijaga. Di step 2.2 pola
   terbalik `[terverifikasi]`.

#### Urutan pembulatan

`[terverifikasi]` Pembulatan **4 desimal dilakukan di dalam setiap `@Math.divide`**, lalu hasilnya
dijumlahkan (step 2.1/2.2) dan dikalikan `CurrencyValue` (step 2.2/3.2) **setelah** pembulatan.
Konsekuensi: `.PremiRp = round(premi,4) × kurs` — bukan `round(premi × kurs, 4)`. Urutan ini harus
ditiru persis; membaliknya menghasilkan selisih sen pada nilai besar.

`[terverifikasi]` `Param.TotalPremi` (step 4.1.1) adalah penjumlahan `.Premium` yang **sudah**
dibulatkan 4 desimal per coverage — akumulasi galat pembulatan per coverage, bukan pembulatan di
akhir.

---

### 1.2 `CountPremi_ACT` — mesin premi utama per `CoverageBasis`

**Berkas:** `NB FacIn\Activity\CountPremi_ACT.xml` (534.732 byte; identik di RNW; Endorsment
534.727 byte). Berkas ini memuat **54** kemunculan konstanta `1000000000`:

```powershell
$c=[IO.File]::ReadAllText("D:\migrasi\RNM\NB FacIn\Activity\CountPremi_ACT.xml")
([regex]::Matches($c,'1000000000')).Count      # -> 54
```

Prorate dihitung lebih dulu `[terverifikasi]` (step `RH_4.pySteps(4)` dan `(8)`):

```
Local.DiffYear     = @DateTime.DateTimeDifference(Start, End, Y)
Local.TotalDayDiff = @DateTime.DateTimeDifference(addCalendar(Start, DiffYear), End, D)
Local.TotalDays    = @DateTime.DateTimeDifference(Local.Year, Local.LastPeriode, D)
Local.Prorate      = (Local.DiffYear*100)
Local.Prorate      = Local.Prorate + (@Math.divide(Local.TotalDayDiff,Local.TotalDays,20)*100)
```

`[terverifikasi]` `Local.TotalDays` dihitung dari 1 Januari tahun-akhir sampai 1 Januari tahun
berikutnya (`Local.Year` = `@year(EndDateTime)+"0101T050000.000 GMT"`), jadi memperhitungkan tahun
kabisat (365/366). **Kontras** dengan §1.5 dan §2.2 yang memakai konstanta 365 keras.

`[terverifikasi]` Timestamp dinormalisasi ke `T050000.000 GMT` (= 12:00 WIB) bila kosong
(step `RH_4.pySteps(3)`), memakai zona `Asia/Jakarta` dan locale `in_ID`.

#### Matriks rumus premi (semua `[terverifikasi]`, dikutip dari step-step di bawah `RH_4`)

Semua memakai presisi **20 desimal**.

| CoverageBasis | PctAdjustment | NetRate | Pembilang | Pembagi |
| --- | --- | --- | --- | --- |
| 1 | 0 | — | `TSI*Rate*Prorate*Indemnity*LossLimit` | `1000000000` (1e9) |
| 1 | ≠0 | — | `TSI*Rate*Prorate*Indemnity*PctAdjustment*LossLimit` | `@toDecimal("100000000000")` (1e11) |
| 1 | 0 | ada | `TSI*NetRate*Prorate*Indemnity*LossLimit` | `1000000000` |
| 1 | ≠0 | ada | `TSI*NetRate*Prorate*Indemnity*PctAdjustment*LossLimit` | `@toDecimal("100000000000")` |
| 2 | 0 | — | `TSI*Rate*FirstScale*Prorate*Indemnity*LossLimit` | `@toDecimal("100000000000")` |
| 2 | ≠0 | — | `TSI*Rate*FirstScale*Prorate*Indemnity*PctAdjustment*LossLimit` | `@toDecimal("10000000000000")` (1e13) |
| 2 | 0 | ada | `TSI*NetRate*FirstScale*Prorate*Indemnity*LossLimit` | `@toDecimal("100000000000")` |
| 2 | ≠0 | ada | `TSI*NetRate*FirstScale*Prorate*Indemnity*PctAdjustment*LossLimit` | `@toDecimal("10000000000000")` |
| 3 | 0 | — | `TSI*Rate*Prorate*Indemnity*LossLimit` | `@toDecimal("1000000000")` |
| 3 | ≠0 | — | `+ PctAdjustment` | `@toDecimal("100000000000")` |
| 4 atau 5 | 0 | — | `TSI*Rate*Prorate*Indemnity*LossLimit` | `@toDecimal("1000000000")` |
| 4 atau 5 | ≠0 | — | `+ PctAdjustment` | `@toDecimal("100000000000")` |

**Rumus balik (`PremiStatus=="amount"`)** `[terverifikasi]` — rate diturunkan dari premi:

```
.Rate = @Math.divide((.Premium * 1000000000),
                     (.TSI * Local.Prorate * .IndemnityPercentage * Local.LossLimit), 20)
```

**TSILiability per CoverageBasis** `[terverifikasi]` — presisi **4 desimal** (berbeda dari premi
yang 20 desimal):

| CoverageBasis | Rumus |
| --- | --- |
| 1 | tidak diset |
| 2 | `.TSILiability = @Math.divide((.TSI * .FirstLoss),100,4)` |
| 3 | `.TSILiability = @Math.divide((.TSI * .EmlPml),100,4)` |
| 4, 5 | `.TSILiability = @Math.divide((.TSI * Local.SubLimit),100,4)` dan di blok lain `@Math.divide((.TSI * .Sublimit),100,4)` |

`[terverifikasi]` Dua ejaan properti dipakai di berkas yang sama: `Local.SubLimit` (diisi dari
`.Sublimit` di step `RH_4.pySteps(10)`) dan `.Sublimit` langsung. `[pertanyaan terbuka]` apakah
keduanya selalu bernilai sama.

**Diskon** `[terverifikasi]`:

```
.Discount           = @Math.divide((.Premium*.DiscountPercentage),100,20)
.DiscountPercentage = @Math.divide(.Discount,.Premium,20)*100
```

**Agregasi** `[terverifikasi]`: `Local.TotalGrossPremi` dijumlahkan dengan cabang terpisah untuk
`.CoverageBasis!=5` dan `.CoverageBasis==5` (layering) lalu disimpan ke
`...PropertyItemList(...).TotalGrossPremi`.

#### Interpretasi skala pembagi

`[dugaan]` Pembagi konsisten dengan tiap faktor persen menyumbang ×100 dan `Rate` menyumbang ×1000:

```
1e9  = 1000 (Rate ‰) × 100 (Prorate %) × 100 (Indemnity %) × 100 (LossLimit %)
1e11 = 1e9 × 100 (PctAdjustment %)        1e13 = 1e11 × 100 (FirstScale %)
```

Artinya `Rate` diperlakukan **per mille (‰)**, bukan persen. Tidak ada pernyataan eksplisit di
korpus — `[pertanyaan terbuka]` untuk konfirmasi bisnis. Lihat §1.5 yang **bertentangan**.

---

### 1.3 `CountPremiFacIn_act` — premi + ASM share + agregasi per mata uang

**Berkas:** `NB FacIn\Activity\CountPremiFacIn_act.xml` (138.784 byte)

`[terverifikasi]` Step `RH_1.pySteps(3)`:

```
.Premium    = @Math.divide((.TSI*.Rate*pyWorkPage.OfferFacIn.ProRatePercent
                           *@if(.FirstLossScale=="",100,.FirstLossScale)
                           *@if(.IndemnityPercentage=="",100,.IndemnityPercentage)),1000000000,4)
.ASMPremium = @Math.divide(.Premium * pyWorkPage.OfferFacIn.ASMSharePersen,100,4)
Local.lPremi = 0
```

**Pembulatan ganda `[terverifikasi]`:** `.Premium` sudah dibulatkan 4 desimal, lalu `.ASMPremium`
membulatkan lagi 4 desimal dari hasil yang sudah dibulatkan. Bukan `round(TSI*Rate*.../1e9 × share,
4)`.

`[terverifikasi]` Agregasi per mata uang (step `RH_1.pySteps(4)`):

```
loop CurrencyMaster.pxResults:
    Local.Currency = .CURR ;  Local.lPremi = 0 ;  Local.lDiscount = 0
    loop LocationList → PropertyItemList → CoverageList:
        when Local.Currency==Local.ObjectCurrency:
            Local.lPremi    = .Premium  + Local.lPremi
            Local.lDiscount = .Discount + Local.lDiscount
    loop OfferFacIn.CurrencyList:
        when .Name==Local.Currency:
            .Policy.Payment.Premium = Local.lPremi
            .Policy.Payment.Diskon  = Local.lDiscount
```

`[terverifikasi]` Step `RH_1.pySteps(5)` menetapkan
`pyWorkPage.Policy.Payment.Premium = Local.lPremi` **di luar** loop mata uang — jadi nilai yang
tersisa adalah nilai mata uang **terakhir** dari `CurrencyMaster.pxResults`. `[pertanyaan terbuka]`
apakah ini disengaja atau bug.

`[terverifikasi]` Validasi rate (step 7 & 8):

```
when (.Coverage==100815||.Coverage==100825||.Coverage==100837)&&(.Rate<MinMax.CARI1||.Rate>MinMax.CARI2)
     Local.ERRMSG = "Policy rate for this coverage must be between " + MinMax.CARI1 + " and " + MinMax.CARI2
```

Kode coverage literal **100815, 100825, 100837** — artinya `belum terverifikasi`.

---

### 1.4 `CountPremiNusareRetro_Act` — premi bagian Nusantara Re atas retro

**Berkas:** `NB FacIn\Activity\CountPremiNusareRetro_Act.xml` (64.707 byte; identik di RNW dan
Endorsment)

`[terverifikasi]` Konstanta share ditanam keras di step `RH_1.pySteps(1)`:

```
local.sharernm1 = 0.45
local.sharernm2 = 0.05
local.tampungpreminet1 = 0
```

`[terverifikasi]` Banding TSI dilakukan **sebagai string** (step `2.1.1` dan `2.1.2`):

```
when .TSILiability<="3000000000"   ("jika qs1")
     .PremiLifeNusantaraRe = .GrossPremiumRetro * local.sharernm1

when .TSILiability>"3000000000"    ("jika qs2, surplus, surplus 2")
     local.tampungpreminet1 = .GrossPremiumRetro * local.sharernm1
     .PremiLifeNusantaraRe  = (.GrossPremiumRetro * local.sharernm2) + local.tampungpreminet1
```

`[terverifikasi]` Step `2.1.3`:

```
.PremiumRetro         = @toDecimal(.GrossPremiumRetro)-@toDecimal(.PremiLifeNusantaraRe)
.PremiLifeNusantaraRe = .Premium-.PremiumRetro
```

Perhatikan: `.PremiLifeNusantaraRe` ditulis dua kali dalam satu step — nilai dari step 2.1.1/2.1.2
dipakai untuk menurunkan `.PremiumRetro`, lalu ditimpa dari `.Premium`.

`[terverifikasi]` Step `2.1.4` — **dua prakondisi terpisah**, keduanya `WhenTrue=2`/`WhenFalse=3`:

```
pyStepsPreCondParamsWhen = pyWorkPage.OfferFacIn.ProRateType!=3
pyStepsPreCondParamsWhen = .pxListSubscript!=1
    .PremiLifeNusantaraRe = 0
```

Arti kode transisi Pega `2` dan `3` di `pyStepsPreCondParamsWhenTrue/False`: **belum terverifikasi**.
Efek: bila `ProRateType != 3`, premi bagian Nusantara Re **dinihilkan**. Ini kandidat perbaikan,
bukan keputusan migrasi.

Ekspresi terkait di berkas yang sama `[terverifikasi]`:
`If(pyWorkPage.OfferFacIn.ProRateType==3,.Premium-.PremiumRetro,0)`

---

### 1.5 `CountProrateExtension_Act` — premi dengan pembagi 1e4 (bertentangan dengan §1.2)

**Berkas:** `NB FacIn\Activity\CountProrateExtension_Act.xml` (75.306 byte)

`[terverifikasi]`:

```
Step 1   Local.dateDifferent = @DateTime.DateTimeDifference(Start, End, D)
         Local.day = 365
Step 2   when param.metodeKalkulasi==1
         .ProRatePercent = @if((@Math.divide(local.dateDifferent,Local.day,8)*100)>100,
                               100,
                               (@Math.divide(local.dateDifferent,Local.day,8)*100))
Step 3   when param.metodeKalkulasi==2  ->  Property-Map-DecisionTree (Tree_ShortPeriod)
Step 5   when param.metodeKalkulasi==3  ->  .ProRatePercent = 100 ; Local.dateDifferent = 100 ; Local.prorate = 100
Step 7   .Premium = @Math.divide(.TSI *.Rate*Local.prorate,10000,4)
```

**Temuan material `[terverifikasi]`:** pembagi di sini **10000** (1e4), presisi **4**, sementara
`CountPremi_ACT` memakai 1e9/1e11/1e13 dengan presisi 20 dan `HitungPremi_FacInDT` memakai 1e9
dengan presisi 4. Skala 1e4 = 100 (Rate %) × 100 (Prorate %) — konsisten hanya bila `Rate`
diperlakukan **persen**, bertentangan dengan interpretasi per-mille di §1.2. Salah satu dari dua
rumus ini salah satu ordo besaran 10, atau `.Rate` di dua konteks itu adalah satuan berbeda.
**`[pertanyaan terbuka]` — MEMBLOKIR.**

`[terverifikasi]` Enumerasi `metodeKalkulasi` / `CalculateMethod`:

| Nilai | Deskripsi step di korpus | Status arti |
| --- | --- | --- |
| `1` | step 6 `pyStepsDescription : Prorate` | `[terverifikasi]` prorata hari/365 |
| `2` | step 3 `pyStepsDescription : shortperiode` | `[terverifikasi]` memakai `Tree_ShortPeriod` |
| `3` | step 5 `pyStepsDescription : fixrate` | `[terverifikasi]` ProRatePercent dipaksa 100 |
| `0`, `≥4` | tidak pernah muncul | celah |

Audit: `& lit.ps1 -Pattern '((?:metodeKalkulasi|CalculateMethod|CalculatedMethod)\s*[!=]=*\s*(?:&quot;|")?\d+)'`
→ hanya `1`, `2`, `3`.

---

### 1.6 `Tree_ShortPeriod` — tabel short-period (DecisionTree, **terekspor penuh**)

**Berkas:** `NB FacIn\DecisionTree\Tree_ShortPeriod.xml` (265.195 byte; **hash/ukuran identik** di
ketiga korpus). Kelas `ASM-FW-GISFW-Data-Coverage`. Label rule: `Master Untuk Short Period`.

`[terverifikasi]` Dua cabang tingkat atas berdasarkan `pyWorkPage.Quotation.GroupPanel`:

**Cabang `GroupPanel = "007"` — hasil dalam skala PERSEN:**

| Kondisi | Hasil |
| --- | --- |
| `dateDifferent < 8` | `12.5` |
| `<= 30` | `15` |
| `<= 60` | `20` |
| `<= 90` | `30` |
| `<= 120` | `40` |
| `<= 150` | `50` |
| `<= 180` | `60` |
| `<= 210` | `70` |
| `<= 240` | `80` |
| `> 240` | `100` |

**Cabang `GroupPanel = "002"` — hasil dalam skala PECAHAN:**

| Kondisi | Hasil |
| --- | --- |
| `dateDifferent < 8` | `0.125` |
| `= 45` *(operator kesamaan, bukan `<=`)* | `0.2` |
| `> 45` … `<= 75` | `0.3` |
| `<= 105` | `0.4` |
| `<= 135` | `0.5` |
| `<= 165` | `0.6` |
| `<= 195` | `0.7` |
| `<= 225` | `0.75` |
| `<= 255` | `0.8` |
| `<= 285` | `0.85` |
| `<= 315` | `0.9` |
| `<= 345` | `0.95` |
| `> 345` | `1` |

**Empat temuan material:**

1. **Skala hasil berbeda 100×** antara cabang `"007"` (12,5…100) dan `"002"` (0,125…1)
   `[terverifikasi]`. Pemanggil di `CountProrateExtension_Act` step 7 memakai hasil sebagai
   `Local.prorate` dengan pembagi 10000 yang sama untuk kedua cabang → salah satu cabang meleset
   faktor 100. **`[pertanyaan terbuka]` — MEMBLOKIR.**
2. **Pita hari berbeda**: `"007"` memakai kelipatan 30 (8/30/60/90/120/150/180/210/240);
   `"002"` memakai kelipatan 30 yang bergeser (8/45/75/105/135/165/195/225/255/285/315/345)
   `[terverifikasi]`.
3. **Baris `param.dateDifferent = 45`** memakai operator kesamaan, bukan `<=`. Nilai 8..44 di cabang
   `"002"` karenanya jatuh melalui ke cabang `> 45` yang `false` → tidak ada hasil.
   `[pertanyaan terbuka]` apa hasil untuk `dateDifferent` antara 8 dan 44 di cabang `"002"`.
4. **Nilai `GroupPanel` selain `"002"` dan `"007"` tidak punya cabang** `[terverifikasi]` — tidak ada
   node `otherwise`. Audit: `& lit.ps1 -Pattern 'GroupPanel\s*[!=]=*\s*(?:&quot;|")(\d*)'` → hanya
   `002` (3 berkas) dan `007` (8 berkas). `[pertanyaan terbuka]` nilai balik default.

Perintah audit hasil tree:

```powershell
$c = Get-Content "D:\migrasi\RNM\NB FacIn\DecisionTree\Tree_ShortPeriod.xml" -Raw -Encoding UTF8
[regex]::Matches($c,'<(pyExpressionString|pyAction|pyResult)>(.*?)</\1>','Singleline') |
  ForEach-Object { $_.Groups[1].Value + ' :: ' + $_.Groups[2].Value }
```

---

### 1.7 Sebaran seluruh pembagi + presisi di korpus

Skrip audit (disimpan sebagai `div.ps1`):

```powershell
$root = "D:\migrasi\RNM"; $tally = @{}
$rx = [regex]'(?i)divide\(.{0,500}?,\s*@?(?:toDecimal\()?"?(\d{2,})"?\)?\s*(?:,\s*(\d+)\s*)?\)'
foreach ($corp in @("NB FacIn","RNW Fac In","Endorsment Fac In")) {
  foreach ($f in [IO.Directory]::EnumerateFiles("$root\$corp","*.xml",'AllDirectories')) {
    $c = [IO.File]::ReadAllText($f)
    if ($c.IndexOf('divide(') -lt 0 -and $c.IndexOf('Divide(') -lt 0) { continue }
    foreach ($m in $rx.Matches($c)) {
      $k = "divisor=" + $m.Groups[1].Value + "  digits=" + $(if($m.Groups[2].Success){$m.Groups[2].Value}else{"(none)"})
      if (-not $tally.ContainsKey($k)) { $tally[$k] = New-Object System.Collections.Generic.HashSet[string] }
      [void]$tally[$k].Add($f) } } }
$tally.Keys | Sort-Object | ForEach-Object { "{0}`tfiles={1}" -f $_, $tally[$_].Count }
```

Hasil: **44 kombinasi pembagi×presisi berbeda**. Yang relevan uang:

| Pembagi | Presisi | Berkas |
| --- | --- | --- |
| 100 | 20 | 149 |
| 100 | 4 | 140 |
| 100 | (tanpa presisi) | 33 |
| 100 | 0 / 1 / 2 / 6 / 8 / 10 / 16 | 8 / 12 / 5 / 2 / 9 / 1 / 3 |
| 1000 | 20 · 4 · 8 · 6 · — | 16 · 12 · 1 · 1 · 6 |
| 10000 | 20 · 4 | 34 · 19 |
| 100000 | 20 · 4 | 15 · 9 |
| 1000000 | 4 | 13 |
| 10000000 | 20 | 3 |
| 100000000 | 4 | 9 |
| **1000000000** | **20 · 4** | **6 · 14** |
| **100000000000** | **20** | **3** |
| **10000000000000** | **20** | **3** |
| 365 | 20 · 8 · 6 | 9 · 9 · 1 |
| 12 | 0/2/4/6/8/12/16 | 3/3/3/3/6/7/7 |

**Temuan:** pembagi `100` muncul dengan **sembilan** presisi berbeda (0,1,2,4,6,8,10,16,20 + tanpa
presisi). Migrasi tidak boleh menyeragamkan — setiap pemanggilan harus ditiru per-rule.
`[terverifikasi]`

---

## 2. Rumus terkait lain

### 2.1 Komisi reasuransi (RI Commission)

**`CountRIComm_Act`** — `NB FacIn\Activity\CountRIComm_Act.xml` (25.199 byte), satu langkah
`[terverifikasi]`:

```
.CoverageList(1).FacOutObjectList(1).commision =
   @divide((@toDecimal(.CoverageList(1).FacOutObjectList(1).ObjectPremi)
           *@toDecimal(.CoverageList(1).FacOutObjectList(1).RiComm)),100,20)
```

Presisi **20**. `@toDecimal(...)` dipakai pada kedua operand → keduanya masuk sebagai **string**.

**`CalculateNetRate_ACT`** — `NB FacIn\Activity\CalculateNetRate_ACT.xml` (60.351 byte), step
`RH_1.pySteps(5)` `[terverifikasi]`:

```
Local.RIComm = @replaceAll(.CoverageList(1).FacOutObjectList(1).RiComm,",",".")
.CoverageList(1).FacOutObjectList(1).ObjectPremi =
   @Math.divide((Local.ShareOffered *Local.Rate * Local.proratepct),100000,20)
.CoverageList(1).FacOutObjectList(1).commision =
   @Math.divide((@Math.divide((Local.ShareOffered *Local.Rate*Local.proratepct),100000,20) * Local.RIComm),100,20)
```

**Pembulatan bersarang `[terverifikasi]`:** `commision` menghitung ulang `ObjectPremi` di dalam
dirinya, membulatkan ke 20 desimal, lalu membulatkan lagi ke 20. Pembagi di sini **100000** (1e5) —
bukan 1e9 maupun 1e4.

### 2.2 Share offered / percent offered

**`CalculateNetRate_ACT`** step 1, 3, 4, 8 `[terverifikasi]`:

```
Local.proratepct = @String.toDate(FacRetroList(i).EndPeriod) - @String.toDate(FacRetroList(i).StartPeriod)
Local.proratepct = Local.proratepct/365*100          <- 365 keras, tanpa pembulatan eksplisit
when IsEdmExtendPeriod:  Local.proratepct = 100

when param.Type=="percent":
    .ShareOffered   = @Math.divide((.PercentOffered * .FacOutTSI),100)        <- TANPA argumen presisi
when param.Type=="amount":
    Local.ShareOffered = @string.toDecimal(.ShareOffered)
    .PercentOffered = @Math.divide((Local.ShareOffered),.FacOutTSI)*100       <- TANPA argumen presisi

.PercentFacOut        = Local.PercentFacOut
.NominalShareOffered  = @Math.divide((.PercentFacOut * .FacOutTSI),100)       <- TANPA argumen presisi
when Local.PercentFacOut<100 || Local.PercentFacOut>100:
    InputParam.HASIL12 = "Total (%)Share Offered can't be lower or more than 100"
```

`[pertanyaan terbuka]` Presisi default `@Math.divide` dua argumen di Pega — **MEMBLOKIR** untuk
rekonsiliasi angka.

**`CountASMShareTotal_ACT`** — `NB FacIn\Activity\CountASMShareTotal_ACT.xml` (29.773 byte)
`[terverifikasi]`:

```
when @String.isDouble(.PercentShare):
    .ASMShareTotal = ((.PercentShare/100)*.SumTotalTSI)*(.ASMSharePersen/(@String.toDecimal("100.0000")))
```

Konstanta `100` ditulis sebagai string `"100.0000"` lalu dikonversi. Prakondisi
`@String.isDouble(.PercentShare)` membuktikan `PercentShare` disimpan sebagai **string**.

### 2.3 Share Nusantara Re atas coverage

**`HitungPremiAndTSINusantaraRe_ACT`** — `NB FacIn\Activity\HitungPremiAndTSINusantaraRe_ACT.xml`
(102.637 byte), step `RH_1.pySteps(1)` `[terverifikasi]`:

```
.TSINusantaraRe  = @Math.divide((@toDecimal(.PctShareNusantaraRe)*.TSILiability),100,20)
.PremiNusantaraRe= @Math.divide((@toDecimal(.PctShareNusantaraRe)*.Premium),100,20)
```

`.TSINusantaraRe` dihitung dari **`.TSILiability`** (hasil §1.2, presisi 4), bukan dari `.TSI`.
Jadi ada rantai pembulatan 4 → 20.

### 2.4 Spreading (penyebaran ke treaty)

**`SetTSIPremiSpreaded_FacIn`** — `NB FacIn\DataTransform\SetTSIPremiSpreaded_FacIn.xml`
(368.903 byte; byte-identik di ketiga korpus).

`[terverifikasi]` Step 1: `WHEN IsGroup → EXIT_MODEL`. Seluruh transform **dilewati** bila `IsGroup`
benar (lihat §6 untuk isi `IsGroup`).

`[terverifikasi]` Inti per cabang lini bisnis (`IsPA` step 5, `IsMBU` 6, `IsFire` 7,
`isGolfInsurance` 8, `IsEngineering` 9, `IsAneka` 10):

```
Param.TSI   = .TSINusantaraRe
Param.Premi = .PremiNusantaraRe
loop .SpreadingList:
    when Param.Type=="percent":
        .TSISpreaded     = @Math.divide(.SharePercentage*Param.TSI,100,20)
        .PremiumSpreaded = @Math.divide(.SharePercentage*Param.Premi,100,20)
    when Param.Type=="amount":
        .SharePercentage = @Math.divide(.TSISpreaded,Param.TSI,20)*100
        .PremiumSpreaded = @Math.divide(.SharePercentage*Param.Premi,100,20)
    Param.Spreaded = Param.Spreaded + .PremiumSpreaded
    Param.Total    = @notEqual(Param.Spreaded,Param.Premi)
    pyWorkPage.Test= Param.Total
```

**Tiga inkonsistensi antar-lini yang material `[terverifikasi]`:**

1. **Pembulatan ulang ke 4 desimal hanya di sebagian cabang.** Setelah loop, langkah
   `.TSISpreaded = @Math.divide(.SharePercentage*Param.TSI,100,4)` dan
   `.PremiumSpreaded = @Math.divide(.SharePercentage*Param.Premi,100,4)` ada di cabang **PA**
   (5.1.1.4.3/4.4), **MBU** (6.1.1.3.4/3.5), dan **Fire non-basis-5** (7.1.1.1.2.3.4/3.5) —
   **tidak ada** di cabang Fire `CoverageBasis==5` (layering), Golf, Engineering, dan Aneka. Lini
   terakhir itu menyimpan 20 desimal.
2. **Pra-pembulatan sebelum uji kesetaraan hanya di PA.** Step 5.1.1.5/5.1.1.6:
   `Param.Premi = @divide(Param.Premi,1,4)` dan `Param.Spreaded = @divide(Param.Spreaded,1,4)`
   sebelum `@notEqual`. Cabang lain membandingkan nilai 20-desimal secara langsung — praktis selalu
   tidak sama karena galat pembulatan.
3. **Reset `Param.Spreaded = 0` per coverage tidak seragam.** Ada di PA (5.1.1.3) dan Aneka
   (10.1.1.1.1.3); **tidak ada** di MBU, Fire, Golf, Engineering → akumulator tidak direset antar
   coverage di lini tersebut.

`[terverifikasi]` Gerbang khusus endorsement spreading di setiap cabang:

```
when pyWorkPage.OfferFacIn.QuotationData.StatusBusiness==3 && pyWorkPage.OfferFacIn.QuotationData.Type=="4":
     Param.Total = @notEqual(Param.Spreaded,0)
```

Perhatikan `StatusBusiness==3` (angka) berdampingan dengan `Type=="4"` (string) **dalam satu
ekspresi**.

**`hitungspreadingotomatis_act`** — `NB FacIn\Activity\hitungspreadingotomatis_act.xml`
(66.067 byte) `[terverifikasi]`:

```
Step 2  local.bataslayer1 = "3000000000"     <- literal STRING
        local.bataslayer2 = "6000000000"
        local.bataslayer3 = "20000000000"
        local.bataslayer4 = "25000000000"
Step 3  when pyWorkPage.OfferFacIn.ShareRNML<=local.bataslayer1
        SpreadingList.pxResults(1).TreatyType = 1001410 ; SharePercentage = 100
Step 4  when ShareRNML<=bataslayer2 && ShareRNML>bataslayer1
        local.tampungTSI     = ShareRNML - bataslayer1
        local.tampungPercent = @divide(local.tampungTSI,ShareRNML,20)*100
        pxResults(1).TreatyType=1001411 ; SharePercentage = local.tampungPercent
        pxResults(2).TreatyType=1001410 ; SharePercentage = @divide(bataslayer1,ShareRNML,20)*100
Step 5  when ShareRNML<=bataslayer3 && ShareRNML>bataslayer2
        pxResults(1).TreatyType=1001412 ; pxResults(2)=1001411 ; pxResults(3)=1001410
```

Kode `TreatyType` **1001410 / 1001411 / 1001412** — artinya `belum terverifikasi`.
Ambang perbandingan adalah string yang dibandingkan terhadap nilai numerik.

### 2.5 Premi cedant (share cedant)

**`CountPremiCedant_Act`** — `NB FacIn\Activity\CountPremiCedant_Act.xml` (62.468 byte)
`[terverifikasi]`:

```
Step 2.3.1  pyStepsDescription : SHARE RNM     when pyWorkPage.OfferFacIn.ShareCedantType==1
    CedingCedantList(i).CurrencyList(<LAST>).TSI =
        @Math.divide((CedingCedantList(i).ShareCeding * @Math.divide(.TSI*pyWorkPage.OfferFacIn.PercentShare,100,20)),100,20)
    CedingCedantList(i).CurrencyList(<LAST>).Premium =
        @Math.divide((CedingCedantList(i).ShareCeding * .Policy.Payment.Premium),100,20)

Step 2.3.2  pyStepsDescription : GROSS         when pyWorkPage.OfferFacIn.ShareCedantType==0
    ...TSI     = @Math.divide((ShareCeding * .TSI),100,20)
    ...Premium = @Math.divide((ShareCeding * (@Math.divide((.Policy.Payment.Premium*100),pyWorkPage.OfferFacIn.PercentShare,20))),100,20)
```

`[terverifikasi]` `ShareCedantType`: `0` = **GROSS**, `1` = **SHARE RNM** (dari `pyStepsDescription`).

### 2.6 Brokerage fee, deduction, PPN/PPh — net payable

`[terverifikasi]` **Empat konvensi tanda berbeda** untuk rumus "net payable" dalam korpus yang sama:

| Ekspresi | Berkas |
| --- | --- |
| `Premium − Commision − BrokerageFee + PPh + PPN` | `NB FacIn\Activity\CountPaymentInstallment_Act.xml`; `Endorsment Fac In\Activity\FillPaymentInstallment.xml` |
| `Premium − Commision − BrokerageFee − Deduction2 + PPh + PPN` | `NB FacIn\DataTransform\CountPremiEDM_DT.xml` |
| `Premium − Commision − BrokerageFee − PPh − PPN` | `Endorsment Fac In\Activity\FillDateInstallment_EDM.xml`; `Endorsment Fac In\Activity\FillPaymentInstallment.xml` |
| `Premium − Commision − BrokerageFee − PPh + PPN` | `Endorsment Fac In\Activity\FillDateInstallment_EDM.xml` |

Juga varian urutan `... + PPN + PPh` (`CountPremiEDM_DT.xml`) dan penempatan `Deduction2` sebelum
`BrokerageFee` (`SaveFillPaymentInstallment.xml`). Untuk penjumlahan desimal urutan tidak mengubah
hasil, tetapi **tanda PPh/PPN mengubahnya**. `[pertanyaan terbuka]` mana yang benar — **MEMBLOKIR**.

Perintah audit:

```powershell
& scan.ps1 -Pattern 'Premi(um)?-\.?(Policy\.Payment\.)?Commision[^<>]{0,130}' -Pad 0 | Sort-Object -Unique
```

`[terverifikasi]` Deduction lain:

```
.Deduction2 (installment)  = @Math.divide(((Local.Deduction2) * .InstallmentPercentage), 100, 4)
.Deduction2 (installment)  = @Math.divide(((Local.Deduction2ins) * .InstallmentPercentage), 100, 20)   <- presisi berbeda
.Deduction2 (dari pct)     = @Math.divide((pyWorkPage.Policy.Payment.PctDeduction2*.Policy.Payment.Premium),100,20)
.PctDeduction2 (balik)     = @Math.divide(Local.Deduction2, Local.PremiumSpreaded)*100                 <- tanpa presisi
Deduction1 (pajak inklusif)= @if(.TypeTax=="Inclusive",@divide(.Deduction1,@divide(102.2,100,8),8),.Deduction1)
BrokerageFee (2 varian)    = .BrokerageFeeSebenarnya* @divide(2,100,8)
                             .BrokerageFeeSebenarnya* @divide(2.2,100,8)
```

Konstanta pajak `102.2`, `2`, `2.2` ditanam keras — artinya `belum terverifikasi`.

`[terverifikasi]` Brokerage pada spreading:

```
(.PremiumSpreaded * pyWorkPage.OfferFacIn.PolicyData.Payment.PctBrokerageFee/100) * pyWorkPage.OfferFacIn.ProrateEDMEnd
(.PremiumSpreaded * pyWorkPage.OfferFacIn.PolicyData.Payment.PctDeduction2/100)   * pyWorkPage.OfferFacIn.ProrateEDMEnd
```

### 2.7 TSI dan Limit of Liability (LOL)

`[terverifikasi]`:

```
Local.TSIObjectItem -> .TSI                              (CountPremiFacIn_act, step 2.1 / 2.1.1)
.TSILiability = TSI*FirstLoss/100 | TSI*EmlPml/100 | TSI*SubLimit/100   (4 desimal, §1.2)
.LimitofLiability = @toDecimal(.TSI) * @toDecimal(.PctLoL)/100
.PctLoL = @if(.TSI==0,0,@if(Local.PctLoL>0,Local.PctLoL,
             @if(.LimitofLiability>0,@divide(.LimitofLiability*100,.TSI,4),0)))
Rate efektif LOL = @if(.NetRate>0,.NetRate * @divide(.LostLimit,100,1),.Rate*@divide(.LostLimit,100,1))
.TSIAdjustment   = @Math.divide((.TSIObjectItem*.PercentageAdjustment),100,4)
                   @Math.divide(.TSIObjectItem*@if(.PercentageAdjustment=="",100,.PercentageAdjustment),100,4)
```

**Temuan `[terverifikasi]`:** `@divide(.LostLimit,100,1)` membulatkan ke **1 desimal**. LostLimit
30% menjadi 0,3 (aman); 33,33% menjadi 0,3 — kehilangan presisi besar pada rate. Bandingkan dengan
pemakaian `.LostLimit` mentah (skala 100) di rumus premi §1.2.

Validasi `PercentageAdjustment` `[terverifikasi]`:
`"%Adjustment can't be less than 60% or more than 100%"`, dan ekspresi
`.PercentageAdjustment!="100" || .PercentageAdjustment!="75"`, `.PercentageAdjustment=="60"`,
`=="70"` — perbandingan **sebagai string**. Ekspresi `!="100" || !="75"` selalu bernilai benar
(tautologi) — kandidat bug `[terverifikasi]`.

### 2.8 Yang belum dibaca

Berkas berikut terkait perhitungan tetapi **isinya belum dibaca** dalam sesi ini (ukuran disebut
supaya prioritas jelas):

| Berkas | Ukuran (byte) |
| --- | --- |
| `NB FacIn\Activity\SumTSIPremiSpreadedRNM_FIRE_Act.xml` | 996.052 |
| `NB FacIn\Activity\SumTSIPremiSpreadedRNM_ANEKA_Act.xml` | 748.342 |
| `NB FacIn\Activity\CountGrossPremiEDM_Act.xml` | 720.953 |
| `NB FacIn\Activity\CountPremiAndTSINusantaraReEDM_ACT.xml` | 713.415 |
| `NB FacIn\Activity\GetLimitAkseptasi_Act.xml` | 544.882 |
| `NB FacIn\Activity\GetLimitAkseptasi_ActFlow.xml` | 539.687 |
| `NB FacIn\Activity\SumTSIPremiSpreadedRNM_Act.xml` | 516.395 |
| `NB FacIn\Activity\CountGrossPremi_Act.xml` | 463.937 |
| `NB FacIn\Activity\CountTotalTSIPremiNusaRe_Act.xml` | 457.219 |
| `NB FacIn\Activity\CountPaymentEdm_Act.xml` | 450.796 |
| `NB FacIn\Activity\SumTotalTSIPremiGross_Act.xml` | 446.578 |
| `NB FacIn\Activity\CountPremiAndTSIRNMFireMBU_ACT.xml` | 398.462 |
| `NB FacIn\Activity\CountPremiAndTSIRNMFireAnekaGolf_ACT.xml` | 386.951 |
| `NB FacIn\DataTransform\CountPremiEDM_DT.xml` | 196.746 |
| `Endorsment Fac In\Activity\CalculateScorLife_Act.xml` | 1.911.659 |

Audit: `Get-ChildItem "D:\migrasi\RNM" -Recurse -File -Filter *.xml | Where-Object { $_.Name -match "Prem|Hitung|Count|Calc" } | Sort-Object Length -Descending`

---

## 3. Enumerasi status

### 3.1 `Quotation.StatusBusiness` / `OfferFacIn.QuotationData.StatusBusiness`

| Nilai | Rule yang menguji | Arti | Bukti |
| --- | --- | --- | --- |
| `1` | `When\IsNB.xml` → `pyWorkPage.Quotation.StatusBusiness = 1` | New Business | `[terverifikasi]` (nama rule + isi kondisi konsisten) |
| `2` | `When\IsRenewal.xml` → `= 2` | Renewal | `[terverifikasi]` |
| `3` | `When\IsEDM.xml` → `= 3`; `When\IsNotEDM.xml` → `!= 3` | Endorsement | `[terverifikasi]` |
| `0`, `≥4` | tidak pernah muncul | — | **celah** |

Audit: `& lit.ps1 -Pattern '(StatusBusiness\s*[!=]=*\s*(?:&quot;|")?\d+)'` → 20 bentuk distinct,
semua hanya memakai literal 1, 2, 3.

**Dua nama properti untuk hal yang sama** `[terverifikasi]`:
`pyWorkPage.Quotation.StatusBusiness` (dipakai `IsNB`, `IsRenewal`, `IsNotEDM`,
`HitungPremi_FacInDT`) vs `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness` (dipakai `IsEDM`,
seluruh keluarga `IsEdm*`). `[pertanyaan terbuka]` apakah keduanya selalu disinkronkan —
**MEMBLOKIR**, karena `HitungPremi_FacInDT` memilih cabang rumus premi berdasarkan yang **pertama**
sementara gerbang endorsement memakai yang **kedua**.

**Perbandingan campur tipe** `[terverifikasi]`: `StatusBusiness==3` (angka, 100 berkas) dan
`StatusBusiness=="3"` (string, 29 berkas) berdampingan; `HitungPremi_FacInDT` memakai `=='3'`.

### 3.2 `OfferFacIn.QuotationData.Type` — tipe endorsement

**Sumber label (VERIFIED):** `Endorsment Fac In\Activity\SetEdmType.xml`, step `RH_1.pySteps(2..4)`,
mengisi `TypePenambahan/TypePengurangan/TypePerubahan.pxResults(<APPEND>).CARI1` (nilai) dan
`.CARI2` (label).

| Nilai | Label di `SetEdmType.xml` | Rule `When` yang menguji | Status arti |
| --- | --- | --- | --- |
| `0` | `"Adjustment Reff. Number"` | `IsEdmAdjRefNo`, **dan** `IsEDMRiSlip` | `[terverifikasi]` |
| `1` | `"Extend Period"` | `IsEdmExtendPeriod` | `[terverifikasi]` |
| `2` | `"Adjustment TSI / Add Object / Rate / Premium"` | `IsEdmAdjTSI` | `[terverifikasi]` |
| `3` | **tidak pernah ditawarkan UI** | `IsEdmAdjRate` | `belum terverifikasi` |
| `4` | `"Adjustment Spreading"` | `IsEdmAdjSpreading` | `[terverifikasi]` |
| `5` | **tidak pernah ditawarkan UI** | `IsEdmAddObject` | `belum terverifikasi` |
| `6` | `"Adjustment Period"` | `IsEdmAdjPeriod` | `[terverifikasi]` |
| `7` | `"Adjustment Currency"` | `IsEdmAdjCurrency` | `[terverifikasi]` |
| `8` | `"Adjustment Deduction"` | `IsEdmAdjRIC` (`C OR D`) | **bentrok** — lihat peringatan |
| `9` | `"Adjustment Insured Name"` | `IsEdmAdjInsured` | `[terverifikasi]` |
| **`10`** | — | — | **CELAH TOTAL: tidak muncul di mana pun** |
| `11` | `"Adjustment Share Cedant"` | `IsEdmAdjShareCedant` | `[terverifikasi]` |
| `12` | `"Adjustment PPN/PPH"` | `IsEdmPPNPPH`, `IsEdmAdjCeding`, `IsEdmAdjRIC` | **bentrok** |

Audit nilai yang muncul:
`& scan.ps1 -Pattern 'QuotationData\.Type' -Pad 60` → literal 0,1,2,3,4,5,6,7,8,9,11,12. Nilai
**10 tidak pernah muncul**.

#### PERINGATAN — nama rule bukan bukti perilaku (bukti langsung)

`[terverifikasi]` Kondisi lengkap ketiga rule berikut (`pyConditionValue1String` + `pyLogic`):

```
IsEdmPPNPPH     LOGIC=A AND B AND C
                EdmType = 4 | StatusBusiness = 3 | Type = 12
IsEdmAdjCeding  LOGIC=A AND B AND C
                EdmType = 4 | StatusBusiness = 3 | Type = 12       <- IDENTIK dengan IsEdmPPNPPH
IsEdmAdjRIC     LOGIC=A AND B AND (C OR D)
                EdmType = 4 | StatusBusiness = 3 | Type = 12 | Type = 8
```

Tiga rule dengan nama bertema berbeda ("PPN/PPH", "Adj Ceding", "Adj RIC") menguji **nilai yang
sama**, sementara label UI untuk `Type = 12` adalah `"Adjustment PPN/PPH"`. `IsEdmAdjCeding` dan
`IsEdmPPNPPH` **selalu bernilai sama** — tidak ada cara membedakan endorsement "ceding" dari
"PPN/PPH" lewat `Type`. **`[pertanyaan terbuka]` — MEMBLOKIR.**

Bentrok kedua `[terverifikasi]`: `Type = 0` diuji oleh `IsEdmAdjRefNo` (label UI
`"Adjustment Reff. Number"`) **dan** oleh `IsEDMRiSlip` (`LOGIC=A OR B`:
`EdmTypeNew = 4 | Type = 0`). Nama "RiSlip" tidak berhubungan dengan label "Reff. Number".

Bentrok ketiga `[terverifikasi]`: `IsEdmAdjRate` (`LOGIC=(A AND B AND C) OR D` dengan
`Type = 3 | Rule IsEdmAdjTSI evaluates to true`) dan `IsEdmAddObject`
(`Type = 5 | Rule IsEdmAdjTSI evaluates to true`). Karena `Type = 3` dan `Type = 5` tidak pernah
ditawarkan UI, **kedua rule ini efektif = `IsEdmAdjTSI`** — ketiganya benar bersamaan. Ini sesuai
label UI `Type = 2` yang memang menggabungkan "TSI / Add Object / Rate / Premium".

#### Tipe yang ditawarkan per `EdmTypeNew` (VERIFIED dari `SetEdmType.xml`)

| `EdmTypeNew` | Nama variabel di rule | `Type` yang ditawarkan |
| --- | --- | --- |
| `1` | `TypePenambahan` | 1, 2, 4, 6, 7, 8, 11, 12 |
| `2` | `TypePengurangan` | 2, 4, 6, 7, 8, 11 — **plus** 12 |
| `3` | `TypePerubahan` | 9, 0, 11 |

**Bug terverifikasi di `SetEdmType.xml` step `RH_1.pySteps(3)`** (`when .Quotation.EdmTypeNew==2`):
sepuluh baris pertama mengisi `TypePengurangan.pxResults`, tetapi dua baris terakhir mengisi
**`TypePenambahan.pxResults`**:

```
SET TypePengurangan.pxResults(<APPEND>).CARI1 = 11
SET TypePengurangan.pxResults(<LAST>).CARI2  = "Adjustment Share Cedant"
SET TypePenambahan.pxResults(<APPEND>).CARI1 = 12          <- salah daftar
SET TypePenambahan.pxResults(<LAST>).CARI2  = "Adjustment PPN/PPH"
```

Efek: pada `EdmTypeNew==2`, opsi "Adjustment PPN/PPH" menempel pada daftar penambahan, bukan
pengurangan. Kandidat perbaikan, bukan keputusan migrasi. `[terverifikasi]`

### 3.3 `OfferFacIn.QuotationData.EdmType` dan `EdmTypeNew`

| Properti | Nilai | Bukti arti |
| --- | --- | --- |
| `EdmType` | `1` | `[terverifikasi]` — `Endorsment Fac In\Activity\SetErrorBatalEndorsement_Act.xml`, deskripsi step: `Get EdmType =1 (Batal Sejak Semula)`; juga komentar di `CountPremiEDM_DT.xml`: `Appending data to currencylist jika EdmType==1 ( batal sejak semula )` |
| | `2` | `[terverifikasi]` — `Endorsment Fac In\Activity\CountPaymentEdm_Act.xml`, komentar: `untuk EDM Batal EDMType = 2` (pembatalan, ragam kedua) |
| | `3` | hanya muncul sebagai kondisi (`EdmType==3`, `EdmType==3&&Datain1.CARI9==0`) — arti **belum terverifikasi** |
| | `4` | `[terverifikasi]` — `NB FacIn\DataTransform\CountPremiEDM_DT.xml`, komentar step: `EdmType = 4 = Penambahan / Pengurangan / Perubahan` |
| | `0`, `≥5` | tidak pernah muncul — **celah** |
| `EdmTypeNew` | `1` | `[terverifikasi]` = Penambahan (`SetEdmType.xml` mengisi `TypePenambahan`) |
| | `2` | `[terverifikasi]` = Pengurangan (`TypePengurangan`) |
| | `3` | `[terverifikasi]` = Perubahan (`TypePerubahan`) |
| | `4` | `[terverifikasi]` diuji oleh `When\IsEDMRiSlip.xml` (`EdmTypeNew = 4`) dan sebagai default — arti `belum terverifikasi` |
| | `0`, `≥5` | tidak pernah muncul — **celah** |

**Default keras `[terverifikasi]`:** `@if(pyWorkPage.Quotation.EdmType=="","4",pyWorkPage.Quotation.EdmType)`
— nilai kosong dianggap `"4"` (string). Namun di tempat lain `EdmType` dibandingkan sebagai angka
(`EdmType==4`, `EdmType!=4`) dan sebagai string (`EdmType=="4"`, `EdmType!="4"`, `EdmType=="2"`).

**PERINGATAN nama rule vs perilaku** `[terverifikasi]`: rule bernama `IsEDMRiSlip` menguji
`EdmTypeNew` (bukan `EdmType`) — satu-satunya rule di korpus yang menguji `EdmTypeNew`. Semua
keluarga `IsEdm*` lain menguji `EdmType = 4`. Nama tidak mengungkapkan perbedaan properti ini.

Section terkait yang ditemukan: `Endorsment Fac In\Section\EdmType1.xml`, `EdmType2.xml`,
`EdmType3.xml` — satu section per nilai `EdmTypeNew`. `[terverifikasi]` (nama berkas + referensi
`TypePenambahan.pxResults` di `EdmType1.xml`).

### 3.4 `EmailType*` — sinyal keputusan underwriting

Sebelas properti berbeda `[terverifikasi]`
(`& lit.ps1 -Pattern '(EmailType\w*)\s*[!=]=*\s*"?(\d+)"?'` → 11 distinct):
`EmailType`, `EmailTypeBinding`, `EmailTypeCeding`, `EmailTypeLife`, `EmailTypeQuotation`,
`EmailTypeRetro`, `EmailTypeRetroSlip`, `EmailTypeTL`, `EmailTypeUW`, `EmailTypeUWEDM`,
`EmailTypeUWPolicy`.

**Pemetaan nilai → status history**, dari `NB FacIn\Activity\InsertHistoryAkseptasiPega.xml`
(prakondisi step 2–7) `[terverifikasi]`:

| Nilai | `InsertHistory.CARI5` | Prakondisi (dikutip sebagian) |
| --- | --- | --- |
| `1` | `"ACCEPT"` | `PolicyTreatyIn.IsApproved=="1" \|\| EmailTypeUWEDM==1 \|\| EmailType==1 \|\| EmailTypeBinding==1 \|\| EmailTypeRetro==1 \|\| EmailTypeRetroSlip==1 \|\| EmailTypeUW==1 \|\| EmailTypeUWPolicy==1 \|\| EmailTypeTL==1` |
| `2` | `"REJECT"` | `PolicyTreatyIn.IsApproved=="0" \|\| EmailTypeUWEDM==2 \|\| EmailType==2 \|\| …` |
| `3` | `"ASK"` | `EmailType==3 \|\| EmailTypeRetro==3 \|\| EmailTypeUW==3 \|\| EmailTypeUWPolicy==3 \|\| EmailTypeTL==3` |
| `4` | `"BANDING"` | `EmailType==4 \|\| EmailTypeBinding==4` |
| `7` | `"DECLINE"` | `EmailTypeQuotation=="2" \|\| EmailTypeUWEDM==7 \|\| EmailType==7 \|\| …` |
| `9` | `"REVISE"` | `EmailTypeBinding==9` |
| — | `"INPUT"` | `Param.Status=="InboxAdmin"` (bukan EmailType) |

**Empat temuan material:**

1. **`EmailTypeQuotation == "2"` memetakan ke `"DECLINE"`, sedangkan semua `EmailType*` lain
   memetakan `2` ke `"REJECT"`** `[terverifikasi]`. Satu nilai literal, dua arti, tergantung
   properti. **`[pertanyaan terbuka]` — MEMBLOKIR.**
2. **Celah nilai:** `5`, `6`, `8` tidak dipetakan ke status apa pun. Namun `EmailTypeCeding==5` dan
   `EmailTypeCeding==6` **memang muncul** di korpus (2 berkas masing-masing) — artinya dua keputusan
   ceding tidak menghasilkan baris history. `[terverifikasi]` `[pertanyaan terbuka]`
3. **`EmailType==9` tidak memetakan ke `"REVISE"`** — hanya `EmailTypeBinding==9`. Tetapi
   `SetAkseptasiProposal_DT` step 4 memetakan `.EmailType==9 → ProposalAcceptStatus=9`. Jadi
   `EmailType==9` mengubah status proposal tanpa jejak history. `[terverifikasi]`
4. **`0` tidak pernah muncul** untuk properti mana pun — celah.

### 3.5 `ProposalAcceptStatus`

Dari `NB FacIn\DataTransform\SetAkseptasiProposal_DT.xml` `[terverifikasi]`:

| Step | Kondisi | Aksi |
| --- | --- | --- |
| 1/1.1 | `.EmailType==1` | `.ProposalAcceptStatus = 1` |
| 2/2.1 | `.EmailType==2` | `= 2` |
| 3/3.1 | `.EmailType==3` | `= 3` |
| 4/4.1 | `.EmailType==9` | `= 9` |
| 5/5.1 | `.EmailType==4 \|\| .EmailTypeBinding==4 \|\| .EmailTypeUWPolicy==4` | `= 4` |
| 6 | — | `pyWorkPage.QuotationList(1).Email = OperatorID.pyAddresses(Email).pyEmailAddress` |

Nilai yang ada: **1, 2, 3, 4, 9**. Celah: `0`, `5`–`8`. Khususnya **`7` (DECLINE) tidak pernah
masuk** `ProposalAcceptStatus` `[terverifikasi]`.

Konsumen `[terverifikasi]`:
- `NB FacIn\Activity\SetBanding_ACT.xml` step 1: prakondisi `pyWorkPage.ProposalAcceptStatus=="4"`
  (dibandingkan **sebagai string**, sedangkan `SetAkseptasiProposal_DT` menulisnya sebagai **angka** `4`).
- `NB FacIn\DecisionTable\IsUWAccepted.xml` — merujuk `ProposalAcceptStatus`, **baris hasil tidak
  terekspor** (lihat §3.9).

Audit: `& lit.ps1 -Pattern '(ProposalAcceptStatus\s*[!=]=*\s*(?:&quot;|")?\d*)'` → 6 bentuk:
`= "1`, `= 4`, `=="4`, `==1`, `==2`, `==3`.

### 3.6 `ProRateType`

| Nilai | Muncul di | Bukti arti |
| --- | --- | --- |
| `1` | `ProRateType=="1"` (4 berkas) | `belum terverifikasi` |
| `2` | `ProRateType=="2"` (4 berkas) | `belum terverifikasi` |
| `3` | 11 berkas (`==3`), 8 berkas (`!=3`) | `[terverifikasi]` — `NB FacIn\Activity\downloadPenawaranFac_act.xml`: `@If(pyWorkPage.OfferFacIn.ProRateType==3,"Dibayarkan Secara Tahunan","Dibayarkan Sekaligus")` dan `If(ProRateType==3,'Tahunan','Sekaligus')` → **3 = pembayaran tahunan**, selain-3 = sekaligus |
| `4` | 11 berkas (`==4`), 9 berkas (`!=4`) | `belum terverifikasi`; dipakai di `setCoveragePeriode_act.xml` bersama `RIRiskCode`/`RIRiskName` dan di `NB FacIn\Section\Periode.xml` bersama `PolicyType=='0'` |
| `0` | tidak pernah muncul | **celah** |

**Perbandingan campur tipe `[terverifikasi]`:** `NB FacIn\Section\Periode.xml` memakai
`.ProRateType=='4' && .QuotationData.PolicyType=='0'` (**string**) sementara berkas yang sama juga
memuat `pyWorkPage.OfferFacIn.ProRateType==4` dan `!=4` (**angka**).

**Dampak uang `[terverifikasi]`:** `CountPremiNusareRetro_Act` step 2.1.4 menihilkan
`.PremiLifeNusantaraRe` bila `ProRateType != 3`. Karena nilai `1`, `2`, `4` semuanya `!= 3`, tiga
dari empat nilai menihilkan premi retro Nusantara Re. `[pertanyaan terbuka]` apakah ini benar —
**MEMBLOKIR**.

### 3.7 Enumerasi lain yang terbaca

| Properti | Nilai literal | Arti | Bukti |
| --- | --- | --- | --- |
| `CoverageBasis` | `1,2,3,4,5` | 2 = First Loss (`TSI*FirstLoss/100`), 3 = EML/PML (`TSI*EmlPml/100`), 4&5 = Sub Limit (`TSI*SubLimit/100`), 5 juga = layering (punya `.LayerList`) | `[terverifikasi]` dari rumus `TSILiability` di `CountPremi_ACT` dan cabang `.LayerList` di `SetTSIPremiSpreaded_FacIn` step 7.1.1.1.1. `1` tidak menetapkan `TSILiability` — arti `belum terverifikasi`. `0` celah. |
| `ShareCedantType` | `0,1` | `0` = GROSS, `1` = SHARE RNM | `[terverifikasi]` — `pyStepsDescription` di `CountPremiCedant_Act` |
| `Param.PremiStatus` | `"percent"`, `"amount"` | `percent` = rate diberikan, premi dihitung; `amount` = premi diberikan, rate diturunkan | `[terverifikasi]` dari arah rumus di `CountPremi_ACT` |
| `Param.Type` (spreading) | `"percent"`, `"amount"` | idem untuk TSISpreaded/SharePercentage | `[terverifikasi]` `SetTSIPremiSpreaded_FacIn` |
| `IsProRate` | `"Prorate"`, `"ShortPeriod"` | `[terverifikasi]` — `CountPremi_ACT` step `RH_4.pySteps(9)` menyetel `"ShortPeriod"` | |
| `TypeTax` | `"Inclusive"` | cabang lain implisit | `[terverifikasi]`; nilai non-Inclusive `belum terverifikasi` |
| `StatusSyariah` | `1` | mengambil `PERCENTTABARUFUND` + `PERCENTUJROHFEE` | `[terverifikasi]` `CountPremiAndTSINusantaraRe_ACT` step 5.3 |
| `PolicyTreatyIn.IsApproved` | `"0"`, `"1"`, `0`, `1` | `"1"` → ACCEPT, `"0"` → REJECT | `[terverifikasi]` dari `InsertHistoryAkseptasiPega`. Dibandingkan sebagai string **dan** angka. |
| `FlagOnGoingPolicy` | `0,1,2` | `0` diset bersama `IsCedingConfirm="Offer"`; `1` bersama `"Policy"` | `[terverifikasi]` `SetBanding_ACT` step 3.x vs 4.x. Nilai `2` `belum terverifikasi`. |
| `IsCedingConfirm` | `"Offer"`, `"Policy"`, `"Binding"`, `"Accepted"`, `"accepted"`, `"notconfirmed"` | `belum terverifikasi` | `[terverifikasi]` bahwa **`"Accepted"` dan `"accepted"` dua ejaan berbeda** dipakai: `IsCedingConfirm="Accepted"` (3 berkas) vs `IsCedingConfirm != "accepted"` / `= "accepted"` (3 berkas). Bila perbandingan case-sensitive, ini bug. |
| `Local.FlagBatal` | `1`, `2` | `1` → komentar `-- Set Nilainya (pct share dkk) --> Batal Awal`; `2` → `--> Batal Prorata` | `[terverifikasi]` `NB FacIn\Activity\SetReinsurerEndorsement_Act.xml` |
| `EdmStatus` | `2` | komentar `Get EdmStatus=2 / Batal Internal` | `[terverifikasi]` `Endorsment Fac In\Activity\SetValueToEDMWork.xml` |
| `TreatyType` | `1001410`, `1001411`, `1001412`, `1000036` | `belum terverifikasi` | `[terverifikasi]` kemunculan di `hitungspreadingotomatis_act` dan `.TreatyType == "1000036"` (dibandingkan string) |
| Kode coverage | `100815`, `100825`, `100837` | `belum terverifikasi` | `[terverifikasi]` `CountPremiFacIn_act` step 7/8 |

### 3.8 Status yang ditulis ke `HISTORYAKSEPTASIPEGA`

**Berkas SQL:** `NB FacIn\RDBList\InsertHistoryAkseptasiPega_Sql.xml` `[terverifikasi]`:

```sql
BEGIN
INSERT INTO HISTORYAKSEPTASIPEGA
(ID_PEGA, Tgl_Transfer, Status, Username, Workbasket, ID_KOMITE)
VALUES
( {InsertHistory.CARI1}, sysdate, {InsertHistory.CARI5}, {InsertHistory.CARI4},
  {InsertHistory.CARI2}, {InsertHistory.CARI6} );
COMMIT;
END;
```

Pemetaan kolom (dari `NB FacIn\Activity\InsertHistoryAkseptasiPega.xml` step 1) `[terverifikasi]`:

| Kolom Oracle | Sumber |
| --- | --- |
| `ID_PEGA` | `CARI1` = `pyWorkPage.pzInsKey` |
| `Status` | `CARI5` ∈ {`ACCEPT`, `REJECT`, `ASK`, `REVISE`, `DECLINE`, `BANDING`, `INPUT`} |
| `Username` | `CARI4` = `OperatorID.pyUserName` |
| `Workbasket` | `CARI2` = **`pyWorkPage.PositionNote`** |
| `ID_KOMITE` | `CARI6` (sumber tidak terbaca di step 1) |
| `Tgl_Transfer` | `sysdate` (bukan `CARI3` = `@CurrentDateTime()` yang diset tapi tidak dipakai) |

**BUKTI KUNCI:** kolom bernama `Workbasket` diisi dari properti bernama `PositionNote`. Ini
membuktikan keduanya adalah **satu ruang nama** (lihat §5).

`[terverifikasi]` Step 8 mem-*override* dua kolom untuk kasus admin:

```
when Param.Status=="InboxAdmin":
     InsertHistory.CARI5 = "INPUT"
     InsertHistory.CARI2 = "ReasFacInAdmin"        <- literal, bukan PositionNote
     InsertHistory.CARI1 = "ASM-FW-GISFW-WORK "+Param.idpega
```

**Pembacaan history** — `NB FacIn\RDBList\GetHistoryAccPega_SQL.xml` `[terverifikasi]`:

```sql
select PXSAVEDATETIME AS TGL_TRANSFER from DATAPEGA.pc_History_ASM_FW_GISFW_Work
where PXHISTORYFORREFERENCE ={pyWorkPage.pzInsKey}
  and PYASSIGNEDTO is not null and pylabel is not null
  and (PXADDEDBYID not in ('<1 identitas operator>','System') or PXADDEDBYID is null)
  order by PXSAVEDATETIME
```

Query ini membaca **tabel internal Pega** `DATAPEGA.pc_History_ASM_FW_GISFW_Work` dan memuat
**1 guard berbasis identitas operator** di klausa `NOT IN`. Tabel ini hilang bersama Pega →
**`[pertanyaan terbuka]` MEMBLOKIR**: sumber pengganti untuk data ini.

### 3.9 DecisionTable — baris hasil TIDAK terekspor

`[terverifikasi]` **Seluruh 12 DecisionTable** di korpus hanya mengekspor properti yang dirujuk,
bukan baris kondisi→hasil. Bukti: `<pyPropertySets REPEATINGTYPE="PageList"/>` kosong dan tidak ada
tag `pyReturn`/`pyResult`.

```powershell
foreach ($f in (Get-ChildItem "D:\migrasi\RNM\NB FacIn\DecisionTable" -File)) {
  $c=[IO.File]::ReadAllText($f.FullName)
  "{0,-40} refs=[{1}]" -f $f.Name,
     ((([regex]::Matches($c,'<pyRuleName>(.*?)</pyRuleName>','Singleline')|%{$_.Groups[1].Value}|Sort-Object -Unique) -join ' | ')) }
```

| DecisionTable | Properti yang dirujuk | Status |
| --- | --- | --- |
| `IsUWAccepted` | `ProposalAcceptStatus` | **baris hilang — MEMBLOKIR** |
| `isApproved` | `PolicyTreatyIn` (kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`) | **baris hilang — MEMBLOKIR** |
| `BusinessType_DeT` | `Quotation.BusinessOldId`, `GroupPanel` | baris hilang |
| `MappingCoverage` | `MappingVariable`, `MappingVariableNote` | baris hilang |
| `MappingAdditionalCoverageIndex` | `Coverage` | baris hilang |
| `MappingOutgoIndex`, `MappingOutgoIndex2` | `ColumnOrder`, `ColumnValue` | baris hilang |
| `LicensePlatRegion_DeT` | (kelas `ASM-FW-GISFW-Data-Vehicle`) | baris hilang |
| `GetMimeType` | `ext` | baris hilang |
| `SetUploadHubAW1/2/3` | `HUB_1`/`HUB_2`/`HUB_3` | baris hilang |

Sebaliknya, satu-satunya `DecisionTree` (`Tree_ShortPeriod`) **terekspor lengkap** — lihat §1.6.

---

## 4. Kode jabatan / token routing approver

### 4.1 Ruang nama `pyWorkPage.PositionNote` — 25 token literal

Perintah audit:

```powershell
# lit.ps1 (lihat lampiran) dengan pola:
& lit.ps1 -Pattern 'PositionNote\s*[!=]=?\s*[&quot;''"]+(Reas[A-Za-z0-9_]*)'
```

| Token | Jml berkas |
| --- | --- |
| `ReasFacInAdmin` | 12 |
| `ReasFacInDepHeadUnderwriting` | 10 |
| `ReasFacInDepHeadUnderwritingLife` | 1 |
| `ReasFacInFacultativeDivHead` | 10 |
| `ReasFacInFinDivHead` | 7 |
| `ReasFacInGroupLeader` | 10 |
| `ReasFacInJuniorUnderwriting` | 8 |
| `ReasFacInJuniorUnderwritingA` | 14 |
| `ReasFacInManagerTeknik` | 10 |
| `ReasFacInMarketing` | 26 |
| `ReasFacInMarketingDirector` | 10 |
| `ReasFacInMedicalLife` | 3 |
| `ReasFacInSeniorUnderwriting` | 12 |
| `ReasFacInTeamLeader` | 28 |
| `ReasFacInTechnicalDirector` | 6 |
| `ReasFacInUnderwriting` | 19 |
| `ReasFacInUnderwritingFinancial` | 13 |
| `ReasFacInUnderwritingLife` | 4 |
| `ReasFacOutAdmin` | 3 |
| `ReasFacOutHead` | 2 |
| `ReasTreatyInAdmin` | 3 |
| `ReasTreatyInDeptHead` | 1 |
| `ReasTreatyInDirector` | 1 |
| `ReasTreatyInGroupLeader` | 1 |
| `ReasTreatyInSecHead` | 2 |

**TOTAL = 25** `[terverifikasi]`

Token ke-26 `ReasFacInDirector` **hanya muncul sebagai nama workbasket**, tidak pernah dibandingkan
dengan `PositionNote` — lihat §5.

### 4.2 Ruang nama `pyWorkPage.LetterNo` — 14 token literal (HURUF BESAR)

```powershell
& lit.ps1 -Pattern 'LetterNo\s*[!=]=*\s*[''"&quot;]+([A-Za-z0-9_ .\-]*)'
```

| Token | Jml berkas |
| --- | --- |
| `DEPHEADUNDERWRITER` | 6 |
| `DEPTHEADUWLIFE` | 1 |
| `DIREKTURMARKETING` | 6 |
| `DIREKTURTEKNIK` | 6 |
| `JUW_A` | 6 |
| `JUW_B` | 3 |
| `KADIVFACULTATIVE` | 6 |
| `KADIVFINANCIAL` | 3 |
| `KADIVTEKNIK` | 6 |
| `MANAGERTEKNIK` | 6 |
| `SENIORUW` | 6 |
| `TREATYINDEPTHEAD` | 3 |
| `UNDERWRITER` | 6 |
| *(string kosong)* | 7 |

**TOTAL non-kosong = 13 + `JUW_B` = 14** `[terverifikasi]`
(`& lit.ps1 -Pattern 'LetterNo\s*=\s*(?:&quot;\|")([^"&<]*)'` → 13 nilai untuk penugasan;
`JUW_B` hanya muncul sebagai perbandingan `==`, tidak pernah ditugaskan — **celah penulis**.)

### 4.3 Pemetaan `LetterNo` → `PositionNote` (terverifikasi penuh)

**Berkas:** `NB FacIn\Activity\SetBanding_ACT.xml`, step `RH_1.pySteps(3).pySteps(1..11)` (cabang
"set tiket UNTUK NB") dan `RH_1.pySteps(4).pySteps(1..11)` (cabang "set tiket UNTUK EDM").

| `LetterNo` | `PositionNote` yang diset | `Local.TicketNext` (NB) | `Local.TicketNext` (EDM) |
| --- | --- | --- | --- |
| `JUW_A` | `ReasFacInJuniorUnderwritingA` | `JUW_AOffer` | `JUW_APolicy` |
| `JUW_B` | `ReasFacInJuniorUnderwriting` | `JUW_BOffer` | `JUW_BPolicy` |
| `UNDERWRITER` | `ReasFacInUnderwriting` | `UWOffer` | `UWPolicy` |
| `SENIORUW` | `ReasFacInSeniorUnderwriting` | `SUWOffer` | `SUWPolicy` |
| `DEPHEADUNDERWRITER` | `ReasFacInDepHeadUnderwriting` | `DepHeadUWOffer` | `DepHeadUWPolicy` |
| *(bukan LetterNo — dipicu `PositionNote=="ReasFacInUnderwritingFinancial"`)* | `ReasFacInUnderwritingFinancial` | `SUWPolicyFinancial` | `SUWPolicyFinancial` |
| `MANAGERTEKNIK` | `ReasFacInManagerTeknik` | `DivHeadUWOffer` | **`DivHeadUWOffer`** ← lihat catatan |
| `KADIVFACULTATIVE` | `ReasFacInFacultativeDivHead` | `DivHeadFacOffer` | `DivHeadFacPolicy` |
| **`KADIVTEKNIK`** | **`ReasFacInGroupLeader`** | `DivHeadOffer` | `DivHeadPolicy` |
| `DIREKTURMARKETING` | `ReasFacInMarketingDirector` | `DirectorOPOffer` | `DirectorOPPolicy` |
| `DIREKTURTEKNIK` | `ReasFacInTechnicalDirector` | `DirectorOffer` | `DirectorPolicy` |

**Empat ketidakkonsistenan ejaan / penamaan — semua `[terverifikasi]`:**

1. **`KADIVTEKNIK` → `ReasFacInGroupLeader`.** Token jabatan "Kepala Divisi Teknik" memetakan ke
   posisi bernama "Group Leader". Tidak ada `ReasFacInKadivTeknik`. Menebak posisi dari nama token
   akan salah.
2. **Cabang EDM untuk `MANAGERTEKNIK` menyetel `DivHeadUWOffer`** (varian *Offer*), bukan
   `DivHeadUWPolicy`, padahal seluruh cabang EDM lain memakai varian *Policy*. Token
   `DivHeadUWPolicy` **ada** dan dipakai di Flow (3 berkas). Kandidat bug salin-tempel.
3. **`ReasFacInUnderwritingFinancial` dipicu oleh `PositionNote` itu sendiri**, bukan oleh
   `LetterNo` — satu-satunya baris yang memakai prakondisi berbeda dari 10 baris lain di blok yang
   sama. Tidak ada `LetterNo` `KADIVFINANCIAL` yang memicunya meski token itu ada (3 berkas).
4. **Empat token `LetterNo` tidak punya baris di `SetBanding_ACT`**: `DEPTHEADUWLIFE`,
   `KADIVFINANCIAL`, `TREATYINDEPTHEAD`, dan string kosong. `[pertanyaan terbuka]` ke mana mereka
   dirutekan.

**Penulis `LetterNo`** — `NB FacIn\DataTransform\SetLetterNo.xml` hanya menulis **2 dari 14** nilai
`[terverifikasi]`:

```
Step 1   when pyWorkPage.PositionNote=="ReasFacInJuniorUnderwritingA"  -> LetterNo = "JUW_A"
Step 2   when pyWorkPage.PositionNote=="ReasFacInUnderwriting"         -> LetterNo = "UNDERWRITER"
```

`[pertanyaan terbuka]` **MEMBLOKIR**: siapa yang menulis 12 nilai `LetterNo` lainnya. Rule `When`
`ToDepHeadUW`, `ToDirMarketing`, `ToDirTeknik`, `ToKadivFacultative`, `ToKadivFin`,
`ToKadivTeknik`, `ToManagerTeknik`, `ToSeniorUW`, `ToTREATYDEPTHEAD`, `ToUW`, `ToJUW_A`,
`LetterNoNull` hanya **membaca** `LetterNo`.

Isi rule `When` routing (semua `[terverifikasi]`, `LOGIC=A`):

| Rule | Kondisi |
| --- | --- |
| `ToUW` | `pyWorkPage.LetterNo = "UNDERWRITER"` |
| `ToJUW_A` | `pyWorkPage.LetterNo = "JUW_A"` |
| `LetterNoNull` | `pyWorkPage.LetterNo = ""` |
| `IsPKSASM` | `pyWorkPage.OfferFacIn.IsB2B = "ASM"` |
| `IsEdmInternalRetro` | `pyWorkPage.OfferFacIn.QuotationData.EndorsementInternalRetro = 1` |

Catatan: kelima rule di atas **terbaca penuh** di korpus `D:\migrasi\RNM\`.

### 4.4 Ejaan `JABATAN` di SQL vs ejaan `PositionNote`

`[terverifikasi]` Semua 13 rule `RDBList` keluarga `GetLimitAkseptasi*` memilih kolom
`JABATAN AS CARI1`:

```sql
-- NB FacIn\RDBList\GetLimitAkseptasi_SQL.xml
SELECT JABATAN AS CARI1, MAX_LIMIT_IDR AS CARI2, MAX_LIMIT_USD AS CARI3,
       BATAS_WAKTU AS CARI4, NAMA AS CARI5
  FROM POOLDATA.M_LIMIT_PROPERTYY
 WHERE team_group = {pyWorkPage.OfferFacIn.QuotationData.TeamGroup}
   AND LIMIT_BOTTOM > (SELECT LIMIT_BOTTOM FROM POOLDATA.M_LIMIT_PROPERTYY
                        WHERE team_group = {pyWorkPage.OfferFacIn.QuotationData.TeamGroup}
                          AND LOGIN = {OperatorID.pyUserIdentifier})
   AND LIMIT_BOTTOM < {DataSearch.CARID2}
 ORDER BY LIMIT_BOTTOM ASC
```

Tabel limit yang dipakai (5 tabel berbeda) `[terverifikasi]`:

| Rule | Tabel | Kolom ambang |
| --- | --- | --- |
| `GetLimitAkseptasi_SQL` | `POOLDATA.M_LIMIT_PROPERTYY` | `LIMIT_BOTTOM` |
| `GetLimitAkseptasiBanding_SQL` | `POOLDATA.M_LIMIT_PROPERTYY` | `LIMIT_BOTTOM2` |
| `GetLimitAkseptasiJUWA_SQL` | `POOLDATA.M_LIMIT_PROPERTYY` | `LIMIT_BOTTOM`, dipilih via `JABATAN = {DataSearch.CARI5}` |
| `GetLimitAkseptasiNonFire*_SQL` | `POOLDATA.M_LIMIT_NONPROPANDENGG` | `LIMIT_BOTTOM` / `LIMIT_BOTTOM2` |
| `GetLimitAkseptasiNonPrefer_SQL` | `M_LIMIT_PROPERTY_NON_PREFERREDD` *(tanpa `POOLDATA.`)* | `LIMIT_BOTTOM` |
| `GetLimitAkseptasiNonPreferJUWA_SQL` | `POOLDATA.M_LIMIT_PROPERTY_NON_PREFERREDD` *(dengan `POOLDATA.`)* | `LIMIT_BOTTOM` |
| `GetLimitAkseptasiPreferedComm_SQL` | `M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL` *(tanpa skema)* | `LIMIT_BOTTOM` |
| `GetLimitAkseptasiPreferedCommJUWA_SQL` | `POOLDATA.M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL` | `LIMIT_BOTTOM` |
| `GetLimitAkseptasiBond_SQL` | `M_LIMIT_FINANCIALINS` | `LIMITBOND_BOTTOM` |
| `GetLimitAkseptasiKreditCL_SQL` | `M_LIMIT_FINANCIALINS` | `LIMITCREDITCL_BOTTOM` |
| `GetLimitAkseptasiKreditNCL_SQL` | `M_LIMIT_FINANCIALINS` | `LIMITCREDITNCL_BOTTOM` |
| `GetLimitAkseptasiLife_SQL` | `M_LIMIT_LIFE` | `LIMIT_BOTTOM` |

**Tiga temuan `[terverifikasi]`:**

1. **`GetLimitAkseptasiBanding_SQL` dan `GetLimitAkseptasiNonFireBanding_SQL` memfilter dengan
   `LIMIT_BOTTOM2` tetapi mengurutkan `ORDER BY LIMIT_BOTTOM ASC`** — kolom filter dan kolom urut
   berbeda. Jika urutan kedua kolom tidak sama, urutan tangga banding salah. Kandidat bug.
2. **Kualifikasi skema tidak konsisten** untuk tabel yang sama: `M_LIMIT_PROPERTY_NON_PREFERREDD`
   tanpa `POOLDATA.` di varian non-JUWA, dengan `POOLDATA.` di varian JUWA. Idem untuk
   `M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL`.
3. **Kriteria pemilihan baris awal berbeda**: varian non-JUWA memakai
   `LOGIN = {OperatorID.pyUserIdentifier}` (identitas login), varian JUWA memakai
   `JABATAN = {DataSearch.CARI5}` (kode jabatan). Dua sumbu yang berbeda.

**Perbandingan ejaan — MEMBLOKIR `[pertanyaan terbuka]`:**

`Endorsment Fac In\Activity\InputAddendumFacIn_PreAct.xml` dan
`NB FacIn\Activity\InputOfferFacInEngineerUW_preACT.xml` `[terverifikasi]`:

```
pyWorkPage.OfferFacIn.IsOccupException!="1" && LimitAkseptasi.pxResults(<LAST>).CARI10==pyWorkPage.PositionNote
```

`NB FacIn\Activity\SendOfferEmail.xml` dan `SetValidateDateUW_PostAct.xml` `[terverifikasi]`:

```
pyWorkPage.PositionNote==.CARI10
LimitAkseptasi.pxResults(Local.Index+1).CARI10
```

Artinya `LimitAkseptasi.pxResults(n).CARI10` **harus** berisi token ruang-nama `Reas*`. Namun:

- Semua SQL `GetLimitAkseptasi*` hanya mengalias **CARI1..CARI5** — tidak ada `CARI10`.
- `CARI10` **tidak pernah ditulis** di mana pun dalam korpus:

```powershell
$root="D:\migrasi\RNM"
foreach ($f in [IO.Directory]::EnumerateFiles("$root\NB FacIn","*.xml",'AllDirectories')) {
  $c=[IO.File]::ReadAllText($f); if($c -match 'LimitAkseptasi[^<>]{0,40}CARI10'){ $f } }
# -> hanya SendOfferEmail.xml (pembacaan), tidak ada penugasan
```

`[pertanyaan terbuka]` **MEMBLOKIR**: dari mana `CARI10` berasal, dan apakah nilai kolom
`M_LIMIT_*.JABATAN` di Oracle memakai ejaan `Reas*` (CamelCase, seperti `PositionNote`) atau ejaan
HURUF BESAR (seperti `LetterNo`). Ketiga korpus tidak memuat satu pun nilai baris tabel `M_LIMIT_*`.
Tanpa jawaban ini, tangga persetujuan tidak dapat diimplementasikan.

---

## 5. Antrean / posisi kerja (workbasket & ticket)

### 5.1 Workbasket dan `PositionNote` adalah SATU ruang nama — dibuktikan

Tiga bukti terpisah `[terverifikasi]`:

1. **Perbandingan langsung.** Ekspresi `pyWorkPage.PositionNote==.pyWorkBasketName` ada di korpus.
2. **Penulisan ke kolom DB.** `InsertHistoryAkseptasiPega` step 1 menyetel
   `InsertHistory.CARI2 = pyWorkPage.PositionNote`, dan SQL menuliskannya ke kolom bernama
   **`Workbasket`** (`InsertHistoryAkseptasiPega_Sql.xml`).
3. **Pembanding email.** `Set TempEmail.CARI28 = OperatorID.pyWorkBasketList(1).pyWorkBasketName`
   lalu diuji `TempEmail.CARI28 != .PositionNote` / `!= pyWorkPage.PositionNote` di banyak
   prakondisi.

Nilai literal yang muncul sebagai `pyWorkBasketName` `[terverifikasi]`:

```
OperatorID.pyWorkBasketList(1).pyWorkBasketName = ""
OperatorID.pyWorkBasketList(1).pyWorkBasketName = "ReasFacInAdmin"
OperatorID.pyWorkBasketList(2).pyWorkBasketName = "ReasFacInJuniorUnderwritingA"
OperatorID.pyWorkBasketList(2).pyWorkBasketName == 'ReasTreatyInAdmin'  /  != 'ReasTreatyInAdmin'
OperatorID.pyWorkBasketList(2).pyWorkBasketName != 'ReasFacInTeamLeader'
.pyWorkBasketName=="ReasFacOutAdmin" || =="ReasFacOutHead" || =="ReasFacOutGroupLeader" || =="ReasFacOutTechnicalDirector"
pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName = ReasFacInDirector      <- TANPA tanda kutip
pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName = ReasFacInGroupLeader   <- TANPA tanda kutip
pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName = ReasFacInMarketing     <- TANPA tanda kutip
pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName = ReasFacInUnderwriting  <- TANPA tanda kutip
```

**Tiga temuan `[terverifikasi]`:**

1. **Empat literal ditulis TANPA tanda kutip** (`= ReasFacInDirector`, dst.) sementara semua
   perbandingan `PositionNote` memakai tanda kutip. Di Pega, literal tanpa kutip di
   `pyConditionValue1String` diperlakukan sebagai referensi properti, bukan string.
   `[pertanyaan terbuka]` apakah empat kondisi ini pernah bernilai benar.
2. **Tiga token workbasket tidak ada di daftar `PositionNote` (§4.1):** `ReasFacInDirector`,
   `ReasFacOutGroupLeader`, `ReasFacOutTechnicalDirector`. Jadi ruang nama workbasket adalah
   **superset** dari nilai `PositionNote` yang pernah diuji.
3. `OperatorID.pyWorkGroup` adalah **sumbu ketiga** — nilai `'ReasLife'` diuji
   (`OperatorID.pyWorkGroup!='ReasLife'`). Berbeda dari workbasket maupun `PositionNote`.

Kesimpulan: **workbasket dan `PositionNote` bukan dua ruang nama terpisah** — keduanya berbagi
kosakata `Reas*` yang sama dan saling dibandingkan. Yang **benar-benar terpisah** adalah:

| Ruang nama | Bentuk | Jumlah | Contoh |
| --- | --- | --- | --- |
| Posisi/workbasket | `Reas…` CamelCase | 25 (+3 hanya-workbasket) | `ReasFacInSeniorUnderwriting` |
| Kode jabatan surat | HURUF BESAR | 14 | `SENIORUW`, `KADIVTEKNIK` |
| Ticket alur | CamelCase berakhiran `Offer`/`Policy` | 27 | `SUWOffer`, `DivHeadFacPolicy` |
| Work group | `Reas…` | ≥1 | `ReasLife` |

### 5.2 Token ticket — 27 nilai

```powershell
& lit.ps1 -Pattern '<PropertiesValue>"([A-Za-z_0-9]*(?:Offer|Policy|PolicyFinancial))"</PropertiesValue>'
```

| Token | Berkas | | Token | Berkas |
| --- | --- | --- | --- | --- |
| `AdminOffer` | 3 | | `DivHeadUWOffer` | 6 |
| `DepHeadUWOffer` | 3 | | `DivHeadUWPolicy` | 3 |
| `DepHeadUWPolicy` | 6 | | `HeadRetroPolicy` | 2 |
| `DirectorOffer` | 3 | | `JUW_AOffer` | 3 |
| `DirectorOPOffer` | 6 | | `JUW_APolicy` | 6 |
| `DirectorOPPolicy` | 6 | | `JUW_BOffer` | 3 |
| `DirectorPolicy` | 6 | | `JUW_BPolicy` | 3 |
| `DivHeadFacOffer` | 6 | | `Offer` | 5 |
| `DivHeadFacPolicy` | 6 | | `Policy` | 8 |
| `DivHeadFinPolicy` | 3 | | `SUWOffer` | 3 |
| `DivHeadOffer` | 6 | | `SUWPolicy` | 6 |
| `DivHeadPolicy` | 6 | | `SUWPolicyFinancial` | 6 |
| | | | `TLPolicy` | 3 |
| | | | `UWOffer` | 3 |
| | | | `UWPolicy` | 6 |

**TOTAL = 27** `[terverifikasi]`

Token yang **tidak** muncul di `SetBanding_ACT`: `AdminOffer`, `DivHeadFinPolicy`,
`HeadRetroPolicy`, `TLPolicy`, `Offer`, `Policy`. `[pertanyaan terbuka]` siapa yang menyetelnya.

`[terverifikasi]` Flow memuat 19 shape ticket (`rowdata REPEATINGINDEX="Ticket1"` … `"Ticket19"`)
dan satu label tombol `Set Ticket to admin offer`. Pemetaan nomor `TicketN` → nama ticket
`belum terverifikasi` (butuh `pyTicketShapes` yang tidak terbaca isinya di ekspor ini).

### 5.3 `LimitAkseptasi` — struktur tangga persetujuan

`[terverifikasi]` Ekspresi yang terbaca:

```
@LengthOfPageList(LimitAkseptasi.pxResults)>1
@LengthOfPageList(LimitAkseptasi.pxResults)>Local.Index
@LengthOfPageList(LimitAkseptasi.pxResults)==Local.Index
@LengthOfPageList(LimitAkseptasi.pxResults)=Local.Index            <- satu tanda '=' (penugasan?)
LimitAkseptasi.pxResults(1).CARI1        (= JABATAN)
LimitAkseptasi.pxResults(1).CARI5        (= NAMA)
LimitAkseptasi.pxResults(Local.Index+1).CARI5
LimitAkseptasi.pxResults(Local.Index+1).CARI10
LimitAkseptasi.pxResults(<LAST>).CARI10
LimitAkseptasi.pxResults(Param.idxlimit)
LimitAkseptasiLast.pxResults
```

**Temuan `[terverifikasi]`:** `@LengthOfPageList(LimitAkseptasi.pxResults)=Local.Index` memakai
**satu** tanda sama dengan di dalam ekspresi kondisi, di tengah `&&`/`||` bersama pembanding lain
(termasuk satu guard identitas operator). `[pertanyaan terbuka]` apakah Pega memperlakukan `=`
sebagai perbandingan atau penugasan dalam konteks ini — **MEMBLOKIR** untuk tangga persetujuan.

---

## 6. Guard berbasis identitas (jumlah + mekanisme saja)

Sesuai aturan anti-halusinasi, **nilai nama tidak disalin**. Hanya jumlah dan mekanisme.

```powershell
$root="D:\migrasi\RNM"; $files=@{}; $occ=0
$d=New-Object System.Collections.Generic.HashSet[string]
foreach ($corp in @("NB FacIn","RNW Fac In","Endorsment Fac In")) {
  foreach ($f in [IO.Directory]::EnumerateFiles("$root\$corp","*.xml",'AllDirectories')) {
    $c=[IO.File]::ReadAllText($f)
    $ms=[regex]::Matches($c,'(?:pyUserIdentifier|pxCreateOperator|pxCreateOpName|pyUserName|PXADDEDBYID)\s*(?:!=|==|=|not in \()\s*[''"]([A-Za-z][A-Za-z0-9_ ]{2,40})[''"]')
    if($ms.Count){ $files[$f]=$ms.Count; $occ+=$ms.Count; foreach($m in $ms){[void]$d.Add($m.Groups[1].Value.ToUpper())} } } }
"files=$($files.Count) occurrences=$occ distinct=$($d.Count)"
```

| Metrik | Nilai |
| --- | --- |
| Berkas yang memuat guard identitas | **56** |
| Total kemunculan | **316** |
| Literal identitas berbeda (setelah normalisasi huruf besar) | **12** — termasuk 1 literal non-orang (`System`), jadi **≤11 identitas orang** |

Sebaran per berkas (5 terbesar) `[terverifikasi]`:

| Berkas | Kemunculan |
| --- | --- |
| `NB FacIn\Section\InputDtlSpreadingCoverage_FacIn.xml` | 26 |
| `NB FacIn\Activity\CheckSpreadingProtectAnekaGolf_ACT.xml` | 8 |
| `NB FacIn\Activity\CountGrossPremiEDM_Act.xml` | 8 |
| `NB FacIn\Activity\CheckSpreadingProtectFire_ACT.xml` | 7 |
| `NB FacIn\Activity\CheckSpreadingProtectMCargoMBU_ACT.xml` | 7 |

`[terverifikasi]` Rule `When\IsGroup.xml` (`LOGIC=A OR B OR C OR D OR E`) bercabang atas:
**1 properti bisnis** (`.OfferFacIn.QuotationData.IsGroup = "Group"`) **+ 3 identitas operator**
(`OperatorID.pyUserIdentifier = <identitas>`) **+ 1 kode kontak marketing**
(`pyWorkPage.OfferFacIn.QuotationData.MarketingCode = "ASM-SFAGIS-WORK-CONTACT CON-54"`).
Rule `When\IsGroupCreate.xml` memuat 3 guard identitas; `When\IsSPVCreate.xml` memuat 4.

**Dampak uang `[terverifikasi]`:** `SetTSIPremiSpreaded_FacIn` step 1 melakukan `EXIT_MODEL` bila
`IsGroup` benar → **seluruh perhitungan spreading dilewati untuk 3 identitas operator tertentu dan
1 kode kontak**. Kandidat perbaikan berprioritas tinggi. `[pertanyaan terbuka]` **MEMBLOKIR**:
aturan bisnis apa yang sebenarnya dimaksud.

**Label status yang memuat nama orang `[terverifikasi]`:**

```powershell
# menghitung tanpa menyalin nilainya
$root="D:\migrasi\RNM"; $tot=0; $files=0
$d=New-Object System.Collections.Generic.HashSet[string]
foreach ($corp in @("NB FacIn","RNW Fac In","Endorsment Fac In")) {
  foreach ($f in [IO.Directory]::EnumerateFiles("$root\$corp","*.xml",'AllDirectories')) {
    $c=[IO.File]::ReadAllText($f)
    $ms=[regex]::Matches($c,"IS IN ([A-Z][A-Z ()]{2,40})&apos;S INBOX|IS IN ([A-Z][A-Z ()]{2,40})'S INBOX")
    if($ms.Count){$files++;$tot+=$ms.Count;foreach($m in $ms){[void]$d.Add($m.Groups[1].Value+$m.Groups[2].Value)}} } }
"files=$files occurrences=$tot distinct=$($d.Count)"
```

| Metrik | Nilai |
| --- | --- |
| Berkas dengan pola `"<prefix> IS IN <X>'S INBOX"` | **19** |
| Total kemunculan | **669** |
| Label berbeda | **23** |

Sebagian besar dari 23 label itu memuat **nama orang**, bukan jabatan (dibandingkan dengan
`"JUNIOR UNDERWRITER (A)"` dan `"UNDERWRITER"` yang generik). Nilai-nilainya **tidak disalin** ke
dokumen ini. Dalam migrasi, label ini harus diturunkan dari `PositionNote` melalui tabel
referensi, bukan dari string literal. `[pertanyaan terbuka]` apakah label ini muncul di layar
pengguna atau hanya di email.

---

## 7. Aturan representasi uang — kesimpulan berbasis bukti

### 7.1 Bukti yang dikumpulkan

**B1 — Nilai desimal masuk sebagai string berkoma desimal.** `[terverifikasi]`

```powershell
$root="D:\migrasi\RNM"; $n=0; $tot=0
foreach ($corp in @("NB FacIn","RNW Fac In","Endorsment Fac In")) {
  foreach ($f in [IO.Directory]::EnumerateFiles("$root\$corp","*.xml",'AllDirectories')) {
    $c=[IO.File]::ReadAllText($f)
    $m=[regex]::Matches($c,'replaceAll\(.{0,140}?,\s*","\s*,\s*"\."\s*\)')
    if($m.Count){ $n++; $tot+=$m.Count } } }
"files=$n occurrences=$tot"     # -> files=91 occurrences=1057
```

**91 berkas, 1.057 kemunculan** konversi koma→titik. Properti uang/persen yang dikonversi meliputi
`.PremiNusantaraRe` (21 berkas), `.TSINusantaraRe` (6), `.TSISpreaded` (6), `.RiCommAllObj` (12),
`.PctLoL` (6), `.SharePercentage` (3), `.RICommision` (6), `pyWorkPage.OfferFacIn.PercentShare` (2),
`.TotalTSINusaRe`, `.TotalTSITopRisk`, `.PCT_RATE`, `.PCT_DISCOUNT`, `.PCT_LOL`, `.PctTreatyLimit`.

**Komentar step yang menyatakannya eksplisit** `[terverifikasi]`:

```
tambah @replaceAll(.PctLoL,",",".") karena property bukan decimal
```

Ini pernyataan langsung dari pengembang sistem lama: properti tersebut **bukan tipe desimal**.

**B2 — Ambang uang dibandingkan sebagai STRING berkutip.** `[terverifikasi]`

```powershell
$root="D:\migrasi\RNM"; $tally=@{}
foreach ($corp in @("NB FacIn","RNW Fac In","Endorsment Fac In")) {
  foreach ($f in [IO.Directory]::EnumerateFiles("$root\$corp","*.xml",'AllDirectories')) {
    $c=[IO.File]::ReadAllText($f); if($c.IndexOf('000"') -lt 0){continue}
    foreach($m in [regex]::Matches($c,'([\w.]{3,40})\s*(&lt;=|&gt;=|&lt;|&gt;|==|!=)\s*"(\d{7,})"')){
      $k=$m.Groups[1].Value+" "+$m.Groups[2].Value+' "'+$m.Groups[3].Value+'"'
      if(-not $tally.ContainsKey($k)){$tally[$k]=New-Object System.Collections.Generic.HashSet[string]}
      [void]$tally[$k].Add($f) } } }
$tally.Keys | Sort-Object | %{ "{0}`tfiles={1}" -f $_, $tally[$_].Count }
```

Hasil — **8 bentuk distinct** `[terverifikasi]`:

| Ekspresi | Berkas |
| --- | --- |
| `.TSILiability > "3000000000"` | 7 |
| `.TSILiability <= "3000000000"` | 7 |
| `.TSILiability <= "6000000000"` | 7 |
| `local.TSIinIDR > "30000000000"` | 3 |
| `local.TSIinIDR <= "30000000000"` | 3 |
| `TotalRNMShare.CARI1 > "30000000000"` | 3 |
| `.TreatyType == "1000036"` | 4 |
| `.TreatyType != "1000036"` | 4 |

**Kontras yang menentukan `[terverifikasi]`:** ambang yang **sama** (30.000.000.000) dibandingkan
dua cara berbeda di korpus yang sama:

| Cara | Berkas |
| --- | --- |
| String mentah: `local.TSIinIDR>"30000000000"` | `NB FacIn\Activity\SetDataScoringRisk_act.xml` |
| String mentah: `TotalRNMShare.CARI1 > "30000000000"` | `NB FacIn\Activity\ScoringResult.xml` |
| Desimal: `pyWorkPage.OfferFacIn.TotalTSINusaRe>toDecimal("30000000000")` | `NB FacIn\Activity\CheckSpreadingProtect_ACT.xml`, `ProtectFIREMBUPA_Act.xml` |
| Desimal: `pyWorkPage.OfferFacIn.TotalTSITopRisk>@toDecimal("30000000000")` | `NB FacIn\Activity\CountTotalTSIPremiNusaRe_Act.xml` |

Perbandingan string leksikografis: `"4000000000" > "30000000000"` bernilai **BENAR** (karena `'4' >
'3'`) padahal 4 miliar < 30 miliar. Ambang ini menentukan tangga akseptasi dan proteksi spreading.
**`[pertanyaan terbuka]` MEMBLOKIR.**

**B3 — Ambang layer disimpan sebagai literal string.** `[terverifikasi]`
`NB FacIn\Activity\hitungspreadingotomatis_act.xml` step 2:
`local.bataslayer1 = "3000000000"` … `local.bataslayer4 = "25000000000"`, lalu dibandingkan
`ShareRNML<=local.bataslayer1`.

**B4 — Arah konversi pemisah desimal BERLAWANAN di dua rule INSERT bersaudara.** `[terverifikasi]`

| Rule | Ekspresi |
| --- | --- |
| `NB FacIn\RDBList\InsertIntoFacinSPreadLife_Sql.xml` | `Replace({DataPolis.CARI4},'.',',')` — **titik → koma** (10 kolom: CARI4,7,8,11-17) |
| `NB FacIn\RDBList\InsertIntoFacinSPreadLifeMonthly_Sql.xml` | `Replace({DataPolis.CARI4},',','.')` — **koma → titik** (11 kolom: CARI4,7,8,11-18) |

Dua rule sibling, kolom yang sama, arah konversi berlawanan. Salah satunya menulis angka dengan
pemisah yang salah ke Oracle. **`[pertanyaan terbuka]` MEMBLOKIR** untuk paralel run.

Rule INSERT lain memakai `To_number(Replace({Datain1.CARI9},',','.'))` `[terverifikasi]` —
konversi eksplisit ke `NUMBER` di sisi SQL, artinya nilai tiba di SQL sebagai **string berkoma**.

**B5 — Presisi pembulatan tidak seragam.** `[terverifikasi]` 44 kombinasi pembagi×presisi (§1.7);
pembagi `100` saja muncul dengan 9 presisi berbeda; rumus premi yang aritmetikanya sama
(`TSI*Rate*Prorate*…/1e9`) dibulatkan **4 desimal** di `HitungPremi_FacInDT` dan **20 desimal** di
`CountPremi_ACT`.

**B6 — Besaran nilai.** `[terverifikasi]` TSI dibandingkan terhadap 25.000.000.000 dan
30.000.000.000 (`hitungspreadingotomatis_act`, `SetDataScoringRisk_act`). Dikalikan rate dan empat
faktor persen sebelum dibagi hingga 1e13 → produk antara mencapai orde 10^23. `float64` punya 15–17
digit signifikan → **pasti** kehilangan presisi.

**B7 — Mata uang adalah dimensi eksplisit.** `[terverifikasi]` `CountPremiFacIn_act` melakukan loop
atas `CurrencyMaster.pxResults` dengan `Local.Currency = .CURR`, membandingkan
`Local.Currency==Local.ObjectCurrency`, dan menyimpan per `OfferFacIn.CurrencyList(.Name)`.
`HitungPremi_FacInDT` step 2.2/3.2 mengalikan `pyWorkPage.Policy.CurrencyValue` (kurs) dengan
default `1` bila kosong. `CountPremiCedant_Act` menghasilkan `CedingCedantList(i).CurrencyList`.

**B8 — Konstanta numerik ditulis sebagai string lalu dikonversi.** `[terverifikasi]`
`@String.toDecimal("100.0000")` (`CountASMShareTotal_ACT`),
`@toDecimal("100000000000")`, `@toDecimal("10000000000000")` (`CountPremi_ACT`),
`toDecimal("1000000000")` tanpa `@` (`CountPremiAndTSIRNMFireMBU_ACT` dsb.).

### 7.2 Aturan mengikat untuk target Go / Oracle / JSON

| Lapisan | Aturan | Dasar |
| --- | --- | --- |
| **Go — tipe** | `float32`/`float64` **dilarang** untuk semua nilai uang, TSI, rate, persen, dan limit. Pakai `decimal.Decimal`. | B5, B6 |
| **Go — mata uang** | Setiap nilai uang membawa kode mata uang sebagai field wajib (`Money{Amount decimal.Decimal; Currency string}`). Korpus mengelola mata uang sebagai dimensi terpisah di `CurrencyList`. | B7 |
| **Go — pembulatan** | Presisi pembulatan adalah **parameter per-rule**, bukan konstanta global. Setiap port menyebut rule Pega asal + presisi asli di komentar. Jangan menyeragamkan 4/8/20. | B5, §1.7 |
| **Go — urutan operasi** | Bagi-dan-bulatkan dilakukan **persis pada titik yang sama** seperti di rule asal. Contoh: `.PremiRp = round(premi,4) * kurs`, bukan `round(premi*kurs,4)`. Pembulatan bersarang (`CountRIComm`, `CalculateNetRate`) dipertahankan bersarang. | §1.1, §1.3, §2.1 |
| **Go — parsing masukan** | Parser menerima string berkoma desimal **dan** berkoma ribuan. Titik masuk: setiap tempat rule lama memanggil `@replaceAll(x,",",".")` atau `@toDecimal(x)` — 1.057 titik di 91 berkas. | B1 |
| **Go — perbandingan ambang** | Setiap ambang yang di Pega dibandingkan sebagai string **dicatat sebagai temuan**, bukan diam-diam dikonversi. Implementasi awal harus **mereproduksi perbandingan string** agar paralel run cocok; perubahan ke perbandingan numerik adalah keputusan bisnis terpisah. | B2, B3 |
| **Oracle** | Kolom numerik tetap `NUMBER`, bukan `BINARY_DOUBLE`/`BINARY_FLOAT`. Aplikasi mengirim string desimal bertitik; **jangan** mereplikasi konversi titik→koma dari `InsertIntoFacinSPreadLife_Sql` tanpa keputusan bisnis. | B4, B6 |
| **JSON API** | Semua nilai uang, TSI, rate, persen, dan limit sebagai **string desimal**, bukan `number` JSON. `number` JSON adalah IEEE-754 double di sebagian besar parser → mengulang B6. | B1, B6 |
| **JSON API — format** | Titik sebagai pemisah desimal, tanpa pemisah ribuan, tanpa tanda kutip di dalam angka. Format tampilan (koma Indonesia) adalah urusan lapisan presentasi. | B1, B4 |
| **Uji rekonsiliasi** | Setiap rumus yang diport punya uji yang membandingkan hasil terhadap keluaran Pega pada angka nyata, termasuk kasus batas 3e9 / 6e9 / 30e9 yang memicu perbedaan string-vs-numerik. | B2 |

---

## 8. Pertanyaan terbuka

### 8.1 MEMBLOKIR implementasi

| # | Pertanyaan | Bukti / lokasi |
| --- | --- | --- |
| **M1** | **Satuan `.Rate`.** `CountPremi_ACT` membagi `TSI*Rate*Prorate*Indemnity*LossLimit` dengan 1e9 (konsisten dengan Rate per mille); `CountProrateExtension_Act` step 7 membagi `TSI*Rate*Prorate` dengan 1e4 (konsisten dengan Rate persen). Satu ordo besaran 10 berbeda. | §1.2 vs §1.5 |
| **M2** | **Skala hasil `Tree_ShortPeriod`.** Cabang `GroupPanel="007"` mengembalikan 12,5…100; cabang `"002"` mengembalikan 0,125…1. Pemanggil memperlakukan keduanya sama. | §1.6 |
| **M3** | **`Tree_ShortPeriod` cabang `"002"` untuk `dateDifferent` 8–44.** Baris memakai `= 45` (kesamaan), bukan `<= 45` → rentang 8–44 tidak punya hasil. | §1.6 |
| **M4** | **`Tree_ShortPeriod` untuk `GroupPanel` selain `"002"`/`"007"`.** Tidak ada cabang default. Nilai balik tidak diketahui. | §1.6 |
| **M5** | **Tanda PPh/PPN dalam net payable.** Empat konvensi berbeda hidup berdampingan: `+PPh +PPN`, `−PPh −PPN`, `−PPh +PPN`, dan varian dengan/tanpa `Deduction2`. | §2.6 |
| **M6** | **Presisi default `@Math.divide` dua argumen.** Dipakai di `CalculateNetRate_ACT` untuk `ShareOffered`, `PercentOffered`, `NominalShareOffered` — semuanya nilai uang. 33 berkas memakai pembagi 100 tanpa argumen presisi. | §2.2, §1.7 |
| **M7** | **Asal `LimitAkseptasi.pxResults(n).CARI10`.** Dibandingkan dengan `PositionNote` untuk menentukan approver berikutnya, tetapi tidak pernah ditulis di korpus dan tidak dialiaskan di SQL mana pun (SQL hanya CARI1–CARI5). | §4.4 |
| **M8** | **Ejaan nilai kolom `M_LIMIT_*.JABATAN` di Oracle.** Ruang nama `Reas*` (seperti `PositionNote`) atau HURUF BESAR (seperti `LetterNo`)? Tidak ada baris data di korpus. | §4.4 |
| **M9** | **Penulis 12 dari 14 nilai `LetterNo`.** `SetLetterNo` hanya menulis `JUW_A` dan `UNDERWRITER`. `JUW_B` tidak pernah ditugaskan sama sekali. | §4.3 |
| **M10** | **`IsEdmAdjCeding` dan `IsEdmPPNPPH` punya kondisi identik.** Dua jenis endorsement tidak dapat dibedakan. | §3.2 |
| **M11** | **Baris hasil 12 DecisionTable tidak terekspor** — khususnya `IsUWAccepted` (menentukan penerimaan UW) dan `isApproved` (menentukan persetujuan treaty). | §3.9 |
| **M12** | **Perbandingan ambang uang sebagai string** (3e9, 6e9, 30e9). Ambang yang sama dibandingkan numerik di rule lain. Menentukan tangga akseptasi dan proteksi spreading. | §7.1 B2 |
| **M13** | **Arah konversi pemisah desimal berlawanan** antara `InsertIntoFacinSPreadLife_Sql` (titik→koma) dan `InsertIntoFacinSPreadLifeMonthly_Sql` (koma→titik) untuk kolom yang sama. | §7.1 B4 |
| **M14** | **`IsGroup` melewati seluruh spreading** untuk 3 identitas operator + 1 kode kontak marketing. Aturan bisnis sebenarnya tidak diketahui. | §6 |
| **M15** | **`ProRateType != 3` menihilkan `PremiLifeNusantaraRe`.** Tiga dari empat nilai (`1`,`2`,`4`) memicunya. | §3.6 |
| **M16** | **`@LengthOfPageList(LimitAkseptasi.pxResults)=Local.Index`** memakai satu `=` di tengah ekspresi kondisi. Perbandingan atau penugasan? | §5.3 |
| **M17** | **Dua properti `StatusBusiness`.** `Quotation.StatusBusiness` memilih cabang rumus premi; `OfferFacIn.QuotationData.StatusBusiness` menggerbangi keluarga `IsEdm*`. Sinkron atau tidak? | §3.1 |
| **M18** | **Sumber pengganti `DATAPEGA.pc_History_ASM_FW_GISFW_Work`.** Dibaca oleh `GetHistoryAccPega_SQL` untuk keputusan alur; tabel internal Pega, hilang bersama Pega. | §3.8 |
| **M19** | **`EmailTypeQuotation == "2"` berarti DECLINE** sementara `2` pada semua `EmailType*` lain berarti REJECT. | §3.4 |

### 8.2 Tidak memblokir, tetapi harus diputuskan bisnis

| # | Pertanyaan | Lokasi |
| --- | --- | --- |
| T1 | Prorate tertukar antara `.Premium` dan `.PremiRp` di cabang endorsement `HitungPremi_FacInDT`. | §1.1 |
| T2 | `FirstLossScale` dijaga `@if(=="",100,…)` di `.Premium` tetapi tidak di `.PremiRp`. | §1.1 |
| T3 | `pyWorkPage.Policy.Payment.Premium` diset di luar loop mata uang → menyimpan mata uang terakhir. | §1.3 |
| T4 | `SetEdmType` menaruh opsi `Type=12` ke `TypePenambahan` pada cabang `EdmTypeNew==2`. | §3.2 |
| T5 | Cabang EDM `MANAGERTEKNIK` di `SetBanding_ACT` menyetel ticket `DivHeadUWOffer`, bukan `DivHeadUWPolicy`. | §4.3 |
| T6 | `GetLimitAkseptasiBanding_SQL` memfilter `LIMIT_BOTTOM2` tetapi `ORDER BY LIMIT_BOTTOM`. | §4.4 |
| T7 | Kualifikasi skema `POOLDATA.` tidak konsisten antar varian tabel limit yang sama. | §4.4 |
| T8 | `.PercentageAdjustment!="100" \|\| .PercentageAdjustment!="75"` adalah tautologi. | §2.7 |
| T9 | `IsCedingConfirm` memakai `"Accepted"` dan `"accepted"` (beda kapitalisasi). | §3.7 |
| T10 | Empat kondisi workbasket memakai literal tanpa tanda kutip (`= ReasFacInDirector`). | §5.1 |
| T11 | `@divide(.LostLimit,100,1)` membulatkan LOL ke 1 desimal. | §2.7 |
| T12 | Pembulatan ulang ke 4 desimal di `SetTSIPremiSpreaded_FacIn` hanya untuk PA/MBU/Fire-non-layer; Golf/Engineering/Aneka/Fire-layer tetap 20 desimal. | §2.4 |
| T13 | `Param.Spreaded` tidak direset per coverage di cabang MBU/Fire/Golf/Engineering. | §2.4 |
| T14 | Pra-pembulatan sebelum `@notEqual(Param.Spreaded,Param.Premi)` hanya ada di cabang PA. | §2.4 |
| T15 | `.Sublimit` vs `Local.SubLimit` di berkas yang sama. | §1.2 |
| T16 | Celah enumerasi: `Type=10`, `StatusBusiness` 0/≥4, `EdmType` 0/≥5, `EmailType*` 0/5/6/8, `ProposalAcceptStatus` 0/5–8, `ProRateType` 0, `metodeKalkulasi` 0/≥4. | §3 |
| T17 | 669 kemunculan / 23 label status yang memuat nama orang di `NBStatus`/`NBStatusNew`. | §6 |
| T18 | `EmailTypeCeding` nilai `5` dan `6` tidak menghasilkan baris `HISTORYAKSEPTASIPEGA`. | §3.4 |
| T19 | `ProposalAcceptStatus` ditulis sebagai angka `4` tetapi dibaca sebagai string `"4"` di `SetBanding_ACT`. | §3.5 |
| T20 | Enam token ticket (`AdminOffer`, `DivHeadFinPolicy`, `HeadRetroPolicy`, `TLPolicy`, `Offer`, `Policy`) tidak disetel di `SetBanding_ACT`. | §5.2 |

---

## Lampiran — skrip audit yang dipakai

Disimpan di direktori scratchpad sesi; disalin di sini agar setiap angka dapat direproduksi.

**`scan.ps1`** — cari pola beserta konteks di korpus:

```powershell
param([string]$Pattern, [string]$Sub = "", [int]$Pad = 120, [string]$Filter = "*.xml")
$root = "D:\migrasi\RNM"
$dirs = @("$root\NB FacIn","$root\RNW Fac In","$root\Endorsment Fac In")
if ($Sub -ne "") { $dirs = $dirs | ForEach-Object { Join-Path $_ $Sub } }
$out = New-Object System.Collections.Generic.List[string]
foreach ($d in $dirs) {
  if (-not (Test-Path $d)) { continue }
  foreach ($f in [IO.Directory]::EnumerateFiles($d, $Filter, 'AllDirectories')) {
    $c = [IO.File]::ReadAllText($f); if ($c -notmatch $Pattern) { continue }
    foreach ($m in [regex]::Matches($c, "[^<>]{0,$Pad}$Pattern[^<>]{0,$Pad}")) {
      $out.Add(("{0} :: {1}" -f $f.Replace("$root\",""), $m.Value)) } } }
$out | Sort-Object -Unique
```

**`lit.ps1`** — hitung literal distinct + jumlah berkas:

```powershell
param([string]$Pattern)
$root = "D:\migrasi\RNM"; $tally = @{}
foreach ($corp in @("NB FacIn","RNW Fac In","Endorsment Fac In")) {
  foreach ($f in [IO.Directory]::EnumerateFiles("$root\$corp", "*.xml", 'AllDirectories')) {
    $c = [IO.File]::ReadAllText($f)
    foreach ($m in [regex]::Matches($c, $Pattern)) {
      $v = $m.Groups[1].Value
      if (-not $tally.ContainsKey($v)) { $tally[$v] = New-Object System.Collections.Generic.HashSet[string] }
      [void]$tally[$v].Add($f.Replace("$root\","")) } } }
foreach ($k in ($tally.Keys | Sort-Object)) { "{0}`tfiles={1}" -f $k, $tally[$k].Count }
"TOTAL DISTINCT = " + $tally.Keys.Count
```

**`act.ps1`** — baca langkah rule `Activity` berurutan:

```powershell
param([string]$Path)
$c = Get-Content $Path -Raw -Encoding UTF8
$tags = 'pyStepPageReference|pyStepsActivityName|pyStepsObjectName|pyStepsDescription|pyStepsPreCondition|pyStepsTransition|PropertiesName|PropertiesValue|pyStepsJavaSource|pyStepsPreCondParamsWhen|pyStepsTransParamsWhen|pyForEachProperty|pyStepsRepeatDefHasRepeat|Param|Value|pyStepsBlockName'
$buf = $null; $pairName = $null
foreach ($m in [regex]::Matches($c, "<($tags)>(.*?)</\1>", 'Singleline')) {
  $t = $m.Groups[1].Value; $v = $m.Groups[2].Value
  if ($v.Trim() -eq '') { continue }
  if ($t -eq 'pyStepPageReference') { if ($buf) { $buf }; $buf = "`n### $v" }
  elseif ($buf) {
    if ($t -eq 'PropertiesName' -or $t -eq 'Param') { $pairName = $v }
    elseif ($t -eq 'PropertiesValue' -or $t -eq 'Value') { $buf += "`n    SET  $pairName  =  $v"; $pairName = $null }
    else { $buf += "`n    $t : $v" } } }
if ($buf) { $buf }
```

**`dt.ps1`** — baca langkah rule `DataTransform` berurutan:

```powershell
param([string]$Path)
$c = Get-Content $Path -Raw -Encoding UTF8
$tags = 'pyPropertyStepId|pyActionName|pyPropertiesName|pyPropertiesValue|pyStepComments'
$cur = $null
foreach ($m in [regex]::Matches($c, "<($tags)>(.*?)</\1>", 'Singleline')) {
  $t = $m.Groups[1].Value; $v = $m.Groups[2].Value
  if ($t -eq 'pyPropertyStepId') { if ($cur) { $cur }; $cur = "`n--- STEP $v" }
  elseif ($cur -and $v.Trim() -ne '') { $cur += "`n    $t = $v" } }
if ($cur) { $cur }
```

**`when.ps1`** — katalog kondisi rule `When` (menghasilkan 226 entri):

```powershell
$root = "D:\migrasi\RNM"; $seen = @{}
foreach ($corp in @("NB FacIn","RNW Fac In","Endorsment Fac In")) {
  foreach ($f in [IO.Directory]::EnumerateFiles("$root\$corp\When", "*.xml")) {
    $name = [IO.Path]::GetFileNameWithoutExtension($f); $c = [IO.File]::ReadAllText($f)
    $conds = [regex]::Matches($c,'<pyConditionValue1String>(.*?)</pyConditionValue1String>','Singleline') | %{ $_.Groups[1].Value.Trim() } | ? { $_ -ne '' }
    $cs    = [regex]::Matches($c,'<pyConditionString>(.*?)</pyConditionString>','Singleline') | %{ $_.Groups[1].Value.Trim() } | ? { $_ -ne '' -and $_ -ne '[Double click to add condition]' }
    $logic = [regex]::Matches($c,'<pyLogic>(.*?)</pyLogic>','Singleline') | %{ $_.Groups[1].Value.Trim() } | ? { $_ -ne '' } | Sort-Object { $_.Length } -Descending | Select-Object -First 1
    if (-not $seen.ContainsKey($name)) { $seen[$name] = @() }
    $seen[$name] += ("LOGIC={0}`t{1}" -f $logic, (($conds + $cs | Sort-Object -Unique) -join ' | ')) } }
foreach ($k in ($seen.Keys | Sort-Object)) { "=== $k"; $seen[$k] | Sort-Object -Unique | %{ "    $_" } }
```

**`sql.ps1`** — dump SQL dari rule `RDBList`:

```powershell
param([string]$Glob = "*.xml", [string]$Corp = "NB FacIn")
foreach ($f in [IO.Directory]::EnumerateFiles("D:\migrasi\RNM\$Corp\RDBList", $Glob)) {
  $c = [IO.File]::ReadAllText($f)
  $m = [regex]::Match($c, '<pyBrowseSQL>(.*?)</pyBrowseSQL>', 'Singleline')
  if ($m.Success -and $m.Groups[1].Value.Trim() -ne '') {
    "===== " + [IO.Path]::GetFileName($f)
    ($m.Groups[1].Value -replace '&lt;','<' -replace '&gt;','>' -replace '&amp;','&') } }
```
