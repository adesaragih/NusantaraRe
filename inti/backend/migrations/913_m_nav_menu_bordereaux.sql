-- 913 - baris menu modul Bordereaux (`bordereaux`), modul DI LUAR dua puluh folder korpus.
--
-- Perintah work owner 04-10-2026: folder korpus `D:\XML\RNM_BRD\Bordereaux` sebagai modul/menu baru, kelompok
-- MASTER TREATY ("Master Treaty"). Padanan Pega: kelas ASM-FW-GISFW-Int-BORDEREAUX (PortalBordereaux,
-- InputBordereaux, UploadCSVBordereaux_Act, BdxSave_Act, ActionSubmit, ...). PANDUAN-TIM-PER-MODUL bab 5: baris modul
-- di luar korpus dibuat langkah inti tersendiri. Bentuk DATAR (sesudah 901); URUTAN 5 sesudah Aggregate.
-- DIMIGRASI '0' - slot menu modulnya (998) yang menyalakannya. Penjaga: `modulLuarKorpus`
-- (inti/backend/penjaga/menu_test.go) dan `MODUL_LUAR_KORPUS` (frontend/katalogKorpus.ts).
--
-- Akses: akun yang sudah ada TIDAK otomatis mendapat menu ini; admin mencentangnya di Kelola User.
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'bordereaux', 'Bordereaux', 'MASTER TREATY', 'bordereaux', 5, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'bordereaux')
/
