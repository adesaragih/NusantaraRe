package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterdata/backend/models"
	"nusantarare/modul/masterdata/backend/repository"
	"nusantarare/modul/masterdata/backend/services"
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

func (t *tiruan) Sisip(_ context.Context, _ *db.Tx, m models.TabelMaster, b models.Baris) error {
	t.isi(m.Nama, b)
	return nil
}

func (t *tiruan) Ubah(_ context.Context, _ *db.Tx, m models.TabelMaster, b models.Baris) error {
	if _, ada := t.baris[m.Nama][b["id"]]; !ada {
		return repository.ErrBarisTidakAda
	}
	t.baris[m.Nama][b["id"]] = b
	return nil
}

func (t *tiruan) UbahStatus(_ context.Context, _ *db.Tx, m models.TabelMaster, id string, aktif bool) error {
	if _, ada := t.baris[m.Nama][id]; !ada {
		return repository.ErrBarisTidakAda
	}
	t.aktif[m.Nama][id] = aktif
	return nil
}

func tanpaOracle(_ context.Context, fn func(tx *db.Tx) error) error { return fn(nil) }

func minta(t *testing.T, svc *services.Service, metode, jalur, badan, pelaku string) (int, string) {
	t.Helper()
	mux := http.NewServeMux()
	DaftarkanRute(mux, svc, true)
	r := httptest.NewRequest(metode, jalur, strings.NewReader(badan))
	if pelaku != "" {
		r.Header.Set("X-Pelaku", pelaku)
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	return w.Code, strings.TrimSpace(w.Body.String())
}

// TestMetadataMaster - GET /api/masterdata: delapan master, kolom dan sifatnya.
func TestMetadataMaster(t *testing.T) {
	kode, isi := minta(t, services.Baru(nil, nil), "GET", "/api/masterdata", "", "")
	var j struct {
		Tabel []struct {
			Kunci      string `json:"kunci"`
			IDOtomatis bool   `json:"idOtomatis"`
			Kolom      []struct {
				Kunci   string `json:"kunci"`
				Turunan bool   `json:"turunan"`
			} `json:"kolom"`
		} `json:"tabel"`
	}
	if err := json.Unmarshal([]byte(isi), &j); err != nil || kode != 200 || len(j.Tabel) != 8 {
		t.Fatalf("%d %s", kode, isi)
	}
	var kunci []string
	for _, x := range j.Tabel {
		kunci = append(kunci, x.Kunci)
	}
	if strings.Join(kunci, ",") != "nation,province,city,district,czone,accumulatedtype,accumulation,objectitemtype" {
		t.Errorf("urutan master %v", kunci)
	}
	if !j.Tabel[6].IDOtomatis || !j.Tabel[1].Kolom[3].Turunan {
		t.Errorf("sifat: %s", isi)
	}
}

// TestTambahUbahStatusMaster - alur tambah / ubah / aktif-nonaktif + galat: 401 tanpa identitas, 404 master / baris,
// 400 wajib / medan asing / rujukan tak ada / lebar, 409 ganda, turunan diisi backend (masukan turunan diabaikan).
func TestTambahUbahStatusMaster(t *testing.T) {
	tr := baruTiruan()
	svc := services.Baru(tr, tanpaOracle)
	kode, isi := minta(t, svc, "POST", "/api/masterdata/province", `{"id":" P2 ","nationId":"INA","note":"UJI JAWA","nationName":"PALSU"}`, "UJI-USER")
	if kode != 201 || isi != `{"id":"P2"}` || tr.baris["PROVINCE"]["P2"]["nationName"] != "INDONESIA" {
		t.Fatalf("tambah: %d %s %v", kode, isi, tr.baris["PROVINCE"]["P2"])
	}
	for nama, u := range map[string]struct {
		metode, jalur, badan, pelaku, pesan string
		kode                                int
	}{
		"tanpa identitas": {"POST", "/api/masterdata/province", `{"id":"P3","note":"X"}`, "", "", 401},
		"master asing":    {"POST", "/api/masterdata/kota", `{"id":"P3"}`, "UJI-USER", "", 404},
		"ID ganda":        {"POST", "/api/masterdata/province", `{"id":"P2","note":"X"}`, "UJI-USER", "sudah ada", 409},
		"wajib":           {"POST", "/api/masterdata/province", `{"id":"P3"}`, "UJI-USER", "note wajib diisi", 400},
		"medan asing":     {"POST", "/api/masterdata/province", `{"id":"P3","note":"X","kota":"Y"}`, "UJI-USER", "kota bukan medan master province", 400},
		"rujukan":         {"POST", "/api/masterdata/province", `{"id":"P3","note":"X","nationId":"ZZ"}`, "UJI-USER", `nationId \"ZZ\" tidak ada di NATION`, 400},
		"lebar":           {"POST", "/api/masterdata/nation", `{"id":"12345678901","note":"X"}`, "UJI-USER", "id paling banyak 10 byte", 400},
		"bukan teks":      {"POST", "/api/masterdata/nation", `{"id":1}`, "UJI-USER", "bernilai teks", 400},
		"ubah tak ada":    {"PUT", "/api/masterdata/province/P9", `{"note":"X"}`, "UJI-USER", "tidak ada", 404},
		"status tak ada":  {"PUT", "/api/masterdata/province/P9/status", `{"aktif":false}`, "UJI-USER", "tidak ada", 404},
		"status badan":    {"PUT", "/api/masterdata/province/P2/status", `{"aktif":"tidak"}`, "UJI-USER", "aktif", 400},
		"halaman":         {"GET", "/api/masterdata/province?halaman=0", "", "", "halaman", 400},
		"status daftar":   {"GET", "/api/masterdata/province?status=x", "", "", "status", 400},
	} {
		if kode, isi := minta(t, svc, u.metode, u.jalur, u.badan, u.pelaku); kode != u.kode || !strings.Contains(isi, u.pesan) {
			t.Errorf("%s: %d %s, mau %d", nama, kode, isi, u.kode)
		}
	}
	if kode, isi := minta(t, svc, "PUT", "/api/masterdata/province/P2", `{"id":"ABAIKAN","note":"UJI JAWA BARAT","nationId":""}`, "UJI-USER"); kode != 200 ||
		isi != `{"id":"P2"}` || tr.baris["PROVINCE"]["P2"]["note"] != "UJI JAWA BARAT" || tr.baris["PROVINCE"]["P2"]["nationName"] != "" {
		t.Errorf("ubah: %d %s %v", kode, isi, tr.baris["PROVINCE"]["P2"])
	}
	if kode, isi := minta(t, svc, "PUT", "/api/masterdata/province/P2/status", `{"aktif":false}`, "UJI-USER"); kode != 200 || isi != `{"id":"P2","aktif":false}` {
		t.Errorf("status: %d %s", kode, isi)
	}
	kode, isi = minta(t, svc, "GET", "/api/masterdata/province?status=nonaktif", "", "")
	if kode != 200 || !strings.Contains(isi, `"aktif":false`) || !strings.Contains(isi, `"id":"P2"`) || !strings.Contains(isi, `"total":1,"halaman":1,"ukuran":50`) {
		t.Errorf("daftar nonaktif: %d %s", kode, isi)
	}
	if kode, _ := minta(t, services.Baru(nil, nil), "GET", "/api/masterdata/province", "", ""); kode != 503 {
		t.Errorf("503: %d", kode)
	}
}

// TestTambahAkumulasi - ID dibuat backend dari negara + zip (prosedur RDBMASTERACCUMULATION), NOTE huruf besar,
// Note + zip ganda 409, negara wajib, turunan ACCUMULATIONNAME / PROVINCE dari rujukan.
func TestTambahAkumulasi(t *testing.T) {
	tr := baruTiruan()
	svc := services.Baru(tr, tanpaOracle)
	badan := `{"negara":"INA","zipCode":"61151","note":"jl uji 1","accumulation":"AT1","provinceId":"P1"}`
	kode, isi := minta(t, svc, "POST", "/api/masterdata/accumulation", badan, "UJI-USER")
	b := tr.baris["ACCUMULATION"]["INA-61151-000001"]
	if kode != 201 || isi != `{"id":"INA-61151-000001"}` || b["note"] != "JL UJI 1" || b["accumulationName"] != "UJI TIPE" || b["province"] != "UJI PROVINSI" ||
		b["negara"] != "" {
		t.Fatalf("tambah: %d %s %v", kode, isi, b)
	}
	if kode, isi := minta(t, svc, "POST", "/api/masterdata/accumulation", strings.Replace(badan, "jl uji 1", "JL UJI 1", 1), "UJI-USER"); kode != 409 ||
		!strings.Contains(isi, "Akumulasi sudah ada") {
		t.Errorf("ganda: %d %s", kode, isi)
	}
	if kode, isi := minta(t, svc, "POST", "/api/masterdata/accumulation", `{"zipCode":"1","note":"X"}`, "UJI-USER"); kode != 400 || !strings.Contains(isi, "negara wajib diisi") {
		t.Errorf("negara: %d %s", kode, isi)
	}
	if kode, isi := minta(t, svc, "PUT", "/api/masterdata/accumulation/INA-61151-000001", `{"negara":"INA","zipCode":"1","note":"X"}`, "UJI-USER"); kode != 400 ||
		!strings.Contains(isi, "negara bukan medan master") {
		t.Errorf("negara saat ubah: %d %s", kode, isi)
	}
}
