// Package backend merakit modul Disease Life (`diseaselife`).
//
// Menu "Disease Life" golongan MASTER TREATY URUTAN 15 (keputusan work owner 08-10-2026 D4; panduan XML
// `D:\NUSARE DEV\Menu Disease\InboxDisease.xml`, section `InboxDisease`, kelas `ASM-FW-GISFW-Int-DISEASE_LIFE`; pola
// modul causeoflosslife / benefitlife). SATU tabel `DISEASE_LIFE` yang SUDAH flat (D1): tidak di-RENAME, kolomnya
// tidak diubah. Migrasinya tinggal di folder modul ini (K0: rentang `080-084` dipinjam dari premiumlistlife, slot menu
// `951` dipinjam dari claimlife) - 080 sequence SEQ_DISEASE_LIFE, 081 PK_DISEASE_LIFE, 951 baris menu.
package backend

import (
	"context"
	"embed"
	"io/fs"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/diseaselife/backend/handlers"
	"nusantarare/modul/diseaselife/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini - 080-081 (sequence + PK DISEASE_LIFE) dan slot menu 951,
// ditanam ke biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// SumberMigrasi menyerahkan folder `migrations/` modul ini kepada pelari migrasi (`inti/backend/migrasi`).
func SumberMigrasi() fs.FS { return berkasMigrasi }

// Nama pengenal modul ini - SAMA dengan nama foldernya.
const Nama = "diseaselife"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`, berkas bangkitan
// `modul_diseaselife_gen.go`). Nol kontrak lintas modul.
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
