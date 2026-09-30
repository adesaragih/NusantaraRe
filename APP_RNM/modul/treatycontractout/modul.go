// Package treatycontractout merakit modul Treaty Contract Out.
//
// Refactor bentuk B (30-09-2026): setiap modul punya satu berkas perakitan
// (`modul.go`) yang menyerahkan miliknya kepada `cmd/api`: rute HTTP-nya dan
// pekerja latar antrean lampirannya. Modul ini NOL migrasi - ia menulis dan
// membaca tabel warisan (tco4).
package treatycontractout

import (
	"context"
	"net/http"
	"time"

	"nusantarare/inti"
	"nusantarare/modul/treatycontractout/handlers"
	"nusantarare/modul/treatycontractout/services"
)

// Nama pengenal modul ini di MODUL_AKTIF dan di GET /api/modul-aktif.
const Nama = "treatycontractout"

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

// Baru merakit modul di atas Service yang sudah disambung `modul.Rakit`.
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
