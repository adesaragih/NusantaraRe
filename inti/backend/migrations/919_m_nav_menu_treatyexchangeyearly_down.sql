-- Mundur 919: buang hak menu Treaty Exchange Yearly dari setiap akun, lalu barisnya.
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'treatyexchangeyearly'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'treatyexchangeyearly'
/
