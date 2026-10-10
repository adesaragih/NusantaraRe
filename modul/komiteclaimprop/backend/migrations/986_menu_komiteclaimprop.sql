-- 986 - menu Komite Claim Prop: modul ini mendapat layar pertamanya (daftar kerja penyetuju, worklist KomiteRouter).
-- Baris modulnya sudah ada sejak 900 (menu datar 901): nol INSERT, satu UPDATE DIMIGRASI.
UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1', TGL_UBAH = SYSDATE
WHERE KODE = 'komiteclaimprop'
/
