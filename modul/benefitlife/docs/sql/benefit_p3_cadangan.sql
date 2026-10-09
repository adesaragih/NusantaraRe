-- LANGKAH-WO-BENEFITLIFE.md (a) P3 - BACA-SAJA. CADANGAN WAJIB SEBELUM -migrate 942-945: ID + JSONDATA UTUH tabel
-- M_BENEFIT_LIFE dan isi view BENEFIT_LIFE, ke folder kerja SQL*Plus yang dibuat lebih dulu:
-- cadangan-benefitlife-<tanggal>\ (mis. cadangan-benefitlife-20261008). Tidak ada cadangan = TIDAK BOLEH -migrate.
-- Dijalankan sebagai BERKAS (@...) supaya TERMOUT OFF berlaku. SQL*Plus 12.2+ atau SQLcl (MARKUP CSV).
-- ⛔ JSONDATA memuat kunci pxCreateOperator / pxCreateOpName (identitas orang): simpan di tempat cadangan DBA, jangan
-- disalin ke dokumen / tiket.
--   M_BENEFIT_LIFE_LENGKAP.csv  - ID + JSONDATA (satu-satunya salinan Number, pxObjClass, px*, dan bentuk JSON asli)
--   BENEFIT_LIFE_SEBELUM.csv    - 2 kolom view (ID, BENEFIT) - acuan perbandingan (c)
SET MARKUP CSV ON QUOTE ON
SET FEEDBACK OFF TERMOUT OFF PAGESIZE 0 LONG 2000000 LONGCHUNKSIZE 32767 LINESIZE 32767 TRIMSPOOL ON
SPOOL M_BENEFIT_LIFE_LENGKAP.csv
SELECT ID, JSONDATA FROM POOLDATA.M_BENEFIT_LIFE ORDER BY ID;
SPOOL OFF
SPOOL BENEFIT_LIFE_SEBELUM.csv
SELECT ID, BENEFIT FROM POOLDATA.BENEFIT_LIFE ORDER BY ID;
SPOOL OFF
SET TERMOUT ON FEEDBACK ON MARKUP CSV OFF
PROMPT cadangan selesai: M_BENEFIT_LIFE_LENGKAP.csv, BENEFIT_LIFE_SEBELUM.csv - periksa keduanya terbuka dan berisi 10 baris
SELECT COUNT(*) AS N_TABEL FROM POOLDATA.M_BENEFIT_LIFE;
SELECT COUNT(*) AS N_VIEW FROM POOLDATA.BENEFIT_LIFE;
