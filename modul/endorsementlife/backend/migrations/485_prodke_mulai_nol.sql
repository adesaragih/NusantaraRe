-- 485 - PRODKE SATU hitungan mulai 0 (keputusan work owner 08-10-2026).
--
-- Aturan bisnis: PRODKE = berapa kali PL Number diubah lewat endorsement. Polis baru (PremiumList, `NBLF-`)
-- = 0; endorsement menarik versi ber-PRODKE terbesar lalu Confirm menulis versi baru = terbesar + 1;
-- NO_ENDORS = `<PL Number>/<PRODKE dua digit>`. Dulu (480, E1) versi 1-based: `PROD_KE DEFAULT 1`, NB = 1,
-- endorsement pertama = 2, nomor dari penghitung terpisah (K3 `UrutanPega`).
--
-- 1. `PROD_KE NUMBER(5) DEFAULT 0` - baris new business PremiumList lahir lewat INSERT yang tidak menyebut
--    `PROD_KE` (`premiumlistlife/backend/repository/polis_kasus.go`), jadi kode PremiumList TIDAK diubah.
--    Tipe diulang (`NUMBER(5)`) supaya penjaga golongan tipe membaca bentuk akhirnya.
-- 2. Polis baru PremiumList yang sudah ada: `NBLF-%` ber-PROD_KE 1 (isi mundur 480) → 0.
-- 3. Kasus Endorsement Life yang sudah ada (`EDML-%`): PROD_KE skema lama → skema baru (`- 1`). Yang tersimpan
--    di sana SELALU angka skema lama - kasus terbuka/ditolak: versi RENCANA saat BuatKasus (versi berjalan
--    1-based + 1, jadi >= 2); kasus Resolved-Completed: angka yang sama, ditetapkan ulang saat Confirm.
--    `>= 1` menjaga baris rusak / kosong tidak menjadi negatif. NO_ENDORS yang sudah terbit TIDAK diubah.
--
-- ⚠️ Pembaca lain `PROD_KE`: Claim Life (`claimlife/backend/repository/pohonklaim.go`) hanya mengurutkan
-- `NVL(PROD_KE, 0) DESC` - urutan relatif tidak berubah oleh pergeseran ini.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029). Batas transaksi milik pelari migrasi.
ALTER TABLE {skema}.T_PREMIUM_LIST MODIFY (
  PROD_KE NUMBER(5) DEFAULT 0
)
/
UPDATE {skema}.T_PREMIUM_LIST SET PROD_KE = 0 WHERE ID LIKE 'NBLF-%' AND PROD_KE = 1
/
UPDATE {skema}.T_PREMIUM_LIST SET PROD_KE = PROD_KE - 1 WHERE ID LIKE 'EDML-%' AND PROD_KE >= 1
/
