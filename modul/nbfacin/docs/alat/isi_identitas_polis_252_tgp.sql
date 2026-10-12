-- =====================================================================================================================
-- BACKFILL identitas polis T_GENERAL_POLIS untuk SEMUA baris T_GENERAL_POLIS (case NB FacIn sistem baru) yang polisnya
-- SUDAH terbit. Permintaan work owner 12-10-2026: "buatin script nya untuk yg ada di T_GENERAL_POLIS aja" (varian sederhana
-- isi_identitas_polis_252.sql: tanpa saringan T_WORK_POLIS / LINI, skema langsung POOLDATA).
--
-- ⛔ DITULIS SAJA. Jalankan SESUDAH migrasi 252 (kolom POLICY_NO / PROD_DATE_TIME / PROD_KE). TANPA COMMIT - COMMIT / ROLLBACK
-- oleh DBA sesudah BAGIAN 2 cocok. Jalankan per pernyataan (atau Run Script).
--
-- Sumber [terverifikasi] = sama dengan penulisan Selesai Policy (repository/selesaipolicy.go):
--   POLICY_NO      = JSON_POLIS.NOPOLIS
--   PROD_DATE_TIME = JSON_POLIS.TGL_PROD (jam WIB) -> teks Pega UTC 'YYYYMMDDTHHMMSS.000 GMT'
--   PROD_KE        = JSON_POLIS.PRODKE (bukan angka -> NULL)
-- Baris JSON_POLIS case = IDPEGA polos ATAU berawalan 'ASM-FW-GISFW-WORK ', TGL_INPUT terbaru.
-- Hanya baris yang POLICY_NO masih NULL (tidak menimpa hasil sistem baru) dan punya JSON_POLIS.
-- =====================================================================================================================

-- ---- BAGIAN 0: hitung sebelum (catat hasilnya) ------------------------------------------------------------------------
SELECT 'Semua baris T_GENERAL_POLIS' AS HAL, COUNT(*) AS N FROM POOLDATA.T_GENERAL_POLIS
UNION ALL
SELECT 'POLICY_NO sudah terisi', COUNT(*) FROM POOLDATA.T_GENERAL_POLIS WHERE POLICY_NO IS NOT NULL
UNION ALL
SELECT 'AKAN DIISI (POLICY_NO NULL + ada JSON_POLIS)', COUNT(*) FROM POOLDATA.T_GENERAL_POLIS g
 WHERE g.POLICY_NO IS NULL
   AND EXISTS (SELECT 1 FROM POOLDATA.JSON_POLIS j WHERE j.IDPEGA IN (g.ID, 'ASM-FW-GISFW-WORK ' || g.ID));

-- Pratinjau nilai yang AKAN ditulis
SELECT g.ID,
       (SELECT MAX(j.NOPOLIS) KEEP (DENSE_RANK LAST ORDER BY j.TGL_INPUT) FROM POOLDATA.JSON_POLIS j
         WHERE j.IDPEGA IN (g.ID, 'ASM-FW-GISFW-WORK ' || g.ID)) AS POLICY_NO_BARU,
       (SELECT TO_CHAR(SYS_EXTRACT_UTC(FROM_TZ(CAST(MAX(j.TGL_PROD) KEEP (DENSE_RANK LAST ORDER BY j.TGL_INPUT) AS TIMESTAMP),
               '+07:00')), 'YYYYMMDD"T"HH24MISS".000 GMT"') FROM POOLDATA.JSON_POLIS j
         WHERE j.IDPEGA IN (g.ID, 'ASM-FW-GISFW-WORK ' || g.ID)) AS PROD_DATE_TIME_BARU,
       (SELECT MAX(j.PRODKE) KEEP (DENSE_RANK LAST ORDER BY j.TGL_INPUT) FROM POOLDATA.JSON_POLIS j
         WHERE j.IDPEGA IN (g.ID, 'ASM-FW-GISFW-WORK ' || g.ID)) AS PRODKE_SUMBER
  FROM POOLDATA.T_GENERAL_POLIS g
 WHERE g.POLICY_NO IS NULL
   AND EXISTS (SELECT 1 FROM POOLDATA.JSON_POLIS j WHERE j.IDPEGA IN (g.ID, 'ASM-FW-GISFW-WORK ' || g.ID))
 ORDER BY g.ID;

-- ---- BAGIAN 1: isi (jumlah baris terubah harus = "AKAN DIISI" di BAGIAN 0) ----------------------------------------------
UPDATE POOLDATA.T_GENERAL_POLIS g
   SET g.POLICY_NO = (SELECT MAX(j.NOPOLIS) KEEP (DENSE_RANK LAST ORDER BY j.TGL_INPUT) FROM POOLDATA.JSON_POLIS j
                       WHERE j.IDPEGA IN (g.ID, 'ASM-FW-GISFW-WORK ' || g.ID)),
       g.PROD_DATE_TIME = (SELECT TO_CHAR(SYS_EXTRACT_UTC(FROM_TZ(CAST(MAX(j.TGL_PROD) KEEP (DENSE_RANK LAST ORDER BY j.TGL_INPUT)
                            AS TIMESTAMP), '+07:00')), 'YYYYMMDD"T"HH24MISS".000 GMT"') FROM POOLDATA.JSON_POLIS j
                            WHERE j.IDPEGA IN (g.ID, 'ASM-FW-GISFW-WORK ' || g.ID)),
       g.PROD_KE = (SELECT CASE WHEN REGEXP_LIKE(MAX(j.PRODKE) KEEP (DENSE_RANK LAST ORDER BY j.TGL_INPUT), '^[0-9]{1,5}$')
                         THEN TO_NUMBER(MAX(j.PRODKE) KEEP (DENSE_RANK LAST ORDER BY j.TGL_INPUT)) END FROM POOLDATA.JSON_POLIS j
                     WHERE j.IDPEGA IN (g.ID, 'ASM-FW-GISFW-WORK ' || g.ID))
 WHERE g.POLICY_NO IS NULL
   AND EXISTS (SELECT 1 FROM POOLDATA.JSON_POLIS j WHERE j.IDPEGA IN (g.ID, 'ASM-FW-GISFW-WORK ' || g.ID));

-- ---- BAGIAN 2: hitung sesudah (sebelum COMMIT) ------------------------------------------------------------------------
-- Harapan: "POLICY_NO sudah terisi" = angka sebelum + jumlah BAGIAN 1; "AKAN DIISI" = 0.
SELECT 'POLICY_NO sudah terisi' AS HAL, COUNT(*) AS N FROM POOLDATA.T_GENERAL_POLIS WHERE POLICY_NO IS NOT NULL
UNION ALL
SELECT 'AKAN DIISI (harus 0)', COUNT(*) FROM POOLDATA.T_GENERAL_POLIS g
 WHERE g.POLICY_NO IS NULL
   AND EXISTS (SELECT 1 FROM POOLDATA.JSON_POLIS j WHERE j.IDPEGA IN (g.ID, 'ASM-FW-GISFW-WORK ' || g.ID))
UNION ALL
SELECT 'PROD_KE NULL tapi POLICY_NO terisi', COUNT(*) FROM POOLDATA.T_GENERAL_POLIS WHERE POLICY_NO IS NOT NULL AND PROD_KE IS NULL;

SELECT ID, POLICY_NO, PROD_DATE_TIME, PROD_KE FROM POOLDATA.T_GENERAL_POLIS WHERE ID = 'NB-151029';

-- Cocok -> COMMIT; tidak cocok -> ROLLBACK.
