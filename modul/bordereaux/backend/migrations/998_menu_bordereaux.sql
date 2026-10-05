-- 998 - menu Bordereaux: modul ini mendapat layar pertamanya.
-- Barisnya dibuat migrasi inti 913 (modul di luar korpus, panduan bab 5): nol INSERT di sini, satu UPDATE DIMIGRASI.
UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1', TGL_UBAH = SYSDATE
WHERE KODE = 'bordereaux'
/
