// Package backend merakit modul Master Contract Retro Life (`mastercontractretrolife`).
//
// Nama paketnya `backend` - sama di setiap modul (nama folder); daftar modul
// bangkitan mengimpornya dengan alias nama modul. Modul ini diserahkan lewat
// `Pendaftaran()` (struktur tim satu folder per modul, 30-09-2026).
//
// ⛔ K1: nol DDL, nol tabel baru - modul ini menulis dan membaca lima tabel
// warisan `POOLDATA` (`docs/STRUKTUR-TABEL-MASTER-CONTRACT-RETRO-LIFE.md`).
// Satu-satunya migrasinya adalah slot menu 958 (`UPDATE DIMIGRASI`).
package backend

import (
	"context"
	"embed"
	"io/fs"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/mastercontractretrolife/backend/handlers"
	"nusantarare/modul/mastercontractretrolife/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini - hanya slot menu 958
// (rentang tabel 100-139 sengaja kosong, K1), ditanam ke biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// SumberMigrasi menyerahkan folder `migrations/` modul ini kepada pelari
// migrasi (`inti/backend/migrasi`).
func SumberMigrasi() fs.FS { return berkasMigrasi }

// Nama pengenal modul ini di MODUL_AKTIF dan di GET /api/modul-aktif.
const Nama = "mastercontractretrolife"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`,
// berkas bangkitan `modul_mastercontractretrolife_gen.go`).
//
// Nol kontrak lintas modul: master rujukan (jenis reasuransi, reinsurer,
// business) dibaca langsung dari tabelnya, seperti RD Pega.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:    Nama,
		Migrasi: berkasMigrasi,
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			return Baru(services.DariDasar(p.Dasar()), p.Config().AuthStub), nil
		},
	}
}

// Modul adalah perakitan modul Master Contract Retro Life untuk `cmd/api` -
// `inti.Modul`.
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
