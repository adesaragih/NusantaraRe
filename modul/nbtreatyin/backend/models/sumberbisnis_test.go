package models

// Uji pemilih Source Of Business (XOL Retro) - nilai harapan dari
// `DataTransform/SearchHierarkiSourceBizAgent_PostDT` dan `btnSOB_DT`.

import "testing"

func TestPilihAgenTanpaAnakMengisiSumberBisnis(t *testing.T) { // PostDT langkah 1.1, 1.2, 4, 5
	h := HalamanBaru()
	TombolSOB(h)
	b := BarisAgen{ID: "UJI-AG-1", ClientName: "UJI SUMBER SATU", Leader0: "UJI-AG-0", ChildCount: "0", ClientID: "UJI-K1"}
	if err := TerapkanSumberBisnis(h, b); err != nil {
		t.Fatal(err)
	}
	for j, harap := range map[string]string{
		"Quotation.btnQuotation":     "SOB",
		"Quotation.SourceOfBusiness": "UJI-AG-1",
		"Quotation.SobName":          "UJI SUMBER SATU",
		"Quotation.SobLeader0":       "UJI-AG-0",
		"Quotation.SobLeader1":       "", // .Leader1 bukan kolom RD
	} {
		if got := h.Ambil(j); got != harap {
			t.Errorf("%s = %q, harap %q", j, got, harap)
		}
	}
}

// `@if(.ChildCount > 0, "", ...)`: simpul yang punya anak MENGOSONGKAN keempat
// medan (bukan membiarkan nilai lama); ChildCount kosong = 0 di ekspresi Pega.
func TestPilihAgenBeranakMengosongkan(t *testing.T) {
	for _, tt := range []struct {
		anak        string
		harapSumber string
	}{
		{"3", ""},
		{"1", ""},
		{"0", "UJI-AG-2"},
		{"", "UJI-AG-2"},
	} {
		h := HalamanBaru()
		h.Setel("Quotation.SourceOfBusiness", "UJI-LAMA")
		h.Setel("Quotation.SobName", "UJI LAMA")
		h.Setel("Quotation.SobLeader0", "UJI-LAMA-0")
		TombolSOB(h)
		b := BarisAgen{ID: "UJI-AG-2", ClientName: "UJI SUMBER DUA", Leader0: "UJI-AG-0", ChildCount: tt.anak}
		if err := TerapkanSumberBisnis(h, b); err != nil {
			t.Fatalf("ChildCount %q: %v", tt.anak, err)
		}
		if got := h.Ambil("Quotation.SourceOfBusiness"); got != tt.harapSumber {
			t.Errorf("ChildCount %q: SourceOfBusiness %q, harap %q", tt.anak, got, tt.harapSumber)
		}
		if tt.harapSumber == "" {
			for _, m := range []string{"SobName", "SobLeader0", "SobLeader1"} {
				if got := h.Ambil("Quotation." + m); got != "" {
					t.Errorf("ChildCount %q: %s %q, harap kosong", tt.anak, m, got)
				}
			}
		}
	}
	h := HalamanBaru()
	TombolSOB(h)
	if err := TerapkanSumberBisnis(h, BarisAgen{ID: "UJI-AG-3", ChildCount: "UJI-BUKAN-ANGKA"}); err == nil {
		t.Fatal("ChildCount bukan angka: ekspresi Pega gagal, di sini galat")
	}
}

// Langkah 1 ber-WHEN `btnQuotation=="SOB"`; langkah 4-5 di luarnya, selalu jalan.
func TestPostDTTanpaTombolSOBHanyaLeader(t *testing.T) {
	h := HalamanBaru()
	h.Setel("Quotation.SourceOfBusiness", "UJI-LAMA")
	h.Setel("Quotation.SobName", "UJI LAMA")
	b := BarisAgen{ID: "UJI-AG-4", ClientName: "UJI SUMBER EMPAT", Leader0: "UJI-AG-0", ChildCount: "0"}
	if err := TerapkanSumberBisnis(h, b); err != nil {
		t.Fatal(err)
	}
	if h.Ambil("Quotation.SourceOfBusiness") != "UJI-LAMA" || h.Ambil("Quotation.SobName") != "UJI LAMA" {
		t.Fatalf("tanpa btnQuotation SOB langkah 1 dilewati: %q %q",
			h.Ambil("Quotation.SourceOfBusiness"), h.Ambil("Quotation.SobName"))
	}
	if h.Ambil("Quotation.SobLeader0") != "UJI-AG-0" {
		t.Fatalf("langkah 4 selalu jalan: %q", h.Ambil("Quotation.SobLeader0"))
	}
}
