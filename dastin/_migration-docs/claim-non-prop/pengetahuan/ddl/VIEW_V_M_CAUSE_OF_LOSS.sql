-- SUMBER: DDL_Script_ClaimNonProp.xls (reverse-engineer Toad)
-- UNTUK DIBACA, BUKAN UNTUK DIJALANKAN. Klausa DROP/ALTER-DROP tidak boleh dieksekusi.


  CREATE OR REPLACE FORCE EDITIONABLE VIEW "POOLDATA"."V_M_CAUSE_OF_LOSS" ("M_COL_ID", "OLD_M_COL_ID", "COL_DESC") AS 
  SELECT M_COL_ID,
          OLD_M_COL_ID,
          json_value (JSON_DATA, '$.COL_DESC')
          
     FROM M_CAUSE_OF_LOSS