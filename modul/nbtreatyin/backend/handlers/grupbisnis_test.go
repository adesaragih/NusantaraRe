package handlers_test

// Perintah work owner 10-10-2026 ("ganti yang di inbox group bisnis jadi treaty bisnis, script yang select group
// bisnis hapus aja"): kolom inbox kini Treaty Business (T_GENERAL_POLIS_TREATY.TREATY_GROUP_NAME), dan Choose
// Business TIDAK lagi mencari nama grup bisnis di tabel BUSINESS (TREATYGROUP.COAID = BUSINESS.BUSINESSGROUPID,
// keputusan 06-10-2026 yang dicabut). Quotation.BusinessName = CLASSOFBUSINESS kontrak apa adanya (preACT langkah 3),
// kosong bila kosong. Insured Name = CEDING kontrak terpilih (keputusan 06-10-2026, tetap).

import (
	"context"
	"net/http"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func TestPilihBisnisTanpaClassOfBusinessTidakMencariGrupBisnis(t *testing.T) {
	u := baru(t)
	id := u.buat()
	u.g.Kontrak["UJI-D9"] = models.BarisKontrak{"ID": "UJI-D9", "TREATYID": "UJI-T9", "CLASSOFBUSINESS": "",
		"TREATYGROUPID": "UJI-TG9", "TREATYGROUP": "UJI GRUP TREATY", "CEDING": "UJI-CEDING-9", "PROPORTIONTYPE": "Proportional"}
	if kode, isi := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D9"}); kode != http.StatusOK {
		t.Fatalf("pilih bisnis: %d %s", kode, isi)
	}
	h := u.g.Halaman[id]
	for j, harap := range map[string]string{
		"Quotation.BusinessName":         "",
		"PolicyTreatyIn.TreatyGroupName": "UJI GRUP TREATY",
		"Quotation.InsuredName":          "UJI-CEDING-9",
		"PolicyTreatyIn.InsuredName":     "UJI-CEDING-9",
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
	if kode, isi := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D8"}); kode != http.StatusOK {
		t.Fatalf("pilih bisnis: %d %s", kode, isi)
	}
	if got := u.g.Halaman[id].Ambil("Quotation.BusinessName"); got != "UJI BISNIS COB" {
		t.Errorf("Quotation.BusinessName = %q, harap Class of Business apa adanya (UJI BISNIS COB)", got)
	}
}

func TestInboxKolomTreatyBusiness(t *testing.T) {
	u := baru(t)
	id := u.buat()
	u.g.Kontrak["UJI-D7"] = models.BarisKontrak{"ID": "UJI-D7", "TREATYID": "UJI-T7", "CLASSOFBUSINESS": "UJI BISNIS COB",
		"TREATYGROUPID": "UJI-TG7", "TREATYGROUP": "UJI GRUP TREATY 7", "PROPORTIONTYPE": "Proportional"}
	if kode, isi := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D7"}); kode != http.StatusOK {
		t.Fatalf("pilih bisnis: %d %s", kode, isi)
	}
	daftar, err := u.g.DaftarKasus(context.Background(), models.SaringanKasus{})
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range daftar {
		if r.ID == id {
			if r.TreatyGroupName != "UJI GRUP TREATY 7" {
				t.Errorf("kolom Treaty Business = %q, harap TreatyGroupName kontrak (UJI GRUP TREATY 7)", r.TreatyGroupName)
			}
			return
		}
	}
	t.Fatalf("kasus %s tidak ada di daftar inbox", id)
}
