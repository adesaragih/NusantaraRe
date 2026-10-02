# Permintaan kepada DBA — modul NB Treaty In

**Kepada:** DBA
**Perihal:** isi empat program basis data, contoh isi kolom JSON, dan tiga hitungan pendek
**Tanggal:** 22 September 2026

---

## Latar singkat

Sistem baru **tidak menyimpan data dalam bentuk dokumen JSON lagi**. Data NB Treaty In akan
disimpan sebagai **tabel datar** yang mengikuti struktur datanya sendiri.

Untuk merancang tabel itu, kami perlu tahu dua hal yang keduanya hanya ada di sisi Anda:
**apa yang selama ini ditulis** oleh program penyimpan, dan **apa isi dokumen JSON** yang selama ini
menampungnya.

Kedua permintaan di bawah karena itu **satu paket**. Salah satunya saja tidak cukup.

---

## 1 · Isi empat program di dalam basis data

Ada empat program yang tersimpan **di dalam basis data**, bukan di aplikasi, dan seluruh penyimpanan
data treaty inward melewatinya:

| Program | Perannya |
| --- | --- |
| `POOLDATA.PEGA_TREATY_IN` | menyimpan data kontrak treaty |
| `POOLDATA.PEGA_JSON_POLIS_TREATYIN` | menyimpan data polis sebagai dokumen JSON |
| `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` | penomoran |
| `POOLDATA.PEGA_DELETE_ERROR_KONVERSI` | **menghapus data produksi bila konversi gagal** |

Yang pertama menerima **24 keterangan sekaligus**.

Kami dapat melihat keterangan apa yang **dikirim masuk** — pemanggilnya ada di dalam berkas yang
kami terima. Yang tidak dapat kami lihat adalah **apa yang dilakukan program itu**: tabel mana yang
diisi, kolom mana, dan pemeriksaan apa yang dijalankan.

**Yang kami minta:** naskah keempatnya.

```sql
SELECT text FROM all_source WHERE name = 'PEGA_TREATY_IN' ORDER BY line;
```

Untuk `PROC_GENERATE_SEQUENCE_NUMBER`, sebagian keterangannya sudah kami terima pada
**15 September 2026**; yang belum adalah naskah lengkapnya.

### Yang keempat paling mendesak

`PEGA_DELETE_ERROR_KONVERSI` **menghapus data produksi** ketika pengiriman ke sistem produksi
gagal. Pemanggilnya:

```
BEGIN POOLDATA.PEGA_DELETE_ERROR_KONVERSI( <kunci kasus>, <lini bisnis>, <keluaran> ); END;
```

Work owner sudah menegaskan bahwa penghapusan itu **memang dikehendaki**, dan bahwa penghapusan
**hanya boleh lewat jalur ini**. Justru karena itu kami perlu naskahnya: tanpa mengetahui **baris
mana saja yang dihapus**, aturan itu tidak dapat ditiru dengan tepat di sistem baru — hanya ditebak.

Dipakai empat modul: `EDM Treaty In`, `Endorsment Fac In`, `NB FacIn`, `RNW Fac In`.

⚠️ Satu keterangan tambahan yang kami perlukan untuk program ini saja: **apakah ia meng-commit
sendiri di dalam?** Blok pemanggilnya tidak memuat `COMMIT`, berbeda dari dua program penyimpan
lainnya. Itu menentukan apakah penghapusan dapat dibatalkan bersama langkah lain yang gagal.

### Kalau naskah tidak boleh dilepas

Itu jawaban yang sah, dan kami sudah menyiapkan jalan kedua supaya tidak perlu bolak-balik.

Yang kami perlukan minimal:

- **pemetaan parameter ke kolom** — parameter ke-berapa masuk ke tabel dan kolom mana
- **daftar pemeriksaan** yang dijalankan di dalamnya, tanpa perlu naskahnya

Itu tidak membuka logika bisnisnya, tetapi cukup bagi kami untuk menulis ulang penyimpanan dengan
hasil yang sama.

---

## 2 · Contoh isi kolom JSON polis

**Ini bahan utama perancangan tabel, bukan pelengkap.**

**Yang kami minta:** **tiga sampai lima isi kolom JSON polis apa adanya**, beserta **ukuran
terbesarnya** di basis data sekarang.

Kolom yang dimaksud adalah yang diisi oleh `PEGA_JSON_POLIS_TREATYIN`. Bila nama tabel dan kolomnya
berbeda dari dugaan kami, mohon sekalian disebutkan — itu justru keterangan yang kami cari.

### Kenapa contoh datanya yang dibutuhkan, bukan penjelasan

Program yang membentuk dokumen JSON itu isinya hanya empat baris. Ia **menyalin seluruh halaman
kerja apa adanya**, tanpa memilih medan apa pun:

```java
ClipboardPage stepPage = tools.getStepPage();
return stepPage.getJSON(false);
```

Artinya tidak ada satu pun aturan di dalam sistem lama yang menyebutkan **medan apa saja** ada di
dalam dokumen itu. Bentuknya bukan rancangan — ia potret halaman kerja pada saat penyimpanan.

Hanya **datanya** yang dapat menunjukkan medan apa saja yang harus ditampung tabel baru, tipenya
apa, dan mana yang bersarang.

Dua modul lain di proyek ini sudah menempuh jalan yang sama, dan tabel barunya dirancang persis dari
contoh data semacam ini.

### Kenapa ukuran terbesarnya juga diminta

Bila dokumennya besar dan bertipe `CLOB`, penyalinan naif ke kolom teks biasa akan **memotong
isinya tanpa pesan galat**. Itu jenis kegagalan yang tidak ketahuan sampai laporan tidak cocok.

---

## 3 · Tiga hitungan pendek

Ketiganya tidak menahan pekerjaan, tetapi jauh lebih murah dijawab sekarang daripada nanti.

### a. Kolom `OPERATORID` pada tabel riwayat

```sql
SELECT COUNT(*) AS total, COUNT(operatorid) AS terisi
  FROM POOLDATA.HISTORYAKSEPTASIPEGA;
```

Kami sudah memastikan bahwa **tidak satu pun perintah di dalam sistem lama menulis ke kolom itu**.
Yang belum kami ketahui adalah apakah ada pengisi **di luar** aplikasi.

Kalau hasilnya nol, tidak ada tindak lanjut. Kalau bukan nol, berarti ada pemicu basis data,
penjadwal, atau sistem lain yang menulis ke tabel ini — dan sesuatu itu belum masuk peta migrasi
kami sama sekali.

### b. Tipe kolom pada view detail treaty

```sql
DESCRIBE POOLDATA.TREATYINDETAILJOINEDM;
```

Kami perlu memastikan tipe dan presisi kolom uang — `LIMITVALUE`, `RETENTIONVALUE`, `EPIVALUE`,
`NETPREMIVALUE`, `MDPVALUE`, `SHAREVALUE` — serta `COMMENCEMENT` dan `TERMINATION`.

Alasannya: di sistem lama nilai-nilai itu diperlakukan sebagai **teks**. Sistem baru harus
membacanya sebagai bilangan dengan presisi penuh, dan kami tidak boleh menebak presisinya.

### c. Berapa kali konversi ke produksi gagal

```sql
SELECT COUNT(*) FROM POOLDATA.JSON_POLIS_MONITORING WHERE sts_konversi = 9;
```

Tiap baris bernilai 9 berarti satu kali penghapusan data produksi pernah dijalankan. Angkanya
menentukan seberapa mendesak pengujian jalur ini di sistem baru — bukan apakah jalurnya benar.

---

## Ringkas

| # | Yang diminta | Bentuk | Menahan apa |
| --- | --- | --- | --- |
| 1 | Naskah empat program basis data | berkas — atau pemetaan parameter ke kolom | penulisan penyimpanan |
| 2 | 3-5 isi kolom JSON polis + ukuran terbesarnya | contoh data | **perancangan tabel datar** |
| 3a | Hitungan kolom `OPERATORID` | dua angka | tidak menahan |
| 3b | Tipe kolom view detail treaty | keluaran `DESCRIBE` | tidak menahan |
| 3c | Berapa kali konversi gagal setahun terakhir | satu angka | tidak menahan |

**Nomor 1 dan 2 adalah satu paket.** Tanpa keduanya, tabel penyimpanan sistem baru hanya dapat
ditebak — dan tebakan pada bentuk penyimpanan berarti data tersimpan di tempat yang berbeda dari
sekarang, tanpa seorang pun tahu sampai laporan tidak cocok.

Kami siap menjelaskan langsung bila ada bagian yang perlu diperjelas.
