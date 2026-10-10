package handlers_test

import (
	"net/http"
	"testing"

	"nusantarare/modul/claimprop/backend/handlers"
	"nusantarare/modul/claimprop/backend/models"
	"nusantarare/modul/claimprop/backend/tiruan"
)

// Keputusan work owner 09-10-2026: menu Komite Claim Prop dibuang; kasus komite tampil di tabel bawah inbox Claim Prop,
// dan tabel itu HANYA muncul bila akun memegang workbasket yang tercantum di roster EMAILKOMITE PROP aktif.
func TestHakKomiteMenurutWorkbasketRosterProp(t *testing.T) {
	u := baruUji(t)
	u.a.Roster = []tiruan.AnggotaRoster{
		{AnggotaKomite: models.AnggotaKomite{ID: "1", OperatorID: "UJI-WB-KOMITE", Jabatan: "UJI-JABATAN"},
			Batas: "-1", Sts: models.STSKlaimProp},
		{AnggotaKomite: models.AnggotaKomite{ID: "2", OperatorID: "UJI-WB-FACIN", Jabatan: "UJI-JABATAN"},
			Batas: "-1", Sts: "FACIN"},
	}
	hak := func(peran string) map[string]any {
		t.Helper()
		kode, out := u.minta(http.MethodGet, handlers.Prefix+"/hak", teknik, peran, nil)
		u.wajib(kode, http.StatusOK, out, "hak "+peran)
		return out
	}
	if out := hak("UJI-LAIN,UJI-WB-KOMITE"); out["komite"] != true {
		t.Fatalf("pemegang workbasket roster PROP: tabel komite tampil: %v", out)
	}
	if out := hak("UJI-WB-FACIN"); out["komite"] != false {
		t.Fatalf("workbasket roster FACIN bukan komite PROP: %v", out)
	}
	if out := hak(models.WorkbasketAcceptation); out["komite"] != false || out["workbasketTeknik"] != true {
		t.Fatalf("tanpa workbasket komite: tabel komite tidak tampil: %v", out)
	}
}
