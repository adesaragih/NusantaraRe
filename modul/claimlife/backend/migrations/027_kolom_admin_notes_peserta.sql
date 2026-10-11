-- T_CLAIMLF_PREMIUMLIST_DETAIL.ADMIN_NOTES - keputusan work owner 11-10-2026 (rincian peserta layar Outstanding mengikuti
-- screenshot Pega CLMLF-9293; persetujuan migrasi diberikan 11-10-2026).
--
-- `[terverifikasi]` `Section/ClaimLifeDetailGCNM.xml` (korpus kelvin/InboxClaimLife (Done)): `ADMIN NOTES` b2588 terikat
-- `.ADMIN_NOTES`, dibekukan `pyDisabledWhen .STS_REJECT=='1' || .STS_REJECT=='2'` (b2628, STRUKTUR "RALAT 27-09-2026").
-- Sebanding NOTES / RECOMMENDATION (VARCHAR2(1000), migrasi 003) - lebar sama.
--
-- Nullable (ADR-U-0027), tanpa pengisian ulang.
--
-- ⛔ NOL COMMIT (ADR-U-0029). Dijalankan work owner sesudah giliran selesai.

ALTER TABLE {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL ADD (
  ADMIN_NOTES VARCHAR2(1000)
)
/
