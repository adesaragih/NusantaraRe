-- T_WORK_POLIS.STATUS -> STATUS_WORK - JARING PENGAMAN, biasanya tanpa efek.
--
-- ⚠️ RALAT 01-10-2026: berkas ini lahir dengan anggapan T_WORK_POLIS di DEV
-- diubah DI LUAR repo. Keliru - perubahan itu milik
-- `059_seragam_kolom_t_work_polis.sql` (tabel kerja T_WORK_* seragam), yang
-- masuk lewat merge origin/dev sesudah pemeriksaan pertama. 059 itu sudah
-- menambah STATUS_WORK, menyalin isinya, dan membuang STATUS; sesudahnya
-- langkah ini TIDAK berbuat apa-apa (kolom `STATUS` sudah tidak ada).
-- Dipertahankan, bukan dihapus: ia mungkin sudah tercatat di T_MIGRASI DEV.
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
