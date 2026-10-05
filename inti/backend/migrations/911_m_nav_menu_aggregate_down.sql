-- Mundur 911: buang hak menu Aggregate dari setiap akun, lalu barisnya.
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'aggregate'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'aggregate'
/
