> Modul  : Treaty In Adjustment · Ronde D · 2026-09-24

# 04 · PARITAS

**Berlaku sebagian.** Paritas menanyakan: apa yang sistem lama dapat lakukan, dan apakah sistem baru
masih dapat melakukannya. Untuk ronde ini pertanyaannya **terbalik pada satu sisi** — sebab salah
satu kemampuan yang diputuskan **tidak pernah ada di sistem lama**.

## 1. `GRL-19` — paritasnya NEGATIF, dan itu keputusan sadar

| Pertanyaan paritas | Jawaban |
|---|---|
| Apa yang sistem lama dapat lakukan dengan dokumen addendum? | **tidak ada.** `TD-01` — nol properti menyimpan nomornya |
| Apakah sistem baru kehilangan sesuatu? | **tidak.** Tidak ada yang hilang dari nol |
| Apakah sistem baru menambah sesuatu? | **ya** — dan penambahan itu **BARU**, bukan pelestarian |

> **Paritas tidak dilanggar; ia tidak berlaku.** Menuliskannya sebagai *"pelestarian"* akan
> menyembunyikan bahwa migrasi **tidak punya sumber** untuk kolom ini.

**Yang justru harus diperiksa paritasnya adalah arah sebaliknya:** apakah ada kemampuan lama yang
**hilang** karena dokumen naik menjadi entitas? **Tidak.** Versi tetap dapat berdiri tanpa dokumen
(butir a), dan persetujuan tetap per versi (butir 3).

## 2. `GRL-20` — paritasnya POSITIF, dan sistem baru LEBIH KUAT

| Kemampuan lama | Di sistem baru |
|---|---|
| memilih material / non material sebelum menyunting | **tetap ada** — `SIFAT_MATERIAL_ADDENDUM` masukan |
| penguncian field menurut pilihan itu | **tetap ada**, dan **ditegakkan saat simpan**, bukan hanya di layar |
| pembedaan dua arah — Material mengunci teks, Non Material mengunci angka | **tetap ada**, dan **baru sekarang tertulis**; sebelumnya tidak pernah dinyatakan siapa pun |

**Satu kemampuan lama yang sengaja TIDAK dilestarikan:** kemampuan **melewati** penguncian itu.
Di sistem lama aturannya hanya di layar (`TDA-10`), sehingga jalur simpan yang tidak melewati layar
dapat menembusnya. Itu **cacat, bukan kemampuan** — dan ia yang membuat `UA-3` perlu dijalankan
sebelum invariannya ditegakkan atas data lama.
