// Package backend merakit modul Cover Life (`coverlife`).
//
// Menu "Cover Life" golongan MASTER TREATY URUTAN 16 (keputusan work owner 08-10-2026 C4; panduan XML
// `D:\NUSARE DEV\Menu Cover\InboxCoverLife.xml`, section `InboxCoverLife`, kelas `ASM-FW-GISFW-Int-COVER_LIFE`; pola
// modul causeoflosslife). SATU tabel `M_COVER_LIFE` (C1): nama TETAP (TANPA RENAME), dijadikan flat di tempat, view
// lamanya dibuang. Migrasinya tinggal di folder modul ini (K0: rentang `085-089` dipinjam dari premiumlistlife, slot
// menu `957` dipinjam dari treatycontractout) - 085-086 tabel, 957 baris menu. ID dari sequence warisan (MODUL.md).
package backend

import (
	"context"
	"embed"
	"io/fs"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/coverlife/backend/handlers"
	"nusantarare/modul/coverlife/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini - 085-086 (satu tabel M_COVER_LIFE) dan slot menu 957, ditanam
// ke biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// SumberMigrasi menyerahkan folder `migrations/` modul ini kepada pelari migrasi (`inti/backend/migrasi`).
func SumberMigrasi() fs.FS { return berkasMigrasi }

// Nama pengenal modul ini - SAMA dengan nama foldernya.
const Nama = "coverlife"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`, berkas bangkitan
// `modul_coverlife_gen.go`). Nol kontrak lintas modul.
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
