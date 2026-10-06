// Package backend merakit modul Adjuster Consultant (`adjusterconsultant`).
//
// Kelola tabel warisan `POOLDATA.ADJUSTERCONSULTANT` - tambah, ubah, aktif/nonaktif, nol hapus (perintah work owner
// 05-10-2026: "buatkan modul/menu CRUD nya ... jangan di inti, buat modul sendiri modul/adjusterconsultant"). Padanan
// Pega: layar master `MstAdjusterConsultant` (folder korpus Claim Fac In dan Claim Prop). Nol tabel baru; migrasi 870
// menambah kolom `IS_ACTIVE`, slot menu 997 menyalakan menunya, baris menunya dibuat migrasi inti 915.
package backend

import (
	"context"
	"embed"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/adjusterconsultant/backend/handlers"
	"nusantarare/modul/adjusterconsultant/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini, ditanam ke biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// Nama pengenal modul ini - SAMA dengan nama foldernya.
const Nama = "adjusterconsultant"

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
