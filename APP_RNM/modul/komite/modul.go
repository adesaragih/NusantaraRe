// Package komite merakit modul Komite Claim Life.
//
// Refactor bentuk B (30-09-2026): setiap modul punya satu berkas perakitan
// (`modul.go`) yang menyerahkan miliknya kepada `cmd/api`. Di paket ini
// memuat migrasinya, rute HTTP-nya, dan pekerja latarnya (bila ada).
package komite

import (
	"context"
	"embed"
	"io/fs"
	"net/http"

	"nusantarare/inti"
	"nusantarare/modul/komite/handlers"
	"nusantarare/modul/komite/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini (rentang 030-049),
// ditanam ke biner. Nama berkas TIDAK berubah dari letak lamanya -
// `T_MIGRASI` mencatat nama, bukan letak.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// SumberMigrasi menyerahkan folder `migrations/` modul ini kepada pelari
// migrasi (`inti/migrasi`).
func SumberMigrasi() fs.FS { return berkasMigrasi }

// Nama pengenal modul ini di MODUL_AKTIF dan di GET /api/modul-aktif.
const Nama = "komite"

// Modul adalah perakitan modul Komite Claim Life untuk `cmd/api` - `inti.Modul`.
type Modul struct {
	svc        *services.Service
	stubPelaku bool
}

// Baru merakit modul di atas Service yang sudah disambung `modul.Rakit`.
func Baru(svc *services.Service, stubPelaku bool) Modul {
	return Modul{svc: svc, stubPelaku: stubPelaku}
}

// Nama menyebut modul ini.
func (Modul) Nama() string { return Nama }

// DaftarkanRute mendaftarkan seluruh rute modul ini ke mux bersama.
func (m Modul) DaftarkanRute(mux *http.ServeMux) { handlers.DaftarkanRute(mux, m.svc, m.stubPelaku) }

// JalankanPekerja - modul ini tidak punya pekerja latar.
func (Modul) JalankanPekerja(context.Context) inti.Pekerja { return inti.TanpaPekerja() }
