# Struktur tabel — R/I Comm Life

Modul `ricommlife` (perintah work owner 06-10-2026). **RALAT R1 (keputusan work owner 08-10-2026, `MODUL.md`)**: SATU
tabel per jenis data, pola R/I Rate Life - ringkasan = kolom **`M_RICOMM_LIFE_SUMMARY`** (migrasi inti `931` + `932`),
rincian = kolom **`M_RICOMM_LIFE`** (`933` + `934`). JSONDATA kedua tabel, view `RICOMM_LIFE_SUMMARY`, dan tabel flat
`RICOMM_LIFE` (924) dibuang; nol tabel baru. Nama objek dan kolom ditulis SEKALI di `backend/repository/ricl_tabel.go`.

| Objek | Jenis | Ditulis | Dibaca | Pembaca lain |
| --- | --- | --- | --- | --- |
| `M_RICOMM_LIFE_SUMMARY` | tabel, PK `ID`, 4 kolom (sesudah 932), 3 baris DEV 08-10-2026 | Add, Edit, Delete, Detail Save / Simpan Upload (OPERATORID, MODIFIEDDATE) | grid, nama kembar, ID terpakai | — (prosedur `PEGA_M_RICOMM_LIFE_SUMMARY` INVALID sesudah 932 - diterima WO) |
| `M_RICOMM_LIFE` | tabel, PK `ID`, 6 kolom (sesudah 934), 2 baris DEV (dari `RICOMM_LIFE`) | Detail tambah / ubah, Edit nama (salinan USEDBY), Delete, Simpan Upload | R/I COMM DETAIL, kembar, jumlah, ID terpakai | — (prosedur `PEGA_M_RICOMM_LIFE` INVALID sesudah 934 - diterima WO) |
| `RICOMM_LIFE` | tabel flat 924 - **DIBUANG 934** (riwayat di bawah) | — | — | — |
| `RICOMM_LIFE_SUMMARY` | view warisan - **DIBUANG 932** (riwayat di bawah) | — | — | — |
| `M_SITE_DATABASE` | tabel warisan (dinyatakan `adjusterconsultant`) | — | site aktif `CURRENT_SITE = '1'` (awalan ID) | `adjusterconsultant` |
| `M_RICOMM_LIFE_SUMMARY_SEQ` | sequence warisan (last 4 DEV) | NEXTVAL | — | `PEGA_M_RICOMM_LIFE_SUMMARY` |
| `M_RICOMM_LIFE_SEQ` | sequence warisan (last 43 DEV) | NEXTVAL | — | `PEGA_M_RICOMM_LIFE` |

**ID baru** (ringkasan dan rincian, butir 4, TETAP) = rumus prosedur `PEGA_M_RICOMM_LIFE` b13/b21:
`site || LPAD(seq.NEXTVAL, 6, '0')`, site = `M_SITE_DATABASE.ID` baris `CURRENT_SITE = '1'` dibaca saat berjalan
(galat bila tidak tepat satu baris). DEV: site 1 → ID 7 karakter (`1000044`). Nomor lebih dari 6 angka atau ID lebih
dari 10 karakter = galat. ID terpakai diperiksa di SATU tabel per jenis. `models.BentukID`, uji `TestBentukIDRumusProsedur`.

## M_RICOMM_LIFE_SUMMARY

Tabel Pega yang BENTUKNYA diubah migrasi inti 931/932 (karena itu tidak lagi dinyatakan "Tabel warisan", preseden
`adjusterconsultant` 870 dan riratelife 927). Tabel di bab ini = kolom yang DIBUAT migrasi
(`931_m_ricomm_life_summary_kolom.sql`, `TestKolomDDLCocokDenganStruktur`); nama = kolom view `RICOMM_LIFE_SUMMARY`
lama. Diisi 932 dari JSONDATA APA ADANYA; semua NULLABLE, wajib-isi ditegakkan Go. Indeks
**`IX_M_RICOMM_LIFE_SUMMARY_NAMA` = `UPPER(TRIM(USEDBY))`** (932): pemeriksa nama kembar (`SqlPemakaiNama`) setiap Save
dan Simpan Upload.

| Kolom | Tipe | DDL | Bukti lebar dan isi |
| --- | --- | --- | --- |
| `USEDBY` | teks | VARCHAR2(200) | **[terverifikasi data DEV 08-10-2026]** panjang maksimum 18 byte. "R/I COMM NAME" `.USEDBY` b1274 (tanpa pyMax). Lebar 200 = salinan nama di rincian (`M_RICOMM_LIFE.USEDBY` 933 = `RICOMM_LIFE.USEDBY` 924) dan `models.BatasNama` - nama ringkasan selalu muat di rinciannya (bukan 500 milik R/I Rate, yang tidak dibatasi salinan 200) |
| `MODIFIEDDATE` | teks | VARCHAR2(50) | bentuk Pega `YYYYMMDDTHHMMSS.mmm GMT` (23 byte; baris DEV 1000003 `20181205T073755.559 GMT`); lebar = riratelife 927. "MODIFY DATE" `.MODIFIEDDATE` b8526 |
| `OPERATORID` | teks | VARCHAR2(200) | lebar = riratelife 927; Go menulis akun login (dipotong `BatasNama` 200 byte). "MODIFY OPERATOR" `.OPERATORID` b8348. Nilai DEV tidak dikutip di dokumen ini (nama orang) |

### Kolom warisan M_RICOMM_LIFE_SUMMARY (tidak dibuat migrasi mana pun)

| Kolom | Tipe | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(10) PK (`SYS_C008861`) | **[terverifikasi data DEV 08-10-2026]** panjang maksimum 7 byte; site + LPAD(`M_RICOMM_LIFE_SUMMARY_SEQ`, 6) |
| `JSONDATA` | CLOB, constraint `ENSURE_M_RICOMM_LIFE_SUMMARY_JSON` (IS JSON, nama 33 byte) | **DIBUANG 932** (`CASCADE CONSTRAINTS`). Kunci `USEDBY`, `MODIFIEDDATE`, `OPERATORID` sudah di kolom; `pxObjClass` (bukan kolom view) hanya di cadangan CSV `ID + JSONDATA` (`LANGKAH-WO-RICOMMLIFE-SATU-TABEL.md` (b)) |

## M_RICOMM_LIFE

Tabel Pega yang BENTUKNYA diubah migrasi inti 933/934 (tidak lagi "Tabel warisan"). Tabel di bab ini = kolom yang DIBUAT
migrasi (`933_m_ricomm_life_kolom.sql`); tipe dan lebar = tabel flat `RICOMM_LIFE` 924 (bukti XML di bawah, tidak
berubah). Diisi 934 dari `RICOMM_LIFE` (sumber kebenaran; baris JSON yang ada diisi dari JSON lebih dulu - DEV 0 baris).
Semua NULLABLE kecuali PK warisan; wajib-isi ditegakkan Go (`models.PeriksaIsianKomisi`). **Indeks
`IX_M_RICOMM_LIFE_IDUSEDBY` (IDUSEDBY)** (934, pengganti `IX_RICOMM_LIFE_IDUSEDBY`): grid R/I COMM DETAIL per ringkasan
(`BrowseRICommLife_RD`, param idusedby b7396), kembar (CONTRACT, YEAR), `UPDATE … USEDBY` saat nama ringkasan diubah,
Delete ringkasan beserta rinciannya.

| Kolom | Tipe | DDL | Bukti XML (`InboxRIComm.xml`, kelas `ASM-FW-GISFW-Int-RI_COMM_LIFE`) dan aturan |
| --- | --- | --- | --- |
| `IDUSEDBY` | teks | VARCHAR2(10) | **[terverifikasi data DEV 08-10-2026]** panjang maksimum 7 byte. `TempIDUsedBy.ID` b1509 (label b1478), tersembunyi `1=2` b1633; = `ID` ringkasan |
| `USEDBY` | teks | VARCHAR2(200) | **[terverifikasi data DEV 08-10-2026]** panjang maksimum 18 byte. `TempIDUsedBy.USEDBY` b1718, disabled b1737-b1738, wajib b1729/b1736; grid "R/I COMM NAME" b7665. **[penyimpangan sadar]** XML tanpa pyMax → batas nama ringkasan 200 byte (924) |
| `CONTRACT` | angka bulat | NUMBER(5) | `.CONTRACT` b1904, **pxNumber** b1907, wajib b1917/b1923, placeholder `0` b1927. **[penyimpangan sadar]** XML tanpa presisi → K6 NUMBER(5); Go: bulat 0-99999 [penyimpangan sadar - menunggu WO] |
| `YEAR` | angka bulat | NUMBER(5) | `.YEAR` b2185, wajib b2196/b2201, placeholder `0` b2205, **pyMax 4** b2206. **[penyimpangan sadar]** K6 tahun NUMBER(5); Go: TEPAT 4 angka [penyimpangan sadar - menunggu WO] |
| `COMM` | angka desimal | NUMBER(38,8) | `.COMM` b2392, wajib b2406/b2412; grid **pxNumber** b9101. **[penyimpangan sadar]** K6 NUMBER(38,8); Go: desimal tidak negatif, koma atau titik, ≤ 30 + 8 angka, tidak pernah dibulatkan |

Angka ditulis tanpa bergantung NLS sesi (`TO_NUMBER(:n)` atas teks angka; COMM `TO_NUMBER(:koef) / POWER(10, :skala)`)
dan dibaca `TO_CHAR(..., 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')` → API memakai titik desimal kanonik (`0.5`); layar
menampilkan koma desimal.

### Kolom warisan M_RICOMM_LIFE (tidak dibuat migrasi mana pun)

| Kolom | Tipe | Isi |
| --- | --- | --- |
| `ID` | VARCHAR2(10) PK (`SYS_C008737`) | **[terverifikasi data DEV 08-10-2026]** panjang maksimum 7 byte; site + LPAD(`M_RICOMM_LIFE_SEQ`, 6) |
| `JSONDATA` | CLOB, constraint `ENSURE_M_RICOMM_LIFE_JSON` (IS JSON) | **DIBUANG 934** (`CASCADE CONSTRAINTS`); 0 baris DEV, kunci `IDUSEDBY`, `USEDBY`, `CONTRACT`, `YEAR`, `COMM` |

## RICOMM_LIFE

RIWAYAT - tabel flat migrasi inti `924_ricomm_life.sql` (sudah jalan di DEV), **DIBUANG**
`934_m_ricomm_life_satu_tabel.sql` (RALAT R1) sesudah isinya disalin ke `M_RICOMM_LIFE`; dibangun ulang PERSIS (kolom,
`PK_RICOMM_LIFE`, `IX_RICOMM_LIFE_IDUSEDBY`) hanya oleh jalur mundur 934. Constraint CHECK sistem `SYS_C0015437` =
`ID` NOT NULL dari DDL 924 (satu-satunya NOT NULL di DDL itu; tidak disalin - ID `M_RICOMM_LIFE` sudah PK);
`LANGKAH-WO-RICOMMLIFE-SATU-TABEL.md` (a) mencatat SEARCH_CONDITION-nya untuk dipastikan WO.

| Kolom | Tipe | DDL | Isi |
| --- | --- | --- | --- |
| `ID` | teks | VARCHAR2(10) NOT NULL, PK | kini `M_RICOMM_LIFE.ID` |
| `IDUSEDBY` | teks | VARCHAR2(10) | kini `M_RICOMM_LIFE.IDUSEDBY` |
| `USEDBY` | teks | VARCHAR2(200) | kini `M_RICOMM_LIFE.USEDBY` |
| `CONTRACT` | angka bulat | NUMBER(5) | kini `M_RICOMM_LIFE.CONTRACT` |
| `YEAR` | angka bulat | NUMBER(5) | kini `M_RICOMM_LIFE.YEAR` |
| `COMM` | angka desimal | NUMBER(38,8) | kini `M_RICOMM_LIFE.COMM` |

## RICOMM_LIFE_SUMMARY

RIWAYAT - view warisan Pega `SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID FROM
POOLDATA.M_RICOMM_LIFE_SUMMARY a`, **DIBUANG** `932_m_ricomm_life_summary_satu_tabel.sql` (RALAT R1); dibangun ulang
hanya oleh jalur mundur 932. Kolomnya kini kolom `M_RICOMM_LIFE_SUMMARY` di atas.
