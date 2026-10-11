// Package backend merakit modul Treaty Contract Out (`treatycontractout`).
//
// Nama paketnya `backend` - sama di setiap modul (nama folder); daftar modul
// bangkitan mengimpornya dengan alias nama modul. Sejak struktur tim satu
// folder per modul (30-09-2026) modul ini diserahkan lewat `Pendaftaran()`.
//
// Refactor bentuk B (30-09-2026): setiap modul punya satu berkas perakitan
// (`modul.go`) yang menyerahkan miliknya kepada `cmd/api`: rute HTTP-nya.
// Modul ini NOL migrasi - ia menulis dan membaca tabel warisan (tco4). Fitur
// lampiran beserta pekerja latarnya dibuang (keputusan work owner 08-10-2026).
package backend

import (
	"context"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatycontractout/backend/handlers"
	"nusantarare/modul/treatycontractout/backend/services"
)

// Nama pengenal modul ini di MODUL_AKTIF dan di GET /api/modul-aktif.
const Nama = "treatycontractout"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`,
// berkas bangkitan `modul_treatycontractout_gen.go`).
//
// Nol kontrak lintas modul, nol migrasi (tco4: tabel warisan). Nama lama
// `treaty` ditolak MODUL_AKTIF dengan kalimat yang menyebut nama ini.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:     Nama,
		NamaLama: []string{"treaty"},
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			return Baru(services.DariDasar(p.Dasar()), p.Config().AuthStub), nil
		},
	}
}

// Modul adalah perakitan modul Treaty Contract Out untuk `cmd/api` -
// `inti.Modul`.
type Modul struct {
	svc        *services.Service
	stubPelaku bool
}

// Baru merakit modul di atas Service yang sudah disambung (`Pendaftaran`).
func Baru(svc *services.Service, stubPelaku bool) Modul {
	return Modul{svc: svc, stubPelaku: stubPelaku}
}

// Nama menyebut modul ini.
func (Modul) Nama() string { return Nama }

// DaftarkanRute mendaftarkan seluruh rute modul ini ke mux bersama.
func (m Modul) DaftarkanRute(mux *http.ServeMux) { handlers.DaftarkanRute(mux, m.svc, m.stubPelaku) }

// JalankanPekerja - tanpa pekerja latar (pekerja antrean lampiran dibuang
// bersama fiturnya, keputusan work owner 08-10-2026).
func (Modul) JalankanPekerja(context.Context) inti.Pekerja { return inti.TanpaPekerja() }
