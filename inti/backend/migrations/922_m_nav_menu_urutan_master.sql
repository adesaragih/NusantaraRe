-- 922 - URUTAN golongan MASTER dirapatkan sesudah merge origin/dev 05-10-2026.
--
-- Dua sisi menambah baris modul MASTER bersamaan dan nomornya bertumpuk: lokal 912-919 (delapan modul master, URUTAN
-- 6-13, keputusan work owner 04-10-2026 "8 modul terpisah") dan GitHub 915-918 (Adjuster Consultant, Treaty Group OJK,
-- Treaty Group, Business Group, URUTAN 6-9, keputusan work owner 05-10-2026). Berkas yang membuat baris-baris itu TIDAK
-- disunting - namanya kunci T_MIGRASI dan isinya mungkin sudah jalan di DEV. Langkah ini menyetel URUTAN 6..17 menurut
-- urutan berkas (urutan sisip), yang ditagih penjaga `TestMenuBersihDuaPuluhBarisSatuPerModul`; MASTER 1-5
-- (mastercontractretrolife, masterproductnamelife, marketingofficer, companydetail, accounts) tetap. Baris yang belum
-- ada: UPDATE nol baris. Bentuk = 909 (`langkahGolonganMenu`). NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
UPDATE {skema}.M_NAV_MENU SET GROUPMENU = 'MASTER', URUTAN = 6, TGL_UBAH = SYSDATE
WHERE KODE = 'masternation'
/
UPDATE {skema}.M_NAV_MENU SET GROUPMENU = 'MASTER', URUTAN = 7, TGL_UBAH = SYSDATE
WHERE KODE = 'masterprovince'
/
UPDATE {skema}.M_NAV_MENU SET GROUPMENU = 'MASTER', URUTAN = 8, TGL_UBAH = SYSDATE
WHERE KODE = 'mastercity'
/
UPDATE {skema}.M_NAV_MENU SET GROUPMENU = 'MASTER', URUTAN = 9, TGL_UBAH = SYSDATE
WHERE KODE = 'adjusterconsultant'
/
UPDATE {skema}.M_NAV_MENU SET GROUPMENU = 'MASTER', URUTAN = 10, TGL_UBAH = SYSDATE
WHERE KODE = 'masterdistrict'
/
UPDATE {skema}.M_NAV_MENU SET GROUPMENU = 'MASTER', URUTAN = 11, TGL_UBAH = SYSDATE
WHERE KODE = 'masterczone'
/
UPDATE {skema}.M_NAV_MENU SET GROUPMENU = 'MASTER', URUTAN = 12, TGL_UBAH = SYSDATE
WHERE KODE = 'treatygroupojk'
/
UPDATE {skema}.M_NAV_MENU SET GROUPMENU = 'MASTER', URUTAN = 13, TGL_UBAH = SYSDATE
WHERE KODE = 'masteraccumulatedtype'
/
UPDATE {skema}.M_NAV_MENU SET GROUPMENU = 'MASTER', URUTAN = 14, TGL_UBAH = SYSDATE
WHERE KODE = 'treatygroup'
/
UPDATE {skema}.M_NAV_MENU SET GROUPMENU = 'MASTER', URUTAN = 15, TGL_UBAH = SYSDATE
WHERE KODE = 'businessgroup'
/
UPDATE {skema}.M_NAV_MENU SET GROUPMENU = 'MASTER', URUTAN = 16, TGL_UBAH = SYSDATE
WHERE KODE = 'masteraccumulation'
/
UPDATE {skema}.M_NAV_MENU SET GROUPMENU = 'MASTER', URUTAN = 17, TGL_UBAH = SYSDATE
WHERE KODE = 'masterobjectitemtype'
/
