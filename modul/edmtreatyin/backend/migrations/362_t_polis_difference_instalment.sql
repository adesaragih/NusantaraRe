-- 362 - T_POLIS_DIFFERENCE_INSTALMENT <- PolicyTreatyIn.TreatyDifference.ListInstallment (SATU tingkat - ID-42:
-- nol rujukan TreatyDifference.ListInstallment(n).InstallmentList di korpus EDM).
-- Asal: Activity/EDMTCalculateTreatyDifference langkah 3.1 - DueDate, InstallmentNo, InstallmentPercentage
-- disalin; PaymentTotal, Premium = baru(n) - lama(n). NOURUT = posisi baris (ID-34).
CREATE TABLE {skema}.T_POLIS_DIFFERENCE_INSTALMENT (
  ID                      VARCHAR2(32) NOT NULL,
  DIFFERENCE_ID           VARCHAR2(32) NOT NULL,
  NOURUT                  NUMBER(5) NOT NULL,
  INSTALLMENT_NO          NUMBER(10),
  DUE_DATE                DATE,
  INSTALLMENT_PERCENTAGE  NUMBER(38,10),
  PREMIUM                 NUMBER(38,10),
  PAYMENT_TOTAL           NUMBER(38,10),
  PASANGAN_BERGESER       VARCHAR2(1),
  RUMUS_BERLAPIS          VARCHAR2(1),
  CONSTRAINT PK_POLIS_DIFF_INSTALMENT PRIMARY KEY (ID),
  CONSTRAINT FK_POLIS_DIFF_INSTALMENT_DIFF FOREIGN KEY (DIFFERENCE_ID) REFERENCES {skema}.T_POLIS_DIFFERENCE (ID),
  CONSTRAINT UQ_POLIS_DIFF_INST_NOURUT UNIQUE (DIFFERENCE_ID, NOURUT)
)
/
