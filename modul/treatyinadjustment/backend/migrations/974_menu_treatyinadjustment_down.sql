-- 974 mundur - menu Treaty In Adjustment kembali "belum dimigrasi".
UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '0', TGL_UBAH = SYSDATE
WHERE KODE = 'treatyinadjustment'
/
