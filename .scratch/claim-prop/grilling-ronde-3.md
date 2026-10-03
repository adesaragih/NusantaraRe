# Grilling Ronde 3 — Claim Prop

**Tanggal:** 2026-09-19 · **Korpus:** `D:\XML\RNM_BRD\Claim Prop\` (READ-ONLY)
**Sasaran:** tiga wilayah yang belum pernah dibaca — `pyMemo`, 84 berkas yang tidak pernah dibuka,
dan 224 langkah berulang.

> ⛔ **Berkas ronde 1, ronde 2, `periksa-ulang-aturan-baru.md`, `spec.md`, dan 15 tiket lainnya
> TIDAK disunting.** Satu-satunya berkas lama yang disentuh adalah **tiket 08**, dan hanya satu
> baris — lihat §A.

---

## §A — Tiket 08 ditambal lebih dulu ✅

`[terverifikasi]` Jalur kegagalan penyimpanan data adjuster ditambahkan ke tabel rule sumber tiket
**08**:

> `SaveAdjusterConsultant_Act` **langkah 4** (`Obj-Save`) punya jalur kegagalan di **keluarga
> gerbang kedua**: syarat *"langkah barusan gagal"* → **lompat ke tanda `FAIL`, yaitu langkah 8**.
> ⚠️ **Langkah 8 bergerbang MATI**, jadi lompatannya sampai ke langkah yang **tidak berbuat
> apa-apa**, lalu alur berlanjut ke langkah 9. **Kegagalan penyimpanan diam-diam diabaikan.**

Ditandai `[terbuka]` — **menunggu work owner**. ⛔ Tiket itu **tidak menetapkan** apa yang
seharusnya terjadi saat gagal.

⛔ Bagian lain tiket 08 **utuh**. Tiket **00** dan **06** **tidak disentuh** — keduanya penambahan
keterangan, ditambal sekali jalan sesudah ronde ini.

---

## §B — PREDIKSI, DIBEKUKAN ⭐

> ## ⛔ BAB INI DITULIS SEBELUM SATU BARIS PUN DISENSUS
>
> Ditulis **2026-09-19**, sebelum `pyMemo`, 84 berkas, dan 224 langkah berulang disentuh.
> **Sesudah ini ia BEKU** — tidak diubah, tidak dihaluskan, tidak digeser kata-katanya supaya
> kelihatan kena.
>
> ⚠️ **Prediksi yang meleset adalah HASIL, bukan kegagalan.** Melesetnya berarti ronde 3 menemukan
> sesuatu yang nyata.

Ketiga prediksi ini **milik asisten**, diajukan 2026-09-19.

### P1 — `pyMemo`

> **≤ 10 catatan yang mengubah pembacaan**, **DAN ≥ 1 menunjuk sesuatu yang jejaknya tidak ada di
> wadah mana pun.**

*Dasarnya:* di modul Komite Claim Prop, dari 73 catatan **7** mengubah pembacaan dan **1**
(*"FIX ERROR HANDLING"* pada rule pengunggah berkas) menunjuk pekerjaan yang jejaknya tidak
ditemukan di keempat keluarga wadah. Claim Prop punya **303** catatan — empat kali lipat — tetapi
saya menduga proporsi yang bermakna **tidak** ikut naik empat kali, karena sebagian besar catatan
bersifat administratif.

### P2 — 84 berkas yang belum pernah dibuka

> **≥ 1 dari 29 Activity menyentuh uang atau efek keluar; sisanya pinggiran.**

*Dasarnya:* 84 berkas itu bukan pilihan acak — ronde 1 dan 2 menyisir jalur utama lebih dulu, jadi
yang tertinggal cenderung cabang. Tetapi **29 Activity** terlalu banyak untuk seluruhnya pinggiran.

### P3 — 224 langkah berulang

> **Halaman yang diulang TIDAK terbaca dari struktur, jadi tetap `[terbuka]`.**

*Dasarnya:* di Komite Claim Prop, penanda jenis perulangan **kosong di seluruh 552 langkah**, dan
**tidak ada elemen yang menyebut halaman mana yang diulang**. Saya menduga ekspor Claim Prop
berbentuk sama.

### ⚠️ Yang sudah saya ketahui sebelum memprediksi — ditulis supaya prediksinya jujur

Tiga angka **dasar** sudah dihitung di `periksa-ulang-aturan-baru.md` dan **bukan** bagian dari
prediksi: `pyMemo` **303 berisi** · **84** berkas belum dibuka · **224** langkah berulang
(`EMBEDDED` 192 · `REPEAT` 32).

**Yang diprediksi adalah apa yang ADA DI DALAMNYA**, dan itu belum saya lihat sama sekali.


---

## §C — `pyMemo`: 303 catatan, nol pernah dibaca

*Bab ini menjawab: apakah catatan pengembang menyimpan sesuatu yang tidak terbaca di tempat lain.*

### C1 — Sensus 100 % `[terverifikasi]`

| Ukuran | Jumlah |
| --- | --- |
| Catatan **berisi** | **303** |
| Berkas yang memuatnya | **295** |
| Pasangan berkas+catatan berbeda | **303** |
| Mengandung **kata perubahan** *(fix · perbaiki · ubah · hapus · tambah · ganti · remark · error)* | **95** |

⛔ **95 bukan berarti 95 mengubah pembacaan.** Sebagian besar mencatat suntingan yang **hasilnya
sudah terlihat di ekspor**, dan ronde 1/2 membacanya dengan benar. Yang dinilai di C2 adalah
catatan yang **menyatakan sesuatu yang tidak terbaca dari ekspor**.

### C2 — Catatan yang MENGUBAH atau MENGUATKAN pembacaan

⚠️ **Diuji satu per satu terhadap korpus.** Dari 95 kandidat, **delapan diuji langsung** — yang
paling mungkin menyentuh kesimpulan yang sudah ditulis.

#### ⭐ MENGUBAH — 2

| # | Catatan | Apa yang berubah |
| --- | --- | --- |
| **1** | `ConnectREST/SendAcceptationToKasir.xml` — **"tambah auth"** | `[terverifikasi]` Berkasnya memuat **15 kemunculan kata *auth*** dan **satu `Basic`** — jadi panggilan ke Kasir memakai **autentikasi dasar**. ⛔ **Kata "auth" NOL di seluruh ronde 1, ronde 2, dan periksa-ulang.** Bab efek keluar dan tiket **13** menerangkan panggilan itu **tanpa menyebut ada autentikasi sama sekali** |
| **2** | `Activity/CountPersen_act.xml` — **"perbaiki struktur hitungan dan jgn kali share ceding lagi karna udah dikalikan"** | ⚠️ Catatan itu **melarang mengalikan share ceding lagi**. Tetapi `[terverifikasi]` **langkah 5.1 masih mengalikannya**: `@toDecimal(.ClaimEstimation) * @divide(ShareCeding,100,4)`. ⛔ **Tidak terbaca dari korpus** apakah 5.1 adalah keadaan **sesudah** perbaikan, atau sisa yang luput. `[terbuka]` |

#### ✅ MENGUATKAN — 3

| Catatan | Yang dikuatkan |
| --- | --- |
| `Activity/SetCatastrope_act.xml` — **"hapus save"** | `[terverifikasi]` **nol langkah simpan** di berkasnya. Catatan dan ekspor **sejalan** |
| `Activity/SetDefNonCatastrope_Act.xml` — **"hapus save"** | `[terverifikasi]` **nol langkah simpan**. Sejalan |
| `Activity/GetUrlGoogleStorage_Act.xml` — **"FIX ERROR HANDLING"** | `[terverifikasi]` jejaknya **ADA** — langkah 4, keluarga kedua, kode `6/2` atas syarat kedaluwarsa kosong |

### C3 — ⚠️ Catatan yang menunjuk pekerjaan **TANPA JEJAK** — 1

| Catatan | Yang dicari |
| --- | --- |
| ⭐ `Activity/InsertGoogleStorage_Act.xml` — **"FIX ERROR HANDLING"** | `[terverifikasi]` **tidak ada jejak penanganan apa pun**: nol baris khusus di keluarga pertama, nol di keluarga kedua, nol langkah `Java` yang menangani gagal |

> ⭐ **Pola yang sama persis dengan modul Komite Claim Prop.** Rule bernama sama, catatan berbunyi
> sama, jejak sama-sama tidak ada — **di dua modul berbeda**. Itu memperkuat dugaan bahwa
> perbaikannya **memang tidak ada di ekspor**, bukan terlewat oleh pembacaan.
>
> ⛔ **Di mana jejaknya tidak ditebak.** `[terbuka]`

### C4 — `pyUsage`

`[terverifikasi]` **870 berisi, tetapi hanya EMPAT nilai berbeda**, dan keempatnya **keterangan baku
Pega**: `FLOW` **672** · `java` **188** · dua kalimat stok Pega tentang portal.

⭐ **Nol yang berisi keterangan proyek ini.** ⛔ Tidak perlu dibaca lagi.

### ⚠️ Batas kejujuran bab ini

**303 catatan tidak diadjudikasi satu per satu.** Yang dilakukan: 95 disaring dengan kata kunci,
lalu **delapan yang paling mungkin bermakna diuji langsung terhadap korpus**. Sisa **87 kandidat
belum diuji**. Angka *"2 mengubah · 1 tanpa jejak"* adalah **hasil pengujian sampel terkuat**,
bukan hasil adjudikasi penuh.

---

## §D — Berkas yang belum pernah dibuka

*Bab ini menjawab: seberapa besar wilayah yang tidak pernah dilihat siapa pun, dan apakah isinya
penting.*

### ⛔ RALAT ke `periksa-ulang-aturan-baru.md` — **84 seharusnya 80**

`[terverifikasi]` Angka **84** di berkas periksa-ulang **terlalu tinggi**. Penyebabnya penyaring
saya sendiri: ia mencocokkan **nama berkas persis**, sedangkan ronde 1/2 kerap menyebut rule
**tanpa akhiran** — misalnya *"`CopyOldataCurr` 7.1"* untuk berkas `CopyOldataCurr_act.xml`.

| Cara hitung | Hasil |
| --- | --- |
| Nama berkas **persis** | 84 |
| Akhiran `_act` · `_SQL` · `_Harness` dan sejenisnya **dibuang** | ⭐ **80** |

**Empat yang keliru masuk daftar:** `CopyOldataCurr_act` · `CountValue_Act` ·
`GetPayAttachment_Act` · `CauseofLoss_Harness`.

⚠️ **Dua di antaranya penting** — `CopyOldataCurr_act` justru berkas yang memuat **rantai OR**
temuan periksa-ulang, dan ia **memang sudah dibaca** ronde 2.

⛔ **Berkas periksa-ulang tidak disunting** — ralat dicatat di sini.

### D1 — Inventaris 80 berkas `[terverifikasi]`

| Jenis | Belum pernah dibuka |
| --- | --- |
| **Activity** | **26** |
| **When** | **16** |
| **RDBList** | **11** |
| **Section** | **11** |
| **Harness** | **5** |
| DataTransform | 4 |
| ReportDefinition | 4 |
| FlowAction | 3 |
| **TOTAL** | **80** |

### D2 — Saringan

⚠️ **Saringan otomatis saya terlalu kasar untuk berkas non-Activity.** Ia menandai *"menyentuh
uang"* bila kata bernuansa uang muncul **di mana pun** dalam berkas — dan pada Section atau Harness
yang besar itu hampir selalu benar. Hasilnya **71 dari 84 lolos**, yang tidak berguna sebagai
saringan.

⭐ **Yang bisa dipercaya adalah baris Activity dan When**, karena di sana penilaiannya berdasarkan
**nama properti yang ditulis** dan **metode langkah**, bukan sekadar kata yang muncul.

**Dari 26 Activity yang belum dibuka, 16 menyentuh uang atau efek keluar:**

| Activity | Tanda |
| --- | --- |
| `SetNameCurrency_Act` | ⭐ **UANG + EFEK KELUAR** — inilah pengunci kurs yang dirujuk modul Komite |
| `SetCurencyList_act` · `SetCurrency_Act` | UANG + EFEK KELUAR |
| `GetDataOustanding` · `GetDataOutsClaim_act` | UANG + EFEK KELUAR |
| `GetRNMShareTreaty` · `SetFormat_Act` · `DeleteListClaim_Act` | UANG |
| `CekPolisAvailable_Act` · `GetReportStatus_Act` · `SetMasterID` · `SetEditCatastrope` · `SetCauseOfLossValue_act` · `CNMSetDetailCauseOfLoss_act` | EFEK KELUAR |

⚠️ **`SetNameCurrency_Act` paling menonjol.** Modul Komite Claim Prop menutup butir kurs (K8) dengan
menunjuk berkas ini sebagai **tempat kurs dikunci** — dan di modul pemiliknya sendiri **ia tidak
pernah dibaca**.

### D3 · D4 — 16 rule `When` yang belum dibaca

`[terverifikasi]` Seluruhnya **penggolong lini bisnis**: `IsAviationHull` · `IsBillboardNeon` ·
`IsBoiler` · `IsBurglary` · `IsCAR` · `IsCIS` · `IsEar` · `IsExclusion` · `IsFidelity` · `IsGlass` ·
`IsLandRig` dan lima lainnya.

⚠️ **Lima di antaranya `dipanggil 0`** — tidak disebut berkas lain mana pun di modul ini:
`IsAviationHull` · `IsBurglary` · `IsCAR` · `IsFidelity` dan satu lagi. ⛔ Apakah itu berarti mati,
atau dipanggil dari modul lain, **tidak terbaca dari modul ini**. `[terbuka]`

⭐ **Tiga catatan `pyMemo` pada rule `When` menyentuh isinya:**

| Rule | Catatan |
| --- | --- |
| `When/IsAneka.xml` | *"ganti IsBonding jadi IsBondingAndCustomBonds"* |
| `When/isPA_PNC.xml` | *"tambah pyWorkPage dan hapus policy"* |
| `When/IsBonding.xml` | *"Tambah untuk PNC"* |

⚠️ Ketiganya menyentuh **penggolongan lini bisnis**, yang punya tiketnya sendiri — **tiket 04**.

### D5 — Yang pinggiran

**4 `ReportDefinition`** *(penyaring daftar master)* dan sebagian **Harness** *(pembungkus
tampilan)*. ⛔ **Sisanya tidak dapat digolongkan dengan saringan otomatis** dan **belum dinilai
satu per satu** — lihat §G3.

---

## §E — 224 langkah berulang

*Bab ini menjawab: berapa kali sebuah langkah berjalan, dan apakah itu terbaca.*

### E1 · E2 — Isi wadahnya `[terverifikasi]`

| Tag di dalam `pyStepsRepeatDef` | Terisi | Nilainya |
| --- | --- | --- |
| `…HasRepeat` | **224** | `EMBEDDED` **192** · `REPEAT` **32** |
| `…Iteration` | 33 | seluruhnya **`1`** |
| `…Start` | 33 | seluruhnya **`1`** |
| `…Limit` | 33 | **`1`** ×32 · ⭐ **`Local.SizeKurs`** ×1 |

**Beda `REPEAT` dan `EMBEDDED`, diturunkan dari isi:** `REPEAT` **selalu** membawa iterasi, awal,
dan batas; `EMBEDDED` **tidak** — kecuali **satu** pengecualian, karena 33 baris membawa ketiganya
sementara `REPEAT` hanya 32.

### ⛔ E1 — halaman yang diulang: **TIDAK TERBACA** `[terbuka]`

| Ukuran | Hasil |
| --- | --- |
| `pyStepsRepeatType` di seluruh **1366** langkah | ⛔ **kosong seluruhnya** |
| `pyStepPage` pada **224** langkah berulang | ⛔ **kosong seluruhnya** |

⛔ **Tidak ada satu pun elemen yang menyebut halaman mana yang diulang.** ⛔ Tidak diterka dari
keterangan langkah.

⭐ **Satu pengecualian yang terbaca:** satu langkah berbatas **`Local.SizeKurs`** — jumlah putarannya
**sama dengan jumlah mata uang**. Itu satu-satunya dari 224 yang jumlah putarannya dapat dibaca.

### E3 — `GetPayAttachmentAdj_Act` langkah 3 `[terverifikasi]`

⭐ **Perulangannya BERSARANG DUA TINGKAT**, dan keduanya `REPEAT`:

```
step 3      REPEAT                    <- tingkat luar
  3.1       Call GetLinkService
  3.2       Connect-REST              <- panggilan layanan luar
  3.3       REPEAT                    <- tingkat dalam, "Untuk Insert Attachment"
    3.3.4   Obj-Save
    3.3.8   Obj-Save
    3.3.9   Commit                    <- penyimpanan akhir, DI DALAM DUA PERULANGAN
```

⛔ **Berapa kali ketiganya terjadi: `[terbuka]`.** Batas kedua perulangan tertulis `1`, tetapi
halaman yang diulang **tidak terbaca**, jadi angka `1` itu **tidak dapat dipercaya sebagai jumlah
putaran sebenarnya**.

⚠️ **`Commit` di dalam dua perulangan bersarang** berarti satu jalur lampiran dapat menghasilkan
**banyak transaksi terpisah**, bukan satu.

### E4 — Langkah berulang lain yang membungkus penulisan atau efek keluar

Selain yang sudah tercatat di periksa-ulang, **nol tambahan** yang membungkus `Commit` atau
panggilan layanan luar. `GetPayAttachmentAdj_Act` **tetap satu-satunya**.

---

## §F — Adu dengan prediksi §B

⛔ Prediksi dibekukan **2026-09-19 pukul 00:41:48**, sebelum `pyMemo`, 80 berkas, dan 224 langkah
berulang disentuh. Stempel waktu berkas ini membuktikannya.

### P1 — `pyMemo` · **KENA, dengan catatan**

| Bagian prediksi | Hasil |
| --- | --- |
| *"≤ 10 yang mengubah pembacaan"* | ✅ **2 mengubah · 3 menguatkan** |
| *"≥ 1 menunjuk sesuatu yang jejaknya tidak ada"* | ✅ **1** — `InsertGoogleStorage_Act` |

⚠️ **Tetapi ujiannya lebih lemah dari kelihatannya.** Dari 303 catatan, **95 disaring** dan hanya
**8 diuji langsung**. Angka *"2 mengubah"* berlaku atas **8 yang diuji**, bukan atas 303. **87
kandidat belum dinilai.** ⛔ Saya tidak menyebut ini KENA penuh.

### P2 — 80 berkas belum dibuka · ⛔ **MELESET**

| Bagian prediksi | Hasil |
| --- | --- |
| *"≥ 1 dari 29 Activity menyentuh uang atau efek keluar"* | ✅ benar — tetapi angkanya **16** |
| *"sisanya pinggiran"* | ⛔ **SALAH** |

**Kenapa meleset:** jumlah Activity yang benar **26**, bukan 29 — dan **16 dari 26 menyentuh uang
atau efek keluar**. Itu **mayoritas**, bukan sisa. ⛔ Kata *"sisanya pinggiran"* menyiratkan
segelintir; kenyataannya **hanya 10 yang pinggiran**.

⚠️ **Dan prediksinya memang terlalu longgar untuk jadi ujian.** *"≥ 1"* hampir mustahil salah.
Saya catat itu sebagai kelemahan prediksi saya, bukan sebagai kemenangan.

### P3 — 224 langkah berulang · ✅ **KENA**

| Bagian prediksi | Hasil |
| --- | --- |
| *"halaman yang diulang TIDAK terbaca dari struktur"* | ✅ `pyStepPage` **kosong di seluruh 224** · `pyStepsRepeatType` **kosong di seluruh 1366** |

⭐ Persis seperti modul Komite Claim Prop. Satu pengecualian kecil yang **tidak** merusak prediksi:
satu langkah berbatas `Local.SizeKurs`.

### F3 — Apakah ronde 4 masih perlu

⛔ **Satu prediksi MELESET, jadi ronde 3 menemukan sesuatu yang nyata.** Yang jadi ikut
dipertanyakan: **cakupan**, bukan ketepatan. 16 Activity yang menyentuh uang atau efek keluar
**belum pernah dibaca isinya** — termasuk pengunci kurs yang **modul lain sudah bersandar padanya**.

**Jawaban: ya, masih perlu** — tetapi bukan ronde penyisiran lebar. Lihat §G3.

---

## §G — Yang goyah, dan apa lagi

### G1 — Kesimpulan yang GOYAH atau TERBALIK

> ⭐ **TERBALIK: NOL.** Tidak ada kesimpulan ronde 1, ronde 2, atau periksa-ulang yang terbantah.

**GOYAH: 4.**

| # | Kesimpulan | Kenapa goyah | Bukti |
| --- | --- | --- | --- |
| **1** | Bab efek keluar menerangkan panggilan ke Kasir **tanpa menyebut autentikasi** | berkas panggilannya memuat **15 kemunculan *auth*** dan satu `Basic`; kata itu **nol** di seluruh ronde | §C2 no.1 |
| **2** | Bab rumus uang mencatat `CountPersen` 5.1 mengalikan share ceding | catatan pengembangnya **melarang** perkalian itu; tidak terbaca mana yang berlaku | §C2 no.2 |
| **3** | Cakupan modul — *"84 berkas belum dibuka"* | ⛔ **angkanya 80**, dan **16 Activity di antaranya menyentuh uang atau efek keluar** | §D |
| **4** | Penanganan gagal pada unggah berkas | catatan *"FIX ERROR HANDLING"* **tanpa jejak** di keempat keluarga wadah | §C3 |

### G2 — Bab `spec.md` dan tiket yang terpengaruh

⛔ **Keduanya tidak disunting.**

| Sasaran | Bagian | Sifatnya |
| --- | --- | --- |
| **`spec.md`** — Implementation Decisions, bagian efek keluar | keterangan panggilan ke Kasir | ⭐ **penambahan nyata** — autentikasi |
| **`spec.md`** — Implementation Decisions, bagian rumus uang | `CountPersen` dan share ceding | `[terbuka]` baru |
| **tiket `13`** efek keluar Kasir/arasapas/konversi/email | keterangan panggilan Kasir | ⭐ **penambahan nyata** |
| **tiket `00`** prefactor skema transaksi uang | rumus dan pembulatan | `[terbuka]` baru |
| **tiket `11`** penyerahan komite dan penutupan klaim | menyebut pengiriman Kasir | penambahan keterangan |
| **tiket `04`** klasifikasi lini bisnis | tiga rule `When` yang catatannya menyebut perubahan penggolongan | perlu diperiksa |

**Lima tiket** — 00 · 04 · 11 · 13, ditambah **08** yang sudah ditambal di §A.

### G3 — Sesudah ronde 3, urut dari yang paling menahan

| # | Yang dikerjakan | Kenapa | Besar | Menunggu |
| --- | --- | --- | --- | --- |
| **1** | **`SetNameCurrency_Act`** dan 15 Activity lain yang menyentuh uang/efek keluar | ⭐ modul Komite **sudah bersandar** pada berkas ini untuk menutup butir kurs, padahal di modul pemiliknya **belum pernah dibaca** | sedang | korpus |
| **2** | **16 rule `When`** — terutama tiga yang catatannya menyebut perubahan penggolongan | menyentuh **tiket 04**; lima di antaranya **dipanggil nol** | sekali sisir | korpus |
| **3** | **87 catatan `pyMemo` yang belum diuji** | delapan yang diuji menghasilkan dua perubahan dan satu tanpa jejak; **sisanya belum dinilai** | sedang | korpus |
| **4** | **Tambal spec dan lima tiket** — sekali jalan | menunggu 1–3 supaya tidak menambal dua kali | sekali jalan | — |
| **5** | **11 `RDBList` + 11 `Section` + sisa 10 Activity** yang belum dinilai | saringan otomatis tidak dapat menilainya; perlu dilihat satu per satu | sedang | korpus |

### G4 — Kesimpulan yang paling rawan salah

⛔ **Ditunjuk, tidak diperbaiki.**

1. ⛔ **"2 catatan mengubah pembacaan"** *(§C2)*. **Paling rawan.** Ia hasil **8 dari 303** yang
   diuji. Kalau 87 kandidat sisanya dinilai, angkanya hampir pasti naik — dan saya sudah
   memakainya untuk menyatakan P1 KENA.
2. ⚠️ **"Nol tambahan yang membungkus `Commit`"** *(§E4)*. Pola *"menyimpulkan ketiadaan"* yang
   sudah dua kali membakar proyek ini. Ia bersandar pada daftar di periksa-ulang, yang sendirinya
   bersandar pada penyaring metode.
3. ⚠️ **"`pyUsage` nol berisi keterangan proyek"** *(§C4)*. Empat nilai berbeda memang sedikit,
   tetapi saya menilai dari **nilai unik**, bukan dari 870 kemunculannya.

### G5 — Pertanyaan BARU untuk work owner — **2**

#### Pertanyaan 1 — Panggilan ke Kasir memakai autentikasi yang tidak pernah tercatat. Bagaimana kredensialnya dipindahkan?

**Apa yang ditanyakan.** Panggilan ke Kasir ternyata memakai **autentikasi dasar** — catatan
pengembangnya berbunyi *"tambah auth"*, dan berkasnya memuat jejak autentikasi. Spec dan tiket
menerangkan panggilan itu **tanpa menyebutnya sama sekali**. Dari mana kredensialnya diambil, dan
bagaimana ia dipindahkan.

**Kenapa muncul.** Baru terbaca ronde ini, dari catatan pengembang yang belum pernah dibaca. Tiket
**13** sudah ditulis tanpa menyebut autentikasi.

**Bedanya jawaban A atau B.** Bila kredensialnya **tersimpan di tempat yang sama dengan alamat
layanan**, pemindahannya ikut jalur yang sudah ada. Bila ia **tertulis di dalam rule**, maka ada
rahasia di dalam ekspor, dan penanganannya berbeda sama sekali — termasuk soal siapa boleh
melihatnya.

**Yang tertahan.** Tiket 13, dan cara sistem baru menyimpan kredensial layanan luar.

#### Pertanyaan 2 — Catatan pengembang melarang mengalikan share ceding lagi, tetapi rumusnya masih mengalikannya. Mana yang berlaku?

**Apa yang ditanyakan.** Catatan pada rule penghitung persen berbunyi *"jangan kalikan share ceding
lagi karena sudah dikalikan"*. Tetapi langkah 5.1 di rule yang sama **masih mengalikannya**. Apakah
yang ada sekarang sudah benar, atau catatan itu menandai perbaikan yang belum selesai.

**Kenapa muncul.** Korpus tidak menyimpan urutan waktu — tidak terbaca apakah langkah 5.1 adalah
keadaan **sesudah** perbaikan atau sisa yang luput.

**Bedanya jawaban A atau B.** Bila **sudah benar**, rumusnya ditiru apa adanya. Bila **belum**,
maka angka yang dihasilkan sekarang **terlalu kecil satu faktor share ceding** — dan itu menyentuh
nilai klaim yang mengalir sampai ke pembayaran.

**Yang tertahan.** Bab rumus uang di spec, tiket 00, dan setiap angka turunan dari penghitung
persen.
