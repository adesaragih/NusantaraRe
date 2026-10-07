-- 920 mundur - baris 'masterdata' kembali (bentuk 911 lama, URUTAN 6) dengan DIMIGRASI '0': kode modul masterdata
-- sudah tidak ada, jadi baris itu TIDAK boleh tampil sebagai menu hidup; menyalakannya = keputusan bila modulnya
-- dikembalikan. ⚠️ Di skema yang tidak pernah menjalankan 911 (dibuang), jalur mundur ini tetap menyisipkan barisnya -
-- mati (DIMIGRASI '0'), tidak dibuang langkah mundur lain. Haknya diberikan kepada setiap akun yang memegang SALAH SATU
-- dari delapan menu master - pendekatan [dugaan], RUGI-INFORMASI: siapa yang dulu memegang 'masterdata' tidak tersimpan
-- setelah 920.
-- Berjalan SEBELUM jalur mundur 919-912 (urutan menurun), jadi hak delapan menu masih ada saat dibaca.
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'masterdata', 'Master Data', 'MASTER', 'masterdata', 6, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'masterdata')
/
INSERT INTO {skema}.M_LOGIN_GO_MENU (LOGIN_ID, MENU_KODE)
SELECT DISTINCT g.LOGIN_ID, 'masterdata'
  FROM {skema}.M_LOGIN_GO_MENU g
 WHERE g.MENU_KODE IN ('masternation', 'masterprovince', 'mastercity', 'masterdistrict', 'masterczone', 'masteraccumulatedtype', 'masteraccumulation', 'masterobjectitemtype')
   AND NOT EXISTS (SELECT 1 FROM {skema}.M_LOGIN_GO_MENU x WHERE x.LOGIN_ID = g.LOGIN_ID AND x.MENU_KODE = 'masterdata')
/
