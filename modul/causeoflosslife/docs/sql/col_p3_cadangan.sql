-- LANGKAH-WO-CAUSEOFLOSSLIFE.md (a) P3 - BACA-SAJA. CADANGAN WAJIB SEBELUM -migrate 090-092 / 955: ID + JSONDATA UTUH
-- tabel M_CAUSEOFLOSS_LIFE dan isi view CAUSEOFLOSS_LIFE, ke folder kerja SQL*Plus yang dibuat lebih dulu:
-- D:\NUSARE DEV\NUSARE\cadangan-causeoflosslife-<tanggal>\ (mis. cadangan-causeoflosslife-20261008). Tidak ada cadangan
-- = TIDAK BOLEH -migrate. Dijalankan sebagai BERKAS (@...) supaya TERMOUT OFF berlaku. SQL*Plus 12.2+ atau SQLcl
-- (MARKUP CSV).
-- ⛔ JSONDATA dapat memuat kunci px* (identitas operator): simpan di tempat cadangan DBA, jangan disalin ke dokumen /
-- tiket.
--   M_CAUSEOFLOSS_LIFE_LENGKAP.csv - ID tabel + JSONDATA (satu-satunya salinan pxObjClass, pyRuleHarness, bentuk JSON
--                                    asli termasuk "CauseofLoss":"" baris 100001)
--   CAUSEOFLOSS_LIFE_SEBELUM.csv   - 2 kolom view (ID, CAUSEOFLOSS) - acuan (c)
SET MARKUP CSV ON QUOTE ON
SET FEEDBACK OFF TERMOUT OFF PAGESIZE 0 LONG 2000000 LONGCHUNKSIZE 32767 LINESIZE 32767 TRIMSPOOL ON
SPOOL M_CAUSEOFLOSS_LIFE_LENGKAP.csv
SELECT ID, JSONDATA FROM POOLDATA.M_CAUSEOFLOSS_LIFE ORDER BY ID;
SPOOL OFF
SPOOL CAUSEOFLOSS_LIFE_SEBELUM.csv
SELECT ID, CAUSEOFLOSS FROM POOLDATA.CAUSEOFLOSS_LIFE ORDER BY ID;
SPOOL OFF
SET TERMOUT ON FEEDBACK ON MARKUP CSV OFF
PROMPT cadangan selesai: M_CAUSEOFLOSS_LIFE_LENGKAP.csv, CAUSEOFLOSS_LIFE_SEBELUM.csv - periksa keduanya terbuka dan berisi 4 baris
SELECT COUNT(*) AS N_TABEL FROM POOLDATA.M_CAUSEOFLOSS_LIFE;
SELECT COUNT(*) AS N_VIEW FROM POOLDATA.CAUSEOFLOSS_LIFE;
