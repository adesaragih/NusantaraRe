// Package backend merakit modul Claim Life (`claimlife`).
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
	"nusantarare/modul/claimlife/backend/handlers"
	"nusantarare/modul/claimlife/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini (rentang 001-029),
// ditanam ke biner. Nama berkas TIDAK berubah dari letak lamanya -
// `T_MIGRASI` mencatat nama, bukan letak.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// SumberMigrasi menyerahkan folder `migrations/` modul ini kepada pelari
// migrasi (`inti/backend/migrasi`).
func SumberMigrasi() fs.FS { return berkasMigrasi }

// Nama pengenal modul ini di MODUL_AKTIF dan di GET /api/modul-aktif.
const Nama = "claimlife"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`,
// berkas bangkitan `modul_claimlife_gen.go`).
//
// Kontrak: MEMBUTUHKAN `kontrak.PembacaPolis` - butir pl4/av, Claim Life
// membaca polis PremiumList; MENYEDIAKAN `kontrak.KlaimKomite` - butir km3,
// Komite membaca dan menuntaskan baris klaim.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:        Nama,
		Migrasi:     berkasMigrasi,
		Membutuhkan: []inti.Kontrak{inti.KontrakDari[kontrak.PembacaPolis]()},
		Menyediakan: []inti.Kontrak{inti.KontrakDari[kontrak.KlaimKomite]()},
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			svc := services.DariDasar(p.Dasar()).DenganPembacaPolis(inti.Ambil[kontrak.PembacaPolis](p))
			inti.Sediakan[kontrak.KlaimKomite](p, services.KlaimUntukKomite(svc))
			return Baru(svc, p.Config().AuthStub), nil
		},
	}
}

// Modul adalah perakitan modul Claim Life untuk `cmd/api` - `inti.Modul`.
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
