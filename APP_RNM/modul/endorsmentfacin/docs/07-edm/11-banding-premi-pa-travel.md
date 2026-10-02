# Banding lanjutan — rumus premi PA & Travel + pengisi prorata

> **Putaran verifikasi, bukan spec.** Melanjutkan `10-banding-premi-edm-nb.md` yang menyisakan dua
> hal sebelum kontrak Seam 3 dikunci: **(1)** rumus PA & Travel, **(2)** pengisi
> `ProratePercentEDMEnd`/`ProratePercentStartEDM`.
>
> **Sumber:** `Endorsment Fac In\` dan `NB FacIn\`. READ-ONLY. `RNM_BRD\` tidak dibaca (K-005).
> Life = **modul 6**, tidak dibahas di sini.

---

## 0. Jawaban singkat

| Pertanyaan | Jawaban |
| --- | --- |
| Rumus **PA** sama dengan NB? | ✅ **SAMA** — berkas bersama `CalculatePremiPA_FacIn` **identik fungsional**. EDM **menambah** jalur pita tarif |
| Rumus **Travel** sama? | ⛔ **BEDA BENTUK** — bukan `Rate × TSI` sama sekali; premi dari **tabel tarif**, dan tabelnya **di luar korpus** |
| `ProratePercent*` = `Prorate*`? | ⛔ **BUKAN** — penyebut, satuan, dan skala **semuanya berbeda**. **Dua pasang rasio, bukan satu** |
| Seam 3 cukup "konteks endorsement opsional"? | ⚠️ **Tidak cukup** — perlu **tiga bentuk perhitungan**, bukan satu bentuk dengan varian |

⛔ **Satu temuan tambahan yang besar** — §4: `SetLocalNonMbuProrate` ternyata bukan sekadar pengisi
prorata, melainkan **kalkulator delta premi tingkat coverage** — bentuk keempat.

---

## 1. PA — rumus dasar SAMA, EDM menambah jalur pita tarif

### 1.1 Berkas bersama identik fungsional

`[terverifikasi]` `CalculatePremiPA_FacIn` ada di **kedua korpus** dan **IDENTIK FUNGSIONAL** di
bawah kontrak 23-tag + `pzIndexes` (K-043).

Rumusnya (`NB FacIn\Activity\CalculatePremiPA_FacIn.xml`):

```
L712   .Premium = (@Math.divide((.TSI × .Rate), 100000, 4) × pyWorkPage.OfferFacIn.ProRatePercent) − .Discount
L857   .Premium = (@Math.divide((.TSI × .Rate), 1000, 4) × @Math.divide(.PctShortPeriod, 100, 4)) − .Discount
L1002  .Premium = (@Math.divide((.TSI × .Rate), 1000, 4) × 1) − .Discount
L1145  .Premium = @Math.divide((.TSI × ProRatePercent × .Rate), 100000, 20) − .Discount
```

✅ **Ini rumus dasar yang sama**, dengan **pembagi komposit K-018 yang berbeda karena faktornya
berbeda**:

| Lini | Faktor | Pembagi |
| --- | --- | ---: |
| **FIRE** | `.Rate` (‰) × `ProRatePercent` × `FirstLossScale` × `IndemnityPercentage` | 1.000 × 100 × 100 × 100 = **10⁹** |
| **PA** | `.Rate` (‰) × `ProRatePercent` | 1.000 × 100 = **100.000** |
| **PA** *(varian periode pendek)* | `.Rate` (‰), `PctShortPeriod` dibagi terpisah | **1.000** |

⛔ **Pembagi berbeda BUKAN rumus berbeda.** Ia konsekuensi langsung aturan **K-018**: pembagi =
hasil kali kontribusi satuan per faktor yang ikut. PA tidak punya `FirstLossScale` dan
`IndemnityPercentage`, jadi pembaginya 10.000× lebih kecil.

📌 Ini justru **penguat K-018**: aturannya terbukti berlaku lintas lini, bukan angka hafalan per lini.

### 1.2 Yang EDM tambahkan — jalur pita tarif

`[terverifikasi]` `Endorsment Fac In\Activity\calculatePremiPA.xml` (EDM-only, 234 KB) **tidak**
memakai rumus `TSI × Rate ÷ 10⁹`. Ia memakai **pencarian pita tarif** lalu akumulasi:

```
L654   Param.MAX_TSI          = .TSI
L675   Param.MIN_TSI          = .TSI
L1047  Primary.Rate           = .PCT_RATE_BAWAH
L1093  Primary.TSI            = .MIN_TSI
L1113  Primary.Premium        = (.PCT_RATE_BAWAH × .MIN_TSI)      ← premi per pita
L1258  Primary.Rate           = .PCT_MIN_RATE
L1499  Param.totalPremiABD    = Param.totalPremiABD + .Premium    ← akumulasi
```

Lalu premi akhir dihitung dari akumulasi itu:

```
L2034  .Premium = @Math.divide((Param.totalPremiABD × .Rate), 1000, 4) × .ProRatePercent
L3151  .Premium = @Math.divide((Param.totalPremiABD × .Rate × local.ProRatePercent), 100, 4)
L3293  .Premium = @Math.divide((Param.totalPremiABD × .Rate × @Math.divide(.PctShortPeriod, 100, 4)), 100, 4)
L3435  .Premium = @Math.divide((Param.totalPremiABD × .Rate), 100, 4) × 1
```

⛔ **Dua tahap:** premi per orang/pita diakumulasi ke `totalPremiABD`, **baru** dikenai rate dan
prorata. Basisnya **bukan TSI**, melainkan **akumulasi premi**.

⚠️ **`[pertanyaan terbuka]`** — `PCT_RATE_BAWAH`, `PCT_MIN_RATE`, `MIN_TSI`, `MAX_TSI` adalah kolom
**pita tarif**; sumber tabelnya **belum ditelusuri**. Jangan ditebak.

⚠️ **Kejanggalan K-046:** `L3059 .Rate = 20` — **rate tertanam literal**. Diport apa adanya;
kandidat perbaikan.

### 1.3 Kesimpulan PA

| | Status |
| --- | --- |
| Rumus dasar PA | ✅ **SAMA** dengan NB — berkas bersama identik fungsional |
| Pembagi | **100.000** (bukan 10⁹) — konsekuensi K-018, bukan rumus lain |
| Jalur pita tarif EDM | **TAMBAHAN**, bentuk berbeda (akumulasi, bukan TSI × Rate) |

---

## 2. Travel — BEDA BENTUK, tarif di luar korpus

### 2.1 Tidak ada `Rate` sama sekali

`[terverifikasi]` `Endorsment Fac In\Activity\CalculatePremiumTravel.xml`:
**`.Rate` = 0 kemunculan · pembagi `1000000000` = 0 kemunculan.**

Premi dibaca **langsung** dari sumber lookup:

```
L809   Primary.Premium              = .PREMIUM                ← premi dari tabel tarif
L901   local.PremiLebihPerminggu    = .PREMIUMMOREWEEK        ← tambahan per minggu
L1005  Primary.Premium              = Primary.Premium + (local.PremiLebihPerminggu × @round((local.SelisihHari − local.MaxHari)/7))
L1358  Primary.Premium              = Primary.Premium + (local.PremiLebihPerminggu × local.SisaMinggu)
L3110  Primary.Premium              = .Premium
```

Plus penyesuaian TSI:

```
L1577  primary.TSI = primary.TSI / 2
L1745  primary.TSI = (primary.TSI × 25) / 100
```

⛔ **Bentuk premi Travel: `premi_tabel + (premi_per_minggu × jumlah_minggu_lebih)`.**
Tidak ada rate, tidak ada TSI × Rate, tidak ada pembagi komposit.

### 2.2 Sumber tarif — di luar korpus

`[terverifikasi]` Rantainya:

```
CalculatePremiumTravel  langkah 2 → Call fillActPremiTravel
  fillActPremiTravel    langkah 9 → RDB-List
                                      RequestType = ViewPremiTravel_SQL
                                      BrowsePage  = OutputData
```

⛔ **`ViewPremiTravel_SQL` NIHIL di korpus EDM.** Ia termasuk **61 rujukan menggantung kategori D**
yang sudah tercatat di E-6 §4.2.

⚠️ **`[pertanyaan terbuka]` — struktur tabel tarif Travel tidak dapat diketahui dari korpus.**
Jangan ditebak. Per `CLAUDE.md` §4.5, rujukan ke rule yang tidak ada di folder mana pun →
**`panic`**, bukan nilai default.

📌 `PREMIUMMOREWEEK` muncul **hanya di satu berkas** (`CalculatePremiumTravel`, 2×) — konsisten
dengan kolom hasil query, bukan properti kasus.

### 2.3 Kesimpulan Travel

⛔ **Travel adalah bentuk ketiga.** Rumus dasar `TSI × Rate ÷ pembagi-komposit` **tidak berlaku**.
Ia **pencarian tabel + tambahan per minggu**. Menyamakannya dengan FIRE/PA akan salah total.

---

## 3. Prorata — DUA pasang rasio, bukan satu

### 3.1 Pengisi `ProratePercent*`

`[terverifikasi]` `Endorsment Fac In\Activity\SetLocalNonMbuProrate.xml`
(`pxObjClass` = `Rule-Obj-Activity`, `pyClassName` = `ASM-FW-GISFW-Data-Coverage`, 297 KB):

```
L6168  .ProratePercentStartEDM = @Math.divide(local.dateDifferentBegEDM, 365, 6) × 100
L6188  .ProratePercentEDMEnd   = @Math.divide(local.dateDifferentEDMEnd,  365, 6) × 100
L6208  .ProRatePercent         = .ProratePercentStartEDM + .ProratePercentEDMEnd
```

### 3.2 ⛔ Banding terhadap `Prorate*` — berbeda di TIGA hal

| | **`ProratePercentStartEDM` / `ProratePercentEDMEnd`** | **`ProrateStartEDM` / `ProrateEDMEnd`** |
| --- | --- | --- |
| **Diisi oleh** | `SetLocalNonMbuProrate` L6168/L6188 | `SetValueToEDMWork` L15 · `CountPremiEDM_DT` (K-048 §8.1) |
| **Penyebut** | **365** — tahun tetap | **periode polis aktual** (`EndDateTime − StartDateTime`) |
| **Satuan** | **persen** (`× 100`) | **pecahan** (tanpa `× 100`) |
| **Skala** | **6** desimal | **20** desimal |
| **Dipakai di** | premi **dasar** per coverage | delta **spreading** per baris produksi |

⛔ **Jawaban tegas: BUKAN nilai yang sama di properti berbeda.** Keduanya **hitungan berbeda** dengan
**satuan berbeda**. Menyatukannya akan menggeser premi **faktor 100** (persen vs pecahan) dan
mengubah basis periode (365 vs periode polis).

📌 `L6208 .ProRatePercent = ProratePercentStartEDM + ProratePercentEDMEnd` menjelaskan varian
dua-bagian di `10-banding` §4.1: untuk endorsement, `ProRatePercent` **adalah jumlah kedua porsi**.

### 3.3 Pemanggil

`[terverifikasi]` `SetLocalNonMbuProrate` dipanggil dari **4 berkas**:
`CalculatePremiFire` · `calculatePremiPA` · `HitungPremiOnChange` · `inputCoveragePA_preAct`.

⛔ Seluruhnya jalur **premi dasar**. Tidak satu pun jalur produksi. Memperkuat pemisahan §5
`10-banding`.

### 3.4 ⚠️ Pertanyaan satuan yang belum terjawab

`[terverifikasi]` Di rumus dasar FIRE, `ProRatePercent` berada **di dalam** pembagi komposit — faktor
100-nya sudah diperhitungkan. Tetapi di `SetLocalNonMbuProrate` L4304 ia dipakai sebagai **pengali
langsung** dengan pembagi hanya **1.000**:

```
local.PremiHitung = … + ((.TSI − .OldCoverage(1).TSI) × @Math.divide(.OldCoverage(1).Rate, 1000, 4)) × .ProRatePercent
```

⚠️ Bila `.ProRatePercent` bernilai persen (mis. `50` untuk setengah tahun), hasilnya **100× lebih
besar** daripada bila ia pecahan. Korpus tidak menjelaskan mana yang dimaksud.

⛔ **`[pertanyaan terbuka]` — milik work owner + Aktuaria.** Salah membaca satuan di sini menggeser
premi faktor 100. **Jangan ditebak**; diport apa adanya sampai dijawab.

---

## 4. ⛔ Temuan tambahan — `SetLocalNonMbuProrate` adalah kalkulator delta premi

`[terverifikasi]` Berkas ini **bukan sekadar pengisi prorata**. Ia memuat **bentuk keempat**
perhitungan premi — delta tingkat **coverage** (bukan baris spreading):

```
L3316  .Premium = (.Rate − .OldCoverage(1).Rate) × (local.selisihHari / local.OldPeriodeDifferent) × .OldCoverage(1).TSI
                + ((.TSI − .OldCoverage(1).TSI) × .Rate) × .ProRatePercent

L3524  .Premium = ((.TSI − .OldCoverage(1).TSI) × .Rate) × .ProRatePercent
                + .OldCoverage(1).TSI × (.Rate − .OldCoverage(1).Rate) × .ProRatePercent

L3882  .Premium = ((.TSI − .OldCoverage(1).TSI) × .OldCoverage(1).Rate) × .ProRatePercent
L4067  .Premium = .OldCoverage(1).TSI × (.Rate − .OldCoverage(1).Rate) × .ProRatePercent
```

⛔ **Pola: premi tambahan = (selisih rate × TSI lama) + (selisih TSI × rate).** Ini **dekomposisi
delta**, berbeda dari:

- **bentuk 1** — rumus dasar `TSI × Rate ÷ pembagi`;
- **bentuk 2** — dua-bagian `(baru × EDMEnd) + (lama × StartEDM)` (`10-banding` §4.1);
- **bentuk 3** — tabel tarif Travel;
- **bentuk 4** — dekomposisi delta di sini.

Ditambah cabang `.CalculateMethod=="2"` yang mengganti `ProRatePercent` dengan
`@Math.divide(.PctShortPeriod, 100, 4)` — **periode pendek**, di L4304/L4559/L4820/L5270.

⚠️ **`[pertanyaan terbuka]`** arti `.CalculateMethod` (nilai `"2"` vs lainnya) — belum terjawab dari
korpus; keluarga enumerasi yang sama dengan K-029.

---

## 5. Dampak final ke kontrak Seam 3

### 5.1 Kesimpulan per lini

| Lini | Rumus dasar | Bentuk tambahan EDM | Status |
| --- | --- | --- | --- |
| **FIRE** | ✅ SAMA (`÷10⁹`) | dua-bagian (bentuk 2) + delta coverage (bentuk 4) | ✅ terverifikasi |
| **PA** | ✅ SAMA (`÷100.000`) — berkas bersama **identik fungsional** | pita tarif (akumulasi `totalPremiABD`) | ✅ terverifikasi |
| **TRAVEL** | ⛔ **TIDAK BERLAKU** | tabel tarif + tambahan per minggu (bentuk 3) | ⚠️ tarif **di luar korpus** |
| Aneka · Golf · MBU · Marine Cargo | ✅ memakai activity **bersama** dengan NB | — | tidak ada activity premi EDM-only |

### 5.2 ⚠️ Kontrak "konteks endorsement opsional" TIDAK CUKUP

Kesimpulan `10-banding` §7.2 mengusulkan Seam 3 diperluas dengan **satu** konteks endorsement
(sepasang rasio + `OldCoverage`). **Setelah PA dan Travel dibedah, itu tidak cukup.**

Yang sebenarnya dibutuhkan:

| # | Bentuk | Masukan khas |
| ---: | --- | --- |
| **1** | Rumus dasar — `TSI × Rate × faktor ÷ pembagi-komposit` | faktor per lini; pembagi diturunkan resolver **K-018** |
| **2** | Dua-bagian endorsement — `(baru × EDMEnd) + (lama × StartEDM)` | sepasang **`ProratePercent*`** (persen, /365) + `OldCoverage(1).*` |
| **3** | Tabel tarif (Travel) | hasil lookup `.PREMIUM` + `.PREMIUMMOREWEEK` × minggu ⚠️ **tabel di luar korpus** |
| **4** | Dekomposisi delta coverage | `(Δrate × TSI_lama) + (ΔTSI × rate)`, dengan cabang `.CalculateMethod` |

⛔ **Tetap TIDAK ada seam baru.** Total seam tetap **5**. Yang berubah: `premium.Calculate` bukan
"satu rumus + varian opsional", melainkan **satu pintu masuk yang memilih di antara empat bentuk
perhitungan** berdasarkan lini bisnis dan konteks endorsement.

### 5.3 Bentuk perluasan yang diusulkan

```
premium.Calculate(konteks) → Money
   konteks memuat:
     · lini bisnis (dari Seam 1, predikat COB)
     · faktor-faktor coverage (TSI, Rate, ProRatePercent, FirstLossScale, IndemnityPercentage)
     · [opsional] konteks endorsement: OldCoverage + sepasang ProratePercent* (persen, /365)
     · [opsional] hasil lookup tarif (Travel) — dari repository, bukan dihitung
     · [opsional] CalculateMethod (periode pendek)
```

⚠️ **Dua pasang rasio prorata tetap terpisah:**

- **`ProratePercent*`** (persen, /365, skala 6) → **modul premi**, masuk Seam 3;
- **`Prorate*`** (pecahan, /periode polis, skala 20) → **modul delta spreading**, tidak masuk Seam 3.

⛔ **Diajukan ke work owner, tidak diputuskan.** Perluasan ini lebih besar dari usulan sebelumnya;
work owner perlu menimbang apakah empat bentuk itu satu fungsi bercabang atau empat implementasi di
balik satu pintu.

---

## 6. Kejanggalan — diport apa adanya (K-046)

| # | Butir | Nama test yang disarankan |
| ---: | --- | --- |
| 1 | **`.Rate = 20` tertanam literal** — `calculatePremiPA` L3059 | `K046_PA_RateLiteral20` |
| 2 | **Satuan `ProRatePercent` ambigu** — di dalam pembagi (FIRE) vs pengali langsung (`SetLocalNonMbuProrate` L4304) | `K046_ProRatePercent_SatuanAmbigu` |
| 3 | **`primary.TSI = (primary.TSI × 25)/100`** dan `primary.TSI / 2` — Travel L1745/L1577, tanpa penjelasan | `K046_Travel_PenyesuaianTSI` |
| 4 | Butir dari `10-banding`: `.PremiRp` rasio tertukar · `FirstLossScale` tanpa guard | (sudah tercatat) |

---

## 7. Pertanyaan terbuka

| # | Pertanyaan | Pemilik | Memblokir kontrak Seam 3? |
| --- | --- | --- | --- |
| **1** | **Satuan `ProRatePercent`** — persen atau pecahan pada pemakaian pengali langsung (§3.4). Salah baca = premi geser 100× | **work owner + Aktuaria** | ⛔ **YA** |
| **2** | **Tabel tarif Travel** — `ViewPremiTravel_SQL` nihil di korpus (§2.2) | **DBA** | ⚠️ ya, untuk lini Travel |
| **3** | **Tabel pita tarif PA** — `PCT_RATE_BAWAH`, `PCT_MIN_RATE`, `MIN_TSI`, `MAX_TSI` (§1.2) | **DBA** | ⚠️ ya, untuk jalur pita PA |
| **4** | **Arti `.CalculateMethod`** — nilai `"2"` vs lainnya (§4) | work owner | tidak |
| **5** | **Empat bentuk: satu fungsi bercabang atau empat implementasi** di balik Seam 3 (§5.3) | **work owner** | ⛔ **YA** |

⛔ Butir **1** dan **5** harus dijawab sebelum kontrak Seam 3 dikunci. Butir 2 dan 3 tidak
memblokir kontrak — keduanya masukan dari **repository**, bentuknya sudah jelas meski isinya belum.

---

## 8. Ringkasan

✅ **Tuntas di putaran ini:**

1. **PA — rumus dasar SAMA** dengan NB; berkas bersama identik fungsional. Pembagi **100.000**
   (bukan 10⁹) adalah **konsekuensi K-018**, bukan rumus lain. EDM menambah jalur pita tarif.
2. **Travel — BEDA BENTUK.** Tidak ada rate; premi dari tabel tarif + tambahan per minggu.
   Tabelnya **di luar korpus**.
3. **Prorata — DUA pasang rasio berbeda**, bukan satu nilai di dua properti. Berbeda penyebut (365
   vs periode polis), satuan (persen vs pecahan), dan skala (6 vs 20).
4. **Temuan tambahan:** `SetLocalNonMbuProrate` memuat **bentuk keempat** — dekomposisi delta premi
   tingkat coverage.

⚠️ **Dampak:** kontrak Seam 3 lebih besar dari perkiraan `10-banding` — **empat bentuk**, bukan satu
bentuk dengan varian. **Tetap tanpa seam baru.**

⛔ **`/to-spec` dan `/to-tickets` TIDAK dijalankan.**

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
