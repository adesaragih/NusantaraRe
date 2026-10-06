-- Jalur mundur 892 - kolom spread Aviation kembali ke NUMBER(10,4) warisan.
--
-- DAPAT GAGAL (ORA-01440) bila kolomnya sudah berisi nilai: Oracle tidak mengizinkan presisi/skala NUMBER dikecilkan
-- pada kolom berisi. Gagalnya terang, dan itu yang benar - memotong nominal diam-diam jauh lebih buruk.
ALTER TABLE {skema}.BORDEREAUX_PREMI_AVIATION MODIFY (
  SPREAD_OF_RISK_OR     NUMBER(10,4),
  SPREAD_OF_RISK_QS     NUMBER(10,4),
  SPREAD_OF_RISK_SPL    NUMBER(10,4),
  SPREAD_OF_RISK_OTHERS NUMBER(10,4)
)
/
ALTER TABLE {skema}.BORDEREAUX_CLAIM_AVIATION MODIFY (
  SPREAD_OF_CLAIM_OR     NUMBER(10,4),
  SPREAD_OF_CLAIM_QS     NUMBER(10,4),
  SPREAD_OF_CLAIM_OTHERS NUMBER(10,4)
)
/
