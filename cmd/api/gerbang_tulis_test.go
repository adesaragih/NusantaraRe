package main

// Gerbang tulis hak menu LIHAT (`M_LOGIN_GO_MENU.HAK`, keputusan work owner 04-10-2026) - TANPA Oracle.
//
// Yang dijaga: modul yang mendaftar (`inti.Pendaftaran.HakLihat`) menolak tulis bagi pemegang menu View only kecuali
// pola yang dibebaskannya; GET tetap lewat; superadmin ber-hak View only ikut ditolak; modul yang TIDAK mendaftar
// tidak tersentuh; dan
// setiap pola bebas memang rute modul pemiliknya.

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/daftar"
	"nusantarare/inti/backend/menu"
)

// kirimLihat - satu permintaan dari sesi pemegang `menuAkun`, dengan `lihat` ber-hak View only.
func kirimLihat(h http.Handler, metode, jalur string, menuAkun, lihat []string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(metode, jalur, strings.NewReader("{}"))
	ctx := inti.DenganMenuLihat(inti.DenganAksesMenu(context.Background(), menuAkun), lihat)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r.WithContext(ctx))
	return w
}

func viewOnly(w *httptest.ResponseRecorder) bool {
	return w.Code == http.StatusForbidden && strings.Contains(w.Body.String(), "View only")
}

func TestGerbangTulisHakLihat(t *testing.T) {
	h := muxGerbang(t, false)
	acc := []string{"accounts"}
	if w := kirimLihat(h, "POST", "/api/accounts", acc, acc); !viewOnly(w) {
		t.Errorf("tulis View only: %d %s", w.Code, w.Body)
	}
	if w := kirimLihat(h, "GET", "/api/accounts", acc, acc); viewOnly(w) {
		t.Errorf("baca View only ditolak: %d", w.Code)
	}
	if w := kirimLihat(h, "POST", "/api/accounts", acc, nil); viewOnly(w) {
		t.Errorf("tulis PENUH ditolak: %d", w.Code)
	}
	super := []string{"accounts", menu.KodeKelolaUser}
	if w := kirimLihat(h, "POST", "/api/accounts", super, acc); !viewOnly(w) {
		t.Errorf("superadmin View only harus ikut ditolak (05-10-2026, B): %d", w.Code)
	}
	ag := []string{"aggregate"}
	if w := kirimLihat(h, "POST", "/api/aggregate/pratinjau", ag, ag); viewOnly(w) {
		t.Errorf("pratinjau (bebas) ditolak: %d", w.Code)
	}
	if w := kirimLihat(h, "POST", "/api/aggregate/simpan", ag, ag); !viewOnly(w) {
		t.Errorf("simpan Aggregate View only: %d %s", w.Code, w.Body)
	}
	bdx := []string{"bordereaux"}
	if w := kirimLihat(h, "POST", "/api/bordereaux/berkas/X/submit", bdx, bdx); viewOnly(w) {
		t.Errorf("submit Bordereaux (bebas, dijaga layanan) ditolak: %d", w.Code)
	}
	// Modul yang belum mendaftar tidak tersentuh, walau barisnya (keliru) LIHAT.
	tco := []string{"treatycontractout"}
	if w := kirimLihat(h, "POST", "/api/treaty-contract-out/tahun", tco, tco); viewOnly(w) {
		t.Errorf("modul tanpa HakLihat ikut ditolak: %d", w.Code)
	}
}

// Tahap 1 (keputusan work owner 04-10-2026): enam modul selesai yang mendaftar, ditambah Adjuster Consultant,
// Business Group, Treaty Group, Treaty Group OJK, Treaty Exchange Yearly, Treaty Description, dan Reinsurance Type
// (05-10-2026).
func TestModulHakLihat(t *testing.T) {
	mau := []string{"accounts", "adjusterconsultant", "aggregate", "bordereaux", "businessgroup", "companydetail",
		"marketingofficer", "masterproductnamelife", "reinsurancetype", "treatydescription", "treatyexchangeyearly",
		"treatygroup", "treatygroupojk"}
	if got := kodeHakLihat(daftar.HakLihat()); strings.Join(got, ",") != strings.Join(mau, ",") {
		t.Errorf("modul HakLihat %v, mau %v", got, mau)
	}
}

// Setiap pola bebas adalah rute TULIS modul pemiliknya - salah ketik tidak boleh diam-diam menolak.
func TestPolaBebasHakLihatTerdaftar(t *testing.T) {
	for modul, bebas := range daftar.HakLihat() {
		for _, pola := range bebas {
			bagian := strings.SplitN(pola, " ", 2)
			if len(bagian) != 2 || bagian[0] == http.MethodGet {
				t.Errorf("%s: pola bebas %q bukan rute tulis", modul, pola)
				continue
			}
			jalur := regexp.MustCompile(`\{[^}]+\}`).ReplaceAllString(bagian[1], "x")
			ada := false
			for _, p := range polaPemilik(t, jalur)[modul] {
				ada = ada || p == pola
			}
			if !ada {
				t.Errorf("%s: pola bebas %q tidak didaftarkan modul itu", modul, pola)
			}
		}
	}
}
