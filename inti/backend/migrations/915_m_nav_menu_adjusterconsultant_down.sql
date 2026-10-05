-- Mundur 915: buang hak menu Adjuster Consultant dari setiap akun, lalu barisnya.
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'adjusterconsultant'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'adjusterconsultant'
/
