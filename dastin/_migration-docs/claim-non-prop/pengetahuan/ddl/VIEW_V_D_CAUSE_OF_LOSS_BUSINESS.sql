-- SUMBER: DDL_Script_ClaimNonProp.xls (reverse-engineer Toad)
-- UNTUK DIBACA, BUKAN UNTUK DIJALANKAN. Klausa DROP/ALTER-DROP tidak boleh dieksekusi.


  CREATE OR REPLACE FORCE EDITIONABLE VIEW "POOLDATA"."V_D_CAUSE_OF_LOSS_BUSINESS" ("D_COL_ID", "OLD_D_COL_ID", "M_COL_ID", "DESCRIPTION", "STS_AKTIF", "LOSS_CODE", "BISNISID") AS 
  SELECT DISTINCT a.D_COL_ID,
                   a.OLD_D_COL_ID,
                   c.M_COL_ID,
                   c.DESCRIPTION,
                   c.STS_AKTIF,
                   c.LOSS_CODE,
                   c.idbisnis AS BISNISID
     FROM pooldata.D_CAUSE_OF_LOSS a, json_table (
a.JSONDATA,'$'    COLUMNS(
                                        D_COL_ID VARCHAR PATH '$.D_COL_ID',
                                        DESCRIPTION VARCHAR PATH '$.DESCRIPTION',
                                        M_COL_ID VARCHAR PATH '$.M_COL_ID',
                                        STS_AKTIF VARCHAR PATH '$.STS_AKTIF',
                                        LOSS_CODE VARCHAR PATH '$.LOSS_CODE',
                                         nested path '$.BISNISID[*]'
                                        columns(idbisnis varchar2 path '$.ID'))) c where  c.D_COL_ID = a.D_COL_ID