-- 880 - SEQ_AGGREGATE, nomor ID baris AGGREGATE baru (`AGG-n`). Pega membentuk ID lewat
-- `pzGenerateUniqueID(tools, "AGG")` (SaveAggregate_Act langkah 6); aplikasi Go memakai sequence ini, mulai dari nomor
-- AGG terbesar + 1 yang dihitung di basis data tempat migrasi berjalan (DEV 04-10-2026: AGG-33686), seperti
-- SEQ_CLIENT_ORG dan SEQ_T_M_ACCOUNT (keputusan work owner 04-10-2026: "buat seq aja, start-nya dari id max+1").
-- Bentuk blok sequence-dari-kueri (`inti/backend/migrasi/sequence_kueri.go`). Nomor yang sudah terpakai dilewati
-- aplikasi. NOCACHE: nomor yang hilang saat instance mati lebih mahal daripada kecepatan.
DECLARE
  n    NUMBER;
  awal NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_SEQUENCES
   WHERE SEQUENCE_OWNER = UPPER('{skema}') AND SEQUENCE_NAME = 'SEQ_AGGREGATE';
  IF n = 0 THEN
    SELECT NVL(MAX(TO_NUMBER(REGEXP_SUBSTR(ID, '^AGG-([0-9]+)$', 1, 1, NULL, 1))), 0) + 1 INTO awal FROM {skema}.AGGREGATE;
    EXECUTE IMMEDIATE 'CREATE SEQUENCE {skema}.SEQ_AGGREGATE START WITH ' || awal || ' INCREMENT BY 1 NOCACHE NOCYCLE';
  END IF;
END;
/
