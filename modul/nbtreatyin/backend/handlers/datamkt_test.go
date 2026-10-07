package handlers_test

// Uji seam 1 - `Activity/CheckDataMkt` langkah 4 lalu 5 (Obj-Save): SEMUA medan marketing yang ditulis
// langkah 4 tersimpan dan terbaca kembali (`[keputusan work owner 06-10-2026]` "CheckDataMkt ikuti aja itu
// semua, tambah ke quotation nya"; TeamGroup dibaca aturan tim Sec Head). Fixture UJI- (tiruan MO).

import (
	"net/http"
	"testing"
)

func TestCheckDataMktMenyimpanSemuaMedanMarketing(t *testing.T) {
	u := baru(t)
	id := u.buat()
	h := halamanLengkap("")
	kode, isi := u.panggil("POST", "/kasus/"+id+"/hitung", admin, map[string]any{"urutan": []map[string]string{{"aksi": "CheckDataMkt"}}, "halaman": h})
	if kode != http.StatusOK {
		t.Fatalf("CheckDataMkt: %d %s", kode, isi)
	}
	harap := map[string]string{"MOID": "UJI-MO-1", "MarketingCode": "UJI-MKT", "MarketingName": "UJI-MO",
		"TeamGroup": "1", "BranchCode": "UJI-CAB", "BranchName": "UJI-CABANG"}
	s := u.g.Halaman[id] // yang TERSIMPAN (proyeksi katalog = kolom Oracle)
	for m, v := range harap {
		if got := s.Ambil("Quotation." + m); got != v {
			t.Errorf("tersimpan Quotation.%s = %q, harap %q", m, got, v)
		}
	}
	_, isi = u.panggil("GET", "/kasus/"+id, admin, nil)
	baca := u.layar(isi).Halaman
	for _, m := range []string{"MOID", "MarketingCode", "TeamGroup"} {
		if got := baca.Ambil("PolicyTreatyIn.QuotationData." + m); got != harap[m] {
			t.Errorf("dibaca kembali PolicyTreatyIn.QuotationData.%s = %q, harap %q", m, got, harap[m])
		}
	}
}
