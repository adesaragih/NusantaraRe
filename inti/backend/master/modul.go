// Package master - mesin master bersama (tim inti). Keputusan work owner 04-10-2026: Master Data dipecah menjadi
// delapan modul menu terpisah ("8 modul terpisah", pola Marketing Officer / Company Detail); mesinnya pindah ke inti
// ("Pindah ke inti"), karena `modul/X` tidak boleh mengimpor `modul/Y`.
//
// Satu modul menu = satu master: `modul/<nama>/backend/modul.go` hanya memanggil `Pendaftaran` di sini dengan nama
// modul, prefix rutenya, dan kunci masternya (`models.DaftarMaster`). Rute dipasang di mux modul itu, jadi gerbang menu
// `cmd/api` (Kelola User) menagih menu modul itu sendiri. Kontrak rute: `modul/masterprovince/docs/issues/03-…`.
package master

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/master/handlers"
	"nusantarare/inti/backend/master/models"
	"nusantarare/inti/backend/master/services"
)

// Pendaftaran - pendaftaran modul menu master `kunci` bernama `nama`, berute di bawah `prefix`. `migrasi` = folder
// migrations modul itu, nil bila tidak bermigrasi.
func Pendaftaran(nama, prefix, kunci string, migrasi fs.FS) inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:    nama,
		Migrasi: migrasi,
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			if _, ada := models.CariMaster(kunci); !ada {
				return nil, fmt.Errorf("master: modul %s menunjuk master %q yang tidak ada", nama, kunci)
			}
			return Modul{nama: nama, prefix: prefix, kunci: kunci, svc: services.DariDasar(p.Dasar()),
				stubPelaku: p.Config().AuthStub}, nil
		},
	}
}

// Modul - satu modul menu master.
type Modul struct {
	nama, prefix, kunci string
	svc                 *services.Service
	stubPelaku          bool
}

// Nama - nama modul (MODUL_AKTIF, GET /api/modul-aktif).
func (m Modul) Nama() string { return m.nama }

// DaftarkanRute memasang rute master modul ini (`handlers.Pasang`).
func (m Modul) DaftarkanRute(mux *http.ServeMux) {
	handlers.Pasang(mux, m.prefix, m.kunci, m.svc, m.stubPelaku)
}

// JalankanPekerja - tanpa pekerja latar.
func (Modul) JalankanPekerja(context.Context) inti.Pekerja { return inti.TanpaPekerja() }
