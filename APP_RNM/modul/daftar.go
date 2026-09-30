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

	"nusantarare/inti"
	"nusantarare/inti/config"
	"nusantarare/modul/claimlife"
	claimlifeservices "nusantarare/modul/claimlife/services"
	"nusantarare/modul/komiteclaimlife"
	komiteservices "nusantarare/modul/komiteclaimlife/services"
	"nusantarare/modul/premiumlistlife"
	premiumlistservices "nusantarare/modul/premiumlistlife/services"
	"nusantarare/modul/treatycontractout"
	treatyservices "nusantarare/modul/treatycontractout/services"
)

// SumberMigrasi mengembalikan folder migrasi SETIAP modul terdaftar.
//
// ⛔ Tidak bergantung pada modul yang aktif: skema selalu utuh, sebab tabel
// satu modul dapat dirujuk tabel modul lain (dan data warisan tidak memilih
// modul). Pelari mengurutkan langkah menurut NAMA berkas, lintas sumber.
// Treaty Contract Out tidak bermigrasi (tco4: tabel warisan).
func SumberMigrasi() []fs.FS {
	return []fs.FS{
		claimlife.SumberMigrasi(),       // 001-029
		komiteclaimlife.SumberMigrasi(), // 030-049
		premiumlistlife.SumberMigrasi(), // 050-079
	}
}

// Rakit membangun SETIAP modul terdaftar di atas satu akar bersama, dan
// menyambung kontrak lintas modulnya (`inti/kontrak`).
//
// Urutannya urutan pendaftaran rute dan urutan GET /api/modul-aktif. `catat`
// adalah pencatat proses (log) untuk modul yang mencatat saat menyala.
func Rakit(dasar *inti.Dasar, cfg config.Config, catat func(string)) []inti.Modul {
	svcPL := premiumlistservices.DariDasar(dasar)
	// Butir pl4/av: Claim Life membaca polis PremiumList lewat inti/kontrak.
	svcCL := claimlifeservices.DariDasar(dasar).DenganPembacaPolis(premiumlistservices.PembacaPolis(svcPL))
	// Butir km3: Komite membaca dan menuntaskan baris klaim lewat inti/kontrak.
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
