-- 523 - T_CLAIM_CLAIM_AMOUNT <- pyWorkPage.ClaimData.ListClaimAmount (Section OutstandingClaim_Est grid Claim Amount).
-- KURS <- .IDR (AddListClaimAmount langkah 6: IDR = KursObjectItem; SetCurencyList_act: IDR = kurs standar).
-- VALUE_IDR <- .USD - nama Pega berbohong, isinya .Value x kurs ke IDR (AC 106).
CREATE TABLE {skema}.T_CLAIM_CLAIM_AMOUNT (
  ID                    VARCHAR2(32) NOT NULL,
  CLAIM_ID              VARCHAR2(32) NOT NULL,
  NOURUT                NUMBER(5) NOT NULL,
  CURRENCY_ID           VARCHAR2(64),
  CURRENCY_NAME         VARCHAR2(64),
  KURS                  NUMBER(38,10),
  CLAIM_AMOUNT          NUMBER(38,10),
  NET_DEDUCTIBLE_VALUE  NUMBER(38,10),
  VALUE                 NUMBER(38,10),
  VALUE_IDR             NUMBER(38,10),
  NOTE                  VARCHAR2(16),
  CONSTRAINT PK_CLAIM_CLAIM_AMOUNT PRIMARY KEY (ID),
  CONSTRAINT FK_CLAIM_CLAIM_AMOUNT_CLAIM FOREIGN KEY (CLAIM_ID) REFERENCES {skema}.T_GENERAL_CLAIM (ID) ON DELETE CASCADE,
  CONSTRAINT UQ_CLAIM_CLAIM_AMOUNT_NOURUT UNIQUE (CLAIM_ID, NOURUT)
)
/
