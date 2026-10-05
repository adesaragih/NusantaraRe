-- 919 - baris menu modul Treaty Exchange Yearly (`treatyexchangeyearly`), modul DI LUAR dua puluh folder korpus.
--
-- Perintah work owner 05-10-2026: menu CRUD tabel warisan POOLDATA.TREATYEXCHANGEYEARLY (kurs tahunan), modul sendiri:
-- "SELECT * FROM TREATYEXCHANGEYEARLY ... BUAT CRUD JUGA"; golongan MASTER TREATY, label "Treaty Exchange Yearly".
-- Pega hanya membaca tabel ini - tidak ada layar master di korpus. PANDUAN-TIM-PER-MODUL bab 5: baris modul di luar
-- korpus dibuat langkah inti tersendiri. Bentuk DATAR (sesudah 901); URUTAN 6 - urutan rapat 1..n menu bersih MASTER
-- TREATY, sesudah Bordereaux (5).
-- DIMIGRASI '1' LANGSUNG: modul ini TANPA migrasi sendiri (MODUL.md `—`) - seluruh nomor modul 001-899 dan slot
-- 950-999 sudah terbagi, dan modul lain tidak boleh disentuh ("JANGAN ADA SENTUH MODUL LAIN"). Penjaga:
-- `modulLuarKorpus` (inti/backend/penjaga/menu_test.go) dan `MODUL_LUAR_KORPUS` (frontend/katalogKorpus.ts).
--
-- Akses: akun yang sudah ada TIDAK otomatis mendapat menu ini; admin mencentangnya di Kelola User.
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'treatyexchangeyearly', 'Treaty Exchange Yearly', 'MASTER TREATY', 'treatyexchangeyearly', 6, '1' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'treatyexchangeyearly')
/
