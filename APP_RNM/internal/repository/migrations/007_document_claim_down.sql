-- Jalur mundur untuk 007_document_claim.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
DROP TABLE {skema}.DOCUMENT_CLAIM CASCADE CONSTRAINTS
/
