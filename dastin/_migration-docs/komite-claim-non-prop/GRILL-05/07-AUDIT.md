> Modul  : Komite Claim Non Prop · Ronde 05 · 2026-09-20
> Peran  : auditor
> Masukan: `00-LINGKUP.md` … `06-PUTUSAN.md` (ronde 5) · `KETETAPAN.md` · `INVENTARIS-BUKTI.md` (dimutakhirkan ronde 5) · `REGISTER-PAGAR.md` · `REGISTER-DEVIASI.md`
> Status : DITUTUP 2026-09-20
> Sifat  : TAMBAH-SAJA

# 07 · AUDIT PENUTUPAN RONDE 5

Empat cara sebuah perkara tertutup: **BUKTI** (berkas·langkah atau tabel·kolom) ·
**KETETAPAN** (nomor) · **LARUT** (nomor ketetapan yang menghapus jalurnya) · **PAGAR**
(nomor pagar + baris `INVENTARIS-BUKTI.md` + pembukanya). Yang tidak memenuhi satu pun
adalah **TERBUKA**, dan wajib muncul di bagian 4 atau bagian 5.

---

## 1. Buku besar

### 1.1 Temuan `N-01 … N-15`

| # | Cara tutup | Rujukan yang dapat dibuka ulang | Status |
|---|---|---|---|
| N-01 | KETETAPAN | `K5-1` | **TERTUTUP** |
| N-02 | KETETAPAN | keputusan beku no. 5 (K-09); tingkat diturunkan `06-PUTUSAN.md` §2 | **TERTUTUP** |
| N-03 | LARUT | `K5-1` — empat langkah keranjang hilang bersama penentu keduanya | **TERTUTUP** |
| N-04 | BUKTI | `CreateChildKomiteCNP_Act.xml` langkah 9; `CreateChildKomiteCloseNP_Act.xml` nol kemunculan; ditutup rancangan oleh `H-2` | **TERTUTUP** |
| N-05 | KETETAPAN | `K5-6` | **TERTUTUP** |
| N-06 | KETETAPAN | `K5-4` | **TERTUTUP** |
| N-07 | KETETAPAN | `K5-2`, memperluas `H-5` | **TERTUTUP** |
| N-08 | KETETAPAN | `K5-3` | **TERTUTUP** |
| N-09 | KETETAPAN | `K5-5`, melengkapi `H-1`; sisa angka batas bawah → QF-1 | **TERTUTUP sebagian** |
| N-10 | — | naik ke P1 di `06-PUTUSAN.md` §2; menunggu pemilik proses | **TERBUKA** → QF-3 |
| N-11 | LARUT | `H-7` | **TERTUTUP** |
| N-12 | LARUT | aturan kerja "deskripsi langkah bukan aturan", `KETETAPAN.md` bagian 7 | **TERTUTUP** |
| N-13 | BUKTI | pencarian nama atas 338 berkas; penyelidikan dihentikan `06-PUTUSAN.md` §2 | **TERTUTUP** |
| N-14 | BUKTI | `KomiteTreaty_Flow.xml`, `Decision1.pyTicketShapes`; sisa penaik → B5-3 | **TERTUTUP sebagian** |
| N-15 | BUKTI | `ReinstatementPremiumDetails.xml` baris 609/2540/4471/6687; sisa arti → B5-2 | **TERTUTUP sebagian** |

### 1.2 Temuan turunan `T5-01 … T5-04`

| # | Cara tutup | Rujukan | Status |
|---|---|---|---|
| T5-01 | KETETAPAN | `K5-1` — penutupan sirkulasi dari keadaan keputusan, bukan dari perbandingan bilangan | **TERTUTUP** |
| T5-02 | KETETAPAN | `K5-1` — keputusan per jenjang sebagai catatan tetap | **TERTUTUP** |
| T5-03 | — | memisahkan "diajukan untuk ditutup" dari "ditutup" belum diperintahkan ketetapan mana pun | **TERBUKA** → bagian 5 |
| T5-04 | LARUT | `D-2`, diperluas ke kedua pembuat anak; K-11 tidak dibuka kembali (`PUTUSAN-01.md` §8) | **TERTUTUP** |

### 1.3 Pertanyaan ronde ini

| # | Cara tutup | Rujukan | Status |
|---|---|---|---|
| Q5-1 | — | klaim bersyarat dan ambang; pemilik proses | **TERBUKA** → QF-3 |
| *(dibatalkan)* urutan giliran jenjang | BUKTI | `KomiteRouter.xml` langkah 6.1 transisi `6/6` | **TERTUTUP sebelum ditanyakan** |

### 1.4 Perkara lama yang berubah status ronde ini

| # | Sebelum | Sesudah | Oleh |
|---|---|---|---|
| K-02 | DIKUATKAN dengan syarat, perlakuan bertanda tafsir | **TERTUTUP · KETETAPAN `K5-1`** | N-01, sidang |
| Q-2 | TERBUKA, menunggu kueri DBA | **TERTUTUP** untuk yang dipakai rule | N-04 |
| Q-3 | TERBUKA | **TERTUTUP sebagian** | N-15 |
| Q-6 | TERBUKA | **TERTUTUP** | N-13 |
| Q-7 | TERBUKA | **TERTUTUP sebagian** | N-14 |
| F-21 | fakta berlaku | **SALAH KAPRAH**, diganti N-07 | sidang |
| F-16 | fakta berlaku | **DIPERLUAS** — dua dimensi, bukan satu | sidang |
| Pagar A-1a | berlaku | **GUGUR** sebagai PG-07 | `05-GERBANG.md` |

---

## 2. Perkara yang tertutup dua kali atau bertabrakan

1. **`H-5` dan `K5-2` memerintahkan hal yang bertingkat, bukan hal yang sama.** `H-5`
   menolak pada pembuatan; `K5-2` menolak seluruh urutan yang bergantung padanya. Keduanya
   berlaku, dan `KETETAPAN.md` bagian 11 mencatat hubungannya sebagai **DIPERLUAS** —
   bukan dua ketetapan yang bersaing.

2. **`H-1`, `H-6`, `K5-5` menyentuh satu perkara dari tiga sisi.** `H-1` menetapkan seleksi
   sebagai data; `H-6` menolak sirkulasi tanpa jenjang; `K5-5` melarang kombinasi tanpa
   aturan. Tanpa `K5-5`, `H-6` akan menolak sirkulasi yang lahir dari lubang konfigurasi
   seolah itu keadaan yang sah. Ketiganya konsisten, tetapi urutan penerapannya penting dan
   dicatat di sini.

3. **Penutupan yang masih bersandar pada `TAFSIR (menunggu C-01)`.** Ronde ini tidak
   menambah satu pun penutupan yang bersandar pada `PROC_GENERATE_SEQUENCE_NUMBER`. Lima
   penutupan lama tetap bersandar padanya, dan kini punya baris inventarisnya sendiri —
   §2.5 baris 4 — sehingga status `TAFSIR` itu dapat ditunjuk, bukan sekadar diingat.

4. **Sembilan pagar A-4/A-5 yang sebelumnya tidak sah kini sah.** `AUDIT-TUTUP-01`
   mencatat bahwa pagar-pagar itu menunjuk ketiadaan data sementara inventaris hanya
   mengenal berkas. §2.5 menutup cacat bentuk itu. Tidak ada pagar yang berubah isi; yang
   berubah hanya keabsahan rujukannya.

5. **Satu tabrakan deret selesai, satu bertahan.** `G-1 … G-4` ronde 3 kini resmi
   `J-1 … J-4` (`KETETAPAN.md`, Peta penomoran), sehingga tabrakan dengan `G-01 … G-20`
   berakhir. Yang bertahan: rujukan silang di dalam badan ketetapan lama sengaja tidak
   disunting agar bunyinya tetap dapat diperiksa — pembaca harus memakai peta.

---

## 3. Putusan

**SELESAI BERSYARAT.**

Sisanya hanya pengambilan bukti dan pagar yang sah. Tidak ada aliran yang menunggu berkas
untuk dapat ditulis spesifikasinya: A-1a dan A-1b keduanya **BOLEH MULAI**, A-2 tidak
tersentuh dan tetap **BOLEH MULAI**, A-3 **TERBATAS** pada satu butir yang dinamai, dan
tiga aliran beku seluruhnya berpagar sah dengan baris inventaris yang dapat ditunjuk.

**Cacah**

| | Temuan `N` | Turunan `T5` | Pertanyaan | Jumlah |
|---|---:|---:|---:|---:|
| TERTUTUP · BUKTI | 2 | 0 | 1 | **3** |
| TERTUTUP · KETETAPAN | 7 | 2 | 0 | **9** |
| TERTUTUP · LARUT | 3 | 1 | 0 | **4** |
| TERTUTUP · PAGAR | 0 | 0 | 0 | **0** |
| TERTUTUP sebagian | 2 | 0 | 0 | **2** |
| **TERBUKA** | 1 | 1 | 1 | **3** |
| Jumlah | 15 | 4 | 2 | **21** |

Nol pagar baru dipasang ronde ini. Dua pagar lama gugur. Enam ketetapan baru lahir. Delapan
deviasi baru masuk register.

---

## 4. Pertanyaan yang hanya manusia dapat jawab

Tiga, seluruhnya milik pemilik proses, seluruhnya dapat dibawa dalam satu pertemuan.

```
QF-1 · Ambang berbasis nilai dan batas bawah yang hilang
  Perkara   : H-1 (pemulihan naik ke pemilik proses) · N-09
  Pertanyaan: Apakah ambang kewenangan komite dipulihkan menjadi fungsi nilai adjustment,
              dan bila dua kelas tetap dipertahankan, berapa batas bawah rentang yang di
              kode tertulis 30 sementara deskripsi langkah 14 menyebut 15?
  Hanya
  manusia   : Keduanya menetapkan siapa menandatangani berapa. Angka 15 tidak ada di kode
              mana pun — hanya di deskripsi langkah — sehingga tidak ada bukti yang dapat
              memutuskannya, dan menebaknya berarti mengarang kewenangan.
  Pilihan   : (a) Pulihkan ambang berbasis nilai; `.LIMIT_BOTTOM` dan `.LIMIT_TOP` pada
                  roster menjadi hidup, dan pertanyaan batas bawah gugur sendiri.
              (b) Pertahankan dua kelas tetap dan tetapkan batas bawah rentang secara tegas.
              (c) Pertahankan dua kelas tetap dan tutup rentangnya; cabang langkah 14
                  dihapus, bukan dihidupkan.
  Rekomendasi: (b) dengan batas bawah 15 — deskripsi langkah adalah satu-satunya jejak niat
              yang tersisa, dan menghidupkan cabang dengan angka itu memulihkan maksudnya
              tanpa mengubah kelasnya.
  Bila tak
  dijawab   : (c) — menutup rentang menghasilkan perilaku tegas tanpa mengarang angka.
  Memblokir : aturan seleksi jenjang pada pembentukan sirkulasi (A-1b). Bagian lain tidak
              terpengaruh.
```

```
QF-2 · Pembatasan hak baca sirkulasi
  Perkara   : J-4 (didaftarkan sebagai pertanyaan terbuka ke pemilik proses)
  Pertanyaan: Apakah hak baca sirkulasi tetap sama luas dengan hak baca klaim, atau
              komentar komite dibatasi hanya untuk pemegang jenjang?
  Hanya
  manusia   : Komentar komite memuat pertimbangan orang atas keputusan uang; seberapa luas
              ia boleh dibaca adalah kebijakan organisasi, bukan akibat teknis.
  Pilihan   : (a) Tetap seluas hak baca klaim.
              (b) Dibatasi pada pemegang jenjang; pembaca lain melihat keputusan tanpa
                  komentar.
  Rekomendasi: (a) — komentar komite sudah menjadi bagian berkas klaim di sistem lama, dan
              membatasinya menghapus konteks dari orang yang hari ini sudah membacanya.
  Bila tak
  dijawab   : (a), sesuai J-4.
  Memblokir : tidak memblokir — J-4 sudah menetapkan default.
```

```
QF-3 · Klaim bersyarat dan ambang
  Perkara   : N-10 (naik ke P1 di 06-PUTUSAN.md §2) · Q5-1
  Pertanyaan: Klaim bersyarat selalu memakai roster terluas, menimpa kelas 25 juta.
              Kebijakan yang dibawa, atau akibat urutan langkah yang tidak disengaja?
  Hanya
  manusia   : Ia menentukan siapa berwenang menyetujui klaim bersyarat di atas 25 juta.
              Berkas hanya dapat menunjukkan bahwa penimpaan itu terjadi, tidak dapat
              menunjukkan apakah ia dimaksudkan.
  Pilihan   : (a) Kebijakan — "bersyarat" menjadi dimensi kedua pada tabel seleksi H-1.
              (b) Bukan kebijakan — bersyarat mengikuti kelas nilainya seperti klaim lain.
  Rekomendasi: (a). Penimpaan itu berdiri di sub-langkah tersendiri dengan pra-syarat
              eksplisit `Local.Subjectivity==true`, bukan di sela-sela langkah lain —
              bentuknya bentuk aturan, bukan bentuk kecelakaan.
  Bila tak
  dijawab   : (a) — mempertahankan perilaku yang berjalan lebih aman daripada mempersempit
              kewenangan tanpa diminta.
  Memblokir : satu baris pada tabel seleksi H-1. Tidak memblokir aliran mana pun.
```

---

## 5. Sisa pengambilan bukti

| # | Bukti | Pemegang | Menutup | Memblokir? |
|---|---|---|---|---|
| B5-2 | ekspor `Rule-Obj-FieldValue` untuk `FlagProrate` dan `.Type` | admin Pega | sisa N-15, sisa N-04 | **tidak** |
| B5-3 | pencarian `komiteAccept_ticket` pada ruleset di luar dua folder | admin Pega | sisa N-14 | **tidak** |
| B5-4 | ekspor `PostEmailKomiteCNP` | admin Pega | G-13 bagian P1; membuka PG-06 | **tidak** — isi A-6 memang beku |
| B5-5 | ekspor ulang dua activity inti dengan parameter metode | admin Pega | sisa G-01, Q-9 | **tidak** |
| B5-6 | ekspor rule `RDB-List` pengisi `TglProd` | admin Pega | sisa Q-4 | **tidak** |
| B5-7 | isi badan `PROC_GENERATE_SEQUENCE_NUMBER` pada basis data berjalan | DBA | C-01; membuka PG-03; mengangkat lima penutupan dari `TAFSIR` | **ya, bersyarat** — menetapkan periode buku sebagai aturan tetap |
| B5-8 | satu `SELECT` distribusi nilai `.Type` | DBA | pelengkap N-04 | **tidak** |
| B5-9 | pemisahan "diajukan untuk ditutup" dari "ditutup" sebagai ketetapan | kamu sendiri | T5-03 | **tidak** — pekerjaan rancangan, tidak menunggu bukti |

Empat pembuka pagar A-4 dan A-5 tidak diulang di sini; statusnya tidak berubah dan
seluruhnya berada di balik pagar yang kini sah (`INVENTARIS-BUKTI.md` §2.5).

**Nol usulan Tracer.** Dua yang sempat berbentuk Tracer diganti pada `03-LUBANG.md` §3,
dan penggantinya disebut di sana.
