package services_test

// Porsi periode tahap pembayaran - `CountPaymentEdm_Act` blok 13-15 (tiket E13,
// lanjutan E06). Dibaca sesudah: services/porsipembayaran.go.

import (
	"errors"
	"testing"

	"nusantarare/inti/backend/uang"
	"nusantarare/modul/endorsmentfacin/backend/models"
	"nusantarare/modul/endorsmentfacin/backend/services"
)

func porsiBayar(t *testing.T, ubah func(*services.MasukanPorsiPembayaran)) services.HasilPorsiPembayaran {
	t.Helper()
	// Tanggal polis 00:00 GMT → langkah 4 → 05:00 GMT; EdmDate (tidak
	// dinormalisasi) diberi 05:00 supaya selisihnya hari bulat.
	m := services.MasukanPorsiPembayaran{
		Predikat:      services.PredikatPembayaran{IsEDM: true},
		EDMDay:        "365",
		StartDateTime: tanggal(t, "2026-01-01 00:00"), EndDateTime: tanggal(t, "2026-01-05 00:00"),
		EdmDate:       tanggal(t, "2026-01-02 05:00"),
		CacahMataUang: 1,
	}
	ubah(&m)
	h, err := services.HitungPorsiPembayaran(m)
	if err != nil {
		t.Fatalf("HitungPorsiPembayaran: %v", err)
	}
	return h
}

func rasioSatu(t *testing.T) uang.Ratio { return rasio(t, "1") }

// TestExtendPeriodPorsiSatu - blok 13 (IsEdmExtendPeriod): datedif = day = 1,
// ProrateEDMEnd = 1; ProrateStartEDM TIDAK ditulis.
func TestExtendPeriodPorsiSatu(t *testing.T) {
	h := porsiBayar(t, func(m *services.MasukanPorsiPembayaran) { m.Predikat.IsEdmExtendPeriod = true })
	if h.ProrateEDMEnd.Kosong() || h.ProrateEDMEnd.Value.Cmp(rasioSatu(t).Value) != 0 || !h.ProrateStartEDM.Kosong() {
		t.Errorf("%+v", h)
	}
	if h.Datedif != 1 || h.Day != 1 {
		t.Errorf("datedif %d day %d", h.Datedif, h.Day)
	}
}

// TestAdjRatePorsiDariSelisihHari - blok 14: datedif = hari Edm→akhir,
// datedifbefore = hari mulai→Edm, day = hari mulai→akhir; porsi HALF_UP 20.
func TestAdjRatePorsiDariSelisihHari(t *testing.T) {
	h := porsiBayar(t, func(m *services.MasukanPorsiPembayaran) { m.Predikat.IsEdmAdjRate = true })
	samaRasio(t, "ProrateEDMEnd", h.ProrateEDMEnd, uang.Ratio{Value: rasio(t, "0.75").Value, Scale: 20})
	samaRasio(t, "ProrateStartEDM", h.ProrateStartEDM, uang.Ratio{Value: rasio(t, "0.25").Value, Scale: 20})
	if h.Datedif != 3 || h.DatedifBefore != 1 || h.Day != 4 {
		t.Errorf("datedif %d before %d day %d", h.Datedif, h.DatedifBefore, h.Day)
	}
}

// TestAdjShareCedantFixRate - langkah 9 (IsMarineCargo ∨ IsEdmAdjShareCedant →
// "FIX RATE") berjalan SEBELUM blok 14, sehingga Adj Share Cedant memasuki 14.3:
// porsi 1 dan 0, bukan dari selisih hari.
func TestAdjShareCedantFixRate(t *testing.T) {
	h := porsiBayar(t, func(m *services.MasukanPorsiPembayaran) { m.Predikat.IsEdmAdjShareCedant = true })
	if !h.FixRate {
		t.Error("langkah 9 tidak menyalakan FIX RATE")
	}
	samaRasio(t, "ProrateEDMEnd", h.ProrateEDMEnd, uang.Ratio{Value: rasio(t, "1").Value, Scale: 20})
	samaRasio(t, "ProrateStartEDM", h.ProrateStartEDM, uang.Ratio{Value: rasio(t, "0").Value, Scale: 20})
}

// TestBukanEDMKeluar - langkah 1 IsEDM F=6: keluar, tidak menulis apa pun.
func TestBukanEDMKeluar(t *testing.T) {
	h := porsiBayar(t, func(m *services.MasukanPorsiPembayaran) {
		m.Predikat = services.PredikatPembayaran{IsEdmExtendPeriod: true}
	})
	if h.Ditulis || !h.StartDateTime.IsZero() {
		t.Errorf("%+v", h)
	}
}

// TestNormalisasiTanggalPolis - langkah 4: tanggal kalender Jakarta pukul 05:00
// GMT. 2026-01-01 18:00 GMT = 2026-01-02 01:00 WIB → 2026-01-02 05:00 GMT.
func TestNormalisasiTanggalPolis(t *testing.T) {
	h := porsiBayar(t, func(m *services.MasukanPorsiPembayaran) {
		m.StartDateTime = tanggal(t, "2025-12-31 18:00")
	})
	if !h.StartDateTime.Equal(tanggal(t, "2026-01-01 05:00")) || !h.EndDateTime.Equal(tanggal(t, "2026-01-05 05:00")) {
		t.Errorf("tanggal %v / %v", h.StartDateTime, h.EndDateTime)
	}
	_, err := services.HitungPorsiPembayaran(services.MasukanPorsiPembayaran{Predikat: services.PredikatPembayaran{IsEDM: true}})
	if !errors.Is(err, services.ErrTanggalPorsiPeriodeKosong) {
		t.Errorf("tanggal kosong: %v", err)
	}
}

// TestEDMDay - kosong sah (15.2 → 365); teks bukan bulat → galat.
func TestEDMDay(t *testing.T) {
	h := porsiBayar(t, func(m *services.MasukanPorsiPembayaran) { m.EDMDay = "" })
	if !h.DayKosong {
		t.Error("EDMDay kosong tidak dibawa sebagai kosong")
	}
	h = porsiBayar(t, func(m *services.MasukanPorsiPembayaran) { m.EDMDay, m.Predikat.IsEdmAdjPeriod = "", true })
	if h.Day != 365 || h.DayKosong {
		t.Errorf("15.2: day %d kosong %v", h.Day, h.DayKosong)
	}
	_, err := services.HitungPorsiPembayaran(services.MasukanPorsiPembayaran{
		Predikat: services.PredikatPembayaran{IsEDM: true}, EDMDay: "abc",
		StartDateTime: tanggal(t, "2026-01-01 00:00"), EndDateTime: tanggal(t, "2026-01-05 00:00"),
	})
	if !errors.Is(err, services.ErrEDMDayBukanBulat) {
		t.Errorf("EDMDay abc: %v", err)
	}
}

// TestAdjRateFixRateAtauMarine - 14.3: FIX RATE atau marine cargo → porsi 1,
// porsi sebelum 0.
func TestAdjRateFixRateAtauMarine(t *testing.T) {
	for nama, ubah := range map[string]func(*services.MasukanPorsiPembayaran){
		"fix rate": func(m *services.MasukanPorsiPembayaran) { m.Predikat.IsEdmAdjRate, m.FixRate = true, true },
		"marine": func(m *services.MasukanPorsiPembayaran) {
			m.Predikat.IsEdmAdjRate, m.Predikat.IsMarineCargo = true, true
		},
	} {
		h := porsiBayar(t, ubah)
		samaRasio(t, nama+" ProrateEDMEnd", h.ProrateEDMEnd, uang.Ratio{Value: rasio(t, "1").Value, Scale: 20})
		samaRasio(t, nama+" ProrateStartEDM", h.ProrateStartEDM, uang.Ratio{Value: rasio(t, "0").Value, Scale: 20})
	}
}

// TestAdjRateTanpaBarisMataUang - 14.4 mengulang CurrencyList; tanpa baris,
// 14.4.5 tidak berjalan: ProrateEDMEnd tetap 1 dari 14.1, ProrateStartEDM tidak
// ditulis.
func TestAdjRateTanpaBarisMataUang(t *testing.T) {
	h := porsiBayar(t, func(m *services.MasukanPorsiPembayaran) { m.Predikat.IsEdmAdjRate, m.CacahMataUang = true, 0 })
	if h.ProrateEDMEnd.Value.Cmp(rasioSatu(t).Value) != 0 || !h.ProrateStartEDM.Kosong() {
		t.Errorf("%+v", h)
	}
}

// TestAdjPeriodSelaluSatu - blok 15: 15.3 ber-`pyStepsPreCondition=false`
// (tetap jalan, P-11) menyetel startdate = edmdate = 1, sehingga 15.4
// `@Math.divide(edmdate, startdate, 20)` SELALU 1, apa pun tanggalnya - K-046,
// diport apa adanya.
func TestAdjPeriodSelaluSatu(t *testing.T) {
	h := porsiBayar(t, func(m *services.MasukanPorsiPembayaran) {
		m.Predikat.IsEdmAdjPeriod = true
		m.StartDateTime, m.EndDateTime = tanggal(t, "2026-01-01 00:00"), tanggal(t, "2026-01-31 00:00")
	})
	samaRasio(t, "ProrateEDMEnd", h.ProrateEDMEnd, uang.Ratio{Value: rasio(t, "1").Value, Scale: 20})
}

// TestHariTakBulatDitolak - `@DateTimeDifference(…,"D")` ke lokal `int`: cara
// Pega memotong pecahan hari belum terverifikasi (A09).
func TestHariTakBulatDitolak(t *testing.T) {
	_, err := services.HitungPorsiPembayaran(services.MasukanPorsiPembayaran{
		Predikat:      services.PredikatPembayaran{IsEDM: true, IsEdmAdjRate: true},
		StartDateTime: tanggal(t, "2026-01-01 00:00"), EndDateTime: tanggal(t, "2026-01-05 00:00"),
		EdmDate: tanggal(t, "2026-01-02 12:00"), CacahMataUang: 1,
	})
	if !errors.Is(err, services.ErrSatuanSelisihWaktuBelumTerverifikasi) {
		t.Errorf("galat %v", err)
	}
}

// TestTanpaCabangTidakMenulis - jenis endorsement lain: blok 13-15 tidak
// menulis porsi (nilainya dari CountPremiEDM_DT / CountPaymentEdmTSIObj_Act,
// belum diport).
func TestTanpaCabangTidakMenulis(t *testing.T) {
	h := porsiBayar(t, func(*services.MasukanPorsiPembayaran) {})
	if !h.ProrateEDMEnd.Kosong() || !h.ProrateStartEDM.Kosong() || h.Ditulis {
		t.Errorf("%+v", h)
	}
}

// TestTSIObjHariInklusif - CountPaymentEdmTSIObj_Act blok 1: (hari Edm→akhir
// + 1) ÷ hari mulai→akhir. Periode 4 hari, endorsement hari ke-1: (3+1)/4 = 1.
// K-046 - porsi dapat menyentuh/melebihi 1.
func TestTSIObjHariInklusif(t *testing.T) {
	h := porsiBayar(t, func(m *services.MasukanPorsiPembayaran) { m.Predikat.IsEdmAdjTSI = true })
	samaRasio(t, "ProrateEDMEnd", h.ProrateEDMEnd, uang.Ratio{Value: rasio(t, "1").Value, Scale: 20})
	// Endorsement pada tanggal mulai: (4+1)/4 = 1.25 - melebihi 1.
	h = porsiBayar(t, func(m *services.MasukanPorsiPembayaran) {
		m.Predikat.IsEdmAddObject, m.EdmDate = true, tanggal(t, "2026-01-01 05:00")
	})
	samaRasio(t, "ProrateEDMEnd > 1", h.ProrateEDMEnd, uang.Ratio{Value: rasio(t, "1.25").Value, Scale: 20})
	// ShortPeriod: 1/1.
	h = porsiBayar(t, func(m *services.MasukanPorsiPembayaran) { m.Predikat.IsEdmAdjRIC, m.IsProRate = true, "ShortPeriod" })
	samaRasio(t, "ShortPeriod", h.ProrateEDMEnd, uang.Ratio{Value: rasio(t, "1").Value, Scale: 20})
}

// TestTSIObjBatalProrata - langkah 3: EdmType 2, per baris mata uang.
func TestTSIObjBatalProrata(t *testing.T) {
	h := porsiBayar(t, func(m *services.MasukanPorsiPembayaran) { m.EdmType = models.EdmBatalProrata })
	samaRasio(t, "ProrateEDMEnd", h.ProrateEDMEnd, uang.Ratio{Value: rasio(t, "1").Value, Scale: 20})
	if h := porsiBayar(t, func(m *services.MasukanPorsiPembayaran) { m.EdmType, m.CacahMataUang = models.EdmBatalProrata, 0 }); h.Ditulis {
		t.Error("tanpa baris mata uang, langkah 3 menulis")
	}
}

// TestBlokSesudahMenimpaTSIObj - urutan CountPaymentEdm_Act: langkah 11
// (TSIObj) SEBELUM blok 13; Extend Period menimpa jadi 1.
func TestBlokSesudahMenimpaTSIObj(t *testing.T) {
	h := porsiBayar(t, func(m *services.MasukanPorsiPembayaran) {
		m.Predikat.IsEdmAddObject, m.Predikat.IsEdmExtendPeriod = true, true
		m.EdmDate = tanggal(t, "2026-01-01 05:00")
	})
	if h.ProrateEDMEnd.Value.Cmp(rasio(t, "1").Value) != 0 {
		t.Errorf("ProrateEDMEnd %s, mau 1 dari blok 13", h.ProrateEDMEnd)
	}
}
