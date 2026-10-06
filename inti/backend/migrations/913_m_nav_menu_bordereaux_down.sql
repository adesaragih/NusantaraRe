-- Mundur 913: buang hak menu Bordereaux dari setiap akun, lalu barisnya.
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'bordereaux'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'bordereaux'
/
