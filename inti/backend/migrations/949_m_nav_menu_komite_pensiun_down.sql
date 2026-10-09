-- 949 mundur - baris menu keempat komite kembali seperti isi awal 900 (KLAIM URUTAN 5-8, DIMIGRASI isi awal; slot
-- menu 986 / 988 yang dulu menyalakannya sudah dihapus). Hak akun TIDAK dikembalikan - berikan lewat Kelola User.
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'komiteclaimfacin', 'Komite Claim FacIn', 'KLAIM', 'komiteclaimfacin', 5, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'komiteclaimfacin')
/
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'komiteclaimlife', 'Komite Claim Life', 'KLAIM', 'komiteclaimlife', 6, '1' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'komiteclaimlife')
/
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'komiteclaimnonprop', 'Komite Claim Non Prop', 'KLAIM', 'komiteclaimnonprop', 7, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'komiteclaimnonprop')
/
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'komiteclaimprop', 'Komite Claim Prop', 'KLAIM', 'komiteclaimprop', 8, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'komiteclaimprop')
/
