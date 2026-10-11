package models

// Uji spreading retro Type QR (keputusan work owner 06-10-2026). TANPA Oracle; nol nama orang (UJI-).

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"
)

func dsp(s string) *apd.Decimal {
	v, _, err := apd.NewFromString(s)
	if err != nil {
		panic(err)
	}
	return v
}

func teks(ds []*apd.Decimal) []string {
	var s []string
	for _, x := range ds {
		s = append(s, TeksSimpan(x))
	}
	return s
}

func TestKaskadeRetro(t *testing.T) {
	maks := []*apd.Decimal{dsp("3000000000"), dsp("6000000000"), dsp("20000000000"), dsp("25000000000")}
	for nama, k := range map[string]struct {
		nilai string
		mau   []string
		lebih bool
	}{
		"contoh 7,5 M":       {"7500000000", []string{"3000000000.0000", "3000000000.0000", "1500000000.0000", "0.0000"}, false},
		"di lapisan 1":       {"1000000000", []string{"1000000000.0000", "0.0000", "0.0000", "0.0000"}, false},
		"tepat di Max 2":     {"6000000000", []string{"3000000000.0000", "3000000000.0000", "0.0000", "0.0000"}, false},
		"tepat Max akhir":    {"25000000000", []string{"3000000000.0000", "3000000000.0000", "14000000000.0000", "5000000000.0000"}, false},
		"melebihi Max akhir": {"25000000000.5", nil, true},
	} {
		b, lebih := KaskadeRetro(dsp(k.nilai), maks)
		if lebih != k.lebih {
			t.Errorf("%s: lebih = %v", nama, lebih)
			continue
		}
		if !lebih && strings.Join(teks(b), "|") != strings.Join(k.mau, "|") {
			t.Errorf("%s: %v, mau %v", nama, teks(b), k.mau)
		}
	}
	// USD dengan celah 1 (231.000 → 231.001): celah tidak berpengaruh.
	b, lebih := KaskadeRetro(dsp("300000"), []*apd.Decimal{dsp("231000"), dsp("500000")})
	if lebih || strings.Join(teks(b), "|") != "231000.0000|69000.0000" {
		t.Errorf("USD celah 1: %v %v", teks(b), lebih)
	}
}

func TestHitungReinsurerRetro(t *testing.T) {
	// Bagian 3.000.000.000, share 50%, rate 2.5, EM 0.5, discount 20%, ovr 10%.
	r := HitungReinsurerRetro("UJI-RE", dsp("50"), dsp("20"), dsp("10"), dsp("3000000000"), dsp("2.5"), dsp("0.5"), true)
	// Amount 1.500.000.000; GROSS = 2.5/1000 × 1.5 × 1.5e9 = 5.625.000
	// Discount = 1.125.000; Comm = (5.625.000 − 1.125.000) × 10% = 450.000; NET = 4.050.000
	for k, v := range map[string]*apd.Decimal{"Amount": r.Amount, "Gross": r.Gross, "Discount": r.Discount, "Comm": r.Comm, "Net": r.Net} {
		mau := map[string]string{"Amount": "1500000000.0000", "Gross": "5625000.0000", "Discount": "1125000.0000", "Comm": "450000.0000", "Net": "4050000.0000"}[k]
		if TeksSimpan(v) != mau {
			t.Errorf("tahun pertama %s = %s, mau %s", k, TeksSimpan(v), mau)
		}
	}
	// Bukan tahun pertama: tanpa Discount; Comm = GROSS × 10% = 562.500; NET = 5.062.500. EM kosong → 0.
	r = HitungReinsurerRetro("UJI-RE", dsp("50"), dsp("20"), dsp("10"), dsp("3000000000"), dsp("2.5"), dsp("0.5"), false)
	if TeksSimpan(r.Discount) != "0.0000" || TeksSimpan(r.Comm) != "562500.0000" || TeksSimpan(r.Net) != "5062500.0000" {
		t.Errorf("bukan tahun pertama: disc %s comm %s net %s", TeksSimpan(r.Discount), TeksSimpan(r.Comm), TeksSimpan(r.Net))
	}
	if r = HitungReinsurerRetro("UJI-RE", dsp("100"), dsp("0"), dsp("0"), dsp("1000"), dsp("2"), nil, false); TeksSimpan(r.Gross) != "2.0000" {
		t.Errorf("EM kosong: gross %s", TeksSimpan(r.Gross))
	}
}

func TestPilihTahunRetro(t *testing.T) {
	tahun := []TahunRetro{
		{ID: "T1", Mulai: "2018-01-01", Akhir: "2024-12-31"},
		{ID: "T2", Mulai: "2025-01-01", Akhir: "2036-12-31"},
	}
	if got, err := PilihTahunRetro(tahun, time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)); err != nil || got.ID != "T2" {
		t.Errorf("di dalam periode: %v %v", got, err)
	}
	if got, err := PilihTahunRetro(tahun, time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)); err != nil || got.ID != "T1" {
		t.Errorf("tepat akhir periode: %v %v", got, err)
	}
	_, err := PilihTahunRetro(tahun, time.Date(2040, 1, 1, 0, 0, 0, 0, time.UTC))
	if !errors.Is(err, ErrSyaratSpreading) || !strings.Contains(err.Error(), "No Contract Retro Life covers begin date 01/01/2040") {
		t.Errorf("di luar semua periode: %v", err)
	}
	ganda := append(tahun, TahunRetro{ID: "T3", Mulai: "2026-01-01", Akhir: "2026-12-31"})
	_, err = PilihTahunRetro(ganda, time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC))
	if !errors.Is(err, ErrSyaratSpreading) || !strings.Contains(err.Error(), "T2, T3") {
		t.Errorf("periode ganda: %v", err)
	}
}

// bahanSpreading - dua lapisan IDR 0-3M (QS) dan 3M-6M (SURPLUS); business UJI-GTL; rate U age 40 contract 1 = 2.5.
func bahanSpreading() (TahunRetro, []KontrakRetro, []ReinsurerRetro, []BusinessRetro, map[string][]BarisRateQR) {
	tahun := TahunRetro{ID: "T2", Mulai: "2025-01-01", Akhir: "2036-12-31"}
	kontrak := []KontrakRetro{
		{ID: "K2", IDTahun: "T2", ReinsTypeID: "R2", ReinsTypeName: "SURPLUS", IDR: "6000000000", USD: "500000", BIDR: "3000000001"},
		{ID: "K1", IDTahun: "T2", ReinsTypeID: "R1", ReinsTypeName: "QS", IDR: "3000000000", USD: "231000", BIDR: "0"},
	}
	reins := []ReinsurerRetro{
		{KontrakID: "K1", Nama: "UJI-RE-A", PctShare: "60", Komisi: "20", OvrComm: "10"},
		{KontrakID: "K1", Nama: "UJI-RE-B", PctShare: "40", Komisi: "20", OvrComm: "10"},
		{KontrakID: "K2", Nama: "UJI-RE-C", PctShare: "100", Komisi: "0", OvrComm: "5"},
	}
	biz := []BusinessRetro{{KontrakID: "K1", BizName: "UJI-GTL", RIRateID: "RR1", RIRate: "UJI RATE QS"},
		{KontrakID: "K2", BizName: " uji-gtl ", RIRateID: "RR2", RIRate: "UJI RATE SP"}}
	rates := map[string][]BarisRateQR{
		"RR1": {{ID: "1", Gender: "U", Age: "40", Contract: "1", Rate: "2,5"}},
		"RR2": {{ID: "2", Gender: "U", Age: "40", Contract: "1", Rate: "4"}},
	}
	return tahun, kontrak, reins, biz, rates
}

func barisSpreading(retro, mu string) []BarisUnggah {
	return []BarisUnggah{{Nomor: 3, Nilai: map[string]string{"CURRENCY": mu},
		QR: &HasilBarisQR{Retro: dsp(retro), Umur: 40, Kontrak: dsp("1"), Sex: "", EM: dsp("0"), TahunPertama: true}}}
}

func TestSpreadingQR(t *testing.T) {
	tahun, kontrak, reins, biz, rates := bahanSpreading()
	m, err := SiapkanSpreading(tahun, kontrak, reins, biz, rates, "UJI-GTL", "1")
	if err != nil {
		t.Fatal(err)
	}
	baris := barisSpreading("4500000000", "IDR")
	tolak, err := SpreadingQR(m, baris, nil)
	if err != nil || len(tolak) != 0 {
		t.Fatalf("%v %+v", err, tolak)
	}
	sp := baris[0].Spreading
	if len(sp) != 2 || sp[0].Kontrak.ID != "K1" || TeksSimpan(sp[0].Bagian) != "3000000000.0000" || TeksSimpan(sp[1].Bagian) != "1500000000.0000" {
		t.Fatalf("lapisan = %+v", sp)
	}
	if len(sp[0].Reinsurer) != 2 || TeksSimpan(sp[0].Reinsurer[0].Amount) != "1800000000.0000" || TeksSimpan(sp[0].Reinsurer[0].Rate) != "2.5000" {
		t.Errorf("reinsurer QS = %+v", sp[0].Reinsurer)
	}
	// Lapisan berbagian 0 tidak disimpan.
	baris = barisSpreading("1000000000", "IDR")
	if _, err := SpreadingQR(m, baris, nil); err != nil || len(baris[0].Spreading) != 1 {
		t.Errorf("lapisan 0 ikut: %+v %v", baris[0].Spreading, err)
	}
	// Melebihi Max terakhir → baris ditolak.
	baris = barisSpreading("6000000001", "IDR")
	tolak, _ = SpreadingQR(m, baris, nil)
	if len(tolak) != 1 || !strings.Contains(tolak[0].Pesan, "exceeds the maximum limit 6.000.000.000 of Contract Retro Life T2 (row 3)") {
		t.Errorf("melebihi: %+v", tolak)
	}
	// Mata uang lain → baris ditolak.
	tolak, _ = SpreadingQR(m, barisSpreading("10", "SGD"), nil)
	if len(tolak) != 1 || tolak[0].Kolom != "CURRENCY" {
		t.Errorf("mata uang lain: %+v", tolak)
	}
	// Rate retro tidak ada (umur 70) → baris ditolak.
	baris = barisSpreading("10", "IDR")
	baris[0].QR.Umur = 70
	tolak, _ = SpreadingQR(m, baris, nil)
	if len(tolak) != 1 || tolak[0].Pesan != "No retro rate found for age 70, contract 1 in R/I Rate UJI RATE QS (QS, Contract Retro Life T2)" {
		t.Errorf("rate retro hilang: %+v", tolak)
	}
	// Baris tanpa hasil QR / terlewati tidak diproses.
	baris = barisSpreading("10", "IDR")
	if tolak, _ := SpreadingQR(m, baris, map[int]bool{3: true}); len(tolak) != 0 || baris[0].Spreading != nil {
		t.Error("baris terlewati diproses")
	}
}

func TestSpreadingSyaratLapisanTerpakai(t *testing.T) {
	cek := func(nama string, ubah func([]KontrakRetro, []ReinsurerRetro, []BusinessRetro) ([]KontrakRetro, []ReinsurerRetro, []BusinessRetro), retro, pesan string) {
		t.Helper()
		tahun, kontrak, reins, biz, rates := bahanSpreading()
		kontrak, reins, biz = ubah(kontrak, reins, biz)
		m, err := SiapkanSpreading(tahun, kontrak, reins, biz, rates, "UJI-GTL", "1")
		if err == nil {
			_, err = SpreadingQR(m, barisSpreading(retro, "IDR"), nil)
		}
		if !errors.Is(err, ErrSyaratSpreading) || !strings.Contains(err.Error(), pesan) {
			t.Errorf("%s: %v, mau memuat %q", nama, err, pesan)
		}
	}
	cek("share ≠ 100", func(k []KontrakRetro, r []ReinsurerRetro, b []BusinessRetro) ([]KontrakRetro, []ReinsurerRetro, []BusinessRetro) {
		r[1].PctShare = "30"
		return k, r, b
	}, "10", "Reinsurer shares of QS (Contract Retro Life T2) total 90%, must be 100%")
	cek("business tidak ada", func(k []KontrakRetro, r []ReinsurerRetro, b []BusinessRetro) ([]KontrakRetro, []ReinsurerRetro, []BusinessRetro) {
		return k, r, b[1:]
	}, "10", "Business UJI-GTL is not in the Business List of QS (Contract Retro Life T2)")
	cek("Max USD kosong", func(k []KontrakRetro, r []ReinsurerRetro, b []BusinessRetro) ([]KontrakRetro, []ReinsurerRetro, []BusinessRetro) {
		k[0].IDR = ""
		return k, r, b
	}, "10", "Maximum Limit (IDR) is empty on SURPLUS (Contract Retro Life T2)")
	// Lapisan yang TIDAK terpakai boleh tidak sah: SURPLUS tanpa business, nilai hanya di QS.
	tahun, kontrak, reins, biz, rates := bahanSpreading()
	m, err := SiapkanSpreading(tahun, kontrak, reins, biz[:1], rates, "UJI-GTL", "1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := SpreadingQR(m, barisSpreading("1000", "IDR"), nil); err != nil {
		t.Errorf("lapisan tak terpakai menolak: %v", err)
	}
}
