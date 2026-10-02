# Bahan spec — Modul 5: 354 Berkas EDM-Only

> **Ini BAHAN untuk `/to-spec`, bukan spec.** `/to-spec` ber-`disable-model-invocation: true` dan
> **belum dijalankan**.
>
> **Sumber:** `07-edm\05-e4-edm-only.md`. Prior art bentuk: `04-spec\05/06/07/08-spec-edm`.
>
> **Keputusan mengikat:** K-005 · K-006 · K-010/K-012 · K-018 · K-043 · K-044 · K-046 · K-047 ·
> K-048 · K-049 · K-050.
>
> ⛔ **Modul ini sengaja TIPIS.** Sebagian besar 354 berkas adalah **pendukung modul lain**, bukan
> modul tersendiri. Tugasnya memilah: mana yang sudah tercakup, mana yang jadi tipe Go sendiri, mana
> yang dibuang.
>
> ⚠️ **Satu koreksi atas penugasan** (§1.2) dan **satu pertanyaan** (§7) — diajukan, tidak diputuskan.

---

## 1. Neraca 354 — apa yang tersisa untuk modul ini

### 1.1 Pemilahan berlapis

`[terverifikasi]`

| Lapis | Jumlah | Ke mana |
| --- | ---: | --- |
| **Total EDM-only** | **354** | — |
| − Fitur pilih-tertanggung | **4** | ⛔ **DIBUANG (K-044)** — tidak diport |
| − Rule `When` | **27** | → **modul 4** (`08-spec-edm-predikat-when.md`) |
| − Sudah dibahas modul 1 / 2 / 3 / E-5 | **28** | → rujuk modul masing-masing |
| **Sisa untuk modul ini** | **295** | §2–§4 |

Dari 295 itu, sebagian besar **berjalan lewat modul yang sudah dispec** — lihat peta pemanggilan §5.

### 1.2 ⛔ Koreksi — keempat berkas pilih-tertanggung MEMANG masuk 354

> Penugasan menduga *"4 berkas itu ada di EDM+RNW, BUKAN EDM-only, jadi mungkin tak masuk 354 ini"*.
> **Dugaan itu tidak bertahan.**

`[terverifikasi]` Keempatnya **ada di dalam himpunan 354**:

| Berkas | Folder | `pxObjClass` |
| --- | --- | --- |
| `ChooseInsured` | `Harness\` | `Rule-HTML-Harness` |
| `ChooseInsuredDtl` | `Section\` | `Rule-HTML-Section` |
| `BrowseAccountInsuredEDM` | `ReportDefinition\` | `Rule-Obj-Report-Definition` |
| `SetDataInsuredEDM_Act` | `Activity\` | `Rule-Obj-Activity` |

**Sebabnya definisi:** "EDM-only" dihitung terhadap **NB saja**. Keempatnya memang ada di RNW, tetapi
**tidak ada di NB** — jadi menurut definisi resmi mereka EDM-only.

⛔ **Konsekuensi:** **K-044 memangkas 354 → 350 berkas yang diport.** Keempatnya ada di korpus,
terbaca, tetapi **tidak diimplementasikan**. Ini harus disebut eksplisit di spec agar tidak
terlihat seperti kelalaian.

### 1.3 Klasifikasi per `pxObjClass`

`[terverifikasi]` — tipe rule dibaca dari `<pxObjClass>`, **bukan** dari nama folder:

| `pxObjClass` | Jml | Folder |
| --- | ---: | --- |
| `Rule-Obj-Activity` | 94 | `Activity\` |
| `Rule-HTML-Section` | 81 | `Section\` |
| `Rule-Obj-FlowAction` | 54 | `FlowAction\` |
| `Rule-Obj-Model` | 29 | `DataTransform\` |
| `Rule-Obj-Report-Definition` | 28 | `ReportDefinition\` |
| `Rule-Obj-When` | 27 | `When\` → **modul 4** |
| **`Rule-Connect-SQL`** | **22** | `RDBList\` ⚠️ |
| `Rule-Declare-Pages` | 14 | `DataPage\` |
| `Rule-HTML-Harness` | 4 | `Harness\` |
| `Rule-Obj-Flow` | 1 | `Flow\` |

⚠️ **Jebakan folder:** seluruh **22** berkas di `RDBList\` ber-`pxObjClass` **`Rule-Connect-SQL`** —
bukan tipe RDBList. Ini **aturan di korpus ini, bukan pengecualian**. Klasifikasi wajib dari
`pxObjClass`.

---

## 2. Sembilan puluh empat Activity — kelompok fungsi

`[terverifikasi]` Dikelompokkan dari **isi** (cacah token properti + metode), bukan nama.
**Ⓟ** = menyentuh premi/perhitungan (prioritas).

| Kelompok | Jml | Status |
| --- | ---: | --- |
| **Z — sudah dibahas modul lain** | **19** | rujuk: modul 1 (9) · modul 2 (4) · modul 3 (5) · E-5 (1) |
| **A — Life / medis** | **10** | → **modul 6** |
| **D — spreading & kapasitas treaty** Ⓟ | **3** | **modul ini** |
| **E — perhitungan premi / TSI / rate** Ⓟ | **9** | **modul ini** |
| **G — lampiran / cetak / notifikasi** | 2 | modul ini, pendukung |
| **H — pendukung lain** | **51** | modul ini, §2.3 |

### 2.1 Kelompok D — spreading & kapasitas treaty Ⓟ

`[terverifikasi]`

| Berkas | Sinyal (premi / TSI / spread / coverage) |
| --- | --- |
| `Activity\CopySpreading_Act` | 52 / 54 / **247** / 265 — terbesar di kelompok |
| `Activity\SpreadingProtection` | 0 / 0 / **96** / 50 |
| `Activity\CheckLimitTreatyType_Act` | 5 / **88** / 20 / 136 |

⛔ Ketiganya menyentuh **kapasitas treaty** — jalur yang sama dengan resolver skala K-018.
`[pertanyaan terbuka]` apakah rumusnya setara dengan jalur NB; **belum dibandingkan**.

### 2.2 Kelompok E — perhitungan premi / TSI / rate Ⓟ

`[terverifikasi]`

| Berkas | premi / TSI / rate / coverage |
| --- | --- |
| `Activity\SetTSIAllCoverageLife_ACT` | 7 / **174** / 66 / 428 — ⚠️ juga bernuansa Life |
| `Activity\CopyTemplateCoverageSQL_PostActEDM` | 31 / 66 / 4 / **522** |
| `Activity\CalculatePremiumTravel` | **71** / 15 / 0 / 114 |
| `Activity\CopyTemplateCoverage` | 54 / 24 / 0 / 234 |
| `Activity\calculatePremiPA` | 37 / 30 / 7 / 190 |
| `Activity\CalculatePremiFire` | 26 / 33 / 17 / 156 |
| `Activity\SetSumTSISpreading_Act` | 0 / **92** / 6 / 96 |
| `Activity\SetTSIAllCoverage_ACT` | 2 / 55 / 2 / 69 |
| `Activity\ChangeTSITJH_Act` | 11 / 39 / 0 / 5 |

⛔ **Empat aktivitas ber-nama `Calculate*`/`calculate*` adalah EDM-only** — perhitungan premi
per-lini di endorsement punya implementasinya sendiri, terpisah dari jalur NB.

⚠️ `[pertanyaan terbuka]` apakah rumusnya setara dengan `services/premium` jalur NB. **Belum
dibandingkan** — bukan lingkup putaran ini, tetapi **wajib dibandingkan sebelum implementasi**,
karena keduanya memakai **Seam 3 yang sama**.

### 2.3 Kelompok H — 51 berkas, dipilah lebih lanjut

`[dugaan]` Bukan "tak terklasifikasi" melainkan **berbobot perhitungan rendah**. Pola dari isi:

| Sub-pola | Contoh | Sifat |
| --- | --- | --- |
| **Salin/terapkan ke seluruh objek** | `CopyDiscountAllObject_act` (premi 46) · `CopyCoverageToALL_HIO_Act` · `CopyDiscountObject` | Ⓟ ringan — menyalin angka |
| **Manipulasi coverage** | `CopyCoverageFormSection_Act` · `ChangeCoveragePremium_FacIn` (cov 207) · `copyCoveragePeriode_act` · `SetTsiCovAct` | inti layar endorsement |
| **Template coverage dari SQL** | `FindTemplateCoverageSQLTEST_PreAct` · `CopyTemplateCoverageSQL_PostAct` | pengambilan data |
| **Kait layar** (`_PreAct` / `_PostAct`) | `InputCoverageAneka_PreAct` · `CoverageList_FacIn_PreAct` | tampilan |
| **Hitung ulang saat ubah** | `HitungPremiOnChange` (rate 14) · `fillActPremiPA` | Ⓟ ringan |
| **Fac out / outgo** | `CopyOutgo_Act` · `CopyOutgoObjectMBU_Act` | jalur retro |

⚠️ Dilabeli `[dugaan]` karena pemilahan dari pola isi agregat, **belum** dibaca langkah-per-langkah.
Untuk spec, kelompok H diperlakukan sebagai **pendukung layar endorsement**, bukan modul tersendiri.

---

## 3. Delapan puluh satu Section — kelompok fungsi

`[terverifikasi]`

| Kelompok | Jml | Contoh + sinyal |
| --- | ---: | --- |
| **F — layar/komponen umum** | **34** | `PlanSection` · `CreateListBenefit` · `CreateListPlan` |
| **E — coverage / objek pertanggungan** | **17** | `InputCoverageTravel` (cov 179) · `InputCoverageGridMBU_FacIn_IsUW` |
| **A — Life / benefit / medis** | **15** | `Medical_Sec` (life 1400) · `InputMedical2_Sec` (1025) → **modul 6** |
| **D — angka premi / TSI** Ⓟ | **7** | `InputCoverageMBU_FacIn_IsUW` (premi 48, cov 505) |
| **C — spreading** Ⓟ | **6** | `InputCoverageMBU_FacIn` (cov 664) · `InputPerCoverageFire_IsUW` (cov 424) · `InputDtlSpreadingCoverage` |
| **B — retro / fac out** | 2 | `ViewCoverageFacOutShow_Is_UW` · `ViewRiCommFacOut1` |

📌 **Pola pasangan `_IsUW`** berulang: banyak section punya kembaran berakhiran `_IsUW` (versi
tampilan underwriter). Sesuai keputusan RNW, keduanya **dua komponen layar terpisah**, mengikuti
sistem lama — bukan satu komponen dengan flag.

---

## 4. Sisa tipe — kelompok ringkas

`[terverifikasi]`

| Tipe | Jml | Sifat | Contoh terbesar |
| --- | ---: | --- | --- |
| `Rule-Obj-FlowAction` | **54** | aksi layar endorsement (tombol, dialog lokal) | `Endorsement_FlowAct` · `Endorsement_FlowAct_IsUW` · `ViewDtlDeductibleFire` |
| `Rule-Obj-Model` | **29** | transformasi data (DataTransform) | `TampilGridPlanPA` · `ChangeUpKendaraan` · `SetTSIPremiSpreaded` |
| `Rule-Obj-Report-Definition` | **28** | daftar/pencarian untuk layar | `crmAccountsList` · `BrowsePremiPA` · `BrowseAccountInsuredEDM` ⛔*dibuang* |
| **`Rule-Connect-SQL`** | **22** | query — ⚠️ semuanya di folder `RDBList\` | `SearchTemplateMainDeductibleSQL` · `SearchTemplateMainCoverageSQL` · `GetStartDate` |
| `Rule-Declare-Pages` | **14** | DataPage master/lookup | `D_BrandList` · `D_ClauseList` · `D_InwardScale` |
| `Rule-HTML-Harness` | **4** | titik masuk layar | `Medical_Harnes` → modul 6 · `ViewOldData` · `SFAPortalEndorsement` ⚠️ · `ChooseInsured` ⛔*dibuang* |
| `Rule-Obj-Flow` | **1** | `InputAddendumFacultativeIn` | → **modul 3** |

⚠️ `ChangeUpKendaraan` adalah salah satu dari **5 berkas** yang bernama sama di NB **tetapi bertipe
rule berbeda** (di EDM `Rule-Obj-Model`, di NB `Rule-Obj-Activity`). Jebakan "nama sama tipe beda" —
implementasi **tidak boleh** memakai ulang yang NB.

---

## 5. Peta pemanggilan

`[terverifikasi]` Seluruh 2.061 berkas EDM dipindai untuk menyebut nama masing-masing dari 354:

| | Jumlah |
| --- | ---: |
| Dipanggil dari modul yang **sudah dispec** | **28** |
| Dipanggil, tetapi dari berkas lain | **321** |
| **Tidak disebut berkas mana pun** | **5** |

Rincian pemanggil bermodul: **modul 1 → 21 sisi rujukan** · **modul 3 → 12** · **modul 2 → 4** ·
**E-5 → 4**.

⛔ **Bacaan yang benar:** 349 dari 354 **memang dipakai**. Yang berdiri di luar jangkauan modul yang
sudah dispec bukan berarti tak terpakai — mereka dipanggil dari **lapisan tampilan** (Section dan
FlowAction), yang konsisten dengan temuan E-4 bahwa Section adalah pemanggil terbesar di EDM.

### 5.1 Lima berkas tidak disebut mana pun (E-Q24)

`[terverifikasi]`

| Folder | Berkas | `pxObjClass` |
| --- | --- | --- |
| `Activity\` | `SaveClausePA_Act` | `Rule-Obj-Activity` |
| `Harness\` | `SFAPortalEndorsement` | `Rule-HTML-Harness` |
| `RDBList\` | `GetMasterKlausulAge(1)` | `Rule-Connect-SQL` |
| `RDBList\` | `SearchClobClauseSQL(1)` | `Rule-Connect-SQL` |
| `When\` | `IsCustomBond` | `Rule-Obj-When` → modul 4 |

⛔ **TIDAK divonis usang (K-006).** "Hilang rujukan" **bukan** bukti tidak terpakai:

- **`SFAPortalEndorsement`** adalah **Harness** — titik masuk layar. `[dugaan]` dipanggil dari
  konfigurasi portal Pega yang **tidak ikut terekspor**, bukan dari rule.
- Dua berkas berakhiran **`(1)`** — `[dugaan]` pola salinan "save as" Pega. ⚠️ **Nama bukan bukti**;
  isinya belum diperiksa.

⚠️ **`[pertanyaan terbuka]` status pemakaian kelimanya** — milik work owner + IT. Sampai dijawab:
**diport apa adanya bila terbaca**, tidak dihapus, tidak di-`panic`.

---

## 6. Berkas Life EDM-only — ditandai untuk modul 6

`[terverifikasi]` **66 berkas** EDM-only bersinyal Life, tersebar di **9 tipe rule**:

| Tipe | Jml |
| --- | ---: |
| Section | 19 |
| Activity | 16 |
| FlowAction | 11 |
| ReportDefinition | 8 |
| DataPage | 6 |
| DataTransform | 3 |
| Harness · Flow · RDBList | 1 masing-masing |

Dari 66 itu, **16 bernama mengandung "Life"** — sisanya hanya terlihat dari isi (`Benefit`,
`Medical`, `Disease`, `Plan`).

⚠️ **Dua angka yang sah, bergantung penyaring:** E-4 §5 mencatat **52** dengan ambang sinyal ≥ 20
tanpa kata `Plan` pada nama; hitungan di sini **66** dengan ambang ≥ 25 **plus** `Plan`. Selisih 14
berkas seputar `Plan*`/`CreateListPlan`. **Modul 6 harus menetapkan batasnya sendiri** — tidak
diputuskan di sini.

⛔ **Detail jalur Life TIDAK dispec di modul ini.** Yang dilakukan hanya **menandai**. Modul 6
(K-044) menangani: `SetOldData` keluar bila `IsLife` · hanya Life yang jalan di `EdmType==1` · Life
melompati tangga akseptasi · Harness, DataPage, dan ReportDefinition sendiri.

---

## 7. Kontrak modul — dan pertanyaan seam

### 7.1 Sebagian besar BUKAN modul tersendiri

⛔ Inti bahan ini: **354 EDM-only bukan satu modul**. Pemetaannya ke struktur Go:

| Kelompok | Jadi apa di sistem baru |
| --- | --- |
| 19 Activity + 27 When + 1 Flow yang sudah dibahas | **bagian modul 1–4** — tidak ada tipe Go baru |
| 10 Activity + 15 Section + 41 lainnya bersinyal Life | **modul 6** |
| 3 Activity spreading Ⓟ | **bagian `services/spreading`** (modul NB yang sudah ada) |
| 9 Activity perhitungan Ⓟ | **bagian `services/premium`** — lewat **Seam 3** |
| 51 Activity pendukung + 34 Section umum + 54 FlowAction + 29 DataTransform | **lapisan tampilan & pendukung** — handler dan komponen React, bukan modul domain |
| 22 Connect-SQL + 14 DataPage + 28 ReportDefinition | **lapisan `repository`** — query dan lookup |
| 4 pilih-tertanggung | ⛔ **tidak diport (K-044)** |

📌 **Tidak ada satu pun kelompok yang menjadi modul domain baru.** Semuanya masuk ke modul yang
sudah ada atau ke lapisan tampilan/repository.

### 7.2 Seam — usulan: TIDAK ada seam baru

**Total seam sekarang 5** (K-050): `rules.Eval` · `acceptance.Next` · `premium.Calculate` ·
`endorsement.PrepareBeforeImage` · `endorsement.OpenCase`.

**Usulan saya: tidak menambah seam.** Alasannya:

- Kelompok Ⓟ (perhitungan premi, spreading) teruji lewat **Seam 3** yang sudah ada — mereka rumus,
  dan Seam 3 adalah pintu perhitungan.
- Kelompok pendukung/tampilan **bukan perilaku domain** — mengujinya lewat seam akan mengunci bentuk
  UI yang belum dirancang.
- Kelompok repository (query, DataPage, ReportDefinition) diuji lewat **seam repository** yang
  sendirinya **belum ada** dan menunggu skema Oracle — kewajiban tertunda yang sudah tercatat sejak
  spec NB, bukan seam baru untuk modul ini.

⚠️ **Satu hal yang perlu keputusan:** kesembilan Activity `Calculate*` kelompok E **belum
dibandingkan** dengan rumus NB. Bila ternyata rumusnya **berbeda**, Seam 3 harus menerima varian
endorsement per lini — perluasan kontrak, bukan seam baru.

⛔ **Diajukan, tidak diputuskan:** apakah perbandingan rumus `Calculate*` EDM vs NB dilakukan
**sebelum** spec ditulis (aman, menunda) atau **saat implementasi** (cepat, berisiko spec berubah)?

---

## 8. Kejanggalan dan kode usang

| # | Butir | Sikap |
| ---: | --- | --- |
| 1 | **5 berkas tanpa rujukan** (§5.1) | ⛔ **tidak divonis usang** (K-006); `[pertanyaan terbuka]` |
| 2 | **`ChangeUpKendaraan`** — nama sama di NB, **tipe rule berbeda** | jangan pakai ulang implementasi NB |
| 3 | **22 berkas `RDBList\` bertipe `Rule-Connect-SQL`** | klasifikasi dari `pxObjClass`, bukan folder |
| 4 | **Dua berkas berakhiran `(1)`** — dugaan salinan "save as" | isinya belum diperiksa; **nama bukan bukti** |
| 5 | **Pasangan `_IsUW`** — dua komponen layar terpisah | ikut sistem lama, keputusan sadar |
| 6 | Kode usang lain yang tersentuh | **diport apa adanya** (K-046) |

---

## 9. Out of Scope

1. **27 rule `When` EDM-only** — **modul 4** (`08-spec-edm-predikat-when.md`).
2. **Jalur produksi EDM** — **E-5**, ditunda menunggu tabel flat + `ALL_SOURCE` dari DBA.
3. **Before-image dan selisih** — modul 1 dan 2.
4. **Alur masuk** — modul 3.
5. **Jalur Life** — **modul 6** (K-044); modul ini hanya menandai.
6. **Fitur pilih-tertanggung** — ⛔ **K-044: DIBUANG**, 4 berkas tidak diport meski masuk 354.
7. **Perbandingan rumus `Calculate*` EDM vs NB** — §7.2, belum dilakukan.
8. **Tangga akseptasi** — klaim arsip §4.3 tetap **`[belum diuji]`**.

---

## 10. Pertanyaan terbuka

| # | Pertanyaan | Pemilik |
| --- | --- | --- |
| **E-Q24** | Status pemakaian 5 berkas tanpa rujukan — khususnya `SFAPortalEndorsement` (Harness, mungkin dari konfigurasi portal di luar ekspor) | work owner + IT |
| **baru** | Kapan rumus `Calculate*` EDM dibandingkan dengan NB — sebelum spec atau saat implementasi (§7.2) | **work owner** |
| **baru** | Batas himpunan Life untuk modul 6 — **52** atau **66**, tergantung penyaring (§6) | **work owner / modul 6** |
| terbuka | Apakah rumus spreading EDM (`CopySpreading_Act`, `CheckLimitTreatyType_Act`) setara jalur NB | belum dibandingkan |

---

## 11. Kesiapan

✅ **Bahan lengkap untuk `/to-spec`:**

| Butir | Status |
| --- | --- |
| **Neraca berlapis** 354 → 4 dibuang → 27 modul 4 → 28 modul 1/2/3/E-5 → **295 sisa** | §1.1 |
| ⛔ **Koreksi:** 4 pilih-tertanggung **MASUK 354**, K-044 memangkas jadi **350 diport** | §1.2 |
| Klasifikasi `pxObjClass` + jebakan folder | §1.3 |
| **94 Activity** — 6 kelompok, prioritas Ⓟ ditandai, yang sudah dibahas ditandai | §2 |
| **81 Section** — 6 kelompok | §3 |
| Sisa tipe (FlowAction 54, Model 29, RD 28, SQL 22, DataPage 14, Harness 4, Flow 1) | §4 |
| **Peta pemanggilan** — 28 dari modul dispec · 321 dari berkas lain · **5 tanpa rujukan** | §5 |
| **66 berkas Life** ditandai untuk modul 6, dengan catatan dua angka sah | §6 |
| **Kontrak:** tidak ada modul domain baru; pemetaan ke struktur Go | §7.1 |
| **Usulan: tidak ada seam baru** + 1 pertanyaan | §7.2 |

⛔ **`/to-spec` dan `/to-tickets` TIDAK dijalankan.**

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
