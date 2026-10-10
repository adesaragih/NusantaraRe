-- 255 - % Pro Rate tingkat CASE tab Coverage Aneka (tiket 73 P2): T_GENERAL_POLIS.PRO_RATE_PERCENT.
--
-- `[terverifikasi]` sel % Pro Rate `NB FacIn\Section\InputCoverageAneka_FacIn.xml` L2713 terikat ke
-- `pyWorkPage.OfferFacIn.ProRatePercent` (satu nilai per case, bukan per coverage); ditulis CountPremiCoverageAneka langkah 4
-- (L964), FillPremiGolf langkah 2 (L399), ChangeCoveragePremium langkah 2, CalculateRate langkah 5. Kolom RANCANGAN
-- T_GENERAL_POLIS (`loader/skema_gen.go`, medan ProRatePercent, NUMBER) -> NUMBER(38,8) (persen, ADR-0016); tidak dibuat 182.
-- Rentang KEDUA nbfacin 255-259 (keputusan work owner 09-10-2026): 214 sudah dijalankan di DEV dan tidak boleh diubah. NULL =
-- belum pernah diisi (aplikasi menghitungnya dari tanggal polis). Ditulis, TIDAK dijalankan agent. Nol COMMIT (ADR-U-0029).
ALTER TABLE {skema}.T_GENERAL_POLIS ADD (
  PRO_RATE_PERCENT NUMBER(38,8)
)
/
