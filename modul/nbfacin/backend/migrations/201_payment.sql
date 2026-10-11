-- 201 - tab Payment kasus FIRE (tiket 51): jumlah cicilan, Total Payment per mata uang, dan
-- `.OfferFacIn.CurrencyList(n).Policy.Payment.ListInstallment`.
--
-- Nama / jalur = workbook RANCANGAN (`loader/skema_gen.go`) `[terverifikasi]`:
--   - T_GENERAL_POLIS.PAY_INSTALLMENT (`.PolicyData.Payment.Installment`, lipatan V-22b; rancangan NUMBER -> NUMBER(5),
--     layar pyMax 2 digit);
--   - T_CURRENCYLIST.SUM_TOTAL_PAYMENT (`.SumTotalPayment`) dan PAY_NET_PREMIUM (`.Policy.Payment.NetPremium`);
--   - T_LISTINSTALLMENT, jalur `CurrencyList/ListInstallment`, induk T_CURRENCYLIST, dibuat SEBAGIAN: kolom sistem +
--     CURRENCY_CODE + kolom grid "Installment Data Information" (InstallmentNo, DueDate, InstallmentPercentage, Premium,
--     RICommision, Deduction2, BrokerageFee, PPh, PPN, Discount, PaymentTotal). AdminFee (hanya IsHealthIndividu), Stamp,
--     Commission, PaymentAmount, RateEM / RateLife, Year tidak dibuat.
-- Tipe: uang NUMBER(38,8) (ADR-0016; rancangan NUMBER polos; DEDUCTION2 rancangan VARCHAR2(50) -> NUMBER(38,8), uang).
-- INSTALLMENT_PERCENTAGE VARCHAR2(50) = teks desimal apa adanya (K51-3): persen Pega berskala 20 (contoh DDL\CONTOH
-- "33.33333333333333333333" x 2 + sisa = tepat 100) - NUMBER(38,8) membuat jumlahnya bukan 100 dan Save meratakan ulang
-- persen pengguna (pola T_TABLEOFLIMIT.PCT_LIMIT, butir 68.1). DUE_DATE VARCHAR2(30) teks Pega `YYYYMMDD` (200 dari 200
-- contoh). ID / PARENT_ID NUMBER(19) (A92; PARENT_ID = T_CURRENCYLIST.ID).
-- Ditulis, TIDAK dijalankan agent. Nol COMMIT (ADR-U-0029).
ALTER TABLE {skema}.T_GENERAL_POLIS ADD (
  PAY_INSTALLMENT NUMBER(5)
)
/
ALTER TABLE {skema}.T_CURRENCYLIST ADD (
  SUM_TOTAL_PAYMENT NUMBER(38,8),
  PAY_NET_PREMIUM   NUMBER(38,8)
)
/
CREATE TABLE {skema}.T_LISTINSTALLMENT (
  ID                     NUMBER(19) NOT NULL,
  IDPEGA                 VARCHAR2(50),
  COB_GROUP              VARCHAR2(20),
  PARENT_ID              NUMBER(19) NOT NULL,
  SEQ_NO                 NUMBER(5) NOT NULL,
  ROW_UID                VARCHAR2(36) NOT NULL,
  CURRENCY_CODE          VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL,
  INSTALLMENT_NO         NUMBER(5),
  DUE_DATE               VARCHAR2(30),
  INSTALLMENT_PERCENTAGE VARCHAR2(50),
  PREMIUM                NUMBER(38,8),
  RI_COMMISION           NUMBER(38,8),
  DEDUCTION2             NUMBER(38,8),
  BROKERAGE_FEE          NUMBER(38,8),
  PPH                    NUMBER(38,8),
  PPN                    NUMBER(38,8),
  DISCOUNT               NUMBER(38,8),
  PAYMENT_TOTAL          NUMBER(38,8),
  CONSTRAINT PK_T_LISTINSTALLMENT PRIMARY KEY (ID),
  CONSTRAINT FK_LISTINSTALLMENT_CURRENCY FOREIGN KEY (PARENT_ID) REFERENCES {skema}.T_CURRENCYLIST (ID)
)
/
CREATE INDEX {skema}.IX_LISTINSTALLMENT_PARENT ON {skema}.T_LISTINSTALLMENT (PARENT_ID, SEQ_NO)
/
CREATE SEQUENCE {skema}.SEQ_T_LISTINSTALLMENT START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
