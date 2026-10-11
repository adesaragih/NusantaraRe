package models_test

// Baris OS_AKSEPTASI_KLAIM Download Claim Face Sheet (`SaveCFS_ACT` langkah 9-19): kunci DATA_JSON per lini = daftar
// Property-Set halaman TempOSAkseptasi XML (+ pxCreateOperator langkah 16, pxObjClass halaman). Perintah work owner
// 11-10-2026: "cek apakah insert ke os akseptasi klaim sudah sesuai dengan SaveCFS_ACT.xml" -> Marine / MBU dilengkapi.

import (
	"encoding/json"
	"sort"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/claimfacin/backend/models"
)

func halamanOSUji(lini string) *models.Halaman {
	h := models.HalamanBaru()
	h.Setel(models.OQ+"BusinessType", lini)
	h.Setel(models.CD+"CauseOfLoss", "UJI SEBAB")
	h.Setel(models.CD+"CauseOfLossID", "UJI-S1")
	h.Setel(models.CD+"NoClaim", "UJI-K.001")
	h.Setel(models.AwalanPolis+".PercentShare", "25")
	h.Setel(models.JalurNoPolis, "UJI-RNM-F.001")
	h.SetelDaftar(models.DaftarObjek, []models.Baris{{"ObjectID": "UJI-O1", "ObjectName": "UJI NAMA", "ObjectLocation": "UJI KEMAS",
		"ObjectJob": "UJI BARANG", "ObjectSurveyLocation": "UJI ANGKUT", "Brand": "UJI-B1", "BrandName": "UJI MEREK",
		"Model": "UJI-M1", "ModelName": "UJI MODEL", "Type": "UJI-T1", "TypeName": "UJI TIPE", "LicensePlate": "UJI 1234"}})
	h.SetelDaftar(models.DaftarItem(1), []models.Baris{{"CoverageID": "UJI-C1", "CoverageNote": "UJI COVER",
		"ObjectItemName": "UJI ITEM", "IsFacretro": "1"}})
	h.SetelDaftar(models.DaftarDiItem(1, 1, "CoverageList"), []models.Baris{{"CoverageNote": "UJI ICC A"}})
	h.SetelDaftar(models.DaftarDiItem(1, 1, models.AnakEstimasi), []models.Baris{{"EstimationValue": "100",
		"GrossEstimationPct": "400", "GrossEstimationPctMBU": "450", "Deductible": "5", "Currency": "IDR", "CurrencyID": "10026"}})
	return h
}

func TestBarisOSEstimasiIkutSaveCFS(t *testing.T) {
	k := &models.Konteks{Pelaku: "UJI-ADMIN", Sekarang: time.Date(2026, 10, 11, 9, 0, 0, 0, models.Jakarta)}
	umum := "CauseOfLoss Currency CurrencyID DeductibleValue EstimationDate GrossValue ObjectID PersenRNM Type Value"
	tambah := "pxCreateOperator pxObjClass"
	for _, c := range []struct {
		lini, kunci string
		nilai       map[string]string
	}{
		{"Fire", umum + " CauseOfLossID CoverageID CoverageName ObjectName " + tambah,
			map[string]string{"ObjectName": "UJI ITEM", "CoverageName": "UJI COVER", "GrossValue": "400"}},
		{"MarineCargo", umum + " TradingName PackingName Goodsname ConveyanceName CoverageName " + tambah,
			map[string]string{"TradingName": "UJI NAMA", "PackingName": "UJI KEMAS", "Goodsname": "UJI BARANG",
				"ConveyanceName": "UJI ANGKUT", "CoverageName": "UJI ICC A"}},
		{"MBUCar", umum + " BrandName BrandID ModelName Model TypeName LicensePlate TypeID CauseOfLossID CoverageID " +
			"CoverageName NoClaim " + tambah,
			map[string]string{"BrandName": "UJI MEREK", "BrandID": "UJI-B1", "ModelName": "UJI MODEL", "Model": "UJI-M1",
				"TypeName": "UJI TIPE", "LicensePlate": "UJI 1234", "TypeID": "UJI-T1", "GrossValue": "450", "Type": "0",
				"NoClaim": "UJI-K.001"}},
	} {
		b := models.BarisOSEstimasi(k, halamanOSUji(c.lini), 1, 1, 1, "CLM-UJI1")
		var isi map[string]string
		if err := json.Unmarshal([]byte(b.DataJSON), &isi); err != nil {
			t.Fatalf("%s: DATA_JSON %v", c.lini, err)
		}
		var dapat []string
		for k := range isi {
			dapat = append(dapat, k)
		}
		mau := strings.Fields(c.kunci)
		sort.Strings(dapat)
		sort.Strings(mau)
		if strings.Join(dapat, " ") != strings.Join(mau, " ") {
			t.Errorf("%s: kunci DATA_JSON\n dapat %v\n mau   %v", c.lini, dapat, mau)
		}
		for k, v := range c.nilai {
			if isi[k] != v {
				t.Errorf("%s: %s = %q, mau %q", c.lini, k, isi[k], v)
			}
		}
		if b.StsReject != models.StsOSOutstanding || b.StsDLA != models.StsDLARetro || b.NoClaim != "UJI-K.001" ||
			b.NoPolis != "UJI-RNM-F.001" || b.CaseID != "CLM-UJI1" {
			t.Errorf("%s: kolom baris OS %+v", c.lini, b)
		}
	}
}
