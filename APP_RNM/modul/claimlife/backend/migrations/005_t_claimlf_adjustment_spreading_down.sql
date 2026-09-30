-- Jalur mundur untuk 005_t_claimlf_adjustment_spreading.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
DROP TABLE {skema}.T_CLAIMLF_ADJUSTMENT_SPREADING CASCADE CONSTRAINTS
/
