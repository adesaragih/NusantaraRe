-- 912 - baris M_NAV_MENU modul luar korpus `masternation` ("Nation", grup MASTER, URUTAN 6).
--
-- Keputusan work owner 04-10-2026 (diteruskan sesi nusantarare-9d): Master Data dipecah menjadi delapan modul menu
-- terpisah ("8 modul terpisah"), nama "Nation"; modul luar korpus TANPA slot menu ("Modul luar korpus tanpa slot") -
-- karena itu barisnya lahir DIMIGRASI '1' di sini, bukan dinyalakan slot. Bentuk datar pola 906-908 (sesudah 901,
-- tanpa PARENT_ID), idempoten lewat NOT EXISTS. Hak menu: 920 menyalinnya dari pemegang menu Master Data; akun lain
-- lewat Kelola User. Ditulis, TIDAK dijalankan agent. Nol COMMIT (ADR-U-0029).
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'masternation', 'Nation', 'MASTER', 'masternation', 6, '1' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'masternation')
/
