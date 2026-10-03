package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// objekTiruan - penyimpan berkeadaan: daftar per case, diganti utuh.
type objekTiruan struct{ ada map[string][]models.ObjekFire }

func (o objekTiruan) BacaObjek(_ context.Context, id string) ([]models.ObjekFire, error) {
	d, ok := o.ada[id]
	if !ok {
		return nil, repository.ErrKasusTidakAda
	}
	return append([]models.ObjekFire{}, d...), nil
}

func (o objekTiruan) GantiObjek(_ context.Context, _ *db.Tx, id string, baris []models.ObjekFire) error {
	if _, ok := o.ada[id]; !ok {
		return repository.ErrKasusTidakAda
	}
	o.ada[id] = append([]models.ObjekFire{}, baris...)
	return nil
}

// TestGantiObjek - tiket 35: urutan dipertahankan, ganti utuh, kosong mengosongkan, dibaca
// ulang; 401 / 400 (indeks baris) / 503 / 404.
func TestGantiObjek(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-USER"}
	svc := Baru(nil).DenganObjek(objekTiruan{ada: map[string][]models.ObjekFire{"UJI-NB-1": {}}}).DenganTransaksi(tanpaTx)
	d, err := svc.GantiObjek(ctx, akun, "UJI-NB-1", []models.ObjekFire{
		{ObjectNo: "2", ObjectType: "UJI TIPE B", IsTopRisk: true, NumberOfFloor: "0"},
		{ObjectNo: "1", ObjectType: "UJI TIPE A", NumberOfFloor: "12", RoofType: "14"},
	})
	if err != nil || len(d) != 2 || d[0].ObjectNo != "2" || !d[0].IsTopRisk || d[1].RoofType != "14" {
		t.Fatalf("%+v (%v)", d, err)
	}
	if d, err := svc.GantiObjek(ctx, akun, "UJI-NB-1", nil); err != nil || len(d) != 0 {
		t.Errorf("kosong: %+v (%v)", d, err)
	}
	if _, err := svc.GantiObjek(ctx, inti.Pelaku{}, "UJI-NB-1", nil); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
	for nama, u := range map[string]struct {
		baris []models.ObjekFire
		pesan string
	}{
		"tipe kosong":    {[]models.ObjekFire{{ObjectType: "UJI"}, {ObjectType: "  "}}, "baris[1].objectType wajib"},
		"lantai minus":   {[]models.ObjekFire{{ObjectType: "UJI", NumberOfFloor: "-1"}}, "baris[0].numberOfFloor"},
		"lantai desimal": {[]models.ObjekFire{{ObjectType: "UJI", NumberOfFloor: "1.5"}}, "baris[0].numberOfFloor"},
		"lantai teks":    {[]models.ObjekFire{{ObjectType: "UJI", NumberOfFloor: "dua"}}, "baris[0].numberOfFloor"},
		"nama 501 bita":  {[]models.ObjekFire{{ObjectType: "UJI", ObjectName: strings.Repeat("U", 501)}}, "baris[0].objectName paling banyak 500"},
		"lokasi 51":      {[]models.ObjekFire{{ObjectType: "UJI", RiskLocation: strings.Repeat("U", 51)}}, "baris[0].riskLocation paling banyak 50"},
	} {
		if _, err := Baru(nil).GantiObjek(ctx, akun, "UJI-NB-1", u.baris); !errors.Is(err, ErrMasukanObjek) || !strings.Contains(err.Error(), u.pesan) {
			t.Errorf("%s: %v, mau 400 %q (sebelum DB)", nama, err, u.pesan)
		}
	}
	if _, err := Baru(nil).GantiObjek(ctx, akun, "UJI-NB-1", nil); !errors.Is(err, ErrObjekTanpaDatabase) {
		t.Errorf("tanpa DB: %v", err)
	}
	if _, err := svc.GantiObjek(ctx, akun, "UJI-NB-TIDAK-ADA", nil); !errors.Is(err, ErrKasusTidakAda) {
		t.Errorf("case tidak ada: %v", err)
	}
	if _, err := svc.BacaObjek(ctx, "UJI-NB-TIDAK-ADA"); !errors.Is(err, ErrKasusTidakAda) {
		t.Errorf("baca case tidak ada: %v", err)
	}
}
