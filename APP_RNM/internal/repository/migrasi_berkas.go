package repository

// Berkas migrasi modul ini, ditanam ke dalam biner.
//
// Refactor bentuk B (30-09-2026): penanaman dipisah dari pelarinya. Pelari
// (migrasi.go) menerima sumber dari luar; berkas ini hanya menyerahkan folder
// `migrations/` di sebelahnya. Nama berkas tidak berubah.

import (
	"embed"
	"io/fs"
)

//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// SumberMigrasi menyerahkan folder `migrations/` paket ini kepada pelari
// migrasi.
func SumberMigrasi() fs.FS { return berkasMigrasi }
