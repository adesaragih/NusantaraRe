package models

import "testing"

func TestUrutkanJenisReasuransiSesuaiWorkOwner(t *testing.T) {
	masuk := []JenisReasuransi{
		{ID: "5", Note: "OR"}, {ID: "9", Note: "LAINNYA"}, {ID: "3", Note: "Surplus"},
		{ID: "1", Note: " QS "}, {ID: "4", Note: "2ND SURPLUS"}, {ID: "2", Note: "2nd QS"},
	}
	got := UrutkanJenisReasuransi(masuk)
	mau := []string{"1", "2", "3", "4", "5", "9"}
	for i, id := range mau {
		if got[i].ID != id {
			t.Fatalf("urutan %v, mau ID %v", got, mau)
		}
	}
	if masuk[0].ID != "5" {
		t.Error("masukan ikut terubah")
	}
}

func TestJenisKontrakGanda(t *testing.T) {
	daftar := []Kontrak{{ID: "K1", ReinsTypeID: "10196"}, {ID: "K2", ReinsTypeID: "10200"}}
	if !JenisKontrakGanda(daftar, "", "10196") {
		t.Error("kontrak baru dengan jenis yang sudah ada tidak terdeteksi ganda")
	}
	if JenisKontrakGanda(daftar, "K1", "10196") {
		t.Error("mengubah kontrak tanpa mengganti jenisnya dianggap ganda")
	}
	if !JenisKontrakGanda(daftar, "K1", "10200") {
		t.Error("mengganti jenis ke jenis milik kontrak lain tidak terdeteksi")
	}
	if JenisKontrakGanda(daftar, "", "10999") {
		t.Error("jenis baru dianggap ganda")
	}
}

func TestUrutkanKontrakMenurutJenis(t *testing.T) {
	masuk := []Kontrak{
		{ID: "1000055", ReinsTypeName: "SURPLUS"}, {ID: "1000056", ReinsTypeName: "QS"},
		{ID: "1000057", ReinsTypeName: "OR"}, {ID: "1000058", ReinsTypeName: "2ND QS"},
	}
	got := UrutkanKontrakMenurutJenis(masuk)
	mau := []string{"1000056", "1000058", "1000055", "1000057"}
	for i, id := range mau {
		if got[i].ID != id {
			t.Fatalf("urutan %v, mau %v", got, mau)
		}
	}
}

func TestBusinessGanda(t *testing.T) {
	daftar := []Business{{ID: "B1", BizCode: "10172"}, {ID: "B2", BizCode: "10180"}}
	if !BusinessGanda(daftar, "", "10172") {
		t.Error("business baru dengan nama yang sudah ada tidak terdeteksi ganda")
	}
	if BusinessGanda(daftar, "B1", "10172") {
		t.Error("mengubah business tanpa mengganti namanya dianggap ganda")
	}
	if !BusinessGanda(daftar, "B1", "10180") {
		t.Error("mengganti ke nama milik business lain tidak terdeteksi")
	}
}

func TestReinsurerGanda(t *testing.T) {
	daftar := []Reinsurer{{ID: "R1", ReinsurerID: "L01"}, {ID: "R2", ReinsurerID: "L02"}}
	if !ReinsurerGanda(daftar, "", "L01") {
		t.Error("reinsurer baru yang sudah ada tidak terdeteksi ganda")
	}
	if ReinsurerGanda(daftar, "R1", "L01") {
		t.Error("mengubah reinsurer tanpa mengganti namanya dianggap ganda")
	}
	if !ReinsurerGanda(daftar, "R1", "L02") {
		t.Error("mengganti ke reinsurer milik baris lain tidak terdeteksi")
	}
}
