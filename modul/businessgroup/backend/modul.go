// Package backend merakit modul Business Group (`businessgroup`).
//
// Kelola tabel warisan `POOLDATA.BUSINESSGROUP` - tambah dan ubah, nol hapus, tanpa grup berakhiran SYARIAH
// (keputusan work owner 05-10-2026). Modul di luar korpus: Pega hanya membaca tabel ini (`BrowseBusinessGroup_RD`).
// Nol tabel baru, nol DDL; slot menu 991 menyalakan menunya, barisnya dibuat migrasi inti 918.
package backend

import (
	"context"
	"embed"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/businessgroup/backend/handlers"
	"nusantarare/modul/businessgroup/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini, ditanam ke biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// Nama pengenal modul ini - SAMA dengan nama foldernya.
const Nama = "businessgroup"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`). Nol kontrak lintas modul.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:    Nama,
		Migrasi: berkasMigrasi,
		// Akses menu LIHAT (keputusan work owner 05-10-2026): ikut gerbang tulis `cmd/api`.
		HakLihat: &inti.HakLihat{},
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
