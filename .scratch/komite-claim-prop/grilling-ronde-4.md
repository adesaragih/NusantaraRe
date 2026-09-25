# Grilling Ronde 4 — Komite Claim Prop

**Tanggal:** 2026-09-18 · **Korpus:** `D:\XML\RNM_BRD\Komite Claim Prop\` (READ-ONLY)
**Lingkup:** Komite Claim Prop saja `[keputusan work owner]` 2026-09-18.
**Sasaran tunggal:** `pyStepsTransParams` — keluarga gerbang kedua, sampai tuntas.
**Metode:** parser XML bersarang. Bukti = path berkas + nama rule + nomor step Pega.
**Tidak ada nomor baris XML di berkas ini.**

> **Ralat dicatat DI SINI**, berkas ronde 1 / 2 / 3 tidak disunting.

---

## 0. Ringkasan — **alarm saya benar dinaikkan, hasilnya lega**

| # | Hasil | Status |
| --- | --- | --- |
| **§A1** | Total baris keluarga kedua | ✅ **552** — **cocok persis** dengan angka ronde 3 |
| **§A2** | Peran keluarga kedua | ⭐ **diuji SESUDAH langkah jalan** — terbaca dari isi syaratnya |
| **§A3** | Langkah punya kedua keluarga | ⭐ **552 dari 552** — **seluruhnya**, bukan sebagian |
| **§A6** | Saklar hidup-mati | ⭐ **`pyStepsTransition`** — logika tiga nilai, sama seperti keluarga 1 |
| **§B** | Kesimpulan yang **terbalik** | ✅ **NOL** — disensus 100% |
| **§B** | **Tambahan** baru | ⚠️ **3** — penanganan **kegagalan**, belum pernah tercatat |
| **§C1** | Sasaran lompatan menggantung | ⚠️ **3 → 5** · ⭐ **kelimanya bergerbang mati** |
| **§D2** | Keluarga elemen **KETIGA** | ⚠️ **ADA** — `pyStepsRepeatDef`, terisi di **85** langkah |

⭐ **Kesimpulan ronde ini, satu kalimat:** keluarga kedua **tidak membatalkan satu pun** kesimpulan
ronde 1 / 2 / 3, tetapi ia menyimpan **seluruh penanganan kegagalan** modul ini — yang selama tiga
ronde **tidak pernah terbaca**.

⚠️ **Ralat penilaian saya sendiri:** di §F1 ronde 3 saya menaikkan sasaran ini ke **puncak** dengan
alasan *"berpotensi tidak lengkap, termasuk K1 dan K2"*. Kekhawatiran itu **tidak terbukti**.
Menaikkannya tetap benar — tidak ada cara tahu tanpa menyisir — tetapi **bobotnya saya lebihkan**
di atas efek keluar. §D1 menyusun ulang urutannya.

---

## §A — Sensus penuh `pyStepsTransParams`

### A1 — Total dan sebaran `[terverifikasi]`

Disisir **35 Activity** modul ini dengan parser bersarang.

| Kode `WhenTrue`/`WhenFalse` | Jumlah | Pembacaan dengan tabel `[data work owner]` |
| --- | --- | --- |
| `2`/`2` | **536** | Continue Whens / Continue Whens — **tidak ada tindakan khusus** |
| `_`/`_` | **7** | kedua dropdown kosong — **tidak ada tindakan dipilih** |
| `1`/`2` | **4** | **Jump to Later Step** bila benar |
| `6`/`2` | 1 | **Exit Activity** bila benar |
| `2`/`6` | 1 | **Exit Activity** bila **salah** |
| `1`/`1` | 1 | **Jump** di kedua sisi, sasaran berbeda |
| `1`/`_` | 1 | **Jump** bila benar; tidak ada tindakan bila salah |
| `6`/`6` | 1 | **Exit Activity** apa pun hasilnya |
| | **TOTAL 552** | |

✅ **Cocok persis dengan 552 yang saya sebut di ronde 3. Selisih NOL.**

⭐ **Hanya 16 baris dari 552 yang membawa tindakan khusus** (552 − 536 `2/2`). Sisanya **polos**.

### A2 — Apa **peran** keluarga kedua `[terverifikasi]`

⛔ **Tidak ditebak dari namanya.** Dibaca dari **isi syaratnya**.

`pxObjClass` keluarga kedua adalah **`Embed-ActivityTransitions`**; keluarga pertama
**`Embed-ActivityPreConditions`**. Elemen syaratnya `pyStepsTransParamsWhen`.

**Bukti yang menentukan — isi syarat keluarga kedua:**

| Syarat | Muncul | Apa yang dirujuknya |
| --- | --- | --- |
| **`StepStatusFail`** | **4×** | ⭐ **hasil langkah yang BARU SAJA dijalankan** |
| `UploadDoc.exp==""` | 1× | isi halaman **setelah** unggah dikerjakan |
| `.StatusKonversi=="1"` | 1× | status **setelah** dibaca |
| `true` | 2× | selalu |
| `pyWorkPage.TransferType=="1" && …` | 1× | keadaan kasus |

⭐ **`StepStatusFail` hanya punya arti setelah langkah dijalankan** — ia menanyakan *"apakah
langkah barusan gagal?"*. Jadi:

| Keluarga | `pxObjClass` | Diuji **kapan** | Terhadap **apa** |
| --- | --- | --- | --- |
| **1 — `pyStepsPreCondParams`** | `Embed-ActivityPreConditions` | **SEBELUM** langkah jalan | keadaan halaman — **boleh jalan atau tidak** |
| **2 — `pyStepsTransParams`** | `Embed-ActivityTransitions` | **SESUDAH** langkah jalan | **hasil langkah itu** — **ke mana lanjutnya** |

⚠️ `[terbuka]` **Urutan itu saya baca dari ISI syarat, bukan dari label struktur.** Tidak ada
elemen yang menuliskan *"diuji sesudah"*. Bukti `StepStatusFail` kuat, tetapi **bukan pernyataan
eksplisit korpus**. Ditandai supaya tidak dianggap lebih pasti dari yang sebenarnya.

### A3 — Langkah yang punya **kedua** keluarga `[terverifikasi]`

| | Jumlah |
| --- | --- |
| Langkah seluruhnya di 35 Activity | **552** |
| Punya **kedua** wadah | **552** |
| Hanya `pyStepsPreCondParams` | **0** |
| Hanya `pyStepsTransParams` | **0** |

⭐ **Setiap langkah selalu membawa kedua wadah.** Jadi keduanya **bukan pilihan**, melainkan
bagian tetap struktur langkah — yang membedakan hanyalah **terisi atau polos**. Terisi:
keluarga 1 **548**, keluarga 2 **545**.

⚠️ `[terbuka]` **Urutan evaluasi bila keduanya terisi tidak terbaca dari struktur.** Secara
logika A2 keluarga 1 mendahului, tetapi **korpus tidak menyatakannya**. Tidak ditebak.

> Catatan angka: **552 langkah** tetapi **566 baris** keluarga 1 (ronde 2) — karena satu langkah
> boleh punya **lebih dari satu** baris syarat. Keduanya benar.

### A4 — Enam baris berkode **1**, berikut sasarannya `[terverifikasi]`

| Activity | Step | Kode | Syarat | Sasaran | Tanda ada di |
| --- | --- | --- | --- | --- | --- |
| `HitServiceToKasirKMT_Act` | **10.7** | `1/2` | `StepStatusFail` | `END` | ✅ **langkah 16** |
| `HitServiceToKasirKMT_Act` | **14.4** | `1/2` | `StepStatusFail` | `END` | ✅ **langkah 16** |
| `KomitePostAdjustment` | **29** | `1/2` | `StepStatusFail` | `ENDSERVICE` | ✅ **langkah 30** |
| `KomitePostAdjustment` | **34** | `1/2` | `StepStatusFail` | `ENDKASIR` | ✅ **langkah 35** |
| `KomitePostAdjustment` | **26** | `1/1` | `true` | `CHK` *(benar)* / `EXT` *(salah)* | ⛔ **`CHK` TIDAK ADA** · `EXT` = langkah 25 |
| `KomitePost_Reject` | **16** | `1/_` | `TransferType=="1" && …` | `SKP` | ⛔ **TIDAK ADA** |

### A5 — Baris berkode **4**, **5**, **6** `[terverifikasi]`

| Kode | Jumlah di keluarga 2 | Rincian |
| --- | --- | --- |
| **4** — Exit Iteration | ⛔ **NOL** | tidak dipakai sama sekali di keluarga ini |
| **5** — Skip Whens | ⛔ **NOL** | tidak dipakai sama sekali di keluarga ini |
| **6** — Exit Activity | **3 baris** | lihat bawah |

| Activity | Step | Kode | Syarat | Artinya |
| --- | --- | --- | --- | --- |
| `GetUrlGoogleStorage_Act` | 4 | `6/2` | `UploadDoc.exp==""` | bila kedaluwarsa **kosong** → **hentikan activity** |
| `HitServiceToKasirKMT_Act` | 3 | `2/6` | `.StatusKonversi=="1"` | bila status konversi **bukan** 1 → **hentikan activity** |
| `KomiteRouter` | 6.1 | `6/6` | `true` | **selalu hentikan activity** sesudah langkah ini |

⭐ **Kode 4 dan 5 ternyata milik keluarga PERTAMA saja.** Itu memperkuat pembagian peran A2:
*Exit Iteration* dan *Skip Whens* adalah tindakan **saat memeriksa syarat**, bukan **sesudah
langkah berjalan**.

### A6 — Saklar hidup-mati: **`pyStepsTransition`** `[terverifikasi]`

Sebaran atas 552 langkah: `(kosong)` **533** · `true` **12** · `false` **3** · `0` **4**.

⭐ **Logikanya SAMA PERSIS dengan `pyStepsPreCondition` di keluarga pertama:**

| Nilai | Arti |
| --- | --- |
| **`true`** | transisi **berlaku** — **HIDUP** |
| **`false`** | transisi **dimatikan** — tersimpan tapi tidak dipakai |
| **kosong / `0`** | tidak pernah ada transisi |

**Keenam belas baris khusus, dipilah hidup-matinya:**

| Activity · step | Kode | `pyStepsTransition` | Keadaan |
| --- | --- | --- | --- |
| `GetUrlGoogleStorage_Act` 4 | `6/2` | `true` | ✅ **HIDUP** |
| `HitServiceToKasirKMT_Act` 3 | `2/6` | `true` | ✅ **HIDUP** |
| `HitServiceToKasirKMT_Act` 10.7 | `1/2` | `true` | ✅ **HIDUP** |
| `HitServiceToKasirKMT_Act` 14.4 | `1/2` | `true` | ✅ **HIDUP** |
| `KomitePostAdjustment` 29 | `1/2` | `true` | ✅ **HIDUP** |
| `KomitePostAdjustment` 34 | `1/2` | `true` | ✅ **HIDUP** |
| `KomiteRouter` 6.1 | `6/6` | `true` | ✅ **HIDUP** |
| **`KomitePostAdjustment` 26** | `1/1` | **`false`** | ⛔ **MATI** |
| **`KomitePost_Reject` 16** | `1/_` | **`false`** | ⛔ **MATI** |
| `HTMLToPDF` 9 · `KomitePost_Close` 12.4 · 12.5 · `KomitePost_Reject` 12.4 · 12.5 · `PrintFileAcceptance_TKMT` 10 · 11 | `_/_` | kosong / `0` | — tidak ada tindakan |

**HIDUP: 7 baris.** Mati: **2**. Tanpa tindakan: **7**.

---

## §B — Akibatnya pada kesimpulan yang bersandar alur langkah

### B1 — K1: langkah **17 · 21 · 34** — ✅ **MASIH BERDIRI**

`[terverifikasi]` Baris keluarga kedua pada ketiganya:

| Step | Keluarga 1 | Keluarga 2 |
| --- | --- | --- |
| **17** *(insert OS akseptasi)* | `2/3` penyetuju terakhir + disetujui · `3/2` `IsSubjectivity==true` | **`2/2` — polos** |
| **21** *(generate PDF akseptasi)* | *(idem)* | **`2/2` — polos** |
| **34** *(HIT SERVICE KE KASIR)* | `3/2` `IsSubjectivity==true` · `2/3` penyetuju terakhir + disetujui | ⚠️ **`1/2` `StepStatusFail` → `ENDKASIR`** |

⭐ **Kesimpulan K1 berdiri utuh.** Keluarga kedua **tidak menyentuh gerbang masuk** ketiga langkah
itu. Untuk 17 dan 21 ia polos; untuk 34 ia **tidak menentukan apakah langkah jalan**, melainkan
**apa yang terjadi kalau langkah itu GAGAL**.

⚠️ **TAMBAHAN 1 — belum pernah tercatat:** bila panggilan ke **Kasir gagal**, alur **melompat ke
langkah 35 (`ENDKASIR`)**, yang keterangannya **`SEND EMAIL`**. Jadi kegagalan transfer **tidak
menghentikan activity** — ia meloncati langkah-langkah di antaranya dan **tetap mengirim email**.

### B2 — K2: `KomiteCount` selalu naik — ✅ **MASIH BERDIRI**

`[terverifikasi]` Langkah **40** *(`+ komite count`)*: `pyStepsPreCondition = false` (gerbang
dimatikan), keluarga 1 `2/3 KomiteCount < KomiteLoop`, **keluarga 2 `2/2` — polos**.
**Tidak ada transisi apa pun pada langkah itu.** Kesimpulan K2 **tidak berubah**.

### B3 — Blok penomoran akseptasi **16.1–16.9** — ✅ **TIDAK BERUBAH**

`[terverifikasi]` Kesebelas langkah **16**, **16.1**–**16.9** seluruhnya berkeluarga-2 **`2/2`**.
**Nol transisi.** Urutan **16.5 → 16.6 → 16.7 → 16.8 → 16.9** tetap seperti tercatat, dan tidak
ada lompatan yang menyelinap ke dalam blok itu.

⭐ Tambahan yang menguatkan: langkah **16** ternyata **langkah berulang**
(`pyStepsRepeatDefHasRepeat = REPEAT`) — cocok dengan perannya sebagai **induk blok penomoran**.

### B4 — Jalur tolak: langkah 12 → `EXT` (25) — ✅ **TIDAK BERUBAH**

`[terverifikasi]` Langkah **12** berkeluarga-2 **`2/2`** — polos. Lompatan ke `EXT` murni milik
keluarga pertama, persis seperti tercatat di ronde 3.

⚠️ **TAMBAHAN 2:** langkah **26** membawa transisi `1/1` dengan syarat harfiah `true` — *"selalu
melompat"*, ke `CHK` bila benar dan `EXT` bila salah. **Tetapi `pyStepsTransition = false`**, jadi
**transisi itu DIMATIKAN** dan tidak pernah berjalan. ⭐ **Itulah sebabnya jalur EXT tidak berubah.**
Seandainya dihidupkan, langkah 26 akan **selalu** melompat — dan ke tanda yang **tidak ada**.

⚠️ **TAMBAHAN 3:** langkah **26** juga **langkah berulang** (`EMBEDDED`). Itu menjelaskan mengapa
kode **4 — Exit Iteration** ada di sana pada keluarga pertama: **perulangannya nyata**, bukan
salah pasang.

### B5 — Kesimpulan yang **berubah**

**NOL yang terbalik.** Yang ada hanyalah **tiga tambahan** di atas dan **dua ralat penomoran**:

| # | Di mana | Yang lama | Yang benar |
| --- | --- | --- | --- |
| 1 | ronde 3 §F1 no.1 · §F3 no.1 | *"`pyStepsTransParams` bisa menggoyang K1 dan K2 — paling menahan"* | ⛔ **tidak terbukti.** K1, K2, blok 16, jalur EXT **seluruhnya berdiri** |
| 2 | ronde 1 §5 | `KomiteRouter` *"yang hidup: step **6** dan step **8**"* | langkahnya **1 · 2 · 3 · 4 · 5 · 6 · 6.1 · 7**; **tidak ada step 8**. Yang hidup **6** dan **6.1** — *"step 8"* adalah **6.1** dibaca datar. **Isi kesimpulan tetap benar**: 6 di-remark, tangga `komitepnc1..4` memang mati |
| 3 | ronde 1 §5 | syarat `.KomiteAproval==0 → param.AssignTo := .KomiteID` ditulis milik **step 6** | ia milik **step 6.1**; step 6 adalah **induknya**, bergerbang `false` |

### B6 — Pernyataan tegas

⭐ **Disensus 100 % — 552 dari 552 baris keluarga kedua dibaca. TIDAK ADA kesimpulan ronde 1, 2,
atau 3 yang terbalik karenanya.** Yang bertambah adalah **penanganan kegagalan**, yang selama tiga
ronde memang tidak terlihat karena berada di wadah yang tidak pernah dibuka.

---

## §C — Dua titik gelap yang sejalur

### C1 — Sasaran lompatan menggantung: **3 → 5**, dan **kelimanya mati** `[terverifikasi]`

Seluruh `pyStepsBlockName` di empat activity kunci, **termasuk yang di-remark**:

| Activity | Tanda yang ada |
| --- | --- |
| `KomitePostAdjustment` | `//` di 16.1 · 16.2 · 16.3 · 16.4 · 39 — lalu **`EXT`=25** · **`ENDSERVICE`=30** · **`ENDKASIR`=35** |
| `KomitePost_Close` | ⛔ **NOL tanda** |
| `KomitePost_Reject` | ⛔ **NOL tanda** |
| `HitServiceToKasirKMT_Act` | `//` di 1 — lalu **`END`=16** |

**Jawaban C1: tandanya memang TIDAK ADA** — bukan tersembunyi di wadah yang belum dibaca.
Dan dengan keluarga kedua ikut terbaca, **dua lagi bertambah**:

| # | Sasaran | Di mana | Keluarga | Gerbangnya |
| --- | --- | --- | --- | --- |
| 1 | `JMP` | `KomitePostAdjustment` step 4 | 1 | `pyStepsPreCondition = false` ⛔ **mati** |
| 2 | `EXT` | `KomitePost_Close` step 1 | 1 | `pyStepsPreCondition = false` ⛔ **mati** |
| 3 | `EXT` | `KomitePost_Reject` step 1 | 1 | `pyStepsPreCondition = false` ⛔ **mati** |
| 4 | **`CHK`** ⭐ baru | `KomitePostAdjustment` step 26 | 2 | `pyStepsTransition = false` ⛔ **mati** |
| 5 | **`SKP`** ⭐ baru | `KomitePost_Reject` step 16 | 2 | `pyStepsTransition = false` ⛔ **mati** |

⭐ **Polanya konsisten dan menenangkan: setiap lompatan yang rusak adalah lompatan yang dimatikan.**
Tidak satu pun lompatan **hidup** menunjuk tanda yang tidak ada. Empat lompatan hidup — `END` ×2,
`ENDSERVICE`, `ENDKASIR` — **seluruhnya ketemu tandanya**.

### C2 — Langkah `Java` di `SetKomiteList_Act` step 6 `[terverifikasi]`

Keterangan Pega-nya: *"Hapus currency yg sama"*. Sumbernya **20 baris**.

**Perilakunya, dalam bahasa biasa:** ia menelusuri daftar sementara **`TempTotalAdj`** — satu-
satunya halaman yang disebutnya — membandingkan isi kolom mata uang antar baris, lalu **membuang
baris yang mata uangnya sudah pernah muncul**. Hasilnya: daftar berisi **setiap mata uang tepat
sekali**, sebagai dasar penjumlahan per mata uang di langkah 7 dan seterusnya.

⛔ **Kode Java tidak disalin ke berkas ini.** Perkakas yang dipakainya hanya pembacaan halaman dan
daftar — mencari halaman, mengambil nilai properti sebagai teks, membaca panjang daftar, membuang
satu baris. **Tidak ada panggilan ke basis data, tidak ada efek keluar, tidak ada penulisan ke
kasus.**

⭐ Butir gelap §F2 no.5 ronde 3 **terjawab** — ia hanya pembersih duplikat, bukan logika
tersembunyi.

---

## §D — Apa lagi yang seharusnya dikerjakan

### D1 — Urutan **disusun ulang**, karena §A datang bersih

> Di ronde 3 sasaran ronde ini saya taruh di puncak. Sekarang terbukti benign, jadi ia **turun
> dari daftar** dan efek keluar naik.

| # | Yang dikerjakan | Kenapa perlu | Besar | Menunggu |
| --- | --- | --- | --- | --- |
| **1** | **Efek keluar** — `SendAcceptationToKasir` (REST) · `HTMLToPDF` · `InsertDocument_Act` · `SendEmailWithAttachments` · `GetLinkService` · SystemSettings `LinkService` | ⭐ **naik ke puncak.** §B justru menambah alasannya: kini terbukti ada **penanganan kegagalan** pada panggilan Kasir dan arasapas, tetapi **apa yang dipanggil belum pernah dibaca** | berat | korpus |
| **2** | **`HitServiceToKasirKMT_Act` dibaca utuh** | ia memuat **3 dari 7** transisi hidup modul ini, termasuk gerbang `.StatusKonversi=="1"` yang **menghentikan activity**. Bersambung langsung dengan butir 1 | sedang | korpus |
| **3** | **Berkas `Flow`** — daur hidup kasus komite | belum pernah dibuka; ia yang menentukan **urutan tahap** | sekali sisir | korpus |
| **4** | **Keluarga ketiga `pyStepsRepeatDef`** — 85 langkah berulang | §D2. Perulangan menentukan **berapa kali** langkah jalan; `Exit Iteration` sudah terbukti nyata di langkah 26 | sedang | korpus |
| **5** | **`ViewDetailInterest`** — Section 115 KB + FlowAction | **layar kedua** modul ini, belum pernah dibuka | sedang | korpus |
| **6** | **15 activity penghitung + jalur konversi** | mengisi angka di layar | sedang | **work owner** *(pertanyaan lingkup yang belum dijawab)* |
| **7** | **Sisa 14 berkas** — 8 RDBList · 3 ReportDefinition · 2 DataTransform · 1 DecisionTable | melengkapi peta | sekali sisir | korpus |

### D2 — Keluarga elemen **KETIGA**: ⚠️ **ADA**

⛔ **Dicari, tidak ditunggu.** Seluruh wadah anak pada satu langkah, di 35 Activity:

| Wadah | Muncul | **Berisi** | Sudah dibaca? |
| --- | --- | --- | --- |
| `pyStepsPreCondParams` | 552 | **548** | ✅ ronde 1–3 |
| `pyStepsTransParams` | 552 | **545** | ✅ **ronde ini** |
| **`pyStepsRepeatDef`** | 552 | **85** | ⚠️ **BELUM** |
| `pyParamArray` | 529 | **482** | ✅ *(dipakai untuk sensus penulis properti)* |
| `pyStepsCallParams` | 529 | **389** | ✅ ronde 2 §C |
| `pyStepsPageAliases` | 456 | **0** | — kosong seluruhnya |
| `pyExpressionGadget` | 6 | 6 | ⚠️ **BELUM** — kecil |

⭐ **Keluarga ketiga = `pyStepsRepeatDef`**, isinya `pyStepsRepeatDefHasRepeat` (**85**),
`…Iteration` · `…Start` · `…Limit` (**21** masing-masing).

`[terverifikasi]` Sebarannya per activity — **21 activity** punya langkah berulang; terbanyak
`CountEstimation_Act` **18** · `CountSpreading_act` **10** · `HitServiceToKasirKMT_Act` **9** ·
`CountPersen_act` **8**. Di `KomitePostAdjustment` hanya **2**: langkah **16** (`REPEAT`) dan
langkah **26** (`EMBEDDED`).

⚠️ `[terbuka]` **BARU** — beda `REPEAT` dan `EMBEDDED` **belum terbaca**, dan **halaman apa yang
diulang belum ditelusuri**. `pyStepsRepeatType` **kosong di seluruh 552 langkah**, jadi jenis
perulangan tidak diambil dari sana.

### D3 — Kesimpulan yang **paling rawan salah** — **3**

1. ⛔ **Ronde 2 §B "85 sel / 55 properti"** — **kini yang paling rawan**, menggantikan yang lama.
   Rawan **karena definisinya**, dan pertanyaannya sudah naik ke work owner tanpa jawaban. Selama
   belum dijawab, **setiap kalimat "tidak tampil di layar"** di ronde 2 berdiri di atas definisi
   yang belum disepakati.
2. ⚠️ **§A2 ronde ini — "keluarga kedua diuji SESUDAH langkah jalan"**. Saya membacanya dari **isi
   syarat** (`StepStatusFail`), bukan dari label struktur. Bila keliru, **Tambahan 1** (kegagalan
   Kasir tetap mengirim email) ikut keliru. Ini kesimpulan **baru** dan **belum teruji silang**.
3. ⚠️ **Ronde 1 §3.7 "9 rule ber-`COMMIT` sendiri"** — diturunkan dengan pencocokan teks, bukan
   parser. Tetap rawan, tidak berubah dari ronde 3.

⛔ **Ketiganya ditunjuk, tidak diperbaiki.**

### D4 — Pertanyaan **BARU** untuk work owner — **3**

> Enam pertanyaan lama masih menggantung dan **tidak diulang di sini**: 14 properti tanpa penulis ·
> `.TreatyType` di balik gerbang mati · nomor klaim di balik `NEVER` · isian jalur tolak/tutup
> hanya dua · tiga lompatan menggantung · definisi kolom layar 55 atau 117.

#### Pertanyaan 7 — Kalau kiriman ke Kasir gagal, email tetap dikirim. Memang begitu?

**Apa yang ditanyakan.** Saat komite menyetujui, sistem mengirim data ke Kasir. Bila pengiriman itu
**gagal**, alur **tidak berhenti** — ia melompat ke langkah yang **mengirim email**, lalu selesai.
Apakah itu memang yang diinginkan.

**Kenapa muncul.** Baru terbaca ronde ini; selama tiga ronde penanganan kegagalan tidak terlihat.
Pola yang sama ada pada kiriman ke arasapas.

**Bedanya jawaban A atau B.** Bila **memang begitu**, sistem baru meniru apa adanya: gagal kirim,
email tetap jalan, dan seseorang harus memeriksa belakangan. Bila **tidak**, ini cacat yang
selama ini tidak ketahuan, dan penanganannya perlu dirancang — apakah dicoba ulang, dibatalkan,
atau ditandai untuk ditinjau manusia.

**Yang tertahan.** Rancangan efek keluar, dan sikap terhadap kegagalan di seluruh modul.

#### Pertanyaan 8 — Kiriman ke Kasir berhenti bila status konversi bukan "1". Apa artinya?

**Apa yang ditanyakan.** Sebelum mengirim ke Kasir, sistem memeriksa sebuah penanda status
konversi. Bila penanda itu **bukan bernilai satu**, seluruh proses pengiriman **dihentikan saat
itu juga**, tanpa pesan. Apa arti penanda itu, dan siapa yang mengisinya.

**Kenapa muncul.** Ia salah satu dari tujuh transisi hidup modul ini, dan satu-satunya yang
**menghentikan** proses berdasarkan keadaan data, bukan kegagalan teknis.

**Bedanya jawaban A atau B.** Bila ia **penanda sah** — misalnya menandai klaim yang sudah
dikonversi — maka ia jadi syarat resmi yang perlu ditegakkan. Bila ia **penanda sementara** dari
masa lalu, menyalinnya justru akan menahan kiriman yang seharusnya jalan.

**Yang tertahan.** Kapan data boleh dikirim ke Kasir.

#### Pertanyaan 9 — Dua transisi sengaja dimatikan: dibiarkan mati, atau pernah dimaksudkan hidup?

**Apa yang ditanyakan.** Dua langkah menyimpan perintah melompat yang **sengaja dimatikan**
saklarnya — satu di jalur adjustment, satu di jalur tolak. Keduanya menunjuk tanda yang **tidak
ada di activity-nya**. Apakah keduanya memang ditinggalkan, atau pernah dimaksudkan hidup dan
tandanya terhapus.

**Kenapa muncul.** Salah satunya bersyarat harfiah *"selalu"*. Kalau saklarnya pernah hidup,
langkah itu **selalu melompat** — dan alurnya akan sangat berbeda dari yang tercatat sekarang.

**Bedanya jawaban A atau B.** Bila **ditinggalkan**, keduanya dicatat sebagai jalur mati dan tidak
dipindahkan. Bila **pernah dimaksudkan hidup**, perlu dipastikan ke mana seharusnya melompat — dan
itu **hanya bisa dijawab orang yang menulisnya**.

**Yang tertahan.** Tidak ada yang tertahan sekarang. Tertahan bila kelak ada yang menghidupkan
saklarnya tanpa tahu tandanya hilang.

---

## §E — Daftar `[terbuka]` modul ini sesudah ronde 4 — **13**

| # | Butir | Menunggu |
| --- | --- | --- |
| 1 | `OPERATORID` `HISTORYAKSEPTASIPEGA` tak pernah diisi jalur komite | korpus / DBA |
| 2 | **F6** — pewarisan kelas kerja *(permintaan sudah dirumuskan, §E2 ronde 3)* | work owner |
| 3 | 14 properti layar tanpa penulis di modul ini | work owner |
| 4 | beda `CONDITION` vs `ExpressionCondition` tidak terbaca | korpus |
| 5 | `.TreatyType` hanya di balik gerbang mati `1=2` | work owner |
| 6 | `pyWorkCover.ClaimData.NoClaim` di balik `NEVER` | work owner |
| 7 | tiga — kini **lima** — sasaran lompatan menggantung | work owner *(pertanyaan 9)* |
| 8 | definisi kolom layar 55 atau 117 | work owner |
| 9 | ⭐ **BARU** — keluarga ketiga `pyStepsRepeatDef`, 85 langkah berulang belum dibaca | korpus |
| 10 | ⭐ **BARU** — beda `REPEAT` vs `EMBEDDED`, dan halaman apa yang diulang | korpus |
| 11 | ⭐ **BARU** — urutan evaluasi bila kedua keluarga terisi, tidak terbaca dari struktur | korpus / work owner |
| 12 | ⭐ **BARU** — kegagalan kirim Kasir tetap mengirim email | work owner *(pertanyaan 7)* |
| 13 | ⭐ **BARU** — arti `.StatusKonversi=="1"` yang menghentikan kiriman Kasir | work owner *(pertanyaan 8)* |

**Yang keluar dari daftar ronde 3:** butir `pyStepsTransParams` *(terjawab §A)* dan butir langkah
`Java` *(terjawab §C2)*. ⛔ **Tidak ada butir dinyatakan tertutup** — itu keputusan work owner.

---

## §F — Ralat, dikumpulkan

| # | Di mana | Yang salah | Yang benar |
| --- | --- | --- | --- |
| 1 | ronde 3 §F1 no.1 · §F3 no.1 | `pyStepsTransParams` *"paling menahan, bisa menggoyang K1 dan K2"* | **tidak terbukti** — nol kesimpulan terbalik (§B6). Menaikkannya tetap benar; **bobotnya dilebihkan** |
| 2 | ronde 1 §5 | `KomiteRouter` *"yang hidup step 6 dan step 8"* | **tidak ada step 8**; yang hidup **6** dan **6.1**. Isi kesimpulan tetap benar |
| 3 | ronde 1 §5 | `.KomiteAproval==0 → param.AssignTo` milik step 6 | milik step **6.1** |
| 4 | ronde 3 §F2 no.5 | langkah `Java` *"bisa melakukan apa saja"* | hanya **membuang baris mata uang duplikat** (§C2) |
| 5 | ronde 3 §A3 | *"tiga sasaran lompatan menggantung"* | **lima** — `CHK` dan `SKP` bertambah (§C1) |

⛔ **Tidak satu pun berkas ronde 1, 2, atau 3 disunting.**
