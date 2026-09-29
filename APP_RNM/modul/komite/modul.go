// Package komite merakit modul Komite Claim Life.
//
// Refactor bentuk B (30-09-2026): setiap modul punya satu berkas perakitan
// (`modul.go`) yang menyerahkan miliknya kepada `cmd/api`. Di paket ini
// dimulai dari migrasinya; rute dan pekerja latar menyusul (paket 6).
package komite

import (
	"embed"
	"io/fs"
)

// berkasMigrasi adalah folder `migrations/` modul ini (rentang 030-049),
// ditanam ke biner. Nama berkas TIDAK berubah dari letak lamanya -
// `T_MIGRASI` mencatat nama, bukan letak.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// SumberMigrasi menyerahkan folder `migrations/` modul ini kepada pelari
// migrasi (`inti/migrasi`).
func SumberMigrasi() fs.FS { return berkasMigrasi }
