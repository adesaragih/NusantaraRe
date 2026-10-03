// Package backend merakit modul Komite Claim Life (`komiteclaimlife`).
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
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimlife/backend/handlers"
	"nusantarare/modul/komiteclaimlife/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini (rentang 030-049),
// ditanam ke biner. Nama berkas TIDAK berubah dari letak lamanya -
// `T_MIGRASI` mencatat nama, bukan letak.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// Nama pengenal modul ini di MODUL_AKTIF dan di GET /api/modul-aktif.
const Nama = "komiteclaimlife"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`,
// berkas bangkitan `modul_komiteclaimlife_gen.go`).
//
// Kontrak: MEMBUTUHKAN `kontrak.KlaimKomite` - butir km3, Komite membaca dan
// menuntaskan baris klaim Claim Life. Nama lama `komite` ditolak MODUL_AKTIF.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:        Nama,
		NamaLama:    []string{"komite"},
		Migrasi:     berkasMigrasi,
		Membutuhkan: []inti.Kontrak{inti.KontrakDari[kontrak.KlaimKomite]()},
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			svc := services.DariDasar(p.Dasar()).DenganKlaim(inti.Ambil[kontrak.KlaimKomite](p))
			return Baru(svc, p.Config().AuthStub), nil
		},
	}
}

// Modul adalah perakitan modul Komite Claim Life untuk `cmd/api` - `inti.Modul`.
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
