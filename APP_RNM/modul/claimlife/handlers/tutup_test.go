package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"nusantarare/modul/claimlife/services"
)

// Uji pintu gerbang `Close Claim`.

func TestRuteBolehTutupTerdaftar(t *testing.T) {
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(isi), `"GET /api/klaim-life/{id}/boleh-tutup"`) {
		t.Error("rute gerbang tutup tidak terdaftar di Router")
	}
	// ⛔ METODENYA GET, dan itu tetap dikunci walau penutupan kini ADA.
	// Dua rute, dua janji: `boleh-tutup` MENJAWAB, `tutup` MENGUBAH. Rute
	// yang menjawab pertanyaan dengan perbuatan adalah rute yang layar
	// panggil untuk memutuskan apakah tombolnya pantas ditawarkan - lalu
	// menutup kasusnya sebagai efek samping.
	if strings.Contains(string(isi), `"POST /api/klaim-life/{id}/boleh-tutup"`) {
		t.Error("gerbang pemeriksa terdaftar sebagai POST; ia tidak mengubah apa pun")
	}
}

func TestRuteTutupTerdaftarSebagaiPOST(t *testing.T) {
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	// ⛔ Butir bb. POST, sebab ia MENGUBAH: STATUS_WORK diisi
	// Resolved-Completed dan TAHAP dikosongkan.
	if !strings.Contains(string(isi), `"POST /api/klaim-life/{id}/tutup"`) {
		t.Error("rute penutupan kasus tidak terdaftar di Router")
	}
	// Dan ia BUKAN GET: penutupan yang dapat dipicu pranala, prefetch, atau
	// tombol kembali adalah penutupan yang terjadi tanpa ada yang menekannya.
	if strings.Contains(string(isi), `"GET /api/klaim-life/{id}/tutup"`) {
		t.Error("penutupan kasus terdaftar sebagai GET")
	}
}

func TestTutupTanpaDatabaseMenjawab503(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/klaim-life/K-1/tutup", nil)
	r.SetPathValue("id", "K-1")
	tutupKlaim(services.New(nil), true)(w, r)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("kode = %d, mau 503", w.Code)
	}
	var isi map[string]any
	if err := json.NewDecoder(w.Body).Decode(&isi); err != nil {
		t.Fatalf("badan bukan JSON: %v", err)
	}
	if _, ada := isi["galat"]; !ada {
		t.Errorf("envelope tanpa kunci \"galat\": %v", isi)
	}
}

// TestBadanPenghalangMemakaiAmplopYangSama mengunci bentuk 409.
//
// ⛔ Kunci galatnya "galat", sama seperti amplop galat lain. Cacat
// `galat` vs `error` yang pernah memakan seluruh pesan backend lahir persis
// dari dua bentuk amplop yang masing-masing benar menurut dirinya sendiri.
func TestBadanPenghalangMemakaiAmplopYangSama(t *testing.T) {
	b, err := json.Marshal(jawabanPenghalang{
		Galat: "klaim belum boleh ditutup: ada peserta yang belum diaksep",
		Penghalang: []services.PenghalangTampil{
			{Urutan: 2, NomorSertifikat: "UJI-002",
				Pesan: "UJI-002 is not approved yet, on list 2"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	var isi map[string]json.RawMessage
	if err := json.Unmarshal(b, &isi); err != nil {
		t.Fatal(err)
	}
	kunci := make([]string, 0, len(isi))
	for k := range isi {
		kunci = append(kunci, k)
	}
	sort.Strings(kunci)
	if !reflect.DeepEqual(kunci, []string{"galat", "penghalang"}) {
		t.Errorf("kunci badan 409 = %v, mau [galat penghalang]", kunci)
	}
	// ⛔ SELURUH penghalang menyeberang, beserta kalimatnya. Pega
	// memasang pesannya di dalam loop, sekali per peserta yang tertandai.
	if !strings.Contains(string(b), "is not approved yet, on list 2") {
		t.Errorf("kalimat rule tidak menyeberang: %s", b)
	}
}

func TestBolehTutupTanpaDatabaseMenjawab503(t *testing.T) {
	// ⛔ 503, BUKAN 200 dengan boleh=false, dan bukan 200 dengan boleh=true.
	// Keduanya berbohong: yang pertama menuduh pesertanya belum diaksep,
	// yang kedua mengizinkan penutupan tanpa memeriksa apa pun.
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/klaim-life/K-1/boleh-tutup", nil)
	r.SetPathValue("id", "K-1")
	bolehTutup(services.New(nil))(w, r)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("kode = %d, mau 503", w.Code)
	}
	var isi map[string]any
	if err := json.NewDecoder(w.Body).Decode(&isi); err != nil {
		t.Fatalf("badan bukan JSON: %v", err)
	}
	if _, ada := isi["galat"]; !ada {
		t.Errorf("envelope tanpa kunci \"galat\": %v", isi)
	}
	if _, ada := isi["boleh"]; ada {
		t.Error("jawaban galat memuat \"boleh\"; ia belum memeriksa apa pun")
	}
}
