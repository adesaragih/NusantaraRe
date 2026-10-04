-- Jalur mundur untuk 413_layer.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
-- Ketiga tabel anak dibongkar sebelum LAYER.
DROP TABLE {skema}.PEMULIHAN_LIMIT CASCADE CONSTRAINTS
/
DROP TABLE {skema}.NILAI_MDP_MINIMUM CASCADE CONSTRAINTS
/
DROP TABLE {skema}.NILAI_MDP CASCADE CONSTRAINTS
/
DROP TABLE {skema}.LAYER CASCADE CONSTRAINTS
/
