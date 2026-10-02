# Audit 08 — Daftar prioritas rule yang perlu diekspor ulang

> **21 September 2026.** Latar: salinan korpus `CountRateRetroCov` ternyata **basi pada rumus uang**
> (K-060), dan ada **526 pasangan lain** yang `pyRuleSetVersion`-nya sama tetapi isinya berbeda.
> Dokumen ini menyaring mana yang layak diekspor ulang **lebih dulu**.

> # ⛔ DIPERBARUI — 39 dari 60 baris DIKELUARKAN
>
> Setelah aturan identitas **`pzInsKey` + `pyClassName`** diterapkan (`10-audit\10-rule-kembar-pzinskey.md`),
> **39 dari 60 baris** di tabel §3 ternyata **G3 — RULE KEMBAR**, bukan drift.
>
> ⛔ **Rule kembar TIDAK perlu diekspor ulang.** Mengekspor ulang tidak menyelesaikan apa pun:
> keduanya memang rule berbeda, dan **keduanya harus diport**. Memilih salah satu **menghilangkan
> perilaku**.
>
> **Daftar ekspor ulang hanya untuk G2** — drift sejati pada satu rule yang sama.
>
> | | Jumlah |
> | --- | ---: |
> | Baris semula di §3 | **60** |
> | ⛔ **Dikeluarkan (G3 — rule kembar)** | **39** |
> | ✅ **Tersisa untuk ekspor ulang (G2)** | **21** |
>
> Yang dikeluarkan pindah ke **`10-audit\10-rule-kembar-pzinskey.md`** — ditangani sebagai
> **dua rule yang sama-sama diport**, bukan sebagai kandidat ekspor ulang.
>
> ⚠️ **Tabel §3 lama tidak dihapus** (`PANDUAN-KERJA` §7); tiap baris kini bertanda golongan.
> **Yang berlaku adalah §3A.**

---

## 0. ⛔ Bukti bahwa masalahnya nyata — divergensi 10× pada premi retro

`[terverifikasi]` Empat salinan `CountRateRetroCov`, **pembagi pada `.PremiumRetro`**:

| Salinan | `pyRuleSetVersion` | Pembagi | Jumlah rumus `.PremiumRetro` |
| --- | --- | ---: | ---: |
| `NB FacIn\Activity\` | `01-01-95` | **100.000** | 2 |
| `RNW Fac In\Activity\` | `01-01-95` | **100.000** | 2 |
| **`Endorsment Fac In\Activity\`** | `01-01-95` | ⛔ **10.000** | **1** |
| **`DDL\` (ekspor ulang, berlaku K-060)** | `01-01-01` | **100.000** | 2 |

⛔ **Salinan Endorsement membagi dengan 10.000 — satu orde lebih kecil**, dan **tidak memuat rumus
koreksi** langkah `10.1` sama sekali. Ketiga salinan korpus ber-**nomor versi identik `01-01-95`**,
sehingga perbedaan sebesar ini **tidak terlihat** dari nomor versi.

```powershell
# menghasilkan tabel di atas
foreach($p in @('D:\migrasi\RNM\NB FacIn\Activity\CountRateRetroCov.xml',
                'D:\migrasi\RNM\RNW Fac In\Activity\CountRateRetroCov.xml',
                'D:\migrasi\RNM\Endorsment Fac In\Activity\CountRateRetroCov.xml',
                'D:\migrasi\RNM\DDL\CountRateRetroCov.xml')){
  $t=[IO.File]::ReadAllText($p)
  $v=[regex]::Match($t,'<pyRuleSetVersion>([^<]*)</pyRuleSetVersion>').Groups[1].Value
  $c5=([regex]::Matches($t,'\),100000,20\)')).Count
  $c4=([regex]::Matches($t,'\),10000,20\)')).Count
  "$([IO.Path]::GetFileName($p)) v=$v  100000=$c5  10000=$c4" }
```

📌 **K-060 menetapkan salinan `DDL\` berlaku**, sehingga pembagi **100.000** yang dipakai. Divergensi
ini **tidak** mengubah keputusan itu — ia **membenarkannya**, dan menunjukkan mengapa daftar ini perlu.

---

## 1. Metode penyaringan

Dari **526** pasangan **versi-sama-isi-beda** (lihat `07-hitung-berkas-dan-drift.md` §3), tiap pasangan
di-`Compare-Object` pada **himpunan baris ter-normalisasi kontrak 23 tag**, lalu baris yang berbeda
diklasifikasikan:

| Kelas | Kriteria |
| :-: | --- |
| **a** | **aritmetika / uang** — baris beda pada `<PropertiesValue>`/`<PropertiesName>`/`<pyValue>` yang memuat `@Math.divide`, `@divide`, atau nama properti ber-`Premi`/`Rate`/`TSI`/`Share`/`Comm`/`Discount`/`Prorate` |
| **b** | **gerbang alur** — baris beda pada `pyStepsPreCondParamsWhen`/`WhenTrue`/`WhenFalse`, atau kondisi rule `When` (`pyConditionValue1String`, `pyConditionString`, `pyLogic`) |
| **c** | **jalur produksi** — nama berkas `Insert*`/`Save*`/`*Produc*`, **atau** baris beda memuat `INSERT INTO` / `UPDATE` / `DELETE FROM` |
| **d** | **tipe rule yang seluruhnya berbeda** — `DecisionTable`, `Flow`, `DecisionTree` |

**Urutan prioritas:** a → b → c → d; di dalam kelas, menurun menurut **jumlah baris berbeda**.

---

## 2. Hasil saringan

`[terverifikasi]`

| Ukuran | Nilai |
| --- | ---: |
| Pasangan versi-sama-isi-beda | **526** |
| **Lolos saringan a/b/c/d** | **87** |
| Ditampilkan di tabel §3 | **60 teratas** |
| **Tidak ditampilkan** | **27** |

**Sebaran kelas** (satu pasangan dapat masuk lebih dari satu kelas):

| Kelas | Pasangan |
| :-: | ---: |
| **a** aritmetika/uang | **22** |
| **b** gerbang alur | **35** |
| **c** jalur produksi | **18** |
| **d** tipe rule | **24** |

**Sebaran per tipe rule:** Activity **32** · DecisionTable **20** · Section **12** · When **11** ·
RDBList **6** · Flow **2** · DecisionTree **2** · ReportDefinition **2**.

⚠️ **Pasangan folder muncul berpasangan.** Karena NB↔RNW **nol berbeda**, hampir setiap berkas muncul
dua kali: `NB vs EDM` **dan** `RNW vs EDM`. Keduanya menunjuk **satu berkas EDM yang sama**, jadi
**satu ekspor ulang menutup dua baris**. Secara berkas unik, 60 baris ≈ **30 berkas**.

---

## 3. Enam puluh teratas

| # | Berkas | Tipe | Pasangan | Kelas | Δbaris | Cuplikan perbedaan |
| ---: | --- | --- | --- | :-: | ---: | --- |
| 1 | `ViewCoverageFire` | Section | RNW↔EDM | ab | 10.703 | `<pyLogic>And</pyLogic>` |
| 2 | `ViewCoverageFire` | Section | NB↔EDM | ab | 10.703 | `<pyLogic>And</pyLogic>` |
| 3 | `ViewCoverage` | Section | NB↔EDM | ab | 4.917 | `<pyLogic>And</pyLogic>` |
| 4 | `ViewCoverage` | Section | RNW↔EDM | ab | 4.917 | `<pyLogic>And</pyLogic>` |
| 5 | `ViewCoveragePA` | Section | RNW↔EDM | a | 3.428 | `<pyValue>.Discount</pyValue>` |
| 6 | `ViewCoveragePA` | Section | NB↔EDM | a | 3.428 | `<pyValue>.Discount</pyValue>` |
| 7 | `ViewCoverageTravel` | Section | RNW↔EDM | a | 2.439 | `<pyValue>.Discount</pyValue>` |
| 8 | `ViewCoverageTravel` | Section | NB↔EDM | a | 2.439 | `<pyValue>.Discount</pyValue>` |
| 9 | `ViewPerilsFacOut` | Section | RNW↔EDM | a | 1.526 | `<pyValue>.Rate</pyValue>` |
| 10 | `ViewPerilsFacOut` | Section | NB↔EDM | a | 1.526 | `<pyValue>.Rate</pyValue>` |
| **11** | **`CountRateRetroCov`** | Activity | RNW↔EDM | ab | 522 | `@Math.divide((Local.ShareOffered*.Rate * Local.prorate),`**`10000`**`,20)` |
| **12** | **`CountRateRetroCov`** | Activity | NB↔EDM | ab | 522 | idem — ⛔ **pembagi 10× berbeda**, lihat §0 |
| 13 | `CallCountPremiFacIn_act` | Activity | NB↔EDM | ab | 511 | `<PropertiesName>.TotalGrossPremi</PropertiesName>` |
| 14 | `CallCountPremiFacIn_act` | Activity | RNW↔EDM | ab | 511 | idem |
| 15 | `ViewGeneralPolisFacOut` | Section | NB↔EDM | a | 505 | `<pyValue>.SumOfTSI</pyValue>` |
| 16 | `ViewGeneralPolisFacOut` | Section | RNW↔EDM | a | 505 | idem |
| 17 | `CalculateNetRate_ACT` | Activity | NB↔EDM | ab | 254 | `@divide(.Rate,Local.TotalRate,4)*Local.TotalNetrate` |
| 18 | `CalculateNetRate_ACT` | Activity | RNW↔EDM | ab | 254 | idem |
| 19 | `CountEdmAdjTSI_Act` | Activity | NB↔EDM | ab | 116 | `<PropertiesName>Local.TSI</PropertiesName>` |
| 20 | `CountEdmAdjTSI_Act` | Activity | RNW↔EDM | ab | 116 | idem |
| 21 | `InputRateFO` | Activity | RNW↔EDM | a | 90 | `@divide((@toDecimal(.CoverageList(1).FacOutObjectList(1).ShareOffered)*…` |
| 22 | `InputRateFO` | Activity | NB↔EDM | a | 90 | idem |
| 23 | `SearchClauseFireSQL_PreAct` | Activity | NB↔EDM | b | 1.334 | `<pyStepsPreCondParamsWhen/>` |
| 24 | `SearchClauseFireSQL_PreAct` | Activity | RNW↔EDM | b | 1.334 | idem |
| 25 | `BrowseMarineCondition_RD` | ReportDefinition | RNW↔EDM | b | 246 | `<pyLogicLabel>A</pyLogicLabel>` |
| 26 | `BrowseMarineCondition_RD` | ReportDefinition | NB↔EDM | b | 246 | idem |
| 27 | `Protection_Act` | Activity | RNW↔EDM | b | 184 | `<pyStepsPreCondParamsWhen/>` |
| 28 | `Protection_Act` | Activity | NB↔EDM | b | 184 | idem |
| 29 | `SetIdxCargo` | Activity | NB↔EDM | b | 179 | `<pyStepsPreCondParamsWhen/>` |
| 30 | `SetIdxCargo` | Activity | RNW↔EDM | b | 179 | idem |
| 31 | `serviceInsertArasapasEDM_act` | Activity | NB↔EDM | b | 178 | `<pyStepsPreCondParamsWhen>IsFacRetro</pyStepsPreCondParamsWhen>` |
| 32 | `serviceInsertArasapasEDM_act` | Activity | RNW↔EDM | b | 178 | idem |
| 33 | `NegativeIsNotAllowed` | Activity | RNW↔EDM | b | 124 | `.SurroundingRisk.BackDistance<0` |
| 34 | `NegativeIsNotAllowed` | Activity | NB↔EDM | b | 124 | idem |
| **35** | **`IsFacRetro`** | When | NB↔EDM | b | 118 | `<pyConditionString>pyWorkPage.OfferFacIn.IsFacRetro = 1</pyConditionString>` |
| **36** | **`IsFacRetro`** | When | RNW↔EDM | b | 118 | idem — ⚠️ predikat pemicu Fac Out (K-057) |
| 37 | `IsMBD` | When | RNW↔EDM | b | 106 | `<pyConditionString>BusinessType = "MBD"</pyConditionString>` |
| 38 | `IsMBD` | When | NB↔EDM | b | 106 | idem |
| 39 | `IsProposalTransfer` | When | NB↔EDM | b | 91 | `<pyConditionString>Transfer Proposal = "1"</pyConditionString>` |
| 40 | `FlagOldData` | When | RNW↔EDM | b | 72 | `<pyConditionValue1StringLabel>.FlagOldData = 1</…>` |
| 41 | `FlagOldData` | When | NB↔EDM | b | 72 | idem |
| 42 | `IsFlagOldData` | When | NB↔EDM | b | 72 | idem |
| 43 | `IsFlagOldData` | When | RNW↔EDM | b | 72 | idem |
| 44 | `IsPASSG` | When | NB↔EDM | b | 51 | `@Utilities.SizeOfPropertyList(.PersonListPA) > 0` |
| 45 | `IsPASSG` | When | RNW↔EDM | b | 51 | idem |
| 46 | `SaveFacinProdFireNB_Act` | Activity | RNW↔EDM | c | 72 | `<pyStepPageReference>RH_1.pySteps(1)</…>` |
| 47 | `SaveFacinProdFireNB_Act` | Activity | NB↔EDM | c | 72 | idem |
| 48 | `InsertFacoutProductionEDM` | Activity | NB↔EDM | c | 21 | `</pyStepsParamUIRemoved>` |
| 49 | `InsertFacoutProductionEDM` | Activity | RNW↔EDM | c | 21 | idem |
| 50 | `SaveFacinProdEDMMarineCargo_Act` | Activity | NB↔EDM | c | 19 | `<pxHighSeverityWarningCount>2</…>` |
| 51 | `SaveFacinProdEDMMarineCargo_Act` | Activity | RNW↔EDM | c | 19 | idem |
| 52 | `InsertFacoutProduction` | Activity | NB↔EDM | c | 18 | `</pyParamArray>` |
| 53 | `InsertFacoutProduction` | Activity | RNW↔EDM | c | 18 | idem |
| 54 | `InsertFacoutProductionEDMCurr` | Activity | NB↔EDM | c | 6 | `<pxWarningJustifiedTime>` |
| 55 | `InsertFacoutProductionEDMCurr` | Activity | RNW↔EDM | c | 6 | idem |
| 56 | `DeleteDataProduction` | RDBList | RNW↔EDM | c | 2 | `<pyStepPageReference>RH_1</…>` |
| 57 | `DeleteDataProduction` | RDBList | NB↔EDM | c | 2 | idem |
| 58 | `InsertHistoryAkseptasiPega_Sql` | RDBList | RNW↔EDM | c | 2 | idem |
| 59 | `InsertHistoryAkseptasiPega_Sql` | RDBList | NB↔EDM | c | 2 | idem |
| 60 | `INSERTJSON_JSONPOLISMONITORING_FACIN` | RDBList | RNW↔EDM | c | 2 | idem |

⚠️ **Δbaris besar tidak selalu berarti perbedaan besar.** Section berukuran puluhan ribu baris dapat
berbeda ribuan baris karena satu blok tata letak. **Kolom "kelas" lebih menentukan daripada Δbaris** —
baris 11–12 hanya Δ522 tetapi memuat divergensi pembagi 10×.

---

## 3A. ✅ DAFTAR EKSPOR ULANG YANG BERLAKU — **21 baris G2**

`[terverifikasi]` Hanya baris ber-**basis `pzInsKey` dan `pyClassName` SAMA** di kedua salinan. Kolom
kelas dicantumkan sebagai bukti.

| # | Berkas | Tipe | `pyClassName` (kedua salinan) | Kelas | Δbaris |
| ---: | --- | --- | --- | :-: | ---: |
| 1 | `CountEdmAdjTSI_Act` | Activity | `ASM-FW-GISFW-Data-Aneka` | ab | 116 |
| 2 | `CountEdmAdjTSI_Act` | Activity | idem *(pasangan RNW↔EDM)* | ab | 116 |
| 3 | `SaveFacinProdFireNB_Act` | Activity | `ASM-FW-GISFW-Work` | c | 72 |
| 4 | `SaveFacinProdFireNB_Act` | Activity | idem *(RNW↔EDM)* | c | 72 |
| 5 | `Protection_Act` | Activity | `ASM-FW-GISFW-Work` | b | 184 |
| 6 | `Protection_Act` | Activity | idem *(RNW↔EDM)* | b | 184 |
| 7 | `serviceInsertArasapasEDM_act` | Activity | `ASM-FW-GISFW-Work` | b | 178 |
| 8 | `serviceInsertArasapasEDM_act` | Activity | idem *(RNW↔EDM)* | b | 178 |
| 9 | `InsertFacoutProductionEDM` | Activity | `ASM-FW-GISFW-Data-FacOffer` | c | 21 |
| 10 | `InsertFacoutProductionEDM` | Activity | idem *(RNW↔EDM)* | c | 21 |
| 11 | `SaveFacinProdEDMMarineCargo_Act` | Activity | `ASM-FW-GISFW-Work` | c | 19 |
| 12 | `SaveFacinProdEDMMarineCargo_Act` | Activity | idem *(RNW↔EDM)* | c | 19 |
| 13 | `InsertFacoutProduction` | Activity | `ASM-FW-GISFW-Data-FacOffer` | c | 18 |
| 14 | `InsertFacoutProduction` | Activity | idem *(RNW↔EDM)* | c | 18 |
| 15 | `InsertFacoutProductionEDMCurr` | Activity | `ASM-FW-GISFW-Data-FacOffer` | c | 6 |
| 16 | `InsertFacoutProductionEDMCurr` | Activity | idem *(RNW↔EDM)* | c | 6 |
| 17 | `DeleteDataProduction` | RDBList | `ASM-FW-GISFW-Int-policyjson` | c | 2 |
| 18 | `DeleteDataProduction` | RDBList | idem *(RNW↔EDM)* | c | 2 |
| 19 | `InsertHistoryAkseptasiPega_Sql` | RDBList | `ASM-FW-GISFW-int-policyjson` | c | 2 |
| 20 | `InsertHistoryAkseptasiPega_Sql` | RDBList | idem *(RNW↔EDM)* | c | 2 |
| 21 | `INSERTJSON_JSONPOLISMONITORING_FACIN` | RDBList | `ASM-FW-GISFW-Int-policyjson` | c | 2 |

**Berkas unik: 11.** Sisanya adalah pasangan folder kedua atas berkas yang sama.

⚠️ `InsertHistoryAkseptasiPega_Sql` ber-kelas **`ASM-FW-GISFW-int-policyjson`** dengan **`i` kecil**,
berbeda dari `Int` pada dua RDBList lain. `[terverifikasi]` Perbandingan kelas dilakukan **peka
huruf**, dan keduanya tetap **sama antar folder**, jadi ia tetap G2. **Ejaan itu diport apa adanya.**

### ⛔ 39 baris yang DIKELUARKAN — rule kembar

Seluruh peringkat 1–18, 21–26, 29–30, dan 33–45 tabel §3 lama. Termasuk:

| Berkas | Kelas NB/RNW | Kelas Endorsment |
| --- | --- | --- |
| `CountRateRetroCov` | `Data-PropertyItem` | `Data-Aneka` |
| `CalculateNetRate_ACT` | `Data-PropertyItem` | `Data-Coverage` |
| `CallCountPremiFacIn_act` | `Data-Aneka` | `Data-PropertyItem` |
| `InputRateFO` | `Data-PropertyItem` | `Data-Aneka` |
| `ViewCoverageFire` | `Data-Property` | `Data-Coverage` |
| `ViewCoverage` | `Data-Cargo` | `Data-Coverage` |
| `ViewCoveragePA` · `ViewCoverageTravel` | `Data-Party-Person` | `Data-Coverage` |
| `ViewPerilsFacOut` | `Data-FacOffer` | `Data-Vehicle` |
| `ViewGeneralPolisFacOut` | `GISFW-Work` | `Data-Policy` |
| `SearchClauseFireSQL_PreAct` | `GISFW-Work` | `Data-Clause` |
| `BrowseMarineCondition_RD` | `Int-MARINECONDITION` | `Int-CONDITION` |
| `SetIdxCargo` | `Int-SHIP` | `Data-Cargo` |
| `NegativeIsNotAllowed` | `Data-OfferFacIn-LocationReinsurance` | `Data-Property` |
| `IsFacRetro` | `Data-OfferFacIn` | `GISFW-Work` |
| `IsMBD` · `IsPASSG` · `IsProposalTransfer` | beragam | `GISFW-Work` |
| `FlagOldData` · `IsFlagOldData` | `Data-Accessory` | `Data-Coverage` |

⛔ **Keempat berkas kelas a (uang) di peringkat teratas — `CountRateRetroCov`, `CalculateNetRate_ACT`,
`CallCountPremiFacIn_act`, `InputRateFO` — seluruhnya RULE KEMBAR.** Rekomendasi §4 butir 2 yang lama
**gugur**: mereka tidak perlu diekspor ulang, mereka perlu **diport dua-duanya**.

---

## 4. Yang disarankan diekspor ulang lebih dulu

> ⛔ **Rekomendasi lama DIBATALKAN** — butir 1, 2 dan 3 seluruhnya menunjuk **rule kembar (G3)**:
>
> ~~1. `CountRateRetroCov` — sudah diekspor ulang, selesai.~~
> ~~2. Kelas a pada Activity — `CallCountPremiFacIn_act`, `CalculateNetRate_ACT`, `CountEdmAdjTSI_Act`, `InputRateFO`.~~
> ~~3. `IsFacRetro` (When) — predikat pemicu Fac Out.~~
>
> Ketiganya **bukan** kandidat ekspor ulang. `CountRateRetroCov` memang sudah diekspor ulang —
> **keduanya**, karena memang dua rule (K-060 amandemen). `CountEdmAdjTSI_Act` tetap G2 dan pindah ke
> daftar baru.

**Urutan yang berlaku — usulan, bukan keputusan:**

1. **`CountEdmAdjTSI_Act`** (G2, kelas a — uang, `Data-Aneka` di kedua sisi, Δ116). Satu-satunya
   kandidat kelas **a** yang benar-benar drift.
2. **Jalur produksi Fac Out** (G2, kelas c) — `InsertFacoutProduction`, `InsertFacoutProductionEDM`,
   `InsertFacoutProductionEDMCurr`, `SaveFacinProdFireNB_Act`, `SaveFacinProdEDMMarineCargo_Act`.
   Menyentuh F10/F11.
3. **Tiga RDBList produksi** (G2, kelas c) — `DeleteDataProduction`,
   `InsertHistoryAkseptasiPega_Sql`, `INSERTJSON_JSONPOLISMONITORING_FACIN`. Δ hanya 2 baris.
4. **`Protection_Act`** dan **`serviceInsertArasapasEDM_act`** (G2, kelas b — gerbang alur).
5. **20 pasangan `DecisionTable`** — `10-audit\07` mencatat tipe ini **100 % berbeda** (10 dari 10
   identitas). ⚠️ **Belum dipilah G2/G3**; lingkupnya tidak diketahui.
6. Sisanya (Section tampilan) — hampir seluruhnya **G3**; ditangani di berkas 10, bukan di sini.

⚠️ Butir 5 dan 6 berlabel **`[dugaan]`** — belum dibedah isinya.

---

## 5. Perintah audit

```powershell
# menghasilkan 526 / 87 dan sebaran kelasnya
# (fungsi Get-Lines/Get-Hash = kontrak 23 tag, lihat 07-hitung-berkas-dan-drift.md §5)
# untuk tiap pasangan versi-sama-isi-beda:
#   $d = Compare-Object $lines_A $lines_B
#   klasifikasi baris $d.InputObject menurut tabel §1
# CATATAN PowerShell: JANGAN memakai $B sebagai nama dictionary -
#   ia bertabrakan dengan $b pada foreach($b in ...) karena variabel tidak peka huruf.
#   Gejalanya: "[System.String] does not contain a method named 'ContainsKey'".
```

---

## 6. Yang TIDAK dapat ditentukan

1. `[pertanyaan terbuka]` **Mengapa** salinan Endorsement `CountRateRetroCov` memakai pembagi 10.000
   dan kehilangan rumus koreksi — versi lama, cabang terpisah, atau kekeliruan. Korpus tidak menjawab.
2. `[pertanyaan terbuka]` Apakah **27 pasangan** yang terpotong dari daftar memuat sesuatu yang lebih
   penting daripada 60 teratas. Pemotongan menurut prioritas kelas dan Δbaris, **bukan** menurut isi.
3. `[dugaan]` Perbedaan pada Section tampilan diduga tata letak — **belum dibedah**.
4. `[pertanyaan terbuka]` Apakah metode deteksi drift diganti ke hash 23 tag — **terbuka di K-058**.

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
