// Package backend merakit modul NB Treaty In (`nbtreatyin`) - realisasi
// penutupan treaty inward: dari kontrak treaty yang disetujui menjadi polis
// treaty, melalui tangga Admin -> Sec Head -> Dept Head (spec.md).
//
// Nama paketnya `backend` - sama di setiap modul (nama folder); daftar modul
// bangkitan (`inti/backend/daftar/modul_nbtreatyin_gen.go`) mengimpornya dengan
// alias nama modul. Padanan Pega: Harness `SFAPortalOpportunities` dan
// `Flow/InputRealizationTreatyIn` (INVENTARIS-XML.md bab 3-4).
package backend

import (
	"context"
	"embed"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/dokumenpolis"
	"nusantarare/modul/nbtreatyin/backend/handlers"
	"nusantarare/modul/nbtreatyin/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini (rentang 320-359, slot
// menu 968-969 - MODUL.md), ditanam ke biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// Nama pengenal modul ini di MODUL_AKTIF dan di GET /api/modul-aktif.
const Nama = "nbtreatyin"

// Pendaftaran menyerahkan modul ini kepada perakit. Nol kontrak disediakan
// dan nol kontrak dipakai: modul ini tidak membaca milik modul lain lewat
// kode (tabel bersama T_WORK_POLIS dibaca dan ditulis lewat SQL-nya sendiri).
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:    Nama,
		Migrasi: berkasMigrasi,
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			m := Baru(services.DariDasar(p.Dasar()), p.Config().AuthStub)
			// Lampiran "Reas" (keputusan work owner 08-10-2026) - penyimpanan bersama, garam token dari config.
			m.lampiran = dokumenpolis.Oracle(p.Dasar(), p.Config().StorageTokenSalt)
			return m, nil
		},
	}
}

// Modul adalah perakitan modul NB Treaty In untuk `cmd/api` - `inti.Modul`.
type Modul struct {
	svc        *services.Layanan
	stubPelaku bool
	// lampiran - lampiran "Reas" kasus (`inti/backend/dokumenpolis`); nil = rutenya tidak dipasang.
	lampiran *dokumenpolis.Layanan
}

// Baru merakit modul di atas layanan yang sudah disambung.
func Baru(svc *services.Layanan, stubPelaku bool) Modul {
	return Modul{svc: svc, stubPelaku: stubPelaku}
}

// Nama menyebut modul ini.
func (Modul) Nama() string { return Nama }

// DaftarkanRute mendaftarkan seluruh rute modul ini ke mux bersama.
func (m Modul) DaftarkanRute(mux *http.ServeMux) {
	handlers.DaftarkanRute(mux, m.svc, m.stubPelaku)
	if m.lampiran != nil {
		handlers.DaftarkanLampiran(mux, m.svc, m.lampiran, m.stubPelaku)
	}
}

// JalankanPekerja - modul ini tidak punya pekerja latar.
func (Modul) JalankanPekerja(context.Context) inti.Pekerja { return inti.TanpaPekerja() }
