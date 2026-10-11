-- Pembongkar `451`. Membuang SATU kolom yang `451` tambahkan — nol yang lain.
--
-- ⚠️ Membongkarnya mengembalikan laporan `kunciTakTersimpan: RevisionDate`
-- pada setiap Save penyesuaian; itu yang DIHARAPKAN, bukan kemunduran diam.

ALTER TABLE {skema}.T_TREATY_REVISION DROP (
  REVISIONDATE
)
/
