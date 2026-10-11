package handlers

import (
	"context"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
	"nusantarare/modul/nbfacin/backend/services"
)

// unggahTiruanH - satu pekerjaan cargo "7" milik UJI-NB-1 berjudul urut berkas; baris 1 = isi kanonik cargo.
type unggahTiruanH struct{}

func (unggahTiruanH) BuatUnggah(context.Context, *db.Tx, models.UnggahCsv) (string, error) {
	return "", nil
}
func (unggahTiruanH) UbahUnggah(context.Context, *db.Tx, models.UnggahCsv) error { return nil }
func (unggahTiruanH) BacaUnggah(_ context.Context, id string) (models.UnggahCsv, error) {
	if id != "7" {
		return models.UnggahCsv{}, repository.ErrUnggahTidakAda
	}
	return models.UnggahCsv{ID: "7", CaseID: "UJI-NB-1", Status: "valid", Jenis: models.JenisCsvCargo, TotalBaris: 1,
		Kolom: []string{"TOP_RISK", "TRADING"}}, nil
}
func (unggahTiruanH) UnggahTerakhir(context.Context, string) (models.UnggahCsv, bool, error) {
	return models.UnggahCsv{}, false, nil
}
func (unggahTiruanH) SisipBarisUnggah(context.Context, *db.Tx, string, string, []models.BarisUnggahCsv) error {
	return nil
}
func (unggahTiruanH) BacaBarisUnggah(_ context.Context, _, jenis string, _, _ int) ([]models.BarisUnggahCsv, error) {
	isi := make([]string, len(models.KolomCsvJenis(jenis)))
	isi[0], isi[len(isi)-1] = "DOMESTIK", "TRUE"
	return []models.BarisUnggahCsv{{No: 1, Isi: isi}}, nil
}
func (unggahTiruanH) SisipPesanUnggah(context.Context, *db.Tx, string, string, int, []models.PesanUnggahCsv) error {
	return nil
}
func (unggahTiruanH) BacaPesanUnggah(context.Context, string, string, int) ([]models.PesanUnggahCsv, error) {
	return nil, nil
}
func (unggahTiruanH) HapusPesanUnggah(context.Context, *db.Tx, string) error { return nil }
func (unggahTiruanH) AdaObjek(context.Context, string, string) (bool, error) { return false, nil }

// TestBarisUnggahCargoKabel - tiket 74 P6: GET …/baris memuat `kolom` (judul urut berkas) dan baris dipetakan per kolom cargo;
// status memuat `jenis`.
func TestBarisUnggahCargoKabel(t *testing.T) {
	svc := layananKasus().DenganUnggahCsv(unggahTiruanH{}, nil)
	kode, isi := minta(t, svc, "GET", "/api/nbfacin/kasus/UJI-NB-1/objek/csv/7/baris?halaman=1&ukuran=10", "", "")
	if kode != 200 || !strings.Contains(isi, `"kolom":["TOP_RISK","TRADING"]`) || !strings.Contains(isi, `"TRADING":"DOMESTIK"`) ||
		!strings.Contains(isi, `"TOP_RISK":"TRUE"`) || strings.Contains(isi, `"NO_OFFER_SLIP"`) {
		t.Errorf("%d %s", kode, isi)
	}
	k := keStatusUnggahKabel(services.StatusUnggahCsv{Unggah: models.UnggahCsv{ID: "7", Jenis: models.JenisCsvCargo}})
	if k.Jenis != "cargo" {
		t.Errorf("jenis: %+v", k)
	}
}
