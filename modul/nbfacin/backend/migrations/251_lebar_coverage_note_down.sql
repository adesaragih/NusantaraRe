-- Jalur mundur 251 - T_COVERAGELIST.COVERAGE_NOTE kembali VARCHAR2(500). ⛔ GAGAL (ORA-01441) bila sudah ada baris dengan
-- COVERAGE_NOTE lebih dari 500 byte - periksa dulu `SELECT COUNT(*) FROM T_COVERAGELIST WHERE LENGTHB(COVERAGE_NOTE) > 500` dan
-- putuskan isinya sebelum mundur. Biner yang memuat 251 (validasi 4000) tidak boleh berjalan di atas skema yang sudah dimundurkan.
ALTER TABLE {skema}.T_COVERAGELIST MODIFY (COVERAGE_NOTE VARCHAR2(500))
/
