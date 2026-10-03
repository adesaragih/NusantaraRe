-- Mundur 906: buang hak menu Marketing Officer dari setiap akun, lalu barisnya. Hak yang tertinggal akan ditolak
-- Kelola User sebagai menu tak dikenal saat akun itu diubah.
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'marketingofficer'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'marketingofficer'
/
