-- Jalur mundur 484: index dan sequence kembali ke bentuk 482 / 481 (`EDMLF-%`, START WITH 1).
--
-- ⚠️ Kasus `EDML-<n>` yang sudah lahir TIDAK ikut mundur dan tidak lagi terlindung index
-- terbuka-ganda; sequence mulai lagi dari 1. Untuk skema uji.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029).
DROP SEQUENCE {skema}.SEQ_WORK_EDM_LIFE
/
CREATE SEQUENCE {skema}.SEQ_WORK_EDM_LIFE START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
DROP INDEX {skema}.UX_PL_EDM_TERBUKA
/
CREATE UNIQUE INDEX {skema}.UX_PL_EDM_TERBUKA ON {skema}.T_PREMIUM_LIST (
  CASE WHEN STATUSS IS NULL AND ID LIKE 'EDMLF-%' THEN OLD_POLICY_NO END
)
/
