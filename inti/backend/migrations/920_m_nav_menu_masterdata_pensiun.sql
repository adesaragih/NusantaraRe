-- 920 - menu Master Data pensiun: hak menunya disalin ke delapan modul master, lalu baris 'masterdata' dibuang.
--
-- Keputusan work owner 04-10-2026 (diteruskan sesi nusantarare-9d): Master Data dipecah menjadi delapan modul
-- (912-919); modul masterdata DIHAPUS ("Dihapus, tabel pindah ke Province"); hak akun "Ya, salin otomatis" - setiap
-- akun yang memegang menu 'masterdata' mendapat kedelapan menu baru. Berkas 911 (baris 'masterdata') dan 996 (slot-nya)
-- dibuang: di skema baru barisnya tidak pernah lahir dan langkah ini tidak mengubah apa pun; di DEV yang sudah
-- menjalankannya, langkah inilah yang membuang barisnya. Urutan: salin hak (sebelum baris lamanya dibuang), buang hak
-- lama, buang baris. Ketiganya aman diulang (NOT EXISTS; DELETE nol baris). Ditulis, TIDAK dijalankan agent. Nol COMMIT.
INSERT INTO {skema}.M_LOGIN_GO_MENU (LOGIN_ID, MENU_KODE)
SELECT g.LOGIN_ID, k.KODE
  FROM {skema}.M_LOGIN_GO_MENU g
 CROSS JOIN (
  SELECT 'masternation' AS KODE FROM DUAL UNION ALL
  SELECT 'masterprovince' FROM DUAL UNION ALL
  SELECT 'mastercity' FROM DUAL UNION ALL
  SELECT 'masterdistrict' FROM DUAL UNION ALL
  SELECT 'masterczone' FROM DUAL UNION ALL
  SELECT 'masteraccumulatedtype' FROM DUAL UNION ALL
  SELECT 'masteraccumulation' FROM DUAL UNION ALL
  SELECT 'masterobjectitemtype' FROM DUAL
 ) k
 WHERE g.MENU_KODE = 'masterdata'
   AND NOT EXISTS (SELECT 1 FROM {skema}.M_LOGIN_GO_MENU x WHERE x.LOGIN_ID = g.LOGIN_ID AND x.MENU_KODE = k.KODE)
/
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE = 'masterdata'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'masterdata'
/
