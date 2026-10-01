-- T_WORK_POLIS.STATUS -> STATUS_WORK, untuk lingkungan yang memasang 050 LAMA.
--
-- `[keputusan work owner 01-10-2026]` ikuti struktur DEV (kolom bernama
-- STATUS_WORK, pola T_WORK_CLAIM). 050 kini membuat STATUS_WORK langsung;
-- langkah ini hanya mengganti nama bila kolom `STATUS` MASIH ada. Di DEV (sudah
-- STATUS_WORK) dan di skema baru ia tidak berbuat apa-apa.
--
-- Berpelindung katalog (bentuk 901 menu datar); nol `COMMIT` (ADR-U-0029).
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_WORK_POLIS' AND COLUMN_NAME = 'STATUS';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.T_WORK_POLIS RENAME COLUMN STATUS TO STATUS_WORK';
  END IF;
END;
/
