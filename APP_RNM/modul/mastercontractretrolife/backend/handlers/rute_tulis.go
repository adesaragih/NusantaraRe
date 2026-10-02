package handlers

// Rute tulis modul Master Contract Retro Life.
//
//	POST /api/master-contract-retro-life/tahun          tahun BARU - `End Period` → `Save` (paket 2)
//	PUT  /api/master-contract-retro-life/tahun/{id}     ubah - `Edit` → `Save` (paket 2)
//	POST /api/master-contract-retro-life/tahun/{id}/kontrak  kontrak BARU - `Add` → `Save` (paket 3)
//	PUT  /api/master-contract-retro-life/kontrak/{id}   ubah kontrak - `Edit` → `Save` (paket 3)
//	POST /api/master-contract-retro-life/kontrak/{id}/reinsurer  reinsurer BARU (paket 4)
//	PUT  /api/master-contract-retro-life/reinsurer/{id} ubah reinsurer (paket 4)
//	POST /api/master-contract-retro-life/reinsurer/{id}/security  security BARU (paket 5)
//	PUT  /api/master-contract-retro-life/security/{id}  ubah security (paket 5)
//	POST /api/master-contract-retro-life/kontrak/{id}/business  business BARU (paket 6)
//	PUT  /api/master-contract-retro-life/business/{id}  ubah business (paket 6)
//	GET  /api/master-contract-retro-life/business/{id}/salin-semua  pratinjau `Copy to all Reinstype`
//	POST /api/master-contract-retro-life/business/{id}/salin-semua  jalankan, badan {"sasaran": [...]}
//	GET  /api/master-contract-retro-life/{kontrak|reinsurer|security|business}/{id}/dampak-hapus  isi popup
//	DELETE /api/master-contract-retro-life/{kontrak|reinsurer|security|business}/{id}  badan {"dampak": {...}}
//
// ⛔ POST dan PUT terpisah walau Pega punya satu `Save` ber-upsert: identitas
// baris baru tidak pernah datang dari klien (ADR-0006), dan badan PUT yang
// membawa `id` berbeda dari jalurnya DITOLAK, bukan salah satu dipilih diam-diam.
// ⛔ Tahun treaty tidak punya rute DELETE (tahun abadi).

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/mastercontractretrolife/backend/models"
	"nusantarare/modul/mastercontractretrolife/backend/services"
)

// batasBadan - badan permintaan terbesar yang dibaca (64 KiB, sama dengan Treaty Contract Out);
// isian terpanjang yang sah (`R/I RATE` 1000 byte) jauh di bawahnya.
const batasBadan = 64 << 10

// bacaBadan mengurai badan JSON SESUDAH identitas diperiksa (401 lebih dulu - permintaan tanpa
// identitas tidak pernah membuat server menampung badannya); terlalu besar = 413, rusak = 400.
// Mengembalikan true bila permintaan SUDAH dijawab.
func bacaBadan(w http.ResponseWriter, r *http.Request, p inti.Pelaku, ke any) bool {
	if jawabGalat(w, inti.WajibIdentitas(p)) {
		return true
	}
	err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadan)).Decode(ke)
	var besar *http.MaxBytesError
	switch {
	case err == nil:
		return false
	case errors.As(err, &besar):
		galat.Tulis(w, http.StatusRequestEntityTooLarge, fmt.Sprintf("request body is larger than %d KiB", batasBadan>>10))
	default:
		galat.Tulis(w, http.StatusBadRequest, "request body is not valid JSON")
	}
	return true
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
		if bacaBadan(w, r, p, &m) {
			return
		}
		hasil, err := l.SimpanTahun(r.Context(), p, m, true)
		tulis(w, hasil, err)
	})
	pasang("PUT "+Prefix+"/tahun/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m services.TahunMasuk
		if bacaBadan(w, r, p, &m) || idJalur(w, r, &m.ID) {
			return
		}
		hasil, err := l.SimpanTahun(r.Context(), p, m, false)
		tulis(w, hasil, err)
	})
	// Kontrak (paket 3): `Add` → `Save` di bawah tahun; `Edit` → `Save` per kontrak.
	pasang("POST "+Prefix+"/tahun/{id}/kontrak", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m services.KontrakMasuk
		if bacaBadan(w, r, p, &m) || tolakIDBaru(w, m.ID) {
			return
		}
		hasil, err := l.SimpanKontrak(r.Context(), p, r.PathValue("id"), m)
		tulis(w, hasil, err)
	})
	pasang("PUT "+Prefix+"/kontrak/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m services.KontrakMasuk
		if bacaBadan(w, r, p, &m) || idJalur(w, r, &m.ID) {
			return
		}
		hasil, err := l.SimpanKontrak(r.Context(), p, "", m)
		tulis(w, hasil, err)
	})
	// Reinsurer (paket 4): `Add` → `Save` di bawah kontrak; `Edit` → `Save`.
	pasang("POST "+Prefix+"/kontrak/{id}/reinsurer", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m services.ReinsurerMasuk
		if bacaBadan(w, r, p, &m) || tolakIDBaru(w, m.ID) {
			return
		}
		hasil, err := l.SimpanReinsurer(r.Context(), p, r.PathValue("id"), m)
		tulis(w, hasil, err)
	})
	pasang("PUT "+Prefix+"/reinsurer/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m services.ReinsurerMasuk
		if bacaBadan(w, r, p, &m) || idJalur(w, r, &m.ID) {
			return
		}
		hasil, err := l.SimpanReinsurer(r.Context(), p, "", m)
		tulis(w, hasil, err)
	})
	// Security (paket 5): `Add` → `Save` di bawah reinsurer; `Edit` → `Save`.
	pasang("POST "+Prefix+"/reinsurer/{id}/security", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m services.SecurityMasuk
		if bacaBadan(w, r, p, &m) || tolakIDBaru(w, m.ID) {
			return
		}
		hasil, err := l.SimpanSecurity(r.Context(), p, r.PathValue("id"), m)
		tulis(w, hasil, err)
	})
	pasang("PUT "+Prefix+"/security/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m services.SecurityMasuk
		if bacaBadan(w, r, p, &m) || idJalur(w, r, &m.ID) {
			return
		}
		hasil, err := l.SimpanSecurity(r.Context(), p, "", m)
		tulis(w, hasil, err)
	})
	// Business (paket 6): `Add` → `Save` di bawah kontrak; `Edit` → `Save`.
	pasang("POST "+Prefix+"/kontrak/{id}/business", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m services.BusinessMasuk
		if bacaBadan(w, r, p, &m) || tolakIDBaru(w, m.ID) {
			return
		}
		hasil, err := l.SimpanBusiness(r.Context(), p, r.PathValue("id"), m)
		tulis(w, hasil, err)
	})
	pasang("PUT "+Prefix+"/business/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m services.BusinessMasuk
		if bacaBadan(w, r, p, &m) || idJalur(w, r, &m.ID) {
			return
		}
		hasil, err := l.SimpanBusiness(r.Context(), p, "", m)
		tulis(w, hasil, err)
	})
	// `Copy to all Reinstype` (paket 6): pratinjau lalu konfirmasi berdaftar sasaran.
	pasang("GET "+Prefix+"/business/{id}/salin-semua", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		hasil, err := l.PratinjauSalinSemua(r.Context(), p, r.PathValue("id"))
		tulis(w, hasil, err)
	})
	pasang("POST "+Prefix+"/business/{id}/salin-semua", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m struct {
			Sasaran []string `json:"sasaran"`
		}
		if bacaBadan(w, r, p, &m) {
			return
		}
		hasil, err := l.SalinSemua(r.Context(), p, r.PathValue("id"), m.Sasaran)
		tulis(w, hasil, err)
	})
	// Hapus berjenjang (paket 7): popup lebih dulu, lalu DELETE berbadan cacahan popup.
	// ⛔ Nol rute untuk tahun treaty.
	for _, jenis := range []services.JenisHapus{services.HapusKontrak, services.HapusReinsurer,
		services.HapusSecurity, services.HapusBusiness} {
		jenis := jenis
		pasang("GET "+Prefix+"/"+string(jenis)+"/{id}/dampak-hapus", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
			hasil, err := l.DampakHapus(r.Context(), p, jenis, r.PathValue("id"))
			tulis(w, hasil, err)
		})
		pasang("DELETE "+Prefix+"/"+string(jenis)+"/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
			var m struct {
				Dampak *models.Dampak `json:"dampak"`
			}
			if bacaBadan(w, r, p, &m) {
				return
			}
			if m.Dampak == nil {
				galat.Tulis(w, http.StatusBadRequest, "delete requires the confirmed impact (dampak) from the confirmation popup")
				return
			}
			hasil, err := l.Hapus(r.Context(), p, jenis, r.PathValue("id"), *m.Dampak)
			tulis(w, hasil, err)
		})
	}
}

// tolakIDBaru - baris baru tidak boleh membawa id (ADR-0006); 400 dan true.
func tolakIDBaru(w http.ResponseWriter, id string) bool {
	if id != "" {
		galat.Tulis(w, http.StatusBadRequest, services.Pesan(services.ErrIDDariKlien))
		return true
	}
	return false
}
