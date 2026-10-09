-- 610 - T_CLAIM_NP_CLAIM_ACCEPT <- `AdjustmentList(n).ListClaimAcceptation` (Claim Acceptation per mata uang,
-- AddAkseptasiCNP_Act langkah 16; Section AdjustmentDetailNP grid pertama).
-- Tabel baru khas XoL (OQ-CNP-06, docs/STRUKTUR-TABEL-CLAIM-NON-PROP.md). ID dari SEQ_T_CLAIM (claimprop 533).
-- Kolom = katalog `models/katalog.go`. NOL COMMIT. -migrate oleh work owner.
CREATE TABLE {skema}.T_CLAIM_NP_CLAIM_ACCEPT (
  ID                    VARCHAR2(32) NOT NULL,
  ADJUSTMENT_ID         VARCHAR2(32) NOT NULL,
  NOURUT                NUMBER(5) NOT NULL,
  CURRENCY_ID           VARCHAR2(64),
  CURRENCY_NAME         VARCHAR2(64),
  KURS                  NUMBER(38,10),
  VALUE                 NUMBER(38,10),
  VALUE_IDR             NUMBER(38,10),
  NET_DEDUCTIBLE_VALUE  NUMBER(38,10),
  NOTE                  VARCHAR2(16),
  TPL                   NUMBER(38,10),
  ADJUSTER_FEE          NUMBER(38,10),
  SALVAGE               NUMBER(38,10),
  OTHERS_FEE            NUMBER(38,10),
  PROPORTION_PCT        NUMBER(38,10),
  CLAIM_AMOUNT_CEDANT   NUMBER(38,10),
  CLAIM_AMOUNT_ADJUST   NUMBER(38,10),
  TSI_VALUE             NUMBER(38,10),
  IS_LOCKED             VARCHAR2(16),
  CONSTRAINT PK_CLAIM_NP_ACC PRIMARY KEY (ID),
  CONSTRAINT FK_CLAIM_NP_ACC_ADJ FOREIGN KEY (ADJUSTMENT_ID) REFERENCES {skema}.T_CLAIM_ADJUSTMENT (ID) ON DELETE CASCADE,
  CONSTRAINT UQ_CLAIM_NP_ACC_NOURUT UNIQUE (ADJUSTMENT_ID, NOURUT)
)
/
