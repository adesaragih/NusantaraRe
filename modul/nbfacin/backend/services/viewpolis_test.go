package services

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// Tiket 71 perluasan - View Policy tab Offer Facultative / Treaty / Life.

// viewTiruan - PembacaViewPolis dalam memori: mencatat "tab|jenis|nilai|offset|batas", satu baris, total 41.
type viewTiruan struct{ cari []string }

func (v *viewTiruan) catat(tab, jenis, nilai string, offset, batas int) {
	v.cari = append(v.cari, fmt.Sprintf("%s|%s|%s|%d|%d", tab, jenis, nilai, offset, batas))
}

func (v *viewTiruan) CariOfferFac(_ context.Context, jenis, nilai string, offset, batas int) ([]models.BarisOfferFac, int, error) {
	v.catat("offer", jenis, nilai, offset, batas)
	return []models.BarisOfferFac{{CaseID: "NB-1"}}, 41, nil
}

func (v *viewTiruan) CariPolisTreaty(_ context.Context, jenis, nilai string, offset, batas int) ([]models.BarisPolisTreaty, int, error) {
	v.catat("treaty", jenis, nilai, offset, batas)
	return []models.BarisPolisTreaty{{CaseID: "NB-1"}}, 41, nil
}

func (v *viewTiruan) CariPolisLife(_ context.Context, jenis, nilai string, offset, batas int) ([]models.BarisPolisLife, int, error) {
	v.catat("life", jenis, nilai, offset, batas)
	return []models.BarisPolisLife{{CaseID: "NBLF-1"}}, 41, nil
}

func (v *viewTiruan) CariOfferLife(_ context.Context, jenis, nilai string, offset, batas int) ([]models.BarisOfferLife, int, error) {
	v.catat("offerlife", jenis, nilai, offset, batas)
	return []models.BarisOfferLife{{CaseID: "NBLF-1"}}, 41, nil
}

// BerkasPolisTreaty - "RNM-T0" NB (prodke 0), "RNM-T2" EDM, lainnya tanpa berkas; nomor yang diminta dicatat.
func (v *viewTiruan) BerkasPolisTreaty(_ context.Context, noPolis string) (models.BerkasPolisTreaty, bool, error) {
	v.cari = append(v.cari, "berkas|"+noPolis)
	switch noPolis {
	case "RNM-T0":
		return models.BerkasPolisTreaty{Modul: models.ModulBerkasPolisTreaty(0), Kasus: "NBT-1"}, true, nil
	case "RNM-T2":
		return models.BerkasPolisTreaty{Modul: models.ModulBerkasPolisTreaty(2), Kasus: "EDMT-2"}, true, nil
	}
	return models.BerkasPolisTreaty{}, false, nil
}

// TestJenisCariOfferFac - View Option tab Offer: 1 hanya NB- / RNW- (9.2), 2 =, 4 %x%, 5 / 6 awalan; 3 / 7 / lain 400 walau filter
// kosong; filter kosong -> tanpa kueri.
func TestJenisCariOfferFac(t *testing.T) {
	for _, c := range []struct{ opsi, filter, jenis, nilai string }{
		{"1", "NB-123", repository.OfferFacIDPega, "NB-123"},
		{"1", "RNW-9", repository.OfferFacIDPega, "RNW-9"},
		{"1", "RNM-F1", "", ""},
		{"2", "RNM-OF-1", repository.OfferFacNoOffer, "RNM-OF-1"},
		{"4", "SLIP", repository.OfferFacReff, "%SLIP%"},
		{"5", "PT A", repository.OfferFacInsured, "PT A%"},
		{"6", "BRK", repository.OfferFacSOB, "BRK%"},
		{"2", "", "", ""},
	} {
		j, n, err := jenisCariOfferFac(c.opsi, c.filter)
		if err != nil || j != c.jenis || n != c.nilai {
			t.Errorf("%s %q: %q %q %v", c.opsi, c.filter, j, n, err)
		}
	}
	for _, opsi := range []string{"", "0", "3", "7", "8", " 1"} {
		for _, f := range []string{"", "X"} {
			if _, _, err := jenisCariOfferFac(opsi, f); !errors.Is(err, ErrMasukanFollowing) {
				t.Errorf("400 %q %q: %v", opsi, f, err)
			}
		}
	}
}

// TestJenisCariTreaty - 10.1-10.7 berurutan: opsi 2 panjang 28 karakter (rune) -> NO ENDORS menimpa NOPOLIS; opsi 1 NB- (bukan
// RNW-) / EDMT-; 3 / 4 / 5; opsi di luar 1-5 400.
func TestJenisCariTreaty(t *testing.T) {
	e28 := strings.Repeat("E", 28)
	e28rune := strings.Repeat("E", 27) + "é" // 28 karakter, 29 byte
	for _, c := range []struct{ opsi, filter, jenis, nilai string }{
		{"2", "RNM-T1", repository.TreatyNoPolis, "RNM-T1"},
		{"2", e28, repository.TreatyNoEndors, e28},
		{"2", e28rune, repository.TreatyNoEndors, e28rune},
		{"2", e28[:27], repository.TreatyNoPolis, e28[:27]},
		{"2", e28 + "E", repository.TreatyNoPolis, e28 + "E"},
		{"1", "NB-7", repository.TreatyIDPega, "NB-7"},
		{"1", "EDMT-7", repository.TreatyEDMT, "EDMT-7"},
		{"1", "RNW-7", "", ""},
		{"1", "7", "", ""},
		{"3", "M-1", repository.TreatyMaster, "M-1"},
		{"4", "PT A", repository.TreatyInsured, "%PT A%"},
		{"5", "BRK", repository.TreatySOB, "%BRK%"},
		{"4", "", "", ""},
	} {
		j, n, err := jenisCariTreaty(c.opsi, c.filter)
		if err != nil || j != c.jenis || n != c.nilai {
			t.Errorf("%s %q: %q %q %v", c.opsi, c.filter, j, n, err)
		}
	}
	for _, opsi := range []string{"", "0", "6", "7", "x"} {
		if _, _, err := jenisCariTreaty(opsi, "X"); !errors.Is(err, ErrMasukanFollowing) {
			t.Errorf("400 %q: %v", opsi, err)
		}
	}
}

// TestJenisCariLife - 11.1-11.6 berurutan (yang belakangan menimpa; SKP `[dugaan]` tak berefek), KECUALI panjang tepat 5 karakter:
// 11.2 menang atas 11.6 (keputusan work owner 11-10-2026) - "00012" / "ABC12" -> digit akhir; RNM 5 karakter -> 11.2; NBLF- -> 11.3.
func TestJenisCariLife(t *testing.T) {
	endors := "RNM-L.10.2026.00001-E.00001" // 27 karakter
	for _, c := range []struct{ filter, jenis, nilai string }{
		{"RNM-L.10.2026.00001", repository.LifeNoPolis, "RNM-L.10.2026.00001"},
		{strings.Repeat("R", 19) + "RNM", repository.LifeNoPolis, strings.Repeat("R", 19) + "RNM"},
		{"RNM12", repository.LifeDigitAkhir, "%RNM12"},
		{"00012", repository.LifeDigitAkhir, "%00012"},
		{"ABC12", repository.LifeDigitAkhir, "%ABC12"},
		{"ABC123", repository.LifeSOB, "%ABC123%"},
		{"ABé1", repository.LifeSOB, "%ABé1%"},
		{"ABé12", repository.LifeDigitAkhir, "%ABé12"},
		{"EDML-", repository.LifeEDML, "EDML-"},
		{"NBLF-12345", repository.LifeIDPega, "NBLF-12345"},
		{"NBLF-", repository.LifeIDPega, "NBLF-"},
		{endors, repository.LifeNoEndors, endors},
		{strings.Repeat("R", 20) + "RNM", repository.LifeNoEndors, strings.Repeat("R", 20) + "RNM"},
		{"EDML-123", repository.LifeEDML, "EDML-123"},
		{"PT ASURANSI", repository.LifeSOB, "%PT ASURANSI%"},
		{"rnm-li", repository.LifeSOB, "%rnm-li%"},
		{"", "", ""},
	} {
		if j, n := jenisCariLife(c.filter); j != c.jenis || n != c.nilai {
			t.Errorf("%q: %q %q", c.filter, j, n)
		}
	}
}

// TestJenisCariOfferLife - 12.1-12.3 berurutan: OFR- menimpa NBLF-; selain keduanya CEDING / POLICYHOLDER.
func TestJenisCariOfferLife(t *testing.T) {
	for _, c := range []struct{ filter, jenis, nilai string }{
		{"NBLF-12", repository.OfferLifeIDPega, "NBLF-12"},
		{"OFR-34", repository.OfferLifeOFR, "OFR-34"},
		{"NBLF-1 OFR-2", repository.OfferLifeOFR, "NBLF-1 OFR-2"},
		{"PT A", repository.OfferLifeCeding, "%PT A%"},
		{"", "", ""},
	} {
		if j, n := jenisCariOfferLife(c.filter); j != c.jenis || n != c.nilai {
			t.Errorf("%q: %q %q", c.filter, j, n)
		}
	}
}

// TestCariViewPolis - semua tab: kutip dibuang, halaman 20 (halaman 2 -> offset 20), filter kosong / kutip saja -> baris [] tanpa
// membaca; opsi tak sah 400 tanpa membaca; 401; halaman < 1 400; tanpa pembaca 503.
func TestCariViewPolis(t *testing.T) {
	ctx := context.Background()
	v := &viewTiruan{}
	svc := Baru(nil).DenganViewPolis(v)
	lain := inti.Pelaku{AkunID: "UJI-LAIN"}
	type hasil struct {
		total, baris int
		err          error
	}
	cari := map[string]func(s *Service, p inti.Pelaku, opsi, filter string, halaman int) hasil{
		"offer": func(s *Service, p inti.Pelaku, opsi, filter string, halaman int) hasil {
			h, err := s.CariOfferFac(ctx, p, opsi, filter, halaman)
			return hasil{h.Total, len(h.Baris), err}
		},
		"treaty": func(s *Service, p inti.Pelaku, opsi, filter string, halaman int) hasil {
			h, err := s.CariPolisTreaty(ctx, p, opsi, filter, halaman)
			return hasil{h.Total, len(h.Baris), err}
		},
		"life": func(s *Service, p inti.Pelaku, _, filter string, halaman int) hasil {
			h, err := s.CariPolisLife(ctx, p, filter, halaman)
			return hasil{h.Total, len(h.Baris), err}
		},
		"offerlife": func(s *Service, p inti.Pelaku, _, filter string, halaman int) hasil {
			h, err := s.CariOfferLife(ctx, p, filter, halaman)
			return hasil{h.Total, len(h.Baris), err}
		},
	}
	mau := map[string]string{
		"offer":     "offer|" + repository.OfferFacInsured + "|PT A%|20|20",
		"treaty":    "treaty|" + repository.TreatyInsured + "|%PT A%|20|20",
		"life":      "life|" + repository.LifeSOB + "|%PT A%|20|20",
		"offerlife": "offerlife|" + repository.OfferLifeCeding + "|%PT A%|20|20",
	}
	opsi := map[string]string{"offer": "5", "treaty": "4", "life": "", "offerlife": ""}
	for tab, f := range cari {
		v.cari = nil
		if h := f(svc, lain, opsi[tab], ` "PT 'A" `, 2); h.err != nil || h.total != 41 || h.baris != 1 || !reflect.DeepEqual(v.cari, []string{mau[tab]}) {
			t.Errorf("%s: %+v %v", tab, h, v.cari)
		}
		for _, kosong := range []string{"", `"'`, "  "} {
			if h := f(svc, lain, opsi[tab], kosong, 1); h.err != nil || h.total != 0 || h.baris != 0 || len(v.cari) != 1 {
				t.Errorf("%s kosong %q: %+v %v", tab, kosong, h, v.cari)
			}
		}
		if h := f(svc, inti.Pelaku{}, opsi[tab], "X", 1); !errors.Is(h.err, inti.ErrTanpaIdentitas) {
			t.Errorf("%s 401: %v", tab, h.err)
		}
		if h := f(svc, lain, opsi[tab], "X", 0); !errors.Is(h.err, ErrMasukanFollowing) {
			t.Errorf("%s halaman 0: %v", tab, h.err)
		}
		if h := f(Baru(nil), lain, opsi[tab], "PT A", 1); !errors.Is(h.err, ErrPolisTanpaDatabase) {
			t.Errorf("%s 503: %v", tab, h.err)
		}
		if h := f(Baru(nil), lain, opsi[tab], "", 1); h.err != nil || h.baris != 0 {
			t.Errorf("%s kosong tanpa basis data: %+v", tab, h)
		}
	}
	v.cari = nil
	for tab, o := range map[string]string{"offer": "7", "treaty": "6"} {
		if h := cari[tab](svc, lain, o, "X", 1); !errors.Is(h.err, ErrMasukanFollowing) || len(v.cari) != 0 {
			t.Errorf("%s opsi %s: %v %v", tab, o, h.err, v.cari)
		}
	}
	if h, _ := svc.CariPolisLife(ctx, lain, "", 1); h.Baris == nil || h.Ukuran != 20 || h.Halaman != 1 {
		t.Errorf("baris [] bukan null: %+v", h)
	}
}

// TestBerkasPolisTreaty - tombol View tab Treaty (work owner 11-10-2026): nomor dipangkas, PRODKE 0 / >= 1 -> modul, tanpa berkas 404,
// kosong 400 tanpa membaca, 401, 503.
func TestBerkasPolisTreaty(t *testing.T) {
	ctx := context.Background()
	v := &viewTiruan{}
	svc := Baru(nil).DenganViewPolis(v)
	lain := inti.Pelaku{AkunID: "UJI-LAIN"}
	for no, mau := range map[string]models.BerkasPolisTreaty{" RNM-T0 ": {Modul: "nbtreatyin", Kasus: "NBT-1"},
		"RNM-T2": {Modul: "edmtreatyin", Kasus: "EDMT-2"}} {
		if b, err := svc.BerkasPolisTreaty(ctx, lain, no); err != nil || b != mau {
			t.Errorf("%q: %+v %v", no, b, err)
		}
	}
	if !reflect.DeepEqual(v.cari, []string{"berkas|RNM-T0", "berkas|RNM-T2"}) && !reflect.DeepEqual(v.cari, []string{"berkas|RNM-T2", "berkas|RNM-T0"}) {
		t.Errorf("dibaca %v", v.cari)
	}
	if _, err := svc.BerkasPolisTreaty(ctx, lain, "RNM-LAMA"); !errors.Is(err, ErrBerkasTreatyTidakAda) {
		t.Errorf("404: %v", err)
	}
	v.cari = nil
	if _, err := svc.BerkasPolisTreaty(ctx, lain, "  "); !errors.Is(err, ErrMasukanPolis) || len(v.cari) != 0 {
		t.Errorf("400: %v %v", err, v.cari)
	}
	if _, err := svc.BerkasPolisTreaty(ctx, inti.Pelaku{}, "RNM-T0"); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("401: %v", err)
	}
	if _, err := Baru(nil).BerkasPolisTreaty(ctx, lain, "RNM-T0"); !errors.Is(err, ErrPolisTanpaDatabase) {
		t.Errorf("503: %v", err)
	}
	for prodke, mau := range map[int64]string{0: "nbtreatyin", 1: "edmtreatyin", 7: "edmtreatyin", -1: "nbtreatyin"} {
		if m := models.ModulBerkasPolisTreaty(prodke); m != mau {
			t.Errorf("prodke %d: %q", prodke, m)
		}
	}
}
