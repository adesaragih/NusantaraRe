-- Jalur mundur untuk 006_t_claimlf_adjustment_spreading_retro.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
DROP TABLE {skema}.T_CLAIMLF_ADJ_SPREADING_RETRO CASCADE CONSTRAINTS
/
