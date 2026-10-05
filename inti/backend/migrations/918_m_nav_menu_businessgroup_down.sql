-- Mundur 918: buang hak menu Business Group dari setiap akun, lalu barisnya.
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'businessgroup'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'businessgroup'
/
