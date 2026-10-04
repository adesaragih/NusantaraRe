package services_test

// Uji baca satu kontrak WARISAN — nol koneksi Oracle.
//
// Bentuk data ujinya dari sapuan 3 Oktober 2026 atas seluruh 1.854 dokumen
// `M_TREATY_IN.JSONDATA`, bukan dikarang:
//
//	Bordeaux        1.851/1.854 — `reporting` 1.164 · `nonreporting` 687
//	AccountingMode  1.851       — `underwriting` 1.110 · `accounting` 741
//	BordereauxNote  1.018
//	ContractRefNo     742
//	TreatyLeader      659       — nilainya teks `"true"` / `"false"`

import (
	"context"
	"errors"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

func kontrakWarisanUji() models.KontrakWarisan {
	return models.KontrakWarisan{
		ID:                  "1000797",
		NamaKontrak:         "PROPERTY QUOTA SHARE",
		LingkupWilayah:      "Republic of Indonesia",
		TahunTreaty:         "2019",
		Cedant:              "SAHABAT INSURANCE",
		AsalBisnis:          "IBS REINSURANCE BROKERS",
		SifatProporsiAsli:   "NonProportional",
		TanggalMulaiAsli:    "20190101",
		TanggalBerakhirAsli: "20191231",
		Bordereaux:          "reporting",
		BordereauxCatatan:   "Not Reporting",
		CaraPembukuan:       "underwriting",
		NomorRujukan:        "N0XM000221",
		PemimpinTreaty:      "false",
		AdaDiJSON: map[string]bool{
			"Bordeaux": true, "BordereauxNote": true, "AccountingMode": true,
			"ContractRefNo": true, "TreatyLeader": true,
		},
	}
}

func TestBacaWarisanMenolakTanpaIdentitasSebelumGudang(t *testing.T) {
	g := &gudangTiruan{kontrakWarisan: kontrakWarisanUji()}
	l := services.LayananDengan(g)

	if _, err := l.BacaKontrakWarisan(context.Background(), inti.Pelaku{}, "1000797"); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Fatalf("mau ErrTanpaIdentitas, dapat %v", err)
	}
}

func TestBacaWarisanMenolakPengenalKosong(t *testing.T) {
	g := &gudangTiruan{kontrakWarisan: kontrakWarisanUji()}
	l := services.LayananDengan(g)

	for _, id := range []string{"", "   "} {
		if _, err := l.BacaKontrakWarisan(context.Background(), pelakuAda, id); !errors.Is(err, services.ErrMasukanTidakSah) {
			t.Errorf("id %q: mau ErrMasukanTidakSah, dapat %v", id, err)
		}
	}
}

// Uji POSITIF: ketiga medan JSON terisi, dan terjemahannya MEMAKAI ULANG
// penerjemah yang sama dengan layar daftar.
func TestBacaWarisanJSONLengkap(t *testing.T) {
	g := &gudangTiruan{kontrakWarisan: kontrakWarisanUji()}
	l := services.LayananDengan(g)

	k, err := l.BacaKontrakWarisan(context.Background(), pelakuAda, "1000797")
	if err != nil {
		t.Fatalf("mau diterima, dapat %v", err)
	}
	for _, p := range []struct{ nama, dapat, mau string }{
		{"Bordeaux", k.Bordereaux, "reporting"},
		{"BordereauxNote", k.BordereauxCatatan, "Not Reporting"},
		{"AccountingMode", k.CaraPembukuan, "underwriting"},
		{"ContractRefNo", k.NomorRujukan, "N0XM000221"},
		{"TreatyLeader", k.PemimpinTreaty, "false"},
		// ⭐ Terjemahan yang SAMA dengan layar daftar — satu penerjemah.
		{"sifat proporsi", k.SifatProporsi, "Non Proportional"},
		{"tanggal mulai", k.TanggalMulai, "01/01/19"},
		{"tanggal berakhir", k.TanggalBerakhir, "31/12/19"},
	} {
		if p.dapat != p.mau {
			t.Errorf("%s = %q, mau %q", p.nama, p.dapat, p.mau)
		}
	}
	// Nilai aslinya tidak hilang.
	if k.SifatProporsiAsli != "NonProportional" || k.TanggalMulaiAsli != "20190101" {
		t.Errorf("nilai asli hilang: %+v", k)
	}
}

// ⭐ Kunci yang HILANG bukan galat — ia medan yang kontrak itu memang tidak
// punya. 1.195 dari 1.854 kontrak tidak punya `TreatyLeader`, dan 1.112
// tidak punya `ContractRefNo`; menolak membukanya akan menolak sebagian
// besar tabel.
func TestKunciJSONHilangBukanGalat(t *testing.T) {
	k := kontrakWarisanUji()
	k.Bordereaux = ""
	k.NomorRujukan = ""
	k.PemimpinTreaty = ""
	k.AdaDiJSON = map[string]bool{"BordereauxNote": true, "AccountingMode": true}
	g := &gudangTiruan{kontrakWarisan: k}
	l := services.LayananDengan(g)

	hasil, err := l.BacaKontrakWarisan(context.Background(), pelakuAda, "1000797")
	if err != nil {
		t.Fatalf("kunci hilang menghasilkan galat %v; ia keadaan yang sah", err)
	}
	if hasil.Bordereaux != "" {
		t.Errorf("Bordeaux terisi %q padahal kuncinya tidak ada", hasil.Bordereaux)
	}
	// ⛔ Dan BEDANYA terbaca: `AdaDiJSON` memisahkan "kuncinya tidak ada"
	// dari "kuncinya ada bernilai kosong". Keduanya berarti hal yang
	// berbeda di layar.
	if hasil.AdaDiJSON["Bordeaux"] {
		t.Error("AdaDiJSON mengaku Bordeaux ada")
	}
	if !hasil.AdaDiJSON["AccountingMode"] {
		t.Error("AdaDiJSON kehilangan AccountingMode yang ada")
	}
}

// Dokumen rusak -> galat yang MENYEBUT kontraknya, bukan panik.
func TestJSONRusakMenghasilkanGalatYangMenyebutKontraknya(t *testing.T) {
	g := &gudangTiruan{galat: services.ErrJSONWarisanRusak}
	l := services.LayananDengan(g)

	_, err := l.BacaKontrakWarisan(context.Background(), pelakuAda, "1000797")
	if !services.WarisanJSONRusak(err) {
		t.Fatalf("mau ErrJSONWarisanRusak, dapat %v", err)
	}
}

// Pengenal yang tidak ada -> 404, bukan form kosong.
func TestWarisanTidakAdaDikenaliSebagai404(t *testing.T) {
	g := &gudangTiruan{galat: services.ErrWarisanTidakAda}
	l := services.LayananDengan(g)

	_, err := l.BacaKontrakWarisan(context.Background(), pelakuAda, "9999999")
	if !services.WarisanTidakAda(err) {
		t.Fatalf("mau ErrWarisanTidakAda, dapat %v", err)
	}
	// Dan ia BUKAN galat JSON - kedua pembedanya tidak boleh tertukar.
	if services.WarisanJSONRusak(err) {
		t.Error("galat 'tidak ada' terbaca sebagai 'JSON rusak'")
	}
}

// `AdaDiJSON` tidak pernah nil, bahkan ketika dokumennya tidak ada sama
// sekali - pemanggil tidak perlu menjaga diri dari peta nil.
func TestAdaDiJSONTidakPernahNil(t *testing.T) {
	k := kontrakWarisanUji()
	k.AdaDiJSON = nil
	g := &gudangTiruan{kontrakWarisan: k}
	l := services.LayananDengan(g)

	hasil, err := l.BacaKontrakWarisan(context.Background(), pelakuAda, "1000797")
	if err != nil {
		t.Fatal(err)
	}
	if hasil.AdaDiJSON == nil {
		t.Error("AdaDiJSON nil")
	}
}
