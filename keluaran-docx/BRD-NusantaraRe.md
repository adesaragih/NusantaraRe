# Dokumen Kebutuhan Bisnis - Nusantara Re

> Migrasi sistem reasuransi dari Pega ke Go + React, dengan Oracle dipertahankan.
> Disusun 25 September 2026 dari spesifikasi, tiket, keputusan arsitektur, dan catatan
> penelusuran proyek. **Nol fakta baru** - seluruh isi berasal dari dokumen yang sudah ada.

## Cara membaca dokumen ini

Bab 1-3 memberi latar, lingkup, dan peta konteks. Bab 4-23 memuat **satu bab per modul**.
Bab 24 adalah yang membuat dokumen ini jujur tanpa perlu ditunda: ia menyatakan **per modul**
apa yang sudah pasti dari korpus dan apa yang masih menunggu pihak lain, beserta pihaknya.
Bab 25-28 memuat model data ringkas, keputusan arsitektur yang paling mengikat, butir terbuka
per pemilik peran, dan risiko yang diterima sadar.

| Penanda | Arti |
| --- | --- |
| `[terverifikasi]` | terbaca dari berkas ekspor sistem lama, dengan perintah ujinya |
| `[keputusan work owner]` | diputuskan pihak yang berwenang - tidak diubah dokumen ini |
| `[penyimpangan sadar]` | sengaja berbeda dari sistem lama, alasannya tertulis |
| `[terbuka]` | belum terjawab - **tidak ditutup** dokumen ini |
| `[data DBA]` | hanya dapat dijawab dari basis data |

---

# 1. Latar dan tujuan

Nusantara Re menjalankan proses reasuransinya di atas **Pega**. Aturan bisnis, tata letak layar,
alur persetujuan, dan sebagian besar perhitungan uang tinggal di dalam aturan Pega, sementara
datanya tinggal di **Oracle** - sebagian sebagai kolom, sebagian sebagai **dokumen teks** di
dalam satu kolom.

**Tujuan migrasi:** memindahkan perilaku sistem ke **Go** di sisi layanan dan **React** di sisi
layar, dengan **Oracle dipertahankan** sebagai basis data.

Tiga ketetapan yang mengikat seluruh pekerjaan:

| # | Ketetapan | Akibatnya |
| --- | --- | --- |
| 1 | **Nol pemanggilan stored procedure** dari sistem baru | logika yang selama ini tinggal di dalam program basis data **ditulis ulang** di lapisan layanan |
| 2 | **Uang tidak pernah bilangan mengambang** | nilai uang memakai tipe desimal berpresisi tetap di seluruh jalur |
| 3 | Arah ketergantungan `handlers -> services -> repository` | lapisan penyimpanan tidak pernah memanggil lapisan layanan |

**Yang tidak berubah:** basis data Oracle, nomor polis yang sudah terbit, dan angka historis -
termasuk galat presisi yang sudah tersimpan di dalamnya, supaya laporan lama tetap dapat
direkonsiliasi.

---

# 2. Lingkup

Pekerjaan mencakup **20 modul** yang dikelompokkan ke dalam **sembilan bounded context**.
Bahannya lahir dari **tiga rangkaian kerja** yang berjalan berdampingan, dan ketiganya dipakai
dokumen ini.

| Rangkaian | Modul | Tiket | Keputusan arsitektur |
| --- | ---: | ---: | ---: |
| utama | 13 | 177 | 42 |
| kedua | 4 | 117 | 56 |
| Fakultatif | 3 *(ditambah satu menu tersendiri)* | 60 | 7 |
| **Jumlah** | **20** | **354** | **105** |

[~] **Cacah tiket ini hasil pengukuran ulang.** Catatan audit sebelumnya menyebut **365**;
selisih **11** seluruhnya adalah **berkas indeks** yang ikut terhitung sebagai tiket. Cacah yang
dipakai di sini adalah **354 berkas tiket sebenarnya**.

## Modul per konteks

| Bounded context | Modul | Tiket |
| --- | --- | ---: |
| **Klaim Jiwa** | Claim Life | 15 |
| **Penawaran Jiwa** | PremiumList Life . Endorsement Life | 21 |
| **Data Induk Jiwa** | Master Product Name Life . Master Contract Retro Life | 22 |
| **Klaim Non-Jiwa** | Claim Fac In . Claim Prop . Claim Non Prop | 68 |
| **Komite Klaim** | Komite Claim Fac In . Komite Claim Life . Komite Claim Prop . Komite Claim Non Prop | 50 |
| **Realisasi Treaty** | NB Treaty In . EDM Treaty In | 40 |
| **Kontrak Treaty** | Treaty In . Treaty In Adjustment | 66 |
| **Penempatan Treaty** | Treaty Contract Out | 12 |
| **Fakultatif** | NB Fac In . RNW Fac In . Endorsment Fac In | 46 |

---

# 3. Peta konteks

Sembilan bounded context terbaca dari korpus, ditambah satu yang **tidak dapat diturunkan**
darinya dan ditandai demikian, bukan diisi tebakan.

| # | Bounded context | Modul | File | Cakupan bukti |
| ---: | --- | --- | ---: | --- |
| 1 | **Facultative Inward** | NB FacIn, RNW Fac In, Endorsment Fac In | 6.071 | **partial** |
| 2 | **Treaty Inward — Master & Akseptasi** | Treaty In, Treaty In Adjustment | 708 | **partial** |
| 3 | **Treaty Inward — Realisasi & Endorsement** | NB Treaty In, EDM Treaty In | 441 | **full** |
| 4 | **Treaty Arrangement (master term)** | Treaty Contract Out | 303 | **partial** |
| 5 | **Claim — Non-Life** | Claim Fac In, Claim Prop, Claim Non Prop | 1.031 | **partial** |
| 6 | **Claim — Life** | Claim Life | 136 | **full** |
| 7 | **Komite (tangga persetujuan klaim)** | Komite Claim FacIn, Prop, Non Prop, Life | 300 | **full** |
| 8 | **Life — Penawaran & Premium List** | PremiumList Life, Endorsement Life | 199 | **full** |
| 9 | **Life — Master** | Master Product Name Life, Master Contract Retro Life | 180 | **partial** |

**Kebocoran batas yang sudah terbukti** dicatat di berkas peta konteks proyek; yang paling
menyentuh pekerjaan ini adalah **tabel akar yang dipakai bersama lebih dari satu lini** - lihat
bab 26 dan dokumen Keputusan Arsitektur, Lampiran 2 Kelompok 2.

---

# 4. Claim Life

*Bounded context: Klaim Jiwa. Bahan dari rangkaian utama.*

| Ukuran | Nilai |
| --- | ---: |
| Tiket | **15** |
| User story pada spec | 60 |
| Pernyataan terverifikasi dari korpus | 75 |
| Keputusan pemilik pekerjaan | 72 |
| Butir terbuka | 23 |
| Penyimpangan sadar | 0 |

## Masalah

Penanganan klaim jiwa (*life*) Nusantara Re hari ini berjalan di Pega. Tiga peran — admin klaim, penasihat medis, dan supervisor — mendaftarkan klaim, menelaah sisi medis, lalu memutuskan akseptasi, dengan Komite sebagai tangga persetujuan di luar sistem ini.

Masalah yang dihadapi pengguna bila status quo dipertahankan:

## Solusi

Membangun konteks **Claim — Life** sebagai layanan Go + antarmuka React yang menulis ke Oracle `POOLDATA` yang sama, dengan:

Perilaku bisnis dibawa apa adanya (paritas), kecuali empat penyimpangan sadar yang sudah diputuskan: jejak audit (**ADR-U-0007**), antre-ulang (**ADR-U-0008**), muatan kontrak Komite yang diperluas (**ADR-U-0001**), dan representasi uang (**ADR-U-0003**).

## Kebutuhan pengguna

Spesifikasi modul ini memuat **60 user story**, dimuat utuh di berkas spec-nya.
Yang dijanjikannya, diringkas: pekerjaan yang hari ini dikerjakan di dalam layar Pega
dapat dikerjakan di sistem baru dengan hasil yang sama, jejak keputusan yang dapat
ditelusuri, dan angka yang dapat direkonsiliasi dengan riwayat.

## Keputusan yang mengikat

- 1. Batas konteks
- 2. Penempatan modul
- 2b. Bentuk penyimpanan — delapan tabel relasional + T_WORK_CLAIM
- RALAT 2026-09-18 · A — T_CLAIM_POLICY dan T_CLAIM_MARKETING DIHAPUS
- RALAT 2026-09-18 · B — kolom yang bertambah dan yang pindah
- RALAT 2026-09-18 · C — yang MASIH TERBUKA. Jangan dijawab sendiri.
- RALAT 2026-09-18 · D — pohon yang berlaku: enam tingkat, berakar di T_WORK_CLAIM
- RALAT 2026-09-18 · E — sebelas relasi

## Yang tidak dibangun

- **Isi dan tangga internal Komite Life.** Konteks luar (**ADR-U-0001**). Spec ini berhenti di tiga
- **Tiga modul Claim lain** (`Claim Prop`, `Claim Non Prop`, `Claim Fac In`). Pola `STS_REJECT` +
- **Merapikan alur.** Non-goal yang sudah disepakati (Ronde 1 Q2): urutan tahap, jalur balik, dan
- **Memperketat kelonggaran wewenang `TP`/`TR`.** Dibawa apa adanya; peninjauan saat Komite/IAM
- **Mengganti penyimpanan berkas.** Google Storage dipertahankan (**ADR-U-0010**).
- **Memindahkan penomoran ke aplikasi.** Tetap di Oracle (**ADR-U-0006**).
- **Identity & Access sebagai konteks.** `[terverifikasi]` `discovery/context-map.md` menandainya
- **Scaffolding kode.** `cmd/`, `internal/`, `pkg/`, `frontend/`, `go.mod`, `Makefile` **belum ada

## Butir terbuka

Modul ini membawa **23 butir terbuka** yang **tidak ditutup** dokumen ini.
- **OQ-020** (sisa)
- **OQ-032**
- **OQ-037**
- **OQ-060**
- **OQ-035**
- ~~**OQ-001**~~

*Rujukan: `.scratch\claim-life\spec.md`*

---

# 5. PremiumList Life

*Bounded context: Penawaran Jiwa. Bahan dari rangkaian utama.*

| Ukuran | Nilai |
| --- | ---: |
| Tiket | **9** |
| User story pada spec | 36 |
| Pernyataan terverifikasi dari korpus | 46 |
| Keputusan pemilik pekerjaan | 35 |
| Butir terbuka | 1 |
| Penyimpangan sadar | 0 |

## Masalah

**Premium List Life adalah hulu domain jiwa** — ia yang menghasilkan `PremiumListSummary` dan `PremiumListDetail` yang kemudian dikonsumsi Claim Life, Komite Claim Life, dan Endorsement. `[terverifikasi]` Class integrasi `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` dirujuk **46 rule** di Claim Life, 25 di Endorsement, 19 di PremiumList, 11 di Komite — ia tulang punggung domain Life.

Masalah yang dihadapi bila status quo dipertahankan:

## Solusi

Membangun konteks **Life — Penawaran & Premium List** sebagai layanan Go + antarmuka React yang menulis ke Oracle `POOLDATA` yang sama, dengan:

Perilaku dibawa apa adanya (paritas), kecuali **tiga penyimpangan sadar**: email diperlakukan sebagai **jalur alarm** bukan notifikasi bisnis; penulisan detail peserta **inline saat simpan, bukan job**; dan Go memegang batas transaksi menggantikan commit implisit Pega.

## Kebutuhan pengguna

Spesifikasi modul ini memuat **36 user story**, dimuat utuh di berkas spec-nya.
Yang dijanjikannya, diringkas: pekerjaan yang hari ini dikerjakan di dalam layar Pega
dapat dikerjakan di sistem baru dengan hasil yang sama, jejak keputusan yang dapat
ditelusuri, dan angka yang dapat direkonsiliasi dengan riwayat.

## Keputusan yang mengikat

- 1. Batas konteks
- 2. Penempatan modul
- 3. Mesin alur
- 4. Periode tutup buku
- 5. Uang — teks di batas, desimal di dalam
- 6. Batas transaksi — ⚠️ DIREVISI 2026-09-16: titik potong lenyap
- 7. Penomoran
- 8. Unggah CSV

## Yang tidak dibangun

- **Endorsement Life — SELURUHNYA.** `[keputusan work owner]` Konteks/menu terpisah: class
- **Claim Life dan Komite Claim Life.** Konteks hilir; spec ini berhenti pada penulisan rekam
- **Lini non-Life.** Treaty, facultative, dan claim non-life punya konteks sendiri.
- **Merapikan alur.** Urutan tahap dan bentuk keputusan dibawa apa adanya; tiga penyimpangan sadar
- **Membangun shape "Input Premium List Summary"** sebelum OQ-023 dijawab.
- **Identity & Access.** `[terverifikasi]` ABSENT dari korpus; dibangun dari nol.
- **Scaffolding kode.** Pekerjaan terpisah yang mendahului tiket mana pun.
- ~~**Migrasi data premium list lama.** Menunggu DDL fisik (OQ-001).~~ ⚠️ **TIDAK LAGI di luar

## Butir terbuka

Modul ini membawa **1 butir terbuka** yang **tidak ditutup** dokumen ini.
- ~~**OQ-001**~~ (sisa)
- **OQ-023**
- **OQ-066**
- **OQ-028**
- **OQ-067**
- ~~**OQ-068**~~

*Rujukan: `.scratch\premiumlist-life\spec.md`*

---

# 6. Endorsement Life

*Bounded context: Penawaran Jiwa. Bahan dari rangkaian utama.*

| Ukuran | Nilai |
| --- | ---: |
| Tiket | **12** |
| User story pada spec | 47 |
| Pernyataan terverifikasi dari korpus | 38 |
| Keputusan pemilik pekerjaan | 39 |
| Butir terbuka | 3 |
| Penyimpangan sadar | 0 |

## Masalah

Polis Life yang sudah berjalan **berubah**. Peserta bertambah atau keluar, datanya salah dan perlu dikoreksi, atau seluruh polis dibatalkan. Hari ini perubahan itu dikerjakan di Pega lewat menu Endorsement Life yang terpisah dari input polis baru — dan pengetahuan tentang cara kerjanya tersebar di 75 berkas rule yang sebagian sudah mati, sebagian menyimpan aturan penting di dalam deskripsi langkah yang **bertentangan dengan kodenya**.

Empat hal membuat konteks ini berisiko dipindahkan secara salah:

## Solusi

Membangun ulang Endorsement Life sebagai **menu tersendiri** di atas Go + React + Oracle, dengan perilaku dibawa apa adanya (paritas) kecuali **tujuh penyimpangan sadar** yang sudah diputus.

Alur inti yang **tidak ada** di konteks PremiumList Life:

## Kebutuhan pengguna

Spesifikasi modul ini memuat **47 user story**, dimuat utuh di berkas spec-nya.
Yang dijanjikannya, diringkas: pekerjaan yang hari ini dikerjakan di dalam layar Pega
dapat dikerjakan di sistem baru dengan hasil yang sama, jejak keputusan yang dapat
ditelusuri, dan angka yang dapat direkonsiliasi dengan riwayat.

## Keputusan yang mengikat

- 1. Batas konteks
- 2. Penempatan modul
- 3. Gerbang kelayakan — lima pemeriksaan, sebelum case dibuat
- 4. Memuat & menyalin data polis lama
- 5. EdmType dan EDMStatus — dua maksud, empat status
- 6. Jurnal balik — × -1, bukan nol
- 7. Penomoran PL_NUMBER_EDM
- 8. Mesin alur

## Yang tidak dibangun

- **PremiumList Life (new business) — seluruhnya.** Konteks/menu terpisah; spec dan tiketnya di
- **Claim Life dan Komite Claim Life.** Konteks hilir. ⚠️ **Kecuali** satu hal: syarat penyaringan
- **Lini non-Life.** Treaty, facultative, dan claim non-life punya konteks sendiri.
- **Mengubah pembacaan Arasapas menjadi layanan.** Perubahan kontrak dengan pihak lain; di luar
- **Merapikan alur.** Urutan tahap dan bentuk keputusan dibawa apa adanya; tujuh penyimpangan sadar
- **Identity & Access.** `[terverifikasi]` ABSENT dari korpus; dibangun dari nol.
- **Scaffolding kode.** Pekerjaan terpisah yang mendahului tiket mana pun.
- **Migrasi data endorsement lama.** Menunggu DDL fisik (OQ-001 sisa).

## Butir terbuka

Modul ini membawa **3 butir terbuka** yang **tidak ditutup** dokumen ini.
- ~~**OQ-001**~~ (sisa)
- **OQ-020** (sisa)
- `[terbuka]`
- **OQ-066**
- **OQ-028**

*Rujukan: `.scratch\endorsement-life\spec.md`*

---

# 7. Master Product Name Life

*Bounded context: Data Induk Jiwa. Bahan dari rangkaian utama.*

| Ukuran | Nilai |
| --- | ---: |
| Tiket | **9** |
| User story pada spec | 36 |
| Pernyataan terverifikasi dari korpus | 22 |
| Keputusan pemilik pekerjaan | 42 |
| Butir terbuka | 5 |
| Penyimpangan sadar | 0 |

## Masalah

Setiap produk asuransi jiwa yang ditanggung ulang punya **definisi master**: siapa cedingnya, siapa pemegang polisnya, mata uangnya, rentang usia yang diterima, batas uang pertanggungan, batas retensi ceding, share Nusantara Re, brokerage, dan syarat-syarat penerimaan lain. Hari ini definisi itu dikelola di satu layar Pega dan disimpan sebagai **satu dokumen JSON** ke dua tabel Oracle.

Lima hal membuat modul ini berisiko dipindahkan secara salah:

## Solusi

Membangun ulang Master Product Name Life sebagai **editor master produk life murni** di atas Go + React + Oracle — bukan proses berjenjang: tidak ada persetujuan, tidak ada Submit/Decline, tidak ada status akseptasi.

`[fakta bisnis — work owner]` **Satu produk**, disimpan ke **dua tabel** yang terhubung lewat **`ID` yang sama**. Bukan dua produk berbeda.

## Kebutuhan pengguna

Spesifikasi modul ini memuat **36 user story**, dimuat utuh di berkas spec-nya.
Yang dijanjikannya, diringkas: pekerjaan yang hari ini dikerjakan di dalam layar Pega
dapat dikerjakan di sistem baru dengan hasil yang sama, jejak keputusan yang dapat
ditelusuri, dan angka yang dapat direkonsiliasi dengan riwayat.

## Keputusan yang mengikat

- 1. Batas konteks
- 2. Penempatan modul
- 3. Entitas — satu produk, SATU tabel induk, lima tabel anak
- 4. Tujuh pemilih master
- 5. Simpan — satu transaksi atomik, tanpa procedure lama
- 6. Identitas dan hasil simpan
- 7. Skema relasional penuh — dan dari mana daftar kolomnya
- Aturan pemilihan kolom [keputusan work owner]

## Yang tidak dibangun

- Jalur
- **Seluruh jalur treaty inward** — `SetTreatyIn_Act` (`DATA-PORTAL` / `SETTREATYIN_ACT`), `SaveTreatyIn` (`ASM-FW-GISFW-INT-TREATY_IN` / `ASM!SAVETREATYIN`), `BrowseTreatyIn`, gerbang `revisionstate==1`, `POOLDATA.PEGA_TREATY_IN`, `M_TREATY_IN` / `M_TREATY_IN_EDM`
- **Jalur klaim** — `GetCountClaim` dan akses langsung ke `DATAPEGA.PC_ASM_FW_GCNMFW_WORK`
- **Jalur `OR`** — `GetReinsTypeOR_Life` (`ASM-FW-GISFW-INT-PRODUCT_LIFE` / `GETREINSTYPEOR_LIFE`), `BrowseReinstypeOR_SQL`, bendera `IsORS`
- **Jalur simpan pintas** — `SaveProductNameLIfeFlat` (`ASM-FW-GISFW-INT-PRODUCT_LIFE` / `ASM!SAVEPRODUCTNAMELIFEFLAT`), `UPDATE` parsial atas `RIRISKID`/`RIRISK`
- **Master Contract Retro Life — seluruhnya.** Konteks/menu terpisah, sudah selesai.
- **Claim Life, PremiumList Life, Endorsement Life.** Konteks hilir yang **membaca** master ini.
- **Master pendukung itu sendiri** — ceding, pemegang polis, mata uang, sumber bisnis, penyebab

## Butir terbuka

Modul ini membawa **5 butir terbuka** yang **tidak ditutup** dokumen ini.
- **OQ-002**
- **OQ-001** (modul ini)
- **OQ-JSON-PRODUCT**
- **OQ-056**
- **OQ-057**
- **OQ-058**

*Rujukan: `.scratch\master-product-name-life\spec.md`*

---

# 8. Master Contract Retro Life

*Bounded context: Data Induk Jiwa. Bahan dari rangkaian utama.*

| Ukuran | Nilai |
| --- | ---: |
| Tiket | **13** |
| User story pada spec | 42 |
| Pernyataan terverifikasi dari korpus | 32 |
| Keputusan pemilik pekerjaan | 42 |
| Butir terbuka | 6 |
| Penyimpangan sadar | 0 |

## Masalah

Retrosesi lini life bersandar pada **data master** yang menentukan berapa risiko ditahan dan berapa yang dilepas: kontrak treaty per tahun, batas proteksi per mata uang, siapa reinsurernya dan berapa persen bagiannya, siapa yang meretrosesi bagian itu lagi, dan jenis bisnis apa yang tercakup. Hari ini seluruhnya dikelola di empat layar grid Pega di atas lima tabel Oracle.

Lima hal membuat modul ini berisiko dipindahkan secara salah:

## Solusi

Membangun ulang Master Contract Retro Life sebagai **menu tersendiri** di atas Go + React + Oracle, dengan perilaku dibawa apa adanya (paritas) kecuali **delapan penyimpangan sadar** yang sudah diputus.

Modul ini **bukan proses berjenjang**: tidak ada persetujuan, tidak ada Submit/Decline, tidak ada status akseptasi. Ia **editor master murni** — buka grid, tambah/ubah/hapus baris, simpan.

## Kebutuhan pengguna

Spesifikasi modul ini memuat **42 user story**, dimuat utuh di berkas spec-nya.
Yang dijanjikannya, diringkas: pekerjaan yang hari ini dikerjakan di dalam layar Pega
dapat dikerjakan di sistem baru dengan hasil yang sama, jejak keputusan yang dapat
ditelusuri, dan angka yang dapat direkonsiliasi dengan riwayat.

## Keputusan yang mengikat

- 1. Batas konteks
- 2. Penempatan modul
- 3. Entitas dan hierarki
- 4. Empat grid
- 5. Batas proteksi dan lebar layer
- 6. Dua tingkat share dan eksposur berjenjang
- 7. Total share dan validasi
- 8. Gerbang konsistensi tahun

## Yang tidak dibangun

- **Master Product Name Life — seluruhnya.** `[keputusan work owner]` Konteks/menu terpisah (114
- **Claim Life, Komite Claim Life.** Konteks hilir yang **membaca** master ini; spec ini berhenti
- **Lini non-Life.** Treaty, facultative, dan claim non-life punya konteks sendiri — termasuk
- **Menulis ulang aturan penyimpanan ke Go.** `[keputusan work owner]` Procedure dipanggil apa
- **Mengubah tipe `RIRATE`.** `[data DBA]` Ia `VARCHAR2(1000)`; mengubahnya menjadi angka menuntut
- **Identity & Access.** `[terverifikasi]` Modul ini **tanpa jejak guard identitas sama sekali** —
- **Efek keluar dan flag lingkungan.** `[terverifikasi]` Modul ini tidak punya integrasi luar;
- **Scaffolding kode.** Pekerjaan terpisah yang mendahului tiket mana pun.

## Butir terbuka

Modul ini membawa **6 butir terbuka** yang **tidak ditutup** dokumen ini.
- **OQ-002**
- **OQ-001**
- **OQ-057**
- **OQ kecil**
- **OQ-054**
- **OQ-066** (keluarga)

*Rujukan: `.scratch\master-contract-retro-life\spec.md`*

---

# 9. Claim Fac In

*Bounded context: Klaim Non-Jiwa. Bahan dari rangkaian utama.*

| Ukuran | Nilai |
| --- | ---: |
| Tiket | **15** |
| User story pada spec | 46 |
| Acceptance criteria | 112 |
| Pernyataan terverifikasi dari korpus | 123 |
| Keputusan pemilik pekerjaan | 41 |
| Butir terbuka | 38 |
| Penyimpangan sadar | 19 |

## Masalah

Klaim **fakultatif masuk** *(facultative inward)* Nusantara Re hari ini berjalan di atas Pega — **482 rule**, **179 activity**, **2.509 langkah**, dan **satu berkas alur** yang memegang seluruh daur hidup kasus. Yang dirasakan orang yang bekerja dengannya:

1. **Aturan bisnisnya tidak terbaca dari layar.** Yang menentukan sebuah klaim boleh dikirim ke komite atau tidak adalah sebuah penanda bernama `IsError` yang bernilai **2**, dan arti angka itu **tidak tertulis sebagai aturan di mana pun** — ia hanya ada di **catatan bebas seorang pengembang** pada rule lain. 2. **Angka uang tidak dijaga tipenya.** Total penyesuaian dideklarasikan sebagai bilangan pecahan biner, tiga nilai uang lain lewat sebagai teks, dan satu di antaranya harus **ditambal koma-ke-titik** di dalam SQL sebelum dapat dihitung. 3. **Nama tidak dapat dipercaya.** Kolom hasil SQL dinamai hal yang bukan isinya — tanggal ke...

## Solusi

Membangun ulang Claim Fac In sebagai **Go + React + Oracle**, dengan perilaku Pega **ditiru apa adanya** kecuali pada **dua belas titik yang sengaja diubah** *(Bab 17)*.

Yang berubah bagi orang yang memakainya:

## Kebutuhan pengguna

Spesifikasi modul ini memuat **46 user story**, dimuat utuh di berkas spec-nya.
Yang dijanjikannya, diringkas: pekerjaan yang hari ini dikerjakan di dalam layar Pega
dapat dikerjakan di sistem baru dengan hasil yang sama, jejak keputusan yang dapat
ditelusuri, dan angka yang dapat direkonsiliasi dengan riwayat.

## Keputusan yang mengikat

- Bab 1 — Lingkup, kelas kerja, dan identitas rule
- Identitas rule EMPAT bagian
- Ruleset ADESAMUEL@
- Tiga sistem sumber
- Bab 2 — Daur hidup kasus
- Tiket: dua bernama, lima cangkang kosong
- Bab 3 — Percabangan per lini produk
- Bab 4 — Klasifikasi lini dan Quotation

## Butir terbuka

Modul ini membawa **38 butir terbuka** yang **tidak ditutup** dokumen ini.
- **1**
- **2**
- **3**
- **4**
- **5**
- ⭐ **6**

*Rujukan: `.scratch\claim-facin\spec.md`*

---

# 10. Claim Prop

*Bounded context: Klaim Non-Jiwa. Bahan dari rangkaian utama.*

| Ukuran | Nilai |
| --- | ---: |
| Tiket | **15** |
| User story pada spec | 66 |
| Acceptance criteria | 124 |
| Pernyataan terverifikasi dari korpus | 180 |
| Keputusan pemilik pekerjaan | 73 |
| Butir terbuka | 42 |
| Penyimpangan sadar | 57 |

## Masalah

Penanganan klaim treaty inward proporsional hari ini berjalan di atas Pega, dan **data klaim intinya tidak pernah ada sebagai baris Oracle**. `[terverifikasi]` Peserta, baris adjustment, estimasi, spreading, dan interest hidup sebagai *page list* di dalam work object Pega dan dipersist lewat `Obj-Save` ke blob; yang menyentuh Oracle hanyalah **proyeksi** — `OS_AKSEPTASI_KLAIM` dan `JSON_KLAIM`, keduanya hibrida satu kolom dokumen `DATA_JSON` plus beberapa kolom kunci.

Akibatnya bagi pengguna dan bagi perusahaan:

## Solusi

Membangun ulang siklus klaim treaty inward proporsional di Go + React di atas Oracle yang sama, dengan **skema relasional yang dirancang baru** (belum pernah ada), **satu jalur perhitungan uang**, **batas transaksi yang eksplisit untuk penomoran**, dan **penegakan wewenang di lapisan layanan**.

Perilaku bisnis **ditiru apa adanya** — termasuk cacat yang sudah diterima sadar — kecuali pada titik-titik yang ditandai ⚠️ **penyimpangan sadar**.

## Kebutuhan pengguna

Spesifikasi modul ini memuat **66 user story**, dimuat utuh di berkas spec-nya.
Yang dijanjikannya, diringkas: pekerjaan yang hari ini dikerjakan di dalam layar Pega
dapat dikerjakan di sistem baru dengan hasil yang sama, jejak keputusan yang dapat
ditelusuri, dan angka yang dapat direkonsiliasi dengan riwayat.

## Keputusan yang mengikat

- 1. ⚠️ PREFACTOR — skema relasional dirancang baru
- 1a. Tanda tangan procedure penulis, dan arti ketujuh kolom penandanya
- 1b. ⚠️ Tabel proyeksi diikuti APA ADANYA — bukan penyimpangan
- 2. ⚠️ Empat schema Oracle, dan mayoritas tabel tanpa prefix
- 3. ⚠️ Batas transaksi penomoran — keputusan eksplisit
- 4. ⚠️ Uang — satu jalur perhitungan
- 5. Mata uang per baris — invariant lama tidak berlaku
- 6. ⚠️ Total share wajib tepat 100 % — penjaga dua sisi

## Yang tidak dibangun

- **Isi konteks Komite Claim Prop.** Modul ini membangun **batasnya**, bukan isinya. Tangga
- **Claim Non Prop dan Claim Fac In.** ⚠️ `[terverifikasi]` Beberapa rule di modul ini mengeksekusi
- **Sistem inti reinsurance non-life, Kasir, dan Arasapas.** Diperlakukan sebagai sistem luar di balik
- **Lini Syariah sebagai lini terpisah.** Varian ber-`S` diperlakukan **identik** dengan induknya;
- **Perbaikan cacat pada butir 14 Implementation Decisions.** Ketiganya **dipertahankan** atas
- **Layar pemeliharaan master** yang hanya terjangkau dari portal admin dan tidak dicapai dari alur
- **Isi fungsi serialisasi JSON Pega.** Tidak ada di korpus; hanya perilakunya yang terbaca dari data.

## Butir terbuka

Modul ini membawa **42 butir terbuka** yang **tidak ditutup** dokumen ini.
- ~~5~~
- # lama
- Kelompok US
- Registrasi dan validasi masuk
- Insured Interest (TSI)
- Estimasi

*Rujukan: `.scratch\claim-prop\spec.md`*

---

# 11. Claim Non Prop

*Bounded context: Klaim Non-Jiwa. Bahan dari rangkaian kedua.*

| Ukuran | Nilai |
| --- | ---: |
| Tiket | **38** |
| User story pada spec | 39 |
| Pernyataan terverifikasi dari korpus | 0 |
| Keputusan pemilik pekerjaan | 0 |
| Butir terbuka | 0 |
| Penyimpangan sadar | 0 |

## Masalah

Petugas klaim hari ini bekerja di atas sistem yang **tidak menyimpan angkanya**. Seluruh nilai uang klaim — estimasi, alokasi per layer, premi pemulihan, nilai adjustment — tersimpan sebagai **teks di dalam satu kolom BLOB** (`PZPVSTREAM`), dan ketika nilai itu disalin keluar ke tabel bisnis, ia tetap teks: view `CLAIMXOL` mengeluarkan `TotalXOLGross`, `CNPReinstatement`, dan `KursIDR` sebagai `varchar2`.

Akibatnya bagi orang yang memakainya:

## Solusi

Model data relasional di mana **setiap nilai punya kolom, tipe, dan presisi sendiri**, dan setiap aturan yang bisa ditegakkan basis data ditegakkan di sana — bukan diserahkan pada ingatan penulis kode.

Bagi orang yang memakainya, yang berubah:

## Kebutuhan pengguna

Spesifikasi modul ini memuat **39 user story**, dimuat utuh di berkas spec-nya.
Yang dijanjikannya, diringkas: pekerjaan yang hari ini dikerjakan di dalam layar Pega
dapat dikerjakan di sistem baru dengan hasil yang sama, jejak keputusan yang dapat
ditelusuri, dan angka yang dapat direkonsiliasi dengan riwayat.

## Keputusan yang mengikat

- 0. Definisi Layer — ditetapkan sekali, dipakai sama di seluruh spec
- Yang ditetapkan untuk spec ini
- Dikuatkan dari tiga arah, bukan satu — EVIDENCED
- "UR" adalah artefak, bukan nilai domain — dan itu menghapus satu kolom
- Bagaimana view kompatibilitas memetakan halus ke kasar
- Akibatnya terhadap ADR-D-CNP-0024 — dinyatakan, tidak diputuskan sendiri
- 1. Entitas yang dimiliki, beserta kunci alaminya
- 2. Relasi dan kardinalitas

*Rujukan: `dastin\_migration-docs\claim-non-prop\2-to-spec\SPEC-MODEL-DATA.md`*

---

# 12. Komite Claim Fac In

*Bounded context: Komite Klaim. Bahan dari rangkaian utama.*

| Ukuran | Nilai |
| --- | ---: |
| Tiket | **13** |
| User story pada spec | 38 |
| Pernyataan terverifikasi dari korpus | 70 |
| Keputusan pemilik pekerjaan | 64 |
| Butir terbuka | 14 |
| Penyimpangan sadar | 15 |

## Masalah

Penyesuaian nilai klaim fakultatif masuk di atas kewenangan seorang penilai harus **disetujui berjenjang** oleh komite klaim. Di sistem lama, proses itu berjalan — **tetapi tiga hal membuatnya tidak dapat dipercaya, tidak dapat diaudit, dan tidak dapat diubah tanpa rilis.**

**Pertama — tidak ada yang benar-benar menjaga siapa boleh memutuskan.** `[terverifikasi]` Sisiran menyeluruh atas 114 berkas menemukan **tepat satu** pemeriksaan pemilik giliran, dan ia menempel pada **tombol Submit di layar**. ⛔ Siapa pun yang dapat memanggil lapisan layanan dapat menyimpan keputusan untuk jenjang mana pun. ⚠️ Di modul saudaranya, `Komite Claim Prop`, **tidak ada pemeriksaan itu sama sekali** — tidak di aktivitas, tidak di layar.

## Solusi

Membangun ulang persetujuan berjenjang komite klaim fakultatif masuk dengan **tiga hal dipisahkan yang di sistem lama menyatu**:

⭐ **Dari pemisahan itu, tiga masalah di atas selesai sekaligus:** wewenang dapat ditegakkan di lapisan layanan karena giliran punya pemilik yang pasti; jejak audit mencatat **akun sebenarnya** karena tak ada lagi penukaran nama; dan tangga berhenti karena **keadaan tiap jenjang**, bukan karena efek samping langkah yang kebetulan hidup.

## Kebutuhan pengguna

Spesifikasi modul ini memuat **38 user story**, dimuat utuh di berkas spec-nya.
Yang dijanjikannya, diringkas: pekerjaan yang hari ini dikerjakan di dalam layar Pega
dapat dikerjakan di sistem baru dengan hasil yang sama, jejak keputusan yang dapat
ditelusuri, dan angka yang dapat direkonsiliasi dengan riwayat.

## Keputusan yang mengikat

- ID-1 — Tiga benda dipisahkan
- ID-2 — Aturan rujuk-atau-salin
- ID-3 — Pemutus tangga adalah keadaan jenjang
- ID-4 — Penegakan wewenang di lapisan layanan
- ID-5 — Larangan menyetujui klaim sendiri, ditiru apa adanya
- ID-6 — Data kutipan disalin utuh
- ID-7 — Penomoran akseptasi
- ID-8 — Penjaga ganda-bayar dibuat eksplisit

## Yang tidak dibangun

- **1**
- **2**
- **3**
- **4**
- **5**
- **6**
- **7**
- **8**

## Butir terbuka

Modul ini membawa **14 butir terbuka** yang **tidak ditutup** dokumen ini.
- ⛔ **1**
- ⛔ **6**
- **11**
- **12**
- **13**
- **15**

*Rujukan: `.scratch\komite-claim-facin\spec.md`*

---

# 13. Komite Claim Life

*Bounded context: Komite Klaim. Bahan dari rangkaian utama.*

| Ukuran | Nilai |
| --- | ---: |
| Tiket | **10** |
| User story pada spec | 42 |
| Pernyataan terverifikasi dari korpus | 54 |
| Keputusan pemilik pekerjaan | 29 |
| Butir terbuka | 7 |
| Penyimpangan sadar | 0 |

## Masalah

Keputusan akseptasi klaim jiwa bernilai besar tidak diambil satu orang — ia menaiki **tangga persetujuan komite**. Hari ini tangga itu berjalan di Pega, dan bentuknya menyulitkan siapa pun yang ingin memahaminya:

1. **Tangganya tidak terlihat di proses.** `[terverifikasi]` Grafnya hanya **1 Assignment + 1 Decision + 4 connector** — satu assignment yang di-loop. Berapa tingkat dan siapa penyetujunya ditentukan **data**, bukan struktur. Membaca flow tidak memberi tahu apa pun tentang tangganya. 2. **Wewenang tidak ditegakkan.** `[terverifikasi]` Pega hanya *menempatkan* kasus di antrean pemilik `KomiteID`. Tidak ada satu pun pemeriksaan bahwa yang menyimpan keputusan adalah orang itu. Keputusan masuk lewat dropdown wajib di layar, dan **tidak ada rule yang menulisnya**. 3. **Efek keluar bisa gagal diam-diam** — termasuk **Kasir**, jalur pembayaran. Klaim dapat disetujui tanpa pernah sampa...

## Solusi

Membangun **Komite Claim Life** sebagai konteks tersendiri di layanan Go + antarmuka React, yang menulis ke Oracle `POOLDATA` yang sama, dengan:

Perilaku bisnis dibawa apa adanya (paritas), kecuali **tiga penyimpangan sadar** yang sudah diputuskan: penegakan wewenang (**ADR-U-0014**), jaminan efek keluar (**ADR-U-0015**), dan penyatuan blok tulis.

## Kebutuhan pengguna

Spesifikasi modul ini memuat **42 user story**, dimuat utuh di berkas spec-nya.
Yang dijanjikannya, diringkas: pekerjaan yang hari ini dikerjakan di dalam layar Pega
dapat dikerjakan di sistem baru dengan hasil yang sama, jejak keputusan yang dapat
ditelusuri, dan angka yang dapat direkonsiliasi dengan riwayat.

## Keputusan yang mengikat

- 1. Batas konteks
- 2. Penempatan modul
- 3. Mesin tangga persetujuan
- 4. Penegakan wewenang
- 5. Satu jalur simpan, bukan dua blok kembar
- 6. Efek keluar — empat, wajib berhasil
- 7. Kode mati yang tidak dimigrasikan
- 8. Data dan transaksi

## Yang tidak dibangun

- **Isi dan alur internal Claim — Life.** Konteks luar (**ADR-U-0001**); spec ini berhenti di tiga
- **Tiga modul Komite lain** (`Komite Claim FacIn`, `Komite Claim Prop`, `Komite Claim Non Prop`).
- **Merapikan alur.** Urutan tangga dan bentuk keputusan dibawa apa adanya. Tiga penyimpangan sadar
- **Identity & Access sebagai konteks.** Dibangun dari nol, di luar spec ini.
- **Scaffolding kode.** `cmd/`, `internal/`, `frontend/`, `Makefile` belum ada — pekerjaan terpisah
- **Kontrak layanan Kasir.** `[terbuka]` OQ-002 — bentuk permintaan dan makna jawabannya tidak ada di

## Butir terbuka

Modul ini membawa **7 butir terbuka** yang **tidak ditutup** dokumen ini.
- **OQ-035**
- **OQ-007 / OQ-021**

*Rujukan: `.scratch\komite-claim-life\spec.md`*

---

# 14. Komite Claim Prop

*Bounded context: Komite Klaim. Bahan dari rangkaian utama.*

| Ukuran | Nilai |
| --- | ---: |
| Tiket | **14** |
| User story pada spec | 39 |
| Acceptance criteria | 86 |
| Pernyataan terverifikasi dari korpus | 110 |
| Keputusan pemilik pekerjaan | 103 |
| Butir terbuka | 37 |
| Penyimpangan sadar | 0 |

## Masalah

Nusantara Re menjalankan persetujuan komite untuk klaim reasuransi **proporsional** di atas Pega. Setiap kali sebuah baris penyesuaian klaim (*adjustment*) perlu disetujui, sistem membuat satu kasus komite tersendiri, lalu mengedarkannya ke para penyetuju **satu per satu, berurutan**. Ketika penyetuju terakhir menyetujui, sistem menerbitkan nomor akseptasi, membuat dokumen PDF, mengirim data pembayaran ke Kasir, memanggil layanan arasapas, menulis beberapa tabel riwayat, dan mengirim email.

Pega akan ditinggalkan. Persoalannya:

## Solusi

Berkas ini menyusun **apa yang benar-benar dikerjakan sistem lama**, dengan setiap pernyataan menyebut buktinya, sehingga orang berikutnya bisa memeriksa ulang tanpa mengulang delapan ronde.

Isinya **sebelas** bab: lingkup dan identitas · daur hidup kasus · tangga penyetuju · layar komite · tiga jalur · nomor akseptasi · angka uang · efek keluar · tabel yang disentuh · kumpulan keputusan work owner · titik yang **sengaja diubah**. Ditambah lampiran aturan baca.

## Kebutuhan pengguna

Spesifikasi modul ini memuat **39 user story**, dimuat utuh di berkas spec-nya.
Yang dijanjikannya, diringkas: pekerjaan yang hari ini dikerjakan di dalam layar Pega
dapat dikerjakan di sistem baru dengan hasil yang sama, jejak keputusan yang dapat
ditelusuri, dan angka yang dapat direkonsiliasi dengan riwayat.

## Keputusan yang mengikat

- Bab 1 — Lingkup dan identitas
- Identitas satu aturan
- ⭐ Aturan awalan — pemilah data komite dan data klaim
- Jejak Save-As
- [terbuka] bab ini
- Bab 2 — Daur hidup kasus
- Tiga aksi lain pada assignment
- Titik lompat darurat

## Butir terbuka

Modul ini membawa **37 butir terbuka** yang **tidak ditutup** dokumen ini.

*Rujukan: `.scratch\komite-claim-prop\spec.md`*

---

# 15. Komite Claim Non Prop

*Bounded context: Komite Klaim. Bahan dari rangkaian kedua.*

| Ukuran | Nilai |
| --- | ---: |
| Tiket | **13** |
| User story pada spec | 43 |
| Pernyataan terverifikasi dari korpus | 0 |
| Keputusan pemilik pekerjaan | 0 |
| Butir terbuka | 0 |
| Penyimpangan sadar | 0 |

## Masalah

Komite klaim non-proporsional hari ini bekerja di dalam Pega, dan cara kerjanya menyimpan enam masalah yang dirasakan langsung oleh orangnya:

**Pemutus tidak tahu apakah gilirannya benar-benar giliran dia.** Pemeriksaan wewenang di sistem lama hanya memasang pesan lalu melanjutkan; siapa pun yang dapat membuka layar dapat mengirim keputusan, dan keputusan itu tersimpan. Pengenal pengguna dicocokkan sebagai potongan teks, sehingga seseorang yang namanya kebetulan terkandung di dalam nama pemegang jenjang ikut lolos.

## Solusi

Modul Komite dibangun ulang sebagai modul di dalam konteks Klaim (`ADR-D-KCNP-0032`), dengan **sirkulasi** sebagai agregat sendiri: satu berkas yang mengedarkan satu **usulan** kepada sejumlah **jenjang** untuk diputuskan.

Dari sudut pandang orang yang memakainya:

## Kebutuhan pengguna

Spesifikasi modul ini memuat **43 user story**, dimuat utuh di berkas spec-nya.
Yang dijanjikannya, diringkas: pekerjaan yang hari ini dikerjakan di dalam layar Pega
dapat dikerjakan di sistem baru dengan hasil yang sama, jejak keputusan yang dapat
ditelusuri, dan angka yang dapat direkonsiliasi dengan riwayat.

## Keputusan yang mengikat

- Batas modul dan transaksi
- Mesin keadaan
- Jenjang aktif — satu penentu
- Tabel seleksi dan roster
- Wewenang
- Pencatatan keputusan
- Maksud dan akibat
- Tulis-balik ke klaim

## Yang tidak dibangun

- PAGAR-01
- PAGAR-02
- PAGAR-03
- PAGAR-04
- PAGAR-05
- PAGAR-06
- PAGAR-07
- PAGAR-08

*Rujukan: `dastin\_migration-docs\komite-claim-non-prop\SPEC-KOMITE-01.md`*

---

# 16. NB Treaty In

*Bounded context: Realisasi Treaty. Bahan dari rangkaian utama.*

| Ukuran | Nilai |
| --- | ---: |
| Tiket | **28** |
| User story pada spec | 41 |
| Acceptance criteria | 160 |
| Pernyataan terverifikasi dari korpus | 144 |
| Keputusan pemilik pekerjaan | 119 |
| Butir terbuka | 34 |
| Penyimpangan sadar | 33 |

## Masalah

Nusantara Re menerima penawaran treaty inward dari ceding company. Ketika sebuah penawaran disetujui, ia harus **direalisasikan** menjadi polis treaty: datanya dilengkapi, diperiksa berjenjang, diberi nomor polis, dan disimpan sebagai kontrak yang mengikat.

Hari ini pekerjaan itu berjalan di atas Pega, dan **empat hal membuatnya mahal dan rapuh**:

## Solusi

Realisasi treaty inward dibangun ulang di Go + React + Oracle, dengan **lima perubahan pokok** terhadap sistem lama — ⭐ seluruhnya `[keputusan work owner]`:

⭐ **Yang TIDAK berubah, dan itu disengaja:** aturan bisnisnya. Nilai penanda persetujuan, penggolongan jenis usaha, rumus pajak brokerage, pengecualian mata uang, dan antrean bersama **ditiru apa adanya** — termasuk ketidakseragamannya.

## Kebutuhan pengguna

Spesifikasi modul ini memuat **41 user story**, dimuat utuh di berkas spec-nya.
Yang dijanjikannya, diringkas: pekerjaan yang hari ini dikerjakan di dalam layar Pega
dapat dikerjakan di sistem baru dengan hasil yang sama, jejak keputusan yang dapat
ditelusuri, dan angka yang dapat direkonsiliasi dengan riwayat.

## Keputusan yang mengikat

- 5.1 ⭐ Sumber data pindah ke view relasional
- 5.2 ⭐ Tangga persetujuan menyusut menjadi tiga
- 5.3 ⭐ Penanda persetujuan dan empat cabang putusan
- 5.4 ⭐ Peran menggantikan nama orang dan nomor telepon
- 5.5 ⭐ Penggolongan jenis usaha
- 5.6 ⭐ Uang, persentase, dan pajak
- 5.7 ⭐ Keutuhan penyimpanan — satu transaksi
- 5.8 ⭐ Tanggal

## Yang tidak dibangun

- ⛔ **28 aturan yatim, 561 langkah `Property-Set`**
- ⛔ **6 aturan pembongkar JSON**
- ⚠️ ~~**Rantai perhitungan uang**~~ — **TIDAK LAGI DIKELUARKAN**
- ⛔ **20 nomor polis di dalam aturan uang**
- ⛔ **Dua aturan uang yang hanya berupa catatan**
- ⛔ **Pencarian antrean menurut nomor urut**
- ⛔ **Rule `When\isApproved`**
- ⛔ **Properti `isApprovedtoDeptHead`**

## Butir terbuka

Modul ini membawa **34 butir terbuka** yang **tidak ditutup** dokumen ini.
- ⭐ **20**
- ⭐ **21**
- ⭐ **22**
- ⛔⛔ **23**

*Rujukan: `.scratch\nb-treaty-in\spec.md` . `.scratch\nb-treaty-in\spec-penyimpanan-relasional.md`*

---

# 17. EDM Treaty In

*Bounded context: Realisasi Treaty. Bahan dari rangkaian utama.*

| Ukuran | Nilai |
| --- | ---: |
| Tiket | **12** |
| User story pada spec | 73 |
| Acceptance criteria | 116 |
| Pernyataan terverifikasi dari korpus | 91 |
| Keputusan pemilik pekerjaan | 94 |
| Butir terbuka | 19 |
| Penyimpangan sadar | 17 |

## Masalah

Hari ini, mengubah kontrak treaty inward yang sudah berjalan dikerjakan di dalam sistem Pega yang akan dihentikan. Endorsemen — perubahan atas polis yang sudah terbit — adalah **jenis kasus tersendiri** di sana, dengan kelas kerjanya sendiri, tangganya sendiri, dan perhitungan selisihnya sendiri.

Yang membuat pemindahannya sulit bukan besarnya, melainkan **empat sifat yang tidak terlihat dari luar**:

## Solusi

Endorsemen treaty inward dibangun ulang di Go + React + Oracle sebagai **jenis kasus tersendiri**, sejajar dengan polis baru, bukan sebagai tahap di dalamnya.

⭐ **Yang membuat ini mungkin sekarang, dan tidak mungkin dua minggu lalu:** isi 380 langkah penetapan nilai **terbaca penuh** dari ekspor — **1.309 pasangan nama=nilai**, **379 dari 380** langkah membawa isinya. Rumusnya tidak perlu ditebak. `[terverifikasi]`

## Kebutuhan pengguna

Spesifikasi modul ini memuat **73 user story**, dimuat utuh di berkas spec-nya.
Yang dijanjikannya, diringkas: pekerjaan yang hari ini dikerjakan di dalam layar Pega
dapat dikerjakan di sistem baru dengan hasil yang sama, jejak keputusan yang dapat
ditelusuri, dan angka yang dapat direkonsiliasi dengan riwayat.

## Keputusan yang mengikat

- 5.1 Bentuk kasus dan alur
- 5.2 Pembekuan data lama — [keputusan work owner] P58
- 5.3 Perhitungan selisih — [keputusan work owner] P57
- 5.4 Penomoran endorsemen — [keputusan work owner] P55
- 5.5 Pembatalan — [keputusan work owner] P56
- 5.6 Penyebaran risiko — [keputusan work owner] P60
- 5.7 Konversi ke produksi — [keputusan work owner] P50, P51, P52, P53
- 5.8 Sumber data

## Yang tidak dibangun

- galat `float64` untuk besaran ini
- galat yang benar-benar tersimpan
- nisbahnya
- pembagian empat dalam desimal
- **a**
- **b**
- **c**
- ⛔ **1**

## Butir terbuka

Modul ini membawa **19 butir terbuka** yang **tidak ditutup** dokumen ini.

*Rujukan: `.scratch\edm-treaty-in\spec.md` . `.scratch\edm-treaty-in\spec-penyimpanan-relasional.md`*

---

# 18. Treaty In

*Bounded context: Kontrak Treaty. Bahan dari rangkaian kedua.*

| Ukuran | Nilai |
| --- | ---: |
| Tiket | **52** |
| Pernyataan terverifikasi dari korpus | 0 |
| Keputusan pemilik pekerjaan | 0 |
| Butir terbuka | 0 |
| Penyimpangan sadar | 0 |

## Masalah

**Tanggal:** 23 September 2026 **Keadaan berkas:** **langkah 1–10 sudah dijalankan.** Bentuk awalnya menuangkan langkah 1–7 ke disk; langkah 8 (invarian, kini **68**), 9 (peta telusur per jalur — **selesai, semestanya kurang**, lihat `4-erd-dan-tabel-datar/AUDIT-PENYEBUT-POHON.md`), dan 10 (ERD) **sudah selesai** dan §11.1 mencatat keadaannya masing-masing.

---

*Rujukan: `dastin\_migration-docs\treaty-in\SPEC-MODEL-DATA.md` . `dastin\_migration-docs\treaty-in\SPEC-INVARIAN.md` . `dastin\_migration-docs\treaty-in\CONTEXT.md`*

---

# 19. Treaty In Adjustment

*Bounded context: Kontrak Treaty. Bahan dari rangkaian kedua.*

| Ukuran | Nilai |
| --- | ---: |
| Tiket | **14** |
| Pernyataan terverifikasi dari korpus | 0 |
| Keputusan pemilik pekerjaan | 0 |
| Butir terbuka | 0 |
| Penyimpangan sadar | 0 |

## Masalah

**Tanggal:** 24 September 2026 · **Fase:** to-spec

---

*Rujukan: `dastin\_migration-docs\treaty-in-adjustment\2-to-spec\STRUKTUR-ADDENDUM.md` . `dastin\_migration-docs\treaty-in-adjustment\STRUKTUR-DATA-TREATY-IN-ADJUSTMENT.md` . `dastin\_migration-docs\treaty-in-adjustment\KEPUTUSAN-GRILLING-ADJUSTMENT.md`*

---

# 20. Treaty Contract Out

*Bounded context: Penempatan Treaty. Bahan dari rangkaian utama.*

| Ukuran | Nilai |
| --- | ---: |
| Tiket | **12** |
| User story pada spec | 40 |
| Acceptance criteria | 73 |
| Pernyataan terverifikasi dari korpus | 20 |
| Keputusan pemilik pekerjaan | 31 |
| Butir terbuka | 10 |
| Penyimpangan sadar | 0 |

## Masalah

Seluruh syarat sebuah kontrak treaty non-life — siapa reinsurernya, berapa sharenya, bisnis apa yang ditanggung, dan dua puluh lima jenis **klausul** yang mengikat (limit, EPI, PLA, Ricomm, profit commission, ex-gratia, cash loss limit, territorial limit, dan seterusnya) — hari ini dikelola di satu layar Pega dan disimpan ke delapan tabel Oracle lewat enam stored procedure.

Tujuh hal membuat modul ini berisiko dipindahkan secara salah:

## Solusi

Membangun ulang Treaty Contract Out sebagai **editor master arrangement kontrak treaty non-life** di atas Go + React + Oracle — **editor master murni**: tidak ada tangga persetujuan, tidak ada Submit/Decline, tidak ada status akseptasi. `[terverifikasi]` Nol rule `Flow`, nol rule `When`, nol `StatusAkseptasi` di seluruh 303 berkas.

---

## Kebutuhan pengguna

Spesifikasi modul ini memuat **40 user story**, dimuat utuh di berkas spec-nya.
Yang dijanjikannya, diringkas: pekerjaan yang hari ini dikerjakan di dalam layar Pega
dapat dikerjakan di sistem baru dengan hasil yang sama, jejak keputusan yang dapat
ditelusuri, dan angka yang dapat direkonsiliasi dengan riwayat.

## Keputusan yang mengikat

- 1. Batas konteks dan kepemilikan tulis
- 2. Entitas dan hierarki [data DBA]
- 3. ⚠️ Satu tabel untuk dua puluh lima jenis klausul (penyimpangan sadar 2)
- 4. [keputusan work owner] Istilah diikuti apa adanya
- 5. Validasi wajib-isi — berbeda per jenis klausul, ditegakkan di Go
- 6. ⚠️ Simpan atomik — procedure tidak commit sendiri (penyimpangan yang menguntungkan)
- 7. Penomoran identitas via sequence [data DBA] — ADR-U-0006
- 8. ⚠️ Kaskade hapus + popup konfirmasi — klausul sengaja dikecualikan (penyimpangan sadar 4)

## Yang tidak dibangun

- Yang dibuang
- **`M_PROPORTIONALARRG` (JSON)** dan **seluruh** kueri `FROM m_PROPORTIONALARRG` — `GetMasterDescriptionEPIParentList`, `GetMasterDescriptionPLAParentList`, `GetMasterDescriptionExGratiaList`, `GetMasterDescriptionFACINParentList`, `GetMasterPortfolioListDetail`, `GetMasterPanggilID`, dan lainnya
- **Fitur salin tahun treaty** — `POOLDATA.PROSESCOPY`, `RDBList/SaveMasterCopyData_SQL.xml`, `Activity/BrowseCopyData.xml`, `Activity/SaveTreatyYearMultiple_Act.xml`, bagian salin `Activity/BrowseDeleteRowTreatyInContract.xml`
- **`SetTreatyArrangementDesc_Act`** (`ASM-FW-GISFW-INT-PROPORTIONALARRG`)
- **Rule jenis reasuransi yang tanpa "Old"** — kecuali sebagai perilaku layar master jenis reasuransi
- **Kolom `PROPORTIONALLIST` dan `OBJECT`**
- **Master jenis reasuransi** (`REINSURANCETYPE`) dan **master jenis klausul** (`TREATYDESC`) —
- **Master grup treaty**, **master mata uang**, **master kategori lampiran** — dibaca saja.

## Butir terbuka

Modul ini membawa **10 butir terbuka** yang **tidak ditutup** dokumen ini.
- **OQ-047**
- **OQ-022**
- **OQ kecil**

*Rujukan: `.scratch\treaty-contract-out\spec.md`*

---

# 21. NB Fac In

*Bounded context: Fakultatif. Bahan dari rangkaian Fakultatif.*

| Ukuran | Nilai |
| --- | ---: |
| Tiket | **16** |
| User story pada spec | 55 |
| Acceptance criteria | 4 |
| Pernyataan terverifikasi dari korpus | 40 |
| Keputusan pemilik pekerjaan | 0 |
| Butir terbuka | 0 |
| Penyimpangan sadar | 0 |

## Masalah

Aplikasi Facultative Inward berjalan di atas Pega yang akan dimatikan. Ketika itu terjadi, sebagian pengetahuan tentang cara sistem ini menghitung dan mengambil keputusan **musnah bersamanya** — bukan "sulit diambil", melainkan tidak ada sumber penggantinya.

Tetapi pekerjaan tidak dapat menunggu sampai semuanya diketahui. `[terverifikasi]` Sebagian besar sistem lama **belum dapat dispesifikasikan sekarang**: isi 35 stored procedure yang dilewati seluruh tulisan produksi tidak ada di korpus, tidak ada satu pun DDL sehingga tipe setiap kolom tidak diketahui, 604 dropdown tidak punya daftar nilai, dan baris tabel keputusan yang menentukan arah setiap keputusan underwriting tidak ikut terekspor.

## Solusi

Lima modul dibangun lebih dulu, dipilih **bukan** karena mudah, melainkan karena **tidak satu pun bagiannya bergantung pada jawaban yang belum masuk**:

Ketiganya diuji lewat **tiga seam** saja, masing-masing di titik tertinggi modulnya, ditambah satu pemeriksaan yang dijalankan **kompilator**, bukan test.

## Kebutuhan pengguna

Spesifikasi modul ini memuat **55 user story**, dimuat utuh di berkas spec-nya.
Yang dijanjikannya, diringkas: pekerjaan yang hari ini dikerjakan di dalam layar Pega
dapat dikerjakan di sistem baru dengan hasil yang sama, jejak keputusan yang dapat
ditelusuri, dan angka yang dapat direkonsiliasi dengan riwayat.

## Keputusan yang mengikat

- Arsitektur umum
- Modul 1 — pkg/money
- Modul 2 — pkg/ratio
- Modul 3 — internal/rules
- Modul 4 — services/acceptance
- Modul 5 — services/premium

*Rujukan: `jefri\OUTPUT FIX\04-spec\03-spec-modul-terverifikasi.md` . `jefri\OUTPUT FIX\04-spec\01-modul-go.md`*

---

# 22. RNW Fac In

*Bounded context: Fakultatif. Bahan dari rangkaian Fakultatif.*

| Ukuran | Nilai |
| --- | ---: |
| Tiket | **8** |
| User story pada spec | 11 |
| Pernyataan terverifikasi dari korpus | 19 |
| Keputusan pemilik pekerjaan | 0 |
| Butir terbuka | 0 |
| Penyimpangan sadar | 0 |

## Masalah

Renewal berjalan di atas **basis kode yang sama dengan New Business**. `[terverifikasi]` **1.907 dari 1.927** berkas RNW **identik byte-per-byte** dengan berkas bernama sama di NB — nol berbeda, metadata ekspor termasuk (K-031).

Yang membedakan renewal hanyalah **cara sebuah kasus masuk dan ditampilkan**.

## Solusi

**Pakai ulang seluruh modul NB tanpa perubahan.** Bangun **14 berkas delta** sebagai lapisan **masuk dan tampilan** di atasnya.

`[terverifikasi]` **Renewal tidak punya mesin hitung ulang.** Hanya empat activity RNW yang menggerbangi `StatusBusiness = 2` — masa berlaku tanggal, validasi tanggal, proteksi spreading, input pembayaran — dan **tidak satu pun perhitungan premi**. Keempatnya berkas bersama yang identik dengan NB.

## Kebutuhan pengguna

Spesifikasi modul ini memuat **11 user story**, dimuat utuh di berkas spec-nya.
Yang dijanjikannya, diringkas: pekerjaan yang hari ini dikerjakan di dalam layar Pega
dapat dikerjakan di sistem baru dengan hasil yang sama, jejak keputusan yang dapat
ditelusuri, dan angka yang dapat direkonsiliasi dengan riwayat.

## Yang tidak dibangun

- **Modul inti** — perhitungan · registry predikat · tangga akseptasi · spreading
- **`EditMarketing`** (Section + FlowAction)
- **Fitur pilih-tertanggung** — `Harness\ChooseInsured` · `Section\ChooseInsuredDtl` · `ReportDefinition\BrowseAccountInsuredEDM` · `Activity\SetDataInsuredEDM_Act`
- **Gerbang masuk teknis** (work type / portal Pega)
- **Rule konversi asli kelas `Work`**
- **Blok `pySaveSQL` `PROSESCOPY`**
- **Isi `InputRenewalDtl` per-field** (4,65 MB bersama `_IsUW`)
- **`Activity\SetOldData` dan properti ber-sufiks `TSIOld`**

*Rujukan: `jefri\OUTPUT FIX\04-spec\04-spec-rnw.md`*

---

# 23. Endorsment Fac In

*Bounded context: Fakultatif. Bahan dari rangkaian Fakultatif.*

| Ukuran | Nilai |
| --- | ---: |
| Tiket | **22** |
| User story pada spec | 35 |
| Pernyataan terverifikasi dari korpus | 84 |
| Keputusan pemilik pekerjaan | 1 |
| Butir terbuka | 0 |
| Penyimpangan sadar | 0 |

## Masalah

Endorsement bukan penerbitan polis baru. Ia **mengubah polis yang sudah berjalan**, dan yang dipertanggungjawabkan ke reasuradur bukan nilai penuh polis melainkan **selisih** antara keadaan sesudah dan keadaan sebelum.

Selisih tidak dapat dihitung tanpa "keadaan sebelum". Di sistem lama, "keadaan sebelum" bukan satu hal — ia **tiga mekanisme berbeda** yang disimpan di tempat berbeda, diisi oleh rule berbeda, pada waktu berbeda, dan dikonsumsi oleh bagian berbeda. Ketiganya sering tertukar ketika dibaca sekilas, karena semuanya "nilai lama".

## Solusi

Bangun modul **before-image** sebagai fondasi siklus endorsement: satu modul yang menyiapkan "keadaan sebelum" dengan **tiga lapis yang tetap terpisah**, persis seperti sistem lama, lalu menyediakan satu kontrak selisih di atasnya.

Ketiga lapis dipertahankan sebagai konsep terpisah di dalam model domain — **bukan** disederhanakan menjadi satu "snapshot":

## Kebutuhan pengguna

Spesifikasi modul ini memuat **35 user story**, dimuat utuh di berkas spec-nya.
Yang dijanjikannya, diringkas: pekerjaan yang hari ini dikerjakan di dalam layar Pega
dapat dikerjakan di sistem baru dengan hasil yang sama, jejak keputusan yang dapat
ditelusuri, dan angka yang dapat direkonsiliasi dengan riwayat.

## Keputusan yang mengikat

- 1. Tiga lapis tetap tiga, bukan satu
- 2. Modul yang dibangun
- 3. Kontrak lapis A — penyalinan yang tidak boleh diringkas
- 4. Kontrak lapis B — 52 penugasan, tujuh lini, pola nilai tunggal
- 5. Dua guard yang sengaja dipertahankan (K-046)
- 6. ⚠️ A.5 — satu-satunya perbaikan sadar di modul ini
- 7. Tipe data
- 8. Kontrak selisih

*Rujukan: `jefri\OUTPUT FIX\04-spec\05-spec-edm-before-image.md` . `jefri\OUTPUT FIX\04-spec\06-spec-edm-selisih.md` . `jefri\OUTPUT FIX\04-spec\09-spec-edm-only.md` . `jefri\OUTPUT FIX\09-facout\01-temuan-dan-rancangan-facout.md`*

---

# 24. Kepastian bahan per modul

[PENTING] **Bab ini yang membuat dokumen ini jujur tanpa perlu ditunda.** Ia menyatakan per
modul apa yang **sudah terverifikasi dari korpus** dan apa yang **masih menunggu pihak lain**,
beserta pihaknya. Pembaca tidak perlu menebak seberapa matang sebuah bab.

| Modul | Terverifikasi | Keputusan pemilik | Butir terbuka | Menunggu |
| --- | ---: | ---: | ---: | --- |
| Claim Life | 75 | 72 | 23 | pemilik pekerjaan |
| PremiumList Life | 46 | 35 | 1 | DBA |
| Endorsement Life | 38 | 39 | 3 | DBA |
| Master Product Name Life | 22 | 42 | 5 | DBA |
| Master Contract Retro Life | 32 | 42 | 6 | DBA |
| Claim Fac In | 123 | 41 | 38 | DBA |
| Claim Prop | 180 | 73 | 42 | DBA |
| Claim Non Prop | 0 | 0 | 0 | - |
| Komite Claim Fac In | 70 | 64 | 14 | DBA |
| Komite Claim Life | 54 | 29 | 7 | DBA |
| Komite Claim Prop | 110 | 103 | 37 | DBA |
| Komite Claim Non Prop | 0 | 0 | 0 | - |
| NB Treaty In | 144 | 119 | 34 | DBA |
| EDM Treaty In | 91 | 94 | 19 | DBA |
| Treaty In | 0 | 0 | 0 | - |
| Treaty In Adjustment | 0 | 0 | 0 | - |
| Treaty Contract Out | 20 | 31 | 10 | DBA |
| NB Fac In | 40 | 0 | 0 | - |
| RNW Fac In | 19 | 0 | 0 | - |
| Endorsment Fac In | 84 | 1 | 0 | - |

## Yang dinyatakan sendiri oleh rangkaian Fakultatif

Spesifikasi rangkaian Fakultatif menyatakan batas bahannya secara terbuka, di bab *Out of
Scope*-nya sendiri. Ini **fakta yang dicatat**, bukan alasan menunda dokumen:

| Yang belum ada | Jumlah | Menunggu |
| --- | ---: | --- |
| isi stored procedure jalur tulis produksi | **35** | DBA |
| definisi tabel fisik | **nol tersedia** | DBA |
| daftar nilai untuk daftar-pilihan layar | **604** tanpa daftar nilai | Product |
| baris tabel keputusan underwriting | tidak terekspor | pemilik ekspor |
| aturan yang perlu diekspor ulang | **77** | pemilik ekspor |
| pertanyaan kuesioner terbuka | **27** | UW . Product . IT . DBA |

[!] **Akibatnya bagi pembaca:** bab modul Fakultatif menyatakan **bentuk** pekerjaan dengan
percaya diri, tetapi **angka dan daftar nilainya** menunggu keenam butir di atas. Tiket
Fakultatif yang tertahan mencerminkan hal yang sama.

---

# 25. Model data ringkas

Penyimpanan berpindah dari **dokumen teks di dalam satu kolom** menjadi **tabel relasional**.
Bab ini menyatakan bentuknya, **tanpa definisi tabel fisik** - presisi dan tipe fisik
dicocokkan dengan DBA di dalam tiket.

## Bentuk yang disepakati untuk polis treaty

| Lapis | Isi | Hubungan |
| --- | --- | --- |
| tabel kerja | akar **lintas-lini**, dipakai bersama lebih dari satu lini | akar |
| tabel generasi | satu baris per **generasi** polis | berbagi kunci utama dengan tabel kerja |
| tabel kuotasi | data penawaran | satu-ke-satu |
| tabel ceding | satu baris per ceding | satu-ke-banyak |
| tabel angsuran | jadwal pembayaran | satu-ke-banyak |
| tabel sebaran | pembagian risiko | satu-ke-banyak |
| tabel lapisan | hanya pada bentuk non-proporsional | satu-ke-banyak |
| tabel riwayat akseptasi | sudah datar di sistem lama | **tidak dibuat ulang** |

**Endorsemen tidak membuat tabel baru.** Ia menempati tabel yang sama dengan nomor generasi
lebih besar dari nol dan penunjuk ke generasi sebelumnya, ditambah **tabel proyeksi selisih**
yang kosong pada polis baru.

## Tiga sifat yang menentukan bentuknya

1. **Dokumen lama tidak berbentuk tetap.** Tiga contoh nyata memberi 37, 64, dan 47 medan
   tingkat atas; hanya **23** ada di ketiganya, sedangkan gabungannya **74**. Bentuknya
   mengikuti jenis kontrak.
2. **Satu nama daftar, dua kedalaman.** Daftar angsuran bersarang pada satu bentuk dan datar
   pada bentuk lain - keduanya nyata di dalam aturan sistem lama. Pengurai yang menganggapnya
   seragam akan patah.
3. **Angka uang di produksi sudah membawa galat.** Selisih sebesar sepersepuluh juta terbukti
   tersimpan permanen, dan ia lahir di rantai perhitungan, bukan di penyimpanan.

[!] **Setiap cacah medan di dokumen ini adalah batas bawah, bukan total.** Sensus berasal dari
rujukan di dalam aturan; medan yang ada di dokumen tetapi tidak pernah dirujuk satu aturan pun
tidak tertangkap cara mana pun.

[AWAS] **Dua rancangan kolom untuk satu tabel fisik** - lihat bab 27 dan dokumen Steering
bagian 5a. Tidak diputuskan di sini.

---

# 26. Keputusan arsitektur yang paling mengikat

Seluruh **105** keputusan dimuat utuh di dokumen *Keputusan Arsitektur*. Yang paling mengikat
pekerjaan sehari-hari:

| Keputusan | Isinya |
| --- | --- |
| uang bukan bilangan mengambang | berlaku seluruh sistem, tanpa pengecualian |
| presisi kolom uang | **satu seri menetapkan tiga puluh delapan digit**, spec penyimpanan treaty menetapkan **dua puluh**; keduanya berskala delapan |
| tabel inti berbagi kunci utama dengan tabel kerja | tidak ada kolom penyambung |
| tabel kerja adalah akar lintas-lini | satu tabel menjadi akar bagi lebih dari satu lini |
| arah ketergantungan | `handlers -> services -> repository`, tidak pernah dibalik |
| nol pemanggilan program tersimpan | logikanya ditulis ulang di lapisan layanan |
| rumus selisih tinggal di lapisan layanan | tampilan basis data yang menghitung selisih **dibatalkan** |

[AWAS] **Keputusan lahir di tiga rangkaian yang masing-masing memulai penomoran dari satu.**
Nomor yang sama dapat berarti keputusan berbeda. Dokumen *Keputusan Arsitektur* memberi awalan
seri dan memuat **konkordansi 105 baris**, serta **pasangan sebidang** - keputusan dari seri
berbeda yang membahas hal yang sama. **Tidak ada yang dimenangkan di sana.**

---

# 27. Butir terbuka per pemilik peran

[AWAS] **Tidak satu pun ditutup dokumen ini.** Penutupan milik pemilik perannya.

Terhitung **239 kemunculan** penanda butir terbuka di seluruh spesifikasi modul,
ditambah kuesioner Fakultatif yang berdiri sendiri.

| Pemilik peran | Yang ditunggu darinya |
| --- | --- |
| **DBA** | isi **35** stored procedure jalur tulis produksi . definisi tabel fisik . presisi kolom uang . panduan bentuk dokumen lama . skema yang dipakai tiap tabel |
| **Product & Underwriting** | daftar nilai untuk **604** daftar-pilihan layar . arti kode lini bisnis . apakah tarif pajak boleh berubah menurut waktu . batas berapa kali satu polis boleh diendorse |
| **Actuarial** | dasar perhitungan yang tidak terbaca dari aturan, tercatat per modul |
| **Finance** | perlakuan atas galat angka yang sudah tersimpan . ketidakseragaman presisi pada rumus yang sama . selisih mana yang dipakai laporan |
| **IAM** | pemetaan nama orang menjadi peran . peran yang belum tersedia untuk sebagian posisi |
| **IT / pemilik ekspor** | **77** aturan perlu diekspor ulang . baris tabel keputusan underwriting . lingkup pemindahan dokumen lama |

## Butir yang memblokir pembangunan

| Butir | Menahan | Pemilik |
| --- | --- | --- |
| isi stored procedure dan definisi tabel | seluruh jalur tulis Fakultatif | DBA |
| daftar nilai daftar-pilihan | layar Fakultatif | Product |
| aturan yang perlu diekspor ulang | perilaku yang belum terbaca | pemilik ekspor |
| lingkup pemindahan dokumen lama | pemuat migrasi di dua rangkaian | pemilik pekerjaan |
| rancangan kolom tabel akar bersama | penyimpanan lintas lini | pemilik pekerjaan |

---

# 28. Risiko yang diterima sadar

Yang di bawah **bukan cacat yang belum diperbaiki**. Ia perilaku yang **sengaja ditiru atau
sengaja tidak ditiru**, dengan alasannya tertulis.

Terhitung **141 kemunculan** penanda penyimpangan sadar di seluruh spesifikasi modul.

| Risiko | Sifatnya | Kenapa diterima |
| --- | --- | --- |
| **galat angka uang di data lama diwariskan** | ditiru | menghitung ulang membuat laporan lama tidak lagi dapat direkonsiliasi |
| **ketidakseragaman presisi pada rumus yang sama** | ditiru | perbedaan nilainya sudah ada di data sekarang, bukan risiko baru |
| **satu transaksi menggantikan commit yang tersebar** | tidak ditiru | sistem lama tidak punya jaminan keutuhan; yang baru lebih ketat |
| **pemeriksaan duplikat berbasis nilai uang tidak ditiru** | tidak ditiru | penjaganya dipatahkan galat presisi yang terbukti ada |
| **pesan galat tanpa markah tampilan** | tidak ditiru | penyajian milik lapisan tampilan |
| **wewenang lewat peran, bukan nama orang** | tidak ditiru | aturan lama berhenti bekerja diam-diam ketika orangnya pindah |
| **rancangan kolom tabel akar dibangun dari satu sisi saja** | diterima | rekonsiliasi antar sisi dinyatakan gugur di rangkaian Fakultatif; irisan kolom dan tabrakan ruang kunci akan muncul saat integrasi |
| **kolom pembeda lini dicabut** | diterima | dinyatakan di luar proyek; satu tabel akan memuat baris dua lini tanpa pembeda |

[AWAS] **Dua baris terakhir belum dipertemukan dengan keputusan arsitektur yang mewajibkan
sebaliknya.** Lihat dokumen Steering bagian 5a dan 5b.

---

*Disusun 25 September 2026 dari spesifikasi, tiket, keputusan arsitektur, dan catatan
penelusuran proyek migrasi Nusantara Re. Nol fakta baru dari korpus; nol butir terbuka ditutup.*