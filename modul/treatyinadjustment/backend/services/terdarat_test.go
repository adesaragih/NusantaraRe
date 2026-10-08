package services_test

// `Penyesuaian.Terdarat` — audit 8 Oktober 2026: penyesuaian yang isi
// tabnya TIDAK ada di pendaratan dikunci layar dari suntingan.

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/treatyinadjustment/backend/models"
	"nusantarare/modul/treatyinadjustment/backend/services"
)

func TestDrafMembawaTandaTerdaratMaster(t *testing.T) {
	m := models.MasukanDraf{ID: "1001001", InternalType: "1", MaterialType: "1"}
	for _, terdarat := range []bool{true, false} {
		dok := models.SisiPenyesuaian{Medan: map[string]string{"ID": "1001001"}, Larik: map[string][]map[string]string{}, Terdarat: terdarat}
		p := services.SusunDraf(dok, m, false, "UJI", time.Now())
		if p.Terdarat != terdarat {
			t.Errorf("dok terdarat %v → draf %v", terdarat, p.Terdarat)
		}
	}
}

// Tanda itu sampai ke layar di tingkat penyesuaian, BUKAN per sisi — sisi
// adalah halaman dokumen, dan `terdarat` bukan properti Pega.
func TestTerdaratDiJSONPenyesuaianBukanSisi(t *testing.T) {
	b, err := json.Marshal(models.Penyesuaian{ID: "X", Baru: models.SisiPenyesuaian{Terdarat: true}})
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if !strings.Contains(s, `"terdarat":false`) || strings.Count(s, "erdarat") != 1 {
		t.Errorf("JSON %s", s)
	}
}
