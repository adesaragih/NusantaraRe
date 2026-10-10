-- Mundur 642 - jenis penyerahan dibuang beserta isinya; baris TT3 / TT4 yang ada kehilangan penandanya (tetap
-- ber-ADJUSTMENT_ID kosong, lihat mundur 641). Untuk skema uji.
--
-- NOL COMMIT (ADR-U-0029).
ALTER TABLE {skema}.T_GENERAL_KOMITE DROP CONSTRAINT CK_GENERAL_KOMITE_TRANSFER
/
ALTER TABLE {skema}.T_GENERAL_KOMITE DROP (TRANSFER_TYPE)
/
