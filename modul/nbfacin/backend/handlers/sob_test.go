package handlers

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

type sobTiruan struct{ baris []models.SOB }

func (s sobTiruan) CariSOB(context.Context, string, int, int) ([]models.SOB, int, error) {
	return s.baris, len(s.baris), nil
}

// TestCariSOB - GET /api/nbfacin/sob (tiket 33): bentuk jawaban persis kontrak frontend
// ({baris:[{id, clientId, name}], total, halaman, ukuran 20}); 400 halaman; 503 tanpa DB.
func TestCariSOB(t *testing.T) {
	svc := services.Baru(nil).DenganSOB(sobTiruan{baris: []models.SOB{{ID: "UJI-G2", ClientID: "UJI-C2", Name: "UJI SOB"}, {ID: "UJI-G1"}}})
	kode, isi := minta(t, svc, "GET", "/api/nbfacin/sob?cari=uji&halaman=1", "", "")
	mau := `{"baris":[{"id":"UJI-G2","clientId":"UJI-C2","name":"UJI SOB"},{"id":"UJI-G1","clientId":"","name":""}],"total":2,"halaman":1,"ukuran":20}`
	if kode != 200 || isi != mau {
		t.Fatalf("%d %s\nmau %s", kode, isi, mau)
	}
	for _, u := range []struct {
		jalur string
		svc   *services.Service
		kode  int
	}{{"/api/nbfacin/sob?halaman=satu", svc, 400}, {"/api/nbfacin/sob?halaman=0", svc, 400}, {"/api/nbfacin/sob", services.Baru(nil), 503}} {
		if kode, isi := minta(t, u.svc, "GET", u.jalur, "", ""); kode != u.kode || !strings.Contains(isi, `"galat"`) {
			t.Errorf("%s: %d %s", u.jalur, kode, isi)
		}
	}
}

// TestSimpanGeneralKodeSOB - PUT .../general meneruskan sourceOfBusinessId; GET mengembalikannya.
func TestSimpanGeneralKodeSOB(t *testing.T) {
	kode, isi := minta(t, layananKasus(), "PUT", "/api/nbfacin/kasus/UJI-NB-1/general", `{"sourceOfBusinessId":"UJI-G1"}`, "UJI-USER")
	if kode != 200 || !strings.Contains(isi, `"sourceOfBusinessId":"UJI-G1"`) {
		t.Errorf("%d %s", kode, isi)
	}
}

// TestSimpanGeneralCeding - PUT .../general meneruskan cedingIds urut; GET mengembalikan
// cedingList (larik, kosong = []).
func TestSimpanGeneralCeding(t *testing.T) {
	kode, isi := minta(t, layananKasus(), "GET", "/api/nbfacin/kasus/UJI-NB-1", "", "")
	if kode != 200 || !strings.Contains(isi, `"cedingList":[]`) {
		t.Fatalf("%d %s", kode, isi)
	}
	kode, isi = minta(t, layananKasusCeding(), "PUT", "/api/nbfacin/kasus/UJI-NB-1/general", `{"cedingIds":["UJI-B","UJI-A"]}`, "UJI-USER")
	if kode != 200 || !strings.Contains(isi, `"cedingList":[{"id":"UJI-B","name":"UJI-B"},{"id":"UJI-A","name":"UJI-A"}]`) {
		t.Errorf("%d %s", kode, isi)
	}
}

// kasusTolakTiruan - simpan General selalu menolak dengan galat masukan services (400).
type kasusTolakTiruan struct{ kasusTiruan }

func (c kasusTolakTiruan) SimpanGeneral(context.Context, *db.Tx, string, models.General) error {
	return fmt.Errorf("%w: cedingIds[1] UJI ditolak", services.ErrMasukanGeneral)
}

// TestSimpanGeneralCedingDitolak - kode ceding yang tidak lolos sampai ke 400 di kabel.
func TestSimpanGeneralCedingDitolak(t *testing.T) {
	svc := services.Baru(nil).DenganKasus(kasusTolakTiruan{kasusTiruan{k: &models.Kasus{CaseID: "UJI-NB-1"}}}).DenganTransaksi(tanpaOracle)
	kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/general", `{"cedingIds":["UJI-A","UJI-TOLAK"]}`, "UJI-USER")
	if kode != 400 || !strings.Contains(isi, "cedingIds[1]") {
		t.Errorf("%d %s", kode, isi)
	}
}
