-- Jalur mundur untuk 412_dokumen_kontrak.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
DROP TABLE {skema}.DOKUMEN_KONTRAK CASCADE CONSTRAINTS
/
