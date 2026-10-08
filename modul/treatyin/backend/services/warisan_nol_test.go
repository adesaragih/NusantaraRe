package services

// Bukti bahwa `nolkanLarik` LENGKAP - bukan bahwa ia menyebut medan yang
// kebetulan sudah ada.
//
// ⛔ Uji ini sengaja memakai REFLEKSI alih-alih menyebut medan satu per satu.
// Uji yang menyebut medan hanya menguji ingatan penulisnya: ia lulus untuk
// medan yang sudah dijaga, dan diam untuk medan yang besok ditambahkan. Yang
// diperlukan justru kebalikannya - GAGAL begitu `KontrakWarisan` tumbuh satu
// larik yang `nolkanLarik` belum sebut.
//
// Galat yang melahirkannya: `Cannot read properties of null (reading
// 'length')`, dilaporkan pemakai 6 Oktober 2026 sesudah menekan `Edit` lalu
// `View`.

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"nusantarare/modul/treatyin/backend/models"
)

// nilYangTersisa menyusuri sebuah nilai dan mengembalikan JALUR setiap iris
// atau peta yang masih `nil`.
//
// Menyusuri struct DI DALAM struct juga - `OpsiKepala` dan kedua
// `TabTeksWarisan` membawa lariknya sendiri, dan larik bersarang merender
// `null` persis sama dengan larik akar.
func nilYangTersisa(v reflect.Value, jalur string) []string {
	var sisa []string
	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		for i := range t.NumField() {
			f := t.Field(i)
			if !f.IsExported() {
				continue
			}
			sisa = append(sisa, nilYangTersisa(v.Field(i), jalur+"."+f.Name)...)
		}
	case reflect.Slice, reflect.Map:
		if v.IsNil() {
			return []string{jalur}
		}
		for i := range v.Len() {
			if v.Kind() == reflect.Slice {
				sisa = append(sisa, nilYangTersisa(v.Index(i), jalur+"[]")...)
			}
		}
	default:
	}
	return sisa
}

// Kontrak NOL-NILAI - keadaan yang paling buruk, dan bukan keadaan yang
// mengada-ada: kontrak yang `JSONDATA`-nya tidak ada pulang lewat jalur
// `return k, nil` lebih awal, dengan seluruh lariknya belum pernah disentuh.
func TestNolkanLarikTidakMenyisakanNil(t *testing.T) {
	var k models.KontrakWarisan
	nolkanLarik(&k)
	if sisa := nilYangTersisa(reflect.ValueOf(k), "KontrakWarisan"); len(sisa) > 0 {
		t.Fatalf("masih ada larik/peta nil sesudah nolkanLarik: %v\n"+
			"tambahkan medan itu ke `nolkanLarik` di services/warisan_nol.go", sisa)
	}
}

// Larik DI DALAM elemen larik - `BarisLayerWarisan.KelasBisnis` hari ini,
// dan apa pun yang kelak menyusul. Tiap larik-berisi-struct diisi SATU
// elemen nol-nilai, lalu elemennya diperiksa.
func TestNolkanLarikMenjangkauLarikBersarang(t *testing.T) {
	var k models.KontrakWarisan
	v := reflect.ValueOf(&k).Elem()
	tp := v.Type()
	diisi := 0
	for i := range tp.NumField() {
		f := v.Field(i)
		if f.Kind() != reflect.Slice || tp.Field(i).Type.Elem().Kind() != reflect.Struct {
			continue
		}
		f.Set(reflect.MakeSlice(f.Type(), 1, 1))
		diisi++
	}
	if diisi == 0 {
		t.Fatal("nol larik-berisi-struct ditemukan - uji ini tidak menguji apa pun")
	}
	nolkanLarik(&k)
	if sisa := nilYangTersisa(reflect.ValueOf(k), "KontrakWarisan"); len(sisa) > 0 {
		t.Fatalf("larik bersarang masih nil sesudah nolkanLarik: %v", sisa)
	}
}

// Dan inilah bentuk yang sungguh sampai ke layar. Uji di atas memeriksa
// struct Go; yang menghentikan halaman adalah JSON-nya.
func TestKontrakWarisanNolTidakMengirimNullKeLayar(t *testing.T) {
	var k models.KontrakWarisan
	nolkanLarik(&k)
	b, err := json.Marshal(k)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(b), ":null") {
		t.Fatalf("jawaban masih memuat `null`; layar membaca `.length` di atasnya:\n%s", b)
	}
}
