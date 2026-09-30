---
status: accepted
tanggal: 2026-09-23
sumber: ronde keputusan dua belas butir — `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`, keputusan work owner
---

# Daftar medan disusun dari aturan, bukan dari panduan bentuk dokumen

Daftar kolom untuk memecah dokumen JSON menjadi baris disusun dengan **menyapu aturan Pega** —
Activity, Section, DataTransform, Flow, RDBList — **bukan** dari keluaran `JSON_DATAGUIDE`.

Berlaku untuk setiap modul yang memecah dokumen JSON tersimpan, bukan hanya modul tempat
keputusan ini lahir.

## Kenapa panduan bentuk dokumen tidak dipakai sebagai dasar

`[terverifikasi]` Panduan yang diterima memuat 378 jalur. **Satu dokumen produksi tunggal memuat
95 jalur yang tidak ada di dalamnya** — dan bukan medan pinggiran: daftar perusahaan ceding,
seluruh daftar pecahan penyebaran, nilai klaim, nilai kelebihan kerugian, dua kolom hasil, nilai
sisa, pengenal induk dokumen.

⭐ Artinya panduan itu **sah** membuktikan sebuah medan **ada**, tetapi **tidak sah** membuktikan
sebuah medan **tidak ada**. Dasar yang hanya bisa membuktikan satu arah tidak dapat dipakai untuk
mengunci daftar kolom.

## Kedudukan panduan bentuk dokumen: penambal

Ia tetap berguna untuk dua hal:

1. **Panjang maksimum per medan** — menghapus kebutuhan menebak ukuran kolom.
2. **Medan yang ditulis sistem**, yang karena itu tidak muncul di aturan mana pun. Pada sapuan
   pertama ada **24 nama** semacam itu.

⇒ Yang dipakai adalah **gabungan**: sapuan aturan sebagai dasar, panduan sebagai penambal.

## Dua titik buta yang wajib ditutup saat menyapu

`[terverifikasi]` Keduanya sudah terbukti meleset **empat kali** dalam proyek ini.

**1 · Dua konvensi tag yang berlawanan.**
`Activity` menyimpan penugasan langkah pada `PropertiesName` dan `PropertiesValue` — **tanpa**
awalan `py`. `DataTransform` dan `Flow` memakai `pyPropertiesName` dan `pyPropertiesValue` —
**dengan** awalan. Sapuan yang hanya mengenal satu konvensi kehilangan separuh isi.

**2 · Rujukan relatif di dalam loop.**
Sapuan yang berjangkar pada nama halaman *(`PolicyTreatyIn.…`)* tidak melihat medan yang ditulis
relatif di dalam perulangan *(`.PremiumSpreaded`, `.CedingCoName`)*. Sapuan wajib menangkap
keduanya.

⚠️ Dua jebakan teknis lain: buang `pyExpressionGadget` sebelum mencocokkan, dan **jangan**
meng-unescape entitas HTML sebelum mencocokkan pola struktur — `&lt;` memecahkan pola tag.

## Akibat

1. Setiap modul yang memecah dokumen JSON menyapu aturannya sendiri lebih dulu, dan menyimpan
   hasilnya sebagai daftar medan yang dapat diperiksa.
2. ⭐ **Cacah hasil sapuan adalah batas bawah, bukan total.** Pemecah dokumen **wajib** menyediakan
   penampung medan tak dikenal, dan penampung itu **wajib kosong** sebelum pekerjaan dinyatakan
   selesai. Itu jaring pengaman yang menggantikan kepastian yang memang tidak ada.
3. Test yang menemukan medan hilang **diam-diam** — tanpa masuk penampung — **gagal**.

## Yang catatan ini TIDAK putuskan

⛔ Bukan tentang tipe kolom. Itu ADR-0016 untuk uang, dan ketetapan modul untuk sisanya.

⛔ Bukan larangan memakai `JSON_DATAGUIDE` sama sekali — ia tetap dipakai sebagai penambal.
