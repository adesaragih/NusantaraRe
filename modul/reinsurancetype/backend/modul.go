// Package backend merakit modul Reinsurance Type (`reinsurancetype`).
//
// Kelola tabel warisan `POOLDATA.REINSURANCETYPE` - master jenis reasuransi; tambah dan ubah, nol hapus (nonaktif =
// Flag inactive; perintah work owner 05-10-2026: "select * from reinsurancetype"). Modul di luar korpus: Pega hanya
// membaca tabel ini. Nol tabel baru, nol DDL, dan NOL MIGRASI SENDIRI - nomor modul sudah habis terbagi dan modul lain
// tidak boleh disentuh; baris menunya dibuat migrasi inti 921 langsung menyala.
package backend

import (
	"context"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/reinsurancetype/backend/handlers"
	"nusantarare/modul/reinsurancetype/backend/services"
)

// Nama pengenal modul ini - SAMA dengan nama foldernya.
const Nama = "reinsurancetype"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`). Nol kontrak lintas modul, nol migrasi.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama: Nama,
		// Akses menu LIHAT (keputusan work owner 05-10-2026): ikut gerbang tulis `cmd/api`.
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
