-- 682 - T_GENERAL_KOMITE: isian Subjectivity tingkat 1 disimpan antar tingkat (keputusan work owner 08-10-2026, OQ-KCP-01
-- jawaban "a"). Section ShowTransfer "Subjectivity ?" (.IsSubjectivity) dan "Subjectivity Note" (.SubjectivityNote) hanya
-- terbuka di tingkat 1 (`pyDisabledWhen .KomiteCount!='1'`) tetapi dibaca di tingkat akhir oleh KomitePostAdjustment
-- S16 / S17 / S21 / S23 / S24 / S34 (`pyWorkPage.IsSubjectivity`, halaman kerja yang bertahan antar tingkat).
--
-- DEFAULT '0' WAJIB: INSERT kelahiran kasus komite (Claim Prop `BuatKasusKomite`, Claim Life `kasuskomite.go`) tidak
-- menyebut kolom ini. Lebar catatan = `T_CLAIM_ADJUSTMENT.SUBJECTIVITY_NOTE` (1000), tujuan salinannya di S24.
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
ALTER TABLE {skema}.T_GENERAL_KOMITE ADD (
  KOMITE_SUBJECTIVITY CHAR(1) DEFAULT '0' NOT NULL,
  KOMITE_SUBJECTIVITY_NOTE VARCHAR2(1000)
)
/
ALTER TABLE {skema}.T_GENERAL_KOMITE ADD CONSTRAINT CK_GENERAL_KOMITE_SUBJ CHECK (KOMITE_SUBJECTIVITY IN ('0', '1'))
/
