# ADR-0045 — Jejak perubahan adalah fakta mesin, terpisah dari catatan manusia

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In, gelombang 1

## Konteks

Satu-satunya jejak di sistem lama adalah daftar komentar, berisi tanggal, pengenal operator, status
akseptasi, dan teks bebas. Ia ditulis **hanya pada perpindahan persetujuan**. Penyimpanan biasa —
termasuk yang mengubah angka — tidak meninggalkan jejak siapa maupun kapan. Tidak ada kolom waktu
sama sekali di tabel kontrak maupun tabel addendum.

## Keputusan

Jejak perubahan masuk **gelombang 1**, dan alasannya bukan kepatuhan:

> **Jejak perubahan adalah bukti bahwa pembekuan benar-benar terjadi.**

ADR-0036 menetapkan angka yang disetujui dibekukan. Pemilik proses menetapkan Adjustment mengambil
data lama lewat SELECT, bukan menyalinnya. Keduanya hanya bisa dipercaya bila ada cara membuktikan
data lama tidak berubah di antara dua pembacaan. Tanpa jejak, pembekuan adalah **janji**, bukan
fakta — dan selisih yang dibukukan Adjustment berdiri di atas janji.

**Isi minimal:** siapa, kapan, apa yang berubah dari nilai apa ke nilai apa, dan **di bawah peran
apa** — peran yang berlaku saat itu, bukan peran orangnya hari ini. Butir terakhir menyambung ke
penugasan peran bertanggal (ADR-0044), dan itu sebabnya keduanya harus lahir bersamaan.

## Jangan digabung dengan catatan komentar

Keduanya benda berbeda:

| | Catatan komentar | Jejak perubahan |
|---|---|---|
| Isinya | narasi manusia: alasan, pertimbangan | fakta mesin |
| Kapan | pada perpindahan persetujuan | pada setiap penulisan |
| Ditulis oleh | orang | sistem |
| Bisa disunting | ya | **tidak** |
| Teks bebas | ya | tidak ada |

Menggabungkannya menghasilkan yang terburuk dari keduanya: catatan yang bisa disunting dan karena
itu tidak membuktikan apa-apa.
