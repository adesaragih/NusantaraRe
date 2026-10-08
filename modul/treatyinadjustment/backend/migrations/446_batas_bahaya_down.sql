-- Pembalikan `446_batas_bahaya.sql`.
--
-- ⛔ DATA HILANG, dan itu memang artinya: kesepuluh kolom diisi pemuat dari
-- dokumen, jadi nilainya dapat dilahirkan kembali dengan menjalankan ulang
-- `pemuat -semua -tabel T_TREATY_HAZARD_LIMIT -ikat` (dan `-penyesuaian`).
-- Nol nilai yang hanya hidup di sini.

DROP SEQUENCE {skema}.SEQ_TT_HAZARD_LIMIT
/
DROP TABLE {skema}.T_TREATY_HAZARD_LIMIT PURGE
/
