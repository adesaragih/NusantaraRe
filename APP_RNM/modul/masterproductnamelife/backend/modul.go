// Package backend merakit modul Master Product Name Life (`masterproductnamelife`).
//
// Nama paketnya `backend` - sama di setiap modul (nama folder); daftar modul
// bangkitan mengimpornya dengan alias nama modul. Modul ini diserahkan lewat
// `Pendaftaran()` (struktur tim satu folder per modul, 30-09-2026).
//
// ⭐ Sejak 02-10-2026 (K5, OQ-MPNL-01 flat; tiket 01 bab bertanggal): produk disimpan di tabel FLAT
// `M_PRODUCTNAME_LIFE` + tujuh anak - migrasi 140–147 (`TestMPNLRentangHanyaTabelFlat`); kedua tabel JSON warisan
// `M_PRODUCT_LIFE`/`M_PRODUCTINWARD_LIFE` hanya dibaca alat pindah `backend/alat/pindahflat`. Lampiran di
// `M_ATTACHMENTPRODUCTNAME` (`docs/STRUKTUR-TABEL-MASTER-PRODUCT-NAME-LIFE.md`). Slot menu 960 (`UPDATE DIMIGRASI`).
package backend

import (
	"context"
	"embed"
	"io/fs"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/masterproductnamelife/backend/handlers"
	"nusantarare/modul/masterproductnamelife/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini - tabel flat 140–147 dan slot menu 960, ditanam ke biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// SumberMigrasi menyerahkan folder `migrations/` modul ini kepada pelari
// migrasi (`inti/backend/migrasi`).
func SumberMigrasi() fs.FS { return berkasMigrasi }

// Nama pengenal modul ini di MODUL_AKTIF dan di GET /api/modul-aktif.
const Nama = "masterproductnamelife"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`,
// berkas bangkitan `modul_masterproductnamelife_gen.go`).
//
// Nol kontrak lintas modul: master rujukan dibaca langsung dari tabelnya,
// seperti RD Pega; outbox `T_LOG_SERVICE_RNM` lewat `inti/backend/outbox`.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:    Nama,
		Migrasi: berkasMigrasi,
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			return Baru(services.DariDasar(p.Dasar()), p.Config().AuthStub), nil
		},
	}
}

// Modul adalah perakitan modul Master Product Name Life untuk `cmd/api` -
// `inti.Modul`.
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

// JalankanPekerja - modul ini tidak punya pekerja latar: efek keluar lampiran
// dikirim seketika sesudah rekam dan diulang lewat rute `ulangi` (stub, P5).
func (Modul) JalankanPekerja(context.Context) inti.Pekerja { return inti.TanpaPekerja() }
