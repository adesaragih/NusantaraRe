-- Jalur mundur untuk 414_jejak_perubahan.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
DROP TABLE {skema}.JEJAK_PERUBAHAN CASCADE CONSTRAINTS
/
