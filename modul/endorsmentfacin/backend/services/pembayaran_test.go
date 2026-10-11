package services_test

// Selisih pembayaran per mata uang - tiket E13. Angka SINTETIS.
// Dibaca sesudah: services/pembayaran.go, porsipembayaran_test.go.
//
// Baris baru  : Premium 1000, Commision 100, BrokerageFee 50, Deduction2 10, PPh 5, PPN 20
//               → netto 865 (= 1000 − 100 − 50 − 10 + 5 + 20, K-046 pajak DITAMBAH)
// Baris lama  : Premium 800, Commision 80, BrokerageFee 40, Deduction2 8, PPh 4, PPN 16
//               → oldpayment 692
// EDMOldDeduction2 masukan baris baru: 8 (blok 14 dan TSIObj tidak menulisnya).

import (
	"errors"
	"testing"

	"nusantarare/inti/backend/uang"
	"nusantarare/modul/endorsmentfacin/backend/models"
	"nusantarare/modul/endorsmentfacin/backend/services"
)

func bayar(t *testing.T, name string, premi, comm, brok, ded2, pph, ppn string) services.BarisBayar {
	t.Helper()
	u := func(s string) uang.Money { return uang.Money{Amount: duit(t, s).Amount, Currency: name} }
	return services.BarisBayar{Name: name, Payment: services.PembayaranBaris{
		Premium: u(premi), Commision: u(comm), BrokerageFee: u(brok), Deduction2: u(ded2), PPh: u(pph), PPN: u(ppn),
		EDMPremiMenjadi: u("0"), EDMOldPremi: u("0"), EDMOldPayment: u("0"), EDMOldCommision: u("0"),
		EDMOldBrokerage: u("0"), EDMOldPPN: u("0"), EDMOldPPh: u("0"), EDMOldDeduction2: u("8"),
		ListInstallment: []services.Cicilan{{InstallmentPercentage: rasio(t, "100")}},
	}}
}

func baruUji(t *testing.T) services.BarisBayar {
	return bayar(t, mataUangUji, "1000", "100", "50", "10", "5", "20")
}

func lamaUji(t *testing.T) services.BarisBayar {
	return bayar(t, mataUangUji, "800", "80", "40", "8", "4", "16")
}

// bayarUji - porsiBayar dengan satu pasangan baris baru/lama.
func bayarUji(t *testing.T, ubah func(*services.MasukanPorsiPembayaran)) services.HasilPorsiPembayaran {
	t.Helper()
	return porsiBayar(t, func(m *services.MasukanPorsiPembayaran) {
		m.MataUang = []services.BarisBayar{baruUji(t)}
		m.MataUangLama = []services.BarisBayar{lamaUji(t)}
		ubah(m)
	})
}

func samaU(t *testing.T, label string, got uang.Money, want string) {
	t.Helper()
	w := duit(t, want)
	if got.Kosong() || got.Amount.Cmp(w.Amount) != 0 {
		t.Errorf("%s: %v, mau %s", label, got.Amount, want)
	}
}

// TestBlok13Pembayaran - 13.2, r = 1: EDMPremiMenjadi = netto; selisih tiap
// komponen terhadap medan EDMOld* yang baru ditulis; cicilan 100%.
func TestBlok13Pembayaran(t *testing.T) {
	// EDMOldDeduction2 masukan 99: 13.2.3 menimpanya dengan Deduction2 lama (8).
	h := bayarUji(t, func(m *services.MasukanPorsiPembayaran) {
		m.Predikat.IsEdmExtendPeriod = true
		m.MataUang[0].Payment.EDMOldDeduction2 = uang.Money{Amount: duit(t, "99").Amount, Currency: mataUangUji}
	})
	b := h.MataUang[0]
	samaU(t, "EDMNewPremi", b.Payment.EDMNewPremi, "865")
	samaU(t, "EDMPremiMenjadi", b.Payment.EDMPremiMenjadi, "865")
	samaU(t, "EDMOldPayment", b.Payment.EDMOldPayment, "692")
	samaU(t, "EDMOldDeduction2 (13.2.3 menulis)", b.Payment.EDMOldDeduction2, "8")
	samaU(t, "SumTotalPayment", b.SumTotalPayment, "173")
	c := b.Payment.ListInstallment[0]
	samaU(t, "cicilan Premium", c.Premium, "200")
	samaU(t, "cicilan RICommision", c.RICommision, "20")
	samaU(t, "cicilan Deduction2", c.Deduction2, "2")
	// 200 − 20 − 10 − 2 + 1 + 4
	samaU(t, "cicilan PaymentTotal", c.PaymentTotal, "173")
	if b.Payment.Premium.Amount.Cmp(duit(t, "1000").Amount) != 0 {
		t.Error("masukan termutasi")
	}
}

// TestK046_NettoPremi_PPhPPNDitambahkan - netto MENAMBAHKAN PPh dan PPN.
func TestK046_NettoPremi_PPhPPNDitambahkan(t *testing.T) {
	h := bayarUji(t, func(m *services.MasukanPorsiPembayaran) { m.Predikat.IsEdmExtendPeriod = true })
	// Bila pajak DIKURANGKAN: 1000 − 100 − 50 − 10 − 5 − 20 = 815.
	samaU(t, "EDMNewPremi", h.MataUang[0].Payment.EDMNewPremi, "865")
}

// TestBlok14Pembayaran - 14.4, rAkhir 3/4, rAwal 1/4: bagian baru + bagian
// lama; Deduction2ins tanpa bagian lama (K-046).
func TestBlok14Pembayaran(t *testing.T) {
	h := bayarUji(t, func(m *services.MasukanPorsiPembayaran) { m.Predikat.IsEdmAdjRate = true })
	b := h.MataUang[0]
	// 865 × 0.75 + 692 × 0.25
	samaU(t, "EDMPremiMenjadi", b.Payment.EDMPremiMenjadi, "821.75")
	samaU(t, "SumTotalPayment", b.SumTotalPayment, "129.75")
	samaU(t, "EDMOldDeduction2 (tidak ditulis 14.4.3)", b.Payment.EDMOldDeduction2, "8")
	c := b.Payment.ListInstallment[0]
	samaU(t, "cicilan Premium", c.Premium, "150")        // 750 + 200 − 800
	samaU(t, "cicilan Deduction2", c.Deduction2, "-0.5") // 10 × 0.75 − 8, tanpa bagian lama
	// 150 − 15 − 7.5 − (−0.5) + 0.75 + 3
	samaU(t, "cicilan PaymentTotal", c.PaymentTotal, "131.75")
	if h.IsEDMShare == nil || *h.IsEDMShare {
		t.Errorf("IsEDMShare %v, mau tertulis false (share sama)", h.IsEDMShare)
	}
}

// TestBlok14Share - 14.4.7: share baru ≠ lama → IsEDMShare 1.
func TestBlok14Share(t *testing.T) {
	h := bayarUji(t, func(m *services.MasukanPorsiPembayaran) {
		m.Predikat.IsEdmAdjRate, m.ShareBaru = true, rasio(t, "6")
	})
	if h.IsEDMShare == nil || !*h.IsEDMShare {
		t.Errorf("IsEDMShare %v", h.IsEDMShare)
	}
}

// TestBlok14MataUangBaruNilaiTerbawa - 14.4.4 menolkan EDMOldPremi/Payment/
// Commision dan menghitung ulang (preminew − Commision) × rAkhir; komponen lama
// lain TERBAWA dari baris sebelumnya - mata uang lain (A20, K-046).
func TestBlok14MataUangBaruNilaiTerbawa(t *testing.T) {
	h := bayarUji(t, func(m *services.MasukanPorsiPembayaran) {
		m.Predikat.IsEdmAdjRate = true
		m.MataUang = append(m.MataUang, bayar(t, "UJI-B", "500", "50", "0", "0", "0", "0"))
	})
	if len(h.BarisTerbawa) != 1 || h.BarisTerbawa[0] != "UJI-B" {
		t.Errorf("BarisTerbawa %v", h.BarisTerbawa)
	}
	b := h.MataUang[1]
	samaU(t, "EDMPremiMenjadi", b.Payment.EDMPremiMenjadi, "337.5") // (500 − 50) × 0.75
	samaU(t, "EDMOldPayment", b.Payment.EDMOldPayment, "0")
	samaU(t, "SumTotalPayment", b.SumTotalPayment, "337.5")
	samaU(t, "EDMOldBrokerage terbawa", b.Payment.EDMOldBrokerage, "40")
	if b.Payment.EDMOldBrokerage.Currency != "UJI-B" {
		t.Errorf("mata uang %q", b.Payment.EDMOldBrokerage.Currency)
	}
}

// TestLokalBelumTerisi - mata uang baru di baris PERTAMA: lokal lama belum
// pernah diisi → galat, bukan nol tebakan.
func TestLokalBelumTerisi(t *testing.T) {
	_, err := services.HitungPorsiPembayaran(masukanBayar(t, func(m *services.MasukanPorsiPembayaran) {
		m.Predikat.IsEdmExtendPeriod = true
		m.MataUangLama = nil
	}))
	if !errors.Is(err, services.ErrLokalBelumTerisi) {
		t.Errorf("galat %v", err)
	}
}

// TestNilaiKosongDanMataUangBaris - medan kosong / mata uang lain → galat.
func TestNilaiKosongDanMataUangBaris(t *testing.T) {
	_, err := services.HitungPorsiPembayaran(masukanBayar(t, func(m *services.MasukanPorsiPembayaran) {
		m.Predikat.IsEdmExtendPeriod = true
		m.MataUang[0].Payment.Commision = uang.Money{}
	}))
	if !errors.Is(err, services.ErrNilaiPembayaranKosong) {
		t.Errorf("kosong: %v", err)
	}
	_, err = services.HitungPorsiPembayaran(masukanBayar(t, func(m *services.MasukanPorsiPembayaran) {
		m.Predikat.IsEdmExtendPeriod = true
		m.MataUang[0].Payment.PPN = uang.Money{Amount: duit(t, "1").Amount, Currency: "LAIN"}
	}))
	if !errors.Is(err, services.ErrMataUangBaris) {
		t.Errorf("mata uang: %v", err)
	}
}

// TestBlok15Pembayaran - 15.5, r = 1; selisih terhadap LOKAL lama.
func TestBlok15Pembayaran(t *testing.T) {
	h := bayarUji(t, func(m *services.MasukanPorsiPembayaran) { m.Predikat.IsEdmAdjPeriod = true })
	b := h.MataUang[0]
	samaU(t, "EDMPremiMenjadi", b.Payment.EDMPremiMenjadi, "865")
	samaU(t, "SumTotalPayment", b.SumTotalPayment, "173")
	samaU(t, "cicilan Deduction2", b.Payment.ListInstallment[0].Deduction2, "2") // 10 − OldDeduction2 lokal 8
	samaU(t, "cicilan PaymentTotal", b.Payment.ListInstallment[0].PaymentTotal, "173")
}

// TestTSIObjBlok1Pembayaran - TSIObj 1.4 (r = (3+1)/4 = 1) lalu langkah 3
// (prorateend = 1): ⛔ K-046 PaymentTotal TSIObj TANPA −Deduction2.
func TestTSIObjBlok1Pembayaran(t *testing.T) {
	h := bayarUji(t, func(m *services.MasukanPorsiPembayaran) { m.Predikat.IsEdmAdjTSI = true })
	b := h.MataUang[0]
	samaU(t, "EDMPremiMenjadi", b.Payment.EDMPremiMenjadi, "865")
	samaU(t, "EDMOldDeduction2 (1.4.3 tidak menulis)", b.Payment.EDMOldDeduction2, "8")
	samaU(t, "NetPremium", b.Payment.NetPremium, "173")
	c := b.Payment.ListInstallment[0]
	samaU(t, "cicilan Deduction2", c.Deduction2, "2")
	// 200 − 20 − 10 + 1 + 4 (TANPA − 2)
	samaU(t, "cicilan PaymentTotal", c.PaymentTotal, "175")
	if h.FlagEDMCurrency {
		t.Error("FlagEDMCurrency padahal jumlah baris sama")
	}
}

// TestTSIObjBlok1MataUangBaru - 1.4.1 menolkan oldpremi/oldcomm tiap baris;
// 1.4.4 menolkan EDMOldPremi/Payment/Commision/Deduction2; 1.5 flag.
func TestTSIObjBlok1MataUangBaru(t *testing.T) {
	h := bayarUji(t, func(m *services.MasukanPorsiPembayaran) {
		m.Predikat.IsEdmAdjTSI = true
		m.MataUang = append(m.MataUang, bayar(t, "UJI-B", "500", "50", "0", "0", "0", "0"))
	})
	b := h.MataUang[1]
	samaU(t, "EDMOldPremi", b.Payment.EDMOldPremi, "0")
	samaU(t, "EDMOldDeduction2", b.Payment.EDMOldDeduction2, "0")
	samaU(t, "SumTotalPayment", b.SumTotalPayment, "450") // EDMPremiMenjadi 450 − 0
	if !h.FlagEDMCurrency {
		t.Error("FlagEDMCurrency tidak menyala (2 baris baru, 1 lama)")
	}
}

// TestTSIObjBlok2 - 2.2: oldpayment dari baris lama menurut CARI24.
func TestTSIObjBlok2(t *testing.T) {
	nb := bayarUji(t, func(m *services.MasukanPorsiPembayaran) { m.Predikat.IsEdmAdjInsured, m.HasilProdKe = true, "0" })
	samaU(t, "OLD NB (Premium − Commision lama)", nb.MataUang[0].Payment.EDMPremiMenjadi, "720")
	samaU(t, "SumTotalPayment", nb.MataUang[0].SumTotalPayment, "0")
	edm := bayarUji(t, func(m *services.MasukanPorsiPembayaran) {
		m.Predikat.IsEdmAdjRefNo, m.HasilProdKe = true, "3"
		m.MataUangLama[0].Payment.EDMPremiMenjadi = uang.Money{Amount: duit(t, "777").Amount, Currency: mataUangUji}
	})
	samaU(t, "OLD EDM (EDMPremiMenjadi lama)", edm.MataUang[0].Payment.EDMPremiMenjadi, "777")
	_, err := services.HitungPorsiPembayaran(masukanBayar(t, func(m *services.MasukanPorsiPembayaran) {
		m.Predikat.IsEdmAdjSpreading, m.HasilProdKe = true, ""
	}))
	if !errors.Is(err, services.ErrNilaiPembayaranKosong) {
		t.Errorf("HASIL2 kosong: %v", err)
	}
}

// TestTSIObjLangkah3SalinCicilan - 3.3: baris tanpa ListInstallment memakai
// salinan cicilan baris PERTAMA; langkah 3 berjalan untuk semua jenis.
func TestTSIObjLangkah3SalinCicilan(t *testing.T) {
	h := bayarUji(t, func(m *services.MasukanPorsiPembayaran) {
		kedua := lamaUji(t)
		kedua.Name = mataUangUji
		kedua.Payment.ListInstallment = nil
		m.MataUang = append(m.MataUang, kedua)
	})
	if n := len(h.MataUang[1].Payment.ListInstallment); n != 1 {
		t.Fatalf("cicilan baris 2: %d", n)
	}
	// prorateend awal 1: 800 − EDMOldPremi 0 = 800.
	samaU(t, "cicilan Premium baris 2", h.MataUang[1].Payment.ListInstallment[0].Premium, "800")
}

// masukanBayar - masukan bayarUji tanpa menjalankan (untuk uji galat).
func masukanBayar(t *testing.T, ubah func(*services.MasukanPorsiPembayaran)) services.MasukanPorsiPembayaran {
	t.Helper()
	m := services.MasukanPorsiPembayaran{
		Predikat:      services.PredikatPembayaran{IsEDM: true},
		EDMDay:        "365",
		StartDateTime: tanggal(t, "2026-01-01 00:00"), EndDateTime: tanggal(t, "2026-01-05 00:00"),
		EdmDate:  tanggal(t, "2026-01-02 05:00"),
		MataUang: []services.BarisBayar{baruUji(t)}, MataUangLama: []services.BarisBayar{lamaUji(t)},
		ProrateEDMEndAwal: rasio(t, "1"), ShareBaru: rasio(t, "5"), ShareLama: rasio(t, "5"),
		HasilProdKe: "0",
	}
	ubah(&m)
	return m
}

// TestRangkaianLangkahSepuluh - langkah 10 (DataTransform) dijalankan di
// tempatnya: TSIObj langkah 3 membaca keluaran DataTransform (NetPremium =
// 865 − 700 dari OldPremi DataTransform), lalu blok 13 menimpa medan EDM*
// (EDMOldPayment akhir 692 dari lokal activity). Urutan terlihat dari dua
// nilai itu.
func TestRangkaianLangkahSepuluh(t *testing.T) {
	h := bayarUji(t, func(m *services.MasukanPorsiPembayaran) {
		m.Predikat.IsEdmExtendPeriod = true
		m.EdmType = models.EdmPerubahan
		m.LangkahSepuluh = &services.LangkahSepuluh{
			LamaStart: tanggal(t, "2026-01-01 05:00"), LamaEnd: tanggal(t, "2026-01-05 05:00"),
			LamaEdmDate: tanggal(t, "2026-01-03 05:00"), RICommision: rasio(t, "10"),
		}
		m.ProrateEDMEndAwal = uang.Ratio{} // tidak dipakai bila langkah 10 jalan
	})
	b := h.MataUang[0]
	samaU(t, "NetPremium (TSIObj langkah 3 atas keluaran DataTransform)", b.Payment.NetPremium, "165")
	samaU(t, "EDMOldPayment (blok 13 menimpa)", b.Payment.EDMOldPayment, "692")
	samaU(t, "SumTotalPayment", b.SumTotalPayment, "173")
}
