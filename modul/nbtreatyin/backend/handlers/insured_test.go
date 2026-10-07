package handlers_test

// Keputusan work owner 06-10-2026: Insured Name = CEDING kontrak terpilih (tombol Choose). Diuji LEWAT HTTP:
// tersimpan (Quotation untuk kolom portal "Insured Name", PolicyTreatyIn untuk layar), terkirim di jawaban
// layar, dan tetap ada saat kasus dibuka ulang.

import (
	"encoding/json"
	"net/http"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/services"
)

func TestPilihBisnisMengisiInsuredNameDariCedingLewatHTTP(t *testing.T) {
	u := baru(t)
	id := u.buat()
	u.g.Kontrak["UJI-D7"] = models.BarisKontrak{"ID": "UJI-D7", "TREATYID": "UJI-T7", "CEDING": "UJI-CEDING-7",
		"CEDINGID": "UJI-C7", "PROPORTIONTYPE": "Proportional"}
	kode, isi := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D7"})
	if kode != http.StatusOK {
		t.Fatalf("pilih bisnis: %d %s", kode, isi)
	}
	var ly services.Layar
	if err := json.Unmarshal([]byte(isi), &ly); err != nil {
		t.Fatal(err)
	}
	if got := ly.Halaman.Ambil("PolicyTreatyIn.InsuredName"); got != "UJI-CEDING-7" {
		t.Errorf("layar PolicyTreatyIn.InsuredName = %q, harap UJI-CEDING-7", got)
	}
	h := u.g.Halaman[id]
	for _, j := range []string{"Quotation.InsuredName", "PolicyTreatyIn.InsuredName"} {
		if got := h.Ambil(j); got != "UJI-CEDING-7" {
			t.Errorf("tersimpan %s = %q, harap UJI-CEDING-7", j, got)
		}
	}
	_, isi = u.panggil("GET", "/kasus/"+id, admin, nil)
	ly = services.Layar{}
	if err := json.Unmarshal([]byte(isi), &ly); err != nil {
		t.Fatal(err)
	}
	if got := ly.Halaman.Ambil("PolicyTreatyIn.InsuredName"); got != "UJI-CEDING-7" {
		t.Errorf("dibuka ulang PolicyTreatyIn.InsuredName = %q, harap UJI-CEDING-7", got)
	}
}
