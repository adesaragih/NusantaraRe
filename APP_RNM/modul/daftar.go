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

	claimliferepository "nusantarare/internal/repository"
	"nusantarare/modul/premiumlist"
)

// SumberMigrasi mengembalikan folder migrasi SETIAP modul terdaftar.
//
// ⛔ Tidak bergantung pada modul yang aktif: skema selalu utuh, sebab tabel
// satu modul dapat dirujuk tabel modul lain (dan data warisan tidak memilih
// modul). Pelari mengurutkan langkah menurut NAMA berkas, lintas sumber.
func SumberMigrasi() []fs.FS {
	return []fs.FS{
		claimliferepository.SumberMigrasi(), // 001-049 (Claim Life, Komite)
		premiumlist.SumberMigrasi(),         // 050-079
	}
}
