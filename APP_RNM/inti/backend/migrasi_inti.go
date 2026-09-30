package backend

// Migrasi milik `inti` - tabel lintas modul, rentang 900-949.
//
// Untuk apa berkas ini: `M_NAV_MENU` (900, `PROMPT-MENU-DARI-TABEL-M_NAV_MENU.md`)
// bukan milik satu modul - ia memuat menu SEMUA modul, termasuk yang belum
// dimigrasi. Karena itu migrasinya tinggal di `inti/backend/migrations/`, dan
// `daftar.SumberMigrasi` (`inti/backend/daftar`) mengumpulkannya bersama
// folder migrasi setiap modul.

import (
	"embed"
	"io/fs"
)

// berkasMigrasi adalah folder `migrations/` inti (rentang 900-949), ditanam ke
// biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// SumberMigrasi menyerahkan folder `migrations/` inti kepada pelari migrasi
// (`inti/backend/migrasi`).
func SumberMigrasi() fs.FS { return berkasMigrasi }
