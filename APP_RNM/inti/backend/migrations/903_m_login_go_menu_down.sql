-- Jalur mundur 903 - menu per akun ikut terbuang bersama tabelnya.
--
-- ⚠️ Sesudahnya gerbang menu menolak setiap akun (tabelnya tidak ada) sampai
-- 903 dijalankan lagi. Jalur mundur ini untuk skema uji.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029).
DROP TABLE {skema}.M_LOGIN_GO_MENU CASCADE CONSTRAINTS
/
