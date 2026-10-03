-- 325 - T_POLIS_SPREADING <- PolicyTreatyIn.SpreadingRiskList (ID-28).
-- Total share dan nilai TIDAK disimpan - turunan baris ini.
-- NOURUT = nomor urut baris di dalam induknya, unik per induk (ID-11, AC 8, 10);
-- kunci pasangan antar generasi, bukan kunci dagang (ID-13, AC 11).
CREATE TABLE {skema}.T_POLIS_SPREADING (
  ID                   VARCHAR2(32) NOT NULL,
  POLIS_ID             VARCHAR2(32) NOT NULL,
  NOURUT               NUMBER(5) NOT NULL,
  TREATY_TYPE          VARCHAR2(64),
  TREATY_NAME          VARCHAR2(255),
  CURRENCY             VARCHAR2(16),
  CURRENCY_ID          VARCHAR2(64),
  SHARE_PERCENTAGE     NUMBER(38,8),
  SPLIT_RNM_SHARE_PCT  NUMBER(38,8),
  CLAIM_PERCENTAGE     NUMBER(38,8),
  PREMIUM_SPREADED     NUMBER(38,8),
  CLAIM_SPREADED       NUMBER(38,8),
  CONSTRAINT PK_POLIS_SPREADING PRIMARY KEY (ID),
  CONSTRAINT FK_POLIS_SPREADING_POLIS FOREIGN KEY (POLIS_ID) REFERENCES {skema}.T_GENERAL_POLIS (ID),
  CONSTRAINT UQ_POLIS_SPREADING_NOURUT UNIQUE (POLIS_ID, NOURUT)
)
/
