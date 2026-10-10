-- 946 - Plan (planlife): view PRODUCT_TYPE_LIFE dibuang, tabel Pega M_PRODUCT_TYPE_LIFE diganti nama menjadi
-- PRODUCT_TYPE_LIFE (langkah 1 dari 3; 947 kolom, 948 isi + pemeriksaan + PK + buang JSONDATA).
--
-- Keputusan work owner 08-10-2026 K1: nama akhir = PRODUCT_TYPE_LIFE, tabel FLAT tanpa JSONDATA, BUKAN tabel baru -
-- DROP VIEW (hanya bila objeknya memang VIEW) lalu ALTER TABLE ... RENAME TO (data dan constraint IS JSON ikut; tabel
-- ini TANPA PK - PK dibuat 948). Pola sama persis dengan Benefit 942 / R/I Risk 935.
--
-- Dua blok berpelindung katalog, SETIAP pernyataan aman DIULANG:
--   1. SYS.ALL_VIEWS - DROP VIEW hanya selama PRODUCT_TYPE_LIFE masih VIEW;
--   2. ALL_TAB_COLUMNS tabel SUMBER - RENAME hanya selama M_PRODUCT_TYPE_LIFE masih ada. "Target belum TABLE" dijaga
--      Oracle sendiri: bila PRODUCT_TYPE_LIFE sudah ada sebagai TABLE (keadaan yang tidak boleh terjadi), RENAME GAGAL
--      KERAS ORA-00955 di dalam blok - pelari hanya menoleransi ORA-00955 pada pernyataan CREATE, jadi langkah berhenti
--      dan tidak tercatat (bentuk penjaga RENAME hanya menanyakan satu katalog; pemeriksaan kedua tidak diperlukan).
-- ⛔ PRASYARAT P1-P3 (modul/planlife/docs/LANGKAH-WO-PLANLIFE.md (a)): penulis dihentikan, acuan dua kali, cadangan
-- ID + JSONDATA, bukti K2. NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_VIEWS
   WHERE OWNER = UPPER('{skema}') AND VIEW_NAME = 'PRODUCT_TYPE_LIFE';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'DROP VIEW {skema}.PRODUCT_TYPE_LIFE';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_PRODUCT_TYPE_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.M_PRODUCT_TYPE_LIFE RENAME TO PRODUCT_TYPE_LIFE';
  END IF;
END;
/
