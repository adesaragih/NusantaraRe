// Package backend merakit modul Treaty Description (`treatydescription`).
//
// Kelola tabel warisan `POOLDATA.TREATYDESC` (master jenis klausul treaty) - tambah dan ubah, nol hapus (keputusan
// work owner 05-10-2026). Modul di luar korpus: Pega hanya membaca tabel ini (`BrowseTreatyDesc_RD`) dan menulisnya
// lewat prosedur `PEGA_TREATYDESC` yang tidak dipanggil di sini. Nol tabel baru, nol DDL, nol migrasi sendiri
// (MODUL.md `—`): baris menunya dibuat dan langsung dinyalakan migrasi inti 920.
package backend

import (
	"context"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatydescription/backend/handlers"
	"nusantarare/modul/treatydescription/backend/services"
)

// Nama pengenal modul ini - SAMA dengan nama foldernya.
const Nama = "treatydescription"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`). Nol kontrak lintas modul, nol migrasi.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama: Nama,
		// Akses menu LIHAT: ikut gerbang tulis `cmd/api` seperti modul master lain.
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
