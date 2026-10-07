package handlers_test

// Keputusan work owner 06-10-2026 (portal kolom "Group Business" = Quotation.BusinessName, RD GetListOpportunity):
// bila CLASSOFBUSINESS kontrak terpilih kosong, Group Business diambil dari tabel BUSINESS lewat ID grup -
// TREATYGROUP.COAID = BUSINESS.BUSINESSGROUPID - lalu BUSINESSGROUPNAME. Class of Business yang terisi tetap
// dipakai apa adanya (preACT langkah 3). Insured Name = CEDING kontrak terpilih.

import (
	"net/http"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func TestPilihBisnisTanpaClassOfBusinessIsiGrupBisnisDariTabelBusiness(t *testing.T) {
	u := baru(t)
	id := u.buat()
	u.g.Kontrak["UJI-D9"] = models.BarisKontrak{"ID": "UJI-D9", "TREATYID": "UJI-T9", "CLASSOFBUSINESS": "",
		"TREATYGROUPID": "UJI-TG9", "CEDING": "UJI-CEDING-9", "PROPORTIONTYPE": "Proportional"}
	u.g.GrupBisnis["UJI-TG9"] = "UJI GRUP BISNIS"
	if kode, isi := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D9"}); kode != http.StatusOK {
		t.Fatalf("pilih bisnis: %d %s", kode, isi)
	}
	h := u.g.Halaman[id]
	for j, harap := range map[string]string{
		"Quotation.BusinessName":     "UJI GRUP BISNIS",
		"Quotation.InsuredName":      "UJI-CEDING-9",
		"PolicyTreatyIn.InsuredName": "UJI-CEDING-9",
	} {
		if got := h.Ambil(j); got != harap {
			t.Errorf("tersimpan %s = %q, harap %q", j, got, harap)
		}
	}
}

func TestPilihBisnisClassOfBusinessTerisiTidakDiganti(t *testing.T) {
	u := baru(t)
	id := u.buat()
	u.g.Kontrak["UJI-D8"] = models.BarisKontrak{"ID": "UJI-D8", "TREATYID": "UJI-T8", "CLASSOFBUSINESS": "UJI BISNIS COB",
		"TREATYGROUPID": "UJI-TG8", "PROPORTIONTYPE": "Proportional"}
	u.g.GrupBisnis["UJI-TG8"] = "UJI GRUP LAIN"
	if kode, isi := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D8"}); kode != http.StatusOK {
		t.Fatalf("pilih bisnis: %d %s", kode, isi)
	}
	if got := u.g.Halaman[id].Ambil("Quotation.BusinessName"); got != "UJI BISNIS COB" {
		t.Errorf("Quotation.BusinessName = %q, harap Class of Business apa adanya (UJI BISNIS COB)", got)
	}
}
