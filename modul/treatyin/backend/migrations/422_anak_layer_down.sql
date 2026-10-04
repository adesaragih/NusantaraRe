-- Jalur mundur untuk 422_anak_layer.sql
--
-- Membongkar ketiga anak cabang LAYER. Urutannya terbalik dari jalur maju:
-- KELAS_BISNIS_KELOMPOK lebih dulu, sebab ia menggantung pada KELOMPOK_LAYER.
--
-- ⚠️ Index tidak dibongkar tersendiri - DROP TABLE membawa index miliknya.
DROP TABLE {skema}.KELAS_BISNIS_KELOMPOK
/
DROP TABLE {skema}.KELOMPOK_LAYER
/
DROP TABLE {skema}.KELAS_BISNIS_LAYER
/
