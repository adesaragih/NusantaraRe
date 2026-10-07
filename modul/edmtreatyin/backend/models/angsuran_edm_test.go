package models

// CountSpreading_Act VERSI EDM (spec-penyimpanan ID-39, AC 35-36): bawaan SplitRNMSharePct 100 pada baris tunggal
// (langkah 4), presisi 20 - BUKAN 100/jumlah baris presisi 10 (NB). Cabang %Share kosong -> pesan baris (keputusan
// work owner 07-10-2026, rekomendasi b), bukan @divide(.SplitRNMSharePct, PolicyTreatyIn.RNMShare, 20) - pembaginya
// tidak diisi rule mana pun. Data UJI- dan angka karangan.

import (
	"testing"
)

// AC 36 langkah 4: satu baris, SplitRNMSharePct kosong / 0 -> 100. %Share kosong -> pesan, tanpa pembagian.
func TestCountSpreadingEDMBawaanSeratus(t *testing.T) {
	for _, awal := range []string{"", "0", "0.00"} {
		h := HalamanBaru()
		h.Setel(pt+"NetPremium", "1000")
		h.SetelDaftar(DaftarSpreading, []Baris{{"TreatyType": "UJI-T", "SplitRNMSharePct": awal}})
		if err := CountSpreading(h, 1); err != nil {
			t.Fatalf("SplitRNMSharePct %q: %v", awal, err)
		}
		b := h.AmbilDaftar(DaftarSpreading)[0]
		if b["SplitRNMSharePct"] != "100" {
			t.Errorf("SplitRNMSharePct %q -> %q, harap 100", awal, b["SplitRNMSharePct"])
		}
		if b["SharePercentage"] != "" || len(h.SemuaPesan()) != 1 {
			t.Errorf("SplitRNMSharePct %q: %%Share %q, pesan %q", awal, b["SharePercentage"], h.SemuaPesan())
		}
	}
}

// Bawaan 100 hanya untuk baris TUNGGAL bernilai kosong/0: nilai lain dan daftar dua baris tidak disentuh.
func TestCountSpreadingEDMBawaanHanyaBarisTunggal(t *testing.T) {
	h := HalamanBaru()
	h.SetelDaftar(DaftarSpreading, []Baris{{"SplitRNMSharePct": "25", "SharePercentage": "30"}})
	if err := CountSpreading(h, 1); err != nil {
		t.Fatal(err)
	}
	if v := h.AmbilDaftar(DaftarSpreading)[0]["SplitRNMSharePct"]; v != "25" {
		t.Errorf("SplitRNMSharePct 25 -> %q", v)
	}
	h = HalamanBaru()
	h.SetelDaftar(DaftarSpreading, []Baris{{"SplitRNMSharePct": "", "SharePercentage": "60"}, {"SplitRNMSharePct": "20", "SharePercentage": "40"}})
	if err := CountSpreading(h, 1); err != nil {
		t.Fatal(err)
	}
	d := h.AmbilDaftar(DaftarSpreading)
	if d[0]["SplitRNMSharePct"] != "" || d[0]["SharePercentage"] != "60" || d[1]["SharePercentage"] != "40" {
		t.Errorf("dua baris disentuh: %v", d)
	}
}

// AC 35: presisi 20 pada `@divide(.SharePercentage,100,20)`. %Share 0.33333333333333333333 -> PremiumSpreaded =
// 300 x 0.00333333333333333333 = 0.99999999999999999900 (NB presisi 10: 0.9999999999).
func TestCountSpreadingEDMPresisiDuaPuluh(t *testing.T) {
	h := HalamanBaru()
	h.Setel(pt+"NetPremium", "300")
	h.SetelDaftar(DaftarSpreading, []Baris{{"SharePercentage": "0.33333333333333333333"}})
	if err := CountSpreading(h, 1); err != nil {
		t.Fatal(err)
	}
	b := h.AmbilDaftar(DaftarSpreading)[0]
	xeuSama(t, "PremiumSpreaded", b["PremiumSpreaded"], "0.99999999999999999900")
	if b["PremiumSpreaded"] == "0.9999999999" {
		t.Error("presisi 10 = versi NB (AC 35)")
	}
}

// Keputusan work owner 07-10-2026 (rekomendasi b): baris ber-%Share kosong - baris yang DITAMBAH lewat Add - tidak
// dibagi `PolicyTreatyIn.RNMShare` (tak diisi rule mana pun) melainkan ditolak dengan pesan di baris itu; baris lain
// tetap dihitung, total hanya dari baris ber-%Share.
func TestCountSpreadingEDMShareKosongDitolakDenganPesan(t *testing.T) {
	h := HalamanBaru()
	h.Setel(pt+"NetPremium", "1000")
	h.SetelDaftar(DaftarSpreading, []Baris{{"SharePercentage": "60"}, {"SplitRNMSharePct": "100"}})
	if err := CountSpreading(h, 1); err != nil {
		t.Fatalf("galat %v, harap pesan baris (bukan galat bagi nol)", err)
	}
	if p := h.SemuaPesan(); len(p) != 1 || p[0] != PesanShareSpreadingKosong(2) {
		t.Fatalf("pesan = %q", p)
	}
	d := h.AmbilDaftar(DaftarSpreading)
	if d[1]["SharePercentage"] != "" || d[1]["PremiumSpreaded"] != "" {
		t.Fatalf("baris 2 tanpa %%Share tidak boleh diisi: %+v", d[1])
	}
	xeuBaris(t, "baris 1", d[0], map[string]string{"PremiumSpreaded": "600"})
	xeuSama(t, "TotalSharePercentagePremium", h.Ambil(pt+"TotalSharePercentagePremium"), "60")
}

func TestCountSpreadingEDMBarisBerShare(t *testing.T) {
	h := HalamanBaru()
	h.Setel(pt+"NetPremium", "1000")
	h.SetelDaftar(DaftarSpreading, []Baris{{"SharePercentage": "40", "ClaimPercentage": "10"}})
	h.Setel(pt+"Claim", "500")
	if err := CountSpreading(h, 1); err != nil {
		t.Fatal(err)
	}
	xeuBaris(t, "baris", h.AmbilDaftar(DaftarSpreading)[0], map[string]string{"PremiumSpreaded": "400", "ClaimSpreaded": "50"})
	xeuSama(t, "TotalClaim", h.Ambil(pt+"TotalClaim"), "50")
}

// Langkah 6 berjalan juga tanpa baris: empat total = 0.
func TestCountSpreadingEDMTanpaBaris(t *testing.T) {
	h := HalamanBaru()
	if err := CountSpreading(h, 0); err != nil {
		t.Fatal(err)
	}
	for _, m := range []string{"TotalSharePercentagePremium", "TotalPremium", "TotalSharePercentageClaim", "TotalClaim"} {
		xeuSama(t, m, h.Ambil(pt+m), "0")
	}
}
