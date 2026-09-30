package handlers

// Rute tulis modul Master Contract Retro Life.
//
//	POST /api/master-contract-retro-life/tahun          tahun BARU - `End Period` → `Save` (paket 2)
//	PUT  /api/master-contract-retro-life/tahun/{id}     ubah - `Edit` → `Save` (paket 2)
//
// ⛔ POST dan PUT terpisah walau Pega punya satu `Save` ber-upsert: identitas
// baris baru tidak pernah datang dari klien (ADR-0006), dan badan PUT yang
// membawa `id` berbeda dari jalurnya DITOLAK, bukan salah satu dipilih diam-diam.
// ⛔ Tahun treaty tidak punya rute DELETE (tahun abadi).

import (
	"encoding/json"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/mastercontractretrolife/backend/services"
)

// bacaBadan mengurai badan JSON; gagal = 400 dan true.
func bacaBadan(w http.ResponseWriter, r *http.Request, ke any) bool {
	if err := json.NewDecoder(r.Body).Decode(ke); err != nil {
		galat.Tulis(w, http.StatusBadRequest, "request body is not valid JSON")
		return true
	}
	return false
}

// idJalur menyamakan id badan dengan id jalur (PUT); beda = 400 dan true.
func idJalur(w http.ResponseWriter, r *http.Request, id *string) bool {
	dariJalur := r.PathValue("id")
	if *id != "" && *id != dariJalur {
		galat.Tulis(w, http.StatusBadRequest, "id in the body differs from id in the path")
		return true
	}
	*id = dariJalur
	return false
}

func daftarkanTulis(pasang func(string, rute)) {
	pasang("POST "+Prefix+"/tahun", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m services.TahunMasuk
		if bacaBadan(w, r, &m) {
			return
		}
		hasil, err := l.SimpanTahun(r.Context(), p, m, true)
		tulis(w, hasil, err)
	})
	pasang("PUT "+Prefix+"/tahun/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m services.TahunMasuk
		if bacaBadan(w, r, &m) || idJalur(w, r, &m.ID) {
			return
		}
		hasil, err := l.SimpanTahun(r.Context(), p, m, false)
		tulis(w, hasil, err)
	})
}
