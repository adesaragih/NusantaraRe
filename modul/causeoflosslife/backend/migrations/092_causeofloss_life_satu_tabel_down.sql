-- Mundur 092: kembalikan JSONDATA dari kolom CAUSEOFLOSS (JSON_OBJECT, kunci "CauseofLoss" - huruf persis kunci Pega,
-- supaya view yang dibuat ulang 090_down (`a.JSONDATA.CauseofLoss`) membacanya; ABSENT ON NULL - baris 100001 kembali
-- tanpa kunci, view tetap membacanya NULL) dan constraint IS JSON bernama asli ENSURE_M_CAUSEOFLOSS_LIFE_JSON.
-- Sesudahnya 091_down membuang CAUSEOFLOSS dan 090_down mengembalikan nama + view.
-- HANYA bila 092 TERCATAT di T_MIGRASI (LANGKAH-WO-CAUSEOFLOSSLIFE.md, Pemulihan). Kunci pxObjClass, pyRuleHarness, dan
-- bentuk JSON asli (termasuk `"CauseofLoss":""` baris 100001) hanya dari cadangan P3. JSONDATA dikembalikan NULLABLE
-- (LANGKAH-WO (a) mencatat NULLABLE aslinya).
-- Aman DIULANG: kolom dan constraint lewat blok berpelindung katalog; UPDATE menulis nilai yang sama.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'CAUSEOFLOSS_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n = 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.CAUSEOFLOSS_LIFE ADD (JSONDATA CLOB)';
  END IF;
END;
/
UPDATE {skema}.CAUSEOFLOSS_LIFE m
   SET JSONDATA = JSON_OBJECT('CauseofLoss' VALUE m.CAUSEOFLOSS ABSENT ON NULL RETURNING CLOB)
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_CONSTRAINTS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'CAUSEOFLOSS_LIFE' AND CONSTRAINT_NAME = 'ENSURE_M_CAUSEOFLOSS_LIFE_JSON';
  IF n = 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.CAUSEOFLOSS_LIFE ADD CONSTRAINT ENSURE_M_CAUSEOFLOSS_LIFE_JSON CHECK (JSONDATA IS JSON)';
  END IF;
END;
/
