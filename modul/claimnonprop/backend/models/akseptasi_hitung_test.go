package models

import "testing"

// halamanAkseptasiUji - satu akseptasi: dua baris Claim Acceptation (IDR, USD) dan satu Loss Allocation USD.
func halamanAkseptasiUji() *Halaman {
	h := HalamanBaru()
	h.Setel(CD+"ShareCeding", "40")
	h.SetelDaftar(DaftarAdjustment, []Baris{{"Type": "1"}})
	h.SetelDaftar(JalurAdj(1, AnakClaimAccept), []Baris{
		{"Currency": "IDR", "CurrencyID": IDMataUangIDR, "AltValue": "1", "Value": "1000", "TPL": "100", "AdjusterFee": "50",
			"Salvage": "30", "CNPOthersFee": "20", "CNPDeductible": "40"},
		{"Currency": "USD", "CurrencyID": "10001", "AltValue": "15000", "Value": "10", "TPL": "0"},
	})
	h.SetelDaftar(JalurAdj(1, AnakLossAlloc), []Baris{{"Currency": "USD", "ClaimPercentage": "50"}})
	return h
}

// TestHitungKlaimAkseptasiRumusCountClaimTNP - CountClaimTNP_Act 13.2.3 atas baris akseptasi.
func TestHitungKlaimAkseptasiRumusCountClaimTNP(t *testing.T) {
	h := halamanAkseptasiUji()
	if err := HitungKlaimAkseptasi(h, 1); err != nil {
		t.Fatal(err)
	}
	r := h.AmbilDaftar(JalurAdj(1, AnakClaimAccept))
	// nilai = 1000 - 40 + 100 + 20 + 50 - 30 = 1100; IDR = 1100 × 1; Cedant = 1100 × 40 % = 440; Adjust = 960 × 40 %.
	for p, harap := range map[string]string{"USD": "1100", "ClaimAmountCedant": "440", "ClaimAmountAdjust": "384"} {
		if !SamaAngka(r[0][p], harap) {
			t.Errorf("IDR %s = %q, harap %s", p, r[0][p], harap)
		}
	}
	// USD: nilai 10 -> IDR 150.000, Cedant 4; Loss USD 50 % -> Claim Amount 4 × 50 % = 2.
	if !SamaAngka(r[1]["USD"], "150000") || !SamaAngka(r[1]["ClaimAmountCedant"], "4") {
		t.Errorf("USD = %v", r[1])
	}
	if l := h.AmbilDaftar(JalurAdj(1, AnakLossAlloc))[0]; !SamaAngka(l["ClaimAmountAdjust"], "2") {
		t.Errorf("loss USD = %v", l)
	}
	idr, cedant, err := TotalKlaimAkseptasi(h, 1)
	if err != nil || !SamaAngka(idr, "151100") || !SamaAngka(cedant, "60440") { // 440 + 4 × 15.000
		t.Errorf("total = %s / %s (%v)", idr, cedant, err)
	}
}

// TestHitungKlaimAkseptasiBarisTerkunciDibiarkan - baris `CNPFlagOuts = 1` tidak dihitung ulang.
func TestHitungKlaimAkseptasiBarisTerkunciDibiarkan(t *testing.T) {
	h := halamanAkseptasiUji()
	h.AmbilDaftar(JalurAdj(1, AnakClaimAccept))[0]["CNPFlagOuts"] = "1"
	h.AmbilDaftar(JalurAdj(1, AnakClaimAccept))[0]["USD"] = "7"
	if err := HitungKlaimAkseptasi(h, 1); err != nil {
		t.Fatal(err)
	}
	if v := h.AmbilDaftar(JalurAdj(1, AnakClaimAccept))[0]["USD"]; v != "7" {
		t.Errorf("baris terkunci berubah: %q", v)
	}
}

// TestTambahLossAkseptasiSatuPerMataUang - Add Loss Allocation akseptasi = satu baris per mata uang Claim Acceptation.
func TestTambahLossAkseptasiSatuPerMataUang(t *testing.T) {
	h := halamanAkseptasiUji()
	h.TambahBaris(JalurAdj(1, AnakClaimAccept), Baris{"Currency": "IDR", "CurrencyID": IDMataUangIDR})
	if err := TambahLossAkseptasi(h, 1); err != nil {
		t.Fatal(err)
	}
	ls := h.AmbilDaftar(JalurAdj(1, AnakLossAlloc))
	if len(ls) != 3 || ls[1]["Currency"] != "IDR" || ls[2]["Currency"] != "USD" {
		t.Fatalf("loss = %v", ls)
	}
	if err := TambahLossAkseptasi(h, 9); err == nil {
		t.Error("akseptasi tak ada mestinya galat")
	}
}
