package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

type tolTiruan struct {
	kode          []string
	kodeDiminta   *[2]string
	daftarDiminta *string
}

func (t tolTiruan) KodeBisnis(_ context.Context, nama, grup string) ([]string, error) {
	*t.kodeDiminta = [2]string{nama, grup}
	return t.kode, nil
}

func (t tolTiruan) DaftarTableOfLimit(_ context.Context, kode, kategori string) ([]models.BarisTableOfLimit, error) {
	*t.daftarDiminta = kode + "|" + kategori
	return []models.BarisTableOfLimit{{Description: "UJI KELAS", PctLimit: "70,000 "}}, nil
}

// TestTableOfLimit - tiket 40 / butir 89: BIZCODE = BUSINESS.ID dari nama Class of Business + group business case;
// tanpa saringan tahun dan Begin date tidak dibutuhkan (A161 menggantikan A153); 0 / > 1 baris atau Class of Business
// case kosong -> 409; kategori dipangkas dan diteruskan; 404 / 400 / 503.
func TestTableOfLimit(t *testing.T) {
	ctx := context.Background()
	var diminta [2]string
	var kat string
	kasus := kasusTiruan{ada: map[string]*models.Kasus{
		"UJI-NB-1": {CaseID: "UJI-NB-1", Opportunity: models.Opportunity{ClassOfBusiness: "UJI KELAS BISNIS", GroupBusinessID: "UJI-G1"},
			General: models.General{StartDateTime: "20251231T170000.000 GMT"}},
		"UJI-NB-2": {CaseID: "UJI-NB-2", General: models.General{StartDateTime: "20260101T050000.000 GMT"}},
		"UJI-NB-3": {CaseID: "UJI-NB-3", Opportunity: models.Opportunity{ClassOfBusiness: "UJI", GroupBusinessID: "UJI-G1"}}}}
	svc := func(kode ...string) *Service {
		return Baru(nil).DenganKasus(kasus).DenganTableOfLimit(tolTiruan{kode: kode, kodeDiminta: &diminta, daftarDiminta: &kat})
	}
	b, err := svc("10048").TableOfLimit(ctx, "UJI-NB-1", " II ")
	if err != nil || len(b) != 1 || b[0].PctLimit != "70,000 " || diminta != [2]string{"UJI KELAS BISNIS", "UJI-G1"} || kat != "10048|II" {
		t.Fatalf("%v %v %v %q", b, err, diminta, kat)
	}
	for nama, u := range map[string]struct {
		svc   *Service
		id    string
		pesan string
	}{
		"tidak ada":      {svc(), "UJI-NB-1", "tidak ada di BUSINESS"},
		"ganda":          {svc("1", "2"), "UJI-NB-1", "ganda di BUSINESS"},
		"case tanpa COB": {svc("1"), "UJI-NB-2", "Class of Business / Group Business case belum diisi"},
	} {
		if _, err := u.svc.TableOfLimit(ctx, u.id, "II"); !errors.Is(err, ErrTableOfLimitTidakSiap) || !strings.Contains(err.Error(), u.pesan) {
			t.Errorf("%s: %v", nama, err)
		}
	}
	// Begin date kosong tidak lagi 409 (A161).
	if b, err := svc("7").TableOfLimit(ctx, "UJI-NB-3", ""); err != nil || len(b) != 1 || kat != "7|" {
		t.Errorf("tanpa Begin: %v %v %q", b, err, kat)
	}
	if _, err := svc("1").TableOfLimit(ctx, "UJI-NB-TIDAK-ADA", ""); !errors.Is(err, ErrKasusTidakAda) {
		t.Errorf("404: %v", err)
	}
	if _, err := svc("1").TableOfLimit(ctx, "UJI-NB-1", strings.Repeat("I", 51)); !errors.Is(err, ErrMasukanTableOfLimit) {
		t.Errorf("400: %v", err)
	}
	if _, err := Baru(nil).TableOfLimit(ctx, "UJI-NB-1", ""); !errors.Is(err, ErrTableOfLimitTanpaDatabase) {
		t.Errorf("503: %v", err)
	}
}
