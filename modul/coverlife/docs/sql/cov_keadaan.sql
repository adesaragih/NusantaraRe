-- LANGKAH-WO-COVERLIFE.md bab Pemulihan - BACA-SAJA. Penentu keadaan sesudah 085-086 / 957 gagal di tengah; bandingkan
-- dengan tabel keputusan K0-K4 di dokumen. Dijalankan sebagai BERKAS (@...).
SET LINESIZE 250 PAGESIZE 100
PROMPT 1. Objek, kolom, catatan T_MIGRASI, dan menu
SELECT
  (SELECT COUNT(*) FROM SYS.ALL_VIEWS WHERE OWNER = 'POOLDATA' AND VIEW_NAME = 'COVER_LIFE')                         AS ADA_VIEW,
  (SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = 'POOLDATA' AND TABLE_NAME = 'M_COVER_LIFE'
      AND COLUMN_NAME = 'COVER')                                                                                   AS ADA_COVER,
  (SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = 'POOLDATA' AND TABLE_NAME = 'M_COVER_LIFE'
      AND COLUMN_NAME = 'NOTE')                                                                                    AS ADA_NOTE,
  (SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = 'POOLDATA' AND TABLE_NAME = 'M_COVER_LIFE'
      AND COLUMN_NAME = 'JSONDATA')                                                                                AS ADA_JSONDATA,
  (SELECT COUNT(*) FROM SYS.ALL_PARTIAL_DROP_TABS WHERE OWNER = 'POOLDATA' AND TABLE_NAME = 'M_COVER_LIFE')          AS SETENGAH_TERBUANG,
  (SELECT COUNT(*) FROM POOLDATA.T_MIGRASI WHERE NAMA = '085_cover_life_kolom')                                     AS T085,
  (SELECT COUNT(*) FROM POOLDATA.T_MIGRASI WHERE NAMA = '086_cover_life_satu_tabel')                                AS T086,
  (SELECT COUNT(*) FROM POOLDATA.T_MIGRASI WHERE NAMA = '957_menu_coverlife')                                       AS T957,
  (SELECT COUNT(*) FROM POOLDATA.M_NAV_MENU WHERE KODE = 'coverlife')                                               AS ADA_MENU
FROM DUAL;
PROMPT 2. Isi kolom (ORA-00904 = kolom belum ada, wajar)
SELECT COUNT(*) AS N, COUNT(COVER) AS N_COVER, COUNT(NOTE) AS N_NOTE FROM POOLDATA.M_COVER_LIFE;
