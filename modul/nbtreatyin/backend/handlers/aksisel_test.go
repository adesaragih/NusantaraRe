package handlers_test

// Uji audit silang putaran 3 (bab 3.2 butir 1) lewat seam HTTP: aksi refresh
// sel layar admin yang membaca tabel acuan, dan langkah bersyarat
// `GeneratePolicyNoTreaty_Act`. Harapan dari langkah XML yang dikutip; nilai
// tabel acuan dari `backend/tiruan` (fixture UJI-).

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/services"
)

func (u *uji) hitung(id string, badan map[string]any) services.Layar {
	u.t.Helper()
	kode, isi := u.panggil("POST", "/kasus/"+id+"/hitung", admin, badan)
	if kode != http.StatusOK {
		u.t.Fatalf("hitung: %d %s", kode, isi)
	}
	var ly services.Layar
	if err := json.Unmarshal([]byte(isi), &ly); err != nil {
		u.t.Fatal(err)
	}
	return ly
}

// Sel `.IDCurrency` `Section/DetailPolicyTreatyIn`: change -> refresh
// `SetCurrency_act` (CURR = .IDCurrency):
//
//	1  InputData.CARI1 = Param.CURR
//	2  RDB-List GetCurrency (BrowsePage Currency)
//	3  .Currency = Currency.pxResults(1).HASIL1
//
// Tiruan GetCurrency: "UJI-ID-USD" -> "USD". Refresh tidak menyimpan.
func TestSetCurrencyMengisiNamaMataUangDariID(t *testing.T) {
	u := baru(t)
	id := u.buat()
	h := halamanLengkap("")
	h.Setel("PolicyTreatyIn.IDCurrency", "UJI-ID-USD")
	h.Setel("PolicyTreatyIn.Currency", "UJI-LAMA")
	ly := u.hitung(id, map[string]any{"urutan": []map[string]string{{"aksi": "SetCurrency"}}, "halaman": h})
	if got := ly.Halaman.Ambil("PolicyTreatyIn.Currency"); got != "USD" {
		t.Fatalf("Currency = %q, harap USD (GetCurrency HASIL1 dari IDCurrency)", got)
	}
	if got := u.g.Halaman[id].Ambil("PolicyTreatyIn.Currency"); got == "USD" {
		t.Fatal("SetCurrency_act tanpa Obj-Save: refresh tidak menyimpan")
	}
}

// Sel `.QuotationData.MOID`: change -> refresh `CheckDataMkt`. Langkah 4 mengisi
// Quotation / QuotationData / MarketingOfficer dari Obj-Browse marketing
// officer (tiruan: ClientID "UJI-MKT", ClientName "UJI-MO", TeamGroup "1"), dan
// langkah 5 `Obj-Save pyWorkPage` MENYIMPAN halaman - satu-satunya refresh
// yang menyimpan. Yang bertahan di penyimpanan hanya medan berkolom
// (MOID, MarketingName, MarketingOfficer); MarketingCode/TeamGroup dibuang
// sadar (`models/katalog.go` TabelQuotation: pembacanya hanya CheckDataMkt
// langkah 6 berlabel `//`).
func TestCheckDataMktMengisiMODanMenyimpan(t *testing.T) {
	u := baru(t)
	id := u.buat()
	h := halamanLengkap("")
	h.Setel("PolicyTreatyIn.QuotationData.MOID", "UJI-MO-7")
	ly := u.hitung(id, map[string]any{"urutan": []map[string]string{{"aksi": "CheckDataMkt"}}, "halaman": h})
	layar := map[string]string{
		"PolicyTreatyIn.MarketingOfficer":            "UJI-MO",
		"PolicyTreatyIn.QuotationData.MOID":          "UJI-MO-7",
		"PolicyTreatyIn.QuotationData.MarketingCode": "UJI-MKT",
		"PolicyTreatyIn.QuotationData.MarketingName": "UJI-MO",
		"PolicyTreatyIn.QuotationData.TeamGroup":     "1",
	}
	for j, w := range layar {
		if got := ly.Halaman.Ambil(j); got != w {
			t.Errorf("layar %s = %q, harap %q", j, got, w)
		}
	}
	for _, j := range []string{"PolicyTreatyIn.MarketingOfficer", "PolicyTreatyIn.QuotationData.MOID",
		"PolicyTreatyIn.QuotationData.MarketingName"} {
		if got := u.g.Halaman[id].Ambil(j); got != layar[j] {
			t.Errorf("tersimpan (Obj-Save langkah 5) %s = %q, harap %q", j, got, layar[j])
		}
	}
}

// `GeneratePolicyNoTreaty_Act` langkah 11 (Call FetchTreatyGroupOldID bila
// `TreatyGroupOldID==""`) dan 13 (Call FetchTreatyGroupOJK bila
// `OJKBusinessID==""`); langkah 28: PolicyNo = CARI4 + ".T" + OJKBusinessID +
// "." + MM.YYYY + "." + urut. OJK yang sudah terisi (preACT langkah 8 saat
// pilih bisnis) TIDAK dibaca ulang; grup lama kosong diisi (tiruan "UJI-OLD").
func TestNomorPolisMemakaiOJKTerisiDanMengisiGrupLama(t *testing.T) {
	u := baru(t)
	id := u.buat()
	if kode, isi := u.kirim(id, admin, halamanLengkap("1")); kode != http.StatusOK {
		t.Fatalf("admin: %d %s", kode, isi)
	}
	if kode, isi := u.kirim(id, secHead, putusan("1")); kode != http.StatusOK {
		t.Fatalf("Sec Head: %d %s", kode, isi)
	}
	u.g.Halaman[id].Setel("PolicyTreatyIn.OJKBusinessID", "UJI-OJK-ADA")
	kode, isi := u.panggil("POST", "/kasus/"+id+"/nomor-polis", deptHead, map[string]any{"halaman": putusan("1")})
	if kode != http.StatusOK {
		t.Fatalf("nomor polis: %d %s", kode, isi)
	}
	var n services.NomorPolis
	if err := json.Unmarshal([]byte(isi), &n); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(n.PolicyNo, "UJI-QR.TUJI-OJK-ADA.") {
		t.Fatalf("PolicyNo %q: langkah 13 tidak boleh menimpa OJKBusinessID yang terisi", n.PolicyNo)
	}
	if got := u.g.Halaman[id].Ambil("PolicyTreatyIn.TreatyGroupOldID"); got != "UJI-OLD" {
		t.Fatalf("TreatyGroupOldID = %q, harap UJI-OLD (FetchTreatyGroupOldID langkah 6)", got)
	}
}
