package services

import (
	"context"
	"errors"
	"reflect"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// Tiket 71 - menu View Policy (Facultative).

// polisTiruan - PembacaPolis dalam memori: "RNM-1" punya dua versi; versi "0" berbaris.
type polisTiruan struct{ minta []string }

func (p *polisTiruan) VersiPolis(_ context.Context, noPolis string) ([]models.VersiPolis, error) {
	p.minta = append(p.minta, "versi|"+noPolis)
	if noPolis != "RNM-1" {
		return nil, repository.ErrPolisTidakAda
	}
	return []models.VersiPolis{{ProdKe: "1", CaseID: "EDM-2", SistemBaru: false}, {ProdKe: "0", CaseID: "NB-1", SistemBaru: true}}, nil
}

func (p *polisTiruan) RingkasanPolis(_ context.Context, noPolis, prodKe string) (models.RingkasanPolis, error) {
	p.minta = append(p.minta, "ringkasan|"+noPolis+"|"+prodKe)
	if noPolis != "RNM-1" || prodKe != "0" {
		return models.RingkasanPolis{}, repository.ErrPolisTidakAda
	}
	return models.RingkasanPolis{Kepala: models.KepalaPolis{PolicyNo: "RNM-1", PolicyType: "Master Policy"},
		Baris: []models.BarisPolis{{Rate: ".15"}, {Rate: "0,15"}, {Rate: ""}, {Rate: " 2.500 "}}}, nil
}

// TestJenisCariPolis - GetListViewPolicy_act blok 8: sama dengan tiket 70 kecuali opsi 5 / 6 awalan dan opsi 7 UPPER awalan;
// langkah 6 kutip dibuang.
func TestJenisCariPolis(t *testing.T) {
	a := repository.AwalanIDPega
	for nama, c := range map[string]struct{ opsi, filter, jenis, nilai string }{
		"nopolis":        {"2", "RNM-F1", "nopolis", "RNM-F1"},
		"id NB":          {"1", "NB-12", "idpega", a + "NB-12"},
		"id EDM":         {"1", "EDM-7", "edm", a + "EDM-7"},
		"RNM-E opsi 2":   {"2", "RNM-E.1", "noendors", "RNM-E.1"},
		"reff tetap %x%": {"4", "SLIP", "reff", "%SLIP%"},
		"insured awalan": {"5", "PT A", "insured", "PT A%"},
		"sob awalan":     {"6", "BRK", "sob", "BRK%"},
		"qq upper":       {"7", "qq", repository.CariPolisQQ, "qq%"},
		"opsi 5 RNM-E":   {"5", "RNM-E1", "insured", "RNM-E1%"},
		"kosong":         {"5", "", "", ""},
	} {
		j, n, err := jenisCariPolis(c.opsi, c.filter)
		if err != nil || j != c.jenis || n != c.nilai {
			t.Errorf("%s: %q %q %v", nama, j, n, err)
		}
	}
	if _, _, err := jenisCariPolis("3", "X"); !errors.Is(err, ErrMasukanFollowing) {
		t.Errorf("opsi 3: %v", err)
	}
	for masuk, mau := range map[string]string{`PT "A'`: "PT A", ` 'x' `: "x", `""`: ""} {
		if got := filterPolis(masuk); got != mau {
			t.Errorf("filter %q: %q", masuk, got)
		}
	}
}

// TestCariPolis - tanpa penjaga case (hanya login), kutip dibuang lalu filter kosong -> tanpa kueri, halaman 20.
func TestCariPolis(t *testing.T) {
	ctx := context.Background()
	f := &followingTiruan{baris: []models.BarisFollowing{{CaseID: "NB-1"}}}
	svc := Baru(nil).DenganFollowing(f)
	lain := inti.Pelaku{AkunID: "UJI-LAIN"} // tanpa peran workbasket apa pun (K71-3)
	h, err := svc.CariPolis(ctx, lain, "6", ` "BRK" `, 2)
	if err != nil || h.Total != 41 || h.Halaman != 2 || !reflect.DeepEqual(f.cari, []string{"sob|BRK%|20|20"}) {
		t.Errorf("%+v %v %v", h, err, f.cari)
	}
	if h, err := svc.CariPolis(ctx, lain, "2", `"'`, 1); err != nil || h.Baris == nil || len(h.Baris) != 0 || len(f.cari) != 1 {
		t.Errorf("kutip saja: %+v %v", h, err)
	}
	if _, err := svc.CariPolis(ctx, inti.Pelaku{}, "2", "X", 1); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("401: %v", err)
	}
	if _, err := Baru(nil).CariPolis(ctx, lain, "2", "X", 1); !errors.Is(err, ErrPolisTanpaDatabase) {
		t.Errorf("503: %v", err)
	}
}

// TestVersiDanRingkasanPolis - versi apa adanya; ringkasan: Policy Type kode, Rate dinormalkan; 404 / 400 / 401.
func TestVersiDanRingkasanPolis(t *testing.T) {
	ctx := context.Background()
	p := &polisTiruan{}
	svc := Baru(nil).DenganPolis(p)
	lain := inti.Pelaku{AkunID: "UJI-LAIN"}
	v, err := svc.VersiPolis(ctx, lain, " RNM-1 ")
	if err != nil || len(v) != 2 || v[1] != (models.VersiPolis{ProdKe: "0", CaseID: "NB-1", SistemBaru: true}) {
		t.Errorf("versi %+v %v", v, err)
	}
	if _, err := svc.VersiPolis(ctx, lain, "RNM-9"); !errors.Is(err, ErrPolisTidakAda) {
		t.Errorf("404: %v", err)
	}
	r, err := svc.RingkasanPolis(ctx, lain, "RNM-1", "0")
	rate := []string{}
	for _, b := range r.Baris {
		rate = append(rate, b.Rate)
	}
	if err != nil || r.Kepala.PolicyType != "1" || !reflect.DeepEqual(rate, []string{"0.15", "0,15", "", "2.500"}) {
		t.Errorf("ringkasan %+v %v %q", r.Kepala, err, rate)
	}
	if _, err := svc.RingkasanPolis(ctx, lain, "RNM-1", "7"); !errors.Is(err, ErrPolisTidakAda) {
		t.Errorf("404 versi: %v", err)
	}
	for nama, c := range map[string]struct{ no, prod string }{"nopolis": {" ", "0"}, "prodke": {"RNM-1", ""}} {
		if _, err := svc.RingkasanPolis(ctx, lain, c.no, c.prod); !errors.Is(err, ErrMasukanPolis) {
			t.Errorf("400 %s: %v", nama, err)
		}
	}
	if _, err := svc.VersiPolis(ctx, inti.Pelaku{}, "RNM-1"); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("401: %v", err)
	}
	if _, err := Baru(nil).VersiPolis(ctx, lain, "RNM-1"); !errors.Is(err, ErrPolisTanpaDatabase) {
		t.Errorf("503: %v", err)
	}
	if !reflect.DeepEqual(p.minta, []string{"versi|RNM-1", "versi|RNM-9", "ringkasan|RNM-1|0", "ringkasan|RNM-1|7"}) {
		t.Errorf("permintaan %q", p.minta)
	}
}
