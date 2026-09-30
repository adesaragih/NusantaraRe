---
status: accepted
tanggal: 2026-09-23
sumber: spec PremiumList Life, keputusan work owner dengan data DBA
---

# Uang melintasi batas stored procedure sebagai teks, dan dikembalikan ke desimal di dalam aplikasi

`[data DBA]` Sebagian stored procedure penulis menerima **seluruh** parameternya bertipe teks,
**termasuk kolom uang**.

ADR-0003 dan ADR-0016 tetap berlaku. Yang berubah hanya **di mana** bentuk desimalnya hidup:

| Tempat | Bentuk uang |
| --- | --- |
| Di dalam aplikasi | **desimal berskala tetap** |
| Di kolom basis data | **desimal berskala tetap** |
| Melintasi batas procedure | **teks**, karena tanda tangannya menuntut demikian |

## Aturannya

1. Perubahan bentuk terjadi **tepat di batas pemanggilan**, bukan lebih awal. Nilai tidak disimpan
   sebagai teks di dalam aplikasi hanya karena nanti akan dikirim sebagai teks.
2. Perubahan bentuk itu **satu fungsi**, dipakai seluruh pemanggil - bukan diulang per pemanggil,
   sebab format pemisah desimal yang berbeda-beda adalah cara paling mudah merusak angka uang.
3. Nilai yang kembali dari procedure diubah balik ke desimal **sebelum** dipakai untuk apa pun.

## Akibat

1. Nol tipe pecahan biner di jalur mana pun, termasuk saat melintasi batas procedure.
2. Test memeriksa perjalanan pulang-pergi satu nilai uang lewat procedure: nilai yang kembali sama
   persis dengan yang dikirim.
