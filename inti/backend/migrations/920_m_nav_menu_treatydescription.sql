-- 920 - baris menu modul Treaty Description (`treatydescription`), modul DI LUAR dua puluh folder korpus.
--
-- Perintah work owner 05-10-2026: CRUD tabel warisan POOLDATA.TREATYDESC sebagai modul baru ("INI MODUL BARU,
-- SELECT * FROM TREATYDESC, PAHAMI TABEL ITU, AKU MAU BUAT CRUD"; "NAMANYA TREATY DESCRIPTION YA"), golongan
-- MASTER TREATY. Layar input Pega-nya tidak ada di korpus. PANDUAN-TIM-PER-MODUL bab 5: baris modul di luar korpus
-- dibuat langkah inti tersendiri. Bentuk DATAR (sesudah 901); URUTAN 7 - urutan rapat 1..n MASTER TREATY, sesudah
-- Treaty Exchange Yearly (6, migrasi inti 919).
-- DIMIGRASI '1' langsung: modul TANPA migrasi sendiri (MODUL.md `—`, `tandaTanpaMigrasi`) - seluruh nomor modul sudah
-- terbagi dan modul lain tidak boleh disentuh ("JANGAN ADA SENTUH MODUL LAIN!"), seperti 919. Penjaga:
-- `modulLuarKorpus` (inti/backend/penjaga/menu_test.go) dan `MODUL_LUAR_KORPUS` (frontend/katalogKorpus.ts).
--
-- Akses: akun yang sudah ada TIDAK otomatis mendapat menu ini; admin mencentangnya di Kelola User.
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'treatydescription', 'Treaty Description', 'MASTER TREATY', 'treatydescription', 7, '1' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'treatydescription')
/
