-- T_CLAIMLF_PREMIUMLIST_DETAIL.MAX_DATA_RECEIVED - keputusan work owner 11-10-2026 ("DATE OF LOSS, CLAIM RECEIVED DATE,
-- DOCUMENT COMPLETE DATE, MAX DATA RECEIVED, CONFIRMATION DATE ini diisi oleh user, dan setelah save to rnm tidak bisa diisi
-- lagi"); persetujuan migrasi diberikan 11-10-2026.
--
-- `[terverifikasi]` `Section/ClaimLifeDetailGCNM.xml` (korpus kelvin/InboxClaimLife (Done)): MAX DATA RECEIVED b13381 terikat
-- `.MAXDATARECEIVED`. Di Pega nilainya DIHITUNG `ValidasiSTNC_Act` b577; kini diisi pemakai (keputusan work owner).
--
-- DATE, seperti keempat tanggal peserta lain. Nullable (ADR-U-0027), tanpa pengisian ulang.
--
-- ⛔ NOL COMMIT (ADR-U-0029). Dijalankan work owner sesudah giliran selesai.

ALTER TABLE {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL ADD (
  MAX_DATA_RECEIVED DATE
)
/
