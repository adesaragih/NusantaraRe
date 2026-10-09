-- 948 - Plan SATU tabel PRODUCT_TYPE_LIFE (langkah 3 dari 3).
--
-- Keputusan work owner 08-10-2026 K1/K2. Lima blok berpelindung katalog, berurutan, SETIAP pernyataan aman DIULANG:
--   1. ISI kelima kolom dari JSONDATA APA ADANYA - ekspresi notasi titik PERSIS teks view lama (`a.JSONDATA.CoverName`,
--      `.Business`, `.BusinessID`, `.Benefit`, `.BenefitID`; 946_down), peka huruf. Nilai yang tidak muat lebar 947
--      menghentikan langkah ORA-12899 (tidak dipotong). Bukti baca-saja: modul/planlife/docs/sql/plan_bukti_k1.sql;
--   2. PEMERIKSAAN K1.4 - BERHENTI bila ada ID NULL, ID ganda, atau ID tabel <> JSONDATA.ID (view lama membaca ID dari
--      JSON, jadi ID itulah yang dirujuk pembaca lain);
--   3. PEMERIKSAAN K2 - BERHENTI bila ada kolom yang NULL / berbeda padahal JSON-nya memuat kunci itu.
--      Bentuk 2 dan 3 = UPDATE yang SENGAJA gagal: baris seperti itu diberi COVERNAME 201 huruf
--      (RPAD(CHR(88), 201, CHR(88))) -> ORA-12899 (COVERNAME VARCHAR2(200)) -> pelari berhenti, 948 tidak tercatat,
--      JSONDATA utuh. Nol baris = nol perubahan. Tanpa tanda kutip (teks EXECUTE IMMEDIATE blok berpelindung tidak boleh
--      memuatnya - polaPerintahKatalog); bentuk ID = NULL milik 944 tidak dipakai karena ID di sini masih NULLABLE;
--   4. PK_PRODUCT_TYPE_LIFE (ID) bila belum ada (ALL_CONSTRAINTS) - PRIMARY KEY sekaligus menjadikan ID NOT NULL (Oracle
--      menambah NOT NULL pada kolom PK), jadi tidak perlu MODIFY terpisah yang mati ORA-01442 bila diulang;
--   5. buang JSONDATA + constraint IS JSON (ENSURE_M_PRODUCT_TYPE_LIFE) lewat CASCADE CONSTRAINTS.
-- Kunci px* / pyRuleHarness (bukan kolom view) TIDAK dipindah - hanya di cadangan P3. Prosedur PEGA_M_PRODUCT_TYPE_LIFE
-- menjadi INVALID (K3, diterima WO); sequence M_PRODUCT_TYPE_LIFE_SEQ TETAP. Tanpa FK (K1).
-- Volume DEV 32 baris. NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'PRODUCT_TYPE_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'UPDATE {skema}.PRODUCT_TYPE_LIFE m SET m.COVERNAME = m.JSONDATA.CoverName, m.BUSINESS = m.JSONDATA.Business, m.BUSINESSID = m.JSONDATA.BusinessID, m.BENEFIT = m.JSONDATA.Benefit, m.BENEFITID = m.JSONDATA.BenefitID';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'PRODUCT_TYPE_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'UPDATE {skema}.PRODUCT_TYPE_LIFE m SET m.COVERNAME = RPAD(CHR(88), 201, CHR(88)) WHERE m.ID IS NULL OR m.JSONDATA.ID IS NULL OR m.ID <> m.JSONDATA.ID OR m.ID IN (SELECT d.ID FROM {skema}.PRODUCT_TYPE_LIFE d GROUP BY d.ID HAVING COUNT(*) > 1)';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'PRODUCT_TYPE_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'UPDATE {skema}.PRODUCT_TYPE_LIFE m SET m.COVERNAME = RPAD(CHR(88), 201, CHR(88)) WHERE (m.JSONDATA.CoverName IS NOT NULL AND (m.COVERNAME IS NULL OR m.COVERNAME <> m.JSONDATA.CoverName)) OR (m.JSONDATA.Business IS NOT NULL AND (m.BUSINESS IS NULL OR m.BUSINESS <> m.JSONDATA.Business)) OR (m.JSONDATA.BusinessID IS NOT NULL AND (m.BUSINESSID IS NULL OR m.BUSINESSID <> m.JSONDATA.BusinessID)) OR (m.JSONDATA.Benefit IS NOT NULL AND (m.BENEFIT IS NULL OR m.BENEFIT <> m.JSONDATA.Benefit)) OR (m.JSONDATA.BenefitID IS NOT NULL AND (m.BENEFITID IS NULL OR m.BENEFITID <> m.JSONDATA.BenefitID))';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_CONSTRAINTS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'PRODUCT_TYPE_LIFE' AND CONSTRAINT_NAME = 'PK_PRODUCT_TYPE_LIFE';
  IF n = 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.PRODUCT_TYPE_LIFE ADD CONSTRAINT PK_PRODUCT_TYPE_LIFE PRIMARY KEY (ID)';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'PRODUCT_TYPE_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.PRODUCT_TYPE_LIFE DROP COLUMN JSONDATA CASCADE CONSTRAINTS';
  END IF;
END;
/
