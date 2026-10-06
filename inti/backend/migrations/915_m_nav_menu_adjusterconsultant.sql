-- 915 - baris menu modul Adjuster Consultant (`adjusterconsultant`), modul DI LUAR dua puluh folder korpus.
--
-- Perintah work owner 05-10-2026: menu CRUD tabel POOLDATA.ADJUSTERCONSULTANT, golongan MASTER, label "Adjuster
-- Consultant"; "baris menu tetap di inti, tapi adjusterconsultant jangan masukkan ke folder inti". Padanan Pega:
-- layar master `MstAdjusterConsultant` (folder korpus Claim Fac In dan Claim Prop). PANDUAN-TIM-PER-MODUL bab 5:
-- baris modul di luar korpus dibuat langkah inti tersendiri. Bentuk DATAR (sesudah 901); URUTAN 6 sesudah Accounts
-- (urutan rapat 1..n menu bersih; baris master* DEV di luar migrasi juga memakai 6-13, urutan seri menurut ID).
-- DIMIGRASI '0' - slot menu modulnya (997) yang menyalakannya. Penjaga: `modulLuarKorpus`
-- (inti/backend/penjaga/menu_test.go) dan `MODUL_LUAR_KORPUS` (frontend/katalogKorpus.ts).
--
-- Akses: akun yang sudah ada TIDAK otomatis mendapat menu ini; admin mencentangnya di Kelola User.
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'adjusterconsultant', 'Adjuster Consultant', 'MASTER', 'adjusterconsultant', 6, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'adjusterconsultant')
/
