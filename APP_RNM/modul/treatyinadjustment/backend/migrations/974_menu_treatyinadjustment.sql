-- 974 - menu Treaty In Adjustment: modul ini mendapat layar pertamanya.
-- Baris modulnya sudah ada sejak 900 (menu datar 901): nol INSERT, satu UPDATE DIMIGRASI.
UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1', TGL_UBAH = SYSDATE
WHERE KODE = 'treatyinadjustment'
/
