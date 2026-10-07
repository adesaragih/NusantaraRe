# Langkah work owner — ringkasan R/I Rate Life menjadi SATU tabel `M_RATE_LIFE_SUMMARY` (RALAT R6)

> Untuk work owner / DBA. **Dijalankan WO, bukan executor** — tidak satu pun perintah di bawah pernah dijalankan saat
> menulis berkas ini. Keputusan work owner 07-10-2026 (`modul/riratelife/MODUL.md`, RALAT R6). Skema DEV: POOLDATA.
> Menggantikan `LANGKAH-WO-RIRATELIFE-FLAT.md` (dihapus; riwayatnya di git). Tidak ada angka baris tetap di sini.

Keadaan awal yang diandaikan: migrasi 926 sudah jalan; `POOLDATA.RATE_LIFE_SUMMARY` = TABLE (sumber kebenaran,
aplikasi menulis ke sana); `POOLDATA.M_RATE_LIFE_SUMMARY` = ID + JSONDATA. Pega tidak lagi menyimpan R/I Rate (WO).

Urutan WAJIB (a) → (d). Kerjakan dalam satu jendela pemeliharaan: hentikan backend :8080 sebelum (b), supaya tidak ada
tulisan aplikasi ke tabel flat di tengah pemindahan.

## (a) Cadangan + angka acuan — baca-saja, SEBELUM apa pun dibuang

Satu-satunya salinan FLAG lama (dan ringkasan yang sudah dihapus aplikasi) ada di `M_RATE_LIFE_SUMMARY.JSONDATA`.
SQL*Plus 12.2+ (atau SQLcl), sebagai pemilik POOLDATA, dari folder tempat cadangan disimpan:

```sql
SET MARKUP CSV ON QUOTE ON
SET FEEDBACK OFF TERMOUT OFF PAGESIZE 0 LONG 2000000 LONGCHUNKSIZE 32767 LINESIZE 32767 TRIMSPOOL ON
SPOOL M_RATE_LIFE_SUMMARY_JSONDATA_20261007.csv
SELECT ID, JSONDATA FROM POOLDATA.M_RATE_LIFE_SUMMARY ORDER BY ID;
SPOOL OFF
SPOOL RATE_LIFE_SUMMARY_FLAT_20261007.csv
SELECT ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID FROM POOLDATA.RATE_LIFE_SUMMARY ORDER BY ID;
SPOOL OFF
SET TERMOUT ON FEEDBACK ON MARKUP CSV OFF
```

(SQL Developer: klik kanan hasil `SELECT ID, JSONDATA FROM POOLDATA.M_RATE_LIFE_SUMMARY ORDER BY ID` → Export → format
csv; ulangi untuk `RATE_LIFE_SUMMARY`.) Periksa kedua berkas terbuka dan cacah barisnya = kueri di bawah.

**Catat angka acuan** (dipakai di (c)):

```sql
SELECT (SELECT COUNT(*) FROM POOLDATA.RATE_LIFE_SUMMARY)   AS FLAT_SEBELUM,
       (SELECT COUNT(*) FROM POOLDATA.M_RATE_LIFE_SUMMARY) AS JSON_SEBELUM FROM DUAL;
SELECT OBJECT_NAME, OBJECT_TYPE, STATUS FROM SYS.ALL_OBJECTS
 WHERE OWNER = 'POOLDATA' AND OBJECT_NAME IN ('PEGA_M_RATE_LIFE_SUMMARY', 'PEGA_M_PLAN_LIFE_SUMMARY');
```

## (b) Migrasi 927 + 928 dari cabang `modul/riratelife/implementasi` — PowerShell, folder akar repo

```powershell
git switch modul/riratelife/implementasi
. .\muat-env.ps1                 # dot-source; ORACLE_SCHEMA=POOLDATA, IS_PEGA_PROD=false (PANDUAN-MENJALANKAN.txt bab 3)
go run ./cmd/api -migrate        # = target Makefile `migrate`
```

Versi cmd.exe: `call muat-env.cmd` lalu `go run ./cmd/api -migrate`. Log yang benar: `dijalankan: 927_m_rate_life_summary_kolom`
lalu `dijalankan: 928_m_rate_life_summary_satu_tabel`. Tidak ada berkas DBA terpisah: seluruh langkah diterima penjaga
migrasi inti dan pelari (927 = ALTER … ADD tunggal; 928 = UPDATE, DELETE, blok berpelindung katalog pembuang JSONDATA,
INSERT `NOT EXISTS`, CREATE INDEX, DROP TABLE flat terakhir).

## (c) Verifikasi — SQL, baca-saja

```sql
-- 1. M_RATE_LIFE_SUMMARY tepat 5 kolom: ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID (nol JSONDATA)
SELECT COLUMN_NAME, DATA_TYPE, DATA_LENGTH, NULLABLE FROM SYS.ALL_TAB_COLUMNS
 WHERE OWNER = 'POOLDATA' AND TABLE_NAME = 'M_RATE_LIFE_SUMMARY' ORDER BY COLUMN_ID;

-- 2. COUNT = FLAT_SEBELUM dari (a)
SELECT COUNT(*) FROM POOLDATA.M_RATE_LIFE_SUMMARY;

-- 3. RATE_LIFE_SUMMARY tidak ada lagi (nol baris); constraint IS JSON hilang (nol baris); indeks nama baru ada (satu baris)
SELECT OBJECT_NAME, OBJECT_TYPE FROM SYS.ALL_OBJECTS WHERE OWNER = 'POOLDATA' AND OBJECT_NAME = 'RATE_LIFE_SUMMARY';
SELECT CONSTRAINT_NAME FROM SYS.ALL_CONSTRAINTS WHERE OWNER = 'POOLDATA' AND CONSTRAINT_NAME = 'ENSURE_M_RATE_LIFE_SUMMARY_JSON';
SELECT INDEX_NAME FROM SYS.ALL_INDEXES WHERE OWNER = 'POOLDATA' AND INDEX_NAME = 'IX_M_RATE_LIFE_SUMMARY_NAMA';

-- 4. T_MIGRASI (kolom NAMA, DIJALANKAN_PADA) - dua baris
SELECT NAMA, DIJALANKAN_PADA FROM POOLDATA.T_MIGRASI
 WHERE NAMA IN ('927_m_rate_life_summary_kolom', '928_m_rate_life_summary_satu_tabel') ORDER BY NAMA;

-- 5. Prosedur Pega yang bergantung - CATAT status (INVALID diharapkan, diterima WO; jangan dikompilasi ulang)
SELECT OBJECT_NAME, OBJECT_TYPE, STATUS FROM SYS.ALL_OBJECTS
 WHERE OWNER = 'POOLDATA' AND OBJECT_NAME IN ('PEGA_M_RATE_LIFE_SUMMARY', 'PEGA_M_PLAN_LIFE_SUMMARY');
```

Isi dapat dibandingkan dengan `RATE_LIFE_SUMMARY_FLAT_20261007.csv` dari (a) (ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID
harus sama baris demi baris).

## (d) Restart backend dan periksa ketiga layar

```powershell
go run ./cmd/api                 # = target Makefile `run-api`, jendela yang sudah memuat env
```

- R/I Rate Life: jumlah ringkasan = COUNT (c)2; Add / Edit / Delete / Rate Detail / Upload berjalan.
- Product Name Life: `Choose R/I Rate` menampilkan daftar yang sama.
- Contract Retro Life: autocomplete `R/I RATE` menampilkan daftar yang sama.

## Pemulihan — 927 / 928 gagal di tengah

Pelari menjalankan pernyataan TANPA transaksi; DDL langsung tetap.

| Keadaan (dari log `-migrate` dan `T_MIGRASI`) | Tindakan |
| --- | --- |
| 927 gagal | 927 hanya satu pernyataan: tidak ada yang berubah. Perbaiki sebabnya (pesan ORA), ulangi (b) |
| 927 tercatat, 928 gagal di pernyataan mana pun SEBELUM `DROP TABLE` | ulangi (b) apa adanya: UPDATE/DELETE menurut flat menghasilkan isi sama, blok pembuang JSONDATA melewati kolom yang sudah tidak ada, INSERT `NOT EXISTS` tidak menggandakan, CREATE INDEX yang sudah ada ditoleransi (ORA-00955) |
| 928 gagal SESUDAH `DROP TABLE` berhasil (mis. koneksi putus sebelum pencatatan; `RATE_LIFE_SUMMARY` sudah tidak ada, `T_MIGRASI` tanpa 928) | JANGAN ulangi (b) (UPDATE akan mati di ORA-00942). Pastikan (c)1-3 sudah benar, lalu catat langkahnya: `INSERT INTO POOLDATA.T_MIGRASI (NAMA, DIJALANKAN_PADA) VALUES ('928_m_rate_life_summary_satu_tabel', SYSDATE); COMMIT;` lalu lanjut (c)4 |

## Jalur mundur

Hanya bila perlu kembali ke keadaan sesudah 926 (tabel flat sebagai penyimpan aplikasi). Pernyataan
`inti/backend/migrations/928_m_rate_life_summary_satu_tabel_down.sql` lalu `927_m_rate_life_summary_kolom_down.sql`
dijalankan DBA berurutan (ganti `{skema}` dengan POOLDATA), lalu `DELETE FROM POOLDATA.T_MIGRASI WHERE NAMA IN
('928_m_rate_life_summary_satu_tabel', '927_m_rate_life_summary_kolom'); COMMIT;` dan backend dari commit sebelum
RALAT R6. JSONDATA dibangun ulang dari kolom (kunci USEDBY, TYPE, MODIFIEDDATE, OPERATORID); **FLAG lama dan ringkasan
yang dibuang 928 hanya dapat dikembalikan dari cadangan CSV (a)**. Prosedur Pega tetap perlu dikompilasi ulang DBA.
⛔ Jangan memakai `-migrate-down` di DEV (dipagari skema uji, membongkar SELURUH migrasi).
