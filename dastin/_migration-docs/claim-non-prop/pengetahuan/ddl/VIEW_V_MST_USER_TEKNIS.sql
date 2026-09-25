-- SUMBER: DDL_Script_ClaimNonProp.xls (reverse-engineer Toad)
-- UNTUK DIBACA, BUKAN UNTUK DIJALANKAN. Klausa DROP/ALTER-DROP tidak boleh dieksekusi.


  CREATE OR REPLACE FORCE EDITIONABLE VIEW "POOLDATA"."V_MST_USER_TEKNIS" ("OPERATOR_ID", "OLD_OPERATOR_ID", "STS_AKTIF", "MCL_NAME", "TEAM_GROUP", "TOTAL_JOB", "COUNTER_QUOTA") AS 
  SELECT OPERATOR_ID,
          OLD_OPERATOR_ID,
          json_value (JSON_DATA, '$.STS_AKTIF'),
          (SELECT MCL_NAME
             FROM hrdasm.v_hrd_mst@asmd.sinarmas.co.id b
            WHERE a.OPERATOR_ID = b.login_aplikasi)
             AS MCL_NAME,
          json_value (JSON_DATA, '$.TEAM_GROUP'),
          TOTAL_JOB,
          json_value (JSON_DATA, '$.COUNTER_QUOTA')
     FROM MST_USER_TEKNIS a