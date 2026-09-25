-- SUMBER: DDL_Script_ClaimNonProp.xls (reverse-engineer Toad)
-- UNTUK DIBACA, BUKAN UNTUK DIJALANKAN. Klausa DROP/ALTER-DROP tidak boleh dieksekusi.


  CREATE OR REPLACE FORCE EDITIONABLE VIEW "POOLDATA"."V_D_CAUSE_OF_LOSS" ("D_COL_ID", "OLD_D_COL_ID", "M_COL_ID", "DESCRIPTION", "STS_AKTIF", "LOSS_CODE", "BISNISID", "BISNISNAME") AS 
  SELECT D_COL_ID,
          OLD_D_COL_ID,
          json_value (jsondata, '$.M_COL_ID'),
          json_value (jsondata, '$.DESCRIPTION'),
          json_value (jsondata, '$.STS_AKTIF'),
          json_value (jsondata, '$.LOSS_CODE'),
          BISNISID,
          (SELECT Note from business where ID = BISNISID)
     FROM D_CAUSE_OF_LOSS, JSON_TABLE (jsondata, '$.BISNISID[*]'
                            COLUMNS( 
                            BISNISID VARCHAR2(100) PATH '$.ID'))