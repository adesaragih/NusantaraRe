-- LANGKAH-WO-CAUSEOFLOSSLIFE.md (d) - MENULIS hak menu (dijalankan WO SESUDAH 955; K6). Hak menu Cause Of Loss Life
-- (causeoflosslife) = SALINAN pemegang Plan (planlife) beserta HAK-nya (LIHAT / PENUH apa adanya). HAK CHECK
-- ('LIHAT','PENUH') (914); PK (LOGIN_ID, MENU_KODE) (903) - NOT EXISTS membuatnya aman diulang.
-- Jalankan: sqlplus <akun>@<db> @col_d_hak.sql
SET PAGESIZE 100 LINESIZE 150 FEEDBACK ON
WHENEVER SQLERROR EXIT FAILURE ROLLBACK
COLUMN MENU_KODE FORMAT A16

INSERT INTO POOLDATA.M_LOGIN_GO_MENU (LOGIN_ID, MENU_KODE, HAK)
SELECT m.LOGIN_ID, 'causeoflosslife', m.HAK
  FROM POOLDATA.M_LOGIN_GO_MENU m
 WHERE m.MENU_KODE = 'planlife'
   AND NOT EXISTS (SELECT 1 FROM POOLDATA.M_LOGIN_GO_MENU x
                    WHERE x.LOGIN_ID = m.LOGIN_ID AND x.MENU_KODE = 'causeoflosslife');
COMMIT;

-- Periksa: cacah akun causeoflosslife = cacah akun planlife; HAK sama per akun (kueri ketiga NOL baris).
SELECT COUNT(*) AS AKUN_PLANLIFE FROM POOLDATA.M_LOGIN_GO_MENU WHERE MENU_KODE = 'planlife';
SELECT COUNT(*) AS AKUN_CAUSEOFLOSSLIFE FROM POOLDATA.M_LOGIN_GO_MENU WHERE MENU_KODE = 'causeoflosslife';
SELECT p.LOGIN_ID, p.HAK AS HAK_PLAN, c.HAK AS HAK_COL FROM POOLDATA.M_LOGIN_GO_MENU p
  LEFT JOIN POOLDATA.M_LOGIN_GO_MENU c ON c.LOGIN_ID = p.LOGIN_ID AND c.MENU_KODE = 'causeoflosslife'
 WHERE p.MENU_KODE = 'planlife' AND (c.HAK IS NULL OR c.HAK <> p.HAK);
EXIT
