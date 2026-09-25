// Package handlers memuat HTTP controller.
//
// Arah ketergantungan: handlers -> services -> repository. Paket ini tidak
// pernah mengimpor repository secara langsung, dan tidak pernah memuat aturan
// dagang - aturan tinggal di services.
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

// Kesehatan adalah pemeriksa yang dapat dipanggil Router tanpa membuat
// handlers bergantung pada repository.
type Kesehatan interface {
	Ping(ctx context.Context) error
}

// Router menyusun seluruh rute.
//
// svc boleh nil di Fase 0; ia sudah diterima di sini supaya tiket pertama yang
// punya endpoint tidak perlu mengubah tanda tangan fungsi ini.
func Router(svc *services.Service, kesehatan Kesehatan) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthz(kesehatan))
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
func healthz(kesehatan Kesehatan) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jawab := jawabanSehat{Status: "sehat", Database: "tidak dikonfigurasi"}

		if kesehatan != nil {
			ctx, batal := context.WithTimeout(r.Context(), 3*time.Second)
			defer batal()
			if err := kesehatan.Ping(ctx); err != nil {
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
