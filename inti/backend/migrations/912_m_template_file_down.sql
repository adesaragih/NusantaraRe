-- Jalur mundur 912 - sequence lalu tabel.
--
-- ⚠️ Seluruh berkas templat yang pernah diunggah HILANG; setiap menu kembali memakai berkas bawaannya. Jalur
-- mundur ini untuk skema uji.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029).
DROP SEQUENCE {skema}.SEQ_M_TEMPLATE_FILE
/
DROP TABLE {skema}.M_TEMPLATE_FILE CASCADE CONSTRAINTS
/
