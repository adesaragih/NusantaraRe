// Package backend merakit modul R/I Rate Life (`riratelife`).
//
// Menu "R/I Rate Life" golongan MASTER TREATY (perintah work owner 05-10-2026: "Buat Menu baru Namanya R/I Rate Life
// pada Master Treaty, menu ini bisa CRUD untuk simpan data ke tabel RATE_LIFE_SUMMARY, panduannya xml yang saya
// berikan"). RALAT R4 (tabel flat, 06-10-2026) =
// riwayat; RALAT R6 (keputusan work owner 07-10-2026): ringkasan dibaca dan ditulis di kolom SATU tabel
// `M_RATE_LIFE_SUMMARY` (migrasi inti 927/928); RALAT R7: rincian dibaca dan ditulis di kolom tabel flat
// `M_RATE_LIFE`. NOL MIGRASI SENDIRI - menu = migrasi inti 922, sequence ID = 923, ringkasan = 926 (riwayat) + 927/928,
// rincian = 929/930.
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
