// Package handlers memuat HTTP controller.
//
// Arah ketergantungan: handlers -> services -> repository. Paket ini tidak
// mengimpor repository, dan tidak pernah memegang koneksi database: seluruh
// pertanyaan ke Oracle - termasuk pemeriksaan kesehatan - lewat services.
// Nol aturan dagang di sini.
//
// [usulan] Router memakai net/http bawaan Go 1.22 (pola "GET /path"),
// sehingga nol ketergantungan pihak ketiga.
package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"nusantarare/internal/services"
)

// Router menyusun seluruh rute.
func Router(svc *services.Service) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthz(svc))
	return mux
}

type jawabanSehat struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

// healthz menjawab tanpa menyentuh aturan dagang mana pun.
//
// Ia tetap menjawab 200 ketika Oracle belum dikonfigurasi: Fase 0 harus dapat
// dijalankan tanpa instance, dan keadaan database dilaporkan apa adanya di
// dalam badan jawaban, bukan disembunyikan.
func healthz(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jawab := jawabanSehat{Status: "sehat", Database: "tidak dikonfigurasi"}

		if svc.PunyaDatabase() {
			ctx, batal := context.WithTimeout(r.Context(), 3*time.Second)
			defer batal()
			if err := svc.CekKesehatan(ctx); err != nil {
				jawab.Database = "tidak terjangkau"
			} else {
				jawab.Database = "terjangkau"
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(jawab)
	}
}
