-- Mundur 922: buang hak menu R/I Rate Life dari setiap akun, lalu barisnya.
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'riratelife'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'riratelife'
/
