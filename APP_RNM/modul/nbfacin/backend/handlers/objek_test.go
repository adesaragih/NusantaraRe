package handlers

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

type objekTiruan struct{ d *[]models.ObjekFire }

func (o objekTiruan) BacaObjek(context.Context, string) ([]models.ObjekFire, error) {
	return append([]models.ObjekFire{}, *o.d...), nil
}

func (o objekTiruan) GantiObjek(_ context.Context, _ *db.Tx, _ string, baris []models.ObjekFire) error {
	*o.d = baris
	return nil
}

// TestObjek - GET/PUT /api/nbfacin/kasus/{caseId}/objek (tiket 35): 23 kunci persis kontrak
// ObjekFire, boolean JSON, `baris` larik walau kosong; PUT 200 baca ulang, 400/401/503.
func TestObjek(t *testing.T) {
	var d []models.ObjekFire
	svc := services.Baru(nil).DenganObjek(objekTiruan{&d}).DenganTransaksi(tanpaOracle)
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/kasus/UJI-NB-1/objek", "", ""); kode != 200 || isi != `{"baris":[]}` {
		t.Fatalf("kosong: %d %s", kode, isi)
	}
	badan := `{"baris":[{"objectNo":"1","objectType":"UJI TIPE","isTopRisk":true,"isMaterialDamage":false,"numberOfFloor":"3","otherType":"UJI LAIN"}]}`
	kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek", badan, "UJI-USER")
	var j struct {
		Baris []map[string]any `json:"baris"`
	}
	kunciMau := []string{"objectNo", "objectType", "objectName", "isMaterialDamage", "isTopRisk", "roadType", "roadName",
		"buildingNo", "zipCode", "country", "riskLocation", "territory", "city", "district", "province", "riskAddressId",
		"numberOfFloor", "roofType", "wallType", "floorType", "partitionType", "supportWallType", "otherType"}
	if err := json.Unmarshal([]byte(isi), &j); err == nil && len(j.Baris) == 1 {
		for _, k := range kunciMau {
			if _, ada := j.Baris[0][k]; !ada {
				t.Errorf("kunci %q tidak ada (kontrak ObjekFire)", k)
			}
		}
	}
	if err := json.Unmarshal([]byte(isi), &j); err != nil || kode != 200 || len(j.Baris) != 1 || len(j.Baris[0]) != len(kunciMau) ||
		j.Baris[0]["isTopRisk"] != true || j.Baris[0]["otherType"] != "UJI LAIN" || j.Baris[0]["numberOfFloor"] != "3" {
		t.Fatalf("%d %s", kode, isi)
	}
	for _, u := range []struct {
		nama, badan, pelaku string
		svc                 *services.Service
		kode                int
	}{
		{"tanpa identitas", badan, "", svc, 401},
		{"JSON rusak", `{"baris":[`, "UJI-USER", svc, 400},
		{"tipe kosong", `{"baris":[{"objectType":""}]}`, "UJI-USER", svc, 400},
		{"tanpa DB", badan, "UJI-USER", services.Baru(nil), 503},
	} {
		if kode, isi := minta(t, u.svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek", u.badan, u.pelaku); kode != u.kode || !strings.Contains(isi, `"galat"`) {
			t.Errorf("%s: %d %s, mau %d", u.nama, kode, isi, u.kode)
		}
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/kasus/UJI-NB-1/objek", "", ""); kode != 503 {
		t.Errorf("GET tanpa DB: %d", kode)
	}
}
