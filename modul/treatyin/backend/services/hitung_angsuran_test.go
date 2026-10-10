package services

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

// kiniUji - 7 Oktober 2026 di Jakarta.
var kiniUji = time.Date(2026, time.October, 7, 9, 0, 0, 0, time.FixedZone("WIB", 7*3600))

func netIDR(v string) []NilaiMataUang {
	return []NilaiMataUang{{Currency: "IDR", CurrencyID: "10026", Value: v}}
}

func hitungAngsuranUji(t *testing.T, m MasukanAngsuran) HasilAngsuran {
	t.Helper()
	h, err := HitungAngsuran(m, kiniUji)
	if err != nil {
		t.Fatalf("HitungAngsuran: %v", err)
	}
	return h
}

func kolomAngsuran(a Angsuran, f func(BarisAngsuran) string) []string {
	out := []string{}
	for _, r := range a.InstallmentList {
		out = append(out, f(r))
	}
	return out
}

// TestAngsuranGambarPega38 - gambar `38-nonprop-menu-installment.png`: empat
// angsuran atas Net Premium IDR 644.674.819,59 → 25,00% dan 161.168.704,90
// per baris, `% Total` 100, `Total` 644.674.819,59.
func TestAngsuranGambarPega38(t *testing.T) {
	h := hitungAngsuranUji(t, MasukanAngsuran{Aksi: AksiAngsuranNilai, InstallmentNo: "4", NetPremium: netIDR("644674819.59")})
	if len(h.Angsuran) != 1 || h.Angsuran[0].Currency != "IDR" {
		t.Fatalf("satu halaman IDR, dapat %+v", h.Angsuran)
	}
	a := h.Angsuran[0]
	if got := kolomAngsuran(a, func(r BarisAngsuran) string { return r.Installment }); !reflect.DeepEqual(got, []string{"1", "2", "3", "4"}) {
		t.Errorf("Installment = %v", got)
	}
	if got := kolomAngsuran(a, func(r BarisAngsuran) string { return r.InstallmentPct }); !reflect.DeepEqual(got, []string{"25", "25", "25", "25"}) {
		t.Errorf("InstallmentPct = %v", got)
	}
	// 0.25 × 644674819.59 = 161168704.8975 → tampil 161.168.704,90.
	if got := kolomAngsuran(a, func(r BarisAngsuran) string { return r.Amount }); !reflect.DeepEqual(got, []string{"161168704.8975", "161168704.8975", "161168704.8975", "161168704.8975"}) {
		t.Errorf("Amount = %v", got)
	}
	if got := kolomAngsuran(a, func(r BarisAngsuran) string { return r.DueDate }); !reflect.DeepEqual(got, []string{"20261007", "20261007", "20261007", "20261007"}) {
		t.Errorf("DueDate = %v (hari ini Jakarta, bentuk tersimpan)", got)
	}
	if a.AmountTotal != "644674819.59" || a.PctTotal != "100" {
		t.Errorf("total = %s / %s", a.AmountTotal, a.PctTotal)
	}
	if !reflect.DeepEqual(h.TotalInstallmentNP, []NilaiMataUang{{Currency: "IDR", Value: "644674819.59"}}) {
		t.Errorf("TotalInstallmentNP = %+v", h.TotalInstallmentNP)
	}
	if h.InstallmentNo != "4" || len(h.Pesan) != 0 {
		t.Errorf("InstallmentNo %q pesan %v", h.InstallmentNo, h.Pesan)
	}
}

// TestAngsuranSisaPersenKeBarisTerakhir - langkah 8.2.2: sisa pembulatan
// masuk baris TERAKHIR, naik maupun turun.
func TestAngsuranSisaPersenKeBarisTerakhir(t *testing.T) {
	kasus := []struct {
		n      string
		persen []string
		amount []string
	}{
		// 100/3 = 33.33; 3 × 33.33 = 99.99 → terakhir 33.34.
		{"3", []string{"33.33", "33.33", "33.34"}, []string{"333.3", "333.3", "333.4"}},
		// 100/7 = 14.29; 7 × 14.29 = 100.03 → terakhir 14.26.
		{"7", []string{"14.29", "14.29", "14.29", "14.29", "14.29", "14.29", "14.26"},
			[]string{"142.9", "142.9", "142.9", "142.9", "142.9", "142.9", "142.6"}},
	}
	for _, k := range kasus {
		h := hitungAngsuranUji(t, MasukanAngsuran{Aksi: AksiAngsuranNilai, InstallmentNo: k.n, NetPremium: netIDR("1000")})
		a := h.Angsuran[0]
		if got := kolomAngsuran(a, func(r BarisAngsuran) string { return r.InstallmentPct }); !reflect.DeepEqual(got, k.persen) {
			t.Errorf("N=%s persen %v, mau %v", k.n, got, k.persen)
		}
		// Amount = @divide(pct,100,4) × value.
		if got := kolomAngsuran(a, func(r BarisAngsuran) string { return r.Amount }); !reflect.DeepEqual(got, k.amount) {
			t.Errorf("N=%s amount %v, mau %v", k.n, got, k.amount)
		}
		if a.PctTotal != "100" || a.AmountTotal != "1000" {
			t.Errorf("N=%s total %s/%s", k.n, a.PctTotal, a.AmountTotal)
		}
	}
}

// TestAngsuranPesanMengosongkanTab - langkah 2 membuang Installment dan
// TotalInstallmentNP SEBELUM pemeriksaan langkah 4/5.
func TestAngsuranPesanMengosongkanTab(t *testing.T) {
	lama := []Angsuran{{Currency: "IDR", AmountTotal: "5", InstallmentList: []BarisAngsuran{{Installment: "1"}}}}
	for _, k := range []struct{ n, pesan string }{
		{"0", pesanAngsuranKurang},
		{"-2", pesanAngsuranKurang},
		{"", pesanAngsuranKurang},
		{"abc", pesanAngsuranKurang},
	} {
		h := hitungAngsuranUji(t, MasukanAngsuran{Aksi: AksiAngsuranNilai, InstallmentNo: k.n, Angsuran: lama, NetPremium: netIDR("10")})
		if !reflect.DeepEqual(h.Pesan, []string{k.pesan}) || len(h.Angsuran) != 0 || h.TotalInstallmentNP == nil || len(h.TotalInstallmentNP) != 0 {
			t.Errorf("N=%q: %+v", k.n, h)
		}
	}
	h := hitungAngsuranUji(t, MasukanAngsuran{Aksi: AksiAngsuranNilai, InstallmentNo: "2", Angsuran: lama})
	if !reflect.DeepEqual(h.Pesan, []string{pesanNetPremiKosong}) || len(h.Angsuran) != 0 {
		t.Errorf("tanpa Net Premium: %+v", h)
	}
}

// TestAngsuranPerMataUang - satu halaman per baris TotalShareNetNP, nilainya
// dari baris bermata uang sama.
func TestAngsuranPerMataUang(t *testing.T) {
	h := hitungAngsuranUji(t, MasukanAngsuran{Aksi: AksiAngsuranNilai, InstallmentNo: "2", NetPremium: []NilaiMataUang{
		{Currency: "IDR", Value: "200"}, {Currency: "USD", Value: "10.5"},
	}})
	if len(h.Angsuran) != 2 || h.Angsuran[0].Currency != "IDR" || h.Angsuran[1].Currency != "USD" {
		t.Fatalf("%+v", h.Angsuran)
	}
	if got := kolomAngsuran(h.Angsuran[1], func(r BarisAngsuran) string { return r.Amount }); !reflect.DeepEqual(got, []string{"5.25", "5.25"}) {
		t.Errorf("USD amount %v", got)
	}
	if got := kolomAngsuran(h.Angsuran[1], func(r BarisAngsuran) string { return r.Currency }); !reflect.DeepEqual(got, []string{"USD", "USD"}) {
		t.Errorf("USD currency baris %v", got)
	}
	if !reflect.DeepEqual(h.TotalInstallmentNP, []NilaiMataUang{{Currency: "IDR", Value: "200"}, {Currency: "USD", Value: "10.5"}}) {
		t.Errorf("TotalInstallmentNP %+v", h.TotalInstallmentNP)
	}
}

// TestAngsuranUpdateMemulihkanTanpaMenghitungUlangAmount - langkah 9 lalu 11:
// persen, tanggal, WPC, dan tanggal bayar dipulihkan dari grid sebelumnya;
// Amount TETAP hasil langkah 8; baris yang tidak punya padanan dipertahankan.
func TestAngsuranUpdateMemulihkanTanpaMenghitungUlangAmount(t *testing.T) {
	sebelum := []Angsuran{{Currency: "IDR", InstallmentList: []BarisAngsuran{
		{Installment: "1", DueDate: "20260118", WPC: "60", PaymentDate: "20260318", InstallmentPct: "60", Amount: "600"},
		{Installment: "2", DueDate: "20260418", WPC: "45", PaymentDate: "20260602", InstallmentPct: "40", Amount: "400"},
	}}}
	h := hitungAngsuranUji(t, MasukanAngsuran{Aksi: AksiAngsuranNilai, Status: "update", InstallmentNo: "3", Angsuran: sebelum, NetPremium: netIDR("1000")})
	a := h.Angsuran[0]
	if got := kolomAngsuran(a, func(r BarisAngsuran) string { return r.InstallmentPct }); !reflect.DeepEqual(got, []string{"60", "40", "33.34"}) {
		t.Errorf("persen %v", got)
	}
	if got := kolomAngsuran(a, func(r BarisAngsuran) string { return r.DueDate + "|" + r.WPC + "|" + r.PaymentDate }); !reflect.DeepEqual(got, []string{"20260118|60|20260318", "20260418|45|20260602", "20261007||"}) {
		t.Errorf("tanggal %v", got)
	}
	if got := kolomAngsuran(a, func(r BarisAngsuran) string { return r.Amount }); !reflect.DeepEqual(got, []string{"333.3", "333.3", "333.4"}) {
		t.Errorf("Amount harus hasil langkah 8: %v", got)
	}
	if a.PctTotal != "133.34" || a.AmountTotal != "1000" {
		t.Errorf("total %s / %s", a.PctTotal, a.AmountTotal)
	}

	// Isian Installment (tanpa status) TIDAK memulihkan.
	h = hitungAngsuranUji(t, MasukanAngsuran{Aksi: AksiAngsuranNilai, InstallmentNo: "2", Angsuran: sebelum, NetPremium: netIDR("1000")})
	if got := kolomAngsuran(h.Angsuran[0], func(r BarisAngsuran) string { return r.InstallmentPct }); !reflect.DeepEqual(got, []string{"50", "50"}) {
		t.Errorf("tanpa status: %v", got)
	}
}

// TestAngsuranAdjustmentMemulihkanDariOld - langkah 10 (EDMState 1/2/3)
// menimpa langkah 9.
func TestAngsuranAdjustmentMemulihkanDariOld(t *testing.T) {
	kini := []Angsuran{{Currency: "IDR", InstallmentList: []BarisAngsuran{{InstallmentPct: "70", DueDate: "20260101"}}}}
	old := []Angsuran{{Currency: "IDR", InstallmentList: []BarisAngsuran{{InstallmentPct: "80", DueDate: "20250101", WPC: "30"}}}}
	for _, k := range []struct{ edm, mau string }{{"", "70"}, {"0", "70"}, {"1", "80"}, {"3", "80"}} {
		h := hitungAngsuranUji(t, MasukanAngsuran{Aksi: AksiAngsuranNilai, Status: "update", InstallmentNo: "2", EDMState: k.edm,
			Angsuran: kini, AngsuranLama: old, NetPremium: netIDR("100")})
		if got := h.Angsuran[0].InstallmentList[0].InstallmentPct; got != k.mau {
			t.Errorf("EDMState %q: persen %s, mau %s", k.edm, got, k.mau)
		}
	}
}

// TestAngsuranTotalBaris - SetTotalInstallment: `editpercentage` menghitung
// ulang Amount (skala 20), tanpa status hanya menjumlahkan.
func TestAngsuranTotalBaris(t *testing.T) {
	ang := []Angsuran{
		{Currency: "USD", AmountTotal: "x"},
		{Currency: "IDR", InstallmentList: []BarisAngsuran{
			{Currency: "IDR", InstallmentPct: "33.333", Amount: "1"},
			{Currency: "IDR", InstallmentPct: "66.667", Amount: "2"},
		}},
	}
	h := hitungAngsuranUji(t, MasukanAngsuran{Aksi: AksiAngsuranTotalBaris, Status: "editpercentage", Indeks: 1, Angsuran: ang, NetPremium: netIDR("3000")})
	a := h.Angsuran[1]
	if got := kolomAngsuran(a, func(r BarisAngsuran) string { return r.Amount }); !reflect.DeepEqual(got, []string{"999.99", "2000.01"}) {
		t.Errorf("editpercentage amount %v", got)
	}
	if a.AmountTotal != "3000" || a.PctTotal != "100" {
		t.Errorf("total %s / %s", a.AmountTotal, a.PctTotal)
	}
	if h.Angsuran[0].AmountTotal != "x" || h.TotalInstallmentNP != nil {
		t.Errorf("halaman lain dan TotalInstallmentNP tidak disentuh: %+v", h)
	}
	if ang[1].AmountTotal != "" {
		t.Errorf("masukan berubah: %+v", ang[1])
	}

	h = hitungAngsuranUji(t, MasukanAngsuran{Aksi: AksiAngsuranTotalBaris, Indeks: 1, Angsuran: ang, NetPremium: netIDR("3000")})
	if got := kolomAngsuran(h.Angsuran[1], func(r BarisAngsuran) string { return r.Amount }); !reflect.DeepEqual(got, []string{"1", "2"}) {
		t.Errorf("tanpa status Amount tetap: %v", got)
	}
	if h.Angsuran[1].AmountTotal != "3" {
		t.Errorf("AmountTotal %s", h.Angsuran[1].AmountTotal)
	}
}

// TestAngsuranUpdateTotal - TreatyInNPSetTotal(installment) langkah 27-28.
func TestAngsuranUpdateTotal(t *testing.T) {
	h := hitungAngsuranUji(t, MasukanAngsuran{Aksi: AksiAngsuranTotal, Angsuran: []Angsuran{
		{Currency: "IDR", AmountTotal: "10"}, {Currency: "USD", AmountTotal: ""},
	}})
	if !reflect.DeepEqual(h.TotalInstallmentNP, []NilaiMataUang{{Currency: "IDR", Value: "10"}, {Currency: "USD", Value: ""}}) {
		t.Errorf("%+v", h.TotalInstallmentNP)
	}
}

// TestAngsuranDitolak - aksi tak dikenal dan pagar 1000 angsuran.
func TestAngsuranDitolak(t *testing.T) {
	for _, m := range []MasukanAngsuran{
		{Aksi: "hapus"},
		{Aksi: AksiAngsuranNilai, InstallmentNo: "1001", NetPremium: netIDR("1")},
	} {
		if _, err := HitungAngsuran(m, kiniUji); !errors.Is(err, ErrMasukanTidakSah) {
			t.Errorf("%+v: err %v", m, err)
		}
	}
	// Batas pagar sendiri masih dihitung.
	h := hitungAngsuranUji(t, MasukanAngsuran{Aksi: AksiAngsuranNilai, InstallmentNo: "1000", NetPremium: netIDR("1")})
	if len(h.Angsuran[0].InstallmentList) != 1000 {
		t.Errorf("1000 baris, dapat %d", len(h.Angsuran[0].InstallmentList))
	}
}

// TestAngsuranBukanBilanganBulat - REPEAT 1..N dengan N pecahan: putaran
// selama i <= N, dan 8.2.2 tidak pernah berlaku (n tidak pernah = N).
func TestAngsuranBukanBilanganBulat(t *testing.T) {
	h := hitungAngsuranUji(t, MasukanAngsuran{Aksi: AksiAngsuranNilai, InstallmentNo: "2.5", NetPremium: netIDR("100")})
	a := h.Angsuran[0]
	if got := kolomAngsuran(a, func(r BarisAngsuran) string { return r.InstallmentPct }); !reflect.DeepEqual(got, []string{"40", "40"}) {
		t.Errorf("%v", got)
	}
	if a.PctTotal != "80" {
		t.Errorf("PctTotal %s - disalin apa adanya", a.PctTotal)
	}
}

// ⭐ `TreatyInUpdatePaymentDate` (sel Due Date): Payment Date = Due Date +
// WPC + 1 hari, SETIAP baris halaman itu saja; halaman lain tidak disentuh.
func TestTanggalBayarSatuHalaman(t *testing.T) {
	m := MasukanAngsuran{Aksi: AksiAngsuranTanggalBayar, Indeks: 0, Angsuran: []Angsuran{
		{Currency: "IDR", InstallmentList: []BarisAngsuran{
			{DueDate: "20240118", WPC: "60"},
			{DueDate: "20240228", WPC: ""},
			{DueDate: "", WPC: "30", PaymentDate: "20990101"},
		}},
		{Currency: "USD", InstallmentList: []BarisAngsuran{{DueDate: "20240118", WPC: "1", PaymentDate: "lama"}}},
	}}
	h := hitungAngsuranUji(t, m)
	got := kolomAngsuran(h.Angsuran[0], func(b BarisAngsuran) string { return b.PaymentDate })
	// 18 Jan + 61 hari = 19 Mar 2024 (kabisat); 28 Feb + 1 = 29 Feb; Due Date kosong → dibiarkan.
	if want := []string{"20240319", "20240229", "20990101"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("PaymentDate %v, mau %v", got, want)
	}
	if h.Angsuran[1].InstallmentList[0].PaymentDate != "lama" {
		t.Fatalf("halaman lain ikut berubah: %+v", h.Angsuran[1])
	}
}

// ⭐ `TreatyInUpdatePaymentDate_Act` (sel WPC): SEMUA halaman, baris ber-Due Date.
func TestTanggalBayarSemuaHalaman(t *testing.T) {
	m := MasukanAngsuran{Aksi: AksiAngsuranTanggalBayarSemua, Angsuran: []Angsuran{
		{Currency: "IDR", InstallmentList: []BarisAngsuran{{DueDate: "20241231", WPC: "0"}}},
		{Currency: "USD", InstallmentList: []BarisAngsuran{{DueDate: "18-01-2024", WPC: "14"}, {DueDate: "", PaymentDate: "x"}}},
	}}
	h := hitungAngsuranUji(t, m)
	if p := h.Angsuran[0].InstallmentList[0].PaymentDate; p != "20250101" {
		t.Fatalf("IDR PaymentDate %q, mau 20250101", p)
	}
	if got := kolomAngsuran(h.Angsuran[1], func(b BarisAngsuran) string { return b.PaymentDate }); !reflect.DeepEqual(got, []string{"20240202", "x"}) {
		t.Fatalf("USD PaymentDate %v", got)
	}
}
