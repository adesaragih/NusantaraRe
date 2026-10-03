-- 480 - tiket 00 Endorsement Life, ralat E1 brief gelombang 2 (01-10-2026).
--
-- ⛔ NOL KOLOM BARU. Kolom EDM `T_PREMIUM_LIST` (`EDM_TYPE` ... `STATUS_OLD`,
-- `PROD_KE`) sudah dibuat migrasi 051 PremiumList Life, seluruhnya nullable;
-- `PARENT_ID` + `FK_PLD_PARENT` + `IDX_PLD_PARENT` sudah dibuat 052. Tabel
-- aplikasi milik kita, bukan warisan.
--
-- ⛔ `PROD_KE NUMBER(5) DEFAULT 1` supaya kode PremiumList TIDAK diubah: baris
-- new business lahir lewat `INSERT INTO T_PREMIUM_LIST (ID, ID_PEGA, TGL_INPUT)`
-- (`premiumlistlife/backend/repository/polis_kasus.go` `sqlSisipPremiumListKosong`)
-- yang tidak menyebut `PROD_KE`, sehingga ia otomatis versi 1. Endorsement
-- pertama menjadi versi 2 dan bernomor `<polis>/02` (RALAT R23, OQ-EDM-008).
-- Tipenya diulang (`NUMBER(5)`) supaya penjaga golongan tipe membaca bentuk
-- akhirnya, bukan kata `DEFAULT`.
--
-- ⛔ Isi mundur IDEMPOTEN: hanya baris yang `PROD_KE`-nya kosong, ke 1. Di DEV
-- tabel ini 0 baris (brief gelombang 2 §5.1). Jalur mundur tidak mengosongkan
-- kembali - nilai 1 sama artinya dengan "versi new business".
--
-- ⚠️ Tiga index pencari versi, bukan kolom (RALAT R21):
--   IDX_PL_NOPOLIS_PRODKE  STRUKTUR PremiumList bab T_PREMIUM_LIST "Index:
--                          NO_POLIS · PROD_KE" - belum dibuat 051. Versi EDM
--                          resmi dicari `NO_POLIS = :1 ORDER BY PROD_KE DESC`.
--   IDX_PL_OLD_POLICY_NO   gerbang 3 "satu endorsement terbuka per polis"
--                          (`FilterProteksiEDMLife` b541 `.PolicyNo`).
--   IDX_PLD_PL_NUMBER      versi new business dicari lewat nomor PL pesertanya:
--                          nomor polis = PL_NUMBER (`InsertJsonPolisEDM` b102).
--
-- ⛔ NOL `COMMIT` (ADR-U-0029). Batas transaksi milik pelari migrasi.
ALTER TABLE {skema}.T_PREMIUM_LIST MODIFY (
  PROD_KE NUMBER(5) DEFAULT 1
)
/
UPDATE {skema}.T_PREMIUM_LIST SET PROD_KE = 1 WHERE PROD_KE IS NULL
/
CREATE INDEX {skema}.IDX_PL_NOPOLIS_PRODKE ON {skema}.T_PREMIUM_LIST (NO_POLIS, PROD_KE)
/
CREATE INDEX {skema}.IDX_PL_OLD_POLICY_NO ON {skema}.T_PREMIUM_LIST (OLD_POLICY_NO)
/
CREATE INDEX {skema}.IDX_PLD_PL_NUMBER ON {skema}.T_PREMIUM_LIST_DETAIL (PL_NUMBER)
/
