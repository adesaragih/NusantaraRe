-- 361 - T_POLIS_DIFFERENCE_SPREADING <- PolicyTreatyIn.TreatyDifference.SpreadingRiskList.
-- Asal: Activity/EDMTCalculateTreatyDifference langkah 2.1 - TreatyType, TreatyName, SharePercentage,
-- ClaimPercentage disalin dari baris baru; PremiumSpreaded, ClaimSpreaded = baru(n) - lama(n).
-- NOURUT = posisi baris (.pxListSubscript): kunci pasangan antar generasi, bukan kunci dagang (ID-34).
-- PASANGAN_BERGESER / RUMUS_BERLAPIS - penanda migrasi, hanya terisi pada baris induk SUMBER = 'PEGA' (ID-35..ID-37).
CREATE TABLE {skema}.T_POLIS_DIFFERENCE_SPREADING (
  ID                 VARCHAR2(32) NOT NULL,
  DIFFERENCE_ID      VARCHAR2(32) NOT NULL,
  NOURUT             NUMBER(5) NOT NULL,
  TREATY_TYPE        VARCHAR2(64),
  TREATY_NAME        VARCHAR2(255),
  SHARE_PERCENTAGE   NUMBER(38,10),
  CLAIM_PERCENTAGE   NUMBER(38,10),
  PREMIUM_SPREADED   NUMBER(38,10),
  CLAIM_SPREADED     NUMBER(38,10),
  PASANGAN_BERGESER  VARCHAR2(1),
  RUMUS_BERLAPIS     VARCHAR2(1),
  CONSTRAINT PK_POLIS_DIFF_SPREADING PRIMARY KEY (ID),
  CONSTRAINT FK_POLIS_DIFF_SPREADING_DIFF FOREIGN KEY (DIFFERENCE_ID) REFERENCES {skema}.T_POLIS_DIFFERENCE (ID),
  CONSTRAINT UQ_POLIS_DIFF_SPREAD_NOURUT UNIQUE (DIFFERENCE_ID, NOURUT)
)
/
