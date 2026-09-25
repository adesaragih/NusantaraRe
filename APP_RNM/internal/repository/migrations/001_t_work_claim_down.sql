-- Jalur mundur untuk 001_t_work_claim.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
DROP TABLE {skema}.T_WORK_CLAIM CASCADE CONSTRAINTS
/
