-- Jalur mundur 254 - membuang T_WORK_POLIS.POSITION_NOTE. ⚠️ Jalankan SESUDAH biner kembali ke versi tanpa tiket 77, atau sesudah
-- POSITION case NB dikembalikan ke kosakata lama (posisi efektif): sesudah konversi, POSITION berisi tahap, bukan workbasket.
ALTER TABLE {skema}.T_WORK_POLIS DROP (POSITION_NOTE)
/
