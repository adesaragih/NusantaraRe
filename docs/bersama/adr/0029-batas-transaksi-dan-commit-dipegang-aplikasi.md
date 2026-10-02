---
status: accepted
tanggal: 2026-09-23
sumber: spec Claim Prop, Claim Fac In, PremiumList Life, Endorsement Life - keputusan work owner
---

# Batas transaksi dipegang aplikasi; `COMMIT` tidak tertanam di teks SQL

`[keputusan work owner]` Aplikasi yang membuka dan menutup transaksi. **Nol `COMMIT` di dalam teks
SQL aplikasi.**

Di sistem lama `COMMIT` tertanam di dalam blok anonim di badan rule - `BEGIN ... ; COMMIT; END;`.
Di sistem baru ia **dikeluarkan oleh aplikasi**, sesudah panggilan.

Objek basis data yang dipanggil **tetap sama**. Logika penulisannya **tidak direplikasi** di
aplikasi - mereplikasinya berarti menebak isi yang tidak terbaca.

## Bahaya yang membuat urutan pemanggilan menjadi wajib

`[data DBA]` Sebagian stored procedure memasang penanganan galat yang menjalankan **pembatalan
transaksi di dalam dirinya sendiri**.

Pembatalan telanjang di Oracle membatalkan **seluruh** transaksi - termasuk apa pun yang aplikasi
tulis **sebelumnya** dalam transaksi yang sama. Pembatalan sampai titik simpan **tidak menolong**,
sebab yang dipasang bukan itu.

Karena itu:

| | Aturan |
| ---: | --- |
| 1 | Panggilan procedure semacam itu adalah **langkah terakhir sebelum penutupan transaksi** |
| 2 | **Tidak boleh ada pekerjaan penting yang masih menggantung** sebelumnya dalam transaksi yang sama |
| 3 | Aplikasi **memeriksa penanda keberhasilan** yang dikembalikan procedure, dan tidak menganggap ketiadaan galat sebagai keberhasilan |

## Akibat

1. Seluruh urutan penyimpanan satu objek kerja dibungkus **satu transaksi**.
2. Test memeriksa bahwa kegagalan di tengah tidak meninggalkan baris separuh jadi.
3. Pemeriksaan penanda keberhasilan adalah bagian dari jalur, bukan pemeriksaan tambahan yang
   boleh dilewati.
