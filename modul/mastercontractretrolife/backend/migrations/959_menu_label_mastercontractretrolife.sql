-- 959 - nama tampilan menu "Contract Retro Life": kata "Master" dihapus (keputusan work owner 03-10-2026, nama
-- tampilan saja). KODE, MODUL, folder, rute, dan MODUL_AKTIF tetap `mastercontractretrolife`.
UPDATE {skema}.M_NAV_MENU SET LABEL = 'Contract Retro Life', TGL_UBAH = SYSDATE
WHERE KODE = 'mastercontractretrolife'
/
