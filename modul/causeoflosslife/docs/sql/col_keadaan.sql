-- LANGKAH-WO-CAUSEOFLOSSLIFE.md bab Pemulihan - BACA-SAJA. Penentu keadaan sesudah 090-092 / 955 gagal di tengah;
-- bandingkan dengan tabel keputusan K0-K4 di dokumen. Dijalankan sebagai BERKAS (@...).
SET LINESIZE 250 PAGESIZE 100
PROMPT 1. Objek, kolom, catatan T_MIGRASI, dan menu
SELECT
  (SELECT COUNT(*) FROM SYS.ALL_VIEWS WHERE OWNER = 'POOLDATA' AND VIEW_NAME = 'CAUSEOFLOSS_LIFE')                  AS ADA_VIEW,
  (SELECT COUNT(*) FROM SYS.ALL_TABLES WHERE OWNER = 'POOLDATA' AND TABLE_NAME = 'M_CAUSEOFLOSS_LIFE')              AS ADA_SUMBER,
  (SELECT COUNT(*) FROM SYS.ALL_TABLES WHERE OWNER = 'POOLDATA' AND TABLE_NAME = 'CAUSEOFLOSS_LIFE')                AS ADA_TARGET,
  (SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS c JOIN SYS.ALL_TABLES t ON t.OWNER = c.OWNER AND t.TABLE_NAME = c.TABLE_NAME
    WHERE c.OWNER = 'POOLDATA' AND c.TABLE_NAME IN ('CAUSEOFLOSS_LIFE', 'M_CAUSEOFLOSS_LIFE')
      AND c.COLUMN_NAME = 'CAUSEOFLOSS')                                                                          AS ADA_KOLOM,
  (SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = 'POOLDATA' AND TABLE_NAME IN ('CAUSEOFLOSS_LIFE', 'M_CAUSEOFLOSS_LIFE')
      AND COLUMN_NAME = 'JSONDATA')                                                                               AS ADA_JSONDATA,
  (SELECT COUNT(*) FROM SYS.ALL_PARTIAL_DROP_TABS WHERE OWNER = 'POOLDATA' AND TABLE_NAME = 'CAUSEOFLOSS_LIFE')     AS SETENGAH_TERBUANG,
  (SELECT COUNT(*) FROM POOLDATA.T_MIGRASI WHERE NAMA = '090_causeofloss_life_ganti_nama')                         AS T090,
  (SELECT COUNT(*) FROM POOLDATA.T_MIGRASI WHERE NAMA = '091_causeofloss_life_kolom')                              AS T091,
  (SELECT COUNT(*) FROM POOLDATA.T_MIGRASI WHERE NAMA = '092_causeofloss_life_satu_tabel')                         AS T092,
  (SELECT COUNT(*) FROM POOLDATA.T_MIGRASI WHERE NAMA = '955_menu_causeoflosslife')                                AS T955,
  (SELECT COUNT(*) FROM POOLDATA.M_NAV_MENU WHERE KODE = 'causeoflosslife')                                        AS ADA_MENU
FROM DUAL;
PROMPT 2. Isi kolom (ORA-00904 / ORA-00942 = kolom / tabel belum ada pada nama itu, wajar)
SELECT COUNT(*) AS N, COUNT(CAUSEOFLOSS) AS N_CAUSEOFLOSS FROM POOLDATA.CAUSEOFLOSS_LIFE;
