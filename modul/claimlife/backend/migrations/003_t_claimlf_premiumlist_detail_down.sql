-- Jalur mundur untuk 003_t_claimlf_premiumlist_detail.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
DROP TABLE {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL CASCADE CONSTRAINTS
/
