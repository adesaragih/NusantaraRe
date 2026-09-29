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
	"claimlife":   "/api/klaim-life",
	"premiumlist": "/api/polis-life",
	"komite":      "/api/komite",
	"treaty":      "/api/treaty-contract-out/tahun",
}

func muxUji(t *testing.T, diminta []string) (http.Handler, []inti.Modul) {
	t.Helper()
	terdaftar := modul.Rakit(inti.NewDasar(nil), config.Config{}, func(string) {})
	aktif, err := pilihModulAktif(terdaftar, diminta)
	if err != nil {
		t.Fatal(err)
	}
	return rakitMux(inti.NewDasar(nil), terdaftar, aktif), aktif
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
		json.NewDecoder(w.Body).Decode(&badan) != nil || !strings.Contains(badan.Galat, "premiumlist") {
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
	mux, _ := muxUji(t, []string{"komite", "claimlife"})
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
	if mau := []string{"claimlife", "komite"}; !reflect.DeepEqual(badan.Modul, mau) {
		t.Errorf("modul aktif = %v, mau %v", badan.Modul, mau)
	}
}

func TestModulAktifTakDikenalDitolak(t *testing.T) {
	_, err := pilihModulAktif(modul.Rakit(inti.NewDasar(nil), config.Config{}, func(string) {}),
		[]string{"claimlife", "klaim"})
	if err == nil || !strings.Contains(err.Error(), `"klaim"`) {
		t.Fatalf("nama tak dikenal tidak ditolak dengan menyebut namanya: %v", err)
	}
}

// Migrasi tidak ikut MODUL_AKTIF: setiap berkas maju di folder migrations/
// modul mana pun ada di daftar pelari - dihitung dari disk, bukan dari angka.
func TestMigrasiTetapLengkapSaatModulNonaktif(t *testing.T) {
	if _, aktif := muxUji(t, []string{"treaty"}); len(aktif) != 1 {
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
	if _, err := os.Stat(filepath.Join("..", "..", "modul", "treaty", "migrations")); err == nil {
		t.Error("modul treaty kini punya folder migrations - tco4 menyatakan nol tabel baru")
	}
}
