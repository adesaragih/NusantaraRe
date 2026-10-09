-- Mundur 081: buang PK_DISEASE_LIFE (indeks unik bernama sama ikut terbuang; NOT NULL yang Oracle tambahkan untuk PK
-- ikut terbuang sehingga ID kembali NULLABLE seperti asli - LANGKAH-WO-DISEASELIFE.md memeriksanya). Data tidak
-- berubah. Blok berpelindung katalog: aman DIULANG (tanpa ORA-02443 sesudah PK terbuang).
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_CONSTRAINTS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'DISEASE_LIFE' AND CONSTRAINT_NAME = 'PK_DISEASE_LIFE';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.DISEASE_LIFE DROP CONSTRAINT PK_DISEASE_LIFE';
  END IF;
END;
/
