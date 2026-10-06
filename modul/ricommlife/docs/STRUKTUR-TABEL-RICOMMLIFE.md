# Struktur tabel — R/I Comm Life

Modul `ricommlife` (perintah work owner 06-10-2026) menulis ringkasan ke tabel JSON warisan `M_RICOMM_LIFE_SUMMARY`
(dibaca lewat view `RICOMM_LIFE_SUMMARY`) dan rinciannya ke **tabel flat baru `RICOMM_LIFE`** (migrasi inti `924`),
yang menggantikan VIEW warisan bernama sama. Keputusan work owner 06-10-2026 butir 1-9: `MODUL.md`. Nama objek, kolom,
dan kunci JSON ditulis SEKALI di `backend/repository/ricl_tabel.go`.

| Objek | Jenis | Ditulis | Dibaca | Pembaca lain |
| --- | --- | --- | --- | --- |
| `RICOMM_LIFE` | **tabel flat baru** (924), PK `ID` | Detail tambah / ubah, Edit nama (salinan USEDBY), Delete, Simpan Upload, alat pindahflat | R/I COMM DETAIL, kembar, jumlah | — (nol kode repo membaca view lamanya, dicek WO) |
| `M_RICOMM_LIFE_SUMMARY` | tabel JSON warisan, 1 baris DEV | Add, Edit, Delete, Detail Save / Simpan Upload (OPERATORID, MODIFIEDDATE) | ID terpakai | prosedur `PEGA_M_RICOMM_LIFE_SUMMARY` |
| `RICOMM_LIFE_SUMMARY` | view warisan, 4 kolom | — | grid, nama kembar | — |
| `M_RICOMM_LIFE` | tabel JSON warisan, 0 baris DEV | **tidak pernah** (butir 3) | alat pindahflat saja | prosedur `PEGA_M_RICOMM_LIFE` |
| `M_SITE_DATABASE` | tabel warisan (dinyatakan `adjusterconsultant`) | — | site aktif `CURRENT_SITE = '1'` (awalan ID) | `adjusterconsultant` |
| `M_RICOMM_LIFE_SUMMARY_SEQ` | sequence warisan (last 4 DEV) | NEXTVAL | — | `PEGA_M_RICOMM_LIFE_SUMMARY` |
| `M_RICOMM_LIFE_SEQ` | sequence warisan (last 43 DEV) | NEXTVAL | — | `PEGA_M_RICOMM_LIFE` |

**ID baru** (ringkasan dan rincian, butir 4) = rumus prosedur `PEGA_M_RICOMM_LIFE` b13/b21:
`site || LPAD(seq.NEXTVAL, 6, '0')`, site = `M_SITE_DATABASE.ID` baris `CURRENT_SITE = '1'` dibaca saat berjalan
(galat bila tidak tepat satu baris). DEV: site 1 → ID 7 karakter (`1000044`). Nomor lebih dari 6 angka atau ID lebih
dari 10 karakter = galat (LPAD Oracle akan memotongnya diam-diam). `models.BentukID`, uji `TestBentukIDRumusProsedur`.

## RICOMM_LIFE

Tabel flat (migrasi inti `924_ricomm_life.sql`). Satu baris = satu rincian komisi milik satu ringkasan. Kolom = PERSIS
kolom VIEW warisan `RICOMM_LIFE` (`SELECT a.ID, a.JSONDATA.IDUSEDBY, a.JSONDATA.USEDBY, a.JSONDATA.CONTRACT,
a.JSONDATA.YEAR, a.JSONDATA.COMM FROM M_RICOMM_LIFE a`). Semua NULLABLE kecuali PK; wajib-isi ditegakkan Go
(`models.PeriksaIsianKomisi`). Indeks `IX_RICOMM_LIFE_IDUSEDBY` (grid, kembar, Delete per ringkasan).

| Kolom | Tipe | DDL | Bukti XML (`InboxRIComm.xml`, kelas `ASM-FW-GISFW-Int-RI_COMM_LIFE`) dan aturan |
| --- | --- | --- | --- |
| `ID` | teks | VARCHAR2(10) NOT NULL, PK | `.ID` b1325 pxTextInput, disabled (label b1294); brief butir 1: VARCHAR2(10); site + LPAD 6 |
| `IDUSEDBY` | teks | VARCHAR2(10) | `TempIDUsedBy.ID` b1509 (label b1478), tersembunyi `1=2` b1633; = `ID` ringkasan (view `RICOMM_LIFE_SUMMARY.ID`); brief butir 1: VARCHAR2(10) |
| `USEDBY` | teks | VARCHAR2(200) | `TempIDUsedBy.USEDBY` b1718 (label b1687), disabled b1737-b1738, wajib b1729/b1736; grid "R/I COMM NAME" (`InboxRIComm` b7665). **[penyimpangan sadar]** XML tanpa pyMax (ringkasan `.USEDBY` b1274 lebar layar 200 px b1284) → lebar = batas nama ringkasan 200 byte (preseden riratelife `BatasNama`); data DEV "RI COMM RETRO" |
| `CONTRACT` | angka bulat | NUMBER(5) | `.CONTRACT` b1904, **pxNumber** b1907, wajib b1917/b1923, placeholder `0` b1927, `ChangeDotToPoint_DT` b1953. **[penyimpangan sadar]** XML tanpa presisi → K6 bilangan kecil NUMBER(5); Go: bulat 0-99999 (preseden riratelife: CONTRACT bulat) |
| `YEAR` | angka bulat | NUMBER(5) | `.YEAR` b2185, pxTextInput b2188, wajib b2196/b2201, placeholder `0` b2205, **pyMax 4** b2206, rata kanan b2207. **[penyimpangan sadar]** kontrol teks tetapi placeholder angka + rata kanan + 4 angka → K6 tahun NUMBER(5) (NUMBER(4) bukan bentuk sah `TestNolNumberTanpaPresisi`); Go: bulat paling banyak 4 angka |
| `COMM` | angka desimal | NUMBER(38,8) | `.COMM` b2392, pxTextInput b2395, wajib b2406/b2412; grid **pxNumber** b9101. **[penyimpangan sadar]** XML tanpa desimal/presisi → K6 persen/rate NUMBER(38,8); Go: desimal tak bertanda, koma atau titik, ≤ 30 angka bulat dan ≤ 8 desimal, tidak pernah dibulatkan |

Angka ditulis tanpa bergantung NLS sesi (`TO_NUMBER(:n)` atas teks angka; COMM `TO_NUMBER(:koef) / POWER(10, :skala)`)
dan dibaca `TO_CHAR(..., 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')` → API memakai titik desimal kanonik (`0.5`); layar
menampilkan koma desimal.

## M_RICOMM_LIFE_SUMMARY

Tabel warisan (dinyatakan di `MODUL.md`, bukan dibuat migrasi), TIDAK diubah strukturnya (butir 6).

| Kolom | Tipe | Null | Dipakai | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | VARCHAR2(10) | ? | "ID" | site + LPAD(`M_RICOMM_LIFE_SUMMARY_SEQ`, 6) |
| `JSONDATA` | CLOB | ? | JSON | sisip `JSON_OBJECT(... ABSENT ON NULL)`; ubah = baca `FOR UPDATE`, ganti kunci di Go, tulis utuh (nol `JSON_MERGEPATCH`) |

Kunci JSON yang ditulis: `USEDBY` (wajib, tidak kembar tanpa beda huruf), `OPERATORID` (akun login), `MODIFIEDDATE`
(Pega `yyyyMMdd'T'HHmmss.SSS 'GMT'`), dan pada sisip `pxObjClass` = `ASM-FW-GISFW-Int-RICOMM_LIFE_SUMMARY`. Bukti
`pxObjClass`: satu-satunya baris DEV (ID 1000003, dicek WO 06-10-2026) memuat
`{"MODIFIEDDATE":"20181205T073755.559 GMT","OPERATORID":"…","pxObjClass":"ASM-FW-GISFW-Int-RICOMM_LIFE_SUMMARY","USEDBY":"RI COMM RETRO"}`
dan kelas section `InboxSummaryRIComm` b84. Edit tidak menyentuh kunci lain.

## RICOMM_LIFE_SUMMARY

View warisan (dinyatakan di `MODUL.md`), dibaca saja: `SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.MODIFIEDDATE,
a.JSONDATA.OPERATORID FROM POOLDATA.M_RICOMM_LIFE_SUMMARY a`.

| Kolom | Tipe | Null | Dipakai | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | VARCHAR2(10) | ? | "ID" b1091 | `M_RICOMM_LIFE_SUMMARY.ID` |
| `USEDBY` | VARCHAR2(4000) | ? | "R/I COMM NAME" b7220 | `JSONDATA.USEDBY` |
| `MODIFIEDDATE` | VARCHAR2(4000) | ? | "MODIFY DATE" b7521 / `.MODIFIEDDATE` b8526 | `JSONDATA.MODIFIEDDATE` |
| `OPERATORID` | VARCHAR2(4000) | ? | "MODIFY OPERATOR" b7377 / `.OPERATORID` b8348 | `JSONDATA.OPERATORID` |

## M_RICOMM_LIFE

Tabel warisan (dinyatakan di `MODUL.md`), 0 baris DEV; TIDAK disentuh (butir 3) - sumber alat `pindahflat` saja.

| Kolom | Tipe | Null | Dipakai | Sumber |
| --- | --- | --- | --- | --- |
| `ID` | VARCHAR2(10) | ? | ID rincian | prosedur `PEGA_M_RICOMM_LIFE` b21 |
| `JSONDATA` | CLOB | ? | kunci `IDUSEDBY`, `USEDBY`, `CONTRACT`, `YEAR`, `COMM` | prosedur `PEGA_M_RICOMM_LIFE` b29 |
