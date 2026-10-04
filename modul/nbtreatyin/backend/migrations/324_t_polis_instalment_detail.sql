-- 324 - T_POLIS_INSTALMENT_DETAIL <- ListInstallment().InstallmentList - hanya
-- non-proporsional (rancangan-tabel-datar §3.2; tiket 19 AC 31).
-- NOURUT = nomor urut baris di dalam induknya, unik per induk (ID-11, AC 8, 10);
-- kunci pasangan antar generasi, bukan kunci dagang (ID-13, AC 11).
CREATE TABLE {skema}.T_POLIS_INSTALMENT_DETAIL (
  ID                      VARCHAR2(32) NOT NULL,
  INSTALMENT_ID           VARCHAR2(32) NOT NULL,
  NOURUT                  NUMBER(5) NOT NULL,
  INSTALLMENT_NO          NUMBER(10),
  DUE_DATE                DATE,
  PAYMENT_DATE            DATE,
  INSTALLMENT_PERCENTAGE  NUMBER(38,10),
  PREMIUM                 NUMBER(38,10),
  PAYMENT_TOTAL           NUMBER(38,10),
  PREMIUM_AFTER_PPH       NUMBER(38,10),
  PREMIUM_AFTER_PPN       NUMBER(38,10),
  PREMIUM_AFTER_TAX       NUMBER(38,10),
  CURRENCY                VARCHAR2(16),
  ID_CURRENCY             VARCHAR2(64),
  CONSTRAINT PK_POLIS_INSTALMENT_DETAIL PRIMARY KEY (ID),
  CONSTRAINT FK_POLIS_INST_DETAIL_INST FOREIGN KEY (INSTALMENT_ID) REFERENCES {skema}.T_POLIS_INSTALMENT (ID),
  CONSTRAINT UQ_POLIS_INST_DETAIL_NOURUT UNIQUE (INSTALMENT_ID, NOURUT)
)
/
