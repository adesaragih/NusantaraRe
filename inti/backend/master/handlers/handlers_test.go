package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/master/models"
	"nusantarare/inti/backend/master/repository"
	"nusantarare/inti/backend/master/services"
)

// tiruan - penyimpan master di memori: tabel -> ID -> baris (+ status).
type tiruan struct {
	baris  map[string]map[string]models.Baris
	aktif  map[string]map[string]bool
	urutan int
}

func baruTiruan() *tiruan {
	t := &tiruan{baris: map[string]map[string]models.Baris{}, aktif: map[string]map[string]bool{}}
	t.isi("NATION", models.Baris{"id": "INA", "note": "INDONESIA"})
	t.isi("ACCUMULATEDTYPE", models.Baris{"id": "AT1", "accumulationType": "UJI TIPE"})
	t.isi("PROVINCE", models.Baris{"id": "P1", "note": "UJI PROVINSI"})
	return t
}

func (t *tiruan) isi(tabel string, b models.Baris) {
	if t.baris[tabel] == nil {
		t.baris[tabel], t.aktif[tabel] = map[string]models.Baris{}, map[string]bool{}
	}
	t.baris[tabel][b["id"]], t.aktif[tabel][b["id"]] = b, true
}

func (t *tiruan) Daftar(_ context.Context, m models.TabelMaster, kata, status string, nomor, ukuran int) (models.Halaman, error) {
	h := models.Halaman{Nomor: nomor, Ukuran: ukuran, Baris: []models.Baris{}, Aktif: []bool{}}
	for id, b := range t.baris[m.Nama] {
		if status == "aktif" && !t.aktif[m.Nama][id] || status == "nonaktif" && t.aktif[m.Nama][id] {
			continue
		}
		h.Baris, h.Aktif = append(h.Baris, b), append(h.Aktif, t.aktif[m.Nama][id])
	}
	h.Total = len(h.Baris)
	return h, nil
}

func (t *tiruan) Ada(_ context.Context, m models.TabelMaster, id string) (bool, error) {
	_, ada := t.baris[m.Nama][id]
	return ada, nil
}

func (t *tiruan) AdaCatatanAkumulasi(_ context.Context, note, zip string) (bool, error) {
	for _, b := range t.baris["ACCUMULATION"] {
		if strings.EqualFold(b["note"], note) && b["zipCode"] == zip {
			return true, nil
		}
	}
	return false, nil
}

func (t *tiruan) NilaiRujukan(_ context.Context, r models.Rujukan, nilai string) (string, bool, error) {
	for _, b := range t.baris[r.Tabel] {
		if b["id"] == nilai {
			return map[string]string{"NOTE": b["note"], "ACCUMULATIONTYPE": b["accumulationType"]}[r.KolomNilai], true, nil
		}
	}
	return "", false, nil
}

func (t *tiruan) IDAkumulasi(_ context.Context, _ *db.Tx, negara, zip string) (string, error) {
	t.urutan++
	return negara + "-" + zip + "-00000" + string(rune('0'+t.urutan)), nil
}

func (t *tiruan) Sisip(_ context.Context, _ *db.Tx, m models.TabelMaster, b models.Baris, akun string) error {
	b["createOp"], b["tglCreate"] = akun, "UJI-TGL"
	t.isi(m.Nama, b)
	return nil
}

func (t *tiruan) Ubah(_ context.Context, _ *db.Tx, m models.TabelMaster, b models.Baris, akun string) error {
	lama, ada := t.baris[m.Nama][b["id"]]
	if !ada {
		return repository.ErrBarisTidakAda
	}
	b["createOp"], b["updateOp"] = lama["createOp"], akun
	t.baris[m.Nama][b["id"]] = b
	return nil
}

func (t *tiruan) UbahStatus(_ context.Context, _ *db.Tx, m models.TabelMaster, id string, aktif bool, akun string) error {
	b, ada := t.baris[m.Nama][id]
	if !ada {
		return repository.ErrBarisTidakAda
	}
	t.aktif[m.Nama][id], b["updateOp"] = aktif, akun
	return nil
}

func tanpaOracle(_ context.Context, fn func(tx *db.Tx) error) error { return fn(nil) }

func minta(t *testing.T, svc *services.Service, metode, jalur, badan, pelaku string) (int, string) {
	t.Helper()
	mux := http.NewServeMux()
	for _, m := range models.DaftarMaster {
		Pasang(mux, "/uji/"+m.Kunci, m.Kunci, svc, true)
	}
	r := httptest.NewRequest(metode, jalur, strings.NewReader(badan))
	if pelaku != "" {
		r.Header.Set("X-Pelaku", pelaku)
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	return w.Code, strings.TrimSpace(w.Body.String())
}

// TestMetaMaster - GET <prefix>/meta setiap master: kunci, sifat, kolom diakhiri keempat kolom jejak (turunan, MD-7),
// dan `rujukan` (selalu larik) = kolom berujukan master, persis RUJUKAN / KOLOM_NAMA frontend lama.
func TestMetaMaster(t *testing.T) {
	type meta struct {
		Kunci      string `json:"kunci"`
		IDOtomatis bool   `json:"idOtomatis"`
		Kolom      []struct {
			Kunci   string `json:"kunci"`
			Turunan bool   `json:"turunan"`
		} `json:"kolom"`
		Rujukan []struct {
			Kunci, Judul, Nilai, Nama string
		} `json:"rujukan"`
	}
	mau := map[string]string{
		"nation":          "",
		"province":        "nationId>Nation:id/note",
		"city":            "provinceId>Province:id/note",
		"district":        "cityId>City:id/note",
		"czone":           "groupOf>CZone:code/description",
		"accumulatedtype": "",
		"accumulation":    "accumulation>Accumulated Type:id/accumulationType,provinceId>Province:id/note",
		"objectitemtype":  "",
	}
	for _, m := range models.DaftarMaster {
		kode, isi := minta(t, services.Baru(nil, nil), "GET", "/uji/"+m.Kunci+"/meta", "", "")
		var j meta
		if err := json.Unmarshal([]byte(isi), &j); err != nil || kode != 200 || j.Kunci != m.Kunci || j.Rujukan == nil {
			t.Fatalf("%s: %d %s", m.Kunci, kode, isi)
		}
		n := len(j.Kolom)
		if n < 4 || j.Kolom[n-4].Kunci != "createOp" || j.Kolom[n-1].Kunci != "tglUpdate" || !j.Kolom[n-3].Turunan {
			t.Errorf("%s tanpa kolom jejak: %+v", m.Kunci, j.Kolom)
		}
		var r []string
		for _, x := range j.Rujukan {
			r = append(r, x.Kunci+">"+x.Judul+":"+x.Nilai+"/"+x.Nama)
		}
		if strings.Join(r, ",") != mau[m.Kunci] {
			t.Errorf("%s rujukan %v, mau %s", m.Kunci, r, mau[m.Kunci])
		}
		if j.IDOtomatis != (m.Kunci == "accumulation") {
			t.Errorf("%s idOtomatis %v", m.Kunci, j.IDOtomatis)
		}
	}
	if len(mau) != len(models.DaftarMaster) {
		t.Errorf("%d master, mau %d", len(models.DaftarMaster), len(mau))
	}
	if _, isi := minta(t, services.Baru(nil, nil), "GET", "/uji/province/meta", "", ""); !strings.Contains(isi, `"turunan":true`) {
		t.Errorf("province nationName turunan: %s", isi)
	}
}

// TestRujukanMaster - GET <prefix>/rujukan/{kolom}: baris master yang dirujuk, AKTIF saja; kolom tanpa rujukan master
// (BRANCHID - isian teks) dan kolom asing -> 404.
func TestRujukanMaster(t *testing.T) {
	tr := baruTiruan()
	tr.isi("PROVINCE", models.Baris{"id": "P9", "note": "NONAKTIF"})
	tr.aktif["PROVINCE"]["P9"] = false
	svc := services.Baru(tr, tanpaOracle)
	kode, isi := minta(t, svc, "GET", "/uji/city/rujukan/provinceId?q=uji", "", "")
	if kode != 200 || !strings.Contains(isi, `"id":"P1"`) || strings.Contains(isi, "P9") || !strings.Contains(isi, `"aktif":true`) {
		t.Errorf("rujukan province: %d %s", kode, isi)
	}
	if kode, isi := minta(t, svc, "GET", "/uji/province/rujukan/nationId", "", ""); kode != 200 || !strings.Contains(isi, "INDONESIA") {
		t.Errorf("rujukan nation: %d %s", kode, isi)
	}
	for _, j := range []string{"/uji/city/rujukan/branchId", "/uji/city/rujukan/kota", "/uji/nation/rujukan/id"} {
		if kode, _ := minta(t, svc, "GET", j, "", ""); kode != 404 {
			t.Errorf("%s: %d, mau 404", j, kode)
		}
	}
	if kode, _ := minta(t, svc, "GET", "/uji/city/rujukan/provinceId?halaman=x", "", ""); kode != 400 {
		t.Errorf("halaman bukan angka: %d", kode)
	}
}

// TestTambahUbahStatusMaster - alur tambah / ubah / aktif-nonaktif + galat: 401 tanpa identitas, 404 master / baris,
// 400 wajib / medan asing / rujukan tak ada / lebar, 409 ganda, turunan diisi backend (masukan turunan diabaikan).
func TestTambahUbahStatusMaster(t *testing.T) {
	tr := baruTiruan()
	svc := services.Baru(tr, tanpaOracle)
	kode, isi := minta(t, svc, "POST", "/uji/province", `{"id":" P2 ","nationId":"INA","note":"UJI JAWA","nationName":"PALSU"}`, "UJI-USER")
	if kode != 201 || isi != `{"id":"P2"}` || tr.baris["PROVINCE"]["P2"]["nationName"] != "INDONESIA" || tr.baris["PROVINCE"]["P2"]["createOp"] != "UJI-USER" {
		t.Fatalf("tambah: %d %s %v", kode, isi, tr.baris["PROVINCE"]["P2"])
	}
	for nama, u := range map[string]struct {
		metode, jalur, badan, pelaku, pesan string
		kode                                int
	}{
		"tanpa identitas": {"POST", "/uji/province", `{"id":"P3","note":"X"}`, "", "", 401},
		"master asing":    {"POST", "/uji/kota", `{"id":"P3"}`, "UJI-USER", "", 404},
		"ID ganda":        {"POST", "/uji/province", `{"id":"P2","note":"X"}`, "UJI-USER", "sudah ada", 409},
		"wajib":           {"POST", "/uji/province", `{"id":"P3"}`, "UJI-USER", "note wajib diisi", 400},
		"medan asing":     {"POST", "/uji/province", `{"id":"P3","note":"X","kota":"Y"}`, "UJI-USER", "kota bukan medan master province", 400},
		"rujukan":         {"POST", "/uji/province", `{"id":"P3","note":"X","nationId":"ZZ"}`, "UJI-USER", `nationId \"ZZ\" tidak ada di NATION`, 400},
		"lebar":           {"POST", "/uji/nation", `{"id":"12345678901","note":"X"}`, "UJI-USER", "id paling banyak 10 byte", 400},
		"bukan teks":      {"POST", "/uji/nation", `{"id":1}`, "UJI-USER", "bernilai teks", 400},
		"ubah tak ada":    {"PUT", "/uji/province/P9", `{"note":"X"}`, "UJI-USER", "tidak ada", 404},
		"status tak ada":  {"PUT", "/uji/province/P9/status", `{"aktif":false}`, "UJI-USER", "tidak ada", 404},
		"status badan":    {"PUT", "/uji/province/P2/status", `{"aktif":"tidak"}`, "UJI-USER", "aktif", 400},
		"halaman":         {"GET", "/uji/province?halaman=0", "", "", "halaman", 400},
		"pelaku panjang":  {"POST", "/uji/province", `{"id":"P3","note":"X"}`, strings.Repeat("U", 65), "akun pelaku paling banyak 64 byte", 400},
		"jejak dikirim":   {"POST", "/uji/province", `{"id":"P4","note":"X","createOp":"PALSU"}`, "UJI-USER", "", 201},
		"status daftar":   {"GET", "/uji/province?status=x", "", "", "status", 400},
	} {
		if kode, isi := minta(t, svc, u.metode, u.jalur, u.badan, u.pelaku); kode != u.kode || !strings.Contains(isi, u.pesan) {
			t.Errorf("%s: %d %s, mau %d", nama, kode, isi, u.kode)
		}
	}
	if kode, isi := minta(t, svc, "PUT", "/uji/province/P2", `{"id":"ABAIKAN","note":"UJI JAWA BARAT","nationId":""}`, "UJI-USER"); kode != 200 ||
		isi != `{"id":"P2"}` || tr.baris["PROVINCE"]["P2"]["note"] != "UJI JAWA BARAT" || tr.baris["PROVINCE"]["P2"]["nationName"] != "" ||
		tr.baris["PROVINCE"]["P2"]["updateOp"] != "UJI-USER" {
		t.Errorf("ubah: %d %s %v", kode, isi, tr.baris["PROVINCE"]["P2"])
	}
	if kode, isi := minta(t, svc, "PUT", "/uji/province/P2/status", `{"aktif":false}`, "UJI-USER"); kode != 200 || isi != `{"id":"P2","aktif":false}` {
		t.Errorf("status: %d %s", kode, isi)
	}
	kode, isi = minta(t, svc, "GET", "/uji/province?status=nonaktif", "", "")
	if kode != 200 || !strings.Contains(isi, `"aktif":false`) || !strings.Contains(isi, `"id":"P2"`) || !strings.Contains(isi, `"total":1,"halaman":1,"ukuran":20`) {
		t.Errorf("daftar nonaktif: %d %s", kode, isi)
	}
	if kode, _ := minta(t, services.Baru(nil, nil), "GET", "/uji/province", "", ""); kode != 503 {
		t.Errorf("503: %d", kode)
	}
}

// TestTambahAkumulasi - ID dibuat backend dari negara + zip (prosedur RDBMASTERACCUMULATION), NOTE huruf besar,
// Note + zip ganda 409, negara wajib, turunan ACCUMULATIONNAME / PROVINCE dari rujukan.
func TestTambahAkumulasi(t *testing.T) {
	tr := baruTiruan()
	svc := services.Baru(tr, tanpaOracle)
	badan := `{"negara":"INA","zipCode":"61151","note":"jl uji 1","accumulation":"AT1","provinceId":"P1"}`
	kode, isi := minta(t, svc, "POST", "/uji/accumulation", badan, "UJI-USER")
	b := tr.baris["ACCUMULATION"]["INA-61151-000001"]
	if kode != 201 || isi != `{"id":"INA-61151-000001"}` || b["note"] != "JL UJI 1" || b["accumulationName"] != "UJI TIPE" || b["province"] != "UJI PROVINSI" ||
		b["negara"] != "" {
		t.Fatalf("tambah: %d %s %v", kode, isi, b)
	}
	if kode, isi := minta(t, svc, "POST", "/uji/accumulation", strings.Replace(badan, "jl uji 1", "JL UJI 1", 1), "UJI-USER"); kode != 409 ||
		!strings.Contains(isi, "Akumulasi sudah ada") {
		t.Errorf("ganda: %d %s", kode, isi)
	}
	if kode, isi := minta(t, svc, "POST", "/uji/accumulation", `{"zipCode":"1","note":"X"}`, "UJI-USER"); kode != 400 || !strings.Contains(isi, "negara wajib diisi") {
		t.Errorf("negara: %d %s", kode, isi)
	}
	if kode, isi := minta(t, svc, "PUT", "/uji/accumulation/INA-61151-000001", `{"negara":"INA","zipCode":"1","note":"X"}`, "UJI-USER"); kode != 400 ||
		!strings.Contains(isi, "negara bukan medan master") {
		t.Errorf("negara saat ubah: %d %s", kode, isi)
	}
}
