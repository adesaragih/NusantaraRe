package handlers_test

// Uji seam 1 - Type Tax / With Tax berubah: pajak dihitung ulang saat itu juga
// (`[keputusan work owner 06-10-2026]` "rumus tax inclusive dan exclusive tidak
// berfungsi"; XML: sel `.TypeTax` hanya postValue). Nilai harapan dihitung tangan
// dari rumus XML; fixture UJI-.
//
//	Proporsional   SetPPNPPH 4: Inclusive 51,1 / 1,022 = 50; Exclusive 51,1 apa adanya
//	               PPH = x 2% ; PPN = x 2,2%
//	NonProp        preACT 16.2.2: Exclusive 306,6 -> PPN 6,7452, PPH 6,132,
//	               BalanceDueTo 2700 + 6,7452 + 6,132 = 2712,8772
//	               preACT 18: layer 102,2 -> PPN 2,2484 -> angsuran IDR .PPN

import (
	"net/http"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func hitungPajak(u *uji, id string, h *models.Halaman) *models.Halaman {
	u.t.Helper()
	kode, isi := u.panggil("POST", "/kasus/"+id+"/hitung", admin, map[string]any{"urutan": []map[string]string{{"aksi": "HitungPajak"}}, "halaman": h})
	if kode != http.StatusOK {
		u.t.Fatalf("hitung pajak: %d %s", kode, isi)
	}
	return u.layar(isi).Halaman
}

func TestTypeTaxBerubahMenghitungUlangPajakProporsional(t *testing.T) {
	for _, tt := range []struct{ typeTax, pph, ppn string }{
		{"Exclusive", "1.022", "1.1242"},
		{models.TypeTaxInclusive, "1", "1.1"},
	} {
		u := baru(t)
		id := u.buat()
		h := halamanLengkap("")
		h.Setel("PolicyTreatyIn.PremiOnp", "600")
		h.Setel("PolicyTreatyIn.Deduction1", "51.1")
		h.Setel("PolicyTreatyIn.FlagPPH", "true")
		h.Setel("PolicyTreatyIn.TypeTax", tt.typeTax)
		got := hitungPajak(u, id, h)
		angkaSama(t, got, "PolicyTreatyIn.PPHValue", tt.pph)
		angkaSama(t, got, "PolicyTreatyIn.PPNValue", tt.ppn)
	}
}

func TestTypeTaxBerubahMenghitungUlangPajakNonProp(t *testing.T) {
	u := baru(t)
	kontrakNP(u)
	id := u.buat()
	awal := models.HalamanBaru()
	awal.Setel("PolicyTreatyIn.FlagPPH", "true")
	awal.Setel("PolicyTreatyIn.TypeTax", models.TypeTaxInclusive)
	if kode, isi := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D-NP", "halaman": awal}); kode != http.StatusOK {
		t.Fatalf("pilih bisnis NonProp: %d %s", kode, isi)
	}
	_, isi := u.panggil("GET", "/kasus/"+id, admin, nil)
	h := u.layar(isi).Halaman
	h.Setel("PolicyTreatyIn.TypeTax", "Exclusive")
	// isian pengguna yang ikut ditulis ulang preACT 16 (13 tanggal, 23 spreading) tetap
	h.Setel("PolicyTreatyIn.StartDate", "2026-10-15")
	sp := h.AmbilDaftar(models.DaftarSpreading)
	sp[0]["SharePercentage"] = "50"
	h.SetelDaftar(models.DaftarSpreading, sp)

	got := hitungPajak(u, id, h)
	for j, harap := range map[string]string{
		"PolicyTreatyIn.PPNValue": "6.7452", "PolicyTreatyIn.PPHValue": "6.132", "PolicyTreatyIn.BalanceDueTo": "2712.8772",
	} {
		angkaSamaTeks(t, j, got.Ambil(j), harap)
	}
	if ang := got.AmbilDaftar(models.DaftarAngsuran); len(ang) != 2 {
		t.Fatalf("ListInstallment %d baris", len(ang))
	} else {
		angkaSamaTeks(t, "Angsuran(IDR).PPN", ang[0]["PPN"], "2.2484")
	}
	if got.Ambil("PolicyTreatyIn.StartDate") != "2026-10-15" {
		t.Errorf("StartDate isian pengguna tertimpa: %q", got.Ambil("PolicyTreatyIn.StartDate"))
	}
	if s := got.AmbilDaftar(models.DaftarSpreading); len(s) != 1 || s[0]["SharePercentage"] != "50" {
		t.Errorf("%%Share spreading isian pengguna tertimpa: %v", s)
	}
}
