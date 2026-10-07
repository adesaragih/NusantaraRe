-- 360 - T_POLIS_DIFFERENCE: selisih satu generasi endorsemen terhadap generasi TEPAT sebelumnya
-- (PolicyTreatyIn.TreatyDifference). Tabel PROYEKSI, bukan tabel sumber (spec-penyimpanan ID-6, ID-23..ID-27):
-- hanya ditulis Go di transaksi yang sama dengan generasinya; bangun ulang hanya menyentuh SUMBER = 'GO'.
--
-- Asal medan: Activity/EDMTCalculateTreatyDifference langkah 1 (korpus EDM Treaty In) - uang = baru - lama,
-- persen dan Installment disalin. Total* (TotalPremium, TotalClaim, TotalSharePercentage*) TIDAK disimpan:
-- turunan baris T_POLIS_DIFFERENCE_SPREADING (pola NB T_POLIS_SPREADING). Kunci saring NOPOLIS, PRODKE, EDM_NO,
-- IDPEGA disalin dari generasinya supaya pembaca SQL tidak perlu join balik (diagram sheet EDM Treaty In Prop J101).
-- Uang/persen NUMBER(38,10) = tabel dasar NB 320-327 (RALAT NUMBER(38,8) spec, keputusan WO 04-10-2026).
CREATE TABLE {skema}.T_POLIS_DIFFERENCE (
  ID                   VARCHAR2(32) NOT NULL,
  POLIS_ID             VARCHAR2(32) NOT NULL,
  NOPOLIS              VARCHAR2(64),
  PRODKE               NUMBER(10),
  EDM_NO               VARCHAR2(64),
  IDPEGA               VARCHAR2(64),
  SUMBER               VARCHAR2(8) NOT NULL,
  INSTALLMENT          VARCHAR2(16),
  GROSS_PREMIUM        NUMBER(38,10),
  PREMI_OGP            NUMBER(38,10),
  RI_COMM_OGP          NUMBER(38,10),
  RESULT_OGP1          NUMBER(38,10),
  OVERIDDING_COMM_OGP  NUMBER(38,10),
  RESULT_OGP2          NUMBER(38,10),
  CLAIM                NUMBER(38,10),
  SALVAGE_VALUE        NUMBER(38,10),
  EXCESS_LOSS          NUMBER(38,10),
  NET_PREMIUM          NUMBER(38,10),
  BALANCE_DUE_TO       NUMBER(38,10),
  PREMI_ONP            NUMBER(38,10),
  RI_COMM_ONP          NUMBER(38,10),
  RESULT_ONP1          NUMBER(38,10),
  OVERIDDING_COMM_ONP  NUMBER(38,10),
  RESULT_ONP2          NUMBER(38,10),
  DEDUCTION1           NUMBER(38,10),
  DEDUCTION2           NUMBER(38,10),
  PPN_VALUE            NUMBER(38,10),
  PPH_VALUE            NUMBER(38,10),
  BALANCE_BEFORE_TAX   NUMBER(38,10),
  BALANCE_BEFORE_PPH   NUMBER(38,10),
  CONSTRAINT PK_POLIS_DIFFERENCE PRIMARY KEY (ID),
  CONSTRAINT FK_POLIS_DIFFERENCE_POLIS FOREIGN KEY (POLIS_ID) REFERENCES {skema}.T_GENERAL_POLIS_TREATY (ID),
  CONSTRAINT UQ_POLIS_DIFFERENCE_POLIS UNIQUE (POLIS_ID),
  CONSTRAINT CK_POLIS_DIFFERENCE_SUMBER CHECK (SUMBER IN ('PEGA', 'GO'))
)
/
