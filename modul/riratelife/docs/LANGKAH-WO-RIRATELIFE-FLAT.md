# Langkah work owner — ringkasan R/I Rate Life menjadi tabel flat (RALAT R4)

> Untuk work owner / DBA. **Dijalankan WO, bukan executor** — tidak satu pun perintah di bawah pernah dijalankan saat
> menulis berkas ini. Keputusan work owner 06-10-2026 K-F1/K-F2 dan cutover delta (`modul/riratelife/MODUL.md`).
> Skema DEV: POOLDATA. `M_RATE_LIFE_SUMMARY` MASIH BERUBAH: tidak ada angka baris tetap di sini - setiap pemeriksaan
> membandingkan kedua tabel SAAT itu.

Urutan WAJIB (a) → (f). Langkah (a) harus selesai SEBELUM (b): selama view `RATE_LIFE_SUMMARY` ada, pra-terbang pelari
migrasi menganggap objek bernama sama "sudah ada" dan melewati CREATE TABLE 926; indeks 926 lalu gagal atas view
(ORA-01702) dan `-migrate` berhenti (bab Pemulihan).

⚠️ Sebaiknya **bekukan penulisan Pega ke R/I Rate** (layar / upload Pega R/I Rate Summary) sejak (a). Bila belum dapat,
langkah (e) menarik delta terakhir. Antara (a) dan (f) layar R/I Rate Life, autocomplete `R/I RATE` (Contract Retro
Life), dan `Choose R/I Rate` (Product Name Life) menjawab galat karena objeknya belum ada / backend belum dimuat ulang —
kerjakan (a)–(f) dalam satu jendela pemeliharaan.

## (a) Lepas view `POOLDATA.RATE_LIFE_SUMMARY` — SQL*Plus / SQL Developer, pemilik POOLDATA atau DBA

Jalankan `modul/riratelife/docs/DBA-LEPAS-VIEW-RATE_LIFE_SUMMARY.sql` bagian demi bagian:

1. Bagian 1 — definisi view harus sama dengan `SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.MODIFIEDDATE,
   a.JSONDATA.OPERATORID, a.JSONDATA.FLAG FROM M_RATE_LIFE_SUMMARY a`.
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

## (c) Verifikasi bentuk — SQL, baca-saja

```sql
-- RATE_LIFE_SUMMARY kini TABLE (+ INDEX IX_RATE_LIFE_SUMMARY_NAMA)
SELECT OBJECT_NAME, OBJECT_TYPE FROM SYS.ALL_OBJECTS
 WHERE OWNER = 'POOLDATA' AND OBJECT_NAME IN ('RATE_LIFE_SUMMARY', 'IX_RATE_LIFE_SUMMARY_NAMA');

-- T_MIGRASI (kolom NAMA, DIJALANKAN_PADA - inti/backend/migrasi/migrasi.go) - harus satu baris
SELECT NAMA, DIJALANKAN_PADA FROM POOLDATA.T_MIGRASI WHERE NAMA = '926_rate_life_summary_flat';

-- Kolom - harus ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID, FLAG
SELECT COLUMN_NAME, DATA_TYPE, DATA_LENGTH, NULLABLE FROM SYS.ALL_TAB_COLUMNS
 WHERE OWNER = 'POOLDATA' AND TABLE_NAME = 'RATE_LIFE_SUMMARY' ORDER BY COLUMN_ID;
```

## (d) Pemindahan penuh — PowerShell, jendela yang sama (env sudah dimuat)

```powershell
go run ./modul/riratelife/backend/alat/pindahflat              # uji kering: SELECT saja
```

Lolos bila: `baru (disisip)` = cacah sumber yang disebut baris `cacah saat dijalankan`, `tidak muat: 0`, `gagal: 0`.
Tabel `panjang maksimum sumber (byte) / lebar kolom` menyebut panjang terbesar tiap kolom; bila ada `tidak muat`:
**BERHENTI** dan laporkan ID + kolomnya (nol pemotongan).

```powershell
go run ./modul/riratelife/backend/alat/pindahflat -jalankan    # satu koneksi, satu transaksi
```

Lolos bila baris terakhir `ditulis: true`. **Catat baris `batas delta berikutnya (-sejak): …`** — dipakai di (e).
Lalu jalankan **Verifikasi isi** (bab di bawah).

## (e) Jalankan ulang alat pindah (DELTA) sesaat sebelum restart backend

Menarik ringkasan yang ditulis / diubah Pega sejak (d). Ganti `<BATAS>` dengan nilai yang dicatat di (d) (atau dari
putaran delta terakhir):

```powershell
go run ./modul/riratelife/backend/alat/pindahflat -sejak="<BATAS>"              # uji kering delta
go run ./modul/riratelife/backend/alat/pindahflat -jalankan -sejak="<BATAS>"    # tulis delta
```

Laporan menyebut `baru`, `berubah` (Pega lebih baru, diperbarui), `sama`, `konflik` (TIDAK ditimpa — ID disebut; periksa
satu per satu bersama pemilik datanya), dan `dilewati` (ada di JSON, tidak di flat, tidak lebih baru dari batas =
dihapus aplikasi). Lalu ulangi **Verifikasi isi**. Aturan lengkap: MODUL.md bab "Cutover delta".

## Verifikasi isi — SQL, baca-saja (sesudah (d) dan sesudah (e))

```sql
-- 1. Cacah harus sama (bila belum ada tulisan aplikasi di tabel flat)
SELECT (SELECT COUNT(*) FROM POOLDATA.M_RATE_LIFE_SUMMARY) AS SUMBER,
       (SELECT COUNT(*) FROM POOLDATA.RATE_LIFE_SUMMARY)   AS FLAT FROM DUAL;

-- 2. MINUS dua arah atas keenam kolom - KEDUANYA harus nol baris. Sisi JSON = notasi titik definisi view lama.
SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID, a.JSONDATA.FLAG
  FROM POOLDATA.M_RATE_LIFE_SUMMARY a
MINUS
SELECT ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID, FLAG FROM POOLDATA.RATE_LIFE_SUMMARY;

SELECT ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID, FLAG FROM POOLDATA.RATE_LIFE_SUMMARY
MINUS
SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID, a.JSONDATA.FLAG
  FROM POOLDATA.M_RATE_LIFE_SUMMARY a;
```

Baris yang muncul hanya boleh ID yang disebut laporan alat sebagai `konflik` atau `dilewati` (atau ringkasan baru
aplikasi bila backend sudah berjalan). Selain itu: BERHENTI dan laporkan.

## (f) Restart backend

Hentikan proses backend :8080 (Ctrl+C di jendelanya), lalu di jendela yang sudah memuat env:

```powershell
go run ./cmd/api                 # = target Makefile `run-api`
```

Periksa: jumlah ringkasan di menu R/I Rate Life = `SELECT COUNT(*) FROM POOLDATA.RATE_LIFE_SUMMARY`; Add / Edit /
Delete / Rate Detail / Upload berjalan; Contract Retro Life (`R/I RATE`) dan Product Name Life (`Choose R/I Rate`)
menampilkan daftar yang sama. Bila Pega belum dibekukan: ulangi (e) berkala dengan `-sejak` terbaru.

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
4. Lanjut ke (c) - (f).

**Jalur mundur** (hanya SEBELUM aplikasi menulis ringkasan baru - tulisan aplikasi tidak ada di JSON dan hilang):
pernyataan `inti/backend/migrations/926_rate_life_summary_flat_down.sql` dijalankan DBA (DROP TABLE lalu CREATE VIEW
asli, ganti `{skema}` dengan POOLDATA), lalu `DELETE FROM POOLDATA.T_MIGRASI WHERE NAMA = '926_rate_life_summary_flat';
COMMIT;` dan backend dari commit sebelum RALAT R4. ⛔ Jangan memakai `-migrate-down` di DEV (dipagari skema uji,
membongkar SELURUH migrasi).
