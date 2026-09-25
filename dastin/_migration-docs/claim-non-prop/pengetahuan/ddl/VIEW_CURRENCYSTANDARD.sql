-- SUMBER: DDL_Script_ClaimNonProp.xls (reverse-engineer Toad)
-- UNTUK DIBACA, BUKAN UNTUK DIJALANKAN. Klausa DROP/ALTER-DROP tidak boleh dieksekusi.


  CREATE OR REPLACE FORCE EDITIONABLE VIEW "POOLDATA"."CURRENCYSTANDARD" ("ID", "CURRENCYVALUE", "CURRENCYDATE", "USERID", "INPUTDATE") AS 
  SELECT a.ID,
          a.CurrencyValue,
          a.CurrencyDate,
          a.UserID,
          a.InputDate
     FROM m_CurrencyStandard a