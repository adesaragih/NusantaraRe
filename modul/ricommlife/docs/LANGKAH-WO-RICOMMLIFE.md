# Langkah work owner — menyalakan R/I Comm Life di DEV

> Untuk work owner / DBA. **Dijalankan WO, bukan executor** — tidak satu pun perintah di bawah pernah dijalankan saat
> menulis berkas ini. Keadaan DEV yang dicek WO 06-10-2026 (baca-saja): backend :8080 sudah memuat `ricommlife`;
> BELUM ada baris `M_NAV_MENU` ricommlife, catatan `T_MIGRASI` 924/925, maupun hak `M_LOGIN_GO_MENU`;
> `POOLDATA.RICOMM_LIFE` masih VIEW.

Urutan WAJIB persis (a) → (e). Langkah (a) harus selesai SEBELUM (b): selama view `RICOMM_LIFE` ada, pra-terbang pelari
migrasi menganggap objek bernama sama "sudah ada" dan **melewati** CREATE TABLE 924 (lihat bab Pemulihan).

## (a) Lepas view `POOLDATA.RICOMM_LIFE` — SQL*Plus / SQL Developer, pemilik POOLDATA atau DBA

Jalankan berkas `modul/ricommlife/docs/DBA-LEPAS-VIEW-RICOMM_LIFE.sql` bagian demi bagian:

1. Bagian 1 — catat definisi view; harus sama dengan
   `SELECT a.ID, a.JSONDATA.IDUSEDBY, a.JSONDATA.USEDBY, a.JSONDATA.CONTRACT, a.JSONDATA.YEAR, a.JSONDATA.COMM FROM M_RICOMM_LIFE a`.
2. Bagian 2 dan 3 — **harus nol baris** (dependensi, sinonim, grant). Bila ada baris: BERHENTI, laporkan.
3. Bagian 4 — `DROP VIEW POOLDATA.RICOMM_LIFE;`
4. Bagian 5 — harus nol baris.

## (b) Migrasi dari cabang `modul/ricommlife/implementasi` — PowerShell, folder akar repo

```powershell
git switch modul/ricommlife/implementasi
. .\muat-env.ps1                 # dot-source; ORACLE_SCHEMA=POOLDATA, IS_PEGA_PROD=false (PANDUAN-MENJALANKAN.txt bab 3)
go run ./cmd/api -migrate        # = target Makefile `migrate`
```

Versi cmd.exe: `call muat-env.cmd` lalu `go run ./cmd/api -migrate`. Log yang benar menyebut `924_ricomm_life` dan
`925_m_nav_menu_ricommlife` sebagai langkah yang dijalankan (`dijalankan: …`), dan **tidak** memuat `dilewati, objeknya sudah ada: RICOMM_LIFE`.

## (c) Verifikasi — SQL, baca-saja

```sql
-- RICOMM_LIFE kini TABLE (bukan VIEW) - harus tepat satu baris, OBJECT_TYPE = 'TABLE'
SELECT OBJECT_NAME, OBJECT_TYPE FROM SYS.ALL_OBJECTS WHERE OWNER = 'POOLDATA' AND OBJECT_NAME = 'RICOMM_LIFE';

-- T_MIGRASI (kolom NAMA, DIJALANKAN_PADA - inti/backend/migrasi/migrasi.go) - harus dua baris
SELECT NAMA, DIJALANKAN_PADA FROM POOLDATA.T_MIGRASI
 WHERE NAMA IN ('924_ricomm_life', '925_m_nav_menu_ricommlife') ORDER BY NAMA;

-- Baris menu - harus satu baris: LABEL 'R/I Comm Life', GROUPMENU 'MASTER TREATY', URUTAN 10, DIMIGRASI '1'
SELECT KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI FROM POOLDATA.M_NAV_MENU WHERE KODE = 'ricommlife';

-- Kolom tabel flat - harus ID, IDUSEDBY, USEDBY, CONTRACT, YEAR, COMM
SELECT COLUMN_NAME, DATA_TYPE, DATA_PRECISION, DATA_SCALE, CHAR_LENGTH FROM SYS.ALL_TAB_COLUMNS
 WHERE OWNER = 'POOLDATA' AND TABLE_NAME = 'RICOMM_LIFE' ORDER BY COLUMN_ID;
```

Opsional (M_RICOMM_LIFE 0 baris di DEV): `go run ./modul/ricommlife/backend/alat/pindahflat` (uji kering) harus
menulis `gagal: 0`.

## (d) Hak menu superadmin — SQL

Superadmin = pemegang menu `kelolauser` (migrasi 914 b9-b10). `HAK` CHECK `('LIHAT','PENUH')` (914); PK
`(LOGIN_ID, MENU_KODE)` (903) - `NOT EXISTS` membuatnya aman diulang.

```sql
INSERT INTO POOLDATA.M_LOGIN_GO_MENU (LOGIN_ID, MENU_KODE, HAK)
SELECT m.LOGIN_ID, 'ricommlife', 'PENUH'
  FROM POOLDATA.M_LOGIN_GO_MENU m
 WHERE m.MENU_KODE = 'kelolauser'
   AND NOT EXISTS (SELECT 1 FROM POOLDATA.M_LOGIN_GO_MENU x
                    WHERE x.LOGIN_ID = m.LOGIN_ID AND x.MENU_KODE = 'ricommlife');
COMMIT;

-- Periksa: cacahnya = cacah pemegang kelolauser
SELECT COUNT(*) FROM POOLDATA.M_LOGIN_GO_MENU WHERE MENU_KODE = 'ricommlife';
```

Akun lain: centang menu "R/I Comm Life" di Kelola User.

## (e) Restart backend

Hentikan proses backend :8080 (Ctrl+C di jendelanya), lalu di jendela yang sudah memuat env:

```powershell
go run ./cmd/api                 # = target Makefile `run-api`
```

Masuk sebagai superadmin → sidebar MASTER TREATY memuat "R/I Comm Life" di urutan 10 → grid R/I COMM SUMMARY tampil.

## Pemulihan — (b) terjalankan sebelum (a)

Yang terjadi bila view masih ada: pra-terbang mencocokkan kolom view dengan DDL 924, CREATE TABLE dijawab ORA-00955 dan
dilewati (log `dilewati, objeknya sudah ada: RICOMM_LIFE`). Pernyataan berikutnya, `CREATE INDEX … ON RICOMM_LIFE`,
gagal atas view (ORA-01702) sehingga **-migrate berhenti** (`migrasi: … 924_ricomm_life …`) dan 924 normalnya TIDAK
tercatat; 925 belum berjalan. Tidak ada data yang rusak (view tidak berubah, M_RICOMM_LIFE tidak disentuh).
Langkah eksak:

```sql
-- 1. Periksa catatan; bila 924 (atau 925) TERNYATA tercatat, hapus catatannya (bila nol baris, lewati DELETE)
SELECT NAMA FROM POOLDATA.T_MIGRASI WHERE NAMA IN ('924_ricomm_life', '925_m_nav_menu_ricommlife');
DELETE FROM POOLDATA.T_MIGRASI WHERE NAMA = '924_ricomm_life';
COMMIT;
-- Periksa juga: indeks tidak boleh tertinggal (harus nol baris)
SELECT INDEX_NAME FROM SYS.ALL_INDEXES WHERE OWNER = 'POOLDATA' AND INDEX_NAME = 'IX_RICOMM_LIFE_IDUSEDBY';
```

2. Jalankan langkah (a) seluruhnya (berkas DBA bagian 1-5, termasuk `DROP VIEW POOLDATA.RICOMM_LIFE;`).
3. Jalankan ulang (b): `. .\muat-env.ps1` lalu `go run ./cmd/api -migrate` — `924_ricomm_life` dan (bila belum
   tercatat) `925_m_nav_menu_ricommlife` dijalankan; log tanpa `dilewati, objeknya sudah ada`.
4. Jalankan (c): `RICOMM_LIFE` = TABLE, `T_MIGRASI` memuat 924 dan 925, baris menu benar.
5. Lanjut ke (d) dan (e).

⛔ Jangan memakai `-migrate-down` di DEV: ia dipagari skema uji dan membongkar SELURUH migrasi yang tercatat.
