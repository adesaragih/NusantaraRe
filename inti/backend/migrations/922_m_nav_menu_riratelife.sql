-- 922 - baris menu modul R/I Rate Life (`riratelife`), modul DI LUAR dua puluh folder korpus.
--
-- Perintah work owner 05-10-2026: "Buat Menu baru Namanya R/I Rate Life pada Master Treaty, menu ini bisa CRUD untuk
-- simpan data ke tabel RATE_LIFE_SUMMARY, panduannya xml yang saya berikan" (section Pega InboxSummaryRIRate, kelas
-- ASM-FW-GISFW-Int-RATE_LIFE_SUMMARY). PANDUAN-TIM-PER-MODUL bab 5: baris modul di luar korpus dibuat langkah inti
-- tersendiri. Bentuk DATAR (sesudah 901); URUTAN 9 - urutan rapat 1..n menu bersih MASTER TREATY, sesudah Reinsurance
-- Type (8).
-- DIMIGRASI '1' LANGSUNG: modul ini TANPA migrasi sendiri (MODUL.md `—`). Penjaga: `modulLuarKorpus`
-- (inti/backend/penjaga/menu_test.go) dan `MODUL_LUAR_KORPUS` (frontend/katalogKorpus.ts).
-- ⚠️ Sebelum menyalakan: WO/DBA memeriksa `SELECT TEXT FROM ALL_VIEWS WHERE VIEW_NAME='RATE_LIFE_SUMMARY'` (ASUMSI
-- kunci JSON M_RATE_LIFE_SUMMARY, modul/riratelife/MODUL.md A1) dan menjalankan 923 (sequence ID).
--
-- Akses: akun yang sudah ada TIDAK otomatis mendapat menu ini; admin mencentangnya di Kelola User.
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'riratelife', 'R/I Rate Life', 'MASTER TREATY', 'riratelife', 9, '1' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'riratelife')
/
