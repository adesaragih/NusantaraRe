-- Mundur 949: buang hak menu Plan dari setiap akun, lalu barisnya.
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'planlife'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'planlife'
/
