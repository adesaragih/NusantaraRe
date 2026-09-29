// Package modul adalah DAFTAR modul aplikasi - satu-satunya paket yang
// mengenal semua modul sekaligus.
//
// Refactor bentuk B (30-09-2026). Modul di `modul/<nama>/` tidak pernah
// mengimpor modul lain; paket ini dan `cmd/api` yang merakitnya. Yang
// memakai daftar ini: `cmd/api` (flag -migrate / -migrate-down) dan skema uji
// test bertag `db`.
package modul

import (
	"io/fs"

	"nusantarare/internal/repository"
	"nusantarare/modul/komite"
	"nusantarare/modul/premiumlist"
)

// SumberMigrasi mengembalikan folder migrasi SETIAP modul terdaftar.
//
// ⛔ Tidak bergantung pada modul yang aktif: skema selalu utuh, sebab tabel
// satu modul dapat dirujuk tabel modul lain (dan data warisan tidak memilih
// modul). Pelari mengurutkan langkah menurut NAMA berkas, lintas sumber.
func SumberMigrasi() []fs.FS {
	return []fs.FS{
		repository.SumberMigrasi(),  // 001-029 (Claim Life)
		komite.SumberMigrasi(),      // 030-049
		premiumlist.SumberMigrasi(), // 050-079
	}
}
