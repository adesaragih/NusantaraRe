-- Jalur mundur 185 - kolom gabungan dan daftar Ceding Co.
-- ⚠️ MODIFY kembali ke VARCHAR2(500) GAGAL (ORA-01441) bila ada CEDING_CO_NAME tersimpan > 500 bita;
-- kosongkan/potong nilai itu lebih dulu (keputusan manusia - data).
ALTER TABLE {skema}.T_QUOTATIONDATA MODIFY (CEDING_CO_NAME VARCHAR2(500))
/
ALTER TABLE {skema}.T_QUOTATIONDATA DROP COLUMN CEDING_CO
/
DROP SEQUENCE {skema}.SEQ_T_CEDINGCOLIST
/
DROP TABLE {skema}.T_CEDINGCOLIST CASCADE CONSTRAINTS
/
