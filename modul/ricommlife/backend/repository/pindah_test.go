package repository

import (
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/ricommlife/backend/models"
)

// Butir 3: -jalankan ditolak bila IS_PEGA_PROD=true; uji kering selalu boleh.
func TestPeriksaMode(t *testing.T) {
	if err := PeriksaMode(true, true); !errors.Is(err, ErrPindahDiProduksi) {
		t.Errorf("jalankan di produksi: %v", err)
	}
	for _, k := range [][2]bool{{false, true}, {false, false}, {true, false}} {
		if err := PeriksaMode(k[0], k[1]); err != nil {
			t.Errorf("%v: %v", k, err)
		}
	}
}

// Rencana: baris sah ditulis; angka JSON (teks atau bilangan) dikanonikkan dan dicatat; kosong menjadi NULL; nilai
// yang tidak muat = gagal TANPA nilai di laporan; ID flat yang sama dilewati, berbeda menahan; flat saja dibiarkan.
func TestRencanaPindah(t *testing.T) {
	sumber := []BarisJSON{
		{ID: "1000001", JSON: `{"IDUSEDBY":"1000003","USEDBY":"UJI COMM","CONTRACT":"1","YEAR":"2026","COMM":"12.5"}`},
		{ID: "1000002", JSON: `{"IDUSEDBY":1000003,"USEDBY":"UJI COMM","CONTRACT":"05","YEAR":2021,"COMM":"0,5"}`},
		{ID: "1000004", JSON: `{"IDUSEDBY":"1000003","USEDBY":"UJI COMM","CONTRACT":"1","YEAR":"2022","COMM":"RAHASIA-UJI"}`},
		{ID: "1000005", JSON: `{"IDUSEDBY":"1000003","USEDBY":"UJI COMM","CONTRACT":"2","YEAR":"2022"}`},
		{ID: "1000006", JSON: `{"IDUSEDBY":"1000003","USEDBY":"UJI COMM","CONTRACT":"3","YEAR":"2023","COMM":"1"}`},
		{ID: "1000007", JSON: `{"IDUSEDBY":"1000003","USEDBY":"UJI COMM","CONTRACT":"4","YEAR":"2024","COMM":"1"}`},
		{ID: "1000008", JSON: `bukan json`},
	}
	flat := []models.Komisi{
		{ID: "1000006", IDUsedBy: "1000003", UsedBy: "UJI COMM", Contract: "3", Year: "2023", Comm: "1"},
		{ID: "1000007", IDUsedBy: "1000003", UsedBy: "UJI COMM", Contract: "4", Year: "2024", Comm: "2"},
		{ID: "1000099", IDUsedBy: "1000003", UsedBy: "UJI COMM", Contract: "9", Year: "2029", Comm: "9"},
	}
	lap, tulis := RencanaPindah(sumber, flat, false)
	if lap.Sumber != 7 || lap.AkanDitulis != 3 || lap.SudahSama != 1 || lap.FlatSaja != 1 ||
		len(lap.Berbeda) != 1 || lap.Berbeda[0] != "1000007" {
		t.Errorf("laporan %+v", lap)
	}
	mau := []models.Komisi{
		{ID: "1000001", IDUsedBy: "1000003", UsedBy: "UJI COMM", Contract: "1", Year: "2026", Comm: "12.5"},
		{ID: "1000002", IDUsedBy: "1000003", UsedBy: "UJI COMM", Contract: "5", Year: "2021", Comm: "0.5"},
		{ID: "1000005", IDUsedBy: "1000003", UsedBy: "UJI COMM", Contract: "2", Year: "2022", Comm: ""},
	}
	if len(tulis) != len(mau) {
		t.Fatalf("tulis %+v", tulis)
	}
	for i := range mau {
		if tulis[i] != mau[i] {
			t.Errorf("tulis[%d] %+v mau %+v", i, tulis[i], mau[i])
		}
	}
	if lap.Normalisasi["CONTRACT"] != 1 || lap.Normalisasi["COMM"] != 1 || lap.Kosong["COMM"] != 1 {
		t.Errorf("normalisasi %v kosong %v", lap.Normalisasi, lap.Kosong)
	}
	if len(lap.Gagal) != 2 || lap.Gagal[0] != "1000004 COMM: tidak muat tipe kolom flat" || lap.Gagal[1] != "1000008 JSONDATA: bukan objek JSON" {
		t.Errorf("gagal %q", lap.Gagal)
	}
	if strings.Contains(lap.Teks(), "RAHASIA-UJI") {
		t.Error("laporan memuat nilai")
	}
	if lap.BolehDitulis() {
		t.Error("gagal/berbeda harus menahan")
	}
	bersih, _ := RencanaPindah(sumber[:2], nil, false)
	if bersih.BolehDitulis() {
		t.Error("normalisasi belum diterima harus menahan")
	}
	if diterima, _ := RencanaPindah(sumber[:2], nil, true); !diterima.BolehDitulis() {
		t.Errorf("normalisasi diterima %+v", diterima)
	}
	if kosong, tulis := RencanaPindah(nil, nil, false); !kosong.BolehDitulis() || len(tulis) != 0 || kosong.Sumber != 0 {
		t.Errorf("sumber kosong (M_RICOMM_LIFE 0 baris DEV) %+v", kosong)
	}
	if ganda, _ := RencanaPindah([]BarisJSON{sumber[0], sumber[0]}, nil, false); len(ganda.Gagal) != 1 ||
		!strings.Contains(ganda.Gagal[0], "kembar") {
		t.Errorf("ID kembar %q", ganda.Gagal)
	}
}
