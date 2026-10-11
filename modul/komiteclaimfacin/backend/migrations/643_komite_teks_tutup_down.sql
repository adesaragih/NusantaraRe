-- Mundur 643 - tiga kolom teks pop-up TT3 / TT4 dibuang beserta isinya (layar komite TT3 / TT4 kembali hanya
-- menampilkan Remarks klaim). Untuk skema uji.
--
-- NOL COMMIT (ADR-U-0029).
ALTER TABLE {skema}.T_GENERAL_KOMITE DROP (KOMITE_CIRCUM_CAUSE_OF_LOSS, KOMITE_EXTENT_OF_LOSS, KOMITE_LEGAL_LIABILITY)
/
