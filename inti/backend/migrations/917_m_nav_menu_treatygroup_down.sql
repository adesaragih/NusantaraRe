-- Mundur 917: buang hak menu Treaty Group dari setiap akun, lalu barisnya.
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'treatygroup'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'treatygroup'
/
