-- 996 - menu Master Data: modul ini mendapat layar pertamanya.
-- Baris modulnya ada sejak migrasi inti 911 (modul di luar korpus): nol INSERT, satu UPDATE DIMIGRASI.
UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1', TGL_UBAH = SYSDATE
WHERE KODE = 'masterdata'
/
