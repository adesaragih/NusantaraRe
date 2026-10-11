-- 202 - T_GENERAL_POLIS.PPN_CHECK NUMBER(5) -> VARCHAR2(10) teks Pega "true" / "false" (tiket 50, ralat 06-10-2026).
--
-- Sebab: butir 68.1 (keputusan work owner 02-10-2026, `loader/aturan.go` kolomTeksMenyimpang) menyimpan `.PPnCheck` TEKS APA
-- ADANYA, dan kode tab Premium Deduction menulis "true" / "false". Migrasi 200 sudah jalan di DEV dengan NUMBER(5) (T_MIGRASI
-- dan ALL_TAB_COLUMNS diperiksa work owner 06-10-2026) sehingga tidak diubah; perubahan tipe lewat langkah baru ini.
--
-- ⚠️ BUKAN `MODIFY` (Oracle menolak NUMBER -> VARCHAR2 pada kolom berisi, ORA-01439) dan BUKAN `RENAME COLUMN` (pembanding
-- STRUKTUR dan penjaga tipe hanya mengenal CREATE TABLE, `ALTER ... ADD (`, dan DROP COLUMN - pola 059 premiumlistlife).
-- Bentuk yang terlihat penjaga: tambah PPN_CHECK_TEKS -> salin (1 -> 'true', 0 -> 'false', NULL tetap NULL) -> buang
-- PPN_CHECK -> tambah PPN_CHECK VARCHAR2(10) -> salin balik -> buang PPN_CHECK_TEKS. UPDATE memakai literal teks, jadi
-- tidak dapat di dalam blok berpelindung katalog (EXECUTE IMMEDIATE tanpa kutip); kedua DROP berpelindung katalog.
--
-- BILA 202 GAGAL DI TENGAH, pengulangan berhenti di ORA-01430 (PPN_CHECK_TEKS sudah ada) dan T_MIGRASI belum mencatat 202.
-- Pulihkan (DBA) menurut keadaan `SELECT COLUMN_NAME, DATA_TYPE FROM ALL_TAB_COLUMNS WHERE OWNER = '<skema>' AND
-- TABLE_NAME = 'T_GENERAL_POLIS' AND COLUMN_NAME IN ('PPN_CHECK', 'PPN_CHECK_TEKS')` - tabelnya di MODUL.md nbfacin bab
-- migrasi 202. Ditulis, TIDAK dijalankan agent. Nol COMMIT (ADR-U-0029).
ALTER TABLE {skema}.T_GENERAL_POLIS ADD (
  PPN_CHECK_TEKS VARCHAR2(10)
)
/
UPDATE {skema}.T_GENERAL_POLIS SET PPN_CHECK_TEKS = CASE PPN_CHECK WHEN 1 THEN 'true' WHEN 0 THEN 'false' END
 WHERE PPN_CHECK IS NOT NULL
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_GENERAL_POLIS' AND COLUMN_NAME = 'PPN_CHECK';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.T_GENERAL_POLIS DROP COLUMN PPN_CHECK';
  END IF;
END;
/
ALTER TABLE {skema}.T_GENERAL_POLIS ADD (
  PPN_CHECK VARCHAR2(10)
)
/
UPDATE {skema}.T_GENERAL_POLIS SET PPN_CHECK = PPN_CHECK_TEKS WHERE PPN_CHECK_TEKS IS NOT NULL
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_GENERAL_POLIS' AND COLUMN_NAME = 'PPN_CHECK_TEKS';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.T_GENERAL_POLIS DROP COLUMN PPN_CHECK_TEKS';
  END IF;
END;
/
