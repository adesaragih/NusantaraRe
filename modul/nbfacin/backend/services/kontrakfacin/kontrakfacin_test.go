package kontrakfacin

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/services/acceptance"
)

type kasus map[string]string

func (k kasus) Nilai(j string) (string, bool) { v, ada := k[j]; return v, ada }

func d(t *testing.T, s string) *apd.Decimal {
	t.Helper()
	v, err := utils.ParseDecimal(s)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

// Data sintetis: satu team group, dua jabatan bentuk A; dua baris bentuk B.
func tangga(t *testing.T) Tangga {
	return Tangga{
		BentukA: acceptance.TabelLimit{acceptance.TabelProperty: {
			{Jabatan: "SENIORUW", TeamGroup: "1", LimitBottom: d(t, "1"), LimitBottom2: d(t, "1")},
			{Jabatan: "DIREKTURTEKNIK", TeamGroup: "1", LimitBottom: d(t, "10"), LimitBottom2: d(t, "10")},
		}},
		BentukB: []acceptance.BarisFinancial{
			{Jabatan: "KADIV KEUANGAN", LimitBond: d(t, "0"), LimitCreditCL: d(t, "0"), LimitCreditNCL: d(t, "0")},
		},
	}
}

var dasar = kasus{
	"pyWorkPage.OfferFacIn.TotalTSINusaRe": "100", "pyWorkPage.OfferFacIn.TotalTSITopRisk": "0",
	"pyWorkPage.OfferFacIn.QuotationData.StatusBusiness": "2", "pyWorkPage.OfferFacIn.QuotationData.TeamGroup": "1",
	"pyWorkPage.OfferFacIn.QuotationData.BusinessType": "Fire", "pyWorkPage.OfferFacIn.QuotationData.BusinessOldId": "22",
	"pyWorkPage.OfferFacIn.IsPreferredRisk": "Preferred Risk",
}

// TestLangkahBentukA - kontrak meneruskan kasus bentuk A ke acceptance.Next apa adanya.
func TestLangkahBentukA(t *testing.T) {
	var tg kontrak.TanggaAkseptasiFacIn = tangga(t)
	got, err := tg.Langkah(dasar, kontrak.PenggunaFacIn{Jabatan: "SENIORUW"})
	mau := kontrak.TransisiFacIn{JabatanTujuan: "DIREKTURTEKNIK", Antrean: "ReasFacInTechnicalDirector", PositionNoteDitulis: true}
	if err != nil || got != mau {
		t.Fatalf("dapat %+v (%v), mau %+v", got, err, mau)
	}
}

// TestLangkahBentukB - kasus Bond dialihkan ke acceptance.NextFinancial.
func TestLangkahBentukB(t *testing.T) {
	bond := kasus{}
	for k, v := range dasar {
		bond[k] = v
	}
	bond["pyWorkPage.OfferFacIn.QuotationData.BusinessOldId"] = "36"
	bond["pyWorkPage.OfferFacIn.QuotationData.BusinessType"] = "Bonding"
	bond["pyWorkPage.PositionNote"] = "ReasFacInUnderwritingFinancial"
	got, err := tangga(t).Langkah(bond, kontrak.PenggunaFacIn{})
	mau := kontrak.TransisiFacIn{JabatanTujuan: "KADIVFINANCIAL", Antrean: "ReasFacInFinDivHead", PositionNoteDitulis: true}
	if err != nil || got != mau {
		t.Fatalf("dapat %+v (%v), mau %+v", got, err, mau)
	}
}

// TestPredikat - kontrak predikat = rules.Eval registry NB.
func TestPredikat(t *testing.T) {
	var p kontrak.PenilaiPredikatFacIn = Predikat{}
	if buka, err := p.Eval("IsFire", dasar); err != nil || !buka {
		t.Fatalf("IsFire: %v (%v)", buka, err)
	}
}

// TestPremi - kontrak premi = premium.Calculate / AsalRumus / LiniDariPredikat.
func TestPremi(t *testing.T) {
	var m kontrak.MesinPremiFacIn = Premi{}
	in := kontrak.MasukanPremiFacIn{LiniBisnis: "PA", CalculateMethod: "3", TSI: "2787500000",
		Rate: "13.47311827957", ProRatePercent: "400.27397260273972602700"}
	got, err := m.Hitung(in)
	if err != nil || utils.FormatDecimal(got.Amount) != "37556317.2043" {
		t.Fatalf("Hitung: %v (%v)", got, err)
	}
	if a := m.AsalRumus(in); a != "CalculatePremiPA_FacIn langkah 6 L1003" {
		t.Fatalf("AsalRumus %q", a)
	}
	k := kasus{"pyWorkPage.Quotation.BusinessCode": "10028"}
	for j, v := range dasar {
		k[j] = v
	}
	h, err := m.LiniDariPredikat(k)
	if err != nil || len(h.Lini) != 1 || h.Lini[0] != (kontrak.LiniFacIn{Lini: "FIRE", PembagiRate: 1000, SimbolSatuan: "‰"}) {
		t.Fatalf("LiniDariPredikat %+v (%v)", h, err)
	}
}

// TestTanggaTanpaTabelDitolak - tabel limit yang belum tersambung tidak boleh
// menjadi "tidak ada kandidat → selesai" diam-diam.
func TestTanggaTanpaTabelDitolak(t *testing.T) {
	if _, err := (Tangga{}).Langkah(dasar, kontrak.PenggunaFacIn{Jabatan: "SENIORUW"}); !errors.Is(err, ErrTabelLimitKosong) {
		t.Fatalf("galat %v, mau ErrTabelLimitKosong", err)
	}
}

// TestMasukanSemuaMedan - setiap medan premium.Input punya kembaran bernama sama
// di kontrak.MasukanPremiFacIn, dan masukan() menyalinnya: medan baru di salah
// satu sisi tanpa pemetaan = merah, bukan nilai nol diam-diam.
func TestMasukanSemuaMedan(t *testing.T) {
	var k kontrak.MasukanPremiFacIn
	kv := reflect.ValueOf(&k).Elem()
	for i := 0; i < kv.NumField(); i++ {
		switch f := kv.Field(i); f.Kind() {
		case reflect.String:
			f.SetString("isi-" + kv.Type().Field(i).Name)
		case reflect.Bool:
			f.SetBool(true)
		default:
			t.Fatalf("%s: jenis %s belum ditangani tes", kv.Type().Field(i).Name, f.Kind())
		}
	}
	in := reflect.ValueOf(masukan(k))
	if in.NumField() != kv.NumField() {
		t.Errorf("premium.Input %d medan, kontrak %d", in.NumField(), kv.NumField())
	}
	for i := 0; i < in.NumField(); i++ {
		nama := in.Type().Field(i).Name
		sumber := kv.FieldByName(nama)
		if !sumber.IsValid() {
			t.Errorf("%s tidak ada di kontrak", nama)
			continue
		}
		if fmt.Sprint(in.Field(i).Interface()) != fmt.Sprint(sumber.Interface()) {
			t.Errorf("%s: %v, mau %v", nama, in.Field(i).Interface(), sumber.Interface())
		}
	}
}
