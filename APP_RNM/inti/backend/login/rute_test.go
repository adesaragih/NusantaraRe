package login

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
)

// serverUji merakit rute login dan satu rute beridentitas di belakang
// middleware - seperti cmd/api.
func serverUji(t *testing.T, g *gudangTiruan, saat time.Time) http.Handler {
	t.Helper()
	r := NewRute(layananUji(g, saat), true)
	mux := http.NewServeMux()
	r.Pasang(mux)
	mux.HandleFunc("GET /api/uji/pelaku", func(w http.ResponseWriter, req *http.Request) {
		_ = json.NewEncoder(w).Encode(inti.PelakuDari(req, false))
	})
	return r.Middleware(mux)
}

func kirim(h http.Handler, metode, jalur, badan string, cookie *http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest(metode, jalur, strings.NewReader(badan))
	if badan != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if cookie != nil {
		req.AddCookie(cookie)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func cookieSesi(t *testing.T, rec *httptest.ResponseRecorder) *http.Cookie {
	t.Helper()
	for _, c := range rec.Result().Cookies() {
		if c.Name == NamaCookie {
			return c
		}
	}
	return nil
}

func TestLoginMemasangCookieAman(t *testing.T) {
	h := serverUji(t, gudangUji(t), saatUji)
	rec := kirim(h, "POST", "/api/auth/login", `{"akun":"UJI-ADMIN","sandi":"Sandi-Benar-01"}`, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("login: %d %s", rec.Code, rec.Body)
	}
	c := cookieSesi(t, rec)
	if c == nil || !c.HttpOnly || !c.Secure || c.SameSite != http.SameSiteStrictMode || c.Path != "/" || c.Value == "" {
		t.Fatalf("cookie sesi: %+v", c)
	}
	badan := rec.Body.String()
	for _, bocor := range []string{"Sandi-Benar-01", "$2a$", "hash", "Versi"} {
		if strings.Contains(badan, bocor) {
			t.Errorf("jawaban login memuat %q: %s", bocor, badan)
		}
	}
	var p Profil
	if err := json.Unmarshal(rec.Body.Bytes(), &p); err != nil || p.AkunID != "UJI-ADMIN" || len(p.Peran) != 2 {
		t.Errorf("profil %+v %v", p, err)
	}
}

func TestLoginGagalTanpaCookie(t *testing.T) {
	for _, k := range []struct {
		badan string
		kode  int
	}{
		{`{"akun":"UJI-ADMIN","sandi":"Sandi-Salah-01"}`, http.StatusUnauthorized},
		{`{"akun":"UJI-TIDAK-ADA","sandi":"Sandi-Benar-01"}`, http.StatusUnauthorized},
		{`{"akun":"UJI-KUNCI","sandi":"Sandi-Benar-01"}`, http.StatusLocked},
		{`bukan json`, http.StatusBadRequest},
	} {
		rec := kirim(serverUji(t, gudangUji(t), saatUji), "POST", "/api/auth/login", k.badan, nil)
		if rec.Code != k.kode || cookieSesi(t, rec) != nil {
			t.Errorf("%s: %d (mau %d), cookie %v", k.badan, rec.Code, k.kode, cookieSesi(t, rec))
		}
	}
	// Tanpa Content-Type JSON: formulir lintas situs tidak dapat mengirimnya.
	req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(`{"akun":"UJI-ADMIN","sandi":"Sandi-Benar-01"}`))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()
	serverUji(t, gudangUji(t), saatUji).ServeHTTP(rec, req)
	if rec.Code != http.StatusUnsupportedMediaType {
		t.Errorf("text/plain: %d", rec.Code)
	}
}

func TestLoginTanpaLayananMenjawab503(t *testing.T) {
	mux := http.NewServeMux()
	r := NewRute(nil, true)
	r.Pasang(mux)
	rec := kirim(r.Middleware(mux), "POST", "/api/auth/login", `{"akun":"UJI-ADMIN","sandi":"x"}`, nil)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("tanpa layanan: %d", rec.Code)
	}
}

func TestSesiMengisiPelaku(t *testing.T) {
	g := gudangUji(t)
	h := serverUji(t, g, saatUji)
	c := cookieSesi(t, kirim(h, "POST", "/api/auth/login", `{"akun":"UJI-ADMIN","sandi":"Sandi-Benar-01"}`, nil))
	if rec := kirim(h, "GET", "/api/auth/saya", "", nil); rec.Code != http.StatusUnauthorized {
		t.Errorf("saya tanpa cookie: %d", rec.Code)
	}
	rec := kirim(h, "GET", "/api/auth/saya", "", c)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"akunId":"UJI-ADMIN"`) {
		t.Errorf("saya: %d %s", rec.Code, rec.Body)
	}
	if cookieSesi(t, rec) == nil {
		t.Error("sesi tidak diperpanjang (cookie tidak diterbitkan ulang)")
	}
	var p inti.Pelaku
	_ = json.Unmarshal(kirim(h, "GET", "/api/uji/pelaku", "", c).Body.Bytes(), &p)
	if p.AkunID != "UJI-ADMIN" || len(p.Peran) != 2 || p.Peran[0] != "ReasLifeAdmin" {
		t.Errorf("inti.PelakuDari dari sesi: %+v", p)
	}
}

// Akun yang wajib ganti sandi belum mendapat pelaku: seluruh jalur modul
// menolaknya, hanya /api/auth/* yang melayani.
func TestWajibGantiSandiBelumBerpelaku(t *testing.T) {
	h := serverUji(t, gudangUji(t), saatUji)
	c := cookieSesi(t, kirim(h, "POST", "/api/auth/login", `{"akun":"UJI-BARU","sandi":"Sandi-Benar-01"}`, nil))
	var p inti.Pelaku
	_ = json.Unmarshal(kirim(h, "GET", "/api/uji/pelaku", "", c).Body.Bytes(), &p)
	if p.AkunID != "" {
		t.Errorf("akun wajib ganti sandi sudah berpelaku: %+v", p)
	}
	if rec := kirim(h, "GET", "/api/auth/saya", "", c); !strings.Contains(rec.Body.String(), `"wajibGantiSandi":true`) {
		t.Errorf("saya: %s", rec.Body)
	}
	rec := kirim(h, "POST", "/api/auth/ganti-sandi", `{"sandiLama":"Sandi-Benar-01","sandiBaru":"Sandi-Baru-0001"}`, c)
	baru := cookieSesi(t, rec)
	if rec.Code != http.StatusOK || baru == nil {
		t.Fatalf("ganti sandi: %d %s", rec.Code, rec.Body)
	}
	_ = json.Unmarshal(kirim(h, "GET", "/api/uji/pelaku", "", baru).Body.Bytes(), &p)
	if p.AkunID != "UJI-BARU" {
		t.Errorf("sesudah ganti sandi belum berpelaku: %+v", p)
	}
	if rec := kirim(h, "GET", "/api/auth/saya", "", c); rec.Code != http.StatusUnauthorized {
		t.Errorf("cookie lama sesudah ganti sandi: %d", rec.Code)
	}
}

func TestGantiSandiSalahBukan401(t *testing.T) {
	h := serverUji(t, gudangUji(t), saatUji)
	c := cookieSesi(t, kirim(h, "POST", "/api/auth/login", `{"akun":"UJI-BARU","sandi":"Sandi-Benar-01"}`, nil))
	for _, badan := range []string{
		`{"sandiLama":"Sandi-Salah-01","sandiBaru":"Sandi-Baru-0001"}`,
		`{"sandiLama":"Sandi-Benar-01","sandiBaru":"pendek"}`,
	} {
		// 400, bukan 401: 401 dibaca layar sebagai "sesi berakhir".
		if rec := kirim(h, "POST", "/api/auth/ganti-sandi", badan, c); rec.Code != http.StatusBadRequest {
			t.Errorf("%s: %d", badan, rec.Code)
		}
	}
}

func TestLogoutMencabutSesi(t *testing.T) {
	h := serverUji(t, gudangUji(t), saatUji)
	c := cookieSesi(t, kirim(h, "POST", "/api/auth/login", `{"akun":"UJI-ADMIN","sandi":"Sandi-Benar-01"}`, nil))
	rec := kirim(h, "POST", "/api/auth/logout", "", c)
	hapus := cookieSesi(t, rec)
	if rec.Code != http.StatusNoContent || hapus == nil || hapus.MaxAge >= 0 {
		t.Fatalf("logout: %d %+v", rec.Code, hapus)
	}
	if rec := kirim(h, "GET", "/api/auth/saya", "", c); rec.Code != http.StatusUnauthorized {
		t.Errorf("cookie sesudah logout: %d", rec.Code)
	}
}

func TestCookiePalsuDibuang(t *testing.T) {
	h := serverUji(t, gudangUji(t), saatUji)
	rec := kirim(h, "GET", "/api/auth/saya", "", &http.Cookie{Name: NamaCookie, Value: "palsu.palsu"})
	if c := cookieSesi(t, rec); rec.Code != http.StatusUnauthorized || c == nil || c.MaxAge >= 0 {
		t.Errorf("cookie palsu: %d %+v", rec.Code, c)
	}
}
