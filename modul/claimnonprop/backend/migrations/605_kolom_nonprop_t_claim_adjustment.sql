-- 605 - kolom akseptasi Non Prop (interim, rekening kedua, galat kasir; AdjustmentDetailNP) di T_CLAIM_ADJUSTMENT.
-- Kolom nullable (wajib-isi di Go), nol MODIFY / DROP kolom yang sudah ada; uang, persen, kurs NUMBER(38,10).
-- Daftar kolom = katalog `models/katalog.go` (properti Pega di komentar katalog). NOL COMMIT. -migrate oleh work owner.
ALTER TABLE {skema}.T_CLAIM_ADJUSTMENT ADD (
  INTERIM_INDEX     VARCHAR2(16),
  BANK_NAME2        VARCHAR2(500),
  BANK_BRANCH2      VARCHAR2(500),
  BANK_ACCOUNT_NO2  VARCHAR2(100),
  SWIFT_CODE2       VARCHAR2(64),
  BANK_ID2          VARCHAR2(64),
  BANK_CURRENCY2    VARCHAR2(64),
  FLAG_ERROR_KASIR  VARCHAR2(16)
)
/
