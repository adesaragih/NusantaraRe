-- 080 - Disease Life (diseaselife): sequence ID aplikasi SEQ_DISEASE_LIFE (keputusan work owner 08-10-2026 D1.1).
--
-- Tabel DISEASE_LIFE (ID VARCHAR2(100), ICD_CODE, DISEASE) SUDAH flat dan TIDAK diubah kolomnya (D1). ID baru =
-- TO_CHAR(SEQ_DISEASE_LIFE.NEXTVAL). Rumus Pega `'1' || LPAD(M_DISEASE_LIFE_SEQ.NEXTVAL, 5, '0')` TIDAK dipakai: di DEV
-- sequence itu (last_number 2052) menghasilkan 102051 yang sudah dipakai data impor (ID kembar) dan 102052, 102053, ...
-- juga sudah terpakai (fakta WO 08-10-2026). M_DISEASE_LIFE_SEQ tidak disentuh.
-- Nilai awal = ID angka tertinggi + 1, DIHITUNG di basis data tempat migrasi berjalan - bentuk PERSIS 923
-- (`migrasi.BacaSequenceDariKueri`): DEV mulai 197586 (ID tertinggi 197585), PROD sesudah ID terbesarnya sendiri. ID
-- bukan angka diabaikan. Penanya katalog dulu: pengulangan sesudah gagal di tengah tidak mati di ORA-00955. NOCACHE:
-- nomor hilang saat instance mati lebih mahal daripada kecepatan. Aplikasi tetap memeriksa ID belum terpakai sebelum
-- INSERT (penulis Pega PEGA_DISEASE_LIFE masih VALID dan dapat menulis ID-nya sendiri).
-- K0: migrasi MODUL (rentang 080-084, dipinjam dari jatah premiumlistlife) - di skema baru berjalan SEBELUM 900; hanya
-- membaca DISEASE_LIFE (objek Pega, sudah ada).
--
-- ⛔ NOL `COMMIT` (ADR-U-0029). `-migrate` dijalankan work owner (LANGKAH-WO-DISEASELIFE.md).
DECLARE
  n    NUMBER;
  awal NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_SEQUENCES
   WHERE SEQUENCE_OWNER = UPPER('{skema}') AND SEQUENCE_NAME = 'SEQ_DISEASE_LIFE';
  IF n = 0 THEN
    SELECT NVL(MAX(TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$'))), 0) + 1 INTO awal FROM {skema}.DISEASE_LIFE;
    EXECUTE IMMEDIATE 'CREATE SEQUENCE {skema}.SEQ_DISEASE_LIFE START WITH ' || awal || ' INCREMENT BY 1 NOCACHE NOCYCLE';
  END IF;
END;
/
