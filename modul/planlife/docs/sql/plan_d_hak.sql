-- LANGKAH-WO-PLANLIFE.md (d) - MENULIS hak menu (dijalankan WO SESUDAH 949; K7). Hak menu Plan (planlife) = GABUNGAN
-- hak R/I Rate Life (riratelife), R/I Comm Life (ricommlife), R/I Risk (ririsklife), dan Benefit (benefitlife): setiap
-- akun yang memegang salah satunya mendapat planlife; PENUH bila salah satunya PENUH, selain itu LIHAT (pola
-- benefit_d_hak.sql). HAK CHECK ('LIHAT','PENUH') (914); PK (LOGIN_ID, MENU_KODE) (903) - NOT EXISTS membuatnya aman
-- diulang. Jalankan: sqlplus <akun>@<db> @plan_d_hak.sql
SET PAGESIZE 100 LINESIZE 150 FEEDBACK ON
WHENEVER SQLERROR EXIT FAILURE ROLLBACK
COLUMN MENU_KODE FORMAT A14

INSERT INTO POOLDATA.M_LOGIN_GO_MENU (LOGIN_ID, MENU_KODE, HAK)
SELECT m.LOGIN_ID, 'planlife',
       CASE WHEN MAX(CASE m.HAK WHEN 'PENUH' THEN 1 ELSE 0 END) = 1 THEN 'PENUH' ELSE 'LIHAT' END
  FROM POOLDATA.M_LOGIN_GO_MENU m
 WHERE m.MENU_KODE IN ('riratelife', 'ricommlife', 'ririsklife', 'benefitlife')
   AND NOT EXISTS (SELECT 1 FROM POOLDATA.M_LOGIN_GO_MENU x
                    WHERE x.LOGIN_ID = m.LOGIN_ID AND x.MENU_KODE = 'planlife')
 GROUP BY m.LOGIN_ID;
COMMIT;

-- Periksa: cacah akun planlife = cacah akun yang memegang salah satu dari keempat menu; HAK sesuai gabungan.
SELECT COUNT(DISTINCT LOGIN_ID) AS AKUN_SUMBER FROM POOLDATA.M_LOGIN_GO_MENU
 WHERE MENU_KODE IN ('riratelife', 'ricommlife', 'ririsklife', 'benefitlife');
SELECT COUNT(*) AS AKUN_PLANLIFE FROM POOLDATA.M_LOGIN_GO_MENU WHERE MENU_KODE = 'planlife';
SELECT LOGIN_ID, MENU_KODE, HAK FROM POOLDATA.M_LOGIN_GO_MENU
 WHERE MENU_KODE IN ('riratelife', 'ricommlife', 'ririsklife', 'benefitlife', 'planlife')
 ORDER BY LOGIN_ID, MENU_KODE;
EXIT
