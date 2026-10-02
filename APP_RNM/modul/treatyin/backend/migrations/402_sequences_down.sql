-- Jalur mundur untuk 402_sequences.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
DROP SEQUENCE {skema}.SEQ_TRIN_KONTRAK
/
DROP SEQUENCE {skema}.SEQ_TRIN_VERSI_KONTRAK
/
DROP SEQUENCE {skema}.SEQ_TRIN_MATA_UANG
/
DROP SEQUENCE {skema}.SEQ_TRIN_JENIS_POTONGAN
/
DROP SEQUENCE {skema}.SEQ_TRIN_JENIS_REASURANSI
/
DROP SEQUENCE {skema}.SEQ_TRIN_BAHAYA
/
DROP SEQUENCE {skema}.SEQ_TRIN_KELOMPOK_TREATY
/
DROP SEQUENCE {skema}.SEQ_TRIN_KELAS_BISNIS
/
