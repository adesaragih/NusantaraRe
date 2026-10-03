// Package backend merakit modul Marketing Officer (`marketingofficer`).
//
// Kelola tabel warisan `POOLDATA.MARKETINGOFFICER` - tambah dan ubah, nol hapus (perintah work owner 03-10-2026:
// "buatkan di modul baru dengan nama marketingofficer, gunanya untuk insert update table marketingofficer; jangan
// di inti"). Padanan Pega: form `InputMarketingOfficer` (NB FacIn, RNW Fac In, Endorsment Fac In). Nol tabel baru,
// nol DDL; satu-satunya migrasi modul ini adalah slot menu 990 (`UPDATE DIMIGRASI`).
package backend

import (
	"context"
	"embed"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/marketingofficer/backend/handlers"
	"nusantarare/modul/marketingofficer/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini - slot menu 990, ditanam ke biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// Nama pengenal modul ini - SAMA dengan nama foldernya.
const Nama = "marketingofficer"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`). Nol kontrak lintas modul: akun
// `M_LOGIN_GO` dan cabang `BRANCH` dibaca langsung dari tabelnya.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:    Nama,
		Migrasi: berkasMigrasi,
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			return Modul{svc: services.DariDasar(p.Dasar()), stubPelaku: p.Config().AuthStub}, nil
		},
	}
}

// Modul memenuhi inti.Modul.
type Modul struct {
	svc        *services.Service
	stubPelaku bool
}

// Nama menyebut modul ini.
func (Modul) Nama() string { return Nama }

// DaftarkanRute mendaftarkan seluruh rute modul ini ke mux bersama.
func (m Modul) DaftarkanRute(mux *http.ServeMux) { handlers.DaftarkanRute(mux, m.svc, m.stubPelaku) }

// JalankanPekerja - modul ini tidak punya pekerja latar.
func (Modul) JalankanPekerja(context.Context) inti.Pekerja { return inti.TanpaPekerja() }
