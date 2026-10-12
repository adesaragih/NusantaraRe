package handlers

// Pendaftaran banyak case (10-10-2026): bentuk badan, pemetaan galat, dan badan jawaban - tanpa Oracle.

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
	"nusantarare/modul/claimlife/backend/services"
)

func TestBadanDaftarDuaBentuk(t *testing.T) {
	var baru permintaanDaftarJSON
	if err := json.Unmarshal([]byte(`{"type":"QP","kodeBisnis":"L1","kelompok":[{"nomorPremiList":"UJI-PL-A","nomorPolis":"UJI-POL-A","mataUang":"IDR","sertifikat":["001"]},{"nomorPremiList":"UJI-PL-B","sertifikat":["101","102"]}]}`), &baru); err != nil {
		t.Fatal(err)
	}
	p := baru.keLayanan()
	if len(p.Kelompok) != 2 || p.Kelompok[1].NomorPremiList != "UJI-PL-B" || len(p.Kelompok[1].Sertifikat) != 2 {
		t.Fatalf("bentuk baru %+v", p)
	}
	var lama permintaanDaftarJSON
	if err := json.Unmarshal([]byte(`{"nomorPremiList":"UJI-PL-A","nomorPolis":"UJI-POL-A","type":"QP","kodeBisnis":"L1","mataUang":"IDR","sertifikat":["001"]}`), &lama); err != nil {
		t.Fatal(err)
	}
	if q := lama.keLayanan(); len(q.Kelompok) != 0 || q.NomorPremiList != "UJI-PL-A" || len(q.Sertifikat) != 1 {
		t.Fatalf("bentuk lama %+v", q)
	}
}

func TestJawabanGalatDaftar(t *testing.T) {
	for _, k := range []struct {
		err    error
		status int
		memuat string
	}{
		{fmt.Errorf("%w: sertifikat 001 pada premium list UJI-PL-A dipilih lebih dari sekali", galat.ErrPermintaanTidakSah), http.StatusBadRequest, "dipilih lebih dari sekali"},
		{fmt.Errorf("%w: nomor polis UJI-POL-B tidak cocok dengan premium list UJI-PL-A", services.ErrKelompokTidakSah), http.StatusUnprocessableEntity, "tidak cocok dengan premium list UJI-PL-A"},
		{&services.GalatPesertaSudahDiklaim{Bentrok: []repository.PesertaAktifInApp{{NomorPremiList: "UJI-PL-A", Sertifikat: "001", NomorKlaim: "UJI-CLMLF-1"}}}, http.StatusConflict, `Certificate No "001" (PL UJI-PL-A) is already in claim UJI-CLMLF-1 (Outstanding)`},
		{inti.ErrTanpaIdentitas, http.StatusUnauthorized, "identity"},
		{inti.ErrTanpaWewenang, http.StatusForbidden, "role"},
		{errors.New("ORA-03113"), http.StatusInternalServerError, "Failed to register the claim."},
	} {
		status, pesan, ada := jawabanGalatDaftar(k.err)
		if !ada || status != k.status || !strings.Contains(pesan, k.memuat) {
			t.Errorf("%v: %d %q", k.err, status, pesan)
		}
		if status == http.StatusUnprocessableEntity && strings.Contains(pesan, "services:") {
			t.Errorf("422 berawalan teknis: %q", pesan)
		}
	}
	if _, _, ada := jawabanGalatDaftar(nil); ada {
		t.Error("nil dianggap galat")
	}
}

func TestJawabanDaftarDaftarCase(t *testing.T) {
	pohon := []models.PohonKlaim{
		{Work: models.WorkClaim{ID: "UJI-W1"}, Klaim: models.Klaim{NomorKlaim: "UJI-K1", Peserta: []models.Peserta{{NomorPremiList: "UJI-PL-A"}, {NomorPremiList: "UJI-PL-A"}}}},
		{Work: models.WorkClaim{ID: "UJI-W2"}, Klaim: models.Klaim{NomorKlaim: "UJI-K2", Peserta: []models.Peserta{{NomorPremiList: "UJI-PL-B"}}}},
	}
	pohon[0].Klaim.ClaimRetro.Currency = "IDR"
	pohon[1].Klaim.ClaimRetro.Currency = "USD"
	b, err := json.Marshal(jawabanDaftar(pohon))
	if err != nil {
		t.Fatal(err)
	}
	mau := `{"kasus":[{"id":"UJI-W1","nomorKlaim":"UJI-K1","nomorPremiList":"UJI-PL-A","mataUang":"IDR","jumlahPeserta":2},{"id":"UJI-W2","nomorKlaim":"UJI-K2","nomorPremiList":"UJI-PL-B","mataUang":"USD","jumlahPeserta":1}]}`
	if string(b) != mau {
		t.Errorf("jawaban %s", b)
	}
	if b, _ := json.Marshal(jawabanDaftar(nil)); string(b) != `{"kasus":[]}` {
		t.Errorf("kosong %s", b)
	}
}
