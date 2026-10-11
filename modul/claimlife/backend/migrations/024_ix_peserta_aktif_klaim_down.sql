-- Jalur mundur 024: buang indeks fungsi pemeriksa peserta aktif (berpelindung, aman diulang).

DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_INDEXES
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_CLAIMLF_PREMIUMLIST_DETAIL' AND INDEX_NAME = 'IX_PLD_TRIM_PL_SERT';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'DROP INDEX {skema}.IX_PLD_TRIM_PL_SERT';
  END IF;
END;
/
