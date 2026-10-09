-- 941 - baris menu modul R/I Risk (`ririsklife`), modul DI LUAR dua puluh folder korpus.
--
-- Keputusan work owner 08-10-2026 K4: KODE / MODUL `ririsklife`, LABEL 'R/I Risk', GROUPMENU MASTER TREATY, URUTAN 11
-- (sesudah R/I Comm Life 10), DIMIGRASI '1' langsung (modul tanpa migrasi sendiri), STATUS_AKTIF bawaan tabel - sama
-- seperti baris ricommlife 925. Panduan: section Pega InboxSummaryRIRisk (kelas ASM-FW-GISFW-Int-RIRISK_LIFE_SUMMARY).
-- LABEL berbeda dari nama modul (bukan "R/I Risk Life"): penjaga `labelTampilDisetujui` + `modulLuarKorpus`
-- (inti/backend/penjaga/menu_test.go) dan `MODUL_LUAR_KORPUS` (frontend/katalogKorpus.ts).
-- Akses: akun yang sudah ada TIDAK otomatis mendapat menu ini - hak SUPERADMIN lewat LANGKAH-WO-RIRISKLIFE.md (d).
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'ririsklife', 'R/I Risk', 'MASTER TREATY', 'ririsklife', 11, '1' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'ririsklife')
/
