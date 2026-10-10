-- 642 - T_GENERAL_KOMITE.TRANSFER_TYPE: jenis penyerahan kasus komite (keputusan work owner 10-10-2026 KCF-03; Pega
-- `pyWorkPage.TransferType`, dibaca `KomitePostAct` S2 / S4 / S5): '2' adjustment (`CreateKMTNo_Act` 6.9), '3' Reject
-- Claim (`SendRejectClaimToKomite2` 7.3), '4' Close Without Payment (`SendCloseClaimToKomite` 7.3). TT1 survey
-- ter-remark di KomitePostAct S3 - tidak diterima.
--
-- DEFAULT '2' WAJIB: INSERT kelahiran kasus komite Claim Life / Prop / Non Prop / Fac In TT2 tidak menyebut kolom ini,
-- dan baris yang sudah ada seluruhnya penyerahan adjustment. Tabelnya sudah ada (claimlife/013): nol CREATE TABLE.
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
ALTER TABLE {skema}.T_GENERAL_KOMITE ADD (
  TRANSFER_TYPE CHAR(1) DEFAULT '2' NOT NULL
)
/
ALTER TABLE {skema}.T_GENERAL_KOMITE ADD CONSTRAINT CK_GENERAL_KOMITE_TRANSFER CHECK (TRANSFER_TYPE IN ('2', '3', '4'))
/
