-- 957 - baris menu modul Cover Life (`coverlife`), modul DI LUAR dua puluh folder korpus, lahir di SLOT MENU modul ini
-- sendiri (K0: slot 957 dipinjam dari jatah treatycontractout).
--
-- Keputusan work owner 08-10-2026 C4: KODE / MODUL `coverlife`, LABEL 'Cover Life', GROUPMENU MASTER TREATY, URUTAN 16
-- (sesudah Disease Life 15), DIMIGRASI '1' langsung, STATUS_AKTIF bawaan tabel - sama seperti baris causeoflosslife
-- 955. Panduan: section Pega InboxCoverLife (kelas ASM-FW-GISFW-Int-COVER_LIFE).
-- Slot menu berjalan SESUDAH 900 / 901 / 909 (nomor tiga digit, pelari mengurut nama berkas), jadi M_NAV_MENU sudah
-- datar dan CHECK GROUPMENU sudah memuat MASTER TREATY. Penjaga: `modulLuarKorpus` + `barisLahirDiSlot`
-- (inti/backend/penjaga) dan `MODUL_LUAR_KORPUS` (frontend/katalogKorpus.ts).
-- Akses: akun yang sudah ada TIDAK otomatis mendapat menu ini - hak lewat LANGKAH-WO-COVERLIFE.md (d) (K5).
-- Aman diulang (NOT EXISTS). NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'coverlife', 'Cover Life', 'MASTER TREATY', 'coverlife', 16, '1' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'coverlife')
/
