> Modul  : Treaty In Adjustment · ronde TDA
> Dibuat : 2026-09-24
> Peran  : ditulis dengan peran penilai. Diputuskan pemilik proses pada 24 September 2026.
> Masukan: `00-LINGKUP.md`, `02-SIDANG.md` (ST-1…ST-7), `PROPERTI-ASLI.md`, ADR induk 0034, 0037,
>          0042, 0044, 0045, 0055, `SPEC-INVARIAN.md` INV-24 dan INV-25, GRL-01…GRL-17
> Status : **TERKUNCI**
> Sifat  : grilling. Tidak ada spesifikasi, DDL, struct Go, komponen React, endpoint, atau tiket.

# Putusan ronde TDA

| # | Judul | Label | Pembatalnya |
|---|---|---|---|
| **GRL-18** | Jenis addendum adalah **masukan**, dapat dibetulkan selama `DRAFT`, dan **beku sejak diajukan** | pilihan orang PELESTARIAN · beku + berjejak PERUBAHAN | `DB-16` |
| — | **Nasib tujuh TDA** ditetapkan; enam tertutup dengan kutipan, satu ditutup oleh GRL-18 | — | — |

---

## GRL-18 — jenis addendum: siapa menulisnya, dan sampai kapan

### Keputusan

Satu kalimat, dan ketiga bagiannya mengikat:

> Jenis sebuah versi **dinyatakan pengisi kontrak** saat versi dibuat; **boleh dibetulkan selama
> versinya `DRAFT`**, dan tiap perubahan meninggalkan baris jejak; **membeku begitu diajukan**.
> Mengubahnya sesudah pengajuan menuntut versinya dikembalikan ke `DRAFT` lebih dulu.

### Bagian 1 — masukan, bukan turunan

**Dasar: ADR-0037 butir 3** — *"Bila sebuah besaran boleh disepakati, ia dinaikkan menjadi
masukan."*

Aturan yang **sama** dipakai `GRL-12` untuk menolak materialitas sebagai masukan, dan ia memberi
jawaban berlawanan di sini. Itu bukan dua keputusan yang kebetulan berbeda:

| Sumbu | Disepakati? | Maka |
|---|---|---|
| materialitas | **tidak** — ia sifat dari apa yang berubah | **turunan** (`GRL-12`) |
| jenis addendum | **ya** — penyesuaian premi adalah perbuatan bisnis yang disepakati dengan cedant | **masukan** |

**Yang mematahkan pilihan "turunan", dan ia skenario, bukan penalaran:**

> Dua versi menaikkan `TreatyIn.EGNPI` sebesar angka yang persis sama. Yang pertama karena cedant
> dan NuRe **menyepakati** penyesuaian premi dari estimasi ke aktual. Yang kedua karena pengisi
> **salah ketik tahun lalu** dan sekarang membetulkannya.
>
> Datanya **identik**. Turunan tidak punya cara memisahkan keduanya — dan tidak ada tempat untuk
> menyatakan bahwa yang kedua bukan kesepakatan.

Ini bentuk yang sama dengan *"uji yang dijawab sama oleh kedua kemungkinan bukan uji"*, hanya
terjadi pada **rumus penurunan** alih-alih pada uji data.

### Bagian 2 — beku sejak `AJUKAN`, bukan sejak `DISETUJUI`

**Dasar: penyetuju harus menyetujui apa yang ia lihat.**

> Bila jenis masih dapat berubah selagi versi menunggu tanda tangan, **tiga tingkat penyetuju dapat
> menyetujui tiga jenis yang berbeda, tanpa satu pun galat, dan tanda tangan pertama tetap tercatat
> sah.**

Itu persis yang sistem lama izinkan, dan itu yang tidak dibawa.

Pembekuan sejak lahir **ditolak**, dan sebabnya bukan bahwa ia salah melainkan bahwa ia **menghukum
tanpa melindungi siapa pun**: pengisi yang membetulkan jenisnya sebelum mengajukan tidak merugikan
satu pihak pun, sebab belum ada satu mata pun yang melihatnya. Ongkosnya nyata — membuang seluruh
isian dan mengulang.

**Tidak ada mekanisme baru yang diminta.** Jalan untuk mengubahnya sesudah pengajuan sudah ada:
`KEMBALIKAN` ke `DRAFT`, perpindahan yang sudah terdaftar ADR-0055 dan sudah berlabel PELESTARIAN di
`GRL-08`.

### Bagian 3 — tiap perubahan berjejak

**Dasar: ADR-0045** — jejak perubahan sebagai fakta mesin, tidak dapat dilewati.

Sistem lama mengubahnya **tanpa jejak sama sekali**: radio menimpa nilainya, tidak ada yang
mencatat nilai sebelumnya, dan tidak ada pemeriksaan di sisi simpan. Maka pertanyaan *"jenis apa
versi ini saat kepala seksi menyetujuinya"* **tidak dapat dijawab** untuk satu pun baris warisan.

### Dalam nama asli — dan nama aslinya tidak diubah

| Properti asli | Sistem lama | Sesudah GRL-18 |
|---|---|---|
| `TreatyIn.EDMState` | disemai tombol (Revisi → `1`, Penyesuaian → `3`), **ditimpa radio picker tanpa syarat apa pun**, kapan saja, nol pemeriksaan di sisi simpan | masukan bernilai **dua** (`GRL-13`); dapat diubah **hanya** selama `DRAFT`; berjejak; beku sejak `AJUKAN` |
| `TreatyIn.EDMMaterialType` | radio kedua, bebas dipasangkan dengan nilai mana pun | **hilang, tanpa pengganti** — materialitas turunan (`GRL-12`) |
| `TreatyIn.AddendumPremi` | bendera yang disetel `TreatyInSetEditPre` untuk `EDMState = 3` | **hilang** — ia turunan dari jenis, bukan fakta tersendiri |

Nilai `TreatyIn.EDMState` **warisan** dibawa apa adanya (`GRL-13`, ADR-0042), termasuk kombinasi
yang tidak lagi sah untuk versi baru.

### Label

| Bagian | Label | Berubah dari |
|---|---|---|
| pengisi yang menyatakan jenis | **PELESTARIAN** | — sistem lama memang begitu |
| beku sejak diajukan | **PERUBAHAN** | *"dapat diubah kapan saja selama layarnya dapat disunting"* |
| tiap perubahan berjejak | **PERUBAHAN** | *"ditimpa tanpa jejak, tanpa nilai sebelumnya, tanpa pemeriksaan"* |

### Bahan to-spec

| # | Bunyi | Kasus yang dapat gagal |
|---|---|---|
| **T-1** | Jenis sebuah versi hanya dapat berubah ketika versinya berkeadaan `DRAFT` | mengubah jenis pada versi yang menunggu persetujuan → **ditolak**, dan pesannya menyebut bahwa versinya harus dikembalikan dulu |
| **T-2** | Tiap perubahan jenis meninggalkan baris `JEJAK_PERUBAHAN` berisi nilai lama, nilai baru, siapa, dan kapan | mengubah jenis dua kali lalu kembali ke nilai semula → **tiga baris jejak**, bukan nol |
| **T-3** | Jenis yang tercatat pada sebuah persetujuan adalah jenis pada saat pengajuan | membaca versi yang sudah disetujui → jenisnya **sama** dengan yang dilihat penyetuju pertama |

### Pembatal — `DB-16`, yang sudah ada

Bila addendum eksternal ternyata **didefinisikan oleh dokumen bertanda tangan cedant**, jenisnya
melekat pada bukti, bukan pada pilihan. Yang berubah hanya **dari mana nilainya datang** — ia tetap
masukan, dan ketiga bagian GRL-18 tetap berlaku. Tidak perlu butir pembatal baru.

---

## Nasib tujuh TDA — ditetapkan

Isi lengkap beserta kutipan pemutusnya ada di [`02-SIDANG.md`](02-SIDANG.md). Di sini nasibnya saja,
dalam kosakata `METODE` §6.3:

| TDA | Nasib | Oleh |
|---|---|---|
| `TDA-01` | **diperbaiki** · sisa warisan **ditunda** | `B-1` + `GRL-10`; sisa diukur `UA-14`, `UA-21` |
| `TDA-07` | **diperbaiki** | `GRL-08` + INV-24 |
| `TDA-09` | **diperbaiki** (KONFIRMASI ADR-0044) | ADR-0044 |
| `TDA-11` | **diperbaiki** untuk model · sisa tampilan **ditunda** | `GRL-10` + `B-4` + INV-25; sisa ke REKOMENDASI TO-SPEC |
| `TDA-12` | **diperbaiki** | `GRL-09` |
| `TDA-13` | **diperbaiki** | `GRL-12` (materialitas) + **`GRL-18`** (jenis) |
| `TDA-15` | **diperbaiki** | ADR-0034 + `GRL-03` + `GRL-14` |

**Lima belas TDA, lima belas nasib.** Tidak ada lagi yang menunggu cabang yang tidak bersidang.

> Enam dari tujuh **tidak menuntut pekerjaan baru.** Keenamnya sudah tertutup berhari-hari oleh
> putusan yang sudah terkunci, dan yang hilang hanya barisnya. Itu yang membuat ronde ini murah —
> dan yang membuat ketiadaannya mahal.
