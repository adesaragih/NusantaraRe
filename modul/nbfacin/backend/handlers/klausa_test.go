package handlers

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
	"nusantarare/modul/nbfacin/backend/services"
)

// klausaTiruan - ClauseList per kasus di memori; "K404" = case tidak ada.
type klausaTiruan struct {
	simpan map[string][]models.KlausaKasus
	idArg  string
	cari   string
}

func (k *klausaTiruan) BacaKlausa(_ context.Context, id string) ([]models.KlausaKasus, error) {
	if id == "K404" {
		return nil, repository.ErrKasusTidakAda
	}
	if d, ada := k.simpan[id]; ada {
		return d, nil
	}
	return []models.KlausaKasus{}, nil
}

func (k *klausaTiruan) GantiKlausa(_ context.Context, _ *db.Tx, id string, baris []models.KlausaKasus) error {
	if id == "K404" {
		return repository.ErrKasusTidakAda
	}
	k.simpan[id] = baris
	return nil
}

func (k *klausaTiruan) CariKlausa(_ context.Context, bahasa, kata string, nomor, ukuran int) ([]models.HasilKlausa, int, error) {
	k.cari = bahasa + "|" + kata + "|" + strconv.Itoa(nomor) + "|" + strconv.Itoa(ukuran)
	if nomor > 1 {
		return []models.HasilKlausa{}, 11, nil
	}
	return []models.HasilKlausa{{ID: "K01", Title: "UJI JUDUL", Info: "UJI INFO", Text: "UJI ISI", ArgumentCount: "2"}}, 11, nil
}

func (k *klausaTiruan) ArgumenKlausa(_ context.Context, id string) ([]models.ArgumenKlausa, error) {
	k.idArg = id
	if id == "TIADA" { // M_ARGCLAUSEFIRE tidak ada (K47-6): repository menjawab daftar kosong
		return []models.ArgumenKlausa{}, nil
	}
	return []models.ArgumenKlausa{{Number: "1", Description: "UJI NAMA", Value: "UJI NILAI"}}, nil
}

func tanpaOracleKlausa(_ context.Context, fn func(tx *db.Tx) error) error { return fn(nil) }

// TestKlausaKasus - tiket 47: GET / PUT ClauseList (bentuk KlausaKasus persis, baca ulang sesudah simpan, argumentList
// selalu larik), pemeriksaan isian (kode wajib / ganda, bahasa 0-2, argumentCount angka, lebar 4000), 401 / 404 / 503.
func TestKlausaKasus(t *testing.T) {
	k := &klausaTiruan{simpan: map[string][]models.KlausaKasus{}}
	svc := services.Baru(nil).DenganKlausa(k).DenganTransaksi(tanpaOracleKlausa)
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/kasus/K1/klausa", "", ""); kode != 200 || isi != `{"baris":[]}` {
		t.Fatalf("kosong: %d %s", kode, isi)
	}
	satu := `{"clauseCode":"UJI-1","clauseTitle":"UJI JUDUL","clauseDescription":"UJI INFO","clauseLanguage":"0",` +
		`"clauseLanguageId":"UJI-BHS","clauseContent":"isi UJI NILAI","clauseContentTemp":"isi _&1","argumentCount":"1",` +
		`"argumentList":[{"argumentNumber":"1","argumentDescription":"UJI NAMA","argumentValue":"UJI NILAI"}]}`
	kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/K1/klausa", `{"baris":[`+satu+`]}`, "UJI-USER")
	// galat.TulisJSON meloloskan `&` sebagai backslash-u0026 (JSON sah; frontend membacanya kembali `_&1`).
	if kode != 200 || isi != `{"baris":[`+strings.ReplaceAll(satu, "&", "\\"+"u0026")+`]}` {
		t.Fatalf("simpan: %d %s", kode, isi)
	}
	tanpaArg := `{"clauseCode":"UJI-2","clauseTitle":"","clauseDescription":"","clauseLanguage":"","clauseLanguageId":"",` +
		`"clauseContent":"","clauseContentTemp":"","argumentCount":"","argumentList":null}`
	if kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/K1/klausa", `{"baris":[`+tanpaArg+`]}`, "UJI-USER"); kode != 200 ||
		!strings.Contains(isi, `"argumentList":[]`) {
		t.Errorf("tanpa argumen: %d %s", kode, isi)
	}
	ganti := func(a, b string) string { return `{"baris":[` + strings.Replace(satu, a, b, 1) + `]}` }
	for nama, c := range map[string]struct {
		badan, pelaku, mau string
		kode               int
	}{
		"kode kosong":     {ganti(`"UJI-1"`, `" "`), "UJI-USER", "clauseCode wajib diisi", 400},
		"kode ganda":      {`{"baris":[` + satu + `,` + satu + `]}`, "UJI-USER", "clauseCode UJI-1 ganda", 400},
		"bahasa":          {ganti(`"clauseLanguage":"0"`, `"clauseLanguage":"3"`), "UJI-USER", "clauseLanguage harus 0, 1, atau 2", 400},
		"jumlah argumen":  {ganti(`"argumentCount":"1"`, `"argumentCount":"satu"`), "UJI-USER", "argumentCount harus angka", 400},
		"isi panjang":     {ganti(`"isi UJI NILAI"`, `"`+strings.Repeat("x", 4001)+`"`), "UJI-USER", "clauseContent paling banyak 4000 bita", 400},
		"argumen panjang": {ganti(`"argumentValue":"UJI NILAI"`, `"argumentValue":"`+strings.Repeat("x", 4001)+`"`), "UJI-USER", "argumentList[0].argumentValue", 400},
		"kunci asing":     {`{"baris":[],"lain":1}`, "UJI-USER", "baris", 400},
		"tanpa baris":     {`{}`, "UJI-USER", "baris", 400},
		"tanpa identitas": {`{"baris":[]}`, "", "", 401},
		"case tidak ada":  {`{"baris":[]}`, "UJI-USER", "", 404},
	} {
		jalur := "/api/nbfacin/kasus/K1/klausa"
		if nama == "case tidak ada" {
			jalur = "/api/nbfacin/kasus/K404/klausa"
		}
		if kode, isi := minta(t, svc, "PUT", jalur, c.badan, c.pelaku); kode != c.kode || !strings.Contains(isi, c.mau) {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
	if kode, _ := minta(t, svc, "GET", "/api/nbfacin/kasus/K404/klausa", "", ""); kode != 404 {
		t.Errorf("GET case tidak ada: %d", kode)
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/kasus/K1/klausa", "", ""); kode != 503 {
		t.Errorf("tanpa basis data: %d", kode)
	}
}

// TestArgumenDanCariKlausa - GET /api/nbfacin/klausa/{id}/argumen (SearchClauseArgFireSQL) dan pencarian
// (RetrieveClauseSQL): bahasa kode bawaan "0", language = label ClauseLanguageID, 10 per halaman, 400 bahasa / halaman.
func TestArgumenDanCariKlausa(t *testing.T) {
	k := &klausaTiruan{simpan: map[string][]models.KlausaKasus{}}
	svc := services.Baru(nil).DenganKlausa(k)
	kode, isi := minta(t, svc, "GET", "/api/nbfacin/klausa/UJI-1/argumen", "", "")
	if kode != 200 || isi != `{"baris":[{"argumentNumber":"1","argumentDescription":"UJI NAMA","argumentValue":"UJI NILAI"}]}` || k.idArg != "UJI-1" {
		t.Errorf("argumen: %d %s %q", kode, isi, k.idArg)
	}
	// K47-6: M_ARGCLAUSEFIRE tidak ada -> 200 daftar kosong (popup Isi Klasula hanya isi klausa, seperti Pega).
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/klausa/TIADA/argumen", "", ""); kode != 200 || isi != `{"baris":[]}` {
		t.Errorf("argumen tiada: %d %s", kode, isi)
	}
	if kode, _ := minta(t, svc, "GET", "/api/nbfacin/klausa/%20/argumen", "", ""); kode != 400 {
		t.Errorf("id kosong: %d", kode)
	}
	kode, isi = minta(t, svc, "GET", "/api/nbfacin/klausa?q=UJI", "", "")
	if kode != 200 || k.cari != "0|UJI|1|10" || isi != `{"baris":[{"id":"K01","title":"UJI JUDUL","info":"UJI INFO","text":"UJI ISI",`+
		`"language":"Indonesia","argumentCount":"2"}],"total":11,"halaman":1,"ukuran":10}` {
		t.Errorf("cari: %d %s %q", kode, isi, k.cari)
	}
	for kodeBhs, label := range map[string]string{"1": "Inggris", "2": "Dual Bahasa"} {
		if _, isi := minta(t, svc, "GET", "/api/nbfacin/klausa?bahasa="+kodeBhs, "", ""); !strings.Contains(isi, `"language":"`+label+`"`) {
			t.Errorf("bahasa %s: %s", kodeBhs, isi)
		}
	}
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/klausa?halaman=2", "", ""); kode != 200 || isi != `{"baris":[],"total":11,"halaman":2,"ukuran":10}` {
		t.Errorf("halaman 2: %d %s", kode, isi)
	}
	for _, j := range []string{"/api/nbfacin/klausa?bahasa=3", "/api/nbfacin/klausa?halaman=0", "/api/nbfacin/klausa?halaman=x"} {
		if kode, _ := minta(t, svc, "GET", j, "", ""); kode != 400 {
			t.Errorf("%s: %d, mau 400", j, kode)
		}
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/klausa", "", ""); kode != 503 {
		t.Errorf("cari tanpa basis data: %d", kode)
	}
}
