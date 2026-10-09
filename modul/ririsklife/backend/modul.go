// Package backend merakit modul R/I Risk (`ririsklife`).
//
// Menu "R/I Risk" golongan MASTER TREATY URUTAN 11 (keputusan work owner 08-10-2026 K4; panduan XML
// `D:\NUSARE DEV\Menu RI Risk\InboxSummaryRIRisk.xml` dan `InboxRIRisk.xml`; pola modul ricommlife). SATU tabel per jenis
// data (K1/K2): ringkasan = `RIRISK_LIFE_SUMMARY`, rincian = `RIRISK_LIFE` - tabel Pega `M_RIRISK_LIFE*` berganti nama
// dan menjadi flat (migrasi inti 935-940; baris menu 941). NOL MIGRASI SENDIRI; ID dari sequence warisan (MODUL.md).
package backend

import (
	"context"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/ririsklife/backend/handlers"
	"nusantarare/modul/ririsklife/backend/services"
)

// Nama pengenal modul ini - SAMA dengan nama foldernya.
const Nama = "ririsklife"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`). Nol kontrak lintas modul, nol migrasi.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama: Nama,
		// Akses menu LIHAT: ikut gerbang tulis `cmd/api`; nol pola bebas (View Upload pun bagi Full saja).
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
