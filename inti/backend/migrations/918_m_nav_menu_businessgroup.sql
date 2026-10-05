-- 918 - baris menu modul Business Group (`businessgroup`), modul DI LUAR dua puluh folder korpus.
--
-- Perintah work owner 05-10-2026: menu CRUD tabel warisan POOLDATA.BUSINESSGROUP, modul sendiri di luar inti:
-- "BUSINESSGROUP buat juga crud nya skalian bisa? modul baru".
-- Golongan MASTER, label "Business Group". Pega hanya membaca tabel ini - tidak ada layar master di korpus.
-- PANDUAN-TIM-PER-MODUL bab 5: baris modul di luar korpus dibuat langkah inti tersendiri. Bentuk DATAR
-- (sesudah 901); URUTAN 9 - urutan rapat 1..n menu bersih MASTER, sesudah Adjuster Consultant (6).
-- DIMIGRASI '0' - slot menu modulnya (991) yang menyalakannya. Penjaga: `modulLuarKorpus`
-- (inti/backend/penjaga/menu_test.go) dan `MODUL_LUAR_KORPUS` (frontend/katalogKorpus.ts).
--
-- Akses: akun yang sudah ada TIDAK otomatis mendapat menu ini; admin mencentangnya di Kelola User.
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'businessgroup', 'Business Group', 'MASTER', 'businessgroup', 9, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'businessgroup')
/
