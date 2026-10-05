-- 917 - baris menu modul Treaty Group (`treatygroup`), modul DI LUAR dua puluh folder korpus.
--
-- Perintah work owner 05-10-2026: menu CRUD tabel warisan POOLDATA.TREATYGROUP, modul sendiri di luar inti:
-- "select * from POOLDATA.TREATYGROUP ... coba kamu pahami".
-- Golongan MASTER, label "Treaty Group". Pega hanya membaca tabel ini - tidak ada layar master di korpus.
-- PANDUAN-TIM-PER-MODUL bab 5: baris modul di luar korpus dibuat langkah inti tersendiri. Bentuk DATAR
-- (sesudah 901); URUTAN 8 - urutan rapat 1..n menu bersih MASTER, sesudah Adjuster Consultant (6).
-- DIMIGRASI '0' - slot menu modulnya (993) yang menyalakannya. Penjaga: `modulLuarKorpus`
-- (inti/backend/penjaga/menu_test.go) dan `MODUL_LUAR_KORPUS` (frontend/katalogKorpus.ts).
--
-- Akses: akun yang sudah ada TIDAK otomatis mendapat menu ini; admin mencentangnya di Kelola User.
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'treatygroup', 'Treaty Group', 'MASTER', 'treatygroup', 8, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'treatygroup')
/
