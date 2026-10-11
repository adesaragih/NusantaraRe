-- 600 - kolom khas lini NONPROP di T_GENERAL_CLAIM (header klaim, tabel bersama milik Claim Life).
--
-- Keputusan work owner 09-10-2026 (OQ-CNP-06 "pola Claim Prop + tabel XoL"; izin menyunting STRUKTUR Claim Life):
-- medan skalar Claim Non Prop (Section OutstandingClaim / InputAcceptation, pyWorkPage, ReceiverClaim bank kedua)
-- disimpan sebagai kolom TAMBAHAN; kolom yang maknanya sama dengan kolom Claim Life / Claim Prop dipakai ulang.
-- Kolom nullable (wajib-isi di Go), nol MODIFY / DROP kolom yang sudah ada; uang, persen, kurs NUMBER(38,10).
-- Daftar kolom = katalog `models/katalog.go` (properti Pega di komentar katalog). NOL COMMIT. -migrate oleh work owner.
ALTER TABLE {skema}.T_GENERAL_CLAIM ADD (
  BUSINESS_OLD_ID         VARCHAR2(64),
  MASTER_ID_TO            VARCHAR2(64),
  POLICY_COB              VARCHAR2(255),
  REINSURANCE_SLIP        VARCHAR2(255),
  CLAIM_NO_CEDING         VARCHAR2(255),
  CIRCUMSTANCES           VARCHAR2(4000),
  SUPPORTING_DOCUMENT     VARCHAR2(4000),
  DEDUCTIBLE_MIN_MAX      VARCHAR2(16),
  IS_TPL                  VARCHAR2(16),
  TPL_FORMAT              VARCHAR2(16),
  TPL_TYPE                VARCHAR2(16),
  TPL_PCT                 NUMBER(38,10),
  IS_SAVE_TO_OS           VARCHAR2(16),
  WAITING_ACTUAL_PREMIUM  VARCHAR2(16),
  STATUS_CASE             VARCHAR2(64),
  IS_REJECT               VARCHAR2(16),
  KOMITE_NO               VARCHAR2(64),
  FLAG_PRINT_PLA          VARCHAR2(16),
  RECEIVER_CURRENCY       VARCHAR2(64),
  RECEIVER_BANK_NAME2     VARCHAR2(500),
  RECEIVER_BANK_BRANCH2   VARCHAR2(500),
  RECEIVER_ACCOUNT_NO2    VARCHAR2(100),
  RECEIVER_SWIFT_CODE2    VARCHAR2(64),
  RECEIVER_BANK_ID2       VARCHAR2(64),
  RECEIVER_CURRENCY2      VARCHAR2(64)
)
/
