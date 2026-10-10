-- 090 - Cause Of Loss Life (causeoflosslife): view CAUSEOFLOSS_LIFE dibuang, tabel Pega M_CAUSEOFLOSS_LIFE diganti nama
-- menjadi CAUSEOFLOSS_LIFE (langkah 1 dari 3; 091 kolom, 092 isi + pemeriksaan K1.4 + buang JSONDATA).
--
-- Keputusan work owner 08-10-2026 K1: nama akhir = CAUSEOFLOSS_LIFE, tabel FLAT tanpa JSONDATA, BUKAN tabel baru - DROP
-- VIEW (hanya bila objeknya memang VIEW) lalu ALTER TABLE ... RENAME TO (data, PK SYS_C008825, dan constraint IS JSON
-- ENSURE_M_CAUSEOFLOSS_LIFE_JSON ikut). Pola sama persis dengan Benefit 942.
-- K0: migrasi MODUL (rentang 090-099, dipinjam dari jatah premiumlistlife) - di skema baru berjalan SEBELUM 900; tidak
-- bergantung pada objek inti mana pun (objek Pega sudah ada).
--
-- Dua blok berpelindung katalog, SETIAP pernyataan aman DIULANG:
--   1. SYS.ALL_VIEWS - DROP VIEW hanya selama CAUSEOFLOSS_LIFE masih VIEW (ALL_TAB_COLUMNS memuat kolom view DAN tabel);
--   2. ALL_TAB_COLUMNS tabel SUMBER - RENAME hanya selama M_CAUSEOFLOSS_LIFE masih ada. Bila target sudah ada sebagai
--      TABLE (keadaan yang tidak boleh terjadi), RENAME GAGAL KERAS ORA-00955 - tidak dilewati diam-diam.
-- Pra-terbang pelari hanya membaca CREATE TABLE; langkah ini tanpa CREATE TABLE.
-- ⛔ PRASYARAT P1-P3 (modul/causeoflosslife/docs/LANGKAH-WO-CAUSEOFLOSSLIFE.md (a)): penulis dihentikan, acuan dua
-- kali, cadangan ID + JSONDATA. NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_VIEWS
   WHERE OWNER = UPPER('{skema}') AND VIEW_NAME = 'CAUSEOFLOSS_LIFE';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'DROP VIEW {skema}.CAUSEOFLOSS_LIFE';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_CAUSEOFLOSS_LIFE' AND COLUMN_NAME = 'ID';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.M_CAUSEOFLOSS_LIFE RENAME TO CAUSEOFLOSS_LIFE';
  END IF;
END;
/
