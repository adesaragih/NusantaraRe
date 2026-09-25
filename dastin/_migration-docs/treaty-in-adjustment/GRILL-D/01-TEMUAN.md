> Modul  : Treaty In Adjustment · Ronde D · 2026-09-24
> Sifat  : TAMBAH-SAJA. Setiap temuan bergolongan **EVIDENCED**, dengan ekspornya disebut.

# 01 · TEMUAN RONDE D

| # | Temuan | `SISI` | Menopang |
|---|---|---|---|
| **TD-01** | **Nol properti menyimpan nomor dokumen addendum eksternal**, atas 28 nama calon dan kelima bentuk penulis, di kedua ekspor | IRISAN | `GRL-19` |
| **TD-02** | **Materialitas mengunci 660 sel layar**, dan ia **sakelar DUA ARAH** — bukan "Non Material lebih sedikit" | 56-KHAS | `GRL-20` |
| **TD-03** | Kondisi penguncinya **18 bentuk berbeda**, bukan 220 sel seperti yang sempat dilaporkan gerbang 0 | 56-KHAS | rekonsiliasi angka |
| **TD-04** | `ContractRefNo` **muncul sekali**, di layar `OldData`, **bukan sebagai penulis** — satu-satunya yang mendekati nomor dokumen | IRISAN | `TD-01` |
| **TD-05** | **`TDA-17` diadili** — dan tiga dari tujuh kolomnya **tidak punya rumah di skema baru sama sekali**, korban `L-8` | IRISAN | `TDA-17` |

---

## TD-01 · Nol properti menyimpan nomor dokumen addendum

**Perkakas:** `tools/sapu-nomor-dokumen.py` · **EVIDENCED(kedua ekspor@ekspor-2026-09)**

**Yang disapu — daftar namanya dicatat di dalam perkakasnya, bukan di sini saja:** 28 nama calon —
`DocNo`, `DocumentNo`, `DocNumber`, `NoDoc`, `NoDokumen`, `NomorDokumen`, `AddendumNo`, `NoAddendum`,
`AddendumNumber`, `AddendumRef`, `AddendumID`, `EndorseNo`, `EndorsementNo`, `NoEndorse`, `RefNo`,
`NoRef`, `ReferenceNo`, `ContractRefNo`, `RefDoc`, `DocRef`, `SlipNo`, `NoSlip`, `LetterNo`,
`NoSurat`, `SuratNo`, `AmendmentNo`, `NoAmendment`, `AddNo`.

**Nama tag diperiksa lebih dulu** terhadap `../treaty-in/alat/datar-nama-elemen.csv`: **8 dari 8 ada**.
Sapuan yang memakai nama tag yang tidak ada di korpus akan melaporkan nol yang tidak berarti apa-apa.

**Kalibrasi — wajib sebelum klaim negatif:**

| Kasus positif yang sudah diketahui | Hasil |
|---|---|
| `EDMState` | **DITEMUKAN** — 20 kemunculan |
| `EDMMaterialType` | **DITEMUKAN** — 40 kemunculan |
| `OLDID` | **DITEMUKAN** — 23 kemunculan |

**LULUS.** Penyapunya bekerja, jadi nolnya bermakna.

**Cakupan:** 708 berkas disapu. **Ditolak:** 2 berkas bukan `.xml`. Salinan terbungkus
(`pyIncludedRuleXML`, `pyRuleVersionsList`) dibuang lebih dulu — ia bukan badan aturan.

### Hasil

> **NOL.** Tidak satu pun dari 28 nama muncul sebagai **penulis** di ekspor mana pun.
> Satu muncul sebagai **bukan penulis** — lihat `TD-04`.

**Menurut tabel putusan yang ditetapkan di muka:** nol di **kedua** ekspor berarti sistem lama
**tidak pernah merekam** nomor dokumen addendum eksternal. Ia hidup di kertas.

### BATAS KLAIM INI, dinyatakan bersama hasilnya

Dua bentuk berada di luar jangkauan sapuan ini, dan keduanya disebut supaya nol tidak terbaca lebih
kuat daripada yang dapat ia tanggung:

* **`Rule-Declare-Expression` dan `Rule-Declare-Trigger` tidak terekspor sama sekali** (`L-10`).
  Keduanya menulis nilai **tanpa dipanggil siapa pun**, dan tidak ada bentuk penulis yang melihatnya.
* **Langkah Java (32) dan langkah SQL/REST (191) tidak terurai** sebagai penugasan properti.

**Maka nol ini berarti: nol di antara bentuk yang terurai.** Ia tetap cukup untuk memutuskan — sebab
sebuah nomor dokumen yang dipakai orang harus muncul **di layar** untuk diketik dan **di penyimpan**
untuk disimpan, dan sapuan ini melihat keduanya.

---

## TD-02 · Materialitas adalah SAKELAR DUA ARAH — dan ini membalik pemahaman lama

**Perkakas:** sapuan `pyDisabledWhen` di badan `Section` kedua ekspor, salinan terbungkus dibuang.
**EVIDENCED(Section@ekspor-2026-09)** · **Ditolak:** 0 sel tanpa `pyValue` dalam jendelanya.

**660 dari 758 kondisi penguncian menyebut `EDMMaterialType`** — 87%. Penguncian layar di modul ini
**memang tentang materialitas**, bukan tentang hal lain yang kebetulan lewat.

Yang tidak diduga: **arahnya dua.**

| Pilihan | Sel terkunci | Nama properti berbeda | Seksi |
|---|---:|---:|---:|
| **Non Material** (`= 2`) | **640** | **93** | 19 |
| **Material** (`= 1`) | **20** | **5** | — |

### Apa yang dikunci saat NON MATERIAL — angka

`Currency` (34) · `Value` (22) · `RNMShare` (13) · `FacultativeShare` (10) ·
`FacultativeShareBrokerage` (10) · `Amount` (6) · `Deduction` (6) · `DeductionPct` (6) ·
`TreatyGroup` (8) · keempat batas bahaya beserta mata uangnya (`Earthquake`, `FloodJab`,
`FloodNation`, `RSMDLimit`, masing-masing 8 + 8) · `InitialDate`, `SubmissionDue`,
`ConfirmationDue`, `SettlementDue` (6 masing-masing) · `TypePortfolio` (6) · dan 70 nama lain.

### Apa yang dikunci saat MATERIAL — teks

| Properti | Sel |
|---|---:|
| `Exclusions` | 6 |
| `SpecialConditions` | 6 |
| `TreatyContractName` | 2 |
| `ContractRefNo` | 2 |
| `pyTemplateButton` | 4 |

> ### Maka aturannya bukan *"Non Material berarti lebih sedikit yang boleh diubah"*
>
> **Kedua jenis menyunting himpunan field yang SALING LEPAS:**
>
> * **Material** → boleh mengubah **angka**, tidak boleh mengubah **pengecualian, ketentuan khusus,
>   nama kontrak, dan nomor rujukan kontrak**.
> * **Non Material** → boleh mengubah **teks dan hal administratif**, tidak boleh mengubah
>   **angka mana pun**.
>
> Ini **lebih tajam** daripada rumusan *"Non Material berarti uang tidak berubah"*, dan ia terbaca
> dari ekspor — bukan dari keterangan.

**Batas yang jujur:** dari 640 sel Non Material, **158 `pyTemplateInputBox`** dan **54
`pyTemplateButton`** adalah **sel templat kisi**, bukan field bisnis bernama. Yang bernama:
**428**. Angka 640 tidak dipakai sebagai jumlah field.

---

## TD-03 · Rekonsiliasi angka dengan laporan gerbang 0

`GERBANG-0-TO-SPEC.md` menyebut **"220 kondisi `pyDisabledWhen` di badan 18 seksi, terbanyak
`DetailLimits` 60 dan `Layers` 38"**. Sapuan ronde ini memberi angka lain:

| | Gerbang 0 | Ronde D |
|---|---:|---:|
| kemunculan | 220 | **758** |
| **kondisi berbeda** | — | **18** |
| seksi | 18 | **26** |
| `DetailLimits` | 60 | 128 |
| `Layers` | 38 | 80 |

**Angka 18 cocok — tetapi bukan untuk hal yang sama:** gerbang 0 menyebut **18 seksi**, sapuan ini
menemukan **18 kondisi berbeda** di **26 seksi**. Dan tiap angka seksi gerbang 0 tepat **separuh**
angka ronde ini, yang konsisten dengan **satu ekspor saja** yang disapu, bukan keduanya.

> **Arah klaim gerbang 0 tidak berubah** — penguncian field memang nyata dan luas, dan itu yang
> menguatkan jawaban C. Yang dikoreksi **angkanya**, dan koreksinya dilaporkan di sini alih-alih
> angkanya dipakai diam-diam.

---

## TD-04 · `ContractRefNo` — satu-satunya yang mendekat, dan ia bukan penulis

| | |
|---|---|
| **Kemunculan** | **1**, di `Section/TreatyInNONProportionalOldData.xml` |
| **Sebagai penulis** | **0** — atas kelima bentuk |
| **`SISI`** | **IRISAN**, dibaca dari `pzInsKey` |
| **Sudah tercatat sebelumnya** | `../treaty-in/KEPUTUSAN-TANPA-VERIFIKASI.md` §4c — `pyReadOnly`, tidak ada di DDL |

Dan `TD-02` menambahkan satu sisi yang belum pernah terbaca: **`ContractRefNo` termasuk yang
dikunci saat Material** (2 sel). Sebuah field yang **tidak pernah ditulis siapa pun** tetap punya
aturan penguncian — mekanisme yang menjaga sesuatu yang tidak pernah ada isinya.

**Ia bukan nomor dokumen addendum.** Namanya menyebut **kontrak**, bukan addendum; ia muncul di
layar non-proporsional; dan nol penulis berarti ia tidak pernah terisi.

---

## TD-05 · `TDA-17` diadili — dan lubangnya lebih dalam dari tujuh kolom

**Ditambahkan 24 September 2026**, menutup calon `TDA-17` yang lahir di sesi induk **sesudah ronde
TDA ditutup** dan karena itu tidak pernah melewati adjudikasi.

### 1. Sisinya — **IRISAN**, dibaca dari `pzInsKey`

`RDBList/SaveTreatyInDetail.xml` ada di **kedua** ekspor dengan `pzInsKey` yang **sama persis**:

```
RULE-CONNECT-SQL ASM-FW-GISFW-INT-TREATYINDETAIL ASM!SAVETREATYINDETAIL #20260828T033641.5
```

Satu instans aturan, bukan dua. **IRISAN** — temuannya milik Treaty In, diadili di sini untuk jalur
addendum saja.

### 2. Adakah aturan HIDUP yang MEMBACA ketujuh kolom itu? — **TIDAK**

**Kalibrasi:** `TREATYID` — kolom yang pasti dirujuk SQL — ditemukan **21 rujukan**. **LULUS.**
**Cakupan:** 708 berkas. **Ditolak:** 2 bukan `.xml`.

| Kolom | Rujukan SQL | Di mana |
|---|---:|---|
| `DEDUCTIBLE`, `DEDUCTIBLE2`, `PREMIUM_EARNED`, `MDP_PCT`, `ROL_PCT`, `MDP` | **2** masing-masing | **hanya `SaveTreatyInDetail.xml`** — dan itu **penulis**, bukan pembaca |
| `ADJ_RATE` | **0** | namanya di parameter berbeda (`P_ADJUSMENT_RATE`) |

> **Tidak satu pun aturan membaca ketujuh kolom itu dari `TREATYINDETAILEDM`.** Ketujuhnya
> **ditulis dan tidak pernah dibaca** di sisi kontrak, dan **tidak ditulis sama sekali** di sisi
> addendum.
>
> Maka akibat langsungnya di dalam Pega: **nol**. Yang terdampak adalah **pembaca di luar sistem
> ini** — ADR-0051, dan itu tidak dapat diukur dari ekspor.

### 3. Nasibnya — **diperbaiki SEBAGIAN**, dan pembuktiannya memunculkan lubang ketiga

Dugaan awal *"diperbaiki, karena struktur baru memuat ketujuhnya"* **dibuktikan, dan ternyata
salah**:

| Kolom datar lama | Rumah di skema baru |
|---|---|
| `DEDUCTIBLE` | `LAYER.DEDUCTIBLE` |
| `MDP` | `LAYER.MDP` |
| `MDP_PCT` | `LAYER.PERSEN_MINIMUM_DEPOSIT` |
| `ADJ_RATE` | `LAYER.PERSEN_PENYESUAIAN` |
| **`DEDUCTIBLE2`** | **TIDAK ADA** |
| **`PREMIUM_EARNED`** | **TIDAK ADA** |
| **`ROL_PCT`** | **TIDAK ADA** |

**Empat diperbaiki. Tiga tidak punya rumah sama sekali.**

### 4. Kenapa ketiganya tidak punya rumah — dan sebabnya bukan kelalaian §10

Ketiganya dicari di `PETA-TELUSUR-JSON.md`. Hasilnya **tegas dan sama untuk ketiganya**: mereka
muncul **hanya di dalam pohon cermin**, tidak pernah di pohon utama.

```
ActualValue.LimitShareSummaryList[].Deductible2      ValueDifference.Limits[].Deductible2
ActualValue.Limits[].PremiumEarnedList[].Value       ValueDifference.Limits[].ROLPct
```

Dan dua di antaranya ada di **daftar titik buta** `datar-titik-buta-pohon.csv`, bertanda `TIDAK`:

```
ASM-FW-GISFW-Data-TreatyInLimits,PremiumEarned,TIDAK
ASM-FW-GISFW-Int-TREATYINDETAIL,DEDUCTIBLE2,TIDAK
```

> **Maka `TDA-17` bukan satu temuan melainkan dua, dan yang kedua lebih berat.**
>
> Yang pertama: tabel datar addendum tertinggal tujuh kolom.
> **Yang kedua: tiga dari tujuh itu tidak pernah terlihat pohon utama sama sekali** — mereka
> korban **`L-8`**, titik buta 414 properti. §10 tidak melewatkannya; §10 **tidak pernah
> diperlihatkan** kepadanya.
>
> Ini persis akibat yang `L-8` peringatkan: *"setiap kesimpulan yang pernah dihitung dari sumber itu
> dihitung atas semesta yang kurang."* Di sini semesta yang kurang itu **menghasilkan tiga kolom
> yang hilang dari model**, dan hilangnya **tidak akan terlihat** — sebab `KAMUS-KOLOM.md` tampak
> rapi.

### 5. Nasib akhir

| Bagian | Nasib | Dasar |
|---|---|---|
| tujuh kolom hilang di `TREATYINDETAILEDM` | **diperbaiki** | skema baru tidak punya dua proyeksi datar; ketimpangannya lenyap karena bentuknya (ADR-0056) |
| `DEDUCTIBLE`, `MDP`, `MDP_PCT`, `ADJ_RATE` | **diperbaiki** | rumahnya ada di `LAYER`, terbukti |
| **`DEDUCTIBLE2`, `PREMIUM_EARNED`, `ROL_PCT`** | **DITUNDA** | tidak punya rumah; sebabnya `L-8`. Pemiliknya **sesi to-spec induk**, dan ia **tidak dapat diputuskan di modul Adjustment** |

**`ROL_PCT` mungkin turunan** — Rate on Line adalah premi dibagi limit, dan §4 memang membuang
turunan. **Mungkin**, dan itu bukan dasar: ia tidak pernah diadili, sehingga tidak ada yang tahu
apakah ia disimpan atau dihitung. **Tidak ditebak.**
