package main

// Uji perakitan modul aktif - TANPA Oracle. Refactor bentuk B paket 6.
//
// Yang dijaga: MODUL_AKTIF memilih modul yang dipasang - rute modul nonaktif
// tidak ada (404), daftar modul aktif dilaporkan GET /api/modul-aktif, dan
// migrasi TIDAK ikut terpilih: skema selalu dari SEMUA modul terdaftar.

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"nusantarare/inti"
	"nusantarare/inti/config"
	"nusantarare/inti/migrasi"
	"nusantarare/modul"
)

// ruteContoh - satu rute GET milik setiap modul terdaftar.
var ruteContoh = map[string]string{
	"claimlife":         "/api/klaim-life",
	"premiumlistlife":   "/api/polis-life",
	"komiteclaimlife":   "/api/komite",
	"treatycontractout": "/api/treaty-contract-out/tahun",
}

func muxUji(t *testing.T, diminta []string) (http.Handler, []inti.Modul) {
	t.Helper()
	terdaftar := modul.Rakit(inti.NewDasar(nil), config.Config{}, func(string) {})
	aktif, err := pilihModulAktif(terdaftar, modul.NamaLama, diminta)
	if err != nil {
		t.Fatal(err)
	}
	return rakitMux(inti.NewDasar(nil), terdaftar, aktif, false), aktif
}

func kode(mux http.Handler, jalur string) int {
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest(http.MethodGet, jalur, nil))
	return w.Code
}

func TestModulNonaktifRutenya404(t *testing.T) {
	mux, _ := muxUji(t, []string{"claimlife"})
	for nama, jalur := range ruteContoh {
		k := kode(mux, jalur)
		if nama == "claimlife" && k == http.StatusNotFound {
			t.Errorf("%s aktif tetapi %s menjawab 404", nama, jalur)
		}
		if nama != "claimlife" && k != http.StatusNotFound {
			t.Errorf("%s NONAKTIF tetapi %s menjawab %d, mau 404", nama, jalur, k)
		}
	}
	// Rute milik aplikasi tetap ada walau modul mana pun nonaktif.
	if k := kode(mux, "/healthz"); k != http.StatusOK {
		t.Errorf("/healthz menjawab %d", k)
	}
	// ⛔ 404 modul nonaktif berbadan JSON `{galat}` yang menyebut modulnya:
	// teks `404 page not found` terbaca frontend sebagai "backend mati".
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/polis-life/ringkas?nomorPolis=UJI-1", nil))
	var badan struct {
		Galat string `json:"galat"`
	}
	if w.Code != http.StatusNotFound || !strings.HasPrefix(w.Header().Get("Content-Type"), "application/json") ||
		json.NewDecoder(w.Body).Decode(&badan) != nil || !strings.Contains(badan.Galat, "premiumlistlife") {
		t.Errorf("rute modul nonaktif: kode %d, Content-Type %q, galat %q", w.Code, w.Header().Get("Content-Type"), badan.Galat)
	}
	// Jalur yang bukan milik modul MANA PUN tetap 404 bawaan mux, seperti dulu.
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/tidak-ada", nil))
	if w.Code != http.StatusNotFound || strings.HasPrefix(w.Header().Get("Content-Type"), "application/json") {
		t.Errorf("jalur tak dikenal: kode %d, Content-Type %q", w.Code, w.Header().Get("Content-Type"))
	}
}

func TestModulAktifKosongBerartiSemua(t *testing.T) {
	mux, aktif := muxUji(t, nil)
	if len(aktif) != len(ruteContoh) {
		t.Fatalf("%d modul aktif, mau %d (semua terdaftar)", len(aktif), len(ruteContoh))
	}
	for nama, jalur := range ruteContoh {
		if k := kode(mux, jalur); k == http.StatusNotFound {
			t.Errorf("MODUL_AKTIF kosong, tetapi rute %s (%s) menjawab 404", jalur, nama)
		}
	}
}

func TestModulAktifDilaporkan(t *testing.T) {
	mux, _ := muxUji(t, []string{"komiteclaimlife", "claimlife"})
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/modul-aktif", nil))
	if w.Code != http.StatusOK || !strings.HasPrefix(w.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("kode %d, Content-Type %q", w.Code, w.Header().Get("Content-Type"))
	}
	var badan struct {
		Modul []string `json:"modul"`
	}
	if err := json.NewDecoder(w.Body).Decode(&badan); err != nil {
		t.Fatal(err)
	}
	// Urutan DAFTAR, bukan urutan env.
	if mau := []string{"claimlife", "komiteclaimlife"}; !reflect.DeepEqual(badan.Modul, mau) {
		t.Errorf("modul aktif = %v, mau %v", badan.Modul, mau)
	}
}

// GET /api/menu milik aplikasi: terpasang walau hanya satu modul aktif, dan
// tanpa Oracle menjawab 503 BERGALAT - bukan 404 modul nonaktif, bukan menu
// kosong diam-diam.
func TestRuteMenuTerpasangDanGagalTerang(t *testing.T) {
	for _, diminta := range [][]string{nil, {"treatycontractout"}} {
		mux, _ := muxUji(t, diminta)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/menu", nil))
		var badan struct {
			Galat string `json:"galat"`
		}
		if err := json.NewDecoder(w.Body).Decode(&badan); err != nil {
			t.Fatalf("MODUL_AKTIF %v: badan bukan JSON: %v", diminta, err)
		}
		if w.Code != http.StatusServiceUnavailable || !strings.Contains(badan.Galat, "M_NAV_MENU") {
			t.Errorf("MODUL_AKTIF %v: kode %d, galat %q - mau 503 yang menyebut M_NAV_MENU", diminta, w.Code, badan.Galat)
		}
	}
}

func TestModulAktifTakDikenalDitolak(t *testing.T) {
	terdaftar := modul.Rakit(inti.NewDasar(nil), config.Config{}, func(string) {})
	_, err := pilihModulAktif(terdaftar, modul.NamaLama, []string{"claimlife", "klaim"})
	if err == nil || !strings.Contains(err.Error(), `"klaim"`) {
		t.Fatalf("nama tak dikenal tidak ditolak dengan menyebut namanya: %v", err)
	}
	// Nama lama (sebelum tabel nama modul 30-09-2026) ditolak dengan kalimat
	// yang menyebut nama BARUnya - bukan diterima diam-diam, bukan "tak dikenal".
	// Diperiksa atas peta YANG DIPAKAI (`modul.NamaLama`), dan setiap nama
	// barunya wajib modul terdaftar: menyuruh operator memakai nama yang juga
	// ditolak lebih buruk daripada tidak menjawab.
	dikenal := map[string]bool{}
	for _, m := range terdaftar {
		dikenal[m.Nama()] = true
	}
	if len(modul.NamaLama) < 3 {
		t.Fatalf("peta nama lama memuat %d nama, mau sekurangnya 3", len(modul.NamaLama))
	}
	for lama, baru := range modul.NamaLama {
		if !dikenal[baru] || dikenal[lama] {
			t.Errorf("nama lama %q -> %q: pengganti terdaftar=%v, nama lama masih terdaftar=%v",
				lama, baru, dikenal[baru], dikenal[lama])
		}
		_, err := pilihModulAktif(terdaftar, modul.NamaLama, []string{"claimlife", lama})
		if err == nil || !strings.Contains(err.Error(), "nama modul lama") ||
			!strings.Contains(err.Error(), `"`+baru+`"`) {
			t.Errorf("nama lama %q: galat %v, mau menyebut %q", lama, err, baru)
		}
	}
}

// Migrasi tidak ikut MODUL_AKTIF: setiap berkas maju di folder migrations/
// modul mana pun ada di daftar pelari - dihitung dari disk, bukan dari angka.
func TestMigrasiTetapLengkapSaatModulNonaktif(t *testing.T) {
	if _, aktif := muxUji(t, []string{"treatycontractout"}); len(aktif) != 1 {
		t.Fatal("prasyarat: hanya satu modul aktif")
	}
	langkah, err := migrasi.Daftar(false, modul.SumberMigrasi()...)
	if err != nil {
		t.Fatal(err)
	}
	var dariPelari []string
	for _, l := range langkah {
		dariPelari = append(dariPelari, l.Nama)
	}
	var dariDisk []string
	berkas, err := filepath.Glob(filepath.Join("..", "..", "modul", "*", "migrations", "*.sql"))
	if err != nil {
		t.Fatal(err)
	}
	// Migrasi lintas modul milik `inti` (900-949, M_NAV_MENU) ikut dihitung.
	milikInti, err := filepath.Glob(filepath.Join("..", "..", "inti", "migrations", "*.sql"))
	if err != nil {
		t.Fatal(err)
	}
	if len(milikInti) == 0 {
		t.Fatal("nol berkas migrasi inti di disk; M_NAV_MENU (900) hilang")
	}
	berkas = append(berkas, milikInti...)
	for _, b := range berkas {
		if n := filepath.Base(b); !strings.HasSuffix(n, "_down.sql") {
			dariDisk = append(dariDisk, n)
		}
	}
	sort.Strings(dariDisk)
	if len(dariDisk) < 20 {
		t.Fatalf("hanya %d berkas migrasi di disk; pembacanya yang rusak", len(dariDisk))
	}
	if !reflect.DeepEqual(dariPelari, dariDisk) {
		t.Errorf("pelari menjalankan %d langkah, disk memuat %d:\npelari %v\ndisk   %v",
			len(dariPelari), len(dariDisk), dariPelari, dariDisk)
	}
	// Prasyarat POSITIF: pemeriksaan "tanpa folder migrations" di bawah lulus
	// dengan sendirinya bila nama folder modulnya salah.
	if _, err := os.Stat(filepath.Join("..", "..", "modul", "treatycontractout")); err != nil {
		t.Fatalf("folder modul Treaty Contract Out tidak ditemukan: %v", err)
	}
	if _, err := os.Stat(filepath.Join("..", "..", "modul", "treatycontractout", "migrations")); err == nil {
		t.Error("modul treatycontractout kini punya folder migrations - tco4 menyatakan nol tabel baru")
	}
}
