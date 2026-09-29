package handlers

// Pintu HTTP kurs USD -> IDR - tiket 11.
//
//	GET /api/treaty-contract-out/tahun/{id}/kurs                                   `testingKurs` (Show b6138/b9256), `NitipKurs` b3882
//	GET /api/treaty-contract-out/tahun/{id}/kurs/konversi?dari=Rp|Usd&nilai=&skala= `HitungRpUsd_depan` / `CalculateTSIExcludeTreaty`
//
// Nomor baris = `Harness/InboxTreatyContractDescription.xml`.
//
// Dibaca sesudah: services/tco_kurs.go.

import (
	"net/http"

	"nusantarare/internal/services"
	"nusantarare/inti"
	"nusantarare/inti/galat"
)

func layananKursTCO(svc *services.Service) *services.KursTCO {
	return svc.KursTCO().
		DenganTahun(services.GudangTahunTreatyOracle(svc)).
		DenganMaster(services.MasterKursOracle(svc)).
		DenganMataUang(services.MataUangOracle(svc))
}

func daftarkanRuteKursTCO(mux *http.ServeMux, svc *services.Service, stub bool) {
	mux.HandleFunc("GET /api/treaty-contract-out/tahun/{id}/kurs", kursTahunTCO(svc, stub))
	mux.HandleFunc("GET /api/treaty-contract-out/tahun/{id}/kurs/konversi", konversiKursTCO(svc, stub))
}

func kursTahunTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		k, err := layananKursTCO(svc).KursTahun(r.Context(), inti.PelakuDari(r, stub), r.PathValue("id"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, k)
	}
}

func konversiKursTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		q := r.URL.Query()
		h, err := layananKursTCO(svc).Konversi(r.Context(), inti.PelakuDari(r, stub), r.PathValue("id"), q.Get("dari"),
			q.Get("nilai"), q.Get("skala"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		galat.TulisJSON(w, h)
	}
}
