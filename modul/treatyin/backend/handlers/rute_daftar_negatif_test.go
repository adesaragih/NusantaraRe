package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"nusantarare/modul/treatyin/backend/handlers"
	"nusantarare/modul/treatyin/backend/services"
)

// Gudang tiruan rute: satu Ceding berstatus "inactive", satu SoB aktif.
func (g gudangTiruan) BacaStatusAktifAgen(_ context.Context, _ []string) (map[string]string, error) {
	return map[string]string{"G0000024": "inactive", "G0000099": "1"}, g.galat
}

func TestRuteDaftarNegatifAgen(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)
	w := minta(t, h, handlers.Prefix+"/agen/daftar-negatif?cedant=G0000024&asalBisnis=G0000099", "AKUN-UJI")
	if w.Code != http.StatusOK {
		t.Fatalf("kode %d: %s", w.Code, w.Body.String())
	}
	var hasil services.HasilDaftarNegatifAgen
	if err := json.Unmarshal(w.Body.Bytes(), &hasil); err != nil {
		t.Fatal(err)
	}
	if !hasil.Cedant.DaftarNegatif || hasil.Cedant.Pesan != "This name is on Agent Negative List" {
		t.Errorf("cedant %+v", hasil.Cedant)
	}
	if hasil.AsalBisnis.DaftarNegatif || hasil.AsalBisnis.Pesan != "" {
		t.Errorf("asal bisnis %+v", hasil.AsalBisnis)
	}

	if w := minta(t, h, handlers.Prefix+"/agen/daftar-negatif?cedant=G0000024", ""); w.Code != http.StatusUnauthorized {
		t.Errorf("tanpa identitas: kode %d", w.Code)
	}
}
