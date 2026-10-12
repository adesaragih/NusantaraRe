package services

import (
	"context"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// TestRingkasanTotalPolis - tiket 76 (keputusan work owner 12-10-2026): total = hitungan Go yang sudah ada (totalTSIUW; ringkasanRNM
// + totalKeIDR tab Cedant), dibulatkan 8 desimal (lolos mask bind PolaAngkaMasuk); kurs tidak ada -> total itu NULL tanpa galat;
// tanpa pembaca spreading -> nil.
func TestRingkasanTotalPolis(t *testing.T) {
	ctx := context.Background()
	teks := func(x *apd.Decimal) string {
		if x == nil {
			return "<nil>"
		}
		return utils.FormatDecimal(x)
	}
	hitung := func(ubah func(*models.KasusSpreading), kurs map[string]*apd.Decimal) (*models.RingkasanTotalPolis, models.KasusSpreading) {
		sp := spreadingSah() // PS 50, IDR, TSILiability 1000
		ubah(&sp)
		svc := layananSubmitUW(submitUji(), posisiJuniorUWB, sp)
		svc.DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{sp}, idUang: map[string]string{"IDR": idMataUangIDR, "USD": "10001"},
			kurs: kurs})
		return svc.ringkasanTotalPolis(ctx, "UJI-NB-1"), sp
	}
	g, sp := hitung(func(*models.KasusSpreading) {}, nil)
	_, jumlah, _ := ringkasanRNM(sp)
	if g == nil || teks(g.TotalTSINusaRe) != "500" || teks(g.TotalTSITopRisk) != "0" || teks(g.TotalTSINusaReSpreading) != "500" ||
		teks(g.TotalPremiNusaRe) != teks(simpan8(jumlah[0].Premi)) || g.TotalPremiNusaRe == nil {
		t.Fatalf("IDR: %+v", g)
	}
	for _, x := range []*apd.Decimal{g.TotalTSINusaRe, g.TotalTSITopRisk, g.TotalTSINusaReSpreading, g.TotalPremiNusaRe} {
		if !repository.PolaAngkaMasuk.MatchString(teks(x)) {
			t.Errorf("mask bind: %s", teks(x))
		}
	}
	top := func(k *models.KasusSpreading) { k.Lokasi[0].IsTopRisk = true; k.PercentShare = d("33.333333333") }
	if g, _ := hitung(top, nil); teks(g.TotalTSITopRisk) != "333.33333333" || teks(g.TotalTSINusaRe) != "333.33333333" {
		t.Errorf("top risk / 8 desimal: %+v", g)
	}
	usd := func(k *models.KasusSpreading) { k.Lokasi[0].Items[0].Currency = "USD" }
	if g, _ := hitung(usd, nil); g == nil || g.TotalTSINusaRe != nil || g.TotalTSITopRisk != nil || g.TotalTSINusaReSpreading != nil ||
		g.TotalPremiNusaRe != nil {
		t.Errorf("kurs tidak ada -> NULL: %+v", g)
	}
	if g := Baru(nil).ringkasanTotalPolis(ctx, "UJI-NB-1"); g != nil {
		t.Errorf("tanpa pembaca: %+v", g)
	}
}

// TestDateOfferSubmit - DATE_OFFER hanya Submit tahap marketing, "MM/dd/yyyy hh:mm a" GMT (ViewOfferFacInUW_PreDT baris 2).
func TestDateOfferSubmit(t *testing.T) {
	jam := func(s string) *Service {
		w, _ := time.Parse(time.RFC3339, s)
		return Baru(nil).DenganJam(func() time.Time { return w })
	}
	for masuk, mau := range map[string]string{"2026-10-07T22:04:05+07:00": "10/07/2026 03:04 PM",
		"2026-10-08T06:30:00+07:00": "10/07/2026 11:30 PM", "2026-10-08T07:05:00+07:00": "10/08/2026 12:05 AM"} {
		if got := jam(masuk).dateOfferSubmit(TahapMarketing); got != mau {
			t.Errorf("%s: %q, mau %q", masuk, got, mau)
		}
	}
	for _, tahap := range []string{TahapTeamLeader, TahapUW, TahapBinding} {
		if got := jam("2026-10-07T22:04:05+07:00").dateOfferSubmit(tahap); got != "" {
			t.Errorf("%s: %q", tahap, got)
		}
	}
}
