-- LANGKAH-WO-DISEASELIFE.md bab Pemulihan - BACA-SAJA. Penentu keadaan sesudah D2 / 080-081 / 951 gagal di tengah;
-- bandingkan dengan tabel keputusan K0-K4 di dokumen. Dijalankan sebagai BERKAS (@...).
SET LINESIZE 250 PAGESIZE 100
PROMPT 1. Objek, catatan T_MIGRASI, menu, dan data
SELECT
  (SELECT COUNT(*) FROM SYS.ALL_SEQUENCES WHERE SEQUENCE_OWNER = 'POOLDATA' AND SEQUENCE_NAME = 'SEQ_DISEASE_LIFE')         AS ADA_SEQ,
  (SELECT COUNT(*) FROM SYS.ALL_CONSTRAINTS WHERE OWNER = 'POOLDATA' AND TABLE_NAME = 'DISEASE_LIFE'
      AND CONSTRAINT_NAME = 'PK_DISEASE_LIFE')                                                                           AS ADA_PK,
  (SELECT COUNT(*) FROM POOLDATA.T_MIGRASI WHERE NAMA = '080_seq_disease_life')                                          AS T080,
  (SELECT COUNT(*) FROM POOLDATA.T_MIGRASI WHERE NAMA = '081_disease_life_pk')                                           AS T081,
  (SELECT COUNT(*) FROM POOLDATA.T_MIGRASI WHERE NAMA = '951_menu_diseaselife')                                          AS T951,
  (SELECT COUNT(*) FROM POOLDATA.M_NAV_MENU WHERE KODE = 'diseaselife')                                                  AS ADA_MENU,
  (SELECT COUNT(*) FROM POOLDATA.DISEASE_LIFE)                                                                           AS N,
  (SELECT COUNT(*) FROM POOLDATA.DISEASE_LIFE WHERE ID IS NULL)                                                          AS N_ID_NULL,
  (SELECT COUNT(*) FROM (SELECT ID FROM POOLDATA.DISEASE_LIFE GROUP BY ID HAVING COUNT(*) > 1))                         AS N_ID_KEMBAR,
  (SELECT COUNT(*) FROM POOLDATA.DISEASE_LIFE WHERE ID = '102051' AND ICD_CODE = 'TEST123')                              AS ADA_BARIS_UJI
FROM DUAL;
