package services

import (
	"reflect"
	"testing"

	"nusantarare/modul/treatyin/backend/models"
)

// Total kosong dihitung dari rinciannya saat kontrak dibuka; total tersimpan
// menang (laporan pemakai 9 Oktober 2026, kontrak HOSPITAL 2025 QS 181M TRT).
func TestLengkapiTotalPenampung(t *testing.T) {
	detail := map[string]any{
		"RNMShareList":      []map[string]any{{"Currency": "IDR", "Value": "56250000"}},
		"RNMSpreadedList":   []map[string]any{{"Currency": "IDR", "Value": "22500000"}},
		"RNMSpreadedListRI": []map[string]any{{"Currency": "IDR", "Value": "33750000"}},
	}
	k := models.KontrakWarisan{
		LimitsPohon: []map[string]any{{"TreatyType": "QUOTA SHARE", "Detail": []map[string]any{detail}}},
		// Total Share tersimpan — TIDAK ditimpa.
		PenampungLarik: map[string][]map[string]any{"TotalShareRnmProp": {{"Currency": "IDR", "Value": "999"}}},
		Egnpi:          []models.BarisEgnpiWarisan{{MataUang: "IDR", Jumlah: "100"}, {MataUang: "USD", Jumlah: "2"}, {MataUang: "IDR", Jumlah: "50"}},
		Angsuran:       []models.BarisAngsuranWarisan{{MataUang: "IDR", Jumlah: "10.5"}, {MataUang: "IDR", Jumlah: "4.5"}},
	}
	lengkapiTotalPenampung(&k)
	ingin := map[string][]map[string]any{
		"TotalShareRnmProp":      {{"Currency": "IDR", "Value": "999"}},
		"TotalSpreadedRnmProp":   {{"Currency": "IDR", "Value": "22500000"}},
		"TotalSpreadedRnmRIProp": {{"Currency": "IDR", "Value": "33750000"}},
		"TotalEgnpiAmountNP":     {{"Currency": "IDR", "Value": "150"}, {"Currency": "USD", "Value": "2"}},
		"TotalInstallmentNP":     {{"Currency": "IDR", "Value": "15"}},
	}
	for n, v := range ingin {
		if !reflect.DeepEqual(k.PenampungLarik[n], v) {
			t.Errorf("%s = %v, ingin %v", n, k.PenampungLarik[n], v)
		}
	}
}

// Kontrak tanpa rincian: nol larik total yang lahir.
func TestLengkapiTotalPenampungTanpaRincian(t *testing.T) {
	k := models.KontrakWarisan{}
	lengkapiTotalPenampung(&k)
	if len(k.PenampungLarik) != 0 {
		t.Errorf("total lahir dari kekosongan: %v", k.PenampungLarik)
	}
}
