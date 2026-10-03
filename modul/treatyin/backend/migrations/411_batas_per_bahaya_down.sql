-- Jalur mundur untuk 411_batas_per_bahaya.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
DROP TABLE {skema}.BATAS_PER_BAHAYA CASCADE CONSTRAINTS
/
