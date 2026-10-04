-- Jalur mundur 911 (dulu 904) - baris menu masterdata dan akses akun atasnya (bila sudah diberikan lewat Kelola User).
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'masterdata'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'masterdata'
/
