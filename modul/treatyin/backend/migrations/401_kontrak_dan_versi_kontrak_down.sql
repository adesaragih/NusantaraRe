-- Jalur mundur untuk 401_kontrak_dan_versi_kontrak.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya. VERSI_KONTRAK menunjuk KONTRAK, jadi ia
-- dibongkar lebih dulu.
DROP TABLE {skema}.VERSI_KONTRAK CASCADE CONSTRAINTS
/
DROP TABLE {skema}.KONTRAK CASCADE CONSTRAINTS
/
