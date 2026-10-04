-- 482 - "hanya ada satu endorsement terbuka per polis" (spec user story 5, AC 3:
-- aturan integritas inti), ditegakkan basis data, bukan hanya gerbang layar.
--
-- Gerbang 3 Pega (`SetErrorBatalEndorsement_Act` 3.4-3.6, RD
-- `FilterProteksiEDMLife` b541 `.PolicyNo = Param.Nopolis`) hanya membaca: dua
-- pembuatan bersamaan atas polis yang sama sama-sama lolos. Index unik
-- berfungsi ini menolak kasus terbuka KEDUA atas `OLD_POLICY_NO` yang sama
-- (ORA-00001, dijawab services dengan pesan gerbang 3 VERBATIM).
--
-- ⚠️ Ekspresinya kosong untuk setiap baris selain kasus EDM terbuka - baris
-- new business PremiumList (`STATUSS` kosong, `ID` `NBLF-`) dan kasus EDM yang
-- sudah diputuskan tidak masuk index sama sekali.
-- ⚠️ Kasus `Resolved-Rejected` TIDAK memblokir (spec AC 31 `[keputusan work
-- owner]`, OQ-EDM-002) - di Pega `FilterProteksiEDMLife` b526 hanya
-- mengecualikan `Resolved-Completed`.
CREATE UNIQUE INDEX {skema}.UX_PL_EDM_TERBUKA ON {skema}.T_PREMIUM_LIST (
  CASE WHEN STATUSS IS NULL AND ID LIKE 'EDMLF-%' THEN OLD_POLICY_NO END
)
/
