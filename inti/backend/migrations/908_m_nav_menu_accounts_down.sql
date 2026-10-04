-- Mundur 908: buang hak menu Accounts dari setiap akun, lalu barisnya. Hak yang tertinggal akan ditolak Kelola User
-- sebagai menu tak dikenal saat akun itu diubah.
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'accounts'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'accounts'
/
