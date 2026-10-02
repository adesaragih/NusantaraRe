package services_test

// Validasi tanggal endorsement - tiket E05. Nomor polis SINTETIS.
//
// Dibaca sesudah: services/validasitanggal.go.

import (
	"errors"
	"testing"
	"time"

	"nusantarare/modul/endorsmentfacin/backend/services"
)

func gmt(t *testing.T, s string) time.Time {
	t.Helper()
	w, err := time.Parse("2006-01-02 15:04", s)
	if err != nil {
		t.Fatal(err)
	}
	return w.UTC()
}

// periodeUji - polis 1 Okt 2025 s.d. 1 Mar 2026 WIB, tersimpan 17:00 GMT hari
// sebelumnya (bentuk yang terlihat di fixture NB-15).
func periodeUji(t *testing.T) *services.PeriodeProduksi {
	t.Helper()
	return &services.PeriodeProduksi{
		BeginDate: gmt(t, "2025-09-30 17:00"),
		EndDate:   gmt(t, "2026-02-28 17:00"),
	}
}

func validasi(t *testing.T, edmWIB string, p *services.PeriodeProduksi, polis string, kecuali ...string) services.HasilValidasiTanggal {
	t.Helper()
	wib := time.FixedZone("WIB", 7*3600)
	w, err := time.ParseInLocation("2006-01-02 15:04", edmWIB, wib)
	if err != nil {
		t.Fatal(err)
	}
	h, err := services.ValidasiTanggalEndorsemen(services.MasukanValidasiTanggal{
		PolicyNo: polis, TanggalEndorsement: w, Periode: p, Pengecualian: services.DaftarPolis(kecuali),
	})
	if err != nil {
		t.Fatalf("galat sistem: %v", err)
	}
	return h
}

func TestTanggalDiDalamPeriodeLolos(t *testing.T) {
	if h := validasi(t, "2025-12-15 10:00", periodeUji(t), "UJI-POLIS-1"); h.Ditolak {
		t.Errorf("ditolak: %+v", h)
	}
}

// TestTanggalDiLuarPeriodeDitolak - galat yang dapat dibaca petugas, sebagai
// keluaran bisnis.
func TestTanggalDiLuarPeriodeDitolak(t *testing.T) {
	for _, edm := range []string{"2025-09-01 10:00", "2026-03-05 10:00"} {
		h := validasi(t, edm, periodeUji(t), "UJI-POLIS-1")
		if !h.Ditolak || h.Pesan != "EDM date cannot be outside the period" || h.Medan != "EndorsementDate" {
			t.Errorf("%s: %+v", edm, h)
		}
	}
}

// TestTanggalSamaDenganTanggalMulaiLolos - langkah 5 baris 1: tanggal EDM sama
// dengan tanggal mulai (format dd/MM/yyyy) dilewati tanpa galat.
//
// ⚠️ Diport apa adanya: tanggal mulai dibaca dari tanggal kalender GMT
// (`substring(8,11) → "T05"`), sehingga untuk polis yang tersimpan 17:00 GMT
// hari SEBELUMNYA, tanggal "mulai" yang dibandingkan adalah sehari lebih awal
// dari tanggal mulai WIB-nya.
func TestTanggalSamaDenganTanggalMulaiLolos(t *testing.T) {
	if h := validasi(t, "2025-09-30 09:00", periodeUji(t), "UJI-POLIS-1"); h.Ditolak {
		t.Errorf("tanggal = tanggal mulai versi GMT ditolak: %+v", h)
	}
}

// TestPengecualianDariKonfigurasi - langkah 2: nomor polis di daftar
// konfigurasi, atau yang memuat "RNML", melompat ke label END.
func TestPengecualianDariKonfigurasi(t *testing.T) {
	if h := validasi(t, "2026-03-05 10:00", periodeUji(t), "UJI-POLIS-X", "UJI-POLIS-X"); h.Ditolak || !h.Dilewati {
		t.Errorf("pengecualian konfigurasi: %+v", h)
	}
	if h := validasi(t, "2026-03-05 10:00", periodeUji(t), "UJI-RNML-1"); h.Ditolak || !h.Dilewati {
		t.Errorf("pola RNML: %+v", h)
	}
	// Pengecualian tidak memerlukan data periode.
	if h := validasi(t, "2026-03-05 10:00", nil, "UJI-POLIS-X", "UJI-POLIS-X"); h.Ditolak {
		t.Errorf("pengecualian tanpa periode: %+v", h)
	}
}

// TestPeriodeTidakAdaBukanTebakan - query periode tanpa baris: perilaku Pega
// (`@toDate("")`) belum terverifikasi → galat sistem, bukan lolos atau tolak.
func TestPeriodeTidakAdaBukanTebakan(t *testing.T) {
	_, err := services.ValidasiTanggalEndorsemen(services.MasukanValidasiTanggal{
		PolicyNo: "UJI-POLIS-1", TanggalEndorsement: gmt(t, "2025-12-15 03:00"),
	})
	if !errors.Is(err, services.ErrPeriodeProduksiTidakAda) {
		t.Errorf("galat %v, mau ErrPeriodeProduksiTidakAda", err)
	}
}
