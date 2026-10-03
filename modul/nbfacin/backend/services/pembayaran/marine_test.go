package pembayaran

import (
	"errors"
	"testing"

	"nusantarare/inti/backend/uang"
	"nusantarare/inti/backend/utils"
)

type kasus map[string]string

func (k kasus) Nilai(j string) (string, bool) { v, ada := k[j]; return v, ada }

const (
	jStatusNB  = "pyWorkPage.Quotation.StatusBusiness"
	jStatusEDM = "pyWorkPage.OfferFacIn.QuotationData.StatusBusiness"
	jBisnis    = ".OfferFacIn.QuotationData.BusinessType"
)

var nbMarine = kasus{jStatusNB: "1", jStatusEDM: "1", jBisnis: "MarineCargo"}

func kargo(cov ...Coverage) [][]Coverage { return [][]Coverage{cov} }

// TestMarineNB - PremiPaymentMarine langkah 1.1: @sum premi dan diskon seluruh
// coverage seluruh kargo.
func TestMarineNB(t *testing.T) {
	daftar := append(kargo(Coverage{"IDR", "1000.0000", "10.0000"}, Coverage{"IDR", "0.5", "0"}),
		[]Coverage{{"IDR", "2.25", "1.5"}})
	got, err := Marine(nbMarine, daftar)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Berlaku || utils.FormatDecimal(got.Premium.Amount) != "1002.7500" || utils.FormatDecimal(got.Diskon.Amount) != "11.5000" || got.Premium.Currency != "IDR" {
		t.Errorf("dapat %+v", got)
	}
}

// TestMarineGerbang - langkah 45 IsMarineCargo; cabang IsEDM belum diport; tanpa
// cabang yang terbuka tidak ditiru diam-diam.
func TestMarineGerbang(t *testing.T) {
	cov := kargo(Coverage{"IDR", "1", "0"})
	if got, err := Marine(kasus{jStatusNB: "1", jBisnis: "Fire"}, cov); err != nil || got.Berlaku {
		t.Errorf("bukan marine: %+v (%v)", got, err)
	}
	// IsNB dan IsEDM membaca halaman berbeda: keduanya bisa benar.
	for nama, k := range map[string]kasus{
		"EDM":        {jStatusNB: "3", jStatusEDM: "3", jBisnis: "MarineCargo"},
		"NB dan EDM": {jStatusNB: "1", jStatusEDM: "3", jBisnis: "MarineCargo"},
	} {
		if _, err := Marine(k, cov); !errors.Is(err, ErrCabangEDM) {
			t.Errorf("%s: galat %v", nama, err)
		}
	}
	if _, err := Marine(kasus{jStatusNB: "2", jStatusEDM: "2", jBisnis: "MarineCargo"}, cov); !errors.Is(err, ErrTanpaCabang) {
		t.Errorf("renewal-tanpa-cabang: galat %v", err)
	}
}

// TestMarineDitolak - masukan kosong, tak terbaca, mata uang kosong (A46),
// mata uang campuran (A45), tanpa coverage - masing-masing dengan galat yang tepat.
func TestMarineDitolak(t *testing.T) {
	for nama, u := range map[string]struct {
		daftar [][]Coverage
		mau    error
	}{
		"premi kosong":            {kargo(Coverage{"IDR", "", "0"}), ErrMasukan},
		"diskon kosong":           {kargo(Coverage{"IDR", "1", ""}), ErrMasukan},
		"tak terbaca":             {kargo(Coverage{"IDR", "1,5", "0"}), ErrMasukan},
		"mata uang kosong":        {kargo(Coverage{"", "1", "0"}), ErrMasukan},
		"mata uang kosong berdua": {kargo(Coverage{"", "1", "0"}, Coverage{"", "2", "0"}), ErrMasukan},
		"mata uang campur":        {kargo(Coverage{"IDR", "1", "0"}, Coverage{"USD", "1", "0"}), uang.ErrMataUangBerbeda},
		"tanpa coverage":          {[][]Coverage{{}}, ErrMasukan},
	} {
		if _, err := Marine(nbMarine, u.daftar); !errors.Is(err, u.mau) {
			t.Errorf("%s: galat %v, mau %v", nama, err, u.mau)
		}
	}
}
