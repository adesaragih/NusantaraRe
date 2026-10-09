-- Pembongkar `452`. Membuang SATU kolom yang `452` tambahkan.
--
-- ⚠️ Membongkarnya mengembalikan laporan
-- `kunciTakTersimpan: ValueDifference.TotalLimitsROL` pada setiap Save
-- penyesuaian; itu yang DIHARAPKAN, bukan kemunduran diam.

ALTER TABLE {skema}.T_TREATY_REVISION DROP (
  VALUEDIFF_TOTALLIMITSROL
)
/
