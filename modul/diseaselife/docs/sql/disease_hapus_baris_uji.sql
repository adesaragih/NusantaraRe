-- LANGKAH-WO-DISEASELIFE.md (a) D2 - MENULIS DATA (satu-satunya SQL tulis data modul Disease Life). Dijalankan WO
-- SESUDAH cadangan P3 dan SEBELUM -migrate, dari folder cadangan-diseaselife-<tanggal>\ :
--   sqlplus POOLDATA@<db> @"D:\NUSARE DEV\NUSARE\NusantaraRe\modul\diseaselife\docs\sql\disease_hapus_baris_uji.sql"
--
-- Keputusan work owner 08-10-2026 D2: baris uji ID 102051 `TEST123 / Sakit` (dibuat lewat prosedur PEGA_DISEASE_LIFE;
-- ID-nya kembar dengan baris impor `C718 ...`) DIHAPUS supaya PK_DISEASE_LIFE (081) dapat berdiri. Baris impor 102051
-- TIDAK disentuh: syarat DELETE memuat KETIGA kolom persis (ID, ICD_CODE, DISEASE huruf persis).
-- Tabel klaim menyimpan TEKS ICD / Disease, bukan ID - menghapus baris ini tidak mematahkan data klaim (fakta WO;
-- dis_bukti.sql kueri G).
--
-- Urutan:
--   1. cadangan baris itu ke DISEASE_LIFE_BARIS_UJI.csv (selain cadangan P3 seluruh tabel);
--   2. blok PL/SQL: DELETE persis satu baris; bila jumlahnya BUKAN 1 -> ROLLBACK dan berhenti dengan galat (ORA-20001)
--      - nol baris berubah; bila tepat 1 -> COMMIT;
--   3. periksa sesudahnya: baris uji NOL, baris impor 102051 tetap SATU, cacah total turun tepat 1 (97.585).
-- Aman diulang: putaran kedua menemukan NOL baris -> berhenti ORA-20001 tanpa mengubah apa pun (itu tanda D2 SUDAH
-- selesai; periksa langkah 3).
SET LINESIZE 200 PAGESIZE 100 FEEDBACK ON
WHENEVER SQLERROR EXIT FAILURE ROLLBACK
PROMPT 1. Cadangan baris uji (DISEASE_LIFE_BARIS_UJI.csv) dan cacah sebelum
SET MARKUP CSV ON QUOTE ON
SET TERMOUT OFF FEEDBACK OFF PAGESIZE 0 LINESIZE 32767 TRIMSPOOL ON
SPOOL DISEASE_LIFE_BARIS_UJI.csv
SELECT ID, ICD_CODE, DISEASE FROM POOLDATA.DISEASE_LIFE WHERE ID = '102051' ORDER BY ICD_CODE;
SPOOL OFF
SET MARKUP CSV OFF
SET TERMOUT ON FEEDBACK ON PAGESIZE 100 LINESIZE 200
SELECT COUNT(*) AS N_SEBELUM FROM POOLDATA.DISEASE_LIFE;
SELECT ID, ICD_CODE, DISEASE FROM POOLDATA.DISEASE_LIFE WHERE ID = '102051' ORDER BY ICD_CODE;
PROMPT 2. DELETE persis satu baris (ROLLBACK + berhenti bila bukan 1)
DECLARE
  n NUMBER;
BEGIN
  DELETE FROM POOLDATA.DISEASE_LIFE
   WHERE ID = '102051' AND ICD_CODE = 'TEST123' AND DISEASE = 'Sakit';
  n := SQL%ROWCOUNT;
  IF n <> 1 THEN
    ROLLBACK;
    RAISE_APPLICATION_ERROR(-20001, 'D2 dibatalkan: DELETE baris uji 102051 TEST123 / Sakit mengenai ' || n ||
      ' baris, mau tepat 1 - nol baris berubah (ROLLBACK). Periksa dis_bukti.sql kueri C / G.');
  END IF;
  COMMIT;
END;
/
PROMPT 3. Sesudah: baris uji NOL, baris impor 102051 tetap SATU, cacah = sebelum - 1 (pembanding 97.585)
SELECT COUNT(*) AS N_BARIS_UJI FROM POOLDATA.DISEASE_LIFE WHERE ID = '102051' AND ICD_CODE = 'TEST123';
SELECT COUNT(*) AS N_ID_102051 FROM POOLDATA.DISEASE_LIFE WHERE ID = '102051';
SELECT COUNT(*) AS N_SESUDAH FROM POOLDATA.DISEASE_LIFE;
SELECT ID, COUNT(*) AS N FROM POOLDATA.DISEASE_LIFE GROUP BY ID HAVING COUNT(*) > 1;
EXIT
