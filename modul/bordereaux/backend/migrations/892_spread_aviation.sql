-- 892 - kolom Spread of Claim / Spread of Risk tabel Aviation dilebarkan dari NUMBER(10,4) ke NUMBER(38,8)
-- (perbaikan DDL warisan, 04-10-2026: Copy Old Data BDX-2026.08.05.59348 ditolak).
--
-- Kolom ini berisi NOMINAL, bukan persen: di 27 tabel detail lain kolom spread bertipe NUMBER bebas dan berisi
-- sampai 11 digit (DEV, agregat), dan pada JSON berkas di atas OR dan QS masing-masing 25% dari Total Incurred.
-- NUMBER(10,4) hanya muat 6 digit bulat, sehingga klaim/premi Aviation bernilai sejuta ke atas tidak pernah bisa
-- disimpan - di Pega pun Resolve-Complete akan gagal ORA-01438. Kedua tabel kosong di DEV saat migrasi ditulis.
-- Melebarkan presisi dan skala NUMBER selalu diizinkan Oracle, berisi atau tidak.
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
ALTER TABLE {skema}.BORDEREAUX_CLAIM_AVIATION MODIFY (
  SPREAD_OF_CLAIM_OR     NUMBER(38,8),
  SPREAD_OF_CLAIM_QS     NUMBER(38,8),
  SPREAD_OF_CLAIM_OTHERS NUMBER(38,8)
)
/
ALTER TABLE {skema}.BORDEREAUX_PREMI_AVIATION MODIFY (
  SPREAD_OF_RISK_OR     NUMBER(38,8),
  SPREAD_OF_RISK_QS     NUMBER(38,8),
  SPREAD_OF_RISK_SPL    NUMBER(38,8),
  SPREAD_OF_RISK_OTHERS NUMBER(38,8)
)
/
