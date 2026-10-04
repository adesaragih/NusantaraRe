-- 907 - baris menu modul Company Detail (`companydetail`), modul DI LUAR dua puluh folder korpus.
--
-- Perintah work owner 03/04-10-2026: modul baru `companydetail`, bentuk menu seperti layar Pega SFAGIS Company
-- Detail (tidak ada di korpus XML), label "Company Detail", kelompok MASTER (disetujui 04-10-2026). PANDUAN-TIM-PER-MODUL
-- bab 5: isi awal 900 hanya dua puluh folder korpus dan slot menu modul tidak boleh membuat kelompok, jadi baris modul
-- di luar korpus dibuat langkah inti tersendiri. Bentuk DATAR (sesudah 901: tanpa PARENT_ID); DIMIGRASI '0' - slot
-- menu modulnya (992) yang menyalakannya. Penjaga: `modulLuarKorpus` (inti/backend/penjaga/menu_test.go) dan
-- `MODUL_LUAR_KORPUS` (frontend/katalogKorpus.ts).
--
-- Akses: akun yang sudah ada TIDAK otomatis mendapat menu ini (903 hanya isi awal); admin mencentangnya di Kelola User.
-- ⛔ NOL `COMMIT` (ADR-U-0029). `-migrate` dijalankan work owner.
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'companydetail', 'Company Detail', 'MASTER', 'companydetail', 5, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'companydetail')
/
