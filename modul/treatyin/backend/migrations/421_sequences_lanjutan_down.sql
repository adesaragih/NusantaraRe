-- Jalur mundur untuk 421_sequences_lanjutan.sql
--
-- ⚠️ Membongkar sequence MENGHILANGKAN nilai terakhirnya. Memasangnya kembali
-- memulai dari 1, dan pengenal yang sudah terpakai akan ditawarkan ulang -
-- pelanggaran INV-03. Jalankan hanya pada skema yang tabelnya ikut dibongkar.
DROP SEQUENCE {skema}.SEQ_TRIN_ARSIP_MUATAN_KELUAR
/
DROP SEQUENCE {skema}.SEQ_TRIN_RINCIAN_ANGSURAN
/
DROP SEQUENCE {skema}.SEQ_TRIN_PENCAPAIAN
/
DROP SEQUENCE {skema}.SEQ_TRIN_KELAS_BISNIS_KLP
/
DROP SEQUENCE {skema}.SEQ_TRIN_KELOMPOK_LAYER
/
DROP SEQUENCE {skema}.SEQ_TRIN_KELAS_BISNIS_LAYER
/
