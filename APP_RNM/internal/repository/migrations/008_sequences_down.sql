-- Jalur mundur untuk 008_sequences.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
DROP SEQUENCE {skema}.SEQ_CLAIMLF_PLD
/
DROP SEQUENCE {skema}.SEQ_CLAIMLF_ADJ
/
DROP SEQUENCE {skema}.SEQ_CLAIMLF_SPR
/
DROP SEQUENCE {skema}.SEQ_CLAIMLF_SPR_RETRO
/
DROP SEQUENCE {skema}.SEQ_T_CLAIMLF_DOCUMENT
/
