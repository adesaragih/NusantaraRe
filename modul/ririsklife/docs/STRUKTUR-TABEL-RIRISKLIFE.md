# Struktur tabel — R/I Risk

Modul `ririsklife` (keputusan work owner 08-10-2026 K1-K4, `MODUL.md`). SATU tabel per jenis data: tabel Pega
`M_RIRISK_LIFE_SUMMARY` / `M_RIRISK_LIFE` **berganti nama** menjadi `RIRISK_LIFE_SUMMARY` / `RIRISK_LIFE` (nama view lama
yang dibuang) dan menjadi tabel flat tanpa JSONDATA - migrasi inti 935-940; nol tabel baru. Nama objek dan kolom ditulis
SEKALI di `backend/repository/rirk_tabel.go`.

| Objek | Jenis | Ditulis | Dibaca | Pembaca lain |
| --- | --- | --- | --- | --- |
| `RIRISK_LIFE_SUMMARY` | tabel (dulu `M_RIRISK_LIFE_SUMMARY`), PK `ID` (`SYS_C009367`), 4 kolom, 117 baris DEV | Save Add / Edit (keputusan work owner 08-10-2026), Simpan Upload (ringkasan baru; OPERATORID / MODIFIEDDATE), Detail Save, Delete | grid, nama kembar, ID terpakai | `masterproductnamelife` (`Choose R/I Risk`, `SELECT ID, USEDBY`) |
| `RIRISK_LIFE` | tabel (dulu `M_RIRISK_LIFE`), PK `PK_RIRISK_LIFE` (940), 8 kolom, 16.398 baris DEV | Detail tambah / EDIT, Edit nama ringkasan (salinan USEDBY), Delete, Simpan Upload | R/I RISK DETAIL, kembar, jumlah, ID terpakai | `premiumlistlife` (View R/I Risk, Hitung QR) |
| `M_RIRISK_LIFE_TEMP` | tabel staging unggah Pega (16.391 baris) | — | — | TIDAK disentuh (K3) |
| `M_SITE_DATABASE` | tabel warisan | — | site aktif (ID ringkasan) | `adjusterconsultant`, `ricommlife` |
| `M_RIRISK_LIFE_SUMMARY_SEQ` / `M_RIRISK_LIFE_SEQ` | sequence warisan (last 166 / 31723 DEV) | NEXTVAL | — | prosedur `PEGA_*` (INVALID sesudah 937/940, K3) |

**ID baru** (dua rumus prosedur Pega, ditiru persis): ringkasan = `site || LPAD(M_RIRISK_LIFE_SUMMARY_SEQ, 6, '0')`
(site = `M_SITE_DATABASE.ID` `CURRENT_SITE='1'`; ≤ 10 karakter); rincian = `'1' || LPAD(M_RIRISK_LIFE_SEQ, 5, '0')`
(huruf `1` tetap, ≤ 6 karakter - nomor 6 angka DITOLAK, LPAD Oracle memotongnya). `models.BentukID` /
`BentukIDRincian`, uji `TestBentukIDRingkasanRumusProsedur`, `TestBentukIDRincianRumusProsedur`.

## RIRISK_LIFE_SUMMARY

Tabel Pega yang BENTUKNYA diubah migrasi inti 935-937 (karena itu tidak dinyatakan "Tabel warisan" di `MODUL.md`).
Tabel di bab ini = kolom yang DIBUAT migrasi (`936_ririsk_life_summary_kolom.sql`, `TestKolomDDLCocokDenganStruktur`);
nama = kolom view lama; lebar = pola ricommlife 931 (K2). Diisi 937 dari JSONDATA apa adanya. Indeks
**`IX_RIRISK_LIFE_SUMMARY_NAMA` = `UPPER(TRIM(USEDBY))`** (937): pemeriksa nama kembar setiap Save dan Simpan Upload
(`SqlPemakaiNama`).

| Kolom | Tipe | DDL | Bukti lebar dan isi |
| --- | --- | --- | --- |
| `USEDBY` | teks | VARCHAR2(200) | **[terverifikasi data DEV 08-10-2026]** panjang maksimum 81 byte. "R/I RISK NAME" `.USEDBY` b1256 (wajib b1274, tanpa pyMax); lebar = ricommlife 931 dan `models.BatasNama` |
| `MODIFIEDDATE` | teks | VARCHAR2(50) | bentuk Pega `20191113T025753.044 GMT` (fakta WO); "MODIFY DATE" `.MODIFIEDDATE` b8358 |
| `OPERATORID` | teks | VARCHAR2(200) | **[terverifikasi data DEV 08-10-2026]** panjang maksimum 17 byte; Go menulis akun login. Nilai DEV tidak dikutip (nama orang) |

### Kolom warisan RIRISK_LIFE_SUMMARY (tidak dibuat migrasi mana pun)

| Kolom | Tipe | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(10) PK `SYS_C009367` (ikut RENAME) | site + LPAD(`M_RIRISK_LIFE_SUMMARY_SEQ`, 6) |
| `JSONDATA` | CLOB, `ENSURE_M_RIRISK_LIFE_SUMMARY_JSON` (33 byte) | **DIBUANG 937** (`CASCADE CONSTRAINTS`); kunci `USEDBY`, `MODIFIEDDATE`, `OPERATORID` di kolom; `pxObjClass` hanya di cadangan CSV P3 |

## RIRISK_LIFE

Tabel Pega yang BENTUKNYA diubah migrasi inti 938-940. Tabel di bab ini = kolom yang DIBUAT migrasi
(`939_ririsk_life_kolom.sql`) - hanya AGE; kolom lain sudah ada (warisan, di bawah). 940 MENIMPA SEMUA kolom view dari
JSONDATA (kolom datar warisan basi: IDUSEDBY NULL di 86 baris, MONTH hanya 1.912 - fakta WO), menambah
`PK_RIRISK_LIFE (ID)` (pola rincian ricommlife; nol ID ganda DEV), dan indeks IDUSEDBY hanya bila belum ada indeks
berkolom pertama IDUSEDBY (`INDEX4`, kolomnya tidak tercatat di repo - LANGKAH-WO (a) kueri 8).

| Kolom | Tipe | DDL | Bukti |
| --- | --- | --- | --- |
| `AGE` | teks | VARCHAR2(10) | **[terverifikasi data DEV 08-10-2026]** kunci `AGE` tidak ada di satu baris JSON pun (view mengeluarkan NULL). Tidak ada medan AGE di XML `InboxRIRisk` - modul tidak membaca / menulisnya. Tipe teks = saudaranya YEAR / MONTH / CONTRACT (VARCHAR2(10)) dan riratelife 929 |

### Kolom warisan RIRISK_LIFE (tidak dibuat migrasi mana pun)

Tipe warisan dipakai apa adanya (K2) - tidak ada ALTER / MODIFY.

| Kolom | Tipe | Isi dan bukti XML (`InboxRIRisk.xml`) |
| --- | --- | --- |
| `ID` | VARCHAR2(6), PK `PK_RIRISK_LIFE` sejak 940 | `.ID` b1347 disabled; `'1' || LPAD(M_RIRISK_LIFE_SEQ, 5)` |
| `IDUSEDBY` | VARCHAR2(100) | `TempIDUsedBy.ID` b1524 (tersembunyi `1=2` b1643); ID ringkasan (maks 7). 2 baris yatim (1000085, 1000087) - tidak dihapus |
| `USEDBY` | VARCHAR2(1000) | `TempIDUsedBy.USEDBY` b1727 disabled wajib; salinan nama ringkasan |
| `YEAR` | VARCHAR2(10) | `.YEAR` b2202 pxTextInput TIDAK wajib, pyMax 4 b2220; NULL di 1.994 baris; tahun polis (bukan kalender) |
| `MONTH` | VARCHAR2(10) | `.MONTH` b2389 pxTextInput TIDAK wajib, pyMax 4 b2407; terisi 1.998 baris (0-180) |
| `RISK` | NUMBER (tanpa presisi / skala) | `.RISK` b2579 pxNumber wajib, label "RISK (PERMIL)" b2548. Dari JSON teks: 11.063 berkoma, 89 bertitik (maks 12 desimal, `580.894351210924`) - NUMBER tanpa skala supaya tidak terpotong. Konversi 940 tanpa NLS; tulis Go `TO_NUMBER(:koef) / POWER(10, :skala)`, baca `TO_CHAR(… 'TM9', NLS eksplisit)` |
| `CONTRACT` | VARCHAR2(10) | `.CONTRACT` b1920 pxTextInput wajib (`ChangeDotToPoint_DT` b1963); seluruhnya angka (DEV) |
| `JSONDATA` | CLOB, `ENSURE_M_RIRISK_LIFE_JSON` | **DIBUANG 940** (`CASCADE CONSTRAINTS`); `pxObjClass` hanya di cadangan CSV P3 |

## Riwayat — view yang dibuang

- `RIRISK_LIFE_SUMMARY` (VIEW): `SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID FROM
  POOLDATA.M_RIRISK_LIFE_SUMMARY a` - dibuang 935, dibuat ulang hanya oleh 935_down.
- `RIRISK_LIFE` (VIEW): `SELECT a.ID, a.JSONDATA.IDUSEDBY, a.JSONDATA.USEDBY, a.JSONDATA.AGE, a.JSONDATA.YEAR,
  a.JSONDATA.MONTH, a.JSONDATA.RISK, a.JSONDATA.CONTRACT FROM M_RIRISK_LIFE a` - dibuang 938, dibuat ulang hanya oleh
  938_down.
