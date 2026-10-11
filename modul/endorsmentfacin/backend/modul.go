// Package backend merakit modul Endorsment Fac In (`endorsmentfacin`) - irisan I0 (pendaftaran) dan I2 (periksa
// sebelum Create), keputusan work owner 11-10-2026.
//
// Nama paketnya `backend` - sama di setiap modul (nama folder); daftar modul bangkitan mengimpornya dengan alias nama
// modul. Nol kontrak lintas modul: modul ini TIDAK mengimpor modul/nbfacin (penjaga impor lintas modul).
package backend

import (
	"context"
	"embed"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/endorsmentfacin/backend/handlers"
	"nusantarare/modul/endorsmentfacin/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini: 261 M_EDM_PENGECUALIAN_POLIS (tabel konfigurasi, diisi DBA)
// dan slot menu 966. Nomor 260 dicadangkan untuk T_WORK_POLIS/sequence kasus EDM- (irisan I3).
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// Nama pengenal modul ini - SAMA dengan nama foldernya.
const Nama = "endorsmentfacin"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`). Tanpa HakLihat: menu berakses penuh.
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
	svc *services.Service
	// stubPelaku - config AuthStub: penunda identitas X-Pelaku (pola modul lain).
	stubPelaku bool
}

// Nama menyebut modul ini.
func (Modul) Nama() string { return Nama }

// DaftarkanRute mendaftarkan seluruh rute modul ini ke mux bersama.
func (m Modul) DaftarkanRute(mux *http.ServeMux) { handlers.DaftarkanRute(mux, m.svc, m.stubPelaku) }

// JalankanPekerja - modul ini tidak punya pekerja latar.
func (Modul) JalankanPekerja(context.Context) inti.Pekerja { return inti.TanpaPekerja() }
