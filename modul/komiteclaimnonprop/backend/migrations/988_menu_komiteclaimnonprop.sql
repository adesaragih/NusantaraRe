-- 988 - menu Komite Claim Non Prop: modul ini mendapat layarnya (layar keputusan komite KMTNP-, dibuka dari inbox Claim
-- Non Prop - perintah work owner 09-10-2026). Baris modulnya sudah ada sejak 900 (menu datar 901): nol INSERT, satu
-- UPDATE DIMIGRASI. Menu Komite Claim Non Prop tidak dipakai (pola Komite Claim Prop): disembunyikan lewat DATA
-- (M_NAV_MENU.STATUS_AKTIF dan kode menu akun), langkah work owner - bukan migrasi slot menu.
UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1', TGL_UBAH = SYSDATE
WHERE KODE = 'komiteclaimnonprop'
/
