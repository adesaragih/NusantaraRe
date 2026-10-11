# 72: Migrasi sekali jalan polis lama Pega → tabel case sistem baru (lanjutan tiket 22–24)

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 09-10-2026, gambar layar View
> Policy polis lama (RNM-F30.02.2026.00038 / NB-150929) yang tampil sebagai ringkasan FACINPRODUCTION: "kenapa berubah jadi
> gini? … HARUSNYA MUNCUL SEMUA". Pilihan work owner: **"Migrasi data lama sekali jalan"**.

**Status:** ⏸ DITUNDA — tahap A (analisis) selesai sesi c3 09-10-2026 (`72-analisis-tahap-a.md`, Q1–Q7, D1–D8). Work owner
09-10-2026: "migrasi data lama dilanjut nnti aja setelah smua cob selesai" (**K72-2**) — dilanjutkan sesudah port semua
Class of Business (tiket 73 dst.); belum ada tulisan data.

## Latar `[terverifikasi]`

- Polis yang lahir di Pega tidak punya baris di tabel case sistem baru (`T_WORK_POLIS` LINI FAC, `T_GENERAL_POLIS`,
  `T_QUOTATIONDATA`, tabel tab Object / Coverage / Spreading / Cedant / Payment / Scoring …); isi lengkapnya hanya di
  `POOLDATA.JSON_POLIS.DATA_JSON` (DDL `D:\migrasi\RNM\DDL\JSON_POLIS.txt`). Layar case / View Policy / View Details membaca
  tabel (K70-1, K71-1, ADR-0028), jadi polis lama hanya dapat diringkas dari `FACINPRODUCTION`.
- Jalur migrasi sudah dirancang: tiket 22 `loader.Flatten` (dibangun, diuji 115 contoh korpus), tiket 23 migrasi tabel
  flat (DITAHAN butir 66: presisi `NUMBER(38,8)` vs nilai Pega hingga 20 desimal / > 38 digit), tiket 24 pemuat
  repository + perkakas (menunggu 23). Sejak itu tabel `T_*` dibuat bertahap per tiket layar (migrasi 182–213).

## Keputusan work owner 09-10-2026

- **K72-1** Polis lama Pega dimigrasikan **sekali jalan** dari `JSON_POLIS.DATA_JSON` ke tabel case sistem baru, supaya
  semua polis tampil sebagai layar case lengkap (mode hanya lihat). Sistem tetap tidak membaca JSON saat berjalan (K70-1
  untuk runtime; JSON hanya sumber migrasi). Pola ADR-0009 (migrasi penuh, Claim Life).

## Tahap

- **A — analisis (sekarang):** kesenjangan keluaran `loader.Flatten` vs tabel `T_*` yang ADA (tabel / kolom hilang,
  tipe & presisi, kunci / ID, kolom sistem), cakupan polis (NB FIRE? semua PRODKE atau versi terakhir? endorsement?),
  pemetaan status case (Resolved-Completed, posisi), nomor case vs `SEQ_WORK_POLIS_NB`, idempotensi, volume, waktu.
- **B — keputusan work owner + DBA** atas temuan A (migrasi tabel tambahan, presisi, cakupan).
- **C — perkakas** (tiket 24): baca `JSON_POLIS` → `loader.Flatten` → tulis per polis satu transaksi; mode uji kering,
  laporan per polis, dapat diulang tanpa ganda. Dijalankan DBA / work owner, bukan oleh aplikasi saat berjalan.
