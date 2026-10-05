# Struktur tabel — Treaty Group

Modul `treatygroup` (keputusan work owner 05-10-2026) menulis tabel warisan `POOLDATA.TREATYGROUP` langsung; nol
tabel baru, nol DDL. Katalog DEV 05-10-2026: 33 baris, tanpa PK, tanpa constraint, tanpa trigger; dua indeks
non-unik (`ID`; `ID, TREATYGROUPNAME`); sequence `TREATYGROUP_SEQ` (nilai terakhir 70, cache 0) dipakai prosedur
`PEGA_TREATYGROUP` dan `PEGA_M_TREATYGROUP`.

Relasi (data DEV 05-10-2026):

- `OJKBUSINESSID` = `TREATYGROUPOJK.ID`; `OJKBUSINESSNAME`, `OJKBUSINESSNAMEIDN`, `ORDERNO` = salinan baris OJK itu
  (33 dari 33 cocok).
- `BUSINESSGROUP.TOPID` = `TREATYGROUP.ID` (induk grup bisnis; 36 dari 45 menunjuk baris tabel ini).
- `COAID` = `BUSINESSGROUP.ID` - arahnya berbalik dan bukan hierarki: sama untuk semua grup se-OJK (13 menunjuk
  anaknya sendiri, 15 anak grup lain se-OJK, 5 anak grup OJK lain). Tidak ada di korpus maupun prosedur Pega.
- `ID` dan salinan `TREATYGROUPNAME` dipakai ±25 tabel lain (TREATYBUSINESS, TREATYYEAR, TREATYINDETAIL,
  CLAIM_MASTER_TREATY, ...) - salinan nama itu tidak diubah modul ini.
- `M_TREATYGROUP` (JSON Pega, 43 baris, 33 cocok ID) tidak disentuh.

## TREATYGROUP

Tabel warisan (dinyatakan di `MODUL.md`, bukan dibuat migrasi).

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | VARCHAR2(100) | ya | | "ID" | situs aktif + `TREATYGROUP_SEQ` 4 digit (`10070` berikutnya di DEV) |
| `OLDID` | VARCHAR2(1000) | ya | | "Old ID" (View) | warisan; tidak ditulis |
| `OJKBUSINESSID` | VARCHAR2(1000) | ya | | "OJK Business" | `TREATYGROUPOJK.ID` |
| `OJKBUSINESSNAME` | VARCHAR2(1000) | ya | | "OJK Business" | salinan `TREATYGROUPOJK.NAME` |
| `OJKBUSINESSNAMEIDN` | VARCHAR2(1000) | ya | | View | salinan `TREATYGROUPOJK.NAMEIDN` |
| `TREATYGROUPNAME` | VARCHAR2(1000) | ya | | "Treaty Group Name" | wajib, huruf besar, tidak kembar |
| `TREATYGROUPSOANAME` | VARCHAR2(1000) | ya | | "SOA Name" | boleh kosong, huruf besar |
| `TGLUPDATE` | VARCHAR2(1000) | ya | | "Last Edited" (WIB) | teks format Pega `YYYYMMDDTHHMMSS.mmm GMT` |
| `USERID` | VARCHAR2(100) | ya | | "Last Edited By" | akun login |
| `COAID` | VARCHAR2(10) | ya | | "COA" (tidak diketik) | `BUSINESSGROUP.ID`; ikut OJK yang dipilih = COAID grup lain se-OJK yang paling sering |
| `ORDERNO` | VARCHAR2(10) | ya | | tidak tampil (urutan daftar) | salinan `TREATYGROUPOJK.ORDERNO` |

## Tabel yang dibaca saja

| Tabel | Kolom | Untuk |
| --- | --- | --- |
| `TREATYGROUPOJK` | `ID`, `NAME`, `NAMEIDN`, `ORDERNO` | pilihan OJK Business dan salinannya |
| `BUSINESSGROUP` | `ID`, `NOTE`, `ALIASNAME`, `TOPID` | nama COA; grup bisnis anak di View (tanpa SYARIAH) |
| `M_SITE_DATABASE` | `ID`, `CURRENT_SITE` | awalan ID baru (situs aktif `CURRENT_SITE` 1) |
