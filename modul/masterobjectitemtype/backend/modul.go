// Package backend merakit modul Object Item Type (`masterobjectitemtype`).
//
// Modul DI LUAR dua puluh folder korpus (PANDUAN-TIM-PER-MODUL bab 5). Keputusan work owner 04-10-2026: Master Data
// dipecah menjadi delapan modul menu terpisah di grup MASTER ("8 modul terpisah"); mesinnya bersama di
// `inti/backend/master` ("Pindah ke inti"). Modul ini = master `objectitemtype` di `/api/master-object-item-type`.
// Baris M_NAV_MENU-nya dari migrasi inti 919 (tanpa slot menu: "Modul luar korpus tanpa slot").
package backend

import (
	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/master"
)

// Nama - nama modul (folder, MODUL_AKTIF, KODE menu).
const Nama = "masterobjectitemtype"

// Pendaftaran - master `objectitemtype` di `/api/master-object-item-type` (kontrak `modul/masterprovince/docs/issues/03-pemecahan-delapan-modul.md`).
func Pendaftaran() inti.Pendaftaran {
	return master.Pendaftaran(Nama, "/api/master-object-item-type", "objectitemtype", nil)
}
