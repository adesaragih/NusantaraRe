-- Jalur mundur 200 - Premium Deduction: CurrencyList pembayaran, persen, dan centang Exclude PPN ikut hilang.
DROP SEQUENCE {skema}.SEQ_T_CURRENCYLIST
/
DROP TABLE {skema}.T_CURRENCYLIST CASCADE CONSTRAINTS
/
ALTER TABLE {skema}.T_GENERAL_POLIS DROP (PPN_CHECK, PAY_RI_COMMISION, PAY_PCT_BROKERAGE_FEE)
/
