-- 976 mundur - menu Endorsement Life kembali "belum dimigrasi".
UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '0', TGL_UBAH = SYSDATE
WHERE KODE = 'endorsementlife'
/
