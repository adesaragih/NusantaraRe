-- 327 - T_POLIS_XOL_LAYER <- TreatyXOLList().ValueList - pemegang penanda
-- layer (ID-29; tiket 19 AC 32-35).
-- NOURUT = nomor urut baris di dalam induknya, unik per induk (ID-11, AC 8, 10);
-- kunci pasangan antar generasi, bukan kunci dagang (ID-13, AC 11).
CREATE TABLE {skema}.T_POLIS_XOL_LAYER (
  ID                        VARCHAR2(32) NOT NULL,
  XOL_ID                    VARCHAR2(32) NOT NULL,
  NOURUT                    NUMBER(5) NOT NULL,
  LAYER                     VARCHAR2(64),
  LAYER_TYPE                VARCHAR2(64),
  LAYER_PART                VARCHAR2(64),
  LAYER_PART_TYPE           VARCHAR2(64),
  CURRENCY                  VARCHAR2(16),
  ID_CURRENCY               VARCHAR2(64),
  GROSS_PREMI               NUMBER(38,10),
  NET_PREMI                 NUMBER(38,10),
  DEDUCTION                 NUMBER(38,10),
  DUE_TO                    VARCHAR2(16),
  DUE_TO_VALUE              NUMBER(38,10),
  BROKERAGE_FEE_SEBENARNYA  NUMBER(38,10),
  PPH_VALUE                 NUMBER(38,10),
  PPN_VALUE                 NUMBER(38,10),
  NET_PREMI_AFTER_PPH       NUMBER(38,10),
  NET_PREMI_AFTER_PPN       NUMBER(38,10),
  NET_PREMI_AFTER_TAX       NUMBER(38,10),
  CONSTRAINT PK_POLIS_XOL_LAYER PRIMARY KEY (ID),
  CONSTRAINT FK_POLIS_XOL_LAYER_XOL FOREIGN KEY (XOL_ID) REFERENCES {skema}.T_POLIS_XOL (ID),
  CONSTRAINT UQ_POLIS_XOL_LAYER_NOURUT UNIQUE (XOL_ID, NOURUT)
)
/
