-- 603 - fee dan nilai IDR Spreading List Non Prop (CountLossAllocation_act langkah 21) di T_CLAIM_SPREADING.
-- Kolom nullable (wajib-isi di Go), nol MODIFY / DROP kolom yang sudah ada; uang, persen, kurs NUMBER(38,10).
-- Daftar kolom = katalog `models/katalog.go` (properti Pega di komentar katalog). NOL COMMIT. -migrate oleh work owner.
ALTER TABLE {skema}.T_CLAIM_SPREADING ADD (
  ADJUSTER_FEE      NUMBER(38,10),
  SALVAGE           NUMBER(38,10),
  OTHERS_FEE        NUMBER(38,10),
  CLAIM_AMOUNT_IDR  NUMBER(38,10),
  IS_LOCKED         VARCHAR2(16)
)
/
