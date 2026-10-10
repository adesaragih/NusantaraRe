-- 923 - sequence ID baru modul R/I Rate Life (K2 keputusan work owner 05-10-2026: "dari sequence Oracle").
--
-- SEQ_M_RATE_LIFE_SUMMARY - ID ringkasan baru M_RATE_LIFE_SUMMARY (Add dan Simpan Upload); SEQ_M_RATE_LIFE - ID baris
-- rate baru M_RATE_LIFE (Simpan Upload). Keduanya ID VARCHAR2(10) teks angka. Nilai awal = ID angka tertinggi tabel
-- masing-masing + 1, DIHITUNG di basis data tempat migrasi berjalan (bentuk `migrasi.BacaSequenceDariKueri`, seperti
-- 910 dan 880 Aggregate) - DEV dan PROD masing-masing mulai sesudah ID terbesarnya sendiri. ID bukan angka diabaikan.
-- Nomor yang tetap terpakai (Pega menulis ID-nya sendiri) dilewati aplikasi. Penanya katalog dulu: pengulangan sesudah
-- gagal di tengah tidak mati di ORA-00955. Tanpa DROP - kedua sequence baru. NOCACHE: nomor hilang saat instance mati
-- lebih mahal daripada kecepatan.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029). `-migrate` dijalankan work owner.
DECLARE
  n    NUMBER;
  awal NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_SEQUENCES
   WHERE SEQUENCE_OWNER = UPPER('{skema}') AND SEQUENCE_NAME = 'SEQ_M_RATE_LIFE_SUMMARY';
  IF n = 0 THEN
    SELECT NVL(MAX(TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$'))), 0) + 1 INTO awal FROM {skema}.M_RATE_LIFE_SUMMARY;
    EXECUTE IMMEDIATE 'CREATE SEQUENCE {skema}.SEQ_M_RATE_LIFE_SUMMARY START WITH ' || awal || ' INCREMENT BY 1 NOCACHE NOCYCLE';
  END IF;
END;
/
DECLARE
  n    NUMBER;
  awal NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_SEQUENCES
   WHERE SEQUENCE_OWNER = UPPER('{skema}') AND SEQUENCE_NAME = 'SEQ_M_RATE_LIFE';
  IF n = 0 THEN
    SELECT NVL(MAX(TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$'))), 0) + 1 INTO awal FROM {skema}.M_RATE_LIFE;
    EXECUTE IMMEDIATE 'CREATE SEQUENCE {skema}.SEQ_M_RATE_LIFE START WITH ' || awal || ' INCREMENT BY 1 NOCACHE NOCYCLE';
  END IF;
END;
/
