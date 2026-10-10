-- Mundur 944: kembalikan JSONDATA dari kolom BENEFIT (JSON_OBJECT, kunci "Benefit" - huruf persis kunci Pega, supaya
-- view yang dibuat ulang 942_down (`a.JSONDATA.Benefit`) membacanya; ABSENT ON NULL) dan constraint IS JSON bernama
-- asli ENSURE_M_BENEFIT_LIFE_JSON. Sesudahnya 943_down membuang BENEFIT dan 942_down mengembalikan nama + view.
-- HANYA bila 944 TERCATAT di T_MIGRASI (LANGKAH-WO-BENEFITLIFE.md, Pemulihan). Kunci Number, pxObjClass, px* lain, dan
-- bentuk JSON asli hanya dari cadangan P3. JSONDATA dikembalikan NULLABLE (LANGKAH-WO (a) mencatat NULLABLE aslinya).
-- Aman DIULANG: kolom dan constraint lewat blok berpelindung katalog; UPDATE menulis nilai yang sama.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'BENEFIT_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n = 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.BENEFIT_LIFE ADD (JSONDATA CLOB)';
  END IF;
END;
/
UPDATE {skema}.BENEFIT_LIFE m
   SET JSONDATA = JSON_OBJECT('Benefit' VALUE m.BENEFIT ABSENT ON NULL RETURNING CLOB)
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_CONSTRAINTS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'BENEFIT_LIFE' AND CONSTRAINT_NAME = 'ENSURE_M_BENEFIT_LIFE_JSON';
  IF n = 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.BENEFIT_LIFE ADD CONSTRAINT ENSURE_M_BENEFIT_LIFE_JSON CHECK (JSONDATA IS JSON)';
  END IF;
END;
/
