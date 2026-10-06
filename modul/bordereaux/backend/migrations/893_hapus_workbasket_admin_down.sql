-- Mundur 893 - baris ReasBordereauxAdmin dikembalikan (tanpa pemegang); untuk skema uji.
INSERT INTO {skema}.M_WORKBASKET (WORKBASKET_ID, NAME, IS_ACTIVE)
SELECT 'ReasBordereauxAdmin', 'Bordereaux Admin', 1 FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_WORKBASKET WHERE WORKBASKET_ID = 'ReasBordereauxAdmin')
/
