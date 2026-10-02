-- Jalur mundur untuk 404_retensi_cedant.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
DROP TABLE {skema}.RETENSI_CEDANT CASCADE CONSTRAINTS
/
