-- Jalur mundur untuk 429_penyebaran.sql
--
-- ⚠️ Membongkar ini mengembalikan `BAGIAN` menjadi induk tanpa cabang
-- penyebaran, dan tiket 38 kembali tanpa tabelnya.
DROP TABLE {skema}.NILAI_PENYEBARAN
/
DROP TABLE {skema}.RINCIAN_PENYEBARAN
/
DROP TABLE {skema}.PENYEBARAN
/
DROP SEQUENCE {skema}.SEQ_TRIN_NILAI_PENYEBARAN
/
DROP SEQUENCE {skema}.SEQ_TRIN_RINCIAN_PENYEBARAN
/
DROP SEQUENCE {skema}.SEQ_TRIN_PENYEBARAN
/
