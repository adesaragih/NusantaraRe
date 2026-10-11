-- Jalur mundur 204 - T_COVERAGELIST.ZONE kembali ke VARCHAR2(1000) (lebar 203).
-- ⚠️ MODIFY kembali GAGAL (ORA-01441) bila ada nilai Zone tersimpan lebih dari 1000 byte - justru itu sebab 204 ditulis;
-- kosongkan/potong nilai itu lebih dulu (keputusan manusia - data).
ALTER TABLE {skema}.T_COVERAGELIST MODIFY (ZONE VARCHAR2(1000))
/
