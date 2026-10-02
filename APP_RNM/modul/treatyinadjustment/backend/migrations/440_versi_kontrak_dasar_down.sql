-- Jalur mundur untuk 440_versi_kontrak_dasar.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya. Index dan kunci asing dibongkar sebelum
-- kolomnya - membuang kolom yang masih dirujuk constraint gagal di Oracle.
DROP INDEX {skema}.IX_VERSI_KONTRAK_DASAR
/
ALTER TABLE {skema}.VERSI_KONTRAK DROP CONSTRAINT FK_VERSI_KONTRAK_DASAR
/
ALTER TABLE {skema}.VERSI_KONTRAK DROP COLUMN ID_VERSI_KONTRAK_DASAR
/
