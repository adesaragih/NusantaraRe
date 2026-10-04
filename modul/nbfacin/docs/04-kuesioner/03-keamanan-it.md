# Kuesioner Keamanan / IT — Alamat Email Tertanam di Kode

**Purpose:** menentukan pemilik dan tujuan pemindahan daftar penerima email yang saat ini
**tertanam sebagai literal di dalam kode** sistem lama. Aturan proyek melarang alamat dan endpoint
menjadi literal; kami perlu persetujuan dan pemilik sebelum memindahkannya ke konfigurasi.

**From:** Tim Migrasi Facultative Inward · **To:** Keamanan / IT ·
**How your answers will be used:** dicatat sebagai keputusan proyek, lalu menentukan ke mana daftar
penerima dipindahkan dan siapa yang berwenang mengubahnya di sistem baru.

---

> ## ✅ TERJAWAB LENGKAP — 16 September 2026
>
> Butir 1 dijawab **"tidak usah dipindahkan"**, dan sudah menjadi keputusan resmi **K-020**:
> diterima sebagai **pengecualian eksplisit tercatat** terhadap `CLAUDE.md` §4.4.
>
> Jawaban ini sempat menjadi **Konflik K-3** karena bertabrakan dengan aturan yang mengikat proyek.
> Ditutup dengan mencatat penyimpangannya secara terbuka — **pengecualian yang tercatat dapat
> ditinjau ulang; pelanggaran yang tidak tercatat tidak.**
>
> Naskah pertanyaan **tidak diubah** — hanya stub jawaban yang diisi.
> **Nilai alamat emailnya tetap tidak disalin ke mana pun**, termasuk ke jawaban di bawah.

## Context

Aplikasi Facultative Inward sedang dipindahkan dari Pega ke sistem baru. Salah satu aturan yang
mengikat proyek ini: **alamat, endpoint, dan kredensial tidak pernah menjadi literal di dalam kode**
— semuanya harus datang dari konfigurasi.

Saat memeriksa jalur pengiriman email polis, kami menemukan bahwa aturan pengirim email memuat
**alamat tujuan dan alamat BCC yang ditulis langsung di dalam kode**, termasuk sedikitnya satu
alamat perorangan (bukan alamat fungsional/jabatan).

**Nilai alamatnya tidak kami salin ke dokumen mana pun**, termasuk dokumen ini — hanya
keberadaannya yang kami catat.

Sistem lama sebetulnya sudah punya tempat yang benar untuk hal semacam ini: sebuah tabel referensi
yang menyimpan alamat layanan dan dibaca saat runtime. Kami menduga daftar penerima email sebaiknya
pindah ke sana atau ke mekanisme konfigurasi setara.

## How to answer

Tenggat: **23 September 2026**. Perkiraan waktu pengisian 10–15 menit.

Isi langsung di bawah pertanyaan. Bila sebagian jawabannya bukan wewenang Anda, sebutkan fungsi mana
yang perlu kami tanyai.

---

## Kepemilikan dan pemindahan

### 1. Siapa pemilik daftar penerima email polis, dan ke mana daftar itu sebaiknya dipindahkan?

Tiga hal yang perlu kami ketahui, dalam satu jawaban:

**(a) Kepemilikan.** Siapa — fungsi atau peran, bukan nama orang — yang berwenang menentukan dan
mengubah daftar penerima email polis? Apakah pemilik itu tahu bahwa saat ini alamatnya tertanam di
kode dan hanya dapat diubah lewat deployment?

**(b) Tujuan pemindahan.** Apakah daftar ini dipindahkan ke tabel referensi layanan yang sudah ada
(tempat alamat layanan lain disimpan), ke variabel lingkungan aplikasi, atau ke mekanisme lain yang
Anda tentukan?

**(c) Alamat perorangan.** Sedikitnya satu penerima adalah alamat perorangan. Apakah itu masih
berlaku, dan apakah sebaiknya diganti alamat fungsional/jabatan agar tidak mengikat ke satu individu?

_Why this matters: butir (c) bukan sekadar kerapian. Alamat perorangan yang tertanam di kode berarti
notifikasi berhenti atau salah alamat ketika orang itu berpindah peran — dan tidak ada yang tahu
sampai ada yang mengeluh. Memindahkannya ke konfigurasi juga berarti perubahan penerima
meninggalkan jejak, yang saat ini tidak ada._

> **Tidak usah dipindahkan.**

✅ **TERTUTUP — dicatat sebagai K-020** (16 September 2026), setelah sempat menjadi **Konflik K-3**.

**Keputusan: jawaban diterima.** Alamat tujuan dan BCC pada jalur notifikasi polis **tetap seperti
keadaan sistem lama** dan **tidak** dipindahkan ke tabel referensi maupun variabel lingkungan pada
tahap migrasi ini.

### Dicatat sebagai pengecualian, bukan kelalaian

`CLAUDE.md` §4.4 mewajibkan alamat dan endpoint menjadi konfigurasi, tidak pernah literal. Keputusan
ini **menyimpang dari aturan itu secara sadar**, dan dicatat justru supaya penyimpangannya terlihat.

⚠️ `CLAUDE.md` **tidak diubah** — §4.4 tetap berbunyi seperti semula. Pengecualian ini hidup di
**register keputusan** (`00-KEPUTUSAN-WORK-OWNER.md` K-020), bukan di piagam proyek.

### Dua risiko yang dinyatakan **diterima**

1. Sedikitnya satu penerima adalah **alamat perorangan** — notifikasi akan berhenti atau salah alamat
   ketika orang itu berpindah peran, dan tidak ada yang tahu sampai ada yang mengeluh.
2. Perubahan daftar penerima **menuntut deployment** dan tidak meninggalkan jejak audit.

Keduanya **bukan keberatan atas keputusannya** — keputusannya jelas dan dijalankan. Keduanya dicatat
supaya bila gejalanya muncul nanti, penyebabnya langsung dikenali alih-alih didiagnosis dari nol.

### 📌 Butir (a) dan (c) belum terisi — dan itu wajar

Jawaban menutup butir **(b)** (tujuan pemindahan: tidak dipindahkan), sehingga **(a) kepemilikan** dan
**(c) status alamat perorangan** menjadi tidak relevan untuk keputusan ini. Keduanya **belum
diketahui** dan sengaja tidak ditebak.

Keduanya baru diperlukan bila K-020 ditinjau ulang — dan K-020 sendiri menyebut pemicunya: daftar
penerima berubah, alamat perorangan diganti alamat fungsional, atau audit keamanan memintanya.

---

## Anything else?

Adakah tempat lain di sistem lama yang Anda ketahui menyimpan alamat, endpoint, atau kredensial
sebagai literal di dalam kode? Kami akan memeriksanya saat porting, tetapi petunjuk dari Anda
mempercepat dan memastikan tidak ada yang terlewat.

Catatan terkait yang sudah kami temukan dan tangani terpisah: sebuah tabel referensi menyimpan
**prompt yang menggerakkan model AI** di alur underwriting — artinya perilaku model dapat diubah
tanpa deployment dan tanpa jejak version control. Itu kami angkat sebagai isu tersendiri; sebutkan
di sini bila menurut Anda keduanya sebaiknya ditangani bersama.

>
