-- 935 - ringkasan R/I Risk: view RIRISK_LIFE_SUMMARY dibuang, tabel Pega M_RIRISK_LIFE_SUMMARY diganti nama menjadi
-- RIRISK_LIFE_SUMMARY (langkah 1 dari 3; 936 kolom, 937 isi + buang JSONDATA).
--
-- Keputusan work owner 08-10-2026 K1: nama akhir = RIRISK_LIFE_SUMMARY, tabel FLAT tanpa JSONDATA, BUKAN tabel baru -
-- DROP VIEW (hanya bila objeknya memang VIEW) lalu ALTER TABLE ... RENAME TO (data, PK SYS_C009367, indeks, dan grant
-- ikut). Pembaca lain (masterproductnamelife, BrowseRIRiskSummary) tetap memakai nama yang sama.
--
-- Dua blok berpelindung katalog, SETIAP pernyataan aman DIULANG:
--   1. SYS.ALL_VIEWS (bentuk kedua migrasi.BacaPerintahKatalog) - DROP VIEW hanya selama RIRISK_LIFE_SUMMARY masih
--      VIEW (ALL_TAB_COLUMNS memuat kolom view DAN tabel, jadi tidak dapat membedakannya);
--   2. ALL_TAB_COLUMNS tabel SUMBER - RENAME hanya selama M_RIRISK_LIFE_SUMMARY masih ada. Bila target sudah ada sebagai
--      TABLE (keadaan yang tidak boleh terjadi), RENAME GAGAL KERAS ORA-00955 - tidak dilewati diam-diam.
-- Pra-terbang pelari (praTerbangBentuk/objekAda) hanya membaca CREATE TABLE; langkah ini tanpa CREATE TABLE.
-- ⛔ PRASYARAT P1-P3 (LANGKAH-WO-RIRISKLIFE.md (a)): penulis dihentikan, pemakaian skema GL diperiksa, cadangan CSV.
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_VIEWS
   WHERE OWNER = UPPER('{skema}') AND VIEW_NAME = 'RIRISK_LIFE_SUMMARY';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'DROP VIEW {skema}.RIRISK_LIFE_SUMMARY';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_RIRISK_LIFE_SUMMARY' AND COLUMN_NAME = 'ID';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.M_RIRISK_LIFE_SUMMARY RENAME TO RIRISK_LIFE_SUMMARY';
  END IF;
END;
/
