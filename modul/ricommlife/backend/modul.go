// Package backend merakit modul R/I Comm Life (`ricommlife`).
//
// Menu "R/I Comm Life" golongan MASTER TREATY (perintah work owner 06-10-2026: "membuat modul baru di Master Treaty
// dengan nama R/I Comm Life, panduannya baca dari xml di D:\NUSARE DEV\Menu RI Comm, konsepnya hampir sama dengan menu
// R/I Rate, membuat CRUD, dan detail bisa di save dan edit"). Ringkasan = JSON warisan `M_RICOMM_LIFE_SUMMARY` (view
// `RICOMM_LIFE_SUMMARY` dibaca); rincian = tabel flat `RICOMM_LIFE` (migrasi inti 924). NOL MIGRASI SENDIRI - tabel
// flat 924 dan baris menu 925 adalah migrasi inti; ID dari sequence warisan (MODUL.md, keputusan work owner
// 06-10-2026).
package backend

import (
	"context"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/ricommlife/backend/handlers"
	"nusantarare/modul/ricommlife/backend/services"
)

// Nama pengenal modul ini - SAMA dengan nama foldernya.
const Nama = "ricommlife"

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
