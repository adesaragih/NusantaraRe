-- Jalur mundur 018 - butir bd.
--
-- ⚠️ `DISEASE` dikembalikan ke VARCHAR2(255), dan itu DAPAT GAGAL bila sudah
-- ada nama penyakit lebih panjang (ORA-01441). Gagalnya terang, dan itu yang
-- benar: memotong nama penyakit diam-diam jauh lebih buruk.
ALTER TABLE {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL MODIFY (
  DISEASE VARCHAR2(255)
)
/

DROP SEQUENCE {skema}.SEQ_CLAIMLF_DIAGNOSE
/

DROP TABLE {skema}.T_CLAIMLF_DIAGNOSE CASCADE CONSTRAINTS
/
