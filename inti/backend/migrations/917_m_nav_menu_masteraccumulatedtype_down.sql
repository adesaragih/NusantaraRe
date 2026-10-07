-- 917 mundur - hak menu lalu baris M_NAV_MENU `masteraccumulatedtype` dibuang (pola 906 mundur).
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'masteraccumulatedtype'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'masteraccumulatedtype'
/
