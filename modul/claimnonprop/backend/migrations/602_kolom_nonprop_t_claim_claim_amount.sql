-- 602 - kolom ListClaimAmount Non Prop (TPL, fee, salvage, proporsi, klaim cedant, CountClaimTNP_Act) di T_CLAIM_CLAIM_AMOUNT.
-- Kolom nullable (wajib-isi di Go), nol MODIFY / DROP kolom yang sudah ada; uang, persen, kurs NUMBER(38,10).
-- Daftar kolom = katalog `models/katalog.go` (properti Pega di komentar katalog). NOL COMMIT. -migrate oleh work owner.
ALTER TABLE {skema}.T_CLAIM_CLAIM_AMOUNT ADD (
  TPL                  NUMBER(38,10),
  ADJUSTER_FEE         NUMBER(38,10),
  SALVAGE              NUMBER(38,10),
  OTHERS_FEE           NUMBER(38,10),
  PROPORTION_PCT       NUMBER(38,10),
  CLAIM_AMOUNT_CEDANT  NUMBER(38,10),
  CLAIM_AMOUNT_ADJUST  NUMBER(38,10),
  TSI_VALUE            NUMBER(38,10),
  IS_LOCKED            VARCHAR2(16)
)
/
