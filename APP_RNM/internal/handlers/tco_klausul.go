package handlers

// Pintu HTTP klausul - tiket 08.
//
//	GET  /api/treaty-contract-out/jenis-klausul?isXol=0|1                 grid For Non XOL / For XOL (BrowseTreatyDesc_RD)
//	GET  /api/treaty-contract-out/klausul-pilihan/{master}?cari=          pemilih ExclutionTreaty (occupation|clause)
//	GET  /api/treaty-contract-out/tahun/{id}/klausul?descId=&induk=       grid per jenis (induk "00" / anak)
//	POST /api/treaty-contract-out/tahun/{id}/klausul                      Add + Save per jenis
//	PUT  /api/treaty-contract-out/tahun/{id}/klausul/{kid}                Edit + Save per jenis
//
// Dibaca sesudah: services/tco_klausul.go.

import (
	"encoding/json"
	"net/http"
	"strings"

	"nusantarare/internal/services"
)

func layananKlausulTCO(svc *services.Service) *services.KlausulTCO {
	return svc.KlausulTCO().
		DenganGudang(services.GudangKlausulOracle(svc)).
		DenganMaster(services.MasterKlausulOracle(svc)).
		DenganTahun(services.PengunciTahunOracle(svc)).
		DenganJenis(services.PembacaJenisReasuransiOracle(svc)).
		DenganKurs(services.PembacaKursOracle(svc))
}

func daftarkanRuteKlausulTCO(mux *http.ServeMux, svc *services.Service, stub bool) {
	mux.HandleFunc("GET /api/treaty-contract-out/jenis-klausul", jenisKlausulTCO(svc, stub))
	mux.HandleFunc("GET /api/treaty-contract-out/klausul-pilihan/{master}", pilihanKlausulTCO(svc, stub))
	mux.HandleFunc("GET /api/treaty-contract-out/tahun/{id}/klausul", daftarKlausulTCO(svc, stub))
	mux.HandleFunc("POST /api/treaty-contract-out/tahun/{id}/klausul", simpanKlausulTCO(svc, stub, false))
	mux.HandleFunc("PUT /api/treaty-contract-out/tahun/{id}/klausul/{kid}", simpanKlausulTCO(svc, stub, true))
}

type jawabanJenisKlausul struct {
	Daftar []services.JenisKlausulTampil `json:"daftar"`
	Total  int                           `json:"total"`
}

func jenisKlausulTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		d, err := layananKlausulTCO(svc).JenisKlausul(r.Context(), pelakuDari(r, stub), r.URL.Query().Get("isXol"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		tulisJSONPolis(w, jawabanJenisKlausul{Daftar: d, Total: len(d)})
	}
}

type jawabanPilihanKlausul struct {
	Daftar []services.PilihanTampil `json:"daftar"`
	Total  int                      `json:"total"`
}

func pilihanKlausulTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		d, err := layananKlausulTCO(svc).Pilihan(r.Context(), pelakuDari(r, stub), r.PathValue("master"),
			r.URL.Query().Get("cari"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		tulisJSONPolis(w, jawabanPilihanKlausul{Daftar: d, Total: len(d)})
	}
}

func daftarKlausulTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		q := r.URL.Query()
		if strings.TrimSpace(q.Get("descId")) == "" {
			galat(w, http.StatusBadRequest, "parameter descId (TreatyDescID) wajib")
			return
		}
		d, err := layananKlausulTCO(svc).Daftar(r.Context(), pelakuDari(r, stub), r.PathValue("id"), q.Get("descId"), q.Get("induk"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		tulisJSONPolis(w, d)
	}
}

func simpanKlausulTCO(svc *services.Service, stub, perbarui bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		var masuk services.KlausulMasuk
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&masuk); err != nil {
			galat(w, http.StatusBadRequest, "badan permintaan bukan JSON klausul yang sah")
			return
		}
		id := strings.TrimSpace(masuk.ID)
		switch {
		case !perbarui && id != "":
			galat(w, http.StatusBadRequest, "identitas klausul dibuat server; POST tidak boleh membawa id")
			return
		case perbarui && id != "" && id != r.PathValue("kid"):
			galat(w, http.StatusBadRequest, "id di badan berbeda dari id di jalur")
			return
		case perbarui:
			masuk.ID = r.PathValue("kid")
		}
		h, err := layananKlausulTCO(svc).Simpan(r.Context(), pelakuDari(r, stub), r.PathValue("id"), masuk)
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		tulisJSONPolis(w, h)
	}
}
