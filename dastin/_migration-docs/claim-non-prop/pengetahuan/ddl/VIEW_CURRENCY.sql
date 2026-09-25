-- SUMBER: DDL_Script_ClaimNonProp.xls (reverse-engineer Toad)
-- UNTUK DIBACA, BUKAN UNTUK DIJALANKAN. Klausa DROP/ALTER-DROP tidak boleh dieksekusi.


  CREATE OR REPLACE FORCE EDITIONABLE VIEW "POOLDATA"."CURRENCY" ("ID", "OLDID", "COUNTRYID", "COUNTRYNAME", "CURRENCYSYMBOL", "CURRENCY", "NOTE", "ISOSYMBOL") AS 
  SELECT b.JSONDATA.ID,
          b.OLDID,
          b.JSONDATA.CountryID,
          (SELECT a.JSONDATA.Note
             FROM m_nation a
            WHERE b.JSONDATA.CountryID = a.ID)
             AS CountryName,
          b.JSONDATA.CurrencySymbol,
          b.JSONDATA.Currency,
          b.JSONDATA.Note,
          b.JSONDATA.ISOSymbol
     FROM m_currency b