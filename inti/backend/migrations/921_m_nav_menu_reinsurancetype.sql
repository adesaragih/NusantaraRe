-- 921 - baris menu modul Reinsurance Type (`reinsurancetype`), modul DI LUAR dua puluh folder korpus.
--
-- Perintah work owner 05-10-2026: menu CRUD tabel warisan POOLDATA.REINSURANCETYPE (master jenis reasuransi), modul
-- sendiri: "select * from reinsurancetype"; golongan MASTER TREATY, label "Reinsurance Type". Pega hanya membaca tabel
-- ini - tidak ada layar master di korpus. PANDUAN-TIM-PER-MODUL bab 5: baris modul di luar korpus dibuat langkah inti
-- tersendiri. Bentuk DATAR (sesudah 901); URUTAN 8 - urutan rapat 1..n menu bersih MASTER TREATY, sesudah Treaty
-- Description (7).
-- DIMIGRASI '1' LANGSUNG: modul ini TANPA migrasi sendiri (MODUL.md `—`) - seluruh nomor modul sudah terbagi dan modul
-- lain tidak boleh disentuh. Penjaga: `modulLuarKorpus` (inti/backend/penjaga/menu_test.go) dan `MODUL_LUAR_KORPUS`
-- (frontend/katalogKorpus.ts).
--
-- Akses: akun yang sudah ada TIDAK otomatis mendapat menu ini; admin mencentangnya di Kelola User.
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'reinsurancetype', 'Reinsurance Type', 'MASTER TREATY', 'reinsurancetype', 8, '1' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'reinsurancetype')
/
