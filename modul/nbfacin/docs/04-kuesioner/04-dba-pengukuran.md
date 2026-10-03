# Permintaan Pengukuran ke DBA — Mata Uang dan Satuan Rate

**Purpose:** dua **pengukuran atas data produksi**, bukan pertanyaan pendapat. Keduanya menutup
pertanyaan yang tidak dapat dijawab dari kode, dan yang kedua bahkan menutup pertanyaan yang sedang
kami ajukan ke Product/Aktuaria — angka dari basis data membuktikannya langsung tanpa perlu
pendapat siapa pun.

**From:** Tim Migrasi Facultative Inward · **To:** DBA ·
**How your answers will be used:** hasil pengukuran dicatat sebagai keputusan proyek dan menentukan
dua hal di sistem baru: bagaimana nilai uang tanpa mata uang ditangani, dan satuan rate per lini
bisnis.

## Context

Aplikasi Facultative Inward sedang dipindahkan dari Pega ke sistem baru. Kami merekam perilaku
sistem lama dari ekspor rule-nya, tetapi **korpus tidak memuat satu baris data pun** — tidak ada
DDL, tidak ada isi tabel. Dua hal karena itu hanya dapat dijawab dengan melihat data produksi.

Yang kami minta di sini adalah **agregat dan sampel kecil**, bukan salinan data. Tidak ada data
pelanggan yang perlu keluar dari basis data: pengukuran pertama hanya butuh hitungan, pengukuran
kedua hanya butuh rentang nilai numerik tanpa pengenal apa pun.

Dokumen ini terpisah dari permintaan besar di `_EKSTRAKSI-PEGA-SELAGI-HIDUP.md` (DDL, `ALL_SOURCE`,
isi tabel limit, ekstraksi riwayat akseptasi). Dua butir di sini **lebih kecil dan lebih cepat**, dan
tidak menunggu keputusan apa pun.

## How to answer

Tenggat: **23 September 2026**. Perkiraan waktu 30–60 menit, tergantung ketersediaan akses baca.

Tuliskan angkanya langsung di bawah tiap permintaan. Bila sebuah kolom yang kami asumsikan ternyata
tidak ada atau bernama lain, **katakan nama sebenarnya** — itu sendiri sudah jawaban yang berguna.

📋 **Daftar kolom persis per tabel ada di Lampiran** di akhir dokumen, ditandai mana yang diminta
datanya dan mana yang hanya konteks.

---

## Pengukuran 1 — nilai uang yang tersimpan tanpa mata uang

### 1. Berapa banyak baris bernilai uang yang tersimpan **tanpa mata uang**?

Kami menemukan bahwa **112 layar menampilkan nilai uang tanpa field mata uang mana pun**, dan hanya
ada 93 tempat di seluruh aplikasi yang mengikat properti mata uang. Sistem lama juga **tidak pernah**
menetapkan mata uang default di kode — ia selalu datang dari data.

Sistem baru mewajibkan setiap nilai uang membawa mata uangnya, dengan keadaan `Tidak Diketahui` yang
eksplisit untuk yang tidak membawanya. Yang perlu kami ketahui: **seberapa besar masalahnya.**

Berikut tabel yang kami maksud, beserta kolom mata uang yang **berhasil kami baca dari SQL sistem
lama**. Untuk tiap tabel: **berapa baris yang kolom mata uangnya `NULL` atau string kosong, dibanding
total baris** — dan bila memungkinkan, sebarannya per tahun.

| Tabel | Kolom uang/rate terbaca | Kolom mata uang terbaca | Yang kami minta |
| --- | ---: | --- | --- |
| `FACINPRODUCTION` | 21 dari 79 kolom | `CURR_ID` | hitungan `NULL`/kosong vs total |
| `FACINOFFER` | 7 dari 50 kolom | `CURR_ID` | hitungan `NULL`/kosong vs total |
| `FACOUTPRODUCTION` | 17 dari 65 kolom | **`CURRENCY` dan `CURRENCYID`** — dua kolom | lihat pertanyaan 1b |
| `FACINLIFE` | 8 dari 28 kolom | `CURRENCY` | hitungan `NULL`/kosong vs total |
| `FACINSPREADLIFE` | 11 dari 18 kolom | **tidak ada** | lihat pertanyaan 1c |
| `JSON_POLIS` | — | tidak terbaca | lihat pertanyaan 1d |

_Why this matters: bila angkanya kecil dan hanya di data lama, sistem baru cukup menandainya dan
jalan terus. Bila besar dan masih terjadi, berarti ada jalur input yang tidak mewajibkan mata uang —
dan itu harus diperbaiki di sumbernya sebelum cutover, bukan ditambal di sistem baru._

>

### 1b. `FACOUTPRODUCTION` punya **dua** kolom mata uang — mana yang berlaku?

Kami membaca dua kolom berbeda pada tabel yang sama: `CURRENCY` dan `CURRENCYID`.

Mana yang menjadi acuan, apakah keduanya selalu konsisten, dan bila tidak — mana yang dipercaya oleh
proses hilir?

_Why this matters: bila keduanya bisa berbeda, sistem baru harus tahu mana yang benar sebelum menulis
satu baris pun. Menyalin keduanya apa adanya justru mewariskan ambiguitasnya._

>

### 1c. `FACINSPREADLIFE` menyimpan 11 kolom nilai uang tetapi **tidak punya kolom mata uang sama sekali**

Ini temuan yang paling kami khawatirkan. Tabel ini menyimpan premi, TSI, diskon, dan komisi hasil
spreading — tetapi tidak ada satu pun kolom yang menyatakan mata uangnya.

Tiga kemungkinan, dan kami tidak dapat membedakannya dari kode:

1. Mata uangnya **diasumsikan selalu sama** (mis. selalu IDR) — bila ya, atas dasar apa?
2. Mata uangnya **diturunkan dari tabel lain** lewat kunci — bila ya, kunci dan tabel mana?
3. Informasi itu memang **hilang** di tabel ini.

_Why this matters: bila jawabannya (3), maka setiap nilai di tabel ini secara teknis tidak
terinterpretasi — dan sistem baru tidak boleh menuliskannya tanpa keputusan eksplisit. Bila (1) atau
(2), kami perlu aturannya agar dapat direproduksi._

>

### 1d. Kolom `JSON_POLIS` tidak dapat kami baca — bisakah Anda menjelaskan bentuknya?

Berbeda dari tabel lain, `JSON_POLIS` **tidak pernah ditulis lewat `INSERT` biasa** di kode yang kami
miliki — seluruh penulisannya melewati stored procedure, yang isinya tidak ada di korpus. Akibatnya
kami tahu tabel ini dipakai 23 rule, tetapi tidak tahu satu pun nama kolomnya selain yang muncul di
klausa `WHERE`.

Yang kami minta: **daftar kolomnya** (DDL cukup), dan khususnya — apakah dokumen JSON di dalamnya
membawa mata uang per nilai, atau hanya satu mata uang untuk seluruh dokumen?

_Why this matters: agregat kasus diserialkan utuh ke tabel ini, jadi ia adalah sumber kebenaran data
polis. Tanpa mengetahui bentuknya, sistem baru tidak dapat membaca data historis sama sekali._

>

---

## Pengukuran 2 — besaran nilai rate

### 2. Berapa rentang nilai `RATE`, `MIN_RATE`, dan `MAX_RATE` yang sebenarnya tersimpan?

Kami perlu memastikan satuan rate per lini bisnis — per mille (‰) atau persen (%). Kode memberi
jawaban yang berbeda-beda per lini, dan tiga label layar bertentangan dengan rumusnya.

**Besaran angkanya membuktikan satuannya secara langsung**, tanpa perlu pendapat: rate tipikal
sekitar 2,5 menunjukkan per mille; sekitar 0,25 menunjukkan persen.

Yang kami minta — **nilai minimum, maksimum, median, dan beberapa nilai tipikal** dari tiap kolom di
bawah. Tidak perlu pengenal polis, nama tertanggung, atau kolom lain: **hanya angka rate-nya**.

**(a) Rate yang benar-benar ditransaksikan** — ini yang paling menentukan, karena inilah angka yang
dipakai menghitung premi. Mohon **dikelompokkan per lini bisnis**, khususnya **MBU**, **PA**, dan
**Fire** sebagai pembanding:

| Tabel | Kolom rate terbaca |
| --- | --- |
| `FACINPRODUCTION` | `RATE` |
| `FACINOFFER` | `RATE` |
| `FACOUTPRODUCTION` | `RATE`, `RATE_COVERAGE` |
| `FACINLIFE` | `RATE` |
| `FACINSPREADLIFE` | `RATE_RETRO` |

**(b) Tabel tarif acuan** — untuk membandingkan apakah rate transaksi memang berasal dari tarif ini
dalam satuan yang sama:

`M_FLEXAS_RATE` (kolom `RATE`, `MIN_RATE`, `MAX_RATE` — tarif kebakaran) · `M_EQS_RATE` ·
`M_FLOOD_RATE` · `M_RSMD_RATE` · `M_TERORISME_RATE` · `RATE_LIFE`

**(c)** Bila ada kolom rate lain di tabel-tabel itu yang tidak kami sebut, mohon disebutkan.

_Why this matters: ini satu-satunya butir dalam seluruh daftar pertanyaan kami yang dapat ditutup
oleh data saja, tanpa menunggu pendapat tim mana pun. Bila datanya jelas, pertanyaan kami ke
Product/Aktuaria menjadi sekadar konfirmasi._

>

---

## Lampiran — daftar kolom per tabel

Daftar ini kami susun dari SQL sistem lama, supaya jelas **apa yang kami minta dan apa yang tidak**.

Tiga kategori, dan hanya dua yang berupa permintaan:

| Kategori | Status |
| --- | --- |
| 🔵 **Kolom mata uang** | **Diminta** — hitungan `NULL`/kosong vs total baris (Pengukuran 1) |
| 🟢 **Kolom rate** | **Diminta** — min / maks / median / beberapa nilai tipikal (Pengukuran 2) |
| ⚪ **Kolom uang lain** | **Tidak diminta datanya** — didaftar hanya agar terlihat apa yang kehilangan mata uang bila kolom 🔵 kosong |

⚠️ **Tidak ada nilai baris per-polis yang kami minta.** Seluruh permintaan berupa hitungan dan
statistik agregat. Kolom ⚪ tidak perlu dikeluarkan sama sekali.

---

### `FACINPRODUCTION` — 79 kolom, 21 di antaranya bernilai uang/rate

- 🔵 **`CURR_ID`**
- 🟢 **`RATE`**
- ⚪ `TSI_MENJADI` · `TSI_SELISIH` · `TSI100_MENJADI` · `TSI100_SELISIH` · `PREMI_MENJADI` ·
  `PREMI_SELISIH` · `RICOMM` · `RICOMM_SELISIH` · `PERCENT_RI_COMM` · `PCT_RI_COMM_SELISIH` ·
  `BROKERAGE_FEE_MENJADI` · `BROKERAGE_FEE_SELISIH` · `PCT_BROKERGARE_FEE` ·
  `PCT_BROKERAGE_FEE_SELISIH` · `DEDUCTION2_MENJADI` · `DEDUCTION2_SELISIH` · `LOL_MENJADI` ·
  `LOL_SELISIH` · `PERCENT_SHARE_SPREADED` · `PRORATE`

⚠️ `PCT_BROKERGARE_FEE` kami tulis **apa adanya** — ejaan itu memang begitu di sistem lama, sementara
pasangannya `PCT_BROKERAGE_FEE_SELISIH` dieja benar. Mohon konfirmasi apakah nama kolom di basis data
memang demikian; bila ya, kami pertahankan dan catat sebagai kandidat perbaikan terpisah.

### `FACINOFFER` — 50 kolom, 7 bernilai uang/rate

- 🔵 **`CURR_ID`**
- 🟢 **`RATE`**
- ⚪ `TSI_MENJADI` · `TSI_SELISIH` · `PREMI_MENJADI` · `PREMI_SELISIH` · `PERCENT_RI_COMM` ·
  `PERCENT_SHARE_SPREADED`

### `FACOUTPRODUCTION` — 65 kolom, 17 bernilai uang/rate

- 🔵 **`CURRENCY`** dan **`CURRENCYID`** — dua kolom; lihat pertanyaan 1b
- 🟢 **`RATE`** dan **`RATE_COVERAGE`**
- ⚪ `TSIRNM` · `TSISPREADED` · `OBJECTPREMI` · `OBJECTPREMI_SELISIH` · `PREMI_COVERAGE_MENJADI` ·
  `PREMI_COVERAGE_SELISIH` · `COMMISION` · `COMMISION_SELISIH` · `COMMISION_COVERAGE_MENJADI` ·
  `COMMISION_COVERAGE_SELISIH` · `COMMISION_COVERAGE_PCT` · `RICOMM` · `SHAREOFFERED` ·
  `SHAREOFFERED_SELISIH` · `PRORATE`

### `FACINLIFE` — 28 kolom, 8 bernilai uang/rate

- 🔵 **`CURRENCY`**
- 🟢 **`RATE`**
- ⚪ `TSI` · `TSI_RNM` · `TSI_CEDING` · `TSI_LIABILITY` · `PREMIUM` · `COMM` · `PCT_COMM`

### `FACINSPREADLIFE` — 18 kolom, 11 bernilai uang/rate, **nol kolom mata uang**

- 🔵 **tidak ada** — inilah inti pertanyaan 1c
- 🟢 **`RATE_RETRO`**
- ⚪ `TSI_SPREADED` · `TSI_LIABILITY` · `PREMIUM_SPREADED` · `NET_PREMI` · `TREATY_SHARE` ·
  `PCT_SHARE` · `OVR_COMM` · `PCT_OVR_COMM` · `DISCOUNT` · `PCT_DISCOUNT`

### `JSON_POLIS` — kolom tidak terbaca

Tidak ada satu pun nama kolom yang dapat kami baca, karena tabel ini **tidak pernah ditulis lewat
`INSERT` biasa** — seluruh penulisannya melewati stored procedure. Yang kami minta di sini adalah
**DDL-nya**, bukan datanya. Lihat pertanyaan 1d.

### Tabel tarif acuan — untuk Pengukuran 2

- 🟢 `M_FLEXAS_RATE` → `RATE`, `MIN_RATE`, `MAX_RATE`
- 🟢 `M_EQS_RATE` · `M_FLOOD_RATE` · `M_RSMD_RATE` · `M_TERORISME_RATE` · `RATE_LIFE` → kolom rate
  masing-masing (mohon sebutkan namanya bila berbeda dari `RATE`)

---

## Anything else?

Dua hal yang membantu bila Anda mengetahuinya:

- Apakah kolom rate atau kolom uang di tabel-tabel itu bertipe **numerik** atau **teks**? Bila teks,
  seluruh perbandingan angka di atasnya sebenarnya perbandingan string — dan itu mengubah cara kami
  mereproduksi perhitungannya.
- Apakah pernah ada perubahan satuan atau format desimal pada kolom-kolom ini di masa lalu, sehingga
  data lama dan data baru tidak sebanding?

>
