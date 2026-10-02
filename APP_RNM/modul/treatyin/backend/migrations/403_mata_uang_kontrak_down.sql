-- Jalur mundur untuk 403_mata_uang_kontrak.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
DROP TABLE {skema}.MATA_UANG_KONTRAK CASCADE CONSTRAINTS
/
