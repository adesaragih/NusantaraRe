// Package backend merakit modul Endorsement Life (`endorsementlife`).
//
// Nama paketnya `backend` - sama di setiap modul (nama folder); daftar modul
// bangkitan mengimpornya dengan alias nama modul. Modul ini diserahkan lewat
// `Pendaftaran()` (struktur tim satu folder per modul, 30-09-2026).
//
// ⛔ Nol tabel baru (spec §16, AC 55): endorsement adalah VERSI baru di tabel
// PremiumList Life (`T_PREMIUM_LIST` + anak-anaknya). Migrasinya: 480
// (`PROD_KE DEFAULT 1` + index pencari versi, tiket 00), 481 (sequence kasus
// `EDMLF-<n>`), dan slot menu 976.
package backend

import (
	"context"
	"embed"
	"io/fs"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/endorsementlife/backend/handlers"
	"nusantarare/modul/endorsementlife/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini, ditanam ke biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// SumberMigrasi menyerahkan folder `migrations/` modul ini kepada pelari
// migrasi (`inti/backend/migrasi`).
func SumberMigrasi() fs.FS { return berkasMigrasi }

// Nama pengenal modul ini di MODUL_AKTIF dan di GET /api/modul-aktif.
const Nama = "endorsementlife"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`,
// berkas bangkitan `modul_endorsementlife_gen.go`).
//
// Nol kontrak lintas modul: tabel PremiumList dibaca dan ditulis langsung
// (tabel bersama, spec §16), tabel warisan disebut namanya sendiri.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:    Nama,
		Migrasi: berkasMigrasi,
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			return Baru(services.DariDasar(p.Dasar()), p.Config().AuthStub), nil
		},
	}
}

// Modul adalah perakitan modul Endorsement Life untuk `cmd/api` - `inti.Modul`.
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

// JalankanPekerja - modul ini tidak punya pekerja latar.
func (Modul) JalankanPekerja(context.Context) inti.Pekerja { return inti.TanpaPekerja() }
