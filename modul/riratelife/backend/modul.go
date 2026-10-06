// Package backend merakit modul R/I Rate Life (`riratelife`).
//
// Menu "R/I Rate Life" golongan MASTER TREATY (perintah work owner 05-10-2026: "Buat Menu baru Namanya R/I Rate Life
// pada Master Treaty, menu ini bisa CRUD untuk simpan data ke tabel RATE_LIFE_SUMMARY, panduannya xml yang saya
// berikan"). RALAT R4 (keputusan work owner 06-10-2026 K-F1/K-F2): ringkasan dibaca dan ditulis di TABEL FLAT
// `RATE_LIFE_SUMMARY` (migrasi inti 926, pengganti view); `M_RATE_LIFE_SUMMARY` (JSON) dibaca saja; rincian `M_RATE_LIFE`
// ditulis, view `RATE_LIFE` dibaca. NOL MIGRASI SENDIRI - menu = migrasi inti 922, sequence ID = 923, tabel flat = 926.
package backend

import (
	"context"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/riratelife/backend/handlers"
	"nusantarare/modul/riratelife/backend/services"
)

// Nama pengenal modul ini - SAMA dengan nama foldernya.
const Nama = "riratelife"

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
