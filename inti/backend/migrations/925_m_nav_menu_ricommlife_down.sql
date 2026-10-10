-- Mundur 925: buang hak menu R/I Comm Life dari setiap akun, lalu barisnya.
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'ricommlife'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'ricommlife'
/
