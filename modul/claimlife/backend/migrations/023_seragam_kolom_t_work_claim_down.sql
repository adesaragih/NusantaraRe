-- Jalur mundur 023 - kebalikan urutan maju, seluruhnya berpelindung katalog.
--
-- CASE_ID kembali berisi ID untuk baris klaim (bentuk yang selalu ditulis
-- pendaftaran); baris Komite kembali kosong, seperti sebelum 023. TYPE
-- kembali ke T_WORK_CLAIM: baris klaim dari header-nya, baris Komite dari
-- header klaim induknya (COVER_KEY) - salinan `pxAddChildWork` yang dulu
-- ditulis saat kasus Komite lahir. Lalu TYPE header dibuang, dan POSITION
-- kembali bernama PY_POSITION. Di jalur mundur kolom lewat blok dan RENAME
-- sah: penjaga bentuk membaca jalur maju (`pelanggaranBlokPLSQL`).
--
-- ⚠️ CASEID warisan yang dipakai sebagai ID oleh migrasi klaim lama kembali
-- sebagai CASE_ID = ID - nilainya sama.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029).
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_WORK_CLAIM' AND COLUMN_NAME = 'CASE_ID';
  IF n = 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.T_WORK_CLAIM ADD (CASE_ID VARCHAR2(64))';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_WORK_CLAIM' AND COLUMN_NAME = 'CASE_ID';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'UPDATE {skema}.T_WORK_CLAIM w SET w.CASE_ID = w.ID WHERE w.CASE_ID IS NULL AND EXISTS (SELECT 1 FROM {skema}.T_GENERAL_CLAIM g WHERE g.ID = w.ID)';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_WORK_CLAIM' AND COLUMN_NAME = 'TYPE';
  IF n = 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.T_WORK_CLAIM ADD (TYPE VARCHAR2(32))';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_WORK_CLAIM' AND COLUMN_NAME = 'TYPE';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'UPDATE {skema}.T_WORK_CLAIM w SET w.TYPE = COALESCE((SELECT g.TYPE FROM {skema}.T_GENERAL_CLAIM g WHERE g.ID = w.ID), (SELECT g.TYPE FROM {skema}.T_GENERAL_CLAIM g WHERE g.ID = w.COVER_KEY)) WHERE w.TYPE IS NULL';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_GENERAL_CLAIM' AND COLUMN_NAME = 'TYPE';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.T_GENERAL_CLAIM DROP COLUMN TYPE';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_WORK_CLAIM' AND COLUMN_NAME = 'POSITION';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.T_WORK_CLAIM RENAME COLUMN POSITION TO PY_POSITION';
  END IF;
END;
/
