-- 760 - perbaiki MARKETINGOFFICER_LOG untuk log perubahan anggota (izin work owner 03-10-2026: "kamu boleh perbaiki
-- MARKETINGOFFICER_LOG jika dibutuhkan").
--
-- Sebelumnya: trigger warisan TRG_MARKETINGOFFICER_LOG mencatat baris LAMA setiap UPDATE MARKETINGOFFICER, tetapi
-- TANPA AKSES_LOGIN (pergantian akun login tidak terekam) dan TANPA waktu log (DEV: 35 baris, TANGGAL hanya 3 terisi,
-- ACTION selalu kosong) - urutan riwayat tidak dapat ditentukan. Trigger TIDAK disentuh. Yang ditambah:
--   - LOG_TIME TIMESTAMP(6): waktu UPDATE yang menghasilkan baris log. Ditambah TANPA default (baris log lama tetap
--     kosong - ditandai "urutan perkiraan" di layar), lalu DEFAULT SYSTIMESTAMP untuk baris BARU, sehingga setiap
--     baris yang trigger sisipkan sesudah ini - juga dari Pega - bercap waktu.
--   - AKSES_LOGIN VARCHAR2(50): AKSES_LOGIN lama. Trigger tidak mengisinya; aplikasi Go mengisinya pada baris log
--     yang baru disisipkan trigger, di transaksi UPDATE yang sama, dan menandainya ACTION = 'UPDATE-GO'. Baris tanpa
--     tanda itu (Pega, atau lama) = AKSES_LOGIN tidak diketahui.
--
-- Skema uji (uji/skemauji) membuat tiruan tabel ini SEBELUM migrasi (`ddlTiruanMarketingOfficer`).
-- ⛔ NOL `COMMIT` (ADR-U-0029). `-migrate` dijalankan work owner.
ALTER TABLE {skema}.MARKETINGOFFICER_LOG ADD (
  LOG_TIME    TIMESTAMP(6),
  AKSES_LOGIN VARCHAR2(50)
)
/
ALTER TABLE {skema}.MARKETINGOFFICER_LOG MODIFY (LOG_TIME DEFAULT SYSTIMESTAMP)
/
