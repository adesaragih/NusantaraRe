-- Jalur mundur untuk 416_bagian.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
-- Kedua tabel anak paket uang dibongkar sebelum BAGIAN.
DROP TABLE {skema}.NILAI_PREMI_BRUTO_MINIMUM CASCADE CONSTRAINTS
/
DROP TABLE {skema}.NILAI_PREMI_BRUTO CASCADE CONSTRAINTS
/
DROP TABLE {skema}.BAGIAN CASCADE CONSTRAINTS
/
