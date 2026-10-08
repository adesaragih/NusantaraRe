package services_test

// Uji modal `View File` — unduh, hapus, ganti kategori (`ShowAttachmentTreaty`).

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

func (g *gudangTiruan) BacaObjekSimpanan(_ context.Context, imageID string) (models.ObjekSimpanan, bool, error) {
	o, ada := g.objekSimpanan[imageID]
	return o, ada, nil
}

func (g *gudangTiruan) PerbaruiObjekSimpanan(_ context.Context, o models.ObjekSimpanan, _ string) error {
	g.objekDiperbarui = append(g.objekDiperbarui, o)
	return nil
}

func (g *gudangTiruan) HapusLampiran(_ context.Context, _, idLampiran, imageID string) error {
	g.lampiranDihapus = append(g.lampiranDihapus, idLampiran+"|"+imageID)
	return nil
}

func (g *gudangTiruan) UbahKategoriLampiran(_ context.Context, _ string, ubah []models.PerubahanKategori) error {
	g.kategoriDiubah = append(g.kategoriDiubah, ubah...)
	return nil
}

func gudangViewFile(exp string) *gudangTiruan {
	g := gudangLampiran()
	g.lampiran = []models.BarisLampiranWarisan{
		{ID: "L1", KodeKategori: "00002", NamaKategori: "Approval Email", NamaBerkas: "Bordero.xlsx", JenisMime: "xlsx", IDSimpanan: "IMG1"},
	}
	g.objekSimpanan = map[string]models.ObjekSimpanan{"IMG1": {
		ImageID: "IMG1", URLPublik: "https://storage.googleapis.com/rnmtest/lama?sig=1",
		AppFolder: "gs://rnmtest/Contract/Doc/2025/08/x - Bordero.xlsx", Exp: exp, App: "rnmtest", NamaObjek: "x - Bordero.xlsx",
	}}
	return g
}

func TestTautanLampiranMemakaiURLTersimpanSelamaBerlaku(t *testing.T) {
	g, s := gudangViewFile("01/01/2099 00:00:00"), &simpananTiruan{}
	u, err := services.LayananDenganSimpanan(g, s).TautanLampiran(context.Background(), admin, "1001001", "L1", false)
	if err != nil {
		t.Fatal(err)
	}
	if u != "https://storage.googleapis.com/rnmtest/lama?sig=1" || len(s.urlBaru) != 0 {
		t.Errorf("url %q, geturl %d", u, len(s.urlBaru))
	}
}

// [6] EXPDATE lewat → Google/geturl (Folder tanpa Namafile & gs://App/), lalu UPDATE.
func TestTautanLampiranKedaluwarsaMemintaURLBaruLaluMemperbarui(t *testing.T) {
	g, s := gudangViewFile("01/01/2020 00:00:00"), &simpananTiruan{}
	u, err := services.LayananDenganSimpanan(g, s).TautanLampiran(context.Background(), admin, "1001001", "L1", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.urlBaru) != 1 || s.urlBaru[0].Folder != "Contract/Doc/2025/08/" || s.urlBaru[0].Namafile != "x - Bordero.xlsx" || s.urlBaru[0].Durasi != 1800 {
		t.Errorf("geturl %+v", s.urlBaru)
	}
	if u != "https://storage.googleapis.com/rnmtest/baru" || len(g.objekDiperbarui) != 1 || g.objekDiperbarui[0].Exp != "01/01/2030 00:00:00" {
		t.Errorf("url %q diperbarui %+v", u, g.objekDiperbarui)
	}
}

func TestTautanLampiranViewOfficeOnline(t *testing.T) {
	g, s := gudangViewFile("01/01/2099 00:00:00"), &simpananTiruan{}
	u, err := services.LayananDenganSimpanan(g, s).TautanLampiran(context.Background(), admin, "1001001", "L1", true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(u, "https://view.officeapps.live.com/op/view.aspx?src=https%3A%2F%2Fstorage.googleapis.com") {
		t.Errorf("office %q", u)
	}
}

// Delete_act: Google/delete (Namafile = APPFOLDER tanpa gs://App/) LALU baris.
func TestHapusLampiranObjekLaluBaris(t *testing.T) {
	g, s := gudangViewFile("01/01/2099 00:00:00"), &simpananTiruan{}
	if _, err := services.LayananDenganSimpanan(g, s).HapusLampiran(context.Background(), admin, "1001001", "L1"); err != nil {
		t.Fatal(err)
	}
	if len(s.dihapus) != 1 || s.dihapus[0].Namafile != "Contract/Doc/2025/08/x - Bordero.xlsx" || s.dihapus[0].App != "rnmtest" {
		t.Errorf("delete %+v", s.dihapus)
	}
	if len(g.lampiranDihapus) != 1 || g.lampiranDihapus[0] != "L1|IMG1" {
		t.Errorf("baris %v", g.lampiranDihapus)
	}
}

// ⚠️ Penyimpangan yang dinyatakan: hapus objek GAGAL → baris TETAP.
func TestHapusLampiranGagalDiStorageBarisTetap(t *testing.T) {
	g, s := gudangViewFile("01/01/2099 00:00:00"), &simpananTiruan{galat: services.ErrSimpananGagal}
	_, err := services.LayananDenganSimpanan(g, s).HapusLampiran(context.Background(), admin, "1001001", "L1")
	if !errors.Is(err, services.ErrSimpananGagal) || len(g.lampiranDihapus) != 0 {
		t.Errorf("galat %v, baris dihapus %v", err, g.lampiranDihapus)
	}
}

func TestUbahKategoriLampiranNamaDariKatalogDanNonProp(t *testing.T) {
	g := gudangViewFile("01/01/2099 00:00:00")
	if _, err := services.LayananDengan(g).UbahKategoriLampiran(context.Background(), admin, "1001001",
		[]services.MasukanUbahKategori{{IDLampiran: "L1", KodeKategori: "00007"}}); err != nil {
		t.Fatal(err)
	}
	u := g.kategoriDiubah[0]
	// Kontrak uji NonProportional → nama Non-Prop untuk 00007.
	if u.IDLampiran != "L1" || u.KodeKategori != "00007" || u.NamaKategori != "Pega Non Proportional Calculation /Perhitungan Pega Non Proportional" {
		t.Errorf("perubahan %+v", u)
	}
}

func TestUbahKategoriLampiranDitolak(t *testing.T) {
	ctx := context.Background()
	g := gudangViewFile("01/01/2099 00:00:00")
	l := services.LayananDengan(g)
	if _, err := l.UbahKategoriLampiran(ctx, admin, "1001001", []services.MasukanUbahKategori{{IDLampiran: "L1", KodeKategori: "99999"}}); !errors.Is(err, services.ErrTombolDitolak) {
		t.Errorf("kategori tak ada: %v", err)
	}
	if _, err := l.UbahKategoriLampiran(ctx, admin, "1001001", []services.MasukanUbahKategori{{IDLampiran: "ASING", KodeKategori: "00002"}}); !errors.Is(err, services.ErrTombolDitolak) {
		t.Errorf("lampiran asing: %v", err)
	}
	// Syarat wadah: bukan Resolve Complete / Decline.
	g.kepalaTreatyIn["1001001"]["StatusAkseptasi"] = "Resolve Complete"
	if _, err := l.UbahKategoriLampiran(ctx, admin, "1001001", []services.MasukanUbahKategori{{IDLampiran: "L1", KodeKategori: "00002"}}); !errors.Is(err, services.ErrTombolDitolak) {
		t.Errorf("tuntas: %v", err)
	}
	if len(g.kategoriDiubah) != 0 {
		t.Errorf("tertulis %v", g.kategoriDiubah)
	}
}
