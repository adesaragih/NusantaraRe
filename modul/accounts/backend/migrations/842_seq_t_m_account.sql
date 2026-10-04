-- 842 - SEQ_T_M_ACCOUNT, nomor ACC akun baru (keputusan work owner 04-10-2026: "buat seq aja, startnya dari id
-- paling tinggi +1"). Nilai awal dihitung di basis data tempat migrasi berjalan: nomor n terbesar dari ID berpola
-- `ASM-SFAGIS-WORK-ACCOUNT ACC-n`, + 1 (DEV 04-10-2026: 4916561 -> 4916562). Bentuk blok sequence-dari-kueri
-- (`inti/backend/migrasi/sequence_kueri.go`). `ACCOUNT_SEQ` lama TIDAK dipakai: nilainya tertinggal +-4 juta.
-- NOCACHE: nomor yang hilang saat instance mati lebih mahal daripada kecepatan (seperti SEQ_WORK_CLAIM).
DECLARE
  n    NUMBER;
  awal NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_SEQUENCES
   WHERE SEQUENCE_OWNER = UPPER('{skema}') AND SEQUENCE_NAME = 'SEQ_T_M_ACCOUNT';
  IF n = 0 THEN
    SELECT NVL(MAX(TO_NUMBER(REGEXP_SUBSTR(ID, '^ASM-SFAGIS-WORK-ACCOUNT ACC-([0-9]+)$', 1, 1, NULL, 1))), 0) + 1 INTO awal FROM {skema}.T_M_ACCOUNT;
    EXECUTE IMMEDIATE 'CREATE SEQUENCE {skema}.SEQ_T_M_ACCOUNT START WITH ' || awal || ' INCREMENT BY 1 NOCACHE NOCYCLE';
  END IF;
END;
/
