-- SUMBER: DDL_Script_ClaimNonProp.xls (reverse-engineer Toad)
-- UNTUK DIBACA, BUKAN UNTUK DIJALANKAN. Klausa DROP/ALTER-DROP tidak boleh dieksekusi.


  CREATE OR REPLACE FORCE EDITIONABLE VIEW "POOLDATA"."PROVINCE" ("ID", "NATIONID", "NOTE", "NATIONNAME") AS 
  SELECT a.JSONDATA.ID AS ID,
  
          a.JSONDATA.NationID AS NationID,
  
          a.JSONDATA.Note AS Note,
  
          
 (Select b.JSONDATA.Note 
          FROM M_Nation b WHERE 
          a.JSONDATA.NationID = b.ID)
          AS NationName
     FROM M_PROVINCE a