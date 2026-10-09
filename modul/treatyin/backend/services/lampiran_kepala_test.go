package services_test

import (
	"context"
	"errors"
	"testing"

	"nusantarare/modul/treatyin/backend/services"
)

// Panel Attachment Adjustment memakai rute lampiran Treaty In: ID penyesuaian
// hanya ada di `TREATY_IN_EDM`, wajib diterima Upload / panel / Delete.
func TestLampiranMenerimaIDPenyesuaian(t *testing.T) {
	g, s := gudangViewFile("01/01/2099 00:00:00"), &simpananTiruan{}
	g.kepalaEDM = map[string]map[string]any{"1001001/R01": {"ProportionType": "Proportional"}}
	l := services.LayananDenganSimpanan(g, s)
	ctx := context.Background()
	if _, err := l.BacaPanelLampiran(ctx, admin, "1001001/R01"); err != nil {
		t.Fatalf("panel penyesuaian: %v", err)
	}
	h, err := l.UnggahLampiran(ctx, admin, services.MasukanUnggahLampiran{
		IDKontrak: "1001001/R01", KodeKategori: "00002",
		Berkas: []services.BerkasUnggah{{Nama: "a.pdf", Isi: []byte("%PDF-1.4 isi berkas uji")}},
	})
	if err != nil || len(h.Berkas) != 1 || !h.Berkas[0].Berhasil {
		t.Fatalf("unggah penyesuaian: %+v %v", h, err)
	}
	if got := g.lampiranBaru[0].IDKontrak; got != "1001001/R01" {
		t.Errorf("lampiran menempel pada %q", got)
	}
	if _, err := l.HapusLampiran(ctx, admin, "1001001/R01", "L1"); err != nil {
		t.Errorf("hapus penyesuaian: %v", err)
	}
}

// ID yang tidak ada di kedua tabel tetap ditolak.
func TestLampiranIDTakDikenalTetapDitolak(t *testing.T) {
	g, s := gudangViewFile("01/01/2099 00:00:00"), &simpananTiruan{}
	if _, err := services.LayananDenganSimpanan(g, s).BacaPanelLampiran(context.Background(), admin, "9999999/R01"); !errors.Is(err, services.ErrKontrakTidakAda) {
		t.Errorf("galat %v", err)
	}
}
