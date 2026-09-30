---
status: accepted
tanggal: 2026-09-23
sumber: spec Claim Life dan spec penyimpanan modul-modul lain, keputusan work owner
---

# Kolom basis data nullable; kewajiban isi ditegakkan di kode

`[keputusan work owner]` **Seluruh kolom dideklarasi nullable.** Kewajiban mengisi ditegakkan di
lapisan aplikasi, bukan oleh batasan `NOT NULL`.

`[penyimpangan sadar]` Ini **penyimpangan yang disadari** dari kebiasaan umum, dan dicatat sebagai
penyimpangan, bukan disamarkan sebagai praktik terbaik.

## Kenapa

1. **Data lama tidak lengkap.** Migrasi membawa baris dari sistem yang tidak pernah mewajibkan
   pengisian. Batasan `NOT NULL` akan menolak data historis yang sah, dan memaksa mengarang nilai
   pengganti - yang lebih berbahaya daripada kosong.
2. **Kewajiban isi bergantung keadaan.** Sebuah kolom wajib pada satu langkah alur dan tidak pada
   langkah lain. Batasan basis data tidak mengenal langkah alur; kode mengenalnya.
3. Pesan galat dari aplikasi dapat dibaca pemakai. Pesan galat dari batasan basis data tidak.

## Harga yang dibayar

Basis data **tidak lagi menjadi jaring pengaman terakhir**. Bila kode lalai, kolom kosong akan
tersimpan tanpa ada yang menolak.

Karena itu dua hal menjadi wajib:

1. Pemeriksaan kewajiban isi terletak di **satu tempat** per objek, bukan tersebar di banyak jalur.
2. Test memeriksa penolakan itu secara langsung - bukan mengandalkan basis data menolaknya.

## Akibat

1. Nol `NOT NULL` pada kolom selain kunci utama.
2. Pemuat migrasi tidak perlu mengarang nilai pengganti untuk data historis yang kosong.
3. Test yang menemukan kolom wajib-isi lolos tersimpan kosong lewat lapisan layanan **gagal**.
