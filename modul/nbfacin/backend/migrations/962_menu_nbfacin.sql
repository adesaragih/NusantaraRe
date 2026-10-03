-- 962 - menu NB FacIn: modul ini mendapat layar pertamanya (tiket 21, keputusan work owner butir 59).
UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1', TGL_UBAH = SYSDATE
WHERE KODE = 'nbfacin'
/
