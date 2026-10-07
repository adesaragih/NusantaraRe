package handlers

import (
	"context"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

type tahapTiruanH struct{ t models.TahapKasus }

func (h tahapTiruanH) TahapKasus(context.Context, string) (models.TahapKasus, error) { return h.t, nil }

// TestBacaSajaHandler - GET kasus: bolehUbah / alasanBacaSaja; PUT general di kasus Team Leader -> 409 (GET tetap 200);
// tahap marketing + workbasket -> bolehUbah true.
func TestBacaSajaHandler(t *testing.T) {
	jalur := "/api/nbfacin/kasus/UJI-NB-1"
	tl := layananKasus().DenganTahap(tahapTiruanH{models.TahapKasus{Position: "ReasFacInTeamLeader", FlagOnGoingPolicy: "0"}})
	kode, isi := mintaPeran(t, tl, "GET", jalur, "", "UJI-MKT", "ReasFacInMarketing")
	if kode != 200 || !strings.Contains(isi, `"bolehUbah":false,"alasanBacaSaja":"Case baca-saja: sudah tidak berada di tahap marketing Offer (posisi \"ReasFacInTeamLeader\")."`) {
		t.Fatalf("GET TL: %d %s", kode, isi)
	}
	if kode, isi := mintaPeran(t, tl, "PUT", jalur+"/general", `{"reffNumber":"UJI"}`, "UJI-MKT", "ReasFacInMarketing"); kode != 409 ||
		!strings.Contains(isi, "baca-saja") {
		t.Errorf("PUT TL: %d %s", kode, isi)
	}
	mkt := layananKasus().DenganTahap(tahapTiruanH{models.TahapKasus{Position: "Offer", FlagOnGoingPolicy: "0"}})
	if kode, isi := mintaPeran(t, mkt, "GET", jalur, "", "UJI-MKT", "ReasFacInMarketing"); kode != 200 ||
		!strings.Contains(isi, `"bolehUbah":true,"alasanBacaSaja":""`) {
		t.Errorf("GET marketing: %d %s", kode, isi)
	}
	if kode, isi := mintaPeran(t, mkt, "PUT", jalur+"/general", `{"reffNumber":"UJI"}`, "UJI-UW", "ReasFacInUnderwriting"); kode != 403 {
		t.Errorf("PUT tanpa workbasket: %d %s", kode, isi)
	}
}
