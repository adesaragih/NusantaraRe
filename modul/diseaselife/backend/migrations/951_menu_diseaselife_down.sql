-- Mundur 951: buang hak menu Disease Life dari setiap akun, lalu barisnya.
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'diseaselife'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'diseaselife'
/
