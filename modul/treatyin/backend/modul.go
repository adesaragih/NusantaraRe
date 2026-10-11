// Package backend merakit modul Treaty In (`treatyin`).
//
// Nama paketnya `backend` - sama di setiap modul (nama folder); daftar modul
// bangkitan mengimpornya dengan alias nama modul. Modul ini diserahkan lewat
// `Pendaftaran()` (struktur tim satu folder per modul, 30-09-2026).
//
// ⛔ LINGKUP HARI INI: tiket 14 dan 15 saja - skema identitas kontrak
// (`KONTRAK`, `VERSI_KONTRAK`) dan keenam tabel acuan, ditambah SATU jalur baca
// atas tabel acuan itu. Papan tiket modul ini (`docs/issues/README.md`)
// menyatakan `L-4`: "tidak ada spesifikasi layar di mana pun", dan melarang
// mengarang layar untuk memenuhi bentuk irisan tegak. Ke-49 tiket sisanya -
// layer, bagian, penyebaran, persetujuan, warisan - LAHIR BERSAMA
// spesifikasinya.
//
// Empat penyelarasan spec dengan repo (skema, presisi NUMBER, nama tabel,
// sequence yang hilang dari `ddl-usulan/`) tercatat di
// `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
package backend

import (
	"context"
	"embed"
	"io/fs"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/handlers"
	"nusantarare/modul/treatyin/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini, ditanam ke biner:
// rentang tabel 400-439 dan slot menu 972.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// SumberMigrasi menyerahkan folder `migrations/` modul ini kepada pelari
// migrasi (`inti/backend/migrasi`).
func SumberMigrasi() fs.FS { return berkasMigrasi }

// Nama pengenal modul ini di MODUL_AKTIF dan di GET /api/modul-aktif.
const Nama = "treatyin"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`,
// berkas bangkitan `modul_treatyin_gen.go`).
//
// Nol kontrak lintas modul hari ini. Tiket 41 ("konsumen hilir membaca
// identitas kontrak dalam bentuk lama") adalah yang pertama menuntutnya, dan
// ia belum dikerjakan.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:    Nama,
		Migrasi: berkasMigrasi,
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			// ⭐ Unggah panel Attachment (8 Oktober 2026) — `ServiceGoogle`
			// seperti XML; garam hanya untuk token penyimpanan BARU.
			// ⛔ Garam tidak pernah dicetak.
			svc := services.DariDasar(p.Dasar()).DenganGaramToken(p.Config().StorageTokenSalt)
			return Baru(svc, p.Config().AuthStub), nil
		},
	}
}

// Modul adalah perakitan modul Treaty In untuk `cmd/api` - `inti.Modul`.
type Modul struct {
	svc        *services.Service
	stubPelaku bool
}

// Baru merakit modul di atas Service yang sudah disambung (`Pendaftaran`).
func Baru(svc *services.Service, stubPelaku bool) Modul {
	return Modul{svc: svc, stubPelaku: stubPelaku}
}

// Nama menyebut modul ini.
func (Modul) Nama() string { return Nama }

// DaftarkanRute mendaftarkan seluruh rute modul ini ke mux bersama.
func (m Modul) DaftarkanRute(mux *http.ServeMux) { handlers.DaftarkanRute(mux, m.svc, m.stubPelaku) }

// JalankanPekerja - modul ini tidak punya pekerja latar.
func (Modul) JalankanPekerja(context.Context) inti.Pekerja { return inti.TanpaPekerja() }
