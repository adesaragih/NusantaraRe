package handlers_test

// Popup Historical Survey Report (keputusan work owner 06-10-2026, membatalkan K7): daftar survei yang
// disunting admin tersimpan lewat Save (T_POLIS_SURVEY) dan kembali saat kasus dibuka.

import (
	"encoding/json"
	"net/http"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/services"
)

func TestSaveMenyimpanDaftarSurveiLewatHTTP(t *testing.T) {
	u := baru(t)
	id := u.buat()
	h := halamanLengkap("1")
	h.Setel("PolicyTreatyIn.QuotationData.IsSurveyReport", "Yes")
	h.SetelDaftar(models.DaftarSurvei, []models.Baris{{"DateofSurvey": "2026-09-01", "SurveyedBy": "UJI-SURVEYOR", "Remarks": "UJI-R"}})
	if kode, isi := u.panggil("PUT", "/kasus/"+id, admin, map[string]any{"halaman": h}); kode != http.StatusOK {
		t.Fatalf("save: %d %s", kode, isi)
	}
	_, isi := u.panggil("GET", "/kasus/"+id, admin, nil)
	var ly services.Layar
	if err := json.Unmarshal([]byte(isi), &ly); err != nil {
		t.Fatal(err)
	}
	b := ly.Halaman.AmbilDaftar(models.DaftarSurvei)
	if len(b) != 1 || b[0]["SurveyedBy"] != "UJI-SURVEYOR" || b[0]["DateofSurvey"] != "2026-09-01" || b[0]["Remarks"] != "UJI-R" {
		t.Fatalf("dibuka ulang daftar survei = %+v", b)
	}
}
