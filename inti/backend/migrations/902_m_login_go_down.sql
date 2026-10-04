-- Jalur mundur 902 - anak dulu, lalu induk. Index ikut terbuang bersama
-- tabelnya.
--
-- ⚠️ Seluruh akun, hash sandi, dan workbasket-nya HILANG. Jalur mundur ini
-- untuk skema uji.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029).
DROP TABLE {skema}.M_LOGIN_GO_WORKBASKET CASCADE CONSTRAINTS
/
DROP TABLE {skema}.M_LOGIN_GO CASCADE CONSTRAINTS
/
