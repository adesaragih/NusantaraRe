package models

// Uji data polis layar Input Premium Detail - tiket 03 bagian 2. MURNI.

import (
	"errors"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"
)

func dataPolisLengkap() IsianDataPolis {
	return IsianDataPolis{
		Type: "QR", ProductNameID: "UJI-P1", ProductName: "UJI-PRODUK",
		SourceOfBusiness: "UJI-S1", SobName: "UJI-SOB",
		ProRateType: "1", MoID: "UJI-M1", MarketingCode: "UJI-MC", MarketingName: "UJI-MARKETING",
		AnnuityInterest: apd.New(5, -2), PremiumRefundFactor: apd.New(1, 0),
	}
}

func TestDataPolisWajibIsiDilaporkanSekaligus(t *testing.T) {
	_, err := SusunDataPolis(IsianDataPolis{Type: "TP"})
	if !errors.Is(err, ErrDataPolisBelumLengkap) {
		t.Fatalf("galat = %v", err)
	}
	for _, p := range []string{PesanProductNameKosong, PesanRISlipKosong, PesanProRateTypeKosong,
		PesanMarketingKosong, PesanAnnuityInterestKosong, PesanPremiumRefundKosong, PesanSOBKosong, PesanBillingKosong,
		PesanRetroKosong} {
		if !strings.Contains(err.Error(), p) {
			t.Errorf("pesan %q tidak dilaporkan: %v", p, err)
		}
	}
	if _, err := SusunDataPolis(dataPolisLengkap()); err != nil {
		t.Errorf("isian lengkap QR ditolak: %v", err)
	}
}

// R/I SLIP dan Billing Name hanya wajib untuk TP/TR; R/I SLIP dibuang untuk QR/QP.
func TestRISlipDanBillingHanyaUntukTPTR(t *testing.T) {
	d := dataPolisLengkap()
	d.RISlipRNM = "UJI-RNML-Q-1"
	got, err := SusunDataPolis(d)
	if err != nil || got.RISlipRNM != "" {
		t.Errorf("QR: R/I SLIP %q, %v - mau dibuang", got.RISlipRNM, err)
	}
	d.Type = "TR"
	if _, err := SusunDataPolis(d); err == nil || !strings.Contains(err.Error(), PesanBillingKosong) {
		t.Errorf("TR tanpa Billing Name: %v", err)
	}
}

// setSecurityReinsurer_act: dua Billing Name mengisi Retrocessionaire otomatis.
func TestRetrocessionaireOtomatis(t *testing.T) {
	d := dataPolisLengkap()
	d.Type, d.RISlipRNM, d.RetroID, d.RetroName = "TP", "UJI-RNML-Q-1", "L0000137", "UJI-BILLING"
	got, err := SusunDataPolis(d)
	if err != nil || got.SecurityReinsurerID != "B0000020" {
		t.Errorf("L0000137 -> %q, %v - mau B0000020", got.SecurityReinsurerID, err)
	}
	if _, _, ada := SecurityReinsurerOtomatis("L0000999"); ada {
		t.Error("Billing Name lain ikut terisi otomatis")
	}
}

func TestDataPolisMenolakPilihanAsing(t *testing.T) {
	for nama, ubah := range map[string]func(*IsianDataPolis){
		"type":      func(d *IsianDataPolis) { d.Type = "QX" },
		"pro rate":  func(d *IsianDataPolis) { d.ProRateType = "9" },
		"marketing": func(d *IsianDataPolis) { d.MarketingCode = "" },
	} {
		d := dataPolisLengkap()
		ubah(&d)
		if _, err := SusunDataPolis(d); !errors.Is(err, ErrDataPolisTidakSah) {
			t.Errorf("%s: %v", nama, err)
		}
	}
}

// SavePremiumList_Act 8.1/8.2 - pesan VERBATIM, nomor urut 1.., SI dilewati TP/TR RNML-FL.
func TestPeriksaBatasProduk(t *testing.T) {
	b := BatasProduk{MinAge: apd.New(18, 0), MaxAge: apd.New(65, 0),
		MinSumInsured: apd.New(1000, 0), MaxSumInsured: apd.New(5000, 0)}
	peserta := []PesertaBatas{
		{NameOfInsured: "UJI-A", EntryAge: apd.New(30, 0), SumInsured: apd.New(2000, 0)},
		{NameOfInsured: "UJI-B", EntryAge: apd.New(70, 0), SumInsured: apd.New(9000, 0)},
		{NameOfInsured: "UJI-C", EntryAge: nil, SumInsured: nil},
	}
	got := PeriksaBatasProduk("QR", "", b, peserta)
	mau := []string{"UJI-B Age exceeds the limit, at list 2", "UJI-B Sum Insured exceeds the limit, at list 2"}
	if strings.Join(got, "|") != strings.Join(mau, "|") {
		t.Errorf("pesan = %q, mau %q", got, mau)
	}
	if got := PeriksaBatasProduk("TP", "UJI/RNML-FL/1", b, peserta); len(got) != 1 || !strings.Contains(got[0], "Age") {
		t.Errorf("TP RNML-FL melewati Sum Insured saja: %q", got)
	}
	if got := PeriksaBatasProduk("QR", "", BatasProduk{}, peserta); len(got) != 0 {
		t.Errorf("batas kosong tetap menuduh: %q", got)
	}
}

// Premium Payment Method: kode 1/2/3 (Calculate1_Act), teks keputusan work owner.
func TestPilihanPremiumPaymentMethod(t *testing.T) {
	mau := map[string]string{"1": "Single", "2": "Annually", "3": "Others"}
	if len(PilihanProRateType) != len(mau) {
		t.Fatalf("pilihan = %v", PilihanProRateType)
	}
	for _, p := range PilihanProRateType {
		if mau[p.Kode] != p.Nama {
			t.Errorf("kode %q = %q, mau %q", p.Kode, p.Nama, mau[p.Kode])
		}
	}
}

// Billing Name dan Retrocessionaire hanya untuk TP/TR; di Type lain dikosongkan.
func TestBillingRetroDikosongkanSelainTPTR(t *testing.T) {
	d := dataPolisLengkap() // Type QR
	d.RetroID, d.RetroName = "UJI-R1", "UJI-BILLING"
	d.SecurityReinsurerID, d.SecurityReinsurer = "UJI-S1", "UJI-RETRO"
	got, err := SusunDataPolis(d)
	if err != nil {
		t.Fatal(err)
	}
	if got.RetroID != "" || got.RetroName != "" || got.SecurityReinsurerID != "" || got.SecurityReinsurer != "" {
		t.Errorf("QR menyimpan Billing/Retro tersembunyi: %+v", got)
	}
}
