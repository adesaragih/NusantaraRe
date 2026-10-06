-- 925 - baris menu modul R/I Comm Life (`ricommlife`), modul DI LUAR dua puluh folder korpus.
--
-- Perintah work owner 06-10-2026: "membuat modul baru di Master Treaty dengan nama R/I Comm Life, panduannya baca dari
-- xml di D:\NUSARE DEV\Menu RI Comm, konsepnya hampir sama dengan menu R/I Rate, membuat CRUD, dan detail bisa di save
-- dan edit" (section Pega InboxSummaryRIComm, kelas ASM-FW-GISFW-Int-RICOMM_LIFE_SUMMARY). Bukti struktur menu Pega
-- (portal/navigasi) TIDAK ditemukan di repo maupun di kedua XML - dasarnya perintah work owner itu (keputusan work
-- owner 06-10-2026 butir 8). PANDUAN-TIM-PER-MODUL bab 5: baris modul di luar korpus dibuat langkah inti tersendiri.
-- Bentuk DATAR (sesudah 901); URUTAN 10 - urutan rapat MASTER TREATY, sesudah R/I Rate Life (9).
-- DIMIGRASI '1' LANGSUNG: modul ini TANPA migrasi sendiri (MODUL.md `—`). Penjaga: `modulLuarKorpus`
-- (inti/backend/penjaga/menu_test.go) dan `MODUL_LUAR_KORPUS` (frontend/katalogKorpus.ts).
-- ⚠️ Berjalan SESUDAH 924 (tabel flat RICOMM_LIFE) - yang menuntut VIEW RICOMM_LIFE sudah dibuang DBA lebih dulu.
--
-- Akses: akun yang sudah ada TIDAK otomatis mendapat menu ini; admin mencentangnya di Kelola User.
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'ricommlife', 'R/I Comm Life', 'MASTER TREATY', 'ricommlife', 10, '1' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'ricommlife')
/
