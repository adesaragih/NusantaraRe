// Package modul adalah DAFTAR modul aplikasi - satu-satunya paket yang
// mengenal semua modul sekaligus.
//
// Refactor bentuk B (30-09-2026). Modul di `modul/<nama>/` tidak pernah
// mengimpor modul lain; paket ini yang merakit dan menyambungnya, dan
// `cmd/api` yang memilih mana yang aktif (MODUL_AKTIF). Yang memakai daftar
// ini: `cmd/api` dan skema uji test bertag `db`.
//
// Menambah modul baru = satu berkas `modul/<nama>/modul.go`, satu baris di
// `Rakit`, dan - bila modul itu bermigrasi - satu baris di `SumberMigrasi`.
package modul

import (
	"io/fs"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/config"
	claimlife "nusantarare/modul/claimlife/backend"
	claimlifeservices "nusantarare/modul/claimlife/backend/services"
	komiteclaimlife "nusantarare/modul/komiteclaimlife/backend"
	komiteservices "nusantarare/modul/komiteclaimlife/backend/services"
	premiumlistlife "nusantarare/modul/premiumlistlife/backend"
	premiumlistservices "nusantarare/modul/premiumlistlife/backend/services"
	treatycontractout "nusantarare/modul/treatycontractout/backend"
	treatyservices "nusantarare/modul/treatycontractout/backend/services"
)

// SumberMigrasi mengembalikan folder migrasi SETIAP modul terdaftar.
//
// ⛔ Tidak bergantung pada modul yang aktif: skema selalu utuh, sebab tabel
// satu modul dapat dirujuk tabel modul lain (dan data warisan tidak memilih
// modul). Pelari mengurutkan langkah menurut NAMA berkas, lintas sumber.
// Treaty Contract Out tidak bermigrasi (tco4: tabel warisan). `inti` membawa
// tabel lintas modul (M_NAV_MENU, 900).
func SumberMigrasi() []fs.FS {
	return []fs.FS{
		claimlife.SumberMigrasi(),       // 001-029
		komiteclaimlife.SumberMigrasi(), // 030-049
		premiumlistlife.SumberMigrasi(), // 050-079
		inti.SumberMigrasi(),            // 900-949
	}
}

// NamaLama memetakan nama modul SEBELUM tabel nama modul (keputusan work
// owner 30-09-2026, `PROMPT-REFACTOR-NAMA-MODUL.md`) ke namanya kini. Hanya
// dipakai untuk MENOLAK nama lama di MODUL_AKTIF dengan kalimat yang menyebut
// nama barunya - tidak pernah untuk menerimanya. Setiap nilainya wajib nama
// modul terdaftar (dijaga `cmd/api/rakit_test.go`).
var NamaLama = map[string]string{
	"premiumlist": "premiumlistlife",
	"komite":      "komiteclaimlife",
	"treaty":      "treatycontractout",
}

// Rakit membangun SETIAP modul terdaftar di atas satu akar bersama, dan
// menyambung kontrak lintas modulnya (`inti/backend/kontrak`).
//
// Urutannya urutan pendaftaran rute dan urutan GET /api/modul-aktif. `catat`
// adalah pencatat proses (log) untuk modul yang mencatat saat menyala.
func Rakit(dasar *inti.Dasar, cfg config.Config, catat func(string)) []inti.Modul {
	svcPL := premiumlistservices.DariDasar(dasar)
	// Butir pl4/av: Claim Life membaca polis PremiumList lewat inti/backend/kontrak.
	svcCL := claimlifeservices.DariDasar(dasar).DenganPembacaPolis(premiumlistservices.PembacaPolis(svcPL))
	// Butir km3: Komite membaca dan menuntaskan baris klaim lewat inti/backend/kontrak.
	svcKM := komiteservices.DariDasar(dasar).DenganKlaim(claimlifeservices.KlaimUntukKomite(svcCL))
	svcTCO := treatyservices.DariDasar(dasar).
		// OQ-TCO-08: bawaan stub; ⛔ garam tidak pernah dicetak.
		DenganPenyimpananLampiranTCO(cfg.PelaksanaStorage == config.PelaksanaStorageNyata, cfg.StorageTokenSalt)
	return []inti.Modul{
		claimlife.Baru(svcCL, cfg.AuthStub),
		premiumlistlife.Baru(svcPL, cfg.AuthStub),
		komiteclaimlife.Baru(svcKM, cfg.AuthStub),
		treatycontractout.Baru(svcTCO, cfg.AuthStub, cfg.IntervalPekerjaLampiranTCO, cfg.PelaksanaStorage, catat),
	}
}
