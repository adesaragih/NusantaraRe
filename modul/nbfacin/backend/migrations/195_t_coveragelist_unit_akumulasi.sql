-- 195 - tab Coverage FIRE (tiket 46): T_COVERAGELIST + tiga kolom rancangan - Indemnity Unit dan akumulasi.
--
-- Ketiganya ADA di rancangan T_COVERAGELIST (skema loader): UNIT VARCHAR2(50) (`.Unit`), ACCUMULATION_CODE VARCHAR2(50)
-- (`.AccumulationCode`), ACCUMULATION_DESCRIPTION VARCHAR2(500) (`.AccumulationDescription`) - tipe dan lebar rancangan
-- apa adanya (pola A109: kolom sebagian 193 ditambah lewat ALTER). Sampel `DDL\P-5 *.txt` [terverifikasi]: Unit satu
-- digit, AccumulationCode `aaa-99999-999999`, AccumulationDescription paling panjang 149 byte. Nol COMMIT (ADR-U-0029).
ALTER TABLE {skema}.T_COVERAGELIST ADD (
  UNIT                     VARCHAR2(50),
  ACCUMULATION_CODE        VARCHAR2(50),
  ACCUMULATION_DESCRIPTION VARCHAR2(500)
)
/
