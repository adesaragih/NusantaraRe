-- Jalur mundur 914 - kolom HAK dibuang; setiap menu kembali berarti akses penuh. Jalur mundur ini untuk skema uji.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029).
ALTER TABLE {skema}.M_LOGIN_GO_MENU DROP CONSTRAINT CK_M_LOGIN_GO_MENU_HAK
/
ALTER TABLE {skema}.M_LOGIN_GO_MENU DROP (HAK)
/
