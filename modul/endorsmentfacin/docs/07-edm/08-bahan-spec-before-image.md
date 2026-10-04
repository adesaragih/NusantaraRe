# Bahan spec — Modul Before-Image EDM (mekanisme "nilai lama")

> **Ini BAHAN untuk `/to-spec`, bukan spec.** `/to-spec` ber-`disable-model-invocation: true` dan
> **belum dijalankan**. Work owner yang menjalankannya manual.
>
> **Sumber:** `D:\migrasi\RNM\Endorsment Fac In\` (banding `NB FacIn\`). READ-ONLY.
> Korpus Treaty `RNM_BRD\` **tidak dibaca** (K-005). Label mengikuti `CLAUDE.md` §3.
>
> **Tidak ada discovery ulang.** Bahan dirapikan dari `07-edm\04-e3-before-image-selisih.md` dan
> arsip `01-flow\03-alur-endorsement.md` §2–§3, dengan angka **dihitung ulang dari korpus** agar
> tidak ada perkiraan yang lolos ke spec.
>
> Keputusan yang mengikat bahan ini: **K-010/K-012** (uang = `Money`, tidak pernah `float`) ·
> **K-018** (skala rasio) · **K-027** (koma desimal) · **K-043** (lingkup 708) ·
> **K-044** (Life ikut; pilih-tertanggung dibuang) · **K-045** · **K-046** (kode usang diport apa
> adanya, kecuali A.5).

---

## 0. Mengapa modul ini didahulukan

`[terverifikasi]` Nilai dasar akseptasi endorsement adalah **selisih**, bukan nilai penuh. Selisih
tidak dapat dihitung tanpa "nilai lama". Karena itu before-image adalah **fondasi** — modul premi,
tangga akseptasi, dan jalur produksi endorsement semuanya bersandar padanya.

⛔ Modul ini **tidak memblokir**: seluruh mekanismenya terbaca dari korpus. Yang ditunda adalah
**jalur produksi** (menunggu tabel flat + `ALL_SOURCE` dari DBA) — bukan bagian modul ini.

⚠️ **Istilah.** `[terverifikasi]` Token `BeforeImage` / `AfterImage` **tidak ada di korpus EDM**
(0 berkas). "Before-image" adalah istilah konsep di dokumen ini, **bukan** nama properti. Jangan
menulis kode yang mengasumsikan properti bernama itu.

---

## 1. Kontrak data — tiga lapis yang harus dipisahkan

`[terverifikasi]` Korpus memuat **tiga mekanisme berbeda** yang sering tertukar. Ketiganya punya
struktur, pengisi, dan **siklus hidup** berbeda.

| | **Lapis A** | **Lapis B** | **Lapis C** |
| --- | --- | --- | --- |
| **Apa** | Dokumen polis versi terakhir, **utuh** | Nilai lama **per baris** | **Penanda** baris warisan |
| **Disimpan di** | `OfferFacIn.OldData` (page, kelas sama dengan `OfferFacIn`) | properti bersaudara `*Old` **di dalam list yang sekarang** | `.IsOldData` di tiap simpul list |
| **Diisi oleh** | `Activity\SetValueToEDMWork` langkah **14** | `Activity\SetOldData` | 7× `Activity\SetOLDValueToEDMWork_<LOB>` |
| **Dipanggil dari** | — (bagian pembuatan case) | `Activity\InputAddendumFacIn_PreAct` langkah **28** | `SetValueToEDMWork` langkah **14.7–14.13** |
| **Kapan** | **sekali**, saat case EDM lahir | **setiap kali layar endorsement dibuka** | **sekali**, saat case EDM lahir |
| **Sumber** | tabel `JSON_POLIS` (kolom `DATA_JSON`) | page `OfferFacIn.OldData` (lapis A) | page `OfferFacIn.OldData` (lapis A) |
| **Dipakai untuk** | **perhitungan SELISIH produksi** | layar + perhitungan per baris | menandai baris asal polis lama |

⛔ **Lapis A dan lapis B melayani konsumen berbeda.** Delta produksi (`facinproduction`) membaca
**lapis A**; layar dan perhitungan per baris membaca **lapis B**. Menyatukan keduanya akan mengubah
angka.

⛔ **Riwayat versi polis bukan salah satu dari ketiganya.** Riwayat ada di tabel `JSON_POLIS` sebagai
**baris bertambah** dengan kolom `PRODKE` sebagai nomor versi (0-based) — `CLAUDE.md` §4.3. Lapis A
*membaca* riwayat itu; lapis B dan C adalah state di dalam case yang sedang berjalan.

### 1.1 Lapis A — sumber datanya

`[terverifikasi]` `Endorsment Fac In\RDBList\GetEDMOldData_SQL.xml` ·
`pxObjClass` = **`Rule-Connect-SQL`** ⚠️ (bukan tipe RDBList meski berada di folder `RDBList\`) ·
`pyClassName` = `ASM-FW-GISFW-Int-OFFERJSON`:

```sql
SELECT a.DATA_JSON AS HASIL1 FROM JSON_POLIS a
 WHERE NOPOLIS={newWorkPage.OfferFacIn.QuotationData.OldPolicyNo}
   and PRODKE=(SELECT COUNT(NOPOLIS)-1 FROM JSON_POLIS a
                WHERE NOPOLIS={newWorkPage.OfferFacIn.QuotationData.OldPolicyNo})
```

Rantai langkah 14 `[terverifikasi]` — **13 sub-langkah**:

| Sub | Isi |
| --- | --- |
| 14.1 | `RDB-List` · `RequestType` = `GetEDMOldData_SQL` · `ClassName` = `ASM-FW-GISFW-Int-OFFERJSON` · `BrowsePage` = **`OldData`** |
| 14.2 | `Java` — adopsi JSON ke page (*"copy NB and remove the comments, and fac retro list"*) |
| **14.3** | `Property-Set` — **54 pasangan** (§2) |
| 14.4–14.5 | ambil ulang `TeamGroup` dari master marketing |
| 14.6 | `Obj-Save` |
| 14.7–14.13 | **7× `Call SetOLDValueToEDMWork_<LOB>`** → lapis C |

⚠️ `pyStepsPreCondition` langkah **14** dan **14.1** = `false`. Semantiknya sudah ditutup lewat
keputusan work owner; dicatat di sini agar terlihat di spec.

---

## 2. Lapis A — daftar field langkah 14.3 (dihitung, bukan diperkirakan)

`[terverifikasi]` Angka persis dari korpus:

| Ukuran | Jumlah |
| --- | ---: |
| Pasangan `Property-Set` di 14.3 | **54** |
| — bersumber `OldData.*` | **51** |
| — **bukan** dari `OldData` | **3** |
| Dari yang 51: mengikuti pola `OfferFacIn.X = OldData.X` | **45** |
| Dari yang 51: **menyimpang** (tulis ke halaman lain) | **6** |

> 🔁 **Koreksi arsip.** `03-alur-endorsement.md` §2.1.1 menyatakan *"54 `Property-Set`, **semuanya**
> berbentuk `newWorkPage.OfferFacIn.X = newWorkPage.OfferFacIn.OldData.X`"*. Itu **tidak berlaku
> untuk 9 dari 54** — 3 tidak bersumber dari `OldData`, dan 6 menulis ke halaman selain `OfferFacIn`.

### 2.1 Empat kelompok field (45 yang berpola `OfferFacIn.X = OldData.X`)

`[terverifikasi]` Nama properti dikutip apa adanya.

**Identitas bisnis — 20 field** (`QuotationData.*`)

```
BusinessOldId · BusinessCode · BusinessFac · BusinessName · BusinessType
CedingCo · CedingCoName · SourceOfBusiness · SobName
MarketingCode · MarketingName · MOID · TeamGroup
InsuredName · InsuredID · NoOfferSlip · IsGroup · QQName · PolicyType · EDMDay
```

**Periode & tanggal — 4 field** (`PolicyData.*`)

```
StartDateTime · EndDateTime · OfferingDate · ProdDateTime
```

**Struktur share & kapasitas — 18 field** (`OfferFacIn.*`)

```
Currency · CurrencyList · Parameters · InwardScale · PercentShare
AdditionalCapital · MaxPctTreatyCapacity · MaxTreatyCapacity · CurrentYear
ProRatePercent · ProRateType · CedingCedantList · ShareCedantType
IsSpecialAcceptance · BinderRNM · PPnCheck · PolicyMasterNumber · IsB2B
```

**Pembayaran — 3 field** (`PolicyData.Payment.*`)

```
Installment · RICommision · PctBrokerageFee
```

### 2.2 Enam penyimpangan — cermin ke halaman lain

`[terverifikasi]` Ini **bukan** salinan `OfferFacIn.X = OldData.X`; targetnya halaman berbeda.
Spec harus menyebutkan keenamnya eksplisit, kalau tidak halaman-halaman ini akan kosong:

| Baris | Target | Sumber |
| ---: | --- | --- |
| L3135 | `newWorkPage.Quotation.BusinessType` | `OldData.QuotationData.BusinessType` |
| L3490 | `newWorkPage.Policy.Payment.Installment` | `OldData.PolicyData.Payment.Installment` |
| L3531 | `newWorkPage.Policy.Payment.RICommision` | `OldData.PolicyData.Payment.RICommision` |
| L3572 | `newWorkPage.Policy.Payment.PctBrokerageFee` | `OldData.PolicyData.Payment.PctBrokerageFee` |
| L3793 | `newWorkPage.IsB2B` *(akar work page)* | `OldData.IsB2B` |
| L3813 | `InputDataCredit.CARI6` | `OldData.QuotationData.MarketingName` |

📌 Tiga field pembayaran ditulis **dua kali** — ke `OfferFacIn.PolicyData.Payment.*` **dan** ke
`Policy.Payment.*`. Itu sebabnya 45 + 6 = 51, bukan 51 field berbeda.

### 2.3 Tiga yang tidak bersumber dari `OldData`

`[terverifikasi]`

| Baris | Target | Sumber |
| ---: | --- | --- |
| L3632 | `pyWorkPage.Quotation.BusinessType` | `newWorkPage.Quotation.BusinessType` |
| L3733 | `newWorkPage.Quotation` | `newWorkPage.OfferFacIn.QuotationData` |
| L3914 | `newWorkPage.PPnCheck` | `newWorkPage.OfferFacIn.PPnCheck` |

### 2.4 ⛔ Akibat yang mengikat desain

`[terverifikasi]` Setelah 14.3, `OfferFacIn` yang "baru" **identik** dengan polis lama. Pengguna lalu
hanya mengubah yang perlu.

> **Karena itu delta dihitung terhadap `OldData`, bukan terhadap "kosong".**

Implementasi Go **tidak boleh** memulai dari struct kosong lalu mengisi perubahan; ia harus
**menyalin lapis A lebih dulu**, persis 54 penugasan di atas, baru menerima perubahan pengguna.

---

## 3. Lapis B — properti `*Old` per lini bisnis

`[terverifikasi]` `Activity\SetOldData.xml` · `pxObjClass` = `Rule-Obj-Activity` ·
`pyClassName` = `ASM-FW-GISFW-Work`. Blok **4** punya **14 sub-langkah** berpasangan
(objek + cedant) untuk **7 lini**.

| Sub | Lini (`pyStepsPreCondParamsWhen`) | Properti `*Old` yang diisi (jalur dari `pyWorkPage.OfferFacIn`) |
| --- | --- | --- |
| **4.1** | `IsFire` | `LocationList(i).Property.PropertyItemList(j).{TSIObjectItemOld, TotalGrossPremiOld, TotalPremiumNusantaraReOld}` · `…Property.TotalTSIList(j).TSIOld` · `…Property.TotalTSIPremiGrossList(k).{TSIOld, PremiumOld, RateOld}` |
| **4.3** | `IsGolfInsurance` | `LocationList(i).Property.RiskLocation.AnekaList(j).TSIOld` · `…AnekaList(j).CoverageList(k).{TSIOld, PremiumOld}` · `…TotalTSIPremiGrossList(j).{TSIOld, PremiumOld, RateOld}` |
| **4.5** | `IsAneka` | `LocationList(i).Property.RiskLocation.OccupationList(j).AnekaList(k).TSIOld` · `…TotalTSIPremiGrossList(j).{TSIOld, PremiumOld, RateOld}` |
| **4.7** | `IsPA` | `PersonList(i).ASMCoverage(j).{TSIOld, PremiumOld}` · `PersonList(i).TotalTSIPremiGrossList(j).{TSIOld, PremiumOld, RateOld}` |
| **4.9** | `IsMarineCargo` | `CargoList(i).CoverageList(j).{TSIOld, PremiumOld, **PremiNusantaraReOld**}` · `CargoList(i).TotalTSIPremiGrossList(j).{TSIOld, PremiumOld, RateOld}` |
| **4.11** | `IsMBU` | `VehicleList(i).CoverageList(j).{TSIOld, PremiumOld, **PremiNusantaraReOld**, PremiRpOld, PremiumGrossDiscountFleetOld}` · `VehicleList(i).TotalTSIPremiGrossList(j).{TSIOld, PremiumOld, RateOld}` |
| **4.13** | `IsTravel` | `PersonList(i).ASMCoverage(j).{TSIOld, PremiumOld}` |
| **4.2/4/6/8/10/12/14** | masing-masing lini | `CedingCedantList(i).CurrencyList(j).{TSIOld, PremiumOld}` — **identik di ketujuh lini** |

`[terverifikasi]` **52 penugasan `*Old`** seluruhnya, dari **89 penugasan** total di `SetOldData`.

📌 **`PremiNusantaraReOld` hanya ada di dua lini**: Marine Cargo (4.9) dan MBU (4.11). Bukan di lima
lini lain.

### 3.1 Pola nilai — kosong menjadi nol

`[terverifikasi]` **52 dari 52** memakai `@If(...)`. **Nol pengecualian.** Bentuk bakunya:

```
TSIOld                     = @If(.TSI!="",                     .TSI,                     0)
PremiumOld                 = @If(.Premium!="",                 .Premium,                 0)
TotalGrossPremiOld         = @If(.TotalGrossPremi!="",         .TotalGrossPremi,         0)
TotalPremiumNusantaraReOld = @If(.TotalPremiumNusantaraRe!="", .TotalPremiumNusantaraRe, 0)
TSIObjectItemOld           = @If(.TSIObjectItem!="",           .TSIObjectItem,           0)
PremiRpOld                 = @If(.PremiRp!="",                 .PremiRp,                 0)
PremiumGrossDiscountFleetOld = @If(.PremiumGrossDiscountFleet!="", .PremiumGrossDiscountFleet, 0)
```

⛔ **Kontrak: nilai kosong dipetakan ke `0`, bukan NULL.** Di Go, `*Old` yang tidak terisi bernilai
nol, bukan pointer nil — kalau tidak, aritmetika selisih akan berbeda.

### 3.2 ⚠️ A.5 — satu-satunya perbaikan sadar

`[terverifikasi]` `PremiNusantaraReOld` memakai fallback **berkutip**:

```
L6294  CargoList(local.idxcargo).CoverageList(local.idx).PremiNusantaraReOld
         = @If(.PremiNusantaraRe!="",.PremiNusantaraRe,"0")      ← STRING "0"
L7527  VehicleList(Local.idxvehicle).CoverageList(Local.idx).PremiNusantaraReOld
         = @If(.PremiNusantaraRe!="",.PremiNusantaraRe,"0")      ← STRING "0"
```

**2 penugasan berkutip · 50 penugasan numerik · total 52.**

⛔ **K-046 A.5 — DI SISTEM BARU PAKAI ANGKA `0`, BUKAN STRING `"0"`.**

> **Ini PERBAIKAN, bukan port.** Satu-satunya di modul ini. `PremiNusantaraRe` adalah **nilai uang**
> (K-010/K-012) sehingga tipenya `Money`, dan `Money` tidak boleh menampung string.

⚠️ **Catatan rekonsiliasi paralel run — wajib masuk spec:** pada kasus `PremiNusantaraRe` kosong di
Marine Cargo atau MBU, sistem lama menyimpan **string** `"0"` dan sistem baru menyimpan **angka** `0`.
Perbandingan mentah bisa tampak berbeda **tipe** meski nilainya sama. Pembanding rekonsiliasi harus
menormalkan sebelum membandingkan, dan selisih semacam ini **dijelaskan oleh A.5**, bukan bug.

📌 `[terverifikasi]` Tambahan yang memperkuat A.5: pada `rowdata` yang sama, cache editor
`<pyExpressionGadget>` (L6297/L7530) justru memuat bentuk **tanpa kutip**. Yang dieksekusi adalah
`<PropertiesValue>`; cache bukan bukti.

### 3.3 Yang **tetap** diport apa adanya (K-046)

| Kejanggalan | Bukti | Sikap |
| --- | --- | --- |
| **`RateOld` guard self-referential** — `@If(.RateOld!="", .Rate, 0)` pada **6 dari 6** kemunculan (4.1, 4.3, 4.5, 4.7, 4.9, 4.11) | korpus | ⛔ **DIPORT APA ADANYA** — work owner menetapkan ini **perilaku yang benar**, bukan bug |
| **`PremiumOld` self-referential di FIRE saja** — `@If(.PremiumOld!="", .Premium, 0)` di 4.1; 17 kemunculan lain memakai `@If(.Premium!="", .Premium, 0)` | korpus | ⛔ **DIPORT APA ADANYA** |

⛔ Implementasi Go **wajib menyalin guard apa adanya**, termasuk yang menguji properti tujuan.
Meluruskannya akan mengubah angka dan memutus rekonsiliasi (`CLAUDE.md` §1).

### 3.4 Tipe data yang diusulkan (turunan K-010/K-012/K-018)

⚠️ Ini **penerapan keputusan yang sudah ada**, bukan keputusan baru.

| Properti | Tipe | Alasan |
| --- | --- | --- |
| `TSIOld` · `TSIObjectItemOld` · `PremiumOld` · `TotalGrossPremiOld` · `TotalPremiumNusantaraReOld` · `PremiNusantaraReOld` · `PremiRpOld` · `PremiumGrossDiscountFleetOld` | **`Money{Amount decimal, Currency string}`** | nilai uang — K-010/K-012, tidak pernah `float` |
| `RateOld` | **`Ratio{Value decimal, Scale}`** | skala per lini mengikuti **K-018** (PA/Layering/FIRE = ‰ ÷1.000; MBU/ANEKA/BONDING/GOLF/MARINE CARGO = % ÷100) |
| `.IsOldData` (lapis C) | `string` — nilai literal `"old"` | penanda, bukan angka |

⛔ `Money` dan `Ratio` **tidak dapat dijumlahkan**. Satu-satunya jembatan: **`Money × Ratio → Money`**.

---

## 4. Lapis C — penanda `.IsOldData`

`[terverifikasi]` Tujuh varian `Activity\SetOLDValueToEDMWork_<LOB>`, **ketujuhnya ADA**:
`_FIRE` · `_MC` · `_Aneka` · `_MBU` · `_LIFE` · `_PA` · `_GOLF`.

Nilai tersimpan **berkutip**: `<PropertiesValue>"old"</PropertiesValue>`.
⚠️ Pencarian `>old<` memberi **0 hasil** — pola cari wajib menyertakan kutip.

| Varian | Cacah `IsOldData` | `IsProRate` | Kedalaman list |
| --- | ---: | ---: | --- |
| `_Aneka` | 14 | 0 | LocationList → OccupationList → AnekaList → CoverageList |
| `_FIRE` | 12 | **2** | LocationList → PropertyItemList → CoverageList → LayerList |
| `_GOLF` | 12 | 0 | LocationList → AnekaList → CoverageList |
| `_MBU` | 9 | 0 | VehicleList → CoverageList |
| `_MC` | 8 | 0 | CargoList → CoverageList |
| `_PA` | 8 | 0 | PersonList → ASMCoverage |
| `_LIFE` | **4** | 0 | **PersonList saja** |

📌 Cacah penanda mengikuti **kedalaman model data**, bukan kelengkapan implementasi. `_LIFE` paling
sedikit karena Life hanya punya satu tingkat list.

⛔ **Hanya `_FIRE` menyetel `newWorkPage.OfferFacIn.IsProRate = "Prorate"`** (L416). Enam varian lain
0×. `IsProRate` adalah properti **tingkat `OfferFacIn`** (menandai seluruh case), dan **tidak**
termasuk 54 field yang disalin lapis A.

⛔ **`_LIFE` menyetel DUA penanda** — `.FlagOldData = "old"` **dan** `.IsOldData = "old"`; enam
varian lain hanya `.IsOldData`.

---

## 5. Kontrak SELISIH

### 5.1 Prorata — dua rasio

`[terverifikasi]` `SetValueToEDMWork` langkah **15** (*"Count StartProRate and EndProRate"*):

```
SET Local.EdmToStart  = QuotationData.EdmDate          − OldData.PolicyData.StartDateTime
SET Local.EdmToEnd    = OldData.PolicyData.EndDateTime − QuotationData.EdmDate
SET Local.TotalPeriod = OldData.PolicyData.EndDateTime − OldData.PolicyData.StartDateTime
SET Local.TotalPeriod = @if(Local.TotalPeriod = 0, 1, Local.TotalPeriod)        ← guard bagi-nol

L5340  SET .OfferFacIn.ProrateStartEDM = @Math.divide(Local.EdmToStart, Local.TotalPeriod, 20)
L5367  SET .OfferFacIn.ProrateEDMEnd   = @Math.divide(Local.EdmToEnd,   Local.TotalPeriod, 20)
```

| Rasio | Arti | Dipakai |
| --- | --- | --- |
| `ProrateStartEDM` | porsi periode **sebelum** tanggal endorsement | hanya bila `QuotationData.Type=="3"` (Adj Rate), lewat `Local.ProrateStart` |
| `ProrateEDMEnd` | porsi periode **sesudah** tanggal endorsement | **pengali seluruh nilai baru** — 19 penugasan di `SaveFacinProdEDMFire_Act` |

⛔ **Skala 20 desimal** → tipe desimal presisi tinggi, **bukan `float`**. Keduanya **`Ratio`**, bukan
`Money`: mereka **mengalikan** Money, tidak pernah dijumlahkan dengannya. Guard `TotalPeriod=0 → 1`
**wajib diport**.

### 5.2 Rumus kanonik

`[terverifikasi]` `Activity\SaveFacinProdEDMFire_Act.xml` — **57 penugasan delta berbentuk
pengurangan**:

```
# 1) nilai baru diprorata lebih dulu
SET Local.TsiSpreadEDM   = .TSISpreaded     × pyWorkPage.OfferFacIn.ProrateEDMEnd
SET Local.PremiSpreadEDM = .PremiumSpreaded × pyWorkPage.OfferFacIn.ProrateEDMEnd

# 2) lalu dikurangi nilai lama dengan TreatyType yang SAMA
SET Local.TsiSpreadEDM   = Local.TsiSpreadEDM   − Local.TsiSpreadNB
SET Local.PremiSpreadEDM = Local.PremiSpreadEDM − Local.PremiSpreadNB
SET Datain1.CARI11       = Local.NewRIComm        − Local.OldRIComm
SET Datain1.CARI15       = Local.NewBrokerageFree − Local.OldBrokerageFree
```

> ### **`SELISIH = (nilai_baru × ProrateEDMEnd) − nilai_lama_dengan_TreatyType_sama`**

⛔ **`nilai_lama` diambil dari LAPIS A (`OfferFacIn.OldData`), bukan dari `*Old` per baris (lapis B).**
Ini pembeda yang paling mudah salah diimplementasikan.

Pencocokan pasangan lama–baru dilakukan per **`TreatyType`**, bukan per indeks posisi — loop
`…4.3.2` mencocokkan `TreatyType` baru terhadap lama sebelum mengambil `Local.TsiSpreadNB` /
`Local.PremiSpreadNB`.

### 5.3 Varian yang menyimpang dari rumus kanonik

`[terverifikasi]` Diport apa adanya:

| Kondisi | Perlakuan |
| --- | --- |
| `Type=="4"` (Adjustment Spreading) | `SELISIH = MENJADI` — seluruh nilai dianggap baru |
| `Type=="7"` + mata uang berubah | idem |
| `Type=="3"` (Adj Rate) | porsi "sebelum" ikut ditambahkan: `(TotalPremiNew + PremiStart) − PremiSpreadNB` |
| coverage hilang dari data baru | `delta = nilai_lama × −1` |
| `EdmType ∈ {1,2,3}` | selisih persentase komisi/brokerage/deduction **dibalik tandanya** (`× −1`) |

### 5.4 Fondasi "nilai dasar akseptasi = selisih"

⚠️ `[belum diuji]` Arsip `03-alur-endorsement.md` §4.3 menyatakan *"nilai dasar akseptasi endorsement
adalah SELISIH TSI, bukan nilai penuh"*. Klaim itu **belum dikonfirmasi ulang ke korpus** — arsip §4
(gerbang akseptasi EDM) memang belum ditinjau di E-1…E-6.

⛔ Spec modul before-image **boleh** menyebut selisih sebagai keluarannya, tetapi **tidak boleh**
menyatakan tangga akseptasi memakainya sampai §4.3 diverifikasi. Itu pekerjaan modul berikutnya.

---

## 6. Batas dan gerbang

`[terverifikasi]` `SetOldData` punya **3 gerbang keluar** di depan blok 4 (kode transisi Pega:
`2` = lanjut, `6` = keluar activity):

| Langkah | Kondisi | WhenTrue | WhenFalse | Akibat |
| ---: | --- | ---: | ---: | --- |
| **1** | `IsLife` | **6** | 2 | Life → **keluar**, tidak mendapat lapis B |
| **2** | `IsEDM` | 2 | **6** | bukan EDM → **keluar** |
| **3** | `@Utilities.SizeOfPropertyList(.OfferFacIn.LocationList) > 100` | **6** | 2 | >100 lokasi → **keluar** |

### 6.1 Konsekuensi ke desain

**Gerbang 1 — Life (K-044: Life IKUT lingkup).**
Life **tidak memakai lapis B sama sekali**. Ia punya jalur tersendiri yang `[terverifikasi]` ada di
korpus: `Activity\SetOLDValueToEDMWork_LIFE` (lapis C, 74.841 B) · `Activity\ReCountPremiLifeEDM`
(174.841 B) · produksi lewat `Activity\SaveFacinLive_Act` (147.233 B) dan
`Activity\SaveFacinSpreadLife_Sql` (113.530 B).

⛔ **Spec harus memodelkan Life sebagai jalur terpisah** — bukan cabang `if` di dalam lapis B,
melainkan alur sendiri yang melewati lapis B.

**Gerbang 2 — non-EDM.** `IsEDM` membaca `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3`
(agregat tersimpan). ⚠️ **`IsNotEDM` BUKAN negasi `IsEDM`** — ia membaca
`pyWorkPage.Quotation.StatusBusiness != 3` (halaman aktif). Keduanya bisa benar atau salah
bersamaan. Jangan implementasikan `IsNotEDM` sebagai `!IsEDM`.

**Gerbang 3 — batas 100 lokasi.** `[terverifikasi]` Polis dengan **lebih dari 100 lokasi tidak
mendapat lapis B sama sekali** — seluruh `*Old` tetap kosong, sehingga layar dan perhitungan per
baris menampilkan nol. ⛔ Diport apa adanya. Alasan ambang 100 tetap `[pertanyaan terbuka]`.

### 6.2 Cakupan lini — lapis B ≠ lapis C

`[terverifikasi]` Himpunan lini berbeda, dan spec harus menyebutkannya eksplisit:

| | Lini yang ditangani |
| --- | --- |
| **Lapis B** (`SetOldData`) | FIRE · GOLF · ANEKA · PA · MARINE CARGO · MBU · **TRAVEL** — tanpa Life |
| **Lapis C** (`SetOLDValueToEDMWork_*`) | FIRE · MC · Aneka · MBU · **LIFE** · PA · GOLF — tanpa Travel |

---

## 7. Di luar lingkup modul ini

- **Jalur produksi** (`facinproduction`, `FACOUTPRODUCTION`, `JSON_POLIS`) — ditunda menunggu tabel
  flat + `ALL_SOURCE` dari DBA
- **Tangga akseptasi EDM** — termasuk klaim §4.3 "nilai dasar akseptasi = selisih" `[belum diuji]`
- **Pilih-tertanggung** — **K-044: DIBUANG**. Keempat berkas `[terverifikasi]` ada di korpus tetapi
  **tidak diport**: `Harness\ChooseInsured` · `Section\ChooseInsuredDtl` ·
  `ReportDefinition\BrowseAccountInsuredEDM` · `Activity\SetDataInsuredEDM_Act`
- **Modul premi endorsement** (`EDMPremiMenjadi`, 11 rumus semantik) — modul tersendiri
- **Lapis A jalur perbaikan** `Activity\GetdataOldEDMError` (dipanggil `InputAddendumFacIn_PreAct`
  langkah 5 bila `QuotationData.BusinessCode==""`) — belum dibedah langkah-per-langkah

---

## 8. Pertanyaan terbuka yang tersisa untuk modul ini

| Pertanyaan | Status | Pemilik |
| --- | --- | --- |
| **Kerapatan `PRODKE`** — `GetEDMOldData_SQL` memakai `PRODKE=(COUNT−1)`, salah bila ada lubang | ✅ **TERTUTUP** — `GetProdKeOldData_SQL` memakai `order by TGL_INPUT desc` yang tidak rentan; jalur itu yang dipakai penomoran versi | — |
| **Alasan ambang 100 lokasi** | `[pertanyaan terbuka]` — korpus tidak menjelaskan; dugaan kinerja | work owner |
| **Mengapa hanya `_FIRE` menyetel `IsProRate`** | `[pertanyaan terbuka]` | work owner |
| **Mengapa `_LIFE` menyetel dua penanda** (`.FlagOldData` + `.IsOldData`); apakah `FlagOldData` dibaca di tempat lain | `[pertanyaan terbuka]` | work owner |
| **Hanya 2 Section menampilkan `*Old`** padahal lapis B mengisi 52 properti — apakah sisanya hanya untuk perhitungan | `[pertanyaan terbuka]` | work owner + Underwriting |
| **Klaim §4.3** "nilai dasar akseptasi = selisih TSI" | `[belum diuji]` — arsip §4 belum ditinjau | pekerjaan modul berikutnya |

⚠️ Tidak satu pun memblokir penulisan spec modul ini. Keempat yang `[pertanyaan terbuka]` menyangkut
**alasan**, bukan **perilaku** — perilakunya sudah terbaca dan dapat diport.

---

## 9. Kesiapan

✅ **Bahan lengkap untuk `/to-spec`:**

| Butir | Status |
| --- | --- |
| Kontrak data 3 lapis (struktur, pengisi, siklus hidup, sumber) | §1 |
| Daftar field lapis A — **54 / 51 / 45 + 6**, dikelompokkan, dihitung dari korpus | §2 |
| Field `*Old` lapis B per 7 lini + pola `@If(...,0)` | §3 |
| **A.5 ditandai** sebagai satu-satunya perbaikan + catatan rekonsiliasi | §3.2 |
| Kode usang yang tetap diport (`RateOld`, `PremiumOld` FIRE) | §3.3 |
| Tipe `Money` / `Ratio` turunan K-010/K-012/K-018 | §3.4 |
| Lapis C + 7 varian + dua kekhususan (`IsProRate`, `FlagOldData`) | §4 |
| Kontrak SELISIH + peran dua rasio prorata + 5 varian menyimpang | §5 |
| Batas & 3 gerbang + konsekuensi desain (Life, >100 lokasi) | §6 |
| Out of scope + pertanyaan terbuka | §7, §8 |

⛔ **`/to-spec` TIDAK dijalankan.** Bahan berhenti di sini; work owner yang menjalankannya manual.

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
