-- LANGKAH-WO-RIRISKLIFE.md (a) P3 - BACA-SAJA. Cadangan CSV kedua tabel SEBELUM 935-941 (ID + JSONDATA + kolom datar),
-- di folder kerja SQL*Plus yang dibuat lebih dulu: cadangan-ririsklife-<tanggal>\ (mis. cadangan-ririsklife-20261008).
-- Tidak ada cadangan = tidak boleh -migrate. Dijalankan sebagai BERKAS (@...) supaya TERMOUT OFF berlaku. SQL*Plus 12.2+
-- atau SQLcl (MARKUP CSV). ⛔ Berkas memuat OPERATORID (nama orang): simpan di tempat cadangan DBA, jangan dibagikan.
--   M_RIRISK_LIFE_SUMMARY_LENGKAP.csv - ID + JSONDATA (satu-satunya salinan pxObjClass dan bentuk JSON asli)
--   M_RIRISK_LIFE_LENGKAP.csv         - ID + JSONDATA + kolom datar warisan (basi) apa adanya
--   RIRISK_LIFE_SUMMARY_SEBELUM.csv   - 4 kolom view ringkasan (acuan perbandingan (c))
--   RIRISK_LIFE_SEBELUM.csv           - 7 kolom view rincian tanpa AGE (kosong seluruhnya), RISK lewat konversi
--                                       ber-format NLS eksplisit lalu TM9 (= bentuk sesudah 940; acuan (c))
SET MARKUP CSV ON QUOTE ON
SET FEEDBACK OFF TERMOUT OFF PAGESIZE 0 LONG 2000000 LONGCHUNKSIZE 32767 LINESIZE 32767 TRIMSPOOL ON
SPOOL M_RIRISK_LIFE_SUMMARY_LENGKAP.csv
SELECT ID, JSONDATA FROM POOLDATA.M_RIRISK_LIFE_SUMMARY ORDER BY ID;
SPOOL OFF
SPOOL M_RIRISK_LIFE_LENGKAP.csv
SELECT ID, JSONDATA, IDUSEDBY, USEDBY, YEAR, MONTH, TO_CHAR(RISK, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''') AS RISK, CONTRACT
  FROM POOLDATA.M_RIRISK_LIFE ORDER BY ID;
SPOOL OFF
SPOOL RIRISK_LIFE_SUMMARY_SEBELUM.csv
SELECT ID, USEDBY, MODIFIEDDATE, OPERATORID FROM POOLDATA.RIRISK_LIFE_SUMMARY ORDER BY ID;
SPOOL OFF
SPOOL RIRISK_LIFE_SEBELUM.csv
SELECT ID, IDUSEDBY, USEDBY, YEAR, MONTH,
       TO_CHAR(TO_NUMBER(REPLACE(TRIM(RISK), ',', '.'), '999999999999999999999999999999D9999999999999999999', 'NLS_NUMERIC_CHARACTERS=''.,'''), 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''') AS RISK,
       CONTRACT FROM POOLDATA.RIRISK_LIFE ORDER BY ID;
SPOOL OFF
SET TERMOUT ON FEEDBACK ON MARKUP CSV OFF
PROMPT cadangan selesai: 4 berkas CSV
