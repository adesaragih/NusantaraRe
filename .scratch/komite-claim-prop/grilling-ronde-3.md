# Grilling Ronde 3 — Komite Claim Prop

**Tanggal:** 2026-09-18 · **Korpus:** `D:\XML\RNM_BRD\Komite Claim Prop\` (READ-ONLY)
**Lingkup:** Komite Claim Prop saja `[keputusan work owner]` 2026-09-18.
**Metode:** parser XML bersarang. Bukti = path berkas + nama rule + nomor step Pega.
**Tidak ada nomor baris XML di berkas ini.** Identitas rule = `pxInsName`, bukan nama berkas.

> **Ralat ke ronde 1 dan ronde 2 dicatat DI SINI**, tidak disunting di tempatnya.

---

## 0. Ringkasan

| # | Hasil | Status |
| --- | --- | --- |
| **§A** | Arti enam kode arah | ✅ **K10 `[tertutup]`** — `[data work owner]` |
| **§A5** | Kesimpulan yang goyah | ⚠️ **4** — satu di antaranya **membatalkan cakupan sensus ronde 2** |
| **§B** | Aturan 1 ronde 2 | ⛔ **DICABUT** — sel **memang** bergerbang sendiri. Famili A **15** · B **9** |
| **§B3** | "Enam properti selalu dapat disunting" | ⛔ **SALAH** — yang benar **dua** |
| **§C** | Definisi kolom layar | 55 *(titik)* · **+62 berawalan halaman** = **117** ≠ 98 |
| **§D** | `SetKomiteList_Act` | ⭐ **tidak mengisi `KomiteList` sama sekali** |
| **§E** | `Data-Comitee` | **7 berkas** |
| **§F** | Pekerjaan berikutnya | **6 butir** · gelap **5** · rawan **3** · pertanyaan baru **4** |

---

## §A — K10 `[tertutup]`: arti kode arah

### A1 — Tabel enam kode `[data work owner]` 2026-09-18

Diambil dari **layar Pega**, bukan tafsiran korpus. Ronde 2 benar bahwa korpus tidak dapat
membuktikannya sendiri; sumber yang sah adalah dropdown Pega.

| Kode | Nama di Pega | Arti |
| --- | --- | --- |
| **1** | **Jump to Later Step** | lompat ke langkah berlabel — sasarannya di **kolom *true param*** |
| **2** | **Continue Whens** | periksa baris syarat berikutnya; **habis baris → langkah dijalankan** |
| **3** | **Skip Step** | langkah **tidak** dijalankan, lanjut ke langkah berikutnya |
| **4** | **Exit Iteration** | putus perulangan |
| **5** | **Skip Whens** | **berhenti memeriksa baris syarat**, langkah **langsung dijalankan** |
| **6** | **Exit Activity** | hentikan activity |
| *(kosong)* | — | dropdown memang kosong, **tidak ada tindakan dipilih** |

⭐ **Baris syarat adalah RANTAI, dibaca dari atas.** `Continue Whens` + `Skip Step` menyusun
**AND**; satu `Skip Whens` mengubahnya menjadi **OR**.

`[terverifikasi]` Letak sasaran lompatan adalah **`pyStepsPreCondParamsWhenTruePrms`** /
**`…WhenFalsePrms`** — bukan `pyStepsPreCondParamsWhen…Param` yang saya cari di ronde 2. Itu sebab
ronde 2 melaporkan "TParam kosong".

### A2 — Enam baris berkode 1 / 4 / 5, dalam bahasa biasa `[terverifikasi]`

| Berkas · step | Kode | Gerbang | Artinya |
| --- | --- | --- | --- |
| `KomitePostAdjustment` **4** | `_/1` → `JMP` | `false` — **mati** | *Seandainya hidup:* bila `TransferType` **bukan** 2, **lompat ke label `JMP`**. Karena `pyStepsPreCondition = false`, gerbangnya **tidak berlaku** dan langkah selalu dijalankan. |
| `KomitePost_Close` **1** | `_/1` → `EXT` | `false` — **mati** | *Seandainya hidup:* bila **bukan** (`TransferType=1` **dan** `TypeComentAnalysis=5`), **lompat ke `EXT`**. Gerbang tidak berlaku. |
| `KomitePost_Reject` **1** | `_/1` → `EXT` | `false` — **mati** | sama persis dengan di atas. |
| **`KomitePostAdjustment` 12** | `1/3` → `EXT` | `true` — **HIDUP** | Bila `AcceptStatus == "2"` (**ditolak**) → **lompat ke label `EXT`, yaitu langkah 25**. Bila tidak ditolak → **Skip Step**, langkah 12 dilewati. |
| **`KomitePostAdjustment` 26** | `_/4` | `true` — **HIDUP** | Bila syaratnya **tidak** terpenuhi → **Exit Iteration**, perulangan penyetuju diputus. Bila terpenuhi → tidak ada tindakan dipilih, langkah dijalankan. |
| **`SendErrorDirectKasir` 2.3** | `5/2` | `true` — **HIDUP** | Baris 1 `IsCLMP`: bila benar → **Skip Whens**, langkah **langsung dijalankan** tanpa memeriksa baris 2. Bila salah → **Continue Whens**, periksa baris 2 `IsCLMNP`. **Inilah OR.** |

⭐ **`SendErrorDirectKasir` step 2.3 membuktikan kaidah OR.** Keterangan Pega-nya sendiri berbunyi
**`CLMP || CLMNP`** — dan susunan kodenya `5/2` lalu `2/3` memang menghasilkan
*"CLMP **atau** CLMNP"*. Kode 5 **satu-satunya** di modul ini.

### A3 — Sasaran lompatan `[terverifikasi]`

Label langkah dibaca dari **`pyStepsBlockName`**; `//` berarti langkah **di-remark**.

| Activity | Label yang ada | Sasaran yang dipakai | Cocok? |
| --- | --- | --- | --- |
| `KomitePostAdjustment` | `EXT`→**25** · `ENDSERVICE`→**30** · `ENDKASIR`→**35** | `EXT` (step 12) | ✅ **ya — langkah 25** |
| `KomitePostAdjustment` | *(idem)* | **`JMP`** (step 4) | ⛔ **TIDAK ADA LABEL `JMP`** |
| `KomitePost_Close` | **NOL label** | `EXT` (step 1) | ⛔ **TIDAK ADA LABELNYA** |
| `KomitePost_Reject` | **NOL label** | `EXT` (step 1) | ⛔ **TIDAK ADA LABELNYA** |

⚠️ **Tiga sasaran lompatan menggantung** — `JMP` sekali, `EXT` dua kali di activity yang tidak
punya label sama sekali. **Ketiganya duduk di baris bergerbang `false`**, jadi **tidak pernah
dijalankan**. Konsisten: lompatan yang rusak justru yang dimatikan. **Tidak ditafsirkan lebih jauh.**

⭐ Label `ENDSERVICE` (30) dan `ENDKASIR` (35) **ada tetapi tidak ada satu pun yang melompat ke
sana** di dalam activity ini. Dicatat sebagai temuan.

### A4 — **K10 `[tertutup]`**

Keenam kode punya arti, dan keenam baris berkode 1/4/5 sudah terbaca. `[data work owner]`.

### A5 — Membaca ulang 566 baris: **empat kesimpulan goyah**

#### ⚠️ GOYAH 1 — **cakupan sensus ronde 2 tidak lengkap** *(paling berat)*

`[terverifikasi]` Selain `pyStepsPreCondParams`, setiap langkah juga punya
**`pyStepsTransParams`** — keluarga kode **kedua**, yang **tidak pernah saya hitung**. Di 35
Activity modul ini ada **552 baris** lagi:

| Kode | Jumlah | | Kode | Jumlah |
| --- | --- | --- | --- | --- |
| `2/2` | **536** | | `2/6` | 1 |
| `_/_` | 7 | | `1/1` | 1 |
| `1/2` | **4** | | `1/_` | 1 |
| `6/2` | 1 | | `6/6` | 1 |

⛔ **Jadi "566 baris syarat" di ronde 2 bukan seluruh baris berkode di modul ini.** Ada **552**
lagi, termasuk **enam baris berkode 1** yang belum pernah dibaca. **Angka K10 "enam baris, tiga
hidup" berlaku untuk `pyStepsPreCondParams` saja.** ⛔ **Tidak saya sisir sekarang** — ia masuk
daftar pekerjaan §F1.

#### ⚠️ GOYAH 2 — `2/2` "tanpa gerbang" benar hasilnya, salah alasannya

Ronde 1 §8 dan ronde 2 §E membaca `2/2` (**298** baris) sebagai *"seluruhnya di step tanpa gerbang"*.
Dengan tabel A1: `2/2` = **Continue Whens / Continue Whens** — artinya *"periksa baris berikutnya,
apa pun hasilnya"*; bila tidak ada baris berikutnya, **langkah dijalankan**. **Hasil akhirnya sama**
(langkah selalu jalan), tetapi **mekanismenya bukan "tidak ada gerbang"** — gerbangnya ada dan
selalu lolos. Perbedaan ini penting bila kelak ada baris kedua.

#### ⚠️ GOYAH 3 — `(kosong)/3` **tidak lagi "belum punya arti"**

Ronde 1 §8 menandai `(kosong)/3` (**14** baris, 4 di `KomitePostAdjustment`) sebagai
**"⚠️ K10 — belum punya arti"**, dan menyimpulkan *"langkah nomor akseptasi belum boleh ditulis"*.
Dengan tabel A1 ia terbaca terang: **bila syarat benar → tidak ada tindakan dipilih → langkah
dijalankan; bila salah → Skip Step.** Itu **gerbang biasa**. `[terverifikasi]` Step **16.9**
(penulis `AcceptedNo` dan `AcceptanceStatus := 1`) bergerbang
`KomiteCount==TotalKomite && AcceptStatus=="1"` — **penyetuju terakhir dan disetujui**. Tidak ada
yang misterius.

#### ⚠️ GOYAH 4 — `_/1` bukan sekadar "kode 1"

Ronde 1 §8 mencatat `(kosong)/1` sebagai kode tanpa arti. Kini terbaca: **tidak ada tindakan bila
benar; lompat bila salah**. Ketiganya bergerbang mati, jadi **tidak mengubah perilaku** — tetapi
pembacaan lamanya salah.

#### ✅ Uji AND untuk langkah 17 / 21 / 34 — **TETAP BERDIRI**

`[terverifikasi]` Ketiganya punya **dua** baris syarat, dan **tidak satu pun berkode 5**:

| Step | Baris 1 | Baris 2 |
| --- | --- | --- |
| **17** *(insert OS akseptasi)* | `2/3` `KomiteCount==TotalKomite && AcceptStatus=="1"` | `3/2` `IsSubjectivity==true` |
| **21** *(generate PDF akseptasi)* | `2/3` *(idem)* | `3/2` `IsSubjectivity==true` |
| **34** *(HIT SERVICE KE KASIR)* | `3/2` `IsSubjectivity==true` | `2/3` `KomiteCount==TotalKomite && AcceptStatus=="1"` |

Tanpa `Skip Whens`, rantai tetap **AND**: *penyetuju terakhir **dan** disetujui **dan** bukan
subjectivity.* ✅ **Kesimpulan K1 ronde 1 berdiri utuh.**
⚠️ Kecil: pada step **34** urutan kedua barisnya **terbalik** dari 17 dan 21. Hasil AND sama.

---

## §B — ⛔ RALAT: aturan baca Section. Sel **memang** bergerbang sendiri

> ### ⛔ ATURAN 1 RONDE 2 DICABUT
>
> Ronde 2 §A2 menulis: *"gerbang tampil hidup di tingkat LAYOUT, bukan di tingkat sel; dari 85 sel
> medan, yang punya `pyVisible` = NOL, `pyCondition` = NOL."* **Itu keliru.**
>
> Sebabnya: saya mencari `pyVisible` sebagai **anak langsung** sel. Ia sebenarnya berada di
> **sub-halaman `<pyUserData>`** milik sel (`pxObjClass = Embed-Harness-HeaderElements`):
>
> ```
> <pyUserData>
>   <pxObjClass>Embed-Harness-HeaderElements</pxObjClass>
>   <pyVisible>OTHER</pyVisible>
>   <pyCondition>.AcceptStatus = 1 &amp;&amp; .TransferType =2</pyCondition>
> </pyUserData>
> ```
>
> ⭐ **Pelajaran, satu baris: "tidak ketemu" hanya sah bila seluruh kedalaman sudah dicari.**
> Sepupu pelajaran §11 ronde 1.

### B1 — Dua famili gerbang, dihitung sendiri `[terverifikasi]`

| Famili | Tingkat | Saklar | Syarat | Hidup |
| --- | --- | --- | --- | --- |
| **A** | **SEL** | `pyUserData/pyVisible = OTHER` | `pyUserData/pyCondition` | **15** |
| **B** | **LAYOUT** | `pyIsVisibilityOption = CONDITION / ExpressionCondition` | `pyContainerVisibleWhen` | **9** |

✅ **Cocok persis: 15 dan 9.**

Sebaran `pyUserData/pyVisible` atas **354** blok: `ALWAYS` **328** · **`OTHER` 15** ·
`NOTBLANK` **5** · kosong **6**.
Sebaran `pyIsVisibilityOption`: `ALWAYS` **32** · `CONDITION` **7** · `ExpressionCondition` **2**.

**Kelima belas gerbang famili A, berikut sel pemiliknya** `[terverifikasi]`:

| Syarat | Properti sel | Jenis |
| --- | --- | --- |
| `.TransferType =='2'` | judul **ADJUSTMENT** | LABEL |
| `.TransferType =='3'` | judul **REJECT** | LABEL |
| `.TransferType =='4'` | judul **CLOSE** | LABEL |
| `.TransferType =2` | `.Adjustment.IsProposeClose` | FIELD |
| `.TransferType =2` | `.Adjustment.IsPropReserved` | FIELD |
| `.AcceptStatus = 1 && .TransferType =2` | `.IsSubjectivity` | FIELD |
| `.IsSubjectivity = true` | `.SubjectivityNote` | FIELD |
| `.Adjustment.SwiftCode != ''` | `.Adjustment.SwiftCode` | FIELD |
| `pyWorkCover.ClaimData.AppointedADJ!=''` | `pyWorkCover.ClaimData.AppointedADJ` | FIELD |
| `pyWorkCover.ClaimData.ConsultantName!=''` | `pyWorkCover.ClaimData.ConsultantName` | FIELD |
| `pyWorkCover.ClaimData.StsKatastrofe ='Non-Catastrophe'` | `pyWorkCover.ClaimData.NonKatastrofeType` | FIELD |
| **`NEVER`** | `pyWorkCover.ClaimData.NoClaim` | FIELD |
| **`NEVER`** | `pyWorkCover.ClaimData.TotalGrossEstimate…` | FIELD |
| **`NEVER`** | `pyWorkCover.ClaimData.TotalEstimasi…` | FIELD |
| **`NEVER`** | *sample text* | LABEL |

⭐ **Tiga judul bagian layar bergerbang `TransferType` 2 / 3 / 4** — layar komite memang
**berganti wajah** menurut arah transfer, dan judulnya ikut berganti. Menguatkan §B1 ronde 2.

### B2 — Kolom "Gerbang" pada tabel 55 properti: **5 baris berubah**

| Properti | Ronde 2 menulis | **Yang benar** |
| --- | --- | --- |
| `.IsSubjectivity` | — | **famili A** `.AcceptStatus = 1 && .TransferType =2` |
| `.SubjectivityNote` | — | **famili A** `.IsSubjectivity = true` |
| `.Adjustment.IsProposeClose` | — | **famili A** `.TransferType =2` |
| `.Adjustment.IsPropReserved` | — | **famili A** `.TransferType =2` |
| `.Adjustment.SwiftCode` | famili B `.TransferType =2` | famili B **+ famili A** `.Adjustment.SwiftCode != ''` |

⚠️ Sepuluh gerbang famili A sisanya menempel pada **label** atau pada properti berawalan
`pyWorkCover.` — **di luar 55** karena definisi ronde 2 hanya menghitung yang berawalan titik (§C).

### B3 — ⛔ RALAT: "enam properti selalu dapat disunting" **SALAH**

Yang benar: **dua**, bukan enam.

| Properti | Dapat disunting | Syarat |
| --- | --- | --- |
| **`.AcceptStatus`** | ✅ **selalu** | — · wajib isi |
| **`.Comment`** | ✅ **selalu** | — · wajib isi |
| `.IsSubjectivity` | ⚠️ **bersyarat** | `.AcceptStatus = 1 && .TransferType =2` |
| `.SubjectivityNote` | ⚠️ **bersyarat** | `.IsSubjectivity = true` |
| `.Adjustment.IsProposeClose` | ⚠️ **bersyarat** | `.TransferType =2` |
| `.Adjustment.IsPropReserved` | ⚠️ **bersyarat** | `.TransferType =2` |

⭐ **Rantai bersyaratnya bertingkat dan masuk akal:** penyetuju hanya melihat kotak *subjectivity*
bila ia **menyetujui** (`AcceptStatus = 1`) **dan** sedang di **jalur adjustment**; catatan
subjectivity baru muncul **setelah kotak itu dicentang**. Dua penanda usul hanya ada di jalur
adjustment.

⚠️ Akibatnya: **di jalur reject dan close, masukan penyetuju tinggal `.AcceptStatus` dan
`.Comment`.** ⛔ Tidak disimpulkan lebih jauh.

### B4 — Sel di balik gerbang mati, dihitung ulang

| Famili | Penanda mati | Sel terdampak |
| --- | --- | --- |
| **B** | `1=2` (1 layout) | **8** — `.ClaimSpreaded` · `.Currency` · `.SharePercentage` ×2 · `.TreatyName` · `.TreatyType` ×2 · *(1 sel lagi)* |
| **A** | `NEVER` (4 sel) | **4** — `pyWorkCover.ClaimData.NoClaim` · `…TotalGrossEstimate…` · `…TotalEstimasi…` · satu LABEL |
| | **TOTAL** | **12 sel** |

⚠️ Ronde 2 melaporkan **8**. Angka yang benar untuk **seluruh** sel adalah **12**. Namun **dalam
definisi 55-properti-berawalan-titik angkanya tetap 8**, karena keempat sel `NEVER` bukan properti
berawalan titik. **Keduanya benar menurut definisinya masing-masing** — itulah sebabnya §C perlu.

⚠️ `[terbuka]` **BARU** — `pyWorkCover.ClaimData.NoClaim` **nomor klaim** duduk di balik `NEVER`.
Nomor klaim tidak tampil di layar komite. Disengaja atau sisa, **tidak ditebak**.

---

## §C — Definisi "kolom layar": **55 atau 117**, bukan 98

### C1 — Dua definisi, dihitung sendiri `[terverifikasi]`

| Definisi | Sel | Properti berbeda |
| --- | --- | --- |
| **(a)** hanya `pyValue` berawalan **titik**, di dalam sel | **85** | **55** |
| **(b)** + seluruh rujukan berawalan **halaman** (`pyWorkCover.` `pyWorkPage.`) di mana pun | +76 kemunculan | **+62** |
| | | **(a)+(b) = 117** |

⚠️ **Angka saya 117, bukan 98.** Selisihnya **definisional, bukan hitungan**: penghitungan saya
ikut memasukkan **halaman**, bukan hanya medan daun — misalnya `pyWorkCover.ClaimData`,
`…EstimationList`, `…PolicyData`, `…InterestList` — dan `pyWorkPage.Edit` yang muncul **8×** sebagai
penanda layar, bukan data. **Saya tulis apa adanya dan tidak saya kejar agar cocok 98.** Siapa pun
yang menghitung ulang harus menyebut saringannya.

⚠️ Catatan yang menjelaskan sebagian selisih: nilai `pyValue` yang **tidak** berawalan titik
ternyata **teks label**, bukan properti — *"Accepted Date"*, *"Claim Amount"*, *"Are you sure to
accept this document?"*. **83 di antaranya.** Jangan dihitung sebagai kolom.

### C2 — Halaman apa `pyWorkCover.` dan `pyWorkPage.` itu `[terverifikasi]`

| Halaman | Isi | Kelas |
| --- | --- | --- |
| **`pyWorkPage`** | **kasus komite itu sendiri** | `ASM-FW-GCNMFW-Work-KomiteTreaty` |
| **`pyWorkCover`** | **kasus klaim induknya** (*cover*) | `ASM-FW-GCNMFW-Work-ClaimTreaty` |

`pyWorkCover` membawa **44** rujukan, hampir seluruhnya `pyWorkCover.ClaimData.*` — kelas
`ASM-FW-GCNMFW-Data-ClaimData`. Itulah **data klaim yang memenuhi layar**: `NoClaim` ·
`InsuredName` · `PolicyNo` · `DateOfLoss` · `TotalEstimasi` · `CauseOfLoss` · `Location` ·
`ReporterName` · `SpreadingClaim` dan seterusnya. `pyWorkPage` membawa **18**, termasuk
`pyWorkPage.Edit` ×8.

⭐ **Layar komite menampilkan dua kasus sekaligus**: kasus komite (`pyWorkPage`) **dan** kasus klaim
induk (`pyWorkCover`). Itu cocok dengan §3.6 ronde 1 yang mencatat `Work-ClaimTreaty` sebagai
*"kasus klaim yang dibuka komite lewat `TempOpenPage`"*.

### C3 — Definisi mana yang dipakai

**Untuk daftar kolom layar dipakai definisi (b), 117** — karena data klaim yang dibaca lewat
`pyWorkCover.` **benar-benar tampil di layar** dan penyetuju membacanya untuk mengambil keputusan.
Definisi (a) hanya mencakup properti kasus komite sendiri, sehingga **menghilangkan seluruh
konteks klaim**. Definisi (a) tetap dicatat karena ia yang menjadi dasar seluruh tabel ronde 2.
⛔ **Ini soal pencatatan, bukan rancangan.**

---

## §D — `SetKomiteList_Act`: ⭐ namanya menyesatkan

### D1 — Identitas `[terverifikasi]`

`pxInsName` = **`ASM-FW-GCNMFW-WORK-KOMITETREATY!SETKOMITELIST_ACT`** · kelas
`ASM-FW-GCNMFW-Work-KomiteTreaty` · `pyActivityType = ACTIVITY`.

### D2 — ⭐ Ia **tidak mengisi `KomiteList` sama sekali**

`[terverifikasi]` Nol penulisan ke `KomiteList` di seluruh langkahnya. Kata `KomiteList` muncul
**5×** di berkasnya, **seluruhnya dibaca, tidak satu pun ditulis.**

**Yang sebenarnya dikerjakannya: menjumlahkan adjustment per mata uang.**

| Step | Tindakan |
| --- | --- |
| **1** | gerbang `T=2 F=6` atas `.TransferType==2` — ⛔ **bila bukan jalur adjustment, `Exit Activity`** |
| 2 · 4 | `Page-Remove` lalu `Page-Copy` — siapkan halaman kerja |
| 3 | `local.TotalAdjRNMinIDR := 0` · `local.TotalAdjGrossinIDR := 0` |
| 5 · 5.1 | kumpulkan mata uang: `TempTotalAdj.pxResults(<APPEND>).Currency := .Currency` |
| 6 | `Java` — *"Hapus currency yg sama"* (buang duplikat) |
| 7 | per mata uang: `local.Currency := .Currency`, nolkan dua penampung |
| **7.1.1** | gerbang **dua baris AND**: `.AcceptanceStatus != 2` **dan** `local.Currency == .Currency` → jumlahkan `GrossAdjustment`, `AdjustmentValue`, dan keduanya **dikali `.KursIDR`** ke IDR |
| 7.2 | `.AdjustmentGross := Local.TotalAdjGross` · `.AdjustmentValue := Local.TotalAdjRNM` |
| 8 | baris penutup `"Total in IDR"` ke `TempTotalAdj` |

⭐ **Tiga akibat:**
1. Layar komite **hanya disiapkan pada jalur adjustment**. Pada reject dan close activity ini
   **keluar di langkah 1** — menguatkan temuan `TransferType = 2` di §B ronde 2 dan §B1 di atas.
2. Baris adjustment **berstatus tolak (`AcceptanceStatus = 2`) tidak ikut dijumlahkan**.
3. Konversi ke IDR memakai **`.KursIDR`** — kurs yang dibekukan, jalur K8. Di sinilah angka itu
   terpakai di layar.

### D3 — Rule lain yang dipanggil

`[terverifikasi]` **Nol `RDB-*`, nol `Call`.** Seluruh langkahnya `Property-Set`, `Page-Remove`,
`Page-Copy`, dan satu `Java`. `pxRuleReferences` hanya menyebut **kelas dan properti**, bukan rule
lain — `Data-Adjustment` (`GrossAdjustment`, `AdjustmentValue`, `KursIDR`, `Currency`,
`AcceptanceStatus`, `AdjustmentGross`), `Data-ClaimData` (`AdjustmentList`), `Work-ClaimTreaty`
(`ClaimData`), `Work-KomiteTreaty` (`TransferType`), `Code-Pega-List` (`pxResults`).

⚠️ Langkah **6** memakai **`Java`** — satu-satunya di activity ini. Isinya **tidak dibaca** di ronde
ini. `[terbuka]` **BARU**.

### D4 — Silang empat kolom grid komite

| Kolom grid | Diisi `SetKomiteList_Act`? | Penulis sebenarnya di modul ini |
| --- | --- | --- |
| `.IDKomite` | ⛔ **tidak** | ⛔ **NOL penulis di modul ini** — hanya dibaca |
| `.KomiteAproval` | ⛔ **tidak** | `KomitePostAdjustment` · `KomitePost_Close` · `KomitePost_Reject` |
| `.DateApprove` | ⛔ **tidak** | *(idem)* |
| `.KomiteComment` | ⛔ **tidak** | *(idem)* |

⭐ **Keempatnya nol.** Grid penyetuju **tidak disiapkan** oleh pra-proses layar; isinya sudah ada
di kasus sejak dibuat, dan hanya **ditimpa saat penyetuju menekan Submit**. `.IDKomite` bahkan
**tidak pernah ditulis di modul ini sama sekali** — menguatkan §9.2 ronde 1 bahwa daftar penyetuju
dibentuk **saat kasus komite dibuat**, di luar modul ini.

---

## §E — `ASM-FW-GCNMFW-Data-Comitee`

### E1 — Sensus pemakaian `[terverifikasi]`

**7 berkas** modul ini menyebutnya:

| Berkas | Kemunculan |
| --- | --- |
| `Section/ShowTransfer.xml` | 30 |
| `Activity/KomitePostAdjustment.xml` | 18 |
| `Activity/KomitePost_Close.xml` | 11 |
| `Activity/KomitePost_Reject.xml` | 11 |
| `Activity/KomiteRouter.xml` | 8 |
| `Activity/SendEmailKlaim_KMT.xml` | 5 |
| `Activity/PrintFileAcceptance_TKMT.xml` | 4 |

Dipakai sebagai **kelas baris daftar penyetuju** — `pyWorkPage.KomiteList(…)` dan
`TempOpenPage.ClaimData.ClaimComitee(…)` keduanya berkelas ini. ⛔ **Tidak disimpulkan pewarisannya.**

⚠️ Ia memang **belum tercatat** di peta kelas §3.6 ronde 1, yang hanya menyebut lima kelas.
**Ralat dicatat di sini, berkas ronde 1 tidak disunting.**

### E2 — F6: apa yang perlu diminta ke work owner

**Satu baris:** ⛔ **Minta ekspor rule `Rule-Obj-Class` untuk lima kelas** —
`ASM-FW-GCNMFW-Work-KomiteTreaty` · `-Work-ClaimTreaty` · `-Data-Adjustment` · `-Data-ClaimData` ·
`-Data-Comitee` — **atau tangkapan layar tab *Class* masing-masing di Pega**, yang menampilkan
**kelas induk** dan **daftar properti yang dideklarasikan di sana**.

---

## §F — Apa lagi yang seharusnya dikerjakan

### F1 — Pekerjaan sebelum spec, urut dari yang paling menahan

| # | Yang dikerjakan | Kenapa perlu | Besar | Menunggu |
| --- | --- | --- | --- | --- |
| **1** | **Sisir `pyStepsTransParams`** — 552 baris berkode yang belum pernah dibaca, termasuk 6 baris berkode 1 | ⛔ **paling menahan**: ia keluarga gerbang kedua. Selama belum dibaca, **setiap kesimpulan alur langkah berpotensi tidak lengkap** — termasuk K1 dan K2 | sedang | korpus |
| **2** | **Efek keluar** — `SendAcceptationToKasir` (REST) · `HTMLToPDF` · `InsertDocument_Act` · `SendEmailWithAttachments` · `GetLinkService` · SystemSettings `LinkService` | apa yang **keluar dari sistem** saat komite menyetujui. Di modul lain ini selalu satu tiket sendiri. Enam berkas **belum pernah dibuka** | berat | korpus |
| **3** | **Baca berkas `Flow`** — daur hidup kasus komite | belum pernah dianalisis; ia yang menentukan **urutan tahap**, bukan urutan step | sekali sisir | korpus |
| **4** | **`ViewDetailInterest`** — Section 115 KB + FlowAction | **layar kedua** modul ini, belum pernah dibuka sama sekali | sedang | korpus |
| **5** | **Delapan activity penghitung + jalur konversi** | mengisi `.EstimationValue`, `.SharePercentage`, `.ClaimSpreaded` di layar | sedang | korpus, lalu work owner *(lingkup)* |
| **6** | **Sisa 8 RDBList · 3 ReportDefinition · 2 DataTransform · 1 DecisionTable · 2 When** | melengkapi peta; `IsPEGAPROD` diduga penanda lingkungan | sekali sisir | korpus |

⚠️ Butir **1** naik ke puncak **karena temuan §A5 GOYAH 1**, bukan karena rencana sebelumnya.

### F2 — Bagian yang **masih gelap** — **5**

1. **`pyStepsTransParams`** — 552 baris, keluarga gerbang kedua. Belum masuk daftar terbuka mana pun.
2. **Berkas `Flow`** — daur hidup kasus. Disebut di inventaris §1 ronde 1, **tidak pernah dibaca**.
3. **`pyActionSets` / `pyBehaviors` di `ShowTransfer`** — aksi sisi-klien pada tombol. `.pyTemplateButton` muncul 2× dan **tidak pernah ditelusuri** ke mana tombolnya menuju.
4. **`pyWorkPage.Edit`** — muncul **8×** di layar sebagai penanda, tidak pernah diperiksa apa yang menyetelnya.
5. **Langkah `Java` di `SetKomiteList_Act` step 6** — satu-satunya Java di activity itu, isinya belum dibaca. Java di Pega bisa melakukan apa saja.

### F3 — Kesimpulan yang **paling rawan salah** — **3**

1. ⛔ **Ronde 2 §E "566 baris syarat"** — **paling rawan**, dan sudah terbukti tidak lengkap (§A5 GOYAH 1). Setiap angka turunannya ikut rawan.
2. ⚠️ **Ronde 2 §B "85 sel / 55 properti"** — rawan **karena definisinya**, bukan karena hitungannya. §C menunjukkan angka sah lain adalah 117. Setiap kalimat ronde 2 yang berbunyi *"tidak ada di layar"* perlu dibaca ulang dengan definisi (b) — termasuk **`TransferType` "tidak tampil"**, yang ternyata **muncul sebagai syarat di 10 gerbang**.
3. ⚠️ **Ronde 1 §3.7 "9 rule ber-`COMMIT` sendiri"** — diturunkan dari penyisiran 23 RDBList dengan pencocokan teks, bukan parser. Bila ada `COMMIT` yang tertulis berbeda ejaannya, ia terlewat.

⛔ **Ketiganya tidak saya perbaiki di sini** — hanya ditunjuk, sesuai perintah.

### F4 — Pertanyaan untuk work owner — **4 baru**

> Dua pertanyaan yang sudah naik dan belum dijawab — **14 properti tanpa penulis**, dan
> **`.TreatyType` di balik gerbang mati** — **dibiarkan terbuka**, tidak diulang di sini.

#### Pertanyaan 1 — Nomor klaim sengaja disembunyikan dari layar komite?

**Apa yang ditanyakan.** Di layar komite, kotak nomor klaim dipasang tetapi diberi penanda yang
membuatnya **tidak pernah muncul**. Apakah itu memang diinginkan, atau sisa dari perubahan lama.

**Kenapa muncul.** Tiga kotak berisi data klaim — nomor klaim dan dua kotak total estimasi — sama-
sama dimatikan dengan cara yang sama.

**Bedanya jawaban A atau B.** Bila **disengaja**, layar baru dibuat tanpa ketiga kotak itu, dan
penyetuju memang tidak perlu melihat nomor klaim. Bila **sisa**, ketiganya perlu dinilai ulang
satu per satu — mungkin justru harus tampil.

**Yang tertahan.** Daftar kolom layar komite belum bisa difinalkan.

#### Pertanyaan 2 — Pada jalur tolak dan tutup, penyetuju memang hanya boleh mengisi dua hal?

**Apa yang ditanyakan.** Di jalur **adjustment**, penyetuju dapat mengisi enam hal. Di jalur
**tolak** dan **tutup**, hanya **dua**: keputusan setuju/tolak, dan catatan. Apakah itu memang
aturan bisnisnya.

**Kenapa muncul.** Empat isian lainnya bergerbang `TransferType = 2`, dan pra-proses layar bahkan
**berhenti di langkah pertama** bila bukan jalur adjustment.

**Bedanya jawaban A atau B.** Bila **benar**, layar baru dibuat tiga wajah dan sisanya
disederhanakan. Bila **tidak**, ada kebutuhan yang selama ini tidak terlayani, dan itu **perubahan
perilaku** — bukan sekadar pemindahan.

**Yang tertahan.** Rancangan layar komite untuk jalur tolak dan tutup.

#### Pertanyaan 3 — Tiga sasaran lompatan yang tidak punya tujuan: dibiarkan?

**Apa yang ditanyakan.** Tiga langkah menyimpan perintah "lompat ke tanda X", tetapi **tanda itu
tidak ada** di activity-nya. Ketiganya kebetulan sedang dimatikan, jadi tidak pernah berjalan.
Apakah dibiarkan apa adanya, atau ditandai sebagai cacat yang perlu diperbaiki lebih dulu.

**Kenapa muncul.** Preseden "tiru apa adanya" sudah dipakai untuk sembilan rule ber-`COMMIT`
sendiri. Pertanyaannya apakah preseden itu juga berlaku untuk lompatan yang rusak.

**Bedanya jawaban A atau B.** Bila **dibiarkan**, ketiganya dicatat sebagai jalur mati dan tidak
dipindahkan. Bila **cacat**, perlu dipastikan dulu ke mana seharusnya melompat — dan itu **hanya
bisa dijawab orang yang menulisnya**, karena tidak ada jejaknya di korpus.

**Yang tertahan.** Tidak ada yang tertahan sekarang; tertahan bila kelak gerbangnya dihidupkan.

#### Pertanyaan 4 — Data klaim di layar komite: ikut dipindahkan, atau dibaca dari modul klaim?

**Apa yang ditanyakan.** Separuh lebih isi layar komite bukan milik kasus komite, melainkan
**milik kasus klaim induknya** — nama tertanggung, nomor polis, tanggal kejadian, total estimasi,
dan sekitar empat puluh hal lain. Apakah itu bagian garapan modul ini, atau cukup dibaca dari
modul klaim.

**Kenapa muncul.** Lingkup sudah diputuskan "hanya Komite Claim Prop". Tetapi layar modul ini
**tidak bisa ditampilkan** tanpa data klaim.

**Bedanya jawaban A atau B.** Bila **ikut**, jumlah kolom layar yang digarap naik dari sekitar 55
menjadi sekitar 117, dan modul ini bergantung pada bentuk data klaim. Bila **dibaca saja**, modul
ini hanya menggarap kolomnya sendiri, tetapi **harus menunggu modul klaim selesai lebih dulu**.

**Yang tertahan.** Besar-kecilnya modul ini, dan urutan pengerjaan terhadap Claim Prop.

---

## §G — Daftar `[terbuka]` modul ini sesudah ronde 3 — **10**

| # | Butir | Menunggu |
| --- | --- | --- |
| 1 | `OPERATORID` `HISTORYAKSEPTASIPEGA` tak pernah diisi jalur komite | korpus / DBA |
| 2 | **F6** — pewarisan kelas kerja *(permintaan ke work owner sudah dirumuskan, §E2)* | work owner |
| 3 | 14 properti layar tanpa penulis di modul ini *(sudah naik, belum dijawab)* | work owner |
| 4 | beda `CONDITION` vs `ExpressionCondition` tidak terbaca dari struktur | korpus |
| 5 | `.TreatyType` hanya di balik gerbang mati `1=2` *(sudah naik, belum dijawab)* | work owner |
| 6 | ⭐ **BARU** — `pyStepsTransParams`, 552 baris berkode belum dibaca | korpus |
| 7 | ⭐ **BARU** — `pyWorkCover.ClaimData.NoClaim` di balik `NEVER` | work owner *(pertanyaan 1)* |
| 8 | ⭐ **BARU** — isi langkah `Java` di `SetKomiteList_Act` step 6 | korpus |
| 9 | ⭐ **BARU** — tiga sasaran lompatan menggantung (`JMP`, `EXT` ×2) | work owner *(pertanyaan 3)* |
| 10 | ⭐ **BARU** — definisi kolom layar 55 atau 117 | work owner *(pertanyaan 4)* |

**K10 keluar dari daftar — `[tertutup]` (§A4).**

---

## §H — Ralat ke ronde sebelumnya, dikumpulkan

| # | Di mana | Yang salah | Yang benar |
| --- | --- | --- | --- |
| 1 | ronde 2 §A2 Aturan 1 | *"sel tidak pernah bergerbang sendiri"* | **famili A, 15 gerbang di `pyUserData`** (§B1) |
| 2 | ronde 2 §B1 | *"enam properti selalu dapat disunting"* | **dua** (§B3) |
| 3 | ronde 2 §B1 | 8 sel di balik gerbang mati | **12 sel** seluruhnya; 8 tetap benar dalam definisi (a) (§B4) |
| 4 | ronde 2 §E | *"566 baris syarat"* | **+552 baris `pyStepsTransParams`** belum terhitung (§A5) |
| 5 | ronde 2 §C | *"TParam kosong"* | sasaran ada di **`…WhenTruePrms`** (§A1) |
| 6 | ronde 1 §8 | `(kosong)/3`, `(kosong)/1` *"belum punya arti"* | terbaca terang dengan tabel A1 (§A5) |
| 7 | ronde 1 §8 · §9.5 | `2/2` *"step tanpa gerbang"* | **Continue Whens/Continue Whens** — hasil sama, mekanisme beda (§A5) |
| 8 | ronde 1 §3.6 | peta kelas menyebut lima kelas | **`Data-Comitee` belum tercatat**, dipakai 7 berkas (§E1) |
| 9 | ronde 2 §G | *(ralat nama `INSERTHISTORYAKSEPTASIPEGA_SQL`)* | tetap berlaku, tidak berubah |

⛔ **Tidak satu pun berkas ronde 1 atau ronde 2 disunting.**
