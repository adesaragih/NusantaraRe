-- 955 - baris menu modul Cause Of Loss Life (`causeoflosslife`), modul DI LUAR dua puluh folder korpus, lahir di SLOT
-- MENU modul ini sendiri (K0: slot 955 dipinjam dari jatah premiumlistlife).
--
-- Keputusan work owner 08-10-2026 K5: KODE / MODUL `causeoflosslife`, LABEL 'Cause Of Loss Life', GROUPMENU MASTER
-- TREATY, URUTAN 14 (sesudah Plan 13), DIMIGRASI '1' langsung, STATUS_AKTIF bawaan tabel - sama seperti baris planlife
-- 949. Panduan: section Pega InboxCauseofLossLife (kelas ASM-FW-GISFW-Int-CAUSEOFLOSS_LIFE).
-- Slot menu berjalan SESUDAH 900 / 901 / 909 (nomor tiga digit, pelari mengurut nama berkas), jadi M_NAV_MENU sudah
-- datar dan CHECK GROUPMENU sudah memuat MASTER TREATY. Penjaga: `modulLuarKorpus` + `barisLahirDiSlot`
-- (inti/backend/penjaga) dan `MODUL_LUAR_KORPUS` (frontend/katalogKorpus.ts).
-- Akses: akun yang sudah ada TIDAK otomatis mendapat menu ini - hak lewat LANGKAH-WO-CAUSEOFLOSSLIFE.md (d) (K6).
-- Aman diulang (NOT EXISTS). NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'causeoflosslife', 'Cause Of Loss Life', 'MASTER TREATY', 'causeoflosslife', 14, '1' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'causeoflosslife')
/
