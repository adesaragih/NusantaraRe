-- Jalur mundur 062. ⚠️ Baris status (ID <> ID_PEGA) ikut dibuang lebih dulu:
-- tanpa kolomnya, baris itu tidak dapat dibedakan dari header biasa.
DELETE FROM {skema}.T_PREMIUM_LIST WHERE ID <> ID_PEGA AND STATUS_PENAWARAN IS NOT NULL
/
ALTER TABLE {skema}.T_PREMIUM_LIST DROP COLUMN STATUS_PENAWARAN
/
