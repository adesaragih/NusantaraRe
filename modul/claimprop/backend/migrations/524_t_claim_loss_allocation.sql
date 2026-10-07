-- 524 - T_CLAIM_LOSS_ALLOCATION <- pyWorkPage.ClaimData.SpreadingRisk (Section OutstandingClaim_Est grid "Loss Allocation").
-- RALAT 07-10-2026 atas STRUKTUR J1 butir 1 ("SpreadingRisk = bacaan master treaty"): di Claim Prop SpreadingRisk
-- adalah DATA KLAIM - baris lahir di AddLossAllocation_act (Page-New/append), diubah CountPersen_act, dihapus
-- RemoveLossAlloction_act (Page-Remove), dibekukan SaveOutstanding_Act langkah 25 (IsOldData); nol RDB-List
-- mengisinya di 112 activity. Karena itu tabel kesebelas, sejajar T_CLAIM_SPREADING. Nama menunggu work owner (OQ).
-- KURS <- .PremiumSpreaded (SetCurrency_Act langkah 5: kurs standar mata uang baris).
CREATE TABLE {skema}.T_CLAIM_LOSS_ALLOCATION (
  ID                VARCHAR2(32) NOT NULL,
  CLAIM_ID          VARCHAR2(32) NOT NULL,
  NOURUT            NUMBER(5) NOT NULL,
  CURRENCY_ID       VARCHAR2(64),
  CURRENCY_NAME     VARCHAR2(64),
  TREATY_TYPE_ID    VARCHAR2(64),
  TREATY_NAME       VARCHAR2(255),
  SHARE_PERCENTAGE  NUMBER(38,10),
  CLAIM_SPREADED    NUMBER(38,10),
  CLAIM_ESTIMATION  NUMBER(38,10),
  KURS              NUMBER(38,10),
  IS_OLD_DATA       VARCHAR2(16),
  CONSTRAINT PK_CLAIM_LOSS_ALLOCATION PRIMARY KEY (ID),
  CONSTRAINT FK_CLAIM_LOSS_ALLOC_CLAIM FOREIGN KEY (CLAIM_ID) REFERENCES {skema}.T_GENERAL_CLAIM (ID) ON DELETE CASCADE,
  CONSTRAINT UQ_CLAIM_LOSS_ALLOC_NOURUT UNIQUE (CLAIM_ID, NOURUT)
)
/
