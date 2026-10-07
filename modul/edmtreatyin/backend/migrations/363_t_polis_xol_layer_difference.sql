-- 363 - T_POLIS_XOL_LAYER_DIFFERENCE <- PolicyTreatyIn.TreatyXOLDifferenceList(c).ValueList(l) (NonProporsional saja).
-- Asal: Activity/CalculateDifferenceEDM_act langkah 1.2.1 / 1.2.3 / 1.2.4 / 1.2.5 - kunci lapisan + Currency /
-- IDCurrency disalin, sepuluh medan uang = selisih menurut rumus XML (batas bawah 0, prorata, pembatalan mentah),
-- DueTo dari tanda selisih DueToValue lapisan itu ("DUE TO US" / "DUE TO YOU").
-- Induk TreatyXOLDifferenceList TIDAK bertabel (ID-7): kuncinya salinan kunci lapisan, nilainya jumlah lapisan
-- (langkah 2) - dibangun ulang dari baris ini, dikelompokkan menurut ID_CURRENCY urut NOURUT.
-- NOURUT = posisi lapisan dalam daftar datar (mata uang ke-c lalu lapisan ke-l).
CREATE TABLE {skema}.T_POLIS_XOL_LAYER_DIFFERENCE (
  ID                        VARCHAR2(32) NOT NULL,
  DIFFERENCE_ID             VARCHAR2(32) NOT NULL,
  NOURUT                    NUMBER(5) NOT NULL,
  LAYER                     VARCHAR2(64),
  LAYER_TYPE                VARCHAR2(64),
  LAYER_PART                VARCHAR2(64),
  LAYER_PART_TYPE           VARCHAR2(64),
  CURRENCY                  VARCHAR2(16),
  ID_CURRENCY               VARCHAR2(64),
  DUE_TO                    VARCHAR2(16),
  GROSS_PREMI               NUMBER(38,10),
  DEDUCTION                 NUMBER(38,10),
  NET_PREMI                 NUMBER(38,10),
  DUE_TO_VALUE              NUMBER(38,10),
  BROKERAGE_FEE_SEBENARNYA  NUMBER(38,10),
  PPH_VALUE                 NUMBER(38,10),
  PPN_VALUE                 NUMBER(38,10),
  NET_PREMI_AFTER_PPH       NUMBER(38,10),
  NET_PREMI_AFTER_PPN       NUMBER(38,10),
  NET_PREMI_AFTER_TAX       NUMBER(38,10),
  PASANGAN_BERGESER         VARCHAR2(1),
  CONSTRAINT PK_POLIS_XOL_LAYER_DIFF PRIMARY KEY (ID),
  CONSTRAINT FK_POLIS_XOL_LAYER_DIFF_DIFF FOREIGN KEY (DIFFERENCE_ID) REFERENCES {skema}.T_POLIS_DIFFERENCE (ID),
  CONSTRAINT UQ_POLIS_XOL_LAYER_DIFF_NOURUT UNIQUE (DIFFERENCE_ID, NOURUT)
)
/
