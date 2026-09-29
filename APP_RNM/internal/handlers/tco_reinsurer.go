package handlers

// Pintu HTTP reinsurer pada kombinasi kontrak - tiket 05.
//
//	GET  /api/treaty-contract-out/reinsurer-master?cari=...                    pemilih b8285 (BrowseAgentReinsSOA_RD)
//	GET  /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/reinsurer           Reinsurer List b11308
//	POST /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/reinsurer           Add b1730 + Save b11405
//	PUT  /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/reinsurer/{rid}     Edit b4491 + Save b11405
//
// Nomor baris = `Section/ViewDetailTreatyReinsurerGrid1.xml` kecuali disebut lain.
//
// Dibaca sesudah: tco_kontrak.go, services/tco_reinsurer.go.

import (
	"encoding/json"
	"net/http"
	"strings"

	"nusantarare/internal/services"
)

func layananReinsurerTCO(svc *services.Service) *services.ReinsurerTCO {
	return svc.ReinsurerTCO().
		DenganGudang(services.GudangReinsurerOracle(svc)).
		DenganKontrak(services.PemegangKontrakOracle(svc)).
		DenganTahun(services.GudangTahunTreatyOracle(svc)).
		DenganMaster(services.MasterReinsurerOracle(svc))
}

func daftarkanRuteReinsurerTCO(mux *http.ServeMux, svc *services.Service, stub bool) {
	mux.HandleFunc("GET /api/treaty-contract-out/reinsurer-master", cariReinsurerMasterTCO(svc, stub))
	mux.HandleFunc("GET /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/reinsurer", daftarReinsurerTCO(svc, stub))
	mux.HandleFunc("POST /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/reinsurer", simpanReinsurerTCO(svc, stub, false))
	mux.HandleFunc("PUT /api/treaty-contract-out/tahun/{id}/kontrak/{kid}/reinsurer/{rid}", simpanReinsurerTCO(svc, stub, true))
}

type jawabanReinsurerMaster struct {
	Daftar []services.ReinsurerMasterTampil `json:"daftar"`
	Total  int                              `json:"total"`
}

func cariReinsurerMasterTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		d, err := layananReinsurerTCO(svc).CariMaster(r.Context(), pelakuDari(r, stub), r.URL.Query().Get("cari"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		tulisJSONPolis(w, jawabanReinsurerMaster{Daftar: d, Total: len(d)})
	}
}

func daftarReinsurerTCO(svc *services.Service, stub bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		d, err := layananReinsurerTCO(svc).Daftar(r.Context(), pelakuDari(r, stub), r.PathValue("id"), r.PathValue("kid"))
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		tulisJSONPolis(w, d)
	}
}

// simpanReinsurerTCO - POST menolak `id` dari klien (AC 5); PUT menolak `id`
// badan yang berbeda dari jalur.
func simpanReinsurerTCO(svc *services.Service, stub, perbarui bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !punyaDBTCO(w, svc) {
			return
		}
		var masuk services.ReinsurerMasuk
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&masuk); err != nil {
			galat(w, http.StatusBadRequest, "badan permintaan bukan JSON reinsurer yang sah")
			return
		}
		id := strings.TrimSpace(masuk.ID)
		switch {
		case !perbarui && id != "":
			galat(w, http.StatusBadRequest, "identitas reinsurer dibuat server; POST tidak boleh membawa id")
			return
		case perbarui && id != "" && id != r.PathValue("rid"):
			galat(w, http.StatusBadRequest, "id di badan berbeda dari id di jalur")
			return
		case perbarui:
			masuk.ID = r.PathValue("rid")
		}
		h, err := layananReinsurerTCO(svc).Simpan(r.Context(), pelakuDari(r, stub), r.PathValue("id"),
			r.PathValue("kid"), masuk)
		if jawabGalatTreatyContractOut(w, err) {
			return
		}
		tulisJSONPolis(w, h)
	}
}
