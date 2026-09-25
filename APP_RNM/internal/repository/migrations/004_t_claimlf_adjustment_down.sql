-- Jalur mundur untuk 004_t_claimlf_adjustment.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
DROP TABLE {skema}.T_CLAIMLF_ADJUSTMENT CASCADE CONSTRAINTS
/
