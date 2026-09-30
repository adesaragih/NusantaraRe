package menu

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"nusantarare/inti"
	"nusantarare/inti/db"
)

// barisUji - potongan isi awal 900, SENGAJA diacak urutannya: pohon tidak
// boleh bergantung pada ORDER BY pembacanya.
func barisUji() []Baris {
	return []Baris{
		{ID: 22, IndukID: 2, Kode: "register", Label: "Register", Golongan: "KLAIM", Modul: "claimlife", Urutan: 2, Dimigrasi: true},
		{ID: 5, Kode: "treatycontractout", Label: "Treaty Contract Out", Golongan: "MASTER", Modul: "treatycontractout", Urutan: 3, Dimigrasi: true},
		{ID: 1, Kode: "claimfacin", Label: "Claim Fac In", Golongan: "KLAIM", Modul: "claimfacin", Urutan: 1},
		{ID: 2, Kode: "claimlife", Label: "Claim Life", Golongan: "KLAIM", Modul: "claimlife", Urutan: 2, Dimigrasi: true},
		{ID: 21, IndukID: 2, Kode: "inbox", Label: "Inbox Claim Life", Golongan: "KLAIM", Modul: "claimlife", Urutan: 1, Dimigrasi: true},
		{ID: 3, Kode: "nbfacin", Label: "NB FacIn", Golongan: "FACULTATIVE", Modul: "nbfacin", Urutan: 1},
		{ID: 4, Kode: "premiumlistlife", Label: "PremiumList Life", Golongan: "TREATY", Modul: "premiumlistlife", Urutan: 5, Dimigrasi: true},
		{ID: 41, IndukID: 4, Kode: "premiumlist", Label: "PremiumList", Golongan: "TREATY", Modul: "premiumlistlife", Urutan: 1, Dimigrasi: true},
		{ID: 51, IndukID: 5, Kode: "tco-tahun", Label: "Treaty Contract Out", Golongan: "MASTER", Modul: "treatycontractout", Urutan: 1, Dimigrasi: true},
		{ID: 6, Kode: "nbtreatyin", Label: "NB Treaty In", Golongan: "TREATY", Modul: "nbtreatyin", Urutan: 1},
	}
}

var semuaAktif = []string{"claimlife", "premiumlistlife", "komiteclaimlife", "treatycontractout"}

// ringkas menulis pohon sebagai teks satu baris per simpul, supaya selisihnya
// terbaca di pesan uji.
func ringkas(m Menu) []string {
	var out []string
	for _, g := range m.Golongan {
		out = append(out, g.Kode)
		for _, k := range g.Kelompok {
			out = append(out, "  "+k.Kode)
			for _, b := range k.Butir {
				out = append(out, "    "+b.Kode)
			}
		}
	}
	return out
}

func TestSusunGolonganKelompokButirBerurutan(t *testing.T) {
	dapat := ringkas(Susun(barisUji(), semuaAktif))
	mau := []string{
		"TREATY", "  nbtreatyin", "  premiumlistlife", "    premiumlist",
		"FACULTATIVE", "  nbfacin",
		"KLAIM", "  claimfacin", "  claimlife", "    inbox", "    register",
		"MASTER", "  treatycontractout", "    tco-tahun",
	}
	if !reflect.DeepEqual(dapat, mau) {
		t.Errorf("pohon:\n%s\nmau:\n%s", strings.Join(dapat, "\n"), strings.Join(mau, "\n"))
	}
}

// Butir milik modul yang TIDAK ada di MODUL_AKTIF tidak dikirim; kelompoknya
// tetap, dengan butir kosong - frontend yang memutuskan cara menampilkannya.
func TestSusunButirModulNonaktifTidakDikirim(t *testing.T) {
	m := Susun(barisUji(), []string{"premiumlistlife"})
	for _, g := range m.Golongan {
		for _, k := range g.Kelompok {
			for _, b := range k.Butir {
				if b.Modul != "premiumlistlife" {
					t.Errorf("butir %s milik modul nonaktif %s ikut dikirim", b.Kode, b.Modul)
				}
			}
			if k.Kode == "claimlife" && (len(k.Butir) != 0 || !k.Dimigrasi) {
				t.Errorf("kelompok claimlife: butir %v, dimigrasi %v - mau tetap ada, tanpa butir", k.Butir, k.Dimigrasi)
			}
		}
	}
	if n := len(ringkas(m)); n != 11 {
		t.Errorf("simpul %d, mau 11 (14 dikurangi inbox, register, tco-tahun)", n)
	}
}

// Butir yang induknya tidak terbaca (induk nonaktif) atau yang induknya BUTIR
// (tingkat ketiga) tidak dikirim - tabel ini dua tingkat.
func TestSusunButirYatimDanTingkatKetigaDibuang(t *testing.T) {
	baris := append(barisUji(),
		Baris{ID: 90, IndukID: 99, Kode: "yatim", Label: "Yatim", Golongan: "KLAIM", Modul: "claimlife", Urutan: 1, Dimigrasi: true},
		Baris{ID: 91, IndukID: 21, Kode: "cucu", Label: "Cucu", Golongan: "KLAIM", Modul: "claimlife", Urutan: 1, Dimigrasi: true},
	)
	for _, s := range ringkas(Susun(baris, semuaAktif)) {
		if strings.Contains(s, "yatim") || strings.Contains(s, "cucu") {
			t.Errorf("simpul %q ikut dikirim", strings.TrimSpace(s))
		}
	}
}

// JSON: daftar kosong adalah `[]`, bukan `null` - frontend membaca panjangnya.
func TestSusunDaftarKosongBukanNull(t *testing.T) {
	isi, err := json.Marshal(Susun(barisUji(), nil))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(isi), "null") {
		t.Errorf("JSON memuat null: %s", isi)
	}
	isi, _ = json.Marshal(Susun(nil, nil))
	if string(isi) != `{"golongan":[]}` {
		t.Errorf("menu kosong = %s", isi)
	}
}

// Titik sambung akses per akun: hari ini meneruskan SEMUA, untuk pelaku siapa pun.
func TestSaringMenuUntukPelakuMeneruskanSemua(t *testing.T) {
	m := Susun(barisUji(), semuaAktif)
	for _, p := range []inti.Pelaku{{}, {AkunID: "UJI-1", Peran: []string{"UjiPeran"}}} {
		if !reflect.DeepEqual(SaringMenuUntukPelaku(p, m), m) {
			t.Errorf("pelaku %+v: menu berubah", p)
		}
	}
}

func TestSQLMenuHanyaBarisAktifBerurutan(t *testing.T) {
	q := sqlMenu("SKEMAUJI.M_NAV_MENU")
	if benderaYa != "1" {
		t.Errorf("bendera aktif %q, mau \"1\" (konvensi data warisan)", benderaYa)
	}
	for _, mau := range []string{"FROM SKEMAUJI.M_NAV_MENU", "WHERE STATUS_AKTIF = :1", "ORDER BY GROUPMENU, URUTAN, ID"} {
		if !strings.Contains(q, mau) {
			t.Errorf("SQL tidak memuat %q: %s", mau, q)
		}
	}
	if err := db.PeriksaSQL(q); err != nil {
		t.Error(err)
	}
}

// pembacaUji memenuhi PembacaMenu tanpa Oracle.
type pembacaUji struct {
	baris []Baris
	err   error
}

func (p pembacaUji) Baca(context.Context) ([]Baris, error) { return p.baris, p.err }

func minta(t *testing.T, h http.HandlerFunc) (*httptest.ResponseRecorder, map[string]any) {
	t.Helper()
	w := httptest.NewRecorder()
	h(w, httptest.NewRequest(http.MethodGet, "/api/menu", nil))
	var badan map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &badan); err != nil {
		t.Fatalf("badan bukan JSON: %q", w.Body.String())
	}
	return w, badan
}

func TestRuteMenjawabPohon(t *testing.T) {
	w, badan := minta(t, Rute(pembacaUji{baris: barisUji()}, semuaAktif, false))
	if w.Code != http.StatusOK {
		t.Fatalf("kode %d: %v", w.Code, badan)
	}
	if g, _ := badan["golongan"].([]any); len(g) != 4 {
		t.Errorf("golongan %d, mau 4: %v", len(g), badan)
	}
}

// Tanpa Oracle: 503 berbadan galat - bukan menu kosong diam-diam.
func TestRuteTanpaDatabase503(t *testing.T) {
	w, badan := minta(t, Rute(nil, semuaAktif, false))
	if w.Code != http.StatusServiceUnavailable || badan["galat"] == nil {
		t.Errorf("kode %d, badan %v - mau 503 dengan galat", w.Code, badan)
	}
}

// Tabel belum ada (migrasi 900 belum dijalankan): 503 yang MENYEBUT sebabnya.
func TestRuteTabelBelumAda503MenyebutMigrasi(t *testing.T) {
	err := errors.New("ORA-00942: table or view does not exist")
	w, badan := minta(t, Rute(pembacaUji{err: err}, semuaAktif, false))
	pesan, _ := badan["galat"].(string)
	if w.Code != http.StatusServiceUnavailable || !strings.Contains(pesan, "900") || !strings.Contains(pesan, "-migrate") {
		t.Errorf("kode %d, galat %q - mau 503 yang menyebut migrasi 900 dan -migrate", w.Code, pesan)
	}
}

func TestRuteGagalBaca500(t *testing.T) {
	w, badan := minta(t, Rute(pembacaUji{err: errors.New("ORA-12541: no listener")}, semuaAktif, false))
	pesan, _ := badan["galat"].(string)
	if w.Code != http.StatusInternalServerError || !strings.Contains(pesan, "M_NAV_MENU") {
		t.Errorf("kode %d, galat %q - mau 500 yang menyebut M_NAV_MENU", w.Code, pesan)
	}
}
