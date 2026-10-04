// Package backend merakit modul Accounts (`accounts`).
//
// Kelola tabel warisan `POOLDATA.T_M_ACCOUNT` (Pega SFAGIS Account) - tambah dan ubah, nol hapus (keputusan work
// owner 04-10-2026: "buat modul/menu baru, nama modulnya accounts, tabelnya T_M_ACCOUNT"; "Delete tidak").
// Migrasi modul: 840 (CREATEDATE, CREATEOP, DESCRIPTION), 841 (PK), 842 (SEQ_T_M_ACCOUNT), slot menu 994.
package backend

import (
	"context"
	"embed"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/accounts/backend/handlers"
	"nusantarare/modul/accounts/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini, ditanam ke biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// Nama pengenal modul ini - SAMA dengan nama foldernya.
const Nama = "accounts"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`). Nol kontrak lintas modul: organisasi
// `CLIENT` dan `BUSINESSGROUP` dibaca langsung dari tabelnya.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:    Nama,
		Migrasi: berkasMigrasi,
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
