-- 601 - TPL per interest (Section InputDtlInterest, SetTPLNote_Act) di T_CLAIM_INTEREST (tabel Claim Prop).
-- Kolom nullable (wajib-isi di Go), nol MODIFY / DROP kolom yang sudah ada; uang, persen, kurs NUMBER(38,10).
-- Daftar kolom = katalog `models/katalog.go` (properti Pega di komentar katalog). NOL COMMIT. -migrate oleh work owner.
ALTER TABLE {skema}.T_CLAIM_INTEREST ADD (
  IS_TPL         VARCHAR2(16),
  TPL_FORMAT     VARCHAR2(16),
  TPL_TYPE       VARCHAR2(16),
  TPL_PCT        NUMBER(38,10),
  TPL_AMOUNT     NUMBER(38,10),
  TPL_MIN_MAX    VARCHAR2(16),
  TPL_CURRENCY   VARCHAR2(64),
  TPL_NOTE       VARCHAR2(255),
  TPL_TYPE2      VARCHAR2(16),
  TPL_PCT2       NUMBER(38,10),
  TPL_AMOUNT2    NUMBER(38,10),
  TPL_MIN_MAX2   VARCHAR2(16),
  TPL_CURRENCY2  VARCHAR2(64)
)
/
