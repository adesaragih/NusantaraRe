// Package backend merakit modul Company Detail (`companydetail`).
//
// Kelola organisasi (Create dan ubah) di tabel datar `POOLDATA.CLIENT`, `CLIENT_PICLIST`, dan `CLIENT_ADDRESS` -
// perintah work owner 03/04-10-2026: modul baru `companydetail`, bentuk menu seperti layar Pega SFAGIS Company Detail,
// tanpa dokumen JSON ("aku tidak mau ada json lagi"), jangan di inti. Migrasi modul: `M_ENUMERASI` dan isinya
// (800-801), `NATION` dari view menjadi tabel datar (802-804), kolom tambahan tiga tabel client (805-807), dua trigger
// `M_CLIENT` dimatikan (808), dan slot menu 992.
package backend

import (
	"context"
	"embed"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/companydetail/backend/handlers"
	"nusantarare/modul/companydetail/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini, ditanam ke biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// Nama pengenal modul ini - SAMA dengan nama foldernya.
const Nama = "companydetail"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`). Nol kontrak lintas modul.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:    Nama,
		Migrasi: berkasMigrasi,
		// Akses menu LIHAT (keputusan work owner 04-10-2026): modul selesai, ikut gerbang tulis `cmd/api`.
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
