package services_test

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

// `Download All` — DownloadAll_Act: semua lampiran kontrak dalam satu zip,
// entri = nama berkas; kembaran diberi akhiran.
func TestUnduhSemuaLampiranMengemasSemuaBerkas(t *testing.T) {
	g, s := gudangViewFile("01/01/2099 00:00:00"), &simpananTiruan{}
	g.lampiran = append(g.lampiran,
		models.BarisLampiranWarisan{ID: "L2", NamaBerkas: "Bordero.xlsx", JenisMime: "xlsx", IDSimpanan: "IMG1"},
		models.BarisLampiranWarisan{ID: "L3", NamaBerkas: "Slip.pdf", JenisMime: "pdf", IDSimpanan: "IMG1"},
	)
	isi, err := services.LayananDenganSimpanan(g, s).UnduhSemuaLampiran(context.Background(), admin, "1001001")
	if err != nil {
		t.Fatal(err)
	}
	zr, err := zip.NewReader(bytes.NewReader(isi), int64(len(isi)))
	if err != nil {
		t.Fatal(err)
	}
	var nama []string
	for _, f := range zr.File {
		nama = append(nama, f.Name)
		r, _ := f.Open()
		b, _ := io.ReadAll(r)
		_ = r.Close()
		if string(b) != "isi-berkas" {
			t.Errorf("%s berisi %q", f.Name, b)
		}
	}
	ingin := []string{"Bordero.xlsx", "Bordero (2).xlsx", "Slip.pdf"}
	if len(nama) != len(ingin) {
		t.Fatalf("entri %v, ingin %v", nama, ingin)
	}
	for i := range ingin {
		if nama[i] != ingin[i] {
			t.Errorf("entri %v, ingin %v", nama, ingin)
		}
	}
	if services.NamaUnduhSemuaLampiran != "AllDocuments.zip" {
		t.Errorf("nama zip %q", services.NamaUnduhSemuaLampiran)
	}
}

func TestUnduhSemuaLampiranTanpaIdentitasDitolak(t *testing.T) {
	g, s := gudangViewFile("01/01/2099 00:00:00"), &simpananTiruan{}
	if _, err := services.LayananDenganSimpanan(g, s).UnduhSemuaLampiran(context.Background(), inti.Pelaku{}, "1001001"); err == nil {
		t.Fatal("tanpa identitas wajib ditolak")
	}
}
