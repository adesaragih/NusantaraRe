-- Mundur 916: buang hak menu Treaty Group OJK dari setiap akun, lalu barisnya.
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'treatygroupojk'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'treatygroupojk'
/
