-- Jalur mundur untuk 418_potongan.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
DROP TABLE {skema}.POTONGAN CASCADE CONSTRAINTS
/
