// Package backend merakit modul Province (`masterprovince`).
//
// Modul DI LUAR dua puluh folder korpus (PANDUAN-TIM-PER-MODUL bab 5). Keputusan work owner 04-10-2026: Master Data
// dipecah menjadi delapan modul menu terpisah di grup MASTER ("8 modul terpisah"); mesinnya bersama di
// `inti/backend/master` ("Pindah ke inti"). Modul ini = master `province` di `/api/master-province`.
// Pemilik migrasi 880-882 (keenam tabel flat master + T_MASTER_STATUS + jejak ubah; keputusan work owner
// 04-10-2026 "Dihapus, tabel pindah ke Province").
// Baris M_NAV_MENU-nya dari migrasi inti 913 (tanpa slot menu: "Modul luar korpus tanpa slot").
package backend

import (
	"embed"
	"io/fs"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/master"
)

// berkasMigrasi - migrasi 880-882 (dulu milik modul masterdata; nama berkas tetap - T_MIGRASI mencatat nama).
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// SumberMigrasi - folder migrations modul ini.
func SumberMigrasi() fs.FS { return berkasMigrasi }

// Nama - nama modul (folder, MODUL_AKTIF, KODE menu).
const Nama = "masterprovince"

// Pendaftaran - master `province` di `/api/master-province` (kontrak `modul/masterprovince/docs/issues/03-pemecahan-delapan-modul.md`).
func Pendaftaran() inti.Pendaftaran {
	return master.Pendaftaran(Nama, "/api/master-province", "province", berkasMigrasi)
}
