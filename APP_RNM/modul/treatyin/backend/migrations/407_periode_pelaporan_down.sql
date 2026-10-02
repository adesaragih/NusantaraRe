-- Jalur mundur untuk 407_periode_pelaporan.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
DROP TABLE {skema}.PERIODE_PELAPORAN CASCADE CONSTRAINTS
/
