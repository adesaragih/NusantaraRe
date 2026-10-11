package handlers

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

type portalTiruan struct {
	baris   []models.BarisPortal
	total   int
	err     error
	selesai *string
}

func (p portalTiruan) CariPortal(_ context.Context, _, _, selesai string, _, _ int) ([]models.BarisPortal, int, error) {
	if p.selesai != nil {
		*p.selesai = selesai
	}
	return p.baris, p.total, p.err
}

// TestCariPortal - GET /api/nbfacin/opportunity (tiket 32): bentuk jawaban persis kontrak
// frontend sesi 0f, `baris` larik walau kosong, 400/401/503/500; POST rute yang sama tetap
// membuat case (tidak tertimpa).
func TestCariPortal(t *testing.T) {
	svc := services.Baru(nil).DenganPortal(portalTiruan{total: 16, baris: []models.BarisPortal{
		{CaseID: "UJI-NB-2", Name: "UJI Prospek", GroupBusiness: "UJI GRUP", InsuredName: "UJI TERTANGGUNG", Marketing: "UJI MO", Status: "Pending-Policy",
			NBStatus: "NB IS IN TEAM LEADER'S INBOX", NoPolis: "RNM-F1.10.2026.00001", TglProduksi: "11-10-2026"},
		{CaseID: "UJI-NB-1"},
	}})
	kode, isi := minta(t, svc, "GET", "/api/nbfacin/opportunity?cari=uji&halaman=2", "", "UJI-USER")
	mau := `{"baris":[{"caseId":"UJI-NB-2","name":"UJI Prospek","groupBusiness":"UJI GRUP","insuredName":"UJI TERTANGGUNG",` +
		`"marketing":"UJI MO","status":"Pending-Policy","nbStatus":"NB IS IN TEAM LEADER'S INBOX","noPolis":"RNM-F1.10.2026.00001",` +
		`"tglProduksi":"11-10-2026"},{"caseId":"UJI-NB-1","name":"",` +
		`"groupBusiness":"","insuredName":"","marketing":"","status":"","nbStatus":"NEW NB","noPolis":"","tglProduksi":""}],"total":16,"halaman":2,"ukuran":15}`
	if kode != 200 || isi != mau {
		t.Fatalf("%d %s\nmau %s", kode, isi, mau)
	}
	if kode, isi := minta(t, services.Baru(nil).DenganPortal(portalTiruan{}), "GET", "/api/nbfacin/opportunity", "", "UJI-USER"); kode != 200 ||
		isi != `{"baris":[],"total":0,"halaman":1,"ukuran":15}` {
		t.Errorf("kosong: %d %s", kode, isi)
	}
	for _, u := range []struct {
		nama, jalur, pelaku string
		svc                 *services.Service
		kode                int
	}{
		{"tanpa identitas", "/api/nbfacin/opportunity", "", svc, 401},
		{"tanpa identitas, halaman teks", "/api/nbfacin/opportunity?halaman=satu", "", svc, 401},
		{"halaman 0", "/api/nbfacin/opportunity?halaman=0", "UJI-USER", svc, 400},
		{"halaman teks", "/api/nbfacin/opportunity?halaman=satu", "UJI-USER", svc, 400},
		{"selesai tak dikenal", "/api/nbfacin/opportunity?selesai=2", "UJI-USER", svc, 400},
		{"tanpa DB", "/api/nbfacin/opportunity", "UJI-USER", services.Baru(nil), 503},
		{"galat Oracle", "/api/nbfacin/opportunity", "UJI-USER", services.Baru(nil).DenganPortal(portalTiruan{err: errors.New("ORA-UJI rincian rahasia")}), 500},
	} {
		if kode, isi := minta(t, u.svc, "GET", u.jalur, "", u.pelaku); kode != u.kode || !strings.Contains(isi, `"galat"`) || strings.Contains(isi, "rahasia") {
			t.Errorf("%s: %d %s, mau %d", u.nama, kode, isi, u.kode)
		}
	}
	// POST /api/nbfacin/opportunity tetap rute pembuat case.
	if kode, _ := minta(t, services.Baru(nil), "POST", "/api/nbfacin/opportunity", badanSah, "UJI-USER"); kode != 503 {
		t.Errorf("POST: %d, mau 503 (buat case tanpa DB)", kode)
	}
}

// TestCariPortalSelesaiKabel - GET /api/nbfacin/opportunity?selesai=1|0 diteruskan; tanpa kueri = "" (semua, perilaku lama).
func TestCariPortalSelesaiKabel(t *testing.T) {
	var selesai string
	svc := services.Baru(nil).DenganPortal(portalTiruan{selesai: &selesai})
	for jalur, mau := range map[string]string{"/api/nbfacin/opportunity?selesai=1": "1", "/api/nbfacin/opportunity?selesai=0&halaman=2": "0",
		"/api/nbfacin/opportunity": ""} {
		selesai = "awal"
		if kode, isi := minta(t, svc, "GET", jalur, "", "UJI-USER"); kode != 200 || selesai != mau {
			t.Errorf("%s: %d %s diteruskan %q", jalur, kode, isi, selesai)
		}
	}
}
