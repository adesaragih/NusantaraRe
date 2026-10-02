-- Jalur mundur untuk 400_tabel_acuan.sql
--
-- Dijalankan MENURUN: langkah bernomor besar dibongkar lebih dulu, supaya
-- anak hilang sebelum induknya.
--
-- CASCADE CONSTRAINTS membuang kunci asing yang menunjuk tabel ini - termasuk
-- FK_JENIS_REASURANSI_INDUK yang menunjuk tabelnya sendiri.
DROP TABLE {skema}.JENIS_REASURANSI CASCADE CONSTRAINTS
/
DROP TABLE {skema}.BAHAYA CASCADE CONSTRAINTS
/
DROP TABLE {skema}.KELOMPOK_TREATY CASCADE CONSTRAINTS
/
DROP TABLE {skema}.KELAS_BISNIS CASCADE CONSTRAINTS
/
DROP TABLE {skema}.JENIS_POTONGAN CASCADE CONSTRAINTS
/
DROP TABLE {skema}.MATA_UANG CASCADE CONSTRAINTS
/
