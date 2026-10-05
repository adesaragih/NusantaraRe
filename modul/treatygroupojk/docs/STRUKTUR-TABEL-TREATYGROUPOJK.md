# Struktur tabel — Treaty Group OJK

Modul `treatygroupojk` (keputusan work owner 05-10-2026) menulis tabel warisan `POOLDATA.TREATYGROUPOJK` langsung;
nol tabel baru, nol DDL. Katalog DEV 05-10-2026: 16 baris, tanpa PK, tanpa constraint, tanpa indeks, tanpa trigger,
tanpa sequence. Pemakai: `TREATYGROUP.OJKBUSINESSID` menunjuk `ID`-nya, dan `TREATYGROUP` menyimpan SALINAN `NAME`,
`NAMEIDN`, `ORDERNO` (33 dari 33 baris DEV cocok) - salinan itu tidak diubah modul ini (keputusan work owner
05-10-2026). Di Pega, `FetchTreatyGroupOJK` hanya membacanya.

## TREATYGROUPOJK

Tabel warisan (dinyatakan di `MODUL.md`, bukan dibuat migrasi).

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | VARCHAR2(50) | ya | | "ID" | `01`..`16`; baris baru = nomor tertinggi + 1, dua digit |
| `NAME` | VARCHAR2(50) | ya | | "Name" | bahasa Inggris (`PROPERTY`); wajib, huruf besar, tidak kembar |
| `NAMEIDN` | VARCHAR2(50) | ya | | "Name (IDN)" | bahasa Indonesia (`HARTA BENDA`); wajib, huruf besar |
| `ORDERNO` | VARCHAR2(50) | ya | | tidak tampil (urutan daftar) | urutan tampil (`1`..`17`, teks angka); Add = tertinggi + 1, Edit tetap |
