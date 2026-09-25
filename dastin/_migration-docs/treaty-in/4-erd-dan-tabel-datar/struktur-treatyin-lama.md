# Struktur Data Treaty In Sistem Lama — pohon `TreatyIn` sampai turunan terdalam

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Treaty In` — **329 berkas XML**, disapu seluruhnya (bukan cuplikan).
> Pembanding diambil dari `D:\XML_NURE\Treaty In Adjustment` — **379 berkas** — **hanya** untuk menandai simpul mana yang juga muncul di sana. **Perilaku modul Adjustment tidak diadili di berkas ini** (embargo sebagian, lihat `KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md`).
> Dibangkitkan [`alat/buat-pohon-treatyin.py`](../alat/buat-pohon-treatyin.py) — sapuan **23 September 2026**.
>
> **Berkas ini memerikan SISTEM LAMA, bukan rancangan sistem baru.** Nama properti ditulis apa adanya seperti di XML, termasuk salah ejanya (`ProfitCommision`, `Bordeaux`, `ClassofBusiness`, `TeritorialScope`, `Reinstatement_List`). Penerjemahan ke istilah `SPEC-MODEL-DATA.md` dan §16 terjadi di berkas lain, bukan di sini.
>
> **TURUNAN.** Ketiga CSV di §1.4 dan `pohon-treatyin.txt` adalah pendamping mesin-baca berkas ini. Bila CSV berbeda dari sumber XML-nya, **alatnya yang salah**.

Melengkapi [`STRUKTUR-DATA.md`](STRUKTUR-DATA.md) yang memerikan **entitas rancangan baru**. Berkas ini memerikan **bentuk lama apa adanya**, sampai kedalaman 4, dengan kamus field per simpul.

---

## 1. Cara baca

### 1.1 Notasi

| Tulisan | Arti |
|---|---|
| `.Prop` | properti tunggal (skalar) |
| `Prop[]` | **Page List** — baris berulang, selalu dirujuk dengan indeks `(...)` di XML |
| `Prop{}` | **Page** — satu halaman bersarang, tanpa indeks |
| `→ KELAS` | kelas Pega yang menaungi isi simpul (awalan `ASM-FW-` dipangkas agar muat) |
| `n=` | cacah rujukan di ekspor — penanda seberapa hidup simpul itu, **bukan** penanda penting |

### 1.2 Kolom `BUKTI` di CSV

| Nilai | Artinya | Cacah |
|---|---|---:|
| `jalur-penuh` | jalur lengkapnya terbaca utuh di XML, mis. `TreatyIn.Limits(1).Detail(1).QSPct` | **927** |
| `hanya-adjustment` | tidak pernah ditulis di ekspor Treaty In; terbaca di ekspor Adjustment — lihat §6 | **58** |
| `alias-primary` | hanya terbaca lewat alias step-page `Primary` | **0** |

**Seluruh pohon ini bersandar pada jalur penuh, tidak satu simpul pun pada alias.** Itu bukan
kebetulan dan perlu dinyatakan: alat penyusunnya **memang** menerima `Primary.*` sebagai `TreatyIn.*`,
tetapi hanya di rule yang applies-to `ASM-FW-GISFW-Int-TREATY_IN` — dan di ekspor ini **tidak ada**
rule semacam itu yang memakai `Primary`.

Tiga aktivitas sempat terlihat memenuhi syarat itu (`FetchQSfromMasterXOL`, `SetSpreadingXOL`,
`SetSpreadName`) karena nama kelas `ASM-FW-GISFW-Int-TREATY_IN` muncul di dalam berkasnya. Yang muncul
itu **entri indeks yang merujuk aturan lain**, bukan kelas aturannya sendiri; ketiganya sebenarnya
applies-to `Data-TreatyInShare` dan `Data-TreatyInLimitsDetail`. Bila syarat applies-to itu tidak
dipasang, `Primary.SpreadingListXOL` akan masuk sebagai simpul palsu `TreatyIn.SpreadingListXOL`
di tingkat akar — padahal `SpreadingListXOL` hanya hidup di bawah `Share[]`.

> Ini kasus khusus aturan yang sudah berlaku di proyek ini: **cabang sebuah properti ditentukan
> pemiliknya, bukan berkas tempat ia ditemukan** (CONTEXT.md §2.8).

Kolom `SUMBER_KELAS` memisahkan kelas yang **dinyatakan sistem** (`pyPagesAndClasses`, `pxPageClass`, `pyStepsClassName`) dari kelas yang **disimpulkan dari bentuk anaknya** (`identitas-kelas (n/m anak cocok)`). Bedanya penting — jangan diratakan.

### 1.3 Tipe data

Pega **tidak mengekspor `Rule-Obj-Property`** di paket ini. Yang ada hanyalah entri indeks rujukan (`pxRuleObjClass = Rule-Obj-Property`) yang menyebut **nama** properti per kelas — 4 675 entri, 61 kelas — tanpa satu pun tipe, panjang, atau keterisian.

Artinya: **tidak ada satu pun tipe yang dinyatakan langsung oleh sistem lama.** Presisi angka diputuskan ADR-0003, bukan dibaca dari sini. Yang bisa dipegang dari berkas ini adalah **bentuk pohon dan nama**, bukan tipenya.

### 1.4 Pendamping mesin-baca

| Berkas | Baris | Isi | Kolom |
|---|---:|---|---|
| [`datar-treatyin-lama.csv`](datar-treatyin-lama.csv) | 985 | **pohon properti** — satu baris per simpul | `JALUR` · `KEDALAMAN` · `NAMA` · `JENIS` · `KELAS_PEGA` · `CACAH_ANAK` · `REF` · `CACAH_BERKAS` · `BUKTI` · `SUMBER_KELAS` · `ADA_DI_ADJUSTMENT` · `BERKAS_CONTOH` |
| [`datar-treatyin-kelas.csv`](datar-treatyin-kelas.csv) | 17 | **peta kelas** — kelas Pega dan simpul yang memakainya | `KELAS_PEGA` · `CACAH_PEMAKAIAN` · `DIPAKAI_SEBAGAI` |
| [`datar-treatyin-muatan-keluar.csv`](datar-treatyin-muatan-keluar.csv) | 175 | **muatan keluar** — tiap argumen ke empat prosedur Oracle (§7) | `PROSEDUR` · `URUT_ARGUMEN` · `DIISI_DARI` · `BERKAS` |
| [`pohon-treatyin.txt`](pohon-treatyin.txt) | 187 | pohon §4 siap tempel | — |

`JALUR` memakai **nama kanonik tanpa indeks**. `TreatyIn.Limits.Detail` berarti `TreatyIn.Limits(n).Detail(m)` di XML; `JENIS = Page List` pada simpul induknya yang menandai perulangan, bukan tanda kurung di jalurnya.

---

## 2. Ringkasan angka

| | |
|---|---:|
| Simpul dalam pohon `TreatyIn` | **985** |
| — di antaranya **inti** (di luar empat salinan §5) | **459** |
| — di antaranya di dalam **empat salinan rekursif** | **526** |
| Kedalaman maksimum | **4** tingkat di bawah `TreatyIn` |
| Properti skalar langsung di `TreatyIn` | **99** |
| Kontainer (Page/Page List) langsung di `TreatyIn` | **43** |
| Kelas Pega berbeda yang menyusun pohon | **17** (+5 tak dideklarasikan) |
| Entri indeks properti tersedia | 4 675 pada 61 kelas — **nama saja, tanpa tipe** |

> **Lebih dari separuh pohon ini adalah salinan dirinya sendiri.** 526 dari 985 simpul berada di dalam `ActualValue`, `ValueDifference`, `OLDDATA`, atau `ValueBeforeProrate` — empat halaman yang berbentuk `TreatyIn` di dalam `TreatyIn`. Angka 985 karena itu **tidak boleh** dibaca sebagai 985 hal yang berbeda.

---

## 3. Peta kelas

Tujuh belas kelas menyusun seluruh pohon. **Satu kelas dipakai ulang di banyak tempat** — inilah sebab pohonnya terlihat berulang-ulang.

| Kelas Pega | Pakai | Dipakai sebagai | Catatan |
|---|---:|---|---|
| `ASM-FW-GISFW-Data-TreatyInLimitsSpreading` | **120** | seluruh daftar `…List[]` bermuatan `Currency` + `Value` | **satu bentuk baris, 120 peran** — lihat peringatan di bawah |
| `ASM-FW-GISFW-Data-TreatyInLimits` | 9 | `Limits`, `Installment`, `Reinstatement_List`, `COBList`, `ClassOfBusinessList`, `FacultativeLimits` | |
| `ASM-FW-GISFW-Data-TreatyInShare` | 7 | `Share`, `FacultativeShareList` (+ salinannya) | dua peran berbeda, satu kelas |
| `ASM-FW-GISFW-Data-TreatyInTotal` | 7 | `IOOLimitList`, `EgnpiTotalList`, `TotalEgnpiAmountNP`, `TotalRetentionAmountNP` | |
| `ASM-FW-GISFW-Data-LimitSummaryList` | 6 | `LimitSummaryList`, `LimitShareSummaryList` | |
| `ASM-FW-GISFW-Int-TREATY_IN` | 5 | `TreatyIn`, `ActualValue`, `OLDDATA`, `ValueBeforeProrate`, `ValueDifference` | **akar, bersarang di dalam dirinya sendiri** — §5 |
| `ASM-FW-GISFW-Data-TreatyInDeduction` | 5 | `DeductionList` | |
| `ASM-FW-GISFW-Data-TreatyInEGNPI` | 4 | `EGNPI`, **`Retention`** | `Retention` memakai kelas EGNPI — nama kelas bukan penanda arti |
| `ASM-FW-GISFW-Data-TreatyInInstallment` | 3 | `InstallmentList`, `Installment` (di dalam salinan) | |
| `ASM-FW-GISFW-Data-TreatyInShareReins` | 2 | `ShareFacultativeReinsurers` | |
| `ASM-FW-GISFW-Data-TreatyInLimitsDetail` | 2 | `Limits.Detail`, `FacultativeLimits.Detail` | |
| `ASM-FW-GISFW-Data-TreatyInAccumulation` | 1 | `AccumulationList` | |
| `ASM-FW-GISFW-Data-SuggestList` | 1 | `CommentList` | kelas pinjaman dari domain saran |
| `ASM-FW-GISFW-Data-TreatyInCurrencyList` | 1 | `CurrencyList` | |
| `ASM-FW-GISFW-Data-TreatyInLimitLayer` | 1 | `Limits.TreatyGroupList` | |
| `ASM-FW-GISFW-Data-TreatyInPortfolio` | 1 | `Portfolio` | |
| `ASM-FW-GISFW-Data-TreatyInAccountReport` | 1 | `ReportingPeriodList` | |

**Lima simpul tidak punya deklarasi kelas di mana pun dalam 329 berkas ini**, dan namanya juga tidak pernah muncul berkelas di tempat lain:

`AchievementLists` · `BrokerageList` · `LimitFacShareSummaryList` · `MDPSummaryList` · `ShareReins`

Simpulnya nyata dan sebagian isinya terbaca; hanya nama kelasnya yang tidak ada. **Jangan tebak nama kelasnya — pakai isinya.**

Selain itu ada **24 nama** yang berkelas jelas di satu tempat tetapi di tempat lain hanya pernah dirujuk satu anaknya, sehingga di lokasi itu bentuknya tidak cukup untuk memastikan kelasnya. Di CSV kolom `KELAS_PEGA` sengaja dibiarkan `(tidak dideklarasikan)` di lokasi tersebut — **bukti tipis tidak dinaikkan pangkatnya menjadi bukti**.

> ### Konsekuensi untuk pemetaan tabel
>
> **Satu kelas Pega ≠ satu tabel.** `ASM-FW-GISFW-Data-TreatyInLimitsSpreading` muncul di **120** tempat, dan hampir seluruhnya berisi pasangan `Currency` + `Value` yang sama persis. Ia bukan sebuah entitas; ia **bentuk pembawa satu nilai uang bermata uang** — persis paket uang ADR-0007. Memetakannya jadi satu tabel akan menumpuk seratus dua puluh besaran berbeda ke dalam satu ember.
>
> Kebalikannya juga berlaku: `Retention` berkelas `TreatyInEGNPI`, dan `CommentList` berkelas `SuggestList`. Nama kelas di ekspor ini **tidak dapat dipakai** sebagai penanda arti bisnis.

---

## 4. Pohon `TreatyIn`

Pohon utuh ada di [`pohon-treatyin.txt`](pohon-treatyin.txt). Di bawah ini bentuk ringkasnya; empat salinan rekursif ditampilkan sebagai satu baris agar isinya tidak tercetak lima kali.

```
TreatyIn{}                            → GISFW-Int-TREATY_IN                 n=142
│
├─ 99 properti skalar (identitas, periode, parameter hitung)  .............  §4.1
│
├─ ValueDifference{}                  → GISFW-Int-TREATY_IN     217 simpul  ◄── SALINAN UTUH, §5
├─ ActualValue{}                      → GISFW-Int-TREATY_IN     152 simpul  ◄── SALINAN UTUH, §5
├─ OLDDATA{}                          → GISFW-Int-TREATY_IN     142 simpul  ◄── SALINAN UTUH, §5
├─ ValueBeforeProrate{}               → GISFW-Int-TREATY_IN      15 simpul  ◄── SALINAN UTUH, §5
│
├─ Limits[]                           → GISFW-Data-TreatyInLimits           n=307  ★ inti layer
│    ├─ 18 field (§4.2)
│    ├─ Detail[]                      → GISFW-Data-TreatyInLimitsDetail     n=104  ★ inti proporsional
│    │    ├─ 14 field (§4.3)
│    │    ├─ EPIList[] · RetentionList[] · CashLossList[]
│    │    │  · ClaimCoopList[] · PLAList[]  → GISFW-Data-TreatyInLimitsSpreading
│    │    │       └─ Currency · Value
│    │    ├─ IOOLimitList[]           → GISFW-Data-TreatyInTotal
│    │    │       └─ Currency · CurrencyID · Value · ID
│    │    ├─ COBList[]                → GISFW-Data-TreatyInLimits
│    │    │       └─ ClassOfBusiness · ClassOfBusinessID
│    │    └─ AchievementLists[]       → (tidak dideklarasikan)              n=1
│    ├─ MDPList[] · EgnpiTotalList[] · PremiumEarnedList[]
│    ├─ Reinstatement_List[]          → GISFW-Data-TreatyInLimits           n=8
│    │       └─ ReinstatementPct · AdditionalAmount1 · AdditionalAmount2
│    └─ TreatyGroupList[]             → GISFW-Data-TreatyInLimitLayer       n=6
│         └─ ClassOfBusinessList[]                                  ◄── KEDALAMAN 4
│
├─ Share[]                            → GISFW-Data-TreatyInShare            n=237  ★ inti penyebaran
│    ├─ 10 field (§4.4)
│    ├─ SpreadingListXOL[]            → GISFW-Data-TreatyInLimitsSpreading  n=45
│    │    ├─ Pct · ReinsTypeName · NetPremiumList · RnmLimitList
│    │    │  · DeductionTotalList · GrossPremiumMinList
│    │    └─ 11 daftar RNMSpreadedList…[]  (Currency · Value)      ◄── KEDALAMAN 4
│    ├─ DeductionList[]               → GISFW-Data-TreatyInDeduction
│    │       └─ Deduction · DeductionPct
│    └─ 17 daftar uang lain (GrossPremiumList, RnmLimitList, RNMSpreadedList…)
│
├─ FacultativeShareList[]             → GISFW-Data-TreatyInShare            n=58
│    ├─ 9 field · DeductionList[] · GrossPremiumList[] · RnmLimitList[]
│    └─ ShareFacultativeReinsurers[]  → GISFW-Data-TreatyInShareReins
│         └─ ReinsName · ReinsID · SharePct · Amount · Amount2
│
├─ ShareFacultativeReinsurers[]       → GISFW-Data-TreatyInShareReins       n=47
│    ├─ ID · Layer · ReinsID · ReinsName · BrokerName · SharePct
│    └─ FacultativeLimits[]           → GISFW-Data-TreatyInLimits
│         └─ Detail[]                 → GISFW-Data-TreatyInLimitsDetail  ◄── KEDALAMAN 4
│
├─ CommentList[]                      → GISFW-Data-SuggestList              n=37
│    └─ OperatorName · Date · Suggest · IsApproved
├─ LimitShareSummaryList[]            → GISFW-Data-LimitSummaryList         n=35
├─ LimitSummaryList[]                 → GISFW-Data-LimitSummaryList         n=18
│    └─ Limit · Limit2 · Deductible · Deductible2 · MDP · MDP2
│       · AggregateLimit · AggregateLimit2
├─ EGNPI[]                            → GISFW-Data-TreatyInEGNPI            n=28
│    └─ Amount · AmountIDR · Currency · CurrencyID · Proportion
│       · AsDate · ID · TreatyGroup · TreatyGroupID
├─ Retention[]                        → GISFW-Data-TreatyInEGNPI            n=19
├─ CurrencyList[]                     → GISFW-Data-TreatyInCurrencyList     n=18
│    └─ Currency · CurrencyID · Conversion · PeriodStart · PeriodEnd
├─ Installment[]                      → GISFW-Data-TreatyInLimits           n=21
│    └─ InstallmentList[]  ← bentuknya HANYA terbaca lewat salinan, §5.2
├─ ReportingPeriodList[]              → GISFW-Data-TreatyInAccountReport    n=15
│    └─ Period · InitialDate · SubmissionDue · ConfirmationDue · SettlementDue
├─ AccumulationList[]                 → GISFW-Data-TreatyInAccumulation     n=14
│    └─ Period · ReportDate
├─ Portfolio[]                        → GISFW-Data-TreatyInPortfolio        n=8
│    └─ Type · TypePortfolio · Description
├─ ShareReins[]                       → (tidak dideklarasikan)              n=8
├─ MDPSummaryList[]                   → (tidak dideklarasikan)              n=6
├─ LimitFacShareSummaryList[]         → (tidak dideklarasikan)              n=5
│
└─ 20 daftar rekap Total…[]           → GISFW-Data-TreatyInLimitsSpreading  §4.5
```

### 4.1 Skalar langsung di `TreatyIn` (99)

Angka dalam kurung adalah cacah rujukan.

**Kendali layar & keadaan proses** (10) — **bukan** data bisnis
`ViewState`(623) · `EDMMaterialType`(472) · `IsEditData`(110) · `EDMState`(67) · `StatusAkseptasi`(61) · `Position`(58) · `ChooseStatusAkseptasi`(21) · `RevisionState`(13) · `PositionUsername`(32) · `pyErrMsg`(6)

> Dua yang paling sering dirujuk di seluruh pohon — `ViewState` dan `EDMMaterialType` — adalah **kendali tampilan**, bukan besaran bisnis. Cacah rujukan mengukur keramaian kode, bukan bobot bisnis.

**Identitas kontrak** (8)
`ID`(62) · `OLDID`(6) · `TreatyContractName`(12) · `ContractRefNo`(3) · `TreatyYear`(21) · `ClassofBusiness`(2) · `TeritorialScope`(6) · `Information`(33)

**Periode & masa berlaku** (5)
`Commencement`(60) · `Termination`(12) · `EDMEffective`(7) · `RevisionDate`(0) ⟦hanya Adjustment⟧ · `AccumulationPeriod`(13)

**Sifat & bentuk kontrak** (6)
`ProportionType`(77) · `AccountingMode`(6) · `AccountingModeNonProp`(3) · `OptionLimit`(6) · `CoInScale`(4) · `IsMultipleRetro`(35)

**Pihak** (10)
`Ceding`(24) · `CedingID`(15) · `CedingStatusActive`(5) · `LeadingReinsSource`(20) · `LeadingReinsSourceID`(14) · `LeadingReinsName`(3) · `LeadingReinsID`(2) · `SourceStatusActive`(5) · `TreatyLeader`(3) · `Comment`(37)

**Bagian & brokerage** (10)
`RNMShare`(19) · `RNMShareP`(17) · `RNMShareAcrossTheBoard`(13) · `RnmShareDeducted`(9) · `NusareSharePct`(2) · `BrokeragePercent`(19) · `BrokeragePercentP`(8) · `BrokeragePct`(2) · `FacultativeShare`(76) · `FacultativeShareBrokerage`(15)

**Fakultatif ringkas** (2)
`FacShare`(6) · `FacShareBrokerage`(6)

> `FacShare` dan `FacShareBrokerage` **diisi dari** `FacultativeShare` dan `FacultativeShareBrokerage` (`TreatyInXOLAddSpreading` langkah 6, **hidup**). Dua pasangan nama untuk satu besaran.

**Pro rate** (4)
`IsProRate`(19) · `ProRatePercent`(49) · `ProRateDays`(8) · `ProRateTotalDays`(7)

**Pelaporan & penagihan** (10)
`ReportingPeriod`(21) · `ReportingStart`(17) · `ReportingEnd`(11) · `ReportingInterval`(6) · `ReportingSubmission`(8) · `ReportingConfirmation`(8) · `ReportingSettlement`(8) · `ReminderDays`(4) · `Bordeaux`(6) · `BordereauxNote`(4)

**Angsuran** (4)
`InstallmentNo`(50) · `TotalInstallmentAmount`(7) · `TotalInstallmentPct`(7) · `CurrencyInstallmentAmount`(4)

**Batas bahaya** (10) — sepasang nilai + mata uang per bahaya
`Earthquake`(4) · `CurrencyEarthquake`(4) · `FloodJab`(4) · `CurrencyFloodJab`(4) · `FloodNation`(4) · `CurrencyFloodNat`(4) · `RSMDLimit`(4) · `CurrencyRSMD`(4) · `MaxCoGroup`(4) · `MaxCoNonGroup`(4)

**Rekap terhitung** (13) — turunan, bukan masukan
`TotalEgnpiAmount`(23) · `TotalEgnpiProportion`(10) · `CurrencyEgnpiAmount`(1) · `TotalLimitsROL`(9) · `TotalLimitsIOOLimit`(4) · `TotalLimitsPremiumEarned`(4) · `TotalLimitsAdjPct`(3) · `TotalLimitsDeductible`(3) · `TotalLimitsMdp`(3) · `TotalRetentionAmount`(3) · `TotalShareGross`(3) · `TotalShareNet`(3) · `TotalShareRnmLimit`(5)

**Naskah kontrak** (4) — berpasangan, bervarian `…P`
`Exclusions`(9) · `ExclusionsP`(7) · `SpecialConditions`(9) · `SpecialConditionsP`(7)

**Sisa** (3)
`Currency`(1) · `RetroList`(4) · `AddendumPremi`(0) ⟦hanya Adjustment⟧

> **Pola sufiks `P`.** `RNMShareP`, `BrokeragePercentP`, `ExclusionsP`, `SpecialConditionsP` — empat nama bersufiks `P` yang berpasangan dengan nama tanpa sufiks. Artinya **tidak dinyatakan di mana pun di ekspor**; ia terbaca dari pemakaian, bukan dari deklarasi. Lihat `PENGETAHUAN.md` untuk pembahasan `RNMShareP`.

### 4.2 `Limits[]` — baris layer (18 field)

Simpul paling hidup kedua di pohon (n=307). **Satu baris = satu layer pada satu kontrak.**

| Kelompok | Field |
|---|---|
| Identitas layer | `ID` · `Layer` · `LayerPart` · `LayerPartType` · `LayerType` · `Cover` · `TreatyType` |
| Batas berpasangan | `Limit` / `Limit2` · `Deductible` / `Deductible2` · `Currency` / `Currency2` |
| Parameter XOL | `MDPPct` · `ROLPct` · `AdjRate` |
| Reinstatement | `ReinstatementPct` · `ReinstatementValue` |

> **Sufiks `2` adalah mata uang kedua, bukan nilai kedua.** `Limit`/`Limit2` berpasangan dengan `Currency`/`Currency2`. Satu baris layer membawa **dua** batas dalam **dua** mata uang, bukan satu batas dengan cadangan.

### 4.3 `Limits[].Detail[]` — rincian proporsional (14 field)

| Kelompok | Field |
|---|---|
| Penggolongan | `TreatyType` · `TreatyGroup` · `TreatyGroupID` · `SpreadingTypeID` |
| Porsi | `QSPct` · `Surplus` · `CessionList` · `RNMShareList` · `SpreadingList` |
| Komisi | `ProfitCommision` ⟦salah eja di sumbernya⟧ · `ProfitME` · `ProfitYDCF` |
| Belum terurai | `RIOGR` · `RIONR` |

`RIOGR` dan `RIONR` **tidak terbaca kepanjangannya** di 329 berkas ini — tercatat sebagai pertanyaan terbuka, tidak ditebak.

### 4.4 `Share[]` — baris penyebaran (10 field + 20 daftar uang)

**Satu baris = satu penyebaran per layer.**

| Kelompok | Field |
|---|---|
| Identitas layer (sama seperti `Limits[]`) | `Cover` · `Layer` · `LayerPart` · `LayerPartType` · `LayerType` |
| Penyebaran | `SpreadingTypeXOL` · `SpreadingTypeIDXOL` · `SpreadingTotalPctXOL` |
| Bagian | `RNMShare` · `ClassofBusinessList` |

Dua puluh daftar uang menggantung di bawahnya, semuanya berbentuk `Currency` + `Value`, dengan nama yang membedakan **peran**, bukan bentuk: `GrossPremiumList`, `NetPremiumList`, `RnmLimitList`, `RnmGrossPremiDisplay`, `RnmLimitListDisplay`, `BrokerageList`, `GrossPremiumMinList`, dan sebelas varian `RNMSpreadedList…` (`XOL`, `RIXOL`, `GrossXOL`, `GrossRIXOL`, `GrossMinXOL`, `GrossRIMinXOL`, `DeductXOL`, `DeductRIXOL`, `NetXOL`, `NetRIXOL`).

> Sebelas varian `RNMSpreadedList…` adalah **matriks tiga sumbu** yang diratakan ke nama: {kotor, potongan, bersih} × {dengan RI, tanpa RI} × {biasa, minimum}. Bentuk barisnya identik; yang membedakan hanya namanya.

### 4.5 Dua puluh daftar rekap `Total…[]` di tingkat akar

Seluruhnya berkelas `TreatyInLimitsSpreading` dan berisi `Currency` + `Value`:

`TotalSpreadedRnmProp` · `TotalSpreadedRnmRIProp` · `TotalSpreadedNetPremi` · `TotalSpreadedNetPremiRI` · `TotalShareNetNP` · `TotalShareGrossNP` · `TotalShareGrossMinNP` · `TotalShareRnmNP` · `TotalShareRnmProp` · `TotalShareDeductionNP` · `TotalLimitIOONP` · `TotalLimitDeductblNP` · `TotalLimitMDPNP` · `TotalLimitMDPMinNP` · `TotalLimitPremiEarnNP` · `TotalInstallmentNP` · `TotalFacShareNetNP` · `TotalFacShareGrossNP` · `TotalFacShareRnmNP` · `TotalFacShareDeductionNP`

Ditambah tiga yang berkelas `TreatyInTotal` (punya `ID` selain `Currency`+`Value`): `TotalEgnpiAmountNP` · `TotalRetentionAmountNP` · (dan `IOOLimitList` di tingkat `Detail`).

Sufiks `NP` menandai cabang **non-proporsional**, `Prop` menandai **proporsional**. Keduanya hidup berdampingan di satu halaman yang sama.

---

## 5. Cabang rekursif: empat salinan `TreatyIn` di dalam `TreatyIn`

Ini bagian yang paling mudah salah kutip, dan paling penting untuk sambungan ke modul Adjustment.

Kelas `ASM-FW-GISFW-Int-TREATY_IN` **bersarang di dalam dirinya sendiri** pada empat tempat:

| Halaman | Simpul | Kelas dinyatakan? | Peran yang terbaca dari kode |
|---|---:|---|---|
| `TreatyIn.ValueDifference` | 217 | ya — `pyPagesAndClasses`, **12 kali** | menerima hasil **pengurangan** |
| `TreatyIn.ActualValue` | 152 | ya — `pyPagesAndClasses`, **27 kali** | lihat peringatan §5.3 |
| `TreatyIn.OLDDATA` | 142 | **tidak, di mana pun** — bentuk anaknya yang sama dengan akar | memegang **nilai lama** |
| `TreatyIn.ValueBeforeProrate` | 15 | ya — `pyPagesAndClasses`, 1 kali | nilai sebelum pro rate |

### 5.1 Mesin selisih, dibaca apa adanya

Lima aktivitas menyentuh `OLDDATA`: `TreatyEDMDifferenceLimits`, `TreatyEDMDifferenceShare`, `TreatyEDMDifferencePremium`, `TreatyEDMDifferenceDeduction`, `TreatyInSetValueInstallment`. Keempat yang pertama dipanggil **`TreatyEDMCalculateDifference`**, yang punya **11 langkah dan seluruhnya hidup** (tidak satu pun ber-`pyStepsBlockName = "//"`).

Bentuk isiannya seragam, dikutip apa adanya:

```
TreatyIn.ValueDifference.Limits(<CURRENT>).Limit
      <=  .Limit - TreatyIn.OLDDATA.Limits(<CURRENT>).Limit

TreatyIn.ValueDifference.LimitSummaryList(<CURRENT>).AggregateLimit2
      <=  .AggregateLimit2 - TreatyIn.OLDDATA.LimitSummaryList(<CURRENT>).AggregateLimit2

TreatyIn.ValueDifference.Share(local.subscript).DeductionList(<CURRENT>).Deduction
      <=  .Deduction - TreatyIn.OLDDATA.Share(local.subscript).DeductionList(<CURRENT>).Deduction

TreatyIn.ValueDifference.TotalLimitsROL
      <=  .TotalLimitsROL - TreatyIn.OLDDATA.TotalLimitsROL
```

Tiga hal yang terbaca langsung, tanpa penafsiran:

1. **Selisih = nilai sekarang − nilai lama**, dihitung per besaran, bukan per baris.
2. **Jalur sumber dan jalur tujuan identik bentuknya.** `X` di `TreatyIn`, `X` di `OLDDATA`, `X` di `ValueDifference` — tiga halaman berbentuk sama, dibaca dengan indeks yang sama (`<CURRENT>`, `local.subscript`).
3. **Selisih hanya dihitung untuk besaran angka.** Dua isian di `TreatyEDMDifferenceShare` menyalin tanpa pengurangan (`SpreadingTypeXOL`, `SpreadingTotalPctXOL` ← `OLDDATA` langsung), jadi tidak semua field di `ValueDifference` berisi selisih.

Butir 3 itu penting dan mudah terlewat: **`ValueDifference` bukan seluruhnya selisih.**

### 5.2 Bentuk sebuah simpul kadang hanya terlihat lewat salinannya

`TreatyIn.Installment` di tingkat akar hanya pernah dirujuk sampai `Currency` dan `ID`. Anaknya `InstallmentList[]` — dengan `DueDate`, `InstallmentPct`, `PaymentDate`, `WPC` — **hanya terbaca di dalam `OLDDATA` dan `ValueDifference`**.

Karena keempat halaman itu berkelas sama, bentuk yang terbaca di salinan **berlaku juga di akarnya**. Di CSV lokasi seperti ini bertanda `SUMBER_KELAS = identitas-kelas`, bukan `dideklarasikan`.

### 5.3 Satu nama yang tidak cocok dengan isinya

Di `TreatyEDMDifferenceShare`, delapan isian menulis ke `TreatyIn.ActualValue.LimitShareSummaryList(<CURRENT>).*` dengan rumus **pengurangan**:

```
TreatyIn.ActualValue.LimitShareSummaryList(<CURRENT>).Limit
      <=  .Limit - TreatyIn.OLDDATA.LimitShareSummaryList(<CURRENT>).Limit
```

Nama halamannya berbunyi "nilai aktual"; yang masuk adalah **selisih**. Ini dicatat sebagai fakta bentuk, **bukan** sebagai penilaian atas perilaku modul Adjustment — adjudikasinya ditangguhkan ke [`TEMUAN-ADJUSTMENT-DITUNDA.md`](TEMUAN-ADJUSTMENT-DITUNDA.md).

> **Konsekuensi untuk rancangan.** Empat halaman berbentuk sama **tidak boleh** menjadi empat tabel bersalin-rupa. Ia satu bentuk dengan empat peran: nilai berjalan, nilai lama, selisih, dan nilai sebelum pro rate. Aturan bisnis yang mengikat sudah ditetapkan: **data lama TIDAK disimpan ulang — ia DISELECT**, dan selisih ditampung tabel tersendiri. Lihat [`KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md`](KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md).

---

## 6. Apa yang ditambahkan modul Adjustment

Dibaca dari `D:\XML_NURE\Treaty In Adjustment` **hanya** untuk menjawab: simpul apa yang ada di sana dan tidak di sini.

**Jawabannya: 58 simpul, dan 56 di antaranya adalah `TreatyIn.OLDDATA.*`.**

| Yang baru | Cacah | Keterangan |
|---|---:|---|
| `TreatyIn.OLDDATA.<field>` | 56 | field yang **sudah ada** di akar `TreatyIn`; yang baru hanyalah bahwa salinan lamanya ikut dirujuk |
| `TreatyIn.AddendumPremi` | 1 | skalar tingkat akar, tidak pernah dirujuk di ekspor Treaty In |
| `TreatyIn.RevisionDate` | 1 | skalar tingkat akar, tidak pernah dirujuk di ekspor Treaty In |

> **Modul Adjustment tidak memperkenalkan satu pun kelas baru, satu pun daftar baru, dan satu pun tingkat kedalaman baru.** Ia memakai pohon yang sama, dan yang berbeda hanyalah **seberapa luas salinan `OLDDATA` dipakai** — ditambah dua skalar.
>
> Ini menguatkan keputusan G1 di `KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md` — **tidak ada entitas baru** — dan kali ini dari arah bukti yang berbeda: bukan dari daftar properti, melainkan dari bentuk pohonnya.

Daftar lengkap 58 simpul ada di `datar-treatyin-lama.csv`, kolom `BUKTI = hanya-adjustment`.

---

## 7. Turunan di luar clipboard

Pohon di atas adalah bentuk **di memori**. Saat keluar dari Pega, ia menyempit tajam.

### 7.1 Empat prosedur Oracle

| Prosedur | Dipanggil dari | Argumen |
|---|---|---:|
| `POOLDATA.PEGA_TREATY_IN` | `RDBList/SaveTreatyIn.xml` | dari `TreatyIn.*` langsung |
| `POOLDATA.PEGA_M_TREATY_IN_EDM` | `RDBList/SaveTreatyInEDM.xml` | dari `TreatyIn.*`, **membawa `TreatyIn.OLDID`** |
| `POOLDATA.PEGA_M_TREATY_IN_DETAIL` | `RDBList/SaveTreatyInDetail.xml` | dari `SaveData.*` (halaman datar antara) |
| `POOLDATA.PEGA_M_TREATY_IN_DETAIL_EDM` | `RDBList/SaveTreatyInDetailEdm.xml` | dari `SaveData.*` |

Seluruh 175 argumen ada di [`datar-treatyin-muatan-keluar.csv`](datar-treatyin-muatan-keluar.csv).

Dua hal yang terbaca dari daftar itu:

- **Dua jalur berbeda bentuk masukannya.** Prosedur kontrak dibaca **langsung** dari `TreatyIn.*`; prosedur detail dibaca dari halaman perantara `SaveData.*` yang sudah diratakan lebih dulu. Keduanya tidak setara sebagai sumber.
- **`OLDID` hanya muncul di jalur EDM.** `PEGA_M_TREATY_IN_EDM` menerima `{TreatyIn.ID}` **dan** `{TreatyIn.OLDID}`; jalur kontrak biasa hanya `{TreatyIn.ID}`. Inilah penanda sambungan addendum di tingkat basis data.

### 7.2 Tabel yang disentuh langsung lewat RDB List

`M_TREATY_IN` · `TREATY_IN_EDM` · `TREATYINDETAIL` · `TREATYINDETAILEDM` · `M_TREATY_IN_EDM` · `M_TREATY_IN_DETAIL` · `M_TREATY_IN_DETAIL_EDM` · `PROPORTIONALARRG` · `TREATYEXCHANGEYEARLY` · `M_ATTACHMENTTREATY_2` · `T_STORAGE_IMAGE` · `CATEGORY_ATTACH_REAS`

Enam di antaranya hanya muncul dalam pernyataan `DELETE` — jalur pembersihan sebelum tulis ulang, bukan jalur baca.

### 7.3 Yang TIDAK ditemukan

`TREATYINOFFER` **tidak punya satu pun penulis yang terjangkau** di ekspor ini (dibuktikan sapuan dua tingkat pada sesi sebelumnya; dicatat di `DAFTAR-ESKALASI-MANAJEMEN.md`). Ia karena itu **bukan** jalur penerbitan yang bisa dipakai sebagai sumber migrasi maupun dasar rekonsiliasi.

---

## 8. Batas bukti

Hal-hal yang **tidak** bisa dijawab berkas ini, dan tidak boleh dikarang dari isinya:

1. **Tipe, panjang, dan presisi properti.** Tidak ada `Rule-Obj-Property` di ekspor — hanya entri indeks bernama. Presisi ditetapkan ADR-0003, bukan dibaca dari sini.
2. **Wajib atau opsional.** Tidak ada satu pun aturan validasi tingkat properti; validasi yang ada hidup di dalam Activity sebagai isian bersyarat.
3. **Nilai enum.** `ProportionType`, `StatusAkseptasi`, `EDMState`, `EDMMaterialType`, `ViewState`, `AccountingMode`, `LayerType`, `LayerPartType` dipakai sebagai kode tanpa daftar nilainya di mana pun. Sebagiannya terbaca dari perbandingan di rule When; sisanya menunggu Uji S-7.
4. **Kelas untuk lima simpul** di §3 (`AchievementLists`, `BrokerageList`, `LimitFacShareSummaryList`, `MDPSummaryList`, `ShareReins`). Isinya sebagian terbaca, namanya tidak ada.
5. **Kepanjangan `RIOGR`, `RIONR`, `RSMD`** dan arti `ProfitME` / `ProfitYDCF`. Tidak terbaca; menunggu jawaban dari luar.
6. **Kardinalitas sebenarnya.** Bentuk `Page List` tidak menjamin lebih dari satu baris pernah ada. Yang sebaliknya juga mungkin dan tidak terlihat dari struktur.
7. **Arti sufiks `P`** pada `RNMShareP`, `BrokeragePercentP`, `ExclusionsP`, `SpecialConditionsP`. Terbaca dari pemakaian, tidak dari deklarasi.

Tujuh butir ini sejalan dengan aturan kerja yang berlaku: **ambil perilaku yang terlihat di XML sebagai sumber kebenaran, tandai asumsinya, jangan berhenti menunggu konfirmasi.**
