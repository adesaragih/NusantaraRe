-- 481 - pengenal kasus Endorsement Life `EDMLF-<n>` (spec §16 AC 72, RALAT R20).
--
-- Kasus endorsement = baris versi `T_PREMIUM_LIST` ber-`ID` `EDMLF-<n>`; ia
-- TIDAK menulis `T_WORK_POLIS` (kotak masuk PremiumList membaca seluruh tabel
-- itu bila posisi kosong - `premiumlistlife/frontend/pages/InboxPremiumList.tsx`
-- `ambilKotakMasukPolis()`), jadi angkanya tidak dapat menumpang
-- `SEQ_WORK_POLIS`. Tanpa padding, sejajar `NBLF-<n>`.
--
-- ⚠️ `START WITH 1` - OQ-EDM-009: DBA menyetel ulang bila pengenal warisan
-- ternyata berbentuk sama (pola migrasi 058 PremiumList).
CREATE SEQUENCE {skema}.SEQ_WORK_EDM_LIFE START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
