-- Jalur mundur untuk 417_detail_proporsional.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
-- Tabel anak paket uang dibongkar sebelum induknya.
DROP TABLE {skema}.NILAI_CADANGAN_PREMI CASCADE CONSTRAINTS
/
DROP TABLE {skema}.DETAIL_PROPORSIONAL CASCADE CONSTRAINTS
/
