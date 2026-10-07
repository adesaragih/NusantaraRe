package models

// Uji data polis layar Input Premium Detail - tiket 03 bagian 2. MURNI.

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"
)

func dataPolisLengkap() IsianDataPolis {
	return IsianDataPolis{
		Type: "QR", ProductNameID: "UJI-P1", ProductName: "UJI-PRODUK",
		SourceOfBusiness: "UJI-S1", SobName: "UJI-SOB",
		ProRateType: "1", MoID: "UJI-M1", MarketingCode: "UJI-MC", MarketingName: "UJI-MARKETING",
		AnnuityInterest: apd.New(5, -2), PremiumRefundFactor: apd.New(1, 0),
		DateReceived: tanggalUji(),
	}
}

func tanggalUji() *time.Time {
	t := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	return &t
}

func TestDataPolisWajibIsiDilaporkanSekaligus(t *testing.T) {
	_, err := SusunDataPolis(IsianDataPolis{Type: "TP"})
	if !errors.Is(err, ErrDataPolisBelumLengkap) {
		t.Fatalf("galat = %v", err)
	}
	for _, p := range []string{PesanProductNameKosong, PesanRISlipKosong, PesanProRateTypeKosong,
		PesanMarketingKosong, PesanAnnuityInterestKosong, PesanPremiumRefundKosong, PesanSOBKosong, PesanBillingKosong,
		PesanRetroKosong, PesanDateReceivedKosong} {
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

// SavePremiumList_Act 8.1/8.2 - penolakan Validate CSV (keputusan work owner
// 05-10-2026): pesan VERBATIM, "at list" = nomor baris CSV, SI dilewati TP/TR RNML-FL.
func TestPeriksaBatasProduk(t *testing.T) {
	b := BatasProduk{MinAge: apd.New(18, 0), MaxAge: apd.New(65, 0),
		MinSumInsured: apd.New(1000, 0), MaxSumInsured: apd.New(5000, 0)}
	baris := func(nomor int, nama, umur, si string) BarisUnggah {
		return BarisUnggah{Nomor: nomor, Nilai: map[string]string{"NAME_OF_INSURED": nama, "ENTRY_AGE": umur, "SUM_INSURED": si}}
	}
	peserta := []BarisUnggah{
		baris(3, "UJI-A", "30", "2000"),
		baris(7, "UJI-B", "70", "9000"),
		baris(9, "UJI-C", "", ""),
		baris(11, "UJI-D", "10", "500"),
		baris(12, "UJI-E", "x", "1,5"), // tidak terurai - urusan validasi bentuk
	}
	pesan := func(p []Penolakan) []string {
		var s []string
		for _, x := range p {
			s = append(s, fmt.Sprintf("%d|%s|%s", x.Baris, x.Kolom, x.Pesan))
		}
		return s
	}
	got := pesan(PeriksaBatasProduk("QR", "", b, peserta, nil))
	mau := []string{
		"7|ENTRY_AGE|UJI-B Age exceeds the limit, at list 7",
		"7|SUM_INSURED|UJI-B Sum Insured exceeds the limit, at list 7",
		"11|ENTRY_AGE|UJI-D Age exceeds the limit, at list 11",
		"11|SUM_INSURED|UJI-D Sum Insured exceeds the limit, at list 11",
	}
	if strings.Join(got, "\n") != strings.Join(mau, "\n") {
		t.Errorf("penolakan =\n%s\nmau\n%s", strings.Join(got, "\n"), strings.Join(mau, "\n"))
	}
	tp := PeriksaBatasProduk("TP", "UJI/RNML-FL/1", b, peserta, nil)
	for _, p := range tp {
		if p.Kolom != "ENTRY_AGE" {
			t.Errorf("TP RNML-FL tetap memeriksa Sum Insured: %+v", p)
		}
	}
	if len(tp) != 2 {
		t.Errorf("TP RNML-FL: %+v", tp)
	}
	if got := PeriksaBatasProduk("QR", "", BatasProduk{}, peserta, nil); len(got) != 0 {
		t.Errorf("batas kosong tetap menuduh: %+v", got)
	}
	// Baris yang sudah ditolak validasi tidak dicek lagi.
	if got := PeriksaBatasProduk("QR", "", b, peserta, map[int]bool{7: true, 11: true}); len(got) != 0 {
		t.Errorf("baris terlewati tetap dicek: %+v", got)
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

// Sebab penolakan batas berpemisah ribuan (permintaan work owner 05-10-2026).
func TestAngkaTampilSebabBatas(t *testing.T) {
	for masuk, mau := range map[string]string{"-27300184": "-27.300.184", "1000000000": "1.000.000.000",
		"0": "0", "70": "70", "1234.5": "1.234,5", "999": "999", "1000": "1.000"} {
		d, _, _ := apd.NewFromString(masuk)
		if got := AngkaTampil(d); got != mau {
			t.Errorf("AngkaTampil(%s) = %q, mau %q", masuk, got, mau)
		}
	}
	b := BatasProduk{MinSumInsured: apd.New(0, 0), MaxSumInsured: apd.New(1000000000, 0)}
	p := PeriksaBatasProduk("QR", "", b, []BarisUnggah{{Nomor: 8, Nilai: map[string]string{
		"NAME_OF_INSURED": "UJI-A", "ENTRY_AGE": "30", "SUM_INSURED": "-27300184"}}}, nil)
	if len(p) != 1 || p[0].Sebab != "SUM_INSURED -27.300.184 outside product limit 0 - 1.000.000.000" {
		t.Errorf("sebab = %+v", p)
	}
}
