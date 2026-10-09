-- 988 mundur - menu Komite Claim Non Prop kembali "belum dimigrasi".
UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '0', TGL_UBAH = SYSDATE
WHERE KODE = 'komiteclaimnonprop'
/
