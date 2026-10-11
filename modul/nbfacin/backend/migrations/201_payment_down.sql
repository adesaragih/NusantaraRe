-- Jalur mundur 201 - Payment: daftar cicilan, Total Payment / Net Premium per mata uang, dan jumlah cicilan ikut hilang.
DROP SEQUENCE {skema}.SEQ_T_LISTINSTALLMENT
/
DROP TABLE {skema}.T_LISTINSTALLMENT CASCADE CONSTRAINTS
/
ALTER TABLE {skema}.T_CURRENCYLIST DROP (SUM_TOTAL_PAYMENT, PAY_NET_PREMIUM)
/
ALTER TABLE {skema}.T_GENERAL_POLIS DROP (PAY_INSTALLMENT)
/
