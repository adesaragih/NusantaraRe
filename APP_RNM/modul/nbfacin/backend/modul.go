// Package backend merakit modul NB FacIn (`nbfacin`) - tiket 20.
//
// Nama paketnya `backend` - sama di setiap modul (nama folder); daftar modul
// bangkitan mengimpornya dengan alias nama modul.
package backend

import (
	"context"
	"embed"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/nbfacin/backend/handlers"
	"nusantarare/modul/nbfacin/backend/services"
	"nusantarare/modul/nbfacin/backend/services/kontrakfacin"
)

// berkasMigrasi adalah folder `migrations/` modul ini. Rentang tabel 180-219 belum
// dipakai (modul ini hanya MEMBACA tabel POOLDATA yang sudah ada); isinya kini hanya
// slot menu 962 (keputusan work owner butir 59).
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// Nama pengenal modul ini - SAMA dengan nama foldernya.
const Nama = "nbfacin"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`).
//
// Kontrak: MENYEDIAKAN `kontrak.PenilaiPredikatFacIn` (registry `When` NB) dan
// `kontrak.MesinPremiFacIn` (premi coverage NB) - keduanya mesin murni, tanpa data.
// `kontrak.TanggaAkseptasiFacIn` BELUM disediakan lewat perakit: tangga butuh tabel
// limit Oracle per permintaan, dan pemakai pertamanya (`rnwfacin`) dibekukan. Modul
// lain belum dapat memakainya - mengimpor `modul/nbfacin` terlarang (CLAUDE.md §5);
// sisa A32 menunggu penyedia berbasis repository.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:    Nama,
		Migrasi: berkasMigrasi,
		Menyediakan: []inti.Kontrak{
			inti.KontrakDari[kontrak.PenilaiPredikatFacIn](),
			inti.KontrakDari[kontrak.MesinPremiFacIn](),
		},
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			inti.Sediakan[kontrak.PenilaiPredikatFacIn](p, kontrakfacin.Predikat{})
			inti.Sediakan[kontrak.MesinPremiFacIn](p, kontrakfacin.Premi{})
			return Modul{svc: services.DariDasar(p.Dasar())}, nil
		},
	}
}

// Modul memenuhi inti.Modul.
type Modul struct {
	svc *services.Service
}

func (Modul) Nama() string                                 { return Nama }
func (m Modul) DaftarkanRute(mux *http.ServeMux)           { handlers.DaftarkanRute(mux, m.svc) }
func (Modul) JalankanPekerja(context.Context) inti.Pekerja { return inti.TanpaPekerja() }
