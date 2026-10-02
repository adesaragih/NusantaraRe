---
status: accepted
---

# Rekonsiliasi paralel run eksak, dicapai bertahap

Ukuran keberhasilan migrasi adalah **rekonsiliasi paralel run dengan nol selisih sampai digit
terakhir** — bukan toleransi. Karena seluruh jalur tulis produksi melewati 35 stored procedure
`POOLDATA` yang isinya tidak kita miliki, membaca kode saja tidak dapat membuktikan port-nya benar;
hanya membandingkan keluaran dua sistem atas masukan yang sama yang bisa. Paralel run **penuh** baru
mungkin setelah ekspor produksi tunggal tiba, sehingga sampai saat itu dijalankan bertahap.

## Considered Options

- **Toleransi seragam (mis. ±0,01)** — ditolak. Bila urutan operasi dan presisi per-langkah
  direproduksi apa adanya, `decimal` bersifat deterministik dan hasilnya **harus** identik. Selisih
  sekecil apa pun berarti ada salah-port, dan toleransi justru menyembunyikannya. Toleransi baru
  masuk akal bila pembulatan sengaja diseragamkan — dan itu perbaikan, bukan migrasi
  (`CLAUDE.md` §1).
- **Tanpa paralel run, cutover langsung dengan UAT** — ditolak. Tanpa isi 35 stored procedure, tidak
  ada dasar lain untuk membuktikan kebenaran jalur produksi.

## Consequences

Tahapan yang disepakati:

1. **Sekarang** — perhitungan murni saja (premi, spreading, nilai dasar akseptasi) atas masukan yang
   direkam. Tidak menyentuh alur, tidak butuh ekspor produksi.
2. **Setelah mesin akseptasi ditulis** (lihat ADR-0003) — per-modul, termasuk tangga akseptasi
   dengan fixture tabel limit.
3. **Setelah `ALL_SOURCE` 35 prosedur tiba** — kasus nyata end-to-end.

Konsekuensi yang mengikat porting: urutan operasi dan presisi pembulatan **tidak boleh** diubah
"agar lebih rapi". Lihat ADR-0005.
