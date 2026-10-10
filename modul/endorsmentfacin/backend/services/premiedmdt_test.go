package services_test

// CountPremiEDM_DT - tiket E14 (sebagian) / E13. Angka SINTETIS; baris baru /
// lama sama dengan pembayaran_test.go (netto baru 865).
// Dibaca sesudah: services/premiedmdt.go.

import (
	"errors"
	"testing"

	"nusantarare/inti/backend/uang"
	"nusantarare/modul/endorsmentfacin/backend/models"
	"nusantarare/modul/endorsmentfacin/backend/services"
)

// premiDT - polis baru = polis lama (mulai 01-01, akhir 01-05, 05:00 GMT),
// EdmDate 01-02: WHEN A saja → startdate 4, edmdate 3+1 = 4 → porsi akhir 1,
// awal 0.
func premiDT(t *testing.T, ubah func(*services.MasukanPremiEDMDT)) services.HasilPremiEDMDT {
	t.Helper()
	h, err := services.HitungPremiEDMDT(masukanDT(t, ubah))
	if err != nil {
		t.Fatalf("HitungPremiEDMDT: %v", err)
	}
	return h
}

func masukanDT(t *testing.T, ubah func(*services.MasukanPremiEDMDT)) services.MasukanPremiEDMDT {
	t.Helper()
	m := services.MasukanPremiEDMDT{
		StartDateTime: tanggal(t, "2026-01-01 05:00"), EndDateTime: tanggal(t, "2026-01-05 05:00"),
		EdmDate:   tanggal(t, "2026-01-02 05:00"),
		LamaStart: tanggal(t, "2026-01-01 05:00"), LamaEnd: tanggal(t, "2026-01-05 05:00"),
		LamaEdmDate: tanggal(t, "2026-01-03 05:00"),
		EdmType:     models.EdmPerubahan, RICommision: rasio(t, "10"), HasilProdKe: "0",
		MataUang: []services.BarisBayar{baruUji(t)}, MataUangLama: []services.BarisBayar{lamaUji(t)},
	}
	ubah(&m)
	return m
}

// TestDTPerubahan - WHEN EdmType==4: OldPremi lama TANPA Deduction2
// (800 − 80 − 40 + 16 + 4 = 700); EDMOldPremi = Premium lama.
func TestDTPerubahan(t *testing.T) {
	h := premiDT(t, func(*services.MasukanPremiEDMDT) {})
	samaRasio(t, "ProrateEDMEnd", h.ProrateEDMEnd, uang.Ratio{Value: rasio(t, "1").Value, Scale: 20})
	b := h.MataUang[0]
	samaU(t, "EDMNewPremi", b.Payment.EDMNewPremi, "865")
	samaU(t, "EDMPremiMenjadi", b.Payment.EDMPremiMenjadi, "865")
	samaU(t, "EDMOldPremi", b.Payment.EDMOldPremi, "800")
	samaU(t, "EDMOldPayment", b.Payment.EDMOldPayment, "700")
	samaU(t, "EDMOldCommision", b.Payment.EDMOldCommision, "80")
	samaU(t, "EDMOldDeduction2", b.Payment.EDMOldDeduction2, "8")
	samaU(t, "SumTotalPayment", b.SumTotalPayment, "165")
	// CARI24 = 1 (bukan polis NB/RNW): OldPremi = EDMPremiMenjadi lama.
	h = premiDT(t, func(m *services.MasukanPremiEDMDT) {
		m.HasilProdKe = "2"
		m.MataUangLama[0].Payment.EDMPremiMenjadi = uang.Money{Amount: duit(t, "600").Amount, Currency: mataUangUji}
	})
	samaU(t, "EDMOldPayment (OLD EDM)", h.MataUang[0].Payment.EDMOldPayment, "600")
}

// TestDTPerubahanMataUangBaru - Check=="0": EDMPremiMenjadi = EDMNewPremi,
// EDMOldPayment/Premi/Commision = 0; Param lain terbawa (A20).
func TestDTPerubahanMataUangBaru(t *testing.T) {
	h := premiDT(t, func(m *services.MasukanPremiEDMDT) {
		m.MataUang = append(m.MataUang, bayar(t, "UJI-B", "500", "50", "0", "0", "0", "0"))
	})
	b := h.MataUang[1]
	samaU(t, "EDMPremiMenjadi", b.Payment.EDMPremiMenjadi, "450")
	samaU(t, "EDMOldPayment", b.Payment.EDMOldPayment, "0")
	samaU(t, "SumTotalPayment", b.SumTotalPayment, "450")
	samaU(t, "EDMOldPPN terbawa", b.Payment.EDMOldPPN, "16")
	if len(h.BarisTerbawa) != 1 {
		t.Errorf("BarisTerbawa %v", h.BarisTerbawa)
	}
}

// TestDTWhenBMenimpaA - mulai sama tetapi akhir berbeda: WHEN A DAN WHEN B
// benar; B (tanggal polis LAMA + EdmDate lama) menimpa: startdate 4, edmdate
// 2+1 = 3 → porsi akhir 0.75, awal 0.25.
func TestDTWhenBMenimpaA(t *testing.T) {
	h := premiDT(t, func(m *services.MasukanPremiEDMDT) { m.EndDateTime = tanggal(t, "2026-01-09 05:00") })
	samaRasio(t, "ProrateEDMEnd", h.ProrateEDMEnd, uang.Ratio{Value: rasio(t, "0.75").Value, Scale: 20})
	h = premiDT(t, func(m *services.MasukanPremiEDMDT) {
		m.EndDateTime, m.EdmType = tanggal(t, "2026-01-09 05:00"), models.EdmBatalProrata
	})
	samaRasio(t, "ProrateEDMEnd batal prorata = awal", h.ProrateEDMEnd, uang.Ratio{Value: rasio(t, "0.25").Value, Scale: 20})
}

// TestDTFixRate - FIX RATE: awal 0, akhir 1.
func TestDTFixRate(t *testing.T) {
	h := premiDT(t, func(m *services.MasukanPremiEDMDT) {
		m.EndDateTime, m.FixRate = tanggal(t, "2026-01-09 05:00"), true
	})
	samaRasio(t, "ProrateEDMEnd", h.ProrateEDMEnd, uang.Ratio{Value: rasio(t, "1").Value, Scale: 0})
}

// TestDTBatalProrata - WHEN EdmType==2 (porsi awal 0.25):
// OldPremiProrate = 800 − 800 × @divide(10,100,4) − 40 + 16 + 4 = 700.
func TestDTBatalProrata(t *testing.T) {
	h := premiDT(t, func(m *services.MasukanPremiEDMDT) {
		m.EndDateTime, m.EdmType = tanggal(t, "2026-01-09 05:00"), models.EdmBatalProrata
	})
	b := h.MataUang[0]
	samaU(t, "EDMNewPremi", b.Payment.EDMNewPremi, "865")
	samaU(t, "EDMOldPremi", b.Payment.EDMOldPremi, "800")
	samaU(t, "EDMOldPayment", b.Payment.EDMOldPayment, "700")
	samaU(t, "EDMPremiMenjadi", b.Payment.EDMPremiMenjadi, "175") // 700 × 0.25
	samaU(t, "SumTotalPayment", b.SumTotalPayment, "-525")
	samaU(t, "EDMOldCommision", b.Payment.EDMOldCommision, "80")
}

// TestDTBatalSejakSemula - WHEN EdmType==1: EDMNewPremi 0, EDMPremiMenjadi =
// netto baru, ProrateEDMEnd 1; ⛔ dihitung untuk lini apa pun (K-046).
func TestDTBatalSejakSemula(t *testing.T) {
	h := premiDT(t, func(m *services.MasukanPremiEDMDT) { m.EdmType = models.EdmBatalSejakSemula })
	samaRasio(t, "ProrateEDMEnd", h.ProrateEDMEnd, uang.Ratio{Value: rasio(t, "1").Value, Scale: 0})
	b := h.MataUang[0]
	samaU(t, "EDMNewPremi", b.Payment.EDMNewPremi, "0")
	samaU(t, "EDMPremiMenjadi", b.Payment.EDMPremiMenjadi, "865")
	samaU(t, "EDMOldPayment", b.Payment.EDMOldPayment, "700")
	samaU(t, "SumTotalPayment", b.SumTotalPayment, "165")
}

// TestDTDivideTidakEksak - `@divide` (mode belum terverifikasi): hasil tak
// eksak pada 4 desimal → galat.
func TestDTDivideTidakEksak(t *testing.T) {
	for nama, ubah := range map[string]func(*services.MasukanPremiEDMDT){
		"RICommision 12.345 / 100": func(m *services.MasukanPremiEDMDT) {
			m.EdmType, m.RICommision = models.EdmBatalProrata, rasio(t, "12.345")
		},
		"Premium 800.12345 / 1 (EdmType 1)": func(m *services.MasukanPremiEDMDT) {
			m.EdmType = models.EdmBatalSejakSemula
			m.MataUangLama[0].Payment.Premium = uang.Money{Amount: duit(t, "800.12345").Amount, Currency: mataUangUji}
		},
	} {
		if _, err := services.HitungPremiEDMDT(masukanDT(t, ubah)); !errors.Is(err, services.ErrModeDivideBelumTerverifikasi) {
			t.Errorf("%s: %v", nama, err)
		}
	}
}

// TestDTTanggalTakBulat - tanggal lama tidak dinormalisasi langkah 4 (mis.
// 17:00 GMT) melawan tanggal baru 05:00 → selisih n,5 hari → galat A09.
func TestDTTanggalTakBulat(t *testing.T) {
	_, err := services.HitungPremiEDMDT(masukanDT(t, func(m *services.MasukanPremiEDMDT) {
		m.LamaStart = tanggal(t, "2025-12-31 17:00")
	}))
	if !errors.Is(err, services.ErrSatuanSelisihWaktuBelumTerverifikasi) {
		t.Errorf("galat %v", err)
	}
}

// TestDTBatalMataUangBaru - cabang batal MENOLKAN Param OldPremi /
// OldPremiProrate / oldcommision tiap baris: baris tanpa pasangan lama
// bernilai lama 0 (bukan terbawa), sedangkan OldPPN dkk. tetap terbawa.
func TestDTBatalMataUangBaru(t *testing.T) {
	h := premiDT(t, func(m *services.MasukanPremiEDMDT) {
		m.EndDateTime, m.EdmType = tanggal(t, "2026-01-09 05:00"), models.EdmBatalProrata
		m.MataUang = append(m.MataUang, bayar(t, "UJI-B", "500", "50", "0", "0", "0", "0"))
	})
	b := h.MataUang[1]
	samaU(t, "EDMOldPremi", b.Payment.EDMOldPremi, "0")
	samaU(t, "EDMOldPayment", b.Payment.EDMOldPayment, "0")
	samaU(t, "EDMOldCommision", b.Payment.EDMOldCommision, "0")
	samaU(t, "EDMPremiMenjadi", b.Payment.EDMPremiMenjadi, "0")
	samaU(t, "EDMOldPPN terbawa", b.Payment.EDMOldPPN, "16")
}

// TestDTBarisPertamaTanpaPasangan - Param lama belum pernah diisi → galat
// (nilai awal Param belum terverifikasi), bukan nol tebakan.
func TestDTBarisPertamaTanpaPasangan(t *testing.T) {
	_, err := services.HitungPremiEDMDT(masukanDT(t, func(m *services.MasukanPremiEDMDT) { m.MataUangLama = nil }))
	if !errors.Is(err, services.ErrLokalBelumTerisi) {
		t.Errorf("galat %v", err)
	}
}
