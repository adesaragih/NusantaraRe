# Audit 10 — Rule kembar: nama sama, **rule berbeda**

> **21 September 2026.** Generalisasi temuan `CountRateRetroCov` (K-060 amandemen) ke **seluruh
> korpus**. ⛔ **Mengubah tafsir angka drift yang sudah tercatat.**

---

## 1. Aturan identitas

⛔ **Identitas rule ditentukan `pzInsKey` dan `pyClassName` — BUKAN nama berkas.**
Dua berkas bernama sama dengan basis `pzInsKey` berbeda adalah **DUA RULE**, bukan dua versi.

⚠️ **`pzInsKey` memuat stempel waktu per-penyimpanan.** Membandingkan kunci utuh akan menandai
hampir semuanya berbeda. Identitas dibaca dari **basisnya**:

```
RULE-OBJ-ACTIVITY ASM-FW-GISFW-DATA-PROPERTYITEM COUNTRATERETROCOV #20260921T095254.245 GMT
└─────────────────────── basis (identitas) ──────────────────────┘ └──── stempel (dibuang) ────┘
```

`[terverifikasi]` Contohnya: salinan `DDL\` dan salinan korpus varian FIRE punya kunci utuh berbeda
(`#20260921…` vs `#20260619…`) tetapi **basis sama** — satu rule, dua salinan.

---

## 2. Tiga golongan

| Golongan | Kriteria | Artinya |
| :-: | --- | --- |
| **G1** | basis `pzInsKey` **sama** + kelas sama + isi sama | satu rule, tidak ada masalah |
| **G2** | basis `pzInsKey` **sama** + kelas sama + **isi BEDA** | satu rule, **drift sejati** — kandidat ekspor ulang |
| **G3** | basis `pzInsKey` **BEDA** atau kelas **BEDA** | ⛔ **RULE KEMBAR** — keduanya harus diport, **bukan dipilih** |

---

## 3. Hasil

`[terverifikasi]` Seluruh nama berkas yang hadir di lebih dari satu folder korpus, per subfolder, per
pasangan folder:

| Golongan | Pasangan | % |
| :-: | ---: | ---: |
| **G1** kunci sama + isi sama | **4.590** | 86,8 % |
| **G2** kunci sama + isi BEDA | **483** | 9,1 % |
| **G3** kunci/kelas BEDA | **215** | **4,1 %** |
| **Total** | **5.288** | 100 % |

Pemeriksaan silang: `4.590 + 483 + 215 = 5.288` ✅ — cocok jumlah pasangan di
`07-hitung-berkas-dan-drift.md` §3. Berkas tanpa `pzInsKey` terbaca: **0**.

**Identitas nama unik dalam G3: 109.**

### ⛔ 3.1 Tafsir angka 526 berubah

`[terverifikasi]` Dari **526** pasangan yang tercatat "versi sama tetapi isi beda":

| | Jumlah |
| --- | ---: |
| Ternyata **G3 — rule kembar**, bukan drift | **93** |
| **Drift sejati** yang tersisa | **433** |

⚠️ **93 pasangan yang selama ini dihitung sebagai drift sebenarnya rule berbeda.** Keduanya
**wajib diport**; mengekspor ulang salah satunya **tidak menyelesaikan apa pun**, dan memilih salah
satu **menghilangkan perilaku**.

📌 Angka **526** di `07-hitung-berkas-dan-drift.md` dan `00-KEPUTUSAN-WORK-OWNER.md` K-058
**tidak salah sebagai pengukuran** — yang berubah adalah **tafsirnya**. Saya tidak menyuntingnya di
sana; dokumen ini yang menjadi rujukan tafsir.

---

## 4. Daftar G3 — 109 nama unik

⚠️ Sebagian besar muncul **dua kali** (pasangan `NB↔EDM` dan `RNW↔EDM`), karena NB dan RNW identik.

### 4.1 `When\` — **42 nama**, hampir seluruh keluarga predikat COB

`FlagOldData` · `IsAneka` · `IsBoiler` · `IsBondingAndCustomBonds` · `IsBondingKBG` · `IsCIS` ·
`IsCIT` · `IsClaim` · `IsContractorsPlantMachinery` · `IsCorporate` · `IsCrime` · `IsCustomBonds` ·
`IsEar` · `IsEngineering` · `IsEnvironmental` · **`IsFacRetro`** · **`IsFire`** · `IsFireStyle1` ·
`IsFireStyle2` · `IsFlagDelete` · `IsFlagOldData` · `IsGlass` · `IsGolfInsurance` · `IsGroupFac` ·
`IsGrowingTrees` · `IsHE` · `IsKPR` · `IsLandRig` · `IsLiability` · **`IsLife`** · `IsMarineCargo` ·
`IsMBD` · **`IsMBU`** · `IsMBUCar` · `IsNotEDM` · `IsOilGas` · **`IsPA`** · `IsPASSG` ·
`IsProposalTransfer` · `IsTravel` · **`IsUW`**

**Contoh terverifikasi:**

| Rule | NB / RNW | Endorsment |
| --- | --- | --- |
| `IsFire` | kelas `ASM-FW-GISFW-Work`<br>kondisi `Rule IsKPR evaluates to true` | kelas **`ASM-SFAGIS-Work-Endorsement`**<br>kondisi `OutData.pxResults(1).CARI2 = "KPR"` |
| `IsPA` | kelas `Data-Party-Person`<br>kondisi `pyWorkPage.Quotation.BusinessType = "PA"` | kelas **`ASM-SFAGIS-Work-Endorsement`**<br>kondisi `OutData.pxResults(1).CARI2 = "PA"` |
| `IsFacRetro` | kelas `ASM-FW-GISFW-Data-OfferFacIn`<br>kondisi `.IsFacRetro = 1` | kelas **`ASM-FW-GISFW-Work`**<br>kondisi `pyWorkPage.OfferFacIn.IsFacRetro = 1` |
| `IsUW` | kelas `ASM-FW-GISFW-Work` | kelas **`ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance`** |

📌 **Ini menjelaskan K-050 dari sudut lain.** K-050 mencatat bahwa EDM membaca COB lewat
`OutData.pxResults(1).CARI2` sedangkan NB membaca properti kasus. Yang dulu dibaca sebagai *"satu
rule dengan dua sumber data"* ternyata **dua rule di dua kelas**.

> ⚠️ **K-050 TIDAK dibatalkan di sini.** Keputusannya (satu registry, sumber data per-rule) adalah
> kewenangan work owner. Yang berubah adalah **premisnya**: "196 predikat identik" perlu diukur ulang
> dengan aturan identitas §1. **Menyerahkan ini sebagai bahan, bukan keputusan.**

### 4.2 `Activity\` — 25 nama

`AddCurencyList_ACT` · `CalculateNetRate_ACT` · `CallCountPremiFacIn_act` · `cekSpreadingFactIn` ·
`CopyCoverage_Act` · **`CountRateRetroCov`** · `fillActPremiTravel` · `FillPaymentInstallment` ·
`GetCurrencyMaster` · `GetLowestPctLimit_ACT` · `GetZipCodeFromAddress` · `GISInitAttach` ·
`InputQuotation_PreAct` · `InputRateFO` · `NegativeIsNotAllowed` · `SearchClauseFireSQL_PostAct` ·
`SearchClauseFireSQL_PreAct` · `SearchRiskAddressAct` · `serviceInsertArasapas_act` ·
`SetCategoryAttach` · `SetCategoryAttachment_Reas` · `SetIdxCargo` · `SetSurveyReport_Act` ·
`SetValidateInstallment_Act` · `SumCurrencyListAllRetro_Act`

⚠️ **Empat di antaranya menyentuh uang** dan ada di daftar ekspor ulang: `CalculateNetRate_ACT`,
`CallCountPremiFacIn_act`, `InputRateFO`, `CountRateRetroCov`.

**Contoh:** `CalculateNetRate_ACT` — NB/RNW kelas `ASM-FW-GISFW-Data-PropertyItem`, Endorsment kelas
**`ASM-FW-GISFW-Data-Coverage`**, versi **sama** (`01-01-73`), isi beda. Ini **rule kembar**, bukan
drift.

### 4.3 `Section\` — 18 nama

`ChooseClauseFire` · `CoverageList` · `CreateListClause` · `HistoricalSurveyReportDtlUW` ·
`InputCoverageAneka_FacIn` · `InputCoverageFire_IsUW` · `InputDtlCoverage_FacIn` ·
`InputDtlCoverage_FacIn_IsUW` · `SelClauseList` · `ShowCoverageFacOut_IsUW` · `UploadMember` ·
`ViewCoverage` · `ViewCoverageFacOutShow` · `ViewCoverageFire` · `ViewCoveragePA` ·
`ViewCoverageTravel` · `ViewDeliveryAddressFacOut` · `ViewGeneralPolisFacOut` ·
`ViewObjectMarineCargoFacOut` · `ViewObjectPAFacOut` · `ViewPerilsFacOut`

### 4.4 Tipe lain — 24 nama

**DataTransform (8):** `CargoPosition_DT` · `CountPremiEDMFacOut_DT` · `InputDtlWarrantyList_DT` ·
`SetCodeRiskExposure_DT` · `SetFlagDelete_DT` · `SetIndexObject_DT` · `SetObjNoDT` ·
`SystemSetOneYear_DT`
**ReportDefinition (6):** `BrowseColor_RD` · `BrowseDistrict_RD` · `BrowseMarineCondition_RD` ·
`BrowseProvince_RD` · `DataTableEditorReport` · `pyDefaultReport`
**Harness (5):** `CedingCompany` · `ChooseClauseFire` · `HistoricalSurveyReportUW` ·
`ViewCoverageFacOut` · `ViewCoverageFacOut_isUW`
**FlowAction (2):** `InputDtlPaymentFlow_FacIn` · `ViewObjectDetailOccupation`
**RDBList (1):** `GetMasterKlausulAge`

---

## 5. Konsekuensi

1. ⛔ **Jangan mengekspor ulang G3.** Mengekspor ulang tidak menyelesaikan apa pun — keduanya memang
   rule berbeda. Daftar ekspor ulang (`08-daftar-ekspor-ulang-prioritas.md`) **hanya untuk G2**.
2. ⛔ **Jangan memilih salah satu G3.** Keduanya diport; pemilihannya terjadi saat runtime lewat
   **kelas halaman**.
3. ⚠️ **Komentar ketertelusuran (`CLAUDE.md` §4.6) wajib menyebut KELAS**, bukan hanya nama rule.
   `// Asal: Activity\CountRateRetroCov` **tidak cukup** — ada dua.
4. ⚠️ **Angka drift lintas folder mana pun yang diukur tanpa memeriksa `pzInsKey` menggabungkan dua
   hal berbeda.** Berlaku bagi angka 526 dan 354.

---

## 6. `[pertanyaan terbuka]`

1. **Apakah 42 rule kembar di `When\` mengubah premis K-050** ("196 predikat identik"). K-050 adalah
   keputusan work owner; **tidak diputuskan di sini**.
2. **Apakah metode deteksi drift diganti** ke hash 23 tag — **tetap terbuka di K-058**. Audit ini
   menambah syarat baru: apa pun metodenya, ia **wajib memeriksa `pzInsKey` + `pyClassName` lebih
   dulu**.
3. `[dugaan]` Sebagian G3 mungkin **bukan** kembar disengaja melainkan rule yang dipindahkan kelasnya
   antar siklus. Membedakan keduanya **tidak dapat dilakukan dari korpus**.

---

## 7. Perintah audit

```powershell
# menghasilkan 5288 / G1 4590 / G2 483 / G3 215 / 109 nama unik / 93 dari 526
# Get-Info = kontrak 23 tag (lihat 07-hitung-berkas-dan-drift.md §5) DITAMBAH:
#   $key   = <pzInsKey> pertama
#   $cls   = <pyClassName> pertama
#   $basis = $key dengan sufiks stempel dibuang:
$rxBasis=[regex]'\s*#\s*\d{8}T\d{6}.*$'
$basis = $rxBasis.Replace($key,'')
# lalu per pasangan berkas bernama sama:
#   if($basisA -cne $basisB -or $clsA -cne $clsB){ G3 }
#   elseif($hashA -eq $hashB){ G1 } else { G2 }
# CATATAN: perbandingan basis dan kelas WAJIB -ceq / -cne (peka huruf).
```

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
