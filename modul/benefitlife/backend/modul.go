// Package backend merakit modul Benefit (`benefitlife`).
//
// Menu "Benefit" golongan MASTER TREATY URUTAN 12 (keputusan work owner 08-10-2026 K4; panduan XML
// `D:\NUSARE DEV\Menu Benefit\InboxBenefit.xml`, section `InboxBenefit`, kelas `ASM-FW-GISFW-Int-BENEFIT_LIFE`; pola
// modul ririsklife). SATU tabel `BENEFIT_LIFE` (K1): tabel Pega `M_BENEFIT_LIFE` berganti nama dan menjadi flat
// (migrasi inti 942-944; baris menu 945). NOL MIGRASI SENDIRI; ID dari sequence warisan (MODUL.md).
package backend

import (
	"context"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/benefitlife/backend/handlers"
	"nusantarare/modul/benefitlife/backend/services"
)

// Nama pengenal modul ini - SAMA dengan nama foldernya.
const Nama = "benefitlife"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`). Nol kontrak lintas modul, nol migrasi.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama: Nama,
		// Akses menu LIHAT: ikut gerbang tulis `cmd/api`; nol pola bebas.
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
