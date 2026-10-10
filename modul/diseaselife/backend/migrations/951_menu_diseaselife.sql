-- 951 - baris menu modul Disease Life (`diseaselife`), modul DI LUAR dua puluh folder korpus, lahir di SLOT MENU modul
-- ini sendiri (K0: slot 951 dipinjam dari jatah claimlife).
--
-- Keputusan work owner 08-10-2026 D4: KODE / MODUL `diseaselife`, LABEL 'Disease Life', GROUPMENU MASTER TREATY,
-- URUTAN 15 (sesudah Cause Of Loss Life 14), DIMIGRASI '1' langsung, STATUS_AKTIF bawaan tabel - sama seperti baris
-- causeoflosslife 955. Panduan: section Pega InboxDisease (kelas ASM-FW-GISFW-Int-DISEASE_LIFE).
-- Slot menu berjalan SESUDAH 900 / 901 / 909 (nomor tiga digit, pelari mengurut nama berkas), jadi M_NAV_MENU sudah
-- datar dan CHECK GROUPMENU sudah memuat MASTER TREATY. Penjaga: `modulLuarKorpus` + `barisLahirDiSlot`
-- (inti/backend/penjaga) dan `MODUL_LUAR_KORPUS` (frontend/katalogKorpus.ts).
-- Akses: akun yang sudah ada TIDAK otomatis mendapat menu ini - hak lewat LANGKAH-WO-DISEASELIFE.md (d) (K5).
-- Aman diulang (NOT EXISTS). NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'diseaselife', 'Disease Life', 'MASTER TREATY', 'diseaselife', 15, '1' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'diseaselife')
/
