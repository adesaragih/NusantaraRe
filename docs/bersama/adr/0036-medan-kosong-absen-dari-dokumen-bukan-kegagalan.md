---
status: accepted
tanggal: 2026-09-23
sumber: spec Master Product Name Life, keputusan work owner
---

# Medan kosong tidak hadir di dokumen JSON, dan ketidakhadirannya bukan kegagalan

`[terverifikasi]` Dokumen JSON sistem lama **tidak memuat medan yang nilainya kosong**. Medan itu
bukan bernilai kosong - ia **tidak ada sama sekali**.

`[keputusan work owner]` Pemuat migrasi membaca setiap medan dari dokumen dan menghasilkan
**kosong bila medan absen**. Ketidakhadiran **bukan kegagalan**.

## Kenapa ini perlu ditulis

Pemuat yang menganggap medan absen sebagai kesalahan akan menolak sebagian besar dokumen lama -
bukan karena dokumennya rusak, melainkan karena pengisiannya memang tidak lengkap sejak awal.

Sebaliknya, pemuat yang diam-diam melewati medan absen tanpa mencatat akan menyamarkan dokumen
yang benar-benar rusak.

## Aturannya

1. Medan absen menghasilkan nilai kosong, dan pemuatan **dilanjutkan**.
2. Medan yang hadir tetapi tidak dikenali masuk **penampung medan tak dikenal** (ADR-0023). Itu
   persoalan yang berbeda dan tidak boleh dicampur.
3. Kewajiban isi tidak ditegakkan oleh pemuat migrasi. Data lama masuk apa adanya; kewajiban isi
   berlaku bagi masukan baru (ADR-0027).

## Akibat

1. Test memuat dokumen yang kehilangan sebagian medan, dan memastikan pemuatan berhasil dengan
   kolom kosong - bukan gagal.
2. Cacah medan yang absen dicatat sebagai keterangan pemuatan, bukan sebagai galat.
