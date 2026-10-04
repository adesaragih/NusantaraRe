-- 904 - baris M_NAV_MENU modul DI LUAR dua puluh folder korpus: `masterdata` (Master Data).
--
-- Keputusan work owner 04-10-2026 (diteruskan sesi nusantarare-0f, `modul/masterdata/docs/issues/01-rencana-master-data.md`):
-- penempatan "Modul baru 'masterdata'". PANDUAN-TIM-PER-MODUL bab 5: tim inti menambah baris M_NAV_MENU modul di luar
-- korpus (isi awal 900 hanya dua puluh folder korpus; slot menu modul tidak boleh membuat baris). 900 / 901 yang sudah
-- terpasang TIDAK disunting.
--
-- Bentuk SESUDAH 901 (tanpa PARENT_ID): satu INSERT idempoten `WHERE NOT EXISTS`. GROUPMENU MASTER, URUTAN 4 (sesudah
-- mastercontractretrolife 1, masterproductnamelife 2, treatycontractout 3), DIMIGRASI '0' - slot menu modul (990) yang menyalakannya saat
-- layarnya ada. Akses per akun (M_LOGIN_GO_MENU) TIDAK diisi: isi awal 903 hanya untuk menu yang ada saat itu; menu ini
-- diberikan lewat Kelola User (M-6). Penjaga: `inti/backend/penjaga/menu_test.go` (skema tiruan menerapkannya) dan
-- `rentang_test.go`. Nol COMMIT.
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'masterdata', 'Master Data', 'MASTER', 'masterdata', 4, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'masterdata')
/
