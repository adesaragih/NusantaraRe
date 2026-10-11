package handlers

import (
	"context"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

type tahapTiruanH struct{ t models.TahapKasus }

func (h tahapTiruanH) TahapKasus(context.Context, string) (models.TahapKasus, error) { return h.t, nil }

// TestBacaSajaHandler - GET kasus: bolehUbah / alasanBacaSaja; PUT general di kasus Junior UW -> 409 (GET tetap 200);
// tahap marketing + workbasket -> bolehUbah true; posisi Team Leader: TL boleh, marketing 403 (tiket 58).
func TestBacaSajaHandler(t *testing.T) {
	jalur := "/api/nbfacin/kasus/UJI-NB-1"
	uw := layananKasus().DenganTahap(tahapTiruanH{models.TahapKasus{Position: "ReasFacInJuniorUnderwritingB", FlagOnGoingPolicy: "0"}})
	kode, isi := mintaPeran(t, uw, "GET", jalur, "", "UJI-MKT", "ReasFacInMarketing")
	if kode != 200 || !strings.Contains(isi, `"bolehUbah":false,"alasanBacaSaja":"Case baca-saja: sudah tidak berada di tahap marketing / Team Leader (posisi \"ReasFacInJuniorUnderwritingB\")."`) {
		t.Fatalf("GET UW: %d %s", kode, isi)
	}
	if kode, isi := mintaPeran(t, uw, "PUT", jalur+"/general", `{"reffNumber":"UJI"}`, "UJI-MKT", "ReasFacInMarketing"); kode != 409 ||
		!strings.Contains(isi, "baca-saja") {
		t.Errorf("PUT UW: %d %s", kode, isi)
	}
	tl := layananKasus().DenganTahap(tahapTiruanH{models.TahapKasus{Position: "ReasFacInTeamLeader", FlagOnGoingPolicy: "0"}})
	if kode, isi := mintaPeran(t, tl, "GET", jalur, "", "UJI-TL", "ReasFacInTeamLeader"); kode != 200 || !strings.Contains(isi, `"bolehUbah":true,"alasanBacaSaja":""`) {
		t.Errorf("GET TL oleh TL: %d %s", kode, isi)
	}
	if kode, isi := mintaPeran(t, tl, "PUT", jalur+"/general", `{"reffNumber":"UJI"}`, "UJI-MKT", "ReasFacInMarketing"); kode != 403 {
		t.Errorf("PUT TL oleh marketing: %d %s", kode, isi)
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
