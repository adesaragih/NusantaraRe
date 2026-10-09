-- 949 - baris menu modul Plan (`planlife`), modul DI LUAR dua puluh folder korpus.
--
-- Keputusan work owner 08-10-2026 K6: KODE / MODUL `planlife`, LABEL 'Plan', GROUPMENU MASTER TREATY, URUTAN 13
-- (sesudah Benefit 12), DIMIGRASI '1' langsung (modul tanpa migrasi sendiri), STATUS_AKTIF bawaan tabel - sama seperti
-- baris benefitlife 945. Panduan: section Pega InboxProductType (kelas ASM-FW-GISFW-Int-PRODUCT_TYPE_LIFE).
-- LABEL berbeda dari nama modul (bukan "Plan Life"): penjaga `labelTampilDisetujui` + `modulLuarKorpus`
-- (inti/backend/penjaga/menu_test.go) dan `MODUL_LUAR_KORPUS` (frontend/katalogKorpus.ts).
-- Akses: hak lewat LANGKAH-WO-PLANLIFE.md (d) (K7). NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'planlife', 'Plan', 'MASTER TREATY', 'planlife', 13, '1' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'planlife')
/
