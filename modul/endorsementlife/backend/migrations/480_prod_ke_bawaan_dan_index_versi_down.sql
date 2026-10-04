-- Jalur mundur 480 - tiket 00 Endorsement Life.
-- Isi mundur `PROD_KE = 1` tidak dibalik: nilai itu berarti versi new business.
DROP INDEX {skema}.IDX_PLD_PL_NUMBER
/
DROP INDEX {skema}.IDX_PL_OLD_POLICY_NO
/
DROP INDEX {skema}.IDX_PL_NOPOLIS_PRODKE
/
ALTER TABLE {skema}.T_PREMIUM_LIST MODIFY (
  PROD_KE NUMBER(5) DEFAULT NULL
)
/
