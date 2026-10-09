package services_test

// Kapan tombol tulis menyertakan baris detail — `SaveTreatyIn_Act` [8] dan
// `SaveTreatyIn_EDM_Act` [14]: HANYA bila dokumen ber-`Resolve Complete`.

import (
	"context"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

// detailSatuBaris - satu Limits, satu Detail tanpa IOOLimitList: tepat satu
// sisipan dini.
func detailSatuBaris(doc map[string]any) {
	doc["ProportionType"] = "Proportional"
	doc["Limits"] = larik(map[string]any{"TreatyType": "QUOTA SHARE",
		"Detail": larik(map[string]any{"TreatyGroup": "PROPERTY", "TreatyGroupID": "10007"})})
}

func TestDirectorAcceptMenulisDetailTreatyIn(t *testing.T) {
	g := gudangSimpan()
	g.dokumenTersimpan["1001001"]["Position"] = models.PosisiDirector
	g.dokumenTersimpan["1001001"]["StatusAkseptasi"] = "Accept"
	detailSatuBaris(g.dokumenTersimpan["1001001"])
	_, err := services.LayananDengan(g).KirimKontrak(context.Background(), director, services.MasukanKirim{
		MasukanSimpan: services.MasukanSimpan{IDKontrak: "1001001"}, Aksi: services.AksiAkseptasi, Pilihan: "Accept",
	})
	if err != nil {
		t.Fatal(err)
	}
	r := g.disimpan[0]
	if r.Detail == nil || len(r.Detail.Baris) != 1 {
		t.Fatalf("detail %+v, mau satu baris", r.Detail)
	}
	if b := r.Detail.Baris[0]; b.Teks["TREATYID"] != "1001001" || b.Teks["TREATYGROUP"] != "PROPERTY" {
		t.Errorf("baris %v", b.Teks)
	}
}

func TestBelumTuntasDetailTidakDisentuh(t *testing.T) {
	g := gudangSimpan()
	detailSatuBaris(g.dokumenTersimpan["1001001"])
	l := services.LayananDengan(g)
	if _, err := l.SimpanKontrak(context.Background(), admin, services.MasukanSimpan{IDKontrak: "1001001"}); err != nil {
		t.Fatal(err)
	}
	if _, err := l.KirimKontrak(context.Background(), admin, services.MasukanKirim{
		MasukanSimpan: services.MasukanSimpan{IDKontrak: "1001001"}, Aksi: services.AksiSubmit,
	}); err != nil {
		t.Fatal(err)
	}
	for i, r := range g.disimpan {
		if r.Detail != nil {
			t.Errorf("tulisan ke-%d (status %v) menyentuh TREATYINDETAIL", i+1, r.Dokumen["StatusAkseptasi"])
		}
	}
}

// Save Force Edit IT atas kontrak tuntas = `SaveTreatyIn_Act` dengan status
// Resolve Complete → detail ditulis ulang, persis syarat [8].
func TestForceEditTuntasMenulisUlangDetail(t *testing.T) {
	g := gudangSimpan()
	g.divisi = map[string]string{"ITDEV": services.DivisiIT}
	g.dokumenTersimpan["1001001"]["StatusAkseptasi"] = models.StatusTuntas
	detailSatuBaris(g.dokumenTersimpan["1001001"])
	if _, err := services.LayananDengan(g).SimpanKontrak(context.Background(), inti.Pelaku{AkunID: "ITDEV"},
		services.MasukanSimpan{IDKontrak: "1001001"}); err != nil {
		t.Fatal(err)
	}
	if r := g.disimpan[0]; r.Detail == nil || len(r.Detail.Baris) != 1 {
		t.Fatalf("detail %+v", r.Detail)
	}
}

func TestPenyesuaianTuntasMenulisDetailEDM(t *testing.T) {
	g := gudangEDM()
	g.kepalaEDM["1001001/R01"]["Position"] = models.PosisiDirector
	g.kepalaEDM["1001001/R01"]["StatusAkseptasi"] = "Accept"
	detailSatuBaris(g.dokumenTersimpan["1001001/R01"])
	if _, err := services.LayananDengan(g).KirimPenyesuaian(context.Background(), director, services.MasukanKirimPenyesuaian{
		MasukanPenyesuaian: services.MasukanPenyesuaian{ID: "1001001/R01"}, Aksi: services.AksiAkseptasi, Pilihan: "Accept",
	}); err != nil {
		t.Fatal(err)
	}
	r := g.disimpanEDM[0]
	if r.Baru["StatusAkseptasi"] != models.StatusTuntas || r.Detail == nil || len(r.Detail.Baris) != 1 {
		t.Fatalf("status %v detail %+v", r.Baru["StatusAkseptasi"], r.Detail)
	}
	if b := r.Detail.Baris[0]; b.Teks["TREATYID"] != "1001001/R01" {
		t.Errorf("TREATYID %q", b.Teks["TREATYID"])
	}
	if _, ada := r.Detail.Baris[0].Angka["SPREAD_RNM_SHARE_PCT"]; ada {
		t.Error("baris EDM memuat kolom khusus TREATYINDETAIL")
	}
}

func TestPenyesuaianBelumTuntasDetailTidakDisentuh(t *testing.T) {
	g := gudangEDM()
	detailSatuBaris(g.dokumenTersimpan["1001001/R01"])
	if _, err := services.LayananDengan(g).SimpanPenyesuaian(context.Background(), admin,
		services.MasukanPenyesuaian{ID: "1001001/R01"}); err != nil {
		t.Fatal(err)
	}
	if g.disimpanEDM[0].Detail != nil {
		t.Error("Save penyesuaian menyentuh TREATYINDETAILEDM")
	}
}
