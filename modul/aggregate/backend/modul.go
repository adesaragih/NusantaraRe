// Package backend merakit modul Aggregate (`aggregate`).
//
// Unggah data aggregate per zona (CSV) ke tabel warisan `POOLDATA.AGGREGATE`, daftar dan rincian unggahan, hapus,
// dan chart RNM Value (USD) Ceding > Treaty Type > Coverage - perintah work owner 04-10-2026: modul Aggregate, kelompok MASTER
// TREATY, langkah XML Pega (folder korpus `Aggregate`) diikuti apa adanya. Migrasi modul: `SEQ_AGGREGATE` (880) dan
// slot menu 996; baris menunya dibuat migrasi inti 911.
package backend

import (
	"context"
	"embed"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/templat"
	"nusantarare/modul/aggregate/backend/handlers"
	"nusantarare/modul/aggregate/backend/models"
	"nusantarare/modul/aggregate/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini, ditanam ke biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// templatBawaan - berkas templat Upload CSV bawaan (header 38 kolom, pemisah `;`): dipakai selama Template Manager
// belum punya versi aktif untuk slot `KodeTemplat` (keputusan work owner 04-10-2026).
//
//go:embed templat/aggregate.csv
var templatBawaan []byte

// KodeTemplat - kode slot templat Upload CSV modul ini di Template Manager. TETAP: tersimpan di M_TEMPLATE_FILE.KODE.
const KodeTemplat = "aggregate.upload"

// jumlahKolomCSV - kolom CSV yang dibaca Upload CSV (`models.KolomGrid` ber-CSV).
func jumlahKolomCSV() int {
	n := 0
	for _, k := range models.KolomGrid {
		if k.CSV > 0 {
			n++
		}
	}
	return n
}

// Nama pengenal modul ini - SAMA dengan nama foldernya.
const Nama = "aggregate"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`). Nol kontrak lintas modul.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:    Nama,
		Migrasi: berkasMigrasi,
		Templat: []templat.Slot{{
			Kode: KodeTemplat, Menu: Nama, Grup: "Aggregate", Nama: "Upload CSV",
			DipakaiDi: "Aggregate › Add Data › Template", Ekstensi: ".csv", Pemisah: ';', JumlahKolom: jumlahKolomCSV(),
			NamaUnduhan: "aggregate.csv", Bawaan: templatBawaan,
		}},
		// Akses menu LIHAT (keputusan work owner 04-10-2026): modul selesai, ikut gerbang tulis `cmd/api`. Pratinjau hanya
		// membaca CSV, jadi bebas.
		HakLihat: &inti.HakLihat{Bebas: []string{"POST " + handlers.Prefix + "/pratinjau"}},
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
