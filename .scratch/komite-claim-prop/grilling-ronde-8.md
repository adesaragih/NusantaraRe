# Grilling Ronde 8 — Komite Claim Prop

**Tanggal:** 2026-09-18 · **Korpus:** `D:\XML\RNM_BRD\Komite Claim Prop\` (READ-ONLY)
**Lingkup:** Komite Claim Prop saja `[keputusan work owner]` 2026-09-18.
**Sasaran:** dua bab terakhir — **rumus angka uang** dan **layar kedua**.
**Metode:** parser XML bersarang, **penyaring tidak peka huruf besar-kecil** (§A2).
**Tidak ada nomor baris XML, tidak ada kode Java, dan tidak ada spec di berkas ini.**

---

## 0. Ringkasan

| # | Hasil | Status |
| --- | --- | --- |
| **§A2** | Ekspor Pega tidak konsisten kapitalnya | ⭐ terbukti · **12 penyaring lain aman** |
| **§B1** | Activity yang benar-benar menghitung uang | **5 dari 10** — ⚠️ **plus satu di luar daftar** |
| **§B4** | Pembulatan | ⭐ **NOL di seluruh modul** · presisi hanya lewat `@divide` |
| **§B4** | Presisi pembagian persen | ⚠️ **TIDAK SERAGAM** — `10` · `4` · polos |
| **§B7** | Aritmetika di dalam perulangan | ⚠️ **91 dari 121** |
| **§C** | Layar kedua `ViewDetailInterest` | ⭐ **3 kolom, nol yang baru** |
| **§C6** | Layar ketiga | ✅ **TIDAK ADA** — keempat berkas tampilan sudah dibaca |
| **§D1** | Sembilan bab spec | ⭐ **LENGKAP bahannya** |

---

## §A — Keputusan dan satu aturan parser

### A1 — Spec ditulis **sesudah** dua bab ini lengkap `[keputusan work owner]` 2026-09-18

⛔ **Tidak ada spec ditulis di ronde ini.** Berkasnya tetap `grilling-ronde-8.md`.

### A2 — ⭐ ATURAN PARSER: **ekspor Pega tidak konsisten kapitalnya**

> **SETIAP penyaringan nilai tag harus TIDAK PEKA huruf besar-kecil. Dan lebih baik menyaring
> dari ISI daripada dari LABEL.**

`[terverifikasi]` Pembuktian silang pada langkah `Java`:

| Cara menyaring | Hasil |
| --- | --- |
| `<pyStepsActivityName>` **== `"Java"`** *(peka huruf)* | **17** |
| `<pyStepsActivityName>` **lower == `"java"`** *(tidak peka)* | **20** |
| **`pyStepsJavaSource` terisi** *(menyaring dari ISI)* | **20** |

⭐ **Ketiga selisihnya di `HTMLToPDF`** — langkah **1**, **5**, dan **6**, metodenya ditulis
**`java` huruf kecil**. Menyaring dari **isi** langsung memberi angka benar tanpa perlu tahu
ejaannya.

`[terverifikasi]` Sebaran ejaan seluruh nilai `pyStepsActivityName` di 35 activity: **hanya satu**
nilai yang punya lebih dari satu ejaan — `Java` **17** dan `java` **3**.

### ⚠️ Audit hitungan ronde 1–7 yang memakai penyaring peka huruf

⛔ **Didaftarkan, tidak diperbaiki.** Dua belas penyaring diperiksa ulang peka vs tidak-peka:

| Penyaring | Peka | Tidak peka | Beda? |
| --- | --- | --- | --- |
| `pxObjClass` = `Embed-Display-Table-Cell` | 367 | 367 | ✅ sama |
| `pxObjClass` = `Embed-ActivityPreConditions` | 570 | 570 | ✅ sama |
| `pxObjClass` = `Embed-Harness-HeaderElements` | 289 | 289 | ✅ sama |
| `pyStepsBlockName` = `//` | 24 | 24 | ✅ sama |
| `pyStepsPreCondition` = `true` | 201 | 201 | ✅ sama |
| `pyStepsTransition` = `true` | 12 | 12 | ✅ sama |
| `pyIsVisibilityOption` = `CONDITION` | 7 | 7 | ✅ sama |
| `pyVisible` = `OTHER` | 15 | 15 | ✅ sama |
| `pyType` = `FIELD` | 132 | 132 | ✅ sama |
| `pyStepsRepeatDefHasRepeat` = `REPEAT` | 20 | 20 | ✅ sama |
| `pyStepsRepeatDefHasRepeat` = `EMBEDDED` | 65 | 65 | ✅ sama |
| `pxRuleObjClass` = `Rule-Obj-Class` | 323 | 323 | ✅ sama |

⭐ **NOL hitungan ronde 1–7 yang terpengaruh.** Satu-satunya korban adalah sensus `Java`, dan itu
**sudah benar** karena ronde 7 menyaring dari `pyStepsJavaSource`, bukan dari label.

---

## §B — Rumus angka uang

### B1 — Mana yang benar-benar menghitung uang `[terverifikasi]`

⛔ **Tidak dipercaya dari namanya.** Diuji dari ada-tidaknya penugasan beraritmetika.

| Activity | Kelas (`pxInsName`) | Menghitung uang? | Alasan |
| --- | --- | --- | --- |
| **`CountEstimation_Act`** | `WORK-CLAIMTREATY` | ✅ **YA** | inti penghitung estimasi — `.EstimationValue` · `.ConvertValue` · total per mata uang |
| **`CountListClaimAmountIDR`** | `WORK-CLAIMTREATY` | ✅ **YA** | `.Value` · `.USD` · dua total |
| **`CountPersen_act`** | `WORK-CLAIMTREATY` | ✅ **YA** | `.ClaimSpreaded` · `.ClaimEstimation` |
| **`CountSpreading_act`** | `WORK-CLAIMTREATY` | ✅ **YA** | `.ClaimSpreaded` di empat tempat |
| **`AddEstimation_Act`** | `WORK-CLAIMTREATY` | ⚠️ **sebagian** | hanya **satu** akumulasi (`ClaimSpreaded + Local.Value`); sisanya menyalin |
| `AddLossAllocation_act` | `WORK-CLAIMTREATY` | ⛔ tidak | hanya menyalin ke `SpreadingRisk` + menghitung **panjang daftar** |
| `CheckEstimateDate_Act` | `WORK-CLAIMTREATY` | ⛔ tidak | membandingkan **tanggal**, bukan uang |
| `CurencyEstimation_Act` | `WORK-CLAIMTREATY` | ⛔ tidak | **nol** aritmetika — hanya menaruh hasil pencarian kurs |
| `SetCurencyList_act` | `WORK-CLAIMTREATY` | ⛔ tidak | idem |
| `SetDefNonCatastrope_Act` | `WORK` | ⛔ tidak | nol aritmetika |

⭐ **5 dari 10.**

> ### ⚠️ TEMUAN BESAR: kesepuluhnya berkelas **CLAIM**, bukan komite
>
> `[terverifikasi]` Kelima penghitung berkelas **`ASM-FW-GCNMFW-Work-ClaimTreaty`** — **kelas kasus
> KLAIM**, bukan `Work-KomiteTreaty`. Menurut aturan awalan §A4 ronde 5, **data klaim dibaca dari
> Claim Prop, tidak digarap modul ini**. Jadi **rumus-rumus ini tercatat sebagai bahan, tetapi
> pemiliknya Claim Prop.**
>
> ⭐ **Dan ada satu penghitung uang yang TIDAK ada di daftar sepuluh** — **`SetKomiteList_Act`**,
> berkelas **`ASM-FW-GCNMFW-Work-KomiteTreaty`**. **Itulah satu-satunya penghitung uang milik
> modul ini.** Sudah dibedah §D ronde 5; rumusnya diulang di B2 karena ia yang relevan.

### B2 — Rumus, apa adanya `[terverifikasi]`

⛔ **Tidak diusulkan tipe data, presisi, atau nama kolom.**

**Milik modul ini — `SetKomiteList_Act` step 7.1.1** *(di dalam perulangan `7` → `7.1`)*:

```
Local.TotalAdjGross      := Local.TotalAdjGross      + .GrossAdjustment
Local.TotalAdjRNM        := Local.TotalAdjRNM        + .AdjustmentValue
Local.TotalAdjGrossinIDR := Local.TotalAdjGrossinIDR + (.GrossAdjustment  * .KursIDR)
Local.TotalAdjRNMinIDR   := Local.TotalAdjRNMinIDR   + (.AdjustmentValue  * .KursIDR)
```
lalu step **7.2**: `.AdjustmentGross := Local.TotalAdjGross` · `.AdjustmentValue := Local.TotalAdjRNM`
Gerbang penjumlahan: **`.AcceptanceStatus != 2`** dan **`local.Currency == .Currency`** — baris
yang **ditolak tidak ikut dijumlahkan**.

**Milik Claim Prop — `CountEstimation_Act` step 3.1**:

```
.PersenRNM            := @toDecimal(pyWorkPage.TreatyInMaster.RNMShareP)
.EstimationValue      := @divide(@toDecimal(RNMShareP),100,10) * @toDecimal(.GrossEstimationPct)
.ConvertValue         := @toDecimal(.EstimationValue)      * @toDecimal(.KursValue)
.ConvertGrossEstimasi := @toDecimal(.GrossEstimationPct)   * @toDecimal(.KursValue)
```
step **3.2** menumpuk: `Local.TotalEstimasiIDR += .ConvertValue` ·
`Local.TotalEstimasiGrossIDR += .ConvertGrossEstimasi`, lalu ditaruh ke
`pyWorkPage.ClaimData.TotalEstimasiIDR` dan `.TotalGrossEstimateIDR`.
step **8.1**: `Local.TotalEstimasi += .EstimationValue` → `ClaimData.TotalEstimasi`;
`Local.TotalEstimasiGross += .GrossEstimationPct` → `ClaimData.TotalGrossEstimateTreaty`.

**`CountListClaimAmountIDR` step 1.1** *(di dalam perulangan)*:

```
.Value         := ( @divide(@toDecimal(pyWorkPage.ClaimData.ShareCeding),100,10) * .ClaimAmount ) - .NetDeductibleValue
.USD           := .Value * .IDR
Local.Total    += .Value        ->  ClaimData.TotalListClaimAmount
Local.TotalIDR += .USD          ->  ClaimData.TotalListClaimAmountIDR
```

⚠️ Nama `.USD` **menyesatkan** — isinya hasil kali dengan `.IDR`, bukan dolar.

**`CountPersen_act`**:

```
step 5.1    Local.Value      := @toDecimal(.ClaimEstimation) * @divide(ShareCeding,100,4)
step 5.2.1  .ClaimSpreaded   := .SharePercentage * Local.Value / 100
step 5.2.1  .ClaimEstimation := .ClaimSpreaded  * .PremiumSpreaded
step 6.2.1  .ClaimSpreaded   := @divide(@toDecimal(.SharePercentage),100,10) * Local.Value
step 6.2.1  .ClaimEstimation := .ClaimSpreaded  * .PremiumSpreaded
```

⚠️ **Step 5.2.1 dan 6.2.1 menghitung hal yang SAMA dengan cara berbeda** — satu memakai
`/100` polos, satu memakai `@divide(…,100,10)`.

**`CountSpreading_act`**:

```
step 7      SpreadingClaim(Param.Index).ClaimSpreaded := @toDecimal(SharePercentage) * @toDecimal(TotalEstimasi) / 100
step 9.2.2  SpreadingClaim(Local.IndexEst).ClaimSpreaded := @toDecimal(SharePercentage) * @toDecimal(Local.Estimation) / 100
step 10.1   .ClaimSpreaded := ( SpreadingClaim(Param.Index).ClaimSpreaded * .SharePercentage ) / 100
step 11.2.1 .ClaimSpreaded := ( SpreadingClaim(Local.IndexEst).ClaimSpreaded * .SharePercentage ) / 100
```

⭐ **Pola tetap di seluruh modul: persen selalu dibagi 100, dan hasil kali selalu tanpa pembulatan.**

### B3 — Kurs: di mana diubah ke IDR, dan kapan dikunci `[terverifikasi]`

| Tempat | Apa yang terjadi |
| --- | --- |
| `SetCurencyList_act` **1** · `CurencyEstimation_Act` **1** | `SearchNilaiKursInput.Currency := Param.CurrID` — **permintaan** kurs |
| `AddEstimation_Act` **5.1** | `SearchNilaiKursInput.Currency := .CurrencyID` |
| `CurencyEstimation_Act` **5** · `SetCurencyList_act` **5** · `AddEstimation_Act` **5.3** | **hasilnya disimpan**: `… := SearchNilaiKursOutput.pxResults(1).NILAIKURS` |
| `CountEstimation_Act` **3.1** | **perkalian ke IDR**: `.ConvertValue := .EstimationValue * .KursValue` |
| `CountListClaimAmountIDR` **1.1** | `.USD := .Value * .IDR` |
| **`SetKomiteList_Act` 7.1.1** | ⭐ **satu-satunya di kelas komite**: `… + (.GrossAdjustment * .KursIDR)` |
| `SaveAcceptation_Act` **1** | `TempOSAkseptasi.KursValue := …AdjustmentList(Param.IdxAdj).KursIDR` — **kurs disalin ke OS akseptasi** |

✅ **Silang dengan K8** *(tertutup ronde 5)*: pengisi `SearchNilaiKursOutput` adalah rule
Connect-SQL **`ASM!CURRENCYSTANDARD`**, masukannya **mata uang + `SYSDATE`**, jadi kurs yang
diambil adalah **kurs tanggal lookup dijalankan**.

⭐ **Penguncian terjadi di luar modul ini** — pola *ambil-hanya-bila-kosong*
(`.KursIDR := @if(.KursIDR=="", …, .KursIDR)`) ada di `Claim Prop/Activity/SetNameCurrency_Act`.
`[terverifikasi]` **Di dalam modul ini nol pola penguncian untuk uang**: keempat pola
*ambil-hanya-bila-kosong* yang ada seluruhnya di `HitServiceToKasirKMT_Act`, dan isinya **nomor
rekening dan operator**, bukan angka uang.

### B4 — Pembulatan: ⭐ **TIDAK ADA, di seluruh modul**

`[terverifikasi]` Sensus **80 berkas**, tidak peka huruf:

| Yang dicari | Ditemukan |
| --- | --- |
| `@round` | ⛔ **0** |
| `@Truncate` / `@trunc` | ⛔ **0** |
| `@FormatNumber` | ⛔ **0** |
| `Math.round` | ⛔ **0** |
| `setScale` | ⛔ **0** |
| `RoundingMode` | ⛔ **0** |
| `@toDecimal` | **54** di 7 berkas |
| `BigDecimal` | 3, hanya `HitServiceToKasirKMT_Act` |

⛔ **Nol pembulatan. Ditegaskan dari sensus seluruh modul, bukan dari satu berkas.**

⭐ **Satu-satunya kendali presisi adalah argumen ketiga `@divide`** — dan ⚠️ **ia tidak seragam.**
`[terverifikasi]` Seluruh **8** pemakaian `@divide`, dibaca dengan kurung berimbang:

| Rumus | Presisi | Di mana |
| --- | --- | --- |
| `@divide(@toDecimal(RNMShareP),100,10)` | **10** | `CountEstimation_Act` |
| `@divide(@toDecimal(ShareCeding),100,10)` | **10** | `CountListClaimAmountIDR` |
| `@divide(@toDecimal(.SharePercentage),100,10)` | **10** | `CountPersen_act` 6.2.1 |
| `@divide(ShareCeding,100,4)` | ⚠️ **4** | `CountPersen_act` 5.1 |
| `@divide(.TotalSpreadAdjustment,1,4)` ×2 | **4** | `SendEmailKlaim_KMT` |
| `@divide(.TotalSharePersen,1,0)` · `@divide(.SharePercentage,1,0)` | **0** | `SendEmailKlaim_KMT` |

⚠️ **Pembagian persen yang sama dikerjakan tiga cara:** `@divide(x,100,10)` · `@divide(x,100,4)` ·
dan **`/100` polos tanpa kendali presisi sama sekali** (`CountPersen_act` 5.2.1,
`CountSpreading_act` 7 · 9.2.2 · 10.1 · 11.2.1). ⛔ **Tidak ditebak mana yang benar.**
`[terbuka]` **BARU**.

### B5 — Disimpan vs dihitung ulang `[terverifikasi]`

| | Jumlah |
| --- | --- |
| Angka uang yang **dikunci** (*ambil-hanya-bila-kosong*) **di modul ini** | ⛔ **0** |
| Angka uang yang **dihitung ulang** tiap kali activity penghitung jalan | **seluruhnya** |

⭐ **Seluruh angka uang dihitung ulang, lalu ditulis ke properti.** Pola penulisannya seragam:
penampung `Local.*` **dinolkan** lebih dulu (`Local.TotalEstimasi := 0`), diisi lewat perulangan,
lalu **ditimpa** ke properti kasus. **Tidak ada satu pun yang diperiksa "sudah terisi atau belum".**

⚠️ **Satu-satunya nilai yang dikunci adalah KURS**, dan penguncinya **di luar modul ini** (B3).

### B6 — Rantai ketergantungan, dalam bahasa biasa

1. **Kurs dicari lebih dulu** — mata uang dikirim ke `CurrencyStandard`, hasilnya disimpan sebagai
   `KursValue` / `KursIDR`. Tanpa ini, tidak ada angka IDR yang bisa dihitung.
2. **Persentase share diambil dari master** — `TreatyInMaster.RNMShareP` dan `ClaimData.ShareCeding`.
3. **Nilai estimasi per baris** dihitung: persen ÷ 100 × nilai kotor.
4. **Nilai IDR per baris** dihitung: nilai estimasi × kurs.
5. **Total per mata uang** ditumpuk dengan menyisir daftar, sesudah baris bermata-uang-sama
   digabung oleh langkah `Java` pembersih duplikat.
6. **Total kasus** ditulis ke `ClaimData.TotalEstimasi` · `TotalEstimasiIDR` ·
   `TotalGrossEstimateIDR` · `TotalGrossEstimateTreaty`.
7. **Spreading** dihitung **sesudah** total ada: share ÷ 100 × total estimasi.
8. ⭐ **Baru di sini komite masuk** — `SetKomiteList_Act` menjumlahkan baris adjustment yang
   **tidak ditolak**, per mata uang, dan mengalikannya dengan `KursIDR` untuk mendapat IDR.

⛔ **Langkah 1–7 milik Claim Prop. Langkah 8 milik modul ini.**

### B7 — ⚠️ Aritmetika di dalam perulangan dan `Java`

`[terverifikasi]` Dari **121 penugasan beraritmetika** di 35 activity:

| | Jumlah |
| --- | --- |
| **Di dalam langkah berulang** (`EMBEDDED` / `REPEAT`) | ⚠️ **91** |
| Di luar perulangan | 30 |
| **Di dalam langkah `Java`** | ⛔ **0** — kedua puluh Java **tidak menghitung uang** |

⭐ **Tiga perempat aritmetika modul ini terjadi di dalam perulangan.** Berapa kali angkanya
berubah **ditentukan oleh berapa kali induknya berputar** — dan itu **masih `[terbuka]`** (§C2
ronde 7: halaman yang diulang tidak terbaca dari struktur).

⭐ **Kabar baiknya: nol perhitungan uang di `Java`.** Kedua puluh langkah `Java` hanya membersihkan
duplikat, merakit teks, dan menangani berkas — **tidak satu pun menyentuh angka uang**. Jadi rumus
uang **seluruhnya terbaca dari `Property-Set`**, tidak ada yang tersembunyi.

---

## §C — Layar kedua: `ViewDetailInterest`

### C1 — ⚠️ Kedua berkas itu **RULE BERBEDA** `[terverifikasi]`

⭐ Persis pelajaran §11 ronde 1 — **nama berkas sama bukan bukti rule sama**:

| Berkas | `pxInsName` | Kelas |
| --- | --- | --- |
| `FlowAction/ViewDetailInterest.xml` | **`ASM-FW-GISFW-DATA-TREATYINTOTAL!VIEWDETAILINTEREST`** | ⚠️ **`ASM-FW-GISFW-Data-TreatyInTotal`** — **bukan kelas komite** |
| `Section/ViewDetailInterest.xml` | **`ASM-FW-GCNMFW-WORK-KOMITETREATY!VIEWDETAILINTEREST`** | ✅ `ASM-FW-GCNMFW-Work-KomiteTreaty` |

**FlowAction — pembungkus, seperti `ViewTransferDtl`:** **2** sel, **nol** medan.
`pySectionReference = ViewDetailInterest`, `pyPreProcessingActivity` **kosong**,
tombol `Submit` / `Cancel`. Rule yang dirujuknya: Section `VIEWDETAILINTEREST` **pada kelas
`Data-TreatyInTotal`**, `CEKINTERESTLISTDTL_DT` (`Rule-Obj-Model`), dan FieldValue judul.

⚠️ **Jadi FlowAction ini merujuk Section pada kelas LAIN**, bukan Section komite yang ada di
folder yang sama. ⛔ **Tidak ditelusuri** — kelasnya di luar modul ini.

### C2 · C3 — Section komite: **3 kolom** `[terverifikasi]`

Meski berkasnya **115 KB**, isinya hampir seluruhnya fragmen HTML dan pembungkus grid bawaan Pega.

| Ukuran | Jumlah |
| --- | --- |
| Sel `Embed-Display-Table-Cell` | **15** |
| **Sel medan** | **3** |
| **Properti berbeda** | **3** |
| — milik **KOMITE** | **3** |
| — milik **KLAIM induk** | ⛔ **0** |

| Properti | n | Hanya-baca | Wajib isi | Gerbang | Label |
| --- | --- | --- | --- | --- | --- |
| `.CurrencyID` | 1 | ✅ ya | ✅ **wajib** | — | |
| `.ObjectName` | 1 | ✅ ya | ✅ **wajib** | — | *Text Input* |
| `.TSIPerObject` | 1 | ✅ ya | ✅ **wajib** | — | |

⭐ **Ketiganya hanya-baca, nol gerbang, nol gerbang mati.** Layar ini **murni menampilkan** —
penyetuju tidak bisa mengubah apa pun di sini.

### C4 — Dari mana layar ini dibuka

⚠️ `[terbuka]` **Tidak terbaca dari modul ini.** Section komite `VIEWDETAILINTEREST` **tidak
dirujuk** oleh `ShowTransfer`, `ViewTransferDtl`, maupun berkas `Flow`. Satu-satunya FlowAction
bernama sama **merujuk Section pada kelas lain**. ⛔ **Tidak ditebak.**

### C5 — Silang dengan `ShowTransfer`: ⭐ **nol yang baru**

| | Jumlah |
| --- | --- |
| Properti `ShowTransfer` | **93** |
| Properti `ViewDetailInterest` | **3** |
| **Irisan** — muncul di kedua layar | **3** — `.CurrencyID` · `.ObjectName` · `.TSIPerObject` |
| ⭐ **HANYA di `ViewDetailInterest`** | ⛔ **0** |

⭐ **Layar kedua tidak menambah satu kolom pun.** Daftar **93 kolom** hasil ronde 6 **sudah
lengkap** untuk modul ini.

### C6 — Layar ketiga: ✅ **TIDAK ADA** `[terverifikasi]`

| Folder | Berkas | Sudah dibaca? |
| --- | --- | --- |
| `Section` | `ShowTransfer.xml` *(348 sel)* · `ViewDetailInterest.xml` *(15 sel)* | ✅ **keduanya** |
| `FlowAction` | `ViewTransferDtl.xml` *(2 sel)* · `ViewDetailInterest.xml` *(2 sel)* | ✅ **keduanya** |

⭐ **Empat berkas tampilan, keempatnya sudah dibaca. Nol yang tersisa.**
Folder `Harness` **ada dan kosong** (§1 ronde 1) — konsisten.

---

## §D — Kesiapan spec

### D1 — Sembilan bab lengkap bahannya? ⭐ **YA**

Kedua bab yang ronde 7 tandai `[terbuka]` **kini punya bahan**:

| Bab yang tertahan | Keadaan sekarang |
| --- | --- |
| **Angka uang** | ✅ **lengkap** — 5 penghitung teridentifikasi, rumusnya tertulis apa adanya, rantai ketergantungan tergambar, nol pembulatan terbukti. ⚠️ **Dengan catatan pemilik**: empat dari lima berkelas **CLAIM**, jadi bab ini **merujuk** rumus Claim Prop; yang milik modul ini hanya `SetKomiteList_Act` |
| **Layar kedua** | ✅ **lengkap** — 3 kolom, nol tambahan; **daftar 93 kolom sudah final** |

### D2 — Kerangka spec yang diusulkan — **9 bab**

⛔ **Kerangkanya saja. Isinya tidak ditulis di sini.**

| # | Bab | Satu kalimat isinya | Bahan dari |
| --- | --- | --- | --- |
| **1** | **Lingkup dan identitas** | Kelas kerja komite, 80 rule, silsilah Save-As, dan aturan awalan yang memilah data komite dari data klaim. | ronde 1 §1 §2a §3.6 · ronde 5 §A4 |
| **2** | **Daur hidup kasus** | Kasus dibuat → masuk kotak kerja penyetuju → keputusan → kembali ke kotak kerja atau selesai. | **ronde 6 §D** |
| **3** | **Tangga penyetuju** | Siapa giliran berikutnya, kapan tangga berhenti, dan bagaimana sisa penyetuju ditolak otomatis. | ronde 1 §3.2 §9.1 §9.2 · ronde 4 §B · ronde 7 §C4 |
| **4** | **Layar komite** | 93 kolom, dua famili gerbang, tiga wajah menurut jalur, dan dua isian yang selalu bisa disunting. | ronde 2 §B · ronde 3 §B · ronde 6 §C · **ronde 8 §C** |
| **5** | **Tiga jalur** | TRANSFER ADJUSTMENT, REJECT, dan CLOSE — apa bedanya di layar dan di mesin. | ronde 1 §3.5 §9.8 · ronde 5 §A5 |
| **6** | **Nomor akseptasi** | Kapan nomor dibuat, dari rule apa, dan mengapa hanya sekali di tingkat akhir. | ronde 2 §C · ronde 4 §B3 · ronde 7 §C3 |
| **7** | **Angka uang** | Dari mana kurs diambil, bagaimana total adjustment dijumlahkan, dan rumus klaim yang dirujuk dari Claim Prop. | **ronde 8 §B** · ronde 5 §D |
| **8** | **Efek keluar** | Delapan hal yang keluar saat komite menyetujui, dan dua sikap berbeda terhadap kegagalan. | ronde 5 §B · **ronde 7 §B4** |
| **9** | **Tabel yang disentuh** | Tabel yang ditulis, sembilan rule ber-`COMMIT` sendiri, dan DDL dua tabel dari DBA. | ronde 1 §3.7 · **ronde 3 §12** |

### D3 — Butir `[terbuka]` yang **ikut terbawa ke spec sebagai tanda terbuka**

⛔ **Tidak ditutup sendiri.**

| # | Butir | Bab yang membawanya |
| --- | --- | --- |
| 1 | `OPERATORID` `HISTORYAKSEPTASIPEGA` tak pernah diisi jalur komite | bab 9 |
| 2 | **F6** — pewarisan kelas kerja, tidak ada di jendela ini | bab 1 |
| 3 | arti `REPEAT` vs `EMBEDDED`, dan halaman apa yang diulang | bab 3 · bab 7 |
| 4 | **berapa kali** perulangan berputar — menentukan berapa kali angka berubah dan berapa kali REST Kasir dipanggil | bab 7 · bab 8 |
| 5 | urutan evaluasi bila kedua keluarga gerbang terisi | bab 3 |
| 6 | unggah Google Storage dan email **tanpa penanganan gagal** | bab 8 |
| 7 | `InsertGoogleStorage_Act` bercatat *"FIX ERROR HANDLING"* tanpa jejak | bab 8 |
| 8 | ⭐ **BARU** — **presisi pembagian persen tidak seragam** (`10` · `4` · polos) dan **nol pembulatan** | bab 7 |
| 9 | ⭐ **BARU** — dari mana layar `ViewDetailInterest` komite dibuka | bab 4 |

⛔ **Dan satu tanda yang bukan butir terbuka, tetapi wajib terbawa** — §A1 ronde 6: **urutan efek
keluar terhadap penyimpanan adalah titik yang sengaja diubah, bukan syarat.**

### D4 — Kesimpulan yang **paling rawan salah menjelang spec** — **3**

1. ⛔ **§C2 ronde 7 — "`EMBEDDED` mengulang daftar, `REPEAT` mengulang sekali"**. **Masih paling
   rawan**, dan sekarang **lebih berat**: §B7 menunjukkan **91 dari 121** aritmetika ada di dalam
   perulangan. Bila pembacaan perulangan saya keliru, **bab 7 ikut keliru**.
2. ⚠️ **§B1 ronde ini — "kelima penghitung berkelas CLAIM, jadi milik Claim Prop"**. Ditarik dari
   `pxInsName`, yang memang aturan sah. **Tetapi** berkasnya duduk di folder ekspor modul ini —
   bila lingkup ternyata mengikuti **folder** dan bukan **kelas**, bab 7 berpindah pemilik.
3. ⚠️ **§B4 — "nol pembulatan"**. Kuat *(sensus 80 berkas, tidak peka huruf)*, tetapi pembulatan
   bisa terjadi **di luar korpus** — di fungsi tersimpan Oracle, atau di layanan Kasir.

⛔ **Ditunjuk, tidak diperbaiki.**

### D5 — Pertanyaan **BARU** untuk work owner — **1**

#### Pertanyaan 15 — Persen dibagi seratus dengan tiga cara berbeda. Mana yang benar?

**Apa yang ditanyakan.** Pembagian persentase — operasi yang sama, dipakai untuk menghitung nilai
klaim — dikerjakan **tiga cara berbeda** di modul ini: dengan presisi **sepuluh angka**, dengan
presisi **empat angka**, dan **tanpa kendali presisi sama sekali**. Ditambah, **tidak ada
pembulatan di mana pun**. Mana yang benar, dan berapa angka di belakang koma yang seharusnya.

**Kenapa muncul.** Baru terbaca ronde ini. Dua langkah yang menghitung **hal yang sama** di
`CountPersen_act` memakai cara berbeda — satu `@divide(…,100,10)`, satu `/100` polos. Selisihnya
kecil per baris, tetapi **angka itu dijumlahkan lewat perulangan** dan akhirnya **dikirim ke
Kasir sebagai nilai pembayaran**.

**Bedanya jawaban A atau B.** Bila **ada satu aturan presisi yang sah**, seluruh rumus disamakan
ke aturan itu, dan selisih terhadap data lama perlu dijelaskan. Bila **ditiru apa adanya**,
ketiga cara itu dipindahkan **persis** termasuk ketidakseragamannya — hasilnya sama dengan Pega,
tetapi sulit diterangkan kepada siapa pun yang memeriksa angka.

**Yang tertahan.** Bab 7 spec, dan setiap angka uang yang dikirim ke Kasir.

---

## §E — Daftar `[terbuka]` modul ini sesudah ronde 8 — **9**

| # | Butir | Menunggu |
| --- | --- | --- |
| 1 | `OPERATORID` tak pernah diisi jalur komite | korpus / DBA |
| 2 | **F6** — pewarisan kelas kerja, dicek habis | work owner |
| 3 | arti `REPEAT` vs `EMBEDDED`, halaman apa yang diulang | korpus |
| 4 | **berapa kali** perulangan berputar | korpus |
| 5 | urutan evaluasi bila kedua keluarga gerbang terisi | korpus |
| 6 | unggah Google Storage dan email tanpa penanganan gagal | work owner |
| 7 | `InsertGoogleStorage_Act` *"FIX ERROR HANDLING"* tanpa jejak | korpus |
| 8 | ⭐ **BARU** — presisi pembagian persen tidak seragam · nol pembulatan | work owner *(pertanyaan 15)* |
| 9 | ⭐ **BARU** — dari mana Section `ViewDetailInterest` komite dibuka | korpus |

⛔ **Tidak ada butir dinyatakan tertutup di ronde ini.**

---

## §F — Ralat, dikumpulkan

| # | Di mana | Yang salah | Yang benar |
| --- | --- | --- | --- |
| 1 | **ronde 7 §E1 no.1** | *"12 activity penghitung"* sebagai satu kelompok kerja modul ini | **5 dari 10** yang benar-benar menghitung, dan **empat di antaranya berkelas CLAIM** — pemiliknya Claim Prop (§B1) |
| 2 | **ronde 7 §E2** | bab angka uang *"bahannya belum ada"* | ✅ **kini ada** (§B) |
| 3 | **ronde 7 §E2** | bab layar kedua *"belum lengkap"* | ✅ **lengkap, dan nol tambahan** (§C5) |
| 4 | **ronde 1 §1 B1** | `ViewDetailInterest.xml` dicatat sebagai *"satu nama muncul dua kali"* | ⚠️ keduanya **rule berbeda pada kelas berbeda** — FlowAction berkelas `Data-TreatyInTotal` (§C1) |

⛔ **Tidak satu pun berkas ronde 1–7 disunting.**
