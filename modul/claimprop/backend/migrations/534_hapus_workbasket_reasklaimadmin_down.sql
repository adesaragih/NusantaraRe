-- Mundur 534 - baris master ReasKlaimAdmin dikembalikan tanpa pemegang (akun pemegangnya tidak dicatat); untuk skema
-- uji.
INSERT INTO {skema}.M_WORKBASKET (WORKBASKET_ID, NAME, IS_ACTIVE)
SELECT 'ReasKlaimAdmin', 'Klaim - Admin', 1 FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_WORKBASKET WHERE WORKBASKET_ID = 'ReasKlaimAdmin')
/
