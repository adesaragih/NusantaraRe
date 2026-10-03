package handlers

// Tab Object FIRE (tiket 35): GET/PUT /api/nbfacin/kasus/{caseId}/objek - kontrak `ObjekFire`
// frontend (modul/nbfacin/frontend/api.ts).

import (
	"encoding/json"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

type objekKabel struct {
	ObjectNo         string `json:"objectNo"`
	ObjectType       string `json:"objectType"`
	ObjectName       string `json:"objectName"`
	IsMaterialDamage bool   `json:"isMaterialDamage"`
	IsTopRisk        bool   `json:"isTopRisk"`
	RoadType         string `json:"roadType"`
	RoadName         string `json:"roadName"`
	BuildingNo       string `json:"buildingNo"`
	ZipCode          string `json:"zipCode"`
	Country          string `json:"country"`
	RiskLocation     string `json:"riskLocation"`
	Territory        string `json:"territory"`
	City             string `json:"city"`
	District         string `json:"district"`
	Province         string `json:"province"`
	RiskAddressID    string `json:"riskAddressId"`
	NumberOfFloor    string `json:"numberOfFloor"`
	RoofType         string `json:"roofType"`
	WallType         string `json:"wallType"`
	FloorType        string `json:"floorType"`
	PartitionType    string `json:"partitionType"`
	SupportWallType  string `json:"supportWallType"`
	OtherType        string `json:"otherType"`
}

type daftarObjek struct {
	Baris []objekKabel `json:"baris"`
}

// batasBadanObjek - badan PUT terbesar yang diterima (A112).
const batasBadanObjek = 1 << 20

func keObjekKabel(d []models.ObjekFire) daftarObjek {
	baris := make([]objekKabel, 0, len(d))
	for _, o := range d {
		baris = append(baris, objekKabel(o))
	}
	return daftarObjek{Baris: baris}
}

// bacaObjek - GET /api/nbfacin/kasus/{caseId}/objek.
func bacaObjek(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.BacaObjek(r.Context(), r.PathValue("caseId"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keObjekKabel(d))
	}
}

// simpanObjek - PUT /api/nbfacin/kasus/{caseId}/objek (tombol Save): ganti utuh, jawab baca ulang.
func simpanObjek(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b daftarObjek
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanObjek)).Decode(&b); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan bukan JSON daftar objek yang sah")
			return
		}
		baris := make([]models.ObjekFire, len(b.Baris))
		for i, o := range b.Baris {
			baris[i] = models.ObjekFire(o)
		}
		d, err := svc.GantiObjek(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"), baris)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keObjekKabel(d))
	}
}
