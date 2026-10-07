-- 914 mundur - hak menu lalu baris M_NAV_MENU `mastercity` dibuang (pola 906 mundur).
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'mastercity'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'mastercity'
/
