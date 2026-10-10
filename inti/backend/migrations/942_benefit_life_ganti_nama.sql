-- 942 - Benefit (benefitlife): view BENEFIT_LIFE dibuang, tabel Pega M_BENEFIT_LIFE diganti nama menjadi BENEFIT_LIFE
-- (langkah 1 dari 3; 943 kolom, 944 isi + pemeriksaan K2 + buang JSONDATA).
--
-- Keputusan work owner 08-10-2026 K1: nama akhir = BENEFIT_LIFE, tabel FLAT tanpa JSONDATA, BUKAN tabel baru - DROP VIEW
-- (hanya bila objeknya memang VIEW) lalu ALTER TABLE ... RENAME TO (data, PK SYS_C009031, dan constraint IS JSON ikut).
-- Pola sama persis dengan R/I Risk 935/938.
--
-- Dua blok berpelindung katalog, SETIAP pernyataan aman DIULANG:
--   1. SYS.ALL_VIEWS - DROP VIEW hanya selama BENEFIT_LIFE masih VIEW (ALL_TAB_COLUMNS memuat kolom view DAN tabel);
--   2. ALL_TAB_COLUMNS tabel SUMBER - RENAME hanya selama M_BENEFIT_LIFE masih ada. Bila target sudah ada sebagai TABLE
--      (keadaan yang tidak boleh terjadi), RENAME GAGAL KERAS ORA-00955 - tidak dilewati diam-diam.
-- Pra-terbang pelari hanya membaca CREATE TABLE; langkah ini tanpa CREATE TABLE.
-- ⛔ PRASYARAT P1-P3 (modul/benefitlife/docs/LANGKAH-WO-BENEFITLIFE.md (a)): penulis dihentikan, acuan dua kali,
-- cadangan ID + JSONDATA. NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_VIEWS
   WHERE OWNER = UPPER('{skema}') AND VIEW_NAME = 'BENEFIT_LIFE';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'DROP VIEW {skema}.BENEFIT_LIFE';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_BENEFIT_LIFE' AND COLUMN_NAME = 'ID';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.M_BENEFIT_LIFE RENAME TO BENEFIT_LIFE';
  END IF;
END;
/
