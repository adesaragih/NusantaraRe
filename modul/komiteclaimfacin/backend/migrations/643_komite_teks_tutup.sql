-- 643 - T_GENERAL_KOMITE: teks pop-up kasus komite TT3 Reject Claim / TT4 Close Without Payment (jawaban work owner
-- 10-10-2026 OQ-KCFI-03, "ADD saja, bukan MODIFY"). `SendRejectClaimToKomite2` / `SendCloseClaimToKomite` 7.2 menyalin
-- `TempCommiteClaim.CircumtansesCouseOfLoss` / `.ExtentOfLoss` / `.LegalLiability` ke `childPageKomite.Komite.*`;
-- Section `ShowTransfer` LS39 menampilkannya. Kasus TT3 / TT4 lahir tanpa adjustment (641), jadi teksnya tidak punya
-- penampung lain (TT2 menyimpannya di baris adjustment `KOMITE_*`, claimprop 528) - nama dan lebar kolom sama.
--
-- Nullable tanpa DEFAULT: baris lama dan kasus Life / Prop / Non Prop / Fac In TT2 tidak menyebut kolom ini. Tabelnya
-- sudah ada (claimlife/013): nol CREATE TABLE.
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
ALTER TABLE {skema}.T_GENERAL_KOMITE ADD (
  KOMITE_CIRCUM_CAUSE_OF_LOSS VARCHAR2(4000),
  KOMITE_EXTENT_OF_LOSS VARCHAR2(4000),
  KOMITE_LEGAL_LIABILITY VARCHAR2(4000)
)
/
