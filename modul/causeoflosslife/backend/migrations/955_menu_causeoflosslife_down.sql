-- Mundur 955: buang hak menu Cause Of Loss Life dari setiap akun, lalu barisnya.
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'causeoflosslife'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'causeoflosslife'
/
