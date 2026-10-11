-- 681 - T_GENERAL_KOMITE.ADJUSTMENT_ID: indeks UNIK UX_GENERAL_KOMITE_ADJ (claimlife/013) diganti indeks biasa
-- (keputusan work owner 08-10-2026). ID adjustment Prop (SEQ_T_CLAIM) dan Life (SEQ_CLAIMLF_ADJ) sama-sama angka polos;
-- indeks unik lintas lini suatu saat menolak penyerahan sah (ORA-00001).
--
-- Keunikan "satu baris adjustment <-> satu kasus komite" tetap dijaga di tabel baris adjustment masing-masing lini, di
-- transaksi kelahiran yang sama:
--   Prop  T_CLAIM_ADJUSTMENT.KOMITE_ID UNIQUE (UQ_CLAIM_ADJUSTMENT_KOMITE) + `sqlSetelKomiteAdjustment ... AND KOMITE_ID
--         IS NULL` + PastikanSatuBaris (claimprop/backend/repository/komite.go);
--   Life  T_CLAIMLF_ADJUSTMENT.KOMITE_ID UNIQUE (UX_ADJ_KOMITE_ID) + `PerbaruiKomiteID ... AND KOMITE_ID IS NULL` +
--         PastikanSatuBaris (claimlife/backend/repository/klaimlife.go), satu transaksi dengan BuatKasusKomite.
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
DROP INDEX {skema}.UX_GENERAL_KOMITE_ADJ
/
CREATE INDEX {skema}.IX_GENERAL_KOMITE_ADJ ON {skema}.T_GENERAL_KOMITE (ADJUSTMENT_ID)
/
