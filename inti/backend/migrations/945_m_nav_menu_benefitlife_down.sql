-- Mundur 945: buang hak menu Benefit dari setiap akun, lalu barisnya.
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'benefitlife'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'benefitlife'
/
