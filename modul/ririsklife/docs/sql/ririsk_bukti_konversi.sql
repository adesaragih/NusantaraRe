-- LANGKAH-WO-RIRISKLIFE.md (a) - BACA-SAJA, SEBELUM -migrate. BUKTI konversi RISK migrasi 940 (K2) TANPA menulis:
-- ekspresi 940 yang SAMA dijalankan sebagai SELECT di DUA setelan NLS sesi (titik lalu koma). Kedua keluaran harus
-- SAMA: 16.398 baris terkonversi, nol gagal, "921,9" -> 921.9, "580.894351210924" -> 580.894351210924.
-- ALTER SESSION hanya mengubah sesi SQL*Plus ini (tidak menulis data). Dijalankan sebagai BERKAS (@...).
SET LINESIZE 200 PAGESIZE 100 NUMWIDTH 30
SPOOL ririsk_bukti_konversi.txt
-- Konversi 940: TRANSLATE koma DAN titik -> pemisah desimal sesi (karakter pertama TO_CHAR(0.5)), lalu TO_NUMBER tanpa
-- format. Hasil ditampilkan lewat TO_CHAR TM9 ber-NLS eksplisit supaya tampilannya tidak ikut NLS.
ALTER SESSION SET NLS_NUMERIC_CHARACTERS = '.,';
PROMPT A. NLS titik
SELECT COUNT(*) AS N, COUNT(RISK) AS N_RISK,
       COUNT(TO_NUMBER(TRANSLATE(TRIM(RISK), ',.', SUBSTR(TO_CHAR(0.5), 1, 1) || SUBSTR(TO_CHAR(0.5), 1, 1))
                       DEFAULT NULL ON CONVERSION ERROR)) AS N_TERKONVERSI,
       SUM(ORA_HASH(TO_CHAR(TO_NUMBER(TRANSLATE(TRIM(RISK), ',.', SUBSTR(TO_CHAR(0.5), 1, 1) || SUBSTR(TO_CHAR(0.5), 1, 1))
                       DEFAULT NULL ON CONVERSION ERROR), 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''))) AS H_HASIL
  FROM POOLDATA.RIRISK_LIFE;
SELECT TO_CHAR(TO_NUMBER(TRANSLATE('921,9', ',.', SUBSTR(TO_CHAR(0.5), 1, 1) || SUBSTR(TO_CHAR(0.5), 1, 1))), 'TM9',
               'NLS_NUMERIC_CHARACTERS=''.,''') AS KOMA,
       TO_CHAR(TO_NUMBER(TRANSLATE('580.894351210924', ',.', SUBSTR(TO_CHAR(0.5), 1, 1) || SUBSTR(TO_CHAR(0.5), 1, 1))), 'TM9',
               'NLS_NUMERIC_CHARACTERS=''.,''') AS TITIK
  FROM DUAL;
ALTER SESSION SET NLS_NUMERIC_CHARACTERS = ',.';
PROMPT B. NLS koma (angka harus SAMA dengan A)
SELECT COUNT(*) AS N, COUNT(RISK) AS N_RISK,
       COUNT(TO_NUMBER(TRANSLATE(TRIM(RISK), ',.', SUBSTR(TO_CHAR(0.5), 1, 1) || SUBSTR(TO_CHAR(0.5), 1, 1))
                       DEFAULT NULL ON CONVERSION ERROR)) AS N_TERKONVERSI,
       SUM(ORA_HASH(TO_CHAR(TO_NUMBER(TRANSLATE(TRIM(RISK), ',.', SUBSTR(TO_CHAR(0.5), 1, 1) || SUBSTR(TO_CHAR(0.5), 1, 1))
                       DEFAULT NULL ON CONVERSION ERROR), 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''))) AS H_HASIL
  FROM POOLDATA.RIRISK_LIFE;
SELECT TO_CHAR(TO_NUMBER(TRANSLATE('921,9', ',.', SUBSTR(TO_CHAR(0.5), 1, 1) || SUBSTR(TO_CHAR(0.5), 1, 1))), 'TM9',
               'NLS_NUMERIC_CHARACTERS=''.,''') AS KOMA,
       TO_CHAR(TO_NUMBER(TRANSLATE('580.894351210924', ',.', SUBSTR(TO_CHAR(0.5), 1, 1) || SUBSTR(TO_CHAR(0.5), 1, 1))), 'TM9',
               'NLS_NUMERIC_CHARACTERS=''.,''') AS TITIK
  FROM DUAL;
PROMPT C. Baris yang TIDAK terkonversi (harus nol baris)
SELECT ID, RISK FROM POOLDATA.RIRISK_LIFE
 WHERE RISK IS NOT NULL
   AND TO_NUMBER(TRANSLATE(TRIM(RISK), ',.', SUBSTR(TO_CHAR(0.5), 1, 1) || SUBSTR(TO_CHAR(0.5), 1, 1))
                 DEFAULT NULL ON CONVERSION ERROR) IS NULL;
SPOOL OFF
-- Catatan: TRANSLATE(x, ',.', ...) di sini = TRANSLATE(x, CHR(44) || CHR(46), ...) di 940 (blok berpelindung tidak
-- boleh memuat tanda kutip). DEFAULT NULL ON CONVERSION ERROR (Oracle 12.2+) hanya untuk MENGHITUNG yang gagal di sini;
-- 940 sendiri gagal keras (ORA-01722) pada nilai bukan angka.
