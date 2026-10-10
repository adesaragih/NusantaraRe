-- Mundur 681 - indeks UNIK UX_GENERAL_KOMITE_ADJ dikembalikan; untuk skema uji.
--
-- ⚠️ GAGAL (ORA-01452) bila sudah ada ADJUSTMENT_ID kembar lintas lini - justru kasus yang membuat 681 perlu.
--
-- NOL COMMIT (ADR-U-0029).
DROP INDEX {skema}.IX_GENERAL_KOMITE_ADJ
/
CREATE UNIQUE INDEX {skema}.UX_GENERAL_KOMITE_ADJ ON {skema}.T_GENERAL_KOMITE (ADJUSTMENT_ID)
/
