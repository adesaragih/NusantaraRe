package main

// Gerbang menu (Kelola User, keputusan work owner 01-10-2026) - TANPA Oracle.
//
// Yang dijaga: rute modul aktif dilayani HANYA bagi akun yang memegang menu
// modul itu, atau menu modul yang meminjam rutenya; tanpa sesi 401 kecuali
// AUTH_STUB; dan setiap panggilan layar ke rute modul LAIN terdaftar di
// `ruteDipinjam` - panggilan yang lupa didaftarkan akan 403 bagi pemegang
// menu layar itu, diam-diam, di layar yang sedang ia pakai.

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/daftar"
)

func muxGerbang(t *testing.T, stub bool) http.Handler {
	t.Helper()
	terdaftar := modulTerdaftar(t)
	aktif, err := pilihModulAktif(terdaftar, daftar.NamaLama(), nil)
	if err != nil {
		t.Fatal(err)
	}
	return rakitMux(inti.NewDasar(nil), terdaftar, aktif, stub, rakitLogin(inti.NewDasar(nil), config.Config{}), daftar.HakLihat())
}

// kodeDengan - kode jawaban satu permintaan GET; `menu` nil = tanpa sesi.
func kodeDengan(h http.Handler, jalur string, menu []string) int {
	r := httptest.NewRequest(http.MethodGet, jalur, nil)
	if menu != nil {
		r = r.WithContext(inti.DenganAksesMenu(context.Background(), menu))
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w.Code
}

func TestGerbangMenuModul(t *testing.T) {
	ditolak := func(k int) bool { return k == http.StatusUnauthorized || k == http.StatusForbidden }
	for _, stub := range []bool{false, true} {
		h := muxGerbang(t, stub)
		// Sesi pemegang Claim Life saja.
		claim := []string{"claimlife"}
		if k := kodeDengan(h, "/api/klaim-life", claim); ditolak(k) {
			t.Errorf("stub=%v: rute menunya sendiri ditolak %d", stub, k)
		}
		for _, j := range []string{"/api/polis-life", "/api/komite", "/api/treaty-contract-out/tahun"} {
			if k := kodeDengan(h, j, claim); k != http.StatusForbidden {
				t.Errorf("stub=%v: %s tanpa menunya %d, mau 403", stub, j, k)
			}
		}
		// Dipinjam Claim Life (panel Data Polis) - BUKAN rute PremiumList lain.
		if k := kodeDengan(h, "/api/polis-life/ringkas?nomorPolis=UJI-1", claim); ditolak(k) {
			t.Errorf("stub=%v: rute pinjaman ditolak %d", stub, k)
		}
		if k := kodeDengan(h, "/api/polis-life/ringkas?nomorPolis=UJI-1", []string{"komiteclaimlife"}); k != http.StatusForbidden {
			t.Errorf("stub=%v: rute pinjaman bagi bukan peminjam %d, mau 403", stub, k)
		}
		// Sesi tanpa satu menu pun: semua rute modul 403.
		if k := kodeDengan(h, "/api/klaim-life", []string{}); k != http.StatusForbidden {
			t.Errorf("stub=%v: sesi tanpa menu %d, mau 403", stub, k)
		}
		// Rute milik aplikasi tidak digerbang menu.
		if k := kodeDengan(h, "/healthz", []string{}); k != http.StatusOK {
			t.Errorf("stub=%v: /healthz %d", stub, k)
		}
		if k := kodeDengan(h, "/api/modul-aktif", []string{}); k != http.StatusOK {
			t.Errorf("stub=%v: /api/modul-aktif %d", stub, k)
		}
	}
	// Tanpa sesi: 401 kecuali AUTH_STUB.
	if k := kodeDengan(muxGerbang(t, false), "/api/klaim-life", nil); k != http.StatusUnauthorized {
		t.Errorf("tanpa sesi, tanpa stub: %d, mau 401", k)
	}
	if k := kodeDengan(muxGerbang(t, true), "/api/klaim-life", nil); k == http.StatusUnauthorized || k == http.StatusForbidden {
		t.Errorf("tanpa sesi, stub: %d, mau dilayani", k)
	}
}

// Delapan modul master (04-10-2026): satu mesin bersama `inti/backend/master`, tetapi gerbangnya PER MODUL - pemegang
// menu City memakai rute City (termasuk saran rujukan Province miliknya sendiri), dan rute Province 403. Rute lama
// `/api/masterdata` sudah tidak ada (404), bukan digerbang.
func TestGerbangMenuMaster(t *testing.T) {
	h := muxGerbang(t, false)
	kota := []string{"mastercity"}
	for _, j := range []string{"/api/master-city/meta", "/api/master-city/rujukan/provinceId"} {
		if k := kodeDengan(h, j, kota); k == http.StatusUnauthorized || k == http.StatusForbidden || k == http.StatusNotFound {
			t.Errorf("%s bagi pemegang City: %d", j, k)
		}
	}
	for _, j := range []string{"/api/master-province", "/api/master-province/meta", "/api/master-nation/meta"} {
		if k := kodeDengan(h, j, kota); k != http.StatusForbidden {
			t.Errorf("%s bagi pemegang City: %d, mau 403", j, k)
		}
	}
	if k := kodeDengan(h, "/api/master-province/meta", []string{"masterprovince"}); k != http.StatusOK {
		t.Errorf("meta Province bagi pemegang Province: %d, mau 200", k)
	}
	if k := kodeDengan(h, "/api/masterdata", []string{"masterprovince"}); k != http.StatusNotFound {
		t.Errorf("/api/masterdata: %d, mau 404", k)
	}
}

// metodeUji - metode yang dicoba saat mengenali pemilik sebuah jalur.
var metodeUji = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}

// polaPemilik - pola tiap modul terdaftar yang cocok dengan jalur itu, per modul.
func polaPemilik(t *testing.T, jalur string) map[string][]string {
	t.Helper()
	out := map[string][]string{}
	for _, m := range modulTerdaftar(t) {
		k := kenali(m)
		for _, metode := range metodeUji {
			if _, p := k.mux.Handler(httptest.NewRequest(metode, jalur, nil)); p != "" {
				out[k.nama] = append(out[k.nama], p)
			}
		}
	}
	return out
}

// Setiap pola `ruteDipinjam` adalah pola SATU modul terdaftar, dan setiap
// peminjamnya modul terdaftar LAIN.
func TestRuteDipinjamTerdaftar(t *testing.T) {
	nama := map[string]bool{}
	for _, m := range modulTerdaftar(t) {
		nama[m.Nama()] = true
	}
	for pola, peminjam := range ruteDipinjam {
		bagian := strings.SplitN(pola, " ", 2)
		jalur := regexp.MustCompile(`\{[^}]+\}`).ReplaceAllString(bagian[1], "x")
		var pemilik []string
		for m, ps := range polaPemilik(t, jalur) {
			for _, p := range ps {
				if p == pola {
					pemilik = append(pemilik, m)
				}
			}
		}
		if len(pemilik) != 1 {
			t.Errorf("%s: dimiliki %v, mau tepat satu modul", pola, pemilik)
			continue
		}
		for _, p := range peminjam {
			if !nama[p] || p == pemilik[0] {
				t.Errorf("%s: peminjam %q bukan modul terdaftar lain (pemilik %s)", pola, p, pemilik[0])
			}
		}
	}
}

// polaJalurAPI - literal teks atau templat yang dimulai `/api/` di kode layar.
var polaJalurAPI = regexp.MustCompile("['\"`](/api/[^'\"`\\s]*)['\"`]")

// ⛔ Setiap jalur `/api/...` yang ditulis layar modul X dan dimiliki modul
// LAIN Y harus terdaftar di `ruteDipinjam` dengan peminjam X. Tanpanya
// pemegang menu X yang tidak memegang Y mendapat 403 di layar X.
func TestPanggilanLintasModulTerdaftar(t *testing.T) {
	akar := filepath.Join("..", "..", "modul")
	diperiksa := 0
	var temuan []string
	err := filepath.WalkDir(akar, func(jalur string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == "node_modules" {
			return filepath.SkipDir
		}
		rel, _ := filepath.Rel(akar, jalur)
		bagian := strings.Split(filepath.ToSlash(rel), "/")
		if d.IsDir() || len(bagian) < 3 || bagian[1] != "frontend" ||
			!(strings.HasSuffix(jalur, ".ts") || strings.HasSuffix(jalur, ".tsx")) || strings.Contains(d.Name(), ".test.") {
			return nil
		}
		isi, err := os.ReadFile(jalur)
		if err != nil {
			return err
		}
		layar := bagian[0]
		for _, m := range polaJalurAPI.FindAllStringSubmatch(string(isi), -1) {
			j := regexp.MustCompile(`\$\{[^}]*\}`).ReplaceAllString(m[1], "x")
			if i := strings.IndexByte(j, '?'); i >= 0 {
				j = j[:i]
			}
			diperiksa++
			for pemilik, pola := range polaPemilik(t, j) {
				if pemilik == layar {
					continue
				}
				boleh := false
				for _, p := range pola {
					for _, pinjam := range ruteDipinjam[p] {
						if pinjam == layar {
							boleh = true
						}
					}
				}
				if !boleh {
					temuan = append(temuan, layar+" -> "+pemilik+" "+m[1]+" "+strings.Join(pola, ","))
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if diperiksa < 50 {
		t.Fatalf("hanya %d jalur /api/ terbaca dari layar modul; pembacanya yang rusak", diperiksa)
	}
	sort.Strings(temuan)
	if len(temuan) > 0 {
		t.Errorf("panggilan layar ke rute modul lain tanpa ruteDipinjam:\n%s", strings.Join(temuan, "\n"))
	}
}
