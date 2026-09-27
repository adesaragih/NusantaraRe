package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"nusantarare/internal/services"
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
	// ⛔ METODENYA GET, dan itu dikunci. Activity Pega memeriksa LALU
	// menyelesaikan penugasan; yang dibangun baru pemeriksaannya. POST akan
	// menjanjikan penutupan yang belum terjadi, dan tombol yang menjanjikan
	// lebih daripada yang ia lakukan adalah cacat yang paling mahal
	// ditemukan belakangan.
	if strings.Contains(string(isi), `"POST /api/klaim-life/{id}/boleh-tutup"`) {
		t.Error("gerbang tutup terdaftar sebagai POST; ia belum menutup apa pun")
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
