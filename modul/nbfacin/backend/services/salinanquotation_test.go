package services

import (
	"context"
	"strings"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
)

// TestSalinanBisnisQuotation - tiket 76b (migrasi 253): Save General membawa salinan bisnis QuotationData dari BUSINESS baris kasus
// (kasusBisnisDari) dan BusinessType_DeT; BusinessType2 LIFE hanya untuk "Life"; BUSINESS tak unik -> kode / old id / panel "" (kolom
// tetap), tipe "UNKNOWN" (bawaan DT); lebih lebar dari kolom rancangan -> ""; tanpa pembaca BUSINESS -> nil (kolom tetap).
func TestSalinanBisnisQuotation(t *testing.T) {
	ctx := context.Background()
	akun := inti.Pelaku{AkunID: "UJI-USER"}
	jam := func() time.Time { return time.Date(2026, 10, 3, 1, 0, 0, 0, time.UTC) }
	simpan := func(cob string, bisnis []models.BarisBisnisLini, tanpaTOL bool) *models.BisnisQuotation {
		ada := kasusUji()
		ada["UJI-NB-1"].Opportunity.ClassOfBusiness, ada["UJI-NB-1"].Opportunity.GroupBusinessID = cob, "3"
		svc := Baru(nil).DenganKasus(kasusTiruan{ada: ada}).DenganTransaksi(tanpaTx).DenganJam(jam)
		if !tanpaTOL {
			svc.DenganTableOfLimit(tolTiruan{bisnis: bisnis})
		}
		if _, err := svc.SimpanGeneral(ctx, akun, "UJI-NB-1", IsianGeneral{PolicyType: "0"}); err != nil {
			t.Fatal(err)
		}
		return ada["UJI-NB-1"].General.Bisnis
	}
	if b := simpan("UJI FIRE", satuBisnis("006", "22"), false); b == nil || *b != (models.BisnisQuotation{Code: "10999",
		Name: "UJI FIRE", OldID: "22", GroupPanel: "006", Type: "FireStyle1", Type2: "NONLIFE"}) {
		t.Errorf("FIRE: %+v", b)
	}
	if b := simpan("UJI LIFE", satuBisnis("009", "L3"), false); b == nil || b.Type != "Life" || b.Type2 != "LIFE" {
		t.Errorf("Life: %+v", b)
	}
	dua := append(satuBisnis("006", "22"), satuBisnis("006", "23", "10998")...)
	if b := simpan("UJI FIRE", dua, false); b == nil || *b != (models.BisnisQuotation{Name: "UJI FIRE", Type: "UNKNOWN", Type2: "NONLIFE"}) {
		t.Errorf("BUSINESS tak unik: %+v", b)
	}
	if b := simpan(strings.Repeat("N", 501), satuBisnis("006", "22", strings.Repeat("9", 51)), false); b == nil || b.Name != "" ||
		b.Code != "" || b.OldID != "22" {
		t.Errorf("lebar: %+v", b)
	}
	if b := simpan("UJI FIRE", nil, true); b != nil {
		t.Errorf("tanpa pembaca BUSINESS: %+v", b)
	}
}
