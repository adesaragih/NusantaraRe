-- Mundur 948: buang PK_PRODUCT_TYPE_LIFE (ID kembali NULLABLE seperti asli - NOT NULL yang ditambahkan Oracle untuk PK
-- ikut terbuang; LANGKAH-WO memeriksanya), kembalikan JSONDATA dari kolom (JSON_OBJECT, kunci huruf PERSIS Pega `ID`,
-- `CoverName`, `Business`, `BusinessID`, `Benefit`, `BenefitID` - view 946_down membacanya; ABSENT ON NULL) dan
-- constraint IS JSON bernama asli ENSURE_M_PRODUCT_TYPE_LIFE. Sesudahnya 947_down membuang kelima kolom dan 946_down
-- mengembalikan nama + view. HANYA bila 948 TERCATAT di T_MIGRASI (LANGKAH-WO-PLANLIFE.md, Pemulihan). Kunci px* dan
-- bentuk JSON asli hanya dari cadangan P3. JSONDATA dikembalikan NULLABLE.
-- Aman DIULANG: PK, kolom, constraint lewat blok berpelindung katalog; UPDATE menulis nilai yang sama.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_CONSTRAINTS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'PRODUCT_TYPE_LIFE' AND CONSTRAINT_NAME = 'PK_PRODUCT_TYPE_LIFE';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.PRODUCT_TYPE_LIFE DROP CONSTRAINT PK_PRODUCT_TYPE_LIFE';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'PRODUCT_TYPE_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n = 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.PRODUCT_TYPE_LIFE ADD (JSONDATA CLOB)';
  END IF;
END;
/
UPDATE {skema}.PRODUCT_TYPE_LIFE m
   SET JSONDATA = JSON_OBJECT('ID' VALUE m.ID, 'CoverName' VALUE m.COVERNAME, 'Business' VALUE m.BUSINESS,
                              'BusinessID' VALUE m.BUSINESSID, 'Benefit' VALUE m.BENEFIT, 'BenefitID' VALUE m.BENEFITID
                              ABSENT ON NULL RETURNING CLOB)
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_CONSTRAINTS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'PRODUCT_TYPE_LIFE' AND CONSTRAINT_NAME = 'ENSURE_M_PRODUCT_TYPE_LIFE';
  IF n = 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.PRODUCT_TYPE_LIFE ADD CONSTRAINT ENSURE_M_PRODUCT_TYPE_LIFE CHECK (JSONDATA IS JSON)';
  END IF;
END;
/
