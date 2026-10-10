-- 945 - baris menu modul Benefit (`benefitlife`), modul DI LUAR dua puluh folder korpus.
--
-- Keputusan work owner 08-10-2026 K4: KODE / MODUL `benefitlife`, LABEL 'Benefit', GROUPMENU MASTER TREATY, URUTAN 12
-- (sesudah R/I Risk 11), DIMIGRASI '1' langsung (modul tanpa migrasi sendiri), STATUS_AKTIF bawaan tabel - sama seperti
-- baris ririsklife 941. Panduan: section Pega InboxBenefit (kelas ASM-FW-GISFW-Int-BENEFIT_LIFE).
-- LABEL berbeda dari nama modul (bukan "Benefit Life"): penjaga `labelTampilDisetujui` + `modulLuarKorpus`
-- (inti/backend/penjaga/menu_test.go) dan `MODUL_LUAR_KORPUS` (frontend/katalogKorpus.ts).
-- Akses: akun yang sudah ada TIDAK otomatis mendapat menu ini - hak lewat LANGKAH-WO-BENEFITLIFE.md (d) (K5).
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'benefitlife', 'Benefit', 'MASTER TREATY', 'benefitlife', 12, '1' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'benefitlife')
/
