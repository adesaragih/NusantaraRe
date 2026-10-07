-- Jalur mundur 483 - sub-jenis EDM Type yang sudah tertulis HILANG bersama kolomnya.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029).
ALTER TABLE {skema}.T_PREMIUM_LIST DROP (
  EDM_TYPE_PERUBAHAN_DATA, EDM_TYPE_BATAL
)
/
