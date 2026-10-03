-- 906 - baris menu modul Marketing Officer (`marketingofficer`), modul DI LUAR dua puluh folder korpus.
--
-- Keputusan work owner 03-10-2026: modul baru `marketingofficer` (insert/update POOLDATA.MARKETINGOFFICER), label
-- "Marketing Officer", kelompok MASTER. PANDUAN-TIM-PER-MODUL bab 5: isi awal 900 hanya dua puluh folder korpus dan
-- slot menu modul tidak boleh membuat kelompok, jadi baris modul di luar korpus dibuat langkah inti tersendiri.
-- Bentuk DATAR (sesudah 901: tanpa PARENT_ID); DIMIGRASI '0' - slot menu modulnya (990) yang menyalakannya.
-- Penjaga: `modulLuarKorpus` (inti/backend/penjaga/menu_test.go) dan `MODUL_LUAR_KORPUS` (frontend/katalogKorpus.ts).
--
-- Akses: akun yang sudah ada TIDAK otomatis mendapat menu ini (903 hanya isi awal); admin mencentangnya di Kelola User.
-- ⛔ NOL `COMMIT` (ADR-U-0029). `-migrate` dijalankan work owner.
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'marketingofficer', 'Marketing Officer', 'MASTER', 'marketingofficer', 4, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'marketingofficer')
/
