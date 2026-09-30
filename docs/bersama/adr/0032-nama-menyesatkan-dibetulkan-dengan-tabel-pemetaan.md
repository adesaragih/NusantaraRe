---
status: accepted
tanggal: 2026-09-23
sumber: spec Claim Fac In Bab 18 dan spec Claim Prop - keputusan work owner atas rekomendasi asisten
---

# Nama yang menyesatkan dibetulkan, disertai tabel pemetaan nama lama ke nama benar

`[keputusan work owner]` Kolom hasil pembacaan **diberi nama sesuai isinya**. Nama lama yang
menyesatkan tidak diwariskan.

Setiap pembetulan **wajib disertai tabel pemetaan** nama lama ke nama benar.

## Masalah yang ditemukan

`[terverifikasi]` Dua rule pembaca riwayat memberi nama samaran pada kolom hasilnya, dan
**memberi nama samaran dengan pemetaan yang berbeda satu sama lain**. Tanggal kejadian dinamai
*tanggal mulai*; nomor klaim dinamai *nama cabang*; penyebab kerugian dinamai *kode bisnis* - dan
pada rule kedua, pemetaannya lain lagi.

`[terverifikasi]` Pemanggilnya satu-satu, dan **nol pemanggil memakai keduanya**, sehingga kedua
pemetaan yang menyesatkan itu tidak pernah bertemu dalam satu jalur. Tetapi **keduanya membaca
tabel yang sama**.

## Kenapa tabel pemetaan itu wajib

Tanpa tabel pemetaan, orang yang membandingkan keluaran lama dan keluaran baru akan mengira
**datanya berubah**, padahal hanya **namanya** yang dibetulkan.

Itu kekeliruan yang mahal: ia memicu penyelidikan atas kerusakan yang tidak pernah terjadi, dan
pada kasus terburuk memicu pembatalan migrasi yang sebenarnya benar.

## Akibat

1. Nama kolom di batas pembacaan mengikuti isinya, bukan mengikuti nama samaran lama.
2. Tabel pemetaan disimpan bersama spec modul yang bersangkutan, bukan di dalam kode.
3. `[penyimpangan sadar]` Pembetulan nama adalah penyimpangan yang disadari dari perilaku lama, dan
   dicatat sebagai penyimpangan.
4. Berapa tepatnya nama samaran yang menyesatkan **belum tentu diketahui seluruhnya**. Setiap yang
   ditemukan kemudian ditambahkan ke tabel pemetaan yang sama.
