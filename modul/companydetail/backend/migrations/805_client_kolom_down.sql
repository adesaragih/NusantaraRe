-- Mundur 805: buang kolom tambahan CLIENT (isinya hilang).
ALTER TABLE {skema}.CLIENT DROP (PARENT_ID, NOTE, CREATED_BY, CREATED_AT, UPDATED_BY, UPDATED_AT)
/
