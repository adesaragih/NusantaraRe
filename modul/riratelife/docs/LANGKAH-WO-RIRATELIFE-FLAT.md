# Langkah work owner — ringkasan R/I Rate Life menjadi tabel flat (RALAT R4)

> Untuk work owner / DBA. **Dijalankan WO, bukan executor** — tidak satu pun perintah di bawah pernah dijalankan saat
> menulis berkas ini. Keputusan work owner 06-10-2026 K-F1/K-F2 (`modul/riratelife/MODUL.md`). Skema DEV: POOLDATA.

Urutan WAJIB (a) → (e). Langkah (a) harus selesai SEBELUM (b): selama view `RATE_LIFE_SUMMARY` ada, pra-terbang pelari
migrasi menganggap objek bernama sama "sudah ada" dan melewati CREATE TABLE 926; indeks 926 lalu gagal atas view
(ORA-01702) dan `-migrate` berhenti (bab Pemulihan).

⚠️ Sebelum (a): **bekukan penulisan Pega ke R/I Rate** (layar Pega R/I Rate Summary / upload Pega). Ringkasan yang
ditulis Pega sesudah langkah (d) tidak masuk tabel flat sampai alat pindah dijalankan ulang (risiko cutover, MODUL.md).
Antara (a) dan (b) layar R/I Rate Life, autocomplete `R/I RATE` (Contract Retro Life), dan `Choose R/I Rate` (Product
Name Life) menjawab galat karena objeknya belum ada — kerjakan (a)–(e) dalam satu jendela pemeliharaan.

## (a) Lepas view `POOLDATA.RATE_LIFE_SUMMARY` — SQL*Plus / SQL Developer, pemilik POOLDATA atau DBA

Jalankan `modul/riratelife/docs/DBA-LEPAS-VIEW-RATE_LIFE_SUMMARY.sql` bagian demi bagian:

1. Bagian 1 — definisi view harus sama dengan `SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.MODIFIEDDATE,
   a.JSONDATA.OPERATORID, a.JSONDATA.FLAG FROM M_RATE_LIFE_SUMMARY a`; catat `COUNT(*)` M_RATE_LIFE_SUMMARY (339 di DEV).
2. Bagian 2 dan 3 — **harus nol baris** (objek yang bergantung, sinonim, grant). Bila ada baris: **BERHENTI**, laporkan.
3. Bagian 4 — `DROP VIEW POOLDATA.RATE_LIFE_SUMMARY;`
4. Bagian 5 — harus nol baris.

## (b) Migrasi dari cabang `modul/riratelife/implementasi` — PowerShell, folder akar repo

```powershell
git switch modul/riratelife/implementasi
. .\muat-env.ps1                 # dot-source; ORACLE_SCHEMA=POOLDATA, IS_PEGA_PROD=false (PANDUAN-MENJALANKAN.txt bab 3)
go run ./cmd/api -migrate        # = target Makefile `migrate`
```

Versi cmd.exe: `call muat-env.cmd` lalu `go run ./cmd/api -migrate`. Log yang benar: `dijalankan: 926_rate_life_summary_flat`,
tanpa `dilewati, objeknya sudah ada: RATE_LIFE_SUMMARY`.

## (c) Verifikasi — SQL, baca-saja

```sql
-- RATE_LIFE_SUMMARY kini TABLE - harus satu baris OBJECT_TYPE = 'TABLE' (+ satu INDEX IX_RATE_LIFE_SUMMARY_NAMA)
SELECT OBJECT_NAME, OBJECT_TYPE FROM SYS.ALL_OBJECTS
 WHERE OWNER = 'POOLDATA' AND OBJECT_NAME IN ('RATE_LIFE_SUMMARY', 'IX_RATE_LIFE_SUMMARY_NAMA');

-- T_MIGRASI (kolom NAMA, DIJALANKAN_PADA - inti/backend/migrasi/migrasi.go) - harus satu baris
SELECT NAMA, DIJALANKAN_PADA FROM POOLDATA.T_MIGRASI WHERE NAMA = '926_rate_life_summary_flat';

-- Kolom - harus ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID, FLAG
SELECT COLUMN_NAME, DATA_TYPE, DATA_LENGTH, NULLABLE FROM SYS.ALL_TAB_COLUMNS
 WHERE OWNER = 'POOLDATA' AND TABLE_NAME = 'RATE_LIFE_SUMMARY' ORDER BY COLUMN_ID;

-- Tabel flat masih kosong; sumber tetap utuh (339 di DEV)
SELECT (SELECT COUNT(*) FROM POOLDATA.RATE_LIFE_SUMMARY) AS FLAT, (SELECT COUNT(*) FROM POOLDATA.M_RATE_LIFE_SUMMARY) AS SUMBER FROM DUAL;
```

## (d) Pindahkan 339 ringkasan — PowerShell, jendela yang sama (env sudah dimuat)

```powershell
go run ./modul/riratelife/backend/alat/pindahflat              # uji kering: SELECT saja
```

Lolos bila: `sumber M_RATE_LIFE_SUMMARY: 339` (atau cacah saat itu), `baris yang (akan) ditulis: 339`, `tidak muat: 0`,
`gagal: 0`, `tabel flat berbeda: 0`. Tabel `panjang maksimum sumber (byte) / lebar kolom` menyebut panjang terbesar tiap
kolom; bila ada `tidak muat`: **BERHENTI** dan laporkan ID + kolomnya (lebar kolom perlu keputusan WO; nol pemotongan).

```powershell
go run ./modul/riratelife/backend/alat/pindahflat -jalankan    # satu koneksi, satu transaksi; aman diulang
```

Lolos bila baris terakhir `ditulis: true`. Verifikasi:

```sql
-- FLAT harus = SUMBER (339 di DEV); baris yang hanya ada di salah satu - harus nol baris
SELECT (SELECT COUNT(*) FROM POOLDATA.RATE_LIFE_SUMMARY) AS FLAT, (SELECT COUNT(*) FROM POOLDATA.M_RATE_LIFE_SUMMARY) AS SUMBER FROM DUAL;
SELECT ID FROM POOLDATA.M_RATE_LIFE_SUMMARY MINUS SELECT ID FROM POOLDATA.RATE_LIFE_SUMMARY;
```

Mengulang `-jalankan` sesudahnya menulis `sudah ada dan sama (dilewati): 339`, `baris yang (akan) ditulis: 0`.

## (e) Restart backend

Hentikan proses backend :8080 (Ctrl+C di jendelanya), lalu di jendela yang sudah memuat env:

```powershell
go run ./cmd/api                 # = target Makefile `run-api`
```

Periksa: menu R/I Rate Life menampilkan 339 ringkasan; Add / Edit / Delete / Rate Detail / Upload berjalan; Contract
Retro Life (autocomplete `R/I RATE`) dan Product Name Life (`Choose R/I Rate`) menampilkan daftar yang sama.

## Pemulihan — (b) terjalankan sebelum (a)

Tanda: log `dilewati, objeknya sudah ada: RATE_LIFE_SUMMARY` lalu `migrasi: … 926_rate_life_summary_flat … ORA-01702`
(-migrate berhenti). Tidak ada data yang rusak: view dan M_RATE_LIFE_SUMMARY tidak berubah. Langkah eksak:

```sql
-- 1. Periksa catatan; bila 926 TERNYATA tercatat, hapus catatannya (bila nol baris, lewati DELETE)
SELECT NAMA FROM POOLDATA.T_MIGRASI WHERE NAMA = '926_rate_life_summary_flat';
DELETE FROM POOLDATA.T_MIGRASI WHERE NAMA = '926_rate_life_summary_flat';
COMMIT;
-- Indeks tidak boleh tertinggal - harus nol baris
SELECT INDEX_NAME FROM SYS.ALL_INDEXES WHERE OWNER = 'POOLDATA' AND INDEX_NAME = 'IX_RATE_LIFE_SUMMARY_NAMA';
```

2. Jalankan (a) seluruhnya (berkas DBA bagian 1-5, termasuk `DROP VIEW POOLDATA.RATE_LIFE_SUMMARY;`).
3. Jalankan ulang (b): `. .\muat-env.ps1` lalu `go run ./cmd/api -migrate` — `926_rate_life_summary_flat` dijalankan.
4. Lanjut ke (c), (d), (e).

**Jalur mundur** (hanya SEBELUM aplikasi menulis ringkasan baru - tulisan aplikasi tidak ada di JSON dan hilang):
pernyataan `inti/backend/migrations/926_rate_life_summary_flat_down.sql` dijalankan DBA (DROP TABLE lalu CREATE VIEW
asli, ganti `{skema}` dengan POOLDATA), lalu `DELETE FROM POOLDATA.T_MIGRASI WHERE NAMA = '926_rate_life_summary_flat';
COMMIT;` dan backend dari commit sebelum RALAT R4. ⛔ Jangan memakai `-migrate-down` di DEV (dipagari skema uji,
membongkar SELURUH migrasi).
