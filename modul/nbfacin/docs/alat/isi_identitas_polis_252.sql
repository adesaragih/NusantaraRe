-- =====================================================================================================================
-- BACKFILL identitas polis TAHAP 1 T_GENERAL_POLIS untuk case NB FacIn yang SUDAH terbit (keputusan work owner 12-10-2026: "ini
-- udh jadi polis masa data yg disimpan tidak lengkap", contoh NB-151029; docs/issues/76-general-polis-lengkap.md).
--
-- ⛔ DITULIS SAJA. BELUM PERNAH DIJALANKAN. Jalankan SESUDAH migrasi 252 (kolom POLICY_NO / PROD_DATE_TIME / PROD_KE). Dijalankan DBA
-- di DEV sesudah persetujuan work owner. Bukan migrasi bernomor. TANPA COMMIT - COMMIT / ROLLBACK oleh DBA sesudah BAGIAN 2 cocok.
--
-- Sebelum menjalankan: ganti `{skema}` dengan skema objek nbfacin (T_WORK_POLIS / T_GENERAL_POLIS). JSON_POLIS ditulis berawalan
-- POOLDATA. seperti di kode; sesuaikan bila sinonim berbeda di DEV.
--
-- SUMBER `[terverifikasi]` = sama dengan penulisan Selesai Policy (repository/selesaipolicy.go sqlIdentitasPolis):
--   - POLICY_NO      = JSON_POLIS.NOPOLIS              (SaveJsonPolicyFacIn_Act langkah 12 -> INSERTJSON_JSONPOLIS_FACIN);
--   - PROD_DATE_TIME = JSON_POLIS.TGL_PROD (jam dinding WIB, CARI21 langkah 13) -> teks Pega UTC 'YYYYMMDDTHHMMSS.000 GMT';
--   - PROD_KE        = JSON_POLIS.PRODKE (NB '0'; bukan angka -> NULL).
-- Baris JSON_POLIS case = IDPEGA dua bentuk (polos / 'ASM-FW-GISFW-WORK ' + ID), TGL_INPUT terbaru (sama dengan sqlCariJSONPolis).
-- FACINPRODUCTION tidak dipakai: Pega selalu menulis JSON_POLIS (langkah 15) sebelum FACINPRODUCTION (langkah 23), dan nilai tiga
-- kolom ini ada utuh di JSON_POLIS.
--
-- CAKUPAN: hanya T_GENERAL_POLIS milik case {skema}.T_WORK_POLIS LINI 'FAC' ber-ID 'NB-%', POLICY_NO masih NULL (tidak menimpa
-- hasil Selesai Policy biner baru), dan punya baris JSON_POLIS. Kolom tahap 1 lain (GENERATE_DATE_POLICY_NO, POLICY_MASTER_*,
-- BINDER_RNM, DATE_OFFER) TIDAK diisi: sumbernya bukan jalur Selesai Policy / tidak ditemukan (lihat issues/76).
-- =====================================================================================================================

-- ---- BAGIAN 0: hitung sebelum (simpan hasilnya) --------------------------------------------------------------------
SELECT 'NB FAC punya T_GENERAL_POLIS' AS HAL, COUNT(*) AS N FROM {skema}.T_GENERAL_POLIS g
 WHERE EXISTS (SELECT 1 FROM {skema}.T_WORK_POLIS w WHERE w.ID = g.ID AND w.LINI = 'FAC' AND w.ID LIKE 'NB-%')
UNION ALL SELECT 'POLICY_NO terisi', COUNT(*) FROM {skema}.T_GENERAL_POLIS g
 WHERE g.POLICY_NO IS NOT NULL
   AND EXISTS (SELECT 1 FROM {skema}.T_WORK_POLIS w WHERE w.ID = g.ID AND w.LINI = 'FAC' AND w.ID LIKE 'NB-%')
UNION ALL SELECT 'AKAN DIISI (POLICY_NO NULL + ada JSON_POLIS)', COUNT(*) FROM {skema}.T_GENERAL_POLIS g
 WHERE g.POLICY_NO IS NULL
   AND EXISTS (SELECT 1 FROM {skema}.T_WORK_POLIS w WHERE w.ID = g.ID AND w.LINI = 'FAC' AND w.ID LIKE 'NB-%')
   AND EXISTS (SELECT 1 FROM POOLDATA.JSON_POLIS j WHERE j.IDPEGA IN (g.ID, 'ASM-FW-GISFW-WORK ' || g.ID));

-- Pratinjau (termasuk NB-151029): nilai yang AKAN ditulis.
SELECT g.ID,
       (SELECT MAX(j.NOPOLIS) KEEP (DENSE_RANK LAST ORDER BY j.TGL_INPUT) FROM POOLDATA.JSON_POLIS j
         WHERE j.IDPEGA IN (g.ID, 'ASM-FW-GISFW-WORK ' || g.ID)) AS POLICY_NO_BARU,
       (SELECT TO_CHAR(SYS_EXTRACT_UTC(FROM_TZ(CAST(MAX(j.TGL_PROD) KEEP (DENSE_RANK LAST ORDER BY j.TGL_INPUT) AS TIMESTAMP),
               '+07:00')), 'YYYYMMDD"T"HH24MISS".000 GMT"') FROM POOLDATA.JSON_POLIS j
         WHERE j.IDPEGA IN (g.ID, 'ASM-FW-GISFW-WORK ' || g.ID)) AS PROD_DATE_TIME_BARU,
       (SELECT MAX(j.PRODKE) KEEP (DENSE_RANK LAST ORDER BY j.TGL_INPUT) FROM POOLDATA.JSON_POLIS j
         WHERE j.IDPEGA IN (g.ID, 'ASM-FW-GISFW-WORK ' || g.ID)) AS PRODKE_SUMBER
  FROM {skema}.T_GENERAL_POLIS g
 WHERE g.POLICY_NO IS NULL
   AND EXISTS (SELECT 1 FROM {skema}.T_WORK_POLIS w WHERE w.ID = g.ID AND w.LINI = 'FAC' AND w.ID LIKE 'NB-%')
   AND EXISTS (SELECT 1 FROM POOLDATA.JSON_POLIS j WHERE j.IDPEGA IN (g.ID, 'ASM-FW-GISFW-WORK ' || g.ID))
 ORDER BY g.ID;

-- ---- BAGIAN 1: isi (jumlah baris terubah = angka "AKAN DIISI" BAGIAN 0) -----------------------------------------------
UPDATE {skema}.T_GENERAL_POLIS g
   SET g.POLICY_NO = (SELECT MAX(j.NOPOLIS) KEEP (DENSE_RANK LAST ORDER BY j.TGL_INPUT) FROM POOLDATA.JSON_POLIS j
                       WHERE j.IDPEGA IN (g.ID, 'ASM-FW-GISFW-WORK ' || g.ID)),
       g.PROD_DATE_TIME = (SELECT TO_CHAR(SYS_EXTRACT_UTC(FROM_TZ(CAST(MAX(j.TGL_PROD) KEEP (DENSE_RANK LAST ORDER BY j.TGL_INPUT)
                            AS TIMESTAMP), '+07:00')), 'YYYYMMDD"T"HH24MISS".000 GMT"') FROM POOLDATA.JSON_POLIS j
                            WHERE j.IDPEGA IN (g.ID, 'ASM-FW-GISFW-WORK ' || g.ID)),
       g.PROD_KE = (SELECT CASE WHEN REGEXP_LIKE(MAX(j.PRODKE) KEEP (DENSE_RANK LAST ORDER BY j.TGL_INPUT), '^[0-9]{1,5}$')
                         THEN TO_NUMBER(MAX(j.PRODKE) KEEP (DENSE_RANK LAST ORDER BY j.TGL_INPUT)) END FROM POOLDATA.JSON_POLIS j
                     WHERE j.IDPEGA IN (g.ID, 'ASM-FW-GISFW-WORK ' || g.ID))
 WHERE g.POLICY_NO IS NULL
   AND EXISTS (SELECT 1 FROM {skema}.T_WORK_POLIS w WHERE w.ID = g.ID AND w.LINI = 'FAC' AND w.ID LIKE 'NB-%')
   AND EXISTS (SELECT 1 FROM POOLDATA.JSON_POLIS j WHERE j.IDPEGA IN (g.ID, 'ASM-FW-GISFW-WORK ' || g.ID));

-- ---- BAGIAN 2: hitung sesudah (sebelum COMMIT) -----------------------------------------------------------------------
-- Harapan: "POLICY_NO terisi" = sebelum + jumlah BAGIAN 1; "AKAN DIISI" = 0; "PROD_KE NULL tapi POLICY_NO terisi" = jumlah PRODKE
-- sumber bukan angka (lihat pratinjau).
SELECT 'POLICY_NO terisi' AS HAL, COUNT(*) AS N FROM {skema}.T_GENERAL_POLIS g
 WHERE g.POLICY_NO IS NOT NULL
   AND EXISTS (SELECT 1 FROM {skema}.T_WORK_POLIS w WHERE w.ID = g.ID AND w.LINI = 'FAC' AND w.ID LIKE 'NB-%')
UNION ALL SELECT 'AKAN DIISI (harus 0)', COUNT(*) FROM {skema}.T_GENERAL_POLIS g
 WHERE g.POLICY_NO IS NULL
   AND EXISTS (SELECT 1 FROM {skema}.T_WORK_POLIS w WHERE w.ID = g.ID AND w.LINI = 'FAC' AND w.ID LIKE 'NB-%')
   AND EXISTS (SELECT 1 FROM POOLDATA.JSON_POLIS j WHERE j.IDPEGA IN (g.ID, 'ASM-FW-GISFW-WORK ' || g.ID))
UNION ALL SELECT 'PROD_KE NULL tapi POLICY_NO terisi', COUNT(*) FROM {skema}.T_GENERAL_POLIS g
 WHERE g.POLICY_NO IS NOT NULL AND g.PROD_KE IS NULL
   AND EXISTS (SELECT 1 FROM {skema}.T_WORK_POLIS w WHERE w.ID = g.ID AND w.LINI = 'FAC' AND w.ID LIKE 'NB-%');

SELECT g.ID, g.POLICY_NO, g.PROD_DATE_TIME, g.PROD_KE FROM {skema}.T_GENERAL_POLIS g WHERE g.ID = 'NB-151029';

-- Cocok -> COMMIT; tidak -> ROLLBACK. (Sengaja tidak ditulis di skrip.)
