// Package backend merakit modul PremiumList Life (`premiumlistlife`).
//
// Nama paketnya `backend` - sama di setiap modul (nama folder); daftar modul
// bangkitan mengimpornya dengan alias nama modul.
//
// Refactor bentuk B (30-09-2026): setiap modul punya satu berkas perakitan
// (`modul.go`) yang menyerahkan miliknya kepada `cmd/api`. Sejak struktur tim
// satu folder per modul (30-09-2026) penyerahannya lewat `Pendaftaran()`:
// kontrak yang disediakan dan dibutuhkan DINYATAKAN di sini, dan perakit
// `inti.Rakit` menyambungnya - bukan baris tangan di daftar bersama. Di paket ini
// memuat migrasinya, rute HTTP-nya, dan pekerja latarnya (bila ada).
package backend

import (
	"context"
	"embed"
	"io/fs"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/premiumlistlife/backend/handlers"
	"nusantarare/modul/premiumlistlife/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini (rentang 050-099, MODUL.md),
// ditanam ke biner. Nama berkas TIDAK berubah dari letak lamanya -
// `T_MIGRASI` mencatat nama, bukan letak.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// SumberMigrasi menyerahkan folder `migrations/` modul ini kepada pelari
// migrasi (`inti/backend/migrasi`).
func SumberMigrasi() fs.FS { return berkasMigrasi }

// Nama pengenal modul ini di MODUL_AKTIF dan di GET /api/modul-aktif.
const Nama = "premiumlistlife"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`,
// berkas bangkitan `modul_premiumlistlife_gen.go`).
//
// Kontrak: MENYEDIAKAN `kontrak.PembacaPolis` - butir pl4/av, Claim Life
// membaca polis modul ini. Nama lama `premiumlist` (sebelum tabel nama modul
// 30-09-2026) ditolak MODUL_AKTIF dengan kalimat yang menyebut nama ini.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:        Nama,
		NamaLama:    []string{"premiumlist"},
		Migrasi:     berkasMigrasi,
		Menyediakan: []inti.Kontrak{inti.KontrakDari[kontrak.PembacaPolis]()},
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			svc := services.DariDasar(p.Dasar())
			inti.Sediakan[kontrak.PembacaPolis](p, services.PembacaPolis(svc))
			return Baru(svc, p.Config().AuthStub), nil
		},
	}
}

// Modul adalah perakitan modul PremiumList Life untuk `cmd/api` - `inti.Modul`.
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
