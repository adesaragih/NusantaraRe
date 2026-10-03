package models

// Uji seam 3 - `Activity/TreatyInputPctCommSpreading` langkah 2-2.1.1.1, yang
// dipanggil `InputPolicyTreatyInDetail_preACT` langkah 17 (tiket 01, 13).
// Nilai harapan dibaca dari XML:
//
//	preACT 17  syarat `pyWorkPage.Quotation.ProportionalType=="NonProportional"`
//	           benar -> lewati langkah; salah -> jalankan
//	2.1        syarat `pyWorkPage.PolicyTreatyIn.TreatyType==.TreatyType`
//	2.1.1.1    syarat `.TreatyGroup==pyWorkPage.PolicyTreatyIn.TreatyGroupName`
//	           RiCommOgp = @replaceAll(.RIOGR,",",".")
//	           RiCommOgp = @replaceAll(.RIONR,",",".")   <- menimpa: hasil akhir RIONR
//
// Baris master diganti baris view TREATYINDETAILJOINEDM (kolom TREATYTYPE,
// TREATYGROUP, RIOGR, RIONR); fixture berawalan UJI-.

import "testing"

func halamanKomisi() *Halaman {
	h := HalamanBaru()
	h.Setel("PolicyTreatyIn.TreatyType", "UJI-QS")
	h.Setel("PolicyTreatyIn.TreatyGroupName", "UJI-GRUP-1")
	h.Setel("PolicyTreatyIn.RiCommOgp", "5")
	h.Setel("PolicyTreatyIn.RiCommOnp", "6")
	h.SetelDaftar(DaftarSpreading, []Baris{{"TreatyType": "UJI-SPR", "SharePercentage": "100"}})
	return h
}

func barisKomisi(id, jenis, grup, riogr, rionr string) BarisKontrak {
	return BarisKontrak{"ID": id, "TREATYTYPE": jenis, "TREATYGROUP": grup, "RIOGR": riogr, "RIONR": rionr}
}

func TestKomisiOgpDiambilDariRIONRBukanRIOGR(t *testing.T) {
	h := halamanKomisi()
	TreatyInputPctCommSpreading(h, []BarisKontrak{barisKomisi("UJI-D1", "UJI-QS", "UJI-GRUP-1", "30", "27,5")})
	// RIOGR ditulis lalu ditimpa RIONR di langkah yang sama; koma -> titik.
	if got := h.Ambil("PolicyTreatyIn.RiCommOgp"); got != "27.5" {
		t.Fatalf("RiCommOgp = %q, harap 27.5 (RIONR, koma jadi titik)", got)
	}
	// RiCommOnp tidak disentuh rule ini.
	if got := h.Ambil("PolicyTreatyIn.RiCommOnp"); got != "6" {
		t.Fatalf("RiCommOnp = %q, harap tetap 6", got)
	}
	// SpreadingRiskList(1) <- SpreadingTotalPct/SpreadingTypeID/SpreadingType:
	// bukan kolom view (alasan c) - baris spreading tidak diubah.
	d := h.AmbilDaftar(DaftarSpreading)
	if len(d) != 1 || d[0]["TreatyType"] != "UJI-SPR" || d[0]["SharePercentage"] != "100" {
		t.Fatalf("SpreadingRiskList berubah: %v", d)
	}
}

func TestKomisiOgpBersyaratJenisDanGrupTeksPersis(t *testing.T) {
	for _, b := range []BarisKontrak{
		barisKomisi("UJI-D1", "UJI-XOL", "UJI-GRUP-1", "30", "40"), // 2.1 salah
		barisKomisi("UJI-D2", "UJI-QS", "UJI-GRUP-2", "30", "40"),  // 2.1.1.1 salah
		barisKomisi("UJI-D3", "uji-qs", "UJI-GRUP-1", "30", "40"),  // pembanding teks persis
		barisKomisi("UJI-D4", "UJI-QS", "UJI-GRUP-1 ", "30", "40"), // spasi ekor
	} {
		h := halamanKomisi()
		TreatyInputPctCommSpreading(h, []BarisKontrak{b})
		if got := h.Ambil("PolicyTreatyIn.RiCommOgp"); got != "5" {
			t.Errorf("%s: RiCommOgp = %q, harap tetap 5 (syarat tidak terpenuhi)", b["ID"], got)
		}
	}
}

func TestKomisiOgpBarisCocokTerakhirMenang(t *testing.T) {
	h := halamanKomisi()
	TreatyInputPctCommSpreading(h, []BarisKontrak{
		barisKomisi("UJI-D1", "UJI-QS", "UJI-GRUP-1", "30", "20"),
		barisKomisi("UJI-D2", "UJI-QS", "UJI-GRUP-2", "30", "99"),
		barisKomisi("UJI-D3", "UJI-QS", "UJI-GRUP-1", "30", "22.5"),
	})
	// kalang Pega menimpa setiap kali syaratnya benar - yang tersisa baris cocok terakhir
	if got := h.Ambil("PolicyTreatyIn.RiCommOgp"); got != "22.5" {
		t.Fatalf("RiCommOgp = %q, harap 22.5", got)
	}
}

func TestKomisiOgpRIONRKosongDitulisKosong(t *testing.T) {
	h := halamanKomisi()
	TreatyInputPctCommSpreading(h, []BarisKontrak{barisKomisi("UJI-D1", "UJI-QS", "UJI-GRUP-1", "30", "")})
	if got := h.Ambil("PolicyTreatyIn.RiCommOgp"); got != "" {
		t.Fatalf("RiCommOgp = %q, harap kosong (Property-Set menulis nilai sumber apa adanya)", got)
	}
}

func TestKomisiOgpTanpaBarisTidakMengubah(t *testing.T) {
	h := halamanKomisi()
	TreatyInputPctCommSpreading(h, nil)
	if got := h.Ambil("PolicyTreatyIn.RiCommOgp"); got != "5" {
		t.Fatalf("RiCommOgp = %q, harap tetap 5", got)
	}
}

func TestLangkahKomisiHanyaBukanNonProporsional(t *testing.T) { // preACT langkah 17
	for _, tt := range []struct {
		jenis string
		jalan bool
	}{
		{JenisProporsional, true},
		{"", true},
		{"nonproportional", true}, // teks persis
		{JenisNonProporsional, false},
	} {
		h := HalamanBaru()
		h.Setel("Quotation.ProportionalType", tt.jenis)
		if got := LangkahKomisiProporsional(h); got != tt.jalan {
			t.Errorf("ProportionalType %q: %v, harap %v", tt.jenis, got, tt.jalan)
		}
	}
}
