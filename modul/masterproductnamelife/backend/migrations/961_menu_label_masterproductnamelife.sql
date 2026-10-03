-- 961 - nama tampilan menu "Product Name Life": kata "Master" dihapus (keputusan work owner 03-10-2026, nama
-- tampilan saja). KODE, MODUL, folder, rute, dan MODUL_AKTIF tetap `masterproductnamelife`.
UPDATE {skema}.M_NAV_MENU SET LABEL = 'Product Name Life', TGL_UBAH = SYSDATE
WHERE KODE = 'masterproductnamelife'
/
