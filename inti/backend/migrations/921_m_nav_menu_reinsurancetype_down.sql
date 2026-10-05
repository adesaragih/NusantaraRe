-- Mundur 921: buang hak menu Reinsurance Type dari setiap akun, lalu barisnya.
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'reinsurancetype'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'reinsurancetype'
/
