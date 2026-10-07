package handlers_test

// Perintah work owner 06-10-2026: saat Choose Business, SpreadingRiskList(1) ikut terset dari view
// (TreatyInputPctCommSpreading 2.1.1.1: TreatyType = SPREADINGTYPEID, TreatyName = SPREADINGTYPE).

import (
	"net/http"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func TestPilihBisnisMengisiSpreadingDariView(t *testing.T) {
	u := baru(t)
	id := u.buat()
	u.g.Kontrak["UJI-D5"] = models.BarisKontrak{"ID": "UJI-D5", "TREATYID": "UJI-T5", "PROPORTIONTYPE": "Proportional",
		"TREATYTYPE": "UJI-QS", "TREATYGROUP": "UJI-GRUP", "RIOGR": "1", "RIONR": "2",
		"SPREADINGTYPEID": "UJI-ST5", "SPREADINGTYPE": "UJI Spreading Lima"}
	if kode, isi := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D5"}); kode != http.StatusOK {
		t.Fatalf("pilih bisnis: %d %s", kode, isi)
	}
	d := u.g.Halaman[id].AmbilDaftar(models.DaftarSpreading)
	if len(d) == 0 || d[0]["TreatyType"] != "UJI-ST5" || d[0]["TreatyName"] != "UJI Spreading Lima" {
		t.Fatalf("SpreadingRiskList tersimpan = %v, harap baris 1 TreatyType UJI-ST5 / TreatyName UJI Spreading Lima", d)
	}
}
