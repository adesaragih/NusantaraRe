// ⚠️ LOKAL SAJA - JANGAN DI-COMMIT. Hapus berkas ini bila M_NAV_MENU sudah dirapikan.
//
// M_NAV_MENU di database memuat GROUPMENU 'MASTER TREATY' (diubah di luar migrasi),
// yang tidak dikenal `Golongan`, sehingga modulnya hilang dari sidebar. Berkas ini
// menambah golongan itu di proses lokal tanpa mengubah database.

package menu

func init() {
	Golongan = append(Golongan, "MASTER TREATY")
}
