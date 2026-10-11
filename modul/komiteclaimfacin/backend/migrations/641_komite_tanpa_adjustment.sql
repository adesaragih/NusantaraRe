-- 641 - T_GENERAL_KOMITE.ADJUSTMENT_ID BOLEH KOSONG (keputusan work owner 10-10-2026 KCF-03, izin MODIFY tabel bersama):
-- kasus komite Claim Fac In TT3 Reject Claim (`SendRejectClaimToKomite2`) dan TT4 Close Without Payment
-- (`SendCloseClaimToKomite`) lahir TANPA baris adjustment. Hanya Fac In yang menulis baris tanpa adjustment; Komite
-- Claim Life / Prop / Non Prop tetap mengisinya (INSERT kelahiran mereka menyebut ADJUSTMENT_ID). Indeks biasa
-- IX_GENERAL_KOMITE_ADJ (komiteclaimprop 681) menerima NULL.
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
ALTER TABLE {skema}.T_GENERAL_KOMITE MODIFY (ADJUSTMENT_ID NULL)
/
