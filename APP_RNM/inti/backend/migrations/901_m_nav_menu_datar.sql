-- 901 - M_NAV_MENU DATAR: satu baris per modul, PARENT_ID dibuang.
--
-- Keputusan work owner 30-09-2026 (`PROMPT-MENU-DATAR-PER-GROUPMENU.md`):
-- "menu jangan ada model seperti child. Buat grouping menu antar GROUPMENU dari
-- tabel M_NAV_MENU. Butir inbox, register, premiumlist, komite, tco-tahun
-- harusnya tidak perlu, karena 1 modul 1 menu."
--
-- Hasil: 20 baris - satu per folder modul korpus, KODE = MODUL = nama modul
-- backend, URUTAN = urutan di dalam GROUPMENU. 900 yang sudah terpasang TIDAK
-- disunting (`T_MIGRASI` mencatat namanya); langkah ini membuang butir anaknya
-- dan tingkat keduanya.
--
-- ⛔ SETIAP LANGKAH DILINDUNGI PEMERIKSAAN KATALOG. Pelari migrasi menjalankan
-- tiap pernyataan tanpa transaksi dan hanya menoleransi ORA-00955 pada CREATE:
-- langkah yang gagal di tengah diulang dari awal. Tanpa pelindung, pengulangan
-- sesudah kolom terbuang mati di ORA-00904 (DELETE ... WHERE PARENT_ID) atau
-- ORA-02443 (DROP CONSTRAINT yang sudah hilang). Karena itu tiap langkah
-- bertanya ke `SYS.ALL_*` lebih dulu, lalu menjalankan perintahnya lewat
-- EXECUTE IMMEDIATE hanya bila masih perlu. Nol COMMIT di teks SQL.
--
-- Backend yang membaca menu (`inti/backend/menu`) menyaring `KODE = MODUL`, jadi
-- ia benar SEBELUM dan SESUDAH langkah ini berjalan.

-- 1. Lima butir anak (inbox, register, premiumlist, komite, tco-tahun).
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_NAV_MENU' AND COLUMN_NAME = 'PARENT_ID';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'DELETE FROM {skema}.M_NAV_MENU WHERE PARENT_ID IS NOT NULL';
  END IF;
END;
/

-- 2. Kunci tamu ke induknya.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_CONSTRAINTS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_NAV_MENU' AND CONSTRAINT_NAME = 'FK_M_NAV_MENU_INDUK';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.M_NAV_MENU DROP CONSTRAINT FK_M_NAV_MENU_INDUK';
  END IF;
END;
/

-- 3. Indeks PARENT_ID (nama dari 900).
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_INDEXES
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_NAV_MENU' AND INDEX_NAME = 'IX_M_NAV_MENU_PARENT';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'DROP INDEX {skema}.IX_M_NAV_MENU_PARENT';
  END IF;
END;
/

-- 4. Kolom PARENT_ID.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_NAV_MENU' AND COLUMN_NAME = 'PARENT_ID';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.M_NAV_MENU DROP COLUMN PARENT_ID';
  END IF;
END;
/
