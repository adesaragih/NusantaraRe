-- Mundur 941: buang hak menu R/I Risk dari setiap akun, lalu barisnya.
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'ririsklife'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'ririsklife'
/
