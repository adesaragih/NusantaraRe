package login

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"nusantarare/inti/backend/menu"
)

// serverKelola - rute login + Kelola User di belakang middleware, dan cookie
// sesi UJI-ADMIN (pemegang Kelola User).
func serverKelola(t *testing.T, g *gudangTiruan) (http.Handler, *http.Cookie) {
	t.Helper()
	r := NewRute(layananUji(g, saatUji), true).DenganKelola(kelolaUji(g))
	mux := http.NewServeMux()
	r.Pasang(mux)
	h := r.Middleware(mux)
	c := cookieSesi(t, kirim(h, "POST", "/api/auth/login", `{"akun":"UJI-ADMIN","sandi":"Sandi-Benar-01"}`, nil))
	if c == nil {
		t.Fatal("login admin tanpa cookie")
	}
	return h, c
}

// ⛔ Tanpa sesi 401, sesi tanpa Kelola User 403 - dan pencabutan menu berlaku
// pada permintaan BERIKUTNYA, tanpa login ulang.
func TestKelolaRuteMenuntutAdmin(t *testing.T) {
	g := gudangUji(t)
	h, c := serverKelola(t, g)
	if rec := kirim(h, "GET", "/api/admin/pengguna", "", nil); rec.Code != http.StatusUnauthorized {
		t.Errorf("tanpa sesi: %d", rec.Code)
	}
	if rec := kirim(h, "GET", "/api/admin/pengguna", "", c); rec.Code != http.StatusOK {
		t.Fatalf("admin: %d %s", rec.Code, rec.Body)
	}
	g.menu["UJI-ADMIN"] = []string{"claimlife"}
	for _, k := range []struct{ metode, jalur, badan string }{
		{"GET", "/api/admin/pengguna", ""},
		{"GET", "/api/admin/pengguna/UJI-KUNCI", ""},
		{"GET", "/api/admin/pilihan-pengguna", ""},
		{"POST", "/api/admin/pengguna", `{"akunId":"UJI-X","nama":"X","sandi":"Sandi-Admin-01"}`},
		{"PUT", "/api/admin/pengguna/UJI-KUNCI", `{"nama":"X"}`},
		{"POST", "/api/admin/pengguna/UJI-KUNCI/aktif", `{"aktif":false}`},
		{"POST", "/api/admin/pengguna/UJI-KUNCI/buka-kunci", ""},
		{"DELETE", "/api/admin/pengguna/UJI-KUNCI", ""},
	} {
		if rec := kirim(h, k.metode, k.jalur, k.badan, c); rec.Code != http.StatusForbidden {
			t.Errorf("%s %s tanpa Kelola User: %d, mau 403", k.metode, k.jalur, rec.Code)
		}
	}
	if _, ada := g.akun["UJI-X"]; ada || !g.akun["UJI-KUNCI"].Aktif || g.akun["UJI-KUNCI"].Nama != "Uji Kunci" {
		t.Error("perubahan terjadi walau ditolak 403")
	}
}

func TestKelolaRuteTanpaLayanan503(t *testing.T) {
	r := NewRute(layananUji(gudangUji(t), saatUji), true)
	mux := http.NewServeMux()
	r.Pasang(mux)
	if rec := kirim(r.Middleware(mux), "GET", "/api/admin/pengguna", "", nil); rec.Code != http.StatusServiceUnavailable {
		t.Errorf("tanpa Kelola: %d", rec.Code)
	}
}

func TestKelolaRuteAlurAdmin(t *testing.T) {
	g := gudangUji(t)
	h, c := serverKelola(t, g)
	badan := `{"akunId":"UJI-KELOLA-3","nama":"Uji Tiga","organisasi":"RNM","divisi":"TECH","unit":"CLM",` +
		`"workbasket":["ReasLifeSPV"],"menu":["claimlife"],"sandi":"Sandi-Admin-01"}`
	rec := kirim(h, "POST", "/api/admin/pengguna", badan, c)
	if rec.Code != http.StatusCreated {
		t.Fatalf("buat: %d %s", rec.Code, rec.Body)
	}
	// ⛔ Sandi dan hash tidak pernah kembali ke layar.
	for _, bocor := range []string{"Sandi-Admin-01", "$2a$", "hash", "sandi"} {
		if strings.Contains(rec.Body.String(), bocor) {
			t.Errorf("jawaban buat memuat %q: %s", bocor, rec.Body)
		}
	}
	var r RinciAkun
	if err := json.Unmarshal(rec.Body.Bytes(), &r); err != nil || r.AkunID != "UJI-KELOLA-3" || !r.WajibGantiSandi ||
		len(r.Menu) != 1 || r.Menu[0] != "claimlife" {
		t.Errorf("rinci akun baru %+v %v", r, err)
	}
	for _, k := range []struct {
		metode, jalur, badan string
		kode                 int
	}{
		{"POST", "/api/admin/pengguna", badan, http.StatusConflict},
		{"POST", "/api/admin/pengguna", strings.Replace(badan, "Sandi-Admin-01", "pendek", 1), http.StatusBadRequest},
		{"POST", "/api/admin/pengguna", strings.Replace(badan, `"claimlife"`, `"tidakada"`, 1), http.StatusBadRequest},
		{"PUT", "/api/admin/pengguna/UJI-KELOLA-3", `{"nama":"Uji Tiga Baru","menu":["premiumlistlife"]}`, http.StatusOK},
		{"PUT", "/api/admin/pengguna/UJI-TIDAK-ADA", `{"nama":"X"}`, http.StatusNotFound},
		{"PUT", "/api/admin/pengguna/UJI-ADMIN", `{"nama":"Uji Admin","menu":["claimlife"]}`, http.StatusConflict},
		{"POST", "/api/admin/pengguna/UJI-ADMIN/aktif", `{"aktif":false}`, http.StatusConflict},
		{"POST", "/api/admin/pengguna/UJI-KELOLA-3/aktif", `{}`, http.StatusBadRequest},
		{"POST", "/api/admin/pengguna/UJI-KELOLA-3/aktif", `{"aktif":false}`, http.StatusOK},
		{"POST", "/api/admin/pengguna/UJI-KUNCI/buka-kunci", "", http.StatusOK},
		{"GET", "/api/admin/pengguna/UJI-KELOLA-3", "", http.StatusOK},
		{"GET", "/api/admin/pilihan-pengguna", "", http.StatusOK},
		{"DELETE", "/api/admin/pengguna/UJI-ADMIN", "", http.StatusConflict},
		{"DELETE", "/api/admin/pengguna/UJI-KELOLA-3", "", http.StatusNoContent},
		{"DELETE", "/api/admin/pengguna/UJI-KELOLA-3", "", http.StatusNotFound},
	} {
		if rec := kirim(h, k.metode, k.jalur, k.badan, c); rec.Code != k.kode {
			t.Errorf("%s %s %s: %d %s, mau %d", k.metode, k.jalur, k.badan, rec.Code, rec.Body, k.kode)
		}
	}
	if g.akun["UJI-KUNCI"].Terkunci {
		t.Error("buka kunci tidak mencabut kunci")
	}
	if _, ada := g.akun["UJI-KELOLA-3"]; ada {
		t.Error("hapus permanen menyisakan akun")
	}
	if !g.akun["UJI-ADMIN"].Aktif || !menuDipegang(g.menu["UJI-ADMIN"], menu.KodeKelolaUser) {
		t.Error("admin mengubah dirinya sendiri walau ditolak")
	}
}

// Tab Security lewat HTTP: password akun lain mencabut sesinya; password
// SENDIRI tidak memutus sesi admin (cookie diterbitkan ulang); "Change
// Password Next Login" saat membuat akun; nol password di jawaban.
func TestKelolaRuteSecurity(t *testing.T) {
	g := gudangUji(t)
	h, c := serverKelola(t, g)
	// Sesi UJI-KUNCI dibuka dulu (kuncinya dibuka admin), lalu passwordnya diganti admin.
	if rec := kirim(h, "POST", "/api/admin/pengguna/UJI-KUNCI/buka-kunci", "", c); rec.Code != http.StatusOK {
		t.Fatalf("buka kunci: %d", rec.Code)
	}
	cLain := cookieSesi(t, kirim(h, "POST", "/api/auth/login", `{"akun":"UJI-KUNCI","sandi":"Sandi-Benar-01"}`, nil))
	if cLain == nil {
		t.Fatal("login UJI-KUNCI")
	}
	rec := kirim(h, "POST", "/api/admin/pengguna/UJI-KUNCI/sandi", `{"sandi":"Sandi-Baru-Admin-1","wajibGanti":true}`, c)
	if rec.Code != http.StatusOK || strings.Contains(rec.Body.String(), "Sandi-Baru-Admin-1") || !strings.Contains(rec.Body.String(), `"wajibGantiSandi":true`) {
		t.Fatalf("atur sandi akun lain: %d %s", rec.Code, rec.Body)
	}
	if rec := kirim(h, "GET", "/api/auth/saya", "", cLain); rec.Code != http.StatusUnauthorized {
		t.Errorf("sesi akun yang passwordnya diganti: %d, mau 401", rec.Code)
	}
	// Password sendiri: sesi admin tetap hidup dengan cookie baru.
	rec = kirim(h, "POST", "/api/admin/pengguna/UJI-ADMIN/sandi", `{"sandi":"Sandi-Baru-Admin-2","wajibGanti":false}`, c)
	cBaru := cookieSesi(t, rec)
	if rec.Code != http.StatusOK || cBaru == nil {
		t.Fatalf("atur sandi sendiri: %d %s, cookie %v", rec.Code, rec.Body, cBaru)
	}
	if rec := kirim(h, "GET", "/api/auth/saya", "", cBaru); rec.Code != http.StatusOK {
		t.Errorf("sesi admin sesudah mengganti password sendiri: %d", rec.Code)
	}
	if rec := kirim(h, "GET", "/api/auth/saya", "", c); rec.Code != http.StatusUnauthorized {
		t.Errorf("cookie LAMA admin sesudah ganti password: %d, mau 401", rec.Code)
	}
	for _, k := range []struct {
		jalur, badan string
		kode         int
	}{
		{"/api/admin/pengguna/UJI-KUNCI/sandi", `{"sandi":"pendek","wajibGanti":true}`, http.StatusBadRequest},
		{"/api/admin/pengguna/UJI-KUNCI/sandi", `{"sandi":"Sandi-Baru-Admin-3"}`, http.StatusBadRequest},
		{"/api/admin/pengguna/UJI-TIDAK-ADA/sandi", `{"wajibGanti":false}`, http.StatusNotFound},
		{"/api/admin/pengguna/UJI-KUNCI/sandi", `{"wajibGanti":false}`, http.StatusOK},
	} {
		if rec := kirim(h, "POST", k.jalur, k.badan, cBaru); rec.Code != k.kode {
			t.Errorf("%s %s: %d %s, mau %d", k.jalur, k.badan, rec.Code, rec.Body, k.kode)
		}
	}
	// Membuat akun tanpa centang wajib ganti; tanpa medan wajibGanti = wajib (aman).
	for _, k := range []struct {
		id, ekor string
		wajib    bool
	}{{"UJI-SEC-1", `,"wajibGanti":false`, false}, {"UJI-SEC-2", ``, true}} {
		badan := `{"akunId":"` + k.id + `","nama":"Sec","sandi":"Sandi-Admin-01"` + k.ekor + `}`
		if rec := kirim(h, "POST", "/api/admin/pengguna", badan, cBaru); rec.Code != http.StatusCreated || g.akun[k.id].WajibGantiSandi != k.wajib {
			t.Errorf("buat %s: %d %s, wajib %v", k.id, rec.Code, rec.Body, g.akun[k.id] != nil && g.akun[k.id].WajibGantiSandi)
		}
	}
}
