-- LANGKAH-WO-COVERLIFE.md (a) P3 - BACA-SAJA. CADANGAN WAJIB SEBELUM -migrate 085-086 / 957: ID + JSONDATA UTUH tabel
-- M_COVER_LIFE dan isi view COVER_LIFE, ke folder kerja SQL*Plus yang dibuat lebih dulu:
-- D:\NUSARE DEV\NUSARE\cadangan-coverlife-<tanggal>\ (mis. cadangan-coverlife-20261008). Tidak ada cadangan = TIDAK BOLEH
-- -migrate. Dijalankan sebagai BERKAS (@...) supaya TERMOUT OFF berlaku. SQL*Plus 12.2+ atau SQLcl (MARKUP CSV).
-- ⛔ JSONDATA dapat memuat kunci px* (identitas operator): simpan di tempat cadangan DBA, jangan disalin ke dokumen /
-- tiket.
--   M_COVER_LIFE_LENGKAP.csv - ID tabel + JSONDATA (satu-satunya salinan pxObjClass, pyRuleHarness, bentuk JSON asli)
--   COVER_LIFE_SEBELUM.csv   - 3 kolom view (ID, COVER, NOTE) - acuan (c)
SET MARKUP CSV ON QUOTE ON
SET FEEDBACK OFF TERMOUT OFF PAGESIZE 0 LONG 2000000 LONGCHUNKSIZE 32767 LINESIZE 32767 TRIMSPOOL ON
SPOOL M_COVER_LIFE_LENGKAP.csv
SELECT ID, JSONDATA FROM POOLDATA.M_COVER_LIFE ORDER BY ID;
SPOOL OFF
SPOOL COVER_LIFE_SEBELUM.csv
SELECT ID, COVER, NOTE FROM POOLDATA.COVER_LIFE ORDER BY ID;
SPOOL OFF
SET TERMOUT ON FEEDBACK ON MARKUP CSV OFF
PROMPT cadangan selesai: M_COVER_LIFE_LENGKAP.csv, COVER_LIFE_SEBELUM.csv - periksa keduanya terbuka dan berisi 4 baris
SELECT COUNT(*) AS N_TABEL FROM POOLDATA.M_COVER_LIFE;
SELECT COUNT(*) AS N_VIEW FROM POOLDATA.COVER_LIFE;
