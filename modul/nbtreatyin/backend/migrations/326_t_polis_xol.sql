-- 326 - T_POLIS_XOL <- PolicyTreatyIn.TreatyXOLList - induk XOL tanpa penanda
-- layer (ID-29); DEDUCTION di sini UANG (ID-30).
-- NOURUT = nomor urut baris di dalam induknya, unik per induk (ID-11, AC 8, 10);
-- kunci pasangan antar generasi, bukan kunci dagang (ID-13, AC 11).
CREATE TABLE {skema}.T_POLIS_XOL (
  ID                        VARCHAR2(32) NOT NULL,
  POLIS_ID                  VARCHAR2(32) NOT NULL,
  NOURUT                    NUMBER(5) NOT NULL,
  CURRENCY                  VARCHAR2(16),
  ID_CURRENCY               VARCHAR2(64),
  GROSS_PREMI               NUMBER(38,8),
  NET_PREMI                 NUMBER(38,8),
  DEDUCTION                 NUMBER(38,8),
  DUE_TO                    VARCHAR2(16),
  DUE_TO_VALUE              NUMBER(38,8),
  BROKERAGE_FEE_SEBENARNYA  NUMBER(38,8),
  PPH_VALUE                 NUMBER(38,8),
  PPN_VALUE                 NUMBER(38,8),
  NET_PREMI_AFTER_PPH       NUMBER(38,8),
  NET_PREMI_AFTER_PPN       NUMBER(38,8),
  NET_PREMI_AFTER_TAX       NUMBER(38,8),
  CONSTRAINT PK_POLIS_XOL PRIMARY KEY (ID),
  CONSTRAINT FK_POLIS_XOL_POLIS FOREIGN KEY (POLIS_ID) REFERENCES {skema}.T_GENERAL_POLIS (ID),
  CONSTRAINT UQ_POLIS_XOL_NOURUT UNIQUE (POLIS_ID, NOURUT)
)
/
