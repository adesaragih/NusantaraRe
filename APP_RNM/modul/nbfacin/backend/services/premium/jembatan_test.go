package premium

import (
	"reflect"
	"strings"
	"testing"
)

// kasusPeta - rules.Kasus dari peta jalur → nilai, jalur persis seperti di rule.
type kasusPeta map[string]string

func (k kasusPeta) Nilai(jalur string) (string, bool) {
	v, ada := k[jalur]
	return v, ada
}

// denganKode - kasus berkode bisnis terisi, seperti data nyata. Beberapa
// predikat di rantai IsAneka (mis. IsYieldShortfall) membandingkan
// pyWorkPage.Quotation.BusinessCode dengan LITERAL ANGKA, sehingga kode kosong
// ditolak (butir 20). "10028" adalah kode bisnis kasus PA nyata (bukan PII).
func denganKode(k kasusPeta) kasusPeta {
	hasil := kasusPeta{"pyWorkPage.Quotation.BusinessCode": "10028"}
	for j, v := range k {
		hasil[j] = v
	}
	return hasil
}

// TestLiniDariPredikat - tiket NB-10: lini bisnis diturunkan dari gerbang
// predikat yang sungguh dipakai sistem lama untuk memilih rumus premi
// (CountGrossPremi_Act langkah 5.1-5.4, CountGPWMarinePAMbu_Act langkah 1.1-1.3).
func TestLiniDariPredikat(t *testing.T) {
	for _, u := range []struct {
		nama  string
		kasus kasusPeta
		mau   []LiniBisnis
	}{
		{"Fire", kasusPeta{"pyWorkPage.OfferFacIn.QuotationData.BusinessType": "Fire"}, []LiniBisnis{LiniFire}},
		{"PA", kasusPeta{"pyWorkPage.Quotation.BusinessType": "PA"}, []LiniBisnis{LiniPA}},
		{"MBU", kasusPeta{"pyWorkPage.OfferFacIn.QuotationData.BusinessType": "MBUCar"}, []LiniBisnis{LiniMBU}},
		{"Golf", kasusPeta{"pyWorkPage.Quotation.BusinessType": "GolfInsurance"}, []LiniBisnis{LiniGolf}},
		{"Marine Cargo", kasusPeta{"pyWorkPage.Quotation.BusinessType": "MarineCargo"}, []LiniBisnis{LiniMarineCargo}},
		{"Aneka", kasusPeta{"pyWorkPage.Quotation.BusinessType": "Aneka"}, []LiniBisnis{LiniAneka}},
		// Bonding tidak punya gerbang sendiri: IsAneka merujuk
		// IsBondingAndCustomBonds, jadi Bonding dihitung lewat jalur ANEKA.
		{"Bonding lewat jalur Aneka", kasusPeta{"pyWorkPage.Quotation.BusinessType": "Bonding"}, []LiniBisnis{LiniAneka}},
		// Spasi di ujung dipangkas di batas input (kandidat perbaikan).
		{"spasi dipangkas", kasusPeta{"pyWorkPage.Quotation.BusinessType": " PA "}, []LiniBisnis{LiniPA}},
	} {
		got, err := LiniDariPredikat(denganKode(u.kasus))
		if err != nil {
			t.Errorf("%s: %v", u.nama, err)
			continue
		}
		var lini []LiniBisnis
		for _, l := range got.Lini {
			lini = append(lini, l.Lini)
			if l.Satuan != SatuanRate(l.Lini) {
				t.Errorf("%s: satuan %s tidak dari resolver", u.nama, l.Lini)
			}
		}
		if !reflect.DeepEqual(lini, u.mau) {
			t.Errorf("%s: lini %v, mau %v", u.nama, lini, u.mau)
		}
	}
}

// TestLiniDariPredikatJalurTakSinkron - jalur baca ganda dipertahankan; bila
// nilainya berbeda, keduanya dihitung seperti sistem lama DAN perbedaannya
// tercatat, bukan dipilih diam-diam.
func TestLiniDariPredikatJalurTakSinkron(t *testing.T) {
	got, err := LiniDariPredikat(denganKode(kasusPeta{
		"pyWorkPage.Quotation.BusinessType":                "PA",
		"pyWorkPage.OfferFacIn.QuotationData.BusinessType": "Fire",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Lini) != 2 || got.Lini[0].Lini != LiniFire || got.Lini[1].Lini != LiniPA {
		t.Errorf("lini %v, mau FIRE lalu PA (urutan langkah)", got.Lini)
	}
	if len(got.Peringatan) == 0 || !strings.Contains(strings.Join(got.Peringatan, " "), "berbeda") {
		t.Errorf("ketidaksinkronan jalur tidak tercatat: %v", got.Peringatan)
	}
}

// TestLiniDariPredikatTakDikenali - tidak satu gerbang pun terbuka → gagal keras
// lewat resolver, bukan lini bawaan.
func TestLiniDariPredikatTakDikenali(t *testing.T) {
	for _, k := range []kasusPeta{
		{"pyWorkPage.Quotation.BusinessType": "Travel"},
		{},
	} {
		harusPanic(t, "lini tak dikenali "+k["pyWorkPage.Quotation.BusinessType"], func() { _, _ = LiniDariPredikat(denganKode(k)) })
	}
}

// TestLiniDariPredikatTakDikenaliWalauNamaLini - BusinessType yang kebetulan
// berisi NAMA lini resolver (mis. "BONDING", "Layering") tetapi tidak membuka
// gerbang mana pun tetap gagal keras, bukan daftar kosong tanpa galat. Temuan
// review: penanda gagal-keras lama memanggil SatuanRate dengan nilai itu, yang
// justru lolos untuk nama yang dikenal.
func TestLiniDariPredikatTakDikenaliWalauNamaLini(t *testing.T) {
	for _, v := range []string{"BONDING", "Layering", "MARINE CARGO"} {
		harusPanic(t, "BusinessType = "+v, func() {
			_, _ = LiniDariPredikat(denganKode(kasusPeta{"pyWorkPage.Quotation.BusinessType": v}))
		})
	}
}

// TestLiniDariPredikatPangkasDicatat - spasi yang dipangkas di batas input
// wajib tercatat sebagai peringatan (kandidat perbaikan), bukan hilang diam-diam.
func TestLiniDariPredikatPangkasDicatat(t *testing.T) {
	got, err := LiniDariPredikat(denganKode(kasusPeta{"pyWorkPage.Quotation.BusinessType": " PA "}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(got.Peringatan, " "), "spasi dipangkas") {
		t.Errorf("pemangkasan tidak tercatat: %v", got.Peringatan)
	}
}
