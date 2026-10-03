// Package backend merakit modul Treaty Contract Out (`treatycontractout`).
//
// Nama paketnya `backend` - sama di setiap modul (nama folder); daftar modul
// bangkitan mengimpornya dengan alias nama modul. Sejak struktur tim satu
// folder per modul (30-09-2026) modul ini diserahkan lewat `Pendaftaran()`.
//
// Refactor bentuk B (30-09-2026): setiap modul punya satu berkas perakitan
// (`modul.go`) yang menyerahkan miliknya kepada `cmd/api`: rute HTTP-nya dan
// pekerja latar antrean lampirannya. Modul ini NOL migrasi - ia menulis dan
// membaca tabel warisan (tco4).
package backend

import (
	"context"
	"net/http"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/config"
	"nusantarare/modul/treatycontractout/backend/handlers"
	"nusantarare/modul/treatycontractout/backend/services"
)

// Nama pengenal modul ini di MODUL_AKTIF dan di GET /api/modul-aktif.
const Nama = "treatycontractout"

// Pendaftaran menyerahkan modul ini kepada perakit (`inti/backend/daftar`,
// berkas bangkitan `modul_treatycontractout_gen.go`).
//
// Nol kontrak lintas modul, nol migrasi (tco4: tabel warisan). Nama lama
// `treaty` ditolak MODUL_AKTIF dengan kalimat yang menyebut nama ini.
func Pendaftaran() inti.Pendaftaran {
	return inti.Pendaftaran{
		Nama:     Nama,
		NamaLama: []string{"treaty"},
		Bangun: func(p *inti.Perakitan) (inti.Modul, error) {
			cfg := p.Config()
			svc := services.DariDasar(p.Dasar()).
				// OQ-TCO-08: bawaan stub; ⛔ garam tidak pernah dicetak.
				DenganPenyimpananLampiranTCO(cfg.PelaksanaStorage == config.PelaksanaStorageNyata, cfg.StorageTokenSalt)
			return Baru(svc, cfg.AuthStub, cfg.IntervalPekerjaLampiranTCO, cfg.PelaksanaStorage, p.Catat), nil
		},
	}
}

// Modul adalah perakitan modul Treaty Contract Out untuk `cmd/api` -
// `inti.Modul`.
type Modul struct {
	svc        *services.Service
	stubPelaku bool
	// interval pekerja latar antrean lampiran (TCO_PEKERJA_LAMPIRAN_INTERVAL);
	// nol = pekerja MATI.
	interval time.Duration
	// pelaksana penyimpanan lampiran (PELAKSANA_STORAGE), untuk dicatat saja.
	pelaksana string
	catat     func(string)
}

// Baru merakit modul di atas Service yang sudah disambung (`Pendaftaran`).
func Baru(svc *services.Service, stubPelaku bool, interval time.Duration, pelaksana string,
	catat func(string)) Modul {
	return Modul{svc: svc, stubPelaku: stubPelaku, interval: interval, pelaksana: pelaksana, catat: catat}
}

// Nama menyebut modul ini.
func (Modul) Nama() string { return Nama }

// DaftarkanRute mendaftarkan seluruh rute modul ini ke mux bersama.
func (m Modul) DaftarkanRute(mux *http.ServeMux) { handlers.DaftarkanRute(mux, m.svc, m.stubPelaku) }

// JalankanPekerja menyalakan pekerja latar antrean lampiran Treaty Contract
// Out (OQ-TCO-09, keputusan work owner 29-09-2026).
//
// Refactor bentuk B (30-09-2026): dipindah apa adanya dari `cmd/api`
// (`jalankanPekerjaLampiranTCO`), kalimat catatannya sama. Mati bila
// TCO_PEKERJA_LAMPIRAN_INTERVAL kosong/0 atau tanpa Oracle. Ia berhenti
// bersama ctx proses; `Selesai` tertutup saat ia benar-benar berhenti
// (langsung tertutup bila tidak dinyalakan).
func (m Modul) JalankanPekerja(ctx context.Context) inti.Pekerja {
	m.catat("treaty contract out attachments: storage executor " + m.pelaksana)
	selesai := make(chan struct{})
	pekerja := inti.Pekerja{
		Selesai:        selesai,
		PesanTerlambat: "treaty contract out attachments: background worker had not stopped by the shutdown deadline",
	}
	if m.interval <= 0 {
		m.catat("treaty contract out attachments: background worker off (interval empty)")
		close(selesai)
		return pekerja
	}
	if !m.svc.PunyaDatabase() {
		m.catat("treaty contract out attachments: background worker off (no oracle)")
		close(selesai)
		return pekerja
	}
	m.catat("treaty contract out attachments: background worker every " + m.interval.String())
	go func() {
		defer close(selesai)
		handlers.LayananLampiranTCO(m.svc).JalankanPekerja(ctx, m.interval, m.catat)
	}()
	return pekerja
}
