// Package backend merakit modul Cause Of Loss Life (`causeoflosslife`).
//
// Menu "Cause Of Loss Life" golongan MASTER TREATY URUTAN 14 (keputusan work owner 08-10-2026 K5; panduan XML
// `D:\NUSARE DEV\Menu Cause Of Loss\InboxCauseofLossLife.xml`, section `InboxCauseofLossLife`, kelas
// `ASM-FW-GISFW-Int-CAUSEOFLOSS_LIFE`; pola modul benefitlife). SATU tabel `CAUSEOFLOSS_LIFE` (K1): tabel Pega
// `M_CAUSEOFLOSS_LIFE` berganti nama dan menjadi flat. Berbeda dengan benefitlife / planlife, migrasinya tinggal di
// folder modul ini (K0: rentang `090-099` dan slot menu `955`, dipinjam dari jatah premiumlistlife) - 090-092 tabel,
// 955 baris menu. ID dari sequence warisan (MODUL.md).
package backend

import (
	"context"
	"embed"
	"io/fs"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/causeoflosslife/backend/handlers"
	"nusantarare/modul/causeoflosslife/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini - 090-092 (satu tabel CAUSEOFLOSS_LIFE) dan slot menu 955,
// ditanam ke biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// SumberMigrasi menyerahkan folder `migrations/` modul ini kepada pelari migrasi (`inti/backend/migrasi`).
func SumberMigrasi() fs.FS { return berkasMigrasi }

// Nama pengenal modul ini - SAMA dengan nama foldernya.
const Nama = "causeoflosslife"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`, berkas bangkitan
// `modul_causeoflosslife_gen.go`). Nol kontrak lintas modul.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:    Nama,
		Migrasi: berkasMigrasi,
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
