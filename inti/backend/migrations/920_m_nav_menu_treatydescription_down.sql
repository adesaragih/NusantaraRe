-- Mundur 920: buang hak menu Treaty Description dari setiap akun, lalu barisnya.
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'treatydescription'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'treatydescription'
/
