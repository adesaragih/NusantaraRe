package handlers

import (
	"encoding/json"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/treatyin/backend/services"
)

// daftarkanSimpan - SATU-SATUNYA jalur tulis form Treaty In (keputusan
// pemilik proses 6–7 Oktober 2026): isian masuk basis data HANYA lewat
// kedua rute ini, yang dipanggil tombolnya — tidak pernah dari `onChange`.
//
//	POST /kontrak/simpan  tombol Save
//	POST /kontrak/kirim   Submit (`aksi: submit`), Actions (`akseptasi` +
//	                      `pilihan`), Decline offer (`decline`)
//	POST /kontrak/revisi  tombol Revision daftar kontrak
//	POST /penyesuaian/simpan  Save layar Adjustment
//	POST /penyesuaian/kirim   Submit / Actions layar Adjustment
//	POST /penyesuaian/hapus   Decline offer layar Adjustment
func daftarkanSimpan(pasang func(string, rute)) {
	pasang("POST "+Prefix+"/kontrak/simpan", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m services.MasukanSimpan
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "request body must be valid JSON")
			return
		}
		hasil, err := l.SimpanKontrak(r.Context(), p, m)
		tulis(w, hasil, err)
	})
	pasang("POST "+Prefix+"/kontrak/kirim", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m services.MasukanKirim
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "request body must be valid JSON")
			return
		}
		hasil, err := l.KirimKontrak(r.Context(), p, m)
		tulis(w, hasil, err)
	})
	// Tombol `Revision` daftar kontrak — `SetTreatyIn_Act` ber-revisionstate.
	pasang("POST "+Prefix+"/kontrak/revisi", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m services.MasukanRevisi
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "request body must be valid JSON")
			return
		}
		hasil, err := l.MulaiRevisi(r.Context(), p, m)
		tulis(w, hasil, err)
	})

	// ⭐ Layar ADJUSTMENT (modul Treaty In Adjustment) — penulis pendaratannya
	// di modul ini, dan penjaga arsitektur melarang modul saling impor.
	pasang("POST "+Prefix+"/penyesuaian/simpan", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m services.MasukanPenyesuaian
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "request body must be valid JSON")
			return
		}
		hasil, err := l.SimpanPenyesuaian(r.Context(), p, m)
		tulis(w, hasil, err)
	})
	pasang("POST "+Prefix+"/penyesuaian/kirim", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m services.MasukanKirimPenyesuaian
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "request body must be valid JSON")
			return
		}
		hasil, err := l.KirimPenyesuaian(r.Context(), p, m)
		tulis(w, hasil, err)
	})
	pasang("POST "+Prefix+"/penyesuaian/hapus", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m services.MasukanHapusPenyesuaian
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "request body must be valid JSON")
			return
		}
		hasil, err := l.HapusPenyesuaian(r.Context(), p, m)
		tulis(w, hasil, err)
	})
}
