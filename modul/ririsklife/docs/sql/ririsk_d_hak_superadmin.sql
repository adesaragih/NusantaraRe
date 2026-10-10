-- LANGKAH-WO-RIRISKLIFE.md (d) - MENULIS hak menu (dijalankan WO, pola LANGKAH-WO ricommlife (d)). Superadmin =
-- pemegang menu `kelolauser` (migrasi 914). HAK CHECK ('LIHAT','PENUH') (914); PK (LOGIN_ID, MENU_KODE) (903) -
-- NOT EXISTS membuatnya aman diulang. Akun lain: centang menu "R/I Risk" di Kelola User.
INSERT INTO POOLDATA.M_LOGIN_GO_MENU (LOGIN_ID, MENU_KODE, HAK)
SELECT m.LOGIN_ID, 'ririsklife', 'PENUH'
  FROM POOLDATA.M_LOGIN_GO_MENU m
 WHERE m.MENU_KODE = 'kelolauser'
   AND NOT EXISTS (SELECT 1 FROM POOLDATA.M_LOGIN_GO_MENU x
                    WHERE x.LOGIN_ID = m.LOGIN_ID AND x.MENU_KODE = 'ririsklife');
COMMIT;
-- Periksa: cacahnya = cacah pemegang kelolauser
SELECT COUNT(*) AS HAK_RIRISKLIFE FROM POOLDATA.M_LOGIN_GO_MENU WHERE MENU_KODE = 'ririsklife';
SELECT COUNT(*) AS PEMEGANG_KELOLAUSER FROM POOLDATA.M_LOGIN_GO_MENU WHERE MENU_KODE = 'kelolauser';
