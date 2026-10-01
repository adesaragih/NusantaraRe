-- Jalur mundur 063 - kembalikan nama STATUS bila STATUS_WORK ada.
--
-- ⚠️ Di skema baru (050 sudah membuat STATUS_WORK) langkah ini MENGHASILKAN
-- bentuk 050 lama; jalur mundur 050 sesudahnya membuang tabelnya juga.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_WORK_POLIS' AND COLUMN_NAME = 'STATUS_WORK';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.T_WORK_POLIS RENAME COLUMN STATUS_WORK TO STATUS';
  END IF;
END;
/
