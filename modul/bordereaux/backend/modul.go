// Package backend merakit modul Bordereaux (`bordereaux`).
//
// Berkas bordereaux premi, klaim, dan subrogasi per Type x Business (29 kombinasi), Upload CSV, Save tanpa JSON,
// persetujuan Maker > Checker > Supervisor dengan riwayat - perintah work owner 04-10-2026: folder korpus
// `D:\XML\RNM_BRD\Bordereaux`, kelompok MASTER TREATY; XML patokan, bug Pega diperbaiki di Go. Migrasi modul:
// `BORDEREAUX_HISTORY` (890), tiga workbasket (891), slot menu 998; baris menunya dibuat migrasi inti 913.
package backend

import (
	"context"
	"embed"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/templat"
	"nusantarare/modul/bordereaux/backend/handlers"
	"nusantarare/modul/bordereaux/backend/models"
	"nusantarare/modul/bordereaux/backend/services"
)

// berkasMigrasi adalah folder `migrations/` modul ini, ditanam ke biner.
//
//go:embed migrations/*.sql
var berkasMigrasi embed.FS

// berkasTemplat - 29 templat Upload CSV bawaan, BARIS HEADER SAJA (data contoh work owner tidak masuk repositori;
// berkas lengkapnya diunggah lewat Template Manager).
//
//go:embed templat/*.csv
var berkasTemplat embed.FS

// Nama pengenal modul ini - SAMA dengan nama foldernya.
const Nama = "bordereaux"

// AwalanTemplat - awalan kode slot Template Manager modul ini (`bordereaux.premi.fire`, ...). TETAP.
const AwalanTemplat = services.AwalanTemplat

// SlotTemplat - satu slot Template Manager per kombinasi Type x Business (`DownloadTemplate_Act` langkah 2.1-2.29).
func SlotTemplat() []templat.Slot {
	slot := make([]templat.Slot, 0, len(models.Kombinasi))
	for _, k := range models.Kombinasi {
		isi, err := berkasTemplat.ReadFile("templat/" + k.BerkasTemplat)
		if err != nil {
			panic("bordereaux: templat bawaan hilang: " + k.BerkasTemplat)
		}
		slot = append(slot, templat.Slot{
			Kode: AwalanTemplat + k.Kode, Menu: Nama, Grup: "Bordereaux", Nama: k.Type + " · " + k.Business,
			DipakaiDi: "Bordereaux › Input Data › Details › Template", Ekstensi: ".csv", Pemisah: ';',
			JumlahKolom: len(k.Kolom), NamaUnduhan: k.BerkasTemplat, Bawaan: isi,
		})
	}
	return slot
}

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`). Nol kontrak lintas modul.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:    Nama,
		Migrasi: berkasMigrasi,
		Templat: SlotTemplat(),
		// Akses menu LIHAT (keputusan work owner 04-10-2026): PENUH = Input Data/Edit/Delete/Submit (pengganti
		// workbasket ReasBordereauxAdmin). Bebas: Upload CSV hanya membaca; Submit dipakai juga Checker/Supervisor
		// ber-hak LIHAT - layanan menegakkan siapa boleh apa (`services.HakAtas`).
		HakLihat: &inti.HakLihat{Bebas: []string{"POST " + handlers.Prefix + "/unggah-csv", "POST " + handlers.Prefix + "/berkas/{id}/submit"}},
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

// Nama memenuhi inti.Modul.
func (Modul) Nama() string { return Nama }

// DaftarkanRute memenuhi inti.Modul.
func (m Modul) DaftarkanRute(mux *http.ServeMux) { handlers.DaftarkanRute(mux, m.svc, m.stubPelaku) }

// JalankanPekerja memenuhi inti.Modul - pekerja kode pos AI menyusul (tahap 6).
func (Modul) JalankanPekerja(context.Context) inti.Pekerja { return inti.TanpaPekerja() }
