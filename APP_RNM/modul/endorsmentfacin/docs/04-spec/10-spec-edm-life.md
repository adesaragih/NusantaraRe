# Bahan spec — Modul 6: Jalur Life Endorsement (EDM)

> **Ini BAHAN untuk `/to-spec`, bukan spec.** `/to-spec` ber-`disable-model-invocation: true` dan
> **belum dijalankan**. **Modul terakhir** rangkaian spec EDM.
>
> **Sumber:** E-1 §7 · E-3 (gerbang `IsLife`) · E-5 §2.1 · K-051 (`10-banding`, `11-banding`).
> Prior art bentuk: `04-spec\05`–`09-spec-edm`.
>
> **Keputusan mengikat:** K-005 · K-010/K-012 · K-018 · K-029 · **K-044 (Life IKUT, jalur terpisah)** ·
> K-046 · K-047 · K-048 · K-049 · K-050 · K-051.
>
> ⚠️ **Tiga temuan yang mengubah gambaran** — §1.3, §3.2, §3.3. Dua pertanyaan diajukan — §7.

---

## 0. Ringkas

| Pertanyaan | Jawaban |
| --- | --- |
| Lingkup Life | **16 inti EDM-only + ±45 subsistem medis/benefit EDM-only**; 60 berkas lain hanya **bercabang** ke Life (§1) |
| Life alur terpisah? | ✅ **Ya** — 4 titik cabang terverifikasi (§2) |
| Rumus premi Life sama? | ⚠️ **Keluarga sama, bentuk berbeda** — `TSILiability × RateLifeAverage × (1 + Rate/100) ÷ 1.000` (§3.1) |
| `CalculateScorLife_Act` = perhitungan premi? | ⛔ **BUKAN** — ia **skoring medis**, keluarannya keputusan, bukan uang (§3.2) |
| Perlu seam baru? | ⚠️ **Mungkin satu** — untuk skoring medis. **Diajukan** (§6.2) |

---

## 1. Lingkup Life — tiga angka, tiga dasar berbeda

### 1.1 Mengapa angkanya pernah 52, pernah 68

`[terverifikasi]` Tiga kriteria memberi tiga angka, dan **ketiganya sah** untuk pertanyaan berbeda:

| Kriteria | Jumlah | Menjawab |
| --- | ---: | --- |
| **A** — nama mengandung `Life` | **68** | "berapa berkas bernama Life" |
| **B** — nama luas (`life\|medic\|benefit\|disease\|plan\|scor\|lab`) **atau** isi bersinyal ≥ 25 | **156** | "berapa yang bernuansa Life/medis" |
| **C** — **menyebut predikat `IsLife`** | **60** | "berapa yang **bercabang** ke Life" |

⛔ **Angka 52 di E-4 bukan salah** — ia memakai ambang isi tanpa kata `Plan`. Angka **68** adalah
kriteria A murni. Keduanya mengukur hal berbeda.

📌 **Kriteria C adalah dasar paling kuat** — ia berbasis **rujukan nyata**, bukan nama. Tetapi C
**bukan** lingkup modul 6: sebagian besar 60 berkas itu **milik modul lain** yang kebetulan punya
cabang Life.

### 1.2 Lingkup yang diusulkan — berlapis

| Lapis | Jumlah | Milik |
| --- | ---: | --- |
| **L1 — Inti Life EDM-only** | **16** | **modul 6** |
| **L2 — Subsistem medis/benefit EDM-only** | **±45** | **modul 6** |
| **L3 — Titik cabang `IsLife`** | 60 | **modul 1–5** — modul 6 hanya menetapkan *perilaku cabangnya* |
| **L4 — Skoring risiko bersama NB** | ±20 | **bukan modul 6** — kemampuan NB yang dipakai ulang (§1.3) |
| L5 — Bernama Life tetapi **bersama NB** | 52 | dipakai ulang; hanya dicek saat integrasi (K-043 pilihan B) |

### 1.3 ⛔ Temuan: subsistem skoring sebagian besar **dibagi dengan NB**

`[terverifikasi]` Berkas skoring terbesar **ada di kedua korpus**:

| Berkas | Ukuran | Status |
| --- | ---: | --- |
| `Section\ScoringRisk` | **8,0 MB** | **bersama NB** |
| `Section\ScoringRiskForm1` … `Form4` | 1,0–2,5 MB | **bersama NB** |
| `Activity\SetScore_Act` | 412 KB | **bersama NB** |
| `Activity\ValueScoringRisk_act` | 267 KB | **bersama NB** |
| `Activity\SetDataScoringRisk_act` | 304 KB | **bersama NB** |
| `Activity\ScoringResult` | 144 KB | **bersama NB** |

⛔ **Skoring risiko BUKAN kemampuan baru Life-EDM.** Ia sudah ada di NB dan dipakai ulang.

Yang **EDM-only** adalah lapisan **pemeriksaan medis**:

| Berkas | Ukuran |
| --- | ---: |
| `Harness\Medical_Harnes` | **4,1 MB** |
| `Section\Medical_Sec` | **4,1 MB** |
| `Section\InputMedical2_Sec` | 3,1 MB |
| `Section\InputMedical1_Sec` | 1,0 MB |
| `Activity\CalculateScorLife_Act` | **1,9 MB** |
| `Activity\SetParamLab_Act` | 484 KB |
| `Activity\SaveMedical` · `CalculatePhysicalExam` · `AddSubHistoryDisease` | 21–56 KB |

📌 **Pemisahan yang benar:** `ScoringRisk` (bersama, risiko/okupasi) ≠ `CalculateScorLife_Act`
(EDM-only, **pemeriksaan medis**). Dua sistem skoring berbeda — jangan disatukan.

### 1.4 Enam belas inti Life EDM-only

`[terverifikasi]`

| Berkas | Ukuran | Peran |
| --- | ---: | --- |
| `Activity\CalculateScorLife_Act` | 1,9 MB | skoring medis — §3.2 |
| `Activity\SetTSIAllCoverageLife_ACT` | 303 KB | TSI seluruh coverage |
| `Activity\setRateLife_act` | 197 KB | penetapan rate Life |
| `Activity\ReCountPremiLifeEDM` | 175 KB | **premi Life** — §3.1 |
| `Activity\SetOLDValueToEDMWork_LIFE` | 75 KB | **before-image lapis C** — §4 |
| `Activity\AddCurrencyListLife_ACT` | 69 KB | daftar mata uang |
| `Activity\CountCoveragePropertyLife_Act` | 64 KB | cacah coverage |
| `Activity\AddAgeLife_ACT` | 32 KB | usia tertanggung |
| `Section\InputCoverageLife_FacIn` (+`_IsUW`) | 528 / 386 KB | layar coverage Life |
| `FlowAction\InputCoverageLife_FacIn` (+`_IsUW`) | 31 KB | aksi layar |
| `FlowAction\InputDtlCoverageLife_FacIn` (+`_IsUW`) | 30 KB | aksi detail |
| `RDBList\BrowseLifeRisk_SQL` | 6 KB | query risiko Life |
| `ReportDefinition\BrowseBenefitLife_RD` | 29 KB | daftar benefit |

---

## 2. Alur terpisah — empat titik cabang

`[terverifikasi]` `When\IsLife` · `pxObjClass` = `Rule-Obj-When` · `pyLogic` = `A OR B OR … OR P`
(**16 cabang**), seluruhnya `pyWorkPage.Quotation.BusinessOldId = "L1"` … `"L16"`.

⚠️ **Predikat Life membaca `BusinessOldId`**, bukan `BusinessType` — berbeda dari keenam predikat COB
yang membaca `OutData.pxResults(1).CARI2` (modul 4).

### 2.1 Titik cabang 1 — before-image lapis B: Life KELUAR

`[terverifikasi]` `Activity\SetOldData.xml` langkah **1**:

```
IF: IsLife    WhenTrue = 6 (keluar activity)   WhenFalse = 2 (lanjut)
```

⛔ **Life tidak mendapat lapis B sama sekali.** Nilai lama per baris tidak diisi lewat jalur umum.

### 2.2 Titik cabang 2 — penomoran endorsement Life

`[terverifikasi]` `Activity\SaveEDMToJsonPolicy_Act.xml`:

```
[12] RDB-List  « Generate No Endorsement »       RequestType = GenerateEndorsementNo
       IF IsLife   [T=3 F=2]   → berjalan saat IsLife SALAH (non-Life)
[13] RDB-List  « Generate No Endorsement LIFE »  RequestType = GenerateEDMNoLife
       IF IsLife   [T=2 F=3]   → berjalan saat IsLife BENAR
[15]   « Nopolis LIFE »                          RequestType = GetKodeProdLife_SQL
       IF …PolicyData.EndorsementNo==""
```

📌 Pasangan **komplementer** — gerbang sama, transisi terbalik (temuan E-5 §1.2).

⚠️ `GenerateEDMNoLife` **NIHIL di korpus**. Per keputusan work owner ini **bukan `panic`** —
penomoran nopolis Life tertanam di `SaveEDMToJsonPolicy_Act` sendiri (langkah 15–17 lewat
`GetKodeProdLife_SQL` + `GetSequenceNumber_SQL`).

### 2.3 Titik cabang 3 — jalur produksi Life

`[terverifikasi]` `SaveEDMToJsonPolicy_Act` langkah **28–29**:

```
[28] Call SaveFacinLive_Act        « -- Untuk LIfe »   (147 KB, Rule-Obj-Activity)
[29] Call SaveFacinSpreadLife_Sql  « -- Untuk LIfe »   (114 KB, Rule-Obj-Activity)
```

⛔ **`SaveFacinProdAllEDM_Act` tidak punya cabang Life** — ketujuh cabangnya `IsFire`, `IsAneka`,
`IsBonding`, `isGolfInsurance`, `IsMarineCargo`, `IsMBU`, `IsPA`. Life berjalan **sejajar**, bukan
di dalamnya.

⚠️ `SaveFacinSpreadLife_Sql` ber-`pxObjClass` **`Rule-Obj-Activity`**, bukan rule SQL — namanya
menyesatkan. Jebakan folder/nama lagi.

### 2.4 Titik cabang 4 — before-image lapis C

`[terverifikasi]` `SetValueToEDMWork` langkah **14.11**: `Call SetOLDValueToEDMWork_LIFE` — §4.

---

## 3. Perhitungan premi Life

### 3.1 `ReCountPremiLifeEDM` — keluarga sama, bentuk berbeda

`[terverifikasi]` `pxObjClass` = `Rule-Obj-Activity` · `pyClassName` = **`Data-Party-Person`** ⚠️
(kelas **orang**, bukan coverage — berbeda dari seluruh kalkulator premi lini lain).

```
L929   .TSI          = @if(.RIRisk>0, @divide(Local.TSI × .RIRisk, 1000, 20), Local.TSI)
L975   .TSILiability = @if(.RIRisk>0, @If(.TSI − Param.TSICeding>0, .TSI − Param.TSICeding, 0), Param.TSILiability)
L1015  .Premium      = @Math.divide((.TSILiability × .RateLifeAverage × (1 + (.Rate/100))), 1000, 20) − 0
L1261  .PremiumSpreaded = @if(.TreatyType=="1000036", 0, Local.Premi × .SharePercentage/100)
L1307  .TSISpreaded     = @if(.TreatyType=="1000036", 0, Local.TSICov × .SharePercentage/100)
```

**Banding terhadap rumus dasar:**

| | Rumus dasar (FIRE/PA) | **Life** |
| --- | --- | --- |
| Basis | `.TSI` | **`.TSILiability`** (TSI dikurangi bagian ceding) |
| Rate | `.Rate` | **`.RateLifeAverage`** — properti rate **berbeda** |
| Peran `.Rate` | rate premi | ⚠️ **faktor pembebanan** `(1 + .Rate/100)` |
| Pembagi | 10⁹ (FIRE) / 100.000 (PA) | **1.000** |
| Skala | 4 | **20** |

⛔ **Keluarga sama** (`TSI × rate ÷ pembagi-komposit`), **bentuk berbeda**: basis liability, rate
dari properti lain, dan **satu faktor pembebanan tambahan**.

✅ Pembagi **1.000** konsisten **K-018** — hanya rate (‰) yang berkontribusi; tidak ada
`ProRatePercent`, `FirstLossScale`, maupun `IndemnityPercentage`.

⚠️ **Kejanggalan K-046:** `.TreatyType=="1000036"` → spreading dipaksa **0**. Kode treaty tertanam
literal, artinya tidak dijelaskan korpus. **Diport apa adanya.**

### 3.2 ⛔ `CalculateScorLife_Act` BUKAN perhitungan premi

`[terverifikasi]` Profil token berkas 1,9 MB itu:

| Token | Cacah |
| --- | ---: |
| `.Rate` · `.TSI` · `.Premium` · `@Math.divide` | **0 · 0 · 0 · 0** |
| `Medical` | **6.975** |
| `Age` | **3.145** |
| `Scor` | 162 |

Keluarannya **keputusan**, bukan uang:

```
L710   .Medical.Medical_Examination.Score.AFP        = @if(AFP  <= 8.78,  "Standard", "Postpone for 12 Months")
L901   .Medical.Medical_Examination.Score.CEA        = @if(CEA  <= 5.00,  "Standard", "Postpone for 12 Months")
L1474  .Medical.Medical_Examination.Score.PSA        = @if(PSA  <= 4.00,  "Standard", "Postpone for 12 Months")
L2193  .Medical.Medical_Examination.Score.HbA1c_NGSP = @if(HbA1c<=5.7, <bergantung .Age>, …)
L3301  .Medical.Medical_Examination.Score.eGFR       = @if(eGFR 20–49, "Decline", …)
```

⛔ **Ini underwriting medis, bukan perhitungan premi.** Domainnya berbeda: masukan hasil
laboratorium, keluaran **`"Standard"` / `"Postpone for 12 Months"` / `"Decline"`**.

⚠️ **Ambang klinis diport apa adanya.** Nilainya (8,78 · 5,00 · 4,00 · 5,7 · …) adalah **aturan
medis/bisnis**, bukan konstanta teknis. **Jangan ditebak, jangan dibulatkan, jangan diseragamkan.**

📌 Beberapa skor bergantung **`.Age`** (mis. `HbA1c_NGSP` bercabang pada rentang usia 17–20, 21–39,
…) — konsisten dengan kelas `Data-Party-Person` dan dengan adanya `AddAgeLife_ACT`.

⛔ **`"Decline"` sebagai keluaran skoring** menyentuh keputusan penerimaan. Kaitannya dengan tangga
akseptasi §5 `[pertanyaan terbuka]`.

### 3.3 Dampak ke Seam 3 — **bukan** bentuk kelima

Setelah dibedah, Life **tidak** menambah bentuk perhitungan premi baru:

| Bentuk (dari `11-banding` §5.2) | Life |
| --- | --- |
| 1 — rumus dasar `TSI × rate ÷ pembagi-komposit` | ✅ **Life masuk di sini**, dengan basis `TSILiability`, rate `RateLifeAverage`, faktor `(1+Rate/100)`, pembagi 1.000 |
| 2 — dua-bagian endorsement | belum terlihat di jalur Life |
| 3 — tabel tarif (Travel) | tidak |
| 4 — dekomposisi delta coverage | tidak |

✅ **Seam 3 tetap empat bentuk.** Life adalah **parameterisasi bentuk 1**, bukan bentuk baru.

⛔ **`CalculateScorLife_Act` sama sekali di luar Seam 3** — ia bukan perhitungan uang.

---

## 4. Before-image Life

### 4.1 Lapis B — Life DILEWATI

`[terverifikasi]` §2.1. Konsekuensi: properti `*Old` per baris **tidak terisi** untuk Life.

⚠️ `[pertanyaan terbuka]` bagaimana layar endorsement Life menampilkan "dari → menjadi" bila lapis B
kosong. Belum ditelusuri.

### 4.2 Lapis C — dua penanda, bukan satu

`[terverifikasi]` `Activity\SetOLDValueToEDMWork_LIFE.xml` (75 KB, terkecil dari tujuh varian):

```
[1]     IF: IsLife  (aktif=true)
  [1.1] Property-Set   « Copy PersonList LIFE »     IF: IsLife (aktif=false)
        SET newWorkPage.OfferFacIn.PersonList = newWorkPage.OfferFacIn.OldData.PersonList
  [1.2]                « Set flag old data LIFE »   IF: IsLife (aktif=false)
    [1.2.1] Property-Set   SET .FlagOldData = "old"     ← penanda TAMBAHAN
    [1.2.2]
      [1.2.2.1] Property-Set   SET .IsOldData  = "old"
```

| | Life | Enam varian lain |
| --- | --- | --- |
| Cacah `IsOldData` | **4** | 8–14 |
| Penanda yang disetel | **`.FlagOldData` DAN `.IsOldData`** | hanya `.IsOldData` |
| `IsProRate` | 0 | hanya `_FIRE` yang menyetel (2×) |
| Kedalaman list | **`PersonList` saja** | 2–4 tingkat bersarang |

⛔ **Cacah 4 bukan kelalaian** — model data Life hanya punya **satu tingkat list**. Jumlah penanda
mengikuti kedalaman, bukan kelengkapan implementasi.

⚠️ `[pertanyaan terbuka]` apakah `.FlagOldData` dibaca di tempat lain — hanya Life yang menyetelnya.

---

## 5. Tangga akseptasi — Life melewatinya

`[terverifikasi]` Korpus EDM **tidak memuat** `GetLimitAkseptasiLife_Act`. Arsip §4.5 mencatat
`Decision19 --When IsLife--> Decision39` — Life melompati seluruh tangga.

⛔ **`[pertanyaan terbuka]` — apakah endorsement Life memang tanpa persetujuan?** Milik **work owner
+ Underwriting**. **Tidak saya putuskan.**

⚠️ Sampai dijawab, perilakunya **diport apa adanya**: Life melewati tangga.

📌 Pertanyaan ini bersinggungan dengan §3.2: `CalculateScorLife_Act` dapat menghasilkan
**`"Decline"`**. Bila Life melewati tangga akseptasi, **apa yang terjadi pada kasus ber-skor
Decline** belum terjawab dari korpus.

⚠️ Klaim arsip §4.3 (*"nilai dasar akseptasi = selisih TSI"*) tetap **`[belum diuji]`** dan tidak
dipakai di sini.

---

## 6. Kontrak modul dan seam

### 6.1 Life memakai seam yang sudah ada

| Kebutuhan Life | Seam |
| --- | --- |
| Predikat `IsLife` + predikat lain | **Seam 1** `rules.Eval` |
| Premi Life (`ReCountPremiLifeEDM`) | **Seam 3** `premium.Calculate` — bentuk 1 terparameterisasi (§3.3) |
| Before-image (lapis A + C; lapis B dilewati) | **Seam 4** `PrepareBeforeImage` — cabang Life |
| Alur masuk (kasus Life lahir sama seperti lini lain) | **Seam 5** `OpenCase` |
| Tangga akseptasi | ⛔ **dilewati** (§5) |

✅ **Empat dari lima seam sudah cukup.** Tidak ada yang perlu ditambah untuk premi maupun
before-image Life.

### 6.2 ⚠️ Satu kandidat seam — skoring medis

`CalculateScorLife_Act` **tidak masuk seam mana pun**:

- bukan predikat (Seam 1) — keluarannya string keputusan, bukan boolean;
- bukan perhitungan uang (Seam 3) — nol `.Rate`/`.TSI`/`.Premium`;
- bukan before-image (Seam 4), bukan alur masuk (Seam 5);
- bukan tangga akseptasi (Seam 2) — Life melewatinya.

**Yang tidak dapat diuji tanpa seam:** ambang klinis menghasilkan skor yang benar, skor
bergantung-usia bercabang benar, dan `"Decline"` muncul pada kondisi yang benar.

| | **Pilihan A — Seam 6 baru** | **Pilihan B — tanpa seam** |
| --- | --- | --- |
| Bentuk | `services/underwriting.ScoreMedical(hasilLab, usia) → SkorMedis` | uji lewat repository tiruan / integrasi |
| Untung | ambang klinis teruji langsung; keluaran `Decline` teruji sebagai perilaku | tetap 5 seam |
| Rugi | seam keenam | ambang klinis tidak teruji — **risiko tinggi**, ini aturan medis |

**Usulan saya: Pilihan A.** Alasannya: ambang klinis adalah **aturan bisnis/medis** dengan
konsekuensi nyata, dan `"Decline"` adalah **keluaran keputusan**. Tidak mengujinya berarti salah
ambang baru ketahuan di produksi. Bentuknya sejajar Seam 4 dan 5: satu pintu untuk satu subsistem.

⛔ **Tidak saya putuskan.** Diajukan ke work owner.

---

## 7. Pertanyaan terbuka

| # | Pertanyaan | Pemilik | Memblokir? |
| --- | --- | --- | --- |
| **1** | **Apakah endorsement Life memang tanpa persetujuan** (§5) | **work owner + UW** | ⚠️ ya, untuk modul akseptasi |
| **2** | **Seam 6 untuk skoring medis?** (§6.2) | **work owner** | ⚠️ ya, menentukan bentuk pengujian |
| **3** | **Apa yang terjadi pada kasus ber-skor `"Decline"`** bila Life melewati tangga (§5) | work owner + UW | tidak |
| **4** | **Bagaimana layar Life menampilkan "dari → menjadi"** bila lapis B kosong (§4.1) | work owner | tidak |
| **5** | **Apakah `.FlagOldData` dibaca di tempat lain** — hanya Life menyetelnya (§4.2) | — verifikasi lanjutan | tidak |
| **6** | **Arti kode treaty `"1000036"`** yang memaksa spreading nol (§3.1) | work owner | tidak |

---

## 8. Kejanggalan — diport apa adanya (K-046)

| # | Butir | Nama test |
| ---: | --- | --- |
| 1 | **Ambang klinis tertanam** di `CalculateScorLife_Act` — aturan medis, bukan konstanta teknis | `K046_SkorMedis_AmbangKlinisDipertahankan` |
| 2 | **`.TreatyType=="1000036"` → spreading 0** (§3.1) | `K046_Life_TreatyType1000036_SpreadingNol` |
| 3 | **`.Rate` sebagai faktor pembebanan** `(1 + Rate/100)`, bukan rate premi (§3.1) | `K046_Life_RateSebagaiPembebanan` |
| 4 | **`_LIFE` menyetel dua penanda** (`.FlagOldData` + `.IsOldData`) (§4.2) | `K046_Life_DuaPenandaOldData` |
| 5 | **`SaveFacinSpreadLife_Sql` ber-`pxObjClass` `Rule-Obj-Activity`** — nama menyesatkan | — catatan penamaan |
| 6 | **`GenerateEDMNoLife` nihil**, penomoran tertanam di pemanggil — **bukan `panic`** | — keputusan work owner |

---

## 9. Out of Scope

1. **Skoring risiko `ScoringRisk`** dan keluarganya — **bersama NB** (§1.3), bukan modul 6.
2. **Jalur produksi Life** (`SaveFacinLive_Act`, `SaveFacinSpreadLife_Sql`) — belum ditelusuri
   isinya; menunggu tabel flat + `ALL_SOURCE` bersama jalur produksi lain.
3. **Tangga akseptasi** — Life melewatinya; klaim §4.3 tetap `[belum diuji]`.
4. **52 berkas bernama Life yang bersama NB** — dipakai ulang, dicek saat integrasi (K-043 pilihan B).
5. **Tabel benefit/plan** (`D_BenefitList`, `BrowseBenefit*`, `D_PkgPlan`, …) — lapisan repository.

---

## 10. Kesiapan

✅ **Bahan lengkap untuk `/to-spec`:**

| Butir | Status |
| --- | --- |
| **Lingkup Life berlapis** L1–L5, dengan penjelasan selisih 52 / 68 / 156 / 60 | §1 |
| ⛔ **Temuan:** skoring risiko **bersama NB**, hanya pemeriksaan medis yang EDM-only | §1.3 |
| 16 inti Life EDM-only | §1.4 |
| **Empat titik cabang** terverifikasi dengan kutipan langkah | §2 |
| Rumus premi Life + banding terhadap rumus dasar | §3.1 |
| ⛔ **`CalculateScorLife_Act` = skoring medis, bukan premi** | §3.2 |
| ✅ **Seam 3 tetap 4 bentuk** — Life parameterisasi bentuk 1 | §3.3 |
| Before-image Life — lapis B dilewati, lapis C dua penanda | §4 |
| Tangga akseptasi — `[pertanyaan terbuka]` | §5 |
| Kontrak seam + **usulan Seam 6 diajukan** | §6 |
| 6 kejanggalan K-046 bernama test | §8 |

⛔ **`/to-spec` dan `/to-tickets` TIDAK dijalankan.**
⚠️ **Pertanyaan 1 dan 2 (§7) sebaiknya dijawab sebelum spec ditulis.**

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan, tanpa data medis pasien.*
