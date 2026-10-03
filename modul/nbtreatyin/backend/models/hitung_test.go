package models

// Uji seam 3 - rantai perhitungan uang (tiket 07: AC 27, 28; tiket 13: AC 79).
//
// ⭐ Nilai harapan DIHITUNG TANGAN dari rumus di INVENTARIS-XML.md bab 6,
// bukan dengan memanggil kode yang diuji:
//
//	PremiOgp 1000, RiCommOgp 10 %, OveriddingCommOgp 5 %, PremiOnp kosong,
//	Deduction1 20, TypeTax "Inclusive", RNMShareP 50
//	  ResultOgp1        = 1000 × 10/100                       = 100
//	  ResultOgp2        = 5/100 × 1000                         = 50
//	  BrokerageFee      = 0,025 × 1000                         = 25
//	  Brokerage sebenar = @divide(20, 1,022, 8)                = 19,56947162
//	  PPHValue          = 19,56947162 × 0,02                   = 0,3913894324
//	  PPNValue          = 19,56947162 × 0,022                  = 0,43052837564
//	  NetPremium        = (1000-100) - 5×1000/100 - 20         = 830
//	  BalanceDueTo      = 830 + PPN + PPH                      = 830,82191780804
//	  BalanceBeforePPH  = 830 + PPN                            = 830,43052837564
//	  BalanceBeforeTax  = 830
//	  GrossPremium      = @divide(1000×100, 50, 4)             = 2000

import (
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"
)

func samaAngka(t *testing.T, h *Halaman, jalur, harap string) {
	t.Helper()
	got, err := h.Angka(jalur)
	if err != nil {
		t.Fatalf("%s: %v", jalur, err)
	}
	want, _, err := apd.NewFromString(harap)
	if err != nil {
		t.Fatalf("harapan %q: %v", harap, err)
	}
	if got.Cmp(want) != 0 {
		t.Errorf("%s = %s, harap %s", jalur, got.Text('f'), harap)
	}
}

func halamanContoh(typeTax string) *Halaman {
	h := HalamanBaru()
	h.Setel("PolicyTreatyIn.QuotationData.ProportionalType", JenisProporsional)
	h.Setel("PolicyTreatyIn.PremiOgp", "1000")
	h.Setel("PolicyTreatyIn.RiCommOgp", "10")
	h.Setel("PolicyTreatyIn.OveriddingCommOgp", "5")
	h.Setel("PolicyTreatyIn.Deduction1", "20")
	h.Setel("PolicyTreatyIn.TypeTax", typeTax)
	h.Setel("TreatyIn.RNMShareP", "50")
	h.Setel("PolicyTreatyIn.FlagPPH", "true") // SetPPNPPH langkah 4 berjalan
	return h
}

// SetPPNPPH langkah 4: berjalan bila FlagPPH "true" ATAU agen PKP (STS_PKP 1);
// selain itu PPH/PPN nol (langkah 2) dan BrokerageFee tidak disentuh.
func TestSetPPNPPHBersyaratFlagAtauPKP(t *testing.T) {
	for _, tt := range []struct {
		flag, pkp string
		jalan     bool
	}{
		{"true", "", true},
		{"", "1", true},
		{"false", "1", true},
		{"false", "0", false},
		{"", "", false},
		{"TRUE", "", false}, // teks persis
	} {
		h := halamanContoh("Inclusive")
		h.Setel("PolicyTreatyIn.FlagPPH", tt.flag)
		h.Setel(JalurStsPKP, tt.pkp)
		h.Setel("PolicyTreatyIn.BrokerageFee", "7")
		if err := SetPPNPPH(h); err != nil {
			t.Fatal(err)
		}
		if tt.jalan {
			samaAngka(t, h, "PolicyTreatyIn.BrokerageFee", "25")
			samaAngka(t, h, "PolicyTreatyIn.PPHValue", "0.3913894324")
		} else {
			samaAngka(t, h, "PolicyTreatyIn.BrokerageFee", "7")
			samaAngka(t, h, "PolicyTreatyIn.PPHValue", "0")
			samaAngka(t, h, "PolicyTreatyIn.PPNValue", "0")
		}
	}
}

func TestPajakBrokerageInclusiveDibagi1022(t *testing.T) { // AC 27
	got, err := HitungBrokerageSebenarnya("Inclusive", "20")
	if err != nil {
		t.Fatal(err)
	}
	if got != "19.56947162" {
		t.Fatalf("Inclusive: dapat %s, harap 19.56947162 (20 / 1,022 dibulatkan 8)", got)
	}
}

func TestPajakBrokerageSelainInclusiveApaAdanya(t *testing.T) { // AC 28
	for _, tt := range []string{"", "inclusive", "Inclusive ", " Inclusive", "Exclusive"} {
		got, err := HitungBrokerageSebenarnya(tt, "20")
		if err != nil {
			t.Fatal(err)
		}
		if got != "20" {
			t.Errorf("TypeTax %q: dapat %s, harap 20 (tidak dibagi 1,022)", tt, got)
		}
	}
}

func TestRantaiNetPremiDanBalance(t *testing.T) {
	h := halamanContoh("Inclusive")
	if err := CountResult1(h, "Pct"); err != nil {
		t.Fatal(err)
	}
	if err := CountResult2Ogp(h, "Pct"); err != nil {
		t.Fatal(err)
	}
	samaAngka(t, h, "PolicyTreatyIn.ResultOgp1", "100")
	samaAngka(t, h, "PolicyTreatyIn.ResultOgp2", "50")
	samaAngka(t, h, "PolicyTreatyIn.BrokerageFee", "25")
	samaAngka(t, h, "PolicyTreatyIn.BrokerageFeeSebenarnya", "19.56947162")
	samaAngka(t, h, "PolicyTreatyIn.PPHValue", "0.3913894324")
	samaAngka(t, h, "PolicyTreatyIn.PPNValue", "0.43052837564")
	samaAngka(t, h, "PolicyTreatyIn.NetPremium", "830")
	samaAngka(t, h, "PolicyTreatyIn.BalanceDueTo", "830.82191780804")
	samaAngka(t, h, "PolicyTreatyIn.BalanceBeforePPH", "830.43052837564")
	samaAngka(t, h, "PolicyTreatyIn.BalanceBeforeTax", "830")
	if got := h.Ambil("PolicyTreatyIn.DueTo"); got != "1" {
		t.Errorf("DueTo = %q, harap 1 (BalanceDueTo >= 0)", got)
	}
}

func TestTypeTaxHurufKecilMenggeserBalance(t *testing.T) { // AC 28, spec §10.2 butir 2
	h := halamanContoh("inclusive")
	if err := CountResult1(h, "Pct"); err != nil {
		t.Fatal(err)
	}
	if err := CountResult2Ogp(h, "Pct"); err != nil {
		t.Fatal(err)
	}
	// 20 tidak dibagi: PPH 0,4 dan PPN 0,44 -> 830,84
	samaAngka(t, h, "PolicyTreatyIn.BalanceDueTo", "830.84")
}

func TestGrossPremiumMemakaiRNMShareP(t *testing.T) {
	h := halamanContoh("Inclusive")
	if err := CountResult1(h, ""); err != nil {
		t.Fatal(err)
	}
	samaAngka(t, h, "PolicyTreatyIn.GrossPremium", "2000")
}

func TestRiCommOgpSeratusTidakDihitung(t *testing.T) {
	// langkah 5 menuntut RiCommOgp < 100: tepat 100 KELUAR tanpa menghitung,
	// dan tanpa pesan (pesan hanya untuk > 100).
	h := halamanContoh("Inclusive")
	h.Setel("PolicyTreatyIn.RiCommOgp", "100")
	h.Setel("PolicyTreatyIn.ResultOgp1", "7")
	if err := CountResult1(h, "Pct"); err != nil {
		t.Fatal(err)
	}
	samaAngka(t, h, "PolicyTreatyIn.ResultOgp1", "7")
	if h.AdaPesan() {
		t.Fatalf("tepat 100 tidak memasang pesan, dapat %v", h.SemuaPesan())
	}
}

func TestRiCommOgpLebihSeratusMemasangPesan(t *testing.T) {
	h := halamanContoh("Inclusive")
	h.Setel("PolicyTreatyIn.RiCommOgp", "150")
	if err := CountResult1(h, "Pct"); err != nil {
		t.Fatal(err)
	}
	if got := h.Pesan["PolicyTreatyIn.RiCommOgp"]; len(got) != 1 || got[0] != PesanErrorMsg1 {
		t.Fatalf("pesan RiCommOgp = %v, harap [ErrorMsg1]", got)
	}
}

func TestPremiKosongMenolkanLimaMedan(t *testing.T) {
	h := HalamanBaru()
	h.Setel("PolicyTreatyIn.QuotationData.ProportionalType", JenisProporsional)
	h.Setel("PolicyTreatyIn.PremiOgp", "")
	h.Setel("PolicyTreatyIn.ResultOgp1", "99")
	h.Setel("TreatyIn.RNMShareP", "50")
	if err := CountOGPONP(h); err != nil {
		t.Fatal(err)
	}
	samaAngka(t, h, "PolicyTreatyIn.ResultOgp1", "0")
	samaAngka(t, h, "PolicyTreatyIn.NetPremium", "0")
	// langkah 8 tanpa klaim/salvage: ClaimType dan ClaimPaymentType dikosongkan
	if h.Ambil("PolicyTreatyIn.ClaimType") != "" || h.Ambil("PolicyTreatyIn.ClaimPaymentType") != "" {
		t.Fatal("tanpa klaim, ClaimType/ClaimPaymentType harus kosong")
	}
}

func TestKlaimMengisiClaimTypeSOA(t *testing.T) {
	h := halamanContoh("Inclusive")
	h.Setel("PolicyTreatyIn.Claim", "10")
	if err := CountOGPONP(h); err != nil {
		t.Fatal(err)
	}
	if h.Ambil("PolicyTreatyIn.ClaimType") != "SOA" || h.Ambil("PolicyTreatyIn.ClaimPaymentType") != "Claim" {
		t.Fatalf("ClaimType=%q ClaimPaymentType=%q, harap SOA / Claim",
			h.Ambil("PolicyTreatyIn.ClaimType"), h.Ambil("PolicyTreatyIn.ClaimPaymentType"))
	}
}

func TestBalanceNegatifDueToNol(t *testing.T) {
	h := HalamanBaru()
	h.Setel("PolicyTreatyIn.BalanceDueTo", "-0.01")
	if err := SetDueTo(h); err != nil {
		t.Fatal(err)
	}
	if h.Ambil("PolicyTreatyIn.DueTo") != "0" {
		t.Fatal("BalanceDueTo < 0 harus DueTo 0")
	}
}

func TestCalculatePremiDariGrossPremium(t *testing.T) {
	h := HalamanBaru()
	h.Setel("PolicyTreatyIn.QuotationData.ProportionalType", JenisProporsional)
	h.Setel("PolicyTreatyIn.GrossPremium", "2000")
	h.Setel("TreatyIn.RNMShareP", "50")
	h.Setel("TreatyIn.BrokeragePercentP", "12.5")
	if err := CalculatePremi(h, "PREMIUM"); err != nil {
		t.Fatal(err)
	}
	// PremiOgp = 2000×50/100 = 1000; ResultOgp1 = 1000×12,5/100 = 125; RiCommOgp = 125/1000×100 = 12,5
	samaAngka(t, h, "PolicyTreatyIn.PremiOgp", "1000")
	samaAngka(t, h, "PolicyTreatyIn.ResultOgp1", "125")
	samaAngka(t, h, "PolicyTreatyIn.RiCommOgp", "12.5")
}

func TestPembagiNolGagalTerang(t *testing.T) {
	h := halamanContoh("Inclusive")
	// Nilai master TERSEDIA tetapi nol (kosong = tidak tersedia, lihat
	// TestMasterTidakTersediaLangkahnyaDilewati).
	h.Setel("TreatyIn.RNMShareP", "0")
	if err := CountResult1(h, ""); err == nil {
		t.Fatal("RNMShareP nol: pembagian dengan nol harus galat, bukan nol diam-diam")
	}
}

func TestFillPaymentInstallmentTigaKali(t *testing.T) {
	h := HalamanBaru()
	h.Setel("PolicyTreatyIn.BalanceDueTo", "1000")
	h.Setel("PolicyTreatyIn.Installment", "3")
	sekarang := time.Date(2026, 10, 3, 9, 0, 0, 0, time.UTC)
	if err := FillPaymentInstallment(h, sekarang); err != nil {
		t.Fatal(err)
	}
	d := h.AmbilDaftar(DaftarAngsuran)
	if len(d) != 3 {
		t.Fatalf("dapat %d angsuran, harap 3", len(d))
	}
	// 100/3 dibulatkan 4 = 33,3333; baris terakhir 33,3333 + (100 - 99,9999) = 33,3334
	harap := []struct{ no, pct, premi string }{
		{"1", "33.3333", "333.3330"}, {"2", "33.3333", "333.3330"}, {"3", "33.3334", "333.3340"},
	}
	for i, w := range harap {
		if d[i]["InstallmentNo"] != w.no || d[i]["InstallmentPercentage"] != w.pct || d[i]["Premium"] != w.premi {
			t.Errorf("baris %d = %v, harap no %s pct %s premi %s", i+1, d[i], w.no, w.pct, w.premi)
		}
		if d[i]["DueDate"] != "2026-10-03" {
			t.Errorf("DueDate baris %d = %q", i+1, d[i]["DueDate"])
		}
	}
}

func TestAngsuranLebihSeratusMemasangPesan(t *testing.T) {
	h := HalamanBaru()
	h.Setel("PolicyTreatyIn.BalanceDueTo", "1000")
	h.SetelDaftar(DaftarAngsuran, []Baris{{"InstallmentPercentage": "60"}, {"InstallmentPercentage": "50"}})
	if err := SetValidateInstallment(h); err != nil {
		t.Fatal(err)
	}
	if got := h.Pesan[HalamanPolis]; len(got) != 1 || got[0] != PesanAngsuranLebih100 {
		t.Fatalf("pesan = %v", got)
	}
	// Premium dihitung ulang dari persentase: 1000 × 0,6 = 600
	if v, _, _ := apd.NewFromString(h.AmbilDaftar(DaftarAngsuran)[0]["Premium"]); v.Cmp(apd.New(600, 0)) != 0 {
		t.Fatalf("Premium baris 1 = %s", h.AmbilDaftar(DaftarAngsuran)[0]["Premium"])
	}
}

func TestSpreadingDuaBarisRata(t *testing.T) {
	h := HalamanBaru()
	h.Setel("PolicyTreatyIn.NetPremium", "830")
	h.Setel("PolicyTreatyIn.Claim", "10")
	h.SetelDaftar(DaftarSpreading, []Baris{{"TreatyType": "UJI-1"}, {"TreatyType": "UJI-2"}})
	if err := CountSpreading(h, 1); err != nil {
		t.Fatal(err)
	}
	d := h.AmbilDaftar(DaftarSpreading)
	for i, b := range d {
		// 100/2 = 50; premi 830 × 0,5 = 415; klaim (0+10-0) × 0,5 = 5
		for m, harap := range map[string]string{"SharePercentage": "50", "ClaimPercentage": "50", "PremiumSpreaded": "415", "ClaimSpreaded": "5"} {
			v, _, err := apd.NewFromString(b[m])
			w, _, _ := apd.NewFromString(harap)
			if err != nil || v.Cmp(w) != 0 {
				t.Errorf("baris %d %s = %q, harap %s", i+1, m, b[m], harap)
			}
		}
	}
	samaAngka(t, h, "PolicyTreatyIn.TotalSharePercentagePremium", "100")
	samaAngka(t, h, "PolicyTreatyIn.TotalPremium", "830")
	samaAngka(t, h, "PolicyTreatyIn.TotalClaim", "10")
}

// P29: halaman TreatyIn tidak lagi diisi JSON master. Tanpa nilai master,
// langkah yang membaginya DILEWATI (GrossClaim tetap), rantai lain berjalan.
func TestMasterTidakTersediaLangkahnyaDilewati(t *testing.T) {
	h := halamanContoh("Inclusive")
	h.Hapus("TreatyIn.RNMShareP")
	h.Setel("PolicyTreatyIn.Claim", "10")
	h.Setel("PolicyTreatyIn.GrossClaim", "77")
	if err := CountNetPremi(h); err != nil {
		t.Fatalf("rantai tidak boleh gagal tanpa master: %v", err)
	}
	samaAngka(t, h, "PolicyTreatyIn.GrossClaim", "77")
	samaAngka(t, h, "PolicyTreatyIn.NetPremium", "930") // (1000-0) - 5x1000/100 - 20; ResultOgp1 belum dihitung
	h.Setel("PolicyTreatyIn.GrossPremium", "2000")
	h.Setel("PolicyTreatyIn.PremiOgp", "1000")
	if err := CalculatePremi(h, "PREMIUM"); err != nil {
		t.Fatal(err)
	}
	samaAngka(t, h, "PolicyTreatyIn.PremiOgp", "1000")
}

// Tanpa master, rantai CountOGPONP utuh tetap berjalan (GrossPremium tidak
// dihitung, pesan validasi tetap terpasang).
func TestCountOGPONPTanpaMaster(t *testing.T) {
	h := halamanContoh("Inclusive")
	h.Hapus("TreatyIn.RNMShareP")
	h.Setel("PolicyTreatyIn.GrossPremium", "5")
	if err := CountOGPONP(h); err != nil {
		t.Fatalf("CountOGPONP tanpa master: %v", err)
	}
	samaAngka(t, h, "PolicyTreatyIn.GrossPremium", "5")
}
