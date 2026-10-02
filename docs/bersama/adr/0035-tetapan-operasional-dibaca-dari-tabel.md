---
status: accepted
tanggal: 2026-09-23
sumber: spec PremiumList Life dan spec Claim Prop, keputusan work owner
---

# Tetapan operasional dibaca dari tabel, bukan ditanam di kode

Tanggal, ambang, dan tetapan operasional lain **dibaca dari tabel** yang dipelihara pengguna -
bukan ditulis tetap di dalam kode atau di dalam teks query.

`[terverifikasi]` Contoh yang menjadi dasar keputusan: tanggal tutup buku dibaca dari tabel
tetapan. Bila tanggal berjalan melewatinya, cap waktu produksi digeser ke awal bulan berikutnya.

## Kenapa

1. Tetapan semacam ini **berubah menurut jadwal dagang**, bukan menurut jadwal rilis. Menanamnya
   di kode berarti setiap pergeseran tanggal menuntut penerbitan versi baru.
2. Orang yang tahu nilainya adalah pengguna, bukan pengembang.

## Batasnya

Ini **bukan** izin memindahkan aturan dagang ke tabel. Yang dibaca dari tabel adalah **nilainya**;
aturan yang memakai nilai itu tetap hidup di kode, tempat ia dapat diuji.

`[terverifikasi]` Di sistem lama ditemukan juga tetapan yang ditanam sebagai nama orang dan nomor
tertentu di dalam teks query. Semua yang seperti itu **dibuang** - lihat ADR-0030.

## Akibat

1. Nol tanggal, ambang, atau nama tetap di dalam kode maupun teks query.
2. Perilaku saat tetapan **belum terisi** ditetapkan eksplisit, bukan dibiarkan menjadi galat yang
   tidak terbaca.
3. Test memakai nilai tetapan yang disuntikkan, bukan nilai produksi.
