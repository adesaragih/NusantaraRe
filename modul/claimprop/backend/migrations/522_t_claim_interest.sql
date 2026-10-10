-- 522 - T_CLAIM_INTEREST <- pyWorkPage.ClaimData.InterestList (Section OutstandingClaim_Intrs grid "Insured Interests 100 %").
-- Objek pertanggungan (AC 41: bukan bunga). Mata uang dan kurs per objek (AC 27).
-- InterestListDtl (salinan tampilan) dan TotalInterestInsured (penjumlahan) tidak disimpan.
CREATE TABLE {skema}.T_CLAIM_INTEREST (
  ID             VARCHAR2(32) NOT NULL,
  CLAIM_ID       VARCHAR2(32) NOT NULL,
  NOURUT         NUMBER(5) NOT NULL,
  OBJECT_NAME    VARCHAR2(1000),
  CURRENCY_ID    VARCHAR2(64),
  CURRENCY_NAME  VARCHAR2(64),
  KURS           NUMBER(38,10),
  TSI_VALUE      NUMBER(38,10),
  TSI_VALUE_IDR  NUMBER(38,10),
  IS_ADJ_VALUE   VARCHAR2(16),
  CONSTRAINT PK_CLAIM_INTEREST PRIMARY KEY (ID),
  CONSTRAINT FK_CLAIM_INTEREST_CLAIM FOREIGN KEY (CLAIM_ID) REFERENCES {skema}.T_GENERAL_CLAIM (ID) ON DELETE CASCADE,
  CONSTRAINT UQ_CLAIM_INTEREST_NOURUT UNIQUE (CLAIM_ID, NOURUT)
)
/
