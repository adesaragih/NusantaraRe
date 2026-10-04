-- Jalur mundur untuk 408_periode_akumulasi.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
DROP TABLE {skema}.PERIODE_AKUMULASI CASCADE CONSTRAINTS
/
