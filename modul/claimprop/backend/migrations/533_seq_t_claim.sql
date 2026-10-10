-- 533 - SEQ_T_CLAIM: identitas baris tabel anak Claim Prop (T_CLAIM_*).
-- Satu sequence untuk seluruh tabel anak: identitas cukup unik, bukan nomor dagang. Pengenal kasus CLMP- / TKMT-
-- memakai SEQ_WORK_CLAIM milik Claim Life (pola PengenalWorkBerikut, LPAD 6).
-- NOCACHE sama dengan seluruh sequence klaim.
CREATE SEQUENCE {skema}.SEQ_T_CLAIM START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
