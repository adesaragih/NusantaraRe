package repository

import (
	"reflect"
	"testing"

	"nusantarare/modul/treatyinadjustment/backend/models"
)

// Total kosong dihitung dari rinciannya; total tersimpan menang (laporan
// pemakai 9 Oktober 2026).
func TestLengkapiTotal(t *testing.T) {
	s := models.SisiPenyesuaian{
		Larik: map[string][]map[string]string{
			"TotalShareRnmProp": {{"Currency": "IDR", "Value": "999"}},
			"EGNPI":             {{"Currency": "IDR", "Amount": "100"}, {"Currency": "USD", "Amount": "2"}},
			"Retention":         {{"Currency": "IDR", "Amount": "3500000000"}},
		},
		Pohon: map[string][]map[string]any{
			"Limits": {{"Detail": []map[string]any{{
				"RNMShareList":      []map[string]any{{"Currency": "IDR", "Value": "56250000"}},
				"RNMSpreadedList":   []map[string]any{{"Currency": "IDR", "Value": "22500000"}},
				"RNMSpreadedListRI": []map[string]any{{"Currency": "IDR", "Value": "33750000"}},
			}}}},
			"Installment": {
				{"Currency": "IDR", "AmountTotal": "40852100.79"},
				{"Currency": "USD", "InstallmentList": []map[string]any{{"Amount": "10"}, {"Amount": "5"}}},
			},
		},
	}
	lengkapiTotal(&s)
	ingin := map[string][]map[string]string{
		"TotalShareRnmProp":      {{"Currency": "IDR", "Value": "999"}},
		"TotalSpreadedRnmProp":   {{"Currency": "IDR", "Value": "22500000"}},
		"TotalSpreadedRnmRIProp": {{"Currency": "IDR", "Value": "33750000"}},
		"TotalEgnpiAmountNP":     {{"Currency": "IDR", "Value": "100"}, {"Currency": "USD", "Value": "2"}},
		"TotalRetentionAmountNP": {{"Currency": "IDR", "Value": "3500000000"}},
		"TotalInstallmentNP":     {{"Currency": "IDR", "Value": "40852100.79"}, {"Currency": "USD", "Value": "15"}},
	}
	for n, v := range ingin {
		if !reflect.DeepEqual(s.Larik[n], v) {
			t.Errorf("%s = %v, ingin %v", n, s.Larik[n], v)
		}
	}
}
