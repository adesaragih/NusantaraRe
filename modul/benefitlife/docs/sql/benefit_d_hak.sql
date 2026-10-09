-- LANGKAH-WO-BENEFITLIFE.md (d) - MENULIS hak menu (dijalankan WO SESUDAH 945; K5). Hak menu Benefit (benefitlife) =
-- GABUNGAN hak R/I Rate Life (riratelife), R/I Comm Life (ricommlife), dan R/I Risk (ririsklife): setiap akun yang
-- memegang salah satunya mendapat benefitlife; PENUH bila salah satunya PENUH, selain itu LIHAT. Pola HAK-RIRISKLIFE.sql.
-- HAK CHECK ('LIHAT','PENUH') (914); PK (LOGIN_ID, MENU_KODE) (903) - NOT EXISTS membuatnya aman diulang (akun yang
-- sudah punya benefitlife tidak diubah). Jalankan: sqlplus <akun>@<db> @benefit_d_hak.sql
SET PAGESIZE 100 LINESIZE 150 FEEDBACK ON
WHENEVER SQLERROR EXIT FAILURE ROLLBACK
COLUMN MENU_KODE FORMAT A14

INSERT INTO POOLDATA.M_LOGIN_GO_MENU (LOGIN_ID, MENU_KODE, HAK)
SELECT m.LOGIN_ID, 'benefitlife',
       CASE WHEN MAX(CASE m.HAK WHEN 'PENUH' THEN 1 ELSE 0 END) = 1 THEN 'PENUH' ELSE 'LIHAT' END
  FROM POOLDATA.M_LOGIN_GO_MENU m
 WHERE m.MENU_KODE IN ('riratelife', 'ricommlife', 'ririsklife')
   AND NOT EXISTS (SELECT 1 FROM POOLDATA.M_LOGIN_GO_MENU x
                    WHERE x.LOGIN_ID = m.LOGIN_ID AND x.MENU_KODE = 'benefitlife')
 GROUP BY m.LOGIN_ID;
COMMIT;

-- Periksa: cacah akun benefitlife = cacah akun yang memegang salah satu dari ketiga menu; HAK sesuai gabungan.
SELECT COUNT(DISTINCT LOGIN_ID) AS AKUN_SUMBER FROM POOLDATA.M_LOGIN_GO_MENU
 WHERE MENU_KODE IN ('riratelife', 'ricommlife', 'ririsklife');
SELECT COUNT(*) AS AKUN_BENEFITLIFE FROM POOLDATA.M_LOGIN_GO_MENU WHERE MENU_KODE = 'benefitlife';
SELECT LOGIN_ID, MENU_KODE, HAK FROM POOLDATA.M_LOGIN_GO_MENU
 WHERE MENU_KODE IN ('riratelife', 'ricommlife', 'ririsklife', 'benefitlife')
 ORDER BY LOGIN_ID, MENU_KODE;
EXIT
