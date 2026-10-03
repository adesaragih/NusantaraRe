-- Mundur 905: buang index email dan index username tanpa beda huruf, kolom CONTACT_ID (index uniknya ikut
-- terbuang), dan sequence-nya.
-- ⚠️ Email yang sudah diubah ke huruf kecil TIDAK dikembalikan - bentuk aslinya tidak disimpan.
DROP INDEX {skema}.UX_M_LOGIN_GO_EMAIL
/
DROP INDEX {skema}.UX_M_LOGIN_GO_LOGIN_ID
/
ALTER TABLE {skema}.M_LOGIN_GO DROP (CONTACT_ID)
/
DROP SEQUENCE {skema}.M_LOGIN_GO_CONTACT_SEQ
/
