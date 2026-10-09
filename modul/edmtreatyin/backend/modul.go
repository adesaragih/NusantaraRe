// Package backend merakit modul EDM Treaty In (`edmtreatyin`) - endorsemen (adendum) polis treaty inward yang
// sudah terbit: generasi baru polis yang sama, dibandingkan dengan generasi tepat sebelumnya, melalui tangga
// Admin -> Sec Head -> Dept Head.
//
// Nama paketnya `backend` - sama di setiap modul (nama folder); daftar modul bangkitan
// (`inti/backend/daftar/modul_edmtreatyin_gen.go`) mengimpornya dengan alias nama modul. Padanan Pega: Harness
// `SFAPortal_Endorsement_Treaty`, `TreatyCreateEdm`, dan `Flow/InputAddendumTreatyIn` kelas
// `ASM-FW-GISFW-Work-EndorsementTreaty` (docs/INVENTARIS-XML.md).
//
// ⛔ Kode yang disalin dari modul NB Treaty In (`modul/nbtreatyin/backend`, snapshot 06-10-2026 di atas HEAD
// 435d3879) menyebut asalnya di kepala berkas: modul tidak saling mengimpor (CLAUDE.md §5, PANDUAN-TIM-PER-MODUL
// bab 7), dan NB tidak menyediakan kontrak.
package backend

import (
	"context"
	"embed"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/dokumenpolis"
	"nusantarare/modul/edmtreatyin/backend/handlers"
	"nusantarare/modul/edmtreatyin/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini (rentang 360-399, slot menu 970-971 - MODUL.md), ditanam
// ke biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// Nama pengenal modul ini di MODUL_AKTIF dan di GET /api/modul-aktif.
const Nama = "edmtreatyin"

// Pendaftaran menyerahkan modul ini kepada perakit. Nol kontrak disediakan dan nol kontrak dipakai: tabel
// generasi polis treaty (milik nbtreatyin, migrasi 320-327) dan T_WORK_POLIS (milik premiumlistlife) dibaca dan
// ditulis lewat SQL modul ini sendiri, seperti NB membaca T_WORK_POLIS.
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

// Modul adalah perakitan modul EDM Treaty In untuk `cmd/api` - `inti.Modul`.
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
