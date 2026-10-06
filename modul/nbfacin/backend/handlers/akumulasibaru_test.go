package handlers

import (
	"context"
	"fmt"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	masterservices "nusantarare/inti/backend/master/services"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

// addTiruan - pembaca form Add New: mencatat parameter.
type addTiruan struct{ catat *string }

func (a addTiruan) CZoneZip(_ context.Context, zip string) (models.CZoneZip, bool, error) {
	*a.catat = "czone:" + zip
	if zip == "00000" {
		return models.CZoneZip{}, false, nil
	}
	return models.CZoneZip{ID: "Z9", Code: "UJI-ZONA"}, true, nil
}

func (a addTiruan) CariZip(_ context.Context, provinsi, kata string, nomor, ukuran int) ([]models.ZipAkumulasi, int, error) {
	*a.catat = fmt.Sprintf("zip:%s|%s|%d|%d", provinsi, kata, nomor, ukuran)
	if nomor > 1 {
		return []models.ZipAkumulasi{}, 31, nil
	}
	return []models.ZipAkumulasi{{ZipCode: "40111", City: "UJI KOTA", Province: "UJI PROV", Nation: "UJI NEGARA", NationInitial: "UJI"}}, 31, nil
}

// penambahTiruan - mesin master: mencatat masukan, menjawab galat yang disetel.
type penambahTiruan struct {
	masukan map[string]string
	kunci   string
	galat   error
}

func (p *penambahTiruan) Tambah(_ context.Context, pl inti.Pelaku, kunci string, m map[string]string) (string, error) {
	if err := inti.WajibIdentitas(pl); err != nil {
		return "", err
	}
	p.kunci, p.masukan = kunci, m
	return "UJI-40111-000001", p.galat
}

// TestTambahAkumulasi - POST /api/nbfacin/akumulasi: wajib SaveAccumulation_Act (pesan verbatim Pega) lebih dulu,
// lalu wajib section; disimpan lewat master "accumulation" dengan kunci master + negara; ganda -> 409 pesan prosedur;
// galat master -> 400 / 503; jawaban {id, note huruf besar}.
func TestTambahAkumulasi(t *testing.T) {
	var catat string
	p := &penambahTiruan{}
	svc := services.Baru(nil).DenganAddAkumulasi(addTiruan{&catat}, p)
	penuh := `{"accumulation":"AT1","accumulationType":"UJI TIPE","note":" jl uji 1 ","keyword":"K","scopeArea":"CITY",` +
		`"cZone":"UJI-ZONA","cZoneId":"Z9","provinceId":"P1","zipCode":"40111","negara":"UJI"}`
	kode, isi := minta(t, svc, "POST", "/api/nbfacin/akumulasi", penuh, "UJI-USER")
	if kode != 201 || isi != `{"id":"UJI-40111-000001","note":"JL UJI 1"}` || p.kunci != "accumulation" ||
		p.masukan["note"] != "jl uji 1" || p.masukan[masterservices.KunciNegara] != "UJI" || p.masukan["cZoneId"] != "Z9" ||
		p.masukan["accumulationType"] != "UJI TIPE" || len(p.masukan) != 10 {
		t.Fatalf("%d %s %v", kode, isi, p.masukan)
	}
	pega := `{"galat":"Postal code, Nation, CZone and Accumulation Description can't be null!"}`
	for nama, k := range map[string]struct {
		badan, pelaku, mau string
		kode               int
	}{
		"zip kosong":    {strings.Replace(penuh, `"40111"`, `" "`, 1), "UJI-USER", pega, 400},
		"negara kosong": {strings.Replace(penuh, `"negara":"UJI"`, `"negara":""`, 1), "UJI-USER", pega, 400},
		"czone kosong":  {strings.Replace(penuh, `"cZone":"UJI-ZONA"`, `"cZone":""`, 1), "UJI-USER", pega, 400},
		"note kosong":   {strings.Replace(penuh, `" jl uji 1 "`, `""`, 1), "UJI-USER", pega, 400},
		"pega dulu":     {`{"keyword":"K"}`, "UJI-USER", pega, 400},
		"wajib section": {strings.Replace(strings.Replace(penuh, `"keyword":"K"`, `"keyword":""`, 1), `"provinceId":"P1"`, `"provinceId":""`, 1),
			"UJI-USER", "provinceId wajib diisi; keyword wajib diisi", 400},
		"kunci asing":     {strings.Replace(penuh, `"negara"`, `"kota"`, 1), "UJI-USER", "objek JSON akumulasi", 400},
		"scope di luar":   {strings.Replace(penuh, `"scopeArea":"CITY"`, `"scopeArea":"area"`, 1), "UJI-USER", "AREA, DISTRICT, CITY, PROVINCE, COUNTRY", 400},
		"tanpa identitas": {penuh, "", "", 401},
	} {
		if kode, isi := minta(t, svc, "POST", "/api/nbfacin/akumulasi", k.badan, k.pelaku); kode != k.kode || !strings.Contains(isi, k.mau) {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
	for nama, k := range map[string]struct {
		galat error
		mau   string
		kode  int
	}{
		"ganda":       {fmt.Errorf("%w: master item", masterservices.ErrSudahAda), `{"galat":"Error master item Akumulasi Sudah Ada"}`, 409},
		"rujukan":     {fmt.Errorf("%w: provinceId tidak ada di PROVINCE", masterservices.ErrMasukanMaster), "provinceId tidak ada di PROVINCE", 400},
		"tanpa basis": {masterservices.ErrMasterTanpaDatabase, "akumulasi tidak terbaca", 503},
	} {
		p.galat = k.galat
		if kode, isi := minta(t, svc, "POST", "/api/nbfacin/akumulasi", penuh, "UJI-USER"); kode != k.kode || !strings.Contains(isi, k.mau) {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
	if kode, _ := minta(t, services.Baru(nil), "POST", "/api/nbfacin/akumulasi", penuh, "UJI-USER"); kode != 503 {
		t.Errorf("tanpa penambah: %d", kode)
	}
}

// TestCZoneDanZipAkumulasi - GET .../czone?zip= (kosong / tanpa CZone = kosong) dan GET .../zipcode?provinceName=&q=.
func TestCZoneDanZipAkumulasi(t *testing.T) {
	var catat string
	svc := services.Baru(nil).DenganAddAkumulasi(addTiruan{&catat}, &penambahTiruan{})
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/akumulasi/czone?zip=%2040111%20", "", ""); kode != 200 ||
		isi != `{"cZoneId":"Z9","cZone":"UJI-ZONA"}` || catat != "czone:40111" {
		t.Errorf("czone: %d %s %q", kode, isi, catat)
	}
	for _, zip := range []string{"", "00000"} {
		if kode, isi := minta(t, svc, "GET", "/api/nbfacin/akumulasi/czone?zip="+zip, "", ""); kode != 200 || isi != `{"cZoneId":"","cZone":""}` {
			t.Errorf("czone %q: %d %s", zip, kode, isi)
		}
	}
	kode, isi := minta(t, svc, "GET", "/api/nbfacin/akumulasi/zipcode?provinceName=jawa&q=401", "", "")
	if kode != 200 || catat != "zip:jawa|401|1|15" ||
		isi != `{"baris":[{"zipCode":"40111","city":"UJI KOTA","province":"UJI PROV","nation":"UJI NEGARA","nationInitial":"UJI"}],"total":31,"halaman":1,"ukuran":15}` {
		t.Errorf("zip: %d %s %q", kode, isi, catat)
	}
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/akumulasi/zipcode?halaman=9", "", ""); kode != 200 ||
		isi != `{"baris":[],"total":31,"halaman":9,"ukuran":15}` || catat != "zip:||9|15" {
		t.Errorf("di luar jangkauan: %d %s %q", kode, isi, catat)
	}
	for _, h := range []string{"0", "x", "-1"} {
		if kode, _ := minta(t, svc, "GET", "/api/nbfacin/akumulasi/zipcode?halaman="+h, "", ""); kode != 400 {
			t.Errorf("halaman %s: %d", h, kode)
		}
	}
	if kode, _ := minta(t, svc, "GET", "/api/nbfacin/akumulasi/zipcode?q="+strings.Repeat("9", 256), "", ""); kode != 400 {
		t.Errorf("q panjang: %d", kode)
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/akumulasi/zipcode", "", ""); kode != 503 {
		t.Errorf("tanpa basis data: %d", kode)
	}
}
