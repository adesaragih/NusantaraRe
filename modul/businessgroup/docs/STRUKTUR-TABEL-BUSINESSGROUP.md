# Struktur tabel — Business Group

Modul `businessgroup` (keputusan work owner 05-10-2026) menulis tabel warisan `POOLDATA.BUSINESSGROUP` langsung; nol
tabel baru, nol DDL. Katalog DEV 05-10-2026: 45 baris (`10001`..`10045`), tanpa PK, tanpa constraint, tanpa indeks,
tanpa trigger; sequence `BUSINESSGROUP_SEQ` (nilai terakhir 46) dipakai prosedur `UPSERT_BUSINESSGROUP`. Objek yang
membacanya: view `BUSINESS_VIEW`, `V_BUSINESSGROUP`, prosedur `RDBMASTERBUSINESSGROUP`.

Relasi (data DEV 05-10-2026):

- `TOPID` = `TREATYGROUP.ID` (induk); `TREATYNAME` = salinan `TREATYGROUPNAME` induknya (36 dari 36 yang induknya ada
  cocok). 9 baris ber-TOPID yang tidak ada di `TREATYGROUP` (8 di antaranya ada di `M_TREATYGROUP` JSON Pega).
- `BUSINESS.BUSINESSGROUPID` menunjuk `ID` (29 grup dipakai) dan menyalin namanya ke `BUSINESSGROUPNAME` - salinan itu
  tidak diubah modul ini. `TREATYGROUP.COAID` juga menunjuk `ID` (16 grup).
- 10 baris berakhiran SYARIAH - tidak dikelola modul ini.
- `M_BUSINESSGROUP` (JSON Pega, 45 baris, 45 cocok ID) tidak disentuh.

## BUSINESSGROUP

Tabel warisan (dinyatakan di `MODUL.md` Accounts, bukan dibuat migrasi).

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | VARCHAR2(4000) | ya | | "ID" | situs aktif + `BUSINESSGROUP_SEQ` 4 digit (`10046` berikutnya di DEV) |
| `NOTE` | VARCHAR2(4000) | ya | | "Name" | wajib, huruf besar, tidak kembar; tidak berakhiran SYARIAH |
| `TOPID` | VARCHAR2(4000) | ya | | "Treaty Group" | `TREATYGROUP.ID` |
| `ALIASNAME` | VARCHAR2(4000) | ya | | "Alias Name" | huruf besar; kosong = Name (44 dari 45 baris DEV sama dengan NOTE) |
| `TREATYNAME` | VARCHAR2(4000) | ya | | "Treaty Group" | salinan `TREATYGROUP.TREATYGROUPNAME` saat Treaty Group dipilih |

## Tabel yang dibaca saja

| Tabel | Kolom | Untuk |
| --- | --- | --- |
| `TREATYGROUP` | `ID`, `TREATYGROUPNAME` | pilihan Treaty Group dan salinan namanya |
| `M_SITE_DATABASE` | `ID`, `CURRENT_SITE` | awalan ID baru (situs aktif `CURRENT_SITE` 1) |
