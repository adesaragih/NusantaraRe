-- 680 - T_GENERAL_KOMITE: dua penanda usul penyetuju (Section ShowTransfer "Propose To Close Case" .Adjustment.IsProposeClose,
-- "Propose To Reserved" .Adjustment.IsPropReserved). Keputusan work owner 18-09 dan 19-09-2026: disimpan di header kasus
-- komite, CHAR(1), wajib isi, hanya '1' / '0'; kotak yang tidak disentuh = '0' (keputusan 21).
--
-- DEFAULT '0' WAJIB: INSERT kelahiran kasus komite (Claim Prop `BuatKasusKomite`, Claim Life `kasuskomite.go`) tidak
-- menyebut kedua kolom ini. Tabelnya sudah ada (claimlife/013, komiteclaimlife/030): nol CREATE TABLE, nol MODIFY.
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
ALTER TABLE {skema}.T_GENERAL_KOMITE ADD (
  KOMITE_USUL_TUTUP CHAR(1) DEFAULT '0' NOT NULL,
  KOMITE_USUL_CADANG CHAR(1) DEFAULT '0' NOT NULL
)
/
ALTER TABLE {skema}.T_GENERAL_KOMITE ADD CONSTRAINT CK_GENERAL_KOMITE_USUL_TUTUP CHECK (KOMITE_USUL_TUTUP IN ('0', '1'))
/
ALTER TABLE {skema}.T_GENERAL_KOMITE ADD CONSTRAINT CK_GENERAL_KOMITE_USUL_CADANG CHECK (KOMITE_USUL_CADANG IN ('0', '1'))
/
