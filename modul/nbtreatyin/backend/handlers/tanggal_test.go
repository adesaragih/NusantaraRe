package handlers_test

// `DataTransform/SystemSetOneYear_DT` - pra-transform aksi refresh `change`
// medan `.StartDate` di `Section/DetailPolicyTreatyIn` (layar admin; di
// `DetailDeptHeadTreatyIn_UW` medannya `pyReadOnly` sehingga tidak pernah
// berubah). Langkah 1: `.EndDate = @DateTime.addCalendar(.StartDate,1,0,0,0,0,0,0)`.

import (
	"encoding/json"
	"net/http"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/services"
)

func TestUbahTanggalMulaiMengisiTanggalAkhirSetahun(t *testing.T) {
	u := baru(t)
	id := u.buat()
	h := halamanLengkap("")
	h.Setel("PolicyTreatyIn.StartDate", "2026-10-15")
	h.Setel("PolicyTreatyIn.EndDate", "2026-10-03")
	kode, isi := u.panggil("POST", "/kasus/"+id+"/hitung", admin, map[string]any{"urutan": []map[string]string{{"aksi": "SystemSetOneYear"}}, "halaman": h})
	if kode != http.StatusOK {
		t.Fatalf("%d %s", kode, isi)
	}
	var ly services.Layar
	if err := json.Unmarshal([]byte(isi), &ly); err != nil {
		t.Fatal(err)
	}
	if got := ly.Halaman.Ambil("PolicyTreatyIn.EndDate"); got != "2027-10-15" {
		t.Fatalf("EndDate = StartDate + 1 tahun: mau 2027-10-15, dapat %q", got)
	}
	if u.g.Halaman[id].Ambil("PolicyTreatyIn.EndDate") == "2027-10-15" {
		t.Fatal("refresh tidak menyimpan")
	}
}
