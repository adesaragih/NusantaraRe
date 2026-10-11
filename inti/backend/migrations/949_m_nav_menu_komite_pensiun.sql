-- 949 (langkah pensiun menu komite) - menu Komite Claim FacIn / Life / Non Prop / Prop DIHAPUS (perintah work owner
-- 09-10-2026: "kode menu nya di hapus dari repo, anggap menu itu tidak pernah ada, karena digabung di menu klaim nya
-- masing-masing"). Baris isi awal 900 (KLAIM URUTAN 5-8, yang terakhir - golongan KLAIM tetap rapat 1-4) dan hak isi
-- awal 903 dibuang. Isi awal 900 / 903 tidak disunting. Hak TIDAK disalin ke menu klaim: pemegang menu komite yang
-- tidak memegang menu klaimnya diberi hak lewat Kelola User (DEV 10-10-2026: nol akun seperti itu).
--
-- Modul komite yang sudah punya backend tetap ada sebagai modul TANPA MENU (frontend `layar.ts`): dipasang bagi pemegang
-- menu klaim peminjamnya (`MODUL_DIPINJAM` frontend/App.tsx) dan rutenya dipinjam (`ruteDipinjam` cmd/api/rakit.go).
-- Diterapkan skema tiruan penjaga menu (`langkahPensiunMenu`). Idempoten. NOL COMMIT. -migrate dijalankan work owner.
DELETE FROM {skema}.M_LOGIN_GO_MENU WHERE MENU_KODE IN ('komiteclaimfacin', 'komiteclaimlife', 'komiteclaimnonprop', 'komiteclaimprop')
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'komiteclaimfacin'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'komiteclaimlife'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'komiteclaimnonprop'
/
DELETE FROM {skema}.M_NAV_MENU WHERE KODE = 'komiteclaimprop'
/
