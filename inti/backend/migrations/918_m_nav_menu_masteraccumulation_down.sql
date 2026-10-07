-- 918 mundur - hak menu lalu baris M_NAV_MENU `masteraccumulation` dibuang (pola 906 mundur).
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'masteraccumulation'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'masteraccumulation'
/
