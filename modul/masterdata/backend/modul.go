// Package backend merakit modul Master Data (`masterdata`).
//
// Modul DI LUAR dua puluh folder korpus (PANDUAN-TIM-PER-MODUL bab 5), keputusan work owner 04-10-2026: menu master
// insert / update / aktif / nonaktif untuk tabel master yang dipakai NB FacIn (`docs/issues/01-rencana-master-data.md`).
// Migrasinya: 880 (view warisan -> tabel flat bernama sama + STS_AKTIF), 881 (T_MASTER_STATUS: status NATION), 882
// (jejak ubah: CREATE_OP / TGL_CREATE / UPDATE_OP / TGL_UPDATE), dan
// slot menu 996. Baris M_NAV_MENU-nya dari migrasi inti 911.
package backend

import (
	"context"
	"embed"
	"io/fs"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/masterdata/backend/handlers"
	"nusantarare/modul/masterdata/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini, ditanam ke biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// SumberMigrasi menyerahkan folder `migrations/` modul ini kepada pelari migrasi (`inti/backend/migrasi`).
func SumberMigrasi() fs.FS { return berkasMigrasi }

// Nama pengenal modul ini di MODUL_AKTIF dan di GET /api/modul-aktif - SAMA dengan nama foldernya.
const Nama = "masterdata"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`, berkas bangkitan
// `modul_masterdata_gen.go`). Nol kontrak lintas modul: tabel master dibaca dan ditulis langsung; modul lain (nbfacin)
// membacanya sebagai tabel.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:    Nama,
		Migrasi: berkasMigrasi,
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			return Baru(services.DariDasar(p.Dasar()), p.Config().AuthStub), nil
		},
	}
}

// Modul adalah perakitan modul Master Data untuk `cmd/api` - `inti.Modul`.
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
