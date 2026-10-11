package services

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// lokasiFacInOfferAneka - satu lokasi: okupasi OKU1 (objek A ber-SerialNo, satu baris spreading), okupasi OKU2 (objek B,
// coverage USD dengan dua baris spreading), objek Golf G (satu baris).
func lokasiFacInOfferAneka() []models.LokasiAneka {
	baris := func(jenis string) models.BarisSpreading {
		return models.BarisSpreading{TreatyType: jenis, SharePercentage: d("50"), TSISpreaded: d("10"), PremiumSpreaded: d("1")}
	}
	return []models.LokasiAneka{{ObjectNo: "1", Okupasi: []models.OkupasiAneka{
		{OccupationName: "OKU1", Objek: []models.ObjekAneka{{ObjectName: "A", Currency: "IDR", SerialNo: "SN-1", Coverages: []models.CoverageAneka{
			{Coverage: "C1", CoverageNote: "N1", Rate: d("0.25"), Spreading: []models.BarisSpreading{baris("10007")}}}}}},
		{OccupationName: "OKU2", Objek: []models.ObjekAneka{{ObjectName: "B", Currency: "IDR", Coverages: []models.CoverageAneka{
			{Coverage: "C2", Currency: "USD", Spreading: []models.BarisSpreading{baris("Q1"), baris("S1")}}}}}}},
		ObjekGolf: []models.ObjekAneka{{ObjectName: "G", Currency: "EUR", SerialNo: "SN-G", Coverages: []models.CoverageAneka{
			{Coverage: "CG", Spreading: []models.BarisSpreading{baris("10007")}}}}}}}
}

func teksFacInOffer(b models.BarisFacInOffer) string {
	return strings.Join([]string{b.OccupationCode, b.ObjectItem, b.CurrID, b.NoSerial, b.Coverage, b.JnReas, b.TopRisk, b.ObjectName}, "|")
}

// TestBarisFacInOfferAneka - SaveOfferProduction_Act cabang 13 / 14 / 15: okupasi bertahap, ObjectName / Currency objek, SerialNo
// (Aneka saja), MBD = Currency coverage, TOP_RISK kosong, alamat dari baris lokasi FIRE yang sama; Bonding satu cabang (K73-35).
func TestBarisFacInOfferAneka(t *testing.T) {
	dasar := models.BarisFacInOffer{IDPega: "UJI-NB-1", BusinessCode: "10048"}
	objek := []models.ObjekFire{{ObjectNo: "1", ObjectName: "LOKASI UJI", Province: "UJI PROV", BuildingNo: "7"}}
	for nama, u := range map[string]struct {
		l   liniFacInOffer
		mau []string
	}{
		"Aneka": {liniFacInOffer{aneka: true}, []string{"OKU1|A|IDR|SN-1|C1|10007||LOKASI UJI", "OKU2 || OKU1|B|IDR||C2|Q1||LOKASI UJI",
			"OKU2 || OKU1|B|IDR||C2|S1||LOKASI UJI"}},
		"MBD": {liniFacInOffer{aneka: true, mbd: true}, []string{"OKU1|A|IDR|SN-1|C1|10007||LOKASI UJI",
			"OKU2 || OKU1|B|USD||C2|Q1||LOKASI UJI", "OKU2 || OKU1|B|USD||C2|S1||LOKASI UJI"}},
		"Bonding": {liniFacInOffer{bonding: true, mbd: true}, []string{"OKU1|A|IDR||C1|10007||LOKASI UJI",
			"OKU2 || OKU1|B|IDR||C2|Q1||LOKASI UJI", "OKU2 || OKU1|B|IDR||C2|S1||LOKASI UJI"}},
		"Golf": {liniFacInOffer{golf: true}, []string{"|G|EUR||CG|10007||LOKASI UJI"}},
	} {
		b, err := barisFacInOfferAneka(dasar, objek, lokasiFacInOfferAneka(), u.l)
		if err != nil {
			t.Fatal(err)
		}
		var got []string
		for _, x := range b {
			got = append(got, teksFacInOffer(x))
		}
		if !reflect.DeepEqual(got, u.mau) {
			t.Errorf("%s:\n%q", nama, got)
		}
		if b[0].IDPega != "UJI-NB-1" || b[0].Province != "UJI PROV" || b[0].BuildingNo != "NO. 7" || t8(b[0].TSI) != "10" {
			t.Errorf("%s dasar / alamat: %+v", nama, b[0])
		}
	}
	if _, err := barisFacInOfferAneka(dasar, nil, lokasiFacInOfferAneka(), liniFacInOffer{aneka: true}); !errors.Is(err, ErrBindingTidakSejalan) {
		t.Errorf("tidak sejalan: %v", err)
	}
}

// TestFacInOfferKasusAneka - case Bonding (IsAneka DAN IsBonding): satu baris per spreading lewat cabang Bonding (K73-35);
// case bukan keluarga Aneka -> cabang FIRE (false).
func TestFacInOfferKasusAneka(t *testing.T) {
	ctx := context.Background()
	svc := layananSubmitAneka(submitUji(), "36", lokasiFacInOfferAneka(), nil)
	b, aneka, err := svc.barisFacInOfferKasusAneka(ctx, "UJI-NB-1", models.BarisFacInOffer{}, []models.ObjekFire{{ObjectNo: "1"}})
	if err != nil || !aneka || len(b) != 3 || b[0].NoSerial != "" {
		t.Errorf("Bonding: %v %v %+v", err, aneka, b)
	}
	if _, aneka, _ := layananSubmit(submitUji()).barisFacInOfferKasusAneka(ctx, "UJI-NB-1", models.BarisFacInOffer{}, nil); aneka {
		t.Error("FIRE")
	}
}
