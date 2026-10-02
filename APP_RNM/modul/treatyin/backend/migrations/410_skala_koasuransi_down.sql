-- Jalur mundur untuk 410_skala_koasuransi.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
DROP TABLE {skema}.SKALA_KOASURANSI CASCADE CONSTRAINTS
/
