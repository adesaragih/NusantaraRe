// Package backend merakit modul Treaty In Adjustment (`treatyinadjustment`).
//
// ⛔ LINGKUP HARI INI: tiket 01 dan 05 saja - kolom `ID_VERSI_KONTRAK_DASAR`
// beserta kunci asingnya, dan pelonggaran `NOMOR_URUT_VERSI` menjadi boleh
// kosong (langkah PERLUAS). Ditambah dua jalur baca atas rantai versi.
//
// ⛔ TIKET 02 DAN 03 TIDAK PUNYA SISA PEKERJAAN SKEMA, dan itu temuan, bukan
// kelalaian: `SIFAT_MATERIAL_ADDENDUM` dan `TANGGAL_BERLAKU_ADDENDUM` sudah
// berdiri sejak migrasi 401 modul `treatyin`, sebab `KAMUS-KOLOM.md` adalah
// SATU model untuk kedua modul - yang dipisahkan 25-09-2026 hanya papan
// tiketnya. Yang tersisa dari keduanya adalah PENEGAKAN di sisi simpan, dan
// jalur simpan belum berspesifikasi (`L-4`).
//
// ⛔ TIKET 04 TERTAHAN `DB-16a`, pertanyaan bisnis yang belum dikirim.
// `DOKUMEN_ADDENDUM` karena itu belum dibuat, dan kunci asing
// `VERSI_KONTRAK.ID_DOKUMEN_ADDENDUM` belum dipasang (lihat migrasi 401 modul
// `treatyin`).
//
// Keputusan atas pertentangan tiket 05 lawan tiket 14 ada di dalam migrasi
// 441; keputusan atas trigger `INV-54` tiket 03 ada di
// `docs/KEPUTUSAN-TIKET-02-03.md`.
package backend

import (
	"context"
	"embed"
	"io/fs"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyinadjustment/backend/handlers"
	"nusantarare/modul/treatyinadjustment/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini, ditanam ke biner:
// rentang 440-479 dan slot menu 974.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// SumberMigrasi menyerahkan folder `migrations/` modul ini kepada pelari
// migrasi (`inti/backend/migrasi`).
func SumberMigrasi() fs.FS { return berkasMigrasi }

// Nama pengenal modul ini di MODUL_AKTIF dan di GET /api/modul-aktif.
const Nama = "treatyinadjustment"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`,
// berkas bangkitan `modul_treatyinadjustment_gen.go`).
//
// Nol kontrak lintas modul: modul ini membaca TABEL yang modul `treatyin`
// buat, dan tabel bukan kontrak. Ia tidak mengimpor satu pun paket modul itu.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:    Nama,
		Migrasi: berkasMigrasi,
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			return Baru(services.DariDasar(p.Dasar()), p.Config().AuthStub), nil
		},
	}
}

// Modul adalah perakitan modul Treaty In Adjustment untuk `cmd/api`.
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
