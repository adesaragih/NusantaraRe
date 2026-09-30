-- Jalur mundur untuk 002_t_general_claim.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
DROP TABLE {skema}.T_GENERAL_CLAIM CASCADE CONSTRAINTS
/
