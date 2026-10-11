-- Membalik `445_kolom_share_revisi.sql` — membuang keempat kolom akar Share
-- dari `T_TREATY_REVISION` (tabel modul `treatyin`, diparkir di sini).
--
-- ⚠️ Nilai yang sudah disimpan di keempat kolom ikut hilang.

ALTER TABLE {skema}.T_TREATY_REVISION DROP (
  RNMSHARE,
  BROKERAGEPERCENT,
  RNMSHAREACROSSTHEBOARD,
  RNMSHAREDEDUCTED
)
/
