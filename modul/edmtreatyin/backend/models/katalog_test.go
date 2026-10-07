package models

// Uji murni ProyeksiKatalog - apa yang tersisa sesudah halaman disimpan ke
// delapan tabel diagram lalu dimuat kembali (tiket 21, spec-penyimpanan
// AC 49-50). Gudang tiruan memakainya supaya uji seam HTTP melihat penyimpanan
// yang sama dengan Oracle: medan tanpa kolom HILANG.

import "testing"

func TestProyeksiKatalogHanyaMedanBerkolom(t *testing.T) {
	h := HalamanBaru()
	h.Setel("PositionNote", PosisiAdmin)
	h.Setel(HalamanPolis+".PremiOgp", "1000")
	h.Setel(HalamanPolis+".CedingCoName", "UJI-A; UJI-B; ")        // R47: tingkat polis, apa adanya
	h.Setel(HalamanPolis+".BrokerageFee", "25")                    // hanya ditulis SetPPNPPH: tanpa kolom
	h.Setel(HalamanPolis+".IsOJKNopolis", "1")                     // hanya ditulis: tanpa kolom
	h.Setel(HalamanPolis+".LayerType", "UJI-L")                    // diagram F26: dicoret
	h.Setel(HalamanPolis+".PolicyNo", "UJI-NOPOL")                 // kunci NOPOLIS, milik SetelNomorPolis
	h.Setel(HalamanQuotation+".GroupPanel", "006")                 // diagram
	h.Setel(HalamanPolis+".QuotationData.GroupPanel", "")          // kosong: Quotation dipakai
	h.Setel(HalamanQuotation+".MOID", "UJI-MO-LAMA")               // ditimpa QuotationData
	h.Setel(HalamanPolis+".QuotationData.MOID", "UJI-MO-1")        // layar menyunting QuotationData
	h.Setel(HalamanQuotation+".BusinessType", "FireStyle2")        // turunan: tanpa kolom
	h.Setel(HalamanQuotation+".SobLeader0", "UJI-LEADER")          // nol pembaca: tanpa kolom
	h.Setel(HalamanPolis+".QuotationData.NoOfferSlip", "UJI-SLIP") // RALAT: tampil di layar
	h.SetelDaftar(DaftarAngsuran, []Baris{{"InstallmentNo": "1", "PaymentDate": "2026-11-01", "PPN": "2"}})
	h.SetelDaftar(JalurAnak(DaftarAngsuran, 1, "InstallmentList"), []Baris{{"PaymentDate": "2026-11-02", "PPN": "1"}})
	h.SetelDaftar(DaftarUsulan, []Baris{{"Suggest": "UJI"}}) // milik HISTORYAKSEPTASIPRODUCTION

	s := ProyeksiKatalog(h)
	for j, harap := range map[string]string{
		"PositionNote":                             PosisiAdmin,
		HalamanPolis + ".PremiOgp":                 "1000",
		HalamanPolis + ".CedingCoName":             "UJI-A; UJI-B; ",
		HalamanPolis + ".BrokerageFee":             "",
		HalamanPolis + ".IsOJKNopolis":             "",
		HalamanPolis + ".LayerType":                "",
		HalamanPolis + ".PolicyNo":                 "",
		HalamanQuotation + ".GroupPanel":           "006",
		HalamanPolis + ".QuotationData.GroupPanel": "006",
		HalamanQuotation + ".MOID":                 "UJI-MO-1",
		HalamanPolis + ".QuotationData.MOID":       "UJI-MO-1",
		HalamanQuotation + ".BusinessType":         "",
		HalamanQuotation + ".SobLeader0":           "",
		HalamanQuotation + ".NoOfferSlip":          "UJI-SLIP",
	} {
		if got := s.Ambil(j); got != harap {
			t.Errorf("%s = %q, harap %q", j, got, harap)
		}
	}
	a := s.AmbilDaftar(DaftarAngsuran)
	if len(a) != 1 || a[0]["InstallmentNo"] != "1" || a[0]["PPN"] != "2" {
		t.Errorf("angsuran %+v", a)
	}
	if _, ada := a[0]["PaymentDate"]; ada {
		t.Error("PAYMENT_DATE hanya di anak (diagram R61)")
	}
	r := s.AmbilDaftar(JalurAnak(DaftarAngsuran, 1, "InstallmentList"))
	if len(r) != 1 || r[0]["PaymentDate"] != "2026-11-02" {
		t.Errorf("rincian %+v", r)
	}
	if _, ada := r[0]["PPN"]; ada {
		t.Error("PPN rincian hanya ditulis: tanpa kolom")
	}
	if len(s.AmbilDaftar(DaftarUsulan)) != 0 {
		t.Error("SuggestList tidak disimpan di tabel polis (K4)")
	}
	if h.Ambil(HalamanPolis+".BrokerageFee") != "25" {
		t.Fatal("halaman asal tidak boleh berubah")
	}
}
