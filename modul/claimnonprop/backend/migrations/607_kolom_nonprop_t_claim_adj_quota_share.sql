-- 607 - kolom Spreading Out akseptasi Non Prop (AdjClaimCNP_Act langkah 15) di T_CLAIM_ADJ_QUOTA_SHARE.
-- Kolom nullable (wajib-isi di Go), nol MODIFY / DROP kolom yang sudah ada; uang, persen, kurs NUMBER(38,10).
-- Daftar kolom = katalog `models/katalog.go` (properti Pega di komentar katalog). NOL COMMIT. -migrate oleh work owner.
ALTER TABLE {skema}.T_CLAIM_ADJ_QUOTA_SHARE ADD (
  PREMIUM_SPREADED  NUMBER(38,10),
  ADJUSTER_FEE      NUMBER(38,10),
  SALVAGE           NUMBER(38,10),
  OTHERS_FEE        NUMBER(38,10),
  TOTAL_CLAIM       NUMBER(38,10),
  NET_CLAIM         NUMBER(38,10)
)
/
