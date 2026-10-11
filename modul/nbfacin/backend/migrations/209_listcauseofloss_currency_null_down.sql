-- Jalur mundur 209 - T_LISTCAUSEOFLOSS.CURRENCY kembali NOT NULL.
--
-- ⚠️ Sejak 209, baris tanpa mata uang tersimpan NULL; NOT NULL ditolak Oracle (ORA-02296) selama baris itu ada. Jalur mundur
-- karena itu MENGISI NULL dengan 'UNKNOWN' (= DEFAULT kolom, rancangan K-069) sebelum mengembalikan NOT NULL - satu-satunya
-- tempat 'UNKNOWN' ditulis, dan hanya oleh jalur mundur (K-012 berlaku untuk aplikasi). Periksa jumlahnya lebih dulu:
-- `SELECT COUNT(*) FROM {skema}.T_LISTCAUSEOFLOSS WHERE CURRENCY IS NULL`.
UPDATE {skema}.T_LISTCAUSEOFLOSS SET CURRENCY = 'UNKNOWN' WHERE CURRENCY IS NULL
/
ALTER TABLE {skema}.T_LISTCAUSEOFLOSS MODIFY (
  CURRENCY VARCHAR2(50) NOT NULL
)
/
