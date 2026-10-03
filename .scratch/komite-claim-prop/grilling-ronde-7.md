# Grilling Ronde 7 — Komite Claim Prop

**Tanggal:** 2026-09-18 · **Korpus:** `D:\XML\RNM_BRD\Komite Claim Prop\` (READ-ONLY)
**Lingkup:** Komite Claim Prop saja `[keputusan work owner]` 2026-09-18.
**Metode:** parser XML bersarang. Bukti = path berkas + nama rule + nomor step Pega.
**Tidak ada nomor baris XML, dan tidak ada kode Java, di berkas ini.**

> **Ralat dicatat DI SINI**, berkas ronde 1–6 tidak disunting.

---

## 0. Ringkasan

| # | Hasil | Status |
| --- | --- | --- |
| **§A1** | `komiteAccept_ticket` | ✅ **dibuang** — jalur mati |
| **§A2** | Penguncian kasus klaim | ⛔ **RALAT saya sendiri** — `ReleaseOnCommit = true`, cocok |
| **§A3** | `pxCreateOpName` · `pxCreateDateTime` | ⚠️ **DATA NYATA** — harus ikut ke spec |
| **§B** | 20 langkah `Java` | **0 sentuh basis data** · **1 layanan luar** · **7 tangani gagal** |
| **§B4** | Ronde 5 §B4 | ⭐ **GUGUR SEBAGIAN** — PDF **ada** penanganannya |
| **§C** | `pyStepsRepeatDef` | **85** — `EMBEDDED` 65 · `REPEAT` 20 |
| **§C5** | Perulangan yang memuat penulisan | ⚠️ **15 langkah** |
| **§D** | Aksi lokal · `pyXMLSignature` · F6 · `CONDITION` | 3 bawaan Pega · **6 parameter terlewat** · F6 buntu |
| **§E2** | **Cukup untuk menulis spec?** | ⭐ **YA** — dengan dua bab ditandai `[terbuka]` |

---

## §A — Satu keputusan, satu ralat, satu peringatan

### A1 — `komiteAccept_ticket` **DIBUANG** `[keputusan work owner]` 2026-09-18

Ticket itu **sudah tidak dipakai**. **Jalur mati — dibuang, tidak dipindahkan, dan tidak
ditelusuri siapa pemicunya.** Perlakuannya **sama dengan lima lompatan menggantung** (§A3 ronde 5).

✅ Butir `[terbuka]` no.8 **TUTUP**.

### A2 — ⛔ RALAT: penguncian kasus klaim **jauh lebih kecil** dari yang saya tulis

> **Pertanyaan 13 saya sendiri salah premisnya.** Ronde 6 §B3 menulis bahwa ketiga `KomitePost*`
> mengunci kasus klaim induk, dan saya menaikkannya sebagai *"selama komite bekerja, klaim itu
> tidak bisa disentuh proses lain"*. **Kalimat itu DICABUT.**

`[terverifikasi]` **Verifikasi sendiri, ketiga berkas:**

| Berkas | Step | `InstanceHandle` | `Lock` | `ReleaseOnCommit` |
| --- | --- | --- | --- | --- |
| `KomitePostAdjustment` | 4 | `pyWorkPage.pxCoverInsKey` | `true` | **`true`** |
| `KomitePost_Close` | 1 | `pyWorkPage.pxCoverInsKey` | `true` | **`true`** |
| `KomitePost_Reject` | 1 | `pyWorkPage.pxCoverInsKey` | `true` | **`true`** |

✅ **COCOK PERSIS di ketiganya.**

⭐ **`ReleaseOnCommit = true` berarti kunci DILEPAS saat `COMMIT`.** Karena `COMMIT` ada di langkah
41 activity yang sama, kunci itu **dipasang saat penyetuju menekan simpan dan dilepas beberapa
detik kemudian** — **bukan** dipegang selama komite menimbang. Penyetuju bisa membuka layar
berjam-jam tanpa mengunci apa pun; penguncian hanya menyelimuti **proses simpan**.

✅ Butir `[terbuka]` no.10 **TUTUP sebagai perilaku biasa**, bukan keputusan.
⛔ **Pertanyaan 13 ditarik.**

### A3 — ⚠️ PERINGATAN: dua yang dibuang dari hitungan kolom **BUKAN sampah**

> ⭐ **`pyWorkPage.pxCreateOpName` ("dibuat oleh") dan `pyWorkPage.pxCreateDateTime` ("tanggal
> dibuat") adalah DATA NYATA yang TAMPIL di layar komite, dan HARUS ikut terbawa ke spec.**
>
> Keduanya berada **di luar hitungan 93 kolom** hanya karena aturan pencatatan §C ronde 6
> memisahkan properti bawaan Pega. **Pemisahan itu soal cara menghitung, bukan soal nilai.**
> Penyetuju melihat keduanya di layar.

⛔ **Nama kolomnya di sistem baru tidak diusulkan di sini.**

⚠️ Tiga sisanya yang dibuang memang benar-benar kerangka: `.pyTemplateInputBox` dan
`.pyTemplateButton` ×2 — kotak kosong dan tombol, bukan data.

---

## §B — Dua puluh langkah `Java`

### B1 · B2 · B3 — Daftar lengkap `[terverifikasi]`

⛔ **Kode Java tidak disalin.** Yang ditulis hanya perilakunya.

**Kelompok 1 — pembersih duplikat (8 langkah).** Pola identik: menelusuri sebuah daftar di
halaman kerja, membandingkan isi kolom mata uang antar baris, lalu **membuang baris yang nilainya
sudah pernah muncul**. Tidak membaca basis data, tidak memanggil apa pun.

| Activity · step | Daftar yang dibersihkan |
| --- | --- |
| `AddEstimation_Act` **3** | `LossAllocation` — kunci `.Currency` `.CurrencyID` `.TreatyType` |
| `AddLossAllocation_act` **3** | `UangList` — kunci `.Currency` `.CurrencyID` `.IDR` |
| `CountEstimation_Act` **5** | `Money` |
| `CountEstimation_Act` **11** | `ClaimData.ListTotalEstimation` |
| `CountPersen_act` **3** | `ListUang` |
| `CountSpreading_act` **4** | `Estimate` — kunci termasuk `.TypeLoss` |
| `SetKomiteList_Act` **6** | `TempTotalAdj` *(sudah dibaca §C2 ronde 4)* |

**Kelompok 2 — perakit pesan ke Kasir (2 langkah).**
`HitServiceToKasirKMT_Act` **10.6** dan **14.3**, masing-masing ±2,4 KB — **yang terbesar di
modul ini**. Keduanya **merakit teks pesan** dari puluhan potongan (`CompanyName`, `AccountNo`,
`Email`, `Kepada`, `LbgID`, `LbuId`, `LdcId`, `AcceptType`, dan seterusnya) lalu menaruhnya ke
halaman utama. ⭐ **Keduanya dibungkus `try`/`catch` dan mencatat galat ke log** bila perakitan
gagal.

**Kelompok 3 — rantai PDF (7 langkah).** ⭐ **Inilah yang membalik ronde 5:**

| Activity · step | Perilaku | Tangani gagal |
| --- | --- | --- |
| `HTMLToPDF` **1** | membuang sisa nilai `PDFDocument` dari pemanggilan sebelumnya | — |
| `HTMLToPDF` **3** | bila markup kosong, isi dari rule stream bernama parameter `StreamName` | — |
| `HTMLToPDF` **5** | membuang rujukan gaya (css) bila markup sudah membawanya sendiri | — |
| `HTMLToPDF` **6** | menyalin markup ke variabel lokal · ⭐ **melempar galat bila kosong** | ✅ `throw` |
| `HTMLToPDF` **9** | membuat berkas PDF lewat API Pega · ⭐ **mencatat galat dan berhenti** bila API tidak mengembalikan isi | ✅ `oLog.error` + `return` |
| `KomitePost_Close` **12.6** · `KomitePost_Reject` **12.6** · `PrintFileAcceptance_TKMT` **12** | membaca berkas PDF dan mengubahnya jadi teks Base64 untuk dilampirkan · ⭐ `try`/`catch`, **melempar galat** *"Can't attach the file to the Work"* | ✅ `throw` + `catch` |

**Kelompok 4 — pengambil lampiran (1 langkah).**
`GetBase64Attachment` **5.3** — ⭐ **satu-satunya Java yang memanggil layanan luar**: ia
**membuka aliran dari sebuah URL**, membacanya, dan mengubahnya jadi Base64. Dibungkus
`try`/`catch` dan mencatat *"Gagal ambil file dari URL"*.

**Kelompok 5 — penentu jenis berkas (2 langkah).**
`InsertDocument_Act` **2** dan `InsertGoogleStorage_Act` **2** — mengambil akhiran nama berkas
sesudah titik terakhir, dipakai bila jenis berkas tidak dikirim. Tidak menyentuh apa pun.

### Jawaban B3, diringkas `[terverifikasi]`

| | Jumlah |
| --- | --- |
| Langkah `Java` seluruhnya | **20** |
| **Menyentuh basis data** | ⭐ **0** — nol tabel, nol `insert`/`update`/`select` |
| **Memanggil layanan luar** | **1** — `GetBase64Attachment` 5.3, membuka aliran dari URL |
| **Menangani kegagalan** | **7** — `GetBase64Attachment` 5.3 · `HTMLToPDF` 6 · `HTMLToPDF` 9 · `HitServiceToKasirKMT_Act` 10.6 · 14.3 · `KomitePost_Close` 12.6 · `KomitePost_Reject` 12.6 · `PrintFileAcceptance_TKMT` 12 |

### B4 — ⭐ UJI ULANG ronde 5 §B4: **GUGUR SEBAGIAN**

Ronde 5 menulis: *"PDF, unggah berkas, dan email NOL penanganan gagal."*

| Efek keluar | Ronde 5 | **Sesudah 20 Java dibaca** |
| --- | --- | --- |
| **PDF** | *"nol penanganan"* | ⛔ **GUGUR.** Ada **empat** titik penanganan: `HTMLToPDF` **6** melempar galat bila markup kosong; `HTMLToPDF` **9** mencatat galat dan berhenti bila API tidak menghasilkan isi; dan **tiga** langkah Base64 melempar galat *"Can't attach the file to the Work"* bila lampiran gagal |
| **Unggah Google Storage** | *"nol penanganan"* | ✅ **MASIH BERDIRI.** `InsertGoogleStorage_Act` tidak punya `Java` penangan, tidak punya transisi berkode khusus, dan langkah `Connect-REST`-nya polos |
| **Email** | *"nol penanganan"* | ✅ **MASIH BERDIRI.** `SendEmailKlaim_KMT` dan `SendEmailWithAttachments` sama-sama nihil |

⭐ **Penanganan PDF berbentuk melempar galat, bukan memulihkan.** Bila PDF gagal, activity
**berhenti dengan galat** — berbeda dari pola Kasir/arasapas yang **melompat dan melanjutkan**.
Jadi modul ini punya **dua sikap berbeda terhadap kegagalan**, dan keduanya disengaja.

⚠️ `[terbuka]` sisa: unggah berkas dan email **memang tanpa penanganan**. Itu kini **kesimpulan
dari empat keluarga wadah yang seluruhnya sudah dibaca** — jauh lebih kuat dari ronde 5.

### B5 — `pyMemo` dan `pyUsage` `[terverifikasi]`

| Wadah | Muncul | **Berisi** |
| --- | --- | --- |
| `pyMemo` | 73 | **73** — terisi seluruhnya |
| `pyUsage` | 298 | **226** |

⭐ **`pyMemo` ternyata catatan perubahan terakhir penulis rule**, dan **tujuh di antaranya
mengubah atau menguatkan pembacaan:**

| Berkas | `pyMemo` | Yang dikuatkan |
| --- | --- | --- |
| `KomitePostAdjustment` | **"add StepStatusFail"** | ⭐ keluarga transisi **sengaja ditambahkan** untuk menangani gagal — menguatkan §A2 ronde 4 |
| `HitServiceToKasirKMT_Act` | **"add StepStatusFail"** | idem |
| `InsertGoogleStorage_Act` | **"FIX ERROR HANDLING"** | ⚠️ **menantang B4**: ada pekerjaan penanganan galat, tetapi **jejaknya tidak terbaca** di empat keluarga wadah |
| `GetUrlGoogleStorage_Act` | **"FIX ERROR HANDLING"** | sejalan — dan di sini jejaknya **ada** (`6/2` Exit Activity) |
| `SetKomiteList_Act` | **"skip selain .TransferType =2"** | menguatkan §D ronde 5 |
| `CurrencyStandard` | **"change browse query"** | menguatkan §B ronde 6 |
| `KomiteTreaty_Flow` | **"Resolved-Completed"** | menguatkan §D2 ronde 6 |

⚠️ `[terbuka]` **BARU** — `InsertGoogleStorage_Act` bercatat *"FIX ERROR HANDLING"* tetapi
**tidak terlihat penanganan apa pun** di keempat keluarga wadah. Entah perbaikannya ada di tempat
lain, entah catatannya menyusul perubahan yang kemudian dicabut. **Tidak ditebak.**

⛔ `pyUsage` **tidak mengubah pembacaan apa pun** — isinya keterangan baku Pega tentang metode,
bukan catatan proyek ini.

---

## §C — Keluarga ketiga `pyStepsRepeatDef`

### C1 — Sensus 100 % `[terverifikasi]`

**85 langkah berulang** di **21 activity** *(dari 552 langkah, 35 activity)*.

| `HasRepeat` | Jumlah |
| --- | --- |
| **`EMBEDDED`** | **65** |
| **`REPEAT`** | **20** |

Terbanyak: `CountEstimation_Act` **18** · `CountSpreading_act` **10** · `HitServiceToKasirKMT_Act`
**9** · `CountPersen_act` **8** · `AddEstimation_Act` **6**. Di `KomitePostAdjustment` hanya **2**.

Tag yang pernah terisi di dalamnya: `…HasRepeat` **85** · `…Iteration` · `…Start` · `…Limit`
masing-masing **21**.

### C2 — Beda `REPEAT` dan `EMBEDDED`, diturunkan dari isi `[terverifikasi]`

⭐ **Satu pola bersih terbaca:**

| | `REPEAT` | `EMBEDDED` |
| --- | --- | --- |
| Jumlah | 20 | 65 |
| Membawa `Iteration` / `Start` / `Limit` | **SELALU**, dan **selalu bernilai `1` / `1` / `1`** | **tidak**, kecuali **satu** |
| Keterangan langkahnya | *"Insert ke table"*, *"Untuk Prop dan Facin"*, *"Set NET untuk CLMP"* | *"looping estimasi"*, *"looping temporary money"*, *"looping Spreading"* |

⭐ Satu-satunya `EMBEDDED` yang membawa `1/1/1` adalah **`KomitePostAdjustment` step 26**.

⚠️ `[terbuka]` **Arti bedanya tidak terbaca dari struktur.** `pyStepsRepeatType` **kosong di
seluruh 552 langkah**, dan **tidak ada satu pun elemen yang menyebut halaman mana yang diulang**.
Pembacaan *"`EMBEDDED` mengulang daftar, `REPEAT` mengulang sekali"* **konsisten dengan keterangan
langkahnya**, tetapi **keterangan bukan struktur**. ⛔ **Tidak ditebak.**

### C3 — Silang dengan kode arah 4 · **langkah 26 TERBUKTI `EMBEDDED`** `[terverifikasi]`

| Baris | Activity · step | Jenis perulangan | Cocok? |
| --- | --- | --- | --- |
| `_`/`4` keluarga 1 | `KomitePostAdjustment` **26** | ⭐ **`EMBEDDED`** (`1`/`1`/`1`) | ✅ **TERBUKTI** — klaim ronde 4 benar |
| `2`/`6` · `6`/`2` · `6`/`6` keluarga 2 | `HitServiceToKasirKMT_Act` **3** · `GetUrlGoogleStorage_Act` **4** · `KomiteRouter` **6.1** | ⛔ **bukan langkah berulang** | konsisten — *Exit Activity* memang tidak butuh perulangan |

⭐ **Hanya kode 4 (*Exit Iteration*) yang duduk di langkah berulang**, dan itu **satu-satunya**
kode yang memang **memerlukan** perulangan. Tidak ada `Exit Iteration` yang tersesat di luar loop.

⭐ Tambahan: `KomitePostAdjustment` step **16** — induk blok penomoran akseptasi — ternyata
**`REPEAT` (`1`/`1`/`1`)**, menguatkan §B3 ronde 4.

### C4 — Tangga penyetuju: ⭐ **bukan dijalankan perulangan langkah**

`[terverifikasi]` `KomiteRouter` punya **satu** langkah berulang (**step 6**, `EMBEDDED`) — dan
itu **induk step 6.1** yang menetapkan `param.AssignTo`. Ia mengulang **daftar penyetuju untuk
memilih yang giliran**, bukan menjalankan tangganya.

✅ **Menguatkan §D5 ronde 6**: tangga penyetuju berputar lewat **jalur balik di Flow**
(`Decision1` → `ASSIGNMENT63`), **bukan** lewat perulangan di dalam activity. Perulangan langkah
hanya dipakai untuk **menyisir daftar**, bukan untuk **menunggu penyetuju**.

Perulangan di `KomitePostAdjustment` step **26** (`EMBEDDED`, anak **26.1** *"auto reject all
komite"*) menyisir `KomiteList` untuk **menolak otomatis sisa penyetuju** — itu **satu putaran di
dalam satu kunjungan**, bukan tangga.

### C5 — ⚠️ Perulangan yang **memuat penulisan atau efek keluar**: **15 langkah**

⛔ **Langkah berulangnya sendiri tidak pernah menulis** — penulisan ada pada **anaknya**. Itu
penting: **berapa kali penulisan terjadi = berapa kali induknya berputar.**

| Activity · langkah berulang | Yang ditulis / dikeluarkan di dalamnya |
| --- | --- |
| **`KomitePostAdjustment` 16** | 4 `RDB-List`: `GETTanggalClosing_SQL` *(remark)* · `GenerateNoAcceptTreaty` *(remark)* · `GetKodeProdNonLife_SQL` · `GetSequenceNumber_SQL` |
| **`KomitePost_Close` 11** | `GCNM!SaveDataToOSAkseptasi` |
| **`KomitePost_Close` 13** · **`KomitePost_Reject` 13** | `GCNM!InsertClaimRejected` — ⚠️ **`POOLDATA.CLAIMREJECTED`**, salah satu dari sembilan rule ber-`COMMIT` sendiri (K7) |
| **`HitServiceToKasirKMT_Act` 10** · **14** | `GetEmailCeding_SQL` · ⭐ **`Connect-REST` ke Kasir** · `InsertLOGDirectKasir_SQL` |
| **`HitServiceToKasirKMT_Act` 14.1 · 14.2 · 14.2.2 · 14.2.2.1 · 14.2.2.1.1** | `GetEmailCeding_SQL` — bersarang lima tingkat |
| **`InsertGoogleStorage_Act` 12** | `GenerateImageID_SQL` · `Insert_T_Storage_SQL` — ⚠️ **`T_STORAGE_IMAGE`** |
| **`GetUrlGoogleStorage_Act` 6 · 6.6** | `GetTokenStorage_SQL` · `Connect-REST` · `Update_T_Storage_SQL` |
| **`AddEstimation_Act` 5** | `CurrencyStandard` *(baca)* |
| **`SaveAcceptationTreaty_TKMT` 4** | `GetLimit…` · `GetList…` *(baca)* |

⭐ **Dua yang paling perlu diperhatikan:** panggilan **REST ke Kasir** dan **penulisan
`CLAIMREJECTED`** keduanya **berada di dalam langkah berulang**. Berapa kali keduanya terjadi
**ditentukan oleh berapa kali induknya berputar** — dan itu **belum terbaca**, karena §C2 tidak
menemukan halaman yang diulang. ⚠️ `[terbuka]`.

---

## §D — Aksi lokal, parameter tersembunyi, dan dua butir lama

### D1 — Tiga aksi lokal `KomiteRouter` `[terverifikasi]`

Diambil dari bentuk `ASSIGNMENT63` di berkas `Flow` (`pyLocalActionsString = (3)`):

| Nama | Asal | Siapa yang menjalankan |
| --- | --- | --- |
| **`pyTransferAssignment`** | ⭐ **bawaan Pega** | pengguna yang memegang assignment — memindahkan tugas ke orang lain |
| **`pyCreateAdhocCase`** | ⭐ **bawaan Pega** | idem — membuat kasus sampingan |
| **`EngageExternal`** | ⭐ **bawaan Pega** | idem — melibatkan pihak luar |

`[terverifikasi]` Ketiganya terdaftar di `pxRuleReferences` berkelas
`ASM-FW-GCNMFW-Work-KomiteTreaty` sebagai `Rule-Obj-FlowAction`, **bersama** `VIEWTRANSFERDTL`.

⭐ **Ketiganya bawaan Pega, bukan buatan Nusantara Re** — tidak ada logika modul ini di dalamnya.
⚠️ Siapa yang **berhak** memakainya tidak terbaca: `pyPrivilegeList` milik Flow memuat **kelas
tanpa nama privilege**, pola yang sama dengan §A2 ronde 6.

### D2 — `pyXMLSignature`: **6 parameter belum pernah tercatat** `[terverifikasi]`

**18 dari 35** activity mendeklarasikan parameter di sana; **50 parameter** seluruhnya. Disilang
dengan seluruh teks ronde 1–6:

| Activity | Parameter yang **terlewat** |
| --- | --- |
| `CheckEstimateDate_Act` | `EstimateDate` |
| `CurencyEstimation_Act` | `CurrID` |
| `SetCurencyList_act` | `CurrID` |
| `HTMLToPDF` | `StreamName` · `PDFDocument` · `UseCompactStylesforPDF` |

⭐ **Hanya 6 dari 50** — 44 sisanya sudah tercatat lewat jalur lain. Tiga di antaranya milik
`HTMLToPDF`, rule **bawaan Pega**, jadi **tidak perlu dipindahkan**. Yang benar-benar baru untuk
modul ini: **`EstimateDate`** dan **`CurrID`** (2 activity).

### D3 — F6: ⛔ **tidak ada di jendela ini** — dicek habis

`[terverifikasi]` Seluruh wadah teks di **80 berkas** modul ini disisir untuk penyebutan kelas
induk (`inherit` · `parent` · `induk` · `ASM-FW-GCNMFW-Work` · `@baseclass`) di `pyMemo`,
`pyUsage`, `pyXMLSignature`, `pyNotes`, `pyDescription`.

**Hasil: satu penyebutan**, dan itu **bukan tentang pewarisan** — `SendEmailWithAttachments`
`pyUsage` berbunyi *"This activity basically gathers all the parameters for email and passes…"*,
keterangan baku Pega.

⛔ **F6 berhenti di sini.** **"Tidak ada di jendela ini"** — dan sekarang itu **bukan lagi
kesimpulan dari ketiadaan berkas saja**, melainkan dari **penyisiran seluruh wadah teks**.
⛔ **Tidak disimpulkan dari pemakaian.** Permintaan ke work owner tetap seperti §E2 ronde 3.

### D4 — `CONDITION` vs `ExpressionCondition`: **terbaca sebagian, artinya tidak**

`[terverifikasi]` Bedanya **ada di struktur**, dan konsisten:

| Elemen | `CONDITION` (7) | `ExpressionCondition` (2) |
| --- | --- | --- |
| `pyIsBodyVisibilityOption` | **`ALWAYS`**, 7 dari 7 | ⛔ **tidak ada sama sekali** |
| `pyIsClientWhen` | `true` 6 · `false` 1 | **`false`**, 2 dari 2 |
| `pyIsClientActiveWhen` | `false` | `false` |
| `pyContainerVisibleWhen` | `.TransferType =2` ×5 · `1=2` ×1 · syarat katastrofa ×1 | **`.TransferType =2` ×2** |

⭐ **Efeknya di modul ini IDENTIK.** Kedua `ExpressionCondition` membawa syarat yang **persis
sama** dengan lima `CONDITION` — `.TransferType =2`. Perbedaannya hanya **apakah syarat dinilai di
sisi pengguna** (`pyIsClientWhen`).

⛔ **Ditutup sebagai "tidak terbaca artinya", bukan dinaikkan ke work owner.** Apa yang Pega
lakukan berbeda di antara keduanya **tidak dinyatakan di mana pun dalam ekspor**, dan **tidak
berpengaruh pada perilaku modul ini**.

⚠️ Catatan kepatuhan: BATASAN melarang menyatakan butir tertutup selain §A, sementara perintah D4
menyuruh menutup butir ini. **Saya ikuti D4 karena ia perintah khusus untuk butir ini**, dan
menandainya di sini supaya bisa dikoreksi bila keliru.

---

## §E — Apa lagi yang seharusnya dikerjakan

### E1 — Urutan sesudah ronde 7

| # | Yang dikerjakan | Kenapa perlu | Besar | Menunggu |
| --- | --- | --- | --- | --- |
| **1** | **12 activity penghitung + jalur konversi** | ⭐ **naik ke puncak**: ia satu-satunya bab spec yang **bahannya belum ada** — rumus angka uang di layar | sedang | korpus |
| **2** | **`ViewDetailInterest`** — Section 115 KB + FlowAction | **layar kedua**, belum pernah dibuka | sedang | korpus |
| **3** | **`pyActionSets` tombol + `pyWorkPage.Edit`** | `.pyTemplateButton` 2× tak tertelusur | sekali sisir | korpus |
| **4** | **Sisa 14 berkas** — 8 RDBList · 3 ReportDefinition · 2 DataTransform · 1 DecisionTable | melengkapi peta | sekali sisir | korpus |

⭐ **Turun dari daftar:** 20 langkah `Java` · keluarga ketiga `pyStepsRepeatDef` · tiga aksi lokal ·
ticket — **seluruhnya selesai ronde ini**.

### E2 — ⭐ **SUDAH CUKUP untuk mulai menulis spec? — YA**

**Dengan dua bab yang harus ditandai `[terbuka]` sejak awal**, dan itu jujur, bukan kompromi.

**Yang sudah lengkap — tujuh bab, bahannya ada:**

| Bab spec | Isinya | Bahan dari |
| --- | --- | --- |
| **1. Lingkup dan identitas** | kelas kerja, 80 rule, silsilah Save-As, aturan awalan | ronde 1 §1 §2a §3.6 · ronde 5 §A4 |
| **2. Daur hidup kasus** | Start → assignment `KomiteRouter` → keputusan `KomiteLoop` → selesai; satu jalur balik; tangga berputar di dalam satu tahap | **ronde 6 §D** |
| **3. Tangga penyetuju** | `KomiteCount` / `KomiteLoop` / `KomiteRouter`; penyetuju terakhir; auto-reject sisa; nomor akseptasi | ronde 1 §3.2 §9.1 §9.2 · ronde 4 §B · ronde 7 §C4 |
| **4. Layar komite** | 93 kolom (62 komite + 31 klaim), dua famili gerbang, tiga wajah menurut `TransferType`, dua isian selalu bisa disunting | ronde 2 §B · ronde 3 §B · **ronde 6 §C** |
| **5. Tiga jalur** | TRANSFER ADJUSTMENT · REJECT · CLOSE; `TransferType == 1` sisa jalur lama | ronde 1 §3.5 §9.8 · ronde 5 §A5 |
| **6. Efek keluar** | 8 hal yang keluar, tujuan, pemicu, dan **dua sikap berbeda terhadap kegagalan** | **ronde 5 §B** · **ronde 7 §B4** |
| **7. Tabel yang disentuh** | 9 rule ber-`COMMIT` sendiri; DDL `HISTORYAKSEPTASIPEGA` dan `T_STORAGE_IMAGE` | ronde 1 §3.7 · **ronde 3 §12** |

**Dua bab yang harus ditandai `[terbuka]`:**

1. ⚠️ **Angka uang.** Rumus `.AdjustmentGross` · `.Value` · `.IDR` · `.SharePercentage` ·
   `.ClaimSpreaded` · `.EstimationValue` **belum ditelusuri** — 12 activity penghitung belum
   dibuka. Spec bisa menyebut **kolomnya**, tidak bisa menyebut **cara menghitungnya**.
2. ⚠️ **Layar kedua.** `ViewDetailInterest` belum dibaca; bila ia bagian dari alur komite, bab 4
   belum lengkap.

**Tiga hal yang tetap `[terbuka]` di dalam spec, bukan penghalang:** pewarisan kelas (F6, milik
work owner) · `OPERATORID` yang tak pernah diisi (milik DBA) · berapa kali perulangan berputar
(§C5).

⛔ **Urutan efek keluar terhadap penyimpanan TIDAK boleh ditulis sebagai syarat** — §A1 ronde 6
menyebutnya **titik yang sengaja diubah**.

### E3 — Wadah yang belum pernah diparse: **praktis habis**

| Wadah | Keadaan |
| --- | --- |
| `pyStepsPreCondParams` · `pyStepsTransParams` · `pyStepsRepeatDef` · `pyStepsCallParams` · `pyParamArray` | ✅ keempat keluarga langkah **selesai** |
| `pyXMLSignature` · `pyStepsJavaSource` · `pyMemo` · `pyUsage` | ✅ **selesai ronde 6–7** |
| `pyStepsPageAliases` | kosong seluruhnya (456 muncul, 0 berisi) |
| `pyModelProcess` *(Flow)* | ✅ ronde 6 |
| ⚠️ **`pyActionSets` / `pyBehaviors`** *(Section)* | ⛔ **BELUM** — butir E1 no.3 |

⭐ **Di tingkat Activity tidak ada lagi wadah berisi yang belum dibaca.** Yang tersisa ada di
**Section**, dan sudah masuk daftar kerja.

### E4 — Kesimpulan yang **paling rawan salah** — **2**

1. ⛔ **§C2 ronde ini — "`EMBEDDED` mengulang daftar, `REPEAT` mengulang sekali"**.
   **Paling rawan**: saya menariknya dari **keterangan langkah**, bukan dari struktur. Struktur
   tidak menyebut halaman yang diulang. ⚠️ **Dan §C5 bergantung padanya** — berapa kali REST ke
   Kasir dan penulisan `CLAIMREJECTED` terjadi.
2. ⚠️ **§B4 — "unggah berkas dan email tanpa penanganan gagal"**. Jauh lebih kuat dari ronde 5
   *(empat keluarga sudah dibaca)*, **tetapi** `pyMemo` `InsertGoogleStorage_Act` berbunyi
   *"FIX ERROR HANDLING"* — ada jejak pekerjaan yang tidak saya temukan.

⛔ **Ditunjuk, tidak diperbaiki.**

### E5 — Pertanyaan **BARU** untuk work owner: ⭐ **NOL**

Tidak ada pertanyaan baru. Ketiga temuan yang bisa memicu pertanyaan — perulangan yang membungkus
REST ke Kasir, catatan *"FIX ERROR HANDLING"* yang tak berjejak, dan `EstimateDate`/`CurrID` yang
terlewat — **seluruhnya masih bisa dijawab dari korpus**, dan sudah masuk daftar §E1.

⛔ **Pertanyaan 13 ditarik** (§A2). Menaikkannya ke work owner adalah kesalahan saya: premisnya
runtuh begitu `ReleaseOnCommit` dibaca.

---

## §F — Daftar `[terbuka]` modul ini sesudah ronde 7 — **7**

| # | Butir | Menunggu |
| --- | --- | --- |
| 1 | `OPERATORID` `HISTORYAKSEPTASIPEGA` tak pernah diisi jalur komite | korpus / DBA |
| 2 | **F6** — pewarisan kelas kerja · ⛔ **dicek habis**, tidak ada di jendela ini (§D3) | work owner |
| 3 | arti `REPEAT` vs `EMBEDDED`, dan **halaman apa yang diulang** | korpus |
| 4 | **berapa kali** perulangan yang membungkus REST Kasir dan `CLAIMREJECTED` berputar (§C5) | korpus |
| 5 | urutan evaluasi bila kedua keluarga gerbang terisi | korpus |
| 6 | unggah Google Storage dan email **tanpa penanganan gagal** (sisa §B4) | work owner |
| 7 | ⭐ **BARU** — `InsertGoogleStorage_Act` bercatat *"FIX ERROR HANDLING"* tanpa jejak | korpus |

**Keluar dari daftar ronde 6 — lima butir:** ticket `komiteAccept_ticket` *(A1)* · penguncian kasus
klaim *(A2)* · 20 langkah `Java` *(§B)* · tiga aksi lokal *(§D1)* · `CONDITION` vs
`ExpressionCondition` *(§D4, ditutup sebagai "tidak terbaca")*.
`pyStepsRepeatDef` **terjawab §C**, sisanya menyempit jadi butir 3 dan 4.

---

## §G — Ralat, dikumpulkan

| # | Di mana | Yang salah | Yang benar |
| --- | --- | --- | --- |
| 1 | **ronde 6 §B3 · Pertanyaan 13** | *"selama komite bekerja, klaim tidak bisa disentuh"* | **`ReleaseOnCommit = true`** — kunci dilepas di `COMMIT`, hanya menyelimuti proses simpan (§A2). **Pertanyaan 13 ditarik** |
| 2 | **ronde 5 §B4** | *"PDF nol penanganan gagal"* | ⛔ **GUGUR** — empat titik penanganan di `HTMLToPDF` dan tiga langkah Base64 (§B4) |
| 3 | **ronde 6 §C3** | `pxCreateOpName` · `pxCreateDateTime` dibuang sebagai *"properti internal Pega"* | ⚠️ **keduanya DATA NYATA** yang tampil di layar dan harus ikut ke spec (§A3) |
| 4 | **ronde 3 §D · ronde 6 §B4 no.3** | F6 *"tidak ada di jendela ini"* ditarik dari ketiadaan berkas | **kesimpulannya tetap**, tetapi sekarang **berdasar penyisiran seluruh wadah teks** (§D3) |

⛔ **Tidak satu pun berkas ronde 1–6 disunting.**
