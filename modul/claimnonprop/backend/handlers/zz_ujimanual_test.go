//go:build ujimanual

package handlers_test

// SEMENTARA - server uji manual layar Claim Non Prop di atas gudang tiruan (fixture UJI-, nol DEV). Dihapus sesudah
// potret layar; tidak pernah di-commit.

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"nusantarare/modul/claimnonprop/backend/handlers"
	"nusantarare/modul/claimnonprop/backend/models"
	"nusantarare/modul/claimnonprop/backend/services"
)

func TestUjiManualServer(t *testing.T) {
	u := baruUji(t)
	baru := u.buat()
	isi := u.buat()
	u.isiOutstanding(isi)
	ia := u.kirimKeTeknik()
	peran := models.WorkbasketAcceptation
	kode, out := u.aksi(ia, teknik, peran, "SaveToOS", 0, "", nil)
	u.wajib(kode, http.StatusOK, out, "save to OS")
	h := u.g.Halaman(ia)
	h.Setel(models.CD+"AppointedADJID", "UJI-ADJ")
	h.Setel(models.CD+"ConsultantID", "UJI-ADJ")
	u.g.SetelHalaman(ia, h)
	kode, out = u.aksi(ia, teknik, peran, "AddAkseptasiCNP", 0, "", nil)
	u.wajib(kode, http.StatusOK, out, "add akseptasi")
	kode, out = u.aksiP(ia, teknik, peran, services.PermintaanAksi{Aksi: "SetPayableTreatyNP:Acc", Indeks: 1,
		Masukan: map[string]string{models.CD + "Payable": "1", models.JalurAdj(1, "PaymentType"): "1",
			models.CD + "Occupation": "UJI"}})
	u.wajib(kode, http.StatusOK, out, "payable")
	t.Logf("kasus baru %s, outstanding terisi %s, input acceptation %s", baru, isi, ia)
	// Semua kasus dibuat UJI-ADMIN; Input Acceptation dipegang workbasket -> pelaku uji memegang keduanya.
	for _, k := range []string{isi, baru} {
		kk := u.g.Kasus[k]
		kk.PembuatID = teknik
		u.g.Kasus[k] = kk
	}

	tulis := func(w http.ResponseWriter, v any) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(v)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/auth/saya", func(w http.ResponseWriter, _ *http.Request) {
		tulis(w, map[string]any{"akunId": teknik, "nama": "UJI TEKNIK", "peran": []string{peran}, "organisasi": "UJI",
			"divisi": "UJI", "unit": "UJI", "wajibGantiSandi": false, "menu": []string{"claimnonprop"}})
	})
	mux.HandleFunc("GET /api/menu", func(w http.ResponseWriter, _ *http.Request) {
		tulis(w, map[string]any{"golongan": []any{map[string]any{"kode": "KLAIM", "modul": []any{map[string]any{
			"kode": "claimnonprop", "label": "Claim Non Prop", "modul": "claimnonprop", "urutan": 1, "dimigrasi": true}}}}})
	})
	mux.HandleFunc("GET /api/modul-aktif", func(w http.ResponseWriter, _ *http.Request) {
		tulis(w, map[string]any{"modul": []string{"claimnonprop"}})
	})
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) { tulis(w, map[string]any{"ok": true}) })
	modul := u.srv
	suntik := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("X-Pelaku", teknik)
		r.Header.Set("X-Peran", peran)
		modul.ServeHTTP(w, r)
	})
	mux.Handle(handlers.Prefix+"/", suntik)
	mux.Handle(handlers.Prefix, suntik)
	srv := &http.Server{Addr: "127.0.0.1:18098", Handler: mux, ReadHeaderTimeout: 10 * time.Second}
	t.Log("server uji manual di http://127.0.0.1:18098")
	if err := srv.ListenAndServe(); err != nil {
		t.Fatal(err)
	}
}
