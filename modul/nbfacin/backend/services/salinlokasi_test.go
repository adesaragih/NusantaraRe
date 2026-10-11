package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
)

var templatLokasiUji = []models.TemplateSpreading{{TreatyType: treatyORS, TreatyName: namaORS, SharePercentage: d("60")},
	{TreatyType: treatyFACOUT, TreatyName: namaFACOUT, SharePercentage: d("40"), FlagSpreading: "Open Fac Market"}}

// TestSalinSpreadingLokasiFire - K73-40 (FIRE juga): hanya lokasi idxLokasi yang spreading-nya diganti template (rumus salin
// case-level: TSI / premi / TSIGross menurut share); lokasi lain tetap; template jawaban []; tidak menulis.
func TestSalinSpreadingLokasiFire(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-UW"}
	k := kasusSatuCoverage("")
	k.Lokasi = append(k.Lokasi, salinKasusSpreading(k).Lokasi...)
	k.Lokasi[0].Items[0].Coverages[0].Spreading = []models.BarisSpreading{{TreatyType: "LAMA", SharePercentage: d("100")}}
	src := &spreadingTiruan{baca: []models.KasusSpreading{k}}
	v, err := Baru(nil).DenganSpreading(src).SalinSpreadingLokasi(ctx, akun, "UJI-NB-1", d("50"), 2, templatLokasiUji)
	if err != nil {
		t.Fatal(err)
	}
	if b := v.Lokasi[0].Items[0].Coverages[0].Spreading; len(b) != 1 || b[0].TreatyType != "LAMA" {
		t.Errorf("lokasi 1 berubah: %v", b)
	}
	b := v.Lokasi[1].Items[0].Coverages[0].Spreading
	if len(b) != 2 || b[0].TreatyType != treatyORS || t8(b[0].TSISpreaded) != "300" || t8(b[0].PremiumSpreaded) != "2.4" ||
		t8(b[0].TSIGrossSpreaded) != "360" || t8(b[1].TSISpreaded) != "200" {
		t.Errorf("lokasi 2: %v", b)
	}
	if len(v.Template) != 0 || len(src.tulis) != 0 {
		t.Errorf("template %v / tulis %d", v.Template, len(src.tulis))
	}
	for nama, c := range map[string]struct {
		idx    int
		templ  []models.TemplateSpreading
		pelaku inti.Pelaku
		mau    error
		pesan  string
	}{
		"idx 0":           {0, templatLokasiUji, akun, ErrMasukanSpreading, "idxLokasi 0 di luar 1..2"},
		"idx 3":           {3, templatLokasiUji, akun, ErrMasukanSpreading, "idxLokasi 3"},
		"flag kosong":     {1, []models.TemplateSpreading{{TreatyType: treatyFACOUT, SharePercentage: d("100")}}, akun, ErrMasukanSpreading, "template[0]." + pesanFlagSpreadingWajib},
		"flag bukan opsi": {1, []models.TemplateSpreading{{TreatyType: treatyFACOUT, SharePercentage: d("100"), FlagSpreading: "UJI"}}, akun, ErrMasukanSpreading, `flagSpreading "UJI" bukan pilihan`},
		"flag huruf":      {1, []models.TemplateSpreading{{TreatyType: treatyFACOUT, SharePercentage: d("100"), FlagSpreading: "sf-mr 2023"}}, akun, ErrMasukanSpreading, "bukan pilihan"},
		"share > 100":     {1, []models.TemplateSpreading{{TreatyType: treatyORS, SharePercentage: d("101")}}, akun, ErrMasukanSpreading, "paling besar 100"},
		"template kosong": {1, nil, akun, ErrMasukanSpreading, "template wajib"},
		"tanpa identitas": {1, templatLokasiUji, inti.Pelaku{}, inti.ErrTanpaIdentitas, ""},
	} {
		_, err := Baru(nil).DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{k}}).SalinSpreadingLokasi(ctx, c.pelaku, "UJI-NB-1",
			d("50"), c.idx, c.templ)
		if !errors.Is(err, c.mau) || !strings.Contains(err.Error(), c.pesan) {
			t.Errorf("%s: %v", nama, err)
		}
	}
}

// TestSalinSpreadingAnekaLokasi - Aneka langkah 6: objek okupasi lokasi idxLokasi saja, TSI / premi menurut share tanpa TSIGross.
func TestSalinSpreadingAnekaLokasi(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-UW"}
	svc, a, tulis, _ := layananSpreadAneka("30")
	a.ada["K"] = append(lokasiSpreadUji(), lokasiSpreadUji()...)
	v, err := svc.SalinSpreadingAnekaLokasi(ctx, akun, "K", d("50"), 2, templatLokasiUji)
	if err != nil {
		t.Fatal(err)
	}
	if b := v.Lokasi[0].Okupasi[0].Objek[0].Coverages[0].Spreading; len(b) != 2 || b[0].TreatyType != "OR" {
		t.Errorf("lokasi 1 berubah: %v", b)
	}
	b := v.Lokasi[1].Okupasi[0].Objek[0].Coverages[0].Spreading
	if len(b) != 2 || barisTeks(b[0]) != "10007|60|300|2.4||240" || barisTeks(b[1]) != "10015|40|200|1.6||160" || len(*tulis) != 0 || len(v.Template) != 0 {
		t.Errorf("lokasi 2: %v (tulis %d)", b, len(*tulis))
	}
	if g := v.Lokasi[1].ObjekGolf[0].Coverages[0].Spreading; len(g) != 1 || g[0].TreatyType != "OR" {
		t.Errorf("objek Golf case bukan Golf tersentuh: %v", g)
	}
	if _, err := svc.SalinSpreadingAnekaLokasi(ctx, akun, "K", d("50"), 3, templatLokasiUji); !errors.Is(err, ErrMasukanSpreading) {
		t.Errorf("idx 3: %v", err)
	}
}

// TestSalinLokasiTahap - tahap selain uw -> 409 (penjaga salin case-level, kontainer IsUWSpread).
func TestSalinLokasiTahap(t *testing.T) {
	svc := Baru(nil).DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{kasusSatuCoverage("")}}).DenganTahap(tahapUji)
	for _, id := range []string{"UJI-MKT2", "UJI-TL"} {
		_, err := svc.SalinSpreadingLokasi(context.Background(), inti.Pelaku{AkunID: "UJI", Peran: []string{"ReasFacInMarketing",
			"ReasFacInTeamLeader"}}, id, d("50"), 1, templatLokasiUji)
		if !errors.Is(err, ErrSpreadingHanyaUW) {
			t.Errorf("%s: %v", id, err)
		}
		if _, err := Baru(nil).DenganTahap(tahapUji).SalinSpreadingAnekaLokasi(context.Background(), inti.Pelaku{AkunID: "UJI",
			Peran: []string{"ReasFacInMarketing", "ReasFacInTeamLeader"}}, id, d("50"), 1, templatLokasiUji); !errors.Is(err, ErrSpreadingHanyaUW) {
			t.Errorf("Aneka %s: %v", id, err)
		}
	}
}
